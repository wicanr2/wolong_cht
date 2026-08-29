package state

import "strconv"

// TalkNoticeVars 把一則事件通知換成 `TALK.DAT` 的變數表。
//
// 回傳 false 表示**這一則不能顯示**：原版 formatter 的位址回查不安全時，
// 顯示半句或猜一個據點都比不顯示糟（fail-closed）。
//
// decode 把記錄裡的原始 Big5 位元組換成可以畫的字。**由呼叫端提供**——
// `state` 存的是原始位元組（存檔要 byte-for-byte 寫回），
// 而怎麼解碼是呈現層的事。
//
// ⚠ 桌面版與手機版共用這一份。每個 marker 的語意都是從原版反推的
//（`\3` 是**君主名**不是勢力編號、`\6` 只調 X 不輸出字元、
// `\2` 在 formatter 有效時要走回查而不是據點名），各寫一份必然會有一邊
// 把某個 marker 印成錯的東西。
func (w *World) TalkNoticeVars(n TalkNotice, decode func(string) string) (map[byte]string, bool) {
	if decode == nil {
		decode = func(s string) string { return s }
	}
	vars := make(map[byte]string, 6)
	// 原版 handler（線性位址 0001097E）只消耗一個 formatter 參數並調整 X
	// 位置，不輸出字元；保留成空字串，不要把排版控制印成「6」。
	vars['6'] = ""
	if w.Player >= 0 && w.Player < len(w.Factions) {
		advisor := w.Factions[w.Player].Advisor
		if advisor >= 0 && advisor < len(w.Generals) && w.Generals[advisor].Alive {
			// marker \4（00010939）取玩家勢力的軍師姓名。
			vars['4'] = decode(w.Generals[advisor].Name)
		} else if name, _ := w.AdvisorNameRaw(); name != "" {
			// 自訂軍師（+0x02 ＝ 0x7F）：名字在區塊 +0x52A2（docs/formats/10 §4）。
			vars['4'] = decode(name)
		}
	}
	if n.RawFormatterWordValid {
		if n.RawFormatterWord < 0 || n.RawFormatterWord > 0xFFFF {
			return nil, false
		}
		raw, ok := w.ResolveTalkFormatter2(uint16(n.RawFormatterWord))
		if !ok {
			return nil, false
		}
		vars['2'] = decode(string(raw))
	} else if n.City >= 0 && n.City < len(w.Cities) {
		// `TALK.DAT` 的 marker 是 ASCII '2'，不是 state 裡的數值欄位 2。
		vars['2'] = decode(w.Cities[n.City].Name)
	}
	if n.General >= 0 && n.General < len(w.Generals) {
		vars['1'] = decode(w.Generals[n.General].Name)
	}
	if n.Faction >= 0 && n.Faction < len(w.Factions) {
		// marker \3 顯示的是該勢力**君主名**（「{3}勢力」），
		// 不是把勢力編號轉成文字。
		vars['3'] = decode(w.LordName(n.Faction))
	}
	if n.Amount >= 0 {
		// marker \7 由 `sub_1062F` 以十進位繪製；這裡只保留數值語意。
		vars['7'] = strconv.Itoa(n.Amount)
	}
	return vars, true
}

// TalkNoticeSeq 是重複標記的依序取值（`text.Table.LinesSeq`，docs/spec/106）。
// 沒有重複標記的訊息回 nil，呼叫端照舊只用 TalkNoticeVars 那張 map。
func (w *World) TalkNoticeSeq(n TalkNotice, decode func(string) string) map[byte][]string {
	if decode == nil {
		decode = func(s string) string { return s }
	}
	var out map[byte][]string
	put := func(marker byte, ids []int, name func(int) string) {
		if len(ids) == 0 {
			return
		}
		vals := make([]string, 0, len(ids))
		for _, id := range ids {
			vals = append(vals, name(id))
		}
		if out == nil {
			out = make(map[byte][]string, 2)
		}
		out[marker] = vals
	}
	put('1', n.SeqGenerals, func(id int) string {
		if id < 0 || id >= len(w.Generals) {
			return ""
		}
		return decode(w.Generals[id].Name)
	})
	put('3', n.SeqFactions, func(id int) string {
		if id < 0 || id >= len(w.Factions) {
			return ""
		}
		return decode(w.LordName(id))
	})
	return out
}
