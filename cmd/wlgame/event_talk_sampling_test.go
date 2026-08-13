package main

import (
	"reflect"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// assertTalkSamplePage 把實際 runtime 的 TALK 展開結果，逐字對回指定的
// DOS/V TALK 槽位，再驗證原始硬行、實測像素寬度與五列分頁。這不是把
// 自己產生的文字拿來測自己，而是固定「分支 → raw index → corrected TALK」
// 的完整鏈。
func assertTalkSamplePage(t *testing.T, lib *library.Library, w *state.World,
	index int, vars map[byte]string, label string) {
	t.Helper()
	g := &game{lib: lib, world: w}
	g.enqueueTalk(index, vars)
	if len(g.messages) != 1 {
		t.Fatalf("%s TALK #%d modal 數 = %d，want 1", label, index, len(g.messages))
	}
	raw, ok := g.talkLines(index, vars)
	if !ok || len(raw) == 0 {
		t.Fatalf("%s TALK #%d raw 展開失敗", label, index)
	}
	want := textdraw.WrapLines(raw, messageContentWidth)
	if !reflect.DeepEqual(g.messages[0].lines, want) {
		t.Fatalf("%s TALK #%d 頁面文字不一致：got=%#v want=%#v",
			label, index, g.messages[0].lines, want)
	}
	pages := (len(want) + messagePageRows - 1) / messagePageRows
	if pages < 1 {
		t.Fatalf("%s TALK #%d 沒有可見頁面", label, index)
	}
	for lineNo, line := range want {
		if width := textdraw.StringWidth(line); width > messageContentWidth+textdraw.GlyphW {
			t.Fatalf("%s TALK #%d 第 %d 列超寬：%d px，內容=%q",
				label, index, lineNo, width, line)
		}
	}
	t.Logf("%s TALK #%d：%d 原始硬行 → %d 實際行／%d 頁",
		label, index, len(raw), len(want), pages)
}

func assertTalkPairSample(t *testing.T, lib *library.Library, w *state.World,
	indices []int, vars []map[byte]string, label string, enqueue func(*game)) {
	t.Helper()
	g := &game{lib: lib, world: w}
	enqueue(g)
	if len(g.messages) != len(indices) {
		t.Fatalf("%s modal 數 = %d，want %d", label, len(g.messages), len(indices))
	}
	for i, index := range indices {
		raw, ok := g.talkLines(index, vars[i])
		if !ok || len(raw) == 0 {
			t.Fatalf("%s TALK #%d raw 展開失敗", label, index)
		}
		want := textdraw.WrapLines(raw, messageContentWidth)
		if !reflect.DeepEqual(g.messages[i].lines, want) {
			t.Fatalf("%s 第 %d 則 TALK #%d 不一致：got=%#v want=%#v",
				label, i, index, g.messages[i].lines, want)
		}
		for lineNo, line := range want {
			if width := textdraw.StringWidth(line); width > messageContentWidth+textdraw.GlyphW {
				t.Fatalf("%s TALK #%d 第 %d 列超寬：%d px，內容=%q",
					label, index, lineNo, width, line)
			}
		}
		pages := (len(want) + messagePageRows - 1) / messagePageRows
		t.Logf("%s TALK #%d：%d 原始硬行 → %d 實際行／%d 頁",
			label, index, len(raw), len(want), pages)
	}
}

// TestEvent2To5FullTalkPageSampling 固定事件 2／3 的前置 prompt、三選一、
// 成功／金額／拒絕／超額回應，以及事件 4／5 的 prompt、三選一與六種撥款
// 結果。共 36 個 raw TALK 頁面／18 組雙頁回應，確保短 parity gate 之外，
// 每一個已接入分支都真的讀到校訂 TALK、保留硬行並能分頁。
func TestEvent2To5FullTalkPageSampling(t *testing.T) {
	lib, err := library.LoadWithOptions("../../workplace/orig/dosv", library.LoadOptions{
		TalkJSON: "../../translations/talk-dosv-corrected.json",
	})
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0

	ceasefire := state.DiplomacyChoice{
		Kind: state.DiplomacyCeasefire, Source: 1, Target: 0,
		InitialAmount: 6000, OfferAmount: 5000,
	}
	cooperation := state.DiplomacyChoice{
		Kind: state.DiplomacyCooperation, Source: 0, Invader: 1, Target: 2,
		InitialAmount: 6000, OfferAmount: 5000,
	}
	g := &game{lib: lib, world: w}
	for _, tc := range []struct {
		label string
		c     state.DiplomacyChoice
		base  int
		vars  map[byte]string
	}{
		{"event3", ceasefire, 0x168, g.diplomacyTalkVars(ceasefire, -1)},
		{"event2", cooperation, 0x175, g.diplomacyTalkVars(cooperation, -1)},
	} {
		for variant := 0; variant < 3; variant++ {
			assertTalkSamplePage(t, lib, w, tc.base+variant, tc.vars,
				tc.label+" prompt variant")
		}
		assertTalkSamplePage(t, lib, w, tc.base+3, tc.vars,
			tc.label+" choice")
	}

	diplomacyCases := []struct {
		label   string
		c       state.DiplomacyChoice
		option  state.DiplomacyOption
		indices []int
	}{
		{"event3-free", ceasefire, state.DiplomacyAcceptFree, []int{364, 43}},
		{"event3-funds", ceasefire, state.DiplomacyOfferFunds, []int{365, 44}},
		{"event3-reject", ceasefire, state.DiplomacyReject, []int{366, 45}},
		{"event3-over", state.DiplomacyChoice{
			Kind: state.DiplomacyCeasefire, Source: 1, Target: 0,
			InitialAmount: 6000, OfferAmount: 7000,
		}, state.DiplomacyOfferFunds, []int{365, 45}},
		{"event2-free", cooperation, state.DiplomacyAcceptFree, []int{377, 47}},
		{"event2-funds", cooperation, state.DiplomacyOfferFunds, []int{378, 48}},
		{"event2-reject", cooperation, state.DiplomacyReject, []int{379, 49}},
		{"event2-over", state.DiplomacyChoice{
			Kind: state.DiplomacyCooperation, Source: 0, Invader: 1, Target: 2,
			InitialAmount: 6000, OfferAmount: 7000,
		}, state.DiplomacyOfferFunds, []int{378, 49}},
	}
	for _, tc := range diplomacyCases {
		t.Run(tc.label, func(t *testing.T) {
			g := &game{lib: lib, world: w}
			choiceAmount := -1
			if tc.option == state.DiplomacyOfferFunds {
				choiceAmount = tc.c.OfferAmount
			}
			response, resultAmount := diplomacyTalkResponse(tc.c, tc.option)
			_ = response
			vars := []map[byte]string{
				g.diplomacyTalkVars(tc.c, choiceAmount),
				g.diplomacyTalkVars(tc.c, resultAmount),
			}
			assertTalkPairSample(t, lib, w, tc.indices, vars, tc.label,
				func(g *game) { g.enqueueDiplomacyTalk(tc.c, tc.option) })
		})
	}

	fundingCases := []struct {
		label   string
		c       state.FundingChoice
		option  state.FundingOption
		indices []int
	}{
		{"event4-full", state.FundingChoice{Kind: state.FundingGovernor, Subject: 0, RequestedAmount: 6000, OfferAmount: 6000}, state.FundingFullAmount, []int{284, 288}},
		{"event4-equal", state.FundingChoice{Kind: state.FundingGovernor, Subject: 0, RequestedAmount: 6000, OfferAmount: 6000}, state.FundingSetAmount, []int{284, 293}},
		{"event4-less", state.FundingChoice{Kind: state.FundingGovernor, Subject: 0, RequestedAmount: 6000, OfferAmount: 5000}, state.FundingSetAmount, []int{285, 293}},
		{"event4-zero", state.FundingChoice{Kind: state.FundingGovernor, Subject: 0, RequestedAmount: 6000, OfferAmount: 0}, state.FundingSetAmount, []int{286, 293}},
		{"event4-more", state.FundingChoice{Kind: state.FundingGovernor, Subject: 0, RequestedAmount: 6000, OfferAmount: 7000}, state.FundingSetAmount, []int{287, 293}},
		{"event4-reject", state.FundingChoice{Kind: state.FundingGovernor, Subject: 0, RequestedAmount: 6000, OfferAmount: 6000}, state.FundingReject, []int{286, 298}},
		{"event5-full", state.FundingChoice{Kind: state.FundingDiplomat, Subject: 1, RequestedAmount: 6000, OfferAmount: 6000}, state.FundingFullAmount, []int{325, 329}},
		{"event5-equal", state.FundingChoice{Kind: state.FundingDiplomat, Subject: 1, RequestedAmount: 6000, OfferAmount: 6000}, state.FundingSetAmount, []int{325, 334}},
		{"event5-less", state.FundingChoice{Kind: state.FundingDiplomat, Subject: 1, RequestedAmount: 6000, OfferAmount: 5000}, state.FundingSetAmount, []int{326, 334}},
		{"event5-zero", state.FundingChoice{Kind: state.FundingDiplomat, Subject: 1, RequestedAmount: 6000, OfferAmount: 0}, state.FundingSetAmount, []int{327, 334}},
		{"event5-more", state.FundingChoice{Kind: state.FundingDiplomat, Subject: 1, RequestedAmount: 6000, OfferAmount: 7000}, state.FundingSetAmount, []int{328, 334}},
		{"event5-reject", state.FundingChoice{Kind: state.FundingDiplomat, Subject: 1, RequestedAmount: 6000, OfferAmount: 6000}, state.FundingReject, []int{327, 339}},
	}
	for _, tc := range fundingCases {
		t.Run(tc.label, func(t *testing.T) {
			g := &game{lib: lib, world: w}
			amount := tc.c.RequestedAmount
			if tc.option == state.FundingSetAmount {
				amount = tc.c.OfferAmount
			}
			vars := g.fundingTalkVars(tc.c, amount)
			assertTalkPairSample(t, lib, w, tc.indices,
				[]map[byte]string{vars, vars}, tc.label,
				func(g *game) { g.enqueueFundingTalk(tc.c, tc.option) })
		})
	}
}
