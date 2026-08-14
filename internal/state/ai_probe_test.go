package state

import (
	"sort"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/strategyai"
	"github.com/wicanr2/wolong_cht/internal/rules/threat"
)

const threatNeutral = threat.Neutral

func TestScenarioOneStrategyNeighbourOrder(t *testing.T) {
	w := load(t, 0)
	got := w.strategyCandidates(0)
	want := []int{13, 12, 11, 21, 2}
	if len(got) != len(want) {
		t.Fatalf("曹操鄰接勢力數 = %d，want %d", len(got), len(want))
	}
	for i, wantFaction := range want {
		if got[i].Faction != wantFaction {
			t.Fatalf("曹操第 %d 個候選 = %d，want %d", i, got[i].Faction, wantFaction)
		}
	}
}

func TestStrategicAIDiplomacyEventGenerators(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	for i := range w.events {
		w.events[i] = QueuedEvent{}
	}
	w.eventCursor = 0

	// sub_12E33：玩家在侵攻方的鄰接清單內，且玩家與侵攻方和平；
	// 被侵攻目標對玩家至少為 0x80|0x23，才會排入事件 2。
	invader, target := 1, 2
	w.Factions[invader].InvasionTarget = target
	w.Friendship[invader][w.Player] = diplomacy.Peace(50)
	w.Friendship[target][w.Player] = diplomacy.Peace(40)
	ordered := []strategyai.Candidate{
		{Faction: w.Player, Friendship: w.Friendship[invader][w.Player]},
	}
	if !w.queueCooperationProposal(rng.NewFixed(1), invader, ordered) {
		t.Fatal("符合四道閘時沒有排入事件 2")
	}
	var cooperation QueuedEvent
	for _, e := range w.events[:eventQueueDispatch] {
		if byte(e.Code) == 2 {
			cooperation = e
			break
		}
	}
	if cooperation.Code != uint16(w.Player)<<8|2 ||
		cooperation.Param != uint16(target)<<8|uint16(invader) {
		t.Fatalf("事件 2 payload = %#v，want code=0x%04X param=0x%04X", cooperation,
			uint16(w.Player)<<8|2, uint16(target)<<8|uint16(invader))
	}

	// sub_12E89：交戰鄰居按交友度排在前面；自身國力不足時，
	// 從造成不足的第二筆起排事件 3，但排除第一筆目前敵人。
	w = load(t, 0)
	w.Player = 0
	source, first, second := 1, 2, 3
	w.Factions[source].Reserves = [3]int{}
	w.Factions[source].Funds = 50_000
	w.Factions[source].Cities = 2
	for _, faction := range []int{first, second} {
		w.Factions[faction].Funds = 50_000
		w.Factions[faction].Cities = 2
		w.Factions[faction].Reserves = [3]int{400, 400, 400}
	}
	ordered = []strategyai.Candidate{
		{Faction: first, Friendship: diplomacy.Friendship(10)},
		{Faction: second, Friendship: diplomacy.Friendship(20)},
	}
	if got := w.queueCeasefireProposals(rng.NewFixed(1), source, ordered); got != 1 {
		t.Fatalf("事件 3 排入數 = %d，want 1", got)
	}
	var ceasefire QueuedEvent
	for _, e := range w.events[:eventQueueDispatch] {
		if byte(e.Code) == 3 {
			ceasefire = e
			break
		}
	}
	if ceasefire.Code != uint16(source)<<8|3 || ceasefire.Param != uint16(0xFF00|second) {
		t.Fatalf("事件 3 payload = %#v，want code=0x%04X param=0x%04X", ceasefire,
			uint16(source)<<8|3, uint16(0xFF00|second))
	}
}

