package tactical

// 尋路：波前擴散 ＋ 回溯，出處是 `loc_1BD46`（docs/re/11 §5.15）。
//
// 原版的流程分兩段：
//
//  1. **擴散**：把成本圖（`ds:0D300`，16 KB）全部填成 `0xFFFF`，
//     從起點放 1，用一個環狀佇列往外長。每走一格加的**不是 1，
//     而是那一格的佔用成本**（`mov al, es:[bx+2000h]`）——
//     所以這是 uniform-cost 搜尋，不是單純的 BFS。
//  2. **回溯**：從終點沿著遞減的波數走回去，
//     **只有轉彎時才寫一個點**（`ch` 記著上一步是橫的還是縱的）。
//
// ⭐ 回溯的計數器是 `mov cl, 40h` ＝ **64 個點**，一個點 2 byte
// ＝ 128 byte，正好是每個兵那塊 `0x1800 + 兵編號 × 128` 的大小（§5.8k）。
// 兩邊對得上，這條路徑才走得完。
//
// ⭐ **那張表是「有沒有兵站著」，不是地形。**
//
// `sub_1BBA6` 先把它歸零，之後由移動落定的 `sub_1B240` 維護：
// 舊格寫 0、新格寫 8（透過 `word_1D2FA` 加 `0x9000` 的位移，
// 也就是 `word_1D2FE` 的開頭，docs/re/63 §1）。
// 所以尋路會**繞開有兵的格子**，繞路的代價是 8 格。

// MaxWaypoints 是一條路徑最多幾個轉彎點（原版 `mov cl, 40h`）。
const MaxWaypoints = 64

// Point 是戰場上的一格。
type Point struct{ X, Y int }

// Penalty 是走進 (x, y) 的**額外**成本（0 ＝ 沒有額外成本）。
//
// 原版從 `es:[bx+2000h]` 讀一個 byte 加上去，值是「這一格有兵 ＝ 8」（見上）。
// **能不能走是導航位元決定的，不是這裡**——這裡只影響繞不繞路。
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
	return f.FindPathForcing(from, to, climb, penalty, nil)
}

