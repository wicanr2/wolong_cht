package tactical

import "testing"

// DOS/V 1A52E／1A560：每筆命令≥4視為0，AH保留最大值，不檢查Alive。
func TestScriptQueriesMaximumNormalizedLeaderCommand(t *testing.T) {
	cases := []struct {
		commands [6]Command
		want     int
	}{
		{[6]Command{ScaleWal, ScaleWal, ScaleWal, wallAttack, wallAttack, wallAttack}, 3},
		{[6]Command{Form, Attack, Charge, ScaleWal, Guard, Retreat}, 3},
		{[6]Command{wallAttack, Holding, Duel, 9, 254, 255}, 0},
		{[6]Command{Attack, Charge, Attack, Charge, Attack, Charge}, 2},
		{[6]Command{Form, Form, Form, Form, Form, Attack}, 1},
		{[6]Command{Form, Form, Form, Form, Form, Form}, 0},
	}
	for k := 0; k < 6; k++ {
		var commands [6]Command
		commands[k] = ScaleWal
		cases = append(cases, struct {
			commands [6]Command
			want     int
		}{commands, 3})
	}
	for player := 0; player < 2; player++ {
		for role := 0; role < 2; role++ {
			for _, c := range cases {
				b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{127}}, 0)
				b.SetPlayerSide(player)
				side := player
				if role == 1 {
					side = 1 - player
				}
				for i := range b.Sides {
					for k := range b.Sides[i].Soldiers {
						b.Sides[i].Soldiers[k] = Soldier{Alive: true, Cmd: ScaleWal, Next: Duel}
					}
				}
				for i := range b.Sides {
					for k := 0; k < Squads; k++ {
						b.Sides[i].Soldiers[k*PerSquad].Cmd = Attack
					}
				}
				for k, cmd := range c.commands {
					s := &b.Sides[side].Soldiers[k*PerSquad]
					s.Alive = false
					s.Cmd = cmd
					s.Next = Duel
				}
				op := byte(opQTheirs)
				if role == 1 {
					op = opQMine
				}
				s := NewScript([]byte{op, 0}, 1-player)
				s.Step(b)
				if s.cond != c.want || s.pc != 2 {
					t.Fatalf("player%d role%d commands%v condition%d want%d", player, role, c.commands, s.cond, c.want)
				}
			}
		}
	}
}

// 1AF9B–1AFC2：原位取得相鄰敵人Step後直接返回，沒有1AFC3的Facing寫入。
func TestAdjacentRetargetPreservesFacing(t *testing.T) {
	for side := 0; side < 2; side++ {
		for _, kind := range []Kind{Cavalry, Infantry} {
			for _, d := range [][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}} {
				b := NewBattle(flatField(), SyntheticFormations(), &fixedRand{seq: []int{127}}, 0)
				s := &b.Sides[side].Soldiers[9]
				*s = Soldier{Alive: true, Kind: kind, HP: 200, Stamina: 80, X: 20, Y: 20, StepX: 20, StepY: 20, GoalX: 20, GoalY: 20, Cmd: ScaleWal, Next: ScaleWal, Target: 9, Facing: South, HitGeneral: true}
				if kind == Cavalry {
					s.StepZ = 4
				}
				e := &b.Sides[1-side].Soldiers[9]
				*e = Soldier{Alive: true, Kind: Infantry, HP: 200, X: 20 + d[0], Y: 20 + d[1]}
				b.moveToward(side, 9)
				if s.X != 20 || s.Y != 20 || s.Facing != South || s.StepX != e.X || s.StepY != e.Y || s.StepZ != e.Z || s.MoveFlag || s.HitGeneral {
					t.Fatalf("side%d kind%d delta%v: %+v", side, kind, d, s)
				}
			}
		}
	}
}

func TestNonAdjacentRetargetStillUpdatesFacing(t *testing.T) {
	for _, c := range []struct {
		dx, dy, goal int
		want         int
	}{{0, 0, 0, West}, {2, 0, 0, East}, {0, -2, 0, North}, {1, -1, 2, East}} {
		b := NewBattle(flatField(), SyntheticFormations(), nil, 0)
		s := &b.Sides[0].Soldiers[9]
		*s = Soldier{Alive: true, Kind: Cavalry, HP: 200, Stamina: 80, X: 20, Y: 20, StepX: 20, StepY: 20, GoalX: 20 + c.goal, GoalY: 20, Target: 9, Facing: South}
		b.Sides[1].Soldiers[9] = Soldier{Alive: true, X: 20 + c.dx, Y: 20 + c.dy}
		b.moveToward(0, 9)
		if s.Facing != c.want || s.StepX != s.GoalX || s.StepY != s.GoalY {
			t.Fatal(c, s)
		}
	}
}
