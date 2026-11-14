package tactical

import "testing"

// 全軍退卻之後**第 120 拍**結束，勝方是沒退卻的那一側（docs/spec/141）。
//
// ⚠ 這是**兜底**那一條出口。原版正常打完走的是「補不出兵」——退卻的兵
// 八拍就走完了（docs/playtest/80）。這裡用硬直留住場上部隊，
// 明確建立「退不出去」的兜底前提。
func TestRetreatEndsBattleAfter120Ticks(t *testing.T) {
	// 原版 `sub_19A33` 的 `mov cs:byte_1D34A, 78h`。**釘住立即值本身**——
	// 只用常數比對的話，把 120 改成別的數也照樣全綠。
	if RetreatCountdown != 0x78 {
		t.Fatalf("倒數 = %d 拍，原版的立即值是 78h（120）", RetreatCountdown)
	}
	for _, side := range []int{0, 1} {
		b := newTestBattle(flatField())
		holdRetreatCountdownUnits(b)
		b.Order(side, -1, Retreat)
		for i := 1; i <= RetreatCountdown; i++ {
			if b.Done {
				t.Fatalf("側 %d 退卻：第 %d 拍就結束了，應該撐到第 %d 拍",
					side, i, RetreatCountdown)
			}
			b.Step()
		}
		if !b.Done {
			t.Errorf("側 %d 退卻：走完 %d 拍還沒結束", side, RetreatCountdown)
			continue
		}
		if b.Winner != 1-side {
			t.Errorf("側 %d 退卻，勝方 = %d，應為 %d", side, b.Winner, 1-side)
		}
	}
}

// 沒有人退卻時倒數不動——原版 `sub_1A6FA` 開頭就 `cmp byte_1D349, 0 / jz`。
func TestRetreatCountdownOnlyRunsWhileRetreating(t *testing.T) {
	b := newTestBattle(flatField())
	for i := 0; i < 40 && !b.Done; i++ {
		b.Step()
	}
	if b.retreat != RetreatCountdown {
		t.Errorf("沒有人退卻卻走了 %d 拍倒數", RetreatCountdown-b.retreat)
	}
	// 個別的兵受傷改退卻**不算全軍退卻**，一樣不起倒數。
	b.Sides[0].Soldiers[3].Next = Retreat
	b.Step()
	if b.retreat != RetreatCountdown {
		t.Error("單一個兵改退卻不該起倒數")
	}
}

// 第二側也下退卻令時倒數不重來——原版 `sub_1A8F6` 開頭
// `cmp cs:byte_1D349, 0 / jnz → stc` 直接把第二次退回去。
func TestSecondRetreatDoesNotRestartCountdown(t *testing.T) {
	b := newTestBattle(flatField())
	holdRetreatCountdownUnits(b)
	b.Order(0, -1, Retreat)
	const half = 60
	for i := 0; i < half; i++ {
		b.Step()
	}
	if b.Done {
		t.Fatalf("第 %d 拍就結束了", half)
	}
	b.Order(1, -1, Retreat) // ← 不受理
	for i := 0; i < RetreatCountdown-half; i++ {
		if b.Done {
			t.Fatalf("倒數重來了：第 %d 拍就結束", half+i)
		}
		b.Step()
	}
	if !b.Done {
		t.Fatalf("走完 %d 拍還沒結束", RetreatCountdown)
	}
	if b.Winner != 1 {
		t.Errorf("勝方 = %d，應為 1（先退卻的是側 0）", b.Winner)
	}
}

func holdRetreatCountdownUnits(b *Battle) {
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			b.Sides[side].Soldiers[k].Stun = 255
		}
	}
}

