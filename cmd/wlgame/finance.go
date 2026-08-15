package main

// 財政、據點、勢力——命令視窗剩下的三格。
//
// **財政的設定值延遲到「次月末」才生效**（說明書：ここにセットされた値は
// 来月末より使用されます，`docs/mechanics/15-realtime.md` §4）。
// 所以這裡改的是 `NextTaxRate`／`NextRecruitCap` 而不是現行值，
// 而畫面要**同時顯示兩欄**——不然玩家看不出「我改的東西還沒生效」。
// 原版的財政視窗也是「今月末」與「來月」兩欄對照。

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// financeState 是財政畫面的狀態。四列：稅率、騎馬、弓兵、步兵
// （原版視窗裡四個綠色按鈕，由上而下就是這個順序）。
type financeState struct {
	active bool
	row    int
}

// TaxMax 是稅率的上限。
//
// ⚠ 這個值**還沒從原版讀出來**，先用 100 當上界。
// 平衡點在 33.75%（`internal/state` 的 TestTaxTippingPoint 有推導），
// 所以上限只要遠高於它就不影響玩法判斷。標成 remake 的暫定值。
const TaxMax = 100

// recruitStep 是募兵數的調整刻度。原版用小算盤直接輸入數字，
// remake 先用固定刻度——**這是操作方式的差異，不是規則的差異**。
const recruitStep = 100

func (g *game) beginFinance() { g.finance = financeState{active: true} }

