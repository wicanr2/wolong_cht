package tactical

import "testing"

func TestClimbStructureHandledBlocksFollowingSwap(t *testing.T) {
	for player := 0; player < 2; player++ {
		for side := 0; side < 2; side++ {
			for _, kind := range []Kind{Archer, Infantry} {
				f, _ := tiledField(32)
				b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
				b.PlayerSide = player
				gate := -1
				for i := range b.Structures {
					if b.Structures[i].Kind == KindGate && !b.Structures[i].Broken {
						gate = i
						break
					}
				}
				if gate < 0 {
					t.Fatal("合成場必須包含未破的門")
				}
				g := &b.Structures[gate]
				s := &b.Sides[side].Soldiers[1]
				*s = Soldier{Alive: true, Kind: kind, HP: 200, Stamina: 10, X: g.X, Y: g.Y,
					StepX: g.X, StepY: g.Y, StepZ: 4, Cmd: ScaleWal, Next: ScaleWal, Target: -1}
				before := g.Durability
				b.moveToward(side, 1)
				if !s.MoveFlag || !s.PathQueued || s.X != g.X || s.Y != g.Y || s.Z != 0 || s.Stamina != 9 {
					t.Fatalf("碰撞已處理須保留移動旗標及位置並排尋路，player%d side%d kind%d：%+v", player, side, kind, *s)
				}
				want := before
				if side == AttackerSide {
					want--
					if s.Facing == b.awayFromCastle() {
						want = 0
					}
				}
				if g.Durability != want {
					t.Fatalf("門耐久=%d，應為%d", g.Durability, want)
				}
				b.Sides[side].Soldiers[2] = Soldier{Alive: true, Kind: Infantry, X: g.X - 1, Y: g.Y, HP: 200}
				if b.swapWith(side, 2, side, 1) {
					t.Fatal("後方步兵不得換走已處理碰撞的前排")
				}
				for _, blockedKind := range []Kind{General, Cavalry} {
					s.Kind = blockedKind
					if b.tryClimb(side, 1) {
						t.Fatal("大將與騎兵仍不得純Z攀爬")
					}
				}
				s.Kind = kind
				b.Field, _ = tiledField(32)
				b.Structures = nil
				if b.tryClimb(side, 1) {
					t.Fatal("無結構物記錄時不得偽造處理成功")
				}
			}
		}
	}
}

func TestSlopeWaypointByteBoundary(t *testing.T) {
	for player := 0; player < 2; player++ {
		for side := 0; side < 2; side++ {
			for _, tc := range []struct {
				up         bool
				step, want int
			}{{false, 0, 255}, {false, 1, 0}, {true, 255, 0}, {true, 254, 255}} {
				rows := make([][]int, Height)
				for y := range rows {
					rows[y] = make([]int, Width)
				}
				rows[10][11] = 1
				b := NewBattle(NewField(rows, 0), SyntheticFormations(), &fixedRand{}, 0)
				b.PlayerSide = player
				x, z, goal := 11, 1, 10
				if tc.up {
					x, z, goal = 10, 0, 11
				}
				s := &b.Sides[side].Soldiers[1]
				*s = Soldier{Alive: true, Kind: Cavalry, X: x, Y: 10, Z: z, StepX: goal, StepY: 10, StepZ: tc.step, Stamina: 10, Target: -1}
				b.moveToward(side, 1)
				if s.X != goal || s.StepZ != tc.want || s.Stamina != 9 || !s.MoveFlag {
					t.Fatalf("水平坡地的中繼高度須按byte溢位，player%d side%d case%+v：%+v", player, side, tc, *s)
				}
			}
		}
	}
}

func TestClimbGateHeadUnitPrecedesStructure(t *testing.T) {
	for player := 0; player < 2; player++ {
		for side := 0; side < 2; side++ {
			for _, blocked := range []bool{false, true} {
				f, _ := tiledField(32)
				b := NewBattle(f, SyntheticFormations(), &fixedRand{}, 0)
				b.PlayerSide = player
				gate := b.structureAt(32, Height/2)
				if gate < 0 {
					t.Fatal("缺少合成門")
				}
				g := &b.Structures[gate]
				before := g.Durability
				s := &b.Sides[side].Soldiers[1]
				*s = Soldier{Alive: true, Kind: Infantry, HP: 200, X: 32, Y: Height / 2, StepZ: 4, Target: -1}
				other := &b.Sides[side].Soldiers[2]
				*other = Soldier{Alive: true, Kind: Infantry, HP: 200, X: 32, Y: Height / 2, Z: 2, PlaneHigh: PlaneHighElevated, OnWall: true, Swapped: blocked}
				b.initUnitCollision()
				handled := b.tryClimb(side, 1)
				if handled == blocked || s.PathQueued || g.Durability != before {
					t.Fatalf("未破門上的單兵優先，不得扣門耐久或排尋路，player%d side%d blocked%v", player, side, blocked)
				}
				if blocked {
					if s.Z != 0 || other.Z != 2 {
						t.Fatal("拒換頭格須保留雙方位置")
					}
				} else if s.Z != 2 || other.Z != 0 || !other.Swapped {
					t.Fatal("頭格換位成功後直接返回，不額外升層")
				}
			}
		}
	}
}
