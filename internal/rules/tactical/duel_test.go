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

// 氣勢公式：武術門檻沒過整個歸零；過了是 戰力×體力＋亂數尾。
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

// 原版1A2EF／1A2F6先玩家後對方。此四次實測取數會接受挑戰，反序會拒戰。
func TestDuelMoraleUsesPlayerOrder(t *testing.T) {
	for _, player := range []int{AttackerSide, DefenderSide} {
		b := duelBattle(t, 198, [2]int{12, 12}, [2]int{6, 6}, []int{16, 40, 56, 39})
		b.SetPlayerSide(player)
		b.Sides[player].Soldiers[0].HP, b.Sides[player].Soldiers[0].Power = 188, 56
		b.Sides[1-player].Soldiers[0].HP, b.Sides[1-player].Soldiers[0].Power = 156, 32
		duelTick(b)
		if b.duel.strong != player || b.duel.hi != 10528 || b.duel.lo != 6784 {
			t.Fatalf("玩家側%d：strong=%d hi=%d lo=%d", player, b.duel.strong, b.duel.hi, b.duel.lo)
		}
		b.TakeDuelTalks()
		duelTick(b)
		talks := b.TakeDuelTalks()
		if len(talks) != 1 || talks[0].Side != 1-player || talks[0].Group != 0x1B8 {
			t.Fatalf("原版同組取數應走應戰：%v", talks)
		}
	}
}

func TestDuelTiedMoraleKeepsPlayerFirst(t *testing.T) {
	b := duelBattle(t, 198, [2]int{12, 12}, [2]int{6, 6}, []int{0})
	b.SetPlayerSide(DefenderSide)
	duelTick(b)
	if b.duel.strong != DefenderSide {
		t.Fatal("原版同分不交換，記錄0為玩家")
	}
}

type duelCountingRand struct {
	values []int
	calls  int
}

func (r *duelCountingRand) Next() int {
	v := r.values[r.calls%len(r.values)]
	r.calls++
	return v
}

func TestDuelMeleeReusesDecisionByteForY(t *testing.T) {
	b := duelBattle(t, 198, [2]int{12, 12}, [2]int{6, 6}, []int{0})
	r := &duelCountingRand{values: []int{5, 14, 255}}
	b.rng = r
	b.duel.phase, b.duel.timer = duelMelee, 48
	b.stepDuel()
	if r.calls != 2 {
		t.Fatalf("1A2A6與1A2B3只取兩次，實際%d", r.calls)
	}
	for side := range b.Sides {
		g := b.Sides[side].Soldiers[0]
		if g.StepX != 38 || g.StepY != 33 || g.GoalX != 38 || g.GoalY != 33 {
			t.Fatalf("判定byte5兼作Y、byte14作X，應指向38,33：%+v", g)
		}
	}
	b.duel.timer = 1
	b.stepDuel()
	if r.calls != 2 {
		t.Fatal("計時歸零立即返回，不可再抽一次")
	}
}

