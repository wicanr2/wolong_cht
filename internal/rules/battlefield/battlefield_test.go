package battlefield

import (
	"os"
	"testing"
)

const exe = "../../../workplace/orig/dosv/KI.EXE"

// 兩張表都是從機器碼抄下來的常數——對著原版的 KI.EXE 驗一次。
// **抄表最容易錯的就是抄錯**，而且錯了不會有任何症狀。
func TestTablesMatchOriginal(t *testing.T) {
	b, err := os.ReadFile(exe)
	if err != nil {
		t.Skip("找不到原版 KI.EXE，跳過")
	}
	// seg000 偏移 → 檔案偏移是 +0x200。
	const terrainOff, pairsOff = 0x200 + 0x982F, 0x200 + 0x97F0

	for i, r := range terrainRanges {
		o := terrainOff + i*3
		if int(b[o]) != r.Kind || b[o+1] != r.Lo || b[o+2] != r.Hi {
			t.Errorf("地形表第 %d 筆抄錯：程式 {%d, %02X, %02X}，原版 {%d, %02X, %02X}",
				i, r.Kind, r.Lo, r.Hi, b[o], b[o+1], b[o+2])
		}
	}
	for i, p := range plainPairs {
		o := pairsOff + i*3
		if int(b[o]) != p.Offset || int(b[o+1]) != p.A || int(b[o+2]) != p.B {
			t.Errorf("配對表第 %d 筆抄錯：程式 {%d, %d, %d}，原版 {%d, %d, %d}",
				i, p.Offset, p.A, p.B, b[o], b[o+1], b[o+2])
		}
	}
}

// 地形分類：幾個已經在 docs/mechanics/30 §2 驗過的代表值。
func TestTerrainClassification(t *testing.T) {
	for _, tc := range []struct {
		tile byte
		want int
	}{
		{0x00, 0}, {0x05, 0}, // 平原
		{0xB8, 1}, {0xBF, 2}, // 跨水構造
		{0x70, 3}, {0xA7, 3}, {0xA9, 3}, // 山地
		{0xA8, 7},            // 山地區間裡挖掉的那一格是水
		{0x06, 4}, {0xB1, 5}, // 兩種林地
		{0x0E, 6}, {0x6F, 6}, {0xCA, 8}, // 水域與碼頭
	} {
		if got := Terrain(tc.tile); got != tc.want {
			t.Errorf("Terrain(0x%02X) ＝ %d，應為 %d", tc.tile, got, tc.want)
		}
	}
}

// 平原的配對：山＋山 → 192、平原＋平原 → 198、配不到 → 198。
// 這幾張的高處格數在 docs/re/11 §4.4 用另一份資料獨立驗過。
func TestPlainPairing(t *testing.T) {
	mount := Neighbours{DownLeft: 3, DownRight: 3}
	if f, rot := plainField(2, mount); f != 192 || rot {
		t.Errorf("山＋山 → %d（轉 %v），應為 192（不轉）", f, rot)
	}
	plain := Neighbours{DownLeft: 0, DownRight: 0}
	if f, _ := plainField(2, plain); f != 198 {
		t.Errorf("平原＋平原 → %d，應為 198", f)
	}
	// ⭐ 換過順序才中 → 轉 180 度。「山 ＋ 平原」有登記、「平原 ＋ 山」沒有。
	fwd, rotF := plainField(2, Neighbours{DownLeft: 3, DownRight: 0})
	rev, rotR := plainField(2, Neighbours{DownLeft: 0, DownRight: 3})
	if fwd != rev {
		t.Errorf("同一組地形換順序選到不同戰場：%d 對 %d", fwd, rev)
	}
	if rotF == rotR {
		t.Error("換順序之後旋轉旗標應該相反")
	}
	if rotF {
		t.Error("表上登記的順序不該轉")
	}
	// 配不到的組合走 fallback。
	if f, _ := plainField(2, Neighbours{DownLeft: 1, DownRight: 2}); f != 198 {
		t.Errorf("配不到的組合 → %d，應為 198（偏移 6）", f)
	}
}

// 行進方向決定取樣哪兩格。
func TestDirectionPicksNeighbours(t *testing.T) {
	n := Neighbours{Centre: 3, TwoDown: 0, DownLeft: 6, DownRight: 6}
	// 方向 0／1 用「中心 ＋ 兩格下方」，順序相反。
	f0, rot0 := plainField(0, n)
	f1, rot1 := plainField(1, n)
	if f0 != f1 {
		t.Errorf("方向 0 與 1 選到不同戰場：%d 對 %d", f0, f1)
	}
	if rot0 == rot1 {
		t.Error("方向 0 與 1 的旋轉旗標應該相反")
	}
	// 方向 2 用「左下 ＋ 右下」＝ 水域 ＋ 水域 → 205。
	if f, _ := plainField(2, n); f != 205 {
		t.Errorf("水域＋水域 → %d，應為 205", f)
	}
}

// 地形類型 1–7 直接對到 0xCF–0xD5。
func TestTerrainFieldsAreContiguous(t *testing.T) {
	for k := 1; k <= 7; k++ {
		f, _ := Select(0, Neighbours{Down: k})
		if want := 0xCE + k; f != want {
			t.Errorf("地形類型 %d → 戰場 %d，應為 %d", k, f, want)
		}
		if f < NumCityFields || f >= NumFields {
			t.Errorf("戰場 %d 超出野戰的區間", f)
		}
	}
	// 類型 8 隨機挑 209–212、類型 9 固定 213。
	for r := 0; r < 4; r++ {
		if f := SelectWater(8, r); f != 209+r {
			t.Errorf("類型 8 亂數 %d → %d，應為 %d", r, f, 209+r)
		}
	}
	if f := SelectWater(9, 0); f != 213 {
		t.Errorf("類型 9 → %d，應為 213", f)
	}
}
