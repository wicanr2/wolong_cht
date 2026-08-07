package army

import "testing"

func full(t TroopType) Corps {
	var c Corps
	c.Alive = true
	c.Morale = DefaultMorale
	for i := range c.Units {
		c.Units[i] = t
		c.Manned[i] = true
	}
	return c
}

// 滿編 6 個位置 × 1000 人 = 6000（說明書 5.5，編成畫面的許褚正是滿編）。
func TestStrength(t *testing.T) {
	if got := full(Infantry).Strength(); got != MaxStrength {
		t.Errorf("滿編 = %d, want %d", got, MaxStrength)
	}
	var half Corps
	for i := 0; i < 3; i++ {
		half.Manned[i] = true
	}
	if got := half.Strength(); got != 3000 {
		t.Errorf("三個位置 = %d, want 3000", got)
	}
	if got := (Corps{}).Strength(); got != 0 {
		t.Errorf("空編 = %d, want 0", got)
	}
}

// 六個位置的名稱（PDF p.13 以 400 dpi 重掃判讀）。
func TestPositionNames(t *testing.T) {
	want := []string{"大將", "先鋒", "左翼", "右翼", "左備", "右備"}
	for i, w := range want {
		if got := Position(i).String(); got != w {
			t.Errorf("位置 %d = %s, want %s", i, got, w)
		}
	}
}

// 節點分三類：據點 0–191、路上節點 192–255、野外 ≥256。
func TestNodeKind(t *testing.T) {
	cases := []struct {
		node int
		want NodeKind
	}{
		{0, CityNode}, {191, CityNode},
		{192, RoadNode}, {255, RoadNode},
		{256, FieldNode}, {9999, FieldNode},
	}
	for _, c := range cases {
		if got := KindOf(c.node); got != c.want {
			t.Errorf("節點 %d → %v, want %v", c.node, got, c.want)
		}
	}
}

// 到達判定比的是節點編號，不是座標 —— 原版就是 `cmp bx, [si+14h]`。
func TestArrived(t *testing.T) {
	c := Corps{Node: 10, TargetNode: 10, X: 5, Y: 5, TargetX: 99, TargetY: 99}
	if !c.Arrived() {
		t.Error("節點相同就算到達，座標不該影響")
	}
	c.TargetNode = 11
	if c.Arrived() {
		t.Error("節點不同不該算到達")
	}
}

// 純騎馬編成爬不上城牆（說明書 5.5）。
func TestCanScaleWalls(t *testing.T) {
	if full(Cavalry).CanScaleWalls() {
		t.Error("純騎馬不該爬得上城牆")
	}
	if !full(Infantry).CanScaleWalls() {
		t.Error("步兵應該爬得上城牆")
	}
	mixed := full(Cavalry)
	mixed.Units[RightGuard] = Archer
	if !mixed.CanScaleWalls() {
		t.Error("混編有弓兵就該爬得上去")
	}
	// 空編不算「純騎馬」，也爬不上去。
	if (Corps{}).AllCavalry() {
		t.Error("空編不該算純騎馬")
	}
}

// 純騎馬有移動速度加成（說明書 5.5）。
func TestAllCavalry(t *testing.T) {
	if !full(Cavalry).AllCavalry() {
		t.Error("六個位置都是騎馬應該算純騎馬")
	}
	mixed := full(Cavalry)
	mixed.Units[0] = Infantry
	if mixed.AllCavalry() {
		t.Error("混了步兵就不算純騎馬")
	}
}

// 士氣 < 100 且戰敗 → 軍團壞滅。100 是說明書直接給的門檻。
func TestDestroyed(t *testing.T) {
	c := Corps{Morale: 99}
	if !c.Destroyed(true) {
		t.Error("士氣 99 戰敗應該壞滅")
	}
	if c.Destroyed(false) {
		t.Error("沒戰敗不該壞滅")
	}
	c.Morale = RoutMoraleGate
	if c.Destroyed(true) {
		t.Error("士氣剛好 100 不該壞滅（門檻是「切った」＝低於）")
	}
}

// 佔用圖是**計數**不是布林 —— 同一格可以疊多個軍團。
func TestOccupancyIsACount(t *testing.T) {
	m := NewOccupancyMap()
	if m.Occupied(10, 10) {
		t.Error("新圖不該有東西")
	}
	m.Enter(10, 10)
	m.Enter(10, 10)
	if m.At(10, 10) != 2 {
		t.Errorf("兩個軍團同格 = %d, want 2", m.At(10, 10))
	}
	m.Leave(10, 10)
	if !m.Occupied(10, 10) {
		t.Error("走了一個還有一個，仍該算佔用")
	}
	m.Leave(10, 10)
	if m.Occupied(10, 10) {
		t.Error("兩個都走了就不該佔用")
	}
}

// ⚠ 原版的 dec 是無條件的，0 會繞回 255。remake 刻意夾在 0 ——
// 繞回 255 會讓那一格永遠看起來有兵，比直接夾住難查得多。
func TestOccupancyDoesNotUnderflow(t *testing.T) {
	m := NewOccupancyMap()
	for i := 0; i < 5; i++ {
		m.Leave(3, 3)
	}
	if m.At(3, 3) != 0 {
		t.Errorf("多減之後 = %d, want 0（不能繞回 255）", m.At(3, 3))
	}
}

// 越界的座標要安靜忽略，不能 panic 也不能寫到別格。
func TestOccupancyBounds(t *testing.T) {
	m := NewOccupancyMap()
	for _, p := range [][2]int{{-1, 0}, {0, -1}, {384, 0}, {0, 256}} {
		m.Enter(p[0], p[1])
		m.Leave(p[0], p[1])
		if m.At(p[0], p[1]) != 0 {
			t.Errorf("越界 (%d,%d) 不該有值", p[0], p[1])
		}
	}
	// 邊界內的極值要正常。
	m.Enter(383, 255)
	if m.At(383, 255) != 1 {
		t.Error("右下角應該可以放")
	}
}

// 行軍節拍：間隔 N 表示每 N 個 tick 走一步。
// 原版是「先減再判斷」（sub_125A3 的 dec / jnz）。
func TestMoveCadence(t *testing.T) {
	c := Corps{MoveTimer: 3, MoveInterval: 3}
	steps := 0
	for i := 0; i < 12; i++ {
		if c.Step() {
			steps++
		}
	}
	if steps != 4 {
		t.Errorf("間隔 3 跑 12 tick 走了 %d 步, want 4", steps)
	}

	// 間隔越小走得越快 —— 純騎馬編成靠的就是這個。
	fast := Corps{MoveTimer: 1, MoveInterval: 1}
	f := 0
	for i := 0; i < 12; i++ {
		if fast.Step() {
			f++
		}
	}
	if f != 12 {
		t.Errorf("間隔 1 跑 12 tick 走了 %d 步, want 12", f)
	}
	if f <= steps {
		t.Error("間隔小的應該走得比較快")
	}
}
