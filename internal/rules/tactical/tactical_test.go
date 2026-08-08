package tactical

import (
	"os"
	"testing"
)

type fixedRand struct {
	seq []int
	i   int
}

func (r *fixedRand) Next() int {
	v := r.seq[r.i%len(r.seq)]
	r.i++
	return v
}

// flatField 是一張平坦的野戰戰場（沒有城，所以 gateX ＝ 0）。
func flatField() *Field {
	stack := make([][]int, Height)
	for y := range stack {
		stack[y] = make([]int, Width)
	}
	return NewField(stack, 0)
}

// walledField 是一張攻城用的戰場：中間一道 X ＝ 32 的城牆（堆疊 4 層），
// 只有 gateX 那一格是通的。
func walledField(gateX int) *Field {
	stack := make([][]int, Height)
	for y := range stack {
		stack[y] = make([]int, Width)
		if y != Height/2 {
			stack[y][32] = 4
		}
	}
	return NewField(stack, gateX)
}

func newTestBattle(f *Field) *Battle {
	b := NewBattle(f, SyntheticFormations(), &fixedRand{seq: []int{1, 7, 3}}, 0)
	for k := 0; k < Squads; k++ {
		b.Deploy(0, k, Infantry, 100)
		b.Deploy(1, k, Infantry, 100)
	}
	b.Place()
	return b
}

// 一側 48 個兵、6 隊 × 8 人，其餘進待機。
// 說明書 4.1：一個編成位置 1,000 人 ＝ 100 個兵，場上只放得下 8 個。
func TestDeploymentSplitsOnFieldAndReserve(t *testing.T) {
	b := newTestBattle(flatField())
	if got := b.Sides[0].Alive(); got != SoldiersOnFoot {
		t.Errorf("場上 %d 個兵，應為 %d", got, SoldiersOnFoot)
	}
	for k, r := range b.Sides[0].Reserve {
		if r != 100-PerSquad {
			t.Errorf("第 %d 隊待機 %d 人，應為 %d", k, r, 100-PerSquad)
		}
	}
	if got := b.Sides[0].Remaining(); got != 600 {
		t.Errorf("總戰力 %d，應為 600", got)
	}
}

// 城壁移動在沒有城的戰場自動變成攻擊。
// 這一條在原版的指令 3 與 13 各有一份，而且腳本作者也從不在野戰段下它。
func TestScaleWallFallsBackOnOpenField(t *testing.T) {
	b := newTestBattle(flatField())
	b.Order(0, -1, ScaleWal)
	if got := b.Sides[0].Soldiers[1].Next; got != Attack {
		t.Errorf("野戰下城壁移動變成 %v，應為攻擊", got)
	}

	b2 := newTestBattle(walledField(32))
	b2.Order(0, -1, ScaleWal)
	if got := b2.Sides[0].Soldiers[1].Next; got != ScaleWal {
		t.Errorf("攻城戰下城壁移動變成 %v，應保持不變", got)
	}
}

// 疲勞度只有「走到陣形位置那一刻」才補滿——下令不會補。
// 說明書 4.2 說陣形是唯一能恢復疲勞的指令，6.1 說它「兵回到指定位置時最小」。
func TestStaminaOnlyRefillsOnArrival(t *testing.T) {
	b := newTestBattle(flatField())
	s := &b.Sides[0].Soldiers[3]
	// 先把他挪走並耗掉體力。
	s.X, s.Y = 10, 10
	s.Stamina = 5
	s.Cmd, s.Next = Attack, Form

	b.updateSoldier(0, 3)
	if s.Stamina > 5 {
		t.Errorf("才剛下令就補到 %d，應該要走回去才補", s.Stamina)
	}
	// 直接放到定位再跑一幀。
	x, y := b.formationSpot(0, 3)
	s.X, s.Y = x, y
	b.updateSoldier(0, 3)
	if s.Stamina != StaminaFull {
		t.Errorf("到位後疲勞度 %d，應補滿 %d", s.Stamina, StaminaFull)
	}
	if s.Cmd != Holding {
		t.Errorf("到位後命令是 %v，應轉成就位", s.Cmd)
	}
}

