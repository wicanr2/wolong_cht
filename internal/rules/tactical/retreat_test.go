package tactical

import "testing"

// 全軍退卻之後**第 120 拍**結束，勝方是沒退卻的那一側（docs/spec/141）。
//
// ⚠ 這是**兜底**那一條出口。原版正常打完走的是「補不出兵」——退卻的兵
// 八拍就走完了（docs/playtest/80）。這裡的測試場兩側各有 552 個待機兵，
// 補得比走得快，所以撐得到倒數走完。
func TestRetreatEndsBattleAfter120Ticks(t *testing.T) {
	// 原版 `sub_19A33` 的 `mov cs:byte_1D34A, 78h`。**釘住立即值本身**——
	// 只用常數比對的話，把 120 改成別的數也照樣全綠。
	if RetreatCountdown != 0x78 {
		t.Fatalf("倒數 = %d 拍，原版的立即值是 78h（120）", RetreatCountdown)
	}
	for _, side := range []int{0, 1} {
		b := newTestBattle(flatField())
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
