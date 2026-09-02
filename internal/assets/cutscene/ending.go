package cutscene

import (
	"encoding/binary"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/assets/palette"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
)

// 結局的常數。出處 docs/re/70、規格 docs/spec/67。
const (
	// Scenes 是結局的幕數，也是 `ENDPAL.BRG` 的色盤組數。
	Scenes = 12
	// TextOffset 是結尾文字在 `D7END.EXE` **載入映像**裡的段內位移。
	TextOffset = 0x5F0
	// TextChars 是字數（`sub_10094` 的 `mov cx, 0C8h`）。
	TextChars = 200
	// TextCols 是一行幾個字：x 從 0x40 每字 +0x10、到 0x180 換行。
	TextCols = (0x180 - 0x40) / 0x10
	// TextX／TextY／TextAdvance／TextLeading 是版面（`sub_10238`）。
	TextX       = 0x40
	TextY       = 0x30
	TextAdvance = 0x10
	TextLeading = 0x20
	// FadeSteps 是淡入淡出的階數（`cx` 0–0x10）。
	FadeSteps = 17
	// FirstSceneWidth／FirstSceneX：第一幕是 320 px 寬，貼在 x = 160
	// （`sub_1016D` 每列每平面 40 B、`di` 從 0x14 起）。
	FirstSceneWidth = 320
	FirstSceneX     = 160
)

// Ending 是播結局要用的全部素材。
type Ending struct {
	// Frames[n] 是第 n 幕，已經上好色。
	Frames [Scenes]*image.RGBA
	// Final 是第一幕**打完字之後**才換上去的那一頁（大大的「終」）。
	Final *image.RGBA
	// Lines 是結尾文字，已照原版版面切成每行 TextCols 個字。
	Lines []string
}

// LoadEnding 從原版素材目錄讀結局要的東西。
//
// 文字**燒在 `D7END.EXE` 裡**，不在 `TALK.DAT`——所以這一支要讀執行檔。
// 讀的是玩家自備的原版目錄，產出不進版控（`CLAUDE.md` §9）。
func LoadEnding(dir string) (*Ending, error) {
	palRaw, err := os.ReadFile(filepath.Join(dir, "ENDPAL.BRG"))
	if err != nil {
		return nil, fmt.Errorf("讀不到 ENDPAL.BRG：%w", err)
	}
	pal, err := palette.Parse(palRaw)
	if err != nil {
		return nil, err
	}
	out := &Ending{}
	for n := 0; n < Scenes; n++ {
		name := fmt.Sprintf("END_S%d.DAT", n+1)
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("讀不到 %s：%w", name, err)
		}
		bank, err := pal.Bank(n) // 第 n 幕用第 n 組（sub_1035F 的 al ＝ 幕號）
		if err != nil {
			return nil, fmt.Errorf("%s 的色盤第 %d 組：%w", name, n, err)
		}
		buf, err := Decode(raw)
		if err != nil {
			return nil, fmt.Errorf("%s：%w", name, err)
		}
		paint := func(px []byte) *image.RGBA {
			img := image.NewRGBA(image.Rect(0, 0, Width, Height))
			for i, v := range px {
				img.Set(i%Width, i/Width, bank[v&0x0F])
			}
			return img
		}
		if n == 0 {
			// 第一幕是三塊拼的，不是整張（docs/formats/09 §6）。
			out.Frames[n] = paint(FirstScenePixels(buf, false))
			out.Final = paint(FirstScenePixels(buf, true))
			continue
		}
		out.Frames[n] = paint(Pixels(buf))
	}
	lines, err := EndingText(dir)
	if err != nil {
		return nil, err
	}
	out.Lines = lines
	return out, nil
}

// EndingText 從 `D7END.EXE` 取出結尾文字，切成每行 TextCols 個全形字。
func EndingText(dir string) ([]string, error) {
	exe, err := os.ReadFile(filepath.Join(dir, "D7END.EXE"))
	if err != nil {
		return nil, fmt.Errorf("讀不到 D7END.EXE：%w", err)
	}
	raw, err := textBytes(exe)
	if err != nil {
		return nil, err
	}
	runes := []rune(text.Decode(raw, text.Big5))
	lines := make([]string, 0, (len(runes)+TextCols-1)/TextCols)
	for i := 0; i < len(runes); i += TextCols {
		j := i + TextCols
		if j > len(runes) {
			j = len(runes)
		}
		lines = append(lines, string(runes[i:j]))
	}
	return lines, nil
}

// textBytes 從 MZ 的載入映像裡切出結尾文字那 400 個 byte。
func textBytes(exe []byte) ([]byte, error) {
	if len(exe) < 0x20 || exe[0] != 'M' || exe[1] != 'Z' {
		return nil, fmt.Errorf("D7END.EXE 不是 MZ 執行檔")
	}
	base := int(binary.LittleEndian.Uint16(exe[8:])) * 16
	lo, hi := base+TextOffset, base+TextOffset+TextChars*2
	if hi > len(exe) {
		return nil, fmt.Errorf("D7END.EXE 只有 %d B，取不到 %d..%d", len(exe), lo, hi)
	}
	return exe[lo:hi], nil
}
