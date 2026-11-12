package library

// Existing local art-study catalog adapter; not a new distribution format. Spec275.
import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"path/filepath"
)

type HDWorldLayer struct {
	ID            string
	Source, Core  image.Rectangle
	Width, Height int
	Feather       int
	file, digest  string
}

type HDWorld struct {
	root   string
	layers []HDWorldLayer
}

type hdWorldRegion struct {
	ID     string `json:"id"`
	Record string `json:"record"`
	Source [4]int `json:"source_world_pixel_rect"`
	Core   [4]int `json:"core_world_pixel_rect"`
}

func hdWorldRect(r [4]int) (image.Rectangle, error) {
	if r[0] < 0 || r[1] < 0 || r[0] > 6144 || r[1] > 4096 || r[2] <= 0 || r[3] <= 0 || r[2] > 6144-r[0] || r[3] > 4096-r[1] {
		return image.Rectangle{}, fmt.Errorf("invalid HD world rectangle: %v", r)
	}
	return image.Rect(r[0], r[1], r[0]+r[2], r[1]+r[3]), nil
}

func LoadHDWorld(dir string) (*HDWorld, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	b, err := readHDAsset(root, "regions.json", 4<<20)
	if err != nil {
		return nil, err
	}
	var catalog struct {
		Regions []hdWorldRegion `json:"regions"`
		Seams   []hdWorldRegion `json:"seams"`
	}
	if err = json.Unmarshal(b, &catalog); err != nil {
		return nil, err
	}
	if len(catalog.Regions) == 0 || len(catalog.Regions)+len(catalog.Seams) > 256 {
		return nil, fmt.Errorf("invalid HD world layer count")
	}
	s := &HDWorld{root: root}
	seen := map[string]bool{}
	for group, rows := range [][]hdWorldRegion{catalog.Regions, catalog.Seams} {
		for _, row := range rows {
			key := fmt.Sprintf("%d:%s", group, row.ID)
			if row.ID == "" || seen[key] {
				return nil, fmt.Errorf("duplicate or empty HD world ID: %s", key)
			}
			seen[key] = true
			if filepath.Ext(row.Record) != ".json" {
				return nil, fmt.Errorf("invalid HD world record")
			}
			b, err = readHDAsset(root, row.Record, 4<<20)
			if err != nil {
				return nil, err
			}
			var rec struct {
				File    string `json:"file"`
				Width   int    `json:"width"`
				Height  int    `json:"height"`
				SHA     string `json:"sha256"`
				Native  bool   `json:"native_unmodified"`
				Source  [4]int `json:"source_world_pixel_rect"`
				Core    [4]int `json:"core_world_pixel_rect"`
				Feather int    `json:"feather_world_pixels"`
			}
			if err = json.Unmarshal(b, &rec); err != nil {
				return nil, err
			}
			if group == 1 {
				row.Source, row.Core = rec.Source, rec.Core
			}
			source, err := hdWorldRect(row.Source)
			if err != nil {
				return nil, err
			}
			core, err := hdWorldRect(row.Core)
			if err != nil {
				return nil, err
			}
			if !core.In(source) || !rec.Native {
				return nil, fmt.Errorf("invalid HD world core or native source: %s", key)
			}
			feather := 0
			if group == 1 {
				feather = rec.Feather
				if feather <= 0 || feather >= (min(core.Dx(), core.Dy())+1)/2 {
					return nil, fmt.Errorf("invalid HD seam feather: %s", key)
				}
			}
			if group == 0 {
				for _, prev := range s.layers {
					if prev.Core.Overlaps(core) {
						return nil, fmt.Errorf("overlapping HD world regions")
					}
				}
			}
			digest, err := hex.DecodeString(rec.SHA)
			if err != nil || len(digest) != 32 || filepath.Ext(rec.File) != ".png" {
				return nil, fmt.Errorf("invalid HD world PNG record: %s", key)
			}
			if _, err = readHDPNG(root, rec.File, rec.Width, rec.Height, rec.SHA); err != nil {
				return nil, err
			}
			s.layers = append(s.layers, HDWorldLayer{ID: key, Source: source, Core: core, Width: rec.Width, Height: rec.Height, Feather: feather, file: rec.File, digest: rec.SHA})
		}
	}
	return s, nil
}

func (s *HDWorld) Count() int {
	if s == nil {
		return 0
	}
	return len(s.layers)
}
func (s *HDWorld) Layers() []HDWorldLayer {
	if s == nil {
		return nil
	}
	return append([]HDWorldLayer(nil), s.layers...)
}
func (s *HDWorld) Image(index int) (image.Image, error) {
	if s == nil || index < 0 || index >= len(s.layers) {
		return nil, fmt.Errorf("HD world layer unavailable: %d", index)
	}
	e := s.layers[index]
	b, err := readHDPNG(s.root, e.file, e.Width, e.Height, e.digest)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(b))
}
