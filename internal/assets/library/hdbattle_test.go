package library

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
)

func TestHDBattleLoadsOnlyRequestedTiles(t *testing.T) {
	dir := t.TempDir()
	encode := func(width, height int) []byte {
		t.Helper()
		var buffer bytes.Buffer
		if err := png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, width, height))); err != nil {
			t.Fatal(err)
		}
		return buffer.Bytes()
	}
	width, height := HDBattleScale*battle.SubTileW, HDBattleScale*battle.SubTileH
	good := encode(width, height)
	wrongSize := encode(width-1, height)
	digest := fmt.Sprintf("%x", sha256.Sum256(good))
	set, shared := 0, -1
	unit := battle.CombinedSourceTerrainTiles
	records := []hdBattleTile{
		{Set: &set, Raw: 1, File: "terrain.png", SHA: strings.ToUpper(digest)},
		{Set: &set, Raw: 2, File: "missing.png", SHA: digest},
		{Set: &set, Raw: 3, File: "wrong-hash.png", SHA: digest},
		{Set: &set, Raw: 4, File: "wrong-size.png", SHA: fmt.Sprintf("%x", sha256.Sum256(wrongSize))},
		{Set: &shared, Raw: unit, File: "sprite.png", SHA: digest},
	}
	manifest, err := json.Marshal(struct {
		Tiles []hdBattleTile `json:"tiles"`
	}{Tiles: records})
	if err != nil {
		t.Fatal(err)
	}
	write := func(name string, contents []byte) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), contents, 0600); err != nil {
			t.Fatal(err)
		}
	}
	write("battle-selection.json", manifest)
	pack, err := LoadHDBattle(dir)
	if err != nil {
		t.Fatalf("index must load before image files exist: %v", err)
	}
	write("terrain.png", good)
	write("wrong-hash.png", wrongSize)
	write("wrong-size.png", wrongSize)
	write("sprite.png", good)
	for raw, reason := range map[int]string{2: "missing.png", 3: "SHA mismatch", 4: "invalid HD dimensions"} {
		if img, err := pack.Image(set, raw); img != nil || err == nil || !strings.Contains(err.Error(), reason) {
			t.Errorf("raw %d: got image %v and error %v, want %q", raw, img != nil, err, reason)
		}
	}
	if img, err := pack.Image(1, 1); img != nil || err != nil {
		t.Errorf("terrain in an unselected set: image %v, error %v", img != nil, err)
	}
	for _, sample := range []hdBattleKey{{set, 1}, {1, unit}} {
		img, err := pack.Image(sample.set, sample.raw)
		if err != nil || img == nil || img.Bounds().Size() != image.Pt(width, height) {
			t.Errorf("usable tile %+v rejected: %v", sample, err)
		}
	}
}

func TestHDBattleWaterOnDemand(t *testing.T) {
	dir := t.TempDir()
	img := image.NewGray(image.Rect(0, 0, HDBattleWaterWidth, HDBattleWaterHeight))
	img.SetGray(0, 0, color.Gray{Y: 170})
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, img); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(encoded.Bytes())
	write := func(name, contents string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0600); err != nil {
			t.Fatal(err)
		}
	}
	water := fmt.Sprintf(`"water":{"file":"water.png","sha256":"%x"}`, sum)
	write("battle-selection.json", "{"+water+"}")
	pack, err := LoadHDBattle(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []int{209, 210, 211, 212} {
		if !pack.WaterAvailable(field, HDBattleWaterSet) {
			t.Fatalf("missing water entry for field %d", field)
		}
	}
	if pack.WaterAvailable(208, HDBattleWaterSet) || pack.WaterAvailable(213, HDBattleWaterSet) ||
		pack.WaterAvailable(210, 1) || pack.Available(HDBattleWaterSet, false) {
		t.Fatal("water entry leaked into the general tile or wrong field catalog")
	}
	if loaded, err := pack.Water(); loaded != nil || err == nil {
		t.Fatalf("missing water image: got %v, %v", loaded != nil, err)
	}
	write("water.png", encoded.String())
	if loaded, err := pack.Water(); err != nil || loaded.At(0, 0) != (color.Gray{Y: 170}) {
		t.Fatalf("on-demand water decode: %v", err)
	}
	write("battle-selection.json", fmt.Sprintf(`{%s,"tiles":[{"set":2,"raw":32,"file":"tile.png","sha256":"%x"}]}`, water, sum))
	pack, err = LoadHDBattle(dir)
	if err != nil {
		t.Fatal(err)
	}
	if pack.WaterAvailable(210, HDBattleWaterSet) || !pack.Available(HDBattleWaterSet, false) {
		t.Fatal("existing tile replacement must retain priority")
	}
}

