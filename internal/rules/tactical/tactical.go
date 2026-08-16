// Package tactical 實作戰術戰鬥。
//
// 規格全部出自 `KI.EXE` 的戰術模組，完整反組譯見 docs/re/11，
// 機制說明見 docs/mechanics/30-combat.md。
//
// 原版的戰術戰鬥是**硬即時**的（說明書 4.1：「戦闘中は絶対に時間を
// 止められません」），所以這一層做成「一幀一步」的純函式狀態機：
// 呼叫端每幀呼叫一次 Step，不含畫面、不含輸入。
// 這樣它可以在無頭環境裡大量重複執行來驗證長期行為。
//
// 規則的每一條都對得上反組譯。少數幾個地方是 remake 的取捨
// （重算繞路的節流、退卻中可穿過自己人），都在原處標出來了。
package tactical

// 戰場的尺寸。格索引是三維的：Z × 0x1000 ＋ Y × 0x40 ＋ X
// （`sub_1BA2E`／`sub_1B047`，docs/re/11 §5.2）。
const (
	Width  = 64
	Height = 62
	Levels = 7

	// 座標會被 `sub_1AACF` 夾在 1–62。
	MinCoord = 1
	MaxCoord = 62

	// ⚠ **戰場只有 Height 列**（`BATTLE.MAP` 一張圖 62 列 × 64 行），
	// 所以 Y 的上限比 `sub_1AACF` 的 62 少一格。原版的格子陣列是
	// 64 × 64（每層 0x1000），最後兩列是空的；本專案只存 62 列，
	// 用 MaxCoord 去索引會越界。**X 用 MaxCoord、Y 用 MaxY。**
	MaxY = Height - 1
)

// 一側的編制：6 隊 × 8 人，每隊 1 隊長 ＋ 7 隊員
// （`sub_1A754` 的 1 + 7 迴圈，docs/re/11 §5.6）。
const (
	Squads         = 6
	PerSquad       = 8
	SoldiersOnFoot = Squads * PerSquad // 48
)

// Kind 是兵種。原版存的是**兵種 × 18**（`sub_1AB7C` 的三個分支門檻
// 都是 18 的倍數，docs/re/11 §5.8e）。
type Kind int

const (
	General  Kind = 0  // 大將：不攻擊、不會陣亡、爬不上城牆
	Cavalry  Kind = 18 // 騎馬：追著敵人打，爬不上城牆
	Archer   Kind = 36 // 弓兵：站著不動放箭
	Infantry Kind = 54 // 步兵：近戰，挨箭只吃四分之一
)

func (k Kind) String() string {
	switch k {
	case General:
		return "大將"
	case Cavalry:
		return "騎馬"
	case Archer:
		return "弓兵"
	case Infantry:
		return "步兵"
	}
	return "?"
}

// PlaneHighGround／PlaneHighElevated 是兵記錄 +0x1E 的原始值。
// 原版不是把 Z 整數直接存進這一 byte，而是以 0／0x10 選擇格索引的
// 高位平面（`sub_1B0D3`／`sub_1B116`，IDA `seg000:1B103`／`1B14A`）。
const (
	PlaneHighGround   byte = 0x00
	PlaneHighElevated byte = 0x10
)

// Command 是六個戰術指令（說明書 4.2）。編號與原版一致，
// 每一個都是從程式行為認出來的，見 docs/re/11 §5.8b。
type Command int

const (
	Form     Command = 0 // 陣形：走到陣形指定的座標
	Attack   Command = 1 // 攻擊：大將以外的兵攻擊
	Charge   Command = 2 // 突擊：攻城戰的守方會開門
	ScaleWal Command = 3 // 城壁移動：野戰時自動變成攻擊
	Guard    Command = 4 // 守陣：敵人靠近就打、走遠就回陣
	Retreat  Command = 5 // 退卻：不可打斷

	// 內部狀態，腳本不會下達（docs/re/11 §5.8b）。
	Holding Command = 7 // 陣形已就位
)

