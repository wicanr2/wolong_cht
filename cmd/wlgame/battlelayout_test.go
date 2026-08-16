package main

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// 戰術外框在原版是模式無關的硬體畫布；攻城與兩軍遭遇只能交換戰場
// 內容、守軍來源與結算，不得各自產生第二套 sidebar／TALK／命中區。
func battleLayoutForMode(mode combat.Mode) dosvBattleLayout {
	switch mode {
	case combat.Field, combat.Siege:
		return dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	default:
		return dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	}
}

func TestFieldAndSiegeShareExactDOSVBattleChrome(t *testing.T) {
	field := battleLayoutForMode(combat.Field)
	siege := battleLayoutForMode(combat.Siege)
	if field != siege {
		t.Fatalf("兩軍遭遇與攻城不得使用不同戰術骨架：\nfield=%#v\nsiege=%#v", field, siege)
	}
	for name, cells := range map[string][]battleRect{
		"field bottom": splitBattleCommandCells(field.BottomCommands),
		"siege bottom": splitBattleCommandCells(siege.BottomCommands),
		"field side":   battleSideCommandCells(field.SideCommands),
		"siege side":   battleSideCommandCells(siege.SideCommands),
	} {
		if len(cells) != len(battleCommandLabels) {
			t.Fatalf("%s 命中格數=%d，want %d", name, len(cells), len(battleCommandLabels))
		}
	}
	if fmt.Sprint(splitBattleCommandCells(field.BottomCommands)) !=
		fmt.Sprint(splitBattleCommandCells(siege.BottomCommands)) ||
		fmt.Sprint(battleSideCommandCells(field.SideCommands)) !=
			fmt.Sprint(battleSideCommandCells(siege.SideCommands)) {
		t.Fatal("兩種模式的指令視覺格／命中格必須逐格相同")
	}
}

func TestDOSVBattleLayoutUsesMeasured640x400Regions(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)

	if got, want := l.Field, (battleRect{0, 0, 480, 368}); got != want {
		t.Fatalf("戰場區 = %#v，預期 %#v", got, want)
	}
	if got, want := l.Sidebar, (battleRect{480, 0, 160, 400}); got != want {
		t.Fatalf("右欄 = %#v，預期 %#v", got, want)
	}
	if got, want := l.BottomCommands, (battleRect{0, 368, 480, 32}); got != want {
		t.Fatalf("底列 = %#v，預期 %#v", got, want)
	}

	if l.Field.overlaps(l.Sidebar) || l.Field.overlaps(l.BottomCommands) ||
		l.Sidebar.overlaps(l.BottomCommands) {
		t.Fatal("base region 不得重疊")
	}
	for name, r := range map[string]battleRect{
		"field":           l.Field,
		"sidebar":         l.Sidebar,
		"bottom commands": l.BottomCommands,
		"top talk":        l.TopTalk,
		"bottom talk":     l.BottomTalk,
		"title":           l.SideTitle,
		"foe":             l.SideFoe,
		"mini map":        l.SideMiniMap,
		"ally":            l.SideAlly,
		"formation":       l.SideFormation,
		"side commands":   l.SideCommands,
		"footer":          l.SideFooter,
	} {
		if r.X < 0 || r.Y < 0 || r.W <= 0 || r.H <= 0 ||
			r.right() > dosvBattleScreenW || r.bottom() > dosvBattleScreenH {
			t.Errorf("%s = %#v 不在 640×400 內", name, r)
		}
	}

	if !l.Field.contains(l.TopTalk) || !l.Field.contains(l.BottomTalk) {
		t.Fatal("TALK overlay 必須留在戰場 viewport 內")
	}
	for _, r := range []battleRect{l.SideTitle, l.SideFoe, l.SideMiniMap, l.SideAlly,
		l.SideFormation, l.SideCommands, l.SideFooter} {
		if !l.Sidebar.contains(r) {
			t.Fatalf("右欄子區域 %#v 必須留在右欄內", r)
		}
	}
	// docs/re/60 §1.2／§2：七格全部是 sub_1C7A9 那一串的直接座標。
	// ⚠ 上格（對方）與小地圖之間**沒有間隙**——原版靠 sub_1CA3B 的
	// 橫帶分隔，不是靠 8 px 框距。
	for _, c := range []struct {
		name string
		got  battleRect
		want battleRect
	}{
		{"標題", l.SideTitle, battleRect{496, 8, 128, 32}},
		{"上格（對方）", l.SideFoe, battleRect{496, 48, 128, 32}},
		{"小地圖", l.SideMiniMap, battleRect{496, 80, 128, 128}},
		{"下格（我方）", l.SideAlly, battleRect{496, 208, 128, 32}},
		{"陣形列", l.SideFormation, battleRect{496, 248, 128, 32}},
		{"命令面板", l.SideCommands, battleRect{496, 280, 128, 96}},
		{"底列", l.SideFooter, battleRect{496, 376, 128, 16}},
	} {
		if c.got != c.want {
			t.Fatalf("%s = %#v，預期 sub_1C863 的直接座標 %#v", c.name, c.got, c.want)
		}
	}

	// 這是正規化原版的 240:80 邏輯比例，避免之後只為了填黑邊而任意放大右欄。
	if l.Field.W/l.Sidebar.W != 3 || l.Field.W%l.Sidebar.W != 0 {
		t.Fatalf("戰場／右欄寬度比例錯誤：%d:%d", l.Field.W, l.Sidebar.W)
	}
}

