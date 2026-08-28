package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/langpack"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// launcherPhase 是一般玩家啟動殼層的純狀態機。它不持有 World，避免
// 尚未確認的劇本／君主被錯誤地綁到 scenario 0。
type launcherPhase uint8

const (
	launcherTitle launcherPhase = iota
	launcherNewGameConfirm
	launcherScenario
	// launcherSelectFaction 是原版的勢力清單（docs/spec/79）。
	// 原版的新遊戲是四層：ＹＥＳ／ＮＯ → 劇本 → **清單** → 君主卡。
	launcherSelectFaction
	launcherSelectPlayer
	launcherGameConfirm
	launcherLoad
	// launcherLanguage 是 remake 自己的一頁（docs/spec/86 §4）。
	// **不加在原版的系統選單裡**——那一頁是逐像素對過的（docs/playtest/39）。
	launcherLanguage
)

type launcherAction uint8

const (
	launcherMoveUp launcherAction = iota + 1
	launcherMoveDown
	launcherConfirm
	launcherCancel
)

type launcherResultKind uint8

const (
	launcherNoResult launcherResultKind = iota
	launcherPreviewScenario
	launcherStartNewGame
	launcherStartLoad
	launcherSetLanguage
)

type launcherResult struct {
	kind     launcherResultKind
	scenario int
	player   int
	slot     int
	// lang 是 launcherSetLanguage 選中的語言（langpack.Choices 的索引）。
	lang int
}

type launcherPlayer struct {
	ID      int
	Lord    string
	Capital string

	// 以下是原版君主選擇卡要用的欄位（docs/spec/27 §1.2）。
	Portrait        int
	Advisor         string
	AdvisorPortrait int
	HasAdvisor      bool // 勢力 +0x02 是 0x7F 就沒有軍師
	Generals        int
	Cities          int
}

type launcherSlot struct {
	Slot      int
	Available bool
	Label     string

	// Title／Year／Month／Day 是原版四槽視窗那四欄（docs/spec/25 §1.2）：
	// 標題來自區塊 +0x40，日期是 +0x06／+0x04／+0x00。
	Title            string
	Year, Month, Day int
}

type launcherModel struct {
	phase  launcherPhase
	cursor  int
	hasSave bool
	scenario        int
	scenarioName    string
	players         []launcherPlayer
	slots           []launcherSlot
	confirmedPlayer int
	confirmLord     string
	notice          string
	// factionTop 是勢力清單的捲動位置（docs/spec/79）。選中的是 cursor，
	// 清單與君主卡**共用同一個 cursor**，不另存一份。
	factionTop      int
	pointerSeen     bool
	pointerX        int
	pointerY        int
}

func newLauncher(hasSave bool, slots []launcherSlot) *launcherModel {
	copySlots := append([]launcherSlot(nil), slots...)
	return &launcherModel{
		phase:   launcherTitle,
		hasSave: hasSave,
		slots:   copySlots,
	}
}

func (l *launcherModel) rowCount() int {
	switch l.phase {
	case launcherTitle:
		// NEW GAME ／（LOAD DATA）／ LANGUAGE
		if l.hasSave {
			return 3
		}
		return 2
	case launcherLanguage:
		return len(langpack.Choices)
	case launcherNewGameConfirm, launcherGameConfirm:
		return 2
	case launcherScenario:
		return 4
	case launcherSelectFaction, launcherSelectPlayer:
		return len(l.players)
	case launcherLoad:
		return len(l.slots)
	default:
		return 0
	}
}

// languageCursor 是目前語言在清單裡的位置——進這一頁時游標停在它上面。
func (l *launcherModel) languageCursor() int {
	// ⚠ **語系狀態不在這裡存第二份**：目前語言的唯一真相是 `uiLang`，
	// launcher 只是把游標停在它上面。
	cur := uiLang.Lang()
	for i, c := range langpack.Choices {
		if c.Lang == cur {
			return i
		}
	}
	return 0
}

func (l *launcherModel) clampCursor() {
	n := l.rowCount()
	if n == 0 {
		l.cursor = 0
		return
	}
	if l.cursor < 0 {
		l.cursor = n - 1
	}
	if l.cursor >= n {
		l.cursor = 0
	}
}

