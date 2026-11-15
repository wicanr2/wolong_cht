package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// 原版196/6/11 23時子刻5：軍團64在據點88，Stage10改回首都82。
func TestHomeReplanVisitAndSave(t *testing.T) {
	for _, stage := range []int{StageHomeResupply, StageDisband} {
		w := load(t, 0)
		const id, faction, here, home = 64, 12, 88, 82
		w.Player = 2
		w.Factions[faction].Capital = home
		c := w.newCorpsRecord(id, faction, here, 82)
		c.Stage, c.Timer, c.Ordered, c.OnPath = stage, 1, here, true
		w.Corps[id] = c
		w.tickOneCorps(id, 23, rng.NewFixed(7))
		got := w.Corps[id]
		if !got.Replan || got.Ordered != home || got.TargetNode != home || got.Stage != stage || got.X != c.X || got.Y != c.Y {
			t.Fatalf("Stage%d 到站改意圖不同：%+v", stage, got)
		}
		loaded, err := LoadBlock(w.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if !loaded.Corps[id].Replan || loaded.Corps[id].Ordered != home {
			t.Fatal("存讀檔遺失重算旗標或意圖")
		}
	}
}

func TestRetargetKeepsLegAndReplanDecision(t *testing.T) {
	w := load(t, 0)
	const id, home = 64, 82
	c := w.newCorpsRecord(id, 12, home, 82)
	c.LinkAddr, c.TargetNode, c.Ordered = 0x0ec0, 88, home
	c.Direction, c.PathPtr, c.Replan = 252, 0x40cc, false
	w.Corps[id] = c
	if w.retarget(id, home) {
		t.Fatal("連結位址不能當成已抵達首都")
	}
	got := w.Corps[id]
	if got.Replan || got.Direction != c.Direction || got.PathPtr != c.PathPtr || got.LinkAddr != c.LinkAddr || got.Ordered != home || got.TargetNode != home {
		t.Fatalf("單純寫目標不應重算路段：%+v", got)
	}
}
