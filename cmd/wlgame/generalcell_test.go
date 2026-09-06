package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// 武將自陳兵種的四條分支（docs/spec/145 §1）。
// ⭐ 平手歸前面那一個——原版問的是「最大值等不等於攻城 → 等不等於野戰」，
// 所以三個一樣高時說城塞戰。只驗三個不同的樣本會漏掉這一格。
func TestGeneralAptitudeTalkBranches(t *testing.T) {
	newGame := func(apt [3]int, captor int) *game {
		w := &state.World{Player: 0}
		w.Generals[5].Alive = true
		w.Generals[5].Faction = 0
		w.Generals[5].Aptitude = apt
		w.Generals[5].Captor = captor
		return &game{world: w}
	}
	cases := []struct {
		name   string
		apt    [3]int
		captor int
		want   int
	}{
		{"攻城最高", [3]int{9, 4, 2}, NoOfficial, siegeBoastTalk},
		{"野戰最高", [3]int{2, 9, 4}, NoOfficial, fieldBoastTalk},
		{"水戰最高", [3]int{2, 4, 9}, NoOfficial, navalBoastTalk},
		{"三個一樣高 → 城塞戰", [3]int{6, 6, 6}, NoOfficial, siegeBoastTalk},
		{"攻城與野戰平手 → 城塞戰", [3]int{6, 6, 1}, NoOfficial, siegeBoastTalk},
		{"野戰與水戰平手 → 野戰", [3]int{1, 6, 6}, NoOfficial, fieldBoastTalk},
		{"俘虜壓過適性", [3]int{0, 0, 9}, 3, captiveRefusesTalk},
	}
	for _, c := range cases {
		if got := newGame(c.apt, c.captor).generalAptitudeTalk(5); got != c.want {
			t.Errorf("%s：組 %#x，want %#x", c.name, got, c.want)
		}
	}
}

// 勢力那一格選完之後鏡頭停在該勢力的首都，而且情報卡開著（§2）。
//
// ⚠ 狀態列要有 `TALK.DAT` 才掛得起來（`setStatusTalk` 查不到就靜靜清空），
// 所以這一支要載原版素材——不載的話那兩條斷言會**恆真地失敗**，
// 而不是恆真地通過。
func TestFactionCellFocusesCapital(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w := &state.World{Player: 0}
	for i := 0; i < 2; i++ {
		w.Factions[i].Alive = true
	}
	w.Factions[1].Capital = 7
	w.Cities[7].X, w.Cities[7].Y = 40, 30
	g := &game{lib: lib, world: w}
	g.openFactionList()

	if !g.statusBoxActive() {
		t.Error("勢力那一格沒有掛狀態列")
	}
	if g.list == nil {
		t.Fatal("勢力一覽沒開")
	}
	g.listPick(1)

	if g.camX != 40-centreCol || g.camY != 30-centreRow {
		t.Errorf("鏡頭 (%d,%d)，want (%d,%d)",
			g.camX, g.camY, 40-centreCol, 30-centreRow)
	}
	if !g.cityInfo.active {
		t.Error("首都的情報卡沒開")
	}
	if !g.statusBoxActive() {
		t.Error("情報卡還開著就把狀態列清掉了（原版是 sub_17E1F 回來才清）")
	}
}

// ⭐ 原版實跑錨點（docs/playtest/86 §2）：`root-saveb` 的第 1 槽，
// 點武將一覽裡的夏侯淵（編號 37，適性 0/4/2、說話類型 7），
// 原版說的是 **#573**「在原野上如果還有人比我強，我倒想和他一戰！」
// ——組 `0x1AA`（野戰）的第 7 格。
//
// 這一支把「挑哪一組」與「組內第幾格」兩段一起釘住；只驗前半段的話，
// 變體展開式改錯了不會紅。
func TestGeneralAptitudeTalkMatchesOriginalCapture(t *testing.T) {
	w, err := state.LoadScenario("../../workplace/dosgolem/root-saveb/SAVE.DAT", 0)
	if err != nil {
		t.Skipf("沒有那份原版存檔：%v", err)
	}
	g := &game{world: w}
	const xiahouYuan = 37
	gen := &w.Generals[xiahouYuan]
	if gen.Aptitude != [3]int{0, 4, 2} || gen.TalkVariant != 7 {
		t.Fatalf("樣本不對：適性 %v、說話類型 %d", gen.Aptitude, gen.TalkVariant)
	}
	base := g.generalAptitudeTalk(xiahouYuan)
	if base != fieldBoastTalk {
		t.Fatalf("組 = %#x，want %#x（野戰）", base, fieldBoastTalk)
	}
	if got := resolveBattleTalkIndex(base, gen.TalkVariant); got != 573 {
		t.Errorf("展開到 #%d，want #573", got)
	}
}
