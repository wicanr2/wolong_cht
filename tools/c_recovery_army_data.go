//go:build matching_army

package main

import (
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

// Army17: IDA9.4 census937f8738, spec252. All patches describe raw input.
// FFFF means CS; FFFD is resolved from that machine's own CS:D2FC.
type armyRecord struct {
	slot                                                                    int
	flags, owner, timer, interval, phase, morale, facing, step, stage, home byte
	men, node, target, path, x, y                                           uint16
}
type armyVector struct {
	active  bool
	steps   int
	cursor  uint16
	hour    byte
	records []armyRecord
	patches []engagementPatch
}

func armyResetStage()                                                  {}
func armyCheckpoint(device *machine.Machine, t uint16, regs registers) {}

func armyPrepare(mem []byte, cs uint16, c resumeCase) {
	code := int(cs) * 16
	graph := c.route.graphSegment
	if graph != 0x8400 && graph != 0x8500 {
		panic("army graph needs isolated8400/8500 bank")
	}
	strategyPutWord(mem, code+0xD18, c.army.cursor)
	mem[code+0xCF3] = c.army.hour
	strategyPutWord(mem, code+0x9874, graph)
	strategyPutWord(mem, code+0xD52, 0x7000)
	// A complete, finite three-point road. routePrepare supplied reciprocal
	// node/link records, but its point slots are not movement coordinates.
	for i, p := range [][4]uint16{{100, 100, 0x44, 0x4000}, {101, 100, 0, 0x4004}, {102, 100, 4, 0x4008}, {103, 100, 0, 0x4104}} {
		_ = i
		a := int(graph)*16 + int(p[3])
		strategyPutWord(mem, a, p[0])
		mem[a+2], mem[a+3] = byte(p[1]), byte(p[2])
	}
	strategyPutWord(mem, int(graph)*16+0x800, 0x4000)
	strategyPutWord(mem, int(graph)*16+0x802, 0x4008)
	for slot := 0; slot < 128; slot++ {
		mem[routeWorld+0x2240+slot*64] = 0
	}
	clear(mem[int(c.route.occupancyBase)*16 : int(c.route.occupancyBase)*16+routeGridSize])
	for _, r := range c.army.records {
		if r.slot < 0 || r.slot >= 128 || r.owner >= 22 || r.x >= 384 || r.y >= 256 {
			panic("army typed record bounds")
		}
		a := routeWorld + 0x2240 + r.slot*64
		mem[a], mem[a+1], mem[a+2], mem[a+3] = r.flags, r.owner, byte(r.slot), r.phase
		strategyPutWord(mem, a+4, r.men)
		mem[a+6], mem[a+8], mem[a+0xA], mem[a+0xB] = r.morale, r.facing, r.step, r.timer
		for _, p := range [][2]uint16{{0xC, r.path}, {0xE, r.node}, {0x10, r.x}, {0x12, r.y}, {0x14, r.target}, {0x1A, r.x}, {0x1C, c.route.occupancyBase + r.y*24}} {
			strategyPutWord(mem, a+int(p[0]), p[1])
		}
		mem[a+0x1E], mem[a+0x20], mem[a+0x21], mem[a+0x23] = r.interval, r.home, 0xA5, r.stage
		for slot, remaining := 0, r.men; slot < 6; slot++ {
			n := remaining
			if n > 100 {
				n = 100
			}
			mem[a+0x29+slot*4], mem[a+0x2A+slot*4] = byte(n), 3
			remaining -= n
		}
		if r.flags >= 0x80 {
			mem[mainPhysical(c.route.occupancyBase+r.y*24, r.x)]++
		}
	}
	// Legal city centers and ownership. Encoded node8 names city1.
	for city := 0; city < 3; city++ {
		a := routeWorld + 0x840 + city*32
		mem[a], mem[a+1], mem[a+0x18] = 0x80, c.player, 1
		strategyPutWord(mem, a+8, uint16(100+city*2))
		strategyPutWord(mem, a+10, 100)
	}
	for owner := 0; owner < 22; owner++ {
		a := routeWorld + owner*64
		mem[a], mem[a+0x1D] = 0x80, 200
		strategyPutWord(mem, a+0x20, 10000)
		mem[a+0x22] = 0
		for other := 0; other < 24; other++ {
			mem[routeWorld+0x600+owner*24+other] = 0
		}
	}
	for _, p := range c.army.patches {
		segment := p.segment
		if segment == 0xFFFF {
			segment = cs
		} else if segment == 0xFFFD {
			segment = strategyWord(mem, code+0xD2FC)
		}
		for i, value := range p.bytes {
			mem[mainPhysical(segment, p.offset+uint16(i))] = value
		}
	}
}

type armyModel struct {
	mem               []byte
	cs, graph         uint16
	touched           map[int]bool
	calls             map[uint16]int
	carry, checkCarry bool
}

func (m *armyModel) byte(a int, v byte)   { m.mem[a] = v; m.touched[a] = true }
func (m *armyModel) word(a int, v uint16) { m.byte(a, byte(v)); m.byte(a+1, byte(v>>8)) }
func (m *armyModel) get(a int) uint16     { return strategyWord(m.mem, a) }
func (m *armyModel) enter(ip uint16)      { m.calls[ip]++ }
func (m *armyModel) road(at uint16) int   { return int(m.graph)*16 + int(at) }
func (m *armyModel) funds(owner byte, cost uint32) {
	a := routeWorld + int(owner)*64 + 0x20
	for i := 0; i < 3; i++ {
		expense := routeWorld + int(owner)*64 + 0x1A + i
		m.byte(expense, m.mem[expense])
	}
	u := uint32(m.mem[a]) | uint32(m.mem[a+1])<<8 | uint32(m.mem[a+2])<<16
	v := int32(u)
	if v&0x800000 != 0 {
		v -= 0x1000000
	}
	v -= int32(cost)
	if v < -655000 {
		v = -655000
	}
	encoded := uint32(v) & 0xFFFFFF
	m.byte(a, byte(encoded))
	m.byte(a+1, byte(encoded>>8))
	m.byte(a+2, byte(encoded>>16))
}
func (m *armyModel) upkeep(a int) {
	m.enter(0x2600)
	m.byte(a+7, m.mem[a+7])
	if m.mem[int(m.cs)*16+0xCF3] != 1 {
		return
	}
	n := m.get(a + 4)
	cost := uint32(n>>5) + 1
	leg := m.get(a+0xE) >= 0x800
	if leg {
		cost = uint32(n>>1) + uint32(n>>2)
	}
	m.enter(0x562B)
	m.funds(m.mem[a+1], cost)
	if !leg {
		v := m.mem[a+6] + 10
		cap := m.mem[routeWorld+int(m.mem[a+1])*64+0x1D]
		if v >= cap {
			v = cap
		}
		m.byte(a+6, v)
	}
}
func (m *armyModel) countdown(a int) {
	m.enter(0x264A)
	if m.mem[a]&0x20 == 0 {
		m.byte(a+3, 0)
		m.byte(a+0x21, 0)
		return
	}
	v := m.mem[a+3] - 1
	if v == 0 {
		v = 1
	}
	m.byte(a+3, v)
}
func (m *armyModel) release(relative uint16) {
	m.enter(0x2A7E)
	a := routeWorld + 0x2240 + int(relative)
	timer := m.mem[a+3] - 1
	m.byte(a+3, timer)
	if timer == 0 {
		m.byte(a, 0)
		general := routeWorld + 0x4240 + int(relative)/2
		m.byte(general+0x17, 0)
		if m.mem[routeWorld+int(m.mem[a+1])*64] < 0x80 {
			m.byte(general+0x1C, 255)
		}
	}
}
func (m *armyModel) routeReplan(a int) bool {
	before := append([]byte(nil), m.mem...)
	r := &routeModel{mem: m.mem, cs: m.cs, graph: m.graph, calls: map[uint16]int{}}
	carry := r.replan((a - routeWorld - 0x2240) / 64)
	for address, value := range m.mem {
		if value != before[address] {
			m.touched[address] = true
		}
	}
	for ip, count := range r.calls {
		m.calls[ip] += count
	}
	return carry
}
func (m *armyModel) direction(a int, step uint16) {
	m.enter(0x2808)
	p := m.get(a+0xC) + step
	dx := m.get(a+0x10) - m.get(m.road(p))
	if dx != 0 {
		m.byte(a+8, byte(dx>>15))
		return
	}
	dy := m.mem[a+0x12] - m.mem[m.road(p)+2]
	if dy != 0 {
		m.byte(a+8, 2+(dy>>7))
	}
}
func (m *armyModel) endpoint(a int) {
	m.enter(0x27A2)
	off := uint16(6)
	if m.mem[a+0xA] == 4 {
		off = 8
	}
	node := m.get(m.road(m.get(a+0xE) + off))
	m.word(a+0xE, node)
	if node < 0x600 {
		city := routeWorld + 0x840 + int(node)*4
		x, y := m.get(city+8), m.get(city+10)
		m.word(a+0x10, x)
		m.word(a+0x12, y)
		m.word(a+0x1A, x)
		m.word(a+0x1C, m.get(int(m.cs)*16+0x9872)+uint16(byte(y))*24)
	}
	if node == m.get(a+0x14) {
		m.byte(a+8, 4)
	}
}
func (m *armyModel) standoff(a int) {
	m.byte(a, m.mem[a]|0x20)
	if m.mem[a+3] == 0 {
		m.byte(a+3, 12)
	}
}
func (m *armyModel) collision(a int, x, y uint16) bool {
	m.enter(0x2831)
	for slot := 0; slot < 127; slot++ {
		other := routeWorld + 0x2240 + slot*64
		if m.mem[other] < 0x80 || m.get(other+0x12) != y || m.get(other+0x10) != x {
			continue
		}
		if m.mem[other+1] == m.mem[a+1] {
			return false
		}
		m.standoff(a)
		return true
	}
	return false
}
func (m *armyModel) cityCollision(a int) bool {
	m.enter(0x2880)
	off := uint16(6)
	if m.mem[a+0xA] == 4 {
		off = 8
	}
	node := m.get(m.road(m.get(a+0xE) + off))
	if m.mem[routeWorld+0x841+int(node)*4] == m.mem[a+1] {
		return false
	}
	m.standoff(a)
	return true
}
func (m *armyModel) walk(a int, p uint16) {
	m.enter(0x2708)
	x := m.get(m.road(p))
	y := m.mem[m.road(p)+2]
	flags := m.mem[m.road(p)+3]
	row := m.get(int(m.cs)*16+0x9872) + uint16(y)*24
	if m.mem[mainPhysical(row, x)] != 0 && m.collision(a, x, uint16(y)) {
		return
	}
	tile := m.mem[mainPhysical(m.get(int(m.cs)*16+0xD44)+uint16(y)*24, x)]
	if m.mem[a]&1 != 0 && tile >= 0xCE && tile <= 0xDD && m.cityCollision(a) {
		return
	}
	m.word(a+0x1C, row)
	m.word(a+0x1A, x)
	m.word(a+0xC, p)
	m.word(a+0x10, x)
	m.byte(a+0x12, y)
	step := int8(m.mem[a+0xA])
	if flags&7 >= 2 && ((flags&0x40 != 0 && step <= 0) || (flags&0x40 == 0 && step >= 0)) {
		m.enter(0x27F6)
		m.direction(a, uint16(int16(-step)))
		m.byte(a+8, m.mem[a+8]^1)
		m.endpoint(a)
	} else {
		m.enter(0x2804)
		m.direction(a, uint16(int16(step)))
	}
}
func (m *armyModel) peace(a int, link uint16) {
	m.enter(0x42AB)
	off := uint16(8)
	if m.mem[a+0xA] == 0xFC {
		off = 6
	}
	node := m.get(m.road(link + off))
	owner := m.mem[routeWorld+0x841+int(node)*4]
	self := m.mem[a+1]
	if owner == self || owner == 24 || m.mem[routeWorld+0x600+int(self)*24+int(owner)] < 0x80 {
		return
	}
	off ^= 2
	node = m.get(m.road(link + off))
	m.word(a+0x14, node)
	m.byte(a+0x20, byte(node>>3))
	m.byte(a, m.mem[a]|2)
}
func (m *armyModel) move(a int) {
	m.enter(0x2662)
	if m.get(a+0xE) == m.get(a+0x14) {
		m.byte(a+8, 4)
		m.enter(0x28F4)
		m.enter(0x4325)
		return
	}
	if m.get(a+0xE) < 0x800 && m.mem[a+1] != m.mem[int(m.cs)*16+0xCFF] {
		m.enter(0x4300)
		city := routeWorld + 0x840 + int(m.get(a+0xE))*4
		if m.mem[city+0x18] <= 1 && m.mem[a+8] != 4 && m.mem[city] >= 0x80 {
			m.byte(a+0x20, byte(m.get(a+0xE)>>3))
			m.byte(a+8, 4)
			m.enter(0x28F4)
			m.enter(0x4325)
			return
		}
	}
	if m.get(a+0xE) >= 0x800 {
		m.peace(a, m.get(a+0xE))
	}
	old := mainPhysical(m.get(a+0x1C), m.get(a+0x1A))
	m.byte(old, m.mem[old]-1)
	if m.mem[a]&2 != 0 {
		m.byte(a, m.mem[a]&^2)
		if m.routeReplan(a) {
			next := mainPhysical(m.get(a+0x1C), m.get(a+0x1A))
			m.byte(next, m.mem[next]+1)
			return
		}
	} else if m.get(a+0xE) < 0x800 {
		if m.routeReplan(a) {
			next := mainPhysical(m.get(a+0x1C), m.get(a+0x1A))
			m.byte(next, m.mem[next]+1)
			return
		}
	}
	if m.mem[a]&1 == 0 {
		m.walk(a, m.get(a+0xC))
		next := mainPhysical(m.get(a+0x1C), m.get(a+0x1A))
		m.byte(next, m.mem[next]+1)
		return
	}
	p := m.get(a + 0xC)
	flags := m.mem[m.road(p)+3]
	step := int8(m.mem[a+0xA])
	if flags&7 >= 2 && ((flags&0x40 != 0 && step <= 0) || (flags&0x40 == 0 && step >= 0)) {
		m.enter(0x27F6)
		m.direction(a, uint16(int16(-step)))
		m.byte(a+8, m.mem[a+8]^1)
		m.endpoint(a)
	} else {
		m.enter(0x26FF)
		m.byte(a, m.mem[a]|1)
		m.walk(a, p+uint16(int16(step)))
	}
	next := mainPhysical(m.get(a+0x1C), m.get(a+0x1A))
	m.byte(next, m.mem[next]+1)
}

func armyAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	m := &armyModel{mem: append([]byte(nil), before...), cs: cs, graph: strategyWord(before, int(cs)*16+0x9874), touched: map[int]bool{}, calls: map[uint16]int{}}
	a := routeWorld + int(c.si)
	switch c.target {
	case 0x125A3:
		cursor := m.get(int(cs)*16 + 0xD18)
		for slot := 0; slot < 16; slot++ {
			a = routeWorld + 0x2240 + int(cursor) + slot*64
			if m.mem[a] >= 0x80 {
				timer := m.mem[a+0xB] - 1
				m.byte(a+0xB, timer)
				if timer == 0 {
					m.byte(a+0xB, m.mem[a+0x1E])
					m.byte(a, m.mem[a]&^0x20)
					m.move(a)
				}
				m.upkeep(a)
				m.countdown(a)
			} else if m.mem[a]&8 != 0 {
				m.release(cursor + uint16(slot)*64)
			}
		}
		cursor += 0x400
		if cursor >= 0x1FC0 {
			cursor = 0
		}
		m.word(int(cs)*16+0xD18, cursor)
	case 0x12600:
		m.upkeep(a)
	case 0x1562B:
		m.funds(byte(c.dx>>8), uint32(c.ax)|uint32(byte(c.dx))<<16)
		if device.CPU.R[cpu.SI] != c.si || device.CPU.Seg[cpu.DS] != 0x7000 {
			panic("army original-first money wrapper SI/DS preservation")
		}
	case 0x1264A:
		m.countdown(a)
	case 0x12662:
		m.move(a)
	case 0x126FF:
		m.enter(0x26FF)
		m.byte(a, m.mem[a]|1)
		m.walk(a, c.bx+uint16(int16(int8(m.mem[a+0xA]))))
	case 0x12708:
		m.walk(a, c.bx)
	case 0x127A2:
		m.endpoint(a)
	case 0x127F6:
		m.direction(a, uint16(int16(-int8(m.mem[a+0xA]))))
		m.byte(a+8, m.mem[a+8]^1)
	case 0x12804:
		m.direction(a, uint16(int16(int8(m.mem[a+0xA]))))
	case 0x12808:
		m.direction(a, c.ax)
	case 0x12831:
		m.checkCarry = true
		m.carry = !m.collision(a, c.dx, c.ax)
	case 0x12880:
		m.checkCarry = true
		m.carry = !m.cityCollision(a)
	case 0x142AB:
		m.peace(a, c.bx)
	case 0x14300:
		city := routeWorld + 0x840 + int(c.bx)*4
		m.checkCarry = true
		m.carry = m.mem[city+0x18] <= 1 && m.mem[a+8] != 4 && m.mem[city] >= 0x80
		if m.carry {
			m.byte(a+0x20, byte(c.bx>>3))
		}
	case 0x128F4:
		city := routeWorld + 0x840 + int(c.bx)*4
		m.checkCarry = true
		m.carry = m.mem[city+1] != m.mem[a+1] && uint16(m.mem[a+0x20])*8 == c.bx
		if m.carry {
			routeInput := c
			routeInput.target = 0x1291A
			routeInput.ax = uint16(m.mem[a+0x20]) * 8
			routeInput.route.general = (int(c.si) - 0x2240) / 64
			extra := routeAudit(before, after, cs, routeInput, observed, device)
			if device.CPU.Flags&cpu.CF == 0 {
				panic("army original-first arrival collapse CF")
			}
			return extra + 1
		}
	case 0x12A7E:
		m.release(c.si)
	default:
		panic("unmodeled Army17 target")
	}
	for address := range m.touched {
		if m.mem[address] != after[address] {
			panic(fmt.Sprintf("army original-first group=%s target=%X offset=%05X want=%02X got=%02X before=%02X", c.group, c.target, address, m.mem[address], after[address], before[address]))
		}
	}
	for target, want := range m.calls {
		actual := len(strategyObserved(observed, target))
		if actual != want {
			panic(fmt.Sprintf("army original-first call%04X group=%s want%d got%d", target, c.group, want, actual))
		}
	}
	if m.checkCarry && (device.CPU.Flags&cpu.CF != 0) != m.carry {
		panic("army original-first CF")
	}
	return len(m.touched) + len(m.calls) + 1
}

