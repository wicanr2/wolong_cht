package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

type fieldRandom struct{ calls, value int }

func (r *fieldRandom) Next() int { r.calls++; return r.value }

func fieldTestSetup() *TacticalSetup {
	return &TacticalSetup{Forms: tactical.SyntheticFormations(), Field: func(int, bool, bool) *tactical.Field {
		rows := make([][]int, tactical.Height)
		for y := range rows {
			rows[y] = make([]int, tactical.Width)
		}
		return tactical.NewField(rows, 0)
	}}
}

// Coordinates differ from the city and attacker. These expectations come from
// original 14B63/14BDD and its pair table, not from another call to SelectWith.
func TestFieldSelectionUsesContactAndParticipatingPlayer(t *testing.T) {
	w := &World{Player: 2}
	w.Cities[8].X, w.Cities[8].Y = 20, 30
	w.Corps[12] = Corps{Faction: 2, Node: 8, X: 60, Y: 70}
	w.Corps[13] = Corps{Faction: 3, X: 100, Y: 110}
	w.Corps[0] = Corps{Alive: true, Faction: 2, Heading: 0} // unrelated moving player
	tiles := map[[2]int]byte{{100, 109}: 0x70, {99, 110}: 0x0E, {101, 110}: 0x0E}
	w.SetTerrain(func(x, y int) byte { return tiles[[2]int{x, y}] })
	for _, tc := range []struct {
		dir, number int
		rotate      bool
	}{
		{0, 193, true}, {1, 193, false}, {2, 205, false}, {3, 205, false}, {4, 205, false},
	} {
		w.Corps[12].Heading = tc.dir
		r := &fieldRandom{value: 2}
		got := w.selectBattlefield(12, 13, r)
		if got.number != tc.number || got.rotate != tc.rotate || r.calls != 0 {
			t.Fatalf("direction %d: %+v calls=%d", tc.dir, got, r.calls)
		}
	}
	// A participating player defender overrides the attacker, including still=4.
	w.Corps[13].Faction, w.Corps[13].Heading = 2, 4
	w.Corps[12].Heading = 1
	if got := w.selectBattlefield(12, 13, &fieldRandom{}); got.number != 205 || got.rotate {
		t.Fatalf("stationary player defender: %+v", got)
	}
	// Nonplain centre must use the actual contact cell, not the city or attacker.
	tiles[[2]int{100, 110}] = 0x70
	if got := w.selectBattlefield(12, 13, &fieldRandom{}); got.number != 209 || got.rotate {
		t.Fatalf("mountain contact: %+v", got)
	}
}

func TestFieldSelectionPrecedesRepresentativeDefender(t *testing.T) {
	w, att, def := standoffFixture(t)
	w.Player = w.Corps[def].Faction
	w.Corps[def].X, w.Corps[def].Y, w.Corps[def].Heading, w.Corps[def].Men = 100, 110, 1, 16
	representative := -1
	for i := range w.Corps {
		if i != att && i != def && !w.Corps[i].Alive {
			representative = i
			break
		}
	}
	if representative < 0 {
		t.Fatal("no free corps slot")
	}
	w.Corps[representative], w.Generals[representative] = w.Corps[def], w.Generals[def]
	w.Corps[representative].Heading, w.Corps[representative].Men = 0, 1000
	w.SetTerrain(func(x, y int) byte {
		if x == 100 && y == 109 {
			return 0x70
		}
		return 0
	})
	w.SetTactical(fieldTestSetup())
	r := &fieldRandom{value: 2}
	w.rng = r
	w.fieldAt(att, def, &CorpsEvent{}, r)
	p := w.PendingBattle()
	if p == nil || p.Defender != representative || p.FieldNumber != 193 || p.FieldRotate {
		t.Fatalf("selection after representative switch: %+v", p)
	}
}

func TestWaterSelectionSharedAcrossCallbacksAndNextEncounter(t *testing.T) {
	w, att, def := standoffFixture(t)
	w.Player = w.Corps[att].Faction
	w.Corps[def].X, w.Corps[def].Y = 100, 110
	w.SetTerrain(func(x, y int) byte {
		if x == 100 && y == 110 {
			return 0xCA
		}
		return 0
	})
	for _, roll := range []int{1, 2} {
		r := &fieldRandom{value: roll}
		w.rng = r
		setup := fieldTestSetup()
		makeField := setup.Field
		categoryCalls, scriptCalls := 0, 0
		setup.Field = func(field int, siege, rotate bool) *tactical.Field {
			if field != 209+roll || siege || rotate || r.calls != 1 {
				t.Fatalf("field=%d rotate=%v calls=%d", field, rotate, r.calls)
			}
			return makeField(field, siege, rotate)
		}
		setup.Category = func(field int, siege bool) int {
			categoryCalls++
			if field != 209+roll || siege {
				t.Fatalf("category used %d", field)
			}
			return 2
		}
		setup.Script = func(field int, siege bool, tactic int) []byte {
			scriptCalls++
			if field != 209+roll || siege {
				t.Fatalf("script used %d", field)
			}
			return nil
		}
		w.SetTactical(setup)
		w.fieldAt(att, def, &CorpsEvent{}, r)
		p := w.PendingBattle()
		if p == nil || p.FieldNumber != 209+roll || p.Node != w.Corps[att].Node || categoryCalls != 2 || scriptCalls != 1 {
			t.Fatalf("shared selection: %+v category=%d script=%d", p, categoryCalls, scriptCalls)
		}
		w.pending = nil
	}
}

func TestWaterSelectionOnceInAutomaticAndFieldFailure(t *testing.T) {
	for _, fieldFailure := range []bool{false, true} {
		counts := [2]int{}
		for variant, tile := range []byte{0, 0xCA} {
			w, att, def := standoffFixture(t)
			w.Player = w.Corps[att].Faction
			w.SetTerrain(func(int, int) byte { return tile })
			if fieldFailure {
				s := fieldTestSetup()
				s.Field = func(int, bool, bool) *tactical.Field { return nil }
				w.SetTactical(s)
			}
			r := &fieldRandom{value: 2}
			w.rng = r
			ev := &CorpsEvent{}
			w.fieldAt(att, def, ev, r)
			if ev.Battle == nil || w.PendingBattle() != nil {
				t.Fatal("automatic fallback did not resolve")
			}
			counts[variant] = r.calls
		}
		if counts[1] != counts[0]+1 {
			t.Fatalf("field failure=%v, RNG calls=%v", fieldFailure, counts)
		}
	}
}

func TestStageBattleSelectionAndSiegeRotation(t *testing.T) {
	w, att, def := standoffFixture(t)
	w.Player = w.Corps[def].Faction
	w.SetTactical(fieldTestSetup())
	w.SetTerrain(func(int, int) byte { return 0xCA })
	r := &fieldRandom{value: 3}
	if err := w.StageBattle(att, def, combat.Field, r); err != nil {
		t.Fatal(err)
	}
	if p := w.PendingBattle(); p.FieldNumber != 212 || p.FieldRotate {
		t.Fatalf("field stage: %+v", p)
	}
	w.pending = nil
	if err := w.StageBattle(att, def, combat.Siege, r); err != nil {
		t.Fatal(err)
	}
	if p := w.PendingBattle(); p.FieldNumber != p.Node || !p.FieldRotate {
		t.Fatalf("defending siege: %+v", p)
	}
}
