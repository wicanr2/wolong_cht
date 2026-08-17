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
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

// Tile 是外框圖塊的邊長（8）。視窗的座標與尺寸最好是它的倍數，
// 不然邊會切在半塊上。
const Tile = gfx.ChromeTile

// 介面顏色的調色盤索引（`GAMEPAL.BRG`，規格 docs/spec/54）。
//
// ⭐ **不要手抄 RGB。** 這五個顏色以前是抄實機截圖抄下來的常數，
// 而 docs/spec/51 把調色盤換算改成走 VGA 的 6 bit DAC 之後，
// 抄下來的值全部差了 2–4——**解碼修好了，常數不會跟著修**。
const (
	MenuIndex   = 8  // 選單／情報視窗的深藍底
	SheetIndex  = 9  // 清單視窗的米色底
	SelectIndex = 5  // 反白條
	InkIndex    = 0  // 米色底上的字（也是命令列的底）
	PaperIndex  = 15 // 深藍底上的字
)

// 內部底色。原版的選單視窗是深藍底 ＋ 龍紋，清單視窗是米色底。
// **`Load` 會照上面的索引從 `GAMEPAL.BRG` 覆寫這幾個值**；
// 這裡的初值只是沒有素材時的 fallback。
var (
	// Menu 是選單／情報視窗的底色。
	Menu = color.RGBA{0, 32, 97, 255}
	// Sheet 是清單視窗的底色。
	Sheet = color.RGBA{243, 211, 146, 255}
	// Select 是反白條的顏色。
	Select = color.RGBA{81, 146, 65, 255}
	// Ink 是米色底上的字色，**也是命令列的底色**（docs/spec/54 §2）。
	Ink = color.RGBA{0, 0, 0, 255}
	// Blank 是「純色平塗」的底。值與 Ink 相同（都是色 0），
	// 名字分開是為了讀得出意圖：Ink 是字色、Blank 是底色。
	// 它**不等於 Menu**，所以 fillInterior 不會鋪龍紋——命令列用這一個。
	Blank = color.RGBA{0, 0, 0, 255}
	// Paper 是深藍底上的字色。
	Paper = color.RGBA{243, 243, 243, 255}

	fallbackEdge = color.RGBA{200, 40, 40, 255}
	fallbackCap  = color.RGBA{240, 200, 80, 255}
)

// Set 是畫外框要用的三塊圖，外加視窗內部的龍紋。
// Load 失敗時全是 nil，畫成純色框。
type Set struct {
	edge    *ebiten.Image // 上下邊
	cap     *ebiten.Image // 柱頭
	shaft   *ebiten.Image // 柱身
	texture *ebiten.Image // 內部的龍紋，32×32
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
	// 顏色跟著調色盤走，不要手抄（docs/spec/54）。取不到就留 fallback。
	for _, c := range []struct {
		dst *color.RGBA
		idx int
	}{{&Menu, MenuIndex}, {&Sheet, SheetIndex}, {&Select, SelectIndex},
		{&Ink, InkIndex}, {&Blank, InkIndex}, {&Paper, PaperIndex}} {
		if col, err := lib.PaletteColor(bank, c.idx); err == nil {
			*c.dst = col
		}
	}
	s.edge, s.cap, s.shaft = get(gfx.ChromeEdge), get(gfx.ChromeCap), get(gfx.ChromeShaft)
	if img, err := lib.RenderWindowTexture(bank); err == nil {
		s.texture = ebiten.NewImageFromImage(img)
	}
	return s
}

// TextureSize 是龍紋磚塊的邊長。**底紋是螢幕對齊的**，
// 所以它同時是水平與垂直的週期（docs/formats/03 §5.5）。
const TextureSize = gfx.WindowTextureSize

// fillInterior 鋪視窗底色。深藍的選單視窗鋪的是**龍紋**，不是純色。
//
// ⭐ **磚塊釘在螢幕上，不是釘在視窗左上角**：
// 像素 (x, y) 取 tile[y mod 32][x mod 32]。
// 先前幾輪都從視窗角落開始平鋪，所以怎麼試都對不上實機
// （docs/formats/03 §5.5）。
func (s *Set) fillInterior(dst *ebiten.Image, x, y, w, h int, fill color.RGBA) {
	if s == nil || s.texture == nil || fill != Menu || w <= 0 || h <= 0 {
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h),
			fill, false)
		return
	}
	clip, ok := dst.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
	if !ok {
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h),
			fill, false)
		return
	}
	floorTo := func(v int) int {
		if v < 0 {
			return -((-v + TextureSize - 1) / TextureSize) * TextureSize
		}
		return v / TextureSize * TextureSize
	}
	for ty := floorTo(y); ty < y+h; ty += TextureSize {
		for tx := floorTo(x); tx < x+w; tx += TextureSize {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tx), float64(ty))
			clip.DrawImage(s.texture, op)
		}
	}
}

// Available 回報有沒有真的拿到原版圖塊。
func (s *Set) Available() bool { return s != nil && s.edge != nil }

// Window 在 (x, y) 畫一個 w×h 的視窗：先鋪底色，再沿四邊貼圖塊。
//
// x、y、w、h 用像素。邊框**畫在矩形內側**，所以內容區是
// (x+Tile, y+Tile) 到 (x+w-Tile, y+h-Tile)。
func (s *Set) Window(dst *ebiten.Image, x, y, w, h int, fill color.RGBA) {
	s.fillInterior(dst, x, y, w, h, fill)
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
