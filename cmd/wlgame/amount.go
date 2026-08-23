package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/state"
)

const (
	amountCursorDiplomacy = 1
	amountCursorFunding   = 2
	amountCursorFinance   = 3
)

// pressedAmountDigit 將一般鍵盤數字鍵映射成狀態層的十進位數字。
// PC-98 原版掃描碼仍是原始 UI 的獨立邊界；這裡只提供跨平台 remake
// 的可重播輸入映射。
func pressedAmountDigit() (int, bool) {
	keys := [...]ebiten.Key{
		ebiten.Key0, ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4,
		ebiten.Key5, ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9,
	}
	for digit, key := range keys {
		if pressed(key) {
			return digit, true
		}
	}
	return 0, false
}

// beginAmountEditor 對應 `sub_17C6E` 開場：錨點由呼叫端給
// （docs/spec/78 §1.2），目前值一律從 0 起（原版 `xor si, si`）。
func (g *game) beginAmountEditor(owner, anchorX, anchorY int) {
	g.amountCursorOwner = owner
	g.amountAnchorX, g.amountAnchorY = anchorX, anchorY
	// 原版游標由滑鼠硬體提供；跨平台第一次進入時停在「0」格，
	// 後續移動／鍵盤操作都會留下最後選取的原生格位。
	g.amountCursorRow, g.amountCursorCol = 1, 3
	g.amountCursorActive = true
}

func (g *game) setAmountCursorAction(edit state.AmountEdit, digit int) {
	row, col, ok := amountPanelButtonCell(edit, digit)
	if !ok {
		return
	}
	g.amountCursorRow, g.amountCursorCol = row, col
	g.amountCursorActive = true
}

func (g *game) amountPointerButton() (amountPanelButton, bool) {
	ax, ay := g.amountAnchor()
	x, y := ebiten.CursorPosition()
	button, row, col, ok := amountPanelButtonAtPoint(ax, ay, x, y)
	if !ok {
		return amountPanelButton{}, false
	}
	g.amountCursorRow, g.amountCursorCol = row, col
	g.amountCursorActive = true
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return amountPanelButton{}, false
	}
	return button, true
}
