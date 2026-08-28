package phone

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 擋住世界的三種決定：遭遇（戰鬥指揮／委任）、外交提案、撥款請求。
//
// ⚠ **不接這三個世界就停住**：`World.tick` 開頭就檢查它們，
// 有一個在就整個 tick 不跑（原版進戰術畫面前也會先問「戰鬥指揮／委任」，
// 那個選單同樣不能讓時鐘在背景偷偷前進）。
// 手機版少接一個的症狀是「玩到某一天畫面不動了」，而且看不出原因。

// modalKind 是現在擋著世界的是哪一種。
type modalKind int

const (
	modalNone modalKind = iota
	modalEncounter
	modalDiplomacy
	modalFunding
)

// ModalKind 回報現在有沒有東西擋著世界。
func (s *Session) ModalKind() modalKind {
	switch {
	case s.world.PendingEncounter() != nil:
		return modalEncounter
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
	case modalEncounter:
		c := s.world.PendingEncounter()
		return fmt.Sprintf("%s 與 %s 遭遇",
			s.corpsName(c.Attacker), s.corpsName(c.Defender))
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
	case modalEncounter:
		return []string{"戰鬥指揮", "委任"}
	case modalDiplomacy:
		// ⚠ 「指定金額」要數值輸入器，手機上還沒做——**先不給那一項**，
		// 給了會是一個點下去什麼都不會發生的選項。
		return []string{"接受", "拒絕"}
	case modalFunding:
		return []string{"照要求撥款", "拒絕"}
	}
	return nil
}

// PickModal 選了第 i 項。
func (s *Session) PickModal(i int) {
	switch s.ModalKind() {
	case modalEncounter:
		if i == 0 {
			s.lastErr = s.world.ChooseBattleCommand()
			if s.lastErr == nil {
				s.battle = battleState{view: s.newBattleView(), tacSpeed: DefaultSpeed}
			}
			return
		}
		s.world.ChooseBattleDelegate(s.rand)
	case modalDiplomacy:
		opt := state.DiplomacyAcceptFree
		if i != 0 {
			opt = state.DiplomacyReject
		}
		s.world.ResolveDiplomacy(opt)
	case modalFunding:
		opt := state.FundingFullAmount
		if i != 0 {
			opt = state.FundingReject
		}
		s.world.ResolveFunding(opt)
	}
}

// corpsName 是一支軍團的名字（＝帶兵武將）。−1 是城兵。
func (s *Session) corpsName(corps int) string {
	if corps < 0 || corps >= len(s.world.Generals) {
		return "城兵"
	}
	return s.Localise(s.world.Generals[corps].Name)
}

// modalOptionRect 是第 i 個選項的矩形。選項貼著主區下緣排。
func modalOptionRect(i, n int) (x, y, w, h int) {
	_, my, mw, mh := MapRect()
	cell := mw / n
	return i * cell, my + mh - BattleRowH, cell, BattleRowH
}

func (s *Session) tapModal(lx, ly float64) bool {
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
	title := s.ModalTitle()
	td.Draw(dst, title, mx+(mw-td.Width(title))/2, my+mh/2-FontH-8, inkText())
	opts := s.ModalOptions()
	for i, o := range opts {
		x, y, w, h := modalOptionRect(i, len(opts))
		s.window(dst, x, y, w, h, inkBar())
		td.Draw(dst, o, x+(w-td.Width(o))/2, y+(h-FontH)/2, inkText())
	}
}
