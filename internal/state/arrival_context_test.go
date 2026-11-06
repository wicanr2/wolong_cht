package state

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// 原版106號據點的抵達分派讀到40號據點的佔用值。
// 42、43與64另覆蓋勢力起點、半筆勢力與表間區域。
func TestArrivalContextControlsStage1AndRNG(t *testing.T) {
	for _, node := range []int{106, 42, 43, 64} {
		for _, value := range []int{2, 3} {
			for _, stale := range []int{-1, runtimeCityBase + citySize} {
				t.Run(fmt.Sprintf("node%d/value%d/stale%d", node, value, stale), func(t *testing.T) {
					w := load(t, 0)
					f := w.AliveFactions()[1]
					i := aiCorps(t, w, f, node)
					c := &w.Corps[i]
					w.Cities[node].Owner = f
					w.Cities[node].Threatened = true
					w.Cities[node].Specific = false
					w.Cities[node].Occupancy = 1
					c.Stage, c.Timer, c.Interval, c.Men = 1, 1, 3, 600
					w.stage1DI = stale

					d := node * citySize
					switch {
					case d >= runtimeCityBase:
						w.Cities[(d-runtimeCityBase)/citySize].Occupancy = value
					case d%factionSize == 0 && d/factionSize < numFactions:
						w.Factions[d/factionSize].Generals = value
					default:
						w.raw[factionBase+d+0x18] = byte(value)
					}

					random, want := rng.NewFixed(7), rng.NewFixed(7)
					stage, timer := 1, 3
					if value == 3 {
						stage, timer = 2, want.Next()&7+1
					}
					w.tickOneCorps(i, 2, random)
					if c.Stage != stage || c.Timer != timer || !bytes.Equal(random.Raw(), want.Raw()) {
						t.Fatalf("Stage/Timer=%d/%d，預期%d/%d；raw RNG一致=%t",
							c.Stage, c.Timer, stage, timer, bytes.Equal(random.Raw(), want.Raw()))
					}
				})
			}
		}
	}
}
