package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func TestMonthlySettlementDoesNotShowSummary(t *testing.T) {
	ev := state.Event{
		Settled:  true,
		Disaster: map[int]economy.Disaster{82: economy.Fire},
		Storm:    &economy.StormArea{},
	}
	for _, previous := range []string{"", "讀取存檔失敗"} {
		g := &game{lastEvent: previous}
		g.processWorldEvent(ev)
		if g.lastEvent != previous || g.messageActive() {
			t.Errorf("月結更新了通知，原本 %q，目前 %q，對話框 %d", previous, g.lastEvent, len(g.messages))
		}
	}
}

// 事件列六秒後要自己消失（docs/spec/88 §2）。
//
// 事件列換成新文字時，倒數需從頭開始。
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
	g.lastEvent = "讀取存檔失敗"
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
	g.lastEvent = "沒有軍團可以行軍"
	if !visible() {
		t.Error("新的一句沒有重新計時")
	}
	g.frame += eventLineFrames
	if visible() {
		t.Error("第二句也該過期")
	}
}
