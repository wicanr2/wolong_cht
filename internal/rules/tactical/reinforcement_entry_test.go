package tactical

import "testing"

// DOS/V 1B413、1B4EA：先移至入口，再測兩層低七位；成功只重設指定欄位。
func TestReinforcementEntryRetainsRecordAndRetries(t *testing.T) {
	for player := 0; player < 2; player++ {
		for side := 0; side < 2; side++ {
			for _, k := range []int{8, 9} {
				b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{127}}, 0)
				b.SetPlayerSide(player)
				s := &b.Sides[side].Soldiers[k]
				*s = Soldier{Kind: Cavalry, Power: 17, Stamina: 80, ProjectileCooldown: 11, Target: 3, PathQueued: true, MoveFlag: true, Swapped: true, TargetRefresh: true, Hurt: true, HitGeneral: true, Stun: 2, PoseStep: 1, X: 25, Y: 63, Z: 3, StepX: 21, StepY: 31, StepZ: 2, GoalX: 22, GoalY: 32, GoalZ: 3, PlaneHigh: PlaneHighGround, Cmd: ScaleWal, Next: wallAttack}
				b.Sides[side].Morale = 73
				b.Sides[side].Reserve[1] = 2
				b.unitCollision = &collisionGrid{}
				x := 1
				if b.Role(side) == 1 {
					x = 62
				}
				cell := 47*Width + x
				b.unitCollision.cells[cell+0x1000] = 100
				b.reinforceSoldier(side, k)
				if s.Alive || s.X != x || s.Y != 47 || s.Z != 0 || s.PlaneHigh != 0 || b.Sides[side].Reserve[1] != 2 || s.Stun != 2 || s.HP != 0 {
					t.Fatalf("blocked entrance must retain rewritten position without consuming stock: %+v", s)
				}
				b.unitCollision.cells[cell+0x1000] = 128
				b.reinforceSoldier(side, k)
				if !s.Alive || s.HP != 73 || s.Stun != 0 || s.Kind != Cavalry || s.Power != 17 || s.Stamina != 80 || s.ProjectileCooldown != 11 || s.Target != 3 || !s.PathQueued || !s.MoveFlag || !s.TargetRefresh || s.Swapped || s.Hurt || s.HitGeneral || s.PoseStep != 1 || s.Cmd != Duel || s.Next != ScaleWal || s.StepX != 21 || s.StepY != 31 || s.StepZ != 2 || s.GoalX != 22 || s.GoalY != 32 || s.GoalZ != 3 || b.Sides[side].Reserve[1] != 1 {
					t.Fatalf("birth must preserve record and normalize command: %+v", s)
				}
				d := b.DrawnSoldier(side, k)
				if d.PoseStep != 0 || d.X != x || d.Y != 47 {
					t.Fatal("birth is drawn before pose flips", d)
				}
			}
		}
	}
}

func TestReinforcementNormalizesOldOrder(t *testing.T) {
	for _, c := range []struct{ old, next, want Command }{{ScaleWal, Retreat, ScaleWal}, {wallAttack, Retreat, ScaleWal}, {Holding, Retreat, Form}, {Attack, wallAttack, ScaleWal}, {Charge, Holding, Form}, {Charge, Guard, Guard}} {
		b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
		b.Sides[0].Morale = 73
		b.Sides[0].Reserve[1] = 1
		b.Sides[0].Soldiers[9] = Soldier{Kind: Infantry, Y: 15, Cmd: c.old, Next: c.next}
		b.reinforceSoldier(0, 9)
		s := b.Sides[0].Soldiers[9]
		if !s.Alive || s.Y != 16 || s.Cmd != Duel || s.Next != c.want {
			t.Fatal(c, s)
		}
	}
}

func TestMovementBirthPrecedesLaterFriendlyMove(t *testing.T) {
	for player := 0; player < 2; player++ {
		for side := 0; side < 2; side++ {
			b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{127}}, 0)
			b.SetPlayerSide(player)
			for i := range b.Sides {
				b.Sides[i].Morale = 200
				for k := range b.Sides[i].Soldiers {
					x := 20
					if b.Role(i) == 1 {
						x = 55
					}
					y := k + 1
					b.Sides[i].Soldiers[k] = Soldier{Alive: true, Kind: Infantry, HP: 200, Stamina: 80, Stun: 255, X: x, Y: y, StepX: x, StepY: y, GoalX: x, GoalY: y, Cmd: Duel, Next: Duel, Target: -1}
				}
			}
			edge, near := 1, 2
			if b.Role(side) == 1 {
				edge, near = 62, 61
			}
			born := &b.Sides[side].Soldiers[9]
			born.Alive = false
			born.HP = 0
			born.Stun = 0
			born.Y = 30
			b.Sides[side].Reserve[1] = 1
			mover := &b.Sides[side].Soldiers[10]
			mover.X, mover.Y, mover.Stun = near, 30, 0
			mover.StepX, mover.StepY = edge, 30
			mover.GoalX, mover.GoalY = edge, 30
			b.Step()
			if !born.Alive || born.X != near || !born.Swapped || mover.X != edge || b.Sides[side].Reserve[1] != 0 {
				t.Fatalf("player%d side%d: later mover must swap with earlier newborn: born%+v mover%+v", player, side, born, mover)
			}
			d := b.DrawnSoldier(side, 9)
			if d.X != edge || d.PoseStep != 0 {
				t.Fatal("later swap cannot rewrite earlier drawing", d)
			}
		}
	}
}

