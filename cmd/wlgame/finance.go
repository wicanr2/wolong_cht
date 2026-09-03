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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// financeState 是財政畫面的狀態。四列：稅率、騎馬、弓兵、步兵
// （原版視窗裡四個綠色按鈕，由上而下就是這個順序）。
type financeState struct {
	active bool
	row    int
	// keyboard：選取框只在用過鍵盤之後畫——原版沒有選取狀態
	// （對拍 playtest/42；與編成視窗的 f.keyboard 同一個先例）。
	keyboard bool
	// editing／value 是開著數值輸入器時的狀態。原版 `sub_17C6E` 的
	// `si` 從 0 起，離開前不寫回任何東西——所以編輯中的值不能直接
	// 塞進 world，取消時會留下痕跡。
	editing bool
	value   int
}

// TaxMax 是稅率的上限：`sub_167CD` 傳給 `sub_17C6E` 的 `ax = 64h`
// （docs/spec/78 §1.2）。
const TaxMax = 100

// recruitMenMax 是三個兵種各自的募集人數上限：`ax = 2710h` ＝ 10,000 人。
// 原版寫回前 `div 10`，所以存進勢力記錄的是**點數**（≤ 1,000）。
const recruitMenMax = 10000

// financeAnchorX／Y 是財政四個熱區共用的數值器錨點
// （`sub_167CD` 等四支的 `dx = 128h`／`bx = 0B8h`）。
const financeAnchorX, financeAnchorY = 296, 184

func (g *game) beginFinance() { g.finance = financeState{active: true} }

// financeRowMax 是第 n 列開數值器時傳進去的上限。
func financeRowMax(row int) int {
	if row == 0 {
		return TaxMax
	}
	return recruitMenMax
}

func (g *game) updateFinance() {
	f := &g.finance
	if f.editing {
		g.updateFinanceAmount()
		return
	}
	switch {
	case g.cancelled():
		f.active = false
		return
	case pressed(ebiten.KeyArrowUp):
		f.row, f.keyboard = (f.row+3)%4, true
	case pressed(ebiten.KeyArrowDown):
		f.row, f.keyboard = (f.row+1)%4, true
	}
	// 原版只能點綠色圖示欄那四格（熱區 0x20–0x23）。
	if row, ok := financeRowAtPointer(); ok {
		f.row = row
		g.beginFinanceAmount(row)
		return
	}
	if pressed(ebiten.KeyEnter) || pressed(ebiten.KeySpace) {
		f.keyboard = true
		g.beginFinanceAmount(f.row)
	}
}

func financeRowAtPointer() (int, bool) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return 0, false
	}
	x, y := ebiten.CursorPosition()
	for row := 0; row < financeRows; row++ {
		if image.Pt(x, y).In(financeRowRect(row)) {
			return row, true
		}
	}
	return 0, false
}

func (g *game) beginFinanceAmount(row int) {
	g.finance.editing, g.finance.value, g.finance.row = true, 0, row
	g.beginAmountEditor(amountCursorFinance, financeAnchorX, financeAnchorY)
}

// updateFinanceAmount 是四個熱區共用的那一段：`sub_17C6E` 的操作迴圈。
// 取消（右鍵／ESC）走原版 `jb locret_167E5` ——**什麼都不寫回**。
func (g *game) updateFinanceAmount() {
	f := &g.finance
	if g.cancelled() {
		f.editing = false
		return
	}
	max := financeRowMax(f.row)
	apply := func(edit state.AmountEdit, digit int) {
		if v, ok := state.EditAmountValue(f.value, max, edit, digit); ok {
			f.value = v
		}
	}
	if button, ok := g.amountPointerButton(); ok {
		if button.edit == state.AmountFinishInput {
			g.commitFinanceAmount()
			return
		}
		apply(button.edit, button.digit)
		return
	}
	if digit, ok := pressedAmountDigit(); ok {
		g.setAmountCursorAction(state.AmountAppendDigit, digit)
		apply(state.AmountAppendDigit, digit)
		return
	}
	edit := state.AmountEdit(255)
	switch {
	case pressed(ebiten.KeyBackspace):
		edit = state.AmountDeleteDigit
	case pressed(ebiten.KeyInsert):
		edit = state.AmountAppendHundred
	case pressed(ebiten.KeyDelete):
		edit = state.AmountClear
	case pressed(ebiten.KeyHome):
		edit = state.AmountSetMax
	case pressed(ebiten.KeyEnter), pressed(ebiten.KeySpace):
		g.setAmountCursorAction(state.AmountFinishInput, 0)
		g.commitFinanceAmount()
		return
	}
	if edit != state.AmountEdit(255) {
		g.setAmountCursorAction(edit, 0)
		apply(edit, 0)
	}
}

