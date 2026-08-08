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
		// 爬不上城牆的兵不該把牆上的敵人當成主要目標——原版是加一個
		// 64 格的距離懲罰讓它排到最後（`sub_1A85B`，docs/re/11 §5.8c）。
		if e.Z > me.Z && !me.CanClimb() {
			d += Width
		}
		if d < bestD {
			best, bestD = i, d
		}
	}
	me.Target = best
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
	return true
}

// updateSoldier 是一個兵的一幀。
func (b *Battle) updateSoldier(side, k int) {
	s := &b.Sides[side].Soldiers[k]
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
	s.GoalX, s.GoalY, s.GoalZ = e.X, e.Y, e.Z
	if abs(s.X-e.X)+abs(s.Y-e.Y) <= 1 && s.Z == e.Z {
		b.hit(side, e, meleePower)
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
// 說明書 4.1：「兵士が戦闘して敵にやられると、**画面外に退却**し、
// 待機中の兵が出陣します」。
func (b *Battle) doRetreat(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	edge := MinCoord
	if b.Sides[side].Mirror {
		edge = MaxCoord
	}
	s.GoalX, s.GoalY, s.GoalZ = edge, s.Y, b.Field.StandLevel(edge, s.Y)
	if s.X == edge {
		s.Alive = false
	}
}

// moveToward 重現 `sub_1AF69`：**三軸依序試**，先 X、再 Y、最後 Z。
//
// ⚠ 原版撞到障礙時會去讀一張預先算好的繞路點清單（`sub_1B00D`，
// docs/re/11 §5.8k），而**算那張清單的程式還沒解出來**。
// 這裡撞到就換下一軸，解開之後要換掉——**標成 remake 差異**。
func (b *Battle) moveToward(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	s.StepX, s.StepY, s.StepZ = s.GoalX, s.GoalY, s.GoalZ

	// 有繞路點就先走繞路點——原版 `sub_1B00D` 每次取一個當中繼點，
	// 取完（`[si+0x17]` 減到 −1）才回頭直接朝目標走（§5.15）。
	if p, ok := s.Path.Next(); ok {
		if s.Path.Len() > 0 || (p.X == s.GoalX && p.Y == s.GoalY) {
			s.StepX, s.StepY = p.X, p.Y
			s.StepZ = b.Field.StandLevel(p.X, p.Y)
		}
	}

	moved := false
	if s.X != s.StepX {
		d := 1
		s.Facing = East
		if s.X > s.StepX {
			d, s.Facing = -1, West
		}
		if b.tryMove(side, k, s.X+d, s.Y, s.Z) {
			moved = true
		}
	}
	if !moved && s.Y != s.StepY {
		d := 1
		s.Facing = South
		if s.Y > s.StepY {
			d, s.Facing = -1, North
		}
		if b.tryMove(side, k, s.X, s.Y+d, s.Z) {
			moved = true
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
		// 爬不上去的兵不能靠這個上牆。
		if z > s.Z && !s.CanClimb() {
			return false
		}
		// 一次只能上下一層。
		if abs(z-s.Z) > 1 {
			// 擋住去路的是城壁或門的話，這一撞要算耐久
			// （原版在同一個碰撞路徑上 `dec [di+18h]`，docs/re/11 §5.9）。
			b.hitStructure(x, y)
			return false
		}
	}
	// 退卻中的兵可以穿過自己人——它正在從陣線裡撤出來。
	// 原版怎麼處理還沒解，**這是 remake 的決定**：不讓它穿過去的話，
	// 後排的兵會把前排堵死，整條輸送帶就卡住了。
	if s.Cmd != Retreat && b.occupied(x, y, z) {
		return false
	}
	s.X, s.Y, s.Z = x, y, z
	return true
}

// occupied 回報那一格上有沒有人。原版把佔用狀況寫在立體格上
// （每格存「兵編號 + 1」，docs/re/11 §5.3）。
func (b *Battle) occupied(x, y, z int) bool {
	for i := range b.Sides {
		for k := range b.Sides[i].Soldiers {
			s := &b.Sides[i].Soldiers[k]
			if s.Alive && s.X == x && s.Y == y && s.Z == z {
				return true
			}
		}
	}
	return false
}
