package main

import (
	"image"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestLauncherNewGameSelectionReachesConfirmedWorldRequest(t *testing.T) {
	l := newLauncher(false, nil)
	if got := l.apply(launcherConfirm); got.kind != launcherNoResult || l.phase != launcherNewGameConfirm {
		t.Fatalf("title confirm = %#v phase=%v", got, l.phase)
	}
	if got := l.apply(launcherConfirm); got.kind != launcherNoResult || l.phase != launcherScenario {
		t.Fatalf("new game confirm = %#v phase=%v", got, l.phase)
	}
	players := []launcherPlayer{{ID: 7, Lord: "曹操", Capital: "許昌"}}
	if !l.setScenarioPlayers(2, "第三章　蜀地偏安", players) {
		t.Fatal("valid scenario player list was rejected")
	}
	// 選完劇本先進**勢力清單**，再進君主卡（docs/spec/79 §1）。
	if got := l.apply(launcherConfirm); got.kind != launcherNoResult || l.phase != launcherSelectPlayer {
		t.Fatalf("faction confirm = %#v phase=%v", got, l.phase)
	}
	if got := l.apply(launcherConfirm); got.kind != launcherNoResult || l.phase != launcherGameConfirm {
		t.Fatalf("player confirm = %#v phase=%v", got, l.phase)
	}
	got := l.apply(launcherConfirm)
	if got.kind != launcherStartNewGame || got.scenario != 2 || got.player != 7 {
		t.Fatalf("start request = %#v, want scenario 2 player 7", got)
	}
}

func TestLauncherCancelReturnsToPreviousScreen(t *testing.T) {
	l := newLauncher(false, nil)
	l.apply(launcherConfirm)
	if got := l.apply(launcherCancel); got.kind != launcherNoResult || l.phase != launcherTitle {
		t.Fatalf("cancel new game = %#v phase=%v", got, l.phase)
	}
	l.apply(launcherConfirm)
	l.apply(launcherConfirm)
	l.setScenarioPlayers(1, "第二章　赤壁之戰", []launcherPlayer{{ID: 3, Lord: "劉備", Capital: "成都"}})
	if got := l.apply(launcherCancel); got.kind != launcherNoResult || l.phase != launcherScenario || l.cursor != 1 {
		t.Fatalf("cancel faction list = %#v phase=%v cursor=%d", got, l.phase, l.cursor)
	}
}

func TestLauncherRejectsIllegalPlayerSelection(t *testing.T) {
	l := newLauncher(false, nil)
	l.phase = launcherScenario
	if l.setScenarioPlayers(0, "劇本 1", []launcherPlayer{{ID: 99, Lord: "不存在", Capital: "不存在"}}) == false {
		t.Fatal("state model should accept the preview payload; legality is checked at data boundary")
	}
	if l.selectPlayer(7) {
		t.Fatal("player id absent from the preview must not be selectable")
	}
	w := &state.World{
		Factions: [22]state.Faction{{Alive: true, Lord: 0, Capital: 0}},
		Generals: [127]state.General{{Alive: true}},
		Cities:   [192]state.City{{Name: "許昌"}},
	}
	if validLauncherPlayer(w, 99) || validLauncherPlayer(w, -1) {
		t.Fatal("out-of-range player was accepted")
	}
	if !validLauncherPlayer(w, 0) {
		t.Fatal("known valid player was rejected")
	}
}

func TestLauncherEmptySaveSlotCannotBeRead(t *testing.T) {
	l := newLauncher(true, []launcherSlot{
		{Slot: 0, Available: false, Label: "空白槽位"},
		{Slot: 1, Available: true, Label: "196年4月1日　曹操"},
	})
	l.selectRow(1)
	if got := l.apply(launcherConfirm); got.kind != launcherNoResult || l.phase != launcherLoad {
		t.Fatalf("load menu = %#v phase=%v", got, l.phase)
	}
	if got := l.apply(launcherConfirm); got.kind != launcherNoResult || l.notice == "" {
		t.Fatalf("empty slot read = %#v notice=%q", got, l.notice)
	}
	l.selectRow(1)
	if got := l.apply(launcherConfirm); got.kind != launcherStartLoad || got.slot != 1 {
		t.Fatalf("valid slot read = %#v", got)
	}
}

func TestLauncherSaveSlotsUseOneOverlayAndRejectEmptyBlocks(t *testing.T) {
	sourcePath := filepath.Join("..", "..", "workplace", "orig", "dosv", "SINARIO.DAT")
	w, err := state.LoadScenario(sourcePath, 0)
	if err != nil {
		t.Fatal(err)
	}
	// 只把第 1 槽做成「玩家存過的」，另外三個區塊維持原版的空玩家資料——
	// 所以 inspectLauncherSlots 必須從這**一個** overlay 檔裡
	// 仍然回報**四個**邏輯槽位。
	w.Player = 0
	overlay := filepath.Join(t.TempDir(), "SAVE.DAT")
	if err := w.SaveInto(sourcePath, overlay, 0); err != nil {
		t.Fatal(err)
	}
	slots := inspectLauncherSlots(overlay)
	if len(slots) != 4 {
		t.Fatalf("slot count = %d, want 4", len(slots))
	}
	if !slots[0].Available || hasAvailableLauncherSlot(slots) == false {
		t.Fatalf("slot 1 should be the only available slot: %#v", slots)
	}
	for i := 1; i < len(slots); i++ {
		if slots[i].Available {
			t.Fatalf("empty slot %d was marked available", i+1)
		}
	}
}

func TestLauncherNewGameAlwaysUsesScenarioSource(t *testing.T) {
	source := filepath.Join("workplace", "orig", "dosv", "SINARIO.DAT")
	overlay := filepath.Join("tmp", "SAVE.DAT")
	if got := launcherNewGamePath(source, overlay); got != source {
		t.Fatalf("new game path = %q, want source %q", got, source)
	}
	if hasAvailableLauncherSlot(nil) {
		t.Fatal("empty slot list must not enable LOAD DATA")
	}
}

func TestLauncherConfirmHitRectsMatchDrawnRows(t *testing.T) {
	checks := []struct {
		phase launcherPhase
		wantY int
	}{
		{launcherTitle, 152},
		{launcherNewGameConfirm, 184},
		{launcherGameConfirm, 192},
	}
	for _, tc := range checks {
		if got := launcherRowRect(tc.phase, 0).Min.Y; got != tc.wantY {
			t.Errorf("phase %v row 0 y = %d, want drawn y %d", tc.phase, got, tc.wantY)
		}
	}
}

// ⚠ 這裡原本有一份 `TestLauncherLayoutContainsTenPlayerListAndHint`：
// 它驗的是啟動殼層自己的「選君主」清單，而那份清單在 2026-08-24 被原版的
// 勢力清單取代了（docs/spec/79）。**留著會驗一個畫不出來的畫面**——
// 那正是這一輪要修掉的那種東西。現在的版面測試在 `factionlist_test.go`。

func assertRectInside(t *testing.T, got, container image.Rectangle) {
	t.Helper()
	if !got.Min.In(container) || !image.Pt(got.Max.X-1, got.Max.Y-1).In(container) {
		t.Fatalf("rectangle %v escapes safe rect %v", got, container)
	}
}
