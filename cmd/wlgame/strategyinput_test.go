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
		if r.Min.X <= chrome.Tile || r.Max.X > strategyMapW-chrome.Tile {
			t.Fatalf("第 %d 格越過命令列內容外框：%v", i, r)
		}
	}

	// 外框、上下留白、格間與地圖都不能命中。
	misses := []image.Point{
		{X: 0, Y: strategyCommandY + strategyCommandH/2},
		{X: strategyMapW - 1, Y: strategyCommandY + strategyCommandH/2},
		{X: strategyCommandCellRect(0).Min.X - 1, Y: strategyCommandHitY + 4},
		{X: strategyCommandCellRect(0).Max.X, Y: strategyCommandHitY + 4},
		{X: strategyCommandCellRect(0).Min.X + 1, Y: strategyCommandHitY - 1},
		{X: strategyCommandCellRect(0).Min.X + 1, Y: strategyCommandHitY + strategyCommandHitH},
		{X: strategyCommandCellRect(0).Max.X + 2, Y: strategyCommandHitY + 4},
		{X: strategyMapW / 2, Y: strategyMapY},
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
