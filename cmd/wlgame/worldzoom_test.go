package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestWorldZoomFocusAndLegacySnap(t *testing.T) {
	for _, percent := range []int{50, 75, 125, 150, 175, 200} {
		g := &game{camX: 186, camY: 102, camSubX: 7, camSubY: 9}
		ox, oy := g.mapOrigin()
		if ox != 2976 || oy != 1632 {
			t.Fatal("legacy includes residual")
		}
		ax, ay := 173, 221
		wx, wy := ox+float64(ax), oy+float64(ay)
		g.zoomMapTo(percent, ax, ay)
		ox, oy = g.mapOrigin()
		if math.Abs(ox+float64(ax)/g.mapZoomScale()-wx) > .51 || math.Abs(oy+float64(ay)/g.mapZoomScale()-wy) > .51 {
			t.Fatal("zoom moved focus", percent)
		}
		wx, wy = ox+float64(ax)/g.mapZoomScale(), oy+float64(ay)/g.mapZoomScale()
		g.zoomMapTo(100, ax, ay)
		ox, oy = g.mapOrigin()
		if math.Abs(ox+float64(ax)-wx) > 8 || math.Abs(oy+float64(ay)-wy) > 8 || g.camSubX != 0 || g.camSubY != 0 {
			t.Fatal("legacy snap", percent)
		}
	}
}

func TestWorldZoomEdgesAndFractionalScroll(t *testing.T) {
	for _, percent := range []int{50, 75, 125, 200} {
		g := &game{worldZoom: worldZoomState{percent: percent}, camX: 100, camY: 100}
		g.scrollMapPixels(-100000, -100000)
		x, y := g.mapOrigin()
		if x != 0 || y != 0 {
			t.Fatal("northwest", percent, x, y)
		}
		g.scrollMapPixels(100000, 100000)
		x, y = g.mapOrigin()
		z := g.mapZoomScale()
		if x+strategyMapW/z > 6144 || y+strategyMapH/z > 4096 || x+strategyMapW/z < 6143 || y+strategyMapH/z < 4095 {
			t.Fatal("southeast", percent, x, y)
		}
		col, row, ok := g.mapPickTileAt(639, 399)
		if !ok || col >= 384 || row >= 256 {
			t.Fatal("edge hit outside world", col, row)
		}
	}
	g := &game{camX: 100, camY: 100, worldZoom: worldZoomState{percent: 200}}
	g.scrollMapPixels(1, 1)
	x, y := g.mapOrigin()
	if x != 1600 || y != 1600 {
		t.Fatal("fraction rounded early")
	}
	g.scrollMapPixels(1, 1)
	x, y = g.mapOrigin()
	if x != 1601 || y != 1601 {
		t.Fatal("fraction lost")
	}
	g.resetMapZoom()
	if g.mapZoomPercent() != 100 || g.worldZoom.fractionX != 0 {
		t.Fatal("new world zoom reset")
	}
}

func TestWorldZoomWheelRespectsUIAndModalOwner(t *testing.T) {
	g := &game{world: &state.World{}, camX: 100, camY: 100, hud: hudCommand | hudMinimap | hudFaction}
	for _, p := range [][2]int{{20, -20}, {300, 15}, {100, 40}, {500, 100}, {500, 250}} {
		if g.zoomMapWheel(1, p[0], p[1], [3]bool{}) {
			t.Fatal("wheel leaked through UI", p)
		}
	}
	if g.zoomMapWheel(1, 200, 150, [3]bool{true}) {
		t.Fatal("wheel stole pressed gesture")
	}
	g.saveUI.active = true
	if g.zoomMapWheel(1, 200, 150, [3]bool{}) || g.desktopEnabled(desktopZoomIn) {
		t.Fatal("modal permits zoom")
	}
	g.saveUI.active = false
	if !g.zoomMapWheel(1, 200, 150, [3]bool{}) || g.mapZoomPercent() != 125 {
		t.Fatal("map wheel unavailable")
	}
	for range 10 {
		g.zoomMapWheel(1, 200, 150, [3]bool{})
	}
	if g.mapZoomPercent() != 200 || g.desktopEnabled(desktopZoomIn) {
		t.Fatal("upper bound")
	}
	for range 10 {
		g.zoomMapWheel(-1, 200, 150, [3]bool{})
	}
	if g.mapZoomPercent() != 50 || g.desktopEnabled(desktopZoomOut) {
		t.Fatal("lower bound")
	}
}

func TestWorldZoomMinimapRetainsCorners(t *testing.T) {
	g := &game{worldZoom: worldZoomState{percent: 50}}
	r := g.zoomMinimapBox()
	if r.Dx() != 44 || r.Dy() != 23 {
		t.Fatal("half scale viewport", r)
	}
	g.worldZoom.percent = 200
	r = g.zoomMinimapBox()
	if r.Dx() != 14 || r.Dy() != 6 {
		t.Fatal("double scale viewport", r)
	}
	for _, size := range []int{6, 14, 23, 44} {
		for _, original := range []int{11, 24} {
			if viewBoxSourceIndex(0, size, original) != 0 || viewBoxSourceIndex(1, size, original) != 1 || viewBoxSourceIndex(size-1, size, original) != original-1 || viewBoxSourceIndex(size-2, size, original) != original-2 {
				t.Fatal("corner modified")
			}
			for i := 0; i < size; i++ {
				s := viewBoxSourceIndex(i, size, original)
				if s < 0 || s >= original {
					t.Fatal("frame index", s)
				}
			}
		}
	}
}

