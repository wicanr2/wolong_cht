package main

// 人事：內政官與外交官的任命／解任（說明書 §3.3「人事」，
// 整理在 docs/mechanics/10-strategy.md）。
//
// **外交官派駐到「勢力」，內政官派駐到「據點」**——兩者的第一階段
// 選的東西不同，這是說明書特別點出來的差異。
//
// 這一格先前在命令視窗上是灰的。內政官的效果解出來之後
// （`sub_14194`，docs/re/07 §19）它才有意義：
// **不派內政官的玩家據點基準是 5，比 AI 的 8 還差**，
// 派一個政治 15 的內政官會拉到 20。這不是錦上添花，是必要的補救。

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
)

// NoOfficial 是「沒有派駐」的哨兵值（原版的 0xFF）。
const NoOfficial = 0xFF

// openPersonnel 開人事的第一層：選要做哪一件事。
//
// 用一覽表當選單而不是另做一個視窗——說明書 3.8 的兩段式選取
// 對所有清單都成立，重用同一個狀態機比較不會走樣。
func (g *game) openPersonnel() {
	const (
		assignGov = iota
		removeGov
		assignDip
		removeDip
	)
	names := []string{"內政官任命", "內政官解任", "外交官任命", "外交官解任"}
	rows := []int{assignGov, removeGov, assignDip, removeDip}

	// 人事是**選單**不是一覽表（原版 3.2 的子選單），所以不套家族欄位。
	g.list = listwin.New(listwin.Generals, []listwin.Column{
		{Title: "人　事", Less: func(a, b int) bool { return a < b }},
	}, rows, listRowsPerPage, &g.sortMem)
	g.listTitle = "人　事"
	g.listCellInk = nil
	g.listRow = func(i int) []string { return []string{names[i]} }
	g.listHint = "↑↓ 移動　Enter 選取／決定　ESC 取消"
	g.listPick = func(i int) bool {
		switch i {
		case assignGov:
			g.pickCityForGovernor()
		case removeGov:
			g.removeGovernor()
		case assignDip:
			g.pickFactionForDiplomat()
		case removeDip:
			g.removeDiplomat()
		}
		return false // 下一層自己換掉 g.list
	}
}

// playerCities 是玩家目前擁有的據點編號。
func (g *game) playerCities() []int {
	var out []int
	for i := range g.world.Cities {
		if g.world.Cities[i].Owner == g.world.Player {
			out = append(out, i)
		}
	}
	return out
}

// freeGenerals 是可以派任的武將：活著、屬於玩家、沒出陣、不是俘虜。
//
// ⚠ 沒有排除「已經在別處當官的」——原版沒有這個限制，
// 而任命本身就會把舊的那一格覆蓋掉。
func (g *game) freeGenerals() []int {
	var out []int
	for i := range g.world.Generals {
		gen := &g.world.Generals[i]
		if gen.Alive && gen.Faction == g.world.Player && !gen.Posted {
			out = append(out, i)
		}
	}
	return out
}

// cityList 開一張據點清單，選完呼叫 pick。欄位照原版家族（docs/spec/38）。
func (g *game) cityList(rows []int, hint string, pick func(int) bool) {
	g.openCityPicker(rows, hint, pick)
}

// generalList 開一張武將清單，選完呼叫 pick。欄位照原版家族。
func (g *game) generalList(rows []int, hint string, pick func(int) bool) {
	g.openGeneralPicker(rows, hint, pick)
}



func (g *game) pickCityForGovernor() {
	rows := g.playerCities()
	if len(rows) == 0 {
		g.lastEvent = "沒有據點"
		g.list = nil
		return
	}
	// **先選據點，再選武將**（說明書：內政官任命的順序）。
	g.cityList(rows, "選要派內政官的據點　Enter 決定　ESC 取消", func(city int) bool {
		free := g.freeGenerals()
		if len(free) == 0 {
			g.lastEvent = "沒有可派任的武將"
			return true
		}
		name := big5(g.world.Cities[city].Name)
		g.generalList(free, "選要派去 "+name+" 的武將　Enter 決定", func(who int) bool {
			g.world.Cities[city].Governor = who
			g.lastEvent = fmt.Sprintf("%s 派任 %s 為內政官",
				name, big5(g.world.Generals[who].Name))
			return true
		})
		return false
	})
}

func (g *game) removeGovernor() {
	var rows []int
	for _, i := range g.playerCities() {
		if gv := g.world.Cities[i].Governor; gv >= 0 && gv < len(g.world.Generals) {
			rows = append(rows, i)
		}
	}
	if len(rows) == 0 {
		g.lastEvent = "沒有派駐中的內政官"
		g.list = nil
		return
	}
	g.cityList(rows, "選要解任的據點　Enter 決定　ESC 取消", func(city int) bool {
		c := &g.world.Cities[city]
		g.lastEvent = fmt.Sprintf("%s 解任內政官 %s",
			big5(c.Name), big5(g.world.Generals[c.Governor].Name))
		c.Governor = NoOfficial
		return true
	})
}

func (g *game) pickFactionForDiplomat() {
	// 外交官派駐到**別的勢力**，所以候選是「活著且不是自己」。
	var rows []int
	for i := range g.world.Factions {
		if g.world.Factions[i].Alive && i != g.world.Player {
			rows = append(rows, i)
		}
	}
	if len(rows) == 0 {
		g.lastEvent = "沒有可派駐的勢力"
		g.list = nil
		return
	}
	g.factionList(rows, "選要派外交官的勢力　Enter 決定　ESC 取消", func(f int) bool {
		free := g.freeGenerals()
		if len(free) == 0 {
			g.lastEvent = "沒有可派任的武將"
			return true
		}
		name := big5(g.world.LordName(f))
		g.generalList(free, "選要派去 "+name+" 的武將　Enter 決定", func(who int) bool {
			g.world.Factions[f].Diplomat = who
			g.lastEvent = fmt.Sprintf("派 %s 出使 %s 軍",
				big5(g.world.Generals[who].Name), name)
			return true
		})
		return false
	})
}

func (g *game) removeDiplomat() {
	var rows []int
	for i := range g.world.Factions {
		f := &g.world.Factions[i]
		if f.Alive && i != g.world.Player &&
			f.Diplomat >= 0 && f.Diplomat < len(g.world.Generals) {
			rows = append(rows, i)
		}
	}
	if len(rows) == 0 {
		g.lastEvent = "沒有派駐中的外交官"
		g.list = nil
		return
	}
	g.factionList(rows, "選要召回外交官的勢力　Enter 決定　ESC 取消", func(f int) bool {
		fa := &g.world.Factions[f]
		g.lastEvent = fmt.Sprintf("召回派駐 %s 軍的 %s",
			big5(g.world.LordName(f)), big5(g.world.Generals[fa.Diplomat].Name))
		fa.Diplomat = NoOfficial
		return true
	})
}

// factionList 開一張勢力清單。交友度是外交官要不要派的主要依據。
func (g *game) factionList(rows []int, hint string, pick func(int) bool) {
	g.openFactionPicker(rows, hint, pick)
}