// 月結的事件 4／5 產生端保留原版的選擇器：事件 4 用據點編號，事件 5
// 用勢力編號；兩筆的 Param 都是要求金額，不是另一個索引。
func TestFundingRequestGenerators(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	for i := range w.events {
		w.events[i] = QueuedEvent{}
	}

	cityID := -1
	for i, c := range w.Cities {
		if c.Owner == w.Player && c.Governor == noFaction {
			cityID = i
			break
		}
	}
	if cityID < 0 {
		t.Fatal("找不到玩家據點")
	}
	officer := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction == w.Player {
			officer = i
			break
		}
	}
	if officer < 0 {
		t.Fatal("找不到玩家武將")
	}
	w.Cities[cityID].Governor = officer
	w.Cities[cityID].Growth = 0       // 原始 +0x10 = 100
	w.Cities[cityID].Prevention = 100 // 原始 +0x11 = 100
	w.Cities[cityID].GarrisonCap = 100
	w.Cities[cityID].Garrison = 0
	w.Generals[officer].Budget = 0
	governors, diplomats := w.queueFundingRequests(rng.NewFixed(1))
	if governors != 1 || diplomats != 0 {
		t.Fatalf("事件 4／5 產生數 = %d/%d，want 1/0", governors, diplomats)
	}
	wantAmount := ((0xB4 - 100) + (0xB4 - 100) + (100 - 0)) >> 1 * 0x32
	var got QueuedEvent
	for _, e := range w.events[:eventQueueDispatch] {
		if byte(e.Code) == 4 {
			got = e
			break
		}
	}
	if got.Code != uint16(cityID)<<8|4 || got.Param != uint16(wantAmount) {
		t.Fatalf("事件 4 payload = %#v，want code=0x%04X param=%d", got,
			uint16(cityID)<<8|4, wantAmount)
	}

	w = load(t, 0)
	w.Player = 0
	for i := range w.events {
		w.events[i] = QueuedEvent{}
	}
	target := 1
	diplomat := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction != w.Player {
			diplomat = i
			break
		}
	}
	if diplomat < 0 {
		t.Fatal("找不到非玩家武將")
	}
	w.Factions[target].Diplomat = diplomat
	w.Generals[diplomat].Budget = 0
	w.Friendship[w.Player][target] = diplomacy.Peace(30).WithWar(true)
	w.Friendship[target][w.Player] = diplomacy.Peace(60)
	governors, diplomats = w.queueFundingRequests(rng.NewFixed(1))
	if governors != 0 || diplomats != 1 {
		t.Fatalf("事件 4／5 產生數 = %d/%d，want 0/1", governors, diplomats)
	}
	got = QueuedEvent{}
	for _, e := range w.events[:eventQueueDispatch] {
		if byte(e.Code) == 5 {
			got = e
			break
		}
	}
	wantAmount = diplomacy.Demand(w.Friendship[w.Player][target], w.Friendship[target][w.Player])
	if got.Code != uint16(target)<<8|5 || got.Param != uint16(wantAmount) {
		t.Fatalf("事件 5 payload = %#v，want code=0x%04X param=%d", got,
			uint16(target)<<8|5, wantAmount)
	}
}

func TestStrategicAIScenarioOneProducesEnemyWarPath(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	w.EnableStrategicAI()
	r := rng.NewFixed(17)
	months, declarations, formed, battles := 0, 0, 0, 0
	for months < 6 {
		ev := w.Tick(r)
		for _, se := range ev.Strategy {
			if se.Corps < 0 {
				declarations++
				t.Logf("宣戰：faction=%d target=%d", se.Faction, se.Target)
			} else {
				formed++
				name := ""
				if se.Destination >= 0 && se.Destination < len(w.Cities) {
					name = w.Cities[se.Destination].Name
				}
				t.Logf("編成：faction=%d target=%d corps=%d destination=%d %q", se.Faction, se.Target, se.Corps, se.Destination, name)
				if se.Target == 0 {
					from := w.Factions[0].Capital
					rows := make([]int, 0, len(w.Cities)-1)
					for city := range w.Cities {
						if city != from {
							rows = append(rows, city)
						}
					}
					sort.SliceStable(rows, func(i, j int) bool {
						di := chebyshev(w.Cities[from], w.Cities[rows[i]])
						dj := chebyshev(w.Cities[from], w.Cities[rows[j]])
						return di < dj
					})
					for pos, city := range rows {
						if city == se.Destination {
							t.Logf("玩家從據點 %d 選距離排序第 %d 列到敵方目的地 %d", from, pos, city)
							break
						}
					}
				}
			}
		}
		for _, ce := range ev.Corps {
			if ce.Battle != nil {
				battles++
			}
		}
		if ev.Settled {
			months++
		}
		if v := w.CheckInvariants(); len(v) > 0 {
			t.Fatalf("第 %d 個月 tick 不變量違反：%s", months, v[0])
		}
	}
	if declarations == 0 {
		t.Fatal("六個月內沒有任何敵方宣戰")
	}
	if formed == 0 {
		t.Fatal("宣戰後沒有敵方編成軍團")
	}
	if battles == 0 {
		t.Fatal("敵方軍團六個月內沒有進入戰鬥")
	}
	t.Logf("6 個月：宣戰 %d、編成 %d、戰鬥 %d、活著軍團 %d", declarations, formed, battles, len(w.AliveCorps()))
}

