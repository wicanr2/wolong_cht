package tactical

import "testing"

// duelBattle 開一場 100 人 × 6 隊的野戰，武裝單挑狀態機。
func duelBattle(t *testing.T, field int, martial, cmdstat [2]int, seq []int) *Battle {
	t.Helper()
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: seq}, 0)
	for k := 0; k < Squads; k++ {
		b.Deploy(0, k, Infantry, 100)
		b.Deploy(1, k, Infantry, 100)
	}
	b.Place()
	b.SetDuelInput(DuelInput{FieldNumber: field, Martial: martial, CommandStat: cmdstat})
	return b
}

// duelTick 讓狀態機把目前的等待走完並推進一個相位。
func duelTick(b *Battle) {
	b.duel.timer = 1
	b.stepDuel()
}

// 閘：只有戰場編號 0xC0–0xD0（byte_1D34B == 1）武裝。
func TestDuelGateByFieldNumber(t *testing.T) {
	for _, tc := range []struct {
		field int
		armed bool
	}{{0xBF, false}, {0xC0, true}, {0xD0, true}, {0xD1, false}} {
		b := duelBattle(t, tc.field, [2]int{90, 90}, [2]int{0, 0}, []int{0})
		if b.duel.armed != tc.armed {
			t.Errorf("編號 %#x：armed=%v，應為 %v", tc.field, b.duel.armed, tc.armed)
		}
	}
}

// 氣勢公式：武術門檻沒過整個歸零；過了是 兵數×體力＋亂數尾。
func TestDuelMoraleFormula(t *testing.T) {
	// 武術 90／統率 0 → 門檻 135，rand(0..7)+8 永遠比不過 → 保留。
	b := duelBattle(t, 0xC0, [2]int{90, 0}, [2]int{0, 0}, []int{0})
	if got := b.duelMorale(0); got != 100*100 {
		t.Errorf("氣勢 %d，應為 戰力 100 × 體力 100 = 10000", got)
	}
	// 武術 0 → 門檻 0 → 一定歸零，只剩亂數尾（rand=3 → 3<<8）。
	b.rng = &fixedRand{seq: []int{3}}
	if got := b.duelMorale(1); got != 3<<8 {
		t.Errorf("歸零側氣勢 %d，應為 3<<8=%d", got, 3<<8)
	}
}

// 兩側都不到 0x12C0：無事發生，正常開打。
func TestDuelNoChallengeWhenBothWeak(t *testing.T) {
	b := duelBattle(t, 0xC0, [2]int{0, 0}, [2]int{0, 0}, []int{0})
	duelTick(b)
	if b.duel.phase != duelIdle {
		t.Fatalf("相位 %d，應直接收尾", b.duel.phase)
	}
	if talks := b.TakeDuelTalks(); len(talks) != 0 {
		t.Errorf("不該有喊話，卻有 %v", talks)
	}
}

// 拒戰分支：強側 0x1B7 → 弱側 0x1B9 → 強側 0x1CC → 立即清命令。
func TestDuelRefuse(t *testing.T) {
	b := duelBattle(t, 0xC0, [2]int{90, 0}, [2]int{0, 0}, []int{0})
	duelTick(b) // 評估：側 0 保留 10000、側 1 歸零 → 挑戰
	talks := b.TakeDuelTalks()
	if len(talks) != 1 || talks[0] != (DuelTalk{Side: 0, Group: 0x1B7}) {
		t.Fatalf("挑戰喊話 %v，應為側 0 組 0x1B7", talks)
	}
	g := &b.Sides[0].Soldiers[0]
	if g.Cmd != Duel {
		t.Errorf("強側大將 cmd=%v，應為單挑", g.Cmd)
	}
	// 挑戰是**寫目標**讓大將騎過去，不搬人（spec/80 §3 第 1 點）。
	if g.GoalX != 0x18 || g.GoalY != 0x20 {
		t.Errorf("強側大將目標 (%d,%d)，應為單挑位 (0x18,0x20)", g.GoalX, g.GoalY)
	}
	if x, y := duelSpot(0); g.X == x && g.Y == y {
		t.Errorf("挑戰當下大將不該已經在單挑位")
	}
	duelTick(b) // 弱側 lo=0 < 0x12C0 → 拒戰
	duelTick(b) // 強側回應＋立即收尾
	got := b.TakeDuelTalks()
	want := []DuelTalk{{1, 0x1B9}, {0, 0x1CC}}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("拒戰喊話 %v，應為 %v", got, want)
	}
	if b.duel.phase != duelIdle || b.Sides[0].Soldiers[0].Cmd == Duel {
		t.Errorf("拒戰後應收尾並解除命令 8")
	}
}

