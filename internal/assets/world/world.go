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

// MapHeader 是**壓縮檔**開頭那 4 byte：小端 u32，值就是地圖本體的長度
// （`00 80 01 00` ＝ 98,304 ＝ 384×256）。原版的載入器 `LSEEK` 跳過它
// 才開始解壓（`sub_1F655`，docs/spec/113）。
//
// ⭐ **它是頭不是尾。** 把它當成圖塊讀，整張地圖會往左移四格——
// 而那四格會以「據點中心在記錄座標 +4」「鏡頭比看到的那一欄小 4」的
// 形式散到各處，看起來像兩個獨立的怪癖（docs/formats/05 §2.1）。
const MapHeader = 4

// Map 是解出來的世界地圖。
type Map struct {
	// Tiles 是 Height×Width 的圖塊編號，逐列排列。
	Tiles []byte
	// Header 是壓縮檔前面那 4 byte（見 MapHeader）。寫回時要原樣放回去。
	Header []byte
	// Extra 是解出來超過本體的部分。**兩版都是 0 B**——檔頭宣告的長度
	// 就是 384×256，解壓器解到那裡剛好結束。留著這個欄位是為了寫回時
	// 不丟掉任何 byte，以及讓「多出來了」變成看得見的異常。
	Extra []byte
}

// ParseMap 解 `MMAP.MAP`（RLE 壓縮的圖塊編號表）。
func ParseMap(data []byte) (*Map, error) {
	if len(data) < MapHeader {
		return nil, fmt.Errorf("world: 只有 %d B，連 %d B 的長度頭都不夠",
			len(data), MapHeader)
	}
	raw, err := rle.DecodeFile(data)
	if err != nil {
		return nil, fmt.Errorf("world: %w", err)
	}
	if len(raw) < Width*Height {
		return nil, fmt.Errorf("world: 解壓後只有 %d B，不足 %d×%d",
			len(raw), Width, Height)
	}
	return &Map{
		Header: data[:MapHeader],
		Tiles:  raw[:Width*Height],
		Extra:  raw[Width*Height:],
	}, nil
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
// ⚠ **這裡畫的是檔案裡的原始圖塊值，不是原版執行時的地圖內容。**
// 原版載入後有一個前處理階段（`sub_1E4CE`）會把道路格（0xCB–0xD3）
// 換成節點流水號並建一張連接關係表，見 docs/formats/05 §3。
// 那一步還沒實作——目前的輸出對「看地形」是對的，
// 對「重現原版的記憶體狀態」還不夠。
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
