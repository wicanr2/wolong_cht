package main

// 新遊戲的勢力清單（docs/spec/79）。
//
// 原版的新遊戲是四層，每一層右鍵退回上一層：
// ＹＥＳ／ＮＯ → 劇本 → **這一張清單** → 君主卡。
// ⭐ 君主卡上沒有換勢力的熱區（docs/spec/27 §2.1），換勢力就是退回這裡。

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 幾何全部出自 `sub_17B3C`／兩個 callback（docs/re/73 §2）。
//
// ⚠ **尺寸與戰略層那七個一覽表相同，左上角不同**：那七個是 (24, 88)，
// 這一個是 (136, 104)。抄一覽表的版面到別的畫面時，抄的是尺寸與列距，
// 不是左上角（docs/re/26 §1）。
const (
	factionListWinX, factionListWinY = 136, 104
	factionListWinW, factionListWinH = 384, 176

	// 列首在視窗左緣 +24（callback 的 `dx = word_181AE + 18h`）。
	factionListRowDX = 24
	// 第一列在視窗上緣 +16——標題佔掉那一列（`bx = word_181B0 + 10h`）。
	factionListRowDY = 16
	factionListRowH  = 16
	// 一頁十列：畫列的 callback 就是十次迴圈。
	factionListRowsPerPage = 10
	// 左邊 16 px 是捲軸，與一覽表同一套（docs/spec/38 §1.1.1）。
	factionListScrollW = 16
)

// 五欄的 X，**相對視窗左緣**（欄位表 `cs:7B12`，docs/spec/79 §1.4）。
const (
	factionColLordX     = 24
	factionColAdvisorX  = 88
	factionColGeneralsX = 168
	factionColCitiesX   = 216
	factionColCapitalX  = 280
	factionColDigits    = 3
)

// 標題與分隔線是原版字串照抄（`cs:7AC6`／`cs:7AEB`）。
// 空槽那一列直接照印分隔線——原版就是同一份資料兩用。
const (
	factionListTitle = "勢力名　軍師名　　武將　據點　　首都"
	factionListDash  = "－－－　－－－　　 --　　-- 　　－－－"
)

func factionListRowFirstY() int { return factionListWinY + factionListRowDY }
func factionListRowX() int      { return factionListWinX + factionListRowDX }
func factionListBodyX() int     { return factionListWinX + factionListScrollW }
func factionListBodyW() int     { return factionListWinW - factionListScrollW }

// factionListTitleRect 是欄位標題那一列。
func factionListTitleRect() image.Rectangle {
	return image.Rect(factionListWinX, factionListWinY,
		factionListWinX+factionListWinW, factionListWinY+factionListRowH)
}

// factionListBodyRect 是清單本體（捲軸右邊、標題下面）。
func factionListBodyRect() image.Rectangle {
	y := factionListRowFirstY()
	return image.Rect(factionListBodyX(), y,
		factionListBodyX()+factionListBodyW(),
		y+factionListRowsPerPage*factionListRowH)
}

func factionListScrollUpRect() image.Rectangle {
	y := factionListRowFirstY()
	return image.Rect(factionListWinX, y,
		factionListWinX+factionListScrollW, y+factionListRowH)
}

func factionListScrollDownRect() image.Rectangle {
	y := factionListWinY + factionListWinH - factionListRowH
	return image.Rect(factionListWinX, y,
		factionListWinX+factionListScrollW, y+factionListRowH)
}

func factionListScrollTrackRect() image.Rectangle {
	return image.Rect(factionListWinX, factionListRowFirstY()+factionListRowH,
		factionListWinX+factionListScrollW,
		factionListWinY+factionListWinH-factionListRowH)
}

// factionListRowRect 是第 i 列（頁內編號 0–9）的矩形。
func factionListRowRect(i int) image.Rectangle {
	y := factionListRowFirstY() + i*factionListRowH
	return image.Rect(factionListBodyX(), y,
		factionListBodyX()+factionListBodyW(), y+factionListRowH)
}

// factionListPageRow 把螢幕座標換成頁內列號。標題列與捲軸都不算。
func factionListPageRow(x, y int) (int, bool) {
	if !image.Pt(x, y).In(factionListBodyRect()) {
		return 0, false
	}
	return (y - factionListRowFirstY()) / factionListRowH, true
}

// clampFactionTop 把捲動位置夾在 [0, 筆數 − 一頁]。
func clampFactionTop(top, total int) int {
	span := total - factionListRowsPerPage
	if span < 0 {
		span = 0
	}
	return clamp(top, 0, span)
}

// factionListTop 是目前的捲動位置。
//
// ⚠ **它不會為了讓游標露出來而自己動。** 原版的 `top` 與選取列是兩個
// 獨立的狀態（`sub_1820E` 的 `word_181BA` 與捲動 handler 各管各的），
// 點 ▲▼ 捲動時選取列可以捲出畫面。要跟著游標走的是**移動游標**那一邊，
// 走 `scrollFactionListToCursor`。
func (l *launcherModel) factionListTop() int {
	return clampFactionTop(l.factionTop, len(l.players))
}

// scrollFactionListToCursor 把捲動位置調到讓選中的那一列露出來。
// 只有「游標動了」才呼叫它——捲軸捲動不呼叫。
func (l *launcherModel) scrollFactionListToCursor() {
	top := l.factionListTop()
	if l.cursor < top {
		top = l.cursor
	}
	if l.cursor >= top+factionListRowsPerPage {
		top = l.cursor - factionListRowsPerPage + 1
	}
	l.factionTop = clampFactionTop(top, len(l.players))
}