func chebyshev(a, b City) int {
	dx, dy := a.X-b.X, a.Y-b.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	if dx > dy {
		return dx
	}
	return dy
}

// 威脅偵測要在據點輪到時算出來（原版 sub_13EFD → sub_13FA9）。
// 這一條用的是原版劇本，所以鄰接關係與交友度都是真的。
func TestCityThreatIsRecomputedOnTick(t *testing.T) {
	w := load(t, 0)
	// 找一個「鄰居裡有交戰勢力」的據點。開局交友度全是和平，
	// 所以要自己把一組關係打成交戰，否則量到的一定是 0——
	// **零值也可能只是前提沒滿足**，不是實作沒接上。
	site, foe := -1, -1
	for i := range w.Cities {
		c := &w.Cities[i]
		if c.Owner < 0 || c.Owner >= numFactions || c.Owner == threatNeutral {
			continue
		}
		for _, n := range c.Neighbours {
			if n < 0 || n >= len(w.Cities) {
				continue
			}
			o := w.Cities[n].Owner
			if o != c.Owner && o >= 0 && o < numFactions {
				site, foe = i, o
				break
			}
		}
		if site >= 0 {
			break
		}
	}
	if site < 0 {
		t.Skip("這個劇本沒有相鄰的敵方據點")
	}
	owner := w.Cities[site].Owner
	w.Friendship[owner][foe] = diplomacy.Friendship(50)

	// 在敵方鄰居那一格擺一支軍團，威脅量才會是非 0。
	var nb int
	for _, n := range w.Cities[site].Neighbours {
		if n >= 0 && n < len(w.Cities) && w.Cities[n].Owner == foe {
			nb = n
			break
		}
	}
	w.Corps[0] = Corps{Alive: true, Faction: foe,
		X: w.Cities[nb].X, Y: w.Cities[nb].Y, Node: nb}

	r := rng.NewFixed(1)
	w.refreshCityThreat(nb, r)  // 先算鄰居的佔用數，威脅量才有東西可加
	w.refreshCityThreat(site, r)

	if got := w.Cities[nb].Occupancy; got != 1 {
		t.Fatalf("敵方據點的佔用數 = %d，want 1", got)
	}
	if w.Cities[site].EnemyNeighbours == 0 {
		t.Fatal("相鄰敵據點數是 0，但剛剛才找到一個敵方鄰居")
	}
	if got := w.Cities[site].Threat; got != 1 {
		t.Fatalf("周邊威脅量 = %d，want 1", got)
	}
}

// 軍團上限是 max(5, 資金 ÷ 8192) 減掉現有軍團數，不是「只能有一支」。
func TestAICorpsCapFollowsFunds(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	w.strategicAI = true
	faction := -1
	for i := 1; i < numFactions; i++ {
		if w.Factions[i].Alive {
			faction = i
			break
		}
	}
	if faction < 0 {
		t.Skip("這個劇本只有一個勢力")
	}
	f := &w.Factions[faction]
	f.InvasionTarget = 0
	f.Funds = 0
	for k := range f.Reserves {
		f.Reserves[k] = 600
	}

	f.Corps = 1
	if ev := w.formAICorps(faction); ev == nil {
		t.Fatal("已有一支軍團就不再編成——上限應該是 5")
	}
	f.Corps = 5
	if ev := w.formAICorps(faction); ev != nil {
		t.Fatal("資金 0 時上限是 5，第六支不該編得出來")
	}
	f.Funds = 81_920 // 上限 10
	if ev := w.formAICorps(faction); ev == nil {
		t.Fatal("資金 81,920 時上限是 10，第六支應該編得出來")
	}
}

