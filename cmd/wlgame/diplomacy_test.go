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
	g.enqueueDiplomacyTalk(c, state.DiplomacyAcceptFree)
	if len(g.messages) < 2 {
		t.Fatalf("外交收尾應該有理由句 ＋ 結果句，得到 %d 則", len(g.messages))
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
	if g.messages[1].lower || g.messages[1].scene >= 0 {
		t.Error("結果句不該進下框，也不該帶插圖")
	}
}

// 沒有軍師時退回一般通知的肖像，不要畫錯人（原版 +0x02 寫 0x7F）。
func TestAdvisorPortraitFallsBackWhenNoAdvisor(t *testing.T) {
	w := &state.World{Player: 0}
	w.Factions[0].Advisor = state.NoAdvisor
	g := &game{world: w}
	if got := g.playerAdvisorPortrait(); got != defaultPortraitPage {
		t.Errorf("沒有軍師時肖像 = %d，want %d", got, defaultPortraitPage)
	}
	var nilGame *game
	if got := nilGame.playerAdvisorPortrait(); got != defaultPortraitPage {
		t.Errorf("沒有 world 時肖像 = %d，want %d", got, defaultPortraitPage)
	}
}

// 結果句由**派駐的外交官**回報，不是一般通知的那張臉（原版 sub_13C3D）。
func TestDiplomacyResultUsesEnvoyPortrait(t *testing.T) {
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w.Player = 0
	g := &game{world: w}
	other := 1
	c := state.DiplomacyChoice{Kind: state.DiplomacyCeasefire, Source: other, Invader: other}

	// 沒有外交官（開局全 0xFF）→ 退回預設，不畫錯人。
	w.Factions[other].Diplomat = 0xFF
	if got := g.diplomacyResultPortrait(c); got != defaultPortraitPage {
		t.Errorf("沒有外交官時 = %d，want %d", got, defaultPortraitPage)
	}

	// 有外交官 → 用那名武將的頭像。
	envoy := -1
	for i := range w.Generals {
		if w.Generals[i].Alive {
			envoy = i
			break
		}
	}
	if envoy < 0 {
		t.Skip("劇本裡沒有活著的武將")
	}
	w.Factions[other].Diplomat = envoy
	if got, want := g.diplomacyResultPortrait(c), w.Generals[envoy].Portrait; got != want {
		t.Errorf("外交官的頭像 = %d，want %d", got, want)
	}

	// 對方就是玩家自己 → 預設。
	c.Source, c.Invader = w.Player, w.Player
	if got := g.diplomacyResultPortrait(c); got != defaultPortraitPage {
		t.Errorf("對方是自己時 = %d，want %d", got, defaultPortraitPage)
	}
}
