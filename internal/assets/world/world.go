// Package world 解大地圖 `MMAP.*`。
//
// 規格：docs/formats/06-mmap-rle.md
// 出處：docs/re/04（`sub_1E48A` / `sub_1E4CE` / `sub_1F5E7`）
package world

import (
	"fmt"
	"image"

	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/assets/palette"
	"github.com/wicanr2/wolong_cht/internal/assets/rle"
)

const (
	// Width、Height 是世界地圖的格數。
	// 出處：sub_1E4CE 的 `cmp cx, 180h`（一列 384 格）與
	// `inc ah / jnz`（ah 從 0 跑滿一圈 = 256 列）。
	Width  = 384
	Height = 256

	// TileSize 是一格的邊長（像素）。
	TileSize = 16

	// TileCount 是 MMAP.MDL 裡的圖塊數：32768 / 128。
	TileCount = 256
)

// Map 是解出來的世界地圖。
type Map struct {
	// Tiles 是 Height×Width 的圖塊編號，逐列排列。
	Tiles []byte
	// Extra 是 RLE 解出來超過 Width*Height 的尾巴。
	// 原版的解壓器不知道目標長度，一路解到檔尾，所以會多出幾個 byte。
	// 保留而不丟棄，是為了寫回時能還原（未解區域一個 byte 都不動）。
	Extra []byte
}

// ParseMap 解 `MMAP.MAP`（RLE 壓縮的圖塊編號表）。
func ParseMap(data []byte) (*Map, error) {
	raw := rle.Decode(data)
	if len(raw) < Width*Height {
		return nil, fmt.Errorf("world: 解壓後只有 %d B，不足 %d B（%d×%d）",
			len(raw), Width*Height, Width, Height)
	}
	return &Map{Tiles: raw[:Width*Height], Extra: raw[Width*Height:]}, nil
}

// Tile 取 (x, y) 的圖塊編號。
func (m *Map) Tile(x, y int) (byte, error) {
	if x < 0 || x >= Width || y < 0 || y >= Height {
		return 0, fmt.Errorf("world: (%d,%d) 超出 %d×%d", x, y, Width, Height)
	}
	return m.Tiles[y*Width+x], nil
}

// TileSet 是 `MMAP.MDL` 的 256 塊地形圖塊。
type TileSet struct {
	spec gfx.Spec
	data []byte
}

// ParseTileSet 解 `MMAP.MDL`。
func ParseTileSet(data []byte) (*TileSet, error) {
	spec := gfx.Spec{Name: "MMAP.MDL", Width: TileSize, Height: TileSize}
	n, rem := spec.Count(data)
	if rem != 0 || n != TileCount {
		return nil, fmt.Errorf("world: MMAP.MDL 解出 %d 塊餘 %d B，預期 %d 塊餘 0",
			n, rem, TileCount)
	}
	return &TileSet{spec: spec, data: data}, nil
}

// Render 畫出以 (x0, y0) 為左上角、cols×rows 格的一塊地圖。
//
// ⚠ **這裡還沒做自動連接。** 原版的道路與河流會依鄰格換圖塊
// （`sub_1E57F`，值域 0xB8–0xDD，見 docs/formats/05 §3），
// 直接查表畫出來的接邊會是斷的。等 `sub_1E68C` 讀完再補。
func (m *Map) Render(ts *TileSet, pal *palette.Palette, bank,
	x0, y0, cols, rows int) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, cols*TileSize, rows*TileSize))
	for ry := 0; ry < rows; ry++ {
		for rx := 0; rx < cols; rx++ {
			t, err := m.Tile(x0+rx, y0+ry)
			if err != nil {
				return nil, err
			}
			tile, err := ts.spec.RenderRGBA(ts.data, int(t), pal, bank)
			if err != nil {
				return nil, err
			}
			for py := 0; py < TileSize; py++ {
				for px := 0; px < TileSize; px++ {
					img.SetRGBA(rx*TileSize+px, ry*TileSize+py,
						tile.RGBAAt(px, py))
				}
			}
		}
	}
	return img, nil
}