func (g *game) updateFinance() {
	f := &g.finance
	switch {
	case pressed(ebiten.KeyEscape):
		f.active = false
		return
	case pressed(ebiten.KeyArrowUp):
		f.row = (f.row + 3) % 4
	case pressed(ebiten.KeyArrowDown):
		f.row = (f.row + 1) % 4
	}
	delta := 0
	if pressed(ebiten.KeyArrowRight) || pressed(ebiten.KeyEqual) {
		delta = 1
	}
	if pressed(ebiten.KeyArrowLeft) || pressed(ebiten.KeyMinus) {
		delta = -1
	}
	if delta == 0 {
		return
	}
	if f.row == 0 {
		g.world.NextTaxRate = clamp(g.world.NextTaxRate+delta, 0, TaxMax)
		return
	}
	t := f.row - 1
	g.world.NextRecruitCap[t] = clamp(
		g.world.NextRecruitCap[t]+delta*recruitStep, 0, 60000)
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// 財政視窗的版面**全部出自原版**（docs/spec/14）：視窗矩形來自
// `sub_1895D(cx=0A15h)`，靜態層是顯示清單場景 1，數值座標由
// `sub_16846` 的 VRAM 位移換算（一列 80 byte）。
const (
	financeWinX, financeWinY = 16, 80
	financeWinW, financeWinH = 336, 160

	financeFundsLabelX, financeFundsLabelY = 32, 96
	financeIncomeLabelX                    = 200
	financeIncomeLabelY                    = 96
	financeExpenseLabelY                   = 112
	financeDividerX, financeDividerY       = 232, 96
	financeDividerH                        = 32

	financeThisLabelX, financeNextLabelX = 48, 208
	financeColLabelY                     = 136

	// 兩欄的左緣：今月底 32、次月 192。列距 16，四列。
	financeThisColX, financeNextColX = 32, 192
	financeRowY                      = 160
	financeRowStep                   = 16
	financeRows                      = 4

	financeIconThisX, financeIconNextX = 96, 256
	financePercentThisX                = 160
	financePercentNextX                = 320

	// 數值的右端由「左緣 ＋ 位數 × 8」決定，座標記的是左緣。
	financeFundsValueX, financeFundsValueY = 104, 112
	financeIncomeValueX                    = 288
	financeValueThisX, financeValueNextX   = 120, 280

	financeFundsDigits  = 7
	financeAmountDigits = 6
	financeRowDigits    = 5

	// 熱區與綠色圖示欄逐格重合：可點的只有「次月」那一欄。
	financeHitW, financeHitH = 24, 16

	// remake 差異：操作提示自己一個框，接在原版視窗下面。
	financeHintY = financeWinY + financeWinH
	financeHintH = 48
)

// financeRowRect 是第 n 列在「次月」欄的可點矩形（原版熱區 0x20+n）。
func financeRowRect(row int) image.Rectangle {
	if row < 0 || row >= financeRows {
		return image.Rectangle{}
	}
	y := financeRowY + row*financeRowStep
	return image.Rect(financeIconNextX, y, financeIconNextX+financeHitW, y+financeHitH)
}

func (g *game) drawFinance(screen *ebiten.Image) {
	if !g.finance.active {
		return
	}
	g.chrome.Window(screen, financeWinX, financeWinY, financeWinW, financeWinH, chrome.Menu)

	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	labelInk := g.paletteInk(strategyInkDim, color.RGBA{255, 223, 154, 255})
	warnInk := g.paletteInk(strategyInkGauge, color.RGBA{210, 48, 40, 255})
	season := int(g.world.Clock.Season())

	// 靜態層：標籤、垂直線、四個值框（顯示清單場景 1）。
	g.td.Draw(screen, "資金", financeFundsLabelX, financeFundsLabelY, ink)
	g.td.Draw(screen, "收入", financeIncomeLabelX, financeIncomeLabelY, ink)
	g.td.Draw(screen, "支出", financeIncomeLabelX, financeExpenseLabelY, ink)
	vector.DrawFilledRect(screen, financeDividerX, financeDividerY, 1, financeDividerH, ink, false)
	g.td.Draw(screen, "今月底", financeThisLabelX, financeColLabelY, ink)
	g.td.Draw(screen, "次月", financeNextLabelX, financeColLabelY, ink)
	for _, box := range []struct{ x, y, w, h int }{
		{32, 112, 128, 16}, {240, 96, 96, 32}, {32, 160, 144, 64}, {192, 160, 144, 64},
	} {
		vector.DrawFilledRect(screen, float32(box.x), float32(box.y),
			float32(box.w), float32(box.h), color.Black, false)
	}
	g.td.Draw(screen, "稅率", financeThisColX, financeRowY, labelInk)
	g.td.Draw(screen, "徵兵數", financeThisColX, financeRowY+financeRowStep, labelInk)
	g.td.Draw(screen, "稅率", financeNextColX, financeRowY, labelInk)
	g.td.Draw(screen, "徵兵數", financeNextColX, financeRowY+financeRowStep, labelInk)
	g.td.Draw(screen, "％", financePercentThisX, financeRowY, labelInk)
	g.td.Draw(screen, "％", financePercentNextX, financeRowY, labelInk)
	for i := 0; i < financeRows; i++ {
		y := financeRowY + i*financeRowStep
		for _, col := range []struct {
			x     int
			green bool
		}{{financeIconThisX, false}, {financeIconNextX, true}} {
			img, err := g.lib.DOSVResourceIcon(i, col.green, season)
			if err != nil {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(col.x), float64(y))
			screen.DrawImage(ebiten.NewImageFromImage(img), op)
		}
	}

	// 動態層：資金／收入／支出，以及兩欄各四列。
	f := g.world.Factions[g.world.Player]
	fundsInk := ink
	if f.Funds < 0 {
		fundsInk = warnInk
	}
	g.td.Draw(screen, strategyHUDNumber(f.Funds, financeFundsDigits),
		financeFundsValueX, financeFundsValueY, fundsInk)
	// ⚠ 收入在原版是全域 `cs:word_10D02`，由誰算未讀（docs/spec/14 §5）。
	// remake 的月結算得出 `res.Income` 但沒有留下來，所以這一格暫時是 0——
	// **留白比填一個自己算的數字誠實**：那會變成看起來對的假值。
	g.td.Draw(screen, strategyHUDNumber(0, financeAmountDigits),
		financeIncomeValueX, financeIncomeLabelY, ink)
	g.td.Draw(screen, strategyHUDNumber(f.Expense, financeAmountDigits),
		financeIncomeValueX, financeExpenseLabelY, ink)

	cols := []struct {
		x    int
		vals [financeRows]int
	}{
		{financeValueThisX, [financeRows]int{
			g.world.TaxRate,
			g.world.RecruitCap[economy.Cavalry] * strategyReserveMenPerPoint,
			g.world.RecruitCap[economy.Archer] * strategyReserveMenPerPoint,
			g.world.RecruitCap[economy.Infantry] * strategyReserveMenPerPoint,
		}},
		{financeValueNextX, [financeRows]int{
			g.world.NextTaxRate,
			g.world.NextRecruitCap[economy.Cavalry] * strategyReserveMenPerPoint,
			g.world.NextRecruitCap[economy.Archer] * strategyReserveMenPerPoint,
			g.world.NextRecruitCap[economy.Infantry] * strategyReserveMenPerPoint,
		}},
	}
	for _, col := range cols {
		for i, v := range col.vals {
			g.td.Draw(screen, strategyHUDNumber(v, financeRowDigits),
				col.x, financeRowY+i*financeRowStep, ink)
		}
	}

	// ↓ 以下是 **remake 差異**，原版沒有：選取標記與操作提示。
	// 原版用滑鼠點綠色圖示，改值走數值輸入器。
	//
	// 提示放在**原版視窗外面**的另一個框裡：擠進去會蓋到兩欄的值框，
	// 而那個框的矩形是原版數值，不該為了塞說明去動它。
	sel := financeRowRect(g.finance.row)
	vector.StrokeRect(screen, float32(sel.Min.X-1), float32(sel.Min.Y-1),
		float32(sel.Dx()+2), float32(sel.Dy()+2), 1, ink, false)
	g.chrome.Window(screen, financeWinX, financeHintY, financeWinW, financeHintH, chrome.Menu)
	g.td.Draw(screen, "設定值於次月末生效", financeWinX+8, financeHintY+8, labelInk)
	g.td.Draw(screen, "↑↓ 選欄　←→ 增減　ESC 關閉",
		financeWinX+8, financeHintY+8+textdraw.GlyphH+2, labelInk)
}

// openCityList 是命令視窗的「據點」：自勢力據點一覽。
func (g *game) openCityList() {
	rows := g.playerCities()
	if len(rows) == 0 {
		g.lastEvent = "沒有據點"
		return
	}
	g.cityList(rows, "↑↓ 移動　Enter 選取／決定　1-5 排序　ESC 取消",
		func(city int) bool {
			// 說明書：「選了游標移過去」。這裡把鏡頭移到該據點，
			// 並開原版的據點情報視窗（docs/spec/23）。
			c := g.world.Cities[city]
			g.camX, g.camY = c.X-viewCols/2, c.Y-viewRows/2
			g.clampCam()
			g.openCityInfo(city)
			return true
		})
}

// openFactionList 是命令視窗的「勢力」：他勢力一覽。
func (g *game) openFactionList() {
	var rows []int
	for i := range g.world.Factions {
		if g.world.Factions[i].Alive {
			rows = append(rows, i)
		}
	}
	g.factionList(rows, "↑↓ 移動　Enter 選取／決定　1-4 排序　ESC 取消",
		func(f int) bool {
			g.lastEvent = big5(g.world.LordName(f)) + " 軍"
			return true
		})
}
