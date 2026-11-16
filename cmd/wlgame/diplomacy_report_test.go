package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestEnvoyReportMissingTalkStillSettles(t *testing.T) {
	g := &game{world: &state.World{Player: -1}}
	settled := 0
	g.enqueueEventMessages(state.Event{
		TalkNotices: []state.TalkNotice{{Index: 57}, {Index: 44}},
		DiplomacyReport: &state.DiplomacyReport{AfterNotice: 1, Confirm: func() []int {
			settled++
			return nil
		}},
	})
	if settled != 1 || g.messageActive() {
		t.Fatal("缺TALK時交易遺失，或留下無法確認的視窗")
	}
}
