package main

import (
	"image"
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 幾何全部照 `sub_17B3C` 與兩個 callback（docs/spec/79 §1）。
func TestFactionListGeometry(t *testing.T) {
	if factionListWinX != 136 || factionListWinY != 104 ||
		factionListWinW != 384 || factionListWinH != 176 {
		t.Fatalf("視窗 = (%d,%d) %dx%d，want (136,104) 384x176",
			factionListWinX, factionListWinY, factionListWinW, factionListWinH)
	}
	// ⚠ 尺寸與戰略層那七個一覽表相同，**左上角不同**（docs/re/26 §1）。
	if factionListWinW != listWinW || factionListWinH != listWinH {
		t.Error("尺寸應與一覽表相同")
	}
	if factionListWinX == listWinX && factionListWinY == listWinY {
		t.Error("左上角不該與一覽表相同——那正是先前被寫成「全部一樣」的那一格")
	}

	for _, tc := range []struct {
		name string
		got  image.Rectangle
		want image.Rectangle
	}{
		{"標題列", factionListTitleRect(), image.Rect(136, 104, 520, 120)},
		{"清單本體", factionListBodyRect(), image.Rect(152, 120, 520, 280)},
		{"▲", factionListScrollUpRect(), image.Rect(136, 120, 152, 136)},
		{"捲軸槽", factionListScrollTrackRect(), image.Rect(136, 136, 152, 264)},
		{"▼", factionListScrollDownRect(), image.Rect(136, 264, 152, 280)},
		{"第一列", factionListRowRect(0), image.Rect(152, 120, 520, 136)},
		{"第十列", factionListRowRect(9), image.Rect(152, 264, 520, 280)},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %v，want %v", tc.name, tc.got, tc.want)
		}
	}
	if factionListRowX() != 160 {
		t.Errorf("列首 X = %d，want 160（視窗 +24）", factionListRowX())
	}
}

// ⭐ 兩份獨立證據要對得上：欄位表 `cs:7B12` 的 X，與**標題字串逐字量出來**
// 的欄位落點。對得上才表示標題與列共用同一個列首（docs/spec/79 §1.4）。
func TestFactionListColumnsMatchTheHeaderString(t *testing.T) {
	for _, tc := range []struct {
		field string
		want  int
	}{
		{"勢力名", factionColLordX},
		{"軍師名", factionColAdvisorX},
		{"武將", factionColGeneralsX},
		{"據點", factionColCitiesX},
		{"首都", factionColCapitalX},
	} {
		i := strings.Index(factionListTitle, tc.field)
		if i < 0 {
			t.Fatalf("標題裡找不到 %q", tc.field)
		}
		got := textdraw.StringWidth(factionListTitle[:i]) + factionListRowDX
		if got != tc.want {
			t.Errorf("%s 的標題落點 = %d，欄位表寫 %d", tc.field, got, tc.want)
		}
	}
	// 反向對照：分隔線的總寬要等於最後一欄的右緣，否則欄寬定義與欄位表不同源。
	if w := textdraw.StringWidth(factionListDash); w != factionColCapitalX-factionListRowDX+48 {
		t.Errorf("分隔線寬 = %d，want %d", w, factionColCapitalX-factionListRowDX+48)
	}
}

func TestFactionListPointerRow(t *testing.T) {
	l := &launcherModel{phase: launcherSelectFaction}
	for i := 0; i < 22; i++ {
		l.players = append(l.players, launcherPlayer{ID: i})
	}

	// 沒捲動時：第 0 列與第 9 列。
	if n, ok := l.factionListSelectAt(200, 121); !ok || n != 0 {
		t.Errorf("第一列 = (%d,%v)，want (0,true)", n, ok)
	}
	if n, ok := l.factionListSelectAt(200, 279); !ok || n != 9 {
		t.Errorf("第十列 = (%d,%v)，want (9,true)", n, ok)
	}

	// 標題列與捲軸都不算選列。
	for _, p := range []image.Point{{X: 200, Y: 110}, {X: 140, Y: 200}, {X: 200, Y: 290}} {
		if n, ok := l.factionListSelectAt(p.X, p.Y); ok {
			t.Errorf("(%d,%d) 不該選到第 %d 列", p.X, p.Y, n)
		}
	}

	// 捲動之後同一個 y 對到不同的勢力。
	l.factionTop = 5
	if n, ok := l.factionListSelectAt(200, 121); !ok || n != 5 {
		t.Errorf("捲到 5 之後第一列 = (%d,%v)，want (5,true)", n, ok)
	}
	// 捲到底就夾住：22 筆、一頁 10 列 ⇒ top 最多 12。
	l.factionTop = 99
	if got := l.factionListTop(); got != 12 {
		t.Errorf("捲到底 = %d，want 12", got)
	}

	// 只有 3 個勢力時，後面七列是空的，點下去選不到東西。
	short := &launcherModel{phase: launcherSelectFaction,
		players: []launcherPlayer{{ID: 0}, {ID: 1}, {ID: 2}}}
	if n, ok := short.factionListSelectAt(200, 200); ok {
		t.Errorf("空列不該選到第 %d 列", n)
	}
}

