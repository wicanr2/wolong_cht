package world

import (
	"bytes"
	"image"
	"image/draw"
	"testing"
)

func TestHDMarkersPreserveOriginalComposition(t *testing.T) {
	m, ts, mch, pal := loadWorldForTest(t)
	x, y := -1, -1
	for yy := 3; yy < Height-3 && x < 0; yy++ {
		for xx := 3; xx < Width-3; xx++ {
			b, _ := m.Tile(xx, yy)
			if IsCityCentre(b) {
				x, y = xx, yy
				break
			}
		}
	}
	if x < 0 {
		t.Fatal("no source city")
	}
	for _, own := range []Ownership{OwnedBySelf, OwnedByOther, Unowned} {
		for _, capital := range []bool{false, true} {
			for _, army := range []bool{false, true} {
				marks := []CityMark{{X: x, Y: y, Own: own, Capital: capital}}
				var corps []CorpsMark
				if army {
					corps = []CorpsMark{{X: x, Y: y, Tile: CorpsTile(0, 4)}}
				}
				ground, err := m.Render(ts, pal, 0, x-3, y-3, 7, 7)
				if err != nil {
					t.Fatal(err)
				}
				original, err := m.RenderMarked(ts, mch, pal, 0, x-3, y-3, 7, 7, marks, corps)
				if err != nil {
					t.Fatal(err)
				}
				over, err := m.RenderMarkers(ts, mch, pal, 0, x-3, y-3, 7, 7, marks, corps)
				if err != nil {
					t.Fatal(err)
				}
				if over.RGBAAt(0, 0).A != 0 {
					t.Fatal("overlay contains unrelated terrain")
				}
				draw.Draw(ground, ground.Bounds(), over, image.Point{}, draw.Over)
				if !bytes.Equal(ground.Pix, original.Pix) {
					t.Fatal("HD overlay changes original ownership/capital/corps composition", own, capital, army)
				}
			}
		}
	}
}
