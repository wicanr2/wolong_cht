package main

import (
	"crypto/sha256"
	"os"
	"testing"
)

// 原版由 dosgolem 匯出的80×50熱區表，每格涵蓋8×8像素。
// 預設不依賴私人原版產物；指定檔案時，缺檔或空表必須失敗。
func TestBattleHitTestsAgainstOriginalHotmap(t *testing.T) {
	path := os.Getenv("WOLONG_TACTICAL_HOTMAP")
	if path == "" {
		t.Skip("以 WOLONG_TACTICAL_HOTMAP 指定原版4000-byte熱區表")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 80*50 {
		t.Fatalf("原版熱區長度=%d，要4000", len(data))
	}
	l := dosvBattleLayoutFor(640, 400)
	sidePixels, bottomPixels := 0, 0
	for y := 0; y < 400; y++ {
		for x := 0; x < 640; x++ {
			id := int(data[(y/8)*80+x/8])
			wantSide := id >= 0x07 && id <= 0x0C
			row, sideOK := battleSideCommandIndexAt(l.SideCommands, x, y)
			if sideOK != wantSide || sideOK && battleSideCommandRowCode[row] != id-0x07 {
				t.Fatalf("側欄(%d,%d)：原版熱區=%02X，remake列=%d 命中=%v", x, y, id, row, sideOK)
			}
			wantBottom := id >= 0x15 && id <= 0x1A
			slot, bottomOK := splitBattleCommandIndexAt(l.BottomCommands, x, y)
			if bottomOK != wantBottom || bottomOK && battleBottomSlotSquad[slot] != id-0x15 {
				t.Fatalf("底列(%d,%d)：原版熱區=%02X，remake格=%d 命中=%v", x, y, id, slot, bottomOK)
			}
			if wantSide {
				sidePixels++
			}
			if wantBottom {
				bottomPixels++
			}
		}
	}
	// 正對照：不能讓空表或尚未進戰場的表得到全零通過。
	if sidePixels != 6*48*16 || bottomPixels != 6*80*32 {
		t.Fatalf("原版有效區域不足：側欄%d 底列%d", sidePixels, bottomPixels)
	}
	t.Logf("原版熱區SHA-256=%x；全畫布256000座標，側欄%d、底列%d像素", sha256.Sum256(data), sidePixels, bottomPixels)
}