func TestDOSVBattleLayoutRejectsAlternativeCanvasWithoutChangingGeometry(t *testing.T) {
	got := dosvBattleLayoutFor(800, 600)
	want := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	if got != want {
		t.Fatal("戰術版面不應因未支援的畫布尺寸偷偷產生第二套幾何")
	}
}

func TestBattleTalkStateHidesEmptyPayload(t *testing.T) {
	var empty battleTalkState
	if empty.visible(0) || empty.visible(1) {
		t.Fatal("沒有 TALK payload 時不得顯示任何 overlay")
	}
	state := battleTalkState{Top: "受控訊息"}
	if !state.visible(0) || state.visible(1) || state.text(0) != "受控訊息" {
		t.Fatalf("TALK 可見狀態錯誤：%#v", state)
	}
}

func TestDOSVBattleCommandLabelsFit80PixelCells(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	cells := splitBattleCommandCells(l.BottomCommands)
	if len(cells) != len(battleCommandLabels) {
		t.Fatalf("底列格數 = %d，預期 %d", len(cells), len(battleCommandLabels))
	}
	for i, cell := range cells {
		wantW := 80
		if cell.W != wantW {
			t.Errorf("底列第 %d 格寬度 = %d，預期 %d", i+1, cell.W, wantW)
		}
		label := fmt.Sprintf("%d %s", i+1, battleCommandLabels[i])
		if !battleCommandLabelFits(label, cell.W) {
			t.Fatalf("底列命令 %q 超出 %d px cell", label, cell.W)
		}
	}
	if cells[0].X != l.BottomCommands.X || cells[len(cells)-1].right() != l.BottomCommands.right() {
		t.Fatal("底列六格沒有完整分配 480 px")
	}
}