func (c Command) String() string {
	switch c {
	case Form:
		return "陣形"
	case Attack:
		return "攻擊"
	case Charge:
		return "突擊"
	case ScaleWal:
		return "城壁移動"
	case Guard:
		return "守陣"
	case Retreat:
		return "退卻"
	case Holding:
		return "就位"
	}
	return "?"
}

// 疲勞度（原版叫「餘力」，兵記錄 +0x19，docs/re/11 §5.8f）。
const (
	// StaminaFull 是走到陣形位置那一刻補到的值。
	// **只有「走到定位」會補**——下令不會（`sub_1AA2C` 的
	// `mov byte ptr [si+19h], 80h`）。
	StaminaFull = 128
	// StaminaFighting 是攻擊時的上限（`sub_1AB7C` 的 `cmp 28h`）。
	StaminaFighting = 40
	// StaminaBackToForm 是守陣時退回陣形的門檻（`sub_1ABB2` 的 `cmp 10h`）。
	StaminaBackToForm = 16
)

const (
	// MaxHP 是兵的體力上限（`sub_1B97E` 的 `cmp [bx+3], 64h`）。
	MaxHP = 100
	// GeneralRetreatHP 是大將自動退卻的門檻（`sub_1AE56` 的 `cmp [di+3], 32h`）。
	GeneralRetreatHP = 50
	// SiegeDrainInterval 是攻城方大將體力遞減的間隔（`sub_1AE56` 的 `mov cs:byte_1D321, 0Ah`）。
	SiegeDrainInterval = 10
)

// 面向。`sub_1B047`／`sub_1B069` 直接寫 0 與 2，飛道具移動
// （`sub_1BA2E`）用同一組編碼。
const (
	West = iota
	North
	East
	South
)

// Soldier 是一個兵。欄位對應原版的 32 byte 記錄（docs/re/11 §5.8l）。
type Soldier struct {
	Alive bool
	Kind  Kind
	HP    int

	X, Y, Z int
	Facing  int

	// PoseStep 是原版兵記錄 +0x02 bit 0。主迴圈每次更新後翻轉它；
	// 除了人物走路圖，也會被特殊投射物拿來選 raw 0x214／0x215。
	// 它是呈現／暫態欄位，不影響戰鬥規則與存檔。
	PoseStep uint8

	// ProjectileCooldown 對應原版兵記錄 +0x13。sub_1AD2D 成功發射
	// 普通投射物時寫入 8；sub_1AD7F 成功發射特殊投射物時寫入 6；
	// 下一次發射嘗試會先遞減非零值並直接返回。
	ProjectileCooldown uint8

	// Stamina 是疲勞度（餘力），越高越好。
	Stamina int

	// PlaneHigh 是原版 +0x1E 的 raw runtime 欄位，只會由地面 0 變成
	// 高位平面 0x10，再回到 0；它不是 Z 高度本身。
	PlaneHigh byte

	// OnWall 是原版的 `[si+0x1E]`（PlaneHigh）：這個兵在不在城牆頂那個
	// 平面上。只有 `sub_1B0D3`／`sub_1B116` 會改它，也就是**在門那一格
	// 上下**（docs/re/63 §4）。
	OnWall bool

	// HighTerrain 是原版 +0x00 bit 1。`sub_1B240` 依目前圖塊堆疊高度
	// 是否至少 4 層設定它，鎖敵與換位都會讀這個旗標。
	HighTerrain bool

	// Climbing 是早期 remake API 留下的相容欄位。新 runtime 以 PlaneHigh
	// 為準，讀取時仍接受手寫測試中的 Climbing=true。
	Climbing bool

	// Cmd 是生效中的命令、Next 是新下達的命令。
	// **兩個欄位是分開的**，所以腳本可以「下令之後等幾幀再問到位了沒」。
	Cmd, Next Command

	// Target 是鎖定的敵人（對方陣營的索引），−1 表示沒有。
	Target int

	// GoalX/GoalY/GoalZ 是最終目標（陣形位置或敵人的位置）。
	// StepX/StepY/StepZ 是這一步要走向的中繼點。
	GoalX, GoalY, GoalZ int
	StepX, StepY, StepZ int

	// Path 是還沒走完的繞路點（原版 `0x1800 + 兵編號 × 128`，§5.15）。
	// PathAt 是上次重算的幀，用來節流。
	Path   *Waypoints
	PathAt int

	// Power 是由士氣算出來的戰力（原版 `+0x18`，`sub_19B6D` 寫）。
	// 打大將時的命中率與傷害都看它。
	Power int

	// HitGeneral 是「這一擊打在大將身上」（原版 `+0x02` 的 bit 3，
	// `sub_1B6BC` 結尾 `or byte ptr [si+2], 8`）——圖號會換一張。
	HitGeneral bool

	// Hurt 是「剛剛被打中」（原版 `+0x02` 的 bit 4）。
	//
	// `sub_1B618` 打中人的時候一起做三件事：設這個位元、**把面向歸零**
	// （`mov byte ptr [di+5], 0`）、設 Swapped。而圖號公式裡
	// bit 4 的作用正是「面向一律當成 0」（§5.13）——**兩邊是同一件事**：
	// 受擊的兵轉向正面、畫受擊的圖，而且那一幀不能被換位置。
	Hurt bool

	// Swapped 是「這一幀已經被別人換過位置了」（原版 `+0x00` 的 bit 6）。
	//
	// `sub_1B732` 換完會對被換的那一個 `or byte ptr [di], 40h`，
	// 而 `sub_1B240` 在自己更新時 `and byte ptr [si], 0BFh` 清掉。
	// **沒有這個旗標，兩個兵會原地互換到天荒地老**——
	// `seg000:B56D` 的 `test byte ptr [di], 61h` 檢查的就是它（bit 6）。
	Swapped bool
}

