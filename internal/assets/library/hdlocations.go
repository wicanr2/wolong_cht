package library

// Adapter for the existing approved five-city scene manifest. Spec275 slice 5.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"path/filepath"
)

type hdLocationEntry struct {
	File   string `json:"file"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	SHA    string `json:"sha256"`
	Native bool   `json:"native_unmodified"`
}

type HDLocations struct {
	root  string
	pages map[int]hdLocationEntry
	Bank  int
}

type hdLocationCity struct {
	ID        int  `json:"id"`
	ScenePage *int `json:"scene_page"`
}

func LoadHDLocations(dir string, pageCount int) (*HDLocations, error) {
	if pageCount <= 0 {
		return nil, fmt.Errorf("invalid original location count")
	}
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	b, err := readHDAsset(root, "originals.json", 1<<20)
	if err != nil {
		return nil, err
	}
	var original struct {
		Count  int              `json:"actual_scene_count"`
		Bank   *int             `json:"palette_bank"`
		Cities []hdLocationCity `json:"cities"`
	}
	if err = json.Unmarshal(b, &original); err != nil {
		return nil, err
	}
	if original.Count != pageCount || original.Bank == nil || *original.Bank < 0 || *original.Bank >= 4 {
		return nil, fmt.Errorf("HD location source archive or palette mismatch")
	}
	mapping := map[int]int{}
	for _, c := range original.Cities {
		if c.ScenePage == nil || *c.ScenePage < 0 || *c.ScenePage >= pageCount {
			return nil, fmt.Errorf("invalid original location mapping")
		}
		if _, exists := mapping[c.ID]; exists {
			return nil, fmt.Errorf("duplicate source city")
		}
		mapping[c.ID] = *c.ScenePage
	}
	b, err = readHDAsset(root, "manifest.json", 1<<20)
	if err != nil {
		return nil, err
	}
	var manifest struct {
		Samples []struct {
			Kind   string          `json:"kind"`
			CityID int             `json:"city_id"`
			City   hdLocationCity  `json:"city"`
			HD     hdLocationEntry `json:"hd"`
			Record hdLocationEntry `json:"hd_record"`
		} `json:"samples"`
	}
	if err = json.Unmarshal(b, &manifest); err != nil {
		return nil, err
	}
	nativeRoot, err := filepath.EvalSymlinks(filepath.Join(root, "hd-generated"))
	if err != nil {
		return nil, err
	}
	if rel, e := filepath.Rel(root, nativeRoot); e != nil || rel != "hd-generated" {
		return nil, fmt.Errorf("HD location directory escapes source root")
	}
	s := &HDLocations{root: nativeRoot, pages: map[int]hdLocationEntry{}, Bank: *original.Bank}
	for _, sample := range manifest.Samples {
		if sample.Kind != "scene" {
			continue
		}
		page, ok := mapping[sample.CityID]
		if !ok || sample.City.ID != sample.CityID || sample.City.ScenePage == nil || *sample.City.ScenePage != page {
			return nil, fmt.Errorf("HD location sample/source mismatch")
		}
		if _, ok = s.pages[page]; ok {
			return nil, fmt.Errorf("duplicate HD location page: %d", page)
		}
		e := sample.Record
		if !e.Native || sample.HD.File != "hd-generated/"+e.File || sample.HD.Width != e.Width || sample.HD.Height != e.Height || sample.HD.SHA != e.SHA {
			return nil, fmt.Errorf("HD location selection mismatch")
		}
		if _, err = readHDPNG(nativeRoot, e.File, e.Width, e.Height, e.SHA); err != nil {
			return nil, err
		}
		s.pages[page] = e
	}
	if len(s.pages) == 0 || len(s.pages) > pageCount {
		return nil, fmt.Errorf("invalid HD location selection count")
	}
	return s, nil
}

func (s *HDLocations) Count() int {
	if s == nil {
		return 0
	}
	return len(s.pages)
}
func (s *HDLocations) Image(page int) (image.Image, error) {
	if s == nil {
		return nil, fmt.Errorf("HD locations unavailable")
	}
	e, ok := s.pages[page]
	if !ok {
		return nil, fmt.Errorf("no HD location for page %d", page)
	}
	b, err := readHDPNG(s.root, e.File, e.Width, e.Height, e.SHA)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(b))
}
