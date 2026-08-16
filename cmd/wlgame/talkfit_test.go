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
// 兩種框（`docs/spec/38` 之外的另一組常數）：
//
//	無肖像 `messageContentWidth` ＝ 22 全形 ＝ 352 px，一頁 5 列
//	有肖像 `talkTextWidth`        ＝ 160 px，框高 80 px
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
		// ① 無肖像的框：折行之後一頁 5 列，超過就要翻頁——翻頁是原版
		//    就有的行為（`sub_18810` 的分頁），所以只看**單行寬度**。
		for _, line := range textdraw.WrapLines(lines, messageContentWidth) {
			if w := textdraw.StringWidth(line); w > messageContentWidth {
				t.Errorf("#%d 折行後仍有 %d px 的行（上限 %d）：%q",
					i, w, messageContentWidth, line)
				overWide++
			}
		}
		// ② 有肖像的框：160 px 寬、框高 80 px。扣掉上下各 8 px 的內縮，
		//    一列 16 px ⇒ **最多 4 列**。
		wrapped := textdraw.WrapLines(lines, talkTextWidth)
		if n := len(wrapped); n > talkBoxRows {
			overRows++
			if overRows <= 5 {
				t.Logf("#%d 在肖像框裡要 %d 列（上限 %d）：%s",
					i, n, talkBoxRows, strings.Join(lines, "／"))
			}
		}
	}
	if overWide > 0 {
		t.Errorf("共 %d 行折行後仍超寬", overWide)
	}
	// **釘住這個數字**：現在是 5 則，其中四則（#77／#102／#166／#230）
	// 是五選一的選單，本來就不走肖像框；剩下 #256 是變數撐長的長台詞。
	// 校訂讓譯文變長時這個數字會漲，那正是要被看見的訊號。
	const knownOverRows = 5
	if overRows > knownOverRows {
		t.Errorf("肖像框需要翻頁的則數 = %d，比已知的 %d 多——"+
			"新校訂讓某則變長了，要回頭看排版", overRows, knownOverRows)
	}
	t.Logf("肖像框需要翻頁的則數：%d", overRows)
}
