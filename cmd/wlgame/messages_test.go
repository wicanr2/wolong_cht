package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

func TestMessagePagePreservesTalkHardBoundaries(t *testing.T) {
	// 框高 80 px、上下各內縮 8、一列 16 px ⇒ 4 列（docs/spec/41 §1）。
	if messagePageRows != 4 || talkLinePitch != 16 {
		t.Fatalf("原版 TALK page contract = rows %d／pitch %d，want 4／16", messagePageRows, talkLinePitch)
	}
	lines := make([]string, messagePageRows+2)
	for i := range lines {
		lines[i] = string(rune('A' + i))
	}

	first, pages, ok := messagePage(lines, 0)
	if !ok || pages != 2 || len(first) != messagePageRows {
		t.Fatalf("第一頁 = (%#v,%d,%v)，want %d 行／2 頁／true", first, pages, ok, messagePageRows)
	}
	if first[0] != "A" || first[len(first)-1] != string(rune('A'+messagePageRows-1)) {
		t.Fatalf("第一頁改動 TALK 硬斷行：%#v", first)
	}
	last, gotPages, ok := messagePage(lines, 1)
	if !ok || gotPages != pages || len(last) != 2 || last[0] != string(rune('A'+messagePageRows)) {
		t.Fatalf("第二頁 = (%#v,%d,%v)，未保留尾端硬斷行", last, gotPages, ok)
	}
}

func TestTalkLinesDropsOnlyStructuralTrailingBlank(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	g := &game{lib: lib}
	lines, ok := g.talkLines(0x168, map[byte]string{'3': "劉璋", '6': ""})
	if !ok || len(lines) != 2 || lines[len(lines)-1] == "" {
		t.Fatalf("TALK #360 結構尾空行未移除：ok=%v lines=%#v", ok, lines)
	}

	// #365 的第一行包含 \6 formatter；它是行內零寬控制，不是結構尾行。
	lines, ok = g.talkLines(0x16D, map[byte]string{'3': "劉璋", '6': "", '7': "500"})
	if !ok || len(lines) != 2 || lines[0] == "" || lines[1] == "" {
		t.Fatalf("TALK 中間行被錯誤移除：ok=%v lines=%#v", ok, lines)
	}
}

func TestLayoutMessageLinesUsesMeasuredWidthAndPreservesHardLines(t *testing.T) {
	lines := layoutMessageLines([]string{
		"甲乙丙丁戊己庚辛壬癸子丑寅卯辰巳午未申酉戌亥甲乙",
		"第二行A2，混排。",
	})
	if len(lines) <= 2 {
		t.Fatalf("密集繁中沒有依實際寬度換行：%#v", lines)
	}
	for i, line := range lines {
		if textdraw.StringWidth(line) > messageContentWidth+textdraw.GlyphW {
			t.Fatalf("第 %d 列超出 modal 寬度：%q (%d)", i, line, textdraw.StringWidth(line))
		}
	}
	if lines[len(lines)-1] == "" {
		t.Fatal("換行不應在尾端製造空列")
	}
}

func TestLayoutMessageLinesKeepsExplicitBlankLine(t *testing.T) {
	got := layoutMessageLines([]string{"甲", "", "乙"})
	if len(got) != 3 || got[0] != "甲" || got[1] != "" || got[2] != "乙" {
		t.Fatalf("TALK hard blank line 遺失：%#v", got)
	}
}

func TestMessagePageClampsInvalidPage(t *testing.T) {
	lines := []string{"甲", "乙"}
	got, pages, ok := messagePage(lines, 99)
	if !ok || pages != 1 || len(got) != 2 || got[0] != "甲" || got[1] != "乙" {
		t.Fatalf("超界頁碼未夾回最後一頁：%#v %d %v", got, pages, ok)
	}
}

func TestReleasedGeneralTalkOnlyTargetsPlayerFaction(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	general := -1
	for i, g := range w.Generals {
		if g.Alive && i != w.Factions[w.Player].Lord {
			general = i
			break
		}
	}
	if general < 0 {
		t.Fatal("找不到可測試的釋放武將")
	}
	g := &game{lib: lib, world: w}
	ev := state.Event{ReleasedGenerals: []int{general}}
	w.Generals[general].Faction = w.Player
	g.enqueueEventMessages(ev)
	if len(g.messages) != 1 || len(g.messages[0].lines) == 0 {
		t.Fatalf("玩家勢力釋放武將沒有產生 #37 通知：%#v", g.messages)
	}

	g.messages = nil
	w.Generals[general].Faction = (w.Player + 1) % len(w.Factions)
	g.enqueueEventMessages(ev)
	if len(g.messages) != 0 {
		t.Fatalf("非玩家勢力釋放武將不應產生 #37：%#v", g.messages)
	}
}

