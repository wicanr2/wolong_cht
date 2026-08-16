package persuasion

import "testing"

// 說明書 3.9 明說「国力の基準は拠点数からチェックされ」——
// 國力比的是**據點數**，不是兵力也不是資金。
func TestPowerUsesCityCount(t *testing.T) {
	s := Situation{OurCities: 14, TheirCities: 8}
	if !s.Applies(Hostility, WeAreStronger) {
		t.Error("14 個據點對 8 個應該算有利")
	}
	if s.Applies(CeaseFire, EnemyIsStronger) {
		t.Error("同一局勢不該兩邊都成立")
	}
	// 資金再多也不影響國力判定。
	s.OurFunds, s.TheirFunds = 0, 655000
	if !s.Applies(Hostility, WeAreStronger) {
		t.Error("國力只看據點數，資金不該影響")
	}
}

// 疲弊的定義是**資金 < 0**（機器碼用 24 位值的符號位判斷）。
func TestExhaustionIsNegativeFunds(t *testing.T) {
	s := Situation{OurFunds: -1, TheirFunds: 1}
	if !s.Applies(CeaseFire, WeAreExhausted) {
		t.Error("資金 −1 應該算疲弊")
	}
	if s.Applies(Hostility, EnemyExhausted) {
		t.Error("資金 +1 不該算疲弊")
	}
	s.TheirFunds = 0
	if s.Applies(Hostility, EnemyExhausted) {
		t.Error("資金剛好 0 不算疲弊（門檻是「マイナス」）")
	}
}

// ⭐ 好戰等級高的君主，「我國有利」的門檻比較寬鬆
// （說明書：好戦レベルが高い場合は多少拠点数が少なくとも有利とみなします）。
func TestAggressionLoosensPowerCheck(t *testing.T) {
	base := Situation{OurCities: 10, TheirCities: 12}
	lu := base
	lu.Aggression = 15 // 呂布
	liu := base
	liu.Aggression = 4 // 劉備
	if !lu.Applies(Hostility, WeAreStronger) {
		t.Error("好戰 15 的君主在 10 對 12 時應該仍覺得有利")
	}
	if liu.Applies(Hostility, WeAreStronger) {
		t.Error("好戰 4 的君主在 10 對 12 時不該覺得有利")
	}
}

// ⭐ 說明書點名的兩個例子：呂布不太差就能說服，劉備要相當差才行。
func TestFriendshipGateMatchesManualExamples(t *testing.T) {
	// 交友度是**含和平位元**的原始值（0x80 ＝ 和平中）。
	const mildlyBad = peaceBit + 25
	lu := Situation{Aggression: 15, Friendship: mildlyBad}
	liu := Situation{Aggression: 4, Friendship: mildlyBad}
	if !lu.Applies(Hostility, FriendshipBad) {
		t.Error("呂布（好戰 15）在交友值 25 時應該接受「交友關係惡」")
	}
	if liu.Applies(Hostility, FriendshipBad) {
		t.Error("劉備（好戰 4）在交友值 25 時不該接受，要更差才行")
	}
	// 劉備要非常差才接受。
	liu.Friendship = peaceBit + 5
	if !liu.Applies(Hostility, FriendshipBad) {
		t.Error("交友值 5 連劉備都該接受了")
	}
}

// ⭐ 「交友關係惡」與「交友關係良」**不是互補的**，中間有一大段兩者皆偽。
//
// 它們是不同指令的理由，門檻也不同：
//
//	敵對「交友關係惡劣」 交友度 <  0x80 ＋ 好戰 ＋ 15
//	協力「交友關係良好」 交友度 ≥ 0x80 ＋ 好戰 × 4 ＋ 60
//
// 好戰 5 時中間那段是 0xA0–0xD0（32–79）。舊版把兩者寫成互補，
// 那一整段的判定必定有一邊是錯的。
func TestFriendshipGatesAreNotComplementary(t *testing.T) {
	const aggr = 5
	middle := Situation{Aggression: aggr,
		Friendship: peaceBit + 50, AllyFriendship: peaceBit + 50}
	if middle.Applies(Hostility, FriendshipBad) {
		t.Error("交友度 50 不該算「關係惡劣」（門檻 20）")
	}
	if middle.Applies(Cooperate, FriendshipGood) {
		t.Error("交友度 50 不該算「關係良好」（門檻 80）")
	}
	bad := Situation{Aggression: aggr, Friendship: peaceBit + 19}
	if !bad.Applies(Hostility, FriendshipBad) {
		t.Error("交友度 19 應該算「關係惡劣」")
	}
	good := Situation{Aggression: aggr, AllyFriendship: peaceBit + 80}
	if !good.Applies(Cooperate, FriendshipGood) {
		t.Error("交友度 80 應該算「關係良好」")
	}
}

