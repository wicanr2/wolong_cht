package phone

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 擋住世界的兩種決定：外交提案、撥款請求。
//
// ⚠ **不接這兩個世界就停住**：`World.tick` 開頭就檢查它們，
// 有一個在就整個 tick 不跑。手機版少接一個的症狀是「玩到某一天畫面不動了」，
// 而且看不出原因。遭遇不在這裡：原版遭遇時**沒有**「戰鬥指揮／委任」選單，
// 玩家那一方沒委任就直接進戰場（docs/spec/105）。

// modalKind 是現在擋著世界的是哪一種。
type modalKind int

const (
	modalNone modalKind = iota
	modalDiplomacy
	modalFunding
)

// ModalKind 回報現在有沒有東西擋著世界。
func (s *Session) ModalKind() modalKind {
	switch {
	case s.world.PendingDiplomacy() != nil:
		return modalDiplomacy
	case s.world.PendingFunding() != nil:
		return modalFunding
	}
	return modalNone
}

// ModalTitle 是這個決定的標題。
func (s *Session) ModalTitle() string {
	switch s.ModalKind() {
	case modalDiplomacy:
		c := s.world.PendingDiplomacy()
		return s.Localise(s.world.LordName(c.Source)) + " 提出外交要求"
	case modalFunding:
		c := s.world.PendingFunding()
		return fmt.Sprintf("%s 請求 %d 資金",
			s.Localise(s.world.Generals[c.Officer].Name), c.RequestedAmount)
	}
	return ""
}

// ModalOptions 是這個決定的選項。
func (s *Session) ModalOptions() []string {
	switch s.ModalKind() {
	case modalDiplomacy:
		// 原版三選一（TALK #283）：答應／提示金額／拒絕，順序＝ state.DiplomacyOption。
		// 「提示金額」開數字鍵盤（docs/spec/103）。
		return []string{"答應", "提示金額", "拒絕"}
	case modalFunding:
		return []string{"照要求撥款", "拒絕"}
	}
	return nil
}

// PickModal 選了第 i 項。
func (s *Session) PickModal(i int) {
	switch s.ModalKind() {
	case modalDiplomacy:
		if i == int(state.DiplomacyOfferFunds) {
			s.amountPad = true
			return
		}
		s.world.ResolveDiplomacy(state.DiplomacyOption(i))
	case modalFunding:
		opt := state.FundingFullAmount
		if i != 0 {
			opt = state.FundingReject
		}
		s.world.ResolveFunding(opt)
	}
}


// modalOptionRect 是第 i 個選項的矩形。選項貼著主區下緣排。
func modalOptionRect(i, n int) (x, y, w, h int) {
	_, my, mw, mh := MapRect()
	cell := mw / n
	return i * cell, my + mh - BattleRowH, cell, BattleRowH
}

func (s *Session) tapModal(lx, ly float64) bool {
	if s.amountPad {
		for i := range amountKeys {
			if x, y, w, h := amountKeyRect(i); inRect(lx, ly, x, y, w, h) {
				s.pressAmountKey(i)
				return true
			}
		}
		return true
	}
	opts := s.ModalOptions()
	for i := range opts {
		if x, y, w, h := modalOptionRect(i, len(opts)); inRect(lx, ly, x, y, w, h) {
			s.PickModal(i)
			return true
		}
	}
	return true // 擋著世界的時候，別的地方點了不算數
}

func (s *Session) drawModal(dst *ebiten.Image, td *textdraw.Drawer) {
	mx, my, mw, mh := MapRect()
	// 半透明壓在地圖上：**地圖還看得見**，玩家才知道這個決定發生在哪。
	fillRect(dst, mx, my, mw, mh, inkOverlay())
	if td == nil || !td.Available() {
		return
	}
	if s.amountPad {
		s.drawAmountPad(dst, td)
		return
	}
	title := s.ModalTitle()
	td.Draw(dst, title, mx+(mw-td.Width(title))/2, my+mh/2-FontH-8, inkText())
	opts := s.ModalOptions()
	for i, o := range opts {
		x, y, w, h := modalOptionRect(i, len(opts))
		s.window(dst, x, y, w, h, inkBar())
		td.Draw(dst, o, x+(w-td.Width(o))/2, y+(h-FontH)/2, inkText())
	}
}

