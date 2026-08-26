package main

import (
	"github.com/wicanr2/wolong_cht/internal/rules/speed"
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
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
	// 視野框直接照 `sub_15C58` 的 `camX/2 + 440`：鏡頭變數就是畫面上
	// x=0 那一欄，沒有偏移（docs/spec/55 §2）。
	minimapCamBias = 0
	// 圖例上兩個君主名的左緣，原版數值（docs/re/62 §4.1）。
	strategyLegendSelfX    = 480
	strategyLegendWatchedX = 576
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
	// 標籤與名字之間那條垂直線：顯示清單 op 06，(560,208) 長 48、顏色 0x0F。
	strategyInfoDividerXOffset = 128 // 560 − 432
	strategyInfoDividerH       = 48
	strategyTrustLabelX      = 16  // 448 − 432
	strategyTrustLabelY      = 80  // 272 − 192
	strategyTrustSlotX       = 16  // 448 − 432
	strategyTrustSlotY       = 96  // 288 − 192
	strategyTrustSlotW       = 176 // 623 − 448 + 1
	strategyTrustSlotH       = 10
	strategyTrustYOffset     = 100 // 292 − 192
	strategyTrustXOffset     = 24  // 456 − 432
	strategyTrustMaxW        = 160
	// 量條高 2 px：`sub_10AAA` 兩次呼叫都帶 `ch = 2`。滿長 160 來自
	// `sub_15F27` 沒有重載的 `ch = 0A0h`（cx ＝ 總長<<8 ｜ 已填長度）。
	strategyTrustBarH = 2
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

// 原版直接給的調色盤索引。**不要把 RGB 常數抄進呈現層**——
// 調色盤有四季四組，抄死的顏色只會在其中一組看起來對。
const (
	strategyInkNormal = 0x0F
	strategyInkDim    = 0x09
	strategyInkGauge  = 0x0A
)

// paletteInk 取原版調色盤的指定色；取不到就用 fallback，不讓畫面消失。
func (g *game) paletteInk(index int, fallback color.RGBA) color.RGBA {
	if g.lib == nil || g.world == nil {
		return fallback
	}
	c, err := g.lib.PaletteColor(int(g.world.Clock.Season()), index)
	if err != nil {
		return fallback
	}
	return c
}

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

// hitTestMinimapLegend 是縮小地圖圖例的右半格（原版熱區 0x17）。
// 原版是 (536, 168, 96, 16)，**只有右半格可點**——左半格是自勢力，
// 沒得選（docs/re/62 §1）。
func hitTestMinimapLegend(x, y int) bool {
	x0 := strategyMinimapX + strategyMinimapW/2
	return x >= x0 && x < strategyMinimapX+strategyMinimapW &&
		y >= strategyMinimapLegendY && y < strategyMinimapLegendY+16
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

// hudOpen／hudSet 是主畫面視窗開關狀態的**唯一**入口。
//
// ⭐ 先前系統視窗另外存在 `g.open[winSystem]`，而 `g.open[]` 那一整套
// 又自己畫了一次命令／自勢力情報／縮小地圖——兩套疊在畫面上。
// 舊那套已整個拿掉（docs/spec/13 §2.5），四個視窗都在這一個 bitmask 裡。
func (g *game) hudOpen(w hudWindow) bool { return g.hud&w != 0 }

func (g *game) hudSet(w hudWindow, open bool) {
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
		if g.hudOpen(hudSystem) {
			g.drawSystemWindow(screen)
		}
		g.drawFactionPicker(screen)
		return
	}
	// 外框與其他視窗相同，**但內部是純黑、沒有龍紋**（docs/spec/54 §2）。
	g.chrome.Window(screen, 0, strategyCommandY, strategyCommandW, strategyCommandH, chrome.Blank)
	for i, label := range naturalCommandLabels {
		x := strategyCommandX + strategyCommandLead + i*strategyCommandCellW
		g.td.Draw(screen, strategyHUDSingleLine(label, strategyCommandTextW), x, strategyCommandY+8, chrome.Paper)
	}
	g.drawHUDSidebar(screen)
	if g.hudOpen(hudSystem) {
		g.drawSystemWindow(screen)
	}
	g.drawFactionPicker(screen)
}

// 系統選單的版面出自原版（docs/spec/13 §2.6）：視窗矩形來自
// `sub_160CC` → `sub_1895D(cx=0C0Dh)`，內容是顯示清單場景 2，
// 六個熱區由那一支的迴圈登記（docs/re/55）。
//
// ⚠ **只有前六列出自原版。** 第七列「主君編成」是 remake 加的開關
// （docs/spec/76），視窗因此比原版高一個列距。
// 原版那六列的座標一格都沒動——新列只加在下面。
const (
	sysWinX, sysWinY = 208, 112
	// ⚠ **原版是 192（六列）。** remake 多了「主君編成」那一列
	// （docs/spec/76，使用者裁定），視窗因此加高一個列距。
	// 代價寫在 docs/playtest/39：系統選單開著時不再與原版逐像素相同。
	sysWinW, sysWinH = 208, 192+sysRowStep

	sysTitleX, sysTitleY = 228, 124
	sysRuleX, sysRuleY   = 216, 142
	sysRuleW             = 191

	// 六列：標籤的黑底 ＋ 標籤 ＋ 值格。列距 24。
	sysLabelBoxX, sysLabelBoxY = 222, 150
	sysLabelBoxW, sysLabelBoxH = 128, 20
	sysLabelX, sysLabelY       = 232, 152
	sysValueX, sysValueY       = 352, 152
	sysValueW, sysValueH       = 48, 16
	sysRowStep                 = 24
	// ⚠ 原版六列；第 7 列是 remake 加的「主君編成」（docs/spec/76）。
	sysRows = 7
)

// sysMenuLabels 是六列的標籤，取自顯示清單場景 2 的字串。
var sysMenuLabels = [sysRows]string{
	"資料儲存", "畫面模式", "音　　效", "戰略速度", "戰術速度", "遊戲結束", "主君編成"}

// 系統選單六列的索引。原版的六個 handler 沒讀（docs/re/55 §4），
// 所以**哪一列做什麼是 remake 自己接的**，照標籤的字面意思。
const (
	sysRowSave = iota
	sysRowVideo
	sysRowSound
	sysRowStrategySpeed
	sysRowTacticalSpeed
	sysRowQuit
	// ⭐ sysRowLordCorps 是 remake 加的第 7 列，原版沒有（docs/spec/76）。
	//
	// ⚠ **加在最後，不是插在「遊戲結束」之前。** 插進去會把原版最後那一列
	// 往下推一格，於是**原版六列沒有一列還在原座標**；加在後面的話那六列
	// 一個 px 都不動，docs/playtest/39 對那六列的比對仍然成立。
	// 代價是「遊戲結束」不再是視覺上的最後一列——換到的是原版版面不被動到。
	sysRowLordCorps
)

// videoModeLabels 是「畫面模式」那一列的兩個選項（原版字串表 `ds:6002h`）。
// 值切換的是 `GAMEPAL.BRG` 的四季組：0 ＝ bank 0–3、1 ＝ bank 4–7
// （`sub_16097` → `sub_109AF`，docs/re/02 §6.2）。
// **remake 只做了第 0 組**——液晶那一組是給 8 階調液晶用的高飽和純色，
// 現代螢幕上沒有對應的顯示器，接上去也沒有可對照的原版畫面。
var videoModeLabels = [2]string{"１６色", " 液晶 "}

// speed.Labels 是原版兩個速度設定的**五個檔位**，字串在 `ds:6033h` 起
// 五筆各 7 bytes（Big5，含前後的半形空白）。
//
// 檔位數不是猜的：系統設定表 `ds:5FF2h` 每筆的第 3 個 byte 就是選項數，
// 四個可調設定分別是 2／5／5／5（畫面模式／音效／戰略速度／戰術速度），
// 而 `sub_16062` 就是「點一下 +1，到頂繞回 0」。
// 值存在 `ds:0CF8h` 起一列一個 byte：`+0` 畫面模式、`+1` 音效、
// `+2` 戰略速度、`+3` 戰術速度（docs/re/55 §4）。
//
// **檔位就是 g.speed／g.tacticalSpeed 的值**：0 ＝ 最高速、4 ＝ 最低速。
// 每一檔實際多快由 speed.Throttle 換算（docs/spec/34）。

// cycleSpeed 重現 `sub_16062`：點一下換下一檔，到頂繞回第 0 檔。
//
// ⚠ 原版只有「往下一檔」一個方向（`inc al / cmp al, ah / jb` 之後
// `xor al, al`）。remake 讓右鍵往回一檔，那是額外的便利，不是原版行為。
func (g *game) cycleSpeed(tactical bool, forward bool) {
	p := &g.speed
	if tactical {
		p = &g.tacticalSpeed
	}
	step := 1
	if !forward {
		step = speed.Levels - 1
	}
	*p = (clamp(*p, 0, speed.Levels-1) + step) % speed.Levels
}

// adjustSpeed 是 remake 額外的細調（＋／− 鍵），不是原版行為。
// delta > 0 ＝ 更快 ＝ 檔位往 0 走（原版的檔位越大越慢）。
func (g *game) adjustSpeed(tactical bool, delta int) {
	p := &g.speed
	if tactical {
		p = &g.tacticalSpeed
	}
	*p = clamp(*p-delta, 0, speed.Levels-1)
}

// dispatchSystemRow 處理點到系統選單某一列。left ＝ 左鍵。
//
// ⚠ **原版那六個 handler 沒讀**（docs/re/55 §4），這裡是照標籤字面
// 接的 remake 行為：兩個速度左鍵加右鍵減，「資料儲存」開既有的四槽視窗，
// 「遊戲結束」走既有的 ＹＥＳ／ＮＯ 離開確認。
func (g *game) dispatchSystemRow(row int, left bool) {
	switch row {
	case sysRowStrategySpeed:
		g.cycleSpeed(false, left)
	case sysRowTacticalSpeed:
		g.cycleSpeed(true, left)
	case sysRowSound:
		g.toggleSound()
	case sysRowLordCorps:
		// ⚠ 左右鍵都是 toggle：這一列只有兩個值，原版那套
		// 「左鍵下一檔／右鍵上一檔」在兩個值上沒有意義。
		g.lordCorps = !g.lordCorps
	case sysRowSave:
		g.beginSaveUI(saveWrite)
	case sysRowQuit:
		g.quitting, g.quitYes = true, false
	}
}

// hitTestSystemRow 回傳點到系統選單的第幾列。
func hitTestSystemRow(x, y int) (int, bool) {
	for k := 0; k < sysRows; k++ {
		if (image.Point{X: x, Y: y}).In(sysRowRect(k)) {
			return k, true
		}
	}
	return 0, false
}

// sysRowRect 是第 k 列的可點矩形（原版熱區 0x20+k）。
func sysRowRect(k int) image.Rectangle {
	if k < 0 || k >= sysRows {
		return image.Rectangle{}
	}
	y := sysValueY + k*sysRowStep
	return image.Rect(sysValueX, y, sysValueX+sysValueW, y+sysValueH)
}

// lordCorpsValue 是「主君編成」那一列的值。
//
// ⚠ 兩個值都是三個半形寬的全形字，與旁邊的速度檔位對齊
// （速度是「 最速 」這種前後帶半形空白的七 byte 字串）。
func lordCorpsValue(allowed bool) string {
	if allowed {
		return " 可 "
	}
	return "不可"
}

// soundValue 是系統選單第 3 列的值。
//
// **沒有音檔時顯示「未接入」而不是「關」**——那是兩件事：
// 「關」是玩家的選擇，「未接入」是玩家還沒跑過 `tools/bgm2ogg.sh`。
// 混成同一個字會讓缺口從畫面上消失（`docs/spec/29` §5）。
// soundValue 是系統選單「音　效」那一格的值。
//
// ⭐ **字串用原版的**：`ds:6010h` 的五個選項是 ＯＦＦ／TYPE 1／2／3／4
// （docs/re/55 §4）。remake 只有開與關，所以對到 TYPE 1 與 ＯＦＦ——
// TYPE 1 就是 DOS/V 的 OPL3 那一組（docs/re/57）。
//
// ⚠ 「未接入」是 **remake 才有的狀態**（跑在沒有音訊裝置的機器上，
// 或沒給 `-audio`）。原版永遠有 YNSOUND，所以它沒有這個選項。
func (g *game) soundValue() string {
	switch {
	case !g.sound.Available():
		return "未接入"
	case g.sound.Enabled():
		return "TYPE 1"
	default:
		return "ＯＦＦ"
	}
}

func (g *game) toggleSound() {
	if !g.sound.Available() {
		return
	}
	g.sound.SetEnabled(!g.sound.Enabled())
}

// drawSystemWindow 畫系統選單（docs/spec/13 §2.6）。
//
// 原版第 1 列與第 6 列的值是清單裡就寫死的「 ＯＫ 」，中間四列由程式填。
// remake 沿用這個形狀，把還沒接的功能寫成值——**不要因為沒做就把列拿掉**，
// 那會讓缺口從畫面上消失。
// sysValueLine 把設定值置中排進 48 px 的值格，**用半形空白補**。
//
// ⭐ 原版就是這樣排的：「 ＯＫ 」的兩個空白是**半形**（8 ＋ 16 ＋ 16 ＋ 8 ＝ 48），
// 「 普通 」同理。用全形空白會變成 64 px，字就往右擠出去
// （實機上量到 remake 的字比原版右 10 px，docs/playtest/39）。
//
// 值本身若已經滿 48 px（「１６色」「TYPE 1」）就原樣回傳。
func sysValueLine(s string) string {
	// ⚠ **先翻再補空白。** 補完之後那一串（`"  可  "`）不會出現在詞表裡，
	// 逐句查一定落空——系統選單的值就會停在中文（docs/spec/87 §3）。
	s = uiLang.Convert(s)
	pad := (sysValueW - textdraw.StringWidth(s)) / 2 / textdraw.HalfW
	if pad <= 0 {
		return s
	}
	return strings.Repeat(" ", pad) + s + strings.Repeat(" ", pad)
}

func (g *game) drawSystemWindow(dst *ebiten.Image) {
	g.chrome.Window(dst, sysWinX, sysWinY, sysWinW, sysWinH, chrome.Menu)
	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	labelInk := g.paletteInk(strategyInkDim, color.RGBA{255, 223, 154, 255})

	g.td.Draw(dst, "　系　 統　 選　 單　", sysTitleX, sysTitleY, ink)
	vector.DrawFilledRect(dst, sysRuleX, sysRuleY, sysRuleW, 1, ink, false)

	// 第 2 列的兩個選項是「１６色」與「 液晶 」（原版字串表 ds:6002h，
	// 對應 GAMEPAL 的 bank 0–3 與 4–7，docs/re/55 §4）。
	// remake 只做了 1６色那一組，所以這一格是固定值。
	values := [sysRows]string{"ＯＫ", videoModeLabels[0], g.soundValue(),
		speed.Labels[clamp(g.speed, 0, speed.Levels-1)],
		speed.Labels[clamp(g.tacticalSpeed, 0, speed.Levels-1)],
		"ＯＫ", lordCorpsValue(g.lordCorps)}
	for k := 0; k < sysRows; k++ {
		dy := k * sysRowStep
		vector.DrawFilledRect(dst, sysLabelBoxX, float32(sysLabelBoxY+dy),
			sysLabelBoxW, sysLabelBoxH, color.Black, false)
		g.td.Draw(dst, sysMenuLabels[k], sysLabelX, sysLabelY+dy, labelInk)

		// 值格只有兩圈框；底色由值自己帶（原版的「 ＯＫ 」是黑字綠底）。
		g.dlValueBox(dst, sysValueX, sysValueY+dy, sysValueW, sysValueH)
		g.dlFill(dst, sysValueX, sysValueY+dy, sysValueW, sysValueH,
			dlValueFill, color.RGBA{45, 110, 55, 255})
		col := g.paletteInk(dlValueText, dlTextFallback)
		if values[k] == "未接入" {
			col = color.RGBA{170, 170, 170, 255}
		}
		g.td.Draw(dst, sysValueLine(values[k]), sysValueX, sysValueY+dy, col)
	}
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
	g.drawMinimapMarkers(screen)
	g.drawMinimapViewBox(screen)
	// 勢力色標是**原版的一張 192×16 圖**（段 3 0x09A0，`sub_15A3A` 貼在
	// (440,168)）：左半紅、右半藍，各帶一個小色塊。圖裡沒有君主名，
	// 那一層由 state 填（原版是 `sub_15DBB`，docs/re/62 §4.1）。
	legendX := strategySidebarX + chrome.Tile
	if img, err := g.lib.DOSVFactionLegend(int(g.world.Clock.Season())); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(legendX), float64(strategyMinimapLegendY))
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	// 兩個君主名的座標是原版數值：(480, 168) 與 (576, 168)（docs/re/62 §4.1）。
	g.td.Draw(screen, strategyHUDSingleLine(big5(g.world.LordName(g.world.Player)), 72),
		strategyLegendSelfX, strategyMinimapLegendY, chrome.Paper)
	watched := ""
	if f := g.watchedFaction(); f >= 0 {
		watched = big5(g.world.LordName(f))
	}
	g.td.Draw(screen, strategyHUDSingleLine(watched, 72),
		strategyLegendWatchedX, strategyMinimapLegendY, chrome.Paper)
}

