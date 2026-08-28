package main

import "testing"

// `▶▶` 是切換（docs/spec/102）：熱區就是側欄的 SideFooter，開→關→開，
// 按過之後才描框。
func TestBattleFastForwardToggle(t *testing.T) {
	l := dosvBattleLayoutFor(screenW, screenH)
	if !l.SideFooter.containsPoint(l.SideFooter.X+64, l.SideFooter.Y+8) {
		t.Fatal("SideFooter 熱區不含自己的中心")
	}
	if l.SideFooter.X != 496 || l.SideFooter.Y != 376 || l.SideFooter.W != 128 || l.SideFooter.H != 16 {
		t.Fatalf("SideFooter = %+v，預期 (496,376) 128×16（熱區 0x0F）", l.SideFooter)
	}
	g := &game{}
	if g.battleFastForward || g.battleFFTouched {
		t.Fatal("初值應該是關、沒按過")
	}
	g.toggleBattleFastForward()
	if !g.battleFastForward || !g.battleFFTouched {
		t.Fatal("第一下沒有打開")
	}
	g.toggleBattleFastForward()
	if g.battleFastForward || !g.battleFFTouched {
		t.Fatal("第二下沒有關掉，或忘了按過")
	}
}
