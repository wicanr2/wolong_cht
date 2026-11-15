package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// 原版 sub_13E11：資金高word相等時取消侵攻；bit6另比較半門檻。
// 經正式 Tick 跨時更新，再存讀檔，避免只驗孤立公式。
func TestHourlyInvasionFinanceAndSave(t *testing.T) {
	for _, tc := range []struct {
		funds       int
		cancel, low bool
	}{
		{45182, true, false}, {45311, true, false}, {45312, false, false},
		{22656, true, true}, {22783, true, true}, {22784, true, false},
	} {
		w := load(t, 0)
		const faction = 11
		f := &w.Factions[faction]
		f.Alive, f.Cities, f.Funds, f.InvasionTarget, f.LowFunds = true, 19, tc.funds, 14, true
		w.hourFaction = faction
		w.Clock.Hour, w.Clock.Subtick = 4, 8
		ev := w.Tick(rng.NewFixed(7))
		if !ev.Clock.Hour || ev.HourFaction != faction || ev.InvasionCancelled != tc.cancel {
			t.Fatalf("資金%d 每時事件=%+v", tc.funds, ev)
		}
		wantTarget := 14
		if tc.cancel {
			wantTarget = 255
		}
		if f.InvasionTarget != wantTarget || f.LowFunds != tc.low {
			t.Fatalf("資金%d target=%d low=%v", tc.funds, f.InvasionTarget, f.LowFunds)
		}
		loaded, err := LoadBlock(w.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		got := loaded.Factions[faction]
		if got.InvasionTarget != wantTarget || got.LowFunds != tc.low {
			t.Fatalf("資金%d 存讀檔遺失財政判定：%+v", tc.funds, got)
		}
	}
}