// 每個指令的選單永遠是五項：四個理由 ＋ 撤回進言。
//
// ⚠ 這條先前寫成「五個理由 ＋ 撤回 ＝ 6 項」，是把說明書那句
// 「常に 5 つの項目」讀成了 5 個理由。原版的選單訊息
// （102／166／230）各正好五行，最後一行就是撤回。
func TestAlwaysFiveOptionsPlusWithdraw(t *testing.T) {
	for _, c := range []Command{Hostility, CeaseFire, Cooperate} {
		opts := Options(c)
		if len(opts) != 5 {
			t.Errorf("%v 的選項有 %d 個, want 5（4 ＋ 撤回）", c, len(opts))
		}
		if opts[len(opts)-1] != Withdraw {
			t.Errorf("%v 的最後一項應該是進言撤回", c)
		}
	}
}

// ⭐ 選到不符合狀況的理由 → 說服失敗、信賴度下降。
func TestWrongReasonCostsTrust(t *testing.T) {
	s := Situation{Aggression: 10, OurCities: 5, TheirCities: 20, Friendship: 90}
	sess := Begin(Hostility, s)
	out, dt := sess.Offer(WeAreStronger) // 5 對 20，不成立
	if out != Failed {
		t.Errorf("選到不成立的理由 → %v, want Failed", out)
	}
	if dt >= 0 {
		t.Errorf("失敗的信賴度變化 = %d, 應該是負的", dt)
	}
}

// ⭐ 進言撤回不損信賴度（說明書：この場合は信頼度は変化しません）。
func TestWithdrawIsFree(t *testing.T) {
	sess := Begin(CeaseFire, Situation{Aggression: 10})
	out, dt := sess.Offer(Withdraw)
	if out != Withdrawn {
		t.Errorf("→ %v, want Withdrawn", out)
	}
	if dt != 0 {
		t.Errorf("撤回的信賴度變化 = %d, want 0", dt)
	}
}

// 講滿足夠的成立理由 → 君主同意，信賴度上升。
func TestAgreementRaisesTrust(t *testing.T) {
	s := Situation{
		Aggression: 15, // 呂布：好戰，敵對提案只要少數理由
		OurCities:  20, TheirCities: 5,
		TheirFunds: -1, Friendship: 5,
		TheyInvadeThirdParty: true,
	}
	sess := Begin(Hostility, s)
	total := 0
	agreed := false
	for _, r := range Options(Hostility) {
		if r == Withdraw {
			break
		}
		out, dt := sess.Offer(r)
		total += dt
		if out == Failed {
			t.Fatalf("理由 %v 應該成立卻失敗", r)
		}
		if out == Agreed {
			agreed = true
			break
		}
	}
	if !agreed {
		t.Fatal("所有成立的理由都講完了君主還是沒同意")
	}
	if total <= 0 {
		t.Errorf("同意後的信賴度變化 = %d, 應該是正的", total)
	}
}

// ⭐ sub_13C1E 的四段信賴度級別直接決定 sub_13BA9 要求的成立理由數。
func TestTrustTierDeterminesRequiredReasons(t *testing.T) {
	for _, tc := range []struct {
		trust int
		need  int
	}{
		{0xFF, 1},
		{0xE0, 1},
		{0xDF, 2},
		{0x90, 2},
		{0x8F, 3},
		{0x20, 3},
		{0x1F, 4},
		{0, 4},
	} {
		for _, c := range []Command{Hostility, CeaseFire, Cooperate} {
			if got := Begin(c, Situation{Trust: tc.trust}).Remaining(); got != tc.need {
				t.Errorf("信賴度 0x%02X、%v：要 %d 個理由，得到 %d",
					tc.trust, c, tc.need, got)
			}
		}
	}
}

