package isoview

import (
	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"testing"
)

func TestProjectileSourceIndex(t *testing.T) {
	for _, tc := range []struct {
		name string
		p    tactical.ProjectileView
		want int
	}{
		{name: "west-east normal", p: tactical.ProjectileView{Direction: tactical.West}, want: 0x210},
		{name: "north-south normal", p: tactical.ProjectileView{Direction: tactical.North}, want: 0x211},
		{name: "special first frame", p: tactical.ProjectileView{Special: true, Direction: tactical.South | 0x80}, want: 0x214},
		{name: "special second frame", p: tactical.ProjectileView{Special: true, SpecialFrame: 1, Direction: tactical.South | 0x80}, want: 0x215},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ProjectileSourceIndex(tc.p); got != tc.want {
				t.Fatalf("raw 圖號 = %#x，預期 %#x", got, tc.want)
			}
		})
	}
}

func TestBattleCameraStartsAtOriginalWorldOrigin(t *testing.T) {
	// 鏡頭存的是**原版的框**（含表頭那一列），所以初值就是 sub_199F3 的
	// (0x24, 0x0E)；地形列號在 isoProject／cellOffset 的入口才換算
	// （docs/spec/57 §2）。
	if battleCamInitX != 0x24 || battleCamInitY != 0x0e {
		t.Fatalf("鏡頭初值 = (%#x,%#x)，預期 (0x24,0x0e)", battleCamInitX, battleCamInitY)
	}
	v := &View{camWorldX: battleCamInitX, camWorldY: battleCamInitY}
	v.applyCameraOrigin()
	if v.camCol != 0x32 || v.camRow != 0x15 {
		t.Fatalf("sub_1DC9D 投影原點 = (%#x,%#x)，預期 (0x32,0x15)",
			v.camCol, v.camRow)
	}
	// 地形列號 13（＝原版的 14）投影出來要與鏡頭同一格——
	// 這是兩個框有對上的正對照。
	if col, row := isoProject(battleCamInitX, battleCamInitY-originalRowBias, 0); col != v.camCol ||
		row != v.camRow {
		t.Fatalf("地形列號 %d 投影 = (%#x,%#x)，預期與鏡頭原點相同 (%#x,%#x)",
			battleCamInitY-originalRowBias, col, row, v.camCol, v.camRow)
	}

	// 相機不再從兵座標推導。改變任意兵的位置，不應改寫保存的 world origin。
	v.applyCameraOrigin()
	if v.camWorldX != battleCamInitX || v.camWorldY != battleCamInitY {
		t.Fatalf("相機 world origin 被繪圖改寫：(%#x,%#x)", v.camWorldX, v.camWorldY)
	}
}

// TestBattleCursorCrossMatchesCameraOffsets 釘住游標十字與鏡頭之間的兩個偏移。
//
// 正對照是原版自己的初值：`sub_199F3` 同時設鏡頭 (0x24,0x0E) 與十字
// (0x20,0x21)，而縮圖點選（`0001C103`／`0001C106`）用的是
// 「X − 4」與「原版 Y ＋ 0x13」。兩者要能互推，偏移才是對的。
func TestBattleCursorCrossMatchesCameraOffsets(t *testing.T) {
	if got := battleCamInitX + cursorBiasX; got != 0x20 {
		t.Fatalf("十字 X 初值 = %#x，預期 0x20", got)
	}
	if got := battleCamInitY + cursorBiasY; got != 0x21 {
		t.Fatalf("十字 Y 初值 = %#x，預期 0x21", got)
	}
	// 點在十字自己身上（原版量到的 562,142，docs/playtest/40 §3.1）
	// 應該把鏡頭與十字都留在初值——這是這條換算的閉環檢查。
	v := &View{}
	v.SetCameraFromMiniMap(562, 142)
	if v.camWorldX != battleCamInitX || v.camWorldY != battleCamInitY {
		t.Fatalf("點十字之後鏡頭 = (%d,%d)，預期 (%d,%d)",
			v.camWorldX, v.camWorldY, battleCamInitX, battleCamInitY)
	}
	if v.cursorX != 0x20 || v.cursorY != 0x21 {
		t.Fatalf("點十字之後十字 = (%#x,%#x)，預期 (0x20,0x21)", v.cursorX, v.cursorY)
	}
}

