package wolongmobile

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	mobileui "github.com/wicanr2/wolong_cht/internal/ui/mobile"
)

const (
	logicalWidth  = 640
	logicalHeight = 496
)

var (
	ink    = color.RGBA{R: 8, G: 18, B: 46, A: 255}
	panel  = color.RGBA{R: 13, G: 30, B: 65, A: 255}
	blue   = color.RGBA{R: 29, G: 73, B: 132, A: 255}
	red    = color.RGBA{R: 162, G: 38, B: 39, A: 255}
	cream  = color.RGBA{R: 230, G: 224, B: 178, A: 255}
	muted  = color.RGBA{R: 169, G: 190, B: 221, A: 255}
	bright = color.RGBA{R: 255, G: 241, B: 188, A: 255}
	black  = color.RGBA{A: 255}
)

// Game is a deliberately small mobile shell prototype. It owns no game
// rules, original data, font, or save file; those remain future integration
// work behind the existing desktop core.
type game struct {
	buttons    []mobileui.Button
	lastAction string
	page       int
	menuOpen   bool
	frame      int
}

func newGame() *game {
	return &game{
		buttons:    mobileui.DockButtons(logicalWidth, logicalHeight),
		lastAction: "WAITING FOR TOUCH",
	}
}

func (g *game) Update() error {
	g.frame++
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		g.handleTap(float64(x), float64(y))
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handleTap(float64(x), float64(y))
	}
	return nil
}

func (g *game) handleTap(x, y float64) {
	for _, button := range g.buttons {
		if !button.Hit(x, y) {
			continue
		}
		switch button.ID {
		case "continue":
			g.page = g.page%3 + 1
			g.lastAction = fmt.Sprintf("TALK page %d/3", g.page)
		case "menu":
			g.menuOpen = !g.menuOpen
			g.lastAction = "COMMAND DRAWER OPEN"
		case "save":
			g.lastAction = "SAVE ENTRY (PROTOTYPE)"
		}
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(ink)
	if len(g.buttons) == 0 {
		g.buttons = mobileui.DockButtons(logicalWidth, logicalHeight)
	}

	fill(screen, 0, 0, logicalWidth, 32, red)
	fill(screen, 0, 32, logicalWidth, 32, blue)
	fill(screen, 0, 64, 432, 336, color.RGBA{R: 44, G: 71, B: 71, A: 255})
	fill(screen, 432, 64, 208, 336, panel)
	fill(screen, 0, 400, logicalWidth, 96, color.RGBA{R: 7, G: 14, B: 31, A: 255})

	for x := 0; x < 27; x++ {
		for y := 0; y < 21; y++ {
			if (x+y+g.frame/30)%5 == 0 {
				fill(screen, float64(x*16+3), float64(67+y*16), 10, 10, color.RGBA{R: 74, G: 111, B: 82, A: 255})
			}
		}
	}
	fill(screen, 46, 214, 172, 10, color.RGBA{R: 31, G: 92, B: 202, A: 255})
	fill(screen, 198, 112, 10, 120, color.RGBA{R: 31, G: 92, B: 202, A: 255})
	fill(screen, 480, 88, 128, 84, color.RGBA{R: 60, G: 93, B: 104, A: 255})
	fill(screen, 480, 174, 128, 16, red)

	ebitenutil.DebugPrintAt(screen, "WOLONG REMAKE  /  MOBILE SHELL", 12, 8)
	ebitenutil.DebugPrintAt(screen, "CLASSIC REVIVAL", 12, 40)
	ebitenutil.DebugPrintAt(screen, "DOS/V 640x400 logical canvas", 14, 76)
	ebitenutil.DebugPrintAt(screen, "NATURAL MAP", 18, 92)
	ebitenutil.DebugPrintAt(screen, "君主  曹操", 450, 202)
	ebitenutil.DebugPrintAt(screen, "首都  許昌", 450, 220)
	ebitenutil.DebugPrintAt(screen, "軍師  荀彧", 450, 238)
	ebitenutil.DebugPrintAt(screen, "信賴度 255", 450, 270)
	ebitenutil.DebugPrintAt(screen, "TOUCH PROTOTYPE", 450, 300)
	ebitenutil.DebugPrintAt(screen, "ACTION: "+g.lastAction, 450, 320)

	for _, button := range g.buttons {
		fill(screen, button.Bounds.X, button.Bounds.Y, button.Bounds.W, button.Bounds.H, color.RGBA{R: 31, G: 48, B: 82, A: 255})
		ebitenutil.DrawRect(screen, button.Bounds.X, button.Bounds.Y, button.Bounds.W, button.Bounds.H, cream)
		ebitenutil.DebugPrintAt(screen, button.Label, int(button.Bounds.X+18), int(button.Bounds.Y+25))
	}
	if g.menuOpen {
		fill(screen, 110, 110, 420, 190, color.RGBA{R: 4, G: 8, B: 18, A: 245})
		ebitenutil.DrawRect(screen, 110, 110, 420, 190, bright)
		ebitenutil.DebugPrintAt(screen, "MOBILE COMMAND DRAWER", 132, 132)
		ebitenutil.DebugPrintAt(screen, "MAP  /  TALK  /  ORDER  /  STATUS", 132, 170)
		ebitenutil.DebugPrintAt(screen, "safe-area + 48dp hitbox prototype", 132, 208)
		ebitenutil.DebugPrintAt(screen, "MENU 再按一次關閉", 132, 246)
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return logicalWidth, logicalHeight
}

func fill(screen *ebiten.Image, x, y, w, h float64, c color.Color) {
	img := ebiten.NewImage(int(w), int(h))
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(img, op)
}

var _ ebiten.Game = (*game)(nil)
