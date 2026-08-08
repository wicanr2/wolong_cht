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

// cityRecords 從 SINARIO.DAT 取 192 個據點的座標與鄰接欄位。
// 這裡直接讀原始位元組而不經過 internal/state —— 資產層不該依賴規則層。
func cityRecords(t *testing.T) (xy [][2]int, known map[[2]int]bool) {
	t.Helper()
	b := read(t, "dosv", "SINARIO.DAT")
	const base, size, n = 0x8C0, 32, 192
	known = map[[2]int]bool{}
	for i := 0; i < n; i++ {
		r := b[base+i*size : base+i*size+size]
		xy = append(xy, [2]int{
			int(r[8]) | int(r[9])<<8,
			int(r[10]) | int(r[11])<<8,
		})
		for k := 0; k < 4; k++ {
			if r[0]>>k&1 == 0 {
				continue
			}
			if m := int(r[0x1C+k]); m < n {
				a, c := i, m
				if a > c {
					a, c = c, a
				}
				known[[2]int{a, c}] = true
			}
		}
	}
	return xy, known
}

// ⭐ 道路圖的三個正對照。這張圖是**推導**出來的（多源 BFS），
// 不是從檔案裡讀出來的，所以它的正確性完全靠這三條：
//
//	① 據點記錄裡那 85 條鄰接必須全部出現
//	② 192 個據點必須全連通
//	③ 每個據點的分支度 ≤ 4 —— 據點記錄只有四個方向槽
//
// ①③ 是原版資料自己說的話，②是「遊戲要能玩」的下限。
// 三條同時成立，推導的參數（道路值域、橋、8-連通、節點偏移）才算對。
func TestRoadGraphAgainstCityRecords(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	xy, known := cityRecords(t)
	edges, err := RoadEdges(m, xy)
	if err != nil {
		t.Fatal(err)
	}

	have := map[[2]int]bool{}
	deg := make([]int, len(xy))
	adj := make([][]int, len(xy))
	for _, e := range edges {
		have[[2]int{e.A, e.B}] = true
		deg[e.A]++
		deg[e.B]++
		adj[e.A] = append(adj[e.A], e.B)
		adj[e.B] = append(adj[e.B], e.A)
	}

	// ① 已知邊全中
	missing := 0
	for k := range known {
		if !have[k] {
			missing++
		}
	}
	if missing != 0 {
		t.Errorf("據點記錄的 %d 條鄰接有 %d 條不在推導的道路圖裡",
			len(known), missing)
	}

	// ② 全連通
	seen := map[int]bool{0: true}
	stack := []int{0}
	for len(stack) > 0 {
		x := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, y := range adj[x] {
			if !seen[y] {
				seen[y] = true
				stack = append(stack, y)
			}
		}
	}
	if len(seen) != len(xy) {
		t.Errorf("只有 %d/%d 個據點連通", len(seen), len(xy))
	}

	// ③ 分支度 ≤ 4
	for i, d := range deg {
		if d > 4 {
			t.Errorf("據點 %d 的分支度 %d > 4（記錄只有四個方向槽）", i, d)
		}
	}
}

// 兩版的 MMAP.MAP 是 byte-for-byte 相同的（CLAUDE.md §3.10），
// 所以推出來的道路圖也必須一模一樣。這條在防「解碼路徑摻進版本相依」。
func TestRoadGraphIdenticalAcrossVersions(t *testing.T) {
	xy, _ := cityRecords(t)
	var got [2][]RoadEdge
	for i, ver := range []string{"dosv", "pc98"} {
		m, err := ParseMap(read(t, ver, "MMAP.MAP"))
		if err != nil {
			t.Fatal(err)
		}
		if got[i], err = RoadEdges(m, xy); err != nil {
			t.Fatal(err)
		}
	}
	if len(got[0]) != len(got[1]) {
		t.Fatalf("兩版邊數不同：%d vs %d", len(got[0]), len(got[1]))
	}
	for i := range got[0] {
		if got[0][i] != got[1][i] {
			t.Fatalf("第 %d 條邊不同：%+v vs %+v", i, got[0][i], got[1][i])
		}
	}
}