func TestReleasedGeneralRawFollowup409IsEmptyNoOp(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	visible := false
	for _, line := range lib.Talk.Messages[0x199].Lines {
		for _, part := range line.Parts {
			if part.Marker != 0 || len(part.Raw) != 0 {
				visible = true
			}
		}
	}
	if visible {
		t.Fatal("原版事件 9 後續 #409 應只有資料上的空行，不應含文字或 marker")
	}

	// sub_150D7 在 #37 後仍呼叫 CX=199h；空槽不應產生空白 modal。
	g := &game{lib: lib}
	g.enqueueTalk(0x199, map[byte]string{'1': "武將", '2': "據點", '3': "君主", '4': "軍師", '6': "", '7': "0"})
	if len(g.messages) != 0 {
		t.Fatalf("空槽 #409 不應排入訊息：%#v", g.messages)
	}
}

func TestSecondaryTalkUsesCapturedRawFormatterWord(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	g := &game{lib: lib, world: w}
	raw, ok := w.ResolveTalkFormatter2(0)
	if !ok {
		t.Fatal("原版 word 0 應能在載入的 DOS/V 動態區解析")
	}
	// 這是呈現層對「未來由 DOS 動態 trace 捕捉的 raw word」的通用驗證。
	// 事件 6 的第二次 sub_13C3D 呼叫讀的是 SS:[DI]，目前沒有合法動態
	// trace 可捕捉它；因此事件 handler 不得將此 synthetic word 0 當作
	// 原版 payload。這裡只驗證已捕捉的 raw DS bytes 能走同一個 Big5
	// 呈現路徑，而不是只測到 modal 被排入。
	g.enqueueTalkNotice(state.TalkNotice{
		Index: 0x48, City: -1, Faction: -1, General: -1, Amount: -1,
		RawFormatterWord: 0, RawFormatterWordValid: true, Secondary: true,
	})
	if len(g.messages) != 1 || !strings.Contains(strings.Join(g.messages[0].lines, "\n"), textDecodeBig5(raw)) {
		t.Fatalf("原版 \\2 raw payload 未進入 #72：raw=%q messages=%#v", textDecodeBig5(raw), g.messages)
	}

	// 沒有顯式 raw word 時，不能把零值誤當成 word 0，也不能猜據點。
	g.messages = nil
	g.enqueueTalkNotice(state.TalkNotice{
		Index: 0x48, City: -1, Faction: -1, General: -1, Amount: -1, Secondary: true,
	})
	if len(g.messages) != 0 {
		t.Fatalf("缺少原版 \\2 payload 的 #72 不應顯示半句：%#v", g.messages)
	}
}

func TestSecondaryTalk76UsesNoPortraitMode(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Fatal(err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Fatal(err)
	}
	g := &game{lib: lib, world: w}
	g.enqueueTalkNotice(state.TalkNotice{Index: 0x4C, NoPortrait: true})
	if len(g.messages) != 1 || g.messages[0].portraitPage != -1 {
		t.Fatalf("#76 次要 TALK 應以無肖像 modal 排入：%#v", g.messages)
	}
}

// 內政官那一句在**八格變體**的範圍裡：`0x1A6` 是組編號不是索引，
// 實際落在 534–541（docs/spec/48 §2）。
//
// ⚠ 把它當索引會拿到 422「．．．．」——語意上也講得通（一句省略號），
// **所以錯了不會被發現**。這條測試就是為了擋那個。
func TestGovernorRegretTalkIndex(t *testing.T) {
	for variant, want := range []int{534, 535, 536, 537, 538, 539, 540, 541} {
		if got := resolveBattleTalkIndex(governorRegretTalkBase, variant); got != want {
			t.Errorf("變體 %d ⇒ %d，要 %d", variant, got, want)
		}
	}
	// `0x1A6` 的十進位就是 422——**組編號與那個誤讀的索引長得一樣**，
	// 所以只能靠算式分辨，不能靠看數字。
	if resolveBattleTalkIndex(governorRegretTalkBase, 0) == governorRegretTalkBase {
		t.Error("組編號被直接當成索引了")
	}
}

// 訊息框那張臉：原版 `sub_18810` 的 40／60 個呼叫點傳固定的 0x93（通報者），
// 只有變體組與 #58 傳說話者的肖像（docs/spec/106）。
func TestNoticePortraitUsesTheReporterFaceExceptForSpeakers(t *testing.T) {
	g := &game{world: &state.World{}}
	g.world.Generals[2].Portrait = 0x11

	plain := state.TalkNotice{Index: 0x1D, City: -1, Faction: -1, General: 2, Amount: -1}
	if got := g.noticePortraitPage(plain); got != reporterPortraitPage {
		t.Fatalf("一般通報用 %#x，預期通報者 %#x", got, reporterPortraitPage)
	}
	variant := state.TalkNotice{Index: 0x197, City: -1, Faction: -1, General: 2, Amount: -1}
	if got := g.noticePortraitPage(variant); got != 0x11 {
		t.Fatalf("變體組要用說話者的肖像，得到 %#x", got)
	}
	lordGone := state.TalkNotice{Index: talkEnemyLordGone, City: -1, Faction: -1, General: 2, Amount: -1}
	if got := g.noticePortraitPage(lordGone); got != 0x11 {
		t.Fatalf("#58 要用說話者的肖像，得到 %#x", got)
	}
}
