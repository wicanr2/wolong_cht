package world

import (
	"fmt"
	"image"

	"github.com/wicanr2/wolong_cht/internal/assets/palette"
)

// MMAP.MCH 的前 0xA000 byte 是 256 個 16×16、每個 160 byte 的平面圖塊。
// 這個 160 不是一般 4bpp 圖的 128：前 32 byte 是遮罩，後 128 byte
// 是四個色平面。出處是 KI.EXE 的 sub_1D804（IDA 線性位址 0001D804），
// 它先讀 16 個 mask word，再讀 64 個 color word。
const (
	MCHTileBytes       = 160
	MCHTileCount       = 256
	MCHTileDataBytes   = MCHTileBytes * MCHTileCount // 0xA000
	MCHPatternTable    = 0xA000
	MCHPatternData     = 0xA100
	MCHPatternEntry    = 4
	MCHPatternEntries  = 0x100 / MCHPatternEntry
	MCHFileSize        = 43058 // DOS/V 與 PC-98 的 MMAP.MCH 均為 0xA832 B
	MCHTransparent     = 0xFF
	MCHTilePlaneBytes  = 32
	MCHTileMaskOffset  = 0
	MCHTileColorOffset = MCHTilePlaneBytes
)

// MCHTile 是一張 16×16 MCH 圖塊。Pix 的 0–15 是調色盤索引，0xFF 是
// 由原版 mask 判定的透明像素。
type MCHTile struct {
	Pix [TileSize * TileSize]byte
}

// MCHPattern 是戰略地圖物件由 MCH 圖塊拼成的矩陣。Tiles 逐列排列，
// 0xFF 表示原版 loc_1D51F 不寫入該格。
type MCHPattern struct {
	Width  int
	Height int
	Tiles  []byte
}

// MCH 是 MMAP.MCH 的可查詢解碼結果。
type MCH struct {
	raw   []byte
	cache [MCHTileCount]*MCHTile
}

// ParseMCH 驗證並載入 MMAP.MCH。
//
// 前 0xA000 byte 的固定大小來自 sub_1D804 的 `ah * 160` 位址算法；
// `word_1987A`（sub_187CC，000187CC）指向檔案 offset 0xA000，後面
// 0x100 byte 是物件 metadata，像素矩陣則從 offset 0xA100 開始。
func ParseMCH(data []byte) (*MCH, error) {
	if len(data) != MCHFileSize {
		return nil, fmt.Errorf("world: MMAP.MCH 是 %d B，預期 %d", len(data), MCHFileSize)
	}
	if len(data) < MCHPatternData {
		return nil, fmt.Errorf("world: MMAP.MCH 不足 metadata／物件資料區")
	}
	return &MCH{raw: data}, nil
}

// Tile 取 MCH 的一張 16×16 平面圖塊。超出範圍回 nil。
func (m *MCH) Tile(id byte) *MCHTile {
	if m == nil {
		return nil
	}
	idx := int(id)
	if idx >= len(m.cache) {
		return nil
	}
	if t := m.cache[idx]; t != nil {
		return t
	}
	base := idx * MCHTileBytes
	if base+MCHTileBytes > len(m.raw) {
		return nil
	}
	t := &MCHTile{}
	for y := 0; y < TileSize; y++ {
		for x := 0; x < TileSize; x++ {
			byteIdx := base + y*2 + x/8
			bit := uint(7 - x%8)
			mask := (m.raw[byteIdx] >> bit) & 1
			if mask == 0 {
				t.Pix[y*TileSize+x] = MCHTransparent
				continue
			}
			var colour byte
			for plane := 0; plane < 4; plane++ {
				planeByte := base + MCHTileColorOffset + plane*MCHTilePlaneBytes + y*2 + x/8
				colour |= ((m.raw[planeByte] >> bit) & 1) << uint(plane)
			}
			t.Pix[y*TileSize+x] = colour
		}
	}
	m.cache[idx] = t
	return t
}

