package tactical

import "sort"

// 傷害與飛道具。公式出自 `sub_1B97E`（命中）與 `sub_1BAB7`（威力隨高度變化），
// 見 docs/re/11 §5.3–§5.5。

const (
	// arrowPower 是 `sub_1AD2D` 建立的普通飛道具初始威力。
	// `sub_1B8AA` 把呼叫端的 CH（1Ch）寫入飛道具 +0x04。
	arrowPower = 0x1C

	// specialProjectilePower 是 `sub_1AD7F` 的 CH=0x20。
	specialProjectilePower = 0x20
)

// hitByArrow 是**飛道具**的傷害（原版 `sub_1B97E`）。
//
// 目標為步兵時傷害右移兩位；大將另有姿態與最低體力閘。
// （`cmp byte ptr [bx+4], 36h / jz` → `shr al,1` × 2）。
// 大將也吃這個減傷，但要多一個條件（`test byte ptr [bx+2], 1`）。
func (b *Battle) hitByArrow(from int, e *Soldier, power int) bool {
	if e == nil || !e.Alive {
		return false
	}
	if e.IsGeneral() && e.PoseStep&1 == 0 {
		return false
	}
	if e.Kind == Infantry || e.IsGeneral() {
		power /= InfantryArrowDivisor
	}
	return b.projectileHit(e, power)
}

// InfantryArrowDivisor 是步兵挨箭的減傷倍數。
//
// 原版 `sub_1B97E`：`cmp byte ptr [bx+4], 36h / jz` → **威力右移兩位**。
// 0x36 ＝ 54 ＝ 步兵。說明書「攻城戦では弓兵、**歩兵**が必要です」的
// 數值依據就在這裡：步兵是唯一扛得住城牆上箭雨的兵種。
const InfantryArrowDivisor = 4

// HitStunFrames 是挨打之後的硬直幀數（原版 `sub_1B618` 的
// `mov byte ptr [di+1], 2`，docs/spec/63）。
//
// 它與同時設下的「這一幀不動」旗標接力：那一幀 ＋ 這裡的兩幀
// ＝ 一共三幀站著不動，第四幀才恢復。
const HitStunFrames = 2

// meleeHit 是一般兵士碰撞命中的處理（原版 `sub_1B618`）。
//
// 原版不是固定威力：先取 `rand & 0x7F`，加上攻擊者 `+0x18` 戰力
// （8-bit）；結果小於 0x46 就未命中。攻擊方有利時把傷害戰力飽和
// 加 0x40；不利時只把命中值減 0x32。突擊（生效命令 2）再把傷害
// 戰力飽和加 0xC8。步兵的四分之一減傷只在 `sub_1B97E` 的飛道具
// 路徑，這裡不套用。
func (b *Battle) meleeHit(side int, attacker, e *Soldier) bool {
	if attacker == nil || e == nil || !attacker.Alive || !e.Alive {
		return false
	}

	power := clampByte(attacker.Power)
	chance := ((b.rng.Next() & 0x7F) + power) & 0xFF
	adv := b.attackerAdvantage(side)
	if adv == Disadvantaged {
		chance -= 0x32
		if chance < 0 {
			chance = 0
		}
	}

	damage := power
	if adv == Advantaged {
		damage = saturatingByteAdd(damage, 0x40)
	}
	if attacker.Cmd == Charge {
		damage = saturatingByteAdd(damage, 0xC8)
	}

	// `sub_1B618` 的尾端不論是否命中都會留下攻擊旗標。
	attacker.HitGeneral = true
	if chance < 0x46 {
		return false
	}
	return b.applyHit(e, damage)
}

// attackCollision 統一一般兵與敵方大將的兩條原版碰撞分支。
func (b *Battle) attackCollision(side int, attacker, e *Soldier) bool {
	if e == nil || !e.Alive {
		return false
	}
	if e.IsGeneral() {
		return b.hitGeneral(attacker, e)
	}
	return b.meleeHit(side, attacker, e)
}

// applyHit 是一般命中後的共用受擊／扣血流程。
func (b *Battle) applyHit(e *Soldier, damage int) bool { return b.applyDamage(e, damage, true) }

