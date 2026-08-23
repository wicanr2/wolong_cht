package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// ⚠ 君主不在編成候選裡（docs/spec/76）。原版讓君主帶兵走的是出陣那條，
// 不是指令列的「編成」。
func TestFormCandidatesExcludeLord(t *testing.T) {
	w := &state.World{Player: 0}
	w.Factions[0].Alive = true
	w.Factions[0].Lord = 3
	for _, i := range []int{3, 4, 5} {
		w.Generals[i].Alive = true
		w.Generals[i].Faction = 0
		w.Generals[i].Posted = false
		w.Generals[i].Captor = 0xFF
	}
	g := &game{world: w}
	rows := g.formCandidates()

	for _, r := range rows {
		if r == 3 {
			t.Error("君主（武將 3）出現在編成候選裡")
		}
	}
	// 反向對照：其他兩位要在——否則「候選永遠是空的」也會讓上面那條通過。
	got := map[int]bool{}
	for _, r := range rows {
		got[r] = true
	}
	for _, want := range []int{4, 5} {
		if !got[want] {
			t.Errorf("武將 %d 應該可以帶兵，卻不在候選裡", want)
		}
	}
}

// 已帶兵、俘虜、別的勢力、死掉的都不算——這些是原本就有的條件，
// 加君主那一行不可以把它們弄壞。
func TestFormCandidatesKeepsExistingFilters(t *testing.T) {
	w := &state.World{Player: 0}
	w.Factions[0].Alive = true
	w.Factions[0].Lord = 0
	w.Generals[1] = state.General{Alive: true, Faction: 0, Posted: true, Captor: 0xFF}
	w.Generals[2] = state.General{Alive: true, Faction: 1, Captor: 0xFF}
	w.Generals[3] = state.General{Alive: false, Faction: 0, Captor: 0xFF}
	w.Generals[4] = state.General{Alive: true, Faction: 0, Captor: 7}
	w.Generals[5] = state.General{Alive: true, Faction: 0, Captor: 0xFF}
	g := &game{world: w}
	rows := g.formCandidates()
	if len(rows) != 1 || rows[0] != 5 {
		t.Errorf("候選 = %v，只應該有武將 5", rows)
	}
}
