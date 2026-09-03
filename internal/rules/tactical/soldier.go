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
	return true
}

// updateSoldier 是一個兵的一幀。
func (b *Battle) updateSoldier(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	// ⭐ **這一幀被別人換走的兵不動**（`sub_1ADC8` 的
	// `0001ADED test al, 40h / jnz loc_1AE26`，docs/spec/62）：
	// 直接跳到重畫，不移動也不攻擊。旗標在那裡清掉，所以壽命是
	// 「從被換走那一刻，到自己這一格輪到為止」。
	//
	// 少了這一條，被推開的兵下一幀馬上又往前擠，把剛換到前排的同伴
	// 再換回去——前排那一格一直換人，而換進去的永遠不是下一個輪到
	// 更新的那個，於是圍著打卻一次也打不到（docs/spec/62；docs/spec/61 §5 講它為什麼與開場體力綁在一起）。
	if s.Swapped {
		s.Swapped = false
		s.PoseStep ^= 1
		return
	}
	// ⭐ **挨打之後的硬直**（`sub_1ADC8` 的 `0001ADF1`–`0001ADFE`，
	// docs/spec/63）：計時器大於 0 就遞減後跳過，**歸零那一幀才清掉
	// `Hurt`**。`Hurt` 同時是換位的擋條件之一，所以硬直期間也不會
	// 被同伴換走——站著挨完。
	if s.Stun > 0 {
		if s.Stun--; s.Stun == 0 {
			s.Hurt = false
		}
		s.PoseStep ^= 1
		return
	}
	s.HitGeneral = false
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
	case Duel:
		// 命令 8：分派是 nullsub（`funcs_1A7E1[8]`）——目標座標由
		// 單挑狀態機寫，移動與接觸命中照走下面的共通路徑
		// （docs/spec/80 §4）。
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
	if b.walkToFormation(side, k) {
		s.Stamina = StaminaFull
		s.Cmd, s.Next = Holding, Holding
	}
}

// walkToFormation 把目標設成陣形位置，回報「是不是已經站在那裡」。
//
// ⚠ **它不碰命令。** 「到位轉狀態 7（就位）」是**命令 0 專屬**的
// （docs/re/11 §5.8 的對照表），而守陣（命令 4）也要走回陣形——
// 借用整支 `doFormation` 會讓守陣在下令的下一幀就被降級成就位
// （docs/spec/96）。
func (b *Battle) walkToFormation(side, k int) bool {
	s := &b.Sides[side].Soldiers[k]
	x, y := b.formationSpot(side, k)
	s.GoalX, s.GoalY = x, y
	s.GoalZ = b.standZ(s, x, y)
	return s.X == x && s.Y == y
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
	// ⭐ 目標 Z 取**高平面**那一格的地面層（原版 `loc_1AB39` 的
	// `bh |= 10h` ＋ `al = es:[bx] & 7`，docs/re/11 §5.8i）——
	// 那是**牆頂**，不是腳下的地面。用 standZ 的話，牆腳有地面的那幾格
	// 會算出 Z ＝ 0，於是「目前 Z ＝ 目標 Z」，純 Z 移動一次都不會試。
	if lv, ok := b.Field.GroundLevel(s.GoalX, s.GoalY, PlaneHigh); ok {
		s.GoalZ = lv
	} else {
		s.GoalZ = b.standZ(s, s.GoalX, s.GoalY)
	}
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
		// ⚠ 還沒鎖到敵人 → 待在陣形位置，**但命令維持守陣**。
		// 開場的兵本來就站在陣形位置上，借用 doFormation 的話
		// 第一幀就會被判成「到位」而降級成就位——等到鎖上敵人時
		// 命令已經不是守陣了，永遠不會反應（docs/spec/96）。
		b.walkToFormation(side, k)
		return
	}
	e := &b.Sides[1-side].Soldiers[s.Target]
	if abs(s.X-e.X)+abs(s.Y-e.Y) <= guardRange {
		b.doAttack(side, k)
		return
	}
	// 敵人走遠了 → 回陣形。原版**只有疲勞** `[si+19h] < 16 時才把命令改回 0
	// （`sub_1ABB2`）；距離只換行為，不換命令。
	if s.Stamina < StaminaBackToForm {
		s.Cmd, s.Next = Form, Form
		b.doFormation(side, k)
		return
	}
	b.walkToFormation(side, k)
}

