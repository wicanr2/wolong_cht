//go:build matching_engagement

package main

// These are raw caller inputs after each machine has executed its own real
// 11B5A initializer. The IDA9.4 source is c-engagement/ida-closed-v5/ida-probe.json
// (32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9).
// FFFF is CS and FFFE is that machine's D30E bank. No executed RAM is imported.
func engagementCombatCases(bases []resumeCase) []resumeCase {
	if len(bases) != 4 {
		panic("engagement combat requires four scenario bases")
	}
	const code, units = uint16(0xFFFF), uint16(0xFFFE)
	const occupancy, terrain = uint16(0x1200), uint16(0x1200 + 0xF80)
	var cases []resumeCase
	makeCase := func(base resumeCase, target uint32, group string) resumeCase {
		c := engagementClone(base)
		c.group, c.target, c.repeat = group, target, 1
		c.engagement.warmup, c.engagement.worldWarmup = true, false
		c.engagement.ds, c.engagement.es = units, occupancy
		c.engagement.postPatches = nil
		c.engagement.phase = "combat-helper"
		c.outcome.escape = false
		c.si, c.di, c.ax, c.bx, c.cx, c.dx = 0, 0x600, 0, 0, 0, 0
		c.events = []listEvent{{100, 120, 0}, {100, 120, 1}}
		return c
	}
	patch := func(c *resumeCase, segment, offset uint16, value []byte) {
		c.engagement.postPatches = append(c.engagement.postPatches, engagementPatch{segment, offset, append([]byte(nil), value...)})
	}
	bytePatch := func(c *resumeCase, segment, offset uint16, value byte) {
		patch(c, segment, offset, []byte{value})
	}
	wordPatch := func(c *resumeCase, segment, offset, value uint16) {
		patch(c, segment, offset, []byte{byte(value), byte(value >> 8)})
	}
	// 19C45/1AFxx maintain current and previous coordinates, cell indices
	// and destination XY together. The caller supplies a legal 64x64 cell.
	position := func(c *resumeCase, record uint16, x, y, layer byte) {
		for _, offset := range []uint16{6, 7} {
			bytePatch(c, units, record+offset, x)
		}
		for _, offset := range []uint16{8, 9} {
			bytePatch(c, units, record+offset, y)
		}
		for _, offset := range []uint16{0xA, 0xB, 0x12} {
			bytePatch(c, units, record+offset, layer)
		}
		cell := uint16(layer)*0x1000 + uint16(y)*64 + uint16(x)
		for _, offset := range []uint16{0xC, 0xE} {
			wordPatch(c, units, record+offset, cell)
		}
		for _, offset := range []uint16{0x10, 0x14} {
			wordPatch(c, units, record+offset, uint16(y)<<8|uint16(x))
		}
	}
	// 1B8AA produces this 20-byte projectile layout. It stores fixed-point
	// XYZ, previous whole coordinates, the occupancy index and graphic ID.
	projectile := func(c *resumeCase, direction byte, velocity uint16) {
		var record [32]byte
		record[0], record[4], record[5] = 0xC0, 0x1C, direction
		record[7], record[9], record[0xB] = 20, 20, 1
		record[0xC], record[0xD], record[0xE] = 20, 20, 1
		cell := uint16(0x1000 + 20*64 + 20)
		for _, offset := range []int{0x10, 0x12} {
			record[offset], record[offset+1] = byte(cell), byte(cell>>8)
		}
		record[0x14], record[0x15] = byte(velocity), byte(velocity>>8)
		record[0x1C], record[0x1D] = 0x10, 2
		patch(c, units, 0x1400, record[:])
	}
	for _, base := range bases {
		// The real duel caller1A1C5 supplies leader0/600 and state8.
		// Keep both live leaders and run the actual animation/update calls.
		for _, target := range []uint32{0x1A298, 0x1A398, 0x1A3C3} {
			c := makeCase(base, target, "duel")
			c.cx = 0x1B7
			for _, record := range []uint16{0, 0x600} {
				position(&c, record, 32, 32, 0)
				bytePatch(&c, units, record+3, 100)
				wordPatch(&c, units, record+0x1A, 0x0808)
			}
			bytePatch(&c, code, 0xD31F, 0)
			cases = append(cases, c)
		}
		for _, target := range []uint32{0x1A6CF, 0x1A6E8} {
			for _, value := range []byte{8, 9} {
				c := makeCase(base, target, "unit-status")
				bytePatch(&c, units, 0x61B, value)
				bytePatch(&c, units, 0x603, value)
				cases = append(cases, c)
			}
		}
		// 1ABFF attacks the target stored at+1C. Both normal ranged and
		// higher-layer special shots use the real 1B8AA allocator.
		for _, elevated := range []bool{false, true} {
			c := makeCase(base, 0x1ABFF, "ranged")
			c.outcome.armies[0].men[0], c.outcome.armies[0].types[0] = 2, 2
			c.si = 0x20
			layer := byte(0)
			if elevated {
				layer = 1
			}
			position(&c, c.si, 20, 20, layer)
			position(&c, 0x600, 21, 20, 0)
			wordPatch(&c, units, c.si+0x1C, 0x600)
			bytePatch(&c, units, c.si+0x13, 0)
			bytePatch(&c, units, 0x1420, 0)
			cases = append(cases, c)
		}
		for _, target := range []uint32{0x1AD2D, 0x1AD7F} {
			for _, cooldown := range []byte{0, 1} {
				c := makeCase(base, target, "ranged")
				c.outcome.armies[0].men[0], c.outcome.armies[0].types[0] = 2, 2
				c.si, c.ax, c.bx, c.dx = 0x20, 1, 1, 2
				position(&c, c.si, 20, 20, 0)
				bytePatch(&c, units, c.si+0x13, cooldown)
				bytePatch(&c, units, 0x1420, 0)
				cases = append(cases, c)
			}
		}
		// B0D3/B116 use a shared ramp tile and the real upper/lower
		// occupancy tests. B186 deliberately reads BX+1000, as original.
		for _, target := range []uint32{0x1B0D3, 0x1B116, 0x1B15D, 0x1B186} {
			for _, blocked := range []bool{false, true} {
				c := makeCase(base, target, "vertical-move")
				layer := byte(0)
				if target == 0x1B116 {
					layer = 1
				}
				position(&c, 0, 10, 10, layer)
				c.bx = 0x28A
				if target == 0x1B186 {
					c.bx += 0x1000
				}
				bytePatch(&c, terrain, 0x28A, 0xF8)
				occupant := byte(0)
				if blocked {
					occupant = 0x31
				}
				bytePatch(&c, occupancy, 0x28A, occupant)
				bytePatch(&c, occupancy, 0x228A, occupant)
				cases = append(cases, c)
			}
		}
		for _, side := range []uint16{0, 0x600} {
			for _, y := range []byte{15, 16, 46, 47} {
				c := makeCase(base, 0x1B4EA, "unit-reposition")
				c.si = side
				position(&c, side, 20, y, 1)
				cases = append(cases, c)
			}
			c := makeCase(base, 0x1B360, "unit-render")
			c.si = side
			position(&c, side, 20, 20, 0)
			bytePatch(&c, units, side+7, 19)
			bytePatch(&c, units, side+1, 2)
			cases = append(cases, c)
		}
		for _, health := range []byte{1, 69, 70, 100} {
			c := makeCase(base, 0x1B618, "melee")
			bytePatch(&c, units, 0x18, 80)
			bytePatch(&c, units, 0x603, health)
			wordPatch(&c, units, 0x61A, 0x0101)
			bytePatch(&c, code, 0xD31E, 1)
			cases = append(cases, c)
		}
		// B732 exchanges two live units and both occupancy layers.
		for _, layer := range []byte{0, 1} {
			c := makeCase(base, 0x1B732, "unit-swap")
			position(&c, 0, 20, 20, layer)
			position(&c, 0x600, 21, 20, layer)
			bytePatch(&c, occupancy, uint16(layer)*0x1000+0x514, 1)
			bytePatch(&c, occupancy, uint16(layer)*0x1000+0x515, 0x31)
			cases = append(cases, c)
		}
		for _, direction := range []byte{0, 1, 2, 3, 0x80, 0x82} {
			c := makeCase(base, 0x1B8AA, "projectile-spawn")
			position(&c, 0, 20, 20, 0)
			c.ax, c.bx, c.cx, c.dx = 20, uint16(direction), 0x1C00, 0x210
			bytePatch(&c, units, 0x1400, 0)
			cases = append(cases, c)
		}
		for _, target := range []uint32{0x1B97E, 0x1BA2E, 0x1BAB7} {
			for _, profile := range []byte{0, 1, 2} {
				c := makeCase(base, target, "projectile")
				c.outcome.armies[1].men[0] = 2
				c.si = 0x1400
				projectile(&c, profile, 20)
				occupant := byte(0x32) // Live620 subordinate, produced by19C45.
				if profile == 1 {
					occupant = 0
				} else if profile == 2 {
					occupant = 0x80
				}
				bytePatch(&c, occupancy, 0x1514, occupant)
				position(&c, 0x620, 20, 20, 1)
				cases = append(cases, c)
			}
		}
		// 1BFF2's first table contains the live cursor body at AL2.
		// 1C01D's separate second table supplies C30D at27/28 and C4A6 at29.
		for _, dispatch := range []struct {
			target uint32
			index  uint16
		}{{0x1BFF2, 2}, {0x1BFF2, 1}, {0x1C01D, 27}, {0x1C01D, 29}} {
			c := makeCase(base, dispatch.target, "combat-dispatch")
			c.engagement.ds, c.engagement.es = code, code
			c.ax = dispatch.index
			wordPatch(&c, code, 0xD330, 0x1F0+40)
			wordPatch(&c, code, 0xD332, 0xCF)
			cases = append(cases, c)
		}
		for _, target := range []uint32{0x1C407, 0x1C4A6, 0x1C4D2} {
			c := makeCase(base, target, "structure-ui")
			c.di, c.bx, c.cx = 0xC00, 16, 0x100
			// 19CE2 creates a live wall record with strength from the city.
			wordPatch(&c, units, 0xC00, 0x0180)
			wordPatch(&c, units, 0xC18, 300)
			wordPatch(&c, code, 0xD326, 0xFFFF)
			cases = append(cases, c)
		}
		// 19CE2 produces a wall record from a contiguous D0..DF map run.
		// Use that complete record shape and two occupied rows, so B824's
		// actual terrain, occupancy, status and minimap writers all run.
		for _, target := range []uint32{0x1B799, 0x1B824} {
			c := makeCase(base, target, "structure-collapse")
			c.di = 0xC00
			var wall [32]byte
			wall[0], wall[1], wall[6], wall[8], wall[0x1A] = 0x80, 1, 20, 20, 2
			wall[0x10], wall[0x11], wall[0x18], wall[0x19] = 0x14, 5, 0x2C, 1
			patch(&c, units, 0xC00, wall[:])
			for _, cell := range []uint16{0x514, 0x554} {
				bytePatch(&c, terrain, cell, 0xD0)
				for layer := uint16(0); layer < 4; layer++ {
					bytePatch(&c, occupancy, cell+layer*0x1000, 0xE1)
				}
				bytePatch(&c, occupancy, cell+0x9000, 100)
			}
			cases = append(cases, c)
		}
		for _, point := range [][2]uint16{{20, 20}, {0, 0}} {
			c := makeCase(base, 0x1DB9B, "projectile-render")
			c.ax, c.dx, c.bx = 1, point[0], point[1]
			cases = append(cases, c)
		}
		for _, tile := range []uint16{0, 15, 31} {
			c := makeCase(base, 0x1DFBB, "tile-render")
			// 19969 passes AX=D302 to1D958, which stores E15A=AX+80.
			// 1DDB4 passes that source bank, a tile0..31 and EGA DI.
			c.engagement.ds, c.engagement.es = 0x1200+0x1180+0x80, 0xA0C8
			c.ax, c.di = tile, 0
			cases = append(cases, c)
		}
		// 1AED2 supplies CL74 for non-infantry, allowing the real BE35
		// vertical edge. The finite adjacent goal contains a live bit8 link.
		c := makeCase(base, 0x1BD46, "vertical-route")
		c.outcome.armies[0].men[0], c.outcome.armies[0].types[0] = 2, 3
		c.si, c.ax, c.dx, c.bx, c.cx, c.bp = 0x20, 0x0A0A, 0x0A0B, 0x1880, 0x0074, 1
		position(&c, c.si, 10, 10, 0)
		wordPatch(&c, units, c.si+0x14, 0x0A0B)
		wordPatch(&c, units, c.si+0x1A, 0x0505)
		// IDA1BD84 loads ES from CS:D2FC; IDA1BE10 reads ES:[BX].
		// FFFD resolves that bank separately from each initialized machine.
		for _, cell := range []uint16{0x28A, 0x28B, 0x128A, 0x128B} {
			bytePatch(&c, 0xFFFD, cell, 0xF8)
		}
		cases = append(cases, c)
	}
	return cases
}
