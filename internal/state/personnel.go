package state

// 內政官與外交官的任命／解任。原版是四支：`sub_16A9B`（任命內政官）、
// `sub_16B4F`（解任內政官）、`sub_16B71`（任命外交官）、`sub_16C2A`
// （解任外交官）。UI 的流程（不過濾、選完回清單、官員說一句）在
// docs/spec/142，這裡只放**狀態怎麼變**——docs/spec/143。
//
// ⭐ 任命寫兩格、解任清三格：
//
//	任命  武將 +0x17 ＝ 2／3      據點 +0x19 ／勢力 +0x2A ＝ 武將編號
//	解任  武將 +0x17 ＝ 0         武將 +0x1A ＝ 0（經費沒收）
//	                              據點 +0x19 ／勢力 +0x2A ＝ 0xFF
//
// **解任連經費一起歸零**是原版的行為（`mov byte [bx+1Ah], 0`，
// docs/re/25 §3.1）——不歸零的話下一任會帶著前任沒花完的錢。

// AssignGovernor 把 who 派到 city 當內政官。那裡已經有人就不動，回 false
// （原版由呼叫端先 `cmp [bx+19h], 0FFh` 擋掉，這裡收在同一支）。
func (w *World) AssignGovernor(city, who int) bool {
	if w == nil || city < 0 || city >= len(w.Cities) ||
		who < 0 || who >= len(w.Generals) {
		return false
	}
	c := &w.Cities[city]
	if c.Governor != noGovernorSlot {
		return false
	}
	w.Generals[who].Duty = DutyGovernor
	c.Governor = who
	return true
}

// DismissGovernor 解任 city 的內政官，回傳被解任的武將編號；
// 那裡本來就沒人就回 noGovernor（−1）。
//
// 照原版 `sub_16B4F` 的順序：**先無條件寫 0xFF，再看舊值**。
func (w *World) DismissGovernor(city int) int {
	if w == nil || city < 0 || city >= len(w.Cities) {
		return noGovernor
	}
	c := &w.Cities[city]
	old := c.Governor
	c.Governor = noGovernorSlot
	if old < 0 || old >= len(w.Generals) {
		return noGovernor
	}
	g := &w.Generals[old]
	g.Budget = 0 // ← 未用完的經費沒收
	g.Duty = DutyNone
	return old
}

// AssignDiplomat 把 who 派駐到 faction 當外交官。與 AssignGovernor 同形。
func (w *World) AssignDiplomat(faction, who int) bool {
	if w == nil || faction < 0 || faction >= len(w.Factions) ||
		who < 0 || who >= len(w.Generals) {
		return false
	}
	f := &w.Factions[faction]
	if f.Diplomat != noFaction {
		return false
	}
	w.Generals[who].Duty = DutyDiplomat
	f.Diplomat = who
	return true
}

// DismissDiplomat 召回派駐 faction 的外交官，回傳武將編號；
// 那裡本來就沒人就回 noGovernor（−1）。與 DismissGovernor 同形。
func (w *World) DismissDiplomat(faction int) int {
	if w == nil || faction < 0 || faction >= len(w.Factions) {
		return noGovernor
	}
	f := &w.Factions[faction]
	old := f.Diplomat
	f.Diplomat = noFaction
	if old < 0 || old >= len(w.Generals) {
		return noGovernor
	}
	g := &w.Generals[old]
	g.Budget = 0
	g.Duty = DutyNone
	return old
}

// NoOfficial 是「沒有派駐」的哨兵（原版的 0xFF），給呈現層判斷用。
const NoOfficial = noGovernorSlot
