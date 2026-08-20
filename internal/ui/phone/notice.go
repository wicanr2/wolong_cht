package phone

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 事件訊息。
//
// ⚠ **不顯示的話事件是無聲的**：災害、戰報、外交結果、月結全部發生在
// 背後，玩家只看得到數字忽然變了。原版是彈出訊息框，手機版改成
// 貼在地圖上緣的一條——**這是 remake 差異**（docs/mobile/android-ux.md §7）。

// noticeHold 是一則訊息停留幾幀。
//
// 原版的訊息框要按鍵才消（`sub_1084A`），手機上每一則都要點會很煩，
// 所以改成自己消失。**六秒**是估的，不是原版的數字——原版沒有這個機制。
const noticeHold = 360

// noticeLines 是最多顯示幾行。
const noticeLines = 3

// notice 是一則等著顯示的訊息。
type notice struct {
	lines []string
	left  int // 還剩幾幀
}

// pushNotices 把這個 tick 的事件訊息收下來。
func (s *Session) pushNotices(ev state.Event) {
	for _, n := range ev.TalkNotices {
		// 擋著世界的外交提案自己會畫，不重複播它的 base 訊息。
		if s.world.PendingDiplomacy() != nil {
			continue
		}
		vars, ok := s.world.TalkNoticeVars(n, big5)
		if !ok {
			continue // fail-closed：寧可不顯示，也不顯示半句
		}
		lines, ok := s.lib.Talk.Lines(n.Index, vars)
		if !ok || len(lines) == 0 {
			continue
		}
		if len(lines) > noticeLines {
			lines = lines[:noticeLines]
		}
		s.notices = append(s.notices, notice{lines: lines, left: noticeHold})
	}
	// 只留最新的幾則：手機的畫面放不下一整串，而舊的那幾則早就過去了。
	if len(s.notices) > 3 {
		s.notices = s.notices[len(s.notices)-3:]
	}
}

// tickNotices 讓訊息自己過期。
func (s *Session) tickNotices() {
	if len(s.notices) == 0 {
		return
	}
	s.notices[0].left--
	if s.notices[0].left <= 0 {
		s.notices = s.notices[1:]
	}
}

// Notice 回傳現在該顯示的那一則，沒有就回 nil。
func (s *Session) Notice() []string {
	if len(s.notices) == 0 {
		return nil
	}
	return s.notices[0].lines
}

// dismissNotice 點一下跳過目前這一則。
func (s *Session) dismissNotice() {
	if len(s.notices) > 0 {
		s.notices = s.notices[1:]
	}
}

func (s *Session) drawNotice(dst *ebiten.Image, td *textdraw.Drawer) {
	lines := s.Notice()
	if len(lines) == 0 || td == nil || !td.Available() {
		return
	}
	mx, my, mw, _ := MapRect()
	h := len(lines)*26 + 16
	fillRect(dst, mx, my, mw, h, inkOverlay)
	strokeRect(dst, mx, my, mw, h, inkEdge)
	for i, l := range lines {
		td.Draw(dst, l, mx+20, my+10+i*26, inkText)
	}
}
