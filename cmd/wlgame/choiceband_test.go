package main

import "testing"

// 期望值取自 sub_10414／1054D／1059B 的原始門檻與回置值。
func TestChoiceBandOriginalBoundaries(t *testing.T) {
	for _, tc := range []struct {
		count, row, y, next, pointerY int
		warp                          bool
	}{
		{2, 0, 100, 0, 100, false}, {2, 0, 131, 0, 131, false},
		{2, 0, 132, 1, 100, true}, {2, 0, 99, 0, 131, true},
		{2, 1, 100, 1, 100, false}, {2, 1, 131, 1, 131, false},
		{2, 1, 132, 1, 100, true}, {2, 1, 99, 0, 131, true},
		{4, 1, 399, 2, 100, true}, {4, 3, 399, 3, 100, true},
		{4, 2, 0, 1, 131, true}, {4, 0, 0, 0, 131, true},
		{5, 2, 132, 3, 100, true}, {5, 4, 132, 4, 100, true},
	} {
		next, y, warp := choiceBandMove(tc.row, tc.count, tc.y)
		if next != tc.next || y != tc.pointerY || warp != tc.warp {
			t.Fatalf("count=%d row=%d y=%d got (%d,%d,%t), want (%d,%d,%t)", tc.count, tc.row, tc.y, next, y, warp, tc.next, tc.pointerY, tc.warp)
		}
	}
}

func choiceTestPointer(t *testing.T) {
	t.Helper()
	old, offset := desktopPointer, desktopGameOffset
	t.Cleanup(func() { desktopPointer, desktopGameOffset = old, offset })
	desktopGameOffset = desktopBarHeight
	desktopPointer.active, desktopPointer.rebase = true, 0
	desktopPointer.x, desktopPointer.y = 240, 47
}

func TestChoiceBandLifecycleAndToolbarOwnership(t *testing.T) {
	choiceTestPointer(t)
	g := &game{}
	g.openPopupMenu(corpsPopupMenu)
	g.syncChoiceBand()
	if g.choiceBand.owner == nil || desktopPointer.x != 0 || desktopPointer.y != 100 {
		t.Fatal("selector did not save and reposition pointer")
	}
	for a := desktopLanguage; a <= desktopQuit; a++ {
		if g.desktopEnabled(a) {
			t.Fatalf("toolbar action %d can steal original selector input", a)
		}
	}
	desktopPointer.x, desktopPointer.y = 500, 132 // 遠在選單文字之外。
	g.talkChoiceInput(&g.cmdMenu.row, 2)
	if g.cmdMenu.row != 1 || desktopPointer.x != 0 || desktopPointer.y != 100 {
		t.Fatal("downward motion did not select second row and reset mouse")
	}
	for i := 0; i < 40; i++ {
		g.talkChoiceInput(&g.cmdMenu.row, 2)
	}
	if g.cmdMenu.row != 1 {
		t.Fatal("stationary mouse repeated selection")
	}
	g.pointerAction = desktopConfirm
	desktopPointer.y = 399
	row, ok := g.talkChoiceInput(&g.cmdMenu.row, 2)
	if !ok || row != 1 {
		t.Fatal("confirmation did not consume current row before movement")
	}
	if desktopPointer.x != 240 || desktopPointer.y != 47 || g.choiceBand.owner != nil {
		t.Fatal("selection did not restore original pointer")
	}
}

func TestChoiceBandCancelRestoresOnlyCurrentOwner(t *testing.T) {
	choiceTestPointer(t)
	g := &game{}
	g.openPopupMenu(corpsPopupMenu)
	g.syncChoiceBand()
	g.cancelFn = func() bool { return true }
	g.updatePopupMenu()
	if g.popupMenuActive() || g.choiceBand.owner != nil || desktopPointer.x != 240 || desktopPointer.y != 47 {
		t.Fatal("right cancel did not close selector and restore cursor")
	}
}

func TestChoiceBandReleasesForChildAndReinitializes(t *testing.T) {
	choiceTestPointer(t)
	g := &game{advise: advisePickCommand}
	g.syncChoiceBand()
	g.adviseCmdRow = 3
	g.advise = advisePickCapital
	g.syncChoiceBand()
	if g.choiceBand.owner != nil || desktopPointer.x != 240 || desktopPointer.y != 47 {
		t.Fatal("child input retained selector mouse ownership")
	}
	g.advise = advisePickCommand
	g.syncChoiceBand()
	if g.choiceBand.owner == nil || g.adviseCmdRow != 0 || desktopPointer.y != 100 {
		t.Fatal("return did not start a fresh original selector")
	}
}
