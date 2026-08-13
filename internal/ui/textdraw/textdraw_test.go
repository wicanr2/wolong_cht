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
