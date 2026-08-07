package economy

import "testing"

// seqRand 依序吐出給定的值，用完就從頭。
func seq(v ...int) Rand { return &fixedRand{seq: v} }

// 防災值 ≥ 64 完全免疫火災；上昇值存值 ≥ 64 完全免疫暴動。
// 亂數是 and 0x3F（0–63），所以門檻是 64 而不是說明書暗示的 1 或 0。
func TestDisasterImmunityThreshold(t *testing.T) {
	// 兩個參數都在門檻上，任何亂數序列都不該觸發。
	for _, r := range []int{0, 1, 23, 63, 255} {
		d := RollCityDisaster(DisasterImmunity, DisasterImmunity, seq(r))
		if d != NoDisaster {
			t.Errorf("亂數 %d：防災 64 上昇 64 卻觸發 %v", r, d)
		}
	}
	// 防災值降到 63 就有機會了。
	if d := RollCityDisaster(63, 100, seq(0, 63)); d != Fire {
		t.Errorf("防災 63 + 亂數(0,63) → %v, want 火災", d)
	}
}

// 開局的實際值：防災 100、上昇值存值 ~100 → 兩者都免疫。
func TestScenarioStartIsImmune(t *testing.T) {
	for _, r := range []int{0, 63, 255} {
		if d := RollCityDisaster(100, 104, seq(r)); d != NoDisaster {
			t.Errorf("開局值卻觸發 %v", d)
		}
	}
}

// 第一道閘：亂數 ≥ 24 就直接不發生，與防災值無關。
func TestDisasterGate(t *testing.T) {
	if d := RollCityDisaster(0, 0, seq(24)); d != NoDisaster {
		t.Errorf("亂數 24 不該過閘，卻觸發 %v", d)
	}
	if d := RollCityDisaster(0, 0, seq(23)); d != Fire {
		t.Errorf("亂數 23 應該過閘並觸發火災，得到 %v", d)
	}
}

// 火災優先：同一個據點同一個月不會同時發生兩種。
func TestFireTakesPrecedence(t *testing.T) {
	d := RollCityDisaster(0, 0, seq(0))
	if d != Fire {
		t.Errorf("防災與上昇值都是 0 時應該先判火災，得到 %v", d)
	}
}

// 防災值滿、上昇值歸零 → 只會有暴動。
func TestRiotOnly(t *testing.T) {
	// 序列：閘(0 過) 比較(0 < 100 → 不火災) 閘(0 過) 比較(0 >= 0 → 暴動)
	if d := RollCityDisaster(100, 0, seq(0)); d != Riot {
		t.Errorf("→ %v, want 暴動", d)
	}
}

// ⚠ 暴風雨的「靠海」判定比的是 X 的低位元組 —— 這是原版的 bug，刻意照抄。
func TestCoastalUsesLowByteOnly(t *testing.T) {
	if !IsCoastalForStorm(200) {
		t.Error("X=200 應判為靠海")
	}
	if IsCoastalForStorm(300) {
		t.Error("X=300 的低位元組是 44，原版會判成內陸（保留這個 bug）")
	}
	if IsCoastalForStorm(370) {
		t.Error("X=370（最東邊）的低位元組是 114，原版會判成內陸")
	}
	if !IsCoastalForStorm(255) {
		t.Error("X=255 應判為靠海")
	}
	if IsCoastalForStorm(191) {
		t.Error("X=191 應判為內陸")
	}
}

// 暴風雨範圍是以據點為中心的 11 × 11 格。
func TestStormArea(t *testing.T) {
	cities := make([]City, 192)
	cities[0] = City{X: 200, Y: 100} // 靠海 → 不需要額外那次擲骰
	// 序列：1（過第一關）、0（據點編號 0）
	s := RollStorm(cities, seq(1, 0))
	if s == nil {
		t.Fatal("靠海據點應該觸發暴風雨")
	}
	if s.MinX != 195 || s.MaxX != 205 || s.MinY != 95 || s.MaxY != 105 {
		t.Errorf("範圍 = %+v, want 195–205 × 95–105", *s)
	}
	if !s.Contains(200, 100) || s.Contains(206, 100) {
		t.Error("Contains 判定錯誤")
	}
}

// 座標小於 10 時不退 5 —— 原版是 `cmp ax,0Ah / jl` 不是 max(0, v−5)。
func TestStormAreaNearEdge(t *testing.T) {
	cities := make([]City, 192)
	cities[0] = City{X: 200, Y: 4}
	s := RollStorm(cities, seq(1, 0))
	if s == nil {
		t.Fatal("應該觸發")
	}
	if s.MinY != 4 || s.MaxY != 14 {
		t.Errorf("Y 範圍 = %d–%d, want 4–14（小於 10 不退）", s.MinY, s.MaxY)
	}
}

// 第一關沒過就不發生。
func TestStormFirstGate(t *testing.T) {
	cities := make([]City, 192)
	if s := RollStorm(cities, seq(0)); s != nil {
		t.Error("第一次擲骰為偶數時不該發生")
	}
}

// 據點編號 ≥ 192 直接不發生（原版 cmp al,0C0h / jnb ret）。
func TestStormIndexOutOfRange(t *testing.T) {
	cities := make([]City, 192)
	if s := RollStorm(cities, seq(1, 192)); s != nil {
		t.Error("據點編號 192 應該直接跳過")
	}
}

// 內陸據點要多過一關，機率減半。
func TestInlandNeedsExtraRoll(t *testing.T) {
	cities := make([]City, 192)
	cities[0] = City{X: 100, Y: 100} // 內陸
	if s := RollStorm(cities, seq(1, 0, 0)); s != nil {
		t.Error("內陸據點在額外擲骰為偶數時不該發生")
	}
	if s := RollStorm(cities, seq(1, 0, 1)); s == nil {
		t.Error("內陸據點在額外擲骰為奇數時應該發生")
	}
}
