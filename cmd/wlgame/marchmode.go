package main

// 行軍指示的第二段：戰鬥指揮／委任／解體（docs/spec/39）。
//
// 原版 `sub_17FDB` 選完目的地之後跳 TALK #21「向{2}移動下。請下達戰鬥指示。」
// 再開一個 `cx = 0x4Ch` 的選單——**那個 `cx` 就是 TALK #76 的索引**，
// 三行字串正是三個選項。項目數平常是 2，**目標據點是自己的首都時才給 3**。

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/talkmenu"
)

const (
	// marchModeTalk 是選單字串的 TALK 索引（原版 `sub_1804E` 的 `cx`）。
	marchModeTalk = 0x4C
	// marchModePromptTalk 是 #21「向{2}移動下。請下達戰鬥指示。」，
	// 由 `sub_17FDB` 在開選單之前掛到左下角的狀態列框（docs/spec/140）。
	marchModePromptTalk = 0x15
)

// marchModeState 是選單開著時的狀態。
type marchModeState struct {
	active bool
	corps  int
	dest   int
	row    int
	// rows 是這一次可以選的項目數：2 或 3。
	rows int
	// x, y 是選單框的左上角。**開選單那一刻由游標算**，之後不再變——
	// 原版讀的 `word_19896`／`word_19898` 也是那一刻的值（docs/spec/39 §3.7）。
	x, y int
}

// 位置的兩道夾制，照抄 `sub_1804E`（docs/spec/39 §3.5）：
//
//	dl > 21h → dl = 21h      ; 33 × 16 + 112 ＝ 640
//	bl > 18h − 項數 → 夾住    ; (24−n) × 16 + (n+1) × 16 ＝ 400
//
// ⭐ 兩個上限都是「右／下緣剛好貼齊畫面」——夾制值自己把框的大小說出來了。
const (
	marchModeMaxCol  = 0x21
	marchModeRowSpan = 0x18
)

// marchModeAnchor 把游標的像素座標換成選單框的左上角。
func marchModeAnchor(cursorX, cursorY, rows int) (int, int) {
	col, row := cursorX/16, cursorY/16
	if col < 0 {
		col = 0
	}
	if row < 0 {
		row = 0
	}
	if col > marchModeMaxCol {
		col = marchModeMaxCol
	}
	if max := marchModeRowSpan - rows; row > max {
		row = max
	}
	return col * 16, row * 16
}

// beginMarchMode 在行軍目標定下來之後開選單。
func (g *game) beginMarchMode(corps, dest int) {
	rows := 2
	if g.world.DisbandAllowed(corps) {
		rows = 3 // ★ 目標是首都才有「解體」
	}
	cx, cy := cursorPosition()
	x, y := marchModeAnchor(cx, cy, rows)
	g.marchMode = marchModeState{
		active: true, corps: corps, dest: dest, rows: rows, x: x, y: y,
	}
	// 訊息 #21 走左下角的狀態列框，與選單是兩個獨立視窗（docs/spec/140）。
	//
	// ⚠ key 是 **ASCII `'2'`**，不是數值 2——`Part.Marker` 存的是原版
	// `\` 後面那個字元（`internal/assets/text/talk.go`）。
	g.setStatusTalk(marchModePromptTalk, map[byte]string{'2': g.cityName(dest)})
}

// cityName 是據點名，編號不合法時回空字串。
func (g *game) cityName(i int) string {
	if g == nil || g.world == nil || i < 0 || i >= len(g.world.Cities) {
		return ""
	}
	return big5(g.world.Cities[i].Name)
}

// marchModeLabels 取 TALK #76 的三行。取不到就退回內建字串——
// 缺原版素材時要能跑，不是整個動不了。
//
// ⛔ 走 `talkmenu.MenuLabels` 不是 `talkLines`：後者會 `TrimRight` 掉行尾的
// 全形空白，而**框寬由第一列的全形字數決定**，砍掉就窄 16 px（docs/spec/124）。
func (g *game) marchModeLabels() []string {
	fallback := []string{"　戰鬥指揮　", "　委　　任　", "　解　　體　"}
	if g == nil || g.lib == nil {
		return fallback
	}
	return talkmenu.MenuLabels(g.lib.Talk, marchModeTalk, nil, fallback)
}

// marchModeRows 是這一次真的畫出來的幾列。
func (g *game) marchModeRows() []string {
	labels := g.marchModeLabels()
	if n := g.marchMode.rows; n >= 0 && n < len(labels) {
		labels = labels[:n]
	}
	return labels
}