func (l *launcherModel) move(delta int) {
	if l.rowCount() == 0 {
		return
	}
	l.cursor += delta
	l.clampCursor()
	// 勢力清單要讓選中的那一列留在畫面上（docs/spec/79 §2）。
	// 君主卡的 ↑↓ 也算——退回清單時才不會停在看不到的位置。
	if l.phase == launcherSelectFaction || l.phase == launcherSelectPlayer {
		l.scrollFactionListToCursor()
	}
}

func (l *launcherModel) selectRow(row int) bool {
	if row < 0 || row >= l.rowCount() {
		return false
	}
	l.cursor = row
	return true
}

func (l *launcherModel) selectPlayer(id int) bool {
	if l.phase != launcherSelectPlayer && l.phase != launcherSelectFaction {
		return false
	}
	for i, p := range l.players {
		if p.ID == id {
			l.cursor = i
			return true
		}
	}
	return false
}

// setScenarioPlayers 是讀取劇本摘要後的唯一接縫。這裡只保存可選玩家的
// 顯示資料，正式 World 仍要在 launcherGameConfirm 確認後重新建立。
func (l *launcherModel) setScenarioPlayers(index int, name string, players []launcherPlayer) bool {
	if l.phase != launcherScenario || index < 0 || index >= 4 {
		return false
	}
	l.scenario = index
	l.scenarioName = name
	l.players = append([]launcherPlayer(nil), players...)
	l.cursor = 0
	l.notice = ""
	if len(l.players) == 0 {
		l.notice = "本劇本沒有可用的玩家勢力"
		return false
	}
	l.phase = launcherSelectFaction
	l.factionTop = 0
	return true
}



func (l *launcherModel) back() {
	l.notice = ""
	switch l.phase {
	case launcherNewGameConfirm, launcherLoad, launcherLanguage:
		l.phase = launcherTitle
		l.cursor = 0
	case launcherScenario:
		l.phase = launcherNewGameConfirm
		l.cursor = 0
	case launcherSelectFaction:
		l.phase = launcherScenario
		l.cursor = l.scenario
	case launcherSelectPlayer:
		// 君主卡退回清單——原版就是這一層（docs/re/73 §1）。
		l.phase = launcherSelectFaction
	case launcherGameConfirm:
		l.phase = launcherSelectPlayer
	}
}

func (l *launcherModel) apply(action launcherAction) launcherResult {
	switch action {
	case launcherMoveUp:
		l.move(-1)
		return launcherResult{}
	case launcherMoveDown:
		l.move(1)
		return launcherResult{}
	case launcherCancel:
		l.back()
		return launcherResult{}
	case launcherConfirm:
		return l.confirm()
	default:
		return launcherResult{}
	}
}

func (l *launcherModel) confirm() launcherResult {
	switch l.phase {
	case launcherTitle:
		if l.cursor == 0 {
			l.phase = launcherNewGameConfirm
			l.cursor = 0
			l.notice = ""
			return launcherResult{}
		}
		if l.cursor == 1 && l.hasSave {
			l.phase = launcherLoad
			l.cursor = 0
			l.notice = ""
			return launcherResult{}
		}
		if l.cursor == l.rowCount()-1 {
			l.phase = launcherLanguage
			l.cursor = l.languageCursor()
			l.notice = ""
		}
	case launcherLanguage:
		if l.cursor < 0 || l.cursor >= len(langpack.Choices) {
			return launcherResult{}
		}
		return launcherResult{kind: launcherSetLanguage, lang: l.cursor}
	case launcherNewGameConfirm:
		if l.cursor == 0 {
			l.phase = launcherScenario
			l.cursor = 0
			l.notice = ""
		} else {
			l.back()
		}
	case launcherScenario:
		return launcherResult{kind: launcherPreviewScenario, scenario: l.cursor}
	case launcherSelectFaction:
		if l.cursor < 0 || l.cursor >= len(l.players) {
			l.notice = "玩家勢力無效"
			return launcherResult{}
		}
		l.phase = launcherSelectPlayer
		l.notice = ""
	case launcherSelectPlayer:
		if l.cursor < 0 || l.cursor >= len(l.players) {
			l.notice = "玩家勢力無效"
			return launcherResult{}
		}
		l.confirmedPlayer = l.players[l.cursor].ID
		l.confirmLord = l.players[l.cursor].Lord
		l.phase = launcherGameConfirm
		l.cursor = 0
		l.notice = ""
	case launcherGameConfirm:
		if l.cursor == 0 {
			if l.scenario < 0 || l.scenario >= 4 || l.confirmedPlayer < 0 {
				l.notice = "玩家勢力無效"
				return launcherResult{}
			}
			return launcherResult{
				kind:     launcherStartNewGame,
				scenario: l.scenario,
				player:   l.confirmedPlayer,
			}
		}
		l.back()
	case launcherLoad:
		if l.cursor < 0 || l.cursor >= len(l.slots) || !l.slots[l.cursor].Available {
			l.notice = "這個槽位沒有可讀取的資料"
			return launcherResult{}
		}
		return launcherResult{kind: launcherStartLoad, slot: l.slots[l.cursor].Slot}
	}
	return launcherResult{}
}

