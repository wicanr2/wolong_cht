package state

import (
	"os"
	"reflect"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/capital"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
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
	// 資金 3 B ＋ 士氣基準 1 ＋ 上昇值 1 ＋ 生產力 2 ＋ 月 1（u16 高位不變）
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

// 信賴度是原版 cs:0D00h（IDA `byte_10D00`）的單一 byte；改寫式存檔
// 必須讀回它，並在超出原版 byte 值域時採 0…255 的邊界。
func TestTrustStorage(t *testing.T) {
	w := load(t, 0)
	if w.Trust != int(w.raw[trustOffset]) {
		t.Fatalf("信賴度 = %d，原始 +0x%X = %d", w.Trust, trustOffset, w.raw[trustOffset])
	}
	for _, in := range []int{-1, 123, 300} {
		w.Trust = in
		got, err := reparse(w.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		want := in
		if want < 0 {
			want = 0
		}
		if want > 0xFF {
			want = 0xFF
		}
		if got.Trust != want {
			t.Errorf("信賴度輸入 %d → %d，want %d", in, got.Trust, want)
		}
	}
}

// Player 是原版全域區段的一組持久化值：`word_10CFD`（區塊 +0x0D）
// 是勢力表位址 faction×0x40，`byte_10CFF`（+0x0F）是勢力編號。
// 新劇本兩者都是無玩家哨兵；有效存檔載入時原版會跳過新遊戲選擇流程，
// 直接使用這組值。
func TestPlayerStorage(t *testing.T) {
	w := load(t, 0)
	if w.Player != -1 {
		t.Fatalf("新劇本 Player = %d，want -1", w.Player)
	}
	w.Player = 3
	b := w.Bytes()
	if got := u16(b, playerPtrOffset); got != 3*factionSize {
		t.Errorf("Player pointer = 0x%X，want 0x%X", got, 3*factionSize)
	}
	if got := int(b[playerOffset]); got != 3 {
		t.Errorf("Player byte = %d，want 3", got)
	}
	got, err := reparse(b)
	if err != nil {
		t.Fatal(err)
	}
	if got.Player != 3 {
		t.Errorf("Player round-trip = %d，want 3", got.Player)
	}

	// 指標與 byte 不一致時 fail-closed，不把一個孤立 byte 猜成玩家勢力。
	b[playerOffset] = 4
	got, err = reparse(b)
	if err != nil {
		t.Fatal(err)
	}
	if got.Player != -1 {
		t.Errorf("不一致的 Player 欄位被解成 %d，want -1", got.Player)
	}
}

// 事件佇列先驗收「原始資料可保存」，不把尚未完全解出的 handler 行為
// 偷換成猜測。Code 的高 byte 也要 round-trip，因為原版用它承載勢力／
// 災害變體（docs/formats/08、docs/re/07）；已知 handler 另由下方測試驗收。
func TestEventQueueStorage(t *testing.T) {
	w := load(t, 0)
	w.events[0] = QueuedEvent{Code: 0x010C, Param: 171}
	w.events[63] = QueuedEvent{Code: 0xFF01, Param: 0x1234}
	w.events[64] = QueuedEvent{Code: 0x020C, Param: 0x00C0}
	w.events[255] = QueuedEvent{Code: 0x000D, Param: 127}

	b := w.Bytes()
	for i, want := range []QueuedEvent{w.events[0], w.events[63], w.events[64], w.events[255]} {
		idx := []int{0, 63, 64, 255}[i]
		off := eventQueueOffset + idx*eventQueueEntrySize
		if got := u16(b, off); got != int(want.Code) {
			t.Errorf("佇列 %d code = 0x%04X, want 0x%04X", idx, got, want.Code)
		}
		if got := u16(b, off+2); got != int(want.Param) {
			t.Errorf("佇列 %d param = 0x%04X, want 0x%04X", idx, got, want.Param)
		}
	}

	got, err := reparse(b)
	if err != nil {
		t.Fatal(err)
	}
	for _, idx := range []int{0, 63, 64, 255} {
		if got.events[idx] != w.events[idx] {
			t.Errorf("佇列 %d round-trip = %#v, want %#v", idx, got.events[idx], w.events[idx])
		}
	}
}

// 原版 sub_131AE 每十次「每時」呼叫才取一格；新遊戲／月結後的第一筆
// 不是立刻處理，而是從 byte_131AD = 7 開始倒數。
func TestEventQueueTiming(t *testing.T) {
	w := load(t, 0)
	w.events[0] = QueuedEvent{Code: 0x010C, Param: 171}
	w.eventCursor, w.eventDelay = 0, 7

	for i := 0; i < 6; i++ {
		if _, ok := w.takeNextQueuedEvent(); ok {
			t.Fatalf("第 %d 次就取到事件", i+1)
		}
	}
	got, ok := w.takeNextQueuedEvent()
	if !ok || got != w.events[0] {
		t.Fatalf("第 7 次取到 %#v／%v，want %#v／true", got, ok, w.events[0])
	}
	if w.eventCursor != eventQueueEntrySize || w.eventDelay != 10 {
		t.Fatalf("取事件後 cursor=%d delay=%d，want %d／10",
			w.eventCursor, w.eventDelay, eventQueueEntrySize)
	}
	for i := 0; i < 9; i++ {
		if _, ok := w.takeNextQueuedEvent(); ok {
			t.Fatalf("事件後第 %d 次又取到事件", i+1)
		}
	}
	if w.eventDelay != 1 {
		t.Fatalf("下一筆前 delay=%d，want 1", w.eventDelay)
	}
}

// 無輸入主迴圈的 state 核心是連續呼叫 World.Tick：據點／軍團先更新，
// 時鐘再按子刻進位，每個新的「時」才進入 hourly／queue dispatcher。事件 10
// 不應成為 clock 的 driver；它只能在原版 byte_131AD 的第 7／之後每 10 次
// 每時節拍被取出。
func TestIdleClockDispatchesQueuedEvent10OnHourlyCadence(t *testing.T) {
	w := load(t, 0)
	w.events = [eventQueueEntries]QueuedEvent{{Code: 0x030A, Param: 0x0042}}
	w.eventCursor, w.eventDelay = 0, 7
	rng := rng.NewFixed(13)

	var hourly int
	var notices []TalkNotice
	for i := 0; i < 7*clock.SubticksPerHour; i++ {
		ev := w.Tick(rng)
		if !ev.Clock.Hour {
			continue
		}
		hourly++
		notices = append(notices, ev.TalkNotices...)
		if hourly < 7 && len(ev.TalkNotices) != 0 {
			t.Fatalf("第 %d 個每時邊界提前產生事件 10：%#v", hourly, ev.TalkNotices)
		}
	}
	if hourly != 7 {
		t.Fatalf("idle fixture 只跑到 %d 個每時邊界，want 7；clock=%+v", hourly, w.Clock)
	}
	if len(notices) != 1 || notices[0] != (TalkNotice{
		Index: 0x42, City: -1, Faction: -1, General: 3, Amount: -1,
	}) {
		t.Fatalf("第 7 個每時邊界的事件 10 = %#v，want 一筆 raw TALK notice", notices)
	}
	if w.Clock.Hour != 8 || w.Clock.Subtick != 0 {
		t.Fatalf("idle clock = %+v，want 從 1:0 前進 7 個每時邊界至 8:0", w.Clock)
	}
}

// sub_12BD9 的月度壓縮是「丟掉前 64 格、後 192 格前移、尾端清零」，
// 不是普通 FIFO 的逐筆刪除；這個差異會影響長期事件軌跡與存檔內容。
func TestEventQueueMonthlyCompaction(t *testing.T) {
	w := load(t, 0)
	for i := range w.events {
		w.events[i] = QueuedEvent{Code: uint16(i + 1), Param: uint16(0xA000 + i)}
	}
	w.eventCursor, w.eventDelay = 124, 2
	w.compactEventQueue()

	for i := 0; i < eventQueueEntries-eventQueueDispatch; i++ {
		want := QueuedEvent{Code: uint16(i + eventQueueDispatch + 1), Param: uint16(0xA000 + i + eventQueueDispatch)}
		if w.events[i] != want {
			t.Fatalf("壓縮後第 %d 格 = %#v，want %#v", i, w.events[i], want)
		}
	}
	for i := eventQueueEntries - eventQueueDispatch; i < eventQueueEntries; i++ {
		if w.events[i] != (QueuedEvent{}) {
			t.Fatalf("壓縮後尾端第 %d 格未清零：%#v", i, w.events[i])
		}
	}
	if w.eventCursor != 0 || w.eventDelay != 7 {
		t.Fatalf("壓縮後 cursor=%d delay=%d，want 0／7", w.eventCursor, w.eventDelay)
	}
}

// 目前已有獨立反組譯證據的事件 handler：事件 1 宣戰、事件 2 合作、
// 事件 3 停戰、事件 6／7 外交官回報、事件 8 遷都、事件 9 釋放指定武將、
// 事件 13 信賴度 −50。
// 未解的事件仍不應因 dispatch 而產生猜測效果。
func TestQueuedEventHandlers(t *testing.T) {
	w := load(t, 0)
	w.Player = -1
	source, target := 0, 13
	w.Factions[source].InvasionTarget = diplomacy.NoTarget
	w.Factions[target].InvasionTarget = diplomacy.NoTarget
	w.Friendship[source][target] = diplomacy.Peace(40)
	w.Friendship[target][source] = diplomacy.Peace(60)
	w.events[0] = QueuedEvent{
		Code:  uint16(source)<<8 | 1,
		Param: uint16(0xFF00 | target),
	}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if got := w.Factions[source].InvasionTarget; got != target {
		t.Fatalf("事件 1 沒有設定發起方目標：got %d want %d", got, target)
	}
	if got := w.Factions[target].InvasionTarget; got != source {
		t.Fatalf("事件 1 沒有套用回頭宣戰：got %d want %d", got, source)
	}
	if !w.Friendship[source][target].AtWar() || !w.Friendship[target][source].AtWar() {
		t.Fatal("事件 1 後雙向交友度仍是和平")
	}

	w = load(t, 0)
	player, invader, invaded := 0, 1, 2
	w.Player = player
	// 事件 2 呼叫 sub_13712(SI=player, DI=被侵攻方，BX=侵攻方)。
	// 讓 player 君主政治 12、被侵攻方最高政治未出陣武將政治 8，
	// 並避開平手亂數；合作金額因此是 (90−60)/2×1000=15000。
	w.Factions[player].Lord = 0
	w.Factions[player].Diplomat = noFaction
	w.Factions[player].Aggression = 1
	w.Factions[invaded].Diplomat = noFaction
	for i := range w.Generals {
		if w.Generals[i].Faction == invaded {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[0] = General{Alive: true, Politics: 12, Posted: true, Faction: player, Captor: noFaction}
	w.Generals[1] = General{Alive: true, Politics: 8, Faction: invaded, Captor: noFaction}
	w.Generals[2] = General{Alive: true, Faction: player, Captor: invaded, Posted: true}
	w.Corps[0] = Corps{Alive: true, Faction: player}
	w.Factions[player].Funds = 10_000
	w.Factions[invaded].Funds = 50_000
	w.Friendship[player][invader] = diplomacy.Peace(40)
	w.Friendship[player][invaded] = diplomacy.Peace(60)
	w.Friendship[invader][player] = diplomacy.Peace(50)
	w.Factions[invader].InvasionTarget = diplomacy.NoTarget
	w.events[0] = QueuedEvent{
		Code:  uint16(player)<<8 | 2,
		Param: uint16(invaded)<<8 | uint16(invader),
	}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if w.PendingDiplomacy() == nil || !w.ResolveDiplomacy(DiplomacyOfferFunds) {
		t.Fatal("事件 2 玩家合作選擇沒有完成狀態收尾")
	}
	if w.Factions[player].Funds != 25_000 || w.Factions[invaded].Funds != 35_000 {
		t.Fatalf("事件 2 金額方向／數值錯誤：player=%d invaded=%d", w.Factions[player].Funds, w.Factions[invaded].Funds)
	}
	if w.Factions[player].InvasionTarget != diplomacy.NoTarget || w.Factions[invader].InvasionTarget != player {
		t.Fatalf("事件 2 宣戰收尾錯誤：player=%d invader=%d", w.Factions[player].InvasionTarget, w.Factions[invader].InvasionTarget)
	}
	if !w.Friendship[player][invader].AtWar() || !w.Friendship[invader][player].AtWar() {
		t.Fatal("事件 2 合作成立後玩家與侵攻方仍是和平")
	}
	if g := w.Generals[2]; g.Faction != invaded || g.Captor != noFaction || g.Posted {
		t.Fatalf("事件 2 沒有釋放合作雙方俘虜：%+v", g)
	}

	w = load(t, 0)
	w.Player = -1
	source, target = 0, 1
	// 事件 3 呼叫 sub_136C4(SI=target, DI=source)。讓 target 的君主
	// 已出陣，這樣 sub_13771 會取君主政治；source 則提供 fallback
	// 的最高政治未出陣武將。數值故意避開平手亂數分支。
	targetLord, sourceGeneral := 0, 1
	w.Factions[target].Lord = targetLord
	w.Factions[target].Diplomat = noFaction
	w.Factions[target].Aggression = 1
	w.Factions[source].Diplomat = noFaction
	for i := range w.Generals {
		if w.Generals[i].Faction == source {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[targetLord] = General{Alive: true, Politics: 12, Posted: true, Faction: target, Captor: noFaction}
	w.Generals[sourceGeneral] = General{Alive: true, Politics: 8, Faction: source, Captor: noFaction}
	w.Corps[targetLord] = Corps{Alive: true, Faction: target}
	w.Friendship[target][source] = diplomacy.Peace(20)
	w.Friendship[source][target] = diplomacy.Peace(50)
	w.Factions[source].Funds = 50_000
	w.Factions[target].Funds = 10_000
	w.Factions[source].InvasionTarget = target
	w.Factions[target].InvasionTarget = 9
	w.Generals[2] = General{Alive: true, Faction: source, Captor: target, Posted: true}
	w.Generals[3] = General{Alive: true, Faction: target, Captor: source, Posted: true}
	w.Generals[4] = General{Alive: true, Faction: source, Captor: 5, Posted: true}
	w.Cities[0].Owner, w.Cities[0].OwnerRecorded = source, target
	w.Cities[1].Owner, w.Cities[1].OwnerRecorded = source, 5
	w.events[0] = QueuedEvent{
		Code:  uint16(source)<<8 | 3,
		Param: uint16(0xFF00 | target),
	}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if got := w.Factions[target].Funds; got != 28_000 {
		t.Fatalf("事件 3 對方資金 = %d，want 28000", got)
	}
	if got := w.Factions[source].Funds; got != 32_000 {
		t.Fatalf("事件 3 提出方資金 = %d，want 32000", got)
	}
	if w.Factions[source].InvasionTarget != diplomacy.NoTarget || w.Factions[target].InvasionTarget != 9 {
		t.Fatalf("事件 3 侵攻目標清理錯誤：source=%d target=%d", w.Factions[source].InvasionTarget, w.Factions[target].InvasionTarget)
	}
	if got := w.Friendship[source][target]; got != diplomacy.Peace(20) || w.Friendship[target][source] != diplomacy.Peace(20) {
		t.Fatalf("事件 3 交友度 = %#v／%#v，want 雙向和平 20", got, w.Friendship[target][source])
	}
	if g := w.Generals[2]; g.Faction != target || g.Captor != noFaction || g.Posted {
		t.Fatalf("事件 3 沒有釋放 source→target 俘虜：%+v", g)
	}
	if g := w.Generals[3]; g.Faction != source || g.Captor != noFaction || g.Posted {
		t.Fatalf("事件 3 沒有釋放 target→source 俘虜：%+v", g)
	}
	if g := w.Generals[4]; g.Faction != source || g.Captor != 5 || !g.Posted {
		t.Fatalf("事件 3 誤改非配對俘虜：%+v", g)
	}
	if w.Cities[0].OwnerRecorded != source || w.Cities[1].OwnerRecorded != 5 {
		t.Fatalf("事件 3 OwnerRecorded 同步邊界錯誤：%d／%d", w.Cities[0].OwnerRecorded, w.Cities[1].OwnerRecorded)
	}

	w = load(t, 0)
	w.Player = -1
	faction, expected := -1, -1
	for i, f := range w.Factions {
		if !f.Alive {
			continue
		}
		old := f.Capital
		if next := w.relocateCapital(i); next != capital.None {
			faction, expected = i, next
			w.Factions[i].Capital = old
			break
		}
	}
	if faction < 0 {
		t.Skip("劇本 1 沒有可供事件 8 驗證的遷都候選")
	}
	oldCapital := w.Factions[faction].Capital
	otherFaction := (faction + 1) % numFactions
	w.Corps[0] = Corps{
		Alive: true, Faction: faction, Ordered: oldCapital, TargetNode: expected,
		TargetX: 0x1234, TargetY: 0x5678,
	}
	w.Corps[1] = Corps{
		Alive: true, Faction: faction, Ordered: oldCapital, TargetNode: oldCapital,
	}
	w.Corps[2] = Corps{
		Alive: true, Faction: otherFaction, Ordered: oldCapital, TargetNode: expected,
	}
	w.events[0] = QueuedEvent{Code: uint16(faction)<<8 | 8, Param: 0xFFFF}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if got := w.Factions[faction].Capital; got != expected {
		t.Fatalf("事件 8 首都 = %d，want %d", got, expected)
	}
	if c := w.Corps[0]; c.Ordered != expected || c.TargetNode != oldCapital || c.TargetX != 0x1234 || c.TargetY != 0x5678 {
		t.Fatalf("事件 8 未照 sub_14502 同步軍團 0：%+v，舊首都=%d 新首都=%d", c, oldCapital, expected)
	}
	if c := w.Corps[1]; c.Ordered != expected || c.TargetNode != oldCapital {
		t.Fatalf("事件 8 不應改寫非新首都目標：%+v", c)
	}
	if c := w.Corps[2]; c.Ordered != oldCapital || c.TargetNode != expected {
		t.Fatalf("事件 8 不應同步其他勢力軍團：%+v", c)
	}

	w = load(t, 0)
	w.Trust = 200
	w.events[0] = QueuedEvent{Code: 13, Param: 0x0196}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if w.Trust != 150 {
		t.Fatalf("事件 13 Trust = %d，want 150", w.Trust)
	}
}

// 玩家進言的三個 producer 必須保留原版的分流：敵對直接收尾，停戰／協力
// 寫入事件 6／7；後兩者從第 20 格提示位置找完整 256 格，而不是套用 AI
// producer 只掃前 64 格的限制。
func TestPlayerDiplomacyProducers(t *testing.T) {
	w := load(t, 0)
	player, target := 0, 1
	w.Player = player
	w.events = [eventQueueEntries]QueuedEvent{}
	w.Factions[target].Diplomat = noFaction
	w.Friendship[player][target] = diplomacy.Peace(40).WithWar(true)
	for i := 0x14; i < eventQueueDispatch; i++ {
		w.events[i] = QueuedEvent{Code: 0x010C}
	}
	w.eventCursor, w.eventDelay = 0, 7
	if !w.QueuePlayerCeasefire(target) {
		t.Fatal("玩家停戰 producer 沒有寫入事件 6")
	}
	if got := w.events[eventQueueDispatch]; got != (QueuedEvent{Code: uint16(target)<<8 | 6}) {
		t.Fatalf("事件 6 payload／完整佇列位置錯誤：%#v", got)
	}
	if w.QueuePlayerCeasefire(target) {
		t.Fatal("同一回報方已有事件 6 時不應重複排入")
	}

	w = load(t, 0)
	w.Player = player
	w.events = [eventQueueEntries]QueuedEvent{}
	w.Factions[player].InvasionTarget = diplomacy.NoTarget
	w.Friendship[player][target] = diplomacy.Peace(40)
	if !w.ApplyPlayerHostility(target) {
		t.Fatal("玩家敵對提案沒有走直接宣戰收尾")
	}
	if !w.Friendship[player][target].AtWar() || !w.Friendship[target][player].AtWar() {
		t.Fatal("玩家敵對提案後雙向交友度仍是和平")
	}

	w = load(t, 0)
	ally, invader := 1, 2
	w.Player = player
	w.events = [eventQueueEntries]QueuedEvent{}
	w.Factions[ally].Diplomat = noFaction
	w.Friendship[player][invader] = diplomacy.Peace(40).WithWar(true)
	if !w.QueuePlayerCooperation(ally, invader) {
		t.Fatal("玩家協力 producer 沒有寫入事件 7")
	}
	want := QueuedEvent{
		Code:  uint16(ally)<<8 | 7,
		Param: uint16(invader)<<8 | uint16(invader),
	}
	if got := w.events[0x14]; got != want {
		t.Fatalf("事件 7 payload／位置錯誤：%#v，want %#v", got, want)
	}
	if w.QueuePlayerCooperation(ally, invader) {
		t.Fatal("同一協力方已有事件 7 時不應重複排入")
	}
}

func TestEvent10ProducerWritesRawTalkPayload(t *testing.T) {
	w := load(t, 0)
	w.events = [eventQueueEntries]QueuedEvent{}
	w.eventCursor = 0
	if !w.QueueEvent10(3, 0x42) {
		t.Fatal("事件 10 producer 沒有寫入有效 raw payload")
	}
	if got := w.events[0]; got != (QueuedEvent{Code: 0x030A, Param: 0x0042}) {
		t.Fatalf("事件 10 raw payload = %#v，want code=030A param=0042", got)
	}

	// producer 使用完整 256 格路徑；前 64 格滿時仍應能保留訊息事件。
	w = load(t, 0)
	w.events = [eventQueueEntries]QueuedEvent{}
	for i := 0; i < eventQueueDispatch; i++ {
		w.events[i] = QueuedEvent{Code: 0x010C}
	}
	w.eventCursor = 0
	if !w.QueueEvent10(4, 0x43) {
		t.Fatal("事件 10 producer 不應被前 64 格滿阻擋")
	}
	if got := w.events[eventQueueDispatch]; got != (QueuedEvent{Code: 0x040A, Param: 0x0043}) {
		t.Fatalf("事件 10 完整佇列位置／payload = %#v", got)
	}

	if w.QueueEvent10(-1, 0x42) || w.QueueEvent10(numGenerals, 0x42) ||
		w.QueueEvent10(3, -1) || w.QueueEvent10(3, 0x10000) {
		t.Fatal("事件 10 producer 應拒絕越界 General／TALK index")
	}
}

// 事件 11／12 的 handler 寫入城市 runtime 的 +0x15 marker；sub_14269
// 在後續據點輪轉才套用持久傷害。事件 12 的高 byte 與 Param 位址、延遲
// 清除事件也要保留。
func TestQueuedDisasterAnimationHandlers(t *testing.T) {
	w := load(t, 0)
	w.events = [eventQueueEntries]QueuedEvent{}
	w.rng = rng.NewFixed(17)
	param := uint16(runtimeCityBase)
	w.events[0] = QueuedEvent{Code: 0x010C, Param: param}
	w.eventCursor, w.eventDelay = 0, 1
	ev := &Event{}
	w.dispatchQueuedEvent(ev)
	if w.disasterMarkers[0] != economy.Fire || w.disasterMarkerLevels[0] == 0 {
		t.Fatalf("事件 0x010C 沒有建立火災 runtime 標記：kind=%v level=%d",
			w.disasterMarkers[0], w.disasterMarkerLevels[0])
	}
	if ev.Disaster[0] != economy.Fire {
		t.Fatalf("事件 0x010C 沒有回報火災據點：%v", ev.Disaster)
	}
	clearSlot := -1
	for i, e := range w.events {
		if e.Code == 0x000C && e.Param == param {
			clearSlot = i
			break
		}
	}
	if clearSlot < 0 {
		t.Fatal("事件 0x010C 沒有排入延遲清除事件")
	}
	w.eventCursor, w.eventDelay = clearSlot*eventQueueEntrySize, 1
	w.dispatchQueuedEvent(&Event{})
	if w.disasterMarkers[0] != economy.NoDisaster || w.disasterMarkerLevels[0] != 0 {
		t.Fatalf("事件 0x000C 沒有清除 runtime 標記：kind=%v level=%d",
			w.disasterMarkers[0], w.disasterMarkerLevels[0])
	}

	w = load(t, 0)
	w.events = [eventQueueEntries]QueuedEvent{}
	w.rng = rng.NewFixed(23)
	c := w.Cities[0]
	w.stormArea = &economy.StormArea{
		MinX: c.X - 5, MinY: c.Y - 5, MaxX: c.X + 5, MaxY: c.Y + 5,
	}
	w.events[0] = QueuedEvent{Code: 0x000B}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if w.disasterMarkers[0] != economy.Storm || w.disasterMarkerLevels[0] == 0 {
		t.Fatalf("事件 0x000B 沒有建立暴風雨 runtime 標記：kind=%v level=%d",
			w.disasterMarkers[0], w.disasterMarkerLevels[0])
	}
}

func TestDisasterMarkerAppliesRawPersistentEffects(t *testing.T) {
	w := load(t, 0)
	c := &w.Cities[0]
	c.Prevention = 3
	c.Growth = 0 // 原始 +0x10 存值 100
	c.Production = 0x1200
	c.Garrison = 20
	w.disasterMarkers[0] = economy.Fire
	w.disasterMarkerLevels[0] = 7

	w.applyCityDisasterEffect(0)

	// deficit = 7 − 3 = 4；生產力損失 = (4 × 0x12) >> 2 = 18。
	if c.Prevention != 0 || c.Growth != -4 || c.Production != 0x11EE || c.Garrison != 18 {
		t.Fatalf("sub_14269 持久效果錯誤：prevention=%d growth=%d production=%#x garrison=%d",
			c.Prevention, c.Growth, c.Production, c.Garrison)
	}
	if w.disasterMarkers[0] != economy.Fire || w.disasterMarkerLevels[0] != 7 {
		t.Fatal("sub_14269 不應自行清除 +0x15 marker")
	}

	// 防災值足夠時只扣護盾，不應碰其他三個欄位。
	c.Prevention = 10
	c.Growth = 12
	c.Production = 0x2345
	c.Garrison = 19
	w.applyCityDisasterEffect(0)
	if c.Prevention != 3 || c.Growth != 12 || c.Production != 0x2345 || c.Garrison != 19 {
		t.Fatalf("防災護盾分支錯誤：prevention=%d growth=%d production=%#x garrison=%d",
			c.Prevention, c.Growth, c.Production, c.Garrison)
	}
}

func TestDisasterMarkerReadOnlySnapshots(t *testing.T) {
	w := load(t, 0)
	w.disasterMarkers[0] = economy.Fire
	w.disasterMarkerLevels[0] = 7

	got, ok := w.DisasterMarkerAt(0)
	if !ok || got != (DisasterMarker{Kind: economy.Fire, Level: 7}) {
		t.Fatalf("火災 marker snapshot 錯誤：%#v，ok=%v", got, ok)
	}
	if _, ok := w.DisasterMarkerAt(-1); ok {
		t.Fatal("負的據點編號不應回傳 marker")
	}
	if _, ok := w.DisasterMarkerAt(len(w.Cities)); ok {
		t.Fatal("超出據點表不應回傳 marker")
	}

	w.stormArea = &economy.StormArea{MinX: 10, MinY: 20, MaxX: 20, MaxY: 30}
	area, ok := w.StormAreaSnapshot()
	if !ok || area != (economy.StormArea{MinX: 10, MinY: 20, MaxX: 20, MaxY: 30}) {
		t.Fatalf("暴風雨範圍 snapshot 錯誤：%#v，ok=%v", area, ok)
	}
	area.MinX = 999
	areaAgain, _ := w.StormAreaSnapshot()
	if areaAgain.MinX != 10 {
		t.Fatal("StormAreaSnapshot 不應讓呼叫端改到 runtime 範圍")
	}
}

// 事件 11／12／13 的通知索引是原版 handler 直接傳入 sub_18810 的值：
// #70 暴風雨、#71 大火、#72 暴動、#51 君主赤字警告。state 只回傳索引
// 與城市目標，文字與 Big5 展開留給呈現層。
func TestQueuedTalkNotices(t *testing.T) {
	w := load(t, 0)
	player := 0
	w.Player = player
	city := w.Factions[player].Capital
	if city < 0 || city >= len(w.Cities) || w.Cities[city].Owner != player {
		t.Fatalf("玩家首都不是有效的玩家據點：player=%d city=%d", player, city)
	}

	param := uint16(runtimeCityBase + city*citySize)
	w.events = [eventQueueEntries]QueuedEvent{{Code: 0x010C, Param: param}}
	w.eventCursor, w.eventDelay = 0, 1
	fire := &Event{}
	w.dispatchQueuedEvent(fire)
	if len(fire.TalkNotices) != 1 || fire.TalkNotices[0] != (TalkNotice{
		Index: 0x47, City: city, Faction: -1, General: -1, Amount: -1,
	}) {
		t.Fatalf("事件 12 火災通知錯誤：%#v", fire.TalkNotices)
	}

	w = load(t, 0)
	w.Player = player
	c := w.Cities[city]
	w.rng = rng.NewFixed(23)
	w.stormArea = &economy.StormArea{
		MinX: c.X - 5, MinY: c.Y - 5, MaxX: c.X + 5, MaxY: c.Y + 5,
	}
	w.events = [eventQueueEntries]QueuedEvent{{Code: 0x000B}}
	w.eventCursor, w.eventDelay = 0, 1
	storm := &Event{}
	w.dispatchQueuedEvent(storm)
	foundStorm := false
	for _, notice := range storm.TalkNotices {
		if notice.Index == 0x46 && notice.City == city && notice.General == -1 {
			foundStorm = true
			break
		}
	}
	if !foundStorm {
		t.Fatalf("事件 11 沒有產生玩家據點 #70 通知：%#v", storm.TalkNotices)
	}

	w = load(t, 0)
	w.events = [eventQueueEntries]QueuedEvent{{Code: 0x000D, Param: 0x0196}}
	w.eventCursor, w.eventDelay = 0, 1
	deficit := &Event{}
	w.dispatchQueuedEvent(deficit)
	if len(deficit.TalkNotices) != 1 || deficit.TalkNotices[0] != (TalkNotice{
		Index: 0x33, City: -1, Faction: -1, General: -1, Amount: -1,
	}) {
		t.Fatalf("事件 13 赤字通知錯誤：%#v", deficit.TalkNotices)
	}
}

// 事件 10（sub_13496）不是「只取出不顯示」：Param 直接是 TALK.DAT
// 索引，高位元組以 FFxx formatter word 提供 \1 的武將索引。
func TestQueuedEvent10TalkNotice(t *testing.T) {
	w := load(t, 0)
	w.events[0] = QueuedEvent{Code: 0x030A, Param: 0x0042}
	w.eventCursor, w.eventDelay = 0, 1
	ev := &Event{}
	w.dispatchQueuedEvent(ev)
	if len(ev.TalkNotices) != 1 || ev.TalkNotices[0] != (TalkNotice{
		Index: 0x42, City: -1, Faction: -1, General: 3, Amount: -1,
	}) {
		t.Fatalf("事件 10 formatter 邊界錯誤：%#v", ev.TalkNotices)
	}

	// 原版仍會呼叫 TALK；來源超出武將表時只拒絕錯誤肖像／\1 代入，
	// 不把整則訊息吞掉。
	w = load(t, 0)
	w.events[0] = QueuedEvent{Code: 0xFF0A, Param: 0x0042}
	w.eventCursor, w.eventDelay = 0, 1
	ev = &Event{}
	w.dispatchQueuedEvent(ev)
	if len(ev.TalkNotices) != 1 || ev.TalkNotices[0].General != -1 ||
		ev.TalkNotices[0].Index != 0x42 {
		t.Fatalf("事件 10 無效 general 應 fail-closed 保留 TALK：%#v", ev.TalkNotices)
	}
}

// sub_123FF／sub_12459／sub_12533 的 runtime 物件時序：建立時 phase=1、
// timer=1；第一次 map update 立即 dirty，但該次 render 仍畫 phase=1，
// 之後每 16 次 map update 才換下一個相位。
func TestDisasterObjectAnimationTiming(t *testing.T) {
	w := load(t, 0)
	w.events = [eventQueueEntries]QueuedEvent{{Code: 0x010C, Param: runtimeCityBase}}
	w.eventCursor, w.eventDelay = 0, 1
	w.rng = rng.NewFixed(17)
	w.dispatchQueuedEvent(&Event{})

	objects := w.RenderDisasterObjects()
	if len(objects) != 1 || objects[0].TypeCode != 1 || objects[0].Phase != 1 {
		t.Fatalf("火災物件初始記錄錯誤：%#v", objects)
	}

	w.AdvanceDisasterObjects()
	objects = w.RenderDisasterObjects()
	if len(objects) != 1 || objects[0].Phase != 1 {
		t.Fatalf("第一次 dirty render 應先畫 phase=1：%#v", objects)
	}

	for i := 0; i < disasterObjectInterval-1; i++ {
		w.AdvanceDisasterObjects()
	}
	objects = w.RenderDisasterObjects()
	if len(objects) != 1 || objects[0].Phase != 2 {
		t.Fatalf("16 次 map update 後應畫 phase=2：%#v", objects)
	}

	// sub_12438 會清掉同城市的所有 runtime object，與 marker 清除同步。
	w.events[1] = QueuedEvent{Code: 0x000C, Param: runtimeCityBase}
	w.eventCursor, w.eventDelay = eventQueueEntrySize, 1
	w.dispatchQueuedEvent(&Event{})
	if got := w.RenderDisasterObjects(); len(got) != 0 {
		t.Fatalf("清除事件後仍有災害物件：%#v", got)
	}
}

func TestSub124FFMatchesRawSignedByteContract(t *testing.T) {
	tests := []struct {
		name      string
		drift     uint16
		random    int
		wantDrift uint16
		wantWhole int
	}{
		{name: "carry positive", drift: 0x000F, random: 7, wantDrift: 0x0000, wantWhole: 1},
		{name: "negative wrap", drift: 0xF100, random: 0, wantDrift: 0xF100, wantWhole: -1},
		{name: "signed fractional no carry", drift: 0x0FF0, random: 7, wantDrift: 0x0FFF, wantWhole: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDrift, gotWhole := sub124FF(tt.drift, tt.random)
			if gotDrift != tt.wantDrift || gotWhole != tt.wantWhole {
				t.Fatalf("sub_124FF(%04X,%d) = (%04X,%d), want (%04X,%d)",
					tt.drift, tt.random, gotDrift, gotWhole, tt.wantDrift, tt.wantWhole)
			}
		})
	}
}

type sequenceRand struct {
	values []int
	pos    int
}

func (r *sequenceRand) Next() int {
	if r == nil || len(r.values) == 0 {
		return 0
	}
	v := r.values[r.pos%len(r.values)]
	r.pos++
	return v
}

func TestMovingDisasterSub1248AUsesOnlyLastHalfOfRawSlots(t *testing.T) {
	w := load(t, 0)
	r := &sequenceRand{values: []int{7, 7}}
	base := disasterObject{
		active: true, kind: economy.Fire, typeCode: 1, city: 0,
		x: 10, y: 20, timer: 1, interval: disasterObjectInterval,
		phase: 1, xDrift: 0x000F, yDrift: 0x000F,
	}
	w.disasterObjects[disasterMovingSlot-1] = base
	w.disasterObjects[disasterMovingSlot] = base

	w.AdvanceMapObjects(r)
	if got := w.disasterObjects[disasterMovingSlot-1]; got.x != 10 ||
		got.y != 20 || got.xDrift != 0x000F || got.yDrift != 0x000F {
		t.Fatalf("sub_1248A 不應作用於前半 slots：%+v", got)
	}
	got := w.disasterObjects[disasterMovingSlot]
	if got.x != 11 || got.y != 21 || got.xDrift != 0 || got.yDrift != 0 ||
		got.timer != disasterObjectInterval || !got.dirty {
		t.Fatalf("後半 slot 未依 raw sub_1248A 移動：%+v", got)
	}
}

func TestMovingDisasterSub1248ARawWrapAndDirectionByte(t *testing.T) {
	w := load(t, 0)
	w.rng = &sequenceRand{values: []int{7, 7}}
	w.stormArea = &economy.StormArea{MinX: 5, MinY: 5, MaxX: 6, MaxY: 6}
	w.disasterObjects[disasterMovingSlot] = disasterObject{
		active: true, kind: economy.Riot, typeCode: 2, city: 0,
		x: disasterWrapMaxX + 1, y: disasterWrapMaxY + 1,
		timer: 1, interval: disasterObjectInterval,
		phase: 1, xDrift: 0x000F, yDrift: 0x000F,
	}

	w.AdvanceDisasterObjects()
	got := w.disasterObjects[disasterMovingSlot]
	if got.x != disasterWrapMinX || got.y != disasterWrapMinY ||
		got.xDrift != 0xFF00 || got.yDrift != 0xFF00 {
		t.Fatalf("sub_1248A 回繞／方向 byte 錯誤：%+v", got)
	}
}

// 事件 6／7 的 TALK 參數來自 handler 的同一組 raw 暫存器：#57 使用回報方
// 的勢力／外交官，#44／#48 的 \7 使用 sub_136C4／sub_13712 留在 DX 的金額。
func TestQueuedDiplomacyReportTalkNotices(t *testing.T) {
	w := load(t, 0)
	player, other := 0, 1
	w.Player = player
	w.Factions[other].Lord = 1
	w.Factions[other].Diplomat = 1
	w.Factions[other].Aggression = 1
	for i := range w.Generals {
		if w.Generals[i].Faction == other {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[1] = General{Alive: true, Politics: 8, Faction: other, Captor: noFaction}
	w.Factions[player].Funds = 50_000
	w.Factions[other].Funds = 10_000
	w.Factions[player].InvasionTarget = other
	w.Factions[other].InvasionTarget = diplomacy.NoTarget
	w.Friendship[player][other] = diplomacy.Peace(50)
	w.Friendship[other][player] = diplomacy.Peace(20)
	w.events[0] = QueuedEvent{Code: uint16(other)<<8 | 6}
	w.eventCursor, w.eventDelay = 0, 1
	event6 := &Event{}
	w.dispatchQueuedEvent(event6)
	want6 := []TalkNotice{
		{Index: 0x39, City: -1, Faction: other, General: 1, Amount: -1},
		{Index: 0x2C, City: -1, Faction: other, General: -1, Amount: 14_000},
	}
	if !reflect.DeepEqual(event6.TalkNotices, want6) {
		t.Fatalf("事件 6 TALK 通知錯誤：got %#v want %#v", event6.TalkNotices, want6)
	}

	w = load(t, 0)
	ally, invader := 1, 2
	w.Player = player
	w.Factions[ally].Lord = 1
	w.Factions[ally].Diplomat = 1
	w.Factions[ally].Aggression = 1
	for i := range w.Generals {
		if w.Generals[i].Faction == ally {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[1] = General{Alive: true, Politics: 8, Faction: ally, Captor: noFaction}
	w.Factions[player].Funds = 50_000
	w.Factions[ally].Funds = 10_000
	w.Factions[ally].InvasionTarget = diplomacy.NoTarget
	w.Factions[invader].InvasionTarget = diplomacy.NoTarget
	w.Friendship[ally][player] = diplomacy.Peace(60)
	w.Friendship[ally][invader] = diplomacy.Peace(20)
	w.Friendship[player][ally] = diplomacy.Peace(50)
	w.Friendship[player][invader] = diplomacy.Peace(50)
	w.events[0] = QueuedEvent{
		Code:  uint16(ally)<<8 | 7,
		Param: uint16(invader)<<8 | uint16(invader),
	}
	w.eventCursor, w.eventDelay = 0, 1
	event7 := &Event{}
	w.dispatchQueuedEvent(event7)
	want7 := []TalkNotice{
		{Index: 0x39, City: -1, Faction: ally, General: 1, Amount: -1},
		{Index: 0x30, City: -1, Faction: ally, General: -1, Amount: 15_000},
	}
	if !reflect.DeepEqual(event7.TalkNotices, want7) {
		t.Fatalf("事件 7 TALK 通知錯誤：got %#v want %#v", event7.TalkNotices, want7)
	}
}

func TestQueuedDiplomacySecondaryTalkConditions(t *testing.T) {
	w := load(t, 0)
	player, other := 0, 1
	w.Player = player
	w.Factions[other].Lord = 1
	w.Factions[other].Diplomat = 1
	w.Factions[other].Aggression = 1
	w.Generals[1] = General{Alive: true, Politics: 8, Faction: other, Captor: noFaction}
	// sub_13138 的第一個方向：回報方曾俘虜玩家方武將。
	w.Generals[2] = General{Faction: player, Captor: other}
	w.Factions[player].Funds = 50_000
	w.Factions[other].Funds = 10_000
	w.Factions[player].InvasionTarget = other
	w.Factions[other].InvasionTarget = diplomacy.NoTarget
	w.Friendship[player][other] = diplomacy.Peace(50)
	w.Friendship[other][player] = diplomacy.Peace(20)
	w.events[0] = QueuedEvent{Code: uint16(other)<<8 | 6}
	w.eventCursor, w.eventDelay = 0, 1
	ev := &Event{}
	w.dispatchQueuedEvent(ev)
	if len(ev.TalkNotices) != 3 {
		t.Fatalf("事件 6 有俘虜關係時應有第二次 TALK：%#v", ev.TalkNotices)
	}
	got := ev.TalkNotices[2]
	if got.Index != 0x48 || !got.Secondary || got.NoPortrait {
		t.Fatalf("事件 6 次要 TALK raw 邊界錯誤：%#v", got)
	}
	if got.RawFormatterWordValid || got.RawFormatterWord != -1 {
		t.Fatalf("事件 6 次要 TALK 未捕捉到原版 SS:[DI] payload 時必須 fail-closed：%#v", got)
	}

	w = load(t, 0)
	ally, invader := 1, 2
	w.Player = player
	w.Factions[ally].Lord = 1
	w.Factions[ally].Diplomat = 1
	w.Factions[ally].Aggression = 1
	w.Generals[1] = General{Alive: true, Politics: 8, Faction: ally, Captor: noFaction}
	w.Generals[2] = General{Faction: player, Captor: ally}
	w.Factions[player].Funds = 50_000
	w.Factions[ally].Funds = 10_000
	w.Factions[ally].InvasionTarget = diplomacy.NoTarget
	w.Factions[invader].InvasionTarget = diplomacy.NoTarget
	w.Friendship[ally][player] = diplomacy.Peace(60)
	w.Friendship[ally][invader] = diplomacy.Peace(20)
	w.Friendship[player][ally] = diplomacy.Peace(50)
	w.Friendship[player][invader] = diplomacy.Peace(50)
	w.events[0] = QueuedEvent{
		Code:  uint16(ally)<<8 | 7,
		Param: uint16(invader)<<8 | uint16(invader),
	}
	w.eventCursor, w.eventDelay = 0, 1
	ev = &Event{}
	w.dispatchQueuedEvent(ev)
	if len(ev.TalkNotices) != 3 {
		t.Fatalf("事件 7 有俘虜關係時應有第二次 TALK：%#v", ev.TalkNotices)
	}
	if got := ev.TalkNotices[2]; got != (TalkNotice{
		Index: 0x4C, City: -1, Faction: -1, General: -1, Amount: -1,
		RawFormatterWord: -1,
		Secondary:        true, NoPortrait: true,
	}) {
		t.Fatalf("事件 7 次要 TALK raw 邊界錯誤：%#v", got)
	}
}

// 事件 6／7 是玩家進言後的延遲外交回報，不是事件 2／3 的同一方向：
// 事件 6 由玩家付停戰金額給回報方；事件 7 由玩家付協力金額給協力方，
// 再由協力方對第三方宣戰。兩個 payload 與原版 handler 的 SI／DI／BX
// 方向都在這裡固定，避免日後把事件字高誤當成付款方。
func TestQueuedDiplomacyReportHandlers(t *testing.T) {
	w := load(t, 0)
	player, other := 0, 1
	w.Player = player
	w.Factions[other].Lord = 1
	w.Factions[other].Diplomat = 1
	w.Factions[other].Aggression = 1
	for i := range w.Generals {
		if w.Generals[i].Faction == other {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[1] = General{
		Alive: true, Politics: 8, Faction: other, Captor: noFaction,
	}
	w.Factions[player].Funds = 50_000
	w.Factions[other].Funds = 10_000
	w.Factions[player].InvasionTarget = other
	w.Factions[other].InvasionTarget = diplomacy.NoTarget
	w.Friendship[player][other] = diplomacy.Peace(50)
	w.Friendship[other][player] = diplomacy.Peace(20)
	w.events[0] = QueuedEvent{Code: uint16(other)<<8 | 6}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if got := w.Factions[player].Funds; got != 36_000 {
		t.Fatalf("事件 6 玩家付款後資金 = %d，want 36000", got)
	}
	if got := w.Factions[other].Funds; got != 24_000 {
		t.Fatalf("事件 6 回報方收款後資金 = %d，want 24000", got)
	}
	if w.Factions[player].InvasionTarget != diplomacy.NoTarget ||
		w.Factions[other].InvasionTarget != diplomacy.NoTarget {
		t.Fatalf("事件 6 停戰後侵攻目標未清除：player=%d other=%d",
			w.Factions[player].InvasionTarget, w.Factions[other].InvasionTarget)
	}
	if w.Friendship[player][other].AtWar() || w.Friendship[other][player].AtWar() {
		t.Fatal("事件 6 停戰回報後雙向交友度仍非和平")
	}

	w = load(t, 0)
	player, ally, invader := 0, 1, 2
	w.Player = player
	w.Factions[ally].Lord = 1
	w.Factions[ally].Diplomat = 1
	w.Factions[ally].Aggression = 1
	for i := range w.Generals {
		if w.Generals[i].Faction == ally {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[1] = General{
		Alive: true, Politics: 8, Faction: ally, Captor: noFaction,
	}
	w.Factions[player].Funds = 50_000
	w.Factions[ally].Funds = 10_000
	w.Factions[ally].InvasionTarget = diplomacy.NoTarget
	w.Factions[invader].InvasionTarget = diplomacy.NoTarget
	w.Friendship[ally][player] = diplomacy.Peace(60)
	w.Friendship[ally][invader] = diplomacy.Peace(20)
	w.Friendship[player][ally] = diplomacy.Peace(50)
	w.Friendship[player][invader] = diplomacy.Peace(50)
	w.events[0] = QueuedEvent{
		Code:  uint16(ally)<<8 | 7,
		Param: uint16(invader)<<8 | uint16(invader),
	}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if got := w.Factions[player].Funds; got != 35_000 {
		t.Fatalf("事件 7 玩家付款後資金 = %d，want 35000", got)
	}
	if got := w.Factions[ally].Funds; got != 25_000 {
		t.Fatalf("事件 7 協力方收款後資金 = %d，want 25000", got)
	}
	if w.Factions[ally].InvasionTarget != invader {
		t.Fatalf("事件 7 協力方未對侵攻目標宣戰：got %d want %d",
			w.Factions[ally].InvasionTarget, invader)
	}
	if !w.Friendship[ally][invader].AtWar() || !w.Friendship[invader][ally].AtWar() {
		t.Fatal("事件 7 協力成立後雙向交友度仍是和平")
	}

	// 回報方外交官在事件抵達前消失時，handler 必須保持 fail-closed。
	w = load(t, 0)
	w.Player = player
	w.Factions[other].Diplomat = noFaction
	w.Factions[player].Funds = 50_000
	w.Factions[other].Funds = 10_000
	w.events[0] = QueuedEvent{Code: uint16(other)<<8 | 6}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if w.Factions[player].Funds != 50_000 || w.Factions[other].Funds != 10_000 {
		t.Fatalf("事件 6 缺少外交官卻改寫資金：player=%d other=%d",
			w.Factions[player].Funds, w.Factions[other].Funds)
	}
}

// sub_17C6E 的數值核心不是每次加固定步長：數字鍵追加一位，另有
// 00、退位、還原初值與清零。事件 2／3、4／5 必須共用同一組上限語意。
func TestRawAmountEditorSemantics(t *testing.T) {
	w := load(t, 0)
	w.diplomacy = &DiplomacyChoice{InitialAmount: 1200, OfferAmount: 12}
	if !w.EditDiplomacyOfferAmount(AmountAppendDigit, 3) || w.diplomacy.OfferAmount != 123 {
		t.Fatalf("外交數字追加 = %d，want 123", w.diplomacy.OfferAmount)
	}
	if !w.EditDiplomacyOfferAmount(AmountAppendHundred, 0) || w.diplomacy.OfferAmount != 12_300 {
		t.Fatalf("外交追加 00 = %d，want 12300", w.diplomacy.OfferAmount)
	}
	if !w.EditDiplomacyOfferAmount(AmountDeleteDigit, 0) || w.diplomacy.OfferAmount != 1_230 {
		t.Fatalf("外交退位 = %d，want 1230", w.diplomacy.OfferAmount)
	}
	if !w.EditDiplomacyOfferAmount(AmountRestoreInitial, 0) || w.diplomacy.OfferAmount != 1_200 {
		t.Fatalf("外交還原 = %d，want 1200", w.diplomacy.OfferAmount)
	}
	if !w.EditDiplomacyOfferAmount(AmountClear, 0) || w.diplomacy.OfferAmount != 0 {
		t.Fatalf("外交清零 = %d，want 0", w.diplomacy.OfferAmount)
	}
	if w.EditDiplomacyOfferAmount(AmountAppendDigit, 10) || w.diplomacy.OfferAmount != 0 {
		t.Fatal("外交非法數字不應改寫狀態")
	}

	w.funding = &FundingChoice{RequestedAmount: 6500, OfferAmount: 29999}
	if !w.EditFundingAmount(AmountAppendDigit, 9) || w.funding.OfferAmount != 30_000 {
		t.Fatalf("撥款上限 = %d，want 30000", w.funding.OfferAmount)
	}
	if !w.EditFundingAmount(AmountRestoreInitial, 0) || w.funding.OfferAmount != 6_500 {
		t.Fatalf("撥款還原 = %d，want 6500", w.funding.OfferAmount)
	}
}

// 事件 4／5 在撥款三選一之前，分別先顯示 sub_132A9／sub_132E9 的
// TALK #56／#57；這裡只固定前置報告，不把 sub_139E8 後續訊息池猜成完整流程。
func TestQueuedFundingInitialTalkNotices(t *testing.T) {
	w := load(t, 0)
	player, city, officer := 0, 0, 1
	w.Player = player
	w.Cities[city].Owner = player
	w.Cities[city].Governor = officer
	w.Generals[officer] = General{Alive: true, Faction: player, Politics: 8, Budget: 0}
	w.events[0] = QueuedEvent{Code: uint16(city)<<8 | 4, Param: 1_000}
	w.eventCursor, w.eventDelay = 0, 1
	governor := &Event{}
	w.dispatchQueuedEvent(governor)
	if len(governor.TalkNotices) != 1 || governor.TalkNotices[0] != (TalkNotice{
		Index: 0x38, City: city, Faction: -1, General: officer, Amount: -1,
	}) {
		t.Fatalf("事件 4 前置 TALK 錯誤：%#v", governor.TalkNotices)
	}

	w = load(t, 0)
	w.Player = player
	other := 1
	w.Factions[other].Diplomat = officer
	w.Generals[officer] = General{Alive: true, Faction: other, Politics: 8, Budget: 0}
	w.events[0] = QueuedEvent{Code: uint16(other)<<8 | 5, Param: 1_000}
	w.eventCursor, w.eventDelay = 0, 1
	diplomat := &Event{}
	w.dispatchQueuedEvent(diplomat)
	if len(diplomat.TalkNotices) != 1 || diplomat.TalkNotices[0] != (TalkNotice{
		Index: 0x39, City: -1, Faction: other, General: officer, Amount: -1,
	}) {
		t.Fatalf("事件 5 前置 TALK 錯誤：%#v", diplomat.TalkNotices)
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

	// 這些區間是**還沒解的**，跑再久也不該被動到。事件佇列已經是
	// 已解出的可變欄位，另由 TestEventQueueStorage／Timing／MonthlyCompaction
	// 驗證，不能再把它當成未知區域。
	regions := []struct {
		name     string
		from, to int
	}{
		{"+0x3B–0x7F 不載入的空隙", 0x3B, 0x80},
		{"軍團表", 0x22C0, 0x42C0},
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

// 原版 idle path 在 sub_11D8E 之前先跑 sub_125A3：時鐘仍在「時 = 1」時，
// 該時的軍費／士氣更新就應該發生；不能先 Advance 成「時 = 2」才傳給軍團。
func TestTickRunsCorpsBeforeClockAdvance(t *testing.T) {
	w := load(t, 0)
	faction := w.AliveFactions()[0]
	w.Factions[faction].Expense = 0
	w.Factions[faction].Reserves = [economy.NumTroopTypes]int{}
	w.Corps[0] = Corps{
		Alive: true, Faction: faction, Men: 32, Morale: 100,
		Timer: 99, Interval: 99, Node: 0, TargetNode: 0,
	}
	w.corpsCursor = 0
	w.hourFaction = (faction + 1) % numFactions
	w.Clock.Hour, w.Clock.Subtick = 1, clock.SubticksPerHour-1

	got := w.Tick(rng.NewFixed(5))
	if !got.Clock.Hour || w.Clock.Hour != 2 || w.Clock.Subtick != 0 {
		t.Fatalf("clock = %+v／event=%+v，want 從 1:8 進位至 2:0", w.Clock, got.Clock)
	}
	wantExpense := combat.Upkeep(32, false)
	if got := w.Factions[faction].Expense; got != wantExpense {
		t.Fatalf("軍團在 clock advance 後才更新：expense=%d，want %d", got, wantExpense)
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
	// 每個兵種兩個槽，各分到上限 100 點，所以池各少 200
	// （docs/spec/21 §2：扣掉的量 ＝ 放進槽裡的量）。
	for tp, want := range map[economy.TroopType]int{
		economy.Cavalry: 5800, economy.Archer: 5800, economy.Infantry: 5800,
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

// 兵不夠不是錯誤——池裡有多少就分多少（docs/spec/21 §2）。
// 唯一的錯誤條件是分配完主將槽還是 0（§4，原版只檢查 [si+29h]）。
func TestFormCorpsDistributesWhateverIsLeft(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	lord := w.Factions[f].Lord
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{6000, 0, 6000}

	kinds := [army.Positions]army.TroopType{
		army.Cavalry, army.Cavalry, army.Archer, army.Archer, army.Cavalry, army.Cavalry,
	}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	if err := w.FormCorps(lord, kinds, manned); err != nil {
		t.Fatalf("弓兵是 0，但主將槽是騎馬，應該編得出來：%v", err)
	}
	c := w.Corps[lord]
	// 四個騎馬槽各 100 點；兩個弓兵槽分不到兵，是空槽。
	for _, k := range []int{0, 1, 4, 5} {
		if c.Units[k].Men != 100 {
			t.Errorf("槽 %d 兵力 %d，應為 100", k, c.Units[k].Men)
		}
	}
	for _, k := range []int{2, 3} {
		if c.Units[k].Men != 0 {
			t.Errorf("槽 %d 應該是空的，卻有 %d 點", k, c.Units[k].Men)
		}
	}
	if got := w.Factions[f].Reserves[economy.Cavalry]; got != 6000-400 {
		t.Errorf("騎馬池剩 %d，應為 %d——扣掉的要等於放進槽裡的", got, 6000-400)
	}
}

// 主將槽分不到兵就不成立：原版按確定時只檢查 [si+29h]。
func TestFormCorpsNeedsMenInGeneralSlot(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	lord := w.Factions[f].Lord
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{0, 6000, 6000}

	kinds := [army.Positions]army.TroopType{
		army.Cavalry, army.Archer, army.Archer, army.Archer, army.Archer, army.Archer,
	}
	manned := [army.Positions]bool{true, true, true, true, true, true}
	if err := w.FormCorps(lord, kinds, manned); err == nil {
		t.Fatal("主將槽分不到兵卻編成成功了")
	}
	if w.Corps[lord].Alive {
		t.Error("失敗卻留下了軍團")
	}
}

// 分配式本身：餘數整個給第一個同型槽，之後的槽對剩下的重分，每槽上限 100。
func TestDistributeReservesFollowsOriginal(t *testing.T) {
	// 池 250 點、三個同型槽：250/3=83 餘 1 → 第一槽 84，
	// 剩 166 分兩槽 → 83，剩 83 分一槽 → 83。合計 250，池歸零。
	pool := [economy.NumTroopTypes]int{250, 0, 0}
	kinds := [army.Positions]army.TroopType{
		army.Cavalry, army.Cavalry, army.Cavalry, army.Cavalry, army.Cavalry, army.Cavalry,
	}
	manned := [army.Positions]bool{true, true, true, false, false, false}
	got := distributeReserves(&pool, kinds, manned)
	want := [army.Positions]int{84, 83, 83, 0, 0, 0}
	if got != want {
		t.Fatalf("分配 = %v，want %v", got, want)
	}
	if pool[economy.Cavalry] != 0 {
		t.Errorf("池剩 %d，應該全部分完", pool[economy.Cavalry])
	}
	// 扣掉的總量必須等於放進槽裡的總量。
	sum := 0
	for _, n := range got {
		sum += n
	}
	if sum != 250 {
		t.Errorf("放進槽裡 %d 點，池少了 250——兩者必須相等", sum)
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

// 玩家的勢力捲進去先問戰鬥指揮／委任，其餘自動判定——原版的分派規則。
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
	for i := 0; i < 200000 && w.PendingBattle() == nil && w.PendingEncounter() == nil; i++ {
		w.Tick(r)
	}
	choice := w.PendingEncounter()
	if choice == nil {
		t.Fatal("玩家的軍團一路走到底都沒有出現戰鬥選擇")
	}
	if choice.Attacker != la || choice.Defender != lb {
		t.Fatalf("遭遇選擇是 %d vs %d，應為 %d vs %d",
			choice.Attacker, choice.Defender, la, lb)
	}
	if err := w.ChooseBattleCommand(); err != nil {
		t.Fatal(err)
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

// 委任要消費同一個遭遇選擇，直接回傳自動判定結果；選單掛著時時鐘不能走。
func TestPlayerBattleCanBeDelegated(t *testing.T) {
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

	r := rng.NewFixed(5)
	queued := &CorpsEvent{}
	w.fight(la, lb, queued, combat.Field, 0, r)
	if w.PendingEncounter() == nil {
		t.Fatal("玩家遭遇沒有進入選擇狀態")
	}
	before := w.Clock
	w.Tick(r)
	if w.Clock != before {
		t.Fatal("遭遇選單掛著時世界仍然前進")
	}
	ev := w.ChooseBattleDelegate(r)
	if ev == nil || ev.Battle == nil {
		t.Fatal("委任沒有回傳自動判定結果")
	}
	if w.PendingEncounter() != nil || w.PendingBattle() != nil {
		t.Fatal("委任完成後仍掛著戰鬥狀態")
	}
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

// 真實劇本的正常玩家切片：編成一槽步兵、沿 MMAP 道路走到汝南，
// 進入沒有敵方軍團的據點時，照原版規則用城兵自動判定攻城。
//
// 這裡刻意不掛 TacticalSetup，也不呼叫任何 demo／StageBattle 捷徑：
// 「城裡只有城兵」本來就不會進戰鬥指揮選單（docs/re/09 §2），
// 但仍必須由正常行軍觸發 `Enemy == -1, Mode == Siege`。
func TestNormalScenarioMarchIntoGarrison(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	xy := make([][2]int, len(w.Cities))
	for i := range w.Cities {
		xy[i] = [2]int{w.Cities[i].X, w.Cities[i].Y}
	}
	edges, err := world.RoadEdges(lib.World, xy)
	if err != nil {
		t.Fatal(err)
	}
	me := make([]march.Edge, len(edges))
	for i, e := range edges {
		me[i] = march.Edge{A: e.A, B: e.B, Steps: e.Steps,
			Path: e.Path, ACell: xy[e.A]}
	}
	w.SetRoads(march.New(len(w.Cities), me))

	leader := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction == w.Player && !g.Posted {
			leader = i
			break
		}
	}
	if leader < 0 {
		t.Fatal("找不到可編成的曹操武將")
	}
	var kinds [army.Positions]army.TroopType
	var manned [army.Positions]bool
	kinds[0] = army.Infantry
	manned[0] = true
	if err := w.FormCorps(leader, kinds, manned); err != nil {
		t.Fatal(err)
	}

	target := -1
	for i, c := range w.Cities {
		if text.Decode([]byte(c.Name), text.Big5) == "汝南" {
			target = i
			if c.Owner == w.Player || c.Garrison <= 0 {
				t.Fatalf("汝南不是可驗證的敵方城兵據點：owner=%d garrison=%d",
					c.Owner, c.Garrison)
			}
			break
		}
	}
	if target < 0 {
		t.Fatal("找不到汝南")
	}
	if err := w.March(leader, target); err != nil {
		t.Fatal(err)
	}

	var fought *CorpsEvent
	for i := 0; i < 100000 && fought == nil; i++ {
		for _, ev := range w.Tick(rng.NewFixed(1)).Corps {
			if ev.Battle != nil {
				copy := ev
				fought = &copy
				break
			}
		}
	}
	if fought == nil {
		t.Fatalf("正常行軍到汝南後沒有攻城：node=%d xy=%d,%d target=%d,%d",
			w.Corps[leader].Node, w.Corps[leader].X, w.Corps[leader].Y,
			w.Corps[leader].TargetX, w.Corps[leader].TargetY)
	}
	if fought.Enemy != -1 || fought.Mode != combat.Siege {
		t.Fatalf("正常城兵遭遇 = enemy=%d mode=%v，應為 enemy=-1、siege",
			fought.Enemy, fought.Mode)
	}
}

// 真實劇本的敵方 AI 遭遇：從正常編成／道路行軍進入戰鬥指揮後，
// 真實 BATTLE.MAP／BATTLE.MDL／BATTLE.DAT 也必須能完成並回寫戰略軍團。
// 這不是 `StageBattle` 或 `demoBattle` 捷徑；只有素材缺席時才跳過。
func TestNormalScenarioTacticalBattleTerminates(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	w.EnableStrategicAI()
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	xy := make([][2]int, len(w.Cities))
	for i := range w.Cities {
		xy[i] = [2]int{w.Cities[i].X, w.Cities[i].Y}
	}
	edges, err := world.RoadEdges(lib.World, xy)
	if err != nil {
		t.Fatal(err)
	}
	me := make([]march.Edge, len(edges))
	for i, e := range edges {
		me[i] = march.Edge{A: e.A, B: e.B, Steps: e.Steps,
			Path: e.Path, ACell: xy[e.A]}
	}
	w.SetRoads(march.New(len(w.Cities), me))

	leader := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction == w.Player && !g.Posted {
			leader = i
			break
		}
	}
	if leader < 0 {
		t.Fatal("找不到可編成的曹操武將")
	}
	var kinds [army.Positions]army.TroopType
	var manned [army.Positions]bool
	kinds[0], manned[0] = army.Infantry, true
	if err := w.FormCorps(leader, kinds, manned); err != nil {
		t.Fatal(err)
	}
	const target = 56 // 濮陽；正常固定種子 17 的敵方 AI 遭遇據點
	if w.Cities[target].Owner != w.Player {
		t.Fatalf("濮陽的正常玩家據點 owner=%d，應為玩家 %d", w.Cities[target].Owner, w.Player)
	}
	if err := w.March(leader, target); err != nil {
		t.Fatal(err)
	}

	read := func(name string) []byte {
		b, err := os.ReadFile("../../workplace/orig/dosv/" + name)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	battleLib, err := battle.Parse(read("BATTLE.MAP"), read("BATTLE.MDL"), read("BATTLE.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	forms, err := tactical.LoadFormations("../../workplace/orig/dosv/KI.EXE")
	if err != nil {
		t.Fatal(err)
	}
	w.SetTactical(&TacticalSetup{
		Forms: forms,
		Field: func(node int, siege bool) *tactical.Field {
			if !siege || node < 0 || node >= battle.NumFields {
				return nil
			}
			return tactical.NewFieldFromTiles(
				battleLib.Tiles(node), battleLib.Heights(node), battleLib.GateX(node))
		},
		Script: func(node int, siege bool, tactic int) []byte {
			if !siege {
				return nil
			}
			return battleLib.Script(tactic, battle.Category(node))
		},
	})

	r := rng.NewFixed(17)
	for i := 0; i < 200000 && w.PendingEncounter() == nil && w.PendingBattle() == nil; i++ {
		w.Tick(r)
	}
	choice := w.PendingEncounter()
	if choice == nil {
		t.Fatal("正常劇本沒有走到敵方 AI 遭遇選單")
	}
	if choice.Mode != combat.Siege || choice.Defender < 0 {
		t.Fatalf("正常敵方遭遇 = attacker=%d defender=%d mode=%v，應為軍團攻城",
			choice.Attacker, choice.Defender, choice.Mode)
	}
	if err := w.ChooseBattleCommand(); err != nil {
		t.Fatal(err)
	}
	p := w.PendingBattle()
	if p == nil {
		t.Fatal("正常敵方遭遇沒有建立戰術戰鬥")
	}
	beforeMen := [2]int{w.Corps[choice.Attacker].Men, w.Corps[choice.Defender].Men}
	if !p.Battle.Run(200000) {
		for side := range p.Battle.Sides {
			alive := 0
			cmds := map[tactical.Command]int{}
			minX, maxX := tactical.Width, 0
			for _, s := range p.Battle.Sides[side].Soldiers {
				if !s.Alive {
					continue
				}
				alive++
				cmds[s.Cmd]++
				if s.X < minX {
					minX = s.X
				}
				if s.X > maxX {
					maxX = s.X
				}
			}
			t.Logf("卡住 side=%d alive=%d reserve=%v cmds=%v x=%d..%d",
				side, alive, p.Battle.Sides[side].Reserve, cmds, minX, maxX)
			for k := 0; k < 4; k++ {
				s := p.Battle.Sides[side].Soldiers[k]
				t.Logf("  soldier=%d alive=%v kind=%v hp=%d xyz=%d,%d,%d goal=%d,%d,%d step=%d,%d,%d target=%d path=%d",
					k, s.Alive, s.Kind, s.HP, s.X, s.Y, s.Z,
					s.GoalX, s.GoalY, s.GoalZ, s.StepX, s.StepY, s.StepZ,
					s.Target, s.Path.Len())
			}
			if side == 0 {
				for k, s := range p.Battle.Sides[side].Soldiers {
					if !s.Alive {
						continue
					}
					pnt, ok := s.Path.Current()
					t.Logf("  retreating soldier=%d xyz=%d,%d,%d edgeZ=%d climb=%v face=%d pathCurrent=%v/%v",
						k, s.X, s.Y, s.Z, p.Battle.Field.StandLevel(tactical.MinCoord, s.Y),
						s.CanClimb(), s.Facing, pnt, ok)
				}
			}
		}
		for i, st := range p.Battle.Structures {
			t.Logf("結構 %d kind=%d x=%d y=%d run=%d hp=%d broken=%v",
				i, st.Kind, st.X, st.Y, st.Run, st.Durability, st.Broken)
		}
		t.Fatalf("正常真實攻城跑了 20 萬幀還沒結束：攻方 %d、守方 %d",
			p.Battle.Sides[0].Remaining(), p.Battle.Sides[1].Remaining())
	}
	ev := w.ResolvePending(r)
	if ev == nil || ev.Battle == nil || w.PendingBattle() != nil {
		t.Fatal("正常真實攻城無法回寫戰略結果")
	}
	if ev.BattleBefore != beforeMen {
		t.Fatalf("戰後事件沒有保留戰前兵力：got=%v want=%v", ev.BattleBefore, beforeMen)
	}
	afterMen := [2]int{w.Corps[choice.Attacker].Men, w.Corps[choice.Defender].Men}
	if ev.BattleAfter != afterMen {
		t.Fatalf("戰後事件兵力與戰略狀態脫鉤：got=%v state=%v", ev.BattleAfter, afterMen)
	}
	if ev.BattleCityDamage != p.Battle.CityDamage(p.CityWall) {
		t.Fatalf("戰後事件城損與戰術結果脫鉤：got=%d", ev.BattleCityDamage)
	}
	t.Logf("正常真實攻城第 %d 幀結束，守方勝 %v；攻方 %d 點、守方 %d 點",
		p.Battle.Frame, ev.Battle.DefenderWins,
		w.Corps[choice.Attacker].Men, w.Corps[choice.Defender].Men)
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

// 事件 9（sub_13485 → sub_150D7）的 Code 高 byte 是 General 索引，
// 不是事件來源勢力；釋放時清出陣／俘虜欄位，並依原俘虜方存亡回寫勢力。
func TestQueuedEventReleaseGeneral(t *testing.T) {
	w := load(t, 0)
	w.Player = -1
	generalID, captor, faction := 7, 1, 2
	w.Generals[generalID] = General{
		Alive: true, Faction: faction, Captor: captor, Posted: true,
	}
	w.Factions[captor].Alive = true
	w.events[0] = QueuedEvent{
		Code: uint16(generalID)<<8 | 9,
	}
	w.eventCursor, w.eventDelay = 0, 1
	ev := &Event{}
	w.dispatchQueuedEvent(ev)
	if got := w.Generals[generalID]; got.Faction != captor || got.Captor != noFaction || got.Posted {
		t.Fatalf("事件 9 未依存活俘虜方釋放武將：%+v", got)
	}
	if len(ev.ReleasedGenerals) != 1 || ev.ReleasedGenerals[0] != generalID {
		t.Fatalf("事件 9 沒有回報釋放武將：%v", ev.ReleasedGenerals)
	}

	w = load(t, 0)
	w.Player = -1
	deadCaptor := 1
	w.Generals[generalID] = General{
		Alive: true, Faction: faction, Captor: deadCaptor, Posted: true,
	}
	w.Factions[deadCaptor].Alive = false
	w.events[0] = QueuedEvent{
		Code: uint16(generalID)<<8 | 9,
	}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if got := w.Generals[generalID]; got.Faction != noFaction || got.Captor != noFaction || got.Posted {
		t.Fatalf("事件 9 未將已滅勢力俘虜放回在野：%+v", got)
	}
}

// 玩家是事件 2 的合作方、事件 3 的停戰對象時，dispatch 只掛起三選一；
// World.Tick 與同一小時後續財政處理都不能穿透模態狀態。
func TestQueuedDiplomacyChoice(t *testing.T) {
	w := load(t, 0)
	player, invader, invaded := 0, 1, 2
	w.Player = player
	w.Factions[player].Lord = 0
	w.Factions[player].Diplomat = noFaction
	w.Factions[player].Aggression = 1
	w.Factions[invaded].Diplomat = noFaction
	for i := range w.Generals {
		if w.Generals[i].Faction == invaded {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[0] = General{Alive: true, Politics: 12, Posted: true, Faction: player, Captor: noFaction}
	w.Generals[1] = General{Alive: true, Politics: 8, Faction: invaded, Captor: noFaction}
	w.Corps[0] = Corps{Alive: true, Faction: player}
	w.Factions[player].Funds = 10_000
	w.Factions[invaded].Funds = 50_000
	w.Friendship[player][invader] = diplomacy.Peace(40)
	w.Friendship[player][invaded] = diplomacy.Peace(60)
	w.Factions[invader].InvasionTarget = diplomacy.NoTarget
	w.events[0] = QueuedEvent{
		Code:  uint16(player)<<8 | 2,
		Param: uint16(invaded)<<8 | uint16(invader),
	}
	w.eventCursor, w.eventDelay = 0, 1
	beforeClock := w.Clock
	w.dispatchQueuedEvent(&Event{})
	c := w.PendingDiplomacy()
	if c == nil || c.Kind != DiplomacyCooperation || c.Source != player || c.Invader != invader || c.Target != invaded {
		t.Fatalf("事件 2 沒有掛起正確外交選擇：%+v", c)
	}
	if w.Factions[player].Funds != 10_000 || w.Factions[invaded].Funds != 50_000 {
		t.Fatalf("事件 2 在選擇前就改了資金：player=%d invaded=%d", w.Factions[player].Funds, w.Factions[invaded].Funds)
	}
	w.Tick(rng.NewFixed(1))
	if w.Clock != beforeClock {
		t.Fatalf("外交選單掛起時時鐘仍前進：got=%+v want=%+v", w.Clock, beforeClock)
	}
	if !w.SetDiplomacyOfferAmount(7_000) || !w.ResolveDiplomacy(DiplomacyOfferFunds) {
		t.Fatal("事件 2 的提供資金選項沒有完成狀態收尾")
	}
	if w.PendingDiplomacy() != nil || w.Factions[player].Funds != 17_000 || w.Factions[invaded].Funds != 43_000 {
		t.Fatalf("事件 2 選擇後狀態錯誤：pending=%+v player=%d invaded=%d", w.PendingDiplomacy(), w.Factions[player].Funds, w.Factions[invaded].Funds)
	}

	w = load(t, 0)
	player, source := 0, 1
	w.Player = player
	w.Factions[player].Lord = player
	w.Factions[player].Diplomat = noFaction
	w.Factions[source].Diplomat = noFaction
	for i := range w.Generals {
		if w.Generals[i].Faction == source {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[player] = General{Alive: true, Politics: 12, Posted: true, Faction: player, Captor: noFaction}
	w.Generals[source+1] = General{Alive: true, Politics: 8, Faction: source, Captor: noFaction}
	w.Corps[player] = Corps{Alive: true, Faction: player}
	w.Friendship[player][source] = diplomacy.Peace(20)
	w.Friendship[source][player] = diplomacy.Peace(50)
	w.Factions[source].InvasionTarget = player
	w.Factions[player].InvasionTarget = 9
	w.events[0] = QueuedEvent{
		Code:  uint16(source)<<8 | 3,
		Param: uint16(0xFF00 | player),
	}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	c = w.PendingDiplomacy()
	if c == nil || c.Kind != DiplomacyCeasefire || c.Source != source || c.Target != player {
		t.Fatalf("事件 3 沒有掛起正確外交選擇：%+v", c)
	}
	if !w.ResolveDiplomacy(DiplomacyAcceptFree) || w.PendingDiplomacy() != nil {
		t.Fatal("事件 3 的無條件同意沒有完成狀態收尾")
	}
	if w.Factions[source].InvasionTarget != diplomacy.NoTarget || w.Friendship[source][player] != diplomacy.Peace(20) {
		t.Fatalf("事件 3 選擇後停戰狀態錯誤：target=%d friendship=%#v", w.Factions[source].InvasionTarget, w.Friendship[source][player])
	}
}

// 事件 2／3 進入玩家三選一前，原版先顯示請求報告：停戰是 #360，協力是
// #373；兩者的 {3} 都是提出請求的勢力君主，而不是事件字高的數字。
func TestQueuedDiplomacyChoiceTalkNotices(t *testing.T) {
	w := load(t, 0)
	player, invader, invaded := 0, 1, 2
	w.Player = player
	w.Factions[player].Lord = player
	w.Factions[player].Diplomat = noFaction
	w.Factions[player].Aggression = 1
	w.Factions[invader].Diplomat = noFaction
	w.Factions[invaded].Diplomat = noFaction
	for i := range w.Generals {
		if w.Generals[i].Faction == invader || w.Generals[i].Faction == invaded {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[0] = General{Alive: true, Politics: 12, Posted: true, Faction: player, Captor: noFaction}
	w.Generals[1] = General{Alive: true, Politics: 8, Faction: invaded, Captor: noFaction}
	w.Corps[0] = Corps{Alive: true, Faction: player}
	w.Factions[player].Funds = 10_000
	w.Factions[invaded].Funds = 50_000
	w.Friendship[player][invader] = diplomacy.Peace(40)
	w.Friendship[player][invaded] = diplomacy.Peace(60)
	w.Factions[invader].InvasionTarget = diplomacy.NoTarget
	w.events[0] = QueuedEvent{
		Code:  uint16(player)<<8 | 2,
		Param: uint16(invaded)<<8 | uint16(invader),
	}
	w.eventCursor, w.eventDelay = 0, 1
	ev := Event{}
	w.dispatchQueuedEvent(&ev)
	want := []TalkNotice{{Index: 0x175, City: -1, Faction: invader, General: -1, Amount: -1}}
	if !reflect.DeepEqual(ev.TalkNotices, want) {
		t.Fatalf("事件 2 前置 TALK 錯誤：got %#v want %#v", ev.TalkNotices, want)
	}

	w = load(t, 0)
	player, source := 0, 1
	w.Player = player
	w.Factions[player].Lord = player
	w.Factions[player].Diplomat = noFaction
	w.Factions[source].Diplomat = noFaction
	for i := range w.Generals {
		if w.Generals[i].Faction == source {
			w.Generals[i].Posted = true
		}
	}
	w.Generals[player] = General{Alive: true, Politics: 12, Posted: true, Faction: player, Captor: noFaction}
	w.Generals[source+1] = General{Alive: true, Politics: 8, Faction: source, Captor: noFaction}
	w.Corps[player] = Corps{Alive: true, Faction: player}
	w.Friendship[player][source] = diplomacy.Peace(20)
	w.Friendship[source][player] = diplomacy.Peace(50)
	w.Factions[source].InvasionTarget = player
	w.events[0] = QueuedEvent{
		Code:  uint16(source)<<8 | 3,
		Param: uint16(0xFF00 | player),
	}
	w.eventCursor, w.eventDelay = 0, 1
	ev = Event{}
	w.dispatchQueuedEvent(&ev)
	want = []TalkNotice{{Index: 0x168, City: -1, Faction: source, General: -1, Amount: -1}}
	if !reflect.DeepEqual(ev.TalkNotices, want) {
		t.Fatalf("事件 3 前置 TALK 錯誤：got %#v want %#v", ev.TalkNotices, want)
	}
}

// 事件 4／5 的玩家撥款在選擇前不改世界；指定金額會照 sub_139E8 寫入
// 官員經費（amount／128）並從玩家資金扣款，拒絕則完全無副作用。
func TestQueuedFundingChoice(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	w.Factions[w.Player].Funds = 50_000
	cityID := -1
	for i, c := range w.Cities {
		if c.Owner == w.Player {
			cityID = i
			break
		}
	}
	if cityID < 0 {
		t.Fatal("找不到玩家據點")
	}
	officer := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction == w.Player {
			officer = i
			break
		}
	}
	if officer < 0 {
		t.Fatal("找不到玩家武將")
	}
	w.Cities[cityID].Governor = officer
	w.Generals[officer].Budget = 0
	w.events[0] = QueuedEvent{Code: uint16(cityID)<<8 | 4, Param: 6_000}
	w.eventCursor, w.eventDelay = 0, 1
	beforeClock := w.Clock
	w.dispatchQueuedEvent(&Event{})
	c := w.PendingFunding()
	if c == nil || c.Kind != FundingGovernor || c.Subject != cityID || c.Officer != officer ||
		c.RequestedAmount != 6_000 || c.OfferAmount != 6_000 {
		t.Fatalf("事件 4 沒有掛起正確撥款選擇：%+v", c)
	}
	if w.Factions[w.Player].Funds != 50_000 || w.Generals[officer].Budget != 0 {
		t.Fatalf("事件 4 在選擇前就改了狀態：funds=%d budget=%d",
			w.Factions[w.Player].Funds, w.Generals[officer].Budget)
	}
	w.Tick(rng.NewFixed(1))
	if w.Clock != beforeClock {
		t.Fatalf("撥款選單掛起時時鐘仍前進：got=%+v want=%+v", w.Clock, beforeClock)
	}
	if !w.SetFundingAmount(7_000) || !w.ResolveFunding(FundingSetAmount) {
		t.Fatal("事件 4 的指定金額選項沒有完成狀態收尾")
	}
	if w.PendingFunding() != nil || w.Factions[w.Player].Funds != 43_000 ||
		w.Generals[officer].Budget != 7_000/128 {
		t.Fatalf("事件 4 選擇後狀態錯誤：pending=%+v funds=%d budget=%d",
			w.PendingFunding(), w.Factions[w.Player].Funds, w.Generals[officer].Budget)
	}

	w = load(t, 0)
	w.Player = 0
	w.Factions[w.Player].Funds = 50_000
	target := 1
	diplomat := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction != w.Player {
			diplomat = i
			break
		}
	}
	if diplomat < 0 {
		t.Fatal("找不到非玩家外交官")
	}
	w.Factions[target].Diplomat = diplomat
	w.Generals[diplomat].Budget = 0
	w.events[0] = QueuedEvent{Code: uint16(target)<<8 | 5, Param: 400}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	c = w.PendingFunding()
	if c == nil || c.Kind != FundingDiplomat || c.Subject != target || c.Officer != diplomat ||
		c.RequestedAmount != fundingMinOffer {
		t.Fatalf("事件 5 沒有套用非零初始金額 500 下限：%+v", c)
	}
	if w.ResolveFunding(FundingReject) || w.PendingFunding() != nil {
		t.Fatal("事件 5 拒絕不應回報成功或留下 pending")
	}
	if w.Factions[w.Player].Funds != 50_000 || w.Generals[diplomat].Budget != 0 {
		t.Fatalf("事件 5 拒絕仍有副作用：funds=%d budget=%d",
			w.Factions[w.Player].Funds, w.Generals[diplomat].Budget)
	}
}

// sub_13902 與 sub_139E8 對「指定金額」的上限語意不同：外交超過
// 初始要求會回傳拒絕碼；撥款超過初始要求則仍會寫入，只有輸入 0 無效。
func TestDiplomacyAndFundingAmountOutcomeBounds(t *testing.T) {
	w := load(t, 0)
	player, invader, target := 0, 1, 2
	w.Player = player
	w.Factions[player].Lord = 0
	w.Factions[target].Lord = 2
	w.Factions[player].Diplomat = noFaction
	w.Factions[target].Diplomat = noFaction
	w.Generals[0] = General{Alive: true, Faction: player, Politics: 12, Posted: true}
	w.Generals[1] = General{Alive: true, Faction: invader, Politics: 8}
	w.Generals[2] = General{Alive: true, Faction: target, Politics: 8, Posted: true}
	w.Corps[0] = Corps{Alive: true, Faction: player}
	w.Factions[player].Funds = 10_000
	w.Factions[target].Funds = 50_000
	w.Friendship[player][invader] = diplomacy.Peace(40)
	w.Friendship[player][target] = diplomacy.Peace(60)
	w.Factions[invader].InvasionTarget = diplomacy.NoTarget
	w.events[0] = QueuedEvent{Code: uint16(player)<<8 | 2,
		Param: uint16(target)<<8 | uint16(invader)}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	c := w.PendingDiplomacy()
	if c == nil || c.InitialAmount != 15_000 {
		t.Fatalf("外交初始金額錯誤：%+v", c)
	}
	trustBeforeOverOffer := w.Trust
	if !w.SetDiplomacyOfferAmount(15_001) || w.ResolveDiplomacy(DiplomacyOfferFunds) {
		t.Fatal("外交超過初始要求不應完成")
	}
	if w.Factions[player].Funds != 10_000 || w.Factions[target].Funds != 50_000 ||
		w.Friendship[player][invader].AtWar() ||
		w.Trust != clampU8(trustBeforeOverOffer-0x1E) {
		t.Fatalf("外交超額輸入仍有副作用：player=%d target=%d", w.Factions[player].Funds, w.Factions[target].Funds)
	}

	w = load(t, 0)
	w.Player = 0
	w.Factions[0].Funds = 50_000
	officer := 1
	w.Cities[0].Owner = 0
	w.Cities[0].Governor = officer
	w.Generals[officer] = General{Alive: true, Faction: 0, Budget: 0}
	w.events[0] = QueuedEvent{Code: 4, Param: 6_000}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if !w.SetFundingAmount(0) || w.ResolveFunding(FundingSetAmount) {
		t.Fatal("撥款指定 0 不應完成")
	}
	if w.Factions[0].Funds != 50_000 || w.Generals[officer].Budget != 0 {
		t.Fatalf("撥款指定 0 仍有副作用：funds=%d budget=%d",
			w.Factions[0].Funds, w.Generals[officer].Budget)
	}

	w = load(t, 0)
	w.Player = 0
	w.Factions[0].Funds = 50_000
	w.Cities[0].Owner = 0
	w.Cities[0].Governor = officer
	w.Generals[officer] = General{Alive: true, Faction: 0, Budget: 0}
	w.events[0] = QueuedEvent{Code: 4, Param: 6_000}
	w.eventCursor, w.eventDelay = 0, 1
	w.dispatchQueuedEvent(&Event{})
	if !w.SetFundingAmount(7_000) || !w.ResolveFunding(FundingSetAmount) {
		t.Fatal("撥款超過初始要求應完成")
	}
	if w.Factions[0].Funds != 43_000 || w.Generals[officer].Budget != 7_000/128 {
		t.Fatalf("撥款超額輸入結果錯誤：funds=%d budget=%d",
			w.Factions[0].Funds, w.Generals[officer].Budget)
	}
}

// 下行軍指令時，原版把目的地同時寫進兩個欄位：`+0x14`（×8，移動用）與
// `+0x20`（無縮放，一覽表顯示與遷都的判準）。`sub_142AB` 直接示範兩者的關係
// ——先 `mov [si+14h],bx` 再 `shr bx,1` ×3 後 `mov [si+20h],bl`。
//
// 這個測試盯的是**存檔內容**而不只是欄位：先前 `Ordered` 只在編成與遷都
// 更新，所以下完行軍指令再存檔，`+0x20` 會停在首都，與原版分歧。
func TestMarchWritesOrderedTargetIntoSave(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	leader := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction == w.Player && !g.Posted {
			leader = i
			break
		}
	}
	if leader < 0 {
		t.Fatal("找不到可編成的武將")
	}
	var kinds [army.Positions]army.TroopType
	var manned [army.Positions]bool
	kinds[0] = army.Infantry
	manned[0] = true
	if err := w.FormCorps(leader, kinds, manned); err != nil {
		t.Fatal(err)
	}
	corps := leader

	// 編成之後目標就是所在地（＝首都），這一點本身也要成立：
	// 原版 `sub_16F26` 把勢力的首都寫進 +0x20，而軍團就編在首都。
	if got := w.Corps[corps].Ordered; got != w.Corps[corps].Node {
		t.Fatalf("編成後 Ordered = %d，應該等於所在地 %d",
			got, w.Corps[corps].Node)
	}

	start := w.Corps[corps].Node
	dest := -1
	for i := range w.Cities {
		if i != start {
			dest = i
			break
		}
	}
	if err := w.March(corps, dest); err != nil {
		// 走不到時原版與 remake 都把目標退回原地，這也要一致。
		if got := w.Corps[corps].Ordered; got != w.Corps[corps].Node {
			t.Fatalf("行軍失敗後 Ordered = %d，應該退回原地 %d",
				got, w.Corps[corps].Node)
		}
		return
	}
	if got := w.Corps[corps].Ordered; got != dest {
		t.Fatalf("Ordered = %d，應該是行軍目的地 %d", got, dest)
	}

	b := w.Bytes()
	r := b[corpsBase+corps*corpsSize:]
	if int(r[0x20]) != dest {
		t.Fatalf("存檔 +0x20 = %d，應該是目的地 %d", r[0x20], dest)
	}
	if int(u16(r, 0x14))/8 != dest {
		t.Fatalf("存檔 +0x14÷8 = %d，應該是目的地 %d", u16(r, 0x14)/8, dest)
	}
}

// 空槽在原版的兵種欄是 **4**，不是「兵種 0 而人數 0」。
//
// 兩個地方靠這個值：`sub_14717` 退兵時看到 4 就跳過，兩個情報視窗取
// 兵種圖示時 `(兵種−1)×0xC0` 也要算到第四張（docs/re/51 §4）。
// 先前 remake 把沒編兵的槽留成騎馬，寫回存檔就變成一支「騎馬空隊」。
func TestFormCorpsMarksEmptySlots(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	leader := w.Factions[f].Lord
	w.Factions[f].Reserves = [economy.NumTroopTypes]int{100, 0, 0}

	var kinds [army.Positions]army.TroopType
	manned := [army.Positions]bool{true}
	if err := w.FormCorps(leader, kinds, manned); err != nil {
		t.Fatalf("編成失敗：%v", err)
	}
	c := w.Corps[leader]
	if c.Units[0].Men != 100 || c.Units[0].Kind != army.Cavalry {
		t.Fatalf("主將槽 = %+v，want 100 點騎馬", c.Units[0])
	}
	for k := 1; k < army.Positions; k++ {
		if c.Units[k].Men != 0 || c.Units[k].Kind != EmptySlotKind {
			t.Errorf("第 %d 槽 = %+v，want 空槽（兵種 %d）", k, c.Units[k], EmptySlotKind)
		}
	}
	// 寫回存檔時是 1-based 的 4。
	if got := byteFromKind(EmptySlotKind); got != 4 {
		t.Errorf("空槽寫回 = %d，want 4", got)
	}
}
