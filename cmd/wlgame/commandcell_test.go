package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
)

// 指令列的反白涵蓋八格，不是只有彈出選單那三格（docs/spec/124 §3.5）。
//
// 原版 `sub_161CA` 在 `call cs:funcs_161FE[bx]` 的**前後各 XOR 一次**，
// 而那一段程式碼不分索引——所以只要那一格的 handler 還沒回來就一直亮著。
func TestActiveCommandCellCoversEveryFlow(t *testing.T) {
	for _, tc := range []struct {
		name string
		cell naturalCommandID
		open func(*game)
	}{
		{"進言", naturalCommandAdvise, func(g *game) { g.advise = advisePickCommand }},
		{"財政", naturalCommandFinance, func(g *game) { g.finance.active = true }},
		{"編成", naturalCommandFormation, func(g *game) { g.form.active = true }},
		{"武將", naturalCommandGeneral, func(g *game) { g.list = &listwin.List{} }},
		{"勢力", naturalCommandFaction, func(g *game) { g.list = &listwin.List{} }},
	} {
		g := &game{cmdCell: -1}
		if got := g.activeCommandCell(); got != -1 {
			t.Errorf("%s：什麼都沒開就亮了第 %d 格", tc.name, got)
		}
		// dispatchNaturalCommand 在動作之前就記下格號（原版也是先反白）。
		g.cmdCell = int(tc.cell)
		tc.open(g)
		if got := g.activeCommandCell(); got != int(tc.cell) {
			t.Errorf("%s：流程開著時亮第 %d 格，want %d", tc.name, got, tc.cell)
		}
		// 流程結束 → syncCommandFlow 把反白與狀態列一起收掉。
		*g = game{cmdCell: int(tc.cell), statusBox: statusBoxState{lines: []string{"x"}}}
		g.syncCommandFlow()
		if got := g.activeCommandCell(); got != -1 {
			t.Errorf("%s：流程結束了還亮著第 %d 格", tc.name, got)
		}
		if g.statusBoxActive() {
			t.Errorf("%s：流程結束了狀態列提示還掛著", tc.name)
		}
	}
}

// 彈出選單那三格走的是另一條路（`popupMenu.cell`），收掉選單就不亮。
func TestActiveCommandCellStillFollowsPopupMenus(t *testing.T) {
	g := &game{cmdCell: -1}
	g.openPopupMenu(cityPopupMenu)
	if got, want := g.activeCommandCell(), int(naturalCommandCity); got != want {
		t.Errorf("據點選單開著時亮第 %d 格，want %d", got, want)
	}
	g.closePopupMenu()
	if got := g.activeCommandCell(); got != -1 {
		t.Errorf("關掉之後還亮著第 %d 格", got)
	}
}

// 選完之後選單框**留在畫面上**，一覽表畫在它上面（docs/spec/126 §1.2）。
//
// ⭐ 原版不擦，只是往上面畫——實測跑兩千萬道指令那條殘影都不會消
// （docs/playtest/83）。
func TestPopupMenuStaysDrawnAfterPick(t *testing.T) {
	g := &game{cmdCell: int(naturalCommandCity)}
	g.openPopupMenu(cityPopupMenu)
	if !g.popupMenuActive() || !g.popupMenuShown() {
		t.Fatal("剛開的選單應該既能操作也畫得出來")
	}
	g.cmdMenu.stale = true // dispatchPopupMenu 做的事
	if g.popupMenuActive() {
		t.Error("選走之後不該再吃輸入")
	}
	if !g.popupMenuShown() {
		t.Error("選走之後框還要留在畫面上")
	}
	// 流程結束 → 連框一起收掉（原版靠重畫地圖擦掉）。
	g.syncCommandFlow()
	if g.popupMenuShown() {
		t.Error("流程結束了框還留著")
	}
}
