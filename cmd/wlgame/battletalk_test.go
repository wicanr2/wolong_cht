package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
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

func TestBattleTalkQueueHasTwoEntryLimit(t *testing.T) {
	var q battleTalkQueue
	if !q.enqueue(BattleTalkEntry{Index: 1, Duration: 10}) ||
		!q.enqueue(BattleTalkEntry{Index: 2, Duration: 10}) {
		t.Fatal("前兩筆 TALK 應可入列")
	}
	if q.enqueue(BattleTalkEntry{Index: 3, Duration: 10}) {
		t.Fatal("戰術開場 TALK 不得超過兩筆")
	}
	if got := len(q.entries); got != 2 {
		t.Fatalf("queue 長度 = %d，want 2", got)
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

func TestBattleTalkQueueAdvanceAndAutoExpire(t *testing.T) {
	var q battleTalkQueue
	q.enqueue(BattleTalkEntry{Index: 1, Duration: 2})
	q.enqueue(BattleTalkEntry{Index: 2, Duration: 3})

	current, ok := q.current()
	if !ok || current.Index != 1 {
		t.Fatalf("初始 TALK = %#v, ok=%v", current, ok)
	}
	if !q.advance() {
		t.Fatal("Enter／Space／滑鼠應能推進到下一句")
	}
	current, ok = q.current()
	if !ok || current.Index != 2 || q.remaining != 3 {
		t.Fatalf("推進後 TALK = %#v, 剩餘=%d, ok=%v", current, q.remaining, ok)
	}
	q.tick(2)
	if current, ok = q.current(); !ok || current.Index != 2 {
		t.Fatalf("尚未到上限時不應自動消失：%#v, ok=%v", current, ok)
	}
	q.tick(1)
	if _, ok = q.current(); ok {
		t.Fatal("達到固定 duration 上限後應自動消失")
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
	if _, ok := g.battleTalkSession.queue.current(); ok {
		t.Fatal("新 battle 不得沿用上一場的 TALK queue")
	}
}

func TestBattleTalkClearRemovesCurrentEntry(t *testing.T) {
	g := &game{}
	b := &tactical.Battle{}
	g.battleTalkSession = battleTalkSession{battle: b}
	g.battleTalkSession.queue.enqueue(BattleTalkEntry{Index: 694, Duration: 10})

	clearBattleTalkSession(g, b)
	if _, ok := g.battleTalkSession.queue.current(); ok {
		t.Fatal("清理後仍有 current TALK entry")
	}
	if g.battleTalkSession.battle != nil {
		t.Fatal("清理後仍保留 battle 指標")
	}
}

func TestBattleTalkStartOnceDoesNotRefillAfterQueueDrains(t *testing.T) {
	first := &tactical.Battle{}
	entries := []BattleTalkEntry{{Index: 694, Duration: 1}, {Index: 702, Duration: 1}}
	var s battleTalkSession

	s.initialize(first, entries)
	if !s.initialized || len(s.queue.entries) != 2 {
		t.Fatalf("第一次初始化 = initialized:%v queue:%d", s.initialized, len(s.queue.entries))
	}
	s.queue.advance()
	s.queue.advance()
	if _, ok := s.queue.current(); ok {
		t.Fatal("兩句推進後 queue 應為空")
	}
	s.initialize(first, entries)
	if _, ok := s.queue.current(); ok {
		t.Fatal("同一 battle 在 queue 耗盡後不得重新補列")
	}

	second := &tactical.Battle{}
	s.initialize(second, entries)
	if s.battle != second || !s.initialized || len(s.queue.entries) != 2 {
		t.Fatalf("新 battle 未重新初始化：battle=%p initialized:%v queue:%d",
			s.battle, s.initialized, len(s.queue.entries))
	}
}
