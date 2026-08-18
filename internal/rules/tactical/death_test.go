package tactical

import "testing"

// 圖號 ＝ 84 ＋ 側 × 90 ＋ 兵種組 × 2 ＋ (計時 ≤ 2 ? 1 : 0)（docs/spec/68 §1.1）。
func TestDeathSpriteNumbers(t *testing.T) {
	cases := []struct {
		side   int
		kind   Kind
		frames int
		want   int
	}{
		{0, Cavalry, 4, 84},   // 騎馬與大將同一組
		{0, General, 3, 84},   //
		{0, Cavalry, 2, 85},   // 後兩幀換第二張
		{0, Cavalry, 1, 85},   //
		{0, Archer, 4, 86},    // 弓 ＋2
		{0, Infantry, 4, 88},  // 步 ＋4
		{1, Cavalry, 4, 174},  // 側 1 ＋90
		{1, Infantry, 1, 179}, //
	}
	for _, c := range cases {
		d := Death{Side: c.side, Kind: c.kind, FramesLeft: c.frames}
		if got := d.Sprite(); got != c.want {
			t.Errorf("側 %d 兵種 %d 計時 %d → 圖號 %d，預期 %d",
				c.side, c.kind, c.frames, got, c.want)
		}
	}
	// 原版的 raw 常數要對得回來：raw ＝ 192 ＋ 圖號 × 2。
	d0 := Death{Side: 0, Kind: Cavalry, FramesLeft: 4}
	if raw := 192 + d0.Sprite()*2; raw != 0x168 {
		t.Errorf("側 0 的 raw ＝ %#x，預期 0x168", raw)
	}
	d1 := Death{Side: 1, Kind: Cavalry, FramesLeft: 4}
	if raw := 192 + d1.Sprite()*2; raw != 0x21C {
		t.Errorf("側 1 的 raw ＝ %#x，預期 0x21C", raw)
	}
}

// 打死一個兵要留下四幀的倒地動畫，而且那四幀不進 Remaining()。
func TestDeathAnimationLastsFourFrames(t *testing.T) {
	b := &Battle{}
	b.Sides[0].Soldiers[0] = Soldier{Alive: true, HP: 1, Kind: Infantry, X: 10, Y: 20}
	s := &b.Sides[0].Soldiers[0]
	b.applyLethalForTest(s)
	deaths := b.Deaths()
	if len(deaths) != 1 {
		t.Fatalf("倒地清單有 %d 筆，預期 1", len(deaths))
	}
	if deaths[0].X != 10 || deaths[0].Y != 20 || deaths[0].Side != 0 {
		t.Fatalf("倒地的位置或側別不對：%+v", deaths[0])
	}
	if b.Sides[0].Remaining() != 0 {
		t.Fatal("倒地中的兵不該算在場上人數裡")
	}
	for i := 0; i < DeathFrames-1; i++ {
		b.stepDeaths()
		if len(b.Deaths()) != 1 {
			t.Fatalf("第 %d 幀就消失了", i+1)
		}
	}
	b.stepDeaths()
	if len(b.Deaths()) != 0 {
		t.Fatalf("第 %d 幀之後還在", DeathFrames)
	}
}
