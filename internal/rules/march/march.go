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

	// LinkAddr、PathAddr 是這條邊在原版執行期道路表裡的段內位址。
	// 只用來還原軍團記錄的 `+0x0C`／`+0x0E`（docs/spec/172），
	// 路由本身不看它們。
	LinkAddr, PathAddr int
}

// CellMark 是路徑上某一格對應到**原版道路表**的位置。
//
// 軍團記錄的三個欄位就是它：`+0x0C` ＝ PathPtr、`+0x0E` ＝ LinkAddr、
// `+0x0A` ＝ Step（`+4` 正向、`−4` 反向，byte 存成 `0xFC`）。
type CellMark struct {
	PathPtr  int
	LinkAddr int
	Step     int
}

// Graph 是據點道路圖。
type Graph struct {
	adj  [][]link
	byLink map[int][2]int // 連結記錄位址 → (A, B)
}

// EdgeByLink 由**原版連結記錄的位址**反查這條邊的兩端（docs/spec/172）。
//
// 軍團記錄的 `+0x0E` 在行軍中放的就是這個位址，所以載入存檔時要靠它
// 才知道那支軍團正走在哪一條路上。
func (g *Graph) EdgeByLink(addr int) (a, b int, ok bool) {
	if g == nil || g.byLink == nil {
		return 0, 0, false
	}
	e, ok := g.byLink[addr]
	return e[0], e[1], ok
}

type link struct {
	to, steps int
	path      [][2]int // 從**這條 link 的起點**走到 to 的格子序列
	// marks 與 path 逐格對應，指回原版道路表的位置（docs/spec/172）。
	marks []CellMark
}

// cellMarks 算出一條 link 每一格對應的原版路徑點。
//
// `Path` 的最後一格是**終點的據點中心**（`withCityEnds`），它不在原版的
// 路徑點表裡——原版到站時 `+0x0C` 留著最後一筆的位址，所以那一格沿用
// 前一格的指標。反向走時指標從最後一筆遞減（`+0x0A` ＝ −4）。
func cellMarks(e Edge, n int, forward bool) []CellMark {
	if n <= 0 {
		return nil
	}
	out := make([]CellMark, n+1)
	step := 4
	if !forward {
		step = -4
	}
	for k := range out {
		idx := k
		if !forward {
			idx = n - 1 - k
		}
		if idx > n-1 {
			idx = n - 1
		}
		if idx < 0 {
			idx = 0
		}
		out[k] = CellMark{
			PathPtr:  e.PathAddr + idx*4,
			LinkAddr: e.LinkAddr,
			Step:     step,
		}
	}
	return out
}

// New 從邊清單建圖。n 是據點數。
func New(n int, edges []Edge) *Graph {
	g := &Graph{adj: make([][]link, n), byLink: map[int][2]int{}}
	for _, e := range edges {
		if e.A < 0 || e.A >= n || e.B < 0 || e.B >= n {
			continue
		}
		n := len(e.Path) - 1 // 原版的路徑點數（`Path` 多了終點城中心）
		if e.LinkAddr != 0 {
			g.byLink[e.LinkAddr] = [2]int{e.A, e.B}
		}
		g.adj[e.A] = append(g.adj[e.A],
			link{e.B, e.Steps, e.Path, cellMarks(e, n, true)})
		// 反向要把格子序列倒過來，而且**最後一格換成起點**：
		// 序列的約定是「不含起點、含終點」，直接反轉會少了 A 的格子、
		// 多出 B 的格子。
		g.adj[e.B] = append(g.adj[e.B],
			link{e.A, e.Steps, reversePath(e.Path, e.ACell), cellMarks(e, n, false)})
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
	out, _ := g.cellRoute(route)
	return out
}

// CellRouteMarked 與 CellRoute 相同，另外回傳每一格在**原版道路表**裡的
// 位置（docs/spec/172）。兩個序列等長。
func (g *Graph) CellRouteMarked(from, to int) ([][2]int, []CellMark) {
	route := g.Route(from, to)
	if len(route) < 2 {
		return nil, nil
	}
	return g.cellRoute(route)
}

func (g *Graph) cellRoute(route []int) ([][2]int, []CellMark) {
	var out [][2]int
	var marks []CellMark
	for i := 0; i+1 < len(route); i++ {
		var seg [][2]int
		var mk []CellMark
		for _, l := range g.adj[route[i]] {
			if l.to == route[i+1] {
				seg, mk = l.path, l.marks
				break
			}
		}
		if seg == nil {
			return nil, nil // 有一段沒有格子序列 → 整條都不用
		}
		out = append(out, seg...)
		if len(mk) == len(seg) {
			marks = append(marks, mk...)
		} else {
			marks = append(marks, make([]CellMark, len(seg))...)
		}
	}
	return out, marks
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
