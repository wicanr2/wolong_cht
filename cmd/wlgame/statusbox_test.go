package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

// newTalkTestGame 起一個只有 TALK 資料的 game——這幾支測的是文字代入，
// 不需要世界狀態以外的東西。
func newTalkTestGame(t *testing.T) *game {
	t.Helper()
	lib, err := library.LoadWithOptions("../../workplace/orig/dosv", library.LoadOptions{
		TalkJSON: "../../translations/talk-dosv-corrected.json",
	})
	if err != nil {
		t.Skipf("缺原版素材：%v", err)
	}
	return &game{lib: lib, world: &state.World{Player: 0}}
}
// 代入欄位的顏色與定寬是**標記的性質**（docs/spec/119 §3.1）：
// `\1`／`\4` 色 9、`\2` 色 0x0B、`\3`／`\5` 色 0x0C，五個都補到三個全形字。
func TestTalkMarkerInkAndPadding(t *testing.T) {
	for marker, want := range map[byte]int{
		'1': talkInkPerson, '4': talkInkPerson,
		'2': talkInkCity,
		'3': talkInkLord, '5': talkInkLord,
		'6': -1, '7': -1,
	} {
		if got := talkMarkerInk(marker); got != want {
			t.Errorf(`\%c 的顏色 ＝ %d，want %d`, marker, got, want)
		}
	}
	vars, fields := padTalkVars(map[byte]string{
		'2': "許昌", '3': "曹操", '6': "不補",
	})
	if got := vars['2']; got != "許昌　" {
		t.Errorf(`\2 補成 %q，want "許昌　"`, got)
	}
	if got := vars['6']; got != "不補" {
		t.Errorf(`\6 不該補，卻變成 %q`, got)
	}
	if len(fields) != 2 {
		t.Fatalf("欄位數 ＝ %d，want 2（\\6 不換色）", len(fields))
	}
	for _, f := range fields {
		if f.ink != talkInkCity && f.ink != talkInkLord {
			t.Errorf("欄位 %q 的顏色 ＝ %d", f.text, f.ink)
		}
	}
}

// 一般訊息框也要帶欄位（docs/spec/119 §3.1）——先前只有狀態列有，
// 所以 TALK #54 的據點名在原版是橘的、remake 是白的。
func TestGeneralMessageBoxCarriesTalkFields(t *testing.T) {
	g := newTalkTestGame(t)
	g.enqueueTalk(0x36, map[byte]string{'2': "濮陽"})
	if len(g.messages) != 1 {
		t.Fatalf("訊息數 ＝ %d", len(g.messages))
	}
	fields := g.messages[0].fields
	if len(fields) != 1 || fields[0].ink != talkInkCity {
		t.Fatalf("欄位 ＝ %#v，want 一個色 %#x 的", fields, talkInkCity)
	}
	if fields[0].text != "濮陽　" {
		t.Errorf("欄位文字 ＝ %q，want 補到三格", fields[0].text)
	}
}

// 換色的儲存格在反白列上是 XOR 12，不是另一個量出來的常數
//（docs/spec/124 §3.6）。⭐ 這一條釘住的是「新增一種換色不必再量一次」。
func TestListWarnInkFollowsHighlightXor(t *testing.T) {
	if got, want := listInkWarn^chrome.HighlightXOR, listInkWarnSelected; got != want {
		t.Fatalf("%d XOR 12 ＝ %d，want %d", listInkWarn, got, want)
	}
	if listInkWarnSelected != 6 {
		t.Errorf("反白列的換色 ＝ %d，docs/playtest/98 量到 6", listInkWarnSelected)
	}
}
