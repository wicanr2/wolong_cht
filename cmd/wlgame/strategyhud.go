package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// DOS/V 自然策略畫面的固定骨架。**每一個數字都出自機器碼**
// （`docs/spec/12-strategy-chrome.md` §1，來源 `docs/re/47`）：
//
//	橫幅        (0,   0, 640,  32)   sub_18755 的熱區 6
//	命令        (0,  32, 432,  32)   sub_1614A → sub_1895D(cx=0x021B)
//	縮小地圖    (432, 32, 208, 160)  sub_15A3A → sub_1895D(cx=0x0A0D)
//	自勢力情報  (432,192, 208, 208)  sub_15E2D → sub_1895D(cx=0x0D0D)
//
// 右欄三段相加 32 + 160 + 208 = 400，剛好鋪滿畫面高度——這是「矩形讀對了」
// 的算術檢查。四個視窗可以逐個開關（docs/spec/13）。
const (
	strategyCommandY = bannerH
	strategyCommandH = 32
	strategyCommandW = 432 // sub_1895D 的 cx 低 byte 0x1B → 27 × 16
	// ⭐ 大地圖鋪滿橫幅以下的**整片畫面**，四個視窗疊在它上面
	// （`sub_1D615` 的迴圈是 40 欄 × 23 列的 16×16 格，docs/re/47 §3.2）。
	// 這不是「左邊地圖、右邊面板」的分割版面——關掉右欄看到的是地圖，不是空白。
	strategyMapY              = bannerH
	strategyMapW              = screenW
	strategyMapH              = screenH - bannerH
	strategySidebarW          = 208 // sub_1895D 的 cx 低 byte 0x0D → 13 × 16
	strategySidebarX          = screenW - strategySidebarW
	strategySidebarInnerX     = strategySidebarX + chrome.Tile
	strategySidebarInnerRight = strategySidebarX + strategySidebarW - chrome.Tile
	strategyMinimapX          = strategySidebarInnerX
	strategyMinimapY          = bannerH + chrome.Tile
	strategyMinimapW          = strategySidebarW - 2*chrome.Tile
	// 縮小地圖視窗的內框是 192×144：上面 128 px 是地圖本體
	// （熱區 0x16 ＝ (440,40,192,128)），下面 16 px 是勢力篩選列
	// （熱區 0x17 ＝ (536,168,96,16)，只有右半格可點）。
	strategyMinimapH       = 160
	strategyMinimapImageH  = 128
	strategyMinimapLegendY = strategyMinimapY + strategyMinimapImageH
	strategyMinimapSwatchY = strategyMinimapLegendY + 4
	strategyFactionY       = bannerH + strategyMinimapH
	strategyFactionH       = screenH - strategyFactionY
	strategyFactionInnerY  = strategyFactionY + chrome.Tile
	strategyFactionInnerW  = strategySidebarW - 2*chrome.Tile

	// 命令列的版面是**從原版讀出來的**，不是從影片量的（docs/re/47 §4.1）：
	//
	//	sub_1614A: sub_106F5(si = cs:6181h, bx = 0x28, dx = 8, ax = 0x0F01)
	//	           sub_1E3D7(al = 0x0C, bx = 0x28, dx = 0x18, cx = 0x0230)
	//
	// **`dx` 是 X、`bx` 是 Y**（`sub_1F6DC` 每畫一個字就 `add dx, di`），
	// 所以字串起點是 X=8、Y=40，熱區是 (24, 40, 384, 16)。
	//
	// `cs:6181h` 是**一整串**，不是八個標籤：
	//
	//	「　 進言　人事　財政　編成　軍團　據點　武將　勢力 　」
	//	  ↑全形 ↑半形      ↑ 兩個字之間沒有空格，詞與詞之間才有一個全形空格
	//
	// 節距 ＝ 32（兩個全形字）＋ 16（一個全形空格）＝ **48**，
	// 第一個字在 8 ＋ 24 ＝ **32**。驗算尾端：8 + 24 + 8×32 + 7×16 + 24 = 424，
	// 落在框寬 432 之內；起點若讀成 40 會算到 456 而溢出框外。
	strategyCommandX     = 8
	strategyCommandLead  = 24 // 開頭的全形空格 ＋ 半形空格
	strategyCommandCellW = 48
	strategyCommandTextW = 32
	// 命中判定照原版的 sub_161CA：索引 ＝ (x − 24) ÷ 48，
	// 所以第 n 格是 [24 + 48n, 72 + 48n)，Y 只命中 40–56 的文字列。
	strategyCommandHitX        = 24
	strategyCommandHitY        = strategyCommandY + chrome.Tile
	strategyCommandHitH        = strategyCommandH - 2*chrome.Tile
	// 自勢力情報視窗的內部座標，出自 docs/re/47 §4.2（相對視窗左上角 432,192）：
	//
	//	底圖（頭像＋標籤）(440,200)   sub_10337(dx=0x1B8, bx=0xC8)
	//	君主／首都／軍師   (576, 208/224/240)  sub_106FD(dx=0x240, bx=0xD0/E0/F0)
	//	信賴度量條         (456,292)   sub_10AAA(dx=0x1C8, bx=0x124)
	//	資金 7 位          (560,312)   sub_1062F(di=0x61C6) → 一列 80 B
	//	預備兵 ×3 6 位     (568, 328/344/360)  di=0x66C7、每列 +0x500
	//
	// 標籤與底層方塊的座標同樣是原版數值，來自**顯示清單**
	// （`sub_10337(al=0)` 的場景 0，docs/re/48 §3）：
	//
	//	信賴度量條的槽  (448,288) 176×10   op 03
	//	資金／預備兵黑底 (448,304) 176×80   op 03
	//	「君主／首都／軍師」(512, 208/224/240)  op 08
	//	「信賴度」        (448,272)          op 08
	//	「資金」「預備兵」 (456, 312/328)     op 08
	//	四張 24×16 圖形   (528, 312/328/344/360)  op 09
	strategyInfoYOffset      = 16 // 208 − 192
	strategyInfoRowStep      = 16 // 0xE0 − 0xD0
	strategyInfoLabelXOffset = 80 // 512 − 432
	strategyInfoLabelW       = 32
	strategyInfoValueXOffset = 144 // 576 − 432
	strategyInfoValueW       = strategySidebarInnerRight - (strategySidebarX + strategyInfoValueXOffset)
	strategyTrustLabelX      = 16  // 448 − 432
	strategyTrustLabelY      = 80  // 272 − 192
	strategyTrustSlotX       = 16  // 448 − 432
	strategyTrustSlotY       = 96  // 288 − 192
	strategyTrustSlotW       = 176 // 623 − 448 + 1
	strategyTrustSlotH       = 10
	strategyTrustYOffset     = 100 // 292 − 192
	strategyTrustXOffset     = 24  // 456 − 432
	strategyTrustMaxW        = 160
	strategyResourceBoxX     = 16  // 448 − 432
	strategyResourceBoxY     = 112 // 304 − 192
	strategyResourceBoxW     = 176
	strategyResourceBoxH     = 80
	strategyResourceLabelX   = 24  // 456 − 432
	strategyFundsYOffset     = 120 // 312 − 192
	strategyFundsXOffset     = 128 // 560 − 432
	strategyReserveYOffset   = 136 // 328 − 192
	strategyReserveXOffset   = 136 // 568 − 432
	strategyIconXOffset      = 96  // 528 − 432：四張 24×16 圖形的欄
	strategyResourceRowStep  = 16
	// 原版是資金 7 位、預備兵 6 位，兩者右端都對齊 x=616。
	strategyFundsDigits   = 7
	strategyReserveDigits = 6
	// 原版顯示預備兵時乘 10（`sub_15F7F` 的 `mul dx(=0x0A)`）：勢力記錄存的是
	// 「點」，一點 10 人。**這裡只換算顯示**——規則層的 `Reserves` 語意
	// 目前與原版不一致（`army.MenPerUnit` 是 1000，等於把存值當人數扣），
	// 那要另外開規格處理，不在畫面這一輪動它（docs/re/47 §5）。
	strategyReserveMenPerPoint = 10
	strategyNumberSlots   = strategyFundsDigits
	strategyNumberW       = strategyNumberSlots * textdraw.HalfW
)

