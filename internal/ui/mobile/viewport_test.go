package mobileui

import (
	"math"
	"testing"
)

func TestViewportKeepsLogicalCanvasAndRejectsBlackBars(t *testing.T) {
	v := NewViewport(2400, 1080, 640, 400, SafeArea{Left: 80, Right: 80, Top: 24, Bottom: 24})
	if got, want := v.Scale, 2.58; math.Abs(got-want) > 1e-9 {
		t.Fatalf("scale = %v, want %v", got, want)
	}
	if _, _, ok := v.ScreenToLogical(80, 24); ok {
		t.Fatal("safe-area corner outside the centered viewport must not hit the game")
	}
	x, y := v.LogicalToScreen(320, 200)
	lx, ly, ok := v.ScreenToLogical(x, y)
	if !ok || math.Abs(lx-320) > 1e-9 || math.Abs(ly-200) > 1e-9 {
		t.Fatalf("round trip = (%v,%v,%v)", lx, ly, ok)
	}
}

func TestViewportMapsDOSVMapCell(t *testing.T) {
	v := NewViewport(1280, 800, 640, 400, SafeArea{})
	x, y := v.LogicalToScreen(16*7+8, 32+32+16*9+8)
	col, row, ok := v.CellAt(x, y, 0, 64, 16, 16, 27, 21)
	if !ok || col != 7 || row != 9 {
		t.Fatalf("cell = (%d,%d,%v), want (7,9,true)", col, row, ok)
	}
	if _, _, ok := v.CellAt(0, 0, 0, 64, 16, 16, 27, 21); ok {
		t.Fatal("tap above the map must not hit a cell")
	}
}

func TestDockButtonsHaveLargeNonOverlappingHitboxes(t *testing.T) {
	buttons := DockButtons(640, 496)
	if len(buttons) != 3 {
		t.Fatalf("button count = %d", len(buttons))
	}
	for i, b := range buttons {
		if b.Bounds.W < 48 || b.Bounds.H < 48 {
			t.Fatalf("button %d is below the touch target: %#v", i, b.Bounds)
		}
		if i > 0 && buttons[i-1].Bounds.X+buttons[i-1].Bounds.W > b.Bounds.X {
			t.Fatalf("button %d overlaps previous button", i)
		}
	}
}
