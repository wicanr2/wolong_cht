package tactical

// 一個兵一幀做的事：鎖敵 → 套用新命令 → 依命令決定目標 → 移動 → 疲勞。
// 順序照原版 `sub_1ADC8` 的迴圈（先 `sub_1A85B` 鎖敵，再跑命令）。

// lockOnNearest 重現 `sub_1A85B`：掃對方那 48 個，用**曼哈頓距離**挑最近的。
//
// 說明書 11.4「中央突破」之所以有效——用窄陣形讓幾個兵鑽過敵陣直取主將——
// 就是因為**鎖敵只看距離，不看陣線**。
func (b *Battle) lockOnNearest(side, k int) {
	me := &b.Sides[side].Soldiers[k]
	foe := &b.Sides[1-side]
	best, bestD := -1, 1<<30
	for i := range foe.Soldiers {
		e := &foe.Soldiers[i]
		if !e.Alive {
			continue
		}
		d := abs(me.X-e.X) + abs(me.Y-e.Y)
		// `sub_1A85B` 先比較雙方 +0x1E，再依攻擊者兵種與候選的
		// +0x00 bit 1 決定是否加 0x40；不是單純比較 Z 高低。
		if targetPlanePenalty(me, e) {
			d += Width
		}
		if d < bestD {
			best, bestD = i, d
		}
	}
	me.Target = best
}

func targetPlanePenalty(me, e *Soldier) bool {
	mePlane, enemyPlane := me.planeHigh(), e.planeHigh()
	if mePlane > enemyPlane {
		return e.HighTerrain
	}
	return me.Kind <= Cavalry && enemyPlane != PlaneHighGround
}

// applyNewOrder 重現 `sub_1A7B7`：**退卻不可打斷**，換命令時記下起點。
func (s *Soldier) applyNewOrder() bool {
	if s.Cmd == Retreat {
		return false
	}
	if s.Cmd == s.Next {
		return false
	}
	s.Cmd = s.Next
	s.StepX, s.StepY, s.StepZ = s.X, s.Y, s.Z
	// 原版換令時把繞路點游標歸零；舊命令的路徑不能帶進退卻或
	// 下一個攻擊目標，否則兵會永遠沿著上一個目標的路走。
	s.Path = nil
	s.PathAt = 0
	return true
}

// updateSoldier 是一個兵的一幀。
func (b *Battle) updateSoldier(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	// 原版每個兵更新時先清掉這一幀的暫態旗標
	// （`and byte ptr [si], 0BFh`、`and byte ptr [si+2], 1`）。
	s.Swapped, s.Hurt, s.HitGeneral = false, false, false
	b.lockOnNearest(side, k)
	s.applyNewOrder()

	switch s.Cmd {
	case Form:
		b.doFormation(side, k)
	case Attack, Charge:
		b.doAttack(side, k)
	case ScaleWal:
		b.doScaleWall(side, k)
	case Guard:
		b.doGuard(side, k)
	case Retreat:
		b.doRetreat(side, k)
	case Holding:
		// 已就位，原地待命。
	}
	b.moveToward(side, k)
	// `sub_1B240` 尾端的 `xor byte ptr [si+2], 1`。特殊投射物在
	// 本幀攻擊時先取舊值，下一幀人物姿勢再翻面。
	s.PoseStep ^= 1
}

// doFormation 是命令 0：走到陣形指定的座標，**到位就轉成「就位」並補滿疲勞**。
//
// 說明書 4.2 說陣形是唯一能恢復疲勞度的指令，6.1 說疲勞
// 「兵回到陣形的指定位置時最小」——原版就是在到位那一刻
// `mov byte ptr [si+19h], 80h`（docs/re/11 §5.8f）。
// **所以恢復是有移動成本的，不是下令就補。**
func (b *Battle) doFormation(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	x, y := b.formationSpot(side, k)
	s.GoalX, s.GoalY = x, y
	s.GoalZ = b.Field.StandLevel(x, y)
	if s.X == x && s.Y == y {
		s.Stamina = StaminaFull
		s.Cmd, s.Next = Holding, Holding
	}
}