// 上下兩格的排列是相反的（docs/re/60 §3／§5）：上格名 y=52、條 y=72／75，
// 下格條 y=211／214、名 y=221。這一項照抄，不要對稱化。
func TestDOSVBattleSideCellLayoutMatchesRawCoordinates(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	for _, c := range []struct {
		name                     string
		cell                     battleRect
		top                      bool
		nameY, menBarY, healthY  int
	}{
		{"上格（對方）", l.SideFoe, true, 52, 72, 75},
		{"下格（我方）", l.SideAlly, false, 221, 211, 214},
	} {
		got := battleSideCellLayoutFor(c.cell, c.top)
		if got.Name.X != 528 || got.Name.Y != c.nameY || got.Name.W != 48 {
			t.Errorf("%s 主將名 = %#v，預期 (528,%d) 寬 48", c.name, got.Name, c.nameY)
		}
		if got.MenBar.X != 498 || got.MenBar.Y != c.menBarY ||
			got.MenBar.W != 124 || got.MenBar.H != 2 {
			t.Errorf("%s 兵力條 = %#v，預期 (498,%d) 124×2", c.name, got.MenBar, c.menBarY)
		}
		if got.HealthBar.X != 498 || got.HealthBar.Y != c.healthY ||
			got.HealthBar.W != 124 || got.HealthBar.H != 2 {
			t.Errorf("%s 體力條 = %#v，預期 (498,%d) 124×2", c.name, got.HealthBar, c.healthY)
		}
	}
}

// sub_1C775 是 `值 >> 2`、sub_1C78E 是 `值 × 3 ÷ 4`，兩者上限都是 0x7C。
func TestDOSVBattleSideBarLengthsMatchRawFormulas(t *testing.T) {
	for _, c := range []struct{ men, health, wantMen, wantHealth int }{
		{0, 0, 0, 0},
		{100, 100, 25, 75},
		{400, 160, 100, 120},
		{9999, 9999, 124, 124}, // 兩條都夾在 0x7C
		{-5, -5, 0, 0},
	} {
		gotMen, gotHealth := battleSideBarLengths(c.men, c.health)
		if gotMen != c.wantMen || gotHealth != c.wantHealth {
			t.Errorf("battleSideBarLengths(%d,%d) = %d,%d，預期 %d,%d",
				c.men, c.health, gotMen, gotHealth, c.wantMen, c.wantHealth)
		}
	}
}

// docs/re/60 §6.1：sub_1C863 在 (496, 280+16k) 註冊的熱區碼是
// 0x09/0x08/0x07/0x0A/0x0B/0x0C，handler 0x1C1B9 算 `命令碼 = 熱區碼 − 7`。
func TestBattleSideCommandRowsMatchRawHotspotCodes(t *testing.T) {
	rawHotspots := [...]int{0x09, 0x08, 0x07, 0x0A, 0x0B, 0x0C}
	wantLabels := [...]string{"突擊", "攻擊", "陣形", "城壁", "守陣", "退却"}
	for row, hot := range rawHotspots {
		if got, want := battleSideCommandRowCode[row], hot-7; got != want {
			t.Errorf("第 %d 列命令碼 = %d，預期熱區 0x%02X − 7 = %d", row, got, hot, want)
		}
		if got := battleSideCommandRowLabel(row); got != wantLabels[row] {
			t.Errorf("第 %d 列 = %q，預期 %q", row, got, wantLabels[row])
		}
		if got := battleSideCommandRowOf(battleSideCommandRowCode[row]); got != row {
			t.Errorf("命令碼 %d 反查列 = %d，預期 %d", battleSideCommandRowCode[row], got, row)
		}
	}
	if got := battleSideCommandRowOf(99); got != -1 {
		t.Errorf("未知命令碼應回 −1，得到 %d", got)
	}
}

// handler 0x1C11A：col = (游標X − 0x1F0) >> 4，(游標Y − 0xF8) >= 0x10 再 +8。
func TestBattleFormationStripIndexMatchesRawHitMath(t *testing.T) {
	r := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH).SideFormation
	for _, c := range []struct {
		x, y, want int
	}{
		{496, 248, 0}, {511, 263, 0},
		{512, 248, 1}, {624 - 16, 248, 7},
		{496, 264, 8}, {624 - 16, 279, 15},
	} {
		got, ok := battleFormationIndexAt(r, c.x, c.y)
		if !ok || got != c.want {
			t.Errorf("(%d,%d) → %d,%v，預期 %d", c.x, c.y, got, ok, c.want)
		}
	}
	if _, ok := battleFormationIndexAt(r, 495, 248); ok {
		t.Error("陣形列左緣外不應命中")
	}
	if _, ok := battleFormationIndexAt(r, 496, 280); ok {
		t.Error("陣形列下緣外不應命中")
	}
	if got := battleFormationCellRect(r, 9); got != (battleRect{512, 264, 16, 16}) {
		t.Errorf("第 9 格 = %#v，預期 (512,264) 16×16", got)
	}
}

