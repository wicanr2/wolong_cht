package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 幾何跟著錨點走（docs/spec/78 §1.2）。**兩個錨點都驗**——
// 只驗事件那一組的話，寫死成常數也會通過。
func TestAmountPanelGeometryFollowsAnchor(t *testing.T) {
	for _, tc := range []struct {
		name             string
		ax, ay           int
		panelX, panelY   int
		firstX, firstY   int
		lastX, lastY     int
	}{
		{"事件 2／3／4／5", amountAnchorEventX, amountAnchorEventY,
			80, 176, 88, 200, 168, 232},
		{"財政四個熱區", financeAnchorX, financeAnchorY,
			288, 176, 296, 200, 376, 232},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rect := amountPanelRectAt(tc.ax, tc.ay)
			if rect.Min.X != tc.panelX || rect.Min.Y != tc.panelY ||
				rect.Dx() != 112 || rect.Dy() != 80 {
				t.Fatalf("面板 = %v，want (%d,%d) 112x80", rect, tc.panelX, tc.panelY)
			}
			first := amountPanelCellRect(tc.ax, tc.ay, 0, 0)
			if first.Min.X != tc.firstX || first.Min.Y != tc.firstY ||
				first.Dx() != 16 || first.Dy() != 16 {
				t.Fatalf("第一格 = %v，want (%d,%d) 16x16", first, tc.firstX, tc.firstY)
			}
			last := amountPanelCellRect(tc.ax, tc.ay,
				amountPanelRows-1, amountPanelCols-1)
			if last.Min.X != tc.lastX || last.Min.Y != tc.lastY {
				t.Fatalf("最後一格 = %v，want (%d,%d)", last, tc.lastX, tc.lastY)
			}
		})
	}
}

func TestDisplayAmountDigitsIsSixColumnsAndClamped(t *testing.T) {
	for _, tt := range []struct {
		value int
		want  string
	}{
		{0, "     0"},
		{1234, "  1234"},
		{-1, "     0"},
		{999999, " 30000"},
	} {
		if got := displayAmountDigits(tt.value); got != tt.want {
			t.Errorf("displayAmountDigits(%d) = %q，want %q", tt.value, got, tt.want)
		}
	}
}

func TestAmountPanelCodesMatchOriginalDispatcher(t *testing.T) {
	wantDigits := [10][2]int{
		{1, 3}, {2, 0}, {2, 1}, {2, 2}, {1, 0},
		{1, 1}, {1, 2}, {0, 0}, {0, 1}, {0, 2},
	}
	for digit, want := range wantDigits {
		row, col, ok := amountPanelButtonCell(state.AmountAppendDigit, digit)
		if !ok || row != want[0] || col != want[1] {
			t.Fatalf("數字 %d 格位 = (%d,%d,%v)，want (%d,%d,true)",
				digit, row, col, ok, want[0], want[1])
		}
		button, ok := amountPanelButtonAt(row, col)
		if !ok || button.digit != digit || button.edit != state.AmountAppendDigit {
			t.Fatalf("數字 %d 分派 = %#v，want digit／append", digit, button)
		}
	}
	for _, tt := range []struct {
		row, col int
		edit     state.AmountEdit
	}{
		{0, 3, state.AmountDeleteDigit},
		{0, 4, state.AmountClear},
		{1, 4, state.AmountSetMax},
		{2, 3, state.AmountAppendHundred},
		{2, 4, state.AmountFinishInput},
	} {
		button, ok := amountPanelButtonAt(tt.row, tt.col)
		if !ok || button.edit != tt.edit {
			t.Fatalf("(%d,%d) = %#v，want edit %d", tt.row, tt.col, button, tt.edit)
		}
	}
}

func TestAmountPanelPointUsesOriginalGridCoordinates(t *testing.T) {
	ax, ay := amountAnchorEventX, amountAnchorEventY
	button, row, col, ok := amountPanelButtonAtPoint(ax, ay, 88, 200)
	if !ok || row != 0 || col != 0 || button.code != 0x59 || button.digit != 7 {
		t.Fatalf("首格點擊 = (%#v,%d,%d,%v)，want raw 59／7 at (0,0)", button, row, col, ok)
	}
	button, row, col, ok = amountPanelButtonAtPoint(ax, ay, 103, 215)
	if !ok || row != 0 || col != 0 || button.edit != state.AmountAppendDigit {
		t.Fatalf("首格內部點擊 = (%#v,%d,%d,%v)，未命中同一 raw 格", button, row, col, ok)
	}
	if _, _, _, ok = amountPanelButtonAtPoint(ax, ay, 87, 200); ok {
		t.Fatal("格子左側一像素不應命中")
	}
}
