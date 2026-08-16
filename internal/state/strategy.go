package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/strategyai"
	"github.com/wicanr2/wolong_cht/internal/rules/threat"
)

// aiFormationTable 是 CS:6C4C 的 18 byte 表，六個槽各有三個候選兵種。
// 兵種值是原版的 1-based byte：1 騎馬、2 弓兵、3 步兵；這裡轉成 army 型別。
// 輸入雜湊與檔案位址見 RESEARCH-LOG.md 的 AI 段落。
var aiFormationTable = [army.Positions][3]army.TroopType{
	{army.Cavalry, army.Infantry, army.Archer},
	{army.Cavalry, army.Infantry, army.Archer},
	{army.Infantry, army.Cavalry, army.Archer},
	{army.Infantry, army.Cavalry, army.Archer},
	{army.Archer, army.Infantry, army.Cavalry},
	{army.Archer, army.Infantry, army.Cavalry},
}

// strategyCandidates 照 sub_12C52 掃據點記錄的四個鄰接槽，按勢力去重。
// Owner 用執行期 +0x01，Neighbours 用原始 +0x1C–+0x1F；這兩個選擇都不能
// 偷換成 OwnerRecorded 或僅由 MMAP 推導的道路圖。
func (w *World) strategyCandidates(faction int) []strategyai.Candidate {
	seen := map[int]bool{}
	var out []strategyai.Candidate
	for _, c := range w.Cities {
		if c.Owner != faction {
			continue
		}
		for k, n := range c.Neighbours {
			if c.Adjacency&(1<<k) == 0 || n < 0 || n >= len(w.Cities) {
				continue
			}
			o := w.Cities[n].Owner
			if o < 0 || o >= numFactions || o == faction || !w.Factions[o].Alive || seen[o] {
				continue
			}
			seen[o] = true
			out = append(out, strategyai.Candidate{
				Faction: o, Friendship: w.Friendship[faction][o],
			})
		}
	}
	return strategyai.SortCandidates(out)
}

func (w *World) strategyFaction(i int) strategyai.Faction {
	f := w.Factions[i]
	return strategyai.Faction{
		Alive:          f.Alive,
		Player:         i == w.Player,
		Cities:         f.Cities,
		Funds:          f.Funds,
		Reserves:       [3]int{f.Reserves[0], f.Reserves[1], f.Reserves[2]},
		Aggression:     f.Aggression,
		InvasionTarget: f.InvasionTarget,
	}
}

