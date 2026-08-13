package main

import (
	"image"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
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
		t.Fatalf("cancel player selection = %#v phase=%v cursor=%d", got, l.phase, l.cursor)
	}
}

func TestLauncherRejectsIllegalPlayerSelection(t *testing.T) {
	l := newLauncher(false, nil)
	l.phase = launcherScenario
	if l.setScenarioPlayers(0, launcherScenarioName(0), []launcherPlayer{{ID: 99, Lord: "不存在", Capital: "不存在"}}) == false {
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
	// Make only slot 1 look like a player save. The other three blocks remain
	// the original empty-player data, so inspectLauncherSlots must still return
	// exactly four logical slots from this one overlay file.
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

func TestLauncherLayoutContainsTenPlayerListAndHint(t *testing.T) {
	safe := launcherTextSafeRect()
	assertRectInside(t, image.Rect(launcherListX, 72,
		launcherListX+textdraw.StringWidth("選擇君主／玩家勢力"), 72+textdraw.GlyphH), safe)
	players := make([]launcherPlayer, 10)
	for i := range players {
		players[i] = launcherPlayer{
			ID:      i,
			Lord:    "司馬懿仲達",
			Capital: "長安郿城",
		}
	}
	l := newLauncher(false, nil)
	l.phase = launcherSelectPlayer
	l.players = players
	l.cursor = len(players) - 1
	start, end := l.visiblePlayers(launcherPlayerMax)
	if got := end - start; got != launcherPlayerMax {
		t.Fatalf("visible player rows = %d, want %d", got, launcherPlayerMax)
	}
	for row, playerIndex := 0, start; playerIndex < end; row, playerIndex = row+1, playerIndex+1 {
		rect := launcherRowRect(launcherSelectPlayer, row)
		assertRectInside(t, rect, safe)
		label := launcherPlayerLabel(players[playerIndex])
		assertRectInside(t, image.Rect(launcherListX, rect.Min.Y+4,
			launcherListX+textdraw.StringWidth(label), rect.Min.Y+4+textdraw.GlyphH), safe)
	}
	hint := image.Rect(launcherListX, launcherHintY,
		launcherListX+textdraw.StringWidth(launcherHint), launcherHintY+textdraw.GlyphH)
	assertRectInside(t, hint, safe)
}

func assertRectInside(t *testing.T, got, container image.Rectangle) {
	t.Helper()
	if !got.Min.In(container) || !image.Pt(got.Max.X-1, got.Max.Y-1).In(container) {
		t.Fatalf("rectangle %v escapes safe rect %v", got, container)
	}
}
