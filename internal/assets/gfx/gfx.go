// Package gfx 解 `*GRF.DAT` 圖庫。
//
// 規格：docs/formats/03-grf-images.md（READY）
// 出處：docs/re/03-image-blitter.md
//
// 這一層是純解碼，回傳 *image.RGBA，不認識 Ebiten。
package gfx

import (
	"fmt"
	"image"

	"github.com/wicanr2/wolong_cht/internal/assets/palette"
)

// Planes 是原版的平面數。兩版都是 16 色 4 平面 planar
// （DOS/V 走 VGA Graphics Controller ＋ Sequencer，PC-98 走 GRCG）。
const Planes = 4

// Spec 描述一個圖庫的排列方式。數值全部出自載入器與繪製常式的參數，
// 不是拿檔案大小湊出來的 —— 見 docs/re/03。
type Spec struct {
	Name   string
	Width  int
	Height int
}

// 四個圖庫的規格。ICONGRF 是四段組合檔，不在這裡（見 IconRegions）。
var (
	Kao   = Spec{"KAOGRF", 64, 64}     // 武將頭像，150 張
	Kyo   = Spec{"KYOGRF", 96, 96}     // 據點景觀，15 張
	Ivent = Spec{"IVENTGRF", 288, 176} // 劇情過場，3 張
)

// IconRegion 是 ICONGRF.DAT 裡的一段。四段長度加起來剛好等於檔案大小。
type IconRegion struct {
	Name   string
	Offset int
	Length int
	Spec   Spec // 寬高為 0 表示該段尚未解出尺寸
}

// ICONGRF.DAT 的四段。段 1 與段 3 還沒解，寬高留 0。
var IconRegions = []IconRegion{
	{"banner", 0x0000, 0x2800, Spec{"ICONGRF/banner", 640, 32}},
	{"tiles", 0x2800, 0x3F00, Spec{Name: "ICONGRF/tiles"}},
	{"minimap", 0x6700, 0x3000, Spec{"ICONGRF/minimap", 192, 128}},
	{"unknown3", 0x9700, 0x23A0, Spec{Name: "ICONGRF/unknown3"}},
}

// 松崗 DOS/V 戰術指令介面位於 ICONGRF 第 1 段。下列 offset／尺寸都由
// sub_1C7F4／sub_1C863 的直接 blit 參數取得；offset 以第 1 段為基準。
const (
	DOSVBattleSideCommandsOffset = 0x1800
	DOSVBattleCommandBaseOffset  = 0x3000
	DOSVBattleCommandGlyphOffset = 0x3900
	DOSVBattleCommandGlyphStride = 0x00C0
	DOSVBattleCommandGlyphCount  = 6
)

var (
	DOSVBattleSideCommands = Spec{Name: "ICONGRF/DOSV battle side commands", Width: 128, Height: 96}
	DOSVBattleCommandBase  = Spec{Name: "ICONGRF/DOSV battle command base", Width: 80, Height: 32}
	DOSVBattleCommandGlyph = Spec{Name: "ICONGRF/DOSV battle command glyph", Width: 24, Height: 16}
)

// FrameBytes 是一張圖佔的位元組數。
func (s Spec) FrameBytes() int {
	return s.Width * s.Height / 2 // 4 bpp
}

// Count 回傳 data 裡有幾張圖，以及餘下幾個位元組。
//
// **餘數不是 0 就代表尺寸錯了**，這是最便宜的檢核 —— 四個圖庫的餘數都是 0。
func (s Spec) Count(data []byte) (count, remainder int) {
	size := s.FrameBytes()
	if size == 0 {
		return 0, len(data)
	}
	return len(data) / size, len(data) % size
}

// Decode 解出第 index 張圖的調色盤索引（值域 0–15）。
//
// 佈局是 plane-major：plane0 整張、plane1 整張、plane2、plane3。
// 這一點是從繪製常式推出來的 —— sub_1FA37 對四個平面各呼叫一次
// sub_1FAA2，而 sub_1FAA2 不重設來源指標，四次呼叫連續吃掉四段資料。
func (s Spec) Decode(data []byte, index int) ([]byte, error) {
	size := s.FrameBytes()
	if size == 0 {
		return nil, fmt.Errorf("gfx: %s 的尺寸還沒解出來", s.Name)
	}
	if index < 0 {
		return nil, fmt.Errorf("gfx: %s 第 %d 張超出範圍（共 %d 張）",
			s.Name, index, len(data)/size)
	}
	return s.DecodeAt(data, index*size)
}

// DecodeAt 解出 data 中指定 byte offset 的一張平面圖。
//
// 大型組合檔不一定把每張圖從檔案開頭連續編號；DOS/V 數值視窗就是
// ICONGRF 第 3 段內的獨立 96×64 資源。因此保留以 byte offset 定位的
// 入口，避免把段內資源硬切成錯誤的 frame index。
func (s Spec) DecodeAt(data []byte, offset int) ([]byte, error) {
	size := s.FrameBytes()
	if size == 0 {
		return nil, fmt.Errorf("gfx: %s 的尺寸還沒解出來", s.Name)
	}
	if offset < 0 || offset+size > len(data) {
		return nil, fmt.Errorf("gfx: %s 的 byte offset 0x%X 超出範圍（需要 0x%X byte）",
			s.Name, offset, size)
	}
	stride := s.Width / 8 // 每平面每列的位元組數
	plane := stride * s.Height
	out := make([]byte, s.Width*s.Height)
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			byteIdx := offset + y*stride + x/8
			bit := uint(7 - x%8) // 最高位在最左
			var v byte
			for p := 0; p < Planes; p++ {
				v |= ((data[byteIdx+p*plane] >> bit) & 1) << uint(p)
			}
			out[y*s.Width+x] = v
		}
	}
	return out, nil
}

// RenderRGBA 解出第 index 張圖並用給定的色組上色。
func (s Spec) RenderRGBA(data []byte, index int, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	idx, err := s.Decode(data, index)
	if err != nil {
		return nil, err
	}
	return s.renderRGBA(idx, p, bankIdx)
}

// RenderRGBAAt 解出 data 中指定 byte offset 的一張平面圖並上色。
func (s Spec) RenderRGBAAt(data []byte, offset int, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	idx, err := s.DecodeAt(data, offset)
	if err != nil {
		return nil, err
	}
	return s.renderRGBA(idx, p, bankIdx)
}

func (s Spec) renderRGBA(idx []byte, p *palette.Palette, bankIdx int) (*image.RGBA, error) {
	bank, err := p.Bank(bankIdx)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, s.Width, s.Height))
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			img.SetRGBA(x, y, bank[idx[y*s.Width+x]])
		}
	}
	return img, nil
}
