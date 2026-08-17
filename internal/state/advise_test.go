package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/combat"
)

// 君主親自出陣的那一支**不是委任**——`sub_16E8F` 設了委任位元，
// `sub_1699E` 緊接著把它清掉（docs/spec/49 §3）。
func TestAdviseSortieFormsUndelegatedCorps(t *testing.T) {
	w := load(t, 0)
	w.Player = w.AliveFactions()[0]
	f := &w.Factions[w.Player]
	lord := f.Lord
	if lord < 0 || lord >= numCorps {
		t.Skip("這個勢力沒有君主")
	}
	f.Aggression = 15 // 資金那道閘的門檻降到 0
	f.Funds = 0
	f.Reserves = [3]int{300, 300, 300} // 總和 900 ≥ 600

	if !w.AdviseSortieAccepted() {
		t.Fatal("兩道閘都過了卻不同意")
	}
	if !w.AdviseSortie() {
		t.Fatal("同意了卻編不出軍團")
	}
	c := &w.Corps[lord]
	if !c.Alive {
		t.Fatal("軍團沒編出來")
	}
	if c.Delegated {
		t.Error("君主親自出陣的軍團不該是委任的")
	}
	if !w.Generals[lord].Posted {
		t.Error("君主沒被標成出陣中")
	}
	// 君主已經帶著軍團就不能再出一次。
	if w.AdviseSortieAccepted() {
		t.Error("君主已經在陣上還能再請一次")
	}
}

// 兵不夠就擋下來，而且**不留下半支軍團**。
func TestAdviseSortieBlockedLeavesNothing(t *testing.T) {
	w := load(t, 0)
	w.Player = w.AliveFactions()[0]
	f := &w.Factions[w.Player]
	f.Aggression, f.Funds = 15, 0
	f.Reserves = [3]int{100, 100, 100} // 總和 300 < 600

	if w.AdviseSortieAccepted() {
		t.Fatal("兵不夠卻同意出陣")
	}
	if w.AdviseSortie() {
		t.Fatal("被擋下卻還是編了軍團")
	}
	if lord := f.Lord; lord >= 0 && lord < numCorps && w.Corps[lord].Alive {
		t.Error("擋下之後不該留著軍團")
	}
}

// 遷都成功之後首都要換，而且掛著舊首都的軍團要跟著改掛（`sub_14502`）。
func TestAdviseRelocateMovesCapitalAndCorps(t *testing.T) {
	w := load(t, 0)
	w.Player = w.AliveFactions()[0]
	f := &w.Factions[w.Player]
	from := w.clampCity(f.Capital)

	// 找一個自己的據點，把它調成「同級或更大、生產力更高」。
	to := -1
	for i := range w.Cities {
		if i != from && w.Cities[i].Owner == w.Player {
			to = i
			break
		}
	}
	if to < 0 {
		t.Skip("這個勢力只有一個據點")
	}
	w.Cities[to].Kind = w.Cities[from].Kind
	w.Cities[to].Production = w.Cities[from].Production + 1

	i := formOne(t, w, w.Player)
	w.Corps[i].Ordered = from

	if !w.AdviseRelocateAccepted(to) {
		t.Fatal("條件都滿足卻不同意遷都")
	}
	if !w.AdviseRelocate(to) {
		t.Fatal("同意了卻沒搬")
	}
	if f.Capital != to {
		t.Errorf("首都 = %d，要 %d", f.Capital, to)
	}
	if got := w.Corps[i].Ordered; got != to {
		t.Errorf("掛舊首都的軍團沒改掛：%d", got)
	}
}

// ⚠ 只看生產力會誤判：更高的生產力配上更小的城，君主還是拒絕。
func TestAdviseRelocateRefusesSmallerCity(t *testing.T) {
	w := load(t, 0)
	w.Player = w.AliveFactions()[0]
	f := &w.Factions[w.Player]
	from := w.clampCity(f.Capital)

	to := -1
	for i := range w.Cities {
		if i != from && w.Cities[i].Owner == w.Player {
			to = i
			break
		}
	}
	if to < 0 {
		t.Skip("這個勢力只有一個據點")
	}
	// 生產力更高（會通過第二個條件），但類型編號更大 ＝ 城更小。
	w.Cities[from].Kind = 0
	w.Cities[to].Kind = 1
	w.Cities[to].Production = w.Cities[from].Production + 100

	if w.AdviseRelocateAccepted(to) {
		t.Error("往更小的城搬不該被接受——只看生產力就會漏掉這一條")
	}
	if f.Capital != from {
		t.Error("被拒絕卻搬了")
	}
}

// 別人的據點與現任首都本身都不能當目標。
func TestAdviseRelocateRejectsBadTargets(t *testing.T) {
	w := load(t, 0)
	w.Player = w.AliveFactions()[0]
	from := w.clampCity(w.Factions[w.Player].Capital)
	if w.AdviseRelocateAccepted(from) {
		t.Error("搬到現在的首都不該被接受")
	}
	for i := range w.Cities {
		if w.Cities[i].Owner != w.Player {
			if w.AdviseRelocateAccepted(i) {
				t.Error("別人的據點不該能當首都")
			}
			break
		}
	}
	if w.AdviseRelocateAccepted(-1) || w.AdviseRelocateAccepted(len(w.Cities)) {
		t.Error("界外的編號要擋掉")
	}
}

// ---------------------------------------------------------------------------
// 帳務：軍費 vs 預備兵維持費（docs/spec/50）
// ---------------------------------------------------------------------------

// ⭐ 軍費當場從資金扣（`sub_12600` → `sub_1562B`），**不進本月支出**。
func TestCorpsUpkeepChargesFundsNotExpense(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	i := formOne(t, w, f)
	men := w.Corps[i].Men
	if men == 0 {
		t.Skip("編不出兵")
	}
	w.Factions[f].Funds = 100000
	w.Factions[f].Expense = 0
	w.corpsCursor = i

	// 直接跑軍團那一輪的「一時」結算。走完整的 Tick 會把補兵、
	// 預備兵維持費一起拉進來，驗不到想驗的那一筆。
	w.tickCorps(upkeepHour, &testRand{s: 7})

	want := 100000 - combat.Upkeep(men, false)
	if got := w.Factions[f].Funds; got != want {
		t.Errorf("資金 = %d，要 %d（兵力 %d ÷ 32 ＋ 1）", got, want, men)
	}
	if got := w.Factions[f].Expense; got != 0 {
		t.Errorf("軍費記進了本月支出：%d", got)
	}
}

// 預備兵維持費仍走「本月支出」（`sub_13E65` → `sub_15673`）。
func TestReserveUpkeepGoesToExpense(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[0]
	w.Factions[f].Expense = 0
	w.Factions[f].Reserves = [3]int{3200, 0, 0} // ÷32 ＝ 100／次
	before := w.Factions[f].Funds

	ev := &Event{}
	w.hourFaction = f
	w.hourly(ev, &testRand{s: 3})

	if got := w.Factions[f].Expense; got != 100 {
		t.Errorf("本月支出 = %d，要 100（3200 ÷ 32）", got)
	}
	if w.Factions[f].Funds != before {
		t.Errorf("預備兵維持費不該當場扣資金：%d → %d", before, w.Factions[f].Funds)
	}
}
