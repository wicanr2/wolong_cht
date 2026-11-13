package tactical

// collisionGrid 保存原版碰撞層低七位與雙平面成本（spec/159、spec/134）。
// 地形仍由 Field 負責；此表不以 Alive 或即時座標重建查詢答案。
type collisionGrid struct {
	cells      [8 * 0x1000]byte
	cost       [numPlanes * 0x1000]byte
	oldPlane   [2 * SoldiersOnFoot]int
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
	// 首次移動批次前只有結構物成本；前置碰撞登記不先計入單兵8。
	b.unitCollision.cost = [numPlanes * 0x1000]byte{}
	b.unitCollision.seedStructureCosts(b.Structures)
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
	s := &b.Sides[side].Soldiers[k]
	// sub_1B240的1B2D3–1B300：換位、硬直或停留也依目前格重設bit1。
	if s.Alive && b.Field != nil {
		s.HighTerrain = b.Field.HighTerrain(s.X, s.Y)
	}
	if b.unitCollision == nil {
		return
	}
	id := b.collisionSlot(side, k)
	if s.Alive {
		b.unitCollision.drawSlot(id, s.X, s.Y, s.Z, int(s.planeHigh())>>4)
	} else if !b.unitDying(side, k) {
		b.clearUnitCollision(side, k)
	}
}

func collisionCell(x, y, z int) (int, bool) {
	if !inBounds(x, y) || z < 0 || z >= 7 {
		return 0, false
	}
	return z*0x1000 + y*Width + x, true
}

// clearSlot 清物理碰撞層；成本由重畫或移除按各自平面清除。重疊時也直接清零。
func (g *collisionGrid) clearSlot(slot int) {
	if slot < 0 || slot >= len(g.old) || !g.registered[slot] {
		return
	}
	p := g.old[slot]
	g.cells[p], g.cells[p+0x1000] = 0, 0
	g.registered[slot] = false
}

// drawSlot 對應sub_1B240：物理層清舊寫新，成本清舊平面、寫目前平面。
func (g *collisionGrid) drawSlot(slot, x, y, z, plane int) {
	if slot < 0 || slot >= len(g.old) {
		return
	}
	p, ok := collisionCell(x, y, z)
	if !ok || plane < 0 || plane >= numPlanes {
		return
	}
	if !g.registered[slot] {
		g.old[slot], g.registered[slot] = p, true
		g.oldPlane[slot] = plane
	}
	g.cost[g.old[slot]&0xFFF+g.oldPlane[slot]*0x1000] = 0
	g.clearSlot(slot)
	g.cells[p] |= byte(slot + 1)
	g.cells[p+0x1000] |= byte(slot + 1)
	g.old[slot], g.registered[slot] = p, true
	g.oldPlane[slot] = plane
	g.cost[p&0xFFF+plane*0x1000] = occupiedPenalty
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

// seedStructureCosts 對應19D6C／19DE8；只在建場寫低平面，不於每次尋路重建。
func (g *collisionGrid) seedStructureCosts(structures []Structure) {
	for _, s := range structures {
		if s.Broken {
			continue
		}
		value := byte(100)
		if s.Kind == KindGate {
			value = 50
		} else if s.Kind != KindWall {
			continue
		}
		for n := 0; n < s.Run; n++ {
			if index, ok := costCell(s.X, s.Y+n, PlaneLow); ok {
				g.cost[index] = value
			}
		}
	}
}

// removeSlot 對應1B3F0–1B3F6，以目前平面清舊格成本；重畫另用oldPlane。
func (g *collisionGrid) removeSlot(slot, plane int) {
	if slot < 0 || slot >= len(g.old) || !g.registered[slot] {
		return
	}
	if plane >= 0 && plane < numPlanes {
		g.cost[g.old[slot]&0xFFF+plane*0x1000] = 0
	}
	g.clearSlot(slot)
}

func (b *Battle) clearUnitCollision(side, k int) {
	if b.unitCollision != nil {
		b.unitCollision.removeSlot(b.collisionSlot(side, k), int(b.Sides[side].Soldiers[k].planeHigh())>>4)
	}
}

// costCell 是完整64×64成本表邊界；不套用單兵1..62的活動範圍。
func costCell(x, y, plane int) (int, bool) {
	if x < 0 || x >= Width || y < 0 || y >= Height || plane < 0 || plane >= numPlanes {
		return 0, false
	}
	return plane*0x1000 + y*Width + x, true
}
