package phone

import (
	"image/color"
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

// 點指令列要開對應的面板，再點一次收掉。
func TestCommandBarOpensAndClosesTheSheet(t *testing.T) {
	s := newTestSession(t)
	x, y, w, h := CommandRect(int(CmdList))
	cx, cy := float64(x+w/2), float64(y+h/2)

	if !s.Tap(cx, cy) || !s.SheetOpen() {
		t.Fatal("點一覽沒有打開面板")
	}
	if s.SheetCommand() != CmdList {
		t.Fatalf("開的是 %v，預期一覽", s.SheetCommand())
	}
	if s.Tap(cx, cy); s.SheetOpen() {
		t.Fatal("再點一次沒有收掉面板")
	}
}

// 面板開著時，點面板上的空白**不可以穿透到地圖**。
func TestSheetSwallowsTapsSoTheMapIsNotTouched(t *testing.T) {
	s := newTestSession(t)
	capital := s.World().Factions[s.World().Player].Capital
	s.Select(capital)
	s.OpenSheet(CmdSystem)

	_, my, _, mh := MapRect()
	s.Tap(10, float64(my+mh-4)) // 面板下緣的空白
	if s.Selected() != capital {
		t.Fatal("點到面板卻改了地圖上的選取")
	}
}

// 返回鍵的順序：關面板 → 收小卡 → 交回系統。
func TestBackClosesOneLayerAtATime(t *testing.T) {
	s := newTestSession(t)
	s.Select(s.World().Factions[s.World().Player].Capital)
	s.OpenSheet(CmdList)

	if !s.Back() || s.SheetOpen() {
		t.Fatal("第一次返回沒有關掉面板")
	}
	if !s.Back() || s.Selected() >= 0 {
		t.Fatal("第二次返回沒有收掉小卡")
	}
	if s.Back() {
		t.Fatal("沒東西可關時返回鍵不可以宣稱吃掉了")
	}
}

// 一覽表**只讀不寫**：翻遍四張表不可以動到 World。
func TestListSheetDoesNotMutateTheWorld(t *testing.T) {
	s := newTestSession(t)
	s.SetPaused(true)
	before := s.World().Fingerprint()
	s.OpenSheet(CmdList)
	for i := range s.Tabs() {
		s.SetSheetTab(i)
		if len(s.sheetRows()) == 0 {
			t.Errorf("第 %d 張表是空的", i)
		}
	}
	if s.World().Fingerprint() != before {
		t.Fatal("只是看一覽表卻改到了 World")
	}
}

// 捲動要夾在範圍內，而且捲不動時不可以留下負的位移。
func TestScrollIsClamped(t *testing.T) {
	s := newTestSession(t)
	s.OpenSheet(CmdList)
	s.ScrollRows(-5)
	if s.sheet.scroll != 0 {
		t.Fatalf("往上捲過頭應該停在 0，得到 %d", s.sheet.scroll)
	}
	s.ScrollRows(9999)
	_, _, _, mh := MapRect()
	visible := (mh - tabH) / rowH
	if want := len(s.sheetRows()) - visible; s.sheet.scroll != max0(want) {
		t.Fatalf("往下捲過頭應該停在 %d，得到 %d", max0(want), s.sheet.scroll)
	}
}

func max0(v int) int {
	if v < 0 {
		return 0
	}
	return v
}

// 速度是**原版檔位 0–4**，不是每畫面幾 tick；夾在範圍內。
func TestSpeedIsTheOriginalLevel(t *testing.T) {
	s := newTestSession(t)
	if s.Speed() != DefaultSpeed {
		t.Fatalf("開局速度是 %d，預期 %d", s.Speed(), DefaultSpeed)
	}
	s.SetSpeed(-3)
	if s.Speed() != 0 {
		t.Fatalf("下限應該是 0，得到 %d", s.Speed())
	}
	s.SetSpeed(99)
	if s.Speed() != 4 {
		t.Fatalf("上限應該是 4，得到 %d", s.Speed())
	}
}

// 進言的五項要取自 TALK.DAT，不是內部術語。
func TestAdviseLabelsComeFromTalkDat(t *testing.T) {
	s := newTestSession(t)
	labels := s.adviseLabels()
	if len(labels) != len(adviseFallbackNames) {
		t.Fatalf("進言選單 %d 列，要 %d", len(labels), len(adviseFallbackNames))
	}
	// 第四列是遷都。原版用字含全形空白，所以比對去空白之後的內容。
	if got := strings.TrimSpace(strings.ReplaceAll(labels[adviseRelocateRow], "　", "")); got != "遷都" {
		t.Errorf("第四列是 %q，預期「遷都」", got)
	}
}

// 遷都要先在地圖上選目的地——沒選就什麼都不做，而且**不可以動到 World**。
func TestRelocateNeedsATargetFirst(t *testing.T) {
	s := newTestSession(t)
	s.SetPaused(true)
	before := s.World().Fingerprint()
	s.PickAdvise(adviseRelocateRow)
	if s.AdviseStage() != adviseIdle {
		t.Fatal("沒選目的地卻進了進言流程")
	}
	if s.World().Fingerprint() != before {
		t.Fatal("沒選目的地卻改到了 World")
	}
	if hint := s.adviseHint(adviseRelocateRow); !strings.Contains(hint, "先在地圖") {
		t.Errorf("提示是 %q，應該要說先選目的地", hint)
	}
}

// 停戰提案：選指令 → 選對象 → 君主開口。**對象清單不含自己**。
func TestCeasefireFlowPicksATargetAndTheLordSpeaks(t *testing.T) {
	s := newTestSession(t)
	s.SetPaused(true)
	s.PickAdvise(1) // 停戰提案
	if s.AdviseStage() != advisePickTarget {
		t.Fatalf("階段是 %v，預期選對象", s.AdviseStage())
	}
	for _, id := range s.adviseFactions() {
		if id == s.World().Player {
			t.Fatal("對象清單裡有自己")
		}
	}
	if len(s.AdviseChoices()) == 0 {
		t.Fatal("沒有可選的對象")
	}
	s.PickAdviseChoice(0)
	if len(s.AdviseLines()) == 0 {
		t.Fatal("選完對象君主一句話都沒說")
	}
	if s.AdviseStage() != advisePersuade && s.AdviseStage() != adviseVerdict {
		t.Fatalf("階段是 %v，預期進說服或已定案", s.AdviseStage())
	}
}

// 請求協助要先選協力方再選對象——兩個都要，順序照原版。
func TestCooperateAsksForAllyBeforeTarget(t *testing.T) {
	s := newTestSession(t)
	s.SetPaused(true)
	s.PickAdvise(2)
	if s.AdviseStage() != advisePickAlly {
		t.Fatalf("階段是 %v，預期先選協力方", s.AdviseStage())
	}
	s.PickAdviseChoice(0)
	if s.AdviseStage() != advisePickTarget {
		t.Fatalf("選完協力方之後是 %v，預期選對象", s.AdviseStage())
	}
	if s.advise.ally < 0 {
		t.Fatal("協力方沒有記下來")
	}
}

// 說服的理由選單要有五列，而且順序**就是索引算式吃的位置**。
func TestReasonLabelsMatchTheOptionOrder(t *testing.T) {
	s := newTestSession(t)
	for _, c := range adviseCommands {
		opts := persuasion.Options(c)
		labels := s.reasonLabels(c)
		if len(labels) != len(opts) {
			t.Fatalf("%v 的理由選單 %d 列，要 %d", c, len(labels), len(opts))
		}
		for i, r := range opts {
			if slot := persuasion.TalkReasonSlot(c, r); slot != i {
				t.Errorf("%v 的第 %d 個理由算成第 %d 格", c, i, slot)
			}
		}
	}
}

// 編成：選武將 → 調六個位置 → 送出。送出之後那支軍團要真的存在。
func TestFormCorpsFlow(t *testing.T) {
	s := newTestSession(t)
	s.SetPaused(true)
	s.OpenSheet(CmdCorps)
	s.SetSheetTab(1)

	ids := s.corpsCandidates()
	if len(ids) == 0 {
		t.Fatal("開局沒有可以帶兵的武將")
	}
	s.tapCorpsFormRow(0) // 選第一個候選人
	if s.form.leader != ids[0] {
		t.Fatalf("選到 %d，預期 %d", s.form.leader, ids[0])
	}
	// 第 1 列是主將那一格；一路點回去也不可以變成空。
	for i := 0; i < 6; i++ {
		s.tapCorpsFormRow(1)
		if !s.form.manned[0] {
			t.Fatal("主將那一格被點成空了")
		}
	}
	s.tapCorpsFormRow(7) // 編成
	if s.lastErr != nil {
		t.Fatalf("編成失敗：%v", s.lastErr)
	}
	if !s.World().Corps[ids[0]].Alive {
		t.Fatal("送出之後軍團不存在")
	}
	if s.SheetTab() != 0 {
		t.Error("編成完沒有跳回現有軍團那一頁")
	}
}

// 派兵要先在地圖上選目的地——沒選就給明確的錯誤，不是靜靜不動。
func TestMarchNeedsATargetFirst(t *testing.T) {
	s := newTestSession(t)
	s.SetPaused(true)
	s.MarchSelected(0)
	if s.lastErr == nil {
		t.Fatal("沒選目的地卻沒有回錯誤")
	}
}

// 非主將的位置可以空，而且「空」是循環的最後一站。
func TestSlotCycleReachesEmptyExceptForTheLeader(t *testing.T) {
	s := newTestSession(t)
	s.form = newCorpsForm()
	// 預設是步兵；再點一次應該變空。
	s.cycleSlot(1)
	if s.form.manned[1] {
		t.Fatal("非主將的位置點不到空")
	}
	s.cycleSlot(0)
	if !s.form.manned[0] {
		t.Fatal("主將的位置變成空了")
	}
}

// 事件訊息要真的出現過——**不顯示的話事件是無聲的**。
//
// ⚠ 這一條也擋「訊息卡住不消失」：跑完之後佇列要是空的。
func TestEventNoticesAppearAndExpire(t *testing.T) {
	s := newTestSession(t)
	seen := 0
	for i := 0; i < 20000 && seen == 0; i++ {
		s.Tick()
		if len(s.Notice()) > 0 {
			seen++
		}
	}
	if seen == 0 {
		t.Fatal("跑了兩萬幀一則事件訊息都沒有")
	}
	for i := 0; i < noticeHold*4; i++ {
		s.SetPaused(true)
		s.tickNotices()
	}
	if len(s.Notice()) > 0 {
		t.Fatal("訊息過了四倍停留時間還在")
	}
}

// 顏色一律取自原版調色盤，**不得是手機層自己的常數**（docs/spec/70）。
//
// ⚠ 判準是「與 `chrome` 相同」而不是「等於某個 RGB」：抄一份 RGB 進來的話
// 這個測試照樣會綠，而那正是 docs/spec/54 §1 記的那個事故。
func TestPhoneUsesOriginalPalette(t *testing.T) {
	s := newTestSession(t)
	if s.ch == nil {
		t.Fatal("沒有載入原版的視窗外框")
	}
	if !s.ch.Available() {
		t.Fatal("外框圖塊沒拿到——面板會畫成純色框")
	}
	for _, c := range []struct {
		name string
		got  color.RGBA
		want color.RGBA
	}{
		{"面板底", inkPanel(), chrome.Menu},
		{"清單底", inkSheet(), chrome.Sheet},
		{"命令列底", inkBar(), chrome.Blank},
		{"深藍底上的字", inkText(), chrome.Paper},
		{"米色底上的字", inkInk(), chrome.Ink},
		{"反白條", inkSelect(), chrome.Select},
	} {
		if c.got != c.want {
			t.Errorf("%s = %v，原版調色盤是 %v", c.name, c.got, c.want)
		}
	}
	// 次要色要真的從調色盤查過，不是留在 fallback。
	if want, err := s.lib.Palette.Bank(0); err == nil {
		if inkDim() != want[secondaryIndex] {
			t.Errorf("次要色 = %v，調色盤第 %d 色是 %v",
				inkDim(), secondaryIndex, want[secondaryIndex])
		}
	}
}
