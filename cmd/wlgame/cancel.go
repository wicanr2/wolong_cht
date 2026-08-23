package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// cancelPressed 是「退回上一層」的統一判定（docs/spec/73）。
//
// ⭐ **原版的取消在輸入層，不在每個視窗裡。** 每個模態畫面都停在同一組
// 等待常式上（`sub_121E7`／`sub_17E1F`／`sub_18B7C`／`sub_18DC8`／`sub_18E5A`），
// 而那些常式**回傳 CF=1 就是取消**——右鍵按在哪裡都一樣，它不是熱區的功能。
//
// ⚠ 先前這條規則散成十二份實作，於是七個面板漏掉右鍵、只認 ESC。
// 漏掉的方式是安靜的：面板打得開、ESC 也能關，只有拿滑鼠玩的人會踩到。
// **一條規則只留一份實作**（CLAUDE.md §7 第 6 條）。
//
// ⚠ **不要在全域攔截。** 面板有巢狀（進言 → 數值輸入），全域攔截會一次
// 關掉兩層；原版的 CF=1 只回到當前那一層的呼叫端，一次退一層。
func cancelPressed() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) ||
		pressed(ebiten.KeyEscape)
}

// cancelled 是面板實際問的那一支。
//
// ⚠ **測試不能直接呼叫 `cancelPressed`**：`inpututil` 讀的是 Ebiten 的
// 全域輸入狀態，無頭測試裡永遠是 false，於是「面板關不關」這件事
// **測起來永遠通過**——那正是這一輪要修的那種安靜失敗。
// 所以留一個可注入的欄位，讓測試驗的是行為不是「有沒有寫這一行」。
func (g *game) cancelled() bool {
	if g.cancelFn != nil {
		return g.cancelFn()
	}
	return cancelPressed()
}
