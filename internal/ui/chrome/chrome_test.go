package chrome

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

// TestChromeColoursComeFromPalette 釘住 docs/spec/54：五個介面顏色
// **跟著 `GAMEPAL.BRG` 走**，不是手抄的 RGB。
//
// 手抄的值會與解碼脫鉤——docs/spec/51 把調色盤換算改成走 VGA 的
// 6 bit DAC 之後，舊的常數每一個都差了 2–4，而畫面看起來照樣正常。
func TestChromeColoursComeFromPalette(t *testing.T) {
	dir := filepath.Join("..", "..", "..", "workplace", "orig", "dosv")
	if _, err := os.Stat(filepath.Join(dir, "GAMEPAL.BRG")); err != nil {
		t.Skipf("找不到原版素材，跳過：%v", err)
	}
	lib, err := library.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	Load(lib, 0)
	for _, c := range []struct {
		name string
		idx  int
		r    uint8
		g    uint8
		b    uint8
	}{
		{"Menu", MenuIndex, Menu.R, Menu.G, Menu.B},
		{"Sheet", SheetIndex, Sheet.R, Sheet.G, Sheet.B},
		{"Select", SelectIndex, Select.R, Select.G, Select.B},
		{"Ink", InkIndex, Ink.R, Ink.G, Ink.B},
		{"Blank", InkIndex, Blank.R, Blank.G, Blank.B},
		{"Paper", PaperIndex, Paper.R, Paper.G, Paper.B},
	} {
		want, err := lib.PaletteColor(0, c.idx)
		if err != nil {
			t.Fatalf("%s：取不到色 %d：%v", c.name, c.idx, err)
		}
		if c.r != want.R || c.g != want.G || c.b != want.B {
			t.Errorf("%s = #%02x%02x%02x，調色盤色 %d 是 #%02x%02x%02x",
				c.name, c.r, c.g, c.b, c.idx, want.R, want.G, want.B)
		}
	}
}