func (b *Battle) applyDamage(e *Soldier, damage int, notifyLeader bool) bool {
	if e == nil || !e.Alive {
		return false
	}

	before := snapshotDeathPose(e)
	// 原版受擊同時設 +0x02 bit 4、面向歸零、+0x00 bit 6，
	// 以及 `+0x01` 的硬直計時（docs/spec/63）。
	e.HitGeneral = false
	e.Hurt = true
	e.Facing = West
	e.Swapped = true
	e.Stun = HitStunFrames

	// 扣血前處理退卻：體力滿 100 的不排退卻；已經是退卻中的也不改。
	if !e.IsGeneral() && e.HP < MaxHP && e.Cmd != Retreat && e.Next != Retreat {
		if e.Cmd != e.Next {
			e.Cmd = e.Next
		}
		e.Next = Retreat
	}

	e.HP -= damage
	if e.HP > 0 {
		return true
	}
	if e.IsGeneral() {
		e.HP = GeneralMinHP
		return true
	}
	b.beginDeath(e, before)
	if notifyLeader {
		b.squadLeaderGoneFor(e)
	}
	return true
}

// squadLeaderGoneFor 把指標回查到戰場中的兵。投射物路徑只有指標，
// 因此在非一般測試兵上仍能套用隊長退卻規則。
func (b *Battle) squadLeaderGoneFor(e *Soldier) {
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			if &b.Sides[side].Soldiers[k] == e {
				b.squadLeaderGone(side, k)
				return
			}
		}
	}
}

func (b *Battle) projectileHit(e *Soldier, damage int) bool {
	return b.applyDamage(e, damage, false)
}

// attackerAdvantage 把原版唯一的 byte_1D31E 轉成攻擊者一側的結果。
// `Advantage[0]` 是原版的全域結果；另一側正好相反。
func (b *Battle) attackerAdvantage(side int) Advantage {
	global := b.Advantage[0]
	if global == Even {
		return Even
	}
	if side == 0 {
		return global
	}
	if global == Advantaged {
		return Disadvantaged
	}
	return Advantaged
}

func saturatingByteAdd(a, b int) int {
	return clampByte(a + b)
}

// shoot 讓弓兵放一箭。
func (b *Battle) shoot(side, k int, e *Soldier) {
	s := &b.Sides[side].Soldiers[k]
	if s.Z <= e.Z {
		b.setArcherAimGoal(s, e)
	}
	if !projectileReady(s) {
		return
	}
	// Archer（+0x04=0x24）的 sub_1ABFF：高於目標且最大軸距不超過
	// 一格時走 sub_1AD7F；不能用曼哈頓距離，斜角相鄰同樣是 1。
	if s.Kind == Archer && s.Z > e.Z && maxAxisDistance(s, e) <= 1 {
		b.launchSpecialProjectile(side, s, e)
		return
	}
	b.launchNormalProjectile(side, s, e)
}

func projectileReady(s *Soldier) bool {
	if s == nil {
		return false
	}
	if s.ProjectileCooldown == 0 {
		return true
	}
	s.ProjectileCooldown--
	return false
}

func maxAxisDistance(s, e *Soldier) int {
	if s == nil || e == nil {
		return 0
	}
	dx, dy := abs(s.X-e.X), abs(s.Y-e.Y)
	if dy > dx {
		return dy
	}
	return dx
}

func (b *Battle) launchNormalProjectile(side int, s, e *Soldier) {
	if b == nil || s == nil || e == nil {
		return
	}
	direction := attackDirection(s.X, s.Y, e.X, e.Y)
	s.Facing = direction
	z := s.Z + 1
	velocity := normalProjectileVelocity(b, s, e)
	slot, owned := b.projectileSlot(side, s)
	if owned && b.projectileSlotBusy(slot) {
		return
	}
	b.emitSFX(SFXNormalLaunch)
	b.projectiles = append(b.projectiles, projectile{
		slot: slot, slotOwned: owned, side: side, x: s.X, y: s.Y, z: z,
		tx: e.X, ty: e.Y, tz: e.Z, power: arrowPower,
		direction: direction, gridIndex: projectileGridIndex(s.X, s.Y, z),
		previousGridIndex: projectileGridIndex(s.X, s.Y, z),
		previousX:         s.X, previousY: s.Y, previousZ: z,
		heightFP: z << 8, velocityFP: velocity,
	})
	// sub_1AD2D 成功配置 +0x1400 區的投射物後才寫入 8；配置失敗時
	// 原版不改冷卻；普通與特殊投射物共用該兵對應的固定槽。
	s.finishProjectileLaunch(8)
}