func (l *launcherModel) playerIndex() int {
	if l.cursor < 0 || l.cursor >= len(l.players) {
		return -1
	}
	return l.players[l.cursor].ID
}

const (
	launcherPanelX      = 112
	launcherPanelY      = 56
	launcherPanelW      = 416
	launcherPanelH      = 288
	launcherTextInset   = 16
	launcherListX       = launcherPanelX + launcherTextInset
	launcherListY       = 112
	launcherListW       = launcherPanelW - launcherTextInset*2
	launcherRowH        = 24
	launcherLoadListY   = 96
	launcherLoadRowH    = 32
	launcherNoticeY     = 288
	launcherHintY       = 312
	launcherHint        = "↑↓ 選擇　Enter 決定　ESC 返回"
)

// launcherTextSafeRect 是外框內、供字與反白列使用的共同安全區。
// Window 的外框佔 panel 四周；再留 8px 內縮，確保最長玩家列與 footer
// 不會蓋到左右柱、上下邊。所有 launcher 的文字座標都必須落在這個矩形。
func launcherTextSafeRect() image.Rectangle {
	return image.Rect(
		launcherPanelX+launcherTextInset,
		launcherPanelY+launcherTextInset,
		launcherPanelX+launcherPanelW-launcherTextInset,
		launcherPanelY+launcherPanelH-launcherTextInset,
	)
}

func launcherRowRect(phase launcherPhase, row int) image.Rectangle {
	switch phase {
	case launcherTitle:
		return image.Rect(launcherListX, 152+row*32, launcherListX+launcherListW, 184+row*32)
	case launcherNewGameConfirm:
		return image.Rect(launcherListX, 184+row*32, launcherListX+launcherListW, 216+row*32)
	case launcherGameConfirm:
		return image.Rect(launcherListX, 192+row*32, launcherListX+launcherListW, 224+row*32)
	case launcherScenario:
		return image.Rect(launcherListX, launcherListY+row*launcherRowH,
			launcherListX+launcherListW, launcherListY+(row+1)*launcherRowH)
	case launcherSelectFaction, launcherSelectPlayer:
		// 這兩頁都不走殼層的清單列：清單有自己的幾何（docs/spec/79），
		// 君主卡只有兩個熱區（docs/spec/27 §2.1）。
		return image.Rectangle{}
	case launcherLoad:
		return image.Rect(launcherListX, launcherLoadListY+row*launcherLoadRowH,
			launcherListX+launcherListW, launcherLoadListY+(row+1)*launcherLoadRowH)
	default:
		return image.Rectangle{}
	}
}

// pointerRow 回傳游標所在的列。
//
// ⚠ **選君主那一頁沒有列可以指。** 它畫的是原版版面的卡片（docs/spec/27），
// 啟動殼層自己的清單一列都沒畫出來——而那些命中格的範圍
// （x 128–512、y 88–280）幾乎蓋滿整張卡片。沿用它們的結果是
// 滑鼠一移過去就換君主、一點下去就直接決定。
// **那一頁的滑鼠改走卡片自己的兩個熱區**（`lordCardHotspotAt`）。
func (l *launcherModel) pointerRow(x, y int) (int, bool) {
	if l.phase == launcherSelectPlayer || l.phase == launcherSelectFaction {
		return 0, false
	}
	for row := 0; row < l.rowCount(); row++ {
		if image.Pt(x, y).In(launcherRowRect(l.phase, row)) {
			return row, true
		}
	}
	return 0, false
}

