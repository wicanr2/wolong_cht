package state

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

func partialTickFixture(t *testing.T, secondBattle bool) *World {
	t.Helper()
	w, att, def := standoffFixture(t)
	a, b := w.Corps[att], w.Corps[def]
	ga, gb := w.Generals[att], w.Generals[def]
	w.Corps = [numCorps]Corps{}
	w.Corps[0], w.Corps[32] = a, b
	w.Generals[0], w.Generals[32] = ga, gb
	w.Generals[0].Sovereign, w.Generals[32].Sovereign = false, false
	w.Player = a.Faction
	w.Factions[a.Faction].Funds = 20000
	w.Cities[0].Owner = combat.NeutralFaction
	w.cityCursor, w.corpsCursor = 0, 0
	w.Clock.Hour, w.Clock.Subtick = upkeepHour, 0
	w.SetTactical(&TacticalSetup{Forms: tactical.SyntheticFormations(),
		Field: func(int, bool, bool) *tactical.Field {
			stack := make([][]int, tactical.Height)
			for y := range stack {
				stack[y] = make([]int, tactical.Width)
			}
			return tactical.NewField(stack, 0)
		}})
	x, y := emptyPair(t, w)
	faceOff(w, 0, 32, x, y)
	w.Corps[0].Timer, w.Corps[0].Interval = 1, 3
	w.Corps[0].Standoff, w.Corps[0].Countdown = true, 1
	for _, i := range []int{4, 16} {
		w.Corps[i] = a
		w.Corps[i].X, w.Corps[i].Y = 20+i, 20
		w.Corps[i].Timer, w.Corps[i].Interval = 7, 7
		w.Corps[i].Standoff, w.Corps[i].Countdown = true, 3
	}
	if secondBattle {
		w.Corps[8], w.Corps[40] = a, b
		w.Generals[8], w.Generals[40] = ga, gb
		w.Generals[8].Sovereign, w.Generals[40].Sovereign = false, false
		faceOff(w, 8, 40, x+12, y+4)
		w.Corps[8].Timer, w.Corps[8].Interval = 1, 3
		w.Corps[8].Standoff, w.Corps[8].Countdown = true, 1
	}
	return w
}

// Controlled battle result isolates the scheduler; normal play is tested separately.
func resolvePartialBattle(t *testing.T, w *World, random *rng.Rand) {
	t.Helper()
	p := w.PendingBattle()
	if p == nil {
		t.Fatal("missing pending battle")
	}
	b := tactical.NewBattle(nil, nil, rng.NewFixed(2), 0)
	b.Done, b.Winner = true, 0
	b.Sides[0].Reserve[0] = w.Corps[p.Attacker].Men / 2
	b.Sides[1].Reserve[0] = w.Corps[p.Defender].Men / 2
	p.Battle = b
	if w.ResolvePending(random) == nil {
		t.Fatal("missing result")
	}
}

