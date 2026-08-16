package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// TestAllTalkLinesFitTheirBox 是 M7 排版 parity 的全量檢查：
// **1,022 則逐則量渲染寬度**，看有沒有塞不進訊息框的。
//
// **原版只有一個訊息框**（`docs/spec/41`）：文字區 `talkTextWidth`
// ＝ 160 px ＝ 10 全形字，框高 80 px 扣掉上下內縮 ⇒ 4 列。
//
// ⚠ **這一支不是「有沒有錯字」的檢查**，是「排版會不會爆框」。
// 校訂過的譯文比原文長一兩個字很常見，而框寬是原版定死的。
func TestAllTalkLinesFitTheirBox(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	if lib.Talk == nil {
		t.Skip("TALK.DAT 沒有載入")
	}
	g := &game{lib: lib}

	// 變數用固定長度的替身：三個全形字是人名／地名的常見長度
	// （`docs/re/25` 的變數表），這樣量到的是「最壞情況的一種」。
	vars := map[byte]string{}
	for _, m := range []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'} {
		vars[m] = "○○○"
	}

	overWide, overRows := 0, 0
	for i := 0; i < len(lib.Talk.Messages); i++ {
		lines, ok := g.talkLines(i, vars)
		if !ok || len(lines) == 0 {
			continue
		}
		wrapped := textdraw.WrapLines(lines, talkTextWidth)
		// ① 單行寬度。折不開的只有「一個變數就吃掉三個全形字」那幾則。
		for _, line := range wrapped {
			if w := textdraw.StringWidth(line); w > talkTextWidth {
				overWide++
				if overWide <= 5 {
					t.Logf("#%d 折行後仍有 %d px 的行（上限 %d）：%q",
						i, w, talkTextWidth, line)
				}
			}
		}
		// ② 列數：框高 80 px 扣掉上下各 8 px 的內縮，一列 16 px ⇒ 最多 4 列。
		if n := len(wrapped); n > talkBoxRows {
			overRows++
			if overRows <= 5 {
				t.Logf("#%d 在肖像框裡要 %d 列（上限 %d）：%s",
					i, n, talkBoxRows, strings.Join(lines, "／"))
			}
		}
	}
	// **釘住這兩個數字**，校訂讓譯文變長時它們會漲，那正是要被看見的訊號。
	//
	// 超寬 4 行：全部是「一個變數 ＋ 一句話」擠在同一行的句型，
	// 而**替身用了三個全形字**（人名的長端）。原版的原文也是這樣寫的，
	// 遇到三字人名時同樣會壓到框邊——這是原版行為，不是校訂造成的。
	//
	// 需要第 5 列的 5 則：四則（#77／#102／#166／#230）是五選一的選單，
	// 由選單常式畫、不進這個框；剩下 #256 是兩個變數撐長的長台詞。
	const knownOverWide, knownOverRows = 4, 5
	if overWide > knownOverWide {
		t.Errorf("超寬的行數 = %d，比已知的 %d 多——"+
			"新校訂讓某一行變長了，要回頭看排版", overWide, knownOverWide)
	}
	if overRows > knownOverRows {
		t.Errorf("需要翻頁的則數 = %d，比已知的 %d 多——"+
			"新校訂讓某則變長了，要回頭看排版", overRows, knownOverRows)
	}
	t.Logf("超寬 %d 行、需要翻頁 %d 則", overWide, overRows)
}

// 訊息框的版面常數逐一對原版（docs/spec/41 §1）。
//
// 這一支存在的理由是**版面常數最容易被「順手調一下」改掉**——
// 它們的出處是機器碼，不是眼睛。
func TestMessageBoxMatchesOriginalGeometry(t *testing.T) {
	for _, tc := range []struct {
		name string
		got  int
		want int
	}{
		{"框 X", talkBoxX, 160},
		{"框 Y", talkBoxY, 160},
		{"框寬", talkBoxW, 256},
		{"框高", talkBoxH, 80},
		{"肖像 X", talkPortraitX, 168},
		{"肖像 Y", talkPortraitY, 168},
		{"文字 X", talkTextX, 240},
		{"文字 Y", talkTextY, 176},
		{"文字寬", talkTextWidth, 160},
		{"行距", talkLinePitch, 16},
		{"列數", talkBoxRows, 4},
		{"一般通知的肖像頁", defaultPortraitPage, 0x93},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %d，原版是 %d", tc.name, tc.got, tc.want)
		}
	}
}

// 三個框內關係：文字要讓開 64 px 的肖像、一整列字不能超出右內緣、
// 四列字不能掉出框底。
//
// ⚠ **這三條不是從常數推出來的，是對常數的獨立約束。**
// 只比對數值的那一支（上面）在三個數字一起被改錯時會一起通過。
func TestTextDoesNotOverlapPortrait(t *testing.T) {
	const portraitSize = 64
	if right := talkPortraitX + portraitSize; talkTextX < right {
		t.Errorf("文字從 %d 開始，肖像畫到 %d", talkTextX, right)
	}
	if end := talkTextX + talkTextWidth; end > talkBoxX+talkBoxW-8 {
		t.Errorf("一整列字畫到 %d，框的右內緣在 %d",
			end, talkBoxX+talkBoxW-8)
	}
	if bottom := talkTextY + talkBoxRows*talkLinePitch; bottom > talkBoxY+talkBoxH {
		t.Errorf("%d 列字畫到 y=%d，框底在 %d",
			talkBoxRows, bottom, talkBoxY+talkBoxH)
	}
}
