package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// aliveCorpsSlot 找一個劇本裡活著的軍團槽；沒有就自己灌一個進原始 bytes。
func aliveCorpsSlot(t *testing.T, w *World) int {
	t.Helper()
	for i := 0; i < numCorps; i++ {
		if w.raw[corpsBase+i*corpsSize] >= aliveFlag {
			return i
		}
	}
	// 劇本開局沒有軍團——灌一個最小的活著的記錄。
	off := corpsBase
	w.raw[off+0x00] = newCorps
	w.raw[off+0x01] = 0
	return 0
}

// TestCountdownRoundTripsBothMeanings 釘住 docs/spec/175 §1.5：
// `+0x03` 是**一個 byte、兩種用途**，用 `+0x00` 分辨。
//
// ⭐ 兩種用途都要 round-trip。只驗其中一種的話，
// 「載入端把敗走倒數讀進對峙倒數」這種錯誤照樣通過——
// 兩個欄位的值在單一狀態下長得一樣。
func TestCountdownRoundTripsBothMeanings(t *testing.T) {
	t.Run("對峙", func(t *testing.T) {
		w := load(t, 0)
		i := aliveCorpsSlot(t, w)
		off := corpsBase + i*corpsSize
		w.raw[off+0x00] |= standoffBit
		w.raw[off+0x03] = 0x0C

		w2 := loadBlock(w.raw)
		c := &w2.Corps[i]
		if !c.Standoff {
			t.Fatal("位元 5 沒讀成 Standoff")
		}
		if c.Routing {
			t.Fatal("活著的軍團被讀成敗走中")
		}
		if c.Countdown != 12 {
			t.Fatalf("對峙倒數 = %d，want 12", c.Countdown)
		}
		out := w2.Bytes()
		if got := out[off+0x00] & standoffBit; got == 0 {
			t.Errorf("寫回的 +0x00 = %#02x，位元 5 掉了", out[off+0x00])
		}
		if got := out[off+0x03]; got != 0x0C {
			t.Errorf("寫回的 +0x03 = %#02x，want 0x0c", got)
		}
	})

	t.Run("敗走", func(t *testing.T) {
		w := load(t, 0)
		i := aliveCorpsSlot(t, w)
		off := corpsBase + i*corpsSize
		w.raw[off+0x00] = 0x08 // 原版 sub_12977 整個 byte 寫 8
		w.raw[off+0x03] = 0x30

		w2 := loadBlock(w.raw)
		c := &w2.Corps[i]
		if c.Alive || !c.Routing {
			t.Fatalf("旗標 08h 應該是「敗走中、不算活著」，得到 Alive=%v Routing=%v",
				c.Alive, c.Routing)
		}
		if c.Standoff {
			t.Fatal("敗走中的軍團不該同時對峙——位元 5 在旗標 08h 裡根本沒設")
		}
		if c.Countdown != 48 {
			t.Fatalf("敗走倒數 = %d，want 48", c.Countdown)
		}
		out := w2.Bytes()
		if got := out[off+0x03]; got != 0x30 {
			t.Errorf("寫回的 +0x03 = %#02x，want 0x30", got)
		}
	})
}

// TestStandoffAndRoutingAreExclusive 是上一支的負對照：
// **兩種用途互斥**，所以「同時成立」在載入端就不該出現。
//
// ⚠ 判準不是「remake 不會自己弄出這種狀態」——存檔可以是任何 bytes。
// 原版靠旗標分辨：敗走是 `08h`（< `80h`），對峙要 `≥ 80h`，
// 所以**同一個 byte 不可能同時滿足兩邊**。
func TestStandoffAndRoutingAreExclusive(t *testing.T) {
	w := load(t, 0)
	i := aliveCorpsSlot(t, w)
	off := corpsBase + i*corpsSize
	// 硬灌一個「旗標 8 又帶位元 5」的壞記錄。
	w.raw[off+0x00] = 0x08 | standoffBit
	w.raw[off+0x03] = 0x0C

	c := &loadBlock(w.raw).Corps[i]
	if c.Standoff && c.Routing {
		t.Fatal("敗走與對峙同時成立——兩種倒數會互相踩到")
	}
	if !c.Routing {
		t.Fatal("旗標含 8 而且不到 0x80，應該是敗走中")
	}
}