// ── 提示金額的數字鍵盤（docs/spec/103）────────────────────────────
//
// 原版的數值輸入器（`sub_17C6E`，docs/spec/78）有數字、「百」、刪一位、
// 最大、清除、確定；這裡一鍵一動作，全部直接呼叫規則層的
// `EditDiplomacyOfferAmount`，上限與清零規則不在這一層。
// **版面是 remake 差異**：手機放不下原版那張 640×400 的數字視窗。

type amountKey struct {
	label string
	edit  state.AmountEdit
	digit int
}

// amountKeyCancel／amountKeyOK 兩顆不是編輯動作，用哨兵值標。
const (
	amountKeyCancel state.AmountEdit = 200
	amountKeyOK     state.AmountEdit = 201
)

var amountKeys = [16]amountKey{
	{"7", state.AmountAppendDigit, 7}, {"8", state.AmountAppendDigit, 8}, {"9", state.AmountAppendDigit, 9}, {"刪除", state.AmountDeleteDigit, 0},
	{"4", state.AmountAppendDigit, 4}, {"5", state.AmountAppendDigit, 5}, {"6", state.AmountAppendDigit, 6}, {"百", state.AmountAppendHundred, 0},
	{"1", state.AmountAppendDigit, 1}, {"2", state.AmountAppendDigit, 2}, {"3", state.AmountAppendDigit, 3}, {"最大", state.AmountSetMax, 0},
	{"清除", state.AmountClear, 0}, {"0", state.AmountAppendDigit, 0}, {"取消", amountKeyCancel, 0}, {"確定", amountKeyOK, 0},
}

// amountPadHeadH 是鍵盤上方標題＋金額那一段的高度（兩列字）。
const amountPadHeadH = LineH*2 + 16

// amountKeyRect 是第 i 顆鍵（4 欄 × 4 列）的矩形，每顆留 6 px 間隙。
func amountKeyRect(i int) (x, y, w, h int) {
	mx, my, mw, mh := MapRect()
	const gap = 6
	cw, ch := mw/4, (mh-amountPadHeadH)/4
	col, row := i%4, i/4
	return mx + col*cw + gap, my + amountPadHeadH + row*ch + gap, cw - gap*2, ch - gap*2
}

// AmountPadOpen 回報數字鍵盤是不是開著（測試與返回鍵用）。
func (s *Session) AmountPadOpen() bool { return s.amountPad }

func (s *Session) pressAmountKey(i int) {
	k := amountKeys[i]
	switch k.edit {
	case amountKeyCancel:
		s.amountPad = false
	case amountKeyOK:
		s.amountPad = false
		s.world.ResolveDiplomacy(state.DiplomacyOfferFunds)
	default:
		s.world.EditDiplomacyOfferAmount(k.edit, k.digit)
	}
}

func (s *Session) drawAmountPad(dst *ebiten.Image, td *textdraw.Drawer) {
	mx, my, mw, _ := MapRect()
	title := s.ModalTitle()
	td.Draw(dst, title, mx+(mw-td.Width(title))/2, my+8, inkText())
	amount := "0"
	if c := s.world.PendingDiplomacy(); c != nil {
		amount = fmt.Sprintf("%d", c.OfferAmount)
	}
	label := "提示金額 " + amount
	td.Draw(dst, label, mx+(mw-td.Width(label))/2, my+8+LineH, inkSelect())
	for i, k := range amountKeys {
		x, y, w, h := amountKeyRect(i)
		bg := inkBar()
		if k.edit == amountKeyOK {
			bg = inkSelect()
		}
		s.window(dst, x, y, w, h, bg)
		ink := inkText()
		if k.edit == amountKeyOK {
			ink = inkInk()
		}
		td.Draw(dst, k.label, x+(w-td.Width(k.label))/2, y+(h-FontH)/2, ink)
	}
}
