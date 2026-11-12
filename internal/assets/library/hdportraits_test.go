package library

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func hdFixture(t *testing.T) (string, map[string]any) {
	t.Helper()
	dir := t.TempDir()
	im := image.NewRGBA(image.Rect(0, 0, 128, 128))
	im.Set(80, 70, color.RGBA{19, 27, 41, 255})
	var b bytes.Buffer
	if err := png.Encode(&b, im); err != nil {
		t.Fatal(err)
	}
	h := sha256.Sum256(b.Bytes())
	entry := map[string]any{"portrait_page": 0, "file": "native.png", "width": 128, "height": 128, "sha256": hex.EncodeToString(h[:])}
	if err := os.WriteFile(filepath.Join(dir, "native.png"), b.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	hdWrite(t, dir, "selected.json", entry)
	hdWrite(t, dir, "selection.json", map[string]any{"records": []string{"selected.json"}})
	return dir, entry
}
func hdWrite(t *testing.T, dir, name string, v any) {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, name), b, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestHDPortraitNativePageZeroAndChangedFile(t *testing.T) {
	dir, _ := hdFixture(t)
	pack, err := LoadHDPortraits(dir, 150)
	if err != nil {
		t.Fatal(err)
	}
	if pack.Count() != 1 {
		t.Fatal(pack.Count())
	}
	im, err := pack.Image(0)
	if err != nil || im.Bounds() != image.Rect(0, 0, 128, 128) {
		t.Fatalf("native image %v %v", im, err)
	}
	if got := color.RGBAModel.Convert(im.At(80, 70)).(color.RGBA); got != (color.RGBA{19, 27, 41, 255}) {
		t.Fatal(got)
	}
	if _, err = pack.Image(1); err == nil {
		t.Fatal("missing page must fall back")
	}
	if err = os.WriteFile(filepath.Join(dir, "native.png"), []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = pack.Image(0); err == nil {
		t.Fatal("changed file accepted")
	}
}

func TestHDPortraitRejectsInvalidGroup(t *testing.T) {
	for _, kind := range []string{"page-count", "missing-page", "page-boundary", "duplicate", "dimensions", "digest", "image-path", "record-path", "symlink"} {
		t.Run(kind, func(t *testing.T) {
			dir, e := hdFixture(t)
			count := 150
			switch kind {
			case "page-count":
				count = 0
			case "missing-page":
				delete(e, "portrait_page")
			case "page-boundary":
				e["portrait_page"] = 150
			case "duplicate":
				hdWrite(t, dir, "selection.json", map[string]any{"records": []string{"selected.json", "selected.json"}})
			case "dimensions":
				e["width"] = 4097
			case "digest":
				e["sha256"] = hex.EncodeToString(make([]byte, 32))
			case "image-path":
				e["file"] = "../native.png"
			case "record-path":
				hdWrite(t, dir, "selection.json", map[string]any{"records": []string{"../selected.json"}})
			case "symlink":
				outside := filepath.Join(t.TempDir(), "other.png")
				if err := os.WriteFile(outside, []byte("outside"), 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(outside, filepath.Join(dir, "escape.png")); err != nil {
					t.Fatal(err)
				}
				e["file"] = "escape.png"
			}
			hdWrite(t, dir, "selected.json", e)
			if pack, err := LoadHDPortraits(dir, count); err == nil || pack != nil {
				t.Fatalf("invalid group accepted %v %v", pack, err)
			}
		})
	}
}

func TestHDPortraitRealSelectedGroup(t *testing.T) {
	dir, orig := os.Getenv("WOLONG_HD_PORTRAITS"), os.Getenv("WOLONG_HD_ORIG")
	if dir == "" || orig == "" {
		t.Skip("local original and selected artwork not requested")
	}
	lib, err := Load(orig)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range lib.Entries {
		if e.Label == "KAOGRF portraits" {
			count = e.Count
		}
	}
	if count != 150 {
		t.Fatal("actual archive count", count)
	}
	pack, err := LoadHDPortraits(dir, count)
	if err != nil {
		t.Fatal(err)
	}
	if pack.Count() != 146 {
		t.Fatal("selected pages", pack.Count())
	}
	for _, page := range []int{0, 50, 138, 134} {
		im, err := pack.Image(page)
		if err != nil || im.Bounds() != image.Rect(0, 0, 1254, 1254) {
			t.Fatalf("page %d: %v", page, err)
		}
	}
	if _, err = pack.Image(149); err == nil {
		t.Fatal("unselected original page must fall back")
	}
}