// TestTickStandoffCountsDown 釘住 `sub_1264A`（docs/spec/175 §1.3）。
//
// ⭐ 判準有三個，缺一個就驗不到：
//   - 對峙中每個巡迴週期減 1
//   - 減到 0 停在 **1**（不是 0）——那一格是「下一次撞上就開打」的訊號
//   - 沒在對峙就歸零
func TestTickStandoffCountsDown(t *testing.T) {
	w := load(t, 0)
	i := aliveCorpsSlot(t, w)
	c := &w.Corps[i]
	c.Alive, c.Standoff, c.Countdown = true, true, 12

	for k := 11; k >= 1; k-- {
		w.tickStandoff(i)
		if c.Countdown != k {
			t.Fatalf("第 %d 次倒數 = %d，want %d", 12-k, c.Countdown, k)
		}
	}
	// 再減一次：dec 到 0 之後停在 1。
	w.tickStandoff(i)
	if c.Countdown != 1 {
		t.Fatalf("減到 0 之後 = %d，want 1", c.Countdown)
	}

	// 擋著的東西不見了 → 歸零。
	c.Standoff = false
	w.tickStandoff(i)
	if c.Countdown != 0 {
		t.Fatalf("沒在對峙時 = %d，want 0", c.Countdown)
	}
}

// TestTickStandoffWrapsLikeDecByte 是照抄原版的那一條：
// `dec byte ptr [si+3]` 在值是 0 時會變 **255**，不是負的。
//
// 正常流程碰不到（設位元 5 的同一支常式當場寫 12），但**存檔可以**——
// 而「看起來比較合理」的 `if > 0 才減` 會讓那種存檔卡在 0 永遠不開打。
func TestTickStandoffWrapsLikeDecByte(t *testing.T) {
	w := load(t, 0)
	i := aliveCorpsSlot(t, w)
	c := &w.Corps[i]
	c.Alive, c.Standoff, c.Countdown = true, true, 0

	w.tickStandoff(i)
	if c.Countdown != 255 {
		t.Fatalf("從 0 減一次 = %d，want 255（byte 運算會繞回去）", c.Countdown)
	}
}

// TestAliveCorpsWritesCountdownByte 釘住寫回端：**活著的軍團也要寫 `+0x03`**。
//
// ⚠ 原版不在對峙時那一格恆為 0（`sub_1264A` 每個週期歸零），
// 所以不寫回等於把存檔裡的舊值留著——而那個舊值會被下一次載入
// 讀成一個已經走到一半的對峙倒數。
func TestAliveCorpsWritesCountdownByte(t *testing.T) {
	w := load(t, 0)
	i := aliveCorpsSlot(t, w)
	off := corpsBase + i*corpsSize
	w.raw[off+0x03] = 0x07 // 存檔裡的殘值

	w2 := loadBlock(w.raw)
	w2.Corps[i].Standoff, w2.Corps[i].Countdown = false, 0
	if got := w2.Bytes()[off+0x03]; got != 0 {
		t.Fatalf("沒在對峙時寫回的 +0x03 = %#02x，want 0", got)
	}
}

// ── 觸發層：`sub_12708` 踏進去之前先問（docs/spec/175 §1.1–§1.35）──

// standoffFixture 編出兩個敵對勢力的軍團，回傳世界與兩支的槽序。
func standoffFixture(t *testing.T) (*World, int, int) {
	t.Helper()
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
	att, def := w.Factions[a].Lord, w.Factions[b].Lord
	if err := w.FormCorps(att, kinds, manned); err != nil {
		t.Fatal(err)
	}
	if err := w.FormCorps(def, kinds, manned); err != nil {
		t.Fatal(err)
	}
	// 交戰中——邊界掉頭（docs/spec/132）會在和平時把軍團彈回去，
	// 那一條先發生，對峙就一次都測不到。
	w.Friendship[a][b] = w.Friendship[a][b].WithWar(true)
	w.Friendship[b][a] = w.Friendship[b][a].WithWar(true)
	return w, att, def
}

