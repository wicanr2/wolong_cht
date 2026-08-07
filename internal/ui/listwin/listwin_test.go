package listwin

import "testing"

// 一組固定的測試資料：值越小排越前面。
func cols(v []int) []Column {
	return []Column{
		{Title: "編號", Less: func(a, b int) bool { return a < b }},
		{Title: "值", Less: func(a, b int) bool { return v[a] < v[b] }},
	}
}

func newList(mem *Memory) (*List, []int) {
	v := []int{30, 10, 20, 10}
	return New(Generals, cols(v), []int{0, 1, 2, 3}, 3, mem), v
}

// ⭐ 兩段式選取：第一次左鍵只反白，第二次才決定（說明書 3.8）。
func TestTwoStageSelection(t *testing.T) {
	l, _ := newList(nil)
	if l.Phase() != Browsing {
		t.Fatal("一開始應該是瀏覽狀態")
	}
	if _, ok := l.Confirm(); ok {
		t.Error("第一次左鍵不該決定")
	}
	if l.Phase() != Selected {
		t.Error("第一次左鍵後應該進入選取狀態")
	}
	got, ok := l.Confirm()
	if !ok {
		t.Fatal("第二次左鍵應該決定")
	}
	if got != 0 {
		t.Errorf("決定的是 %d, want 0", got)
	}
	if l.Phase() != Browsing {
		t.Error("決定之後應該回到瀏覽狀態")
	}
}

// ⭐ 反白狀態下取消只退回選取層，不關視窗（說明書 3.8）。
func TestCancelFromSelectedDoesNotClose(t *testing.T) {
	l, _ := newList(nil)
	l.Confirm() // 進入選取
	if l.Cancel() {
		t.Error("反白狀態下取消不該關視窗")
	}
	if l.Phase() != Browsing {
		t.Error("取消之後應該回到瀏覽狀態")
	}
	if !l.Cancel() {
		t.Error("瀏覽狀態下取消才該關視窗")
	}
}

// 移動游標會取消選取 —— 選取的那一列不再是游標所指的列。
func TestMoveClearsSelection(t *testing.T) {
	l, _ := newList(nil)
	l.Confirm()
	l.Move(1)
	if l.Phase() != Browsing {
		t.Error("移動游標後不該還在選取狀態")
	}
}

// 排序：點欄位名排序，再點同一欄反向。
func TestSortAndReverse(t *testing.T) {
	var mem Memory
	l, v := newList(&mem)

	l.SortBy(1) // 依值升冪：10(1), 10(3), 20(2), 30(0)
	want := []int{1, 3, 2, 0}
	for i, w := range want {
		if l.Rows[i] != w {
			t.Fatalf("升冪 = %v, want %v（值 %v）", l.Rows, want, v)
		}
	}

	l.SortBy(1) // 再點一次反向
	want = []int{0, 2, 1, 3}
	for i, w := range want {
		if l.Rows[i] != w {
			t.Fatalf("降冪 = %v, want %v", l.Rows, want)
		}
	}
}

// 穩定排序：同鍵值要保持原本的相對次序，畫面才不會跳。
func TestSortIsStable(t *testing.T) {
	var mem Memory
	l, _ := newList(&mem)
	l.SortBy(1)
	// 值同為 10 的是 1 與 3，原本 1 在 3 前面。
	i1, i3 := -1, -1
	for i, r := range l.Rows {
		if r == 1 {
			i1 = i
		}
		if r == 3 {
			i3 = i
		}
	}
	if i1 > i3 {
		t.Errorf("同值的相對次序被打亂：%v", l.Rows)
	}
}

// ⭐ 排序狀態以視窗種類為單位記住，重開不用再排（說明書 3.8）。
func TestSortMemoryPerKind(t *testing.T) {
	var mem Memory
	l, v := newList(&mem)
	l.SortBy(1)

	// 關掉再開同一種視窗 → 應該已經排好。
	l2 := New(Generals, cols(v), []int{0, 1, 2, 3}, 3, &mem)
	want := []int{1, 3, 2, 0}
	for i, w := range want {
		if l2.Rows[i] != w {
			t.Fatalf("重開後 = %v, want %v（排序狀態應該記住）", l2.Rows, want)
		}
	}

	// 不同種類的視窗不受影響。
	l3 := New(Cities, cols(v), []int{0, 1, 2, 3}, 3, &mem)
	for i := range l3.Rows {
		if l3.Rows[i] != i {
			t.Fatalf("不同種類的視窗不該被影響：%v", l3.Rows)
		}
	}
}

// 捲動：游標移出可視範圍時要跟著捲。
func TestScroll(t *testing.T) {
	l, _ := newList(nil) // Height = 3，共 4 列
	if _, first := l.Visible(); first != 0 {
		t.Errorf("一開始 first = %d, want 0", first)
	}
	l.Move(3) // 移到最後一列
	rows, first := l.Visible()
	if first != 1 || len(rows) != 3 {
		t.Errorf("捲動後 first = %d 顯示 %d 列, want 1 / 3", first, len(rows))
	}
	l.Move(-3)
	if _, first := l.Visible(); first != 0 {
		t.Errorf("捲回去 first = %d, want 0", first)
	}
}

// 游標不會跑出範圍。
func TestCursorBounds(t *testing.T) {
	l, _ := newList(nil)
	l.Move(-99)
	if l.Cursor != 0 {
		t.Errorf("往上超出 → %d, want 0", l.Cursor)
	}
	l.Move(99)
	if l.Cursor != len(l.Rows)-1 {
		t.Errorf("往下超出 → %d, want %d", l.Cursor, len(l.Rows)-1)
	}
}

// 空清單不能 panic，而且不該回傳有效選擇。
func TestEmptyList(t *testing.T) {
	l := New(Corps, cols(nil), nil, 5, nil)
	l.Move(1)
	if _, ok := l.Confirm(); ok {
		t.Error("空清單不該決定出東西")
	}
	if l.Selection() != -1 {
		t.Errorf("空清單的選擇 = %d, want -1", l.Selection())
	}
	if !l.Cancel() {
		t.Error("空清單取消應該關視窗")
	}
	if rows, _ := l.Visible(); rows != nil {
		t.Error("空清單不該有可視列")
	}
}

// 排序之後游標回到最上面 —— 停在原位會指到語意不同的資料。
func TestSortResetsCursor(t *testing.T) {
	var mem Memory
	l, _ := newList(&mem)
	l.Move(3)
	l.SortBy(1)
	if l.Cursor != 0 || l.Top != 0 {
		t.Errorf("排序後 cursor=%d top=%d, want 0/0", l.Cursor, l.Top)
	}
}
