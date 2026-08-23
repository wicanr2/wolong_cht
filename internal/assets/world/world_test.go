package world

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/palette"
	"github.com/wicanr2/wolong_cht/internal/assets/rle"
)

// loadWorldForTest 把畫一張大地圖需要的四樣東西讀進來。
// 缺任何一樣就 skip——原版素材不進版控。
func loadWorldForTest(t *testing.T) (*Map, *TileSet, *MCH, *palette.Palette) {
	t.Helper()
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	ts, err := ParseTileSet(read(t, "dosv", "MMAP.MDL"))
	if err != nil {
		t.Fatal(err)
	}
	mch, err := ParseMCH(read(t, "dosv", "MMAP.MCH"))
	if err != nil {
		t.Fatal(err)
	}
	pal, err := palette.Parse(read(t, "dosv", "GAMEPAL.BRG"))
	if err != nil {
		t.Fatal(err)
	}
	return m, ts, mch, pal
}

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

// TestMCHObjectLayout 固定 sub_1D804 的 256×160 物件圖塊區，以及
// sub_12533 的 CS:985Ah object type/frame 查表。這不是把檔案尾端湊成圖，
// 而是同時檢查 IDA 的 metadata 讀法與實際 MMAP.MCH source byte。
func TestMCHObjectLayout(t *testing.T) {
	m, err := ParseMCH(read(t, "dosv", "MMAP.MCH"))
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []byte{0, 0x80, 0xFF} {
		tile := m.Tile(id)
		if tile == nil {
			t.Fatalf("MCH 圖塊 0x%02X 解不出來", id)
		}
		opaque, transparent := 0, 0
		for _, px := range tile.Pix {
			if px == MCHTransparent {
				transparent++
			} else if px < 16 {
				opaque++
			} else {
				t.Fatalf("MCH 圖塊 0x%02X 出現非法色號 0x%02X", id, px)
			}
		}
		if opaque == 0 || transparent == 0 {
			t.Fatalf("MCH 圖塊 0x%02X 沒有同時呈現 mask 的不透明／透明像素", id)
		}
	}

	for _, tc := range []struct {
		objectType, frame, index, width, height int
	}{
		{1, 0, 0x18, 16, 9},
		{1, 4, 0x1C, 16, 9},
		{1, 7, 0x1A, 16, 9},
		{2, 0, 0x20, 5, 5},
		{2, 7, 0x23, 5, 5},
		{3, 0, 0x28, 5, 5},
	} {
		index, ok := ObjectPatternIndex(tc.objectType, tc.frame)
		if !ok || index != tc.index {
			t.Fatalf("object type %d frame %d 查到 0x%X，預期 0x%X",
				tc.objectType, tc.frame, index, tc.index)
		}
		pattern, ok := m.Pattern(index)
		if !ok || pattern.Width != tc.width || pattern.Height != tc.height {
			t.Fatalf("pattern 0x%X = %dx%d，預期 %dx%d",
				tc.index, pattern.Width, pattern.Height, tc.width, tc.height)
		}
	}
	p, ok := m.PatternFor(1, 0)
	if !ok || len(p.Tiles) != 16*9 || p.Tiles[4] != 0xD0 {
		t.Fatalf("火災第 0 相位的 source 矩陣不符：ok=%v len=%d tile[4]=0x%02X",
			ok, len(p.Tiles), p.Tiles[4])
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

// 兩版的 MMAP.MAP 是 byte-for-byte 相同的（docs/re/01 §2，兩邊都是 80,716 B），
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
		a, b := got[0][i], got[1][i]
		if a.A != b.A || a.B != b.B || a.Steps != b.Steps ||
			len(a.Path) != len(b.Path) {
			t.Fatalf("第 %d 條邊不同：%d–%d/%d 格 vs %d–%d/%d 格",
				i, a.A, a.B, a.Steps, b.A, b.B, b.Steps)
		}
		for k := range a.Path {
			if a.Path[k] != b.Path[k] {
				t.Fatalf("第 %d 條邊的第 %d 格不同：%v vs %v",
					i, k, a.Path[k], b.Path[k])
			}
		}
	}
}

