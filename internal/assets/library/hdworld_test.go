package library

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func hdWorldFixture(t *testing.T) (string, map[string]any, map[string]any) {
	t.Helper()
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "image.png"))
	if err != nil {
		t.Fatal(err)
	}
	if err = png.Encode(f, image.NewRGBA(image.Rect(0, 0, 8, 8))); err != nil {
		t.Fatal(err)
	}
	f.Close()
	b, _ := os.ReadFile(filepath.Join(dir, "image.png"))
	digest := sha256.Sum256(b)
	rec := map[string]any{"file": "image.png", "width": 8, "height": 8, "sha256": hex.EncodeToString(digest[:]), "native_unmodified": true}
	cat := map[string]any{"regions": []any{map[string]any{"id": "a", "record": "record.json", "source_world_pixel_rect": [4]int{0, 0, 32, 32}, "core_world_pixel_rect": [4]int{0, 0, 16, 16}}}}
	return dir, rec, cat
}

func writeHDWorldJSON(t *testing.T, dir, name string, v any) {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, name), b, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestHDWorldCatalogValidation(t *testing.T) {
	for _, name := range []string{"valid", "core outside source", "world overflow", "overlapping cores", "duplicate ID", "record traversal", "PNG traversal", "bad hash", "bad size", "not native", "missing PNG", "seam outside source"} {
		t.Run(name, func(t *testing.T) {
			dir, rec, cat := hdWorldFixture(t)
			row := cat["regions"].([]any)[0].(map[string]any)
			switch name {
			case "core outside source":
				row["core_world_pixel_rect"] = [4]int{24, 0, 16, 16}
			case "world overflow":
				row["source_world_pixel_rect"] = [4]int{6140, 0, 32, 32}
			case "overlapping cores":
				cat["regions"] = []any{row, map[string]any{"id": "b", "record": "record.json", "source_world_pixel_rect": [4]int{0, 0, 32, 32}, "core_world_pixel_rect": [4]int{8, 8, 16, 16}}}
			case "duplicate ID":
				cat["regions"] = []any{row, row}
			case "record traversal":
				row["record"] = "../record.json"
			case "PNG traversal":
				rec["file"] = "../image.png"
			case "bad hash":
				rec["sha256"] = "bad"
			case "bad size":
				rec["width"] = 4097
			case "not native":
				rec["native_unmodified"] = false
			case "missing PNG":
				rec["file"] = "missing.png"
			case "seam outside source":
				rec["source_world_pixel_rect"] = [4]int{0, 0, 32, 32}
				rec["core_world_pixel_rect"] = [4]int{30, 0, 16, 16}
				rec["feather_world_pixels"] = 2
				cat["seams"] = []any{map[string]any{"id": "s", "record": "record.json"}}
			}
			writeHDWorldJSON(t, dir, "record.json", rec)
			writeHDWorldJSON(t, dir, "regions.json", cat)
			pack, err := LoadHDWorld(dir)
			if name != "valid" {
				if err == nil {
					t.Fatal("accepted malformed world pack")
				}
				return
			}
			if err != nil || pack.Count() != 1 {
				t.Fatal(pack, err)
			}
			if _, err = pack.Image(0); err != nil {
				t.Fatal(err)
			}
			if _, err = pack.Image(1); err == nil {
				t.Fatal("accepted absent layer")
			}
			layers := pack.Layers()
			layers[0].Core = image.Rectangle{}
			if pack.Layers()[0].Core.Empty() {
				t.Fatal("public metadata mutates validated index")
			}
			if err = os.WriteFile(filepath.Join(dir, "image.png"), []byte("damaged"), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err = pack.Image(0); err == nil {
				t.Fatal("did not recheck PNG at decode")
			}
		})
	}
}

func TestHDWorldActualSelection(t *testing.T) {
	dir := os.Getenv("WOLONG_HD_WORLD")
	if dir == "" {
		t.Skip("local selected world art not requested")
	}
	pack, err := LoadHDWorld(dir)
	if err != nil {
		t.Fatal(err)
	}
	if pack.Count() != 57 {
		t.Fatalf("selected layer count %d", pack.Count())
	}
	area := 0
	for i, l := range pack.Layers() {
		if i < 35 {
			area += l.Core.Dx() * l.Core.Dy()
		} else if l.Feather != 8 {
			t.Fatal("seam policy changed")
		}
	}
	if area != 6144*4096 {
		t.Fatal("world coverage changed", area)
	}
	for _, i := range []int{0, 34, 56} {
		if _, err = pack.Image(i); err != nil {
			t.Fatal(err)
		}
	}
}
