package state

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/capital"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
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
		if w.Factions[i].MoraleBase != 200 {
			t.Errorf("勢力 %d 的士氣基準 = %d, want 200", i, w.Factions[i].MoraleBase)
		}
		if w.Factions[i].InvasionTarget != 0xFF {
			t.Errorf("勢力 %d 開局就在侵攻 %d？", i, w.Factions[i].InvasionTarget)
		}
		if w.Factions[i].Corps != 0 {
			t.Errorf("勢力 %d 開局就有 %d 個軍團？", i, w.Factions[i].Corps)
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

// ⭐ 端到端：稅率低於平衡點十年後生產力會長到上限，高於平衡點會崩到 0。
// 純單元測試驗不出來，因為經濟是複利模型，要跑幾十個月才看得出方向。
//
// ⚠ **平衡點從 22.5% 移到 33.75%**（2026-08-08）。原因不是調參，
// 是補上了一整層機制：內政官（`sub_14194`，`docs/re/07` §19）。
//
//	舊：上昇值/月 = −(稅率−30) − 7.5            → 平衡點 22.5%
//	新：上昇值/月 = −(稅率−30) − 7.5 + 11.25    → 平衡點 33.75%
//	                                  ↑ 據點整備（玩家據點無內政官：
//	                                    rate=5 → 成功 6/16、gain=1 → +11.25/月）
//
// 舊版這條測試斷言「25% 十年後崩到 0」，那是**缺了整備那一層的產物**。
// 這種測試很危險：它會把模型的缺陷釘成規格，讓補上缺陷的那一次改動
// 看起來像「弄壞了東西」。改的時候要問的是**哪一邊才對得上機器碼**。
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
	// 18% 與 25% 都在平衡點以下 → 都長到生產力上限（所以會相等，不比大小）。
	for _, tax := range []int{18, 25} {
		if got := run(tax); got == 0 {
			t.Errorf("稅率 %d%% 在平衡點以下，十年後不該崩到 0", tax)
		}
	}
	// 40% 在平衡點以上 → 上昇值每月淨 −6.25，撐不過十年。
	below, above := run(18), run(40)
	if below <= above {
		t.Errorf("稅率 18%% 的平均生產力 %d 應該遠高於 40%% 的 %d", below, above)
	}
	if above != 0 {
		t.Errorf("稅率 40%% 十年後平均生產力 = %d, want 0（應該崩潰）", above)
	}
}

type testRand struct{ s uint32 }

func (r *testRand) Next() int {
	r.s = r.s*1664525 + 1013904223
	return int(r.s >> 16)
}

var _ economy.Rand = (*testRand)(nil)

// ---------------------------------------------------------------------------
// 存檔寫回
// ---------------------------------------------------------------------------

// ⭐ Round-trip：載入之後原封不動寫回，必須與原始位元組**完全相同**。
//
// 這是「改寫而非重建」策略的驗收條件。只要有任何一個未解欄位被誤寫成 0，
// 或任何一個已解欄位的編碼寫錯（偏移、位元組序、+100 偏移…），
// 這條測試就會抓到。
func TestSaveRoundTrip(t *testing.T) {
	raw, err := os.ReadFile(origPath)
	if err != nil {
		t.Skip("找不到原版 SINARIO.DAT，跳過")
	}
	for idx := 0; idx < 4; idx++ {
		w, err := LoadScenario(origPath, idx)
		if err != nil {
			t.Fatal(err)
		}
		got := w.Bytes()
		want := raw[idx*blockSize : (idx+1)*blockSize]
		if len(got) != len(want) {
			t.Fatalf("劇本 %d：長度 %d, want %d", idx+1, len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("劇本 %d：偏移 0x%X 不同（得到 %02X，原始 %02X）",
					idx+1, i, got[i], want[i])
			}
		}
	}
}