// runStrategicAI 是 sub_12BD9 的事件佇列轉接層。
//
// 原版先建 22 個排序緩衝區，再把遷都、協力、停戰、宣戰事件丟進佇列；
// remake 已保存 raw 資料、月度壓縮與事件 1／2／3／4／5／8／9／13 的已證實處理端。
// 排序、友好度漂移、目標尾端驗證、三道宣戰閘與雙向交戰值都仍由已讀證據
// 約束。
func (w *World) runStrategicAI(rng economy.Rand) ([]StrategyEvent, map[int]int) {
	orders := make([][]strategyai.Candidate, numFactions)
	for i := range w.Factions {
		if w.Factions[i].Alive {
			orders[i] = w.strategyCandidates(i)
		}
	}

	// sub_12D3A 的事件 8 是月結中「無侵攻目標時」的 25% 主動遷都。
	// 先抽這個事件，再跑 sub_12D58 的決策鏈，保持原版呼叫順序；
	// 實際首都改寫交給每小時 dispatch 的 sub_133EA → sub_133FD。
	for i := range w.Factions {
		if !w.Factions[i].Alive || w.Factions[i].InvasionTarget != diplomacy.NoTarget {
			continue
		}
		if rng.Next()&0xFF >= 0x40 {
			continue
		}
		_ = w.queueEvent(rng, i, 8, 0xFFFF, 0xFF)
	}

	// 先做每個勢力的月度交友度漂移。原版緩衝區已建好，所以後面的
	// 目標排序不重新排列；宣戰事件也延後到本輪評估完才套用。
	for i := range w.Factions {
		if !w.Factions[i].Alive {
			continue
		}
		if i == w.Player {
			w.driftPlayerFriendship(i, orders[i])
		} else if len(orders[i]) > 0 {
			w.driftAIFriendship(i, orders[i][0].Faction)
		}

		// sub_12D58 的呼叫順序是漂移後先跑 sub_12E33／sub_12E89，
		// 再跑 sub_12EFB。候選清單仍是月結前建立的排序結果，但兩個
		// 產生器讀取的交友度與國力是此刻的 live state。
		w.queueCooperationProposal(rng, i, orders[i])
		w.queueCeasefireProposals(rng, i, orders[i])
	}

	type declaration struct{ faction, target int }
	var declarations []declaration
	for i := range w.Factions {
		if !w.Factions[i].Alive || len(orders[i]) == 0 {
			continue
		}
		first := orders[i][0]
		fr := w.Friendship[i][first.Faction]
		f := &w.Factions[i]

		// 已有目標時，原版尾段要求排序後第一筆確實是交戰中的
		// 最敵對鄰居，否則把 +0x19 快取清回 0xFF。沒有目標時，
		// sub_12EFB 會直接用第一筆候選跑三道閘；事件稍後才把和平
		// 轉成交戰，所以這裡不能先用 AtWar() 把新宣戰擋掉。
		if f.InvasionTarget != diplomacy.NoTarget {
			if f.InvasionTarget != first.Faction || !fr.AtWar() {
				f.InvasionTarget = diplomacy.NoTarget
			}
			continue
		}
		if i == w.Player {
			continue
		}
		self, target := w.strategyFaction(i), w.strategyFaction(first.Faction)
		if strategyai.ShouldDeclareWar(self, target, strategyai.Candidate{
			Faction: first.Faction, Friendship: fr,
		}) {
			declarations = append(declarations, declaration{faction: i, target: first.Faction})
		}
	}

	var out []StrategyEvent
	for _, d := range declarations {
		// sub_12EFB 發的是事件 1；不要在月結邊界直接改寫 +0x19
		// 或交友度，否則會跳過 sub_131AE 的每十次節拍。
		if w.queueEvent(rng, d.faction, 1, uint16(0xFF00|d.target), 0xFF) {
			out = append(out, StrategyEvent{Faction: d.faction, Target: d.target, Corps: -1, Destination: -1})
		}
	}
	return out, nil
}

// queueCooperationProposal 重現 sub_12E33 的事件 2 產生端。
//
// 事件 Code 的高 byte 是玩家勢力；Param 低 byte 是正在侵攻的勢力，
// 高 byte 是它的侵攻目標。候選清單由 sub_12C52 建好，所以這裡只在清單
// 中尋找玩家，不用 MMAP 或幾何距離重新猜鄰接關係。
func (w *World) queueCooperationProposal(rng economy.Rand, faction int, ordered []strategyai.Candidate) bool {
	if faction < 0 || faction >= numFactions || faction == w.Player || !w.Factions[faction].Alive {
		return false
	}
	target := w.Factions[faction].InvasionTarget
	if target == diplomacy.NoTarget || target == w.Player || target < 0 || target >= numFactions ||
		!w.Factions[target].Alive {
		return false
	}
	playerAdjacent := false
	for _, c := range ordered {
		if c.Faction == w.Player {
			playerAdjacent = true
			break
		}
	}
	if !playerAdjacent || w.Friendship[faction][w.Player].Raw() < 0x80 ||
		w.Friendship[target][w.Player].Raw() < 0xA3 {
		return false
	}

	param := uint16(target)<<8 | uint16(faction)
	return w.queueEvent(rng, w.Player, 2, param, 0xFF)
}

