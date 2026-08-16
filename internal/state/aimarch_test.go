package state

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
