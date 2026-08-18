package tactical

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
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

// 近戰不是固定威力：命中值取 rand&0x7f 加攻擊者戰力，
// 有利與突擊再分別套用原版的 0x40／0xc8 飽和加成。
func TestMeleeUsesOriginalPowerAndAdvantage(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
	attacker := &Soldier{Alive: true, Kind: Cavalry, HP: MaxHP, Power: 80,
		Cmd: Attack, Next: Attack}
	target := &Soldier{Alive: true, Kind: Infantry, HP: MaxHP}
	b.Advantage[0] = Even
	if !b.meleeHit(0, attacker, target) {
		t.Fatal("rand 0 + 戰力 80 應達到 0x46 命中門檻")
	}
	if target.HP != 20 {
		t.Errorf("一般近戰扣血 %d，應為 80", MaxHP-target.HP)
	}

	b.rng = &fixedRand{seq: []int{0}}
	target = &Soldier{Alive: true, Kind: Infantry, HP: MaxHP}
	b.Advantage[0] = Disadvantaged
	if b.meleeHit(0, attacker, target) {
		t.Fatal("不利時 (0 + 80 - 0x32) 不應達到命中門檻")
	}
	if target.HP != MaxHP {
		t.Errorf("未命中卻扣了 %d 點血", MaxHP-target.HP)
	}

	b.rng = &fixedRand{seq: []int{0}}
	target = &Soldier{Alive: true, Kind: Infantry, HP: MaxHP}
	attacker.Power = 100
	b.Advantage[0] = Advantaged
	if !b.meleeHit(0, attacker, target) || target.Alive {
		t.Fatal("有利時 100 + 0x40 應以飽和威力擊倒普通兵")
	}

	b.rng = &fixedRand{seq: []int{100}}
	target = &Soldier{Alive: true, Kind: Infantry, HP: MaxHP}
	attacker.Power, attacker.Cmd = 40, Charge
	b.Advantage[0] = Even
	if !b.meleeHit(0, attacker, target) || target.Alive {
		t.Fatal("突擊的 40 + 0xc8 應擊倒普通兵")
	}
}