func TestBattleDisplayListSplitsSoldiersIntoOriginalRawUnits(t *testing.T) {
	v := &View{camCol: 0, camRow: 0}
	b := &tactical.Battle{}
	b.Sides[0].Soldiers[0] = tactical.Soldier{Alive: true, Kind: tactical.General,
		X: 20, Y: 24, Z: 0}
	b.Sides[1].Soldiers[0] = tactical.Soldier{Alive: true, Kind: tactical.General,
		X: 22, Y: 18, Z: 0}

	entries := v.buildDisplayList(b)
	var units []battleDisplayEntry
	for _, e := range entries {
		if e.kind == displayRawUnit {
			units = append(units, e)
		}
	}
	if len(units) != 4 {
		t.Fatalf("兵 raw-unit entries = %d，預期 4", len(units))
	}
	for i := 0; i < len(units); i += 2 {
		upper, lower := units[i], units[i+1]
		if upper.raw != lower.raw+1 {
			t.Fatalf("sub_1DA1C unit 配對 = (%#x,%#x)，預期奇數上半／偶數下半",
				upper.raw, lower.raw)
		}
		if lower.cellRow != upper.cellRow-1 || lower.layer != upper.layer+1 {
			t.Fatalf("sub_1DA1C 鄰列寫入錯誤：upper=%+v lower=%+v", upper, lower)
		}
		if upper.lane != 1 || lower.lane != 1 {
			t.Fatalf("人物應落在 lane 1：upper=%d lower=%d", upper.lane, lower.lane)
		}
	}
}

func TestBattleDisplayGridOriginalDimensionsAndSlotFormula(t *testing.T) {
	if got := battleDisplayGridRows * battleDisplayRowSize; got != 0x7800 {
		t.Fatalf("顯示格總長 = %#x，預期 sub_1D971 清除的 0x7800", got)
	}
	if got := battleDisplayGridCols * battleDisplayCellSize; got != battleDisplayRowSize {
		t.Fatalf("顯示格列寬 = %#x，預期 %#x", got, battleDisplayRowSize)
	}
	for _, tc := range []struct {
		col, row, z, lane, want int
	}{
		{col: 0, row: 0, z: 0, lane: 0, want: 4},
		{col: 0, row: 0, z: 0, lane: 1, want: 6},
		{col: 1, row: 0, z: 1, lane: 1, want: 0x2a},
		{col: 3, row: 2, z: 6, lane: 0, want: 0x87c},
	} {
		if got := battleDisplaySlotOffset(tc.col, tc.row, tc.z, tc.lane); got != tc.want {
			t.Errorf("slot(%d,%d,%d,%d)=%#x，預期 %#x",
				tc.col, tc.row, tc.z, tc.lane, got, tc.want)
		}
	}
}

func TestBattleDisplayGridKeepsFirstProducerInOccupiedSlot(t *testing.T) {
	a := battleDisplayEntry{kind: displayRawUnit, cellCol: 2, cellRow: 3,
		layer: 1, lane: 1, raw: 0x123}
	b := a
	b.raw = 0x456
	grid := makeDisplayGrid([]battleDisplayEntry{a, b})
	got := grid[3][2][1][1]
	if !got.set || got.entry.raw != a.raw {
		t.Fatalf("已佔用槽被後一個 producer 覆蓋：%+v", got)
	}
}

func TestUnfoldDisplayTileMatchesSub1E011Quadrants(t *testing.T) {
	want := [][2]int{
		{0, 0}, {battle.SubTileW, 0},
		{0, battle.SubTileH / 4}, {battle.SubTileW, battle.SubTileH / 4},
	}
	for i, p := range want {
		x, y := displayUnfoldDestination(i)
		if x != p[0] || y != p[1] {
			t.Errorf("sub_1E011 part %d destination=(%d,%d)，預期 (%d,%d)",
				i, x, y, p[0], p[1])
		}
	}
}

func TestSoldierPoseUsesPerRecordBit(t *testing.T) {
	for _, pose := range []uint8{0, 1} {
		s := tactical.Soldier{Kind: tactical.Infantry, Facing: tactical.East, PoseStep: pose}
		flags := int(s.PoseStep) & battle.PoseFlagStep
		got := battle.SpriteFor(int(s.Kind), s.Facing, flags)
		want := int(s.Kind) + tactical.East*battle.FacingStride + int(pose)
		if got != want {
			t.Fatalf("PoseStep=%d 圖號=%d，預期 %d", pose, got, want)
		}
	}
}

