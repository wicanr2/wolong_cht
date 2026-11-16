package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 渲染測試的已採納結果fixture；規則另由state及真實原版分支驗證。
func adoptedDiplomacyResultForTest(c state.DiplomacyChoice, option state.DiplomacyOption) state.DiplomacyResult {
	r := state.DiplomacyResult{AdviceOption: option, AdviceAdopted: true, Response: int(option), Amount: c.InitialAmount}
	if option == state.DiplomacyOfferFunds {
		r.Amount = c.OfferAmount
		if c.OfferAmount == 0 {
			r.AdviceOption, r.Response = state.DiplomacyAcceptFree, 0
		}
		if c.OfferAmount > c.InitialAmount {
			r.Response = 3
		}
	}
	return r
}

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
			gotResponse, _ := diplomacyTalkResponse(state.DiplomacyResult{Response: tt.wantResponse, Amount: 6000})
			if gotResponse != tt.wantResponse {
				t.Fatalf("response = %d, want %d", gotResponse, tt.wantResponse)
			}
		})
	}

	if got, _ := diplomacyTalkResponse(state.DiplomacyResult{Response: 3, Amount: 7000}); got != 2 {
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
		Kind: state.DiplomacyCooperation, Source: 0, Invader: 2, Target: requester,
		InitialAmount: 6000, OfferAmount: 5000,
	}, state.DiplomacyResult{AdviceOption: state.DiplomacyOfferFunds, AdviceAdopted: true, Response: 1, Amount: 5000})
	if len(g.messages) != 3 {
		t.Fatalf("事件2軍師／君主／結果訊息數 = %d，want 3", len(g.messages))
	}
	choice := strings.Join(g.messages[0].lines, "\n")
	result := strings.Join(g.messages[2].lines, "\n")
	if strings.Contains(choice, "{") || !strings.Contains(choice, "5000") ||
		strings.Contains(result, "{") || !strings.Contains(result, "5000") ||
		!strings.Contains(result, big5(w.LordName(requester))) {
		t.Fatalf("事件 2 選擇／結果 TALK 展開錯誤：choice=%q result=%q", choice, result)
	}
}

// 玩家挑的理由由**軍師**在事件場景的下框說出來（docs/spec/42 §2）。
func TestAdvisorLineUsesLowerBox(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	g := &game{lib: lib, world: w}

	c := state.DiplomacyChoice{
		Kind: state.DiplomacyCeasefire, Source: 1, Invader: 1,
		InitialAmount: 100, OfferAmount: 100,
	}
	g.enqueueDiplomacyTalk(c, state.DiplomacyResult{AdviceOption: state.DiplomacyAcceptFree, AdviceAdopted: true})
	if len(g.messages) < 3 {
		t.Fatalf("外交收尾應有軍師、君主與結果，得到%d則", len(g.messages))
	}

	reason := g.messages[0]
	if !reason.lower {
		t.Error("理由句沒有進下框")
	}
	if reason.scene != 0 {
		t.Errorf("理由句的插圖頁 = %d，事件 2／3 是第 0 頁", reason.scene)
	}
	if want := g.playerAdvisorPortrait(); reason.portraitPage != want {
		t.Errorf("理由句的肖像 = %d，軍師是 %d", reason.portraitPage, want)
	}
	// 結果句仍走一般通知框。
	if g.messages[2].lower || g.messages[2].scene >= 0 {
		t.Error("結果句不該進下框，也不該帶插圖")
	}
	if !g.messages[1].upper || g.messages[1].underlay == nil || reason.underlay == nil {
		t.Error("君主回覆或保留的上下框遺失")
	}
}

// 索引127走自訂軍師的DS:5221；真正無效的索引才退回通報者。
func TestAdvisorPortraitUsesCustomRecordOrFallback(t *testing.T) {
	w := &state.World{Player: 0}
	w.Factions[0].Advisor = state.NoAdvisor
	g := &game{world: w}
	w.AdvisorPortrait = 23
	if got := g.playerAdvisorPortrait(); got != 23 {
		t.Errorf("自訂軍師肖像 = %d，want23", got)
	}
	w.Factions[0].Advisor = 255
	if got := g.playerAdvisorPortrait(); got != defaultPortraitPage {
		t.Errorf("無效軍師肖像 = %d，want %d", got, defaultPortraitPage)
	}
	var nilGame *game
	if got := nilGame.playerAdvisorPortrait(); got != defaultPortraitPage {
		t.Errorf("沒有 world 時肖像 = %d，want %d", got, defaultPortraitPage)
	}
}

