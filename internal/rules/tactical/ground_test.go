package tactical

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
)

// realLib 載入原版戰場資料；沒有素材就跳過。
func realLib(t *testing.T) *battle.Library {
	t.Helper()
	const dir = "../../../workplace/orig/dosv"
	if _, err := os.Stat(dir + "/BATTLE.MAP"); err != nil {
		t.Skip("找不到原版 BATTLE.MAP，跳過")
	}
	read := func(name string) []byte {
		b, err := os.ReadFile(dir + "/" + name)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	lib, err := battle.Parse(read("BATTLE.MAP"), read("BATTLE.MDL"), nil)
	if err != nil {
		t.Fatal(err)
	}
	return lib
}

func realField(t *testing.T, lib *battle.Library, n int) *Field {
	t.Helper()
	return NewFieldFromTileLayers(lib.Tiles(n), lib.Heights(n), lib.TileLayers(n), lib.GateX(n))
}

// 地面層要與機器碼演算法算出來的一樣（docs/re/63 §2）：
// 低平面幾乎全是層 0，高平面（牆頂）正規化後只有一個高度。
func TestGroundPlanesMatchRawAlgorithm(t *testing.T) {
	lib := realLib(t)
	for _, n := range []int{0, 5, 20} {
		f := realField(t, lib, n)
		lowHist := map[int]int{}
		highHist := map[int]int{}
		gates := 0
		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				if lv, ok := f.GroundLevel(x, y, PlaneLow); ok {
					lowHist[lv]++
				}
				if lv, ok := f.GroundLevel(x, y, PlaneHigh); ok {
					highHist[lv]++
				}
				if f.IsGateCell(x, y) {
					gates++
				}
			}
		}
		// 存的是兵腳下佔的那一層，平地是 0。
		if lowHist[0] < 2000 {
			t.Errorf("戰場 %d：低平面平地只有 %d 格，太少", n, lowHist[0])
		}
		if len(highHist) != 1 {
			t.Errorf("戰場 %d：高平面有 %d 種高度 %v，正規化之後應該只剩一種",
				n, len(highHist), highHist)
		}
		if gates == 0 {
			t.Errorf("戰場 %d 是攻城圖，卻一個門格都沒有", n)
		}
		t.Logf("戰場 %d：低平面 %v／高平面 %v／門格 %d", n, lowHist, highHist, gates)
	}
}

// 牆頂只有 4 或 5 兩個可能（`sub_1BBA6` 的多數決）。
func TestHighPlaneIsNormalisedToOneLevel(t *testing.T) {
	lib := realLib(t)
	for _, n := range []int{0, 5, 20} {
		f := realField(t, lib, n)
		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				lv, ok := f.GroundLevel(x, y, PlaneHigh)
				if ok && lv != 4 && lv != 5 {
					t.Fatalf("戰場 %d (%d,%d) 高平面 ＝ %d，只能是 4 或 5", n, x, y, lv)
				}
			}
		}
	}
}

// 上下城牆只在門格（原版圖塊 ≥ 0xF0）。城壁本體不行。
func TestOnlyGateCellsAllowClimb(t *testing.T) {
	lib := realLib(t)
	f := realField(t, lib, 0)
	walls, gates := 0, 0
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			tile := f.Tile(x, y)
			switch {
			case tile >= 0xF0:
				gates++
				if !f.IsGateCell(x, y) {
					t.Fatalf("(%d,%d) 圖塊 %#x 是門，卻不能上下", x, y, tile)
				}
			case tile >= 0xD0 && tile <= 0xDF:
				walls++
				if f.IsGateCell(x, y) {
					t.Fatalf("(%d,%d) 圖塊 %#x 是城壁，不該能上下", x, y, tile)
				}
			}
		}
	}
	if walls == 0 || gates == 0 {
		t.Fatalf("戰場 0 應該同時有城壁與門，量到 %d／%d", walls, gates)
	}
}

// 導航位元的兩條規則：低平面 |Δ| ≤ 1、高平面要相等。
func TestNavBitsFollowLevelRules(t *testing.T) {
	lib := realLib(t)
	f := realField(t, lib, 0)
	// 只檢查可站範圍內的格子：導航表建在整張 64×62 上，
	// 而 GroundLevel 照原版把座標夾在 1–62（inBounds）。
	for y := MinCoord + 1; y < MaxY; y++ {
		for x := MinCoord + 1; x < MaxCoord; x++ {
			for _, d := range navSteps {
				if !f.Linked(x, y, PlaneLow, d.dx, d.dy) {
					continue
				}
				a, okA := f.GroundLevel(x, y, PlaneLow)
				b, okB := f.GroundLevel(x+d.dx, y+d.dy, PlaneLow)
				if !okA || !okB {
					t.Fatalf("(%d,%d)→(%+d,%+d) 連通但有一邊沒有地面", x, y, d.dx, d.dy)
				}
				if abs(a-b) > 1 {
					t.Fatalf("(%d,%d)→(%+d,%+d) 連通但高度差 %d", x, y, d.dx, d.dy, abs(a-b))
				}
			}
		}
	}
}

// ⭐ 城壁擋人的是**通行層**不是地面層表。
//
// 地面層表是拿「打壞後的圖塊」算的（`sub_1BCA6` 把 0xD0–0xDF 換成 +0x10，
// docs/re/63 §2），所以尋路會直接規劃穿過城壁——原版也是這樣。
// 真正擋住的是 `sub_1BB6D` 展開的通行層，**打壞之後才會通**。
func TestSiegeWallBlocksUntilBroken(t *testing.T) {
	lib := realLib(t)
	for _, n := range []int{0, 5, 20} {
		f := realField(t, lib, n)
		gx := f.GateX()
		y := Height / 2
		// 尋路穿得過去（地面層表已經是打壞後的樣子）。
		if pts := f.FindPath(Point{X: gx - 8, Y: y}, Point{X: gx + 8, Y: y}, true, nil); len(pts) == 0 {
			t.Errorf("戰場 %d：尋路應該規劃得出穿過城門那一列的路", n)
		}
		// 但站在牆外那一格，往牆裡走會被通行層擋住。
		blocked := false
		for x := gx - 3; x <= gx+3; x++ {
			if f.Blocked(x, y, 1) {
				blocked = true
				break
			}
		}
		if !blocked {
			t.Errorf("戰場 %d：城門那一列在打壞前應該有一格擋著", n)
		}
	}
}

