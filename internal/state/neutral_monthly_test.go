package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

func TestNeutralMonthlyDeclarationRawPayloadAndThreshold(t *testing.T) {
	for _, funds := range []int{112 << 8, 113 << 8} {
		w := &World{}
		w.Factions[0] = Faction{Alive: true, Cities: 1, Funds: funds, InvasionTarget: diplomacy.NoTarget}
		w.Cities[0] = City{Owner: 0, Adjacency: 1, Neighbours: [4]int{1, -1, -1, -1}}
		w.Cities[1].Owner = 24
		var out []StrategyEvent
		queued := w.queueNeutralDeclaration(rng.NewFixed(0), 0, &out)
		if queued != (funds > 112<<8) {
			t.Fatalf("funds %d queued %v", funds, queued)
		}
		if queued {
			found := false
			for _, e := range w.events {
				found = found || e.Code == 1 && e.Param == 0xFF18
			}
			if !found || len(out) != 1 || out[0].Target != 24 {
				t.Fatal("neutral raw declaration not written")
			}
		}
	}
}

func TestNeutralMonthlyDeclarationMissingBorderAndExistingTarget(t *testing.T) {
	w := &World{}
	w.Factions[0] = Faction{Alive: true, Cities: 1, Funds: 500000, InvasionTarget: 24}
	var out []StrategyEvent
	g := rng.NewFixed(0)
	before := append([]byte(nil), g.Raw()...)
	if w.queueNeutralDeclaration(g, 0, &out) || w.Factions[0].InvasionTarget != diplomacy.NoTarget {
		t.Fatal("missing border did not clear neutral target")
	}
	w.Cities[0] = City{Owner: 0, Adjacency: 1, Neighbours: [4]int{1, -1, -1, -1}}
	w.Cities[1].Owner = 24
	w.Factions[0].InvasionTarget = 24
	if w.queueNeutralDeclaration(g, 0, &out) || len(out) != 0 {
		t.Fatal("existing neutral target emitted duplicate")
	}
	for i, value := range before {
		if g.Raw()[i] != value {
			t.Fatal("non-producing neutral path consumed RNG")
		}
	}
}
