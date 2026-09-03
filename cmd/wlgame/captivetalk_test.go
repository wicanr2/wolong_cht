package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func captiveGame(t *testing.T, lib *library.Library, player int) *game {
	t.Helper()
	w := &state.World{Player: player}
	w.Generals[7].Alive = true
	w.Generals[7].TalkVariant = 4
	w.Generals[7].Portrait = 12
	return &game{lib: lib, world: w}
}

func captiveEvent(fate combat.Fate) state.CorpsEvent {
	return state.CorpsEvent{
		Corps:     7,
		Destroyed: []int{7},
		Fate:      map[int]combat.Fate{7: fate},
		FateSides: map[int]state.FateSide{7: {Winner: 1, Loser: 0}},
	}
}

// 三種下場 × 敗方／勝方／局外人。**九種都要驗**——只驗一邊的話，
// 「永遠出同一則」也會通過（docs/spec/123 §1）。
//
// ⚠ 比對的字取自 **`workplace/orig/dosv` 的未校訂原文**（`library.Load`
// 讀的就是它）。#32／#34 有校訂（句首的殘字「配對」，docs/spec/123 §1.4），
// 所以這裡用的是校訂前的字——**換成讀語系包的話這幾條要跟著改**。
func TestCaptiveTalkPicksViewerSide(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name   string
		fate   combat.Fate
		player int
		want   []string // 每一則要出現的字；nil ＝ 不出訊息
	}{
		{"脫身·敗方", combat.Escaped, 0, []string{"逃過敵軍"}},
		{"脫身·勝方", combat.Escaped, 1, []string{"沒能將"}},
		{"脫身·局外", combat.Escaped, 2, nil},
		{"被擒·敗方", combat.Captured, 0, []string{"遭敵軍所擒"}},
		{"被擒·勝方", combat.Captured, 1, []string{"捉到", "放過我的話"}},
		{"被擒·局外", combat.Captured, 2, nil},
		{"自刎·勝方", combat.Suicide, 1, []string{"自刎"}},
		{"自刎·敗方", combat.Suicide, 0, nil}, // 舊主已滅，敗方沒有那一則
		{"自刎·局外", combat.Suicide, 2, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := captiveGame(t, lib, c.player)
			g.reportCaptives(captiveEvent(c.fate))
			if len(g.messages) != len(c.want) {
				t.Fatalf("排了 %d 則，want %d", len(g.messages), len(c.want))
			}
			for i, w := range c.want {
				if got := strings.Join(g.messages[i].lines, ""); !strings.Contains(got, w) {
					t.Errorf("第 %d 則 ＝ %q，應含 %q", i, got, w)
				}
			}
		})
	}
}

// 被擒那一組是 `0x19A`（438–445），**不是**敗走的 `0x198`（422–429）。
// 兩組形狀一樣、句子不同，拿錯不會炸，只會講別人的台詞。
func TestCaptiveRegretTalkIndex(t *testing.T) {
	for v := 0; v < 8; v++ {
		if got, want := resolveBattleTalkIndex(captiveRegretTalkBase, v), 438+v; got != want {
			t.Errorf("變體 %d ＝ %d，want %d", v, got, want)
		}
		if resolveBattleTalkIndex(captiveRegretTalkBase, v) ==
			resolveBattleTalkIndex(routRegretTalkBase, v) {
			t.Errorf("變體 %d 與敗走那一組撞了", v)
		}
	}
}

// 被擒之後武將的勢力已經是勝方——訊息要比的是舊主。
func TestCaptiveTalkUsesFormerLordNotCurrent(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	g := captiveGame(t, lib, 0)
	g.world.Generals[7].Faction = 1 // 已經改隸勝方（corpsPerishes 做的）
	g.reportCaptives(captiveEvent(combat.Captured))
	if len(g.messages) != 1 {
		t.Fatalf("排了 %d 則，want 1（敗方視角只有一則）", len(g.messages))
	}
	if got := strings.Join(g.messages[0].lines, ""); !strings.Contains(got, "遭敵軍所擒") {
		t.Errorf("拿到勝方那一則了：%q", got)
	}
}