// TestCellOffsetFollowsOriginalWalk 釘住 `sub_1DC9D` 的走法：
// 顯示格第 r 列從 (camX−r, camY+r) 起，交替走 y+1 與 x+1，一步一格。
// `cellOffset` 對那條路徑上的每一格都必須回 (s, r)。
//
// ⭐ **鏡頭的奇偶要一起測。** 原版縮圖點選的公式
// `(((x − 0x1F0) >> 1) | 1) − 0x13` 把 `camWorldY` 強制成奇數，
// 所以奇數是常態不是邊角；而「先各自 floorDiv2 再相減」在鏡頭是奇數時
// 會對一半的格子差一列——只測偶數的話這個 bug 測不出來。
func TestCellOffsetFollowsOriginalWalk(t *testing.T) {
	// 鏡頭與走訪都用**原版的框**（含表頭那一列）；cellOffset 吃的是
	// 地形列號，所以傳進去之前要減掉 originalRowBias。
	for _, cam := range [][2]int{{36, 14}, {36, 11}, {35, 14}, {35, 11}, {0, 0}, {12, 41}} {
		v := &View{camWorldX: cam[0], camWorldY: cam[1]}
		for r := 0; r < 30; r++ {
			x, y := cam[0]-r, cam[1]+r
			for s := 0; s <= 30; s++ {
				dcol, drow := v.cellOffset(x, y-originalRowBias, 0)
				if dcol != s || drow != r {
					t.Fatalf("鏡頭 %v 的第 %d 列第 %d 格是 (%d,%d)，cellOffset 回 (%d,%d)",
						cam, r, s, x, y, dcol, drow)
				}
				if s%2 == 0 { // 走法：先 y+1，再 x+1
					y++
				} else {
					x++
				}
			}
		}
	}
}

// TestCellOffsetIsNotTwoProjections 明寫「不能各自投影再相減」這件事，
// 免得之後有人為了「看起來對稱」把 cellOffset 改回去。
func TestCellOffsetIsNotTwoProjections(t *testing.T) {
	v := &View{camWorldX: 36, camWorldY: 11} // (camY−camX) 是奇數
	camCol, camRow := isoProjectOriginal(v.camWorldX, v.camWorldY, 0)
	diff := 0
	for y := 0; y < 62; y++ { // y 是地形列號
		for x := 0; x < 64; x++ {
			_, drow := v.cellOffset(x, y, 0)
			_, row := isoProject(x, y, 0)
			if drow != row-camRow {
				diff++
			}
			if dcol, _ := v.cellOffset(x, y, 0); dcol != (x+y+originalRowBias)-camCol {
				t.Fatalf("(%d,%d) 的欄位移不該有差", x, y)
			}
		}
	}
	if diff == 0 {
		t.Error("鏡頭是奇數時兩種算法竟然一致——這個測試沒有鑑別力")
	}
}

// TestDisplayBandIsEightRows 釘住 `sub_1E0E1` 的 `mov cx, 8`：一帶是 **8 列**。
//
// 這是實際踩到的 bug：先前用 `isoRowPx`（16）當帶高，下一帶的內容會被
// 上一帶蓋掉一半，城壁的面上多出一排亮邊。四帶要**剛好**鋪滿 32 列——
// 這個算術檢查比「看起來對」可靠（docs/spec/58 §5）。
func TestDisplayBandIsEightRows(t *testing.T) {
	if displayBandRows != 8 {
		t.Fatalf("一帶 %d 列，預期 8（sub_1E0E1 的 mov cx, 8）", displayBandRows)
	}
	if got := displayBandRows * 4; got != battle.SubTileH {
		t.Fatalf("四帶共 %d 列，預期鋪滿 %d", got, battle.SubTileH)
	}
	// 四帶的來源列與目的列：0x30/0x20/0x10/0 → 源 24／16／8／0 → 目 0／8／16／24
	for i, tc := range [4][2]int{{24, 0}, {16, 8}, {8, 16}, {0, 24}} {
		if tc[0]+tc[1] != battle.SubTileH-displayBandRows {
			t.Fatalf("第 %d 帶 (源 %d, 目 %d) 不符 `di = 0x30 − dx` 的鏡射", i, tc[0], tc[1])
		}
	}
}

