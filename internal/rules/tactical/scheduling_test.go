package tactical

import "testing"

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
