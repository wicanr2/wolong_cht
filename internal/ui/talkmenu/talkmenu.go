// Package talkmenu 把 `TALK.DAT` 裡的「一則多行訊息」當成選單用。
//
// 原版有好幾個選單就是這樣存的：五個選項放在同一則訊息的五行裡
//（進言的五項是 #77、三組說服理由是 #102／#166／#230）。
// **行的順序就是資料**——索引算式直接吃那個位置。
package talkmenu

import "github.com/wicanr2/wolong_cht/internal/assets/text"

// Labels 取一則訊息的每一行當選單列，**行尾的全形空白已經被切掉**。
//
// 給「把標籤接進另一段文字」的地方用（手機版的標題是
// `名稱 + "：選對象"`，帶著補位空白會多出一個洞）。
// 要畫原版那種**框寬由字數決定**的選單框，用 MenuLabels。
//
// 行數對不上 fallback 時整份退回 fallback——半份原文半份備援
// 會讓畫面看起來正常，而錯的那幾列沒人會發現。
func Labels(t *text.Table, index int, vars map[byte]string, fallback []string) []string {
	return pick(t.Lines(index, vars))(fallback)
}

// MenuLabels 與 Labels 相同，但**保留行尾的全形空白**。
//
// ⛔ 那些空白是版面的一部分：原版把每一列補到等寬，而框寬就是
// 「第一列有幾個全形字 ＋ 1」格（`docs/spec/45` §2.2）。
// 砍掉的話 `TALK #79` 的「　位置確認　」變成 5 個字，框窄 16 px
// （`docs/spec/124`）。
func MenuLabels(t *text.Table, index int, vars map[byte]string, fallback []string) []string {
	return pick(t.MenuLines(index, vars))(fallback)
}

func pick(lines []string, ok bool) func([]string) []string {
	return func(fallback []string) []string {
		out := append([]string(nil), fallback...)
		if !ok || len(lines) != len(out) {
			return out
		}
		for i, l := range lines {
			if l != "" {
				out[i] = l
			}
		}
		return out
	}
}