// 改過的欄位要真的寫進去，而且只動該動的地方。
func TestSaveWritesChanges(t *testing.T) {
	w := load(t, 0)
	before := w.Bytes()

	w.Factions[0].Funds = -12345
	w.Factions[0].MoraleBase = 77
	w.Cities[0].Growth = -50
	w.Cities[0].Production = 4242
	w.Clock.Month = 7
	w.TaxRate = 42

	after := w.Bytes()
	if len(before) != len(after) {
		t.Fatal("長度變了")
	}
	diff := 0
	for i := range before {
		if before[i] != after[i] {
			diff++
		}
	}
	// 資金 3 B ＋ 信賴度 1 ＋ 上昇值 1 ＋ 生產力 2 ＋ 月 1（u16 高位不變）
	// ＋ 該月天數 1 ＋ 稅率 1 ＝ 10 個 byte。
	if diff != 10 {
		t.Errorf("改了 %d 個 byte，預期 10（多出來的表示動到不該動的地方）", diff)
	}

	// 讀回來要一致。
	reloaded, err := reparse(after)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Factions[0].Funds != -12345 {
		t.Errorf("資金 = %d, want −12345", reloaded.Factions[0].Funds)
	}
	if reloaded.Cities[0].Growth != -50 {
		t.Errorf("上昇值 = %d, want −50", reloaded.Cities[0].Growth)
	}
	if reloaded.TaxRate != 42 {
		t.Errorf("稅率 = %d, want 42", reloaded.TaxRate)
	}
	if reloaded.Clock.Month != 7 {
		t.Errorf("月 = %d, want 7", reloaded.Clock.Month)
	}
	// 換月之後該月天數要跟著更新，否則進位判斷會用到舊值。
	if after[0x01] != 31 {
		t.Errorf("7 月的天數欄位 = %d, want 31", after[0x01])
	}
}

