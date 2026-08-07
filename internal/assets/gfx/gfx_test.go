package gfx

import (
	"os"
	"path/filepath"
	"testing"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	p := filepath.Join("..", "..", "..", "workplace", "orig", "dosv", name)
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("找不到原版素材 %s，跳過", p)
	}
	return raw
}

// TestNoRemainder 是圖庫尺寸最便宜的檢核：餘數不是 0 就代表尺寸錯了。
// 四個圖庫的尺寸全部出自反組譯（docs/re/03），四個的餘數都該是 0。
func TestNoRemainder(t *testing.T) {
	for _, c := range []struct {
		spec  Spec
		file  string
		count int
	}{
		{Kao, "KAOGRF.DAT", 150},
		{Kyo, "KYOGRF.DAT", 15},
		{Ivent, "IVENTGRF.DAT", 3},
	} {
		got, rem := c.spec.Count(read(t, c.file))
		if rem != 0 {
			t.Errorf("%s 餘 %d byte —— 尺寸 %dx%d 是錯的",
				c.file, rem, c.spec.Width, c.spec.Height)
		}
		if got != c.count {
			t.Errorf("%s 解出 %d 張，預期 %d 張", c.file, got, c.count)
		}
	}
}

// TestIconRegionsCoverFile 釘住 ICONGRF 四段的長度加起來等於檔案大小。
func TestIconRegionsCoverFile(t *testing.T) {
	raw := read(t, "ICONGRF.DAT")
	total := 0
	for _, r := range IconRegions {
		if r.Offset != total {
			t.Fatalf("段 %s 的位移 0x%X 與前一段的結尾 0x%X 對不上",
				r.Name, r.Offset, total)
		}
		total += r.Length
	}
	if total != len(raw) {
		t.Errorf("四段合計 %d byte，檔案 %d byte", total, len(raw))
	}
}

// TestDecodeBounds 確認越界會回 error，不是 panic 或靜悄悄回錯的圖。
func TestDecodeBounds(t *testing.T) {
	raw := read(t, "KAOGRF.DAT")
	if _, err := Kao.Decode(raw, 150); err == nil {
		t.Error("第 150 張（超出範圍）應該回 error")
	}
	if _, err := Kao.Decode(raw, -1); err == nil {
		t.Error("負數索引應該回 error")
	}
}
