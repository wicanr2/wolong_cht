package combat

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
)

// fixedRand 回傳固定序列，讓擲骰的結果可預測。
type fixedRand struct {
	seq []int
	i   int
}

func (r *fixedRand) Next() int {
	v := r.seq[r.i%len(r.seq)]
	r.i++
	return v
}

// fullCorps 是六槽滿員（各 100 ＝ 1000 人）的軍團。
func fullCorps(kind army.TroopType, l Leader, morale int) Corps {
	c := Corps{Leader: l, Morale: morale}
	for i := range c.Units {
		c.Units[i] = Unit{Men: 100, Kind: kind}
		c.Men += 100
	}
	return c
}

// 地形係數表是從 KI.EXE 檔案偏移 0x5320 讀出來的 16 個 byte。
// 這一條把「讀到的位元組」與「程式用的列」釘在一起——
// 表整個換位置或欄序改了，這裡會先炸。
func TestTerrainFactorTable(t *testing.T) {
	want := [16]int{
		2, 3, 3, 0,
		3, 2, 1, 0,
		1, 3, 2, 0,
		2, 1, 2, 0,
	}
	i := 0
	for _, row := range terrainFactor {
		for _, v := range row {
			if v != want[i] {
				t.Errorf("係數表第 %d 個 byte ＝ %d，應為 %d", i, v, want[i])
			}
			i++
		}
	}
	// 攻守用的列：野戰雙方同列，攻城分開。
	for _, tc := range []struct {
		m        Mode
		attacker bool
		want     int
	}{
		{Field, true, 1}, {Field, false, 1},
		{Siege, true, 3}, {Siege, false, 0},
	} {
		if got := rowFor(tc.m, tc.attacker); got != tc.want {
			t.Errorf("rowFor(%v, attacker=%v) ＝ %d，應為 %d", tc.m, tc.attacker, got, tc.want)
		}
	}
}

// 騎兵野外強、城牆邊弱——這是說明書「騎馬のみの編成では城壁に登れない」
// 在數值上的體現。同一支騎兵軍團換個場合，戰力就要換方向。
func TestCavalryStrongInFieldWeakAtWalls(t *testing.T) {
	l := Leader{Martial: 10, Command: 10, SiegeAptitude: 5, FieldAptitude: 5}
	rng := &fixedRand{seq: []int{1}} // 武力 ≥ 統率 且 rand&3 != 0 → 一律用 武力×2

	cav := fullCorps(army.Cavalry, l, 200)
	arc := fullCorps(army.Archer, l, 200)

	cavField := Power(cav, Field, true, 0, rng)
	arcField := Power(arc, Field, true, 0, rng)
	if cavField <= arcField {
		t.Errorf("野戰騎兵 %d 應強於弓兵 %d", cavField, arcField)
	}

	cavHold := Power(cav, Siege, false, 0, rng)
	arcHold := Power(arc, Siege, false, 0, rng)
	if cavHold >= arcHold {
		t.Errorf("守城騎兵 %d 應弱於弓兵 %d", cavHold, arcHold)
	}
}

// 將領值的三個分支。原版擲 `rand & 3`：非 0（75%）用 武力×2。
func TestLeaderValue(t *testing.T) {
	for _, tc := range []struct {
		name    string
		l       Leader
		roll    int
		want    int
		comment string
	}{
		{"武力高，多數情況", Leader{Martial: 12, Command: 6}, 1, 24, "12 × 2"},
		{"武力高，少數情況", Leader{Martial: 12, Command: 6}, 0, 21, "12 − 3 ＋ 12"},
		{"統率高", Leader{Martial: 4, Command: 12}, 0, 21, "12 − 3 ＋ 12"},
	} {
		got := leaderValue(tc.l, &fixedRand{seq: []int{tc.roll}})
		if got != tc.want {
			t.Errorf("%s：leaderValue ＝ %d，應為 %d（%s）", tc.name, got, tc.want, tc.comment)
		}
	}
	// 武力 < 統率時不擲骰——原版 `cmp dl, dh / jb` 直接跳過。
	rng := &fixedRand{seq: []int{0}}
	leaderValue(Leader{Martial: 1, Command: 9}, rng)
	if rng.i != 0 {
		t.Errorf("統率較高時不該擲骰，卻用掉了 %d 次", rng.i)
	}
}

