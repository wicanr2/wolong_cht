package march

import "testing"

//   0 ──2── 1 ──2── 2
//   │               │
//   └───────9───────┘
func sample() *Graph {
	return New(4, []Edge{
		{A: 0, B: 1, Steps: 2}, {A: 1, B: 2, Steps: 2}, {A: 0, B: 2, Steps: 9},
	})
}

func TestRoutePrefersShorterTotal(t *testing.T) {
	// 直達 9 格 vs 繞經 #1 共 4 格 —— 邊數多但總距離短的那條才對。
	got := sample().Route(0, 2)
	want := []int{0, 1, 2}
	if len(got) != len(want) {
		t.Fatalf("路線 %v，期望 %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("路線 %v，期望 %v", got, want)
		}
	}
	if d := sample().Distance(0, 2); d != 4 {
		t.Fatalf("距離 %d，期望 4", d)
	}
}

// 「已經到了」與「走不到」是兩件事。
func TestSameNodeIsNotUnreachable(t *testing.T) {
	if r := sample().Route(1, 1); len(r) != 1 || r[0] != 1 {
		t.Fatalf("同一個據點應回長度 1 的序列，得到 %v", r)
	}
	if r := sample().Route(0, 3); r != nil {
		t.Fatalf("孤立的據點應回 nil，得到 %v", r)
	}
	if d := sample().Distance(0, 3); d != -1 {
		t.Fatalf("走不到的距離應是 −1，得到 %d", d)
	}
}


func TestNilGraphIsSafe(t *testing.T) {
	// 沒有原版素材時圖是 nil，呼叫端會退回直線移動——**不能 panic**。
	var g *Graph
	if g.Route(0, 1) != nil || g.CellRoute(0, 1) != nil || g.Distance(0, 1) != -1 {
		t.Fatal("nil 圖應該安靜地回空值")
	}
}

// ⭐ 反向的格子序列要正確：從 B 走回 A 應該是同一條路倒著走，
// 而且結尾要落在 A 的所在格。
//
// 序列的約定是「不含起點、含終點」，所以**直接反轉是錯的**——
// 那會少掉 A 的格子、多出 B 的格子。這條把那個 off-by-one 釘住。
func TestReversePathEndsAtOrigin(t *testing.T) {
	a := [2]int{10, 10}
	b := [2]int{13, 10}
	g := New(2, []Edge{{
		A: 0, B: 1, Steps: 3, ACell: a,
		Path: [][2]int{{11, 10}, {12, 10}, b},
	}})

	fwd := g.CellRoute(0, 1)
	if len(fwd) != 3 || fwd[len(fwd)-1] != b {
		t.Fatalf("正向 %v，應以 %v 結尾", fwd, b)
	}
	back := g.CellRoute(1, 0)
	if len(back) != 3 || back[len(back)-1] != a {
		t.Fatalf("反向 %v，應以 %v 結尾", back, a)
	}
	// 中間那兩格要是同一組，順序相反。
	if back[0] != fwd[1] || back[1] != fwd[0] {
		t.Fatalf("反向的中段不是正向倒著走：%v vs %v", back, fwd)
	}
}

// 沒有格子序列時 CellRoute 要回 nil，讓呼叫端退回直線。
func TestCellRouteWithoutPaths(t *testing.T) {
	if r := sample().CellRoute(0, 2); r != nil {
		t.Fatalf("沒有 Path 應回 nil，得到 %v", r)
	}
}
