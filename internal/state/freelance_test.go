package state

import "testing"

// 挑一個在場的在野武將當 fixture，並把倒數與心向的勢力設成指定值。
func freelanceFixture(t *testing.T, affinity int) (*World, int) {
	t.Helper()
	w := load(t, 0)
	for id, g := range w.Generals {
		if !g.Alive || g.Faction != noFaction {
			continue
		}
		g.Timer = 0
		g.Affinity = affinity
		g.VanishIfAffinityGone = false
		w.Generals[id] = g
		return w, id
	}
	t.Fatal("劇本一沒有在場的在野武將，fixture 不成立")
	return nil, -1
}

// 開局在場的在野武將**全部**有心向的勢力，這是 docs/re/77 §2.4 的判準：
// 一個例外都沒有。這一條同時證明「隨機投靠」在開局那一輪輪不到。
func TestEveryFreelanceGeneralHasAffinity(t *testing.T) {
	for idx := 0; idx < 4; idx++ {
		w := load(t, idx)
		n, missing := 0, 0
		for _, g := range w.Generals {
			if !g.Alive || g.Faction != noFaction {
				continue
			}
			n++
			if g.Affinity == noFaction {
				missing++
			}
		}
		// 正對照：沒有在野武將的話上面那個檢查是空的，會無聲通過。
		if n == 0 {
			t.Errorf("劇本 %d 一個在野武將都沒有，這一輪什麼都沒驗到", idx+1)
		}
		if missing != 0 {
			t.Errorf("劇本 %d：%d 名在野武將裡有 %d 名沒有心向的勢力", idx+1, n, missing)
		}
	}
}

// 亂數落在 0x40 以下就兌現：出仕過去，欄位清成 0xFF，勢力的武將數 +1。
func TestFreelanceJoinsAffinityFaction(t *testing.T) {
	w, id := freelanceFixture(t, 1)
	if !w.Factions[1].Alive {
		t.Skip("劇本一的勢力 1 不在場，換一個 fixture 才驗得了")
	}
	before := w.Factions[1].Generals
	w.recruitFreelanceGenerals(&sequenceRand{values: []int{0x3F}})

	g := w.Generals[id]
	if g.Faction != 1 {
		t.Errorf("Faction = %d，want 1", g.Faction)
	}
	if g.Affinity != noFaction {
		t.Errorf("Affinity = %d，兌現之後應該清成 0xFF", g.Affinity)
	}
	if got := w.Factions[1].Generals; got != before+1 {
		t.Errorf("勢力 1 的武將數 = %d，want %d", got, before+1)
	}
}

// 亂數 ≥ 0x40（75%）什麼都不動——連欄位都不清。
func TestFreelanceHoldsOnHighRoll(t *testing.T) {
	w, id := freelanceFixture(t, 1)
	w.recruitFreelanceGenerals(&sequenceRand{values: []int{0x40}})
	if g := w.Generals[id]; g.Faction != noFaction || g.Affinity != 1 {
		t.Errorf("0x40 不該兌現：%+v", g)
	}
}

// 心向的勢力已經滅了：旗標 bit 5 設著的整筆歸零，沒設的留在原地。
func TestFreelanceVanishesOnlyWithFlag(t *testing.T) {
	// ⚠ 不要去「找一個已滅的勢力」——劇本一開局一個都沒有，
	// 那樣寫這一條會整個跳過而測試報告是綠的。自己弄倒一個。
	const dead = numFactions - 1

	for _, flag := range []bool{false, true} {
		w, id := freelanceFixture(t, dead)
		w.Factions[dead].Alive = false
		g := w.Generals[id]
		g.VanishIfAffinityGone = flag
		w.Generals[id] = g

		w.recruitFreelanceGenerals(&sequenceRand{values: []int{0x00}})

		got := w.Generals[id]
		if got.Alive == flag {
			t.Errorf("bit 5 = %v 時 Alive = %v，want %v", flag, got.Alive, !flag)
		}
		if got.Faction != noFaction {
			t.Errorf("勢力已滅不該出仕：Faction = %d", got.Faction)
		}
		if got.Affinity != noFaction {
			t.Errorf("兌現之後 Affinity 應該清成 0xFF，得到 %d", got.Affinity)
		}
	}
}

// 倒數沒歸零的那個月只遞減，不做別的（sub_1585F 的第二個閘）。
func TestFreelanceTimerGate(t *testing.T) {
	w, id := freelanceFixture(t, 1)
	g := w.Generals[id]
	g.Timer = 2
	w.Generals[id] = g

	w.recruitFreelanceGenerals(&sequenceRand{values: []int{0x00}})
	if got := w.Generals[id]; got.Timer != 1 || got.Faction != noFaction {
		t.Fatalf("倒數 2 → 1 且不出仕，得到 %+v", got)
	}
	w.recruitFreelanceGenerals(&sequenceRand{values: []int{0x00}})
	if got := w.Generals[id]; got.Timer != 0 || got.Faction != noFaction {
		t.Fatalf("倒數 1 → 0 且仍不出仕，得到 %+v", got)
	}
	w.recruitFreelanceGenerals(&sequenceRand{values: []int{0x00}})
	if got := w.Generals[id]; got.Faction != 1 {
		t.Fatalf("倒數歸零之後該出仕，得到 %+v", got)
	}
}

// +0x19 與旗標的四個位元要 byte-for-byte 寫得回去（CLAUDE.md §9）。
func TestGeneralAffinityAndFlagsRoundTrip(t *testing.T) {
	for idx := 0; idx < 4; idx++ {
		w := load(t, idx)
		out := w.Bytes()
		src := w.raw
		for i := range w.Generals {
			off := generalBase + i*generalSize
			if out[off+0x19] != src[off+0x19] {
				t.Fatalf("劇本 %d 武將 %d 的 +0x19：%#02x != %#02x",
					idx+1, i, out[off+0x19], src[off+0x19])
			}
			// bit 0 沒解，一起比整個 byte 才擋得住「改寫變重建」。
			if out[off] != src[off] {
				t.Fatalf("劇本 %d 武將 %d 的旗標：%#02x != %#02x",
					idx+1, i, out[off], src[off])
			}
		}
	}
}
