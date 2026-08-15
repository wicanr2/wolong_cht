package main

// 君主選擇視窗（docs/spec/27）。
//
// 原版是新遊戲流程的一頁（`sub_18E5A`）：一次顯示**一個**勢力的君主、軍師、
// 首都、武將數、據點數，右下兩顆是「自定」與「確定」。
// remake 的啟動殼層保留自己的流程，只把這一頁換成原版版面。

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

// 版面**全部出自原版**（docs/spec/27）：視窗矩形來自 `sub_1895D(cx=0C0Fh)`，
// 靜態層是顯示清單場景 8，數值座標由 `sub_18EA0` 的 VRAM 位移換算。
const (
	lordCardX, lordCardY = 160, 112
	lordCardW, lordCardH = 240, 192

	lordPortraitX, lordPortraitY   = 184, 128
	lordNameX, lordNameY           = 200, 216
	lordAdvPortraitX, lordAdvPortraitY = 312, 168
	lordAdvNameX, lordAdvNameY     = 328, 144

	lordLabelX                 = 184
	lordCapitalY               = 240
	lordGeneralsY              = 256
	lordCitiesY                = 272
	lordCapitalNameX           = 264
	lordCountX                 = 272
	lordCountDigits            = 3

	lordDividerX, lordDividerY = 247, 240
	lordDividerH               = 48

	// ⚠ 熱區編號與畫面上下顛倒：0x20 是下面的「確定」，0x21 是上面的「自定」。
	lordCustomX, lordCustomY = 328, 248
	lordOKX, lordOKY         = 328, 272
	lordButtonW, lordButtonH = 48, 16
	lordButtonTextX          = 336
)

// drawLordCard 畫一張原版版面的君主選擇卡。
func (g *game) drawLordCard(screen *ebiten.Image, p launcherPlayer, season int) {
	g.chrome.Window(screen, lordCardX, lordCardY, lordCardW, lordCardH, chrome.Menu)
	ink := g.paletteInk(strategyInkNormal, chrome.Paper)

	// 靜態層（顯示清單場景 8）。
	vector.DrawFilledRect(screen, 176, 128, 208, 104, color.Black, false)
	vector.DrawFilledRect(screen, 176, 240, 144, 48, color.Black, false)
	vector.DrawFilledRect(screen, lordDividerX, lordDividerY, 1, lordDividerH, ink, false)
	for _, b := range []struct {
		x, y  int
		label string
	}{{lordCustomX, lordCustomY, "自定"}, {lordOKX, lordOKY, "確定"}} {
		// 兩顆是**按鈕**：底色 7 ＋ 一圈 9／6（docs/re/48 §2.1）。
		g.dlButton(screen, b.x, b.y, lordButtonW, lordButtonH)
		g.td.Draw(screen, b.label, lordButtonTextX, b.y, ink)
	}
	g.td.Draw(screen, "君主", lordLabelX, 200, ink)
	g.td.Draw(screen, "軍師", 312, 128, ink)
	g.td.Draw(screen, "首　都", lordLabelX, lordCapitalY, ink)
	g.td.Draw(screen, "武將數", lordLabelX, lordGeneralsY, ink)
	g.td.Draw(screen, "據點數", lordLabelX, lordCitiesY, ink)

	// 動態層。
	if img, err := g.lib.Portrait(p.Portrait, season); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(lordPortraitX, lordPortraitY)
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	g.td.Draw(screen, p.Lord, lordNameX, lordNameY, ink)
	// 軍師是 0x7F（無）時**連頭像都不畫**，原版就是這樣。
	if p.HasAdvisor {
		if img, err := g.lib.Portrait(p.AdvisorPortrait, season); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(lordAdvPortraitX, lordAdvPortraitY)
			screen.DrawImage(ebiten.NewImageFromImage(img), op)
		}
		g.td.Draw(screen, p.Advisor, lordAdvNameX, lordAdvNameY, ink)
	}
	g.td.Draw(screen, p.Capital, lordCapitalNameX, lordCapitalY, ink)
	g.td.Draw(screen, strategyHUDNumber(p.Generals, lordCountDigits),
		lordCountX, lordGeneralsY, ink)
	g.td.Draw(screen, strategyHUDNumber(p.Cities, lordCountDigits),
		lordCountX, lordCitiesY, ink)
}
