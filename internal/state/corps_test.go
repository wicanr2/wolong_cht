package state

import (
	"bytes"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

func TestCapturedCorpsLeaderCannotBeSelectedForRelief(t *testing.T) {
	for _, sovereign := range []bool{false, true} {
		w := load(t, 0)
		const id, loser, winner = 7, 12, 13
		g := &w.Generals[id]
		g.Alive, g.Faction, g.Captor, g.Duty = true, loser, noFaction, DutyCorpsLeader
		g.Martial, g.Command, g.Politics = 15, 0, 0
		g.Sovereign, g.TalkVariant = sovereign, 1
		w.Factions[loser].Lord, w.Factions[loser].Capital = 3, 82
		w.Corps[id] = Corps{Alive: true, Faction: loser}
		ev := &CorpsEvent{Corps: id}
		w.corpsPerishes(ev, id, winner, &fixedRand{v: 127})
		if ev.Fate[id] != combat.Captured || g.Duty != DutyCaptive || g.Captor != loser || g.Faction != winner {
			t.Fatalf("戰敗被俘狀態錯誤：fate=%v general=%+v", ev.Fate[id], g)
		}
		variant := 1
		if sovereign {
			variant += 3
		}
		if g.Sovereign || g.TalkVariant != variant {
			t.Fatal("被俘主公型沒有套用既有轉換")
		}
		for i := range w.Generals {
			if i != id {
				w.Generals[i].Martial = 0
			}
		}
		w.Factions[winner].Alive, w.Factions[winner].Corps = true, 0
		w.Factions[winner].Reserves = [3]int{900, 900, 900}
		w.Cities[97].Owner = winner
		w.requestRelief(97, 1, rng.NewFixed(7))
		if w.Corps[id].Alive || g.Duty != DutyCaptive {
			t.Fatal("AI 把俘虜重新編成軍團")
		}
		saved := w.Bytes()
		restored, err := LoadBlock(saved)
		if err != nil || !bytes.Equal(saved, restored.Bytes()) || restored.Generals[id].Duty != DutyCaptive {
			t.Fatal("被俘狀態存讀檔改變", err)
		}
	}
}

// 應戰軍團照 sub_14C72 計分挑（docs/spec/82）：
// 分數＝(兵數>>4 低 byte)×(士氣>>4)×(評價>>4＋1)，嚴格最大、同分取先。
func TestDefenderPickedByScore(t *testing.T) {
	w := load(t, 0)
	// 造兩支同勢力守軍疊在同一個據點：0 號兵多但士氣低、1 號兵少士氣高。
	f := (w.Player + 1) % 2
	node := w.Factions[f].Capital
	for j := 0; j < 2; j++ {
		w.Corps[j] = Corps{Alive: true, Faction: f, Node: node,
			X: w.Cities[node].X, Y: w.Cities[node].Y}
	}
	w.Corps[0].Men, w.Corps[0].Morale = 600, 60
	w.Corps[1].Men, w.Corps[1].Morale = 300, 240
	// 分數：0 → (600>>4)&0xFF=37 × 3 = 111×(r+1)；1 → 18 × 15 = 270×(r+1)。
	got := w.pickDefender(-1, f, func(d *Corps) bool { return d.Node == node })
	if got != 1 {
		t.Fatalf("挑到軍團 %d，兵少士氣高的 1 分數較高應出戰", got)
	}
	// 同分取先者。
	w.Corps[1] = w.Corps[0]
	if got := w.pickDefender(-1, f, func(d *Corps) bool { return d.Node == node }); got != 0 {
		t.Fatalf("同分挑到 %d，應取先者 0", got)
	}
}

// 下場判定要把當下的兩個勢力記進事件（docs/spec/123 §3）。
//
// ⭐ **被擒之後 `Generals[i].Faction` 已經是勝方**，事後反推拿不到舊主——
// 這條測試擋的就是「事件不記、UI 自己去猜」那條路。
func TestCorpsPerishesRecordsFateSides(t *testing.T) {
	w := &World{}
	w.Factions[0].Alive = true
	w.Factions[0].Capital = 5
	w.Factions[0].Corps = 1
	w.Factions[1].Alive = true
	w.Corps[7].Alive = true
	w.Corps[7].Faction = 0
	w.Generals[7].Alive = true
	w.Generals[7].Faction = 0

	ev := &CorpsEvent{Corps: 7}
	w.corpsPerishes(ev, 7, 1, &testRand{})

	side, ok := ev.FateSides[7]
	if !ok {
		t.Fatal("沒有記下勝敗兩方")
	}
	if side.Winner != 1 || side.Loser != 0 {
		t.Errorf("勝方 %d／敗方 %d，want 1／0", side.Winner, side.Loser)
	}
	if _, ok := ev.Fate[7]; !ok {
		t.Error("Fate 沒記")
	}
}