// naturalCommandLabels 是原版 `cs:6181h` 那一串裡的八個詞。
// **兩個字之間沒有空格**——先前寫成「進　言」是照影片猜的，
// 那讓每個詞多佔 16 px，八個詞累積起來就與原版對不上（docs/re/46 §1）。
var naturalCommandLabels = [...]string{
	"進言", "人事", "財政", "編成",
	"軍團", "據點", "武將", "勢力",
}

// naturalCommandID 依畫面標籤命名，不依快捷鍵順序命名。快捷鍵只是另一個
// 輸入來源；滑鼠 cell 與鍵盤都必須先得到同一個 ID，再進入同一個 action。
type naturalCommandID uint8

const (
	naturalCommandAdvise naturalCommandID = iota
	naturalCommandPersonnel
	naturalCommandFinance
	naturalCommandFormation
	naturalCommandCorps
	naturalCommandCity
	naturalCommandGeneral
	naturalCommandFaction
)

var naturalCommandActions = [...]func(*game){
	(*game).openAdvise,
	(*game).openPersonnel,
	(*game).beginFinance,
	(*game).beginForm,
	(*game).openCorpsList,
	(*game).openCityList,
	(*game).openGeneralList,
	(*game).openFactionList,
}

// strategyCommandCellRect 是自然策略頂端命令列的純幾何契約。
// 不讀 game、不讀輸入狀態；外框、地圖與右側 HUD 都不在矩形內。
//
// 分格照原版 `sub_161CA` 的 `索引 ＝ (x − 24) ÷ 48`：**格與格之間沒有間隙**，
// 24–408 這一段每一個像素都屬於某一格。原版另外用 40 px 寬畫高亮，
// 那是視覺不是命中範圍——照高亮寬度做命中會在每一格右緣留 8 px 死區。
func strategyCommandCellRect(index int) image.Rectangle {
	if index < 0 || index >= len(naturalCommandLabels) {
		return image.Rectangle{}
	}
	x := strategyCommandHitX + index*strategyCommandCellW
	return image.Rect(x, strategyCommandHitY, x+strategyCommandCellW,
		strategyCommandHitY+strategyCommandHitH)
}

