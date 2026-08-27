package phone

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// sheet 是指令列四個入口打開的全螢幕面板。
//
// ⭐ 「全螢幕」指的是**佔滿主區**，上下兩條留著。手機上把入口列也蓋掉的話，
// 換一個入口要先關再開；留著等於一次點擊就換頁（docs/mobile/android-ux.md §4）。
type sheet struct {
	open bool
	cmd  Command
	tab  int
	// scroll 是列表捲到第幾列。以**列**為單位，不是像素——
	// 半列的位移在點陣字上只會讓字看起來糊。
	scroll int
}

// 列表的版面。
const (
	sheetPadX   = 20
	tabH        = 44
	rowH        = 30
	rowTextDX   = 12
	sheetHeadDY = 10
)

// Tabs 是目前 sheet 的分頁標題，沒有分頁時回 nil。
func (s *Session) Tabs() []string {
	if !s.sheet.open {
		return nil
	}
	switch s.sheet.cmd {
	case CmdList:
		return []string{"武將", "據點", "勢力", "軍團"}
	case CmdCorps:
		return []string{"現有", "編成"}
	case CmdSystem:
		return []string{"速度", "存檔", "語言", "關於"}
	}
	return nil
}

// OpenSheet 打開一個入口。再點同一個入口會關掉它。
func (s *Session) OpenSheet(c Command) {
	if s.sheet.open && s.sheet.cmd == c {
		s.sheet = sheet{}
		return
	}
	s.sheet = sheet{open: true, cmd: c}
}

// SheetOpen 回報現在有沒有開著面板。
func (s *Session) SheetOpen() bool { return s.sheet.open }

// SheetCommand 回報開著的是哪一個入口。沒開時回 -1。
func (s *Session) SheetCommand() Command {
	if !s.sheet.open {
		return -1
	}
	return s.sheet.cmd
}

// SheetTab 回報目前的分頁。
func (s *Session) SheetTab() int { return s.sheet.tab }

// SetSheetTab 換分頁，並把捲動歸零——換頁還留在第 40 列會讓人以為表是空的。
func (s *Session) SetSheetTab(i int) {
	if tabs := s.Tabs(); i >= 0 && i < len(tabs) {
		s.sheet.tab, s.sheet.scroll = i, 0
	}
}

// Back 是 Android 返回鍵的行為：關面板 → 收小卡 → 交給呼叫端決定離開。
// 回傳 true 表示這一次按鍵已經被吃掉。
func (s *Session) Back() bool {
	switch {
	case s.sheet.open:
		s.sheet = sheet{}
		return true
	case s.selected >= 0:
		s.selected = -1
		return true
	}
	return false
}

// ScrollRows 捲動列表。捲不動時什麼都不做（不是錯誤）。
func (s *Session) ScrollRows(d int) {
	n := s.scrollableRows()
	if n == 0 {
		return
	}
	_, _, _, mh := MapRect()
	visible := (mh - tabH) / rowH
	max := n - visible
	if max < 0 {
		max = 0
	}
	s.sheet.scroll += d
	if s.sheet.scroll > max {
		s.sheet.scroll = max
	}
	if s.sheet.scroll < 0 {
		s.sheet.scroll = 0
	}
}

// scrollableRows 是目前這一頁有幾列可以捲。
// 進言的選對象也是一張長列表，用同一個捲動狀態。
func (s *Session) scrollableRows() int {
	if s.advise.stage == advisePickAlly || s.advise.stage == advisePickTarget {
		return len(s.AdviseChoices())
	}
	if !s.sheet.open {
		return 0
	}
	return len(s.sheetRows())
}

// sheetRow 是列表的一列：左邊是名稱，右邊是若干欄。
type sheetRow struct {
	name string
	cols []string
	// dim 表示這一列不是玩家的東西，畫得暗一點。
	dim bool
}

// sheetRows 產生目前分頁的內容。
//
// ⚠ **只讀不寫**。一覽表在原版是唯讀視窗，手機版照這條；
// 要改東西一律走進言或軍團（docs/mobile/android-ux.md §4）。
func (s *Session) sheetRows() []sheetRow {
	if !s.sheet.open {
		return nil
	}
	switch s.sheet.cmd {
	case CmdList:
		switch s.sheet.tab {
		case 0:
			return s.generalRows()
		case 1:
			return s.cityRows()
		case 2:
			return s.factionRows()
		default:
			return s.corpsRows()
		}
	case CmdCorps:
		if s.sheet.tab == 1 {
			return s.corpsFormRows()
		}
		return s.corpsRows()
	case CmdAdvise:
		return s.adviseRows()
	case CmdSystem:
		return s.systemRows()
	}
	return nil
}