// directStartFlagWasPassed 是既有驗收入口的明確白名單。它只放會要求
// 「進入遊戲後的固定畫面」的旗標；-orig、-font、-seed、-speed 與
// -save-file 單獨使用仍會進一般玩家 launcher。
func directStartFlagWasPassed() bool {
	const name = true
	directFlags := map[string]bool{
		"scenario-file":      name,
		"scenario":           name,
		"player":             name,
		"shot":               name,
		"open-window":        name,
		"open-list":          name,
		"open-advise":        name,
		"open-form":          name,
		"open-corps":         name,
		"open-battle":        name,
		"open-siege":         name,
		"open-message":       name,
		"open-talk-index":    name,
		"open-outcome":       name,
	}
	found := false
	flagVisit(func(flagName string) {
		if directFlags[flagName] {
			found = true
		}
	})
	return found
}

// flagVisit 是 main.go 對 flag.CommandLine.Visit 的窄包裝，讓純狀態機測試
// 不必改動全域 flag set，也讓 direct 判斷保持一個可替換的接縫。
var flagVisit = func(fn func(string)) {
	// 由 main.go 在 flag.Parse 後安裝；未安裝時沒有直接旗標。
}

func validLauncherPlayer(w *state.World, id int) bool {
	if w == nil || id < 0 || id >= len(w.Factions) {
		return false
	}
	f := w.Factions[id]
	if !f.Alive || f.Lord < 0 || f.Lord >= len(w.Generals) || !w.Generals[f.Lord].Alive {
		return false
	}
	return f.Capital >= 0 && f.Capital < len(w.Cities)
}

func launcherPlayers(w *state.World) []launcherPlayer {
	if w == nil {
		return nil
	}
	players := make([]launcherPlayer, 0, len(w.Factions))
	for id, f := range w.Factions {
		if !validLauncherPlayer(w, id) {
			continue
		}
		p := launcherPlayer{
			ID:       id,
			Lord:     big5(w.LordName(id)),
			Capital:  big5(w.Cities[f.Capital].Name),
			Generals: f.Generals,
			Cities:   f.Cities,
		}
		if f.Lord >= 0 && f.Lord < len(w.Generals) {
			p.Portrait = w.Generals[f.Lord].Portrait
		}
		// 軍師欄 0x7F ＝ 無（docs/formats/08 勢力 +0x02）。
		if f.Advisor >= 0 && f.Advisor < len(w.Generals) && f.Advisor != state.NoAdvisor {
			p.HasAdvisor = true
			p.Advisor = big5(w.Generals[f.Advisor].Name)
			p.AdvisorPortrait = w.Generals[f.Advisor].Portrait
		}
		players = append(players, p)
	}
	return players
}

// launcherScenarioName 回傳劇本標題。
//
// **從資料取，不要憑印象編**（CLAUDE.md §7 第 5 條）：標題就在區塊 +0x40，
// 原版四槽視窗畫的也是它（docs/re/52 §4）。先前這裡硬編四個名字，
// 而原版寫的是「第一章・　　「呂布歸天」之卷」——連格式都不一樣。
func (g *game) launcherScenarioName(index int) string {
	if g != nil && g.scenarioTitles != nil {
		if t, ok := g.scenarioTitles[index]; ok && t != "" {
			return t
		}
	}
	return fmt.Sprintf("劇本 %d", index+1)
}

func inspectLauncherSlots(path string) []launcherSlot {
	slots := make([]launcherSlot, 4)
	for i := range slots {
		slots[i] = launcherSlot{Slot: i, Label: "空白槽位"}
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return slots
	}
	for i := range slots {
		w, err := state.LoadScenario(path, i)
		if err != nil || !validLauncherPlayer(w, w.Player) {
			continue
		}
		slots[i].Available = true
		slots[i].Title = big5(w.Title)
		slots[i].Year, slots[i].Month, slots[i].Day = w.Clock.Year, w.Clock.Month, w.Clock.Day
		slots[i].Label = fmt.Sprintf("%d年%d月%d日　%s", w.Clock.Year, w.Clock.Month,
			w.Clock.Day, big5(w.LordName(w.Player)))
	}
	return slots
}

