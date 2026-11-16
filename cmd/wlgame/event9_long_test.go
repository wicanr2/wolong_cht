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

	w.Generals[ids[0]].TalkVariant = 7
	g := &game{lib: lib, world: w}
	for _, tc := range []struct {
		name      string
		general   int
		faction   int
		wantDelta int
		wantName  bool
	}{
		{"player-release", ids[0], w.Player, 2, true},
		{"other-faction-release", ids[1], (w.Player + 1) % len(w.Factions), 0, false},
		{"unaffiliated-release", ids[2], -1, 0, false},
		{"later-player-release", ids[0], w.Player, 2, true},
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
			assertTalkModalContract(t, g, before+2)
			if tc.wantName {
				name := big5(w.Generals[tc.general].Name)
				content := strings.Join(g.messages[before].lines, "")
				if name != "" && !strings.Contains(content, name) {
					t.Fatalf("#37 未展開武將名 %q：%q", name, content)
				}
			}
		})
	}
	if len(g.messages) != 4 {
		t.Fatalf("長程事件 9 最終 modal 數 = %d，want 4：%#v", len(g.messages), g.messages)
	}
}

// 原版15107／15110：通報者37之後接武將本人的199h組，不可直接讀空槽409。
func TestEvent9ReturnedGeneralReplyVariants(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	w.Generals[37].Faction = 0
	for variant := 0; variant < 8; variant++ {
		w.Generals[37].TalkVariant = variant
		g := &game{lib: lib, world: w}
		g.enqueueEventMessages(state.Event{ReleasedGenerals: []int{37}})
		count := 1
		if variant >= 3 {
			count = 2
		}
		if len(g.messages) != count {
			t.Fatalf("raw%d messages=%d want%d", variant, len(g.messages), count)
		}
		if g.messages[0].portraitPage != 147 {
			t.Fatal("return notice must use reporter")
		}
		if count == 2 {
			m := g.messages[1]
			lines, ok := g.talkLines(430+variant, nil)
			if !ok || strings.Join(m.lines, "\n") != strings.Join(layoutMessageLines(lines), "\n") {
				t.Fatalf("raw%d does not match original TALK%d", variant, 430+variant)
			}
			if m.portraitPage != 62 || m.scene != -1 || m.upper || m.lower {
				t.Fatalf("raw%d reply uses wrong speaker or box: %+v", variant, m)
			}
		}
	}
	// Return to a nonplayer faction or unaffiliated state still has no notification.
	for _, faction := range []int{1, 255} {
		w.Generals[37].Faction = faction
		g := &game{lib: lib, world: w}
		g.enqueueEventMessages(state.Event{ReleasedGenerals: []int{37}})
		if len(g.messages) != 0 {
			t.Fatal("nonplayer release must stay silent")
		}
	}
}