// 這三個事件邊界出自受控原版290／310／330拍，兩段等待不可合併。
func TestDuelVerdictWaitsAfterBothSpeakers(t *testing.T) {
	b := duelBattle(t, 198, [2]int{12, 12}, [2]int{6, 6}, []int{0})
	b.Frame, b.duel.phase, b.duel.timer = 290, duelMelee, 40
	b.Sides[1].Soldiers[0].HP = 64
	for side := range b.Sides {
		b.Sides[side].Soldiers[0].Cmd, b.Sides[side].Soldiers[0].Next = Duel, Duel
	}
	b.stepDuel()
	talks := b.TakeDuelTalks()
	if len(talks) != 1 || talks[0] != (DuelTalk{Side: 1, Group: 0x1CC, Frame: 290}) {
		t.Fatalf("低於70當拍喊敗方：%v", talks)
	}
	for b.Frame < 310 {
		b.Frame++
		b.stepDuel()
	}
	talks = b.TakeDuelTalks()
	if len(talks) != 1 || talks[0] != (DuelTalk{Side: 0, Group: 0x1CD, Frame: 310}) || !b.OpeningActive() {
		t.Fatalf("310拍喊勝方後仍封鎖主迴圈：%v", talks)
	}
	for b.Frame < 329 {
		b.Frame++
		b.stepDuel()
	}
	if !b.OpeningActive() || b.Sides[0].Soldiers[0].Next != Duel {
		t.Fatal("勝方必須完整再等20拍")
	}
	b.Frame++
	b.stepDuel()
	if b.OpeningActive() {
		t.Fatal("330拍等待完成應返回主迴圈")
	}
	for side := range b.Sides {
		g := b.Sides[side].Soldiers[0]
		if g.Next != Form || g.Cmd != Duel {
			t.Fatal("收尾只清Next，目前Cmd留待下一批更新")
		}
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

// 第50拍先完成一般移動，挑戰命令從第51拍才影響單兵。
func TestDuelChallengeWaitCompletesBeforeMovementChanges(t *testing.T) {
	for _, player := range []int{AttackerSide, DefenderSide} {
		b := NewBattle(flatField(), &Formations{}, &fixedRand{seq: []int{0}}, 0)
		b.SetPlayerSide(player)
		in := DuelInput{FieldNumber: 198}
		for side := range b.Sides {
			x, hp, power, martial, command := 60, 156, 32, 7, 4
			if side == player {
				x, hp, power, martial, command = 3, 188, 56, 12, 6
			}
			b.Sides[side].Line = x // 零偏移陣形，維持此等待邊界的原地目標。
			b.Sides[side].Soldiers[0] = Soldier{Alive: true, Kind: General, X: x, Y: 32, HP: hp, Power: power, Stamina: 127,
				Cmd: Holding, Next: Holding, GoalX: x, GoalY: 32, StepX: x, StepY: 32, Target: -1}
			in.Martial[side], in.CommandStat[side] = martial, command
		}
		b.SetDuelInput(in)
		// 原版最後一拍等待的邊界；整段正常開場另由dosgolem重播驗證。
		b.Frame, b.duel.timer = 49, 1
		b.Step()
		g := b.Sides[player].Soldiers[0]
		if g.X != 3 || g.Y != 32 || g.Stamina != 127 || g.StepX != 24 {
			t.Fatalf("第50拍仍應先完成原地等待：%+v", g)
		}
		talks := b.TakeDuelTalks()
		if len(talks) != 1 || talks[0] != (DuelTalk{Side: player, Group: 0x1B7, Frame: 50}) || g.Next != Duel || g.Cmd != Holding {
			t.Fatalf("完成50拍後應立即掛挑戰及命令：%v", talks)
		}
		b.Step()
		if b.Sides[player].Soldiers[0].Stamina != 126 {
			t.Fatal("第51拍才應按單挑目標開始移動並消耗體力")
		}
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
	if g.Next != Duel {
		t.Errorf("強側大將 next=%v，應為單挑", g.Next)
	}
	// 挑戰是**寫目標**讓大將騎過去，不搬人（spec/80 §3 第 1 點）。
	if g.GoalX != 0x18 || g.GoalY != 0x20 {
		t.Errorf("強側大將目標 (%d,%d)，應為單挑位 (0x18,0x20)", g.GoalX, g.GoalY)
	}
	if x, y := b.duelSpot(0); g.X == x && g.Y == y {
		t.Errorf("挑戰當下大將不該已經在單挑位")
	}
	duelTick(b) // 弱側 lo=0 < 0x12C0 → 拒戰
	duelTick(b) // 強側回應＋立即收尾
	got := b.TakeDuelTalks()
	want := []DuelTalk{{Side: 1, Group: 0x1B9}, {Side: 0, Group: 0x1CC}}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("拒戰喊話 %v，應為 %v", got, want)
	}
	if b.duel.phase != duelIdle || b.Sides[0].Soldiers[0].Next == Duel {
		t.Errorf("拒戰後應收尾並解除命令 8")
	}
}

// 應戰分支：0x1B8 → 第 0 回合互嗆 0x1BA/0x1BB → 會合點 → 決著。
func TestDuelAcceptAndVerdict(t *testing.T) {
	b := duelBattle(t, 0xC0, [2]int{90, 90}, [2]int{0, 0}, []int{0})
	duelTick(b) // 評估：同分不換 → 強側 = 玩家，0x1B7
	duelTick(b) // 應戰 0x1B8
	duelTick(b) // 互嗆第一句 0x1BA
	duelTick(b) // 互嗆第二句 0x1BB → 會合 → 對打段
	got := b.TakeDuelTalks()
	want := []DuelTalk{{Side: 0, Group: 0x1B7}, {Side: 1, Group: 0x1B8}, {Side: 0, Group: 0x1BA}, {Side: 1, Group: 0x1BB}}
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
		if g.Next != Duel {
			t.Errorf("側 %d 大將未收到下一命令 8", i)
		}
		// 互嗆完兩大將的目標都是會合點 (0x20,0x20)——騎過去，不瞬移。
		if g.GoalX != 0x20 || g.GoalY != 0x20 {
			t.Errorf("側 %d 大將目標 (%d,%d)，應為會合點 (0x20,0x20)", i, g.GoalX, g.GoalY)
		}
	}
	// 守側大將被打到 0x40（< 0x46）→ 決著：敗方 0x1CC、勝方 0x1CD。
	b.Sides[1].Soldiers[0].HP = 0x40
	b.stepDuel() // 體力檢查當拍喊敗方
	duelTick(b)  // 敗方等20拍後，勝方喊話
	duelTick(b)  // 勝方再等20拍才收尾
	got = b.TakeDuelTalks()
	want = []DuelTalk{{Side: 1, Group: 0x1CC}, {Side: 0, Group: 0x1CD}}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("決著喊話 %v，應為 %v", got, want)
	}
	if b.duel.phase != duelIdle {
		t.Errorf("決著後未收尾")
	}
	for i := range b.Sides {
		if b.Sides[i].Soldiers[0].Next == Duel {
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
	b.stepDuel() // 敗方退卻中不喊；此測試只驗喊話選擇
	duelTick(b)  // 勝方0x1CD
	duelTick(b)  // 收尾
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
	// 目標在東邊：首拍取目標，下一拍走一格，命令仍是 8。
	s.GoalX = s.X + 4
	s.Path = nil
	b.updateSoldier(0, 0)
	if s.X != x {
		t.Fatal("取新目標的同拍不可直接位移")
	}
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

// 玩家即使是守方，也使用原版記錄 0 的單挑位。
func TestDuelPositionsFollowPlayerRole(t *testing.T) {
	for _, player := range []int{AttackerSide, DefenderSide} {
		b := newTestBattle(flatField())
		b.SetPlayerSide(player)
		b.duelFace()
		for side := range b.Sides {
			wantX := 40
			if side == player {
				wantX = 24
			}
			g := b.Sides[side].Soldiers[0]
			if g.GoalX != wantX || g.GoalY != 32 || g.StepX != wantX || g.StepY != 32 {
				t.Fatalf("玩家側 %d，大將側 %d：目標 %d,%d，應為 %d,32", player, side, g.GoalX, g.GoalY, wantX)
			}
		}
	}
}
