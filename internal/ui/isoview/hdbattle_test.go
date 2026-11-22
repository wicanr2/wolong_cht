package isoview

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

func TestShipTimberAndWaterShareSource(t *testing.T) {
	const raw = 82
	const scale = library.HDBattleScale
	water := color.RGBA{R: 12, G: 90, B: 170, A: 255}
	wood := color.RGBA{R: 220, G: 150, B: 60, A: 255}
	material := battleMaterials{water: image.NewUniform(water), ship: image.NewUniform(wood)}
	material.palette[1] = color.RGBA{R: 95, G: 111, B: 111, A: 255}
	material.palette[3] = color.RGBA{R: 48, G: 65, B: 211, A: 255}
	material.palette[7] = color.RGBA{R: 180, G: 110, B: 20, A: 255}
	sub := &battle.SubTile{}
	for i := range sub.Pix {
		sub.Pix[i] = battle.Transparent
	}
	sub.Pix[0], sub.Pix[1], sub.Pix[2] = 3, 7, 1
	material.sources[raw].ground = groundShip
	material.addSource(raw, sub)
	out := image.NewRGBA(image.Rect(0, 0, battle.SubTileW*scale, displayBandRows*scale))
	material.paintBand(out, raw, image.Point{}, 0, 0)
	want := []color.RGBA{water, {R: 200, G: 130, B: 40, A: 255}, material.palette[1], {}}
	for i, ink := range want {
		if got := out.RGBAAt(i*scale, 0); got != ink {
			t.Errorf("source pixel %d: got %v, want %v", i, got, ink)
		}
	}
	material.sources[raw].waterBands = 0
	material.water = nil
	material.paintBand(out, raw, image.Point{}, 0, 0)
	if got := out.RGBAAt(0, 0); got != material.palette[3] {
		t.Errorf("water without a texture: got %v, want %v", got, material.palette[3])
	}
}

func TestBankStoneAndWaterKeepSeparateBands(t *testing.T) {
	const raw = 31
	const scale = library.HDBattleScale
	blue := color.RGBA{R: 13, G: 44, B: 201, A: 255}
	stone := color.RGBA{R: 140, G: 80, B: 60, A: 255}
	w := battleMaterials{water: image.NewUniform(blue), wallFace: image.NewUniform(stone), wallTextureParts: 10}
	w.palette[7] = color.RGBA{R: 100, G: 65, B: 30, A: 255}
	tile := &battle.SubTile{}
	for _, y := range []int{0, displayBandRows} {
		tile.Pix[y*battle.SubTileW], tile.Pix[y*battle.SubTileW+1] = 3, 7
	}
	w.sources[raw] = battleMaterialSource{tile: tile, bands: 1, wallBands: 1}
	w.addSource(raw, tile)
	out := image.NewRGBA(image.Rect(0, 0, battle.SubTileW*scale, battle.SubTileH*scale))
	w.paintBand(out, raw, image.Point{}, 0, 0)
	w.paintBand(out, raw, image.Point{}, displayBandRows, displayBandRows)
	for i, want := range [2][2]color.RGBA{{blue, stone}, {blue, w.palette[7]}} {
		for x, expected := range want {
			if got := out.RGBAAt(x*scale, i*displayBandRows*scale); got != expected {
				t.Errorf("band %d pixel %d: got %v, want %v", i, x, got, expected)
			}
		}
	}
}

func TestLadderFootMaterialPreservesFloorAndTransparentPixels(t *testing.T) {
	const raw = 178
	const scale = library.HDBattleScale
	wood := image.NewUniform(color.RGBA{R: 220, G: 160, B: 100, A: 255})
	w := battleMaterials{wood: wood, wallFace: image.White, wallTop: image.White, wallTextureParts: 5}
	w.palette[1] = color.RGBA{R: 50, G: 50, B: 50, A: 255}
	w.palette[7] = color.RGBA{R: 180, G: 100, B: 20, A: 255}
	w.palette[9] = color.RGBA{R: 230, G: 210, B: 90, A: 255}
	w.palette[15] = color.RGBA{R: 240, G: 235, B: 230, A: 255}
	tile, floor := &battle.SubTile{}, &battle.SubTile{}
	copy(tile.Pix[:], []int8{1, 7, 9, 15, battle.Transparent})
	floor.Pix[2] = 9
	w.sources[raw] = battleMaterialSource{tile: tile, bands: 0b1011, ladder: true, footGround: floor}
	tile.Pix[16*battle.SubTileW] = 7
	out := image.NewRGBA(image.Rect(0, 0, battle.SubTileW*scale, battle.SubTileH*scale))
	w.paintBand(out, raw, image.Point{}, 0, 0)
	w.paintBand(out, raw, image.Point{}, 16, 16)

	want := []color.RGBA{
		{R: 53, G: 53, B: 53, A: 255}, // Gray grain stays gray.
		{R: 192, G: 118, B: 44, A: 255},
		w.palette[9], w.palette[15], {},
	}
	for x, ink := range want {
		if got := out.RGBAAt(x*scale, 0); got != ink {
			t.Errorf("source pixel %d: got %v, want %v", x, got, ink)
		}
	}
	if got := out.RGBAAt(0, 16*scale); got != w.palette[7] {
		t.Errorf("ground band: got %v, want %v", got, w.palette[7])
	}
}