func TestFullRetreatStopsOnlyRetiringSideReinforcements(t *testing.T) {
	for player := 0; player < 2; player++ {
		for retiring := 0; retiring < 2; retiring++ {
			b := newTestBattle(flatField())
			b.SetPlayerSide(player)
			for side := range b.Sides {
				b.Sides[side].Reserve[0] = 17
				b.Sides[side].Soldiers[1].Alive = false
			}
			b.Order(retiring, -1, Retreat)
			b.reinforce()
			if b.Sides[retiring].Soldiers[1].Alive || b.Sides[retiring].Reserve[0] != 17 {
				t.Fatalf("player%d retiring%d: 全軍撤退不得補兵或消耗待機", player, retiring)
			}
			other := &b.Sides[1-retiring]
			if !other.Soldiers[1].Alive || other.Reserve[0] != 16 {
				t.Fatalf("player%d retiring%d: 另一側仍應正常補兵", player, retiring)
			}
		}
	}
	// 個別受傷撤退和隊長倒下，不開全軍撤退閘。
	b := newTestBattle(flatField())
	b.Sides[0].Soldiers[0].Alive = false
	b.Sides[0].Soldiers[1].Alive = false
	b.Sides[0].Soldiers[2].Cmd, b.Sides[0].Soldiers[2].Next = Retreat, Retreat
	b.Sides[0].Reserve[0] = 17
	b.reinforce()
	if !b.Sides[0].Soldiers[1].Alive || b.Sides[0].Reserve[0] != 16 {
		t.Fatal("隊長不在與個別退卻仍應補兵")
	}
}

func TestVictoryUsesMovementCountAndPlayerRoleOrder(t *testing.T) {
	for player := 0; player < 2; player++ {
		for empty := 0; empty < 2; empty++ {
			b := newTestBattle(flatField())
			b.SetPlayerSide(player)
			// 場上計數歸零仍可能有大量待機，不能繼續等待它們歸零。
			counts := [2]int{48, 48}
			counts[empty] = 0
			reserve := b.Sides[empty].Reserve
			b.checkVictory(counts)
			if !b.Done || b.Winner != 1-empty || b.Sides[empty].Reserve != reserve {
				t.Fatalf("player%d empty%d: 勝負應讀移動計數並保留待機", player, empty)
			}
		}
		b := newTestBattle(flatField())
		b.SetPlayerSide(player)
		b.checkVictory([2]int{})
		if !b.Done || b.Winner != 1-player {
			t.Fatalf("player%d: 同時空場須先判玩家輸", player)
		}
	}
}

func TestVictoryDoesNotRecountUnitsAfterMovement(t *testing.T) {
	b := newTestBattle(flatField())
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			b.Sides[side].Soldiers[k].Alive = false
		}
	}
	// 移動分支處理過的單兵仍計入當拍，不能用最終Alive取代。
	b.checkVictory([2]int{1, 1})
	if b.Done {
		t.Fatal("不得在計數完成後重數Alive")
	}
}

func TestStepCountsOnlyUnitsEnteringMovementBatch(t *testing.T) {
	b := newTestBattle(flatField())
	holdRetreatCountdownUnits(b)
	b.Step()
	if b.Done {
		t.Fatal("場上兩軍仍在，移動批次不得漏計")
	}
	b = newTestBattle(flatField())
	for k := range b.Sides[0].Soldiers {
		b.Sides[0].Soldiers[k].Alive = false
	}
	b.Step()
	if b.Sides[0].Alive() == 0 {
		t.Fatal("前提：待機兵應補進場")
	}
	if !b.Done || b.Winner != 1 {
		t.Fatal("補兵分支當拍不進場上計數，不得改用補完後的Alive")
	}
}

func TestGlobalRetreatQueuesOnlySixLeadersAndRejectsRepeat(t *testing.T) {
	for side := 0; side < 2; side++ {
		b := newTestBattle(flatField())
		for i := range b.Sides {
			for k := range b.Sides[i].Soldiers {
				b.Sides[i].Soldiers[k].Cmd, b.Sides[i].Soldiers[k].Next = Guard, Guard
			}
		}
		b.Sides[side].Soldiers[PerSquad].Alive = false
		b.Order(side, -1, Retreat)
		for k, s := range b.Sides[side].Soldiers {
			want := Guard
			if k%PerSquad == 0 {
				want = Retreat
			}
			if s.Cmd != Guard || s.Next != want {
				t.Fatalf("side%d slot%d：全軍令應只寫六個隊長Next", side, k)
			}
		}
		before := b.Sides
		b.Order(1-side, -1, Retreat)
		for i := range b.Sides {
			for k, s := range b.Sides[i].Soldiers {
				if s.Next != before[i].Soldiers[k].Next {
					t.Fatal("第二次全軍撤退應整次拒絕，不得改另一側命令")
				}
			}
		}
		if b.endPhase != side+1 {
			t.Fatal("第二次撤退不得改變退卻側")
		}
	}
}

