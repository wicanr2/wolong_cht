package main

import (
	"encoding/binary"
	"os"
	"testing"
)

// 原版那兩張表在 `KI.EXE` 的段內偏移 `0x9309`（docs/re/58 §2）。
// MZ 的檔頭長度寫在 `+0x08`（段數），所以檔案位移算得出來——
// **不要把位移寫死**，換一份執行檔就會錯。
const seasonTableSegOffset = 0x9309

func originalSeasonTables(t *testing.T) (music, palette []byte) {
	t.Helper()
	raw, err := os.ReadFile("../../workplace/orig/dosv/KI.EXE")
	if err != nil {
		t.Skip("找不到原版 KI.EXE，跳過")
	}
	hdr := int(binary.LittleEndian.Uint16(raw[8:])) * 16
	off := hdr + seasonTableSegOffset
	if off+24 > len(raw) {
		t.Fatalf("KI.EXE 只有 %d bytes，讀不到 0x%X 的表", len(raw), off)
	}
	return raw[off : off+12], raw[off+12 : off+24]
}

// 四季配樂表要逐 byte 等於原版。
func TestSeasonMusicMatchesOriginalTable(t *testing.T) {
	music, _ := originalSeasonTables(t)
	for m := 0; m < 12; m++ {
		if int(music[m]) != seasonMusic[m] {
			t.Errorf("%d 月：原版是曲 %d，remake 寫 %d", m+1, music[m], seasonMusic[m])
		}
	}
}

// ⭐ 交叉驗證：緊接著的第二張表是季節調色盤對（`sub_19336` 用同一個
// 月份索引讀），低 nibble ＝ 目標盤。**兩張表必須指向同一組季節切分**。
//
// 這一條才是「曲 2–5 是四季」的證據——單看第一張表只知道
// 「12 個月分成四組」，不知道那四組是不是季節。
func TestSeasonMusicAgreesWithPaletteTable(t *testing.T) {
	music, palette := originalSeasonTables(t)
	// 調色盤索引 → 季節（docs/re/06 §6）→ 該季的曲號。
	wantSong := map[byte]byte{0: 2, 1: 3, 2: 4, 3: 5} // 春 夏 秋 冬
	for m := 0; m < 12; m++ {
		target := palette[m] & 0x0F
		want, ok := wantSong[target]
		if !ok {
			t.Fatalf("%d 月的目標調色盤索引 %d 超出 0–3", m+1, target)
		}
		if music[m] != want {
			t.Errorf("%d 月：調色盤指向季節 %d（曲 %d），音樂表卻是曲 %d",
				m+1, target, want, music[m])
		}
	}
}

// 換季那三個月，來源盤要是前一季——順便確認 nibble 的高低沒讀反。
func TestSeasonTransitionMonthsCarryPreviousSeason(t *testing.T) {
	_, palette := originalSeasonTables(t)
	for _, m := range []int{3, 6, 9, 12} { // 轉入春／夏／秋／冬的月份
		v := palette[m-1]
		from, to := v>>4, v&0x0F
		if from == to {
			t.Errorf("%d 月應該是換季月，來源與目標卻都是 %d", m, to)
		}
		if want := byte((int(to) + 3) % 4); from != want {
			t.Errorf("%d 月的來源盤是 %d，前一季應該是 %d", m, from, want)
		}
	}
}
