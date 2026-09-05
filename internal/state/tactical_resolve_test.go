package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// TestTacticalBattleAlwaysResolves 是那個卡死的迴歸閘（docs/spec/62、63）。
//
// 兵的開場體力 ＝ 軍團士氣（docs/spec/61）之後前線會擠在一起，而原版的
// 攻擊是走進敵人格子撞出來的。少了「被換位的兵這一幀不動」，前排那一格
// 會一直換人，48 個兵圍著 2 個打卻一次也打不到——量到 95 萬 tick 零傷害。
//
// ⚠ 這一條要靠**跑完**來守，不能靠放寬上限：卡死時場面一個 byte 都不動，
// 所以「再多跑一點就好了」永遠不成立。
func TestTacticalBattleAlwaysResolves(t *testing.T) {
	w := load(t, 0)
	w.SetTactical(&TacticalSetup{
		Forms: tactical.SyntheticFormations(),
		Field: func(int, bool, bool) *tactical.Field {
			stack := make([][]int, tactical.Height)
			for y := range stack {
				stack[y] = make([]int, tactical.Width)
			}
			return tactical.NewField(stack, 0)
		},
	})
	alive := w.AliveFactions()
	a, b := alive[0], alive[1]
	w.Player = a
	for _, f := range []int{a, b} {
		w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
	}
	kinds := [army.Positions]army.TroopType{}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	la, lb := w.Factions[a].Lord, w.Factions[b].Lord
	w.FormCorps(la, kinds, manned)
	w.FormCorps(lb, kinds, manned)
	// ⚠ 先交戰，否則軍團在邊界掉頭（`docs/spec/132`）。
	w.Friendship[a][b] = w.Friendship[a][b].WithWar(true)
	w.Friendship[b][a] = w.Friendship[b][a].WithWar(true)
	w.March(la, w.Factions[b].Capital)
	w.March(lb, w.Factions[a].Capital)
	r := rng.NewFixed(5)
	for i := 0; i < 200000 && w.PendingBattle() == nil; i++ {
		w.Tick(r)
	}
	p := w.PendingBattle()
	if p == nil {
		t.Fatal("沒開戰術")
	}
	for n := 1; n <= 40; n++ {
		if p.Battle.Run(5000) {
			t.Logf("★ %d 千 tick 內打完，勝方 %d", n*5, p.Battle.Winner)
			return
		}
	}
	s0, s1 := &p.Battle.Sides[0], &p.Battle.Sides[1]
	t.Fatalf("20 萬 tick 沒打完：存活 %d/%d 待機 %v/%v",
		s0.Alive(), s1.Alive(), s0.Reserve, s1.Reserve)
}
