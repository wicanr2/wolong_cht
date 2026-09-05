package tactical

import "testing"

// 開場擺位（`docs/spec/133`）：邊界那一欄 ＋ 亂數 Y。
//
// ⚠ **判準是值域與落點，不是逐兵全等**——Y 是亂數，兩邊的亂數不同源。
func spawnBattle(t *testing.T) *Battle {
	t.Helper()
	stack := make([][]int, Height)
	for y := range stack {
		stack[y] = make([]int, Width)
	}
	var r seqRand
	b := NewBattle(NewField(stack, 0), SyntheticFormations(), &r, 0)
	for side := 0; side < 2; side++ {
		for sq := 0; sq < Squads; sq++ {
			b.Deploy(side, sq, Infantry, 100)
		}
	}
	return b
}

// TestSpawnPutsEveryoneOnTheEdgeColumn 釘住 X 的兩個立即值。
func TestSpawnPutsEveryoneOnTheEdgeColumn(t *testing.T) {
	b := spawnBattle(t)
	var r seqRand
	b.Spawn(&r)
	want := [2]int{spawnXSide0, spawnXSide1}
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			s := &b.Sides[side].Soldiers[k]
			if !s.Alive {
				continue
			}
			if s.X != want[side] {
				t.Fatalf("側 %d 第 %d 個兵在 X=%d，應為 %d", side, k, s.X, want[side])
			}
		}
	}
}

// TestSpawnYStaysInRange 釘住 `rnd & 0x1F + 0x10` 的值域。
func TestSpawnYStaysInRange(t *testing.T) {
	b := spawnBattle(t)
	var r seqRand
	b.Spawn(&r)
	seen := map[int]bool{}
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			s := &b.Sides[side].Soldiers[k]
			if !s.Alive {
				continue
			}
			if s.Y < spawnYBase || s.Y > spawnYBase+spawnYMask {
				t.Fatalf("側 %d 第 %d 個兵的 Y=%d 落在 %d–%d 之外",
					side, k, s.Y, spawnYBase, spawnYBase+spawnYMask)
			}
			seen[s.Y] = true
		}
	}
	// 值域只有 32 格而場上有 96 個兵——**一定**會重複，
	// 所以「沒有重複」表示某處偷偷加了避讓。
	if len(seen) > spawnYMask+1 {
		t.Fatalf("用到了 %d 種 Y，值域只有 %d 種", len(seen), spawnYMask+1)
	}
}

// TestSpawnAllowsStacking 釘住「不查佔用」。
//
// ⚠ 加一道「那一格有人嗎」不是保險，是改行為——原版實測同一隊
// 三個兵在同一格（`docs/re/87` §3）。
func TestSpawnAllowsStacking(t *testing.T) {
	b := spawnBattle(t)
	var r seqRand
	b.Spawn(&r)
	cells := map[[3]int]int{}
	dup := 0
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			s := &b.Sides[side].Soldiers[k]
			if !s.Alive {
				continue
			}
			key := [3]int{side, s.X, s.Y}
			cells[key]++
			if cells[key] > 1 {
				dup++
			}
		}
	}
	if dup == 0 {
		t.Fatal("96 個兵擠 32 種 Y 卻沒有任何重疊——擺位偷偷避讓了")
	}
}

// seqRand 是遞增的假亂數：值域與重疊都驗得到，而且結果可重現。
type seqRand int

func (r *seqRand) Next() int { *r += 7; return int(*r) }

// TestSpawnFollowsTheFormationLine 釘住「擺哪一邊看陣形線，不看側號」。
//
// ⚠ 這是踩過的坑：remake 的 `Sides[0]` 恆為**攻方**，而原版的側 0 恆為
// **玩家**（玩家守城時整個戰場轉 180 度）。照側號擺會把兩邊各丟到對面
// 那一端——走 59 格、體力耗光，整場停住（`docs/spec/133` §3.5）。
func TestSpawnFollowsTheFormationLine(t *testing.T) {
	b := spawnBattle(t)
	// 把兩側的陣形線對調：側 0 在遠端、側 1 在近端。
	b.Sides[0].Line, b.Sides[1].Line = LineFor(1, 0), LineFor(0, 0)
	var r seqRand
	b.Spawn(&r)
	want := [2]int{spawnXSide1, spawnXSide0} // 跟著陣形線一起對調
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			s := &b.Sides[side].Soldiers[k]
			if !s.Alive {
				continue
			}
			if s.X != want[side] {
				t.Fatalf("陣形線在 %d 的那一側擺在 X=%d，應為 %d",
					b.Sides[side].Line, s.X, want[side])
			}
		}
	}
}
