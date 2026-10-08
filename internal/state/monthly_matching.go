//go:build matching

package state

import "github.com/wicanr2/wolong_cht/internal/rules/economy"

// MatchingMonthly 只供原版/C/Go 研究比較；正式建置不含此直接入口。
func (w *World) MatchingMonthly(rng economy.Rand) Event {
	w.rng = rng
	ev := Event{HourFaction: -1, Settled: true}
	w.monthlyRules(&ev, rng)
	return ev
}
