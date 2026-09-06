package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 原版 `sub_18853` 有 27 個帶索引的呼叫點（docs/spec/140 §1.1）。
// 這一支把「哪一步掛哪一則」逐格釘住——常數寫錯一個就紅。
//
// ⭐ 只驗常數不驗流程是不夠的：**先前接了的三條也是常數對、流程沒接**。
// 所以下面那幾支是真的走進流程再看框裡是哪一則。
func TestStatusTalkIndexesMatchOriginal(t *testing.T) {
	cases := []struct {
		name string
		got  int
		want int
	}{
		{"#0 編成選武將 sub_16C5E", formStatusTalk, 0},
		{"#1 編成指示 sub_16C92", formOrderTalk, 1},
		{"#2 行軍選軍團 sub_1628F", marchCorpsTalk, 2},
		{"#3 行軍目標 sub_17FDB", marchTargetTalk, 3},
		{"#4 軍團情報 sub_17F90", corpsInfoTalk, 4},
		{"#5 敵對提案 sub_16405", adviseHostilityTalk, 5},
		{"#6 停戰提案 sub_164F1", adviseCeaseFireTalk, 6},
		{"#7 協同進攻 sub_16623 第二步", adviseTargetTalk, 7},
		{"#8 請求協助 sub_16623 第一步", adviseAllyTalk, 8},
		{"#9 任命選武將", pickOfficialTalk, 9},
		{"#11 內政官任命", governorAssignTalk, 11},
		{"#12 外交官任命", diplomatAssignTalk, 12},
		{"#13 內政官解任", governorRemoveTalk, 13},
		{"#14 外交官解任", diplomatRemoveTalk, 14},
		{"#15 遷都 sub_16909", adviseRelocateTalk, 15},
		{"#16 財政 sub_1678D", financeStatusTalk, 16},
		{"#17 稅率 sub_167CD", financeAmountTalk[0], 17},
		{"#18 騎兵 sub_167E6", financeAmountTalk[1], 18},
		{"#19 弓兵 sub_16806", financeAmountTalk[2], 19},
		{"#20 步兵 sub_16826", financeAmountTalk[3], 20},
		{"#21 行軍三選一 sub_17FDB", marchModePromptTalk, 21},
		{"#22 位置確認 sub_1628F", locateCorpsTalk, 22},
		{"#23 據點一覽 sub_162FB", cityListStatusTalk, 23},
		{"#24 武將 sub_16366", generalCellTalk, 24},
		{"#25 勢力 sub_163BF", factionCellTalk, 25},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s：接的是 #%d，want #%d", c.name, c.got, c.want)
		}
	}
}

// 請求協助是**兩步兩則**：`sub_16623` 先 `cx = 8` 選要拜託誰，
// 選完才換 `cx = 7` 選要一起打誰，取消退回上一步又換回 #8。
//
// ⭐ 這一支**真的走流程再看框裡的字**——只驗常數的話，
// 「常數對、流程沒接」這種錯誤一個都擋不住（先前那三條就是這樣）。
func TestAdviseCooperateSwapsStatusBetweenSteps(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	g := &game{lib: lib, world: w}
	want := func(index int) []string {
		lines, ok := g.talkLines(index, nil)
		if !ok {
			t.Fatalf("讀不到 TALK #%d", index)
		}
		return lines
	}
	same := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	g.openAdvise()
	g.pickAdviseCommand(2) // 請求協助
	if g.advise != advisePickAlly {
		t.Fatalf("第 2 項應該先選協助勢力，現在是 %v", g.advise)
	}
	if !same(g.statusBox.lines, want(adviseAllyTalk)) {
		t.Errorf("第一步掛的不是 #8：%q", g.statusBox.lines)
	}

	g.advise = advisePickTarget
	g.setAdviseStatus()
	if !same(g.statusBox.lines, want(adviseTargetTalk)) {
		t.Errorf("第二步掛的不是 #7：%q", g.statusBox.lines)
	}

	// 退回第一步要換回 #8。
	g.advise = advisePickAlly
	g.setAdviseStatus()
	if !same(g.statusBox.lines, want(adviseAllyTalk)) {
		t.Errorf("退回第一步沒有換回 #8：%q", g.statusBox.lines)
	}

	// 敵對／停戰只有一則，不會因為階段而變。
	for _, c := range []struct {
		row  int
		talk int
	}{{0, adviseHostilityTalk}, {1, adviseCeaseFireTalk}} {
		g.openAdvise()
		g.pickAdviseCommand(c.row)
		if !same(g.statusBox.lines, want(c.talk)) {
			t.Errorf("第 %d 項掛的不是 #%d：%q", c.row, c.talk, g.statusBox.lines)
		}
	}
}
