package phone

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 手機版的顏色與外框**全部取自原版**（docs/spec/70）：
// 底色查 `GAMEPAL.BRG`、外框是 `ICONGRF.DAT` 的 8×8 點陣圖塊，
// 兩者都由 `internal/ui/chrome` 提供——**與桌面版同一份**。
//
// ⛔ **不要在這一層抄 RGB。** `chrome` 的顏色是載入素材時查調色盤覆寫的；
// 抄一份常數就會在調色盤換算改動時悄悄脫鉤（docs/spec/54 §1 記的就是那個事故）。
// 下面每一支都是**現查**，不是快取。
//
// 對應關係（哪一塊算原版的哪一種視窗）見 docs/spec/70 §2。

// inkVoid() 是畫面最底層。原版的命令列與四周都是色 0。
func inkVoid() color.RGBA { return chrome.Blank }

// inkBar() 是狀態列與指令列的底：原版的命令列是**純黑，沒有龍紋**
//（docs/spec/54 §2）。
func inkBar() color.RGBA { return chrome.Blank }

// inkPanel() 是選單／情報視窗的底：深藍 ＋ 龍紋。
func inkPanel() color.RGBA { return chrome.Menu }

// inkSheet 是清單視窗的底：米色。一覽的四張表用它。
func inkSheet() color.RGBA { return chrome.Sheet }

// inkText() 是深藍底上的字（色 15）。
func inkText() color.RGBA { return chrome.Paper }

// inkInk 是米色底上的字（色 0）。
func inkInk() color.RGBA { return chrome.Ink }

// inkSelect() 是反白條（色 5）。
func inkSelect() color.RGBA { return chrome.Select }

// inkDim() 是次要文字。原版側欄的「君主名與『對』」用色 11
//（docs/spec/31 §2.1），這裡沿用它當次要色。
func inkDim() color.RGBA { return palSecondary }

// inkEdge() 是分隔線。沒有原版對應（原版用圖塊框，不畫線），
// 取次要色讓它與外框同一個色系。
func inkEdge() color.RGBA { return palSecondary }

// inkOverlay() 是擋住世界的決定壓在地圖上那一層：**半透明的選單底色**，
// 地圖還看得見，玩家才知道這個決定發生在哪。
func inkOverlay() color.RGBA {
	c := chrome.Menu
	c.A = 214
	return c
}

// palSecondary 是調色盤索引 11。載入素材時更新，取不到就留 fallback。
var palSecondary = color.RGBA{130, 150, 190, 255}

// secondaryIndex 是次要文字的調色盤索引（docs/spec/31 §2.1）。
const secondaryIndex = 11

// setPalette 把跟著調色盤走的顏色更新一次。**與 `chrome.Load` 同一個做法**：
// 顏色是查出來的，不是抄的。
func setPalette(lib *library.Library, bank int) {
	if lib == nil {
		return
	}
	if c, err := lib.PaletteColor(bank, secondaryIndex); err == nil {
		palSecondary = c
	}
}

// window 畫一個原版樣式的視窗：底色 ＋ `ICONGRF` 的外框圖塊。
//
// ⚠ **外框圖塊是 8×8**，寬高不是 8 的倍數時邊會切在半塊上。底色照原來的
// 矩形填滿，框則向下取到 8 的倍數——寧可框小一圈，不要切半塊（docs/spec/70 §2）。
func (s *Session) window(dst *ebiten.Image, x, y, w, h int, fill color.RGBA) {
	fillRect(dst, x, y, w, h, fill)
	if s.ch == nil {
		return
	}
	fw, fh := w/chrome.Tile*chrome.Tile, h/chrome.Tile*chrome.Tile
	if fw < chrome.Tile*2 || fh < chrome.Tile*2 {
		return
	}
	s.ch.Window(dst, x, y, fw, fh, fill)
}

