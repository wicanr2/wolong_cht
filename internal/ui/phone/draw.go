package phone

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 手機版的配色。⚠ **這不是原版調色盤**——原版的介面色是從
// `GAMEPAL.BRG` 查出來的（docs/spec/54），那一套綁在 640×400 的版面上。
// 這裡是重畫的外殼，用自己的色票，並且刻意壓低彩度讓原版美術站在前面。
var (
	inkBar    = color.RGBA{16, 22, 38, 255}
	inkPanel  = color.RGBA{24, 33, 56, 255}
	inkEdge   = color.RGBA{92, 116, 168, 255}
	inkText   = color.RGBA{232, 238, 250, 255}
	inkDim    = color.RGBA{150, 168, 205, 255}
	inkSelect = color.RGBA{255, 214, 102, 255}
	inkVoid   = color.RGBA{8, 10, 18, 255}
	// inkOverlay 是擋住世界的決定壓在地圖上那一層。**半透明**：
	// 地圖還看得見，玩家才知道這個決定發生在哪。
	inkOverlay = color.RGBA{8, 10, 18, 210}
)

// Draw 把一局畫進 960×540 的邏輯畫布。
func (s *Session) Draw(dst *ebiten.Image, td *textdraw.Drawer) {
	dst.Fill(inkVoid)
	// ⭐ 戰場開著時**整個主區換成戰場**，大地圖不畫——
	// 原版進戰術畫面時戰略畫面也整個換掉。
	if s.BattleActive() {
		s.drawBattle(dst, td)
		s.drawStatusBar(dst, td)
		return
	}
	s.drawMap(dst)
	// 擋住世界的決定優先：不選它時間不會走（modal.go）。
	if s.ModalKind() != modalNone {
		s.drawModal(dst, td)
		s.drawStatusBar(dst, td)
		s.drawCommandBar(dst, td)
		return
	}
	s.drawNotice(dst, td)
	switch {
	case s.advise.stage != adviseIdle:
		s.drawAdvise(dst, td)
	case s.sheet.open:
		s.drawSheet(dst, td)
	case s.selected >= 0:
		s.drawCityCard(dst, td)
	}
	s.drawStatusBar(dst, td)
	s.drawCommandBar(dst, td)
}

// drawSheet 畫指令列打開的面板。
func (s *Session) drawSheet(dst *ebiten.Image, td *textdraw.Drawer) {
	mx, my, mw, mh := MapRect()
	vector.DrawFilledRect(dst, float32(mx), float32(my), float32(mw), float32(mh), inkPanel, false)
	s.drawTabs(dst, td)
	rows := s.sheetRows()
	if td == nil || !td.Available() {
		return
	}
	top := my + tabH
	visible := (mh - tabH) / rowH
	for i := 0; i < visible && i+s.sheet.scroll < len(rows); i++ {
		r := rows[i+s.sheet.scroll]
		y := top + i*rowH
		if (i+s.sheet.scroll)%2 == 1 {
			vector.DrawFilledRect(dst, float32(mx), float32(y),
				float32(mw), float32(rowH), inkBar, false)
		}
		ink := inkText
		if r.dim {
			ink = inkDim
		}
		td.Draw(dst, r.name, mx+sheetPadX, y+sheetHeadDY, ink)
		// 欄位從右往左排：欄數不固定，靠右對齊才不會因為某一頁少一欄
		// 就整排跟著位移。
		x := mx + mw - sheetPadX
		for j := len(r.cols) - 1; j >= 0; j-- {
			c := r.cols[j]
			if c == "" {
				continue
			}
			x -= td.Width(c)
			td.Draw(dst, c, x, y+sheetHeadDY, inkDim)
			x -= rowTextDX * 2
		}
	}
	if f := s.sheetFooter(); f != "" {
		td.Draw(dst, f, mx+sheetPadX, my+mh-rowH, inkDim)
	}
	if s.lastErr != nil {
		td.Draw(dst, s.lastErr.Error(), mx+sheetPadX, my+mh-rowH*2, inkSelect)
	}
}

