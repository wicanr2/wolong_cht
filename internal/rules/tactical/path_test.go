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
	f := walledField(32)   // X=32 一道 4 層高的牆，只有中間那一格通
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
	f := flatField()
	// 成本函式一律回 0 ＝ 全部不可通行。
	if got := f.FindPath(Point{X: 5, Y: 5}, Point{X: 40, Y: 40}, false,
		func(int, int) int { return 0 }); got != nil {
		t.Errorf("全不可通行時回了 %v，應為 nil", got)
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
