package tactical

// 兩個平面的地面圖與導航位元。出處 docs/re/63、規格 docs/spec/36。
//
// 原版的戰場不是一張高度圖，是**兩個平面**：低平面（地面，層 0–3）
// 與高平面（城牆頂，層 4–6）。每一格每個平面存一個 byte：
//
//	bit 0–2  地面層          bit 3  這一格是門（可以上下）
//	bit 4 ←  bit 5 →  bit 6 ↑  bit 7 ↓     ← 這四個方向走得通
//
// 尋路直接吃這張表（`0001BE10` 的 `test al, 10h/20h/40h/80h/08h`），
// 移動則另外用地面層做高度同步（`sub_1B1B1`）。

// 平面編號。原版是位址的 bit 12（`bx | 0x1000`），這裡用索引。
const (
	PlaneLow  = 0
	PlaneHigh = 1
	numPlanes = 2
)

// 導航 byte 的位元。方向與原版的 `sub_1BD07` 傳進來的 `bl` 相同。
//
// 原版把地面層擠在同一個 byte 的 bit 0–2，這裡把層數另外存成 lvl，
// **因為 remake 的 Z 從 0 起算**（原版靠「整個 byte 是 0」表示沒有地面，
// 那在 Z 可以是 0 的座標系裡會與「地面在第 0 層」撞號）。
const (
	navGate  = 0x08 // 這一格是門：可以上下換平面
	navLeft  = 0x10
	navRight = 0x20
	navUp    = 0x40
	navDown  = 0x80
)

// noLevel 是 lvl 裡的「這個平面沒有地面」。
const noLevel = -1

// noGround 是 `sub_1BCA6` 的「這個平面沒有地面」（原版用 0x80）。
const noGround = 0x80

// topFace 是「可以站的頂面」的子圖塊門檻。同一個門檻 `sub_1BB6D`
// 也在用：0 與 ≥ 0x70 不擋，1–0x6F 是實體。
const topFace = 0x70

// lowPlaneLevels 是低平面在哪幾層找地面（`sub_1BCA6` 的 `cmp dl, 4`）。
const lowPlaneLevels = 4

// gateTile 是「可以上下」的圖塊下界，gateSolid 是「未破、高平面走不過去」的上界。
const (
	gateTile  = 0xF0
	gateSolid = 0xF8
)

// wallTile 是城壁的圖塊範圍；建表時改讀 +0x10 那個圖塊的屬性
// （`sub_1BCA6` 的 `cmp bx, 680h/700h`）。
const (
	wallTileLo = 0xD0
	wallTileHi = 0xDF
)

// groundOf 重現 `sub_1BCA6`：回傳 (低平面層, 高平面層)，noGround ＝ 沒有。
// layers[圖塊] 是那個圖塊由下往上七層的子圖塊。
//
// ⚠ 回傳的是**兵腳下佔的那一層**（原版的 L，0 起算）：頂面那一層
// 本身不擋人（`sub_1BB6D` 對 ≥ 0x70 不設 bit 7），兵就站在它上面。
// noGround ＝ 這個平面沒有地面。
func groundOf(tile byte, layers *[256][Levels]byte) (int, int) {
	rec := layers[tile]
	if tile >= wallTileLo && tile <= wallTileHi {
		rec = layers[int(tile)+0x10]
	}
	l := 0
	for l < lowPlaneLevels && rec[l] < topFace {
		l++
	}
	if l >= lowPlaneLevels {
		return noGround, noGround
	}
	lo := l
	if tile >= gateTile {
		// 門在高平面標 8；`sub_1BD07` 會把它換成鄰格牆頂的高度。
		return lo, navGate
	}
	for l = lowPlaneLevels; l < Levels && rec[l] < topFace; l++ {
	}
	if l >= Levels {
		return lo, noGround
	}
	return lo, l
}

// buildNav 重現 `sub_1BBA6` → `sub_1BC39` → `sub_1BD07`：
// 算每格兩個平面的地面層、連四個鄰格設方向位元，最後把高平面正規化成
// 單一高度。tiles 為 nil 時退化成「低平面 ＝ 堆疊高度、沒有高平面」，
// **演算法本身是同一套**。
func (f *Field) buildNav() {
	// 有子圖塊表時，先照原版的規則重算每一層擋不擋人。
	if f.layers != nil && f.tiles != nil {
		for y := 0; y < Height && y < len(f.tiles); y++ {
			for x := 0; x < Width && x < len(f.tiles[y]); x++ {
				f.setSolidFromLayers(x, y)
			}
		}
	}
	var lo, hi [Height][Width]int
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			lo[y][x], hi[y][x] = f.cellGround(x, y)
		}
	}
	// gateMark 是 `sub_1BCA6` 回的「門」標記，在 hi 裡用一個不可能的層號表示。
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			l, h := lo[y][x], hi[y][x]
			var nl, nh uint8
			ll, lh := noLevel, noLevel
			if l != noGround {
				ll = l
			}
			switch {
			case h == navGate:
				// `sub_1BC39`：門在兩個平面都設 bit 3。高度先留空，
				// 下面再跟鄰格的牆頂對齊（`sub_1BD07` 的第一行）。
				nl |= navGate
				nh |= navGate
			case h != noGround:
				lh = h
			}
			if h == navGate {
				for _, d := range navSteps {
					nx, ny := x+d.dx, y+d.dy
					if nx < 0 || nx >= Width || ny < 0 || ny >= Height {
						continue
					}
					if ah := hi[ny][nx]; ah != noGround && ah != navGate {
						lh = ah
						break
					}
				}
			}
			for _, d := range navSteps {
				nx, ny := x+d.dx, y+d.dy
				if nx < 0 || nx >= Width || ny < 0 || ny >= Height {
					continue
				}
				al, ah := lo[ny][nx], hi[ny][nx]
				if l != noGround && al != noGround && abs(al-l) <= 1 {
					nl |= d.bit
				}
				if lh != noLevel && ah != noGround && (ah == navGate || ah == lh) {
					nh |= d.bit
				}
			}
			f.lvl[PlaneLow][y][x] = int8(ll)
			f.lvl[PlaneHigh][y][x] = int8(lh)
			f.nav[PlaneLow][y][x] = nl
			f.nav[PlaneHigh][y][x] = nh
		}
	}
	f.normaliseHighPlane()
}

