package tactical

import "testing"

func TestWallShootingUsesBothTerrainConditions(t *testing.T) {
	for side := 0; side < 2; side++ {
		for _, k := range []int{0, 1} {
			for _, terrain := range []bool{false, true} {
				for _, high := range []bool{false, true} {
					b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
					s := &b.Sides[side].Soldiers[k]
					*s = Soldier{Alive: true, Kind: Archer, HP: 100, X: 44, Y: 38, Z: 4,
						GoalX: 44, GoalY: 38, GoalZ: 4, StepX: 42, StepY: 38, StepZ: 4,
						HighTerrain: terrain, OnWall: high, Stamina: 128, Cmd: ScaleWal, Next: ScaleWal, Target: -1}
					if high {
						s.PlaneHigh = PlaneHighElevated
					}
					b.Sides[1-side].Soldiers[0] = Soldier{Alive: true, Kind: Infantry, HP: 100, X: 5, Y: 38}
					b.updateSoldierCommand(side, k)
					if !terrain || !high {
						if s.Cmd != ScaleWal || s.Next != ScaleWal || s.Stamina != 128 || len(b.projectiles) != 0 {
							t.Fatalf("缺少城上條件不可提前轉攻擊：%+v", *s)
						}
						continue
					}
					if s.Cmd != ScaleWal || s.Next != wallAttack || s.Stamina != 40 || s.ProjectileCooldown != 8 || s.Stun != 2 || !s.Swapped || !s.HitGeneral || s.Hurt || len(b.projectiles) != 1 {
						t.Fatalf("登城同拍射箭應只改Next6並設發射硬直：%+v", *s)
					}
					b.updateSoldierMovement(side, k)
					if s.Stamina != 40 || s.Stun != 2 || s.X != 44 {
						t.Fatal("發射當拍應略過移動，不扣移動疲勞")
					}
					b.updateSoldierCommand(side, k)
					if s.Cmd != wallAttack || s.Next != wallAttack || s.ProjectileCooldown != 7 || len(b.projectiles) != 1 {
						t.Fatal("命令6須繼續走同一處理器，冷卻只遞減")
					}
				}
			}
		}
	}
}

func TestWallShootingArcherGoalRefreshBranches(t *testing.T) {
	for _, refresh := range []bool{false, true} {
		b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
		path := &Waypoints{pts: []Point{{X: 12, Y: 13}}}
		s := &b.Sides[0].Soldiers[1]
		*s = Soldier{Alive: true, Kind: Archer, X: 44, Y: 38, Z: 4, Target: 1,
			StepX: 12, StepY: 13, StepZ: 2, GoalZ: 4, Path: path, TargetRefresh: refresh, ProjectileCooldown: 1}
		b.Sides[1].Soldiers[1] = Soldier{Alive: true, Kind: Infantry, X: 5, Y: 38}
		b.doAttackAction(0, 1)
		if refresh {
			if s.GoalX != 44 || s.GoalY != 38 || s.StepX != 44 || s.StepY != 38 || s.StepZ != 4 || s.Path != nil || s.TargetRefresh {
				t.Fatal("刷新弓兵須用自身Goal／Step並清路")
			}
		} else if s.GoalX != 5 || s.GoalY != 38 || s.StepX != 12 || s.StepY != 13 || s.StepZ != 2 || s.Path != path {
			t.Fatal("未刷新弓兵只改敵方GoalXY，保留Step與路徑")
		}
	}
}

func TestWallShootingAimGoalBounds(t *testing.T) {
	for _, tc := range []struct {
		sx, sy, tx, ty, wantX, wantY int
	}{
		{60, 20, 54, 20, 67, 20}, {60, 20, 55, 20, 55, 20},
		{1, 20, 10, 20, 253, 20}, {1, 20, 9, 20, 9, 20},
		{20, 60, 20, 54, 20, 67}, {20, 60, 20, 55, 20, 55},
		{20, 1, 20, 10, 20, 253}, {20, 1, 20, 9, 20, 9},
	} {
		b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
		s := Soldier{X: tc.sx, Y: tc.sy, StepX: 17, StepY: 18, StepZ: 4, GoalZ: 6}
		e := Soldier{X: tc.tx, Y: tc.ty, Z: 3}
		b.setArcherAimGoal(&s, &e)
		if s.GoalX != tc.wantX || s.GoalY != tc.wantY || s.StepX != 17 || s.StepY != 18 || s.StepZ != 4 || s.GoalZ != 6 {
			t.Fatalf("四向瞄準的10／54邊界或byte結果不符：%+v，案例%+v", s, tc)
		}
	}
}

