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
// ⚠ **180 個裡面誰是誰還沒解。** 已知的只有「相隔偶數個的單位彼此相似、
// 奇數個的不相似」，也就是**兩兩一組**；那一組是「兩張動畫幀」還是
// 「圖形 ＋ 陰影」還沒有證據。所以下面只提供「第 n 個」的取用，
// 兵種與方向怎麼對應是**呼叫端的暫定選擇**，不是原版行為。
const (
	SpriteUnit     = SubTileSize // 320，與子圖塊同格式
	NumSprites     = 360
	SpritesPerSide = 180
	// EmptySprite 是兩側都空著的那一格（170 與 350）。
	EmptySprite = 170
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

// Side 取第 side 側的第 n 張（n 是 0–179）。
func (s *Sprites) Side(side, n int) *SubTile {
	if n < 0 || n >= SpritesPerSide {
		return nil
	}
	return s.At(side*SpritesPerSide + n)
}