// 說明書 10.5：「武術と統率はその合計が同じであれば強さは同じ」。
// 那是**評價**的性質（internal/rules/general 已釘住）；戰力不是這樣。
// 這一條記錄差異：同樣的和，武力偏重的一方戰力較高。
func TestMartialBeatsCommandAtEqualSum(t *testing.T) {
	rng := &fixedRand{seq: []int{1}}
	base := Leader{SiegeAptitude: 5, FieldAptitude: 5}

	m := base
	m.Martial, m.Command = 12, 4
	c := base
	c.Martial, c.Command = 4, 12

	pm := Power(fullCorps(army.Infantry, m, 200), Field, true, 0, rng)
	pc := Power(fullCorps(army.Infantry, c, 200), Field, true, 0, rng)
	if pm <= pc {
		t.Errorf("武 12 統 4 的戰力 %d 應高於武 4 統 12 的 %d（評價相同但戰力不同）", pm, pc)
	}
}

// 適性讓分母從 16 掉到 6，將領值從 2 升到 30。
// 兩者相乘，最強與最弱的差距要在數十倍這個量級。
func TestLeaderDominatesPower(t *testing.T) {
	rng := &fixedRand{seq: []int{1}}
	best := Leader{Martial: 15, Command: 15, FieldAptitude: 10}
	worst := Leader{Martial: 1, Command: 1, FieldAptitude: 0}

	hi := Power(fullCorps(army.Infantry, best, 200), Field, true, 0, rng)
	lo := Power(fullCorps(army.Infantry, worst, 200), Field, true, 0, rng)
	if lo == 0 || hi/lo < 20 {
		t.Errorf("最強 %d ÷ 最弱 %d ＝ %d 倍，應在數十倍量級", hi, lo, hi/max(lo, 1))
	}
}

// 城市損傷的形狀：苦戰傷得多、輾壓傷得少——直到比值超過 63 繞回去。
// **這個溢位是原版行為**，測試把它釘住，免得日後被當成 bug「修掉」。
func TestCityDamageOverflow(t *testing.T) {
	for _, tc := range []struct{ ratio, want int }{
		{8, 13},  // 1 : 1
		{16, 11}, // 2 : 1
		{32, 7},  // 4 : 1
		{63, 0},  // 約 8 : 1
		{64, 63}, // ★ 減出負數，繞回最大
		{100, 54},
	} {
		if got := CityDamage(tc.ratio); got != tc.want {
			t.Errorf("CityDamage(%d) ＝ %d，應為 %d", tc.ratio, got, tc.want)
		}
	}
}

// 勝負不擲骰：戰力大的一方直接贏。同樣的雙方跑很多次，結果不該變。
func TestOutcomeIsDeterministic(t *testing.T) {
	strong := Leader{Martial: 14, Command: 10, FieldAptitude: 8}
	weak := Leader{Martial: 3, Command: 3, FieldAptitude: 1}

	for i := 0; i < 32; i++ {
		a := fullCorps(army.Cavalry, strong, 200)
		d := fullCorps(army.Infantry, weak, 200)
		rng := &fixedRand{seq: []int{i, i * 7, i * 13, 3, 200, 17}}
		if r := Resolve(&a, &d, Field, 0, rng); r.DefenderWins {
			t.Fatalf("第 %d 次：強方輸了", i)
		}
	}
}

// 敗方的士氣被重設成 100 × 兵力比，戰前有多高都一樣。
// 所以輸過一場的軍團必定低於 100，下一場不論勝負都會壞滅。
func TestLoserMoraleIsResetThenNextBattleRouts(t *testing.T) {
	strong := Leader{Martial: 14, Command: 10, FieldAptitude: 8}
	weak := Leader{Martial: 3, Command: 3, FieldAptitude: 1}

	a := fullCorps(army.Cavalry, strong, 200)
	d := fullCorps(army.Infantry, weak, 200) // 戰前士氣 200
	rng := &fixedRand{seq: []int{1, 4, 30, 2, 90, 5}}

	r := Resolve(&a, &d, Field, 0, rng)
	if r.DefenderWins {
		t.Fatal("前提錯了：這一場應該由攻方贏")
	}
	if d.Morale >= army.RoutMoraleGate {
		t.Errorf("敗方戰後士氣 %d，應低於 %d", d.Morale, army.RoutMoraleGate)
	}
	if r.DefenderDestroyed {
		t.Error("第一場不該壞滅——士氣掉下來但還沒歸零")
	}

	// 同一支軍團再打一場：戰前士氣已經不足 100 → 歸零 → 壞滅。
	a2 := fullCorps(army.Cavalry, strong, 200)
	r2 := Resolve(&a2, &d, Field, 0, rng)
	if !r2.DefenderDestroyed {
		t.Errorf("士氣 %d 的軍團再戰應壞滅", d.Morale)
	}
}

