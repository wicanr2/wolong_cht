package phone

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/speed"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/ui/isoview"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 戰場。畫面版面是手機自己的（docs/mobile/android-ux.md §5），
// 但**等角投影與規則層都與桌面版共用**：`internal/ui/isoview`、
// `internal/rules/tactical`。

// battleCommandRow 是命令列由左到右第 i 個送出的命令碼。
//
// ⚠ **畫面順序不是命令碼順序**（第 0 個送命令 2）。這張表是原版資料，
// 與桌面版共用（`battle.SideCommandRowCode`）；抄一份就會把命令送錯。
func battleCommandRow(i int) tactical.Command {
	return tactical.Command(battle.SideCommandRowCode[i])
}

// squadSlot 是底列第 i 格對應哪一個編成位置。原版底列是**空間排列**
// （左翼 左備 主將 前鋒 右備 右翼），不是記錄順序。
func squadSlot(i int) int { return battle.BottomSlotSquad[i] }

// battleState 是手機端的一場戰術戰鬥。
type battleState struct {
	view             *isoview.View
	orderIcons       [len(battle.SideCommandRowCode)]*ebiten.Image
	orderIconsLoaded bool

	// squad 是選中的編成位置（0–5）。**選中是視覺狀態**：
	// 下令要再點一次底列的命令（兩段式，docs/mobile/android-ux.md §3）。
	squad int

	// throttle 是戰術層的速度節流。戰略與戰術**各有一個獨立的檔位**，
	// 這是這一款的特性（docs/spec/34）。
	throttle speed.Throttle
	tacSpeed int
}

// BattleActive 回報現在是不是在戰場上。
func (s *Session) BattleActive() bool { return s.world.PendingBattle() != nil }

// Battle 回傳等著打的那一場，沒有就回 nil。
func (s *Session) Battle() *tactical.Battle {
	p := s.world.PendingBattle()
	if p == nil {
		return nil
	}
	return p.Battle
}

// tickBattle 推進戰場。
//
// ⚠ **戰術層沒有暫停**（說明書 4.1：戦闘中は絶対に時間を止められません）。
// 手機版照這條走，不因為開著選單就停。
func (s *Session) tickBattle() {
	b := s.Battle()
	if b == nil {
		return
	}
	if s.battle.view == nil {
		s.battle = battleState{view: s.newBattleView(), tacSpeed: DefaultSpeed}
		if s.battle.view != nil {
			s.battle.view.Observe(b)
		}
	}
	if s.paused {
		return // 只有驗收路徑會停，那不是遊戲內的暫停
	}
	n := s.battle.throttle.Steps(s.battle.tacSpeed, speed.TacticalMul, speed.HighSpeedTactical)
	for ; n > 0; n-- {
		b.Step()
	}
	if b.Done {
		s.finishBattle()
	}
}

// finishBattle 把打完的結果交回戰略層，並丟掉這一場的繪圖資源。
func (s *Session) finishBattle() {
	s.world.ResolvePending(s.rand)
	s.battle = battleState{tacSpeed: DefaultSpeed}
}

func (s *Session) newBattleView() *isoview.View {
	p := s.world.PendingBattle()
	if p == nil || s.setup == nil || s.lib == nil || s.lib.Palette == nil {
		return nil
	}
	bank, err := s.lib.Palette.Bank(0)
	if err != nil {
		return nil
	}
	lib := s.setup.Library()
	if lib == nil {
		return nil
	}
	return isoview.New(isoview.Options{
		Lib:     lib,
		Palette: bank,
		Sprites: s.setup.Sprites(),
		Field:   p.FieldNumber,
		// 翻轉要用**建場當時**那個值，不能重算——重算會在玩家守城時
		// 得到相反的答案，旗與地形就對不上（docs/playtest/40 §11）。
		Rotate: p.FieldRotate,
		Rand:   s.rand.Next,
	})
}

// 戰場的版面。
//
// ⭐ **戰場上不畫四個入口那一條**：進言／一覽／軍團／系統在戰術畫面上
// 本來就不作用（原版進戰術畫面時整個畫面都換掉）。省下那 64 px 之後
// 原版的 480×368 視野才塞得進來——不然戰場會被裁掉一截，
// 而「看得到多少戰場」是會影響決策的。
const (
	battleHeaderH = 32
	// BattleRowH 是兩排的觸控高度，可見外框在格內置中。
	BattleRowH        = 48
	battleFaceH       = 33
	battleFaceInset   = 6
	battleMinimapSize = 96
)

