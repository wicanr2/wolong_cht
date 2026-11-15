package tactical

import (
	"fmt"
	"testing"
)

// Native 1A564 measurements: columns are fresh counts 0,31,32,33,48;
// rows are global bytes 0,1,2. The query reads the prior movement snapshot.
func TestAdvantageQueryUsesNativeGlobalAndSnapshot(t *testing.T) {
	counts := []int{0, 31, 32, 33, 48}
	want := [3][5]int{{2, 2, 2, 0, 0}, {2, 2, 2, 1, 1}, {2, 2, 2, 2, 2}}
	for player := 0; player < 2; player++ {
		for global := 0; global < 3; global++ {
			for col, fresh := range counts {
				t.Run(fmt.Sprintf("player%d-global%d-fresh%d", player, global, fresh), func(t *testing.T) {
					random := &fixedRand{seq: []int{17}}
					b := NewBattle(flatField(), SyntheticFormations(), random, 0)
					b.PlayerSide = player
					enemy := 1 - player
					b.Advantage[player] = Advantage(global)
					b.Advantage[enemy] = Advantaged - Advantage(global)
					b.freshCounts[enemy] = fresh
					// Live units deliberately disagree with the cached count.
					for k := range b.Sides[enemy].Soldiers {
						b.Sides[enemy].Soldiers[k] = Soldier{Alive: true, Stamina: 128}
					}
					s := NewScript([]byte{opQAdv, 0}, enemy)
					s.Step(b)
					if s.cond != want[global][col] || s.pc != 2 || random.i != 0 {
						t.Fatalf("condition %d pc%d rng%d", s.cond, s.pc, random.i)
					}
				})
			}
		}
	}
}

func freshSnapshotBattle(counts [2]int, player int) *Battle {
	b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{127}}, 0)
	b.SetPlayerSide(player)
	for side := 0; side < 2; side++ {
		b.Sides[side].Morale = 200
		for k := 0; k < 48; k++ {
			if k > 0 && k > counts[side] {
				continue
			}
			stamina := 0
			if k > 0 || counts[side] == 48 {
				stamina = 1
			}
			x := 1
			if side == 1 {
				x = 62
			}
			y := k + 1
			if k == 0 {
				y = 60
			}
			kind := Infantry
			if k == 0 {
				kind = General
			}
			b.Sides[side].Soldiers[k] = Soldier{Alive: true, Kind: kind, HP: 200, Stamina: stamina, Stun: 9,
				X: x, Y: y, StepX: x, StepY: y, GoalX: x, GoalY: y, Cmd: Duel, Next: Duel, Target: -1}
		}
	}
	return b
}

func TestAdvantageCountsMapCanonicalPlayerRoles(t *testing.T) {
	cases := []struct {
		counts [2]int
		want   [2][2]Advantage
	}{
		{[2]int{48, 43}, [2][2]Advantage{{Even, Even}, {Advantaged, Disadvantaged}}},
		{[2]int{43, 48}, [2][2]Advantage{{Disadvantaged, Advantaged}, {Even, Even}}},
		{[2]int{48, 48}, [2][2]Advantage{{Even, Even}, {Even, Even}}},
		{[2]int{0, 0}, [2][2]Advantage{{Even, Even}, {Even, Even}}},
	}
	for _, c := range cases {
		for player := 0; player < 2; player++ {
			b := freshSnapshotBattle(c.counts, player)
			b.Step()
			if b.freshCounts != c.counts || b.Advantage != c.want[player] {
				t.Fatalf("player%d counts%v got%v advantage%v want%v", player, c.counts, b.freshCounts, b.Advantage, c.want[player])
			}
			for side := 0; side < 2; side++ {
				if b.attackerAdvantage(side) != c.want[player][side] {
					t.Fatal("melee and stored advantage disagree")
				}
			}
		}
	}
}

func TestFreshSnapshotExcludesBirthAndKeepsEarlierCount(t *testing.T) {
	b := freshSnapshotBattle([2]int{}, 0)
	b.Sides[0].Reserve[1] = 1
	b.Sides[0].Kinds[1] = Infantry
	b.Sides[0].Soldiers[8] = Soldier{Kind: Infantry, Stamina: 2, Y: 40, Cmd: Duel, Next: Duel, Target: -1}
	b.Step()
	if b.freshCounts[0] != 0 || b.Sides[0].Fresh() != 1 {
		t.Fatal("newborn must wait until its own movement turn to enter the count", b.freshCounts, b.Sides[0].Fresh())
	}
	b.Step()
	if b.freshCounts[0] != 1 {
		t.Fatal("newborn must enter the following movement snapshot")
	}

	b = freshSnapshotBattle([2]int{}, 0)
	victim := &b.Sides[0].Soldiers[1]
	*victim = Soldier{Alive: true, Kind: Infantry, HP: 30, Stamina: 1, Stun: 9, X: 21, Y: 20, StepX: 21, StepY: 20, GoalX: 21, GoalY: 20, Cmd: Duel, Next: Duel, Target: -1}
	b.Sides[1].Soldiers[1] = Soldier{Alive: true, Kind: Infantry, HP: 200, Power: 255, Stamina: 128, X: 20, Y: 20, StepX: 21, StepY: 20, GoalX: 21, GoalY: 20, Cmd: Holding, Next: Holding, Target: -1}
	b.Step()
	if victim.Alive || b.freshCounts[0] != 1 || b.Sides[0].Fresh() != 0 {
		t.Fatal("later collision must not rewrite the earlier counted slot", victim.Alive, b.freshCounts, b.Sides[0].Fresh())
	}
}

