package main

import "testing"

// 以原版回讀的絕對座標驗證，不由被測幾何函式產生預期值。
func TestLauncherDatesSelectTheDisplayedSlot(t *testing.T) {
	for _, phase := range []launcherPhase{launcherScenario, launcherLoad} {
		l := &launcherModel{phase: phase, slots: make([]launcherSlot, 4)}
		for slot, y := range []int{144, 192, 240, 288} {
			for _, p := range [][2]int{{256, y}, {300, y + 8}, {375, y + 15}} {
				row, ok := l.pointerRow(p[0], p[1])
				if !ok || row != slot {
					t.Errorf("phase=%v (%d,%d) 選到 %d/%v，預期槽 %d", phase, p[0], p[1], row, ok, slot)
				}
			}
			for _, p := range [][2]int{{255, y}, {376, y}, {300, y - 1}, {300, y + 16}, {220, y - 16}} {
				if row, ok := l.pointerRow(p[0], p[1]); ok {
					t.Errorf("日期欄外 (%d,%d) 不應選到 %d", p[0], p[1], row)
				}
			}
		}
	}
}

func TestLauncherFactionSelectionMatchesOracle(t *testing.T) {
	l := &launcherModel{phase: launcherScenario}
	l.setScenarioPlayers(0, "第一章", []launcherPlayer{{ID: 0}, {ID: 4}, {ID: 7}})
	if l.factionSelected {
		t.Fatal("剛進清單不應反白")
	}
	l.clickFactionRow(0)
	if l.phase != launcherSelectFaction || !l.factionSelected || l.cursor != 0 {
		t.Fatal("第一次點擊應只反白第一列")
	}
	l.apply(launcherCancel)
	if l.phase != launcherSelectFaction || l.factionSelected {
		t.Fatal("第一次右鍵應只解除反白")
	}
	l.clickFactionRow(1)
	l.clickFactionRow(2)
	if l.phase != launcherSelectPlayer || l.playerIndex() != 4 {
		t.Fatal("反白後點另一列應確認原選取的勢力")
	}
	l.apply(launcherCancel)
	if l.phase != launcherSelectFaction || l.factionSelected {
		t.Fatal("君主卡右鍵應回未反白的清單")
	}
	l.clickFactionRow(99)
	if l.phase != launcherSelectFaction || l.factionSelected {
		t.Fatal("空列不應選取或確認")
	}
	l.apply(launcherCancel)
	if l.phase != launcherScenario {
		t.Fatal("瀏覽狀態右鍵應回劇本")
	}
}
