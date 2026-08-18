package main

import (
	"os"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
)

// 打壞城壁或門之後，畫出來的地形要跟著換（docs/spec/66）。
//
// 原版把圖塊值改掉就重新展開那一格的堆疊（`sub_1B824` → `sub_1BB6D`），
// 而繪圖端每一幀都從同一個緩衝區重建；remake 的 `v.subs` 是進戰場那一幀
// 展開的靜態副本，少了這一步就會「走得過去但牆還立著」。
func TestBrokenStructureRepaintsTerrain(t *testing.T) {
	read := func(n string) []byte {
		b, err := os.ReadFile("../../workplace/orig/dosv/" + n)
		if err != nil {
			t.Skip("找不到原版 " + n + "，跳過")
		}
		return b
	}
	lib, err := battle.Parse(read("BATTLE.MAP"), read("BATTLE.MDL"), read("BATTLE.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	// 找一張攻城戰場，並在裡面找一格城壁或門。
	field, wx, wy := -1, -1, -1
	for n := 0; n < battle.NumFields && field < 0; n++ {
		if !lib.IsSiege(n) {
			continue
		}
		tiles := lib.Tiles(n)
		for y := 0; y < battle.Height && field < 0; y++ {
			for x := 0; x < battle.Width; x++ {
				if v := tiles[y][x]; v >= tactical.TileWallLo && v <= tactical.TileGateHi {
					field, wx, wy = n, x, y
					break
				}
			}
		}
	}
	if field < 0 {
		t.Fatal("214 張戰場裡找不到任何城壁或門")
	}

	f := tactical.NewFieldFromTileLayers(lib.Tiles(field), lib.Heights(field),
		lib.TileLayers(field), lib.GateX(field))
	b := &tactical.Battle{Field: f}
	v := &battleView{lib: lib, set: lib.TileSet(field), field: field,
		subs: lib.SubTiles(field)}

	before := append([]byte(nil), v.subs[wy][wx]...)
	if rev := f.Revision(); rev != 0 {
		t.Fatalf("還沒打壞任何東西，版本應該是 0，得到 %d", rev)
	}
	v.syncTiles(b)
	if got := v.subs[wy][wx]; string(got) != string(before) {
		t.Fatalf("沒有東西被打壞卻重建了地形：%v → %v", before, got)
	}

	delta := 0x10
	if f.Tile(wx, wy) >= tactical.TileGateLo {
		delta = 8
	}
	f.Retile(wx, wy, delta)
	if f.Revision() == 0 {
		t.Fatal("Retile 之後版本沒有變，繪圖層不會知道要重畫")
	}
	v.syncTiles(b)
	got := v.subs[wy][wx]
	if string(got) == string(before) {
		t.Fatalf("(%d,%d) 的城壁／門打壞了，畫面上的堆疊還是 %v", wx, wy, before)
	}
	want := lib.SubTilesFor(field, tilesFromField(f))[wy][wx]
	if string(got) != string(want) {
		t.Fatalf("重建出來的堆疊 %v，預期 %v", got, want)
	}
}

func tilesFromField(f *tactical.Field) [][]byte {
	out := make([][]byte, tactical.Height)
	for y := range out {
		out[y] = make([]byte, tactical.Width)
		for x := range out[y] {
			out[y][x] = f.Tile(x, y)
		}
	}
	return out
}
