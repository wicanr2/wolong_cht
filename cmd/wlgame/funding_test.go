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
		{"governor-equal", state.FundingGovernor, state.FundingSetAmount, 6000, 284, 288},
		{"governor-less", state.FundingGovernor, state.FundingSetAmount, 5000, 285, 293},
		{"governor-zero", state.FundingGovernor, state.FundingSetAmount, 0, 286, 298},
		{"governor-more", state.FundingGovernor, state.FundingSetAmount, 7000, 287, 303},
		{"governor-reject", state.FundingGovernor, state.FundingReject, 6000, 286, 298},
		{"diplomat-full", state.FundingDiplomat, state.FundingFullAmount, 6000, 325, 329},
		{"diplomat-equal", state.FundingDiplomat, state.FundingSetAmount, 6000, 325, 329},
		{"diplomat-less", state.FundingDiplomat, state.FundingSetAmount, 5000, 326, 334},
		{"diplomat-zero", state.FundingDiplomat, state.FundingSetAmount, 0, 327, 339},
		{"diplomat-more", state.FundingDiplomat, state.FundingSetAmount, 7000, 328, 344},
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

// 原版13C99只減一次3；夏侯淵37在原版樣本的raw7必須使用base+4。
func TestFundingOfficerVariantUsesSingleOriginalFold(t *testing.T) {
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	g := &game{world: w}
	c := state.FundingChoice{Officer: 37}
	for _, sample := range []struct{ raw, want int }{{0, 278}, {2, 280}, {3, 278}, {5, 280}, {6, 281}, {7, 282}} {
		w.Generals[37].TalkVariant = sample.raw
		if got := g.fundingOfficerTalkIndex(c, 278); got != sample.want {
			t.Fatalf("raw%d: index=%d want%d", sample.raw, got, sample.want)
		}
	}
}

// 原版零要求的順序是引見→官員報告→軍師回答→官員答覆，沒有選列。
func TestZeroFundingRequestQueuesOriginalAcknowledgements(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	w.Generals[37].TalkVariant = 7
	for _, sample := range []struct {
		kind                                state.FundingKind
		subject, intro, prompt, first, last int
	}{
		{state.FundingGovernor, 82, 56, 312, 313, 318},
		{state.FundingDiplomat, 1, 57, 353, 354, 359},
	} {
		c := state.FundingChoice{Kind: sample.kind, Subject: sample.subject, Officer: 37}
		n := state.TalkNotice{Index: sample.intro, City: -1, Faction: -1, General: 37, Amount: -1}
		if sample.kind == state.FundingGovernor {
			n.City = sample.subject
		} else {
			n.Faction = sample.subject
		}
		g := &game{lib: lib, world: w}
		g.enqueueEventMessages(state.Event{Funding: &c, TalkNotices: []state.TalkNotice{n}})
		if len(g.messages) != 4 {
			t.Fatalf("kind%d: messages=%d want4", sample.kind, len(g.messages))
		}
		vars, _ := padTalkVars(g.fundingTalkVars(c, 0))
		for i, index := range []int{sample.prompt, sample.first, sample.last} {
			m := g.messages[i+1]
			want := g.legacyTalkLines(index, vars, talkTextWidth)
			if strings.Join(m.lines, "\n") != strings.Join(want, "\n") || m.scene != 1 {
				t.Fatalf("kind%d message%d does not match native TALK%d", sample.kind, i, index)
			}
			if i == 1 {
				if !m.lower || m.portraitPage != g.playerAdvisorPortrait() {
					t.Fatal("advisor reply not in lower box")
				}
			} else if !m.upper || m.portraitPage != 62 {
				t.Fatal("officer reply not in upper box")
			}
		}
		if g.messages[3].then == nil || g.messages[2].underlay == nil || g.messages[3].underlay == nil {
			t.Fatal("acknowledgement lost final continuation or previous reply")
		}
	}
}
