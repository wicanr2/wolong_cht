package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
)

// 兵記錄 `+0x18`（戰力）的算式。出處 `KI.EXE` 的 `sub_19B6D`（一般部隊）、
// `sub_19B40`（大將）、`sub_19C13`（取主將能力）——docs/re/78、docs/spec/115。

// kindBonus 是每兵種係數，原始 bytes 在 `seg000:9C0F`：`1E 04 0C 00`。
//
// ⚠ 原版索引的是**兵種編碼 1／2／3**（`bx = k − 1` 之後查表），
// 而 `army.TroopType` 是 0 起算的騎馬／弓兵／步兵——**同一個順序，差一**。
// 這裡直接用 Go 的值當索引。
var kindBonus = [economy.NumTroopTypes]int{30, 4, 12}

// 戰場類別。`sub_19A33` 把 `byte_10D34`（**戰場編號**）分三段寫進
// `byte_1D34B`，而那三段正好對上武將記錄的三個適性欄：
//
//	類別 0（攻城；戰場編號 ＝ 據點編號 0–191）→ +0x0E 攻城適性
//	類別 1（陸上；0xC0–0xD0）                 → +0x0F 陸戰適性
//	類別 2（海上；≥ 0xD1）                    → +0x10 海戰適性
//
// ⭐ 「陸上」「海上」是**原版自己的字**：`sub_1C955` 用 `byte_1D34B`
// 選側欄標題的第一段，字串在 `CS:0xD45E`／`0xD464`（docs/re/60 §11）。
//
// ⛔ **戰場編號不是大地圖的圖塊值。** 野戰的編號由 `sub_14B63` 從五格
// 地形算出來（docs/re/05、docs/re/58 §4）——橋（圖塊 `0xCA`）的編號是
// `0xD1`–`0xD4` ＝ 海上，而圖塊 `0xCA` 自己 < `0xD1`。
const (
	siegeAptitude = 0
	fieldAptitude = 1
	waterAptitude = 2

	// fieldNumberLow／fieldNumberHigh 是 `sub_19A33` 的兩道門檻。
	fieldNumberLow  = 0xC0
	fieldNumberHigh = 0xD1
)

// aptitudeIndex 是這一場戰鬥要讀主將的哪一個適性欄。
// 參數是**戰場編號**（`byte_10D34`），不是圖塊。
func aptitudeIndex(field int) int {
	switch {
	case field < fieldNumberLow:
		return siegeAptitude
	case field < fieldNumberHigh:
		return fieldAptitude
	default:
		return waterAptitude
	}
}

// soldierPower 是一般部隊每個兵的戰力：`((統率 + 適性) × 3 + 兵種係數) ÷ 4`。
//
// ⚠ `sub_19B6D` 裡還有一段依戰場類別調整騎兵的碼，**那是死碼**——
// `cmp al, 1` 的 `al` 在比較之前已經被換成「兵種 × 18」（docs/re/78 §3.1）。
// 照抄原版就是不要那兩條。
func soldierPower(command, aptitude int, kind army.TroopType) int {
	bonus := 0
	if k := int(kind); k >= 0 && k < len(kindBonus) {
		bonus = kindBonus[k]
	}
	p := ((command+aptitude)*3 + bonus) / 4
	if p < 0 {
		p = 0
	}
	return p
}

// leaderPower 是大將那一格的戰力：`(武力 × 2 + 適性) × 2`。
func leaderPower(martial, aptitude int) int {
	p := (martial*2 + aptitude) * 2
	if p < 0 {
		p = 0
	}
	return p
}

// leaderHP 是大將那一格的開場體力：`max(70, (武力 × 4 + 50) × 士氣 ÷ 100)`。
//
// ⚠ **不是軍團士氣。** 一般兵才是（docs/spec/61）；`sub_19AF4` 先填滿
// 48 個兵，再用 `sub_19B40` 把第 0 號蓋掉。
func leaderHP(martial, morale int) int {
	hp := (martial*4 + 50) * morale / 100
	if hp < leaderHPFloor {
		hp = leaderHPFloor
	}
	return hp
}

// leaderHPFloor 是 `sub_19B40` 的 `cmp al, 46h`。
const leaderHPFloor = 0x46

// squadPowers 算出六個編成位置的戰力，以及大將那一格的戰力與體力。
//
// 軍團編號就是主將的武將編號（`sub_1291A`／`sub_16F26` 都直接換算），
// 所以 `w.Generals[corps]` 就是這一支軍團的主將。
func (w *World) squadPowers(corps, category int) (
	squads [army.Positions]int, lp, lhp int) {
	if w == nil || corps < 0 || corps >= len(w.Corps) || corps >= len(w.Generals) {
		return
	}
	c := &w.Corps[corps]
	g := &w.Generals[corps]
	apt := 0
	if i := category; i >= 0 && i < len(g.Aptitude) {
		apt = g.Aptitude[i]
	}
	for k, u := range c.Units {
		if u.Men == 0 {
			continue
		}
		squads[k] = soldierPower(g.Command, apt, u.Kind)
	}
	return squads, leaderPower(g.Martial, apt), leaderHP(g.Martial, c.Morale)
}
