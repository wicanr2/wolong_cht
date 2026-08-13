package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestFundingTalkIndicesMatchRaw139E8Branches(t *testing.T) {
	tests := []struct {
		name         string
		kind         state.FundingKind
		option       state.FundingOption
		amount       int
		wantFirst    int
		wantFollowUp int
	}{
		{"governor-full", state.FundingGovernor, state.FundingFullAmount, 6000, 284, 288},
		{"governor-equal", state.FundingGovernor, state.FundingSetAmount, 6000, 284, 293},
		{"governor-less", state.FundingGovernor, state.FundingSetAmount, 5000, 285, 293},
		{"governor-zero", state.FundingGovernor, state.FundingSetAmount, 0, 286, 293},
		{"governor-more", state.FundingGovernor, state.FundingSetAmount, 7000, 287, 293},
		{"governor-reject", state.FundingGovernor, state.FundingReject, 6000, 286, 298},
		{"diplomat-full", state.FundingDiplomat, state.FundingFullAmount, 6000, 325, 329},
		{"diplomat-equal", state.FundingDiplomat, state.FundingSetAmount, 6000, 325, 334},
		{"diplomat-less", state.FundingDiplomat, state.FundingSetAmount, 5000, 326, 334},
		{"diplomat-zero", state.FundingDiplomat, state.FundingSetAmount, 0, 327, 334},
		{"diplomat-more", state.FundingDiplomat, state.FundingSetAmount, 7000, 328, 334},
		{"diplomat-reject", state.FundingDiplomat, state.FundingReject, 6000, 327, 339},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := state.FundingChoice{
				Kind:            tt.kind,
				RequestedAmount: 6000,
				OfferAmount:     tt.amount,
			}
			gotFirst, gotFollowUp := fundingTalkIndices(c, tt.option)
			if gotFirst != tt.wantFirst || gotFollowUp != tt.wantFollowUp {
				t.Fatalf("TALK index = (%d, %d), want (%d, %d)",
					gotFirst, gotFollowUp, tt.wantFirst, tt.wantFollowUp)
			}
		})
	}
}

func TestFundingTalkExpansionUsesOriginalMarkers(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	city := -1
	for i, c := range w.Cities {
		if c.Owner == w.Player {
			city = i
			break
		}
	}
	if city < 0 {
		t.Fatal("找不到玩家據點")
	}
	advisor := w.Factions[w.Player].Advisor
	if advisor < 0 || advisor >= len(w.Generals) || !w.Generals[advisor].Alive {
		t.Fatalf("劇本 1 玩家軍師無效：%d", advisor)
	}

	g := &game{lib: lib, world: w}
	g.enqueueTalkNotice(state.TalkNotice{
		Index: 0x116, City: city, Faction: -1, General: -1, Amount: 6000,
	})
	if len(g.messages) != 1 {
		t.Fatalf("事件 4 前置要求訊息數 = %d，want 1", len(g.messages))
	}
	initial := strings.Join(g.messages[0].lines, "\n")
	if strings.Contains(initial, "{") || !strings.Contains(initial, big5(w.Cities[city].Name)) ||
		!strings.Contains(initial, "6000") || !strings.Contains(initial, big5(w.Generals[advisor].Name)) {
		t.Fatalf("事件 4 marker 展開錯誤：%q", initial)
	}

	g.messages = nil
	g.enqueueFundingTalk(state.FundingChoice{
		Kind: state.FundingGovernor, Subject: city,
		RequestedAmount: 6000, OfferAmount: 7000,
	}, state.FundingSetAmount)
	if len(g.messages) != 2 {
		t.Fatalf("事件 4 結果／收尾訊息數 = %d，want 2", len(g.messages))
	}
	result := strings.Join(g.messages[0].lines, "\n")
	if strings.Contains(result, "{") || !strings.Contains(result, "7000") {
		t.Fatalf("事件 4 指定超額 TALK 展開錯誤：%q", result)
	}
}