func armyCases(base func(string, uint32, int) resumeCase) []resumeCase {
	templates := outcomeCases(base)
	var bases [4]resumeCase
	var found [4]bool
	for _, c := range templates {
		if c.scenario >= 0 && c.scenario < 4 && !found[c.scenario] {
			bases[c.scenario] = c
			found[c.scenario] = true
		}
	}
	var cases []resumeCase
	makeCase := func(scenario int, group string, target uint32) resumeCase {
		if !found[scenario] {
			panic("army lacks original scenario")
		}
		c := engagementClone(bases[scenario])
		c.group, c.target, c.repeat = group, target, 1
		c.engagement = engagementVector{}
		c.army = armyVector{active: true, steps: 1, hour: 0, records: []armyRecord{{slot: 0, flags: 0xC1, owner: c.player, timer: 2, interval: 3, morale: 100, facing: 1, step: 4, men: 1, node: 0x800, target: 8, path: 0x4000, x: 100, y: 100}}}
		c.si, c.di, c.ax, c.bx, c.cx, c.dx = 0x2240, 0x2280, 0, 0x4004, 0, 0
		c.outcome.escape = false
		c.outcome.rngSeed = 3
		c.route.rngSeed = 3
		c.route.graphSegment = 0x8400
		c.route.occupancyBase = 0xB000
		c.route.nodes, c.route.links = routeGraph([]uint16{0, 8, 16}, [][2]int{{0, 1}, {1, 2}}, []byte{2}, c.player)
		c.events = []listEvent{{100, 120, 0}, {100, 120, 1}}
		c.menuChoices = []int{-1}
		return c
	}
	patch := func(c *resumeCase, segment, offset uint16, values ...byte) {
		c.army.patches = append(c.army.patches, engagementPatch{segment: segment, offset: offset, bytes: append([]byte(nil), values...)})
	}
	word := func(c *resumeCase, segment, offset, value uint16) {
		patch(c, segment, offset, byte(value), byte(value>>8))
	}
	for scenario := 0; scenario < 4; scenario++ {
		// Every original batch, including4200 slot127; no movement during the first cycle.
		c := makeCase(scenario, "batch", 0x125A3)
		c.army.steps = 8
		c.army.records = nil
		for slot := 0; slot < 128; slot++ {
			r := armyRecord{slot: slot, flags: 0xC0, owner: c.player, timer: 2, interval: 3, morale: 100, facing: 4, step: 4, men: 1, node: 0, target: 0, path: 0x4000, x: uint16(20 + slot%64), y: uint16(20 + slot/64)}
			c.army.records = append(c.army.records, r)
		}
		cases = append(cases, c)
		c = makeCase(scenario, "batch", 0x125A3)
		inactive := c.army.records[0]
		inactive.slot, inactive.flags, inactive.phase, inactive.owner = 15, 8, 1, byte(c.owner)
		c.army.records = append(c.army.records, inactive)
		patch(&c, 0x7000, 0x4240+15*32+0x17, 7)
		patch(&c, 0x7000, 0x4240+15*32+0x1C, byte(c.owner))
		cases = append(cases, c)
		c = makeCase(scenario, "batch", 0x125A3)
		c.army.hour = 1
		c.army.records[0].men = 3
		station := c.army.records[0]
		station.slot, station.owner, station.node, station.morale, station.men = 1, byte(c.owner), 0, 246, 33
		station.x = 102
		c.army.records = append(c.army.records, station)
		patch(&c, 0x7000, 0x2247, 0x5A)
		patch(&c, 0x7000, 0x2287, 0xA5)
		cases = append(cases, c)
		for _, timer := range []byte{0, 1, 2, 255} {
			for _, interval := range []byte{2, 3} {
				c = makeCase(scenario, "batch", 0x125A3)
				c.army.records[0].timer = timer
				c.army.records[0].interval = interval
				cases = append(cases, c)
			}
		}
		for _, hour := range []byte{1, 0, 2} {
			for _, node := range []uint16{0x5F8, 0x800} {
				for _, n := range []uint16{0, 1, 2, 3, 31, 32, 33, 600} {
					c = makeCase(scenario, "upkeep", 0x12600)
					c.army.hour = hour
					c.army.records[0].node = node
					c.army.records[0].men = n
					cases = append(cases, c)
				}
			}
		}
		for _, pair := range [][2]byte{{0, 200}, {189, 200}, {190, 200}, {199, 200}, {245, 200}, {246, 200}, {255, 255}} {
			c = makeCase(scenario, "upkeep", 0x12600)
			c.army.hour = 1
			c.army.records[0].node = 0
			c.army.records[0].morale = pair[0]
			patch(&c, 0x7000, uint16(c.player)*64+0x1D, pair[1])
			cases = append(cases, c)
		}
		for _, owner := range []byte{0, 7, 21} {
			for _, amount := range []uint16{0, 600} {
				c = makeCase(scenario, "upkeep", 0x1562B)
				c.si = 0x4240
				c.ax, c.dx = amount, uint16(owner)<<8
				cases = append(cases, c)
			}
		}
		for _, profile := range []struct {
			owner, high byte
			amount      uint16
			funds       int32
		}{{7, 0, 1, -655000}, {7, 0, 2, -654999}, {21, 1, 1, 10000}} {
			c = makeCase(scenario, "upkeep", 0x1562B)
			c.si = 0x4240
			c.ax, c.dx = profile.amount, uint16(profile.owner)<<8|uint16(profile.high)
			encoded := uint32(profile.funds) & 0xFFFFFF
			patch(&c, 0x7000, uint16(profile.owner)*64+0x20, byte(encoded), byte(encoded>>8), byte(encoded>>16))
			cases = append(cases, c)
		}
		for _, on := range []byte{0, 0x20} {
			for _, phase := range []byte{0, 1, 2, 12, 255} {
				c = makeCase(scenario, "encounter", 0x1264A)
				c.army.records[0].flags |= on
				c.army.records[0].phase = phase
				cases = append(cases, c)
			}
		}
		for _, target := range []uint32{0x12662, 0x126FF, 0x12708, 0x127A2, 0x127F6, 0x12804} {
			for _, step := range []byte{4, 0xFC} {
				c = makeCase(scenario, "movement", target)
				if target == 0x127F6 || target == 0x12804 {
					c.group = "direction"
				}
				c.army.records[0].step = step
				c.bx = 0x4000
				if step == 0xFC {
					c.army.records[0].path = 0x4008
					c.army.records[0].x = 102
					c.army.records[0].facing = 0
					c.bx = 0x4008
				}
				if target == 0x12708 {
					c.bx = 0x4004
				}
				if target == 0x127A2 {
					c.army.records[0].path = 0x4004
				}
				cases = append(cases, c)
			}
		}
		for _, p := range [][3]uint16{{100, 100, 4}, {102, 100, 4}, {101, 99, 4}, {101, 101, 4}, {101, 100, 4}} {
			c = makeCase(scenario, "direction", 0x12808)
			c.army.records[0].x, c.army.records[0].y = p[0], p[1]
			c.ax = p[2]
			cases = append(cases, c)
		}
		for _, slot := range []int{1, 126, 127} {
			for _, phase := range []byte{0, 2} {
				c = makeCase(scenario, "encounter", 0x12831)
				c.ax, c.dx = 100, 101
				r := c.army.records[0]
				r.slot, r.owner, r.x, r.phase = slot, byte(c.owner), 101, 0
				c.army.records = append(c.army.records, r)
				c.army.records[0].phase = phase
				cases = append(cases, c)
			}
		}
		c = makeCase(scenario, "encounter", 0x12831)
		c.ax, c.dx = 100, 101
		r := c.army.records[0]
		r.slot, r.x = 1, 101
		c.army.records = append(c.army.records, r)
		r.slot, r.owner = 2, byte(c.owner)
		c.army.records = append(c.army.records, r)
		cases = append(cases, c)
		for _, same := range []bool{true, false} {
			for _, phase := range []byte{0, 2} {
				c = makeCase(scenario, "encounter", 0x12880)
				c.army.records[0].phase = phase
				if !same {
					patch(&c, 0x7000, 0x861, byte(c.owner))
				}
				cases = append(cases, c)
			}
		}
		for _, step := range []byte{4, 0xFC} {
			for _, kind := range []int{0, 1, 2, 3} {
				c = makeCase(scenario, "peace", 0x142AB)
				c.bx = 0x800
				c.army.records[0].step = step
				owner := c.player
				relation := byte(0)
				if kind == 1 {
					owner = 24
				}
				if kind >= 2 {
					owner = byte(c.owner)
					relation = 0x7F
					if kind == 3 {
						relation = 0x80
					}
				}
				node := uint16(8)
				if step == 0xFC {
					node = 0
				}
				patch(&c, 0x7000, 0x841+node*4, owner)
				if owner < 24 {
					patch(&c, 0x7000, 0x600+uint16(c.player)*24+uint16(owner), relation)
				}
				cases = append(cases, c)
			}
		}
		for _, profile := range [][3]byte{{0, 0x80, 1}, {1, 0x80, 1}, {2, 0x80, 1}, {1, 0x7F, 1}, {1, 0x80, 4}} {
			c = makeCase(scenario, "arrival", 0x14300)
			c.bx = 8
			c.army.records[0].facing = profile[2]
			patch(&c, 0x7000, 0x860, profile[1])
			patch(&c, 0x7000, 0x878, profile[0])
			cases = append(cases, c)
		}
		for _, same := range []bool{true, false} {
			c = makeCase(scenario, "arrival", 0x128F4)
			c.bx = 8
			c.army.records[0].home = 0
			if !same {
				patch(&c, 0x7000, 0x861, byte(c.owner))
			}
			cases = append(cases, c)
		}
		c = makeCase(scenario, "arrival", 0x128F4)
		c.bx = 8
		c.army.records[0].home = 1
		c.army.records[0].node, c.army.records[0].target, c.army.records[0].x = 8, 8, 102
		patch(&c, 0x7000, 0x861, byte(c.owner))
		cases = append(cases, c)
		c = makeCase(scenario, "arrival", 0x12662)
		c.army.records[0].node, c.army.records[0].target, c.army.records[0].home = 0, 0, 2
		patch(&c, 0x7000, uint16(c.player)*64+3, 0)
		cases = append(cases, c)
		c = makeCase(scenario, "arrival", 0x12662)
		c.army.records[0].owner = byte(c.owner)
		c.army.records[0].node, c.army.records[0].target = 8, 16
		patch(&c, 0x7000, 0x861, byte(c.owner))
		cases = append(cases, c)
		for _, friendly := range []bool{false, true} {
			c = makeCase(scenario, "movement", 0x12662)
			other := c.army.records[0]
			other.slot, other.x, other.owner = 1, 101, byte(c.owner)
			if friendly {
				other.owner = c.player
			}
			c.army.records = append(c.army.records, other)
			patch(&c, 0x7000, 0x861, byte(c.owner))
			patch(&c, 0x5000+100*24, 101, 0xCE)
			cases = append(cases, c)
		}
		c = makeCase(scenario, "movement", 0x12662)
		c.army.records[0].flags, c.army.records[0].node = 0xC0, 0
		patch(&c, 0x7000, 0x861, byte(c.owner))
		patch(&c, 0x5000+100*24, 100, 0xCE)
		cases = append(cases, c)
		for _, phase := range []byte{0, 1, 2} {
			for _, alive := range []byte{0, 0x80} {
				c = makeCase(scenario, "cleanup", 0x12A7E)
				c.si = 0
				c.army.records[0].flags, c.army.records[0].phase = 8, phase
				patch(&c, 0x7000, uint16(c.player)*64, alive)
				patch(&c, 0x7000, 0x4257, 7)
				cases = append(cases, c)
			}
		}
		// Poison map points explicitly; random original tiles must not select a city path.
		for i := range cases {
			if cases[i].scenario == scenario {
				word(&cases[i], 0xFFFF, 0xD44, 0x5000)
				// Preserve explicit gate cases; otherwise supply an ordinary road.
				for _, x := range []uint16{100, 101, 102, 103} {
					found := false
					for _, p := range cases[i].army.patches {
						if p.segment == 0x5000+100*24 && p.offset == x {
							found = true
						}
					}
					if !found {
						patch(&cases[i], 0x5000+100*24, x, 0xBA)
					}
				}
			}
		}
	}
	return cases
}
