package gfx

import (
	"os"
	"path/filepath"
	"testing"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	p := filepath.Join("..", "..", "..", "workplace", "orig", "dosv", name)
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("找不到原版素材 %s，跳過", p)
	}
	return raw
}

// TestNoRemainder 是圖庫尺寸最便宜的檢核：餘數不是 0 就代表尺寸錯了。
// 四個圖庫的尺寸全部出自反組譯（docs/re/03），四個的餘數都該是 0。
func TestNoRemainder(t *testing.T) {
	for _, c := range []struct {
		spec  Spec
		file  string
		count int
	}{
		{Kao, "KAOGRF.DAT", 150},
		{Kyo, "KYOGRF.DAT", 15},
		{Ivent, "IVENTGRF.DAT", 3},
	} {
		got, rem := c.spec.Count(read(t, c.file))
		if rem != 0 {
			t.Errorf("%s 餘 %d byte —— 尺寸 %dx%d 是錯的",
				c.file, rem, c.spec.Width, c.spec.Height)
		}
		if got != c.count {
			t.Errorf("%s 解出 %d 張，預期 %d 張", c.file, got, c.count)
		}
	}
}

// TestIconRegionsCoverFile 釘住 ICONGRF 四段的長度加起來等於檔案大小。
func TestIconRegionsCoverFile(t *testing.T) {
	raw := read(t, "ICONGRF.DAT")
	total := 0
	for _, r := range IconRegions {
		if r.Offset != total {
			t.Fatalf("段 %s 的位移 0x%X 與前一段的結尾 0x%X 對不上",
				r.Name, r.Offset, total)
		}
		total += r.Length
	}
	if total != len(raw) {
		t.Errorf("四段合計 %d byte，檔案 %d byte", total, len(raw))
	}
}

func TestDOSVBattleCommandAssetsDecodeAtIDAOffsets(t *testing.T) {
	raw := read(t, "ICONGRF.DAT")
	r := IconRegions[1]
	seg := raw[r.Offset : r.Offset+r.Length]
	if _, err := DOSVBattleSideCommands.DecodeAt(seg, DOSVBattleSideCommandsOffset); err != nil {
		t.Fatal(err)
	}
	if _, err := DOSVBattleCommandBase.DecodeAt(seg, DOSVBattleCommandBaseOffset); err != nil {
		t.Fatal(err)
	}
	seen := make(map[string]bool)
	for i := 0; i < DOSVBattleCommandGlyphCount; i++ {
		off := DOSVBattleCommandGlyphOffset + i*DOSVBattleCommandGlyphStride
		pixels, err := DOSVBattleCommandGlyph.DecodeAt(seg, off)
		if err != nil {
			t.Fatalf("glyph %d: %v", i, err)
		}
		key := string(pixels)
		if seen[key] {
			t.Fatalf("glyph %d 與前一張重複；sub_1F888 的 SI 前進可能被解錯", i)
		}
		seen[key] = true
	}
	last := DOSVBattleCommandGlyphOffset +
		(DOSVBattleCommandGlyphCount-1)*DOSVBattleCommandGlyphStride
	if got, want := last+DOSVBattleCommandGlyph.FrameBytes(), 0x3D80; got != want {
		t.Fatalf("六張 glyph 末端 = 0x%X，預期 0x%X", got, want)
	}
}

// TestDecodeBounds 確認越界會回 error，不是 panic 或靜悄悄回錯的圖。
func TestDecodeBounds(t *testing.T) {
	raw := read(t, "KAOGRF.DAT")
	if _, err := Kao.Decode(raw, 150); err == nil {
		t.Error("第 150 張（超出範圍）應該回 error")
	}
	if _, err := Kao.Decode(raw, -1); err == nil {
		t.Error("負數索引應該回 error")
	}
}

// DOS/V sub_17D0D 以 DS:SI=word_10D50:0600h、AX=4006h 複製數值視窗
// 內框；由 sub_100DF 的段指標換算，它是 ICONGRF 第 3 段的 0x14A0h。
func TestDOSVAmountPanelResource(t *testing.T) {
	raw := read(t, "ICONGRF.DAT")
	segment := raw[IconRegions[3].Offset : IconRegions[3].Offset+IconRegions[3].Length]
	idx, err := DOSVAmountPanel.DecodeAt(segment, DOSVAmountPanelOffset)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(idx), DOSVAmountPanel.Width*DOSVAmountPanel.Height; got != want {
		t.Fatalf("DOS/V 數值內框像素數 = %d，want %d", got, want)
	}
	for i, v := range idx {
		if v >= 16 {
			t.Fatalf("DOS/V 數值內框第 %d 像素超出 4bpp：%d", i, v)
		}
	}
}

// sub_17D0D 的 96×64 blit 已包含固定的 3×6 按鈕圖像；sub_17D5F
// 另外只寫入同位置的 hit-test byte。每個 16×16 cell 都應有明顯
// 非背景像素，不能再用文字／向量框替代原版 glyph。
func TestDOSVAmountPanelContainsStaticButtonGlyphs(t *testing.T) {
	raw := read(t, "ICONGRF.DAT")
	segment := raw[IconRegions[3].Offset : IconRegions[3].Offset+IconRegions[3].Length]
	idx, err := DOSVAmountPanel.DecodeAt(segment, DOSVAmountPanelOffset)
	if err != nil {
		t.Fatal(err)
	}
	for row := 0; row < 3; row++ {
		for col := 0; col < 6; col++ {
			cell := make([]byte, 0, 16*16)
			for y := 0; y < 16; y++ {
				start := (16+row*16+y)*96 + col*16
				cell = append(cell, idx[start:start+16]...)
			}
			background := cell[0]
			different := 0
			for _, v := range cell {
				if v != background {
					different++
				}
			}
			if different < 128 {
				t.Fatalf("DOS/V 按鈕 glyph (%d,%d) 只有 %d 個非背景像素", row, col, different)
			}
		}
	}
}
