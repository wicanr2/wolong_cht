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

// 事件 2／3 的畫面位置以 DOS/V `sub_13D09`／`sub_13D68` 的
// IVENTGRF 288×176 blit 與 640×400 畫布為主；既有 PC-98 fixture 只作
// 歷史交叉驗證，不是本輪畫面 oracle。
const (
	talkBoxX = 24
	talkBoxY = 80
	talkBoxW = 232
	talkBoxH = 80

	talkPortraitX = 32
	talkPortraitY = 88
	talkTextX     = 80
	talkTextY     = 96
	talkTextWidth = 160
	talkLinePitch = 16 // sub_18810／sub_1895D 的 CL=10h

	talkSceneX = 64
	talkSceneY = 144

	talkChoiceX = 80
	talkChoiceY = 176
	talkChoiceW = 128
	talkChoiceH = 64
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
		if i >= 3 {
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

func (g *game) talkChoiceClick() (int, bool) {
	x, y := ebiten.CursorPosition()
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		x < talkChoiceX+chrome.Tile || x >= talkChoiceX+talkChoiceW-chrome.Tile ||
		y < talkChoiceY+chrome.Tile || y >= talkChoiceY+talkChoiceH-chrome.Tile {
		return 0, false
	}
	row := (y - (talkChoiceY + chrome.Tile)) / talkLinePitch
	if row < 0 || row >= 3 {
		return 0, false
	}
	return row, true
}

func (g *game) drawLegacyHint(screen *ebiten.Image, s string) {
	// Hint 放在場景下方的細框，不覆蓋 IVENTGRF 或原版 TALK 內容。
	if s == "" {
		return
	}
	w := g.td.Width(s) + 2*chrome.Tile
	if w > screenW-2*chrome.Tile {
		w = screenW - 2*chrome.Tile
	}
	x := (screenW - w) / 2 / chrome.Tile * chrome.Tile
	y := screenH - 2*chrome.Tile
	g.chrome.Window(screen, x, y, w, 2*chrome.Tile, chrome.Menu)
	g.td.Draw(screen, s, x+chrome.Tile, y+1, color.RGBA{170, 170, 180, 255})
}

func diplomacyTalkPromptIndex(c state.DiplomacyChoice, variant int) int {
	return diplomacyTalkBase(c.Kind) + talkVariant(variant)
}