// ⭐ 路徑必須真的踩在路上。
//
// 推導只保證「兩端是據點、長度是最短」，**不保證中間每一格都是道路**——
// 如果 parent 鏈接錯邊，路徑會穿過山river，而總長度看起來還是對的。
// 這條逐格檢查圖塊值。
func TestRoadPathStaysOnRoad(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	xy, _ := cityRecords(t)
	edges, err := RoadEdges(m, xy)
	if err != nil {
		t.Fatal(err)
	}
	offRoad, checked := 0, 0
	for _, e := range edges {
		// 兩端各有一小段「城中心 → 節點格 → 城門格」，那幾格踩在
		// 城池圖形上而不是道路上。StubA／StubB 記著長度，照它排除——
		// **不要用固定的邊界值**，那會在城圖大小不同時失準。
		if len(e.Path) <= e.StubA+e.StubB {
			continue
		}
		for _, c := range e.Path[e.StubA : len(e.Path)-e.StubB] {
			checked++
			if !isRoad(m.Tiles[c[1]*Width+c[0]]) {
				offRoad++
			}
		}
	}
	if checked == 0 {
		t.Fatal("沒有檢查到任何格子")
	}
	if offRoad != 0 {
		t.Errorf("%d/%d 格不在道路上", offRoad, checked)
	}
}

// 路徑要是**連續**的：相鄰兩格必須 8-連通。
// 斷開的話軍團會瞬移，而總步數看起來還是對的。
func TestRoadPathIsContiguous(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	xy, _ := cityRecords(t)
	edges, err := RoadEdges(m, xy)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range edges {
		prev := xy[e.A]
		for i, c := range e.Path {
			dx, dy := c[0]-prev[0], c[1]-prev[1]
			if dx < -1 || dx > 1 || dy < -1 || dy > 1 || (dx == 0 && dy == 0) {
				t.Fatalf("邊 %d–%d 的第 %d 格從 %v 跳到 %v", e.A, e.B, i, prev, c)
			}
			prev = c
		}
		if prev != xy[e.B] {
			t.Fatalf("邊 %d–%d 的終點是 %v，應該是 %v", e.A, e.B, prev, xy[e.B])
		}
	}
}

// ⭐ 城門格的數量必須等於「類 4」圖塊的數量。
//
// 城門格是用「節點格的四個方向、先 1 格再 2 格」找出來的（sub_1E57F），
// 而類 4（0xD4–0xDD）是 sub_1E961 獨立分出來的一類。
// **兩邊算法完全無關，數字卻應該一樣**——
// 對得上就代表「類 4 ＝ 城門格」這個讀法是對的，也代表沒有漏找。
func TestGateCellsMatchClassFour(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	xy, _ := cityRecords(t)

	classFour := 0
	for _, v := range m.Tiles[:Width*Height] {
		if tileClass(v) == 4 {
			classFour++
		}
	}

	// 照 RoadEdges 的方式數一次城門格。
	gates := map[int]bool{}
	for _, c := range xy {
		node := -1
		for _, dx := range [...]int{nodeDX, -nodeDX} {
			if v := m.Tiles[c[1]*Width+c[0]+dx]; v >= nodeLo && v <= nodeHi {
				node = c[1]*Width + c[0] + dx
				break
			}
		}
		if node < 0 {
			t.Fatalf("據點 (%d,%d) 找不到節點格", c[0], c[1])
		}
		for _, probe := range gateProbe {
			for _, off := range probe {
				if q := node + off; q >= 0 && q < Width*Height && isRoad(m.Tiles[q]) {
					gates[q] = true
					break
				}
			}
		}
	}
	if len(gates) != classFour {
		t.Errorf("城門格 %d 個、類 4 圖塊 %d 格，兩者應相等",
			len(gates), classFour)
	}
}

