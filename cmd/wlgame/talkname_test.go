package main

import (
	"strings"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// TestRoutTalkUsesAlias 走完整條路：世界狀態 → `reportRout` → 訊息佇列，
// 斷言 `\1` 代進去的是**呼び名**（docs/spec/119）。
//
// ⭐ **不要在測試裡自己呼叫 `TalkName()` 再比對**——那樣兩邊都用同一支，
// 代入點改回 `Name` 也照樣綠。這裡比的是**寫死的字串「孔明」**，
// 而諸葛亮的姓名是「諸葛亮」，兩者一眼分得出來。
func TestRoutTalkUsesAlias(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("找不到原版素材：%v", err)
	}
	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skipf("讀不到劇本：%v", err)
	}
	w.Player = 0

	lead := -1
	for i := range w.Generals {
		if text.Decode([]byte(w.Generals[i].Name), text.Big5) == "諸葛亮" {
			lead = i
			break
		}
	}
	if lead < 0 {
		t.Skip("劇本一裡找不到諸葛亮")
	}
	g0 := &w.Generals[lead]
	g0.Alive, g0.Posted, g0.Faction = true, false, w.Player
	w.Factions[w.Player].Reserves = [3]int{9000, 9000, 9000}

	kinds := [6]army.TroopType{
		army.Cavalry, army.Cavalry, army.Archer,
		army.Archer, army.Infantry, army.Infantry,
	}
	manned := [6]bool{true, true, true, true, true, true}
	if err := w.FormCorps(lead, kinds, manned); err != nil {
		t.Skipf("編不出軍團：%v", err)
	}
	corps := -1
	for i := range w.Corps {
		if w.Corps[i].Alive && w.Leader(i) == lead {
			corps = i
			break
		}
	}
	if corps < 0 {
		t.Fatal("編成之後找不到那支軍團")
	}

	g := &game{lib: lib, world: w}
	g.reportRout(state.CorpsEvent{Corps: corps, Routed: true})
	if len(g.messages) == 0 {
		t.Fatal("敗走沒有排出訊息")
	}
	line := strings.Join(g.messages[0].lines, "")
	if !strings.Contains(line, "孔明") {
		t.Errorf("敗走訊息 %q 裡沒有呼び名「孔明」——`\\1` 代的還是姓名", line)
	}
	if strings.Contains(line, "諸葛亮") {
		t.Errorf("敗走訊息 %q 代進了姓名「諸葛亮」，原版代的是 +0x08 呼び名", line)
	}
}