// 資金寫回時要照原版的上下限鉗住。
func TestSaveClampsFunds(t *testing.T) {
	w := load(t, 0)
	w.Factions[0].Funds = 99_999_999
	got, err := reparse(w.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if got.Factions[0].Funds != 655000 {
		t.Errorf("資金 = %d, want 655000（上限）", got.Factions[0].Funds)
	}
	w.Factions[0].Funds = -99_999_999
	got, _ = reparse(w.Bytes())
	if got.Factions[0].Funds != -655000 {
		t.Errorf("資金 = %d, want −655000（下限）", got.Factions[0].Funds)
	}
}

// 跑過一段時間之後存檔，未解區域仍然必須與原始檔完全相同。
// **這才是「改寫而非重建」真正要防的事**——遊戲跑了一陣子之後
// 才存檔，是最容易把沒理解的區域寫壞的時機。
func TestSaveAfterSimulationPreservesUnknownRegions(t *testing.T) {
	raw, err := os.ReadFile(origPath)
	if err != nil {
		t.Skip("找不到原版 SINARIO.DAT，跳過")
	}
	w := load(t, 0)
	w.Player = 0
	rng := &testRand{s: 999}
	for months := 0; months < 24; {
		if w.Tick(rng).Settled {
			months++
		}
	}
	got := w.Bytes()
	orig := raw[:blockSize]

	// 這些區間是**還沒解的**，跑再久也不該被動到。
	regions := []struct {
		name     string
		from, to int
	}{
		{"+0x3B–0x7F 不載入的空隙", 0x3B, 0x80},
		{"軍團表", 0x22C0, 0x42C0},
		{"事件佇列", 0x52C0, 0x56C0},
	}
	for _, r := range regions {
		for i := r.from; i < r.to; i++ {
			if got[i] != orig[i] {
				t.Fatalf("%s 的偏移 0x%X 被動到了（%02X → %02X）",
					r.name, i, orig[i], got[i])
			}
		}
	}
}

// reparse 把一個區塊的位元組重新解析成 World，測試用。
func reparse(block []byte) (*World, error) {
	full := make([]byte, blockSize*numBlocks)
	copy(full, block)
	f, err := os.CreateTemp("", "wolong-*.dat")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(full); err != nil {
		return nil, err
	}
	f.Close()
	return LoadScenario(f.Name(), 0)
}

// ⭐ 交友度矩陣：位址算法出自 sub_13119（0x600 + 觀察者×24 + 對象）。
//
// 驗收用的是**歷史事實**，不是我推的數字：
//
//	劉備 ↔ 公孫瓚 = 100（同門，公孫瓚曾收留劉備）
//	孫策 → 劉表  = 20（孫堅死於劉表部將黃祖之手）
func TestFriendshipMatrix(t *testing.T) {
	w := load(t, 0)

	name := func(i int) string {
		return text.Decode([]byte(w.LordName(i)), text.Big5)
	}
	idx := map[string]int{}
	for _, i := range w.AliveFactions() {
		idx[name(i)] = i
	}

	liubei, gongsun := idx["劉備"], idx["公孫瓚"]
	if got := w.Friendship[liubei][gongsun].Value(); got != 100 {
		t.Errorf("劉備看公孫瓚 = %d, want 100", got)
	}
	if got := w.Friendship[gongsun][liubei].Value(); got != 100 {
		t.Errorf("公孫瓚看劉備 = %d, want 100", got)
	}

	sunce, liubiao := idx["孫策"], idx["劉表"]
	if got := w.Friendship[sunce][liubiao].Value(); got != 20 {
		t.Errorf("孫策看劉表 = %d, want 20（父仇）", got)
	}
	// **但劉表看孫策不是 20** —— 交友度是有向的。
	if w.Friendship[liubiao][sunce].Value() == 20 {
		t.Error("交友度應該是有向的，兩個方向不該一樣")
	}

	// 對角線是自己，值 127（超過一般上限 100）。
	for _, i := range w.AliveFactions() {
		if got := w.Friendship[i][i].Value(); got != 127 {
			t.Errorf("勢力 %d 對自己 = %d, want 127", i, got)
		}
	}
}

// ⭐ 開局沒有任何交戰：每一格都帶著「和平」位元，
// 而且沒有任何勢力有侵攻目標。這兩件事必須一致 ——
// 它們就是「最高位元 = 和平而不是交戰」的證據。
func TestNoWarAtScenarioStart(t *testing.T) {
	for idx := 0; idx < 4; idx++ {
		w := load(t, idx)
		for _, i := range w.AliveFactions() {
			if w.Factions[i].InvasionTarget != 0xFF {
				t.Errorf("劇本 %d 勢力 %d 開局就有侵攻目標", idx+1, i)
			}
			for j := range w.Friendship[i] {
				if w.Friendship[i][j].AtWar() {
					t.Errorf("劇本 %d：%d 對 %d 開局就交戰", idx+1, i, j)
				}
			}
		}
	}
}

// 每「時」只處理一個勢力，22 個勢力輪一圈——所以每個勢力大約
// 每天被處理一次，不是每 tick 一次。這個節奏決定了外交官的效率，
// 寫錯會讓交友度以 9 倍速度上漲（一「時」有 9 個子刻）。
func TestHourlyRotation(t *testing.T) {
	w := load(t, 0)
	rng := rng.NewFixed(7)

	seen := map[int]int{}
	hours := 0
	for hours < 44 { // 兩圈
		ev := w.Tick(rng)
		if !ev.Clock.Hour {
			if ev.HourFaction != -1 {
				t.Fatalf("沒有進位到新的「時」卻輪到了勢力 %d", ev.HourFaction)
			}
			continue
		}
		hours++
		seen[ev.HourFaction]++
	}
	if len(seen) != 22 {
		t.Errorf("兩圈只輪到 %d 個勢力，應為 22", len(seen))
	}
	for i, n := range seen {
		if n != 2 {
			t.Errorf("勢力 %d 在兩圈裡被處理 %d 次，應為 2", i, n)
		}
	}
}

// 財政撐不住就自動取消侵攻。門檻與據點數掛鉤，
// 這是政略 AI 的硬性煞車（docs/re/08 §1）。
func TestBrokeFactionAbandonsInvasion(t *testing.T) {
	w := load(t, 0)
	rng := rng.NewFixed(3)

	// 挑一個活著的勢力，塞給它一個侵攻目標再把錢抽乾。
	victim := w.AliveFactions()[0]
	w.Factions[victim].InvasionTarget = (victim + 1) % 22
	w.Factions[victim].Funds = 0

	for i := 0; i < 22*clock.SubticksPerHour+9; i++ {
		ev := w.Tick(rng)
		if ev.HourFaction == victim {
			break
		}
	}
	if got := w.Factions[victim].InvasionTarget; got != diplomacy.NoTarget {
		t.Errorf("沒錢的勢力仍在侵攻 %d，應被清成 0xFF", got)
	}
	if !w.Factions[victim].LowFunds {
		t.Error("資金 0 遠低於門檻的一半，bit 6 應設起")
	}
}

// 預備兵維持費每小時累加，月結時才付。預備兵越多扣得越多——
// 說明書 5.2 說的是「月単位」，實際上是每小時扣、月末結算。
func TestReserveUpkeepAccumulatesHourly(t *testing.T) {
	w := load(t, 0)
	rng := rng.NewFixed(11)

	f := w.AliveFactions()[0]
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{3200, 1600, 1600}
	w.Factions[f].Expense = 0

	for i := 0; i < 22*clock.SubticksPerHour+9; i++ {
		if ev := w.Tick(rng); ev.HourFaction == f {
			break
		}
	}
	// (3200 + 1600 + 1600) ÷ 32 = 200
	if got := w.Factions[f].Expense; got != 200 {
		t.Errorf("累計支出 ＝ %d，應為 200", got)
	}
}

// 編成一支軍團：兵從預備兵扣、士氣繼承勢力的基準、位置在首都、
// 武將標成出陣中、勢力的軍團數 +1。
func TestFormCorps(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	lord := w.Factions[f].Lord
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
	before := w.Factions[f].Corps

	kinds := [army.Positions]army.TroopType{
		army.Cavalry, army.Cavalry, army.Archer,
		army.Archer, army.Infantry, army.Infantry,
	}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	if err := w.FormCorps(lord, kinds, manned); err != nil {
		t.Fatalf("編成失敗：%v", err)
	}

	c := w.Corps[lord]
	if !c.Alive {
		t.Fatal("軍團沒有建立")
	}
	if c.Morale != w.Factions[f].MoraleBase {
		t.Errorf("士氣 %d，應繼承勢力基準 %d", c.Morale, w.Factions[f].MoraleBase)
	}
	if c.Node != w.Factions[f].Capital {
		t.Errorf("位置在據點 %d，應在首都 %d", c.Node, w.Factions[f].Capital)
	}
	if c.Men != 600 { // 六槽 × 100 點 = 6,000 人
		t.Errorf("兵力 %d 點，應為 600（＝6,000 人）", c.Men)
	}
	if !w.Generals[lord].Posted {
		t.Error("武將沒有被標成出陣中")
	}
	if w.Factions[f].Corps != before+1 {
		t.Errorf("勢力軍團數 %d，應為 %d", w.Factions[f].Corps, before+1)
	}
	// 騎馬扣 2,000、弓兵 2,000、步兵 2,000。
	for tp, want := range map[economy.TroopType]int{
		economy.Cavalry: 4000, economy.Archer: 4000, economy.Infantry: 4000,
	} {
		if got := w.Factions[f].Reserves[tp]; got != want {
			t.Errorf("%v 預備兵剩 %d，應為 %d", tp, got, want)
		}
	}
	// 同一個武將不能編兩支。
	if err := w.FormCorps(lord, kinds, manned); err == nil {
		t.Error("同一個武將編了第二支軍團")
	}
}

// 預備兵不足就整批不做——原版沒有「編一半」這回事。
func TestFormCorpsAllOrNothing(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	lord := w.Factions[f].Lord
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 0, 6000}

	kinds := [army.Positions]army.TroopType{
		army.Cavalry, army.Cavalry, army.Archer, army.Archer, army.Cavalry, army.Cavalry,
	}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	if err := w.FormCorps(lord, kinds, manned); err == nil {
		t.Fatal("弓兵不足卻編成成功了")
	}
	if w.Corps[lord].Alive {
		t.Error("失敗卻留下了軍團")
	}
	if w.Factions[f].Reserves[economy.Cavalry] != 6000 {
		t.Error("失敗卻扣了騎馬預備兵")
	}
}

