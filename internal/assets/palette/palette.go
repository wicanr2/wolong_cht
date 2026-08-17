// Package palette 解 `.BRG` 調色盤。
//
// 規格：docs/formats/02-brg-palette.md（READY）
// 出處：docs/re/02-palette-routine.md（兩版 KI.EXE 的機器碼）
//
// 這一層是純解碼，不認識 Ebiten，回傳標準函式庫的型別。
package palette

import (
	"fmt"
	"image/color"
)

const (
	// BankSize 是一組的色數。原版一次只把一組推進硬體
	// （`cmp ch, 10h` / `jb`）。
	BankSize = 16

	// FullBrightness 是亮度的滿值。原版的淡入從 0 走到這裡，
	// 淡出反過來。**不是 255** —— 呼叫端寫的是 `mov cl, 10h`。
	FullBrightness = 16

	bytesPerColor = 3
)

// Palette 是一個 `.BRG` 檔解出來的全部色組。
type Palette struct {
	// Banks[i] 是第 i 組的 16 色。
	Banks [][BankSize]color.RGBA
}

// Scale 是原版的亮度縮放，回傳 0–15 的通道值。
//
//	al = v << 4 ; ax = al * bl ; ax += 0x80 ; 取 ah
//
// PC-98 把結果直接寫進 4 bit 的類比調色盤；DOS/V 再左移 2 位變成
// VGA DAC 的 6 bit。兩者是同一個顏色，只是硬體位寬不同。
func Scale(v byte, brightness int) byte {
	return byte(((int(v) << 4 * brightness) + 0x80) >> 8)
}

// toSRGB 把 4 bit 通道轉成 8 bit，走的是 **DOS/V 的硬體路徑**。
//
// 規格 docs/spec/51。原版在 DOS/V 多做兩次 `shl ah,1` 才寫進 VGA DAC
// （docs/re/02 §5），所以 DAC 值是 `4v`——**上限 60，不是 63**。
// 類比輸出因此永遠到不了滿刻度：色號 15 出來是 `#F3`，不是 `#FF`。
//
// 6 → 8 bit 用 VGA 慣用的位元複製。這不是美感選擇：拿松崗 DOS/V 的
// 實機截圖比對，`GAMEPAL.BRG` 春組 16 色**整組**都出現在畫面上，
// 而舊的 `v*255/15` 只中黑色一色（tools/palette_dacscale.py）。
//
// ⚠ PC-98 的類比調色盤本來就是 4 bit／通道，那一版 15 就是滿刻度。
// **兩版的顏色本來就不一樣**，這條只管 DOS/V。
func toSRGB(v byte, brightness int) uint8 {
	dac := Scale(v, brightness) << 2 // 4 bit → VGA DAC 的 6 bit（0–60）
	return uint8(dac<<2 | dac>>4)
}

// Parse 解一個 `.BRG` 檔。每色 3 byte，順序是 B、R、G
// —— 副檔名就是通道順序。
func Parse(data []byte) (*Palette, error) {
	if len(data) == 0 || len(data)%(BankSize*bytesPerColor) != 0 {
		return nil, fmt.Errorf("palette: %d byte 不是 %d 的倍數（每組 %d 色）",
			len(data), BankSize*bytesPerColor, BankSize)
	}
	p := &Palette{}
	for off := 0; off < len(data); off += BankSize * bytesPerColor {
		var bank [BankSize]color.RGBA
		for i := range bank {
			c := data[off+i*bytesPerColor:]
			bank[i] = color.RGBA{
				R: toSRGB(c[1], FullBrightness),
				G: toSRGB(c[2], FullBrightness),
				B: toSRGB(c[0], FullBrightness),
				A: 0xff,
			}
		}
		p.Banks = append(p.Banks, bank)
	}
	return p, nil
}

// Bank 取第 i 組。超出範圍回傳 error，不要靜悄悄回第 0 組
// —— 那會讓「選錯組」變成看不出來的顏色錯誤。
func (p *Palette) Bank(i int) ([BankSize]color.RGBA, error) {
	if i < 0 || i >= len(p.Banks) {
		return [BankSize]color.RGBA{},
			fmt.Errorf("palette: 第 %d 組不存在（共 %d 組）", i, len(p.Banks))
	}
	return p.Banks[i], nil
}

// Season 是 GAMEPAL.BRG 前四組的語意：整個季節效果只換色號 14
// （地表植被），其餘 15 色四季共用。見 docs/formats/02 §4。
type Season int

const (
	Spring Season = iota // 色 14 = #88aa66 灰綠
	Summer               // #55aa11 鮮綠
	Autumn               // #dd8800 橙褐
	Winter               // #ffffff 雪白
)
