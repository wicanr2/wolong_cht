package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 身分欄照原版 `sub_1770C`：**職務值優先，君主墊底**（docs/spec/143 §1）。
//
//	al = [si+17h] ; al == 0 && [si] & 40h → al = 5
//
// ⭐ 君主兼軍團長要顯示「軍團長」不是「君主」——先前從
// `Factions[].Lord` 反查的寫法在這一格會反過來。
func TestGeneralRankFollowsDutyFirst(t *testing.T) {
	newGame := func(duty int, sovereign bool) *game {
		w := &state.World{Player: 0}
		w.Factions[0].Alive = true
		w.Factions[0].Lord = 3
		w.Generals[3].Alive = true
		w.Generals[3].Faction = 0
		w.Generals[3].Captor = 0xFF
		w.Generals[3].Duty = duty
		w.Generals[3].Sovereign = sovereign
		return &game{world: w}
	}
	cases := []struct {
		duty      int
		sovereign bool
		want      int
		name      string
	}{
		{state.DutyNone, false, 0, "無職"},
		{state.DutyNone, true, 5, "君主"},
		{state.DutyCorpsLeader, true, 1, "君主兼軍團長 → 軍團長"},
		{state.DutyGovernor, false, 2, "內政官"},
		{state.DutyDiplomat, false, 3, "外交官"},
		{state.DutyCaptive, false, 4, "俘虜"},
		{state.DutyCaptive, true, 4, "被俘的君主 → 俘虜"},
	}
	for _, c := range cases {
		if got := newGame(c.duty, c.sovereign).generalRank(3); got != c.want {
			t.Errorf("%s：generalRank = %d（%s），want %d（%s）",
				c.name, got, listRankNames[got], c.want, listRankNames[c.want])
		}
	}
}

// 軍師既不在編成候選也不在人事候選（docs/spec/144）。
// ⭐ **兩張都要驗**：先前只有 `formCandidates` 比對 `Advisor` 編號，
// `freeGenerals`（人事任命）漏了——同一個機制貼兩張補丁，貼漏一張。
func TestAdvisorInNeitherCandidateList(t *testing.T) {
	w := &state.World{Player: 0}
	w.Factions[0].Alive = true
	w.Factions[0].Lord = 3
	w.Factions[0].Advisor = 7
	for _, i := range []int{3, 4, 7} {
		w.Generals[i].Alive = true
		w.Generals[i].Faction = 0
		w.Generals[i].Duty = state.DutyNone
		w.Generals[i].Captor = 0xFF
	}
	w.Factions[0].Generals = 3
	g := &game{world: w, lordCorps: true}

	has := func(rows []int, n int) bool {
		for _, r := range rows {
			if r == n {
				return true
			}
		}
		return false
	}
	if !has(g.formCandidates(), 7) || !has(g.freeGenerals(), 7) {
		t.Fatal("還沒選軍師，武將 7 就不在候選裡")
	}

	w.TakeAdvisor(0, 7)

	if has(g.formCandidates(), 7) {
		t.Error("軍師還在編成候選裡")
	}
	if has(g.freeGenerals(), 7) {
		t.Error("軍師還在人事候選裡")
	}
}
