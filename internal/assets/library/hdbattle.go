package library

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
)

const HDBattleScale = 4

const (
	HDBattleWaterWidth  = 256
	HDBattleWaterHeight = 128
	HDBattleWaterSet    = 2
	HDBattleWaterRaw    = 32
)

type hdBattleKey struct {
	set int
	raw int
}

type hdBattleTile struct {
	Set  *int   `json:"set"`
	Raw  int    `json:"raw"`
	File string `json:"file"`
	SHA  string `json:"sha256"`
}

type hdBattleMaterial struct {
	File string `json:"file"`
	SHA  string `json:"sha256"`
}

type hdBattleWall struct {
	Cap  hdBattleMaterial `json:"cap"`
	Face hdBattleMaterial `json:"face"`
}

type hdBattleGround struct {
	Grass hdBattleMaterial `json:"grass"`
	Road  hdBattleMaterial `json:"road"`
}

// HDBattle indexes optional native-resolution replacements for the battle compositor.
type HDBattle struct {
	root    string
	tiles   map[hdBattleKey]hdBattleTile
	sets    [battle.NumTileSets]bool
	units   bool
	water   *hdBattleMaterial
	wall    *hdBattleWall
	ground  *hdBattleGround
	wood    *hdBattleMaterial
	foliage *hdBattleMaterial
	slope   *hdBattleMaterial
	roof    *hdBattleMaterial
	ship    *hdBattleMaterial
}

