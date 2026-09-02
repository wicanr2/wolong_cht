// Package cutscene 解 `OPEN_S*.DAT`／`END_S*.DAT` 這一批過場畫面。
//
// 格式見 `docs/formats/09-cutscene-images.md`。摘要：
//
//	外層   4 byte 長度頭 ＋ 與 `MMAP.MAP` 同一種 RLE（`internal/assets/rle`）
//	內層   640 × 400、16 色、EGA 四平面
//	版面   上下**兩半各 200 列**；一半之內是**平面優先**（4 × 16,000 B）
//	色盤   `ENDPAL.BRG`／`OPENPAL.BRG` 的第 n 組 ＝ 第 n 幕
package cutscene

import (
	"fmt"
	"image"
	"image/color"

	"github.com/wicanr2/wolong_cht/internal/assets/rle"
)

// 一張過場畫面的尺寸與版面（`sub_1178A` 的兩次呼叫，docs/re/70 §5）。
const (
	Width  = 640
	Height = 400
	// Planes 是 EGA 的四個位元平面；色號 ＝ 四個平面同一個位元湊出來的 0–15。
	Planes = 4
	// Stride 是一列一個平面佔幾個 byte。
	Stride = Width / 8
	// HalfRows 是「一半」的列數。原版分兩次貼，第二次的目的位址是
	// `0x3E80` ＝ 200 × 80，所以切點在第 200 列。
	HalfRows = Height / 2
	// PlaneHalf 是一半之內一個平面的長度。
	PlaneHalf = Stride * HalfRows
	// HalfSize 是一半的長度（四個平面）。
	HalfSize = PlaneHalf * Planes
	// Size 是一張完整畫面解壓後的長度。
	Size = HalfSize * 2
)

// Decode 把 `.DAT` 解壓成畫面 buffer。
//
// ⭐ 走 `rle.DecodeFile`：檔案前 4 byte 是解壓長度，原版的載入器
// `LSEEK` 跳過它才開始解（docs/spec/113）。**解出來一定等於宣告值**——
// 從 offset 0 解會掉相位，畫面整體位移。
//
// 回傳長度是檔頭宣告的值，不一定等於 `Size`：`END_S1` 宣告 91,200 B
// （那一幕的版面本來就不是整張 640 × 400，見 `FirstScenePixels`），
// 而開場的 `OPEN_S2`–`S4` 是 384,000 B ＝ 12 幀 320 × 200。
// **比 `Size` 短才補 0，長的一律原樣回傳**——截掉就只剩第一幀。
func Decode(src []byte) ([]byte, error) {
	out, err := rle.DecodeFile(src)
	if err != nil {
		return nil, err
	}
	if len(out) >= Size {
		return out, nil
	}
	buf := make([]byte, Size)
	copy(buf, out)
	return buf, nil
}

// Pixels 回傳每個像素的色號（0–15），逐列排列。
func Pixels(buf []byte) []byte {
	px := make([]byte, Width*Height)
	for y := 0; y < Height; y++ {
		base := 0
		row := y
		if y >= HalfRows {
			base, row = HalfSize, y-HalfRows
		}
		for p := 0; p < Planes; p++ {
			off := base + p*PlaneHalf + row*Stride
			for b := 0; b < Stride; b++ {
				v := buf[off+b]
				if v == 0 {
					continue
				}
				for bit := 0; bit < 8; bit++ {
					if v&(0x80>>bit) != 0 {
						px[y*Width+b*8+bit] |= 1 << p
					}
				}
			}
		}
	}
	return px
}

// Image 用給定的 16 色色盤畫成 RGBA。
func Image(buf []byte, pal []color.RGBA) (*image.RGBA, error) {
	if len(pal) < 16 {
		return nil, fmt.Errorf("cutscene: 色盤只有 %d 色，需要 16", len(pal))
	}
	px := Pixels(buf)
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	for i, v := range px {
		img.Set(i%Width, i/Width, pal[v&0x0F])
	}
	return img, nil
}

// 第一幕的三塊（`sub_1016D`／`sub_101B5`／`sub_10204`，docs/re/70 §5）。
// 位移是**解壓後**緩衝區裡的 byte 位置。
const (
	// FirstPlane01Off 是平面 0／1 交錯那一塊（每列 40 B ＋ 40 B）。
	FirstPlane01Off = 0x3E80 // 16,000
	// FirstPlane2Off／FirstPlane2Rows 是平面 2 那一塊，從第 120 列起。
	FirstPlane2Off  = 0xBB8 * 16 // 48,000
	FirstPlane2Row  = 120
	FirstPlane2Rows = 280
	// FirstPlane3Off 是平面 3 那一塊，整張 640 × 400。
	FirstPlane3Off = 0xE74 * 16 // 59,200
)

// FirstScenePixels 是第一幕的色號圖。
//
// ⭐ 它**不是** §2 那種整張 640×400：平面 0／1 是 320 px 寬貼在 x = 160，
// 平面 2 只蓋下面 280 列，平面 3 才是整張。拿整張的版面去畫，
// 畫面上的亭子會左右各出現一份。
//
// final 決定要不要疊平面 3。原版的平面 3 是**打完字之後**才換上去的那一頁
// （`sub_10204` 在淡出之後才跑），內容是一個大大的「終」——
// 打字的時候疊上去會把畫面蓋掉一半。
func FirstScenePixels(buf []byte, final bool) []byte {
	px := make([]byte, Width*Height)
	set := func(x, y, p int, v byte) {
		if v == 0 || x < 0 || x >= Width || y < 0 || y >= Height {
			return
		}
		for bit := 0; bit < 8; bit++ {
			if v&(0x80>>bit) != 0 {
				px[y*Width+x+bit] |= 1 << p
			}
		}
	}
	half := FirstSceneWidth / 8 // 40 B
	for y := 0; y < Height; y++ {
		off := FirstPlane01Off + y*half*2
		for b := 0; b < half; b++ {
			if off+half+b >= len(buf) {
				break
			}
			set(FirstSceneX+b*8, y, 0, buf[off+b])
			set(FirstSceneX+b*8, y, 1, buf[off+half+b])
		}
	}
	for r := 0; r < FirstPlane2Rows; r++ {
		off := FirstPlane2Off + r*half
		for b := 0; b < half; b++ {
			if off+b >= len(buf) {
				break
			}
			set(FirstSceneX+b*8, FirstPlane2Row+r, 2, buf[off+b])
		}
	}
	if final {
		for y := 0; y < Height; y++ {
			off := FirstPlane3Off + y*Stride
			for b := 0; b < Stride; b++ {
				if off+b >= len(buf) {
					break
				}
				set(b*8, y, 3, buf[off+b])
			}
		}
	}
	return px
}
