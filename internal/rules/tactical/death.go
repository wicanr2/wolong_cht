package tactical

// 倒地動畫。出處 `sub_1B360`（`0001B360`），規格 docs/spec/68。
//
// 血歸零時原版只清掉「在場」那個位元、把同一個計時器設成 4，
// 之後每幀由 `sub_1B360` 換一組圖重畫，歸零才由 `sub_1B4B8(ah=1)` 收掉。
// 四幀不算存活單位，但碰撞格保留到 sub_1B3B2 清除（spec/159）。
// 不能由 bit 7 已清推定碰撞格為空。

const (
	// DeathFrames 是倒地要幾幀（`0001B69D` 的 `mov byte ptr [di+1], 4`）。
	DeathFrames = 4
	// DeathSpriteBase 是倒地圖的圖號起點。原版寫的是 raw `0x168`，
	// 而 raw ＝ 192 ＋ 圖號 × 2，所以圖號是 84。
	DeathSpriteBase = 84
	// DeathSpritePerSide 是換到另一側要加多少（`sub_1B240` 的 `add cx, 5Ah`）。
	DeathSpritePerSide = 90
	// deathLateFrames 是「後兩幀換第二張」的門檻（`cmp [si+1], 2 / ja`）。
	deathLateFrames = 2
)

// Death 是一筆倒地動畫。位置固定，只有計時在走。
type Death struct {
	Slot                         int
	Side                         int
	X, Y, Z                      int
	Kind                         Kind
	FramesLeft                   int
	Pending                      bool
	BeforeFacing                 int
	BeforeStep                   uint8
	BeforeHurt, BeforeHitGeneral bool
}

// Sprite 是這一幀該畫哪一張。
//
//	圖號 ＝ 84 ＋ 側 × 90 ＋ 兵種組 × 2 ＋ (計時 ≤ 2 ? 1 : 0)
//
// 兵種只分三組：大將與騎馬共用第 0 組（`sub_1B360` 的 `cmp [si+4], 24h`
// 只切三段），弓是第 1 組、步是第 2 組。
func (d Death) Sprite() int {
	n := DeathSpriteBase + d.Side*DeathSpritePerSide + deathArmGroup(d.Kind)*2
	if d.FramesLeft <= deathLateFrames {
		n++
	}
	return n
}

// deathArmGroup 重現 `cmp byte ptr [si+4], 24h` 的三分。
// `+0x04` 存的是**已經乘過 18** 的兵種值，所以門檻 36 就是「第 2 種」。
func deathArmGroup(kind Kind) int {
	switch v := int(kind); {
	case v < 36:
		return 0
	case v == 36:
		return 1
	default:
		return 2
	}
}

// Deaths 回傳目前還在播的倒地動畫，供呈現層畫。
func (b *Battle) Deaths() []Death {
	if b == nil || len(b.deaths) == 0 {
		return nil
	}
	return append([]Death(nil), b.deaths...)
}

// addDeath 在兵離場的同一處記一筆。
func (b *Battle) addDeath(s *Soldier, side int, before deathPose) {
	if b == nil || s == nil {
		return
	}
	slot := -1
	for k := range b.Sides[side].Soldiers {
		if s == &b.Sides[side].Soldiers[k] {
			slot = k
			break
		}
	}
	b.deaths = append(b.deaths, Death{Side: side, Slot: slot, X: s.X, Y: s.Y, Z: s.Z,
		Kind: s.Kind, FramesLeft: DeathFrames, Pending: true, BeforeFacing: before.facing, BeforeStep: before.step, BeforeHurt: before.hurt, BeforeHitGeneral: before.hitGeneral})
}

// stepDeath只推進輪到的單兵，已輪過後才被擊倒者等下一拍。
func (b *Battle) stepDeath(side, k int) bool {
	for i, d := range b.deaths {
		if d.Side == side && d.Slot == k {
			b.advanceDeath(i)
			return true
		}
	}
	return false
}
func (b *Battle) advanceDeath(i int) bool {
	d := &b.deaths[i]
	d.FramesLeft--
	d.Pending = false
	if d.Side >= 0 && d.Side < len(b.Sides) && d.Slot >= 0 && d.Slot < SoldiersOnFoot {
		b.Sides[d.Side].Soldiers[d.Slot].Stun = d.FramesLeft
	}
	if d.FramesLeft > 0 {
		return false
	}
	if b.unitCollision != nil && d.Slot >= 0 {
		b.clearUnitCollision(d.Side, d.Slot)
	}
	b.deaths = append(b.deaths[:i], b.deaths[i+1:]...)
	return true
}

// stepDeaths供獨立動畫測試逐拍使用；正式Step沿單兵槽序呼叫stepDeath。
func (b *Battle) stepDeaths() {
	for i := 0; i < len(b.deaths); {
		if !b.advanceDeath(i) {
			i++
		}
	}
}

type deathPose struct {
	facing           int
	step             uint8
	hurt, hitGeneral bool
}

func snapshotDeathPose(s *Soldier) deathPose {
	return deathPose{int(s.Facing), uint8(s.PoseStep), s.Hurt, s.HitGeneral}
}
func (b *Battle) beginDeath(s *Soldier, before deathPose) {
	s.HP, s.Alive = 0, false
	s.TargetRefresh, s.Swapped, s.MoveFlag, s.HighTerrain = false, false, false, false
	s.Stun = DeathFrames
	b.addDeathFor(s, before)
}

// addDeathFor 從指標回查是哪一側的兵。投射物那條路只拿得到指標。
func (b *Battle) addDeathFor(e *Soldier, before deathPose) {
	for side := range b.Sides {
		for k := range b.Sides[side].Soldiers {
			if &b.Sides[side].Soldiers[k] == e {
				b.addDeath(e, side, before)
				return
			}
		}
	}
}

// applyLethalForTest 走與致命傷同一條收尾路徑，給測試用。
func (b *Battle) applyLethalForTest(s *Soldier) {
	b.beginDeath(s, snapshotDeathPose(s))
}