// watchedFaction 是圖例第二格盯著的勢力（原版的 `cs:byte_198A7`）。
//
// ⭐ **開局時它是 0，而且可以等於自勢力**——原版沒有初始化那個 byte，
// 資料段開機是 0，所以第一次進主畫面時圖例右格顯示的是**勢力 0**
// （曹操局的實機截圖上兩格都是「曹操」）。**選單擋自勢力，初值不擋**
// （`sub_15AFC` 的 `cmp al, cs:byte_10CFF / jz 忽略`，docs/re/62 §4.2）。
//
// 已滅亡的勢力往後找；全都不行才回 −1。
func (g *game) watchedFaction() int {
	n := len(g.world.Factions)
	for i := 0; i < n; i++ {
		f := (g.minimapFaction + i) % n
		if g.world.Factions[f].Alive {
			return f
		}
	}
	return -1
}

// minimapMarkerColours 挑一個據點的（外框, 中心）色。
// 四種所屬的屬性出自 `sub_15CE0`（docs/re/62 §2）：
// 無所屬 0x0F、自勢力 0xAC、盯著的勢力 0xF3、其餘 0x83。
func (g *game) minimapMarkerColours(owner int) (border, centre color.RGBA) {
	switch {
	case owner < 0 || owner >= len(g.world.Factions):
		return g.minimapInk[0], g.minimapInk[1] // 0x0F
	case owner == g.world.Player:
		return g.minimapInk[2], g.minimapInk[3] // 0xAC
	case owner == g.watchedFaction():
		return g.minimapInk[1], g.minimapInk[4] // 0xF3
	default:
		return g.minimapInk[5], g.minimapInk[4] // 0x83
	}
}

