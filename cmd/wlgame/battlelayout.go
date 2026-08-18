package main

import (
	"github.com/wicanr2/wolong_cht/internal/assets/cjk"
	"github.com/wicanr2/wolong_cht/internal/assets/gfx"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/ui/textdraw"
)

// 松崗 DOS/V 戰術畫面的固定版面。
//
// **原版的戰術畫面就是 640×400**：sub_1F1A3 的裁切邊界是 0x280／0x190，
// 顯示記憶體一列 80 bytes（640 ÷ 8），sub_1C863 的貼點換算回來也全部
// 落在 640×400 上（docs/re/60 §0）。所以 remake 與原版是 1:1，
// 下面 dosvBattleNative* 那一組「原生 × 2」的常數只是同一組數字的
// 另一種寫法，不代表原版是 320×200。
//
// 戰場投影本身使用 docs/re/11 §5.12 的 8 px 原生欄列；這個檔案只決定
// 它在完整畫面的位置，以及 TALK／右欄／底列的可見邊界。
//
// DOS/V 量測基準（confirmed）：
//   - 畫布 640×400；
//   - 戰場左欄 480×368，底部 32 px 是六個部隊 slot（sub_1C7F4，y=368）；
//   - 側欄 x 480–640，左右柱各 16 px，內容寬 128（sub_1CA3B）。
//
// TALK 會覆蓋戰場 viewport，這是原版畫面層級，而不是 base region
// 重疊錯誤；contract 只驗證 field／sidebar／bottom 三個 base region。

type battleRect struct {
	X, Y, W, H int
}

func (r battleRect) right() int  { return r.X + r.W }
func (r battleRect) bottom() int { return r.Y + r.H }

func (r battleRect) contains(other battleRect) bool {
	return other.X >= r.X && other.Y >= r.Y &&
		other.right() <= r.right() && other.bottom() <= r.bottom()
}

func (r battleRect) overlaps(other battleRect) bool {
	return r.X < other.right() && other.X < r.right() &&
		r.Y < other.bottom() && other.Y < r.bottom()
}

// dosvBattleLayout 是 640×400 戰術畫面的唯一幾何來源。
type dosvBattleLayout struct {
	Field          battleRect
	Sidebar        battleRect
	BottomCommands battleRect

	// overlay regions：原版 TALK 疊在戰場上，不算 base region。
	TopTalk    battleRect
	BottomTalk battleRect

	// 側欄七格，出處 docs/spec/31 §2.1（sub_1C7A9 那一串）。
	// ⭐ 上格是**對方**、下格是**我方**，與攻守無關：原版 sub_14ED7
	// 依玩家攻守互換 word_10D2E／word_10D30，側欄固定把 word_10D2E
	// （我方）畫在下面。
	SideTitle     battleRect
	SideFoe       battleRect
	SideMiniMap   battleRect
	SideAlly      battleRect
	SideFormation battleRect
	SideCommands  battleRect
	// SideLines 是陣形線那三格（熱區 0x04–0x06，docs/spec/37 §2.2）。
	SideLines [3]battleRect
	SideFooter    battleRect
}

// battleTalkState 是呈現層的最小 TALK contract；文案來自 TALK.DAT，
// 不在 state 儲存字串。
type battleTalkState struct {
	Top, Bottom                 string
	TopPortrait, BottomPortrait int
}

func (s battleTalkState) text(side int) string {
	if side == 0 {
		return s.Top
	}
	return s.Bottom
}

func (s battleTalkState) visible(side int) bool { return s.text(side) != "" }

func (s battleTalkState) portrait(side int) int {
	if side == 0 {
		return s.TopPortrait
	}
	return s.BottomPortrait
}

// battleCommandLabels 以**命令碼**為索引，不是畫面列序。
var battleCommandLabels = [...]string{"陣形", "攻擊", "突擊", "城壁", "守陣", "退却"}

// battleSideCommandRowCode 是指令面板由上而下第 row 列送出的命令碼。
//
// 出處 docs/re/60 §6.1：sub_1C863 在 (496, 280+16k) 註冊的熱區碼依序是
// 0x09／0x08／0x07／0x0A／0x0B／0x0C，而 handler 0x1C1B9 算的是
// `命令碼 = 熱區碼 − 7`。**畫面順序不是命令碼順序**——第 0 列送命令 2。
// 台詞把三個名字釘死了（TALK 0x1B1+命令碼，×8 展開後 622／630／638）：
// 命令 0「擺出陣形！！」、命令 1「前進！轉為攻擊！」、命令 2「好啊！衝啊！！！」。
var battleSideCommandRowCode = [...]int{2, 1, 0, 3, 4, 5}

