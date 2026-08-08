package march

import "testing"

//   0 ──2── 1 ──2── 2
//   │               │
//   └───────9───────┘
func sample() *Graph {
	return New(4, []Edge{{0, 1, 2}, {1, 2, 2}, {0, 2, 9}})
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

func TestAdjacencyIsUndirected(t *testing.T) {
	g := sample()
	if !g.Adjacent(0, 1) || !g.Adjacent(1, 0) {
		t.Fatal("邊應該是雙向的")
	}
	if g.Adjacent(0, 3) {
		t.Fatal("沒有邊的兩點不該相鄰")
	}
}

func TestNilGraphIsSafe(t *testing.T) {
	// 沒有原版素材時圖是 nil，呼叫端會退回直線移動——**不能 panic**。
	var g *Graph
	if g.Route(0, 1) != nil || g.Adjacent(0, 1) || g.Neighbours(0) != nil {
		t.Fatal("nil 圖應該安靜地回空值")
	}
}
