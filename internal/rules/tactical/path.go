package tactical

// 尋路：波前擴散 ＋ 回溯，出處是 `loc_1BD46`（docs/re/11 §5.15）。
//
// 原版的流程分兩段：
//
//  1. **擴散**：把成本圖（`ds:0D300`，16 KB）全部填成 `0xFFFF`，
//     從目標放 1，用一個環狀佇列往士兵位置長。每走一格加的**不是 1，
//     而是那一格的佔用成本**（`mov al, es:[bx+2000h]`）——
//     所以這是 uniform-cost 搜尋，不是單純的 BFS。
//  2. **回溯**：從士兵位置依成本及軸順序取點，不沿擴散前驅，
//     **只有轉彎時才寫一個點**（`ch` 記著上一步是橫的還是縱的）。
//
// 正式單兵明示終點平面；goalPlane=-1的舊通用helper保留前向相容語意。
// 證據與範圍見spec/134的2026-11-12訂正。
//
// ⭐ 回溯的計數器是 `mov cl, 40h` ＝ **64 個點**，一個點 2 byte
// ＝ 128 byte，正好是每個兵那塊 `0x1800 + 兵編號 × 128` 的大小（§5.8k）。
// 兩邊對得上，這條路徑才走得完。
//
// 額外成本包含單兵8、城壁100與門50，分低／高兩個平面。
// 建場、逐兵重畫／移除及破城維護同一張表，不由目前兵座標重建（spec/134）。

// MaxWaypoints 是一條路徑最多幾個轉彎點（原版 `mov cl, 40h`）。
const MaxWaypoints = 64

// Point 是戰場上的一格，或一個只改高度的繞路點。
//
// 原版的 128-byte 路徑緩衝區不是每一筆都代表 XY：`sub_1B00D`
// 看到第一 byte >= 0x80 時，只把第二 byte 寫進 StepZ，XY 保持不變
// （docs/spec/159）。HasZ 用來保存這個非破壞性差異；一般 XY 點維持
// HasZ=false，避免把普通路徑誤當成高度切換。
type Point struct {
	X, Y int
	Z    int
	HasZ bool
}

// Penalty 是走進 (x, y) 的**額外**成本（0 ＝ 沒有額外成本）。
//
// 原版從 `es:[bx+2000h]` 讀一個 byte 加上去，值是「這一格有兵 ＝ 8」（見上）。
// **能不能走是導航位元決定的，不是這裡**——這裡只影響繞不繞路。
type Penalty func(x, y int) int

// planePenalty 是正式單兵使用的兩平面成本，舊XY介面由包裝器保留。
type planePenalty func(x, y, plane int) int

// pathNode 是擴散時每一格的狀態。
type pathNode struct {
	cost int // 0xFFFF ＝ 還沒走到
	from int // 舊通用helper的前驅；正式單兵依成本回溯，不讀這一欄
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
	return f.FindPathForPlanes(from, to, f.planeAt(from.X, from.Y, PlaneLow), -1, climb, penalty)
}

// FindPathForPlanes 與 FindPath 相同，但明確指定起點／終點平面。
//
// 舊呼叫端只給 XY 時，原版會從低平面開始、在 XY 目的地先到就停止；
// 戰術單位則可能已在牆頂，或明確要走到敵人的高平面，因此正式移動路徑
// 不能再把目前平面遺失。goalPlane=-1 保留「任一平面到達 XY 即可」的
// 舊語意。
func (f *Field) FindPathForPlanes(from, to Point, startPlane, goalPlane int,
	climb bool, penalty Penalty) []Point {
	return f.findPathForPlanes(from, to, startPlane, goalPlane, climb, penalty, nil)
}