func hasAvailableLauncherSlot(slots []launcherSlot) bool {
	for _, slot := range slots {
		if slot.Available {
			return true
		}
	}
	return false
}

// launcherNewGamePath **刻意忽略既有的 overlay**：開新遊戲一律從
// `SINARIO.DAT` 建，overlay 只是「讀取存檔」的來源或之後的寫入目標。
func launcherNewGamePath(sourceFile, overlay string) string {
	_ = overlay
	return sourceFile
}

func (g *game) updateLauncher() error {
	if g.launcher == nil {
		return nil
	}
	if g.naming != nil {
		// 命名視窗開著時輸入全部歸它（原版 `sub_18FC9` 自己的等待迴圈）。
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.naming.click(ebiten.CursorPosition())
		}
		if pressed(ebiten.KeyEscape) || g.cancelled() {
			g.naming.cancel = true
		}
		g.settleNaming()
		return nil
	}
	x, y := ebiten.CursorPosition()
	if !g.launcher.pointerSeen {
		g.launcher.pointerSeen = true
		g.launcher.pointerX, g.launcher.pointerY = x, y
	} else if x != g.launcher.pointerX || y != g.launcher.pointerY {
		g.launcher.pointerX, g.launcher.pointerY = x, y
		if row, ok := g.launcher.pointerRow(x, y); ok {
			g.launcher.selectRow(row)
		}
	}
	switch g.launcher.phase {
	case launcherSelectPlayer:
		if handled, err := g.updateLordCardPointer(); handled {
			return err
		}
	case launcherSelectFaction:
		if handled, err := g.updateFactionListPointer(); handled {
			return err
		}
	}
	if g.launcher.phase != launcherSelectPlayer &&
		g.launcher.phase != launcherSelectFaction &&
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if row, ok := g.launcher.pointerRow(ebiten.CursorPosition()); ok {
			g.launcher.selectRow(row)
			if result := g.launcher.apply(launcherConfirm); result.kind != launcherNoResult {
				return g.applyLauncherResult(result)
			}
		}
		return nil
	}
	var action launcherAction
	switch {
	case pressed(ebiten.KeyArrowUp):
		action = launcherMoveUp
	case pressed(ebiten.KeyArrowDown):
		action = launcherMoveDown
	case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
		action = launcherConfirm
	case pressed(ebiten.KeyEscape):
		action = launcherCancel
	}
	// 右鍵在原版就是「退回上一層」（`sub_121E7` 回 CF=1）。
	if action == 0 && g.cancelled() {
		action = launcherCancel
	}
	if action == 0 {
		return nil
	}
	return g.applyLauncherResult(g.launcher.apply(action))
}

// updateLordCardPointer 是選君主那一頁的滑鼠：**只認卡片自己的兩個熱區**
// （docs/spec/27 §2.1）。回傳 handled=true 表示這一幀的滑鼠已經處理完，
// 不要再走鍵盤那一段。
func (g *game) updateLordCardPointer() (bool, error) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false, nil
	}
	return true, g.applyLordCardHotspot(lordCardHotspotAt(ebiten.CursorPosition()))
}

// applyLordCardHotspot 是上面那一支的決定部分，**與讀滑鼠分開**。
//
// ⚠ 理由與 `cancelled()` 那一支相同：`inpututil` 讀的是 Ebiten 的全域輸入
// 狀態，無頭測試裡永遠是 false。分開之後測試驗的是行為，
// 不是「有沒有寫這一行」。
func (g *game) applyLordCardHotspot(h lordCardHotspot) error {
	switch h {
	case lordCardConfirm:
		return g.applyLauncherResult(g.launcher.apply(launcherConfirm))
	case lordCardCustom:
		// 「自定」＝ 軍師命名（docs/spec/104）。
		if err := g.openNaming(); err != nil {
			g.launcher.notice = fmt.Sprintf("開不了命名視窗：%v", err)
		}
		return nil
	default:
		// 原版：其他位置的點擊回到等待迴圈，什麼都不做。
		return nil
	}
}

