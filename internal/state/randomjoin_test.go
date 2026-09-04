package state

import "testing"

// fixedRand 讓每一次 Next() 回固定值，用來釘住骰面對應的分支。
type fixedRand struct{ v int }

func (f fixedRand) Next() int { return f.v }

// 骰面 1 與 24–47 都投靠**武將最少**的勢力（docs/spec/130 §2.1、§2.3）。
func TestRandomJoinPrefersFewestGenerals(t *testing.T) {
	for _, roll := range []int{0, 0x17, 0x2E} { // +1 之後 ＝ 1、24、47
		w := &World{Player: 5}
		for i := 0; i < numFactions; i++ {
			w.Factions[i].Alive = true
			w.Factions[i].Generals = 10
		}
		w.Factions[7].Generals = 2 // 最少
		w.Generals[3] = General{Alive: true, Faction: noFaction, Affinity: noFaction}

		got := w.randomJoin(3, fixedRand{roll})
		if got != 7 {
			t.Errorf("骰面 %d → 勢力 %d，want 7（武將最少）", roll+1, got)
		}
		if w.Generals[3].Faction != 7 || w.Factions[7].Generals != 3 {
			t.Errorf("沒有寫回：Faction=%d、武將數=%d",
				w.Generals[3].Faction, w.Factions[7].Generals)
		}
	}
}

// 骰面 n（1–23）從最少那一個往後數第 n 個**存在**的勢力，會繞回。
func TestRandomJoinCountsFromFewest(t *testing.T) {
	w := &World{Player: 21}
	for i := 0; i < numFactions; i++ {
		w.Factions[i].Alive = false
	}
	// 只有 3、9、20 活著；9 的武將最少 ⇒ 起點是 9。
	for _, i := range []int{3, 9, 20} {
		w.Factions[i].Alive = true
		w.Factions[i].Generals = 10
	}
	w.Factions[9].Generals = 1
	w.Generals[4] = General{Alive: true, Faction: noFaction, Affinity: noFaction}

	// 骰面 3 ⇒ 從 9 起數三個存在的：9(1)、20(2)、3(3，繞回) ⇒ 勢力 3。
	if got := w.randomJoin(4, fixedRand{2}); got != 3 {
		t.Errorf("骰面 3 → 勢力 %d，want 3（繞回）", got)
	}
}

// 骰面 ≥ 48 走玩家救濟：據點數 ÷ 4 ＋ 1 > 武將數才送人。
func TestRandomJoinPlayerReliefThreshold(t *testing.T) {
	mk := func(cities, generals int) *World {
		w := &World{Player: 2}
		for i := 0; i < numFactions; i++ {
			w.Factions[i].Alive = true
			w.Factions[i].Generals = 10
		}
		w.Factions[2].Cities, w.Factions[2].Generals = cities, generals
		w.Generals[6] = General{Alive: true, Faction: noFaction, Affinity: noFaction}
		return w
	}
	// 20/4+1 = 6 > 3 ⇒ 送。
	if got := mk(20, 3).randomJoin(6, fixedRand{0x2F}); got != 2 {
		t.Errorf("缺人時 → %d，want 2（玩家）", got)
	}
	// 20/4+1 = 6 ≤ 6 ⇒ 不動。
	if got := mk(20, 6).randomJoin(6, fixedRand{0x2F}); got != -1 {
		t.Errorf("不缺人時 → %d，want −1（什麼都不做）", got)
	}
}