func hitTestNaturalCommand(x, y int) (naturalCommandID, bool) {
	for i := range naturalCommandLabels {
		if image.Pt(x, y).In(strategyCommandCellRect(i)) {
			return naturalCommandID(i), true
		}
	}
	return 0, false
}

// strategyCommandForShortcut 僅負責把既有鍵盤快捷鍵映射到畫面標籤 ID。
// 特別保留畫面順序與快捷鍵順序的分離：M（行軍）不屬於頂端八格。
func strategyCommandForShortcut(key ebiten.Key) (naturalCommandID, bool) {
	switch key {
	case ebiten.KeyP:
		return naturalCommandAdvise, true
	case ebiten.KeyJ:
		return naturalCommandPersonnel, true
	case ebiten.KeyF:
		return naturalCommandFinance, true
	case ebiten.KeyA:
		return naturalCommandFormation, true
	case ebiten.KeyC:
		return naturalCommandCorps, true
	case ebiten.KeyT:
		return naturalCommandCity, true
	case ebiten.KeyG:
		return naturalCommandGeneral, true
	case ebiten.KeyK:
		return naturalCommandFaction, true
	default:
		return 0, false
	}
}

func (g *game) dispatchNaturalCommand(command naturalCommandID) bool {
	if int(command) < 0 || int(command) >= len(naturalCommandActions) {
		return false
	}
	naturalCommandActions[command](g)
	return true
}