// 玩家的據點被侵攻目標貼著、而且那一格沒有軍團時，會跳訊息 #38。
// 再跑一次不會重複跳——冷卻計時器擋著。
func TestPlayerCityAsksForRelief(t *testing.T) {
	w := load(t, 0)
	site, foe := -1, -1
	for i := range w.Cities {
		c := &w.Cities[i]
		if c.Owner < 0 || c.Owner >= numFactions {
			continue
		}
		for _, n := range c.Neighbours {
			if n < 0 || n >= len(w.Cities) {
				continue
			}
			if o := w.Cities[n].Owner; o != c.Owner && o >= 0 && o < numFactions {
				site, foe = i, o
				break
			}
		}
		if site >= 0 {
			break
		}
	}
	if site < 0 {
		t.Skip("這個劇本沒有相鄰的敵方據點")
	}
	w.Player = w.Cities[site].Owner
	w.Factions[w.Player].InvasionTarget = foe
	w.Friendship[w.Player][foe] = diplomacy.Friendship(50) // 交戰
	w.Cities[site].ReliefCooldown = 0
	for i := range w.Corps { // 這一格不能有軍團，否則走的是另一條路
		w.Corps[i].Alive = false
	}

	r := rng.NewFixed(3)
	got := w.refreshCityThreat(site, r)
	if len(got) != 1 || got[0].Index != 38 || got[0].City != site {
		t.Fatalf("求援訊息 ＝ %+v，期望一則 #38 指向據點 %d", got, site)
	}
	if w.Cities[site].ReliefCooldown < 24 || w.Cities[site].ReliefCooldown > 39 {
		t.Fatalf("玩家的冷卻 ＝ %d，期望 24–39", w.Cities[site].ReliefCooldown)
	}
	if w.Factions[w.Player].ReliefSite != site {
		t.Errorf("勢力 +0x16 ＝ %d，期望 %d", w.Factions[w.Player].ReliefSite, site)
	}
	if again := w.refreshCityThreat(site, r); len(again) != 0 {
		t.Fatalf("冷卻中還是又求援了：%+v", again)
	}
}

// 求援只調得動「委任」中的軍團（原版 sub_14155 的 test byte [di],4）。
// 玩家自己指揮的、以及等著解體的，都不該被調走。
func TestReliefOnlyMovesDelegatedCorps(t *testing.T) {
	w := load(t, 0)
	site := 0
	owner := w.Cities[site].Owner
	if owner < 0 || owner >= numFactions {
		t.Skip("第一個據點沒有勢力")
	}
	target := -1
	for i := range w.Cities {
		if i != site && w.Cities[i].Owner != owner {
			target = i
			break
		}
	}
	if target < 0 {
		t.Skip("找不到別的勢力的據點")
	}
	for i := range w.Corps {
		w.Corps[i] = Corps{}
	}
	at := func(i int, delegated bool, stage int) {
		w.Corps[i] = Corps{Alive: true, Faction: owner, Node: site,
			X: w.Cities[site].X, Y: w.Cities[site].Y,
			TargetNode: site, Delegated: delegated, Stage: stage}
	}
	at(0, false, 0) // 玩家自己指揮
	at(1, true, 11) // 待解體
	at(2, true, 0)  // 委任中 ✓

	w.dispatchGarrison(site, target, 3, 0, rng.NewFixed(1))

	if w.Corps[0].TargetNode != site {
		t.Error("玩家自己指揮的軍團被調走了")
	}
	if w.Corps[1].TargetNode != site {
		t.Error("等著解體的軍團被調走了")
	}
	if w.Corps[2].TargetNode == site {
		t.Error("委任中的軍團沒有被調走")
	}
}
