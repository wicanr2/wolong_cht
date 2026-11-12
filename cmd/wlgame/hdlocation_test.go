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
)

func TestHDLocationGPUNativeSizeOcclusionAndCache(t *testing.T) {
	if os.Getenv("WOLONG_HD_LOCATION_GPU") == "" {
		t.Skip("separate RunGame process required")
	}
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "hd-generated"), 0700); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 192, 192))
	for y := 0; y < 192; y++ {
		for x := 0; x < 192; x++ {
			c := color.RGBA{255, 0, 0, 255}
			if x%2 != 0 {
				c = color.RGBA{0, 255, 0, 255}
			}
			img.SetRGBA(x, y, c)
		}
	}
	originals := []any{}
	samples := []any{}
	for page := 0; page < 3; page++ {
		name := fmt.Sprintf("%d.png", page)
		f, err := os.Create(filepath.Join(dir, "hd-generated", name))
		if err != nil {
			t.Fatal(err)
		}
		if err = png.Encode(f, img); err != nil {
			t.Fatal(err)
		}
		f.Close()
		b, _ := os.ReadFile(filepath.Join(dir, "hd-generated", name))
		h := sha256.Sum256(b)
		meta := map[string]any{"file": name, "width": 192, "height": 192, "sha256": hex.EncodeToString(h[:]), "native_unmodified": true}
		city := map[string]any{"id": 80 + page, "scene_page": page}
		originals = append(originals, city)
		samples = append(samples, map[string]any{"kind": "scene", "city_id": 80 + page, "city": city, "hd": map[string]any{"file": "hd-generated/" + name, "width": 192, "height": 192, "sha256": meta["sha256"]}, "hd_record": meta})
	}
	for name, data := range map[string]any{"originals.json": map[string]any{"actual_scene_count": 3, "palette_bank": 0, "cities": originals}, "manifest.json": map[string]any{"samples": samples}} {
		b, _ := json.Marshal(data)
		if err := os.WriteFile(filepath.Join(dir, name), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	pack, err := library.LoadHDLocations(dir, 3)
	if err != nil {
		t.Fatal(err)
	}
	p := &hdLocationGPUProbe{dir: dir, g: &game{hdDisplay: true, hdLocation: hdLocationState{pack: pack}, cityInfo: cityInfoState{active: true, city: 82}}}
	if err = ebiten.RunGame(p); err != nil && err != ebiten.Termination {
		t.Fatal(err)
	}
	if p.err != nil {
		t.Fatal(p.err)
	}
}

type hdLocationGPUProbe struct {
	g    *game
	dir  string
	done bool
	err  error
}

func (p *hdLocationGPUProbe) Layout(int, int) (int, int) { return screenW, screenH }
func (p *hdLocationGPUProbe) Update() error {
	if p.done {
		return ebiten.Termination
	}
	return nil
}
func (p *hdLocationGPUProbe) Draw(_ *ebiten.Image) {
	if p.done {
		return
	}
	p.done = true
	g := p.g
	src := ebiten.NewImage(screenW, screenH)
	defer src.Dispose()
	src.Fill(color.RGBA{255, 0, 255, 255})
	frame := &hdComposite{}
	frame.begin(2)
	defer frame.surface.Dispose()
	g.hdFrame = frame
	if !g.drawHDLocation(src, 0, 0) || !g.hdAssetsAvailable() {
		p.err = fmt.Errorf("page-zero native scene not drawn")
		return
	}
	rgba := func(x, y int) color.RGBA { return color.RGBAModel.Convert(frame.surface.At(x, y)).(color.RGBA) }
	if rgba(100, 620) == rgba(101, 620) {
		p.err = fmt.Errorf("native scene flattened through 96px raster")
		return
	}
	if rgba(30, 620) != (color.RGBA{255, 0, 255, 255}) || rgba(224, 620) != (color.RGBA{255, 0, 255, 255}) {
		p.err = fmt.Errorf("scene not restricted to original 96px rectangle")
		return
	}
	vector.DrawFilledRect(src, 21, 293, 8, 8, color.RGBA{0, 0, 255, 255}, false)
	frame.flush(src)
	if rgba(44, 588) != (color.RGBA{0, 0, 255, 255}) {
		p.err = fmt.Errorf("later UI behind native scene")
		return
	}
	frame.portrait(src, g.hdLocation.cache[0], 300, 10)
	if rgba(727, 40) == (color.RGBA{255, 0, 255, 255}) || rgba(728, 40) != (color.RGBA{255, 0, 255, 255}) {
		p.err = fmt.Errorf("shared native rectangle changed 64px portrait size")
		return
	}
	if !g.drawHDLocation(src, 1, 0) || !g.drawHDLocation(src, 2, 0) || len(g.hdLocation.cache) != 2 || g.hdLocation.cache[0] != nil {
		p.err = fmt.Errorf("location LRU exceeded two pages or retained oldest page")
		return
	}
	if err := os.Remove(filepath.Join(p.dir, "hd-generated/0.png")); err != nil {
		p.err = err
		return
	}
	if g.drawHDLocation(src, 0, 0) || !g.hdLocation.failed[0] {
		p.err = fmt.Errorf("damaged scene did not request original fallback")
		return
	}
	if g.drawHDLocation(src, 1, 1) {
		p.err = fmt.Errorf("unsupported palette used native scene")
		return
	}
	g.videoLCD = true
	if g.drawHDLocation(src, 1, 0) {
		p.err = fmt.Errorf("LCD used native scene")
		return
	}
	g.videoLCD = false
	g.hdDisplay = false
	if g.drawHDLocation(src, 1, 0) {
		p.err = fmt.Errorf("original display used native scene")
	}
}
