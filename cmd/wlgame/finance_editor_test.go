package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 四列各自的上限來自四支 handler 的 `ax`（docs/spec/78 §1.2）。
func TestFinanceRowMaxComesFromTheHandlers(t *testing.T) {
	if got := financeRowMax(0); got != 100 {
		t.Errorf("稅率上限 = %d，want 100（sub_167CD 的 ax=64h）", got)
	}
	for row := 1; row < financeRows; row++ {
		if got := financeRowMax(row); got != 10000 {
			t.Errorf("第 %d 列上限 = %d，want 10000（ax=2710h）", row, got)
		}
	}
}

// 徵兵那三列輸入的是**人數**，寫回勢力記錄前要除以 10（原版 `div bx, bx=0Ah`）。
func TestFinanceCommitConvertsMenToPoints(t *testing.T) {
	g := &game{world: &state.World{}}
	g.finance = financeState{active: true, editing: true, row: 2, value: 9_876}
	g.commitFinanceAmount()
	if got := g.world.NextRecruitCap[economy.Archer]; got != 987 {
		t.Errorf("弓兵 = %d 點，want 987（9876 ÷ 10 取整）", got)
	}
	if g.finance.editing {
		t.Error("決定之後還開著數值器")
	}

	g.finance = financeState{active: true, editing: true, row: 0, value: 37}
	g.commitFinanceAmount()
	if g.world.NextTaxRate != 37 {
		t.Errorf("稅率 = %d，want 37", g.world.NextTaxRate)
	}
}

// 數字鍵一路按下去要被夾在上限，不能長成六位數。
func TestFinanceAmountClampsToRowMax(t *testing.T) {
	v := 0
	for i := 0; i < 6; i++ {
		got, ok := state.EditAmountValue(v, financeRowMax(0), state.AmountAppendDigit, 9)
		if !ok {
			t.Fatal("數字鍵被拒絕了")
		}
		v = got
	}
	if v != 100 {
		t.Errorf("連按六次 9 之後 = %d，want 100", v)
	}
}
