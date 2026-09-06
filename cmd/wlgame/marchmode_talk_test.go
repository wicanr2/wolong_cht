package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 行軍指示的抬頭要走 TALK #21，不是內建 fallback（docs/spec/87 §6）。
//
// ⭐ **fallback 與 TALK #21 的中文一字不差**，所以「畫面上有字」證明不了
// 任何事——繁中之下兩條路的輸出完全相同。這一支直接比對 `talkLines`
// 的回傳：ok 必須為真、內容必須含代進去的據點名。
//
// 這個 bug 原本的樣子是 `map[byte]string{2: dest}`——`Part.Marker` 存的是
// ASCII `'2'`，給數值 2 會讓 `Lines` fail-closed，靜靜走 fallback。
// 只有英文版才看得出來（fallback 是中文的）。
func TestMarchModePromptComesFromTalk(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("找不到原版素材：%v", err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	g := &game{lib: lib, world: w}

	const dest = "ZZTESTCITY"
	lines, ok := g.talkLines(0x15, map[byte]string{'2': dest})
	if !ok {
		t.Fatal("TALK #21 取不到——變數 key 給錯的話這裡就是 false，畫面會靜靜退回 fallback")
	}
	if len(lines) == 0 {
		t.Fatal("TALK #21 是空的")
	}
	if !strings.Contains(strings.Join(lines, ""), dest) {
		t.Errorf("據點名沒有代進去：%q", lines)
	}
	// ⭐ **這一則沒有 fallback**（docs/spec/140）：取不到就不掛框，
	// 與其他訊息路徑一樣 fail-closed。從前它是選單框的標題、還帶一份
	// 與原文一字不差的中文備援，於是「畫面上有字」永遠為真，
	// 走的是哪一條看不出來。現在唯一的判準就是上面那個 `ok`。
	g.setStatusTalk(0x15, map[byte]string{'2': dest})
	if !g.statusBoxActive() {
		t.Fatal("setStatusTalk 之後狀態列框應該掛著")
	}
	if !strings.Contains(strings.Join(g.statusBox.lines, ""), dest) {
		t.Errorf("狀態列框裡沒有據點名：%q", g.statusBox.lines)
	}
	g.setStatusTalk(0x15, map[byte]string{2: dest}) // 給錯 key
	if g.statusBoxActive() {
		t.Error("marker key 給錯時應該什麼都不掛，不是掛一個半殘的框")
	}

	// 給錯 key 就是這個 bug 原本的樣子——fail-closed，什麼都不說。
	if _, ok := g.talkLines(0x15, map[byte]string{2: dest}); ok {
		t.Error("數值 key 竟然也通得過——marker 的定義變了，這一支要重寫")
	}
}
