package tactical

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

func TestClimbHeadCellReadsOneLayerIncludingTopBody(t *testing.T) {
	g := &collisionGrid{}
	g.drawSlot(2, 20, 30, 3, PlaneHigh)
	if g.atLayer(20, 30, 2) != -1 || g.at(20, 30, 2) != 2 {
		t.Fatal("頭格查詢不可連帶查更高一格")
	}
	g.drawSlot(3, 21, 30, 6, PlaneHigh)
	if g.atLayer(21, 30, 7) != 3 || g.atLayer(21, 30, 8) != -1 || g.atLayer(21, 30, -1) != -1 {
		t.Fatal("第7層保留上半身，第8層及負層拒絕")
	}
	if g.atLayer(Width, 30, 1) != -1 {
		t.Fatal("格座標越界須安全返回")
	}
}

func climbHeadFixture(side, player int, stored bool, otherZ int, enemy bool) *Battle {
	f, _ := tiledField(32)
	f.Retile(32, Height/2, brokenGateDelta)
	b := NewBattle(f, SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
	b.PlayerSide = player
	b.Advantage[0] = Even
	b.Sides[side].Soldiers[1] = Soldier{Alive: true, Kind: Infantry, HP: 200, Power: 32,
		X: 32, Y: Height / 2, StepZ: 4, Cmd: ScaleWal, Next: ScaleWal, Target: -1}
	otherSide := side
	if enemy {
		otherSide = 1 - side
	}
	b.Sides[otherSide].Soldiers[2] = Soldier{Alive: true, Kind: Infantry, HP: 100, Power: 10,
		X: 32, Y: Height / 2, Z: otherZ, PlaneHigh: PlaneHighElevated, OnWall: true,
		StepZ: 4, Cmd: ScaleWal, Next: ScaleWal, Target: -1}
	if stored {
		b.initUnitCollision()
	}
	return b
}

func TestClimbHeadCellSkipsRemoteSoldierAndTracksCurrentZ(t *testing.T) {
	for side := 0; side < 2; side++ {
		for player := 0; player < 2; player++ {
			for _, stored := range []bool{false, true} {
				b := climbHeadFixture(side, player, stored, 3, false)
				s, other := &b.Sides[side].Soldiers[1], &b.Sides[side].Soldiers[2]
				if !b.tryClimb(side, 1) || s.Z != 1 || s.StepZ != 4 || s.PlaneHigh != PlaneHighElevated {
					t.Fatalf("Z2空、遠層Z3／4有兵仍可升到1，side%d player%d stored%v：%+v", side, player, stored, *s)
				}
				other.Swapped = true
				b.drawUnitCollision(side, 1)
				if b.tryClimb(side, 1) || s.Z != 1 {
					t.Fatal("第二次須查新頭格Z3，禁止換位的兵應阻擋")
				}
			}
		}
	}
}

func TestClimbHeadCellRecognizesTwoCellBody(t *testing.T) {
	for _, stored := range []bool{false, true} {
		for _, z := range []int{1, 2} {
			b := climbHeadFixture(1, 1, stored, z, false)
			s, other := &b.Sides[1].Soldiers[1], &b.Sides[1].Soldiers[2]
			other.Swapped = true
			if side, k := b.climbHeadOccupant(s); side != 1 || k != 2 {
				t.Fatalf("Z%d單兵的身體須覆蓋Z2頭格，stored%v：%d/%d", z, stored, side, k)
			}
			if b.tryClimb(1, 1) || s.Z != 0 {
				t.Fatal("真正頭格被拒換單兵佔用時不得升層")
			}
		}
	}
	b := climbHeadFixture(0, 0, true, 2, false)
	b.Sides[0].Soldiers[2].Alive = false
	if _, k := b.climbHeadOccupant(&b.Sides[0].Soldiers[1]); k != 2 {
		t.Fatal("正式查詢應讀持續格表，不依Alive重建")
	}
	b = climbHeadFixture(0, 0, false, 7, false)
	b.Sides[0].Soldiers[1].Z = 6
	if _, k := b.climbHeadOccupant(&b.Sides[0].Soldiers[1]); k != -1 {
		t.Fatal("無表fallback也須拒絕第8層查詢")
	}
}

func TestClimbHeadCellCollisionReturnsWithoutExtraAscent(t *testing.T) {
	for side := 0; side < 2; side++ {
		for player := 0; player < 2; player++ {
			for _, stored := range []bool{false, true} {
				for _, enemy := range []bool{false, true} {
					b := climbHeadFixture(side, player, stored, 2, enemy)
					s := &b.Sides[side].Soldiers[1]
					otherSide := side
					if enemy {
						otherSide = 1 - side
					}
					other := &b.Sides[otherSide].Soldiers[2]
					if !b.tryClimb(side, 1) {
						t.Fatal("上爬撞兵應回傳共用處理成功")
					}
					if enemy {
						if s.Z != 0 || !s.HitGeneral || other.HP != 100 {
							t.Fatal("固定亂數未命中仍應結束碰撞，不得升層")
						}
						b.rng = &fixedRand{seq: []int{64}}
						if !b.tryClimb(side, 1) || s.Z != 0 || other.HP >= 100 {
							t.Fatal("命中時沿共用傷害入口，仍不額外升層")
						}
					} else if s.Z != 2 || s.PlaneHigh != PlaneHighElevated || other.Z != 0 || other.PlaneHigh != PlaneHighGround || !other.Swapped {
						t.Fatal("友兵換位後直接返回，不得再由Z2升成3")
					}
				}
			}
		}
	}
}

func TestCollisionGridAgainstOriginalOpeningReceipt(t *testing.T) {
	const root = "../../../workplace/six-fixes-20260908/original-collision/"
	data, err := os.ReadFile(root + "tick-0001.json")
	if os.IsNotExist(err) {
		t.Skip("未安裝本機 dosgolem 收據；不宣稱原版對拍通過")
	}
	if err != nil {
		t.Fatal(err)
	}
	var snapshot struct {
		Units string `json:"unit_bytes"`
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	raw, err := hex.DecodeString(snapshot.Units)
	if err != nil || len(raw) < 96*32 {
		t.Fatal("原版單位資料長度或編碼錯誤")
	}
	want, err := os.ReadFile(root + "collision-tick2-0000.bin")
	if err != nil || len(want) < 8192 {
		t.Fatal("原版碰撞表缺漏")
	}
	var g collisionGrid
	for i := 0; i < 96; i++ {
		u := raw[i*32 : (i+1)*32]
		g.drawSlot(i, int(u[6]), int(u[8]), int(u[10]), int(u[30])>>4)
	}
	for i := 0; i < 8192; i++ {
		if g.cells[i] != want[i]&127 {
			t.Fatalf("ES:%04x low7=%d，原版=%d", i, g.cells[i], want[i]&127)
		}
	}
}

func TestCollisionGridOverlappingDrawAndClear(t *testing.T) {
	var g collisionGrid
	g.drawSlot(3, 1, 38, 0, PlaneLow)
	g.drawSlot(8, 1, 38, 0, PlaneLow)
	if got := g.at(1, 38, 0); got != 8 {
		t.Fatalf("重疊重畫取槽 %d，want 8", got)
	}
	g.drawSlot(3, 2, 38, 0, PlaneLow)
	if got := g.at(1, 38, 0); got != -1 {
		t.Fatalf("原格應清空，got %d", got)
	}
	if got := g.at(2, 38, 0); got != 3 {
		t.Fatalf("新格 = %d，want 3", got)
	}
}

func TestCollisionGridSwapPreservesOriginalUpperLayerAsymmetry(t *testing.T) {
	var g collisionGrid
	g.drawSlot(8, 1, 38, 0, PlaneLow)
	g.drawSlot(20, 1, 37, 0, PlaneLow)
	a, _ := collisionCell(1, 38, 0)
	b, _ := collisionCell(1, 37, 0)
	g.swapSlots(8, 20, a, b)
	if g.cells[a] != 9 || g.cells[b] != 21 || g.cells[a+0x1000] != 9 || g.cells[b+0x1000] != 9 {
		t.Fatalf("換格不符原始寫入：%d %d / %d %d", g.cells[a], g.cells[b], g.cells[a+0x1000], g.cells[b+0x1000])
	}
	if g.old[8] != b || g.old[20] != a {
		t.Fatal("未交換原格指標")
	}
}
