package main

import (
	"image"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
)

func testInputList() *listwin.List {
	return listwin.New(listwin.Generals, []listwin.Column{
		{Title: "名", Less: func(a, b int) bool { return a < b }},
	}, []int{10, 11, 12, 13, 14, 15, 16}, 3, nil)
}

func TestSaveSlotHitTargetsCoverExactlyFourRows(t *testing.T) {
	for slot := 0; slot < 4; slot++ {
		r := saveSlotRect(slot)
		if r.Empty() {
			t.Fatalf("slot %d has empty hit rectangle", slot)
		}
		if !(image.Point{X: r.Min.X + r.Dx()/2, Y: r.Min.Y + r.Dy()/2}).In(r) {
			t.Fatalf("slot %d center does not hit %v", slot, r)
		}
		// 原版四個熱區高 16、**間距 48**（sub_18BD1 的 `add bx, 30h`），
		// 所以它們之間有 32 px 的空隙——不是相鄰的四列。
		if prev := saveSlotRect(slot - 1); slot > 0 {
			if r.Min.Y-prev.Min.Y != saveSlotStep {
				t.Fatalf("slot %d 的間距 = %d，want %d", slot, r.Min.Y-prev.Min.Y, saveSlotStep)
			}
			if prev.Overlaps(r) {
				t.Fatalf("slot %d 與上一列重疊", slot)
			}
		}
	}
	misses := []image.Point{
		{X: saveSlotRect(0).Min.X - 1, Y: saveSlotRect(0).Min.Y + 2},
		{X: saveSlotRect(3).Max.X, Y: saveSlotRect(3).Min.Y + 2},
		{X: saveSlotRect(0).Min.X + 2, Y: saveSlotRect(0).Min.Y - 1},
		{X: saveSlotRect(3).Min.X + 2, Y: saveSlotRect(3).Max.Y},
		// 兩列之間的空隙也不能命中。
		{X: saveSlotRect(0).Min.X + 2, Y: saveSlotRect(0).Max.Y + 4},
	}
	for _, p := range misses {
		for slot := 0; slot < 4; slot++ {
			if p.In(saveSlotRect(slot)) {
				t.Fatalf("outside point %v hit save slot %d", p, slot)
			}
		}
	}
}

func TestListRowHitTargetsRejectMarginsAndGaps(t *testing.T) {
	l := testInputList()
	rows, first := l.Visible()
	for visible := range rows {
		r := listRowRect(l, visible)
		row, ok := listRowAt(l, r.Min.X+r.Dx()/2, r.Min.Y+r.Dy()/2)
		if !ok || row != first+visible {
			t.Fatalf("visible row %d hit = (%d, %v), want %d", visible, row, ok, first+visible)
		}
	}
	misses := []image.Point{
		{X: listWindowX, Y: listRowRect(l, 0).Min.Y + 2},
		{X: listWindowX + listWindowW, Y: listRowRect(l, 0).Min.Y + 2},
		{X: listRowRect(l, 0).Min.X + 2, Y: listRowRect(l, 0).Min.Y - 1},
		{X: listRowRect(l, 2).Min.X + 2, Y: listRowRect(l, 2).Max.Y},
	}
	for _, p := range misses {
		if row, ok := listRowAt(l, p.X, p.Y); ok {
			t.Fatalf("margin point %v hit row %d", p, row)
		}
	}
}

func TestListMouseDispatcherUsesTwoStageConfirm(t *testing.T) {
	chosen := -1
	g := &game{list: testInputList(), listPick: func(id int) bool {
		chosen = id
		return true
	}}
	g.dispatchListAction(listUIAction{kind: listActionClickRow, value: 1})
	if g.list == nil || g.list.Phase() != listwin.Selected || g.list.Cursor != 1 {
		t.Fatalf("first row click did not select: list=%v", g.list)
	}
	g.dispatchListAction(listUIAction{kind: listActionClickRow, value: 1})
	if g.list != nil || chosen != 11 {
		t.Fatalf("second row click did not confirm: list=%v chosen=%d", g.list, chosen)
	}
}

func TestListPageAndFooterCancelAreSharedActions(t *testing.T) {
	g := &game{list: testInputList()}
	g.dispatchListAction(listUIAction{kind: listActionPage, value: 1})
	if g.list.Cursor != 3 || g.list.Top != 1 {
		t.Fatalf("page down cursor/top=%d/%d, want 3/1", g.list.Cursor, g.list.Top)
	}
	g.dispatchListAction(listUIAction{kind: listActionCancel})
	if g.list != nil {
		t.Fatal("cancel from browsing should close the list")
	}
}

func TestSaveModalActionsDoNotTouchBackgroundWindows(t *testing.T) {
	g := &game{open: [numWindows]bool{winCommand: true}, saveUI: saveUIState{
		active: true,
		action: saveWrite,
		slot:   0,
	}}
	g.dispatchSaveUI(saveUIAction{kind: saveActionNext})
	if g.saveUI.slot != 1 || !g.open[winCommand] {
		t.Fatalf("save selection changed background: slot=%d command=%v", g.saveUI.slot, g.open[winCommand])
	}
	g.dispatchSaveUI(saveUIAction{kind: saveActionCancel})
	if g.saveUI.active || !g.open[winCommand] {
		t.Fatalf("save cancel leaked into background: save=%v command=%v", g.saveUI.active, g.open[winCommand])
	}
}
