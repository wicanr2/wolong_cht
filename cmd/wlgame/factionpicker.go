package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

// 縮小地圖那兩個熱區：點地圖捲鏡頭（`0x16`）與 22 勢力的選擇視窗（`0x17`）。
// 版面與規則全部照原版，出處在 docs/spec/35 §2.5、docs/re/71 §3–§4。

// 熱區 0x16：縮小地圖的地圖區（440, 40, 192×128）。
//
// ⭐ **1 px 對 2 格**——世界 384×256 格畫成 192×128 px（docs/re/62 §1）。
const minimapCellsPerPixel = 2

// minimapWorldAt 把螢幕座標換成世界格，並回報有沒有落在地圖區裡。
// 原版 `sub_15AB6` 減完原點會**借位夾成 0**，這裡照抄。
func minimapWorldAt(x, y int) (col, row int, ok bool) {
	dx := x - strategyMinimapX
	dy := y - strategyMinimapY
	if dx < 0 || dy < 0 || dx >= strategyMinimapW || dy >= strategyMinimapImageH {
		return 0, 0, false
	}
	return dx * minimapCellsPerPixel, dy * minimapCellsPerPixel, true
}

// centreCamOn 把鏡頭移到某一格，用的是**與開局定位同一支**的偏移。
//
// 原版 `sub_121B2` 呼叫 `sub_12151(ax=0x14, cx=0x0C)`，而 `−(20,12)`
// 正是開局鏡頭那一組立即值（docs/spec/52）。共用同一段換算，
// 不要各寫一次——那是 CLAUDE.md §7 第 6 條講的「一條規則只留一份實作」。
func (g *game) centreCamOn(col, row int) {
	g.camX = col - centreCol
	g.camY = row - centreRow
	g.clampCam()
}

// 勢力選擇視窗的版面（docs/spec/35 §2.5.2，全部是原版直接座標）。
const (
	pickerWinX = 512 // sub_1895D 的 dx=0x20 × 16
	pickerWinY = 192 // (bx=0x0A + 2) × 16
	pickerWinW = 128 // cx 低 byte 0x08 × 2 × 8
	pickerWinH = 208 // cx 高 byte 0x0D × 2 × 8

	// 兩塊填色與熱區都比外框四邊各內縮 8 px——那 8 px 就是框邊。
	pickerTitleX  = 544 // sub_106F5 的 dx=0x220
	pickerTitleY  = 200 // bx=0xC8
	pickerListX   = 520
	pickerListY   = 216 // 0xD8
	pickerListW   = 112 // sub_1E3D7 的 cx 低 byte 0x0E × 8
	pickerListH   = 176 // cx 高 byte 0x16 × 8
	pickerRowH    = 16
	pickerRows    = 11  // 22 個勢力分兩欄
	pickerLeftX   = 524 // sub_15C14 的 dx=0x20C
	pickerRightX  = 580 // dx=0x244
	pickerSplitX  = 576 // cmp cx, 0x240：X ≥ 576 ⇒ 右欄
	pickerSlots   = pickerRows * 2
	pickerEmptyZh = "－－－" // cs:5C0D
	pickerTitleZh = "勢力一覽" // cs:5C04

	// 字色的調色盤索引（docs/re/62 §4.2 的屬性 0x90／0x9A／0x93）。
	pickerInkNormal   = 0
	pickerInkOwn      = 10
	pickerInkSelected = 3
	// 兩塊底色（sub_1F1A3 的 ah）。
	pickerTitleBand = 7
	pickerListBand  = 9
)

// pickerRowRect 回傳第 n 個槽位的文字起點。左欄 0–10、右欄 11–21。
func pickerRowOrigin(n int) (x, y int) {
	x = pickerLeftX
	if n >= pickerRows {
		x, n = pickerRightX, n-pickerRows
	}
	return x, pickerListY + n*pickerRowH
}

// pickerSlotAt 把點擊換成槽位編號。原版先減清單頂端、右欄再加 11 列，
// 然後整個 `>> 4`——所以**欄的分界是 X ≥ 576，不是視窗中線**。
func pickerSlotAt(x, y int) (int, bool) {
	if x < pickerListX || x >= pickerListX+pickerListW ||
		y < pickerListY || y >= pickerListY+pickerListH {
		return 0, false
	}
	n := (y - pickerListY) / pickerRowH
	if x >= pickerSplitX {
		n += pickerRows
	}
	if n < 0 || n >= pickerSlots {
		return 0, false
	}
	return n, true
}

// pickerSelectable 是原版的兩條規則：**勢力要存在**（記錄 +0x00 ≥ 0x80）、
// **不能選自己**（`cmp al, cs:byte_10CFF / jz 忽略`）。
func (g *game) pickerSelectable(n int) bool {
	if g.world == nil || n < 0 || n >= len(g.world.Factions) {
		return false
	}
	return g.world.Factions[n].Alive && n != g.world.Player
}

// updateFactionPicker 是視窗開著時的輸入：左鍵選、右鍵關。
func (g *game) updateFactionPicker() bool {
	if !g.factionPicker {
		return false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || pressed(ebiten.KeyEscape) {
		g.factionPicker = false
		return true
	}
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	x, y := ebiten.CursorPosition()
	if n, ok := pickerSlotAt(x, y); ok && g.pickerSelectable(n) {
		g.minimapFaction = n
	}
	// ⚠ 點在視窗裡但不合規則就**什麼都不做**，視窗不關——照原版
	// （`jb 忽略` 之後 `jmp loc_15B7B` 回到等待迴圈）。
	return true
}

// drawFactionPicker 畫視窗。三塊：框、標題列、清單。
func (g *game) drawFactionPicker(dst *ebiten.Image) {
	if !g.factionPicker || g.world == nil {
		return
	}
	g.chrome.Window(dst, pickerWinX, pickerWinY, pickerWinW, pickerWinH, chrome.Menu)

	band := func(x, y, w, h, index int, fallback color.RGBA) {
		vector.DrawFilledRect(dst, float32(x), float32(y), float32(w), float32(h),
			g.paletteInk(index, fallback), false)
	}
	band(pickerListX, pickerTitleY, pickerListW, pickerRowH,
		pickerTitleBand, color.RGBA{0, 0, 0, 255})
	band(pickerListX, pickerListY, pickerListW, pickerListH,
		pickerListBand, color.RGBA{255, 221, 153, 255})

	g.td.Draw(dst, pickerTitleZh, pickerTitleX, pickerTitleY,
		g.paletteInk(pickerInkNormal, chrome.Paper))

	for n := 0; n < pickerSlots; n++ {
		x, y := pickerRowOrigin(n)
		if n >= len(g.world.Factions) || !g.world.Factions[n].Alive {
			g.td.Draw(dst, pickerEmptyZh, x, y,
				g.paletteInk(pickerInkNormal, color.RGBA{0, 0, 0, 255}))
			continue
		}
		ink := pickerInkNormal
		switch {
		case n == g.world.Player:
			ink = pickerInkOwn
		case n == g.minimapFaction:
			ink = pickerInkSelected
		}
		// 原版畫的是 `sub_188B0`（沒讀，docs/spec/35 §5），remake 沿用
		// 既有的「勢力 ＝ 君主名」慣例，與外交視窗同一支。
		g.td.Draw(dst, g.diplomacyFactionName(n), x, y,
			g.paletteInk(ink, color.RGBA{0, 0, 0, 255}))
	}
}
