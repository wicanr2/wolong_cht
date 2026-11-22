package isoview

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

const waterPaletteIndex = 3

var waterlineRaw = [...]int{50, 57, 65, 66, 67, 68, 69, 74, 82, 90, 92, 93, 97, 98, 99, 100, 101, 106, 107, 108, 109}

// DOS/V set 2, waterlineRaw order. Each decoded source index is one byte.
const waterlineSignature = "a8e4386a36f86ebe49726d95f12163904593e0701fbb8664f351b79c05a42e0b"

// Water-colored ground, bridges and banks in each verified DOS/V land atlas.
var landWaterStyles = [...]struct {
	rawIDs    []int
	signature string
}{
	{[]int{12, 13, 14, 15, 20, 24, 25, 26, 27, 33, 37, 41, 45, 48, 57, 58, 60, 64, 65, 67, 68, 72, 73, 75, 80, 82, 83, 155, 163, 171},
		"7d0b373598d7d73456326adfe65735049be47613ce635be8e2349f842286a2f6"},
	{[]int{13, 14, 15, 24, 25, 27, 28, 29, 30, 31, 33, 34, 41, 45, 48, 49, 56, 65, 90, 91, 92, 93, 96, 99, 105, 106, 107, 108, 109, 153},
		"ef27ae0973889db35264a0410e17387a8fde50696e826f141b24299af8f3fdfc"},
	{[]int{11, 12, 13, 14, 32, 35, 36, 37, 41, 43, 49, 50, 57, 65, 66, 67, 68, 69, 72, 74, 82, 86, 90, 92, 94, 99, 100, 107, 108, 145, 147, 168},
		"b971328bf2c7da3bc965d896569d2fbdf69eca3f4bcc8632559291af82d5c307"},
}

type waterBandKey struct {
	raw   int
	band  int
	phase image.Point
}

type battleMaterialSource struct {
	tile       *battle.SubTile
	bands      uint8
	waterBands uint8
	wallBands  uint8
	wallCap    bool
	ground     groundSurface
	ladder     bool
	footGround *battle.SubTile
}

type battleMaterials struct {
	water          image.Image
	wallFace       image.Image
	wallTop        image.Image
	palette        [16]color.RGBA
	sources        [battle.NumSubTiles]battleMaterialSource
	bands          map[waterBandKey]*ebiten.Image
	origin         image.Point
	ground         [groundSurfaceCount]image.Image
	wood           image.Image
	foliage        image.Image
	slope          image.Image
	roof           image.Image
	ship           image.Image
	ladderRevision int

	wallCapColors    uint16
	wallTextureParts uint32
}

func (v *View) waterAvailable(pack *library.HDBattle) bool {
	if v.lib == nil || !pack.HasWaterMaterial() {
		return false
	}
	if v.field >= 209 && v.field <= 212 {
		return pack.WaterAvailable(v.field, v.set)
	}
	return uint(v.set) < uint(len(landWaterStyles))
}

func (v *View) waterSources(pack *library.HDBattle) []int {
	if !v.waterAvailable(pack) {
		return nil
	}
	if v.field < 209 || v.field > 212 {
		style := landWaterStyles[v.set]
		used := v.terrainSourcesUsed()
		var selected []int
		for _, raw := range style.rawIDs {
			if used[raw] && !pack.HasTile(v.set, raw) {
				selected = append(selected, raw)
			}
		}
		if len(selected) == 0 || !v.matchesTerrainSources(style.rawIDs, style.signature) {
			return nil
		}
		return selected
	}

	sub := v.lib.SubTile(v.set, library.HDBattleWaterRaw)
	if sub == nil {
		return nil
	}
	count := 0
	for _, pixel := range sub.Pix {
		switch pixel {
		case waterPaletteIndex:
			count++
		case battle.Transparent:
		default:
			return nil
		}
	}
	if count != len(sub.Pix)/2 {
		return nil
	}
	selected := []int{library.HDBattleWaterRaw}
	if v.matchesTerrainSources(waterlineRaw[:], waterlineSignature) {
		for _, raw := range waterlineRaw {
			if !pack.HasTile(v.set, raw) {
				selected = append(selected, raw)
			}
		}
	}
	return selected
}

func (v *View) loadWater(pack *library.HDBattle) {
	selected := v.waterSources(pack)
	if len(selected) == 0 {
		return
	}
	texture, err := pack.Water()
	if err != nil {
		log.Printf("HD water material fallback: %v", err)
		return
	}
	v.materials.water, v.materials.palette = texture, v.pal
	if v.materials.bands == nil {
		v.materials.bands = make(map[waterBandKey]*ebiten.Image)
	}
	for _, raw := range selected {
		v.materials.addSource(raw, v.lib.SubTile(v.set, raw))
	}
}

