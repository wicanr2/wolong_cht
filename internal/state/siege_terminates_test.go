package state_test

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

type siegeRand struct{}

func (siegeRand) Next() int { return 0 }

// battleFixture 建一場戰鬥，回傳規則層的 Battle。找不到素材就跳過。
func battleFixture(t *testing.T, siege bool, node, attacker, defender int) *tactical.Battle {
	t.Helper()
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
		Siege: siege, Node: node, Attacker: attacker, Defender: defender,
	})
	if w.PendingEncounter() != nil {
		if err := w.ChooseBattleCommand(); err != nil {
			t.Fatalf("ChooseBattleCommand: %v", err)
		}
	}
	pb := w.PendingBattle()
	if pb == nil || pb.Battle == nil {
		t.Skipf("這組參數（攻城=%v 節點=%d 攻=%d 守=%d）沒有開出戰鬥",
			siege, node, attacker, defender)
	}
	return pb.Battle
}

// TestSiegeFixtureTerminates 釘住「攻城戰會結束」（docs/spec/94）。
//
// ⭐ **這是會抓到那個死鎖的那一支測試。** 唯一的結束條件是「補不出兵」
// （`sub_1A6FA`），而攻方大將的體力被攻城計時器耗到 50 以下就全軍退卻——
// 退卻的兵彼此不能對調，只能靠繞路點繞過去。繞路點被每幀清掉、或是
// 「手上有路就不重算」，整排就卡在半路，兵一個都退不出去，
// `Done` 永遠是 false。**規則層的單元測試全綠，因為沒有人跑完整場。**
//
// 這條 fixture（據點 82、軍團 81 攻 39 守、玩家守方）現在在第 1,192 幀
// 結束。上限取 6,000 是留餘裕，不是期望值——**這一支只斷言「會結束」**，
// 誰贏、第幾幀都不是它管的（那兩個會隨規則層變）。
func TestSiegeFixtureTerminates(t *testing.T) {
	b := battleFixture(t, true, 82, 81, 39)

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

// TestSpawnHeightMatchesGroundPlane 釘住 docs/spec/95：
// 開場擺兵的 Z 要落在**移動用的那一層**。
//
// `Place()` 本來用 `StandLevel`（圖塊堆疊高度），而 `tryMove` 走一格時
// 把 Z 同步成 `GroundLevel`（地面層表）。兩個表不一樣，於是
// **沒動過的兵停在移動層之上**：攻方走到守方腳下卻差一層，
// `doAttack` 的碰撞與 `anyoneAt` 都比 Z，永遠打不到對方。
// 量到的後果是守方從頭到尾一兵未損（docs/playtest/51 §3）。
//
// 這條 fixture 修正前 96 個兵**每一個**的兩個值都不一樣。
func TestSpawnHeightMatchesGroundPlane(t *testing.T) {
	b := battleFixture(t, true, 82, 81, 39)

	bad := 0
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			s := &b.Sides[side].Soldiers[k]
			if !s.Alive {
				continue
			}
			lv, ok := b.Field.GroundLevel(s.X, s.Y, s.Plane())
			if !ok {
				continue // 那一格在這個平面沒有地面，退回堆疊高度是對的
			}
			if s.Z != lv {
				if bad < 3 {
					t.Errorf("側%d 兵%d (%d,%d) 站在 Z=%d，而地面層是 %d",
						side, k, s.X, s.Y, s.Z, lv)
				}
				bad++
			}
		}
	}
	if bad > 0 {
		t.Errorf("共 %d 個兵的開場高度不在移動層上", bad)
	}
}

// TestFieldBattleTerminates 是野戰那一半：照自然流程撞出一場遭遇，
// 然後跑到結束。退卻的閘與繞路點是兩種戰場共用的（docs/spec/94），
// 所以攻城那條死鎖在野戰同樣會發生。
//
// 存檔的產生方式在 docs/playtest/43 §2。
func TestFieldBattleTerminates(t *testing.T) {
	const save = "../../workplace/parity/SAVE-FIELD.DAT"
	if _, err := os.Stat(save); err != nil {
		t.Skipf("找不到 %s（產生方式見 docs/playtest/43 §2），跳過", save)
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

	rng := siegeRand{}
	for i := 0; i < 4000 && w.PendingEncounter() == nil; i++ {
		w.Tick(rng)
	}
	if w.PendingEncounter() == nil {
		t.Fatal("推了 4,000 tick 都沒有撞出遭遇——存檔的行軍狀態可能不對")
	}
	if err := w.ChooseBattleCommand(); err != nil {
		t.Fatalf("ChooseBattleCommand: %v", err)
	}
	pb := w.PendingBattle()
	if pb == nil || pb.Battle == nil {
		t.Fatal("選了戰鬥指揮卻沒有戰場")
	}
	b := pb.Battle

	const limit = 8000
	for b.Frame < limit && !b.Done {
		b.Step()
	}
	if !b.Done {
		t.Fatalf("野戰跑了 %d 幀還沒結束：側 0 剩 %d 兵、側 1 剩 %d 兵",
			b.Frame, b.Sides[0].Remaining(), b.Sides[1].Remaining())
	}
	t.Logf("第 %d 幀結束，勝方 側%d", b.Frame, b.Winner)
}