func TestFormationPreservesNativeRegisterAndGoalOrder(t *testing.T) {
	cases := []struct {
		name                   string
		k                      int
		cmd, next              Command
		gx, gy, step, wantStep int
		oldAtSelf, newAtSelf   bool
		wantCmd                Command
	}{
		{"leader compares old goal bytes even without new order", 0, Form, Form, 20, 21, 2, 0, true, false, Holding},
		{"leader equal old goal bytes preserve height", 0, Form, Form, 20, 20, 2, 2, true, false, Holding},
		{"leader equal goal bytes preserve common changed height", 0, wallAttack, Form, 20, 20, 2, 4, true, false, Holding},
		{"leader old goal is checked before new goal", 0, wallAttack, Form, 20, 21, 2, 0, true, false, Holding},
		{"leader reaching new goal waits for next command pass", 0, Form, Form, 20, 21, 2, 0, false, true, Form},
		{"member new form clears height", 1, wallAttack, Form, 20, 21, 2, 0, false, false, Form},
		{"member same form preserves height", 1, Form, Form, 20, 21, 2, 2, false, false, Form},
		{"leader new holding clears height", 0, wallAttack, Holding, 20, 20, 2, 0, false, false, Holding},
		{"leader same holding preserves height", 0, Holding, Holding, 20, 21, 2, 2, false, false, Holding},
		{"leader new guard clears height", 0, wallAttack, Guard, 20, 21, 2, 0, false, false, Guard},
		{"member holding uses common changed height", 1, wallAttack, Holding, 20, 21, 2, 4, false, false, Holding},
	}
	for player := 0; player < 2; player++ {
		for side := 0; side < 2; side++ {
			for _, c := range cases {
				t.Run(fmt.Sprintf("%d-%d-%s", player, side, c.name), func(t *testing.T) {
					b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
					b.SetPlayerSide(player)
					x, y := 22, 24
					if c.oldAtSelf {
						x, y = c.gx, c.gy
					}
					if c.newAtSelf {
						x, y = b.formationSpot(side, c.k)
					}
					s := &b.Sides[side].Soldiers[c.k]
					*s = Soldier{Alive: true, Kind: Infantry, HP: 100, Stamina: 1, Stun: 9, X: x, Y: y, Z: 4, StepX: 12, StepY: 13, StepZ: c.step, GoalX: c.gx, GoalY: c.gy, Cmd: c.cmd, Next: c.next, Target: -1}
					b.updateSoldierCommand(side, c.k)
					if s.StepZ != c.wantStep || s.Cmd != c.wantCmd {
						t.Fatalf("StepZ%d Cmd%d; want StepZ%d Cmd%d", s.StepZ, s.Cmd, c.wantStep, c.wantCmd)
					}
				})
			}
		}
	}
}

func TestNewFormationLeaderPropagatesNativeTargetOrder(t *testing.T) {
	for player := 0; player < 2; player++ {
		for side := 0; side < 2; side++ {
			for _, changed := range []bool{false, true} {
				b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{0}}, 0)
				b.SetPlayerSide(player)
				leader := &b.Sides[side].Soldiers[8]
				*leader = Soldier{Alive: true, Kind: Infantry, HP: 100, Cmd: Form, Next: Form, X: 20, Y: 30, GoalX: 22, GoalY: 31}
				if changed {
					leader.Cmd = wallAttack
				}
				for k := 9; k < 16; k++ {
					b.Sides[side].Soldiers[k] = Soldier{Alive: k != 10, Cmd: Guard, Next: Guard}
				}
				b.Sides[side].Soldiers[9].Cmd = Retreat
				b.Sides[side].Soldiers[16].Next = Charge
				b.updateSoldierCommand(side, 8)
				for k := 9; k < 16; k++ {
					s := b.Sides[side].Soldiers[k]
					want := Guard
					if changed {
						want = Form
					}
					if s.Next != want || s.TargetRefresh != changed {
						t.Fatal("new Form must propagate to exactly seven members", changed, k, s.Next, s.TargetRefresh)
					}
				}
				if b.Sides[side].Soldiers[16].Next != Charge || b.Sides[side].Soldiers[9].Cmd != Retreat {
					t.Fatal("propagation must preserve the following squad and active retreat")
				}
			}
		}
	}
}