// queueCeasefireProposals 重現 sub_12E89 的事件 3 產生端。
//
// sub_12C52 已將交戰中的鄰居排在和平鄰居之前，並以高位元作標記；因此
// 這裡只掃排序清單的交戰前綴。當自身國力被前綴累減到零以下時，原版從
// 造成不足的那一筆起發事件，但保留第一筆（通常是目前的侵攻目標）不談。
func (w *World) queueCeasefireProposals(rng economy.Rand, faction int, ordered []strategyai.Candidate) int {
	if faction < 0 || faction >= numFactions || faction == w.Player || !w.Factions[faction].Alive ||
		len(ordered) == 0 {
		return 0
	}

	remaining := strategyai.Power(w.strategyFaction(faction))
	hostilePrefix := 0
	trigger := -1
	for i, c := range ordered {
		if !c.Friendship.AtWar() {
			break
		}
		hostilePrefix = i + 1
		if trigger < 0 {
			remaining -= strategyai.Power(w.strategyFaction(c.Faction))
			if remaining <= 0 {
				trigger = i
			}
		}
	}
	if trigger < 0 {
		return 0
	}

	queued := 0
	for i := trigger; i < hostilePrefix; i++ {
		if i == 0 {
			continue
		}
		param := uint16(0xFF00 | ordered[i].Faction)
		if w.queueEvent(rng, faction, 3, param, 0xFF) {
			queued++
		}
	}
	return queued
}

// queueFundingRequests 重現 sub_15715／sub_1578F。這兩支在月結壓縮事件佇列
// 後、其他政略事件前執行；事件 4 的字高是據點編號，事件 5 的字高是勢力
// 編號，兩者的 Param 都是 sub_139E8 的初始要求金額。
//
// 原版只用記錄欄位作掃描條件：它不另外確認官員存活或所屬勢力，處理端
// （sub_132A9／sub_132E9）才再次確認指標仍有效。因此這裡也只做 Go 陣列
// 邊界檢查，不把額外的「看起來合理」條件混進發送端。
func (w *World) queueFundingRequests(rng economy.Rand) (governors, diplomats int) {
	if w.Player < 0 || w.Player >= numFactions {
		return 0, 0
	}

	// sub_15715：玩家據點中，有內政官且官員 +0x1A 為零者。
	for cityID, c := range w.Cities {
		if c.Owner != w.Player || c.Governor == noFaction ||
			c.Governor < 0 || c.Governor >= numGenerals || w.Generals[c.Governor].Budget != 0 {
			continue
		}
		// +0x10 是「上昇值 + 100」的存檔 byte；state.Growth 已扣掉
		// 100，所以不能直接拿 Growth 去和原版的 0xB4 比較。
		growthStored := c.Growth + 100
		gap := positiveGap(0xB4, growthStored) +
			positiveGap(0xB4, c.Prevention) +
			positiveGap(c.GarrisonCap, c.Garrison)
		amount := (gap >> 1) * 0x32
		if w.queueEvent(rng, cityID, 4, uint16(amount), 0xFF) {
			governors++
		}
	}

	// sub_1578F：除了玩家自身外，所有「派駐該勢力的外交官」經費為零
	// 的記錄。要求金額的 min 與和平／交戰旗標都交給 diplomacy.Demand，
	// 其輸入保留原始 byte，不改成只看低七位。
	for factionID, f := range w.Factions {
		if factionID == w.Player || f.Diplomat == noFaction ||
			f.Diplomat < 0 || f.Diplomat >= numGenerals || w.Generals[f.Diplomat].Budget != 0 {
			continue
		}
		amount := diplomacy.Demand(w.Friendship[w.Player][factionID], w.Friendship[factionID][w.Player])
		if w.queueEvent(rng, factionID, 5, uint16(amount), 0xFF) {
			diplomats++
		}
	}
	return governors, diplomats
}

func positiveGap(want, got int) int {
	if want <= got {
		return 0
	}
	return want - got
}

func (w *World) driftAIFriendship(faction, target int) {
	fr := w.Friendship[faction][target]
	if fr.AtWar() {
		if target != w.Player && fr.Value() < 50 {
			w.Friendship[faction][target] = fr.WithValue(fr.Value() + 1)
		}
		return
	}
	v := fr.Value()
	if v < 22 {
		v = 22
	}
	w.Friendship[faction][target] = fr.WithValue(v - 2)
}

func (w *World) driftPlayerFriendship(player int, ordered []strategyai.Candidate) {
	if len(ordered) > 0 {
		target := ordered[0].Faction
		fr := w.Friendship[player][target]
		wasWar := fr.AtWar()
		fr = fr.WithValue(fr.Value() - 1)
		if wasWar {
			fr = fr.WithValue(fr.Value() - 7)
		}
		w.Friendship[player][target] = fr
	}
	for i := range w.Factions {
		if i == player || !w.Factions[i].Alive {
			continue
		}
		w.Friendship[i][player] = w.Friendship[i][player].WithValue(
			w.Friendship[i][player].Value() - 1)
	}
}

