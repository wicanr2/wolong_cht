// Package mobileui contains platform-neutral layout and touch geometry for the
// Android shell. It deliberately does not import Android or Ebiten APIs.
package mobileui

import "math"

// SafeArea is expressed in physical screen pixels by the platform shell.
type SafeArea struct {
	Left, Top, Right, Bottom float64
}

// Rect is a floating-point rectangle used by the input adapter.
type Rect struct {
	X, Y, W, H float64
}

func (r Rect) Contains(x, y float64) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

// Viewport maps a fixed logical DOS/V canvas into a safe screen rectangle.
// The logical canvas is never stretched and black bars are not interactive.
type Viewport struct {
	LogicalW, LogicalH float64
	ScreenW, ScreenH   float64
	Scale              float64
	OffsetX, OffsetY   float64
}

func NewViewport(screenW, screenH, logicalW, logicalH float64, safe SafeArea) Viewport {
	availableW := math.Max(1, screenW-safe.Left-safe.Right)
	availableH := math.Max(1, screenH-safe.Top-safe.Bottom)
	logicalW = math.Max(1, logicalW)
	logicalH = math.Max(1, logicalH)
	scale := math.Min(availableW/logicalW, availableH/logicalH)
	contentW := logicalW * scale
	contentH := logicalH * scale
	return Viewport{
		LogicalW: logicalW,
		LogicalH: logicalH,
		ScreenW:  screenW,
		ScreenH:  screenH,
		Scale:    scale,
		OffsetX:  safe.Left + (availableW-contentW)/2,
		OffsetY:  safe.Top + (availableH-contentH)/2,
	}
}

func (v Viewport) ScreenRect() Rect {
	return Rect{X: v.OffsetX, Y: v.OffsetY, W: v.LogicalW * v.Scale, H: v.LogicalH * v.Scale}
}

func (v Viewport) ScreenToLogical(x, y float64) (float64, float64, bool) {
	if v.Scale <= 0 || !v.ScreenRect().Contains(x, y) {
		return 0, 0, false
	}
	return (x - v.OffsetX) / v.Scale, (y - v.OffsetY) / v.Scale, true
}

func (v Viewport) LogicalToScreen(x, y float64) (float64, float64) {
	return v.OffsetX + x*v.Scale, v.OffsetY + y*v.Scale
}

// CellAt converts a screen tap to a logical map cell. The returned cell is
// zero-based and false means the tap landed outside the map.
func (v Viewport) CellAt(screenX, screenY, originX, originY, cellW, cellH float64, cols, rows int) (int, int, bool) {
	x, y, ok := v.ScreenToLogical(screenX, screenY)
	if !ok || cellW <= 0 || cellH <= 0 {
		return 0, 0, false
	}
	col := int(math.Floor((x - originX) / cellW))
	row := int(math.Floor((y - originY) / cellH))
	if col < 0 || col >= cols || row < 0 || row >= rows {
		return 0, 0, false
	}
	return col, row, true
}

// Button uses an intentionally larger logical hitbox than the DOS/V glyph.
type Button struct {
	ID, Label string
	Bounds    Rect
}

func (b Button) Hit(x, y float64) bool { return b.Bounds.Contains(x, y) }

// DockButtons returns the prototype's three bottom-drawer controls. Each
// control is 64 logical pixels high, which is above the 48 dp minimum target
// before the platform scales the logical canvas.
func DockButtons(logicalW, logicalH float64) []Button {
	const margin = 12.0
	const gap = 8.0
	const height = 64.0
	width := (logicalW - 2*margin - 2*gap) / 3
	y := logicalH - margin - height
	return []Button{
		{ID: "continue", Label: "CONTINUE", Bounds: Rect{X: margin, Y: y, W: width, H: height}},
		{ID: "menu", Label: "MENU", Bounds: Rect{X: margin + width + gap, Y: y, W: width, H: height}},
		{ID: "save", Label: "SAVE", Bounds: Rect{X: margin + 2*(width+gap), Y: y, W: width, H: height}},
	}
}
