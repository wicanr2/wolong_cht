package main

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// sub_13D91／sub_13DC9 對 byte_10D00 在 0 與 255 飽和；GUI 的直接反應與
// Session 理由路徑都必須共用這個邊界。
func TestAdjustTrustSaturatesOriginalByte(t *testing.T) {
	for _, tc := range []struct {
		name  string
		start int
		delta int
		want  int
	}{
		{"increase", 250, 20, 255},
		{"decrease", 5, -20, 0},
		{"middle", 100, -20, 80},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := &game{world: &state.World{Trust: tc.start}}
			g.adjustTrust(tc.delta)
			if g.world.Trust != tc.want {
				t.Fatalf("Trust %d %+d → %d, want %d",
					tc.start, tc.delta, g.world.Trust, tc.want)
			}
		})
	}
}
