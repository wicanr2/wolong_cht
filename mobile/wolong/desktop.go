//go:build !android

package wolongmobile

import "github.com/hajimehoshi/ebiten/v2"

// RunDesktop is excluded from the Android binding. Keeping the game type
// private prevents gomobile from trying to expose Ebiten's Layout signature.
func RunDesktop() error {
	ebiten.SetWindowSize(1280, 992)
	ebiten.SetWindowTitle("臥龍傳 Remake — Android touch prototype")
	return ebiten.RunGame(newGame())
}
