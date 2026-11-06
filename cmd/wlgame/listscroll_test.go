package main

import (
	"image"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
)

func scrollTestList() *listwin.List {
	rows := make([]int, 22)
	for i := range rows {
		rows[i] = i
	}
	return listwin.New(listwin.Factions, nil, rows, 10, nil)
}

func TestListScrollOriginalSamples(t *testing.T) {
	// dosgolem正常22列勢力表；原版18503及放開後181BC，spec274。
	for _, s := range []struct{ y, top, offset int }{{128, 0, 0}, {152, 0, 0}, {184, 6, 35}, {240, 12, 70}, {-20, 0, 0}, {400, 12, 70}} {
		top := listScrollTopAt(s.y, listScrollTrackRect(), 22, 10)
		y, h := listScrollThumb(listScrollTrackRect(), top, 22, 10)
		if top != s.top || y != 120+s.offset || h != 58 {
			t.Fatalf("y=%d: top/offset/height=%d/%d/%d", s.y, top, y-120, h)
		}
	}
	for _, total := range []int{0, 5, 10} {
		if listScrollTopAt(240, listScrollTrackRect(), total, 10) != 0 {
			t.Fatal("short list scrolled")
		}
	}
	if listScrollTopAt(240, image.Rectangle{}, 22, 10) != 0 {
		t.Fatal("empty track scrolled")
	}
}

func TestListScrollPreviewsUntilAllButtonsReleased(t *testing.T) {
	l := scrollTestList()
	g := &game{list: l, listDrag: listScrollDrag{list: l}}
	for _, down := range [][3]bool{{true, false, false}, {false, true, false}, {false, false, true}} {
		if !g.updateListScrollDrag(184, down) || l.Top != 0 || g.listDrag.top != 6 || l.Phase() != listwin.Browsing {
			t.Fatalf("preview changed committed state: %v", l)
		}
	}
	if !g.updateListScrollDrag(240, [3]bool{}) || l.Top != 12 || g.listScrollOwnsInput() || l.Phase() != listwin.Browsing {
		t.Fatal("release did not commit without selecting")
	}
}

func TestListScrollRejectsReplacedOwner(t *testing.T) {
	l := scrollTestList()
	g := &game{list: scrollTestList(), listDrag: listScrollDrag{list: l, top: 12}}
	if g.updateListScrollDrag(240, [3]bool{}) || g.list.Top != 0 {
		t.Fatal("old release changed replacement list")
	}
	l = g.list
	g.listDrag = listScrollDrag{list: l, top: 12}
	l.Select(1)
	l.Confirm()
	if g.updateListScrollDrag(240, [3]bool{}) || l.Top != 0 {
		t.Fatal("selected list accepted drag")
	}
}

func TestLauncherScrollUsesSamePreviewAndCommit(t *testing.T) {
	l := &launcherModel{phase: launcherSelectFaction, players: make([]launcherPlayer, 22)}
	g := &game{launcher: l, listDrag: listScrollDrag{launcher: l}}
	if !g.updateListScrollDrag(200, [3]bool{true, false, false}) || l.factionTop != 0 || g.listDrag.top != 6 {
		t.Fatal("launcher preview differs")
	}
	g.updateListScrollDrag(256, [3]bool{})
	if l.factionTop != 12 || l.factionSelected || l.cursor != 0 {
		t.Fatal("launcher release changed selection")
	}
	g.listDrag = listScrollDrag{launcher: l, top: 12}
	l.phase = launcherScenario
	if g.updateListScrollDrag(200, [3]bool{}) || l.factionTop != 12 {
		t.Fatal("closed launcher accepted old drag")
	}
}
