package tactical

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
)

func TestGroundScansHighPlaneWithoutLowGround(t *testing.T) {
	var layers [256][Levels]byte
	layers[0x20][4] = topFace
	layers[0xE0][5] = topFace
	layers[0xF8][6] = topFace
	layers[0xF9][2] = topFace
	layers[0xF9][4] = topFace
	for _, tt := range []struct {
		tile      byte
		low, high int
	}{
		{0x20, noGround, 4},
		{0xD0, noGround, 5}, // 城壁讀破壞後的0xE0。
		{0xF8, noGround, 6}, // 無低層時不設門標記。
		{0xF9, 2, navGate},
		{0x21, noGround, noGround},
	} {
		low, high := groundOf(tt.tile, &layers)
		if low != tt.low || high != tt.high {
			t.Fatalf("圖塊%#x的兩平面=%d/%d，應為%d/%d", tt.tile, low, high, tt.low, tt.high)
		}
	}
}

func TestOriginalPlayableRowsInclude62(t *testing.T) {
	if Height != 64 || MaxY != 62 || !inBounds(10, 62) || inBounds(10, 63) || clampY(63) != 62 {
		t.Fatalf("原版有效Y應為1–62，緩衝64列：Height=%d MaxY=%d", Height, MaxY)
	}
	rows := make([][]int, 64)
	for y := range rows {
		rows[y] = make([]int, 64)
	}
	rows[62][10] = 1
	f := NewField(rows, 0)
	if z, ok := f.GroundLevel(10, 62, PlaneLow); !ok || z != 1 {
		t.Fatalf("原版最後地形列的高度=%d/%v，應為1", z, ok)
	}
}

func TestLowSlopeKeepsPlaneAndRelativeWaypointHeight(t *testing.T) {
	rows := make([][]int, Height)
	for y := range rows {
		rows[y] = make([]int, Width)
	}
	rows[10][11] = 1
	b := &Battle{Field: NewField(rows, 0)}
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, Kind: Infantry, X: 10, Y: 10, Z: 0, StepZ: 3}
	if moved, _ := b.tryMove(0, 1, 11, 10, 0); !moved || s.Z != 1 || s.StepZ != 4 || s.PlaneHigh != 0 || s.OnWall || s.Climbing {
		t.Fatalf("低平面上坡應同步中繼高度並保持平面：%+v", *s)
	}
	if moved, _ := b.tryMove(0, 1, 10, 10, 0); !moved || s.Z != 0 || s.StepZ != 3 || s.PlaneHigh != 0 {
		t.Fatalf("低平面下坡應保留相對中繼高度：%+v", *s)
	}
}

// realLib 載入原版戰場資料；沒有素材就跳過。
func realLib(t *testing.T) *battle.Library {
	t.Helper()
	const dir = "../../../workplace/orig/dosv"
	if _, err := os.Stat(dir + "/BATTLE.MAP"); err != nil {
		t.Skip("找不到原版 BATTLE.MAP，跳過")
	}
	read := func(name string) []byte {
		b, err := os.ReadFile(dir + "/" + name)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	lib, err := battle.Parse(read("BATTLE.MAP"), read("BATTLE.MDL"), nil)
	if err != nil {
		t.Fatal(err)
	}
	return lib
}

func realField(t *testing.T, lib *battle.Library, n int) *Field {
	t.Helper()
	return NewFieldFromTileLayers(lib.MinimapRows(n), lib.Heights(n), lib.TileLayers(n), lib.GateX(n))
}

// 地面層要與機器碼演算法算出來的一樣（docs/re/63 §2）：
// 低平面幾乎全是層0；高平面只有保留方向的格共用正規化高度，異高格仍保留高度。
func TestGroundPlanesMatchRawAlgorithm(t *testing.T) {
	lib := realLib(t)
	for _, n := range []int{0, 5, 20} {
		f := realField(t, lib, n)
		lowHist := map[int]int{}
		highHist := map[int]int{}
		connectedHigh := map[int]int{}
		gates := 0
		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				if lv, ok := f.GroundLevel(x, y, PlaneLow); ok {
					lowHist[lv]++
				}
				if lv, ok := f.GroundLevel(x, y, PlaneHigh); ok {
					highHist[lv]++
					if f.nav[PlaneHigh][y][x]&0xF0 != 0 {
						connectedHigh[lv]++
					}
				}
				if f.IsGateCell(x, y) {
					gates++
				}
			}
		}
		// 存的是兵腳下佔的那一層，平地是 0。
		if lowHist[0] < 2000 {
			t.Errorf("戰場 %d：低平面平地只有 %d 格，太少", n, lowHist[0])
		}
		if len(connectedHigh) != 1 {
			t.Errorf("戰場%d：保留導航方向的高層應共用一種高度，得到%v", n, connectedHigh)
		}
		if gates == 0 {
			t.Errorf("戰場 %d 是攻城圖，卻一個門格都沒有", n)
		}
		t.Logf("戰場 %d：低平面 %v／高平面 %v／門格 %d", n, lowHist, highHist, gates)
	}
}