// 隊長離場後，原版的七名隊員退卻，該隊不能再用沒有隊長的待機兵補場。
func TestLeaderLossRetreatsSquadAndDropsReserve(t *testing.T) {
	b := newTestBattle(flatField())
	leader := &b.Sides[0].Soldiers[PerSquad]
	leader.HP = 1
	if !b.applyHit(leader, 1) {
		t.Fatal("隊長應該被命中")
	}
	if leader.Alive {
		t.Fatal("非大將隊長被擊倒後仍在場")
	}
	if b.Sides[0].Reserve[1] != 0 {
		t.Errorf("隊長離場後仍有 %d 個待機兵，應為 0", b.Sides[0].Reserve[1])
	}
	for k := PerSquad + 1; k < 2*PerSquad; k++ {
		if b.Sides[0].Soldiers[k].Next != Retreat {
			t.Errorf("第 %d 個隊員命令是 %v，應為退卻", k, b.Sides[0].Soldiers[k].Next)
		}
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
		b.applyHit(g, 50)
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
//
// ⭐ **只有箭**。那條規則在 `sub_1B97E`（飛道具的命中），
// 近戰的 `sub_1B618` 沒有——本專案一度把它套在兩邊，
// 結果近戰打步兵幾乎不掉血，戰鬥卡住（下面那條測試釘住這件事）。
func TestInfantryResistsArrows(t *testing.T) {
	b := newTestBattle(flatField())
	inf := &Soldier{Alive: true, Kind: Infantry, HP: MaxHP}
	cav := &Soldier{Alive: true, Kind: Cavalry, HP: MaxHP}
	b.hitByArrow(0, inf, 40)
	b.hitByArrow(0, cav, 40)
	if MaxHP-inf.HP != 10 {
		t.Errorf("步兵挨箭掉了 %d，應為 10（40 ÷ 4）", MaxHP-inf.HP)
	}
	if MaxHP-cav.HP != 40 {
		t.Errorf("騎馬挨箭掉了 %d，應為 40", MaxHP-cav.HP)
	}
}

func TestSpecialProjectileUsesCH20AndFallsVertically(t *testing.T) {
	b := newTestBattle(flatField())
	attacker := &b.Sides[0].Soldiers[1]
	*attacker = Soldier{Alive: true, Kind: Archer, HP: MaxHP, Power: 40,
		X: 10, Y: 20, Z: 2, Target: 1, Cmd: Attack, Next: Attack}
	target := &b.Sides[1].Soldiers[1]
	*target = Soldier{Alive: true, Kind: Infantry, HP: MaxHP,
		X: 11, Y: 20, Z: 1, Target: -1, Cmd: Attack, Next: Attack}
	b.shoot(0, 1, target)
	if len(b.projectiles) != 1 || !b.projectiles[0].special {
		t.Fatalf("高處近距離弓兵沒有建立 CH=0x20 效果：%+v", b.projectiles)
	}
	p := b.projectiles[0]
	if p.power != specialProjectilePower || p.x != 11 || p.y != 20 || p.z != 3 {
		t.Fatalf("特殊效果初值錯誤：%+v", p)
	}
	if attacker.ProjectileCooldown != 6 {
		t.Fatalf("sub_1AD7F 成功後 +0x13 應為 6，got %d", attacker.ProjectileCooldown)
	}
	for i := 0; i < 2; i++ {
		b.stepProjectiles()
		if target.HP != MaxHP {
			t.Fatalf("第 %d 幀提前命中：HP=%d", i+1, target.HP)
		}
	}
	b.stepProjectiles()
	expectedPower := specialProjectilePower
	for i := 0; i < 2; i++ {
		expectedPower += expectedPower>>2 + 1 // sub_1BAB7 的兩次下降加成
	}
	if target.HP != MaxHP-expectedPower/InfantryArrowDivisor {
		t.Fatalf("CH=0x20 垂直效果傷害=%d，應為 %d", MaxHP-target.HP,
			expectedPower/InfantryArrowDivisor)
	}
}

func TestSpecialProjectileCarriesPoseBitIntoSecondRawFrame(t *testing.T) {
	b := newTestBattle(flatField())
	attacker := &b.Sides[0].Soldiers[1]
	*attacker = Soldier{Alive: true, Kind: Archer, HP: MaxHP, Power: 40,
		PoseStep: 1, X: 10, Y: 20, Z: 2, Target: 1, Cmd: Attack, Next: Attack}
	target := &b.Sides[1].Soldiers[1]
	*target = Soldier{Alive: true, Kind: Infantry, HP: MaxHP,
		X: 11, Y: 20, Z: 1, Target: -1, Cmd: Attack, Next: Attack}
	b.shoot(0, 1, target)
	views := b.Projectiles()
	if len(views) != 1 || !views[0].Special || views[0].SpecialFrame != 1 {
		t.Fatalf("特殊投射物沒有保留 +0x02 bit 0：%+v", views)
	}
}

func TestClimbingInfantryCanUseSpecialProjectile(t *testing.T) {
	b := newTestBattle(flatField())
	attacker := &b.Sides[0].Soldiers[1]
	*attacker = Soldier{Alive: true, Kind: Infantry, HP: MaxHP, Power: 40,
		X: 10, Y: 20, Z: 1, Climbing: true, Target: 1, Cmd: Attack, Next: Attack}
	target := &b.Sides[1].Soldiers[1]
	*target = Soldier{Alive: true, Kind: Infantry, HP: MaxHP,
		X: 11, Y: 20, Z: 1, Target: -1, Cmd: Attack, Next: Attack}
	b.doAttack(0, 1)
	if len(b.projectiles) != 1 || !b.projectiles[0].special {
		t.Fatalf("高層步兵近距離未走 CH=0x20 分支：%+v", b.projectiles)
	}
}

func TestProjectileCooldownMatchesRawLaunchBranches(t *testing.T) {
	b := newTestBattle(flatField())
	shooter := &b.Sides[0].Soldiers[1]
	*shooter = Soldier{Alive: true, Kind: Archer, HP: MaxHP, Power: 40,
		X: 10, Y: 20, Z: 1, Target: 1, Cmd: Attack, Next: Attack}
	target := &b.Sides[1].Soldiers[1]
	*target = Soldier{Alive: true, Kind: Infantry, HP: MaxHP,
		X: 14, Y: 20, Z: 1, Target: -1, Cmd: Attack, Next: Attack}

	b.shoot(0, 1, target)
	if len(b.projectiles) != 1 || shooter.ProjectileCooldown != 8 {
		t.Fatalf("sub_1AD2D 成功後投射物／冷卻錯誤：count=%d cooldown=%d", len(b.projectiles), shooter.ProjectileCooldown)
	}
	for want := uint8(7); ; want-- {
		b.shoot(0, 1, target)
		if len(b.projectiles) != 1 || shooter.ProjectileCooldown != want {
			t.Fatalf("普通投射物冷卻遞減錯誤：count=%d cooldown=%d want=%d", len(b.projectiles), shooter.ProjectileCooldown, want)
		}
		if want == 0 {
			break
		}
	}
	b.shoot(0, 1, target)
	if len(b.projectiles) != 2 || shooter.ProjectileCooldown != 8 {
		t.Fatalf("冷卻歸零後應再次發射普通投射物：count=%d cooldown=%d", len(b.projectiles), shooter.ProjectileCooldown)
	}

	shooter.ProjectileCooldown = 0
	shooter.Z = 2
	target.X, target.Y, target.Z = 11, 20, 1
	b.shoot(0, 1, target)
	if len(b.projectiles) != 3 || !b.projectiles[2].special || shooter.ProjectileCooldown != 6 {
		t.Fatalf("sub_1AD7F 成功後投射物／冷卻錯誤：%+v cooldown=%d", b.projectiles, shooter.ProjectileCooldown)
	}
}

func TestRawPlaneHighAndTerrainFlag(t *testing.T) {
	b := newTestBattle(walledField(32))
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, Kind: Infantry, X: 32, Y: 10, Z: 4}
	s.syncTerrain(b.Field, s.X, s.Y, s.Z)
	if s.PlaneHigh != PlaneHighElevated || !s.Climbing {
		t.Fatalf("高位平面欄位錯誤：PlaneHigh=%#x climbing=%v", s.PlaneHigh, s.Climbing)
	}
	if !s.HighTerrain {
		t.Fatal("堆疊 4 層的格子沒有設定 +0x00 bit 1 對應旗標")
	}
	s.syncTerrain(b.Field, 10, 10, 0)
	if s.PlaneHigh != PlaneHighGround || s.Climbing || s.HighTerrain {
		t.Fatalf("回到地面後 raw terrain 沒清除：%+v", *s)
	}
}

func TestSpecialProjectileUsesPlaneHighAndMaxAxisDistance(t *testing.T) {
	attacker := &Soldier{Kind: Infantry, X: 10, Y: 20, Z: 2,
		PlaneHigh: PlaneHighElevated}
	target := &Soldier{Kind: Infantry, X: 12, Y: 22, Z: 0}
	if !specialAttackAvailable(attacker, target) {
		t.Fatal("兩軸差值都是 2 時，原版 sub_1ACA4 應允許特殊投射物")
	}
	target.X, target.Y = 13, 21
	if specialAttackAvailable(attacker, target) {
		t.Fatal("較大軸差值為 3 時，不應進入特殊投射物")
	}
}

func TestNormalProjectileVelocityMatchesRawSub1AD2D(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{3}}, 0)
	shooter := &Soldier{X: 10, Y: 20, Z: 3}
	target := &Soldier{X: 14, Y: 22, Z: 1}

	// max(|dx|, |dy|) = 4；4 >> 1 + (1 - 3) + (3 & 3) = 3，
	// 對應 raw sub_1AD2D 的 `shr bx,1`、高度差、`and ax,3`。
	if got, want := normalProjectileVelocity(b, shooter, target), 3*0x14; got != want {
		t.Fatalf("普通箭初始速度=%#x，應為 %#x", got, want)
	}
}

