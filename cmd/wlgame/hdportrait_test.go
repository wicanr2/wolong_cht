package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

func TestHDCompositeKeepsNativeDetailAndLaterOcclusion(t *testing.T) {
	if os.Getenv("WOLONG_HD_GPU_PROBE") == "" {
		t.Skip("Run separately with WOLONG_HD_GPU_PROBE=1 so GPU reads occur inside RunGame")
	}
	probe := &hdPortraitGPUProbe{}
	if err := ebiten.RunGame(probe); err != nil && err != ebiten.Termination {
		t.Fatal(err)
	}
	if probe.err != nil {
		t.Fatal(probe.err)
	}
}

type hdPortraitGPUProbe struct {
	done bool
	err  error
}

func (p *hdPortraitGPUProbe) Layout(int, int) (int, int) { return screenW, screenH }
func (p *hdPortraitGPUProbe) Update() error {
	if p.done {
		return ebiten.Termination
	}
	return nil
}
func (p *hdPortraitGPUProbe) Draw(_ *ebiten.Image) {
	if p.done {
		return
	}
	p.done = true
	legacy := ebiten.NewImage(screenW, screenH)
	defer legacy.Dispose()
	legacy.Fill(color.Black)
	native := image.NewRGBA(image.Rect(0, 0, 128, 128))
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			if x%2 == 0 {
				native.Set(x, y, color.RGBA{255, 0, 0, 255})
			} else {
				native.Set(x, y, color.RGBA{0, 255, 0, 255})
			}
		}
	}
	sprite := ebiten.NewImageFromImage(native)
	defer sprite.Dispose()
	var frame hdComposite
	frame.begin(2)
	defer frame.surface.Dispose()
	frame.portrait(legacy, sprite, 20, 30)
	// An opaque modal and a cursor drawn later must cover the native portrait.
	vector.DrawFilledRect(legacy, 20, 30, 16, 16, color.RGBA{0, 0, 255, 255}, false)
	vector.DrawFilledRect(legacy, 60, 70, 1, 1, color.White, false)
	frame.flush(legacy)
	rgba := func(x, y int) color.RGBA { return color.RGBAModel.Convert(frame.surface.At(x, y)).(color.RGBA) }
	if rgba(45, 65) != (color.RGBA{0, 0, 255, 255}) {
		p.err = fmt.Errorf("later modal behind HD portrait: %v", rgba(45, 65))
		return
	}
	if rgba(120, 140) != (color.RGBA{255, 255, 255, 255}) {
		p.err = fmt.Errorf("cursor behind HD portrait")
		return
	}
	if rgba(80, 100) == rgba(81, 100) {
		p.err = fmt.Errorf("native one-pixel stripes collapsed through a 64px raster")
		return
	}
}

func TestHDDisplayPreferencesBackwardCompatible(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "preferences.json")
	defaults := desktopPreferences{Version: 1, Language: "zh-hant"}
	if defaults.HDDisplay {
		t.Fatal("first launch must be original")
	}
	if err := os.WriteFile(path, []byte(`{"version":1,"language":"zh-hant"}`), 0600); err != nil {
		t.Fatal(err)
	}
	p, err := readDesktopPreferences(path, defaults)
	if err != nil || p.HDDisplay {
		t.Fatalf("legacy preferences: %v %v", p, err)
	}
	hdDefaults := defaults
	hdDefaults.HDDisplay = true
	if old, err := readDesktopPreferences(path, hdDefaults); err != nil || old.HDDisplay {
		t.Fatal("missing HD field must still mean original", err)
	}
	p.HDDisplay = true
	if err = writeDesktopPreferences(path, p); err != nil {
		t.Fatal(err)
	}
	p, err = readDesktopPreferences(path, defaults)
	if err != nil || !p.HDDisplay {
		t.Fatal("HD choice not restored", err)
	}
	g := &game{}
	if err = g.applyDesktopPreferences(p, map[string]bool{"lang": true}); err != nil {
		t.Fatal(err)
	}
	if !g.hdDisplay || g.desktopEnabled(desktopDisplay) {
		t.Fatal("missing assets must preserve request and disable switching")
	}
}

func TestHDDesktopSingleRowAllLanguages(t *testing.T) {
	for i, a := range desktopActions {
		r := desktopButton(a)
		if r.Min.Y != 2 || r.Max.X > screenW || (i > 0 && r.Min.X < desktopButton(desktopActions[i-1]).Max.X) {
			t.Fatal("toolbar wraps or overlaps", a, r)
		}
		for _, lang := range []string{"zh-hant", "zh-hans", "ja", "en"} {
			labels := []string{}
			if a == desktopDisplay {
				labels = desktopDisplayLabels[lang]
			} else if a == desktopZoomReset {
				labels = []string{fmt.Sprintf(desktopZoomLabels[lang][1], 200)}
			} else if a == desktopZoomOut {
				labels = desktopZoomLabels[lang][:1]
			} else if a == desktopZoomIn {
				labels = desktopZoomLabels[lang][2:]
			} else if a == desktopReport {
				labels = []string{battleReportLabels[lang]["button"]}
			} else {
				labels = []string{desktopLabels[lang][int(a)-1]}
			}
			for _, label := range labels {
				width := 0
				for _, ch := range label {
					if ch < 128 {
						width += 8
					} else {
						width += 16
					}
				}
				if width > r.Dx()-8 {
					t.Fatalf("%s toolbar text exceeds safe rectangle: %q", lang, label)
				}
			}
		}
	}
}

func TestHDPortraitLCDRetainsOriginalVariant(t *testing.T) {
	orig, dir := os.Getenv("WOLONG_HD_ORIG"), os.Getenv("WOLONG_HD_PORTRAITS")
	if orig == "" || dir == "" {
		t.Skip("local original and selected artwork not requested")
	}
	lib, err := library.Load(orig)
	if err != nil {
		t.Fatal(err)
	}
	g := &game{lib: lib, hdDisplay: true, videoLCD: true, hdFrame: &hdComposite{}}
	g.hdPortraits, err = library.LoadHDPortraits(dir, originalPortraitCount(g))
	if err != nil {
		t.Fatal(err)
	}
	if g.desktopEnabled(desktopDisplay) {
		t.Fatal("LCD has no native HD variant")
	}
	canvas := ebiten.NewImage(screenW, screenH)
	defer canvas.Dispose()
	g.drawPortrait(canvas, 50, 0, 0, 4)
	if g.hdNativeDraws != 0 || len(g.hdPortraitCache) != 0 || !g.hdDisplay {
		t.Fatal("LCD changed HD preference or loaded color portrait")
	}
}
