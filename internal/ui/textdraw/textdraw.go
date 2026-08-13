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

// RuneWidth 回傳單一字元在這套點陣字上的像素寬度。
// ASCII 走 8×15 半形字；其他 Unicode 字元走 16×15 全形字。TALK
// 排版與 Drawer 必須共用這個契約，否則「量過不溢位」只會在其中一條
// 路徑成立。
func RuneWidth(ch rune) int {
	if ch < 0x80 {
		return HalfW
	}
	return GlyphW
}

// StringWidth 回傳一列字串的實際繪製寬度（像素）。
func StringWidth(s string) int {
	w := 0
	for _, ch := range s {
		if ch == '\n' || ch == '\r' {
			continue
		}
		w += RuneWidth(ch)
	}
	return w
}

// WrapLines 依據實際點陣字寬度斷行，但保留呼叫端傳入的硬斷行。
//
// 這是呈現層的 formatter，不改 TALK.DAT，也不把換行寫回規則層。每個
// 空字串仍保留成一列；換行時捨棄剛好位於斷點的前導空白，避免英文／數字
// 混排在下一列留下難以察覺的縮排。
func WrapLines(lines []string, maxPixels int) []string {
	if len(lines) == 0 {
		return nil
	}
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, WrapLine(line, maxPixels)...)
	}
	return out
}

// WrapLine 對一條 TALK 硬斷行做測量式換行。關閉標點不會被單獨放到
// 下一行；若它正好超出邊界，寧可讓該列多出一個全形字，也不讓中文標點
// 出現在列首。maxPixels<=0 時維持原列，作為 fail-safe。
func WrapLine(line string, maxPixels int) []string {
	if maxPixels <= 0 {
		return []string{line}
	}
	if line == "" {
		return []string{""}
	}

	var out []string
	current := make([]rune, 0, len([]rune(line)))
	width := 0
	flush := func() {
		out = append(out, string(current))
		current = current[:0]
		width = 0
	}
	for _, ch := range line {
		if ch == '\r' {
			continue
		}
		if ch == '\n' {
			flush()
			continue
		}
		if ch == ' ' && len(current) == 0 {
			continue
		}
		w := RuneWidth(ch)
		if width > 0 && width+w > maxPixels {
			if isClosingPunctuation(ch) {
				current = append(current, ch)
				width += w
				continue
			}
			flush()
			if ch == ' ' {
				continue
			}
		}
		current = append(current, ch)
		width += w
	}
	if len(current) > 0 || len(out) == 0 {
		flush()
	}
	return out
}

func isClosingPunctuation(ch rune) bool {
	switch ch {
	case ',', '.', '!', '?', ':', ';', ')', ']', '}', '%',
		'，', '。', '、', '！', '？', '：', '；', '）', '］', '｝',
		'》', '」', '』', '】', '〕', '〉', '”', '’':
		return true
	default:
		return false
	}
}

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
	return StringWidth(s)
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
		x += RuneWidth(ch)
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
