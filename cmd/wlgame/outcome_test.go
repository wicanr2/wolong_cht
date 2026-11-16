package main

import (
	"reflect"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 原版13DC9傳的是019E群組。414是褒獎；逐出文字在470+君主raw型。
func TestTrustOutcomeUsesLordGroupAndAdvisorName(t *testing.T) {
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
	w.Player, w.Trust = 0, 1
	w.AdjustTrust(-1)
	lord := w.Factions[w.Player].Lord
	g := &game{lib: lib, world: w}
	notice := state.TalkNotice{Index: 0x019E, General: lord, City: -1, Faction: -1, Amount: -1}
	vars, ok := w.TalkNoticeVars(notice, big5)
	if !ok {
		t.Fatal("missing advisor parameters")
	}
	for variant := 0; variant < 3; variant++ {
		w.Generals[lord].TalkVariant = variant
		want, ok := lib.Talk.Lines(470+variant, vars)
		if !ok || !reflect.DeepEqual(g.outcomeLines(), want) {
			t.Fatalf("raw%d: got %q want %q", variant, g.outcomeLines(), want)
		}
	}
}

func TestOutcomeModalReasonAndHitRect(t *testing.T) {
	if got := outcomeReason(state.DefeatTrustZero); got != "信賴度歸零，已被逐出勢力。" {
		t.Fatalf("trust reason = %q", got)
	}
	if got := outcomeReason(state.DefeatFactionEliminated); got != "玩家勢力失去最後可替代據點。" {
		t.Fatalf("faction reason = %q", got)
	}
	if !outcomeConfirmRect().In(outcomePanel) {
		t.Fatal("outcome confirm rect 超出 modal")
	}
	g := &game{world: &state.World{Trust: 1}}
	g.world.AdjustTrust(-1)
	lines := g.outcomeLines()
	if len(lines) != 1 || lines[0] != "信賴度歸零，已被逐出勢力。" {
		t.Fatalf("缺少 TALK 資產時應使用 fallback：%q", lines)
	}
}

func TestOutcomeFreezesGameAndConfirmReturnsLauncher(t *testing.T) {
	g := &game{world: &state.World{Trust: 1}, saveFile: "/definitely/not/a/save"}
	g.world.AdjustTrust(-1)
	if g.timeRuns() {
		t.Fatal("outcome 後 timeRuns 不應為 true")
	}
	if err := g.returnToLauncher(); err != nil {
		t.Fatal(err)
	}
	if g.world != nil || g.launcher == nil {
		t.Fatalf("確認 outcome 後未回 launcher：world=%v launcher=%v", g.world, g.launcher)
	}
}
