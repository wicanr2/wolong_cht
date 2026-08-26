// HZK16：UCDOS 的 GB2312 16×16 點陣字型，簡體語系用（docs/spec/84）。
//
// 與倚天同一個政策：**字型檔不隨本專案散布**，由使用者自備放進字型目錄。
// 格式：區位碼索引，glyph = ((區−0xA1)×94 ＋ (位−0xA1)) × 32，
// 每字 32 byte、16 列 × 每列 2 byte、MSB-first。
//
// ⚠ 字高 16 比倚天的 15 多一列。版面常數（textdraw 的 GlyphH）不動，
// 第 16 列畫進行距裡——視覺上簡體字比繁體字低一像素收尾，
// 不影響排版計算（docs/spec/84 §2）。
package cjk

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	hzkGlyphH      = 16
	hzkRowBytes    = 2
	hzkGlyphStride = hzkGlyphH * hzkRowBytes // 32
	hzkZoneSize    = 94
)

// HZK16Font 是載入好的 GB2312 點陣字。與 *Font 一樣以 rune 取字模，
// textdraw 透過 GlyphSource 介面使用，不認識底層是哪一種字型。
type HZK16Font struct {
	data []byte
	bold bool
}

// LoadHZK16 讀取 HZK16 檔。
//
// 自我檢查：GB2312 漢字區第一個字「啊」（B0A1）的字模不得全空——
// 空的代表整體偏移或檔案不對，會呈現成「有字但都不對」。
func LoadHZK16(path string) (*HZK16Font, error) {
	return LoadHZK16WithOptions(path, Options{})
}

// LoadHZK16WithOptions 載入並套用顯示選項（Bold 與倚天同一種加粗）。
func LoadHZK16WithOptions(path string, opts Options) (*HZK16Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cjk: 讀取 HZK16 %s 失敗: %w", path, err)
	}
	if len(data) == 0 || len(data)%hzkGlyphStride != 0 {
		return nil, fmt.Errorf("cjk: %s 長度 %d 不是 %d 的整數倍，格式不符",
			path, len(data), hzkGlyphStride)
	}
	f := &HZK16Font{data: data, bold: opts.Bold}
	if a, ok := f.Glyph('啊'); !ok || countAlpha(a) == 0 {
		return nil, fmt.Errorf("cjk: %s 的「啊」(B0A1) 取不到字模，索引公式或檔案版本不符", path)
	}
	return f, nil
}

// LoadHZK16Dir 從字型目錄找 HZK16（大小寫兩種拼法都接受）。
func LoadHZK16Dir(dir string, opts Options) (*HZK16Font, error) {
	path, err := fontPath(dir, "HZK16", "hzk16")
	if err != nil {
		return nil, err
	}
	return LoadHZK16WithOptions(path, opts)
}

// Glyph 取一個字的 16×16 字模。非 GB2312 範圍（含簡體檔沒有的繁體字）
// 回 false，由呈現層畫缺字框——**缺字要看得出來**。
func (f *HZK16Font) Glyph(ch rune) (*image.Alpha, bool) {
	if f == nil {
		return nil, false
	}
	enc, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(string(ch)))
	if err != nil || len(enc) != 2 {
		return nil, false
	}
	qu, wei := enc[0], enc[1]
	// 只收 GB2312 的區位範圍；GBK 的擴充區（0x40–0xA0 第二 byte）不在 HZK16 裡。
	if qu < 0xA1 || qu > 0xF7 || wei < 0xA1 || wei > 0xFE {
		return nil, false
	}
	idx := (int(qu)-0xA1)*hzkZoneSize + (int(wei) - 0xA1)
	off := idx * hzkGlyphStride
	if off < 0 || off+hzkGlyphStride > len(f.data) {
		return nil, false
	}
	img := image.NewAlpha(image.Rect(0, 0, GlyphWidth, hzkGlyphH))
	for y := 0; y < hzkGlyphH; y++ {
		w := uint16(f.data[off+y*2])<<8 | uint16(f.data[off+y*2+1])
		if f.bold {
			w = emboldenRow(w)
		}
		for x := 0; x < GlyphWidth; x++ {
			if w>>(15-x)&1 == 1 {
				img.SetAlpha(x, y, color.Alpha{A: 0xff})
			}
		}
	}
	return img, true
}

func countAlpha(a *image.Alpha) int {
	n := 0
	b := a.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if a.AlphaAt(x, y).A != 0 {
				n++
			}
		}
	}
	return n
}