// factionListSelectAt 把點擊換成勢力編號。回傳 false 表示點在空列上——
// 原版空列印的是分隔線，點下去沒有東西可選。
func (l *launcherModel) factionListSelectAt(x, y int) (int, bool) {
	row, ok := factionListPageRow(x, y)
	if !ok {
		return 0, false
	}
	n := l.factionListTop() + row
	if n < 0 || n >= len(l.players) {
		return 0, false
	}
	return n, true
}

// drawFactionList 畫清單。三塊：框、標題列、十列。
func (g *game) drawFactionList(screen *ebiten.Image) {
	l := g.launcher
	if l == nil {
		return
	}
	// 框畫在內容區外面：往左上各退 8 px、寬高各多 16 px（docs/spec/38 §1.1）。
	g.chrome.Window(screen, factionListWinX-chrome.Tile, factionListWinY-chrome.Tile,
		factionListWinW+2*chrome.Tile, factionListWinH+2*chrome.Tile, chrome.Menu)
	vector.DrawFilledRect(screen, float32(factionListWinX), float32(factionListWinY),
		float32(factionListWinW), float32(factionListWinH), color.Black, false)

	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	dim := g.paletteInk(strategyInkDim, color.RGBA{200, 200, 210, 255})
	g.td.Draw(screen, factionListTitle, factionListRowX(), factionListWinY, ink)

	top := l.factionListTop()
	for i := 0; i < factionListRowsPerPage; i++ {
		n := top + i
		y := factionListRowFirstY() + i*factionListRowH
		if n >= len(l.players) {
			// 空列照印分隔線——原版拿的就是同一份字串（docs/re/73 §3）。
			g.td.Draw(screen, factionListDash, factionListRowX(), y, dim)
			continue
		}
		if n == l.cursor {
			r := factionListRowRect(i)
			vector.DrawFilledRect(screen, float32(r.Min.X), float32(r.Min.Y),
				float32(r.Dx()), float32(r.Dy()), chrome.Select, false)
		}
		col := dim
		if n == l.cursor {
			col = ink
		}
		p := l.players[n]
		x := factionListWinX
		g.td.Draw(screen, p.Lord, x+factionColLordX, y, col)
		if p.HasAdvisor {
			g.td.Draw(screen, p.Advisor, x+factionColAdvisorX, y, col)
		}
		g.td.Draw(screen, strategyHUDNumber(p.Generals, factionColDigits),
			x+factionColGeneralsX-factionColDigits*textdraw.HalfW, y, col)
		g.td.Draw(screen, strategyHUDNumber(p.Cities, factionColDigits),
			x+factionColCitiesX-factionColDigits*textdraw.HalfW, y, col)
		g.td.Draw(screen, p.Capital, x+factionColCapitalX, y, col)
	}

	g.drawFactionListScrollbar(screen, top, len(l.players), ink)
}

func (g *game) drawFactionListScrollbar(screen *ebiten.Image, top, total int, ink color.RGBA) {
	for _, r := range []image.Rectangle{
		factionListScrollUpRect(), factionListScrollTrackRect(), factionListScrollDownRect(),
	} {
		vector.StrokeRect(screen, float32(r.Min.X), float32(r.Min.Y),
			float32(r.Dx()), float32(r.Dy()), 1, ink, false)
	}
	// 三角形用三條水平線疊出來，與一覽表同一畫法。
	tri := func(r image.Rectangle, up bool) {
		cx := float32(r.Min.X+r.Dx()/2) - 0.5
		for i := 0; i < 5; i++ {
			w := float32(2*i + 1)
			y := float32(r.Min.Y + 5 + i)
			if !up {
				y = float32(r.Max.Y - 6 - i)
			}
			vector.DrawFilledRect(screen, cx-w/2, y, w, 1, ink, false)
		}
	}
	tri(factionListScrollUpRect(), true)
	tri(factionListScrollDownRect(), false)

	// ⚠ 滑塊按比例畫，**原版那一支沒讀**（docs/spec/38 §4）。
	track := factionListScrollTrackRect()
	span := total - factionListRowsPerPage
	if span <= 0 {
		return
	}
	h := track.Dy() * factionListRowsPerPage / total
	if h < factionListRowH {
		h = factionListRowH
	}
	y := track.Min.Y + (track.Dy()-h)*clamp(top, 0, span)/span
	vector.DrawFilledRect(screen, float32(track.Min.X+3), float32(y),
		float32(track.Dx()-6), float32(h), ink, false)
}

// updateFactionListPointer 是清單的滑鼠：點一列選它並進君主卡，
// 點上下箭頭捲一列。回傳 handled=true 表示這一幀的滑鼠處理完了。
func (g *game) updateFactionListPointer() (bool, error) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false, nil
	}
	l := g.launcher
	x, y := ebiten.CursorPosition()
	p := image.Pt(x, y)
	switch {
	case p.In(factionListScrollUpRect()):
		l.factionTop = clampFactionTop(l.factionListTop()-1, len(l.players))
		return true, nil
	case p.In(factionListScrollDownRect()):
		l.factionTop = clampFactionTop(l.factionListTop()+1, len(l.players))
		return true, nil
	}
	if n, ok := l.factionListSelectAt(x, y); ok {
		l.cursor = n
		return true, g.applyLauncherResult(l.apply(launcherConfirm))
	}
	// 點在別的地方什麼都不做——與君主卡同一條規則。
	return true, nil
}
