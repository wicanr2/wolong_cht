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

// InfantryArrowDivisor 是步兵挨箭的減傷倍數。
//
// 原版 `sub_1B97E`：`cmp byte ptr [bx+4], 36h / jz` → **威力右移兩位**。
// 0x36 ＝ 54 ＝ 步兵。說明書「攻城戦では弓兵、**歩兵**が必要です」的
// 數值依據就在這裡：步兵是唯一扛得住城牆上箭雨的兵種。
const InfantryArrowDivisor = 4

// hit 對一個兵造成傷害。
//
// 三條規則全部照原版：
//
//   - **步兵挨箭只吃四分之一**
//   - **大將不會陣亡**：體力扣到 1 就停住（`cmp [bx+4], 0 / jz` 跳過死亡分支）
//   - 受傷的兵**自己把命令改成退卻**，但體力還滿的不會
//     （`cmp byte ptr [bx+3], 64h / jnb` 跳過）
func (b *Battle) hit(from int, e *Soldier, power int) {
	if !e.Alive {
		return
	}
	if e.Kind == Infantry && power > 0 {
		power >>= 2
	}
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
			b.hit(p.side, hitSoldier, p.power)
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
