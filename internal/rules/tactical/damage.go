package tactical

// 傷害與飛道具。公式出自 `sub_1B97E`（命中）與 `sub_1BAB7`（威力隨高度變化），
// 見 docs/re/11 §5.3–§5.5。

const (
	// meleePower 是近戰一次的威力。
	//
	// ⚠ **實際數值還沒反組譯出來**（傷害是從飛道具記錄的 `+0x04` 讀的，
	// 而那個欄位是誰填的還沒解）。這是 remake 的暫定值。
	meleePower = 6
	// arrowPower 是箭的初始威力。同樣是暫定值。
	arrowPower = 12

	// arrowRange 是弓兵的射程（曼哈頓距離）。暫定值。
	arrowRange = 20
	// arrowCooldown 是兩箭之間的間隔（幀）。暫定值。
	arrowCooldown = 24
)

// hitByArrow 是**飛道具**的傷害（原版 `sub_1B97E`）。
//
// 與近戰的差別只有一條：**目標是步兵就右移兩位**
// （`cmp byte ptr [bx+4], 36h / jz` → `shr al,1` × 2）。
// 大將也吃這個減傷，但要多一個條件（`test byte ptr [bx+2], 1`）。
func (b *Battle) hitByArrow(from int, e *Soldier, power int) {
	if e.Kind == Infantry && power > 0 {
		power /= InfantryArrowDivisor
		if power < 1 {
			power = 1
		}
	}
	b.hit(from, e, power)
}

// InfantryArrowDivisor 是步兵挨箭的減傷倍數。
//
// 原版 `sub_1B97E`：`cmp byte ptr [bx+4], 36h / jz` → **威力右移兩位**。
// 0x36 ＝ 54 ＝ 步兵。說明書「攻城戦では弓兵、**歩兵**が必要です」的
// 數值依據就在這裡：步兵是唯一扛得住城牆上箭雨的兵種。
const InfantryArrowDivisor = 4

// hit 是**近戰**傷害（原版 `sub_1B618`）。
//
//   - **大將不會陣亡**：體力扣到 1 就停住（`cmp [bx+4], 0 / jz` 跳過死亡分支）
//   - 受傷的兵**自己把命令改成退卻**，但體力還滿的不會
//     （`cmp byte ptr [bx+3], 64h / jnb` 跳過）
//
// ⚠ **這裡沒有「步兵吃四分之一」**——那條在 `sub_1B97E`，
// 也就是**飛道具**的命中常式，近戰的 `sub_1B618` 沒有。
// 本專案一度把它套在兩邊，結果近戰威力 6 右移兩位變成 1，
// 步兵幾乎打不動，戰鬥卡住（見 hitByArrow）。
func (b *Battle) hit(from int, e *Soldier, power int) {
	if !e.Alive {
		return
	}
	// 受擊硬直：原版 `sub_1B618` 一起做三件事——設 `+0x02` bit 4、
	// **把面向歸零**、設「這一幀被動過」。面向歸零之後圖號公式會
	// 一律畫正面（§5.13 的 bit 4），所以三件事其實是同一件：
	// 被打中的兵轉過來、畫受擊的圖、那一幀不能被換位置。
	e.Hurt = true

	wasFull := e.HP >= MaxHP
	e.HP -= power
	if e.HP <= 0 {
		e.HP = 1
		if !e.IsGeneral() {
			e.HP = 0
			e.Alive = false
			return
		}
	}
	// 受傷就退卻——但體力原本是滿的那一下不算。
	if !wasFull && e.Cmd != Retreat {
		e.Next = Retreat
	}
}

// shoot 讓弓兵放一箭。
func (b *Battle) shoot(side, k int, e *Soldier) {
	s := &b.Sides[side].Soldiers[k]
	if abs(s.X-e.X)+abs(s.Y-e.Y) > arrowRange {
		return
	}
	// 用幀數錯開，免得 8 個弓兵同一幀齊射。
	if (b.Frame+k)%arrowCooldown != 0 {
		return
	}
	b.projectiles = append(b.projectiles, projectile{
		side: side, x: s.X, y: s.Y, z: s.Z,
		tx: e.X, ty: e.Y, tz: e.Z, power: arrowPower,
	})
}

// stepProjectiles 推進所有飛道具一格。
//
// ⭐ **威力隨高度改變**（`sub_1BAB7`）：
//
//	往上飛 → 威力 −25%
//	往下落 → 威力 +25%
//
// 這解釋了說明書為什麼強調登上城牆攻擊——站在高處往下射，一路都在加成。
func (b *Battle) stepProjectiles() {
	live := b.projectiles[:0]
	for _, p := range b.projectiles {
		switch {
		case p.x != p.tx:
			p.x += sign(p.tx - p.x)
		case p.y != p.ty:
			p.y += sign(p.ty - p.y)
		}
		// 高度線性逼近目標，並依方向調整威力。
		if p.z != p.tz {
			d := sign(p.tz - p.z)
			p.z += d
			q := p.power >> 2
			if d > 0 {
				p.power -= q // 往上
			} else {
				p.power += q + 1 // 往下
			}
		}
		if !inBounds(p.x, p.y) || p.power <= 0 {
			continue
		}
		if hitSoldier := b.soldierAt(1-p.side, p.x, p.y, p.z); hitSoldier != nil {
			b.hitByArrow(p.side, hitSoldier, p.power)
			continue // 箭消失
		}
		if p.x == p.tx && p.y == p.ty {
			continue // 落地
		}
		live = append(live, p)
	}
	b.projectiles = live
}

// soldierAt 回傳某一側站在那一格的兵。**誤傷自己人不判定**——
// 原版是用兵編號的區間分陣營（0x01–0x2F 對 0x30–0x5F）。
func (b *Battle) soldierAt(side, x, y, z int) *Soldier {
	for k := range b.Sides[side].Soldiers {
		s := &b.Sides[side].Soldiers[k]
		if s.Alive && s.X == x && s.Y == y && s.Z == z {
			return s
		}
	}
	return nil
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
	g.Hurt, g.Swapped = true, true
	dmg := attacker.Power >> generalDamageShift
	if dmg < 1 {
		dmg = 1
	}
	g.HP -= dmg
	if g.HP < GeneralMinHP {
		g.HP = GeneralMinHP
	}
	attacker.HitGeneral = true // 圖號 bit 3
	return true
}