// minimapMarkerPos 把格座標換算成標記左上角。原版是
// `dx = [si+8] >> 1 + 1B6h`／`bx = [si+0Ah] >> 1 + 26h`，
// 那個 −2 是把 4×4 的方塊對準格子中心（docs/re/62 §2）。
func minimapMarkerPos(cx, cy int) (int, int) {
	return strategyMinimapX + cx/2 - 2, strategyMinimapY + cy/2 - 2
}

// drawMinimapMarkers 畫 192 個據點（原版 `sub_15CC6` → `sub_15CE0`）。
// 圖形是 4×4：外框一圈背景色、中心 2×2 前景色（`sub_15D19`）。
// 座標是 `原點 + 格/2 − 2`，那個 −2 就是把 4×4 置中。
func (g *game) drawMinimapMarkers(dst *ebiten.Image) {
	for i := range g.world.Cities {
		c := &g.world.Cities[i]
		x, y := minimapMarkerPos(c.X, c.Y)
		if x < strategyMinimapX-2 || y < strategyMinimapY-2 ||
			x > strategyMinimapX+strategyMinimapW || y > strategyMinimapY+strategyMinimapImageH {
			continue
		}
		border, centre := g.minimapMarkerColours(c.Owner)
		vector.DrawFilledRect(dst, float32(x), float32(y), 4, 4, border, false)
		vector.DrawFilledRect(dst, float32(x+1), float32(y+1), 2, 2, centre, false)
	}
}

