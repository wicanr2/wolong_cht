package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestDiplomacyTalkIndicesMatchRaw13902Branches(t *testing.T) {
	tests := []struct {
		name         string
		kind         state.DiplomacyKind
		option       state.DiplomacyOption
		wantChoice   int
		wantResponse int
	}{
		{"ceasefire-free", state.DiplomacyCeasefire, state.DiplomacyAcceptFree, 364, 0},
		{"ceasefire-funds", state.DiplomacyCeasefire, state.DiplomacyOfferFunds, 365, 1},
		{"ceasefire-reject", state.DiplomacyCeasefire, state.DiplomacyReject, 366, 2},
		{"cooperation-free", state.DiplomacyCooperation, state.DiplomacyAcceptFree, 377, 0},
		{"cooperation-funds", state.DiplomacyCooperation, state.DiplomacyOfferFunds, 378, 1},
		{"cooperation-reject", state.DiplomacyCooperation, state.DiplomacyReject, 379, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := state.DiplomacyChoice{
				Kind: tt.kind, Source: 1, Invader: 1,
				InitialAmount: 6000, OfferAmount: 6000,
			}
			if got := diplomacyTalkChoiceIndex(c, tt.option); got != tt.wantChoice {
				t.Fatalf("choice TALK index = %d, want %d", got, tt.wantChoice)
			}
			gotResponse, _ := diplomacyTalkResponse(c, tt.option)
			if gotResponse != tt.wantResponse {
				t.Fatalf("response = %d, want %d", gotResponse, tt.wantResponse)
			}
		})
	}

	c := state.DiplomacyChoice{
		Kind: state.DiplomacyCeasefire, Source: 1,
		InitialAmount: 6000, OfferAmount: 7000,
	}
	if got, _ := diplomacyTalkResponse(c, state.DiplomacyOfferFunds); got != 2 {
		t.Fatalf("超額外交金額 response = %d, want 2", got)
	}
}

func TestDiplomacyPromptUsesGeneralTalkVariant(t *testing.T) {
	c := state.DiplomacyChoice{Kind: state.DiplomacyCeasefire}
	if got := diplomacyTalkPromptIndex(c, 0); got != 360 {
		t.Fatalf("variant 0 prompt = %d，want 360", got)
	}
	if got := diplomacyTalkPromptIndex(c, 1); got != 361 {
		t.Fatalf("variant 1 prompt = %d，want 361", got)
	}
	if got := diplomacyTalkPromptIndex(c, 4); got != 361 {
		t.Fatalf("raw variant 4 prompt = %d，want 361", got)
	}
}

func TestDiplomacyTalkExpansionUsesOriginalRequestMarkers(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	requester := 1
	g := &game{lib: lib, world: w}

	g.enqueueTalkNotice(state.TalkNotice{
		Index: 0x168, City: -1, Faction: requester, General: -1, Amount: -1,
	})
	if len(g.messages) != 1 {
		t.Fatalf("事件 3 前置訊息數 = %d，want 1", len(g.messages))
	}
	ceasefire := strings.Join(g.messages[0].lines, "\n")
	if strings.Contains(ceasefire, "{") || !strings.Contains(ceasefire, "停戰") ||
		!strings.Contains(ceasefire, big5(w.LordName(requester))) {
		t.Fatalf("事件 3 前置 TALK 展開錯誤：%q", ceasefire)
	}

	g.messages = nil
	g.enqueueTalkNotice(state.TalkNotice{
		Index: 0x175, City: -1, Faction: requester, General: -1, Amount: -1,
	})
	if len(g.messages) != 1 {
		t.Fatalf("事件 2 前置訊息數 = %d，want 1", len(g.messages))
	}
	cooperation := strings.Join(g.messages[0].lines, "\n")
	if strings.Contains(cooperation, "{") || !strings.Contains(cooperation, "協助") ||
		!strings.Contains(cooperation, big5(w.LordName(requester))) {
		t.Fatalf("事件 2 前置 TALK 展開錯誤：%q", cooperation)
	}

	g.messages = nil
	g.enqueueDiplomacyTalk(state.DiplomacyChoice{
		Kind: state.DiplomacyCooperation, Source: 0, Invader: requester,
		InitialAmount: 6000, OfferAmount: 5000,
	}, state.DiplomacyOfferFunds)
	if len(g.messages) != 2 {
		t.Fatalf("事件 2 選擇／結果訊息數 = %d，want 2", len(g.messages))
	}
	choice := strings.Join(g.messages[0].lines, "\n")
	result := strings.Join(g.messages[1].lines, "\n")
	if strings.Contains(choice, "{") || !strings.Contains(choice, "5000") ||
		strings.Contains(result, "{") || !strings.Contains(result, "5000") ||
		!strings.Contains(result, big5(w.LordName(requester))) {
		t.Fatalf("事件 2 選擇／結果 TALK 展開錯誤：choice=%q result=%q", choice, result)
	}
}