func TestBlankCustomAdvisorStillReceivesSummons(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	w.SetCustomAdvisor(23, nil)
	g := &game{lib: lib, world: w}
	g.enqueueTalkNotice(state.TalkNotice{Index: 42, City: -1, Faction: -1, General: -1, Amount: -1})
	if len(g.messages) != 1 || !strings.Contains(strings.Join(g.messages[0].lines, ""), "主公有事") {
		t.Fatal("空白自訂姓名讓原版召見通報消失")
	}
	g.messages = nil
	c := &state.DiplomacyChoice{Kind: state.DiplomacyCooperation, Source: 0, Invader: 12, Target: 2}
	g.enqueueEventMessages(state.Event{Diplomacy: c, TalkNotices: []state.TalkNotice{
		{Index: 42, City: -1, Faction: -1, General: -1, Amount: -1},
	}})
	if len(g.messages) != 2 || g.messages[0].scene >= 0 || !g.messages[1].upper || g.messages[1].scene != 0 {
		t.Fatal("召見與君主提問沒有各自等待確認")
	}
}

func TestIgnoredCooperationTalkUsesActualTermsAndVictim(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player, w.Factions[12].Diplomat = 0, 0
	g := &game{lib: lib, world: w}
	g.enqueueDiplomacyTalk(state.DiplomacyChoice{Kind: state.DiplomacyCooperation, Source: 0, Invader: 12, Target: 2},
		state.DiplomacyResult{AdviceOption: state.DiplomacyReject, AdviceAdopted: false, Response: 1, Amount: 15000})
	if len(g.messages) != 3 {
		t.Fatalf("忽略拒絕建言後的訊息數=%d，want3", len(g.messages))
	}
	text := strings.Join(g.messages[2].lines, "")
	if !strings.Contains(text, "劉備") || !strings.Contains(text, "15000") || strings.Contains(text, "袁術") {
		t.Fatalf("結果沿用了建言或侵攻方：%q", text)
	}
	if g.messages[2].portraitPage != defaultPortraitPage {
		t.Fatal("玩家合作報告誤用了侵攻方外交官肖像")
	}
}

// 原版13262在13C3D入口使用接收方SI，玩家收到來使時仍由147報告。
func TestIncomingCeasefireResultUsesReporterPortrait(t *testing.T) {
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	w.Factions[3].Diplomat = 37
	g := &game{world: w}
	c := state.DiplomacyChoice{Kind: state.DiplomacyCeasefire, Source: 3, Target: 0}
	if got := g.diplomacyResultPortrait(c); got != 147 {
		t.Fatalf("incoming ceasefire portrait=%d want147", got)
	}
	c = state.DiplomacyChoice{Kind: state.DiplomacyCooperation, Source: 0, Invader: 12, Target: 3}
	if got := g.diplomacyResultPortrait(c); got != 147 {
		t.Fatalf("player cooperation portrait=%d want147", got)
	}
}

func TestDiplomacyExcessOfferWarningUsesRawLordVariant(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	lord := w.Factions[0].Lord
	c := state.DiplomacyChoice{Kind: state.DiplomacyCeasefire, Source: 3, Target: 0, InitialAmount: 29000, OfferAmount: 30000}
	for _, v := range []int{0, 1} {
		w.Generals[lord].TalkVariant = v
		g := &game{lib: lib, world: w}
		g.enqueueDiplomacyTalk(c, state.DiplomacyResult{AdviceOption: state.DiplomacyOfferFunds, AdviceAdopted: true, Response: 3, Amount: 30000})
		want := 4
		if len(g.messages) != want {
			t.Fatalf("raw%d messages=%d want%d", v, len(g.messages), want)
		}
		{
			last := g.messages[3]
			vars, ok := w.TalkNoticeVars(state.TalkNotice{Index: 0x1A5, City: -1, Faction: 3, General: lord, Amount: -1}, big5)
			if !ok {
				t.Fatal("warning parameters missing")
			}
			padded, _ := padTalkVars(vars)
			lines, ok := g.talkLines(526+v, padded)
			if !ok || strings.Join(last.lines, "\n") != strings.Join(layoutMessageLines(lines), "\n") || last.portraitPage != w.Generals[lord].Portrait {
				t.Fatal("raw1 excess-offer warning does not match original527/lord portrait")
			}
		}
	}
}
