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

// SpriteFor 回傳某個兵種第 pose 個姿勢的圖號。
// kind 直接傳兵種的**儲存值**（0／18／36／54），因為它本身就是索引。
func SpriteFor(kind, pose int) int {
	return kind + pose%PosesPerKind
}
