package main

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
)

// 四個一覽表家族的欄位與逐列格式。出處 docs/re/26 §4.1／§9、
// docs/re/27 §2–§4，規格 docs/spec/38。
//
// ⭐ **「看」與「選」用同一組欄位**：原版八個呼叫端只分五組
// `(bx, si, di)`，兩種取法只差在建清單的 callback（docs/re/26 §4.2）。
// 所以這裡一個家族只寫一次，picker 與一覽表共用。

// listNum 是數字欄的格式化：原版走 `sub_1062F(值, 位數)`，位數不足補空白。
func listNum(v, digits int) string { return fmt.Sprintf("%*d", digits, v) }

// generalRank 回傳武將的身分（docs/re/26 §9 的名稱表索引）。
//
// ⚠ 原版存在武將記錄 `+0x17`（0–5），remake 只留了 `Posted bool`，
// 所以這裡**從狀態反推**——君主那一格原版本來也是顯示時反查的。
// 俘虜推不出來（docs/spec/38 §4）。
func (g *game) generalRank(id int) int {
	w := g.world
	if w == nil || id < 0 || id >= len(w.Generals) {
		return 0
	}
	gen := w.Generals[id]
	if gen.Faction >= 0 && gen.Faction < len(w.Factions) {
		f := w.Factions[gen.Faction]
		if f.Lord == id {
			return 5 // 君主
		}
	}
	for i := range w.Factions {
		if w.Factions[i].Diplomat == id {
			return 3 // 外交官
		}
	}
	for i := range w.Cities {
		if w.Cities[i].Governor == id {
			return 2 // 內政官
		}
	}
	if gen.Posted {
		return 1 // 軍團長
	}
	return 0
}

// listColumnsGenerals 是武將家族的六欄（順序 ＝ 標題字串的欄序）。
func (g *game) listColumnsGenerals() []listwin.Column {
	gs := g.world.Generals
	return []listwin.Column{
		{Title: "武將名", Less: func(a, b int) bool { return gs[a].Name < gs[b].Name }},
		{Title: "武術", Less: func(a, b int) bool { return gs[a].Martial > gs[b].Martial }},
		{Title: "統率", Less: func(a, b int) bool { return gs[a].Command > gs[b].Command }},
		{Title: "政治", Less: func(a, b int) bool { return gs[a].Politics > gs[b].Politics }},
		{Title: "勢力", Less: func(a, b int) bool { return gs[a].Faction < gs[b].Faction }},
		{Title: "身分", Less: func(a, b int) bool { return g.generalRank(a) > g.generalRank(b) }},
	}
}

func (g *game) listRowGeneral(id int) []string {
	gen := g.world.Generals[id]
	faction := "－－－"
	if gen.Faction >= 0 && gen.Faction < len(g.world.Factions) {
		faction = big5(g.world.LordName(gen.Faction))
	}
	return []string{
		big5(gen.Name),
		listNum(gen.Martial, 2), listNum(gen.Command, 2), listNum(gen.Politics, 2),
		faction, listRankNames[g.generalRank(id)],
	}
}

// listColumnsCities 是據點家族的六欄。
func (g *game) listColumnsCities() []listwin.Column {
	cs := g.world.Cities
	return []listwin.Column{
		{Title: "據點名", Less: func(a, b int) bool { return cs[a].Name < cs[b].Name }},
		{Title: "生產力", Less: func(a, b int) bool { return cs[a].Production > cs[b].Production }},
		{Title: "上昇率", Less: func(a, b int) bool { return cs[a].Growth > cs[b].Growth }},
		{Title: "防災", Less: func(a, b int) bool { return cs[a].Prevention > cs[b].Prevention }},
		{Title: "城兵", Less: func(a, b int) bool { return cs[a].Garrison > cs[b].Garrison }},
		{Title: "內政官", Less: func(a, b int) bool { return cs[a].Governor < cs[b].Governor }},
	}
}

