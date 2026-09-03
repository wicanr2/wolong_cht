package state

import "testing"

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
