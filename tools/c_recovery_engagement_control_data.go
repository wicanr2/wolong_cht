//go:build matching_engagement

package main

// Warmup is independently executed on each machine before these controlled
// inputs. FFFF denotes CS and FFFE denotes the real CS:D30E unit bank.
func engagementControlCases(bases []resumeCase) []resumeCase {
	const code, units = uint16(0xFFFF), uint16(0xFFFE)
	var cases []resumeCase
	makeCase := func(base resumeCase, group string, target uint32) resumeCase {
		c := engagementClone(base)
		c.group, c.target, c.repeat = group, target, 1
		c.engagement.warmup = true
		c.engagement.ds, c.engagement.es = code, code
		c.engagement.postPatches = nil
		c.outcome.rngSeed = 3
		return c
	}
	bytePatch := func(c *resumeCase, segment, offset uint16, value byte) {
		c.engagement.postPatches = append(c.engagement.postPatches, engagementPatch{segment: segment, offset: offset, bytes: []byte{value}})
	}
	wordPatch := func(c *resumeCase, segment, offset, value uint16) {
		c.engagement.postPatches = append(c.engagement.postPatches, engagementPatch{segment: segment, offset: offset, bytes: []byte{byte(value), byte(value >> 8)}})
	}
	moveRecord := func(c *resumeCase, record uint16, x, y byte) {
		bytePatch(c, units, record+6, x)
		bytePatch(c, units, record+8, y)
		bytePatch(c, units, record+0x0A, 0)
		wordPatch(c, units, record+0x0C, uint16(y)*64+uint16(x))
		wordPatch(c, units, record+0x10, uint16(y)<<8|uint16(x))
		wordPatch(c, units, record+0x14, uint16(y)<<8|uint16(x))
	}
	for _, base := range bases {
		// 1A065 has DS=CS in its original 19FA0 caller. The live EB producer
		// is 1C23C, and 199F5 supplies the legitimate dirty value one.
		c := makeCase(base, "patch", 0x1A065)
		bytePatch(&c, code, 0xA06A, 0xEB)
		bytePatch(&c, code, 0xD348, 1)
		cases = append(cases, c)

		// 19969..19974 calls 1D958 with BX=D306, so its real 1D971 grid
		// bank is arena1200+3D20. Poison the last word of 3C00 words.
		c = makeCase(base, "grid", 0x1D971)
		wordPatch(&c, 0x1200+0x3D20, 0x77FE, 0x5AA5)
		cases = append(cases, c)

		// 15FAA calls 160CC with DS=CS before drawing all six AL20..25 keys.
		c = makeCase(base, "settings", 0x160CC)
		cases = append(cases, c)

		// 1ABFF obtains DL0..3 and signed AX from 1ACA4, then 1ACD6
		// consumes those values and [SI+1C]. Both live leader records exist.
		for direction := uint16(0); direction < 4; direction++ {
			c = makeCase(base, "direction", 0x1ACD6)
			c.engagement.ds = units
			c.si, c.ax, c.dx = 0, 0, direction
			wordPatch(&c, units, 0x1C, 0x600)
			bytePatch(&c, units, 0x606, 20)
			bytePatch(&c, units, 0x608, 20)
			cases = append(cases, c)
		}
		// Dominant negative X forces the original ACB4 assignment DL=2.
		c = makeCase(base, "direction", 0x1ACA4)
		c.engagement.ds = units
		c.si = 0
		wordPatch(&c, units, 0x1C, 0x600)
		bytePatch(&c, units, 0x06, 10)
		bytePatch(&c, units, 0x08, 10)
		bytePatch(&c, units, 0x0A, 5)
		bytePatch(&c, units, 0x606, 20)
		bytePatch(&c, units, 0x608, 10)
		bytePatch(&c, units, 0x60A, 5)
		cases = append(cases, c)

		// DB34 is the original signed isometric projection entry. The
		// coordinate difference includes both signs without changing a bank.
		for _, point := range [][2]uint16{{10, 8}, {8, 10}} {
			c = makeCase(base, "projection", 0x1DB34)
			c.ax, c.cx, c.dx, c.bx = 0, 1, point[0], point[1]
			cases = append(cases, c)
		}

		// 19A33 really produces C00/600 for flags01 and 600/000 for
		// flags81. 1B1B1 returns every nonzero occupancy code unchanged,
		// including the current object's code; it has no self-identity filter.
		// Preserve the real live leader and test its occupied adjacent cell.
		for _, profile := range []struct {
			flags, occupancy    byte
			actor, patch, bound uint16
		}{{0x01, 0x31, 0x600, 0xB561, 0x0C00}, {0x81, 1, 0, 0xB569, 0}} {
			c = makeCase(base, "patch", 0x1B533)
			c.engagement.flags = profile.flags
			c.engagement.ds, c.engagement.es = units, 0x1200
			c.si, c.ax, c.bx = profile.actor, uint16(profile.occupancy), 0x28B
			moveRecord(&c, profile.actor, 10, 10)
			bytePatch(&c, units, profile.actor+5, 2) // Original 1B069 east step.
			wordPatch(&c, units, profile.actor+0x1A, 0x0101)
			bytePatch(&c, 0x1200, 0x28A, profile.occupancy)
			bytePatch(&c, 0x1200, 0x28B, profile.occupancy)
			wordPatch(&c, code, profile.patch, profile.bound)
			cases = append(cases, c)
		}

		// BATTLE.MAP field0's real header has limit36. The warmup producer
		// CB13 writes that byte into AB4F. Enter AB39's AL!=AH arm with no
		// active high-layer movement, as 1A988 does after a state transition.
		c = makeCase(base, "patch", 0x1AB39)
		c.engagement.ds = units
		c.si, c.ax = 0, 0x0100
		bytePatch(&c, units, 0, 0x80)
		bytePatch(&c, units, 0x1E, 0)
		cases = append(cases, c)

		// C5D7 derives its immediate from live D32E. D32C=5 is a legal
		// cursor column and matches D33C=5, so the second case exercises
		// the original C5F0 producer of EB at C604. Both loops have64 rows.
		for _, sameColumn := range []bool{false, true} {
			c = makeCase(base, "patch", 0x1C5D7)
			wordPatch(&c, code, 0xD32E, 0x21)
			bytePatch(&c, code, 0xD33C, 5)
			column := uint16(0x20)
			if sameColumn {
				column = 5
			}
			wordPatch(&c, code, 0xD32C, column)
			cases = append(cases, c)
		}

		// 1AED2 supplies AX=current XY, DX=goal XY, BX=1800+SI*4,
		// CH=layer and CL=EB for an infantry unit (type18). The original
		// four horizontal links and bit8 vertical link are explicit terrain
		// inputs; the adjacent goal makes the search short and bounded.
		c = makeCase(base, "patch", 0x1BD46)
		c.outcome.armies[0].types[0] = 1
		c.engagement.ds = units
		c.si, c.ax, c.dx, c.bx, c.cx, c.bp = 0, 0x0A0A, 0x0A0B, 0x1800, 0x00EB, 1
		moveRecord(&c, 0, 10, 10)
		wordPatch(&c, units, 0x14, 0x0A0B)
		bytePatch(&c, 0xFFFD, 0x28A, 0xF0)
		bytePatch(&c, 0xFFFD, 0x28B, 0xF8)
		bytePatch(&c, 0xFFFD, 0x128B, 0xF1)
		cases = append(cases, c)

		// A426 masks the command to five bits and its high three bits are
		// the AL argument. These19 commands use argument0 and AH=1; all
		// following words are bounded script data in the real D30C bank.
		const script = uint16(0x1200 + 0x4A2E)
		for opcode := uint16(0); opcode < 19; opcode++ {
			c = makeCase(base, "script", 0x1A426)
			c.events = []listEvent{{100, 120, 0}, {100, 120, 1}}
			wordPatch(&c, code, 0xD311, 0)
			wordPatch(&c, code, 0xD313, 0)
			wordPatch(&c, script, 0, 0x0100|opcode)
			for offset := uint16(2); offset < 32; offset += 2 {
				wordPatch(&c, script, offset, 0x0100)
			}
			cases = append(cases, c)
		}
		// A591's five legal AL cases read one operand at SI. A zero low
		// byte with high byte1 denotes the valid jump to script word1.
		// Comparison values force each conditional case's distinct arm.
		for kind, comparison := range []byte{0, 0, 1, 0, 2} {
			c = makeCase(base, "script", 0x1A591)
			c.engagement.ds = script
			c.si, c.ax = 0x10, 0x0100|uint16(kind)
			bytePatch(&c, code, 0xD315, comparison)
			wordPatch(&c, script, 2, 0x0100)
			wordPatch(&c, script, 0x10, 0x0100)
			cases = append(cases, c)
		}

		// A754/A785 call A7B7 at each100-byte formation leader and A7FD
		// on its seven20-byte subordinates. Two real men guarantee the
		// first subordinate is created alive by19C45, rather than testing a
		// fabricated zero record. 1A85B's surrounding ES is D2F6.
		for _, family := range []struct {
			target        uint32
			member, count uint16
		}{{0x1A7B7, 0, 11}, {0x1A7FD, 0x20, 9}} {
			for state := uint16(0); state < family.count; state++ {
				c = makeCase(base, "unit", family.target)
				if family.member != 0 {
					c.outcome.armies[0].men[0] = 2
				}
				c.events = []listEvent{{100, 120, 0}, {100, 120, 1}}
				c.engagement.ds, c.engagement.es = units, 0x1200+0xF80
				c.si = family.member
				current := uint16(0)
				if state == 0 {
					current = 1
				}
				wordPatch(&c, units, family.member+0x1A, state<<8|current)
				wordPatch(&c, units, family.member+0x1C, 0x600)
				cases = append(cases, c)
			}
		}
	}
	return cases
}