// 純騎馬編成走得快（說明書 5.5）。
func TestAllCavalryMarchesFaster(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{20000, 20000, 20000}

	cav := [army.Positions]army.TroopType{}
	mixed := [army.Positions]army.TroopType{
		army.Cavalry, army.Cavalry, army.Cavalry, army.Cavalry, army.Cavalry, army.Infantry,
	}
	manned := [army.Positions]bool{true, true, true, true, true, true}

	a := w.Factions[f].Lord
	b := -1
	for i := range w.Generals {
		if w.Generals[i].Alive && w.Generals[i].Faction == f && i != a {
			b = i
			break
		}
	}
	if b < 0 {
		t.Skip("這個勢力只有一名武將")
	}
	if err := w.FormCorps(a, cav, manned); err != nil {
		t.Fatal(err)
	}
	if err := w.FormCorps(b, mixed, manned); err != nil {
		t.Fatal(err)
	}
	if w.Corps[a].Interval >= w.Corps[b].Interval {
		t.Errorf("純騎馬間隔 %d 應小於混編的 %d",
			w.Corps[a].Interval, w.Corps[b].Interval)
	}
}

// 軍團表要能原樣寫回：載入 → 編成 → 寫回 → 再載入，欄位一致。
func TestCorpsRoundTrip(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	lord := w.Factions[f].Lord
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
	kinds := [army.Positions]army.TroopType{
		army.Cavalry, army.Archer, army.Infantry,
		army.Cavalry, army.Archer, army.Infantry,
	}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	if err := w.FormCorps(lord, kinds, manned); err != nil {
		t.Fatal(err)
	}
	if err := w.March(lord, 42); err != nil {
		t.Fatal(err)
	}

	b := w.Bytes()
	w2 := &World{raw: b}
	w2.loadCorps(b)

	a, c := w.Corps[lord], w2.Corps[lord]
	if a != c {
		t.Errorf("寫回後不一致：\n 原 %+v\n 後 %+v", a, c)
	}
	// 沒有軍團的槽一個 byte 都不該動。
	orig := load(t, 0)
	for i := range orig.Corps {
		if i == lord {
			continue
		}
		if w2.Corps[i] != orig.Corps[i] {
			t.Fatalf("軍團 %d 被動到了", i)
		}
	}
}

