package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
)

// PersuasionSituation 把目前的世界狀態換成說服判定要的局勢。
//
// 每一欄都對應已解出的資料：國力用據點數、疲弊用資金 < 0、
// 侵攻狀態用勢力記錄 +0x19（docs/re/08 §1）。
//
// ally 只有「請求協助」用得到，其餘傳 -1。
//
// ⚠ 這一份是**桌面版與手機版共用**的。局勢的組法有好幾個容易錯的細節
//（有向的交友度、含和平位元的原始值、掃侵攻時要跳過交涉對象本身），
// 各寫一份必然會長出差異，而差異只會在特定局面出現（CLAUDE.md §7 第 6 條）。
func (w *World) PersuasionSituation(cmd persuasion.Command, target, ally int) persuasion.Situation {
	if target < 0 || target >= len(w.Factions) ||
		w.Player < 0 || w.Player >= len(w.Factions) {
		return persuasion.Situation{}
	}
	me := w.Factions[w.Player]
	them := w.Factions[target]

	// 「有別的勢力正在侵攻我方」——停戰提案的「我正在防禦戰」用這個，
	// 原版是 `sub_16577` 裡掃 22 個勢力的迴圈，**而且跳過交涉對象本身**。
	anyone := false
	for i := range w.Factions {
		f := &w.Factions[i]
		if f.Alive && i != target && f.InvasionTarget == w.Player {
			anyone = true
			break
		}
	}

	s := persuasion.Situation{
		Trust:       w.Trust,
		Aggression:  me.Aggression,
		OurCities:   me.Cities,
		TheirCities: them.Cities,
		OurFunds:    me.Funds,
		TheirFunds:  them.Funds,

		// 交友度是**有向的**，而且判定式要的是**含和平位元的原始值**
		// （`0x80` ＝ 和平）——門檻常數本身就帶著它。
		Friendship: w.Friendship[w.Player][target].Raw(),

		TheyInvadeThirdParty: them.InvasionTarget != diplomacy.NoTarget &&
			them.InvasionTarget != w.Player,
		TheyInvadeUs:    them.InvasionTarget == w.Player,
		AnyoneInvadesUs: anyone,
	}
	if cmd == persuasion.Cooperate && ally >= 0 && ally < len(w.Factions) {
		s.AllyCities = w.Factions[ally].Cities
		s.AllyFriendship = w.Friendship[w.Player][ally].Raw()
		s.SameFactionPicked = ally == target
	}
	return s
}
