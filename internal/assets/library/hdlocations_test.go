package library

import (
	"os"
	"path/filepath"
	"testing"
)

func locationFixture(t *testing.T) (string, map[string]any, map[string]any, map[string]any) {
	t.Helper()
	dir, rec, _ := hdWorldFixture(t)
	if err := os.Mkdir(filepath.Join(dir, "hd-generated"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(dir, "image.png"), filepath.Join(dir, "hd-generated/image.png")); err != nil {
		t.Fatal(err)
	}
	original := map[string]any{"actual_scene_count": 15, "palette_bank": 0, "cities": []any{map[string]any{"id": 82, "scene_page": 0}}}
	hd := map[string]any{"file": "hd-generated/image.png", "width": rec["width"], "height": rec["height"], "sha256": rec["sha256"]}
	sample := map[string]any{"kind": "scene", "city_id": 82, "city": map[string]any{"id": 82, "scene_page": 0}, "hd": hd, "hd_record": rec}
	return dir, original, sample, rec
}

func TestHDLocationsSelectionAndPageZero(t *testing.T) {
	for _, name := range []string{"valid", "archive count", "missing page", "out of range", "page mismatch", "city mismatch", "duplicate page", "record mismatch", "not native", "file traversal", "directory escape", "bad image"} {
		t.Run(name, func(t *testing.T) {
			dir, original, sample, rec := locationFixture(t)
			samples := []any{sample, map[string]any{"kind": "map", "hd": map[string]any{"file": "missing.png"}}}
			switch name {
			case "archive count":
				original["actual_scene_count"] = 14
			case "missing page":
				sample["city"] = map[string]any{"id": 82}
			case "out of range":
				original["cities"] = []any{map[string]any{"id": 82, "scene_page": 15}}
			case "page mismatch":
				sample["city"] = map[string]any{"id": 82, "scene_page": 1}
			case "city mismatch":
				sample["city"] = map[string]any{"id": 63, "scene_page": 0}
			case "duplicate page":
				samples = append(samples, sample)
			case "record mismatch":
				sample["hd"] = map[string]any{"file": "hd-generated/image.png", "width": 9, "height": 8, "sha256": rec["sha256"]}
			case "not native":
				rec["native_unmodified"] = false
			case "file traversal":
				rec["file"] = "../image.png"
				sample["hd"].(map[string]any)["file"] = "hd-generated/../image.png"
			case "directory escape":
				if err := os.Rename(filepath.Join(dir, "hd-generated"), filepath.Join(dir, "old")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(t.TempDir(), filepath.Join(dir, "hd-generated")); err != nil {
					t.Fatal(err)
				}
			case "bad image":
				if err := os.WriteFile(filepath.Join(dir, "hd-generated/image.png"), []byte("broken"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			writeHDWorldJSON(t, dir, "originals.json", original)
			writeHDWorldJSON(t, dir, "manifest.json", map[string]any{"samples": samples})
			pack, err := LoadHDLocations(dir, 15)
			if name != "valid" {
				if err == nil {
					t.Fatal("invalid scene selection accepted")
				}
				return
			}
			if err != nil || pack.Count() != 1 || pack.Bank != 0 {
				t.Fatal(pack, err)
			}
			if _, err = pack.Image(0); err != nil {
				t.Fatal("legitimate page zero rejected", err)
			}
			for _, page := range []int{-1, 1, 15} {
				if _, err = pack.Image(page); err == nil {
					t.Fatal("unselected or invalid page accepted", page)
				}
			}
			if _, err = LoadHDLocations(dir, 0); err == nil {
				t.Fatal("no original archive accepted")
			}
			if err = os.WriteFile(filepath.Join(dir, "hd-generated/image.png"), []byte("corrupt"), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err = pack.Image(0); err == nil {
				t.Fatal("native bytes not rechecked on decode")
			}
		})
	}
}

func TestHDLocationsActualApprovedSelection(t *testing.T) {
	dir, orig := os.Getenv("WOLONG_HD_LOCATIONS"), os.Getenv("WOLONG_HD_ORIG")
	if dir == "" || orig == "" {
		t.Skip("local approved scenes not requested")
	}
	lib, err := Load(orig)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range lib.Entries {
		if e.Label == "KYOGRF locations" {
			count = e.Count
		}
	}
	if count != 15 {
		t.Fatal("unexpected original archive shape", count)
	}
	pack, err := LoadHDLocations(dir, count)
	if err != nil || pack.Count() != 5 {
		t.Fatal(pack, err)
	}
	for _, page := range []int{0, 2, 5, 6, 14} {
		img, err := pack.Image(page)
		if err != nil {
			t.Fatal(err)
		}
		if img.Bounds().Dx() != 1254 || img.Bounds().Dy() != 1254 {
			t.Fatal("native scene resolution changed")
		}
	}
	if _, err = pack.Image(13); err == nil {
		t.Fatal("unpainted source page silently substituted")
	}
}
