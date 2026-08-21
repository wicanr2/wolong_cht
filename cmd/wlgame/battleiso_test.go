package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/isoview"
)

// 原生畫布要正好等於松崗 DOS/V 的可見 viewport。
//
// ⚠ 這一條是**跨層**的：畫布大小在 `internal/ui/isoview`，viewport 在
// 版面表。兩邊各改各的就會出現「畫得出來但裁切在錯的圖層」，
// 所以斷言留在看得到兩者的這一側。
func TestBattleViewBufferMatchesDOSVVisibleViewport(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	if got, want := isoview.NativeW*isoScale, l.Field.W; got != want {
		t.Fatalf("戰場 buffer 寬度 = %d，DOS/V viewport = %d", got, want)
	}
	if got, want := isoview.NativeH*isoScale, l.Field.H; got != want {
		t.Fatalf("戰場 buffer 高度 = %d，DOS/V viewport = %d", got, want)
	}
}

// 指令面板的熱區只有左邊 48 px；右邊 80 px 是陣線三格。
//
// ⚠ 算成整條 128 的話，點陣線的右半也會送出命令——**兩個熱區疊在一起**，
// 而畫面上看不出來（docs/spec/31 §2.1）。
func TestSideCommandCellsAre48Wide(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	cells := battleSideCommandCells(l.SideCommands)
	if len(cells) != 6 {
		t.Fatalf("命令格 %d 個，要 6", len(cells))
	}
	for i, c := range cells {
		if c.W != 48 || c.H != 16 {
			t.Errorf("第 %d 格 %d×%d，原版是 48×16", i, c.W, c.H)
		}
		if want := 280 + i*16; c.Y != want {
			t.Errorf("第 %d 格 y=%d，原版是 %d", i, c.Y, want)
		}
	}
	// 陣線三格不可以落在命令格裡。
	for _, line := range l.SideLines {
		for i, c := range cells {
			if line.X < c.right() && c.X < line.right() &&
				line.Y < c.bottom() && c.Y < line.bottom() {
				t.Errorf("陣線 %+v 與第 %d 個命令格 %+v 重疊", line, i, c)
			}
		}
	}
}
