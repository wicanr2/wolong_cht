package state

// 「自定」軍師（docs/formats/10 §4、docs/spec/104）。
//
// 原版的自訂軍師不是任何一位武將：勢力 +0x02 寫 0x7F，名字與肖像另外放在
// 區塊 +0x52A1／+0x52A2（`ds:5221h`／`5222h`）。六個 Big5 字前三個是軍師名、
// 後三個是別號；沒改過的格子是 `A1 D0`（空標記），清掉的格子是全形空白 `A1 40`。

const (
	advisorNameChars = advisorNameLen / 2
	advisorNameHalf  = advisorNameChars / 2 // 3：軍師名與別號各三字
)

// advisorEmptyMark 是原版判「還沒取名」的第一格（`cmp word ptr [si], 0D0A1h`）。
var advisorEmptyMark = [2]byte{0xA1, 0xD0}

// HasCustomAdvisor 回報玩家勢力是不是用自訂軍師（+0x02 ＝ 0x7F 而且取過名）。
func (w *World) HasCustomAdvisor() bool {
	if w.Player < 0 || w.Player >= len(w.Factions) || w.Factions[w.Player].Advisor != NoAdvisor {
		return false
	}
	return w.AdvisorName[0] != advisorEmptyMark[0] || w.AdvisorName[1] != advisorEmptyMark[1]
}

// AdvisorNameRaw 回傳自訂軍師的「軍師名」與「別號」，都是原始 Big5 位元組
// （與 `General.Name` 同一種表示，交給呼叫端的 `big5()` 解碼），
// 尾端的全形空白已去掉。沒有自訂軍師時兩個都是空字串。
func (w *World) AdvisorNameRaw() (name, alias string) {
	if !w.HasCustomAdvisor() {
		return "", ""
	}
	trim := func(b []byte) string {
		for len(b) >= 2 && ((b[len(b)-2] == 0xA1 && b[len(b)-1] == 0x40) ||
			(b[len(b)-2] == advisorEmptyMark[0] && b[len(b)-1] == advisorEmptyMark[1])) {
			b = b[:len(b)-2]
		}
		return string(b)
	}
	return trim(w.AdvisorName[:advisorNameHalf*2]), trim(w.AdvisorName[advisorNameHalf*2:])
}

// SetCustomAdvisor 把玩家勢力的軍師換成自訂的：+0x02 寫 0x7F（原版「確定」
// 那一步 `mov byte ptr [bx+2], 7Fh`），名字與肖像寫進區塊欄位。
// name 是 12 個 Big5 位元組（不足補全形空白）。
func (w *World) SetCustomAdvisor(portrait int, name []byte) {
	if w.Player < 0 || w.Player >= len(w.Factions) {
		return
	}
	w.Factions[w.Player].Advisor = NoAdvisor
	w.AdvisorPortrait = portrait & 0xFF
	for i := 0; i < advisorNameLen; i += 2 {
		if i+1 < len(name) {
			w.AdvisorName[i], w.AdvisorName[i+1] = name[i], name[i+1]
		} else {
			w.AdvisorName[i], w.AdvisorName[i+1] = 0xA1, 0x40
		}
	}
}