// generalRows 列玩家勢力的武將。
//
// ⚠ 只列**自己的**：原版的武將一覽也是一次看一個勢力，
// 全部 255 名一起列在手機上捲不完，而且那不是遊戲給玩家的資訊。
func (s *Session) generalRows() []sheetRow {
	w := s.world
	rows := make([]sheetRow, 0, 32)
	for i := range w.Generals {
		g := &w.Generals[i]
		if !g.Alive || g.Faction != w.Player {
			continue
		}
		rows = append(rows, sheetRow{
			name: s.Localise(g.Name),
			cols: []string{
				fmt.Sprintf("武 %d", g.Martial),
				fmt.Sprintf("統 %d", g.Command),
				fmt.Sprintf("政 %d", g.Politics),
				s.generalPost(i),
			},
		})
	}
	return rows
}

// generalPost 回報這名武將現在在做什麼。空白表示待命。
func (s *Session) generalPost(n int) string {
	w := s.world
	for i := range w.Cities {
		if w.Cities[i].Governor == n {
			return "內政 " + s.Localise(w.Cities[i].Name)
		}
	}
	if c := &w.Corps[n]; c.Alive {
		return "軍團"
	}
	if w.Factions[w.Player].Lord == n {
		return "君主"
	}
	if w.Factions[w.Player].Advisor == n {
		return "軍師"
	}
	// ⚠ 待命要畫成 `－` 而不是空字串。欄位是**靠右往左**排的，
	// 空字串會讓這一列少一欄，右邊三欄跟著往右移一格——
	// 整張表看起來像沒對齊。
	return "－"
}

func (s *Session) cityRows() []sheetRow {
	w := s.world
	rows := make([]sheetRow, 0, 32)
	for i := range w.Cities {
		c := &w.Cities[i]
		if c.Owner != w.Player {
			continue
		}
		rows = append(rows, sheetRow{
			name: s.Localise(c.Name),
			cols: []string{
				fmt.Sprintf("生產 %d", c.Production),
				fmt.Sprintf("防災 %d", c.Prevention),
				fmt.Sprintf("城兵 %d", c.Garrison*MenPerPoint),
			},
		})
	}
	return rows
}

// factionRows 列所有還活著的勢力。這一張要看別人——外交要用。
func (s *Session) factionRows() []sheetRow {
	w := s.world
	rows := make([]sheetRow, 0, 16)
	for i := range w.Factions {
		f := &w.Factions[i]
		if !f.Alive {
			continue
		}
		rows = append(rows, sheetRow{
			name: s.Localise(w.LordName(i)),
			cols: []string{
				fmt.Sprintf("據點 %d", f.Cities),
				fmt.Sprintf("武將 %d", f.Generals),
				s.diplomacyLabel(i),
			},
			dim: i != w.Player,
		})
	}
	return rows
}

// diplomacyLabel 是「我方對這個勢力」的關係。對自己回「本國」。
//
// ⚠ 交戰與交友度是**兩個獨立的位元**：交戰旗標在最高位，交友度在低七位
//（`internal/rules/diplomacy`）。畫面上要先看交戰——正在打的時候
// 交友度是多少不影響玩家的判斷。
func (s *Session) diplomacyLabel(other int) string {
	if other == s.world.Player {
		return "本國"
	}
	return s.world.Friendship[s.world.Player][other].Level().String()
}

func (s *Session) corpsRows() []sheetRow {
	w := s.world
	rows := make([]sheetRow, 0, 16)
	for i := range w.Corps {
		c := &w.Corps[i]
		if !c.Alive || c.Faction != w.Player {
			continue
		}
		rows = append(rows, sheetRow{
			name: s.Localise(w.Generals[i].Name),
			cols: []string{
				fmt.Sprintf("兵 %d", c.Men*MenPerPoint),
				fmt.Sprintf("士氣 %d", c.Morale),
				s.corpsWhere(c),
			},
		})
	}
	if len(rows) == 0 {
		rows = append(rows, sheetRow{name: "尚未編成軍團", dim: true})
	}
	return rows
}

// corpsWhere 說這個軍團在哪、往哪去。
func (s *Session) corpsWhere(c *state.Corps) string {
	w := s.world
	at := "行軍中"
	if c.Node >= 0 && c.Node < len(w.Cities) && c.X == w.Cities[c.Node].X && c.Y == w.Cities[c.Node].Y {
		at = s.Localise(w.Cities[c.Node].Name)
	}
	if c.Ordered >= 0 && c.Ordered < len(w.Cities) && c.Ordered != c.Node {
		return at + " → " + s.Localise(w.Cities[c.Ordered].Name)
	}
	return at
}

