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

// DOS/V sub_17D0D 以 DS:SI=word_10D50:0600h、AX=4006h 複製數值視窗內框。
// word_10D50 ＝ 段 3 的載入位址 ＋ 0x9A0（docs/re/48 §6），所以段內位移
// 是 0x0FA0。
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

// 顯示清單 op 09 取的四張 24×16 圖示（docs/re/48 §4、§6）。
// 這裡驗的是**位址換算**：`word_10D50` 的 0x1200 ＝ 段 3 的 0x1BA0。
// 內容檢查只驗「不是一整片單色」：紅色那組是黑剪影配紅底（兩色），
// 綠色那組多一圈邊框（更多色）。**兩組不是同一張圖的換色版**——
// 剪影逐像素比不相等，所以這裡不驗那個。
func TestDOSVResourceIcons(t *testing.T) {
	raw := read(t, "ICONGRF.DAT")
	segment := raw[IconRegions[3].Offset : IconRegions[3].Offset+IconRegions[3].Length]
	for _, base := range []int{DOSVResourceIconOffset, DOSVResourceIconGreenOffset} {
		for i := 0; i < DOSVResourceIconCount; i++ {
			off := base + i*DOSVResourceIconStride
			idx, err := DOSVResourceIcon.DecodeAt(segment, off)
			if err != nil {
				t.Fatalf("0x%X 第 %d 張：%v", base, i, err)
			}
			seen := map[byte]int{}
			for _, v := range idx {
				if v >= 16 {
					t.Fatalf("0x%X 第 %d 張超出 4bpp：%d", base, i, v)
				}
				seen[v]++
			}
			if len(seen) < 2 {
				t.Fatalf("0x%X 第 %d 張只有 %d 種顏色——偏移八成落在空白區",
					base, i, len(seen))
			}
		}
	}
}

// 編成畫面的「兵種 4 ＝ 空槽」圖示（docs/re/49 §3）。
//
// 位移是從 `sub_16DA8` 的 `0x15C0 + (兵種−1)×0xC0` 推出來的，兵種 1–3 落在
// 綠色那一組裡，兵種 4 落在**綠組後面一張**。所以這一張要單獨驗：
// 它必須解得開、而且不是綠組最後一張的重複。
func TestDOSVEmptySlotIcon(t *testing.T) {
	raw := read(t, "ICONGRF.DAT")
	segment := raw[IconRegions[3].Offset : IconRegions[3].Offset+IconRegions[3].Length]
	empty, err := DOSVResourceIcon.DecodeAt(segment, DOSVEmptySlotIconOffset)
	if err != nil {
		t.Fatalf("空槽圖示：%v", err)
	}
	last, err := DOSVResourceIcon.DecodeAt(segment,
		DOSVResourceIconGreenOffset+(DOSVResourceIconCount-1)*DOSVResourceIconStride)
	if err != nil {
		t.Fatalf("綠組最後一張：%v", err)
	}
	same := true
	seen := map[byte]int{}
	for i := range empty {
		if empty[i] >= 16 {
			t.Fatalf("空槽圖示超出 4bpp：%d", empty[i])
		}
		if empty[i] != last[i] {
			same = false
		}
		seen[empty[i]]++
	}
	if same {
		t.Error("空槽圖示與綠組最後一張逐像素相同——位移八成算錯了")
	}
	// 空槽本來就該是一格空白，所以**不驗顏色數 ≥ 2**——那條檢查對其他
	// 圖示是「偏移有沒有落在空白區」的護欄，用在這一張會反過來擋住正解。
	// 實測是 384 個像素全 0，也就是原版在空槽那一格貼一張全黑圖把前一張擦掉。
	//
	// ⚠ 但「全黑」與「位移落在段尾之外」長得一模一樣。分辨的方法是看
	// **這一張後面還有沒有內容**：有，就表示這片空白是刻意留的。
	if len(seen) == 1 {
		tail := segment[DOSVEmptySlotIconOffset+DOSVResourceIconStride:]
		nonzero := false
		for _, v := range tail {
			if v != 0 {
				nonzero = true
				break
			}
		}
		if !nonzero {
			t.Errorf("空槽圖示全 0 而它之後也全 0——0x%X 可能只是段尾，不是一張圖",
				DOSVEmptySlotIconOffset)
		}
	}
}

// TestDOSVOrderIcons 驗戰術底列那六張「目前命令」的圖示
// （`sub_1C673` 取 `碼 × 0xC0`，段內位移從 0 起算）。
//
// 六張要**互不相同**——它們是六個命令的圖示，任兩張一樣就表示
// stride 或起點算錯了。
func TestDOSVOrderIcons(t *testing.T) {
	raw := read(t, "ICONGRF.DAT")
	segment := raw[IconRegions[3].Offset : IconRegions[3].Offset+IconRegions[3].Length]
	icons := make([][]byte, DOSVOrderIconCount)
	for i := range icons {
		px, err := DOSVResourceIcon.DecodeAt(segment,
			DOSVOrderIconOffset+i*DOSVOrderIconStride)
		if err != nil {
			t.Fatalf("命令圖示 %d：%v", i, err)
		}
		seen := map[byte]bool{}
		for _, v := range px {
			if v >= 16 {
				t.Fatalf("命令圖示 %d 超出 4bpp：%d", i, v)
			}
			seen[v] = true
		}
		// 六張都是有內容的圖示（不像空槽那張是刻意的全黑）。
		if len(seen) < 2 {
			t.Errorf("命令圖示 %d 只有 %d 種顏色——位移可能落在空白區", i, len(seen))
		}
		icons[i] = px
	}
	for i := range icons {
		for j := i + 1; j < len(icons); j++ {
			same := true
			for k := range icons[i] {
				if icons[i][k] != icons[j][k] {
					same = false
					break
				}
			}
			if same {
				t.Errorf("命令圖示 %d 與 %d 逐像素相同", i, j)
			}
		}
	}
	// 第六張（退卻）之後接的是外框 motif（ChromeEdge），
	// 六張 × 0xC0 ＝ 0x480 還在 0x6C0 之前，不會撞到它。
	if DOSVOrderIconOffset+DOSVOrderIconCount*DOSVOrderIconStride > ChromeEdge {
		t.Errorf("六張命令圖示的尾端 0x%X 撞進外框圖塊 0x%X",
			DOSVOrderIconOffset+DOSVOrderIconCount*DOSVOrderIconStride, ChromeEdge)
	}
}
