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
// 清零、**設成上限**、完成輸入。這裡保留原始 byte，讓畫面格與實際操作
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
		button.edit = state.AmountSetMax
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

func amountPanelButtonAtPoint(ax, ay, x, y int) (amountPanelButton, int, int, bool) {
	for row := 0; row < amountPanelRows; row++ {
		for col := 0; col < amountPanelCols; col++ {
			if image.Pt(x, y).In(amountPanelCellRect(ax, ay, row, col)) {
				button, ok := amountPanelButtonAt(row, col)
				return button, row, col, ok
			}
		}
	}
	return amountPanelButton{}, 0, 0, false
}

// amountPanelButtonLabel 只在原始 96×64 資源缺失時用得到。
// 字樣照原版圖庫上解出來的那幾個（docs/re/48 §6），不要自己另取名字。
func amountPanelButtonLabel(button amountPanelButton) string {
	if button.edit == state.AmountAppendDigit {
		return fmt.Sprintf("%d", button.digit)
	}
	switch button.edit {
	case state.AmountAppendHundred:
		return "00"
	case state.AmountDeleteDigit:
		return "◀"
	case state.AmountClear:
		return "消"
	case state.AmountSetMax:
		return "大"
	case state.AmountFinishInput:
		return "定"
	default:
		return "?"
	}
}

// 幾何全部由**錨點**推出來（docs/spec/78 §1.2）：
//
//	存／還原區  (錨點 − 8)        112 × 80    sub_19796 ／ sub_197C3
//	外框 blit   錨點               96 × 64     sub_17D0D
//	3×6 格      (錨點X, 錨點Y+16) 每格 16×16  sub_17D5F 的 add bx,10h
//
// ⚠ **不要寫死 (88,184)。** 那只是事件 2／3／4／5 的錨點；
// 財政的四個熱區傳的是 (296,184)。
const (
	amountAnchorEventX, amountAnchorEventY = 88, 184
	amountPanelMargin                      = 8
	amountPanelW                           = 112
	amountPanelH                           = 80
	amountGridDY                           = 16
	amountPanelCols                        = 6
	amountPanelRows                        = 3
	amountPanelCellW                       = 16
	amountPanelCellH                       = 16
	amountDisplayMax                       = 0x7530 // 事件 2／3／4／5 的上限 30,000
)

// amountAnchor 回傳目前這一次輸入的錨點。沒設過就是事件那一組。
func (g *game) amountAnchor() (int, int) {
	if g.amountAnchorX == 0 && g.amountAnchorY == 0 {
		return amountAnchorEventX, amountAnchorEventY
	}
	return g.amountAnchorX, g.amountAnchorY
}

func amountPanelRectAt(ax, ay int) image.Rectangle {
	x, y := ax-amountPanelMargin, ay-amountPanelMargin
	return image.Rect(x, y, x+amountPanelW, y+amountPanelH)
}

func amountPanelCellRect(ax, ay, row, col int) image.Rectangle {
	gx, gy := ax, ay+amountGridDY
	return image.Rect(
		gx+col*amountPanelCellW,
		gy+row*amountPanelCellH,
		gx+(col+1)*amountPanelCellW,
		gy+(row+1)*amountPanelCellH,
	)
}

// displayAmountValue 只做顯示夾限（負值歸零、上限 30,000）。
func displayAmountValue(value int) int {
	if value < 0 {
		return 0
	}
	if value > amountDisplayMax {
		return amountDisplayMax
	}
	return value
}

func displayAmountDigits(value int) string {
	// 原版值列不補零：值 0 畫成一個「0」右靠（2026-08-24 實機對拍
	// workplace/promo-live/parity-windows/p6-amount.png，單一白字
	// 落在欄位最右一格）。六格欄位以空白補位。
	return fmt.Sprintf("%*d", amountPanelCols, displayAmountValue(value))
}

// drawAmountPanel 是兩種事件選單共用的原版 3×6 數值選取器。
func (g *game) drawAmountPanel(screen *ebiten.Image, current int, selected bool) {
	ax, ay := g.amountAnchor()
	rect := amountPanelRectAt(ax, ay)
	// 原版在 96×64 資源外面還有一圈標準視窗框（112×80，內部剛好
	// 112−16 × 80−16 ＝ 96×64 被資源蓋滿）——2026-08-24 實機對拍
	// p6-amount.png：紅紋上下邊 ＋ 黃柱兩側。先畫框再貼資源。
	g.chrome.Window(screen, rect.Min.X, rect.Min.Y,
		rect.Dx(), rect.Dy(), chrome.Menu)
	if g.amountFrame != nil {
		// sub_17D0D 把 96×64 的資源貼在錨點上；sub_19796 另外保存外圍
		// 112×80 的背景，供 modal 結束時還原。
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(ax), float64(ay))
		screen.DrawImage(g.amountFrame, op)
	} else {
		g.drawAmountPanelFallbackButtons(screen, selected)
	}
	// 原版 sub_17C6E 在每次等待輸入前以 sub_1062F 重繪目前 SI；
	// 這個動態值仍疊在靜態內框上。上限是 state 語意，
	// 不在此新增原版沒有的面板外文字。
	// 數字用原版 8×16 字模、不補零（2026-08-24 對拍 p6-amount.png：
	// 值 0 是單一白字右靠，墨水列 185..198 ＝ mask topY 184 ＝ 錨點 Y）。
	g.drawOriginalNumber(screen, displayAmountValue(current), ax, ay,
		amountPanelCols, chrome.Paper)
	if selected {
		g.drawDOSVAmountCursor(screen)
	}
}

// drawAmountPanelFallbackButtons 只在原始 96×64 資源缺失時使用；有資源
// 時絕不覆蓋它，因為 DOS/V 內框的下半部已包含真正的 3×6 button glyph。
func (g *game) drawAmountPanelFallbackButtons(screen *ebiten.Image, selected bool) {
	ax, ay := g.amountAnchor()
	for row := 0; row < amountPanelRows; row++ {
		for col := 0; col < amountPanelCols; col++ {
			cell := amountPanelCellRect(ax, ay, row, col)
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
	ax, ay := g.amountAnchor()
	x, y := ebiten.CursorPosition()
	if !image.Pt(x, y).In(amountPanelRectAt(ax, ay)) {
		cell := amountPanelCellRect(ax, ay, g.amountCursorRow, g.amountCursorCol)
		x, y = cell.Min.X, cell.Min.Y
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.cursorImage, op)
}
