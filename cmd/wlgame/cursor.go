package main

// 原版自己畫的滑鼠游標（docs/spec/154）。
//
// ⭐ 圖樣是**逐像素從原版擷取抽出來的**：拿有游標與沒游標的兩張同狀態
// 截圖相減，兩個不同背景抽出來的結果逐格相同。
//
// 桌面捕捉模式也用此圖樣顯示虛擬游標；`-cursor X,Y` 可固定對拍位置。原版在
// 「清單等待點選」時旗標是 0（畫面上沒有游標），而那到底是原版行為
// 還是 oracle 的限制還沒定案（docs/spec/154 §4）。

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

// arrowCursor 是 14×14 的箭頭：`W` ＝ 色 15 描邊、`#` ＝ 色 10 填色。
// **熱點在左上角 (0, 0)**。
var arrowCursor = [14]string{
	"WW............",
	"W#WW..........",
	".W##WW........",
	".W####WW......",
	"..W#####WW....",
	"..W#######WW..",
	"...W########WW",
	"...W#######W..",
	"....W#####W...",
	"....W######W..",
	".....W##W###W.",
	".....W#W.W###W",
	"......W...W##W",
	"......W....WW.",
}

const (
	arrowCursorEdge = 0x0F // 描邊：色 15
	arrowCursorFill = 0x0A // 填色：色 10
)

// drawMouseCursor 在 `cursorAt` 指定的位置畫那個箭頭。
// 未指定固定位置時，捕捉模式使用正常玩家游標位置。
func (g *game) drawMouseCursor(screen *ebiten.Image) {
	if g == nil {
		return
	}
	at := g.cursorAt
	if at == nil && desktopPointer.active {
		if g.desktopMapCursorVisible() {
			return
		}
		x, y := cursorPosition()
		at = &image.Point{X: x, Y: y}
	}
	if at == nil {
		return
	}
	edge := g.paletteInk(arrowCursorEdge, chrome.Paper)
	fill := g.paletteInk(arrowCursorFill, listWarnInk)
	for dy, row := range arrowCursor {
		for dx, c := range row {
			var col = edge
			switch c {
			case 'W':
			case '#':
				col = fill
			default:
				continue
			}
			vector.DrawFilledRect(screen,
				float32(at.X+dx), float32(at.Y+dy),
				1, 1, col, false)
		}
	}
}

// desktopMapCursorVisible 與可點選地圖共用熱區界線；其他 UI 維持箭頭。
func (g *game) desktopMapCursorVisible() bool {
	if g == nil || !desktopPointer.active || g.cursorAt != nil || g.world == nil ||
		g.launcher != nil || g.battleActive() || g.hudOpen(hudSystem) ||
		g.quitting || g.quitMenu.active || g.endingActive() || g.messageActive() {
		return false
	}
	if !g.mapPickActive() && g.mapClickBusy() {
		return false
	}
	x, y := cursorPosition()
	if g.hudOpen(hudCommand) && x < strategyCommandW && y < strategyCommandY+strategyCommandH ||
		g.hudOpen(hudMinimap) && x >= strategySidebarX && y < strategyFactionY ||
		g.hudOpen(hudFaction) && x >= strategySidebarX && y >= strategyFactionY {
		return false
	}
	_, _, ok := g.mapPickTileAt(x, y)
	return ok
}

// cursorPoint 解析 `-cursor X,Y`。空字串回 nil（不畫）。
func cursorPoint(spec string) *image.Point {
	var x, y int
	if spec == "" {
		return nil
	}
	if n, err := fmt.Sscanf(spec, "%d,%d", &x, &y); err != nil || n != 2 {
		return nil
	}
	return &image.Point{X: x, Y: y}
}
