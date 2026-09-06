package main

// 大地圖選點（docs/spec/149）。
//
// 原版下行軍指示時，選目標據點走的是 `sub_1703C`：畫一個空心框游標，
// 等玩家點大地圖上的據點；右鍵取消。**不是一覽表。**

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/world"
)

// mapPickState 是「現在在大地圖上選一格」。
// **pick 是 nil 就是沒在選**——零值安全，不必另外記一個 active 旗標。
type mapPickState struct {
	// pick 在點到據點時呼叫。
	pick func(city int)
	// cancel 是右鍵：原版 `sub_1703C` 回 CF=1，整條流程結束。
	cancel func()
	// tile 是**對拍 fixture 用**的固定格：headless 的指標位置不可控
	// （同 `hideAmountCursor`），釘住之後游標才畫得到指定的那一格。
	// nil ＝ 跟著滑鼠走，一般遊玩一律是 nil。
	tile *image.Point
}

func (g *game) mapPickActive() bool { return g != nil && g.mapPick.pick != nil }

// beginMapPick 進入選點狀態。
func (g *game) beginMapPick(pick func(city int), cancel func()) {
	g.mapPick = mapPickState{pick: pick, cancel: cancel}
}

func (g *game) endMapPick() { g.mapPick = mapPickState{} }

// mapPickTileAt 把螢幕座標換成格座標。
//
// ⭐ remake 的鏡頭本來就是**以格為單位**（`camX`／`camY` 是欄與列），
// 所以沒有原版那個「原點與畫面座標各捨去一次、各吃掉一格」的問題
// （docs/re/85 §2）。地圖區以外回 false。
func (g *game) mapPickTileAt(x, y int) (col, row int, ok bool) {
	if x < 0 || x >= strategyMapW || y < strategyMapY || y >= screenH {
		return 0, 0, false
	}
	return g.camX + x/world.TileSize,
		g.camY + (y-strategyMapY)/world.TileSize, true
}

// cityAtTile 是原版 `sub_1709E` 的第二道閘：**據點表登記的 X／Y 與那一格
// 完全相等**才算命中（docs/spec/149 §1）。
//
// ⭐ 一座城的圖形有 4×4 格，但只有登記的那一格按得到——
// 「點圖示中間沒反應」是原版行為。
//
// 第一道閘（圖塊編號落在 `0CBh`–`0D3h`）不必另外查：圖塊是據點就一定在
// 表裡，反之亦然，所以兩道閘在這裡同義。
func (g *game) cityAtTile(col, row int) (int, bool) {
	if g == nil || g.world == nil {
		return 0, false
	}
	for i := range g.world.Cities {
		c := &g.world.Cities[i]
		if int(c.X) == col && int(c.Y) == row {
			return i, true
		}
	}
	return 0, false
}

// updateMapPick 吃選點期間的滑鼠。回 true 表示這一格輸入被吃掉了。
//
// ⚠ **熱區優先**：原版 `sub_1703C` 先問 `sub_1E453`，游標壓在熱區上時
// 這一圈根本不問據點（docs/re/85 §3）。所以呼叫端要排在 HUD 的命中判定
// **之後**。
func (g *game) updateMapPick() bool {
	if !g.mapPickActive() {
		return false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		cancel := g.mapPick.cancel
		g.endMapPick()
		if cancel != nil {
			cancel()
		}
		return true
	}
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}
	x, y := ebiten.CursorPosition()
	col, row, ok := g.mapPickTileAt(x, y)
	if !ok {
		return false
	}
	city, hit := g.cityAtTile(col, row)
	if !hit {
		return true // 點在地圖上但不是據點：原版就是沒反應
	}
	pick := g.mapPick.pick
	g.endMapPick()
	if pick != nil {
		pick(city)
	}
	return true
}

// 游標的兩層遮罩，**逐像素取自原版擷取**（docs/spec/149 §1.1）：
// 15×15 的白色圓角空心框，加一層黑影。
//
// ⚠ **兩層各存一份，不要拿白的位移一格當黑的。** 影子的圓角步法與本體
// 不一樣（左下、右上各差一格），位移法會差 6 px——而那 6 px 只有
// 逐像素對拍看得到。**照抄就是照抄。**
var (
	// 白：色 15 `(243,243,243)`。原點 ＝ 那一格的左上角。
	mapPickCursorWhite = [15]string{
		".#############..",
		"##...........##.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"#.............#.",
		"##...........##.",
		".#############..",
	}
	// 黑影：比本體多一列、多一行。
	mapPickCursorBlack = [16]string{
		"..............K.",
		"..KKKKKKKKKKK..K",
		".KK..........K.K",
		".K.............K",
		".K.............K",
		".K.............K",
		".K.............K",
		".K.............K",
		".K.............K",
		".K.............K",
		".K.............K",
		".K.............K",
		".K.............K",
		"..K............K",
		"K.............KK",
		".KKKKKKKKKKKKKK.",
	}
)

// mapPickCursorTile 是游標現在停在哪一格。fixture 釘住的話用釘的那一格。
func (g *game) mapPickCursorTile() (col, row int, ok bool) {
	if g.mapPick.tile != nil {
		return g.mapPick.tile.X, g.mapPick.tile.Y, true
	}
	return g.mapPickTileAt(ebiten.CursorPosition())
}

// drawMapPickCursor 把游標畫在**游標所在那一格的左上角**——
// 原版是格對齊的，不跟著像素走（docs/spec/149 §1.1）。
func (g *game) drawMapPickCursor(screen *ebiten.Image) {
	if !g.mapPickActive() {
		return
	}
	col, row, ok := g.mapPickCursorTile()
	if !ok {
		return
	}
	x := (col - g.camX) * world.TileSize
	y := strategyMapY + (row-g.camY)*world.TileSize
	white := g.paletteInk(15, color.RGBA{243, 243, 243, 255})
	black := color.RGBA{0, 0, 0, 255}
	draw := func(rows []string, c color.RGBA) {
		for dy, line := range rows {
			for dx, ch := range line {
				if ch == '.' {
					continue
				}
				vector.DrawFilledRect(screen,
					float32(x+dx), float32(y+dy), 1, 1, c, false)
			}
		}
	}
	draw(mapPickCursorBlack[:], black)
	draw(mapPickCursorWhite[:], white)
}