// shootSpecial 重現 `sub_1AD7F` 的 CH=0x20 分支。
//
// 這不是另一種會沿 X/Y 飛很遠的箭：呼叫端把方向加上 0x80，
// `sub_1BA2E` 因此不再移動 X/Y，只在攻擊者朝向的相鄰格做短促的
// 垂直效果；建立時的威力是 0x20，初始垂直速度是 -0x100。
func (b *Battle) shootSpecial(side, k int, e *Soldier) {
	s := &b.Sides[side].Soldiers[k]
	if !projectileReady(s) {
		return
	}
	b.launchSpecialProjectile(side, s, e)
}

func (b *Battle) launchSpecialProjectile(side int, s, e *Soldier) {
	if b == nil || s == nil || e == nil {
		return
	}
	direction := attackDirection(s.X, s.Y, e.X, e.Y)
	s.Facing = direction
	x, y := s.X, s.Y
	switch direction {
	case West:
		x--
	case North:
		y--
	case East:
		x++
	case South:
		y++
	}
	if !inBounds(x, y) {
		return
	}
	slot, owned := b.projectileSlot(side, s)
	if owned && b.projectileSlotBusy(slot) {
		return
	}
	b.emitSFX(SFXSpecialLaunch)
	b.projectiles = append(b.projectiles, projectile{
		slot: slot, slotOwned: owned, side: side, x: x, y: y, z: s.Z + 1,
		tx: e.X, ty: e.Y, tz: e.Z, power: specialProjectilePower,
		special: true, specialFrame: s.PoseStep & 1, direction: direction | 0x80,
		gridIndex:         projectileGridIndex(x, y, s.Z+1),
		previousGridIndex: projectileGridIndex(x, y, s.Z+1),
		previousX:         x, previousY: y, previousZ: s.Z + 1,
		heightFP: (s.Z + 1) << 8, velocityFP: -0x100,
	})
	// sub_1AD7F 成功配置後寫入兵記錄 +0x13=6。
	s.finishProjectileLaunch(6)
}

// specialAttackAvailable 對應 `sub_1AC55` 的 raw `+0x1E` 比較；
// 只有目前兵在較高位平面、且兩軸差值的較大者不超過原版回傳的 2
// 才會進入 `sub_1AD7F`。
func specialAttackAvailable(s, e *Soldier) bool {
	return s != nil && e != nil && s.Kind == Infantry &&
		s.planeHigh() > e.planeHigh() && maxAxisDistance(s, e) <= 2
}

func attackDirection(x, y, tx, ty int) int {
	dx, dy := abs(x-tx), abs(y-ty)
	if dx >= dy {
		if x >= tx {
			return West
		}
		return East
	}
	if y > ty {
		return North
	}
	return South
}

// stepProjectiles 重現 `sub_1B941` 的投射物順序：先檢查目前格，再移動，
// 最後套用 `sub_1BAB7` 的高度威力變化。固定槽升序；中間清旗標繪製見spec159限制。
func (b *Battle) stepProjectiles() {
	sort.SliceStable(b.projectiles, func(i, j int) bool {
		a, c := b.projectiles[i], b.projectiles[j]
		if a.slotOwned != c.slotOwned {
			return a.slotOwned
		}
		return a.slotOwned && a.slot < c.slot
	})
	g := b.unitCollision
	if g == nil {
		g = &collisionGrid{}
		for _, side := range [2]int{b.PlayerSide, 1 - b.PlayerSide} {
			for k, s := range b.Sides[side].Soldiers {
				if s.Alive {
					g.drawSlot(b.collisionSlot(side, k), s.X, s.Y, s.Z, int(s.planeHigh())>>4)
				}
			}
		}
	}
	live := b.projectiles[:0]
	for _, p := range b.projectiles {
		if b.projectileBlocked(p) {
			continue
		}
		cell := p.gridIndex
		side, k := b.collisionOwner(g.atLayer(cell%Width, cell/Width%Height, cell/(Width*Height)))
		if k >= 0 && side != p.side {
			e := &b.Sides[side].Soldiers[k]
			if e.Alive {
				if b.hitByArrow(p.side, e, p.power) {
					b.emitSFX(SFXProjectileHit)
				}
				continue
			}
		}
		if !b.advanceProjectile(&p) {
			continue
		}
		live = append(live, p)
	}
	b.projectiles = live
}

