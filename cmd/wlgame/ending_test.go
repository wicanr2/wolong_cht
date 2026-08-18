package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/cutscene"
)

// 原版的延遲單位是四分之一個 BIOS tick（INT 1Ch 的處理常式一次加 4，
// docs/re/70 §4）。換成 60 fps 之後每個常數的幀數要落在合理範圍，
// 而且**不能歸零**——歸零的話整段動畫會一幀跑完。
func TestEndingDelayConversion(t *testing.T) {
	cases := []struct{ units, min, max int }{
		{endingFadeStepUnits, 7, 9},      // ≈ 0.14 s ＝ 8 幀
		{endingTypeUnits, 14, 16},        // ≈ 0.25 s ＝ 15 幀
		{endingSceneHoldUnits, 40, 42},   // ≈ 0.69 s ＝ 41 幀
		{endingFinalHoldUnits, 490, 500}, // ≈ 8.2 s ＝ 494 幀
	}
	for _, c := range cases {
		n := endingFrames(c.units)
		if n < c.min || n > c.max {
			t.Errorf("延遲 %d 單位換成 %d 幀，預期 %d–%d", c.units, n, c.min, c.max)
		}
	}
	if endingFrames(0) != 1 {
		t.Error("延遲 0 也要至少一幀，否則那一步等於沒發生")
	}
}

// 十二幕要依序走完，第一幕會多一段打字，最後一幕停住不自己收掉。
func TestEndingPlaysAllScenesInOrder(t *testing.T) {
	e := &endingState{art: &cutscene.Ending{Lines: []string{"一二三"}}}
	seen := map[int]bool{}
	typed := false
	for i := 0; i < 200000 && !e.update(); i++ {
		seen[e.scene] = true
		if e.scene == 0 && e.phase == endingType {
			typed = true
		}
	}
	if !typed {
		t.Error("第一幕沒有走到打字那一段")
	}
	for n := 0; n < cutscene.Scenes; n++ {
		if !seen[n] {
			t.Errorf("第 %d 幕沒有播到", n)
		}
	}
	if !e.finished() {
		t.Fatal("跑完了卻沒有停在最後一幕")
	}
	// 停住之後再推幀不能繼續往前跑。
	if !e.update() || e.scene != cutscene.Scenes-1 {
		t.Errorf("停住之後又動了：scene = %d", e.scene)
	}
}
