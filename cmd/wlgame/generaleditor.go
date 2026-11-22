package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wicanr2/wolong_cht/internal/state"
	"github.com/wicanr2/wolong_cht/internal/ui/chrome"
)

func desktopToolLabel(lang string, action desktopAction, assist bool) string {
	labels := desktopLabels["tools-"+lang]
	if action == desktopGeneralEditor {
		return labels[2]
	}
	if assist {
		return labels[1]
	}
	return labels[0]
}

const (
	editorPageSize = 12
	editorPrev     = 1
	editorNext     = 2
	editorCancel   = 3
	editorApply    = 4
	editorSelect   = 5
	editorAdjust   = editorSelect + editorPageSize
)

var (
	editorRect       = image.Rect(60, 44, 580, 368)
	editorPrevRect   = image.Rect(80, 324, 117, 348)
	editorNextRect   = image.Rect(219, 324, 256, 348)
	editorCancelRect = image.Rect(367, 324, 446, 348)
	editorApplyRect  = image.Rect(461, 324, 550, 348)
)

type generalEditor struct {
	ids                   []int
	page, selected, armed int
	draft                 state.GeneralAbilities
}

func (e *generalEditor) fields() [7]*int {
	return [7]*int{&e.draft.Martial, &e.draft.Command, &e.draft.Politics, &e.draft.Tactic,
		&e.draft.Aptitude[0], &e.draft.Aptitude[1], &e.draft.Aptitude[2]}
}

func (g *game) openGeneralEditor() {
	e := &generalEditor{selected: -1}
	lord := -1
	if p := g.world.Player; p >= 0 && p < len(g.world.Factions) {
		lord = g.world.Factions[p].Lord
	}
	initial := 0
	for id, gen := range g.world.Generals {
		if gen.Name == "" {
			continue
		}
		if id == lord {
			initial = len(e.ids)
		}
		e.ids = append(e.ids, id)
	}
	if len(e.ids) == 0 {
		return
	}
	e.page = initial / editorPageSize
	e.selectGeneral(g.world, e.ids[initial])
	g.generalEditor = e
}

func (e *generalEditor) selectGeneral(w *state.World, id int) {
	e.selected, e.draft = id, w.Generals[id].Abilities()
}

func (e *generalEditor) pageIDs() []int {
	return e.ids[e.page*editorPageSize : min((e.page+1)*editorPageSize, len(e.ids))]
}

func editorRow(row int) image.Rectangle {
	return image.Rect(80, 78+20*row, 256, 98+20*row)
}

func editorAdjustment(field, direction int) image.Rectangle {
	x, y := 438+74*direction, 98+30*field
	return image.Rect(x, y, x+38, y+24)
}

func editorBounds(field int) (int, int) {
	switch {
	case field < 3:
		return 1, 15
	case field == 3:
		return 0, 7
	default:
		return 0, 10
	}
}

func (e *generalEditor) adjustmentEnabled(field, direction int) bool {
	low, high := editorBounds(field)
	value := *e.fields()[field]
	if direction == 0 {
		return value > low
	}
	return value < high
}

func (e *generalEditor) hit(p image.Point, dirty bool) int {
	switch {
	case p.In(editorPrevRect) && e.page > 0:
		return editorPrev
	case p.In(editorNextRect) && (e.page+1)*editorPageSize < len(e.ids):
		return editorNext
	case p.In(editorCancelRect):
		return editorCancel
	case p.In(editorApplyRect) && dirty:
		return editorApply
	}
	for row := range e.pageIDs() {
		if p.In(editorRow(row)) {
			return editorSelect + row
		}
	}
	for field := range e.fields() {
		for direction := range 2 {
			if p.In(editorAdjustment(field, direction)) && e.adjustmentEnabled(field, direction) {
				return editorAdjust + field*2 + direction
			}
		}
	}
	return 0
}

func (g *game) updateGeneralEditor(x, y int, down, previous [3]bool) {
	e := g.generalEditor
	if previous[1] && !down[1] {
		g.generalEditor = nil
		return
	}
	hit := e.hit(image.Pt(x, y), e.draft != g.world.Generals[e.selected].Abilities())
	if down[0] && !previous[0] {
		e.armed = hit
	}
	if down[0] || !previous[0] {
		return
	}
	armed := e.armed
	e.armed = 0
	if armed == 0 || armed != hit {
		return
	}
	switch hit {
	case editorPrev:
		e.page--
	case editorNext:
		e.page++
	case editorCancel:
		g.generalEditor = nil
	case editorApply:
		if g.world.ApplyGeneralAbilities(e.selected, e.draft) {
			g.generalEditor = nil
		}
	default:
		if hit < editorAdjust {
			id := e.pageIDs()[hit-editorSelect]
			if id != e.selected {
				e.selectGeneral(g.world, id)
			}
			return
		}
		field, direction := (hit-editorAdjust)/2, (hit-editorAdjust)%2
		value := e.fields()[field]
		low, high := editorBounds(field)
		*value = clamp(*value+2*direction-1, low, high)
	}
}

func (g *game) drawGeneralEditor(dst *ebiten.Image) {
	e := g.generalEditor
	labels := desktopLabels["tools-"+desktopLocale()]
	g.chrome.Window(dst, editorRect.Min.X, editorRect.Min.Y, editorRect.Dx(), editorRect.Dy(), chrome.Blank)
	draw := func(label string, x, y int) {
		g.td.Draw(dst, label, x, y, chrome.Paper)
	}
	x, y := cursorPosition()
	p := image.Pt(x, y)
	button := func(r image.Rectangle, label string, enabled, selected bool) {
		ink := chrome.Paper
		if !enabled {
			ink = color.RGBA{126, 128, 132, 255}
		}
		if selected || enabled && p.In(r) {
			vector.DrawFilledRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()),
				color.RGBA{42, 55, 83, 255}, false)
		}
		vector.StrokeRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), 1, ink, false)
		g.td.Draw(dst, label, r.Min.X+(r.Dx()-g.td.Width(label))/2, r.Min.Y+(r.Dy()-16)/2, ink)
	}
	draw(labels[3], 80, 56)
	for row, id := range e.pageIDs() {
		button(editorRow(row), big5(g.world.Generals[id].Name), true, id == e.selected)
	}
	gen := g.world.Generals[e.selected]
	draw(big5(gen.Name), 284, 72)
	for i, value := range e.fields() {
		draw(labels[4+i], 284, 102+i*30)
		text := strconv.Itoa(*value)
		draw(text, 494-g.td.Width(text)/2, 102+i*30)
		button(editorAdjustment(i, 0), "-", e.adjustmentEnabled(i, 0), false)
		button(editorAdjustment(i, 1), "+", e.adjustmentEnabled(i, 1), false)
	}
	button(editorPrevRect, "<", e.page > 0, false)
	button(editorNextRect, ">", (e.page+1)*editorPageSize < len(e.ids), false)
	draw(fmt.Sprintf("%d/%d", e.page+1, (len(e.ids)+editorPageSize-1)/editorPageSize), 145, 328)
	button(editorCancelRect, labels[12], true, false)
	button(editorApplyRect, labels[11], e.draft != gen.Abilities(), false)
	if p.In(dst.Bounds()) && desktopPointer.active {
		drawArrowCursor(dst, p, chrome.Paper, color.RGBA{240, 210, 149, 255})
	}
}
