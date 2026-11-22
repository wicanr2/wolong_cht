package tactical

import "testing"

func TestCommandAssistWalkingAndAscending(t *testing.T) {
	for _, tc := range []struct {
		name    string
		enabled bool
		rating  [2]int
		want    [2]int
	}{
		{"off", false, [2]int{15, 15}, [2]int{2, 2}},
		{"opponents", true, [2]int{15, 1}, [2]int{4, 2}},
		{"mirrored", true, [2]int{1, 15}, [2]int{2, 4}},
		{"midpoint", true, [2]int{8, 8}, [2]int{3, 3}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for side, want := range tc.want {
				for _, climb := range []bool{false, true} {
					f := flatField()
					if climb {
						f, _ = tiledField(32)
					}
					b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
					b.Sides[side].CommanderRating = tc.rating[side]
					b.SetCommandAssist(tc.enabled)
					s := &b.Sides[side].Soldiers[1]
					*s = Soldier{Alive: true, Kind: Infantry, X: 10, Y: 10, StepX: 20, StepY: 10, Stamina: 10, Target: -1}
					if climb {
						b.breakRow(Height / 2)
						*s = Soldier{Alive: true, Kind: Infantry, X: 32, Y: Height / 2, StepX: 32, StepY: Height / 2,
							StepZ: 4, Cmd: ScaleWal, Stamina: 10, Target: -1}
					}
					b.initUnitCollision()
					for range 2 {
						b.updateSoldierMovement(side, 1)
						b.drawUnitCollision(side, 1)
					}
					got := s.X - 10
					if climb {
						got = s.Z
					}
					if got != want || s.Stamina != 8 {
						t.Fatalf("side %d climb %t: moved %d (want %d), stamina %d (want 8)", side, climb, got, want, s.Stamina)
					}
				}
			}
		})
	}
}

func TestCommandAssistPathBudgetPreservesBaseFIFO(t *testing.T) {
	b := queueFixture(t)
	for k := range 3 {
		b.requestPath(0, k)
		b.requestPath(1, k)
	}
	b.Sides[1].CommanderRating = 15
	b.SetCommandAssist(true)
	b.drainPathQueue()
	want := []pathRequest{{0, 1}, {0, 2}, {1, 2}}
	for _, expected := range want {
		if got, ok := b.paths.pop(); !ok || got != expected {
			t.Fatalf("remaining queue: got %+v, want %+v", got, expected)
		}
	}
	if b.paths.n != 0 || b.Sides[1].Soldiers[1].PathQueued {
		t.Fatal("boost should service only one extra request from the high-command side")
	}
}

func TestCommandAssistReachesGateInTwoMovementFrames(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		f, _ := tiledField(32)
		b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
		b.breakRow(Height / 2)
		b.SetCommandAssist(enabled)
		b.Sides[AttackerSide].CommanderRating = 15
		s := &b.Sides[AttackerSide].Soldiers[1]
		*s = Soldier{Alive: true, Kind: Infantry, X: 28, Y: Height / 2, StepX: 32, StepY: Height / 2,
			StepZ: 4, Cmd: ScaleWal, Stamina: 10, Target: -1}
		b.initUnitCollision()
		for range 2 {
			b.updateSoldierMovement(AttackerSide, 1)
			b.drawUnitCollision(AttackerSide, 1)
		}
		want := 30
		if enabled {
			want = 32
		}
		if s.X != want {
			t.Fatalf("assist %t: soldier reached x=%d, want %d", enabled, s.X, want)
		}
	}
}

func TestCommandAssistBonusStepStopsBeforeEnemy(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{}, 0)
	b.Sides[0].CommanderRating = 15
	b.SetCommandAssist(true)
	a, e := &b.Sides[0].Soldiers[1], &b.Sides[1].Soldiers[1]
	*a = Soldier{Alive: true, Kind: Infantry, X: 10, Y: 10, HP: 100, StepX: 15, StepY: 10, Target: 1}
	*e = Soldier{Alive: true, Kind: Infantry, X: 12, Y: 10, HP: 100, Target: -1}
	b.initUnitCollision()
	b.updateSoldierMovement(0, 1)
	if a.X != 11 || e.HP != 100 || e.Stun != 0 {
		t.Fatalf("bonus should stop at occupied cell: attacker x=%d, enemy hp=%d stun=%d", a.X, e.HP, e.Stun)
	}
}

