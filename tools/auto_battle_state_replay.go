//go:build ignore

// 原版自動判定入口資料的規則層重播。
package main

import (
	"encoding/json"
	"flag"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/state"
	"os"
	"path/filepath"
)

func main() {
	original := flag.String("original", "", "原版收據")
	out := flag.String("out", "", "輸出")
	flag.Parse()
	must := func(e error) {
		if e != nil {
			panic(e)
		}
	}
	w, e := state.LoadScenario(filepath.Join(*original, "entry-SAVE.DAT"), 0)
	must(e)
	raw, e := os.ReadFile(filepath.Join(*original, "rng.bin"))
	must(e)
	r, ok := rng.FromRaw(raw)
	if !ok {
		panic("亂數格式")
	}
	view := func(i int) combat.Corps {
		c := w.Corps[i]
		g := w.Generals[i]
		return combat.Corps{Faction: c.Faction, Leader: combat.Leader{Martial: g.Martial, Command: g.Command, SiegeAptitude: g.Aptitude[0], FieldAptitude: g.Aptitude[1], Rating: g.Rules().Rating()}, Units: c.Units, Morale: c.Morale, Men: c.Men}
	}
	a, d := view(35), view(39)
	pa, pd := combat.Power(a, combat.Field, true, 0, r), combat.Power(d, combat.Field, false, 0, r)
	r, _ = rng.FromRaw(raw)
	before := [2]combat.Corps{a, d}
	result := combat.Resolve(&a, &d, combat.Field, 0, r)
	c, s := r.State()
	b, e := json.MarshalIndent(map[string]any{"input": before, "power": [2]int{pa, pd}, "result": result, "after": [2]combat.Corps{a, d}, "rng_prefix": [2]uint8{c, s}, "delegated": [2]bool{w.Corps[35].Delegated, w.Corps[39].Delegated}, "classification": "共用自動判定規則；不是戰術戰鬥"}, "", "  ")
	must(e)
	must(os.WriteFile(*out, b, 0644))
}
