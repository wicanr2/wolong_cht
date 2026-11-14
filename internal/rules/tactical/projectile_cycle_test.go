package tactical

import (
	"fmt"
	"reflect"
	"testing"
)

// These literal trajectories were measured from DOS/V KI.EXE with dosgolem,
// calling 1AD2D/1AD7F then 1B941 until removal. See spec159's complete-cycle
// receipt for binary identity, seed, input matrix and the independent comparison.
// Only launch and flight are advanced here; soldier cooldown consumers are not.
func TestProjectileCompleteNativeTrajectories(t *testing.T) {
	type sample struct{ height, velocity, power int }
	normal := []sample{
		{1280, 100, 28}, {1380, 80, 28}, {1460, 60, 28}, {1520, 40, 28},
		{1304, 20, 28}, {1324, 0, 28}, {1324, -20, 28}, {1304, -40, 28},
		{1264, -60, 36}, {1204, -80, 36}, {1124, -100, 36}, {1024, -120, 36},
		{904, -140, 46}, {764, -160, 58}, {604, -180, 58}, {424, -200, 73},
		{224, -220, 92}, {4, -240, 92},
	}
	special := []sample{
		{1280, -256, 32}, {1024, -276, 41}, {768, -296, 52},
		{512, -316, 66}, {256, -336, 83}, {0, -356, 104},
	}
	for player := 0; player < 2; player++ {
		for role := 0; role < 2; role++ {
			for direction := West; direction <= South; direction++ {
				for _, vertical := range []bool{false, true} {
					phases := 1
					if vertical {
						phases = 2
					}
					for phase := 0; phase < phases; phase++ {
						for _, finish := range []string{"empty", "body", "terrain"} {
							t.Run(fmt.Sprintf("player%d-role%d-dir%d-special%t-phase%d-%s", player, role, direction, vertical, phase, finish), func(t *testing.T) {
								random := &fixedRand{seq: []int{1}}
								b := NewBattle(flatField(), SyntheticFormations(), random, 0)
								b.PlayerSide = player
								b.unitCollision = &collisionGrid{}
								side := player
								if role == 1 {
									side = 1 - player
								}
								dx, dy := [4]int{-1, 0, 1, 0}[direction], [4]int{0, -1, 0, 1}[direction]
								trajectory, distance, targetZ, end, launch, cooldown, hp := normal, 8, 4, 18, byte(12), uint8(8), 64
								if vertical {
									trajectory, distance, targetZ, end, launch, cooldown, hp = special, 1, 0, 6, 10, 6, 17
								}
								s := &b.Sides[side].Soldiers[6]
								*s = Soldier{Alive: true, Kind: Archer, HP: 100, X: 30, Y: 30, Z: 4, PoseStep: uint8(phase)}
								target := Soldier{Alive: true, Kind: Cavalry, HP: 100, X: 30 + dx*distance, Y: 30 + dy*distance, Z: targetZ, Cmd: Guard, Next: Attack}
								if finish == "body" {
									b.Sides[1-side].Soldiers[1] = target
									b.unitCollision.drawSlot((1-role)*48+1, target.X, target.Y, target.Z, 0)
									end = 9
									if vertical {
										end = 5
									}
								}
								if finish == "terrain" {
									x, y := 30+dx*2, 30+dy*2
									end = 2
									if vertical {
										x, y, end = target.X, target.Y, 1
									}
									for z := 0; z < 7; z++ {
										b.Field.solid[z][y][x] = true
									}
								}
								if vertical {
									b.launchSpecialProjectile(side, s, &target)
								} else {
									b.launchNormalProjectile(side, s, &target)
								}
								for step := 0; step <= end; step++ {
									if step > 0 {
										b.stepProjectiles()
									}
									var sounds []byte
									if step == 0 {
										sounds = []byte{launch}
									} else if step == end && finish == "body" {
										sounds = []byte{11}
									}
									if got := b.TakeSoundEffects(); !reflect.DeepEqual(got, sounds) {
										t.Fatalf("step %d sounds %v, want %v", step, got, sounds)
									}
									if s.ProjectileCooldown != cooldown || s.PoseStep != uint8(phase) || s.Stun != 2 {
										t.Fatal("flight batch changed shooter timing or phase")
									}
									if step == end {
										if len(b.projectiles) != 0 || len(b.Projectiles()) != 0 {
											t.Fatal("projectile survived native removal boundary")
										}
										break
									}
									if len(b.projectiles) != 1 {
										t.Fatalf("step %d projectile removed too early", step)
									}
									q, want := b.projectiles[0], trajectory[step]
									offset, dir := step, direction
									if vertical {
										offset, dir = 1, direction|0x80
									}
									if q.x != 30+dx*offset || q.y != 30+dy*offset || q.z != want.height>>8 || q.heightFP != want.height || q.velocityFP != want.velocity || q.power != want.power || q.direction != dir || q.slot != role*16+6 || !q.slotOwned {
										t.Fatalf("step %d: %+v, want %+v", step, q, want)
									}
									view := b.Projectiles()[0]
									if view.Special != vertical || view.SpecialFrame != phase || view.X != q.x || view.Y != q.y || view.Z != q.z {
										t.Fatal("presentation lost projectile type, phase or position")
									}
								}
								if finish == "body" && b.Sides[1-side].Soldiers[1].HP != hp {
									t.Fatal("native hit damage differs")
								}
								if want := 1; vertical && random.i != 0 || !vertical && random.i != want {
									t.Fatal("launch consumed a different RNG count")
								}
							})
						}
					}
				}
			}
		}
	}
}
