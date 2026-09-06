package main

// 左下角的狀態列提示框（docs/spec/140）。
//
// 原版 `sub_18853` 在每個指令流程進行中掛一個 `(0, 320, 256, 80)` 的框，
// **與一般訊息框是同一個框，只是換了位置**（`docs/re/66` §1.2）。
// 肖像固定 `0x93`（通報者），`cx = 0FFFFh` 是清掉。
//
// ⚠ 它**不是模態的**：框掛著的同時流程繼續（選據點、跳選單）。

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

const (
	// 框 ＝ `sub_189A4(dx=0, bx=12h, cx=510h)`：
	//   X ＝ dx × 16 ＝ 0、Y ＝ (bx + 2) × 16 ＝ 320
	//   寬 ＝ cl × 16 ＝ 256、高 ＝ ch × 16 ＝ 80  ← 與訊息框同一組立即值
	statusBoxX = 0
	statusBoxY = 320

	// talkInkCity 是 `\2`（據點名）代入之後畫的顏色。
	//
	// ⭐ **每個標記各有自己的顏色**，七支 handler 的表在
	// `docs/formats/01` §3／`docs/re/79` §2：`\1`／`\4` 色 9、
	// `\2` 色 `0x0B`、`\3`／`\5` 色 `0x0C`。
	// 原版擷取上量到的正是 `(243,162,0)` ＝ 第 11 色。
	talkInkCity = 0x0B

	// talkFieldCells 是 `\1`–`\5` 代入時原版畫幾個全形字。
	//
	// ⭐ 五支 handler 一律 `al = 3` 再 `call loc_10701`（docs/re/79 §2），
	// 也就是**照欄位寬度畫滿三格，補白也畫出來**——李暹的呼び名整格是
	// 三個全形空白，原版真的畫一片空白。「照抄就是照抄。」
	talkFieldCells = 3
)

// padTalkField 把代入的名字補到 `talkFieldCells` 個全形字。
//
// ⚠ **要在語系轉換之後補**：`uitext.Convert` 的覆寫詞表是整串比對的，
// 帶著補白會查不到，名字會原封不動留在畫面上。
func padTalkField(s string) string {
	want := talkFieldCells * talkLinePitch
	for textdraw.StringWidth(s) < want {
		s += "　"
	}
	return s
}

// statusBoxState 是現在掛著哪一則。**lines 是空的就是沒掛**——零值安全。
type statusBoxState struct {
	lines []string
	// city 是 `\2` 代入的據點名——句中那一段要換色（見 talkInkCity）。
	city string
}

// setStatusTalk 掛一則 TALK 訊息（原版 `sub_18853(cx = 索引)`）。
//
// 取不到那一則就當作沒掛：fail-closed 比畫半句好。
func (g *game) setStatusTalk(index int, vars map[byte]string) {
	// `\1`–`\5` 都是定寬欄位（見 talkFieldCells）。`\6` 不畫字、
	// `\7` 走數值繪製，兩個都不補。
	padded := make(map[byte]string, len(vars))
	for k, v := range vars {
		if k >= '1' && k <= '5' && v != "" {
			v = padTalkField(v)
		}
		padded[k] = v
	}
	vars = padded
	city := vars['2']
	lines, ok := g.talkLines(index, vars)
	if !ok || len(lines) == 0 {
		g.statusBox = statusBoxState{}
		return
	}
	// ⚠ **TALK 的斷行是照原版版面排的**，那個版面是全形字。半形語系
	// 照抄會畫出框外（docs/spec/87 §8），所以重折一次；繁中原本就放得下。
	g.statusBox = statusBoxState{lines: layoutMessageLines(lines), city: city}
}

// clearStatusTalk 清掉（原版 `sub_18853(cx = 0FFFFh)`）。
func (g *game) clearStatusTalk() { g.statusBox = statusBoxState{} }

// statusBoxActive 回報有沒有掛著。
func (g *game) statusBoxActive() bool { return len(g.statusBox.lines) > 0 }

func (g *game) drawStatusBox(screen *ebiten.Image) {
	if !g.statusBoxActive() {
		return
	}
	lines := g.statusBox.lines
	if len(lines) > messagePageRows {
		lines = lines[:messagePageRows]
	}
	// 框與肖像走一般訊息框那一支——原版也是同一支（`sub_1075B`）。
	// 文字自己畫，因為 `\2` 代入的據點名要換色。
	g.drawLegacyTalkBox(screen, statusBoxX, statusBoxY, talkBoxW, talkBoxH,
		nil, defaultPortraitPage)
	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	cityInk := g.paletteInk(talkInkCity, color.RGBA{243, 162, 0, 255})
	for i, line := range lines {
		g.drawTalkLineWithName(screen, line, g.statusBox.city,
			statusBoxX+talkTextX-talkBoxX,
			statusBoxY+talkTextY-talkBoxY+i*talkLinePitch, ink, cityInk)
	}
}
