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
	"github.com/wicanr2/wolong_cht/internal/state"
)

// NoOfficial 是「沒有派駐」的哨兵值（原版的 0xFF）。規則層是同一個。
const NoOfficial = state.NoOfficial

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
// ⭐ **人事四條都不寫事件列**（docs/spec/145 §1.1 的判準）：
// 那位官員自己會說一句（`officialSays`），事件列再登記一次就是
// 同一件事講兩遍，而且那條列蓋在地圖上，逐像素對拍時是最大的差異
// （docs/playtest/92：5,808 px）。留著的只有「沒有據點」這種
// **原版走不到的狀態**——那是 remake 自己的防呆，不是重複的回報。
func (g *game) playerCities() []int {
	var out []int
	for i := range g.world.Cities {
		if g.world.Cities[i].Owner == g.world.Player {
			out = append(out, i)
		}
	}
	return out
}

// freeGenerals 是人事任命的候選。⭐ **與編成共用同一份過濾**
// （原版三條流程都走 `sub_17663`，docs/spec/148），差別只有：
// **人事一律排除君主，沒有開關**——`-lord-corps` 那個使用者裁定的差異
// （docs/spec/76 §3）只管編成。
func (g *game) freeGenerals() []int {
	return g.candidateGenerals(true)
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

	// 任命時選到「那裡已經有人」的兩則（`sub_16A9B` 的 `cx = 34h`／
	// `sub_16B71` 的 `cx = 35h`）。
	cityStaffedTalk    = 0x34 // #52「{2}已有{1}大人前去赴任了。」
	factionStaffedTalk = 0x35 // #53「{3}勢力已有{1}大人前去赴任了。」

	// 選武將那一步的狀態列（`mov cx, 9`）。
	pickOfficialTalk = 0x09 // #9「請選擇任命之武將。」

	// 那位官員自己說的一句：**八格一組**，由武將記錄 +0x1E 選組內第幾個，
	// 肖像取 +0x01（`sub_18810` 的 `ah`／`al`，docs/spec/142）。
	governorAssignedTalk  = 0x19C
	diplomatAssignedTalk  = 0x19D
	governorDismissedTalk = 0x1A2
	diplomatDismissedTalk = 0x1A3
)

// officialSays 把那位官員的一句掛出來（變體組 ＋ 他自己的肖像）。
func (g *game) officialSays(base, who int) {
	if g == nil || g.world == nil || who < 0 || who >= len(g.world.Generals) {
		return
	}
	gen := &g.world.Generals[who]
	g.enqueueTalkWithPortrait(
		resolveBattleTalkIndex(base, gen.TalkVariant), nil, gen.Portrait)
}

func (g *game) pickCityForGovernor() {
	g.setStatusTalk(governorAssignTalk, nil)
	rows := g.playerCities()
	if len(rows) == 0 {
		g.lastEvent = "沒有據點"
		g.list = nil
		return
	}
	// **先選據點，再選武將**（說明書：內政官任命的順序）。
	//
	// ⭐ **選完回到據點清單**（原版 `sub_16A9B` 的 `jmp loc_16AA4`，
	// docs/spec/142）：不論成功、那裡已經有人、還是選武將時取消，
	// 都回到這一張清單，**右鍵才離開**。
	g.cityList(rows, "選要派內政官的據點　Enter 決定　ESC 取消", func(city int) bool {
		c := &g.world.Cities[city]
		if c.Governor != NoOfficial {
			g.enqueueTalk(cityStaffedTalk, map[byte]string{
				'2': big5(c.Name),
				'1': big5(g.world.Generals[c.Governor].Name),
			})
			return false // 回清單
		}
		free := g.freeGenerals()
		if len(free) == 0 {
			g.lastEvent = "沒有可派任的武將"
			return false
		}
		name := big5(c.Name)
		g.setStatusTalk(pickOfficialTalk, nil)
		g.generalList(free, "選要派去 "+name+" 的武將　Enter 決定", func(who int) bool {
			// 職務（武將 +0x17 ＝ 2）與據點 +0x19 是同一步寫的，
			// 收在規則層（docs/spec/143 §2）。
			g.world.AssignGovernor(city, who)
			g.officialSays(governorAssignedTalk, who)
			// ★ 回到據點清單——**等那一句被按掉之後**（docs/spec/142 §3.1）。
			// 原版的 `sub_18810` 擋在這裡，訊息還掛著時畫面上是
			// 武將一覽 ＋ 狀態列 #9，不是回去以後的據點一覽。
			g.afterTalk(g.pickCityForGovernor)
			return false
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
		// 原版 `sub_16B4F`：先寫 0xFF、再看舊值是不是 0xFF；
		// 有人的話連職務與經費餘額一起清（docs/spec/143 §2）。
		old := g.world.DismissGovernor(city)
		if old < 0 {
			g.enqueueTalk(nobodyPostedTalk,
				map[byte]string{'2': big5(c.Name)})
			return false // ★ 回清單（原版 `jmp loc_16B11`）
		}
		g.officialSays(governorDismissedTalk, old)
		return false // ★ 成功也回清單
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
	// 與內政官任命同形（原版 `sub_16B71` 的 `jmp loc_16B7A`，docs/spec/142）。
	g.factionList(rows, "選要派外交官的勢力　Enter 決定　ESC 取消", func(f int) bool {
		fa := &g.world.Factions[f]
		name := big5(g.world.LordName(f))
		if fa.Diplomat != NoOfficial {
			g.enqueueTalk(factionStaffedTalk, map[byte]string{
				'3': name,
				'1': big5(g.world.Generals[fa.Diplomat].Name),
			})
			return false
		}
		free := g.freeGenerals()
		if len(free) == 0 {
			g.lastEvent = "沒有可派任的武將"
			return false
		}
		g.setStatusTalk(pickOfficialTalk, nil)
		g.generalList(free, "選要派去 "+name+" 的武將　Enter 決定", func(who int) bool {
			g.world.AssignDiplomat(f, who)
			g.officialSays(diplomatAssignedTalk, who)
			// ★ 回到勢力清單——同樣等那一句被按掉（docs/spec/142 §3.1）。
			g.afterTalk(g.pickFactionForDiplomat)
			return false
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
		// 原版 `sub_16C2A`：先寫 0xFF、再看舊值；有人的話連職務與
		// 經費餘額一起清（docs/spec/143 §2）。
		old := g.world.DismissDiplomat(f)
		if old < 0 {
			g.enqueueTalk(factionNobodyPostedTalk,
				map[byte]string{'3': big5(g.world.LordName(f))})
			return false // ★ 回清單
		}
		g.officialSays(diplomatDismissedTalk, old)
		return false // ★ 成功也回清單
	})
}

// factionList 開一張勢力清單。交友度是外交官要不要派的主要依據。
func (g *game) factionList(rows []int, hint string, pick func(int) bool) {
	g.openFactionPicker(rows, hint, pick)
}


