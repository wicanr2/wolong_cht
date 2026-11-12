package tactical

import "testing"

func TestNewScaleOrderSetsWaypointAndPreservesLaterPath(t *testing.T) {
	lib := realLib(t)
	for side := 0; side < 2; side++ {
		for _, k := range []int{32, 33} { // 隊長及隊員。
			for _, kind := range []Kind{Infantry, Archer} {
				b := NewBattle(realField(t, lib, 76), SyntheticFormations(), nil, 0)
				s := &b.Sides[side].Soldiers[k]
				*s = Soldier{Alive: true, Kind: kind, X: 60, Y: 35, Cmd: Holding, Next: ScaleWal,
					StepX: 59, StepY: 34, StepZ: 1, Target: -1}
				b.updateSoldierCommand(side, k)
				if s.Cmd != ScaleWal || s.StepX != 42 || s.StepY != 35 || s.StepZ != 4 || s.GoalX != 42 || s.GoalY != 35 {
					t.Fatalf("側%d槽%d新登城中繼點不符原版第137拍：%+v", side, k, *s)
				}
				b.updateSoldierMovement(side, k)
				if s.X != 59 {
					t.Fatalf("新登城命令首拍未朝城牆移動：%+v", *s)
				}
				s.StepX, s.StepY, s.StepZ = 58, 34, 3
				path := &Waypoints{pts: []Point{{X: 57, Y: 34}}}
				s.Path = path
				b.updateSoldierCommand(side, k)
				if s.StepX != 58 || s.StepY != 34 || s.StepZ != 3 || s.Path != path || s.Path.Len() != 1 {
					t.Fatalf("持續登城不可重設中繼點與繞路：%+v", *s)
				}
			}
		}
	}
}

func TestSwapWaypointHeightAndTerrainFlags(t *testing.T) {
	// 原版1B783–1B795；同高坡地也交換，不能改成只補Z差。
	tests := []struct {
		name              string
		z, otherZ         int
		plane, otherPlane byte
		want, otherWant   int
	}{
		{"flat low plane", 0, 0, 0, 0, 2, 4},
		{"ascending low slope", 0, 1, 0, 0, 4, 2},
		{"descending low slope", 1, 0, 0, 0, 4, 2},
		{"equal height low slope", 1, 1, 0, 0, 4, 2},
		{"both high plane", 4, 5, 0x10, 0x10, 2, 4},
		{"first high plane", 4, 1, 0x10, 0, 2, 4},
		{"second high plane", 1, 4, 0, 0x10, 2, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for side := 0; side < 2; side++ {
				b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
				me, other := &b.Sides[side].Soldiers[1], &b.Sides[side].Soldiers[2]
				*me = Soldier{Alive: true, Kind: Infantry, X: 10, Y: 20, Z: tt.z,
					StepZ: 2, PlaneHigh: tt.plane, HighTerrain: true}
				*other = Soldier{Alive: true, Kind: Infantry, X: 11, Y: 20, Z: tt.otherZ,
					StepZ: 4, PlaneHigh: tt.otherPlane}
				if !b.swapWith(side, 1, side, 2) {
					t.Fatal("可換位的步兵被拒絕")
				}
				if me.StepZ != tt.want || other.StepZ != tt.otherWant {
					t.Fatalf("中繼高度=%d/%d，原版=%d/%d", me.StepZ, other.StepZ, tt.want, tt.otherWant)
				}
				if me.X != 11 || other.X != 10 || me.Z != tt.otherZ || other.Z != tt.z ||
					me.PlaneHigh != tt.otherPlane || other.PlaneHigh != tt.plane || !other.Swapped {
					t.Fatal("既有位置、平面或被換位旗標未正確交換")
				}
				if !me.HighTerrain || other.HighTerrain {
					t.Fatal("換位函式不應交換地形旗標")
				}
			}
		})
	}
}

