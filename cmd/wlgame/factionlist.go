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

// 新遊戲那幾層的背景鏡頭：`sub_11A6E` 在 `sub_1D615` 前設的兩個立即值
// （`dx=0AAh`／`bx=62h`），單位是**格**（docs/spec/79 §1.1.1）。
// 調色盤是 `sub_10241(al=0)` 的第 0 組。
const (
	launcherCamX, launcherCamY = 170, 98
	launcherSeason             = 0
)

// 標題與分隔線是原版字串照抄（`cs:7AC6`／`cs:7AEB`）。
// 空槽那一列直接照印分隔線——原版就是同一份資料兩用。
const (
	factionListTitle = "勢力名　軍師名　　武將　據點　　首都"
	factionListDash  = "－－－　－－－　　 --　　-- 　　－－－"
	// factionColDash 是單一名字欄的分隔線（三個全形破折號）。
	factionColDash = "－－－"
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
	// 清單本體是**米色底**（色 9）配黑字（色 0），與遊戲中那七張一覽表
	// 同一套（docs/spec/38 §1.1）。先前這裡填 color.Black ＋ 淺灰字，
	// 與原版整片相反（docs/spec/107 §2）。`chrome` 在殼層之前就以第 0 組
	// 載好了，所以這兩個色票在殼層也是對的。
	vector.DrawFilledRect(screen, float32(factionListWinX), float32(factionListWinY),
		float32(factionListWinW), float32(factionListWinH), chrome.Sheet, false)
	// 標題列是**黑底白字**，與遊戲中那七張一覽表一樣（cmd/wlgame/main.go
	// 的 drawList）。實機量到的就是這兩塊：標題列 5,301 px 黑 ＋ 843 px 白。
	vector.DrawFilledRect(screen, float32(factionListWinX), float32(factionListWinY),
		float32(factionListWinW), float32(factionListRowH),
		color.RGBA{0, 0, 0, 255}, false)

	// 標題列是黑底白字，所以標題那一行仍用亮色。
	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	// 列上的字落在米色底上，用色 0。
	rowInk := chrome.Ink
	dim := chrome.Ink
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
		// 反白列的字色**不換**——原版的一覽表整頁都用色 0，
		// 選中那一列只是底下多鋪一條綠（cmd/wlgame/main.go 的 drawList
		// 走的是同一套）。
		col := rowInk
		p := l.players[n]
		x := factionListWinX
		g.td.Draw(screen, p.Lord, x+factionColLordX, y, col)
		// 沒有軍師（勢力 +0x02 ＝ 0x7F）那一格印分隔線，不是留白——
		// 原版實機的劉表／劉度／韓玄／袁術四列都是「－－－」。
		advisor := factionColDash
		if p.HasAdvisor {
			advisor = p.Advisor
		}
		g.td.Draw(screen, advisor, x+factionColAdvisorX, y, col)
		// ⚠ 兩個數字欄的 X 是**三位數欄位的左緣**，不是右緣
		// （docs/spec/79 §1.4 的絕對 X 304／352）；先前多減了一個欄寬，
		// 整欄往左偏 24 px。字模也要用**原版的 8×16 數字**
		// （docs/spec/38 §1.5），文字字型的數字形狀與原版不一樣。
		g.drawOriginalNumber(screen, p.Generals,
			x+factionColGeneralsX, y, factionColDigits, col)
		g.drawOriginalNumber(screen, p.Cities,
			x+factionColCitiesX, y, factionColDigits, col)
		g.td.Draw(screen, p.Capital, x+factionColCapitalX, y, col)
	}

	// 捲軸與戰略層那七張一覽表共用同一份畫法（docs/spec/38 §1.6）：
	// 純黑槽 ＋ 3D 綠鈕 ＋ 比例式滑塊。先前這裡自己描了三個框。
	g.drawScrollbarAt(screen,
		factionListTitleScrollRect(), factionListScrollUpRect(),
		factionListScrollTrackRect(), factionListScrollDownRect(),
		top, len(l.players), factionListRowsPerPage)
}

// factionListTitleScrollRect 是標題列左邊那一格（捲軸欄的頂端，純黑）。
func factionListTitleScrollRect() image.Rectangle {
	return image.Rect(factionListWinX, factionListWinY,
		factionListWinX+factionListScrollW, factionListWinY+factionListRowH)
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

// updateFactionListPointer 是清單的滑鼠：滾輪捲、點一列選它並進君主卡、
// 點上下箭頭捲一列。
// 回傳 handled=true 表示這一幀的**點擊**處理完了，不要再走鍵盤那一段。
//
// ⚠ **滑鼠移動不改變選中的勢力**——那是上一輪修掉的那個 bug 的規則
// （docs/spec/27 §2.1），這一頁照同一條走。滾輪捲的是**視野**，
// 與 ▲▼ 同一個語意，也不動選取列（原版的 `top` 與選取列本來就是
// 兩個獨立狀態）。滾輪本身是 remake 加的，原版只有 ▲▼。
func (g *game) updateFactionListPointer() (bool, error) {
	l := g.launcher
	if _, dy := ebiten.Wheel(); dy != 0 {
		step := 1
		if dy > 0 {
			step = -1
		}
		l.factionTop = clampFactionTop(l.factionListTop()+step, len(l.players))
	}
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false, nil
	}
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
