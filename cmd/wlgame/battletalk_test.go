package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestResolveBattleTalkIndex(t *testing.T) {
	for _, tc := range []struct {
		name          string
		base, variant int
		want          int
	}{
		{name: "direct below threshold", base: 0x120, variant: 2, want: 0x120},
		{name: "opening attacker", base: 0x1BA, variant: 0, want: 694},
		{name: "opening variant", base: 0x1BA, variant: 2, want: 696},
		{name: "opening defender", base: 0x1BB, variant: 1, want: 703},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := resolveBattleTalkIndex(tc.base, tc.variant); got != tc.want {
				t.Fatalf("resolveBattleTalkIndex(%#x,%d) = %d, want %d",
					tc.base, tc.variant, got, tc.want)
			}
		})
	}
}

// TestBattleTalkSlotsAreIndexedBySide 釘住「每一側各一個框」。
//
// 原版 `sub_1C315` 依側別把到期時刻寫進 `word_1D322`（側 0）或
// `word_1D324`（側 1），兩邊各自計時、可以同時掛著（docs/spec/60）。
// 同一側再來一句就整個換掉，不排隊。
func TestBattleTalkSlotsAreIndexedBySide(t *testing.T) {
	var q battleTalkSlots
	if !q.set(BattleTalkEntry{Index: 1, Side: 0, Duration: 10}) ||
		!q.set(BattleTalkEntry{Index: 2, Side: 1, Duration: 10}) {
		t.Fatal("兩側各一句應可掛上")
	}
	for side, want := range map[int]int{0: 1, 1: 2} {
		got, ok := q.current(side)
		if !ok || got.Index != want {
			t.Fatalf("側 %d 的框 = %#v, ok=%v，want Index=%d", side, got, ok, want)
		}
	}
	if q.set(BattleTalkEntry{Index: 3, Side: 2, Duration: 10}) {
		t.Fatal("只有兩側，side=2 不得掛上")
	}
	if !q.set(BattleTalkEntry{Index: 4, Side: 0, Duration: 10}) {
		t.Fatal("同一側再來一句應可覆寫")
	}
	if got, _ := q.current(0); got.Index != 4 {
		t.Fatalf("同側新句未覆寫舊句：Index=%d", got.Index)
	}
	if got, _ := q.current(1); got.Index != 2 {
		t.Fatalf("覆寫側 0 不應動到側 1：Index=%d", got.Index)
	}
}

func TestBattleTalkUnknownMarkerFailsClosed(t *testing.T) {
	var table text.Table
	table.Messages[694] = text.Message{Lines: []text.Line{{Parts: []text.Part{{Marker: '1'}}}}}
	g := &game{lib: &library.Library{Talk: &table}}
	if got, ok := g.battleTalkText(BattleTalkEntry{Index: 694}); ok || got != "" {
		t.Fatalf("未知 marker 不得顯示：%q, ok=%v", got, ok)
	}
}

// TestBattleTalkSlotsExpireIndependently 釘住兩側的計時器互不干擾。
func TestBattleTalkSlotsExpireIndependently(t *testing.T) {
	var q battleTalkSlots
	q.set(BattleTalkEntry{Index: 1, Side: 0, Duration: 2})
	q.set(BattleTalkEntry{Index: 2, Side: 1, Duration: 5})

	q.tick(2)
	if _, ok := q.current(0); ok {
		t.Fatal("側 0 到期後應消失")
	}
	if got, ok := q.current(1); !ok || got.Index != 2 {
		t.Fatalf("側 1 還沒到期不應被側 0 帶走：%#v, ok=%v", got, ok)
	}
	if !q.active() {
		t.Fatal("還有一側掛著時 active 應為 true")
	}
	q.tick(3)
	if q.active() {
		t.Fatal("兩側都到期後 active 應為 false")
	}

	// 玩家按鍵一次收掉兩側。
	q.set(BattleTalkEntry{Index: 3, Side: 0, Duration: 60})
	q.set(BattleTalkEntry{Index: 4, Side: 1, Duration: 60})
	if !q.clearAll() || q.active() {
		t.Fatal("Enter／Space／滑鼠應一次收掉兩側")
	}
}

