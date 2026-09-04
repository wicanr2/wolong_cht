package state

import "testing"

// fixedRand 讓每一次 Next() 回固定值，用來釘住骰面對應的分支。
type fixedRand struct{ v int }

func (f fixedRand) Next() int { return f.v }

// 骰面 1 與 24–47 都投靠**武將最少**的勢力（docs/spec/130 §2.1、§2.3）。
func TestRandomJoinPrefersFewestGenerals(t *testing.T) {
	for _, roll := range []int{0, 0x17, 0x2E} { // +1 之後 ＝ 1、24、47
		w := &World{Player: 5}
		for i := 0; i < numFactions; i++ {
			w.Factions[i].Alive = true
			w.Factions[i].Generals = 10
		}
		w.Factions[7].Generals = 2 // 最少
		w.Generals[3] = General{Alive: true, Faction: noFaction, Affinity: noFaction}

		got := w.randomJoin(3, fixedRand{roll})
		if got != 7 {
			t.Errorf("骰面 %d → 勢力 %d，want 7（武將最少）", roll+1, got)
		}
		if w.Generals[3].Faction != 7 || w.Factions[7].Generals != 3 {
			t.Errorf("沒有寫回：Faction=%d、武將數=%d",
				w.Generals[3].Faction, w.Factions[7].Generals)
		}
	}
}

// 骰面 n（1–23）從最少那一個往後數第 n 個**存在**的勢力，會繞回。
func TestRandomJoinCountsFromFewest(t *testing.T) {
	w := &World{Player: 21}
	for i := 0; i < numFactions; i++ {
		w.Factions[i].Alive = false
	}
	// 只有 3、9、20 活著；9 的武將最少 ⇒ 起點是 9。
	for _, i := range []int{3, 9, 20} {
		w.Factions[i].Alive = true
		w.Factions[i].Generals = 10
	}
	w.Factions[9].Generals = 1
	w.Generals[4] = General{Alive: true, Faction: noFaction, Affinity: noFaction}

	// 骰面 3 ⇒ 從 9 起數三個存在的：9(1)、20(2)、3(3，繞回) ⇒ 勢力 3。
	if got := w.randomJoin(4, fixedRand{2}); got != 3 {
		t.Errorf("骰面 3 → 勢力 %d，want 3（繞回）", got)
	}
}

// 骰面 ≥ 48 走玩家救濟：據點數 ÷ 4 ＋ 1 > 武將數才送人。
func TestRandomJoinPlayerReliefThreshold(t *testing.T) {
	mk := func(cities, generals int) *World {
		w := &World{Player: 2}
		for i := 0; i < numFactions; i++ {
			w.Factions[i].Alive = true
			w.Factions[i].Generals = 10
		}
		w.Factions[2].Cities, w.Factions[2].Generals = cities, generals
		w.Generals[6] = General{Alive: true, Faction: noFaction, Affinity: noFaction}
		return w
	}
	// 20/4+1 = 6 > 3 ⇒ 送。
	if got := mk(20, 3).randomJoin(6, fixedRand{0x2F}); got != 2 {
		t.Errorf("缺人時 → %d，want 2（玩家）", got)
	}
	// 20/4+1 = 6 ≤ 6 ⇒ 不動。
	if got := mk(20, 6).randomJoin(6, fixedRand{0x2F}); got != -1 {
		t.Errorf("不缺人時 → %d，want −1（什麼都不做）", got)
	}
}

// ⭐ 接線：沒有心向的武將**每月都跑**，不受 25% 閘限制
// （原版 `cmp bh, 0FFh / jz` 在擲那顆骰之前，docs/spec/130）。
func TestRandomJoinSkipsTwentyFivePercentGate(t *testing.T) {
	w := &World{Player: 5}
	for i := 0; i < numFactions; i++ {
		w.Factions[i].Alive = true
		w.Factions[i].Generals = 10
	}
	// 玩家缺人：40/4+1 = 11 > 10 ⇒ 骰面 64 那一條會送人過來。
	w.Factions[5].Cities = 40
	w.Generals[3] = General{Alive: true, Faction: noFaction, Affinity: noFaction}

	// ⭐ 0xFF 遠大於 affinityRollLimit（0x40）——**有心向的話這一輪不會動**，
	// 而這一位沒有心向，所以照跑。骰面 (0xFF&0x3F)+1 = 64 ⇒ 玩家救濟。
	joined := w.recruitFreelanceGenerals(fixedRand{0xFF})
	if len(joined) != 1 || joined[0] != 3 {
		t.Fatalf("出仕清單 %v，want [3]——隨機投靠不該被 25%% 閘擋住", joined)
	}
	if w.Generals[3].Faction != 5 {
		t.Errorf("投靠的是勢力 %d，want 5（玩家救濟）", w.Generals[3].Faction)
	}

	// 正對照：同一顆骰、同一個人，**有心向**就會被 25% 閘擋下來。
	w2 := &World{Player: 5}
	for i := 0; i < numFactions; i++ {
		w2.Factions[i].Alive = true
		w2.Factions[i].Generals = 10
	}
	w2.Generals[3] = General{Alive: true, Faction: noFaction, Affinity: 7}
	if got := w2.recruitFreelanceGenerals(fixedRand{0xFF}); len(got) != 0 {
		t.Errorf("有心向的在 0xFF 這一輪出仕了 %v，25%% 閘沒作用", got)
	}
}
