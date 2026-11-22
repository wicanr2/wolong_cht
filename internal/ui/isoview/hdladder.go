package isoview

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

type ladderStyle struct {
	signature  string
	firstRaw   int
	footGround [2]int // Background under foot sources 178 and 181.
}

// The DOS/V gate stacks replace the feet, rails and tips as a unit.
// Set 2 has a second ladder using sources 176–178.
var ladderStyles = [...]ladderStyle{
	{"bd04ea015ae642aefc3246fae1f5dd9b7348fb733cc597ac2299fbd7d7e04383", 178, [2]int{129, 176}},
	{"0c2c1e99bf0ad360ecde0896f2745f30407196d323d2c7cb9953d3b9880ffa06", 178, [2]int{113, 176}},
	{"2b0c41de38ec830fcf01c0b7e1a84fc6ff45d41c0e536ffb286a9216418d9b1f", 176, [2]int{114, 114}},
}

func (v *View) loadLadders(pack *library.HDBattle) {
	style := ladderStyles[v.set]
	if !v.matchesTerrainSources([]int{176, 177, 178, 179, 180, 181, style.footGround[0]}, style.signature) {
		return
	}
	used := v.terrainSourcesUsed()
	for raw := style.firstRaw; raw <= 181; raw++ {
		if !used[raw] || pack.HasTile(v.set, raw) || v.materials.sources[raw].tile != nil {
			continue
		}
		if v.materials.wood == nil {
			wood, err := pack.Wood()
			if err != nil {
				log.Printf("HD ladder material fallback: %v", err)
				return
			}
			v.materials.wood = wood
		}
		var footGround *battle.SubTile
		bands := uint8(0b1111)
		if raw == 178 || raw == 181 {
			footGround = v.lib.SubTile(v.set, style.footGround[(raw-178)/3])
			bands = 0b1011 // The other band shows only the original ground.
		}
		v.materials.sources[raw] = battleMaterialSource{
			tile: v.lib.SubTile(v.set, raw), bands: bands, ladder: true, footGround: footGround,
		}
	}
	if v.materials.wood != nil {
		v.materials.palette = v.pal
		if v.materials.bands == nil {
			v.materials.bands = make(map[waterBandKey]*ebiten.Image)
		}
	}
}

func (w *battleMaterials) ladderColor(index int, pos image.Point) color.RGBA {
	original := w.palette[index]
	r, g, b, _ := w.wood.At(repeatWater(pos.Y, library.HDBattleWaterWidth),
		repeatWater(pos.X, library.HDBattleWaterHeight)).RGBA()
	switch index {
	case 6, 7, 9, 11:
		// Leave the contrast between the ladder posts and rungs visible.
		return color.RGBA{blendMaterial(original.R, r, 3), blendMaterial(original.G, g, 3),
			blendMaterial(original.B, b, 3), original.A}
	case 1, 2:
		// The cool palette variant keeps its original hue.
		luma := int((54*r + 183*g + 19*b) >> 16)
		factor := 256 + (luma-128)/3
		scale := func(c uint8) uint8 { return uint8(min(255, (int(c)*factor+128)/256)) }
		return color.RGBA{scale(original.R), scale(original.G), scale(original.B), original.A}
	default:
		return original
	}
}
