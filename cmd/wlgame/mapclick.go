package main

// 大地圖左鍵（docs/spec/151）。
//
// 原版 `sub_11E46` 是一條完整的分派：那一格是據點就開情報卡、有軍團就開
// 軍團情報面板、兩者都有就先跳一張兩項的選單（TALK #80）。
// remake 的大地圖點擊先前什麼都不做。

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/wicanr2/wolong_cht/internal/ui/talkmenu"
)

// mapChoiceTalk 是那張兩項選單的字串索引：#80「　據　　點　」「　軍　團　」。
const mapChoiceTalk = 0x50

// 位置的兩道夾制，照抄 `sub_11F0E`（docs/spec/151 §1.2）。
// ⚠ 與行軍三選一的 `sub_1804E` **不是同一組數字**——框窄一格、少一列。
const (
	mapChoiceMaxCol = 0x23
	mapChoiceMaxRow = 0x16
)

// mapChoiceState 是那張兩項選單開著時的狀態。
// **active 是 false 就是沒開**——零值安全。
//
// ⭐ 原版這一族選單全部走同一支 `sub_193E9`：只有「字串的 TALK 索引、
// 項數、框的左上角」三個參數不同（行軍三選一 `sub_1804E`、
// 這一張 `sub_11F0E`、離開確認 `docs/spec/153`）。
type mapChoiceState struct {
	active bool
	city   int
	corps  []int
	row    int
	x, y   int
}

// mapChoiceAnchor 把游標的像素座標換成選單框的左上角。
func mapChoiceAnchor(cursorX, cursorY int) (int, int) {
	col, row := cursorX/16, cursorY/16
	if col < 0 {
		col = 0
	}
	if row < 0 {
		row = 0
	}
	if col > mapChoiceMaxCol {
		col = mapChoiceMaxCol
	}
	if row > mapChoiceMaxRow {
		row = mapChoiceMaxRow
	}
	return col * 16, row * 16
}

// corpsAtTile 是那一格上的軍團（原版建表 callback `0x17217`，docs/spec/151 §1.1）。
//
// ⭐ **不分敵我**：原版那一支讀了玩家勢力卻從頭到尾沒比對它，
// 所以點得到別人的軍團——「別人的」由狀態列 #4 那條路處理。
func (g *game) corpsAtTile(col, row int) []int {
	if g == nil || g.world == nil {
		return nil
	}
	var out []int
	for i := range g.world.Corps {
		c := &g.world.Corps[i]
		if c.Alive && int(c.X) == col && int(c.Y) == row {
			out = append(out, i)
		}
	}
	return out
}

// mapClickBusy 回答「現在有沒有別的東西在吃輸入」。
func (g *game) mapClickBusy() bool {
	return g.list != nil || g.adviseActive() || g.form.active ||
		g.marchMode.active || g.finance.active || g.saveUI.active ||
		g.messageActive() || g.mapPickActive() || g.mapChoice.active ||
		g.cityInfo.active || g.corpsInfo.active || g.popupMenuActive()
}

// updateMapClick 吃大地圖的左鍵。回 true 表示這一格輸入被吃掉了。
//
// ⚠ **排在熱區與各視窗之後**，與地圖選點同一個位置（docs/spec/149 §2）。
func (g *game) updateMapClick() bool {
	if g.world == nil || g.mapClickBusy() {
		return false
	}
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}
	x, y := ebiten.CursorPosition()
	col, row, ok := g.mapPickTileAt(x, y)
	if !ok {
		return false
	}
	return g.dispatchMapClick(col, row, x, y)
}

// dispatchMapClick 是 `sub_11E46` 的分派本體。抽出來是為了讓驗收 fixture
// （`-map-click`）走**同一條路**，而不是自己再實作一次。
//
// cursorX／cursorY 只用來決定那張兩項選單開在哪裡（原版讀的也是那一刻的
// 游標位置，docs/spec/151 §1.2）。
func (g *game) dispatchMapClick(col, row, cursorX, cursorY int) bool {
	city, hasCity := g.cityAtTile(col, row)
	corps := g.corpsAtTile(col, row)
	switch {
	case hasCity && len(corps) > 0:
		// 兩者都有：先問要看哪一個（原版 `sub_11F0E`）。
		cx, cy := mapChoiceAnchor(cursorX, cursorY)
		g.mapChoice = mapChoiceState{
			active: true, city: city, corps: corps, x: cx, y: cy,
		}
	case hasCity:
		g.openCityInfo(city)
	case len(corps) > 0:
		g.openCorpsOnTile(corps)
	default:
		return false // 空地：原版也是沒反應
	}
	return true
}

// openCorpsOnTile 開「那一格上的軍團」一覽。
//
// ⭐ **選完回一覽表**（原版 `sub_11E46` 的 `jmp loc_11EE8` 迴圈）：
// 同一格上有好幾支時可以連著看，右鍵才離開。
func (g *game) openCorpsOnTile(rows []int) {
	if len(rows) == 0 {
		return
	}
	g.openCorpsListWith(rows, "選擇要看的軍團　Enter 選取／決定　ESC 取消",
		func(i int) bool {
			g.enterCorpsFromMap(i)
			return false
		})
}

// enterCorpsFromMap 是原版 `sub_17F90` 的兩條路（docs/spec/149 §1.3）：
// 自己的軍團接著進行軍指示（狀態列 #3），別人的只是看（#4）。
func (g *game) enterCorpsFromMap(corps int) {
	if corps < 0 || corps >= len(g.world.Corps) {
		return
	}
	if int(g.world.Corps[corps].Faction) == g.world.Player {
		g.showCorpsPanel(corps)
		g.pickDestination(corps)
		return
	}
	g.openCorpsInfo(corps)
}

// mapChoiceRows 取 TALK #80 的兩行。
//
// ⛔ 走 `talkmenu.MenuLabels` 不是 `talkLines`：後者會 `TrimRight` 掉行尾的
// 全形空白，而**框寬由第一列的全形字數決定**（docs/spec/124）。
func (g *game) mapChoiceRows() []string {
	fallback := []string{"　據　　點　", "　軍　團　"}
	if g == nil || g.lib == nil {
		return fallback
	}
	return talkmenu.MenuLabels(g.lib.Talk, mapChoiceTalk, nil, fallback)
}

func (g *game) updateMapChoice() {
	m := &g.mapChoice
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) ||
		pressed(ebiten.KeyEscape) {
		m.active = false
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if row, ok := g.talkChoiceClick(m.x, m.y, g.mapChoiceRows()); ok {
			m.row = row
			g.commitMapChoice()
		}
		return
	}
	switch {
	case pressed(ebiten.KeyArrowUp), pressed(ebiten.KeyArrowDown):
		m.row = 1 - m.row
	case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
		g.commitMapChoice()
	case pressed(ebiten.Key1):
		m.row = 0
		g.commitMapChoice()
	case pressed(ebiten.Key2):
		m.row = 1
		g.commitMapChoice()
	}
}

func (g *game) commitMapChoice() {
	m := &g.mapChoice
	city, corps, row := m.city, m.corps, m.row
	m.active = false
	if row == 1 {
		g.openCorpsOnTile(corps)
		return
	}
	g.openCityInfo(city)
}

func (g *game) drawMapChoice(screen *ebiten.Image) {
	if !g.mapChoice.active {
		return
	}
	g.drawLegacyChoiceBox(screen, g.mapChoice.x, g.mapChoice.y,
		g.mapChoiceRows(), g.mapChoice.row)
}
