// Package army 實作軍團與行軍。
//
// 結構出自 `KI.EXE` 的 `sub_16F26`（編成）、`sub_16F86`（位置）、
// `sub_12662`／`sub_127A2`（行軍），反組譯見 docs/re/08 §4–§6。
// 機制說明見 docs/mechanics/10-strategy.md §3–§4。
//
// ⚠ 連結表的實際格式還沒定案（docs/re/08 §6 的警告），
// 所以這一層只實作**已經確定**的部分：軍團的組成、佔用圖、
// 到達判定與節點分類。實際的路徑推進留給呼叫端餵節點。
package army

// 一個軍團最多六個部隊，每個部隊固定 1000 人（說明書 5.5）。
const (
	Positions      = 6
	MenPerUnit     = 1000
	MaxStrength    = Positions * MenPerUnit // 6000
	DefaultMorale  = 200                    // 勢力記錄 +0x1D
	RoutMoraleGate = 100                    // 士氣低於此值戰敗 → 軍團壞滅
)

// Position 是六個編成位置（說明書 5.5 的編成畫面，400 dpi 判讀）。
type Position int

const (
	General    Position = iota // 大將
	Vanguard                   // 先鋒
	LeftWing                   // 左翼
	RightWing                  // 右翼
	LeftGuard                  // 左備
	RightGuard                 // 右備
)

func (p Position) String() string {
	return [...]string{"大將", "先鋒", "左翼", "右翼", "左備", "右備"}[p]
}

// TroopType 是三個兵種。與 economy 的順序一致（勢力記錄 +0x04/+0x06/+0x08）。
type TroopType int

const (
	Cavalry  TroopType = iota // 騎馬
	Archer                    // 弓兵
	Infantry                  // 步兵
)

func (t TroopType) String() string {
	return [...]string{"騎馬", "弓兵", "步兵"}[t]
}

// 節點索引的兩個邊界（原版比的是索引 × 8，這裡用索引本身）。
//
//	0 – 191   據點（與據點表一一對應）
//	192 – 255 路上的中間節點
//	≥ 256     野外
const (
	NumCityNodes = 192 // 原版 cmp bx, 600h → 索引 × 8 < 0x600
	NumRoadNodes = 256 // 原版 cmp bx, 800h
)

// NodeKind 是節點的三種類型。
type NodeKind int

const (
	CityNode NodeKind = iota
	RoadNode
	FieldNode
)

// KindOf 依節點索引分類。
func KindOf(node int) NodeKind {
	switch {
	case node < NumCityNodes:
		return CityNode
	case node < NumRoadNodes:
		return RoadNode
	default:
		return FieldNode
	}
}

// Corps 是一個軍團。欄位對應軍團記錄（docs/re/08 §4）。
type Corps struct {
	Alive   bool
	Faction int
	Leader  int // 武將編號

	// Units 是六個編成位置各放什麼兵種。空位用 -1。
	Units  [Positions]TroopType
	Manned [Positions]bool

	Morale int // 記錄 +0x06，編成時 = 勢力的士氣基準

	// 現在與目標是**兩組獨立的三元組**。行軍就是把前者推向後者。
	Node, X, Y                   int // 記錄 +0x0E（÷8）、+0x10、+0x12
	TargetNode, TargetX, TargetY int // 記錄 +0x14（÷8）、+0x16、+0x18

	Direction int // 記錄 +0x0A，決定走哪一條連結

	// MoveTimer 與 MoveInterval 是行軍的節拍（記錄 +0x0B／+0x1E）。
	// 每個 tick 計時器減 1，歸零就走一步並從 MoveInterval 重載。
	//
	// **所以「速度」是間隔的倒數**：間隔越小走得越快。
	// 說明書 5.5「騎馬隊のみの軍団は移動速度が速くなります」
	// 就是給純騎馬編成一個較小的間隔。
	MoveTimer    int
	MoveInterval int
}

