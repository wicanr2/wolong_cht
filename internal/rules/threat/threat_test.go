package threat

import "testing"

// 上限的兩段在交界處必須連續——不連續就表示我把門檻讀錯了。
func TestCorpsCapIsContinuousAtTheBoundary(t *testing.T) {
	for _, c := range []struct {
		funds, want int
	}{
		{0, 5},
		{40_960, 5},      // 交界：40960 ÷ 8192 ＝ 5，兩式同值
		{40_959, 5},      // 交界下一格仍是 5
		{49_152, 6},      // 6 × 8192
		{81_920, 10},
		{-1, 5},          // 資金是有號 24 位，負值不能讓上限變負
	} {
		if got := CorpsCap(c.funds); got != c.want {
			t.Errorf("CorpsCap(%d) ＝ %d，期望 %d", c.funds, got, c.want)
		}
	}
}

func TestBudgetSubtractsExistingCorps(t *testing.T) {
	if got := Budget(0, 2); got != 3 {
		t.Errorf("Budget(0,2) ＝ %d，期望 3", got)
	}
	// 額度用完就是 0，不是負數——原版走的是 CF=1 那條路，不派。
	if got := Budget(0, 9); got != 0 {
		t.Errorf("Budget(0,9) ＝ %d，期望 0", got)
	}
	if got := Budget(81_920, 4); got != 6 {
		t.Errorf("Budget(81920,4) ＝ %d，期望 6", got)
	}
}

// 和平的鄰居不算威脅。這是交友度 bit 7 極性的第二條獨立證據，
// 弄反的話 AI 開局就會把所有人當敵人。
func TestPeaceBitMakesNeighbourHarmless(t *testing.T) {
	ns := []Neighbour{{Site: 7, Owner: 3, Occupancy: 4, Friendship: 0x80 | 55}}
	r := Scan(1, NoTarget, 1, ns)
	if r.Threatened || r.Level != 0 {
		t.Fatalf("和平的鄰居被當成威脅：%+v", r)
	}
	ns[0].Friendship = 55 // 交戰
	r = Scan(1, NoTarget, 1, ns)
	if !r.Threatened || r.Level != 4 {
		t.Fatalf("交戰的鄰居沒有被算進威脅：%+v", r)
	}
}

// +0x1B ＝ 0 時原版直接 retn，連掃都不掃。
func TestNoEnemyNeighboursSkipsTheScanEntirely(t *testing.T) {
	ns := []Neighbour{{Site: 7, Owner: 3, Occupancy: 9, Friendship: 0}}
	if r := Scan(1, 3, 0, ns); r.Threatened || r.Level != 0 || r.Targets != nil {
		t.Fatalf("+0x1B ＝ 0 還是掃了：%+v", r)
	}
}

// 中立鄰居不看交友度（勢力表沒有第 24 筆），但仍然要比對侵攻目標。
func TestNeutralNeighbourIsNotWeighedByFriendship(t *testing.T) {
	ns := []Neighbour{{Site: 7, Owner: Neutral, Occupancy: 3, Friendship: 0}}
	r := Scan(1, NoTarget, 1, ns)
	if r.Threatened || r.Level != 0 {
		t.Fatalf("中立鄰居被算成威脅：%+v", r)
	}
}

// bit 6 的分野是「有沒有具體目標」，不是遠近。
func TestSpecificMeansTheInvasionTargetOwnsANeighbour(t *testing.T) {
	ns := []Neighbour{
		{Site: 7, Owner: 3, Occupancy: 2, Friendship: 0},  // 交戰，但不是侵攻目標
		{Site: 9, Owner: 5, Occupancy: 6, Friendship: 0},  // 侵攻目標
	}
	r := Scan(1, 5, 2, ns)
	if !r.Threatened {
		t.Fatal("兩個交戰鄰居卻沒有受威脅")
	}
	if !r.Specific {
		t.Fatal("侵攻目標就在隔壁卻沒有設 Specific")
	}
	if r.Level != 8 {
		t.Errorf("威脅量 ＝ %d，期望 2+6 ＝ 8", r.Level)
	}
	if len(r.Targets) != 1 || r.Targets[0] != 9 {
		t.Errorf("目標清單 ＝ %v，期望 [9]", r.Targets)
	}
}

