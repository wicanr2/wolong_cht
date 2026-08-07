package economy

import "testing"

// fixedRand 讓赤字懲罰在測試裡是決定性的。
type fixedRand struct {
	seq []int
	i   int
}

func (r *fixedRand) Next() int {
	if len(r.seq) == 0 {
		return 0
	}
	v := r.seq[r.i%len(r.seq)]
	r.i++
	return v
}

func zeroRand() Rand { return &fixedRand{} }

// 距離除數表照抄 ds:5532h／ds:5535h：≤80→2、≤200→3、其餘→4。
// 距離用的是**切比雪夫距離**，不是歐氏也不是曼哈頓。
func TestDistanceDivisor(t *testing.T) {
	cap0 := City{X: 100, Y: 100}
	cases := []struct {
		x, y int
		want int
	}{
		{100, 100, 2}, // 距離 0
		{180, 100, 2}, // 距離 80，邊界內
		{181, 100, 3}, // 距離 81
		{100, 300, 3}, // 距離 200，邊界內
		{100, 301, 4}, // 距離 201
	}
	for _, c := range cases {
		if got := distanceDivisor(City{X: c.x, Y: c.y}, cap0); got != c.want {
			t.Errorf("(%d,%d) → %d, want %d", c.x, c.y, got, c.want)
		}
	}
	// 切比雪夫：取兩軸較大者，不是相加。
	if got := distanceDivisor(City{X: 20, Y: 40}, cap0); got != 2 {
		t.Errorf("|Δx|=80 |Δy|=60 → %d, want 2（切比雪夫取 80，非曼哈頓的 140）", got)
	}
	// 超過 255 的距離被夾住，仍落在最後一級。
	if got := distanceDivisor(City{X: 383, Y: 255}, City{X: 0, Y: 0}); got != 4 {
		t.Errorf("地圖對角 → %d, want 4", got)
	}
}

// 募兵配比要照抄原版的移位序列，不能寫成 a*19/32 這種分數式。
// 北方與中部的三份加起來必須**正好等於 a**（餘數被第三份吸收）。
func TestSplitRecruitsExactSum(t *testing.T) {
	for a := 0; a < 2000; a++ {
		for _, y := range []int{0, 79, 80, 149} { // 北方與中部
			cav, arc, inf := splitRecruits(a, y)
			if cav+arc+inf != a {
				t.Fatalf("y=%d a=%d → (%d,%d,%d) 合計 %d, want %d",
					y, a, cav, arc, inf, cav+arc+inf, a)
			}
		}
	}
}

// 南方那一支**會少掉奇數位**：合計是 a 去掉最低位。
// 這不是 bug，是原版的移位序列本來就這樣（騎馬 a>>5、弓 a>>1、
// 步 (a>>1)−(a>>5)，合計 2×(a>>1)）。照抄，不要「修好」它。
func TestSplitRecruitsSouthLosesOddBit(t *testing.T) {
	for a := 0; a < 2000; a++ {
		cav, arc, inf := splitRecruits(a, 200)
		want := a &^ 1
		if cav+arc+inf != want {
			t.Fatalf("南方 a=%d → 合計 %d, want %d", a, cav+arc+inf, want)
		}
	}
}

// 這一組是文件裡記下來的分岔案例：a=33 時
// 原版給 (20,1,12)，寫成分數式會給 (19,1,12)。
func TestSplitRecruitsKnownValues(t *testing.T) {
	cases := []struct {
		a, y          int
		cav, arc, inf int
	}{
		{32, 0, 19, 1, 12},   // 北方，剛好整除
		{33, 0, 20, 1, 12},   // 北方，餘數給騎馬 ← 分數式會算成 19
		{32, 100, 4, 4, 24},  // 中部
		{33, 100, 4, 4, 25},  // 中部，餘數給步兵
		{32, 200, 1, 16, 15}, // 南方
	}
	for _, c := range cases {
		cav, arc, inf := splitRecruits(c.a, c.y)
		if cav != c.cav || arc != c.arc || inf != c.inf {
			t.Errorf("a=%d y=%d → (%d,%d,%d), want (%d,%d,%d)",
				c.a, c.y, cav, arc, inf, c.cav, c.arc, c.inf)
		}
	}
}