// ⚠ 與說明書的差異：機器碼對**勝方**一樣清零。
// 士氣不足 100 的軍團打贏也會散掉，說明書只寫了「負けると」。
func TestLowMoraleWinnerAlsoRouts(t *testing.T) {
	strong := Leader{Martial: 14, Command: 10, FieldAptitude: 8}
	weak := Leader{Martial: 3, Command: 3, FieldAptitude: 1}

	a := fullCorps(army.Cavalry, strong, 60) // 戰前士氣 60
	d := fullCorps(army.Infantry, weak, 200)
	rng := &fixedRand{seq: []int{1, 4, 30, 2, 90, 5}}

	r := Resolve(&a, &d, Field, 0, rng)
	if r.DefenderWins {
		t.Fatal("前提錯了：士氣 60 但將領與兵種佔優，這一場應該贏")
	}
	if !r.AttackerDestroyed {
		t.Errorf("戰前士氣不足 100 的勝方也該壞滅，實得士氣 %d", a.Morale)
	}
}

// 大將槽永遠留 1：自動判定打不死大將的部隊。
func TestGeneralSlotNeverEmpties(t *testing.T) {
	l := Leader{Martial: 8, Command: 8, FieldAptitude: 5}
	a := fullCorps(army.Cavalry, l, 200)
	d := fullCorps(army.Infantry, l, 200)
	// 每一槽只有 3 人，敗方每槽至少扣 8。
	for i := range d.Units {
		d.Units[i].Men = 3
	}
	d.Men = 18

	Resolve(&a, &d, Field, 0, &fixedRand{seq: []int{1, 200, 7, 3, 90, 250}})
	if d.Units[0].Men != 1 {
		t.Errorf("大將槽剩 %d 人，應保底 1", d.Units[0].Men)
	}
	for i := 1; i < army.Positions; i++ {
		if d.Units[i].Men != 0 {
			t.Errorf("第 %d 槽剩 %d 人，應被打光", i, d.Units[i].Men)
		}
	}
}

// 城兵：六隊步兵、士氣 255，餘數散給前幾槽。
func TestGarrison(t *testing.T) {
	g := Garrison(3, 100)
	if g.Morale != GarrisonMorale {
		t.Errorf("城兵士氣 %d，應為 %d", g.Morale, GarrisonMorale)
	}
	total := 0
	for i, u := range g.Units {
		if u.Kind != army.Infantry {
			t.Errorf("第 %d 槽是 %v，城兵應全是步兵", i, u.Kind)
		}
		total += u.Men
	}
	if total != 100 {
		t.Errorf("城兵合計 %d，應為 100", total)
	}
	// 100 ÷ 6 ＝ 16 餘 4 → 前四槽 17、後兩槽 16。
	want := [army.Positions]int{17, 17, 17, 17, 16, 16}
	for i, u := range g.Units {
		if u.Men != want[i] {
			t.Errorf("第 %d 槽 %d 人，應為 %d", i, u.Men, want[i])
		}
	}
}

// 守城方的城兵加成只有攻城戰的守方拿得到。
func TestGarrisonBonusOnlyForCityDefender(t *testing.T) {
	l := Leader{Martial: 8, Command: 8, SiegeAptitude: 5, FieldAptitude: 5}
	c := fullCorps(army.Infantry, l, 200)
	rng := &fixedRand{seq: []int{1}}

	if Power(c, Siege, false, 0, rng) >= Power(c, Siege, false, 150, rng) {
		t.Error("攻城的守方應該吃到城兵加成")
	}
	if Power(c, Siege, true, 0, rng) != Power(c, Siege, true, 150, rng) {
		t.Error("攻城的攻方不該吃到城兵加成")
	}
	if Power(c, Field, false, 0, rng) != Power(c, Field, false, 150, rng) {
		t.Error("野戰不該吃到城兵加成")
	}
}

// 君主親征絕不被擒——原版在擲骰之前就跳掉了。
func TestRulerNeverCaptured(t *testing.T) {
	c := Captive{Rating: 10, IsRuler: true, HasCapital: true}
	// 骰到最差的值（127 遠大於任何門檻）也一樣。
	if got := RollFate(c, 1, 2, &fixedRand{seq: []int{127}}); got != Escaped {
		t.Errorf("君主的下場是 %v，應為 Escaped", got)
	}
}

