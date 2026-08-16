package state

// 電腦勢力的行軍決策鏈：軍團走到目標據點時，Stage 0–3 各有一支 handler。
//
// 規格 `docs/spec/40`，機器碼出處 `docs/re/65`
// （`sub_1439D`／`sub_143AF`／`sub_1440F`／`sub_14466`）。
//
// ⭐ 這條鏈是勢力 `LostSite`（+0x17）與 `ReliefSite`（+0x16）唯一的**讀取端**：
// 軍團走到據點、發現沒事做，就去把這兩格待辦領走。沒有它，AI 的軍團
// 編出來、走到第一個目標之後就停在那裡不動。

import "github.com/wicanr2/wolong_cht/internal/rules/army"

const (
	// aiHomeThreshold 是 `cmp word [si+4], 12Ch`：兵力 ≤ 300 點（3,000 人）
	// 就放下手邊的事回首都補兵。
	aiHomeThreshold = 300
	// aiSlotThreshold 是 `mov al, 1Eh`：補完兵之後任一槽不到 30 點
	// （300 人）就整團解散。
	aiSlotThreshold = 30
	// aiCrowded 是 `cmp byte [di+18h], 2` 的比較值：這一格已經站了
	// 超過兩支軍團，就分一支出去打。
	aiCrowded = 2
	// aiHoldAlone 是 Stage 2 的 `cmp byte [di+18h], 1`：受威脅而且只有
	// 自己這一支，就退回 Stage 1 留守。
	aiHoldAlone = 1
)

// aiArrive 是非玩家軍團的 Stage 0–3 分派（原版 `funcs_1434F` 的第 4–7 筆）。
func (w *World) aiArrive(i int, rng rander) {
	c := &w.Corps[i]
	switch c.Stage {
	case 0:
		w.aiStage0(i)
	case 1:
		w.aiStage1(i, rng)
	case 2:
		w.aiStage2(i)
	default:
		// 原版的第 7 筆是 Stage 3；玩家那一半把 0–3 全走同一支，
		// AI 這一半只有 3 落到這裡。
		w.aiStage3(i)
	}
}

// rander 只要 Next()，讓 Stage 1 的延遲不必牽進整個 combat 套件。
type rander interface{ Next() int }

// aiStage0 是 `sub_1439D`：到了目標就換檔，除非這一格的威脅還有具體目標。
func (w *World) aiStage0(i int) {
	c := &w.Corps[i]
	if w.citySpecific(c.TargetNode) {
		return // 威脅有具體目標 → 原地駐守
	}
	c.Stage = 1
}

// aiStage1 是 `sub_143AF`：站在據點上決定下一步。
func (w *World) aiStage1(i int, rng rander) {
	c := &w.Corps[i]
	// 還在野外（節點編號 ≥ 256）就退回 Stage 0 繼續走。
	if army.KindOf(c.Node) == army.FieldNode {
		c.Stage = StageNormal
		return
	}
	if c.Men <= aiHomeThreshold {
		c.Stage = StageHomeResupply
		return
	}
	if w.citySpecific(c.TargetNode) {
		c.Stage = StageNormal
		return
	}
	// ⚠ **remake 差異**：原版這一行的 `di` 沒設就用，讀到的位址比據點
	// 記錄少 0x840（`docs/re/65` §3.2）。這裡實作作者意圖的版本——
	// 同一家族的 `sub_1440F` 就是這樣讀的。
	if !w.cityThreatened(c.TargetNode) || w.cityOccupancy(c.TargetNode) > aiCrowded {
		// 出擊前隨機等 1–8 個 tick，讓同一批軍團不會一起動。
		c.Timer = rng.Next()&7 + 1
		c.Stage = 2
		return
	}
	// 留守。順便看要不要補兵——只有站在首都而且未滿編才補。
	if c.Men >= army60000Points || c.Faction < 0 || c.Faction >= numFactions {
		return
	}
	if c.Node == w.Factions[c.Faction].Capital {
		c.Stage = StageResupply
	}
}

