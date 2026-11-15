package tactical

import "testing"

// DOS/V 1B824、1BB6D；原版隔離樣本與完整重播見spec159。
func TestBrokenStructureClearsSixPhysicalLayers(t *testing.T) {
	lib := realLib(t)
	for _, player := range []int{0, 1} {
		for _, sortie := range []bool{false, true} {
			b := NewBattle(realField(t, lib, 76), SyntheticFormations(), nil, 0)
			b.SetPlayerSide(player)
			b.Structures = []Structure{{Kind: KindWall, X: 42, Y: 30, Run: 5, Durability: 3000}, {Kind: KindGate, X: 39, Y: 30, Run: 1, Durability: 80}}
			b.unitCollision = &collisionGrid{}
			for _, s := range b.Structures {
				for y := s.Y; y < s.Y+s.Run; y++ {
					tile := byte(0xd3)
					if s.Kind == KindGate {
						tile = 0xf0
					}
					b.Field.tiles[y][s.X] = tile
					for z := 0; z < 8; z++ {
						b.unitCollision.cells[z*Width*Height+y*Width+s.X] = 37
						if z < Levels {
							b.Field.solid[z][y][s.X] = true
						}
					}
					b.unitCollision.cost[y*Width+s.X] = 77
					b.unitCollision.cost[Width*Height+y*Width+s.X] = 88
				}
			}
			b.unitCollision.cells[29*Width+42] = 43
			nav, lvl := b.Field.nav, b.Field.lvl
			units := [2][SoldiersOnFoot]Soldier{b.Sides[0].Soldiers, b.Sides[1].Soldiers}
			if sortie {
				b.SortieBreaksWalls()
			} else {
				b.Structures[0].Durability = 0
				b.breakRow(30)
			}
			for _, s := range b.Structures {
				broken := !sortie || s.Kind == KindWall
				if s.Broken != broken {
					t.Fatalf("結構物狀態錯誤：%+v", s)
				}
				for y := s.Y; y < s.Y+s.Run; y++ {
					for z := 0; z < 8; z++ {
						want := byte(37)
						if broken && z < 6 {
							want = 0
						}
						if got := b.unitCollision.cells[z*Width*Height+y*Width+s.X]; got != want {
							t.Fatalf("側%d出擊%v格%d,%d,%d佔位%d，應為%d", player, sortie, s.X, y, z, got, want)
						}
						if broken && z < 6 && b.Field.Blocked(s.X, y, z) {
							t.Fatalf("破城後第%d層仍擋投射物", z)
						}
						if broken && z == 6 && b.Field.Blocked(s.X, y, z) != (s.Kind == KindWall) {
							t.Fatal("第6層應保留原版樣本E3實心／F8非實心的高位")
						}
					}
					want := byte(77)
					if broken {
						want = 0
					}
					if b.unitCollision.cost[y*Width+s.X] != want || b.unitCollision.cost[Width*Height+y*Width+s.X] != 88 {
						t.Fatal("只應清低平面成本")
					}
				}
			}
			if b.unitCollision.cells[29*Width+42] != 43 || b.Field.nav != nav || b.Field.lvl != lvl {
				t.Fatal("不可改鄰格或重建導航")
			}
			if units != [2][SoldiersOnFoot]Soldier{b.Sides[0].Soldiers, b.Sides[1].Soldiers} {
				t.Fatal("清格不可刪除兵士記錄")
			}
		}
	}
}

func TestLaterRetilePreservesEarlierBreach(t *testing.T) {
	b := NewBattle(realField(t, realLib(t), 76), SyntheticFormations(), nil, 0)
	b.Structures = []Structure{{Kind: KindWall, X: 42, Y: 30, Run: 1}, {Kind: KindWall, X: 43, Y: 40, Run: 3}}
	b.Field.tiles[30][42] = 0xd3
	for y := 40; y < 43; y++ {
		b.Field.tiles[y][43] = 0xd3
	}
	nav := b.Field.nav
	b.breakRow(30)
	b.breakRow(40)
	if b.Field.Tile(42, 30) != 0xe3 || b.Field.Blocked(42, 30, 3) || b.Field.nav != nav {
		t.Fatal("後續破城不可重新填回先前清空的E3實體")
	}
}

func TestBreakingStructuresPreservesUntouchedDurability(t *testing.T) {
	for _, sortie := range []bool{false, true} {
		f, _ := tiledField(32)
		b := NewBattle(f, SyntheticFormations(), nil, 250)
		original := append([]Structure(nil), b.Structures...)
		damage := b.CityDamage(250)
		if sortie {
			b.SortieBreaksWalls()
			b.SortieBreaksWalls()
		} else {
			b.Structures = append(b.Structures, Structure{Kind: KindGate, X: 33, Y: b.Structures[0].Y, Run: 1, Durability: 80})
			b.Structures[0].Durability = 0
			b.breakRow(b.Structures[0].Y)
			if s := b.Structures[len(b.Structures)-1]; !s.Broken || s.Durability != 80 {
				t.Fatal("同排未受擊門應保留耐久80")
			}
		}
		if sortie {
			for i, s := range b.Structures {
				if s.Durability != original[i].Durability {
					t.Fatal("出擊只開牆，不能歸零耐久")
				}
			}
			if b.CityDamage(250) != damage {
				t.Fatal("出擊不可憑空增加戰後城損")
			}
		}
	}
}

func TestAbsentGeneralRetreatAndSiegeDrain(t *testing.T) {
	for _, player := range []int{0, 1} {
		for _, alive := range []bool{false, true} {
			for _, hp := range []int{0, 1, 49, 50, 70} {
				for _, siege := range []bool{false, true} {
					for _, due := range []bool{false, true} {
						f := flatField()
						if siege {
							f.gateX = 32
						}
						b := NewBattle(f, SyntheticFormations(), nil, 0)
						b.SetPlayerSide(player)
						for side := range b.Sides {
							for k := 0; k < SoldiersOnFoot; k += PerSquad {
								b.Sides[side].Soldiers[k] = Soldier{HP: 200, Cmd: Holding, Next: Holding}
							}
						}
						g := &b.Sides[0].Soldiers[0]
						g.HP, g.Alive = hp, alive
						b.siegeTick = SiegeDrainInterval - 2
						if due {
							b.siegeTick++
						}
						b.checkGeneralRetreat()
						b.drainSiegeGeneral()
						phase := 0
						if hp < 50 {
							phase = 1
						}
						want := hp
						if siege && due && hp > 0 {
							want--
						}
						if g.HP != want || g.Alive != alive || b.endPhase != phase {
							t.Fatalf("玩家%d在場%v HP%d攻城%v到期%v：HP%d phase%d", player, alive, hp, siege, due, g.HP, b.endPhase)
						}
					}
				}
			}
		}
	}
}

func TestEmptyEnemyTargetRetainsOriginalDefault(t *testing.T) {
	for _, side := range []int{0, 1} {
		for _, old := range []int{-1, 0, 7} {
			for _, refresh := range []bool{false, true} {
				b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
				path := &Waypoints{pts: []Point{{X: 1, Y: 2}, {X: 3, Y: 4}}, i: 1}
				s := &b.Sides[side].Soldiers[1]
				*s = Soldier{Alive: true, Kind: Cavalry, X: 20, Y: 20, Target: old, TargetRefresh: refresh, Path: path}
				b.lockOnNearest(side, 1)
				if s.Target != 0 || s.TargetRefresh != (refresh || old != 0) || s.Path != path || path.i != 1 {
					t.Fatalf("空敵陣鎖敵應保留槽0、刷新及繞路：%+v", *s)
				}
			}
		}
	}
}