// 流程：劇本 → 清單 → 君主卡，取消各退一層（docs/re/73 §1）。
func TestLauncherFlowGoesThroughFactionList(t *testing.T) {
	l := &launcherModel{phase: launcherScenario, confirmedPlayer: -1}
	players := []launcherPlayer{{ID: 0, Lord: "曹操"}, {ID: 4, Lord: "孫策"}}
	if !l.setScenarioPlayers(0, "第一章", players) {
		t.Fatal("setScenarioPlayers 失敗")
	}
	if l.phase != launcherSelectFaction {
		t.Fatalf("選完劇本應該進清單，得到 %v", l.phase)
	}

	l.cursor = 1
	l.apply(launcherConfirm)
	if l.phase != launcherSelectPlayer {
		t.Fatalf("在清單上決定應該進君主卡，得到 %v", l.phase)
	}
	if l.cursor != 1 {
		t.Errorf("進卡片之後選中的勢力變了：cursor = %d", l.cursor)
	}

	l.apply(launcherCancel)
	if l.phase != launcherSelectFaction {
		t.Fatalf("卡片取消應該回清單，得到 %v", l.phase)
	}
	if l.cursor != 1 {
		t.Errorf("退回清單之後選中的勢力變了：cursor = %d", l.cursor)
	}

	l.apply(launcherCancel)
	if l.phase != launcherScenario {
		t.Fatalf("清單取消應該回劇本，得到 %v", l.phase)
	}

	// 再走一次到底，確認確定會帶著正確的勢力編號。
	l.setScenarioPlayers(0, "第一章", players)
	l.cursor = 1
	l.apply(launcherConfirm)
	l.apply(launcherConfirm)
	if l.phase != launcherGameConfirm || l.confirmedPlayer != 4 {
		t.Errorf("最後 = %v／player %d，want launcherGameConfirm／4",
			l.phase, l.confirmedPlayer)
	}
}

// 背景的鏡頭就是 `sub_11A6E` 在 `sub_1D615` 前設的那兩個立即值
// （docs/spec/79 §1.1.1）。**不帶 marks**，所以據點沒有勢力徽記。
func TestLauncherCameraMatchesOriginal(t *testing.T) {
	if launcherCamX != 0xAA || launcherCamY != 0x62 {
		t.Errorf("鏡頭 = (%d,%d)，want (0xAA,0x62) ＝ (170,98)",
			launcherCamX, launcherCamY)
	}
	// 鏡頭要落在可視範圍內：世界 384×256、視野 40×23。
	if launcherCamX < 0 || launcherCamX > 384-viewCols {
		t.Errorf("camX %d 超出可視範圍", launcherCamX)
	}
	if launcherCamY < 0 || launcherCamY > 256-viewRows {
		t.Errorf("camY %d 超出可視範圍", launcherCamY)
	}
}

// 滾輪捲的是**視野**，不動選取列（原版的 top 與選取列是兩個獨立狀態）。
func TestFactionListWheelMovesViewNotCursor(t *testing.T) {
	l := &launcherModel{phase: launcherSelectFaction}
	for i := 0; i < 22; i++ {
		l.players = append(l.players, launcherPlayer{ID: i})
	}
	l.cursor = 0

	l.factionTop = clampFactionTop(l.factionListTop()+1, len(l.players))
	if l.factionListTop() != 1 {
		t.Fatalf("捲一格之後 top = %d，want 1", l.factionListTop())
	}
	if l.cursor != 0 {
		t.Errorf("捲動不該動到選取列，cursor = %d", l.cursor)
	}
	// 捲動之後第一列對到的是第 1 個勢力。
	if n, ok := l.factionListSelectAt(200, 121); !ok || n != 1 {
		t.Errorf("捲動後第一列 = (%d,%v)，want (1,true)", n, ok)
	}
	// ↑↓ 那一邊才把視野拉回來。
	l.move(0)
	if l.factionListTop() != 0 {
		t.Errorf("移動游標之後 top = %d，want 0（要把 cursor 拉回畫面）", l.factionListTop())
	}
}
