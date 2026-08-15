package main

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

func TestNaturalCommandHitTestUsesOnlyInnerCells(t *testing.T) {
	for i := range naturalCommandLabels {
		r := strategyCommandCellRect(i)
		if got, ok := hitTestNaturalCommand(r.Min.X+r.Dx()/2, r.Min.Y+r.Dy()/2); !ok || got != naturalCommandID(i) {
			t.Fatalf("第 %d 格中心命中 = (%d, %v)，want (%d, true)", i, got, ok, i)
		}
		if r.Min.X <= chrome.Tile || r.Max.X > strategyCommandW-chrome.Tile {
			t.Fatalf("第 %d 格越過命令列內容外框：%v", i, r)
		}
	}

	// 原版的八格之間**沒有間隙**（索引 ＝ (x−24)÷48），所以這裡驗的是
	// 左右兩端的邊界與上下留白，不是格間死區。
	misses := []image.Point{
		{X: 0, Y: strategyCommandY + strategyCommandH/2},
		{X: strategyCommandW - 1, Y: strategyCommandY + strategyCommandH/2},
		{X: strategyCommandCellRect(0).Min.X - 1, Y: strategyCommandHitY + 4},
		{X: strategyCommandCellRect(7).Max.X, Y: strategyCommandHitY + 4},
		{X: strategyCommandCellRect(0).Min.X + 1, Y: strategyCommandHitY - 1},
		{X: strategyCommandCellRect(0).Min.X + 1, Y: strategyCommandHitY + strategyCommandHitH},
		{X: strategyCommandW / 2, Y: strategyMapY + strategyCommandH},
		{X: strategySidebarX, Y: strategyCommandY + strategyCommandH/2},
	}
	for _, p := range misses {
		if _, ok := hitTestNaturalCommand(p.X, p.Y); ok {
			t.Fatalf("非命令內容點 %v 意外命中", p)
		}
	}
}

func TestNaturalCommandShortcutSharesLabelMapping(t *testing.T) {
	cases := []struct {
		key ebiten.Key
		id  naturalCommandID
	}{
		{ebiten.KeyP, naturalCommandAdvise},
		{ebiten.KeyJ, naturalCommandPersonnel},
		{ebiten.KeyF, naturalCommandFinance},
		{ebiten.KeyA, naturalCommandFormation},
		{ebiten.KeyC, naturalCommandCorps},
		{ebiten.KeyT, naturalCommandCity},
		{ebiten.KeyG, naturalCommandGeneral},
		{ebiten.KeyK, naturalCommandFaction},
	}
	for i, tc := range cases {
		if naturalCommandID(i) != tc.id {
			t.Fatalf("畫面第 %d 格 ID = %d，want %d", i, i, tc.id)
		}
		got, ok := strategyCommandForShortcut(tc.key)
		if !ok || got != tc.id {
			t.Errorf("快捷鍵 %v 映射 = (%d, %v)，want (%d, true)", tc.key, got, ok, tc.id)
		}
	}
	if _, ok := strategyCommandForShortcut(ebiten.KeyM); ok {
		t.Fatal("M（行軍）不應被映射到頂端八格")
	}
}

func TestNaturalCommandActionTableMatchesEightLabels(t *testing.T) {
	if len(naturalCommandLabels) != 8 || len(naturalCommandActions) != len(naturalCommandLabels) {
		t.Fatalf("畫面標籤／共享 action 數量 = %d/%d，want 8/8", len(naturalCommandLabels), len(naturalCommandActions))
	}
}

// 橫幅右側五格開關的幾何與語意（docs/spec/13、docs/re/47 §2）。
func TestHUDSwitchGeometryAndSemantics(t *testing.T) {
	// 五格：336/368/400/432/464，各 32×32，Y 落在 0–32 的橫幅內。
	for i := 0; i < hudSwitchN; i++ {
		x := hudSwitchX0 + i*hudSwitchW
		for _, p := range []image.Point{{X: x, Y: 0}, {X: x + hudSwitchW - 1, Y: bannerH - 1}} {
			got, ok := hitTestHUDSwitch(p.X, p.Y)
			if !ok || got != i {
				t.Fatalf("%v 命中 = (%d, %v)，want (%d, true)", p, got, ok, i)
			}
		}
	}
	// 橫幅左段（熱區 6）與橫幅以下都不是開關。
	for _, p := range []image.Point{
		{X: hudSwitchX0 - 1, Y: 16}, {X: 0, Y: 0},
		{X: hudSwitchX0 + hudSwitchN*hudSwitchW, Y: 16},
		{X: hudSwitchX0, Y: bannerH},
	} {
		if _, ok := hitTestHUDSwitch(p.X, p.Y); ok {
			t.Fatalf("%v 不該是開關", p)
		}
	}
	// 第五格原版接 nullsub_1，這裡也不接東西。
	if w := hudSwitchWindow(4); w != 0 {
		t.Fatalf("第五格對到視窗 %d，原版是 nullsub_1", w)
	}

	// 左鍵開、右鍵關；重複同一個動作不翻轉。
	g := &game{hud: hudCommand}
	if !g.hudOpen(hudCommand) || g.hudOpen(hudMinimap) {
		t.Fatal("初始狀態不對")
	}
	g.hudSet(hudCommand, false)
	g.hudSet(hudCommand, false)
	if g.hudOpen(hudCommand) {
		t.Fatal("右鍵關了兩次反而變成開著——關／開被寫成 toggle")
	}
	g.hudSet(hudMinimap, true)
	g.hudSet(hudMinimap, true)
	if !g.hudOpen(hudMinimap) {
		t.Fatal("左鍵開了兩次反而變成關著")
	}
	// 系統視窗與另外三個一樣，狀態就在 g.hud 裡（docs/spec/13 §2.5）。
	g.hudSet(hudSystem, true)
	if !g.hudOpen(hudSystem) {
		t.Fatal("點第四格沒有開系統視窗")
	}
	g.hudSet(hudSystem, false)
	if g.hudOpen(hudSystem) {
		t.Fatal("右鍵沒有關掉系統視窗")
	}
}
