package tactical

import "testing"

// queueFixture 建一場只用來看佇列的戰鬥。
func queueFixture(t *testing.T) *Battle {
	t.Helper()
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{}, 0)
	for side := 0; side < 2; side++ {
		for k := 0; k < Squads; k++ {
			b.Deploy(side, k, Infantry, 8)
		}
	}
	return b
}

// ⭐ 去重：同一個兵排兩次只佔一格（原版 `sub_1C653` 的 `test [si], 10h`）。
func TestPathQueueDedupes(t *testing.T) {
	b := queueFixture(t)
	b.requestPath(0, 0)
	b.requestPath(0, 0)
	b.requestPath(0, 0)
	if b.paths.n != 1 {
		t.Errorf("同一個兵排了三次，佇列長度 %d，want 1", b.paths.n)
	}
	if !b.Sides[0].Soldiers[0].PathQueued {
		t.Error("排進去了卻沒有設 PathQueued")
	}
	// 消化之後旗標要清掉，下一次才排得進去。
	b.drainPathQueue()
	if b.Sides[0].Soldiers[0].PathQueued {
		t.Error("出隊之後 PathQueued 沒清——這個兵再也排不進佇列")
	}
	b.requestPath(0, 0)
	if b.paths.n != 1 {
		t.Errorf("清掉旗標之後排不進去，佇列長度 %d", b.paths.n)
	}
}

// ⭐ 每幀兩筆（`sub_1AED2` 的 `mov cx, 2`）。三筆要兩幀才消化得完。
func TestPathQueueDrainsTwoPerFrame(t *testing.T) {
	b := queueFixture(t)
	for k := 0; k < 3; k++ {
		b.requestPath(0, k)
	}
	if b.paths.n != 3 {
		t.Fatalf("排了三筆，佇列長度 %d", b.paths.n)
	}
	b.drainPathQueue()
	if b.paths.n != 1 {
		t.Errorf("一幀消化之後剩 %d 筆，want 1（每幀只算兩個兵）", b.paths.n)
	}
	b.drainPathQueue()
	if b.paths.n != 0 {
		t.Errorf("兩幀之後還剩 %d 筆", b.paths.n)
	}
	// 正對照：如果一次全消化完，上面第一個斷言會拿到 0——那是另一個實作。
}

// FIFO：先卡住的先算。
func TestPathQueueIsFIFO(t *testing.T) {
	q := &pathQueue{}
	for i := 0; i < 5; i++ {
		if !q.push(pathRequest{side: 0, k: i}) {
			t.Fatalf("第 %d 筆推不進去", i)
		}
	}
	for i := 0; i < 5; i++ {
		r, ok := q.pop()
		if !ok {
			t.Fatalf("第 %d 筆取不出來", i)
		}
		if r.k != i {
			t.Errorf("第 %d 次取出 k=%d，want %d——不是 FIFO", i, r.k, i)
		}
	}
	if _, ok := q.pop(); ok {
		t.Error("空佇列還取得出東西")
	}
}

// ⭐ **結構上塞不滿**：去重旗標讓在飛的請求最多 96 筆（兵的上限），
// 而佇列有 128 格。這一支同時證明環狀游標繞得回去。
func TestPathQueueNeverOverflows(t *testing.T) {
	b := queueFixture(t)
	total := 0
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			if b.Sides[side].Soldiers[k].Alive {
				b.requestPath(side, k)
				total++
			}
		}
	}
	if total == 0 {
		t.Fatal("fixture 一個活著的兵都沒有——這一支在驗空氣")
	}
	if total > pathQueueSize {
		t.Fatalf("活著的兵 %d 個，比佇列的 %d 格還多——前提不成立",
			total, pathQueueSize)
	}
	if b.paths.n != total {
		t.Errorf("%d 個兵各排一次，佇列長度 %d", total, b.paths.n)
	}

	// 環狀：全部消化完再排一輪，游標要繞回去而不是滿出來。
	for b.paths.n > 0 {
		b.drainPathQueue()
	}
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			if b.Sides[side].Soldiers[k].Alive {
				b.requestPath(side, k)
			}
		}
	}
	if b.paths.n != total {
		t.Errorf("第二輪只排進 %d 筆，want %d——環狀游標繞不回去", b.paths.n, total)
	}
}

// 排隊期間陣亡的兵要作廢，不能拿死人去算路。
func TestPathQueueSkipsDeadSoldiers(t *testing.T) {
	b := queueFixture(t)
	b.requestPath(0, 0)
	b.Sides[0].Soldiers[0].Alive = false
	b.drainPathQueue()
	if b.paths.n != 0 {
		t.Error("死掉的那一筆沒有從佇列裡拿掉")
	}
	if b.Sides[0].Soldiers[0].Path != nil {
		t.Error("幫死人算了一條路")
	}
}
