package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// TestEvent9LongNotificationRoute 驗證 state 的長程事件 9 結果交給 GUI 後，
// 只有釋放回玩家勢力的武將建立 TALK #37；非玩家與在野結果不會污染 modal
// queue，且後續玩家通知仍依原始事件順序追加。
func TestEvent9LongNotificationRoute(t *testing.T) {
	lib, err := library.LoadWithOptions("../../workplace/orig/dosv", library.LoadOptions{
		TalkJSON: "../../translations/talk-dosv-corrected.json",
	})
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	ids := []int{0, 1, 2}
	for _, id := range ids {
		if id >= len(w.Generals) {
			t.Fatalf("事件 9 fixture General #%d 不存在", id)
		}
	}

	g := &game{lib: lib, world: w}
	for _, tc := range []struct {
		name      string
		general   int
		faction   int
		wantDelta int
		wantName  bool
	}{
		{"player-release", ids[0], w.Player, 1, true},
		{"other-faction-release", ids[1], (w.Player + 1) % len(w.Factions), 0, false},
		{"unaffiliated-release", ids[2], -1, 0, false},
		{"later-player-release", ids[0], w.Player, 1, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w.Generals[tc.general].Faction = tc.faction
			before := len(g.messages)
			g.enqueueEventMessages(state.Event{ReleasedGenerals: []int{tc.general}})
			if got := len(g.messages) - before; got != tc.wantDelta {
				t.Fatalf("#37 modal 增量 = %d，want %d；queue=%#v",
					got, tc.wantDelta, g.messages)
			}
			if tc.wantDelta == 0 {
				return
			}
			assertTalkModalContract(t, g, before+1)
			if tc.wantName {
				name := big5(w.Generals[tc.general].Name)
				content := strings.Join(g.messages[len(g.messages)-1].lines, "")
				if name != "" && !strings.Contains(content, name) {
					t.Fatalf("#37 未展開武將名 %q：%q", name, content)
				}
			}
		})
	}
	if len(g.messages) != 2 {
		t.Fatalf("長程事件 9 最終 modal 數 = %d，want 2：%#v", len(g.messages), g.messages)
	}
}