func TestMultipleBirthsUseDistinctEntrancesInOneSquad(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
	b.unitCollision = &collisionGrid{}
	b.Sides[0].Reserve[1] = 2
	b.Sides[0].Morale = 73
	b.Sides[0].Soldiers[8] = Soldier{Kind: Cavalry, Y: 20}
	b.Sides[0].Soldiers[9] = Soldier{Kind: Infantry, Y: 21}
	b.reinforce()
	if !b.Sides[0].Soldiers[8].Alive || !b.Sides[0].Soldiers[9].Alive || b.Sides[0].Reserve[1] != 0 {
		t.Fatal("same squad births are limited by occupancy and reserve")
	}
}

// 1A953／1A96D／1A988 → 1A8DE；重複令不重傳，不跳過退卻或不在場隊員。
func TestChangedLeaderAttackOrdersReachAllSevenMembers(t *testing.T) {
	for _, cmd := range []Command{Attack, Charge, ScaleWal} {
		for _, same := range []bool{false, true} {
			b := newTestBattle(walledField(32))
			s := &b.Sides[0].Soldiers[8]
			s.Kind = Cavalry
			s.Cmd = Duel
			s.Next = cmd
			if same {
				s.Cmd = cmd
			}
			for k := 9; k < 16; k++ {
				m := &b.Sides[0].Soldiers[k]
				m.Next = Holding
				m.TargetRefresh = false
			}
			b.Sides[0].Soldiers[9].Alive = false
			b.Sides[0].Soldiers[10].Cmd = Retreat
			b.updateSoldierCommand(0, 8)
			for k := 9; k < 16; k++ {
				m := b.Sides[0].Soldiers[k]
				want := cmd
				if same {
					want = Holding
				}
				if m.Next != want || m.TargetRefresh == same {
					t.Fatal(cmd, same, k, m)
				}
			}
			if b.Sides[0].Soldiers[10].Cmd != Retreat {
				t.Fatal("broadcast cannot interrupt retreat")
			}
		}
	}
}

func TestAttackGeneralReturnsToFormationAndKeepsCommand(t *testing.T) {
	b := newTestBattle(flatField())
	s := &b.Sides[0].Soldiers[0]
	s.Cmd, s.Next = Attack, Attack
	s.Stamina = 80
	s.TargetRefresh = true
	s.X, s.Y = b.formationSpot(0, 0)
	b.updateSoldierCommand(0, 0)
	if s.Cmd != Attack || s.Next != Attack || s.Stamina != StaminaFull || !s.TargetRefresh || s.GoalX != s.X || s.GoalY != s.Y {
		t.Fatal("general attack uses formation without consuming refresh", s)
	}
}

func TestExitAndDeathCompletionClearPresenceFlagsOnly(t *testing.T) {
	b := newTestBattle(flatField())
	s := &b.Sides[0].Soldiers[9]
	s.X, s.Y = MinCoord, 30
	s.Cmd, s.Next = Retreat, ScaleWal
	s.PathQueued = true
	s.TargetRefresh = true
	s.Swapped = true
	s.MoveFlag = true
	s.HighTerrain = true
	s.Stun = 2
	s.Hurt = true
	s.HitGeneral = true
	s.PoseStep = 1
	b.doRetreat(0, 9)
	if s.Alive || s.TargetRefresh || s.Swapped || s.MoveFlag || s.HighTerrain || !s.PathQueued || s.Stun != 2 || !s.Hurt || !s.HitGeneral || s.PoseStep != 1 || s.Cmd != Retreat || s.Next != ScaleWal || b.Sides[0].Escaped[1] != 1 {
		t.Fatal("exit masks only byte0", s)
	}
	b = newTestBattle(flatField())
	s = &b.Sides[0].Soldiers[9]
	b.beginDeath(s, snapshotDeathPose(s))
	s.TargetRefresh = true
	s.MoveFlag = true
	s.Swapped = true
	s.PathQueued = true
	for i := 0; i < 4; i++ {
		b.stepDeath(0, 9)
	}
	if s.Alive || s.TargetRefresh || s.Swapped || s.MoveFlag || s.HighTerrain || !s.PathQueued {
		t.Fatal("death completion reapplies presence mask", s)
	}
}