func TestSwappedTerrainFlagRefreshesAtEachDraw(t *testing.T) {
	rows := make([][]int, Height)
	for y := range rows {
		rows[y] = make([]int, Width)
	}
	rows[20][10] = 4
	b := NewBattle(NewField(rows, 0), SyntheticFormations(), nil, 0)
	me, other := &b.Sides[0].Soldiers[1], &b.Sides[0].Soldiers[2]
	*me = Soldier{Alive: true, Kind: Infantry, X: 10, Y: 20, HighTerrain: true}
	*other = Soldier{Alive: true, Kind: Infantry, X: 11, Y: 20}
	if !b.swapWith(0, 1, 0, 2) {
		t.Fatal("換位失敗")
	}
	b.drawUnitCollision(0, 1)
	if me.HighTerrain || other.HighTerrain {
		t.Fatal("主動換位者應先重算，被換者尚未輪到時保留原旗標")
	}
	b.updateSoldierMovement(0, 2)
	b.drawUnitCollision(0, 2)
	if !other.HighTerrain || other.Swapped || other.X != 10 {
		t.Fatal("被換者跳過移動後仍應重畫並重算地形旗標")
	}
}

func TestLockOnUsesOriginalYDirectionAndScanOrder(t *testing.T) {
	tests := []struct {
		name     string
		meX, meY int
		foes     [2][2]int
		want     int
	}{
		{"larger Y adds distance twice", 1, 28, [2][2]int{{62, 29}, {62, 27}}, 1},
		{"weighted tie retains earlier slot", 62, 37, [2][2]int{{1, 35}, {1, 38}}, 0},
		{"X and Y both affect selection", 10, 20, [2][2]int{{12, 19}, {10, 22}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for side := 0; side < 2; side++ {
				b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
				b.Sides[side].Soldiers[0] = Soldier{Alive: true, Kind: Infantry, X: tt.meX, Y: tt.meY}
				for i, xy := range tt.foes {
					b.Sides[1-side].Soldiers[i] = Soldier{Alive: true, X: xy[0], Y: xy[1]}
				}
				b.lockOnNearest(side, 0)
				if got := b.Sides[side].Soldiers[0].Target; got != tt.want {
					t.Fatalf("側%d鎖定槽%d，原版應選槽%d", side, got, tt.want)
				}
			}
		})
	}
}

func TestHoldingLeaderKeepsStaminaAndTracksFormationLine(t *testing.T) {
	for _, side := range []int{AttackerSide, DefenderSide} {
		for k := 0; k < SoldiersOnFoot; k += PerSquad {
			b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
			s := &b.Sides[side].Soldiers[k]
			*s = Soldier{Alive: true, Kind: Infantry, Cmd: Holding, Next: Holding, Stamina: 37, Target: -1}
			b.Place()
			b.updateSoldierCommand(side, k)
			if s.Stamina != StaminaFull || s.Cmd != Holding || s.Next != Holding {
				t.Fatalf("側%d隊長%d就位後應保持命令並補餘力：%+v", side, k, *s)
			}
			x := s.GoalX
			if side == AttackerSide {
				b.Sides[side].Line += 5
			} else {
				b.Sides[side].Line -= 5
			}
			b.updateSoldierCommand(side, k)
			if s.GoalX == x {
				t.Fatalf("側%d隊長%d未追隨新的陣形線", side, k)
			}
			start := s.X
			b.updateSoldierMovement(side, k)
			b.updateSoldierMovement(side, k)
			if s.X == start {
				t.Fatalf("側%d隊長%d未走向新的陣形位置", side, k)
			}
		}
	}
}

func TestHoldingFollowerKeepsOriginalEmptyDispatch(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, Kind: Infantry, Cmd: Holding, Next: Holding, Stamina: 37, GoalX: 22, GoalY: 23}
	b.updateSoldierCommand(0, 1)
	if s.Stamina != 37 || s.GoalX != 22 || s.GoalY != 23 {
		t.Fatalf("隊員命令7應保持原版空分派：%+v", *s)
	}
}

