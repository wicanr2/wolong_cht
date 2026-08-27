package world

import "testing"

// TestRoadGraphCountsArePinned 釘住道路圖的三個總量。
//
// ⭐ **為什麼需要這一支**：這三個數字曾經在文件裡待了十天而沒有人發現它們
// 已經不對。`docs/re/08` §7.7 記的是 253 條邊、總長 7,822 格、最長 94 格，
// 那是 2026-08-17「`MMAP.MAP` 開頭四個 byte 是長度欄位」修好**之前**量的；
// 那次修正把整張地圖右移四格，道路圖跟著變，而沒有任何檢查會紅。
//
// 邊數尤其危險：`CLAUDE.md` §7 第 19 條——**「多找到一條」比「少找到一條」
// 更容易被當成好消息**，所以更難被質疑。
//
// 數字變了不代表壞了；要改這裡的期望值，先確認變化有解釋，再一起更新
// `docs/re/08` §7.7。
func TestRoadGraphCountsArePinned(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	xy, _ := cityRecords(t)
	edges, err := RoadEdges(m, xy)
	if err != nil {
		t.Fatal(err)
	}

	total, longest := 0, 0
	for _, e := range edges {
		total += e.Steps
		if e.Steps > longest {
			longest = e.Steps
		}
	}
	if len(edges) != 254 {
		t.Errorf("邊數 = %d，期望 254", len(edges))
	}
	if total != 5852 {
		t.Errorf("路徑總長（含頭尾城內那一小段）= %d，期望 5852", total)
	}
	if longest != 86 {
		t.Errorf("最長一條 = %d 格，期望 86", longest)
	}
}

// TestHuiJiZhangAnEdgeIsRealNotWrapAround 釘住那條被誤判過的邊。
//
// 移植時踩過平面陣列的坑：`si ± 1` 在列邊界會繞到下一列的另一端，
// 因而**憑空生出一條邊**。當時的結論是「會稽–章安就是那條假的」，
// 但加上 x 位移限制之後它仍然在——因為它是真的。
//
// 判準不是「它在不在」，是**路徑有沒有繞行的指紋**：單步 x 位移超過 1。
func TestHuiJiZhangAnEdgeIsRealNotWrapAround(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	xy, _ := cityRecords(t)
	edges, err := RoadEdges(m, xy)
	if err != nil {
		t.Fatal(err)
	}

	const huiJi, zhangAn = 153, 173
	var found *RoadEdge
	for i := range edges {
		a, b := edges[i].A, edges[i].B
		if a > b {
			a, b = b, a
		}
		if a == huiJi && b == zhangAn {
			found = &edges[i]
			break
		}
	}
	if found == nil {
		t.Fatal("會稽–章安不在道路圖裡")
	}
	for i := 1; i < len(found.Path); i++ {
		if d := abs(found.Path[i][0] - found.Path[i-1][0]); d > 1 {
			t.Fatalf("第 %d 步的 x 位移 %d > 1，這是列邊界繞行的指紋", i, d)
		}
	}
	if end := found.Path[len(found.Path)-1]; end != xy[zhangAn] {
		t.Errorf("終點 %v 不是章安的城中心 %v", end, xy[zhangAn])
	}
}
