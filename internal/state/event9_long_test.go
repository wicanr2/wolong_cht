package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// TestEvent9LongNaturalRoute 以正常 World.Tick 的「每時 queue consumer」跑過
// 三筆事件 9：玩家勢力可見通知、非玩家釋放，以及俘虜方已滅亡而回到在野。
// 這個 27 小時 fixture 比單次 dispatch 更接近原版長程路徑，但仍是有界測試，
// 不把完整長時間遊玩重新引入驗收範圍。
func TestEvent9LongNaturalRoute(t *testing.T) {
	w := load(t, 0)
	w.Player = 0
	ids := make([]int, 0, 3)
	for i, general := range w.Generals {
		if !general.Alive {
			continue
		}
		ids = append(ids, i)
		if len(ids) == 3 {
			break
		}
	}
	if len(ids) != 3 {
		t.Fatal("找不到三名事件 9 fixture 武將")
	}

	// 第一筆釋放回玩家勢力；第二筆回非玩家勢力；第三筆的 Captor
	// 已滅亡，照 sub_150D7 回到在野。保留每名武將原始名稱／肖像，
	// 讓後續呈現層可以用同一個事件觀測值接 TALK #37。
	for i, id := range ids {
		g := w.Generals[id]
		g.Alive = true
		g.Posted = true
		g.Faction = 2
		g.Captor = 0
		if i == 1 {
			g.Captor = 1
		} else if i == 2 {
			g.Captor = 2
		}
		w.Generals[id] = g
	}
	w.Factions[0].Alive = true
	w.Factions[1].Alive = true
	w.Factions[2].Alive = false

	w.events = [eventQueueEntries]QueuedEvent{
		{Code: uint16(ids[0])<<8 | 9},
		{Code: uint16(ids[1])<<8 | 9},
		{Code: uint16(ids[2])<<8 | 9},
	}
	w.eventCursor = 0
	w.eventDelay = 7
	// 每九個子刻進入一次 hourly；第一個邊界從下一個 tick 開始。
	w.Clock.Subtick = clock.SubticksPerHour - 1
	random := rng.NewFixed(17)
	wantHour := map[int]int{7: ids[0], 17: ids[1], 27: ids[2]}

	for hour := 1; hour <= 27; hour++ {
		var hourEvent Event
		for subtick := 0; subtick < clock.SubticksPerHour; subtick++ {
			ev := w.Tick(random)
			if ev.Clock.Hour {
				hourEvent = ev
			}
		}
		want, shouldRelease := wantHour[hour]
		if shouldRelease {
			if len(hourEvent.ReleasedGenerals) != 1 ||
				hourEvent.ReleasedGenerals[0] != want {
				t.Fatalf("第 %d 小時事件 9 = %v，want [%d]",
					hour, hourEvent.ReleasedGenerals, want)
			}
			t.Logf("第 %d 小時取出事件 9：General #%d", hour, want)
			continue
		}
		if len(hourEvent.ReleasedGenerals) != 0 {
			t.Fatalf("第 %d 小時提前釋放武將：%v", hour, hourEvent.ReleasedGenerals)
		}
	}

	if got := w.Generals[ids[0]]; got.Faction != w.Player || got.Captor != noFaction || got.Posted {
		t.Fatalf("玩家通知武將釋放後狀態錯誤：%+v", got)
	}
	if got := w.Generals[ids[1]]; got.Faction != 1 || got.Captor != noFaction || got.Posted {
		t.Fatalf("非玩家通知武將釋放後狀態錯誤：%+v", got)
	}
	if got := w.Generals[ids[2]]; got.Faction != noFaction || got.Captor != noFaction || got.Posted {
		t.Fatalf("滅亡勢力俘虜釋放後狀態錯誤：%+v", got)
	}
}