// FindPathForcing 與 FindPath 相同，但把 force(x, y) 為真的格子**當成可通行**，
// 即使地形擋著。
//
// 它存在的理由只有一個：攻城時城體多半是打不壞的地形，可破壞的只有
// gateX 那一小段城壁與 2–14 格門。純地形尋路在城封死時回空，兵就會
// 整團停在城牆前不動，直到攻城計時器把大將耗光——**攻方永遠攻不進去**。
//
// ⚠ **`force` 這個參數是 remake 加的，原版沒有。** 正式明示平面的單兵
// 路徑採loc_1BD46的反向擴散／成本回溯（spec/134）；原版把點寫進
// `0x1800 + 兵編號 × 128`，remake存進Soldier.Path。
// 此通用helper未明示終點平面，保留既有前向相容語意，不作逐點parity聲明。
//
// 加 `force` 的理由是**驗收路徑要能終止**：AI 對 AI 的自動戰鬥沒有玩家
// 換陣形，城封死時純地形尋路回空，兵會整團停在牆前直到攻城計時器把大將
// 耗光。⭐ **原版不需要這一條，因為破城本來就是玩家的工作**
// （說明書第 11 章，docs/spec/37 §4）——所以這是 remake 差異，
// 不是待補的缺口。
func (f *Field) FindPathForcing(from, to Point, climb bool, penalty Penalty,
	force func(x, y int) bool) []Point {
	return f.FindPathForPlanesForcing(from, to, f.planeAt(from.X, from.Y, PlaneLow), -1,
		climb, penalty, force)
}

// FindPathForPlanesForcing 是 FindPathForPlanes 的攻城 fallback 版本。
func (f *Field) FindPathForPlanesForcing(from, to Point, startPlane, goalPlane int,
	climb bool, penalty Penalty, force func(x, y int) bool) []Point {
	return f.findPathForPlanes(from, to, startPlane, goalPlane, climb, penalty, force)
}

func (f *Field) findPathForPlanes(from, to Point, startPlane, goalPlane int,
	climb bool, penalty Penalty, force func(x, y int) bool) []Point {
	var cost planePenalty
	if penalty != nil {
		cost = func(x, y, _ int) int { return penalty(x, y) }
	}
	return f.findPathForPlaneCosts(from, to, startPlane, goalPlane, climb, cost, force)
}

