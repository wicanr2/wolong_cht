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
