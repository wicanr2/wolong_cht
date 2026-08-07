package diplomacy

import "testing"

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

func seq(v ...int) Rand { return &fixedRand{seq: v} }

// 交友度是一個 byte：低 7 位是值，最高位元是交戰旗標。
func TestFriendshipEncoding(t *testing.T) {
	f := Friendship(50)
	if f.Value() != 50 || f.AtWar() {
		t.Errorf("50 → 值 %d 交戰 %v", f.Value(), f.AtWar())
	}
	w := f.WithWar(true)
	if w.Value() != 50 || !w.AtWar() {
		t.Errorf("設交戰後 → 值 %d 交戰 %v（值不該變）", w.Value(), w.AtWar())
	}
	// 改值要保留交戰旗標。
	w2 := w.WithValue(80)
	if w2.Value() != 80 || !w2.AtWar() {
		t.Errorf("改值後 → 值 %d 交戰 %v（旗標應保留）", w2.Value(), w2.AtWar())
	}
	// 上限 100。
	if got := f.WithValue(999).Value(); got != MaxFriendship {
		t.Errorf("超過上限 → %d, want %d", got, MaxFriendship)
	}
}

// 交戰壓過一切：說明書說交戰中固定顯示「交戰」，但值仍在漲。
func TestLevelAtWarOverrides(t *testing.T) {
	f := Friendship(95).WithWar(true)
	if f.Level() != AtWar {
		t.Errorf("交戰中顯示 %v, want 交戰", f.Level())
	}
	if f.Value() != 95 {
		t.Errorf("交戰中的值 = %d，值本身不該被蓋掉", f.Value())
	}
	if f.WithWar(false).Level() != Intimate {
		t.Errorf("停戰後應該顯示親密，得到 %v", f.WithWar(false).Level())
	}
}

// 外交官：12.5% 機率動作，動了就扣經費（23 − 政治）。
func TestDiplomatCostAndChance(t *testing.T) {
	// 亂數 0x20 = 32 → 不小於 32，不動作。
	d := &Diplomat{Politics: 10, Budget: 100}
	f := Friendship(0)
	if d.Tick(&f, seq(0x20)) || d.Budget != 100 {
		t.Errorf("機率沒過卻動作了，經費 %d", d.Budget)
	}

	// 亂數 0x1F = 31 → 動作。政治 10 → 消耗 13。
	d = &Diplomat{Politics: 10, Budget: 100}
	d.Tick(&f, seq(0x1F, 0))
	if d.Budget != 87 {
		t.Errorf("經費 = %d, want 87（100 − (23 − 10)）", d.Budget)
	}
}

// 政治同時決定消耗與成功率 —— 說明書只講了錢。
func TestPoliticsAffectsBothCostAndSuccess(t *testing.T) {
	// 政治 15：消耗 8，rand(0..15) 一定 ≤ 15 → 必成功。
	hi := &Diplomat{Politics: 15, Budget: 100}
	f := Friendship(0)
	if !hi.Tick(&f, seq(0, 15)) {
		t.Error("政治 15 遇到亂數 15 應該成功")
	}
	if hi.Budget != 92 {
		t.Errorf("政治 15 的消耗 = %d, want 8", 100-hi.Budget)
	}

	// 政治 1：消耗 22，rand=2 > 1 → 失敗，但錢照扣。
	lo := &Diplomat{Politics: 1, Budget: 100}
	g := Friendship(0)
	if lo.Tick(&g, seq(0, 2)) {
		t.Error("政治 1 遇到亂數 2 不該成功")
	}
	if lo.Budget != 78 {
		t.Errorf("政治 1 的消耗 = %d, want 22", 100-lo.Budget)
	}
	if g.Value() != 0 {
		t.Errorf("失敗卻加了交友度：%d", g.Value())
	}
}

// 經費歸零就完全停工。
func TestDiplomatStopsWhenBroke(t *testing.T) {
	d := &Diplomat{Politics: 15, Budget: 0}
	f := Friendship(0)
	if d.Tick(&f, seq(0, 0)) {
		t.Error("沒錢還在工作")
	}
}

// 成功一次加 1，而且不會超過 100。
func TestFriendshipIncrementAndCap(t *testing.T) {
	d := &Diplomat{Politics: 15, Budget: 10000}
	f := Friendship(0)
	rng := seq(0, 0)
	for i := 0; i < 500; i++ {
		d.Tick(&f, rng)
	}
	if f.Value() != MaxFriendship {
		t.Errorf("跑很久之後 = %d, want %d", f.Value(), MaxFriendship)
	}
}

// 要價公式，並確認取的是兩個方向的較小者。
func TestDemand(t *testing.T) {
	// 我方看對方 60、對方看我方 30 → 取 30 → (125 − 30) × 200
	if got := Demand(Friendship(60), Friendship(30)); got != 95*200 {
		t.Errorf("要價 = %d, want %d", got, 95*200)
	}
	// 反過來一樣。
	if got := Demand(Friendship(30), Friendship(60)); got != 95*200 {
		t.Errorf("方向對調後 = %d, 應該一樣", got)
	}
	// 交戰旗標讓基準從 125 降到 100。
	if got := Demand(Friendship(30).WithWar(true), Friendship(60)); got != 70*200 {
		t.Errorf("交戰中要價 = %d, want %d", got, 70*200)
	}
}

// ⭐ 協力是確定性的布林判定，不是機率 —— remake 不要在這裡擲骰子。
func TestCooperationIsDeterministic(t *testing.T) {
	allyToUs, allyToTarget := Friendship(60), Friendship(30)

	if !CanRequestCooperation(false, allyToUs, allyToTarget, 0) {
		t.Error("條件滿足卻不成立")
	}
	// 同樣條件呼叫一百次，結果必須完全一致。
	for i := 0; i < 100; i++ {
		if !CanRequestCooperation(false, allyToUs, allyToTarget, 0) {
			t.Fatal("同樣條件下結果不一致 —— 這裡不該有隨機性")
		}
	}
	// 協力方正在打別人 → 不成立。
	if CanRequestCooperation(true, allyToUs, allyToTarget, 0) {
		t.Error("協力方正在侵攻時不該成立")
	}
	// 交友度反過來 → 不成立。
	if CanRequestCooperation(false, Friendship(30), Friendship(60), 0) {
		t.Error("對我方比對目標差時不該成立")
	}
	// 對方君主的額外門檻。
	if CanRequestCooperation(false, Friendship(35), Friendship(30), 10) {
		t.Error("差距不足以跨過門檻時不該成立")
	}
}

// 停戰的條件：對方沒有正在侵攻我方（說明書 7.3）。
func TestCeaseFire(t *testing.T) {
	if CanCeaseFire(true) {
		t.Error("對方正在侵攻時不該能停戰")
	}
	if !CanCeaseFire(false) {
		t.Error("對方沒在侵攻時應該能停戰")
	}
}
