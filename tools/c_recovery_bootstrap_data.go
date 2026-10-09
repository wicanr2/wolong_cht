//go:build matching_bootstrap

package main

import (
	"bytes"
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

type bootstrapArmy struct {
	slot                  int
	flags                 byte
	x, y, offset, segment uint16
}

const bootstrapWorld = 0x70000
const bootstrapGridSize = 384 * 256

func bootstrapCityLocation(mem []byte, cs uint16, record int) (uint16, uint16) {
	x, y := strategyWord(mem, record+8), strategyWord(mem, record+10)
	segment := strategyWord(mem, int(cs)*16+0x0D44) + (y-2)*24
	return segment, x + 0x300
}

func bootstrapOffsets(kind byte) [4]int {
	switch kind & 15 {
	case 0:
		return [4]int{-0x302, -0x2FE, 0x302, 0x2FE}
	case 3:
		return [4]int{-0x180, -1, 1, 0x180}
	default:
		return [4]int{-0x181, -0x17F, 0x181, 0x17F}
	}
}

func bootstrapPrepare(mem []byte, cs uint16, c resumeCase) {
	if c.repeat < 1 {
		panic("bootstrap fixture requires repeat >=1")
	}
	code := int(cs) * 16
	occupancy, decoded := c.occupancyBase, c.decodedMapBase
	if occupancy == 0 {
		occupancy = 0xB000
	}
	if decoded == 0 {
		decoded = 0x5000
	}
	for _, v := range [][2]uint16{{0x0D44, 0x5000}, {0x0D52, 0x7000}, {0x9872, occupancy}, {0xD850, decoded}} {
		strategyPutWord(mem, code+int(v[0]), v[1])
	}
	strategyPutWord(mem, code+0x0CFD, uint16(c.player)*64)
	mem[code+0x0CFF] = c.player
	strategyPutWord(mem, code+0x0CF6, c.year)
	mem[code+0x0CF4], mem[code+0x0CF0] = c.month, c.day
	copy(mem[code+0x1964:code+0x1994], c.palette[:])
	occupancyStart := int(occupancy) * 16
	for at := occupancyStart; at < occupancyStart+bootstrapGridSize; at++ {
		mem[at] = c.count
	}
	for _, army := range c.armies {
		if army.slot < 0 || army.slot > 127 {
			panic("bootstrap army fixture slot outside 0..127")
		}
		record := bootstrapWorld + 0x2240 + army.slot*64
		mem[record] = army.flags
		strategyPutWord(mem, record+0x10, army.x)
		strategyPutWord(mem, record+0x12, army.y)
		strategyPutWord(mem, record+0x1A, army.offset)
		strategyPutWord(mem, record+0x1C, army.segment)
	}
	if c.group == "city" {
		record := bootstrapWorld + 0x840 + c.city*32
		strategyPutWord(mem, record+8, c.mapX)
		strategyPutWord(mem, record+10, c.mapY)
		mem[record+1], mem[record+0x16] = c.cityOwner, c.cityKind
		segment, offset := bootstrapCityLocation(mem, cs, record)
		mem[mainPhysical(segment, offset)] = c.centerTile
		for index, delta := range bootstrapOffsets(c.cityKind) {
			mem[mainPhysical(segment, offset+uint16(delta))] = c.neighbors[index]
		}
	}
	if c.group == "fringe" {
		segment := strategyWord(mem, code+0x0D44) + c.mapY*24
		mem[mainPhysical(segment, c.mapX)] = c.centerTile
	}
	if c.group == "score" || c.group == "prefix" || c.group == "chain" {
		record := bootstrapWorld + 0x4240 + c.general*32
		mem[record] = c.generalFlags
		copy(mem[record+0x0E:record+0x11], c.aptitudes[:])
		mem[record+0x11], mem[record+0x12] = c.war, c.lead
		mem[record+0x1F] = c.scorePoison
	}
	if c.group == "score" && c.altDS {
		strategyPutWord(mem, mainPhysical(0x2600, 0x0D52), 0x7000)
	}
	if c.group == "prefix" || c.group == "chain" {
		if c.repeat != 1 || c.returnIP != 0x006A || occupancy != 0xB000 || decoded != 0x5000 {
			panic("bootstrap original prefix requires one CALL10067 with return006A")
		}
		// Poison proves that original save instructions and the real CALL must
		// establish the frame. No original prefix output is copied to native C.
		strategyPutWord(mem, code+0x9901, 0xA55A)
		strategyPutWord(mem, code+0x9903, 0x5AA5)
		strategyPutWord(mem, mainPhysical(c.savedSS, c.savedSP), 0xCC33)
	}
}

type bootstrapModel struct {
	mem     []byte
	cs      uint16
	touched map[int]bool
}

func (m *bootstrapModel) store(at int, value byte) {
	at &= 0xFFFFF
	m.mem[at] = value
	m.touched[at] = true
}

func (m *bootstrapModel) fringe(segment, offset uint16, shade byte) {
	at := mainPhysical(segment, offset)
	value := int(m.mem[at]) - 0xDE
	if value < 0 || value >= 20 {
		return
	}
	if value >= 10 {
		value -= 10
	}
	m.store(at, byte(0xDE+value+int(shade)))
}

func (m *bootstrapModel) city(index int) {
	record := bootstrapWorld + 0x840 + index*32
	segment, offset := bootstrapCityLocation(m.mem, m.cs, record)
	at := mainPhysical(segment, offset)
	owner, player := m.mem[record+1], m.mem[int(m.cs)*16+0x0CFF]
	shade := byte(1)
	if owner == 24 {
		shade = 2
	} else if owner == player {
		shade = 0
	}
	base := byte(((uint16(byte(m.mem[at]-0xCB)) / 3) * 3) + 0xCB)
	m.store(at, base+shade)
	shade = 0
	if owner == player {
		shade = 10
	}
	for _, delta := range bootstrapOffsets(m.mem[record+0x16]) {
		m.fringe(segment, offset+uint16(delta), shade)
	}
}

func (m *bootstrapModel) army(slot int) {
	record := bootstrapWorld + 0x2240 + slot*64
	if m.mem[record] < 0x80 {
		return
	}
	x, y := strategyWord(m.mem, record+0x10), strategyWord(m.mem, record+0x12)
	segment := strategyWord(m.mem, int(m.cs)*16+0x9872) + y*24
	strategyPutWord(m.mem, record+0x1A, x)
	strategyPutWord(m.mem, record+0x1C, segment)
	at := mainPhysical(segment, x)
	m.store(at, m.mem[at]+1)
}

func (m *bootstrapModel) bulk() {
	for city := 0; city < 192; city++ {
		m.city(city)
	}
	for army := 0; army < 127; army++ {
		m.army(army)
	}
}

func (m *bootstrapModel) score() {
	for general := 0; general < 127; general++ {
		record := bootstrapWorld + 0x4240 + general*32
		if m.mem[record] < 0x80 {
			continue
		}
		value := byte(0)
		for offset := 0x0E; offset <= 0x10; offset++ {
			value += m.mem[record+offset] >> 4
		}
		value += m.mem[record+0x11] * 2
		value += m.mem[record+0x12] * 2
		m.mem[record+0x1F] = value
	}
}

func bootstrapAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent bootstrap audit group=%s target=%X: %s", c.group, c.target, detail))
		}
	}
	code := int(cs) * 16
	check(strategyWord(before, code+0x0D52) == 0x7000, "input world segment mapping")
	mapStart := int(strategyWord(before, code+0x0D44)) * 16
	occupancyStart := int(strategyWord(before, code+0x9872)) * 16
	m := bootstrapModel{mem: append([]byte(nil), before...), cs: cs, touched: make(map[int]bool)}
	scorePasses, bulkPasses, cityPasses, fringePasses, armyPasses := 0, 0, 0, 0, 0
	for pass := 0; pass < c.repeat; pass++ {
		switch c.group {
		case "city":
			m.city(c.city)
			cityPasses++
		case "fringe":
			m.fringe(strategyWord(before, code+0x0D44)+c.mapY*24, c.mapX, byte(c.ax>>8))
			fringePasses++
		case "army":
			m.army(c.general)
			armyPasses++
		case "bulk":
			m.bulk()
			bulkPasses++
		case "score":
			m.score()
			scorePasses++
		case "prefix", "chain":
			m.bulk()
			m.score()
			bulkPasses++
			scorePasses++
		default:
			panic("unknown bootstrap phase-one group")
		}
	}
	for _, bank := range []struct {
		name       string
		start, end int
	}{
		{"complete D44 map", mapStart, mapStart + bootstrapGridSize},
		{"complete 9872 occupancy grid", occupancyStart, occupancyStart + bootstrapGridSize},
		{"complete world/descriptor/score records", bootstrapWorld, bootstrapWorld + 0x5240},
	} {
		check(bytes.Equal(m.mem[bank.start:bank.end], after[bank.start:bank.end]), bank.name)
	}
	// Effective offsets are 16-bit and physical addresses are 20-bit. Check
	// every modeled store as well, including explicit raw boundary fixtures.
	for at := range m.touched {
		check(after[at] == m.mem[at], fmt.Sprintf("original effective store %05X", at))
	}
	for _, bank := range []struct {
		name       string
		start, end int
	}{
		{"ICONGRF segment3", 0x22000, 0x243A0},
		{"MDL", 0x30000, 0x38000},
		{"MCH", 0x38000, 0x42832},
		{"TALK", 0x45000, 0x4D586},
		{"ICONGRF segment2", 0x68000, 0x6B000},
	} {
		check(bytes.Equal(before[bank.start:bank.end], after[bank.start:bank.end]), "readonly bank "+bank.name)
	}
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "date header remains unchanged")
	cities := strategyObserved(observed, 0x8A1E)
	armies := strategyObserved(observed, 0x8AEA)
	check(len(cities) == bulkPasses*192+cityPasses, "city component count")
	check(len(armies) == bulkPasses*127+armyPasses, "army component count including inactive slots")
	check(len(strategyObserved(observed, 0x8AD1)) == (bulkPasses*192+cityPasses)*4+fringePasses, "four original fringe calls per city")
	if bulkPasses != 0 {
		for index, r := range cities {
			check(r.DS == 0x7000 && r.SI == uint16(0x840+(index%192)*32), "original city scan order")
		}
		for index, r := range armies {
			check(r.DS == 0x7000 && r.SI == uint16(0x2240+(index%127)*64), "original army scan ends before slot127")
		}
	}
	numbers := strategyObserved(observed, 0x062F)
	check(len(numbers) == scorePasses*3, "three original date numeric calls per score pass")
	date := [...]registers{
		{AX: strategyWord(before, code+0x0CF6), DX: 0, BX: 0x0903, DI: 0x02BB},
		{AX: uint16(before[code+0x0CF4]), DX: 0, BX: 0x0902, DI: 0x02C0},
		{AX: uint16(before[code+0x0CF0]), DX: 0, BX: 0x0902, DI: 0x02C4},
	}
	for index, r := range numbers {
		want := date[index%3]
		check(r.AX == want.AX && r.DX == want.DX && r.BX == want.BX && r.DI == want.DI && r.DS == cs, "date raw numeric parameters")
	}
	check(len(strategyObserved(observed, 0x533D)) == scorePasses, "score wrapper pass count")
	if c.group == "score" {
		entryDS := cs
		if c.altDS {
			entryDS = 0x2600
		}
		check(strategyWord(before, mainPhysical(entryDS, 0x0D52)) == 0x7000, "score input DS:D52")
		check(device.CPU.Seg[cpu.DS] == cs, "score returns DS=CS rather than arbitrary entry DS")
	}
	if c.group == "prefix" || c.group == "chain" {
		check(c.repeat == 1, "original prefix executes once")
		check(strategyWord(after, code+0x9901) == c.savedSP && strategyWord(after, code+0x9903) == c.savedSS, "original prefix saved SS/SP")
		check(strategyWord(after, mainPhysical(c.savedSS, c.savedSP)) == 0x006A, "real CALL10067 outer return word")
		check(device.CPU.Seg[cpu.CS] == cs && device.CPU.IP == 0x1BF6 && device.CPU.Seg[cpu.SS] == c.savedSS &&
			device.CPU.R[cpu.SP] == c.savedSP, "prefix stops before the main loop without a synthetic RET")
		check(device.CPU.Seg[cpu.DS] == cs && device.CPU.Seg[cpu.ES] == cs, "prefix DS/ES baseline")
	}
	return checks
}

func bootstrapMainCase(c resumeCase) resumeCase {
	if c.group != "chain" {
		return c
	}
	switch c.target {
	case 0x13DC9, 0x13D91:
		c.group = "trust"
	case 0x13830:
		c.group = "scene"
	case 0x13BA9:
		c.group = "reason"
	case 0x13B5A:
		c.group = "reason-loop"
	case 0x11CB1:
		c.group = "escape"
	default:
		panic("bootstrap chain main target")
	}
	return c
}

// Each side calls this on its own prefix result. Snapshot state from the
// original executor is never used to manufacture the native prefix output.
func bootstrapMainPrepare(mem []byte, cs uint16, c resumeCase) {
	mainPrepare(mem, cs, bootstrapMainCase(c))
}

func bootstrapMainAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	return mainAudit(before, after, cs, bootstrapMainCase(c), observed, device)
}