// 保留方向的高層只用4或5；1BC18保留的其他高度不得被當成無地面。
func TestHighPlaneDirectionsUseNormalizedLevel(t *testing.T) {
	lib := realLib(t)
	for _, n := range []int{0, 5, 20} {
		f := realField(t, lib, n)
		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				lv, ok := f.GroundLevel(x, y, PlaneHigh)
				if ok && f.nav[PlaneHigh][y][x]&0xF0 != 0 && lv != 4 && lv != 5 {
					t.Fatalf("戰場 %d (%d,%d) 高平面 ＝ %d，只能是 4 或 5", n, x, y, lv)
				}
			}
		}
	}
}

func TestRealFieldRetainsNonNavigableHeightSix(t *testing.T) {
	f := realField(t, realLib(t), 0)
	if lv, ok := f.GroundLevel(34, 4, PlaneHigh); !ok || lv != 6 || f.nav[PlaneHigh][4][34] != 0 {
		t.Fatal("戰場0的(34,4)高度6須保留，但沒有導航方向")
	}
}

// 上下城牆只在門格（原版圖塊 ≥ 0xF0）。城壁本體不行。
func TestOnlyGateCellsAllowClimb(t *testing.T) {
	lib := realLib(t)
	f := realField(t, lib, 0)
	walls, gates := 0, 0
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			tile := f.Tile(x, y)
			switch {
			case tile >= 0xF0:
				gates++
				if !f.IsGateCell(x, y) {
					t.Fatalf("(%d,%d) 圖塊 %#x 是門，卻不能上下", x, y, tile)
				}
			case tile >= 0xD0 && tile <= 0xDF:
				walls++
				if f.IsGateCell(x, y) {
					t.Fatalf("(%d,%d) 圖塊 %#x 是城壁，不該能上下", x, y, tile)
				}
			}
		}
	}
	if walls == 0 || gates == 0 {
		t.Fatalf("戰場 0 應該同時有城壁與門，量到 %d／%d", walls, gates)
	}
}

// 導航位元的兩條規則：低平面 |Δ| ≤ 1、高平面要相等。
func TestNavBitsFollowLevelRules(t *testing.T) {
	lib := realLib(t)
	f := realField(t, lib, 0)
	// 只檢查可站範圍內的格子：導航表建在整張 64×62 上，
	// 而 GroundLevel 照原版把座標夾在 1–62（inBounds）。
	for y := MinCoord + 1; y < MaxY; y++ {
		for x := MinCoord + 1; x < MaxCoord; x++ {
			for _, d := range navSteps {
				if !f.Linked(x, y, PlaneLow, d.dx, d.dy) {
					continue
				}
				a, okA := f.GroundLevel(x, y, PlaneLow)
				b, okB := f.GroundLevel(x+d.dx, y+d.dy, PlaneLow)
				if !okA || !okB {
					t.Fatalf("(%d,%d)→(%+d,%+d) 連通但有一邊沒有地面", x, y, d.dx, d.dy)
				}
				if abs(a-b) > 1 {
					t.Fatalf("(%d,%d)→(%+d,%+d) 連通但高度差 %d", x, y, d.dx, d.dy, abs(a-b))
				}
			}
		}
	}
}

