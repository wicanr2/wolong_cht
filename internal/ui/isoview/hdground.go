package isoview

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

type groundSurface uint8

const (
	groundNone groundSurface = iota
	groundGrass
	groundRoad
	groundSurfaceCount
	woodSlantDown
	woodSlantUp
	groundFoliage
	groundSlope
	groundRoof
	groundShip
)

type groundStyle struct {
	rawIDs     []int // Grass sources followed by roads, in fingerprint order.
	grassCount int
	signature  string
}

var groundStyles = [...]groundStyle{
	{[]int{112, 128, 136, 172, 174, 182, 183, 129, 144}, 7, "226f657703ccf16e78b8e406b13bbe5e61d68516dff63954b612db702879a49e"},
	{[]int{112, 137, 145, 156, 171, 174, 144}, 6, "d3a1b8f36f95b74ef4c1d279fa161289bf6ec9623a27ff788d323db1fd9d7e9e"},
	{[]int{117, 118, 119, 153, 161, 167, 171}, 6, "c45c7adcb32cb8fe4d5e2f06f8189111a09330d272ba817510496bd92208937f"},
}

type terrainMaterialStyle struct {
	rawIDs    []int
	surface   groundSurface
	signature string
}

var woodStyles = [...]terrainMaterialStyle{
	{[]int{95, 87}, woodSlantDown, "b001c9b0f94506b6ac8fedfc773128b03d018b0a21667bbfced189c098362ebd"},
	{[]int{120}, woodSlantUp, "8615c7fc2ca92f958d4b23106bb56a134d13a4cf545eff00e305745c5173918a"},
	{[]int{15, 47, 55}, woodSlantDown, "9c2918fc571edd41ef3c698be88288978dc5ef4bbd9873ce4d333dda40f2dde3"},
}

var foliageStyles = [...]terrainMaterialStyle{
	{[]int{38, 46, 53, 69}, groundFoliage, "8eaa47e149d7c520d3a1fb553ea6c2b34f71414defef1fd14afef44ccd7ed089"},
	{[]int{38, 46, 53, 61}, groundFoliage, "74924a6f21851cb79a8780defa2a8b08375173a9de3b4a200c155e773e202822"},
	{[]int{38, 46}, groundFoliage, "2a039f04740a96dca4aa4998e33da979ae55af982ddb0c4aac4c9d4db7347778"},
}

var slopeStyle = terrainMaterialStyle{
	[]int{27, 97, 98, 100, 105, 106, 108, 160},
	groundSlope,
	"583b60866c38b35a09ffdc1840e4f4ff4d85826e7e8ee0b4f668464ac37b0d87",
}

var roofStyles = [...]terrainMaterialStyle{
	{
		[]int{77, 78, 79, 88, 89, 90, 91, 92, 101, 107, 108, 109, 110, 111, 116, 118, 126, 130, 135, 138, 142, 143, 151},
		groundRoof, "a72ef89ef370c8c28c51f6e547ab25cce053d4a6b05777c7d211063eb2674f66",
	},
	{
		[]int{111, 119, 124, 127, 132, 134, 135, 142, 143, 151, 159},
		groundRoof, "5afca898ea96449e29a8c43dbe755197711d95cec93d44932a4e4f113a12a817",
	},
	{
		[]int{133, 141, 143, 149, 151, 154, 157},
		groundRoof, "265089dbf1d12c68f949001aa795d261646a490cfa897ef8fe84d2663fb9a634",
	},
}

var shipStyle = terrainMaterialStyle{
	[]int{
		7, 9, 15, 17, 18, 23, 39, 51, 52, 53, 54, 57, 58, 59, 60, 61, 65, 66, 67, 68, 69, 74, 76, 77, 80,
		82, 84, 85, 88, 90, 92, 93, 96, 98, 99, 100, 101, 106, 107, 108, 109, 112, 120, 124, 125, 127, 128,
		136, 142, 146, 162, 170, 183, 188, 189,
	},
	groundShip, "57201c5eac5b1409a863bc4d35db7a2f0251caad08740e6a57beb35322abe240",
}