// Pattern 取 metadata 表中的第 index 項。
//
// metadata 的四 byte 是寬、 高、 相對於 0xA100 的矩陣位移（little-endian）。
// 這個欄位排列直接對應 sub_12533 讀取 `es:[bx]`、`es:[bx+1]`、
// `es:[bx+2]`，再由 loc_1D51F 把 source base 加上 0x100。
func (m *MCH) Pattern(index int) (MCHPattern, bool) {
	if m == nil || index < 0 || index >= MCHPatternEntries {
		return MCHPattern{}, false
	}
	entry := MCHPatternTable + index*MCHPatternEntry
	width, height := int(m.raw[entry]), int(m.raw[entry+1])
	offset := int(m.raw[entry+2]) | int(m.raw[entry+3])<<8
	count := width * height
	start := MCHPatternData + offset
	if width == 0 || height == 0 || start < MCHPatternData || start+count > len(m.raw) {
		return MCHPattern{}, false
	}
	tiles := make([]byte, count)
	copy(tiles, m.raw[start:start+count])
	return MCHPattern{Width: width, Height: height, Tiles: tiles}, true
}

// ObjectPatternIndex 是原版 CS:[bx-67A6h]（16-bit wrap 後為 CS:985Ah）
// 的固定查表：`bx = 物件 type × 8 + 相位`，所以 **`bx = 0` 落在 985Ah**
// ——那八個 byte 是 **type 0**，不是 type 1（docs/spec/146 §1）。
//
//	type 0：18 19 1A 1B 1C 18 19 1A   16×9 格，會飄的雲（常駐 16 個）
//	type 1：20 21 22 23 20 21 22 23   5×5，火災（sub_134B1 的 ah=1）
//	type 2：28 29 2A 2B 28 29 2A 2B   5×5，暴動（ah=2）
//	type 3：00 …                      3×3，目前沒有可見的產生端
//
// ⚠ 這張表先前整體位移一格（docs/re/14 §4 的初版），
// 於是火災畫成 16×9 的雲、暴動畫成火災的圖。
func ObjectPatternIndex(objectType, frame int) (int, bool) {
	if objectType < 0 || objectType > 3 || frame < 0 || frame >= 8 {
		return 0, false
	}
	table := [...]byte{
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x18, 0x19, 0x1A,
		0x20, 0x21, 0x22, 0x23, 0x20, 0x21, 0x22, 0x23,
		0x28, 0x29, 0x2A, 0x2B, 0x28, 0x29, 0x2A, 0x2B,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	return int(table[objectType*8+frame]), true
}

// PatternFor 依原版 object type／phase 取物件矩陣。
func (m *MCH) PatternFor(objectType, frame int) (MCHPattern, bool) {
	index, ok := ObjectPatternIndex(objectType, frame)
	if !ok {
		return MCHPattern{}, false
	}
	return m.Pattern(index)
}

// RenderPattern 把 MCH 物件矩陣合成 RGBA。每個 source byte 是一張 MCH
// 16×16 圖塊編號；0xFF 是透明格，圖塊自己的 mask 再提供像素透明度。
func (m *MCH) RenderPattern(pattern MCHPattern, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	if m == nil || p == nil || pattern.Width <= 0 || pattern.Height <= 0 ||
		len(pattern.Tiles) != pattern.Width*pattern.Height {
		return nil, fmt.Errorf("world: 無效的 MMAP.MCH 物件矩陣")
	}
	bank, err := p.Bank(bankIdx)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, pattern.Width*TileSize, pattern.Height*TileSize))
	for row := 0; row < pattern.Height; row++ {
		for col := 0; col < pattern.Width; col++ {
			id := pattern.Tiles[row*pattern.Width+col]
			if id == MCHTransparent {
				continue
			}
			tile := m.Tile(id)
			if tile == nil {
				return nil, fmt.Errorf("world: MMAP.MCH 圖塊 0x%02X 超出範圍", id)
			}
			for y := 0; y < TileSize; y++ {
				for x := 0; x < TileSize; x++ {
					colour := tile.Pix[y*TileSize+x]
					if colour == MCHTransparent {
						continue
					}
					rgba := bank[colour]
					img.SetRGBA(col*TileSize+x, row*TileSize+y, rgba)
				}
			}
		}
	}
	return img, nil
}
