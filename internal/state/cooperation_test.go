package state

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/strategyai"
)

func cooperationTestWorld(t *testing.T, scenario int) *World {
	t.Helper()
	w := load(t, scenario)
	w.Player, w.Trust = 0, 255
	for i := range w.Generals {
		w.Generals[i].Faction = noFaction
		w.Generals[i].Duty = DutyNone
	}
	for i := 0; i < 3; i++ {
		w.Factions[i].Alive = true
		w.Factions[i].Lord, w.Factions[i].Diplomat = i, noFaction
		w.Factions[i].Funds = 50_000
		w.Factions[i].InvasionTarget = diplomacy.NoTarget
		w.Generals[i] = General{Alive: true, Faction: i, Politics: 7, Captor: noFaction}
	}
	w.events = [eventQueueEntries]QueuedEvent{}
	w.eventCursor, w.eventDelay = 0, 1
	w.rng = rng.NewFixed(7)
	return w
}

// 已由原版 sub_13712 的 12,288 組收據驗證；不是由正式程式產生期望值。
func TestCooperationWarValueAndPeaceBoundaries(t *testing.T) {
	cases := []struct {
		raw, aggression, response, amount int
	}{
		{0, 0, 2, 30000}, {31, 0, 2, 30000}, {90, 0, 2, 30000}, {127, 15, 2, 30000},
		{128 + 39, 0, 2, 25000}, {128 + 40, 0, 1, 25000},
		{128 + 69, 15, 2, 10000}, {128 + 70, 15, 1, 10000},
		{128 + 89, 15, 0, 0}, {128 + 90, 15, 0, 0}, {255, 15, 0, 0},
	}
	for scenario := 0; scenario < 4; scenario++ {
		for _, tc := range cases {
			t.Run(fmt.Sprintf("s%d-raw%d-a%d", scenario, tc.raw, tc.aggression), func(t *testing.T) {
				w := cooperationTestWorld(t, scenario)
				w.Friendship[0][1] = diplomacy.Peace(0)
				w.Friendship[0][2] = diplomacy.Friendship(tc.raw)
				w.Factions[0].Aggression = tc.aggression
				before := w.Bytes()
				response, amount, ok := w.cooperationTerms(0, 1, 2)
				if !ok || response != tc.response || amount != tc.amount {
					t.Fatalf("response=%d amount=%d ok=%v, want %d/%d", response, amount, ok, tc.response, tc.amount)
				}
				if !bytes.Equal(before, w.Bytes()) {
					t.Fatal("條件函式改寫世界資料")
				}
			})
		}
	}
}

func TestCooperationWarValueThroughQueueChoiceAndSave(t *testing.T) {
	w := cooperationTestWorld(t, 0)
	w.Factions[0].Aggression = 0
	w.Factions[1].InvasionTarget = 2
	w.Friendship[1][0] = diplomacy.Peace(50)
	w.Friendship[2][0] = diplomacy.Peace(40)
	w.Friendship[0][1] = diplomacy.Peace(0)
	w.Friendship[0][2] = diplomacy.Friendship(90)
	ordered := []strategyai.Candidate{{Faction: 0, Friendship: w.Friendship[1][0]}}
	if !w.queueCooperationProposal(rng.NewFixed(7), 1, ordered) {
		t.Fatal("反向交友度符合 gate ④b 時沒有產生事件")
	}
	slot := -1
	for i, e := range w.events {
		if e.Code == 2 {
			if e.Param != 0x0201 {
				t.Fatalf("事件的侵攻方／被侵攻方錯誤：%#v", e)
			}
			slot = i
			break
		}
	}
	if slot < 0 {
		t.Fatal("找不到協力佇列")
	}
	w.eventCursor = slot * eventQueueEntrySize
	w.dispatchQueuedEvent(&Event{})
	c := w.PendingDiplomacy()
	if c == nil || c.Source != 0 || c.Invader != 1 || c.Target != 2 || c.InitialAmount != 30_000 {
		t.Fatalf("玩家協力預設金額／角色錯誤：%+v", c)
	}
	clock := w.Clock
	w.Tick(rng.NewFixed(7))
	if w.Clock != clock {
		t.Fatal("外交選單未暫停時間")
	}
	if !w.SetDiplomacyOfferAmount(10_000) || !w.ResolveDiplomacy(DiplomacyOfferFunds) {
		t.Fatal("玩家選擇原版允許的較低金額未完成")
	}
	if w.Factions[0].Funds != 60_000 || w.Factions[2].Funds != 40_000 ||
		!w.Friendship[0][1].AtWar() || !w.Friendship[1][0].AtWar() {
		t.Fatal("被侵攻方付款或合作方宣戰方向錯誤")
	}
	restored, err := LoadBlock(w.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if restored.Factions[0].Funds != 60_000 || restored.Factions[2].Funds != 40_000 ||
		restored.Friendship[0][1] != w.Friendship[0][1] {
		t.Fatal("協力結果未保存")
	}
}

func TestCooperationAIRejectsHighWarValueWithoutPayment(t *testing.T) {
	w := cooperationTestWorld(t, 0)
	w.Factions[1].Diplomat = 0
	w.Factions[1].Aggression = 0
	w.Friendship[1][2] = diplomacy.Peace(0)
	w.Friendship[1][0] = diplomacy.Friendship(90)
	before := w.Bytes()
	if w.applyQueuedCooperationReport(1, 2) {
		t.Fatal("交戰關係被誤判成免費協力")
	}
	if !bytes.Equal(before, w.Bytes()) {
		t.Fatal("AI 拒絕後改寫資金、交戰狀態或佇列")
	}
}