func TestLockOnPlanePenaltyMatchesRawBranches(t *testing.T) {
	for _, tc := range []struct {
		name string
		me   Soldier
		e    Soldier
		want bool
	}{
		{"地面騎兵看高位平面", Soldier{Kind: Cavalry}, Soldier{PlaneHigh: PlaneHighElevated}, true},
		{"地面步兵看高位平面", Soldier{Kind: Infantry}, Soldier{PlaneHigh: PlaneHighElevated}, false},
		{"高位步兵看普通地面", Soldier{Kind: Infantry, PlaneHigh: PlaneHighElevated}, Soldier{}, false},
		{"高位步兵看高地面旗標", Soldier{Kind: Infantry, PlaneHigh: PlaneHighElevated}, Soldier{HighTerrain: true}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := targetPlanePenalty(&tc.me, &tc.e); got != tc.want {
				t.Fatalf("targetPlanePenalty=%v，預期 %v", got, tc.want)
			}
		})
	}
}

func TestProjectileRawDirectionGridAndHeightPower(t *testing.T) {
	b := newTestBattle(flatField())
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			b.Sides[side].Soldiers[k].Alive = false
		}
	}
	p := projectile{
		side: 0, x: 10, y: 20, z: 1, direction: East, power: arrowPower,
		heightFP: 1 << 8, velocityFP: 0x100,
		gridIndex:         projectileGridIndex(10, 20, 1),
		previousGridIndex: projectileGridIndex(10, 20, 1),
	}
	b.projectiles = []projectile{p}
	b.stepProjectiles()
	if len(b.projectiles) != 1 {
		t.Fatalf("普通投射物在空氣層不應消失：%+v", b.projectiles)
	}
	got := b.projectiles[0]
	if got.x != 11 || got.y != 20 || got.z != 2 {
		t.Fatalf("raw 方向／高度錯誤：%+v", got)
	}
	if got.previousX != 10 || got.previousY != 20 || got.previousZ != 1 {
		t.Fatalf("sub_1BAB7 前一格錯誤：%+v", got)
	}
	if got.gridIndex != projectileGridIndex(11, 20, 2) ||
		got.previousGridIndex != projectileGridIndex(10, 20, 1) {
		t.Fatalf("+0x10／+0x12 格索引錯誤：%+v", got)
	}
	if got.power != arrowPower-arrowPower>>2 || got.velocityFP != 0x100-0x14 {
		t.Fatalf("上升時威力／速度錯誤：power=%d velocity=%#x", got.power, got.velocityFP)
	}
}

func TestProjectileChecksCurrentCellBeforeMoving(t *testing.T) {
	b := newTestBattle(flatField())
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			b.Sides[side].Soldiers[k].Alive = false
		}
	}
	target := &b.Sides[1].Soldiers[1]
	*target = Soldier{Alive: true, Kind: Infantry, HP: MaxHP, X: 10, Y: 20, Z: 1}
	b.projectiles = []projectile{{
		side: 0, x: 10, y: 20, z: 1, direction: East, power: arrowPower,
		heightFP: 1 << 8, velocityFP: 0x100,
		gridIndex: projectileGridIndex(10, 20, 1),
	}}
	b.stepProjectiles()
	if len(b.projectiles) != 0 {
		t.Fatalf("命中目前格後投射物沒有消失：%+v", b.projectiles)
	}
	if target.HP != MaxHP-arrowPower/InfantryArrowDivisor {
		t.Fatalf("目前格命中傷害=%d，應為 %d", MaxHP-target.HP,
			arrowPower/InfantryArrowDivisor)
	}
}

