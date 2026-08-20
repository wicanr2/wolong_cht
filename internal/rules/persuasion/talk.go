package persuasion

// 進言那幾段對白在 `TALK.DAT` 裡的位置。
//
// 這些算式是從 `sub_13830`／`sub_13B5A`／`sub_13BA9` 反組譯出來的
//（docs/spec/44），不是排版選擇——**索引怎麼算是原版的規則**，
// 換一個呈現層不會換算式。所以放在規則層讓兩個 UI 共用；
// 各自抄一份的話，某一種結果的措辭會在其中一邊悄悄錯掉。

// TalkBase 是三個進言指令的 TALK 起點——`sub_16405`／`sub_164F1`／
// `sub_16623` 傳給 `sub_13830` 的 `cx`。
//
// ⚠ **三組措辭各不相同**，不能共用一組。
func TalkBase(c Command) int {
	switch c {
	case CeaseFire:
		return 0x96 // 150
	case Cooperate:
		return 0xD6 // 214
	default:
		return 0x56 // 86，敵對（進兵）
	}
}

// TalkReplyIndex 是君主回答的 TALK 索引：`base + 4 + 結果碼×3`
//（`sub_13830` 的 `cx = base + al×3 + 4`）。結果碼 ≥ 4 一律用 #83。
//
// 每個位置佔三則，因為 `sub_13C99` 會把君主的**說話型**加進索引。
func TalkReplyIndex(base int, r Reaction, variant int) int {
	if r >= SameFaction {
		return 0x53 + variant // #83「我想軍師並不是來談笑的。」
	}
	return base + 4 + 3*int(r) + variant
}

// TalkReasonBase 是說服迴圈的起點。原版在進迴圈前 `add [bp+0], 10h`
//（`sub_13830`），所以理由那一段全部相對於 base + 16 —— 而 base + 16
// 正好是那三則五選一的選單（#102／#166／#230）。
func TalkReasonBase(c Command) int { return TalkBase(c) + 0x10 }

// TalkReasonSlot 是理由在這個指令的選單裡排第幾（0–4，4 是撤回）。
// **順序也是資料**——原版的索引算式直接吃這個位置。
func TalkReasonSlot(c Command, r Reason) int {
	for i, o := range Options(c) {
		if o == r {
			return i
		}
	}
	return len(Options(c)) - 1 // 找不到就當撤回，不要算出界
}

// TalkReasonReply 是君主對一個理由的反應（原版 `sub_13BA9` 的結尾）：
//
//	撤回（第 5 項）      base + 42
//	同一個理由講第二次   base + 45
//	否則                 base + 位置×9 + 結果×3 + 6
//
// 結果碼 0 ＝ 理由不成立、1 ＝ 湊夠了、2 ＝ 還要再一個。
// 每個位置佔三則是君主的**說話型**變體（`sub_13C99` 的 `add cx, ax`）。
func TalkReasonReply(base, slot int, out Outcome, repeat bool, variant int) int {
	switch {
	case out == Withdrawn:
		return base + 42 + variant
	case repeat:
		return base + 45 + variant
	}
	code := 2 // Continue：還要再一個
	switch out {
	case Failed:
		code = 0
	case Agreed:
		code = 1
	}
	return base + slot*9 + code*3 + 6 + variant
}

// 遷都與請求出陣**不走說服迴圈**——君主只做一次驗收，
// 三句話定案（`sub_16909`／`sub_1699E` 傳給 `sub_13B08` 的 `cx`，
// docs/spec/49 §1）。
const (
	TalkRelocateBase = 0x182 // 386
	TalkSortieBase   = 0x18C // 396
)

// TalkMenu 是進言那五項的選單訊息（`sub_16224` 的 `cx = 4Dh`）。
const TalkMenu = 0x4D
