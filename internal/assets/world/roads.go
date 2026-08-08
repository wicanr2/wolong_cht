package world

import "fmt"

// 道路網。**這是原版建表常式的直接移植**，不是等價的重寫：
//
//	sub_1E4CE  逐格掃地圖找據點的節點格（row-major）
//	sub_1E57F  對節點格的四個方向找「城門格」（先試 1 格再試 2 格）
//	sub_1E81C  從城門格沿路走到對面的城門格，逐格記下座標
//	sub_1E961  判定一格能不能走，並回傳它的類別
//
// 出處與推導見 `docs/re/08` §7.1–§7.7。
//
// ⚠ 先前這裡是**多源 BFS 的等價重寫**。逐條比對之後推翻：
// BFS 少找到一條邊（會稽–章安），而且共同的 253 條裡有 **109 條**
// 路徑長度差超過 2 格——BFS 會抄近路，原版照著畫出來的路走。
// **「拓樸一樣」不等於「路一樣」。**

const (
	// 圖塊的四個類別，出自 `sub_1E961` 的連續比較。
	// **走訪繼續的條件是 `(類別 & 7) ≤ 1`**，所以類 0／1 繼續、類 3／4 停。
	//
	//	類 1  0xB8–0xB9  可走
	//	類 0  0xBA–0xCA  可走（一般道路）
	//	類 3  0xCB–0xD3  終端 —— **全圖正好 192 格 ＝ 據點的節點格**
	//	類 4  0xD4–0xDD  終端 —— **全圖正好 508 格 ＝ 城門格**
	//
	// 值域外（< 0xB8 或 > 0xDD）一律不可走。
	roadLo, roadHi = 0xB8, 0xDD
	nodeLo, nodeHi = 0xCB, 0xD3
	gateLo         = 0xD4

	// nodeDX 是節點格相對於據點座標的偏移。據點記錄的 X 是城圖左緣，
	// 節點在右邊四格（191/192 個是 `+4`，剩下一個是 `−4`）。
	nodeDX = 4
)

// RoadEdge 是兩個據點之間的一條路。
type RoadEdge struct {
	A, B  int
	Steps int // 路徑格數，可以直接當距離

	// Path 是從 A 到 B 要經過的地圖格，**不含 A 的所在格、含 B 的所在格**。
	// 這樣把多條邊接起來時中繼據點不會重複一格。
	//
	// 頭尾各有一小段是「城中心 → 節點格 → 城門格」的直線，
	// 那幾格踩在城池圖形上而不是道路上（Stub 記著長度）。
	Path [][2]int

	// StubA、StubB 是頭尾那兩段的格數。要逐格檢查「有沒有踩在路上」時
	// 得把它們排除——**城池圖形本身不是道路圖塊**。
	StubA, StubB int
}

// tileClass 是 `sub_1E961` 的分類。回 −1 表示不可走。
func tileClass(v byte) int {
	switch {
	case v < roadLo || v > roadHi:
		return -1
	case v < 0xBA:
		return 1
	case v < nodeLo:
		return 0
	case v < gateLo:
		return 3
	default:
		return 4
	}
}

func isRoad(v byte) bool { return tileClass(v) >= 0 }

// 八個方向的偏移，**順序就是 `sub_1E81C` 的嘗試順序**：
// 西 → 東 → 北 → 南 → 西北 → 東北 → 西南 → 東南。
//
// 順序是行為的一部分：走到分岔時第一個能走的方向就是它要走的方向。
// 換個順序會得到一條不同的路，而長度看起來還是差不多。
var stepOrder = [8]int{-1, 1, -Width, Width, -Width - 1, -Width + 1, Width - 1, Width + 1}

// 四個起步方向的偏移（`sub_1E57F`）。**每個方向先試 1 格再試 2 格**——
// 城門格不一定緊貼節點格。
var gateProbe = [4][2]int{
	{-1, -2}, {1, 2}, {-Width, -2 * Width}, {Width, 2 * Width},
}

// 起步方向對應的第一步（`sub_1E81C` 的 `dh & 0x0F`）。
var firstStep = [4]int{-1, 1, -Width, Width}

