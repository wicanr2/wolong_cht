package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestHDWorldGeometryAndFeather(t *testing.T) {
	l := library.HDWorldLayer{Source: image.Rect(1600, 1600, 2240, 1968), Core: image.Rect(1600, 1600, 1920, 1968), Width: 1280, Height: 736}
	for _, percent := range []int{50, 100, 200} {
		v := worldLayerVertices(l, 1600, 1600, float64(percent)/100, 2)
		if len(v) != 4 || v[0].DstX != 0 || v[0].DstY != 64 || v[1].SrcX != 640 || v[3].DstY > 800 {
			t.Fatal("world/core/viewport mapping", percent, v)
		}
	}
	if v := worldLayerVertices(l, 3000, 3000, 1, 2); len(v) != 0 {
		t.Fatal("invisible region drawn")
	}
	src := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			src.SetNRGBA(x, y, color.NRGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	l = library.HDWorldLayer{Source: image.Rect(0, 0, 32, 32), Core: image.Rect(8, 8, 24, 24), Feather: 4}
	mask := featherWorldLayer(src, l)
	if _, _, _, a := mask.At(0, 0).RGBA(); a != 0 {
		t.Fatal("seam outside core remains opaque")
	}
	if _, _, _, a := mask.At(16, 16).RGBA(); a != 65535 {
		t.Fatal("seam center is not opaque")
	}
	if _, _, _, a := mask.At(8, 16).RGBA(); a == 0 || a == 65535 {
		t.Fatal("seam edge lacks opacity transition")
	}
	if src.NRGBAAt(8, 16).A != 255 {
		t.Fatal("native source image mutated")
	}
}

func TestHDWorldGPUKeepsDetailAndFallback(t *testing.T) {
	if os.Getenv("WOLONG_HD_WORLD_GPU") == "" {
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
	dir := t.TempDir()
	rows := []map[string]any{}
	for i := 0; i < 2; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 640, 736))
		for y := 0; y < 736; y++ {
			for x := 0; x < 640; x++ {
				c := color.RGBA{255, 0, 0, 255}
				if x%2 != 0 {
					c = color.RGBA{0, 255, 0, 255}
				}
				img.SetRGBA(x, y, c)
			}
		}
		name := fmt.Sprintf("%d.png", i)
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if err = png.Encode(f, img); err != nil {
			t.Fatal(err)
		}
		f.Close()
		data, _ := os.ReadFile(filepath.Join(dir, name))
		h := sha256.Sum256(data)
		record := fmt.Sprintf("%d.json", i)
		b, _ := json.Marshal(map[string]any{"file": name, "width": 640, "height": 736, "sha256": hex.EncodeToString(h[:]), "native_unmodified": true})
		if err = os.WriteFile(filepath.Join(dir, record), b, 0600); err != nil {
			t.Fatal(err)
		}
		rows = append(rows, map[string]any{"id": name, "record": record, "source_world_pixel_rect": [4]int{1600 + i*320, 1600, 320, 368}, "core_world_pixel_rect": [4]int{1600 + i*320, 1600, 320, 368}})
	}
	b, _ := json.Marshal(map[string]any{"regions": rows})
	if err = os.WriteFile(filepath.Join(dir, "regions.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	pack, err := library.LoadHDWorld(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(filepath.Join(dir, "1.png")); err != nil {
		t.Fatal(err)
	}
	p := &hdWorldGPUProbe{g: &game{world: w, lib: lib, camX: 100, camY: 100, hdDisplay: true, hdWorld: hdWorldState{pack: pack}}}
	if err = ebiten.RunGame(p); err != nil && err != ebiten.Termination {
		t.Fatal(err)
	}
	if p.err != nil {
		t.Fatal(p.err)
	}
}

type hdWorldGPUProbe struct {
	g    *game
	done bool
	err  error
}

func (p *hdWorldGPUProbe) Layout(int, int) (int, int) { return screenW, screenH }
func (p *hdWorldGPUProbe) Update() error {
	if p.done {
		return ebiten.Termination
	}
	return nil
}
func (p *hdWorldGPUProbe) Draw(_ *ebiten.Image) {
	if p.done {
		return
	}
	p.done = true
	g := p.g
	legacy := ebiten.NewImage(screenW, screenH)
	defer legacy.Dispose()
	legacy.Fill(color.RGBA{255, 0, 255, 255})
	frame := &hdComposite{}
	frame.begin(2)
	defer frame.surface.Dispose()
	g.hdFrame = frame
	if !g.drawHDWorld(legacy) || g.hdWorld.drawn != 1 || !g.hdWorld.failed[1] {
		p.err = fmt.Errorf("per-region fallback failed")
		return
	}
	rgba := func(x, y int) color.RGBA { return color.RGBAModel.Convert(frame.surface.At(x, y)).(color.RGBA) }
	if rgba(200, 200) == rgba(201, 200) {
		p.err = fmt.Errorf("native stripes flattened through the logical map raster")
		return
	}
	if rgba(20, 20) != (color.RGBA{255, 0, 255, 255}) {
		p.err = fmt.Errorf("HD map overwrote fixed banner")
		return
	}
	if !g.hdAssetsAvailable() {
		p.err = fmt.Errorf("world-only pack must enable HD mode")
		return
	}
	g.videoLCD = true
	if g.drawHDWorld(legacy) {
		p.err = fmt.Errorf("LCD used unsupported HD world art")
		return
	}
	g.videoLCD = false
	month := g.world.Clock.Month
	for m := 1; m <= 12; m++ {
		g.world.Clock.Month = m
		if g.paletteBank() != 0 && g.drawHDWorld(legacy) {
			p.err = fmt.Errorf("unsupported seasonal palette used HD world art")
			return
		}
	}
	g.world.Clock.Month = month
	vector.DrawFilledRect(legacy, 100, 100, 16, 16, color.RGBA{0, 0, 255, 255}, false)
	frame.flush(legacy)
	if rgba(220, 220) != (color.RGBA{0, 0, 255, 255}) {
		p.err = fmt.Errorf("later UI behind native map")
		return
	}
	original, err := g.lib.RenderWorldMarked(100, 100, 40, 23, 0, g.cityMarks(), g.corpsMarks())
	if err != nil {
		p.err = err
		return
	}
	if rgba(1100, 200) != original.RGBAAt(550, 68) {
		p.err = fmt.Errorf("failed region did not retain original map: %v versus %v", rgba(1100, 200), original.RGBAAt(550, 68))
		return
	}
	if len(g.hdWorld.cache) > 16 || g.hdWorld.bytes > 128<<20 {
		p.err = fmt.Errorf("world cache exceeded limits")
	}
}