// formAICorps 是 sub_145C1 → sub_16E8F 的保真垂直切片：選同勢力且未出陣
// 的最高武力武將，照 CS:6C4C 的六槽表挑兵種，再按原版的「每型剩餘槽平均
// 分配、單槽上限 100」扣預備兵。它刻意不呼叫玩家 UI 的 FormCorps，因為
// 玩家介面採固定 1,000 人一槽的現代輸入適配，而原版 AI 編成的預備兵尺度
// 是 0–100 的軍團槽兵力。
// formAICorps 編一支新軍團，目標由侵攻對象決定（原版的一般路徑）。
func (w *World) formAICorps(faction int) *StrategyEvent {
	return w.formAICorpsTo(faction, -1)
}

// formAICorpsTo 編一支新軍團。dest ≥ 0 時直接指定目標據點——
// 求援走的就是這一條（原版 `sub_14575` 把 `ch` 寫進軍團 `+0x20`）。
func (w *World) formAICorpsTo(faction, dest int) *StrategyEvent {
	if faction < 0 || faction >= numFactions || faction == w.Player {
		return nil
	}
	f := &w.Factions[faction]
	// 上限是 `max(5, 資金 ÷ 8192)` 減掉現有軍團數（原版 `sub_14575`，
	// docs/re/44 §4.1）。**先前這裡寫的是 `f.Corps != 0`**——那讓 AI
	// 永遠只有一支軍團，資金再多也不會擴軍。
	if !f.Alive || f.InvasionTarget == diplomacy.NoTarget ||
		threat.Budget(f.Funds, f.Corps) == 0 {
		return nil
	}

	leader := -1
	for i, g := range w.Generals {
		if !g.Alive || g.Faction != faction || g.Posted || g.Captor != noFaction {
			continue
		}
		if leader < 0 || g.Martial > w.Generals[leader].Martial {
			leader = i
		}
	}
	if leader < 0 || leader >= numCorps {
		return nil
	}

	var kinds [army.Positions]army.TroopType
	probe := f.Reserves
	for slot, choices := range aiFormationTable {
		found := false
		for _, kind := range choices {
			if probe[kind] < 0x32 {
				continue
			}
			kinds[slot] = kind
			probe[kind] -= 0x32
			found = true
			break
		}
		if !found {
			return nil
		}
	}

	var counts [economy.NumTroopTypes]int
	for _, kind := range kinds {
		counts[kind]++
	}
	var c Corps
	c.Alive, c.Faction, c.Morale = true, faction, f.MoraleBase
	// AI 自己編的軍團預設是委任（`sub_16E8F` 的 `or [si], 4`）。
	c.Delegated = true
	c.Ordered = f.Capital
	home := w.clampCity(f.Capital)
	c.Node, c.X, c.Y = home, w.Cities[home].X, w.Cities[home].Y
	c.TargetNode, c.TargetX, c.TargetY = c.Node, c.X, c.Y
	for slot, kind := range kinds {
		remaining := counts[kind]
		available := f.Reserves[kind]
		men := available / remaining
		if remainder := available % remaining; remainder > 0 {
			men += remainder
		}
		if men > 100 {
			men = 100
		}
		f.Reserves[kind] -= men
		counts[kind]--
		c.Units[slot] = combat.Unit{Men: men, Kind: kind}
		c.Men += men
	}
	c.Interval = IntervalMixed
	allCavalry := true
	for _, kind := range kinds {
		if kind != army.Cavalry {
			allCavalry = false
			break
		}
	}
	if allCavalry {
		c.Interval = IntervalCavalry
	}
	c.Timer = c.Interval

	w.Corps[leader] = c
	w.Generals[leader].Posted = true
	f.Corps++

	destination := dest
	if destination < 0 {
		destination = w.nearestFactionCity(home, f.InvasionTarget)
	}
	if destination >= 0 {
		_ = w.March(leader, destination)
	}
	return &StrategyEvent{
		Faction: faction, Target: f.InvasionTarget,
		Corps: leader, Destination: destination,
	}
}

