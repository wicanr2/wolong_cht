package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/namechars"
)

func namingTestTable(t *testing.T) *namechars.Table {
	t.Helper()
	path := filepath.Join("..", "..", "workplace", "orig", "dosv", "END_S15.DAT")
	if _, err := os.Stat(path); err != nil {
		t.Skip("沒有原版素材")
	}
	tb, err := namechars.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return tb
}

// 選字：點格寫進目前格、游標前進；格間空隙不算；名字位元組是 Big5。
func TestNamingPickWritesCellAndAdvances(t *testing.T) {
	m := newNamingModel(namingTestTable(t), nil)
	// 第 0 頁第 1 格（列 0 欄 1）＝「八」A44B。
	m.click(namingGridX+namingGridPitch+3, namingGridY+3)
	if m.cells[0] != 0xA44B || m.cursor != 1 {
		t.Fatalf("cells[0]=%04X cursor=%d", m.cells[0], m.cursor)
	}
	// 點在格與格之間的 4 px 空隙：不算。
	m.click(namingGridX+16+1, namingGridY+3)
	if m.cursor != 1 {
		t.Fatal("空隙也算成選字")
	}
	if got := m.nameBytes(); got[0] != 0xA4 || got[1] != 0x4B || got[2] != 0xA1 || got[3] != 0x40 {
		t.Fatalf("nameBytes = % X", got)
	}
	if !m.hasName() {
		t.Fatal("第一格有字卻說沒名字")
	}
}

// 重來＝清掉目前格退一格；繼續＝清掉目前格跳一格；到第六格不再前進。
func TestNamingRedoAndContinue(t *testing.T) {
	m := newNamingModel(namingTestTable(t), nil)
	for i := 0; i < 8; i++ {
		m.click(namingGridX+3, namingGridY+3) // 「ㄅ」
	}
	if m.cursor != namingCells-1 {
		t.Fatalf("八次選字後游標在 %d", m.cursor)
	}
	m.click(namingRedoX+4, namingBtnY+4)
	if m.cells[namingCells-1] != 0 || m.cursor != namingCells-2 {
		t.Fatal("重來沒有清掉目前格並退一格")
	}
	m.click(namingContX+4, namingBtnY+4)
	if m.cells[namingCells-2] != 0 || m.cursor != namingCells-1 {
		t.Fatal("繼續沒有清掉目前格並跳一格")
	}
}

// 聲母列跳到分段標記那一頁、翻頁夾在 0..0x13BA。
func TestNamingInitialJumpAndPaging(t *testing.T) {
	m := newNamingModel(namingTestTable(t), nil)
	m.click(namingInitialHotX+32*4+8, namingInitialsY+4) // ㄐ
	if m.table.Runes[m.page] != 'ㄐ' {
		t.Fatalf("跳到 %q", m.table.Runes[m.page])
	}
	m.click(200+8, namingPagerY+4) // 上一頁
	if m.table.Runes[m.page] == 'ㄐ' {
		t.Fatal("上一頁沒有動")
	}
	for i := 0; i < 40; i++ {
		m.click(368+8, namingPagerY+4)
	}
	if m.page != namingMaxPage {
		t.Fatalf("下一頁夾在 %d，預期 %d", m.page, namingMaxPage)
	}
}

// 肖像 ◀▶ 循環且跳過武將在用的號碼。
func TestNamingPortraitSkipsUsedOnes(t *testing.T) {
	m := newNamingModel(namingTestTable(t), nil)
	m.used[0x92] = true
	m.portrait = 0x91
	m.click(272+8, 176+8) // 後 ▼
	if m.portrait != 0 {
		t.Fatalf("0x91 → 後一張應跳過 0x92 回到 0，得到 %#x", m.portrait)
	}
	m.click(272+8, 152+8) // 前 ▲
	if m.portrait != 0x91 {
		t.Fatalf("回前一張應是 0x91，得到 %#x", m.portrait)
	}
}

// 空名字不能「確定」。
func TestNamingRejectsEmptyConfirm(t *testing.T) {
	m := newNamingModel(namingTestTable(t), nil)
	g := &game{naming: m}
	m.click(namingOKX+4, namingBtnY+4)
	g.settleNaming()
	if g.naming == nil || g.customAdvisor != nil {
		t.Fatal("空名字被放行了")
	}
	m.click(namingGridX+3, namingGridY+3)
	m.click(namingOKX+4, namingBtnY+4)
	g.settleNaming()
	if g.naming != nil || g.customAdvisor == nil || g.customAdvisor.portrait != 0x91 {
		t.Fatal("有名字卻沒有收尾")
	}
}
