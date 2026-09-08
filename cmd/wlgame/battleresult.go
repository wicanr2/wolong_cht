package main

import (
	"fmt"
	"time"

	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

var battleResultDurations = [...]int{0, 3, 5, 10, 15, 30}

// battleResultTimer 從真正顯示頁面開始計時，不使用戰術 tick 或顯示更新數。
type battleResultTimer struct {
	battle  *tactical.Battle
	shownAt time.Time
}

func (r *battleResultTimer) show(b *tactical.Battle, now time.Time) {
	if r.battle != b || r.shownAt.IsZero() {
		r.battle, r.shownAt = b, now
	}
}

func (r *battleResultTimer) ready(b *tactical.Battle, seconds int, now time.Time, dismiss bool) bool {
	if seconds <= 0 {
		return true
	}
	if r.battle != b || r.shownAt.IsZero() {
		return false
	}
	return dismiss || now.Sub(r.shownAt) >= time.Duration(seconds)*time.Second
}

func (g *game) cycleBattleResultDuration() {
	for i, seconds := range battleResultDurations {
		if seconds == g.battleResultSeconds {
			g.battleResultSeconds = battleResultDurations[(i+1)%len(battleResultDurations)]
			return
		}
	}
	g.battleResultSeconds = 0
}

func battleResultValue(seconds int) string {
	if seconds <= 0 {
		return " 關 "
	}
	return fmt.Sprintf("%d秒", seconds)
}
