package main

import "testing"

// ⛔ 不認得的值要**回錯誤**，不是靜靜當成留白。
//
// 少了這一條，`-shot-when gatebar`（少一個橫線）會退回「照 `-shot-frames`
// 截圖」——輸出看起來完全正常，拍到的卻是別的局面（docs/spec/118 §2.2）。
func TestParseShotWhenRejectsUnknown(t *testing.T) {
	for _, bad := range []string{
		"gatebar", "gate bar", "battle-frame", "battle-frame:", "battle-frame:x",
		"battle-frame:-1", "門強度",
	} {
		if c, err := parseShotWhen(bad); err == nil {
			t.Errorf("-shot-when %q 應該回錯誤，卻拿到 %+v", bad, c)
		}
	}
	// 留白是合法的：那就是「照 -shot-frames」，也是預設。
	c, err := parseShotWhen("")
	if err != nil || c != nil {
		t.Errorf("留白應該回 (nil, nil)，拿到 (%+v, %v)", c, err)
	}
	// 前後空白不算打錯字。
	if c, err := parseShotWhen("  gate-bar "); err != nil || c == nil || c.name != "gate-bar" {
		t.Errorf("前後空白應該吃掉，拿到 (%+v, %v)", c, err)
	}
}

// 三個條件在「沒有戰鬥」時一律不成立——這一端與成立那一端一樣重要：
// 只斷言成立的那一端，一個永遠回 true 的實作也會通過。
func TestShotConditionsFalseWithoutBattle(t *testing.T) {
	g := &game{}
	for _, name := range []string{"battle", "battle-frame:0", "gate-bar"} {
		c, err := parseShotWhen(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if c.check(g) {
			t.Errorf("%s：沒有世界也沒有戰鬥，不該成立", name)
		}
	}
	// 沒帶 `-shot-when` 時 shotReady 恆真（原本的行為一個字都沒變）。
	if !g.shotReady() {
		t.Error("沒帶 -shot-when 卻不肯截圖")
	}
}

// `battle-frame:N` 的邊界是 **≥**，`gate-bar` 跟著 `StructureBar()` 走。
//
// 用攻城那條 fixture，因為門強度條**只在攻城戰、而且只對城壁**顯示
// （docs/spec/32 §2）——野戰那一條驗不到它。
func TestShotConditionsOnALiveBattle(t *testing.T) {
	g := siegeShotFixture(t)
	b := g.shotBattle()
	if b == nil {
		t.Fatal("fixture 沒有開出戰鬥")
	}

	holds := func(name string) bool {
		c, err := parseShotWhen(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		return c.check(g)
	}

	if !holds("battle") {
		t.Error("戰鬥開著，`battle` 應該成立")
	}
	if !holds("battle-frame:0") {
		t.Error("第 0 幀對 `battle-frame:0` 應該成立（判準是 ≥）")
	}
	if holds("battle-frame:1") {
		t.Error("第 0 幀對 `battle-frame:1` 不該成立")
	}
	if holds("gate-bar") {
		t.Error("開場條還沒亮，`gate-bar` 不該成立")
	}

	// 推到條顯示中為止（這條 fixture 的城壁第一次挨打在第 148 幀附近）。
	for b.Frame < 400 && !b.Done {
		b.Step()
		if _, shown := b.StructureBar(); shown {
			break
		}
	}
	if _, shown := b.StructureBar(); !shown {
		t.Fatal("400 幀內城壁一次都沒挨打——條件的另一端驗不到，" +
			"這一支就只證明了「永遠 false」")
	}
	if !holds("gate-bar") {
		t.Error("條亮著，`gate-bar` 應該成立")
	}
	if !holds("battle-frame:1") {
		t.Error("推過 1 幀之後 `battle-frame:1` 應該成立")
	}
}