func TestWorldTickBattleResumePreservesVisitAndBudget(t *testing.T) {
	for _, mapObjects := range []bool{false, true} {
		t.Run(map[bool]string{false: "TickThenTickMap", true: "TickMapThenTick"}[mapObjects], func(t *testing.T) {
			w := partialTickFixture(t, false)
			random := rng.NewFixed(7)
			clockBefore, fundsBefore := w.Clock, w.Factions[w.Player].Funds
			objectsBefore := append([]byte(nil), w.Bytes()[0x21C0:0x22C0]...)
			if mapObjects {
				w.TickMap(random)
			} else {
				w.Tick(random)
			}
			if w.PendingBattle() == nil || !w.HasIncompleteTick() {
				t.Fatal("did not suspend")
			}
			if w.Clock != clockBefore || w.corpsCursor != 0 || w.cityCursor != 1 {
				t.Fatal("time or cursors advanced past battle")
			}
			if w.Factions[w.Player].Funds != fundsBefore || w.Corps[4].Timer != 7 || w.Corps[4].Countdown != 3 {
				t.Fatal("tail or remaining corps ran before battle")
			}
			if !bytes.Equal(objectsBefore, w.Bytes()[0x21C0:0x22C0]) {
				t.Fatal("objects advanced before battle")
			}
			frozen, frozenRNG := w.Bytes(), random.Raw()
			w.Tick(random)
			w.TickMap(random)
			w.ResumeTick(random)
			if !bytes.Equal(frozen, w.Bytes()) || !bytes.Equal(frozenRNG, random.Raw()) {
				t.Fatal("pending reentry changed world")
			}

			resolvePartialBattle(t, w, random)
			if !w.HasIncompleteTick() {
				t.Fatal("result lost unfinished visit")
			}
			fundsAfterBattle := w.Factions[w.Player].Funds
			timerAfterBattle, moraleAfterBattle := w.Corps[0].Timer, w.Corps[0].Morale
			if w.Corps[0].Men != 300 {
				t.Fatalf("post battle men=%d", w.Corps[0].Men)
			}
			if mapObjects {
				w.Tick(random)
			} else {
				w.TickMap(random)
			}
			wantClock := clockBefore
			wantClock.Advance()
			if w.Clock != wantClock || w.corpsCursor != 16 || w.cityCursor != 1 || w.HasIncompleteTick() {
				t.Fatal("did not complete exactly the original update")
			}
			// 300/32+1 = 10 for the winner; 600/32+1 = 19 for corps4.
			if got := fundsAfterBattle - w.Factions[w.Player].Funds; got != 29 {
				t.Fatalf("upkeep=%d, want29 from post-battle men", got)
			}
			if w.Corps[0].Timer != timerAfterBattle || w.Corps[0].Morale != moraleAfterBattle+10 {
				t.Fatal("winner visit repeated or upkeep missing")
			}
			if w.Corps[4].Timer != 6 || w.Corps[4].Countdown != 2 || w.Corps[16].Timer != 7 {
				t.Fatal("remaining batch quota or countdown changed")
			}
			objectsAfter := w.Bytes()[0x21C0:0x22C0]
			if bytes.Equal(objectsBefore, objectsAfter) == mapObjects {
				t.Fatal("resume used the new caller map-object policy")
			}
			finished, finishedRNG := w.Bytes(), random.Raw()
			w.ResumeTick(random)
			if !bytes.Equal(finished, w.Bytes()) || !bytes.Equal(finishedRNG, random.Raw()) {
				t.Fatal("empty resume started another update")
			}
		})
	}
}

func TestWorldTickBattleResumeCanSuspendAgain(t *testing.T) {
	w := partialTickFixture(t, true)
	random := rng.NewFixed(7)
	before := w.Clock
	w.TickMap(random)
	resolvePartialBattle(t, w, random)
	w.ResumeTick(random)
	if p := w.PendingBattle(); p == nil || p.Attacker != 8 || p.Defender != 40 {
		t.Fatal("second battle was skipped")
	}
	if w.Clock != before || w.corpsCursor != 0 || w.cityCursor != 1 || w.Corps[4].Timer != 6 || w.Corps[16].Timer != 7 {
		t.Fatal("second suspension reset or overran the batch")
	}
	resolvePartialBattle(t, w, random)
	w.ResumeTick(random)
	want := before
	want.Advance()
	if w.Clock != want || w.corpsCursor != 16 || w.cityCursor != 1 || w.Corps[4].Timer != 6 || w.Corps[16].Timer != 7 || w.HasIncompleteTick() {
		t.Fatal("second result resumed a fresh batch")
	}
}

func TestWorldTickBattleResumeStopsForOutcome(t *testing.T) {
	w := partialTickFixture(t, false)
	random := rng.NewFixed(7)
	w.TickMap(random)
	w.Trust = 1
	w.AdjustTrust(-1)
	frozen := w.Bytes()
	w.ResumeTick(random)
	if w.HasIncompleteTick() || !bytes.Equal(frozen, w.Bytes()) {
		t.Fatal("outcome continued the unfinished update")
	}
}

