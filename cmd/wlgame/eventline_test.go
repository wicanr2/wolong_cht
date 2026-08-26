package main

import "testing"

// 事件列六秒後要自己消失（docs/spec/88 §2）。
//
// ⭐ 這一支釘的是**重計時的機制**，不是某一個設定點：`lastEvent` 有
// 三十幾個設定點，其中幾個是 `+=`。靠每個設定點都記得記時間遲早會漏，
// 而漏掉的症狀是「那一種事件的框永遠不消失」——那正是使用者看到的
// 「月結一直出現」。
func TestEventLineExpiresAndRestartsOnNewText(t *testing.T) {
	g := &game{}
	visible := func() bool {
		// drawStrategy 的判斷條件抄在這裡：內容變了就重計時，
		// 然後看有沒有超過六秒。
		if g.lastEvent != g.lastEventShown {
			g.lastEventShown, g.lastEventAt = g.lastEvent, g.frame
		}
		return g.lastEvent != "" && g.frame-g.lastEventAt < eventLineFrames
	}

	if visible() {
		t.Fatal("沒有事件時不該畫")
	}
	g.frame = 100
	g.lastEvent = "月結"
	if !visible() {
		t.Fatal("剛設定就該看得到")
	}
	g.frame = 100 + eventLineFrames - 1
	if !visible() {
		t.Error("六秒內就消失了")
	}
	g.frame = 100 + eventLineFrames
	if visible() {
		t.Error("六秒之後還在")
	}
	// 同一句再來一次也要重新出現——月結每個月都會發生。
	g.lastEvent = "月結　災害1"
	if !visible() {
		t.Error("新的一句沒有重新計時")
	}
	g.frame += eventLineFrames
	if visible() {
		t.Error("第二句也該過期")
	}
}