// IsGeneral 回報這是不是大將。原版用 `+0x04 == 0` 判。
func (s *Soldier) IsGeneral() bool { return s.Kind == General }

// planeHigh 回傳原版 +0x1E；Climbing 只保留給未經 Place 的舊測試資料。
func (s *Soldier) planeHigh() byte {
	if s.PlaneHigh != PlaneHighGround {
		return s.PlaneHigh
	}
	if s.Climbing {
		return PlaneHighElevated
	}
	return PlaneHighGround
}

// syncTerrain 把目前座標同步成原版的 +0x1E 與 +0x00 bit 1。
// Plane 回傳這個兵所在的平面（PlaneLow／PlaneHigh）。
func (s *Soldier) Plane() int {
	if s.OnWall {
		return PlaneHigh
	}
	return PlaneLow
}

func (s *Soldier) syncTerrain(f *Field, x, y, z int) {
	s.PlaneHigh = PlaneHighGround
	s.Climbing = false
	if z > 0 {
		s.PlaneHigh = PlaneHighElevated
		s.Climbing = true
	}
	s.HighTerrain = f != nil && f.HighTerrain(x, y)
}

// CanClimb 回報這個兵爬不爬得上城牆。
//
// 原版是 `cmp byte ptr [si+4], 12h / jbe` ——**大將與騎馬跳過 Z 軸移動**
// （docs/re/11 §5.8j）。說明書 5.5「騎馬のみの編成では城壁に登れない」
// 在機器碼裡就是這一行。
func (s *Soldier) CanClimb() bool { return s.Kind > Cavalry }

// Field 是戰場的立體格。Solid[z] 為真表示那一層是實心的（站不上去）。
//
// 原版把它存成 7 張 64×64 的圖（`sub_1BB6D` 每層相隔 0x1000），
// 一格的可站立層由地圖圖塊的堆疊決定（docs/re/11 §4.2、§5.2）。
type Field struct {
	solid [Levels][Height][Width]bool
	// gateX 是命令 3（城壁移動）要走過去的那一格 X
	// （`BATTLE.MAP` 索引的第二欄，docs/re/11 §5.8i）。0 表示這張圖沒有城。
	gateX int
	// top[y][x] 是那一格最高的可站立層。
	top [Height][Width]int

	// tiles[y][x] 是原始圖塊值、heights[圖塊] 是它的堆疊層數。
	// 兩個都只有「從圖塊建的戰場」才有——城壁打壞時要換圖塊再重算高度，
	// 沒有它們就退化成「打壞了但地形不變」。
	tiles   [][]byte
	heights *[256]int
	// layers[圖塊] 是那個圖塊由下往上七層的子圖塊（`BATTLE.MDL` 的 [1..7]）。
	// 兩個平面的地面層與導航位元都是從它算的（docs/re/63 §2）。
	layers *[256][Levels]byte

	// nav[平面][y][x] 是導航 byte（門旗標 ＋ 四個方向位元），
	// lvl[平面][y][x] 是站得上去的 Z（noLevel ＝ 沒有地面）。
	// 建法見 ground.go，出處 `sub_1BC39`。
	nav [numPlanes][Height][Width]uint8
	lvl [numPlanes][Height][Width]int8
}

