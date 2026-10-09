//go:build matching_tick

package main

import (
	"bytes"
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

// Independent spec/248 model. IDA Pro 9.4 linear base10000; fixed probe
// 1922512894237831508aa8d2ba041d8f6b389a9f473775010e3d1fb2d1d542f5.
// KI.EXE fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868.
type tickCity struct {
	index                                                         int
	flags, owner, oldOwner, cooldown, occupancy, threat, governor byte
	marker, growth, prevention, capacity, garrison, neighborCount byte
	x, y, production                                              uint16
	neighbors                                                     [4]byte
}

type tickObject struct {
	slot                 int
	flags, timer, reload byte
	x, y, vx, vy         uint16
}

type tickArmy struct {
	slot                        int
	flags, stage, owner, target byte
	node                        uint16
}

type tickVector struct {
	city                                                                    int
	cursor, occupancyBase                                                   uint16
	rngSeed                                                                 byte
	candidates                                                              [16]byte
	bounds                                                                  [4]int16
	cities                                                                  []tickCity
	objects                                                                 []tickObject
	armies                                                                  []tickArmy
	ally, relation, armyCount, governorBudget, governorWar, governorAbility byte
	capacityWord                                                            uint16
	reserves                                                                [3]uint16
	dispatch                                                                bool
}

const tickWorld = 0x70000
const tickWorldSize = 0x5240
const tickGridSize = 384 * 256

func tickPrepare(mem []byte, cs uint16, c resumeCase) {
	code, v := int(cs)*16, c.tick
	if v.city < 0 || v.city >= 192 || v.occupancyBase == 0 {
		panic("tick fixture requires legal city and explicit occupancy segment")
	}
	strategyPutWord(mem, code+0x0D52, 0x7000)
	strategyPutWord(mem, code+0x0D1E, v.cursor)
	strategyPutWord(mem, code+0x9872, v.occupancyBase)
	strategyPutWord(mem, code+0x0CFD, uint16(c.player)*64)
	mem[code+0x0CFF], mem[code+0x98A6] = c.player, 0
	for i, value := range v.bounds {
		strategyPutWord(mem, code+0x0D22+i*2, uint16(value))
	}
	// The same raw258-byte state is installed independently on both machines.
	// Seeds are enumerated before execution; no result is used to select one.
	mem[code+0xECFC], mem[code+0xECFD] = v.rngSeed^0xA5, v.rngSeed
	for i := 0; i < 256; i++ {
		mem[code+0xECFE+i] = byte(i*73 + i/2 + int(v.rngSeed))
	}
	occupancy := int(v.occupancyBase) * 16
	clear(mem[occupancy : occupancy+tickGridSize])
	for _, city := range v.cities {
		if city.index < 0 || city.index >= 192 || city.x > 383 || city.y > 255 {
			panic("tick city fixture bounds")
		}
		a := tickWorld + 0x840 + city.index*32
		mem[a], mem[a+1] = city.flags, city.owner
		strategyPutWord(mem, a+8, city.x)
		strategyPutWord(mem, a+10, city.y)
		strategyPutWord(mem, a+14, city.production)
		mem[a+0x10], mem[a+0x11] = city.growth, city.prevention
		mem[a+0x12], mem[a+0x13] = city.capacity, city.garrison
		mem[a+0x14], mem[a+0x15], mem[a+0x17] = city.threat, city.marker, city.cooldown
		mem[a+0x18], mem[a+0x19], mem[a+0x1A], mem[a+0x1B] = city.occupancy, city.governor, city.oldOwner, city.neighborCount
		copy(mem[a+0x1C:a+0x20], city.neighbors[:])
		mem[mainPhysical(v.occupancyBase+city.y*24, city.x)] = city.occupancy
	}
	for _, city := range v.cities {
		if city.governor != 255 {
			g := tickWorld + 0x4240 + int(city.governor)*32
			mem[g+0x1A], mem[g+0x11], mem[g+0x13] = v.governorBudget, v.governorWar, v.governorAbility
		}
	}
	owner := c.owner
	if owner < 0 || owner >= 22 {
		panic("tick fixture foreign owner")
	}
	for f := 0; f < 22; f++ {
		mem[tickWorld+f*64+0x19] = v.ally
		for other := 0; other < 24; other++ {
			mem[tickWorld+0x600+f*24+other] = v.relation
		}
	}
	faction := tickWorld + owner*64
	mem[faction+0x14] = v.armyCount
	strategyPutWord(mem, faction+0x21, v.capacityWord)
	for i, value := range v.reserves {
		strategyPutWord(mem, faction+4+i*2, value)
	}
	for i := 0; i < 128; i++ {
		mem[tickWorld+0x2240+i*64] = 0
	}
	for _, army := range v.armies {
		if army.slot < 0 || army.slot >= 128 {
			panic("tick army fixture slot")
		}
		a := tickWorld + 0x2240 + army.slot*64
		mem[a], mem[a+1], mem[a+0x20], mem[a+0x23] = army.flags, army.owner, army.target, army.stage
		strategyPutWord(mem, a+0x0E, army.node)
	}
	if c.group == "reinforce" || v.dispatch {
		// 145C1 checks owner/duty/war, not an existence flag or ruler identity.
		for i := 0; i < 127; i++ {
			mem[tickWorld+0x4240+i*32+0x17] = 255
		}
		if v.dispatch {
			mem[faction+3] = 191
			strategyPutWord(mem, tickWorld+0x840+191*32+8, 40)
			strategyPutWord(mem, tickWorld+0x840+191*32+10, 160)
			for i, war := range []byte{100, 80} {
				g := tickWorld + 0x4240 + i*32
				if c.profile == 1 {
					war = 100
					mem[g] = 0
				}
				mem[g+0x1C], mem[g+0x17], mem[g+0x11] = byte(owner), 0, war
			}
			if c.profile == 2 {
				mem[tickWorld+0x4240+0x17], mem[tickWorld+0x4260+0x17] = 255, 255
				g := tickWorld + 0x4240 + 127*32
				mem[g+0x1C], mem[g+0x17], mem[g+0x11] = byte(owner), 0, 255
			}
		}
	}
	for i := 0; i < 32; i++ {
		mem[tickWorld+0x2040+i*16] = 0
	}
	for _, object := range v.objects {
		if object.slot < 0 || object.slot >= 32 {
			panic("tick object fixture slot")
		}
		a := tickWorld + 0x2040 + object.slot*16
		mem[a], mem[a+12], mem[a+13] = object.flags, object.timer, object.reload
		for i, value := range []uint16{object.x, object.y, object.vx, object.vy} {
			offset := []int{2, 4, 8, 10}[i]
			strategyPutWord(mem, a+offset, value)
		}
	}
	copy(mem[mainPhysical(c.entrySS, c.bp):mainPhysical(c.entrySS, c.bp)+16], v.candidates[:])
	if c.altDS {
		strategyPutWord(mem, 0x26000+0x0D52, 0x7000)
	}
}

type tickModel struct {
	mem                    []byte
	cs, ss                 uint16
	calls                  map[uint16]int
	touched, legacy        map[int]bool
	spawns                 []int
	carry, checkCarry      bool
	ax, dx                 uint16
	ch                     byte
	checkVelocity, checkCH bool
}

func (m *tickModel) enter(ip uint16)      { m.calls[ip]++ }
func (m *tickModel) put(a int, b byte)    { m.mem[a] = b; m.touched[a] = true }
func (m *tickModel) word(a int, v uint16) { m.put(a, byte(v)); m.put(a+1, byte(v>>8)) }
func (m *tickModel) random() byte {
	m.enter(0xECE0)
	a := int(m.cs)*16 + 0xECFC
	value := m.mem[a+2+int(m.mem[a+1])] + m.mem[a]
	m.mem[a] += 0x89
	m.mem[a+1] = value
	return value
}

func (m *tickModel) velocity(value uint16) (uint16, uint16) {
	m.enter(0x24FF)
	lo, hi := byte(value), byte(value>>8)
	r := m.random() & 7
	if r <= 2 {
		hi += r - 1
	}
	if int8(hi) < -15 {
		hi = 0xF1
	}
	if int8(hi) > 15 {
		hi = 15
	}
	lo += hi
	step := uint16(0)
	if int8(lo) >= 15 {
		lo -= 15
		step = 1
	} else if int8(lo) <= -15 {
		lo += 15
		step = 0xFFFF
	}
	return uint16(hi)<<8 | uint16(lo), step
}

func (m *tickModel) object(slot int) {
	m.enter(0x248A)
	a := tickWorld + 0x2040 + slot*16
	vx, sx := m.velocity(strategyWord(m.mem, a+8))
	m.word(a+8, vx)
	m.word(a+2, strategyWord(m.mem, a+2)+sx)
	vy, sy := m.velocity(strategyWord(m.mem, a+10))
	m.word(a+10, vy)
	m.word(a+4, strategyWord(m.mem, a+4)+sy)
	code := int(m.cs) * 16
	for axis, offset := range []int{2, 4} {
		position := strategyWord(m.mem, a+offset)
		direction := byte(0)
		if int16(position) < int16(strategyWord(m.mem, code+0x0D22+axis*2)) {
			direction = 1
		}
		if int16(position) > int16(strategyWord(m.mem, code+0x0D26+axis*2)) {
			direction = 255
		}
		max := int16(400)
		if axis == 1 {
			max = 272
		}
		if int16(position) < -16 {
			position = uint16(max)
		} else if int16(position) > max {
			position = 0xFFF0
		}
		m.word(a+offset, position)
		m.put(a+9+axis*2, m.mem[a+9+axis*2]+direction)
	}
}

func (m *tickModel) objects() {
	m.enter(0x2459)
	for slot := 0; slot < 32; slot++ {
		a := tickWorld + 0x2040 + slot*16
		if m.mem[a] < 0x80 {
			continue
		}
		m.put(a+12, m.mem[a+12]-1)
		if m.mem[a+12] != 0 {
			continue
		}
		m.put(a+12, m.mem[a+13])
		m.put(a, m.mem[a]|1)
		if slot >= 16 {
			m.object(slot)
		}
	}
}

func (m *tickModel) growth(index int) {
	m.enter(0x4194)
	a := tickWorld + 0x840 + index*32
	chance, men := byte(8), byte(4)
	if m.mem[a+1] == m.mem[int(m.cs)*16+0x0CFF] {
		chance, men = 5, 1
		if g := m.mem[a+0x19]; g != 255 {
			general := tickWorld + 0x4240 + int(g)*32
			if m.mem[general+0x1A] != 0 {
				m.put(general+0x1A, m.mem[general+0x1A]-1)
				chance += m.mem[general+0x13]
				men = (men + m.mem[general+0x11]) >> 1
			}
		}
	}
	amount := chance - 15
	if chance <= 15 {
		amount = 1
	}
	if m.random()&15 <= chance {
		value := m.mem[a+0x10] + amount
		if value > 200 {
			value = 200
		}
		m.put(a+0x10, value)
	}
	if m.random()&15 <= chance {
		amount = amount/2 + 1
		value := m.mem[a+0x11] + amount
		if value > 200 {
			value = 200
		}
		m.put(a+0x11, value)
	}
	capacity, garrison := m.mem[a+0x12], m.mem[a+0x13]
	if garrison < capacity {
		if m.random() >= 24 {
			return
		}
		growth := m.mem[a+0x10]
		if growth < men {
			growth = 0
		} else {
			growth -= men
		}
		m.put(a+0x10, growth)
		if int(garrison)+int(men) > 255 {
			garrison = 255
		} else {
			garrison += men
		}
	}
	if garrison > capacity {
		garrison = capacity
	}
	m.put(a+0x13, garrison)
}

func (m *tickModel) disaster(index int) {
	m.enter(0x4269)
	a := tickWorld + 0x840 + index*32
	marker, prevention := m.mem[a+0x15], m.mem[a+0x11]
	if prevention >= marker {
		m.put(a+0x11, prevention-marker)
		return
	}
	deficit := marker - prevention
	m.put(a+0x11, 0)
	growth := m.mem[a+0x10]
	if growth < deficit {
		growth = 0
	} else {
		growth -= deficit
	}
	m.put(a+0x10, growth)
	loss := (uint16(deficit) * uint16(m.mem[a+0x0F])) >> 2
	m.word(a+0x0E, strategyWord(m.mem, a+0x0E)-loss)
	garrison := m.mem[a+0x13]
	if garrison < deficit/2 {
		garrison = 0
	} else {
		garrison -= deficit / 2
	}
	m.put(a+0x13, garrison)
}

func (m *tickModel) neighbors(index int, frame uint16) byte {
	m.enter(0x3FA9)
	a := tickWorld + 0x840 + index*32
	if m.mem[a+0x1B] == 0 {
		return 0
	}
	owner := m.mem[a+1]
	ally := m.mem[tickWorld+int(owner)*64+0x19]
	ch := byte(0)
	at := mainPhysical(m.ss, frame)
	for slot := 0; slot < 4; slot++ {
		id := m.mem[a+0x1C+slot]
		if id == 255 {
			break
		}
		if id >= 192 {
			panic("tick neighbor fixture has invalid record")
		}
		n := tickWorld + 0x840 + int(id)*32
		other := m.mem[n+1]
		if other == owner {
			continue
		}
		if other != 24 {
			if m.mem[tickWorld+0x600+int(owner)*24+int(other)] >= 0x80 {
				continue
			}
			m.put(at, 0xFE)
			ch += m.mem[n+0x18]
		}
		if ally == other {
			m.put(at, id)
			m.put(at+1, other)
			m.put(at+2, m.mem[n+0x18]+1)
			at += 4
		}
	}
	return ch
}

func (m *tickModel) queue(index int) {
	m.enter(0x40B3)
	a := tickWorld + 0x840 + index*32
	m.put(tickWorld+int(m.mem[a+1])*64+0x16, byte(index))
}

func (m *tickModel) spawn(owner byte) (uint16, bool) {
	m.enter(0x45C1)
	best, war := -1, byte(0)
	for g := 0; g < 127; g++ {
		a := tickWorld + 0x4240 + g*32
		if m.mem[a+0x1C] == owner && m.mem[a+0x17] == 0 && m.mem[a+0x11] > war {
			best, war = g, m.mem[a+0x11]
		}
	}
	if best < 0 {
		return 0, false
	}
	m.enter(0x6E8F)
	m.spawns = append(m.spawns, best)
	army := tickWorld + 0x2240 + best*64
	// The old16E8F reserve selector needs six50-unit admissions. It may write
	// partial types on failure. That old six-slot layout is outside this new
	// model, and is still compared in full between original and native C.
	for i := 0; i < 64; i++ {
		m.legacy[army+i] = true
	}
	m.put(army+2, byte(best))
	faction := tickWorld + int(owner)*64
	available := 0
	for i := 0; i < 3; i++ {
		available += int(strategyWord(m.mem, faction+4+i*2) / 50)
	}
	if available < 6 {
		return uint16(army - tickWorld), false
	}
	for i := 0; i < 6; i++ {
		m.legacy[faction+4+i] = true
	}
	m.put(tickWorld+0x4240+best*32+0x17, 1)
	if m.mem[army] < 0x80 {
		m.put(faction+0x14, m.mem[faction+0x14]+1)
	}
	m.put(army, 0xC4)
	m.put(army+1, owner)
	capital := m.mem[faction+3]
	m.put(army+0x20, capital)
	m.put(army+0x23, 1)
	city := tickWorld + 0x840 + int(capital)*32
	x, y := strategyWord(m.mem, city+8), strategyWord(m.mem, city+10)
	segment := strategyWord(m.mem, int(m.cs)*16+0x9872) + uint16(byte(y))*24
	at := mainPhysical(segment, x)
	m.put(at, m.mem[at]+1)
	return uint16(army - tickWorld), true
}

func (m *tickModel) reinforce(owner byte, need, target byte) (bool, byte) {
	m.enter(0x4575)
	faction := tickWorld + int(owner)*64
	value := strategyWord(m.mem, faction+0x21)
	limit := byte(5)
	if int16(value) > 160 {
		limit = byte((value << 3) >> 8)
	}
	count := m.mem[faction+0x14]
	if limit <= count {
		return false, need
	}
	limit -= count
	if need > limit {
		need = limit
	}
	remaining := need
	made := false
	for tries := 0; tries < 128; tries++ {
		army, ok := m.spawn(owner)
		if !ok {
			return made, remaining
		}
		made = true
		m.put(tickWorld+int(army)+0x20, target)
		m.put(tickWorld+int(army)+0x23, 0)
		remaining--
		if remaining == 0 {
			return true, 0
		}
	}
	panic("tick reinforce model exceeded original127-general candidate bound")
}

func (m *tickModel) report(index int, need byte) {
	m.enter(0x40C9)
	a := tickWorld + 0x840 + index*32
	if m.mem[a+0x17] != 0 {
		return
	}
	owner := m.mem[a+1]
	if owner == m.mem[int(m.cs)*16+0x0CFF] {
		m.enter(0x0CDE)
		m.enter(0xEB11)
		m.enter(0x8810)
		m.put(a+0x17, (m.random()&15)+24)
		return
	}
	ok, _ := m.reinforce(owner, need, byte(index))
	if !ok {
		return
	}
	capital := int(m.mem[tickWorld+int(owner)*64+3])
	other := tickWorld + 0x840 + capital*32
	abs := func(a, b uint16) uint16 {
		if a < b {
			return b - a
		}
		return a - b
	}
	x := abs(strategyWord(m.mem, a+8), strategyWord(m.mem, other+8))
	// Original14137 really reads the capital's +848 X again for the Y term.
	y := abs(strategyWord(m.mem, a+10), strategyWord(m.mem, other+8))
	delay := (x + y) >> 3
	if delay > 30 {
		delay = 30
	}
	m.put(a+0x17, byte(delay))
}

func (m *tickModel) threatGate(index int, frame uint16) bool {
	m.enter(0x4028)
	a := tickWorld + 0x840 + index*32
	m.put(a, m.mem[a]&0x3F)
	first := m.mem[mainPhysical(m.ss, frame)]
	if first > 0xFE {
		m.put(a+0x17, 0)
		return false
	}
	flag := byte(0x80)
	if first < 0xFE {
		flag |= 0x40
	}
	m.put(a, m.mem[a]|flag)
	if m.mem[a+0x18] >= 1 {
		return true
	}
	m.report(index, 1)
	m.queue(index)
	return false
}

func (m *tickModel) assign(index int, target, stage, need, total byte) {
	m.enter(0x4155)
	extra := total - need
	for slot := 0; slot < 128; slot++ {
		a := tickWorld + 0x2240 + slot*64
		if strategyWord(m.mem, a+0x0E) != uint16(index*8) || m.mem[a] < 0x80 {
			continue
		}
		if extra != 0 && m.random() < 64 {
			extra--
			continue
		}
		if m.mem[a]&4 != 0 && m.mem[a+0x23] < 8 {
			m.put(a+0x20, target)
			m.put(a+0x23, stage)
		}
		need--
		if need == 0 {
			return
		}
	}
	panic("tick14155 fixture cannot fulfill its bounded original scan")
}

func (m *tickModel) threatChoice(index int, frame uint16) {
	m.enter(0x4057)
	start := mainPhysical(m.ss, frame)
	if m.mem[start] >= 0xFE {
		return
	}
	n := 0
	for n < 4 && m.mem[start+n*4] < 0xFE {
		n++
	}
	if n == 0 {
		panic("tick candidate model empty")
	}
	choice := m.random() & 3
	steps := int(choice)
	if steps == 0 {
		steps = 256
	}
	selected := start + ((steps-1)%n)*4
	a := tickWorld + 0x840 + index*32
	if m.mem[a+0x18] <= 1 {
		need := m.mem[a+0x14] + 2
		borrow := need < m.mem[a+0x18]
		need -= m.mem[a+0x18]
		if borrow || need == 0 || m.mem[a+1] == m.mem[int(m.cs)*16+0x0CFF] {
			return
		}
		m.report(index, need)
		m.queue(index)
		return
	}
	// The original decrements AL until zero while choosing its candidate.
	// This branch therefore converts the resulting zero to a request of one.
	m.assign(index, m.mem[selected], 0, 1, m.mem[a+0x18])
	m.put(a+0x17, 0)
}

func (m *tickModel) threat(index int, frame uint16) {
	m.enter(0x3F74)
	start := mainPhysical(m.ss, frame)
	for i := 0; i < 16; i++ {
		m.put(start+i, 255)
	}
	value := m.neighbors(index, frame)
	m.put(tickWorld+0x840+index*32+0x14, value)
	if m.threatGate(index, frame) {
		m.threatChoice(index, frame)
	}
}

func (m *tickModel) city(index int, frame uint16) {
	m.enter(0x3EFD)
	a := tickWorld + 0x840 + index*32
	if m.mem[a+0x17] != 0 {
		m.put(a+0x17, m.mem[a+0x17]-1)
	}
	if old := m.mem[a+0x1A]; old != m.mem[a+1] {
		m.put(tickWorld+int(old)*64+0x17, byte(index))
	}
	segment := strategyWord(m.mem, int(m.cs)*16+0x9872) + strategyWord(m.mem, a+10)*24
	m.put(a+0x18, m.mem[mainPhysical(segment, strategyWord(m.mem, a+8))]&0x7F)
	if m.mem[a+1] != 24 {
		m.threat(index, frame)
	}
	m.growth(index)
	m.disaster(index)
	next := uint16(index*32) + 32
	if next >= 0x1800 {
		next = 0
	}
	m.word(int(m.cs)*16+0x0D1E, next)
}

func tickExpected(before []byte, cs uint16, c resumeCase) *tickModel {
	m := &tickModel{mem: append([]byte(nil), before...), cs: cs, ss: c.entrySS, calls: map[uint16]int{}, touched: map[int]bool{}, legacy: map[int]bool{}}
	v := c.tick
	switch c.target {
	case 0x13EFD:
		m.city(int(strategyWord(before, int(cs)*16+0x0D1E))/32, c.entrySP-32)
	case 0x13F74:
		m.threat(v.city, c.entrySP-30)
	case 0x13FA9:
		m.ch = m.neighbors(v.city, c.bp)
		m.checkCH = true
	case 0x14028:
		m.carry = m.threatGate(v.city, c.bp)
		m.checkCarry = true
	case 0x14057:
		m.threatChoice(v.city, c.bp)
	case 0x140B3:
		m.queue(v.city)
	case 0x140C9:
		m.report(v.city, byte(c.ax))
	case 0x14155:
		m.assign(v.city, byte(c.cx), byte(c.cx>>8), byte(c.dx), byte(c.dx>>8))
	case 0x14194:
		m.growth(v.city)
	case 0x14269:
		m.disaster(v.city)
	case 0x14575:
		ok, _ := m.reinforce(byte(c.owner), byte(c.cx), byte(c.cx>>8))
		m.carry, m.checkCarry = !ok, true
	case 0x145C1:
		_, ok := m.spawn(byte(c.ax))
		m.carry, m.checkCarry = !ok, true
	case 0x12459:
		m.objects()
	case 0x1248A:
		m.object(v.objects[0].slot)
	case 0x124FF:
		m.dx, m.ax = m.velocity(c.dx)
		m.checkVelocity = true
	case 0x10CDE:
		m.enter(0x0CDE)
		m.enter(0xEB11)
	case 0x1EB11:
		m.enter(0xEB11)
	case 0x1EB5E:
	default:
		panic("unknown independent tick root")
	}
	return m
}

// IN values are an observed, explicit platform-fixture input. Native C never
// consumes tickOriginalPortReads. RNG/world expectations use only before RAM.
// This checks pulse and poll order; it makes no hardware wall-clock claim.
func tickAlertAudit(c resumeCase, m *tickModel, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, what string) {
		checks++
		if !ok {
			panic("independent tick alert: " + what)
		}
	}
	type io struct {
		port  uint16
		value byte
	}
	var reads, writes []io
	check(len(tickOriginalPortReads)%3 == 0, "complete IN fixture records")
	for i := 0; i < len(tickOriginalPortReads); i += 3 {
		port := uint16(tickOriginalPortReads[i]) | uint16(tickOriginalPortReads[i+1])<<8
		if port == 0x3DA || port == 0x61 {
			reads = append(reads, io{port, tickOriginalPortReads[i+2]})
		}
	}
	at, polls := 0, 0
	read := func(port uint16) byte {
		check(at < len(reads), "missing original platform IN")
		r := reads[at]
		at++
		check(r.port == port, fmt.Sprintf("IN order expected%03X got%03X", port, r.port))
		if port == 0x3DA {
			polls++
		}
		return r.value
	}
	poll := func(high bool) {
		for tries := 0; tries < 1024; tries++ {
			if (read(0x3DA)&8 != 0) == high {
				return
			}
		}
		panic("fixed tick platform did not terminate vertical-status poll")
	}
	if c.target == 0x1EB5E {
		value := read(0x3DA)
		check(byte(device.CPU.R[cpu.AX]) == value, "status helper AL")
		check((device.CPU.Flags&cpu.CF != 0) == (value&8 != 0), "status helper CF")
	} else {
		for call := 0; call < m.calls[0xEB11]; call++ {
			value := uint16(0x0101)
			if c.target == 0x1EB11 {
				value = c.ax
			}
			cl, ch := byte(value), byte(value>>8)
			for pulse := 0; pulse < 256; pulse++ {
				poll(true)
				writes = append(writes, io{0x61, read(0x61) | 3})
				for _, high := range []bool{false, true, false, true} {
					poll(high)
				}
				writes = append(writes, io{0x61, read(0x61) & 0xFC})
				if cl == 1 {
					break
				}
				delays := int(ch)
				if delays == 0 {
					delays = 256
				}
				for delay := 0; delay < delays; delay++ {
					for _, high := range []bool{false, true, false, true} {
						poll(high)
					}
				}
				cl--
				if cl == 0 {
					break
				}
			}
		}
	}
	check(at == len(reads), "extra IN beyond original alert sequence")
	check(len(strategyObserved(observed, 0xEB5E)) == polls, "status helper count equals original IN polls")
	var got []io
	for _, w := range device.PortLog {
		if w.Port == 0x61 {
			got = append(got, io{w.Port, w.Val})
		}
	}
	check(len(got) == len(writes), "speaker OUT sequence length")
	for i, want := range writes {
		check(got[i] == want, fmt.Sprintf("speaker OUT%d", i))
	}
	for _, r := range strategyObserved(observed, 0xEB11) {
		want := uint16(0x0101)
		if c.target == 0x1EB11 {
			want = c.ax
		}
		check(r.AX == want, "original alert AX")
	}
	return checks
}

func tickAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, what string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent tick group=%s target=%X: %s", c.group, c.target, what))
		}
	}
	m := tickExpected(before, cs, c)
	code := int(cs) * 16
	var raw [258]byte
	raw[0], raw[1] = c.tick.rngSeed^0xA5, c.tick.rngSeed
	for i := 0; i < 256; i++ {
		raw[2+i] = byte(i*73 + i/2 + int(c.tick.rngSeed))
	}
	check(bytes.Equal(before[code+0xECFC:code+0xEDFE], raw[:]), "enumerated raw RNG initial-state contract")
	first := -1
	for i := 0; i < tickWorldSize; i++ {
		a := tickWorld + i
		if (!m.legacy[a] || m.touched[a]) && after[a] != m.mem[a] {
			first = i
			break
		}
	}
	check(first < 0, fmt.Sprintf("world model first differing byte%04X", first))
	occupancy := int(strategyWord(before, code+0x9872)) * 16
	check(bytes.Equal(after[occupancy:occupancy+tickGridSize], m.mem[occupancy:occupancy+tickGridSize]), "complete occupancy grid")
	check(strategyWord(after, code+0x9872) == strategyWord(before, code+0x9872), "occupancy pointer preserved")
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "game date changed")
	check(strategyWord(after, code+0x0D1E) == strategyWord(m.mem, code+0x0D1E), "city cursor and wrap")
	check(bytes.Equal(after[code+0xECFC:code+0xEDFE], m.mem[code+0xECFC:code+0xEDFE]), "full258-byte RNG recurrence and immutable table")
	for _, ip := range []uint16{0x0CDE, 0x2459, 0x248A, 0x24FF, 0x3EFD, 0x3F74, 0x3FA9, 0x4028, 0x4057, 0x40B3, 0x40C9, 0x4155, 0x4194, 0x4269, 0x4575, 0x45C1, 0xEB11, 0xECE0, 0x6E8F} {
		check(len(strategyObserved(observed, ip)) == m.calls[ip], fmt.Sprintf("entry%04X count want%d got%d", ip, m.calls[ip], len(strategyObserved(observed, ip))))
	}
	if m.checkCarry {
		check((device.CPU.Flags&cpu.CF != 0) == m.carry, "original root carry")
	}
	if m.checkVelocity {
		check(device.CPU.R[cpu.AX] == m.ax && device.CPU.R[cpu.DX] == m.dx, "velocity AX/DX signed thresholds")
	}
	if m.checkCH {
		check(byte(device.CPU.R[cpu.CX]>>8) == m.ch, "neighbor total byte wrap")
	}
	if c.target == 0x13FA9 || c.target == 0x14028 || c.target == 0x14057 {
		frame := mainPhysical(c.entrySS, c.bp)
		check(bytes.Equal(after[frame:frame+16], m.mem[frame:frame+16]), "direct helper candidate frame")
	}
	for i, g := range m.spawns {
		r := strategyObserved(observed, 0x6E8F)[i]
		check(r.SI == uint16(0x4240+g*32) && r.DS == 0x7000, "original maximal-war first-tie candidate")
	}
	warnings := strategyObserved(observed, 0x8810)
	check(len(warnings) == m.calls[0x8810], "real player warning call count")
	for _, r := range warnings {
		check(byte(r.AX) == 0x93 && r.CX == 0x26 && r.DS == 0x7000, "real TALK26 portrait93/world DS")
	}
	if c.target == 0x13EFD || c.target == 0x12459 {
		check(device.CPU.Seg[cpu.DS] == cs, "scheduler helper final DS=CS")
	}
	checks += tickAlertAudit(c, m, observed, device)
	return checks
}

