package main

import "testing"

func TestDesktopSaveClickIgnoresEmptyReadSlot(t *testing.T) {
	g := &game{saveUI: saveUIState{
		active: true, action: saveRead, slot: 0,
		slots: []launcherSlot{{Available: true}, {}, {}, {}},
	}}
	for _, slot := range []int{-1, 1, 2, 3, 4} {
		g.clickSaveSlot(slot)
		if !g.saveUI.active || g.saveUI.slot != 0 || g.saveUI.touched || g.lastEvent != "" {
			t.Fatalf("空槽或無效槽 %d 不應執行讀取或改選：%#v，%s", slot, g.saveUI, g.lastEvent)
		}
	}
}
