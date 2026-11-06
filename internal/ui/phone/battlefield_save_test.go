package phone

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
)

func TestPhoneLoadReattachesTacticalAndTerrain(t *testing.T) {
	s := newSessionInTempRoot(t)
	w := s.world
	for _, f := range []int{0, 12} {
		w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
		if err := w.FormCorps(w.Factions[f].Lord, [army.Positions]army.TroopType{}, [army.Positions]bool{true, true, true, true, true, true}); err != nil {
			t.Fatal(err)
		}
	}
	a, d := w.Factions[0].Lord, w.Factions[12].Lord
	w.Corps[d].X, w.Corps[d].Y = 100, 110
	m := &world.Map{Tiles: make([]byte, world.Width*world.Height)}
	m.Tiles[110*world.Width+100] = 0xCA
	s.lib.World = m
	s.attachBattlefield()
	if err := s.SaveSlot(0); err != nil {
		t.Fatal(err)
	}
	w.Corps[d].X, w.Corps[d].Y = 20, 30
	if err := s.LoadSlot(0); err != nil {
		t.Fatal(err)
	}
	if s.world == w {
		t.Fatal("load did not replace world")
	}
	if err := s.world.StageBattle(a, d, combat.Field, s.rand); err != nil {
		t.Fatal(err)
	}
	p := s.world.PendingBattle()
	if p.FieldNumber < 209 || p.FieldNumber > 212 {
		t.Fatalf("loaded field %d", p.FieldNumber)
	}
	if b := s.MusicScene().Battle; b == nil || b.Field != p.FieldNumber {
		t.Fatalf("music did not use selected field: %+v", b)
	}
}