// BattleFieldRect 是戰場畫面的可視區。
func BattleFieldRect() (x, y, w, h int) {
	return 0, battleHeaderH, LogicalW, LogicalH - battleHeaderH - BattleRowH*2
}

func battleViewport() image.Rectangle {
	x, y, w, h := BattleFieldRect()
	viewportH := min(h, w*isoview.NativeH/isoview.NativeW)
	viewportW := viewportH * isoview.NativeW / isoview.NativeH
	left, top := x+(w-viewportW)/2, y+(h-viewportH)/2
	return image.Rect(left, top, left+viewportW, top+viewportH)
}

// BattleViewportContains reports whether a touch is on the live battlefield.
func (s *Session) BattleViewportContains(x, y int) bool {
	r := battleViewport()
	return s.BattleActive() && s.battle.view != nil && image.Pt(x, y).In(r)
}

// PinchZoom applies one level of zoom in the scene where the gesture began.
func (s *Session) PinchZoom(direction, x, y int) {
	if !s.BattleActive() {
		s.SetZoom(s.zoom + direction)
		return
	}
	if s.BattleViewportContains(x, y) {
		r := battleViewport()
		v := s.battle.view
		focusX := (x - r.Min.X) * isoview.NativeW / r.Dx()
		focusY := (y - r.Min.Y) * isoview.NativeH / r.Dy()
		v.ZoomAt(v.ZoomPercent()+25*direction, focusX, focusY)
	}
}

// SquadRect 是上排六格的第 i 格。
func SquadRect(i int) (x, y, w, h int) {
	cell := LogicalW / army.Positions
	return i * cell, LogicalH - BattleRowH*2, cell, BattleRowH
}

// BattleCommandRect 是底列六個命令的第 i 個。
func BattleCommandRect(i int) (x, y, w, h int) {
	cell := LogicalW / len(battle.SideCommandRowCode)
	return i * cell, LogicalH - BattleRowH, cell, BattleRowH
}

// tapBattle 處理戰場上的一次點擊。
func (s *Session) tapBattle(lx, ly float64) bool {
	for i := 0; i < army.Positions; i++ {
		if x, y, w, h := SquadRect(i); inRect(lx, ly, x, y, w, h) {
			s.battle.squad = squadSlot(i)
			return true
		}
	}
	for i := range battle.SideCommandRowCode {
		if x, y, w, h := BattleCommandRect(i); inRect(lx, ly, x, y, w, h) {
			s.orderSquad(battleCommandRow(i))
			return true
		}
	}
	return true // 戰場上點別的地方不做事，但也不要穿透到底下的大地圖
}

func inRect(lx, ly float64, x, y, w, h int) bool {
	return lx >= float64(x) && lx < float64(x+w) &&
		ly >= float64(y) && ly < float64(y+h)
}

// orderSquad 對選中的那一隊下令。
func (s *Session) orderSquad(c tactical.Command) {
	b := s.Battle()
	if b == nil {
		return
	}
	b.Order(s.playerSide(), s.battle.squad, c)
}

// playerSide 回報玩家在這一場是攻方（0）還是守方（1）。
//
// ⚠ **側別是攻守不是玩家**（`internal/rules/tactical` 的約定，
// 與原版相反），所以要從軍團的勢力反查。
func (s *Session) playerSide() int {
	p := s.world.PendingBattle()
	if p == nil {
		return tactical.AttackerSide
	}
	if p.Attacker >= 0 && p.Attacker < len(s.world.Corps) &&
		s.world.Corps[p.Attacker].Faction == s.world.Player {
		return tactical.AttackerSide
	}
	return tactical.DefenderSide
}

