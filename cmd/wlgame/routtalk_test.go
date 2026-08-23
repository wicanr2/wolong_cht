package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 敗走的訊息只對玩家自己的軍團出（docs/spec/77 §1.1）。
// **兩個方向都要驗**——只驗一邊的話，「永遠出」與「永遠不出」都會通過。
func TestRoutTalkOnlyForPlayer(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	newGame := func(faction int) *game {
		w := &state.World{Player: 0}
		w.Corps[7].Faction = faction
		w.Generals[7].Alive = true
		w.Generals[7].Faction = faction
		w.Generals[7].TalkVariant = 4
		w.Generals[7].Portrait = 12
		return &game{lib: lib, world: w}
	}

	mine := newGame(0)
	mine.reportRout(state.CorpsEvent{Corps: 7, Routed: true})
	if len(mine.messages) != 1 {
		t.Fatalf("自己的軍團敗走排了 %d 則，want 1", len(mine.messages))
	}
	if got := strings.Join(mine.messages[0].lines, ""); !strings.Contains(got, "遭殲") {
		t.Errorf("排到的不是 #1F：%q", got)
	}

	theirs := newGame(1)
	theirs.reportRout(state.CorpsEvent{Corps: 7, Routed: true})
	if len(theirs.messages) != 0 {
		t.Errorf("別人的軍團敗走也出訊息了：%#v", theirs.messages)
	}
}

// 倒數歸零那一刻是**兩則**：一般通知 ＋ 主將自己的檢討句（docs/spec/77 §1.2）。
func TestRoutEndEnqueuesReturnAndRegret(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w := &state.World{Player: 0}
	w.Corps[7].Faction = 0
	w.Generals[7].Alive = true
	w.Generals[7].TalkVariant = 4
	w.Generals[7].Portrait = 12
	g := &game{lib: lib, world: w}

	g.reportRout(state.CorpsEvent{Corps: 7, RoutEnded: true})
	if len(g.messages) != 2 {
		t.Fatalf("排了 %d 則，want 2", len(g.messages))
	}
	if got := strings.Join(g.messages[0].lines, ""); !strings.Contains(got, "平安歸來") {
		t.Errorf("第一則不是 #23：%q", got)
	}
	if got := strings.Join(g.messages[1].lines, ""); !strings.Contains(got, "打回來") {
		t.Errorf("第二則不是變體 4（426）：%q", got)
	}
	if g.messages[0].portraitPage != -1 {
		t.Errorf("#23 應走一般通知肖像（al=93h），得到 %d", g.messages[0].portraitPage)
	}
	if g.messages[1].portraitPage != 12 {
		t.Errorf("檢討句應用主將的肖像，得到 %d", g.messages[1].portraitPage)
	}
}

// 組編號 `0x198` 展開成 422–429。
//
// ⚠ 內政官那一組的組編號 `0x1A6` 十進位剛好是 422 ——**把組編號當索引
// 用會落到這一組的第 0 格**「．．．．」，而那句話語意上也講得通，
// 所以錯了不會被發現。這條測試就是為了擋那個。
func TestRoutRegretTalkIndex(t *testing.T) {
	for variant, want := range []int{422, 423, 424, 425, 426, 427, 428, 429} {
		if got := resolveBattleTalkIndex(routRegretTalkBase, variant); got != want {
			t.Errorf("變體 %d ⇒ %d，要 %d", variant, got, want)
		}
	}
	// ⭐ 這一組的第 0 格（422）**數值上就等於** `governorRegretTalkBase`
	// 的十進位。兩組要靠算式分開，不能靠看數字。
	if resolveBattleTalkIndex(routRegretTalkBase, 0) !=
		int(governorRegretTalkBase) {
		t.Error("前提變了：422 不再同時是敗走組的第 0 格與內政官的組編號")
	}
	if resolveBattleTalkIndex(routRegretTalkBase, 0) ==
		resolveBattleTalkIndex(governorRegretTalkBase, 0) {
		t.Error("敗走與內政官展開後落到同一句")
	}
}
