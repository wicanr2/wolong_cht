package state

import (
	"os"
	"reflect"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

func TestOutcomeTrustZeroLatchesAtMutationBoundary(t *testing.T) {
	w := &World{Trust: 1}
	before := w.Clock
	if got := w.AdjustTrust(-1); got != 0 {
		t.Fatalf("Trust = %d，want 0", got)
	}
	if got := w.Outcome(); got != DefeatTrustZero {
		t.Fatalf("Outcome = %v，want DefeatTrustZero", got)
	}
	if selector, ok := w.OutcomeMessageSelector(); !ok || selector != 0x019E {
		t.Fatalf("selector = %#x, %v，want 0x019E, true", selector, ok)
	}
	// 已經離開主循環後，Tick 不可再改時鐘或觸發 AI／據點副作用。
	w.Tick(rng.NewFixed(1))
	if !reflect.DeepEqual(w.Clock, before) {
		t.Fatalf("outcome 後 Clock 改變：%+v → %+v", before, w.Clock)
	}

	// 已為零但沒有 1→0 transition 時不應自行建立敗北結果。
	zero := &World{Trust: 0}
	zero.AdjustTrust(-1)
	if got := zero.Outcome(); got != InProgress {
		t.Fatalf("已為 0 的 Trust 不應觸發 outcome：%v", got)
	}
}

func outcomeCaptureWorld(player, old int, withAlternative bool) *World {
	w := &World{Player: player}
	for i := range w.Cities {
		w.Cities[i].Owner = combat.NeutralFaction
		w.Cities[i].OwnerRecorded = combat.NeutralFaction
	}
	w.Factions[old] = Faction{Alive: true, Capital: 0, Cities: 1}
	w.Factions[player].Alive = true
	w.Cities[0].Owner = old
	w.Cities[0].OwnerRecorded = old
	if withAlternative {
		w.Cities[1].Owner = old
		w.Cities[1].OwnerRecorded = old
		w.Factions[old].Cities = 2
	}
	attacker := player
	if attacker == old {
		attacker = (old + 1) % numFactions
	}
	w.Factions[attacker].Alive = true
	w.Corps[0] = Corps{Alive: true, Faction: attacker, Node: 0}
	return w
}

func TestOutcomeCaptureLastCapitalDefeatsPlayer(t *testing.T) {
	w := outcomeCaptureWorld(1, 1, false)
	ev := &CorpsEvent{Captured: -1}
	w.capture(0, ev, &testRand{s: 1})
	if w.Factions[1].Alive || w.Factions[1].Capital != noCity {
		t.Fatalf("玩家勢力未按最後首都邊界清除：%+v", w.Factions[1])
	}
	if got := w.Outcome(); got != DefeatFactionEliminated {
		t.Fatalf("Outcome = %v，want DefeatFactionEliminated", got)
	}
}

func TestOutcomeCaptureWithAlternativeCapitalDoesNotDefeat(t *testing.T) {
	w := outcomeCaptureWorld(1, 1, true)
	ev := &CorpsEvent{Captured: -1}
	w.capture(0, ev, &testRand{s: 1})
	if !w.Factions[1].Alive || w.Factions[1].Capital != 1 {
		t.Fatalf("仍有替代據點卻未遷都：%+v", w.Factions[1])
	}
	if got := w.Outcome(); got != InProgress {
		t.Fatalf("仍有替代據點不應敗北：%v", got)
	}
}

func TestOutcomeNonPlayerEliminationDoesNotSetPlayerOutcome(t *testing.T) {
	w := outcomeCaptureWorld(1, 2, false)
	ev := &CorpsEvent{Captured: -1}
	w.capture(0, ev, &testRand{s: 1})
	if w.Factions[2].Alive || w.Factions[2].Capital != noCity {
		t.Fatalf("非玩家勢力未清除：%+v", w.Factions[2])
	}
	if got := w.Outcome(); got != InProgress {
		t.Fatalf("非玩家消滅不應設定玩家 outcome：%v", got)
	}
}

func TestOutcomeLatchIsNotOverwritten(t *testing.T) {
	w := &World{Trust: 1}
	w.AdjustTrust(-1)
	if !w.latchOutcome(DefeatFactionEliminated) {
		// 第二次 latch 必須拒絕；這個分支本身是預期結果。
	} else {
		t.Fatal("已 latch 的 outcome 不應允許第二個原因覆蓋")
	}
	if got := w.Outcome(); got != DefeatTrustZero {
		t.Fatalf("Outcome 被覆蓋成 %v", got)
	}
}

// 存活勢力數要從劇本區塊 +0x3A 載入，四個劇本的值是 22／11／6／4。
//
// 這一條直接讀原版檔案——**它是「+0x3A 是存活勢力數」的證據本身**
// （`docs/re/59` §3），不是把實作反過來寫成期望值。
func TestLivingFactionsComesFromScenarioBlock(t *testing.T) {
	raw, err := os.ReadFile("../../workplace/orig/dosv/SINARIO.DAT")
	if err != nil {
		t.Skip("找不到原版 SINARIO.DAT，跳過")
	}
	want := []int{22, 11, 6, 4}
	for i, n := range want {
		if got := int(raw[i*blockSize+livingFactionsOffset]); got != n {
			t.Errorf("劇本 %d 的 +0x3A ＝ %d，該劇本有 %d 個勢力", i, got, n)
		}
		w := loadBlock(raw[i*blockSize : (i+1)*blockSize])
		if w.LivingFactions != n {
			t.Errorf("劇本 %d 載入後 LivingFactions ＝ %d，want %d", i, w.LivingFactions, n)
		}
	}
}

// 滅到剩一個勢力就是結局。
func TestEliminatingDownToOneFactionIsVictory(t *testing.T) {
	w := &World{Player: 0, LivingFactions: 3}
	for i := range w.Factions {
		w.Factions[i].Alive = i < 3
	}
	w.eliminateFaction(2, noFaction)
	if got := w.Outcome(); got != InProgress {
		t.Fatalf("還剩兩個勢力就判 %v", got)
	}
	w.eliminateFaction(1, noFaction)
	if got := w.Outcome(); got != Victory {
		t.Fatalf("剩一個勢力時 Outcome ＝ %v，want Victory", got)
	}
}

// ⚠ 玩家自己滅亡要走敗北，不能因為「剩一個」變成結局。
//
// 原版的順序（先判玩家、再減計數器）就是為了這件事（`docs/re/59` §4）。
func TestPlayerEliminationBeatsVictory(t *testing.T) {
	w := &World{Player: 0, LivingFactions: 2}
	for i := range w.Factions {
		w.Factions[i].Alive = i < 2
	}
	w.eliminateFaction(0, noFaction)
	if got := w.Outcome(); got != DefeatFactionEliminated {
		t.Fatalf("玩家勢力滅亡時 Outcome ＝ %v，want DefeatFactionEliminated", got)
	}
}

// 存活勢力數要能 round-trip 回同一個 byte。
func TestLivingFactionsRoundTrips(t *testing.T) {
	raw, err := os.ReadFile("../../workplace/orig/dosv/SINARIO.DAT")
	if err != nil {
		t.Skip("找不到原版 SINARIO.DAT，跳過")
	}
	block := append([]byte(nil), raw[:blockSize]...)
	w := loadBlock(block)
	out := w.Bytes()
	if out[livingFactionsOffset] != block[livingFactionsOffset] {
		t.Errorf("寫回後 +0x3A ＝ %d，原本是 %d",
			out[livingFactionsOffset], block[livingFactionsOffset])
	}
}

// 結局的第二則是**組編號** `0x197`，說話者是玩家所仕勢力的君主
// （`sub_11CD0` 取君主記錄的說話型與肖像，docs/re/59 §2）。
func TestVictoryLordTalkIndex(t *testing.T) {
	w := &World{Player: 3}
	w.Factions[3].Lord = 7
	if _, _, ok := w.VictoryLordTalkIndex(); ok {
		t.Fatal("還沒結局就給了結局訊息")
	}
	w.outcome = Victory
	index, lord, ok := w.VictoryLordTalkIndex()
	if !ok {
		t.Fatal("結局了卻拿不到君主那一句")
	}
	if index != VictoryLordTalkBase {
		t.Fatalf("組編號 = %#x，預期 %#x", index, VictoryLordTalkBase)
	}
	if lord != 7 {
		t.Fatalf("說話者 = %d，預期君主 7", lord)
	}
	// 沒有君主（例如剛滅亡）就不硬給。
	w.Factions[3].Lord = noFaction
	if _, _, ok := w.VictoryLordTalkIndex(); ok {
		t.Fatal("沒有君主卻給了訊息")
	}
}
