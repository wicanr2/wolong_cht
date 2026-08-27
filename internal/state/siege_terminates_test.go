package state_test

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/state"
)

type siegeRand struct{}

func (siegeRand) Next() int { return 0 }

// TestSiegeFixtureTerminates 釘住「攻城戰會結束」（docs/spec/94）。
//
// ⭐ **這是會抓到那個死鎖的那一支測試。** 唯一的結束條件是「補不出兵」
// （`sub_1A6FA`），而攻方大將的體力被攻城計時器耗到 50 以下就全軍退卻——
// 退卻的兵彼此不能對調，只能靠繞路點繞過去。繞路點被每幀清掉、或是
// 「手上有路就不重算」，整排就卡在半路，兵一個都退不出去，
// `Done` 永遠是 false。**規則層的單元測試全綠，因為沒有人跑完整場。**
//
// 這條 fixture（據點 82、軍團 81 攻 39 守、玩家守方）修好之後在
// 第 1,416 幀結束，守方勝。上限取 6,000 是留餘裕，不是期望值。
func TestSiegeFixtureTerminates(t *testing.T) {
	const save = "../../workplace/promo-live/parity-battle4/SAVE-E.DAT"
	if _, err := os.Stat(save); err != nil {
		t.Skipf("找不到 %s，跳過", save)
	}
	w, err := state.LoadScenario(save, 0)
	if err != nil {
		t.Skipf("讀不到存檔：%v", err)
	}
	w.Player = 0
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("讀不到原版素材：%v", err)
	}
	_, setup, err := battlesetup.Load(battlesetup.Options{
		Dir: "../../workplace/orig/dosv", World: w, Map: lib.World,
	})
	if err != nil {
		t.Fatalf("battlesetup.Load: %v", err)
	}
	w.SetTactical(setup)
	battlesetup.StageEncounter(w, siegeRand{}, battlesetup.StageOptions{
		Siege: true, Node: 82, Attacker: 81, Defender: 39,
	})
	if w.PendingEncounter() != nil {
		if err := w.ChooseBattleCommand(); err != nil {
			t.Fatalf("ChooseBattleCommand: %v", err)
		}
	}
	b := w.PendingBattle().Battle

	const limit = 6000
	for b.Frame < limit && !b.Done {
		b.Step()
	}
	if !b.Done {
		t.Fatalf("跑了 %d 幀還沒結束：側 0 剩 %d 兵、側 1 剩 %d 兵——"+
			"攻城戰卡住了", b.Frame, b.Sides[0].Remaining(), b.Sides[1].Remaining())
	}
	t.Logf("第 %d 幀結束，勝方 側%d", b.Frame, b.Winner)
}
