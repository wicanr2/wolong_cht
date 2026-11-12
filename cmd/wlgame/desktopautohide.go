package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const desktopRevealHeight = 6

func (s *desktopGame) barHeight() int {
	if s.barVisible {
		return desktopBarHeight
	}
	return desktopRevealHeight
}

// No geometry transition may split a held or just-released mouse gesture.
func autoBarVisible(current bool, x, y int, down, previous, owned [3]bool, locked bool) bool {
	if locked || down != [3]bool{} || previous != [3]bool{} || owned != [3]bool{} {
		return current
	}
	if current {
		return x >= 0 && x < screenW && y >= 0 && y < desktopBarHeight
	}
	return x >= 0 && x < screenW && y >= 0 && y < desktopRevealHeight
}

func remapDesktopPoint(x, y float64, old ebiten.GeoM, bounds image.Rectangle, newBar int) (float64, float64) {
	px, py := old.Apply(x, y)
	scale := math.Min(float64(bounds.Dx())/screenW, float64(bounds.Dy())/float64(screenH+newBar))
	if scale <= 0 {
		return x, y
	}
	ox := (float64(bounds.Dx()) - screenW*scale) / 2
	oy := (float64(bounds.Dy()) - float64(screenH+newBar)*scale) / 2
	return (px - ox) / scale, (py - oy) / scale
}

func (s *desktopGame) setBarVisible(visible bool) {
	if s.barVisible == visible {
		return
	}
	old := s.barHeight()
	s.barVisible = visible
	next := s.barHeight()
	if desktopPointer.active {
		x, y := float64(desktopPointer.x), float64(desktopPointer.y+old)
		if !ebiten.IsFullscreen() && !ebiten.IsWindowMaximized() {
			w, h := ebiten.WindowSize()
			scale := math.Min(float64(w)/screenW, float64(h)/float64(screenH+old))
			ebiten.SetWindowSize(w, h+int(math.Round(float64(next-old)*scale)))
		} else if s.hasFinalGeometry {
			x, y = remapDesktopPoint(x, y, s.finalGeometry, s.finalBounds, next)
		}
		desktopPointer.x = max(0, min(screenW-1, int(math.Round(x))))
		desktopPointer.y = max(0, min(screenH+next-1, int(math.Round(y)))) - next
		desktopPointer.rebase = 2
	} else if !ebiten.IsFullscreen() && !ebiten.IsWindowMaximized() {
		w, h := ebiten.WindowSize()
		scale := math.Min(float64(w)/screenW, float64(h)/float64(screenH+old))
		ebiten.SetWindowSize(w, h+int(math.Round(float64(next-old)*scale)))
	}
	desktopGameOffset = next
}