// battleSideCommandRowLabel 回傳指令面板第 row 列該顯示的字。
func battleSideCommandRowLabel(row int) string {
	if row < 0 || row >= len(battleSideCommandRowCode) {
		return ""
	}
	return battleCommandLabels[battleSideCommandRowCode[row]]
}

// battleBottomSlotSquad 是底列由左到右第 i 格對應哪一個編成位置
// （原版 `cs:0xD2E4`，docs/spec/33 §1.1）。
//
// 六個編成位置是 0 主將／1 前鋒／2 左翼／3 右翼／4 左備／5 右備，
// 所以畫面上是「左翼 左備 主將 前鋒 右備 右翼」——**空間排列**。
var battleBottomSlotSquad = [...]int{2, 4, 0, 1, 5, 3}

// battleSquadSlotX 是第 k 個編成位置畫在哪個 X（原版 `cs:0xD2EA`）。
// 它與 battleBottomSlotSquad 互為反排列。
var battleSquadSlotX = [...]int{160, 240, 0, 400, 80, 320}

// battleSquadSlot 回傳第 k 個編成位置在底列的第幾格。
func battleSquadSlot(squad int) int {
	if squad < 0 || squad >= len(battleSquadSlotX) {
		return -1
	}
	return battleSquadSlotX[squad] / battleBottomSlotW
}

const battleBottomSlotW = 80

// 底列每一格內的四樣東西（docs/spec/33 §1.2）。
// 三張圖示都是 24×16，橫向接續填滿 80 px 的格子。
const (
	battleSlotGlyphX = 4 // 位置名 glyph
	battleSlotGlyphY = 6
	battleSlotArmX   = 29 // 兵種圖示（sub_19B6D 的 `add dx, 1Dh`）
	battleSlotArmY   = 6
	battleSlotOrderX = 54 // 目前命令的圖示（sub_1C673 的 `add dx, 36h`）
	battleSlotOrderY = 6
	battleSlotBarX   = 2 // 待機兵條
	battleSlotBarY   = 396 // sub_1C74C 的 bx=0x18C
	battleSlotBarLen = 0x4C
	battleSlotBarH   = 2
)

// battleSlotSelectRects 是 sub_1C6BF 的兩個同心矩形，
// 相對於格左緣 x0：外框 (x0+2,372)-(x0+77,392)、內框各縮 1 px。
func battleSlotSelectRects(x0 int) (outer, inner battleRect) {
	outer = battleRect{X: x0 + 2, Y: 372, W: 77 - 2 + 1, H: 392 - 372 + 1}
	inner = battleRect{X: x0 + 3, Y: 373, W: 76 - 3 + 1, H: 391 - 373 + 1}
	return outer, inner
}

// battleSideCommandRowOf 是 battleSideCommandRowCode 的反查：命令碼在第幾列。
func battleSideCommandRowOf(code int) int {
	for row, c := range battleSideCommandRowCode {
		if c == code {
			return row
		}
	}
	return -1
}

const (
	battleCommandMinPad = 8
	battlePanelInset    = 8
	battleStatusTextPad = 4
)

func battleCommandLabelFits(label string, cellW int) bool {
	return battleCommandTextWidth(label)+2*battleCommandMinPad <= cellW
}

// battleCommandTextWidth mirrors textdraw.RuneWidth without importing Ebiten's
// Drawer into the geometry-only test. This is the same fixed DOS/V contract:
// ASCII 8 px, CJK／全形 16 px。
func battleCommandTextWidth(s string) int {
	w := 0
	for _, ch := range s {
		if ch < 0x80 {
			w += cjk.GlyphWidth / 2
			continue
		}
		w += cjk.GlyphWidth
	}
	return w
}

// splitBattleCommandCells 將底列的 480 px 寬度完整分配給六格，
// 每格正好 80 px；最後一格的右緣與右欄框線精確接齊。
func splitBattleCommandCells(r battleRect) []battleRect {
	n := len(battleCommandLabels)
	base, extra := r.W/n, r.W%n
	out := make([]battleRect, n)
	x := r.X
	for i := range out {
		w := base
		if i < extra {
			w++
		}
		out[i] = battleRect{X: x, Y: r.Y, W: w, H: r.H}
		x += w
	}
	return out
}

