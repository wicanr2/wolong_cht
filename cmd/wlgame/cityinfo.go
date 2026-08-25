package main

// 據點情報視窗（docs/spec/23）。
//
// 原版是在大地圖上點一個據點就跳出來（`sub_11E46`），純顯示、沒有熱區，
// 按右鍵關掉。remake 走一覽表選取——那是既有的操作差異，版面照原版。

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// cityInfoState 是據點情報視窗的狀態。
type cityInfoState struct {
	active bool
	city   int
}

// 版面**全部出自原版**（docs/spec/23）：視窗矩形來自 `sub_1895D(cx=810h)`，
// 靜態層是顯示清單場景 3，數值座標由 `sub_17E4A` 的 VRAM 位移換算。
const (
	cityWinX, cityWinY = 0, 272
	cityWinW, cityWinH = 256, 128

	cityViewX, cityViewY = 16, 288 // KYOGRF 96×96
	cityViewSize         = 96

	cityNameX, cityNameY = 128, 288
	cityKindX            = 192
	cityLordX, cityLordY = 192, 304

	cityLabelX = 128
	cityRowY   = 320
	cityRowStep = 16
	cityRows    = 4

	cityGarrisonX, cityProductionX = 208, 192
	cityGrowthX, cityPreventionX   = 208, 216

	cityGarrisonDigits   = 4
	cityProductionDigits = 6
	cityGrowthDigits     = 4
	cityPreventionDigits = 3

	// remake 差異：操作提示自己一個框，接在原版視窗右邊
	// （下面已經沒有空間了——視窗底就是畫面底）。
	cityHintX = cityWinX + cityWinW
	cityHintW = 240
	cityHintH = 32
)

// cityKindLabels 是 `cs:7DF5` 的六個詞。索引就是 City.Kind，
// **但首都覆寫成第 5 個**（原版比對勢力 +0x03，見 docs/re/50 §2.1）。
var cityKindLabels = [...]string{"大都市", "中都市", "小都市", "關卡", "戰場", "首都"}

// cityKindLabel 回傳該據點在情報視窗上印的類型。
func cityKindLabel(kind int, capital bool) string {
	if capital {
		return cityKindLabels[5]
	}
	if kind < 0 || kind >= len(cityKindLabels)-1 {
		return ""
	}
	return cityKindLabels[kind]
}

// openCityInfo 開據點情報視窗。
func (g *game) openCityInfo(city int) {
	if city < 0 || city >= len(g.world.Cities) {
		return
	}
	g.cityInfo = cityInfoState{active: true, city: city}
}

func (g *game) updateCityInfo() {
	if g.cancelled() || pressed(ebiten.KeyEnter) {
		g.cityInfo.active = false
	}
}

func (g *game) drawCityInfo(screen *ebiten.Image) {
	if !g.cityInfo.active {
		return
	}
	c := g.world.Cities[g.cityInfo.city]
	g.chrome.Window(screen, cityWinX, cityWinY, cityWinW, cityWinH, chrome.Menu)

	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	labelInk := g.paletteInk(strategyInkDim, color.RGBA{255, 223, 154, 255})
	warnInk := g.paletteInk(strategyInkGauge, color.RGBA{210, 48, 40, 255})
	season := int(g.world.Clock.Season())

	// 靜態層（顯示清單場景 3）。
	vector.DrawFilledRect(screen, cityLabelX, cityRowY, 112, 64, color.Black, false)
	for i, label := range []string{"城兵數", "生產力", "上昇值", "防災值"} {
		g.td.Draw(screen, label, cityLabelX, cityRowY+i*cityRowStep, labelInk)
	}

	// 左半那張 96×96 景觀圖，張號是據點記錄 +0x16 的高 4 位。
	// **解不出來就留白**——換一張頂替會讓「圖對不對」永遠測不出來。
	if img, err := g.lib.Location(c.KindHigh, season); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cityViewX), float64(cityViewY))
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}

	capital := c.Owner >= 0 && c.Owner < len(g.world.Factions) &&
		g.world.Factions[c.Owner].Capital == g.cityInfo.city
	g.td.Draw(screen, big5(c.Name), cityNameX, cityNameY, ink)
	g.td.Draw(screen, cityKindLabel(c.Kind, capital), cityKindX, cityNameY, ink)

	lord := "中立"
	if c.Owner >= 0 && c.Owner < len(g.world.Factions) && g.world.Factions[c.Owner].Alive {
		lord = big5(g.world.LordName(c.Owner))
	}
	g.td.Draw(screen, lord, cityLordX, cityLordY, ink)

	// 四個值用 **8×16 原版數字字模**（實機 m0 逐格：倚天 ASCII 的「1」
	// 墨寬 6px、原版是 4px，spec/23 §4）。
	// 城兵數與預備兵同一個單位：存的是點，畫面上乘 10。
	g.drawOriginalNumber(screen, c.Garrison*strategyReserveMenPerPoint,
		cityGarrisonX, cityRowY, cityGarrisonDigits, labelInk)
	g.drawOriginalNumber(screen, c.Production,
		cityProductionX, cityRowY+cityRowStep, cityProductionDigits, labelInk)
	growthInk := labelInk
	if c.Growth < 0 {
		growthInk = warnInk
	}
	// Growth 已經是減過 100 的實際值，這裡不再減一次。
	g.drawOriginalNumber(screen, c.Growth,
		cityGrowthX, cityRowY+2*cityRowStep, cityGrowthDigits, growthInk)
	g.drawOriginalNumber(screen, c.Prevention,
		cityPreventionX, cityRowY+3*cityRowStep, cityPreventionDigits, labelInk)

	// ↓ remake 差異：原版按右鍵關掉，沒有這行字。
	g.chrome.Window(screen, cityHintX, cityWinY, cityHintW, cityHintH, chrome.Menu)
	g.td.Draw(screen, "ESC 關閉", cityHintX+8, cityWinY+(cityHintH-textdraw.GlyphH)/2, labelInk)
}