// 沒有首都就無處可退，直接被擒——連骰都不擲。
func TestNoCapitalMeansCaptured(t *testing.T) {
	c := Captive{Rating: 66, HasCapital: false, LordSurvives: true}
	rng := &fixedRand{seq: []int{0}}
	if got := RollFate(c, 1, 2, rng); got != Captured {
		t.Errorf("沒有首都時下場是 %v，應為 Captured", got)
	}
	if rng.i != 0 {
		t.Error("沒有首都時不該擲骰")
	}
}

// 能力越高越容易脫身：門檻是 評價 ÷ 2 ＋ 40。
func TestEscapeChanceRisesWithRating(t *testing.T) {
	// 呂布、趙雲（66）與文官（10）的門檻。
	if got := escapeThreshold(66); got != 73 {
		t.Errorf("評價 66 的門檻 ＝ %d，應為 73", got)
	}
	if got := escapeThreshold(10); got != 45 {
		t.Errorf("評價 10 的門檻 ＝ %d，應為 45", got)
	}
	// 骰出 50：高評價逃得掉，低評價逃不掉。
	base := Captive{HasCapital: true, LordSurvives: true}
	hi, lo := base, base
	hi.Rating, lo.Rating = 66, 10
	if got := RollFate(hi, 1, 2, &fixedRand{seq: []int{50}}); got != Escaped {
		t.Errorf("評價 66 骰 50 的下場是 %v，應為 Escaped", got)
	}
	if got := RollFate(lo, 1, 2, &fixedRand{seq: []int{50}}); got != Captured {
		t.Errorf("評價 10 骰 50 的下場是 %v，應為 Captured", got)
	}
}

// 不事二主：舊主已滅 ＋ 武將旗標 bit 4 → 自刎（訊息 0x43）。
func TestSuicideWhenLordIsGone(t *testing.T) {
	c := Captive{Rating: 10, HasCapital: true, LoyalToDeath: true, LordSurvives: false}
	if got := RollFate(c, 1, 2, &fixedRand{seq: []int{127}}); got != Suicide {
		t.Errorf("下場是 %v，應為 Suicide", got)
	}
	// 舊主還在就只是被擒。
	c.LordSurvives = true
	if got := RollFate(c, 1, 2, &fixedRand{seq: []int{127}}); got != Captured {
		t.Errorf("舊主還在時下場是 %v，應為 Captured", got)
	}
	// 沒有那個旗標的武將照樣改隸。
	c.LordSurvives, c.LoyalToDeath = false, false
	if got := RollFate(c, 1, 2, &fixedRand{seq: []int{127}}); got != Captured {
		t.Errorf("沒有不事二主旗標時下場是 %v，應為 Captured", got)
	}
}

// 勝方是無主勢力（0x18）就不抓人。
func TestNeutralVictorTakesNoPrisoners(t *testing.T) {
	c := Captive{Rating: 10, HasCapital: true}
	if got := RollFate(c, NeutralFaction, 2, &fixedRand{seq: []int{127}}); got != Escaped {
		t.Errorf("無主勢力獲勝時下場是 %v，應為 Escaped", got)
	}
}

// 野外駐留的軍費是據點的 24 倍，而且士氣不回。
// 這是原版逼軍團回城的手段，數字要釘住。
func TestUpkeepPunishesTheField(t *testing.T) {
	const men = 3200
	inTown, inField := Upkeep(men, false), Upkeep(men, true)
	if inTown != men/32+1 {
		t.Errorf("據點軍費 ＝ %d，應為 %d", inTown, men/32+1)
	}
	if inField != men*3/4 {
		t.Errorf("野外軍費 ＝ %d，應為 %d", inField, men*3/4)
	}
	if inField/inTown < 20 {
		t.Errorf("野外 %d ÷ 據點 %d ＝ %d 倍，應在 20 倍以上", inField, inTown, inField/inTown)
	}

	c := Corps{Morale: 150}
	Recover(&c, 200, true)
	if c.Morale != 150 {
		t.Errorf("野外士氣變成 %d，不該回復", c.Morale)
	}
	Recover(&c, 200, false)
	if c.Morale != 160 {
		t.Errorf("據點士氣 ＝ %d，應為 160", c.Morale)
	}
	c.Morale = 195
	Recover(&c, 200, false)
	if c.Morale != 200 {
		t.Errorf("士氣 ＝ %d，應被壓在勢力基準 200", c.Morale)
	}
}