func LoadHDBattle(dir string) (*HDBattle, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	data, err := readHDAsset(root, "battle-selection.json", 1<<20)
	if err != nil {
		return nil, err
	}
	var manifest struct {
		Tiles   []hdBattleTile    `json:"tiles"`
		Water   *hdBattleMaterial `json:"water"`
		Wall    *hdBattleWall     `json:"wall"`
		Ground  *hdBattleGround   `json:"ground"`
		Wood    *hdBattleMaterial `json:"wood"`
		Foliage *hdBattleMaterial `json:"foliage"`
		Slope   *hdBattleMaterial `json:"slope"`
		Roof    *hdBattleMaterial `json:"roof"`
		Ship    *hdBattleMaterial `json:"ship"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	if len(manifest.Tiles) == 0 && manifest.Water == nil && manifest.Wall == nil && manifest.Ground == nil &&
		manifest.Wood == nil && manifest.Foliage == nil && manifest.Slope == nil && manifest.Roof == nil && manifest.Ship == nil {
		return nil, fmt.Errorf("empty HD battle selection")
	}
	pack := &HDBattle{root: root, tiles: make(map[hdBattleKey]hdBattleTile, len(manifest.Tiles))}
	if water := manifest.Water; water != nil {
		if err := validateHDPNGRecord(water.File, HDBattleWaterWidth, HDBattleWaterHeight, water.SHA); err != nil {
			return nil, fmt.Errorf("HD water material: %w", err)
		}
		pack.water = water
	}
	if wall := manifest.Wall; wall != nil {
		for _, material := range []hdBattleMaterial{wall.Cap, wall.Face} {
			if err := validateHDPNGRecord(material.File, HDBattleWaterWidth, HDBattleWaterHeight, material.SHA); err != nil {
				return nil, fmt.Errorf("HD wall material: %w", err)
			}
		}
		pack.wall = wall
	}
	if ground := manifest.Ground; ground != nil {
		for _, material := range []hdBattleMaterial{ground.Grass, ground.Road} {
			if err := validateHDPNGRecord(material.File, HDBattleWaterWidth, HDBattleWaterHeight, material.SHA); err != nil {
				return nil, fmt.Errorf("HD ground material: %w", err)
			}
		}
		pack.ground = ground
	}
	for _, entry := range []struct {
		name   string
		record *hdBattleMaterial
	}{
		{"wood", manifest.Wood},
		{"foliage", manifest.Foliage},
		{"slope", manifest.Slope},
		{"roof", manifest.Roof},
		{"ship", manifest.Ship},
	} {
		if material := entry.record; material != nil {
			if err := validateHDPNGRecord(material.File, HDBattleWaterWidth, HDBattleWaterHeight, material.SHA); err != nil {
				return nil, fmt.Errorf("HD %s material: %w", entry.name, err)
			}
		}
	}
	pack.wood, pack.foliage, pack.slope, pack.roof = manifest.Wood, manifest.Foliage, manifest.Slope, manifest.Roof
	pack.ship = manifest.Ship
	for _, tile := range manifest.Tiles {
		if tile.Set == nil {
			return nil, fmt.Errorf("missing HD battle set for raw %d", tile.Raw)
		}
		set := *tile.Set
		switch {
		case tile.Raw > 0 && tile.Raw < battle.CombinedSourceTerrainTiles && set >= 0 && set < len(pack.sets):
			pack.sets[set] = true
		case tile.Raw >= battle.CombinedSourceTerrainTiles && tile.Raw < battle.CombinedSourceTerrainTiles+battle.NumSprites && set == -1:
			pack.units = true
		default:
			return nil, fmt.Errorf("invalid HD battle tile set %d raw %d", set, tile.Raw)
		}
		key := hdBattleKey{set, tile.Raw}
		if _, ok := pack.tiles[key]; ok {
			return nil, fmt.Errorf("duplicate HD battle tile set %d raw %d", set, tile.Raw)
		}
		if err := validateHDPNGRecord(tile.File, HDBattleScale*battle.SubTileW, HDBattleScale*battle.SubTileH, tile.SHA); err != nil {
			return nil, err
		}
		pack.tiles[key] = tile
	}
	return pack, nil
}

func (p *HDBattle) HasWaterMaterial() bool {
	return p != nil && p.water != nil
}

// WaterAvailable reports whether the legacy naval surface may be used in this field.
// An explicit replacement tile retains its existing behavior.
func (p *HDBattle) WaterAvailable(field, set int) bool {
	if !p.HasWaterMaterial() || set != HDBattleWaterSet || field < 209 || field > 212 {
		return false
	}
	return !p.HasTile(set, HDBattleWaterRaw)
}

// HasTile reports whether an explicit tile has priority, even before its PNG is loaded.
func (p *HDBattle) HasTile(set, raw int) bool {
	if p == nil {
		return false
	}
	if raw >= battle.CombinedSourceTerrainTiles {
		set = -1
	}
	_, explicit := p.tiles[hdBattleKey{set, raw}]
	return explicit
}

func (p *HDBattle) Water() (image.Image, error) {
	if p == nil || p.water == nil {
		return nil, nil
	}
	return p.readMaterial(*p.water)
}

func (p *HDBattle) WallAvailable() bool {
	return p != nil && p.wall != nil
}

func (p *HDBattle) Wall() (cap, face image.Image, err error) {
	if p == nil || p.wall == nil {
		return nil, nil, nil
	}
	return p.readMaterialPair(p.wall.Cap, p.wall.Face)
}

func (p *HDBattle) GroundAvailable() bool {
	return p != nil && p.ground != nil
}

func (p *HDBattle) Ground() (grass, road image.Image, err error) {
	if !p.GroundAvailable() {
		return nil, nil, nil
	}
	return p.readMaterialPair(p.ground.Grass, p.ground.Road)
}

func (p *HDBattle) WoodAvailable() bool {
	return p != nil && p.wood != nil
}

func (p *HDBattle) Wood() (image.Image, error) {
	if !p.WoodAvailable() {
		return nil, nil
	}
	return p.readMaterial(*p.wood)
}

func (p *HDBattle) FoliageAvailable() bool {
	return p != nil && p.foliage != nil
}

func (p *HDBattle) Foliage() (image.Image, error) {
	if !p.FoliageAvailable() {
		return nil, nil
	}
	return p.readMaterial(*p.foliage)
}

func (p *HDBattle) SlopeAvailable() bool {
	return p != nil && p.slope != nil
}

func (p *HDBattle) Slope() (image.Image, error) {
	if !p.SlopeAvailable() {
		return nil, nil
	}
	return p.readMaterial(*p.slope)
}

func (p *HDBattle) RoofAvailable() bool {
	return p != nil && p.roof != nil
}

func (p *HDBattle) Roof() (image.Image, error) {
	if !p.RoofAvailable() {
		return nil, nil
	}
	return p.readMaterial(*p.roof)
}

func (p *HDBattle) ShipAvailable() bool {
	return p != nil && p.ship != nil
}

func (p *HDBattle) Ship() (image.Image, error) {
	if !p.ShipAvailable() {
		return nil, nil
	}
	return p.readMaterial(*p.ship)
}

func (p *HDBattle) readMaterialPair(first, second hdBattleMaterial) (image.Image, image.Image, error) {
	a, err := p.readMaterial(first)
	if err != nil {
		return nil, nil, err
	}
	b, err := p.readMaterial(second)
	if err != nil {
		return nil, nil, err
	}
	return a, b, nil
}

func (p *HDBattle) readMaterial(material hdBattleMaterial) (image.Image, error) {
	data, err := readHDPNG(p.root, material.File, HDBattleWaterWidth, HDBattleWaterHeight, material.SHA)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	for y := 0; y < HDBattleWaterHeight; y++ {
		for x := 0; x < HDBattleWaterWidth; x++ {
			if _, _, _, alpha := img.At(x, y).RGBA(); alpha != 0xffff {
				return nil, fmt.Errorf("HD material must be opaque: %s", material.File)
			}
		}
	}
	return img, nil
}

func (p *HDBattle) Available(set int, units bool) bool {
	return p != nil && (set >= 0 && set < len(p.sets) && p.sets[set] || units && p.units)
}

// Image returns nil for an unselected tile. Selected files are opened on demand.
func (p *HDBattle) Image(set, raw int) (image.Image, error) {
	if raw >= battle.CombinedSourceTerrainTiles {
		set = -1
	}
	tile, ok := p.tiles[hdBattleKey{set, raw}]
	if !ok {
		return nil, nil
	}
	data, err := readHDPNG(p.root, tile.File, HDBattleScale*battle.SubTileW, HDBattleScale*battle.SubTileH, tile.SHA)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(data))
}
