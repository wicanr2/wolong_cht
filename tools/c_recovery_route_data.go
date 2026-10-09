//go:build matching_route

package main

import (
	"bytes"
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

// Independent spec249 model; IDA9.4 linear base10000, fixed probe
// ff99cfbd72b7211d14ac5d7544ad7ffa5af9fa8ac8c421d2f1d4311a438d1a66.
// KI.EXE fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868.
type routeNode struct {
	offset uint16
	ports  [4]uint16
	owner  byte
}
type routeLink struct {
	offset, a, b uint16
	length       byte
}
type routeVisit struct {
	offset uint16
	value  byte
}
type routeQueue struct{ offset, node, cost, step, link uint16 }
type routeVector struct {
	graphSegment, occupancyBase                                                                                           uint16
	nodes                                                                                                                 []routeNode
	links                                                                                                                 []routeLink
	visited                                                                                                               []routeVisit
	queue                                                                                                                 []routeQueue
	general                                                                                                               int
	armyFlags, armyOwner, armyStage, armyLeader, ruler                                                                    byte
	generalFlags, generalDuty, generalOwner, captor, temperament, score, factionFlags, capital, occupancy, panel, rngSeed byte
	current, target                                                                                                       uint16
}

const routeWorld = 0x70000
const routeGridSize = 384 * 256

func routePrepare(mem []byte, cs uint16, c resumeCase) {
	v, code := c.route, int(cs)*16
	if v.graphSegment != 0x8400 && v.graphSegment != 0x8500 {
		panic("route fixture must use isolated8400/8500 graph bank")
	}
	if v.general < 0 || v.general >= 128 || v.armyOwner >= 22 || v.generalOwner >= 22 {
		panic("route fixture record/owner bounds")
	}
	graph := int(v.graphSegment) * 16
	clear(mem[graph : graph+65536])
	for i := 0x8000; i < 0x8800; i++ {
		mem[graph+i] = 0xA5
	}
	for i := 0x8800; i < 0x8C00; i++ {
		mem[graph+i] = byte(i*37 + i/8)
	}
	for _, n := range v.nodes {
		if n.offset&7 != 0 || n.offset >= 0x800 {
			panic("route fixture node must be aligned below800")
		}
		for port, link := range n.ports {
			strategyPutWord(mem, graph+int(n.offset)+port*2, link)
		}
		if n.offset < 0x600 {
			mem[routeWorld+0x840+int(n.offset)*4+1] = n.owner
		}
	}
	for i, link := range v.links {
		if link.offset < 0x800 || link.offset >= 0x4000 || link.a&7 != 0 || link.b&7 != 0 {
			panic("route fixture reciprocal link bounds")
		}
		a := graph + int(link.offset)
		point := uint16(0x4000 + i*8)
		strategyPutWord(mem, a, point)
		strategyPutWord(mem, a+2, point+uint16(link.length)*4)
		mem[a+4] = link.length
		strategyPutWord(mem, a+6, link.a)
		strategyPutWord(mem, a+8, link.b)
	}
	for _, mark := range v.visited {
		mem[graph+int(mark.offset)] = mark.value
	}
	for _, q := range v.queue {
		for i, value := range []uint16{q.node, q.cost, q.step, q.link} {
			strategyPutWord(mem, graph+int(q.offset)+i*2, value)
		}
	}
	strategyPutWord(mem, code+0x9874, v.graphSegment)
	strategyPutWord(mem, code+0x9872, v.occupancyBase)
	strategyPutWord(mem, code+0x0D52, 0x7000)
	strategyPutWord(mem, code+0x0CFD, uint16(c.player)*64)
	mem[code+0x0CFF] = c.player
	mem[code+0x98A6], mem[code+0x2919] = v.panel, v.captor
	strategyPutWord(mem, code+0x49B8, 0xA55A)
	strategyPutWord(mem, code+0x49BE, 0x5AA5)
	mem[code+0x49D2] = 0xE3
	mem[code+0xECFC], mem[code+0xECFD] = v.rngSeed^0xA5, v.rngSeed
	for i := 0; i < 256; i++ {
		mem[code+0xECFE+i] = byte(i*73 + i/2 + int(v.rngSeed))
	}
	occupancy := int(v.occupancyBase) * 16
	clear(mem[occupancy : occupancy+routeGridSize])
	a := routeWorld + 0x2240 + v.general*64
	mem[a], mem[a+1], mem[a+2], mem[a+0x23] = v.armyFlags, v.armyOwner, v.armyLeader, v.armyStage
	strategyPutWord(mem, a+0x0E, v.current)
	strategyPutWord(mem, a+0x14, v.target)
	strategyPutWord(mem, a+0x10, 100)
	strategyPutWord(mem, a+0x12, 100)
	strategyPutWord(mem, a+0x1A, 100)
	strategyPutWord(mem, a+0x1C, v.occupancyBase+100*24)
	mem[mainPhysical(v.occupancyBase+100*24, 100)] = v.occupancy
	g := routeWorld + 0x4240 + v.general*32
	mem[g], mem[g+0x17], mem[g+0x1C], mem[g+0x1D], mem[g+0x1E], mem[g+0x1F] = v.generalFlags, v.generalDuty, v.generalOwner, 255, v.temperament, v.score
	for _, owner := range []byte{v.armyOwner, v.generalOwner} {
		f := routeWorld + int(owner)*64
		mem[f], mem[f+1], mem[f+3], mem[f+0x14], mem[f+0x18] = v.factionFlags, v.ruler, v.capital, 1, 1
	}
}

type routeResult struct {
	ax, bx, cx uint16
	carry      bool
}
type routeModel struct {
	mem                    []byte
	cs, graph              uint16
	calls                  map[uint16]int
	searches               []registers
	warnings               []registers
	result                 routeResult
	checkResult, clearDF   bool
	edgeBX, edgeCX, edgeBP uint16
	checkEdge              bool
	minimaps               [][2]uint16
}

func (m *routeModel) enter(ip uint16)           { m.calls[ip]++ }
func (m *routeModel) at(offset uint16) int      { return mainPhysical(m.graph, offset) }
func (m *routeModel) word(offset uint16) uint16 { return strategyWord(m.mem, m.at(offset)) }
func (m *routeModel) put(offset, value uint16)  { strategyPutWord(m.mem, m.at(offset), value) }
func routeRing(offset uint16) uint16            { return (offset + 8) & 0xFBFF }

// Port visited offsets wrap as16-bit offsets. No node-keyed visited set or
// heap/Dijkstra replacement is used. Reciprocal search retains original order.
func (m *routeModel) edge(port, encoded, cost, tail uint16, count byte, bp uint16) (uint16, byte, uint16) {
	m.enter(0x4A0F)
	mark := port - 0x8000
	if m.mem[m.at(mark)] != 0 {
		return tail, count, bp
	}
	m.mem[m.at(mark)] = 1
	link, direction := encoded&0x3FFF, encoded&0xC000
	if direction < 0x4000 {
		return tail, count, bp
	}
	other := m.word(link + 6)
	bp = 4
	if direction == 0x4000 {
		other = m.word(link + 8)
		bp = 0xFFFC
	}
	cost += uint16(m.mem[m.at(link+4)])
	portBack := other
	found := false
	for tries := 0; tries < 4; tries++ {
		if m.word(portBack)&0x3FFF == link {
			found = true
			break
		}
		portBack += 2
	}
	if !found {
		panic("route fixture has no reciprocal port within original node")
	}
	if m.mem[m.at(portBack-0x8000)] != 0 {
		return tail, count, bp
	}
	m.mem[m.at(portBack-0x8000)] = 1
	for i, value := range []uint16{portBack & 0xFFF8, cost, bp, m.word(portBack) & 0x3FFF} {
		m.put(tail+uint16(i*2), value)
	}
	return routeRing(tail), count + 1, bp
}

func (m *routeModel) search(start, stopA, stopB uint16, owner byte) routeResult {
	m.enter(0x491B)
	m.searches = append(m.searches, registers{AX: start, BX: stopA, CX: stopB, DX: uint16(owner), DS: m.graph})
	if start == stopA || start == stopB {
		return routeResult{start, stopA, stopB, true}
	}
	m.clearDF = true
	code := int(m.cs) * 16
	strategyPutWord(m.mem, code+0x49B8, stopA)
	strategyPutWord(m.mem, code+0x49BE, stopB)
	m.mem[code+0x49D2] = owner
	for i := 0x8000; i < 0x8800; i++ {
		m.mem[m.at(uint16(i))] = 0
	}
	tail, head := uint16(0x8800), uint16(0x87F8)
	minimum, node, cost := uint16(0), start, uint16(0)
	active, pending := byte(1), byte(0)
	initial := true
	for iterations := 0; iterations < 8192; iterations++ {
		if !initial {
			node, cost = m.word(head), m.word(head+2)
		}
		initial = false
		if cost == minimum {
			if node == stopA || node == stopB {
				return routeResult{m.word(head + 4), m.word(head + 6), cost, false}
			}
			if node < 0x600 {
				if m.mem[routeWorld+0x840+int(node)*4+1] != owner {
					cost += 0xA6
					cost |= 0x8000
				}
				cost += 4
			}
			for port := uint16(0); port < 8; port += 2 {
				tail, pending, _ = m.edge(node+port, m.word(node+port), cost, tail, pending, 0)
			}
		} else {
			for word := uint16(0); word < 8; word += 2 {
				m.put(tail+word, m.word(head+word))
			}
			tail = routeRing(tail)
			pending++
		}
		head = routeRing(head)
		active--
		if active != 0 {
			continue
		}
		if pending == 0 {
			return routeResult{minimum, 0x800, 0, true}
		}
		minimum = 0xFFFF
		read := head
		for i := 0; i < int(pending); i++ {
			if value := m.word(read + 2); value < minimum {
				minimum = value
			}
			read = routeRing(read)
		}
		active, pending = pending, 0
	}
	panic("route model exceeded bounded legal graph work")
}

func (m *routeModel) random() byte {
	m.enter(0xECE0)
	a := int(m.cs)*16 + 0xECFC
	value := m.mem[a+2+int(m.mem[a+1])] + m.mem[a]
	m.mem[a] += 0x89
	m.mem[a+1] = value
	return value
}

func (m *routeModel) miniDirty(army int) {
	m.enter(0x2BA8)
	m.mem[army] &^= 0x10
	if m.mem[int(m.cs)*16+0x98A6]&4 != 0 {
		m.enter(0x9656)
		m.minimaps = append(m.minimaps, [2]uint16{strategyWord(m.mem, army+0x10), strategyWord(m.mem, army+0x12)})
		m.calls[0x96CF] += 4
		m.enter(0x96ED)
		m.clearDF = true
	}
}

func (m *routeModel) warning(ax, cx uint16) {
	m.enter(0x8810)
	m.warnings = append(m.warnings, registers{AX: ax, CX: cx, DS: 0x7000})
}
func (m *routeModel) removeCount(army int) {
	m.enter(0x4689)
	m.mem[routeWorld+int(m.mem[army+1])*64+0x14]--
}
func (m *routeModel) decrementOccupancy(army int) {
	segment, offset := strategyWord(m.mem, army+0x1C), strategyWord(m.mem, army+0x1A)
	m.mem[mainPhysical(segment, offset)]--
}

func (m *routeModel) dissolve(general int) {
	m.enter(0x2977)
	army := routeWorld + 0x2240 + general*64
	m.decrementOccupancy(army)
	m.mem[army], m.mem[army+3] = 8, 48
	m.removeCount(army)
	player, captor := m.mem[int(m.cs)*16+0x0CFF], m.mem[int(m.cs)*16+0x2919]
	if m.mem[army+1] == player {
		m.warning(0x0093, 0x1F)
	} else if captor == player {
		m.warning(0x0093, 0x20)
	}
}

func (m *routeModel) capture(general int) {
	m.enter(0x29C3)
	army, g := routeWorld+0x2240+general*64, routeWorld+0x4240+general*32
	if m.mem[army] >= 0x80 {
		m.decrementOccupancy(army)
		m.removeCount(army)
	}
	old := m.mem[g+0x1C]
	m.enter(0x2AD2)
	if old != 255 {
		m.mem[routeWorld+int(old)*64+0x18]--
	}
	m.mem[army], m.mem[g+0x17] = 0, 4
	captor, player := m.mem[int(m.cs)*16+0x2919], m.mem[int(m.cs)*16+0x0CFF]
	m.mem[g+0x1C], m.mem[g+0x1D] = captor, old
	if m.mem[routeWorld+int(old)*64] < 0x80 && m.mem[g]&0x10 != 0 {
		m.mem[g], m.mem[g+0x1C], m.mem[g+0x1D] = 0, 255, 255
		if captor == player {
			m.warning(0x0093, 0x43)
		}
		return
	}
	if m.mem[g]&0x40 != 0 {
		m.mem[g] &^= 0x40
		m.mem[g+0x1E] += 3
	}
	if old == player {
		m.warning(0x0093, 0x21)
	} else if captor == player {
		m.warning(0x0093, 0x22)
		m.warning(uint16(m.mem[g+0x1E])<<8|uint16(m.mem[g+1]), 0x19A)
	}
}

func (m *routeModel) collapse(general int, captor byte) {
	m.enter(0x291A)
	a := routeWorld + 0x2240 + general*64
	if m.mem[a] < 0x80 {
		return
	}
	m.mem[int(m.cs)*16+0x2919] = captor
	m.miniDirty(a)
	f := routeWorld + int(m.mem[a+1])*64
	retain := false
	if m.mem[f+3] != 255 {
		if m.mem[f+1] == m.mem[a+2] || captor == m.mem[a+1] || captor == 24 {
			retain = true
		} else {
			r := m.random() & 127
			threshold := (m.mem[routeWorld+0x4240+general*32+0x1F] >> 1) + 40
			retain = r <= threshold
		}
	}
	if retain {
		m.dissolve(general)
	} else {
		m.capture(general)
	}
}

func (m *routeModel) replan(general int) bool {
	m.enter(0x47BB)
	a := routeWorld + 0x2240 + general*64
	target, current := strategyWord(m.mem, a+0x14), strategyWord(m.mem, a+0x0E)
	owner := m.mem[a+1]
	var step byte
	if current >= 0x800 {
		endA, endB := m.word(current+6), m.word(current+8)
		if target == endB {
			step = 4
			m.mem[a] |= 1
		} else if target == endA {
			step = 0xFC
			m.mem[a] |= 1
		} else {
			r := m.search(target, endB, endA, owner)
			if r.cx >= 0x8000 && m.mem[a+0x23] >= 10 {
				m.collapse(general, owner)
				return true
			}
			destination := m.word(r.bx + 8)
			if byte(r.ax) == 4 {
				destination = m.word(r.bx + 6)
			}
			step = 4
			m.mem[a] |= 1
			if destination != endB {
				step = 0xFC
			}
		}
	} else {
		if target == current {
			return false
		}
		r := m.search(target, current, current, owner)
		if r.cx >= 0x8000 && m.mem[a+0x23] >= 10 {
			m.collapse(general, owner)
			return true
		}
		point := m.word(r.bx + 2)
		if byte(r.ax) == 4 {
			point = m.word(r.bx)
		}
		strategyPutWord(m.mem, a+0x0C, point)
		strategyPutWord(m.mem, a+0x0E, r.bx)
		m.mem[a] &^= 1
		step = byte(r.ax)
	}
	m.mem[a+0x0A] = step
	return false
}

func (m *routeModel) retreat(general int) routeResult {
	m.enter(0x487B)
	a := routeWorld + 0x2240 + general*64
	owner := m.mem[a+1]
	capital := m.mem[routeWorld+int(owner)*64+3]
	if capital == 255 {
		return routeResult{carry: true}
	}
	start := uint16(capital) * 8
	current := strategyWord(m.mem, a+0x0E)
	endA, endB := current, current
	onLink := current >= 0x800
	if onLink {
		endA, endB = m.word(current+6), m.word(current+8)
		if m.mem[routeWorld+0x840+int(endA)*4+1] != owner {
			endA = endB
		}
		if m.mem[routeWorld+0x840+int(endB)*4+1] != owner {
			if endA == endB {
				return routeResult{carry: true}
			}
			endB = endA
		}
	}
	r := m.search(start, endB, endA, owner)
	if r.carry {
		return routeResult{ax: r.ax << 2, bx: r.ax << 2}
	}
	link := r.bx
	if onLink && r.ax == 0xFFFC || !onLink && r.ax == 4 {
		link += 2
	}
	node := m.word(link + 6)
	if m.mem[routeWorld+0x840+int(node)*4+1] != owner {
		return routeResult{carry: true}
	}
	return routeResult{ax: r.ax, bx: node << 2}
}

func routeExpected(before []byte, cs uint16, c resumeCase) *routeModel {
	m := &routeModel{mem: append([]byte(nil), before...), cs: cs, graph: c.route.graphSegment, calls: map[uint16]int{}}
	v := c.route
	switch c.target {
	case 0x1491B:
		m.result = m.search(c.ax, c.bx, c.cx, byte(c.dx))
		m.checkResult = true
	case 0x14A0F:
		var ch byte
		m.edgeBX, ch, m.edgeBP = m.edge(c.si, c.ax, c.dx, c.bx, byte(c.cx>>8), c.bp)
		m.edgeCX = uint16(ch)<<8 | uint16(byte(c.cx))
		m.checkEdge = true
	case 0x147BB:
		m.result.carry = m.replan(v.general)
	case 0x1487B:
		m.result = m.retreat(v.general)
	case 0x1291A:
		m.collapse(v.general, byte(c.ax))
	case 0x12977:
		m.dissolve(v.general)
	case 0x129C3:
		m.capture(v.general)
	case 0x12BA8:
		m.miniDirty(routeWorld + 0x2240 + v.general*64)
	case 0x19656:
		m.enter(0x9656)
		m.minimaps = append(m.minimaps, [2]uint16{c.dx, c.bx})
		m.calls[0x96CF] = 4
		m.clearDF = true
	case 0x196CF:
		m.enter(0x96CF)
	default:
		panic("unknown independent route target")
	}
	return m
}

// Before planes are a separate fixture input, copied before either test body.
// Only the new blitter is modeled.12BA8's existing196ED overlay is compared by
// the full original/C state check, with its exact entry arguments checked here.
func routeMinimapAudit(before []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, s string) {
		checks++
		if !ok {
			panic("independent route minimap: " + s)
		}
	}
	if c.target != 0x19656 && c.target != 0x196CF {
		return 0
	}
	want := routeBeforePlanes
	type copyCall struct {
		si, bx uint16
		mask   byte
		df     bool
	}
	var calls []copyCall
	if c.target == 0x196CF {
		calls = append(calls, copyCall{c.si, c.bx, byte(c.ax), c.df})
	} else {
		x, y := c.dx>>1, c.bx>>1
		if x > 178 {
			x = 178
		}
		x = (x + 438) >> 3
		y += 38
		source := y*24 + x - 55 - 960
		destination := y*80 + x
		for plane := 0; plane < 4; plane++ {
			calls = append(calls, copyCall{source + uint16(plane)*0xC00, destination, byte(1 << plane), false})
		}
	}
	entries := strategyObserved(observed, 0x96CF)
	check(len(entries) == len(calls), "blitter entry count")
	sourceSegment := strategyWord(before, int(cs)*16+0x0D3A)
	if c.target == 0x196CF {
		sourceSegment = 0x6800
	}
	for index, call := range calls {
		r := entries[index]
		check(r.SI == call.si && r.BX == call.bx && byte(r.AX) == call.mask && r.DS == sourceSegment && r.ES == 0xA0C8, "original source/destination/plane arguments")
		si, di := call.si, call.bx
		for row := 0; row < 4; row++ {
			for column := 0; column < 2; column++ {
				value := before[mainPhysical(sourceSegment, si)]
				at := uint16(0x0C80 + di)
				for plane := 0; plane < 4; plane++ {
					if call.mask&(1<<plane) != 0 {
						want[plane][at] = value
					}
				}
				if call.df {
					si--
					di--
				} else {
					si++
					di++
				}
			}
			di += 0x4E
			si += 0x16
		}
	}
	for plane := 0; plane < 4; plane++ {
		check(bytes.Equal(device.VGA.Planes[plane][:], want[plane][:]), fmt.Sprintf("complete plane%d", plane))
	}
	return checks
}

func routeAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, s string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent route group=%s target=%X: %s", c.group, c.target, s))
		}
	}
	m := routeExpected(before, cs, c)
	code := int(cs) * 16
	graph := int(c.route.graphSegment) * 16
	var raw [258]byte
	raw[0], raw[1] = c.route.rngSeed^0xA5, c.route.rngSeed
	for i := 0; i < 256; i++ {
		raw[i+2] = byte(i*73 + i/2 + int(c.route.rngSeed))
	}
	check(bytes.Equal(before[code+0xECFC:code+0xEDFE], raw[:]), "enumerated raw RNG input contract")
	check(bytes.Equal(after[graph:graph+0x8C00], m.mem[graph:graph+0x8C00]), "graph,port-visited and entire ring queue bytes")
	check(bytes.Equal(after[routeWorld:routeWorld+0x5240], m.mem[routeWorld:routeWorld+0x5240]), "complete world transition")
	occupancy := int(c.route.occupancyBase) * 16
	check(bytes.Equal(after[occupancy:occupancy+routeGridSize], m.mem[occupancy:occupancy+routeGridSize]), "complete occupancy including FF/0 wrap")
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "game date changed")
	check(bytes.Equal(after[code+0xECFC:code+0xEDFE], m.mem[code+0xECFC:code+0xEDFE]), "fixed258-byte RNG recurrence")
	for _, offset := range []int{0x49B8, 0x49BE, 0x49D2} {
		size := 2
		if offset == 0x49D2 {
			size = 1
		}
		check(bytes.Equal(after[code+offset:code+offset+size], m.mem[code+offset:code+offset+size]), "live original code patch")
	}
	check(after[code+0x2919] == m.mem[code+0x2919], "collapse captor patch/early return")
	for _, ip := range []uint16{0x491B, 0x4A0F, 0x47BB, 0x487B, 0x291A, 0x2977, 0x29C3, 0x2BA8, 0x9656, 0x96CF, 0x96ED, 0x4689, 0x2AD2, 0xECE0, 0x8810} {
		check(len(strategyObserved(observed, ip)) == m.calls[ip], fmt.Sprintf("entry%04X expected%d", ip, m.calls[ip]))
	}
	for i, r := range strategyObserved(observed, 0x491B) {
		w := m.searches[i]
		check(r.AX == w.AX && r.BX == w.BX && r.CX == w.CX && byte(r.DX) == byte(w.DX) && r.DS == w.DS, "raw search original caller arguments")
	}
	if m.checkResult {
		check(device.CPU.R[cpu.AX] == m.result.ax && device.CPU.R[cpu.BX] == m.result.bx && device.CPU.R[cpu.CX] == m.result.cx, "raw search AX/BX/CX")
		check(device.CPU.R[cpu.DX] == c.dx && device.CPU.R[cpu.SI] == c.si && device.CPU.R[cpu.DI] == c.di && device.CPU.R[cpu.BP] == c.bp, "raw search preserved DX/SI/DI/BP")
	}
	if m.checkEdge {
		check(device.CPU.R[cpu.BX] == m.edgeBX && device.CPU.R[cpu.CX] == m.edgeCX && device.CPU.R[cpu.BP] == m.edgeBP && device.CPU.Flags&cpu.CF == 0, "edge tail/count/direction and CLC")
		check(device.CPU.R[cpu.AX] == c.ax && device.CPU.R[cpu.DX] == c.dx && device.CPU.R[cpu.SI] == c.si && device.CPU.R[cpu.DI] == c.di, "edge preserved AX/DX/SI/DI")
	}
	if m.checkResult || c.target == 0x147BB || c.target == 0x1487B {
		check((device.CPU.Flags&cpu.CF != 0) == m.result.carry, "original root CF")
	}
	if c.target == 0x1487B && !m.result.carry {
		check(device.CPU.R[cpu.BX] == m.result.bx && device.CPU.R[cpu.AX] == m.result.ax, "retreat original AX/city-relative BX")
	}
	wantDS := uint16(0x7000)
	if c.target == 0x1491B || c.target == 0x14A0F {
		wantDS = c.route.graphSegment
	}
	if c.target == 0x196CF {
		wantDS = 0x6800
	}
	check(device.CPU.Seg[cpu.DS] == wantDS, "original graph/world/source DS contract")
	if c.target == 0x1491B || c.target == 0x14A0F || c.target == 0x147BB || c.target == 0x1487B || c.target == 0x19656 || c.target == 0x196CF {
		wantDF := c.df && !m.clearDF
		check((device.CPU.Flags&cpu.DF != 0) == wantDF, "original CLD/early-return boundary")
	}
	for i, r := range strategyObserved(observed, 0x8810) {
		w := m.warnings[i]
		check(r.CX == w.CX && r.DS == w.DS, "player notification index/DS")
		if w.CX == 0x19A {
			check(r.AX == w.AX, "new temperament/name byte notice")
		} else {
			check(byte(r.AX) == 0x93, "original notification portrait")
		}
	}
	for i, r := range strategyObserved(observed, 0x9656) {
		w := m.minimaps[i]
		check(r.DX == w[0] && r.BX == w[1], "real minimap restore world coordinates")
	}
	for _, r := range strategyObserved(observed, 0x96ED) {
		check(r.DX == (strategyWord(before, code+0x988E)>>1)+440 && r.BX == (strategyWord(before, code+0x9890)>>1)+40, "existing minimap overlay camera arguments")
	}
	checks += routeMinimapAudit(before, cs, c, observed, device)
	return checks
}

