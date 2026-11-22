package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/ui/isoview"
)

// 戰場的等角投影畫在 `internal/ui/isoview`——**手機版用同一份**。
// 這裡只做桌面版的版面：把那張原生畫布放到 DOS/V 版面的戰場區。

// newBattleView 準備一張戰場的繪圖資源。缺素材就回 nil，
// 呼叫端會退回沒有美術的畫法。
func (g *game) newBattleView(field int, rotate bool) *isoview.View {
	if g.battle == nil || g.lib == nil || g.lib.Palette == nil {
		return nil
	}
	lib := g.battle.Library()
	if lib == nil {
		return nil
	}
	bank, err := g.lib.Palette.Bank(0)
	if err != nil {
		return nil
	}
	opt := isoview.Options{
		Lib: lib, Palette: bank, Sprites: g.battle.Sprites(),
		Field: field, Rotate: rotate, Rand: g.rng.Next,
	}
	if g.battleCamAt != nil {
		at := [2]int{g.battleCamAt[0], g.battleCamAt[1]}
		opt.CamAt = &at
	}
	return isoview.New(opt)
}

func (g *game) drawBattleIso(screen *ebiten.Image, b *tactical.Battle, field battleRect) {
	if g.hdFrame != nil {
		if native := g.view.RenderHD(b, g.hdBattle); native != nil {
			g.hdFrame.nativeRect(screen, native, field.X, field.Y, field.W, field.H)
			g.hdBattleDrawn = true
			return
		}
	}
	buf := g.view.Render(b)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(field.W)/float64(buf.Bounds().Dx()), float64(field.H)/float64(buf.Bounds().Dy()))
	op.GeoM.Translate(float64(field.X), float64(field.Y))
	screen.DrawImage(buf, op)
}

func (g *game) battleZoomReady() bool {
	return g.view != nil && g.desktopBattleVisible() && !g.battleFastForward && !g.world.PendingBattle().Battle.Done
}

func (g *game) zoomBattleWheel(dy float64, x, y int, down [3]bool) bool {
	field := dosvBattleLayoutFor(screenW, screenH).Field
	if dy == 0 || down != [3]bool{} || !g.battleZoomReady() || !field.containsPoint(x, y) {
		return false
	}
	g.view.ZoomAt(g.view.ZoomPercent()+wheelZoomStep(dy), x-field.X, y-field.Y)
	return true
}

// viewCursorX／viewCursorY 是小地圖十字的位置，沒有戰場畫面時回 0。
func (g *game) viewCursorX() int {
	if g.view == nil {
		return 0
	}
	x, _ := g.view.Cursor()
	return x
}

func (g *game) viewCursorY() int {
	if g.view == nil {
		return 0
	}
	_, y := g.view.Cursor()
	return y
}