// adviseRows 是進言的五項。用字取自 `TALK.DAT` #77，不是內部術語
//（桌面版 `cmd/wlgame/advise.go` 同一份出處）。
func (s *Session) adviseRows() []sheetRow {
	names := s.adviseLabels()
	rows := make([]sheetRow, 0, len(names))
	for i, n := range names {
		rows = append(rows, sheetRow{name: n, cols: []string{s.adviseHint(i)}})
	}
	return rows
}

func (s *Session) systemRows() []sheetRow {
	switch s.sheet.tab {
	case 0:
		return []sheetRow{
			{name: "戰略速度", cols: []string{fmt.Sprintf("%d", s.speed)}},
			{name: "慢", cols: []string{"⟨ 點這一列變慢"}},
			{name: "快", cols: []string{"⟩ 點這一列變快"}},
			{name: "音效", cols: []string{s.soundValue()}},
		}
	case 1:
		return s.saveRows()
	case 2:
		return s.languageRows()
	default:
		return []sheetRow{
			{name: "臥龍傳 Remake", cols: []string{"手機版"}},
			{name: "原版", cols: []string{"NEO･GETEN 1994 / 松崗 1995"}},
			{name: "資料來源", cols: []string{"使用者自備，不隨程式散布"}, dim: true},
		}
	}
}


// soundValue 是「音效」那一列的值（docs/spec/92 §2.3）。
//
// ⚠ **「未接入」與「關」是兩件事**：關是玩家的選擇，未接入是這一台
// 根本沒有音檔（沒跑過 `tools/bgm2ogg.sh`，或 APK 不是完整版）。
// 混成同一個字會讓缺口從畫面上消失。
func (s *Session) soundValue() string {
	switch {
	case s.music == nil || !s.music.Available():
		return "未接入"
	case s.music.Enabled():
		return "開"
	default:
		return "關"
	}
}

// languageRows 是可切的語言（docs/spec/86 §4）。
//
// **手機沒有命令列**，這一頁是 Android 唯一能換語言的地方。
// 每一列用該語言自己的寫法寫，目前選中的那一列打勾——
// 換過去之後畫面全是那個語言，用中文列出來反而認不出自己選了什麼。
func (s *Session) languageRows() []sheetRow {
	cur := s.Language()
	rows := make([]sheetRow, 0, len(LanguageChoices))
	for _, l := range LanguageChoices {
		// ⚠ 記號要挑**倚天 Big5 有的字**：`✓` 不在裡面，畫出來是方框。
		mark := ""
		if l.Lang == cur {
			mark = "●"
		}
		rows = append(rows, sheetRow{name: l.Name, cols: []string{mark}})
	}
	return rows
}

// saveRows 是四個存檔槽。
//
// ⚠ **來源與去處分開**：讀的是使用者匯入的 `orig/SAVE.DAT`，
// 寫的是 `save/SAVE.DAT`。原版資產一律唯讀（CLAUDE.md §9）。
func (s *Session) saveRows() []sheetRow {
	rows := make([]sheetRow, 0, SaveSlots)
	for i := 0; i < SaveSlots; i++ {
		rows = append(rows, sheetRow{
			name: fmt.Sprintf("第 %d 槽", i+1),
			cols: []string{s.slotLabel(i)},
		})
	}
	return rows
}

// sheetFooter 是面板底下那一行提示。
//
// ⭐ 「左半存、右半讀」這種**看不見的分割**一定要寫出來。
// 手機上沒有滑鼠指標可以 hover，玩家不會自己發現一列有兩個動作。
func (s *Session) sheetFooter() string {
	if s.sheet.cmd == CmdSystem && s.sheet.tab == 1 {
		return "點一列的左半＝存檔，右半＝讀檔"
	}
	if s.sheet.cmd == CmdList && s.sheet.tab == 1 {
		return "點一列把鏡頭移到那個據點"
	}
	if s.sheet.cmd == CmdCorps && s.sheet.tab == 0 {
		return "左半＝把鏡頭移過去，右半＝派往地圖上選中的據點"
	}
	if s.sheet.cmd == CmdCorps && s.sheet.tab == 1 {
		return "點位置換兵種：騎馬→弓兵→步兵→空（主將不能空）"
	}
	return ""
}
