package state

import (
	"sort"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/strategyai"
)

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
