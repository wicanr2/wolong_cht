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
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
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
//（左翼 左備 主將 前鋒 右備 右翼），不是記錄順序。
func squadSlot(i int) int { return battle.BottomSlotSquad[i] }

// battleState 是手機端的一場戰術戰鬥。
type battleState struct {
	view *isoview.View

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
	return isoview.New(isoview.Options{
		Lib:     s.setup.Library(),
		Palette: bank,
		Sprites: s.setup.Sprites(),
		Field:   s.setup.FieldNumber(p.Node, p.Mode == combat.Siege),
		// 翻轉要用**建場當時**那個值，不能重算——重算會在玩家守城時
		// 得到相反的答案，旗與地形就對不上（docs/playtest/40 §11）。
		Rotate: s.setup.Rotate(),
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
	// BattleRowH 是六格與六命令那兩列的高度。
	BattleRowH = 56
)

// BattleFieldRect 是戰場畫面的可視區。
func BattleFieldRect() (x, y, w, h int) {
	return 0, StatusH, LogicalW, LogicalH - StatusH - BattleRowH*2
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
	fx, fy, fw, fh := BattleFieldRect()
	fillRect(dst, fx, fy, fw, fh, inkVoid())

	if s.battle.view != nil {
		buf := s.battle.view.Render(b)
		// 等比放大到塞得進戰場區為止。**只用整數倍**——
		// 原版是點陣圖，非整數倍會糊（§3 的同一條理由）。
		scale := 1
		for (scale+1)*isoview.NativeW <= fw && (scale+1)*isoview.NativeH <= fh {
			scale++
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(scale), float64(scale))
		op.GeoM.Translate(float64(fx+(fw-scale*isoview.NativeW)/2),
			float64(fy+(fh-scale*isoview.NativeH)/2))
		dst.SubImage(rect(fx, fy, fw, fh)).(*ebiten.Image).DrawImage(buf, op)
	}
	s.drawBattleSides(dst, td, b, fx, fy, fw, fh)
	s.drawSquadStrip(dst, td, b)
	s.drawBattleCommands(dst, td)
}

// drawBattleSides 用戰場左右剩下的兩塊放原版側欄有的東西：
// 兩軍的主將與兵力、戰場縮圖。
//
// ⚠ 原版的側欄在**右邊一整條**（docs/spec/31）。手機是橫向 16:9，
// 480 px 的戰場置中之後左右各剩一塊，擺兩邊比擠成一條好讀——
// **這是 remake 差異**，內容仍照原版那幾個欄位。
func (s *Session) drawBattleSides(dst *ebiten.Image, td *textdraw.Drawer,
	b *tactical.Battle, fx, fy, fw, fh int) {

	margin := (fw - isoview.NativeW) / 2
	if margin < 120 || td == nil || !td.Available() {
		return
	}
	p := s.world.PendingBattle()
	if p == nil {
		return
	}
	name := func(corps int) string {
		if corps < 0 || corps >= len(s.world.Generals) {
			return "城兵"
		}
		return big5(s.world.Generals[corps].Name)
	}
	// 左：攻守雙方。玩家那一側標出來——原版靠側欄換邊表達（docs/spec/56）。
	rows := [][2]string{
		{"攻", name(p.Attacker)},
		{"", fmt.Sprintf("%d", b.Sides[tactical.AttackerSide].Alive()*tactical.MenPerSoldier)},
		{"守", name(p.Defender)},
		{"", fmt.Sprintf("%d", b.Sides[tactical.DefenderSide].Alive()*tactical.MenPerSoldier)},
	}
	for i, r := range rows {
		y := fy + 24 + i*30
		// 玩家那一側用色 15（白），對方用次要色——**不要用反白條的色 5**：
		// 那是清單選取用的綠，拿來標「哪一側是我」會讀成「選中了它」。
		ink := inkDim()
		if (i < 2) == (s.playerSide() == tactical.AttackerSide) {
			ink = inkText()
		}
		if r[0] != "" {
			td.Draw(dst, r[0], fx+16, y, inkDim())
		}
		td.Draw(dst, r[1], fx+48, y, ink)
	}
	// 右：戰場縮圖。原版點它可以移動鏡頭；手機版目前只顯示。
	if mm := s.battle.view.Minimap(); mm != nil {
		w, h := mm.Bounds().Dx(), mm.Bounds().Dy()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(fx+fw-margin+(margin-w)/2), float64(fy+(fh-h)/2))
		dst.DrawImage(mm, op)
	}
}

func (s *Session) drawSquadStrip(dst *ebiten.Image, td *textdraw.Drawer, b *tactical.Battle) {
	side := s.playerSide()
	for i := 0; i < army.Positions; i++ {
		x, y, w, h := SquadRect(i)
		// 選中的那一格用**反白條**（色 5），其餘是命令列的黑底——
		// 原版底列六格的選取也是靠框與反白標的（docs/spec/33 §1.1）。
		bg, ink := inkBar(), inkText()
		if squadSlot(i) == s.battle.squad {
			bg, ink = inkSelect(), inkInk()
		}
		s.window(dst, x, y, w, h, bg)
		if td == nil || !td.Available() {
			continue
		}
		label := unitLabel(squadSlot(i))
		td.Draw(dst, label, x+(w-td.Width(label))/2, y+6, ink)
		men := fmt.Sprintf("%d", squadMen(b, side, squadSlot(i)))
		td.Draw(dst, men, x+(w-td.Width(men))/2, y+h-24, ink)
	}
}

func (s *Session) drawBattleCommands(dst *ebiten.Image, td *textdraw.Drawer) {
	for i := range battle.SideCommandRowCode {
		c := battleCommandRow(i)
		x, y, w, h := BattleCommandRect(i)
		s.window(dst, x, y, w, h, inkBar())
		if td == nil || !td.Available() {
			continue
		}
		label := c.String()
		td.Draw(dst, label, x+(w-td.Width(label))/2, y+(h-16)/2, inkText())
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
	if s.world.PendingEncounter() != nil {
		if err := s.world.ChooseBattleCommand(); err != nil {
			return err
		}
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
