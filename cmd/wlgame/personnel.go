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

)

// NoOfficial 是「沒有派駐」的哨兵值（原版的 0xFF）。
const NoOfficial = 0xFF

// openPersonnel 開人事的第一層：選要做哪一件事。
//
// 用一覽表當選單而不是另做一個視窗——說明書 3.8 的兩段式選取
// 對所有清單都成立，重用同一個狀態機比較不會走樣。
// openPersonnel 是指令列第 2 格（「人事」）。
//
// ⚠ 原版是**四列彈出選單**（`sub_193E9(ax=4, cx=4Eh, dx=403h)`，
// docs/spec/126），不是一覽表。四個項目與順序照 `funcs_16279`。
func (g *game) openPersonnel() { g.openPopupMenu(personnelPopupMenu) }

// dispatchPersonnelMenu 是那四列各自接到哪。
func (g *game) dispatchPersonnelMenu(row int) {
	switch row {
	case 0:
		g.pickCityForGovernor()
	case 1:
		g.removeGovernor()
	case 2:
		g.pickFactionForDiplomat()
	case 3:
		g.removeDiplomat()
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



// 人事的四條出口各自先掛一則狀態列提示（`sub_16A9B`／`sub_16B08`／
// `sub_16B71`／`sub_16BE3` 開頭的 `sub_18853`，docs/re/25 §3、docs/spec/140）。
const (
	governorAssignTalk = 0x0B // #11「要派遣內政官到哪個據點？」
	diplomatAssignTalk = 0x0C // #12「要派遣外交官到哪個勢力？」
	governorRemoveTalk = 0x0D // #13「要解任哪個據點的內政官？」
	diplomatRemoveTalk = 0x0E // #14「要解任哪個勢力的外交官？」

	// 解任時選到「那裡本來就沒人」的兩則（`sub_16B08` 的 `cx = 36h`／
	// `sub_16BE3` 的 `cx = 37h`，docs/spec/142）。
	nobodyPostedTalk        = 0x36 // #54「是否有所差錯？{2}並未派遣任何人。」
	factionNobodyPostedTalk = 0x37 // #55「{3}勢力仍未派遣任何人．．．」
)

func (g *game) pickCityForGovernor() {
	g.setStatusTalk(governorAssignTalk, nil)
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

// removeGovernor 是「內政官解任」。
//
// ⭐ **不先過濾**（docs/spec/142）：原版 `sub_16B08` 開的是全部據點的清單，
// 選到沒派人的那一個才跳 TALK #54，**沒有「沒有派駐中的內政官」這種前置檢查**。
// 先過濾會讓玩家看不到「那座城本來就沒人」這件事。
func (g *game) removeGovernor() {
	g.setStatusTalk(governorRemoveTalk, nil)
	rows := g.playerCities()
	if len(rows) == 0 {
		g.lastEvent = "沒有據點"
		g.list = nil
		return
	}
	g.cityList(rows, "選要解任的據點　Enter 決定　ESC 取消", func(city int) bool {
		c := &g.world.Cities[city]
		// 原版 `sub_16B4F`：先寫 0xFF、再看舊值是不是 0xFF。
		old := c.Governor
		c.Governor = NoOfficial
		if old < 0 || old >= len(g.world.Generals) || old == NoOfficial {
			g.enqueueTalk(nobodyPostedTalk,
				map[byte]string{'2': padTalkField(big5(c.Name))})
			return true
		}
		g.lastEvent = fmt.Sprintf("%s 解任內政官 %s",
			big5(c.Name), big5(g.world.Generals[old].Name))
		return true
	})
}

func (g *game) pickFactionForDiplomat() {
	g.setStatusTalk(diplomatAssignTalk, nil)
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

// removeDiplomat 是「外交官解任」。與 removeGovernor 同形（docs/spec/142）：
// **不先過濾**，選到沒派人的才跳 TALK #55。
func (g *game) removeDiplomat() {
	g.setStatusTalk(diplomatRemoveTalk, nil)
	var rows []int
	for i := range g.world.Factions {
		f := &g.world.Factions[i]
		if f.Alive && i != g.world.Player {
			rows = append(rows, i)
		}
	}
	if len(rows) == 0 {
		g.lastEvent = "沒有可解任的勢力"
		g.list = nil
		return
	}
	g.factionList(rows, "選要召回外交官的勢力　Enter 決定　ESC 取消", func(f int) bool {
		fa := &g.world.Factions[f]
		// 原版 `sub_16C2A`：先寫 0xFF、再看舊值。
		old := fa.Diplomat
		fa.Diplomat = NoOfficial
		if old < 0 || old >= len(g.world.Generals) || old == NoOfficial {
			g.enqueueTalk(factionNobodyPostedTalk,
				map[byte]string{'3': padTalkField(big5(g.world.LordName(f)))})
			return true
		}
		g.lastEvent = fmt.Sprintf("召回派駐 %s 軍的 %s",
			big5(g.world.LordName(f)), big5(g.world.Generals[old].Name))
		return true
	})
}

// factionList 開一張勢力清單。交友度是外交官要不要派的主要依據。
func (g *game) factionList(rows []int, hint string, pick func(int) bool) {
	g.openFactionPicker(rows, hint, pick)
}


