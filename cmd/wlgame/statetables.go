package main

import (
	"encoding/hex"
	"log"
	"strings"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 狀態層對拍：把勢力表／據點表／武將表印成逐欄可比的表
// （`docs/spec/138`）。軍團表走 `logAliveCorps`，那一張先做過
// （`docs/playtest/71`）。
//
// ⭐ **名字印十六進位。** 原版側是 dosgolem，而它刻意零相依，
// 解不了 Big5；兩邊都印 hex 就逐 byte 可比，也不必猜編碼。
//
// ⚠ **補白是全形空格 `A1 40`，不是 ASCII 空白**——`TrimSpace`
// 砍不掉，名字後面會拖著看不見的兩個 byte。

// dumpStateTables 依旗標印出三張表，**印過一次就不再印**。
//
// ⚠ **有 `-shot-when` 時要跟著取樣點印**，不是啟動時印：兩邊要在同一個
// 遊戲時刻比表，而內政每小時都會動幾個據點的上昇值、防災值與城兵數
// （`docs/spec/138` §4.1）。
func (g *game) dumpStateTables() {
	if g.tablesDumped {
		return
	}
	if !g.listFactions && !g.listCities && !g.listGenerals {
		return
	}
	g.tablesDumped = true
	c := g.world.Clock
	log.Printf("狀態表取樣點：%d年%d月%d日 %d時", c.Year, c.Month, c.Day, c.Hour)
	if g.listFactions {
		logFactions(g)
	}
	if g.listCities {
		logCities(g)
	}
	if g.listGenerals {
		logGenerals(g)
	}
}

// nameHex 把定長的 Big5 名字去掉全形空格補白，回傳大寫十六進位。
func nameHex(s string) string {
	for strings.HasSuffix(s, "\xa1\x40") {
		s = s[:len(s)-2]
	}
	return strings.ToUpper(hex.EncodeToString([]byte(s)))
}

// logFactions 印勢力表（`docs/formats/08` §1.5）。
func logFactions(g *game) {
	if g.world == nil {
		return
	}
	n := 0
	for i := range g.world.Factions {
		f := &g.world.Factions[i]
		if !f.Alive {
			continue
		}
		n++
		log.Printf("勢力 %2d 君%3d 師%3d 都%3d 備%5d/%5d/%5d 將%3d 城%3d 團%2d 金%9d 戰%2d 氣%3d 敵%3d",
			i, f.Lord, f.Advisor, f.Capital,
			f.Reserves[0], f.Reserves[1], f.Reserves[2],
			f.Generals, f.Cities, f.Corps, f.Funds,
			f.Aggression, f.MoraleBase, f.InvasionTarget)
	}
	log.Printf("共 %d 個勢力", n)
}

// logCities 印據點表（`docs/formats/08` §1.6）。
//
// ⚠ 兩個欄位在 remake 是**換算過的值**，比對時要換回去
// （`docs/spec/138` §4）：上昇值原版存「實際值 ＋ 100」，
// 城兵數原版存「人數 ÷ 10」而 remake 直接留存值。
func logCities(g *game) {
	if g.world == nil {
		return
	}
	for i := range g.world.Cities {
		c := &g.world.Cities[i]
		adj := c.Neighbours
		log.Printf("據點 %3d 主%3d 名%-12s (%3d,%3d) 產%5d/%5d 昇%3d 災%3d 兵%3d/%3d 類%d 官%3d 原%3d 鄰%3d,%3d,%3d,%3d",
			i, c.Owner, nameHex(c.Name), c.X, c.Y,
			c.Production, c.ProductionCap,
			c.Growth+100, c.Prevention, c.Garrison, c.GarrisonCap,
			c.Kind, c.Governor, c.OwnerRecorded,
			adj[0], adj[1], adj[2], adj[3])
	}
	log.Printf("共 %d 個據點", len(g.world.Cities))
}

// generalDuty 把武將的職務還原成原版 `+0x17` 的五個值
// （0 無／1 出陣中／2 內政官／3 外交官／4 捕虜）。
//
// ⚠ **remake 沒有這個欄位。** 載入時 `+0x17` 被壓成
// `General.Posted bool`（`r[0x17] != 0`），職務本身改記在別的地方：
// 內政官記在據點的 `Governor`、外交官記在**被派駐**勢力的 `Diplomat`。
// 所以要比這一欄就得從那兩張表反查回來——資訊沒丟，只是換了地方。
func generalDuty(g *game, id int) int {
	if g.world.Generals[id].Captor != 0xFF {
		return 4
	}
	for i := range g.world.Cities {
		if g.world.Cities[i].Governor == id {
			return 2
		}
	}
	for i := range g.world.Factions {
		if g.world.Factions[i].Diplomat == id {
			return 3
		}
	}
	if g.world.Generals[id].Posted {
		return 1
	}
	return 0
}

// generalFlags 把 remake 拆開的四個 bool 組回記錄 `+0x00` 那個 byte。
//
// ⚠ **bit 0 沒有欄位可以組回來**：四個劇本裡只有劇本三的張衛設著它，
// 語意未解，remake 載入時就丟掉了（`docs/formats/08` §3）。
// 逐欄比會在那一筆報出來——**那是誠實的結果，不是工具壞掉**。
func generalFlags(g state.General) uint8 {
	var f uint8
	if g.Alive {
		f |= 0x80
	}
	if g.Sovereign {
		f |= 0x40
	}
	if g.VanishIfAffinityGone {
		f |= 0x20
	}
	if g.LoyalToDeath {
		f |= 0x10
	}
	return f
}

// logGenerals 印武將表，只印在場的（`docs/formats/08` §3）。
func logGenerals(g *game) {
	if g.world == nil {
		return
	}
	n := 0
	for i := range g.world.Generals {
		gen := &g.world.Generals[i]
		if !gen.Alive {
			continue
		}
		n++
		post := generalDuty(g, i)
		log.Printf("武將 %3d 名%-12s 呼%-12s 勢%3d 職%d 武%2d 統%2d 政%2d 適%2d/%2d/%2d 本%d 說%d 向%3d 舊%3d 評%3d 旗%02X",
			i, nameHex(gen.Name), nameHex(gen.Alias), gen.Faction, post,
			gen.Martial, gen.Command, gen.Politics,
			gen.Aptitude[0], gen.Aptitude[1], gen.Aptitude[2],
			gen.Tactic, gen.TalkVariant, gen.Affinity, gen.Captor,
			gen.Rules().Rating(), generalFlags(*gen))
	}
	log.Printf("共 %d 人在場", n)
}
