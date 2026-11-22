package isoview

import (
	"image"
	"math"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
)

func (v *View) ZoomPercent() int {
	if v == nil || v.zoom == 0 {
		return 100
	}
	return v.zoom
}

// ZoomAt changes the viewport around a point in the logical battlefield.
func (v *View) ZoomAt(percent, x, y int) {
	percent = max(50, min(200, percent))
	old := v.ZoomPercent()
	if old == percent {
		return
	}
	v.zoomX, v.zoomY = v.zoomOffset()
	change := 100/float64(old) - 100/float64(percent)
	v.zoomX += float64(x) * change
	v.zoomY += float64(y) * change
	v.zoom = percent
	if percent == 100 {
		v.zoomX, v.zoomY = 0, 0
	}
	v.applyCameraOrigin()
	v.compositeCache = nil
	v.refreshZoomedDrawings()
}

// Reproject the last drawn poses before the next tactical tick.
func (v *View) refreshZoomedDrawings() {
	if b := v.observed; b != nil {
		v.rebuildRetained(b)
		v.restoreDrawings(b, false)
	}
}

func (v *View) zoomOffset() (float64, float64) {
	if v.ZoomPercent() == 100 {
		return 0, 0
	}
	col, row := isoProjectOriginal(v.camWorldX, v.camWorldY, 0)
	factor := 100 / float64(v.ZoomPercent())
	x := float64(col * isoColPx)
	y := float64(row * isoRowPx)
	left := -float64(isoColPx * 2)
	top := -float64(battle.MaxStack * isoRowPx)
	right := float64(2*battle.Width*isoColPx) - NativeW*factor
	bottom := float64((battle.Width+2)*isoRowPx) - NativeH*factor
	return max(left, min(right, v.zoomX+x)) - x, max(top, min(bottom, v.zoomY+y)) - y
}

func (v *View) zoomOrigin() (int, int) {
	x, y := v.zoomOffset()
	return int(math.Floor(x / (2 * isoColPx))), int(math.Floor(y / isoRowPx))
}

func (v *View) visibleScene() image.Rectangle {
	x, y := v.zoomOffset()
	left, top := v.zoomOrigin()
	minX := int(math.Floor(x)) - left*2*isoColPx
	minY := int(math.Floor(y)) - top*isoRowPx
	percent := float64(v.ZoomPercent())
	w, h := int(math.Ceil(NativeW*100/percent)), int(math.Ceil(NativeH*100/percent))
	return image.Rect(minX, minY, minX+w, minY+h)
}

func (v *View) sourceUnits() (cols, rows int) {
	if v.ZoomPercent() == 100 {
		return battleDisplayAnchors, battleDisplayScanRows
	}
	r := v.visibleScene()
	// Reserve room for the crop offset within the first cell. The grid then
	// keeps the same dimensions while panning at a fixed zoom.
	return (r.Dx() + 4*isoColPx - 2) / (2 * isoColPx), (r.Dy() + 2*isoRowPx - 2) / isoRowPx
}
