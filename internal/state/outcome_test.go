package state

import (
	"reflect"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

func TestOutcomeTrustZeroLatchesAtMutationBoundary(t *testing.T) {
	w := &World{Trust: 1}
	before := w.Clock
	if got := w.AdjustTrust(-1); got != 0 {
		t.Fatalf("Trust = %d，want 0", got)
	}
	if got := w.Outcome(); got != DefeatTrustZero {
		t.Fatalf("Outcome = %v，want DefeatTrustZero", got)
	}
	if selector, ok := w.OutcomeMessageSelector(); !ok || selector != 0x019E {
		t.Fatalf("selector = %#x, %v，want 0x019E, true", selector, ok)
	}
	// 已經離開主循環後，Tick 不可再改時鐘或觸發 AI／據點副作用。
	w.Tick(rng.NewFixed(1))
	if !reflect.DeepEqual(w.Clock, before) {
		t.Fatalf("outcome 後 Clock 改變：%+v → %+v", before, w.Clock)
	}

	// 已為零但沒有 1→0 transition 時不應自行建立敗北結果。
	zero := &World{Trust: 0}
	zero.AdjustTrust(-1)
	if got := zero.Outcome(); got != InProgress {
		t.Fatalf("已為 0 的 Trust 不應觸發 outcome：%v", got)
	}
}

func outcomeCaptureWorld(player, old int, withAlternative bool) *World {
	w := &World{Player: player}
	for i := range w.Cities {
		w.Cities[i].Owner = combat.NeutralFaction
		w.Cities[i].OwnerRecorded = combat.NeutralFaction
	}
	w.Factions[old] = Faction{Alive: true, Capital: 0, Cities: 1}
	w.Factions[player].Alive = true
	w.Cities[0].Owner = old
	w.Cities[0].OwnerRecorded = old
	if withAlternative {
		w.Cities[1].Owner = old
		w.Cities[1].OwnerRecorded = old
		w.Factions[old].Cities = 2
	}
	attacker := player
	if attacker == old {
		attacker = (old + 1) % numFactions
	}
	w.Factions[attacker].Alive = true
	w.Corps[0] = Corps{Alive: true, Faction: attacker, Node: 0}
	return w
}

func TestOutcomeCaptureLastCapitalDefeatsPlayer(t *testing.T) {
	w := outcomeCaptureWorld(1, 1, false)
	ev := &CorpsEvent{Captured: -1}
	w.capture(0, ev)
	if w.Factions[1].Alive || w.Factions[1].Capital != noCity {
		t.Fatalf("玩家勢力未按最後首都邊界清除：%+v", w.Factions[1])
	}
	if got := w.Outcome(); got != DefeatFactionEliminated {
		t.Fatalf("Outcome = %v，want DefeatFactionEliminated", got)
	}
}

func TestOutcomeCaptureWithAlternativeCapitalDoesNotDefeat(t *testing.T) {
	w := outcomeCaptureWorld(1, 1, true)
	ev := &CorpsEvent{Captured: -1}
	w.capture(0, ev)
	if !w.Factions[1].Alive || w.Factions[1].Capital != 1 {
		t.Fatalf("仍有替代據點卻未遷都：%+v", w.Factions[1])
	}
	if got := w.Outcome(); got != InProgress {
		t.Fatalf("仍有替代據點不應敗北：%v", got)
	}
}

func TestOutcomeNonPlayerEliminationDoesNotSetPlayerOutcome(t *testing.T) {
	w := outcomeCaptureWorld(1, 2, false)
	ev := &CorpsEvent{Captured: -1}
	w.capture(0, ev)
	if w.Factions[2].Alive || w.Factions[2].Capital != noCity {
		t.Fatalf("非玩家勢力未清除：%+v", w.Factions[2])
	}
	if got := w.Outcome(); got != InProgress {
		t.Fatalf("非玩家消滅不應設定玩家 outcome：%v", got)
	}
}

func TestOutcomeLatchIsNotOverwritten(t *testing.T) {
	w := &World{Trust: 1}
	w.AdjustTrust(-1)
	if !w.latchOutcome(DefeatFactionEliminated) {
		// 第二次 latch 必須拒絕；這個分支本身是預期結果。
	} else {
		t.Fatal("已 latch 的 outcome 不應允許第二個原因覆蓋")
	}
	if got := w.Outcome(); got != DefeatTrustZero {
		t.Fatalf("Outcome 被覆蓋成 %v", got)
	}
}