func (v *View) groundAvailable(pack *library.HDBattle) bool {
	return pack.GroundAvailable() && v.lib != nil && uint(v.set) < uint(len(groundStyles)) &&
		(v.field < 209 || v.field > 212)
}

func (v *View) woodAvailable(pack *library.HDBattle) bool {
	return pack.WoodAvailable() && v.lib != nil && uint(v.set) < uint(len(woodStyles)) &&
		(v.field < 209 || v.field > 212)
}

func (v *View) foliageAvailable(pack *library.HDBattle) bool {
	return pack.FoliageAvailable() && v.lib != nil && uint(v.set) < uint(len(foliageStyles)) &&
		(v.field < 209 || v.field > 212)
}

func (v *View) slopeAvailable(pack *library.HDBattle) bool {
	return pack.SlopeAvailable() && v.lib != nil && v.set == 1
}

func (v *View) roofAvailable(pack *library.HDBattle) bool {
	siege := v.field >= 0 && battle.Category(v.field) == 0
	ship := v.field >= 209 && v.field <= 212
	return pack.RoofAvailable() && v.lib != nil && uint(v.set) < uint(len(roofStyles)) &&
		(siege || ship)
}

func (v *View) shipAvailable(pack *library.HDBattle) bool {
	return pack.ShipAvailable() && v.lib != nil && v.set == library.HDBattleWaterSet && v.field >= 209 && v.field <= 212
}

func (v *View) terrainSourcesUsed() (used [256]bool) {
	for _, row := range v.subs {
		for _, stack := range row {
			for _, raw := range stack {
				used[raw] = true
			}
		}
	}
	return used
}

func (v *View) loadGround(pack *library.HDBattle) {
	if !v.groundAvailable(pack) {
		return
	}
	style := groundStyles[v.set]
	if !v.matchesTerrainSources(style.rawIDs, style.signature) {
		return
	}
	used := v.terrainSourcesUsed()
	var selected []int
	for i, raw := range style.rawIDs {
		if used[raw] && !pack.HasTile(v.set, raw) && v.materials.sources[raw].tile == nil {
			selected = append(selected, i)
		}
	}
	if len(selected) == 0 {
		return
	}
	grass, road, err := pack.Ground()
	if err != nil {
		log.Printf("HD ground material fallback: %v", err)
		return
	}
	v.materials.ground[groundGrass], v.materials.ground[groundRoad] = grass, road
	v.materials.palette = v.pal
	if v.materials.bands == nil {
		v.materials.bands = make(map[waterBandKey]*ebiten.Image)
	}
	for _, i := range selected {
		raw := style.rawIDs[i]
		surface := groundGrass
		if i >= style.grassCount {
			surface = groundRoad
		}
		v.materials.sources[raw] = battleMaterialSource{tile: v.lib.SubTile(v.set, raw), bands: 0b1111, ground: surface}
	}
}

func (v *View) loadWood(pack *library.HDBattle) {
	if v.woodAvailable(pack) {
		v.materials.wood = v.loadTerrainMaterial(pack, woodStyles[v.set], "wood", pack.Wood)
	}
}

func (v *View) loadFoliage(pack *library.HDBattle) {
	if v.foliageAvailable(pack) {
		v.materials.foliage = v.loadTerrainMaterial(pack, foliageStyles[v.set], "foliage", pack.Foliage)
	}
}

func (v *View) loadSlope(pack *library.HDBattle) {
	if v.slopeAvailable(pack) {
		v.materials.slope = v.loadTerrainMaterial(pack, slopeStyle, "slope", pack.Slope)
	}
}

func (v *View) loadRoof(pack *library.HDBattle) {
	if v.roofAvailable(pack) {
		v.materials.roof = v.loadTerrainMaterial(pack, roofStyles[v.set], "roof", pack.Roof)
	}
}