func TestDOSVBattleSideCommandsUseOriginalSingleColumnRows(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	cells := battleSideCommandCells(l.SideCommands)
	if len(cells) != len(battleCommandLabels) {
		t.Fatalf("右欄命令格數 = %d，預期 6", len(cells))
	}
	for i, cell := range cells {
		if !l.SideCommands.contains(cell) {
			t.Fatalf("右欄第 %d 列超出原版面板：%#v", i+1, cell)
		}
		if cell.X != l.SideCommands.X || cell.W != 128 || cell.H != 16 ||
			cell.Y != l.SideCommands.Y+i*16 {
			t.Fatalf("右欄第 %d 列 = %#v，預期單欄 128×16", i+1, cell)
		}
	}
}

func TestBattleCommandHitTestsRespectGlyphEdgesAndGaps(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	bottom := splitBattleCommandCells(l.BottomCommands)
	for i, cell := range bottom {
		hit, ok := battleHitRect(cell)
		if !ok {
			t.Fatalf("底列第 %d 格沒有可命中內框", i+1)
		}
		if got, ok := splitBattleCommandIndexAt(l.BottomCommands, hit.X+hit.W/2, hit.Y+hit.H/2); !ok || got != i {
			t.Fatalf("底列第 %d 格中心命中 = %d, %v", i+1, got, ok)
		}
		if _, ok := splitBattleCommandIndexAt(l.BottomCommands, cell.X+1, cell.Y+cell.H/2); ok {
			t.Fatalf("底列第 %d 格外框不應命中", i+1)
		}
		if _, ok := splitBattleCommandIndexAt(l.BottomCommands, cell.right()-1, cell.Y+cell.H/2); ok {
			t.Fatalf("底列第 %d 格右側邊界不應命中", i+1)
		}
	}
	if _, ok := splitBattleCommandIndexAt(l.BottomCommands, -1, l.BottomCommands.Y+1); ok {
		t.Fatal("畫面外不應命中底列命令")
	}

	side := battleSideCommandCells(l.SideCommands)
	for i, cell := range side {
		hit, ok := battleHitRect(cell)
		if !ok {
			t.Fatalf("右欄第 %d 格沒有可命中內框", i+1)
		}
		if got, ok := battleSideCommandIndexAt(l.SideCommands, hit.X+hit.W/2, hit.Y+hit.H/2); !ok || got != i {
			t.Fatalf("右欄第 %d 格中心命中 = %d, %v", i+1, got, ok)
		}
	}
	if _, ok := battleSideCommandIndexAt(l.SideCommands, l.SideCommands.X+1, l.SideCommands.Y+1); ok {
		t.Fatal("右欄命令外框不應命中")
	}
}

func TestBattleChoiceRowsHoverOnlyInsideGlyphBands(t *testing.T) {
	l := battleChoiceLayoutFor()
	for i, row := range l.Rows {
		if got, ok := battleChoiceRowAt(row.X+row.W/2, row.Y+row.H/2); !ok || got != i {
			t.Fatalf("遭遇第 %d 列中心命中 = %d, %v", i+1, got, ok)
		}
	}
	if _, ok := battleChoiceRowAt(l.Rows[0].X, l.Rows[0].bottom()); ok {
		t.Fatal("遭遇列距空白不應命中")
	}
	if _, ok := battleChoiceRowAt(l.Window.X+1, l.Window.Y+1); ok {
		t.Fatal("遭遇視窗框線不應命中")
	}
	if _, ok := battleChoiceRowAt(-1, -1); ok {
		t.Fatal("畫面外不應命中遭遇列")
	}
}

