package main

import "github.com/wicanr2/wolong_cht/internal/assets/cjk"

// 松崗 DOS/V 戰術畫面的固定版面。
//
// 原版戰術是 320×200 邏輯畫布；本 remake 的正式畫布是 640×400，
// 所以這裡只保留一份「原版邏輯像素 × 2」的幾何契約。戰場投影本身仍
// 使用 docs/re/11 §5.12 的 8 px 原生欄列；這個檔案只決定它在完整畫面
// 的位置，以及 TALK／右欄／底列的可見邊界。
//
// DOS/V 量測基準（confirmed／layout evidence）：
//   - 320×200 邏輯畫布；
//   - 戰場左欄 240×184（30×23 個 8 px 欄列，底部 16 px 為命令列）；
//   - 右欄 80×200；
//   - remake 以 2 倍整數放大到 640×400。
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

	SideAttacker battleRect
	SideMiniMap  battleRect
	SideDefender battleRect
	SideCommands battleRect
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

var battleCommandLabels = [...]string{"陣形", "攻擊", "突擊", "城壁", "守陣", "退卻"}

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

type battleStatusColumns struct {
	HP, Advantage battleRect
}

// battleSideStatusColumns 固定第三列的左右安全欄。HP 欄可容納
// 「體力 100」64 px，優劣欄保留兩個全形字 32 px，中間至少 8 px。
func battleSideStatusColumns(r battleRect) battleStatusColumns {
	innerX := r.X + battlePanelInset + battleStatusTextPad
	innerRight := r.right() - battlePanelInset - battleStatusTextPad
	const advantageW = 32
	advX := innerRight - advantageW
	return battleStatusColumns{
		HP: battleRect{X: innerX, Y: r.Y, W: advX - battleCommandMinPad - innerX,
			H: r.H},
		Advantage: battleRect{X: advX, Y: r.Y, W: advantageW, H: r.H},
	}
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

		SideAttacker: battleRect{X: sideX + 8, Y: 8, W: sideW - 16, H: 64},
		// sub_1C51E 的 DOS/V 原點是 x=0x1F0、y=0x50；在 640×400
		// 畫布上這裡保留 x=496、y=80 的完整 128×128 raw image。
		// 它不再被 8 px 面板內縮吞掉，右側剩餘 16 px 作為 sidebar 邊界。
		SideMiniMap:  battleRect{X: sideX + 16, Y: 80, W: 128, H: 128},
		SideDefender: battleRect{X: sideX + 8, Y: 216, W: sideW - 16, H: 64},
		// sub_1C863：segment1+0x1800，AX=0x6008，目的 (496,280)。
		SideCommands: battleRect{X: sideX + 16, Y: 280, W: 128, H: 96},
	}
}
