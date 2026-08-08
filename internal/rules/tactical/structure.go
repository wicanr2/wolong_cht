package tactical

// 城壁與門。
//
// 原版把它們**當成兵一樣的實體**放進同一個陣列：兵在 `word_1D30E` 的
// 0x0000–0x0BFF（96 個 × 32 byte），城壁與門接在 0x0C00–0x0DFF，
// 一樣是 32 byte 一筆，最多 16 筆（`sub_19CB3` 把 es 設成
// `word_1D30E + 0xC0` 段，兩個建立常式都在 `di >= 0x200` 就停）。
//
// 建立的來源不是另外一張表，而是**戰場圖塊值本身**：
//
//	0xD0–0xDF  城壁（`sub_19CE2`）
//	0xF0–0xF7  門　（`sub_19DA1`）
//
// 掃描是逐行往下走（si 每次 +0x40 ＝ 往下一列，走完一整行才換下一行），
// **一「行連續的城壁圖塊」只產生一筆實體**，長度記在 `+0x1A`；
// 打壞它時整段一起消失（`sub_1B824` 用 `+0x1A` 當迴圈次數）。

// 圖塊值的兩段區間（`sub_19CE2` 與 `sub_19DA1` 的上下界）。
const (
	TileWallLo = 0xD0
	TileWallHi = 0xDF
	TileGateLo = 0xF0
	TileGateHi = 0xF7

	// 打壞之後圖塊值的變化（`sub_1B824`）：**未滿 0xF0 加 0x10，否則加 8**。
	// 城壁 0xD0–0xDF → 0xE0–0xEF（瓦礫），門 0xF0–0xF7 → 0xF8–0xFF（破門）。
	brokenWallDelta = 0x10
	brokenGateDelta = 8
)

// MaxStructures 是城壁加門的上限。原版兩個常式共用 0x0C00–0x0DFF，
// 一筆 32 byte，所以是 16 筆——**滿了就不再建**，多出來的城壁打不壞。
const MaxStructures = 16

// 耐久的三個來源（`sub_19CE2`／`sub_19DA1`）。
const (
	// FieldWallDurability 是**野戰**戰場上城壁的耐久（`mov word ptr es:[di+18h], 12Ch`）。
	FieldWallDurability = 300
	// GateDurability 是門的耐久，攻城野戰都一樣（`mov word ptr es:[di+18h], 50h`）。
	GateDurability = 80

	// 攻城戰的城壁耐久 ＝（據點記錄 +0x13 ＋ 50）× 10。
	wallDurabilityBase  = 50
	wallDurabilityScale = 10
)

// 實體的類型碼（`es:[di]` 那個 word 的高位元組）。
const (
	KindWall = 1 // `mov word ptr es:[di], 180h`
	KindGate = 2 // `mov word ptr es:[di], 280h`
)

// SiegeWallDurability 回傳攻城戰的城壁耐久。
// cityWall 是據點記錄的 `+0x13`（docs/formats/08）。
func SiegeWallDurability(cityWall int) int {
	return (cityWall + wallDurabilityBase) * wallDurabilityScale
}

// Structure 是一段城壁或一道門。
type Structure struct {
	Kind int // KindWall／KindGate
	X, Y int // 這一段的最上面那一格
	Run  int // 往下連續幾格（原版 `+0x1A`）

	// Durability 是剩餘耐久（原版 `+0x18`，word）。
	Durability int
	// Broken 為真表示已經打壞了（原版 `+0x00` 的 bit 0）。
	Broken bool
}

// buildStructures 重現 `sub_19CE2` ＋ `sub_19DA1`。
//
// tiles[y][x] 是戰場的原始圖塊值。siege 決定城壁的耐久要用哪一個公式，
// cityWall 是攻城時的據點 `+0x13`（野戰用不到）。
//
// **先掃城壁再掃門**，兩者共用 16 筆的額度——與原版一致
// （`sub_19CB3` 連著呼叫，di 沒有歸零）。
func buildStructures(tiles [][]byte, siege bool, cityWall int) []Structure {
	if tiles == nil {
		return nil
	}
	wallHP := FieldWallDurability
	if siege {
		wallHP = SiegeWallDurability(cityWall)
	}

	var out []Structure
	// 城壁：逐行往下，連續的算同一段。
	for x := 0; x < Width; x++ {
		run := 0
		for y := 0; y < Height; y++ {
			t := tileAt(tiles, x, y)
			if t >= TileWallLo && t <= TileWallHi {
				if run == 0 {
					if len(out) >= MaxStructures {
						return out
					}
					out = append(out, Structure{
						Kind: KindWall, X: x, Y: y, Durability: wallHP,
					})
				}
				run++
				continue
			}
			if run > 0 {
				out[len(out)-1].Run = run
				run = 0
			}
		}
		if run > 0 {
			out[len(out)-1].Run = run
		}
	}

	// 門：一格一筆，不合併。
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			t := tileAt(tiles, x, y)
			if t < TileGateLo || t > TileGateHi {
				continue
			}
			if len(out) >= MaxStructures {
				return out
			}
			out = append(out, Structure{
				Kind: KindGate, X: x, Y: y, Run: 1, Durability: GateDurability,
			})
		}
	}
	return out
}

func tileAt(tiles [][]byte, x, y int) byte {
	if y < 0 || y >= len(tiles) || x < 0 || x >= len(tiles[y]) {
		return 0
	}
	return tiles[y][x]
}

