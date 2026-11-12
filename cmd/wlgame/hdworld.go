package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

type hdWorldState struct {
	pack       *library.HDWorld
	cache      map[int]*ebiten.Image
	failed     map[int]bool
	lru        []int
	bytes      int
	drawn      int
	drawOrigin [2]int
	drawZoom   int
}

func (g *game) hdAssetsAvailable() bool {
	return g.hdPortraits.Count() > 0 || g.hdWorld.pack.Count() > 0 || g.hdBanner.pack != nil || g.hdLocation.pack.Count() > 0
}

func (g *game) nativeWorldLayer(index int, layer library.HDWorldLayer) *ebiten.Image {
	q := &g.hdWorld
	if q.failed[index] {
		return nil
	}
	if native := q.cache[index]; native != nil {
		q.touch(index)
		return native
	}
	img, err := q.pack.Image(index)
	if err != nil {
		if q.failed == nil {
			q.failed = map[int]bool{}
		}
		q.failed[index] = true
		log.Printf("HD world fallback (%s): %v", layer.ID, err)
		return nil
	}
	if layer.Feather > 0 {
		img = featherWorldLayer(img, layer)
	}
	if q.cache == nil {
		q.cache = map[int]*ebiten.Image{}
	}
	size := layer.Width * layer.Height * 4
	for len(q.lru) > 0 && (len(q.cache) >= 16 || q.bytes+size > 128<<20) {
		old := q.lru[0]
		q.lru = q.lru[1:]
		bounds := q.cache[old].Bounds()
		q.bytes -= bounds.Dx() * bounds.Dy() * 4
		q.cache[old].Dispose()
		delete(q.cache, old)
	}
	native := ebiten.NewImageFromImage(img)
	q.cache[index] = native
	q.bytes += size
	q.touch(index)
	return native
}

func (q *hdWorldState) touch(index int) {
	for i, v := range q.lru {
		if v == index {
			q.lru = append(q.lru[:i], q.lru[i+1:]...)
			break
		}
	}
	q.lru = append(q.lru, index)
}

// Opacity is a runtime view mask; the selected native PNG bytes stay untouched.
func featherWorldLayer(src image.Image, layer library.HDWorldLayer) image.Image {
	b := src.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	fade := func(v, lo, hi float64) float64 {
		return max(0, min(1, min((v-lo)/float64(layer.Feather), (hi-v)/float64(layer.Feather))))
	}
	for y := 0; y < b.Dy(); y++ {
		wy := float64(layer.Source.Min.Y) + (float64(y)+.5)*float64(layer.Source.Dy())/float64(b.Dy())
		fy := fade(wy, float64(layer.Core.Min.Y), float64(layer.Core.Max.Y))
		for x := 0; x < b.Dx(); x++ {
			wx := float64(layer.Source.Min.X) + (float64(x)+.5)*float64(layer.Source.Dx())/float64(b.Dx())
			c := color.NRGBAModel.Convert(src.At(b.Min.X+x, b.Min.Y+y)).(color.NRGBA)
			c.A = uint8(math.Round(float64(c.A) * fy * fade(wx, float64(layer.Core.Min.X), float64(layer.Core.Max.X))))
			out.SetNRGBA(x, y, c)
		}
	}
	return out
}

func worldLayerVertices(layer library.HDWorldLayer, ox, oy, z, scale float64) []ebiten.Vertex {
	x0 := max(float64(layer.Core.Min.X), ox)
	y0 := max(float64(layer.Core.Min.Y), oy)
	x1 := min(float64(layer.Core.Max.X), ox+strategyMapW/z)
	y1 := min(float64(layer.Core.Max.Y), oy+strategyMapH/z)
	if x1 <= x0 || y1 <= y0 {
		return nil
	}
	points := [][2]float64{{x0, y0}, {x1, y0}, {x0, y1}, {x1, y1}}
	vertices := make([]ebiten.Vertex, 4)
	for i, p := range points {
		vertices[i] = ebiten.Vertex{DstX: float32((p[0] - ox) * z * scale), DstY: float32((strategyMapY + (p[1]-oy)*z) * scale),
			SrcX: float32((p[0] - float64(layer.Source.Min.X)) * float64(layer.Width) / float64(layer.Source.Dx())),
			SrcY: float32((p[1] - float64(layer.Source.Min.Y)) * float64(layer.Height) / float64(layer.Source.Dy())), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1}
	}
	return vertices
}

