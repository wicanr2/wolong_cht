package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
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
		if got := persuasion.TalkBase(tc.cmd); got != tc.want {
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
		if got := persuasion.TalkReplyIndex(base, tc.r, 0); got != tc.want {
			t.Errorf("碼 %d 的索引 = %d，原版是 %d", tc.r, got, tc.want)
		}
	}
	// 說話型直接加進索引（sub_13C99 的 add cx, ax）。
	if got := persuasion.TalkReplyIndex(base, persuasion.Refuse, 2); got != 92 {
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
		g.clearAdviseBoxes()
		g.adviseSay(adviseAdvisor, persuasion.TalkBase(c)+3)
		g.adviseAdvance()
		line := strings.Join(g.adviseAdvisorSaid, "")
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
	if got := persuasion.TalkReasonBase(persuasion.Hostility); got != 102 {
		t.Fatalf("敵對的說服起點 = %d，要 102", got)
	}
	base := persuasion.TalkReasonBase(persuasion.Hostility)
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
			if got := persuasion.TalkReasonReply(base, slot, out, false, 0); got != row[code] {
				t.Errorf("第 %d 列的 %v = %d，要 %d", slot, out, got, row[code])
			}
		}
	}
	if got := persuasion.TalkReasonReply(base, 4, persuasion.Withdrawn, false, 0); got != 144 {
		t.Errorf("撤回 = %d，要 144", got)
	}
	if got := persuasion.TalkReasonReply(base, 1, persuasion.Continue, true, 0); got != 147 {
		t.Errorf("重複同一個理由 = %d，要 147", got)
	}
	// 說話型變體佔三則，所以每個位置的下一則就是變體 1。
	if got := persuasion.TalkReasonReply(base, 0, persuasion.Failed, false, 2); got != 110 {
		t.Errorf("變體 2 = %d，要 110", got)
	}
	// ⭐ 一個指令佔 16 + 48 則，所以下一個指令的起點正好接在後面。
	// 這一條同時把 §2 的 16 與 §5 的 48 釘住。
	for _, tc := range [][2]persuasion.Command{
		{persuasion.Hostility, persuasion.CeaseFire},
		{persuasion.CeaseFire, persuasion.Cooperate},
	} {
		if got := persuasion.TalkReasonBase(tc[0]) + 48; got != persuasion.TalkBase(tc[1]) {
			t.Errorf("%v 之後接 %v：%d ≠ %d",
				tc[0], tc[1], got, persuasion.TalkBase(tc[1]))
		}
	}
}

