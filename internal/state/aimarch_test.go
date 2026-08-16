package state

import (
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
)


import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
)

// aiCorps 編一支非玩家的軍團，並把它擺在指定據點上（已經抵達的狀態）。
func aiCorps(t *testing.T, w *World, faction, node int) int {
	t.Helper()
	// 玩家設成別人，這一支才會走 AI 那一半的分派。
	for _, f := range w.AliveFactions() {
		if f != faction {
			w.Player = f
			break
		}
	}
	i := formOne(t, w, faction)
	c := &w.Corps[i]
	c.Node, c.TargetNode, c.Ordered = node, node, node
	c.X, c.Y = w.Cities[node].X, w.Cities[node].Y
	c.TargetX, c.TargetY = c.X, c.Y
	return i
}

// Stage 0：目標據點的「威脅有具體目標」旗標亮著就原地駐守（sub_1439D）。
func TestAIStage0StaysWhileThreatIsSpecific(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	w.Corps[i].Stage = StageNormal

	w.Cities[node].Specific = true
	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Stage != StageNormal {
		t.Errorf("旗標亮著就不該換檔，Stage = %d", w.Corps[i].Stage)
	}

	w.Cities[node].Specific = false
	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Stage != 1 {
		t.Errorf("旗標滅了應該轉 Stage 1，得到 %d", w.Corps[i].Stage)
	}
}

// Stage 1：兵力 ≤ 300 點就放下手邊的事回首都補兵（sub_143AF 的 12Ch）。
func TestAIStage1SendsWeakCorpsHome(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	c := &w.Corps[i]
	c.Stage = 1
	c.Men = aiHomeThreshold

	w.arriveCorps(i, &testRand{})
	if c.Stage != StageHomeResupply {
		t.Fatalf("殘兵應該轉 Stage %d，得到 %d", StageHomeResupply, c.Stage)
	}
	// 已經在首都上 → 下一次分派直接轉補兵。
	w.arriveCorps(i, &testRand{})
	if c.Stage != StageResupply {
		t.Errorf("人在首都就該轉 Stage %d，得到 %d", StageResupply, c.Stage)
	}
}

// Stage 2：先領失土（+0x17）再領求援（+0x16），兩格都是取走就清空。
func TestAIStage2TakesLostSiteFirst(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	w.Corps[i].Stage = 2
	w.Cities[node].Threatened, w.Cities[node].Specific = false, false

	lost, relief := pickTwoCities(w, node)
	w.Factions[f].LostSite, w.Factions[f].ReliefSite = lost, relief

	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Ordered != lost {
		t.Errorf("第一次該領走失土 %d，卻去了 %d", lost, w.Corps[i].Ordered)
	}
	if w.Factions[f].LostSite != noSite {
		t.Errorf("領走之後失土那一格要清空，得到 %d", w.Factions[f].LostSite)
	}
	if w.Factions[f].ReliefSite != relief {
		t.Error("求援那一格不該一起被清掉")
	}
	if w.Corps[i].Stage != StageNormal {
		t.Errorf("領到目標就該轉 Stage 0，得到 %d", w.Corps[i].Stage)
	}

	// 第二次只剩求援。
	w.Corps[i].Stage = 2
	w.Corps[i].TargetNode = node
	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Ordered != relief {
		t.Errorf("第二次該領走求援 %d，卻去了 %d", relief, w.Corps[i].Ordered)
	}
	if w.Factions[f].ReliefSite != noSite {
		t.Errorf("求援那一格要清空，得到 %d", w.Factions[f].ReliefSite)
	}
}

// 資金吃緊（勢力 +0x00 位元 6）的勢力跳過失土，只接求援。
func TestAIStage2SkipsLostSiteWhenLowFunds(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	w.Corps[i].Stage = 2
	w.Cities[node].Threatened, w.Cities[node].Specific = false, false

	lost, relief := pickTwoCities(w, node)
	w.Factions[f].LostSite, w.Factions[f].ReliefSite = lost, relief
	w.Factions[f].LowFunds = true

	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Ordered != relief {
		t.Errorf("沒錢時該只接求援 %d，卻去了 %d", relief, w.Corps[i].Ordered)
	}
	if w.Factions[f].LostSite != lost {
		t.Error("沒錢時失土那一格不該被領走")
	}
}