// 鍵盤 1–6 直接下命令是 **remake 加的**（原版這一層只有滑鼠）；
// 底列六格在原版是**選部隊**，兩者不再是同一個索引（docs/spec/33）。
func TestBattleKeyboardIssuesCommandsAndBottomRowPicksSquads(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	cells := splitBattleCommandCells(l.BottomCommands)
	keys := [...]ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3,
		ebiten.Key4, ebiten.Key5, ebiten.Key6,
	}
	for i, key := range keys {
		fromKey, ok := battleCommandIndexForKey(key)
		if !ok || fromKey != i {
			t.Fatalf("按鍵 %v 映射 = %d, %v；預期命令碼 %d", key, fromKey, ok, i)
		}
		hit, ok := battleHitRect(cells[i])
		if !ok {
			t.Fatalf("底列第 %d 格沒有 hit rect", i+1)
		}
		slot, ok := splitBattleCommandIndexAt(l.BottomCommands,
			hit.X+hit.W/2, hit.Y+hit.H/2)
		if !ok || slot != i {
			t.Fatalf("底列第 %d 格命中 = %d, %v", i+1, slot, ok)
		}
	}
}

// cs:0xD2E4 與 cs:0xD2EA 互為反排列——六格全部成立才算對得上。
func TestBottomSlotSquadTablesAreInverse(t *testing.T) {
	if len(battleBottomSlotSquad) != 6 || len(battleSquadSlotX) != 6 {
		t.Fatal("兩張順序表都必須是 6 筆")
	}
	seen := map[int]bool{}
	for slot, squad := range battleBottomSlotSquad {
		if squad < 0 || squad > 5 || seen[squad] {
			t.Fatalf("第 %d 格對到的隊 %d 不合法或重複", slot, squad)
		}
		seen[squad] = true
		if got, want := battleSquadSlotX[squad], slot*battleBottomSlotW; got != want {
			t.Errorf("隊 %d 的 X = %d，預期第 %d 格的 %d", squad, got, slot, want)
		}
		if got := battleSquadSlot(squad); got != slot {
			t.Errorf("隊 %d 反查格 = %d，預期 %d", squad, got, slot)
		}
	}
	// 畫面由左到右：左翼 左備 主將 前鋒 右備 右翼。
	want := []string{"左翼", "左備", "主將", "前鋒", "右備", "右翼"}
	for slot, squad := range battleBottomSlotSquad {
		if battleSquadLabels[squad] != want[slot] {
			t.Errorf("第 %d 格 = %q，預期 %q", slot, battleSquadLabels[squad], want[slot])
		}
	}
}

// sub_1C6BF 的兩個同心矩形：外框 (x+2,372)-(x+77,392)、內框各縮 1 px。
func TestSquadSelectRectsMatchRawGeometry(t *testing.T) {
	outer, inner := battleSlotSelectRects(80)
	if outer != (battleRect{82, 372, 76, 21}) {
		t.Errorf("外框 = %#v，預期 (82,372) 76×21", outer)
	}
	if inner != (battleRect{83, 373, 74, 19}) {
		t.Errorf("內框 = %#v，預期 (83,373) 74×19", inner)
	}
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	if !l.BottomCommands.contains(outer) {
		t.Errorf("選取框 %#v 必須留在底列內", outer)
	}
}

const chromeTileForContract = 8