// navSteps 是四個方向與它們的位元。順序與原版的檢查順序相同。
var navSteps = [4]struct {
	dx, dy int
	bit    uint8
}{{-1, 0, navLeft}, {1, 0, navRight}, {0, -1, navUp}, {0, 1, navDown}}

// cellGround 回傳一格兩個平面的地面層。有圖塊層表就照原版算，
// 沒有就用堆疊高度退化（合成戰場）。
func (f *Field) cellGround(x, y int) (int, int) {
	if f.layers != nil && f.tiles != nil && y < len(f.tiles) && x < len(f.tiles[y]) {
		return groundOf(f.tiles[y][x], f.layers)
	}
	// 合成戰場（NewField／沒有子圖塊表的 NewFieldFromTiles）只有堆疊高度，
	// 而堆疊高度就是 remake 的 Z，所以直接用。**堆到 4 層以上就是城牆**
	// （docs/re/11 §4.3 的同一條判準），歸到高平面。
	h := f.top[y][x]
	if h >= Levels {
		return noGround, noGround
	}
	// 有圖塊值就照原版認門：門是唯一能上下的格子（docs/re/63 §4）。
	if f.tiles != nil && y < len(f.tiles) && x < len(f.tiles[y]) &&
		f.tiles[y][x] >= gateTile {
		return h, navGate
	}
	if h >= lowPlaneLevels {
		return noGround, h
	}
	return h, noGround
}

// normaliseHighPlane 重現 `sub_1BBA6` 尾端那一段：整面牆頂壓成單一高度。
//
//	統計有地面的格子，層 ≤ 4 與 > 4 哪邊多 → 目標高度 4 或 5
//	高度不同的：是門就改成目標高度，不是門就整格清 0
//
// 常數在這裡是 +1 之後的（5 與 6），因為 byte 存的是層號 +1。
// 原版用 bit 3（門）當「這一格有高平面」的判準；byte 存 +1 之後
// **有地面就一定非 0**，兩種寫法等價。
func (f *Field) normaliseHighPlane() {
	// 原版寫死在 4 與 5 之間二選一（`cmp al, 4 / ja`），因為它的牆頂
	// 只可能落在那兩層。這裡改成**取多數**：在原版資料上結果相同
	// （門格旁邊的牆頂都是同一層），而合成戰場的牆可以是別的高度。
	count := map[int8]int{}
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			// 原版只統計**門格**（`test al, 8`）：牆頂那些格子沒有 bit 3。
			if f.nav[PlaneHigh][y][x]&navGate == 0 {
				continue
			}
			if lv := f.lvl[PlaneHigh][y][x]; lv != noLevel {
				count[lv]++
			}
		}
	}
	if len(count) == 0 {
		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				if lv := f.lvl[PlaneHigh][y][x]; lv != noLevel {
					count[lv]++
				}
			}
		}
	}
	if len(count) == 0 {
		return
	}
	want, best := int8(noLevel), -1
	for lv, n := range count {
		if n > best || (n == best && lv < want) {
			want, best = lv, n
		}
	}
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if f.lvl[PlaneHigh][y][x] == want {
				continue
			}
			if f.nav[PlaneHigh][y][x]&navGate != 0 {
				f.lvl[PlaneHigh][y][x] = want
				continue
			}
			f.lvl[PlaneHigh][y][x] = noLevel
			f.nav[PlaneHigh][y][x] = 0
		}
	}
}

// GroundLevel 回傳 (x, y) 在某個平面站得上去的 Z，第二個值是「有沒有地面」。
func (f *Field) GroundLevel(x, y, plane int) (int, bool) {
	if !inBounds(x, y) || plane < 0 || plane >= numPlanes {
		return 0, false
	}
	v := f.lvl[plane][y][x]
	if v == noLevel {
		return 0, false
	}
	return int(v), true
}

// Linked 回報從 (x, y) 在某個平面往 (dx, dy) 走得通不通。
func (f *Field) Linked(x, y, plane, dx, dy int) bool {
	if !inBounds(x, y) || plane < 0 || plane >= numPlanes {
		return false
	}
	for _, d := range navSteps {
		if d.dx == dx && d.dy == dy {
			return f.nav[plane][y][x]&d.bit != 0
		}
	}
	return false
}

// IsGateCell 回報 (x, y) 是不是可以上下的門格（原版圖塊 ≥ 0xF0）。
func (f *Field) IsGateCell(x, y int) bool {
	if !inBounds(x, y) {
		return false
	}
	return f.nav[PlaneLow][y][x]&navGate != 0
}

// GateBlocksHighPlane 回報 (x, y) 是不是**未破的門**：高平面上橫向穿不過去
// （`sub_1B1B1` 的 `cmp dl, 0F0h / cmp dl, 0F8h`）。
func (f *Field) GateBlocksHighPlane(x, y int) bool {
	if !f.HasTiles() {
		return false
	}
	t := f.Tile(x, y)
	return t >= gateTile && t < gateSolid
}
