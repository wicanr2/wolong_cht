package state

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// 原版1170案例的各資料區、邊界與16-bit繞回；來源值2留守、3出擊。
func TestAIHoldContextRebuiltAtEveryArrival(t *testing.T) {
	for _, node := range []int{0, 4, 16, 20, 21, 22, 29, 30, 42, 63, 64, 106, 129, 191} {
		for _, hold := range []bool{false, true} {
			for _, value := range []int{2, 3} {
				for _, stale := range []int{-1, runtimeCityBase + citySize} {
					t.Run(fmt.Sprintf("node%d/hold%v/value%d/stale%d", node, hold, value, stale), func(t *testing.T) {
						w, i := aiContextFixture(t, node, hold, value)
						w.stage1DI = stale
						random, expected := rng.NewFixed(7), rng.NewFixed(7)
						stage, timer := 1, 3
						if value == 3 {
							stage, timer = 2, expected.Next()&7+1
						}
						w.tickOneCorps(i, 2, random)
						c := &w.Corps[i]
						if c.Stage != stage || c.Timer != timer || c.Heading != HeadingStill || !bytes.Equal(random.Raw(), expected.Raw()) {
							t.Fatalf("Stage/Timer/Heading=%d/%d/%d，預期%d/%d/4；RNG相同=%t", c.Stage, c.Timer, c.Heading, stage, timer, bytes.Equal(random.Raw(), expected.Raw()))
						}
					})
				}
			}
		}
	}
}

func aiContextFixture(t *testing.T, node int, hold bool, value int) (*World, int) {
	t.Helper()
	w := load(t, 0)
	f := w.AliveFactions()[1]
	i := aiCorps(t, w, f, node)
	c := &w.Corps[i]
	w.Cities[node].Owner, w.Cities[node].Threatened, w.Cities[node].Specific, w.Cities[node].Occupancy = f, true, false, 1
	c.Stage, c.Timer, c.Interval, c.Men, c.Heading = 1, 1, 3, 600, 0
	c.Ordered = node
	if hold {
		c.TargetNode = (node + 1) % numCities
	}
	d := node * citySize
	if hold {
		d = (node * aiHoldStride) & 0xffff
	}
	if d+24 < dosvWorldToIcons {
		b := w.Bytes()
		b[factionBase+d+24] = byte(value)
		w = loadBlock(b)
	}
	icon := bytes.Repeat([]byte{byte(value)}, dosvIconBankSize)
	w.SetAIHoldSources(icon, func(x, y int) byte { return byte(value) })
	return w, i
}

func TestAIHoldContextRebuiltAfterSaveRestore(t *testing.T) {
	for _, node := range []int{16, 21, 22, 30, 64, 106, 191} {
		for _, value := range []int{2, 3} {
			w, i := aiContextFixture(t, node, true, value)
			before := w.Bytes()
			restored, err := LoadBlock(before)
			if err != nil {
				t.Fatal(err)
			}
			if err := restored.Restore(w.TakeSnapshot()); err != nil {
				t.Fatal(err)
			}
			restored.SetAIHoldSources(bytes.Repeat([]byte{byte(value)}, dosvIconBankSize), func(x, y int) byte { return byte(value) })
			if !bytes.Equal(before, restored.Bytes()) {
				t.Fatal("重新掛來源改寫了存檔")
			}
			w.stage1DI, restored.stage1DI = 0, -1
			a, b := rng.NewFixed(19), rng.NewFixed(19)
			w.tickOneCorps(i, 2, a)
			restored.tickOneCorps(i, 2, b)
			if !bytes.Equal(w.Bytes(), restored.Bytes()) || !bytes.Equal(a.Raw(), b.Raw()) {
				t.Fatalf("node%d/value%d：讀檔後到站與未讀檔不同", node, value)
			}
		}
	}
}

func TestAIHoldContextReadsLiveGeneralBytes(t *testing.T) {
	w := load(t, 0)
	w.stage1DI = 0x5000 // 武將110的+18h是Timer，並非據點佔用。
	for _, value := range []int{2, 3} {
		w.Generals[110].Timer = value
		if got := w.residualField18(); got != value {
			t.Fatalf("Timer%d：讀到%d", value, got)
		}
	}
}

func TestAIHoldSourcesCopiedAndClearedWhenMissing(t *testing.T) {
	w := load(t, 0)
	icon := bytes.Repeat([]byte{2}, dosvIconBankSize)
	w.SetAIHoldSources(icon, func(x, y int) byte { return 2 })
	icon[21*aiHoldStride+24-dosvWorldToIcons] = 3
	w.stage1DI = 21 * aiHoldStride
	if got := w.residualField18(); got != 2 {
		t.Fatalf("保留了可變圖庫參考：%d", got)
	}
	w.SetAIHoldSources(icon[:len(icon)-1], nil)
	for _, n := range []int{21, 30} {
		w.stage1DI = n * aiHoldStride
		if got := w.residualField18(); got != aiCrowded+1 {
			t.Fatalf("缺來源時仍讀舊值：node%d=%d", n, got)
		}
	}
}
