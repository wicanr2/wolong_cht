package savefile

import (
	"errors"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestNativeSaveRejectsUnfinishedBattleTick(t *testing.T) {
	w := load(t)
	factions := w.AliveFactions()
	a, b := factions[0], factions[1]
	w.Player = a
	w.SetTactical(&state.TacticalSetup{Forms: tactical.SyntheticFormations(),
		Field: func(int, bool, bool) *tactical.Field {
			stack := make([][]int, tactical.Height)
			for y := range stack {
				stack[y] = make([]int, tactical.Width)
			}
			return tactical.NewField(stack, 0)
		}})
	for _, f := range []int{a, b} {
		w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
		if err := w.FormCorps(w.Factions[f].Lord, [army.Positions]army.TroopType{}, [army.Positions]bool{true, true, true, true, true, true}); err != nil {
			t.Fatal(err)
		}
	}
	w.Friendship[a][b] = w.Friendship[a][b].WithWar(true)
	w.Friendship[b][a] = w.Friendship[b][a].WithWar(true)
	if err := w.March(w.Factions[b].Lord, w.Factions[a].Capital); err != nil {
		t.Fatal(err)
	}
	random := rng.NewFixed(5)
	for n := 0; n < 200000 && w.PendingBattle() == nil; n++ {
		w.Tick(random)
	}
	p := w.PendingBattle()
	if p == nil || !w.HasIncompleteTick() {
		t.Fatal("fixture never entered normal tactical battle")
	}
	if _, err := Encode(w, Origin{}); !errors.Is(err, state.ErrIncompleteTick) {
		t.Fatalf("pending encode error=%v", err)
	}
	result := tactical.NewBattle(nil, nil, rng.NewFixed(2), 0)
	result.Done, result.Winner = true, 1
	result.Sides[0].Reserve[0] = w.Corps[p.Attacker].Men / 2
	result.Sides[1].Reserve[0] = w.Corps[p.Defender].Men / 2
	p.Battle = result
	if w.ResolvePending(random) == nil {
		t.Fatal("missing result")
	}
	if _, err := Encode(w, Origin{}); !errors.Is(err, state.ErrIncompleteTick) {
		t.Fatalf("result-to-resume gap encode error=%v", err)
	}
	w.ResumeTick(random)
	data, err := Encode(w, Origin{})
	if err != nil {
		t.Fatal(err)
	}
	back, _, err := Decode(data)
	if err != nil || back.Clock != w.Clock {
		t.Fatal("completed native save lost resumed time")
	}
}
