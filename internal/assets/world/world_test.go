package world

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/rle"
)

func read(t *testing.T, ver, name string) []byte {
	t.Helper()
	p := filepath.Join("..", "..", "..", "workplace", "orig", ver, name)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("找不到原版素材 %s，跳過", p)
	}
	return b
}

// TestMapSize 釘住 384×256 —— 這個數字是從 sub_1E4CE 的迴圈邊界讀出來的，
// 不是湊的。解壓後不足這個大小就代表 RLE 解錯了。
func TestMapSize(t *testing.T) {
	for _, ver := range []string{"dosv", "pc98"} {
		m, err := ParseMap(read(t, ver, "MMAP.MAP"))
		if err != nil {
			t.Fatalf("%s: %v", ver, err)
		}
		if len(m.Tiles) != Width*Height {
			t.Errorf("%s: 解出 %d 格，預期 %d", ver, len(m.Tiles), Width*Height)
		}
		// 尾巴幾個 byte 是原版解壓器「解到檔尾」的副產品，不該多到離譜。
		if len(m.Extra) > 64 {
			t.Errorf("%s: 尾巴 %d B 太多，RLE 可能解錯", ver, len(m.Extra))
		}
	}
}

// TestTileSet 釘住 MMAP.MDL = 256 塊 16×16，餘 0。
func TestTileSet(t *testing.T) {
	if _, err := ParseTileSet(read(t, "dosv", "MMAP.MDL")); err != nil {
		t.Error(err)
	}
}

// TestTilesLookSane 檢查解出來的地圖不是雜訊：
// 用到的圖塊種類要夠多（真實地圖會用上百種），
// 而且最常出現的那一種不該佔壓倒性多數（那代表解壓爆掉變成一片同值）。
func TestTilesLookSane(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	count := map[byte]int{}
	for _, v := range m.Tiles {
		count[v]++
	}
	if len(count) < 100 {
		t.Errorf("只用到 %d 種圖塊，太少，RLE 可能解錯", len(count))
	}
	max := 0
	for _, c := range count {
		if c > max {
			max = c
		}
	}
	if max > len(m.Tiles)/3 {
		t.Errorf("最常見的圖塊佔 %d/%d，超過三分之一 —— 解壓可能爆掉",
			max, len(m.Tiles))
	}
}

// TestRLERoundsToKnownSizes 兩版的 MMAP.MAP 不同（dosv 與 pc98 都有），
// 但解出來的格數必須一樣 —— 那是同一張地圖。
func TestRLESameShapeBothVersions(t *testing.T) {
	a := rle.Decode(read(t, "dosv", "MMAP.MAP"))
	b := rle.Decode(read(t, "pc98", "MMAP.MAP"))
	if len(a) != len(b) {
		t.Errorf("兩版解出的長度不同：dosv %d B、pc98 %d B", len(a), len(b))
	}
}
