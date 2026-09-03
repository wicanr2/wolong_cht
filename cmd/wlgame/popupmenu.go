package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/ui/talkmenu"
)

// 指令列的彈出選單（docs/spec/126）。
//
// 進言、人事、軍團、據點四格點下去都是先跳一張 `sub_193E9` 的選單，
// 差別只有 TALK 索引、位置與選中之後往哪走——所以**只留一份實作**
// （CLAUDE.md §7 第 6 條）。進言那一張還沒併進來（spec/126 §4）。

// popupMenu 是一張彈出選單的定義。
type popupMenu struct {
	// talk 是選項文字的 TALK 索引（`sub_193E9` 的 `cx`）。
	talk int
	// x, y 是框的左上角，由 `dx` 的粗格換算：高 byte 是列、低 byte 是欄。
	x, y int
	// cell 是這張選單屬於指令列哪一格——反白要用（docs/spec/124）。
	cell naturalCommandID
	// fallback 只在讀不到 `TALK.DAT` 時用。
	//
	// ⭐ 兩邊的**全形空白要留著**：框寬由字數決定，原版把每一列補到
	// 等寬（docs/spec/125）。
	fallback []string
	// dispatch 是選中第 row 項之後往哪走。
	dispatch func(*game, int)
}

// popupMenuCol 把指令列的格號換成選單的粗格欄號。
//
// 指令格寬 48 px ＝ 3 個 16 px 的粗格，而四張選單的 `dx` 低 byte
// 分別是 0／3／12／15 ＝ 索引 × 3（docs/spec/126 §1）。
func popupMenuCol(cell naturalCommandID) int { return int(cell) * 3 }

// 四張裡的三張。`y` 固定是粗格第 4 列（`dx` 的高 byte 都是 4）。
const popupMenuRow = 4

var (
	// personnelPopupMenu 是「人事」（`sub_16265`：ax=4, cx=4Eh, dx=403h）。
	personnelPopupMenu = &popupMenu{
		talk: 0x4E, cell: naturalCommandPersonnel,
		fallback: []string{"　內政官任命　", "　內政官解任　",
			"　外交官任命　", "　外交官解任　"},
		dispatch: (*game).dispatchPersonnelMenu,
	}
	// corpsPopupMenu 是「軍團」（`sub_1628F`：ax=2, cx=4Fh, dx=40Ch）。
	corpsPopupMenu = &popupMenu{
		talk: 0x4F, cell: naturalCommandCorps,
		fallback: []string{"　位置確認　", "　行軍指示　"},
		dispatch: (*game).dispatchCorpsMenu,
	}
	// cityPopupMenu 是「據點」（`sub_162FB`：ax=2, cx=52h, dx=40Fh）。
	cityPopupMenu = &popupMenu{
		talk: 0x52, cell: naturalCommandCity,
		fallback: []string{"　首都確認　", "　據點一覽　"},
		dispatch: (*game).dispatchCityMenu,
	}
)

// popupMenusByName 給驗收旗標 `-open-command-menu` 用。
var popupMenusByName = map[string]*popupMenu{
	"corps": corpsPopupMenu, "city": cityPopupMenu,
	"personnel": personnelPopupMenu,
}

func init() {
	// 位置由格號算出來，不要三處各抄一組座標。
	for _, m := range []*popupMenu{personnelPopupMenu, corpsPopupMenu,
		cityPopupMenu} {
		m.x, m.y = popupMenuCol(m.cell)*16, popupMenuRow*16
	}
}

// popupMenuState 是「現在開著哪一張」。**menu 是 nil 就是沒開**——
// 零值安全，不必另外記一個 active 旗標。
type popupMenuState struct {
	menu *popupMenu
	row  int
}

// popupMenuActive 回報有沒有選單開著。
func (g *game) popupMenuActive() bool { return g != nil && g.cmdMenu.menu != nil }

// openPopupMenu 開一張選單。
func (g *game) openPopupMenu(m *popupMenu) { g.cmdMenu.menu, g.cmdMenu.row = m, 0 }

// closePopupMenu 收掉目前那一張。
func (g *game) closePopupMenu() { g.cmdMenu.menu = nil }

// popupMenuLabels 是那張選單的每一列，直接取自 `TALK.DAT`。
func (g *game) popupMenuLabels(m *popupMenu) []string {
	if g == nil || g.lib == nil || m == nil {
		return append([]string(nil), m.fallback...)
	}
	return talkmenu.MenuLabels(g.lib.Talk, m.talk, nil, m.fallback)
}

// dispatchPopupMenu 收掉選單再走那一項。
//
// ⚠ **先收再走**：下一層可能自己開別的視窗，收在後面會把它一起關掉。
func (g *game) dispatchPopupMenu(row int) {
	m := g.cmdMenu.menu
	if m == nil {
		return
	}
	g.closePopupMenu()
	m.dispatch(g, row)
}

// updatePopupMenu 處理選單的輸入。回傳 true 表示它吃掉了這一幀。
func (g *game) updatePopupMenu() bool {
	if !g.popupMenuActive() {
		return false
	}
	m := g.cmdMenu.menu
	labels := g.popupMenuLabels(m)
	if row, ok := g.talkChoiceClick(m.x, m.y, labels); ok {
		g.dispatchPopupMenu(row)
		return true
	}
	switch {
	case pressed(ebiten.KeyArrowUp):
		g.cmdMenu.row = (g.cmdMenu.row + len(labels) - 1) % len(labels)
	case pressed(ebiten.KeyArrowDown):
		g.cmdMenu.row = (g.cmdMenu.row + 1) % len(labels)
	case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
		g.dispatchPopupMenu(g.cmdMenu.row)
	case g.cancelled():
		g.closePopupMenu()
	}
	// 數字鍵是 remake 加的捷徑；原版只有游標選取。
	for i := range labels {
		if pressed(ebiten.Key1 + ebiten.Key(i)) {
			g.dispatchPopupMenu(i)
			break
		}
	}
	return true
}

// drawPopupMenu 畫目前開著的那一張。
func (g *game) drawPopupMenu(screen *ebiten.Image) {
	if !g.popupMenuActive() {
		return
	}
	m := g.cmdMenu.menu
	g.drawLegacyChoiceBox(screen, m.x, m.y, g.popupMenuLabels(m), g.cmdMenu.row)
}
