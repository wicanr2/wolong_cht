package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
)

type diplomacyRoll struct{ value, calls int }

func (r *diplomacyRoll) Next() int { r.calls++; return r.value }

func TestCeasefireDiscountUsesSignedTotal(t *testing.T) {
	for _, tc := range []struct{ friendship, amount int }{{20, 14000}, {46, 1000}, {48, 0}, {60, 0}} {
		w := load(t, 0)
		w.Player = 0
		w.Factions[1].Lord, w.Factions[1].Diplomat = 1, 1
		w.Factions[1].Aggression, w.Factions[1].InvasionTarget = 0, 255
		for i := range w.Generals {
			if w.Generals[i].Faction == 1 {
				w.Generals[i].Duty = DutyCorpsLeader
			}
		}
		w.Generals[1] = General{Alive: true, Politics: 8, Faction: 1, Captor: 255}
		w.Friendship[1][0] = diplomacy.Peace(tc.friendship)
		response, amount, ok := w.ceasefireTerms(0, 1)
		if !ok || amount != tc.amount || (response == 0) != (tc.amount == 0) {
			t.Fatalf("交友%d：回應%d、金額%d，原版136F0需先加折扣再鉗零", tc.friendship, response, amount)
		}
	}
}

func TestImportedEventDelayAndNativeSnapshot(t *testing.T) {
	w := cooperationTestWorld(t, 0)
	loaded, err := LoadBlock(w.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	// 原版冷啟動正常讀檔：196/5/6 12h 載入，13h 進入事件。
	if loaded.TakeSnapshot().EventDelay != 1 {
		t.Fatal("DOS 存檔誤用月結等待初值")
	}
	loaded.compactEventQueue()
	if loaded.TakeSnapshot().EventDelay != 7 {
		t.Fatal("月結必須保留原版七小時等待")
	}
	snapshot := loaded.TakeSnapshot()
	snapshot.EventDelay = 4
	if err := loaded.Restore(snapshot); err != nil {
		t.Fatal(err)
	}
	if loaded.TakeSnapshot().EventDelay != 4 {
		t.Fatal("原生快照的相位被 DOS 初值覆蓋")
	}
}

func TestIgnoredDiplomacyKeepsInitialResponse(t *testing.T) {
	for _, tc := range []struct{ friendship, response, amount int }{
		{128 + 90, 0, 0}, {128 + 60, 1, 15000}, {90, 2, 30000},
	} {
		for _, option := range []DiplomacyOption{DiplomacyAcceptFree, DiplomacyOfferFunds, DiplomacyReject} {
			w := cooperationTestWorld(t, 0)
			w.Trust, w.Factions[0].Aggression = 0, 0
			w.Friendship[0][1] = diplomacy.Peace(40)
			w.Friendship[0][2] = diplomacy.Friendship(tc.friendship)
			random := &diplomacyRoll{value: 29}
			w.rng = random
			if !w.beginDiplomacy(&Event{}, DiplomacyChoice{Kind: DiplomacyCooperation, Source: 0, Invader: 1, Target: 2}) {
				t.Fatal("無法建立協力視窗")
			}
			if w.PendingDiplomacy().InitialResponse != tc.response {
				t.Fatal("初始回應未保存")
			}
			w.SetDiplomacyOfferAmount(30000)
			before := random.calls
			result, applied := w.ResolveDiplomacyWithResult(option)
			if result.AdviceAdopted || result.Response != tc.response || result.Amount != tc.amount || applied != (tc.response < 2) {
				t.Fatalf("initial=%d option=%d result=%+v applied=%v", tc.response, option, result, applied)
			}
			if random.calls != before+1 || w.Trust != 0 {
				t.Fatal("多取亂數或誤扣不採納建議的信賴度")
			}
			want := 50000
			if tc.response == 1 {
				want += tc.amount
			}
			if w.Factions[0].Funds != want {
				t.Fatal("未按初始條件付款")
			}
			restored, err := LoadBlock(w.Bytes())
			if err != nil || restored.Factions[0].Funds != want {
				t.Fatal("外交結算未保存", err)
			}
		}
	}
}

func TestIgnoredCeasefireKeepsOriginalRefusal(t *testing.T) {
	w := cooperationTestWorld(t, 0)
	w.Trust = 0
	w.Factions[0].InvasionTarget = 1
	w.Friendship[0][1] = diplomacy.Friendship(20)
	w.Friendship[1][0] = diplomacy.Friendship(20)
	w.rng = &diplomacyRoll{value: 29}
	if !w.beginDiplomacy(&Event{}, DiplomacyChoice{Kind: DiplomacyCeasefire, Source: 1, Target: 0}) {
		t.Fatal("無法建立停戰視窗")
	}
	result, applied := w.ResolveDiplomacyWithResult(DiplomacyAcceptFree)
	if result.AdviceAdopted || result.Response != 2 || applied || !w.Friendship[0][1].AtWar() {
		t.Fatalf("原始拒絕被改成免費停戰：%+v", result)
	}
}

func TestDiplomacyAmountZeroAndCancelledInput(t *testing.T) {
	for _, option := range []DiplomacyOption{DiplomacyOfferFunds, DiplomacyAcceptFree} {
		w := cooperationTestWorld(t, 0)
		w.Factions[0].Aggression = 0
		w.Friendship[0][1], w.Friendship[0][2] = diplomacy.Peace(40), diplomacy.Peace(60)
		w.rng = &diplomacyRoll{value: 29}
		w.beginDiplomacy(&Event{}, DiplomacyChoice{Kind: DiplomacyCooperation, Source: 0, Invader: 1, Target: 2})
		amount := 0
		if option == DiplomacyAcceptFree {
			amount = 30000 // 數值器取消後，舊編輯值不應影響免費建言。
		}
		w.SetDiplomacyOfferAmount(amount)
		result, applied := w.ResolveDiplomacyWithResult(option)
		if !applied || result.Response != 0 || result.AdviceOption != DiplomacyAcceptFree || w.Trust != 255 {
			t.Fatalf("免費／取消輸入的回應錯誤：%+v", result)
		}
		if len(result.TalkNotices) != 1 || result.TalkNotices[0].Faction != 1 {
			t.Fatal("漏掉原版後續宣戰通知")
		}
	}
}
