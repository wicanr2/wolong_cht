package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

type hdLocationState struct {
	pack       *library.HDLocations
	cache      map[int]*ebiten.Image
	failed     map[int]bool
	lru        []int
	drawn      bool
	page, city int
}

func originalLocationCount(g *game) int {
	if g == nil || g.lib == nil {
		return 0
	}
	for _, e := range g.lib.Entries {
		if e.Label == "KYOGRF locations" {
			return e.Count
		}
	}
	return 0
}

func (q *hdLocationState) touch(page int) {
	for i, p := range q.lru {
		if p == page {
			q.lru = append(q.lru[:i], q.lru[i+1:]...)
			break
		}
	}
	q.lru = append(q.lru, page)
}

func (g *game) drawHDLocation(screen *ebiten.Image, page, bank int) bool {
	q := &g.hdLocation
	if !g.hdDisplay || g.hdFrame == nil || g.videoLCD || q.pack == nil || bank != q.pack.Bank || q.failed[page] {
		return false
	}
	native := q.cache[page]
	if native == nil {
		img, err := q.pack.Image(page)
		if err != nil {
			if q.failed == nil {
				q.failed = map[int]bool{}
			}
			q.failed[page] = true
			log.Printf("HD location fallback (%d): %v", page, err)
			return false
		}
		if q.cache == nil {
			q.cache = map[int]*ebiten.Image{}
		}
		if len(q.lru) >= 2 {
			old := q.lru[0]
			q.lru = q.lru[1:]
			q.cache[old].Dispose()
			delete(q.cache, old)
		}
		native = ebiten.NewImageFromImage(img)
		q.cache[page] = native
	}
	q.touch(page)
	g.hdFrame.nativeRect(screen, native, cityViewX, cityViewY, cityViewSize, cityViewSize)
	q.drawn, q.page, q.city = true, page, g.cityInfo.city
	return true
}
