package battle

import "fmt"

// `BATTLE.SCH` 的人物圖形。
//
// 規格與驗證見 docs/formats/07 §10。
//
// 115,200 B ＝ **360 個 320 B 的單位**，格式與 `BATTLE.MDL` 的子圖塊
// **完全一樣**：五個 64 B 的位元平面（遮罩 ＋ 4bpp 的四張），16 × 32。
// 同一條不變量（遮罩為 0 處色平面全 0）在這裡也是 **100.00%**。
//
// ⭐ **360 ＝ 兩側 × 180**。兩條證據：
//
//   - 整份檔案裡**只有兩個單位是全空的：170 與 350**，正好相差 180
//   - 前半是紅色系、後半是藍色系（攻方／守方）
//
// ⭐ **一張圖是 16 × 64，由兩個單位疊起來：奇數在上、偶數在下。**
// 量遮罩的接合率就看得出來——「單位 i+1 的底列」接「單位 i 的頂列」
// 是 27.6%，其餘三種接法都在 5–17%（隨機基準約 6%）。
// 所以一側 180 個單位 ＝ **90 張圖**。
//
// ⭐ **90 張分成五組，每組 18 張，而分界正好是 0／18／36／54／72。**
// 那就是兵種的儲存值——**兵種之所以存成「× 18」，是因為它就是圖形表的索引**
// （`sub_1AB7C` 的門檻全是 18 的倍數，同一個理由）。
//
//	兵種  0（大將） → 圖 0–17
//	兵種 18（騎馬） → 圖 18–35   ← 這一組 37–42% 的像素是棕色（馬）
//	兵種 36（弓兵） → 圖 36–53
//	兵種 54（步兵） → 圖 54–71
//	　　　　軍旗　　 → 圖 72–89   ← 白桿紅旗／藍旗，佔滿 64 px 的高度
//
// ⚠ 一組裡的 18 張是什麼（方向 × 動作幀？）**還沒解**。
const (
	SpriteUnit     = SubTileSize // 320，與子圖塊同格式
	NumSprites     = 360
	UnitsPerSide   = 180
	SpritesPerSide = UnitsPerSide / 2 // 90
	SpriteW        = SubTileW         // 16
	SpriteH        = SubTileH * 2     // 64
	// PosesPerKind 是一個兵種有幾張圖。**與兵種的儲存倍率同一個數**。
	PosesPerKind = 18
	// BannerSprite 是軍旗那一組的起點。
	BannerSprite = 72
	// EmptyUnit 是兩側都空著的那一格（170 與 350）。
	EmptyUnit = 170
)

// Sprites 是解好的人物圖形。
type Sprites struct {
	raw   []byte
	cache [NumSprites]*SubTile
}

// ParseSprites 解 `BATTLE.SCH`。
func ParseSprites(b []byte) (*Sprites, error) {
	if n := NumSprites * SpriteUnit; len(b) != n {
		return nil, fmt.Errorf("battle: BATTLE.SCH 是 %d B，預期 %d", len(b), n)
	}
	return &Sprites{raw: b}, nil
}

// At 取第 n 個單位（0–359）。超出範圍回 nil。
func (s *Sprites) At(n int) *SubTile {
	if s == nil || n < 0 || n >= NumSprites {
		return nil
	}
	if t := s.cache[n]; t != nil {
		return t
	}
	t := decodePlanar(s.raw[n*SpriteUnit : (n+1)*SpriteUnit])
	s.cache[n] = t
	return t
}

// Frame 是一張 16 × 64 的人物圖形，Pix[y*16+x] 是色號或 Transparent。
type Frame struct {
	Pix [SpriteW * SpriteH]int8
}

// At 回傳 (x, y) 的色號，透明處回傳 Transparent。
func (f *Frame) At(x, y int) int {
	if x < 0 || x >= SpriteW || y < 0 || y >= SpriteH {
		return Transparent
	}
	return int(f.Pix[y*SpriteW+x])
}

// Sprite 取第 side 側的第 n 張圖（n 是 0–89）。
//
// **奇數單位在上、偶數在下**——別把順序寫反了，寫反會變成
// 「旗子在腳邊、桿子在頭上」。
func (s *Sprites) Sprite(side, n int) *Frame {
	if s == nil || side < 0 || side > 1 || n < 0 || n >= SpritesPerSide {
		return nil
	}
	base := side*UnitsPerSide + n*2
	top, bottom := s.At(base+1), s.At(base)
	if top == nil || bottom == nil {
		return nil
	}
	f := &Frame{}
	copy(f.Pix[:SubTileW*SubTileH], top.Pix[:])
	copy(f.Pix[SubTileW*SubTileH:], bottom.Pix[:])
	return f
}

