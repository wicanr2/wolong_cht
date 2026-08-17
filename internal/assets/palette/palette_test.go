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
		season  Season
		r, g, b uint8
	}{
		// 走 DOS/V 的 VGA DAC（docs/spec/51）：4 bit 的 f 是 0xF3 不是 0xFF。
		{Spring, 0x82, 0xa2, 0x61},
		{Summer, 0x51, 0xa2, 0x10},
		{Autumn, 0xd3, 0x82, 0x00},
		{Winter, 0xf3, 0xf3, 0xf3},
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

// TestDACScaleNeverReachesFullScale 釘住 docs/spec/51：DOS/V 的 4 bit 通道
// 走 `shl ah,1` 兩次進 VGA DAC，最大只到 60/63，**畫面上沒有純白**。
//
// 這條同時擋住兩種寫法：`v*255/15`（15 → 255，整張畫面亮 4%）
// 與 `v<<4`（15 → 240，白色發灰）。兩者都不是原版。
func TestDACScaleNeverReachesFullScale(t *testing.T) {
	c := load(t, "GAMEPAL.BRG").Banks[0][15]
	if c.R != 0xf3 || c.G != 0xf3 || c.B != 0xf3 {
		t.Errorf("色 15 是 #%02x%02x%02x，預期 #f3f3f3", c.R, c.G, c.B)
	}
	for v := byte(0); v <= 15; v++ {
		if got := toSRGB(v, FullBrightness); got > 0xf3 {
			t.Errorf("toSRGB(%d) = %d，超過 DAC 60 能到的 243", v, got)
		}
	}
}