// strategyHUDSingleLine 是常駐欄位的 fail-safe。正常 DOS/V 名稱不會進入截斷路徑；
// 若測試 fixture 或損壞資料帶入過長字串，不能讓字模穿過下一欄或下一列。
func strategyHUDSingleLine(s string, maxPixels int) string {
	if maxPixels <= 0 {
		return ""
	}
	if textdraw.StringWidth(s) <= maxPixels {
		return s
	}
	width := 0
	runes := make([]rune, 0, len([]rune(s)))
	for _, ch := range s {
		w := textdraw.RuneWidth(ch)
		if width+w > maxPixels {
			break
		}
		runes = append(runes, ch)
		width += w
	}
	if len(runes) == 0 {
		return "?"
	}
	return string(runes)
}

// strategyHUDNumber 把數值排進固定的半形數字槽位。原版是資金 7 位、
// 預備兵 6 位（`sub_1062F` 的 `bl`，docs/re/47 §4.2），兩者右端都對齊。
// 超出槽位時飽和顯示——這不是原版數值規則，只是不讓數字撐破欄位；
// 真實 state 不會被改寫。
func strategyHUDNumber(value, digits int) string {
	max := 1
	for i := 0; i < digits; i++ {
		max *= 10
	}
	if value >= max {
		value = max - 1
	}
	if min := -(max/10 - 1); value < min {
		value = min
	}
	return fmt.Sprintf("%*d", digits, value)
}

// hudWindow 是主畫面上四個**可開關**的常駐視窗，對應原版 `byte_198A6`
// 的 bit 0–3（`docs/spec/13-main-window-toggles.md`）。
type hudWindow uint8

const (
	hudCommand hudWindow = 1 << iota // bit 0：命令
	hudFaction                       // bit 1：自勢力情報
	hudMinimap                       // bit 2：縮小地圖
	hudSystem                        // bit 3：系統（開著時時間停止）
)

// hudSwitchRect 是橫幅右側五格開關的矩形。原版由 `sub_18755` 的迴圈
// 登記：起點 x=336、每格 32×32、共五格，Y 是 0–32（docs/re/47 §2）。
//
// 第五格原版接 `nullsub_1`，這裡照樣不接東西——**留著它是刻意的**，
// 少一格會讓熱區編號與原版對不起來，之後想驗就得重推一次。
const (
	hudSwitchX0 = 336
	hudSwitchW  = 32
	hudSwitchN  = 5
)

// hudSwitchWindow 把開關編號（0 起）換成它控制的視窗；第五格回 0。
func hudSwitchWindow(index int) hudWindow {
	switch index {
	case 0:
		return hudCommand
	case 1:
		return hudFaction
	case 2:
		return hudMinimap
	case 3:
		return hudSystem
	}
	return 0
}

// hitTestHUDSwitch 回傳游標落在第幾格開關。
func hitTestHUDSwitch(x, y int) (int, bool) {
	if y < 0 || y >= bannerH || x < hudSwitchX0 {
		return 0, false
	}
	i := (x - hudSwitchX0) / hudSwitchW
	if i >= hudSwitchN {
		return 0, false
	}
	return i, true
}

// hudOpen／hudSet 是四個視窗開關狀態的唯一入口。
//
// 系統視窗**不另存一個位元**：remake 早就有 `g.open[winSystem]`
// 這個模態視窗，再開一個 bit 會變成兩份狀態各自漂移
// （`CLAUDE.md` §7 第 6 條）。
func (g *game) hudOpen(w hudWindow) bool {
	if w == hudSystem {
		return g.open[winSystem]
	}
	return g.hud&w != 0
}

func (g *game) hudSet(w hudWindow, open bool) {
	if w == hudSystem {
		g.open[winSystem] = open
		return
	}
	if open {
		g.hud |= w
	} else {
		g.hud &^= w
	}
}

