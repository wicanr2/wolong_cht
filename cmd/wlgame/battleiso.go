package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/ui/isoview"
)

// 戰場的等角投影畫在 `internal/ui/isoview`——**手機版用同一份**。
// 這裡只做桌面版的版面：把那張原生畫布放到 DOS/V 版面的戰場區。

// isoScale 是桌面版的放大倍率。原生畫布正好是 DOS/V 的可見 viewport，
// 所以是 1；改這個值之前先看 `TestBattleViewBufferMatchesDOSVVisibleViewport`。
const isoScale = 1

// newBattleView 準備一張戰場的繪圖資源。缺素材就回 nil，
// 呼叫端會退回沒有美術的畫法。
func (g *game) newBattleView(field int) *isoview.View {
	if g.battleLib == nil || g.lib == nil || g.lib.Palette == nil {
		return nil
	}
	bank, err := g.lib.Palette.Bank(0)
	if err != nil {
		return nil
	}
	opt := isoview.Options{
		Lib: g.battleLib, Palette: bank, Sprites: g.battleSprites,
		Field: field, Rotate: g.battle.Rotate(), Rand: g.rng.Next,
	}
	if g.battleCamAt != nil {
		at := [2]int{g.battleCamAt[0], g.battleCamAt[1]}
		opt.CamAt = &at
	}
	return isoview.New(opt)
}

func (g *game) drawBattleIso(screen *ebiten.Image, b *tactical.Battle, me *tactical.Soldier) {
	l := dosvBattleLayoutFor(screenW, screenH)
	buf := g.view.Render(b)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(isoScale, isoScale)
	op.GeoM.Translate(float64(l.Field.X), float64(l.Field.Y))
	screen.DrawImage(buf, op)
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
