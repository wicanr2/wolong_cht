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

// tick 讓狀態機把目前的等待走完並推進一個相位。
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
		t.Errorf("氣勢 %d，應為 100 兵 × 100 體力 = 10000", got)
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

// 拒戰分支：強側 0x1B7 → 弱側 0x1B9 → 強側 0x1CC → 命令歸 0。
func TestDuelRefuse(t *testing.T) {
	b := duelBattle(t, 0xC0, [2]int{90, 0}, [2]int{0, 0}, []int{0})
	duelTick(b) // 評估：側 0 保留 10000、側 1 歸零 → 挑戰
	talks := b.TakeDuelTalks()
	if len(talks) != 1 || talks[0] != (DuelTalk{Side: 0, Group: 0x1B7}) {
		t.Fatalf("挑戰喊話 %v，應為側 0 組 0x1B7", talks)
	}
	if got := b.Sides[0].Soldiers[0].Cmd; got != Duel {
		t.Errorf("強側大將 cmd=%v，應為單挑", got)
	}
	duelTick(b) // 弱側 lo=0 < 0x12C0 → 拒戰
	duelTick(b) // 強側回應
	duelTick(b) // 收尾
	got := b.TakeDuelTalks()
	want := []DuelTalk{{1, 0x1B9}, {0, 0x1CC}}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("拒戰喊話 %v，應為 %v", got, want)
	}
	if b.duel.phase != duelIdle || b.Sides[0].Soldiers[0].Cmd == Duel {
		t.Errorf("拒戰後應收尾並解除命令 8")
	}
}

// 應戰分支：0x1B8 → 第 0 回合互嗆 0x1BA/0x1BB → 對打段 → 決著。
func TestDuelAcceptAndVerdict(t *testing.T) {
	b := duelBattle(t, 0xC0, [2]int{90, 90}, [2]int{0, 0}, []int{0})
	duelTick(b) // 評估：同分不換 → 強側 = 攻方
	duelTick(b) // 應戰
	duelTick(b) // 互嗆第一句
	duelTick(b) // 互嗆第二句 → 對打段
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
		x, y := duelSpot(i)
		if g.X != x || g.Y != y {
			t.Errorf("側 %d 大將 (%d,%d)，應戰後應在對峙位 (%d,%d)", i, g.X, g.Y, x, y)
		}
	}
	// 守側大將被打到 0x40（< 0x46）→ 決著：敗方 0x1CC、勝方 0x1CD。
	b.Sides[1].Soldiers[0].HP = 0x40
	b.stepDuel() // 體力檢查 → 敗方喊話
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
	b.stepDuel()
	duelTick(b)
	got := b.TakeDuelTalks()
	if len(got) != 1 || got[0] != (DuelTalk{Side: 0, Group: 0x1CD}) {
		t.Fatalf("退卻中敗方仍喊話：%v，應只剩勝方 0x1CD", got)
	}
}

// 命令 8 的凍結：updateSoldier 對單挑中的兵是 no-op（原版 nullsub）。
func TestDuelCommandFreezesSquad(t *testing.T) {
	b := newTestBattle(flatField())
	s := &b.Sides[0].Soldiers[0]
	s.Cmd, s.Next = Duel, Duel
	x, y := s.X, s.Y
	b.updateSoldier(0, 0)
	if s.X != x || s.Y != y {
		t.Errorf("單挑中的兵動了：(%d,%d)→(%d,%d)", x, y, s.X, s.Y)
	}
	if s.Cmd != Duel {
		t.Errorf("命令 8 被 updateSoldier 改掉：%v", s.Cmd)
	}
}