// MinWallDurability 重現 `sub_1A65D`：掃 16 筆，**只看城壁不看門**，
// 回傳最小耐久與「有沒有任何一段已經被打壞」。
//
// 一段都沒有時回傳 (0xFFFF, true)——原版的 dx 初值就是 0xFFFF。
func (b *Battle) MinWallDurability() (min int, intact bool) {
	min, intact = 0xFFFF, true
	for i := range b.Structures {
		s := &b.Structures[i]
		if s.Kind != KindWall {
			continue
		}
		if s.Durability < min {
			min = s.Durability
		}
		if s.Broken {
			intact = false
		}
	}
	return
}

// WallQuery 是腳本指令 15 查到的值（`sub_1A65D` 寫進 byte_1D315 的那個）。
//
//	任何一段城壁破了　→　0
//	否則　　　　　　　→　最小耐久 × 4 的高位元組（＝ 最小耐久 ÷ 64）
func (b *Battle) WallQuery() int {
	min, intact := b.MinWallDurability()
	if !intact {
		return 0
	}
	return (min * 4) >> 8 & 0xFF
}

// CityDamage 回傳這一場攻城戰打掉據點多少點，重現 `sub_19FDC` 的前半：
//
//	損失 ＝（據點+0x13 ＋ 50 −　最小耐久 ÷ 10）÷ 8
//
// 這個值會同時從據點的 `+0x10`、`+0x11`、`+0x13` 各扣一次（都夾在 0）。
func (b *Battle) CityDamage(cityWall int) int {
	min, _ := b.MinWallDurability()
	d := (cityWall + wallDurabilityBase - min/wallDurabilityScale) >> 3
	if d < 0 {
		return 0
	}
	return d
}

// hitStructure 是兵撞上城壁／門時的處理，重現 `seg000:B5B7`–`B5FE`：
//
//	守方（side 1）→ 直接離開，**碰不壞城壁**（`cmp si, 600h / jnb loc_1B612`）
//	攻方面向「背對城的方向」→ 耐久**直接歸零**
//	其餘 → 耐久減 1；歸零之後再撞一次，同一列一起垮
//
// ⭐ 那個「直接歸零」的分支每一個運算元都查證過了：
//
//   - `byte_1D34B == 0` ＝ 攻城戰（§3.2）
//   - `byte_10D35` bit 7 ＝ **玩家是守方**（`sub_14E5C` 只在互換
//     `word_10D2E`／`word_10D30` 那一支設它）
//   - `si < 0x600` ＝ side 0；配上 bit 7 之後**兩種組合指的都是攻方**
//   - `[si+05]` ＝ 面向（`sub_1B047` 等四支常式寫 0／1／2／3 的同時
//     各自 `dec [si+6]`／`dec [si+8]`／`inc [si+6]`／`inc [si+8]`）
//
// ⚠ **語意存疑，但行為照抄。** 186 張攻城圖的城壁一律在 X 33–46，
// 攻方從 X 5 出發，所以「朝城走」是 East；條件要的卻是 West——
// 而面向只有走成功才更新，**攻方唯一會面向 West 的時機是退卻**。
// 兩種鏡射組合下這個不對稱是一致的，所以機制是確定的；
// 它究竟是設計還是原版的 bug，沒有資料可以判斷。**照抄。**
func (b *Battle) hitStructure(side, facing, x, y int) bool {
	i := b.structureAt(x, y)
	if i < 0 {
		return false
	}
	s := &b.Structures[i]
	if s.Broken {
		return false
	}
	// 守方碰不壞自己的城壁。
	if side != 0 {
		return false
	}
	// 攻方背對城的方向 → 直接歸零。
	if b.Field.IsSiege() && facing == awayFromCastle {
		s.Durability = 0
	}
	if s.Durability > 0 {
		s.Durability--
		return true
	}
	b.breakRow(s.Y)
	return true
}

// awayFromCastle 是攻方「背對城」的面向。
//
// 攻城圖的城壁一律在 X 33–46（186 張零例外），攻方的陣形原點在 X 5，
// 所以朝城是 East、背對是 West。
const awayFromCastle = West

// structureAt 找出蓋住 (x, y) 的那一段。一段城壁蓋 Run 格。
func (b *Battle) structureAt(x, y int) int {
	for i := range b.Structures {
		s := &b.Structures[i]
		if s.X == x && y >= s.Y && y < s.Y+s.Run {
			return i
		}
	}
	return -1
}

// breakRow 重現 `sub_1B799`：**打壞的不只是撞到的那一段，
// 而是所有 `+0x08`（Y）相同的實體**——同一排的城壁一起垮。
func (b *Battle) breakRow(y int) {
	for i := range b.Structures {
		s := &b.Structures[i]
		if s.Y != y || s.Broken {
			continue
		}
		s.Broken, s.Durability = true, 0
		delta := brokenWallDelta
		if s.Kind == KindGate {
			delta = brokenGateDelta
		}
		for r := 0; r < s.Run; r++ {
			b.Field.Retile(s.X, s.Y+r, delta)
		}
		b.Log = append(b.Log, structureName(s.Kind)+"被打壞了")
	}
}

// OpenGates 重現 `sub_1B7CB`：突擊時守方把門全部打開。
// 說明書 4.2 說開了的門這場戰鬥不能再關，所以這裡也不還原。
func (b *Battle) OpenGates() {
	for i := range b.Structures {
		s := &b.Structures[i]
		if s.Kind != KindGate || s.Broken {
			continue
		}
		s.Broken, s.Durability = true, 0
		for r := 0; r < s.Run; r++ {
			b.Field.Retile(s.X, s.Y+r, brokenGateDelta)
		}
	}
}

func structureName(kind int) string {
	if kind == KindGate {
		return "門"
	}
	return "城壁"
}
