package textdraw

import "testing"

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