func (s *Session) drawTabs(dst *ebiten.Image, td *textdraw.Drawer) {
	tabs := s.Tabs()
	mx, my, mw, _ := MapRect()
	if len(tabs) == 0 {
		// 沒有分頁的面板仍然要有一條標題列，否則第一列會貼著狀態列，
		// 看起來像地圖的一部分。
		vector.DrawFilledRect(dst, float32(mx), float32(my), float32(mw), tabH, inkBar, false)
		if td != nil && td.Available() && s.sheet.cmd >= 0 {
			td.Draw(dst, s.sheet.cmd.Label(), mx+sheetPadX, my+sheetHeadDY+2, inkText)
		}
		return
	}
	cell := mw / len(tabs)
	for i, t := range tabs {
		x := mx + i*cell
		bg := inkBar
		if i == s.sheet.tab {
			bg = inkPanel
		}
		vector.DrawFilledRect(dst, float32(x), float32(my), float32(cell), tabH, bg, false)
		if i == s.sheet.tab {
			vector.DrawFilledRect(dst, float32(x), float32(my+tabH-3),
				float32(cell), 3, inkSelect, false)
		}
		if td == nil || !td.Available() {
			continue
		}
		ink := inkDim
		if i == s.sheet.tab {
			ink = inkText
		}
		td.Draw(dst, t, x+(cell-td.Width(t))/2, my+sheetHeadDY+2, ink)
	}
}

// drawAdvise 畫進言：上面是對白，下面是這一刻要選的那幾列。
//
// ⚠ 原版是**上下兩個講話框配肖像**（docs/spec/45）。手機的主區放不下
// 兩個框加肖像，改成一列一句、靠邊區分誰在講——**這是 remake 差異**，
// 記在 docs/mobile/android-ux.md §7。
func (s *Session) drawAdvise(dst *ebiten.Image, td *textdraw.Drawer) {
	mx, my, mw, mh := MapRect()
	vector.DrawFilledRect(dst, float32(mx), float32(my), float32(mw), float32(mh), inkPanel, false)
	vector.DrawFilledRect(dst, float32(mx), float32(my), float32(mw), tabH, inkBar, false)
	if td == nil || !td.Available() {
		return
	}
	td.Draw(dst, s.adviseTitle(), mx+sheetPadX, my+sheetHeadDY+2, inkText)

	if s.advise.stage == advisePickAlly || s.advise.stage == advisePickTarget {
		// 對象可能有二十幾個，底部的選項條放不下——用可捲的清單。
		choices := s.AdviseChoices()
		top := my + tabH
		visible := (mh - tabH) / rowH
		for i := 0; i < visible && i+s.sheet.scroll < len(choices); i++ {
			y := top + i*rowH
			if (i+s.sheet.scroll)%2 == 1 {
				vector.DrawFilledRect(dst, float32(mx), float32(y),
					float32(mw), float32(rowH), inkBar, false)
			}
			td.Draw(dst, choices[i+s.sheet.scroll], mx+sheetPadX, y+sheetHeadDY, inkText)
		}
		return
	}

	choices := s.AdviseChoices()
	bottom := my + mh - len(choices)*rowH

	// 對白只留最後幾句——面板高度固定，講太多就從上面捲掉。
	said := s.advise.said
	if n := (bottom - my - tabH) / rowH; n > 0 && len(said) > n {
		said = said[len(said)-n:]
	}
	for i, l := range said {
		y := my + tabH + sheetHeadDY + i*rowH
		if l.lord {
			td.Draw(dst, l.text, mx+sheetPadX, y, inkSelect)
			continue
		}
		// 軍師（玩家）的話靠右，與君主分開。
		td.Draw(dst, l.text, mx+mw-sheetPadX-td.Width(l.text), y, inkText)
	}
	for i, c := range choices {
		y := bottom + i*rowH
		vector.DrawFilledRect(dst, float32(mx), float32(y), float32(mw), float32(rowH-2), inkBar, false)
		td.Draw(dst, c, mx+sheetPadX, y+4, inkText)
	}
	if len(choices) == 0 {
		td.Draw(dst, "點畫面繼續", mx+mw-sheetPadX-td.Width("點畫面繼續"),
			my+mh-rowH, inkDim)
	}
}

