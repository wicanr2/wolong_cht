package main

// 戰略地圖呈現差異與驗收契約見 docs/spec/275。100%保留原版整格分支。
import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/translations"
)

type worldZoomState struct {
	percent              int
	fractionX, fractionY float64
	layer                *ebiten.Image
	drawing              bool
	sourceW, sourceH     int
}

var desktopZoomLabels = func() map[string][]string {
	b, err := translations.Files.ReadFile("desktop-zoom.json")
	var m map[string][]string
	if err == nil {
		err = json.Unmarshal(b, &m)
	}
	if err != nil {
		log.Panicf("desktop zoom labels: %v", err)
	}
	for _, lang := range []string{"zh-hant", "zh-hans", "ja", "en"} {
		if len(m[lang]) != 3 {
			log.Panicf("desktop zoom labels: %s", lang)
		}
	}
	return m
}()

func (g *game) mapZoomPercent() int {
	if g.worldZoom.percent == 0 {
		return 100
	}
	return g.worldZoom.percent
}

func (g *game) mapZoomScale() float64 { return float64(g.mapZoomPercent()) / 100 }

func (g *game) mapOrigin() (float64, float64) {
	x, y := g.camX*16, g.camY*16
	if g.mapZoomPercent() != 100 {
		x += g.camSubX
		y += g.camSubY
	}
	return float64(x), float64(y)
}

func (g *game) setZoomOrigin(x, y float64) {
	z := g.mapZoomScale()
	if g.mapZoomPercent() == 100 {
		g.camX, g.camY = int(math.Round(x/16)), int(math.Round(y/16))
		g.camSubX, g.camSubY = 0, 0
		g.clampCam()
		return
	}
	px := max(0, min(int(math.Floor(384*16-strategyMapW/z)), int(math.Round(x))))
	py := max(0, min(int(math.Floor(256*16-strategyMapH/z)), int(math.Round(y))))
	g.camX, g.camSubX = px/16, px%16
	g.camY, g.camSubY = py/16, py%16
}

// anchorY是地圖內座標，不含橫幅。按倍率鈕使用視野中心，滾輪使用游標位置。
func (g *game) zoomMapTo(percent, anchorX, anchorY int) {
	percent = max(50, min(200, percent))
	if percent == g.mapZoomPercent() {
		return
	}
	x, y := g.mapOrigin()
	old := g.mapZoomScale()
	x += float64(anchorX) / old
	y += float64(anchorY) / old
	g.worldZoom.percent = percent
	g.worldZoom.fractionX, g.worldZoom.fractionY = 0, 0
	z := g.mapZoomScale()
	g.setZoomOrigin(x-float64(anchorX)/z, y-float64(anchorY)/z)
}

func (g *game) resetMapZoom() {
	g.worldZoom.percent = 0
	g.worldZoom.fractionX, g.worldZoom.fractionY = 0, 0
}

func (g *game) strategyZoomReady() bool {
	return g.desktopHostReady() && !g.hudOpen(hudSystem)
}

func (g *game) strategyMapPointVisible(x, y int) bool {
	if x < 0 || x >= strategyMapW || y < strategyMapY || y >= screenH {
		return false
	}
	return !(g.hudOpen(hudCommand) && x < strategyCommandW && y < strategyCommandY+strategyCommandH ||
		g.hudOpen(hudMinimap) && x >= strategySidebarX && y < strategyFactionY ||
		g.hudOpen(hudFaction) && x >= strategySidebarX && y >= strategyFactionY)
}

func (g *game) zoomMapWheel(dy float64, x, y int, down [3]bool) bool {
	if dy == 0 || down != [3]bool{} || !g.strategyZoomReady() || !g.strategyMapPointVisible(x, y) {
		return false
	}
	step := 25
	if dy < 0 {
		step = -25
	}
	g.zoomMapTo(g.mapZoomPercent()+step, x, y-strategyMapY)
	return true
}

func (g *game) desktopZoomLabel(a desktopAction) string {
	rows := desktopZoomLabels[desktopLocale()]
	switch a {
	case desktopZoomOut:
		return rows[0]
	case desktopZoomReset:
		return fmt.Sprintf(rows[1], g.mapZoomPercent())
	case desktopZoomIn:
		return rows[2]
	}
	return ""
}

func (g *game) drawStrategyMap(screen *ebiten.Image) {
	if g.drawHDWorld(screen) {
		return
	}
	season := g.paletteBank()
	if g.mapZoomPercent() == 100 {
		if img, err := g.lib.RenderWorldMarked(g.camX, g.camY, viewCols, viewRows, season, g.cityMarks(), g.corpsMarks()); err == nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, float64(strategyMapY))
			screen.DrawImage(ebiten.NewImageFromImage(img), op)
		}
		g.drawDisasterOverlay(screen)
		return
	}
	z := g.mapZoomScale()
	cols := int(math.Ceil((strategyMapW/z + float64(g.camSubX)) / 16))
	rows := int(math.Ceil((strategyMapH/z + float64(g.camSubY)) / 16))
	w, h := cols*16, strategyMapY+rows*16
	q := &g.worldZoom
	if q.layer == nil || q.layer.Bounds().Dx() != w || q.layer.Bounds().Dy() != h {
		if q.layer != nil {
			q.layer.Deallocate()
		}
		q.layer = ebiten.NewImage(w, h)
	}
	q.layer.Clear()
	if img, err := g.lib.RenderWorldMarked(g.camX, g.camY, cols, rows, season, g.cityMarks(), g.corpsMarks()); err == nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(strategyMapY))
		q.layer.DrawImage(ebiten.NewImageFromImage(img), op)
	}
	q.drawing, q.sourceW, q.sourceH = true, w, h
	g.drawDisasterOverlay(q.layer)
	g.drawMapPickCursor(q.layer)
	q.drawing = false
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	op.GeoM.Scale(z, z)
	op.GeoM.Translate(-float64(g.camSubX)*z, strategyMapY-float64(g.camSubY)*z)
	src := q.layer.SubImage(image.Rect(0, strategyMapY, w, h)).(*ebiten.Image)
	dst := screen.SubImage(image.Rect(0, strategyMapY, strategyMapW, screenH)).(*ebiten.Image)
	dst.DrawImage(src, op)
}

func (g *game) mapSourceBounds() (int, int) {
	if g.worldZoom.drawing {
		return g.worldZoom.sourceW, g.worldZoom.sourceH
	}
	return strategyMapW, screenH
}

// 小地圖原有24×11框在100%完全沿用。其他倍率保留原邊角，平鋪內側。
func (g *game) zoomMinimapBox() image.Rectangle {
	x, y := g.mapOrigin()
	z := g.mapZoomScale()
	left := strategyMinimapX + int(math.Floor(x/32)) - 2
	top := strategyMinimapY + int(math.Floor(y/32))
	w := int(math.Round(strategyMapW/(32*z))) + 4
	h := max(4, int(math.Round(strategyMapH/(32*z))))
	return image.Rect(left, top, left+w, top+h)
}

func viewBoxSourceIndex(p, size, original int) int {
	if p < 2 {
		return p
	}
	if p >= size-2 {
		return original - (size - p)
	}
	return 2 + (p-2)%(original-4)
}