func TestProjectileStopsAtSolidLayerAfterMoving(t *testing.T) {
	b := newTestBattle(walledField(32))
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			b.Sides[side].Soldiers[k].Alive = false
		}
	}
	b.projectiles = []projectile{{
		side: 0, x: 31, y: 10, z: 0, direction: East, power: arrowPower,
		heightFP: 0, velocityFP: 0,
		gridIndex: projectileGridIndex(31, 10, 0),
	}}
	b.stepProjectiles()
	if len(b.projectiles) != 0 {
		t.Fatalf("投射物進入實心城壁層後仍存活：%+v", b.projectiles)
	}
}

// ⭐ 近戰**不吃**那個四分之一。
func TestMeleeIgnoresInfantryArrowResistance(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{100}}, 0)
	b.Advantage[0] = Even
	attacker := &Soldier{Alive: true, Kind: Cavalry, Power: 40, HP: MaxHP,
		Cmd: Attack, Next: Attack}
	inf := &Soldier{Alive: true, Kind: Infantry, HP: MaxHP}
	cav := &Soldier{Alive: true, Kind: Cavalry, HP: MaxHP}
	if !b.meleeHit(0, attacker, inf) {
		t.Fatal("第一下近戰應該命中")
	}
	attacker.HitGeneral = false
	if !b.meleeHit(0, attacker, cav) {
		t.Fatal("第二下近戰應該命中")
	}
	if MaxHP-inf.HP != 40 {
		t.Errorf("步兵近戰掉了 %d，應為 40——四分之一只在飛道具那一支",
			MaxHP-inf.HP)
	}
	if MaxHP-cav.HP != 40 {
		t.Errorf("騎馬近戰掉了 %d，應為 40", MaxHP-cav.HP)
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

	// 城牆本身走不上去——不論兵種。牆頂只能從**門那一格**爬
	// （`sub_1B0D3` 的 `cmp al, 0F0h`，docs/re/63 §4）。
	b := newTestBattle(walledField(32))
	s := &b.Sides[0].Soldiers[1]
	s.Kind = Cavalry
	s.X, s.Y, s.Z = 31, 10, 0
	s.GoalX, s.GoalY, s.GoalZ = 32, 10, 4
	b.moveToward(0, 1)
	if s.X == 32 {
		t.Error("騎馬爬上城牆了")
	}
	s.Kind = Infantry
	s.X, s.Y, s.Z = 31, 10, 0
	b.moveToward(0, 1)
	if s.X == 32 {
		t.Error("步兵也不該直接走上城牆——城壁那一格沒有低平面地面")
	}
}

// 上下城牆只在門那一格，而且大將與騎馬不做（docs/re/63 §4）。
func TestClimbOnlyAtGateCells(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
	b.Deploy(AttackerSide, 0, Infantry, 8)
	b.Place()
	gateY := Height / 2
	if !f.IsGateCell(32, gateY) {
		t.Fatalf("(32,%d) 應該是門格", gateY)
	}
	if f.IsGateCell(32, gateY-1) {
		t.Fatalf("(32,%d) 是城壁，不該是門格", gateY-1)
	}
	top, ok := f.GroundLevel(32, gateY, PlaneHigh)
	if !ok {
		t.Fatal("門格的高平面應該有地面（跟旁邊的牆頂同高）")
	}

	s := &b.Sides[AttackerSide].Soldiers[0]
	s.Kind = Infantry
	s.X, s.Y, s.Z, s.OnWall = 32, gateY, 4, false
	if !b.tryClimb(AttackerSide, 0) {
		t.Fatal("步兵站在門格上應該爬得上去")
	}
	if !s.OnWall || s.Z != top {
		t.Errorf("爬上去之後 OnWall=%v Z=%d，預期 true／%d", s.OnWall, s.Z, top)
	}
	if !b.tryClimb(AttackerSide, 0) {
		t.Error("站在牆頂的門格上應該下得來")
	}
	if s.OnWall {
		t.Error("下來之後 OnWall 還是 true")
	}

	// 騎馬與大將不做 Z 移動。
	for _, kind := range []Kind{Cavalry, General} {
		s.Kind, s.OnWall, s.Z = kind, false, 4
		if b.tryClimb(AttackerSide, 0) {
			t.Errorf("%v 不該爬得上城牆", kind)
		}
	}

	// 不是門的那一格，步兵也上不去。
	s.Kind, s.OnWall = Infantry, false
	s.X, s.Y = 32, gateY-1
	if b.tryClimb(AttackerSide, 0) {
		t.Error("城壁那一格不該爬得上去")
	}
}

