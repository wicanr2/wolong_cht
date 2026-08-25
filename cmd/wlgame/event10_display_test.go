package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 事件 10 的訊息要真的走到畫面訊息框，不能停在 state 的 TalkNotice。
//
// 整條鏈：QueueEvent10 → dispatch（sub_13496 邊界：AH=General、
// DX=TALK index）→ Event.TalkNotices → enqueueEventMessages →
// 訊息框逐句顯示（\1 代入武將名）。這支測的是「接線存在」，
// 不是原版 producer——原版自然 producer 仍是 unknown（docs/re/15）。
func TestEvent10NoticeReachesMessageBox(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	general := -1
	for i, gen := range w.Generals {
		if gen.Alive {
			general = i
			break
		}
	}
	if general < 0 {
		t.Fatal("找不到活著的武將")
	}
	// #66（0x42）是近似 producer 的來投句型，\1 代入武將名。
	if !w.QueueEvent10(general, 0x42) {
		t.Fatal("事件 10 排不進佇列")
	}

	g := &game{lib: lib, world: w}
	// 佇列的 dispatch 有節流（每 7／10 個每時邊界取一筆），跑到取出為止。
	r := rng.NewFixed(1)
	for i := 0; i < 200000 && len(g.messages) == 0; i++ {
		g.enqueueEventMessages(w.Tick(r))
	}
	if len(g.messages) == 0 {
		t.Fatal("事件 10 的訊息沒有走到訊息框")
	}
	name := strings.TrimSpace(big5(w.Generals[general].Name))
	joined := strings.Join(g.messages[0].lines, "")
	if name == "" || !strings.Contains(joined, name) {
		t.Fatalf("訊息裡沒有代入武將名 %q：%#v", name, g.messages[0].lines)
	}
}