// 攻擊時疲勞度被壓到 40 —— 打起來就回不到滿的。
func TestAttackCapsStamina(t *testing.T) {
	b := newTestBattle(flatField())
	s := &b.Sides[0].Soldiers[2]
	s.Stamina = StaminaFull
	s.Cmd, s.Next = Attack, Attack
	b.updateSoldier(0, 2)
	if s.Stamina > StaminaFighting {
		t.Errorf("攻擊中疲勞度 %d，上限應為 %d", s.Stamina, StaminaFighting)
	}
}

// 大將不會陣亡：體力扣到 1 就停住。
// 說明書 6.1「將軍體力…低於一定值自動退卻」——退卻不是陣亡。
func TestGeneralNeverDies(t *testing.T) {
	b := newTestBattle(flatField())
	g := &b.Sides[0].Soldiers[0]
	if !g.IsGeneral() {
		t.Fatal("第 0 隊的隊長應該是大將")
	}
	for i := 0; i < 100; i++ {
		b.hit(1, g, 50)
	}
	if !g.Alive {
		t.Error("大將陣亡了")
	}
	if g.HP != 1 {
		t.Errorf("大將體力 %d，應停在 1", g.HP)
	}
}

// 大將體力低於 50 → 全軍退卻。
func TestGeneralRetreatOrdersWholeSide(t *testing.T) {
	b := newTestBattle(flatField())
	b.Sides[0].Soldiers[0].HP = GeneralRetreatHP - 1
	b.checkGeneralRetreat()
	for k := range b.Sides[0].Soldiers {
		if s := &b.Sides[0].Soldiers[k]; s.Alive && s.Next != Retreat {
			t.Fatalf("第 %d 個兵的命令是 %v，應為退卻", k, s.Next)
		}
	}
}

// 步兵挨箭只吃四分之一 —— 說明書「攻城戦では弓兵、歩兵が必要です」的
// 數值依據。
func TestInfantryResistsArrows(t *testing.T) {
	b := newTestBattle(flatField())
	inf := &Soldier{Alive: true, Kind: Infantry, HP: MaxHP}
	cav := &Soldier{Alive: true, Kind: Cavalry, HP: MaxHP}
	b.hit(0, inf, 40)
	b.hit(0, cav, 40)
	if MaxHP-inf.HP != 10 {
		t.Errorf("步兵掉了 %d，應為 10（40 ÷ 4）", MaxHP-inf.HP)
	}
	if MaxHP-cav.HP != 40 {
		t.Errorf("騎馬掉了 %d，應為 40", MaxHP-cav.HP)
	}
}

// 有利／不利：差距 ≤ 8 判成普通。
// 說明書 6.1「敵も同じ状態の場合は通常と変わりません」就是這個門檻。
func TestAdvantageEvenBand(t *testing.T) {
	for _, tc := range []struct {
		mine, theirs int
		want         Advantage
	}{
		{40, 40, Even},       // 40 − 47 ＝ −7，在 ±8 內
		{48, 40, Even},       // 48 − 47 ＝ +1
		{55, 40, Even},       // 55 − 47 ＝ +8，剛好在邊界上
		{56, 40, Advantaged}, // 56 − 47 ＝ +9，出帶
		{38, 40, Disadvantaged},
		{20, 40, Disadvantaged},
	} {
		if got := computeAdvantage(tc.mine, tc.theirs); got != tc.want {
			t.Errorf("computeAdvantage(%d, %d) ＝ %v，應為 %v",
				tc.mine, tc.theirs, got, tc.want)
		}
	}
}