func TestHorizontalStepAdjustsOneLevelForNonClimber(t *testing.T) {
	stack := make([][]int, Height)
	for y := range stack {
		stack[y] = make([]int, Width)
	}
	stack[20][11] = 2
	b := NewBattle(NewField(stack, 0), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, Kind: Cavalry, HP: MaxHP, Stamina: StaminaFull,
		Cmd: Attack, Next: Attack, X: 10, Y: 20, Z: 1, GoalX: 11, GoalY: 20, GoalZ: 2}
	b.moveToward(0, 1)
	if s.X != 11 || s.Z != 2 {
		t.Fatalf("水平跨一層未同步高度：位置=%d,%d,%d", s.X, s.Y, s.Z)
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

func TestRetreatDropsOldPath(t *testing.T) {
	b := newTestBattle(flatField())
	s := &b.Sides[0].Soldiers[1]
	s.Cmd, s.Next = Retreat, Retreat
	s.Path = &Waypoints{pts: []Point{{X: 20, Y: 20}}}
	s.PathAt = 123
	b.doRetreat(0, 1)
	if s.Path != nil || s.PathAt != 0 {
		t.Fatal("退卻仍保留舊繞路點")
	}
}

func TestRetreatUsesOriginalExitTarget(t *testing.T) {
	b := newTestBattle(flatField())
	low := &b.Sides[0].Soldiers[1]
	*low = Soldier{Alive: true, Cmd: Retreat, Next: Retreat, X: 20, Y: 5, Z: 3}
	b.doRetreat(0, 1)
	if low.GoalX != 1 || low.GoalY != 0x10 || low.GoalZ != 0 {
		t.Fatalf("低處退卻目標錯誤：%d,%d,%d", low.GoalX, low.GoalY, low.GoalZ)
	}

	high := &b.Sides[1].Soldiers[1]
	*high = Soldier{Alive: true, Cmd: Retreat, Next: Retreat, X: 20, Y: 60, Z: 3}
	b.doRetreat(1, 1)
	if high.GoalX != 0x3E || high.GoalY != 0x2F || high.GoalZ != 0 {
		t.Fatalf("高處退卻目標錯誤：%d,%d,%d", high.GoalX, high.GoalY, high.GoalZ)
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

// 真實 BATTLE.MAP 的攻城地形也不能讓核心戰鬥卡死。
// 原始素材不隨專案散布；沒有使用者提供的 dosv 資產時跳過。
func TestRealSiegeFieldBattleTerminates(t *testing.T) {
	const dir = "../../../workplace/orig/dosv"
	if _, err := os.Stat(dir + "/BATTLE.MAP"); err != nil {
		t.Skip("找不到原版 BATTLE.MAP，跳過")
	}
	read := func(name string) []byte {
		b, err := os.ReadFile(dir + "/" + name)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	lib, err := battle.Parse(read("BATTLE.MAP"), read("BATTLE.MDL"), read("BATTLE.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	forms, err := LoadFormations(dir + "/KI.EXE")
	if err != nil {
		t.Fatal(err)
	}
	const fieldNumber = 56 // 劇本 1 的濮陽，正常 AI 截圖所進的攻城場
	f := NewFieldFromTiles(lib.Tiles(fieldNumber), lib.Heights(fieldNumber), lib.GateX(fieldNumber))
	b := NewBattle(f, forms, &fixedRand{seq: []int{1, 7, 3}}, 127)
	b.Deploy(0, 0, Infantry, 16)
	b.Deploy(1, 0, Infantry, 49)
	b.Place()
	b.Order(0, -1, Attack)
	b.Order(1, -1, Attack)
	if code := lib.Script(0, battle.Category(fieldNumber)); code != nil {
		b.SetScript(1, NewScript(code, 1))
	}
	if !b.Run(200000) {
		t.Fatalf("真實攻城場跑了 20 萬幀還沒結束（攻方剩 %d、守方剩 %d）",
			b.Sides[0].Remaining(), b.Sides[1].Remaining())
	}
	t.Logf("真實攻城場第 %d 幀結束，勝方 %d（剩 %d 對 %d）",
		b.Frame, b.Winner, b.Sides[0].Remaining(), b.Sides[1].Remaining())
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

// 補進場的兵要用**那一隊的兵種**，不能用隊長的。
//
// 第 0 隊的隊長是大將（兵種 0），照隊長補會讓整隊都變成大將——
// 而大將不攻擊、不會陣亡，一整隊站著不動，戰鬥就永遠打不完。
func TestReinforcementsKeepSquadKind(t *testing.T) {
	b := newTestBattle(flatField())
	// 把第 0 隊打光，讓待機的補上來。
	for i := 0; i < PerSquad; i++ {
		b.Sides[0].Soldiers[i].Alive = false
	}
	for i := 0; i < 40; i++ {
		b.reinforce()
	}
	generals := 0
	for i := 0; i < PerSquad; i++ {
		s := &b.Sides[0].Soldiers[i]
		if !s.Alive {
			continue
		}
		if s.Kind == General {
			generals++
		} else if s.Kind != Infantry {
			t.Errorf("第 0 隊補進來的兵種是 %v，應為步兵", s.Kind)
		}
	}
	if generals > 0 {
		t.Errorf("補進來 %d 個大將——補兵抄了隊長的兵種", generals)
	}
}

// ⭐ 撞到自己人是**對調位置**，不是擋下來也不是穿過去（`sub_1B732`）。
func TestFriendlyCollisionSwaps(t *testing.T) {
	b := newTestBattle(flatField())
	a := &b.Sides[0].Soldiers[10] // 不是隊長，所以不是大將
	c := &b.Sides[0].Soldiers[11]
	a.Kind, c.Kind = Infantry, Infantry
	a.X, a.Y, a.Z = 20, 20, 0
	c.X, c.Y, c.Z = 21, 20, 0
	a.Cmd, c.Cmd = Attack, Attack

	if ok, _ := b.tryMove(0, 10, 21, 20, 0); !ok {
		t.Fatal("撞到自己人應該換位成功")
	}
	if a.X != 21 || c.X != 20 {
		t.Errorf("換位之後 a.X=%d c.X=%d，應為 21 與 20", a.X, c.X)
	}
	if !c.Swapped {
		t.Error("被換的那一個要標記（原版 `or byte ptr [di], 40h`）")
	}

	// 標記的作用：**同一幀裡第三個人不能再跟它換**
	// （`test byte ptr [di], 61h` 的 bit 6）。
	d := &b.Sides[0].Soldiers[12]
	d.Kind = Infantry
	d.X, d.Y, d.Z = 19, 20, 0
	if ok, _ := b.tryMove(0, 12, 20, 20, 0); ok {
		t.Error("這一幀已經被換過的兵，不該再被第三個人換走")
	}
	// 那個兵自己更新時會把旗標清掉（`and byte ptr [si], 0BFh`）。
	b.updateSoldier(0, 11)
	if c.Swapped {
		t.Error("兵自己更新時應該清掉「被換過」的旗標")
	}
}

// 撞到敵人是**打他**，自己不動（`loc_1B5A1` → `sub_1B618`）。
func TestEnemyCollisionAttacks(t *testing.T) {
	b := newTestBattle(flatField())
	b.Advantage[0] = Even
	a := &b.Sides[0].Soldiers[10]
	e := &b.Sides[1].Soldiers[10]
	a.Kind, e.Kind = Infantry, Infantry
	a.X, a.Y, a.Z = 20, 20, 0
	e.X, e.Y, e.Z = 21, 20, 0
	hp := e.HP

	if ok, _ := b.tryMove(0, 10, 21, 20, 0); ok {
		t.Error("敵人擋著不該走得過去")
	}
	if a.X != 20 {
		t.Errorf("撞到敵人卻移動了，X ＝ %d", a.X)
	}
	if e.HP >= hp {
		t.Errorf("撞到敵人沒扣血（%d → %d）", hp, e.HP)
	}
}

// 大將換不了位置（`cmp byte ptr [di+4], 0 / jz` ＝ 失敗）。
func TestGeneralNeverSwaps(t *testing.T) {
	b := newTestBattle(flatField())
	g := &b.Sides[0].Soldiers[0] // 第 0 隊隊長 ＝ 大將
	a := &b.Sides[0].Soldiers[10]
	a.Kind = Infantry
	g.X, g.Y, g.Z = 21, 20, 0
	a.X, a.Y, a.Z = 20, 20, 0
	if ok, _ := b.tryMove(0, 10, 21, 20, 0); ok {
		t.Error("不該跟大將換位置")
	}
	if a.X != 20 || g.X != 21 {
		t.Error("大將被換走了")
	}
}

// ⭐ 六個指令下下去，部隊要真的動。
//
// 這一條釘的是「戰術畫面能不能指揮」，不是「攻城打不打得下來」——
// 破城要靠玩家換陣形、挑時機（說明書第 11 章），不是 AI 自己會。
func TestSixCommandsMakeSoldiersAct(t *testing.T) {
	type outcome struct {
		moved, attacked, retreating bool
	}
	run := func(t *testing.T, cmd Command, siege bool) outcome {
		t.Helper()
		f := walledField(32)
		if !siege {
			f = flatField()
		}
		b := newTestBattle(f)
		// Place() 已經把兵擺在陣形位置上，所以「陣形」「守陣」下去
		// 本來就不必動。先把他們往前推三格，才看得出有沒有回位。
		for k := range b.Sides[AttackerSide].Soldiers {
			if s := &b.Sides[AttackerSide].Soldiers[k]; s.Alive {
				s.X = clamp(s.X + 3)
			}
		}
		b.Order(AttackerSide, -1, cmd)
		b.Order(DefenderSide, -1, Guard)
		before := make([][3]int, SoldiersOnFoot)
		for k := range b.Sides[AttackerSide].Soldiers {
			s := &b.Sides[AttackerSide].Soldiers[k]
			before[k] = [3]int{s.X, s.Y, s.Z}
		}
		var out outcome
		for i := 0; i < 400 && !b.Done; i++ {
			b.Step()
			for k := range b.Sides[AttackerSide].Soldiers {
				s := &b.Sides[AttackerSide].Soldiers[k]
				if !s.Alive {
					continue
				}
				if [3]int{s.X, s.Y, s.Z} != before[k] {
					out.moved = true
				}
				if s.Cmd == Retreat {
					out.retreating = true
				}
			}
			if b.Sides[DefenderSide].Remaining() < Squads*PerSquad {
				out.attacked = true
			}
		}
		return out
	}

	for _, tc := range []struct {
		cmd   Command
		siege bool
		want  string
	}{
		{Form, false, "moved"},
		{Attack, false, "moved"},
		{Charge, false, "moved"},
		{ScaleWal, true, "moved"},
		{Guard, false, "moved"},
		{Retreat, false, "retreating"},
	} {
		got := run(t, tc.cmd, tc.siege)
		switch tc.want {
		case "moved":
			if !got.moved {
				t.Errorf("命令「%v」下了 400 幀，部隊一格都沒動", tc.cmd)
			}
		case "retreating":
			if !got.retreating {
				t.Errorf("命令「%v」下了 400 幀，沒有任何兵在退卻", tc.cmd)
			}
		}
	}
}

// 陣形與陣形線是玩家可以改的（說明書 4.4／4.5，docs/spec/37）。
// 改完之後下「陣形」令，部隊要走到新的位置。
func TestFormationAndLineChangeWhereSoldiersGather(t *testing.T) {
	gather := func(formation, line int) (sumX int) {
		b := newTestBattle(flatField())
		b.Sides[AttackerSide].Formation = formation
		b.Sides[AttackerSide].Line = LineFor(AttackerSide, line)
		b.Order(AttackerSide, -1, Form)
		for i := 0; i < 600 && !b.Done; i++ {
			b.Step()
		}
		for k := range b.Sides[AttackerSide].Soldiers {
			if s := &b.Sides[AttackerSide].Soldiers[k]; s.Alive {
				sumX += s.X
			}
		}
		return sumX
	}
	own, centre := gather(0, 0), gather(0, 1)
	if own >= centre {
		t.Errorf("陣形線「自軍側」的 X 總和 %d 應該小於「中央」的 %d", own, centre)
	}
	if a, b2 := gather(0, 1), gather(5, 1); a == b2 {
		t.Errorf("換陣形之後部隊位置沒變（兩次都是 %d）", a)
	}
}

// TestSetPlayerSideFollowsPlayerNotAttacker 釘住 docs/spec/56 §2.1：
// **陣形原點與鏡射跟著「玩家／腳本」走，不跟「攻方／守方」走。**
//
// 原版在玩家守方時互換 `word_10D2E`／`word_10D30`，互換之後 side 0 永遠是
// 玩家；而原點是 `word_1D33C`（玩家，X=5）與 `word_1D33E`（腳本，X=58）。
// remake 的 `Sides[0]` 固定是攻方，所以玩家守城時要把這兩樣換過來。
//
// ⭐ 這一條擋的是「照 side 索引取值」那個很自然的寫法——它在玩家攻城時
// 完全正確，只有玩家守城時錯，而守城的畫面平常不會拿來對拍。
func TestSetPlayerSideFollowsPlayerNotAttacker(t *testing.T) {
	player, script := LineFor(0, 0), LineFor(1, 0)
	if player == script {
		t.Fatal("玩家與腳本的陣形原點相同，這個測試沒有鑑別力")
	}
	for _, tc := range []struct {
		name string
		side int
	}{{"玩家攻城", AttackerSide}, {"玩家守城", DefenderSide}} {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBattle(nil, SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
			b.SetPlayerSide(tc.side)
			if got := b.Sides[tc.side].Line; got != player {
				t.Errorf("玩家那一側的原點是 %d，應為 %d", got, player)
			}
			if got := b.Sides[1-tc.side].Line; got != script {
				t.Errorf("腳本那一側的原點是 %d，應為 %d", got, script)
			}
			if b.Sides[tc.side].Mirror {
				t.Error("玩家那一側不該鏡射")
			}
			if !b.Sides[1-tc.side].Mirror {
				t.Error("腳本那一側要鏡射")
			}
		})
	}
	// 沒呼叫過就等同「玩家是攻方」——NewBattle 的預設。
	b := NewBattle(nil, SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	if b.PlayerSide != AttackerSide || b.Sides[0].Line != player {
		t.Error("預設不是「玩家在攻方」")
	}
}

// TestSwappedSoldierSkipsItsTurn 釘住「被換位的兵這一幀不動」。
//
// 原版 `sub_1ADC8` 的 `0001ADED test al, 40h / jnz loc_1AE26`：
// 旗標立著就直接跳到重畫，不移動也不攻擊，旗標在那裡才清掉
// （docs/spec/62）。少了它，被推開的兵下一幀馬上又往前擠，
// 前排那一格一直換人，圍著打卻打不到（docs/spec/61 §5.1）。
func TestSwappedSoldierSkipsItsTurn(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1, 7, 3}}, 0)
	s := &b.Sides[0].Soldiers[0]
	*s = Soldier{Alive: true, Kind: Infantry, HP: MaxHP, Stamina: StaminaFull,
		X: 10, Y: 20, Z: 0, Target: -1, Cmd: Form, Next: Form}
	s.GoalX, s.GoalY = 14, 20
	s.Swapped = true

	b.updateSoldier(0, 0)
	if s.X != 10 || s.Y != 20 {
		t.Fatalf("被換位的兵這一幀不該移動，卻走到 (%d,%d)", s.X, s.Y)
	}
	if s.Swapped {
		t.Fatal("跳過那一幀之後旗標要清掉，否則它會永遠不動")
	}

	// 旗標清掉之後，下一幀照常走。
	b.updateSoldier(0, 0)
	if s.X == 10 && s.Y == 20 {
		t.Fatal("旗標清掉後應該恢復移動")
	}
}

// TestDeployStartsHPAtMorale 釘住「開場體力 ＝ 軍團士氣」（docs/spec/61）。
func TestDeployStartsHPAtMorale(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	b.Sides[0].Morale = 200
	b.Deploy(0, 1, Infantry, 8)
	if got := b.Sides[0].Soldiers[PerSquad].HP; got != 200 {
		t.Fatalf("士氣 200 的隊開場體力 = %d，預期 200", got)
	}
	if MaxHP != 100 {
		t.Fatalf("MaxHP = %d：那是 sub_1B97E 的回復上限，不該跟著開場值走", MaxHP)
	}

	// 沒給士氣時退回預設，否則兵一上場就是 0 點。
	b2 := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	b2.Deploy(0, 1, Infantry, 8)
	if got := b2.Sides[0].Soldiers[PerSquad].HP; got != DefaultPower {
		t.Fatalf("沒給士氣時開場體力 = %d，預期 %d", got, DefaultPower)
	}
}

// TestReinforcementUsesMoraleHP 釘住增援與開場走同一個值。
func TestReinforcementUsesMoraleHP(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	b.Sides[0].Morale = 150
	b.Deploy(0, 1, Infantry, 20) // 8 上場、12 待機
	b.Place()
	k := PerSquad + 3
	b.Sides[0].Soldiers[k].Alive = false
	b.reinforce()
	if got := b.Sides[0].Soldiers[k].HP; got != 150 {
		t.Fatalf("補進來的兵體力 = %d，預期與開場同值 150", got)
	}
}

// TestHitStunSkipsThreeFrames 釘住挨打之後的硬直（docs/spec/63）。
//
// 原版 `sub_1B618` 命中時把 `+0x01` 寫成 2 並設 `+0x00` bit 6，
// 兩者接力擋掉三幀；`Hurt` 撐到硬直歸零那一幀才清，
// 因為它同時是換位的擋條件（`docs/re/11` §5.16）。
func TestHitStunSkipsThreeFrames(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1, 7, 3}}, 0)
	s := &b.Sides[0].Soldiers[0]
	*s = Soldier{Alive: true, Kind: Infantry, HP: MaxHP, Stamina: StaminaFull,
		X: 10, Y: 20, Z: 0, Target: -1, Cmd: Form, Next: Form}
	s.GoalX, s.GoalY = 14, 20

	b.applyHit(s, 1)
	if !s.Swapped || s.Stun != HitStunFrames || !s.Hurt {
		t.Fatalf("挨打後 = 不動:%v 硬直:%d 受擊:%v", s.Swapped, s.Stun, s.Hurt)
	}

	for frame := 1; frame <= 3; frame++ {
		b.updateSoldier(0, 0)
		if s.X != 10 || s.Y != 20 {
			t.Fatalf("硬直第 %d 幀就動了：(%d,%d)", frame, s.X, s.Y)
		}
		if frame < 3 && !s.Hurt {
			t.Fatalf("硬直還沒結束（第 %d 幀）就清掉受擊旗標，換位會擋不住", frame)
		}
	}
	if s.Hurt {
		t.Fatal("硬直歸零那一幀要清掉受擊旗標")
	}
	b.updateSoldier(0, 0)
	if s.X == 10 && s.Y == 20 {
		t.Fatal("第四幀應該恢復移動")
	}
}

// TestRetreatedSoldiersCountAsSurvivors 釘住「退卻算生還、戰死不算」。
//
// 原版兩種離場走同一支 `sub_1B4B8`，差別只有 `ah`：退卻那條
// （`sub_1AAED`）傳 0 會把該隊的存活數 +1，倒地數完四幀那條傳 1 不加
// （docs/spec/65）。打完的寫回是「Σ（存活 ＋ 待機）」。
//
// ⚠ 生還**不能**進 `Remaining()`：那個是勝負判定用的「還補得出兵嗎」。
func TestRetreatedSoldiersCountAsSurvivors(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	b.Sides[0].Morale = 100
	b.Deploy(0, 0, Infantry, 8)
	base := b.Sides[0].Remaining()

	// 一個兵走到畫面邊緣：離場，但算生還。
	s := &b.Sides[0].Soldiers[1]
	s.Cmd, s.Next = Retreat, Retreat
	s.X, s.Y = MinCoord, 0x20
	b.doRetreat(0, 1)
	if s.Alive {
		t.Fatal("走到邊緣應該離場")
	}
	if got := b.Sides[0].Remaining(); got != base-1 {
		t.Fatalf("離場之後 Remaining = %d，預期 %d（生還不能算進去）", got, base-1)
	}
	if got := b.Sides[0].Survivors(); got != base {
		t.Fatalf("退卻的兵要算生還：Survivors = %d，預期 %d", got, base)
	}

	// 戰死：兩邊都不算。
	d := &b.Sides[0].Soldiers[2]
	d.Alive = false
	if got := b.Sides[0].Survivors(); got != base-1 {
		t.Fatalf("戰死的不該算生還：Survivors = %d，預期 %d", got, base-1)
	}
}
