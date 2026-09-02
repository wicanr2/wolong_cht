package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/speed"
)

// resumeFrames 是門檻換算成畫面更新次數：160 × 600 ÷ 2913，進位。
var resumeFrames = (idleResumeUnits + speed.UnitsPerFrame - 1) / speed.UnitsPerFrame

// still 送 n 個「游標沒動、也沒有輸入」的畫面，回報最後一個是否放行。
func still(g *idleClockGate, n int) bool {
	ok := false
	for i := 0; i < n; i++ {
		ok = g.Allows(12, 34, false)
	}
	return ok
}

func TestIdleClockGateRequiresStablePointerAndNoCommand(t *testing.T) {
	var gate idleClockGate
	cases := []struct {
		name        string
		x, y        int
		inputActive bool
		want        bool
	}{
		{name: "first observation is not idle", x: 12, y: 34, want: false},
		// ⚠ 停一幀**不夠**：原版還要等 160 個回呼（docs/spec/112 §1）。
		{name: "one still frame is not enough", x: 12, y: 34, want: false},
		{name: "pointer movement pauses clock", x: 13, y: 34, want: false},
		{name: "command pauses clock even when pointer is stable", x: 13, y: 34, inputActive: true, want: false},
		{name: "still frame after a command is still not enough", x: 13, y: 34, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := gate.Allows(tc.x, tc.y, tc.inputActive); got != tc.want {
				t.Fatalf("Allows(%d, %d, %t) = %t, want %t", tc.x, tc.y, tc.inputActive, got, tc.want)
			}
		})
	}
}

// 游標一直動就一直被重設，世界永遠不前進——原版 `sub_11F7F` 只在座標
// 變了那一支寫 `ds:98A5h`，所以移動期間是**完全停住**不是變慢
// （docs/spec/112 §1）。
func TestIdleClockGateHoldsWhileCursorMoves(t *testing.T) {
	var gate idleClockGate
	for i := 0; i < 500; i++ {
		if gate.Allows(i%40, i%30, false) {
			t.Fatalf("第 %d 幀游標還在動就放行了", i)
		}
	}
}

// 停下來之後要等滿 160 個回呼（`ds:0D2Dh` 的 `0A0h`）才恢復。
func TestIdleClockGateResumesAfterDelay(t *testing.T) {
	if resumeFrames != 33 {
		t.Fatalf("門檻換算 ＝ %d 幀，docs/spec/112 §4 記的是 33 幀", resumeFrames)
	}
	var gate idleClockGate
	gate.Allows(12, 34, false) // 第一次一定 false（座標還沒穩定）
	if still(&gate, resumeFrames-1) {
		t.Errorf("第 %d 個靜止幀就放行，太早", resumeFrames-1)
	}
	if !still(&gate, 1) {
		t.Errorf("第 %d 個靜止幀還沒放行", resumeFrames)
	}
	// 放行之後只要游標不動就持續放行，不會每幀重新等一次。
	if !still(&gate, 5) {
		t.Error("恢復之後又被擋住")
	}
}

// 有命令輸入的那一幀不前進，而且把等待重新開始。
func TestIdleClockGateInputRestartsDelay(t *testing.T) {
	var gate idleClockGate
	gate.Allows(12, 34, false)
	still(&gate, resumeFrames)
	if gate.Allows(12, 34, true) {
		t.Fatal("有輸入的那一幀不該放行")
	}
	if still(&gate, resumeFrames-1) {
		t.Error("輸入之後沒有重新等滿")
	}
	if !still(&gate, 1) {
		t.Error("輸入之後等滿了卻沒恢復")
	}
}

// 訊息框收掉之後也要重新等——原版 `sub_18810` 在擦掉框之後設 `8`。
func TestIdleClockGatePauseRestartsDelay(t *testing.T) {
	var gate idleClockGate
	gate.Allows(12, 34, false)
	if !still(&gate, resumeFrames) {
		t.Fatal("等夠了卻沒放行，測試前提就不成立")
	}
	gate.Pause()
	if still(&gate, resumeFrames-1) {
		t.Error("Pause 之後沒有重新等滿")
	}
	if !still(&gate, 1) {
		t.Error("Pause 之後等滿了卻沒恢復")
	}
}
