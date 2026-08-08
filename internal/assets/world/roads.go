package world

import "fmt"

// 道路網的推導。原版是在載入 `MMAP` 時現建的（`sub_1E4CE` 掃節點格、
// `sub_1E717` 建連結表、`sub_1E81C` 沿路走），這裡從同一份地圖資料
// 推出**等價的圖**：節點是 192 個據點，邊是「兩個據點之間有一條路直達」。
//
// ⚠ **這不是照抄那三支常式。** 原版還會建路上（192–255）與野外（≥256）
// 節點，並把每條路的路徑點與外接矩形算出來（`docs/re/08` §7.1–§7.3）。
// 這裡只取「據點與據點的連通關係 ＋ 步數」——行軍要的就是這個。
// 路徑點層級的忠實重現還沒做，**標成 remake 差異**。

const (
	// roadLo、roadHi 是道路格的圖塊值域。出自 `sub_1E4CE`：
	// 它檢查鄰格是否落在 `0xB8`–`0xDD`（`docs/formats/05` §3）。
	roadLo, roadHi = 0xB8, 0xDD

	// bridgeTile 是橋。**這個值不在上面的值域裡**，是量出來的：
	// 只用 `0xB8`–`0xDD` 時 192 個據點散成 27 群，把 `0xFF` 加進去
	// 就收成 1 群。全圖只有 84 格，都貼著 `0xFE` 成組出現在河上。
	bridgeTile = 0xFF

	// nodeLo、nodeHi 是「據點節點格」的圖塊值域（`sub_1E4CE` 的掃描條件）。
	// 全圖正好 192 格 —— 與據點數相同。
	nodeLo, nodeHi = 0xCB, 0xD3

	// nodeDX 是節點格相對於據點座標的偏移。
	// 據點記錄的 X 是城圖左緣，節點在右邊四格（191/192 個據點都是 +4，
	// 剩下一個是 −4）。
	nodeDX = 4
)

// RoadEdge 是兩個據點之間的一條路。
type RoadEdge struct {
	A, B  int
	Steps int // 沿路的格數，可以直接當距離

	// Path 是從 A 到 B 要經過的地圖格，**不含 A 的所在格、含 B 的所在格**。
	// 這樣把多條邊接起來時不會在中繼據點重複一格。
	//
	// 頭尾各有一小段是「城門到城中心」：據點記錄的座標是城圖左緣，
	// 節點格在 `(X+4, Y)`，所以路的兩端各接一段 ≤ 4 格的水平走法。
	Path [][2]int
}

func isRoad(v byte) bool {
	return (v >= roadLo && v <= roadHi) || v == bridgeTile
}

// nodeCell 找出據點 i 的節點格。回 −1 表示找不到。
func (m *Map) nodeCell(x, y int) int {
	for _, dx := range [...]int{nodeDX, -nodeDX} {
		nx := x + dx
		if nx < 0 || nx >= Width || y < 0 || y >= Height {
			continue
		}
		if v := m.Tiles[y*Width+nx]; v >= nodeLo && v <= nodeHi {
			return y*Width + nx
		}
	}
	return -1
}

// 八個方向。**道路是 8-連通的**——`sub_1E81C` 找下一格時依序試
// 西／東／北／南／西北／東北／西南／東南。
// 只用四方向的話道路的階梯狀轉折會斷開（實測 192 個據點散成 233 群）。
var neighbours8 = [8][2]int{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
}

