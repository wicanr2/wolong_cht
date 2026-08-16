package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 訊息框的版面常數。**原版只有一個框**，有沒有講話的人不改變版面
// （docs/spec/41、docs/re/66）。
//
//	框     sub_1895D(bx=8, dx=0Ah, cx=0510h)  ⇒ (160, 160, 256, 80)
//	內容區 sub_10BCD 四邊各內縮 8            ⇒ (168, 168) 240×64
//	肖像   dx×16 + 72 是右緣，64×64          ⇒ (168, 168)
//	文字   dx×16 + 80 ／ bx×16 + 16          ⇒ (240, 176)
//
// 事件 2／3 的 IVENTGRF 288×176 場景（`sub_13D09`／`sub_13D68`）另計，
// 那兩個講話位置目前只有機器碼、沒有影格覆驗。
const (
	talkBoxX = 160
	talkBoxY = 160
	talkBoxW = 256
	talkBoxH = 80

	talkPortraitX = 168
	talkPortraitY = 168
	// 文字讓開 64 px 的肖像：框內 +8 起、+72 結束，字從 +80 開始。
	talkTextX = 240
	talkTextY = 176
	// 每列 10 個全形字。TALK.DAT 的原文就折在這個寬度上（825 行剛好
	// 160 px），遊戲本身不換行。
	talkTextWidth = 160
	talkLinePitch = 16 // 全形字高；sub_106F9 每畫一個字前進 16 px
	// talkBoxRows 是框裡放得下幾列：框高 80、上下各內縮 8、一列 16 px ⇒ 4 列。
	talkBoxRows = (talkBoxH - 16) / talkLinePitch

	// defaultPortraitPage 是一般通知的肖像（原版 `sub_18810` 的呼叫端
	// 幾乎清一色 `mov al, 93h`）。KAOGRF 第 147 張。
	defaultPortraitPage = 0x93

	// 事件場景（IVENTGRF 插圖）上的兩個講話框。**寬高與訊息框完全一樣**
	// ——`sub_1075B` 把 cx = 0510h 寫死，位置才由呼叫端決定。
	//
	//	sub_13C99（講話的武將）dx=0, bx=5   ⇒ 框 (0, 80)
	//	sub_13CDC（君主）      dx=8, bx=12h ⇒ 框 (128, 288)
	//
	// 兩個位置都在原版實錄影格上量過（docs/re/66 §5.1）。
	talkUpperBoxX, talkUpperBoxY = 0, 80
	talkLowerBoxX, talkLowerBoxY = 128, 288

	talkSceneX = 64
	talkSceneY = 144

	// 選單框：外交三選一、撥款、說服五選一共用 `sub_13B7E`，
	// 那一組座標與尺寸全是寫死的（docs/spec/45 §2）：
	//
	//	sub_19796(dx=50h, bx=0B0h)  ⇒ 像素 (80, 176)
	//	cx = 600Ah → sub_1FB11      ⇒ 每列 0Ah×2 ＝ 20 B ＝ 160 px、96 列
	//
	// 高 96 ＝ 上下各內縮 8 ＋ 五列 × 16。**選項少的時候框不縮**——
	// `sub_13B5A` 傳 5、`sub_13902` 傳 3，而框三處相同。
	talkChoiceX    = 80
	talkChoiceY    = 176
	talkChoiceW    = 160
	talkChoiceH    = 96
	talkChoiceRows = (talkChoiceH - 2*chrome.Tile) / talkLinePitch
)

// playerLordPortrait 是 composite TALK 的 fallback speaker。事件 3 fixture
// 已證實 sub_187FF 取玩家君主 record，再由 sub_13C99 傳入 General +0x01。
func (g *game) playerLordPortrait() int {
	if g == nil || g.world == nil || g.world.Player < 0 ||
		g.world.Player >= len(g.world.Factions) {
		return -1
	}
	lord := g.world.Factions[g.world.Player].Lord
	if lord < 0 || lord >= len(g.world.Generals) {
		return -1
	}
	return g.world.Generals[lord].Portrait
}

func (g *game) playerTalkVariant() int {
	if g == nil || g.world == nil || g.world.Player < 0 ||
		g.world.Player >= len(g.world.Factions) {
		return 0
	}
	lord := g.world.Factions[g.world.Player].Lord
	if lord < 0 || lord >= len(g.world.Generals) {
		return 0
	}
	return talkVariant(g.world.Generals[lord].TalkVariant)
}

func talkVariant(raw int) int {
	// 原版只做一次「>=3 減 3」；對外部 raw fixture 也以三種變體
	// 收斂，避免負值或異常值把 TALK 索引帶出範圍。
	if raw < 0 {
		return 0
	}
	for raw >= 3 {
		raw -= 3
	}
	return raw
}

func (g *game) legacyTalkLines(index int, vars map[byte]string, width int) []string {
	lines, ok := g.talkLines(index, vars)
	if !ok {
		return nil
	}
	if width <= 0 {
		width = talkTextWidth
	}
	return textdraw.WrapLines(lines, width)
}

func (g *game) drawPortrait(screen *ebiten.Image, page, x, y, bank int) {
	if g == nil || g.lib == nil || page < 0 {
		return
	}
	img, err := g.lib.Portrait(page, bank)
	if err != nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(ebiten.NewImageFromImage(img), op)
}