// Stage 2 領不到待辦、這一格又不受威脅 → 解體回收兵力。
func TestAIStage2DisbandsWithNothingToDo(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	w.Corps[i].Stage = 2
	w.Cities[node].Threatened, w.Cities[node].Specific = false, false
	w.Factions[f].LostSite, w.Factions[f].ReliefSite = noSite, noSite

	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Stage != StageDisband {
		t.Fatalf("沒事做該轉 Stage %d，得到 %d", StageDisband, w.Corps[i].Stage)
	}

	// 受威脅時反而要留著。這一格要站得下第二支，否則會先走「留守」那條
	// （受威脅且只有自己一支 → 退回 Stage 1）。
	w.Corps[i].Stage = 2
	w.Cities[node].Threatened = true
	w.Cities[node].Occupancy = aiHoldAlone + 1
	w.arriveCorps(i, &testRand{})
	if w.Corps[i].Stage != 2 {
		t.Errorf("受威脅時該留在 Stage 2，得到 %d", w.Corps[i].Stage)
	}
}

// Stage 3：六槽任一不到 30 點就整團解散，否則去等士氣（sub_14466）。
func TestAIStage3DisbandsUnderstrengthCorps(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	c := &w.Corps[i]
	for k := range c.Units {
		if c.Units[k].Men < aiSlotThreshold {
			t.Skipf("這一支第 %d 槽只有 %d 點，測不到滿額那條路", k, c.Units[k].Men)
		}
	}

	c.Stage = StageDone
	w.arriveCorps(i, &testRand{})
	if c.Stage != StageWaitMorale {
		t.Fatalf("六槽都夠該轉 Stage %d，得到 %d", StageWaitMorale, c.Stage)
	}

	c.Stage = StageDone
	c.Men -= c.Units[2].Men - (aiSlotThreshold - 1)
	c.Units[2].Men = aiSlotThreshold - 1
	w.arriveCorps(i, &testRand{})
	if c.Stage != StageDisband {
		t.Errorf("有一槽不足額該轉 Stage %d，得到 %d", StageDisband, c.Stage)
	}
}

// Stage 8：士氣回到勢力基準才繼續行動（sub_14483）。
func TestAIStage8WaitsForMorale(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	c := &w.Corps[i]
	c.Stage = StageWaitMorale
	c.Morale = w.Factions[f].MoraleBase - 1

	w.arriveCorps(i, &testRand{})
	if c.Stage != StageWaitMorale {
		t.Fatalf("士氣沒回滿不該動，Stage = %d", c.Stage)
	}
	c.Morale = w.Factions[f].MoraleBase
	w.arriveCorps(i, &testRand{})
	if c.Stage != 1 {
		t.Errorf("士氣達標該回 Stage 1，得到 %d", c.Stage)
	}
}

// 整條鏈：站在不受威脅的據點上 → Stage 1 → Stage 2 → 領走求援並出發。
func TestAICorpsLifecycleReachesRelief(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	c := &w.Corps[i]
	c.Stage = StageNormal
	w.Cities[node].Threatened, w.Cities[node].Specific = false, false
	w.Cities[node].Occupancy = 0

	_, relief := pickTwoCities(w, node)
	w.Factions[f].LostSite = noSite
	w.Factions[f].ReliefSite = relief

	for n := 0; n < 4 && c.Ordered != relief; n++ {
		w.arriveCorps(i, &testRand{})
		c.TargetNode = c.Node // 還沒真的走，維持「停在原地」
	}
	if c.Ordered != relief {
		t.Fatalf("跑完四次分派還沒去救援，Ordered = %d、Stage = %d",
			c.Ordered, c.Stage)
	}
}

// 勢力滅亡時不會留下沒有主人的軍團（sub_14FCE）。
func TestEliminatedFactionLeavesNoCorps(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)

	w.eliminateFaction(f, noFaction)
	if w.Corps[i].Alive {
		t.Error("滅亡的勢力還留著軍團")
	}
	if w.Factions[f].Corps != 0 {
		t.Errorf("滅亡的勢力軍團數 = %d", w.Factions[f].Corps)
	}
	if w.Generals[i].Posted {
		t.Error("軍團沒了，主將還標著出陣")
	}
	// 這裡不查全域不變量：測試是直接呼叫滅亡入口，該勢力的據點還在，
	// 而原版只有在「最後一個據點也沒了」時才走到這條路。
	if w.Factions[f].Generals != 0 {
		t.Errorf("滅亡的勢力武將數 = %d", w.Factions[f].Generals)
	}
}

