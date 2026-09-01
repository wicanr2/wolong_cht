package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 指令列「軍團」是兩列選單，不是直接開一覽（docs/spec/110）。
// 兩個錨點來自 `sub_1628F` 的 `dx = 40Ch` ＝ 粗格 (12, 4)。
func TestCorpsCommandMenuRows(t *testing.T) {
	if len(corpsMenuLabels) != 2 {
		t.Fatalf("選單 %d 列，原版是 2 列（TALK #79）", len(corpsMenuLabels))
	}
	for i, want := range [...]string{"位置確認", "行軍指示"} {
		if corpsMenuLabels[i] != want {
			t.Errorf("第 %d 列 ＝ %q，want %q", i, corpsMenuLabels[i], want)
		}
	}
	if corpsMenuX != 192 || corpsMenuY != 64 {
		t.Errorf("選單錨點 ＝ (%d, %d)，want (192, 64)＝粗格 (12, 4)",
			corpsMenuX, corpsMenuY)
	}

	// ⭐ 用**實機量到的**列位置釘死幾何：原版那張選單的第一列
	// （位置確認）佔遊戲座標 y 72–87（`parity-tap5/menu.png`，
	// 螢幕 y 112–127 減掉 40 px 黑邊，docs/spec/110）。
	bx, by, w, h := legacyChoiceRect(corpsMenuX, corpsMenuY, corpsMenuLabels[:])
	if bx != 192 || by != 64 {
		t.Errorf("選單框左上角 ＝ (%d, %d)，want (192, 64)", bx, by)
	}
	if firstRow := by + talkLinePitch/2; firstRow != 72 {
		t.Errorf("第一列 y ＝ %d，實機量到 72", firstRow)
	}
	if h != (len(corpsMenuLabels)+1)*talkLinePitch {
		t.Errorf("框高 ＝ %d，兩列應該是 %d", h, (len(corpsMenuLabels)+1)*talkLinePitch)
	}
	if w <= 0 {
		t.Errorf("框寬 ＝ %d", w)
	}
}

// 選單開著的時候要吃掉輸入，選完就收掉。
func TestCorpsCommandMenuLifecycle(t *testing.T) {
	g := &game{}
	if g.corpsMenuActive() {
		t.Fatal("還沒開就是 active")
	}
	g.openCorpsCommandMenu()
	if !g.corpsMenuActive() || g.corpsMenu.row != 0 {
		t.Fatalf("開起來之後 active=%v row=%d", g.corpsMenuActive(), g.corpsMenu.row)
	}
	g.closeCorpsCommandMenu()
	if g.corpsMenuActive() {
		t.Error("關掉之後還是 active")
	}
}

// 位置確認的鏡頭 ＝ (軍團 X − 20, Y − 12)，與 `sub_12151(ax=14h, cx=0Ch)`
// 的立即值同一組（docs/spec/110 §1.1）。
func TestLocateCorpsMovesCamera(t *testing.T) {
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w.Player = 0
	f := &w.Factions[w.Player]
	f.Reserves = [3]int{600, 600, 600}

	// 找一個能帶兵的武將編一支軍團出來。
	leader := -1
	for i, gen := range w.Generals {
		if gen.Alive && gen.Faction == w.Player && !gen.Posted && gen.Captor == 0xFF {
			leader = i
			break
		}
	}
	if leader < 0 {
		t.Skip("這個勢力沒有可以帶兵的武將")
	}
	manned := [army.Positions]bool{}
	manned[0] = true
	if err := w.FormCorps(leader, [army.Positions]army.TroopType{}, manned); err != nil {
		t.Fatalf("編成失敗：%v", err)
	}

	g := &game{world: w, camX: 0, camY: 0}
	c := w.Corps[leader]
	g.focusCorps(leader)

	wantX, wantY := c.X-centreCol, c.Y-centreRow
	// clampCam 會把超出邊界的夾回來，所以拿夾過的期望值比。
	want := &game{world: w, camX: wantX, camY: wantY}
	want.clampCam()
	if g.camX != want.camX || g.camY != want.camY {
		t.Errorf("鏡頭 ＝ (%d, %d)，want (%d, %d)（軍團在 (%d, %d)）",
			g.camX, g.camY, want.camX, want.camY, c.X, c.Y)
	}
}

// 越界的軍團編號不能動到鏡頭。
func TestLocateCorpsIgnoresOutOfRange(t *testing.T) {
	g := &game{world: &state.World{}, camX: 7, camY: 9}
	for _, n := range []int{-1, 1 << 20} {
		g.focusCorps(n)
		if g.camX != 7 || g.camY != 9 {
			t.Fatalf("軍團編號 %d 動到了鏡頭：(%d, %d)", n, g.camX, g.camY)
		}
	}
}