// battleSideCommandCells 對應 sub_1C863 貼在 (496,280) 的 128×96
// 原版複合面板：六個命令是單欄六列，每列 16 px。
func battleSideCommandCells(r battleRect) []battleRect {
	out := make([]battleRect, len(battleCommandLabels))
	for row := range out {
		y0 := r.Y + row*r.H/len(out)
		y1 := r.Y + (row+1)*r.H/len(out)
		out[row] = battleRect{X: r.X, Y: y0, W: r.W, H: y1 - y0}
	}
	return out
}

// battleHitRect 是可點擊的 glyph 內框。面板外框與每格之間的安全邊界
// 不屬於命中區，避免點到 chrome 邊線或分隔空隙時誤下令。
func battleHitRect(r battleRect) (battleRect, bool) {
	const inset = 2
	if r.W <= 2*inset || r.H <= 2*inset {
		return battleRect{}, false
	}
	return battleRect{X: r.X + inset, Y: r.Y + inset,
		W: r.W - 2*inset, H: r.H - 2*inset}, true
}

func (r battleRect) containsPoint(x, y int) bool {
	return x >= r.X && x < r.right() && y >= r.Y && y < r.bottom()
}

func battleCommandIndexAt(cells []battleRect, x, y int) (int, bool) {
	for i, cell := range cells {
		hit, ok := battleHitRect(cell)
		if ok && hit.containsPoint(x, y) {
			return i, true
		}
	}
	return 0, false
}

// splitBattleCommandIndexAt 是底列六格的純 hit-test；它不讀取 game 狀態，
// 也不執行任何規則。繪圖與輸入都以 splitBattleCommandCells 的同一格座標為準。
func splitBattleCommandIndexAt(r battleRect, x, y int) (int, bool) {
	return battleCommandIndexAt(splitBattleCommandCells(r), x, y)
}

// battleSideCommandIndexAt 是右欄 2×3 命令格的純 hit-test；面板外框與
// 內縮空白不命中，命令索引仍固定為 0..5、逐列由左至右排列。
func battleSideCommandIndexAt(r battleRect, x, y int) (int, bool) {
	return battleCommandIndexAt(battleSideCommandCells(r), x, y)
}

// battleChoiceLayout 集中遭遇選單視窗與兩列的可見文字區。列區只涵蓋
// glyph 高度，列距中的空白與視窗外框不命中。
type battleChoiceLayout struct {
	Window battleRect
	Rows   [2]battleRect
}

func battleChoiceLayoutFor() battleChoiceLayout {
	const (
		x, y, w, h = 104, 112, 432, 160
		rowX       = x + 8 + 4
		rowY       = y + 3*15 + 24
		rowW       = w - 2*(8+4)
		rowH       = 15
		rowPitch   = 15 + 4
	)
	return battleChoiceLayout{
		Window: battleRect{X: x, Y: y, W: w, H: h},
		Rows: [2]battleRect{
			{X: rowX, Y: rowY, W: rowW, H: rowH},
			{X: rowX, Y: rowY + rowPitch, W: rowW, H: rowH},
		},
	}
}

func battleChoiceRowAt(x, y int) (int, bool) {
	rows := battleChoiceLayoutFor().Rows
	for i, row := range rows {
		if row.containsPoint(x, y) {
			return i, true
		}
	}
	return 0, false
}

// battleSideCellLayout 是一格將旗裡三件東西的位置。
//
// 出處 docs/re/60 §3／§5：主將名走 sub_106FD（固定三個全形字），
// 兩條計量條走 sub_10AAA（x=498、高 2、總長上限 124）。
// ⚠ **上下兩格的排列是相反的**：上格是「名在上、條在下」（名 y=52、
// 條 y=72／75），下格是「條在上、名在下」（條 y=211／214、名 y=221）。
// 照抄，不要對稱化。
type battleSideCellLayout struct {
	Name      battleRect // 主將名的左上角，寬 48（三個全形字）
	MenBar    battleRect // 兵力條，色 12
	HealthBar battleRect // 大將體力條，色 11
}