func (w *World) nearestFactionCity(from, faction int) int {
	best, bestDistance := -1, int(^uint(0)>>1)
	for i, c := range w.Cities {
		if c.Owner != faction {
			continue
		}
		d := -1
		if w.roads != nil {
			d = w.roads.Distance(from, i)
		}
		if d < 0 {
			dx, dy := c.X-w.Cities[from].X, c.Y-w.Cities[from].Y
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			d = dx
			if dy > d {
				d = dy
			}
		}
		if d < bestDistance || (d == bestDistance && i < best) {
			best, bestDistance = i, d
		}
	}
	return best
}

// capitalNone 避免 strategy.go 直接依賴 capital 套件的 None 命名，並讓
// relocateCapital 的哨兵仍由 state 內部統一處理。
const capitalNone = -1

// cityNeighbours 把據點的四個鄰接槽翻成規則層的形狀。
//
// 槽的有效性看 +0x1C–+0x1F 的哨兵，**不看 +0x00 的低 4 位**——
// 那四位是「哪幾個鄰居是敵方」，拿來當「有沒有鄰居」會在據點換手之後
// 把真的鄰居濾掉（docs/re/44 §5）。
func (w *World) cityNeighbours(i int) []threat.Neighbour {
	c := &w.Cities[i]
	out := make([]threat.Neighbour, 0, threat.Slots)
	for _, n := range c.Neighbours {
		if n < 0 || n >= len(w.Cities) {
			out = append(out, threat.Neighbour{Site: -1})
			continue
		}
		nb := &w.Cities[n]
		f := 0
		if c.Owner >= 0 && c.Owner < numFactions && nb.Owner >= 0 && nb.Owner < numFactions {
			f = w.Friendship[c.Owner][nb.Owner].Raw()
		}
		out = append(out, threat.Neighbour{
			Site: n, Owner: nb.Owner, Occupancy: nb.Occupancy, Friendship: f,
		})
	}
	return out
}

// occupancyAt 數停在 (x, y) 那一格的軍團數。
//
// 原版存的是一張 384 × 256 的計數器陣列（`word_19872`），軍團出發時
// `dec`、到站時 `inc`。remake 直接數——**佔用數是導出值，記帳會漂**。
func (w *World) occupancyAt(x, y int) int {
	n := 0
	for i := range w.Corps {
		if w.Corps[i].Alive && w.Corps[i].X == x && w.Corps[i].Y == y {
			n++
		}
	}
	return n
}

// refreshCityThreat 重算一個據點的佔用數、敵方鄰接遮罩與威脅量
// （原版 `sub_13EFD` 的佔用圖抄寫 ＋ `sub_13FA9`）。docs/re/44 §1–§2。
func (w *World) refreshCityThreat(i int, rng economy.Rand) []TalkNotice {
	if i < 0 || i >= len(w.Cities) {
		return nil
	}
	c := &w.Cities[i]
	// 求援冷卻每輪 −1（原版 sub_13EFD 的第一件事）。
	if c.ReliefCooldown > 0 {
		c.ReliefCooldown--
	}
	// 據點易主時，原主的勢力記錄留下「被佔走的據點編號」。
	if c.OwnerRecorded != c.Owner &&
		c.OwnerRecorded >= 0 && c.OwnerRecorded < numFactions {
		w.Factions[c.OwnerRecorded].LostSite = i
	}
	c.Occupancy = w.occupancyAt(c.X, c.Y)
	ns := w.cityNeighbours(i)
	c.Adjacency, c.EnemyNeighbours = threat.EnemyMask(c.Owner, ns)
	if c.Owner == threat.Neutral || c.Owner < 0 || c.Owner >= numFactions {
		c.Threat = 0
		c.Threatened, c.Specific = false, false
		return nil
	}
	r := threat.Scan(c.Owner, w.invasionTarget(c.Owner), c.EnemyNeighbours, ns)
	c.Threat = r.Level
	// 記錄 +0x00 的位元 7／6。原版 `sub_14028` 先清再設，所以是每輪
	// 重算的派生旗標，不是狀態——AI 軍團的 Stage 0–2 讀它們。
	c.Threatened, c.Specific = r.Threatened, r.Specific
	return w.relieve(i, r, rng)
}