// drawBattle 畫戰場。
func (s *Session) drawBattle(dst *ebiten.Image, td *textdraw.Drawer) {
	b := s.Battle()
	if b == nil {
		return
	}
	td.SetScale(1)
	fx, fy, fw, fh := BattleFieldRect()
	fillRect(dst, fx, fy, fw, fh, inkVoid())

	if s.battle.view != nil {
		buf := s.battle.view.Render(b)
		// Render returns the visible source crop, whose size changes with zoom.
		// Project it into the same destination rectangle at every zoom level.
		viewport := battleViewport()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(viewport.Dx())/float64(buf.Bounds().Dx()),
			float64(viewport.Dy())/float64(buf.Bounds().Dy()))
		op.GeoM.Translate(float64(viewport.Min.X), float64(viewport.Min.Y))
		dst.SubImage(viewport).(*ebiten.Image).DrawImage(buf, op)
	}
	s.drawBattleHeader(dst, td, b)
	s.drawBattleMinimap(dst, fy)
	fillRect(dst, 0, LogicalH-BattleRowH*2, LogicalW, BattleRowH*2, color.RGBA{26, 25, 20, 255})
	fillRect(dst, 0, LogicalH-BattleRowH*2, LogicalW, 1, color.RGBA{146, 83, 71, 255})
	s.drawSquadStrip(dst, td, b)
	s.drawBattleCommands(dst, td)
}

func (s *Session) drawBattleHeader(dst *ebiten.Image, td *textdraw.Drawer, b *tactical.Battle) {
	fillRect(dst, 0, 0, LogicalW, battleHeaderH, color.RGBA{27, 27, 22, 255})
	fillRect(dst, 0, battleHeaderH-1, LogicalW, 1, color.RGBA{146, 83, 71, 255})
	p := s.world.PendingBattle()
	if p == nil || !td.Available() {
		return
	}
	for side, corps := range [2]int{p.Attacker, p.Defender} {
		name := "城兵"
		if corps >= 0 && corps < len(s.world.Generals) {
			name = s.Localise(s.world.Generals[corps].Name)
		}
		label := [2]string{"攻", "守"}[side]
		label += " " + name + fmt.Sprintf(" %d", b.Sides[side].Alive()*tactical.MenPerSoldier)
		x := 14
		if side == tactical.DefenderSide {
			x = LogicalW - 14 - td.Width(label)
		}
		ink := color.RGBA{221, 186, 118, 255}
		if side == s.playerSide() {
			ink = battleButtonInk
		}
		td.Draw(dst, label, x, (battleHeaderH-td.FontH())/2, ink)
	}
}

func (s *Session) drawBattleMinimap(dst *ebiten.Image, y int) {
	if s.battle.view == nil {
		return
	}
	if mm := s.battle.view.Minimap(); mm != nil {
		bounds := mm.Bounds()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(battleMinimapSize)/float64(bounds.Dx()), float64(battleMinimapSize)/float64(bounds.Dy()))
		x, top := LogicalW-16-battleMinimapSize, y+13
		op.GeoM.Translate(float64(x), float64(top))
		dst.DrawImage(mm, op)
		strokeRect(dst, x-1, top-1, battleMinimapSize+2, battleMinimapSize+2, color.RGBA{191, 167, 123, 255})
	}
}

var battleButtonInk = color.RGBA{241, 230, 203, 255}

func battleButtonFace(x, y, w, h int) image.Rectangle {
	top := y + (h-battleFaceH)/2
	return image.Rect(x+battleFaceInset, top, x+w-battleFaceInset, top+battleFaceH)
}

func drawBattleFace(dst *ebiten.Image, r image.Rectangle, selected, retreat bool) {
	fill, border := color.RGBA{34, 34, 27, 255}, color.RGBA{81, 67, 51, 255}
	if selected {
		fill, border = color.RGBA{119, 47, 40, 255}, color.RGBA{220, 144, 116, 255}
	} else if retreat {
		border = color.RGBA{169, 107, 91, 255}
	}
	fillRect(dst, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), fill)
	strokeRect(dst, r.Min.X, r.Min.Y, r.Dx(), r.Dy(), border)
}

func (s *Session) drawSquadStrip(dst *ebiten.Image, td *textdraw.Drawer, b *tactical.Battle) {
	side := s.playerSide()
	for i := 0; i < army.Positions; i++ {
		face := battleButtonFace(SquadRect(i))
		squad := squadSlot(i)
		drawBattleFace(dst, face, squad == s.battle.squad, false)
		if !td.Available() {
			continue
		}
		label := unitLabel(squad) + " " + fmt.Sprintf("%d", squadMen(b, side, squad))
		td.Draw(dst, label, face.Min.X+(face.Dx()-td.Width(label))/2, face.Min.Y+(face.Dy()-td.FontH())/2, battleButtonInk)
	}
}

