package isoview

import (
	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"testing"
)

func TestBannerProducerKeepsBothNativeCells(t *testing.T) {
	for _, tc := range []struct {
		side, z, phase int
		raw            []int
	}{{0, 5, 0, []int{337, 336}}, {1, 5, 3, []int{523, 522}}, {0, 6, 0, []int{344}}, {1, 6, 3, []int{527}}} {
		banner := battle.Banner{X: 20, Y: 40, Z: tc.z, Side: tc.side, Phase: tc.phase}
		col, row := isoProject(banner.X, banner.Y, banner.Z)
		v := &View{banners: []battle.Banner{banner}, camCol: col - 15, camRow: row - 12}
		entries := v.buildDisplayList(nil, &tactical.Battle{})
		if len(entries) != len(tc.raw) {
			t.Fatalf("高度%d旗應登錄%d格，得到%d", tc.z, len(tc.raw), len(entries))
		}
		for i, e := range entries {
			if e.raw != tc.raw[i] || e.cellCol != 15 || e.cellRow != 12-i || e.layer != tc.z+i || e.lane != 0 {
				t.Fatalf("旗格%d不符原版1DC03：%+v", i, e)
			}
		}
	}
}

func TestBannerProducerClipsBeforeRegisteringUpperCell(t *testing.T) {
	for _, z := range []int{5, 6} {
		for _, tc := range []struct {
			col, row int
			visible  bool
		}{{-1, 12, false}, {0, 12, true}, {30, 12, true}, {31, 12, false}, {15, -1, false}, {15, 0, true}, {15, 23, true}, {15, 24, false}} {
			banner := battle.Banner{X: 20, Y: 40, Z: z}
			col, row := isoProject(banner.X, banner.Y, z)
			v := &View{banners: []battle.Banner{banner}, camCol: col - tc.col, camRow: row - tc.row}
			entries := v.buildDisplayList(nil, &tactical.Battle{})
			count := 0
			if tc.visible {
				count = 1
				if z != 6 && tc.row > 0 {
					count++
				}
			}
			if len(entries) != count {
				t.Fatalf("高度%d錨點%d,%d：格數%d，原版%d", z, tc.col, tc.row, len(entries), count)
			}
		}
	}
}
