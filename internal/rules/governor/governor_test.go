package governor

import "testing"

// fixed 回一串指定的亂數值，用完就從頭來。
// 這一層要驗的是**公式**，不是亂數分布，所以要把亂數釘死。
func fixed(vs ...int) func() int {
	i := 0
	return func() int {
		v := vs[i%len(vs)]
		i++
		return v
	}
}

// TestGainFromRate 釘住 ch = max(1, rate − 15)。
// `sub ch, 0Fh / ja` 是**無號**比較，所以 rate ≤ 15 一律得到 1。
func TestGainFromRate(t *testing.T) {
	for _, c := range []struct{ rate, want int }{
		{5, 1}, {8, 1}, {15, 1}, {16, 1}, {20, 5}, {25, 10},
	} {
		city := City{Growth: 0, GarrisonCap: 0}
		gov := &Official{Politics: c.rate - PlayerRate, Budget: 1}
		// 亂數固定成 0：兩個判定都必定成功，徵兵那一段被 GarrisonCap 擋掉。
		Tick(&city, gov, true, fixed(0))
		if city.Growth != c.want {
			t.Errorf("rate %d：上昇值增量 %d，預期 %d", c.rate, city.Growth, c.want)
		}
	}
}

// TestAIBaselineBeatsDecay 是這一支存在的理由。
//
// 月結每月扣 rand(0..15)（期望 −7.5），而 AI 據點每天有 9/16 的機率 +1。
// **一個月的期望值必須是正的**，否則所有 AI 據點都會單調掉到暴動
// ——那正是實作這一支之前模擬跑出「120 個月 1872 次暴動」的原因。
func TestAIBaselineBeatsDecay(t *testing.T) {
	// 30 天 × 16 種亂數值輪一圈，數成功次數。
	city := City{Growth: 0, GarrisonCap: 0}
	seq := make([]int, 0, 16)
	for i := 0; i < 16; i++ {
		seq = append(seq, i)
	}
	rnd := fixed(seq...)
	for d := 0; d < 30; d++ {
		Tick(&city, nil, false, rnd)
	}
	// rate=8 → rand(0..15) ≤ 8 成立 9/16；30 天期望 16.9。
	if city.Growth < 10 {
		t.Errorf("AI 據點 30 天只長了 %d，撐不住每月 −7.5 的衰減", city.Growth)
	}
	t.Logf("AI 據點 30 天上昇值 +%d（月衰減期望 −7.5）", city.Growth)
}

// TestPlayerWorseThanAI 釘住那個反直覺的設計：
// **沒派內政官的玩家據點比 AI 的據點還糟。**
// 這不是 bug，是原版的常數（cl 預設 8，玩家分支 5）。
func TestPlayerWorseThanAI(t *testing.T) {
	seq := []int{6, 6, 200} // 6 ≤ 8 成立、6 > 5 不成立；200 讓徵兵不觸發
	ai := City{GarrisonCap: 0}
	Tick(&ai, nil, false, fixed(seq...))
	pl := City{GarrisonCap: 0}
	Tick(&pl, nil, true, fixed(seq...))
	if ai.Growth <= pl.Growth {
		t.Errorf("AI %d 應該要 > 玩家 %d", ai.Growth, pl.Growth)
	}
}

// TestBudgetConsumed 確認經費每次扣 1，且歸零之後內政官就不再加成。
func TestBudgetConsumed(t *testing.T) {
	gov := &Official{Politics: 15, Martial: 10, Budget: 2}
	city := City{GarrisonCap: 0}
	for i := 0; i < 3; i++ {
		Tick(&city, gov, true, fixed(0))
	}
	if gov.Budget != 0 {
		t.Errorf("經費剩 %d，預期 0", gov.Budget)
	}
	// 前兩次 rate=20（gain 5），第三次沒錢 rate=5（gain 1）→ 5+5+1 = 11
	if city.Growth != 11 {
		t.Errorf("上昇值 %d，預期 11（5+5+1）", city.Growth)
	}
}

// TestGarrisonTrimmedWhenOverCap 釘住尾段那個容易漏掉的分支：
// **城兵超過上限時會被修剪回上限，即使這一輪沒徵兵。**
func TestGarrisonTrimmedWhenOverCap(t *testing.T) {
	city := City{Garrison: 200, GarrisonCap: 100}
	Tick(&city, nil, false, fixed(0))
	if city.Garrison != 100 {
		t.Errorf("城兵 %d，預期被修剪回上限 100", city.Garrison)
	}
}

// TestDraftCostsGrowth 徵兵要拿上昇值去換。
func TestDraftCostsGrowth(t *testing.T) {
	city := City{Growth: 50, Garrison: 10, GarrisonCap: 100}
	// 前兩個亂數讓兩個判定失敗（15 > 8），第三個 0 < 24 觸發徵兵。
	Tick(&city, nil, false, fixed(15, 15, 0))
	if city.Garrison != 10+BaseDraft {
		t.Errorf("城兵 %d，預期 %d", city.Garrison, 10+BaseDraft)
	}
	if city.Growth != 50-BaseDraft {
		t.Errorf("上昇值 %d，預期 %d", city.Growth, 50-BaseDraft)
	}
}

// TestGrowthBoundsAreOnStoredValue 釘住 docs/spec/165：上昇值的兩個夾值
// 在原版是做在**存值**上（上限 0C8h、下限 0），而 `City.Growth` 是實際值
// （存值 − 100），所以界是 +100 與 −100，不是 200 與 0。
func TestGrowthBoundsAreOnStoredValue(t *testing.T) {
	// 上限：實際 +100 ＝ 存值 200，再成功也不會長。
	c := City{Growth: MaxGrowth, Prevention: 0, Garrison: 5, GarrisonCap: 5}
	Tick(&c, nil, false, func() int { return 0 }) // 0 ⇒ 兩次判定都成功
	if c.Growth != MaxGrowth {
		t.Fatalf("上昇值上限：得到 %d，要 %d", c.Growth, MaxGrowth)
	}

	// 下限：夾在實際 −100（存值 0）。從 −100 出發，上昇 +1 再徵兵 −4，
	// 結果停在 −100。
	c = City{Growth: MinGrowth, Prevention: 0, Garrison: 0, GarrisonCap: 99}
	Tick(&c, nil, false, func() int { return 0 })
	if c.Growth != MinGrowth {
		t.Fatalf("上昇值下限：得到 %d，要 %d", c.Growth, MinGrowth)
	}

	// ⚠ 負對照：地板若還套在實際值 0 上，下面這一格會停在 0 而不是 −99。
	c = City{Growth: -96, Prevention: 0, Garrison: 0, GarrisonCap: 99}
	Tick(&c, nil, false, func() int { return 0 })
	if c.Growth != -99 {
		t.Fatalf("上昇值要能掉到負值：得到 %d，要 -99", c.Growth)
	}
}