func (g *game) listRowCity(id int) []string {
	c := g.world.Cities[id]
	who := "　　　"
	if c.Governor >= 0 && c.Governor < len(g.world.Generals) {
		who = big5(g.world.Generals[c.Governor].Name)
	}
	// 城兵**先 ×10**（docs/re/27 §3）；上昇率的存值是實際值 ＋100，
	// remake 的 Growth 已經是實際值，所以直接印。
	return []string{
		big5(c.Name),
		listNum(c.Production, 5), listNum(c.Growth, 4),
		listNum(c.Prevention, 3), listNum(c.Garrison*10, 4), who,
	}
}

// listColumnsFactions 是勢力家族的六欄。
func (g *game) listColumnsFactions() []listwin.Column {
	fs := &g.world.Factions
	return []listwin.Column{
		{Title: "勢力名", Less: func(a, b int) bool {
			return g.world.LordName(a) < g.world.LordName(b)
		}},
		{Title: "武將", Less: func(a, b int) bool { return fs[a].Generals > fs[b].Generals }},
		{Title: "據點", Less: func(a, b int) bool { return fs[a].Cities > fs[b].Cities }},
		{Title: "首都", Less: func(a, b int) bool { return fs[a].Capital < fs[b].Capital }},
		{Title: "外交", Less: func(a, b int) bool {
			return g.factionDiplomacy(a) > g.factionDiplomacy(b)
		}},
		{Title: "外交官", Less: func(a, b int) bool { return fs[a].Diplomat < fs[b].Diplomat }},
	}
}

// factionDiplomacy 回傳玩家對某個勢力的關係等級（listDiplomacyNames 的索引）。
func (g *game) factionDiplomacy(id int) int {
	w := g.world
	if w == nil || id == w.Player {
		return 6 // －－：自己那一格
	}
	f := w.Friendship[w.Player][id]
	return listDiplomacyLevel(f.Value(), f.AtWar())
}

func (g *game) listRowFaction(id int) []string {
	f := g.world.Factions[id]
	capital := "－－－"
	if f.Capital >= 0 && f.Capital < len(g.world.Cities) {
		capital = big5(g.world.Cities[f.Capital].Name)
	}
	who := "　　　"
	if f.Diplomat >= 0 && f.Diplomat < len(g.world.Generals) {
		who = big5(g.world.Generals[f.Diplomat].Name)
	}
	return []string{
		big5(g.world.LordName(id)),
		listNum(f.Generals, 3), listNum(f.Cities, 3), capital,
		listDiplomacyNames[g.factionDiplomacy(id)], who,
	}
}

// listColumnsCorps 是軍團家族的五欄。
func (g *game) listColumnsCorps() []listwin.Column {
	cp := g.world.Corps
	gs := g.world.Generals
	return []listwin.Column{
		{Title: "武將名", Less: func(a, b int) bool {
			return gs[g.world.Leader(a)].Name < gs[g.world.Leader(b)].Name
		}},
		{Title: "總兵數", Less: func(a, b int) bool { return cp[a].Men > cp[b].Men }},
		{Title: "士氣值", Less: func(a, b int) bool { return cp[a].Morale > cp[b].Morale }},
		{Title: "現在位置", Less: func(a, b int) bool { return cp[a].Node < cp[b].Node }},
		{Title: "目標據點", Less: func(a, b int) bool { return cp[a].Ordered < cp[b].Ordered }},
	}
}

func (g *game) listRowCorps(id int) []string {
	c := g.world.Corps[id]
	name := func(node int) string {
		if node < 0 || node >= len(g.world.Cities) {
			return "－　　" // 不在據點上（原版 +0x0E ≥ 0x800）
		}
		return big5(g.world.Cities[node].Name)
	}
	// 總兵數**先 ×10**（docs/re/27 §3）。
	return []string{
		big5(g.world.Generals[g.world.Leader(id)].Name),
		listNum(c.Men*10, 4), listNum(c.Morale, 3),
		name(c.Node), name(c.Ordered),
	}
}
