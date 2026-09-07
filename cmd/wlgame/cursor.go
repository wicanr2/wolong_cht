package main

// 原版自己畫的滑鼠游標（docs/spec/154）。
//
// ⭐ 圖樣是**逐像素從原版擷取抽出來的**：拿有游標與沒游標的兩張同狀態
// 截圖相減，兩個不同背景抽出來的結果逐格相同。
//
// ⚠ **只在對拍時畫**（`-cursor X,Y`）。遊玩端維持系統游標——原版在
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
// **沒有指定就什麼都不畫**——這是對拍用的旗標，不是遊玩路徑。
func (g *game) drawMouseCursor(screen *ebiten.Image) {
	if g == nil || g.cursorAt == nil {
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
				float32(g.cursorAt.X+dx), float32(g.cursorAt.Y+dy),
				1, 1, col, false)
		}
	}
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
