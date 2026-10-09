//go:build matching_march

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const marchWorld = 0x70000
const marchWorldSize = 0x5240

var marchHandlers = [12]uint16{
	0x4370, 0x4370, 0x4370, 0x4370, 0x439D, 0x43AF,
	0x440F, 0x4466, 0x4483, 0x4499, 0x44A9, 0x44D6,
}

func marchWord(mem []byte, at int) uint16 {
	return binary.LittleEndian.Uint16(mem[at : at+2])
}

func marchPutWord(mem []byte, at int, value uint16) {
	binary.LittleEndian.PutUint16(mem[at:at+2], value)
}

func marchForeign(world []byte) int {
	for faction := 1; faction < 22; faction++ {
		if world[faction*64] >= 0x80 {
			return faction
		}
	}
	return -1
}

func marchPrepare(mem []byte, cs uint16, c resumeCase) {
	code := int(cs) * 16
	marchPutWord(mem, code+0x0D44, 0x5000)
	marchPutWord(mem, code+0x9872, 0xB000)
	marchPutWord(mem, code+0x0CFD, uint16(c.player)*64)
	mem[code+0x0CFF] = c.player
	army := marchWorld + 0x2240 + c.general*64
	faction := marchWorld + c.owner*64
	mem[army] = c.armyFlags
	mem[army+1] = byte(c.owner)
	mem[army+2] = byte(c.general)
	marchPutWord(mem, army+4, c.total)
	mem[army+6] = c.morale
	mem[army+0x23] = c.armyStage
	mem[army+0x20] = byte(c.targetCity)
	mem[faction] = c.factionFlags
	mem[faction+0x14] = 1
	mem[faction+0x1D] = 200
	mem[faction+0x16] = c.queue16
	mem[faction+0x17] = c.queue17
	for kind, reserve := range c.reserves {
		marchPutWord(mem, faction+4+kind*2, reserve)
	}
	for slot := 0; slot < 6; slot++ {
		mem[army+0x29+slot*4] = c.slotMen[slot]
		mem[army+0x2A+slot*4] = c.slotTypes[slot]
	}
	target := marchWorld + 0x840 + c.targetCity*32
	mem[target] = c.cityFlags
	mem[target+0x18] = c.threat
	current := c.currentCity
	if c.atTarget {
		current = c.targetCity
	}
	currentRecord := marchWorld + 0x840 + current*32
	x, y := marchWord(mem, currentRecord+8), marchWord(mem, currentRecord+10)
	node := uint16(current * 8)
	switch c.mismatch {
	case 1:
		x++
	case 2:
		y++
	case 3:
		node++
	}
	marchPutWord(mem, army+0x0E, node)
	marchPutWord(mem, army+0x10, x)
	marchPutWord(mem, army+0x12, y)
	segment := uint16(0xB000 + y*24)
	marchPutWord(mem, army+0x1A, x)
	marchPutWord(mem, army+0x1C, segment)
	mem[(int(segment)*16+int(x))&0xFFFFF] = 1
	mem[marchWorld+0x4240+c.general*32+0x17] = 1
	mem[code+0xECFC], mem[code+0xECFD] = c.rngSeed^0xA5, c.rngSeed
	for i := 0; i < 256; i++ {
		mem[code+0xECFE+i] = byte(i*73 + i/2 + int(c.rngSeed))
	}
	marchPutWord(mem, code+0x9882, c.cameraPixelX)
	marchPutWord(mem, code+0x9884, c.cameraPixelY)
	marchPutWord(mem, code+0x987E, c.lastMouseX)
	marchPutWord(mem, code+0x9880, c.lastMouseY)
	marchPutWord(mem, code+0x9886, c.screenX)
	marchPutWord(mem, code+0x9888, c.screenY)
	marchPutWord(mem, code+0x988E, c.cameraPixelX>>4)
	marchPutWord(mem, code+0x9890, c.cameraPixelY>>4)
	marchPutWord(mem, code+0x9896, c.screenX>>4)
	marchPutWord(mem, code+0x9898, c.screenY>>4)
	marchPutWord(mem, code+0x989A, (c.cameraPixelX>>4)+(c.screenX>>4))
	row := (c.cameraPixelY >> 4) + (c.screenY >> 4)
	if row < 2 {
		row = 0
	} else {
		row -= 2
	}
	marchPutWord(mem, code+0x989C, row)
}