// ⭐ 城壁擋人的是**通行層**不是地面層表。
//
// 地面層表是拿「打壞後的圖塊」算的（`sub_1BCA6` 把 0xD0–0xDF 換成 +0x10，
// docs/re/63 §2），所以尋路會直接規劃穿過城壁——原版也是這樣。
// 真正擋住的是 `sub_1BB6D` 展開的通行層，**打壞之後才會通**。
func TestSiegeWallBlocksUntilBroken(t *testing.T) {
	lib := realLib(t)
	for _, n := range []int{0, 5, 20} {
		f := realField(t, lib, n)
		gx := f.GateX()
		y := Height / 2
		// 尋路穿得過去（地面層表已經是打壞後的樣子）。
		if pts := f.FindPath(Point{X: gx - 8, Y: y}, Point{X: gx + 8, Y: y}, true, nil); len(pts) == 0 {
			t.Errorf("戰場 %d：尋路應該規劃得出穿過城門那一列的路", n)
		}
		// 但站在牆外那一格，往牆裡走會被通行層擋住。
		blocked := false
		for x := gx - 3; x <= gx+3; x++ {
			if f.Blocked(x, y, 1) {
				blocked = true
				break
			}
		}
		if !blocked {
			t.Errorf("戰場 %d：城門那一列在打壞前應該有一格擋著", n)
		}
	}
}

// ⭐ 打壞城門那一格就開得出通道。
//
// `docs/playtest/28` §2 曾經量到「城牆四格厚、可破壞的只有一格，
// 撞穿也開不出通道」——那是舊模型（「堆疊高度以下全實心」）的結論。
// 照原版的判準（子圖塊 ≥ 0x70 不擋人）重算之後，城門那一列**只有一格擋著**。
func TestBreakingWallOpensCorridor(t *testing.T) {
	lib := realLib(t)
	for _, n := range []int{0, 5, 20} {
		f := realField(t, lib, n)
		gx, y := f.GateX(), 32
		blocked := 0
		for x := gx - 4; x <= gx+4; x++ {
			if f.Blocked(x, y, 1) {
				blocked++
			}
		}
		if blocked != 1 {
			t.Errorf("戰場 %d：城門那一列有 %d 格擋著，預期 1", n, blocked)
		}
		for x := gx - 4; x <= gx+4; x++ {
			switch tl := f.Tile(x, y); {
			case tl >= 0xD0 && tl <= 0xDF:
				f.Retile(x, y, 0x10)
			case tl >= 0xF0 && tl < 0xF8:
				f.Retile(x, y, 8)
			}
		}
		for x := gx - 4; x <= gx+4; x++ {
			if f.Blocked(x, y, 1) {
				t.Errorf("戰場 %d：打壞之後 (%d,%d) 還擋著", n, x, y)
			}
		}
	}
}

