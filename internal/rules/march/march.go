// Package march 是行軍的路徑規劃：在據點道路圖上找路。
//
// 圖從 `MMAP` 推導出來（`internal/assets/world` 的 RoadEdges），
// 這一層只負責找路，不認識地圖也不認識檔案格式。
//
// ⭐ **原版也找路，只是一次只找一步**：軍團記錄存著目標
// （`+0x14`／`+0x16`／`+0x18`），每次輪到移動就跑一次 `loc_1491B`
// ——一個**從目標往回**的 uniform-cost 搜尋——決定「往這條邊的哪一端走」
// （`docs/spec/192`）。
//
// remake 一次算好整條路再照著走。兩者在成本模型一致時等價，
// 所以 `RouteCost` 必須照原版的三項成本算（Σ邊長 ＋ 4×節點 ＋
// 0xA6×非己方據點），不能只用格數。
//
// remake 讓玩家可以直接點遠處的據點，這一層也負責那個路由——
// **多段規劃是操作方式的差異，不是規則的差異**：走的還是同一條路。
package march

import "container/heap"

// nodeCost 是原版 `loc_1491B` 的 `add dx, 4`：**每個節點固定 4**。
// 它讓「經過的節點少」勝過「格數少一點」——兩種度量會選出不同的路
// （docs/spec/192）。
const nodeCost = 4

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

	// BGate 是 **B 那一端的城門格**。`Path` 刻意不含它（走到它的同一拍
	// 就換成 B 中心了），但**反向走法從它開始**——
	// 反向序列不是正向序列的倒轉（docs/spec/169 §3.1.1）。
	BGate [2]int

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

	// Next 是**下一個路徑點**的座標。朝向（`+0x08`）看的是它，不是剛走過
	// 的那一格——原版 `sub_12804` 用 `+0x0C ＋ 步進` 取下一筆再比座標
	// （docs/spec/173 §1.2）。用「上一格到現在」算會整整晚一步，
	// 而**位置完全正確**，所以只有逐欄比對看得見。
	//
	// 最後一格（據點中心）沒有下一筆，是零值——到站時朝向本來就寫 4。
	Next [2]int

	// Point 是**這一格對應的原版路徑點**（`es:[bx]`／`es:[bx+2]`）。
	//
	// ⭐ 它與走過的格子只差終點那一筆：格子序列的最後一格是據點中心，
	// 路徑點序列的最後一筆是**那一端的城門格**。`sub_12708` 踏進去之前
	// 問的是路徑點，不是據點中心——守軍站在城門格上時，原版撞到的是
	// 軍團（野戰），而拿據點中心去問會撞到據點（攻城，docs/spec/186）。
	Point [2]int
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
// `Path` 與原版的路徑點**一對一**：最後一格的座標已經換成據點中心
// （`withCityEnds`），但它對應的仍是最後那一筆路徑點——原版在那一拍
// 先寫路徑點的座標、再被 `sub_127A2` 的換節點覆寫成據點中心。
// 反向走時指標從最後一筆遞減（`+0x0A` ＝ −4）。
func cellMarks(e Edge, n int, forward bool) []CellMark {
	if n <= 0 {
		return nil
	}
	// ⭐ **真正的路徑點序列**（不是走過的格子序列）。兩者只差終點那一格：
	// 格子序列的最後一格是據點中心，路徑點序列的最後一筆是那一端的城門格。
	// 朝向要用它算——`sub_12804` 看的是**下一個路徑點**（docs/spec/173 §1.2）。
	//   正向 q ＝ p0 … p(n-2), p(n-1)      ＝ Path[0…n-2] ＋ BGate
	//   反向 q ＝ p(n-1) … p1, p0           ＝ BGate ＋ Path[n-2…0]
	pts := make([][2]int, n)
	if forward {
		copy(pts, e.Path[:n-1])
		pts[n-1] = e.BGate
	} else {
		pts[0] = e.BGate
		for k := 1; k < n; k++ {
			pts[k] = e.Path[n-1-k]
		}
	}
	out := make([]CellMark, n)
	step := 4
	if !forward {
		step = -4
	}
	for k := range out {
		idx := k
		if !forward {
			idx = n - 1 - k
		}
		var next [2]int
		if k+1 < n {
			next = pts[k+1]
		}
		out[k] = CellMark{
			PathPtr:  e.PathAddr + idx*4,
			LinkAddr: e.LinkAddr,
			Step:     step,
			Next:     next,
			Point:    pts[k],
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
		n := len(e.Path) // 與原版的路徑點一對一（docs/spec/169 §3.1）
		if e.LinkAddr != 0 {
			g.byLink[e.LinkAddr] = [2]int{e.A, e.B}
		}
		g.adj[e.A] = append(g.adj[e.A],
			link{e.B, e.Steps, e.Path, cellMarks(e, n, true)})
		// ⚠ **反向不是正向的倒轉。** 正向序列是 `p0 … p(n-2) ＋ B 中心`，
		// 反向要的是 `p(n-1) … p1 ＋ A 中心`——兩個方向各吃掉自己終點
		// 那一格（docs/spec/169 §3.1.1）。倒轉會以 B 中心開頭、以 `p0`
		// 結尾，整條差一格，而**路徑指標照樣對得上**，所以只有逐格
		// 對拍看得見。
		g.adj[e.B] = append(g.adj[e.B],
			link{e.A, e.Steps, reversePath(e.Path, e.ACell, e.BGate), cellMarks(e, n, false)})
	}
	return g
}

// reversePath 把 A→B 的序列翻成 B→A：反轉之後去掉頭（原本的 B 格），
// 再把 A 的格子接到尾巴。
// reversePath 從正向序列 `p0 … p(n-2) ＋ B 中心` 造出反向序列
// `p(n-1) … p1 ＋ A 中心`（docs/spec/169 §3.1.1）。
//
// b 是 B 那一端的城門格（`p(n-1)`），正向序列裡沒有它；
// a 是 A 的所在格，反向的終點。
func reversePath(p [][2]int, a, b [2]int) [][2]int {
	if len(p) == 0 {
		return nil
	}
	out := make([][2]int, 0, len(p))
	out = append(out, b)                  // 反向的第一格 ＝ B 那一端的城門格
	for i := len(p) - 2; i >= 1; i-- {    // 中段：p(n-2) … p1
		out = append(out, p[i])
	}
	return append(out, a)                 // 反向的終點 ＝ A 中心
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
	return g.CellRouteMarkedCost(from, to, nil)
}

// CellRouteMarkedCost 與 CellRouteMarked 相同，但用 RouteCost 的成本模型。
func (g *Graph) CellRouteMarkedCost(from, to int, penalty func(int) int) ([][2]int, []CellMark) {
	route := g.RouteCost(from, to, penalty)
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
func (g *Graph) Route(from, to int) []int { return g.RouteCost(from, to, nil) }

// RouteCost 是 Route 帶上**每個節點的額外成本**的版本。
//
// ⭐ 原版的成本不是格數（`loc_1491B`，docs/spec/192）：
//
//	成本 ＝ Σ(邊長) ＋ 4 × 節點數 ＋ penalty(節點)
//
// 那個 `4` 是 `add dx, 4`，對**每個展開的節點**加一次；`penalty` 對應
// `add dx, 0A6h`——非己方的據點加 166，大到足以蓋過任何合理的邊長差。
//
// ⚠ **節點成本歸給「進入該節點的邊」**：原版是從**目標**往回搜、在
// 終止檢查**之前**累加，所以起點（軍團現在的位置）不算、目標算。
// 反過來搜要把它掛在 `l.to` 上才等價（docs/spec/192 §3）。
//
// penalty 是 nil 就只有那個 4。
func (g *Graph) RouteCost(from, to int, penalty func(city int) int) []int {
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
			w := l.steps + nodeCost
			if penalty != nil {
				w += penalty(l.to)
			}
			if n := it.cost + w; n < dist[l.to] {
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