// doRetreat 是命令 5：往自軍側的邊緣走，走出畫面就離場。
//
// 說明書 4.1：「兵士が戦闘して敵にやられると、**画面外に退却**し，
// 待機中の兵が出陣します」。原版 `sub_1AAED` 會把出口寫成
// X=1／62、Y 夾在 0x10..0x2F，Z 固定為 0。
func (b *Battle) doRetreat(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	edge := MinCoord
	if b.Sides[side].Mirror {
		edge = MaxCoord
	}
	// 退卻是新的、固定朝向側邊的出口；若兵在受傷前已經有攻擊繞路點，
	// 那個點可能把它留在城門／敵陣附近。`applyNewOrder` 通常會清掉路徑，
	// 但已經處於退卻命令的兵也可能從受擊或隊長連鎖進入這裡，所以在出口
	// 再判一次。
	//
	// ⚠ **只清「不是通往出口的那條路」**（docs/spec/94）。這裡本來是
	// 每一幀無條件清掉，而一幀之內的順序是 doRetreat → moveToward（用路）
	// → 走不動 → replan（算路）——清除在最前面，於是算好的繞路點
	// **下一幀開頭就被丟掉**，兵永遠只朝終點直線推。配上「退卻中的同伴
	// 不能對調」（sub_1B732 的閘），一整排就鎖死，攻城戰永遠不會結束。
	//
	// 判準不能用「目標相等」：退卻目標的 Y 是 clamp(s.Y)，會跟著兵走，
	// 繞路一走 Y 就又被清掉了。
	if p, ok := s.Path.Last(); !ok || p.X != edge {
		s.Path = nil
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
		// ⭐ 走到畫面外的兵**算生還**（原版 `sub_1B4B8` 的 `ah = 0`
		// 那條路加一，docs/spec/65）。戰死的走同一支但不加。
		b.Sides[side].Escaped[k/PerSquad]++
		b.squadLeaderGone(side, k)
	}
}

// squadLeaderGone 重現 `sub_1A83F`：某隊第一格的隊長不在場時，該隊
// 剩下的七人全部改排退卻。
//
// ⚠ **不清那一隊的待機數。** 原版沒有這一步——碰結算緩衝區的 11 支函式
// 逐一看完，待機數只有 `dec`（開場取用、補兵），沒有任何地方寫 0
// （docs/re/83 §1、docs/spec/128）。那些兵在原版會補進場、下一幀又被
// 改成退卻、走出畫面，最後**被算成生還**（`sub_1B4B8` 的 `ah = 0`）。
// 先前清成 0 的版本讓它們從帳上直接消失，方向是少算兵力。
func (b *Battle) squadLeaderGone(side, k int) {
	if k < 0 || k >= SoldiersOnFoot || k%PerSquad != 0 {
		return
	}
	squad := k / PerSquad
	s := &b.Sides[side]
	for j := squad * PerSquad; j < (squad+1)*PerSquad; j++ {
		if j == k || !s.Soldiers[j].Alive {
			continue
		}
		if s.Soldiers[j].Cmd != Retreat {
			s.Soldiers[j].Next = Retreat
		}
	}
}