func TestWorldZoomOriginalCityCoordinates(t *testing.T) {
	orig := os.Getenv("WOLONG_HD_ORIG")
	if orig == "" {
		t.Skip("local original not requested")
	}
	w, err := state.LoadScenario(filepath.Join(orig, "SINARIO.DAT"), 0)
	if err != nil {
		t.Fatal(err)
	}
	checks := 0
	for _, percent := range []int{50, 75, 100, 125, 200} {
		g := &game{world: w, worldZoom: worldZoomState{percent: percent}}
		for id, c := range w.Cities {
			g.centreCamOn(c.X, c.Y)
			ox, oy := g.mapOrigin()
			z := g.mapZoomScale()
			x := int(math.Floor((float64(c.X*16+8) - ox) * z))
			y := strategyMapY + int(math.Floor((float64(c.Y*16+8)-oy)*z))
			col, row, ok := g.mapPickTileAt(x, y)
			if !ok || col != c.X || row != c.Y {
				t.Fatalf("city %d at %d%% hit %d,%d", id, percent, col, row)
			}
			checks++
		}
	}
	t.Logf("%d registered city/zoom coordinate checks", checks)
}

func TestWorldZoomGPUMapSamplingAndClip(t *testing.T) {
	if os.Getenv("WOLONG_ZOOM_GPU_PROBE") == "" {
		t.Skip("separate RunGame process required")
	}
	orig := os.Getenv("WOLONG_HD_ORIG")
	lib, err := library.Load(orig)
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario(filepath.Join(orig, "SINARIO.DAT"), 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	full, err := lib.RenderWorldMarked(0, 0, 384, 256, 0, (&game{world: w}).cityMarks(), nil)
	if err != nil {
		t.Fatal(err)
	}
	// The source scenario contains clouds, including one outside the legacy 640px source.
	// Compose its original MCH pixels independently in full world coordinates.
	for _, o := range w.RenderDisasterObjects() {
		pattern, ok := lib.MCH.PatternFor(int(o.TypeCode), int(o.Phase))
		if !ok {
			t.Fatal("missing cloud pattern")
		}
		sprite, err := lib.MCH.RenderPattern(pattern, lib.Palette, 0)
		if err != nil {
			t.Fatal(err)
		}
		x := (o.X - sprite.Bounds().Dx()/16/2) * 16
		y := (o.Y - sprite.Bounds().Dy()/16/2) * 16
		draw.Draw(full, sprite.Bounds().Add(image.Pt(x, y)), sprite, sprite.Bounds().Min, draw.Over)
	}
	p := &worldZoomGPUProbe{g: &game{world: w, lib: lib, camX: 186, camY: 102, camSubX: 7, camSubY: 9}}
	p.source = full
	if err = ebiten.RunGame(p); err != nil && err != ebiten.Termination {
		t.Fatal(err)
	}
	if p.err != nil {
		t.Fatal(p.err)
	}
}

type worldZoomGPUProbe struct {
	g      *game
	source *image.RGBA
	done   bool
	err    error
}

func (p *worldZoomGPUProbe) Layout(int, int) (int, int) { return screenW, screenH }
func (p *worldZoomGPUProbe) Update() error {
	if p.done {
		return ebiten.Termination
	}
	return nil
}
func (p *worldZoomGPUProbe) Draw(_ *ebiten.Image) {
	if p.done {
		return
	}
	p.done = true
	dst := ebiten.NewImage(screenW, screenH)
	defer dst.Dispose()
	for _, percent := range []int{100, 50, 75, 125, 200} {
		g := p.g
		g.worldZoom.percent = percent
		dst.Fill(color.RGBA{255, 0, 255, 255})
		g.drawStrategyMap(dst)
		if dst.At(20, 20) != (color.RGBA{255, 0, 255, 255}) {
			p.err = fmt.Errorf("zoom %d overwrote fixed banner", percent)
			return
		}
		x0, y0 := g.mapOrigin()
		z := g.mapZoomScale()
		for _, xy := range [][2]int{{0, 32}, {639, 32}, {0, 399}, {639, 399}, {100, 150}, {320, 200}, {601, 383}} {
			wx := x0 + (float64(xy[0])+.5)/z
			wy := y0 + (float64(xy[1]-32)+.5)/z
			actual := color.RGBAModel.Convert(dst.At(xy[0], xy[1])).(color.RGBA)
			// Nearest sampling may select either texel at an exact boundary.
			match := false
			for _, sx := range []int{int(math.Floor(wx)), int(math.Ceil(wx) - 1)} {
				for _, sy := range []int{int(math.Floor(wy)), int(math.Ceil(wy) - 1)} {
					if actual == color.RGBAModel.Convert(p.source.At(sx, sy)).(color.RGBA) {
						match = true
					}
				}
			}
			if !match {
				p.err = fmt.Errorf("zoom %d sample %v is %v at world %v,%v", percent, xy, actual, wx, wy)
				return
			}
		}
	}
}
