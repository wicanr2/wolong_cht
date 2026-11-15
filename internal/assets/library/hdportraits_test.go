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
	for _, kind := range []string{"page-count", "missing-page", "page-boundary", "duplicate", "dimensions", "digest", "image-path", "record-path"} {
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
				e["sha256"] = "invalid"
			case "image-path":
				e["file"] = "../native.png"
			case "record-path":
				hdWrite(t, dir, "selection.json", map[string]any{"records": []string{"../selected.json"}})
			}
			hdWrite(t, dir, "selected.json", e)
			if pack, err := LoadHDPortraits(dir, count); err == nil || pack != nil {
				t.Fatalf("invalid group accepted %v %v", pack, err)
			}
		})
	}
}

// All catalogs must work with no PNG present. This detects eager image I/O,
// including hash/header scans that do not decode pixels or populate the GPU.
func TestHDCatalogsDeferImageReads(t *testing.T) {
	for _, kind := range []string{"portrait", "world", "location", "banner"} {
		t.Run(kind, func(t *testing.T) {
			dir, imageFile, load := hdLazyFixture(t, kind)
			path := filepath.Join(dir, imageFile)
			if err := os.Rename(path, path+".saved"); err != nil {
				t.Fatal(err)
			}
			decode, err := load()
			if err != nil {
				t.Fatalf("catalog opened an image before display: %v", err)
			}
			if err := decode(); err == nil {
				t.Fatal("missing requested image must fall back")
			}
			if err := os.Rename(path+".saved", path); err != nil {
				t.Fatal(err)
			}
			if err := decode(); err != nil {
				t.Fatalf("restored image must decode on demand: %v", err)
			}
			if err := os.WriteFile(path, []byte("changed"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := decode(); err == nil {
				t.Fatal("changed image must be rejected when requested")
			}
		})
	}
}

func TestHDCatalogsRejectImageSymlinkOnDemand(t *testing.T) {
	for _, kind := range []string{"portrait", "world", "location", "banner"} {
		t.Run(kind, func(t *testing.T) {
			dir, file, load := hdLazyFixture(t, kind)
			path := filepath.Join(dir, file)
			outside := filepath.Join(t.TempDir(), "outside.png")
			if err := os.Rename(path, outside); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, path); err != nil {
				t.Fatal(err)
			}
			decode, err := load()
			if err != nil {
				t.Fatalf("catalog should only read metadata: %v", err)
			}
			if err := decode(); err == nil {
				t.Fatal("image symlink escaped its asset directory")
			}
		})
	}
}

func hdLazyFixture(t *testing.T, kind string) (string, string, func() (func() error, error)) {
	t.Helper()
	if kind == "portrait" {
		dir, _ := hdFixture(t)
		return dir, "native.png", func() (func() error, error) {
			pack, err := LoadHDPortraits(dir, 150)
			return func() error { _, err := pack.Image(0); return err }, err
		}
	}
	if kind == "location" {
		dir, original, sample, _ := locationFixture(t)
		hdWrite(t, dir, "originals.json", original)
		hdWrite(t, dir, "manifest.json", map[string]any{"samples": []any{sample}})
		return dir, "hd-generated/image.png", func() (func() error, error) {
			pack, err := LoadHDLocations(dir, 15)
			return func() error { _, err := pack.Image(0); return err }, err
		}
	}
	dir, rec, cat := hdWorldFixture(t)
	if kind == "world" {
		hdWrite(t, dir, "record.json", rec)
		hdWrite(t, dir, "regions.json", cat)
		return dir, "image.png", func() (func() error, error) {
			pack, err := LoadHDWorld(dir)
			return func() error { _, err := pack.Image(0); return err }, err
		}
	}
	rec["content_rect"], rec["palette_bank"] = [4]int{0, 0, 8, 8}, 0
	hdWrite(t, dir, "banner-selection.json", rec)
	return dir, "image.png", func() (func() error, error) {
		pack, err := LoadHDBanner(dir)
		return func() error { _, err := pack.Image(); return err }, err
	}
}

func TestHDPortraitUnavailablePageDoesNotDisableOtherPages(t *testing.T) {
	dir, entry := hdFixture(t)
	entry["portrait_page"], entry["file"] = 1, "missing.png"
	hdWrite(t, dir, "missing.json", entry)
	hdWrite(t, dir, "selection.json", map[string]any{"records": []string{"selected.json", "missing.json"}})
	pack, err := LoadHDPortraits(dir, 150)
	if err != nil || pack.Count() != 2 {
		t.Fatal("unseen missing PNG disabled the whole group", err)
	}
	if _, err := pack.Image(0); err != nil {
		t.Fatal("valid page must remain usable", err)
	}
	if _, err := pack.Image(1); err == nil {
		t.Fatal("only the missing page must fall back")
	}
}

func TestHDPortraitChecksHashAndActualDimensionsAtDisplay(t *testing.T) {
	for _, kind := range []string{"hash", "dimensions"} {
		t.Run(kind, func(t *testing.T) {
			dir, entry := hdFixture(t)
			if kind == "hash" {
				entry["sha256"] = hex.EncodeToString(make([]byte, 32))
			} else {
				entry["width"] = 127
			}
			hdWrite(t, dir, "selected.json", entry)
			pack, err := LoadHDPortraits(dir, 150)
			if err != nil {
				t.Fatal("structurally valid metadata must not open PNG", err)
			}
			if _, err := pack.Image(0); err == nil {
				t.Fatal("requested image did not validate actual bytes", kind)
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