// applySquadLeaderGone 是 `sub_1A754`／`sub_1A785` 的逐隊檢查：**每一幀**
// 看每一隊隊長的在場位元，不在就對那一隊施加 `sub_1A83F`。
//
// ⭐ 原版是**每幀重新施加**，不是死亡當下的一次性事件——所以隊長倒下之後
// 補進場的兵，下一幀也會被改成退卻（docs/re/83 §4）。
func (b *Battle) applySquadLeaderGone() {
	for side := range b.Sides {
		for squad := 0; squad < Squads; squad++ {
			k := squad * PerSquad
			if b.Sides[side].Soldiers[k].Alive {
				continue
			}
			b.squadLeaderGone(side, k)
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
		pz := b.standZ(s, p.X, p.Y)
		if s.X == p.X && s.Y == p.Y && s.Z == pz {
			s.Path.Advance()
			p, ok = s.Path.Current()
		}
		if ok {
			s.StepX, s.StepY, s.StepZ = p.X, p.Y, b.standZ(s, p.X, p.Y)
		}
	}

	// ⭐ **面向只在真的走成功那一步才更新。**
	// 原版把 `[si+5]` 寫在四個移動常式裡（`sub_1B047`／`1B069`／`1B08B`／
	// `1B0AF`），而那些常式只有在走得動時才被呼叫——所以**被牆擋住的兵
	// 保持原本的面向**。差別看得見：面向決定畫哪一張圖，也決定
	// §5.9 那個城壁分支成不成立。
	moved, walled := false, false
	if s.X != s.StepX {
		d, face := 1, East
		if s.X > s.StepX {
			d, face = -1, West
		}
		ok, blocked := b.tryMove(side, k, s.X+d, s.Y, s.Z)
		if ok {
			s.Facing, moved = face, true
		}
		walled = walled || blocked
	}
	if !moved && s.Y != s.StepY {
		d, face := 1, South
		if s.Y > s.StepY {
			d, face = -1, North
		}
		ok, blocked := b.tryMove(side, k, s.X, s.Y+d, s.Z)
		if ok {
			s.Facing, moved = face, true
		}
		walled = walled || blocked
	}
	// ★ 純 Z 移動：**只在門那一格**，而且大將與騎馬不做
	// （`sub_1B0D3` 開頭的 `cmp al, 0F0h`，加上 `sub_1AF69` 的
	// `cmp [si+4], 12h / jbe`，docs/re/63 §4）。
	//
	// ⚠ **觸發條件只有「X 與 Y 這一幀都沒走成」**（docs/spec/97）。
	// `sub_1AF69` 的 `0001AF78`（Y）與 `0001AF82`（Z）是同一條 `jb` 鏈的
	// 下一站——走失敗會跳過來，已經等於目標也會直接落下來。
	// 這裡本來還要求 `s.X == s.StepX && s.Y == s.StepY`，於是登城的兵
	// 站在門格上被前面的城壁擋住時**不會試著往上爬**，改成一直撞牆
	// 把城壁磨穿——而原版是爬上去的。
	if !moved && s.Z != s.StepZ {
		if b.tryClimb(side, k) {
			moved = true
		}
	}
	// 三個軸都走不動 → 算一條繞路。原版在 `sub_1AED2` 就是這樣
	// 補上 `0x1800 + 兵編號 × 128` 那塊繞路點清單的（§5.15）。
	//
	// ⚠ **被地形擋住時也要重算，即使這一幀靠跟同伴對調而「動了」。**
	// 對調（`sub_1B732`）會回報成功，於是前排卡在城牆上的兵靠著跟後排
	// 換位就一直算「有移動」，永遠不重算路——整團會在城牆前churn 到
	// 攻城計時器把大將耗光。
	if !moved || walled {
		b.requestPath(side, k)
	}
	if moved && s.Stamina > 0 {
		s.Stamina-- // 移動每幀 −1（`sub_1ADC8`）
	}
}

// occupancyCost 是尋路的額外成本：**有兵站著的格子加 8**。
//
// 原版把它存在另一張表（`word_1D2FE`），移動落定時舊格寫 0、新格寫 8，
// 尋路擴散時 `mov al, es:[bx+2000h]` 讀出來加進波數
// （docs/re/63 §1、docs/re/11 §5.15）。**這不是「不能走」，是「繞得過就繞」。**
//
// 回傳的是一個查表函式：一次掃完 96 個兵建表，不要每查一格掃一次——
// 尋路一次會查上千格，逐格掃會讓無頭模擬慢一個數量級。
func (b *Battle) occupancyCost() Penalty {
	var occupied [Height][Width]bool
	for i := range b.Sides {
		for k := range b.Sides[i].Soldiers {
			s := &b.Sides[i].Soldiers[k]
			if s.Alive && s.X >= 0 && s.X < Width && s.Y >= 0 && s.Y < Height {
				occupied[s.Y][s.X] = true
			}
		}
	}
	return func(x, y int) int {
		if x < 0 || x >= Width || y < 0 || y >= Height || !occupied[y][x] {
			return 0
		}
		return occupiedPenalty
	}
}

// occupiedPenalty 是原版寫進佔用表的值（`sub_1B240` 的 `mov byte …, 8`）。
const occupiedPenalty = 8

// standZ 是一個兵走到 (x, y) 之後會站在哪一層。
//
// 先看**自己所在的平面**；那個平面在那一格沒有地面時看另一個平面——
// 這正是「走到門那一格然後爬上去」的訊號：`moveToward` 的 X／Y 都到位
// 之後 Z 還差，就會去叫 `tryClimb`（docs/re/63 §4）。
func (b *Battle) standZ(s *Soldier, x, y int) int {
	if lv, ok := b.Field.GroundLevel(x, y, s.Plane()); ok {
		return lv
	}
	other := PlaneHigh
	if s.OnWall {
		other = PlaneLow
	}
	if lv, ok := b.Field.GroundLevel(x, y, other); ok {
		return lv
	}
	return b.Field.StandLevel(x, y)
}

// computePath 幫一個兵算一條繞開障礙的路（原版 `sub_1AED2` 的出隊段）。
//
// ⚠ **手上有路不是不重算的理由。** 進到佇列的前提是「這一幀走不動或
// 撞到地形」，所以有路又走不動代表那條路現在不通——原版也是這樣
// （`docs/re/80` §3 的四個入隊點）。這裡本來會因為 `Path.Len() > 0`
// 直接返回，於是被同伴擋住的兵抱著一條穿過同伴的直線路永遠不重算
// （`docs/spec/94` §2.1）。
//
// ⭐ **節流不在這裡**：原版的預算是全域的「每幀兩筆」，
// 由 `pathQueue` 管（`docs/spec/120`）。
func (b *Battle) computePath(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	from, to := Point{X: s.X, Y: s.Y}, Point{X: s.GoalX, Y: s.GoalY}
	occupied := b.occupancyCost()
	pts := b.Field.FindPath(from, to, s.CanClimb(), occupied)
	if len(pts) == 0 {
		// 地形走不通 → 改成「可以拆的就穿過去」。兵會走到那一格前面
		// 撞上去，`tryMove` 把它算成一次耐久損傷，撞穿了地形就通了。
		// 沒有這一步的話，攻城時攻方會整團卡在打不壞的城體前面。
		pts = b.Field.FindPathForcing(from, to, s.CanClimb(),
			func(x, y int) int { return b.breachCost(x, y) + occupied(x, y) },
			b.breakableAt)
	}
	if len(pts) == 0 {
		return
	}
	s.Path = &Waypoints{pts: pts}
}

// tryMove 試著走到一格。走得上去才動。
//
// 第二個回傳值是「**被地形擋住**」——擋路的是牆或高低差，不是別的兵。
// 呼叫端要靠它決定要不要重算繞路：撞到兵可以靠對調解決，撞到地形不行。
func (b *Battle) tryMove(side, k, x, y, z int) (moved, walled bool) {
	s := &b.Sides[side].Soldiers[k]
	if !inBounds(x, y) {
		return false, true
	}
	// 水平跨格照 `sub_1B1B1`（docs/re/63 §5）：讀**自己所在平面**的
	// 地面層，沒有地面就是撞到牆；有地面就把 Z 同步一層。
	// 「大將／騎馬不能爬牆」限制的是純 Z 移動（tryClimb），不在這裡。
	lv, ok := b.Field.GroundLevel(x, y, s.Plane())
	if !ok {
		// 擋住去路的是城壁或門的話，這一撞要算耐久
		// （原版在同一個碰撞路徑上 `dec [di+18h]`，docs/re/11 §5.9）。
		b.hitStructure(side, s.Facing, x, y)
		return false, true
	}
	// `sub_1B1B1` 的第一段：目標格**自己那一層之上**有地形就擋
	// （`ah = es:[bx+1000h] / and ah, ah / jnz 失敗`）。城壁就是這樣擋人的
	// ——它的地面層表是拿「打壞後的圖塊」算的，本來就有地面（docs/re/63 §2）。
	if b.Field.Blocked(x, y, s.Z+1) {
		b.hitStructure(side, s.Facing, x, y)
		return false, true
	}
	if s.Plane() == PlaneHigh {
		// 高平面要高度完全相等，而且未破的門橫向穿不過去。
		if lv != s.Z || b.Field.GateBlocksHighPlane(x, y) {
			b.hitStructure(side, s.Facing, x, y)
			return false, true
		}
		z = lv
	} else {
		switch {
		case lv > s.Z:
			z = s.Z + 1
		case lv < s.Z:
			z = s.Z - 1
		default:
			z = s.Z
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
			return false, false
		}
		return b.swapWith(side, k, side2, k2), false
	}
	s.X, s.Y, s.Z = x, y, z
	s.syncTerrain(b.Field, x, y, z)
	return true, false
}

// tryClimb 重現 `sub_1B0D3`／`sub_1B116` ＋ 它們呼叫的
// `sub_1B186`／`sub_1B15D`：在**打破的**門那一格上下城牆。
//
// 四個條件缺一不可（docs/spec/36 §1.4）：
//
//	兵種 > 0x12（大將與騎馬不做 Z 移動）
//	腳下那一格的圖塊 ≥ 0xF0，也就是門          （呼叫端擋的）
//	同一格的圖塊 ≥ 0xF8，也就是**已經打破**    （被呼叫的那一支擋的）
//	要去的那個平面在這一格有地面、而且沒有別人站著
//
// ⚠ **未破的門爬不上去。** 兩層檢查讀的是同一格的同一個 byte
// （原版 `di = bx & 0FFFh` 把平面位元遮掉了），所以 `≥ 0xF0` 那一關
// 其實被 `≥ 0xF8` 蓋過去。尋路那邊**沒有**這一關，路徑會規劃穿過未破的
// 門，走到才被擋——這個不對稱是原版行為，照抄。
func (b *Battle) tryClimb(side, k int) bool {
	s := &b.Sides[side].Soldiers[k]
	if !s.CanClimb() || !b.Field.IsGateCell(s.X, s.Y) {
		return false
	}
	if b.Field.GateBlocksHighPlane(s.X, s.Y) {
		// ⭐ **爬不上去的那一下要打門**（docs/spec/98）。原版
		// `sub_1B186` 回報「上一層有實體」時把實體編號留在 al，
		// `sub_1B0D3` 的 `and al, al / jnz` 就據此走 `loc_1B533`
		// 碰撞處理——**未破的門就是這樣被打開的**（耐久只有 80）。
		// 少了這一下，登城的兵站在門格上永遠卡住，攻方只能改去磨
		// 城壁（耐久上千），而那會讓門強度條一路亮著——原版不會。
		if b.breakableAt(s.X, s.Y) {
			b.hitStructure(side, s.Facing, s.X, s.Y)
		}
		return false
	}
	other := PlaneHigh
	if s.OnWall {
		other = PlaneLow
	}
	lv, ok := b.Field.GroundLevel(s.X, s.Y, other)
	if !ok {
		return false
	}
	// 目標層有別人就上不去（原版 `sub_1B186` 檢查上一層）。
	// **自己不算**：低平面與高平面同高時，站著的就是自己。
	if side2, k2 := b.anyoneAt(s.X, s.Y, lv); k2 >= 0 && !(side2 == side && k2 == k) {
		return false
	}
	s.OnWall = other == PlaneHigh
	s.Z = lv
	s.syncTerrain(b.Field, s.X, s.Y, s.Z)
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
