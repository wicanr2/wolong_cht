// Package threat 是據點的威脅偵測與 AI 出兵額度（原版 `sub_13FA9`／`sub_14575`）。
//
// 純規則層，不認識畫面也不認識 World。證據在 docs/re/44。
//
// 原版每 tick 只處理**一個**據點（`sub_13EFD`，游標 `word_10D1E`），
// 192 個 tick 掃完一輪——**AI 的反應速度是被這個掃描週期限制的**，
// 不是被判斷邏輯限制的。呼叫端要照這個節奏跑，不要一次掃全圖。
package threat

// Neutral 是「中立／無所屬」的勢力編號（原版寫死的 0x18 ＝ 24）。
const Neutral = 0x18

// NoTarget 是「沒有侵攻目標」（勢力記錄 +0x19 的 0xFF）。
const NoTarget = 0xFF

// Slots 是一個據點的鄰接槽數（記錄 +0x1C–+0x1F）。
const Slots = 4

// PeaceBit 是交友度的最高位元。**1 ＝ 和平**，不是交戰。
//
// 這個極性是從資料定下來的（四個劇本每一格都是 1，而開局沒有人有侵攻目標），
// 而 `sub_13FA9` 對 `>= 0x80` 的鄰居直接跳過又獨立印證了一次——
// **和平的鄰居不算威脅**。見 docs/formats/08 §1.55 與 docs/re/44 §2.2。
const PeaceBit = 0x80

// Neighbour 是一個鄰接槽。
type Neighbour struct {
	// Site 是鄰居的據點編號。負數表示這個槽沒有鄰居（原版的 0xFF 哨兵）。
	Site int
	// Owner 是鄰居的所屬勢力（記錄 +0x01），Neutral ＝ 中立。
	Owner int
	// Occupancy 是停在鄰居那一格的軍團數（記錄 +0x18）。
	Occupancy int
	// Friendship 是「本據點的勢力」對「鄰居的勢力」的交友度
	// （交友度表的一格，有向）。
	Friendship int
}

// Result 是掃完四個鄰接槽的結果。
type Result struct {
	// Threatened 對應據點記錄 +0x00 的 bit 7。
	Threatened bool
	// Specific 對應 bit 6：威脅裡有**本勢力侵攻目標**的據點。
	//
	// ⚠ 分野是「有沒有具體目標」，不是遠近。
	Specific bool
	// Level 是據點 +0x14 ＝ 鄰接敵方據點的軍團數總和。
	Level int
	// Targets 是侵攻目標勢力的鄰居據點編號，最多四個（原版那 16 B 緩衝）。
	Targets []int
}

// Scan 重現 `sub_13FA9` ＋ `sub_14028` 的判斷（docs/re/44 §2）。
//
// enemyNeighbours 是據點記錄 +0x1B；它是 0 就直接回，原版連掃都不掃。
// invasionTarget 是本勢力的侵攻目標（勢力記錄 +0x19）。
func Scan(owner, invasionTarget, enemyNeighbours int, ns []Neighbour) Result {
	var r Result
	if enemyNeighbours == 0 {
		return r
	}
	for i, n := range ns {
		if i >= Slots || n.Site < 0 {
			break
		}
		if n.Owner == owner {
			continue
		}
		// 中立的鄰居不看交友度（勢力表只有 22 筆，24 是哨兵），
		// 但仍然要比對侵攻目標——原版是 `jz` 跳進比較，不是跳過整輪。
		if n.Owner != Neutral {
			if n.Friendship&PeaceBit != 0 {
				continue
			}
			r.Threatened = true
			r.Level += n.Occupancy
		}
		if invasionTarget != NoTarget && n.Owner == invasionTarget {
			r.Threatened = true
			r.Specific = true
			r.Targets = append(r.Targets, n.Site)
		}
	}
	return r
}

// Requested 是據點想要幾支援軍：`+0x14 ＋ 2 − +0x18`（原版 `sub_14057`）。
// 不足 1 就回 0 ＝ 不求援。
func Requested(level, occupancy int) int {
	n := level + 2 - occupancy
	if n < 1 {
		return 0
	}
	return n
}

// CorpsCap 是一個勢力的軍團數上限（原版 `sub_14575`）。
//
//	資金 >> 8 <= 0xA0（資金 <= 40,960）  →  5
//	否則                                →  資金 >> 13 ＝ 資金 ÷ 8,192
//
// 兩段在交界處連續（40,960 ÷ 8,192 ＝ 5），所以整條規則就是
// `max(5, 資金 ÷ 8192)`。原版在大額那一支還有一行 `add bl, 2`，
// **下一個指令就把 bl 覆蓋掉**，對結果沒有影響——不要照抄死碼。
func CorpsCap(funds int) int {
	if funds < 0 {
		funds = 0
	}
	if cap := funds >> 13; cap > 5 {
		return cap
	}
	return 5
}

// Budget 是「還能新編幾支」＝ 上限減掉現有軍團數。負數回 0。
func Budget(funds, corps int) int {
	n := CorpsCap(funds) - corps
	if n < 0 {
		return 0
	}
	return n
}

// EnemyMask 由四個鄰接槽算出據點記錄 +0x00 的低 4 位
// （**哪幾個鄰接槽屬於別的勢力**，不是「哪幾個方向有鄰接」）。
//
// 原版是 `sub_1890A` 在據點換手時逐位設／清，並同步加減 +0x1B；
// 這裡直接由現況算，結果一樣而且不會漂。有沒有鄰居由 +0x1C–+0x1F 的
// 哨兵決定，與這四位無關（docs/re/44 §5）。
func EnemyMask(owner int, ns []Neighbour) (mask, count int) {
	for i, n := range ns {
		if i >= Slots || n.Site < 0 || n.Owner == owner {
			continue
		}
		mask |= 1 << uint(i)
		count++
	}
	return mask, count
}