// 騎馬與大將爬不上城牆 —— 說明書 5.5「騎馬のみの編成では城壁に登れない」。
func TestCavalryCannotClimb(t *testing.T) {
	for _, tc := range []struct {
		k    Kind
		want bool
	}{
		{General, false}, {Cavalry, false}, {Archer, true}, {Infantry, true},
	} {
		if got := (&Soldier{Kind: tc.k}).CanClimb(); got != tc.want {
			t.Errorf("%v 爬牆 ＝ %v，應為 %v", tc.k, got, tc.want)
		}
	}

	b := newTestBattle(walledField(32))
	// 把一個騎馬放在牆邊，要它走到牆上。
	s := &b.Sides[0].Soldiers[1]
	s.Kind = Cavalry
	s.X, s.Y, s.Z = 31, 10, 0
	s.GoalX, s.GoalY, s.GoalZ = 32, 10, 4
	b.moveToward(0, 1)
	if s.X == 32 {
		t.Error("騎馬爬上城牆了")
	}
	// 換成步兵就上得去（一次一層，所以要走四幀）。
	s.Kind = Infantry
	for i := 0; i < 8 && s.X != 32; i++ {
		s.Z = b.Field.StandLevel(s.X, s.Y)
		b.moveToward(0, 1)
		if s.Z == 4 {
			break
		}
		s.Z++ // 逐層往上
	}
	if !s.CanClimb() {
		t.Error("步兵應該爬得上去")
	}
}

// 退卻不可打斷 —— 說明書 4.2「一旦執行不能取消」。
func TestRetreatCannotBeInterrupted(t *testing.T) {
	s := &Soldier{Alive: true, Cmd: Retreat, Next: Attack}
	if s.applyNewOrder() {
		t.Error("退卻中卻接受了新命令")
	}
	if s.Cmd != Retreat {
		t.Errorf("命令變成 %v，應維持退卻", s.Cmd)
	}
}

// 跑一場完整的戰鬥：一定會結束，而且勝方是還有兵的那一側。
func TestBattleTerminates(t *testing.T) {
	b := newTestBattle(flatField())
	b.Order(0, -1, Attack)
	b.Order(1, -1, Attack)
	for i := 0; i < 200000 && !b.Done; i++ {
		b.Step()
	}
	if !b.Done {
		t.Fatalf("跑了 20 萬幀還沒結束（攻方剩 %d、守方剩 %d）",
			b.Sides[0].Remaining(), b.Sides[1].Remaining())
	}
	if b.Sides[1-b.Winner].Remaining() != 0 {
		t.Errorf("判給第 %d 側，但對方還剩 %d",
			b.Winner, b.Sides[1-b.Winner].Remaining())
	}
	t.Logf("第 %d 幀結束，勝方 %d（剩 %d 對 %d）", b.Frame, b.Winner,
		b.Sides[b.Winner].Remaining(), b.Sides[1-b.Winner].Remaining())
}

// 原版的陣形表載得進來，而且性質與 docs/re/11 §5.8d 對得上。
func TestRealFormationTable(t *testing.T) {
	const exe = "../../../workplace/orig/dosv/KI.EXE"
	if _, err := os.Stat(exe); err != nil {
		t.Skip("找不到原版 KI.EXE，跳過")
	}
	f, err := LoadFormations(exe)
	if err != nil {
		t.Fatal(err)
	}
	// 陣形 0 是 3 格寬 × 49 格高的縱列（中央突破用的窄陣形）。
	minX, maxX, minY, maxY := f.Bounds(0)
	if w, h := maxX-minX+1, maxY-minY+1; w != 3 || h != 49 {
		t.Errorf("陣形 0 是 %d × %d，應為 3 × 49", w, h)
	}
	// 陣形 15 是最密集的一個。
	minX, maxX, minY, maxY = f.Bounds(15)
	if w, h := maxX-minX+1, maxY-minY+1; w != 8 || h != 9 {
		t.Errorf("陣形 15 是 %d × %d，應為 8 × 9", w, h)
	}
	// ⭐ 陣形 4／5／6（同形狀的上／中／下三個位置）把**六個隊長全排在最後面**。
	//
	// ⚠ 這只對這三個成立。本專案一度從陣形 5 的圖推廣成「所有陣形都這樣」，
	// 跑這條測試才發現十六個裡只有三個是（docs/re/11 §5.8d）。
	for _, form := range []int{4, 5, 6} {
		lo, _, _, _ := f.Bounds(form)
		for k := 0; k < SoldiersOnFoot; k += PerSquad {
			if x, _ := f.Offset(form, k); x != lo {
				t.Errorf("陣形 %d 第 %d 隊的隊長在 X=%d，應為最後面的 %d",
					form, k/PerSquad, x, lo)
			}
		}
	}
	// 其餘陣形不該通過同一條檢查——否則就是我又把特例當成通則了。
	allBack := 0
	for form := 0; form < NumFormations; form++ {
		lo, _, _, _ := f.Bounds(form)
		ok := true
		for k := 0; k < SoldiersOnFoot; k += PerSquad {
			if x, _ := f.Offset(form, k); x != lo {
				ok = false
			}
		}
		if ok {
			allBack++
		}
	}
	if allBack != 3 {
		t.Errorf("有 %d 個陣形把隊長全排在最後面，應為 3（4／5／6）", allBack)
	}
}