func (g *game) drawWorldRaster(screen *ebiten.Image, src image.Image) {
	z := g.mapZoomScale()
	sx, sy := 0, 0
	if g.mapZoomPercent() != 100 {
		sx, sy = g.camSubX, g.camSubY
	}
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	op.GeoM.Scale(z, z)
	op.GeoM.Translate(-float64(sx)*z, strategyMapY-float64(sy)*z)
	dst := screen.SubImage(image.Rect(0, strategyMapY, strategyMapW, screenH)).(*ebiten.Image)
	if gpu, ok := src.(*ebiten.Image); ok {
		dst.DrawImage(gpu, op)
	} else {
		dst.DrawImage(ebiten.NewImageFromImage(src), op)
	}
}

func (g *game) drawHDWorld(screen *ebiten.Image) bool {
	if g.hdFrame == nil || !g.hdDisplay || g.hdWorld.pack.Count() == 0 || g.paletteBank() != 0 || g.videoLCD {
		return false
	}
	ox, oy := g.mapOrigin()
	z := g.mapZoomScale()
	sx, sy := int(ox)-g.camX*16, int(oy)-g.camY*16
	cols := int(math.Ceil((strategyMapW/z + float64(sx)) / 16))
	rows := int(math.Ceil((strategyMapH/z + float64(sy)) / 16))
	ground, err := g.lib.RenderWorld(g.camX, g.camY, cols, rows, 0)
	if err != nil {
		return false
	}
	g.drawWorldRaster(screen, ground)
	g.hdFrame.flush(screen)
	for index, layer := range g.hdWorld.pack.Layers() {
		vertices := worldLayerVertices(layer, ox, oy, z, g.hdFrame.scale)
		if len(vertices) == 0 {
			continue
		}
		if native := g.nativeWorldLayer(index, layer); native != nil {
			g.hdFrame.surface.DrawTriangles(vertices, []uint16{0, 1, 2, 1, 3, 2}, native, &ebiten.DrawTrianglesOptions{Filter: ebiten.FilterLinear})
			g.hdWorld.drawn++
		}
	}
	markers, err := g.lib.World.RenderMarkers(g.lib.Tiles, g.lib.MCH, g.lib.Palette, 0, g.camX, g.camY, cols, rows, g.cityMarks(), g.corpsMarks())
	if err != nil {
		return false
	}
	w, h := cols*16, strategyMapY+rows*16
	q := &g.worldZoom
	if q.layer == nil || q.layer.Bounds().Dx() != w || q.layer.Bounds().Dy() != h {
		if q.layer != nil {
			q.layer.Dispose()
		}
		q.layer = ebiten.NewImage(w, h)
	}
	q.layer.Clear()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, strategyMapY)
	q.layer.DrawImage(ebiten.NewImageFromImage(markers), op)
	q.drawing, q.sourceW, q.sourceH = true, w, h
	g.drawDisasterOverlay(q.layer)
	if g.mapZoomPercent() != 100 {
		g.drawMapPickCursor(q.layer)
	}
	q.drawing = false
	g.drawWorldRaster(screen, q.layer.SubImage(image.Rect(0, strategyMapY, w, h)))
	g.hdWorld.drawOrigin = [2]int{int(ox), int(oy)}
	g.hdWorld.drawZoom = g.mapZoomPercent()
	return true
}
