package gfx

import (
	"os"
	"path/filepath"
	"testing"
)

// iconSegment 讀原版 `ICONGRF.DAT` 的某一段。沒有素材就跳過——
// 這個 repo 不收原版資產（CLAUDE.md §9）。
func iconSegment(t *testing.T, name string) []byte {
	t.Helper()
	p := filepath.Join("..", "..", "..", "workplace", "orig", "dosv", "ICONGRF.DAT")
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("找不到原版素材 %s，跳過", p)
	}
	for _, r := range IconRegions {
		if r.Name != name {
			continue
		}
		if r.Offset+r.Length > len(raw) {
			t.Fatalf("%s 段超出檔案（%d B）", name, len(raw))
		}
		return raw[r.Offset : r.Offset+r.Length]
	}
	t.Fatalf("沒有 %s 這一段", name)
	return nil
}

// 視窗底紋：32×32、只有兩個索引（0 黑、8 深藍），檔尾那 128 byte。
func TestDecodeWindowTextureIsTwoColour(t *testing.T) {
	seg := make([]byte, WindowTextureOffset+windowTextureBytes)
	for i := 0; i < windowTextureBytes; i++ {
		seg[WindowTextureOffset+i] = byte(i * 37) // 隨便一組位元
	}
	idx, err := DecodeWindowTexture(seg, WindowTextureOffset)
	if err != nil {
		t.Fatal(err)
	}
	if len(idx) != WindowTextureSize*WindowTextureSize {
		t.Fatalf("解出 %d 點，要 %d", len(idx), WindowTextureSize*WindowTextureSize)
	}
	for i, v := range idx {
		if v != windowTexturePaper && v != windowTextureInk {
			t.Fatalf("第 %d 點是索引 %d，只能是 %d 或 %d",
				i, v, windowTexturePaper, windowTextureInk)
		}
	}
	// 位元序是 MSB 先：第 0 個 byte 的最高位是 (0, 0)。
	seg[WindowTextureOffset] = 0x80
	idx, _ = DecodeWindowTexture(seg, WindowTextureOffset)
	if idx[0] != windowTextureInk || idx[1] != windowTexturePaper {
		t.Errorf("位元序反了：前兩點 %d %d", idx[0], idx[1])
	}
	// 段不夠長要**明確報錯**，不要解出半塊。
	if _, err := DecodeWindowTexture(seg[:WindowTextureOffset+4], WindowTextureOffset); err == nil {
		t.Error("段太短卻沒有報錯")
	}
}

// TestDigitFontShape 釘住 docs/spec/52 §5：數字字模在 `ICONGRF` 段 3 的
// `+0x840`，8×16、11 格，而且**墨水只有 14 列**（上下各留一列空白）。
//
// 那 14 列是這一條的重點：倚天的 ASCII 數字只有 9 列，
// 同樣的位置同樣的顏色，字形不同還是逐像素差得出來。
func TestDigitFontShape(t *testing.T) {
	seg := iconSegment(t, "unknown3")
	for i := 0; i < DigitMinus; i++ { // 0–9；負號另外驗
		mask, err := DecodeDigit(seg, i)
		if err != nil {
			t.Fatalf("第 %d 格：%v", i, err)
		}
		if len(mask) != DigitWidth*DigitHeight {
			t.Fatalf("第 %d 格解出 %d 點，預期 %d", i, len(mask), DigitWidth*DigitHeight)
		}
		top, bottom := -1, -1
		for y := 0; y < DigitHeight; y++ {
			ink := false
			for x := 0; x < DigitWidth; x++ {
				if mask[y*DigitWidth+x] != 0 {
					ink = true
				}
			}
			if ink {
				if top < 0 {
					top = y
				}
				bottom = y
			}
		}
		if top != 1 || bottom != 14 {
			t.Errorf("第 %d 格的墨水在第 %d–%d 列，預期 1–14", i, top, bottom)
		}
	}
	// 第 10 格是負號：只有正中間**一列**有墨水。
	minus, err := DecodeDigit(seg, DigitMinus)
	if err != nil {
		t.Fatal(err)
	}
	rows := 0
	for y := 0; y < DigitHeight; y++ {
		for x := 0; x < DigitWidth; x++ {
			if minus[y*DigitWidth+x] != 0 {
				rows++
				break
			}
		}
	}
	if rows != 1 {
		t.Errorf("負號有 %d 列墨水，預期 1", rows)
	}
}