func (g *game) updateMarchMode() {
	m := &g.marchMode
	cancel := func() {
		// 原版右鍵是**回去重選據點**，不是整條流程取消。
		m.active = false
		g.clearStatusTalk()
		g.pickDestination(m.corps)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		cancel()
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if row, ok := g.talkChoiceClick(m.x, m.y, g.marchModeRows()); ok {
			m.row = row
			g.commitMarchMode()
		}
		return
	}
	switch {
	case pressed(ebiten.KeyEscape):
		cancel()
	case pressed(ebiten.KeyArrowUp):
		m.row = (m.row + m.rows - 1) % m.rows
	case pressed(ebiten.KeyArrowDown):
		m.row = (m.row + 1) % m.rows
	case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
		g.commitMarchMode()
	default:
		for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
			if i < m.rows && pressed(k) {
				m.row = i
				g.commitMarchMode()
				return
			}
		}
	}
}

func (g *game) commitMarchMode() {
	m := &g.marchMode
	mode := state.MarchMode(m.row)
	if err := g.world.SetMarchMode(m.corps, mode); err != nil {
		g.setEvent(err.Error())
		return
	}
	name := big5(g.world.Generals[g.world.Leader(m.corps)].Name)
	dest := g.cityName(m.dest)
	if dest == "" {
		dest = "原地"
	}
	g.lastEvent = name + " 向 " + dest + " 行軍（" + mode.String() + "）"
	g.finishMarchOrder()
}

func (g *game) drawMarchMode(screen *ebiten.Image) {
	m := &g.marchMode
	if !m.active {
		return
	}
	// 與指令列那三張彈出選單同一份幾何（`sub_193E9`，docs/spec/126）：
	// 框寬 ＝ (第一列的全形字數 + 1) × 16 ＝ 112、框高 ＝ (項目數 + 1) × 16。
	g.drawLegacyChoiceBox(screen, m.x, m.y, g.marchModeRows(), m.row)
}

// demoMarchMode 是**驗收用**的捷徑：編一支軍團、對首都下行軍，
// 停在三選一那一格（第三項「解體」因此會出現）。正常玩不會走到這裡。
// demoMarchList 是**驗收用**的捷徑：編一支軍團，停在行軍目的地一覽。
// 與 demoMarchMode 差在少走最後一步（選完目的地才有三選一）。
func (g *game) demoMarchList() {
	rows := g.formCandidates()
	if len(rows) == 0 {
		return
	}
	kinds, manned := g.affordable()
	leader := rows[0]
	if err := g.world.FormCorps(leader, kinds, manned); err != nil {
		g.setEvent(err.Error())
		return
	}
	g.pickDestination(leader)
}

// demoCorpsOnMap 編一支軍團之後**停在大地圖**（docs/spec/74 §4.1）。
//
// marchTo ≥ 0 就再下一道行軍指示，讓 `-shot-frames` 推進的 tick
// 把它帶出城——軍團待在自己城裡時疊在據點中心徽記上，
// 兩者都是紅色系，肉眼分不出來（原版也一樣）。
//
// ⚠ 走 `formCandidates()` 的真實資格判定，不要自己抄一份：
// 驗收路徑與遊戲跑的是同一條規則，才驗得到規則本身。
func (g *game) demoCorpsOnMap(marchTo int) {
	rows := g.formCandidates()
	if len(rows) == 0 {
		log.Print("⚠ -corps-on-map：沒有可編成的武將，這一張截圖上不會有軍團")
		return
	}
	kinds, manned := g.affordable()
	leader := rows[0]
	if err := g.world.FormCorps(leader, kinds, manned); err != nil {
		log.Printf("⚠ -corps-on-map：編成失敗（%v）", err)
		return
	}
	if marchTo < 0 {
		return
	}
	if marchTo >= len(g.world.Cities) {
		log.Printf("⚠ -march-to %d 超出據點範圍（0–%d），只編成不下令",
			marchTo, len(g.world.Cities)-1)
		return
	}
	if err := g.world.March(leader, marchTo); err != nil {
		log.Printf("⚠ -march-to %d：下不了行軍指示（%v）", marchTo, err)
	}
}

func (g *game) demoMarchMode(to int) {
	rows := g.formCandidates()
	if len(rows) == 0 {
		return
	}
	kinds, manned := g.affordable()
	leader := rows[0]
	if err := g.world.FormCorps(leader, kinds, manned); err != nil {
		g.setEvent(err.Error())
		return
	}
	dest := g.world.Factions[g.world.Player].Capital
	if to >= 0 && to < len(g.world.Cities) {
		dest = to // ← 不是首都 ⇒ 只有兩項
	}
	if err := g.world.March(leader, dest); err != nil {
		g.setEvent(err.Error())
		return
	}
	g.beginMarchMode(leader, dest)
}