// RoadEdges 從地圖推出據點之間的道路圖。
//
// cities 是 192 個據點的 (X, Y)，順序就是據點編號。
//
// 作法是**多源 BFS**：192 個節點格同時往外擴，每一格記下「離哪個據點最近」
// 與距離；兩個不同據點的領域碰在一起的地方就是一條邊，長度是
// `兩側距離相加 ＋ 1`。這等價於把道路網收縮成據點圖。
//
// 三個獨立的驗證（`internal/assets/world` 的測試釘住）：
//
//	據點記錄裡那 85 條鄰接（+0x00 遮罩 ＋ +0x1C–+0x1F）**全部**出現在結果裡
//	192 個據點**全連通**
//	每個據點的分支度 ≤ 4 —— 正好對上據點記錄只有四個方向槽
func RoadEdges(m *Map, cities [][2]int) ([]RoadEdge, error) {
	if len(m.Tiles) < Width*Height {
		return nil, fmt.Errorf("world: 地圖只有 %d 格", len(m.Tiles))
	}
	owner := make([]int32, Width*Height)
	dist := make([]int32, Width*Height)
	parent := make([]int32, Width*Height)
	for i := range owner {
		owner[i], parent[i] = -1, -1
	}

	nodeOf := make([]int, len(cities))
	queue := make([]int, 0, Width*Height)
	for ci, c := range cities {
		cell := m.nodeCell(c[0], c[1])
		if cell < 0 {
			return nil, fmt.Errorf("world: 據點 %d (%d,%d) 找不到節點格",
				ci, c[0], c[1])
		}
		nodeOf[ci] = cell
		owner[cell] = int32(ci)
		queue = append(queue, cell)
	}

	type meeting struct{ w, p, q int }
	best := map[[2]int]meeting{}
	for head := 0; head < len(queue); head++ {
		p := queue[head]
		px := p % Width
		for _, d := range neighbours8 {
			qx, qy := px+d[0], p/Width+d[1]
			if qx < 0 || qx >= Width || qy < 0 || qy >= Height {
				continue
			}
			q := qy*Width + qx
			if !isRoad(m.Tiles[q]) {
				continue
			}
			switch {
			case owner[q] < 0:
				owner[q] = owner[p]
				dist[q] = dist[p] + 1
				parent[q] = int32(p)
				queue = append(queue, q)
			case owner[q] != owner[p]:
				a, b := int(owner[p]), int(owner[q])
				pp, qq := p, q
				if a > b {
					a, b = b, a
					pp, qq = q, p
				}
				w := int(dist[p]+dist[q]) + 1
				k := [2]int{a, b}
				if old, ok := best[k]; !ok || w < old.w {
					best[k] = meeting{w: w, p: pp, q: qq}
				}
			}
		}
	}

	// 從相遇的兩格各自沿 parent 回到起點，接成一條完整的格子路徑。
	chain := func(cell int) [][2]int {
		var out [][2]int
		for c := int32(cell); c >= 0; c = parent[c] {
			out = append(out, [2]int{int(c) % Width, int(c) / Width})
		}
		return out
	}
	reverse := func(a [][2]int) {
		for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
			a[i], a[j] = a[j], a[i]
		}
	}

	out := make([]RoadEdge, 0, len(best))
	for k, mt := range best {
		// A 側：從 A 的節點格走到相遇點 p。
		left := chain(mt.p)
		reverse(left)
		// B 側：從相遇點 q 走到 B 的節點格。
		right := chain(mt.q)

		cells := append(append([][2]int{}, left...), right...)
		out = append(out, RoadEdge{
			A: k[0], B: k[1], Steps: mt.w,
			Path: withCityEnds(cells, cities[k[0]], cities[k[1]]),
		})
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && less(out[j], out[j-1]); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}

// withCityEnds 把「節點格到節點格」的路徑補成「據點座標到據點座標」。
//
// 節點格在 `(X+4, Y)`，而軍團的位置用的是據點記錄的 `(X, Y)`，
// 兩者差一小段水平距離。少了這一段，軍團出城與進城時會**跳 4 格**。
//
// 回傳的序列**不含起點格、含終點格**——接起來時中繼據點不會重複一格。
func withCityEnds(cells [][2]int, from, to [2]int) [][2]int {
	out := make([][2]int, 0, len(cells)+8)
	for x := from[0]; x != cells[0][0]; x += sign(cells[0][0] - x) {
		if x != from[0] {
			out = append(out, [2]int{x, from[1]})
		}
	}
	out = append(out, cells...)
	last := cells[len(cells)-1]
	for x := last[0]; x != to[0]; x += sign(to[0] - x) {
		if x != last[0] {
			out = append(out, [2]int{x, to[1]})
		}
	}
	out = append(out, to)
	return out
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