const (
	battleSideBarX      = 2   // 498 − 496
	battleSideBarMaxLen = 124 // sub_1C775／sub_1C78E 的 CH=0x7C
	battleSideBarH      = 2   // sub_10AAA 的 CH=2
	battleSideNameX     = 32  // 528 − 496
	battleSideNameW     = 48  // 三個全形字
)

func battleSideCellLayoutFor(r battleRect, top bool) battleSideCellLayout {
	nameDY, menDY, healthDY := 13, 3, 6 // 下格：221／211／214
	if top {
		nameDY, menDY, healthDY = 4, 24, 27 // 上格：52／72／75
	}
	bar := func(dy int) battleRect {
		return battleRect{X: r.X + battleSideBarX, Y: r.Y + dy,
			W: battleSideBarMaxLen, H: battleSideBarH}
	}
	return battleSideCellLayout{
		Name: battleRect{X: r.X + battleSideNameX, Y: r.Y + nameDY,
			W: battleSideNameW, H: textdraw.GlyphH},
		MenBar:    bar(menDY),
		HealthBar: bar(healthDY),
	}
}

// battleSideBarLengths 照抄 sub_1C6F6 的兩條算式：
// 兵力 >> 2、大將體力 × 3 ÷ 4，各自上限 124。
func battleSideBarLengths(men, health int) (menLen, healthLen int) {
	clamp := func(v int) int {
		if v < 0 {
			return 0
		}
		if v > battleSideBarMaxLen {
			return battleSideBarMaxLen
		}
		return v
	}
	return clamp(men >> 2), clamp(health * 3 / 4)
}

const (
	dosvBattleScreenW = 640
	dosvBattleScreenH = 400
	dosvBattleScale   = 2

	// 正規化松崗 DOS/V 錄影後，右欄左界在 logical x=240；先前的
	// 248 是把縮圖自身原點誤當成右欄起點。
	dosvBattleNativeFieldW   = 240
	dosvBattleNativeFieldH   = 184
	dosvBattleNativeSideW    = 80
	dosvBattleNativeCommandH = 16
)

func dosvBattleLayoutFor(w, h int) dosvBattleLayout {
	// 目前遊戲的邏輯畫布固定為 640×400。保留 w／h 參數是讓 contract
	// 明確檢查呼叫端，不讓未來縮放時偷偷產生第二套魔術數字。
	if w != dosvBattleScreenW || h != dosvBattleScreenH {
		w, h = dosvBattleScreenW, dosvBattleScreenH
	}
	fieldW := dosvBattleNativeFieldW * dosvBattleScale
	fieldH := dosvBattleNativeFieldH * dosvBattleScale
	sideW := dosvBattleNativeSideW * dosvBattleScale
	commandH := dosvBattleNativeCommandH * dosvBattleScale
	sideX := w - sideW
	return dosvBattleLayout{
		Field:          battleRect{X: 0, Y: 0, W: fieldW, H: fieldH},
		Sidebar:        battleRect{X: sideX, Y: 0, W: sideW, H: h},
		BottomCommands: battleRect{X: 0, Y: fieldH, W: fieldW, H: commandH},

		TopTalk:    battleRect{X: 0, Y: 0, W: 256, H: 80},
		BottomTalk: battleRect{X: 224, Y: 288, W: 256, H: 80},

		// 側欄的七格全部是 sub_1C7A9 那一串的直接座標（docs/re/60 §2）。
		// 內容寬一律 128 px，起點 x=496——側欄的左右柱子各佔 16 px。
		//
		// sub_1F1A3 在 (496,8)-(623,39) 清出標題格；兩行文字在 y=8 與 y=24。
		SideTitle: battleRect{X: sideX + 16, Y: 8, W: 128, H: 32},
		// sub_1C863：segment1+0x0800，AX=0x2008，目的 (496,48)。
		SideFoe: battleRect{X: sideX + 16, Y: 48, W: 128, H: 32},
		// sub_1C51E 的 DOS/V 原點是 x=0x1F0、y=0x50；在 640×400
		// 畫布上這裡保留 x=496、y=80 的完整 128×128 raw image。
		SideMiniMap: battleRect{X: sideX + 16, Y: 80, W: 128, H: 128},
		// sub_1C863：segment1+0x1000，AX=0x2008，目的 (496,208)。
		SideAlly: battleRect{X: sideX + 16, Y: 208, W: 128, H: 32},
		// sub_1C863：segment1+0x0000，AX=0x2008，目的 (496,248)。
		// 8 欄 × 2 列的 16×16 格 ＝ 十六個陣形（docs/re/60 §8）。
		SideFormation: battleRect{X: sideX + 16, Y: 248, W: 128, H: 32},
		// sub_1C863：segment1+0x1800，AX=0x6008，目的 (496,280)。
		SideCommands: battleRect{X: sideX + 16, Y: 280, W: 128, H: 96},
		SideLines: battleLineRects(sideX),
		// sub_1C863：segment1+0x3500，AX=0x1008，目的 (496,376)。
		SideFooter: battleRect{X: sideX + 16, Y: 376, W: 128, H: 16},
	}
}

