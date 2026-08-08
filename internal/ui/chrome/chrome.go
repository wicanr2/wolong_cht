// Package chrome 畫原版樣式的視窗外框。
//
// 外框不是自己設計的，是原版的美術：`ICONGRF` 段 3 裡的三塊 8×8
// （見 `internal/assets/gfx/chrome.go` 的出處與找法）。
//
//	上下邊　　紅框綠心的 motif，每 8 px 一個
//	左右邊　　金色圓柱，最上面一塊是柱頭
//	內部　　　深藍底
//
// ⚠ 原版的內部還有龍形花紋。**那不是可重複貼的圖塊**——
// 量過任何尺度都沒有週期（`docs/formats/03` §5.5），是每個視窗一張的
// 裝飾底圖。要補得先反組譯視窗繪製常式，不是再試一次拼貼尺寸。
//
// 沒有原版素材時退回純色框，**不會整個畫不出來**——
// 素材是玩家自備的，缺了要能降級跑。
package chrome

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

// Tile 是外框圖塊的邊長（8）。視窗的座標與尺寸最好是它的倍數，
// 不然邊會切在半塊上。
const Tile = gfx.ChromeTile

// 內部底色。原版的選單視窗是深藍底 ＋ 龍紋，清單視窗是米色底。
var (
	// Menu 是選單／情報視窗的底色（原版量到的 (0,32,101)）。
	Menu = color.RGBA{0, 32, 101, 255}
	// Sheet 是清單視窗的底色（原版量到的 (255,223,154)）。
	Sheet = color.RGBA{255, 223, 154, 255}
	// Select 是反白條的顏色（原版量到的 (85,154,69)）。
	Select = color.RGBA{85, 154, 69, 255}
	// Ink 是米色底上的字色。
	Ink = color.RGBA{0, 0, 0, 255}
	// Paper 是深藍底上的字色。
	Paper = color.RGBA{255, 255, 255, 255}

	fallbackEdge = color.RGBA{200, 40, 40, 255}
	fallbackCap  = color.RGBA{240, 200, 80, 255}
)

// Set 是畫外框要用的三塊圖。Load 失敗時三塊都是 nil，畫成純色框。
type Set struct {
	edge  *ebiten.Image // 上下邊
	cap   *ebiten.Image // 柱頭
	shaft *ebiten.Image // 柱身
}

// Load 從素材庫取出外框圖塊。**缺素材不算錯**——回傳的 Set 仍可用，
// 只是畫成純色框。bank 是調色盤組號（跟著季節走）。
func Load(lib *library.Library, bank int) *Set {
	s := &Set{}
	if lib == nil {
		return s
	}
	get := func(off int) *ebiten.Image {
		img, err := lib.RenderChrome(off, bank)
		if err != nil {
			return nil
		}
		return ebiten.NewImageFromImage(img)
	}
	s.edge, s.cap, s.shaft = get(gfx.ChromeEdge), get(gfx.ChromeCap), get(gfx.ChromeShaft)
	return s
}

// Available 回報有沒有真的拿到原版圖塊。
func (s *Set) Available() bool { return s != nil && s.edge != nil }

// Window 在 (x, y) 畫一個 w×h 的視窗：先鋪底色，再沿四邊貼圖塊。
//
// x、y、w、h 用像素。邊框**畫在矩形內側**，所以內容區是
// (x+Tile, y+Tile) 到 (x+w-Tile, y+h-Tile)。
func (s *Set) Window(dst *ebiten.Image, x, y, w, h int, fill color.RGBA) {
	vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h), fill, false)
	if !s.Available() {
		vector.StrokeRect(dst, float32(x), float32(y), float32(w), float32(h),
			2, fallbackCap, false)
		vector.StrokeRect(dst, float32(x)+2, float32(y)+2, float32(w)-4, float32(h)-4,
			1, fallbackEdge, false)
		return
	}
	put := func(img *ebiten.Image, px, py int) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(px), float64(py))
		dst.DrawImage(img, op)
	}
	// 上下邊。從左柱右邊開始鋪到右柱左邊，中間不留缺口。
	for px := x + Tile; px+Tile <= x+w-Tile; px += Tile {
		put(s.edge, px, y)
		put(s.edge, px, y+h-Tile)
	}
	// 左右邊。最上面一塊是柱頭，其餘是柱身。
	put(s.cap, x, y)
	put(s.cap, x+w-Tile, y)
	for py := y + Tile; py+Tile <= y+h-Tile; py += Tile {
		put(s.shaft, x, py)
		put(s.shaft, x+w-Tile, py)
	}
	put(s.cap, x, y+h-Tile)
	put(s.cap, x+w-Tile, y+h-Tile)
}
