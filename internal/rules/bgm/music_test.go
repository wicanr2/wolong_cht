package bgm

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/battlefield"
)

// 原版那兩張表在 `KI.EXE` 的段內偏移 `0x9309`（docs/re/58 §2）。
// MZ 的檔頭長度寫在 `+0x08`（段數），所以檔案位移算得出來——
// **不要把位移寫死**，換一份執行檔就會錯。
const seasonTableSegOffset = 0x9309

func originalSeasonTables(t *testing.T) (music, palette []byte) {
	t.Helper()
	raw, err := os.ReadFile("../../../workplace/orig/dosv/KI.EXE")
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

// 戰場音樂的三個門檻要跟 `internal/rules/battlefield` 的常數對齊。
//
// 原版是 `cmp al, 0C0h` 與 `cmp al, 0D1h`（docs/re/58 §4）——
// **門檻吃的是戰場編號，不是「攻城／野戰」這個布林值**，
// 所以山地／林地／水域那一組（0xD1 起）才分得出來。
func TestBattleMusicThresholdsMatchBattlefieldConstants(t *testing.T) {
	if battlefield.FieldBase != 0xC0 {
		t.Errorf("平原野戰的基底是 0x%X，原版是 0xC0", battlefield.FieldBase)
	}
	// 地形類型 1–7 → 0xCF–0xD5；門檻 0xD1 落在類型 2 與 3 之間。
	if battlefield.TerrainBase+2 != 0xD0 || battlefield.TerrainBase+3 != 0xD1 {
		t.Errorf("地形戰場基底 0x%X 讓 0xD1 的門檻切錯位置", battlefield.TerrainBase)
	}
	// 攻城戰的編號是據點編號，一定落在門檻以下。
	if battlefield.NumCityFields > battlefield.FieldBase {
		t.Errorf("據點戰場有 %d 張，會撞進野戰的編號區間",
			battlefield.NumCityFields)
	}
}

// 逐條對 docs/re/58 的表。**規則搬家之後這一支是它的護欄**——
// 桌面與手機共用同一份，改壞了兩邊一起壞。
func TestTrackCoversEveryScene(t *testing.T) {
	for _, tc := range []struct {
		name  string
		scene Scene
		want  string
	}{
		{"啟動殼層", Scene{Launcher: true}, "bgm-0"},
		{"結局蓋過勝負", Scene{Ending: true, GameOver: true, Month: 4}, "endbgm-0"},
		{"還沒開局", Scene{}, ""},
		{"勝負已定", Scene{GameOver: true, Month: 4}, "overbgm-0"},
		{"事件與對話", Scene{Message: true, Month: 4}, "bgm-6"},
		{"春", Scene{Month: 3}, "bgm-2"},
		{"夏", Scene{Month: 6}, "bgm-3"},
		{"秋", Scene{Month: 9}, "bgm-4"},
		{"冬（跨年）", Scene{Month: 12}, "bgm-5"},
		{"冬（一月）", Scene{Month: 1}, "bgm-5"},
		{"月份越界", Scene{Month: 13}, ""},
		{"平原野戰", Scene{Month: 4, Battle: &Battle{Field: 0xC5}}, "bgm-9"},
		{"山林水戰", Scene{Month: 4, Battle: &Battle{Field: 0xD2}}, "bgm-10"},
		{"攻城・玩家攻", Scene{Month: 4, Battle: &Battle{Field: 0x52, PlayerAttacks: true}}, "bgm-7"},
		{"攻城・玩家守", Scene{Month: 4, Battle: &Battle{Field: 0x52}}, "bgm-8"},
		// ⭐ 戰術畫面比事件訊息優先：原版進戰術先停曲再依戰場挑。
		{"戰術蓋過訊息", Scene{Month: 4, Message: true, Battle: &Battle{Field: 0xC5}}, "bgm-9"},
	} {
		if got := Track(tc.scene); got != tc.want {
			t.Errorf("%s：Track ＝ %q，預期 %q", tc.name, got, tc.want)
		}
	}
}
