package cjk

import (
	"os"
	"path/filepath"
	"testing"
)

// 合成一份最小區位字型：填滿指定 (區,位) 的格子。
// 真字型是別人的資產不進版控，測試以格式面驗證索引公式。
func writeFakeKuTen(t *testing.T, ku, ten int, fill bool) string {
	t.Helper()
	idx := (ku-1)*kutenZoneSize + (ten - 1)
	data := make([]byte, (idx+1)*kutenGlyphStride)
	if fill {
		for y := 0; y < kutenGlyphH; y++ {
			data[idx*kutenGlyphStride+y*2] = 0xFF
		}
	}
	p := filepath.Join(t.TempDir(), "font")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestKuTenIndexFormulaGB2312(t *testing.T) {
	// 「啊」＝ B0A1 ⇒ 區 16、位 1。
	f, err := LoadKuTen16(writeFakeKuTen(t, 16, 1, true), GB2312, Options{})
	if err != nil {
		t.Fatal(err)
	}
	a, ok := f.Glyph('啊')
	if !ok || countAlpha(a) != kutenGlyphH*8 {
		t.Fatalf("「啊」的字模不對：ok=%v 亮點=%d", ok, countAlpha(a))
	}
	if b := a.Bounds(); b.Dy() != 16 || b.Dx() != 16 {
		t.Fatalf("字模尺寸 %v，預期 16×16", b)
	}
	// 繁體字不在 GB2312 → 缺字（呈現層換下一份字型或畫方框）。
	if _, ok := f.Glyph('龜'); ok {
		t.Fatal("GB2312 範圍外的字不該有字模")
	}
	// 超出檔長 → 缺字，不是越界。
	if _, ok := f.Glyph('魑'); ok {
		t.Fatal("超出檔長的字不該有字模")
	}
}

func TestKuTenIndexFormulaJIS(t *testing.T) {
	// 「あ」＝ EUC-JP A4A2 ⇒ 區 4、位 2。
	f, err := LoadKuTen16(writeFakeKuTen(t, 4, 2, true), JISX0208, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if a, ok := f.Glyph('あ'); !ok || countAlpha(a) == 0 {
		t.Fatalf("「あ」取不到字模：ok=%v", ok)
	}
	// ⭐ JIS X 0212 的補助漢字（EUC-JP 三位元組、`0x8F` 前綴）**不在**
	// JISKAN16 裡。不擋掉的話會拿 `0x8F` 當區碼算出不相干的字模——
	// 畫面上會是「有字但都不對」，比缺字難查得多。
	// 「汜」正是這一類（PC-98 版把它放外字，見 tools/namepack.py）。
	if _, ok := f.Glyph('汜'); ok {
		t.Fatal("JIS X 0212 的字不該從 JISKAN16 取到字模")
	}
}

// 自我檢查要能擋「檔案對不上索引」：全零檔的探針字是空的。
func TestKuTenRejectsBlankProbe(t *testing.T) {
	if _, err := LoadKuTen16(writeFakeKuTen(t, 16, 1, false), GB2312, Options{}); err == nil {
		t.Fatal("全零檔應該被自我檢查擋下")
	}
	if _, err := LoadKuTen16(writeFakeKuTen(t, 4, 2, false), JISX0208, Options{}); err == nil {
		t.Fatal("全零檔應該被自我檢查擋下（JIS）")
	}
}
