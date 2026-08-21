package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

// 數值面板的幾何位置來自 IDA 線性位址 00017C6E 的直接 caller：
// caller 的 DX=0058h、BX=00B8h，視窗保存區從 (80,176) 開始，
// sub_17D5F 的 16×16 格位從 (88,200) 開始，迴圈為 3 列 × 6 欄。
//
// CS:7D93h 的 18 bytes 與 sub_17C6E 的 AL-52h 分派已由 IDA／原始
// DOS/V KI.EXE 證實。AL=0..9 是逐位數字；10..14 依序是乘 100、退位、
// 清零、還原初值、完成輸入。這裡保留原始 byte，讓畫面格與實際操作
// 使用同一張可回查的表，而不是把 3×6 格畫成無關的數值欄。
var amountPanelCodes = [amountPanelRows][amountPanelCols]byte{
	{0x59, 0x5A, 0x5B, 0x5D, 0x5E, 0x5E},
	{0x56, 0x57, 0x58, 0x52, 0x5F, 0x5F},
	{0x53, 0x54, 0x55, 0x5C, 0x60, 0x60},
}

type amountPanelButton struct {
	code  byte
	edit  state.AmountEdit
	digit int
}

func amountPanelButtonAt(row, col int) (amountPanelButton, bool) {
	if row < 0 || row >= amountPanelRows || col < 0 || col >= amountPanelCols {
		return amountPanelButton{}, false
	}
	code := amountPanelCodes[row][col]
	value := int(code - 0x52)
	button := amountPanelButton{code: code, digit: -1}
	switch {
	case value >= 0 && value <= 9:
		button.edit = state.AmountAppendDigit
		button.digit = value
	case value == 10:
		button.edit = state.AmountAppendHundred
	case value == 11:
		button.edit = state.AmountDeleteDigit
	case value == 12:
		button.edit = state.AmountClear
	case value == 13:
		button.edit = state.AmountRestoreInitial
	case value == 14:
		button.edit = state.AmountFinishInput
	default:
		return amountPanelButton{}, false
	}
	return button, true
}

func amountPanelButtonCell(edit state.AmountEdit, digit int) (int, int, bool) {
	for row := 0; row < amountPanelRows; row++ {
		for col := 0; col < amountPanelCols; col++ {
			button, ok := amountPanelButtonAt(row, col)
			if !ok || button.edit != edit {
				continue
			}
			if edit != state.AmountAppendDigit || button.digit == digit {
				return row, col, true
			}
		}
	}
	return 0, 0, false
}

func amountPanelButtonAtPoint(x, y int) (amountPanelButton, int, int, bool) {
	for row := 0; row < amountPanelRows; row++ {
		for col := 0; col < amountPanelCols; col++ {
			if image.Pt(x, y).In(amountPanelCellRect(row, col)) {
				button, ok := amountPanelButtonAt(row, col)
				return button, row, col, ok
			}
		}
	}
	return amountPanelButton{}, 0, 0, false
}

func amountPanelButtonLabel(button amountPanelButton) string {
	if button.edit == state.AmountAppendDigit {
		return fmt.Sprintf("%d", button.digit)
	}
	switch button.edit {
	case state.AmountAppendHundred:
		return "百"
	case state.AmountDeleteDigit:
		return "退"
	case state.AmountClear:
		return "清"
	case state.AmountRestoreInitial:
		return "初"
	case state.AmountFinishInput:
		return "完"
	default:
		return "?"
	}
}

const (
	amountPanelX     = 80
	amountPanelY     = 176
	amountPanelW     = 112
	amountPanelH     = 80
	amountFrameX     = 88
	amountFrameY     = 184
	amountPanelGridX = 88
	amountPanelGridY = 200
	amountPanelCols  = 6
	amountPanelRows  = 3
	amountPanelCellW = 16
	amountPanelCellH = 16
	amountDisplayMax = 0x7530 // sub_17C6E 的已證實輸入上限 30,000
)

var amountPanelRect = image.Rect(
	amountPanelX,
	amountPanelY,
	amountPanelX+amountPanelW,
	amountPanelY+amountPanelH,
)