func (s *Session) drawBattleCommands(dst *ebiten.Image, td *textdraw.Drawer) {
	if !s.battle.orderIconsLoaded {
		for i, code := range battle.SideCommandRowCode {
			if icon, err := s.lib.DOSVOrderIcon(code, 0); err == nil {
				s.battle.orderIcons[i] = ebiten.NewImageFromImage(icon)
			}
		}
		s.battle.orderIconsLoaded = true
	}
	for i := range battle.SideCommandRowCode {
		c := battleCommandRow(i)
		face := battleButtonFace(BattleCommandRect(i))
		drawBattleFace(dst, face, false, c == tactical.Retreat)
		label := c.String()
		width := 0
		if td.Available() {
			width = td.Width(label)
		}
		icon := s.battle.orderIcons[i]
		if icon != nil {
			width += icon.Bounds().Dx()
			if td.Available() {
				width += 5
			}
		}
		x := face.Min.X + (face.Dx()-width)/2
		if icon != nil {
			bounds := icon.Bounds()
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(face.Min.Y+(face.Dy()-bounds.Dy())/2))
			dst.DrawImage(icon, op)
			x += bounds.Dx() + 5
		}
		if td.Available() {
			td.Draw(dst, label, x, face.Min.Y+(face.Dy()-td.FontH())/2, battleButtonInk)
		}
	}
}

// squadMen 是一隊還活著的兵數（人）。一個兵代表 10 人（說明書 4.1）。
func squadMen(b *tactical.Battle, side, squad int) int {
	n := 0
	lo, hi := squad*tactical.PerSquad, (squad+1)*tactical.PerSquad
	for i := lo; i < hi && i < len(b.Sides[side].Soldiers); i++ {
		if b.Sides[side].Soldiers[i].Alive {
			n++
		}
	}
	return n * tactical.MenPerSoldier
}

// fillRect／strokeRect 是版面用的兩個小工具，讓畫面碼讀起來短一點。
func fillRect(dst *ebiten.Image, x, y, w, h int, c color.RGBA) {
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h), c, false)
}

func strokeRect(dst *ebiten.Image, x, y, w, h int, c color.RGBA) {
	vector.StrokeRect(dst, float32(x), float32(y), float32(w), float32(h), 1, c, false)
}

// rect 是 image.Rect 的短寫法（SubImage 要用）。
func rect(x, y, w, h int) image.Rectangle { return image.Rect(x, y, x+w, y+h) }

// OpenDemoBattle 是**驗收鉤子**：直接擺一場攻城戰出來。
//
// 擺位走 `internal/battlesetup`，與桌面版的 `-open-siege` 同一支——
// 兩邊的驗收要能開在同一個局面，否則畫面比不了。
func (s *Session) OpenDemoBattle(node int) error {
	att, def := -1, -1
	for i := range s.world.Corps {
		if !s.world.Corps[i].Alive {
			continue
		}
		if s.world.Corps[i].Faction == s.world.Player {
			def = i
			continue
		}
		att = i
	}
	if att < 0 || def < 0 {
		// 開局沒有現成軍團就自己編兩支：一支玩家的、一支鄰國的。
		var err error
		if att, def, err = s.formDemoCorps(); err != nil {
			return err
		}
	}
	if !battlesetup.StageEncounter(s.world, s.rand, battlesetup.StageOptions{
		Siege: true, Node: node, Attacker: att, Defender: def,
	}) {
		return fmt.Errorf("擺不出戰鬥（攻 %d 守 %d）", att, def)
	}
	s.battle = battleState{view: s.newBattleView(), tacSpeed: DefaultSpeed}
	return nil
}

// formDemoCorps 編兩支軍團給驗收用：玩家一支、另一個勢力一支。
func (s *Session) formDemoCorps() (att, def int, err error) {
	pick := func(faction int) (int, error) {
		for i := range s.world.Generals {
			g := &s.world.Generals[i]
			if !g.Alive || g.Faction != faction || s.world.Corps[i].Alive {
				continue
			}
			f := newCorpsForm()
			if e := s.world.FormCorps(i, f.kinds, f.manned); e == nil {
				return i, nil
			}
		}
		return -1, fmt.Errorf("勢力 %d 編不出軍團", faction)
	}
	if def, err = pick(s.world.Player); err != nil {
		return -1, -1, err
	}
	for i := range s.world.Factions {
		if i == s.world.Player || !s.world.Factions[i].Alive {
			continue
		}
		if att, err = pick(i); err == nil {
			return att, def, nil
		}
	}
	return -1, -1, fmt.Errorf("找不到第二個勢力的軍團")
}