// commitFinanceAmount 是四支 handler 的寫回：稅率直接存，
// 三個兵種的**人數要先除以 10**（原版 `mov bx, 0Ah / div bx`）。
func (g *game) commitFinanceAmount() {
	f := &g.finance
	f.editing = false
	if f.row == 0 {
		g.world.NextTaxRate = clamp(f.value, 0, TaxMax)
		return
	}
	t := f.row - 1
	g.world.NextRecruitCap[t] = f.value / strategyReserveMenPerPoint
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
	// 數值一律用色 9（米黃）：2026-08-24 實機對拍 orig-w4-finance.png，
	// 資金／收入／支出／兩欄各列／％ 全是 (243,211,146) ＝ 調色盤索引 9；
	// 白色（索引 15）只用在「資金」「收入」等標籤。
	valueInk := labelInk
	f := g.world.Factions[g.world.Player]
	fundsInk := valueInk
	if f.Funds < 0 {
		fundsInk = warnInk
	}
	// 數字用原版 8×16 字模（ICONGRF 段 3 +0x840）：2026-08-24 對拍
	// orig-w4-finance.png 的墨水列 113..126 ＝ mask topY 112，倚天 ASCII
	// 只有 9 列、字形不同，同位置同色仍逐像素差得出來。
	g.drawOriginalNumber(screen, f.Funds,
		financeFundsValueX, financeFundsValueY, financeFundsDigits, fundsInk)
	// ⚠ 收入在原版是全域 `cs:word_10D02`，由誰算未讀（docs/spec/14 §5）。
	// remake 的月結算得出 `res.Income` 但沒有留下來，所以這一格暫時是 0——
	// **留白比填一個自己算的數字誠實**：那會變成看起來對的假值。
	// 收入與支出是**月結快照**（原版 sub_1548F 寫的顯示用全域，
	// docs/spec/14 §5）：讀檔後、月結前顯示 0。
	g.drawOriginalNumber(screen, g.world.IncomeSnap,
		financeIncomeValueX, financeIncomeLabelY, financeAmountDigits, valueInk)
	g.drawOriginalNumber(screen, g.world.ExpenseSnap,
		financeIncomeValueX, financeExpenseLabelY, financeAmountDigits, valueInk)

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
			g.drawOriginalNumber(screen, v,
				col.x, financeRowY+i*financeRowStep, financeRowDigits, valueInk)
		}
	}

	// ↓ 以下是 **remake 差異**，原版沒有：選取標記與操作提示。
	// 原版用滑鼠點綠色圖示，改值走數值輸入器。
	//
	// 提示放在**原版視窗外面**的另一個框裡：擠進去會蓋到兩欄的值框，
	// 而那個框的矩形是原版數值，不該為了塞說明去動它。
	if g.finance.keyboard {
		sel := financeRowRect(g.finance.row)
		vector.StrokeRect(screen, float32(sel.Min.X-1), float32(sel.Min.Y-1),
			float32(sel.Dx()+2), float32(sel.Dy()+2), 1, ink, false)
	}
	g.chrome.Window(screen, financeWinX, financeHintY, financeWinW, financeHintH, chrome.Menu)
	g.td.Draw(screen, "設定值於次月末生效", financeWinX+8, financeHintY+8, labelInk)
	g.td.Draw(screen, "↑↓ 選欄　Enter 輸入　ESC 關閉",
		financeWinX+8, financeHintY+8+textdraw.GlyphH+2, labelInk)

	// 數值輸入器疊在最上面（原版 sub_17C6E 會先存下底下的畫面）。
	if g.finance.editing {
		g.drawAmountPanel(screen, g.finance.value, true)
	}
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
			// **與「首都確認」共用這一支**（docs/spec/126 §1.1）。
			g.focusCity(city)
			return true
		})
}

// openCityCommandMenu 是指令列第 6 格（「據點」）。
//
// ⚠ 原版點下去**不是**直接開一覽，而是先跳兩列選單
// （`sub_193E9(ax=2, cx=52h, dx=40Fh)`，docs/spec/126）。
func (g *game) openCityCommandMenu() { g.openPopupMenu(cityPopupMenu) }

// dispatchCityMenu 是「據點」那兩列各自接到哪。
func (g *game) dispatchCityMenu(row int) {
	switch row {
	case 0:
		g.beginLocateCapital()
	case 1:
		g.openCityList()
	}
}

// beginLocateCapital 是「首都確認」：鏡頭移到自己的首都並開情報視窗。
//
// ⭐ 原版與「據點一覽」**共用尾段**（`loc_1633C`）——差別只在
// 城是自己挑的還是從勢力記錄 `+3` 直接取的（docs/spec/126 §1.1）。
func (g *game) beginLocateCapital() {
	p := g.world.Player
	if p < 0 || p >= len(g.world.Factions) {
		return
	}
	capital := g.world.Factions[p].Capital
	if capital < 0 || capital >= len(g.world.Cities) {
		g.lastEvent = "沒有首都"
		return
	}
	g.focusCity(capital)
}

// focusCity 是那條共用尾段：鏡頭 ＝ 據點 −(20,12)，再開情報視窗。
func (g *game) focusCity(city int) {
	c := g.world.Cities[city]
	g.camX, g.camY = c.X-centreCol, c.Y-centreRow
	g.clampCam()
	g.openCityInfo(city)
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
