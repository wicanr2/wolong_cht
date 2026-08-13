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
// ⭐ **那張地形成本表在出貨版裡永遠是 0。**
//
// `sub_1BBA6` 最後 `mov es, cs:word_1D2FE / xor ax, ax / mov cx, 1000h /
// rep stosw` 把它全部歸零，而**整支程式再也沒有人寫它**——
// `word_1D2FE` 這個符號總共只出現兩次（資料定義 ＋ 那一次歸零），
// 從其他基底也搆不到（`+0x2000`／`+0x3000` 那些寫入的 `ds`／`es`
// 都是 `word_1D2FA` 的七層通行圖）。
//
// 所以 `[bx] = 地形成本 + dx = dx`，**實際行為就是純 BFS**。
// 這個欄位是留著沒用的設計。

// MaxWaypoints 是一條路徑最多幾個轉彎點（原版 `mov cl, 40h`）。
const MaxWaypoints = 64

// Point 是戰場上的一格。
type Point struct{ X, Y int }

// Penalty 是走進 (x, y) 的**額外**成本（0 ＝ 沒有額外成本）。
//
// 原版從 `es:[bx+2000h]` 讀一個 byte 加上去，但那張表永遠是 0（見上）。
// 留這個鉤子是為了「解出誰該寫它」時不必動演算法；
// **能不能走是地形決定的，不是這裡**。
type Penalty func(x, y int) int

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
func (f *Field) FindPath(from, to Point, climb bool, penalty Penalty) []Point {
	if !inBounds(from.X, from.Y) || !inBounds(to.X, to.Y) {
		return nil
	}
	if from == to {
		return nil
	}
	if penalty == nil {
		penalty = func(int, int) int { return 0 }
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
			n := idx(nx, ny)
			// 原版是 `[bx] = 地形成本 + dx`（dx 是波數）；成本恆為 0
			// 的話就等於每走一格加 1。
			if v := nodes[cur].cost + 1 + penalty(nx, ny); v < nodes[n].cost {
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
			back = append(back, Point{X: cur % Width, Y: cur / Width})
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
// 與 tryMove 同一組規則：水平跨格允許同步一層高度；只有在同一格做
// 純 Z 軸移動時，才由兵種能力限制爬牆（`cmp byte ptr [si+4], 12h / jbe`，
// docs/re/11 §5.8j）。
func (f *Field) stepOK(fromZ, x, y int, climb bool) bool {
	// 尋路的自我修改分支仍依 `climb` 決定能否向上跨越高度；這讓非爬牆
	// 兵繞過整道一層牆。實際逐格移動則另由 tryMove 重現
	// `sub_1B1B1` 的水平一層高度同步，兩者不能混成同一個檢查。
	z := f.StandLevel(x, y)
	if abs(z-fromZ) > 1 {
		return false
	}
	return z <= fromZ || climb
}

// Waypoints 讓一個兵沿著算好的路徑走。
//
// 對應原版：`sub_1AED2` 算好之後把點數寫進 `[si+0x17]`、游標歸零，
// `sub_1B00D` 每次取一個點（`[si+0x16] += 2`）當下一個中繼點。
type Waypoints struct {
	pts []Point
	i   int
}

// Current 回傳目前的中繼點，但不前進游標。
// 原版 `sub_1B00D` 只有在兵已抵達目前目標後才會遞減點數並前進
// `+0x16`；每幀先取下一個點會把尚未走完的轉角跳掉。
func (w *Waypoints) Current() (Point, bool) {
	if w == nil || w.i >= len(w.pts) {
		return Point{}, false
	}
	return w.pts[w.i], true
}

// Advance 抵達目前中繼點後才消費它。
func (w *Waypoints) Advance() {
	if w == nil || w.i >= len(w.pts) {
		return
	}
	w.i++
}

// Len 回傳還剩幾個點。
func (w *Waypoints) Len() int {
	if w == nil {
		return 0
	}
	return len(w.pts) - w.i
}