// 撤回那一列的軍師句是「既然不和主公之意…」，位置在選單第 5 列。
func TestAdviseReasonSlotPutsWithdrawLast(t *testing.T) {
	for _, c := range adviseCommands {
		if got := persuasion.TalkReasonSlot(c, persuasion.Withdraw); got != 4 {
			t.Errorf("%v 的撤回排第 %d，要第 4", c, got)
		}
		for i, r := range persuasion.Options(c) {
			if got := persuasion.TalkReasonSlot(c, r); got != i {
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

// 選單框：位置是立即值，**大小由內容算**（docs/spec/45 §2）。
func TestChoiceBoxMatchesOriginalRect(t *testing.T) {
	for _, tc := range []struct {
		name      string
		got, want int
	}{
		{"X（sub_193E9 的 dl=5 ⇒ 5×16）", talkChoiceX, 80},
		{"Y（dh=0Bh ⇒ 11×16）", talkChoiceY, 176},
		{"進言選單 X（sub_16224 的 dx=400h）", adviseMenuX, 0},
		{"進言選單 Y", adviseMenuY, 64},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %d，要 %d", tc.name, tc.got, tc.want)
		}
	}
	// 五列 6 個全形字（#102／#77 都是這個形狀）⇒ 112 × 96。
	five := []string{"外交關係惡劣", "我國較有利　", "敵正侵攻他國", "敵勢力疲乏　", "撤回進言　　"}
	if _, _, w, h := legacyChoiceRect(80, 176, five); w != 112 || h != 96 {
		t.Errorf("五列選單 = %d×%d，要 112×96", w, h)
	}
	// 三列 5 個字（#363）⇒ 96 × 64。**列少了框要跟著縮。**
	three := []string{"無條件同意", "提供資金　", "拒　　絕　"}
	if _, _, w, h := legacyChoiceRect(80, 176, three); w != 96 || h != 64 {
		t.Errorf("三列選單 = %d×%d，要 96×64", w, h)
	}
	// ⚠ 原版只看第一列，remake 取最大值（text.Decode 會砍掉行尾的
	// 全形空白，拿不回 padding）。四則實際的選單兩種算法同值，
	// 只有這種人造的參差資料分得出來。
	ragged := []string{"拒　　絕", "很長很長的一列文字"}
	if _, _, w, _ := legacyChoiceRect(0, 0, ragged); w != 10*16 {
		t.Errorf("寬 = %d，取最大值的話要 %d", w, 10*16)
	}
	if _, _, w, h := legacyChoiceRect(0, 0, nil); w != 32 || h != 16 {
		t.Errorf("空選單 = %d×%d，要 32×16（不要算出 0 或負數）", w, h)
	}
}

// 誰說話就進誰的框（docs/spec/45 §1）：君主 → 上框、軍師 → 下框，
// 而且**兩個框各自只換掉最新的那一句**。
func TestAdviseSceneLinesGoToTheRightBox(t *testing.T) {
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
	g.adviseCmd = persuasion.Hostility

	base := persuasion.TalkBase(g.adviseCmd)
	g.clearAdviseBoxes()
	g.adviseSay(adviseLord, base) // #86「{4}啊，是怎麼了？」
	g.adviseSay(adviseAdvisor, base+3)
	// ⭐ 排隊之後畫面上還是空的——**原版每一句都要等按鍵**。
	if len(g.adviseLordSaid) != 0 || len(g.adviseAdvisorSaid) != 0 {
		t.Fatal("還沒推進就把句子寫進框裡了")
	}
	if !g.adviseAdvance() || !g.adviseAdvance() {
		t.Fatal("排了兩句卻推不動兩次")
	}
	if g.adviseTalking() {
		t.Error("兩句都演完了還說在講話")
	}
	lord, advisor := strings.Join(g.adviseLordSaid, ""), strings.Join(g.adviseAdvisorSaid, "")
	if lord == "" || advisor == "" {
		t.Fatalf("兩個框要各有一句：上 %q 下 %q", lord, advisor)
	}
	if lord == advisor {
		t.Errorf("兩個框拿到同一句：%q", lord)
	}
	// 君主再說一句：上框換掉，下框不動。
	g.adviseSay(adviseLord, persuasion.TalkReplyIndex(base, persuasion.Agree, 0))
	g.adviseAdvance()
	if got := strings.Join(g.adviseLordSaid, ""); got == lord {
		t.Errorf("上框沒有換句：%q", got)
	}
	if got := strings.Join(g.adviseAdvisorSaid, ""); got != advisor {
		t.Errorf("下框被動到了：%q ≠ %q", got, advisor)
	}
	// 索引查不到時保留原句，不顯示半句（fail-closed）。
	g.adviseSay(adviseLord, 1<<20)
	g.adviseAdvance()
	if len(g.adviseLordSaid) == 0 {
		t.Error("查不到的索引把上框清空了")
	}
	// 每列都在框裡放得下。
	for _, line := range append(append([]string{}, g.adviseLordSaid...), g.adviseAdvisorSaid...) {
		if w := textdraw.StringWidth(line); w > talkTextWidth {
			t.Errorf("%q 寬 %d，超過框內 %d px", line, w, talkTextWidth)
		}
	}
}

// 第四、五項的六個位置（docs/spec/49 §1）。
func TestVerdictTalkIndicesMatchOriginal(t *testing.T) {
	for _, tc := range []struct {
		name                            string
		base, open, advisor, yes, no    int
	}{
		{"遷都", adviseRelocateTalkBase, 386, 389, 390, 393},
		{"請求出陣", adviseSortieTalkBase, 396, 399, 400, 403},
	} {
		if tc.base != tc.open {
			t.Errorf("%s 的起點 = %d，要 %d", tc.name, tc.base, tc.open)
		}
		if got := tc.base + 3; got != tc.advisor {
			t.Errorf("%s 的軍師句 = %d，要 %d", tc.name, got, tc.advisor)
		}
		if got := tc.base + 4; got != tc.yes {
			t.Errorf("%s 的接受句 = %d，要 %d", tc.name, got, tc.yes)
		}
		if got := tc.base + 7; got != tc.no {
			t.Errorf("%s 的拒絕句 = %d，要 %d", tc.name, got, tc.no)
		}
	}
	// 兩項共用 `sub_13B08`，所以區塊等寬（10 則）。
	if adviseSortieTalkBase-adviseRelocateTalkBase != 10 {
		t.Errorf("兩個起點差 %d，要 10", adviseSortieTalkBase-adviseRelocateTalkBase)
	}
}

// 第四、五項在選單上排第 4、5，而且**沒有說服迴圈**。
func TestAdviseMenuHasFiveRows(t *testing.T) {
	if got := len(adviseFallbackNames); got != 5 {
		t.Fatalf("進言選單 %d 列，要 5", got)
	}
	if adviseRelocateRow != len(adviseCommands) || adviseSortieRow != len(adviseCommands)+1 {
		t.Errorf("第四、五項的列號不對：%d／%d", adviseRelocateRow, adviseSortieRow)
	}
}

// 君主定案那三句進對的框（上／下／上），而且接受與拒絕拿到不同的第三句。
func TestVerdictLinesGoToTheRightBoxes(t *testing.T) {
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
	g.target = -1

	g.sayVerdict(adviseSortieTalkBase, true)
	if g.advise != adviseVerdict {
		t.Errorf("階段 = %d，要 adviseVerdict", g.advise)
	}
	// sayVerdict 排三句、先演第一句；剩下兩句要按鍵才會出來。
	if !g.adviseTalking() {
		t.Fatal("三句應該還沒演完")
	}
	for g.adviseAdvance() {
	}
	yes := strings.Join(g.adviseLordSaid, "")
	advisor := strings.Join(g.adviseAdvisorSaid, "")
	if yes == "" || advisor == "" {
		t.Fatalf("兩個框要各有一句：上 %q 下 %q", yes, advisor)
	}
	g.sayVerdict(adviseSortieTalkBase, false)
	for g.adviseAdvance() {
	}
	if no := strings.Join(g.adviseLordSaid, ""); no == yes {
		t.Errorf("接受與拒絕拿到同一句：%q", no)
	}
	if got := strings.Join(g.adviseAdvisorSaid, ""); got != advisor {
		t.Errorf("軍師那一句不該跟著結果變：%q ≠ %q", got, advisor)
	}
}

// 進言選單的五列取自 TALK #77（`sub_16224` 的 `cx = 4Dh`），
// **不是拿內部術語頂替**——原版第三項寫「請求協助」不是「協力要請」。
func TestAdviseCommandLabelsComeFromTalk(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	g := &game{lib: lib}
	labels := g.adviseCommandLabels()
	if len(labels) != 5 {
		t.Fatalf("%d 列，要 5", len(labels))
	}
	for _, want := range []struct {
		row  int
		text string
		// 行首的全形空白**照原樣保留**（行尾的被 text.Decode 砍掉了）。
	}{{2, "\u3000請求協助"}, {3, "\u3000遷\u3000\u3000都"}, {4, "請求君主出陣"}} {
		if labels[want.row] != want.text {
			t.Errorf("第 %d 列 = %q，要 %q", want.row+1, labels[want.row], want.text)
		}
	}
	if labels[2] == persuasion.Cooperate.String() {
		t.Error("第三列用了內部術語，不是原版選單的用字")
	}
	// 五列等寬 ⇒ 框寬 (6+1)×16；最長的一列（6 個字）放得進去。
	if _, _, w, _ := legacyChoiceRect(adviseMenuX, adviseMenuY, labels); w != 112 {
		t.Errorf("框寬 = %d，要 112——padding 被 trim 掉了嗎？", w)
	}
	for _, l := range labels {
		if textdraw.StringWidth(l) > 112-2*8 {
			t.Errorf("%q 寬 %d，超出框內 %d px", l, textdraw.StringWidth(l), 112-16)
		}
	}
}

// AskReason 也要演君主的回答，而且是**四個反應碼共用的那個算式**
// （docs/spec/108）。先前只有 default 那一支叫得到 TalkReplyIndex，
// 於是最常走的那條路少演一句。
func TestAskReasonPlaysLordReplyBeforeMenu(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	// 劇本一的呂布（勢力 13）對曹操（勢力 0）：四個理由都不成立，
	// 而第一反應是 AskReason —— 正是實機拍到的那個局面
	// （docs/playtest/56 §4.3）。
	w.Player = 13
	g := &game{lib: lib, world: w}
	g.adviseCmd, g.target = persuasion.Hostility, 0

	s := g.situation(g.target)
	if got := persuasion.FirstReaction(persuasion.Hostility, s, false); got != persuasion.AskReason {
		t.Fatalf("第一反應 ＝ %v，這個測試要的是 AskReason", got)
	}

	g.beginPersuasion()
	if g.sess == nil {
		t.Fatal("AskReason 應該開說服迴圈")
	}

	// 佇列裡君主那一側要有兩句：開場 ＋ 回答。演完之後上框停在回答。
	for g.adviseTalking() {
		g.adviseAdvance()
	}
	base := persuasion.TalkBase(persuasion.Hostility)
	want := persuasion.TalkReplyIndex(base, persuasion.AskReason, g.playerTalkVariant())
	lines, ok := g.talkLines(want, nil)
	if !ok || len(lines) == 0 {
		t.Fatalf("取不到 TALK #%d", want)
	}
	got := strings.Join(g.adviseLordSaid, "")
	if wantText := strings.Join(lines, ""); got != wantText {
		t.Errorf("上框 ＝ %q，want %q（TALK #%d）", got, wantText, want)
	}
}

// 君主帶著軍團的時候，進言整個開不起來（docs/spec/111）。
// **兩個方向都要驗**——只驗擋下來那一邊的話，把 openAdvise 寫死成
// 「永遠不開」也會通過。
func TestLordLeadsCorpsBlocksAdvise(t *testing.T) {
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	w.Player = 0
	g := &game{world: w}

	g.openAdvise()
	if !g.adviseActive() {
		t.Fatal("君主在朝堂上卻開不了進言")
	}
	g.closeAdvise()

	// 手動把君主編成軍團長（remake 差異，docs/spec/76）。
	f := &w.Factions[w.Player]
	f.Reserves = [3]int{600, 600, 600}
	manned := [army.Positions]bool{}
	manned[0] = true
	if err := w.FormCorps(f.Lord, [army.Positions]army.TroopType{}, manned); err != nil {
		t.Fatalf("君主編成失敗：%v", err)
	}
	if !w.LordLeadsCorps() {
		t.Fatal("君主帶兵了，判準卻沒認出來")
	}

	g.lastEvent = ""
	g.openAdvise()
	if g.adviseActive() {
		t.Error("君主在領軍，進言還開得起來")
	}
	if g.lastEvent != adviseLordAwayEvent {
		t.Errorf("事件列 ＝ %q，want %q", g.lastEvent, adviseLordAwayEvent)
	}
}