// 地域偏差要能重現說明書那句「騎馬は北方でよく集まり、弓は南方でよく集まります」。
func TestRegionalBias(t *testing.T) {
	const a = 3200
	nc, na, _ := splitRecruits(a, 10)
	_, sa, _ := splitRecruits(a, 200)
	mc, ma, mi := splitRecruits(a, 100)

	if nc <= na {
		t.Errorf("北方應該騎馬 > 弓兵，得到 騎馬%d 弓%d", nc, na)
	}
	if sa <= sc(a) {
		t.Errorf("南方應該弓兵 > 騎馬，得到 弓%d 騎馬%d", sa, sc(a))
	}
	if mi <= mc || mi <= ma {
		t.Errorf("中部應該步兵最多，得到 騎馬%d 弓%d 步%d", mc, ma, mi)
	}
}

func sc(a int) int { c, _, _ := splitRecruits(a, 200); return c }

// 收入公式：Σ(生產力 ÷ 距離除數) × 稅率 ÷ 100。
func TestIncome(t *testing.T) {
	capital := City{X: 100, Y: 100}
	cities := []City{
		{X: 100, Y: 100, Production: 4000, Owner: 0}, // 距離 0 → ÷2 → 2000
		{X: 250, Y: 100, Production: 3000, Owner: 0}, // 距離 150 → ÷3 → 1000
		{X: 100, Y: 100, Production: 9999, Owner: 1}, // 別人的，不算
	}
	f := &Faction{Funds: 0, Capital: capital, TaxRate: 50,
		RecruitCap: [NumTroopTypes]int{9999, 9999, 9999}}
	res := Settle(f, cities, 0, zeroRand())

	if res.GrossBase != 3000 {
		t.Errorf("套稅率前的收入 = %d, want 3000", res.GrossBase)
	}
	if res.Income != 1500 {
		t.Errorf("收入 = %d, want 1500（3000 × 50%%）", res.Income)
	}
	if res.Cities != 2 {
		t.Errorf("據點數 = %d, want 2", res.Cities)
	}
	if f.Cities != 2 {
		t.Errorf("據點數沒寫回勢力：%d", f.Cities)
	}
}

// AI 不看稅率，收入固定除以 2（原版 sub_15456 的 shr/rcr）。
func TestAIIgnoresTaxRate(t *testing.T) {
	cities := []City{{X: 0, Y: 0, Production: 4000, Owner: 0}}
	for _, tax := range []int{0, 25, 100} {
		f := &Faction{Capital: City{X: 0, Y: 0}, TaxRate: tax, AI: true}
		res := Settle(f, cities, 0, zeroRand())
		if res.Income != 1000 { // 4000/2（距離）/2（AI）
			t.Errorf("稅率 %d%% 時 AI 收入 = %d, want 1000", tax, res.Income)
		}
	}
}

// 募兵數設定是上限不是目標。
func TestRecruitCapIsCeiling(t *testing.T) {
	cities := []City{{X: 0, Y: 0, Production: 60000, Owner: 0}}
	f := &Faction{Capital: City{X: 0, Y: 0}, TaxRate: 100,
		RecruitCap: [NumTroopTypes]int{10, 10, 10}}
	res := Settle(f, cities, 0, zeroRand())
	for tt := TroopType(0); tt < NumTroopTypes; tt++ {
		if res.Recruited[tt] != 10 {
			t.Errorf("兵種 %d 募到 %d, want 10（被上限擋住）", tt, res.Recruited[tt])
		}
	}

	// 反過來：據點供不上時，募到的比設定少。
	poor := []City{{X: 0, Y: 0, Production: 200, Owner: 0}}
	g := &Faction{Capital: City{X: 0, Y: 0}, TaxRate: 100,
		RecruitCap: [NumTroopTypes]int{9999, 9999, 9999}}
	res2 := Settle(g, poor, 0, zeroRand())
	total := res2.Recruited[0] + res2.Recruited[1] + res2.Recruited[2]
	if total >= 9999 {
		t.Errorf("窮據點卻募到 %d 人，上限不該是保證", total)
	}
}

