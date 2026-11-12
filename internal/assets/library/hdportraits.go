package library

// HD頭像沿既有本機選稿表讀取；不將美術嵌入引擎。契約見spec275第一段接線。
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

type hdPortraitEntry struct {
	PortraitPage *int   `json:"portrait_page"`
	File         string `json:"file"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	SHA256       string `json:"sha256"`
}

// HDPortraits holds a validated page index. Decoding stays outside Ebiten.
type HDPortraits struct {
	root  string
	pages map[int]hdPortraitEntry
}

func (s *HDPortraits) read(name string, limit int64) ([]byte, error) {
	return readHDAsset(s.root, name, limit)
}

func readHDAsset(root, name string, limit int64) ([]byte, error) {
	if name == "" || filepath.Base(name) != name || name == "." || name == ".." {
		return nil, fmt.Errorf("invalid HD filename %q", name)
	}
	path, err := filepath.EvalSymlinks(filepath.Join(root, name))
	if err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(root, path)
	if err != nil || rel != name {
		return nil, fmt.Errorf("HD path escapes asset directory: %q", name)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("HD input is not regular: %q", name)
	}
	b, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("HD input exceeds limit: %q", name)
	}
	return b, nil
}

func (s *HDPortraits) pngBytes(e hdPortraitEntry) ([]byte, error) {
	return readHDPNG(s.root, e.File, e.Width, e.Height, e.SHA256)
}

func readHDPNG(root, file string, width, height int, digest string) ([]byte, error) {
	b, err := readHDAsset(root, file, 16<<20)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(b)
	if hex.EncodeToString(h[:]) != digest {
		return nil, fmt.Errorf("HD SHA mismatch: %s", file)
	}
	c, err := png.DecodeConfig(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	if c.Width != width || c.Height != height || c.Width <= 0 || c.Height <= 0 || c.Width > 4096 || c.Height > 4096 || c.Width*c.Height > 4194304 {
		return nil, fmt.Errorf("invalid HD dimensions: %s", file)
	}
	return b, nil
}

// LoadHDPortraits validates the whole selected group before it becomes available.
// pageCount must come from the loaded original archive, not a version string.
func LoadHDPortraits(dir string, pageCount int) (*HDPortraits, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	if pageCount <= 0 {
		return nil, fmt.Errorf("invalid original portrait count")
	}
	s := &HDPortraits{root: root, pages: map[int]hdPortraitEntry{}}
	b, err := s.read("selection.json", 1<<20)
	if err != nil {
		return nil, err
	}
	var selection struct {
		Records []string `json:"records"`
	}
	if err = json.Unmarshal(b, &selection); err != nil {
		return nil, err
	}
	if len(selection.Records) == 0 || len(selection.Records) > pageCount {
		return nil, fmt.Errorf("invalid HD selection count")
	}
	for _, name := range selection.Records {
		if filepath.Ext(name) != ".json" {
			return nil, fmt.Errorf("invalid HD record %q", name)
		}
		b, err = s.read(name, 1<<20)
		if err != nil {
			return nil, err
		}
		var e hdPortraitEntry
		if err = json.Unmarshal(b, &e); err != nil {
			return nil, err
		}
		if e.PortraitPage == nil || *e.PortraitPage < 0 || *e.PortraitPage >= pageCount {
			return nil, fmt.Errorf("HD page out of original archive bounds: %s", name)
		}
		if _, ok := s.pages[*e.PortraitPage]; ok {
			return nil, fmt.Errorf("duplicate HD page: %d", *e.PortraitPage)
		}
		digest, err := hex.DecodeString(e.SHA256)
		if err != nil || len(digest) != sha256.Size || filepath.Ext(e.File) != ".png" {
			return nil, fmt.Errorf("invalid HD image record: %s", name)
		}
		if _, err = s.pngBytes(e); err != nil {
			return nil, err
		}
		s.pages[*e.PortraitPage] = e
	}
	return s, nil
}

func (s *HDPortraits) Count() int {
	if s == nil {
		return 0
	}
	return len(s.pages)
}

// Image rechecks the selected native bytes on each decode. Missing pages return
// an error so the caller keeps the original portrait.
func (s *HDPortraits) Image(page int) (image.Image, error) {
	if s == nil {
		return nil, fmt.Errorf("HD portraits unavailable")
	}
	e, ok := s.pages[page]
	if !ok {
		return nil, fmt.Errorf("no HD portrait for page %d", page)
	}
	b, err := s.pngBytes(e)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(b))
}
