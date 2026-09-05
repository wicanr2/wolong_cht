package tactical

import "testing"

// 平地上直線走：一路同方向，**中間不該有轉彎點**，只留終點。
func TestPathStraightLineHasNoTurns(t *testing.T) {
	f := flatField()
	got := f.FindPath(Point{X: 5, Y: 30}, Point{X: 20, Y: 30}, false, nil)
	if len(got) != 1 {
		t.Fatalf("直線走出 %d 個點 %v，應該只有終點一個"+
			"（原版只在轉彎時 stosw）", len(got), got)
	}
	if got[0] != (Point{X: 20, Y: 30}) {
		t.Errorf("終點是 %v，應為 (20, 30)", got[0])
	}
}

// 要轉一次彎的路徑，點數要 ≥ 2（轉角 ＋ 終點）。
func TestPathRecordsTurns(t *testing.T) {
	f := flatField()
	got := f.FindPath(Point{X: 5, Y: 30}, Point{X: 20, Y: 40}, false, nil)
	if len(got) < 2 {
		t.Fatalf("需要轉彎的路徑只有 %d 個點 %v", len(got), got)
	}
	if got[len(got)-1] != (Point{X: 20, Y: 40}) {
		t.Errorf("最後一個點是 %v，應為終點 (20, 40)", got[len(got)-1])
	}
	// 每一個點都要在戰場內，而且相鄰兩點必須同一軸（轉彎點的定義）。
	prev := Point{X: 5, Y: 30}
	for _, p := range got {
		if !inBounds(p.X, p.Y) {
			t.Errorf("點 %v 超出戰場", p)
		}
		if p.X != prev.X && p.Y != prev.Y {
			t.Errorf("%v → %v 不在同一軸上——轉彎點之間應該是直線", prev, p)
		}
		prev = p
	}
}

// ⭐ 爬不上去的兵要繞過城牆，爬得上去的可以直接翻。
func TestPathRespectsClimb(t *testing.T) {
	f := walledField(32) // X=32 一道 4 層高的牆，只有中間那一格通
	gate := Height / 2

	// 目標挑在牆的另一側、且不是門那一列。
	from := Point{X: 20, Y: gate - 10}
	to := Point{X: 40, Y: gate - 10}

	foot := f.FindPath(from, to, false, nil)
	if foot == nil {
		t.Fatal("爬不上去的兵找不到路——應該要繞到門那一格")
	}
	// 繞路一定會經過門那一列。
	viaGate := false
	for _, p := range foot {
		if p.Y == gate {
			viaGate = true
		}
	}
	if !viaGate {
		t.Errorf("步兵的路徑 %v 沒有經過門所在的那一列 %d", foot, gate)
	}

	// 爬得上去的兵可以直接翻牆，路徑不必繞。
	climb := f.FindPath(from, to, true, nil)
	if climb == nil {
		t.Fatal("爬得上去的兵反而找不到路")
	}
	if len(climb) > len(foot) {
		t.Errorf("爬牆的路徑（%d 點）比繞路的（%d 點）還長", len(climb), len(foot))
	}
}

// 走不到的目標要回 nil，不能回半條路。
func TestPathUnreachable(t *testing.T) {
	// 一道只有一層高的坎：**所有兵種都跨得過**。
	//
	// ⚠ 高度差 ≤ 1 的連通與兵種無關（`sub_1BD07`，docs/re/63 §3）；
	// 兵種只擋**純 Z 移動**（`sub_1AF69` 的 `cmp [si+4], 12h`），
	// 而純 Z 只在門那一格發生。四層的牆連爬得上去的兵也翻不過。
	stack := make([][]int, Height)
	for y := range stack {
		stack[y] = make([]int, Width)
		stack[y][32] = 1
	}
	f := NewField(stack, 32)
	for _, climb := range []bool{false, true} {
		if got := f.FindPath(Point{X: 20, Y: 30}, Point{X: 40, Y: 30}, climb, nil); got == nil {
			t.Errorf("一層高的坎應該跨得過（climb=%v）", climb)
		}
	}

	// 四層的牆連爬得上去的兵也過不去（一次只能上下一層）。
	for y := range stack {
		stack[y][32] = 4
	}
	tall := NewField(stack, 32)
	if got := tall.FindPath(Point{X: 20, Y: 30}, Point{X: 40, Y: 30}, true, nil); got != nil {
		t.Errorf("四層的牆不該翻得過去，卻回了 %v", got)
	}
}

// ⭐ 成本表是「有沒有兵站著」（有兵 ＝ 8），不是地形——docs/re/63 §1。
//
// 這條測試釘住「額外成本 0 ＝ 每走一格加 1 ＝ 波數」這個等價關係：
// 給任何一格加懲罰都只會讓路徑變長或不變，不會變短。
func TestPathPenaltyIsInertByDefault(t *testing.T) {
	f := flatField()
	base := f.FindPath(Point{X: 5, Y: 30}, Point{X: 30, Y: 40}, false, nil)
	zero := f.FindPath(Point{X: 5, Y: 30}, Point{X: 30, Y: 40}, false,
		func(int, int) int { return 0 })
	if len(base) != len(zero) {
		t.Errorf("預設（%d 點）與明寫 0（%d 點）不一致", len(base), len(zero))
	}
	// 對整條直線加重懲罰，路徑不該變短。
	heavy := f.FindPath(Point{X: 5, Y: 30}, Point{X: 30, Y: 40}, false,
		func(x, y int) int {
			if y == 30 {
				return 50
			}
			return 0
		})
	if heavy == nil {
		t.Fatal("加了懲罰就找不到路")
	}
	if len(heavy) < len(base) {
		t.Errorf("加懲罰之後轉彎點反而變少（%d < %d）", len(heavy), len(base))
	}
}

