//go:build matching_engagement

package main

import (
	"bytes"
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

// These inputs are applied after each side has independently completed the
// real 11B5A world restore. They never import the original execution's output.
func engagementFrontCases(bases []resumeCase) []resumeCase {
	if len(bases) != 4 {
		panic("engagement front requires four scenario bases")
	}
	const code, world = uint16(0xFFFF), uint16(0x7000)
	var cases []resumeCase
	patch := func(c *resumeCase, segment, offset uint16, value []byte) {
		c.engagement.postPatches = append(c.engagement.postPatches, engagementPatch{segment, offset, append([]byte(nil), value...)})
	}
	bytePatch := func(c *resumeCase, segment, offset uint16, value byte) {
		patch(c, segment, offset, []byte{value})
	}
	wordPatch := func(c *resumeCase, segment, offset, value uint16) {
		patch(c, segment, offset, []byte{byte(value), byte(value >> 8)})
	}
	makeCase := func(base resumeCase, group string, target uint32, profile int) resumeCase {
		c := engagementClone(base)
		c.group, c.target, c.profile, c.repeat = group, target, profile, 1
		c.engagement.warmup, c.engagement.worldWarmup = true, true
		c.engagement.fullEntrySP = c.entrySP
		c.engagement.phase = "front"
		c.engagement.ds, c.engagement.es = world, 0x1200
		c.engagement.postPatches = nil
		c.outcome.escape = false
		c.si, c.di, c.bx, c.ax, c.cx, c.dx = 0x2240, 0x2280, 0x2280, 100, 0, 100
		ownerA, ownerB := c.player, c.outcome.armies[1].owner
		if profile == 0 {
			ownerA = ownerB
			ownerB = (ownerA + 1) % 22
			if ownerB == c.player {
				ownerB = (ownerB + 1) % 22
			}
		} else if profile == 3 {
			ownerA, ownerB = ownerB, c.player
		}
		flags := byte(0xF0) // Both bit20 and bit10 distinguish the two clear masks.
		if profile <= 1 {
			flags |= 4
		}
		for side, owner := range []byte{ownerA, ownerB} {
			a := uint16(0x2240 + side*64)
			var record [64]byte
			record[0], record[1], record[2], record[3] = flags, owner, byte(side), 0xA5
			record[4], record[6], record[8], record[0x23] = 1, 100, 1, 8
			node := uint16(0)
			capital := byte(0)
			if side == 0 {
				node, capital = 8, 1
			}
			record[0x0E], record[0x14], record[0x20] = byte(node), byte(node), capital
			record[0x10], record[0x12], record[0x1A] = 100, 100, 100
			occupancyRow := uint16(0x1200+0x1800) + 100*24
			record[0x1C], record[0x1D] = byte(occupancyRow), byte(occupancyRow>>8)
			for slot := 0; slot < 6; slot++ {
				record[0x2A+slot*4] = 3
			}
			record[0x29] = 1
			patch(&c, world, a, record[:])
			g := uint16(0x4240 + side*32)
			for _, v := range [][2]byte{{0, 0x80}, {0x11, 8}, {0x12, 8}, {0x17, 1}, {0x1C, owner}, {0x1D, 255}, {0x1E, 3}, {0x1F, 0}} {
				bytePatch(&c, world, g+uint16(v[0]), v[1])
			}
			patch(&c, world, g+0x0E, []byte{0, 0, 0})
			f := uint16(owner) * 64
			for _, v := range [][2]byte{{0, 0x80}, {1, byte(side)}, {3, capital}, {0x14, 1}, {0x18, 2}, {0x23, 2}, {0x2A, 255}} {
				bytePatch(&c, world, f+uint16(v[0]), v[1])
			}
		}
		for index, owner := range []byte{ownerB, ownerA, ownerB} {
			a := uint16(0x840 + index*32)
			for _, v := range [][2]byte{{0, 0x80}, {1, owner}, {0x19, 255}, {0x1A, owner}, {0x1B, 0}} {
				bytePatch(&c, world, a+uint16(v[0]), v[1])
			}
			patch(&c, world, a+0x1C, []byte{255, 255, 255, 255})
		}
		bytePatch(&c, world, 0x840+0x13, 7)
		bytePatch(&c, 0x1200+0x1800+100*24, 100, 2)
		bytePatch(&c, 0x1200+0x1800+100*24, 120, 0)
		wordPatch(&c, code, 0x0D2E, 0x2240)
		wordPatch(&c, code, 0x0D30, 0x2280)
		wordPatch(&c, code, 0x0D32, 0)
		bytePatch(&c, code, 0x0D34, 0xC0)
		bytePatch(&c, code, 0x0D35, 0)
		// 14B63 subtracts one row, then samples left/right/up/down/center.
		mapRow := uint16(0x1200 + 99*24)
		for _, v := range []struct {
			offset uint16
			value  byte
		}{{100 + 0x17F, 6}, {100 + 0x181, 14}, {100, 0xB8}, {100 + 0x300, 0xA8}, {100 + 0x180, 0x14}} {
			bytePatch(&c, mapRow, v.offset, v.value)
		}
		return c
	}
	for _, base := range bases {
		for profile := 0; profile < 4; profile++ {
			cases = append(cases, makeCase(base, "field", 0x14A7B, profile))
			c := makeCase(base, "siege", 0x14ADE, profile)
			c.di = 0x840
			cases = append(cases, c)
		}
		c := makeCase(base, "siege", 0x14ADE, 1)
		c.di = 0x840
		bytePatch(&c, world, 0x2280, 0x7F) // Original selector falls through to true14F8A.
		cases = append(cases, c)
		for _, tile := range []byte{0x14, 0xA8, 0xCA, 0xC0} {
			c = makeCase(base, "terrain", 0x14B63, 1)
			bytePatch(&c, 0x1200+99*24, 100+0x180, tile)
			cases = append(cases, c)
		}
		for direction := uint16(0); direction < 4; direction++ {
			c = makeCase(base, "terrain", 0x14BDD, 1)
			c.ax, c.cx, c.dx = direction, 0x0503, 0x0406
			cases = append(cases, c)
		}
		for _, input := range [][2]byte{{8, 0x14}, {9, 0x14}, {9, 0xCA}} {
			c = makeCase(base, "terrain", 0x14C1A, 1)
			c.bx = uint16(input[0])
			bytePatch(&c, 0x1200+0x1800+100*24, 100, input[1])
			cases = append(cases, c)
		}
		for _, tile := range []byte{0x14, 0xB8, 0xBA, 0x70, 6, 0xB1, 0xA8, 0xCA} {
			c = makeCase(base, "terrain", 0x14C4C, 1)
			c.ax = 0x5500 | uint16(tile)
			cases = append(cases, c)
		}
		for _, match := range []bool{true, false} {
			c = makeCase(base, "front-helper", 0x14C72, 1)
			owner := c.outcome.armies[1].owner
			if !match {
				owner = 24
			}
			c.cx = uint16(owner) << 8
			cases = append(cases, c)
		}
		for _, target := range []uint32{0x14E5C, 0x14ED7, 0x14EB9, 0x14F58, 0x14F71} {
			c = makeCase(base, "front-helper", target, 1)
			c.cx = 0x1D
			if target == 0x14ED7 || target == 0x14F58 || target == 0x14F71 {
				c.di, c.cx = 0x840, 0x1B
				wordPatch(&c, code, 0x0D32, 0x840)
				bytePatch(&c, code, 0x0D34, 0)
			}
			if target == 0x14F71 {
				c.cx = 0x1A
			}
			cases = append(cases, c)
		}
		for _, garrison := range []byte{7, 255} {
			c = makeCase(base, "front-helper", 0x14F8A, 1)
			c.di = 0x840
			bytePatch(&c, world, 0x840+0x13, garrison)
			bytePatch(&c, world, 0x4200, 0x5A)
			cases = append(cases, c)
		}
		c = makeCase(base, "front-helper", 0x14FC8, 1)
		c.bx = 0x4200
		bytePatch(&c, world, 0x4200, 0xF0)
		cases = append(cases, c)
	}
	return cases
}

type engagementFrontModel struct {
	*outcomeModel
	cx, bx                    uint16
	checkCX, checkBX, checkCF bool
	field, flags              byte
	tactical, compareWorld    bool
	list                      []uint16
	frame                     uint16
}

func (m *engagementFrontModel) classify(tile byte) byte {
	base := int(m.cs)*16 + 0x982F
	for row := 0; row < 14; row++ {
		at := base + row*3
		if tile >= m.mem[at+1] && tile <= m.mem[at+2] {
			return m.mem[at]
		}
	}
	return m.mem[base+42]
}

func (m *engagementFrontModel) pair(direction, left, right, up, down byte) uint16 {
	a, b := down, up
	if direction == 1 {
		a, b = up, down
	} else if direction >= 2 {
		a, b = left, right
		if direction == 3 {
			a, b = right, left
		}
	}
	base := int(m.cs)*16 + 0x97F0
	for row := 0; row < 21; row++ {
		at := base + row*3
		if a == m.mem[at+1] && b == m.mem[at+2] {
			return uint16(byte(m.mem[at] + 0xC0))
		}
		if b == m.mem[at+1] && a == m.mem[at+2] {
			return 0x4000 | uint16(byte(m.mem[at]+0xC0))
		}
	}
	return 0xC6
}

func (m *engagementFrontModel) water(kind byte, si uint16) uint16 {
	if kind != 9 {
		return uint16((m.random() & 3) + 0xD1)
	}
	a := outcomeWorld + int(si)
	segment, offset := strategyWord(m.mem, a+0x1C), strategyWord(m.mem, a+0x1A)
	rotation := byte(0)
	if m.mem[mainPhysical(segment, offset)] == 0xCA {
		rotation = 0x40
	}
	// IDA14C33 LDS changes DS.14C41's [SI+1] therefore reads this same
	// occupancy segment; treating it as a world owner would repair the game.
	if m.mem[int(m.cs)*16+0x0CFF] != m.mem[mainPhysical(segment, si+1)] {
		rotation ^= 0x40
	}
	return uint16(rotation)<<8 | 0xD5
}

func (m *engagementFrontModel) terrain(si, di, ax uint16) uint16 {
	code := int(m.cs) * 16
	direction, player := byte(ax), m.mem[code+0x0CFF]
	if m.mem[outcomeWorld+int(si)+1] == player {
		direction = m.mem[outcomeWorld+int(si)+8]
	}
	if m.mem[outcomeWorld+int(di)+1] == player {
		direction = m.mem[outcomeWorld+int(di)+8]
	}
	a := outcomeWorld + int(di)
	segment := strategyWord(m.mem, a+0x1C) - strategyWord(m.mem, code+0x9872) + strategyWord(m.mem, code+0x0D44) - 24
	offset := strategyWord(m.mem, a+0x1A)
	read := func(delta uint16) byte { return m.classify(m.mem[mainPhysical(segment, offset+delta)]) }
	left, right, up, down, center := read(0x17F), read(0x181), read(0), read(0x300), read(0x180)
	if center == 0 {
		return m.pair(direction, left, right, up, down)
	}
	if center < 8 {
		return uint16(center + 0xCE)
	}
	return m.water(center, si)
}

func (m *engagementFrontModel) selectArmy(x, y uint16, owner byte, frame uint16) uint16 {
	m.frame = frame
	best, strength := uint16(0), uint16(0)
	for slot := 0; slot < 127; slot++ {
		pointer := uint16(0x2240 + slot*64)
		a := outcomeWorld + int(pointer)
		if m.mem[a] < 0x80 || strategyWord(m.mem, a+0x12) != y || strategyWord(m.mem, a+0x10) != x || m.mem[a+1] != owner {
			continue
		}
		m.list = append(m.list, pointer)
		strategyPutWord(m.mem, mainPhysical(m.c.entrySS, frame+uint16(len(m.list)-1)*2), pointer)
		if len(m.list) == 1 {
			best = pointer
		}
		//14CBE..14CD9 overwrites AH after SHR AX. MUL AH consumes only AL.
		power := uint16(byte(strategyWord(m.mem, a+4)>>4)) * uint16(m.mem[a+6]>>4)
		power *= uint16(m.mem[outcomeWorld+0x4240+slot*32+0x1F]>>4) + 1
		if power > strength {
			best, strength = pointer, power
		}
	}
	strategyPutWord(m.mem, mainPhysical(m.c.entrySS, frame+0xFE), uint16(len(m.list)))
	m.c.bp = frame
	return best
}

func (m *engagementFrontModel) garrison(city uint16) {
	a := outcomeWorld + 0x4200
	value := m.mem[outcomeWorld+int(city)+0x13]
	m.mem[a+1], m.mem[a+2], m.mem[a+6] = m.mem[outcomeWorld+int(city)+1], 127, 255
	strategyPutWord(m.mem, a+4, uint16(value))
	for slot := 0; slot < 6; slot++ {
		men := value / 6
		if byte(slot) < value%6 {
			men++
		}
		m.mem[a+0x29+slot*4], m.mem[a+0x2A+slot*4] = men, 3
	}
}

func (m *engagementFrontModel) fieldBattle(first, second uint16) uint16 {
	player := m.mem[int(m.cs)*16+0x0CFF]
	a, b := outcomeWorld+int(first), outcomeWorld+int(second)
	m.tactical = m.mem[a+1] == player && m.mem[a]&4 == 0 || m.mem[a+1] != player && m.mem[b+1] == player && m.mem[b]&4 == 0
	if m.tactical {
		if m.mem[a+1] != player {
			m.flags |= 0x80
		}
		m.compareWorld = false
		return 0
	}
	return m.automatic(first, second, 1)
}

func (m *engagementFrontModel) siegeBattle(first, second, city uint16) uint16 {
	player := m.mem[int(m.cs)*16+0x0CFF]
	a, b := outcomeWorld+int(first), outcomeWorld+int(second)
	if second != 0x4200 {
		if m.mem[a+1] == player {
			m.tactical = m.mem[a]&4 == 0
		} else if m.mem[outcomeWorld+int(city)+1] == player {
			m.tactical = m.mem[b]&4 == 0
			if m.tactical {
				m.flags |= 0xC0
			}
		}
	}
	if m.tactical {
		m.compareWorld = false
		return 0
	}
	return m.automatic(first, second, 0)
}

func engagementFrontExpected(before []byte, cs uint16, c resumeCase) *engagementFrontModel {
	code := int(cs) * 16
	r := &routeModel{mem: append([]byte(nil), before...), cs: cs, graph: strategyWord(before, code+0x9874), calls: map[uint16]int{}}
	m := &engagementFrontModel{outcomeModel: &outcomeModel{routeModel: r, c: c, unknown: map[int]bool{}}, compareWorld: true, field: before[code+0x0D34], flags: before[code+0x0D35]}
	switch c.target {
	case 0x14C4C:
		kind := m.classify(byte(c.ax))
		m.ax, m.checkAX = uint16(kind)*0x101, true
	case 0x14BDD:
		m.cx, m.checkCX = m.pair(byte(c.ax), byte(c.cx), byte(c.cx>>8), byte(c.dx>>8), byte(c.dx)), true
	case 0x14C1A:
		m.cx, m.checkCX = m.water(byte(c.bx), c.si), true
	case 0x14B63:
		m.cx, m.checkCX = m.terrain(c.si, c.di, c.ax), true
	case 0x14C72:
		m.bx = m.selectArmy(c.dx, c.ax, byte(c.cx>>8), c.bp)
		m.checkBX, m.checkCF, m.carry = len(m.list) != 0, true, len(m.list) == 0
	case 0x14F8A:
		m.garrison(c.di)
		m.bx, m.checkBX = 0x4200, true
	case 0x14FC8:
		m.mem[outcomeWorld+int(c.bx)] = 0
	case 0x14EB9, 0x14F58, 0x14F71:
		// Their live stack arguments and real TALK calls are audited below.
		// Full state comparison owns the old renderer's incidental caches.
		m.compareWorld = false
	case 0x14E5C:
		m.ax = m.fieldBattle(c.si, c.di)
		m.checkAX = !m.tactical
	case 0x14ED7:
		m.ax = m.siegeBattle(c.si, c.bx, c.di)
		m.checkAX = !m.tactical
	case 0x14A7B:
		rotation := m.terrain(c.si, c.di, c.ax)
		m.field, m.flags = byte(rotation), byte(rotation>>8)
		strategyPutWord(m.mem, code+0x0D32, 0)
		m.mem[code+0x0D34], m.mem[code+0x0D35] = m.field, m.flags
		selected := m.selectArmy(c.dx, c.ax, before[outcomeWorld+int(c.di)+1], c.entrySP-12-0x100)
		if selected == 0 {
			break
		}
		for _, army := range []uint16{c.si, selected} {
			a := outcomeWorld + int(army)
			m.mem[a] &^= 0x20
			m.mem[a+3] = 0
		}
		result := m.fieldBattle(c.si, selected)
		if !m.tactical {
			if byte(result>>8) == 2 {
				m.collapse((int(selected)-0x2240)/64, m.mem[outcomeWorld+int(c.si)+1])
			} else if byte(result>>8) != 0 {
				m.collapse((int(c.si)-0x2240)/64, m.mem[outcomeWorld+int(selected)+1])
			}
		}
	case 0x14ADE:
		m.field, m.flags = byte((c.di-0x840)/32), 0
		strategyPutWord(m.mem, code+0x0D32, c.di)
		m.mem[code+0x0D34], m.mem[code+0x0D35] = m.field, m.flags
		a := outcomeWorld + int(c.si)
		m.mem[a] &^= 0x20
		m.mem[a+3] = 0
		city := outcomeWorld + int(c.di)
		selected := m.selectArmy(strategyWord(m.mem, city+8), strategyWord(m.mem, city+10), m.mem[city+1], c.entrySP-14-0x100)
		placeholder := selected == 0
		if placeholder {
			m.garrison(c.di)
			selected = 0x4200
		}
		result := m.siegeBattle(c.si, selected, c.di)
		if !m.tactical {
			if placeholder {
				m.mem[outcomeWorld+0x4200] = 0
			}
			if byte(result) == 0 {
				if !placeholder && byte(result>>8)&2 != 0 {
					m.collapse((int(selected)-0x2240)/64, m.mem[a+1])
				}
				m.captureCity((int(c.di)-0x840)/32, m.mem[a+1])
			} else if byte(result>>8)&1 != 0 {
				m.collapse((int(c.si)-0x2240)/64, m.mem[city+1])
			}
		}
	default:
		panic("unknown engagement front target")
	}
	return m
}

// The old auto-battle/outcome models are reused. The tactical branch checks
// the front-end decision and real inner return frames; it does not independently
// derive tactical combat arithmetic or initialize C from executed original RAM.
func engagementFrontAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent engagement front group=%s target=%X scenario=%d profile=%d: %s", c.group, c.target, c.scenario, c.profile, detail))
		}
	}
	m := engagementFrontExpected(before, cs, c)
	code := int(cs) * 16
	check(c.engagement.worldWarmup && c.engagement.warmup, "own-side true world prefix required")
	check(strategyWord(before, code+0x9872) == strategyWord(before, code+0x0D46)+0x1800 && strategyWord(before, code+0x987C) == strategyWord(before, code+0x0D46)+0x50D6, "restored main banks before front entry")
	entries := strategyObserved(observed, uint16(c.target))
	check(len(entries) >= 1, "actual original front entry")
	entry := entries[0]
	check(entry.SS == c.entrySS && entry.SP == c.entrySP && entry.DS == engagementResolveSegment(before, cs, c.engagement.ds), "front original ABI")
	check(!m.escape && !c.outcome.escape && len(strategyObserved(observed, 0x1CB1)) == 0, "bounded front fixtures return without eliminating player")
	check(device.CPU.Seg[cpu.CS] == cs && device.CPU.Seg[cpu.SS] == c.entrySS && device.CPU.R[cpu.SP] == c.entrySP+2 && device.CPU.IP == 0xF000, "complete front outer return")
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "front strategic date preserved")
	if m.compareWorld {
		for offset := 0; offset < outcomeWorldSize; offset++ {
			if after[outcomeWorld+offset] != m.mem[outcomeWorld+offset] {
				check(false, fmt.Sprintf("world first=%04X before=%02X expected=%02X actual=%02X", offset, before[outcomeWorld+offset], m.mem[outcomeWorld+offset], after[outcomeWorld+offset]))
			}
		}
		check(true, "complete before-derived front world")
	}
	if m.checkAX {
		check(device.CPU.R[cpu.AX] == m.ax, "front AX")
	}
	if m.checkCX {
		check(device.CPU.R[cpu.CX] == m.cx, "original terrain lookup/rotation CX")
	}
	if m.checkBX {
		check(device.CPU.R[cpu.BX] == m.bx, "original chosen or placeholder BX")
	}
	if m.checkCF {
		check((device.CPU.Flags&cpu.CF != 0) == m.carry, "original selector carry")
		frame := mainPhysical(c.entrySS, c.bp)
		check(strategyWord(after, frame+0xFE) == uint16(len(m.list)), "127-slot selector list count")
		for i, pointer := range m.list {
			check(strategyWord(after, frame+i*2) == pointer, "natural selector list order")
		}
	}
	if c.target == 0x14EB9 || c.target == 0x14F58 || c.target == 0x14F71 {
		messages := strategyObserved(observed, 0x8810)
		check(len(messages) == 1 && messages[0].AX == 0xFF93 && messages[0].CX == c.cx && messages[0].DI == c.entrySP-8 && messages[0].SS == c.entrySS, "real TALK caller live argument frame")
	}
	battle := c.target == 0x14A7B || c.target == 0x14ADE || c.target == 0x14E5C || c.target == 0x14ED7
	if battle {
		count := 0
		if m.tactical {
			count = 1
		}
		check(len(strategyObserved(observed, 0x1B5A)) == count, "player delegation chooses real tactical branch")
		auto := 1 - count
		if c.target == 0x14A7B && len(m.list) == 0 {
			auto = 0
		}
		check(len(strategyObserved(observed, 0x5130)) == auto, "AI/delegation chooses old automatic branch")
		check(after[code+0x0D34] == m.field && after[code+0x0D35] == m.flags, "front battlefield and original rotation flags")
		if m.tactical {
			inner := strategyObserved(observed, 0x1B5A)[0]
			saved := inner.SP - 18
			for _, ip := range []uint16{0x9946, 0xCC31, 0x9A33, 0x9FA0, 0x9FDC, 0x1B76, 0x19CA} {
				check(len(strategyObserved(observed, ip)) == 1, fmt.Sprintf("front true tactical closure%04X", ip))
			}
			loop := strategyObserved(observed, 0x9FA0)[0]
			cleanup := strategyObserved(observed, 0x1B76)[0]
			check(loop.SP == saved && loop.SS == c.entrySS && strategyWord(after, code+0xD340) == saved, "inner19FA6 uses real front stack")
			check(cleanup.SP == saved+2 && cleanup.SS == c.entrySS, "inner19FDC resumes actual11B76")
		}
	}
	if m.tactical || !m.compareWorld {
		// As in the common flow oracle, only recurrence is independent here;
		// the tactical/renderer RNG call count comes from observed true entries.
		rng := &routeModel{mem: append([]byte(nil), before...), cs: cs, calls: map[uint16]int{}}
		for range strategyObserved(observed, 0xECE0) {
			rng.random()
		}
		check(bytes.Equal(after[code+0xECFC:code+0xEDFE], rng.mem[code+0xECFC:code+0xEDFE]), "front observed RNG recurrence")
	} else {
		check(bytes.Equal(after[code+0xECFC:code+0xEDFE], m.mem[code+0xECFC:code+0xEDFE]), "independent front RNG recurrence and count")
		check(len(strategyObserved(observed, 0xECE0)) == m.calls[0xECE0], "independent front RNG entry count")
	}
	return checks
}
