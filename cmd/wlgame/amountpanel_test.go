package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestAmountPanelKeepsVerifiedNativeGeometry(t *testing.T) {
	if got, want := amountPanelRect.Min.X, 80; got != want {
		t.Fatalf("面板 X = %d，want %d", got, want)
	}
	if got, want := amountPanelRect.Min.Y, 176; got != want {
		t.Fatalf("面板 Y = %d，want %d", got, want)
	}
	if got, want := amountPanelRect.Dx(), 112; got != want {
		t.Fatalf("面板寬 = %d，want %d", got, want)
	}
	if got, want := amountPanelRect.Dy(), 80; got != want {
		t.Fatalf("面板高 = %d，want %d", got, want)
	}

	first := amountPanelCellRect(0, 0)
	if first.Min.X != 88 || first.Min.Y != 200 || first.Dx() != 16 || first.Dy() != 16 {
		t.Fatalf("第一格幾何 = %v，want (88,200) 16x16", first)
	}
	last := amountPanelCellRect(amountPanelRows-1, amountPanelCols-1)
	if last.Min.X != 168 || last.Min.Y != 232 || last.Dx() != 16 || last.Dy() != 16 {
		t.Fatalf("最後一格幾何 = %v，want (168,232) 16x16", last)
	}
}

func TestDisplayAmountDigitsIsSixColumnsAndClamped(t *testing.T) {
	for _, tt := range []struct {
		value int
		want  string
	}{
		{0, "000000"},
		{1234, "001234"},
		{-1, "000000"},
		{999999, "030000"},
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
		{1, 4, state.AmountRestoreInitial},
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
	button, row, col, ok := amountPanelButtonAtPoint(88, 200)
	if !ok || row != 0 || col != 0 || button.code != 0x59 || button.digit != 7 {
		t.Fatalf("首格點擊 = (%#v,%d,%d,%v)，want raw 59／7 at (0,0)", button, row, col, ok)
	}
	button, row, col, ok = amountPanelButtonAtPoint(103, 215)
	if !ok || row != 0 || col != 0 || button.edit != state.AmountAppendDigit {
		t.Fatalf("首格內部點擊 = (%#v,%d,%d,%v)，未命中同一 raw 格", button, row, col, ok)
	}
	if _, _, _, ok = amountPanelButtonAtPoint(amountPanelGridX-1, amountPanelGridY); ok {
		t.Fatal("格子左側一像素不應命中")
	}
}