// 應戰分支：0x1B8 → 第 0 回合互嗆 0x1BA/0x1BB → 會合點 → 決著。
func TestDuelAcceptAndVerdict(t *testing.T) {
	b := duelBattle(t, 0xC0, [2]int{90, 90}, [2]int{0, 0}, []int{0})
	duelTick(b) // 評估：同分不換 → 強側 = 攻方，0x1B7
	duelTick(b) // 應戰 0x1B8
	duelTick(b) // 互嗆第一句 0x1BA
	duelTick(b) // 互嗆第二句 0x1BB → 會合 → 對打段
	got := b.TakeDuelTalks()
	want := []DuelTalk{{0, 0x1B7}, {1, 0x1B8}, {0, 0x1BA}, {1, 0x1BB}}
	if len(got) != len(want) {
		t.Fatalf("喊話 %v，應為 %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("喊話 %v，應為 %v", got, want)
		}
	}
	if b.duel.phase != duelMelee {
		t.Fatalf("相位 %d，應在對打段", b.duel.phase)
	}
	for i := range b.Sides {
		g := &b.Sides[i].Soldiers[0]
		if g.Cmd != Duel {
			t.Errorf("側 %d 大將未進命令 8", i)
		}
		// 互嗆完兩大將的目標都是會合點 (0x20,0x20)——騎過去，不瞬移。
		if g.GoalX != 0x20 || g.GoalY != 0x20 {
			t.Errorf("側 %d 大將目標 (%d,%d)，應為會合點 (0x20,0x20)", i, g.GoalX, g.GoalY)
		}
	}
	// 守側大將被打到 0x40（< 0x46）→ 決著：敗方 0x1CC、勝方 0x1CD。
	b.Sides[1].Soldiers[0].HP = 0x40
	b.stepDuel() // 體力檢查 → 決著相位
	duelTick(b)  // 敗方喊話
	duelTick(b)  // 勝方喊話、收尾
	got = b.TakeDuelTalks()
	want = []DuelTalk{{1, 0x1CC}, {0, 0x1CD}}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("決著喊話 %v，應為 %v", got, want)
	}
	if b.duel.phase != duelIdle {
		t.Errorf("決著後未收尾")
	}
	for i := range b.Sides {
		if b.Sides[i].Soldiers[0].Cmd == Duel {
			t.Errorf("側 %d 決著後命令未解除", i)
		}
	}
}

// 回合 ≥1 的互嗆先講側是**體力高側**；體力差 < 0x14 用「勢均」pair。
func TestDuelBanterSpeakerAndPair(t *testing.T) {
	b := duelBattle(t, 0xC0, [2]int{90, 90}, [2]int{0, 0}, []int{0})
	b.duel.strong = 0
	b.duel.round = 1
	b.Sides[0].Soldiers[0].HP = 80
	b.Sides[1].Soldiers[0].HP = 100
	if pair := b.duelBanterPair(); pair != 0x1BC || b.duel.first != 1 {
		t.Errorf("round1 差 20：pair=%#x first=%d，應為 0x1BC／體力高側 1", pair, b.duel.first)
	}
	b.Sides[0].Soldiers[0].HP = 95
	if pair := b.duelBanterPair(); pair != 0x1BE || b.duel.first != 1 {
		t.Errorf("round1 差 5：pair=%#x first=%d，應為勢均 0x1BE／側 1", pair, b.duel.first)
	}
	b.duel.round = 3
	if pair := b.duelBanterPair(); pair != 0x1C6 {
		t.Errorf("round3 勢均 pair=%#x，應為 0x1BC+2×4+2=0x1C6", pair)
	}
}

// 敗方已在退卻（命令 5）就不喊 0x1CC，只有勝方的 0x1CD。
func TestDuelVerdictSkipsRetreatingLoser(t *testing.T) {
	b := duelBattle(t, 0xC0, [2]int{90, 90}, [2]int{0, 0}, []int{0})
	duelTick(b)
	duelTick(b)
	duelTick(b)
	duelTick(b)
	b.TakeDuelTalks()
	b.Sides[1].Soldiers[0].HP = 0x40
	b.Sides[1].Soldiers[0].Cmd = Retreat
	b.stepDuel() // 體力檢查 → 決著相位
	duelTick(b)  // 敗方段（退卻中：不喊）
	duelTick(b)  // 勝方 0x1CD、收尾
	got := b.TakeDuelTalks()
	if len(got) != 1 || got[0] != (DuelTalk{Side: 0, Group: 0x1CD}) {
		t.Fatalf("退卻中敗方仍喊話：%v，應只剩勝方 0x1CD", got)
	}
}

// 命令 8：分派 no-op（nullsub），但**照常走向目標**——移動在共通路徑。
func TestDuelCommandMovesTowardGoalOnly(t *testing.T) {
	b := newTestBattle(flatField())
	s := &b.Sides[0].Soldiers[0]
	s.Cmd, s.Next = Duel, Duel
	// 目標＝原地：不動。
	s.GoalX, s.GoalY, s.GoalZ = s.X, s.Y, s.Z
	x, y := s.X, s.Y
	b.updateSoldier(0, 0)
	if s.X != x || s.Y != y {
		t.Errorf("目標原地卻動了：(%d,%d)→(%d,%d)", x, y, s.X, s.Y)
	}
	// 目標在東邊：走一格，命令仍是 8。
	s.GoalX = s.X + 4
	s.Path = nil
	b.updateSoldier(0, 0)
	if s.X != x+1 {
		t.Errorf("朝目標走了 %d 格，應為 1", s.X-x)
	}
	if s.Cmd != Duel {
		t.Errorf("命令 8 被 updateSoldier 改掉：%v", s.Cmd)
	}
}

// 開場 50 tick ＋ 單挑期間腳本不跑（OpeningActive）。
func TestOpeningBlocksScripts(t *testing.T) {
	b := newTestBattle(flatField())
	if !b.OpeningActive() {
		t.Fatal("開場第 0 tick 就該是 opening")
	}
	for b.Frame <= duelOpeningTicks {
		b.Step()
	}
	if b.OpeningActive() {
		t.Fatalf("第 %d tick 還在 opening（未武裝單挑）", b.Frame)
	}
	b.SetDuelInput(DuelInput{FieldNumber: 0xC0, Martial: [2]int{90, 90}})
	duelTick(b)
	if !b.OpeningActive() {
		t.Fatal("單挑進行中應算 opening（腳本與輸入被擋）")
	}
}