func (g *game) applyLauncherResult(result launcherResult) error {
	switch result.kind {
	case launcherPreviewScenario:
		w, err := state.LoadScenario(g.sourceFile, result.scenario)
		if err != nil {
			g.launcher.notice = fmt.Sprintf("讀取劇本失敗：%v", err)
			return nil
		}
		if !g.launcher.setScenarioPlayers(result.scenario, g.launcherScenarioName(result.scenario), launcherPlayers(w)) {
			return nil
		}
		g.launcherPreviewWorld = w
		g.customAdvisor = nil
	case launcherStartNewGame:
		if err := g.startWorld(launcherNewGamePath(g.sourceFile, g.saveFile), result.scenario, result.player, true, true); err != nil {
			g.launcher.notice = fmt.Sprintf("開始新遊戲失敗：%v", err)
			return nil
		}
		if c := g.customAdvisor; c != nil && g.world != nil {
			// 原版「確定」那一步：勢力 +0x02 寫 0x7F，名字與肖像進區塊欄位。
			g.world.SetCustomAdvisor(c.portrait, c.name)
		}
		g.launcher = nil
	case launcherStartLoad:
		if err := g.startWorld(g.saveFile, result.slot, -1, false, false); err != nil {
			g.launcher.notice = fmt.Sprintf("讀取存檔失敗：%v", err)
			return nil
		}
		g.launcher = nil
	case launcherSetLanguage:
		if result.lang < 0 || result.lang >= len(langpack.Choices) {
			return nil
		}
		// 換完停在原地：**下一列就是換好的樣子**，玩家看得到自己選了什麼。
		if err := g.setLanguage(langpack.Choices[result.lang].Lang); err != nil {
			g.launcher.notice = fmt.Sprintf("換語言失敗：%v", err)
		}
	}
	return nil
}

// drawLauncherCaption 畫一行疊在大地圖上的字。
//
// ⚠ **一定要先鋪底。** 這幾行（劇本名、操作提示、錯誤訊息）都是 remake 加的，
// 原版沒有；背景換成大地圖之後，直接畫上去在草地與河流上完全讀不出來。
func (g *game) drawLauncherCaption(screen *ebiten.Image, text string, x, y int, col color.RGBA) {
	if text == "" {
		return
	}
	w := textdraw.StringWidth(text)
	vector.DrawFilledRect(screen, float32(x-4), float32(y-2),
		float32(w+8), float32(textdraw.GlyphH+4), color.Black, false)
	g.td.Draw(screen, text, x, y, col)
}

