package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 一座城的圖形有 4×4 格，但**只有據點表登記的那一格按得到**
// （原版 `sub_1709E` 的第二道閘，docs/spec/149 §1）。
//
// ⭐ 「點圖示中間沒反應」是原版行為，不是 bug。
func TestCityAtTileNeedsExactRegisteredTile(t *testing.T) {
	w := &state.World{}
	w.Cities[7].X, w.Cities[7].Y = 206, 114
	g := &game{world: w}
	if got, ok := g.cityAtTile(206, 114); !ok || got != 7 {
		t.Fatalf("登記格 (206,114) ＝ %d,%v，want 7,true", got, ok)
	}
	for _, tc := range [][2]int{{207, 114}, {206, 115}, {205, 113}, {208, 116}} {
		if _, ok := g.cityAtTile(tc[0], tc[1]); ok {
			t.Errorf("(%d,%d) 不該命中——原版只認登記的那一格", tc[0], tc[1])
		}
	}
}

// 螢幕座標 → 格：鏡頭本來就是以格為單位，地圖區以外回 false。
func TestMapPickTileAt(t *testing.T) {
	g := &game{camX: 100, camY: 40}
	col, row, ok := g.mapPickTileAt(192, 112)
	if !ok || col != 112 || row != 45 {
		t.Errorf("(192,112) ＝ (%d,%d,%v)，want (112,45,true)", col, row, ok)
	}
	// 指令列（y < 32）與地圖右緣以外都不算。
	if _, _, ok := g.mapPickTileAt(192, 16); ok {
		t.Error("指令列上的座標不該算成地圖格")
	}
	if _, _, ok := g.mapPickTileAt(strategyMapW, 112); ok {
		t.Error("地圖右緣以外不該算成地圖格")
	}
}

// 游標的兩層遮罩逐像素取自原版擷取（docs/spec/149 §1.1）。
//
// ⭐ **黑影不是白框位移一格**：左下、右上的圓角步法不一樣，
// 位移法會差 6 px。這一支釘住那幾格。
func TestMapPickCursorMask(t *testing.T) {
	for i, line := range mapPickCursorWhite {
		if len(line) != 16 {
			t.Fatalf("白第 %d 列長度 %d，want 16", i, len(line))
		}
	}
	for i, line := range mapPickCursorBlack {
		if len(line) != 16 {
			t.Fatalf("黑第 %d 列長度 %d，want 16", i, len(line))
		}
	}
	// 白框的外框是 15×15：頭尾兩列縮一格，中間兩側各一格。
	if mapPickCursorWhite[0] != ".#############.." {
		t.Errorf("白第 0 列 ＝ %q", mapPickCursorWhite[0])
	}
	if mapPickCursorWhite[1] != "##...........##." {
		t.Errorf("白第 1 列 ＝ %q", mapPickCursorWhite[1])
	}
	// 影子與白框位移一格**不相等**的那幾格（左下 dy14 dx0、右上 dy0 dx14）。
	if mapPickCursorBlack[0][14] != 'K' {
		t.Error("右上角那一格黑影不見了")
	}
	if mapPickCursorBlack[14][0] != 'K' {
		t.Error("左下角那一格黑影不見了")
	}
}

// 選點期間時間停著（原版 `sub_1703C` 的迴圈不推進世界，docs/spec/149 §1.2）。
func TestMapPickStopsTime(t *testing.T) {
	g := &game{world: &state.World{}}
	if !g.timeRuns() {
		t.Fatal("什麼都沒開的時候時間本來就該跑")
	}
	g.beginMapPick(func(int) {}, nil)
	if g.timeRuns() {
		t.Error("選點期間時間不該跑")
	}
	g.endMapPick()
	if !g.timeRuns() {
		t.Error("選完之後時間要回來")
	}
}

// 選點期間指令格還亮著：畫面上一張視窗都沒有，但流程還在跑
// （docs/spec/149 §2）。漏掉的症狀是反白與狀態列整組消失。
func TestMapPickKeepsCommandCellLit(t *testing.T) {
	g := &game{world: &state.World{}, cmdCell: int(naturalCommandCorps)}
	// 對照：什麼都沒開的時候本來就會熄——`commandFlowRunning` 回 false。
	if got := g.activeCommandCell(); got != -1 {
		t.Fatalf("什麼都沒開卻還亮著：%d", got)
	}
	g.beginMapPick(func(int) {}, nil)
	if got := g.activeCommandCell(); got != int(naturalCommandCorps) {
		t.Errorf("選點期間指令格 ＝ %d，want %d", got, naturalCommandCorps)
	}
}
