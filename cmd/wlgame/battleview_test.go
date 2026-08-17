package main

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
			if got := projectileSourceIndex(tc.p); got != tc.want {
				t.Fatalf("raw 圖號 = %#x，預期 %#x", got, tc.want)
			}
		})
	}
}

func TestBattleViewBufferMatchesDOSVVisibleViewport(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	if got, want := isoNativeW*isoScale, l.Field.W; got != want {
		t.Fatalf("戰場 buffer 寬度 = %d，DOS/V viewport = %d", got, want)
	}
	if got, want := isoNativeH*isoScale, l.Field.H; got != want {
		t.Fatalf("戰場 buffer 高度 = %d，DOS/V viewport = %d", got, want)
	}
	// 原版投影仍檢查 31×24 欄列；最後一欄／列由 240×184 viewport 裁切。
	if isoCols*isoColPx <= isoNativeW || isoRows*isoRowPx <= isoNativeH {
		t.Fatal("投影測試範圍必須大於可見 viewport，才能保留原版邊界裁切")
	}
}

func TestBattleCameraStartsAtOriginalWorldOrigin(t *testing.T) {
	v := &battleView{camWorldX: 0x24, camWorldY: 0x0e}
	v.applyCameraOrigin()
	if v.camCol != 0x32 || v.camRow != 0x15 {
		t.Fatalf("sub_1DC9D 投影原點 = (%#x,%#x)，預期 (0x32,0x15)",
			v.camCol, v.camRow)
	}

	// 相機不再從兵座標推導。改變任意兵的位置，不應改寫保存的 world origin。
	v.applyCameraOrigin()
	if v.camWorldX != 0x24 || v.camWorldY != 0x0e {
		t.Fatalf("相機 world origin 被繪圖改寫：(%#x,%#x)", v.camWorldX, v.camWorldY)
	}
}

func TestBattleCameraMiniMapFormula(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	if l.SideMiniMap.X != 0x1f0 || l.SideMiniMap.Y != 0x50 ||
		l.SideMiniMap.right()-1 != 0x26f || l.SideMiniMap.bottom()-1 != 0xcf {
		t.Fatalf("縮圖命中矩形=%+v，未對齊原版 0x1f0..0x26f／0x50..0xcf", l.SideMiniMap)
	}
	for _, tc := range []struct {
		name           string
		x, y           int
		worldX, worldY int
	}{
		{name: "left bottom", x: 0x1f0, y: 0xcf, worldX: 4, worldY: -18},
		{name: "right top", x: 0x26f, y: 0x50, worldX: 66, worldY: 44},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := &battleView{}
			v.setCameraFromMiniMap(tc.x, tc.y)
			if v.camWorldX != tc.worldX || v.camWorldY != tc.worldY {
				t.Fatalf("縮圖相機=(%d,%d)，預期 (%d,%d)",
					v.camWorldX, v.camWorldY, tc.worldX, tc.worldY)
			}
			wantCol, wantRow := isoProject(tc.worldX, tc.worldY, 0)
			if v.camCol != wantCol || v.camRow != wantRow {
				t.Fatalf("投影原點=(%d,%d)，預期 (%d,%d)",
					v.camCol, v.camRow, wantCol, wantRow)
			}
		})
	}
}

func TestBattleDisplayListSplitsSoldiersIntoOriginalRawUnits(t *testing.T) {
	v := &battleView{camCol: 0, camRow: 0}
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
	for _, cam := range [][2]int{{36, 14}, {36, 11}, {35, 14}, {35, 11}, {0, 0}, {12, 41}} {
		v := &battleView{camWorldX: cam[0], camWorldY: cam[1]}
		for r := 0; r < 30; r++ {
			x, y := cam[0]-r, cam[1]+r
			for s := 0; s <= 30; s++ {
				dcol, drow := v.cellOffset(x, y, 0)
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
	v := &battleView{camWorldX: 36, camWorldY: 11} // (camY−camX) 是奇數
	camCol, camRow := isoProject(v.camWorldX, v.camWorldY, 0)
	diff := 0
	for y := 0; y < 62; y++ {
		for x := 0; x < 64; x++ {
			_, drow := v.cellOffset(x, y, 0)
			_, row := isoProject(x, y, 0)
			if drow != row-camRow {
				diff++
			}
			if dcol, _ := v.cellOffset(x, y, 0); dcol != (x+y)-camCol {
				t.Fatalf("(%d,%d) 的欄位移不該有差", x, y)
			}
		}
	}
	if diff == 0 {
		t.Error("鏡頭是奇數時兩種算法竟然一致——這個測試沒有鑑別力")
	}
}
