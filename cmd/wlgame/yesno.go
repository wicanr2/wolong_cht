package main

// ＹＥＳ／ＮＯ 對話框（docs/spec/26）。
//
// 原版是**可移動位置**的元件：`sub_18DC8` 的位置由呼叫端的 dx／bx 決定，
// 問題文字也由呼叫端給。而且它**不登記熱區**——直接拿滑鼠座標減去左上角
// 再除 8 來分辨 ＹＥＳ／ＮＯ，中間那 8 px 的縫是「兩個都不算」。

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 版面出自原版（docs/spec/26）：框的大小來自 `sub_10C14(cx=60Dh)`，
// 內容是顯示清單場景 7，貼在框內 (+8, +8)。
const (
	yesNoW, yesNoH = 208, 96

	// 場景 7 的原點相對於框，以及它的內容。
	yesNoSceneDX, yesNoSceneDY = 8, 8
	yesNoRuleDY                = 31 // op 05：(0,23) 長 191
	yesNoRuleW                 = 191
	yesNoQuestionDY            = 12 // sub_18DC8 的 `add bx, 0Ch`

	yesNoBoxDX  = 40 // 兩個選項框，相對於框的左上角
	yesNoYesDY  = 40
	yesNoNoDY   = 64
	yesNoBoxW   = 128
	yesNoBoxH   = 16
	yesNoTextDX = 80

	// 命中判定的範圍：相對 (框+40, 框+40)，X < 128、Y < 72。
	yesNoHitW, yesNoHitH = 128, 72
)

// 離開確認的位置：畫面正中。原版這個對話框的位置由呼叫端決定，
// remake 的離開確認沒有別的參考點，就置中。
const (
	quitDialogX = (screenW - yesNoW) / 2
	quitDialogY = (screenH - yesNoH) / 2
)

// yesNoRect 是整個對話框的矩形。
func yesNoRect(x, y int) image.Rectangle {
	return image.Rect(x, y, x+yesNoW, y+yesNoH)
}

// hitTestYesNo 照原版 `sub_18DC8` 的算式判定點在哪一邊。
//
// 回傳 (是否命中, 是不是 ＹＥＳ)。**Y 除以 8 等於 2 的那一條是縫**，
// 兩個選項都不算——原版就是這樣寫的，不是四捨五入到最近的按鈕。
func hitTestYesNo(x, y, px, py int) (hit, yes bool) {
	dx, dy := px-(x+yesNoBoxDX), py-(y+yesNoYesDY)
	if dx < 0 || dx >= yesNoHitW || dy < 0 || dy >= yesNoHitH {
		return false, false
	}
	switch row := dy / 8; {
	case row == 2:
		return false, false
	case row < 2:
		return true, true
	default:
		return true, false
	}
}

// drawYesNo 畫一個原版版面的 ＹＥＳ／ＮＯ 對話框。
// selected 為 true 時 ＹＥＳ 那格反白——**原版沒有選取狀態**（只有滑鼠），
// remake 保留鍵盤操作才需要它。
func (g *game) drawYesNo(screen *ebiten.Image, x, y int, question string, yesSelected bool) {
	g.chrome.Window(screen, x, y, yesNoW, yesNoH, chrome.Menu)
	ink := g.paletteInk(strategyInkNormal, chrome.Paper)

	// 問題文字置中：原版由呼叫端算好 X 偏移再傳進來。
	qx := x + (yesNoW-textdraw.StringWidth(question))/2
	g.td.Draw(screen, question, qx, y+yesNoQuestionDY, ink)
	vector.DrawFilledRect(screen, float32(x+yesNoSceneDX), float32(y+yesNoRuleDY),
		yesNoRuleW, 1, ink, false)

	for _, row := range []struct {
		dy    int
		label string
		on    bool
	}{
		{yesNoYesDY, "ＹＥＳ", yesSelected},
		{yesNoNoDY, "Ｎ　Ｏ", !yesSelected},
	} {
		// 兩格是**凹槽**：底色 5、外圈 2／0、內圈 D／4（docs/re/48 §2.1）。
		g.dlSunken(screen, x+yesNoBoxDX, y+row.dy, yesNoBoxW, yesNoBoxH)
		// remake 的選取標記：文字換成選取色。**不用紅色**——
		// 紅色在這個專案裡固定表示負值（資金、上昇值）。
		col := ink
		if row.on {
			col = chrome.Select
		}
		g.td.Draw(screen, row.label, x+yesNoTextDX, y+row.dy, col)
	}
}
