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

func TestBothDamageEntriesInitializeDeathAndKeepPreviousPose(t *testing.T) {
	for _, projectile := range []bool{false, true} {
		b := newTestBattle(flatField())
		s := &b.Sides[0].Soldiers[1]
		*s = Soldier{Alive: true, HP: 1, Kind: Archer, Facing: East, PoseStep: 1,
			HitGeneral: true, TargetRefresh: true, Swapped: true, MoveFlag: true, HighTerrain: true, PathQueued: true}
		if projectile {
			b.projectileHit(s, 5)
		} else {
			b.applyHit(s, 5)
		}
		if s.Alive || s.HP != 0 || s.Stun != 4 || s.TargetRefresh || s.Swapped || s.MoveFlag || s.HighTerrain || !s.PathQueued {
			t.Fatalf("projectile%v：致命傷未依原版清旗標／設4拍", projectile)
		}
		d := b.Deaths()[0]
		if !d.Pending || d.BeforeFacing != East || d.BeforeStep != 1 || d.BeforeHurt || !d.BeforeHitGeneral {
			t.Fatal("首次倒地處理前須保留受傷前姿態")
		}
	}
}

func TestDeathTimingFollowsNativeUnitOrder(t *testing.T) {
	for _, late := range []bool{false, true} {
		for player := 0; player < 2; player++ {
			b := newTestBattle(flatField())
			b.SetPlayerSide(player)
			b.rng = &fixedRand{seq: []int{0}}
			b.Advantage = [2]Advantage{Even, Even}
			holdRetreatCountdownUnits(b)
			targetSide := 1 - player
			if late {
				targetSide = player
			}
			attackerSide := 1 - targetSide
			target := &b.Sides[targetSide].Soldiers[1]
			*target = Soldier{Alive: true, HP: 1, Kind: Archer, X: 20, Y: 32, StepX: 20, StepY: 32,
				GoalX: 20, GoalY: 32, Cmd: Holding, Next: Holding, Target: -1}
			attacker := &b.Sides[attackerSide].Soldiers[1]
			*attacker = Soldier{Alive: true, HP: 100, Power: 200, Kind: Infantry, X: 21, Y: 32,
				StepX: 20, StepY: 32, GoalX: 20, GoalY: 32, Cmd: Holding, Next: Holding, Target: -1}
			b.Step()
			if target.Alive || len(b.Deaths()) != 1 {
				t.Fatalf("late%v player%d：前提應由實際近戰擊倒目標", late, player)
			}
			want := 3
			if late {
				want = 4
			}
			d := b.Deaths()[0]
			if d.FramesLeft != want || target.Stun != want || d.Pending != late {
				t.Fatalf("late%v player%d：本拍已輪過才被擊倒，應保留4至下拍；實際%+v", late, player, d)
			}
		}
	}
}

func TestDeathCleanupRetainsCellsAndDefersReplacement(t *testing.T) {
	for side := 0; side < 2; side++ {
		b := newTestBattle(flatField())
		holdRetreatCountdownUnits(b)
		s := &b.Sides[side].Soldiers[1]
		s.X, s.Y, s.Z = 20, 32, 0
		b.Sides[side].Reserve[0] = 3
		b.initUnitCollision()
		b.applyLethalForTest(s)
		for frame := 1; frame <= 4; frame++ {
			b.Step()
			if s.Alive || b.Sides[side].Reserve[0] != 3 || b.Sides[side].Escaped[0] != 0 {
				t.Fatal("倒地與清除當拍不得補兵或記生還")
			}
			if s.Stun != 4-frame {
				t.Fatal("單兵計時須與倒地清單同步")
			}
			want := b.collisionSlot(side, 1)
			if frame == 4 {
				want = -1
			}
			for _, z := range []int{0, 1} {
				if got := b.unitCollision.atLayer(20, 32, z); got != want {
					t.Fatalf("frame%d z%d：持續格=%d，應為%d", frame, z, got, want)
				}
			}
		}
		b.Step()
		if !s.Alive || b.Sides[side].Reserve[0] != 2 {
			t.Fatal("清除的下一拍才應補兵")
		}
	}
}
