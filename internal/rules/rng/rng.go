// Package rng 是原版的亂數產生器。
//
// 出自 `KI.EXE` 的 `sub_1ECE0`（取數）與 `sub_1EC82`（建表與播種）。
// 規則層每個套件都收一個 `Rand` 介面，在這之前 cmd/wlsim 用的是
// 一個線性同餘產生器充數——**那只保證可重現，不保證分佈與原版一樣**。
// 這個套件把它換掉。
//
// 為什麼值得花力氣：傷亡量、災害、外交官成效、武將逃脫，
// 全部是「拿一個 byte 去比門檻」。門檻抄對了但分佈不同，
// 長期行為還是會偏掉，而且偏得很難查。
package rng

import "time"

// step 是計數器每次的增量（原版 `add byte_1ECFC, 89h`）。
// 0x89 是奇數，所以計數器本身會走完 256 個值才回頭。
const step = 0x89

// 播種時兩個索引的步長（原版 `add bl, 4Fh` 與 `add dl, 89h`）。
const (
	shuffleStepI = 0x4F
	shuffleStepJ = 0x89
)

// Rand 是原版的產生器狀態：一張 256 byte 的置換表，加上兩個 byte。
//
//	s ← T[s] + c
//	c ← c + 0x89
//	回傳 s
//
// 兩個 byte 合起來 65,536 種狀態，實測週期約 6.3 萬，
// 輸出在 0–255 上均勻，低位元也均勻。
type Rand struct {
	table [256]byte
	c, s  byte
}

// Next 取一個 0–255 的亂數。
func (r *Rand) Next() int {
	r.s = r.table[r.s] + r.c
	r.c += step
	return int(r.s)
}

// New 用時分秒播種，與原版 `sub_1EC82` 一致。
//
// 原版讀的是 BIOS 的即時時鐘（`int 1Ah`，ah=2），三個值都是 **BCD**，
// 所以「37 分」在種子裡是 `0x37` ＝ 55 而不是 37。這裡照做——
// 不是為了相容存檔（亂數狀態不存檔），而是因為 BCD 讓某些值
// （個位 A–F）永遠不出現，種子的分佈本來就不是均勻的。
// 想重現原版的行為就得連這個一起重現。
func New(hour, minute, second int) *Rand {
	return newRaw(bcd(hour), bcd(minute), bcd(second))
}

// newRaw 收的是原版真正拿到的那三個 byte（已經是 BCD 了）。
func newRaw(h, m, s byte) *Rand {
	r := &Rand{}
	for i := range r.table {
		r.table[i] = byte(i) // 先填成 T[i] = i
	}

	// 256 次交換，兩個索引各自以固定步長前進。
	// 起點都由「秒」決定——**時與分不影響洗牌，只進種子**。
	i, j := s, s+1
	for n := 0; n < 256; n++ {
		r.table[i], r.table[j] = r.table[j], r.table[i]
		i += shuffleStepI
		j += shuffleStepJ
	}

	r.c = s + m + h<<2 // 原版 `shl ch,1` 兩次之後才加
	r.s = r.c ^ s
	return r
}

// Now 用系統時間播種，等同原版開機時做的事。
func Now() *Rand {
	t := time.Now()
	return New(t.Hour(), t.Minute(), t.Second())
}

// NewFixed 用一個固定的 byte 播種，給測試與可重現的長跑用。
//
// 原版沒有這個入口——它只有時鐘。但長跑驗證需要同一個 seed 跑兩次
// 結果一樣，所以這裡把「秒」那個 byte 拉出來當種子參數。
// **不過 BCD**：這裡要的是 256 個可用的種子，不是像原版那樣
// 只有時鐘走得到的那幾十個。
func NewFixed(seed int) *Rand { return newRaw(0, 0, byte(seed)) }

// bcd 把 0–99 的十進位值轉成 BCD 的那個 byte。
// 超出範圍的值照原版的 byte 行為截斷。
func bcd(v int) byte { return byte(v/10<<4 | v%10) }