// Step 把行軍的計時器推進一個 tick，回傳這個 tick 是否要走一步。
//
// 照原版 `sub_125A3`：先減，歸零才動作並重載。
// **注意是「先減再判斷」**——間隔設 N 表示每 N 個 tick 走一步，
// 不是 N+1。
func (c *Corps) Step() bool {
	c.MoveTimer--
	if c.MoveTimer > 0 {
		return false
	}
	c.MoveTimer = c.MoveInterval
	return true
}

// Strength 回傳總兵力。每個有兵的位置固定 1000 人（說明書 5.5）。
func (c Corps) Strength() int {
	n := 0
	for _, ok := range c.Manned {
		if ok {
			n++
		}
	}
	return n * MenPerUnit
}

// Arrived 回報軍團是否已經到達目標節點。
// 原版就是直接比兩個節點編號（`cmp bx, [si+14h]`），不比座標。
func (c Corps) Arrived() bool { return c.Node == c.TargetNode }

// CanScaleWalls 回報這個軍團爬不爬得上城牆。
//
// 說明書 5.5：「騎馬のみの編成では城壁に登れないため
// 攻城戦では弓兵、歩兵が必要です」。
func (c Corps) CanScaleWalls() bool {
	for i, ok := range c.Manned {
		if ok && c.Units[i] != Cavalry {
			return true
		}
	}
	return false
}

// AllCavalry 回報是不是純騎馬編成。原版給它移動速度加成
// （說明書 5.5：「騎馬隊のみの軍団は移動速度が速くなります」）。
func (c Corps) AllCavalry() bool {
	any := false
	for i, ok := range c.Manned {
		if !ok {
			continue
		}
		any = true
		if c.Units[i] != Cavalry {
			return false
		}
	}
	return any
}

// Destroyed 回報這次戰敗會不會讓軍團壞滅。
//
// 說明書 5.5：「100 を切った状態で戦闘して負けると軍団が壊滅し、
// 場合によっては武将が敵に捕まります」。**100 是說明書直接給的門檻。**
func (c Corps) Destroyed(lostBattle bool) bool {
	return lostBattle && c.Morale < RoutMoraleGate
}

// OccupancyMap 是 384×256 的單位計數圖（docs/re/08 §5）。
//
// 原版把它放在一塊 98,304 B 的區域，位址是 `Y × 384 + X`，
// 軍團進入時 `inc`、離開時 `dec`。**每格是計數不是布林**——
// 同一格可以疊多個軍團。
type OccupancyMap struct {
	w, h  int
	cells []uint8
}

// NewOccupancyMap 建一張與大地圖同尺寸的佔用圖。
func NewOccupancyMap() *OccupancyMap {
	const w, h = 384, 256
	return &OccupancyMap{w: w, h: h, cells: make([]uint8, w*h)}
}

func (m *OccupancyMap) idx(x, y int) (int, bool) {
	if x < 0 || x >= m.w || y < 0 || y >= m.h {
		return 0, false
	}
	return y*m.w + x, true
}

// Enter 讓一個單位進入某格。
func (m *OccupancyMap) Enter(x, y int) {
	if i, ok := m.idx(x, y); ok && m.cells[i] < 255 {
		m.cells[i]++
	}
}

// Leave 讓一個單位離開某格。
//
// 原版是無條件 `dec`，會在 0 的時候繞回 255。這裡夾在 0——
// **這是刻意的 remake 差異**：原版那個下溢只有在配對錯誤時才會發生，
// 而繞回 255 會讓那一格永遠看起來有兵，比直接夾住難查得多。
func (m *OccupancyMap) Leave(x, y int) {
	if i, ok := m.idx(x, y); ok && m.cells[i] > 0 {
		m.cells[i]--
	}
}

// At 回傳某格上的單位數。
func (m *OccupancyMap) At(x, y int) int {
	if i, ok := m.idx(x, y); ok {
		return int(m.cells[i])
	}
	return 0
}

// Occupied 回報某格上有沒有單位。
func (m *OccupancyMap) Occupied(x, y int) bool { return m.At(x, y) > 0 }