// 轉彎點最多 64 個——原版的緩衝區只有 128 byte（`mov cl, 40h`）。
func TestPathCapsAtBufferSize(t *testing.T) {
	if MaxWaypoints != 0x40 {
		t.Errorf("上限是 %d，原版是 0x40 ＝ 64", MaxWaypoints)
	}
	if MaxWaypoints*2 != 128 {
		t.Errorf("64 個點 × 2 byte 應為 128，與 0x1800 + 兵編號 × 128 對不上")
	}
}

// 原版只有抵達目前中繼點後才消費下一個點；若每幀直接取下一點，
// 兵會在尚未走完第一段時跳過轉角。
func TestWaypointsAdvanceOnlyAfterArrival(t *testing.T) {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{127}}, 0)
	s := &b.Sides[0].Soldiers[0]
	*s = Soldier{
		Alive: true, Kind: Infantry, HP: MaxHP, Power: DefaultPower,
		X: 10, Y: 20, Z: 0, GoalX: 20, GoalY: 20, GoalZ: 0,
		StepX: 10, StepY: 20, StepZ: 0, Cmd: Attack, Next: Attack,
		Path: &Waypoints{pts: []Point{{X: 12, Y: 20}, {X: 12, Y: 22}}},
	}

	b.moveToward(0, 0)
	if s.X != 11 || s.Path.Len() != 2 {
		t.Fatalf("第一幀走到 (%d,%d)，剩 %d 點；應為 (11,20) 與 2 點",
			s.X, s.Y, s.Path.Len())
	}
	b.moveToward(0, 0)
	if s.X != 12 || s.Y != 20 || s.Path.Len() != 2 {
		t.Fatalf("抵達第一中繼點後座標 (%d,%d)、剩 %d 點錯誤",
			s.X, s.Y, s.Path.Len())
	}
	b.moveToward(0, 0)
	if s.X != 12 || s.Y != 21 || s.Path.Len() != 1 {
		t.Fatalf("下一幀沒有前進到第二段：座標 (%d,%d)、剩 %d 點",
			s.X, s.Y, s.Path.Len())
	}
}

// ⭐ Y 的上限要跟實際的列數一致。
//
// `sub_1AACF` 把座標夾在 1–62，但一張戰場只有 62 列（0–61）——
// 拿 62 去索引會越界。**這個 bug 一直在，是尋路開始探邊界才炸出來的**：
// 在那之前沒有任何路徑會走到最後一列。
func TestFieldBoundsMatchRowCount(t *testing.T) {
	if MaxY != Height-1 {
		t.Fatalf("MaxY ＝ %d，應為 %d", MaxY, Height-1)
	}
	if inBounds(1, Height) {
		t.Error("Y ＝ Height 不該算在範圍內")
	}
	f := flatField()
	// 掃過整個合法範圍不能炸。
	for y := MinCoord; y <= MaxY; y++ {
		for x := MinCoord; x <= MaxCoord; x++ {
			f.StandLevel(x, y)
			f.Walkable(x, y, 0)
		}
	}
	if got := clampY(999); got != MaxY {
		t.Errorf("clampY(999) ＝ %d，應為 %d", got, MaxY)
	}
	// 沿著最後一列走得完，不會 panic。
	if p := f.FindPath(Point{X: 2, Y: MaxY}, Point{X: 30, Y: MaxY}, false, nil); p == nil {
		t.Error("沿著最後一列找不到路")
	}
}

// ⭐ 有兵擋在直線上就繞過去（`docs/spec/134`）。
//
// 判準是「**路徑不經過那一格**」，不是「點數變多」——繞路的轉彎點
// 可能只有兩個，而點數相同時舊的錯誤實作照樣通過。
func TestFindPathGoesAroundAnOccupiedCell(t *testing.T) {
	f := flatField()
	from, to := Point{X: 30, Y: 30}, Point{X: 30, Y: 33}
	block := Point{X: 30, Y: 31}
	got := f.FindPath(from, to, false, func(x, y int) int {
		if x == block.X && y == block.Y {
			return 8 // `sub_1B240` 對有兵的格子寫 8
		}
		return 0
	})
	if len(got) == 0 {
		t.Fatal("繞得過去卻回了空路徑")
	}
	for _, p := range walk(from, got) {
		if p == block {
			t.Fatalf("路徑 %v 穿過有兵的 %v——繞路成本沒有生效", got, block)
		}
	}
	if got[len(got)-1] != to {
		t.Errorf("最後一個點是 %v，應為終點 %v", got[len(got)-1], to)
	}
}

// 反向對照：沒有東西擋路時**不准**平白繞路，否則上面那條測試
// 用「一律繞路」的實作也會通過。
func TestFindPathStaysStraightWhenNothingBlocks(t *testing.T) {
	f := flatField()
	from, to := Point{X: 30, Y: 30}, Point{X: 30, Y: 33}
	got := f.FindPath(from, to, false, func(int, int) int { return 0 })
	if len(got) != 1 || got[0] != to {
		t.Fatalf("沒有人擋路卻走出 %v，應該只有終點一個點", got)
	}
}

// walk 把轉彎點清單攤成逐格的路線。
func walk(from Point, pts []Point) []Point {
	out := []Point{from}
	cur := from
	for _, p := range pts {
		for cur != p {
			switch {
			case cur.X < p.X:
				cur.X++
			case cur.X > p.X:
				cur.X--
			case cur.Y < p.Y:
				cur.Y++
			case cur.Y > p.Y:
				cur.Y--
			}
			out = append(out, cur)
		}
	}
	return out
}
