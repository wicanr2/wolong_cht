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

// 反白不是「換一個底色」，是把色號 XOR 12（`sub_10B46`，docs/spec/124）。
// 黑底白字那一族因此變成黃底藍字。
func TestHighlightIsXorOfInkAndPaper(t *testing.T) {
	if got, want := InkIndex^HighlightXOR, HighlightIndex; got != want {
		t.Errorf("底 %d XOR %d ＝ %d，want %d", InkIndex, HighlightXOR, got, want)
	}
	if got, want := PaperIndex^HighlightXOR, HighlightInkIndex; got != want {
		t.Errorf("字 %d XOR %d ＝ %d，want %d", PaperIndex, HighlightXOR, got, want)
	}
	// 實機量到的兩個顏色（parity-tap5/menu.png，docs/playtest/60）。
	if HighlightIndex != 12 || HighlightInkIndex != 3 {
		t.Errorf("反白色號 ＝ (%d, %d)，實機量到 (12, 3)",
			HighlightIndex, HighlightInkIndex)
	}
}

// 反白是**整塊 XOR 12**，清單的反白條也不例外（docs/spec/124 §3.6）。
// 先前 SelectIndex 是一個獨立的量測值，於是「同一條規則」看起來像三條。
func TestSelectIsSheetXorHighlight(t *testing.T) {
	if got, want := SheetIndex^HighlightXOR, SelectIndex; got != want {
		t.Errorf("清單底 %d XOR %d ＝ %d，want %d",
			SheetIndex, HighlightXOR, got, want)
	}
	// docs/playtest/98 §3 量到的三組，逐組驗一次。
	for _, tc := range []struct{ name string; normal, selected int }{
		{"清單底色", SheetIndex, SelectIndex},
		{"一般文字", InkIndex, HighlightIndex},
		{"換色的儲存格", 10, 6},
	} {
		if got := tc.normal ^ HighlightXOR; got != tc.selected {
			t.Errorf("%s：%d XOR 12 ＝ %d，實機量到 %d",
				tc.name, tc.normal, got, tc.selected)
		}
	}
}
