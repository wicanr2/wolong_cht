package main

// 軍團情報視窗（docs/spec/24）。
//
// 原版在大地圖上點軍團或走指令列「軍團」進來（`sub_17F90`）：
// **自己的軍團**看完接著進指令流程，別人的只是看。remake 這一輪只做顯示。
//
// 視窗矩形與自勢力情報那一格相同——開著的時候蓋掉它，原版就是這樣。

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/army"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

type corpsInfoState struct {
	active bool
	corps  int
}

// 版面**全部出自原版**（docs/spec/24）：視窗矩形來自 `sub_1895D(cx=0D0Dh)`，
// 靜態層是顯示清單場景 4，數值座標由 `sub_1807B`／`sub_1812A` 的
// VRAM 位移換算。
const (
	corpsWinX, corpsWinY = 432, 192
	corpsWinW, corpsWinH = 208, 208

	corpsPortraitX, corpsPortraitY = 440, 200

	corpsHeadLabelX = 512
	corpsHeadValueX = 576
	corpsLeaderY    = 208
	corpsLordY      = 224
	corpsCapitalY   = 240
	corpsDividerX   = 560
	corpsDividerH   = 48

	corpsTotalLabelX = 464
	corpsTotalY      = 272
	corpsTotalX      = 536
	corpsSlashX      = 568
	corpsMoraleX     = 584
	corpsTotalDigits = 4
	corpsMoraleDigits = 3

	// 六個槽：標籤 → 兵種圖示 → 兵力。
	corpsSlotLabelX = 464
	corpsSlotIconX  = 520
	corpsSlotValueX = 576
	corpsSlotY      = 288
	corpsSlotStep   = 16
	corpsSlotDigits = 4

	// remake 差異：提示自己一個框，接在原版視窗左邊
	// （視窗右緣就是畫面右緣，下緣就是畫面下緣）。
	corpsHintW = 176
	corpsHintX = corpsWinX - corpsHintW
	corpsHintH = 32
	corpsHintY = corpsWinY
)

func (g *game) openCorpsInfo(corps int) {
	if corps < 0 || corps >= len(g.world.Corps) || !g.world.Corps[corps].Alive {
		return
	}
	g.corpsInfo = corpsInfoState{active: true, corps: corps}
}

func (g *game) updateCorpsInfo() {
	if pressed(ebiten.KeyEscape) || pressed(ebiten.KeyEnter) {
		g.corpsInfo.active = false
	}
}

func (g *game) drawCorpsInfo(screen *ebiten.Image) {
	if !g.corpsInfo.active {
		return
	}
	c := g.world.Corps[g.corpsInfo.corps]
	g.chrome.Window(screen, corpsWinX, corpsWinY, corpsWinW, corpsWinH, chrome.Menu)

	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	labelInk := g.paletteInk(strategyInkDim, color.RGBA{255, 223, 154, 255})
	season := int(g.world.Clock.Season())

	// 靜態層（顯示清單場景 4）。
	vector.DrawFilledRect(screen, 456, corpsTotalY, 160, 112, color.Black, false)
	vector.DrawFilledRect(screen, corpsDividerX, corpsLeaderY, 1, corpsDividerH, ink, false)
	g.td.Draw(screen, "將軍", corpsHeadLabelX, corpsLeaderY, ink)
	g.td.Draw(screen, "君主", corpsHeadLabelX, corpsLordY, ink)
	g.td.Draw(screen, "首都", corpsHeadLabelX, corpsCapitalY, ink)
	g.td.Draw(screen, "總兵力", corpsTotalLabelX, corpsTotalY, ink)
	g.td.Draw(screen, "／", corpsSlashX, corpsTotalY, ink)
	for k := 0; k < army.Positions; k++ {
		g.td.Draw(screen, formSlotLabels[k],
			corpsSlotLabelX, corpsSlotY+k*corpsSlotStep, labelInk)
	}

	// 上半：頭像與三個名字。軍團編號 ＝ 主將的武將編號。
	leader := g.world.Generals[g.world.Leader(g.corpsInfo.corps)]
	if img, err := g.lib.Portrait(leader.Portrait, season); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(corpsPortraitX), float64(corpsPortraitY))
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	g.td.Draw(screen, big5(leader.Name), corpsHeadValueX, corpsLeaderY, ink)
	if c.Faction >= 0 && c.Faction < len(g.world.Factions) {
		f := g.world.Factions[c.Faction]
		g.td.Draw(screen, big5(g.world.LordName(c.Faction)), corpsHeadValueX, corpsLordY, ink)
		if f.Capital >= 0 && f.Capital < len(g.world.Cities) {
			g.td.Draw(screen, big5(g.world.Cities[f.Capital].Name),
				corpsHeadValueX, corpsCapitalY, ink)
		}
	}

	// 「總兵力 6000／200」是一列：值、斜線、士氣。
	g.td.Draw(screen, strategyHUDNumber(c.Men*strategyReserveMenPerPoint,
		corpsTotalDigits), corpsTotalX, corpsTotalY, ink)
	g.td.Draw(screen, strategyHUDNumber(c.Morale, corpsMoraleDigits),
		corpsMoraleX, corpsTotalY, ink)

	// 六個槽。**空槽照原版畫天秤**——原版的圖庫基底是紅色那一組，
	// 而兵種 4 算出來剛好越界到綠色組的第一張（docs/re/51 §4）。
	for k := 0; k < army.Positions; k++ {
		y := corpsSlotY + k*corpsSlotStep
		kind := int(c.Units[k].Kind) + 1
		if img := g.corpsSlotIcon(kind, season); img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(corpsSlotIconX), float64(y))
			screen.DrawImage(ebiten.NewImageFromImage(img), op)
		}
		g.td.Draw(screen, strategyHUDNumber(c.Units[k].Men*strategyReserveMenPerPoint,
			corpsSlotDigits), corpsSlotValueX, y, labelInk)
	}

	// ↓ remake 差異：原版按右鍵關掉，沒有這行字。
	g.chrome.Window(screen, corpsHintX, corpsHintY, corpsHintW, corpsHintH, chrome.Menu)
	g.td.Draw(screen, "ESC 關閉", corpsHintX+8,
		corpsHintY+(corpsHintH-textdraw.GlyphH)/2, labelInk)
}

// corpsSlotIcon 取軍團情報那一欄的兵種圖示：**紅色組**，兵種 4 越界到
// 綠色組第一張（原版行為，docs/spec/24 §1.3）。
func (g *game) corpsSlotIcon(kind, season int) *image.RGBA {
	if kind < 1 || kind > 4 {
		return nil
	}
	if kind == 4 {
		img, err := g.lib.DOSVResourceIcon(0, true, season)
		if err != nil {
			return nil
		}
		return img
	}
	img, err := g.lib.DOSVResourceIcon(kind, false, season)
	if err != nil {
		return nil
	}
	return img
}