// sub_13830 的 AL=2 路徑完成後 +10，錯選理由則立即 −20。
func TestReasonPathUsesOriginalTrustDeltas(t *testing.T) {
	s := Situation{Trust: 0xE0, OurCities: 10, TheirCities: 1}
	sess := Begin(Hostility, s)
	if out, dt := sess.Offer(WeAreStronger); out != Agreed || dt != TrustOnReasonSuccess {
		t.Fatalf("高信賴度選到成立理由：(%v, %d)，want (Agreed, %d)",
			out, dt, TrustOnReasonSuccess)
	}

	sess = Begin(Hostility, Situation{Trust: 0xE0})
	if out, dt := sess.Offer(WeAreStronger); out != Failed || dt != TrustOnFailure {
		t.Fatalf("錯選理由：(%v, %d)，want (Failed, %d)", out, dt, TrustOnFailure)
	}
}

// sub_13830 的第一反應碼：AL=1 為 +20；AL=0、3 為 −20；AL=4 不變。
func TestFirstReactionUsesOriginalTrustDeltas(t *testing.T) {
	for _, tc := range []struct {
		reaction Reaction
		delta    int
	}{
		{Agree, TrustOnImmediateSuccess},
		{Refuse, TrustOnFailure},
		{AlreadyAtWar, TrustOnFailure},
		{AskReason, 0},
		{SameFaction, 0},
	} {
		if got := ReactionTrustDelta(tc.reaction); got != tc.delta {
			t.Errorf("反應碼 %d：信賴度變化 %d，want %d", tc.reaction, got, tc.delta)
		}
	}
}

// 同一個理由重複講不算數，也不該罰。
func TestRepeatedReasonIsNeutral(t *testing.T) {
	s := Situation{Aggression: 10, TheyInvadeUs: true}
	sess := Begin(Hostility, s)
	sess.Offer(WeAreDefending)
	before := sess.Remaining()
	out, dt := sess.Offer(WeAreDefending)
	if out == Failed || dt != 0 {
		t.Errorf("重複同一個理由 → %v (%d), 不該被罰", out, dt)
	}
	if sess.Remaining() != before {
		t.Error("重複的理由不該讓進度前進")
	}
}

// ⭐ 成立的理由都講完了君主還不點頭 → 應該用撤回收手，而不是硬選。
func TestExhaustedSuggestsWithdraw(t *testing.T) {
	// 只有一個理由成立（敵勢力疲乏），但君主要的不只一個。
	// 敵對提案的池是 交友惡／我國有利／敵侵他國／敵疲乏，
	// 這裡刻意讓前三個都不成立。
	// ⚠ 交友度要帶和平位元。少了它，`peaceBit + 好戰 + 15` 那個門檻
	// 會把任何裸值都判成「關係惡劣」，於是這裡就不只一個理由成立。
	s := Situation{Aggression: 4, TheirFunds: -1,
		OurCities: 1, TheirCities: 99, Friendship: peaceBit + 100}
	sess := Begin(Hostility, s)
	if sess.Exhausted() {
		t.Fatal("一開始還有成立的理由")
	}
	sess.Offer(EnemyExhausted)
	if !sess.Exhausted() {
		t.Error("成立的理由講完之後 Exhausted 應該為 true")
	}
	// 此時撤回不損信賴度。
	if _, dt := sess.Offer(Withdraw); dt != 0 {
		t.Errorf("撤回卻扣了 %d 信賴度", dt)
	}
}

// ⭐ 佇列裡已經有那筆事件 → 君主當場同意，而且**這一關排在最前面**：
// 同一個局面沒有佇列事件時是拒絕的。
// ⭐ 佇列裡已經有那筆事件 → 君主當場同意，而且**這一關排在最前面**。
func TestFirstReactionQueuedBeatsThreshold(t *testing.T) {
	s := Situation{Aggression: 5, Friendship: peaceBit + 45} // 45 ≥ 5×2+20
	if got := FirstReaction(Hostility, s, true); got != Agree {
		t.Fatalf("佇列裡有事件應同意，得到 %d", got)
	}
	if got := FirstReaction(Hostility, s, false); got != Refuse {
		t.Fatalf("沒有佇列事件應拒絕，得到 %d", got)
	}
}

