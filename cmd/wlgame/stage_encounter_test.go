package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 野戰對拍那條 fixture：`workplace/parity/SAVE-FIELD.DAT` 裡
// 軍團 39（夏侯惇，玩家）與軍團 35（呂布）在同一格上遭遇
// （產生方式見 docs/playtest/43 §2）。
const (
	fieldParitySave = "../../workplace/parity/SAVE-FIELD.DAT"
	fieldParityMine = 39
	fieldParityFoe  = 35
)

// 攻城對拍那條 fixture：`parity-battle4/SAVE-E.DAT` 的
// 軍團 81（張遼，攻）對軍團 39（夏侯惇，守），據點 82
// （docs/playtest/51 §1）。門強度條只在攻城戰顯示，所以它是
// `-shot-when gate-bar` 唯一驗得到的局面（docs/spec/32 §2）。
const (
	siegeParitySave     = "../../workplace/promo-live/parity-battle4/SAVE-E.DAT"
	siegeParityNode     = 82
	siegeParityAttacker = 81
	siegeParityDefender = 39
)

// battleGameFixture 建一個只夠跑 `stageEncounter` 的 game。
func battleGameFixture(t *testing.T, save string, seed int, corps ...int) *game {
	t.Helper()
	const dir = "../../workplace/orig/dosv"
	lib, err := library.Load(dir)
	if err != nil {
		t.Skipf("找不到原版素材：%v", err)
	}
	w, err := state.LoadScenario(save, 0)
	if err != nil {
		t.Skipf("讀不到 %s：%v", save, err)
	}
	w.Player = 0
	p, setup, err := battlesetup.Load(battlesetup.Options{
		Dir: dir, World: w, Map: lib.World,
	})
	if err != nil {
		t.Fatalf("battlesetup.Load: %v", err)
	}
	w.SetTactical(setup)
	for _, c := range corps {
		if c >= len(w.Corps) || !w.Corps[c].Alive {
			t.Skipf("%s 裡沒有軍團 %d", save, c)
		}
	}
	return &game{lib: lib, world: w, battle: p, rng: rng.NewFixed(seed)}
}

// stageEncounterFixture 是野戰那一條（夏侯惇對呂布）。
func stageEncounterFixture(t *testing.T, seed int) *game {
	t.Helper()
	return battleGameFixture(t, fieldParitySave, seed, fieldParityMine, fieldParityFoe)
}

// siegeShotFixture 開一場攻城並停在第 0 拍，給 `-shot-when` 的條件用。
func siegeShotFixture(t *testing.T) *game {
	t.Helper()
	g := battleGameFixture(t, siegeParitySave, 7, siegeParityAttacker, siegeParityDefender)
	g.stageEncounter(true, siegeParityNode, 0,
		&g.world.Corps[siegeParityAttacker], &g.world.Corps[siegeParityDefender])
	if g.world.PendingBattle() == nil {
		t.Skip("攻城 fixture 沒有開出戰鬥")
	}
	return g
}

// TestStageEncounterArmsDuelBeforeStepping 釘住 docs/spec/117：
// 驗收捷徑要**先武裝開場喊話再推戰場**。
//
// ⭐ 這一支原本的失敗長相不是「錯的畫面」而是「少一塊」：
// `stageEncounter` 先推完 N 個 tick 才呼叫 `startBattleTalk`，
// 而 `SetDuelInput` 在那支裡面——單挑的開場只有 50 tick，
// 武裝的時候該喊話的那一刻已經過去，於是對拍看到的是一張
// **沒有挑戰對白框**的戰場，八個區照樣 PASS，只有 `field` 差 11%。
//
// 判準取 `DuelActive()`（挑戰喊話起、決著收尾止）：
// 第 0 拍必須是 false（開場 50 tick 還沒走完），第 52 拍必須是 true。
// **兩端都斷言**，少了前半這一支就分不出「真的武裝了」與「永遠回 true」。
func TestStageEncounterArmsDuelBeforeStepping(t *testing.T) {
	for _, tc := range []struct {
		steps int
		want  bool
	}{
		{0, false},
		{52, true},
	} {
		g := stageEncounterFixture(t, 7)
		g.stageEncounter(false, -1, tc.steps,
			&g.world.Corps[fieldParityMine], &g.world.Corps[fieldParityFoe])
		p := g.world.PendingBattle()
		if p == nil || p.Battle == nil {
			t.Fatalf("第 %d 拍：fixture 沒有開出戰鬥", tc.steps)
		}
		if got := p.Battle.DuelActive(); got != tc.want {
			t.Errorf("第 %d 拍 DuelActive() = %v，want %v——"+
				"false 表示單挑狀態機沒有在開場那 50 tick 之內武裝",
				tc.steps, got, tc.want)
		}
	}
}
