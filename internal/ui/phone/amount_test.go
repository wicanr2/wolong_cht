package phone

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// 外交提案的「提示金額」要有數字鍵盤（docs/spec/103）：三選一有三項、
// 按鍵改到規則層的 OfferAmount、確定後決定消失。
func TestDiplomacyAmountKeypad(t *testing.T) {
	s := newTestSession(t)
	s.world.BeginDiplomacyChoice(state.DiplomacyChoice{
		Kind: state.DiplomacyCeasefire, Source: 1, Target: s.world.Player,
		InitialAmount: 1200, OfferAmount: 1200,
	})
	if s.ModalKind() != modalDiplomacy {
		t.Fatal("沒有停在外交提案")
	}
	if got := s.ModalOptions(); len(got) != 3 || got[1] != "提示金額" {
		t.Fatalf("三選一是 %v", got)
	}
	s.PickModal(int(state.DiplomacyOfferFunds))
	if !s.AmountPadOpen() {
		t.Fatal("點「提示金額」沒有開鍵盤")
	}
	press := func(label string) {
		for i, k := range amountKeys {
			if k.label == label {
				x, y, w, h := amountKeyRect(i)
				if w < 48 || h < 48 {
					t.Fatalf("鍵 %s 只有 %d×%d，比觸控下限 48 小", label, w, h)
				}
				s.Tap(float64(x+w/2), float64(y+h/2))
				return
			}
		}
		t.Fatalf("沒有 %s 這顆鍵", label)
	}
	press("清除")
	press("1")
	press("2")
	press("百")
	if got := s.world.PendingDiplomacy().OfferAmount; got != 1200 {
		t.Fatalf("清除→1→2→百 之後是 %d，預期 1200", got)
	}
	press("刪除")
	if got := s.world.PendingDiplomacy().OfferAmount; got != 120 {
		t.Fatalf("刪除之後是 %d，預期 120", got)
	}
	press("取消")
	if s.AmountPadOpen() || s.ModalKind() != modalDiplomacy {
		t.Fatal("取消沒有回到三選一")
	}
	s.PickModal(int(state.DiplomacyOfferFunds))
	press("最大")
	press("確定")
	if s.AmountPadOpen() || s.ModalKind() != modalNone || s.world.PendingDiplomacy() != nil {
		t.Fatal("確定之後決定沒有消失")
	}
}

// 返回鍵先關鍵盤、再關別的。
func TestBackClosesAmountPadFirst(t *testing.T) {
	s := newTestSession(t)
	s.world.BeginDiplomacyChoice(state.DiplomacyChoice{Kind: state.DiplomacyCeasefire, InitialAmount: 100, OfferAmount: 100})
	s.PickModal(int(state.DiplomacyOfferFunds))
	if !s.Back() || s.AmountPadOpen() || s.ModalKind() != modalDiplomacy {
		t.Fatal("返回鍵沒有只關掉鍵盤")
	}
}