// 資金上下限對稱鉗制在 ±655,000。
func TestFundsClamp(t *testing.T) {
	f := &Faction{Funds: MaxFunds - 10, Capital: City{},
		TaxRate: 100, RecruitCap: [NumTroopTypes]int{9999, 9999, 9999}}
	Settle(f, []City{{Production: 100000, Owner: 0}}, 0, zeroRand())
	if f.Funds != MaxFunds {
		t.Errorf("資金 = %d, want 上限 %d", f.Funds, MaxFunds)
	}

	g := &Faction{Funds: 0, Expense: 10_000_000, Capital: City{}}
	Settle(g, nil, 0, zeroRand())
	if g.Funds != MinFunds {
		t.Errorf("資金 = %d, want 下限 %d", g.Funds, MinFunds)
	}
}

// 赤字懲罰：每個兵種各扣 (|資金|>>8)×16 ＋ 亂數(0–31)。
// 注意是**每個兵種各扣一次**，總損失是三倍。
func TestDeficitPenalty(t *testing.T) {
	f := &Faction{
		Funds:    0,
		Expense:  16000, // 結算後資金 = −16000
		Capital:  City{},
		Reserves: [NumTroopTypes]int{5000, 5000, 5000},
	}
	res := Settle(f, nil, 0, &fixedRand{seq: []int{0}})
	want := ((16000) >> 8) * 16 // = 62 × 16 = 992
	for tt := TroopType(0); tt < NumTroopTypes; tt++ {
		if res.Deficit[tt] != want {
			t.Errorf("兵種 %d 扣了 %d, want %d", tt, res.Deficit[tt], want)
		}
		if f.Reserves[tt] != 5000-want {
			t.Errorf("兵種 %d 剩 %d, want %d", tt, f.Reserves[tt], 5000-want)
		}
	}
	if !f.Exhausted() {
		t.Error("資金為負時 Exhausted() 應為 true")
	}
}

// 預備兵扣到負要歸零，不能變成負數。
func TestDeficitFloorsAtZero(t *testing.T) {
	f := &Faction{Funds: -600000, Capital: City{},
		Reserves: [NumTroopTypes]int{10, 10, 10}}
	Settle(f, nil, 0, zeroRand())
	for tt := TroopType(0); tt < NumTroopTypes; tt++ {
		if f.Reserves[tt] != 0 {
			t.Errorf("兵種 %d = %d, want 0", tt, f.Reserves[tt])
		}
	}
}

// 資金非負時完全不觸發懲罰。
func TestNoPenaltyWhenSolvent(t *testing.T) {
	f := &Faction{Funds: 1, Capital: City{},
		Reserves: [NumTroopTypes]int{5000, 5000, 5000}}
	res := Settle(f, nil, 0, &fixedRand{seq: []int{31}})
	for tt := TroopType(0); tt < NumTroopTypes; tt++ {
		if res.Deficit[tt] != 0 || f.Reserves[tt] != 5000 {
			t.Errorf("資金為正卻扣了兵：%+v", res.Deficit)
		}
	}
	if f.Exhausted() {
		t.Error("資金為正時 Exhausted() 應為 false")
	}
}

// 月結的順序：先扣支出 → 再算收入 → 收入入帳 → 赤字懲罰。
// 支出扣完之後要歸零，不能累積到下個月重扣。
func TestExpenseClearedAfterSettle(t *testing.T) {
	f := &Faction{Funds: 10000, Expense: 3000, Capital: City{}}
	Settle(f, nil, 0, zeroRand())
	if f.Expense != 0 {
		t.Errorf("支出沒歸零：%d", f.Expense)
	}
	if f.Funds != 7000 {
		t.Errorf("資金 = %d, want 7000", f.Funds)
	}
}

