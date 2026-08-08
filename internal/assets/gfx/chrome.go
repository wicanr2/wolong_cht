package gfx

import (
	"fmt"
	"image"

	"github.com/wicanr2/wolong_cht/internal/assets/palette"
)

// 視窗外框的圖塊。**在 `ICONGRF` 段 3，不在段 1。**
//
// 段 3（`0x9700`，9,120 B）原本整段標「未解」，只知道它走 `sub_1F888`
// ——那支**位元對齊**的繪製常式，可以把圖畫在非 8 倍數的 x 上。
// 這正是介面元件需要的能力（視窗可以開在任意位置），
// 而地圖圖塊、頭像那些永遠對齊 8 的東西都走另一支。
// **「為什麼要用位元對齊的繪製常式」本身就指出了這一段是什麼。**
//
// 找法：拿 PC-98 實機截圖上的邊框像素，回頭在段裡逐 byte 掃相符的圖塊。
// 比的不是顏色（實機是 PC-98 調色盤，這裡用 DOS/V 的 `GAMEPAL`），
// 而是**「顏色的等價關係」**——把每個像素換成「它是這塊裡第幾種出現的顏色」，
// 得到與調色盤無關的簽章。詳見 docs/formats/03 §5.4。
//
// 尺寸是 **8×8**：畫面上量到的紅框 motif 週期正好 8 px，
// 邊柱也是 8 px 寬（黑 橙 黃 白 黃 橙 褐 黑）。
// 4 bpp × 8×8 ＝ 32 B／塊，而 9,120 ÷ 32 ＝ 285 整除。
//
// ⚠ **段 3 不是整段同尺寸的圖塊陣列**（與段 1 一樣是混合內容），
// 所以這裡只認已經對上實機畫面的那幾塊，不假裝整段都解得開。
const (
	// ChromeTile 是外框圖塊的邊長。
	ChromeTile = 8

	// ChromeEdge 是上下邊那條「紅框綠心」的 motif。
	ChromeEdge = 0x06C0
	// ChromeCap 是左右邊柱最上面那一塊（金色柱頭）。
	ChromeCap = 0x06E0
	// ChromeShaft 是邊柱的柱身，往下重複貼。
	// 段內 `0x0700`–`0x0720` 每 8 byte 都相符 ——
	// 那是因為柱身在垂直方向同紋，不是有五塊不同的圖。
	ChromeShaft = 0x0700
)

// chromeBytes 是一塊 8×8 4bpp 的大小。
const chromeBytes = ChromeTile * ChromeTile * Planes / 8

// DecodeChrome 解段 3 裡位於 off 的那一塊 8×8，回傳調色盤索引。
//
// seg 要傳**段 3 本身**（`IconRegions[3]` 那一段），不是整個 `ICONGRF.DAT`。
func DecodeChrome(seg []byte, off int) ([]byte, error) {
	if off < 0 || off+chromeBytes > len(seg) {
		return nil, fmt.Errorf("外框圖塊偏移 0x%X 超出段 3（%d byte）", off, len(seg))
	}
	out := make([]byte, ChromeTile*ChromeTile)
	for y := 0; y < ChromeTile; y++ {
		for x := 0; x < ChromeTile; x++ {
			var v byte
			for pl := 0; pl < Planes; pl++ {
				b := seg[off+pl*ChromeTile+y]
				v |= ((b >> (7 - uint(x))) & 1) << uint(pl)
			}
			out[y*ChromeTile+x] = v
		}
	}
	return out, nil
}

// RenderChrome 解一塊外框圖塊並上色。
func RenderChrome(seg []byte, off int, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	idx, err := DecodeChrome(seg, off)
	if err != nil {
		return nil, err
	}
	bank, err := p.Bank(bankIdx)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, ChromeTile, ChromeTile))
	for y := 0; y < ChromeTile; y++ {
		for x := 0; x < ChromeTile; x++ {
			img.SetRGBA(x, y, bank[idx[y*ChromeTile+x]])
		}
	}
	return img, nil
}