// 每 tick 只更新 16 支軍團，一輪要 8 個 tick——不是全部一起動。
func TestCorpsCursorRotates(t *testing.T) {
	w := load(t, 0)
	seen := map[int]bool{}
	for n := 0; n < 8; n++ {
		start := w.corpsCursor
		for k := 0; k < 16; k++ {
			seen[(start+k)%127] = true
		}
		w.tickCorps(0, rng.NewFixed(1))
	}
	if len(seen) != 127 {
		t.Errorf("8 個 tick 掃到 %d 支軍團，應為全部 127 支", len(seen))
	}
}

// 端對端：兩個勢力的軍團往對方的首都走，途中撞上就打一場。
//
// 這條把「編成 → 行軍 → 遭遇 → 自動判定 → 傷亡／壞滅／敗將下場」
// 整條鏈接起來跑。單元測試各自驗過每一段，這裡驗的是**接得起來**。
func TestCorpsMeetAndFight(t *testing.T) {
	w := load(t, 0)
	alive := w.AliveFactions()
	if len(alive) < 2 {
		t.Skip("這個劇本只有一個勢力")
	}
	a, b := alive[0], alive[1]
	for _, f := range []int{a, b} {
		w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
	}
	kinds := [army.Positions]army.TroopType{}
	manned := [army.Positions]bool{true, true, true, true, true, true}

	la, lb := w.Factions[a].Lord, w.Factions[b].Lord
	if err := w.FormCorps(la, kinds, manned); err != nil {
		t.Fatal(err)
	}
	if err := w.FormCorps(lb, kinds, manned); err != nil {
		t.Fatal(err)
	}
	// 互相往對方的首都走。
	if err := w.March(la, w.Factions[b].Capital); err != nil {
		t.Fatal(err)
	}
	if err := w.March(lb, w.Factions[a].Capital); err != nil {
		t.Fatal(err)
	}

	r := rng.NewFixed(5)
	var fought *CorpsEvent
	for i := 0; i < 200000 && fought == nil; i++ {
		ev := w.Tick(r)
		for k := range ev.Corps {
			// 要的是兩支軍團的野戰遭遇，不是走到城下打城兵。
			if ev.Corps[k].Battle != nil && ev.Corps[k].Enemy >= 0 {
				fought = &ev.Corps[k]
			}
		}
	}
	if fought == nil {
		t.Fatalf("兩支軍團一路走到底都沒有交戰\n a:(%d,%d)→(%d,%d)  b:(%d,%d)→(%d,%d)",
			w.Corps[la].X, w.Corps[la].Y, w.Corps[la].TargetX, w.Corps[la].TargetY,
			w.Corps[lb].X, w.Corps[lb].Y, w.Corps[lb].TargetX, w.Corps[lb].TargetY)
	}
	if fought.Battle.Ratio < 8 {
		t.Errorf("戰力比值 %d，最小應為 8（勢均力敵）", fought.Battle.Ratio)
	}
	// 打完之後至少有一方掉了兵。
	if w.Corps[la].Men == 600 && w.Corps[lb].Men == 600 {
		t.Error("打了一場卻兩邊都沒有傷亡")
	}
	t.Logf("交戰：軍團 %d vs %d，比值 %d，守方勝 %v，壞滅 %v",
		fought.Corps, fought.Enemy, fought.Battle.Ratio,
		fought.Battle.DefenderWins, fought.Destroyed)

	// 敗方的士氣被重設成 100 × 兵力比，所以打輸一場之後必定低於 100
	// （docs/re/09 §4.4）。這是接起來之後才看得到的行為。
	loser := fought.Corps
	if !fought.Battle.DefenderWins {
		loser = fought.Enemy
	}
	if m := w.Corps[loser].Morale; w.Corps[loser].Alive && m >= 100 {
		t.Errorf("敗方軍團 %d 戰後士氣 %d，應低於 100", loser, m)
	}
}

