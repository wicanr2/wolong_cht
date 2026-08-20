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