// ⭐ 「已經在交戰」與「原本就沒交戰」是**同一個反應碼、相反的條件**。
//
//	敵對  和平位元沒設（在打）→ 3「不是已經在交戰狀態中了嗎！」
//	停戰  和平位元設著（沒打）→ 3「原本就沒有和\3交戰啊！」
func TestAlreadyAtWarIsMirroredBetweenCommands(t *testing.T) {
	atWar := Situation{Aggression: 5, Friendship: 25} // 沒有和平位元
	atPeace := Situation{Aggression: 5, Friendship: peaceBit + 25}
	if got := FirstReaction(Hostility, atWar, false); got != AlreadyAtWar {
		t.Errorf("敵對：交戰中應回 3，得到 %d", got)
	}
	if got := FirstReaction(CeaseFire, atPeace, false); got != AlreadyAtWar {
		t.Errorf("停戰：和平中應回 3，得到 %d", got)
	}
	// 反過來各自不該回 3。
	if got := FirstReaction(Hostility, atPeace, false); got == AlreadyAtWar {
		t.Error("敵對：和平中不該回 3")
	}
	if got := FirstReaction(CeaseFire, atWar, false); got == AlreadyAtWar {
		t.Error("停戰：交戰中不該回 3")
	}
}

// ⭐ 好戰等級在三個指令的拒絕門檻上有三個不同的係數。
func TestRefusalThresholdsDifferPerCommand(t *testing.T) {
	// 敵對：交友度 ≥ 好戰×2+20 → 拒絕。好戰 5 → 30。
	for _, c := range []struct {
		f    int
		want Reaction
	}{{29, AskReason}, {30, Refuse}} {
		s := Situation{Aggression: 5, Friendship: peaceBit + c.f}
		if got := FirstReaction(Hostility, s, false); got != c.want {
			t.Errorf("敵對 交友度 %d：得到 %d，期望 %d", c.f, got, c.want)
		}
	}
	// 停戰：交友度 < 好戰÷2 → 拒絕。好戰 15 → 7。交戰中才走得到。
	for _, c := range []struct {
		f    int
		want Reaction
	}{{6, Refuse}, {7, AskReason}} {
		s := Situation{Aggression: 15, Friendship: c.f}
		if got := FirstReaction(CeaseFire, s, false); got != c.want {
			t.Errorf("停戰 交友度 %d：得到 %d，期望 %d", c.f, got, c.want)
		}
	}
	// 協力：對協力對象的交友度 < 0x80+好戰×4+30 → 拒絕。好戰 5 → 50。
	for _, c := range []struct {
		f    int
		want Reaction
	}{{49, Refuse}, {50, AskReason}} {
		s := Situation{Aggression: 5, AllyFriendship: peaceBit + c.f, Friendship: 10}
		if got := FirstReaction(Cooperate, s, false); got != c.want {
			t.Errorf("協力 交友度 %d：得到 %d，期望 %d", c.f, got, c.want)
		}
	}
}

// 協力：兩邊選成同一家 → 反應碼 4（訊息 83「軍師並不是來談笑的」）。
// **這一關排在所有門檻之前**，所以就算其他條件都不成立也是回 4。
func TestCooperateSameFactionIsCheckedFirst(t *testing.T) {
	s := Situation{Aggression: 5, SameFactionPicked: true, AllyFriendship: 0}
	if got := FirstReaction(Cooperate, s, false); got != SameFaction {
		t.Fatalf("應回 4，得到 %d", got)
	}
}

// 協力：被侵攻對象打，而且我方國力不到它的一半 → 不必說理由，直接同意。
func TestCooperateAgreesWhenBadlyOutmatched(t *testing.T) {
	base := Situation{Aggression: 5, AllyFriendship: peaceBit + 100,
		Friendship: 10, TheyInvadeUs: true, OurCities: 1, TheirCities: 100}
	if got := FirstReaction(Cooperate, base, false); got != Agree {
		t.Fatalf("國力懸殊又被打，應直接同意，得到 %d", got)
	}
	// 差距不夠大就要說理由。
	base.OurCities = 60
	if got := FirstReaction(Cooperate, base, false); got != AskReason {
		t.Fatalf("差距不夠大應問理由，得到 %d", got)
	}
}