// NewField 從每格的堆疊高度建一張戰場。
//
// stack[y][x] 是那一格的圖塊堆疊層數（0–7）：**堆疊的部分是實心的**，
// 兵站在它上面。堆疊 ≥ 4 的格子原版會給站上去的兵設一個旗標
// （docs/re/11 §4.3），在攻城圖上那些格子構成城牆。
func NewField(stack [][]int, gateX int) *Field {
	f := &Field{gateX: gateX}
	for y := 0; y < Height && y < len(stack); y++ {
		for x := 0; x < Width && x < len(stack[y]); x++ {
			h := stack[y][x]
			if h > Levels {
				h = Levels
			}
			for z := 0; z < h; z++ {
				f.solid[z][y][x] = true
			}
			f.top[y][x] = h
		}
	}
	f.buildNav()
	return f
}

// NewFieldFromTiles 從戰場的**原始圖塊值**建一張戰場。
//
// heights[圖塊] 是那個圖塊的堆疊層數（`BATTLE.MDL` 的圖塊定義第一個 byte）。
// 與 NewField 的差別是這個版本記得住圖塊值，所以城壁被打壞時
// 可以換成瓦礫的圖塊再重算高度——原版 `sub_1B824` 做的正是這件事。
func NewFieldFromTiles(tiles [][]byte, heights *[256]int, gateX int) *Field {
	return NewFieldFromTileLayers(tiles, heights, nil, gateX)
}

// NewFieldFromTileLayers 與 NewFieldFromTiles 相同，但多吃一張七層子圖塊表
// （`internal/assets/battle` 的 `TileLayers`）。**有它才算得出兩個平面的
// 地面圖與導航位元**（docs/spec/36）；傳 nil 會退化成單平面。
func NewFieldFromTileLayers(tiles [][]byte, heights *[256]int, layers *[256][Levels]byte, gateX int) *Field {
	f := &Field{gateX: gateX, tiles: tiles, heights: heights, layers: layers}
	for y := 0; y < Height && y < len(tiles); y++ {
		for x := 0; x < Width && x < len(tiles[y]); x++ {
			f.setCell(x, y, heights[tiles[y][x]])
		}
	}
	f.buildNav()
	return f
}

// setSolidFromLayers 重現 `sub_1BB6D`：一層擋不擋人看它的子圖塊——
// **0 與 ≥ 0x70 不擋，1–0x6F 擋**。`setCell` 那套「堆疊高度以下全實心」
// 只在沒有子圖塊表時才對；有表的時候差很多（例如圖塊 0xC9 堆七層，
// 實際只有第 3 層擋人）。
func (f *Field) setSolidFromLayers(x, y int) {
	rec := f.layers[f.tiles[y][x]]
	for z := 0; z < Levels; z++ {
		v := rec[z]
		f.solid[z][y][x] = v != 0 && v < topFace
	}
}

func (f *Field) setCell(x, y, h int) {
	if h > Levels {
		h = Levels
	}
	if h < 0 {
		h = 0
	}
	for z := 0; z < Levels; z++ {
		f.solid[z][y][x] = z < h
	}
	f.top[y][x] = h
}