// 順序驗證：支出大到會讓資金轉負時，募兵仍在懲罰**之前**發生，
// 所以那個月是「先募到再被扣掉」。
func TestRecruitHappensBeforeDeficitPenalty(t *testing.T) {
	cities := []City{{X: 0, Y: 0, Production: 40000, Owner: 0}}
	f := &Faction{
		Funds: 0, Expense: 300000, Capital: City{X: 0, Y: 0},
		TaxRate: 1, RecruitCap: [NumTroopTypes]int{9999, 9999, 9999},
	}
	res := Settle(f, cities, 0, zeroRand())
	total := res.Recruited[0] + res.Recruited[1] + res.Recruited[2]
	if total == 0 {
		t.Fatal("赤字月完全沒募到兵——順序被寫反了")
	}
	if res.Deficit[0] == 0 {
		t.Fatal("赤字月沒有觸發懲罰")
	}
}

// 預備兵上限 65,500（每個兵種各自）。
func TestReserveCap(t *testing.T) {
	f := &Faction{
		Capital: City{X: 0, Y: 0}, TaxRate: 100,
		Reserves:   [NumTroopTypes]int{MaxReserve - 5, 0, 0},
		RecruitCap: [NumTroopTypes]int{99999, 99999, 99999},
	}
	Settle(f, []City{{X: 0, Y: 0, Production: 60000, Owner: 0}}, 0, zeroRand())
	if f.Reserves[Cavalry] != MaxReserve {
		t.Errorf("騎馬 = %d, want 上限 %d", f.Reserves[Cavalry], MaxReserve)
	}
}

// ---------------------------------------------------------------------------
// 生產力與上昇值
// ---------------------------------------------------------------------------

// 變化量與生產力本身成正比 —— 這是複利模型不是線性模型。
// 說明書：「大きい数値の方が変化が大きくなります」。
func TestGrowthIsProportional(t *testing.T) {
	small := CityState{Production: 1000, ProductionCap: 999999, Growth: 10}
	big := CityState{Production: 20000, ProductionCap: 999999, Growth: 10}
	before := [2]int{small.Production, big.Production}

	GrowCity(&small, 0, false, zeroRand())
	GrowCity(&big, 0, false, zeroRand())

	ds := small.Production - before[0]
	db := big.Production - before[1]
	if db <= ds {
		t.Errorf("大據點的增量 %d 應該大於小據點的 %d", db, ds)
	}
	// (生產力>>8) × 上昇值 / 2
	if want := (1000 >> 8) * 10 / 2; ds != want {
		t.Errorf("小據點增量 = %d, want %d", ds, want)
	}
	if want := (20000 >> 8) * 10 / 2; db != want {
		t.Errorf("大據點增量 = %d, want %d", db, want)
	}
}

// 生產力不會超過該據點的上限。
func TestProductionCap(t *testing.T) {
	c := CityState{Production: 20714, ProductionCap: 21000, Growth: 11}
	GrowCity(&c, 0, false, zeroRand())
	if c.Production != 21000 {
		t.Errorf("生產力 = %d, want 被上限 21000 鉗住", c.Production)
	}
}

// 上昇值為負時生產力下降，而且扣的是**全額**不是一半。
func TestNegativeGrowthShrinks(t *testing.T) {
	c := CityState{Production: 20000, ProductionCap: 999999, Growth: -10}
	GrowCity(&c, 0, false, zeroRand())
	want := 20000 - (20000>>8)*10 // 全額，不是 /2
	if c.Production != want {
		t.Errorf("生產力 = %d, want %d（負成長扣全額）", c.Production, want)
	}
}

// 生產力不會掉到負數。
func TestProductionFloorsAtZero(t *testing.T) {
	c := CityState{Production: 100, ProductionCap: 999999, Growth: -100}
	for i := 0; i < 50; i++ {
		GrowCity(&c, 0, false, zeroRand())
	}
	if c.Production < 0 {
		t.Errorf("生產力 = %d, 不該為負", c.Production)
	}
}

