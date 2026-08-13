package state

import "github.com/wicanr2/wolong_cht/internal/rules/economy"

// 這兩個 TALK index 來自已讀出的 DOS/V 月結俘虜處理鄰近路徑：
// sub_15940 以 CX=41h／42h 顯示「俘虜逃走／歸降」。它們是可重用的
// 文字證據，不代表低碼 0x0A producer 已被原版反組譯證實。
const (
	approximateEvent10EscapeTalk = 0x41
	approximateEvent10JoinTalk   = 0x42
)

// produceApproximateEvent10 是事件 10 的**替代自然 producer**。
//
// 已證實的部分只有 `sub_13496`：事件字高是 General、Param 是 TALK index，
// 每時 dispatcher 才會消費。原版到底在哪裡寫入低碼 0x0A 仍 unknown；這裡
// 不偽造原版呼叫者，而是把已存在的玩家俘虜狀態接到相同 raw queue contract：
// 月結時每月最多選一名玩家勢力目前收容的武將，依固定 RNG 邊界近似逃走／
// 歸降，寫入一筆事件 10，下一個每時節拍再轉成 TALK。狀態變更與 queue 寫入
// 以單筆為原子單位，滿 queue 時不改動武將。
func (w *World) produceApproximateEvent10(rng economy.Rand) bool {
	if w == nil || !w.approximateEvent10 || rng == nil ||
		w.Player < 0 || w.Player >= numFactions || !w.Factions[w.Player].Alive {
		return false
	}

	for id := range w.Generals {
		g := &w.Generals[id]
		if !g.Alive || g.Faction != w.Player || g.Captor == noFaction ||
			w.hasQueuedEvent10(id) {
			continue
		}
		// +0x18 的倒數是原版 sub_1585F 的月度閘；近似路徑保留
		// 「未到行動月只遞減」的可觀測部分。
		if g.Timer > 0 {
			g.Timer--
			continue
		}

		var talk int
		next := *g
		switch roll := rng.Next() & 0xFF; {
		case roll < 0x20:
			// 原版鄰近路徑把逃走者移到無主狀態；remake 的
			// noFaction 是同一層的在野 sentinel。
			next.Faction = noFaction
			next.Captor = noFaction
			next.Posted = false
			talk = approximateEvent10EscapeTalk
		case roll < 0x40:
			// 近似「歸降我軍」：保留目前玩家勢力，只清掉
			// 俘虜來源，讓下一輪可正常編成／派駐。
			next.Captor = noFaction
			next.Posted = false
			talk = approximateEvent10JoinTalk
		default:
			// 原版 sub_15940 的其餘亂數區間不產生這兩種結果。
			continue
		}

		if !w.queueFullEvent(id, 0x0A, uint16(talk), 0) {
			return false
		}
		*g = next
		return true
	}
	return false
}

func (w *World) hasQueuedEvent10(general int) bool {
	for _, e := range w.events {
		if byte(e.Code) == 0x0A && int(e.Code>>8) == general {
			return true
		}
	}
	return false
}
