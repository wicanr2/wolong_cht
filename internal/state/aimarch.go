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

// aiStage0 是 `sub_1439D`：**先把意圖落實成行軍目標**，已經站在那裡才換檔。
//
// ⭐ 第一件事是 `sub_14548`（`retarget`）：把 `+0x14`／`+0x16`／`+0x18`
// 設成 `+0x20` 指的據點。AI 挑目標的那一支（`sub_1440F`／`aiStage2`）
// **只寫意圖**，落實是這裡的事——所以從「決定去哪」到「踏出第一步」
// 要跨三個更新週期（docs/spec/170）。
func (w *World) aiStage0(i int) {
	c := &w.Corps[i]
	if !w.retarget(i, c.Ordered) {
		return // 還沒到，這一拍只設目標
	}
	if w.citySpecific(c.Node) {
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
	// ⚠ **「在首都」要腳踩在上面**，不是 `Node` 記著它。行軍中 `Node`
	// 留著出發那一站——從首都出發的軍團會被讀成「已經到家」
	// （原版 `sub_14548` 比三個欄位：座標兩格 ＋ `+0x0E`）。
	if w.onCity(i) && c.Node == w.Factions[c.Faction].Capital {
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
	// ⚠ **只寫意圖，不碰行軍目標。** 原版 `sub_1440F` 寫的是 `+0x20`
	// 與位元 1，`+0x14` 要等下一拍的 `sub_1439D` 才落實（docs/spec/170）。
	// 兩件事併成一步會讓 AI 早一個更新週期出發，整條行軍軌跡跟著平移。
	c.Ordered = dest
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
		// 原版這裡還設位元 1 ＝「下一步要重算」（remake 未建模）。
		c.Ordered = capital
	}
	// ⭐ `sub_14548` 是**無條件**呼叫的（`loc_144C0` 落下去就是），
	// 不管意圖有沒有變——三個目標欄位每一輪都會被重寫，
	// 而它回 CF=1（座標兩格 ＋ `+0x0E` 都相同）才轉 Stage 9。
	if w.retarget(i, capital) {
		c.Stage = StageResupply
		return
	}
	w.routIfBlocked(i)
}

// retargetAndReplan 是 AI 改行軍目標那一刻該做的事（docs/spec/177）。
//
// ⚠ **與玩家下指令的 `March` 是兩回事**，兩點都要照抄：
//
//   - `coords` ＝ 要不要一起寫 `+0x16`／`+0x18`。`sub_144A9`／`sub_144D6`
//     經 `sub_14548` **會**寫（三個欄位一起），`sub_1474A`（戰後退卻）
//     **只寫 `+0x14`／`+0x20`**，座標留著上一次寫進去的值（§1.4）。
//   - 軍團走在邊上時，原版設位元 1、下一拍 `sub_147BB` 只決定「往這條邊的
//     哪一端」（`replanOnLeg`），不從出發據點重排整條路線。照 `March` 走
//     會排出一條從據點起算的路線，與軍團現在的座標對不上，`step` 就
//     走不動——**軍團原地凍住**。
func (w *World) retargetAndReplan(i, node int, coords bool) {
	if i < 0 || i >= numCorps || !validCity(node) {
		return
	}
	c := &w.Corps[i]
	if c.LinkAddr != 0 {
		c.TargetNode, c.Ordered = node, node
		if coords {
			c.TargetX, c.TargetY = w.Cities[node].X, w.Cities[node].Y
		}
		w.replanOnLeg(i)
		return
	}
	x, y := c.TargetX, c.TargetY
	_ = w.March(i, node)
	if !coords {
		c.TargetX, c.TargetY = x, y
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
	// ⭐ **走在路上時，下一站就是這條邊的端點**，不必再往前算一步
	// （§2.1）。站在節點上才走下面的「往首都的第一站」。
	if end, onRoad := w.retreatEndpoint(i, capital); onRoad {
		return end
	}
	from := c.Node
	if from == capital || w.roads == nil {
		return capital
	}
	path := w.roads.Route(from, capital)
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

// retreatEndpoint 是 `sub_1487B` 走在路上時的那一半：
// **回家的下一站就是這條邊的某一個端點**。
//
// `+0x0E` ≥ `800h` 表示軍團走在路上，那一欄是連結記錄的位址。原版把這條
// 邊的兩端寫進廣度優先的**終止條件**（`loc_1491B` 的
// `mov cs:word_149B8, bx` ／ `mov cs:word_149BE, cx`，那一段是自我修改碼），
// 然後**從首都往外搜**（`mov si, ax`，ax ＝ 首都 × 8）；搜到之後
// `mov bx, [bx+6]` 取的就是被搜到的那一端。
//
// ⇒ 兩端都屬於自己時退到**離首都近的那個**，只有一端屬於自己就退到那一端，
// 兩端都不是自己的地就退不了。回傳的第二個值是「軍團是不是走在路上」——
// false 表示它站在節點上，那時走的是另一條路（往首都的第一站）。
//
// ⚠ 兩端到首都**等距**時原版由廣度優先的展開順序決定，還沒解
// （docs/spec/46 §6）。
func (w *World) retreatEndpoint(i, capital int) (int, bool) {
	c := &w.Corps[i]
	if c.LinkAddr == 0 || w.roads == nil {
		return 0, false
	}
	a, b, ok := w.roads.EdgeByLink(c.LinkAddr)
	if !ok {
		// 缺道路圖或這條邊不在圖裡：當成站在節點上，不要整支停擺。
		return 0, false
	}
	oa, ob := w.nodeOwnedBy(a, c.Faction), w.nodeOwnedBy(b, c.Faction)
	switch {
	case oa && ob:
		da, db := w.roads.Distance(a, capital), w.roads.Distance(b, capital)
		if db >= 0 && (da < 0 || db < da) {
			return b, true
		}
		return a, true
	case ob:
		return b, true
	case oa:
		return a, true
	}
	return -1, true // `loc_14903` 的 STC
}

// nodeOwnedBy 是 `cmp dl, es:[di+841h]`：那個節點的所屬欄等於這個勢力。
func (w *World) nodeOwnedBy(node, faction int) bool {
	return node >= 0 && node < len(w.Cities) && w.Cities[node].Owner == faction
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
	//
	// ⚠ **判準是「腳踩在據點上」，不是「`Node` 記著哪個據點」**：
	// 行軍中 `Node` 留著出發那一站，光看它會把野外遭遇與城外對峙的
	// 攻方讀成「站在自家城裡」而不退（原版 `+0x0E` 這時 ≥ `800h`，
	// docs/spec/46 §2、docs/spec/175 §3）。
	if w.onCity(i) && w.Cities[c.Node].Owner == c.Faction {
		c.Stage = StageWaitMorale
		return false
	}
	next := w.nextHopHome(i)
	if next < 0 {
		return true
	}
	// ⚠ `sub_1474A` 只寫 `+0x14`／`+0x20` 與位元 1，**不碰 `+0x16`／`+0x18`**
	// （docs/spec/177 §1.4）——它不經過 `sub_14548`。
	w.retargetAndReplan(i, next, false)
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