func (v *View) matchesTerrainSources(rawIDs []int, signature string) bool {
	hash := sha256.New()
	var pixels [battle.SubTileW * battle.SubTileH]byte
	for _, raw := range rawIDs {
		sub := v.lib.SubTile(v.set, raw)
		if sub == nil {
			return false
		}
		for i, index := range sub.Pix {
			pixels[i] = byte(index)
		}
		hash.Write(pixels[:])
	}
	return hex.EncodeToString(hash.Sum(nil)) == signature
}

func (w *battleMaterials) addSource(raw int, sub *battle.SubTile) {
	var bands uint8
	for i, pixel := range sub.Pix {
		if pixel == waterPaletteIndex {
			bands |= 1 << (i / (battle.SubTileW * displayBandRows))
		}
	}
	if bands != 0 {
		source := &w.sources[raw]
		source.tile = sub
		source.waterBands = bands
		source.bands |= bands
	}
}

func (w *battleMaterials) reset() {
	for _, band := range w.bands {
		band.Dispose()
	}
	*w = battleMaterials{}
}

func (w *battleMaterials) enabled(raw, scale int) bool {
	return scale == library.HDBattleScale &&
		uint(raw) < uint(len(w.sources)) && w.sources[raw].tile != nil
}

func (w *battleMaterials) active() bool {
	return w.water != nil || w.wallFace != nil || w.ground[groundGrass] != nil ||
		w.wood != nil || w.foliage != nil || w.slope != nil || w.roof != nil || w.ship != nil
}

func repeatWater(position, period int) int {
	remainder := position % period
	if remainder < 0 {
		remainder += period
	}
	return remainder
}

// drawBand samples the position this band will occupy after unfolding.
// Adjacent cells share texture coordinates even when their source bands differ.
func (w *battleMaterials) drawBand(dst *ebiten.Image, origin image.Point, raw, srcY, dstY int) bool {
	source := w.sources[raw]
	if source.bands&(1<<(srcY/displayBandRows)) == 0 {
		return false
	}
	const scale = library.HDBattleScale
	x, y := displayUnfoldDestination(dstY / displayBandRows)
	phase := waterPhase(origin.Add(image.Pt(x, y)))
	key := waterBandKey{raw: raw, band: srcY, phase: phase}
	band, found := w.bands[key]
	if !found {
		rgba := image.NewRGBA(image.Rect(0, 0, battle.SubTileW*scale, displayBandRows*scale))
		w.paintBand(rgba, raw, phase, srcY, 0)
		band = ebiten.NewImageFromImage(rgba)
		w.bands[key] = band
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(dstY*scale))
	dst.DrawImage(band, op)
	return true
}

func waterPhase(origin image.Point) image.Point {
	const scale = library.HDBattleScale
	return image.Pt(
		repeatWater(origin.X*scale, library.HDBattleWaterWidth),
		repeatWater(origin.Y*scale, library.HDBattleWaterHeight),
	)
}

func (w *battleMaterials) drawFull(dst *ebiten.Image, origin image.Point, raw int) {
	const scale = library.HDBattleScale
	key := waterBandKey{raw: raw, band: -1, phase: waterPhase(origin)}
	full, found := w.bands[key]
	if !found {
		rgba := image.NewRGBA(image.Rect(0, 0, battle.SubTileW*scale, battle.SubTileH*scale))
		for part := range 4 {
			sourceY := part * displayBandRows
			x, y := displayUnfoldDestination(part)
			w.paintBand(rgba, raw, waterPhase(origin.Add(image.Pt(x, y))), sourceY, sourceY)
		}
		full = ebiten.NewImageFromImage(rgba)
		w.bands[key] = full
	}
	dst.DrawImage(full, nil)
}

func (w *battleMaterials) paintBand(dst *image.RGBA, raw int, phase image.Point, srcY, dstY int) {
	const scale = library.HDBattleScale
	source := w.sources[raw]
	band := uint8(1 << (srcY / displayBandRows))
	wall := source.wallBands&band != 0
	water := source.waterBands&band != 0
	for dy := 0; dy < displayBandRows*scale; dy++ {
		for dx := 0; dx < battle.SubTileW*scale; dx++ {
			sx, sy := dx/scale, srcY+dy/scale
			index := source.tile.At(sx, sy)
			pos := phase.Add(image.Pt(dx, dy))
			y := dstY*scale + dy
			switch {
			case index == battle.Transparent:
				continue
			case water && index == waterPaletteIndex:
				dst.Set(dx, y, w.water.At(
					pos.X%library.HDBattleWaterWidth,
					pos.Y%library.HDBattleWaterHeight,
				))
			case source.ladder && source.bands&band != 0 &&
				(source.footGround == nil || source.footGround.At(sx, sy) != index):
				dst.SetRGBA(dx, y, w.ladderColor(index, pos))
			case source.ground != groundNone:
				dst.SetRGBA(dx, y, w.groundColor(source.ground, index, pos))
			case wall && w.isWallColor(source.wallCap, srcY, index):
				dst.SetRGBA(dx, y, w.wallColor(source.wallCap, srcY, index, pos))
			default:
				dst.SetRGBA(dx, y, w.palette[index&15])
			}
		}
	}
}
