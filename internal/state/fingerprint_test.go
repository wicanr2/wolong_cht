package state

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

// 同一個 seed 跑同樣的 tick 數，兩次要得到同一個指紋。
// 這是 Android 里程碑 A 的判準，也是桌面端的決定性迴歸（docs/spec/69）。
func TestFingerprintIsDeterministic(t *testing.T) {
	run := func(seed, ticks int) string {
		w := load(t, 0)
		r := rng.NewFixed(seed)
		for i := 0; i < ticks; i++ {
			w.Tick(r)
		}
		return w.FingerprintHex()
	}
	a, b := run(7, 500), run(7, 500)
	if a != b {
		t.Fatalf("同一個 seed 跑兩次得到不同指紋：%s vs %s", a, b)
	}
	if c := run(8, 500); c == a {
		t.Fatalf("換了 seed 指紋卻一樣（%s）——指紋沒有吃到亂數", c)
	}
	if d := run(7, 400); d == a {
		t.Fatalf("少跑 100 tick 指紋卻一樣（%s）——指紋沒有吃到時間推進", d)
	}
}

// 同一個世界連算兩次要相同。抓的是「走 map 造成的不決定性」——
// 那種錯誤不會讓別的測試變紅，只會讓指紋在跨平台比對時無故不同。
func TestFingerprintIsStableForTheSameWorld(t *testing.T) {
	w := load(t, 0)
	first := w.Fingerprint()
	for i := 0; i < 8; i++ {
		if w.Fingerprint() != first {
			t.Fatalf("第 %d 次重算就變了", i+2)
		}
	}
}

// ⭐ 正對照：每一個進指紋的欄位各改一格，指紋一定要變。
//
// **這一條擋的是「漏掉某個欄位」**——漏掉不會讓任何測試變紅，
// 只會讓指紋安靜地失去偵測力（docs/spec/69 §5）。
func TestFingerprintCoversEveryRecordedField(t *testing.T) {
	cases := []struct {
		name string
		poke func(w *World)
	}{
		{"存檔位元組（時鐘）", func(w *World) { w.Clock.Day++ }},
		{"存檔位元組（存活勢力數）", func(w *World) { w.LivingFactions++ }},
		{"據點整備游標", func(w *World) { w.cityCursor++ }},
		{"事件游標", func(w *World) { w.eventCursor++ }},
		{"事件延遲", func(w *World) { w.eventDelay++ }},
		{"軍團游標", func(w *World) { w.corpsCursor++ }},
		{"據點偏好", func(w *World) { w.cityBias[3]++ }},
		{"災害 marker", func(w *World) { w.disasterMarkers[5] = economy.Fire }},
		{"災害 marker 等級", func(w *World) { w.disasterMarkerLevels[5]++ }},
		{"災害物件 active", func(w *World) { w.disasterObjects[2].active = true }},
		{"災害物件 timer", func(w *World) { w.disasterObjects[2].timer++ }},
		{"災害物件位置", func(w *World) { w.disasterObjects[2].x++ }},
		{"暴風區", func(w *World) { w.stormArea = &economy.StormArea{} }},
		{"待決狀態", func(w *World) { w.pending = &Pending{} }},
		{"亂數狀態", func(w *World) { w.rng = rng.NewFixed(99) }},
	}
	for _, c := range cases {
		base := load(t, 0)
		base.rng = rng.NewFixed(1)
		before := base.Fingerprint()
		c.poke(base)
		if base.Fingerprint() == before {
			t.Errorf("改了「%s」指紋卻沒變——那個欄位沒進指紋", c.name)
		}
	}
}

// 「問不到亂數狀態」與「亂數狀態剛好是 0, 0」不可以算出同一個指紋。
//
// 兩者的意思差很多：前者是**指紋少涵蓋了一塊**，後者是一個正常的狀態。
// 混在一起的話，一個沒有 State() 的假亂數會讓指紋看起來仍在運作。
func TestFingerprintDistinguishesMissingRandomSource(t *testing.T) {
	a := load(t, 0)
	a.rng = noState{} // 問不到
	b := load(t, 0)
	b.rng = zeroState{} // 問得到，答案是 0, 0
	if a.Fingerprint() == b.Fingerprint() {
		t.Fatal("「問不到亂數狀態」與「狀態是 0, 0」算出同一個指紋")
	}
}

// noState 不提供 State()；zeroState 提供但回 0, 0。
type noState struct{}

func (noState) Next() int { return 0 }

type zeroState struct{}

func (zeroState) Next() int              { return 0 }
func (zeroState) State() (byte, byte)    { return 0, 0 }
