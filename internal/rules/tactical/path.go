package tactical

// 尋路：波前擴散 ＋ 回溯，出處是 `loc_1BD46`（docs/re/11 §5.15）。
//
// 原版的流程分兩段：
//
//  1. **擴散**：把成本圖（`ds:0D300`，16 KB）全部填成 `0xFFFF`，
//     從起點放 1，用一個環狀佇列往外長。每走一格加的**不是 1，
//     而是那一格的地形成本**（`mov al, es:[bx+2000h]`）——
//     所以這是 uniform-cost 搜尋，不是單純的 BFS。
//  2. **回溯**：從終點沿著遞減的波數走回去，
//     **只有轉彎時才寫一個點**（`ch` 記著上一步是橫的還是縱的）。
//
// ⭐ 回溯的計數器是 `mov cl, 40h` ＝ **64 個點**，一個點 2 byte
// ＝ 128 byte，正好是每個兵那塊 `0x1800 + 兵編號 × 128` 的大小（§5.8k）。
// 兩邊對得上，這條路徑才走得完。
//
// ⚠ **地形成本圖還沒解**（原版在 `ds:0D2FE`，由 `sub_1BBA6` 建）。
// 這裡的 Cost 預設一律回 1，行為等於 BFS；解出來之後換掉 Cost 就好。

// MaxWaypoints 是一條路徑最多幾個轉彎點（原版 `mov cl, 40h`）。
const MaxWaypoints = 64

// Point 是戰場上的一格。
type Point struct{ X, Y int }

// Cost 是走進 (x, y) 要付的代價。原版從一張每格一 byte 的表讀
// （`es:[bx+2000h]`）；回 0 表示不可通行。
type Cost func(x, y int) int

// pathNode 是擴散時每一格的狀態。
type pathNode struct {
	cost int // 0xFFFF ＝ 還沒走到
	from int // 從哪一格過來的（回溯用；−1 ＝ 起點）
}

const unreached = 0xFFFF

// FindPath 從 from 走到 to，回傳**轉彎點**的清單（不含起點）。
//
// climb 為真表示這個兵爬得上城牆（`cmp [si+4], 12h / jbe` 的另一側，
// 見 Soldier.CanClimb）——原版是靠**自我修改碼**切換的：
// `mov cs:byte_1BE33, cl` 把一個 `jz`（0x74）改成 `jmp short`（0xEB），
// 直接跳過 Z 軸那一支。
//
// 走不到就回 nil。點數超過 MaxWaypoints 時只回前 64 個——與原版一致
// （原版的緩衝區就只有那麼大，`dec cl / jz` 到了就停）。
func (f *Field) FindPath(from, to Point, climb bool, cost Cost) []Point {
	if !inBounds(from.X, from.Y) || !inBounds(to.X, to.Y) {
		return nil
	}
	if from == to {
		return nil
	}
	if cost == nil {
		cost = func(int, int) int { return 1 }
	}

	idx := func(x, y int) int { return y*Width + x }
	nodes := make([]pathNode, Width*Height)
	for i := range nodes {
		nodes[i] = pathNode{cost: unreached, from: -1}
	}

	// ① 擴散。原版是環狀佇列 ＋ 逐層推進；這裡用同樣的「先進先出、
	// 成本較低才覆蓋」語意，結果一致。
	start := idx(from.X, from.Y)
	goal := idx(to.X, to.Y)
	nodes[start].cost = 1
	queue := []int{start}
	found := false
	for len(queue) > 0 && !found {
		cur := queue[0]
		queue = queue[1:]
		cx, cy := cur%Width, cur/Width
		cz := f.StandLevel(cx, cy)
		for _, d := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			nx, ny := cx+d[0], cy+d[1]
			if !inBounds(nx, ny) {
				continue
			}
			if !f.stepOK(cz, nx, ny, climb) {
				continue
			}
			c := cost(nx, ny)
			if c <= 0 {
				continue // 0 ＝ 不可通行
			}
			n := idx(nx, ny)
			if v := nodes[cur].cost + c; v < nodes[n].cost {
				nodes[n] = pathNode{cost: v, from: cur}
				if n == goal {
					found = true
					break
				}
				queue = append(queue, n)
			}
		}
	}
	if !found {
		return nil
	}

	// ② 回溯。從終點往回走，**只有轉彎才記一個點**。
	var back []Point
	cur := goal
	lastDX, lastDY := 0, 0
	for cur != start {
		prev := nodes[cur].from
		if prev < 0 {
			break
		}
		dx := cur%Width - prev%Width
		dy := cur/Width - prev/Width
		// 方向變了 → 前一格是轉彎點。
		if (dx != lastDX || dy != lastDY) && lastDX|lastDY != 0 {
			back = append(back, Point{X: cur%Width, Y: cur / Width})
		}
		lastDX, lastDY = dx, dy
		cur = prev
	}
	// 終點一定要在清單裡（原版最後一個 stosw 就是它）。
	out := make([]Point, 0, len(back)+1)
	for i := len(back) - 1; i >= 0; i-- {
		out = append(out, back[i])
	}
	out = append(out, to)
	if len(out) > MaxWaypoints {
		out = out[:MaxWaypoints]
	}
	return out
}

// stepOK 回報從高度 fromZ 的格子踏進 (x, y) 行不行。
//
// 與 tryMove 同一組規則：一次只能上下一層，爬不上去的兵不能上牆
// （`cmp byte ptr [si+4], 12h / jbe`，docs/re/11 §5.8j）。
func (f *Field) stepOK(fromZ, x, y int, climb bool) bool {
	z := f.StandLevel(x, y)
	if z > fromZ && !climb {
		return false
	}
	return abs(z-fromZ) <= 1
}

// Waypoints 讓一個兵沿著算好的路徑走。
//
// 對應原版：`sub_1AED2` 算好之後把點數寫進 `[si+0x17]`、游標歸零，
// `sub_1B00D` 每次取一個點（`[si+0x16] += 2`）當下一個中繼點。
type Waypoints struct {
	pts []Point
	i   int
}

// Next 取下一個中繼點。取完回 false。
func (w *Waypoints) Next() (Point, bool) {
	if w == nil || w.i >= len(w.pts) {
		return Point{}, false
	}
	p := w.pts[w.i]
	w.i++
	return p, true
}

// Len 回傳還剩幾個點。
func (w *Waypoints) Len() int {
	if w == nil {
		return 0
	}
	return len(w.pts) - w.i
}
