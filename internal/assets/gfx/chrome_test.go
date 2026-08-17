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

// TestViewBoxShape 釘住 docs/spec/55：縮小地圖的視野框是 24×11 的點陣，
// 只用兩個顏色，而且**只有左邊 20 px 有圖**。
func TestViewBoxShape(t *testing.T) {
	seg := iconSegment(t, "unknown3")
	pix, err := DecodeViewBox(seg, ViewBoxOffset)
	if err != nil {
		t.Fatal(err)
	}
	if len(pix) != ViewBoxWidth*ViewBoxRows {
		t.Fatalf("解出 %d 點，預期 %d", len(pix), ViewBoxWidth*ViewBoxRows)
	}
	for i, v := range pix {
		if v != ViewBoxTransparent && v != 0 && v != 15 {
			t.Fatalf("第 %d 點是色 %d，框只該有色 0、色 15 與透明", i, v)
		}
	}
	// 第 0 列：第 1–17 欄是白邊，這 17 px 與實機量到的白線等長。
	white := 0
	for x := 0; x < ViewBoxWidth; x++ {
		if pix[x] == 15 {
			white++
		}
	}
	if white != 17 {
		t.Errorf("第 0 列有 %d 個白點，預期 17", white)
	}
	// 右邊四欄整列透明——圖只有 20 px 寬。
	for y := 0; y < ViewBoxRows; y++ {
		for x := 20; x < ViewBoxWidth; x++ {
			if pix[y*ViewBoxWidth+x] != ViewBoxTransparent {
				t.Fatalf("(%d,%d) 不是透明，框應該只有 20 px 寬", x, y)
			}
		}
	}
}

// TestBattleFrameTilesFillSegment 釘住 docs/spec/31 §1.1：側欄外框的四塊
// （橫帶／角／左柱／右柱）在 `ICONGRF` 段 1 的最後，**四塊剛好用完整段**。
//
// 這一條擋的是「位移抄錯一格」：位移各差 0x40／0x40／0x80，
// 而最後一塊的結尾必須等於段長，錯一格就對不起來。
// 另外驗四塊都不是空的——全 0 會畫成一條黑柱，而黑柱在藍底上看起來
// 就像「沒畫」，與真的沒畫分不出來。
func TestBattleFrameTilesFillSegment(t *testing.T) {
	seg := iconSegment(t, "tiles")
	tiles := []struct {
		name string
		off  int
		spec Spec
	}{
		{"橫帶", DOSVBattleFrameBandOffset, DOSVBattleFrameHalf},
		{"角", DOSVBattleFrameCornerOffset, DOSVBattleFrameHalf},
		{"左柱", DOSVBattleFrameLeftOffset, DOSVBattleFrameCell},
		{"右柱", DOSVBattleFrameRightOffset, DOSVBattleFrameCell},
	}
	end := 0
	for _, tc := range tiles {
		idx, err := tc.spec.DecodeAt(seg, tc.off)
		if err != nil {
			t.Fatalf("%s 解不出來：%v", tc.name, err)
		}
		nonZero := 0
		for _, v := range idx {
			if v != 0 {
				nonZero++
			}
		}
		if nonZero == 0 {
			t.Errorf("%s（+%#x）整塊都是索引 0", tc.name, tc.off)
		}
		if e := tc.off + tc.spec.FrameBytes(); e > end {
			end = e
		}
	}
	if end != len(seg) {
		t.Errorf("四塊結束在 %#x，段長 %#x——位移或尺寸有一個抄錯了", end, len(seg))
	}
}
