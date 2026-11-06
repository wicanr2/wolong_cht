package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/ui/listwin"
)

// sub_18713：以總列數換算指標，扣半頁，最後夾在有效頁首範圍。
func listScrollTopAt(y int, track image.Rectangle, total, page int) int {
	if total <= page || page <= 0 || track.Dy() <= 0 {
		return 0
	}
	return clamp(clamp(y-track.Min.Y, 0, track.Dy())*total/track.Dy()-page/2, 0, total-page)
}

// sub_18607：偏移上取整，高度下取整（spec/274）。
func listScrollThumb(track image.Rectangle, top, total, page int) (y, height int) {
	if total <= page || page <= 0 {
		return track.Min.Y, track.Dy()
	}
	top = clamp(top, 0, total-page)
	return track.Min.Y + (top*track.Dy()+total-1)/total, track.Dy() * page / total
}

// 一次拖曳只有一個清單擁有者。正式Top在全部按鍵放開後才提交。
type listScrollDrag struct {
	list     *listwin.List
	launcher *launcherModel
	top      int
}

func (g *game) listScrollOwnsInput() bool {
	d := g.listDrag
	return d.list != nil && d.list == g.list && g.launcher == nil && d.list.Phase() == listwin.Browsing ||
		d.launcher != nil && d.launcher == g.launcher && d.launcher.phase == launcherSelectFaction && !d.launcher.factionSelected
}

func listScrollButtons() [3]bool {
	return [3]bool{
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight),
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle),
	}
}

func (g *game) updateListScrollDrag(y int, down [3]bool) bool {
	if !g.listScrollOwnsInput() {
		g.listDrag = listScrollDrag{}
		return false
	}
	d := &g.listDrag
	if d.list != nil {
		d.top = listScrollTopAt(y, listScrollTrackRect(), len(d.list.Rows), d.list.Height)
	} else {
		d.top = listScrollTopAt(y, factionListScrollTrackRect(), len(d.launcher.players), factionListRowsPerPage)
	}
	if down == [3]bool{} {
		if d.list != nil {
			g.scrollListTo(d.top)
		} else {
			d.launcher.factionTop = d.top
		}
		g.listDrag = listScrollDrag{}
	}
	return true // 包含放開那一幀，不能再命中列或取消。
}