// emptyPair 找兩格相鄰的野外空地（沒有據點、也沒有別的軍團站著）。
func emptyPair(t *testing.T, w *World) (int, int) {
	t.Helper()
	occupied := func(x, y int) bool {
		if w.cityAt(x, y) >= 0 {
			return true
		}
		for i := range w.Corps {
			if w.Corps[i].Alive && w.Corps[i].X == x && w.Corps[i].Y == y {
				return true
			}
		}
		return false
	}
	for y := 1; y < 240; y++ {
		for x := 2; x < 320; x++ {
			if !occupied(x, y) && !occupied(x-1, y) {
				return x, y
			}
		}
	}
	t.Fatal("找不到兩格相鄰的空地")
	return 0, 0
}

// faceOff 把守方擺在 (x,y)、攻方擺在它左邊一格，攻方下一步就要踏上去。
// **每個巡迴週期都是移動拍**（`Interval = 1`），這樣「第幾個週期開打」
// 直接等於原版的倒數次數。
func faceOff(w *World, att, def, x, y int) {
	w.ClearMarchRoute(att)
	w.ClearMarchRoute(def)
	d := &w.Corps[def]
	d.X, d.Y = x, y
	d.TargetX, d.TargetY = x, y
	d.TargetNode = d.Node

	c := &w.Corps[att]
	c.X, c.Y = x-1, y
	c.TargetX, c.TargetY = x, y
	c.Interval, c.Timer = 1, 1
	c.Standoff, c.Countdown = false, 0
	// ⚠ 攻方要處在「**還沒到行軍目標**」的狀態，否則 `tickOneCorps` 走的是
	// 抵達分派那一條而不是移動——`atTargetNode` 的判準是
	// 「`LinkAddr == 0` 而且 `Node == TargetNode`」，而 `ClearMarchRoute`
	// 剛把 `LinkAddr` 清成 0（docs/spec/179 §2）。借守方的據點當目標。
	c.TargetNode = d.Node
	if c.TargetNode == c.Node {
		c.TargetNode = (c.Node + 1) % len(w.Cities)
	}
}

// TestStandoffStopsBeforeEnemyCell 釘住 §1.1–§1.2：下一格站著敵方軍團時
// **這一拍不動**，設 `+0x00` 位元 5、`+0x03` ← 12。
//
// ⭐ 「不動」是這一條的重點。只驗旗標的話，「照樣走進去但順便設個旗標」
// 也會通過——而那正是 remake 原本的行為（一走到就結算）。
func TestStandoffStopsBeforeEnemyCell(t *testing.T) {
	w, att, def := standoffFixture(t)
	x, y := emptyPair(t, w)
	faceOff(w, att, def, x, y)

	ev := w.tickOneCorps(att, 0, rng.New(0, 0, 0))
	c := &w.Corps[att]
	if c.X != x-1 || c.Y != y {
		t.Errorf("軍團走到 (%d,%d)，應該停在 (%d,%d)——撞上敵人的那一格不進去",
			c.X, c.Y, x-1, y)
	}
	if !c.Standoff {
		t.Error("位元 5 沒設")
	}
	if c.Countdown != standoffTicks {
		t.Errorf("+0x03 = %d，want %d", c.Countdown, standoffTicks)
	}
	if ev != nil && ev.Battle != nil {
		t.Error("第一次撞上就開打了——原版要先對峙 12 個週期")
	}
}

// TestStandoffFightsOnTwelfthCycle 釘住 §1.3：倒數由**每次巡到**減 1，
// 減到 1 的那一次才結算 ⇒ 第 12 個巡迴週期，一圈 8 拍就是 96 拍。
//
// ⭐ 這一支同時是「一撞上就打」與「用移動拍計數」兩種寫法的負對照：
// 前者會落在第 1 個週期，後者在 `Interval > 1` 時會晚好幾倍。
func TestStandoffFightsOnTwelfthCycle(t *testing.T) {
	w, att, def := standoffFixture(t)
	x, y := emptyPair(t, w)
	faceOff(w, att, def, x, y)

	r := rng.New(0, 0, 0)
	fought := 0
	for cycle := 1; cycle <= 24 && fought == 0; cycle++ {
		// 原版 `sub_125A3` 對一支軍團做的就是這兩件事：
		// 移動／軍費（`sub_12662`／`sub_12600`）→ 倒數（`sub_1264A`）。
		ev := w.tickOneCorps(att, 0, r)
		w.tickStandoff(att)
		if ev != nil && ev.Battle != nil {
			fought = cycle
		}
	}
	if fought != 12 {
		t.Fatalf("第 %d 個巡迴週期開打，want 12（12 × 8 拍 ＝ 96 拍）", fought)
	}
}