// drawNaturalStrategyHUD 畫主畫面的四個常駐視窗。
//
// 版面常數全部出自機器碼（`docs/spec/12-strategy-chrome.md` §1）；
// 視窗是**疊在大地圖上的不透明層**，關掉哪個就露出底下的地圖
// （docs/re/47 §3.2）。
func (g *game) drawNaturalStrategyHUD(screen *ebiten.Image) {
	if !g.hudOpen(hudCommand) {
		g.drawHUDSidebar(screen)
		return
	}
	// 命令列使用與原版視窗相同的深藍底／紅金外框；文字沿 8 px 邊界排列。
	g.chrome.Window(screen, 0, strategyCommandY, strategyCommandW, strategyCommandH, chrome.Menu)
	for i, label := range naturalCommandLabels {
		x := strategyCommandX + strategyCommandLead + i*strategyCommandCellW
		g.td.Draw(screen, strategyHUDSingleLine(label, strategyCommandTextW), x, strategyCommandY+8, chrome.Paper)
	}
	g.drawHUDSidebar(screen)
}

// drawHUDSidebar 畫右欄的兩個視窗。**兩個各自可以關掉**，
// 所以不能靠「上面那個蓋住下面那個的邊」來省一條分隔線——
// 原版的兩個矩形是背靠背相鄰的（192 ＝ 32 + 160）。
func (g *game) drawHUDSidebar(screen *ebiten.Image) {
	if g.hudOpen(hudFaction) {
		g.chrome.Window(screen, strategySidebarX, strategyFactionY, strategySidebarW, strategyFactionH, chrome.Menu)
		g.drawNaturalFactionHUD(screen, strategySidebarX, strategyFactionY)
	}
	if !g.hudOpen(hudMinimap) {
		return
	}

	// 右上縮小地圖。ICONGRF 段 2 本身就是 192×128，不能用大地圖降採樣替代。
	g.chrome.Window(screen, strategySidebarX, bannerH, strategySidebarW, strategyMinimapH, chrome.Menu)
	if img, err := g.lib.Render(g.minimapAsset(), 0, int(g.world.Clock.Season())); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(strategyMinimapX), float64(strategyMinimapY))
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	// 原版縮小地圖下方是兩個半欄寬的勢力色標；圖像本身不含 runtime
	// 君主名，所以這一列由 state 填入。先前只有 8×8 色點，會讓右欄骨架
	// 看起來少掉一整列 16 px 的紅／藍帶。
	red := color.RGBA{210, 48, 40, 255}
	blue := color.RGBA{45, 105, 210, 255}
	legendX := strategySidebarX + chrome.Tile
	vector.DrawFilledRect(screen, float32(legendX), float32(strategyMinimapLegendY), 96, 16, red, false)
	vector.DrawFilledRect(screen, float32(legendX+96), float32(strategyMinimapLegendY), 96, 16, blue, false)
	vector.DrawFilledRect(screen, float32(legendX+8), float32(strategyMinimapSwatchY), 8, 8,
		color.RGBA{85, 154, 69, 255}, false)
	vector.DrawFilledRect(screen, float32(legendX+104), float32(strategyMinimapSwatchY), 8, 8,
		color.RGBA{85, 154, 69, 255}, false)
	g.td.Draw(screen, strategyHUDSingleLine(big5(g.world.LordName(g.world.Player)), 72), legendX+24, strategyMinimapLegendY, chrome.Paper)
	enemyName := "敵"
	for faction, f := range g.world.Factions {
		if faction != g.world.Player && f.Alive {
			enemyName = big5(g.world.LordName(faction))
			break
		}
	}
	g.td.Draw(screen, strategyHUDSingleLine(enemyName, 72), legendX+120, strategyMinimapLegendY, chrome.Paper)
}

