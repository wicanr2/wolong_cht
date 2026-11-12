package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// hdComposite flushes the legacy raster between native portraits. Later UI
// remains above the portrait without reconstructing modal visibility rules.
type hdComposite struct {
	surface *ebiten.Image
	scale   float64
}

func (c *hdComposite) begin(scale float64) {
	if math.IsNaN(scale) || math.IsInf(scale, 0) || scale <= 0 {
		scale = 2
	}
	c.scale = math.Max(1, math.Min(8, scale))
	w, h := int(math.Ceil(screenW*c.scale)), int(math.Ceil(screenH*c.scale))
	if c.surface == nil || c.surface.Bounds().Dx() != w || c.surface.Bounds().Dy() != h {
		if c.surface != nil {
			c.surface.Dispose()
		}
		c.surface = ebiten.NewImage(w, h)
	}
	c.surface.Clear()
}

func (c *hdComposite) flush(src *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.scale, c.scale)
	c.surface.DrawImage(src, op)
	src.Clear()
}

func (c *hdComposite) portrait(src, native *ebiten.Image, x, y int) {
	c.nativeRect(src, native, x, y, 64, 64)
}

func (c *hdComposite) nativeRect(src, native *ebiten.Image, x, y, w, h int) {
	c.flush(src)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(w)*c.scale/float64(native.Bounds().Dx()), float64(h)*c.scale/float64(native.Bounds().Dy()))
	op.GeoM.Translate(float64(x)*c.scale, float64(y)*c.scale)
	op.Filter = ebiten.FilterLinear
	c.surface.DrawImage(native, op)
}

func (g *game) nativePortrait(page int) *ebiten.Image {
	if g.hdPortraits == nil || g.hdPortraitFailed[page] {
		return nil
	}
	if img := g.hdPortraitCache[page]; img != nil {
		g.touchHDPortrait(page)
		return img
	}
	img, err := g.hdPortraits.Image(page)
	if err != nil {
		if g.hdPortraitFailed == nil {
			g.hdPortraitFailed = map[int]bool{}
		}
		g.hdPortraitFailed[page] = true
		log.Printf("HD portrait fallback (%d): %v", page, err)
		return nil
	}
	if g.hdPortraitCache == nil {
		g.hdPortraitCache = map[int]*ebiten.Image{}
	}
	if len(g.hdPortraitLRU) >= 8 {
		old := g.hdPortraitLRU[0]
		g.hdPortraitLRU = g.hdPortraitLRU[1:]
		g.hdPortraitCache[old].Dispose()
		delete(g.hdPortraitCache, old)
	}
	native := ebiten.NewImageFromImage(img)
	g.hdPortraitCache[page] = native
	g.touchHDPortrait(page)
	return native
}

func (g *game) touchHDPortrait(page int) {
	for i, p := range g.hdPortraitLRU {
		if p == page {
			g.hdPortraitLRU = append(g.hdPortraitLRU[:i], g.hdPortraitLRU[i+1:]...)
			break
		}
	}
	g.hdPortraitLRU = append(g.hdPortraitLRU, page)
}

// The original-mode output uses Ebiten's own default final-screen path.
func (s *desktopGame) DrawFinalScreen(dst ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	s.finalGeometry, s.finalBounds, s.hasFinalGeometry = geoM, dst.Bounds(), true
	s.finalScale = math.Abs(geoM.Element(0, 0))
	ebiten.DefaultDrawFinalScreen(dst, offscreen, geoM)
	if !s.hdFrameDrawn || s.hd.surface == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1/s.hd.scale, 1/s.hd.scale)
	op.GeoM.Translate(0, float64(s.barHeight()))
	op.GeoM.Concat(geoM)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(s.hd.surface, op)
}

func originalPortraitCount(g *game) int {
	if g == nil || g.lib == nil {
		return 0
	}
	for _, e := range g.lib.Entries {
		if e.Label == "KAOGRF portraits" {
			return e.Count
		}
	}
	return 0
}