// 姿勢的位元編碼（`sub_1B240` 尾段，docs/re/11 §5.13）。
//
//	cl  = [si+05] × 2          面向 0–3 → 0／2／4／6
//	ch  = [si+02] & 0x19       狀態旗標的 bit 0、3、4
//	     bit 4 設了 → **面向歸零**
//	cl |= ch ; cl += [si+04]   ＋ 兵種（已經 × 18）
//	side 1 → cx += 0x5A（90）
//	cx = cx × 2 + 0xC0（192）  → 合併表裡的單位編號
//
// 最後那個 `+192` 正是**地形子圖塊的張數**——原版把地形與人物放在
// 同一張表裡（`word_1E15A`），前 192 個單位是地形（`BATTLE.MDL` 的
// 61,440 B ÷ 320），後面接 `BATTLE.SCH`。兩個檔案在記憶體裡是連著的
// （`sub_1CC31` 把 `ds:0D304` 排在 `ds:0D302` 之後）。
//
// ⭐ **bit 0 是走路的動畫幀**：`sub_1B240` 每次更新完就 `xor [si+2], 1`。
const (
	// FacingStride 是面向在圖號裡的間隔。
	FacingStride = 2
	// PoseFlagMask 是狀態旗標裡會進圖號的位元。
	PoseFlagMask = 0x19
	// PoseFlagStep 是動畫幀那一位（bit 0），每次更新翻面。
	PoseFlagStep = 0x01
	// PoseFlagFront 是「面向歸零」那一位（bit 4）。
	PoseFlagFront = 0x10
)

// SpriteFor 回傳一個兵要用第幾張圖（0–89，同一側內）。
//
// kind 傳兵種的**儲存值**（0／18／36／54），因為它本身就是索引；
// facing 是 0–3；flags 是兵記錄 `+0x02` 的狀態旗標。
func SpriteFor(kind, facing, flags int) int {
	f := flags & PoseFlagMask
	pose := facing * FacingStride
	if f&PoseFlagFront != 0 {
		pose = 0 // bit 4 設了就一律用正面
	}
	n := kind + (pose | f)
	if n < 0 || n >= SpritesPerSide {
		return 0
	}
	return n
}

// ---------------------------------------------------------------------------
// 場上的軍旗
// ---------------------------------------------------------------------------

// 戰場上插的旗（`sub_19E10`，docs/re/11 §5.14）。
//
// 載入時掃過 64 × 64 每一格，**看那一格圖塊的最頂層子圖塊**：
// 落在 `0xBA`–`0xBF` 就在那裡插一支旗，記錄放在單位區的 `0x0E00` 起
// （接在 96 個兵與 16 段城壁門後面），類型碼 3。
//
//	es:[di]    = 0x3C0   旗標 0xC0、類型 3
//	es:[di+6]  = X
//	es:[di+8]  = Y
//	es:[di+0A] = 堆疊高度（＝ 插在地面上）
//	es:[di+1B] = 亂數 0–3
//	es:[di+1C] = 0x150 或 0x204   ← 最頂層子圖塊的**最低位**決定
//
// ⭐ 那兩個常數就是圖號：`(0x150 − 192) ÷ 2 = 72` ＝ **軍旗那一組的起點**，
// `(0x204 − 192) ÷ 2 = 162 = 72 + 90` ＝ **另一側的軍旗**。
// 換句話說旗子的顏色是**圖塊編號的最低位**選的，與交戰雙方無關。
//
// 拿 214 張戰場驗過：每張 0–48 支，**沒有一張超過 `0x0E00`–`0x1800`
// 放得下的 80 筆**；194 張有旗。
const (
	TopSubTileFlagLo = 0xBA
	TopSubTileFlagHi = 0xBF
	// MaxBanners 是 `0x0E00` 到繞路點區 `0x1800` 之間放得下的筆數。
	MaxBanners = (0x1800 - 0x0E00) / 32 // 80
)

// Banner 是場上的一支旗。
type Banner struct {
	X, Y, Z int
	// Side 是旗色，由最頂層子圖塊的最低位決定（0 → 圖 72、1 → 圖 162）。
	Side int
	// Variant 是原版寫進 `+0x1B` 的亂數 0–3。
	// ⚠ 它的用途還沒解（旗號是 `+0x1C` 直接給的，不經過這個值）。
	Variant int
}

// Banners 回傳第 n 張戰場上所有的旗。
//
// rand 每支旗會被呼叫一次，對應原版的 `call sub_1ECE0 / and al, 3`。
// 傳 nil 就一律用 0。
func (l *Library) Banners(n int, rand func() int) []Banner {
	t := l.TileSet(n)
	off := FieldsBase + n*FieldSize + CellsOff
	cells := l.mapData[off : off+NumCells]
	var out []Banner
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			st := l.stacks[t][cells[y*Width+x]]
			if len(st) == 0 {
				continue
			}
			top := st[len(st)-1]
			if top < TopSubTileFlagLo || top > TopSubTileFlagHi {
				continue
			}
			if len(out) >= MaxBanners {
				return out
			}
			v := 0
			if rand != nil {
				v = rand() & 3
			}
			out = append(out, Banner{
				X: x, Y: y, Z: len(st), Side: int(top & 1), Variant: v,
			})
		}
	}
	return out
}