// ⭐ 接線：沒有心向的武將**每月都跑**，不受 25% 閘限制
// （原版 `cmp bh, 0FFh / jz` 在擲那顆骰之前，docs/spec/130）。
func TestRandomJoinSkipsTwentyFivePercentGate(t *testing.T) {
	w := &World{Player: 5}
	for i := 0; i < numFactions; i++ {
		w.Factions[i].Alive = true
		w.Factions[i].Generals = 10
	}
	// 玩家缺人：40/4+1 = 11 > 10 ⇒ 骰面 64 那一條會送人過來。
	w.Factions[5].Cities = 40
	w.Generals[3] = General{Alive: true, Faction: noFaction, Affinity: noFaction}

	// ⭐ 0xFF 遠大於 affinityRollLimit（0x40）——**有心向的話這一輪不會動**，
	// 而這一位沒有心向，所以照跑。骰面 (0xFF&0x3F)+1 = 64 ⇒ 玩家救濟。
	joined := w.recruitFreelanceGenerals(fixedRand{0xFF})
	if len(joined) != 1 || joined[0] != 3 {
		t.Fatalf("出仕清單 %v，want [3]——隨機投靠不該被 25%% 閘擋住", joined)
	}
	if w.Generals[3].Faction != 5 {
		t.Errorf("投靠的是勢力 %d，want 5（玩家救濟）", w.Generals[3].Faction)
	}

	// 正對照：同一顆骰、同一個人，**有心向**就會被 25% 閘擋下來。
	w2 := &World{Player: 5}
	for i := 0; i < numFactions; i++ {
		w2.Factions[i].Alive = true
		w2.Factions[i].Generals = 10
	}
	w2.Generals[3] = General{Alive: true, Faction: noFaction, Affinity: 7}
	if got := w2.recruitFreelanceGenerals(fixedRand{0xFF}); len(got) != 0 {
		t.Errorf("有心向的在 0xFF 這一輪出仕了 %v，25%% 閘沒作用", got)
	}
}

// ⭐ 接線：投靠玩家時要排一則 TALK #41「{1}加入麾下了。」
// （`sub_15899` 的 `loc_1591A`，兩條路共用），投靠別人時**不排**。
func TestFreelanceJoinNotifiesOnlyPlayer(t *testing.T) {
	mk := func(affinity, player int) *World {
		w := &World{Player: player}
		for i := 0; i < numFactions; i++ {
			w.Factions[i].Alive = true
			w.Factions[i].Generals = 10
		}
		w.Generals[3] = General{Alive: true, Faction: noFaction, Affinity: affinity}
		return w
	}
	notices := func(w *World) []TalkNotice {
		var ev Event
		for _, id := range w.recruitFreelanceGenerals(fixedRand{0}) {
			if w.Generals[id].Faction == w.Player {
				ev.TalkNotices = append(ev.TalkNotices,
					TalkNotice{Index: freelanceJoinTalk, General: id})
			}
		}
		return ev.TalkNotices
	}
	// 心向 ＝ 玩家 ⇒ 一則 #41，帶那位武將。
	got := notices(mk(4, 4))
	if len(got) != 1 || got[0].Index != freelanceJoinTalk || got[0].General != 3 {
		t.Errorf("投靠玩家 → %+v，want 一則 #%d／武將 3", got, freelanceJoinTalk)
	}
	// 心向 ＝ 別人 ⇒ 不通知。
	if got := notices(mk(4, 9)); len(got) != 0 {
		t.Errorf("投靠別人卻通知了：%+v", got)
	}
}

// ⭐⭐ **真正的接線測試**：跑 `Tick` 到月結，看那一則 #41 有沒有掛在
// 那個 tick 的 `Event` 上。
//
// ⚠ 這一支存在的理由是**上面那支測不到接線**——它自己組 `TalkNotice`，
// 沒有走 `Tick`。把 `state.go` 月結段那幾行拔掉，**全套照樣綠**
// （2026-09-04 的突變測試發現，同一個 session 第三次踩到）。
func TestTickNotifiesFreelanceJoin(t *testing.T) {
	w := load(t, 0)
	w.Player = w.AliveFactions()[0]

	// 挑一個在野武將，心向設成玩家、倒數歸零。
	target := -1
	for i := range w.Generals {
		g := &w.Generals[i]
		if g.Alive && g.Faction == noFaction && g.Captor == noFaction {
			g.Affinity, g.Timer = w.Player, 0
			target = i
			break
		}
	}
	if target < 0 {
		t.Skip("這個劇本沒有在野武將")
	}

	// fixedRand{0}：25% 那道閘一定過（0 < affinityRollLimit）。
	r := fixedRand{0}
	for i := 0; i < 100000; i++ {
		ev := w.Tick(r)
		if !ev.Settled {
			continue
		}
		for _, n := range ev.TalkNotices {
			if n.Index == freelanceJoinTalk && n.General == target {
				return // 通過
			}
		}
		t.Fatalf("月結跑了，但沒有 #%d 的出仕通知（武將 %d 現在屬於勢力 %d）",
			freelanceJoinTalk, target, w.Generals[target].Faction)
	}
	t.Fatal("跑了十萬個 tick 都沒到月結")
}
