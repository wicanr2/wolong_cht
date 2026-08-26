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

	// 選單框的左上角：外交三選一、撥款、說服五選一共用 `sub_13B7E`，
	// 那一組座標是寫死的（docs/spec/45 §2）：
	//
	//	sub_193E9(dl=5, dh=0Bh)    粗格 ⇒ (5×16, 11×16) ＝ (80, 176)
	//	sub_19796(dx=50h, bx=0B0h) 像素 ⇒ 同一個點
	//
	// **大小不是常數**——`sub_19479` 由內容算（見 legacyChoiceRect）。
	talkChoiceX = 80
	talkChoiceY = 176

	// adviseMenuX／Y 是進言那五項的選單（`sub_16224` 的 `dx = 400h`，
	// 粗格 (0, 4)）。
	adviseMenuX = 0
	adviseMenuY = 64
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
	// ⭐ **TALK 框的底是黑的，不是選單的藍底龍紋**（docs/spec/88 §1）。
	// 實機量過：框內 6,298 px 全是 (0,0,0)。`chrome.Menu` 會鋪龍紋
	//（`fillInterior` 只在 fill == Menu 時鋪），所以這裡要傳 `Blank`。
	// ⚠ 系統選單確實是藍底龍紋（playtest/39 逐像素 PASS），別一起改。
	g.chrome.Window(screen, x, y, w, h, chrome.Blank)
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

// legacyChoiceRect 是選單框的大小。**由內容算，不是常數**
// （`sub_19479`，docs/spec/45 §2.2）：
//
//	cl ＝ 第一列的全形字數、ch ＝ 列數（呼叫端傳給 sub_193E9 的 al）
//	inc cl / inc ch  ⇒  寬 (字數+1)×16、高 (列數+1)×16
//
// 那個 +1 就是四邊各 8 px 的框邊。
//
// ⭐ **原版只看第一列**——它靠選單文字全部用全形空白補到等寬
// （#102 五列都是 6 個字，#77 也是）來保證第一列就是最寬的。
//
// ⚠ **remake 取所有列的最大值**，因為 `text.Decode` 會 `TrimRight`
// 掉行尾的全形空白（那對武將名、據點名是對的），拿不回原本的 padding。
// 四則選單（#77／#102／#363／#376）兩種算法的結果**相同**，
// 差別只在原版會切掉「比第一列長的後面幾列」，這裡不會。
func legacyChoiceRect(x, y int, lines []string) (int, int, int, int) {
	cells := 0
	for _, l := range lines {
		if n := textdraw.StringWidth(l) / talkLinePitch; n > cells {
			cells = n
		}
	}
	if cells < 1 {
		cells = 1
	}
	return x, y, (cells + 1) * talkLinePitch, (len(lines) + 1) * talkLinePitch
}

func (g *game) drawLegacyChoiceBox(screen *ebiten.Image, x, y int,
	lines []string, selected int) {

	bx, by, w, h := legacyChoiceRect(x, y, lines)
	// 選項框與 TALK 框同一族，底一樣是黑的（docs/spec/88 §1）。
	g.chrome.Window(screen, bx, by, w, h, chrome.Blank)
	for i, line := range lines {
		ly := by + chrome.Tile + i*talkLinePitch
		if i == selected {
			vector.DrawFilledRect(screen, float32(bx+chrome.Tile),
				float32(ly-1), float32(w-2*chrome.Tile),
				float32(talkLinePitch), chrome.Select, false)
			vector.StrokeRect(screen, float32(bx+chrome.Tile),
				float32(ly-1), float32(w-2*chrome.Tile),
				float32(talkLinePitch), 1,
				color.RGBA{255, 223, 154, 255}, false)
		}
		g.td.Draw(screen, line, bx+chrome.Tile, ly, chrome.Paper)
	}
}

// talkChoiceClick 回報滑鼠點到第幾列。命中範圍與畫出來的框同一個算式。
func (g *game) talkChoiceClick(x0, y0 int, lines []string) (int, bool) {
	if len(lines) == 0 {
		return 0, false
	}
	bx, by, w, h := legacyChoiceRect(x0, y0, lines)
	x, y := ebiten.CursorPosition()
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		x < bx+chrome.Tile || x >= bx+w-chrome.Tile ||
		y < by+chrome.Tile || y >= by+h-chrome.Tile {
		return 0, false
	}
	row := (y - (by + chrome.Tile)) / talkLinePitch
	if row < 0 || row >= len(lines) {
		return 0, false
	}
	return row, true
}

func (g *game) drawLegacyHint(screen *ebiten.Image, s string, y int) {
	if s == "" {
		return
	}
	w := g.td.Width(s) + 2*chrome.Tile
	if w > screenW-2*chrome.Tile {
		w = screenW - 2*chrome.Tile
	}
	x := (screenW - w) / 2 / chrome.Tile * chrome.Tile
	// 按鍵提示是 remake 加的，但要跟旁邊那兩個框同一個樣子。
	g.chrome.Window(screen, x, y, w, 4*chrome.Tile, chrome.Blank)
	g.td.Draw(screen, s, x+chrome.Tile, y+chrome.Tile,
		color.RGBA{170, 170, 180, 255})
}

func diplomacyTalkPromptIndex(c state.DiplomacyChoice, variant int) int {
	return diplomacyTalkBase(c.Kind) + talkVariant(variant)
}