// 腳本直譯器：等待、下令、分支都要照原版的編碼跑。
func TestScriptBasics(t *testing.T) {
	b := newTestBattle(flatField())
	// e3 00 ＝ 指令 3、參數 7（全軍）、運算元 0（陣形）
	// 00 05 ＝ 等待 5 幀
	// 63 01 ＝ 指令 3、參數 3（第 3 隊）、運算元 1（攻擊）
	code := make([]byte, ScriptCodeSize)
	copy(code, []byte{0xe3, 0x00, 0x00, 0x05, 0x63, 0x01})
	s := NewScript(code, 0)

	s.Step(b)
	if got := b.Sides[0].Soldiers[1].Next; got != Form {
		t.Errorf("全軍命令是 %v，應為陣形", got)
	}
	// 第二個 Step 讀到「等待 5」，接下來五幀都在扣計時器。
	for i := 0; i < 6; i++ {
		s.Step(b)
		if b.Sides[0].Soldiers[3*PerSquad+1].Next == Attack {
			t.Fatalf("第 %d 幀就執行了等待後面的指令", i)
		}
	}
	s.Step(b)
	if got := b.Sides[0].Soldiers[3*PerSquad+1].Next; got != Attack {
		t.Errorf("第 3 隊的命令是 %v，應為攻擊", got)
	}
	if got := b.Sides[0].Soldiers[1].Next; got == Attack {
		t.Error("指定第 3 隊的命令卻影響到第 0 隊")
	}
}

// 分支指令是 4 byte：後面那個 word 是跳躍目標，低位元組必須是 0。
func TestScriptBranch(t *testing.T) {
	b := newTestBattle(flatField())
	code := make([]byte, ScriptCodeSize)
	// 0: 09 02  q.rand 2        固定亂數的第一個是 1 → cond ＝ 1
	// 2: 4a 00  branch != 0 → 目標（成立）
	// 4: 00 04  目標 ＝ 第 4 個 word ＝ byte 8
	// 6: e3 05  order 全軍 退卻   ← 跳過去就不該執行到
	// 8: e3 01  order 全軍 攻擊
	copy(code, []byte{0x09, 0x02, 0x4a, 0x00, 0x00, 0x04, 0xe3, 0x05, 0xe3, 0x01})
	s := NewScript(code, 0)
	s.Step(b) // q.rand → cond
	s.Step(b) // branch
	s.Step(b) // 目標處的指令
	if got := b.Sides[0].Soldiers[1].Next; got != Attack {
		t.Errorf("跳過去之後的命令是 %v，應為攻擊（分支沒跳對）", got)
	}
}

// 原版的腳本跑得起來，而且不會炸。
func TestRealScriptsRun(t *testing.T) {
	const dat = "../../../workplace/orig/dosv/BATTLE.DAT"
	raw, err := os.ReadFile(dat)
	if err != nil {
		t.Skip("找不到原版 BATTLE.DAT，跳過")
	}
	for seg := 0; seg < 32; seg++ {
		b := newTestBattle(flatField())
		b.SetScript(0, NewScript(raw[seg*256:(seg+1)*256], 0))
		b.SetScript(1, NewScript(raw[seg*256:(seg+1)*256], 1))
		for i := 0; i < 20000 && !b.Done; i++ {
			b.Step()
		}
		if b.Sides[0].Remaining() < 0 || b.Sides[1].Remaining() < 0 {
			t.Fatalf("段 %d 跑出負的兵數", seg)
		}
	}
}
