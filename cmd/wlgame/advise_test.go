package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// sub_13D91／sub_13DC9 對 byte_10D00 在 0 與 255 飽和；GUI 的直接反應與
// Session 理由路徑都必須共用這個邊界。
func TestAdjustTrustSaturatesOriginalByte(t *testing.T) {
	for _, tc := range []struct {
		name  string
		start int
		delta int
		want  int
	}{
		{"increase", 250, 20, 255},
		{"decrease", 5, -20, 0},
		{"middle", 100, -20, 80},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := &game{world: &state.World{Trust: tc.start}}
			g.adjustTrust(tc.delta)
			if g.world.Trust != tc.want {
				t.Fatalf("Trust %d %+d → %d, want %d",
					tc.start, tc.delta, g.world.Trust, tc.want)
			}
		})
	}
}

// 三個進言指令的 TALK 起點（docs/spec/44 §1）。
func TestAdviseTalkBasesMatchOriginal(t *testing.T) {
	for _, tc := range []struct {
		cmd  persuasion.Command
		want int
	}{
		{persuasion.Hostility, 86},
		{persuasion.CeaseFire, 150},
		{persuasion.Cooperate, 214},
	} {
		if got := adviseTalkBase(tc.cmd); got != tc.want {
			t.Errorf("%v 的起點 = %d，原版是 %d", tc.cmd, got, tc.want)
		}
	}
}

// 君主回答的索引 ＝ base + 4 + 結果碼×3（＋說話型），碼 ≥ 4 用 83。
func TestAdviseReplyIndexFollowsReactionCode(t *testing.T) {
	const base = 86
	for _, tc := range []struct {
		r    persuasion.Reaction
		want int
	}{
		{persuasion.Refuse, 90},
		{persuasion.Agree, 93},
		{persuasion.AskReason, 96},
		{persuasion.AlreadyAtWar, 99},
		{persuasion.SameFaction, 83},
	} {
		if got := adviseReplyIndex(base, tc.r, 0); got != tc.want {
			t.Errorf("碼 %d 的索引 = %d，原版是 %d", tc.r, got, tc.want)
		}
	}
	// 說話型直接加進索引（sub_13C99 的 add cx, ax）。
	if got := adviseReplyIndex(base, persuasion.Refuse, 2); got != 92 {
		t.Errorf("說話型 2 的索引 = %d，want 92", got)
	}
}

// 三個指令拿到的是**不同**的原文，而且展開得出對象的名字。
func TestAdviseLinesComeFromTalkDat(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	g := &game{lib: lib, world: w}
	g.target = 1

	seen := map[string]persuasion.Command{}
	for _, c := range adviseCommands {
		line := g.adviseLine("軍師", adviseTalkBase(c)+3)
		if line == "" {
			t.Fatalf("%v 的進言句是空的", c)
		}
		if prev, dup := seen[line]; dup {
			t.Errorf("%v 與 %v 拿到同一句：%q", c, prev, line)
		}
		seen[line] = c
		if !strings.Contains(line, big5(w.LordName(1))) {
			t.Errorf("%v 的進言句沒有展開 {3}：%q", c, line)
		}
	}
}

// 說服迴圈的索引算式（docs/spec/44 §5）。原版在進迴圈前 `add [bp+0], 10h`，
// 所以整段相對於 base+16，而那一格正好是五選一的選單。
func TestAdviseReasonIndicesFollowOriginalArithmetic(t *testing.T) {
	if got := adviseReasonBase(persuasion.Hostility); got != 102 {
		t.Fatalf("敵對的說服起點 = %d，要 102", got)
	}
	base := adviseReasonBase(persuasion.Hostility)
	for slot, want := range []int{103, 104, 105, 106, 107} {
		if got := base + slot + 1; got != want {
			t.Errorf("第 %d 列的軍師句 = %d，要 %d", slot, got, want)
		}
	}
	for slot, row := range [][3]int{
		{108, 111, 114}, {117, 120, 123}, {126, 129, 132}, {135, 138, 141},
	} {
		for code, out := range []persuasion.Outcome{
			persuasion.Failed, persuasion.Agreed, persuasion.Continue,
		} {
			if got := adviseReasonReply(base, slot, out, false, 0); got != row[code] {
				t.Errorf("第 %d 列的 %v = %d，要 %d", slot, out, got, row[code])
			}
		}
	}
	if got := adviseReasonReply(base, 4, persuasion.Withdrawn, false, 0); got != 144 {
		t.Errorf("撤回 = %d，要 144", got)
	}
	if got := adviseReasonReply(base, 1, persuasion.Continue, true, 0); got != 147 {
		t.Errorf("重複同一個理由 = %d，要 147", got)
	}
	// 說話型變體佔三則，所以每個位置的下一則就是變體 1。
	if got := adviseReasonReply(base, 0, persuasion.Failed, false, 2); got != 110 {
		t.Errorf("變體 2 = %d，要 110", got)
	}
	// ⭐ 一個指令佔 16 + 48 則，所以下一個指令的起點正好接在後面。
	// 這一條同時把 §2 的 16 與 §5 的 48 釘住。
	for _, tc := range [][2]persuasion.Command{
		{persuasion.Hostility, persuasion.CeaseFire},
		{persuasion.CeaseFire, persuasion.Cooperate},
	} {
		if got := adviseReasonBase(tc[0]) + 48; got != adviseTalkBase(tc[1]) {
			t.Errorf("%v 之後接 %v：%d ≠ %d",
				tc[0], tc[1], got, adviseTalkBase(tc[1]))
		}
	}
}

// 撤回那一列的軍師句是「既然不和主公之意…」，位置在選單第 5 列。
func TestAdviseReasonSlotPutsWithdrawLast(t *testing.T) {
	for _, c := range adviseCommands {
		if got := adviseReasonSlot(c, persuasion.Withdraw); got != 4 {
			t.Errorf("%v 的撤回排第 %d，要第 4", c, got)
		}
		for i, r := range persuasion.Options(c) {
			if got := adviseReasonSlot(c, r); got != i {
				t.Errorf("%v 的 %v 排第 %d，要第 %d", c, r, got, i)
			}
		}
	}
}

// 選單那五列直接取自 TALK.DAT，不是寫在 Go 原始碼裡的字串（CLAUDE.md §6）。
func TestAdviseReasonLabelsComeFromTalkMenu(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	g := &game{lib: lib, world: w}
	g.target = 1

	seen := map[string]persuasion.Command{}
	for _, c := range adviseCommands {
		labels := g.adviseReasonLabels(c)
		if len(labels) != len(persuasion.Options(c)) {
			t.Fatalf("%v 的選單 %d 列，要 %d 列",
				c, len(labels), len(persuasion.Options(c)))
		}
		if labels[4] != "撤回進言" {
			t.Errorf("%v 的第 5 列是 %q，要「撤回進言」", c, labels[4])
		}
		for _, l := range labels[:4] {
			if strings.HasSuffix(l, "　") || l == "" {
				t.Errorf("%v 的選單列沒有去掉尾端全形空白：%q", c, l)
			}
		}
		key := strings.Join(labels[:4], "/")
		if prev, dup := seen[key]; dup {
			t.Errorf("%v 與 %v 的選單一樣：%s", c, prev, key)
		}
		seen[key] = c
	}
}
