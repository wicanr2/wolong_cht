package main

// 「自定」軍師命名視窗（docs/spec/104；原版 `sub_18FC9`，顯示清單場景 9）。
//
// 版面全部出自原版（docs/re/54 §1–§2）；選字表是 `END_S15.DAT`
// （docs/formats/10）。模型與畫面分開：`namingModel` 不認識 Ebiten，
// 無頭測試驗的是行為。

import (
	"image"
	"image/color"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/assets/namechars"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

// 版面常數（docs/re/54 §1；`sub_1928A`／`sub_19223`／`sub_190C0` 的立即值）。
const (
	namingWinX, namingWinY, namingWinW, namingWinH = 192, 128, 352, 256

	namingInputX, namingInputY, namingInputW, namingInputH = 272, 144, 256, 64
	namingLineX, namingLineY, namingLineLen                = 273, 176, 54
	namingLabelNameX, namingLabelAliasX, namingLabelY      = 336, 144, 144
	// 六格名字：`sub_19223` 前三格從 x=352、後三格從 x=448，y=168；
	// 目前格的底線畫在 y=186。
	namingCellsX, namingCellsAliasX, namingCellsY, namingCellW = 352, 448, 168, 16
	namingCursorY                                              = 186

	namingPortraitX, namingPortraitY = 208, 144 // 假說：輸入區左邊的空位
	namingPrevX, namingPrevY         = 280, 152 // 「前 ▲」（熱區 0x27：(272,152) 56×32）
	namingNextX, namingNextY         = 280, 184 // 「後 ▼」（熱區 0x28：(272,176) 56×32）

	namingBtnY                         = 192
	namingRedoX, namingContX, namingOKX = 352, 408, 464
	namingBtnW, namingOKW, namingBtnH  = 48, 64, 16

	namingInitialsY   = 216 // 聲母列（`cs:1871` 的 42 bytes，屬性 0F01）
	namingInitialsX   = 200
	namingInitialHotX = 216 // 熱區 0x25：(216,216) 320×16，每個聲母 32 px

	// 選字格：`sub_1928A` 起點 (210,238)、格距 20、16 欄 × 6 列。
	namingGridX, namingGridY, namingGridPitch = 210, 238, 20
	namingGridCols, namingGridRows            = 16, 6
	namingGridBoxX, namingGridBoxY            = 200, 232 // 底 336×128
	namingGridBoxW, namingGridBoxH            = 336, 128

	namingPagerY, namingPagerH = 360, 16
	namingPagerLineX           = 367
	namingPrevPageX            = 232
	namingNextPageX            = 400

	namingPageChars = namingGridCols * namingGridRows // 96
	namingMaxPage   = 0x13BA / 2                     // `sub_1908D` 的上限（字索引）
	namingMaxIndex  = 0x1478 / 2                     // `sub_192E3` 的上限
	namingCells     = 6
	namingPortraits = 0x93 // 肖像 0..0x92（`sub_1912D`／`sub_19144`）
)

// namingInitials 是聲母列的十個跳點（字本身就在表裡當分段標記，docs/formats/10 §3）。
var namingInitials = []rune("ㄅㄈㄋㄍㄐㄑㄓㄕㄙㄨ")

// namingHotspot 是原版的九個熱區（docs/re/54 §2）。
type namingHotspot int

const (
	namingNone namingHotspot = iota
	namingOK                 // 0x20 確定
	namingRedo               // 0x21 重來：清掉目前格、退一格（`sub_1905E`）
	namingCont               // 0x22 繼續：清掉目前格、跳下一格（`sub_1906E`）
	namingPrevPage           // 0x23
	namingNextPage           // 0x24
	namingInitial            // 0x25 聲母列
	namingPick               // 0x26 選字
	namingPortraitPrev       // 0x27 前 ▲
	namingPortraitNext       // 0x28 後 ▼
)

// namingModel 是視窗的狀態。
type namingModel struct {
	table *namechars.Table
	// cells 是六格，每格一個 Big5 碼（0 ＝ 還沒填，畫成全形空白）。
	cells    [namingCells]uint16
	cursor   int // 0..5
	page     int // 本頁第一個字的索引
	portrait int
	used     map[int]bool // 武將表裡有人用的肖像，◀▶ 跳過（`sub_192F5`）
	done     bool
	cancel   bool
}

func newNamingModel(table *namechars.Table, w *state.World) *namingModel {
	m := &namingModel{table: table, portrait: 0x91, used: map[int]bool{}}
	if w != nil {
		for _, g := range w.Generals {
			m.used[g.Portrait] = true
		}
		if w.AdvisorPortrait != 0xFF {
			m.portrait = w.AdvisorPortrait
		}
	}
	return m
}

// hotspotAt 是 `sub_18FC9` 的 `sub al, 20h` 分派。
func namingHotspotAt(x, y int) namingHotspot {
	p := image.Pt(x, y)
	in := func(x0, y0, w, h int) bool { return p.In(image.Rect(x0, y0, x0+w, y0+h)) }
	switch {
	case in(namingOKX, namingBtnY, namingOKW, namingBtnH):
		return namingOK
	case in(namingRedoX, namingBtnY, namingBtnW, namingBtnH):
		return namingRedo
	case in(namingContX, namingBtnY, namingBtnW, namingBtnH):
		return namingCont
	case in(200, namingPagerY, 168, namingPagerH):
		return namingPrevPage
	case in(368, namingPagerY, 168, namingPagerH):
		return namingNextPage
	case in(namingInitialHotX, namingInitialsY, 320, 16):
		return namingInitial
	case in(namingGridBoxX, namingGridBoxY, namingGridBoxW, namingGridBoxH):
		return namingPick
	case in(272, 152, 56, 32):
		return namingPortraitPrev
	case in(272, 176, 56, 32):
		return namingPortraitNext
	}
	return namingNone
}

// gridIndexAt 是 `sub_190C0`：(x−210)÷20、(y−238)÷20 的商是欄／列，
// 餘數 ≥ 16 是點在格間的空隙。
func namingGridIndexAt(x, y int) (int, bool) {
	dx, dy := x-namingGridX, y-namingGridY
	if dx < 0 || dy < 0 {
		return 0, false
	}
	col, cr := dx/namingGridPitch, dx%namingGridPitch
	row, rr := dy/namingGridPitch, dy%namingGridPitch
	if col >= namingGridCols || row >= namingGridRows || cr >= 16 || rr >= 16 {
		return 0, false
	}
	return row*namingGridCols + col, true
}

// click 處理一次左鍵。
func (m *namingModel) click(x, y int) {
	switch namingHotspotAt(x, y) {
	case namingOK:
		m.done = true
	case namingRedo:
		m.cells[m.cursor] = 0
		if m.cursor > 0 {
			m.cursor--
		}
	case namingCont:
		m.cells[m.cursor] = 0
		if m.cursor < namingCells-1 {
			m.cursor++
		}
	case namingPrevPage:
		m.setPage(m.page - namingPageChars)
	case namingNextPage:
		m.setPage(m.page + namingPageChars)
	case namingInitial:
		k := (x - namingInitialHotX) / 32
		if k >= 0 && k < len(namingInitials) {
			m.jumpTo(namingInitials[k])
		}
	case namingPick:
		if i, ok := namingGridIndexAt(x, y); ok {
			m.pick(m.page + i)
		}
	case namingPortraitPrev:
		m.stepPortrait(-1)
	case namingPortraitNext:
		m.stepPortrait(+1)
	}
}

func (m *namingModel) setPage(p int) {
	if p < 0 {
		p = 0
	}
	if p > namingMaxPage {
		p = namingMaxPage
	}
	m.page = p
}

// jumpTo 把頁翻到某個聲母的分段標記（那個字本身在表裡）。
func (m *namingModel) jumpTo(r rune) {
	for i, c := range m.table.Runes {
		if c == r {
			m.setPage(i)
			return
		}
	}
}

// pick 是 `sub_190C0` 的後半：把字寫進目前格，游標前進（最後一格不動）。
func (m *namingModel) pick(index int) {
	if index < 0 || index > namingMaxIndex || index >= len(m.table.Big5) {
		return
	}
	m.cells[m.cursor] = m.table.Big5[index]
	if m.cursor < namingCells-1 {
		m.cursor++
	}
}

// stepPortrait 是 `sub_1912D`／`sub_19144`：在 0..0x92 循環，跳過武將在用的。
func (m *namingModel) stepPortrait(d int) {
	for i := 0; i < namingPortraits; i++ {
		m.portrait = (m.portrait + d + namingPortraits) % namingPortraits
		if !m.used[m.portrait] {
			return
		}
	}
}

// nameBytes 是要寫進區塊 +0x52A2 的 12 個 Big5 位元組（空格是全形空白）。
func (m *namingModel) nameBytes() []byte {
	out := make([]byte, 0, namingCells*2)
	for _, c := range m.cells {
		if c == 0 {
			c = 0xA140
		}
		out = append(out, byte(c>>8), byte(c))
	}
	return out
}

// hasName 回報第一格填了沒——原版「確定」要 +0x02 ≠ 0x7F 或名字非空才放行
// （`sub_18E5A` 的 `cmp word ptr ds:5222h, 0D0A1h`）。
func (m *namingModel) hasName() bool { return m.cells[0] != 0 }

// cellText 把一格轉成畫面用的 Big5 字串。
func (m *namingModel) cellText(i int) string {
	c := m.cells[i]
	if c == 0 {
		return "\xA1\x40"
	}
	return string([]byte{byte(c >> 8), byte(c)})
}

// ── 畫面 ─────────────────────────────────────────────────────────

func (g *game) openNaming() error {
	table, err := namechars.Load(filepath.Join(g.origDir, "END_S15.DAT"))
	if err != nil {
		return err
	}
	g.naming = newNamingModel(table, g.launcherPreviewWorld)
	return nil
}

func (g *game) drawNaming(screen *ebiten.Image, season int) {
	m := g.naming
	g.chrome.Window(screen, namingWinX, namingWinY, namingWinW, namingWinH, chrome.Menu)
	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	line9 := g.paletteInk(0x09, chrome.Paper)

	// 靜態層（場景 9）。
	g.dlFill(screen, namingInputX, namingInputY, namingInputW, namingInputH, 0x00, color.RGBA{0, 0, 0, 255})
	g.dlFill(screen, namingLineX, namingLineY, namingLineLen, 1, strategyInkNormal, chrome.Paper)
	g.td.Draw(screen, "軍師名", namingLabelNameX, namingLabelY, ink)
	g.td.Draw(screen, "別　號", namingLabelAliasX, namingLabelY, ink)
	g.td.Draw(screen, "前 ▲", namingPrevX, namingPrevY, ink)
	g.td.Draw(screen, "後 ▼", namingNextX, namingNextY, ink)
	for _, b := range []struct {
		x, w  int
		label string
	}{{namingRedoX, namingBtnW, "重來"}, {namingContX, namingBtnW, "繼續"}, {namingOKX, namingOKW, "確定"}} {
		g.dlButton(screen, b.x, namingBtnY, b.w, namingBtnH)
		g.td.Draw(screen, b.label, b.x+8, namingBtnY, g.dlButtonInk())
	}
	g.dlFill(screen, namingGridBoxX, namingGridBoxY, namingGridBoxW, namingGridBoxH, 0x00, color.RGBA{0, 0, 0, 255})
	g.dlFill(screen, namingGridBoxX, namingPagerY, namingGridBoxW, namingPagerH, 0x00, color.RGBA{0, 0, 0, 255})
	g.dlFill(screen, namingPagerLineX, namingPagerY+1, 1, 14, 0x09, line9)
	g.td.Draw(screen, "上一頁　▲　", namingPrevPageX, namingPagerY, ink)
	g.td.Draw(screen, "下一頁　▼　", namingNextPageX, namingPagerY, ink)
	initials := "　"
	for _, r := range namingInitials {
		initials += string(r) + "　"
	}
	g.td.Draw(screen, initials, namingInitialsX, namingInitialsY, ink)

	// 動態層：肖像、六格、游標底線、本頁的字。
	if img, err := g.lib.Portrait(m.portrait, season); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(namingPortraitX, namingPortraitY)
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	for i := 0; i < namingCells; i++ {
		x := namingCellsX + i*namingCellW
		if i >= 3 {
			x = namingCellsAliasX + (i-3)*namingCellW
		}
		g.td.Draw(screen, big5(m.cellText(i)), x, namingCellsY, g.paletteInk(0x0F, chrome.Paper))
		col := 0x01
		if i == m.cursor {
			col = 0x0F
		}
		g.dlFill(screen, x, namingCursorY, namingCellW-2, 1, col, chrome.Paper)
	}
	for row := 0; row < namingGridRows; row++ {
		for c := 0; c < namingGridCols; c++ {
			idx := m.page + row*namingGridCols + c
			if idx >= len(m.table.Runes) {
				break
			}
			g.td.Draw(screen, string(m.table.Runes[idx]),
				namingGridX+c*namingGridPitch, namingGridY+row*namingGridPitch,
				g.paletteInk(0x09, chrome.Paper))
		}
	}
}

// customAdvisor 是命名視窗「確定」之後留給開局用的結果。
type customAdvisor struct {
	portrait int
	name     []byte
}

// settleNaming 收尾：確定 → 留下結果、把卡片上的軍師換成自訂的；取消 → 什麼都不留
// （原版取消時 `sub_18F7C` 把 +0x02 從 +0x3F 抄回、緩衝區清成空標記）。
func (g *game) settleNaming() {
	m := g.naming
	switch {
	case m.cancel:
		g.naming = nil
	case m.done:
		if !m.hasName() {
			// 原版「確定」要第一格有字才放行（`sub_18E5A`）。
			m.done = false
			return
		}
		g.customAdvisor = &customAdvisor{portrait: m.portrait, name: m.nameBytes()}
		if l := g.launcher; l != nil && l.cursor >= 0 && l.cursor < len(l.players) {
			p := &l.players[l.cursor]
			p.HasAdvisor, p.AdvisorPortrait = true, m.portrait
			p.Advisor = big5(string(m.nameBytes()[:advisorNameHalfBytes]))
		}
		g.naming = nil
	}
}

// advisorNameHalfBytes 是「軍師名」那三格的位元組數。
const advisorNameHalfBytes = 6
