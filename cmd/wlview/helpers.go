package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// pressed 是「這一幀剛按下」，不是「按著」。
// 用持續按下的話換頁會一次跳幾十張。
func pressed(k ebiten.Key) bool { return inpututil.IsKeyJustPressed(k) }
