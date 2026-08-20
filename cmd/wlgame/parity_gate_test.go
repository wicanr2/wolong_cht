package main

import (
	"github.com/wicanr2/wolong_cht/internal/ui/isoview"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// assertTalkModalContract 是事件／M7 的共同呈現層 gate：marker 必須已展開、
// 每一頁最多五列、換行後不超過目前 remake modal 的安全寬度。它不把
// lossy 影片或不同狀態的畫面誤寫成逐像素 parity。
func assertTalkModalContract(t *testing.T, g *game, want int) {
	t.Helper()
	if len(g.messages) != want {
		t.Fatalf("TALK modal 數量 = %d，want %d：%#v", len(g.messages), want, g.messages)
	}
	for i, dialog := range g.messages {
		if len(dialog.lines) == 0 {
			t.Fatalf("第 %d 則 TALK 沒有可見列", i)
		}
		for lineNo, line := range dialog.lines {
			if width := textdraw.StringWidth(line); width > messageContentWidth+textdraw.GlyphW {
				t.Fatalf("第 %d 則 TALK 第 %d 列超出安全寬度：%d px，內容=%q",
					i, lineNo, width, line)
			}
		}
		pages := (len(dialog.lines) + messagePageRows - 1) / messagePageRows
		for page := 0; page < pages; page++ {
			rows, gotPages, ok := messagePage(dialog.lines, page)
			if !ok || gotPages != pages || len(rows) > messagePageRows {
				t.Fatalf("第 %d 則 TALK 第 %d 頁違反五列契約：rows=%#v pages=%d ok=%v",
					i, page, rows, gotPages, ok)
			}
		}
	}
}

// TestEvent2To5TalkBranchParityGate 一次覆蓋事件 2／3 的三選一與事件
// 4／5 的全額、等額、低額、零額、超額、拒絕。索引映射另由 raw branch
// tests 固定；本 gate 只驗證每個已證實 branch 都能進入實際 TALK 展開、
// 分頁與寬度契約。
func TestEvent2To5TalkBranchParityGate(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0

	diplomacyCases := []struct {
		name            string
		kind            state.DiplomacyKind
		option          state.DiplomacyOption
		source, invader int
		offer           int
	}{
		{"event3-free", state.DiplomacyCeasefire, state.DiplomacyAcceptFree, 1, 1, 6000},
		{"event3-funds", state.DiplomacyCeasefire, state.DiplomacyOfferFunds, 1, 1, 5000},
		{"event3-reject", state.DiplomacyCeasefire, state.DiplomacyReject, 1, 1, 6000},
		{"event2-free", state.DiplomacyCooperation, state.DiplomacyAcceptFree, 0, 1, 6000},
		{"event2-funds", state.DiplomacyCooperation, state.DiplomacyOfferFunds, 0, 1, 5000},
		{"event2-reject", state.DiplomacyCooperation, state.DiplomacyReject, 0, 1, 6000},
	}
	for _, tc := range diplomacyCases {
		t.Run(tc.name, func(t *testing.T) {
			g := &game{lib: lib, world: w}
			g.enqueueDiplomacyTalk(state.DiplomacyChoice{
				Kind: tc.kind, Source: tc.source, Invader: tc.invader,
				InitialAmount: 6000, OfferAmount: tc.offer,
			}, tc.option)
			assertTalkModalContract(t, g, 2)
		})
	}

	city := -1
	for i, c := range w.Cities {
		if c.Owner == w.Player {
			city = i
			break
		}
	}
	if city < 0 {
		t.Fatal("找不到事件 4 的玩家據點 fixture")
	}

	fundingCases := []struct {
		name    string
		kind    state.FundingKind
		option  state.FundingOption
		amount  int
		subject int
	}{
		{"event4-full", state.FundingGovernor, state.FundingFullAmount, 6000, city},
		{"event4-equal", state.FundingGovernor, state.FundingSetAmount, 6000, city},
		{"event4-less", state.FundingGovernor, state.FundingSetAmount, 5000, city},
		{"event4-zero", state.FundingGovernor, state.FundingSetAmount, 0, city},
		{"event4-more", state.FundingGovernor, state.FundingSetAmount, 7000, city},
		{"event4-reject", state.FundingGovernor, state.FundingReject, 6000, city},
		{"event5-full", state.FundingDiplomat, state.FundingFullAmount, 6000, 1},
		{"event5-equal", state.FundingDiplomat, state.FundingSetAmount, 6000, 1},
		{"event5-less", state.FundingDiplomat, state.FundingSetAmount, 5000, 1},
		{"event5-zero", state.FundingDiplomat, state.FundingSetAmount, 0, 1},
		{"event5-more", state.FundingDiplomat, state.FundingSetAmount, 7000, 1},
		{"event5-reject", state.FundingDiplomat, state.FundingReject, 6000, 1},
	}
	for _, tc := range fundingCases {
		t.Run(tc.name, func(t *testing.T) {
			g := &game{lib: lib, world: w}
			g.enqueueFundingTalk(state.FundingChoice{
				Kind: tc.kind, Subject: tc.subject,
				RequestedAmount: 6000, OfferAmount: tc.amount,
			}, tc.option)
			assertTalkModalContract(t, g, 2)
		})
	}
}

type correctionManifest struct {
	Corrections []struct {
		ID  int     `json:"id"`
		Fix *string `json:"fix"`
	} `json:"corrections"`
}

// TestM7CorrectedTalkLayoutGate 逐筆走過所有已定案修正；Python selftest
// 負責 byte-level 產出／round-trip，本測試則讓 runtime 的 Big5 marker 展開、
// 實測字寬、五列分頁再走一次，避免校訂表只在工具層正確。
func TestM7CorrectedTalkLayoutGate(t *testing.T) {
	manifestBytes, err := os.ReadFile("../../translations/corrections.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest correctionManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatal(err)
	}
	lib, err := library.LoadWithOptions("../../workplace/orig/dosv", library.LoadOptions{
		TalkJSON: "../../translations/talk-dosv-corrected.json",
	})
	if err != nil {
		t.Fatal(err)
	}
	g := &game{lib: lib}
	vars := map[byte]string{
		'1': "武將", '2': "據點", '3': "君主", '4': "軍師",
		'5': "目標", '6': "", '7': "1234",
	}
	fixed := 0
	for _, correction := range manifest.Corrections {
		if correction.Fix == nil {
			continue
		}
		fixed++
		lines, ok := g.talkLines(correction.ID, vars)
		if !ok || len(lines) == 0 {
			t.Fatalf("M7 #%-3d runtime marker 展開失敗：%#v", correction.ID, correction.Fix)
		}
		wrapped := layoutMessageLines(lines)
		if len(wrapped) == 0 || strings.Contains(strings.Join(wrapped, "\n"), "{") {
			t.Fatalf("M7 #%-3d 仍留下未展開 marker：%#v", correction.ID, wrapped)
		}
		for lineNo, line := range wrapped {
			if width := textdraw.StringWidth(line); width > messageContentWidth+textdraw.GlyphW {
				t.Fatalf("M7 #%-3d 第 %d 列超寬：%d px，內容=%q",
					correction.ID, lineNo, width, line)
			}
		}
		pages := (len(wrapped) + messagePageRows - 1) / messagePageRows
		for page := 0; page < pages; page++ {
			rows, gotPages, ok := messagePage(wrapped, page)
			if !ok || gotPages != pages || len(rows) > messagePageRows {
				t.Fatalf("M7 #%-3d 第 %d 頁不符合五列契約：%#v", correction.ID, page, rows)
			}
		}
	}
	if fixed != 60 {
		t.Fatalf("M7 已定案修正數 = %d，want 60", fixed)
	}
}

// TestProjectileParityGate 將 raw BATTLE.SCH 圖號與規則層的三個可見
// 投射物狀態綁在同一個命名 gate：普通兩方向、特殊第一幀、特殊第二幀。
// 速度／高度／碰撞的逐步規則由 internal/rules/tactical 的 raw tests
// 分別驗證，這裡只驗證呈現層不會丟掉 frame bit。
func TestProjectileParityGate(t *testing.T) {
	cases := []struct {
		name string
		view tactical.ProjectileView
		want int
	}{
		{"normal-horizontal", tactical.ProjectileView{Direction: tactical.West}, 0x210},
		{"normal-vertical", tactical.ProjectileView{Direction: tactical.North}, 0x211},
		{"special-frame-0", tactical.ProjectileView{Special: true, Direction: tactical.South | 0x80}, 0x214},
		{"special-frame-1", tactical.ProjectileView{Special: true, SpecialFrame: 1, Direction: tactical.South | 0x80}, 0x215},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isoview.ProjectileSourceIndex(tc.view); got != tc.want {
				t.Fatalf("raw projectile image = %#x，want %#x", got, tc.want)
			}
		})
	}
}

