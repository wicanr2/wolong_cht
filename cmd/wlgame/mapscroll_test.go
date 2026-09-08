package main

import "testing"

func TestPointerCaptureRebaseWaitsForFreshFrame(t *testing.T) {
	saved := desktopPointer
	t.Cleanup(func() { desktopPointer = saved })
	desktopPointer.x, desktopPointer.y = 50, 50
	desktopPointer.rawX, desktopPointer.rawY = 50, 50
	desktopPointer.rebase = 2
	g := &game{camX: 186, camY: 102}
	// 低更新率會在同一張輸入快照補跑多次 Update。
	for range 8 {
		dx, dy := consumeDesktopPointer(50, 50)
		g.scrollMapPixels(dx, dy)
	}
	finishDesktopPointerFrame()
	// rapid-focus-diagnostic 的原始重定位序列，並非玩家移動。
	dx, dy := consumeDesktopPointer(-220, -100)
	g.scrollMapPixels(dx, dy)
	if g.camX != 186 || g.camY != 102 || g.camSubX != 0 || g.camSubY != 0 {
		t.Fatal("重新捕捉的座標跳變被當成外推")
	}
	if desktopPointer.x != 50 || desktopPointer.y != 50 {
		t.Fatal("重新捕捉改變畫面游標位置")
	}
	// 新基準之後的大幅真實移動不可被距離閾值吞掉。
	dx, dy = consumeDesktopPointer(420, 300)
	g.scrollMapPixels(dx, dy)
	if dx != 51 || dy != 51 {
		t.Fatalf("真實外推應保留 51,51，得到 %d,%d", dx, dy)
	}
	dx, dy = consumeDesktopPointer(420, 300)
	if dx != 0 || dy != 0 {
		t.Fatal("停手後仍捲動")
	}
}

func TestMapScrollOriginalPushStopAndReverse(t *testing.T) {
	g := &game{camX: 186, camY: 102}
	x, over := pushPointer(320, 448, 640)
	g.scrollMapPixels(over, 0)
	if x != 639 || g.camX*16+g.camSubX != 3105 {
		t.Fatalf("原版原點應為 3105，得到 %d；游標 %d", g.camX*16+g.camSubX, x)
	}
	for i := 0; i < 100; i++ {
		x, over = pushPointer(x, 0, 640)
		g.scrollMapPixels(over, 0)
	}
	if g.camX*16+g.camSubX != 3105 {
		t.Fatal("停手後仍捲動")
	}
	x, over = pushPointer(x, -16, 640)
	if x != 623 || over != 0 {
		t.Fatal("反向移動應先離開邊緣")
	}
	g.scrollMapPixels(-100000, 100000)
	if g.camX != 0 || g.camSubX != 0 || g.camY != 256-viewRows || g.camSubY != 0 {
		t.Fatal("地圖邊界沒有夾住")
	}
}

func TestMapScrollPickRetainsOriginalSeparateTruncation(t *testing.T) {
	g := &game{camX: 185, camY: 84, camSubX: 1, camSubY: 1}
	x, y, ok := g.mapPickTileAt(639, 399)
	if !ok || x != 224 || y != 106 {
		t.Fatalf("re/85 原版案例：得到 %d,%d", x, y)
	}
}