func TestHDBattleMaterialPairsLoadOnDemand(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, HDBattleWaterWidth, HDBattleWaterHeight))
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, img); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(encoded.Bytes())
	for _, tc := range []struct {
		section, first, second string
		available              func(*HDBattle) bool
		read                   func(*HDBattle) (image.Image, image.Image, error)
	}{
		{"wall", "cap", "face", (*HDBattle).WallAvailable, (*HDBattle).Wall},
		{"ground", "grass", "road", (*HDBattle).GroundAvailable, (*HDBattle).Ground},
	} {
		t.Run(tc.section, func(t *testing.T) {
			dir := t.TempDir()
			write := func(name string, contents []byte) {
				t.Helper()
				if err := os.WriteFile(filepath.Join(dir, name), contents, 0600); err != nil {
					t.Fatal(err)
				}
			}
			manifest := fmt.Sprintf(`{"%s":{"%s":{"file":"a.png","sha256":"%x"},"%s":{"file":"b.png","sha256":"%x"}}}`,
				tc.section, tc.first, sum, tc.second, sum)
			write("battle-selection.json", []byte(manifest))
			pack, err := LoadHDBattle(dir)
			if err != nil {
				t.Fatalf("index must load before material files exist: %v", err)
			}
			if !tc.available(pack) {
				t.Fatal("missing material metadata")
			}
			write("a.png", encoded.Bytes())
			if first, second, err := tc.read(pack); first != nil || second != nil || err == nil {
				t.Fatalf("missing second file: got first %v, second %v, error %v", first != nil, second != nil, err)
			}
			write("b.png", encoded.Bytes())
			first, second, err := tc.read(pack)
			if err != nil || first == nil || second == nil || first.Bounds() != img.Bounds() || second.Bounds() != img.Bounds() {
				t.Fatalf("on-demand decode: got first %v, second %v, error %v", first != nil, second != nil, err)
			}
		})
	}
}

func TestHDBattleSingleMaterialOnDemand(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, HDBattleWaterWidth, HDBattleWaterHeight))
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, img); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(encoded.Bytes())
	for _, tc := range []struct {
		name      string
		available func(*HDBattle) bool
		read      func(*HDBattle) (image.Image, error)
	}{
		{"wood", (*HDBattle).WoodAvailable, (*HDBattle).Wood},
		{"foliage", (*HDBattle).FoliageAvailable, (*HDBattle).Foliage},
		{"slope", (*HDBattle).SlopeAvailable, (*HDBattle).Slope},
		{"roof", (*HDBattle).RoofAvailable, (*HDBattle).Roof},
		{"ship", (*HDBattle).ShipAvailable, (*HDBattle).Ship},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			manifest := fmt.Sprintf(`{"%s":{"file":"material.png","sha256":"%x"}}`, tc.name, sum)
			if err := os.WriteFile(filepath.Join(dir, "battle-selection.json"), []byte(manifest), 0600); err != nil {
				t.Fatal(err)
			}
			pack, err := LoadHDBattle(dir)
			if err != nil || !tc.available(pack) {
				t.Fatalf("material metadata must load before the PNG exists: %v", err)
			}
			if loaded, err := tc.read(pack); loaded != nil || err == nil {
				t.Fatalf("missing material image: got %v, %v", loaded != nil, err)
			}
			if err := os.WriteFile(filepath.Join(dir, "material.png"), encoded.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			loaded, err := tc.read(pack)
			if err != nil || loaded == nil || loaded.Bounds() != img.Bounds() {
				t.Fatalf("on-demand material decode: got %v, %v", loaded != nil, err)
			}
		})
	}
}
