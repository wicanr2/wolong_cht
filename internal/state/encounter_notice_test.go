package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
)

// 進戰術畫面前那一則訊息（docs/spec/105 §1）：野戰 #0x1D 的兩個 `{1}` 是
// 攻守兩個主將；攻城看玩家站哪一邊選 #0x1C／#0x1B。
func TestEncounterNoticePicksOriginalTalkIndex(t *testing.T) {
	w := &World{Player: 3}
	w.Corps[1] = Corps{Alive: true, Faction: 3, Node: 42} // 玩家的軍團
	w.Corps[2] = Corps{Alive: true, Faction: 5, Node: 42} // 敵方

	field := w.encounterNotice(1, 2, combat.Field)
	if field.Index != talkFieldEncounter {
		t.Fatalf("野戰用 %#x，預期 %#x", field.Index, talkFieldEncounter)
	}
	if len(field.SeqGenerals) != 2 ||
		field.SeqGenerals[0] != 1 || field.SeqGenerals[1] != 2 {
		t.Fatalf("野戰的兩個 {1} = %v，預期 [玩家 對方]", field.SeqGenerals)
	}
	// 玩家是守方時原版先 xchg si,di：第一個 {1} 仍是玩家那一方。
	held := w.encounterNotice(2, 1, combat.Field)
	if len(held.SeqGenerals) != 2 ||
		held.SeqGenerals[0] != 1 || held.SeqGenerals[1] != 2 {
		t.Fatalf("玩家守方的兩個 {1} = %v，預期 [玩家 對方]", held.SeqGenerals)
	}
	if field.City != -1 {
		t.Fatalf("野戰不該帶據點，得到 %d", field.City)
	}

	out := w.encounterNotice(1, 2, combat.Siege) // 玩家是攻方
	if out.Index != talkSiegeOutgoing || out.General != 1 || out.City != 42 {
		t.Fatalf("玩家攻城 = %#v", out)
	}
	in := w.encounterNotice(2, 1, combat.Siege) // 玩家是守方
	if in.Index != talkSiegeIncoming || in.General != 2 || in.City != 42 {
		t.Fatalf("玩家守城 = %#v（{1} 一律是攻方主將）", in)
	}
}

// 訊息要真的送到 Event.TalkNotices——只掛在 CorpsEvent 上呈現層看不到。
func TestEncounterNoticeReachesTheEventQueue(t *testing.T) {
	ev := &CorpsEvent{}
	w := &World{Player: 3}
	w.Corps[1] = Corps{Alive: true, Faction: 3, Node: 7}
	ev.TalkNotices = append(ev.TalkNotices, w.encounterNotice(1, 2, combat.Field))
	e := Event{Corps: []CorpsEvent{*ev}}
	for i := range e.Corps {
		e.TalkNotices = append(e.TalkNotices, e.Corps[i].TalkNotices...)
	}
	if len(e.TalkNotices) != 1 || e.TalkNotices[0].Index != talkFieldEncounter {
		t.Fatalf("Event.TalkNotices = %#v", e.TalkNotices)
	}
}