// drawMinimapViewBox 畫視野框（原版 `sub_15C58` → `sub_196ED`）。
// 大小 ＝ 畫面格數的一半：40×23 格 → 20×12 px。
//
// ⚠ 原版那是一張美術（`word_10D4C`），remake 畫線——**remake 差異**。
// drawMinimapViewBox 畫縮小地圖上的視野框。
//
// ⭐ **框是原版的點陣**（`ICONGRF` 段 3 `+0x8F0`，24×11 的 4 bpp 小圖），
// 不是描一個矩形——它有白邊 ＋ 右下黑影的立體感（docs/spec/55）。
//
// ⚠ **X 要減 4 格。** 原版的鏡頭變數 `word_1988E` 比畫面上實際看到的
// 那一欄小 4（docs/spec/55 §2），而 remake 的 camX 存的是**看到的那一欄**。
func (g *game) drawMinimapViewBox(dst *ebiten.Image) {
	pix, err := g.lib.ViewBox()
	if err != nil {
		return
	}
	x := strategyMinimapX + (g.camX-minimapCamBias)/2
	y := strategyMinimapY + g.camY/2
	season := int(g.world.Clock.Season())
	for dy := 0; dy < gfx.ViewBoxRows; dy++ {
		for dx := 0; dx < gfx.ViewBoxWidth; dx++ {
			c := pix[dy*gfx.ViewBoxWidth+dx]
			if c == gfx.ViewBoxTransparent {
				continue
			}
			col, err := g.lib.PaletteColor(season, int(c))
			if err != nil {
				continue
			}
			dst.Set(x+dx, y+dy, col)
		}
	}
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

	// 顏色也是原版數值，不再自己挑：顯示清單 op 08 的屬性 byte 與繪製
	// 常式的 `ax`／`bx` 高 byte 都是調色盤索引（docs/re/48 §4）。
	//
	//	0x0F  君主／首都／軍師的標籤與名字、資金與預備兵的數字
	//	0x09  「資金」「預備兵」兩個標籤
	//	0x0A  信賴度量條，以及資金為負時的數字
	ink := g.paletteInk(strategyInkNormal, chrome.Paper)
	labelInk := g.paletteInk(strategyInkDim, color.RGBA{255, 223, 154, 255})
	gaugeInk := g.paletteInk(strategyInkGauge, color.RGBA{85, 154, 69, 255})
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
	vector.DrawFilledRect(dst, float32(x+strategyInfoDividerXOffset), float32(infoY),
		1, strategyInfoDividerH, ink, false)
	// 信賴度：原版先畫一個 176×10 的槽（顯示清單 op 03），再在裡面畫量條
	// （`sub_15F27`，長度 (信賴度×100 + 0x9F) ÷ 0xA0，滿長 160）。
	// 量條本體高 2 px，填色 0x0A、未填色 0x00——`sub_10AAA` 分兩段畫，
	// 第二段把 `ah` 換成 `al`（0x00）。
	// **槽是純黑（色 0）**，不是深色調——實機上量到的（docs/spec/54 §2 同源）。
	g.td.Draw(dst, "信賴度", x+strategyTrustLabelX, y+strategyTrustLabelY, ink)
	vector.DrawFilledRect(dst, float32(x+strategyTrustSlotX), float32(y+strategyTrustSlotY),
		strategyTrustSlotW, strategyTrustSlotH, chrome.Blank, false)
	trustW := (g.world.Trust*100 + 0x9F) / 0xA0
	if trustW > strategyTrustMaxW {
		trustW = strategyTrustMaxW
	}
	if trustW > 0 {
		vector.DrawFilledRect(dst, float32(x+strategyTrustXOffset), float32(y+strategyTrustYOffset),
			float32(trustW), strategyTrustBarH, gaugeInk, false)
	}

	// 資源區的黑底也是顯示清單畫的（op 03，(448,304) 176×80）。
	vector.DrawFilledRect(dst, float32(x+strategyResourceBoxX), float32(y+strategyResourceBoxY),
		strategyResourceBoxW, strategyResourceBoxH, color.Black, false)
	fundsY := y + strategyFundsYOffset
	g.td.Draw(dst, "資金", x+strategyResourceLabelX, fundsY, labelInk)
	g.td.Draw(dst, "預備兵", x+strategyResourceLabelX, y+strategyReserveYOffset, labelInk)
	// 資金為負時原版把顏色從 0x0F 換成 0x0A（`sub_15F5D` 的 `cmp dh, 80h`）。
	fundsInk := ink
	if f.Funds < 0 {
		fundsInk = gaugeInk
	}
	g.drawOriginalNumber(dst, f.Funds, x+strategyFundsXOffset, fundsY,
		strategyFundsDigits, fundsInk)

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
		g.drawOriginalNumber(dst, value*strategyReserveMenPerPoint,
			x+strategyReserveXOffset, iconY, strategyReserveDigits, ink)
	}
}

