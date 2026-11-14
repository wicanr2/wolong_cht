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
	if s.Facing != West || s.GoalX != 30 || s.GoalY != 20 || len(b.projectiles) != 1 || b.projectiles[0].x != 20 {
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
			if len(b.projectiles) != 1 || s.ProjectileCooldown != 0 {
				t.Fatal("冷卻結束仍須等待自己的投射物槽釋放")
			}
		}
	}
}

func TestWallBatchProjectileSlotsAndFailedLaunch(t *testing.T) {
	for player := 0; player < 2; player++ {
		b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
		b.PlayerSide = player
		e := Soldier{Alive: true, X: 40, Y: 20}
		for side := 0; side < 2; side++ {
			for _, k := range []int{6, 22, 38} {
				s := &b.Sides[side].Soldiers[k]
				*s = Soldier{Alive: true, Kind: Archer, X: 20, Y: 20, Z: 4, Facing: West, PoseStep: 1, Hurt: true, Stun: 9}
				slot, owned := b.projectileSlot(side, s)
				if !owned || slot != b.Role(side)*16+6 {
					t.Fatal("槽位須依玩家／腳本角色及側內k%16")
				}
				rng := b.rng.(*fixedRand)
				n := rng.i
				b.TakeSoundEffects()
				b.launchNormalProjectile(side, s, &e)
				if rng.i != n+1 || s.Facing != East {
					t.Fatal("普通箭失敗前也須更新面向與取亂數")
				}
				if k == 6 {
					if s.ProjectileCooldown != 8 || s.Stun != 2 || !s.Swapped || len(b.TakeSoundEffects()) != 1 {
						t.Fatal("空槽應成功發射")
					}
				} else {
					if s.ProjectileCooldown != 0 || s.Stun != 9 || !s.Hurt || s.Swapped || s.HitGeneral || len(b.TakeSoundEffects()) != 0 {
						t.Fatal("共槽被佔時不得設定成功旗標／硬直／冷卻或音效")
					}
					b.launchSpecialProjectile(side, s, &e)
					if rng.i != n+1 || s.ProjectileCooldown != 0 || len(b.TakeSoundEffects()) != 0 {
						t.Fatal("特殊投射物共槽且不消耗普通箭亂數")
					}
				}
			}
		}
		if len(b.projectiles) != 2 {
			t.Fatal("雙側同編號使用獨立槽")
		}
		for i := range b.projectiles {
			b.projectiles[i].power = 0
		}
		b.stepProjectiles()
		if len(b.projectiles) != 0 {
			t.Fatal("沿既有更新入口移除投射物")
		}
		b.launchSpecialProjectile(0, &b.Sides[0].Soldiers[22], &e)
		if len(b.projectiles) != 1 || b.Sides[0].Soldiers[22].ProjectileCooldown != 6 {
			t.Fatal("釋放後普通槽可由特殊投射物重用")
		}
	}
}

func TestWallBatchNormalProjectileBirthAndFirstStep(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
	s := &b.Sides[0].Soldiers[6]
	*s = Soldier{Alive: true, Kind: Archer, X: 20, Y: 20, Z: 4}
	e := Soldier{Alive: true, X: 40, Y: 20, Z: 4}
	b.launchNormalProjectile(0, s, &e)
	if len(b.projectiles) != 1 {
		t.Fatal("應發射普通箭")
	}
	p := b.projectiles[0]
	if p.x != 20 || p.y != 20 || p.z != 5 || p.previousX != 20 || p.gridIndex != projectileGridIndex(20, 20, 5) {
		t.Fatal("普通箭須出生於射手XY與Z+1")
	}
	b.stepProjectiles()
	if len(b.projectiles) != 1 || b.projectiles[0].x != 21 || b.projectiles[0].previousX != 20 {
		t.Fatal("下一次投射物更新才走第一格")
	}
	c := &b.Sides[0].Soldiers[7]
	*c = *s
	c.ProjectileCooldown = 0
	b.launchSpecialProjectile(0, c, &e)
	if len(b.projectiles) != 2 || !b.projectiles[1].special || b.projectiles[1].x != 21 || b.projectiles[1].y != 20 {
		t.Fatal("特殊投射物仍出生於相鄰格")
	}
}

