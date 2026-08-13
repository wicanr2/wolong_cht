package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

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
