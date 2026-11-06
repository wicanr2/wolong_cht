package main

import (
	"encoding/json"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
	"github.com/wicanr2/wolong_cht/translations"
)

func TestDesktopGestureOwnsReleaseOutsideBar(t *testing.T) {
	var bar desktopGesture
	enabled := func(desktopAction) bool { return true }
	x := desktopButton(desktopSave).Min.X + 10
	for _, p := range []struct {
		x, y int
		down bool
	}{{x, 12, true}, {x, 160, true}, {x, 160, false}} {
		a, blocked := bar.pointer(p.x, p.y, [3]bool{p.down}, enabled)
		if !blocked || a != desktopNone {
			t.Fatalf("gesture leaked: action=%v blocked=%v", a, blocked)
		}
	}
	if a, blocked := bar.pointer(x, 160, [3]bool{}, enabled); blocked || a != desktopNone {
		t.Fatal("release ownership persisted")
	}
}

func TestDesktopGestureRechecksAvailability(t *testing.T) {
	var bar desktopGesture
	allowed := true
	enabled := func(desktopAction) bool { return allowed }
	x := desktopButton(desktopLoad).Min.X + 10
	bar.pointer(x, 12, [3]bool{true}, enabled)
	allowed = false
	if a, blocked := bar.pointer(x, 12, [3]bool{}, enabled); !blocked || a != desktopNone {
		t.Fatal("disabled action dispatched")
	}
	allowed = true
	bar.pointer(x, 12, [3]bool{true}, enabled)
	if a, blocked := bar.pointer(x, 12, [3]bool{}, enabled); !blocked || a != desktopLoad {
		t.Fatal("fresh click was not dispatched")
	}
	if a, _ := bar.pointer(x, 12, [3]bool{}, enabled); a != desktopNone {
		t.Fatal("release repeated action")
	}
}

func TestDesktopRightGestureCannotDispatch(t *testing.T) {
	var bar desktopGesture
	enabled := func(desktopAction) bool { return true }
	bar.pointer(20, 12, [3]bool{false, true}, enabled)
	if a, blocked := bar.pointer(200, 160, [3]bool{}, enabled); !blocked || a != desktopNone {
		t.Fatal("right release leaked into game")
	}
}

func TestDesktopModalDisablesHostActions(t *testing.T) {
	g := &game{world: &state.World{}}
	if !g.desktopEnabled(desktopSave) {
		t.Fatal("idle game cannot save")
	}
	g.saveUI.active = true
	for _, a := range []desktopAction{desktopLanguage, desktopSettings, desktopSave, desktopLoad} {
		if g.desktopEnabled(a) {
			t.Fatalf("modal permits host action %v", a)
		}
	}
	if !g.desktopEnabled(desktopBack) || !g.desktopEnabled(desktopConfirm) {
		t.Fatal("modal navigation unavailable")
	}
}

func TestDesktopBackClosesOnlySaveOwner(t *testing.T) {
	g := &game{list: &listwin.List{}, pointerAction: desktopBack}
	g.saveUI.active = true
	g.updateSaveUI()
	if g.saveUI.active || g.list == nil {
		t.Fatal("back did not stop at the current modal")
	}
}

func TestDesktopBackClosesFormationBeforeList(t *testing.T) {
	g := &game{list: &listwin.List{}, pointerAction: desktopBack}
	g.form.active = true
	g.updateForm()
	if g.form.active || g.list == nil {
		t.Fatal("back escaped past the formation owner")
	}
	g.updateListUI()
	if g.list != nil {
		t.Fatal("second back did not close the list")
	}
}

func TestDesktopSaveOwnerCanConfirmOverSystem(t *testing.T) {
	g := &game{world: &state.World{}, hud: hudSystem}
	g.saveUI.active = true
	if !g.desktopEnabled(desktopConfirm) {
		t.Fatal("background system panel blocked the save owner")
	}
}

func TestDesktopToolbarPointerKeepsGameCoordinates(t *testing.T) {
	old, offset := desktopPointer, desktopGameOffset
	defer func() { desktopPointer, desktopGameOffset = old, offset }()
	desktopGameOffset = desktopBarHeight
	desktopPointer.active, desktopPointer.rebase = true, 0
	desktopPointer.x, desktopPointer.y = 100, 0
	desktopPointer.rawX, desktopPointer.rawY = 100, 100
	if x, y := consumeDesktopPointer(100, 80); x != 0 || y != 0 || desktopPointer.y != -20 {
		t.Fatal("toolbar cannot be reached without scrolling")
	}
	if x, y := consumeDesktopPointer(100, 100); x != 0 || y != 0 || desktopPointer.y != 0 {
		t.Fatal("game origin shifted")
	}
}

func TestDesktopLabelsFitAllLanguages(t *testing.T) {
	data, err := translations.Files.ReadFile("desktop-controls.json")
	if err != nil {
		t.Fatal(err)
	}
	var labels map[string][]string
	if err := json.Unmarshal(data, &labels); err != nil {
		t.Fatal(err)
	}
	for _, lang := range []string{"zh-hant", "zh-hans", "ja", "en"} {
		if len(labels[lang]) != 7 {
			t.Fatalf("missing labels: %s", lang)
		}
		for _, label := range labels[lang] {
			width := 0
			for _, r := range label {
				width += textdraw.RuneWidth(r)
			}
			if width == 0 || width > 80 {
				t.Fatalf("%s label overflows: %q (%d)", lang, label, width)
			}
		}
	}
}
