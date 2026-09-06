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
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

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
	// talkInkPerson 是 `\1`／`\4`（武將呼び名）、talkInkLord 是
	// `\3`／`\5`（君主姓名）的顏色（docs/spec/119 §1）。
	talkInkPerson = 9
	talkInkLord   = 0x0C

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

// talkField 是一個代入的欄位：畫出來的字，以及它自己的顏色
// （docs/spec/119 §3.1）。**顏色是標記的性質，不是框的性質**。
type talkField struct {
	text string
	ink  int // 原版調色盤索引
}

// talkMarkerInk 是每個標記代入之後畫的顏色（docs/spec/119 §1）。
// 不是 `\1`–`\5` 就回 −1（`\6` 是排版控制、`\7` 走數值繪製，都不換色）。
func talkMarkerInk(marker byte) int {
	switch marker {
	case '1', '4':
		return talkInkPerson
	case '2':
		return talkInkCity
	case '3', '5':
		return talkInkLord
	}
	return -1
}

// padTalkVars 把 `\1`–`\5` 補到定寬，並收集成 []talkField
// （docs/spec/119 §3.1）。⚠ **呼叫端不要再自己補**。
//
// 欄位**由長到短**排序：短的先比會把長的切一半（「曹操」對上「曹操軍」）。
func padTalkVars(vars map[byte]string) (map[byte]string, []talkField) {
	if len(vars) == 0 {
		return vars, nil
	}
	out := make(map[byte]string, len(vars))
	var fields []talkField
	for k, v := range vars {
		ink := talkMarkerInk(k)
		if ink >= 0 && v != "" {
			v = padTalkField(v)
			fields = append(fields, talkField{text: v, ink: ink})
		}
		out[k] = v
	}
	sort.SliceStable(fields, func(i, j int) bool {
		return len(fields[i].text) > len(fields[j].text)
	})
	return out, fields
}

// statusBoxState 是現在掛著哪一則。**lines 是空的就是沒掛**——零值安全。
type statusBoxState struct {
	lines []string
	// fields 是代入的欄位，每個各有自己的顏色（docs/spec/119 §3.1）。
	fields []talkField
}

// setStatusTalk 掛一則 TALK 訊息（原版 `sub_18853(cx = 索引)`）。
//
// 取不到那一則就當作沒掛：fail-closed 比畫半句好。
func (g *game) setStatusTalk(index int, vars map[byte]string) {
	// `\1`–`\5` 都是定寬欄位，而且各有各的顏色（docs/spec/119 §3.1）。
	// `\6` 不畫字、`\7` 走數值繪製，兩個都不補也不換色。
	vars, fields := padTalkVars(vars)
	lines, ok := g.talkLines(index, vars)
	if !ok || len(lines) == 0 {
		g.statusBox = statusBoxState{}
		return
	}
	// ⚠ **TALK 的斷行是照原版版面排的**，那個版面是全形字。半形語系
	// 照抄會畫出框外（docs/spec/87 §8），所以重折一次；繁中原本就放得下。
	g.statusBox = statusBoxState{lines: layoutMessageLines(lines), fields: fields}
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
	// 框、肖像與文字全部走一般訊息框那一支——原版也是同一支（`sub_1075B`），
	// 連代入欄位的換色與定寬都是共用的（docs/spec/119 §3.1）。
	g.drawLegacyTalkBoxFields(screen, statusBoxX, statusBoxY,
		talkBoxW, talkBoxH, lines, defaultPortraitPage, g.statusBox.fields)
}

// drawTalkLineFields 畫一行 TALK 文字，句中等於某個代入欄位的那一段
// 換成那個欄位自己的顏色（docs/spec/119 §3.1）。
//
// ⭐ 欄位已經**由長到短**排序（`padTalkVars`）：每一步取「最早出現的那一個」，
// 同一個位置有兩個相符時長的贏，短的才不會把長的切一半。
func (g *game) drawTalkLineFields(screen *ebiten.Image, line string,
	fields []talkField, x, y int, ink color.RGBA) {

	for line != "" {
		best, at := -1, -1
		for i, f := range fields {
			if f.text == "" {
				continue
			}
			idx := strings.Index(line, f.text)
			if idx < 0 {
				continue
			}
			if at < 0 || idx < at {
				best, at = i, idx
			}
		}
		if best < 0 {
			g.td.Draw(screen, line, x, y, ink)
			return
		}
		if at > 0 {
			g.td.Draw(screen, line[:at], x, y, ink)
			x += textdraw.StringWidth(line[:at])
		}
		f := fields[best]
		g.td.Draw(screen, f.text, x, y, g.paletteInk(f.ink, ink))
		x += textdraw.StringWidth(f.text)
		line = line[at+len(f.text):]
	}
}
