package state

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
)

// 原版資產不隨本專案散布（CLAUDE.md §10），所以這些測試在沒有
// SINARIO.DAT 的環境會跳過，而不是失敗。
const origPath = "../../workplace/orig/dosv/SINARIO.DAT"

func load(t *testing.T, idx int) *World {
	t.Helper()
	if _, err := os.Stat(origPath); err != nil {
		t.Skip("找不到原版 SINARIO.DAT，跳過")
	}
	w, err := LoadScenario(origPath, idx)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

// 四個劇本的起始日必須與 docs/formats/08 §1.1 一致。
// 特別注意劇本 2 是 9 月、劇本 4 是 225 年 —— 與日文說明書的截圖不同。
func TestLoadClock(t *testing.T) {
	want := []struct{ y, m, d int }{
		{196, 4, 1}, {208, 9, 1}, {212, 6, 1}, {225, 5, 1},
	}
	for i, c := range want {
		w := load(t, i)
		if w.Clock.Year != c.y || w.Clock.Month != c.m || w.Clock.Day != c.d {
			t.Errorf("劇本 %d = %d/%d/%d, want %d/%d/%d",
				i+1, w.Clock.Year, w.Clock.Month, w.Clock.Day, c.y, c.m, c.d)
		}
		if w.Clock.Hour != 1 {
			t.Errorf("劇本 %d 的「時」= %d, want 1", i+1, w.Clock.Hour)
		}
	}
}

// 稅率預設 18%，低於 22.5% 的平衡點 —— 所以預設設定下據點會自然繁榮。
func TestDefaultTaxRate(t *testing.T) {
	for i := 0; i < 4; i++ {
		w := load(t, i)
		if w.TaxRate != 18 {
			t.Errorf("劇本 %d 的稅率 = %d, want 18", i+1, w.TaxRate)
		}
		for _, c := range w.RecruitCap {
			if c != 0 {
				t.Errorf("劇本 %d 的募兵數預設不是 0：%d", i+1, c)
			}
		}
	}
}

// 勢力數：22 / 11 / 6 / 4，隨著歷史推進而合併。
func TestFactionCounts(t *testing.T) {
	want := []int{22, 11, 6, 4}
	for i, n := range want {
		if got := len(load(t, i).AliveFactions()); got != n {
			t.Errorf("劇本 %d 的勢力數 = %d, want %d", i+1, got, n)
		}
	}
}

// 好戰等級：呂布最高、劉表最低，而且同一君主跨劇本是常數。
func TestAggressionKnownValues(t *testing.T) {
	w := load(t, 0)
	got := map[string]int{}
	for _, i := range w.AliveFactions() {
		got[text.Decode([]byte(w.LordName(i)), text.Big5)] = w.Factions[i].Aggression
	}
	want := map[string]int{"呂布": 15, "曹操": 14, "劉備": 4, "劉表": 1}
	for name, v := range want {
		if got[name] != v {
			t.Errorf("%s 的好戰等級 = %d, want %d", name, got[name], v)
		}
	}
	// 跨劇本常數。
	for _, idx := range []int{1, 2} {
		w2 := load(t, idx)
		for _, i := range w2.AliveFactions() {
			n := text.Decode([]byte(w2.LordName(i)), text.Big5)
			if v, ok := want[n]; ok && w2.Factions[i].Aggression != v {
				t.Errorf("劇本 %d 的 %s 好戰等級 = %d, want %d（應為個人常數）",
					idx+1, n, w2.Factions[i].Aggression, v)
			}
		}
	}
}

// 勢力記錄的據點數（+0x23）是用 **+0x1A** 那一欄算出來的。
func TestCityCountsMatchRecordedOwner(t *testing.T) {
	for idx := 0; idx < 4; idx++ {
		w := load(t, idx)
		counted := map[int]int{}
		for _, c := range w.Cities {
			counted[c.OwnerRecorded]++
		}
		for _, i := range w.AliveFactions() {
			if w.Factions[i].Cities != counted[i] {
				t.Errorf("劇本 %d 勢力 %d：記錄 %d 個據點，依 +0x1A 算出 %d 個",
					idx+1, i, w.Factions[i].Cities, counted[i])
			}
		}
	}
}

// ⚠ 原版的資料瑕疵：兩個據點的 +0x01 與 +0x1A 不一致。
//
// 執行期讀的是 +0x01（月結 sub_153C6 的 `cmp cl, [di+1]`），
// 所以遊戲實際會把武陵判給劉璋、南昌判給孫權以外的人 ——
// 與 +0x1A、與勢力記錄的據點數、與歷史都不符。
//
// 這條測試**釘住這個瑕疵**，而不是修掉它：remake 預設照抄原版行為。
// 如果哪天兩欄變成完全一致，代表資料被動過，要回頭確認。
func TestKnownOwnerDiscrepancies(t *testing.T) {
	type want struct{ scenario, city, runtime, recorded int }
	known := []want{
		{2, 171, 3, 2}, // 劇本 3 武陵：執行期劉璋，記錄劉備
		{3, 175, 2, 1}, // 劇本 4 南昌：執行期劉禪，記錄孫權
	}
	found := 0
	for idx := 0; idx < 4; idx++ {
		w := load(t, idx)
		for i, c := range w.Cities {
			if c.Owner == c.OwnerRecorded {
				continue
			}
			found++
			ok := false
			for _, k := range known {
				if k.scenario == idx && k.city == i &&
					k.runtime == c.Owner && k.recorded == c.OwnerRecorded {
					ok = true
				}
			}
			if !ok {
				t.Errorf("劇本 %d 據點 %d（%s）出現未記錄的不一致：+0x01=%d +0x1A=%d",
					idx+1, i, text.Decode([]byte(c.Name), text.Big5),
					c.Owner, c.OwnerRecorded)
			}
		}
	}
	if found != len(known) {
		t.Errorf("不一致的據點共 %d 個，預期 %d 個", found, len(known))
	}
}

// 城市名稱要解得出來。
func TestCityNames(t *testing.T) {
	w := load(t, 0)
	want := map[int]string{0: "北京", 2: "涿郡"}
	for i, n := range want {
		if got := text.Decode([]byte(w.Cities[i].Name), text.Big5); got != n {
			t.Errorf("據點 %d 的名稱 = %q, want %q", i, got, n)
		}
	}
	for i, c := range w.Cities {
		if c.Name == "" {
			t.Errorf("據點 %d 沒有名稱", i)
		}
	}
}

// 武將的所屬勢力分佈要等於勢力記錄的武將數。
func TestGeneralCountsConsistent(t *testing.T) {
	w := load(t, 0)
	counted := map[int]int{}
	for _, g := range w.Generals {
		if g.Alive {
			counted[g.Faction]++
		}
	}
	for _, i := range w.AliveFactions() {
		if w.Factions[i].Generals != counted[i] {
			t.Errorf("勢力 %d：記錄 %d 名武將，實際 %d 名", i, w.Factions[i].Generals, counted[i])
		}
	}
}

// 座標與生產力的值域，以及「生產力不超過上限」這條不變式。
func TestCityInvariants(t *testing.T) {
	for idx := 0; idx < 4; idx++ {
		w := load(t, idx)
		for i, c := range w.Cities {
			if c.X < 0 || c.X >= 384 || c.Y < 0 || c.Y >= 256 {
				t.Errorf("劇本 %d 據點 %d 座標 (%d,%d) 超出 384×256", idx+1, i, c.X, c.Y)
			}
			if c.Production > c.ProductionCap {
				t.Errorf("劇本 %d 據點 %d 生產力 %d > 上限 %d",
					idx+1, i, c.Production, c.ProductionCap)
			}
			if c.Garrison > c.GarrisonCap {
				t.Errorf("劇本 %d 據點 %d 城兵 %d > 上限 %d",
					idx+1, i, c.Garrison, c.GarrisonCap)
			}
		}
	}
}

// 信賴度開局全部 200；官員全部未派駐。
func TestInitialAssignments(t *testing.T) {
	w := load(t, 0)
	for _, i := range w.AliveFactions() {
		if w.Factions[i].Trust != 200 {
			t.Errorf("勢力 %d 的信賴度 = %d, want 200", i, w.Factions[i].Trust)
		}
		if w.Factions[i].Diplomat != 0xFF {
			t.Errorf("勢力 %d 開局就有外交官？", i)
		}
	}
	for i, c := range w.Cities {
		if c.Governor != 0xFF {
			t.Errorf("據點 %d 開局就有內政官？", i)
		}
	}
}

// ⭐ 端到端：稅率 18% 十年後生產力應該成長，25% 應該崩到 0。
// 這條驗的是「從機器碼推出的平衡點 22.5%」——
// 純單元測試驗不出來，因為經濟是複利模型，要跑幾十個月才看得出方向。
func TestTaxTippingPoint(t *testing.T) {
	run := func(tax int) int {
		w := load(t, 0)
		w.Player = 0
		w.TaxRate, w.NextTaxRate = tax, tax
		rng := &testRand{s: 12345}
		for months := 0; months < 120; {
			if w.Tick(rng).Settled {
				months++
			}
		}
		sum, n := 0, 0
		for _, c := range w.Cities {
			if c.Owner == 0 {
				sum += c.Production
				n++
			}
		}
		if n == 0 {
			return 0
		}
		return sum / n
	}
	below, above := run(18), run(25)
	if below <= above {
		t.Errorf("稅率 18%% 的平均生產力 %d 應該遠高於 25%% 的 %d", below, above)
	}
	if above != 0 {
		t.Errorf("稅率 25%% 十年後平均生產力 = %d, want 0（應該崩潰）", above)
	}
}

type testRand struct{ s uint32 }

func (r *testRand) Next() int {
	r.s = r.s*1664525 + 1013904223
	return int(r.s >> 16)
}

var _ economy.Rand = (*testRand)(nil)
