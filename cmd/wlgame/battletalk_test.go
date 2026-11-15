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

// 原版1C315的角色0是玩家，玩家進攻或守城都在下框。
func TestBattleTalkPlacementFollowsPlayerRole(t *testing.T) {
	var table text.Table
	for index, literal := range map[int]string{1: "attacker", 2: "defender"} {
		table.Messages[index] = text.Message{Lines: []text.Line{{Parts: []text.Part{{Raw: []byte(literal)}}}}}
	}
	for _, player := range []int{tactical.AttackerSide, tactical.DefenderSide} {
		b := &tactical.Battle{PlayerSide: player}
		p := &state.Pending{Battle: b}
		g := &game{world: &state.World{}, lib: &library.Library{Talk: &table}}
		g.battleTalkSession.initialize(b, []BattleTalkEntry{
			{Index: 1, Side: 0, Portrait: 13}, {Index: 2, Side: 1, Portrait: 27},
		})
		got := g.battleTalkState(b, p)
		wantTop, wantBottom := "defender", "attacker"
		wantTopPortrait, wantBottomPortrait := 27, 13
		if player == tactical.DefenderSide {
			wantTop, wantBottom = wantBottom, wantTop
			wantTopPortrait, wantBottomPortrait = wantBottomPortrait, wantTopPortrait
		}
		if got.Top != wantTop || got.Bottom != wantBottom || got.TopPortrait != wantTopPortrait || got.BottomPortrait != wantBottomPortrait {
			t.Fatalf("玩家側%d的台詞或頭像位置錯誤：%+v", player, got)
		}
	}
}

// re60 §6.1的展開後索引以說話型0為基準，退卻與兩種城壁分支各不相同。
func TestPlayerBattleCommandTalkIndices(t *testing.T) {
	for _, tc := range []struct {
		command           int
		defends, accepted bool
		want              int
	}{
		{0, false, true, 622}, {1, false, true, 630}, {2, false, true, 638},
		{3, false, true, 646}, {4, false, true, 654}, {5, false, true, 606},
		{3, true, true, 662}, {3, false, false, 582}, {3, true, false, 582},
	} {
		group := playerBattleTalkGroup(tc.command, tc.defends, tc.accepted)
		if got := resolveBattleTalkIndex(group, 0); got != tc.want {
			t.Errorf("指令%d 守方%v 接受%v：索引%d，應為%d", tc.command, tc.defends, tc.accepted, got, tc.want)
		}
	}
	if got := resolveBattleTalkIndex(playerBattleTalkGroup(3, false, true), 7); got != 653 {
		t.Fatalf("正常曹操242拍登城，夏侯淵說話型7：索引%d，應為653", got)
	}
	for _, command := range []int{-1, 6, 100} {
		if playerBattleTalkGroup(command, false, true) != -1 {
			t.Fatalf("無效指令%d不應顯示喊話", command)
		}
	}
}

func TestBattleTalkLifetimeUsesOriginalEventFrame(t *testing.T) {
	var table text.Table
	table.Messages[621] = text.Message{Lines: []text.Line{{Parts: []text.Part{{Raw: []byte("retreat")}}}}}
	table.Messages[861] = text.Message{Lines: []text.Line{{Parts: []text.Part{{Raw: []byte("script")}}}}}
	for _, tc := range []struct{ frame, event, group, remaining int }{
		{51, 50, 0x1CE, 59}, {68, 67, 0x1CE, 59},
		{573, 573, 0x1B0, 60}, {575, 573, 0x1B0, 58},
		{632, 573, 0x1B0, 1}, {633, 573, 0x1B0, 0},
		{573, 574, 0x1B0, 0},
	} {
		b := &tactical.Battle{Frame: tc.frame}
		p := &state.Pending{Battle: b, Attacker: 0, Defender: 1}
		w := &state.World{}
		w.Corps[0].Alive = true
		w.Generals[0].TalkVariant = 7
		g := &game{world: w, lib: &library.Library{Talk: &table}}
		ok := g.queueBattleTalk(p, 0, tc.group, tc.event)
		if ok != (tc.remaining > 0) || g.battleTalkSession.queue.remaining[0] != tc.remaining {
			t.Fatalf("frame%d，事件%d：剩餘%v，應為%d", tc.frame, tc.event, g.battleTalkSession.queue.remaining, tc.remaining)
		}
		if tc.remaining > 0 {
			g.tickBattleTalk(tc.remaining)
			if g.battleTalkSession.queue.active() {
				t.Fatal("到期後不能仍保留喊話")
			}
		}
	}
}

func TestBattleTalkMarkerRetainsThreeCellNameSlot(t *testing.T) {
	var table text.Table
	table.Messages[100] = text.Message{Lines: []text.Line{{Parts: []text.Part{{Marker: '1'}, {Raw: []byte("!")}}}}}
	g := &game{lib: &library.Library{Talk: &table}}
	for _, tc := range []struct{ name, want string }{
		{"呂布", "呂布　!"}, {"夏侯淵", "夏侯淵!"},
	} {
		got, name, ok := g.battleTalkText(BattleTalkEntry{Index: 100, Other: tc.name})
		if !ok || got != tc.want || name != tc.name {
			t.Fatalf("三格姓名欄：文字%q、換色名字%q，應為%q、%q", got, name, tc.want, tc.name)
		}
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
	if got, _, ok := g.battleTalkText(BattleTalkEntry{Index: 694}); ok || got != "" {
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
		t.Fatal("Enter／Space 應一次收掉兩側")
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