// RoadEdges 從地圖建出據點之間的道路圖，含每條路的逐格路徑。
//
// cities 是 192 個據點的 (X, Y)，順序就是據點編號。
//
// 五個獨立的驗證（見 `world_test.go`）：
//
//	據點記錄裡那 85 條鄰接全部出現       85 / 85
//	192 個據點全連通                     192 / 192
//	每個據點的分支度 ≤ 4（記錄只有四槽）  最大 4
//	城門格數 ＝ 類 4 的格數              508 / 508
//	路徑逐格連續、且中段全踩在道路圖塊上
func RoadEdges(m *Map, cities [][2]int) ([]RoadEdge, error) {
	if len(m.Tiles) < Width*Height {
		return nil, fmt.Errorf("world: 地圖只有 %d 格", len(m.Tiles))
	}
	t := m.Tiles

	// ① 每個據點的節點格。
	nodeOf := make([]int, len(cities))
	for ci, c := range cities {
		cell := -1
		for _, dx := range [...]int{nodeDX, -nodeDX} {
			nx := c[0] + dx
			if nx < 0 || nx >= Width || c[1] < 0 || c[1] >= Height {
				continue
			}
			if v := t[c[1]*Width+nx]; v >= nodeLo && v <= nodeHi {
				cell = c[1]*Width + nx
				break
			}
		}
		if cell < 0 {
			return nil, fmt.Errorf("world: 據點 %d (%d,%d) 找不到節點格",
				ci, c[0], c[1])
		}
		nodeOf[ci] = cell
	}

	// ② 城門格。`sub_1E57F` 對四個方向各找一格，先 1 格再 2 格。
	type start struct {
		city, gate, dir int
	}
	gateCity := map[int]int{}
	var starts []start
	for ci, node := range nodeOf {
		for d, probe := range gateProbe {
			for _, off := range probe {
				q := node + off
				if q < 0 || q >= Width*Height || !isRoad(t[q]) {
					continue
				}
				gateCity[q] = ci
				starts = append(starts, start{city: ci, gate: q, dir: d})
				break
			}
		}
	}

	// ③ 從每個城門格走一次。
	best := map[[2]int][][2]int{}
	stub := map[[2]int][2]int{}
	for _, s := range starts {
		cells, end := walkRoad(t, s.gate, s.dir)
		other, ok := gateCity[end]
		if !ok || other == s.city {
			continue // 走到死路或繞回自己
		}
		a, b := s.city, other
		reversed := false
		if a > b {
			a, b = b, a
			reversed = true
		}
		k := [2]int{a, b}
		if old, ok := best[k]; ok && len(old) <= len(cells) {
			continue
		}
		if reversed {
			for i, j := 0, len(cells)-1; i < j; i, j = i+1, j-1 {
				cells[i], cells[j] = cells[j], cells[i]
			}
		}
		full, sa, sb := withCityEnds(cells, nodeOf[a], nodeOf[b], cities[a], cities[b])
		best[k] = full
		stub[k] = [2]int{sa, sb}
	}

	out := make([]RoadEdge, 0, len(best))
	for k, path := range best {
		out = append(out, RoadEdge{
			A: k[0], B: k[1], Steps: len(path), Path: path,
			StubA: stub[k][0], StubB: stub[k][1],
		})
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && less(out[j], out[j-1]); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}

// walkRoad 是 `sub_1E81C`：從城門格出發沿路走，回傳走過的格子（含起點）
// 與停下來的那一格。
//
// 兩條規則都是行為的一部分，不能簡化：
//
//	**不能走回上一格**（`cmp si, di / jz`）。少了它會原地來回。
//	**八個方向有固定優先序**。分岔時走第一個能走的，不是走最近的。
func walkRoad(t []byte, gate, dir int) ([]int, int) {
	si, prev := gate, -1
	var path []int

	// `cmp al, 0D4h / jnb` —— 起步格是城門格（類 4）就先記下再踏出第一步。
	if t[si] >= gateLo {
		path = append(path, si)
		prev = si
		si += firstStep[dir]
	}
	for range Width * Height {
		if si < 0 || si >= len(t) || si == prev {
			return path, si
		}
		c := tileClass(t[si])
		path = append(path, si)
		if c < 0 || c > 1 {
			return path, si // 踩到終端（節點格或城門格）
		}
		next := -1
		for _, off := range stepOrder {
			q := si + off
			// ⚠ **x 位移不能超過 1。** 原版的地圖定址是分段的（一列一段），
			// 天然就是二維；這裡用平面陣列模擬，`si ± 1` 在列邊界會
			// **繞到下一列的另一端**。實測有一條路因此從 (383,190)
			// 跳到 (0,192)。加這個限制不是偏離原版，是把平面模型
			// 修正回原版本來的二維語意。
			if q < 0 || q >= len(t) || q == prev || abs(q%Width-si%Width) > 1 {
				continue
			}
			if isRoad(t[q]) {
				next = q
				break
			}
		}
		if next < 0 {
			return path, si
		}
		prev, si = si, next
	}
	return path, si
}

// withCityEnds 把「城門格到城門格」的路徑補成「據點座標到據點座標」，
// 並回傳頭尾兩段的長度。
//
// 軍團的位置用的是據點記錄的 (X, Y)，而路是從城門格開始的。
// 少了這兩段，軍團出城與進城時會**跳好幾格**。
//
// 回傳的序列**不含起點格、含終點格**——接起來時中繼據點不會重複一格。
func withCityEnds(cells []int, nodeA, nodeB int, from, to [2]int) ([][2]int, int, int) {
	// ⚠ `between` 與 `straight` 都是「不含起點、含終點」，
	// 所以 `節點格 → 城門格` 已經含了城門格，而 `cells[0]` 也是城門格。
	// **接的時候要跳過 `cells[0]`**，不然路徑裡會出現重複的一格——
	// 那在畫面上看不出來（軍團原地停一拍），但連續性檢查會抓到。
	head := straight(from, cellXY(nodeA))            // 城中心 → 節點格
	head = append(head, between(nodeA, cells[0])...) // 節點格 → 城門格（含）
	tail := between(cells[len(cells)-1], nodeB)      // 城門格 → 節點格
	tail = append(tail, straight(cellXY(nodeB), to)...)

	out := make([][2]int, 0, len(head)+len(cells)+len(tail))
	out = append(out, head...)
	for _, c := range cells[1:] {
		out = append(out, cellXY(c))
	}
	out = append(out, tail...)
	return out, len(head), len(tail)
}

func cellXY(c int) [2]int { return [2]int{c % Width, c / Width} }

// straight 產生 from（不含）到 to（含）的直線格子。兩點共用一個座標軸時
// 就是一條直線；不共用時走對角再補直——城門那一小段只會是前者。
func straight(from, to [2]int) [][2]int {
	var out [][2]int
	x, y := from[0], from[1]
	for x != to[0] || y != to[1] {
		x += sign(to[0] - x)
		y += sign(to[1] - y)
		out = append(out, [2]int{x, y})
	}
	return out
}

func between(a, b int) [][2]int { return straight(cellXY(a), cellXY(b)) }

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func sign(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

func less(a, b RoadEdge) bool {
	if a.A != b.A {
		return a.A < b.A
	}
	return a.B < b.B
}
