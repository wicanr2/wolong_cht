package phone

import "github.com/wicanr2/wolong_cht/internal/assets/world"

// 輸入的路由。**所有點擊只有這一個入口**——平台殼（`mobile/wolong`）
// 只負責把觸控與滑鼠換成邏輯座標，判斷點到什麼是這裡的事。
// 這樣桌面驗收與手機走的是同一條路徑（docs/mobile/android-plan.md §6）。

// Tap 處理一次點擊。回傳 true 表示它吃掉了這一次點擊。
//
// 順序就是遮蔽順序：指令列永遠在最上面，其次是開著的面板，最後才是地圖。
// **反過來的話**，面板底下的地圖會跟著被點到，而玩家看不到那件事發生。
func (s *Session) Tap(lx, ly float64) bool {
	if i, ok := commandAt(lx, ly); ok {
		s.OpenSheet(Command(i))
		return true
	}
	if s.advise.stage != adviseIdle {
		return s.tapAdvise(lx, ly)
	}
	if s.sheet.open {
		return s.tapSheet(lx, ly)
	}
	return s.SelectAt(lx, ly)
}

// commandAt 回報點到第幾個指令鈕。
func commandAt(lx, ly float64) (int, bool) {
	for i := 0; i < int(numCommands); i++ {
		x, y, w, h := CommandRect(i)
		if lx >= float64(x) && lx < float64(x+w) &&
			ly >= float64(y) && ly < float64(y+h) {
			return i, true
		}
	}
	return 0, false
}

// rowAt 回報點到列表的第幾列（已經算進捲動）。-1 表示沒點到列。
func (s *Session) rowAt(ly float64) int {
	_, my, _, mh := MapRect()
	top := my + tabH
	if ly < float64(top) || ly >= float64(my+mh) {
		return -1
	}
	i := int(ly-float64(top))/rowH + s.sheet.scroll
	if i >= len(s.sheetRows()) {
		return -1
	}
	return i
}

// tabAt 回報點到第幾個分頁。-1 表示沒點到。
func (s *Session) tabAt(lx, ly float64) int {
	tabs := s.Tabs()
	if len(tabs) == 0 {
		return -1
	}
	_, my, mw, _ := MapRect()
	if ly < float64(my) || ly >= float64(my+tabH) {
		return -1
	}
	cell := mw / len(tabs)
	i := int(lx) / cell
	if i >= len(tabs) {
		return -1
	}
	return i
}

func (s *Session) tapSheet(lx, ly float64) bool {
	if i := s.tabAt(lx, ly); i >= 0 {
		s.SetSheetTab(i)
		return true
	}
	row := s.rowAt(ly)
	if row < 0 {
		return true // 點在面板的空白處：吃掉，不要穿透到地圖
	}
	switch s.sheet.cmd {
	case CmdAdvise:
		s.PickAdvise(row)
	case CmdList:
		s.tapListRow(row)
	case CmdCorps:
		if s.sheet.tab == 1 {
			s.tapCorpsFormRow(row)
			return true
		}
		s.tapCorpsRow(row, lx)
	case CmdSystem:
		s.tapSystemRow(row, lx)
	}
	return true
}

// tapListRow 點一覽表的一列。據點那一頁會把鏡頭移過去——
// 一覽表在原版是唯讀的，這是手機版加的便利（remake 差異）。
func (s *Session) tapListRow(row int) {
	if s.sheet.tab != 1 {
		return
	}
	n := 0
	for i := range s.world.Cities {
		if s.world.Cities[i].Owner != s.world.Player {
			continue
		}
		if n == row {
			s.focusCity(i)
			return
		}
		n++
	}
}

// tapCorpsRow 點軍團清單的一列：左半把鏡頭移過去，右半派往選中的據點。
func (s *Session) tapCorpsRow(row int, lx float64) {
	n := 0
	for i := range s.world.Corps {
		c := &s.world.Corps[i]
		if !c.Alive || c.Faction != s.world.Player {
			continue
		}
		if n == row {
			if lx >= float64(LogicalW/2) {
				s.MarchSelected(i)
				return
			}
			cols, rows := s.viewTiles()
			s.setCamera(c.X-cols/2, c.Y-rows/2)
			s.sheet = sheet{}
			return
		}
		n++
	}
}

// focusCity 把鏡頭移到某個據點並選中它，然後收掉面板。
func (s *Session) focusCity(i int) {
	c := &s.world.Cities[i]
	cols, rows := s.viewTiles()
	s.setCamera(c.X+world.CityCentreDX-cols/2, c.Y-rows/2)
	s.selected = i
	s.sheet = sheet{}
}

// tapSystemRow 點系統頁的一列。
//
// ⭐ 存檔那一頁**左半是存、右半是讀**。手機上一列只能有一個動作太浪費，
// 而「長按＝另一個動作」在 Android 上與系統手勢打架（docs/mobile/android-ux.md §3）。
func (s *Session) tapSystemRow(row int, lx float64) {
	switch s.sheet.tab {
	case 0:
		switch row {
		case 1:
			s.SetSpeed(s.speed + 1) // 檔位越大越慢
		case 2:
			s.SetSpeed(s.speed - 1)
		}
	case 1:
		if row < 0 || row >= SaveSlots {
			return
		}
		if lx < float64(LogicalW/2) {
			s.lastErr = s.SaveSlot(row)
		} else {
			s.lastErr = s.LoadSlot(row)
		}
	}
}

func (s *Session) tapAdvise(lx, ly float64) bool {
	choices := s.AdviseChoices()
	if len(choices) == 0 {
		// 君主已經講完了：點哪裡都是「知道了」。
		s.CloseAdvise()
		return true
	}
	_, my, _, mh := MapRect()
	if s.advise.stage == advisePickAlly || s.advise.stage == advisePickTarget {
		top := my + tabH
		if ly < float64(top) {
			return true
		}
		i := int(ly-float64(top))/rowH + s.sheet.scroll
		if i < len(choices) {
			s.sheet.scroll = 0
			s.PickAdviseChoice(i)
		}
		return true
	}
	// 說服的五個理由貼著面板下緣排，對白在上面。
	top := my + mh - len(choices)*rowH
	if ly < float64(top) {
		return true
	}
	i := int(ly-float64(top)) / rowH
	s.PickAdviseChoice(i)
	return true
}
