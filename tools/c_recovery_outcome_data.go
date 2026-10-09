//go:build matching_outcome

package main

import (
	"bytes"
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

// Independent spec250 model, IDA9.4 linear base10000, fixed probe
// b20f688289b06ca5384f3ac810cd4b9b93cc5ddb82d5be4d12556701a8988ce9.
// KI.EXE fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868.
type outcomeArmy struct {
	slot                        int
	owner, flags, morale, stage byte
	men, types                  [6]byte
	current, target             uint16
}
type outcomeCity struct {
	index                                                                                    int
	owner, flags, oldOwner, governor, neighborCount, garrison, growth, prevention, kind, art byte
	x, y, production                                                                         uint16
	neighbors                                                                                [4]byte
}
type outcomeGeneral struct {
	index                                                             int
	flags, owner, duty, captor, temperament, war, lead, score, budget byte
	aptitudes                                                         [3]byte
}
type outcomeFaction struct {
	index                                                  int
	flags, ruler, capital, cityCount, foreignOfficer, ally byte
}
type outcomeVector struct {
	escape                                                             bool
	city, general                                                      int
	oldOwner, newOwner, mode, ratio, reportCity, rngSeed, liveFactions byte
	power                                                              uint16
	armies                                                             []outcomeArmy
	cities                                                             []outcomeCity
	generals                                                           []outcomeGeneral
	factions                                                           []outcomeFaction
	redirectSlots                                                      []int
	panel, selectedFaction                                             byte
	markX, markY                                                       uint16
	coefficients                                                       [16]byte
	customCoefficients                                                 bool
}

const outcomeWorld = 0x70000
const outcomeWorldSize = 0x5240

func outcomePrepare(mem []byte, cs uint16, c resumeCase) {
	routePrepare(mem, cs, c)
	o, code := c.outcome, int(cs)*16
	if o.city < 0 || o.city >= 192 || o.general < 0 || o.general >= 128 || len(o.redirectSlots) > 127 {
		panic("outcome typed fixture bounds")
	}
	if c.group == "capture" || c.group == "fall" || c.group == "history" || c.group == "neighbors" {
		for i := 0; i < 192; i++ {
			a := outcomeWorld + 0x840 + i*32
			mem[a+1], mem[a+0x1A], mem[a+0x19] = 24, 24, 255
			mem[a+0x1B] = 0
			for j := 0; j < 4; j++ {
				mem[a+0x1C+j] = 255
			}
		}
		for i := 0; i < 128; i++ {
			mem[outcomeWorld+0x4240+i*32] = 0
			mem[outcomeWorld+0x2240+i*64] = 0
		}
	}
	for _, f := range o.factions {
		if f.index < 0 || f.index >= 22 {
			panic("outcome faction bounds")
		}
		a := outcomeWorld + f.index*64
		mem[a], mem[a+1], mem[a+3], mem[a+0x23], mem[a+0x2A], mem[a+0x19] = f.flags, f.ruler, f.capital, f.cityCount, f.foreignOfficer, f.ally
	}
	for _, city := range o.cities {
		if city.index < 0 || city.index >= 192 || city.x > 383 || city.y > 255 {
			panic("outcome city fixture bounds")
		}
		a := outcomeWorld + 0x840 + city.index*32
		mem[a], mem[a+1], mem[a+0x1A], mem[a+0x19], mem[a+0x1B] = city.flags, city.owner, city.oldOwner, city.governor, city.neighborCount
		mem[a+0x13], mem[a+0x10], mem[a+0x11], mem[a+0x16] = city.garrison, city.growth, city.prevention, city.art<<4|city.kind
		strategyPutWord(mem, a+8, city.x)
		strategyPutWord(mem, a+10, city.y)
		strategyPutWord(mem, a+14, city.production)
		copy(mem[a+0x1C:a+0x20], city.neighbors[:])
		mapBase := strategyWord(mem, code+0x0D44)
		segment := mapBase + (city.y-2)*24
		offset := city.x + 0x300
		mem[mainPhysical(segment, offset)] = 0xCB
		deltas := [4]int{-0x181, -0x17F, 0x181, 0x17F}
		if city.kind&15 == 0 {
			deltas = [4]int{-0x302, -0x2FE, 0x302, 0x2FE}
		} else if city.kind&15 == 3 {
			deltas = [4]int{-0x180, -1, 1, 0x180}
		}
		for i, delta := range deltas {
			mem[mainPhysical(segment, offset+uint16(delta))] = byte(0xDE + i)
		}
	}
	for _, g := range o.generals {
		if g.index < 0 || g.index >= 128 {
			panic("outcome general bounds")
		}
		a := outcomeWorld + 0x4240 + g.index*32
		mem[a], mem[a+0x1C], mem[a+0x17], mem[a+0x1D], mem[a+0x1E], mem[a+0x11], mem[a+0x12], mem[a+0x1F], mem[a+0x1A] = g.flags, g.owner, g.duty, g.captor, g.temperament, g.war, g.lead, g.score, g.budget
		copy(mem[a+0x0E:a+0x11], g.aptitudes[:])
		army := outcomeWorld + 0x2240 + g.index*64
		strategyPutWord(mem, army+0x1A, uint16(100+g.index))
		strategyPutWord(mem, army+0x1C, c.route.occupancyBase+100*24)
	}
	for _, army := range o.armies {
		if army.slot < 0 || army.slot >= 128 || army.owner >= 22 {
			panic("outcome army bounds")
		}
		a := outcomeWorld + 0x2240 + army.slot*64
		mem[a], mem[a+1], mem[a+2], mem[a+6], mem[a+0x23] = army.flags, army.owner, byte(army.slot), army.morale, army.stage
		strategyPutWord(mem, a+0x0E, army.current)
		strategyPutWord(mem, a+0x14, army.target)
		mem[a+0x20] = byte(army.target / 8)
		strategyPutWord(mem, a+0x10, 100)
		strategyPutWord(mem, a+0x12, 100)
		strategyPutWord(mem, a+0x1A, uint16(100+army.slot))
		strategyPutWord(mem, a+0x1C, c.route.occupancyBase+100*24)
		total := uint16(0)
		for slot := 0; slot < 6; slot++ {
			if army.types[slot] < 1 || army.types[slot] > 4 {
				panic("outcome legal troop type")
			}
			mem[a+0x29+slot*4], mem[a+0x2A+slot*4] = army.men[slot], army.types[slot]
			total += uint16(army.men[slot])
		}
		strategyPutWord(mem, a+4, total)
		mem[mainPhysical(c.route.occupancyBase+100*24, uint16(100+army.slot))] = c.route.occupancy
	}
	for i, slot := range o.redirectSlots {
		if slot < 0 || slot >= 128 {
			panic("outcome redirect slot bounds")
		}
		strategyPutWord(mem, mainPhysical(c.entrySS, c.bp+uint16(i*2)), uint16(0x2240+slot*64))
	}
	mem[mainPhysical(c.entrySS, c.bp+0xFE)] = byte(len(o.redirectSlots))
	strategyPutWord(mem, code+0x0D32, uint16(0x840+o.city*32))
	mem[code+0x0D34] = o.reportCity
	mem[code+0x98A6], mem[code+0x98A7], mem[code+0x0D2A] = o.panel, o.selectedFaction, o.liveFactions
	strategyPutWord(mem, code+0x9892, o.markX)
	strategyPutWord(mem, code+0x9894, o.markY)
	mem[code+0xECFC], mem[code+0xECFD] = o.rngSeed^0xA5, o.rngSeed
	for i := 0; i < 256; i++ {
		mem[code+0xECFE+i] = byte(i*73 + i/2 + int(o.rngSeed))
	}
	if o.customCoefficients {
		copy(mem[code+0x5120:code+0x5130], o.coefficients[:])
	}
	copy(mem[code+0x1964:code+0x1994], c.palette[:])
	// Explicit isolated nonlocal fixture, not a claim that bootstrap was run.
	if c.returnIP != 0x006A || c.savedSS == 0 || c.savedSP == 0 {
		panic("outcome fixture needs explicit outer006A saved frame")
	}
	strategyPutWord(mem, code+0x9901, c.savedSP)
	strategyPutWord(mem, code+0x9903, c.savedSS)
	strategyPutWord(mem, mainPhysical(c.savedSS, c.savedSP), c.returnIP)
}

type outcomeModel struct {
	*routeModel
	c                                           resumeCase
	ax, dx                                      uint16
	checkAX, checkDX, checkCarry, carry, escape bool
	unknown                                     map[int]bool
	casualtyCalls                               []registers
	redirects                                   []registers
	clearOrder                                  []uint16
	sounds                                      []uint16
	marks                                       []registers
	saves                                       []registers
}

func (m *outcomeModel) coefficient(army uint16, mode byte) uint16 {
	m.enter(0x5285)
	a := outcomeWorld + int(army)
	sum := uint16(0)
	for i := 0; i < 6; i++ {
		kind := m.mem[a+0x2A+i*4]
		at := int(m.cs)*16 + 0x5120 + int(mode)*4 + int(kind-1)
		sum += uint16(m.mem[a+0x29+i*4]) * uint16(m.mem[at])
	}
	if mode == 0 {
		city := strategyWord(m.mem, int(m.cs)*16+0x0D32)
		sum += uint16(m.mem[outcomeWorld+int(city)+0x13])
	}
	return uint16(sum * uint16(m.mem[a+6]>>3))
}

func (m *outcomeModel) ability(army uint16, mode byte, power uint16) uint16 {
	m.enter(0x52D7)
	a := outcomeWorld + int(army)
	g := outcomeWorld + 0x4240 + int(m.mem[a+2])*32
	war, lead := m.mem[g+0x11], m.mem[g+0x12]
	value := war * 2
	if war < lead {
		value = lead - (lead >> 2) + m.mem[g+0x12]
	} else if m.random()&3 == 0 {
		value = war - (war >> 2) + m.mem[g+0x12]
	}
	nibble := m.mem[g+0x0E+int(mode)] >> 4
	denominator := uint16(16 - nibble)
	factor := (uint16(value) << 4) / denominator
	return uint16((uint32(factor) * uint32(power)) >> 10)
}

func (m *outcomeModel) casualties(winner, loser uint16, ratio byte) {
	m.enter(0x51B3)
	m.casualtyCalls = append(m.casualtyCalls, registers{SI: winner, DI: loser, AX: uint16(ratio)})
	code := int(m.cs) * 16
	if m.mem[code+0x0D34] < 192 {
		a := outcomeWorld + int(strategyWord(m.mem, code+0x0D32))
		damage := byte(63-ratio) >> 2
		for _, offset := range []int{0x13, 0x10, 0x11} {
			if m.mem[a+offset] < damage {
				m.mem[a+offset] = 0
			} else {
				m.mem[a+offset] -= damage
			}
		}
	}
	a, b := outcomeWorld+int(winner), outcomeWorld+int(loser)
	oldA, oldB := strategyWord(m.mem, a+4), strategyWord(m.mem, b+4)
	newA, newB := uint16(0), uint16(0)
	for slot := 0; slot < 6; slot++ {
		lossA := (m.random() & 7) + 2
		menA := m.mem[a+0x29+slot*4]
		if menA <= lossA {
			menA = 0
		} else {
			menA -= lossA
		}
		if slot == 0 && menA == 0 {
			menA = 1
		}
		m.mem[a+0x29+slot*4] = menA
		newA += uint16(menA)
		r := m.random()
		divisor := uint16(ratio) + uint16(slot) + 1
		if divisor > 255 || divisor == 0 {
			panic("outcome casualty fixture invalid byte DIV")
		}
		lossB := byte(uint16(r)%divisor) + 8
		menB := m.mem[b+0x29+slot*4]
		if menB <= lossB {
			menB = 0
		} else {
			menB -= lossB
		}
		if slot == 0 && menB == 0 {
			menB = 1
		}
		m.mem[b+0x29+slot*4] = menB
		newB += uint16(menB)
	}
	strategyPutWord(m.mem, a+4, newA)
	strategyPutWord(m.mem, b+4, newB)
	moraleA, moraleB := byte(0), byte(0)
	if oldA != 0 && m.mem[a+6] >= 100 {
		moraleA = byte(uint32(m.mem[a+6]) * uint32(newA) / uint32(oldA))
	}
	if oldB != 0 && m.mem[b+6] >= 100 {
		moraleB = byte(uint32(100) * uint32(newB) / uint32(oldB))
	}
	m.mem[a+6], m.mem[b+6] = moraleA, moraleB
}

func (m *outcomeModel) recalculate(army uint16) {
	m.enter(0x6FD2)
	a := outcomeWorld + int(army)
	sum := uint16(0)
	mixed := false
	for i := 0; i < 6; i++ {
		sum += uint16(m.mem[a+0x29+i*4])
		if m.mem[a+0x2A+i*4] != 1 {
			mixed = true
		}
	}
	strategyPutWord(m.mem, a+4, sum)
	m.mem[a+0x1E] = 2
	if mixed {
		m.mem[a+0x1E] = 3
	}
	// IDA17012..1701A: preserve x, SHL twice, add x; AL truncates to byte.
	m.mem[a+9] = m.mem[outcomeWorld+int(m.mem[a+1])*64+0x3E] * 5
	m.mem[a+0x0B] = 1
}

func (m *outcomeModel) spent(army uint16, retreat byte) bool {
	m.enter(0x474A)
	m.recalculate(army)
	a := outcomeWorld + int(army)
	if m.mem[a+6] == 0 || m.mem[a+0x29] == 0 {
		return true
	}
	stage := byte(8)
	if retreat != 0 {
		node := strategyWord(m.mem, a+0x0E)
		if node >= 0x600 || m.mem[outcomeWorld+0x840+int(node)*4+1] != m.mem[a+1] {
			result := m.routeModel.retreat((int(army) - 0x2240) / 64)
			if result.carry {
				return true
			}
			target := result.bx >> 2
			strategyPutWord(m.mem, a+0x14, target)
			m.mem[a+0x20] = byte(target >> 3)
			m.mem[a] |= 2
			if strategyWord(m.mem, a+4) <= 300 || m.mem[a+0x20] == m.mem[outcomeWorld+int(m.mem[a+1])*64+3] {
				stage = 10
			}
		}
	}
	m.mem[a+0x23] = stage
	return false
}

func (m *outcomeModel) automatic(first, second uint16, mode byte) uint16 {
	m.enter(0x5130)
	firstMode := mode
	if firstMode == 0 {
		firstMode = 3
	}
	a := m.ability(first, mode, m.coefficient(first, firstMode)) + 8
	b := m.ability(second, mode, m.coefficient(second, mode)) + 8
	if a == 0 || b == 0 {
		panic("outcome auto fixture has zero wrapped denominator")
	}
	winner, loser := first, second
	lo, hi := b, a
	result := byte(0)
	if a < b {
		winner, loser, lo, hi, result = second, first, a, b, 1
	}
	ratio := uint32(hi) * 8 / uint32(lo)
	if ratio > 65535 {
		panic("outcome auto DIV overflows AX")
	}
	if ratio >= 100 {
		ratio = 100
	}
	m.casualties(winner, loser, byte(ratio))
	broken := byte(0)
	if m.spent(first, result) {
		broken |= 1
	}
	if m.spent(second, result^1) {
		broken |= 2
	}
	return uint16(broken)<<8 | uint16(result)
}

func (m *outcomeModel) sound(short bool) {
	if short {
		m.enter(0x0CDE)
		m.sounds = append(m.sounds, 0x0101)
	} else {
		m.enter(0x0CE7)
		m.sounds = append(m.sounds, 0x0202)
	}
	m.enter(0xEB11)
}
func (m *outcomeModel) officerCity(index int) {
	m.enter(0x4D63)
	a := outcomeWorld + 0x840 + index*32
	g := m.mem[a+0x19]
	if g == 255 {
		return
	}
	m.mem[a+0x19] = 255
	general := outcomeWorld + 0x4240 + int(g)*32
	m.mem[general+0x17] = 0
	m.sound(false)
	m.warning(0x93, 0x44)
	m.warning(uint16(m.mem[general+0x1E])<<8|uint16(m.mem[general+1]), 0x1A6)
}
func (m *outcomeModel) officerFaction(owner byte) {
	m.enter(0x5074)
	f := outcomeWorld + int(owner)*64
	g := m.mem[f+0x2A]
	if g == 255 {
		return
	}
	m.mem[f+0x2A] = 255
	general := outcomeWorld + 0x4240 + int(g)*32
	m.mem[general+0x17] = 0
	m.warning(0x93, 0x45)
	m.warning(uint16(m.mem[general+0x1E])<<8|uint16(m.mem[general+1]), 0x1A7)
}
func (m *outcomeModel) clearArmy(index int) {
	m.enter(0x7028)
	a := outcomeWorld + 0x2240 + index*64
	flags := m.mem[a]
	m.mem[a] = 0
	if flags&8 == 0 {
		m.decrementOccupancy(a)
	}
}
func (m *outcomeModel) detach(index int) {
	m.enter(0x50B4)
	g := outcomeWorld + 0x4240 + index*32
	if m.mem[g+0x17] != 0 {
		m.miniDirty(outcomeWorld + 0x2240 + index*64)
		m.clearArmy(index)
		m.mem[g+0x17] = 0
	}
	m.mem[g+0x1C] = 255
}
func (m *outcomeModel) release(index int) {
	m.enter(0x50D7)
	g := outcomeWorld + 0x4240 + index*32
	owner := m.mem[g+0x1D]
	m.mem[g+0x17], m.mem[g+0x1D] = 0, 255
	if m.mem[outcomeWorld+int(owner)*64] < 0x80 {
		owner = 255
	}
	m.mem[g+0x1C] = owner
	m.enter(0x2AD2)
	if owner != 255 {
		m.mem[outcomeWorld+int(owner)*64+0x18]++
	}
	if owner == m.c.player {
		m.sound(true)
		m.warning(0x93, 0x25)
		m.warning(uint16(m.mem[g+0x1E])<<8|uint16(m.mem[g+1]), 0x199)
	}
}
func (m *outcomeModel) history(a, b byte) {
	m.enter(0x4236)
	for i := 0; i < 192; i++ {
		record := outcomeWorld + 0x840 + i*32
		owner, old := m.mem[record+1], m.mem[record+0x1A]
		if (owner == a || owner == b) && (old == a || old == b) {
			m.mem[record+0x1A] = owner
		}
	}
}
func (m *outcomeModel) capital(owner byte) int {
	m.enter(0x6A3D)
	kind, production, chosen, safe := byte(255), uint16(0), -1, false
	for i := 0; i < 192; i++ {
		a := outcomeWorld + 0x840 + i*32
		k, p := m.mem[a+0x16]&15, strategyWord(m.mem, a+14)
		if m.mem[a+1] != owner || kind < k || production > p {
			continue
		}
		if !safe {
			kind, production, chosen = k, p, i
		}
		if m.mem[a]&31 == 0 {
			safe = true
			kind, production, chosen = k, p, i
		}
	}
	return chosen
}
func (m *outcomeModel) relocate(index int, owner byte) bool {
	m.enter(0x4DF0)
	f := outcomeWorld + int(owner)*64
	old := m.mem[f+3]
	if old != byte(index) {
		return false
	}
	chosen := m.capital(owner)
	if chosen < 0 {
		m.mem[f+3] = 255
		m.mem[f] &= 0x7F
		return true
	}
	m.mem[f+3] = byte(chosen)
	m.enter(0x4502)
	//14502's literal operands are retained, including its new-node comparison
	// before writing the old-node word. No field-name-based correction is made.
	for slot := 0; slot < 127; slot++ {
		a := outcomeWorld + 0x2240 + slot*64
		if m.mem[a+1] == owner && m.mem[a] >= 0x80 && m.mem[a+0x20] == old {
			m.mem[a+0x20] = byte(chosen)
			if strategyWord(m.mem, a+0x14) == uint16(chosen*8) {
				strategyPutWord(m.mem, a+0x14, uint16(old)*8)
				m.mem[a] |= 2
			}
		}
	}
	if owner == m.c.player {
		m.sound(false)
		m.warning(0x93, 0x1E)
		m.enter(0x5E60)
	}
	return false
}
func (m *outcomeModel) redirect(index int) {
	m.enter(0x4DA4)
	frame := mainPhysical(m.c.entrySS, m.c.bp)
	count := int(m.mem[frame+0xFE])
	if count == 0 {
		panic("direct14DA4 requires nonempty original list")
	}
	first := strategyWord(m.mem, frame)
	result := m.routeModel.retreat((int(first) - 0x2240) / 64)
	newOwner := m.mem[outcomeWorld+0x840+index*32+1]
	for i := 0; i < count; i++ {
		army := strategyWord(m.mem, frame+i*2)
		m.redirects = append(m.redirects, registers{SI: army})
		if result.carry {
			m.routeModel.collapse((int(army)-0x2240)/64, newOwner)
		} else {
			a := outcomeWorld + int(army)
			target := result.bx >> 2
			m.mem[a+0x20] = byte(target >> 3)
			strategyPutWord(m.mem, a+0x14, target)
			m.mem[a+0x0B] = 1
			m.mem[a] |= 2
		}
	}
}
func (m *outcomeModel) eliminate(index int, owner byte) {
	m.enter(0x4FCE)
	m.clearOrder = append(m.clearOrder, 0x4FCE)
	city := outcomeWorld + 0x840 + index*32
	newOwner := m.mem[city+1]
	code := int(m.cs) * 16
	f := outcomeWorld + int(owner)*64
	m.mem[code+0x2919] = newOwner
	m.mem[f] &= 0x7F
	if uint16(owner)*64 == strategyWord(m.mem, code+0x0CFD) {
		m.escape = true
		m.enter(0x1CB1)
		m.enter(0x0A1C)
		m.calls[0xEBDC] += 272
		return
	}
	m.mem[code+0x0D2A]--
	m.officerFaction(owner)
	m.history(owner, newOwner)
	ruler := m.mem[f+1]
	for g := 0; g < 127; g++ {
		a := outcomeWorld + 0x4240 + g*32
		if m.mem[a] < 0x80 || m.mem[a+0x1C] != owner {
			continue
		}
		if m.mem[a+0x1D] != 255 {
			m.release(g)
		} else if byte(g) != ruler && m.mem[a+0x17] != 0 {
			m.detach(g)
		} else {
			m.routeModel.capture(g)
		}
	}
	m.sound(true)
	m.warning(0x93, 0x24)
	for faction := 0; faction < 22; faction++ {
		a := outcomeWorld + faction*64
		if m.mem[a] >= 0x80 && m.mem[a+0x19] == owner {
			m.mem[a+0x19] = 255
			m.enter(0x3669)
			//The actual1504E loop clears AX, so13669 receives AH=0, not CH.
			left, right := outcomeWorld+0x600+faction*24, outcomeWorld+0x600+faction
			value := m.mem[left]
			if m.mem[right] < value {
				value = m.mem[right]
			}
			value |= 0x80
			m.mem[left], m.mem[right] = value, value
		}
	}
}
func (m *outcomeModel) neighbour(index, other int, owner, bit byte) {
	m.enter(0x890A)
	a, b := outcomeWorld+0x840+index*32, outcomeWorld+0x840+other*32
	back := byte(0)
	for i := 0; i < 4; i++ {
		if m.mem[b+0x1C+i] == byte(index) {
			back = 1 << i
			break
		}
	}
	if back == 0 {
		return
	}
	if owner != m.mem[b+1] {
		if m.mem[b]&back == 0 {
			m.mem[b+0x1B]++
			m.mem[b] |= back
		}
		if m.mem[a]&bit == 0 {
			m.mem[a+0x1B]++
			m.mem[a] |= bit
		}
	} else {
		if m.mem[b]&back != 0 {
			m.mem[b+0x1B]--
			m.mem[b] &^= back
		}
		if m.mem[a]&bit != 0 {
			m.mem[a+0x1B]--
			m.mem[a] &^= bit
		}
	}
}
func (m *outcomeModel) neighbours(index int) {
	m.enter(0x88CC)
	a := outcomeWorld + 0x840 + index*32
	bit := byte(1)
	for i := 0; i < 4; i++ {
		other := m.mem[a+0x1C+i]
		if other == 255 {
			continue
		}
		m.neighbour(index, int(other), m.mem[a+1], bit)
		bit <<= 1
	}
}
func (m *outcomeModel) colourCity(index int) {
	m.enter(0x8A1E)
	a := outcomeWorld + 0x840 + index*32
	code := int(m.cs) * 16
	segment := strategyWord(m.mem, code+0x0D44) + (strategyWord(m.mem, a+10)-2)*24
	offset := strategyWord(m.mem, a+8) + 0x300
	at := mainPhysical(segment, offset)
	owner := m.mem[a+1]
	shade := byte(1)
	if owner == 24 {
		shade = 2
	} else if owner == m.c.player {
		shade = 0
	}
	m.mem[at] = byte(uint16(byte(m.mem[at]-0xCB))/3*3+0xCB) + shade
	shade = 0
	if owner == m.c.player {
		shade = 10
	}
	deltas := [4]int{-0x181, -0x17F, 0x181, 0x17F}
	switch m.mem[a+0x16] & 15 {
	case 0:
		deltas = [4]int{-0x302, -0x2FE, 0x302, 0x2FE}
	case 3:
		deltas = [4]int{-0x180, -1, 1, 0x180}
	}
	for _, delta := range deltas {
		m.enter(0x8AD1)
		p := mainPhysical(segment, offset+uint16(delta))
		t := m.mem[p]
		if t >= 0xDE && t <= 0xF1 {
			n := t - 0xDE
			if n >= 10 {
				n -= 10
			}
			m.mem[p] = 0xDE + n + shade
		}
	}
}
func (m *outcomeModel) restoreMark() {
	m.enter(0x5CA4)
	code := int(m.cs) * 16
	x, y := strategyWord(m.mem, code+0x9892), strategyWord(m.mem, code+0x9894)
	if x == 65535 {
		return
	}
	m.enter(0x9541)
	m.marks = append(m.marks, registers{DX: (x >> 1) + 440, BX: (y >> 1) + 40})
	strategyPutWord(m.mem, code+0x9892, 65535)
}
func (m *outcomeModel) paintCity(index int) {
	m.enter(0x5CE0)
	a := outcomeWorld + 0x840 + index*32
	owner := m.mem[a+1]
	color := byte(0x83)
	if owner == 24 {
		color = 15
	} else if owner == m.c.player {
		color = 0xAC
	} else if owner == m.mem[int(m.cs)*16+0x98A7] {
		color = 0xF3
	}
	m.enter(0x5D19)
	m.marks = append(m.marks, registers{AX: uint16(color) << 8, DX: (strategyWord(m.mem, a+8) >> 1) + 438, BX: (strategyWord(m.mem, a+10) >> 1) + 38})
}
func (m *outcomeModel) saveMinimap(x, y uint16, independent bool) {
	m.enter(0x95C9)
	xx, yy := x>>1, y>>1
	if xx > 178 {
		xx = 178
	}
	xx = (xx + 438) >> 3
	yy += 38
	source, dest := yy*80+xx, yy*24+xx-55-960
	segment := strategyWord(m.mem, int(m.cs)*16+0x0D3A)
	for plane := 0; plane < 4; plane++ {
		m.savePlane(segment, source, dest+uint16(plane)*0xC00, byte(plane), false, independent)
	}
}
func (m *outcomeModel) savePlane(segment, source, dest uint16, plane byte, df, independent bool) {
	m.enter(0x963F)
	m.saves = append(m.saves, registers{AX: uint16(plane), BX: source, DI: dest, DS: 0xA0C8, ES: segment})
	for row := 0; row < 4; row++ {
		for column := uint16(0); column < 2; column++ {
			//19648 MOVSW uses CPU.write16: each byte's offset wraps atFFFF.
			at := mainPhysical(segment, dest+column)
			if independent {
				m.mem[at] = outcomeBeforePlanes[plane&3][uint16(0xC80+source+column)]
			} else {
				m.unknown[at] = true
			}
		}
		if df {
			source -= 2
			dest -= 2
		} else {
			source += 2
			dest += 2
		}
		source += 0x4E
		dest += 0x16
	}
}
func (m *outcomeModel) captureCity(index int, newOwner byte) {
	m.enter(0x4CF3)
	m.clearOrder = append(m.clearOrder, 0x4CF3)
	a := outcomeWorld + 0x840 + index*32
	old := m.mem[a+1]
	m.mem[a+1], m.mem[a+0x1A] = newOwner, old
	if old != 24 {
		m.officerCity(index)
		m.mem[outcomeWorld+int(old)*64+0x23]--
		failed := m.relocate(index, old)
		m.clearOrder = append(m.clearOrder, 0x4DF0)
		if m.mem[mainPhysical(m.c.entrySS, m.c.bp+0xFE)] != 0 {
			m.redirect(index)
			m.clearOrder = append(m.clearOrder, 0x4DA4)
		}
		if failed {
			m.eliminate(index, old)
			if m.escape {
				return
			}
		}
	}
	m.mem[outcomeWorld+int(newOwner)*64+0x23]++
	m.colourCity(index)
	m.neighbours(index)
	if m.mem[int(m.cs)*16+0x98A6]&4 != 0 {
		m.restoreMark()
		m.enter(0x9656)
		m.minimaps = append(m.minimaps, [2]uint16{strategyWord(m.mem, a+8), strategyWord(m.mem, a+10)})
		m.calls[0x96CF] += 4
		m.paintCity(index)
		m.saveMinimap(strategyWord(m.mem, a+8), strategyWord(m.mem, a+10), false)
		m.enter(0x5C58)
		code := int(m.cs) * 16
		x, y := strategyWord(m.mem, code+0x988E), strategyWord(m.mem, code+0x9890)
		if strategyWord(m.mem, code+0x9892) != x || strategyWord(m.mem, code+0x9894) != y {
			strategyPutWord(m.mem, code+0x9892, x)
			strategyPutWord(m.mem, code+0x9894, y)
			m.enter(0x96ED)
		}
	}
}

func outcomeExpected(before []byte, cs uint16, c resumeCase) *outcomeModel {
	r := &routeModel{mem: append([]byte(nil), before...), cs: cs, graph: strategyWord(before, int(cs)*16+0x9874), calls: map[uint16]int{}}
	m := &outcomeModel{routeModel: r, c: c, unknown: map[int]bool{}}
	city := func() int { return (int(c.si) - 0x840) / 32 }
	switch c.target {
	case 0x10CE7:
		m.sound(false)
	case 0x14236:
		m.history(byte(c.ax), byte(c.ax>>8))
	case 0x1474A:
		m.carry = m.spent(c.si, byte(c.cx))
		m.checkCarry = true
	case 0x14CF3:
		m.captureCity(city(), byte(c.ax))
	case 0x14D63:
		m.officerCity(city())
	case 0x14DA4:
		m.redirect(city())
	case 0x14DF0:
		m.carry = m.relocate(city(), byte(c.bx/64))
		m.checkCarry = true
	case 0x14FCE:
		m.eliminate(city(), byte(c.bx/64))
	case 0x15074:
		m.officerFaction(byte(c.bx / 64))
	case 0x150B4:
		m.detach((int(c.di) - 0x4240) / 32)
	case 0x15130:
		m.ax = m.automatic(c.si, c.di, byte(c.ax))
		m.checkAX = true
	case 0x151B3:
		m.casualties(c.si, c.di, byte(c.ax))
	case 0x15285:
		m.dx = m.coefficient(c.si, byte(c.ax))
		m.checkDX = true
	case 0x152D7:
		m.dx = m.ability(c.si, byte(c.ax), c.dx)
		m.checkDX = true
	case 0x15CA4:
		m.restoreMark()
	case 0x15CE0:
		m.paintCity(city())
	case 0x17028:
		m.clearArmy((int(c.si) - 0x2240) / 64)
	case 0x188CC:
		m.neighbours(city())
	case 0x1890A:
		m.neighbour(city(), int(c.bx)/32, byte(c.ax), byte(c.dx))
	case 0x195C9:
		m.saveMinimap(c.dx, c.bx, true)
	case 0x1963F:
		m.savePlane(0x6800, c.bx, c.di, byte(c.ax), c.df, true)
	default:
		panic("unknown independent outcome target")
	}
	return m
}

// Observed IN is the explicit platform fixture input. Native C never consumes
// this original stream. RAM/RNG expectations are derived exclusively from before.
func outcomeSoundAudit(m *outcomeModel, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, s string) {
		checks++
		if !ok {
			panic("independent outcome sound: " + s)
		}
	}
	type io struct {
		port  uint16
		value byte
	}
	var reads, writes []io
	check(len(outcomeOriginalPortReads)%3 == 0, "whole platform IN records")
	for i := 0; i < len(outcomeOriginalPortReads); i += 3 {
		port := uint16(outcomeOriginalPortReads[i]) | uint16(outcomeOriginalPortReads[i+1])<<8
		if port == 0x3DA || port == 0x61 {
			reads = append(reads, io{port, outcomeOriginalPortReads[i+2]})
		}
	}
	at, polls := 0, 0
	read := func(port uint16) byte {
		check(at < len(reads), "missing IN")
		r := reads[at]
		at++
		check(r.port == port, "platform IN order")
		if port == 0x3DA {
			polls++
		}
		return r.value
	}
	poll := func(high bool) {
		for i := 0; i < 1024; i++ {
			if (read(0x3DA)&8 != 0) == high {
				return
			}
		}
		panic("fixed outcome status poll did not terminate")
	}
	for _, duration := range m.sounds {
		cl, ch := byte(duration), byte(duration>>8)
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
			for i := 0; i < delays; i++ {
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
	soundPolls := polls
	if m.escape {
		for level := 0; level < 17; level++ {
			for _, high := range []bool{false, true, false, true} {
				poll(high)
			}
		}
	}
	check(at == len(reads), "extra status/speaker IN")
	setters := strategyObserved(observed, 0xEB11)
	check(len(setters) == len(m.sounds), "real speaker entry count")
	for i, r := range setters {
		check(r.AX == m.sounds[i], "real speaker duration AX")
	}
	check(len(strategyObserved(observed, 0xEB5E)) == soundPolls, "status helper calls before optional fade")
	var got []io
	for _, w := range device.PortLog {
		if w.Port == 0x61 {
			got = append(got, io{w.Port, w.Val})
		}
	}
	check(len(got) == len(writes), "speaker OUT count")
	for i, w := range writes {
		check(got[i] == w, fmt.Sprintf("speaker OUT%d", i))
	}
	return checks
}

func outcomeAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, s string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent outcome group=%s target=%X: %s", c.group, c.target, s))
		}
	}
	m := outcomeExpected(before, cs, c)
	code := int(cs) * 16
	check(m.escape == c.outcome.escape, "independent nonlocal decision versus declared fixture")
	worldEqual := bytes.Equal(after[outcomeWorld:outcomeWorld+outcomeWorldSize], m.mem[outcomeWorld:outcomeWorld+outcomeWorldSize])
	if !worldEqual {
		first := 0
		for first < outcomeWorldSize && after[outcomeWorld+first] == m.mem[outcomeWorld+first] {
			first++
		}
		lo, hi := first-16, first+32
		if lo < 0 {
			lo = 0
		}
		if hi > outcomeWorldSize {
			hi = outcomeWorldSize
		}
		var params []string
		for _, r := range observed {
			switch r.IP {
			case 0x5130, 0x5285, 0x52D7, 0x51B3, 0x474A, 0x6FD2, 0x487B, 0x491B, 0xECE0, 0x4CF3, 0x4DF0, 0x4DA4, 0x4FCE:
				if len(params) < 128 {
					params = append(params, fmt.Sprintf("IP=%04X AX=%04X BX=%04X CX=%04X DX=%04X SI=%04X DI=%04X BP=%04X DS=%04X ES=%04X SS=%04X FLAGS=%04X", r.IP, r.AX, r.BX, r.CX, r.DX, r.SI, r.DI, r.BP, r.DS, r.ES, r.SS, r.Flags))
				}
			}
		}
		check(false, fmt.Sprintf("complete world transition first offset=%04X expected=%02X actual=%02X before=%02X window[%04X:%04X] before=%X expected=%X actual=%X input AX=%04X BX=%04X CX=%04X DX=%04X SI=%04X DI=%04X BP=%04X scenario=%d outcome=%+v route=%+v expectedCasualties=%+v expectedAX=%04X actualAX=%04X observed=%v", first, m.mem[outcomeWorld+first], after[outcomeWorld+first], before[outcomeWorld+first], lo, hi, before[outcomeWorld+lo:outcomeWorld+hi], m.mem[outcomeWorld+lo:outcomeWorld+hi], after[outcomeWorld+lo:outcomeWorld+hi], c.ax, c.bx, c.cx, c.dx, c.si, c.di, c.bp, c.scenario, c.outcome, c.route, m.casualtyCalls, m.ax, device.CPU.R[cpu.AX], params))
	} else {
		check(true, "complete world transition")
	}
	for _, v := range []struct{ start, size int }{{int(c.route.graphSegment) * 16, 0x8C00}, {int(c.route.occupancyBase) * 16, routeGridSize}, {int(strategyWord(before, code+0x0D44)) * 16, routeGridSize}} {
		check(bytes.Equal(after[v.start:v.start+v.size], m.mem[v.start:v.start+v.size]), "graph/occupancy/decoded-map before-derived model")
	}
	var raw [258]byte
	raw[0], raw[1] = c.outcome.rngSeed^0xA5, c.outcome.rngSeed
	for i := 0; i < 256; i++ {
		raw[i+2] = byte(i*73 + i/2 + int(c.outcome.rngSeed))
	}
	check(bytes.Equal(before[code+0xECFC:code+0xEDFE], raw[:]), "raw RNG initial fixture")
	check(bytes.Equal(after[code+0xECFC:code+0xEDFE], m.mem[code+0xECFC:code+0xEDFE]), "full RNG recurrence/count")
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "date changed")
	for _, offset := range []int{0x0D2A, 0x2919, 0x98A6} {
		check(after[code+offset] == m.mem[code+offset], fmt.Sprintf("CS:%04X model", offset))
	}
	for _, offset := range []int{0x49B8, 0x49BE, 0x49D2} {
		size := 2
		if offset == 0x49D2 {
			size = 1
		}
		check(bytes.Equal(after[code+offset:code+offset+size], m.mem[code+offset:code+offset+size]), "route live patch")
	}
	if c.target == 0x15CA4 {
		check(strategyWord(after, code+0x9892) == strategyWord(m.mem, code+0x9892), "marker FFFF sentinel")
	}
	for _, ip := range []uint16{0x0CE7, 0x4236, 0x474A, 0x4CF3, 0x4D63, 0x4DA4, 0x4DF0, 0x4FCE, 0x5074, 0x50B4, 0x5130, 0x51B3, 0x5285, 0x52D7, 0x5CA4, 0x5CE0, 0x7028, 0x88CC, 0x890A, 0x95C9, 0x963F, 0xECE0, 0x487B, 0x491B, 0x4A0F, 0x6FD2, 0x6A3D, 0x4502, 0x291A, 0x2977, 0x29C3, 0x2BA8, 0x50D7, 0x8810, 0x1CB1} {
		check(len(strategyObserved(observed, ip)) == m.calls[ip], fmt.Sprintf("entry%04X count expected%d got%d", ip, m.calls[ip], len(strategyObserved(observed, ip))))
	}
	if m.checkAX {
		check(device.CPU.R[cpu.AX] == m.ax, "automatic AL/AH including tie/broken bits")
	}
	if m.checkDX {
		check(device.CPU.R[cpu.DX] == m.dx, "word/byte/XLAT battle arithmetic")
	}
	if m.checkCarry {
		check((device.CPU.Flags&cpu.CF != 0) == m.carry, "original root carry")
	}
	for i, r := range strategyObserved(observed, 0x51B3) {
		w := m.casualtyCalls[i]
		check(r.SI == w.SI && r.DI == w.DI && byte(r.AX) == byte(w.AX), "winner/loser six-slot casualty inputs")
	}
	for i, r := range strategyObserved(observed, 0x491B) {
		w := m.searches[i]
		check(r.AX == w.AX && r.BX == w.BX && r.CX == w.CX && byte(r.DX) == byte(w.DX) && r.DS == w.DS, "real retreat search inputs")
	}
	warnings := strategyObserved(observed, 0x8810)
	check(len(warnings) == len(m.warnings), "real report count")
	for i, r := range warnings {
		w := m.warnings[i]
		check(r.CX == w.CX && r.DS == w.DS, "report index/DS")
		if w.CX >= 0x199 {
			check(r.AX == w.AX, "report temperament/portrait raw bytes")
		} else {
			check(byte(r.AX) == 0x93, "report portrait93")
		}
	}
	markers := append(strategyObserved(observed, 0x9541), strategyObserved(observed, 0x5D19)...)
	check(len(markers) == len(m.marks), "real minimap mark calls")
	for i, r := range markers {
		w := m.marks[i]
		check(r.DX == w.DX && r.BX == w.BX, "minimap marker coordinates")
		if w.AX != 0 {
			check(byte(r.AX>>8) == byte(w.AX>>8), "owner marker color")
		}
	}
	saves := strategyObserved(observed, 0x963F)
	check(len(saves) == len(m.saves), "four-plane cache save calls")
	for i, r := range saves {
		w := m.saves[i]
		check(byte(r.AX) == byte(w.AX) && r.BX == w.BX && r.DI == w.DI && r.DS == w.DS && r.ES == w.ES, "VGA read-plane source/destination")
		dest := w.DI
		for row := 0; row < 4; row++ {
			for column := uint16(0); column < 2; column++ {
				at := mainPhysical(w.ES, dest+column)
				if !m.unknown[at] {
					stride := uint16(0x50)
					if c.target == 0x1963F && c.df {
						stride = 0x4C
					}
					source := w.BX + uint16(row)*stride + column
					check(after[at] == m.mem[at], fmt.Sprintf("before-plane byte saved to RAM index=%d plane=%d row=%d col=%d source=A0C8:%04X dest=%04X:%04X physical=%05X expected=%02X actual=%02X", i, byte(w.AX), row, column, source, w.ES, dest+column, at, m.mem[at], after[at]))
				}
			}
			if c.target == 0x1963F && c.df {
				dest -= 2
			} else {
				dest += 2
			}
			dest += 0x16
		}
	}
	if m.escape {
		old := c
		old.group = "escape"
		old.target = 0x11CB1
		checks += mainAudit(before, after, cs, old, observed, device)
		seen := false
		for _, r := range observed {
			if r.IP == 0x1CB1 {
				seen = true
				continue
			}
			if seen && r.CS == cs {
				check(r.IP == 0x0A1C || r.IP == 0xEBDC, "abandoned caller suffix resumed")
			}
		}
	} else {
		check(device.CPU.Seg[cpu.SS] == c.entrySS && device.CPU.R[cpu.SP] == c.entrySP+2 && device.CPU.IP == 0xF000, "normal return stack")
	}
	checks += outcomeSoundAudit(m, observed, device)
	return checks
}