// adviseTitle 說現在這一頁在做什麼。**選對象與說服長得很像**
//（都是一排可點的列），沒有標題會分不出來。
func (s *Session) adviseTitle() string {
	name := adviseFallbackNames[s.advise.row]
	if labels := s.adviseLabels(); s.advise.row < len(labels) {
		name = labels[s.advise.row]
	}
	switch s.advise.stage {
	case advisePickAlly:
		return name + "：選協力方"
	case advisePickTarget:
		return name + "：選對象"
	case advisePersuade:
		return name + "：說服君主"
	}
	return name
}

func (s *Session) drawMap(dst *ebiten.Image) {
	mx, my, _, _ := MapRect()
	cols, rows := s.viewTiles()
	// 多畫一格，右下角才不會在鏡頭落在半格時露出底色。
	img, err := s.lib.RenderWorldMarked(s.camX, s.camY, cols+1, rows+1,
		int(s.world.Clock.Season()), s.cityMarks())
	if err != nil {
		return
	}
	tex := ebiten.NewImageFromImage(img)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(s.zoom), float64(s.zoom))
	op.GeoM.Translate(float64(mx), float64(my))
	// 多畫的那一格會落到上下兩條底下，而那兩條是不透明的——
	// 所以不必為了裁切另外開 SubImage。
	dst.DrawImage(tex, op)

	if s.selected >= 0 {
		s.drawSelectionRing(dst)
	}
}

// drawSelectionRing 在選中的據點外圍畫一圈。**選中是視覺狀態，不是指令**。
//
// ⚠ **一定要畫兩圈**：外圈深色、內圈亮色。原版的據點圖形是土黃色的，
// 單畫一圈亮金色疊上去**在城上完全看不出來**——不是沒畫，是沒對比。
// 深色墊底之後在草地、土黃城與河面上都看得見。
func (s *Session) drawSelectionRing(dst *ebiten.Image) {
	c := &s.world.Cities[s.selected]
	mx, my, _, _ := MapRect()
	px := float32(TilePx * s.zoom)
	// ⭐ 圈要**比城大**：大城的圖形是 5×5 格（`world.applyDecor` 的
	// 222–225 是 5×5、226–229 是 3×3），畫 3 格的圈會落在城的**內部**，
	// 看起來像城的一部分而不是選取狀態。7 格才框得住最大的城還留一格邊。
	const ringTiles = 7
	half := float32(ringTiles / 2)
	x := float32(mx) + (float32(c.X+world.CityCentreDX-s.camX)-half)*px
	y := float32(my) + (float32(c.Y-s.camY)-half)*px
	side := px * ringTiles
	vector.StrokeRect(dst, x-1, y-1, side+2, side+2, 5, inkVoid, false)
	vector.StrokeRect(dst, x, y, side, side, 3, inkSelect, false)
}

func (s *Session) cityMarks() []world.CityMark {
	marks := make([]world.CityMark, 0, len(s.world.Cities))
	capital := -1
	if p := s.world.Player; p >= 0 && p < len(s.world.Factions) {
		capital = s.world.Factions[p].Capital
	}
	for i := range s.world.Cities {
		c := &s.world.Cities[i]
		marks = append(marks, world.CityMark{
			X: c.X + world.CityCentreDX, Y: c.Y,
			Own:     world.OwnershipOf(c.Owner, s.world.Player),
			Capital: i == capital,
		})
	}
	return marks
}

func (s *Session) drawStatusBar(dst *ebiten.Image, td *textdraw.Drawer) {
	vector.DrawFilledRect(dst, 0, 0, LogicalW, StatusH, inkBar, false)
	vector.StrokeLine(dst, 0, StatusH, LogicalW, StatusH, 1, inkEdge, false)
	if td == nil || !td.Available() {
		return
	}
	c := s.world.Clock
	td.Draw(dst, fmt.Sprintf("%d年%d月%d日", c.Year, c.Month, c.Day), 16, 18, inkText)

	p := s.world.Player
	if p < 0 || p >= len(s.world.Factions) {
		return
	}
	f := &s.world.Factions[p]
	td.Draw(dst, "資金", LogicalW-330, 18, inkDim)
	td.Draw(dst, fmt.Sprintf("%d", f.Funds), LogicalW-268, 18, inkText)
	td.Draw(dst, "預備兵", LogicalW-150, 18, inkDim)
	td.Draw(dst, fmt.Sprintf("%d", totalReserves(f)*MenPerPoint), LogicalW-62, 18, inkText)
}

