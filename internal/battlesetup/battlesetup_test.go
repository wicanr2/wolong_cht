package battlesetup

import (
	"os"
	"reflect"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

const original = "../../workplace/orig/dosv"

func TestProviderLoadsExplicitSelectionForEveryEncounter(t *testing.T) {
	if _, err := os.Stat(original + "/KI.EXE"); err != nil {
		t.Skip("缺少原版素材")
	}
	w, err := state.LoadScenario(original+"/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	p, s, err := Load(Options{Dir: original, World: w})
	if err != nil {
		t.Fatal(err)
	}
	if p.lib == nil {
		t.Fatal("missing battle library")
	}
	for _, tc := range []struct {
		number int
		rotate bool
	}{{193, true}, {210, false}, {82, true}, {211, false}, {193, false}} {
		tiles, gate := p.lib.Tiles(tc.number), p.lib.GateX(tc.number)
		if tc.rotate {
			tiles, gate = battle.Rotate180(tiles), battle.RotateGateX(gate)
		}
		want := tactical.NewFieldFromTileLayers(tiles, p.lib.Heights(tc.number), p.lib.TileLayers(tc.number), gate)
		got := s.Field(tc.number, tc.number < 192, tc.rotate)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("wrong geometry for explicit field %d rotate %v", tc.number, tc.rotate)
		}
		if got := s.Category(tc.number, tc.number < 192); got != battle.Category(tc.number) {
			t.Fatalf("category %d", got)
		}
		if !reflect.DeepEqual(s.Script(tc.number, false, 2), p.lib.Script(2, battle.Category(tc.number))) {
			t.Fatalf("script field %d", tc.number)
		}
	}
}

func TestMissingFormationsStillAttachTerrainForAutomaticBattles(t *testing.T) {
	if _, err := os.Stat(original + "/SINARIO.DAT"); err != nil {
		t.Skip("缺少原版素材")
	}
	w, err := state.LoadScenario(original+"/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
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
	if _, _, err := Load(Options{Dir: t.TempDir(), World: w, Map: m}); err == nil {
		t.Fatal("missing forms should report error")
	}
	// Supply a simple field afterwards. StageBattle selects from the terrain
	// attached even though asset loading failed, without touching city coords.
	w.SetTactical(&state.TacticalSetup{Forms: tactical.SyntheticFormations(), Field: func(int, bool, bool) *tactical.Field { return SyntheticField(false) }})
	if err := w.StageBattle(a, d, combat.Field, rng.NewFixed(3)); err != nil {
		t.Fatal(err)
	}
	if n := w.PendingBattle().FieldNumber; n < 209 || n > 212 {
		t.Fatalf("terrain was lost after asset error: %d", n)
	}
}