// doAttack 是命令 1／2。大將不攻擊（說明書「大將以外的兵攻擊」），
// 弓兵站著放箭，其餘追過去。
func (b *Battle) doAttack(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	if s.Stamina > StaminaFighting {
		s.Stamina = StaminaFighting
	}
	if s.Target < 0 {
		return
	}
	e := &b.Sides[1-side].Soldiers[s.Target]

	if s.IsGeneral() && s.Cmd == Attack {
		// 命令 1 時大將只移動不攻擊。
		s.GoalX, s.GoalY, s.GoalZ = e.X, e.Y, e.Z
		return
	}
	if s.Kind == Archer {
		s.GoalX, s.GoalY, s.GoalZ = s.X, s.Y, s.Z // 弓兵不動
		b.shoot(side, k, e)
		return
	}
	if specialAttackAvailable(s, e) {
		b.shootSpecial(side, k, e)
		return
	}
	s.GoalX, s.GoalY, s.GoalZ = e.X, e.Y, e.Z
	if abs(s.X-e.X)+abs(s.Y-e.Y) <= 1 && s.Z == e.Z {
		b.attackCollision(side, s, e)
	}
}

// doScaleWall 是命令 3：走到登城點，目標高度取那一格的可站立層。
func (b *Battle) doScaleWall(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	x := b.Field.GateX()
	s.GoalX, s.GoalY = clamp(x), s.Y
	s.GoalZ = b.Field.StandLevel(s.GoalX, s.GoalY)
	if s.X == s.GoalX && s.Z == s.GoalZ {
		s.Cmd, s.Next = Attack, Attack // 上去了就轉成攻擊
	}
}

// doGuard 是命令 4：**敵人曼哈頓距離 ≤ 16 就打，否則回陣形**。
// 門檻出自 `sub_1A99C` 的 `cmp dl, 10h`（docs/re/11 §5.8b）。
const guardRange = 16

func (b *Battle) doGuard(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	if s.Target < 0 {
		b.doFormation(side, k)
		return
	}
	e := &b.Sides[1-side].Soldiers[s.Target]
	if abs(s.X-e.X)+abs(s.Y-e.Y) <= guardRange {
		b.doAttack(side, k)
		return
	}
	// 敵人走遠了 → 回陣形。原版是 `[si+19h] < 16` 時把命令改回 0。
	if s.Stamina < StaminaBackToForm {
		s.Cmd, s.Next = Form, Form
	}
	b.doFormation(side, k)
}

// doRetreat 是命令 5：往自軍側的邊緣走，走出畫面就離場。
//
// 說明書 4.1：「兵士が戦闘して敵にやられると、**画面外に退却**し，
// 待機中の兵が出陣します」。原版 `sub_1AAED` 會把出口寫成
// X=1／62、Y 夾在 0x10..0x2F，Z 固定為 0。
func (b *Battle) doRetreat(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	// 退卻是新的、固定朝向側邊的出口；若兵在受傷前已經有攻擊繞路點，
	// 那個點可能把它留在城門／敵陣附近。`applyNewOrder` 通常會清掉路徑，
	// 但已經處於退卻命令的兵也可能從受擊或隊長連鎖進入這裡，必須在出口
	// 再清一次，否則正常攻城會有兵永遠走不到畫面邊緣。
	s.Path = nil
	s.PathAt = 0
	edge := MinCoord
	if b.Sides[side].Mirror {
		edge = MaxCoord
	}
	retreatY := s.Y
	if retreatY < 0x10 {
		retreatY = 0x10
	}
	if retreatY > 0x2F {
		retreatY = 0x2F
	}
	s.GoalX, s.GoalY, s.GoalZ = edge, retreatY, 0
	if s.X == edge {
		s.Alive = false
		b.squadLeaderGone(side, k)
	}
}

// squadLeaderGone 重現 `sub_1A83F`：某隊第一格的隊長不在場時，該隊
// 剩下的七人全部改排退卻。隊長已經不在，畫面外的預備兵也不能再補成
// 一個沒有隊長的隊伍；清掉該隊待機數才能讓 §5.9 的「補不出兵」成立。
func (b *Battle) squadLeaderGone(side, k int) {
	if k < 0 || k >= SoldiersOnFoot || k%PerSquad != 0 {
		return
	}
	squad := k / PerSquad
	s := &b.Sides[side]
	s.Reserve[squad] = 0
	for j := squad * PerSquad; j < (squad+1)*PerSquad; j++ {
		if j == k || !s.Soldiers[j].Alive {
			continue
		}
		if s.Soldiers[j].Cmd != Retreat {
			s.Soldiers[j].Next = Retreat
		}
	}
}