func totalReserves(f *state.Faction) int {
	n := 0
	for _, v := range f.Reserves {
		n += v
	}
	return n
}

func (s *Session) drawCommandBar(dst *ebiten.Image, td *textdraw.Drawer) {
	vector.DrawFilledRect(dst, 0, LogicalH-CommandH, LogicalW, CommandH, inkBar, false)
	vector.StrokeLine(dst, 0, LogicalH-CommandH, LogicalW, LogicalH-CommandH, 1, inkEdge, false)
	for i := 0; i < int(numCommands); i++ {
		x, y, w, h := CommandRect(i)
		bg := inkPanel
		// 開著的那個入口要看得出來——手機上沒有游標，
		// 不標的話玩家分不出「面板是誰開的」。
		if s.sheet.open && s.sheet.cmd == Command(i) {
			bg = inkEdge
		}
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h), bg, false)
		vector.StrokeRect(dst, float32(x), float32(y), float32(w), float32(h), 1, inkEdge, false)
		if td == nil || !td.Available() {
			continue
		}
		label := Command(i).Label()
		td.Draw(dst, label, x+(w-td.Width(label))/2, y+(h-16)/2, inkText)
	}
}

func (s *Session) drawCityCard(dst *ebiten.Image, td *textdraw.Drawer) {
	_, my, _, mh := MapRect()
	x := LogicalW - CardW - CardMargin
	y := my + mh - CardH - CardMargin
	vector.DrawFilledRect(dst, float32(x), float32(y), CardW, CardH, inkPanel, false)
	vector.StrokeRect(dst, float32(x), float32(y), CardW, CardH, 2, inkSelect, false)
	if td == nil || !td.Available() {
		return
	}
	c := &s.world.Cities[s.selected]
	td.Draw(dst, big5(c.Name), x+16, y+14, inkSelect)
	rows := [][2]string{
		{"歸屬", s.ownerName(c.Owner)},
		{"生產力", fmt.Sprintf("%d", c.Production)},
		{"防災", fmt.Sprintf("%d", c.Prevention)},
		{"城兵", fmt.Sprintf("%d", c.Garrison*MenPerPoint)},
	}
	for i, r := range rows {
		ry := y + 46 + i*28
		td.Draw(dst, r[0], x+16, ry, inkDim)
		td.Draw(dst, r[1], x+CardW-16-td.Width(r[1]), ry, inkText)
	}
}

// ownerName 用**君主姓名**當勢力名——原版的勢力就是以君主識別的
// （一覽表的「勢力」欄印的也是君主）。
func (s *Session) ownerName(owner int) string {
	if owner < 0 || owner >= len(s.world.Factions) {
		return "無主"
	}
	// 桌面版對「勢力不在了」顯示「中立」，手機版照同一個字——
	// 兩邊用不同的詞會讓人以為那是兩種狀態。
	if !s.world.Factions[owner].Alive {
		return "中立"
	}
	if n := s.world.LordName(owner); n != "" {
		return big5(n)
	}
	return "中立"
}

// big5 把存檔裡的**原始 Big5 位元組**解成可以畫的字。
//
// ⚠ `state` 的 `Name` 欄位刻意保留原始位元組（存檔要 byte-for-byte 寫回），
// 直接丟給 textdraw 會整串變成方框。桌面版同樣要先過這一層。
func big5(raw string) string { return text.Decode([]byte(raw), text.Big5) }

// MenPerPoint 是「點」換算成人數的倍率。
//
// ⚠ 勢力記錄與據點記錄存的是**點**，一點 10 人（原版 `sub_15F7F` 的
// `mul dx(=0x0A)`）。桌面版同樣只在顯示時乘——**規則層的語意不動**。
// 兩邊的換算必須一致，否則同一個城在兩個畫面上會有兩個城兵數。
const MenPerPoint = 10
