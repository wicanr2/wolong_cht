// Package cutscene 解 `OPEN_S*.DAT`／`END_S*.DAT` 這一批過場畫面。
//
// 格式見 `docs/formats/09-cutscene-images.md`。摘要：
//
//	外層   與 `MMAP.MAP` 同一種 RLE（`internal/assets/rle`）
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
// ⚠ **長度可能略短於 `Size`**：原版的緩衝區是先配好的，檔案只編到最後一個
// 非零的 byte（實測 `END_S3` 解出 127,749 B，差 251 B）。缺的那一段當成 0，
// 因為原版那一塊在整張貼圖前已經被前一幕蓋成黑色。
func Decode(src []byte) []byte {
	out := rle.Decode(src)
	if len(out) >= Size {
		return out[:Size]
	}
	buf := make([]byte, Size)
	copy(buf, out)
	return buf
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