// pickTwoCities 挑兩個不等於 skip 的據點編號當測試目標。
func pickTwoCities(w *World, skip int) (int, int) {
	var got []int
	for i := range w.Cities {
		if i == skip || w.Cities[i].Owner == diplomacy.NoTarget {
			continue
		}
		got = append(got, i)
		if len(got) == 2 {
			break
		}
	}
	return got[0], got[1]
}

// 敗走：倒數歸零之後軍團才真的消失，主將回到無職（docs/spec/43）。
func TestRoutedCorpsDisappearsAfterTimer(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	before := w.Factions[f].Corps

	w.routCorps(i)
	if w.Corps[i].Alive {
		t.Error("敗走中的軍團不該算活著")
	}
	if !w.Corps[i].Routing || w.Corps[i].RoutTimer != routDuration {
		t.Fatalf("敗走狀態 = %v／%d，want true／%d",
			w.Corps[i].Routing, w.Corps[i].RoutTimer, routDuration)
	}
	if w.Factions[f].Corps != before-1 {
		t.Errorf("勢力軍團數 = %d，want %d", w.Factions[f].Corps, before-1)
	}
	if !w.Generals[i].Posted {
		t.Error("倒數還沒歸零，主將不該解職")
	}

	for n := 0; n < routDuration; n++ {
		w.tickRout(i)
	}
	if w.Corps[i].Routing {
		t.Error("倒數歸零之後還在敗走狀態")
	}
	if w.Generals[i].Posted {
		t.Error("倒數歸零之後主將沒解職")
	}
}

// 敗走**不回收兵員**——這是它與解體最大的差別。
func TestRoutLosesTheMen(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	pool := w.Factions[f].Reserves

	w.routCorps(i)
	if w.Factions[f].Reserves != pool {
		t.Errorf("敗走之後預備兵變成 %v，原本 %v——兵不該回池",
			w.Factions[f].Reserves, pool)
	}
	if w.Corps[i].Men == 0 {
		t.Error("敗走不該把六槽清空，原版只改旗標與倒數")
	}
}

// 旗標 8 與倒數要寫得回存檔（原版 `sub_12977` 只改這兩個 byte）。
func TestRoutingSurvivesSaveRoundTrip(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)
	w.routCorps(i)
	w.Corps[i].RoutTimer = 17

	raw := make([]byte, corpsBase+numCorps*corpsSize)
	w.saveCorps(raw)
	r := raw[corpsBase+i*corpsSize:]
	if r[0x00] != 0x08 || r[0x03] != 17 {
		t.Fatalf("寫回去的旗標／倒數 = %#x／%d，want 0x08／17", r[0x00], r[0x03])
	}

	var back World
	back.loadCorps(raw)
	if !back.Corps[i].Routing || back.Corps[i].RoutTimer != 17 {
		t.Errorf("讀回來 = %v／%d，want true／17",
			back.Corps[i].Routing, back.Corps[i].RoutTimer)
	}
	if back.Corps[i].Alive {
		t.Error("旗標 8 不到 0x80，不該算活著")
	}
}

// 回家的路穿過別人的地 → 敗走；全是自己的地 → 不敗走（docs/spec/43 §1）。
//
// ⚠ 這一條**只有掛了道路圖才會生效**：`returnBlocked` 走的是
// `w.routes[i]`，而 `wlsim` 沒有 `SetRoads`，路徑永遠是空的。
func TestReturnBlockedNeedsForeignCityOnTheRoute(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	node := w.clampCity(w.Factions[f].Capital)
	i := aiCorps(t, w, f, node)

	// 自己的地就用首都；另外找一個別人的據點。
	own, foreign := node, -1
	for n := range w.Cities {
		if w.Cities[n].Owner != f {
			foreign = n
			break
		}
	}
	if foreign < 0 {
		t.Skip("這個劇本裡沒有別人的據點")
	}

	w.routes[i] = [][2]int{{w.Cities[own].X, w.Cities[own].Y}}
	if w.returnBlocked(i) {
		t.Error("整條路都是自己的地，不該算被擋")
	}
	w.routes[i] = [][2]int{
		{w.Cities[own].X, w.Cities[own].Y},
		{w.Cities[foreign].X, w.Cities[foreign].Y},
	}
	if !w.returnBlocked(i) {
		t.Fatal("路上有別人的據點，應該算被擋")
	}

	w.routIfBlocked(i)
	if !w.Corps[i].Routing {
		t.Error("被擋住之後沒有進敗走狀態")
	}
}