// cityMarks 把世界狀態換成大地圖上要改的據點格（docs/spec/53）。
//
// ⭐ **中心格在據點記錄座標的右邊四格**（`world.CityCentreDX`）。
// 原版 `MMAP.MAP` 的 192 座據點逐座驗過，所以這裡不做例外處理；
// 真的對不上時 `RenderMarked` 會跳過那一座，並由 `checkCityCentres`
// 在載入時就記一筆 warning——**安靜畫錯比畫不出來難發現**。
func (g *game) cityMarks() []world.CityMark {
	if g.world == nil {
		return nil
	}
	marks := make([]world.CityMark, 0, len(g.world.Cities))
	capital := -1
	if p := g.world.Player; p >= 0 && p < len(g.world.Factions) {
		capital = g.world.Factions[p].Capital
	}
	for i := range g.world.Cities {
		c := &g.world.Cities[i]
		marks = append(marks, world.CityMark{
			X:       c.X + world.CityCentreDX,
			Y:       c.Y,
			Own:     world.OwnershipOf(c.Owner, g.world.Player),
			Capital: i == capital,
		})
	}
	return marks
}

// corpsMarks 回傳要疊在大地圖上的軍團圖塊（docs/spec/74）。
//
// ⚠ 原版是掃整張軍團表、`[si+0] >= 0xC0` 才畫；remake 的對應欄位是
// `Corps.Alive`（`internal/state/corps.go`），不要另外發明存在判定。
func (g *game) corpsMarks() []world.CorpsMark {
	if g.world == nil {
		return nil
	}
	marks := make([]world.CorpsMark, 0, 8)
	for i := range g.world.Corps {
		c := &g.world.Corps[i]
		if !c.Alive {
			continue
		}
		marks = append(marks, world.CorpsMark{
			X: c.X, Y: c.Y, Tile: world.CorpsTile(c.Faction, c.Heading),
		})
	}
	return marks
}

