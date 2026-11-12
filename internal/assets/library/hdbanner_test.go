package library

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHDBannerValidationAndDecodeRecheck(t *testing.T) {
	for _, name := range []string{"valid", "outside", "negative", "bad palette", "not native", "path escape", "bad SHA"} {
		t.Run(name, func(t *testing.T) {
			dir, rec, _ := hdWorldFixture(t)
			rec["content_rect"] = [4]int{0, 1, 8, 4}
			rec["palette_bank"] = 0
			switch name {
			case "outside":
				rec["content_rect"] = [4]int{0, 4, 8, 8}
			case "negative":
				rec["content_rect"] = [4]int{-1, 0, 8, 4}
			case "bad palette":
				rec["palette_bank"] = 4
			case "not native":
				rec["native_unmodified"] = false
			case "path escape":
				rec["file"] = "../image.png"
			case "bad SHA":
				rec["sha256"] = "wrong"
			}
			writeHDWorldJSON(t, dir, "banner-selection.json", rec)
			pack, err := LoadHDBanner(dir)
			if name != "valid" {
				if err == nil {
					t.Fatal("invalid banner accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err = pack.Image(); err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(filepath.Join(dir, "image.png"), []byte("corrupt"), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err = pack.Image(); err == nil {
				t.Fatal("banner decode did not recheck source")
			}
		})
	}
}

func TestHDBannerActualNativeSelection(t *testing.T) {
	dir := os.Getenv("WOLONG_HD_BANNER")
	if dir == "" {
		t.Skip("local selected banner not requested")
	}
	pack, err := LoadHDBanner(dir)
	if err != nil {
		t.Fatal(err)
	}
	if pack.Bank != 0 || pack.Content.Dx() != 2172 || pack.Content.Dy() != 106 {
		t.Fatal("selected native content changed", pack.Content)
	}
	if _, err = pack.Image(); err != nil {
		t.Fatal(err)
	}
}