// 路徑不能比兩點直線距離短——短了代表某處穿牆或跳格。
//
// 這條抓的是「連續性檢查放過去、但整體走向錯了」的情況：
// 逐格連續只保證每一步合法，不保證整條路合理。
func TestRoadPathNotShorterThanStraightLine(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	xy, _ := cityRecords(t)
	edges, err := RoadEdges(m, xy)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range edges {
		dx := xy[e.A][0] - xy[e.B][0]
		dy := xy[e.A][1] - xy[e.B][1]
		if dx < 0 {
			dx = -dx
		}
		if dy < 0 {
			dy = -dy
		}
		straightLine := dx
		if dy > dx {
			straightLine = dy // 切比雪夫距離 ＝ 8 方向的下限
		}
		if len(e.Path) < straightLine {
			t.Errorf("邊 %d–%d 只有 %d 格，直線下限是 %d",
				e.A, e.B, len(e.Path), straightLine)
		}
	}
}

// ⭐ 圖塊算式釘死：勢力 × 5 ＋ 朝向（docs/spec/74 §3，原版 sub_12B2A）。
// 22 個勢力 × 5 張 ＝ 110 張，最後一張是 109。
func TestCorpsTile(t *testing.T) {
	cases := []struct {
		faction, heading int
		want             byte
	}{
		{0, 0, 0}, {0, 4, 4},
		{1, 0, 5}, {1, 3, 8},
		{21, 4, 109},
	}
	for _, c := range cases {
		if got := CorpsTile(c.faction, c.heading); got != c.want {
			t.Errorf("CorpsTile(%d,%d) = %d，want %d", c.faction, c.heading, got, c.want)
		}
	}
	// ⚠ 朝向越界要退回「靜止」，**不可以**溢位成別的勢力那五張。
	for _, h := range []int{-1, 5, 99} {
		if got := CorpsTile(3, h); got != CorpsTile(3, CorpsHeadings-1) {
			t.Errorf("朝向 %d 沒有退回靜止：得到 %d", h, got)
		}
	}
	// 兩個勢力的圖塊區間不可以重疊。
	for f := 0; f < 21; f++ {
		if CorpsTile(f, CorpsHeadings-1) >= CorpsTile(f+1, 0) {
			t.Fatalf("勢力 %d 與 %d 的圖塊撞了", f, f+1)
		}
	}
}

// 軍團疊圖真的會蓋在地形上，而且死掉的軍團不畫。
func TestRenderMarkedDrawsCorps(t *testing.T) {
	m, ts, mch, pal := loadWorldForTest(t)
	base, err := m.RenderMarked(ts, mch, pal, 0, 0, 0, 4, 4, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	with, err := m.RenderMarked(ts, mch, pal, 0, 0, 0, 4, 4, nil,
		[]CorpsMark{{X: 1, Y: 1, Tile: CorpsTile(0, 4)}})
	if err != nil {
		t.Fatal(err)
	}
	diff := 0
	for i := range base.Pix {
		if base.Pix[i] != with.Pix[i] {
			diff++
		}
	}
	if diff == 0 {
		t.Fatal("疊了軍團卻一個像素都沒變")
	}
	// ⚠ 只能改到那一格：16×16×4 ＝ 1024 個 byte 是上限。
	if diff > TileSize*TileSize*4 {
		t.Errorf("軍團疊圖改了 %d 個 byte，超出一格的範圍", diff)
	}
	// ⭐ 像素有變不等於畫對了。設 WOLONG_DUMP_DIR 就把圖存出來用眼睛看——
	// 畫面的錯測試看不到（CLAUDE.md §7 第 13 條）。
	dumpPNG(t, with, "corps-overlay.png")
}

// dumpPNG 在 WOLONG_DUMP_DIR 有設時把圖寫出來，方便肉眼複驗。
func dumpPNG(t *testing.T, img *image.RGBA, name string) {
	t.Helper()
	dir := os.Getenv("WOLONG_DUMP_DIR")
	if dir == "" {
		return
	}
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}