// invasionTarget 回傳勢力的侵攻目標，翻成規則層的哨兵值。
func (w *World) invasionTarget(faction int) int {
	if faction < 0 || faction >= numFactions {
		return threat.NoTarget
	}
	t := w.Factions[faction].InvasionTarget
	if t == diplomacy.NoTarget {
		return threat.NoTarget
	}
	return t
}

// requestRelief 是據點求援（原版 `sub_140C9` ＋ `sub_140B3`，docs/re/40 §4）。
//
// 冷卻中就不重複求援。玩家的據點跳訊息 #38；AI 的據點直接編軍團過來，
// 冷卻寫「離首都的距離」。兩條路最後都在勢力記錄留下求援的據點編號。
func (w *World) requestRelief(site, want int, rng economy.Rand) []TalkNotice {
	c := &w.Cities[site]
	if c.ReliefCooldown != 0 {
		return nil
	}
	f := &w.Factions[c.Owner]
	var notices []TalkNotice
	if c.Owner == w.Player {
		notices = append(notices, TalkNotice{
			Index: threat.ReliefMessage, City: site, Faction: c.Owner,
		})
		c.ReliefCooldown = threat.PlayerCooldown(rng.Next())
	} else {
		for n := threat.Budget(f.Funds, f.Corps); n > 0 && want > 0; n-- {
			if w.formAICorpsTo(c.Owner, site) == nil {
				break
			}
			want--
		}
		cap := w.clampCity(f.Capital)
		c.ReliefCooldown = threat.AICooldown(c.X, c.Y, w.Cities[cap].X)
	}
	f.ReliefSite = site
	return notices
}

// dispatchGarrison 把已經停在這個據點的軍團調去 target（原版 `sub_14155`）。
func (w *World) dispatchGarrison(site, target, want, skip int, rng economy.Rand) {
	gs := make([]threat.Garrison, len(w.Corps))
	for i := range w.Corps {
		cp := &w.Corps[i]
		// 位元 2 ＝「委任」（docs/re/45）：求援只調得動交給電腦指揮的軍團，
		// 玩家自己指揮的不會被搶走。`+0x23 >= 8` 是待解體，也調不動。
		gs[i] = threat.Garrison{
			At:    cp.Node,
			Ready: cp.Alive && cp.Faction == w.Cities[site].Owner && cp.Delegated,
			Stage: cp.Stage,
		}
		if !cp.Alive {
			gs[i].At = -1
		}
	}
	for _, i := range threat.Dispatch(gs, site, want, skip, func() int { return rng.Next() }) {
		_ = w.March(i, target)
	}
}

// relieve 跑一個據點的威脅反應（原版 `sub_13F74` → `sub_14028`／`sub_14057`）。
//
// 順序照抄：貼身威脅走立刻求援那條，其餘走機率路徑；機率路徑再依
// 「這個據點站了幾支軍團」分成求援與派兵兩支。
func (w *World) relieve(site int, r threat.Result, rng economy.Rand) []TalkNotice {
	c := &w.Cities[site]
	if !r.Threatened {
		c.ReliefCooldown = 0
		return nil
	}
	// sub_14028：有具體目標而且這一格沒有軍團 → 立刻求援一支。
	if r.Specific && c.Occupancy == 0 {
		return w.requestRelief(site, 1, rng)
	}
	// sub_14057：從最多四個威脅裡隨機挑一個。
	if len(r.Targets) == 0 {
		return nil
	}
	pick := r.Targets[rng.Next()&3%len(r.Targets)]
	want := threat.Requested(c.Threat, c.Occupancy)
	if c.Occupancy <= 1 {
		// ⚠ 這一條**玩家的據點不走**（原版 `cmp cl, [si+841h] / jz 結束`）。
		// 玩家只從上面那條貼身威脅的路徑收到求援訊息。
		if want == 0 || c.Owner == w.Player {
			return nil
		}
		return w.requestRelief(site, want, rng)
	}
	if want < 1 {
		want = 1
	}
	w.dispatchGarrison(site, pick, want, c.Occupancy-want, rng)
	c.ReliefCooldown = 0
	return nil
}
