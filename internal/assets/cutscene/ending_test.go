package cutscene

import (
	"strings"
	"testing"
)

// 結尾文字燒在 `D7END.EXE` 裡（段內 0x5F0，200 個全形字），
// 版面是一行 20 字（`sub_10238` 的 x 從 0x40 到 0x180、每字 0x10）。
func TestEndingTextLayout(t *testing.T) {
	lines, err := EndingText(dir)
	if err != nil {
		t.Skip("找不到原版 D7END.EXE，跳過：" + err.Error())
	}
	if len(lines) != TextChars/TextCols {
		t.Fatalf("切成 %d 行，預期 %d", len(lines), TextChars/TextCols)
	}
	// 尾端的全形空白會被 text.Decode 修掉，所以只有最後一行可以短。
	for i, l := range lines {
		n := len([]rune(l))
		if i < len(lines)-1 && n != TextCols {
			t.Errorf("第 %d 行 %d 字，預期 %d：%q", i, n, TextCols, l)
		}
		if n > TextCols {
			t.Errorf("第 %d 行 %d 字，超過一行 %d 字：%q", i, n, TextCols, l)
		}
	}
	// 第一行與最後一行是可辨識的定錨——版面或編碼錯了會先壞在這裡。
	if !strings.Contains(lines[0], "英雄") {
		t.Errorf("第 0 行 = %q，預期含「英雄」", lines[0])
	}
	if !strings.Contains(lines[len(lines)-1], "天下第一軍師") {
		t.Errorf("最後一行 = %q，預期含「天下第一軍師」", lines[len(lines)-1])
	}
}

// 十二幕都要載得起來，而且各自用自己那一組色盤。
func TestLoadEndingFrames(t *testing.T) {
	e, err := LoadEnding(dir)
	if err != nil {
		t.Skip("找不到原版素材，跳過：" + err.Error())
	}
	for n, f := range e.Frames {
		if f == nil {
			t.Fatalf("第 %d 幕沒有畫面", n)
		}
		if b := f.Bounds(); b.Dx() != Width || b.Dy() != Height {
			t.Fatalf("第 %d 幕是 %v，預期 %dx%d", n, b, Width, Height)
		}
	}
	// 相鄰兩幕不該畫出一模一樣的東西——色盤組接錯時最常見的症狀。
	if string(e.Frames[1].Pix) == string(e.Frames[2].Pix) {
		t.Fatal("第 1 幕與第 2 幕逐像素相同")
	}
}