// 門強度條的版面（docs/spec/32 §2）。這幾個數字全部是原版直接座標，
// 而且 264 正好是 TALK 框右緣 256 再 +8——條不會被 TALK 蓋掉。
func TestGateBarGeometryMatchesRawConstants(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	if gateBarLabelX != 264 || gateBarLabelY != 8 {
		t.Errorf("標籤 = (%d,%d)，預期 sub_1C407 的 (264,8)", gateBarLabelX, gateBarLabelY)
	}
	if gateBarX != 320 || gateBarY != 16 || gateBarLen != 0x97 || gateBarH != 2 {
		t.Errorf("條 = (%d,%d) %d×%d，預期 sub_1C4D2 的 (320,16) 151×2",
			gateBarX, gateBarY, gateBarLen, gateBarH)
	}
	// sub_10BCD 清到 471（含），條蓋 320..470——原版就差這 1 px，照抄。
	if right := gateBarX + gateBarLen - 1; right != gateBarClearX1-1 {
		t.Errorf("條的右緣 = %d，預期清除區右緣 %d 再往左 1 px",
			right, gateBarClearX1)
	}
	if gateBarLabelX != l.TopTalk.right()+8 {
		t.Errorf("標籤 x=%d 應在 TALK 框右緣 %d 之後 8 px",
			gateBarLabelX, l.TopTalk.right())
	}
	bar := battleRect{X: gateBarX, Y: gateBarY, W: gateBarLen, H: gateBarH}
	label := battleRect{X: gateBarLabelX, Y: gateBarLabelY,
		W: gateBarClearX1 - gateBarLabelX + 1, H: 15}
	if !l.Field.contains(bar) || !l.Field.contains(label) {
		t.Error("門強度條與標籤必須落在戰場 viewport 內")
	}
	if bar.overlaps(l.Sidebar) || label.overlaps(l.Sidebar) {
		t.Error("門強度條不得壓到側欄")
	}
}

// 陣形選單是 8 欄 × 2 列的 16×16 格（原版 handler 0x1C11A，docs/spec/37 §2.1）。
func TestBattleFormationPickerMatchesRawGeometry(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	r := l.SideFormation
	if r.X != 496 || r.Y != 248 || r.W != 128 || r.H != 32 {
		t.Fatalf("陣形選單 = (%d,%d,%d×%d)，原版是 (496,248,128×32)", r.X, r.Y, r.W, r.H)
	}
	// 左上第一格、右上第八格、左下第九格、右下第十六格。
	for _, tc := range []struct {
		x, y, want int
	}{
		{496, 248, 0}, {496 + 7*16, 248, 7},
		{496, 264, 8}, {496 + 7*16, 279, 15},
	} {
		got, ok := battleFormationIndexAt(r, tc.x, tc.y)
		if !ok || got != tc.want {
			t.Errorf("(%d,%d) → %d/%v，預期 %d", tc.x, tc.y, got, ok, tc.want)
		}
	}
	if _, ok := battleFormationIndexAt(r, 495, 248); ok {
		t.Error("左緣外一格不該命中")
	}
	// 選取框：格內縮 1 px 畫 14×14。
	cell := battleFormationCellRect(r, 9)
	if cell.X != 496+16 || cell.Y != 264 || cell.W != 16 || cell.H != 16 {
		t.Errorf("第 9 格 = (%d,%d,%d×%d)", cell.X, cell.Y, cell.W, cell.H)
	}
}

// 陣形線三格（熱區 0x04–0x06）的矩形照原版（docs/spec/37 §2.2）。
func TestBattleFormationLineHitBoxes(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	want := [3]battleRect{
		{X: 552, Y: 288, W: 64, H: 24},
		{X: 552, Y: 312, W: 64, H: 32},
		{X: 552, Y: 344, W: 64, H: 24},
	}
	if l.SideLines != want {
		t.Fatalf("陣形線三格 = %v，原版是 %v", l.SideLines, want)
	}
	for i, r := range want {
		got, ok := battleLineIndexAt(l.SideLines, r.X, r.Y)
		if !ok || got != i {
			t.Errorf("第 %d 格左上角命中 %d/%v", i, got, ok)
		}
	}
	if _, ok := battleLineIndexAt(l.SideLines, 552, 287); ok {
		t.Error("第一格上緣外不該命中")
	}
	// 由上而下寫進去的是 48／28／5（敵軍側／中央／自軍側）。
	for k, want := range [3]int{48, 28, 5} {
		if got := tactical.LineFor(tactical.AttackerSide, 2-k); got != want {
			t.Errorf("第 %d 格 → 陣形線 %d，預期 %d", k, got, want)
		}
	}
}