func tickCases(base func(string, uint32, int) resumeCase) []resumeCase {
	var all []resumeCase
	city := func(index int, owner byte) tickCity {
		return tickCity{index: index, flags: 0xC3, owner: owner, oldOwner: owner, governor: 255,
			x: 100, y: 200, production: 1024, growth: 100, prevention: 100, capacity: 30, garrison: 30, neighbors: [4]byte{255, 255, 255, 255}}
	}
	makeCase := func(group string, target uint32, scenario int) resumeCase {
		c := base(group, target, scenario)
		c.profile, c.altDS = 0, false
		c.tick = tickVector{occupancyBase: 0xB000, rngSeed: 3, ally: 255, relation: 0x80,
			bounds: [4]int16{-16, -16, 400, 272}, cities: []tickCity{city(0, c.player)}}
		for i := range c.tick.candidates {
			c.tick.candidates[i] = 255
		}
		c.ax, c.cx, c.dx = 1, 0x0101, 0x0301
		c.events = []listEvent{{100, 120, 1}}
		c.mapEvents = nil
		c.menuChoices = nil
		return c
	}
	add := func(c resumeCase) { all = append(all, c) }
	armyFixture := func(c *resumeCase) {
		for i := 0; i < 16; i++ {
			c.tick.armies = append(c.tick.armies, tickArmy{slot: i, flags: 0xC4, stage: byte(i % 10), owner: byte(c.owner), target: 0xAA, node: uint16(c.tick.city * 8)})
		}
	}
	for scenario := 0; scenario < 4; scenario++ {
		for index := 0; index < 192; index++ {
			c := makeCase("city", 0x13EFD, scenario)
			c.tick.city, c.tick.cursor = index, uint16(index*32)
			c.tick.cities = []tickCity{city(index, 24)}
			c.tick.cities[0].occupancy = 0x80 | byte(index&127)
			c.tick.rngSeed = byte(index)
			add(c)
		}
		for profile := 0; profile < 12; profile++ {
			c := makeCase("city", 0x13EFD, scenario)
			index := 0
			if profile%2 != 0 {
				index = 191
			}
			c.tick.city, c.tick.cursor = index, uint16(index*32)
			owner := []byte{c.player, byte(c.owner), 24}[profile%3]
			c.tick.cities = []tickCity{city(index, owner)}
			p := &c.tick.cities[0]
			p.oldOwner = byte(c.owner)
			p.cooldown = []byte{0, 1, 255}[profile%3]
			p.occupancy = []byte{0, 127, 128, 255}[profile%4]
			if profile >= 8 {
				c.tick.occupancyBase = 0xC000
			}
			add(c)
		}
		for _, relation := range []byte{0, 0x7F, 0x80, 0xFF} {
			for _, allyKind := range []int{0, 1, 2} {
				for rotation := 0; rotation < 4; rotation++ {
					c := makeCase("neighbors", 0x13FA9, scenario)
					owner := byte(c.owner)
					c.tick.cities = []tickCity{city(0, owner)}
					c.tick.cities[0].neighborCount = 1
					c.tick.relation = relation
					c.tick.ally = []byte{c.player, owner, 24}[allyKind]
					for i, other := range []byte{c.player, owner, 24, c.player} {
						id := i + 1
						n := city(id, other)
						n.flags = 0
						n.occupancy = []byte{127, 255, 1, 130}[i]
						c.tick.cities = append(c.tick.cities, n)
						c.tick.cities[0].neighbors[i] = byte(1 + (i+rotation)%4)
					}
					add(c)
				}
			}
		}
		for _, flags := range []byte{0, 0xC0} {
			for _, count := range []byte{0, 1} {
				for stop := 0; stop < 4; stop++ {
					c := makeCase("neighbors", 0x13FA9, scenario)
					c.tick.cities = []tickCity{city(0, byte(c.owner))}
					c.tick.cities[0].flags, c.tick.cities[0].neighborCount = flags, count
					c.tick.ally, c.tick.relation = c.player, 0x7F
					for i := 0; i < 4; i++ {
						n := city(i+1, c.player)
						n.occupancy = 255
						n.flags = flags
						c.tick.cities = append(c.tick.cities, n)
						c.tick.cities[0].neighbors[i] = byte(i + 1)
					}
					c.tick.cities[0].neighbors[stop] = 255
					add(c)
				}
			}
		}
		for _, target := range []uint32{0x13F74, 0x14028, 0x14057, 0x140B3, 0x140C9, 0x14155} {
			for profile := 0; profile < 8; profile++ {
				c := makeCase("threat", target, scenario)
				c.tick.city = profile * 23
				owner, enemy := c.player, byte(c.owner)
				if profile%2 != 0 {
					owner, enemy = byte(c.owner), c.player
				}
				c.tick.cities = []tickCity{city(c.tick.city, owner), city(1, enemy), city(2, enemy)}
				p := &c.tick.cities[0]
				p.neighborCount = 1
				p.neighbors = [4]byte{1, 2, 255, 255}
				p.occupancy = byte(profile % 4)
				p.threat = byte(profile * 41)
				p.cooldown = byte(profile % 2)
				c.tick.cities[1].occupancy, c.tick.cities[2].occupancy = 127, 255
				c.tick.ally, c.tick.relation = enemy, 0x7F
				switch profile % 3 {
				case 0:
					c.tick.candidates[0] = 255
				case 1:
					c.tick.candidates[0] = 0xFE
				case 2:
					c.tick.candidates[0], c.tick.candidates[1], c.tick.candidates[2] = 1, enemy, 255
				}
				c.tick.rngSeed = byte(profile * 37)
				armyFixture(&c)
				if target == 0x140C9 && profile == 7 {
					p.owner, p.cooldown, p.x, p.y = byte(c.owner), 0, 100, 200
					c.tick.dispatch = true
					c.tick.reserves = [3]uint16{1000, 1000, 1000}
					c.tick.armies = nil
				}
				add(c)
			}
		}
		for _, ownerKind := range []int{0, 1, 2} {
			for _, budget := range []byte{0, 1, 2, 255} {
				for profile := 0; profile < 8; profile++ {
					c := makeCase("growth", 0x14194, scenario)
					owner := []byte{c.player, byte(c.owner), 24}[ownerKind]
					c.tick.cities = []tickCity{city(0, owner)}
					p := &c.tick.cities[0]
					p.governor = byte(c.general)
					if profile == 7 {
						p.governor = 255
					}
					p.growth = []byte{0, 1, 199, 200, 250, 255, 127, 254}[profile]
					p.prevention = byte(profile*37 + 250)
					pair := [][2]byte{{0, 255}, {1, 0}, {10, 9}, {200, 0}, {255, 254}, {255, 0}, {100, 200}, {255, 255}}[profile]
					p.capacity, p.garrison = pair[0], pair[1]
					c.tick.governorBudget = budget
					c.tick.governorAbility = []byte{0, 10, 11, 250, 251, 255, 127, 200}[profile]
					c.tick.governorWar = []byte{0, 1, 15, 255, 254, 127, 128, 200}[profile]
					c.tick.rngSeed = byte(profile * 31)
					add(c)
				}
			}
		}
		for seed := 0; seed < 256; seed++ {
			c := makeCase("growth", 0x14194, scenario)
			c.tick.rngSeed = byte(seed)
			c.tick.cities[0].capacity, c.tick.cities[0].garrison = 255, 0
			if seed%2 != 0 {
				p := &c.tick.cities[0]
				p.governor, p.growth, p.prevention, p.garrison = byte(c.general), 255, 255, 254
				c.tick.governorBudget, c.tick.governorWar = 1, 254
				c.tick.governorAbility = []byte{250, 251, 255, 200}[(seed/2)%4]
			}
			add(c)
		}
		for _, marker := range []byte{0, 1, 2, 4, 11, 127, 254, 255} {
			for _, prevention := range []byte{0, 1, 2, 4, 11, 127, 254, 255} {
				for _, production := range []uint16{0, 1, 255, 256, 65535} {
					c := makeCase("disaster", 0x14269, scenario)
					p := &c.tick.cities[0]
					p.marker, p.prevention, p.production = marker, prevention, production
					p.growth, p.garrison = prevention, marker/3
					add(c)
				}
			}
		}
		for _, target := range []uint32{0x14575, 0x145C1} {
			for _, dispatch := range []bool{false, true} {
				for _, capacity := range []uint16{0, 160, 161, 0x7FFF, 0xFFFF} {
					for _, count := range []byte{0, 4, 5, 255} {
						for _, need := range []byte{1, 2} {
							c := makeCase("reinforce", target, scenario)
							c.tick.dispatch, c.tick.capacityWord, c.tick.armyCount = dispatch, capacity, count
							c.tick.reserves = [3]uint16{1000, 1000, 1000}
							c.ax = uint16(c.owner)
							c.cx = uint16(need) | 0x0500
							add(c)
						}
					}
				}
			}
		}
		for _, target := range []uint32{0x14575, 0x145C1} {
			for _, profile := range []int{1, 2} {
				c := makeCase("reinforce", target, scenario)
				c.profile, c.tick.dispatch = profile, true
				c.tick.reserves = [3]uint16{1000, 1000, 1000}
				c.ax, c.cx = uint16(c.owner), 0x0502
				add(c)
			}
			for _, reserves := range [][3]uint16{{0, 0, 0}, {50, 0, 0}} {
				c := makeCase("reinforce", target, scenario)
				c.tick.dispatch, c.tick.reserves = true, reserves
				c.ax, c.cx = uint16(c.owner), 0x0501
				add(c)
			}
		}
		c := makeCase("reinforce", 0x14575, scenario)
		c.tick.dispatch = true
		c.tick.reserves = [3]uint16{1000, 1000, 1000}
		c.cx = 0x0503
		add(c)
		for slot := 0; slot < 32; slot++ {
			for _, flags := range []byte{0, 0x7F, 0x80, 0xC0, 0xFF} {
				for _, timer := range []byte{0, 1, 2, 255} {
					c := makeCase("objects", 0x12459, scenario)
					c.tick.objects = []tickObject{{slot: slot, flags: flags, timer: timer, reload: byte(slot * 9), x: uint16(slot * 17), y: uint16(slot * 11), vx: 0x0F0E, vy: 0xF1F2}}
					add(c)
				}
			}
		}
		c = makeCase("objects", 0x12459, scenario)
		for slot := 0; slot < 32; slot++ {
			c.tick.objects = append(c.tick.objects, tickObject{slot: slot, flags: 0x80 | byte(slot&1), timer: 1, reload: byte(slot * 7), x: uint16(slot * 17), y: uint16(slot * 11), vx: 0x0F0E, vy: 0xF1F2})
		}
		add(c)
		for _, xy := range [][2]uint16{{0, 0}, {0xFFEF, 0xFFEF}, {0xFFF0, 0xFFF0}, {400, 272}, {401, 273}, {0x7FFF, 0x8000}} {
			for _, velocity := range []uint16{0, 0x0F0F, 0xF1F1, 0x1000, 0xF000, 0x7F80} {
				c := makeCase("objects", 0x1248A, scenario)
				c.tick.bounds = [4]int16{-10, -10, 380, 250}
				c.tick.objects = []tickObject{{slot: 16, flags: 0x7F, x: xy[0], y: xy[1], vx: velocity, vy: velocity}}
				add(c)
			}
		}
		for _, hi := range []byte{0, 1, 15, 16, 127, 128, 240, 241} {
			for _, lo := range []byte{0, 1, 14, 15, 16, 127, 240, 241} {
				for _, seed := range []byte{0, 3, 255} {
					c := makeCase("objects", 0x124FF, scenario)
					c.dx = uint16(hi)<<8 | uint16(lo)
					c.tick.rngSeed = seed
					add(c)
				}
			}
		}
		add(makeCase("alert", 0x10CDE, scenario))
		for _, duration := range []uint16{0x0101, 0x0201, 0x0102, 0x0202, 0x0001} {
			c := makeCase("alert", 0x1EB11, scenario)
			c.ax = duration
			add(c)
		}
		for _, ax := range []uint16{0, 0xFFFF, 0x5500, 0xAAFF} {
			c := makeCase("alert", 0x1EB5E, scenario)
			c.ax = ax
			add(c)
		}
	}
	return all
}