func outcomeCases(base func(string, uint32, int) resumeCase) []resumeCase {
	var all []resumeCase
	city := func(index int, owner byte) outcomeCity {
		return outcomeCity{index: index, owner: owner, oldOwner: owner, flags: 0xC0, governor: 255, x: uint16(100 + index*10), y: 100, production: uint16(1000 + index*500), garrison: 80, growth: 100, prevention: 100, neighbors: [4]byte{255, 255, 255, 255}}
	}
	general := func(index int, owner byte) outcomeGeneral {
		return outcomeGeneral{index: index, flags: 0x80, owner: owner, duty: 1, captor: 255, temperament: 3, war: 8, lead: 10, score: 50, budget: 255, aptitudes: [3]byte{0x50, 0x60, 0x70}}
	}
	army := func(slot int, owner byte) outcomeArmy {
		return outcomeArmy{slot: slot, owner: owner, flags: 0xC0, morale: 100, stage: 9, men: [6]byte{50, 50, 50, 50, 50, 50}, types: [6]byte{1, 1, 1, 1, 1, 1}, current: 0, target: 16}
	}
	makeCase := func(group string, target uint32, scenario int) resumeCase {
		c := base(group, target, scenario)
		foreign := byte(c.owner)
		c.savedSS, c.savedSP, c.entrySS, c.entrySP, c.returnIP, c.bp = 0xF100, 0x7FFE, 0xF000, 0x6FFE, 0x006A, 0x5000
		c.df, c.altDS, c.mode = false, false, 1
		c.route = routeVector{graphSegment: 0x8400, occupancyBase: 0xB000, general: 0, armyFlags: 0xC0, armyOwner: c.player, generalOwner: c.player, armyLeader: 0, ruler: 0, factionFlags: 0x80, capital: 0, occupancy: 1, current: 0, target: 16, rngSeed: 3}
		c.route.nodes, c.route.links = routeGraph([]uint16{0, 8, 16}, [][2]int{{0, 1}, {1, 2}}, nil, c.player)
		o := outcomeVector{city: 0, general: 2, oldOwner: foreign, newOwner: c.player, mode: 1, ratio: 8, reportCity: 255, rngSeed: 3, liveFactions: 5, power: 256, markX: 65535, markY: 100, selectedFaction: foreign}
		o.cities = []outcomeCity{city(0, c.player), city(1, foreign), city(2, foreign)}
		o.armies = []outcomeArmy{army(0, c.player), army(1, foreign)}
		o.generals = []outcomeGeneral{general(0, c.player), general(1, foreign), general(2, foreign), general(3, foreign), general(4, foreign)}
		o.generals[2].duty, o.generals[3].duty, o.generals[4].duty, o.generals[4].captor = 2, 3, 4, c.player
		o.factions = []outcomeFaction{{index: int(c.player), flags: 0x80, ruler: 0, capital: 0, cityCount: 1, foreignOfficer: 255, ally: 255}, {index: int(foreign), flags: 0x80, ruler: 1, capital: 2, cityCount: 2, foreignOfficer: 255, ally: 255}}
		c.outcome = o
		c.events = []listEvent{{100, 120, 1}}
		c.mapEvents = nil
		c.menuChoices = nil
		return c
	}
	add := func(c resumeCase) {
		o := c.outcome
		c.si, c.di, c.ax, c.bx, c.cx, c.dx = 0x2240, 0x2280, 0xA500|uint16(o.mode), uint16(o.oldOwner)*64, uint16(c.mode), o.power
		cityRecord := uint16(0x840 + o.city*32)
		switch c.target {
		case 0x151B3:
			c.ax = 0xA500 | uint16(o.ratio)
		case 0x1474A:
			c.si = uint16(0x2240 + o.armies[0].slot*64)
		case 0x14CF3, 0x14D63, 0x14DA4, 0x14DF0, 0x14FCE, 0x15CE0, 0x188CC:
			c.si = cityRecord
			c.ax = 0xA500 | uint16(o.newOwner)
		case 0x150B4:
			c.di = uint16(0x4240 + o.general*32)
		case 0x17028:
			c.si = uint16(0x2240 + o.general*64)
		case 0x14236:
			c.ax = uint16(o.newOwner)<<8 | uint16(o.oldOwner)
		case 0x1890A:
			c.si = cityRecord
			c.bx = 32
			c.ax = uint16(o.city)<<8 | uint16(o.cities[0].owner)
			c.dx = uint16(c.mask)
		case 0x195C9:
			c.dx, c.bx = c.queryX, c.queryY
		case 0x1963F:
			c.ax = uint16(o.mode)
			c.bx = 0x1200
			c.di = 0x100 + uint16(o.mode)*0xC00
		}
		all = append(all, c)
	}
	for scenario := 0; scenario < 4; scenario++ {
		for mode := byte(0); mode < 4; mode++ {
			for _, morale := range []byte{0, 1, 7, 8, 99, 100, 200, 255} {
				for profile := 0; profile < 6; profile++ {
					for _, custom := range []bool{false, true} {
						c := makeCase("battle", 0x15285, scenario)
						o := &c.outcome
						o.mode = mode
						o.armies[0].morale = morale
						o.customCoefficients = custom
						for i := 0; i < 16; i++ {
							o.coefficients[i] = 255
						}
						for slot := 0; slot < 6; slot++ {
							o.armies[0].men[slot] = []byte{0, 1, 100, 255, byte(slot * 41), 0}[profile]
							o.armies[0].types[slot] = byte(slot%4 + 1)
						}
						if profile == 5 {
							o.armies[0].men[0] = 1
						}
						add(c)
					}
				}
			}
		}
		pairs := [][2]byte{{0, 0}, {15, 1}, {1, 15}, {15, 15}, {127, 128}, {128, 127}, {255, 255}, {255, 0}}
		for mode := byte(0); mode < 3; mode++ {
			for profile, pair := range pairs {
				for nibble, aptitude := range []byte{0, 0x10, 0xF0, 0xFF} {
					for _, seed := range []byte{0, 3, 255} {
						c := makeCase("battle", 0x152D7, scenario)
						o := &c.outcome
						o.mode, o.rngSeed = mode, seed
						o.generals[0].war, o.generals[0].lead = pair[0], pair[1]
						o.generals[0].aptitudes = [3]byte{aptitude, aptitude, aptitude}
						o.power = []uint16{0, 1, 255, 256, 0x8000, 0xFFFF}[(profile+nibble)%6]
						add(c)
					}
				}
			}
		}
		for _, ratio := range []byte{0, 1, 7, 8, 15, 63, 64, 100} {
			for profile := 0; profile < 5; profile++ {
				for moraleProfile := 0; moraleProfile < 4; moraleProfile++ {
					for _, seed := range []byte{0, 3, 255} {
						c := makeCase("battle", 0x151B3, scenario)
						o := &c.outcome
						o.ratio, o.rngSeed = ratio, seed
						o.reportCity = 0
						if profile%2 != 0 {
							o.reportCity = 255
						}
						pair := [][2]byte{{99, 99}, {100, 100}, {255, 100}, {100, 255}}[moraleProfile]
						o.armies[0].morale, o.armies[1].morale = pair[0], pair[1]
						for side := 0; side < 2; side++ {
							for slot := 0; slot < 6; slot++ {
								o.armies[side].men[slot] = []byte{0, 0, 10, 100, byte(slot*17 + 1)}[profile]
							}
							if profile == 1 {
								o.armies[side].men[0] = 1
							}
						}
						o.cities[0].garrison, o.cities[0].growth, o.cities[0].prevention = 1, 15, 255
						add(c)
					}
				}
			}
		}
		for _, mode := range []byte{0, 1} {
			for profile := 0; profile < 12; profile++ {
				for _, seed := range []byte{0, 1, 3, 7, 29, 106, 127, 255} {
					c := makeCase("battle", 0x15130, scenario)
					o := &c.outcome
					o.mode, o.rngSeed = mode, seed
					o.reportCity = byte(profile%2) * 255
					o.cities[0].garrison = 0
					if profile < 2 {
						o.generals[0].war, o.generals[0].lead = 1, 15
						o.generals[1].war, o.generals[1].lead = 1, 15
						o.generals[0].aptitudes, o.generals[1].aptitudes = [3]byte{}, [3]byte{}
					}
					if profile >= 2 {
						side := profile % 2
						o.armies[side].morale = 255
						o.generals[side].war, o.generals[side].lead = 15, 1
						o.generals[side].aptitudes = [3]byte{0xA0, 0xA0, 0xA0}
						for slot := 0; slot < 6; slot++ {
							o.armies[side].men[slot] = 100
							o.armies[1-side].men[slot] = byte(profile%3 + 1)
						}
					}
					if profile >= 8 {
						o.armies[0].morale, o.armies[1].morale = []byte{0, 99, 100, 255}[profile-8], []byte{255, 99, 0, 100}[profile-8]
					}
					if profile == 7 {
						o.armies[0].men = [6]byte{}
						o.armies[1].men = [6]byte{}
					}
					add(c)
				}
			}
		}
		for _, total := range []uint16{0, 1, 299, 300, 301, 600} {
			for _, morale := range []byte{0, 1, 99, 100} {
				for _, ownNode := range []bool{false, true} {
					for _, retreat := range []byte{0, 1} {
						c := makeCase("retreat", 0x1474A, scenario)
						c.mode = retreat
						o := &c.outcome
						o.armies[0].morale = morale
						remaining := total
						for slot := 0; slot < 6; slot++ {
							men := remaining
							if men > 100 {
								men = 100
							}
							o.armies[0].men[slot] = byte(men)
							remaining -= men
						}
						o.factions[0].capital = 2
						o.cities[1].owner, o.cities[2].owner = c.player, c.player
						if !ownNode {
							o.cities[0].owner = byte(c.owner)
						}
						if total == 1 && !ownNode {
							o.armies[0].men = [6]byte{0, 1, 0, 0, 0, 0}
						}
						add(c)
					}
				}
			}
		}
		for _, variant := range []int{0, 1, 2, 3} {
			for remaining := 0; remaining < 3; remaining++ {
				for _, officer := range []byte{255, 2} {
					for _, redirect := range []bool{false, true} {
						for _, panel := range []byte{0, 4} {
							for _, count := range []byte{0, 1, 255} {
								c := makeCase("capture", 0x14CF3, scenario)
								o := &c.outcome
								old, newOwner := byte(c.owner), c.player
								if variant == 1 {
									old, newOwner = c.player, byte(c.owner)
								} else if variant >= 2 {
									old = 24
									if variant == 3 {
										newOwner = byte(c.owner)
									}
								}
								o.oldOwner, o.newOwner, o.panel = old, newOwner, panel
								o.cities = []outcomeCity{city(0, old), city(1, newOwner), city(2, newOwner)}
								o.cities[0].governor = officer
								if old == 24 {
									o.cities[0].governor = 255
								}
								capital := byte(2)
								if remaining == 0 {
									o.cities[2].owner = old
								} else {
									capital = 0
									if remaining == 1 {
										o.cities[1].owner = old
									}
								}
								if old != 24 {
									o.factions = append(o.factions, outcomeFaction{index: int(old), flags: 0x80, ruler: 0, capital: capital, cityCount: count, foreignOfficer: 3, ally: 255})
								}
								ally := byte(255)
								if old != 24 {
									ally = old
								}
								o.factions = append(o.factions, outcomeFaction{index: int(newOwner), flags: 0x80, ruler: 5, capital: 2, cityCount: count, foreignOfficer: 255, ally: ally})
								for i := range o.armies {
									owner := old
									if owner == 24 {
										owner = newOwner
									}
									o.armies[i].owner = owner
									o.armies[i].current = 0
									o.armies[i].target = 8
								}
								for i := range o.generals {
									owner := old
									if owner == 24 {
										owner = newOwner
									}
									o.generals[i].owner = owner
								}
								o.generals[0].flags = 0xD0
								o.generals[0].temperament = 253
								o.generals[4].captor = newOwner
								if redirect && old != 24 {
									o.redirectSlots = []int{0, 1}
								}
								o.escape = old == c.player && remaining == 2
								add(c)
							}
						}
					}
				}
			}
		}
		for _, target := range []uint32{0x14DF0, 0x14DA4} {
			for _, playerOld := range []bool{false, true} {
				for profile := 0; profile < 6; profile++ {
					for _, panel := range []byte{0, 4} {
						c := makeCase("capture", target, scenario)
						o := &c.outcome
						o.panel = panel
						old, newOwner := byte(c.owner), c.player
						if playerOld {
							old, newOwner = c.player, byte(c.owner)
						}
						o.oldOwner, o.newOwner = old, newOwner
						o.cities = []outcomeCity{city(0, newOwner), city(1, old), city(2, old)}
						o.cities[1].kind = byte(profile % 5)
						o.cities[1].production = 500
						o.cities[2].production = 4000
						if profile%2 != 0 {
							o.cities[1].flags |= 1
						}
						capital := byte(0)
						if profile == 0 {
							capital = 2
						}
						if profile == 5 {
							o.cities[1].owner, o.cities[2].owner = newOwner, newOwner
							capital = 255
							if target == 0x14DF0 {
								capital = 0
							}
						}
						o.factions = append(o.factions, outcomeFaction{index: int(old), flags: 0x80, ruler: 0, capital: capital, cityCount: 2, foreignOfficer: 255, ally: 255})
						for i := range o.armies {
							o.armies[i].owner = old
							o.armies[i].target = 8
						}
						if target == 0x14DA4 {
							o.redirectSlots = []int{0, 1}
						}
						for i := range o.generals {
							o.generals[i].owner = old
						}
						add(c)
					}
				}
			}
		}
		for _, playerOld := range []bool{false, true} {
			for _, diplomat := range []byte{255, 3} {
				for _, panel := range []byte{0, 4} {
					for _, live := range []byte{0, 1, 255} {
						c := makeCase("fall", 0x14FCE, scenario)
						o := &c.outcome
						old, newOwner := byte(c.owner), c.player
						if playerOld {
							old, newOwner = c.player, byte(c.owner)
						}
						o.oldOwner, o.newOwner, o.panel, o.liveFactions = old, newOwner, panel, live
						o.cities = []outcomeCity{city(0, newOwner), city(1, newOwner)}
						o.factions = append(o.factions, outcomeFaction{index: int(old), flags: 0x80, ruler: 0, capital: 255, cityCount: 0, foreignOfficer: diplomat, ally: 255})
						for i := range o.factions {
							if o.factions[i].index == int(newOwner) {
								o.factions[i].ally = old
							}
						}
						for i := range o.generals {
							o.generals[i].owner = old
						}
						o.generals[0].flags = 0xD0
						o.generals[4].captor = newOwner
						for i := range o.armies {
							o.armies[i].owner = old
						}
						o.escape = playerOld
						add(c)
					}
				}
			}
		}
		for _, target := range []uint32{0x14D63, 0x15074} {
			for _, index := range []byte{255, 2} {
				for _, budget := range []byte{0, 255} {
					for _, owner := range []byte{0, byte(makeCase("diplomat", target, scenario).owner)} {
						group := "capture"
						if target == 0x15074 {
							group = "diplomat"
						}
						c := makeCase(group, target, scenario)
						o := &c.outcome
						o.oldOwner = owner
						o.generals[2].owner, o.generals[2].budget = owner, budget
						if target == 0x14D63 {
							o.cities[0].governor = index
							o.cities[0].owner, o.cities[0].oldOwner = owner, owner
						} else {
							o.generals[2].duty = 3
							o.factions = append(o.factions, outcomeFaction{index: int(owner), flags: 0x80, ruler: 0, capital: 0, cityCount: 1, foreignOfficer: index, ally: 255})
						}
						add(c)
					}
				}
			}
		}
		for _, target := range []uint32{0x150B4, 0x17028} {
			for _, duty := range []byte{0, 1, 2, 3, 4} {
				for _, flags := range []byte{0, 8, 0x80, 0xD0} {
					for _, occupancy := range []byte{0, 255} {
						c := makeCase("army", target, scenario)
						o := &c.outcome
						o.generals[2].duty = duty
						c.route.occupancy = occupancy
						a := army(2, byte(c.owner))
						a.flags = flags
						o.armies = append(o.armies, a)
						add(c)
					}
				}
			}
		}
		for _, swap := range []bool{false, true} {
			c := makeCase("history", 0x14236, scenario)
			o := &c.outcome
			if swap {
				o.oldOwner, o.newOwner = o.newOwner, o.oldOwner
			}
			o.cities = nil
			for i := 0; i < 192; i++ {
				owner := []byte{0, byte(c.owner), 24}[i%3]
				v := city(i, owner)
				v.x, v.y = uint16(10+i%30*10), uint16(10+i/30*10)
				v.oldOwner = []byte{byte(c.owner), 0, 24}[(i/3)%3]
				o.cities = append(o.cities, v)
			}
			add(c)
		}
		for _, target := range []uint32{0x188CC, 0x1890A} {
			for mask := byte(0); mask < 16; mask++ {
				for _, same := range []bool{false, true} {
					for holes := 0; holes < 2; holes++ {
						c := makeCase("neighbors", target, scenario)
						o := &c.outcome
						o.cities = []outcomeCity{city(0, c.player), city(1, byte(c.owner)), city(2, byte(c.owner)), city(3, byte(c.owner)), city(4, byte(c.owner))}
						o.cities[0].flags, o.cities[0].neighborCount = 0xC0|mask, 255
						o.cities[0].neighbors = [4]byte{1, 2, 3, 4}
						if holes != 0 {
							o.cities[0].neighbors[1] = 255
						}
						for i := 1; i < 5; i++ {
							if same {
								o.cities[i].owner = c.player
							}
							o.cities[i].flags, o.cities[i].neighborCount = 0xC0|mask, 0
							o.cities[i].neighbors = [4]byte{255, 0, 255, 255}
						}
						c.mask = 1 << uint(mask%4)
						add(c)
					}
				}
			}
		}
		for _, mark := range [][2]uint16{{65535, 0}, {0, 0}, {100, 100}, {383, 255}} {
			c := makeCase("minimap", 0x15CA4, scenario)
			c.outcome.markX, c.outcome.markY = mark[0], mark[1]
			add(c)
		}
		for _, owner := range []byte{0, byte(makeCase("minimap", 0x15CE0, scenario).owner), 24, 21} {
			c := makeCase("minimap", 0x15CE0, scenario)
			c.outcome.cities[0].owner = owner
			add(c)
		}
		for _, x := range []uint16{0, 355, 356, 357, 383, 65535} {
			for _, y := range []uint16{0, 1, 100, 255} {
				for _, df := range []bool{false, true} {
					c := makeCase("minimap", 0x195C9, scenario)
					c.queryX, c.queryY, c.df = x, y, df
					add(c)
				}
			}
		}
		for plane := byte(0); plane < 4; plane++ {
			for _, df := range []bool{false, true} {
				c := makeCase("minimap", 0x1963F, scenario)
				c.outcome.mode, c.df = plane, df
				add(c)
			}
		}
		add(makeCase("alert", 0x10CE7, scenario))
	}
	return all
}
