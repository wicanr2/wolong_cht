package main

import (
	"image"
	"testing"
)

// 選君主那一頁的滑鼠**只認卡片自己的兩個熱區**（docs/spec/27 §2.1）。
//
// ⚠ 這一條擋的是一個實際發生過的 bug：那一頁畫的是卡片，
// 但命中測試沿用啟動殼層的清單列，而那些列**一列都沒畫出來**、
// 範圍又幾乎蓋滿整張卡片——滑鼠一移過去就換君主、一點下去就直接決定。
func TestLordCardIgnoresPointerOutsideItsTwoHotspots(t *testing.T) {
	l := &launcherModel{
		phase: launcherSelectPlayer,
		players: []launcherPlayer{
			{ID: 0, Lord: "曹操"}, {ID: 1, Lord: "劉備"}, {ID: 2, Lord: "孫策"},
		},
	}

	// ① 卡片上任何一點都不該被當成「清單的第 N 列」。
	//    取樣點刻意含舊的幽靈區（x 128–512、y 88–280）四個角與中央。
	for _, p := range []image.Point{
		{X: 130, Y: 90}, {X: 500, Y: 90}, {X: 130, Y: 270}, {X: 500, Y: 270},
		{X: 280, Y: 200}, // 卡片正中央
		{X: 330, Y: 274}, // 連「確定」那一格也不該回列號
	} {
		if row, ok := l.pointerRow(p.X, p.Y); ok {
			t.Errorf("(%d,%d) 仍被當成清單第 %d 列", p.X, p.Y, row)
		}
	}

	// ② 兩個熱區各自命中，位置與原版一致（0x20 在下、0x21 在上）。
	for _, tc := range []struct {
		name string
		x, y int
		want lordCardHotspot
	}{
		{"確定", lordOKX + 1, lordOKY + 1, lordCardConfirm},
		{"自定", lordCustomX + 1, lordCustomY + 1, lordCardCustom},
		{"兩顆按鈕之間", lordOKX + 1, lordCustomY + lordButtonH + 1, lordCardNone},
		{"卡片中央", 280, 200, lordCardNone},
		{"按鈕右邊一格", lordOKX + lordButtonW, lordOKY + 1, lordCardNone},
	} {
		if got := lordCardHotspotAt(tc.x, tc.y); got != tc.want {
			t.Errorf("%s (%d,%d) = %v，want %v", tc.name, tc.x, tc.y, got, tc.want)
		}
	}

	// ③ 反向對照：**別的階段仍然指得到列**，否則「永遠回 false」也會通過 ①。
	l.phase = launcherTitle
	r := launcherRowRect(launcherTitle, 0)
	if _, ok := l.pointerRow(r.Min.X+1, r.Min.Y+1); !ok {
		t.Error("標題頁的第一列指不到了——這條修正不該影響其他階段")
	}
}

// 「自定」按下去只留一句提示，不會把君主決定下去（docs/spec/27 §5：卡在存檔）；
// 卡片其他位置一點反應都沒有；只有「確定」會往下走。
func TestLordCardHotspotActions(t *testing.T) {
	newGame := func() *game {
		return &game{launcher: &launcherModel{
			phase:           launcherSelectPlayer,
			players:         []launcherPlayer{{ID: 3, Lord: "曹操"}},
			confirmedPlayer: -1,
		}}
	}

	g := newGame()
	if err := g.applyLordCardHotspot(lordCardCustom); err != nil {
		t.Fatal(err)
	}
	if g.launcher.phase != launcherSelectPlayer || g.launcher.confirmedPlayer != -1 {
		t.Error("「自定」把君主決定下去了")
	}
	if g.launcher.notice == "" {
		t.Error("「自定」沒有留下任何提示，玩家會以為按鈕壞了")
	}

	g = newGame()
	if err := g.applyLordCardHotspot(lordCardNone); err != nil {
		t.Fatal(err)
	}
	if g.launcher.phase != launcherSelectPlayer || g.launcher.notice != "" {
		t.Error("點在卡片空白處不該有任何反應")
	}

	g = newGame()
	if err := g.applyLordCardHotspot(lordCardConfirm); err != nil {
		t.Fatal(err)
	}
	if g.launcher.confirmedPlayer != 3 || g.launcher.phase != launcherGameConfirm {
		t.Errorf("「確定」沒有往下走：player=%d phase=%v",
			g.launcher.confirmedPlayer, g.launcher.phase)
	}
}
