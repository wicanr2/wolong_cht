package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

type hdBannerState struct {
	pack          *library.HDBanner
	native        *ebiten.Image
	failed, drawn bool
}

func (g *game) drawHDBanner(screen *ebiten.Image, bank int) bool {
	q := &g.hdBanner
	if g.hdFrame == nil || !g.hdDisplay || g.videoLCD || q.pack == nil || q.failed || q.pack.Bank != bank {
		return false
	}
	if q.native == nil {
		img, err := q.pack.Image()
		if err != nil {
			q.failed = true
			log.Printf("HD banner fallback: %v", err)
			return false
		}
		q.native = ebiten.NewImageFromImage(img)
	}
	g.hdFrame.flush(screen)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(screenW*g.hdFrame.scale/float64(q.pack.Content.Dx()), bannerH*g.hdFrame.scale/float64(q.pack.Content.Dy()))
	g.hdFrame.surface.DrawImage(q.native.SubImage(q.pack.Content).(*ebiten.Image), op)
	q.drawn = true
	return true
}
