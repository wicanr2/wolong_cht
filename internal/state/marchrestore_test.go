package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
)

// roadsFor 建這個世界的道路圖並注入——**正式路徑的最小版本**。
//
// ⭐ 規則層不讀檔案，道路圖是呼叫端注入的；少了這一步 `w.step` 走直線退路，
// 而畫面與存檔欄位看起來都正常（docs/spec/169 §7）。
func roadsFor(t *testing.T, w *World) {
	t.Helper()
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skip(err)
	}
	xy := make([][2]int, len(w.Cities))
	for i := range w.Cities {
		xy[i] = [2]int{w.Cities[i].X, w.Cities[i].Y}
	}
	edges, err := world.RoadEdges(lib.World, xy)
	if err != nil {
		t.Fatal(err)
	}
	w.SetRoads(march.New(len(w.Cities), world.MarchEdges(edges, xy)))
}

func scenarioWithRoads(t *testing.T) *World {
	t.Helper()
	w, err := LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skip(err)
	}
	roadsFor(t, w)
	return w
}

// formOneCorps 照正常玩家路徑編一支軍團——**劇本開局沒有現成的軍團**。
func formOneCorps(t *testing.T, w *World) int {
	t.Helper()
	w.Player = 0
	leader := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction == w.Player && !g.Posted() {
			leader = i
			break
		}
	}
	if leader < 0 {
		t.Skip("找不到可編成的武將")
	}
	var kinds [army.Positions]army.TroopType
	var manned [army.Positions]bool
	kinds[0], manned[0] = army.Infantry, true
	if err := w.FormCorps(leader, kinds, manned); err != nil {
		t.Skip(err)
	}
	return leader
}

// TestMarchRoutesRestoreOnLoad 釘住 docs/spec/172 §4.5：存檔裡「正在行軍」的
// 軍團，載入之後要接回**格子路徑**。
//
// 判準是「剩下要走的每一格都與存檔前相同」，不是「有沒有抵達」——
// 走直線一樣會抵達，用終點當斷言等於什麼都沒驗。
func TestMarchRoutesRestoreOnLoad(t *testing.T) {
	w := scenarioWithRoads(t)

	corps := formOneCorps(t, w)
	c := &w.Corps[corps]
	target := -1
	for i := range w.Cities {
		if i != c.Node && w.roads.Route(c.Node, i) != nil {
			target = i
			break
		}
	}
	if target < 0 {
		t.Skip("找不到走得到的據點")
	}
	if err := w.March(corps, target); err != nil {
		t.Fatal(err)
	}
	// 走幾格，讓它離開據點、停在 leg 中間。
	for k := 0; k < 3; k++ {
		w.step(corps)
	}
	want := append([][2]int(nil), w.routes[corps]...)
	if len(want) == 0 {
		t.Fatal("走三格之後就沒有路徑了，換一個案例")
	}
	if c.LinkAddr == 0 {
		t.Fatal("行軍中的 +0x0E 應該是連結記錄位址，得到 0")
	}

	// 重新載入一個世界，把存檔會保存的那幾個欄位搬過去，再注入道路圖。
	w2, err := LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skip(err)
	}
	w2.Corps[corps].Alive = true
	w2.Corps[corps].Faction = c.Faction
	d := &w2.Corps[corps]
	d.LinkAddr, d.PathPtr, d.Direction = c.LinkAddr, c.PathPtr, c.Direction
	d.X, d.Y = c.X, c.Y
	d.TargetNode = c.TargetNode
	d.TargetX, d.TargetY = c.TargetX, c.TargetY
	d.Node = c.TargetNode // 載入端的佔位值（docs/spec/172 §4.5）
	roadsFor(t, w2)

	if got := w2.UnresolvedMarches(); got != 0 {
		t.Fatalf("有 %d 支軍團的行軍狀態還原不了", got)
	}
	if d.Node != c.Node {
		t.Fatalf("出發據點還原成 %d，要 %d", d.Node, c.Node)
	}
	got := w2.routes[corps]
	if len(got) != len(want) {
		t.Fatalf("剩餘格數 = %d，要 %d", len(got), len(want))
	}
	for k := range want {
		if got[k] != want[k] {
			t.Fatalf("第 %d 格 = %v，要 %v", k, got[k], want[k])
		}
	}
}

// TestUnresolvedMarchDoesNotMove 是上一支的負對照：**沒有道路圖時軍團不准動**。
//
// ⚠ 走直線比不動更難發現——不動看得出來（`UnresolvedMarches()` 會報），
// 走錯路只有逐格對拍才看得見。而且載入端的 `Node` 是佔位值，
// 這時候跑一步會被判成抵達目標，連帶在錯的地方觸發攻城。
func TestUnresolvedMarchDoesNotMove(t *testing.T) {
	w, err := LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skip(err)
	}
	corps := formOneCorps(t, w)
	c := &w.Corps[corps]
	// 假裝這支軍團在存檔裡正走在某條 leg 上，而道路圖沒有注入。
	c.LinkAddr, c.PathPtr, c.Direction = 0x0800, 0x2000, 4
	c.Node = c.TargetNode
	c.Timer = 1
	before := [2]int{c.X, c.Y}

	w.tickOneCorps(corps, 1, fixedRand{})

	if now := [2]int{c.X, c.Y}; now != before {
		t.Fatalf("沒有道路圖時軍團動了：%v → %v", before, now)
	}
}
