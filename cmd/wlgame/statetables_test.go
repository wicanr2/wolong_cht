package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 名字欄位的補白是**全形空格 `A1 40`**，不是 ASCII 空白
// （`docs/spec/138` §3）。用 TrimSpace 砍不掉，名字後面會拖著
// 看不見的兩個 byte，逐欄比就一路紅。
func TestNameHexStripsIdeographicPadding(t *testing.T) {
	// 「曹操」＝ B1E4 BEDE，後面補一個全形空格。
	got := nameHex("\xb1\xe4\xbe\xde\xa1\x40")
	if got != "B1E4BEDE" {
		t.Errorf("nameHex = %q，應為 B1E4BEDE（補白要砍掉）", got)
	}
	// 三個字滿格的不該被動到。
	full := nameHex("\xb1\xe4\xbe\xde\xa5\x40")
	if full != "B1E4BEDEA540" {
		t.Errorf("nameHex = %q，滿三個字不該砍", full)
	}
	// ASCII 空白**不是**補白，要留著——砍掉會讓兩邊長度對不上。
	if got := nameHex("\xb1\xe4 "); got != "B1E420" {
		t.Errorf("nameHex = %q，ASCII 空白不是補白", got)
	}
}

// 旗標 byte 拆成四個 bool 之後要組得回去（`docs/formats/08` §3）。
// 組得回來表示拆對了；bit 0 組不回來，那是已知的缺口。
func TestGeneralFlagsRoundTrips(t *testing.T) {
	for _, raw := range []uint8{0x80, 0x90, 0xA0, 0xC0, 0xD0, 0xE0} {
		g := state.General{
			Alive:                raw&0x80 != 0,
			Sovereign:            raw&0x40 != 0,
			VanishIfAffinityGone: raw&0x20 != 0,
			LoyalToDeath:         raw&0x10 != 0,
		}
		if got := generalFlags(g); got != raw {
			t.Errorf("旗標 %02X 拆開再組回來變成 %02X", raw, got)
		}
	}
	// bit 0 沒有欄位接得住——四個劇本裡只有劇本三的張衛設著它。
	if got := generalFlags(state.General{Alive: true}); got&0x01 != 0 {
		t.Errorf("bit 0 憑空多出來：%02X", got)
	}
}
