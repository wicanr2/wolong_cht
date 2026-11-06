package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestDesktopLoadReattachesBattlefieldToRestoredWorld(t *testing.T) {
	if _, err := os.Stat(nativeScenario); err != nil {
		t.Skip("缺少原版素材")
	}
	w, err := state.LoadScenario(nativeScenario, 0)
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
	g := &game{world: w, lib: &library.Library{World: m}, origDir: filepath.Dir(nativeScenario), sourceFile: nativeScenario, saveBase: nativeScenario, saveFile: filepath.Join(t.TempDir(), "SAVE.DAT")}
	g.installTactical(g.origDir)
	if err := g.writeSave(0); err != nil {
		t.Fatal(err)
	}
	w.Corps[d].X, w.Corps[d].Y = 20, 30
	if err := g.readSave(0); err != nil {
		t.Fatal(err)
	}
	if g.world == w {
		t.Fatal("load did not replace world")
	}
	if err := g.world.StageBattle(a, d, combat.Field, rng.NewFixed(7)); err != nil {
		t.Fatal(err)
	}
	p := g.world.PendingBattle()
	if p.FieldNumber < 209 || p.FieldNumber > 212 {
		t.Fatalf("restored world lost terrain or used old position: %d", p.FieldNumber)
	}
	if g.battleFieldName(p) != "海上" {
		t.Fatal("field name did not use pending selection")
	}
}
