package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// 戰術戰鬥打完，**兩側**的士氣都按兵力比縮（docs/spec/129、`sub_19F58`）。
func TestPostBattleMoraleScalesWithMen(t *testing.T) {
	// 勝方：戰前士氣 × 新 ÷ 舊。
	if got, want := postBattleMorale(200, 100, 60, true), 120; got != want {
		t.Errorf("勝方 = %d，want %d", got, want)
	}
	// 敗方：base 99（不是戰前值，也不是 100）× 新 ÷ 舊。
	if got, want := postBattleMorale(200, 100, 60, false), 59; got != want {
		t.Errorf("敗方 = %d，want %d", got, want)
	}
}

// 戰前士氣 < 100 的敗方歸零（`sub_19EBD` 的 `cmp al, 63h / ja`）。
func TestPostBattleLoserBelowGateZeroes(t *testing.T) {
	if got := postBattleMorale(99, 100, 60, false); got != 0 {
		t.Errorf("戰前 99 的敗方 = %d，want 0", got)
	}
	if got := postBattleMorale(100, 100, 60, false); got == 0 {
		t.Error("戰前 100 的敗方不該歸零")
	}
}

// ⭐ 戰前士氣 < 100 的**勝方不歸零**——這是戰術層與自動判定的差別
// （`sub_1474A` 對勝方有 `cmp [si+6], 64h / jb ⇒ 歸零`，`sub_19F58` 沒有）。
func TestPostBattleWinnerBelowGateSurvives(t *testing.T) {
	if got, want := postBattleMorale(80, 100, 50, true), 40; got != want {
		t.Errorf("戰前 80 的勝方 = %d，want %d", got, want)
	}
}

// 兵力歸零 ⇒ 士氣 0（原版 `and cx,cx / jz`）。
func TestPostBattleMoraleZeroWhenNoMen(t *testing.T) {
	if got := postBattleMorale(200, 100, 0, true); got != 0 {
		t.Errorf("新兵力 0 ⇒ %d，want 0", got)
	}
	if got := postBattleMorale(200, 0, 60, true); got != 0 {
		t.Errorf("舊兵力 0 ⇒ %d，want 0", got)
	}
}

// ⭐ 接線本身：`ResolvePending` 要真的用 postBattleMorale。
//
// ⚠ 這一支存在的理由是**上面四支測不到接線**——它們直接呼叫那個函式，
// 所以把 `apply()` 裡那一行拔掉照樣全綠（2026-09-03 的突變測試發現）。
func TestResolvePendingScalesMorale(t *testing.T) {
	w := load(t, 0)
	alive := w.AliveFactions()
	fa, fb := alive[0], alive[1]
	na := w.clampCity(w.Factions[fa].Capital)
	nb := w.clampCity(w.Factions[fb].Capital)
	att := aiCorps(t, w, fa, na)
	def := aiCorps(t, w, fb, nb)

	// 兩支軍團的兵力與士氣設成好算的值。
	setMen := func(i, men, morale int) {
		c := &w.Corps[i]
		total := 0
		for k := range c.Units {
			c.Units[k].Men = men / len(c.Units)
			total += c.Units[k].Men
		}
		c.Men, c.Morale = total, morale
	}
	setMen(att, 120, 200)
	setMen(def, 120, 200)
	beforeAtt := w.Corps[att].Men

	// 手工掛一場打完的戰鬥：攻方勝，兩側各留一半的兵。
	b := tactical.NewBattle(nil, nil, rng.NewFixed(1), 0)
	b.Done, b.Winner = true, 0
	half := beforeAtt / 2
	b.Sides[0].Reserve[0] = half * tactical.MenPerSoldier / tactical.MenPerSoldier
	b.Sides[1].Reserve[0] = half * tactical.MenPerSoldier / tactical.MenPerSoldier
	w.pending = &Pending{Battle: b, Attacker: att, Defender: def,
		Node: na, Mode: combat.Field}

	if ev := w.ResolvePending(rng.NewFixed(1)); ev == nil {
		t.Fatal("結算不出結果")
	}
	// 勝方：戰前士氣 × 新兵力 ÷ 舊兵力。舊 120、新 60 ⇒ 200 × 60 ÷ 120 ＝ 100。
	want := postBattleMorale(200, beforeAtt, w.Corps[att].Men, true)
	if got := w.Corps[att].Morale; got != want {
		t.Errorf("勝方戰後士氣 = %d，want %d（舊 %d、新 %d）",
			got, want, beforeAtt, w.Corps[att].Men)
	}
	if w.Corps[att].Morale == 200 {
		t.Error("士氣完全沒被縮——接線沒接上")
	}
}
