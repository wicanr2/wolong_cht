package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 君主在不在編成候選裡由系統選單那一列決定（docs/spec/76）。
// **兩個方向都要驗**——只驗一邊的話，開關寫死成任一值都會通過。
func TestFormCandidatesLordToggle(t *testing.T) {
	newGame := func(allow bool) *game {
		w := &state.World{Player: 0}
		w.Factions[0].Alive = true
		w.Factions[0].Lord = 3
		for _, i := range []int{3, 4, 5} {
			w.Generals[i].Alive = true
			w.Generals[i].Faction = 0
			w.Generals[i].Posted = false
			w.Generals[i].Captor = 0xFF
		}
		return &game{world: w, lordCorps: allow}
	}
	has := func(rows []int, n int) bool {
		for _, r := range rows {
			if r == n {
				return true
			}
		}
		return false
	}

	t.Run("切成不可＝原版行為", func(t *testing.T) {
		rows := newGame(false).formCandidates()
		if has(rows, 3) {
			t.Error("關掉之後君主還在候選裡")
		}
		// 反向對照：其他兩位要在，否則「候選永遠是空的」也會通過。
		for _, want := range []int{4, 5} {
			if !has(rows, want) {
				t.Errorf("武將 %d 應該可以帶兵，卻不在候選裡", want)
			}
		}
	})

	t.Run("預設可＝remake 差異", func(t *testing.T) {
		rows := newGame(true).formCandidates()
		if !has(rows, 3) {
			t.Error("開著卻選不到君主")
		}
	})
}

// ⚠ 預設值是 true（放行）。這一條擋的是「哪天有人把預設改掉卻沒改文件」。
func TestSystemMenuHasLordCorpsRow(t *testing.T) {
	if sysMenuLabels[sysRowLordCorps] != "主君編成" {
		t.Errorf("第 %d 列的標籤是 %q", sysRowLordCorps, sysMenuLabels[sysRowLordCorps])
	}
	// 七列要放得進視窗：最後一列的標籤底框不可以掉出下緣。
	bottom := sysLabelBoxY + (sysRows-1)*sysRowStep + sysLabelBoxH
	if bottom > sysWinY+sysWinH {
		t.Errorf("第 %d 列的下緣 %d 超出視窗下緣 %d", sysRows-1, bottom, sysWinY+sysWinH)
	}
	// ⭐ 原版那六列的**位置一格都不能動**，新列只能加在後面。
	want := [6]string{"資料儲存", "畫面模式", "音　　效", "戰略速度", "戰術速度", "遊戲結束"}
	for k, w := range want {
		if sysMenuLabels[k] != w {
			t.Errorf("原版第 %d 列變成 %q，應該是 %q", k, sysMenuLabels[k], w)
		}
	}
	g := &game{lordCorps: true}
	g.dispatchSystemRow(sysRowLordCorps, true)
	if g.lordCorps {
		t.Error("點了那一列卻沒切換")
	}
	g.dispatchSystemRow(sysRowLordCorps, false)
	if !g.lordCorps {
		t.Error("右鍵也要能切回來")
	}
}

// 已帶兵、俘虜、別的勢力、死掉的都不算——這些是原本就有的條件，
// 加君主那一行不可以把它們弄壞。
func TestFormCandidatesKeepsExistingFilters(t *testing.T) {
	w := &state.World{Player: 0}
	w.Factions[0].Alive = true
	w.Factions[0].Lord = 0
	w.Generals[1] = state.General{Alive: true, Faction: 0, Posted: true, Captor: 0xFF}
	w.Generals[2] = state.General{Alive: true, Faction: 1, Captor: 0xFF}
	w.Generals[3] = state.General{Alive: false, Faction: 0, Captor: 0xFF}
	w.Generals[4] = state.General{Alive: true, Faction: 0, Captor: 7}
	w.Generals[5] = state.General{Alive: true, Faction: 0, Captor: 0xFF}
	g := &game{world: w}
	rows := g.formCandidates()
	if len(rows) != 1 || rows[0] != 5 {
		t.Errorf("候選 = %v，只應該有武將 5", rows)
	}
}
