package tactical

import "testing"

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
