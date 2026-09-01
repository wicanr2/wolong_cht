package main

import (
	"image/color"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

// colorSentinel 是一個原版調色盤不會出現的值，用來分辨「查到了」與
// 「走了 fallback」。
var colorSentinel = color.RGBA{1, 2, 3, 255}

// 啟動殼層沒有 World，UI 顏色仍然要查調色盤——原版在那一層固定切第 0 組
// （`sub_11A6E` 的 `mov al, 0`，docs/spec/107）。
func TestUIPaletteBankFallsBackToLauncherSeason(t *testing.T) {
	g := &game{}
	if got := g.uiPaletteBank(); got != launcherSeason {
		t.Errorf("沒有世界時組別 ＝ %d，want %d", got, launcherSeason)
	}

	w, err := state.LoadScenario("../../workplace/orig/dosv/SINARIO.DAT", 0)
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	g.world = w
	if got, want := g.uiPaletteBank(), int(w.Clock.Season()); got != want {
		t.Errorf("有世界時組別 ＝ %d，want %d（跟著季節）", got, want)
	}
}

// ⭐ 正對照：**沒有世界的時候也要拿到真的顏色**，不能退回 fallback。
// 這一條擋的正是先前的缺口——君主卡的兩顆鈕整片變灰（docs/spec/107 §2）。
func TestPaletteInkWithoutWorldUsesRealPalette(t *testing.T) {
	lib, err := library.Load("../../workplace/orig/dosv")
	if err != nil {
		t.Skipf("沒有原版素材：%v", err)
	}
	g := &game{lib: lib} // 世界是 nil —— 這就是啟動殼層的狀態

	for _, idx := range []int{chrome.SheetIndex, chrome.InkIndex, 0x06, 0x07} {
		want, err := lib.PaletteColor(launcherSeason, idx)
		if err != nil {
			t.Fatalf("取不到第 0 組的色 %d：%v", idx, err)
		}
		// fallback 故意給一個不可能出現的值，取到它就表示走了退路。
		sentinel := colorSentinel
		if got := g.paletteInk(idx, sentinel); got != want {
			t.Errorf("色 %d ＝ %v，want %v（取到 fallback 就是又退回去了）",
				idx, got, want)
		}
	}
}

// 沒有素材時仍要回 fallback，不能讓畫面消失。
func TestPaletteInkWithoutLibraryKeepsFallback(t *testing.T) {
	g := &game{}
	if got := g.paletteInk(chrome.SheetIndex, colorSentinel); got != colorSentinel {
		t.Errorf("沒有素材時 ＝ %v，want fallback %v", got, colorSentinel)
	}
}
