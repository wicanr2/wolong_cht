// Package talkmenu 把 `TALK.DAT` 裡的「一則多行訊息」當成選單用。
//
// 原版有好幾個選單就是這樣存的：五個選項放在同一則訊息的五行裡
//（進言的五項是 #77、三組說服理由是 #102／#166／#230）。
// **行的順序就是資料**——索引算式直接吃那個位置。
package talkmenu

import "github.com/wicanr2/wolong_cht/internal/assets/text"

// Labels 取一則訊息的每一行當選單列。
//
// ⛔ **不要 trim。** 那些全形空白是版面的一部分：框寬由字數決定，
// 每一列要一樣寬。`text.Decode` 已經處理過行尾，這裡再砍就連行首也沒了。
//
// 行數對不上 fallback 時整份退回 fallback——半份原文半份備援
// 會讓畫面看起來正常，而錯的那幾列沒人會發現。
func Labels(t *text.Table, index int, vars map[byte]string, fallback []string) []string {
	out := append([]string(nil), fallback...)
	lines, ok := t.Lines(index, vars)
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
