package persuasion

import "testing"

// 說明書 3.9 明說「国力の基準は拠点数からチェックされ」——
// 國力比的是**據點數**，不是兵力也不是資金。
func TestPowerUsesCityCount(t *testing.T) {
	s := Situation{OurCities: 14, TheirCities: 8}
	if !s.Applies(WeAreStronger) {
		t.Error("14 個據點對 8 個應該算有利")
	}
	if s.Applies(EnemyIsStronger) {
		t.Error("同一局勢不該兩邊都成立")
	}
	// 資金再多也不影響國力判定。
	s.OurFunds, s.TheirFunds = 0, 655000
	if !s.Applies(WeAreStronger) {
		t.Error("國力只看據點數，資金不該影響")
	}
}

// 疲弊的定義是**資金 < 0**（機器碼用 24 位值的符號位判斷）。
func TestExhaustionIsNegativeFunds(t *testing.T) {
	s := Situation{OurFunds: -1, TheirFunds: 1}
	if !s.Applies(WeAreExhausted) {
		t.Error("資金 −1 應該算疲弊")
	}
	if s.Applies(EnemyExhausted) {
		t.Error("資金 +1 不該算疲弊")
	}
	s.TheirFunds = 0
	if s.Applies(EnemyExhausted) {
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
	if !lu.Applies(WeAreStronger) {
		t.Error("好戰 15 的君主在 10 對 12 時應該仍覺得有利")
	}
	if liu.Applies(WeAreStronger) {
		t.Error("好戰 4 的君主在 10 對 12 時不該覺得有利")
	}
}

// ⭐ 說明書點名的兩個例子：呂布不太差就能說服，劉備要相當差才行。
func TestFriendshipGateMatchesManualExamples(t *testing.T) {
	const mildlyBad = 25 // 交友值偏低但不算最惡
	lu := Situation{Aggression: 15, Friendship: mildlyBad}
	liu := Situation{Aggression: 4, Friendship: mildlyBad}
	if !lu.Applies(FriendshipBad) {
		t.Error("呂布（好戰 15）在交友值 25 時應該接受「交友關係惡」")
	}
	if liu.Applies(FriendshipBad) {
		t.Error("劉備（好戰 4）在交友值 25 時不該接受，要更差才行")
	}
	// 劉備要非常差才接受。
	liu.Friendship = 5
	if !liu.Applies(FriendshipBad) {
		t.Error("交友值 5 連劉備都該接受了")
	}
}

// 交友關係惡與良是互斥的：同一個局勢只有一個成立。
func TestFriendshipReasonsAreExclusive(t *testing.T) {
	for _, f := range []int{0, 25, 50, 100} {
		s := Situation{Aggression: 10, Friendship: f}
		if s.Applies(FriendshipBad) == s.Applies(FriendshipGood) {
			t.Errorf("交友值 %d：兩個理由同時 %v", f, s.Applies(FriendshipBad))
		}
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

// ⭐ 好戰的君主容易被說服去打人，消極的容易被說服停戰。
func TestAggressionShapesWhatLordAccepts(t *testing.T) {
	luHostile := requiredReasons(Hostility, 15)
	liuHostile := requiredReasons(Hostility, 4)
	if luHostile >= liuHostile {
		t.Errorf("好戰的君主對敵對提案要 %d 個理由，消極的要 %d —— 方向反了",
			luHostile, liuHostile)
	}
	luPeace := requiredReasons(CeaseFire, 15)
	liuPeace := requiredReasons(CeaseFire, 4)
	if luPeace <= liuPeace {
		t.Errorf("好戰的君主對停戰要 %d 個理由，消極的要 %d —— 方向反了",
			luPeace, liuPeace)
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
	s := Situation{Aggression: 4, TheirFunds: -1,
		OurCities: 1, TheirCities: 99, Friendship: 100}
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
func TestFirstReactionQueuedBeatsThreshold(t *testing.T) {
	s := Situation{Aggression: 5, Friendship: 45} // 45 ≥ 5×2+20
	if got := FirstReaction(s, true, false); got != Agree {
		t.Fatalf("佇列裡有事件應同意，得到 %d", got)
	}
	if got := FirstReaction(s, false, false); got != Refuse {
		t.Fatalf("沒有佇列事件應拒絕，得到 %d", got)
	}
}

func TestFirstReactionAlreadyAtWar(t *testing.T) {
	s := Situation{Aggression: 5, Friendship: 10}
	if got := FirstReaction(s, false, true); got != AlreadyAtWar {
		t.Fatalf("交戰中應回 3，得到 %d", got)
	}
}

// 好戰等級直接決定君主肯不肯聽：同一個交友度，換個君主結果相反。
func TestFirstReactionThresholdMovesWithAggression(t *testing.T) {
	for _, c := range []struct {
		aggression int
		want       Reaction
	}{{5, Refuse}, {6, AskReason}} { // 門檻 30 vs 32，交友度 30
		s := Situation{Aggression: c.aggression, Friendship: 30}
		if got := FirstReaction(s, false, false); got != c.want {
			t.Fatalf("好戰 %d 應回 %d，得到 %d", c.aggression, c.want, got)
		}
	}
}

// 拒絕門檻（好戰×2+20）與「外交關係惡劣」門檻（好戰+15）是兩條線，
// 中間那一段就是「君主願意聽，但這個理由不成立」。
func TestTwoThresholdsAreDistinct(t *testing.T) {
	s := Situation{Aggression: 5, Friendship: 25} // 20 ≤ 25 < 30
	if got := FirstReaction(s, false, false); got != AskReason {
		t.Fatalf("應該進說服迴圈，得到 %d", got)
	}
	if s.Applies(FriendshipBad) {
		t.Error("交友度 25 ≥ 好戰+15 ＝ 20，「外交關係惡劣」不該成立")
	}
}
