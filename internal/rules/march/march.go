// Package march 是行軍的路徑規劃：在據點道路圖上找路。
//
// 圖從 `MMAP` 推導出來（`internal/assets/world` 的 RoadEdges），
// 這一層只負責找路，不認識地圖也不認識檔案格式。
//
// ⚠ **原版沒有「找路」這一步。** 它的軍團記錄裡直接存著目標
// （`+0x14`／`+0x16`／`+0x18`），沿著載入時建好的連結表一段一段走
// （`docs/re/08` §7）。玩家下指令時選的是**相鄰的據點**，
// 所以原版不需要跨多段的規劃。
//
// remake 讓玩家可以直接點遠處的據點，中間的路由這裡算——
// **這是操作方式的差異，不是規則的差異**：走的還是同一條路、
// 同一個距離，只是不必一段一段點。
package march

import "container/heap"

// Edge 是一條路。
type Edge struct {
	A, B  int
	Steps int
	// Path 是從 A 到 B 的地圖格序列，**不含 A 的所在格、含 B 的**。
	// 可以是 nil —— 那時 CellRoute 回 nil，呼叫端退回直線移動。
	Path [][2]int

	// ACell 是 A 的所在格。Path 刻意不含它（不然接兩段路時中繼點會重複），
	// 但反向的序列需要它當結尾，所以要另外帶進來。
	ACell [2]int
}

// Graph 是據點道路圖。
type Graph struct {
	adj [][]link
}

type link struct {
	to, steps int
	path      [][2]int // 從**這條 link 的起點**走到 to 的格子序列
}

// New 從邊清單建圖。n 是據點數。
func New(n int, edges []Edge) *Graph {
	g := &Graph{adj: make([][]link, n)}
	for _, e := range edges {
		if e.A < 0 || e.A >= n || e.B < 0 || e.B >= n {
			continue
		}
		g.adj[e.A] = append(g.adj[e.A], link{e.B, e.Steps, e.Path})
		// 反向要把格子序列倒過來，而且**最後一格換成起點**：
		// 序列的約定是「不含起點、含終點」，直接反轉會少了 A 的格子、
		// 多出 B 的格子。
		g.adj[e.B] = append(g.adj[e.B], link{e.A, e.Steps, reversePath(e.Path, e.ACell)})
	}
	return g
}

// reversePath 把 A→B 的序列翻成 B→A：反轉之後去掉頭（原本的 B 格），
// 再把 A 的格子接到尾巴。
func reversePath(p [][2]int, a [2]int) [][2]int {
	if len(p) == 0 {
		return nil
	}
	out := make([][2]int, 0, len(p))
	for i := len(p) - 2; i >= 0; i-- {
		out = append(out, p[i])
	}
	return append(out, a)
}

// CellRoute 回傳 from 走到 to 要經過的**每一格**，不含 from 的所在格。
// 走不到、或圖裡沒有格子序列時回 nil。
func (g *Graph) CellRoute(from, to int) [][2]int {
	route := g.Route(from, to)
	if len(route) < 2 {
		return nil
	}
	var out [][2]int
	for i := 0; i+1 < len(route); i++ {
		var seg [][2]int
		for _, l := range g.adj[route[i]] {
			if l.to == route[i+1] {
				seg = l.path
				break
			}
		}
		if seg == nil {
			return nil // 有一段沒有格子序列 → 整條都不用
		}
		out = append(out, seg...)
	}
	return out
}

// Neighbours 回傳與 node 直接相連的據點。
func (g *Graph) Neighbours(node int) []int {
	if g == nil || node < 0 || node >= len(g.adj) {
		return nil
	}
	out := make([]int, 0, len(g.adj[node]))
	for _, l := range g.adj[node] {
		out = append(out, l.to)
	}
	return out
}

// Adjacent 回報兩個據點是不是直接相連。
func (g *Graph) Adjacent(a, b int) bool {
	if g == nil || a < 0 || a >= len(g.adj) {
		return false
	}
	for _, l := range g.adj[a] {
		if l.to == b {
			return true
		}
	}
	return false
}

// Route 找 from 到 to 的最短路，回傳**含頭尾**的據點序列。
//
// 走不到回 nil。from == to 回長度 1 的序列——
// **不是 nil**：「已經到了」與「走不到」是兩件事，呼叫端要分得出來。
func (g *Graph) Route(from, to int) []int {
	if g == nil || from < 0 || from >= len(g.adj) || to < 0 || to >= len(g.adj) {
		return nil
	}
	if from == to {
		return []int{from}
	}
	const inf = int(^uint(0) >> 1)
	dist := make([]int, len(g.adj))
	prev := make([]int, len(g.adj))
	for i := range dist {
		dist[i], prev[i] = inf, -1
	}
	dist[from] = 0

	pq := &queue{{node: from}}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(item)
		if it.cost > dist[it.node] {
			continue
		}
		if it.node == to {
			break
		}
		for _, l := range g.adj[it.node] {
			if n := it.cost + l.steps; n < dist[l.to] {
				dist[l.to] = n
				prev[l.to] = it.node
				heap.Push(pq, item{node: l.to, cost: n})
			}
		}
	}
	if dist[to] == inf {
		return nil
	}
	var path []int
	for n := to; n != -1; n = prev[n] {
		path = append(path, n)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// Distance 回傳最短路的總格數，走不到回 −1。
func (g *Graph) Distance(from, to int) int {
	path := g.Route(from, to)
	if path == nil {
		return -1
	}
	total := 0
	for i := 0; i+1 < len(path); i++ {
		for _, l := range g.adj[path[i]] {
			if l.to == path[i+1] {
				total += l.steps
				break
			}
		}
	}
	return total
}

type item struct{ node, cost int }

type queue []item

func (q queue) Len() int            { return len(q) }
func (q queue) Less(i, j int) bool  { return q[i].cost < q[j].cost }
func (q queue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *queue) Push(x any)         { *q = append(*q, x.(item)) }
func (q *queue) Pop() any           { old := *q; n := len(old); it := old[n-1]; *q = old[:n-1]; return it }