// ---------------------------------------------------------------------------
// 戰後退一站回家（docs/spec/46）
// ---------------------------------------------------------------------------

// threeInARow 建一條 A—B—C 的線形道路圖，回傳三個據點編號。
// A 當首都，C 是軍團現在站的地方。
func threeInARow(t *testing.T, w *World, faction int) (a, b, c int) {
	t.Helper()
	picked := []int{}
	for n := range w.Cities {
		if len(picked) == 3 {
			break
		}
		picked = append(picked, n)
	}
	if len(picked) < 3 {
		t.Skip("據點不夠三個")
	}
	a, b, c = picked[0], picked[1], picked[2]
	for _, n := range picked {
		w.Cities[n].Owner = faction
	}
	w.Factions[faction].Capital = a
	w.SetRoads(march.New(len(w.Cities), []march.Edge{
		{A: a, B: b, Steps: 1}, {A: b, B: c, Steps: 1},
	}))
	return a, b, c
}


// standOnForeignGround 把軍團腳下那一格換成別人的，
// 好讓它走「不在自家據點上」那一支。
func standOnForeignGround(w *World, node, faction int) {
	for other := range w.Factions {
		if other != faction {
			w.Cities[node].Owner = other
			return
		}
	}
}

// `sub_1487B` 的第二個歸屬檢查：下一站不是自己的地就退不了。
func TestNextHopHomeStopsAtForeignGround(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	a, b, c := threeInARow(t, w, f)
	i := aiCorps(t, w, f, c)

	if got := w.nextHopHome(i); got != b {
		t.Fatalf("下一站 = %d，要 %d（首都是 %d）", got, b, a)
	}
	// 中間那一站換人 ⇒ 退不了。
	for other := range w.Factions {
		if other != f {
			w.Cities[b].Owner = other
			break
		}
	}
	if got := w.nextHopHome(i); got != -1 {
		t.Errorf("下一站是別人的地，應該回 −1，實際 %d", got)
	}
}

// `cmp al, 0FFh`：沒有首都就無處可退。
func TestNextHopHomeWithoutCapital(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	_, _, c := threeInARow(t, w, f)
	i := aiCorps(t, w, f, c)
	w.Factions[f].Capital = noCity
	if got := w.nextHopHome(i); got != -1 {
		t.Errorf("沒有首都應該回 −1，實際 %d", got)
	}
}

// 敗方退**一站**，不是直接回首都。
func TestLoserRetreatsOneHop(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	a, b, c := threeInARow(t, w, f)
	i := aiCorps(t, w, f, c)
	standOnForeignGround(w, c, f)
	w.Corps[i].Men = army60000Points // 兵力夠 ⇒ 不轉 Stage 10

	if dead := w.retreatOrPerish(i, false); dead {
		t.Fatal("有退路卻判成壞滅")
	}
	if got := w.Corps[i].Ordered; got != b {
		t.Errorf("退到 %d，要退一站到 %d（首都在 %d）", got, b, a)
	}
	if got := w.Corps[i].Stage; got != StageWaitMorale {
		t.Errorf("Stage = %d，兵力還夠時要 %d", got, StageWaitMorale)
	}
}

// 兵力 ≤ 300 就轉 Stage 10（`cmp word [si+4], 12Ch`）。
func TestWeakLoserSwitchesToHomeResupply(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	_, _, c := threeInARow(t, w, f)
	i := aiCorps(t, w, f, c)
	standOnForeignGround(w, c, f)
	w.Corps[i].Men = aiHomeThreshold

	if dead := w.retreatOrPerish(i, false); dead {
		t.Fatal("有退路卻判成壞滅")
	}
	if got := w.Corps[i].Stage; got != StageHomeResupply {
		t.Errorf("Stage = %d，兵力 ≤ %d 時要 %d", got, aiHomeThreshold, StageHomeResupply)
	}
}