func TestWallShootingSameCellKeepsNativeWestDirection(t *testing.T) {
	if attackDirection(20, 20, 20, 20) != West {
		t.Fatal("1ACA4兩軸距皆0時保留初始西向")
	}
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, Kind: Archer, X: 20, Y: 20}
	e := Soldier{Alive: true, X: 20, Y: 20}
	b.shoot(0, 1, &e)
	if s.Facing != West || s.GoalX != 30 || s.GoalY != 20 || len(b.projectiles) != 1 || b.projectiles[0].x != 19 {
		t.Fatal("同格瞄準及普通箭須沿原版西向分支")
	}
}

func TestWallShootingSuccessPoseAndCooldown(t *testing.T) {
	for _, special := range []bool{false, true} {
		b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
		s := &b.Sides[0].Soldiers[1]
		*s = Soldier{Alive: true, Kind: Archer, X: 44, Y: 38, Z: 4, PoseStep: 1, Hurt: true, Stun: 9}
		e := &b.Sides[1].Soldiers[1]
		*e = Soldier{Alive: true, Kind: Infantry, X: 5, Y: 38}
		cooldown := uint8(8)
		if special {
			b.launchSpecialProjectile(0, s, e)
			cooldown = 6
		} else {
			b.launchNormalProjectile(0, s, e)
		}
		if s.PoseStep != 1 || s.Hurt || !s.HitGeneral || !s.Swapped || s.Stun != 2 || s.ProjectileCooldown != cooldown || s.Facing != West {
			t.Fatalf("發射成功應保留姿勢低位、清受擊圖號並設動作硬直：%+v", *s)
		}
		if !special {
			s.ProjectileCooldown = 2
			b.shoot(0, 1, e)
			b.shoot(0, 1, e)
			if s.ProjectileCooldown != 0 || len(b.projectiles) != 1 || !s.HitGeneral {
				t.Fatal("冷卻呼叫只遞減，不重發、不提前清攻擊圖號")
			}
			b.shoot(0, 1, e)
			if len(b.projectiles) != 2 || s.ProjectileCooldown != 8 {
				t.Fatal("距離39仍可在冷卻結束後發普通箭，不套額外20射程")
			}
		}
	}
}

func TestWallShootingAttackPoseClearingConsumers(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, Kind: Archer, X: 20, Y: 20, Cmd: Holding, Next: Holding, Stun: 2, Hurt: true, HitGeneral: true, PoseStep: 1}
	b.updateSoldierCommand(0, 1)
	b.updateSoldierMovement(0, 1)
	if !s.HitGeneral || !s.Hurt || s.Stun != 1 {
		t.Fatal("命令入口及尚未歸零的硬直不得清圖號")
	}
	b.updateSoldierMovement(0, 1)
	if s.HitGeneral || s.Hurt || s.Stun != 0 || s.PoseStep != 1 {
		t.Fatal("硬直歸零須清圖號並沿既有每拍翻面")
	}
	s.HitGeneral = true
	s.StepX, s.StepY, s.StepZ = s.X, s.Y, s.Z
	b.moveToward(0, 1)
	if s.HitGeneral {
		t.Fatal("路徑耗盡須清攻擊圖號")
	}
	for _, general := range []bool{false, true} {
		e := Soldier{Alive: true, HP: 100, Kind: Archer, Cmd: Attack, PoseStep: 1, HitGeneral: true}
		if general {
			e.Kind = General
			attacker := Soldier{Alive: true, Power: 16}
			if !b.hitGeneral(&attacker, &e) {
				t.Fatal("固定亂數應命中大將")
			}
		} else {
			b.applyHit(&e, 1)
		}
		if e.HitGeneral || !e.Hurt || e.PoseStep != 1 {
			t.Fatal("受擊先清舊攻擊圖號，保留姿勢低位再設Hurt")
		}
	}
}
