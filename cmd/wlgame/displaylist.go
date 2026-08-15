package main

// 顯示清單那幾個 opcode 的共用畫法（docs/re/48 §2.1）。
//
// ⭐ 記錄的第六個 word **不是保留欄位**，是顏色：
//   - `op 03`（填矩形）用它的**低 byte**當填色。
//   - `op 04`（立體框）用整個 word：低 byte 畫**上邊與左邊**、
//     高 byte 畫**下邊與右邊**（`sub_1F465` 每條線之前 `xchg ah,al`），
//     立體感就是這樣做出來的。
//
// 十個場景只用了兩套配色，剛好對應兩種控制項。

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// 顯示清單實際用到的調色盤索引。
const (
	// 凹槽：可輸入／可選的欄位（四槽視窗的日期欄、ＹＥＳ／ＮＯ 的兩格）。
	dlSunkenFill              = 0x05
	dlSunkenOuterLight        = 0x02
	dlSunkenOuterDark         = 0x00
	dlSunkenInnerLight        = 0x0D
	dlSunkenInnerDark         = 0x04

	// 按鈕：「確定」「自定」「重來」「繼續」。
	dlButtonFill  = 0x07
	dlButtonLight = 0x09
	dlButtonDark  = 0x06

	// ⭐ 按鈕上的字是**黑的**（`op 08` 的 arg2 ＝ 0003，docs/re/55 §3）。
	// 先前四個視窗都畫成白字。
	dlButtonText = 0x00

	// 系統選單值格的「 ＯＫ 」：黑字綠底（arg2 ＝ 5001）。
	dlValueFill = 0x05
	dlValueText = 0x00
)

// 取不到原版調色盤時的替代色。只在沒有素材的環境（測試、無 lib）出現。
var (
	dlTextFallback   = color.RGBA{20, 20, 20, 255}
	dlSunkenFallback = color.RGBA{40, 40, 60, 255}
	dlButtonFallback = color.RGBA{90, 90, 90, 255}
	dlLightFallback  = color.RGBA{210, 210, 210, 255}
	dlDarkFallback   = color.RGBA{60, 60, 60, 255}
)

// dlFill 是 op 03：一塊實心矩形。
func (g *game) dlFill(dst *ebiten.Image, x, y, w, h, index int, fallback color.RGBA) {
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h),
		g.paletteInk(index, fallback), false)
}

// dlBevel 是 op 04 的一圈：上邊與左邊用 light，下邊與右邊用 dark。
func (g *game) dlBevel(dst *ebiten.Image, x, y, w, h, light, dark int) {
	l := g.paletteInk(light, dlLightFallback)
	d := g.paletteInk(dark, dlDarkFallback)
	fx, fy, fw, fh := float32(x), float32(y), float32(w), float32(h)
	vector.DrawFilledRect(dst, fx, fy, fw, 1, l, false)      // 上
	vector.DrawFilledRect(dst, fx, fy, 1, fh, l, false)      // 左
	vector.DrawFilledRect(dst, fx, fy+fh-1, fw, 1, d, false) // 下
	vector.DrawFilledRect(dst, fx+fw-1, fy, 1, fh, d, false) // 右
}

// dlSunken 畫一個凹槽：底色 ＋ 兩圈立體邊。
//
// (x, y, w, h) 是**底色那一塊**；兩圈框畫在它外面，與原版一樣
// （場景 6 的日期欄底 (152,56) 120×16，兩筆 op 04 分別從 −2 與 −1 開始）。
func (g *game) dlSunken(dst *ebiten.Image, x, y, w, h int) {
	g.dlFill(dst, x, y, w, h, dlSunkenFill, dlSunkenFallback)
	g.dlBevel(dst, x-2, y-2, w+4, h+4, dlSunkenOuterLight, dlSunkenOuterDark)
	g.dlBevel(dst, x-1, y-1, w+2, h+2, dlSunkenInnerLight, dlSunkenInnerDark)
}

// dlButton 畫一顆按鈕：底色 7 ＋ 一圈 9／6。
func (g *game) dlButton(dst *ebiten.Image, x, y, w, h int) {
	g.dlFill(dst, x, y, w, h, dlButtonFill, dlButtonFallback)
	g.dlBevel(dst, x-1, y-1, w+2, h+2, dlButtonLight, dlButtonDark)
}

// dlButtonInk 是按鈕上那行字的顏色。
func (g *game) dlButtonInk() color.RGBA { return g.paletteInk(dlButtonText, dlTextFallback) }

// dlValueBox 是系統選單那種值格：**只有兩圈框沒有底**，
// 底色由「 ＯＫ 」那個字串自己帶（docs/re/55 §2）。
func (g *game) dlValueBox(dst *ebiten.Image, x, y, w, h int) {
	g.dlBevel(dst, x-2, y-2, w+4, h+4, dlSunkenOuterLight, dlSunkenOuterDark)
	g.dlBevel(dst, x-1, y-1, w+2, h+2, dlSunkenInnerLight, dlSunkenInnerDark)
}