// moveToward 重現 `sub_1AF69`：**三軸依序試**，先 X、再 Y、最後 Z。
//
// 撞到障礙時會去讀一張預先算好的繞路點清單（`sub_1B00D`，
// docs/re/11 §5.8k），那張清單由 `loc_1BD46` 算（§5.15）。
func (b *Battle) moveToward(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	s.StepX, s.StepY, s.StepZ = s.GoalX, s.GoalY, s.GoalZ

	// 有繞路點就先走目前的中繼點。原版 `sub_1B00D` 只有在抵達
	// 目前的 X/Y/Z 後才消費它；不能每幀直接取下一個點，否則兵只
	// 走一步就會跳過轉角（§5.15）。
	if p, ok := s.Path.Current(); ok {
		pz := b.Field.StandLevel(p.X, p.Y)
		if s.X == p.X && s.Y == p.Y && s.Z == pz {
			s.Path.Advance()
			p, ok = s.Path.Current()
		}
		if ok {
			s.StepX, s.StepY, s.StepZ = p.X, p.Y, b.Field.StandLevel(p.X, p.Y)
		}
	}

	// ⭐ **面向只在真的走成功那一步才更新。**
	// 原版把 `[si+5]` 寫在四個移動常式裡（`sub_1B047`／`1B069`／`1B08B`／
	// `1B0AF`），而那些常式只有在走得動時才被呼叫——所以**被牆擋住的兵
	// 保持原本的面向**。差別看得見：面向決定畫哪一張圖，也決定
	// §5.9 那個城壁分支成不成立。
	moved := false
	if s.X != s.StepX {
		d, face := 1, East
		if s.X > s.StepX {
			d, face = -1, West
		}
		if b.tryMove(side, k, s.X+d, s.Y, s.Z) {
			s.Facing, moved = face, true
		}
	}
	if !moved && s.Y != s.StepY {
		d, face := 1, South
		if s.Y > s.StepY {
			d, face = -1, North
		}
		if b.tryMove(side, k, s.X, s.Y+d, s.Z) {
			s.Facing, moved = face, true
		}
	}
	// ★ 大將與騎馬不做 Z 移動 —— 爬不上城牆（`cmp [si+4], 12h / jbe`）。
	if !moved && s.CanClimb() && s.Z != s.StepZ {
		d := 1
		if s.Z > s.StepZ {
			d = -1
		}
		if b.tryMove(side, k, s.X, s.Y, s.Z+d) {
			moved = true
		}
	}
	if !moved {
		// 三個軸都走不動 → 算一條繞路。原版在 `sub_1AED2` 就是這樣
		// 補上 `0x1800 + 兵編號 × 128` 那塊繞路點清單的（§5.15）。
		b.replan(side, k)
	}
	if moved && s.Stamina > 0 {
		s.Stamina-- // 移動每幀 −1（`sub_1ADC8`）
	}
}

// replanInterval 是重算繞路的最短間隔（幀）。
//
// ⚠ 原版沒有這個節流——它是在「命令生效」那一刻算一次（`sub_1AED2`）。
// 這裡加上是因為本專案的兵每幀都可能被別人擋住，不節流的話 48 × 2 個兵
// 每幀各跑一次波前擴散，無頭模擬會慢到跑不完。**這是 remake 的取捨。**
const replanInterval = 30

// replan 幫一個兵算一條繞開障礙的路。
func (b *Battle) replan(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	if s.Path.Len() > 0 || b.Frame-s.PathAt < replanInterval {
		return
	}
	s.PathAt = b.Frame
	pts := b.Field.FindPath(Point{X: s.X, Y: s.Y},
		Point{X: s.GoalX, Y: s.GoalY}, s.CanClimb(), nil)
	if len(pts) == 0 {
		return
	}
	s.Path = &Waypoints{pts: pts}
}

