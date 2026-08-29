package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
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

// 空城被攻下：原版不進戰術畫面，自動判定完攻方贏就跳 #26
// （`sub_14ED7` 的 `loc_14EF1`，docs/spec/105 §4.5）。
// **只在被攻陷的是玩家的城時**——玩家自己去打空城那條路（`loc_14F2B` 的
// `cmp bx, 4200h` → `jz loc_14EE7`）判定完就回，一則訊息都不發。
func TestEmptyCityFallReportsTalk26(t *testing.T) {
	// 城裡沒有守軍軍團 → 走 fightGarrison（打城兵），不是 resolveCorpsBattle。
	stage := func(player, attacker int) (*World, *CorpsEvent) {
		w := &World{Player: player}
		w.Corps[1] = Corps{
			Alive: true, Faction: attacker, Morale: 100,
			Men: 6000, Node: 7, X: 1, Y: 1,
		}
		for i := range w.Corps[1].Units {
			w.Corps[1].Units[i] = combat.Unit{Men: 1000, Kind: army.Infantry}
		}
		w.Generals[1] = General{Martial: 90, Command: 90}
		w.Cities[7] = City{Owner: 3, Garrison: 1, Prevention: 50}
		w.Cities[8] = City{Owner: 3}
		w.Factions[3].Cities, w.Factions[3].Capital = 2, 8
		w.Factions[5].Cities = 1
		return w, &CorpsEvent{}
	}

	// 敵軍（5）攻下玩家（3）的空城。
	w, ev := stage(3, 5)
	if w.wantsTactical(1, -1) {
		t.Fatal("打空城不該進戰術畫面")
	}
	w.fightGarrison(1, ev, rng.NewFixed(5))
	if ev.Captured != 7 {
		t.Fatalf("空城沒被攻下（Captured=%d），這一場的佈局有問題", ev.Captured)
	}
	if len(ev.TalkNotices) != 1 {
		t.Fatalf("空城陷落應該剛好一則訊息，得到 %#v", ev.TalkNotices)
	}
	got := ev.TalkNotices[0]
	if got.Index != talkSiegeCityFallen {
		t.Errorf("索引 %#x，預期 %#x（#26）", got.Index, talkSiegeCityFallen)
	}
	// #26「{2}受到{1}兵馬的攻擊」：{2} 是據點、{1} 是攻方主將。
	if got.City != 7 || got.General != 1 {
		t.Errorf("變數 = 據點 %d／主將 %d，預期 7／1", got.City, got.General)
	}

	// 玩家（5）自己攻下敵方的空城 → 原版是靜的。
	w2, ev2 := stage(5, 5)
	w2.fightGarrison(1, ev2, rng.NewFixed(5))
	if ev2.Captured != 7 {
		t.Fatalf("玩家沒攻下空城（Captured=%d）", ev2.Captured)
	}
	if len(ev2.TalkNotices) != 0 {
		t.Errorf("玩家自己攻空城不該有訊息，得到 %#v", ev2.TalkNotices)
	}
}
