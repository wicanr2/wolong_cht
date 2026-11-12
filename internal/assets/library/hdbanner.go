package library

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"path/filepath"
)

// HDBanner adapts the local selected native banner. Spec275 fourth slice.
type HDBanner struct {
	root          string
	file, digest  string
	width, height int
	Content       image.Rectangle
	Bank          int
}

func LoadHDBanner(dir string) (*HDBanner, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	b, err := readHDAsset(root, "banner-selection.json", 1<<20)
	if err != nil {
		return nil, err
	}
	var rec struct {
		File    string `json:"file"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		SHA     string `json:"sha256"`
		Native  bool   `json:"native_unmodified"`
		Content [4]int `json:"content_rect"`
		Bank    int    `json:"palette_bank"`
	}
	if err = json.Unmarshal(b, &rec); err != nil {
		return nil, err
	}
	c := rec.Content
	if !rec.Native || rec.Bank < 0 || rec.Bank >= 4 || c[0] < 0 || c[1] < 0 || c[2] <= 0 || c[3] <= 0 || c[0] > rec.Width || c[1] > rec.Height || c[2] > rec.Width-c[0] || c[3] > rec.Height-c[1] {
		return nil, fmt.Errorf("invalid HD banner content or palette")
	}
	if _, err = readHDPNG(root, rec.File, rec.Width, rec.Height, rec.SHA); err != nil {
		return nil, err
	}
	return &HDBanner{root: root, file: rec.File, digest: rec.SHA, width: rec.Width, height: rec.Height, Content: image.Rect(c[0], c[1], c[0]+c[2], c[1]+c[3]), Bank: rec.Bank}, nil
}

func (s *HDBanner) Image() (image.Image, error) {
	if s == nil {
		return nil, fmt.Errorf("HD banner unavailable")
	}
	b, err := readHDPNG(s.root, s.file, s.width, s.height, s.digest)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(b))
}
