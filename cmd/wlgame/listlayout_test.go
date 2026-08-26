package main

import "testing"

// 羅馬化人名比原版的三個全形字寬，欄寬要裁——但要裁得讀得出來
// （docs/spec/84 §6）。
func TestFitCellTruncatesRomanisedNames(t *testing.T) {
	const room = 56 // 7 個半形格
	for _, tc := range []struct{ in, want string }{
		{"CAO-CAO", "CAO-CAO"},    // 塞得下就原樣
		{"XIAHOU-YUAN", "XIAHOU"}, // 切在連字號上 → 連字號也去掉
		{"ZHUGE-LIANG", "ZHUGE"},  // 名只剩一個字母 → 只留姓
		{"XU-HUANG", "XU-HUAN"},   // 名還剩得夠多 → 照切
		{"曹操", "曹操"},              // 中文照原樣（3 字 48px 本來就塞得下）
	} {
		if got := fitCell(tc.in, room); got != tc.want {
			t.Errorf("fitCell(%q, %d) = %q，預期 %q", tc.in, room, got, tc.want)
		}
	}
	if got := fitCell("CAO-CAO", 0); got != "CAO-CAO" {
		t.Errorf("欄寬 0（未知）時不該裁：%q", got)
	}
}