// TestStandoffAtEnemyCityIsSiege 釘住 §1.35：下一格是**別人的據點**時
// 走的是同一套倒數，倒數完打的是攻城不是野戰。
func TestStandoffAtEnemyCityIsSiege(t *testing.T) {
	w, att, def := standoffFixture(t)
	node := w.Corps[def].Node
	if army.KindOf(node) != army.CityNode {
		t.Skip("守方不在據點上")
	}
	w.Cities[node].Owner = w.Corps[def].Faction
	w.ClearMarchRoute(att)
	c := &w.Corps[att]
	c.X, c.Y = w.Cities[node].X-1, w.Cities[node].Y
	c.TargetNode = node
	c.TargetX, c.TargetY = w.Cities[node].X, w.Cities[node].Y
	c.Interval, c.Timer = 1, 1
	c.Standoff, c.Countdown = false, 0

	r := rng.New(0, 0, 0)
	var got *CorpsEvent
	for cycle := 1; cycle <= 24 && got == nil; cycle++ {
		ev := w.tickOneCorps(att, 0, r)
		w.tickStandoff(att)
		if ev != nil && ev.Battle != nil {
			got = ev
		}
		if cycle < 12 && (c.X != w.Cities[node].X-1 || c.Y != w.Cities[node].Y) {
			t.Fatalf("第 %d 個週期就走進城了 (%d,%d)——對峙期間軍團不動",
				cycle, c.X, c.Y)
		}
	}
	if got == nil {
		t.Fatal("對峙走完沒有打起來")
	}
	if got.Mode != combat.Siege {
		t.Errorf("打成 %v，want 攻城——據點那一條走的是 sub_12880", got.Mode)
	}
}

// TestStandoffClearsWhenBlockerLeaves 釘住「位元 5 是每個週期重算的狀態」
// （`sub_125A3` 的 `and [si],0DFh`）：擋路的走了，對峙就散了，`+0x03` 歸零。
func TestStandoffClearsWhenBlockerLeaves(t *testing.T) {
	w, att, def := standoffFixture(t)
	x, y := emptyPair(t, w)
	faceOff(w, att, def, x, y)

	r := rng.New(0, 0, 0)
	w.tickOneCorps(att, 0, r)
	w.tickStandoff(att)
	if !w.Corps[att].Standoff {
		t.Fatal("第一次撞上沒進對峙，後面驗不到解除")
	}

	// 擋路的那一支讓開。
	w.Corps[def].X, w.Corps[def].Y = x+5, y+5
	w.tickOneCorps(att, 0, r)
	w.tickStandoff(att)
	c := &w.Corps[att]
	if c.Standoff {
		t.Error("擋路的走了，位元 5 還留著")
	}
	if c.Countdown != 0 {
		t.Errorf("+0x03 = %d，want 0——`sub_1264A` 沒卡住就歸零", c.Countdown)
	}
	if c.X != x || c.Y != y {
		t.Errorf("軍團停在 (%d,%d)，路空了應該走到 (%d,%d)", c.X, c.Y, x, y)
	}
}

// TestOwnCorpsDoesNotBlock 釘住 `sub_12831` 的放行條件：
// 佔著那一格的是**自己人**就照常走進去（原版允許軍團疊同格）。
//
// ⚠ 這一條也擋住「掃出所有敵人」那種寫法——原版找到第一支就停，
// 是自己人就 `stc` 放行，不會再往後看。
func TestOwnCorpsDoesNotBlock(t *testing.T) {
	w, att, def := standoffFixture(t)
	x, y := emptyPair(t, w)
	faceOff(w, att, def, x, y)
	// 讓擋路的那一支改成自己人。
	w.Corps[def].Faction = w.Corps[att].Faction

	w.tickOneCorps(att, 0, rng.New(0, 0, 0))
	c := &w.Corps[att]
	if c.Standoff {
		t.Error("自己人擋不住，不該進對峙")
	}
	if c.X != x || c.Y != y {
		t.Errorf("軍團停在 (%d,%d)，應該疊到 (%d,%d) 上", c.X, c.Y, x, y)
	}
}