func TestWorldTickBattleResumeRejectsPartialOriginalSave(t *testing.T) {
	w := partialTickFixture(t, false)
	random := rng.NewFixed(7)
	w.TickMap(random)
	dst := filepath.Join(t.TempDir(), "SAVE.DAT")
	sentinel := []byte("existing save")
	if err := os.WriteFile(dst, sentinel, 0600); err != nil {
		t.Fatal(err)
	}
	for _, afterResult := range []bool{false, true} {
		if afterResult {
			resolvePartialBattle(t, w, random)
		}
		if err := w.SaveInto("../../workplace/orig/dosv/SINARIO.DAT", dst, 0); !errors.Is(err, ErrIncompleteTick) {
			t.Fatalf("partial save error=%v", err)
		}
		got, err := os.ReadFile(dst)
		if err != nil || !bytes.Equal(got, sentinel) {
			t.Fatal("partial save touched existing file")
		}
	}
	w.ResumeTick(random)
	if err := w.SaveInto("../../workplace/orig/dosv/SINARIO.DAT", dst, 0); err != nil {
		t.Fatal(err)
	}
	back, err := LoadScenario(dst, 0)
	if err != nil || back.Clock != w.Clock || back.corpsCursor != w.corpsCursor {
		t.Fatal("completed save lost resumed progress")
	}
}

func TestWorldTickBattleResumeDeliversEventsOnce(t *testing.T) {
	w := partialTickFixture(t, false)
	w.Corps[6] = Corps{Faction: w.Player, Routing: true, Countdown: 1}
	w.Generals[6].Alive, w.Generals[6].Faction = true, w.Player
	random := rng.NewFixed(7)
	first := w.TickMap(random)
	for _, event := range first.Corps {
		if event.Corps == 6 {
			t.Fatal("later corps event escaped before battle")
		}
	}
	resolvePartialBattle(t, w, random)
	resumed := w.ResumeTick(random)
	count := 0
	for _, event := range resumed.Corps {
		if event.Corps == 6 && event.RoutEnded {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("rout-ended events=%d, want1", count)
	}
	again := w.ResumeTick(random)
	if len(again.Corps) != 0 || len(again.TalkNotices) != 0 || len(again.Strategy) != 0 {
		t.Fatal("resume re-sent events")
	}
}

func TestWorldTickBattleResumeVictoryDiscardsRemainder(t *testing.T) {
	w := partialTickFixture(t, false)
	random := rng.NewFixed(7)
	w.TickMap(random)
	p := w.PendingBattle()
	a, b := w.Corps[p.Attacker].Faction, w.Corps[p.Defender].Faction
	for i := range w.Cities {
		w.Cities[i].Owner = combat.NeutralFaction
	}
	w.Cities[w.Factions[a].Capital].Owner = a
	w.Cities[p.Node].Owner = b
	w.Cities[p.Node].OwnerRecorded = b
	w.Factions[b].Capital, w.Factions[b].Cities = p.Node, 1
	for i := range w.Factions {
		w.Factions[i].Alive = i == a || i == b
	}
	w.LivingFactions = 2
	p.Mode = combat.Siege
	before := w.Clock
	resolvePartialBattle(t, w, random)
	if w.Outcome() != Victory || w.HasIncompleteTick() {
		t.Fatalf("outcome=%v incomplete=%v", w.Outcome(), w.HasIncompleteTick())
	}
	frozen := w.Bytes()
	w.ResumeTick(random)
	w.TickMap(random)
	if w.Clock != before || w.corpsCursor != 0 || w.Corps[4].Timer != 7 || !bytes.Equal(frozen, w.Bytes()) {
		t.Fatal("victory ran the abandoned tail")
	}
}