func TestHighPlaneNormalizationKeepsHeightAndUsesByteGateVote(t *testing.T) {
	for _, tc := range []struct {
		name      string
		low, high int
		want      int8
	}{
		{"no gates", 0, 0, 4}, {"lower", 2, 1, 4}, {"tie", 1, 1, 4}, {"higher", 1, 2, 5},
		{"lower byte wraps", 256, 1, 5}, {"higher byte wraps", 1, 256, 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := flatField()
			for i := 0; i < tc.low+tc.high; i++ {
				x, y := i%64, i/64
				lv := int8(4)
				if i >= tc.low {
					lv = 5
				}
				f.lvl[PlaneHigh][y][x] = lv
				f.nav[PlaneHigh][y][x] = navGate | navLeft
			}
			f.lvl[PlaneHigh][60][60] = 4
			f.nav[PlaneHigh][60][60] = 0xF0
			f.lvl[PlaneHigh][60][61] = 5
			f.nav[PlaneHigh][60][61] = 0xF0
			f.lvl[PlaneHigh][60][62] = 5
			f.nav[PlaneHigh][60][62] = 0xF0
			f.normaliseHighPlane()
			for i := 0; i < tc.low+tc.high; i++ {
				x, y := i%64, i/64
				if f.lvl[PlaneHigh][y][x] != tc.want || f.nav[PlaneHigh][y][x] != navGate|navLeft {
					t.Fatal("門格只依byte票數選4／5，保留導航旗標")
				}
			}
			for x, lv := range map[int]int8{60: 4, 61: 5, 62: 5} {
				wantFlags := uint8(0)
				if lv == tc.want {
					wantFlags = 0xF0
				}
				if f.lvl[PlaneHigh][60][x] != lv || f.nav[PlaneHigh][60][x] != wantFlags {
					t.Fatal("異高只清旗標，不得把高度當成無地面")
				}
			}
		})
	}
}
func TestHighPlaneAxisRefusalDoesNotQueueBeforeAnotherAxisSucceeds(t *testing.T) {
	cases := []struct {
		name           string
		tile           byte
		level          int8
		head           bool
		cmd            Command
		allowed, queue bool
	}{
		{"height mismatch", 0x28, 5, false, Holding, false, false},
		{"matching height", 0xEF, 4, false, Holding, true, false},
		{"broken gate allowed", 0xF8, 4, false, Holding, true, false},
		{"missing ground still queues", 0x28, noLevel, false, Holding, false, true},
		{"head terrain still queues", 0x28, 4, true, Holding, false, true},
	}
	for tile := 0xF0; tile <= 0xF7; tile++ {
		cases = append(cases, struct {
			name           string
			tile           byte
			level          int8
			head           bool
			cmd            Command
			allowed, queue bool
		}{"intact gate", byte(tile), 4, false, Holding, false, false})
	}
	for tile := 0xF8; tile <= 0xFF; tile++ {
		cases = append(cases, struct {
			name           string
			tile           byte
			level          int8
			head           bool
			cmd            Command
			allowed, queue bool
		}{"Cmd6 broken gate", byte(tile), 4, false, wallAttack, false, false})
	}
	for player := 0; player < 2; player++ {
		for _, tc := range cases {
			b := newTestBattle(flatField())
			b.SetPlayerSide(player)
			holdRetreatCountdownUnits(b)
			f := b.Field
			f.tiles = make([][]byte, 64)
			for y := range f.tiles {
				f.tiles[y] = make([]byte, 64)
			}
			f.lvl[PlaneHigh][32][21] = tc.level
			f.lvl[PlaneHigh][33][20] = 4
			f.tiles[32][21] = tc.tile
			f.solid[5][32][21] = tc.head
			target := &b.Sides[player].Soldiers[1]
			*target = Soldier{Alive: true, HP: 100, Kind: Infantry, X: 20, Y: 32, Z: 4, PlaneHigh: 16, OnWall: true,
				StepX: 21, StepY: 33, StepZ: 4, GoalX: 21, GoalY: 33, GoalZ: 4, Cmd: tc.cmd, Next: tc.cmd, Target: -1, Stamina: 128}
			b.Step()
			wantX, wantY := 20, 33
			if tc.allowed {
				wantX, wantY = 21, 32
			}
			if target.X != wantX || target.Y != wantY || target.PathQueued != tc.queue {
				t.Fatalf("player%d %s tile%X：位置(%d,%d)、排隊%v；應(%d,%d)、%v", player, tc.name, tc.tile, target.X, target.Y, target.PathQueued, wantX, wantY, tc.queue)
			}
		}
	}
}
func TestHighPlaneAllAxesRefusedQueuesOnce(t *testing.T) {
	b := newTestBattle(flatField())
	holdRetreatCountdownUnits(b)
	b.Field.lvl[PlaneHigh][32][21] = 5
	b.Field.lvl[PlaneHigh][33][20] = 5
	s := &b.Sides[0].Soldiers[1]
	*s = Soldier{Alive: true, HP: 100, Kind: Infantry, X: 20, Y: 32, Z: 4, PlaneHigh: 16, OnWall: true, StepX: 21, StepY: 33, StepZ: 4, GoalX: 21, GoalY: 33, Cmd: Holding, Next: Holding, Target: -1, Stamina: 128}
	b.Step()
	if s.X != 20 || s.Y != 32 || !s.PathQueued || b.paths.n != 1 {
		t.Fatal("所有軸都拒絕時，才由共用出口排一次尋路")
	}
}
