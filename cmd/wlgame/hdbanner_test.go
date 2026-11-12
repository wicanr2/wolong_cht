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

func TestHDBannerGPUCropNativeDetailAndDateOrder(t *testing.T) {
	if os.Getenv("WOLONG_HD_BANNER_GPU") == "" {
		t.Skip("separate RunGame process required")
	}
	dir := t.TempDir()
	img := image.NewRGBA(image.Rect(0, 0, 1280, 128))
	for y := 0; y < 128; y++ {
		for x := 0; x < 1280; x++ {
			c := color.RGBA{0, 0, 255, 255}
			if y >= 32 && y < 96 {
				c = color.RGBA{255, 0, 0, 255}
				if x%2 != 0 {
					c = color.RGBA{0, 255, 0, 255}
				}
			}
			img.SetRGBA(x, y, c)
		}
	}
	f, err := os.Create(filepath.Join(dir, "banner.png"))
	if err != nil {
		t.Fatal(err)
	}
	if err = png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	f.Close()
	b, _ := os.ReadFile(filepath.Join(dir, "banner.png"))
	h := sha256.Sum256(b)
	rec := map[string]any{"file": "banner.png", "width": 1280, "height": 128, "sha256": hex.EncodeToString(h[:]), "native_unmodified": true, "content_rect": [4]int{0, 32, 1280, 64}, "palette_bank": 0}
	b, _ = json.Marshal(rec)
	if err = os.WriteFile(filepath.Join(dir, "banner-selection.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	pack, err := library.LoadHDBanner(dir)
	if err != nil {
		t.Fatal(err)
	}
	p := &hdBannerGPUProbe{g: &game{hdDisplay: true, hdBanner: hdBannerState{pack: pack}}}
	if err = ebiten.RunGame(p); err != nil && err != ebiten.Termination {
		t.Fatal(err)
	}
	if p.err != nil {
		t.Fatal(p.err)
	}
}

type hdBannerGPUProbe struct {
	g    *game
	done bool
	err  error
}

func (p *hdBannerGPUProbe) Layout(int, int) (int, int) { return screenW, screenH }
func (p *hdBannerGPUProbe) Update() error {
	if p.done {
		return ebiten.Termination
	}
	return nil
}
func (p *hdBannerGPUProbe) Draw(_ *ebiten.Image) {
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
	if !g.drawHDBanner(legacy, 0) || !g.hdBanner.drawn || !g.hdAssetsAvailable() {
		p.err = fmt.Errorf("native banner did not draw")
		return
	}
	rgba := func(x, y int) color.RGBA { return color.RGBAModel.Convert(frame.surface.At(x, y)).(color.RGBA) }
	if rgba(100, 20) == rgba(101, 20) || rgba(100, 20) == (color.RGBA{0, 0, 255, 255}) {
		p.err = fmt.Errorf("banner flattened or native margins included")
		return
	}
	if rgba(100, 100) != (color.RGBA{255, 0, 255, 255}) {
		p.err = fmt.Errorf("banner changed map area")
		return
	}
	vector.DrawFilledRect(legacy, 470, 8, 8, 16, color.RGBA{220, 200, 150, 255}, false)
	frame.flush(legacy)
	if rgba(944, 22) != (color.RGBA{220, 200, 150, 255}) {
		p.err = fmt.Errorf("later live date hidden behind banner")
		return
	}
	if g.drawHDBanner(legacy, 1) {
		p.err = fmt.Errorf("wrong palette used HD banner")
		return
	}
	g.videoLCD = true
	if g.drawHDBanner(legacy, 0) {
		p.err = fmt.Errorf("LCD used color HD banner")
	}
}
