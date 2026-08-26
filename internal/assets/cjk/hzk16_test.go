package cjk

import (
	"os"
	"path/filepath"
	"testing"
)

// 合成一份最小 HZK16：只填「啊」（B0A1，idx 1410）的格子。
// 真字型是商業軟體不進版控，測試以格式面驗證索引公式。
func writeFakeHZK16(t *testing.T, fill bool) string {
	t.Helper()
	const idxA = (0xB0-0xA1)*94 + (0xA1 - 0xA1) // 1410
	data := make([]byte, (idxA+1)*hzkGlyphStride)
	if fill {
		for y := 0; y < hzkGlyphH; y++ {
			data[idxA*hzkGlyphStride+y*2] = 0xFF
		}
	}
	p := filepath.Join(t.TempDir(), "HZK16")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestHZK16IndexFormula(t *testing.T) {
	f, err := LoadHZK16(writeFakeHZK16(t, true))
	if err != nil {
		t.Fatal(err)
	}
	a, ok := f.Glyph('啊')
	if !ok || countAlpha(a) != hzkGlyphH*8 {
		t.Fatalf("「啊」的字模不對：ok=%v 亮點=%d", ok, countAlpha(a))
	}
	if b := a.Bounds(); b.Dy() != 16 || b.Dx() != 16 {
		t.Fatalf("字模尺寸 %v，預期 16×16", b)
	}
	// 繁體字不在 GB2312 → 缺字（呈現層畫方框）。
	if _, ok := f.Glyph('龜'); ok {
		t.Fatal("GB2312 範圍外的字不該有字模")
	}
	// 檔案範圍外的字（次常用區之後）→ 缺字而不是越界。
	if _, ok := f.Glyph('魑'); ok {
		t.Fatal("超出檔長的字不該有字模")
	}
}

// 自我檢查要能擋「檔案對不上索引」：全零檔的「啊」是空的。
func TestHZK16RejectsBlankProbe(t *testing.T) {
	if _, err := LoadHZK16(writeFakeHZK16(t, false)); err == nil {
		t.Fatal("全零檔應該被自我檢查擋下")
	}
}
