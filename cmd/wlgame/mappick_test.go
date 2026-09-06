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

// 那一格上的軍團**不分敵我**，而且只認完全相等的格
// （原版建表 callback `0x17217`，docs/spec/151 §1.1）。
func TestCorpsAtTileListsBothSides(t *testing.T) {
	w := &state.World{Player: 0}
	w.Corps[3].Alive, w.Corps[3].Faction, w.Corps[3].X, w.Corps[3].Y = true, 0, 206, 114
	w.Corps[9].Alive, w.Corps[9].Faction, w.Corps[9].X, w.Corps[9].Y = true, 7, 206, 114
	w.Corps[5].Alive, w.Corps[5].Faction, w.Corps[5].X, w.Corps[5].Y = true, 0, 207, 114
	w.Corps[6].Faction, w.Corps[6].X, w.Corps[6].Y = 0, 206, 114 // 沒活著
	g := &game{world: w}
	got := g.corpsAtTile(206, 114)
	if len(got) != 2 || got[0] != 3 || got[1] != 9 {
		t.Fatalf("那一格上的軍團 ＝ %v，want [3 9]（含別人的那一支）", got)
	}
	if n := len(g.corpsAtTile(208, 114)); n != 0 {
		t.Errorf("隔壁格不該有 %d 支", n)
	}
}

// 兩項選單的位置照 `sub_11F0E` 的兩道夾制（docs/spec/151 §1.2）。
// ⚠ 與行軍三選一**不是同一組數字**。
func TestMapChoiceAnchorClamps(t *testing.T) {
	if x, y := mapChoiceAnchor(100, 100); x != 96 || y != 96 {
		t.Errorf("(100,100) ＝ (%d,%d)，want (96,96)", x, y)
	}
	x, y := mapChoiceAnchor(639, 399)
	if x != mapChoiceMaxCol*16 || y != mapChoiceMaxRow*16 {
		t.Errorf("右下角 ＝ (%d,%d)，want (%d,%d)",
			x, y, mapChoiceMaxCol*16, mapChoiceMaxRow*16)
	}
	// 直的那道夾住之後下緣剛好貼齊畫面（兩項 ⇒ 框高 48）。
	if y+48 != screenH {
		t.Errorf("夾住之後下緣 ＝ %d，want %d", y+48, screenH)
	}
	// ⚠ **橫的那道不貼齊**：0x23 × 16 ＋ 112 ＝ 672，超出 640 共 32 px。
	// 行軍三選一的 0x21 才剛好 640（docs/spec/39 §3.5）。
	// 兩個立即值都是機器碼直接讀到的，這裡只釘住「不一樣」這件事，
	// 不替原版補一個它沒寫的規則（docs/spec/151 §4）。
	if mapChoiceMaxCol == marchModeMaxCol {
		t.Error("兩張選單的橫向夾制在原版是不同的值")
	}
}

// 大地圖點擊的三條分支（docs/spec/151 §1）：只有據點 → 情報卡、
// 只有軍團 → 那一格上的軍團一覽、兩者都有 → 先跳兩項選單。
func TestMapClickDispatch(t *testing.T) {
	newGame := func() *game {
		w := &state.World{Player: 0}
		w.Cities[7].X, w.Cities[7].Y = 206, 114
		w.Cities[8].X, w.Cities[8].Y = 100, 50
		w.Corps[3].Alive, w.Corps[3].Faction = true, 0
		w.Corps[3].X, w.Corps[3].Y = 206, 114
		w.Corps[4].Alive, w.Corps[4].Faction = true, 3
		w.Corps[4].X, w.Corps[4].Y = 60, 60
		return &game{world: w, cmdCell: -1}
	}

	// 只有據點。
	g := newGame()
	if !g.dispatchMapClick(100, 50, 0, 0) {
		t.Fatal("點到據點卻回 false")
	}
	if !g.cityInfo.active || g.cityInfo.city != 8 {
		t.Errorf("據點情報卡 ＝ %#v", g.cityInfo)
	}

	// 只有軍團（別人的）。
	g = newGame()
	if !g.dispatchMapClick(60, 60, 0, 0) {
		t.Fatal("點到軍團卻回 false")
	}
	if g.list == nil || len(g.list.Rows) != 1 || g.list.Rows[0] != 4 {
		t.Fatalf("那一格上的軍團一覽 ＝ %#v", g.list)
	}

	// 兩者都有 → 兩項選單。
	g = newGame()
	if !g.dispatchMapClick(206, 114, 0, 0) {
		t.Fatal("點到又是據點又有軍團的那一格卻回 false")
	}
	if !g.mapChoice.active || g.mapChoice.city != 7 ||
		len(g.mapChoice.corps) != 1 || g.mapChoice.corps[0] != 3 {
		t.Fatalf("兩項選單 ＝ %#v", g.mapChoice)
	}
	if g.cityInfo.active || g.list != nil {
		t.Error("跳選單的那一刻不該同時開別的東西")
	}

	// 空地：沒反應。
	g = newGame()
	if g.dispatchMapClick(1, 1, 0, 0) {
		t.Error("空地不該吃掉點擊")
	}
}
