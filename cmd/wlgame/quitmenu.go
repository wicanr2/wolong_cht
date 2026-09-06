package main

// 系統選單「遊戲結束」的兩項確認選單（docs/spec/153）。
//
// 原版是寫死位置的 `sub_193E9`：
//
//	mov ax, 102h        ; ah = 1、al = 2 項
//	mov cx, 51h         ; TALK #81「　終　　了　」「　取　　消　」
//	mov dx, 1116h       ; ★ dl = 16h → X = 352、dh = 11h → Y = 272
//	call sub_193E9
//	jb  取消            ; 右鍵
//	and al, al / jnz 取消 ; 選第 1 列也是取消
//	xor al, al / call sub_11CB1   ; 第 0 列 ＝ 終了

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/ui/talkmenu"
)

const (
	// quitMenuTalk 是選單字串的 TALK 索引（原版 `cx = 51h`）。
	quitMenuTalk = 0x51
	// 框的左上角是**寫死的**（原版 `dx = 1116h`），不跟著游標走——
	// 與行軍三選一、大地圖那兩張都不一樣。
	quitMenuX = 0x16 * 16
	quitMenuY = 0x11 * 16
)

// quitMenuState：**active 是 false 就是沒開**——零值安全。
type quitMenuState struct {
	active bool
	row    int
}

func (g *game) openQuitMenu() { g.quitMenu = quitMenuState{active: true} }

// quitMenuRows 取 TALK #81 的兩行。
//
// ⛔ 走 `talkmenu.MenuLabels` 不是 `talkLines`：後者會 `TrimRight` 掉行尾的
// 全形空白，而**框寬由第一列的全形字數決定**（docs/spec/125）。
func (g *game) quitMenuRows() []string {
	fallback := []string{"　終　　了　", "　取　　消　"}
	if g == nil || g.lib == nil {
		return fallback
	}
	return talkmenu.MenuLabels(g.lib.Talk, quitMenuTalk, nil, fallback)
}

// quitMenuLeaves：原版是 `and al, al / jnz 取消`——**只有第 0 列會結束**，
// 其他任何列（含右鍵回來的 CF）都當成取消。
func quitMenuLeaves(row int) bool { return row == 0 }

// updateQuitMenu 回 true 表示要結束遊戲。
func (g *game) updateQuitMenu() bool {
	m := &g.quitMenu
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) ||
		pressed(ebiten.KeyEscape) {
		m.active = false
		return false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if row, ok := g.talkChoiceClick(quitMenuX, quitMenuY, g.quitMenuRows()); ok {
			m.active = false
			return quitMenuLeaves(row)
		}
		return false
	}
	switch {
	case pressed(ebiten.KeyArrowUp), pressed(ebiten.KeyArrowDown):
		m.row = 1 - m.row
	case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
		m.active = false
		return quitMenuLeaves(m.row)
	case pressed(ebiten.Key1):
		m.active = false
		return true
	case pressed(ebiten.Key2):
		m.active = false
	}
	return false
}

func (g *game) drawQuitMenu(screen *ebiten.Image) {
	if !g.quitMenu.active {
		return
	}
	g.drawLegacyChoiceBox(screen, quitMenuX, quitMenuY,
		g.quitMenuRows(), g.quitMenu.row)
}