func (v *View) loadShip(pack *library.HDBattle) {
	if v.shipAvailable(pack) {
		v.materials.ship = v.loadTerrainMaterial(pack, shipStyle, "ship", pack.Ship)
	}
}

func (v *View) loadTerrainMaterial(pack *library.HDBattle, style terrainMaterialStyle, name string, read func() (image.Image, error)) image.Image {
	if !v.matchesTerrainSources(style.rawIDs, style.signature) {
		return nil
	}
	used := v.terrainSourcesUsed()
	var selected []int
	for _, raw := range style.rawIDs {
		if used[raw] && !pack.HasTile(v.set, raw) && v.materials.sources[raw].tile == nil {
			selected = append(selected, raw)
		}
	}
	if len(selected) == 0 {
		return nil
	}
	texture, err := read()
	if err != nil {
		log.Printf("HD %s material fallback: %v", name, err)
		return nil
	}
	v.materials.palette = v.pal
	if v.materials.bands == nil {
		v.materials.bands = make(map[waterBandKey]*ebiten.Image)
	}
	for _, raw := range selected {
		source := &v.materials.sources[raw]
		source.tile = v.lib.SubTile(v.set, raw)
		source.bands, source.ground = 0b1111, style.surface
	}
	return texture
}

func (w *battleMaterials) groundColor(surface groundSurface, index int, pos image.Point) color.RGBA {
	original := w.palette[index]
	parts := uint32(6)
	var texture image.Image
	switch surface {
	case groundGrass:
		if index != 4 && index != 5 && index != 14 {
			return original
		}
		texture = w.ground[groundGrass]
	case groundRoad:
		if index != 2 && index != 9 && index != 12 {
			return original
		}
		texture = w.ground[groundRoad]
		parts = 7
	case woodSlantDown, woodSlantUp:
		if index != 6 && index != 7 && index != 9 && index != 11 {
			return original
		}
		texture = w.wood
		parts = 5
		// Align the grain with the projected bridge planks or fence slats.
		if surface == woodSlantUp {
			pos.Y += pos.X / 2
		} else {
			pos.Y -= pos.X / 2
		}
	case groundFoliage:
		if index != 4 && index != 5 && index != 14 {
			return original
		}
		texture = w.foliage
		parts = 4
	case groundSlope:
		if index != 6 && index != 7 && index != 11 {
			return original
		}
		texture = w.slope
		parts = 4
	case groundRoof:
		if index != 8 {
			return original
		}
		texture = w.roof
		parts = 6
	case groundShip:
		switch index {
		case 6, 7, 9, 11, 12:
			texture = w.ship
			parts = 5
			pos.Y -= pos.X / 2
		default:
			return original
		}
	default:
		return original
	}
	r, g, b, _ := texture.At(repeatWater(pos.X, library.HDBattleWaterWidth), repeatWater(pos.Y, library.HDBattleWaterHeight)).RGBA()
	if surface == groundFoliage && index == 14 {
		return shadeSeasonalFoliage(original, r, g, b)
	}
	return color.RGBA{
		R: blendMaterial(original.R, r, parts),
		G: blendMaterial(original.G, g, parts),
		B: blendMaterial(original.B, b, parts),
		A: original.A,
	}
}

func shadeSeasonalFoliage(original color.RGBA, r, g, b uint32) color.RGBA {
	// Modulate brightness while keeping the palette's seasonal hue.
	luma := int((54*r + 183*g + 19*b) >> 16)
	factor := 256 + 3*(luma-108)/4
	if peak := int(max(original.R, original.G, original.B)); peak > 0 {
		factor = min(factor, 255*256/peak)
	}
	return color.RGBA{
		R: uint8((int(original.R)*factor + 128) / 256),
		G: uint8((int(original.G)*factor + 128) / 256),
		B: uint8((int(original.B)*factor + 128) / 256),
		A: original.A,
	}
}
