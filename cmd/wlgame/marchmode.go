package main

// 行軍指示的第二段：戰鬥指揮／委任／解體（docs/spec/39）。
//
// 原版 `sub_17FDB` 選完目的地之後跳 TALK #21「向{2}移動下。請下達戰鬥指示。」
// 再開一個 `cx = 0x4Ch` 的選單——**那個 `cx` 就是 TALK #76 的索引**，
// 三行字串正是三個選項。項目數平常是 2，**目標據點是自己的首都時才給 3**。

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// marchModeTalk 是選單字串的 TALK 索引（原版 `sub_1804E` 的 `cx`）。
const marchModeTalk = 0x4C

// marchModeState 是選單開著時的狀態。
type marchModeState struct {
	active bool
	corps  int
	dest   int
	row    int
	// rows 是這一次可以選的項目數：2 或 3。
	rows int
}

// 選單的版面。**原版的 `sub_193E9` 版面沒解**（docs/spec/39 §4），
// 這一組是 remake 自訂的，沿用兩選一選單那一套比例。
const (
	marchModeWinX, marchModeWinY = 168, 128
	marchModeWinW, marchModeWinH = 304, 128
	marchModeRowX                = marchModeWinX + chrome.Tile + 4
	marchModeRowY                = marchModeWinY + chrome.Tile + 2*(textdraw.GlyphH+4) + 4
	marchModeRowW                = marchModeWinW - 2*(chrome.Tile+4)
	marchModeRowH                = textdraw.GlyphH + 2
	marchModeRowPitch            = textdraw.GlyphH + 6
)

func marchModeRowRect(i int) image.Rectangle {
	y := marchModeRowY + i*marchModeRowPitch
	return image.Rect(marchModeRowX, y, marchModeRowX+marchModeRowW, y+marchModeRowH)
}

func marchModeRowAt(x, y, rows int) (int, bool) {
	p := image.Pt(x, y)
	for i := 0; i < rows; i++ {
		if p.In(marchModeRowRect(i)) {
			return i, true
		}
	}
	return 0, false
}

// beginMarchMode 在行軍目標定下來之後開選單。
func (g *game) beginMarchMode(corps, dest int) {
	rows := 2
	if g.world.DisbandAllowed(corps) {
		rows = 3 // ★ 目標是首都才有「解體」
	}
	g.marchMode = marchModeState{active: true, corps: corps, dest: dest, rows: rows}
}

// marchModeLabels 取 TALK #76 的三行。取不到就退回內建字串——
// 缺原版素材時要能跑，不是整個動不了。
func (g *game) marchModeLabels() []string {
	if lines, ok := g.talkLines(marchModeTalk, nil); ok && len(lines) >= 3 {
		return lines[:3]
	}
	return []string{"　戰鬥指揮　", "　委　　任　", "　解　　體　"}
}

func (g *game) updateMarchMode() {
	m := &g.marchMode
	cancel := func() {
		// 原版右鍵是**回去重選據點**，不是整條流程取消。
		m.active = false
		g.pickDestination(m.corps)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		cancel()
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if row, ok := marchModeRowAt(x, y, m.rows); ok {
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
	dest := "原地"
	if m.dest >= 0 && m.dest < len(g.world.Cities) {
		dest = big5(g.world.Cities[m.dest].Name)
	}
	g.lastEvent = name + " 向 " + dest + " 行軍（" + mode.String() + "）"
	m.active = false
}

func (g *game) drawMarchMode(screen *ebiten.Image) {
	m := &g.marchMode
	if !m.active {
		return
	}
	g.chrome.Window(screen, marchModeWinX, marchModeWinY,
		marchModeWinW, marchModeWinH, chrome.Menu)
	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	dim := g.paletteInk(strategyInkDim, color.RGBA{255, 223, 154, 255})

	// 標題就是原版跳的那一則：TALK #21「向{2}移動下。請下達戰鬥指示。」
	dest := ""
	if m.dest >= 0 && m.dest < len(g.world.Cities) {
		dest = big5(g.world.Cities[m.dest].Name)
	}
	head := marchModePromptFallback(dest)
	if lines, ok := g.talkLines(0x15, map[byte]string{2: dest}); ok && len(lines) > 0 {
		head = lines
	}
	for i, line := range head {
		if i >= 2 {
			break
		}
		g.td.Draw(screen, line, marchModeRowX,
			marchModeWinY+chrome.Tile+i*(textdraw.GlyphH+4), dim)
	}

	for i, label := range g.marchModeLabels() {
		if i >= m.rows {
			break
		}
		r := marchModeRowRect(i)
		if i == m.row {
			vector.DrawFilledRect(screen, float32(r.Min.X), float32(r.Min.Y),
				float32(r.Dx()), float32(r.Dy()), chrome.Select, false)
		}
		x := r.Min.X + (r.Dx()-textdraw.StringWidth(label))/2
		g.td.Draw(screen, label, x, r.Min.Y+1, ink)
	}
}

// marchModePromptFallback 是取不到 TALK #21 時的替代。
func marchModePromptFallback(dest string) []string {
	return []string{"向" + dest + "移動下。請下", "達戰鬥指示。"}
}

// demoMarchMode 是**驗收用**的捷徑：編一支軍團、對首都下行軍，
// 停在三選一那一格（第三項「解體」因此會出現）。正常玩不會走到這裡。
func (g *game) demoMarchMode() {
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
	if err := g.world.March(leader, dest); err != nil {
		g.setEvent(err.Error())
		return
	}
	g.beginMarchMode(leader, dest)
}