// TestEvent9ShortFixtureGate 是事件 9 的最小正常通知路徑：state 已先
// 證實釋放結果；這裡驗證只有玩家勢力收到 #37，#409 空槽不會生成空白 modal。
func TestEvent9ShortFixtureGate(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	general := -1
	for i, candidate := range w.Generals {
		if candidate.Alive && i != w.Factions[w.Player].Lord {
			general = i
			break
		}
	}
	if general < 0 || general >= len(w.Generals) || !w.Generals[general].Alive {
		t.Fatal("找不到事件 9 釋放武將 fixture")
	}

	for _, tc := range []struct {
		name       string
		faction    int
		wantNotice int
	}{
		{"player-faction", w.Player, 1},
		{"other-faction", (w.Player + 1) % len(w.Factions), 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w.Generals[general].Faction = tc.faction
			g := &game{lib: lib, world: w}
			g.enqueueEventMessages(state.Event{ReleasedGenerals: []int{general}})
			if len(g.messages) != tc.wantNotice {
				t.Fatalf("事件 9 #37 modal = %d，want %d：%#v",
					len(g.messages), tc.wantNotice, g.messages)
			}
			if tc.wantNotice != 0 {
				assertTalkModalContract(t, g, 1)
			}
		})
	}

	g := &game{lib: lib}
	g.enqueueTalk(0x199, map[byte]string{
		'1': "武將", '2': "據點", '3': "君主", '4': "軍師", '6': "", '7': "0",
	})
	if len(g.messages) != 0 {
		t.Fatalf("事件 9 #409 空槽不應建立空白 modal：%#v", g.messages)
	}
}