// tryMove 試著走到一格。走得上去才動。
func (b *Battle) tryMove(side, k, x, y, z int) bool {
	s := &b.Sides[side].Soldiers[k]
	if !inBounds(x, y) {
		return false
	}
	// Z 沒指定成那一格的可站立層時，用那一格的頂。
	if !b.Field.Walkable(x, y, z) {
		z = b.Field.StandLevel(x, y)
		// 水平跨格時，原版 `sub_1B1B1` 會把兵的 Z 同步調整一層；
		// 「大將／騎馬不能爬牆」只限制純 Z 軸移動（`sub_1AF69`
		// 在 X/Y 已到位後的分支），不能拿來阻擋退卻時跨過一格高度 1
		// 的邊界地形。
		if x == s.X && y == s.Y && z > s.Z && !s.CanClimb() {
			return false
		}
		// 一次只能上下一層。
		if abs(z-s.Z) > 1 {
			// 擋住去路的是城壁或門的話，這一撞要算耐久
			// （原版在同一個碰撞路徑上 `dec [di+18h]`，docs/re/11 §5.9）。
			b.hitStructure(side, s.Facing, x, y)
			return false
		}
	}
	// 有人擋著 → 走碰撞處理（`loc_1B533`）。**敵我分兩條路**：
	//
	//	敵人擋路 → 打他（`loc_1B5A1` → `sub_1B618`），自己不動
	//	自己人擋路 → **兩個對調位置**（`loc_1B56D` → `sub_1B732`）
	//
	// 分敵我那一行是 `cmp di, <立即值>`，而那個立即值是
	// **自我修改碼**寫進去的（`byte_1B562`／`byte_1B56A`，由 `sub_19A33`
	// 依雙方的編號範圍填）——這是本作第四處自我修改碼。
	if side2, k2 := b.anyoneAt(x, y, z); k2 >= 0 {
		if side2 != side {
			e := &b.Sides[side2].Soldiers[k2]
			// 大將走 `sub_1B6BC`，其餘敵人走 `sub_1B618`。
			b.attackCollision(side, s, e)
			return false
		}
		return b.swapWith(side, k, side2, k2)
	}
	s.X, s.Y, s.Z = x, y, z
	s.syncTerrain(b.Field, x, y, z)
	return true
}

// swapWith 重現 `seg000:B56D`–`B598` ＋ `sub_1B732`：
// 條件符合就把兩個兵的座標對調，回傳 true（＝「走成功了」）。
//
// 擋路的一方符合下列任一條就不換（原版 `jz/jnz loc_1B59E` ＝ 失敗）：
//
//	[di+04] == 0     擋路的是**大將**
//	[di+1A] == 5     擋路的正在**退卻**
//	[di] & 0x61      bit 0 ＝ **陣亡**（`sub_1B618` 把體力扣到 0 時設，
//	                 同時 `[di+1] = 4`）；bit 5 ＝ **正在移動常式裡面**
//	                 （`sub_1AF69` 進入時 `or [si], 20h`、離開時清掉，
//	                 是重入保護，本專案不需要）；bit 6 ＝ **這一幀已經被換過**
//	[di+02] & 0x10   **剛剛被打中**（`sub_1B618` 設，同時把面向歸零）
//
// 通過之後還分兩種：
//
//	同一層（`[si+1E] == [di+1E]`，那是 Z 平面的位址高位）→ 直接換
//	不同層                                              → **兩邊都要是弓兵或步兵**
//	                                                      （`cmp [si+4], 12h / jbe` 兩次）
//
// 所以**大將與騎馬跨不了層換位**——它們本來就爬不上城牆。
func (b *Battle) swapWith(side, k, side2, k2 int) bool {
	_ = side2 // 走到這裡時一定同側
	me := &b.Sides[side].Soldiers[k]
	other := &b.Sides[side2].Soldiers[k2]

	// bit 6：這一幀已經被換過了就不能再換，否則兩個兵會原地互換不停。
	// bit 4 of +0x02：剛被打中的兵這一幀不能被換走。
	if other.Swapped || other.Hurt || other.IsGeneral() || other.Cmd == Retreat {
		return false
	}
	// 跨層對調只有弓兵與步兵做得到。
	if me.Z != other.Z && (!me.CanClimb() || !other.CanClimb()) {
		return false
	}
	me.X, other.X = other.X, me.X
	me.Y, other.Y = other.Y, me.Y
	me.Z, other.Z = other.Z, me.Z
	me.PlaneHigh, other.PlaneHigh = other.PlaneHigh, me.PlaneHigh
	me.HighTerrain, other.HighTerrain = other.HighTerrain, me.HighTerrain
	me.Climbing, other.Climbing = other.Climbing, me.Climbing
	other.Swapped = true
	return true
}

// anyoneAt 回傳站在那一格的兵（側別與索引），沒有人回 (0, −1)。
func (b *Battle) anyoneAt(x, y, z int) (int, int) {
	for i := range b.Sides {
		for k := range b.Sides[i].Soldiers {
			s := &b.Sides[i].Soldiers[k]
			if s.Alive && s.X == x && s.Y == y && s.Z == z {
				return i, k
			}
		}
	}
	return 0, -1
}
