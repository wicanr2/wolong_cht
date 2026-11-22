package isoview

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

func TestRetainedDrawingKeepsFirstUntilAnotherUnitErases(t *testing.T) {
	v := retainedTestView(45, 28)
	v.drawRetained(20, 40, 2, 228, false)
	v.drawRetained(20, 40, 2, 408, false)
	if got := v.retained[12][15][2][1].raw; got != 229 {
		t.Fatalf("已佔用時應保留第一筆，得到%d", got)
	}
	v.EraseSoldierDrawing(20, 40, 2)
	v.drawRetained(20, 40, 2, 408, false)
	if got := v.retained[12][15][2][1].raw; got != 409 {
		t.Fatalf("另一人擦除後應可登錄，得到%d", got)
	}
}

func TestRetainedDrawingClipsAnchorBeforeUpperHalf(t *testing.T) {
	v := retainedTestView(45, 16)
	v.retained[23][15][3][1] = battleDisplaySlot{set: true, raw: 999}
	v.drawRetained(20, 40, 2, 228, false)
	v.EraseSoldierDrawing(20, 40, 2)
	if got := v.retained[23][15][3][1].raw; got != 999 {
		t.Fatalf("錨點列24越界，不能改列23的上半，得到%d", got)
	}
}

func TestRetainedDrawingRejectsOnlyBuriedHalf(t *testing.T) {
	v := retainedTestView(45, 28)
	v.retainedInfo[12][15].start = 6
	v.drawRetained(20, 40, 2, 228, false)
	if !displaySlotEmpty(v.retained[12][15][2][1]) || v.retained[11][15][3][1].raw != 228 {
		t.Fatal("下半低於起始深度時，仍可單獨登錄上半")
	}
}

func TestCameraRebuildClearsWholeRetainedGrid(t *testing.T) {
	b := &tactical.Battle{}
	v := retainedTestView(0, 0)
	v.observed, v.retainedCamera = b, [2]int{9, 9}
	v.retained[1][2][3][1] = battleDisplaySlot{set: true, raw: 228}
	v.retained[20][25][1][1] = battleDisplaySlot{set: true, raw: 408}
	v.BeginDrawingFrame(b)
	if !displaySlotEmpty(v.retained[1][2][3][1]) || !displaySlotEmpty(v.retained[20][25][1][1]) {
		t.Fatal("換鏡頭的地形重建應清掉整張人物格表")
	}
}

func TestCompositeKeyTracksEveryNativeContribution(t *testing.T) {
	grid := newDisplayGrid[battleDisplayCell](battleDisplayGridRows, battleDisplayGridCols)
	base := displayCompositeKey(grid, 10, 15, 0, 6)
	for _, at := range [5][2]int{{10, 14}, {10, 16}, {10, 15}, {11, 14}, {11, 16}} {
		for z := range grid[at[0]][at[1]] {
			for lane := 0; lane < 2; lane++ {
				grid[at[0]][at[1]][z][lane] = battleDisplaySlot{set: true, raw: 228}
				if displayCompositeKey(grid, 10, 15, 0, 6) == base {
					t.Fatalf("遺漏原版來源%v，深度%d，lane%d", at, z, lane)
				}
				grid[at[0]][at[1]][z][lane] = battleDisplaySlot{}
			}
		}
	}
	if displayCompositeKey(grid, 10, 15, 1, 6) == base || displayCompositeKey(grid, 10, 15, 0, 5) == base {
		t.Fatal("起訖深度必須參與合成鍵值")
	}
}

func retainedTestView(col, row int) *View {
	return &View{camCol: col, camRow: row,
		retained:     newDisplayGrid[battleDisplayCell](battleDisplayGridRows, battleDisplayGridCols),
		retainedInfo: newDisplayGrid[battleDisplaySlotInfo](battleDisplayGridRows, battleDisplayGridCols),
	}
}