func TestProjectileHeightCapWritesBackFixedPoint(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
	for _, fraction := range []int{0, 0xe0, 0xff} {
		p := projectile{x: 20, y: 20, z: 5, heightFP: 5<<8 | fraction, velocityFP: 0x180, power: 28, direction: 0x80}
		if !b.advanceProjectile(&p) || p.z != 5 || p.heightFP != 5<<8|fraction {
			t.Fatal("1BA96須保留高度小數並寫回上限，不能只改顯示Z")
		}
		if p.velocityFP != 0x16c {
			t.Fatal("只夾本拍位移，保存速度仍由原值減14h")
		}
		p.velocityFP = -0x180
		if !b.advanceProjectile(&p) || p.heightFP != 4<<8|fraction || p.velocityFP != -0x194 {
			t.Fatal("負速度也只夾本拍位移")
		}
	}
}

func TestProjectileNativeSlotReleaseAfterLanding(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{2}}, 0)
	s := &b.Sides[1].Soldiers[38]
	*s = Soldier{Alive: true, Kind: Archer, HP: 200, X: 44, Y: 25, Z: 4}
	e := Soldier{Alive: true, Kind: Infantry, HP: 200, X: 5, Y: 24}
	b.launchNormalProjectile(1, s, &e)
	if len(b.projectiles) != 1 || b.projectiles[0].slot != 22 || b.projectiles[0].velocityFP != 340 {
		t.Fatal("原版158拍槽22的發射條件不符")
	}
	for n := 0; n < 28; n++ {
		b.stepProjectiles()
	}
	if len(b.projectiles) != 1 {
		t.Fatal("原版185拍仍在飛行")
	}
	p := b.projectiles[0]
	if p.x != 16 || p.y != 25 || p.z != 0 || p.heightFP != 204 || p.velocityFP != -220 || p.power != 92 {
		t.Fatalf("原版185拍飛行記錄不同：%+v", p)
	}
	s.ProjectileCooldown = 0
	b.shoot(1, 38, &e)
	if len(b.projectiles) != 1 || s.ProjectileCooldown != 0 {
		t.Fatal("落地前仍須等待固定槽釋放")
	}
	b.stepProjectiles()
	if len(b.projectiles) != 0 {
		t.Fatal("原版186拍落地須於同次更新釋放槽22")
	}
	b.shoot(1, 38, &e)
	if len(b.projectiles) != 1 || s.ProjectileCooldown != 8 {
		t.Fatal("原版187拍可由同一槽正常重發")
	}
}

func TestProjectileResolvesHitBeforeMovement(t *testing.T) {
	for player := 0; player < 2; player++ {
		b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
		b.PlayerSide = player
		for side := range b.Sides {
			x := 3 + side*50
			b.Sides[side].Soldiers[0] = Soldier{Alive: true, Kind: General, HP: 200, Stamina: 128, X: x, Y: 10, StepX: x, StepY: 10, Cmd: Duel, Next: Duel, Target: -1}
		}
		s := &b.Sides[0].Soldiers[1]
		*s = Soldier{Alive: true, Kind: Infantry, HP: 100, Stamina: 128, X: 20, Y: 20, StepX: 21, StepY: 20, Cmd: Holding, Next: Holding, Target: -1}
		b.Sides[1].Soldiers[1] = Soldier{Alive: true, Kind: Infantry, HP: 100, X: 40, Y: 40, StepX: 40, StepY: 40, Cmd: Holding, Next: Holding, Target: -1}
		b.projectiles = []projectile{{side: 1, x: 20, y: 20, z: 0, heightFP: 0, velocityFP: 0, power: 28, direction: West, gridIndex: projectileGridIndex(20, 20, 0)}}
		b.Step()
		if s.HP != 93 || s.X != 20 || s.Y != 20 || !s.Hurt || len(b.projectiles) != 0 {
			t.Fatalf("飛行須在士兵走離前命中並阻止當拍移動，player%d：%+v", player, *s)
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
