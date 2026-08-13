//go:build android

package wolongmobile

import "github.com/hajimehoshi/ebiten/v2/mobile"

func init() {
	mobile.SetGame(newGame())
}

// Initialize gives the generated Java wrapper one stable, bindable entry
// point. The actual Ebiten game is registered by init above.
func Initialize() {}
