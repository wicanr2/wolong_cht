package main

import (
	"testing"
)

// 縮圖的命中矩形要對齊原版的 0x1f0..0x26f／0x50..0xcf。
//
// ⚠ 這一條是**跨層**的：公式在 `internal/ui/isoview`（它直接寫死原版的
// 螢幕座標），矩形在版面表。兩邊對不上的話點縮圖會跳到別的地方，
// 而畫面上看起來只是「點不準」。
func TestBattleMiniMapRectMatchesOriginal(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	if l.SideMiniMap.X != 0x1f0 || l.SideMiniMap.Y != 0x50 ||
		l.SideMiniMap.right()-1 != 0x26f || l.SideMiniMap.bottom()-1 != 0xcf {
		t.Fatalf("縮圖命中矩形=%+v，未對齊原版 0x1f0..0x26f／0x50..0xcf", l.SideMiniMap)
	}
}
