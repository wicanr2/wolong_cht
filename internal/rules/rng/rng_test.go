package rng

import (
	"math"
	"testing"
)

// 播種只是一連串交換，所以表在任何時候都是 0–255 的置換。
// 這一條擋的是「洗牌寫成賦值而不是交換」那類錯誤——
// 那會讓某些值永遠取不到，而輸出看起來仍然很亂。
func TestTableStaysAPermutation(t *testing.T) {
	for s := 0; s < 256; s++ {
		r := NewFixed(s)
		var seen [256]bool
		for _, v := range r.table {
			if seen[v] {
				t.Fatalf("種子 %d：表裡有重複值 %d，不是置換", s, v)
			}
			seen[v] = true
		}
	}
}

// 從 KI.EXE 讀出來的演算法，用獨立寫的參照模型算出前幾個輸出。
// 這是整個套件唯一的「對照原版」錨點，改動實作時它會先炸。
func TestMatchesReference(t *testing.T) {
	got := take(New(20, 35, 27), 12)
	want := []int{40, 52, 82, 172, 235, 220, 70, 116, 55, 140, 89, 77}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("New(20,35,27) 第 %d 個 ＝ %d，應為 %d\n實得 %v", i, got[i], want[i], got)
		}
	}
}

// 時分秒是 BCD。「35 分」在種子裡是 0x35 ＝ 53，不是 35。
func TestSeedIsBCD(t *testing.T) {
	if got := bcd(35); got != 0x35 {
		t.Errorf("bcd(35) ＝ 0x%02X，應為 0x35", got)
	}
	// 洗牌只吃「秒」，時與分只進計數器 —— 換掉分鐘不該改變表。
	a, b := New(0, 0, 27), New(0, 59, 27)
	if a.table != b.table {
		t.Error("換掉分鐘不該影響置換表")
	}
	if a.c == b.c {
		t.Error("換掉分鐘應該改變計數器初值")
	}
}

// 輸出要蓋滿 0–255 而且大致均勻。傷亡、災害、逃脫全是拿它比門檻，
// 分佈偏掉的話長期行為會跟著偏，而且很難查。
func TestOutputIsUniform(t *testing.T) {
	const n = 65536
	var hist [256]int
	r := NewFixed(0x27)
	for i := 0; i < n; i++ {
		hist[r.Next()]++
	}
	sum := 0
	for v, c := range hist {
		if c == 0 {
			t.Fatalf("值 %d 一次都沒出現", v)
		}
		sum += v * c
	}
	if mean := float64(sum) / n; math.Abs(mean-127.5) > 2 {
		t.Errorf("平均 %.2f，應接近 127.5", mean)
	}
	// 卡方：256 格、期望值 256，自由度 255。臨界值取寬一點的 400。
	chi := 0.0
	for _, c := range hist {
		d := float64(c) - float64(n)/256
		chi += d * d / (float64(n) / 256)
	}
	if chi > 400 {
		t.Errorf("卡方 %.1f 過高，分佈不夠均勻", chi)
	}
}

// 低位元也要均勻。規則層大量使用 `rand & 7`、`rand & 3`、`rand & 0x0F`
// ——高位元再亂，低位元有偏差一樣會壞事。
func TestLowBitsAreUniform(t *testing.T) {
	const n = 1 << 16
	for _, mask := range []int{1, 3, 7, 0x0F} {
		buckets := make([]int, mask+1)
		r := NewFixed(0x11)
		for i := 0; i < n; i++ {
			buckets[r.Next()&mask]++
		}
		want := n / (mask + 1)
		for v, c := range buckets {
			if d := c - want; d < -want/16 || d > want/16 {
				t.Errorf("mask 0x%X：值 %d 出現 %d 次，期望約 %d", mask, v, c, want)
			}
		}
	}
}

// 週期要遠大於一場遊戲會用掉的次數。狀態是兩個 byte，上限 65,536。
func TestPeriodIsLong(t *testing.T) {
	r := NewFixed(0x27)
	type st struct{ c, s byte }
	seen := map[st]int{}
	for i := 0; ; i++ {
		k := st{r.c, r.s}
		if first, ok := seen[k]; ok {
			if p := i - first; p < 40000 {
				t.Errorf("週期只有 %d，太短", p)
			}
			return
		}
		seen[k] = i
		r.Next()
	}
}

// 同一個種子跑兩次要一模一樣——長跑回歸比對靠這個。
func TestReproducible(t *testing.T) {
	a, b := take(NewFixed(9), 500), take(NewFixed(9), 500)
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("第 %d 個不同：%d ≠ %d", i, a[i], b[i])
		}
	}
	if c := take(NewFixed(10), 500); equal(a, c) {
		t.Error("不同種子產生了相同的序列")
	}
}

func take(r *Rand, n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = r.Next()
	}
	return out
}

func equal(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