// Retile 把 (x, y) 的圖塊值加上 delta 再重算高度。
//
// 原版打壞城壁時圖塊值**未滿 0xF0 加 0x10、否則加 8**，然後重新展開
// 那個圖塊的堆疊（`sub_1B824` → `sub_1BB6D`）。沒有圖塊資料時什麼都不做。
func (f *Field) Retile(x, y, delta int) {
	if f.tiles == nil || f.heights == nil || !inBounds(x, y) {
		return
	}
	if y >= len(f.tiles) || x >= len(f.tiles[y]) {
		return
	}
	t := int(f.tiles[y][x]) + delta
	if t > 0xFF {
		t = 0xFF
	}
	f.tiles[y][x] = byte(t)
	f.setCell(x, y, f.heights[t])
	if f.layers != nil {
		f.setSolidFromLayers(x, y)
	}
	// 原版打壞城壁時**不重算地面表**（`sub_1B824` 只重跑 `sub_1BB6D`），
	// 而那不影響結果：城壁的地面層本來就是拿**打壞後那個圖塊**算的
	// （`sub_1BCA6` 把 0xD0–0xDF 換成 +0x10，docs/re/63 §2），
	// 所以重算是恆等變換。真正需要它的是沒有子圖塊表的合成戰場——
	// 那種戰場的高度直接來自 heights，不重算就會留著舊的連通性。
	f.buildNav()
}

// Tile 回傳 (x, y) 目前的圖塊值。沒有圖塊資料時回 0。
func (f *Field) Tile(x, y int) byte {
	if f.tiles == nil || y < 0 || y >= len(f.tiles) || x < 0 || x >= len(f.tiles[y]) {
		return 0
	}
	return f.tiles[y][x]
}

// HasTiles 回報這張戰場帶不帶原始圖塊值。
func (f *Field) HasTiles() bool { return f.tiles != nil }

// GateX 回傳登城點的 X。0 表示這是野戰用的戰場。
func (f *Field) GateX() int { return f.gateX }

// IsSiege 回報這是不是攻城戰的戰場。
//
// 判準與原版一致：**索引第二欄為 0 就是野戰**
// （214 張裡零例外，docs/re/11 §4.5）。
func (f *Field) IsSiege() bool { return f.gateX != 0 }

// StandLevel 回傳 (x, y) 那一格站得上去的層。
func (f *Field) StandLevel(x, y int) int {
	if !inBounds(x, y) {
		return 0
	}
	return f.top[y][x]
}

// HighTerrain 對應原版 `sub_1B240` 對 +0x00 bit 1 的設定：
// 目前格子的堆疊高度至少 4 層。
func (f *Field) HighTerrain(x, y int) bool {
	return f.StandLevel(x, y) >= 4
}

// Blocked 回報 (x, y) 的第 z 層有沒有地形擋著。
//
// 對應原版通行圖的 bit 7（`sub_1BB6D`）。移動時檢查的是**兵頭頂那一層**
// 在目標格的值（原版的 `es:[bx+1000h]`，docs/re/63 §5）——
// 兵腳下那一層是頂面，本來就不擋人。
func (f *Field) Blocked(x, y, z int) bool {
	if !inBounds(x, y) || z < 0 || z >= Levels {
		return false
	}
	return f.solid[z][y][x]
}

// Walkable 回報 (x, y, z) 站不站得上去：那一層本身不能是實心的，
// 而且要正好站在堆疊的頂上（不能浮空）。
func (f *Field) Walkable(x, y, z int) bool {
	if !inBounds(x, y) || z < 0 || z >= Levels {
		return false
	}
	return f.top[y][x] == z
}

func inBounds(x, y int) bool {
	return x >= MinCoord && x <= MaxCoord && y >= MinCoord && y <= MaxY
}

// clamp 重現 `sub_1AACF`：X 夾在 1–62。
func clamp(v int) int {
	if v <= 0 {
		return MinCoord
	}
	if v >= 63 {
		return MaxCoord
	}
	return v
}

// clampY 是 Y 的版本，上限是實際的列數（見 MaxY）。
func clampY(v int) int {
	if v <= 0 {
		return MinCoord
	}
	if v > MaxY {
		return MaxY
	}
	return v
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
