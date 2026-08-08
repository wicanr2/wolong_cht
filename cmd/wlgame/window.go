package main

// 四個視窗。原版的版面**沒有常駐側欄**——情報、命令、縮小地圖、系統
// 全部是浮在大地圖上的視窗，而且開哪幾個會決定時間走不走
// （docs/mechanics/15-realtime.md §2）。
//
// 外框用的是原版美術（`ICONGRF` 段 3，見 internal/ui/chrome），
// 不是自己畫的線框。

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 四個視窗的位置與大小。都對齊 8 px，不然外框會切在半塊上。
var windowRect = [numWindows]struct{ X, Y, W, H int }{
	winCommand: {8, 48, 152, 216},
	winFaction: {168, 48, 240, 264},
	// 縮小地圖底圖是 192×128，加上兩側各 8 px 外框與一列標題正好 208×176。
	// 視窗小於底圖的話圖會蓋掉右邊與下邊的框——那個錯只有截圖看得出來。
	winMinimap: {archorRight(208), 48, 208, 176},
	winSystem:  {archorRight(216), 216, 216, 152},
}

// commandMenu 是原版命令視窗的八組指令（說明書 §3.2、§3.3，
// 整理在 docs/mechanics/10-strategy.md）。key 空字串 ＝ **還沒實作**。
//
// 保留未實作的項目是刻意的：選單本身就是一份「還缺什麼」的清單，
// 把沒做的藏起來會讓缺口從畫面上消失。
var commandMenu = []struct{ key, name string }{
	{"P", "進　言"},
	{"", "人　事"},
	{"", "財　政"},
	{"A", "編　成"},
	{"M", "行　軍"},
	{"C", "軍　團"},
	{"", "據　點"},
	{"G", "武　將"},
	{"", "勢　力"},
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
		w.line("系　統")
		w.line("")
		w.lineC("此視窗開著時間停止", color.RGBA{255, 200, 120, 255})
		w.line(fmt.Sprintf("速度　%d　（−／＝ 調整）", g.speed))
		w.line("方向鍵　捲動地圖")
		w.line("1-4　開關視窗")
		w.line("ESC　關閉　F10　離開")
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
