package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
)

func TestMonthlyStormRawGlobals(t *testing.T) {
	w := &World{raw: make([]byte, blockSize)}
	for _, area := range []*economy.StormArea{nil, {MinX: 29, MinY: 86, MaxX: 39, MaxY: 96}} {
		w.stormArea = area
		w.saveMonthlyStormGlobals()
		want := [4]int{0xFFF0, 0xFFF0, 400, 400}
		if area != nil {
			want = [4]int{29, 86, 39, 96}
		}
		for i, value := range want {
			if got := u16(w.Bytes(), 0x32+i*2); got != value {
				t.Fatalf("storm word %d: got %d want %d", i, got, value)
			}
		}
		if u16(w.RawBlock(), 0x32) != 0 {
			t.Fatal("original raw anchor mutated")
		}
		reloaded, err := LoadBlock(w.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		for i, value := range want {
			if u16(reloaded.Bytes(), 0x32+i*2) != value {
				t.Fatal("storm globals lost in original-format round trip")
			}
		}
	}
}
