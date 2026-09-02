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

// stageEncounterFixture 建一個只夠跑 `stageEncounter` 的 game。
func stageEncounterFixture(t *testing.T, seed int) *game {
	t.Helper()
	const dir = "../../workplace/orig/dosv"
	lib, err := library.Load(dir)
	if err != nil {
		t.Skipf("找不到原版素材：%v", err)
	}
	w, err := state.LoadScenario(fieldParitySave, 0)
	if err != nil {
		t.Skipf("讀不到 %s：%v", fieldParitySave, err)
	}
	w.Player = 0
	p, setup, err := battlesetup.Load(battlesetup.Options{
		Dir: dir, World: w, Map: lib.World,
	})
	if err != nil {
		t.Fatalf("battlesetup.Load: %v", err)
	}
	w.SetTactical(setup)
	if fieldParityMine >= len(w.Corps) || fieldParityFoe >= len(w.Corps) ||
		!w.Corps[fieldParityMine].Alive || !w.Corps[fieldParityFoe].Alive {
		t.Skipf("%s 裡沒有軍團 %d／%d", fieldParitySave, fieldParityMine, fieldParityFoe)
	}
	return &game{lib: lib, world: w, battle: p, rng: rng.NewFixed(seed)}
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
