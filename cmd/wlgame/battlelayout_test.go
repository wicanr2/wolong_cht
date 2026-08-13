package main

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
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
		"attacker":        l.SideAttacker,
		"mini map":        l.SideMiniMap,
		"defender":        l.SideDefender,
		"side commands":   l.SideCommands,
	} {
		if r.X < 0 || r.Y < 0 || r.W <= 0 || r.H <= 0 ||
			r.right() > dosvBattleScreenW || r.bottom() > dosvBattleScreenH {
			t.Errorf("%s = %#v 不在 640×400 內", name, r)
		}
	}

	if !l.Field.contains(l.TopTalk) || !l.Field.contains(l.BottomTalk) {
		t.Fatal("TALK overlay 必須留在戰場 viewport 內")
	}
	if !l.Sidebar.contains(l.SideAttacker) || !l.Sidebar.contains(l.SideMiniMap) ||
		!l.Sidebar.contains(l.SideDefender) || !l.Sidebar.contains(l.SideCommands) {
		t.Fatal("右欄子區域必須留在右欄內")
	}
	if l.SideMiniMap.Y-l.SideAttacker.bottom() != chromeTileForContract ||
		l.SideDefender.Y-l.SideMiniMap.bottom() != chromeTileForContract {
		t.Fatal("右欄狀態／縮圖區之間必須保留 8 px 原版框距")
	}
	if got, want := l.SideCommands, (battleRect{496, 280, 128, 96}); got != want {
		t.Fatalf("右欄命令面板 = %#v，預期 sub_1C863 直接坐標 %#v", got, want)
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

func TestDOSVBattleSideStatusColumnsDoNotOverlap(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	for _, panel := range []battleRect{l.SideAttacker, l.SideDefender} {
		cols := battleSideStatusColumns(panel)
		if cols.HP.overlaps(cols.Advantage) {
			t.Fatal("體力與優劣欄不得重疊")
		}
		if gap := cols.Advantage.X - cols.HP.right(); gap < battleCommandMinPad {
			t.Fatalf("體力與優劣欄間距 = %d，至少需要 %d", gap, battleCommandMinPad)
		}
		if got := battleCommandTextWidth("體力 100"); got > cols.HP.W {
			t.Fatalf("最大體力文字寬 %d 超出安全欄 %d", got, cols.HP.W)
		}
		if got := battleCommandTextWidth("不利"); got > cols.Advantage.W {
			t.Fatalf("優劣文字寬 %d 超出安全欄 %d", got, cols.Advantage.W)
		}
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

func TestBattleKeyboardAndPointerResolveSameCommandIndex(t *testing.T) {
	l := dosvBattleLayoutFor(dosvBattleScreenW, dosvBattleScreenH)
	cells := splitBattleCommandCells(l.BottomCommands)
	keys := [...]ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3,
		ebiten.Key4, ebiten.Key5, ebiten.Key6,
	}
	for i, key := range keys {
		fromKey, ok := battleCommandIndexForKey(key)
		if !ok || fromKey != i {
			t.Fatalf("按鍵 %v 映射 = %d, %v；預期 %d", key, fromKey, ok, i)
		}
		hit, ok := battleHitRect(cells[i])
		if !ok {
			t.Fatalf("命令 %d 沒有 hit rect", i+1)
		}
		fromPointer, ok := splitBattleCommandIndexAt(l.BottomCommands,
			hit.X+hit.W/2, hit.Y+hit.H/2)
		if !ok || fromPointer != fromKey {
			t.Fatalf("命令 %d 鍵鼠結果不同：key=%d pointer=%d", i+1, fromKey, fromPointer)
		}
	}
}

const chromeTileForContract = 8
