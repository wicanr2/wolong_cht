package gfx

import (
	"encoding/binary"
	"fmt"
	"image"

	"github.com/wicanr2/wolong_cht/internal/assets/palette"
)

// DOS/V KI.EXE 的滑鼠游標是 seg002 內建的兩層 16×16 mask，不在
// MOUSE.MCH／MOUSE.SCH。sub_201E4 先以 set/reset=0x0F 畫白色外框，
// 再以 set/reset=0x0A 畫紅色填色；兩次 sub_2020C 各消費 32 byte。
//
// 位址基準：MZ header 後的 image 以 0x200 開始，IDA 的 seg002 對應
// image+0x10000，原始 KI.EXE 的 `mov si,31Bh` 因此落在 file 0x1051B。
// 這些常數保留原始 IDA 位址與檔案偏移的可回查關係，不把游標資料
// 改名成沒有證據的通用滑鼠格式。
const (
	DOSVKIImageSegmentOffset = 0x10000
	DOSVCursorSegmentOffset  = 0x031B  // IDA: seg002:031B
	DOSVCursorFileOffset     = 0x1051B // DOS/V KI.EXE，header=0x200
	DOSVCursorWidth          = 16
	DOSVCursorHeight         = 16
	DOSVCursorMaskBytes      = DOSVCursorWidth * DOSVCursorHeight / 8
	DOSVCursorSourceBytes    = DOSVCursorMaskBytes * 2
	DOSVCursorTransparent    = 0xFF
	DOSVCursorOutline        = 0x0F // EGA set/reset=0x0F
	DOSVCursorFill           = 0x0A // EGA set/reset=0x0A
)

// dosvCursorSourceOffset 依 MZ header 算出 KI.EXE 內 seg002:031B 的檔案位址。
// DOS/V 本版驗證值為 0x1051B；若 header 或 segment 佈局改變，拒絕靜默
// 讀錯位址。
func dosvCursorSourceOffset(exe []byte) (int, error) {
	if len(exe) < 0x0A || exe[0] != 'M' || exe[1] != 'Z' {
		return 0, fmt.Errorf("DOS/V KI.EXE 不是可辨識的 MZ 檔")
	}
	headerBytes := int(binary.LittleEndian.Uint16(exe[0x08:0x0A])) * 16
	off := headerBytes + DOSVKIImageSegmentOffset + DOSVCursorSegmentOffset
	if off != DOSVCursorFileOffset {
		return 0, fmt.Errorf("DOS/V KI.EXE image layout = 0x%X，預期 cursor file offset 0x%X",
			off, DOSVCursorFileOffset)
	}
	if off+DOSVCursorSourceBytes > len(exe) {
		return 0, fmt.Errorf("DOS/V KI.EXE 不足 cursor bytes：0x%X..0x%X／檔案 0x%X",
			off, off+DOSVCursorSourceBytes, len(exe))
	}
	return off, nil
}

// DecodeDOSVCursor 解出 DOS/V KI.EXE 內建游標的 palette index。
//
// 回傳長度為 16×16；0xFF 透明、0x0F 白色外框、0x0A 紅色填色。
// 原始每列是 little-endian word，sub_2020C 寫回 EGA 時先寫 high byte，
// 所以 decoder 必須反轉每列兩個 source byte，而不是把它當一般 MSB-first
// 連續 16-bit bitmap 直接讀取。
func DecodeDOSVCursor(exe []byte) ([]byte, error) {
	off, err := dosvCursorSourceOffset(exe)
	if err != nil {
		return nil, err
	}
	source := exe[off : off+DOSVCursorSourceBytes]
	out := make([]byte, DOSVCursorWidth*DOSVCursorHeight)
	for i := range out {
		out[i] = DOSVCursorTransparent
	}
	for layer, colour := range [...]byte{DOSVCursorOutline, DOSVCursorFill} {
		base := layer * DOSVCursorMaskBytes
		for y := 0; y < DOSVCursorHeight; y++ {
			// `mov ax,[si]` followed by `[bx]=ah, [bx+1]=al`。
			row := source[base+y*2 : base+y*2+2]
			for x := 0; x < DOSVCursorWidth; x++ {
				bit := (row[1-x/8] >> uint(7-x%8)) & 1
				if bit != 0 {
					out[y*DOSVCursorWidth+x] = colour
				}
			}
		}
	}
	return out, nil
}

// RenderDOSVCursorPixels 將 DecodeDOSVCursor 的 palette index 轉成透明 RGBA。
func RenderDOSVCursorPixels(idx []byte, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	if len(idx) != DOSVCursorWidth*DOSVCursorHeight {
		return nil, fmt.Errorf("DOS/V cursor 像素數 %d，預期 %d",
			len(idx), DOSVCursorWidth*DOSVCursorHeight)
	}
	bank, err := p.Bank(bankIdx)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, DOSVCursorWidth, DOSVCursorHeight))
	for y := 0; y < DOSVCursorHeight; y++ {
		for x := 0; x < DOSVCursorWidth; x++ {
			v := idx[y*DOSVCursorWidth+x]
			if v == DOSVCursorTransparent {
				continue
			}
			img.SetRGBA(x, y, bank[v])
		}
	}
	return img, nil
}

// RenderDOSVCursor 從使用者提供的原始 KI.EXE 直接解碼游標並上色。
func RenderDOSVCursor(exe []byte, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	idx, err := DecodeDOSVCursor(exe)
	if err != nil {
		return nil, err
	}
	return RenderDOSVCursorPixels(idx, p, bankIdx)
}
