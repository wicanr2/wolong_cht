package state

import (
	"bytes"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"testing"
)

// Native normal trace: corps56 disbands at one o'clock, retains aggregate303,
// pays10, and morale82 becomes92. Reusing its record must retain old path bytes.
func TestDisbandVisitAndReformationPreserveOriginalRecord(t *testing.T) {
	w := load(t, 0)
	const id, faction, home = 56, 12, 82
	f := &w.Factions[faction]
	f.Alive, f.Capital, f.Corps, f.Funds = true, home, 1, 1000
	f.Reserves = [3]int{}
	f.MoraleBase = 200
	g := &w.Generals[id]
	g.Alive, g.Faction, g.Duty, g.Captor = true, faction, DutyCorpsLeader, 255
	w.Cities[home].Owner = faction
	c := &w.Corps[id]
	*c = Corps{Alive: true, Faction: faction, Men: 303, Morale: 82, Node: home, TargetNode: home, Ordered: home,
		X: w.Cities[home].X, Y: w.Cities[home].Y, TargetX: w.Cities[home].X, TargetY: w.Cities[home].Y,
		Stage: StageDisband, Direction: 4, PathPtr: 0x57fc, Timer: 1, Interval: 3, Heading: HeadingStill}
	for k, kind := range []army.TroopType{0, 0, 2, 2, 1, 1} {
		men := 50
		if k == 5 {
			men = 53
		}
		c.Units[k] = combat.Unit{Kind: kind, Men: men}
	}
	ev := w.tickOneCorps(id, 1, rng.NewFixed(7))
	if ev == nil || !ev.Disbanded || c.Alive || f.Funds != 990 || c.Men != 303 || c.Morale != 92 || c.Stage != StageDisband {
		t.Fatalf("解體尾段不同：事件%+v 資金%d 軍團%+v", ev, f.Funds, *c)
	}
	b := w.Bytes()
	off := corpsBase + id*corpsSize
	if b[off] != 0 || u16(b, off+4) != 303 || b[off+6] != 92 || b[off+0x23] != StageDisband {
		t.Fatal("解體欄位未保存")
	}
	loaded, e := LoadBlock(b)
	if e != nil || !bytes.Equal(b, loaded.Bytes()) {
		t.Fatal("解體存讀檔", e)
	}
	kinds := [6]army.TroopType{0, 0, 2, 2, 1, 1}
	manned := [6]bool{true, true, true, true, true, true}
	if e := loaded.FormCorps(id, kinds, manned); e != nil {
		t.Fatal(e)
	}
	reformed := loaded.Corps[id]
	if reformed.Direction != 4 || reformed.PathPtr != 0x57fc || reformed.Stage != 1 || reformed.Timer != 1 {
		t.Fatalf("重編成遺失原欄位：%+v", reformed)
	}
	saved := loaded.Bytes()
	again, e := LoadBlock(saved)
	if e != nil || !bytes.Equal(saved, again.Bytes()) {
		t.Fatal("重編成存讀檔", e)
	}
}
