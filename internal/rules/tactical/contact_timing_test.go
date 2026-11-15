package tactical

import "testing"

type contactPhaseRand struct {
	calls int
	value int
}

func (r *contactPhaseRand) Next() int { r.calls++; return r.value }

func contactPhaseBattle(side, player int, kind Kind, cmd Command, distance int) (*Battle, *contactPhaseRand, int) {
	rows := make([][]int, Height)
	for y := range rows {
		rows[y] = make([]int, Width)
	}
	r := &contactPhaseRand{value: 127}
	b := NewBattle(NewField(rows, 0), SyntheticFormations(), r, 0)
	b.PlayerSide = player
	b.Advantage[0] = Even
	k := 1
	if kind == General {
		k = 0
	}
	b.Sides[side].Soldiers[k] = Soldier{Alive: true, Kind: kind, HP: 200, Power: 32, Stamina: 80, X: 20, Y: 20, StepX: 20, StepY: 20, GoalX: 20, GoalY: 20, Cmd: cmd, Next: cmd, Target: 0}
	b.Sides[1-side].Soldiers[0] = Soldier{Alive: true, Kind: Infantry, HP: 200, X: 20 + distance, Y: 20, Cmd: Guard, Next: Guard, Target: -1}
	return b, r, k
}

func TestContactCommandPhaseDoesNotStrike(t *testing.T) {
	for side := 0; side < 2; side++ {
		for player := 0; player < 2; player++ {
			for _, kind := range []Kind{General, Cavalry, Infantry} {
				for _, cmd := range []Command{Attack, Charge, Guard} {
					if kind == General && cmd == Guard {
						continue
					}
					for _, distance := range []int{0, 1, 3} {
						b, r, k := contactPhaseBattle(side, player, kind, cmd, distance)
						b.updateSoldierCommand(side, k)
						if r.calls != 0 || b.Sides[1-side].Soldiers[0].HP != 200 || b.Sides[side].Soldiers[k].HitGeneral {
							t.Fatalf("命令批次只更新追敵目標，不得額外近戰：side%d player%d kind%d cmd%d distance%d", side, player, kind, cmd, distance)
						}
					}
				}
			}
		}
	}
}

func TestContactStrikeWaitsForMovementAttempt(t *testing.T) {
	for side := 0; side < 2; side++ {
		for player := 0; player < 2; player++ {
			b, r, k := contactPhaseBattle(side, player, Infantry, Attack, 1)
			s := &b.Sides[side].Soldiers[k]
			e := &b.Sides[1-side].Soldiers[0]
			b.updateSoldierCommand(side, k)
			b.updateSoldierMovement(side, k)
			if r.calls != 0 || e.HP != 200 || s.StepX != 21 || s.X != 20 {
				t.Fatal("原位取下一目標的當拍不可偷打")
			}
			b.updateSoldierMovement(side, k)
			if r.calls != 1 || e.HP != 168 || s.X != 20 || !s.MoveFlag || !s.HitGeneral || e.ProjectileCooldown != 8 {
				t.Fatalf("真正移動碰撞才打一次並重設受擊冷卻：calls%d actor%+v target%+v", r.calls, *s, *e)
			}
			b, r, k = contactPhaseBattle(side, player, Infantry, Attack, 1)
			b.Sides[side].Soldiers[k].TargetRefresh = true
			b.updateSoldierCommand(side, k)
			if r.calls != 0 {
				t.Fatal("刷新中繼點也不得在命令階段近戰")
			}
			b.updateSoldierMovement(side, k)
			if r.calls != 1 || b.Sides[1-side].Soldiers[0].HP != 168 {
				t.Fatal("刷新已對準敵格時，同拍移動碰撞仍可攻擊一次")
			}
			for _, swapped := range []bool{false, true} {
				b, r, k = contactPhaseBattle(side, player, Infantry, Attack, 1)
				s = &b.Sides[side].Soldiers[k]
				s.Stun = 2
				s.Swapped = swapped
				b.updateSoldierCommand(side, k)
				b.updateSoldierMovement(side, k)
				if r.calls != 0 || b.Sides[1-side].Soldiers[0].HP != 200 {
					t.Fatal("命令更新不能繞過硬直／換位的移動閘偷打")
				}
			}
		}
	}
}

func TestContactCommandStillLaunchesProjectiles(t *testing.T) {
	for _, kind := range []Kind{Archer, Infantry} {
		b, r, k := contactPhaseBattle(0, 0, kind, Attack, 1)
		s := &b.Sides[0].Soldiers[k]
		wantCalls := 1
		if kind == Infantry {
			s.Z = 4
			s.PlaneHigh = PlaneHighElevated
			s.OnWall = true
			wantCalls = 0
		}
		b.updateSoldierCommand(0, k)
		views := b.Projectiles()
		if len(views) != 1 || r.calls != wantCalls || s.Stun != 2 || views[0].Special != (kind == Infantry) {
			t.Fatal("命令階段仍須保留普通箭與特殊投射物")
		}
	}
}

func TestContactMeleeCooldownDiffersFromArrowAndGeneral(t *testing.T) {
	for _, kind := range []Kind{Cavalry, Archer, Infantry} {
		for _, initial := range []uint8{0, 3, 8} {
			for _, hit := range []bool{false, true} {
				b, r, k := contactPhaseBattle(0, 0, Infantry, Attack, 1)
				e := &b.Sides[1].Soldiers[0]
				e.Kind = kind
				e.ProjectileCooldown = initial
				if !hit {
					r.value = 0
					b.Sides[0].Soldiers[k].Power = 0
				}
				b.meleeHit(0, &b.Sides[0].Soldiers[k], e)
				want := initial
				if hit {
					want = 8
				}
				if e.ProjectileCooldown != want {
					t.Fatalf("一般近戰命中才重設冷卻8：kind%d hit%v initial%d got%d", kind, hit, initial, e.ProjectileCooldown)
				}
			}
			b, _, _ := contactPhaseBattle(0, 0, Infantry, Attack, 1)
			e := &b.Sides[1].Soldiers[0]
			e.Kind = kind
			e.ProjectileCooldown = initial
			b.hitByArrow(0, e, 32)
			if e.ProjectileCooldown != initial {
				t.Fatal("箭傷保留原發射冷卻")
			}
		}
	}
	b, r, k := contactPhaseBattle(0, 0, Infantry, Attack, 1)
	r.value = 0
	e := &b.Sides[1].Soldiers[0]
	e.Kind = General
	e.Cmd = Attack
	e.ProjectileCooldown = 3
	b.hitGeneral(&b.Sides[0].Soldiers[k], e)
	if e.ProjectileCooldown != 3 || e.HP >= 200 {
		t.Fatal("大將近戰受擊不重設發射冷卻")
	}
}
