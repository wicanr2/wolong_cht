package gfx

import "testing"

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