// TestDisplaySlotHeightAndStart 釘住 `sub_1DD22` 寫在表頭的兩個欄位：
// 高度看**任何**非零 unit、起始只看**小 unit**（< 0x20），兩者都取最大的 z，
// 而且單位是 2z（docs/spec/58 §2）。
func TestDisplaySlotHeightAndStart(t *testing.T) {
	var grid battleDisplayGrid
	set := func(z, raw int) {
		grid[3][5][z][0] = battleDisplaySlot{set: true,
			entry: battleDisplayEntry{raw: raw}}
	}
	set(0, 0x40) // 大
	set(1, 0x10) // 小
	set(2, 0)    // 空的不算
	set(3, 0x80) // 大
	info := makeDisplayInfo(&grid)
	if got := info[3][5].height; got != 2*3 {
		t.Fatalf("高度 = %d，預期 %d（最高的非零 unit 在 z=3）", got, 2*3)
	}
	if got := info[3][5].start; got != 2*1 {
		t.Fatalf("起始 = %d，預期 %d（唯一的小 unit 在 z=1）", got, 2*1)
	}
	// 一格都沒有 → 兩個欄位都是 0（原版是 `mov word ptr [si+1], 0`）
	if info[0][0].height != 0 || info[0][0].start != 0 {
		t.Fatalf("空格的表頭不是 0：%+v", info[0][0])
	}
}

// TestDisplayDepthRangeFollowsNeighbours 釘住範圍算式：從自己的起始深度畫到
// **自己與四個斜鄰格之中最高的**那一層（含頭含尾，docs/spec/58 §3）。
func TestDisplayDepthRangeFollowsNeighbours(t *testing.T) {
	var info [battleDisplayGridRows][battleDisplayGridCols]battleDisplaySlotInfo
	const row, anchor = 5, 9
	info[row][anchor] = battleDisplaySlotInfo{height: 2 * 1, start: 2 * 1}
	info[row][anchor+1] = battleDisplaySlotInfo{height: 2 * 4} // 最高的鄰居
	info[row+1][anchor-1] = battleDisplaySlotInfo{height: 2 * 2}
	z0, z1 := displayDepthRange(&info, row, anchor)
	if z0 != 1 || z1 != 4 {
		t.Fatalf("深度範圍 = %d..%d，預期 1..4", z0, z1)
	}
	// 自己最高時就只看自己
	info[row][anchor].height = 2 * 6
	if _, z1 := displayDepthRange(&info, row, anchor); z1 != 6 {
		t.Fatalf("自己最高時 z1 = %d，預期 6", z1)
	}
	// 沒有鄰居也不會越界
	info[row][anchor] = battleDisplaySlotInfo{}
	if a, b := displayDepthRange(&info, 0, 1); a != 0 || b != 0 {
		t.Fatalf("空格的範圍 = %d..%d，預期 0..0", a, b)
	}
}

// 點縮圖換算成鏡頭的公式（原版 `0x1C0C6`，docs/re/60 §7）。
// 座標是**原版的螢幕座標**，公式裡就寫死了，所以這裡不必知道版面。
func TestCameraFromMiniMapFormula(t *testing.T) {
	for _, tc := range []struct {
		name           string
		x, y           int
		worldX, worldY int
	}{
		{name: "left bottom", x: 0x1f0, y: 0xcf, worldX: 4, worldY: -18},
		{name: "right top", x: 0x26f, y: 0x50, worldX: 66, worldY: 44},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := &View{}
			v.SetCameraFromMiniMap(tc.x, tc.y)
			if v.camWorldX != tc.worldX || v.camWorldY != tc.worldY {
				t.Fatalf("縮圖相機=(%d,%d)，預期 (%d,%d)",
					v.camWorldX, v.camWorldY, tc.worldX, tc.worldY)
			}
			wantCol, wantRow := isoProjectOriginal(tc.worldX, tc.worldY, 0)
			if v.camCol != wantCol || v.camRow != wantRow {
				t.Fatalf("投影原點=(%d,%d)，預期 (%d,%d)",
					v.camCol, v.camRow, wantCol, wantRow)
			}
		})
	}
}
