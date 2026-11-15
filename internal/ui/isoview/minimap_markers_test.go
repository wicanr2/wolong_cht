package isoview

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

func minimapMarkerView(player int) (*View, *tactical.Battle) {
	rows := make([][]byte, 64)
	for y := range rows {
		rows[y] = make([]byte, 64)
	}
	rows[26][30] = 0xD0
	heights := new([256]int)
	f := tactical.NewFieldFromTiles(rows, heights, 27)
	b := tactical.NewBattle(f, nil, nil, 116)
	b.SetPlayerSide(player)
	b.Sides[player].Line = 5
	v := &View{observed: b, cursorX: 32, cursorY: 33}
	return v, b
}

func TestMinimapLaterEraseCanHideEarlierSoldier(t *testing.T) {
	for _, player := range []int{0, 1} {
		v, _ := minimapMarkerView(player)
		// 先處理的兵移入(5,10)，後處理的兵仍須擦除自己的舊(5,10)。
		v.EraseSoldierDrawing(4, 10, 0)
		v.DrawSoldierDrawing(0, tactical.SoldierDrawing{X: 5, Y: 10})
		v.EraseSoldierDrawing(5, 10, 0)
		v.DrawSoldierDrawing(1, tactical.SoldierDrawing{X: 6, Y: 10})
		marks := v.MinimapMarkers()
		if marks[10][5] != 0 || marks[10][6] != 2 {
			t.Fatalf("後續擦除不應被存活名單補回，玩家%d：%v", player, marks[10][4:7])
		}
		// 停在原格時，擦除後仍會畫回本兵的角色色。
		v.EraseSoldierDrawing(6, 10, 0)
		v.DrawSoldierDrawing(0, tactical.SoldierDrawing{X: 6, Y: 10})
		if v.MinimapMarkers()[10][6] != 1 {
			t.Fatal("停在原格也要完成登錄")
		}
	}
}

func TestMinimapBrokenTerrainRedrawsCrossAndLineInOrder(t *testing.T) {
	v, b := minimapMarkerView(1)
	for _, xy := range [][2]int{{30, 26}, {5, 7}, {32, 20}, {20, 33}, {12, 16}} {
		v.DrawSoldierDrawing(0, tactical.SoldierDrawing{X: xy[0], Y: xy[1]})
	}
	b.Field.Retile(30, 26, 0x10)
	// 原版先重畫破損格、十字及陣形線，才畫攻擊者自己的兵點。
	v.DrawSoldierDrawing(1, tactical.SoldierDrawing{X: 30, Y: 26})
	marks := v.MinimapMarkers()
	if marks[26][30] != 2 || marks[7][5] != 0 || marks[20][32] != 0 || marks[33][20] != 0 || marks[16][12] != 1 {
		t.Fatal("破損後的線條重畫必須覆蓋已畫兵點，其他位置保留")
	}
	v.DrawSoldierDrawing(0, tactical.SoldierDrawing{X: 5, Y: 7})
	if v.MinimapMarkers()[7][5] != 1 {
		t.Fatal("同拍較晚處理的兵應能重新蓋回陣形線")
	}
}

func TestMinimapDeathDrawingKeepsMarkerUntilRemoval(t *testing.T) {
	v, _ := minimapMarkerView(0)
	v.DrawSoldierDrawing(0, tactical.SoldierDrawing{X: 20, Y: 25})
	v.EraseSoldierDrawing(20, 25, 0)
	v.DrawDeathDrawing(0, tactical.Death{X: 20, Y: 25, FramesLeft: 3})
	if v.MinimapMarkers()[25][20] != 1 {
		t.Fatal("倒地繪圖不得提早取代1B3B2的兵點清除")
	}
	v.EraseSoldierDrawing(20, 25, 0)
	if v.MinimapMarkers()[25][20] != 0 {
		t.Fatal("最終移除必須擦除兵點")
	}
}
