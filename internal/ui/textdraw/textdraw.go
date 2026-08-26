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
// isWordRune 判斷這個字元是不是「不可從中間切開」的單字組成字元。
// 只有拉丁字母與數字算——中日文每個字都能斷。
func isWordRune(ch rune) bool {
	return ch < 0x80 && (ch >= '0' && ch <= '9' ||
		ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' ||
		ch == '\'' || ch == '-')
}

// lastSpace 回傳最後一個空白的索引，沒有回 -1。
func lastSpace(rs []rune) int {
	for i := len(rs) - 1; i >= 0; i-- {
		if rs[i] == ' ' {
			return i
		}
	}
	return -1
}

func runesWidth(rs []rune) int {
	w := 0
	for _, ch := range rs {
		w += RuneWidth(ch)
	}
	return w
}

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
			// ⭐ 英文要斷在空白，不能從字中間切開。中文每個字都能斷，
			// 拉丁字母不行——`strategist` 被切成 `strategi` ＋ `st`
			// 讀不出來。把整個未完成的單字搬到下一列
			// （單字本身就比一列長時沒得搬，照舊硬切）。
			if isWordRune(ch) {
				if k := lastSpace(current); k > 0 {
					tail := append([]rune(nil), current[k+1:]...)
					out = append(out, string(current[:k]))
					current = current[:0]
					current = append(current, tail...)
					current = append(current, ch)
					width = runesWidth(current)
					continue
				}
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

// GlyphSource 是全形字模來源：倚天（*cjk.Font，Big5）或
// HZK16（*cjk.HZK16Font，GB2312，docs/spec/84）都滿足。
type GlyphSource interface {
	Glyph(ch rune) (*image.Alpha, bool)
}

// Chain 把多份字型串起來：前面的取不到字模就換下一份。
//
// 用途是**語系的字集不會剛好等於一份字型的字集**：日文版的人名裡有
// PC-98 外字（`汜`、`瓚`、`繡`），那些字不在 JIS X 0208 裡但在倚天 Big5
// 裡有，所以 `JISKAN16 → 倚天` 這條鏈才畫得全（docs/spec/84 §2）。
// nil 會被濾掉，全空回 nil。
func Chain(sources ...GlyphSource) GlyphSource {
	var out chain
	for _, s := range sources {
		if s == nil {
			continue
		}
		// 呼叫端手上常是具體型別的 nil 指標；包進介面之後 s != nil
		// 仍成立（typed-nil），要靠各字型自己的 nil receiver 判斷擋住。
		out = append(out, s)
	}
	if len(out) == 0 {
		return nil
	}
	if len(out) == 1 {
		return out[0]
	}
	return out
}

type chain []GlyphSource

func (c chain) Glyph(ch rune) (*image.Alpha, bool) {
	for _, s := range c {
		if a, ok := s.Glyph(ch); ok {
			return a, true
		}
	}
	return nil, false
}

// Drawer 把字串畫成 Ebiten 圖片。
type Drawer struct {
	font  GlyphSource
	ascii *cjk.ASCIIFont
	cache map[cacheKey]*ebiten.Image
	// runes 是呈現前的字級替換（簡體語系的繁→簡表，docs/spec/84）。
	// 這是「同一個字選哪個字形」的層：涵蓋 Go 內 literal、人名與
	// talk fallback；已是簡體的文字過表是恆等。不改排版寬度。
	runes map[rune]rune
	// translate 是 UI 詞的語系轉換（uitext.Table.Convert）。
	translate func(string) string
}

type cacheKey struct {
	ch rune
	c  color.RGBA
}

// New 建一個 Drawer。兩個字型都可以是 nil——那樣字會全部畫成方框，
// 但程式仍然跑得起來（**沒有字型不該讓遊戲開不了**）。
// ⚠ 呼叫端若手上是具體型別的 nil 指標，**傳 nil 字面值**，
// 不要把 typed-nil 塞進介面——那會讓 Available() 誤報有字型。
func New(font GlyphSource, ascii *cjk.ASCIIFont) *Drawer {
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
	for _, ch := range d.text(s) {
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

// SetRuneMap 設定字級替換表（nil 表示不替換）。
func (d *Drawer) SetRuneMap(m map[rune]rune) {
	if d != nil {
		d.runes = m
	}
}

// SetTranslator 設定 UI 詞的翻譯函式（nil 表示不翻）。
//
// 掛在**畫的那一刻**是刻意的：介面字串散在幾十個檔案的 literal 裡，
// 一個一個包起來會漏，而每一個都會經過這裡（docs/spec/84 §2）。
// 代價是**版面已經先用原文算好了**——英文比中文窄，框會偏大不會破版；
// 真正要重算版面是第三期的事。
func (d *Drawer) SetTranslator(fn func(string) string) {
	if d != nil {
		d.translate = fn
	}
}

func (d *Drawer) text(s string) string {
	if d == nil || d.translate == nil {
		return s
	}
	return d.translate(s)
}

func (d *Drawer) glyph(ch rune, c color.RGBA) *ebiten.Image {
	if v, ok := d.runes[ch]; ok {
		ch = v
	}
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
