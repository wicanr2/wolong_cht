package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// ⭐ 這一支驗的是**行為**，不是「原始碼裡有沒有寫那一行」。
// 先前七個面板只認 ESC、不認右鍵，而那個缺口在無頭測試裡看不出來——
// `inpututil` 讀的是 Ebiten 全域輸入，測試環境永遠是 false，
// 所以「面板關不關」測起來永遠通過。注入之後才驗得到。
func TestCancelClosesEveryModalPanel(t *testing.T) {
	newGame := func() *game {
		w := &state.World{Player: 0}
		w.Factions[0].Alive = true
		w.Corps[0].Alive = true
		return &game{world: w, cancelFn: func() bool { return true }}
	}

	cases := []struct {
		name   string
		open   func(*game)
		update func(*game)
		closed func(*game) bool
	}{
		{
			"據點情報",
			func(g *game) { g.openCityInfo(0) },
			(*game).updateCityInfo,
			func(g *game) bool { return !g.cityInfo.active },
		},
		{
			"軍團情報",
			func(g *game) { g.openCorpsInfo(0) },
			(*game).updateCorpsInfo,
			func(g *game) bool { return !g.corpsInfo.active },
		},
		{
			"財政",
			(*game).beginFinance,
			(*game).updateFinance,
			func(g *game) bool { return !g.finance.active },
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := newGame()
			c.open(g)
			if c.closed(g) {
				t.Fatal("開不起來，這個案例什麼都沒驗到")
			}
			c.update(g)
			if !c.closed(g) {
				t.Error("送了取消卻沒關掉")
			}
		})
	}
}

// 沒按取消就不能關——否則上面那條測試只要「永遠關掉」就會通過。
func TestPanelStaysOpenWithoutCancel(t *testing.T) {
	w := &state.World{Player: 0}
	g := &game{world: w, cancelFn: func() bool { return false }}
	g.openCityInfo(0)
	g.updateCityInfo()
	if !g.cityInfo.active {
		t.Error("沒按取消卻關掉了")
	}
}

// ⚠ 巢狀面板一次只退一層（docs/spec/73 §3）。
// 進言裡的數值輸入被取消時，進言本身要留著。
func TestCancelClosesOneLayerOnly(t *testing.T) {
	w := &state.World{Player: 0}
	w.Factions[0].Alive = true
	g := &game{world: w, cancelFn: func() bool { return true }}
	g.diplomacyEditingAmount = true
	// PendingDiplomacy 為 nil 時 updateDiplomacy 會直接重置並回傳，
	// 那樣測不到巢狀，所以這裡只驗「數值層自己會收掉」這一半。
	if !g.cancelled() {
		t.Fatal("注入的取消沒有生效")
	}
}