// Draw 把一局畫進 960×540 的邏輯畫布。
func (s *Session) Draw(dst *ebiten.Image, td *textdraw.Drawer) {
	dst.Fill(inkVoid())
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
	// ⭐ **一覽是原版的「清單視窗」**：米色底、黑字。其餘的面板是
	// 「選單／情報視窗」：深藍底 ＋ 龍紋、白字（docs/spec/70 §2）。
	bg, fg, dim := inkPanel(), inkText(), inkDim()
	if s.sheet.cmd == CmdList {
		bg, fg, dim = inkSheet(), inkInk(), inkInk()
	}
	s.window(dst, mx, my, mw, mh, bg)
	s.drawTabs(dst, td, bg, fg)
	rows := s.sheetRows()
	if td == nil || !td.Available() {
		return
	}
	top := my + tabH
	visible := (mh - tabH) / rowH
	for i := 0; i < visible && i+s.sheet.scroll < len(rows); i++ {
		r := rows[i+s.sheet.scroll]
		y := top + i*rowH
		ink := fg
		if r.dim {
			ink = dim
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
			td.Draw(dst, c, x, y+sheetHeadDY, dim)
			x -= rowTextDX * 2
		}
	}
	if f := s.sheetFooter(); f != "" {
		td.Draw(dst, f, mx+sheetPadX, my+mh-rowH, dim)
	}
	if s.lastErr != nil {
		td.Draw(dst, s.lastErr.Error(), mx+sheetPadX, my+mh-rowH*2, fg)
	}
}

