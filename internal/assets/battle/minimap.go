package battle

import (
	"image"
	"image/color"
)

// TacticalMinimapWidth／Height 是 DOS/V sub_1C51E 的 64×64 格、每格
// 2×2 畫素所形成的原始縮圖尺寸。這裡保留 palette index，不在資產層
// 猜測戰術畫面專用調色盤。
const (
	TacticalMinimapWidth  = 128
	TacticalMinimapHeight = 128
)

// TacticalMinimap 是一次初始化後可重複繪製的 raw palette-index 縮圖。
// Pixels[y*128+x] 的值域是 EGA set/reset 實際使用的低 4 bit。
type TacticalMinimap struct {
	Pixels [TacticalMinimapWidth * TacticalMinimapHeight]uint8
}

// RenderTacticalMinimap 重現已由 IDA 證實的 DOS/V 戰術縮圖 producer：
//
//   BATTLE.MAP[mapY*64+mapX]
//     → BATTLE.MDL attribute[tile]
//     → 2×2 palette-index block
//
// `sub_1C51E` 的暫存器命名不是一般影像 API 的 x/y：呼叫端把
// DX=mapX、BX=mapY 傳入，routine 交換後寫到
// screenX=2*mapY、screenY=2*(63-mapX)。因此這裡保留原版的轉置與
// Y 反轉，不能改成常見的 screenX=mapX、screenY=mapY。
//
// BATTLE.MAP 實際檔案資料是 64×62；缺少的最後兩列以 tile 0 處理，
// 以保留原版 64×64 緩衝區／sub_1C83E 迴圈的完整輸出範圍。
func RenderTacticalMinimap(tiles [][]byte, attributes []byte) TacticalMinimap {
	var out TacticalMinimap
	for mapY := 0; mapY < Width; mapY++ {
		for mapX := 0; mapX < Width; mapX++ {
			tile := byte(0)
			if mapY < len(tiles) && mapX < len(tiles[mapY]) {
				tile = tiles[mapY][mapX]
			}
			attr := byte(0)
			if int(tile) < len(attributes) {
				// EGA set/reset consumes plane select bits 0–3. The upper
				// bits are retained in the source table but do not become a
				// palette index in sub_1C51E.
				attr = attributes[tile] & 0x0f
			}
			x0 := mapY * 2
			y0 := (Width - 1 - mapX) * 2
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					out.Pixels[(y0+dy)*TacticalMinimapWidth+x0+dx] = attr
				}
			}
		}
	}
	return out
}

// RGBA 將 palette index 縮圖轉成標準庫影像。palette 必須由目前已載入
// 的 GAMEPAL.BRG bank 提供；這個 helper 不改變 raw mapping。
func (m TacticalMinimap) RGBA(palette [16]color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, TacticalMinimapWidth, TacticalMinimapHeight))
	for y := 0; y < TacticalMinimapHeight; y++ {
		for x := 0; x < TacticalMinimapWidth; x++ {
			img.SetRGBA(x, y, palette[m.Pixels[y*TacticalMinimapWidth+x]&0x0f])
		}
	}
	return img
}