func (g *game) drawLauncher(screen *ebiten.Image) {
	// ⭐ **背景是大地圖**，鏡頭固定在 (170, 98)——原版 `sub_11A6E` 在進
	// 新遊戲選單前就把 40×23 的顯示清單畫好了（docs/spec/79 §1.1.1）。
	// **不帶 marks**：`sub_1D615` 只複製圖塊編號，所以據點沒有勢力徽記。
	// 底色先鋪一層，地圖畫不出來（缺素材）時才看得到。
	vector.DrawFilledRect(screen, 0, 0, screenW, screenH, color.RGBA{12, 14, 26, 255}, false)
	if g.lib != nil {
		if img, err := g.lib.RenderWorld(launcherCamX, launcherCamY,
			viewCols, viewRows, launcherSeason); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, strategyMapY)
			screen.DrawImage(ebiten.NewImageFromImage(img), op)
		}
		if banner, err := g.lib.Banner(launcherSeason); err == nil {
			screen.DrawImage(ebiten.NewImageFromImage(banner), nil)
		}
	}
	if g.td == nil || g.launcher == nil {
		return
	}
	white := chrome.Paper
	dim := color.RGBA{200, 200, 210, 255}
	amber := color.RGBA{240, 200, 120, 255}

	// 選君主這一頁換成原版版面的卡片（docs/spec/27）：一次一個勢力，
	// 上下鍵換。**不畫 launcher 自己的大框**——兩個框疊起來很難看，
	// 而原版這一頁本來就只有那一個框。
	l := g.launcher
	if l.phase == launcherSelectFaction {
		g.drawLauncherCaption(screen, g.launcherScenarioName(l.scenario),
			factionListWinX, factionListWinY-textdraw.GlyphH-8, amber)
		g.drawFactionList(screen)
		g.drawLauncherCaption(screen, "點一列選君主　↑↓ 移動　Enter 決定　ESC 返回",
			factionListWinX, factionListWinY+factionListWinH+8, dim)
		g.drawLauncherCaption(screen, l.notice, factionListWinX,
			factionListWinY+factionListWinH+8+textdraw.GlyphH+2,
			color.RGBA{255, 180, 180, 255})
		return
	}
	if l.phase == launcherSelectPlayer && l.cursor >= 0 && l.cursor < len(l.players) {
		g.drawLordCard(screen, l.players[l.cursor], launcherSeason)
		if g.naming != nil {
			g.drawNaming(screen, launcherSeason)
			return
		}
		g.drawLauncherCaption(screen, g.launcherScenarioName(l.scenario),
			lordCardX, lordCardY-textdraw.GlyphH-8, amber)
		g.drawLauncherCaption(screen, fmt.Sprintf("↑↓ 換君主（%d／%d）　Enter 決定　ESC 回清單",
			l.cursor+1, len(l.players)), lordCardX-88, lordCardY+lordCardH+8, dim)
		g.drawLauncherCaption(screen, l.notice, lordCardX-88,
			lordCardY+lordCardH+8+textdraw.GlyphH+2,
			color.RGBA{255, 180, 180, 255})
		return
	}
	g.chrome.Window(screen, launcherPanelX, launcherPanelY, launcherPanelW, launcherPanelH, chrome.Menu)

	drawRows := func(rows []string, startY int, selected int, rowH int) {
		for i, label := range rows {
			y := startY + i*rowH
			if i == selected {
				vector.DrawFilledRect(screen, float32(launcherListX), float32(y),
					float32(launcherListW), float32(rowH), chrome.Select, false)
			}
			col := white
			if i != selected {
				col = dim
			}
			g.td.Draw(screen, label, launcherListX+16, y+4, col)
		}
	}

	switch l.phase {
	case launcherTitle:
		rows := []string{"NEW GAME"}
		if l.hasSave {
			rows = append(rows, "LOAD DATA")
		}
		rows = append(rows, "LANGUAGE")
		drawRows(rows, 152, l.cursor, 32)
	case launcherLanguage:
		g.td.Draw(screen, "LANGUAGE", launcherListX+16, 88, amber)
		rows := make([]string, len(langpack.Choices))
		for i, c := range langpack.Choices {
			// ⚠ 記號要挑**倚天 Big5 有的字**：`▶` 不在裡面，
			// 畫出來是一個方框，看起來像缺字不像游標。
			mark := "　"
			if c.Lang == uiLang.Lang() {
				mark = "●"
			}
			rows[i] = mark + c.Name
		}
		drawRows(rows, launcherListY, l.cursor, launcherRowH)
	case launcherNewGameConfirm:
		g.td.Draw(screen, "NEW GAME", launcherListX+16, 120, amber)
		g.td.Draw(screen, "開始新遊戲？", launcherListX+16, 136, dim)
		drawRows([]string{"YES", "NO"}, 184, l.cursor, 32)
	case launcherScenario:
		g.td.Draw(screen, "選擇劇本", launcherListX+16, 88, amber)
		rows := make([]string, 4)
		for i := range rows {
			rows[i] = g.launcherScenarioName(i)
		}
		drawRows(rows, launcherListY, l.cursor, launcherRowH)
	case launcherGameConfirm:
		g.td.Draw(screen, "確認新遊戲", launcherListX+16, 112, amber)
		lord := l.confirmLord
		if l.confirmedPlayer < 0 {
			lord = "（無效玩家）"
		}
		g.td.Draw(screen, l.scenarioName, launcherListX+16, 136, white)
		g.td.Draw(screen, lord, launcherListX+16, 152, white)
		drawRows([]string{"開始", "返回"}, 192, l.cursor, 32)
	case launcherLoad:
		g.td.Draw(screen, "LOAD DATA", launcherListX+16, 72, amber)
		rows := make([]string, len(l.slots))
		for i, slot := range l.slots {
			rows[i] = fmt.Sprintf("第 %d 槽　%s", i+1, slot.Label)
		}
		drawRows(rows, launcherLoadListY, l.cursor, launcherLoadRowH)
	}
	if l.notice != "" {
		g.td.Draw(screen, l.notice, launcherListX, launcherNoticeY, color.RGBA{255, 180, 180, 255})
	}
	g.td.Draw(screen, launcherHint, launcherListX, launcherHintY, dim)
}

