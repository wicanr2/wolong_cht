package state_test

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
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
	provider, setup, err := battlesetup.Load(battlesetup.Options{
		Dir: "../../workplace/orig/dosv", World: w, Map: lib.World,
	})
	if err != nil {
		t.Fatalf("battlesetup.Load: %v", err)
	}
	w.SetTactical(setup)
	// ⭐ **道路圖是正式路徑的一部分。** 少了它 `w.step` 走直線退路，
	// 而存檔裡「正在行軍」的軍團連路徑都還原不了（docs/spec/172 §4.5）——
	// 這個 fixture 先前就是靠直線碰巧撞出遭遇的。
	xy := make([][2]int, len(w.Cities))
	for i := range w.Cities {
		xy[i] = [2]int{w.Cities[i].X, w.Cities[i].Y}
	}
	edges, err := world.RoadEdges(lib.World, xy)
	if err != nil {
		t.Fatalf("RoadEdges: %v", err)
	}
	w.SetRoads(march.New(len(w.Cities), world.MarchEdges(edges, xy)))
	if n := w.UnresolvedMarches(); n != 0 {
		t.Fatalf("%d 支軍團的行軍狀態還原不了", n)
	}
	// ⭐ **正對照**：上面那一句在「沒有軍團需要還原」時也會通過，
	// 兩者在斷言上長得一樣。這一句把假零擋掉——這份存檔就是為了
	// 「有軍團正走在路上」而做的（docs/playtest/43 §2）。
	if n := w.RestoredMarches(); n == 0 {
		t.Fatal("這份存檔一支行軍中的軍團都沒有，" +
			"那上面那句 UnresolvedMarches()==0 什麼都沒驗到")
	}

	rng := siegeRand{}
	// ⚠ **上限要夠緊才擋得住「還是會撞到，只是晚了很多」**——那正是
	// 「靠直線退路碰巧撞出來」那種缺陷的樣子，而寬鬆的上限對它完全無感。
	// 實測是第 134 拍（道路圖接上、兩支行軍中的軍團還原之後）；
	// 600 是四倍餘裕。要改這個數字先確認變化有解釋。
	const meetLimit = 600
	met := -1
	for i := 0; i < meetLimit && w.PendingBattle() == nil; i++ {
		w.Tick(rng)
		if w.PendingBattle() != nil {
			met = i + 1
		}
	}
	if w.PendingBattle() == nil {
		t.Fatalf("推了 %d 拍都沒有撞出遭遇——存檔的行軍狀態可能不對，"+
			"或是軍團沒有走在道路圖上", meetLimit)
	}
	t.Logf("第 %d 拍撞出遭遇（還原了 %d 支行軍中的軍團）",
		met, w.RestoredMarches())
	pb := w.PendingBattle()
	if pb == nil || pb.Battle == nil {
		t.Fatal("選了戰鬥指揮卻沒有戰場")
	}
	b := pb.Battle
	// 正常玩家入口會由 startBattleTalk 武裝單挑，這裡照做。
	// ⚠ 原版的野戰單挑是**依武將個性自動觸發**的（武將 `+0x16` ＝ `Tactic`，
	// 同一個值也選 `BATTLE.DAT` 的腳本段）。remake 還沒解那個觸發判定，
	// 所以這裡是無條件武裝——它是 fixture 的簡化，不是原版行為。
	duel := tactical.DuelInput{FieldNumber: provider.FieldNumber(pb.Node, false)}
	for side, corps := range [2]int{pb.Attacker, pb.Defender} {
		g := w.Generals[w.Leader(corps)]
		duel.Martial[side], duel.CommandStat[side] = g.Martial, g.Command
	}
	b.SetDuelInput(duel)

	const limit = 8000
	// ⭐ **玩家側開場站在陣形線上等下令**（internal/state/tactical.go 的
	// `b.Order(side, -1, tactical.Form)`），而「攻擊」的定義是**大將以外的兵**
	// 攻擊、大將除非「突擊」不主動出擊。所以這一場要打得起來，必須補上
	// 原版流程裡的那一步：單挑結束後玩家下命令。沒有這一步，兩軍對峙到天荒
	// 地老才是**正確**的原版行為，不是規則層的缺陷。
	ordered := false
	for b.Frame < limit && !b.Done {
		b.Step()
		if !ordered && !b.OpeningActive() {
			b.OrderSelected(b.PlayerSide, tactical.Attack)
			ordered = true
		}
	}
	if !b.Done {
		t.Fatalf("野戰跑了 %d 幀還沒結束：側 0 剩 %d 兵、側 1 剩 %d 兵",
			b.Frame, b.Sides[0].Remaining(), b.Sides[1].Remaining())
	}
	t.Logf("第 %d 幀結束，勝方 側%d", b.Frame, b.Winner)
}
