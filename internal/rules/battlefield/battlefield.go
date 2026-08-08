// Package battlefield 決定一場戰鬥用哪一張戰場。
//
// 公式出自 `KI.EXE` 的 `sub_14B63`（取樣周圍地形）、`sub_14C4C`（圖塊 → 地形類型）、
// `sub_14BDD`（平原的配對查表）與 `sub_14C1A`（水域），
// 反組譯見 docs/re/05 與 docs/re/11 §4.4。
//
//   - **攻城戰**：戰場編號就是據點編號（0–191）
//   - **野戰**：從大地圖上**即時算出來**，不是隨機挑一張
//
// 兩張表都從機器碼抄成 Go 常數，並用測試對著原版的 `KI.EXE` 驗
// （沒有原版檔就跳過那條測試）。
package battlefield

// 野戰戰場的編號區間。
const (
	NumCityFields = 192  // 0–191 是攻城用的，與據點一一對應
	FieldBase     = 0xC0 // 平原那一組的基底
	TerrainBase   = 0xCE // 地形類型 1–7 的基底 → 0xCF–0xD5
	NumFields     = 214
)

// terrainRanges 是 `sub_14C4C` 的 14 筆範圍表（`cs:982Fh`）。
// 逐筆比對，命中就回傳類型；全部不中回傳 0（平原）。
var terrainRanges = [...]struct {
	Kind   int
	Lo, Hi byte
}{
	{1, 0xB8, 0xB9}, // 跨水構造（橋／渡口）
	{2, 0xBA, 0xBF}, // 跨水構造（另一組）
	{3, 0x70, 0xA7}, // 山地／丘陵
	{3, 0xA9, 0xAF},
	{4, 0x06, 0x06}, // 林地
	{4, 0x1D, 0x1D},
	{4, 0xB0, 0xB0},
	{4, 0xB6, 0xB7},
	{5, 0xB1, 0xB3}, // 林地（另一種）
	{6, 0x0E, 0x0F}, // 水域
	{6, 0x1E, 0x6F},
	{7, 0xA8, 0xA8}, // 水（純藍網點）
	{8, 0xCA, 0xCA}, // 水上木構造（碼頭）
	{9, 0xC0, 0xC3}, // 水上木構造（橋／棧道）
}

// Terrain 把大地圖的圖塊值換成地形類型（0 ＝ 平原）。
func Terrain(tile byte) int {
	for _, r := range terrainRanges {
		if tile >= r.Lo && tile <= r.Hi {
			return r.Kind
		}
	}
	return 0
}

// plainPairs 是 `sub_14BDD` 的 21 筆配對表（`cs:97F0h`）。
//
// 拿**兩格的地形類型**去配；配不到就用偏移 6。
// ⭐ **換過順序才配上時要把戰場轉 180 度**——`xor dh, 40h` 設的
// 正是 `byte_10D35` 的 bit 6，而載入器看的也是同一個位元。
// 戰場轉不轉，取決於兩格地形誰在前，不是另外算方向的（docs/re/11 §4.4）。
var plainPairs = [...]struct{ Offset, A, B int }{
	{0, 3, 3}, {1, 3, 0}, {2, 3, 5}, {3, 3, 6}, {4, 3, 4}, {5, 3, 7},
	{6, 0, 0}, {6, 0, 4}, {7, 0, 5}, {8, 0, 6}, {8, 5, 6}, {9, 5, 5},
	{10, 5, 4}, {11, 4, 4}, {12, 4, 6}, {13, 6, 6}, {14, 7, 7},
	{6, 0, 7}, {9, 5, 7}, {14, 6, 7}, {11, 4, 7},
}

// plainFallback 是配不到任何一組時用的偏移（原版 `mov dl, 6`）。
const plainFallback = 6

// Neighbours 是 `sub_14B63` 取樣的五格，全部相對於軍團所在格。
//
//	Centre     +0
//	Down       +384      （一列 384 格）
//	DownLeft   +384 − 1
//	DownRight  +384 + 1
//	TwoDown    +768
type Neighbours struct {
	Centre, Down, DownLeft, DownRight, TwoDown int
}

// Select 依周圍地形挑一張野戰的戰場。
//
// dir 是軍團的行進方向（軍團記錄 `+0x08`，0–3；4 表示靜止）。
// 回傳戰場編號與**要不要轉 180 度**。
func Select(dir int, n Neighbours) (field int, rotate bool) {
	// 正下方那一格決定走哪一條路（`sub_14B63` 的三分）。
	switch b := n.Down; {
	case b == 0:
		return plainField(dir, n)
	case b >= 8:
		// 類型 8／9 走 `sub_14C1A`：8 從 209–212 隨機挑、9 固定 213。
		// 隨機那一支要亂數，交給 SelectWater。
		return TerrainBase + b, false
	default:
		// 類型 1–7 → 戰場 0xCF–0xD5。
		return TerrainBase + b, false
	}
}

// plainField 重現 `sub_14BDD`：拿兩格的地形類型去配 21 筆的表。
//
// **哪兩格由行進方向決定**：
//
//	0     兩格下方 ＋ 中心
//	1     中心 ＋ 兩格下方   （順序相反）
//	2     左下 ＋ 右下
//	3     右下 ＋ 左下       （順序相反）
//	其餘  左下 ＋ 右下
func plainField(dir int, n Neighbours) (int, bool) {
	var a, b int
	switch dir {
	case 0:
		a, b = n.TwoDown, n.Centre
	case 1:
		a, b = n.Centre, n.TwoDown
	case 3:
		a, b = n.DownRight, n.DownLeft
	default:
		a, b = n.DownLeft, n.DownRight
	}
	for _, p := range plainPairs {
		if p.A == a && p.B == b {
			return FieldBase + p.Offset, false
		}
		// ⭐ 換過順序才中 → 戰場轉 180 度。
		if p.A == b && p.B == a {
			return FieldBase + p.Offset, true
		}
	}
	return FieldBase + plainFallback, false
}

// SelectWater 處理地形類型 8 的隨機挑選（`sub_14C1A`）。
// 類型 8 從 209–212 隨機挑一張，類型 9 固定 213。
func SelectWater(kind int, rand int) int {
	if kind == 9 {
		return 213
	}
	return 209 + rand&3
}
