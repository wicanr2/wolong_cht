package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 四張選單的粗格欄 ＝ 指令索引 × 3（指令格 48 px ＝ 3 個 16 px 粗格）。
// 值來自四個 handler 的 `dx` 立即值（docs/spec/126 §1）。
func TestCommandMenuAnchorsFollowCommandCell(t *testing.T) {
	for _, c := range []struct {
		name string
		menu *popupMenu
		dx   int // 原版的 dx 立即值
		x    int
	}{
		{"人事", personnelPopupMenu, 0x403, 48},
		{"軍團", corpsPopupMenu, 0x40C, 192},
		{"據點", cityPopupMenu, 0x40F, 240},
	} {
		if got := popupMenuCol(c.menu.cell); got != c.dx&0xFF {
			t.Errorf("%s 的欄 ＝ %d，原版 dx 低 byte 是 %d", c.name, got, c.dx&0xFF)
		}
		if c.menu.x != c.x || c.menu.y != 64 {
			t.Errorf("%s 的錨點 ＝ (%d, %d)，want (%d, 64)",
				c.name, c.menu.x, c.menu.y, c.x)
		}
		if got := c.dx >> 8; got != popupMenuRow {
			t.Errorf("%s 的列 ＝ %d，原版 dx 高 byte 是 %d", c.name, popupMenuRow, got)
		}
	}
	// 進言那一張還沒併進來，但欄號同一條規則（dx = 400h ⇒ 欄 0）。
	if got := popupMenuCol(naturalCommandAdvise); got != 0 {
		t.Errorf("進言的欄 ＝ %d，原版 dx 低 byte 是 0", got)
	}
}

// 三張的標籤都取自 `TALK.DAT`，而且**每一列的全形字數與原版相同**
// ——框寬就是靠它算的（docs/spec/125）。
func TestCommandMenuLabelsComeFromTalk(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	g := &game{lib: lib}
	for _, c := range []struct {
		name  string
		menu  *popupMenu
		rows  int
		cells int // 每一列的全形字數
	}{
		{"人事", personnelPopupMenu, 4, 7},
		{"軍團", corpsPopupMenu, 2, 6},
		{"據點", cityPopupMenu, 2, 6},
	} {
		labels := g.popupMenuLabels(c.menu)
		if len(labels) != c.rows {
			t.Errorf("%s 有 %d 列，want %d", c.name, len(labels), c.rows)
			continue
		}
		for i, l := range labels {
			if n := len([]rune(l)); n != c.cells {
				t.Errorf("%s 第 %d 列 %q 是 %d 個全形字，原版是 %d",
					c.name, i, l, n, c.cells)
			}
		}
		_, _, w, h := legacyChoiceRect(c.menu.x, c.menu.y, labels)
		if want := (c.cells + 1) * talkLinePitch; w != want {
			t.Errorf("%s 框寬 ＝ %d，want %d", c.name, w, want)
		}
		if want := (c.rows + 1) * talkLinePitch; h != want {
			t.Errorf("%s 框高 ＝ %d，want %d", c.name, h, want)
		}
		// fallback 也要同寬，否則沒有 TALK.DAT 時框會跳。
		for i, l := range c.menu.fallback {
			if n := len([]rune(l)); n != c.cells {
				t.Errorf("%s 的 fallback 第 %d 列 %q 是 %d 個全形字，want %d",
					c.name, i, l, n, c.cells)
			}
		}
	}
}

// 「首都確認」＝ 鏡頭移到首都（docs/spec/126 §1.1，與據點一覽共用尾段）。
func TestLocateCapitalMovesCamera(t *testing.T) {
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w.Player = 0
	capital := w.Factions[0].Capital
	if capital < 0 || capital >= len(w.Cities) {
		t.Skip("這個勢力沒有首都")
	}
	g := &game{world: w}
	g.beginLocateCapital()

	c := w.Cities[capital]
	want := &game{world: w, camX: c.X - centreCol, camY: c.Y - centreRow}
	want.clampCam()
	if g.camX != want.camX || g.camY != want.camY {
		t.Errorf("鏡頭 ＝ (%d, %d)，want (%d, %d)（首都在 (%d, %d)）",
			g.camX, g.camY, want.camX, want.camY, c.X, c.Y)
	}
}

// 三張選單開著時，指令列反白的是各自那一格（docs/spec/124）。
func TestPopupMenuHighlightsOwnCell(t *testing.T) {
	for _, m := range []*popupMenu{personnelPopupMenu, corpsPopupMenu, cityPopupMenu} {
		g := &game{}
		g.openPopupMenu(m)
		if got, want := g.activeCommandCell(), int(m.cell); got != want {
			t.Errorf("開著 TALK #%02X 那一張時亮第 %d 格，want %d", m.talk, got, want)
		}
		g.closePopupMenu()
		if got := g.activeCommandCell(); got != -1 {
			t.Errorf("關掉之後還亮著第 %d 格", got)
		}
	}
}