// aiStage2 是 `sub_1440F`：從勢力層的兩格待辦裡挑下一個目標。
func (w *World) aiStage2(i int) {
	c := &w.Corps[i]
	node := c.TargetNode
	if w.citySpecific(node) {
		c.Stage = 1
		return
	}
	if w.cityThreatened(node) && w.cityOccupancy(node) <= aiHoldAlone {
		c.Stage = 1 // 受威脅而且只有自己這一支 → 留守
		return
	}
	dest := w.takeAIErrand(c.Faction)
	if dest < 0 {
		// 沒事做：這一格還受威脅就留著，否則解散回收兵力。
		if w.cityThreatened(node) {
			return
		}
		c.Stage = StageDisband
		w.leaveCell(node)
		return
	}
	if c.Ordered != dest {
		_ = w.March(i, dest)
	}
	c.Stage = StageNormal
	w.leaveCell(node)
}

// aiStage3 是 `sub_14466`：補完兵之後的體檢。
//
// **六個槽逐一比 30 點**，任一不足就整團解散——空槽的兵力是 0，
// 一樣不過關，所以湊不齊六槽的軍團 AI 不留。
func (w *World) aiStage3(i int) {
	c := &w.Corps[i]
	for _, u := range c.Units {
		if u.Men < aiSlotThreshold {
			c.Stage = StageDisband
			return
		}
	}
	c.Stage = StageWaitMorale
}

// waitMorale 是 `sub_14483`：士氣沒回到勢力基準就留著不動。
func (w *World) waitMorale(i int) {
	c := &w.Corps[i]
	if c.Faction < 0 || c.Faction >= numFactions {
		c.Stage = 1
		return
	}
	if c.Morale < w.Factions[c.Faction].MoraleBase {
		return
	}
	c.Stage = 1
}

// headHomeResupply 是 `sub_144A9`：目標校正成首都，到了就轉補兵。
func (w *World) headHomeResupply(i int) {
	c := &w.Corps[i]
	if c.Faction < 0 || c.Faction >= numFactions {
		return
	}
	capital := w.clampCity(w.Factions[c.Faction].Capital)
	if c.Ordered != capital {
		_ = w.March(i, capital)
		w.routIfBlocked(i)
		return
	}
	if c.Node == capital {
		c.Stage = StageResupply
	}
}

// nextHopHome 是 `sub_1487B`：回首都的路上「下一個自己的據點」。
// 回傳 −1 表示**退不了**（沒有首都、沒有下一站，或下一站不是自己的地）。
//
// ⚠ 這支不是 Stage 10／11 用的。那兩檔走 `sub_144A9`／`sub_144D6`，
// 直接把目標設成首都、不逐站（[`docs/spec/46`](../../docs/spec/46-post-battle-retreat.md) §3）。
// 用到它的是戰後（`sub_1474A`）與據點失守（`sub_14DA4`）。
func (w *World) nextHopHome(i int) int {
	if i < 0 || i >= numCorps {
		return -1
	}
	c := &w.Corps[i]
	if c.Faction < 0 || c.Faction >= numFactions {
		return -1
	}
	capital := w.Factions[c.Faction].Capital
	if !validCity(capital) {
		return -1 // `cmp al, 0FFh` ⇒ 沒有首都就無處可退
	}
	if c.Node == capital || w.roads == nil {
		return capital
	}
	path := w.roads.Route(c.Node, capital)
	if len(path) < 2 {
		// 廣度優先回 carry 時原版拿首都本身當下一站，而且**不再檢查歸屬**
		// （`loc_1490C` 直接 CLC）。
		return capital
	}
	next := path[1]
	if !validCity(next) || w.Cities[next].Owner != c.Faction {
		return -1 // `cmp dl, [bx+841h]` 不符 ⇒ STC
	}
	return next
}

