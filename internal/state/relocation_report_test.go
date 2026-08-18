package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/capital"
)

// pickRelocatableFaction 找一個「跑事件 8 會真的搬首都」的勢力，
// 並把首都還原成搬之前的樣子。
func pickRelocatableFaction(t *testing.T, w *World) (faction, next int) {
	t.Helper()
	for i, f := range w.Factions {
		if !f.Alive {
			continue
		}
		old := f.Capital
		if n := w.relocateCapital(i); n != capital.None {
			w.Factions[i].Capital = old
			return i, n
		}
	}
	return -1, -1
}

// TestAIRelocationReportsOnlyWithEnvoy 釘住「沒有外交官就一句話都不報」。
//
// 原版 `sub_133FD` 在非玩家勢力那條路先讀 `[si+2Ah]`（派駐該勢力的
// 外交官），是 `0xFF` 就直接返回（docs/spec/64 §1.2）。
// **別的勢力遷都是一則情報，不是公開事實。**
func TestAIRelocationReportsOnlyWithEnvoy(t *testing.T) {
	run := func(t *testing.T, withEnvoy bool) (*World, Event, int, int) {
		t.Helper()
		w := load(t, 0)
		w.Player = -1
		faction, next := pickRelocatableFaction(t, w)
		if faction < 0 {
			t.Skip("這個劇本沒有可供事件 8 驗證的遷都候選")
		}
		envoy := noFaction
		if withEnvoy {
			for i, g := range w.Generals {
				if g.Alive && g.Faction != faction {
					envoy = i
					break
				}
			}
			if envoy == noFaction {
				t.Skip("找不到可當外交官的武將")
			}
		}
		w.Factions[faction].Diplomat = envoy
		w.events[0] = QueuedEvent{Code: uint16(faction)<<8 | 8, Param: 0xFFFF}
		w.eventCursor, w.eventDelay = 0, 1
		var ev Event
		w.dispatchQueuedEvent(&ev)
		if got := w.Factions[faction].Capital; got != next {
			t.Fatalf("首都沒搬：%d，want %d", got, next)
		}
		return w, ev, faction, envoy
	}

	t.Run("沒有外交官", func(t *testing.T) {
		_, ev, _, _ := run(t, false)
		if len(ev.TalkNotices) != 0 {
			t.Fatalf("沒有外交官卻報了 %d 則：%+v", len(ev.TalkNotices), ev.TalkNotices)
		}
	})

	t.Run("有外交官", func(t *testing.T) {
		w, ev, faction, envoy := run(t, true)
		if len(ev.TalkNotices) != 2 {
			t.Fatalf("應該是「通報 ＋ 內容」兩則，得到 %d 則：%+v",
				len(ev.TalkNotices), ev.TalkNotices)
		}
		// #57「駐{3}勢力的外交官{1}大人前來報告。」
		if n := ev.TalkNotices[0]; n.Index != 0x39 || n.Faction != faction || n.General != envoy {
			t.Fatalf("第一則不是 #57 通報：%+v", n)
		}
		// 0x1A4 是組編號，展開由呈現層做；這裡只釘住欄位齊全。
		n := ev.TalkNotices[1]
		if n.Index != CapitalMovedTalkBase {
			t.Fatalf("第二則的組編號 = %#x，want %#x", n.Index, CapitalMovedTalkBase)
		}
		if n.Faction != faction || n.General != envoy ||
			n.City != w.Factions[faction].Capital {
			t.Fatalf("第二則缺欄位（{3} 勢力、{2} 新首都、說話者外交官）：%+v", n)
		}
	})
}