func TestBattleTalkNewBattleResetsSession(t *testing.T) {
	g := &game{}
	first := &tactical.Battle{}
	second := &tactical.Battle{}
	g.battleTalkSession.initialize(first, []BattleTalkEntry{{Index: 694, Duration: 10}})

	g.bindBattleTalk(second)
	if g.battleTalkSession.battle != second {
		t.Fatal("新 battle 未取代舊 session")
	}
	if g.battleTalkSession.initialized {
		t.Fatal("新 battle 不應沿用舊 initialized 狀態")
	}
	if g.battleTalkSession.queue.active() {
		t.Fatal("新 battle 不得沿用上一場的對白框")
	}
}

func TestBattleTalkClearRemovesCurrentEntry(t *testing.T) {
	g := &game{}
	b := &tactical.Battle{}
	g.battleTalkSession = battleTalkSession{battle: b}
	g.battleTalkSession.queue.set(BattleTalkEntry{Index: 694, Duration: 10})

	clearBattleTalkSession(g, b)
	if g.battleTalkSession.queue.active() {
		t.Fatal("清理後仍有掛著的對白框")
	}
	if g.battleTalkSession.battle != nil {
		t.Fatal("清理後仍保留 battle 指標")
	}
}

func TestBattleTalkStartOnceDoesNotRefillAfterBoxesExpire(t *testing.T) {
	first := &tactical.Battle{}
	entries := []BattleTalkEntry{
		{Index: 694, Side: 0, Duration: 1},
		{Index: 702, Side: 1, Duration: 1},
	}
	var s battleTalkSession

	s.initialize(first, entries)
	if !s.initialized || !s.queue.active() {
		t.Fatalf("第一次初始化 = initialized:%v active:%v", s.initialized, s.queue.active())
	}
	s.queue.tick(1)
	if s.queue.active() {
		t.Fatal("兩側都到期後不應還有框")
	}
	s.initialize(first, entries)
	if s.queue.active() {
		t.Fatal("同一 battle 在對白框結束後不得重新掛上")
	}

	second := &tactical.Battle{}
	s.initialize(second, entries)
	if s.battle != second || !s.initialized || !s.queue.active() {
		t.Fatalf("新 battle 未重新初始化：battle=%p initialized:%v active:%v",
			s.battle, s.initialized, s.queue.active())
	}
}

// TestBattleTalkDurationMatchesOriginal 釘住對白框的壽命 ＝ 60 tick。
//
// 原版 `sub_1C315` 設的到期時刻是「目前節拍 ＋ 0x3C」，由 `sub_1A12A`
// 每 tick 比對後擦掉（docs/spec/60）。門強度條走同一個機制但常數是 0x14——
// 兩個值不同，寫成同一個就會有一邊錯。
func TestBattleTalkDurationMatchesOriginal(t *testing.T) {
	if battleTalkDuration != 0x3C {
		t.Fatalf("對白框壽命 = %d，預期 0x3C（60）", battleTalkDuration)
	}
	if battleTalkDuration == 0x14 {
		t.Fatal("對白框用了門強度條的 0x14——那是另一個計時器")
	}
}

// TestNoticeExpandsGroupIndexByTalkVariant 釘住通知的組編號展開。
//
// TALK.DAT 索引 ≥ 0x196 是八格一組，`sub_18810` 用說話者的**原始**
// `+0x1E` 選組內第幾個（docs/re/25 §1）。遷都那一組（0x1A4 → 518–525）
// 剛好橫跨主公型（0–2）與臣下型（3–7）——收斂成 0–2 會把外交官的
// 回報變成君主的命令句（docs/spec/64 §3）。
func TestNoticeExpandsGroupIndexByTalkVariant(t *testing.T) {
	g := &game{world: &state.World{}}
	g.world.Generals[0].TalkVariant = 1 // 主公型
	g.world.Generals[1].TalkVariant = 5 // 臣下型

	for _, tc := range []struct {
		name   string
		notice state.TalkNotice
		want   int
	}{
		{"門檻以下不展開", state.TalkNotice{Index: 0x39, General: 1}, 0x39},
		{"君主下令", state.TalkNotice{Index: state.CapitalMovedTalkBase, General: 0}, 519},
		{"外交官回報", state.TalkNotice{Index: state.CapitalMovedTalkBase, General: 1}, 523},
		{"沒有說話者取第 0 格", state.TalkNotice{Index: state.CapitalMovedTalkBase, General: -1}, 518},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := g.noticeTalkIndex(tc.notice); got != tc.want {
				t.Fatalf("展開 = %d，want %d", got, tc.want)
			}
		})
	}
}