// 只有一個交戰鄰居而它不是侵攻目標：受威脅但沒有具體目標。
func TestThreatWithoutSpecificTarget(t *testing.T) {
	ns := []Neighbour{{Site: 7, Owner: 3, Occupancy: 2, Friendship: 0}}
	r := Scan(1, 5, 1, ns)
	if !r.Threatened || r.Specific || len(r.Targets) != 0 {
		t.Fatalf("%+v", r)
	}
}

func TestRequested(t *testing.T) {
	// 威脅量 5、自己這格站了 2 支 → 要 5
	if got := Requested(5, 2); got != 5 {
		t.Errorf("Requested(5,2) ＝ %d，期望 5", got)
	}
	// 守得住就不求援
	if got := Requested(0, 3); got != 0 {
		t.Errorf("Requested(0,3) ＝ %d，期望 0", got)
	}
}

// 低 4 位是「哪幾個鄰居是敵方」，不是「哪幾個方向有鄰接」——
// 空槽與同勢力的鄰居都不設位元。
func TestEnemyMaskCountsOnlyForeignNeighbours(t *testing.T) {
	ns := []Neighbour{
		{Site: 1, Owner: 1},  // 同勢力
		{Site: 2, Owner: 3},  // 敵方
		{Site: -1},           // 沒有這個鄰居
		{Site: 4, Owner: Neutral},
	}
	mask, n := EnemyMask(1, ns)
	if mask != 0b1010 || n != 2 {
		t.Fatalf("mask ＝ %04b、count ＝ %d，期望 1010／2", mask, n)
	}
}

func TestPlayerCooldownRange(t *testing.T) {
	lo, hi := 0xFF, 0
	for r := 0; r < 256; r++ {
		v := PlayerCooldown(r)
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	if lo != 24 || hi != 39 {
		t.Fatalf("玩家冷卻值域 ＝ %d–%d，期望 24–39", lo, hi)
	}
}

// 距離算式的 Y 分量減的是首都的 X。這是原版行為（兩版 byte 相同），
// 照抄不修正——所以這一條測的是「不對稱有沒有被好心改掉」。
func TestAICooldownKeepsTheOriginalAsymmetry(t *testing.T) {
	// 首都 X ＝ 100。據點 (100, 100)：對稱算式會得 0，原版也是 0。
	if got := AICooldown(100, 100, 100); got != 0 {
		t.Errorf("AICooldown(100,100,100) ＝ %d，期望 0", got)
	}
	// 據點 (100, 20)：|100−100| ＋ |20−100| ＝ 80 → 80÷8 ＝ 10。
	// 若第二項改成減首都 Y，這個值就會不同——不對稱被改掉會在這裡爆。
	if got := AICooldown(100, 20, 100); got != 10 {
		t.Errorf("AICooldown(100,20,100) ＝ %d，期望 10", got)
	}
	// 上限 30
	if got := AICooldown(370, 248, 4); got != 30 {
		t.Errorf("遠距離沒有被夾到 30，得到 %d", got)
	}
}

func TestDispatchOnlyTakesCorpsStandingAtTheSite(t *testing.T) {
	gs := []Garrison{
		{At: 9, Ready: true},          // 別的據點
		{At: 5, Ready: true},          // ✓
		{At: 5, Ready: false},         // 位元 2 沒設
		{At: 5, Ready: true, Stage: 8}, // +0x23 到上限
		{At: 5, Ready: true},          // ✓
	}
	got := Dispatch(gs, 5, 4, 0, func() int { return 0xFF })
	want := []int{1, 4}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("選出 %v，期望 %v", got, want)
	}
}

// 還有跳過額度時，亂數 < 0x40 的那一支被跳過，而且**額度只在跳過時才減**。
func TestDispatchSkipBudgetIsSpentOnlyWhenItSkips(t *testing.T) {
	gs := []Garrison{{At: 5, Ready: true}, {At: 5, Ready: true}, {At: 5, Ready: true}}
	rolls := []int{0x00, 0xFF, 0xFF} // 第一支被跳過，之後兩支都派
	i := 0
	got := Dispatch(gs, 5, 2, 1, func() int { r := rolls[i]; i++; return r })
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("選出 %v，期望 [1 2]", got)
	}
	if i != 1 {
		t.Fatalf("抽了 %d 次亂數，期望 1——額度用完之後不該再抽", i)
	}
}
