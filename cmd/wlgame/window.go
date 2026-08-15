package main

// 四個可切換視窗。DOS/V 自然主畫面已由 strategyhud.go 畫出命令列與右側
// 常駐 HUD；這裡負責玩家切換的命令／情報／縮小地圖／系統暫存層，以及
// 開哪幾個會決定時間走不走（docs/mechanics/15-realtime.md §2）。
//
// 外框用的是原版美術（`ICONGRF` 段 3，見 internal/ui/chrome），
// 不是自己畫的線框。

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

const (
	listWindowX       = 98
	listWindowY       = 82
	listWindowW       = 448
	listFooterButtonW = 80
)

// listWindowHeight、listRowRect 與 drawList 共用同一組幾何，讓滑鼠／觸控
// 命中區不會漂離畫面上實際畫出的列。
func listWindowHeight(l *listwin.List) int {
	if l == nil {
		return 0
	}
	return (5+(l.Height+2)*(textdraw.GlyphH+2))/chrome.Tile*chrome.Tile + 2*chrome.Tile
}

func listRowRect(l *listwin.List, visible int) image.Rectangle {
	if l == nil || visible < 0 || visible >= l.Height {
		return image.Rectangle{}
	}
	const x, y, w = listWindowX, listWindowY, listWindowW
	ry := y + chrome.Tile + 2 + textdraw.GlyphH + 6 + visible*(textdraw.GlyphH+2)
	return image.Rect(x+chrome.Tile, ry, x+w-chrome.Tile, ry+textdraw.GlyphH+2)
}

func listFooterY(l *listwin.List) int {
	return listWindowY + listWindowHeight(l) - chrome.Tile - textdraw.GlyphH
}

func listFooterRect(l *listwin.List, button int) image.Rectangle {
	if l == nil || button < 0 || button > 3 {
		return image.Rectangle{}
	}
	x := listWindowX + chrome.Tile + button*listFooterButtonW
	return image.Rect(x, listFooterY(l), x+listFooterButtonW,
		listFooterY(l)+textdraw.GlyphH+2)
}

func listHeaderRect(l *listwin.List, col int) image.Rectangle {
	if l == nil || col < 0 || col >= len(l.Columns) {
		return image.Rectangle{}
	}
	x := listWindowX + chrome.Tile + 4 + col*80
	y := listWindowY + chrome.Tile
	return image.Rect(x, y, x+80, y+textdraw.GlyphH+4)
}

func listRowAt(l *listwin.List, x, y int) (int, bool) {
	if l == nil {
		return -1, false
	}
	rows, first := l.Visible()
	for visible := range rows {
		if (image.Point{X: x, Y: y}).In(listRowRect(l, visible)) {
			return first + visible, true
		}
	}
	return -1, false
}

type listUIActionKind uint8

const (
	listActionMove listUIActionKind = iota + 1
	listActionPage
	listActionClickRow
	listActionConfirm
	listActionCancel
	listActionSort
)

type listUIAction struct {
	kind  listUIActionKind
	value int
}

// dispatchListAction 是列表的單一輸入分派器。鍵盤、滑鼠與觸控（映射成
// 滑鼠）都只產生這些動作，不直接碰背景命令或 World。
func (g *game) dispatchListAction(action listUIAction) {
	if g.list == nil {
		return
	}
	switch action.kind {
	case listActionMove:
		g.list.Move(action.value)
	case listActionPage:
		g.list.Page(action.value)
	case listActionClickRow:
		wasSelected := g.list.Phase() == listwin.Selected && g.list.Cursor == action.value
		if !g.list.Select(action.value) {
			return
		}
		if !wasSelected {
			g.list.Confirm()
			return
		}
		g.confirmListSelection()
	case listActionConfirm:
		g.confirmListSelection()
	case listActionCancel:
		if g.list.Cancel() {
			g.list = nil
		}
	case listActionSort:
		g.list.SortBy(action.value)
	}
}

func (g *game) confirmListSelection() {
	if g.list == nil {
		return
	}
	if id, ok := g.list.Confirm(); ok && g.listPick != nil && g.listPick(id) {
		g.list = nil
	}
}

func (g *game) updateListUI() {
	if g.list == nil {
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.dispatchListAction(listUIAction{kind: listActionCancel})
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if row, ok := listRowAt(g.list, x, y); ok {
			g.dispatchListAction(listUIAction{kind: listActionClickRow, value: row})
			return
		}
		for button, kind := range []listUIActionKind{
			listActionPage, listActionPage, listActionConfirm, listActionCancel,
		} {
			if (image.Point{X: x, Y: y}).In(listFooterRect(g.list, button)) {
				value := 0
				if button == 0 {
					value = -1
				} else if button == 1 {
					value = 1
				}
				g.dispatchListAction(listUIAction{kind: kind, value: value})
				return
			}
		}
		for col := range g.list.Columns {
			if (image.Point{X: x, Y: y}).In(listHeaderRect(g.list, col)) {
				g.dispatchListAction(listUIAction{kind: listActionSort, value: col})
				return
			}
		}
		return
	}
	switch {
	case pressed(ebiten.KeyArrowUp):
		g.dispatchListAction(listUIAction{kind: listActionMove, value: -1})
	case pressed(ebiten.KeyArrowDown):
		g.dispatchListAction(listUIAction{kind: listActionMove, value: 1})
	case pressed(ebiten.KeyPageUp):
		g.dispatchListAction(listUIAction{kind: listActionPage, value: -1})
	case pressed(ebiten.KeyPageDown):
		g.dispatchListAction(listUIAction{kind: listActionPage, value: 1})
	case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
		g.dispatchListAction(listUIAction{kind: listActionConfirm})
	case pressed(ebiten.KeyEscape):
		g.dispatchListAction(listUIAction{kind: listActionCancel})
	default:
		for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3,
			ebiten.Key4, ebiten.Key5} {
			if pressed(k) {
				g.dispatchListAction(listUIAction{kind: listActionSort, value: i})
				return
			}
		}
	}
}

// commandMenu 是原版命令視窗的八組指令（說明書 §3.2、§3.3，
// 整理在 docs/mechanics/10-strategy.md）。key 空字串 ＝ **還沒實作**。
//
// 保留未實作的項目是刻意的：選單本身就是一份「還缺什麼」的清單，
// 把沒做的藏起來會讓缺口從畫面上消失。
var commandMenu = []struct{ key, name string }{
	{"P", "進　言"},
	{"J", "人　事"},
	{"F", "財　政"},
	{"A", "編　成"},
	{"M", "行　軍"},
	{"C", "軍　團"},
	{"T", "據　點"},
	{"G", "武　將"},
	{"K", "勢　力"},
}

// archorRight 讓視窗靠右邊，留 8 px 邊。
func archorRight(w int) int { return screenW - w - 8 }


// minimapAsset 找出 ICONGRF 段 2（192×128 縮小地圖底圖）在素材清單裡的位置。
func (g *game) minimapAsset() int {
	for i, e := range g.lib.Entries {
		if e.Label == "ICONGRF/minimap" {
			return i
		}
	}
	return -1
}