func routeGraph(offsets []uint16, pairs [][2]int, lengths []byte, owner byte) ([]routeNode, []routeLink) {
	nodes := make([]routeNode, len(offsets))
	for i, offset := range offsets {
		nodes[i] = routeNode{offset: offset, owner: owner}
	}
	var links []routeLink
	for i, pair := range pairs {
		link := uint16(0x800 + i*16)
		a, b := pair[0], pair[1]
		length := byte(1)
		if len(lengths) != 0 {
			length = lengths[i%len(lengths)]
		}
		links = append(links, routeLink{link, offsets[a], offsets[b], length})
		for _, end := range []struct {
			node  int
			value uint16
		}{{a, 0x4000 | link}, {b, 0x8000 | link}} {
			placed := false
			for port := 0; port < 4; port++ {
				if nodes[end.node].ports[port] == 0 {
					nodes[end.node].ports[port] = end.value
					placed = true
					break
				}
			}
			if !placed {
				panic("route synthetic graph exceeds four ports")
			}
		}
	}
	return nodes, links
}

func routeCases(base func(string, uint32, int) resumeCase) []resumeCase {
	var all []resumeCase
	makeCase := func(group string, target uint32, scenario int) resumeCase {
		c := base(group, target, scenario)
		c.df = false
		c.altDS = false
		c.profile = 0
		c.route = routeVector{graphSegment: 0x8400, occupancyBase: 0xB000, general: c.general,
			armyFlags: 0xD5, armyOwner: byte(c.owner), armyLeader: byte(c.general), ruler: byte((c.general + 1) % 128),
			generalFlags: 0xC0, generalDuty: 1, generalOwner: byte(c.owner), captor: c.player,
			temperament: 253, score: 0, factionFlags: 0x80, capital: 2, occupancy: 1, rngSeed: 3, current: 0, target: 16}
		c.route.nodes, c.route.links = routeGraph([]uint16{0, 8, 16}, [][2]int{{0, 1}, {1, 2}}, nil, byte(c.owner))
		c.ax, c.bx, c.cx, c.dx = 0, 8, 16, 0xA500|uint16(c.owner)
		c.events = []listEvent{{100, 120, 1}}
		c.mapEvents = nil
		c.menuChoices = nil
		return c
	}
	add := func(c resumeCase) { all = append(all, c) }
	for scenario := 0; scenario < 4; scenario++ {
		patterns := [][][2]int{{{0, 1}, {1, 2}}, {{0, 1}, {0, 2}}, {{0, 2}, {0, 1}}, {{0, 1}, {0, 2}, {1, 3}, {2, 3}}, {{0, 1}}}
		for pattern, pairs := range patterns {
			for ownerProfile := 0; ownerProfile < 3; ownerProfile++ {
				for ends := 0; ends < 3; ends++ {
					for lengthProfile := 0; lengthProfile < 3; lengthProfile++ {
						for _, df := range []bool{false, true} {
							c := makeCase("search", 0x1491B, scenario)
							c.df = df
							offsets := []uint16{0, 8, 16}
							if pattern == 3 {
								offsets = append(offsets, 24)
							}
							lengths := [][]byte{{1}, {1, 2, 3, 4}, {255, 1, 2, 255}}[lengthProfile]
							c.route.nodes, c.route.links = routeGraph(offsets, pairs, lengths, c.route.armyOwner)
							if ownerProfile == 1 {
								c.route.nodes[0].owner = 24
							} else if ownerProfile == 2 {
								c.route.nodes[1].owner = c.player
							}
							c.ax = 0
							switch ends {
							case 0:
								c.bx, c.cx = 8, 16
							case 1:
								c.bx, c.cx = 16, 8
							case 2:
								c.bx, c.cx = 16, 16
								if pattern == 3 {
									c.bx, c.cx = 24, 24
								}
							}
							add(c)
						}
					}
				}
			}
		}
		for _, count := range []int{2, 3, 128, 129, 130, 192} {
			for _, reverse := range []bool{false, true} {
				for _, df := range []bool{false, true} {
					c := makeCase("search", 0x1491B, scenario)
					c.df = df
					var offsets []uint16
					var pairs [][2]int
					for i := 0; i < count; i++ {
						offsets = append(offsets, uint16(i*8))
						if i != 0 {
							pairs = append(pairs, [2]int{i - 1, i})
						}
					}
					c.route.nodes, c.route.links = routeGraph(offsets, pairs, nil, c.route.armyOwner)
					c.ax, c.bx, c.cx = 0, uint16((count-1)*8), uint16((count-1)*8)
					if reverse {
						c.ax, c.bx, c.cx = c.bx, 0, 0
					}
					if count == 192 && reverse {
						c.route.graphSegment = 0x8500
					}
					add(c)
				}
			}
		}
		for _, df := range []bool{false, true} {
			for _, ends := range [][2]uint16{{0, 8}, {8, 0}} {
				c := makeCase("search", 0x1491B, scenario)
				c.df = df
				c.bx, c.cx = ends[0], ends[1]
				add(c)
			}
			c := makeCase("search", 0x1491B, scenario)
			c.df = df
			c.route.nodes, c.route.links = routeGraph([]uint16{0x600, 0x608, 0x7F8}, [][2]int{{0, 1}, {1, 2}}, []byte{255}, c.route.armyOwner)
			c.ax, c.bx, c.cx = 0x600, 0x7F8, 0x7F8
			add(c)
		}
		for _, direction := range []uint16{0, 0x4000, 0x8000, 0xC000} {
			for _, visited := range []byte{0, 1} {
				for _, backVisited := range []byte{0, 1} {
					for _, tail := range []uint16{0x8800, 0x8BF8} {
						for _, cost := range []uint16{0, 255, 0x7FFF, 0xFFFE} {
							c := makeCase("edge", 0x14A0F, scenario)
							c.ax, c.si, c.bx, c.cx, c.dx = direction|0x800, 0, tail, 0x0055, cost
							other := uint16(8)
							if direction > 0x4000 {
								c.si, other = 8, 0
							}
							c.route.visited = []routeVisit{{c.si - 0x8000, visited}, {other - 0x8000, backVisited}}
							add(c)
						}
					}
				}
			}
		}
		for _, count := range []byte{127, 255} {
			for _, tail := range []uint16{0x8800, 0x8BF8} {
				c := makeCase("edge", 0x14A0F, scenario)
				c.ax, c.si, c.bx, c.cx = 0x4800, 0, tail, uint16(count)<<8|0x55
				c.route.visited = []routeVisit{{0x8000, 0}, {0x8008, 0}}
				add(c)
			}
		}
		for _, visited := range []byte{0, 1} {
			c := makeCase("edge", 0x14A0F, scenario)
			c.si, c.ax, c.bx = 0x9000, 0, 0x8BF8
			c.route.visited = []routeVisit{{0x1000, visited}}
			c.route.graphSegment = 0x8500
			add(c)
		}
		for _, current := range []uint16{0, 0x800} {
			for _, target := range []uint16{0, 8, 16} {
				for _, stage := range []byte{0, 9, 10, 11} {
					for ownerProfile := 0; ownerProfile < 3; ownerProfile++ {
						for _, df := range []bool{false, true} {
							c := makeCase("replan", 0x147BB, scenario)
							c.df = df
							c.route.current, c.route.target, c.route.armyStage = current, target, stage
							if ownerProfile != 0 {
								owner := c.player
								if ownerProfile == 2 {
									owner = 24
								}
								for i := range c.route.nodes {
									c.route.nodes[i].owner = owner
								}
							}
							add(c)
						}
					}
				}
			}
		}
		for _, current := range []uint16{0, 0x800} {
			for _, capital := range []byte{255, 0, 1, 2} {
				for ownership := 0; ownership < 4; ownership++ {
					for _, df := range []bool{false, true} {
						c := makeCase("retreat", 0x1487B, scenario)
						c.df = df
						c.route.current, c.route.capital = current, capital
						if ownership&1 != 0 {
							c.route.nodes[0].owner = c.player
						}
						if ownership&2 != 0 {
							c.route.nodes[1].owner = c.player
						}
						add(c)
					}
				}
			}
		}
		for seed := 0; seed < 256; seed++ {
			c := makeCase("collapse", 0x1291A, scenario)
			c.route.rngSeed = byte(seed)
			c.ax = uint16(c.route.captor)
			add(c)
		}
		for _, score := range []byte{1, 79, 80, 175, 255} {
			c := makeCase("collapse", 0x1291A, scenario)
			c.route.score = score
			c.ax = uint16(c.route.captor)
			add(c)
		}
		for profile := 0; profile < 12; profile++ {
			c := makeCase("collapse", 0x1291A, scenario)
			c.ax = uint16(c.route.captor)
			switch profile % 6 {
			case 0:
				c.route.armyFlags = 0x7F
				c.ax = 24
			case 1:
				c.route.capital = 255
			case 2:
				c.route.ruler = c.route.armyLeader
			case 3:
				c.ax = uint16(c.route.armyOwner)
			case 4:
				c.ax = 24
			case 5:
				c.route.armyOwner, c.route.generalOwner = c.player, c.player
			}
			if profile >= 6 {
				c.route.panel = 4
				c.route.occupancy = 255
			}
			add(c)
		}
		template := makeCase("collapse", 0x129C3, scenario)
		for _, faction := range []byte{0, 0x80} {
			for _, flags := range []byte{0x80, 0x90, 0xC0, 0xD0} {
				for _, armyFlags := range []byte{0, 0x80} {
					for _, oldOwner := range []byte{template.player, byte(template.owner)} {
						for _, captor := range []byte{template.player, byte(template.owner), 24} {
							c := makeCase("collapse", 0x129C3, scenario)
							c.route.factionFlags, c.route.generalFlags, c.route.armyFlags = faction, flags, armyFlags
							c.route.armyOwner, c.route.generalOwner, c.route.captor = oldOwner, oldOwner, captor
							c.route.temperament = []byte{0, 253, 255}[int(flags+captor)%3]
							c.route.occupancy = byte(0 - int(flags&0x10)/16)
							add(c)
						}
					}
				}
			}
		}
		for _, flags := range []byte{0, 0x7F, 0x80, 0xD5} {
			for _, own := range []bool{false, true} {
				for kind := 0; kind < 3; kind++ {
					for _, occupancy := range []byte{0, 255} {
						c := makeCase("collapse", 0x12977, scenario)
						c.route.armyFlags = flags
						c.route.occupancy = occupancy
						if own {
							c.route.armyOwner = c.player
						}
						c.route.captor = []byte{c.player, byte(c.owner), 24}[kind]
						add(c)
					}
				}
			}
		}
		for _, flags := range []byte{0, 0x10, 0x90, 0xFF} {
			for _, panel := range []byte{0, 4} {
				for _, df := range []bool{false, true} {
					c := makeCase("minimap", 0x12BA8, scenario)
					c.df = df
					c.route.armyFlags, c.route.panel = flags, panel
					add(c)
				}
			}
		}
		for _, x := range []uint16{0, 355, 356, 357, 383, 65535} {
			for _, y := range []uint16{0, 1, 100, 255} {
				for _, df := range []bool{false, true} {
					c := makeCase("minimap", 0x19656, scenario)
					c.dx, c.bx, c.df = x, y, df
					add(c)
				}
			}
		}
		for _, mask := range []byte{0, 1, 2, 4, 8, 15} {
			for _, si := range []uint16{0, 0xFFFE} {
				for _, df := range []bool{false, true} {
					c := makeCase("minimap", 0x196CF, scenario)
					c.ax, c.si, c.bx, c.df = uint16(mask), si, 0x1200, df
					add(c)
				}
			}
		}
	}
	return all
}