func (g *game) drawIventScene(screen *ebiten.Image, page int) {
	if g == nil || g.lib == nil {
		return
	}
	asset := -1
	for i, entry := range g.lib.Entries {
		if entry.Spec.Name == gfx.Ivent.Name {
			asset = i
			break
		}
	}
	if asset < 0 {
		return
	}
	season := 0
	if g.world != nil {
		season = int(g.world.Clock.Season())
	}
	img, err := g.lib.Render(asset, page, season)
	if err != nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(talkSceneX, talkSceneY)
	screen.DrawImage(ebiten.NewImageFromImage(img), op)
}

// drawLegacyTalkBox 對齊 sub_1075B 的 64×64 肖像、16 px TALK 行與
// sub_1895D 的五行框。頁面切割在呼叫端完成；這裡不把結構尾空行畫出來。
func (g *game) drawLegacyTalkBox(screen *ebiten.Image, x, y, w, h int,
	lines []string, portraitPage int) {
	g.chrome.Window(screen, x, y, w, h, chrome.Menu)
	bank := 0
	if g.world != nil {
		bank = int(g.world.Clock.Season())
	}
	g.drawPortrait(screen, portraitPage, x+talkPortraitX-talkBoxX,
		y+talkPortraitY-talkBoxY, bank)
	for i, line := range lines {
		if i >= messagePageRows {
			break
		}
		g.td.Draw(screen, line, x+talkTextX-talkBoxX,
			y+talkTextY-talkBoxY+i*talkLinePitch,
			chrome.Paper)
	}
}

func (g *game) drawLegacyChoiceBox(screen *ebiten.Image, lines []string, selected int) {
	g.chrome.Window(screen, talkChoiceX, talkChoiceY,
		talkChoiceW, talkChoiceH, chrome.Menu)
	for i, line := range lines {
		if i >= talkChoiceRows {
			break
		}
		y := talkChoiceY + chrome.Tile + i*talkLinePitch
		if i == selected {
			vector.DrawFilledRect(screen, float32(talkChoiceX+chrome.Tile),
				float32(y-1), float32(talkChoiceW-2*chrome.Tile),
				float32(talkLinePitch), chrome.Select, false)
			vector.StrokeRect(screen, float32(talkChoiceX+chrome.Tile),
				float32(y-1), float32(talkChoiceW-2*chrome.Tile),
				float32(talkLinePitch), 1,
				color.RGBA{255, 223, 154, 255}, false)
			// ASCII `>` 與原版滑鼠／反白狀態一起保留，避免缺少中文字型
			// 時把游標畫成不可辨識的空方框。
			g.td.Draw(screen, ">", talkChoiceX+chrome.Tile, y, chrome.Paper)
			g.td.Draw(screen, line, talkChoiceX+chrome.Tile+16, y, chrome.Paper)
			continue
		}
		g.td.Draw(screen, line, talkChoiceX+chrome.Tile, y, chrome.Paper)
	}
}

// talkChoiceClick 回報滑鼠點到第幾列。`rows` 是這一次有幾個選項——
// 原版把它當 `al` 傳給 `sub_193E9`（`sub_13B5A` 給 5、`sub_13902` 給 3），
// **框的大小不跟著變，能選的範圍才變**（docs/spec/45 §2.2）。
func (g *game) talkChoiceClick(rows int) (int, bool) {
	x, y := ebiten.CursorPosition()
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		x < talkChoiceX+chrome.Tile || x >= talkChoiceX+talkChoiceW-chrome.Tile ||
		y < talkChoiceY+chrome.Tile || y >= talkChoiceY+talkChoiceH-chrome.Tile {
		return 0, false
	}
	if rows > talkChoiceRows {
		rows = talkChoiceRows
	}
	row := (y - (talkChoiceY + chrome.Tile)) / talkLinePitch
	if row < 0 || row >= rows {
		return 0, false
	}
	return row, true
}

// drawLegacyHint 是 remake 自己加的按鍵提示，原版沒有這一條。
//
// 框高要 4 塊：`chrome.Window` 的邊框是 8 px，兩塊高的框整個都是邊，
// 字會被切成看不出是什麼的東西。內容高 16 px ＝ 一個全形字。
//
// y 由呼叫端指定，因為**畫面底部 360–392 是事件視窗的位置**
// （`main.go` 的 `lastEvent`），兩者同時出現會疊在一起。
func (g *game) drawLegacyHint(screen *ebiten.Image, s string, y int) {
	if s == "" {
		return
	}
	w := g.td.Width(s) + 2*chrome.Tile
	if w > screenW-2*chrome.Tile {
		w = screenW - 2*chrome.Tile
	}
	x := (screenW - w) / 2 / chrome.Tile * chrome.Tile
	g.chrome.Window(screen, x, y, w, 4*chrome.Tile, chrome.Menu)
	g.td.Draw(screen, s, x+chrome.Tile, y+chrome.Tile,
		color.RGBA{170, 170, 180, 255})
}

func diplomacyTalkPromptIndex(c state.DiplomacyChoice, variant int) int {
	return diplomacyTalkBase(c.Kind) + talkVariant(variant)
}
