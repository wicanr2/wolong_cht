package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/capital"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
)

// 進言的第四項與第五項（說明書 3.2）。
//
// 這兩項與敵對／停戰／協力不同：**沒有說服迴圈**。君主只做一次驗收，
// 不問理由、也不動信賴度——原版兩支都直接走 `sub_13B08`
// （插圖 ＋ 上框 ＋ 下框 ＋ 上框，兩種結果），沒有 `sub_13B5A`。
//
// 規格 `docs/spec/49`，機制 `docs/mechanics/70-ai.md`。

// AdviseRelocateAccepted 回報君主會不會答應遷都到 node（原版 `sub_16909`）。
//
// **只判斷，不動狀態**——畫面要先把君主的回答演完才真的搬。
func (w *World) AdviseRelocateAccepted(node int) bool {
	if w == nil || w.Player < 0 || w.Player >= numFactions {
		return false
	}
	from := w.clampCity(w.Factions[w.Player].Capital)
	if node < 0 || node >= len(w.Cities) || node == from {
		return false
	}
	if w.Cities[node].Owner != w.Player {
		return false // 原版由選單只列自己的據點擋掉
	}
	site := func(i int) capital.Site {
		c := &w.Cities[i]
		return capital.Site{Owner: c.Owner, Kind: c.Kind,
			Production: c.Production, Adjacency: c.Adjacency}
	}
	return capital.AcceptRelocation(site(from), site(node))
}

// AdviseRelocate 把首都搬到 node，回傳有沒有搬成。
//
// 搬成之後要跑 `sub_14502`（`syncCorpsAfterCapitalChange`）：
// 目標還掛著舊首都的軍團一律改掛新首都。
func (w *World) AdviseRelocate(node int) bool {
	if !w.AdviseRelocateAccepted(node) {
		return false
	}
	old := w.Factions[w.Player].Capital
	w.Factions[w.Player].Capital = node
	w.syncCorpsAfterCapitalChange(w.Player, w.clampCity(old), node)
	return true
}

// sortie 把玩家勢力目前的狀況換成出陣判定要的三個數。
func (w *World) sortie() persuasion.Sortie {
	f := &w.Factions[w.Player]
	reserves := 0
	for _, n := range f.Reserves {
		reserves += n
	}
	return persuasion.Sortie{
		Funds: f.Funds, Aggression: f.Aggression, Reserves: reserves,
	}
}

// LordLeadsCorps 回報玩家的君主現在是不是帶著軍團在外面。
//
// ⭐ **這是「君主在不在朝堂上」的唯一判準。** 原版只有出陣那一條路會問它
// （`sub_16EC9`「君主已經帶著軍團就不能再出一次」）；remake 允許把君主
// 編成軍團長（docs/spec/76），於是同一個狀態也擋住整個進言
// （docs/spec/111）。兩個呼叫端共用這一支，避免各寫一份而行為分岔。
//
// 君主編號同時是軍團編號——原版 `sub_1291A` 直接把軍團記錄位址換算成
// 武將記錄位址，兩張表是同一組索引（docs/formats/08 §2）。
func (w *World) LordLeadsCorps() bool {
	if w == nil || w.Player < 0 || w.Player >= numFactions {
		return false
	}
	lord := w.Factions[w.Player].Lord
	if lord < 0 || lord >= numCorps {
		return false
	}
	return w.Corps[lord].Alive
}

// AdviseSortieAccepted 回報君主會不會親自出陣（原版 `sub_1699E` 的兩道閘）。
func (w *World) AdviseSortieAccepted() bool {
	if w == nil || w.Player < 0 || w.Player >= numFactions {
		return false
	}
	if lord := w.Factions[w.Player].Lord; lord < 0 || lord >= numCorps {
		return false
	}
	if w.LordLeadsCorps() {
		return false // 君主已經帶著軍團了（原版由 `sub_16EC9` 擋）
	}
	return persuasion.AcceptSortie(w.sortie())
}

// AdviseSortie 讓君主本人編一支軍團，回傳有沒有編成。
//
// ⭐ **這一支不是委任的。** `sub_16E8F` 一律把委任位元設起來，
// 而 `sub_1699E` 緊接著 `and byte ptr [di], 0FBh` 把它清掉——
// 君主親自出陣，指揮權在玩家手上。
func (w *World) AdviseSortie() bool {
	if !w.AdviseSortieAccepted() {
		return false
	}
	return w.autoFormCorps(w.Player, w.Factions[w.Player].Lord, false)
}
