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

// RoadEdge 是兩個據點之間的一條路。Steps 是沿路的格數，可以直接當距離。
type RoadEdge struct {
	A, B  int
	Steps int
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
	for i := range owner {
		owner[i] = -1
	}

	queue := make([]int, 0, Width*Height)
	for ci, c := range cities {
		cell := m.nodeCell(c[0], c[1])
		if cell < 0 {
			return nil, fmt.Errorf("world: 據點 %d (%d,%d) 找不到節點格",
				ci, c[0], c[1])
		}
		owner[cell] = int32(ci)
		queue = append(queue, cell)
	}

	best := map[[2]int]int{}
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
				queue = append(queue, q)
			case owner[q] != owner[p]:
				a, b := int(owner[p]), int(owner[q])
				if a > b {
					a, b = b, a
				}
				w := int(dist[p]+dist[q]) + 1
				k := [2]int{a, b}
				if old, ok := best[k]; !ok || w < old {
					best[k] = w
				}
			}
		}
	}

	// 排序輸出，讓結果與 map 的走訪順序無關——不然同一份地圖每次
	// 跑出來的邊順序都不一樣，測試會忽好忽壞。
	out := make([]RoadEdge, 0, len(best))
	for k, w := range best {
		out = append(out, RoadEdge{A: k[0], B: k[1], Steps: w})
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && less(out[j], out[j-1]); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}

func less(a, b RoadEdge) bool {
	if a.A != b.A {
		return a.A < b.A
	}
	return a.B < b.B
}