// Only documented field writes are modeled. This model does not execute or
// decode guest opcodes and never consumes native-C output.
type marchModel struct {
	mem  []byte
	cs   uint16
	army int
	di   uint16
}

func (m *marchModel) faction() int {
	return marchWorld + int(m.mem[m.army+1])*64
}

func (m *marchModel) route(target uint16) bool {
	city := marchWorld + int(target)
	x, y := marchWord(m.mem, city+8), marchWord(m.mem, city+10)
	node := (target - 0x840) >> 2
	arrived := marchWord(m.mem, m.army+0x10) == x && marchWord(m.mem, m.army+0x12) == y &&
		marchWord(m.mem, m.army+0x0E) == node
	marchPutWord(m.mem, m.army+0x16, x)
	marchPutWord(m.mem, m.army+0x18, y)
	marchPutWord(m.mem, m.army+0x14, node)
	return arrived
}

func (m *marchModel) random() byte {
	at := int(m.cs)*16 + 0xECFC
	value := m.mem[at+2+int(m.mem[at+1])] + m.mem[at]
	m.mem[at] += 0x89
	m.mem[at+1] = value
	return value
}

func (m *marchModel) returnMen() {
	faction := m.faction()
	for slot := 0; slot < 6; slot++ {
		kind := m.mem[m.army+0x2A+slot*4]
		if kind == 4 {
			continue
		}
		if kind < 1 || kind > 3 {
			panic("independent march reserve model requires original slot kinds 1..4")
		}
		pool := faction + 4 + int(kind-1)*2
		amount := uint32(marchWord(m.mem, pool)) + uint32(m.mem[m.army+0x29+slot*4])
		if amount > 65500 {
			amount = 65500
		}
		marchPutWord(m.mem, pool, uint16(amount))
		m.mem[m.army+0x29+slot*4] = 0
	}
}

func (m *marchModel) refill() {
	m.returnMen()
	faction := m.faction()
	var remaining [3]uint16
	for slot := 0; slot < 6; slot++ {
		kind := m.mem[m.army+0x2A+slot*4]
		if kind != 4 {
			remaining[kind-1]++
		}
	}
	for slot := 0; slot < 6; slot++ {
		kind := m.mem[m.army+0x2A+slot*4]
		if kind == 4 {
			continue
		}
		pool := faction + 4 + int(kind-1)*2
		reserve, count := marchWord(m.mem, pool), remaining[kind-1]
		amount := reserve/count + reserve%count
		remaining[kind-1]--
		if amount > 100 {
			amount = 100
		}
		marchPutWord(m.mem, pool, reserve-amount)
		m.mem[m.army+0x29+slot*4] = byte(amount)
	}
	var total uint16
	interval := byte(1)
	for slot := 0; slot < 6; slot++ {
		total += uint16(m.mem[m.army+0x29+slot*4])
		if m.mem[m.army+0x2A+slot*4] != 1 {
			interval = 2
		}
	}
	if uint16(m.army-marchWorld+0x40) > 0x12C {
		interval++
	}
	marchPutWord(m.mem, m.army+4, total)
	m.mem[m.army+0x1E] = interval
	m.mem[m.army+9] = m.mem[faction+0x3E] * 5
	m.mem[m.army+0x0B] = 1
}

func (m *marchModel) dissolve(counterOnly bool) {
	m.mem[m.faction()+0x14]--
	if counterOnly {
		return
	}
	m.returnMen()
	m.mem[m.army] = 0
	general := (m.army - marchWorld - 0x2240) >> 1
	m.mem[marchWorld+general+0x4257] = 0
	offset, segment := marchWord(m.mem, m.army+0x1A), marchWord(m.mem, m.army+0x1C)
	cell := (int(segment)*16 + int(offset)) & 0xFFFFF
	m.mem[cell]--
}

