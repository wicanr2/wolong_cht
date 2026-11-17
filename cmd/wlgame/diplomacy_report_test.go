package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestEnvoyReportMissingTalkStillSettles(t *testing.T) {
	g := &game{world: &state.World{Player: -1}}
	settled := 0
	g.enqueueEventMessages(state.Event{
		TalkNotices: []state.TalkNotice{{Index: 57}, {Index: 44}},
		DiplomacyReport: &state.DiplomacyReport{AfterNotice: 1, Confirm: func() *state.DiplomacySettlement {
			settled++
			return nil
		}},
	})
	if settled != 1 || g.messageActive() {
		t.Fatal("缺TALK時交易遺失，或留下無法確認的視窗")
	}
}

func TestDiplomacyStagesMissingTalkComplete(t *testing.T) {
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player, w.Trust = 0, 255
	w.Friendship[0][12], w.Friendship[12][0] = diplomacy.Peace(40), diplomacy.Peace(50)
	w.SetDiplomacyReportConfirmation(true)
	w.Generals[39].Faction, w.Generals[39].Captor, w.Generals[39].Duty = 2, 0, state.DutyCaptive
	w.Factions[0].Generals--
	c := state.DiplomacyChoice{Kind: state.DiplomacyCooperation, Source: 0, Invader: 12, Target: 2,
		InitialResponse: 0, InitialAmount: 0}
	w.BeginDiplomacyChoice(c)
	g := &game{world: w}
	g.resolveDiplomacy(c, state.DiplomacyAcceptFree)
	if g.messageActive() || w.Generals[39].Faction != 0 || w.Generals[39].Captor != 255 || !w.Friendship[0][12].AtWar() {
		t.Fatal("缺TALK時未完整返還或宣戰，或留下無法確認的視窗")
	}
}