func TestOriginalFirstStepSelectsGoalBeforeMoving(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	s := &b.Sides[0].Soldiers[0]
	*s = Soldier{Alive: true, Kind: Infantry, X: 1, Y: 46, StepX: 1, StepY: 46,
		GoalX: 3, GoalY: 32, Stamina: 128, Target: -1}
	b.moveToward(0, 0)
	if s.X != 1 || s.Y != 46 || s.StepX != 3 || s.StepY != 32 || s.Stamina != 127 {
		t.Fatalf("原版首拍欄位不符：%+v", *s)
	}
	b.moveToward(0, 0)
	if s.X != 2 || s.Y != 46 || s.Stamina != 126 {
		t.Fatalf("原版第二拍欄位不符：%+v", *s)
	}
}

func TestEnemyCollisionEndsStepWithoutTryingOtherAxis(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, Kind: Infantry, Cmd: Attack, X: 10, Y: 20,
		StepX: 12, StepY: 22, GoalX: 12, GoalY: 22, HP: 100, Power: 10, Stamina: 128}
	e := &b.Sides[1].Soldiers[1]
	*e = Soldier{Alive: true, Kind: Infantry, Cmd: Attack, X: 11, Y: 20, HP: 100, Power: 10}
	b.moveToward(0, 1)
	if s.X != 10 || s.Y != 20 || s.PathQueued || !s.MoveFlag {
		t.Fatalf("撞敵後不可再試 Y 或排尋路：%+v", *s)
	}
}

func TestGeneralReturnBypassesOrdinarySwapGuards(t *testing.T) {
	for _, command := range []Command{Form, Retreat} {
		for _, targetSide := range []int{0, 1} {
			b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
			g := &b.Sides[0].Soldiers[0]
			*g = Soldier{Alive: true, Kind: 0, Cmd: command, X: 10, Y: 20}
			e := &b.Sides[targetSide].Soldiers[1]
			*e = Soldier{Alive: true, Kind: Infantry, X: 11, Y: 20, MoveFlag: true, Swapped: true, Hurt: true, Cmd: Retreat}
			if ok, _ := b.tryMove(0, 0, 11, 20, 0); !ok || g.X != 11 || e.X != 10 {
				t.Fatalf("大將命令 %v 對側 %d 未按原版例外換位", command, targetSide)
			}
		}
	}
}

func TestStationaryMovementConsumesStamina(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	s := &b.Sides[0].Soldiers[0]
	*s = Soldier{Alive: true, Kind: Infantry, X: 10, Y: 20, GoalX: 10, GoalY: 20,
		StepX: 10, StepY: 20, Stamina: 128, Cmd: Holding, Next: Holding, Target: -1}
	b.updateSoldierMovement(0, 0)
	if s.X != 10 || s.Y != 20 || s.Stamina != 127 {
		t.Fatalf("未位移更新後位置／體力 = %d,%d / %d，原版要求 10,20 / 127", s.X, s.Y, s.Stamina)
	}
	s.Stamina = 0
	b.updateSoldierMovement(0, 0)
	if s.Stamina != 0 {
		t.Fatal("零體力不可下溢")
	}
}

func TestStunnedSoldierReceivesOrderBeforeMovement(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{1}}, 0)
	s := &b.Sides[0].Soldiers[0]
	*s = Soldier{Alive: true, Kind: Infantry, X: 10, Y: 20, Stamina: 128,
		Cmd: Holding, Next: Duel, Stun: 2, Target: -1}
	b.updateSoldier(0, 0)
	if s.Cmd != Duel || s.Stun != 1 || s.X != 10 || s.Y != 20 || s.Stamina != 128 {
		t.Fatalf("硬直時應接收命令且跳過移動與耗體力：%+v", *s)
	}
}
