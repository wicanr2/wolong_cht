package textdraw

import (
	"image"
	"image/color"
	"testing"
)

func TestStringWidthMatchesMixedHalfAndFullWidth(t *testing.T) {
	if got, want := StringWidth("A中2"), HalfW+GlyphW+HalfW; got != want {
		t.Fatalf("混排寬度 = %d，want %d", got, want)
	}
}

func TestWrapLineKeepsClosingPunctuationWithPreviousLine(t *testing.T) {
	got := WrapLine("甲乙。", GlyphW)
	if len(got) != 2 || got[0] != "甲" || got[1] != "乙。" {
		t.Fatalf("中文標點換行 = %#v，預期標點留在前一字後", got)
	}
}

func TestWrapLinesPreservesHardBlankLines(t *testing.T) {
	got := WrapLines([]string{"甲乙", "", "A2"}, GlyphW)
	if len(got) != 4 || got[0] != "甲" || got[1] != "乙" || got[2] != "" ||
		got[3] != "A2" {
		// A2 兩個半形字可同列，故總列數應為四列；保留這個
		// 明確訊息，讓日後調整半形規格時不會靜默改變 hard line。
		t.Fatalf("硬斷行／空列 = %#v", got)
	}
}

// 英文要斷在空白：從字中間切開的話讀不出來（docs/spec/84 §2）。
func TestWrapLineBreaksEnglishAtSpaces(t *testing.T) {
	got := WrapLine("The strategist reports", 20*HalfW)
	if len(got) != 2 || got[0] != "The strategist" || got[1] != "reports" {
		t.Fatalf("英文折行 = %#v，預期斷在空白", got)
	}
	// 單字本身比一列長時沒得搬，照舊硬切——不能因此整列不折。
	long := WrapLine("AAAAAAAAAAAAAAAAAAAAAAAAA", 10*HalfW)
	if len(long) != 3 {
		t.Fatalf("超長單字 = %#v，預期照舊硬切成三列", long)
	}
	// 中文不受影響：每個字都能斷。
	zh := WrapLine("甲乙丙丁", 2*GlyphW)
	if len(zh) != 2 || zh[0] != "甲乙" {
		t.Fatalf("中文折行被改壞了 = %#v", zh)
	}
}

// Scale2x 要把斜線的階梯補成 45 度，直線與空隙不動（docs/spec/101）。
func TestScale2xFillsDiagonalStairs(t *testing.T) {
	// 一條 2×2 的斜線：(0,0)、(1,1)。
	a := image.NewAlpha(image.Rect(0, 0, 3, 3))
	a.SetAlpha(0, 0, color.Alpha{255})
	a.SetAlpha(1, 1, color.Alpha{255})
	out := Scale2x(a)
	if out.Bounds().Dx() != 6 || out.Bounds().Dy() != 6 {
		t.Fatalf("尺寸 %v，預期 6×6", out.Bounds())
	}
	// 斜線的內角：(0,0) 的右下角與 (1,1) 的左上角應被補上。
	if out.AlphaAt(1, 1).A == 0 || out.AlphaAt(2, 2).A == 0 {
		t.Fatal("斜線的內角沒有補上")
	}
	// 孤立的一點只放大成 2×2，不長出去。
	dot := image.NewAlpha(image.Rect(0, 0, 3, 3))
	dot.SetAlpha(1, 1, color.Alpha{255})
	do := Scale2x(dot)
	n := 0
	for y := 0; y < 6; y++ {
		for x := 0; x < 6; x++ {
			if do.AlphaAt(x, y).A != 0 {
				n++
			}
		}
	}
	if n != 4 {
		t.Fatalf("孤立點放大後有 %d 格，預期 4", n)
	}
	// 一條直線放大後仍是實心矩形。
	l := image.NewAlpha(image.Rect(0, 0, 3, 1))
	for x := 0; x < 3; x++ {
		l.SetAlpha(x, 0, color.Alpha{255})
	}
	lo := Scale2x(l)
	for y := 0; y < 2; y++ {
		for x := 0; x < 6; x++ {
			if lo.AlphaAt(x, y).A != 255 {
				t.Fatalf("直線在 (%d,%d) 破了", x, y)
			}
		}
	}
}