func (g *game) drawNaturalFactionHUD(dst *ebiten.Image, x, y int) {
	p := g.world.Player
	f := g.world.Factions[p]
	season := int(g.world.Clock.Season())

	if lord := f.Lord; lord >= 0 && lord < len(g.world.Generals) {
		if img, err := g.lib.Portrait(g.world.Generals[lord].Portrait, season); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(strategySidebarInnerX), float64(strategyFactionInnerY))
			dst.DrawImage(ebiten.NewImageFromImage(img), op)
		}
	}

	ink := chrome.Paper
	// 原版頭像右側是「君主／首都／軍師」三列。值的座標是原版數值
	// （576, 208/224/240）；標籤與分隔線在原版屬於底圖，位置是這裡自己定的。
	labelX := x + strategyInfoLabelXOffset
	valueX := x + strategyInfoValueXOffset
	infoY := y + strategyInfoYOffset
	capital := "－"
	if f.Capital >= 0 && f.Capital < len(g.world.Cities) {
		capital = big5(g.world.Cities[f.Capital].Name)
	}
	info := [...]struct{ label, value string }{
		{"君主", big5(g.world.LordName(p))},
		{"首都", capital},
		{"軍師", big5(g.advisorName())},
	}
	for i, row := range info {
		py := infoY + i*strategyInfoRowStep
		g.td.Draw(dst, strategyHUDSingleLine(row.label, strategyInfoLabelW), labelX, py, ink)
		g.td.Draw(dst, strategyHUDSingleLine(row.value, strategyInfoValueW), valueX, py, ink)
	}
	// 信賴度：原版先畫一個 176×10 的槽（顯示清單 op 03），再在裡面畫量條
	// （`sub_15F27`，長度 (信賴度×100 + 0x9F) ÷ 0xA0）。顏色是自己挑的，
	// 原版的 `ax=0x0A00` 還沒對過調色盤。
	g.td.Draw(dst, "信賴度", x+strategyTrustLabelX, y+strategyTrustLabelY, ink)
	vector.DrawFilledRect(dst, float32(x+strategyTrustSlotX), float32(y+strategyTrustSlotY),
		strategyTrustSlotW, strategyTrustSlotH, color.RGBA{24, 24, 32, 255}, false)
	trustW := (g.world.Trust*100 + 0x9F) / 0xA0
	if trustW > strategyTrustMaxW {
		trustW = strategyTrustMaxW
	}
	if trustW > 0 {
		vector.DrawFilledRect(dst, float32(x+strategyTrustXOffset), float32(y+strategyTrustYOffset),
			float32(trustW), 8, color.RGBA{85, 154, 69, 255}, false)
	}

	// 資源區的黑底也是顯示清單畫的（op 03，(448,304) 176×80）。
	vector.DrawFilledRect(dst, float32(x+strategyResourceBoxX), float32(y+strategyResourceBoxY),
		strategyResourceBoxW, strategyResourceBoxH, color.Black, false)
	resourceInk := color.RGBA{255, 223, 154, 255}
	fundsY := y + strategyFundsYOffset
	g.td.Draw(dst, "資金", x+strategyResourceLabelX, fundsY, resourceInk)
	g.td.Draw(dst, "預備兵", x+strategyResourceLabelX, y+strategyReserveYOffset, resourceInk)
	g.td.Draw(dst, strategyHUDNumber(f.Funds, strategyFundsDigits), x+strategyFundsXOffset, fundsY, ink)

	// 顯示換算見常數區的 strategyReserveMenPerPoint。
	reserveValues := [...]int{
		f.Reserves[economy.Cavalry],
		f.Reserves[economy.Archer],
		f.Reserves[economy.Infantry],
	}
	// 原版在 (528, 312/328/344/360) 各貼一張 24×16 圖形（顯示清單 op 09，
	// 圖庫位移 0x1200 起連號）。那四張是 ICONGRF 段 3 的天秤／馬／弓／步，
	// 位址換算見 docs/re/48 §6——**用原版素材，不自繪**。
	for i := 0; i < 1+len(reserveValues); i++ {
		iconY := y + strategyFundsYOffset + i*strategyResourceRowStep
		img, err := g.lib.DOSVResourceIcon(i, false, season)
		if err != nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+strategyIconXOffset), float64(iconY))
		dst.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	for i, value := range reserveValues {
		iconY := y + strategyReserveYOffset + i*strategyResourceRowStep
		g.td.Draw(dst, strategyHUDNumber(value*strategyReserveMenPerPoint, strategyReserveDigits),
			x+strategyReserveXOffset, iconY, ink)
	}
}
