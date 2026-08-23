package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/cutscene"
)

// 原版的延遲單位是四分之一個 BIOS tick（INT 1Ch 的處理常式一次加 4，
// docs/re/70 §4）。換成 60 fps 之後每個常數的幀數要落在合理範圍，
// 而且**不能歸零**——歸零的話整段動畫會一幀跑完。
func TestEndingDelayConversion(t *testing.T) {
	cases := []struct {
		name        string
		units       int
		min, max    int
	}{
		{"第一幕淡入", endingFirstFadeInUnits, 7, 9},        // ≈ 0.14 s
		{"第一幕淡出", endingFirstFadeOutUnits, 4, 7},       // ≈ 0.10 s
		{"「終」淡出", endingFinalFadeOutUnits, 10, 13},      // ≈ 0.19 s
		{"打字", endingTypeUnits, 14, 16},                // ≈ 0.25 s
		{"幕間黑畫面", endingInterludeUnits, 163, 166},      // ≈ 2.7 s
		{"第 2–11 幕淡入淡出", endingSceneFadeUnits, 1, 2},   // ≈ 0.03 s，會被夾到 1 幀
		{"第 2–11 幕黑幕停留", endingSceneHoldUnits, 40, 42}, // ≈ 0.69 s
		{"第 2–11 幕看圖", endingSceneViewUnits, 655, 662}, // ≈ 11.0 s
		{"「終」停留", endingFinalHoldUnits, 490, 500},      // ≈ 8.2 s
		{"最後一幕淡入", endingLastFadeUnits, 5, 8},          // ≈ 0.11 s
	}
	for _, c := range cases {
		n := endingFrames(c.units)
		if n < c.min || n > c.max {
			t.Errorf("%s：%d 單位換成 %d 幀，預期 %d–%d", c.name, c.units, n, c.min, c.max)
		}
	}
	if endingFrames(0) != 1 {
		t.Error("延遲 0 也要至少一幀，否則那一步等於沒發生")
	}
}

// ⭐ 三段各有各的一組常數，**不可以互相套用**。
// 先前把第一幕的 0x0A 套到第 2–11 幕，每幕就快了約 10 秒
// （docs/spec/67 §8）。這一條把每一幕的三個參數釘住。
func TestScenePaceMatchesOriginal(t *testing.T) {
	cases := []struct {
		scene              int
		finalPage          bool
		pre, fade, hold    int
		fadeOut            int
	}{
		// 第一幕第一趟：淡入 0x0A、沒有黑幕停留、hold 由打字取代，淡出 0x07
		{0, false, 0, 0x0A, 0, 0x07},
		// 「終」那一頁：淡入 0x0A、停 0x258、淡出 0x0E
		{0, true, 0, 0x0A, 0x258, 0x0E},
		// 幕 1：黑幕停留要含 start 的 delay(0xC8)
		{1, false, 0xC8 + 0x32, 0x02, 0x320, 0x02},
		{5, false, 0x32, 0x02, 0x320, 0x02},
		{10, false, 0x32, 0x02, 0x320, 0x02},
		// 最後一幕：淡入 8、沒有停留
		{cutscene.Scenes - 1, false, 0, 0x08, 0, 0x08},
	}
	for _, c := range cases {
		pre, fade, hold := scenePace(c.scene, c.finalPage)
		if pre != c.pre || fade != c.fade || hold != c.hold {
			t.Errorf("scenePace(%d, %v) = (%d, %d, %d)，預期 (%d, %d, %d)",
				c.scene, c.finalPage, pre, fade, hold, c.pre, c.fade, c.hold)
		}
		if got := fadeOutUnits(c.scene, c.finalPage); got != c.fadeOut {
			t.Errorf("fadeOutUnits(%d, %v) = %d，預期 %d",
				c.scene, c.finalPage, got, c.fadeOut)
		}
	}
	// ⭐ 看圖時間必須遠大於黑幕停留——先前的 bug 是兩者都是 0x32。
	if endingSceneViewUnits <= endingSceneHoldUnits*4 {
		t.Errorf("看圖時間 %d 沒有明顯大於黑幕停留 %d，回到 2026-08-21 之前的 bug 了",
			endingSceneViewUnits, endingSceneHoldUnits)
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

// ⭐ 第一幕打完字要**先淡出到黑再換頁**，不是在全亮的狀態下直接換
// （原版 sub_10204 在淡出之後才跑，docs/re/70 §4）。
func TestFirstSceneFadesOutBeforeFinalPage(t *testing.T) {
	e := &endingState{art: &cutscene.Ending{Lines: []string{"一二三"}}}
	sawTyping, sawDark := false, false
	for i := 0; i < 200000; i++ {
		if e.scene != 0 {
			break
		}
		if e.phase == endingType {
			sawTyping = true
			if e.finalPage {
				t.Fatal("還在打字就翻到「終」那一頁了")
			}
		}
		// 換頁的那一刻畫面必須是全黑的（step 回到 0）。
		if e.finalPage && !sawDark {
			if e.step != 0 {
				t.Fatalf("換到「終」那一頁時亮度是 %d／%d，應該是全黑",
					e.step, cutscene.FadeSteps-1)
			}
			sawDark = true
		}
		if e.update() {
			break
		}
	}
	if !sawTyping {
		t.Error("沒有走到打字那一段")
	}
	if !sawDark {
		t.Error("第一幕從頭到尾沒有翻到「終」那一頁")
	}
}

// 節拍改對之後整段結局應該長很多：第 2–11 幕各有一段 11 秒的看圖時間。
// 這一條擋「看圖時間被拿掉」那一類回歸——只看總長，不綁死實作。
func TestEndingIsLongEnoughToRead(t *testing.T) {
	e := &endingState{art: &cutscene.Ending{Lines: []string{"一二三"}}}
	frames := 0
	for ; frames < 200000 && !e.update(); frames++ {
	}
	// 光是第 2–11 幕的看圖時間就是 10 × 11.0 s ＝ 110 s。
	const minFrames = 10 * 655
	if frames < minFrames {
		t.Errorf("整段結局只有 %d 幀（約 %d 秒），比第 2–11 幕的看圖時間總和還短——"+
			"看圖時間大概又被拿掉了", frames, frames/endingFPS)
	}
}

// ⭐ 結局要放 `endbgm-0`，不是 `overbgm-0`。
// 結局播的時候勝負早就定了，所以這一格若排在 Outcome 之後，
// 整段結局會放成遊戲結束曲——而那是另一支執行檔的配樂（docs/re/58 §6）。
func TestEndingUsesEndingMusic(t *testing.T) {
	g := &game{}
	g.ending = &endingState{art: &cutscene.Ending{Lines: []string{"一"}}}
	if got := g.musicTrack(); got != "endbgm-0" {
		t.Errorf("結局播 %q，應該是 endbgm-0", got)
	}
	// 反向對照：結局收掉之後不可以還黏在 endbgm。
	g.ending = nil
	if got := g.musicTrack(); got == "endbgm-0" {
		t.Error("結局結束了還在播 endbgm-0")
	}
}