func TestSlopeMaterialLeavesGrassWaterAndShadow(t *testing.T) {
	texture := image.NewRGBA(image.Rect(0, 0, library.HDBattleWaterWidth, library.HDBattleWaterHeight))
	texture.SetRGBA(0, 0, color.RGBA{R: 223, G: 143, B: 48, A: 255})
	w := battleMaterials{slope: texture}
	for i := range w.palette {
		w.palette[i] = color.RGBA{R: 100, G: uint8(5 * i), B: 30, A: 255}
	}
	for i, original := range w.palette {
		got := w.groundColor(groundSlope, i, image.Point{})
		earth := i == 6 || i == 7 || i == 11
		if changed := got != original; changed != earth || got.A != original.A {
			t.Errorf("palette index %d: got %v with baseline %v", i, got, original)
		}
	}
}

func TestSeasonalFoliageKeepsPaletteHue(t *testing.T) {
	autumn := color.RGBA{R: 211, G: 130, A: 255}
	winter := color.RGBA{R: 243, G: 243, B: 243, A: 255}
	for _, shade := range []uint32{0x2020, 0x8080, 0xffff} {
		leaf := shadeSeasonalFoliage(autumn, shade, shade, shade)
		ratioError := int(leaf.R)*int(autumn.G) - int(leaf.G)*int(autumn.R)
		if leaf.B != 0 || leaf.A != autumn.A || ratioError < -211 || ratioError > 211 {
			t.Fatalf("autumn shade %x changed hue: %v", shade, leaf)
		}
		snow := shadeSeasonalFoliage(winter, shade, shade, shade)
		if snow.R != snow.G || snow.G != snow.B || snow.A != winter.A {
			t.Fatalf("winter shade %x lost its neutral color: %v", shade, snow)
		}
	}
	if shadeSeasonalFoliage(winter, 0x2020, 0x2020, 0x2020).R >= winter.R {
		t.Fatal("dark material must shade the original winter palette")
	}
}