// 勝方（原版的 `cl == 0`）原地不動。
func TestWinnerStandsStill(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	_, _, c := threeInARow(t, w, f)
	i := aiCorps(t, w, f, c)

	if dead := w.retreatOrPerish(i, true); dead {
		t.Fatal("勝方被判成壞滅")
	}
	if got := w.Corps[i].Ordered; got != c {
		t.Errorf("勝方動了：Ordered = %d，要 %d", got, c)
	}
	if got := w.Corps[i].Stage; got != StageWaitMorale {
		t.Errorf("Stage = %d，要 %d", got, StageWaitMorale)
	}
}

// 站在自家據點上就不退——攻城時守方走的正是這一支。
func TestDefenderInOwnCityDoesNotRetreat(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	_, _, c := threeInARow(t, w, f)
	i := aiCorps(t, w, f, c) // c 已經是自己的地

	if dead := w.retreatOrPerish(i, false); dead {
		t.Fatal("守在自家城裡卻被判成壞滅")
	}
	if got := w.Corps[i].Ordered; got != c {
		t.Errorf("守方動了：Ordered = %d，要 %d", got, c)
	}
}

// ⭐ 壞滅的第三個入口：退不了。
func TestNoRetreatMeansDestroyed(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	_, b, c := threeInARow(t, w, f)
	i := aiCorps(t, w, f, c)
	// 現在站的那一格換成別人的（才不會走「站在自家城裡」那一支），
	// 回家的下一站也換成別人的。
	for other := range w.Factions {
		if other != f {
			w.Cities[c].Owner, w.Cities[b].Owner = other, other
			break
		}
	}
	if dead := w.retreatOrPerish(i, false); !dead {
		t.Fatal("退不了卻沒判成壞滅")
	}
}

// ---------------------------------------------------------------------------
// 據點易主之後的調頭（docs/spec/47）
// ---------------------------------------------------------------------------

// stackedDefender 在 node 上放一支 faction 的軍團（不參戰，只是疊在那裡）。
func stackedDefender(t *testing.T, w *World, faction, node int) int {
	t.Helper()
	i := formOne(t, w, faction)
	c := &w.Corps[i]
	c.Node, c.TargetNode, c.Ordered = node, node, node
	c.X, c.Y = w.Cities[node].X, w.Cities[node].Y
	c.TargetX, c.TargetY = c.X, c.Y
	return i
}

// 疊在那一格上、沒被捲進戰鬥的守軍，易主之後退一站回家。
func TestFallenCityCorpsRetreatOneHop(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	a, b, c := threeInARow(t, w, f)
	i := stackedDefender(t, w, f, c)

	attacker := -1
	for _, g := range w.AliveFactions() {
		if g != f {
			attacker = g
			break
		}
	}
	att := stackedDefender(t, w, attacker, c)

	ev := &CorpsEvent{Captured: -1}
	w.capture(att, ev, &testRand{s: 1})

	if w.Cities[c].Owner != attacker {
		t.Fatalf("據點沒換手：%d", w.Cities[c].Owner)
	}
	if !w.Corps[i].Alive {
		t.Fatal("有退路的守軍不該消失")
	}
	if got := w.Corps[i].Ordered; got != b {
		t.Errorf("調頭到 %d，要退一站到 %d（首都在 %d）", got, b, a)
	}
	if got := w.Corps[i].Timer; got != 1 {
		t.Errorf("計時器 = %d，原版寫 1（下一個 tick 起步）", got)
	}
}

// 退不了就走壞滅同一個出口（原版 `sub_1291A`）。
func TestFallenCityCorpsWithNoRetreatPerish(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	_, b, c := threeInARow(t, w, f)
	i := stackedDefender(t, w, f, c)

	attacker := -1
	for _, g := range w.AliveFactions() {
		if g != f {
			attacker = g
			break
		}
	}
	w.Cities[b].Owner = attacker // 回家的下一站是別人的地
	att := stackedDefender(t, w, attacker, c)

	before := w.Factions[f].Corps
	ev := &CorpsEvent{Captured: -1}
	w.capture(att, ev, &testRand{s: 1})

	if w.Corps[i].Alive {
		t.Fatal("退不了的守軍還在")
	}
	if w.Factions[f].Corps != before-1 {
		t.Errorf("勢力軍團數 %d，要 %d", w.Factions[f].Corps, before-1)
	}
	found := false
	for _, d := range ev.Destroyed {
		if d == i {
			found = true
		}
	}
	if !found {
		t.Error("沒有進 ev.Destroyed，事件層看不到這支軍團沒了")
	}
}

