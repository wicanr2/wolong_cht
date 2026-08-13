package main

// 四個可切換視窗。DOS/V 自然主畫面已由 strategyhud.go 畫出命令列與右側
// 常駐 HUD；這裡負責玩家切換的命令／情報／縮小地圖／系統暫存層，以及
// 開哪幾個會決定時間走不走（docs/mechanics/15-realtime.md §2）。
//
// 外框用的是原版美術（`ICONGRF` 段 3，見 internal/ui/chrome），
// 不是自己畫的線框。

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
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

// 四個視窗的位置與大小。都對齊 8 px，不然外框會切在半塊上。
var windowRect = [numWindows]struct{ X, Y, W, H int }{
	winCommand: {208, 114, 208, 194},
	winFaction: {168, 48, 240, 264},
	// 縮小地圖底圖是 192×128，加上兩側各 8 px 外框與一列標題正好 208×176。
	// 視窗小於底圖的話圖會蓋掉右邊與下邊的框——那個錯只有截圖看得出來。
	winMinimap: {archorRight(208), 48, 208, 176},
	winSystem:  {208, 114, 208, 194},
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

// lineWriter 是「一列一列往下寫」的小工具。視窗內容幾乎都是這種形狀。
type lineWriter struct {
	g    *game
	dst  *ebiten.Image
	x, y int
	col  color.RGBA
}

func (w *lineWriter) line(s string) {
	w.g.td.Draw(w.dst, s, w.x, w.y, w.col)
	w.y += textdraw.GlyphH + textdraw.LineGap
}

func (w *lineWriter) lineC(s string, c color.RGBA) {
	w.g.td.Draw(w.dst, s, w.x, w.y, c)
	w.y += textdraw.GlyphH + textdraw.LineGap
}

func (g *game) drawWindow(dst *ebiten.Image, k windowKind) {
	r := windowRect[k]
	fill := chrome.Menu
	ink := chrome.Paper
	if k == winCommand {
		fill, ink = chrome.Sheet, chrome.Ink
	}
	g.chrome.Window(dst, r.X, r.Y, r.W, r.H, fill)
	w := &lineWriter{g: g, dst: dst, x: r.X + chrome.Tile + 4, y: r.Y + chrome.Tile + 2, col: ink}

	switch k {
	case winCommand:
		// 原版的命令視窗是**八組**：進言 ＋ 說明書 §3 的其餘七個
		// （docs/mechanics/10-strategy.md）。沒實作的照樣列出來並標灰，
		// **不要因為還沒做就從選單上消失**——那會讓「還缺什麼」看不見。
		w.line("命　令")
		w.line("")
		for _, c := range commandMenu {
			col := ink
			key := c.key
			if key == "" {
				col, key = color.RGBA{150, 140, 120, 255}, "－"
			}
			w.lineC(fmt.Sprintf("%s %s", key, c.name), col)
		}
	case winFaction:
		g.drawFactionWindow(dst, r.X, r.Y, r.W, w)
	case winMinimap:
		w.line("縮小地圖")
		if img, err := g.lib.Render(g.minimapAsset(), 0, int(g.world.Clock.Season())); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(r.X+chrome.Tile), float64(r.Y+chrome.Tile+18))
			dst.DrawImage(ebiten.NewImageFromImage(img), op)
		}
	case winSystem:
		// 松崗 DOS/V 550s 是中央五列設定面板，右側以綠底值格顯示
		// 目前值；不能用一整頁 remake 快捷鍵說明取代。
		w.line("系　統　設　定")
		labels := []string{"資料儲存", "畫面模式", "音　　效", "戰略速度", "戰術速度"}
		values := []string{"S／L", "TYPE 1", "未接入", fmt.Sprintf("%d", g.speed), fmt.Sprintf("%d", g.speed)}
		rowY := r.Y + chrome.Tile + textdraw.GlyphH + 8
		valueX := r.X + r.W - chrome.Tile - 72
		for i := range labels {
			y := rowY + i*(textdraw.GlyphH+4)
			g.td.Draw(dst, labels[i], r.X+chrome.Tile+4, y, ink)
			vector.DrawFilledRect(dst, float32(valueX), float32(y-1), 64,
				float32(textdraw.GlyphH+2), color.RGBA{45, 110, 55, 255}, false)
			col := ink
			if values[i] == "未接入" {
				col = color.RGBA{170, 170, 170, 255}
			}
			g.td.Draw(dst, values[i], valueX+4, y, col)
		}
	}
}

// drawFactionWindow 是「自勢力情報」。**左上角放君主頭像**——
// 原版的君主選擇與情報畫面都有頭像，先前這裡只有數字。
//
// 頭像的頁碼是武將記錄的 `+0x01`（state.General.Portrait），
// **不是武將編號**：曹操是第 16 個武將但頭像在第 50 頁。
func (g *game) drawFactionWindow(dst *ebiten.Image, x, y, w int, lw *lineWriter) {
	p := g.world.Player
	f := g.world.Factions[p]
	season := int(g.world.Clock.Season())

	const kaoW = 64
	drawn := false
	if lord := f.Lord; lord >= 0 && lord < len(g.world.Generals) {
		if img, err := g.lib.Portrait(g.world.Generals[lord].Portrait, season); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x+chrome.Tile), float64(y+chrome.Tile))
			dst.DrawImage(ebiten.NewImageFromImage(img), op)
			drawn = true
		}
	}
	tx := x + chrome.Tile
	if drawn {
		tx += kaoW + 6
	}
	lw.x = tx
	lw.line(big5(g.world.LordName(p)) + " 軍")
	lw.lineC("軍師 "+big5(g.advisorName()), color.RGBA{200, 200, 210, 255})
	lw.line(fmt.Sprintf("信賴度 %3d", g.world.Trust))
	lw.line(fmt.Sprintf("據點   %3d", f.Cities))

	// 頭像下面接著寫，這一段用整個視窗寬度。
	lw.x = x + chrome.Tile + 4
	if drawn && lw.y < y+chrome.Tile+kaoW+2 {
		lw.y = y + chrome.Tile + kaoW + 2
	}
	lw.line(fmt.Sprintf("武將 %3d　軍團 %3d", f.Generals, f.Corps))
	lw.line(fmt.Sprintf("資金 %8d", f.Funds))
	lw.line("預備兵")
	lw.line(fmt.Sprintf("　騎馬 %6d", f.Reserves[economy.Cavalry]))
	lw.line(fmt.Sprintf("　弓兵 %6d", f.Reserves[economy.Archer]))
	lw.line(fmt.Sprintf("　步兵 %6d", f.Reserves[economy.Infantry]))
	lw.line(fmt.Sprintf("稅率 %6d%%", g.world.TaxRate))
	if g.timeRuns() {
		lw.lineC("時間 進行中", color.RGBA{140, 230, 140, 255})
	} else {
		lw.lineC("時間 停止", color.RGBA{240, 140, 140, 255})
	}
}

// minimapAsset 找出 ICONGRF 段 2（192×128 縮小地圖底圖）在素材清單裡的位置。
func (g *game) minimapAsset() int {
	for i, e := range g.lib.Entries {
		if e.Label == "ICONGRF/minimap" {
			return i
		}
	}
	return -1
}