// advanceProjectile 對應 `sub_1BA2E` 加 `sub_1BAB7`。
func (b *Battle) advanceProjectile(p *projectile) bool {
	p.previousX, p.previousY, p.previousZ = p.x, p.y, p.z
	p.previousGridIndex = p.gridIndex

	if p.direction < 0x80 {
		p.x = stepProjectileX(p.x, p.direction)
		p.y = stepProjectileY(p.y, p.direction)
	}
	step := p.velocityFP
	if step > 0x100 {
		step = 0x100
	}
	if step < -0x100 {
		step = -0x100
	}
	p.heightFP += step
	if p.heightFP < 0 {
		return false
	}
	p.velocityFP -= 0x14
	p.z = p.heightFP >> 8
	if p.z > Levels-2 {
		p.z = Levels - 2
		// 1BA96寫回高度高byte；只夾顯示Z會讓固定點高度繼續累積。
		p.heightFP = p.z<<8 | p.heightFP&0xff
	}
	p.gridIndex = projectileGridIndex(p.x, p.y, p.z)
	if b.projectileBlocked(*p) {
		return false
	}

	// `sub_1BAB7` 比較新舊整數 Z：上升扣 power/4，下降加 power/4+1。
	if p.z != p.previousZ {
		q := p.power >> 2
		if p.z > p.previousZ {
			p.power -= q
		} else {
			p.power += q + 1
		}
	}
	p.power &= 0xff
	return true
}

func (b *Battle) projectileBlocked(p projectile) bool {
	if p.gridIndex < 0 || p.gridIndex >= Levels*Width*Height {
		return true
	}
	x, y, z := p.gridIndex%Width, p.gridIndex/Width%Height, p.gridIndex/(Width*Height)
	if !inBounds(x, y) {
		return true
	}
	return b.Field != nil && b.Field.solid[z][y][x]
}

func projectileGridIndex(x, y, z int) int {
	return z*0x1000 + y*0x40 + x
}

func stepProjectileX(x, direction int) int {
	switch direction & 3 {
	case West:
		return x - 1
	case East:
		return x + 1
	default:
		return x
	}
}

func stepProjectileY(y, direction int) int {
	switch direction & 3 {
	case North:
		return y - 1
	case South:
		return y + 1
	default:
		return y
	}
}

// normalProjectileVelocity 重現 `sub_1AD2D`：兩軸最大距離右移一位，
// 加上目標／射手高度差與 `sub_1ECE0 & 3`，最後乘 0x14。`b.rng` 使用
// internal/rules/rng 對應 `sub_1ECE0`／`sub_1EC82` 的同一個取數與播種演算法；
// 逐幀 `sub_1BA2E` 的夾限與 −0x14 衰減也已接入。
func normalProjectileVelocity(b *Battle, s, e *Soldier) int {
	dx, dy := abs(s.X-e.X), abs(s.Y-e.Y)
	if dy > dx {
		dx = dy
	}
	base := dx/2 + e.Z - s.Z
	jitter := 0
	if b != nil && b.rng != nil {
		jitter = b.rng.Next() & 3
	}
	return (base + jitter) * 0x14
}

func sign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