// ⭐ 「我正在防禦戰」同一個選項、兩個條件。
//
// 停戰提案掃**全部 22 個勢力**（`sub_16577` 的 `mov cx, 16h` 迴圈，
// 而且跳過交涉對象本身）；請求協助只看**侵攻對象**
// （`sub_166D9` 的 `cmp al, byte_10CFF`）。
//
// 差別看得見：被第三方入侵時，停戰進言可以用這個理由，協力進言不行。
func TestDefendingMeansDifferentThingsPerCommand(t *testing.T) {
	byThird := Situation{AnyoneInvadesUs: true, TheyInvadeUs: false}
	if !byThird.Applies(CeaseFire, WeAreDefending) {
		t.Error("停戰：被第三方入侵就算防禦戰")
	}
	if byThird.Applies(Cooperate, WeAreDefending) {
		t.Error("協力：看的是侵攻對象，第三方不算")
	}
	byTarget := Situation{AnyoneInvadesUs: false, TheyInvadeUs: true}
	if !byTarget.Applies(Cooperate, WeAreDefending) {
		t.Error("協力：侵攻對象正在打我們就算")
	}
}

// 疲弊看的是誰的資金，兩個指令相反。
func TestExhaustionSideDependsOnCommand(t *testing.T) {
	weBroke := Situation{OurFunds: -1, TheirFunds: 1000}
	if !weBroke.Applies(CeaseFire, WeAreExhausted) {
		t.Error("停戰「我國力疲乏」看我方資金")
	}
	if weBroke.Applies(Hostility, EnemyExhausted) {
		t.Error("敵對「敵勢力疲乏」看對方資金")
	}
}

// 協力的兩個「強大」比的是不同勢力：協力對象 vs 侵攻對象。
func TestCooperateComparesTwoDifferentFactions(t *testing.T) {
	s := Situation{Aggression: 5, OurCities: 10,
		AllyCities: 100, TheirCities: 1}
	if !s.Applies(Cooperate, AllyIsStronger) {
		t.Error("協力對象 100 個據點應該算強大")
	}
	if s.Applies(Cooperate, InvaderIsStronger) {
		t.Error("侵攻對象只有 1 個據點，不該算強大")
	}
}

// 國力相等時**兩邊都不成立**——不是其中一邊。
func TestPowerTieFavoursNeither(t *testing.T) {
	// 好戰 5：我方係數 25，對方係數 25，據點數相同即打平。
	s := Situation{Aggression: 5, OurCities: 8, TheirCities: 8}
	if s.Applies(Hostility, WeAreStronger) {
		t.Error("打平不該算「我國較有利」")
	}
	if s.Applies(CeaseFire, EnemyIsStronger) {
		t.Error("打平不該算「對我國較不利」")
	}
}

// 兩道閘各自擋得住，都過才出陣（docs/spec/49 §3）。
func TestAcceptSortieNeedsBothGates(t *testing.T) {
	base := Sortie{Funds: 20000, Aggression: 5, Reserves: 700}
	if !AcceptSortie(base) {
		t.Fatal("兩道都過卻不出陣")
	}
	poor := base
	poor.Funds = SortieFundsGate(base.Aggression) - 1
	if AcceptSortie(poor) {
		t.Error("錢不夠還出陣")
	}
	thin := base
	thin.Reserves = SortieReserveGate - 1
	if AcceptSortie(thin) {
		t.Error("兵不夠還出陣")
	}
	// 邊界：兩個都是「大於等於」。
	edge := Sortie{Funds: SortieFundsGate(5), Aggression: 5, Reserves: SortieReserveGate}
	if !AcceptSortie(edge) {
		t.Error("剛好踩線應該過")
	}
}

// 好戰等級越高門檻越低（`(15 − 好戰) × 1,024`）。
func TestSortieFundsGateScalesWithAggression(t *testing.T) {
	for _, tc := range []struct{ aggression, want int }{
		{15, 0}, {14, 1024}, {0, 15360},
		{-3, 15360}, {99, 0}, // 值域外要收斂，不要算出負門檻
	} {
		if got := SortieFundsGate(tc.aggression); got != tc.want {
			t.Errorf("好戰 %d ⇒ %d，要 %d", tc.aggression, got, tc.want)
		}
	}
}
