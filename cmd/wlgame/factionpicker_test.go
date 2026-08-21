package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// ⭐ 版面全部釘死在原版的直接座標（docs/spec/35 §2.5.2）。
// 三塊互相印證：外框 (512,192)–(639,399)，兩塊填色與熱區都比外框
// **四邊各內縮 8 px**——那 8 px 就是框邊。這一條擋「某一個數字被改掉」。
func TestFactionPickerGeometryMatchesRawConstants(t *testing.T) {
	const border = 8
	if pickerListX != pickerWinX+border {
		t.Errorf("清單左緣 %d，預期外框 %d + %d", pickerListX, pickerWinX, border)
	}
	if right := pickerListX + pickerListW; right != pickerWinX+pickerWinW-border {
		t.Errorf("清單右緣 %d，預期 %d", right, pickerWinX+pickerWinW-border)
	}
	if pickerTitleY != pickerWinY+border {
		t.Errorf("標題列上緣 %d，預期 %d", pickerTitleY, pickerWinY+border)
	}
	if bottom := pickerListY + pickerListH; bottom != pickerWinY+pickerWinH-border {
		t.Errorf("清單下緣 %d，預期 %d", bottom, pickerWinY+pickerWinH-border)
	}
	// 標題列剛好一列高，清單接在它下面。
	if pickerListY != pickerTitleY+pickerRowH {
		t.Errorf("清單上緣 %d，預期標題列 %d + %d", pickerListY, pickerTitleY, pickerRowH)
	}
	// 11 列 × 16 px 剛好填滿清單。
	if pickerRows*pickerRowH != pickerListH {
		t.Errorf("%d 列 × %d ＝ %d，清單高 %d",
			pickerRows, pickerRowH, pickerRows*pickerRowH, pickerListH)
	}
	// ⭐ 欄的分界是 X ≥ 576（原版 `cmp cx, 0x240`），剛好把 112 px 切一半。
	if pickerSplitX != pickerListX+pickerListW/2 {
		t.Errorf("欄分界 %d，預期 %d", pickerSplitX, pickerListX+pickerListW/2)
	}
	// 兩欄的文字各離欄左緣 4 px。
	if pickerLeftX-pickerListX != pickerRightX-pickerSplitX {
		t.Errorf("兩欄的文字內縮不同：左 %d、右 %d",
			pickerLeftX-pickerListX, pickerRightX-pickerSplitX)
	}
	if pickerSlots != 22 {
		t.Errorf("槽位數 %d，預期 22 個勢力", pickerSlots)
	}
}

// pickerSlotAt 與 pickerRowOrigin 必須互相對得起來：
// 每一個槽位的文字起點，送回去要換回同一個槽位。
func TestFactionPickerHitTestRoundTrips(t *testing.T) {
	for n := 0; n < pickerSlots; n++ {
		x, y := pickerRowOrigin(n)
		got, ok := pickerSlotAt(x, y)
		if !ok || got != n {
			t.Errorf("槽 %d 的起點 (%d,%d) 換回 %d（ok=%v）", n, x, y, got, ok)
		}
		// 同一列的下緣（+15）仍是同一格。
		if got, ok := pickerSlotAt(x, y+pickerRowH-1); !ok || got != n {
			t.Errorf("槽 %d 的列底換回 %d（ok=%v）", n, got, ok)
		}
	}
	// ⭐ 分界的兩側是不同欄：575 是左欄第 0 列、576 是右欄第 0 列。
	if n, ok := pickerSlotAt(pickerSplitX-1, pickerListY); !ok || n != 0 {
		t.Errorf("X=%d 應該是左欄第 0 列，得到 %d", pickerSplitX-1, n)
	}
	if n, ok := pickerSlotAt(pickerSplitX, pickerListY); !ok || n != pickerRows {
		t.Errorf("X=%d 應該是右欄第 0 列（%d），得到 %d", pickerSplitX, pickerRows, n)
	}
	// 視窗外不接。
	for _, pt := range [][2]int{
		{pickerListX - 1, pickerListY},
		{pickerListX + pickerListW, pickerListY},
		{pickerListX, pickerListY - 1},
		{pickerListX, pickerListY + pickerListH},
	} {
		if _, ok := pickerSlotAt(pt[0], pt[1]); ok {
			t.Errorf("(%d,%d) 在清單外卻被接走", pt[0], pt[1])
		}
	}
}

// 原版的兩條規則：勢力要存在、**不能選自己**。
func TestFactionPickerSelectionRules(t *testing.T) {
	w := &state.World{Player: 3}
	w.Factions[3].Alive = true
	w.Factions[7].Alive = true
	g := &game{world: w}

	if g.pickerSelectable(3) {
		t.Error("自勢力不該可選（sub_15AFC 的 cmp al, cs:byte_10CFF）")
	}
	if g.pickerSelectable(5) {
		t.Error("沒活著的勢力不該可選")
	}
	if !g.pickerSelectable(7) {
		t.Error("勢力 7 活著又不是自己，應該可選")
	}
	if g.pickerSelectable(-1) || g.pickerSelectable(9999) {
		t.Error("越界的槽位不該可選")
	}
}

// ⭐ 熱區 0x16：點縮小地圖 ＝ 1 px 對 2 格，鏡頭再減 (20,12)——
// **那個偏移與開局定位共用**（docs/spec/52），不可以各寫一份。
func TestMinimapClickScrollsCamera(t *testing.T) {
	if _, _, ok := minimapWorldAt(strategyMinimapX-1, strategyMinimapY); ok {
		t.Error("地圖區左邊界外不該接")
	}
	if _, _, ok := minimapWorldAt(strategyMinimapX, strategyMinimapY+strategyMinimapImageH); ok {
		t.Error("地圖區下邊界外不該接（那一列是圖例）")
	}
	// 左上角 ⇒ 世界 (0,0)
	col, row, ok := minimapWorldAt(strategyMinimapX, strategyMinimapY)
	if !ok || col != 0 || row != 0 {
		t.Fatalf("左上角換成 (%d,%d)，預期 (0,0)", col, row)
	}
	// 右下角 ⇒ 世界 (382,254)：192×128 px × 2 減最後一格
	col, row, _ = minimapWorldAt(strategyMinimapX+strategyMinimapW-1,
		strategyMinimapY+strategyMinimapImageH-1)
	if col != 382 || row != 254 {
		t.Errorf("右下角換成 (%d,%d)，預期 (382,254)", col, row)
	}

	g := &game{}
	// 地圖中央附近：鏡頭要是那一格減 (20,12)。
	g.centreCamOn(200, 150)
	if g.camX != 200-centreCol || g.camY != 150-centreRow {
		t.Errorf("鏡頭 (%d,%d)，預期 (%d,%d)",
			g.camX, g.camY, 200-centreCol, 150-centreRow)
	}
	// 邊界要夾住，不能捲出世界外。
	g.centreCamOn(0, 0)
	if g.camX != 0 || g.camY != 0 {
		t.Errorf("左上角的鏡頭 (%d,%d)，預期夾成 (0,0)", g.camX, g.camY)
	}
	g.centreCamOn(383, 255)
	if g.camX != 384-viewCols || g.camY != 256-viewRows {
		t.Errorf("右下角的鏡頭 (%d,%d)，預期夾成 (%d,%d)",
			g.camX, g.camY, 384-viewCols, 256-viewRows)
	}
}
