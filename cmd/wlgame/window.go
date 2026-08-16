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

	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
)

// 一覽表的幾何全部來自原版（docs/spec/38 §1.1）：視窗 (24,88,384,176)、
// 一列 16 px、一頁 10 列。**滑鼠命中與繪製共用這一組**，不會漂開。

// listRowRect 是第 visible 列的可點矩形。**從清單本體的左緣起算**——
// 左邊那 16 px 是捲軸（docs/re/26 §10）。
func listRowRect(l *listwin.List, visible int) image.Rectangle {
	if l == nil || visible < 0 || visible >= l.Height {
		return image.Rectangle{}
	}
	y := listRowY(visible)
	return image.Rect(listBodyX(), y, listBodyX()+listBodyW(), y+listRowH)
}

// listHeaderRect 是點標題排序的命中區：**一欄一格**，寬度就是分隔線
// 定義的欄寬（docs/re/26 §4.1）。
func listHeaderRect(l *listwin.List, col int) image.Rectangle {
	if l == nil || col < 0 || col >= len(l.Columns) {
		return image.Rectangle{}
	}
	f := listFieldsFor(l)
	if col >= len(f) {
		return image.Rectangle{}
	}
	x := listFieldX(f, col)
	return image.Rect(x, listWinY, x+f[col].W, listWinY+listRowH)
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
	// listActionScroll 是 ▲／▼：捲一列（原版 sub_1851A／sub_18546）。
	listActionScroll
	// listActionScrollTo 是點捲軸槽：把游標 Y 換算成新的 top
	// （原版 sub_184DD 按住不放時每輪做一次，remake 只在按下時做）。
	listActionScrollTo
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
	case listActionScroll:
		g.scrollListTo(g.list.Top + action.value)
	case listActionScrollTo:
		track := listScrollTrackRect()
		span := len(g.list.Rows) - listRowsPerPage
		if span <= 0 || track.Dy() <= 0 {
			return
		}
		g.scrollListTo((action.value - track.Min.Y) * span / track.Dy())
	}
}

// scrollListTo 把 top 夾在 0–(筆數 − 一頁) 之間，並把游標拉進可見範圍。
// **原版的 top 上限就是「筆數 − (高格 − 1)」**（sub_18546，docs/re/26 §10）。
func (g *game) scrollListTo(top int) {
	if g.list == nil {
		return
	}
	span := len(g.list.Rows) - listRowsPerPage
	if span < 0 {
		span = 0
	}
	g.list.Top = clamp(top, 0, span)
	g.list.Cursor = clamp(g.list.Cursor, g.list.Top,
		g.list.Top+listRowsPerPage-1)
	if g.list.Cursor >= len(g.list.Rows) {
		g.list.Cursor = len(g.list.Rows) - 1
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
		// 捲軸三段（原版熱區 0x3F–0x41，docs/re/26 §10）。
		p := image.Point{X: x, Y: y}
		if p.In(listScrollUpRect()) {
			g.dispatchListAction(listUIAction{kind: listActionScroll, value: -1})
			return
		}
		if p.In(listScrollDownRect()) {
			g.dispatchListAction(listUIAction{kind: listActionScroll, value: 1})
			return
		}
		if p.In(listScrollTrackRect()) {
			g.dispatchListAction(listUIAction{kind: listActionScrollTo, value: y})
			return
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