const (
	// battleFormationCols／Rows：熱區 0x03 那一格是 8 欄 × 2 列的 16×16 格。
	battleFormationCols = 8
	battleFormationRows = 2
	battleFormationCell = 16
)

// battleFormationIndexAt 對應 handler 0x1C11A：
//
//	col = (游標X − 0x1F0) >> 4 ；(游標Y − 0xF8) >= 0x10 時再 +8
func battleFormationIndexAt(r battleRect, x, y int) (int, bool) {
	if !r.containsPoint(x, y) {
		return 0, false
	}
	col := (x - r.X) / battleFormationCell
	if col >= battleFormationCols {
		col = battleFormationCols - 1
	}
	if y-r.Y >= battleFormationCell {
		col += battleFormationCols
	}
	return col, true
}

// battleLineRects 是陣形線那三格（熱區 0x04／0x05／0x06，docs/re/60 §10）。
// 由上而下是**敵軍側／中央／自軍側**，寫進 word_1D33C 的值是 48／28／5。
func battleLineRects(sideX int) [3]battleRect {
	x := sideX + 72 // 原版 552 ＝ 側欄 480 + 72
	return [3]battleRect{
		{X: x, Y: 288, W: 64, H: 24},
		{X: x, Y: 312, W: 64, H: 32},
		{X: x, Y: 344, W: 64, H: 24},
	}
}

// battleLineIndexAt 回傳點到第幾格陣形線。
func battleLineIndexAt(rects [3]battleRect, x, y int) (int, bool) {
	for i, r := range rects {
		if r.containsPoint(x, y) {
			return i, true
		}
	}
	return 0, false
}

// battleFormationCellRect 是第 idx 格的 16×16 矩形；選取框再內縮 1 px
// 畫成 14×14（sub_1C61F）。
func battleFormationCellRect(r battleRect, idx int) battleRect {
	if idx < 0 || idx >= battleFormationCols*battleFormationRows {
		return battleRect{}
	}
	return battleRect{
		X: r.X + (idx%battleFormationCols)*battleFormationCell,
		Y: r.Y + (idx/battleFormationCols)*battleFormationCell,
		W: battleFormationCell,
		H: battleFormationCell,
	}
}

// battleSlotOrderIcon 把一隊隊長的命令換成底列要畫的圖示編號。
//
// ⭐ 原版的圖示是**下令當時畫一次**就不再更新（docs/spec/33 §1.6）：
// `sub_1A92E` 在全隊走到定位時把單位記錄改成 7（就位）卻不重畫，
// 所以格子裡留著的是「陣形」那一張。照著把目前的命令碼直接畫，
// 整列會在就位之後變空——7 不在 0–5 裡。
//
// 回傳 −1 表示這一格不畫圖示。
func battleSlotOrderIcon(cmd tactical.Command) int {
	if cmd == tactical.Holding {
		return int(tactical.Form)
	}
	if cmd < 0 || int(cmd) >= gfx.DOSVOrderIconCount {
		return -1
	}
	return int(cmd)
}

// battleSlotArmIcon 把一隊的兵種換成兵種圖示編號（0 馬／1 弓／2 步）。
// 原版存的是兵種 × 18，圖示索引是 `兵種 − 1`；大將（0）沒有圖示。
func battleSlotArmIcon(kind tactical.Kind) int {
	if kind <= tactical.General {
		return -1
	}
	idx := int(kind)/18 - 1
	if idx < 0 || idx >= gfx.DOSVSquadArmIconCount {
		return -1
	}
	return idx
}
