package cjk

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
)

// 倚天的半形字型：256 個字 × 15 列 × 1 byte，MSB-first。
// 與 16×15 的全形字**同高**，混排才不會參差。
const (
	ASCIIWidth  = 8
	ASCIIHeight = GlyphHeight // 15
	asciiStride = ASCIIHeight // 每字 15 byte
	asciiCount  = 256
	asciiSize   = asciiCount * asciiStride // 3840
)

// ASCIIFont 是倚天的 ASCFONT.15。
type ASCIIFont struct{ data []byte }

// LoadASCII 讀取 ASCFONT.15。
func LoadASCII(path string) (*ASCIIFont, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// 大小固定。不吻合就直接拒絕——**寧可載不起來，也不要畫出一堆亂碼**，
	// 亂碼比錯誤訊息難查得多。
	if len(b) != asciiSize {
		return nil, fmt.Errorf("%s 大小 %d，預期 %d（256 字 × 15 列）",
			filepath.Base(path), len(b), asciiSize)
	}
	return &ASCIIFont{data: b}, nil
}

// LoadASCIIDir 在字型目錄裡找 ASCFONT.15（大小寫都試）。
func LoadASCIIDir(dir string) (*ASCIIFont, error) {
	for _, n := range []string{"ASCFONT.15", "ascfont.15"} {
		p := filepath.Join(dir, n)
		if _, err := os.Stat(p); err == nil {
			return LoadASCII(p)
		}
	}
	return nil, fmt.Errorf("%s 裡找不到 ASCFONT.15", dir)
}

// Glyph 取一個半形字的字模。
func (f *ASCIIFont) Glyph(ch rune) (*image.Alpha, bool) {
	if f == nil || ch < 0 || ch >= asciiCount {
		return nil, false
	}
	src := f.data[int(ch)*asciiStride:]
	img := image.NewAlpha(image.Rect(0, 0, ASCIIWidth, ASCIIHeight))
	ink := false
	for y := 0; y < ASCIIHeight; y++ {
		row := src[y]
		for x := 0; x < ASCIIWidth; x++ {
			if row&(0x80>>x) != 0 {
				img.SetAlpha(x, y, color.Alpha{A: 0xff})
				ink = true
			}
		}
	}
	// 空白字（例如 0x20）沒有筆劃是正常的，仍然回 true——
	// 回 false 會讓呼叫端誤以為缺字而畫方框。
	_ = ink
	return img, true
}