func TestHDBattleCompositorGPU(t *testing.T) {
	if os.Getenv("WOLONG_HD_BATTLE_GPU") == "" {
		t.Skip("separate RunGame process required")
	}
	dir := t.TempDir()
	const scale = library.HDBattleScale
	colors := [4]color.RGBA{{220, 12, 12, 255}, {12, 220, 12, 255}, {12, 12, 220, 255}, {220, 220, 12, 255}}
	art := image.NewRGBA(image.Rect(0, 0, battle.SubTileW*scale, battle.SubTileH*scale))
	for y := 0; y < art.Bounds().Dy(); y++ {
		for x := 0; x < art.Bounds().Dx(); x++ {
			if x >= 2 { // The hole reveals terrain in the layer below.
				art.SetRGBA(x, y, colors[y/(displayBandRows*scale)])
			}
		}
	}
	var data bytes.Buffer
	if err := png.Encode(&data, art); err != nil {
		t.Fatal(err)
	}
	unit := battle.CombinedSourceTerrainTiles
	manifest := fmt.Sprintf(`{"tiles":[{"set":-1,"raw":%d,"file":"unit.png","sha256":"%x"},{"set":0,"raw":1,"file":"missing.png","sha256":"%x"}]}`, unit, sha256.Sum256(data.Bytes()), sha256.Sum256(data.Bytes()))
	if err := os.WriteFile(filepath.Join(dir, "battle-selection.json"), []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}
	pack, err := library.LoadHDBattle(dir)
	if err != nil {
		t.Fatal(err) // Neither PNG exists at index load time.
	}
	if err := os.WriteFile(filepath.Join(dir, "unit.png"), data.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	probe := &hdBattleProbe{draw: func() error {
		if pack.Available(1, false) || !pack.Available(0, false) || !pack.Available(1, true) {
			return fmt.Errorf("terrain set or shared unit lookup changed")
		}
		base := color.RGBA{10, 70, 110, 255}
		legacy := ebiten.NewImage(battle.SubTileW, battle.SubTileH)
		defer legacy.Dispose()
		legacy.Fill(base)
		v := &View{
			hdPack: pack, hdTiles: map[int]*ebiten.Image{},
			cache: map[int]*ebiten.Image{1: legacy}, sourceCache: map[int]*ebiten.Image{unit: legacy},
			sprites: &battle.Sprites{},
		}
		entries := []battleDisplayEntry{{cellCol: 1, raw: 1}, {cellCol: 1, raw: unit, lane: 1}}
		for _, multiplier := range []int{1, scale, 1} {
			v.render(nil, multiplier)
			v.drawDisplayGrid(v.buf, entries, multiplier)
			for part := range colors {
				x, y := displayUnfoldDestination(part)
				for offset := range 3 {
					want := base
					if multiplier == scale && offset >= 2 {
						want = colors[part]
					}
					if px := v.buf.At(x*multiplier+offset, y*multiplier); px != want {
						return fmt.Errorf("scale %d part %d offset %d = %v, want %v", multiplier, part, offset, px, want)
					}
				}
			}
		}
		return nil
	}}
	if err := ebiten.RunGame(probe); err != nil {
		t.Fatal(err)
	}
}

func TestHDWaterBandsGPU(t *testing.T) {
	if os.Getenv("WOLONG_HD_WATER_GPU") == "" {
		t.Skip("separate RunGame process required")
	}
	probe := &hdBattleProbe{draw: func() error {
		const scale = library.HDBattleScale
		const width = library.HDBattleWaterWidth
		const height = library.HDBattleWaterHeight
		texture := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := range height {
			for x := range width {
				texture.SetRGBA(x, y, color.RGBA{uint8(x), uint8(y), 42, 255})
			}
		}
		mask := &battle.SubTile{}
		for i := range mask.Pix {
			mask.Pix[i] = battle.Transparent
			if i%battle.SubTileW >= 4 {
				mask.Pix[i] = 3
			}
			if i%battle.SubTileW == 3 {
				mask.Pix[i] = 7
			}
		}
		water := battleMaterials{water: texture, bands: make(map[waterBandKey]*ebiten.Image)}
		hull := color.RGBA{72, 49, 21, 255}
		water.palette[7] = hull
		water.addSource(library.HDBattleWaterRaw, mask)
		defer water.reset()
		encoded := ebiten.NewImage(battle.SubTileW*scale, battle.SubTileH*scale)
		defer encoded.Dispose()
		unfolded := ebiten.NewImage(2*battle.SubTileW*scale, 2*displayBandRows*scale)
		defer unfolded.Dispose()
		bg := color.RGBA{19, 8, 6, 255}
		for _, origin := range []image.Point{image.Pt(-32, -16), image.Pt(16, 16), image.Pt(-32, -16)} {
			for _, whole := range []bool{false, true} {
				encoded.Fill(bg)
				if whole {
					water.drawFull(encoded, origin, library.HDBattleWaterRaw)
				} else {
					for part := range 4 {
						water.drawBand(encoded, origin, library.HDBattleWaterRaw, (3-part)*displayBandRows, part*displayBandRows)
					}
				}
				unfoldDisplayTile(unfolded, encoded, scale)
				pixels := image.NewRGBA(unfolded.Bounds())
				unfolded.ReadPixels(pixels.Pix)
				for y := range pixels.Bounds().Dy() {
					for x := range pixels.Bounds().Dx() {
						want := bg
						column := (x / scale) % battle.SubTileW
						if column == 3 {
							want = hull
						}
						if column >= 4 {
							want = color.RGBA{uint8((origin.X*scale + x) & 255), uint8((origin.Y*scale + y) & 127), 42, 255}
						}
						if got := pixels.RGBAAt(x, y); got != want {
							return fmt.Errorf("origin %v pixel %d,%d = %v, want %v", origin, x, y, got, want)
						}
					}
				}
			}
		}
		return nil
	}}
	if err := ebiten.RunGame(probe); err != nil {
		t.Fatal(err)
	}
}

type hdBattleProbe struct {
	draw func() error
	err  error
}

func (p *hdBattleProbe) Layout(int, int) (int, int) { return NativeW, NativeH }
func (p *hdBattleProbe) Update() error {
	if p.draw != nil {
		return nil
	}
	if p.err != nil {
		return p.err
	}
	return ebiten.Termination
}
func (p *hdBattleProbe) Draw(*ebiten.Image) {
	if p.draw != nil {
		p.err, p.draw = p.draw(), nil
	}
}