func (s *Session) drawTabs(dst *ebiten.Image, td *textdraw.Drawer, bg, fg color.RGBA) {
	tabs := s.Tabs()
	mx, my, mw, _ := MapRect()
	if len(tabs) == 0 {
		// 沒有分頁的面板仍然要有一條標題列，否則第一列會貼著狀態列，
		// 看起來像地圖的一部分。
		fillRect(dst, mx+chrome.Tile, my+chrome.Tile, mw-chrome.Tile*2, tabH-chrome.Tile, inkBar())
		if td != nil && td.Available() && s.sheet.cmd >= 0 {
			td.Draw(dst, s.sheet.cmd.Label(), mx+sheetPadX, my+sheetHeadDY+2, inkText())
		}
		return
	}
	cell := mw / len(tabs)
	for i, t := range tabs {
		x := mx + i*cell
		// 選中的分頁用**反白條**（色 5），其餘留視窗底色——
		// 原版清單的選取就是這樣標的。
		ink := fg
		if i == s.sheet.tab {
			fillRect(dst, x, my+chrome.Tile, cell, tabH-chrome.Tile, inkSelect())
			ink = inkInk()
		} else {
			fillRect(dst, x, my+chrome.Tile, cell, tabH-chrome.Tile, bg)
		}
		if td == nil || !td.Available() {
			continue
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
	// ⭐ **選對象是一張清單**（就是勢力一覽的內容），用原版的清單視窗：
	// 米色底黑字。對白那一段是原版的**對白框**，深藍 ＋ 龍紋。
	picking := s.advise.stage == advisePickAlly || s.advise.stage == advisePickTarget
	bg, fg := inkPanel(), inkText()
	if picking {
		bg, fg = inkSheet(), inkInk()
	}
	s.window(dst, mx, my, mw, mh, bg)
	if td == nil || !td.Available() {
		return
	}
	td.Draw(dst, s.adviseTitle(), mx+sheetPadX, my+sheetHeadDY+2, fg)

	if picking {
		// 對象可能有二十幾個，底部的選項條放不下——用可捲的清單。
		choices := s.AdviseChoices()
		top := my + tabH
		visible := (mh - tabH) / rowH
		for i := 0; i < visible && i+s.sheet.scroll < len(choices); i++ {
			y := top + i*rowH
			td.Draw(dst, choices[i+s.sheet.scroll], mx+sheetPadX, y+sheetHeadDY, fg)
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
			td.Draw(dst, l.text, mx+sheetPadX, y, inkSelect())
			continue
		}
		// 軍師（玩家）的話靠右，與君主分開。
		td.Draw(dst, l.text, mx+mw-sheetPadX-td.Width(l.text), y, inkText())
	}
	for i, c := range choices {
		y := bottom + i*rowH
		fillRect(dst, mx+chrome.Tile, y, mw-chrome.Tile*2, rowH-2, inkBar())
		td.Draw(dst, c, mx+sheetPadX, y+4, inkText())
	}
	if len(choices) == 0 {
		td.Draw(dst, "點畫面繼續", mx+mw-sheetPadX-td.Width("點畫面繼續"),
			my+mh-rowH, inkDim())
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
	// ⚠ 內圈用**色 15**（白）不是反白條的色 5：原版沒有「選中的據點」這個
	// 東西，色 5 是清單反白用的綠，疊在土黃城與草地上分不出來。
	// 白配黑在大地圖的每一種地形上都看得見（docs/spec/70 §2 的例外）。
	vector.StrokeRect(dst, x-1, y-1, side+2, side+2, 5, inkVoid(), false)
	vector.StrokeRect(dst, x, y, side, side, 3, inkText(), false)
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
	s.window(dst, 0, 0, LogicalW, StatusH, inkBar())
	if td == nil || !td.Available() {
		return
	}
	c := s.world.Clock
	td.Draw(dst, fmt.Sprintf("%d年%d月%d日", c.Year, c.Month, c.Day), 16, 18, inkText())

	p := s.world.Player
	if p < 0 || p >= len(s.world.Factions) {
		return
	}
	f := &s.world.Factions[p]
	td.Draw(dst, "資金", LogicalW-330, 18, inkDim())
	td.Draw(dst, fmt.Sprintf("%d", f.Funds), LogicalW-268, 18, inkText())
	td.Draw(dst, "預備兵", LogicalW-150, 18, inkDim())
	td.Draw(dst, fmt.Sprintf("%d", totalReserves(f)*MenPerPoint), LogicalW-62, 18, inkText())
}

func totalReserves(f *state.Faction) int {
	n := 0
	for _, v := range f.Reserves {
		n += v
	}
	return n
}

func (s *Session) drawCommandBar(dst *ebiten.Image, td *textdraw.Drawer) {
	fillRect(dst, 0, LogicalH-CommandH, LogicalW, CommandH, inkBar())
	for i := 0; i < int(numCommands); i++ {
		x, y, w, h := CommandRect(i)
		// 開著的那個入口用**反白條的顏色**——手機上沒有游標，
		// 不標的話玩家分不出「面板是誰開的」。原版的反白就是色 5。
		bg := inkBar()
		ink := inkText()
		if s.sheet.open && s.sheet.cmd == Command(i) {
			bg, ink = inkSelect(), inkInk()
		}
		s.window(dst, x, y, w, h, bg)
		if td == nil || !td.Available() {
			continue
		}
		label := Command(i).Label()
		td.Draw(dst, label, x+(w-td.Width(label))/2, y+(h-16)/2, ink)
	}
}

func (s *Session) drawCityCard(dst *ebiten.Image, td *textdraw.Drawer) {
	_, my, _, mh := MapRect()
	x := LogicalW - CardW - CardMargin
	y := my + mh - CardH - CardMargin
	s.window(dst, x, y, CardW, CardH, inkPanel())
	if td == nil || !td.Available() {
		return
	}
	c := &s.world.Cities[s.selected]
	td.Draw(dst, big5(c.Name), x+16, y+14, inkText())
	rows := [][2]string{
		{"歸屬", s.ownerName(c.Owner)},
		{"生產力", fmt.Sprintf("%d", c.Production)},
		{"防災", fmt.Sprintf("%d", c.Prevention)},
		{"城兵", fmt.Sprintf("%d", c.Garrison*MenPerPoint)},
	}
	for i, r := range rows {
		ry := y + 46 + i*28
		td.Draw(dst, r[0], x+16, ry, inkDim())
		td.Draw(dst, r[1], x+CardW-16-td.Width(r[1]), ry, inkText())
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