func (m *marchModel) handler(ip, target uint16) {
	city := marchWorld + int(target)
	stage := m.army + 0x23
	switch ip {
	case 0x4370:
		m.mem[stage] = 0
		if m.route(target) && marchWord(m.mem, m.army+4) < 600 {
			player := marchWorld + int(marchWord(m.mem, int(m.cs)*16+0x0CFD))
			if byte((target-0x840)/32) == m.mem[player+3] {
				m.mem[stage] = 9
			}
		}
	case 0x439D:
		if m.route(target) && m.mem[city]&0x40 == 0 {
			m.mem[stage] = 1
		}
	case 0x43AF:
		if marchWord(m.mem, m.army+0x0E) >= 0x800 {
			m.mem[stage] = 0
		} else if marchWord(m.mem, m.army+4) <= 300 {
			m.mem[stage] = 10
		} else if m.mem[city]&0x40 != 0 {
			m.mem[stage] = 0
		} else if m.mem[city] < 0x80 || m.mem[marchWorld+int(m.di+0x18)] > 2 {
			m.mem[m.army+0x0B] = (m.random() & 7) + 1
			m.mem[stage] = 2
		} else if marchWord(m.mem, m.army+4) < 600 && byte((target-0x840)/32) == m.mem[m.faction()+3] {
			m.mem[stage] = 9
		}
	case 0x440F:
		if m.mem[city]&0x40 != 0 || m.mem[city] >= 0x80 && m.mem[city+0x18] <= 1 {
			m.mem[stage] = 1
			return
		}
		faction := m.faction()
		candidate := byte(0xFF)
		if m.mem[faction]&0x40 == 0 {
			candidate, m.mem[faction+0x17] = m.mem[faction+0x17], 0xFF
		}
		if candidate == 0xFF {
			candidate, m.mem[faction+0x16] = m.mem[faction+0x16], 0xFF
		}
		if candidate != 0xFF {
			if m.mem[m.army+0x20] != candidate {
				m.mem[m.army+0x20] = candidate
				m.mem[m.army] |= 2
			}
			m.mem[stage] = 0
			m.mem[city+0x18]--
		} else if m.mem[city] < 0x80 {
			m.mem[stage] = 11
			m.mem[city+0x18]--
		}
	case 0x4466:
		m.mem[stage] = 8
		for slot := 0; slot < 6; slot++ {
			if m.mem[m.army+0x29+slot*4] < 30 {
				m.mem[stage] = 11
				break
			}
		}
	case 0x4483:
		if m.mem[m.army+6] >= m.mem[m.faction()+0x1D] {
			m.mem[stage] = 1
		}
	case 0x4499:
		m.refill()
		m.mem[stage] = 3
	case 0x44A9, 0x44D6:
		capital := m.mem[m.faction()+3]
		if m.mem[m.army+0x20] != capital {
			m.mem[m.army+0x20] = capital
			m.mem[m.army] |= 2
		}
		if m.route(0x840 + uint16(capital)*32) {
			if ip == 0x44A9 {
				m.mem[stage] = 9
			} else {
				m.dissolve(false)
			}
		}
	default:
		panic(fmt.Sprintf("independent march handler has no contract for %04X", ip))
	}
}

func marchEntry(observed []registers, ip uint16) (registers, bool) {
	for _, r := range observed {
		if r.IP == ip {
			return r, true
		}
	}
	return registers{}, false
}