// 大將的位置一定要有兵——原版的壞滅判定直接看第一槽，
// 空著的話軍團一編出來就會被判掉（docs/re/09 §5）。
func TestFormCorpsNeedsGeneralSlot(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	lord := w.Factions[f].Lord
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}

	kinds := [army.Positions]army.TroopType{}
	// 第一槽空著，其餘有兵。
	manned := [army.Positions]bool{false, true, true, true, true, true}
	if err := w.FormCorps(lord, kinds, manned); err == nil {
		t.Fatal("大將空著卻編成成功了")
	}
	// 全空更不行。
	if err := w.FormCorps(lord, kinds, [army.Positions]bool{}); err == nil {
		t.Fatal("六個位置全空卻編成成功了")
	}
	if w.Corps[lord].Alive {
		t.Error("失敗卻留下了軍團")
	}
}

// 玩家的勢力捲進去就開戰術畫面，其餘自動判定——原版的分派規則。
func TestPlayerBattleGoesTactical(t *testing.T) {
	w := load(t, 0)
	w.SetTactical(&TacticalSetup{
		Forms: tactical.SyntheticFormations(),
		Field: func(int, bool) *tactical.Field {
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
	if err := w.FormCorps(la, kinds, manned); err != nil {
		t.Fatal(err)
	}
	if err := w.FormCorps(lb, kinds, manned); err != nil {
		t.Fatal(err)
	}
	if err := w.March(la, w.Factions[b].Capital); err != nil {
		t.Fatal(err)
	}
	if err := w.March(lb, w.Factions[a].Capital); err != nil {
		t.Fatal(err)
	}

	r := rng.NewFixed(5)
	for i := 0; i < 200000 && w.PendingBattle() == nil; i++ {
		w.Tick(r)
	}
	p := w.PendingBattle()
	if p == nil {
		t.Fatal("玩家的軍團一路走到底都沒有開戰術畫面")
	}
	if got := p.Battle.Sides[0].Alive(); got != tactical.SoldiersOnFoot {
		t.Errorf("攻方場上 %d 個兵，應為 %d", got, tactical.SoldiersOnFoot)
	}

	// 有戰鬥掛著的時候世界不前進。
	before := w.Clock
	w.Tick(r)
	if w.Clock != before {
		t.Error("戰術戰鬥還沒打完，世界卻繼續走了")
	}

	if !p.Battle.Run(200000) {
		t.Fatal("戰術戰鬥跑不完")
	}
	ev := w.ResolvePending(r)
	if ev == nil || ev.Battle == nil {
		t.Fatal("結算不出結果")
	}
	if w.PendingBattle() != nil {
		t.Error("結算完了卻還掛著")
	}
	// 結算完世界要能繼續走。
	before = w.Clock
	w.Tick(r)
	if w.Clock == before {
		t.Error("結算完世界還是停著")
	}
	t.Logf("戰術戰鬥 %d 幀，守方勝 %v；攻方剩 %d 點、守方剩 %d 點",
		p.Battle.Frame, ev.Battle.DefenderWins, w.Corps[la].Men, w.Corps[lb].Men)
}

// ⭐ 走進**有守軍的敵方據點**要打攻城，不是野戰。
//
// 原版的判定順序其實反過來（先看佔用圖再看據點圖塊），但那建立在
// 「一個據點佔好幾格地圖」上：攻方踏進的通常是據點的**別格**，
// 那幾格沒有人，所以走攻城，再由 `sub_14C72` 用據點座標把守軍找出來。
// 本專案的據點是一個點，照抄順序的話攻城那條路永遠走不到——
// 所以 `resolveContact` 把據點放在野戰前面（docs/re/09 §2）。
func TestMarchIntoDefendedCityIsSiege(t *testing.T) {
	w := load(t, 0)
	alive := w.AliveFactions()
	a, b := alive[0], alive[1]
	for _, f := range []int{a, b} {
		w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
	}
	kinds := [army.Positions]army.TroopType{}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	att, def := w.Factions[a].Lord, w.Factions[b].Lord
	if err := w.FormCorps(att, kinds, manned); err != nil {
		t.Fatal(err)
	}
	if err := w.FormCorps(def, kinds, manned); err != nil {
		t.Fatal(err)
	}

	// 守方待在自己的城裡；攻方從隔壁一格走進去。
	node := w.Corps[def].Node
	if w.Cities[node].Owner != b {
		w.Cities[node].Owner = b
	}
	c := &w.Corps[att]
	c.Node = node
	c.X, c.Y = w.Cities[node].X-1, w.Cities[node].Y
	c.TargetNode = node
	c.TargetX, c.TargetY = w.Cities[node].X, w.Cities[node].Y
	c.Timer = 1

	var got *CorpsEvent
	for i := 0; i < 64 && got == nil; i++ {
		for _, ev := range w.tickCorps(0, rng.New(0, 0, 0)) {
			if ev.Battle != nil {
				e := ev
				got = &e
				break
			}
		}
	}
	if got == nil {
		t.Fatal("走進敵方據點沒有打起來")
	}
	if got.Enemy != def {
		t.Errorf("對手是軍團 %d，應為 %d", got.Enemy, def)
	}
	if got.Mode != combat.Siege {
		t.Errorf("打成 %v，應為攻城——據點的判定被野戰搶先了", got.Mode)
	}
}

// TestCapitalPickMatchesScenarioData 拿四個劇本的初始資料當黃金對照：
// 把 `sub_16A3D` 照抄的選首都演算法跑過每一個活著的勢力，
// 與檔案裡的首都欄位（勢力記錄 +0x03）比。
//
// **命中 41/43，兩個例外是作者刻意填的**：
//   - 劇本 1 勢力 19：三個據點同為小城，演算法選生產力較高的「黃」，
//     作者填「北海」。
//   - 劇本 2 勢力 2：作者把首都填成**長阪**——那是戰場（類型 4，
//     城兵上限 0）。長阪坡的劉備。
//
// 首都是作者填的，這支只在首都易主時才跑，所以不必吻合。
// 這條測試釘住的是**那個數字**：如果它變了，代表選首都的邏輯或
// 據點類型的解讀被動過。
func TestCapitalPickMatchesScenarioData(t *testing.T) {
	if _, err := os.Stat(origPath); err != nil {
		t.Skip("找不到原版 SINARIO.DAT，跳過")
	}
	hit, total := 0, 0
	for idx := 0; idx < 4; idx++ {
		w, err := LoadScenario(origPath, idx)
		if err != nil {
			t.Fatalf("劇本 %d 載入失敗：%v", idx+1, err)
		}
		for f := range w.Factions {
			if !w.Factions[f].Alive {
				continue
			}
			total++
			if w.relocateCapital(f) == capital.None {
				hit++ // 回 None 代表「選出來的就是現在這個」
			}
		}
	}
	if total != 43 || hit != 41 {
		t.Fatalf("四個劇本應是 41/43 吻合，得到 %d/%d", hit, total)
	}
}

// TestCityKindDistribution 釘住據點類型的分佈。
// 三條獨立證據（城名、生產力上限、城兵上限）都指向這個讀法，
// 見 docs/formats/08 §1.6。
func TestCityKindDistribution(t *testing.T) {
	w := load(t, 0)
	var n [5]int
	for i := range w.Cities {
		c := &w.Cities[i]
		if c.Kind < 0 || c.Kind > 4 {
			t.Fatalf("據點 %d 的類型 %d 超出 0–4", i, c.Kind)
		}
		n[c.Kind]++
		// 戰場不能駐兵，關的守備最厚——類型解讀對不對，這兩條最敏感。
		if c.Kind == 4 && c.GarrisonCap != 0 {
			t.Errorf("戰場 %s 的城兵上限應是 0，得到 %d", c.Name, c.GarrisonCap)
		}
		if c.Kind == 3 && c.GarrisonCap < 194 {
			t.Errorf("關 %s 的城兵上限應 ≥ 194，得到 %d", c.Name, c.GarrisonCap)
		}
	}
	want := [5]int{21, 72, 86, 9, 4}
	if n != want {
		t.Fatalf("類型分佈 %v，期望 %v（大/中/小/關/戰場）", n, want)
	}
}

// ⭐ 行軍要沿著道路走，不是直線穿過山河。
//
// 用原版的地圖推出道路圖掛上去，然後叫一支軍團走到一個**不相鄰**的據點，
// 檢查它真的經過中繼據點。沒有這條測試的話，「Route 永遠是空的」
// 與「行軍正常」在畫面上看起來一樣——軍團照樣會抵達目的地。
func TestMarchFollowsRoads(t *testing.T) {
	raw, err := os.ReadFile("../../workplace/orig/dosv/MMAP.MAP")
	if err != nil {
		t.Skip("找不到原版 MMAP.MAP，跳過")
	}
	m, err := world.ParseMap(raw)
	if err != nil {
		t.Fatal(err)
	}
	w := load(t, 0)
	xy := make([][2]int, len(w.Cities))
	for i := range w.Cities {
		xy[i] = [2]int{w.Cities[i].X, w.Cities[i].Y}
	}
	edges, err := world.RoadEdges(m, xy)
	if err != nil {
		t.Fatal(err)
	}
	me := make([]march.Edge, len(edges))
	for i, e := range edges {
		me[i] = march.Edge{A: e.A, B: e.B, Steps: e.Steps,
			Path: e.Path, ACell: xy[e.A]}
	}
	g := march.New(len(w.Cities), me)
	w.SetRoads(g)

	// 找一對距離最遠的據點當測試對象——最遠的那一對必然要經過中繼點。
	from, to, hops := 0, 0, 0
	for a := 0; a < len(w.Cities); a += 17 { // 抽樣，不必全掃
		for b := 0; b < len(w.Cities); b += 17 {
			if p := g.Route(a, b); len(p) > hops {
				from, to, hops = a, b, len(p)
			}
		}
	}
	if hops < 3 {
		t.Fatalf("抽樣找不到需要中繼的路線（最長 %d 段）", hops)
	}

	f := w.AliveFactions()[0]
	lord := w.Factions[f].Lord
	kinds := [army.Positions]army.TroopType{
		army.Cavalry, army.Cavalry, army.Cavalry,
		army.Cavalry, army.Cavalry, army.Cavalry,
	}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
	if err := w.FormCorps(lord, kinds, manned); err != nil {
		t.Fatal(err)
	}
	c := &w.Corps[lord]
	c.Node, c.X, c.Y = from, w.Cities[from].X, w.Cities[from].Y

	if err := w.March(lord, to); err != nil {
		t.Fatal(err)
	}
	if len(w.routes[lord]) == 0 {
		t.Fatal("沒有算出格子路徑")
	}
	// ⭐ 路徑必須逐格連續。斷開的話軍團會瞬移，而「有沒有抵達」看不出來。
	prev := [2]int{c.X, c.Y}
	for k, cell := range w.routes[lord] {
		dx, dy := cell[0]-prev[0], cell[1]-prev[1]
		if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
			t.Fatalf("第 %d 格從 %v 跳到 %v", k, prev, cell)
		}
		prev = cell
	}

	// 走完全程，記下經過的據點與每一格的座標。
	seen := map[int]bool{}
	cells := 0
	for i := 0; i < 200000 && c.Node != to; i++ {
		w.step(lord)
		seen[c.Node] = true
		cells++
	}
	// 走的格數要與路徑長度相符 —— 差太多代表某處抄了近路。
	if cells < hops {
		t.Errorf("只走了 %d 格，路線有 %d 段，太短了", cells, hops)
	}
	if c.Node != to {
		t.Fatal("走不到目的地")
	}
	route := g.Route(from, to)
	for _, n := range route[1 : len(route)-1] {
		if !seen[n] {
			t.Errorf("沒有經過中繼據點 %s", big5Name(w.Cities[n].Name))
		}
	}
}

func big5Name(s string) string { return s }

// 走不到的目的地要回錯誤，不能默默走直線。
func TestMarchRejectsUnreachable(t *testing.T) {
	w := load(t, 0)
	// 只有一條邊 0–1 的圖：從 0 到 2 走不到。
	w.SetRoads(march.New(len(w.Cities), []march.Edge{{A: 0, B: 1, Steps: 1}}))
	f := w.AliveFactions()[0]
	lord := w.Factions[f].Lord
	kinds := [army.Positions]army.TroopType{
		army.Cavalry, army.Cavalry, army.Cavalry,
		army.Cavalry, army.Cavalry, army.Cavalry,
	}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 6000, 6000}
	if err := w.FormCorps(lord, kinds, manned); err != nil {
		t.Fatal(err)
	}
	w.Corps[lord].Node = 0
	if err := w.March(lord, 2); err == nil {
		t.Fatal("走不到卻沒有回錯誤")
	}
	if w.Corps[lord].TargetNode != 0 {
		t.Errorf("失敗後目標應留在原地，卻是 %d", w.Corps[lord].TargetNode)
	}
}