func (f *Field) findPathForPlaneCosts(from, to Point, startPlane, goalPlane int,
	climb bool, penalty planePenalty, force func(x, y int) bool) []Point {
	if !inBounds(from.X, from.Y) || !inBounds(to.X, to.Y) {
		return nil
	}
	if startPlane < 0 || startPlane >= numPlanes ||
		(goalPlane < -1 || goalPlane >= numPlanes) {
		return nil
	}
	if from.X == to.X && from.Y == to.Y && (goalPlane < 0 || startPlane == goalPlane) {
		return nil
	}
	if penalty == nil {
		penalty = func(int, int, int) int { return 0 }
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

	// ① 擴散。⭐ **這是波數佇列（dial queue）版的 Dijkstra，不是 BFS**
	// （`docs/spec/134`）：佇列之外另外帶一個**波數**，彈出來的格子
	// 成本還沒被波數追過就**推回佇列**，這一波不展開；一波消化完
	// 波數才 +1。有兵的格子成本 +8，等於在佇列裡多躺八波。
	//
	// ⚠ **少了「推回」那一步，`penalty` 就完全不起作用**——純先進先出
	// 之下直線一定比繞路先碰到終點，於是被大將擋住（大將不能對調，
	// §5.16）的兵會永久卡死在原地。
	reverse := goalPlane >= 0
	searchFrom, searchTo := from, to
	rootPlane, stopPlane := startPlane, goalPlane
	if reverse {
		searchFrom, searchTo = to, from
		rootPlane, stopPlane = goalPlane, startPlane
	}
	start := idx(searchFrom.X, searchFrom.Y, rootPlane)
	nodes[start].cost = 1

	// 原版的環狀佇列：`si` 寫、`di` 讀、`cx` 記本波的結尾。
	// 起點不進佇列——`loc_1BDD3` 設好成本就直接跳進展開那一段。
	queue := []int{}
	read, waveEnd := 0, 0
	wave := 2 // 原版的 `mov dx, 2`
	cur := start
	goal := -1

	for {
		// `cmp bx, cs:word_1BD44 / jz loc_1BE4A`：彈到終點就收工。
		if cur%cells%Width == searchTo.X && cur%cells/Width == searchTo.Y &&
			(stopPlane < 0 || cur/cells == stopPlane) {
			goal = cur
			break
		}
		// `cmp dx, [bx] / ja`：波數還沒追過這一格 → 推回佇列。
		if wave <= nodes[cur].cost {
			queue = append(queue, cur)
		} else {
			plane := cur / cells
			cx, cy := cur%cells%Width, cur%cells/Width
			push := func(nx, ny, np, extra int) {
				n := idx(nx, ny, np)
				// `cmp dx, [bx-2] / jnb ret`
				if wave >= nodes[n].cost {
					return
				}
				// `add ax, dx`：新成本 ＝ **波數** ＋ 那一格的成本。
				// 平面切換的高度差走 extra（原版是另一支的加法）。
				nodes[n] = pathNode{cost: wave + extra + penalty(nx, ny, np), from: cur}
				queue = append(queue, n)
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
				push(nx, ny, plane, 0)
			}
			// 換平面：只有門格可以，而且只有爬得上去的兵種
			// （`sub_1AF69` 的 `cmp [si+4], 12h`，docs/re/63 §4）。
			if climb && f.IsGateCell(cx, cy) {
				other := PlaneHigh
				if plane == PlaneHigh {
					other = PlaneLow
				}
				zHere, okHere := f.GroundLevel(cx, cy, plane)
				zThere, okThere := f.GroundLevel(cx, cy, other)
				if okHere && okThere {
					push(cx, cy, other, abs(zThere-zHere))
				}
			}
		}
		// `loc_1BE38`：本波消化完就推波數，佇列空了才算走不到。
		if read == waveEnd {
			wave++
			waveEnd = len(queue)
			if read == waveEnd {
				break
			}
		}
		cur = queue[read]
		read++
	}
	if goal < 0 {
		return nil
	}
	if reverse {
		// 1BE4A：用停止時波數減2，停止格自己的佔用成本可能還沒到期。
		return f.originalPathPoints(nodes, goal, wave-2)
	}

	// ② 回溯。從終點往回走，**只有轉彎才記一個點**。
	var back []Point
	cur = goal
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
			// 不然兵走到門那一格就會直接跳過爬牆。這不是普通 XY
			// 點：`sub_1B00D` 的高位標記只更新 StepZ。
			x, y := cur%cells%Width, cur%cells/Width
			z, ok := f.GroundLevel(x, y, cur/cells)
			if !ok {
				// 正常的平面邊界只發生在門格；若資料不完整，保留
				// 座標點也比捏造高度安全，並讓後續驗證看見缺口。
				back = append(back, Point{X: x, Y: y})
			} else {
				back = append(back, Point{X: x, Y: y, Z: z, HasZ: true})
			}
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

// originalPathPoints 對應1BE4A–1BF23；只讀成本，不追擴散前驅。
func (f *Field) originalPathPoints(nodes []pathNode, cur, remaining int) []Point {
	const cells = Width * Height
	points := make([]Point, 0, 8)
	appendXY := func() bool {
		points = append(points, Point{X: cur % cells % Width, Y: cur % cells / Width})
		return len(points) >= MaxWaypoints
	}
	axis, first := 0, true
	for remaining > 0 {
		plane, x, y := cur/cells, cur%cells%Width, cur%cells/Width
		move := func(a int) bool {
			for _, delta := range []int{-1, 1} {
				nx, ny := x, y
				if a == 0 {
					nx += delta
				} else {
					ny += delta
				}
				if !inBounds(nx, ny) {
					continue
				}
				n := plane*cells + ny*Width + nx
				if nodes[n].cost == remaining {
					cur = n
					return true
				}
			}
			return false
		}
		moved := move(axis)
		if !moved {
			axis = 1 - axis
			// 第一次選方向不寫點，其後換軸才寫當前XY。
			if !first && appendXY() {
				return points
			}
			moved = move(axis)
		}
		if moved {
			remaining--
			if remaining == 0 {
				appendXY()
				return points
			}
			first = false
			continue
		}
		// 1BEE0：門格可換平面；其餘情況遞減預算，不猜一個前驅。
		other := (1-plane)*cells + y*Width + x
		if f.nav[plane][y][x]&navGate != 0 && nodes[other].cost <= remaining {
			z := 0
			if 1-plane == PlaneHigh {
				if v, ok := f.GroundLevel(x, y, PlaneHigh); ok {
					z = v
				}
			}
			points = append(points, Point{X: x, Y: y, Z: z, HasZ: true})
			if len(points) >= MaxWaypoints {
				return points
			}
			cur, remaining = other, nodes[other].cost
		}
		remaining--
		axis, first = 0, true
	}
	return points
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

// Last 回傳這條路的終點。沒有路（或已經走完）時回 false。
//
// 用途是分辨「這條路是為了哪個目標算的」——`doRetreat` 靠它判斷
// 手上這條是不是退卻用的繞路（docs/spec/94 §3）。
func (w *Waypoints) Last() (Point, bool) {
	if w == nil || len(w.pts) == 0 {
		return Point{}, false
	}
	return w.pts[len(w.pts)-1], true
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
