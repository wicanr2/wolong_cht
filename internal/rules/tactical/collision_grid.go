package tactical

// collisionGrid 保存原版碰撞層的單位低七位（spec/159）。
// 地形仍由 Field 負責；此表不以 Alive 或即時座標重建查詢答案。
type collisionGrid struct {
	cells      [8 * 0x1000]byte
	old        [2 * SoldiersOnFoot]int
	registered [2 * SoldiersOnFoot]bool
}

func (b *Battle) collisionSlot(side, k int) int {
	if side == b.PlayerSide {
		return k
	}
	return SoldiersOnFoot + k
}

func (b *Battle) initUnitCollision() {
	b.unitCollision = &collisionGrid{}
	for _, side := range []int{b.PlayerSide, 1 - b.PlayerSide} {
		for k := range b.Sides[side].Soldiers {
			b.drawUnitCollision(side, k)
		}
	}
}

func (b *Battle) unitDying(side, k int) bool {
	for _, d := range b.deaths {
		if d.Side == side && d.Slot == k {
			return true
		}
	}
	return false
}

func (b *Battle) drawUnitCollision(side, k int) {
	if b.unitCollision == nil {
		return
	}
	s := &b.Sides[side].Soldiers[k]
	id := b.collisionSlot(side, k)
	if s.Alive {
		b.unitCollision.drawSlot(id, s.X, s.Y, s.Z)
	} else if !b.unitDying(side, k) {
		b.unitCollision.clearSlot(id)
	}
}

func collisionCell(x, y, z int) (int, bool) {
	if !inBounds(x, y) || z < 0 || z >= 7 {
		return 0, false
	}
	return z*0x1000 + y*Width + x, true
}

// clearSlot 對應 sub_1B3B2；重疊佔位時也清零，不能自行保護其他編號。
func (g *collisionGrid) clearSlot(slot int) {
	if slot < 0 || slot >= len(g.old) || !g.registered[slot] {
		return
	}
	p := g.old[slot]
	g.cells[p], g.cells[p+0x1000] = 0, 0
	g.registered[slot] = false
}

// drawSlot 對應 sub_1B240：清原格，OR 寫目前格與上一層。
func (g *collisionGrid) drawSlot(slot, x, y, z int) {
	if slot < 0 || slot >= len(g.old) {
		return
	}
	p, ok := collisionCell(x, y, z)
	if !ok {
		return
	}
	if !g.registered[slot] {
		g.old[slot], g.registered[slot] = p, true
	}
	g.clearSlot(slot)
	g.cells[p] |= byte(slot + 1)
	g.cells[p+0x1000] |= byte(slot + 1)
	g.old[slot], g.registered[slot] = p, true
}

// at 對應 sub_1B1B1：先腳下層，再上一層。回傳原版槽序，-1 表示空。
func (g *collisionGrid) at(x, y, z int) int {
	p, ok := collisionCell(x, y, z)
	if !ok {
		return -1
	}
	v := g.cells[p]
	if v == 0 {
		v = g.cells[p+0x1000]
	}
	if v == 0 {
		return -1
	}
	return int(v) - 1
}

// swapSlots 對應 sub_1B732 的格內容與原格指標交換。
func (g *collisionGrid) swapSlots(a, b, currentA, currentB int) {
	if a < 0 || b < 0 || a >= len(g.old) || b >= len(g.old) {
		return
	}
	if currentA < 0 || currentB < 0 || currentA+0x1000 >= len(g.cells) || currentB+0x1000 >= len(g.cells) {
		return
	}
	// 0001B761 交換後 AH 未寫回來源上層；保留原版不對稱。
	// 0001B769 從交換後的 SI+0x0C 取到 B，0001B76C 將 B 原值寫回 B。
	// 低層最後未改；不得從 xchg 名稱推定為對稱交換。
	g.cells[currentB+0x1000] = g.cells[currentA+0x1000]
	g.old[a], g.old[b] = g.old[b], g.old[a]
	g.registered[a], g.registered[b] = g.registered[b], g.registered[a]
}
