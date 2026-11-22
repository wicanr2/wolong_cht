package isoview

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

type wallStyle struct {
	rawIDs       []int // Wall faces followed by caps, in fingerprint order.
	signature    string
	capCount     int
	capColors    uint16
	textureParts uint32 // Out of ten.
}

// DOS/V decoded source fingerprints in rawIDs order, one byte per pixel.
var wallStyles = [...]wallStyle{
	{[]int{28, 29, 30, 31, 125}, "fc3cfce321a7fd37eab31ebd7cbe4e4a59355c79b5caea7f46c552e2f4e57a75", 1, 1<<7 | 1<<9 | 1<<11 | 1<<12, 6},
	{[]int{1, 2, 3, 8, 9, 10, 11, 114, 116, 117}, "a9b7d7f88eea34f31281daf537c57507d88e7759832a4ce4c94929fd22f4012f", 3, 1<<7 | 1<<9 | 1<<11 | 1<<12, 3},
	{[]int{19, 20, 21, 22, 166}, "b562ca15b1d5ac9e6a54852ab57abd051dd78ad17b05c8e0aa08edbed0138d7b", 1, 1<<6 | 1<<7 | 1<<9, 3},
}

func (v *View) wallAvailable(pack *library.HDBattle) bool {
	return pack.WallAvailable() && v.lib != nil && uint(v.set) < uint(len(wallStyles)) && v.lib.IsSiege(v.field)
}

func (v *View) loadWall(pack *library.HDBattle) {
	if !v.wallAvailable(pack) {
		return
	}
	style := wallStyles[v.set]
	if !v.matchesTerrainSources(style.rawIDs, style.signature) {
		return
	}
	used := v.terrainSourcesUsed()
	var selected []int
	for i, raw := range style.rawIDs {
		if used[raw] && !pack.HasTile(v.set, raw) {
			selected = append(selected, i)
		}
	}
	if len(selected) == 0 {
		return
	}
	cap, face, err := pack.Wall()
	if err != nil {
		log.Printf("HD wall material fallback: %v", err)
		return
	}
	if cap == nil || face == nil {
		return
	}
	v.materials = battleMaterials{
		wallFace: face, wallTop: cap, palette: v.pal, bands: make(map[waterBandKey]*ebiten.Image),
		wallCapColors: style.capColors, wallTextureParts: style.textureParts,
	}
	for _, i := range selected {
		raw := style.rawIDs[i]
		bands := uint8(0b1111)
		if v.set == 0 && raw == 31 {
			// Side bands of 31 also appear in the gate opening.
			bands = 0b0101
		}
		v.materials.sources[raw] = battleMaterialSource{
			tile: v.lib.SubTile(v.set, raw), bands: bands, wallBands: bands, wallCap: i >= len(style.rawIDs)-style.capCount,
		}
	}
}

func (w *battleMaterials) isWallColor(cap bool, srcY, index int) bool {
	if cap {
		return index >= 0 && index < 16 && w.wallCapColors&(1<<index) != 0
	}
	if srcY/displayBandRows%2 == 1 {
		return index == 6 || index == 7
	}
	return index == 6 || index == 7 || index == 9 || index == 11
}

func (w *battleMaterials) wallColor(cap bool, srcY, index int, position image.Point) color.RGBA {
	texture := w.wallFace
	x, y := position.X, position.Y
	shade := false
	switch {
	case cap:
		texture = w.wallTop
		x, y = x+2*y, y-x/2
	case srcY/displayBandRows%2 == 1:
		y += x / 2
		shade = true
	default:
		y -= x / 2
	}
	r, g, b, _ := texture.At(repeatWater(x, library.HDBattleWaterWidth),
		repeatWater(y, library.HDBattleWaterHeight)).RGBA()
	if shade {
		r, g, b = r*3/4, g*3/4, b*3/4
	}
	original := w.palette[index]
	return color.RGBA{
		R: blendMaterial(original.R, r, w.wallTextureParts),
		G: blendMaterial(original.G, g, w.wallTextureParts),
		B: blendMaterial(original.B, b, w.wallTextureParts),
		A: original.A,
	}
}

func blendMaterial(original uint8, material, parts uint32) uint8 {
	return uint8(((10-parts)*uint32(original) + parts*(material>>8) + 5) / 10)
}