func marchAudit(before, after []byte, cs uint16, c resumeCase, observed []registers) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent march audit group=%s target=%X: %s", c.group, c.target, detail))
		}
	}
	code := int(cs) * 16
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "date changed")
	army := marchWorld + 0x2240 + c.general*64
	checkCommand := func(result []byte) {
		flags, stage, timer, ordered := before[army], before[army+0x23], before[army+0x0B], before[army+0x20]
		if c.mode != 3 && c.mode != 4 {
			choice := c.mode
			if choice == 5 {
				choice = 0
			}
			if choice == 0 {
				flags &^= 4
				stage = 0
			} else if choice == 1 {
				flags |= 4
				stage = 0
			} else if choice == 2 {
				stage = 11
			} else {
				panic("independent march command mode outside 0..5")
			}
			timer, ordered = 1, byte(c.targetCity)
		}
		check(result[army] == flags, "command flags including unchanged cancellation")
		check(result[army+0x23] == stage, "command stage")
		check(result[army+0x0B] == timer, "command timer")
		check(result[army+0x20] == ordered, "command ordered target")
	}
	if c.group == "command" {
		checkCommand(after)
		return checks
	}
	modelInput := before
	if c.group == "controller" {
		if before[army+1] != c.player {
			check(!hasMarchEntry(observed, 0x4325), "foreign controller unexpectedly ran dispatch")
			return checks
		}
		check(len(originalStageBefore) == len(before), "missing original-only pre-dispatch snapshot")
		checkCommand(originalStageBefore)
		modelInput = originalStageBefore
	}
	if c.group != "route" && c.group != "dispatch" && c.group != "handler" && c.group != "dissolve" && c.group != "controller" {
		// UI redraw may consume object animation flags. No world-preservation
		// claim is made here for cursor/picker/hit/menu/camera/minimap/redraw.
		return checks
	}
	m := marchModel{mem: append([]byte(nil), modelInput...), cs: cs, army: army}
	if c.group == "dispatch" || c.group == "controller" {
		entry, ok := marchEntry(observed, 0x4325)
		check(ok, "missing original dispatcher entry")
		check(entry.DS == 0x7000 && entry.SI == uint16(army-marchWorld), "dispatcher DS/SI input")
		m.di = entry.DI
		for index, handler := range marchHandlers {
			check(marchWord(modelInput, code+0x4358+index*2) == handler, "fixed original dispatch table")
		}
		index := int(modelInput[army+0x23])
		if index < 8 && modelInput[army+1] != modelInput[code+0x0CFF] {
			index += 4
		}
		check(index < len(marchHandlers), "fixture stage outside documented dispatch table")
		handler := marchHandlers[index]
		check(hasMarchEntry(observed, handler), fmt.Sprintf("selected handler %04X never entered", handler))
		target := 0x840 + uint16(modelInput[army+0x20])*32
		m.handler(handler, target)
		if c.group == "controller" {
			m.mem[army] |= 2
		}
	} else if c.group == "route" {
		entry, ok := marchEntry(observed, 0x4548)
		check(ok, "missing original route entry")
		check(entry.BX == uint16(0x840+c.targetCity*32), "route target pointer")
		m.route(entry.BX)
	} else if c.group == "handler" {
		entry, ok := marchEntry(observed, uint16(c.target))
		check(ok, "missing original handler entry")
		m.di = entry.DI
		m.handler(uint16(c.target), entry.AX)
	} else {
		m.dissolve(c.target == 0x14689)
	}
	for off, want := range m.mem[marchWorld : marchWorld+marchWorldSize] {
		if after[marchWorld+off] != want {
			check(false, fmt.Sprintf("world DS:%04X want=%02X got=%02X", off, want, after[marchWorld+off]))
		}
	}
	check(true, "complete state-handler world writes")
	check(bytes.Equal(m.mem[code+0xECFC:code+0xEDFE], after[code+0xECFC:code+0xEDFE]), "raw deterministic RNG state")
	offset, segment := marchWord(modelInput, army+0x1A), marchWord(modelInput, army+0x1C)
	cell := (int(segment)*16 + int(offset)) & 0xFFFFF
	check(after[cell] == m.mem[cell], "occupancy descriptor byte")
	for _, off := range []int{0x0E, 0x10, 0x12} {
		check(marchWord(after, army+off) == marchWord(modelInput, army+off), "current node/position remains unchanged")
	}
	return checks
}

func hasMarchEntry(observed []registers, ip uint16) bool {
	_, ok := marchEntry(observed, ip)
	return ok
}
