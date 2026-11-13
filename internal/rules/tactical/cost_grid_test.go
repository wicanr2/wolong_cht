package tactical

import "testing"

func TestCostBeforeFirstMovementOnlyContainsStructures(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), nil, 0)
	s := &b.Sides[0].Soldiers[0]
	*s = Soldier{Alive: true, X: 5, Y: 20}
	b.initUnitCollision()
	if b.occupancyCost()(5, 20, PlaneLow) != 0 {
		t.Fatal("首次移動批次前不可先計入單兵成本")
	}
	b.drawUnitCollision(0, 0)
	if b.occupancyCost()(5, 20, PlaneLow) != 8 {
		t.Fatal("第一次重畫後應寫入單兵成本")
	}
}

func TestCostPreservesDrawHistory(t *testing.T) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), nil, 0)
	b.initUnitCollision()
	wall := b.Structures[0]
	var gate Structure
	for _, s := range b.Structures {
		if s.Kind == KindGate {
			gate = s
			break
		}
	}
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, X: wall.X, Y: wall.Y}
	cost := b.occupancyCost()
	if cost(wall.X, wall.Y, PlaneLow) != 100 || cost(gate.X, gate.Y, PlaneLow) != 50 || cost(wall.X, wall.Y, PlaneHigh) != 0 {
		t.Fatal("應讀保存的城壁100／門50，不能用目前兵位置重建")
	}
	b.drawUnitCollision(0, 1)
	if cost(wall.X, wall.Y, PlaneLow) != 8 {
		t.Fatal("重畫應直接覆寫結構物成本")
	}
	s.X--
	b.drawUnitCollision(0, 1)
	if cost(wall.X, wall.Y, PlaneLow) != 0 {
		t.Fatal("離開舊格只清零，不恢復結構物成本")
	}
	b.unitCollision.drawSlot(2, s.X, s.Y, 0, PlaneLow)
	s.X--
	b.drawUnitCollision(0, 1)
	if cost(s.X+1, s.Y, PlaneLow) != 0 {
		t.Fatal("其他兵仍佔舊格時也直接清零")
	}
}

func TestCostRedrawAndRemovalUseDifferentPlaneFields(t *testing.T) {
	var g collisionGrid
	g.drawSlot(0, 10, 20, 2, PlaneLow)
	low, high := 20*Width+10, 0x1000+20*Width+10
	if g.cost[low] != 8 || g.cost[high] != 0 {
		t.Fatal("坡地Z2仍是低平面成本")
	}
	g.drawSlot(0, 10, 20, 4, PlaneHigh)
	if g.cost[low] != 0 || g.cost[high] != 8 {
		t.Fatal("重畫要清舊平面並寫新平面")
	}
	// 1B3F3讀目前+1E；若移除前換回低平面，不能誤用舊+1F。
	g.cost[low] = 50
	g.removeSlot(0, PlaneLow)
	if g.cost[low] != 0 || g.cost[high] != 8 || g.registered[0] {
		t.Fatal("移除須使用目前平面")
	}
}

func TestCostSwapKeepsOldPlaneInOriginalSlot(t *testing.T) {
	var g collisionGrid
	g.drawSlot(0, 10, 20, 0, PlaneLow)
	g.drawSlot(1, 11, 20, 4, PlaneHigh)
	a, _ := collisionCell(10, 20, 0)
	b, _ := collisionCell(11, 20, 4)
	before := g.cost
	g.swapSlots(0, 1, a, b)
	if g.old[0] != b || g.old[1] != a || g.oldPlane[0] != PlaneLow || g.oldPlane[1] != PlaneHigh || g.cost != before {
		t.Fatal("換位只交換舊格，舊平面不交換，也不直接改成本")
	}
}

func TestStructureBreakClearsOnlyLowCost(t *testing.T) {
	for _, mode := range []int{KindWall, KindGate, 3} {
		f, _ := tiledField(32)
		b := NewBattle(f, SyntheticFormations(), nil, 0)
		b.initUnitCollision()
		var target Structure
		for _, s := range b.Structures {
			if s.Kind == mode || mode == 3 && s.Kind == KindWall {
				target = s
				break
			}
		}
		index := target.Y*Width + target.X
		b.unitCollision.cost[index], b.unitCollision.cost[index+0x1000] = 8, 8
		if mode == 3 {
			b.SortieBreaksWalls()
		} else {
			b.breakRow(target.Y)
		}
		for n := 0; n < target.Run; n++ {
			if b.unitCollision.cost[index+n*Width] != 0 {
				t.Fatal("破壞須按Run清低平面")
			}
		}
		if b.unitCollision.cost[index+0x1000] != 8 {
			t.Fatal("破壞不可清高平面")
		}
		if mode == 3 {
			for _, s := range b.Structures {
				if s.Kind == KindGate && b.unitCollision.cost[s.Y*Width+s.X] != 50 {
					t.Fatal("守方出擊不可清未破門的成本")
				}
			}
		}
	}
}

func TestCostIncludesStructureBorders(t *testing.T) {
	var g collisionGrid
	structures := []Structure{{Kind: KindWall, X: 0, Y: 0, Run: 2}, {Kind: KindGate, X: 63, Y: 63, Run: 1}}
	g.seedStructureCosts(structures)
	if g.cost[0] != 100 || g.cost[Width] != 100 || g.cost[4095] != 50 {
		t.Fatal("結構成本要包含完整64×64邊框")
	}
	b := &Battle{unitCollision: &g}
	for i := range structures {
		b.clearStructureCosts(&structures[i])
	}
	if g.cost[0] != 0 || g.cost[Width] != 0 || g.cost[4095] != 0 {
		t.Fatal("邊框破壞也須清成本")
	}
}

func TestRetreatExitClearsNativeCostAndCollision(t *testing.T) {
	for side := 0; side < 2; side++ {
		b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
		edge := MinCoord
		b.Sides[side].Mirror = side == 1
		if side == 1 {
			edge = MaxCoord
		}
		s := &b.Sides[side].Soldiers[1]
		*s = Soldier{Alive: true, X: edge, Y: 20, Cmd: Retreat, Next: Retreat}
		b.initUnitCollision()
		b.drawUnitCollision(side, 1)
		b.doRetreat(side, 1)
		if s.Alive || b.unitCollision.registered[b.collisionSlot(side, 1)] || b.occupancyCost()(edge, 20, PlaneLow) != 0 || b.unitCollision.at(edge, 20, 0) != -1 {
			t.Fatal("退卻到出口須立即清成本與碰撞，不等死亡單位的移動批次")
		}
	}
}

func TestPathCostLookupUsesDestinationPlane(t *testing.T) {
	f := flatField()
	f.nav[PlaneHigh], f.lvl[PlaneHigh] = f.nav[PlaneLow], f.lvl[PlaneLow]
	from, to := Point{X: 5, Y: 20}, Point{X: 10, Y: 20}
	cost := func(x, y, plane int) int {
		if plane == PlaneHigh && x == 7 && y == 20 {
			return 50
		}
		return 0
	}
	low := f.findPathForPlaneCosts(from, to, PlaneLow, PlaneLow, false, cost, nil)
	high := f.findPathForPlaneCosts(from, to, PlaneHigh, PlaneHigh, false, cost, nil)
	if len(low) != 1 || len(high) <= 1 {
		t.Fatalf("同XY高平面應繞開成本50：低%v，高%v", low, high)
	}
}
