package tactical

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

func TestCollisionGridAgainstOriginalOpeningReceipt(t *testing.T) {
	const root = "../../../workplace/six-fixes-20260908/original-collision/"
	data, err := os.ReadFile(root + "tick-0001.json")
	if os.IsNotExist(err) {
		t.Skip("未安裝本機 dosgolem 收據；不宣稱原版對拍通過")
	}
	if err != nil {
		t.Fatal(err)
	}
	var snapshot struct {
		Units string `json:"unit_bytes"`
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	raw, err := hex.DecodeString(snapshot.Units)
	if err != nil || len(raw) < 96*32 {
		t.Fatal("原版單位資料長度或編碼錯誤")
	}
	want, err := os.ReadFile(root + "collision-tick2-0000.bin")
	if err != nil || len(want) < 8192 {
		t.Fatal("原版碰撞表缺漏")
	}
	var g collisionGrid
	for i := 0; i < 96; i++ {
		u := raw[i*32 : (i+1)*32]
		g.drawSlot(i, int(u[6]), int(u[8]), int(u[10]))
	}
	for i := 0; i < 8192; i++ {
		if g.cells[i] != want[i]&127 {
			t.Fatalf("ES:%04x low7=%d，原版=%d", i, g.cells[i], want[i]&127)
		}
	}
}

func TestCollisionGridOverlappingDrawAndClear(t *testing.T) {
	var g collisionGrid
	g.drawSlot(3, 1, 38, 0)
	g.drawSlot(8, 1, 38, 0)
	if got := g.at(1, 38, 0); got != 8 {
		t.Fatalf("重疊重畫取槽 %d，want 8", got)
	}
	g.drawSlot(3, 2, 38, 0)
	if got := g.at(1, 38, 0); got != -1 {
		t.Fatalf("原格應清空，got %d", got)
	}
	if got := g.at(2, 38, 0); got != 3 {
		t.Fatalf("新格 = %d，want 3", got)
	}
}

func TestCollisionGridSwapPreservesOriginalUpperLayerAsymmetry(t *testing.T) {
	var g collisionGrid
	g.drawSlot(8, 1, 38, 0)
	g.drawSlot(20, 1, 37, 0)
	a, _ := collisionCell(1, 38, 0)
	b, _ := collisionCell(1, 37, 0)
	g.swapSlots(8, 20, a, b)
	if g.cells[a] != 9 || g.cells[b] != 21 || g.cells[a+0x1000] != 9 || g.cells[b+0x1000] != 9 {
		t.Fatalf("換格不符原始寫入：%d %d / %d %d", g.cells[a], g.cells[b], g.cells[a+0x1000], g.cells[b+0x1000])
	}
	if g.old[8] != b || g.old[20] != a {
		t.Fatal("未交換原格指標")
	}
}
