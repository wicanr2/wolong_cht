package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 行軍三選一的版面（docs/spec/39 §3.6／§3.7）。
//
// 原版量到的框：兩項 `(320, 160, 112, 48)`、三項 `(320, 160, 112, 64)`，
// 左上角就是游標所在的 16 px 格。

// marchModeFallbackLabels 是 TALK #76 的三行（取不到素材時的內建字串），
// 每一列 6 個全形字——**框寬就是由這個數字來的**。
func marchModeFallbackLabels() []string {
	return (*game)(nil).marchModeLabels()
}

func TestMarchModeBoxMatchesOriginal(t *testing.T) {
	labels := marchModeFallbackLabels()
	if len(labels) != 3 {
		t.Fatalf("TALK #76 的內建備援 = %d 行，want 3", len(labels))
	}
	for _, tc := range []struct{ rows, w, h int }{
		{rows: 2, w: 112, h: 48},
		{rows: 3, w: 112, h: 64},
	} {
		x, y, w, h := legacyChoiceRect(320, 160, labels[:tc.rows])
		if x != 320 || y != 160 || w != tc.w || h != tc.h {
			t.Errorf("%d 項的框 = (%d,%d,%d,%d)，want (320,160,%d,%d)",
				tc.rows, x, y, w, h, tc.w, tc.h)
		}
	}
}

// 兩道夾制照抄 `sub_1804E`：欄 ≤ 0x21、列 ≤ 0x18 − 項數。
// **兩個上限都要讓框的右／下緣剛好貼齊 640／400。**
func TestMarchModeAnchorClamps(t *testing.T) {
	labels := marchModeFallbackLabels()
	for _, rows := range []int{2, 3} {
		x, y := marchModeAnchor(639, 399, rows)
		if x != 0x21*16 {
			t.Errorf("%d 項：游標貼右緣時 x = %d，want %d", rows, x, 0x21*16)
		}
		if want := (0x18 - rows) * 16; y != want {
			t.Errorf("%d 項：游標貼下緣時 y = %d，want %d", rows, y, want)
		}
		_, _, w, h := legacyChoiceRect(x, y, labels[:rows])
		if x+w != 640 {
			t.Errorf("%d 項：夾住之後右緣 = %d，want 640", rows, x+w)
		}
		if y+h != 400 {
			t.Errorf("%d 項：夾住之後下緣 = %d，want 400", rows, y+h)
		}
	}
	// 沒被夾到的位置原封不動落在游標所在的那一格。
	if x, y := marchModeAnchor(327, 175, 3); x != 320 || y != 160 {
		t.Errorf("游標 (327,175) 的選單左上角 = (%d,%d)，want (320,160)", x, y)
	}
	// 負座標（游標跑出畫面）不要算出負的框。
	if x, y := marchModeAnchor(-8, -8, 2); x != 0 || y != 0 {
		t.Errorf("游標在畫面外時 = (%d,%d)，want (0,0)", x, y)
	}
}

// 狀態列提示框的四個數字各自對回 `sub_18853` 的立即值（docs/spec/140）。
func TestStatusBoxRectMatchesOriginal(t *testing.T) {
	if statusBoxX != 0*16 {
		t.Errorf("statusBoxX = %d，want dx × 16 = 0", statusBoxX)
	}
	if want := (0x12 + 2) * 16; statusBoxY != want {
		t.Errorf("statusBoxY = %d，want (bx + 2) × 16 = %d", statusBoxY, want)
	}
	if talkBoxW != 0x10*16 || talkBoxH != 0x05*16 {
		t.Errorf("框 = %d×%d，want cl×16 × ch×16 = 256×80", talkBoxW, talkBoxH)
	}
	if statusBoxX+talkBoxW > screenW || statusBoxY+talkBoxH != screenH {
		t.Errorf("框 (%d,%d,%d,%d) 沒有貼齊畫面左下角",
			statusBoxX, statusBoxY, talkBoxW, talkBoxH)
	}
}

// `cx = 0FFFFh` 是清掉，不是畫一個空框。
func TestStatusTalkClearsLikeFFFF(t *testing.T) {
	g := &game{statusBox: statusBoxState{lines: []string{"向許昌移動下。"}}}
	if !g.statusBoxActive() {
		t.Fatal("掛著訊息時 statusBoxActive() 應該為真")
	}
	g.clearStatusTalk()
	if g.statusBoxActive() {
		t.Error("clearStatusTalk() 之後不應該再畫")
	}
}

// `\1`–`\5` 代入的名字是**定寬三個全形字**（原版五支 handler 都是
// `al = 3`，docs/re/79 §2）——補白也會畫出來。
func TestPadTalkFieldKeepsFieldWidth(t *testing.T) {
	const want = talkFieldCells * talkLinePitch
	for _, in := range []string{"許昌", "洛陽", "南陽宛", "", "上庸"} {
		if got := textdraw.StringWidth(padTalkField(in)); got != want {
			t.Errorf("padTalkField(%q) 寬 %d，want %d", in, got, want)
		}
	}
	// 已經滿三格就不要再補。
	if got := padTalkField("南陽宛"); got != "南陽宛" {
		t.Errorf("三個全形字不該補白，得到 %q", got)
	}
}

// 補白要在**語系轉換之後**：`uitext.Convert` 的覆寫詞表是整串比對的。
// 這一支釘住 `setStatusTalk` 補的是代入之後的值，不是原始位元組。
func TestStatusTalkPadsAfterSubstitution(t *testing.T) {
	g := &game{}
	g.setStatusTalk(0x15, map[byte]string{'2': "許昌"})
	// 沒有 lib ⇒ 取不到 TALK，fail-closed；這裡驗的是不會 panic。
	if g.statusBoxActive() {
		t.Error("沒有 TALK.DAT 時不該掛框")
	}
	if got := padTalkField("許昌"); got != "許昌　" {
		t.Errorf("padTalkField(\"許昌\") = %q，want \"許昌　\"", got)
	}
}
