package battle

import (
	"image/color"
	"testing"
)

func TestRenderTacticalMinimapRawMappingAndYFlip(t *testing.T) {
	tiles := make([][]byte, Height)
	for y := range tiles {
		tiles[y] = make([]byte, Width)
	}
	attrs := make([]byte, TileDefs)
	attrs[0x07] = 0xab
	attrs[0x08] = 0x04
	tiles[3][2] = 0x07
	tiles[0][0] = 0x08

	m := RenderTacticalMinimap(tiles, attrs)
	if got := m.Pixels[122*TacticalMinimapWidth+6]; got != 0x0b {
		t.Fatalf("map (2,3) 的 attribute 低 4 bit = %#x，預期 0xb", got)
	}
	for y := 122; y < 124; y++ {
		for x := 6; x < 8; x++ {
			if got := m.Pixels[y*TacticalMinimapWidth+x]; got != 0x0b {
				t.Fatalf("map (2,3) 的 2×2 block 在 (%d,%d) = %#x", x, y, got)
			}
		}
	}
	// DX=mapX=0、BX=mapY=0 → x=0、y=126；這釘住 Y 反轉。
	if got := m.Pixels[126*TacticalMinimapWidth]; got != 0x04 {
		t.Fatalf("map (0,0) 沒落在反轉後的左下角：%#x", got)
	}
	if got := m.Pixels[0*TacticalMinimapWidth+126]; got != 0 {
		t.Fatalf("map (63,63) 的空白角落不應覆蓋錯誤位置：%#x", got)
	}
}

func TestRenderTacticalMinimapIsFull128By128(t *testing.T) {
	attrs := make([]byte, TileDefs)
	attrs[0] = 0x0d
	m := RenderTacticalMinimap(nil, attrs)
	if len(m.Pixels) != TacticalMinimapWidth*TacticalMinimapHeight {
		t.Fatalf("縮圖像素數 = %d，預期 %d", len(m.Pixels), TacticalMinimapWidth*TacticalMinimapHeight)
	}
	for i, got := range m.Pixels {
		if got != 0x0d {
			t.Fatalf("缺少 raw tile 的第 %d 像素 = %#x，預期 0xd", i, got)
		}
	}
}

func TestTacticalMinimapRGBAUsesPaletteIndex(t *testing.T) {
	var palette [16]color.RGBA
	palette[3] = color.RGBA{1, 2, 3, 255}
	m := TacticalMinimap{}
	m.Pixels[0] = 3
	img := m.RGBA(palette)
	if got := img.RGBAAt(0, 0); got != palette[3] {
		t.Fatalf("RGBA palette index = %#v，預期 %#v", got, palette[3])
	}
}

func TestTileAttributesReturnsRawMDLTableCopy(t *testing.T) {
	mapData := make([]byte, IndexSize+NumFields*FieldSize)
	mapData[0] = 1 // 第 0 張戰場使用圖塊組 1。
	mdl := make([]byte, MDLHeader+NumTileSets*TileSetSize)
	mdl[TileDefs+0x23] = 0x9e
	l, err := Parse(mapData, mdl, nil)
	if err != nil {
		t.Fatal(err)
	}
	attrs := l.TileAttributes(0)
	if len(attrs) != TileDefs || attrs[0x23] != 0x9e {
		t.Fatalf("raw attribute table 未保留：len=%d value=%#x", len(attrs), attrs[0x23])
	}
	attrs[0x23] = 0
	if l.TileAttributes(0)[0x23] != 0x9e {
		t.Fatal("TileAttributes 回傳內部 slice，修改外部資料不應影響 Library")
	}
}
