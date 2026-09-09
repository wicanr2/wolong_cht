package state

import "testing"

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
