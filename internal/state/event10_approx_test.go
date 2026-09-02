package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/clock"
)

func event10PrisonerFixture(t *testing.T) (*World, int) {
	t.Helper()
	w := load(t, 0)
	w.Player = 0
	w.SetApproximateEvent10(true)
	for id, g := range w.Generals {
		if !g.Alive {
			continue
		}
		g.Faction = w.Player
		g.Captor = 1
		g.Posted = true
		g.Timer = 0
		// 歸降那一支要 `+0x1C == +0x19`（docs/re/77 §2.2）：
		// 關押他的勢力剛好是他心向的那一個。fixture 讓條件成立，
		// 條件不成立的情形由 TestCaptiveSurrenderNeedsAffinity 蓋。
		g.Affinity = w.Player
		w.Generals[id] = g
		return w, id
	}
	t.Fatal("找不到可用的事件 10 俘虜 fixture 武將")
	return nil, -1
}

func TestApproximateEvent10ProducerUsesKnownRawContract(t *testing.T) {
	w, id := event10PrisonerFixture(t)
	w.events = [eventQueueEntries]QueuedEvent{}
	if !w.produceApproximateEvent10(&sequenceRand{values: []int{0x01}}) {
		t.Fatal("近似事件 10 producer 沒有排入逃走消息")
	}
	if got := w.events[0]; got != (QueuedEvent{
		Code: uint16(id)<<8 | 0x0A, Param: approximateEvent10EscapeTalk,
	}) {
		t.Fatalf("逃走 raw payload = %#v，want General<<8|0A／TALK 41h", got)
	}
	if g := w.Generals[id]; g.Faction != noFaction || g.Captor != noFaction || g.Posted {
		t.Fatalf("逃走近似狀態錯誤：%+v", g)
	}

	w, id = event10PrisonerFixture(t)
	w.events = [eventQueueEntries]QueuedEvent{}
	if !w.produceApproximateEvent10(&sequenceRand{values: []int{0x20}}) {
		t.Fatal("近似事件 10 producer 沒有排入歸降消息")
	}
	if got := w.events[0]; got != (QueuedEvent{
		Code: uint16(id)<<8 | 0x0A, Param: approximateEvent10JoinTalk,
	}) {
		t.Fatalf("歸降 raw payload = %#v，want General<<8|0A／TALK 42h", got)
	}
	if g := w.Generals[id]; g.Faction != w.Player || g.Captor != noFaction || g.Posted {
		t.Fatalf("歸降近似狀態錯誤：%+v", g)
	}
}

func TestApproximateEvent10ProducerIsBoundedAndDisableable(t *testing.T) {
	w, id := event10PrisonerFixture(t)
	w.events = [eventQueueEntries]QueuedEvent{}
	w.SetApproximateEvent10(false)
	if w.produceApproximateEvent10(&sequenceRand{values: []int{0}}) {
		t.Fatal("關閉近似 producer 後不應排入事件 10")
	}
	if w.events[0].Code != 0 || w.Generals[id].Captor == noFaction {
		t.Fatal("關閉近似 producer 不應改動 queue／俘虜")
	}

	w.SetApproximateEvent10(true)
	w.Generals[id].Timer = 2
	if w.produceApproximateEvent10(&sequenceRand{values: []int{0}}) {
		t.Fatal("倒數未到時不應產生事件 10")
	}
	if w.Generals[id].Timer != 1 || w.events[0].Code != 0 {
		t.Fatalf("倒數閘錯誤：timer=%d event=%#v", w.Generals[id].Timer, w.events[0])
	}
}

func TestApproximateEvent10ReentersIdleClockConsumer(t *testing.T) {
	w, id := event10PrisonerFixture(t)
	w.events = [eventQueueEntries]QueuedEvent{}
	if !w.produceApproximateEvent10(&sequenceRand{values: []int{0x20}}) {
		t.Fatal("近似 producer 沒有建立 raw 事件")
	}
	w.eventCursor, w.eventDelay = 0, 1
	w.Clock.Subtick = clock.SubticksPerHour - 1
	ev := w.Tick(&sequenceRand{values: []int{0xFF}})
	if !ev.Clock.Hour || len(ev.TalkNotices) != 1 || !sameNotice(ev.TalkNotices[0], TalkNotice{
		Index: approximateEvent10JoinTalk, City: -1, Faction: -1, General: id, Amount: -1,
	}) {
		t.Fatalf("近似事件 10 未沿 idle clock consumer 顯示：clock=%+v notices=%#v",
			ev.Clock, ev.TalkNotices)
	}
}

// 關押他的勢力不是他心向的那一個時，歸降那一段不成立——
// 原版 sub_15940 先比 `[si+1Ch] == [si+19h]` 才走歸降（docs/re/77 §2.2）。
func TestCaptiveSurrenderNeedsAffinity(t *testing.T) {
	w, id := event10PrisonerFixture(t)
	w.events = [eventQueueEntries]QueuedEvent{}
	g := w.Generals[id]
	g.Affinity = 5 // 不是玩家勢力
	w.Generals[id] = g

	if w.produceApproximateEvent10(&sequenceRand{values: []int{0x20}}) {
		t.Fatal("心向的勢力對不上，不該產生歸降事件")
	}
	if got := w.Generals[id]; got.Captor != 1 || !got.Posted {
		t.Fatalf("狀態不該被改：%+v", got)
	}

	// 逃走那一支不看這個欄位，照樣成立。
	if !w.produceApproximateEvent10(&sequenceRand{values: []int{0x01}}) {
		t.Fatal("逃走那一支不該受心向的勢力影響")
	}
}
