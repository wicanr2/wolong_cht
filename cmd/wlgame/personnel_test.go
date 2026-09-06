package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 解任**不先過濾**：清單列的是全部的據點／勢力，不是「有派人的那些」
// （原版 `sub_16B08`／`sub_16BE3`，docs/spec/142）。
//
// ⚠ 先過濾會讓「那座城本來就沒人」這件事在畫面上消失——原版是選下去
// 才跳 TALK #54／#55 告訴你。
func TestDismissListsEveryTarget(t *testing.T) {
	w := &state.World{Player: 0}
	// 表是定長的：全部標成「不是玩家的」，只留三座城／三個勢力。
	for i := range w.Cities {
		w.Cities[i].Owner = 9
		w.Cities[i].Governor = NoOfficial
	}
	for i := 0; i < 3; i++ {
		w.Cities[i].Owner = 0 // 一個內政官都沒有
	}
	for i := range w.Factions {
		w.Factions[i].Alive = false
		w.Factions[i].Diplomat = NoOfficial
	}
	for i := 0; i < 3; i++ {
		w.Factions[i].Alive = true
	}
	const cities, factions = 3, 3

	g := &game{world: w, cmdCell: -1}
	g.removeGovernor()
	if g.list == nil {
		t.Fatal("一個內政官都沒有時清單還是要開（原版不過濾）")
	}
	if got := len(g.list.Rows); got != cities {
		t.Errorf("內政官解任列了 %d 座城，want %d（全部）", got, cities)
	}

	g2 := &game{world: w, cmdCell: -1}
	g2.removeDiplomat()
	if g2.list == nil {
		t.Fatal("一個外交官都沒有時清單還是要開")
	}
	// 候選是「活著且不是自己」⇒ 三個勢力扣掉玩家 ＝ 兩個。
	if got := len(g2.list.Rows); got != factions-1 {
		t.Errorf("外交官解任列了 %d 個勢力，want %d", got, factions-1)
	}
}

// 選到「那裡本來就沒人」的那一個：跳訊息，而且不寫壞資料
// （原版 `sub_16B4F` 是先寫 0xFF 再看舊值，等於寫了一次同樣的值）。
func TestDismissEmptySlotReportsNobody(t *testing.T) {
	w := &state.World{Player: 0}
	for i := range w.Cities {
		w.Cities[i].Owner = 9
		w.Cities[i].Governor = NoOfficial
	}
	w.Cities[0].Owner = 0

	g := &game{world: w, cmdCell: -1}
	g.removeGovernor()
	if g.list == nil {
		t.Fatal("清單沒開")
	}
	g.listPick(0)
	if w.Cities[0].Governor != NoOfficial {
		t.Errorf("內政官欄被寫成 %d，want %d", w.Cities[0].Governor, NoOfficial)
	}
	// 沒有 TALK.DAT 時 enqueueTalk 是 fail-closed，所以這裡驗的是「不會爆」
	// 與「資料沒被寫壞」；訊息本身的對拍在 docs/playtest/84。
}

// 選完**回到清單繼續選**，右鍵才離開（原版 `sub_16A9B`／`sub_16B08` 的
// `jmp` 迴圈，docs/spec/142 §1）。
func TestPersonnelFlowsLoopBackToTheList(t *testing.T) {
	w := &state.World{Player: 0}
	for i := range w.Cities {
		w.Cities[i].Owner = 9
		w.Cities[i].Governor = NoOfficial
	}
	w.Cities[0].Owner, w.Cities[1].Owner = 0, 0
	w.Cities[0].Governor = 5 // 這一座已經有人
	for i := range w.Generals {
		w.Generals[i].Alive = false
	}
	w.Generals[5] = state.General{Alive: true, Faction: 0}

	// 任命：選到「已經有人」的那一座 → 不關清單。
	g := &game{world: w, cmdCell: -1}
	g.pickCityForGovernor()
	if g.list == nil {
		t.Fatal("清單沒開")
	}
	if g.listPick(0) {
		t.Error("選到已有內政官的據點時關掉了清單，原版是回清單再選")
	}

	// 解任：選到沒人的那一座 → 也不關清單。
	g2 := &game{world: w, cmdCell: -1}
	g2.removeGovernor()
	if g2.listPick(1) {
		t.Error("解任選到沒人的據點時關掉了清單，原版是回清單再選")
	}
	// 解任成功 → 一樣不關。
	if g2.listPick(0) {
		t.Error("解任成功之後關掉了清單，原版是回清單再選")
	}
	if w.Cities[0].Governor != NoOfficial {
		t.Errorf("解任之後內政官欄 = %d，want %d", w.Cities[0].Governor, NoOfficial)
	}
}

// 那位官員說的一句是**八格一組**，由武將 +0x1E 選組內第幾個
// （原版 `sub_18810` 的 `ah`，docs/spec/142）。
func TestOfficialLineUsesVariantGroup(t *testing.T) {
	for _, tc := range []struct {
		base, variant, want int
	}{
		{governorAssignedTalk, 3, 457},  // 「遵命。」
		{governorAssignedTalk, 7, 461},  // 「我立刻前往。」
		{diplomatAssignedTalk, 3, 465},
		{governorDismissedTalk, 0, 502},
		{diplomatDismissedTalk, 0, 510},
	} {
		if got := resolveBattleTalkIndex(tc.base, tc.variant); got != tc.want {
			t.Errorf("組 %#x 變體 %d → #%d，want #%d",
				tc.base, tc.variant, got, tc.want)
		}
	}
}
