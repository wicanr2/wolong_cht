package main

import (
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

// desktopPointer 保存捕捉模式的畫面游標；原始座標可無限延伸，畫面座標不可。
// 外推只消費本次位移，邊緣停留不產生任何位移（spec/149 §5）。
var desktopPointer struct {
	active           bool
	x, y, rawX, rawY int
	// 2：等捕捉後的 Draw；1：下一份輸入快照只用來重建基準。
	rebase int
}

func cursorPosition() (int, int) {
	if desktopPointer.active {
		return desktopPointer.x, desktopPointer.y
	}
	return ebiten.CursorPosition()
}

func pushPointer(pos, delta, limit int) (next, overflow int) {
	next = pos + delta
	if next < 0 {
		return 0, next
	}
	if next >= limit {
		return limit - 1, next - (limit - 1)
	}
	return next, 0
}

func (g *game) updateDesktopPointer() (int, int) {
	enabled := runtime.GOOS != "android" && runtime.GOOS != "ios" &&
		g.world != nil && g.launcher == nil && ebiten.IsFocused()
	p := &desktopPointer
	if !enabled {
		if p.active {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
			p.active = false
		}
		return 0, 0
	}
	x, y := ebiten.CursorPosition()
	if !p.active {
		p.x, _ = pushPointer(0, x, screenW)
		p.y, _ = pushPointer(0, y, screenH)
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		p.active = true
		p.rawX, p.rawY = x, y
		p.rebase = 2
		return 0, 0
	}
	return consumeDesktopPointer(x, y)
}

func consumeDesktopPointer(x, y int) (int, int) {
	p := &desktopPointer
	if p.rebase != 0 {
		p.rawX, p.rawY = x, y
		if p.rebase == 1 {
			p.rebase = 0
		}
		return 0, 0
	}
	dx, dy := x-p.rawX, y-p.rawY
	p.rawX, p.rawY = x, y
	var ox, oy int
	p.x, ox = pushPointer(p.x, dx, screenW)
	p.y, oy = pushPointer(p.y, dy, screenH)
	return ox, oy
}

// 捕捉切換後，跨過 Draw 才會取得下一份底層游標快照。
// 不能只略過一次 Update：低更新率時同一畫面可能補跑多次 Update。
func finishDesktopPointerFrame() {
	if desktopPointer.rebase == 2 {
		desktopPointer.rebase = 1
	}
}

func (g *game) scrollMapPixels(dx, dy int) {
	x := max(0, min((384-viewCols)*16, g.camX*16+g.camSubX+dx))
	y := max(0, min((256-viewRows)*16, g.camY*16+g.camSubY+dy))
	g.camX, g.camSubX = x/16, x%16
	g.camY, g.camSubY = y/16, y%16
}