func amountPanelCellRect(row, col int) image.Rectangle {
	return image.Rect(
		amountPanelGridX+col*amountPanelCellW,
		amountPanelGridY+row*amountPanelCellH,
		amountPanelGridX+(col+1)*amountPanelCellW,
		amountPanelGridY+(row+1)*amountPanelCellH,
	)
}

func displayAmountDigits(value int) string {
	if value < 0 {
		value = 0
	}
	if value > amountDisplayMax {
		value = amountDisplayMax
	}
	return fmt.Sprintf("%0*d", amountPanelCols, value)
}

// drawAmountPanel 是兩種事件選單共用的原版 3×6 數值選取器。
func (g *game) drawAmountPanel(screen *ebiten.Image, current, initial int, selected bool) {
	if g.amountFrame != nil {
		// DOS/V sub_17D0D 的目的座標是 (88,184)，資源尺寸為 96×64；
		// sub_19796 另保存外圍 (80,176) 的 112×80 背景，供 modal 結束時還原。
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(amountFrameX, amountFrameY)
		screen.DrawImage(g.amountFrame, op)
	} else {
		g.chrome.Window(screen, amountPanelX, amountPanelY, amountPanelW, amountPanelH, chrome.Menu)
		g.drawAmountPanelFallbackButtons(screen, selected)
	}
	// 原版 sub_17C6E 在每次等待輸入前以 sub_1062F 重繪目前 SI；
	// 這個動態值仍疊在 DOS/V 靜態內框上。初值／上限是 state 語意，
	// 不在此新增原版沒有的面板外文字。
	g.td.Draw(screen, displayAmountDigits(current), amountPanelGridX,
		amountPanelY+chrome.Tile+1, chrome.Paper)
	if selected {
		g.drawDOSVAmountCursor(screen)
	}
}

// drawAmountPanelFallbackButtons 只在原始 96×64 資源缺失時使用；有資源
// 時絕不覆蓋它，因為 DOS/V 內框的下半部已包含真正的 3×6 button glyph。
func (g *game) drawAmountPanelFallbackButtons(screen *ebiten.Image, selected bool) {
	for row := 0; row < amountPanelRows; row++ {
		for col := 0; col < amountPanelCols; col++ {
			cell := amountPanelCellRect(row, col)
			button, _ := amountPanelButtonAt(row, col)
			fill := color.RGBA{0, 20, 70, 255}
			border := color.RGBA{120, 150, 180, 255}
			if selected && g.amountCursorActive &&
				row == g.amountCursorRow && col == g.amountCursorCol {
				fill = color.RGBA{20, 110, 70, 255}
				border = color.RGBA{255, 223, 154, 255}
			}
			vector.DrawFilledRect(screen, float32(cell.Min.X), float32(cell.Min.Y),
				float32(cell.Dx()), float32(cell.Dy()), fill, false)
			vector.StrokeRect(screen, float32(cell.Min.X), float32(cell.Min.Y),
				float32(cell.Dx()), float32(cell.Dy()), 1, border, false)
			if selected && g.amountCursorActive &&
				row == g.amountCursorRow && col == g.amountCursorCol {
				vector.StrokeRect(screen, float32(cell.Min.X+1), float32(cell.Min.Y+1),
					float32(cell.Dx()-2), float32(cell.Dy()-2), 1,
					color.RGBA{255, 223, 154, 255}, false)
			}
			g.td.Draw(screen, amountPanelButtonLabel(button), cell.Min.X+4,
				cell.Min.Y+1, chrome.Paper)
		}
	}
}

// drawDOSVAmountCursor 使用原版 16×16 白框／紅填箭頭。若目前沒有有效
// 的畫面滑鼠位置，鍵盤 fallback 會把箭頭放到同一個 raw action 格位；
// 這只決定呈現位置，不改變 amountPanelButtonAtPoint 的命中規則。
func (g *game) drawDOSVAmountCursor(screen *ebiten.Image) {
	if g.cursorImage == nil || !g.amountCursorActive {
		return
	}
	x, y := ebiten.CursorPosition()
	if !image.Pt(x, y).In(amountPanelRect) {
		cell := amountPanelCellRect(g.amountCursorRow, g.amountCursorCol)
		x, y = cell.Min.X, cell.Min.Y
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.cursorImage, op)
}
