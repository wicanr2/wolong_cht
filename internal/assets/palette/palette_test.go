package palette

import (
	"os"
	"path/filepath"
	"testing"
)

func load(t *testing.T, name string) *Palette {
	t.Helper()
	p := filepath.Join("..", "..", "..", "workplace", "orig", "dosv", name)
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("找不到原版素材 %s，跳過", p)
	}
	pal, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return pal
}

// TestBankCounts 釘住 docs/formats/02 §1 的組數。
func TestBankCounts(t *testing.T) {
	for name, want := range map[string]int{
		"GAMEPAL.BRG": 8, "OPENPAL.BRG": 6, "ENDPAL.BRG": 12, "OVERPAL.BRG": 1,
	} {
		if got := len(load(t, name).Banks); got != want {
			t.Errorf("%s 有 %d 組，預期 %d 組", name, got, want)
		}
	}
}

// TestSeasonColor 釘住四季只換色號 14 這件事（docs/formats/02 §4）。
// 這條規則是 remake 的實作約束，漂掉了季節效果就不對。
func TestSeasonColor(t *testing.T) {
	pal := load(t, "GAMEPAL.BRG")
	want := []struct {
		season     Season
		r, g, b    uint8
	}{
		{Spring, 0x88, 0xaa, 0x66},
		{Summer, 0x55, 0xaa, 0x11},
		{Autumn, 0xdd, 0x88, 0x00},
		{Winter, 0xff, 0xff, 0xff},
	}
	for _, w := range want {
		c := pal.Banks[w.season][14]
		if c.R != w.r || c.G != w.g || c.B != w.b {
			t.Errorf("季 %d 的色 14 是 #%02x%02x%02x，預期 #%02x%02x%02x",
				w.season, c.R, c.G, c.B, w.r, w.g, w.b)
		}
	}
	// 除了色 14，四季應該完全一樣。
	for s := Summer; s <= Winter; s++ {
		for i := 0; i < BankSize; i++ {
			if i == 14 {
				continue
			}
			if pal.Banks[s][i] != pal.Banks[Spring][i] {
				t.Errorf("季 %d 的色 %d 與春不同 —— 季節應該只換色 14", s, i)
			}
		}
	}
}

// TestScale 釘住亮度公式：滿值 16 時原值不變，0 時全黑。
func TestScale(t *testing.T) {
	for v := byte(0); v <= 15; v++ {
		if got := Scale(v, FullBrightness); got != v {
			t.Errorf("Scale(%d, 滿值) = %d，預期 %d", v, got, v)
		}
		if got := Scale(v, 0); got != 0 {
			t.Errorf("Scale(%d, 0) = %d，預期 0", v, got)
		}
	}
}

// TestWhiteIsPureWhite 防止用 v<<4 —— 那會讓白色變成 #f0f0f0。
func TestWhiteIsPureWhite(t *testing.T) {
	c := load(t, "GAMEPAL.BRG").Banks[0][15]
	if c.R != 0xff || c.G != 0xff || c.B != 0xff {
		t.Errorf("色 15 是 #%02x%02x%02x，預期純白", c.R, c.G, c.B)
	}
}