// 一覽視窗的幾何照原版：(24,88,384,176)、一列 16 px、一頁 10 列
// （docs/re/26 §2，規格 docs/spec/38 §1.1）。
func TestListWindowGeometryMatchesRaw(t *testing.T) {
	if listWinX != 24 || listWinY != 88 || listWinW != 384 || listWinH != 176 {
		t.Fatalf("一覽視窗 = (%d,%d,%d×%d)，原版是 (24,88,384×176)",
			listWinX, listWinY, listWinW, listWinH)
	}
	if listRowH != 16 || listRowsPerPage != 10 {
		t.Fatalf("列高 %d／一頁 %d 列，原版是 16／10", listRowH, listRowsPerPage)
	}
	// 第 0 列緊接在標題那一列之下。
	if got := listRowY(0); got != listWinY+listRowH {
		t.Errorf("第 0 列 y = %d，預期 %d", got, listWinY+listRowH)
	}
}

// 四個家族的標題與欄數照原版（docs/re/26 §4.1）。
func TestListFamilyHeadersMatchRawStrings(t *testing.T) {
	for _, tc := range []struct {
		fam   listFamily
		title string
		cols  int
	}{
		{listFamilyCorps, "武將名　總兵數　士氣值 現在位置 目標據點", 5},
		{listFamilyCities, "據點名　生產力　上昇率　防災　城兵　內政官", 6},
		{listFamilyGenerals, "武將名　武術 統率 政治　　勢力　　　身分", 6},
		{listFamilyFactions, "勢力名　武將　據點　首都　　外交　　外交官", 6},
	} {
		if tc.fam.Title != tc.title {
			t.Errorf("標題 = %q，原版是 %q", tc.fam.Title, tc.title)
		}
		if got := len(tc.fam.fields()); got != tc.cols {
			t.Errorf("%q 切出 %d 欄，原版的欄數是 %d", tc.title, got, tc.cols)
		}
	}
}

// 欄的 x 與寬度是**從分隔線算出來的**，不是自己編的。
func TestListFieldsComeFromSeparator(t *testing.T) {
	f := listFamilyGenerals.fields()
	// 武將名是三個全形 ＝ 48 px，從 0 開始。
	if f[0].X != 0 || f[0].W != 48 || f[0].Numeric {
		t.Errorf("第 0 欄 = %+v，預期 x0 w48 文字欄", f[0])
	}
	// 武術是兩個半形 ＝ 16 px 的數字欄。
	if f[1].W != 16 || !f[1].Numeric {
		t.Errorf("第 1 欄 = %+v，預期 w16 數字欄", f[1])
	}
	// 最後一欄不能超出視窗內緣。
	last := f[len(f)-1]
	if listWinX+last.X+last.W > listWinX+listWinW {
		t.Errorf("最後一欄右緣 %d 超出視窗", last.X+last.W)
	}
}

// 外交六級的換算照 `sub_17A7A`（docs/re/27 §4）。
func TestListDiplomacyLevelsFollowRaw(t *testing.T) {
	for _, tc := range []struct {
		raw   int
		atWar bool
		want  int
	}{
		{50, true, 0}, {0, false, 1}, {19, false, 1}, {20, false, 2},
		{59, false, 3}, {60, false, 4}, {80, false, 5}, {100, false, 5},
		{101, false, 6},
	} {
		if got := listDiplomacyLevel(tc.raw, tc.atWar); got != tc.want {
			t.Errorf("交友度 %d（交戰 %v）→ %d（%s），預期 %d（%s）",
				tc.raw, tc.atWar, got, listDiplomacyNames[got],
				tc.want, listDiplomacyNames[tc.want])
		}
	}
}
