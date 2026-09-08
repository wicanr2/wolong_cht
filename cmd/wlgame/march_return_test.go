package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 驗證入口的續行，不以這個窄狀態測試代替正常下令／存檔重播。
func TestMarchOrderReturnsToCallingView(t *testing.T) {
	for _, fromMap := range []bool{true, false} {
		w := &state.World{Player: 0}
		w.Corps[39].Alive = true
		w.Corps[39].Faction = 0
		w.Corps[39].TargetNode = 1
		w.Corps[39].Ordered = 1
		w.Corps[40].Alive = true // 另一格的友軍不應混進地圖入口清單。
		g := &game{world: w}
		if !fromMap {
			g.beginMarch()
		} else {
			g.openCorpsOnTile([]int{39})
		}
		g.listPick(39)
		g.endMapPick() // 地圖命中完成；路線規則由另一層測試與 GUI 重播覆蓋。
		g.beginMarchMode(39, 1)
		g.commitMarchMode()
		if g.marchMode.active || g.mapPickActive() {
			t.Fatal("完成後仍留在行軍子視窗")
		}
		if g.list == nil {
			t.Fatal("下令後沒有回到父軍團清單")
		}
		wantRows := 2
		if fromMap {
			wantRows = 1
		}
		if len(g.list.Rows) != wantRows {
			t.Fatalf("返回清單的範圍錯誤：地圖入口=%v，列數=%d", fromMap, len(g.list.Rows))
		}
		g.dispatchListAction(listUIAction{kind: listActionCancel})
		if g.list != nil || g.corpsInfo.active || g.marchReturn != nil {
			t.Fatalf("返回地圖後仍有阻擋輸入的視窗：地圖入口=%v", fromMap)
		}
	}
}

func TestMarchTargetCancelReturnsToCallingView(t *testing.T) {
	w := &state.World{Player: 0}
	w.Corps[39].Alive = true
	g := &game{world: w}
	g.beginMarch()
	g.listPick(39)
	cancel := g.mapPick.cancel
	g.endMapPick()
	cancel()
	if g.list == nil || g.marchMode.active {
		t.Fatal("取消目的地沒有返回原軍團清單")
	}
}