// 首都被打下來時，那一格上的守軍最後是朝**新首都**走。
//
// ⚠ **這一條驗的是結果，不是順序。** 原版 `sub_14CF3` 先遷都
// （`sub_14DF0`）再調頭（`sub_14DA4`），remake 照抄了這個順序；
// 但把兩者對調也會得到同一個答案——遷都會呼叫 `sub_14502`
// （`syncCorpsAfterCapitalChange`），把「目標還是舊首都」的軍團
// 一律改掛新首都。**兩條路在這裡收斂**，所以測試擋不住順序寫反。
// 實際跑過負對照確認：把調頭提前，這一條照樣綠。
func TestFallenCapitalRedirectsTowardTheNewCapital(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	picked := []int{}
	for n := range w.Cities {
		if len(picked) == 3 {
			break
		}
		picked = append(picked, n)
	}
	if len(picked) < 3 {
		t.Skip("據點不夠三個")
	}
	a, b, c := picked[0], picked[1], picked[2]
	for _, n := range picked {
		w.Cities[n].Owner = f
	}
	w.Factions[f].Capital = c // ★ 首都就是等一下要陷落的那一格
	w.SetRoads(march.New(len(w.Cities), []march.Edge{
		{A: c, B: a, Steps: 1}, {A: c, B: b, Steps: 1},
	}))
	i := stackedDefender(t, w, f, c)

	attacker := -1
	for _, g := range w.AliveFactions() {
		if g != f {
			attacker = g
			break
		}
	}
	att := stackedDefender(t, w, attacker, c)

	ev := &CorpsEvent{Captured: -1}
	w.capture(att, ev, &testRand{s: 1})

	newCapital := w.Factions[f].Capital
	if newCapital == c || newCapital == noCity {
		t.Fatalf("首都被打下來卻沒遷都：%d", newCapital)
	}
	if !w.Corps[i].Alive {
		t.Fatal("遷都之後應該還有地方可去")
	}
	if got := w.Corps[i].Ordered; got != newCapital {
		t.Errorf("調頭到 %d，新首都是 %d（星形拓樸下下一站就是它）", got, newCapital)
	}
}

// ---------------------------------------------------------------------------
// 內政官被遣回（docs/spec/48）
// ---------------------------------------------------------------------------

// 據點被攻陷 → 內政官槽清空、那名武將不再「出陣中」、事件帶出編號。
func TestCityFallReturnsGovernor(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	_, _, node := threeInARow(t, w, f)

	gov := -1
	for i, g := range w.Generals {
		if g.Alive && g.Faction == f {
			gov = i
			break
		}
	}
	if gov < 0 {
		t.Skip("這個勢力沒有武將")
	}
	w.Cities[node].Governor = gov
	w.Generals[gov].Posted = true

	attacker := -1
	for _, g := range w.AliveFactions() {
		if g != f {
			attacker = g
			break
		}
	}
	att := stackedDefender(t, w, attacker, node)

	ev := &CorpsEvent{Captured: -1, GovernorReturned: noGovernor}
	w.capture(att, ev, &testRand{s: 1})

	if got := w.Cities[node].Governor; got != noGovernorSlot {
		t.Errorf("內政官槽 = %d，要 0xFF", got)
	}
	if w.Generals[gov].Posted {
		t.Error("被遣回的內政官還掛著「出陣中」")
	}
	if ev.GovernorReturned != gov {
		t.Errorf("事件帶出 %d，要 %d", ev.GovernorReturned, gov)
	}
}

// 無主的據點被佔下時整段跳過（原版 `cmp bh, 18h / jz`）。
func TestNeutralCityFallKeepsNoGovernor(t *testing.T) {
	w := load(t, 0)
	f := w.AliveFactions()[1]
	_, _, node := threeInARow(t, w, f)
	w.Cities[node].Owner = combat.NeutralFaction
	w.Cities[node].Governor = 3

	att := stackedDefender(t, w, f, node)
	ev := &CorpsEvent{Captured: -1, GovernorReturned: noGovernor}
	w.capture(att, ev, &testRand{s: 1})

	if ev.GovernorReturned != noGovernor {
		t.Errorf("無主的據點不該跑遣回那一段，卻回了 %d", ev.GovernorReturned)
	}
	if got := w.Cities[node].Governor; got != 3 {
		t.Errorf("內政官槽被動了：%d", got)
	}
}