func TestRetreatInitializesStepAndPreservesPreviousCommand(t *testing.T) {
	for player := 0; player < 2; player++ {
		for side := 0; side < 2; side++ {
			for _, k := range []int{0, 1} {
				b := newTestBattle(flatField())
				b.SetPlayerSide(player)
				s := &b.Sides[side].Soldiers[k]
				*s = Soldier{Alive: true, X: 20, Y: 5, Z: 3, Cmd: Guard, Next: Retreat,
					Target: -1, Path: &Waypoints{pts: []Point{{X: 30, Y: 20}}}}
				b.updateSoldierCommand(side, k)
				edge := MinCoord
				if b.Role(side) == 1 {
					edge = MaxCoord
				}
				if s.Cmd != Retreat || s.Next != Guard || s.GoalX != edge || s.GoalY != 16 || s.GoalZ != 0 ||
					s.StepX != edge || s.StepY != 16 || s.StepZ != 0 || s.Path != nil || !s.PathQueued {
					t.Fatalf("player%d side%d slot%d：首次退卻未保存舊命令及出口／尋路請求", player, side, k)
				}
				if k == 0 {
					for j := 1; j < PerSquad; j++ {
						if b.Sides[side].Soldiers[j].Next != Retreat || !b.Sides[side].Soldiers[j].TargetRefresh {
							t.Fatal("新退卻隊長須向七隊員傳令並要求目標刷新")
						}
					}
				}
				// 繞路後Y改變，仍保留接令時的出口Y與正在消費的中繼點。
				s.Y, s.StepX, s.StepY = 20, 19, 20
				s.Path = &Waypoints{pts: []Point{{X: edge, Y: 16}}}
				path := s.Path
				b.updateSoldierCommand(side, k)
				if s.GoalY != 16 || s.StepX != 19 || s.StepY != 20 || s.Path != path || s.Next != Guard {
					t.Fatal("持續退卻不得重設出口、中繼點、舊命令或路徑")
				}
			}
		}
	}
}

func TestGeneralRetreatChecksBeforeSiegeDrainAndAfterCommands(t *testing.T) {
	f, _ := tiledField(32)
	b := newTestBattle(f)
	b.Sides[0].Soldiers[0].HP = GeneralRetreatHP
	b.Sides[1].Soldiers[0].HP = 100
	b.siegeTick = SiegeDrainInterval - 1
	b.Step()
	if b.Sides[0].Soldiers[0].HP != 49 || b.endPhase != 0 {
		t.Fatal("體力50先過門檻檢查，再扣成49，當拍尚未全軍退卻")
	}
	b.Step()
	if b.endPhase != 1 || b.Sides[0].Soldiers[0].Cmd == Retreat {
		t.Fatal("下一拍移動前才下令；已完成的命令批次不能提早再跑")
	}
	b.Step()
	if b.Sides[0].Soldiers[0].Cmd != Retreat {
		t.Fatal("再下一次命令批次應開始退卻")
	}
	for player := 0; player < 2; player++ {
		b := newTestBattle(flatField())
		b.SetPlayerSide(player)
		for side := range b.Sides {
			b.Sides[side].Soldiers[0].HP = 49
		}
		b.checkGeneralRetreat()
		if b.endPhase != player+1 {
			t.Fatal("雙方體力皆低時，須先受理玩家角色的全軍退卻")
		}
	}
}
