package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/wicanr2/wolong_cht/internal/assets/world"
)

// 大地圖的可視範圍。原版主畫面最上方有一條 32 px 的橫幅
// （ICONGRF 段 0，docs/formats/03 §5.1），所以地圖區是 640×368，
// 也就是 40 × 23 格 —— 這個數字與原版的畫面圖塊快取一致
// （docs/re/05 提到的 (bx×40 + dx)×8 索引算式）。
const (
	viewCols   = screenW / world.TileSize             // 40
	viewRows   = (screenH - bannerH) / world.TileSize // 23
	bannerH    = 32
	scrollFast = 8 // 按住 Shift 時一次捲幾格
)

// worldView 是大地圖的捲動狀態。
type worldView struct {
	x, y   int // 左上角所在的格
	season int
}

func (w *worldView) clamp() {
	if w.x < 0 {
		w.x = 0
	}
	if w.y < 0 {
		w.y = 0
	}
	if max := world.Width - viewCols; w.x > max {
		w.x = max
	}
	if max := world.Height - viewRows; w.y > max {
		w.y = max
	}
}

func (w *worldView) update() {
	step := 1
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		step = scrollFast
	}
	// 方向鍵用「按著就持續捲」，不是「按一下捲一格」——
	// 384×256 格的地圖，一格一格點會點到天亮。
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		w.x += step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		w.x -= step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		w.y += step
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		w.y -= step
	}
	if pressed(ebiten.KeyHome) {
		w.x, w.y = 0, 0
	}
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if pressed(k) {
			w.season = i
		}
	}
	w.clamp()
}

func (a *app) drawWorld(screen *ebiten.Image) {
	img, err := a.lib.RenderWorld(a.world.x, a.world.y, viewCols, viewRows, a.world.season)
	if err != nil {
		ebitenutil.DebugPrint(screen, err.Error())
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, bannerH)
	screen.DrawImage(ebiten.NewImageFromImage(img), op)

	// 原版的橫幅就畫在最上方那 32 px（ICONGRF 段 0）。
	// 找得到就畫上去，讓版面與原版一致。
	if i := a.bannerIndex(); i >= 0 {
		if b, err := a.lib.Render(i, 0, a.world.season); err == nil {
			screen.DrawImage(ebiten.NewImageFromImage(b), &ebiten.DrawImageOptions{})
		}
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf(
		"WORLD %dx%d  view (%d,%d)  season=%s  [arrows]scroll [shift]fast [Home]origin [tab]assets [F10]quit",
		world.Width, world.Height, a.world.x, a.world.y,
		seasonNames[a.world.season]), 0, screenH-16)
}

// bannerIndex 找出 ICONGRF 段 0（標題橫幅）在素材清單裡的位置。
func (a *app) bannerIndex() int {
	for i, e := range a.lib.Entries {
		if e.Label == "ICONGRF/banner" {
			return i
		}
	}
	return -1
}
