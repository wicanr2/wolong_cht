// 區位碼 16×16 點陣字型：簡體用 `HZK16`（GB2312）、日文用 `JISKAN16`
// （JIS X 0208）。兩者是同一種佈局，只有字集不同（docs/spec/84）。
//
// 與倚天同一個政策：**字型檔不隨本專案散布**，由使用者自備放進字型目錄。
//
// 格式：glyph = (區−1)×94 ＋ (位−1)，每字 32 byte、16 列 × 每列 2 byte、
// MSB-first。GB2312 與 JIS X 0208 都是 94×94 的區位字集，所以索引公式共用；
// 差別只在「怎麼把一個 rune 換成區位」——GBK 與 EUC-JP 的雙位元組**就是**
// 區位加 0xA0，所以直接拿編碼器的輸出減 0xA0。
//
// ⚠ 字高 16 比倚天的 15 多一列。版面常數（textdraw 的 GlyphH）不動，
// 第 16 列畫進行距裡——視覺上比倚天低一像素收尾，不影響排版計算。
package cjk

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	kutenGlyphH      = 16
	kutenRowBytes    = 2
	kutenGlyphStride = kutenGlyphH * kutenRowBytes // 32
	kutenZoneSize    = 94
)

// Charset 是區位字型的字集。
type Charset int

const (
	// GB2312 是簡體中文（`HZK16`）。
	GB2312 Charset = iota
	// JISX0208 是日文（`JISKAN16`）。含假名與新字體漢字。
	JISX0208
)

func (c Charset) fileNames() []string {
	if c == JISX0208 {
		return []string{"JISKAN16", "jiskan16", "JISKAN16.F16", "jiskan16.f16"}
	}
	return []string{"HZK16", "hzk16"}
}

// probe 是自我檢查用的字：索引公式或檔案版本不對時，它會是空的。
func (c Charset) probe() rune {
	if c == JISX0208 {
		return 'あ' // 區 4 位 2
	}
	return '啊' // 區 16 位 1
}

// kuten 把一個 rune 換成 (區, 位)；不在這個字集裡就回 false。
func (c Charset) kuten(ch rune) (int, int, bool) {
	var enc interface{ Bytes([]byte) ([]byte, error) }
	if c == JISX0208 {
		enc = japanese.EUCJP.NewEncoder()
	} else {
		enc = simplifiedchinese.GBK.NewEncoder()
	}
	b, err := enc.Bytes([]byte(string(ch)))
	// EUC-JP 的三位元組形式（`0x8F` 前綴）是 JIS X 0212 補助漢字，
	// **不在 JISKAN16 裡**；GBK 的擴充區（第二 byte < 0xA1）也不在 HZK16 裡。
	// 兩種都要擋掉，否則會算出越界或不相干的字模。
	if err != nil || len(b) != 2 || b[0] < 0xA1 || b[0] > 0xFE || b[1] < 0xA1 || b[1] > 0xFE {
		return 0, 0, false
	}
	return int(b[0]) - 0xA0, int(b[1]) - 0xA0, true
}

// KuTen16Font 是載入好的區位點陣字。以 rune 取字模，
// textdraw 透過 GlyphSource 介面使用，不認識底層是哪一種字型。
type KuTen16Font struct {
	data []byte
	cs   Charset
	bold bool
}

// LoadKuTen16 讀取一份區位字型並做自我檢查。
func LoadKuTen16(path string, cs Charset, opts Options) (*KuTen16Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cjk: 讀取區位字型 %s 失敗: %w", path, err)
	}
	if len(data) == 0 || len(data)%kutenGlyphStride != 0 {
		return nil, fmt.Errorf("cjk: %s 長度 %d 不是 %d 的整數倍，格式不符",
			path, len(data), kutenGlyphStride)
	}
	f := &KuTen16Font{data: data, cs: cs, bold: opts.Bold}
	probe := cs.probe()
	if a, ok := f.Glyph(probe); !ok || countAlpha(a) == 0 {
		return nil, fmt.Errorf("cjk: %s 取不到「%c」的字模，索引公式或檔案版本不符",
			path, probe)
	}
	return f, nil
}

// LoadKuTen16Dir 從字型目錄找對應字集的檔案（大小寫兩種拼法都接受）。
func LoadKuTen16Dir(dir string, cs Charset, opts Options) (*KuTen16Font, error) {
	path, err := fontPath(dir, cs.fileNames()...)
	if err != nil {
		return nil, err
	}
	return LoadKuTen16(path, cs, opts)
}

// Glyph 取一個字的 16×16 字模。不在字集裡或超出檔長就回 false，
// 由呈現層接手（字型鏈的下一份，或畫缺字框）——**缺字要看得出來**。
func (f *KuTen16Font) Glyph(ch rune) (*image.Alpha, bool) {
	if f == nil {
		return nil, false
	}
	ku, ten, ok := f.cs.kuten(ch)
	if !ok {
		return nil, false
	}
	off := ((ku-1)*kutenZoneSize + (ten - 1)) * kutenGlyphStride
	if off < 0 || off+kutenGlyphStride > len(f.data) {
		return nil, false
	}
	img := image.NewAlpha(image.Rect(0, 0, GlyphWidth, kutenGlyphH))
	for y := 0; y < kutenGlyphH; y++ {
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
