package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// ⭐ 開局的四個劇本本身就要滿足所有不變量。
//
// 這是最基本的一條：如果原版的檔案讀進來就不自洽，
// 那不是規則錯，是**解析錯**。
func TestScenariosStartConsistent(t *testing.T) {
	for s := 0; s < 4; s++ {
		w := load(t, s)
		if v := w.CheckInvariants(); len(v) > 0 {
			for _, x := range v {
				t.Errorf("劇本 %d 開局就違反：%s", s, x)
			}
		}
	}
}

// ⭐ 跑久了也要一直自洽。
//
// **這一條驗的是規則「組合起來」對不對。** 佔領、壞滅、招降、陣亡、
// 編成、月結會互相牽動同一批欄位；單元測試只保證單次呼叫正確，
// 保證不了跑三年之後那些冗餘計數還對得上。
//
// 一旦有一天不自洽，測試會報出**第幾天、哪一條、差多少**——
// 那正是往回查的入口。
func TestWorldStaysConsistent(t *testing.T) {
	const days = 3 * 360
	for s := 0; s < 4; s++ {
		w := load(t, s)
		r := rng.New(0, 0, s)
		for d := 0; d < days; d++ {
			for h := 0; h < 24; h++ {
				w.Tick(r)
			}
			if v := w.CheckInvariants(); len(v) > 0 {
				for i, x := range v {
					if i >= 3 {
						t.Errorf("劇本 %d：還有 %d 條沒列", s, len(v)-3)
						break
					}
					t.Errorf("劇本 %d 第 %d 天違反：%s", s, d, x)
				}
				break // 一天就夠，繼續跑只會噪音
			}
		}
	}
}

// ⭐ 把戰鬥那條鏈也跑進去。
//
// 上一條測試只跑內政與月結——**開局沒有軍團，佔領／壞滅／招降／陣亡
// 那幾條路一次都沒走到**，而那正是最會互相踩到的地方。
// 這裡主動編成軍團、反覆派去打對方的據點，每天檢查一次。
func TestInvariantsUnderWar(t *testing.T) {
	w := load(t, 0)
	r := rng.New(0, 0, 7)

	alive := w.AliveFactions()
	if len(alive) < 4 {
		t.Skip("勢力太少")
	}
	// 幫前幾個勢力各編一支軍團。
	var formed []int
	for _, f := range alive[:6] {
		w.Factions[f].Reserves = [economy.NumTroopTypes]int{9000, 9000, 9000}
		lord := w.Factions[f].Lord
		kinds := [army.Positions]army.TroopType{}
		manned := [army.Positions]bool{true, true, true, true, true, true}
		if err := w.FormCorps(lord, kinds, manned); err != nil {
			continue
		}
		formed = append(formed, lord)
	}
	if len(formed) < 2 {
		t.Fatal("一支軍團都編不出來")
	}
	if v := w.CheckInvariants(); len(v) > 0 {
		t.Fatalf("編成之後就不自洽：%v", v[0])
	}

	battles, captures := 0, 0
	for d := 0; d < 5*360; d++ {
		// 每隔一陣子把還活著的軍團派去打最近的敵方據點。
		if d%20 == 0 {
			for _, c := range formed {
				if !w.Corps[c].Alive {
					continue
				}
				if node := nearestEnemyCity(w, c); node >= 0 {
					_ = w.March(c, node)
				}
			}
		}
		for h := 0; h < 24; h++ {
			ev := w.Tick(r)
			for _, ce := range ev.Corps {
				if ce.Battle != nil {
					battles++
				}
				if ce.Captured >= 0 {
					captures++
				}
			}
		}
		if v := w.CheckInvariants(); len(v) > 0 {
			t.Fatalf("第 %d 天（打過 %d 場、佔了 %d 城）違反：%s",
				d, battles, captures, v[0])
		}
	}
	if battles == 0 {
		t.Error("五年下來一場都沒打起來——這條測試沒驗到戰鬥那條鏈")
	}
	t.Logf("五年、%d 場戰鬥、%d 次佔領，不變量全程成立", battles, captures)
}

// nearestEnemyCity 找離軍團最近的敵方據點。
func nearestEnemyCity(w *World, corps int) int {
	c := w.Corps[corps]
	best, bestD := -1, 1<<30
	for i := range w.Cities {
		city := &w.Cities[i]
		if city.Owner == c.Faction {
			continue
		}
		d := abs(city.X-c.X) + abs(city.Y-c.Y)
		if d < bestD {
			best, bestD = i, d
		}
	}
	return best
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// ⭐ 長跑觀察到「火災一次都沒有」——查下去是**資料決定的，不是規則錯**。
//
// `RollCityDisaster` 的火災閘是 `rng & 0x3F >= 防災值`，左邊只到 63；
// 而**開局 192 個據點的防災值全部是 100**（docs/formats/08 §1.6），
// 所以在沒有攻城的世界裡火災**在數學上不可能發生**。
// 防災值只有被攻城打掉才會降——降到 63 以下才輪得到火災。
//
// 這條測試把那個推論釘住：沒打仗 → 防災值不動 → 火災不可能；
// 打過城 → 防災值降下來 → 火災變得可能。
func TestFireNeedsSiegeDamageFirst(t *testing.T) {
	w := load(t, 0)
	const gate = 0x3F // rng & 0x3F 的上限

	min := 255
	for i := range w.Cities {
		if p := w.Cities[i].Prevention; p < min {
			min = p
		}
	}
	if min <= gate {
		t.Fatalf("開局最低防災值 %d 已經 ≤ %d，前提不成立", min, gate)
	}
	t.Logf("開局最低防災值 %d > %d → 火災不可能", min, gate)

	// 打幾場攻城，防災值就會被打下來。
	for i := range w.Cities {
		w.Cities[i].Prevention = 100
	}
	r := combat.Result{CityDamage: 20}
	for n := 0; n < 3; n++ {
		w.damageCity(0, combat.Siege, r)
	}
	if got := w.Cities[0].Prevention; got > gate {
		t.Errorf("打了三次攻城，防災值還有 %d（> %d），火災仍然不可能", got, gate)
	} else {
		t.Logf("三次攻城之後防災值 %d ≤ %d → 火災開始可能", got, gate)
	}
}