func TestCommandAssistLadderPreparation(t *testing.T) {
	for _, tc := range []struct {
		name     string
		rating   int
		enabled  bool
		contacts int
	}{
		{"off", 15, false, 81},
		{"minimum", 1, true, 81},
		{"midpoint", 8, true, 55},
		{"maximum", 15, true, 41},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, ladder := ladderPreparationFixture()
			s := &b.Sides[AttackerSide].Soldiers[1]
			b.Sides[AttackerSide].CommanderRating = tc.rating
			b.SetCommandAssist(tc.enabled)
			for step := 1; step < tc.contacts; step++ {
				if !b.tryClimb(AttackerSide, 1) || ladder.Broken || s.Z != 0 {
					t.Fatalf("premature ladder completion at contact %d", step)
				}
			}
			if ladder.Durability != 0 || b.Field.Tile(s.X, s.Y) != TileGateLo {
				t.Fatal("zero progress must retain the old tile until the next contact")
			}
			if !b.tryClimb(AttackerSide, 1) || !ladder.Broken || s.Z != 0 ||
				b.Field.Tile(s.X, s.Y) != TileGateLo+brokenGateDelta {
				t.Fatal("final contact must reveal the ladder without moving the soldier")
			}
			if !b.tryClimb(AttackerSide, 1) || s.Z != 1 {
				t.Fatal("the revealed ladder must permit climbing")
			}
		})
	}
}

func TestCommandAssistLadderEligibilityAndToggle(t *testing.T) {
	b, ladder := ladderPreparationFixture()
	b.Sides[AttackerSide].CommanderRating = 8
	b.Sides[DefenderSide].CommanderRating = 15
	b.SetCommandAssist(true)
	s := &b.Sides[AttackerSide].Soldiers[1]
	s.Kind = Cavalry
	if b.tryClimb(AttackerSide, 1) || ladder.Durability != GateDurability {
		t.Fatal("cavalry cannot prepare a ladder through climbing")
	}
	if !b.hitStructure(AttackerSide, East, ladder.X, ladder.Y) || ladder.Durability != GateDurability-1 ||
		ladder.ladderCredit != 0 {
		t.Fatal("ordinary contact must not receive ladder credit")
	}
	s.Kind = Infantry
	b.Sides[DefenderSide].Soldiers[1] = *s
	if !b.tryClimb(DefenderSide, 1) || ladder.Durability != GateDurability-1 {
		t.Fatal("defender contact must still leave the intact structure unchanged")
	}
	b.tryClimb(AttackerSide, 1)
	if ladder.Durability != GateDurability-2 || ladder.ladderCredit != 7 {
		t.Fatal("first eligible contact must accumulate half a bonus increment")
	}
	b.SetCommandAssist(false)
	b.SetCommandAssist(true)
	b.tryClimb(AttackerSide, 1)
	if ladder.Durability != GateDurability-3 || ladder.ladderCredit != 7 {
		t.Fatal("switching off must preserve progress while clearing fractional credit")
	}
}

func TestCommandAssistUsesEnemyAttackerRatingWhenPlayerDefends(t *testing.T) {
	b, ladder := ladderPreparationFixture()
	b.SetPlayerSide(DefenderSide)
	// 翻轉後的梯位用奇數圖格，攻方的上爬面向朝西。
	b.Field.Retile(ladder.X, ladder.Y, 1)
	b.Sides[AttackerSide].CommanderRating = 15
	b.Sides[DefenderSide].CommanderRating = 1
	b.SetCommandAssist(true)
	if !b.tryClimb(AttackerSide, 1) || ladder.Durability != GateDurability-2 {
		t.Fatal("computer-controlled attacker should receive its own command bonus")
	}
}

func ladderPreparationFixture() (*Battle, *Structure) {
	f, _ := tiledField(32)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
	ladder := &b.Structures[b.structureAt(32, Height/2)]
	b.Sides[AttackerSide].Soldiers[1] = Soldier{Alive: true, Kind: Infantry, X: ladder.X, Y: ladder.Y,
		StepX: ladder.X, StepY: ladder.Y, StepZ: 4, Cmd: ScaleWal, Target: -1}
	return b, ladder
}