// FindPathForcing 與 FindPath 相同，但把 force(x, y) 為真的格子**當成可通行**，
// 即使地形擋著。
//
// 它存在的理由只有一個：攻城時城體多半是打不壞的地形，可破壞的只有
// gateX 那一小段城壁與 2–14 格門。純地形尋路在城封死時回空，兵就會
// 整團停在城牆前不動，直到攻城計時器把大將耗光——**攻方永遠攻不進去**。
//
// ⚠ **這是 remake 的近似。** 原版走的是 `0x1800 + 兵編號 × 128` 那張
// 預先算好的繞路點清單，而那個演算法還沒解出來
// （docs/mechanics/30-combat.md「還沒解的」）。解出來之後這一支要換掉。
func (f *Field) FindPathForcing(from, to Point, climb bool, penalty Penalty,
	force func(x, y int) bool) []Point {
	if !inBounds(from.X, from.Y) || !inBounds(to.X, to.Y) {
		return nil
	}
	if from == to {
		return nil
	}
	if penalty == nil {
		penalty = func(int, int) int { return 0 }
	}

	// 節點是 (格, 平面)。原版就是這樣走的：位址的 bit 12 是平面，
	// 擴散時 `test al, 08h` 那一支（`loc_1BFBF`）用 `xor bh, 20h`
	// 換到另一個平面，成本是兩邊地面層的差（docs/re/63 §5）。
	const cells = Width * Height
	idx := func(x, y, plane int) int { return plane*cells + y*Width + x }
	nodes := make([]pathNode, numPlanes*cells)
	for i := range nodes {
		nodes[i] = pathNode{cost: unreached, from: -1}
	}

	// ① 擴散。原版是環狀佇列 ＋ 逐層推進；這裡用同樣的「先進先出、
	// 成本較低才覆蓋」語意，結果一致。
	startPlane := f.planeAt(from.X, from.Y, PlaneLow)
	start := idx(from.X, from.Y, startPlane)
	nodes[start].cost = 1
	queue := []int{start}
	goal := -1
	for len(queue) > 0 && goal < 0 {
		cur := queue[0]
		queue = queue[1:]
		plane := cur / cells
		cx, cy := cur%cells%Width, cur%cells/Width
		push := func(nx, ny, np, extra int) bool {
			n := idx(nx, ny, np)
			v := nodes[cur].cost + extra + penalty(nx, ny)
			if v >= nodes[n].cost {
				return false
			}
			nodes[n] = pathNode{cost: v, from: cur}
			if nx == to.X && ny == to.Y {
				goal = n
				return true
			}
			queue = append(queue, n)
			return false
		}
		for _, d := range navSteps {
			nx, ny := cx+d.dx, cy+d.dy
			if !inBounds(nx, ny) {
				continue
			}
			// 導航位元就是原版的通行圖（`0001BE10` 的四個 test）。
			// force 是 remake 的補丁：可破壞的城壁當成走得過去。
			if !f.Linked(cx, cy, plane, d.dx, d.dy) &&
				!(plane == PlaneLow && f.breachLinked(cx, cy, nx, ny, force)) {
				continue
			}
			if push(nx, ny, plane, 1) {
				break
			}
		}
		if goal >= 0 {
			break
		}
		// 換平面：只有門格可以，而且只有爬得上去的兵種
		// （`sub_1AF69` 的 `cmp [si+4], 12h`，docs/re/63 §4）。
		if !climb || !f.IsGateCell(cx, cy) {
			continue
		}
		other := PlaneHigh
		if plane == PlaneHigh {
			other = PlaneLow
		}
		zHere, okHere := f.GroundLevel(cx, cy, plane)
		zThere, okThere := f.GroundLevel(cx, cy, other)
		if !okHere || !okThere {
			continue
		}
		push(cx, cy, other, abs(zThere-zHere))
	}
	if goal < 0 {
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
		dx := cur%cells%Width - prev%cells%Width
		dy := cur%cells/Width - prev%cells/Width
		if dx == 0 && dy == 0 {
			// 換平面那一步：位置沒變，一定要留一個點，
			// 不然兵走到門那一格就會直接跳過爬牆。
			back = append(back, Point{X: cur % cells % Width, Y: cur % cells / Width})
			lastDX, lastDY = 0, 0
			cur = prev
			continue
		}
		// 方向變了 → 前一格是轉彎點。
		if (dx != lastDX || dy != lastDY) && lastDX|lastDY != 0 {
			back = append(back, Point{X: cur % cells % Width, Y: cur % cells / Width})
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

// planeAt 挑一個位置該用哪個平面：優先用 want，那個平面沒有地面就換另一個。
func (f *Field) planeAt(x, y, want int) int {
	if _, ok := f.GroundLevel(x, y, want); ok {
		return want
	}
	other := PlaneHigh
	if want == PlaneHigh {
		other = PlaneLow
	}
	if _, ok := f.GroundLevel(x, y, other); ok {
		return other
	}
	return want
}

// breachLinked 是 remake 的補丁：把「可以撞穿的城壁」當成走得過去。
//
// 導航位元是照**沒打壞**的地形算的，所以城壁那一格與兩邊都不連通。
// 撞穿之後那一格會變瓦礫（高度 0），這裡就照打壞後的高度判連通——
// 兵才會走到牆前撞上去，`tryMove` 把它算成一次耐久損傷。
//
// ⚠ 原版不需要這一段：它的攻方走的是門那一格爬上牆頂（docs/re/63 §4）。
// 這條路徑在 remake 也接上了，這個補丁只剩「連門都沒有的圖」在用。
func (f *Field) breachLinked(cx, cy, nx, ny int, force func(x, y int) bool) bool {
	if force == nil {
		return false
	}
	a, aok := f.breachLevel(cx, cy, force)
	b, bok := f.breachLevel(nx, ny, force)
	return aok && bok && abs(a-b) <= 1
}

// breachLevel 回傳「撞穿之後」那一格的可站立層。
func (f *Field) breachLevel(x, y int, force func(x, y int) bool) (int, bool) {
	if force != nil && force(x, y) {
		return 0, true
	}
	return f.GroundLevel(x, y, PlaneLow)
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
