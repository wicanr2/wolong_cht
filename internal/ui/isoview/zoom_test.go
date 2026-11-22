package isoview

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

func TestBattleZoomExpandsViewportAndResetsOriginalCamera(t *testing.T) {
	v := &View{}
	v.SetCameraWorld(36, 14)
	if _, _, shown := v.ScreenPos(0, 0, 30, 12, 0); shown {
		t.Fatal("test point starts outside the original view")
	}
	v.ZoomAt(50, NativeW/2, NativeH/2)
	if x, y, shown := v.ScreenPos(0, 0, 30, 12, 0); !shown || x < 0 || x >= NativeW || y < 0 || y >= NativeH {
		t.Fatalf("zoom out did not expose the point: (%d,%d) %v", x, y, shown)
	}
	col, row := v.cellOffset(30, 12-originalRowBias, 0)
	if x, y, shown := v.scenePosition(30, 12, 0); !shown || col*isoColPx != x || row*isoRowPx != y {
		t.Fatalf("terrain and unit separated after zoom: cell %d,%d, pixel %d,%d", col, row, x, y)
	}
	v.ZoomAt(100, NativeW/2, NativeH/2)
	if cols, rows := v.sourceUnits(); cols != battleDisplayAnchors || rows != battleDisplayScanRows {
		t.Fatalf("reset altered the original display grid: %d x %d", cols, rows)
	}
	if v.camCol != 50 || v.camRow != 21 {
		t.Fatalf("reset altered the original projected camera: %d, %d", v.camCol, v.camRow)
	}
	v.ZoomAt(200, 160, 96)
	if x, y, shown := v.ScreenPos(0, 0, 35, 25, 0); !shown || x != 160 || y != 96 {
		t.Fatalf("wheel anchor moved: (%d,%d) %v", x, y, shown)
	}
}

func TestZoomedBattleMinimapRefreshesUnitOnClick(t *testing.T) {
	v := &View{}
	v.SetCameraWorld(36, 14)
	b := &tactical.Battle{Frame: 65}
	b.Sides[0].Soldiers[0] = tactical.Soldier{X: 36, Y: 31, Alive: true}
	v.Observe(b)
	v.ZoomAt(50, NativeW/2, NativeH/2)
	unit := b.DrawnSoldier(0, 0)
	col, row := isoProjectOriginal(unit.X, unit.Y, unit.Z)
	oldCol, oldRow := col-v.camCol, row-v.camRow
	if displaySlotEmpty(v.retained[oldRow][oldCol][0][1]) {
		t.Fatal("precondition: unit has drawn at the old camera")
	}
	minimap := v.MinimapMarkers()
	v.SetCameraFromMiniMap(582, 151)
	newCol, newRow := col-v.camCol, row-v.camRow
	if newCol == oldCol && newRow == oldRow {
		t.Fatal("precondition: click must move the camera")
	}
	if !displaySlotEmpty(v.retained[oldRow][oldCol][0][1]) || displaySlotEmpty(v.retained[newRow][newCol][0][1]) {
		t.Fatal("first frame after minimap click still contains old unit position")
	}
	if b.Frame != 65 || v.ZoomPercent() != 50 || v.MinimapMarkers() != minimap {
		t.Fatal("camera move affected battle time, zoom or minimap markers")
	}
}