// 撞到**敵方大將**的處理，重現 `sub_1B6BC`（`loc_1B5B1` 呼叫的那一支）。
//
// 大將不是「打不到」——是**要先過一關**：
//
//	大將體力 ≤ 1        → 打不動
//	大將的命令是 0 或 5  → 打不到（陣形中／退卻中）
//	亂數 < 0x19          → **直接命中**（25/256 ≈ 10%）
//	否則                 → 差值 ＝ 攻擊者戰力 − 大將戰力（負的當 0、上限 24）
//	                       亂數 & 0x7F **小於**差值才命中
//
// 命中的傷害是**攻擊者戰力 ÷ 8**（至少 1），而且**大將體力最低留 1**
// （`sub [di+3], al / ja / mov byte ptr [di+3], 1`）——大將不會被打死，
// 但會被打到 50 以下觸發全軍退卻（§5.8h）。
//
// ⭐ 最後 `or byte ptr [si+2], 8` 給**攻擊者**設 bit 3——
// 那正是圖號公式裡剩下的那一位（§5.13）：**打在大將身上會換一張圖**。
const (
	// generalAutoHit 是直接命中的亂數門檻（`cmp al, 19h / jb`）。
	generalAutoHit = 0x19
	// generalEdgeCap 是戰力差的上限（`cmp ah, 18h`）。
	generalEdgeCap = 24
	// generalDamageShift 是傷害的位移（`shr al, 1` × 3 ＝ ÷ 8）。
	generalDamageShift = 3
	// GeneralMinHP 是大將被打到的下限——**不會死**。
	GeneralMinHP = 1
)

// hitGeneral 打一次敵方大將。回傳有沒有命中。
func (b *Battle) hitGeneral(attacker *Soldier, g *Soldier) bool {
	if attacker == nil || g == nil || !attacker.Alive {
		return false
	}
	// `sub_1B6BC` 的尾端不論命中與否都設攻擊者的圖號旗標。
	attacker.HitGeneral = true
	if !g.Alive || g.HP <= GeneralMinHP {
		return false
	}
	// 陣形中或退卻中的大將打不到。
	if g.Cmd == Form || g.Cmd == Retreat {
		return false
	}
	if b.rng.Next()&0xFF >= generalAutoHit {
		edge := attacker.Power - g.Power
		if edge < 0 {
			edge = 0
		}
		if edge > generalEdgeCap {
			edge = generalEdgeCap
		}
		if b.rng.Next()&0x7F >= edge {
			return false
		}
	}
	g.HitGeneral = false
	g.Hurt, g.Swapped, g.Facing = true, true, West
	g.Stun = HitStunFrames
	dmg := attacker.Power >> generalDamageShift
	if dmg < 1 {
		dmg = 1
	}
	g.HP -= dmg
	if g.HP < GeneralMinHP {
		g.HP = GeneralMinHP
	}
	return true
}

// setArcherAimGoal 對應1ACD6的四方向目標調整，不改中繼點。
func (b *Battle) setArcherAimGoal(s, e *Soldier) {
	x, y := e.X, e.Y
	margin := 10
	if dz := e.Z - s.Z; dz >= 0 {
		margin += dz
	} else {
		margin = 0
	}
	switch attackDirection(s.X, s.Y, e.X, e.Y) {
	case West:
		if x <= 54 {
			x = (x + margin) & 255
		}
	case North:
		if y <= 54 {
			y = (y + margin) & 255
		}
	case East:
		if x >= 10 {
			x = (x - margin) & 255
		}
	case South:
		if y >= 10 {
			y = (y - margin) & 255
		}
	}
	s.GoalX, s.GoalY = x, y
}

// finishProjectileLaunch 對應普通／特殊發射的共同成功出口。
func (s *Soldier) finishProjectileLaunch(cooldown uint8) {
	s.Swapped = true
	s.Hurt = false
	s.HitGeneral = true
	s.Stun = 2
	s.ProjectileCooldown = cooldown
}

func (b *Battle) projectileSlot(side int, s *Soldier) (int, bool) {
	for k := range b.Sides[side].Soldiers {
		if &b.Sides[side].Soldiers[k] == s {
			return b.Role(side)*16 + k%16, true
		}
	}
	return 0, false
}
func (b *Battle) projectileSlotBusy(slot int) bool {
	for _, p := range b.projectiles {
		if p.slotOwned && p.slot == slot {
			return true
		}
	}
	return false
}
