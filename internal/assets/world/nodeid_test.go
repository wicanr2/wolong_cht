package world

import "testing"

// 行軍節點編號恆在 0–191（據點）——「路上的中間節點 192–255」不存在。
//
// 原版把終點節點編號直接讀「walk 停下那一格在影子緩衝區的蓋章值」
// （`sub_1E77D` 的 `mov bl, es:[si]`），而活著的蓋章路徑只有
// `sub_1E57F` → `sub_1E68C`：蓋的是據點流水號 0–191、蓋在城門格。
// 會蓋出其他值的那段程式碼（`0x1E50D`–`0x1E566`）沒有任何呼叫者，
// 是死碼。所以只要**每一條 walk 都停在某城的城門格**，
// 節點編號就恆在 0–191——這支測試釘住這個前提。
//
// 若這支開始失敗（某條 walk 停在未蓋章的節點格），代表原版會把
// 圖塊值 0xCB–0xD3 當節點編號用，`docs/mechanics/20` §「行軍」的
// 節點結論要重新查。
func TestRoadWalksAllEndAtGateCells(t *testing.T) {
	m, err := ParseMap(read(t, "dosv", "MMAP.MAP"))
	if err != nil {
		t.Fatal(err)
	}
	xy, _ := cityRecords(t)
	tiles := m.Tiles

	nodeOf := make([]int, len(xy))
	for ci, c := range xy {
		nodeOf[ci] = c[1]*Width + c[0]
	}
	gateCity := map[int]int{}
	type start struct{ gate, dir int }
	var starts []start
	for ci, node := range nodeOf {
		for d, probe := range gateProbe {
			for _, off := range probe {
				q := node + off
				if q < 0 || q >= Width*Height || !isRoad(tiles[q]) {
					continue
				}
				gateCity[q] = ci
				starts = append(starts, start{gate: q, dir: d})
				break
			}
		}
	}

	bad := 0
	for _, s := range starts {
		_, end := walkRoad(tiles, s.gate, s.dir)
		if end < 0 || end >= len(tiles) {
			continue // 走出界＝死路，原版同樣不建這條連結
		}
		if _, ok := gateCity[end]; ok {
			continue
		}
		// 停在非城門格：只有「未蓋章的節點格／野地」才會走到這裡。
		bad++
		if bad <= 5 {
			t.Errorf("walk 從格 %d 出發停在非城門格 %d（圖塊 0x%02X）",
				s.gate, end, tiles[end])
		}
	}
	if bad != 0 {
		t.Fatalf("%d 條 walk 停在非城門格——節點編號 0–191 的前提被打破", bad)
	}
}