// retreatOrPerish 是 `sub_1474A` 的後半段：敗方退一站回家。
// 回傳 true 表示**退不了 ⇒ 壞滅**——那是壞滅的第三個入口，
// 士氣 0 與大將槽 0 之外的那一個（`docs/spec/46` §1）。
//
// won 對應原版的 `cl`：勝方（`cl == 0`）原地不動。
func (w *World) retreatOrPerish(i int, won bool) bool {
	if i < 0 || i >= numCorps || !w.Corps[i].Alive {
		return false
	}
	c := &w.Corps[i]
	if won {
		c.Stage = StageWaitMorale
		return false
	}
	// 站在自家據點上就不退（`cmp al, [si+1] / jz .stay`）。
	// 攻城時守方正是這一支，所以據點易主不受這條規則影響。
	if army.KindOf(c.Node) == army.CityNode && c.Node >= 0 && c.Node < len(w.Cities) &&
		w.Cities[c.Node].Owner == c.Faction {
		c.Stage = StageWaitMorale
		return false
	}
	next := w.nextHopHome(i)
	if next < 0 {
		return true
	}
	_ = w.March(i, next)
	// 兵力 ≤ 300 或退到首都就轉回首都補兵，否則先等士氣。
	if c.Men <= aiHomeThreshold || next == w.Factions[c.Faction].Capital {
		c.Stage = StageHomeResupply
		return false
	}
	c.Stage = StageWaitMorale
	return false
}

// routIfBlocked 是 `sub_147BB` 的 `0x8000` 分支（`docs/spec/43`）：
// **回家的路要穿過別人的地，這支軍團就敗走。**
//
// ⚠ 條件裡的 `Stage ≥ 10` 由呼叫端保證——只有 `headHomeResupply`
// 與 `arriveDisband` 會叫它，那兩支正是 Stage 10 與 11。
func (w *World) routIfBlocked(i int) {
	if w.Corps[i].Alive && w.returnBlocked(i) {
		w.routCorps(i)
	}
}

// takeAIErrand 領走勢力層的一件待辦，回傳目標據點編號（沒有就 −1）。
//
// **兩格都是「取走就清空」的一格佇列**（原版的 `xchg al, [bx+17h]`），
// 所以同一件待辦只會派出一支軍團。失土優先於求援，
// 但**資金吃緊的勢力跳過失土**——沒錢就不主動反攻。
func (w *World) takeAIErrand(faction int) int {
	if faction < 0 || faction >= numFactions {
		return -1
	}
	f := &w.Factions[faction]
	if !f.LowFunds {
		if site := f.LostSite; validCity(site) {
			f.LostSite = noSite
			return site
		}
	}
	if site := f.ReliefSite; validCity(site) {
		f.ReliefSite = noSite
		return site
	}
	return -1
}

// noSite 是「這一格待辦是空的」（原版的 0xFF 哨兵）。
const noSite = 0xFF

func validCity(site int) bool { return site >= 0 && site < army.NumCityNodes }

// leaveCell 是 `dec byte [di+18h]`：軍團決定離開時先手動把這一格的
// 佔用數減一，同一個 tick 裡後面才被處理的軍團才看得到。
//
// `Occupancy` 每輪由位置重算（是快取不是狀態），所以這個扣減只在
// 同一個 tick 內生效——語意與原版相同。
func (w *World) leaveCell(node int) {
	if !validCity(node) || node >= len(w.Cities) {
		return
	}
	if w.Cities[node].Occupancy > 0 {
		w.Cities[node].Occupancy--
	}
}

// dropGeneralCount／raiseGeneralCount 是勢力記錄 +0x18 的維護
// （原版 `sub_12AD2`：舊勢力 −1、新勢力 +1，`0xFF` 那一側不動）。
//
// **被俘、戰死、釋放都要走這兩支**，否則勢力的武將數會與實際人數脫節，
// 而那個數字是政略 AI 的輸入之一。
func (w *World) dropGeneralCount(faction int) {
	if faction < 0 || faction >= numFactions {
		return
	}
	if w.Factions[faction].Generals > 0 {
		w.Factions[faction].Generals--
	}
}

func (w *World) raiseGeneralCount(faction int) {
	if faction < 0 || faction >= numFactions {
		return
	}
	w.Factions[faction].Generals++
}

func (w *World) cityThreatened(node int) bool {
	if !validCity(node) || node >= len(w.Cities) {
		return false
	}
	return w.Cities[node].Threatened
}

func (w *World) citySpecific(node int) bool {
	if !validCity(node) || node >= len(w.Cities) {
		return false
	}
	return w.Cities[node].Specific
}

func (w *World) cityOccupancy(node int) int {
	if !validCity(node) || node >= len(w.Cities) {
		return 0
	}
	return w.Cities[node].Occupancy
}
