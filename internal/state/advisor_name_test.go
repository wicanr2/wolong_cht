package state

import (
	"os"
	"path/filepath"
	"testing"
)

// 自定軍師的肖像與六個字在區塊 +0x52A1／+0x52A2（docs/formats/10 §4）：
// 原版四個劇本都是 FF ＋ 六個 D0A1，改過之後要能 byte-for-byte 寫回。
func TestAdvisorNameRoundTripsThroughBlock(t *testing.T) {
	path := filepath.Join("..", "..", "workplace", "orig", "dosv", "SINARIO.DAT")
	if _, err := os.Stat(path); err != nil {
		t.Skip("沒有原版素材")
	}
	w, err := LoadScenario(path, 0)
	if err != nil {
		t.Fatal(err)
	}
	if w.AdvisorPortrait != 0xFF {
		t.Fatalf("原版肖像欄是 %#x，預期 FF", w.AdvisorPortrait)
	}
	for i := 0; i < advisorNameLen; i += 2 {
		if w.AdvisorName[i] != 0xA1 || w.AdvisorName[i+1] != 0xD0 {
			t.Fatalf("原版名字欄第 %d 字是 %02X%02X，預期空標記 A1 D0", i/2, w.AdvisorName[i], w.AdvisorName[i+1])
		}
	}
	w.AdvisorPortrait = 0x91
	name := []byte{0xA4, 0xD5, 0xAE, 0xD4, 0xA1, 0x40, 0xA1, 0x40, 0xA1, 0x40, 0xA1, 0x40} // 孔明＋四個全形空白
	copy(w.AdvisorName[:], name)
	b := w.Bytes()
	if b[advisorPortraitOffset] != 0x91 || string(b[advisorNameOffset:advisorNameOffset+advisorNameLen]) != string(name) {
		t.Fatal("Bytes() 沒把肖像與名字寫回區塊")
	}
	again := loadBlock(b)
	if again.AdvisorPortrait != 0x91 || again.AdvisorName != w.AdvisorName {
		t.Fatal("重新載入後肖像或名字變了")
	}
}

// 自訂軍師：+0x02 變 0x7F、名字與別號各三字、{4} 取得到名字。
func TestCustomAdvisorNameFeedsTalkVars(t *testing.T) {
	path := filepath.Join("..", "..", "workplace", "orig", "dosv", "SINARIO.DAT")
	if _, err := os.Stat(path); err != nil {
		t.Skip("沒有原版素材")
	}
	w, err := LoadScenario(path, 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	if w.HasCustomAdvisor() {
		t.Fatal("原版劇本不該有自訂軍師")
	}
	name := []byte{0xA4, 0xD5, 0xAE, 0xD4, 0xA1, 0x40, 0xA7, 0xB5, 0xAB, 0xC2} // 孔明　／ 臥龍
	w.SetCustomAdvisor(0x91, name)
	if w.Factions[w.Player].Advisor != NoAdvisor || !w.HasCustomAdvisor() {
		t.Fatal("SetCustomAdvisor 沒把 +0x02 寫成 0x7F")
	}
	n, a := w.AdvisorNameRaw()
	if n != "\xA4\xD5\xAE\xD4" || a != "\xA7\xB5\xAB\xC2" {
		t.Fatalf("名字 %x 別號 %x", n, a)
	}
	vars, ok := w.TalkNoticeVars(TalkNotice{}, nil)
	if !ok || vars['4'] != n {
		t.Fatalf("{4} = %x，預期自訂軍師名", vars['4'])
	}
}
