// Package textdraw 把倚天點陣字畫到 Ebiten 的畫面上。
//
// 為什麼不用 Ebiten 內建的除錯字型：它只有 ASCII，中文會被**靜靜吃掉**，
// 畫面上看起來像排版 bug，很難查。這一層一律走倚天 16×15 點陣，
// 取不到字模時畫一個明顯的方框，**寧可醜也不要不見**。
//
// 字型檔不隨本專案散布，執行時由使用者以路徑指定
// （與原版資料同一個處理方式）。
package textdraw

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/assets/cjk"
)

// 半形字用內建的 8×15 點陣，與倚天的 16×15 同高，混排才不會參差。
const (
	GlyphW  = cjk.GlyphWidth  // 16
	GlyphH  = cjk.GlyphHeight // 15
	HalfW   = GlyphW / 2      // 8
	LineGap = 2
)

// Drawer 把字串畫成 Ebiten 圖片。
type Drawer struct {
	font  *cjk.Font
	ascii *cjk.ASCIIFont
	cache map[cacheKey]*ebiten.Image
}

type cacheKey struct {
	ch rune
	c  color.RGBA
}

// New 建一個 Drawer。兩個字型都可以是 nil——那樣字會全部畫成方框，
// 但程式仍然跑得起來（**沒有字型不該讓遊戲開不了**）。
func New(font *cjk.Font, ascii *cjk.ASCIIFont) *Drawer {
	return &Drawer{font: font, ascii: ascii, cache: map[cacheKey]*ebiten.Image{}}
}

// Available 回報有沒有載到全形字型。
func (d *Drawer) Available() bool { return d != nil && d.font != nil }

// Width 回傳一段字串畫出來會佔多寬（像素）。
func (d *Drawer) Width(s string) int {
	w := 0
	for _, ch := range s {
		if ch < 0x80 {
			w += HalfW
		} else {
			w += GlyphW
		}
	}
	return w
}

// Draw 從 (x, y) 開始畫一段字串，回傳結束時的 x。
// y 是**字的上緣**，不是基線——點陣字沒有基線的概念。
func (d *Drawer) Draw(dst *ebiten.Image, s string, x, y int, c color.RGBA) int {
	for _, ch := range s {
		if ch == '\n' {
			continue // 換行由呼叫端處理，這裡只畫一列
		}
		img := d.glyph(ch, c)
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			dst.DrawImage(img, op)
		}
		if ch < 0x80 {
			x += HalfW
		} else {
			x += GlyphW
		}
	}
	return x
}

// DrawLines 畫多列，列距是 GlyphH + LineGap。
func (d *Drawer) DrawLines(dst *ebiten.Image, lines []string, x, y int, c color.RGBA) {
	for _, ln := range lines {
		d.Draw(dst, ln, x, y, c)
		y += GlyphH + LineGap
	}
}

func (d *Drawer) glyph(ch rune, c color.RGBA) *ebiten.Image {
	k := cacheKey{ch, c}
	if img, ok := d.cache[k]; ok {
		return img
	}
	img := d.render(ch, c)
	d.cache[k] = img
	return img
}

func (d *Drawer) render(ch rune, c color.RGBA) *ebiten.Image {
	if ch == ' ' {
		return nil
	}
	if d.font != nil {
		if a, ok := d.font.Glyph(ch); ok {
			return tint(a, c)
		}
	}
	if ch < 0x80 && d.ascii != nil {
		if a, ok := d.ascii.Glyph(ch); ok {
			return tint(a, c)
		}
	}
	// 取不到字模：畫一個空心方框。**不要什麼都不畫**——
	// 缺字要看得出來，否則會被誤判成排版問題查很久。
	return missingBox(ch, c)
}

func tint(a *image.Alpha, c color.RGBA) *ebiten.Image {
	b := a.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if a.AlphaAt(x, y).A != 0 {
				rgba.SetRGBA(x, y, c)
			}
		}
	}
	return ebiten.NewImageFromImage(rgba)
}

func missingBox(ch rune, c color.RGBA) *ebiten.Image {
	w := GlyphW
	if ch < 0x80 {
		w = HalfW
	}
	rgba := image.NewRGBA(image.Rect(0, 0, w, GlyphH))
	for x := 0; x < w; x++ {
		rgba.SetRGBA(x, 0, c)
		rgba.SetRGBA(x, GlyphH-1, c)
	}
	for y := 0; y < GlyphH; y++ {
		rgba.SetRGBA(0, y, c)
		rgba.SetRGBA(w-1, y, c)
	}
	return ebiten.NewImageFromImage(rgba)
}