// checkCityCentres 在載入後驗一次「記錄座標 +4 是據點中心」。
//
// 這是 docs/spec/53 §5 的假設，**假設要有現形的機制**：
// 對不上就在 log 留一行，而不是讓畫面上少幾個徽記。
func (g *game) checkCityCentres() {
	if g.world == nil || g.lib == nil || g.lib.World == nil {
		return
	}
	bad := 0
	for i := range g.world.Cities {
		c := &g.world.Cities[i]
		t, err := g.lib.World.Tile(c.X+world.CityCentreDX, c.Y)
		if err != nil || !world.IsCityCentre(t) {
			bad++
		}
	}
	if bad > 0 {
		log.Printf("⚠ %d 座據點的 (X+4, Y) 不是據點中心圖塊，徽記會少畫（docs/spec/53 §5）", bad)
	}
}

// drawOriginalNumber 用**原版的數字字模**把一個數字右對齊畫進 digits 個格子。
//
// ⭐ 字模在 `ICONGRF` 段 3 的 `+0x840`，8×16、11 格（0–9 ＋ 負號），
// 規格 docs/spec/52 §4。倚天的 ASCII 數字墨水只有 9 列、字形也不同，
// **同樣的位置同樣的顏色還是逐像素差得出來**，所以數字不共用文字層。
//
// leftX 是**最左一格**的左緣（原版 `sub_1062F` 的 `di` 就是這一格），
// topY 是格子頂端——墨水落在 topY+1 … topY+14。
//
// 原版把前導的空位用**背景色**填滿；這裡只畫墨水，因為目前用到的地方
// 底本來就是那個背景色。哪天底不是純色了要回來補這一段。
func (g *game) drawOriginalNumber(dst *ebiten.Image, value, leftX, topY, digits int,
	ink color.RGBA) {
	if g.lib == nil || digits <= 0 {
		return
	}
	text := fmt.Sprintf("%d", value)
	if len(text) > digits {
		text = text[len(text)-digits:] // 溢位保留低位，與原版的位數上限一致
	}
	x := leftX + (digits-len(text))*gfx.DigitWidth
	for _, r := range text {
		idx := int(r - '0')
		if r == '-' {
			idx = gfx.DigitMinus
		} else if r < '0' || r > '9' {
			continue
		}
		mask, err := g.lib.DigitMask(idx)
		if err != nil {
			return
		}
		for dy := 0; dy < gfx.DigitHeight; dy++ {
			for dx := 0; dx < gfx.DigitWidth; dx++ {
				if mask[dy*gfx.DigitWidth+dx] != 0 {
					dst.Set(x+dx, topY+dy, ink)
				}
			}
		}
		x += gfx.DigitWidth
	}
}

// drawBannerNumber 把日期填進橫幅。right 是數字欄的右緣。
func (g *game) drawBannerNumber(screen *ebiten.Image, value, right int, ink color.RGBA) {
	digits := len(fmt.Sprintf("%d", value))
	g.drawOriginalNumber(screen, value, right-digits*gfx.DigitWidth, bannerTextY, digits, ink)
}