// ⭐ 稅率的中性點是 30%：低於就繁榮，高於就荒廢。
func TestTaxNeutralPoint(t *testing.T) {
	base := CityState{Production: 20000, ProductionCap: 999999, Growth: 0}

	// 稅率 30% → 上昇值不受稅率影響（只被自然衰減扣掉）。
	c := base
	GrowCity(&c, TaxNeutral, true, zeroRand())
	if c.Growth != 0 {
		t.Errorf("稅率 30%% 時上昇值 = %d, want 0", c.Growth)
	}

	// 稅率 10% → 上昇值 +20。
	c = base
	GrowCity(&c, 10, true, zeroRand())
	if c.Growth != 20 {
		t.Errorf("稅率 10%% 時上昇值 = %d, want +20", c.Growth)
	}

	// 稅率 60% → 上昇值 −30。
	c = base
	GrowCity(&c, 60, true, zeroRand())
	if c.Growth != -30 {
		t.Errorf("稅率 60%% 時上昇值 = %d, want −30", c.Growth)
	}
}

// 稅率修正只對玩家的據點生效，AI 的據點完全不受影響。
func TestTaxOnlyAffectsPlayerCities(t *testing.T) {
	ai := CityState{Production: 20000, ProductionCap: 999999, Growth: 5}
	GrowCity(&ai, 90, false, zeroRand()) // 極高稅率，但 applyTax=false
	if ai.Growth != 5 {
		t.Errorf("AI 據點的上昇值 = %d, want 5（不受稅率影響）", ai.Growth)
	}
}

// 上昇值每月自然衰減 rand(0..15)，並鉗在 ±100。
func TestGrowthDecayAndClamp(t *testing.T) {
	c := CityState{Production: 5000, ProductionCap: 999999, Growth: 50}
	GrowCity(&c, 0, false, &fixedRand{seq: []int{15}})
	if c.Growth != 35 {
		t.Errorf("上昇值 = %d, want 35（50 − 15）", c.Growth)
	}

	hi := CityState{Production: 5000, ProductionCap: 999999, Growth: 100}
	GrowCity(&hi, 0, true, zeroRand()) // 稅率 0 → +30
	if hi.Growth != MaxGrowth {
		t.Errorf("上昇值 = %d, want 上限 %d", hi.Growth, MaxGrowth)
	}

	lo := CityState{Production: 5000, ProductionCap: 999999, Growth: -100}
	GrowCity(&lo, 100, true, &fixedRand{seq: []int{15}})
	if lo.Growth != MinGrowth {
		t.Errorf("上昇值 = %d, want 下限 %d", lo.Growth, MinGrowth)
	}
}

// 暴動的免疫門檻是「上昇值存值 ≥ 64」＝ 實際上昇值 ≥ −36，
// 比說明書講的「不為負」寬鬆（docs/re/07 §17）。
func TestRiotGate(t *testing.T) {
	for _, g := range []int{0, 1, 100, -1, -36} {
		if (CityState{Growth: g}).RiotRisk() {
			t.Errorf("上昇值 %d 不該有暴動風險（門檻是 −36）", g)
		}
	}
	for _, g := range []int{-37, -100} {
		if !(CityState{Growth: g}).RiotRisk() {
			t.Errorf("上昇值 %d 應該有暴動風險", g)
		}
	}
}

// 長期行為：稅率低於平衡點（30 − 7.5 ≈ 22.5）時據點會長大，
// 高於就會萎縮。這是攻略章「通常は税率を下げるだけで、
// 内政の必要はありません」那句話的可驗證版本。
func TestLongRunTaxBehaviour(t *testing.T) {
	// 用固定的平均衰減 8 來代表 rand(0..15) 的期望值。
	avg := func() Rand { return &fixedRand{seq: []int{8}} }

	low := CityState{Production: 5000, ProductionCap: 42800, Growth: 4}
	high := CityState{Production: 5000, ProductionCap: 42800, Growth: 4}
	for i := 0; i < 120; i++ { // 十年
		GrowCity(&low, 10, true, avg())
		GrowCity(&high, 50, true, avg())
	}
	if low.Production <= 5000 {
		t.Errorf("低稅十年後生產力 = %d, 應該成長", low.Production)
	}
	if high.Production >= 5000 {
		t.Errorf("高稅十年後生產力 = %d, 應該萎縮", high.Production)
	}
	if !high.RiotRisk() {
		t.Errorf("高稅十年後上昇值 = %d，應該已經跌破 −36（暴動風險）", high.Growth)
	}
}
