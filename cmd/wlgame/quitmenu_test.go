package main

import "testing"

// 系統選單第 6 列開的是兩項選單，不是 ＹＥＳ／ＮＯ 對話框（docs/spec/153 §1）。
func TestQuitMenuOpensFromSystemRow(t *testing.T) {
	g := &game{}
	g.dispatchSystemRow(sysRowQuit, true)
	if !g.quitMenu.active {
		t.Fatal("點「遊戲結束」沒有開兩項選單")
	}
	if g.quitting {
		t.Error("不該同時掀起 F10 那條的 ＹＥＳ／ＮＯ 對話框")
	}
	if g.quitMenu.row != 0 {
		t.Errorf("剛開的選取列 ＝ %d，want 0", g.quitMenu.row)
	}
}

// 位置是**寫死**的 `dx = 1116h`，不跟著游標（docs/spec/153 §1.1）。
func TestQuitMenuAnchorIsFixed(t *testing.T) {
	if quitMenuX != 0x16*16 || quitMenuY != 0x11*16 {
		t.Fatalf("框的左上角 ＝ (%d,%d)，want (352,272)", quitMenuX, quitMenuY)
	}
}

// TALK #81 兩列，且 fallback 也要補位——框寬由第一列的全形字數決定。
func TestQuitMenuRowsComeFromTalk(t *testing.T) {
	g := newTalkTestGame(t)
	rows := g.quitMenuRows()
	if len(rows) != 2 {
		t.Fatalf("列數 ＝ %d，want 2", len(rows))
	}
	for i, want := range []string{"　終　　了　", "　取　　消　"} {
		if rows[i] != want {
			t.Errorf("第 %d 列 ＝ %q，want %q", i, rows[i], want)
		}
	}
	// 沒有素材時的 fallback 走同一個形狀。
	empty := &game{}
	if got := empty.quitMenuRows(); len(got) != 2 || got[0] != "　終　　了　" {
		t.Errorf("fallback ＝ %q", got)
	}
}

// 三條出路裡只有第 0 列會離開（docs/spec/153 §1）。
func TestQuitMenuRowSemantics(t *testing.T) {
	if !quitMenuLeaves(0) {
		t.Error("第 0 列「終了」應該要離開")
	}
	for _, row := range []int{1, 2, -1} {
		if quitMenuLeaves(row) {
			t.Errorf("第 %d 列不該離開", row)
		}
	}
}

// 框是六個全形字兩列 ⇒ 112×48，右下角落在 (464, 320) 都還在畫面內
// （docs/spec/153 §1.2）。
func TestQuitMenuBoxRect(t *testing.T) {
	g := newTalkTestGame(t)
	bx, by, w, h := legacyChoiceRect(quitMenuX, quitMenuY, g.quitMenuRows())
	if bx != 352 || by != 272 || w != 112 || h != 48 {
		t.Fatalf("框 ＝ (%d,%d,%d,%d)，want (352,272,112,48)", bx, by, w, h)
	}
	if bx+w > screenW || by+h > screenH {
		t.Errorf("框超出畫面：右緣 %d、下緣 %d", bx+w, by+h)
	}
}
