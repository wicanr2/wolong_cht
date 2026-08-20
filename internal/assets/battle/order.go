package battle

// 戰場 UI 的兩張順序表。**它們是原版資料，不是版面選擇**，
// 所以兩個 UI 共用同一份——各抄一份的話其中一邊會把命令送錯。

// SideCommandRowCode 是指令面板由上而下第 row 列送出的命令碼。
//
// 出處：`sub_1C863` 在 (496, 280+16k) 註冊的熱區碼依序是
// 0x09／0x08／0x07／0x0A／0x0B／0x0C，而 handler `0x1C1B9` 算的是
// `命令碼 = 熱區碼 − 7`（docs/re/60 §6.1）。
//
// ⚠ **畫面順序不是命令碼順序**——第 0 列送的是命令 2。
var SideCommandRowCode = [6]int{2, 1, 0, 3, 4, 5}

// BottomSlotSquad 是底列由左到右第 i 格對應哪一個編成位置
//（原版 `cs:0xD2E4`，docs/spec/33 §1.1）。
//
// 六個編成位置是 0 主將／1 前鋒／2 左翼／3 右翼／4 左備／5 右備，
// 所以畫面上是「左翼 左備 主將 前鋒 右備 右翼」——**空間排列**。
var BottomSlotSquad = [6]int{2, 4, 0, 1, 5, 3}

// SquadSlot 是 BottomSlotSquad 的反排列：第 k 個編成位置排在第幾格。
func SquadSlot(squad int) int {
	for i, s := range BottomSlotSquad {
		if s == squad {
			return i
		}
	}
	return -1
}
