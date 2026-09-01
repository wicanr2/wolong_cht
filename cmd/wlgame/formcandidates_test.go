package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
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

// 編成成功之後主將那一句：組編號 `0x19B` 由**原始的說話類型**展開
// （docs/spec/109 §2）。
func TestFormationLeaderTalkIndex(t *testing.T) {
	if formLeaderTalkGroup != 0x19B {
		t.Fatalf("組編號 ＝ %#x，原版是 0x19B（sub_16C5E 的確定分支）", formLeaderTalkGroup)
	}
	for _, tc := range []struct{ variant, want int }{
		{0, 446}, {3, 449}, {7, 453},
	} {
		if got := resolveBattleTalkIndex(formLeaderTalkGroup, tc.variant); got != tc.want {
			t.Errorf("說話類型 %d → #%d，want #%d", tc.variant, got, tc.want)
		}
	}
}

// ⭐ 主公型（說話類型 0–2）那三格在 `TALK.DAT` 裡是空的——原版的編成候選
// 排除君主，所以那條路走不到。remake 允許君主編成，於是會取到空字串，
// **不可以因此跳一張空框**（docs/spec/109 §2）。
func TestFormationLeaderTalkSkipsEmpty(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	g := &game{lib: lib, world: w}

	for v := 0; v <= 2; v++ {
		idx := resolveBattleTalkIndex(formLeaderTalkGroup, v)
		before := len(g.messages)
		g.enqueueTalkWithPortrait(idx, nil, 0)
		if len(g.messages) != before {
			t.Errorf("說話類型 %d（#%d）是空的，卻開了框", v, idx)
		}
	}
	// 臣下型要開得起來，否則上面那一條會被「什麼都不開」蒙混過去。
	idx := resolveBattleTalkIndex(formLeaderTalkGroup, 7)
	before := len(g.messages)
	g.enqueueTalkWithPortrait(idx, nil, 0)
	if len(g.messages) == before {
		t.Errorf("說話類型 7（#%d）有內容，卻沒開框", idx)
	}
}
