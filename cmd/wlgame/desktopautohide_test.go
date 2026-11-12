package main

import (
	"image"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestDesktopAutoHideRespectsGestureAndOriginalOwners(t *testing.T) {
	z := [3]bool{}
	left := [3]bool{true}
	if !autoBarVisible(false, 320, 3, z, z, z, false) {
		t.Fatal("top edge does not reveal")
	}
	if autoBarVisible(true, 320, 28, z, z, z, false) {
		t.Fatal("leaving does not hide")
	}
	if autoBarVisible(false, 320, 20, z, z, z, false) {
		t.Fatal("hidden bar steals original banner")
	}
	for _, v := range []struct {
		down, previous, owned [3]bool
		locked                bool
	}{{left, z, z, false}, {z, left, z, false}, {z, z, left, false}, {z, z, z, true}} {
		if !autoBarVisible(true, 320, 100, v.down, v.previous, v.owned, v.locked) {
			t.Fatal("hid during owned gesture")
		}
		if autoBarVisible(false, 320, 3, v.down, v.previous, v.owned, v.locked) {
			t.Fatal("revealed during original/game gesture")
		}
	}
	var b desktopGesture
	for _, down := range [][3]bool{left, z} {
		a, blocked := b.pointerWithHeight(20, 3, down, func(desktopAction) bool { return false }, desktopRevealHeight)
		if a != desktopNone || !blocked {
			t.Fatal("collapsed handle dispatched action or leaked release")
		}
	}
}

func TestDesktopAutoHideRemapsPhysicalPointer(t *testing.T) {
	for _, size := range []image.Point{{1280, 812}, {1920, 1080}, {800, 1000}} {
		for _, heights := range [][2]int{{6, 28}, {28, 6}} {
			scale := math.Min(float64(size.X)/640, float64(size.Y)/float64(400+heights[0]))
			var old ebiten.GeoM
			old.Scale(scale, scale)
			old.Translate((float64(size.X)-640*scale)/2, (float64(size.Y)-float64(400+heights[0])*scale)/2)
			for _, point := range [][2]float64{{100, 3}, {320, 30}, {600, 390}} {
				x, y := remapDesktopPoint(point[0], point[1], old, image.Rect(0, 0, size.X, size.Y), heights[1])
				s := math.Min(float64(size.X)/640, float64(size.Y)/float64(400+heights[1]))
				px := x*s + (float64(size.X)-640*s)/2
				py := y*s + (float64(size.Y)-float64(400+heights[1])*s)/2
				ox, oy := old.Apply(point[0], point[1])
				if math.Abs(px-ox) > .001 || math.Abs(py-oy) > .001 {
					t.Fatal("pointer moved during viewport reflow")
				}
			}
		}
	}
}
