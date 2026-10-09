//go:build matching_formation

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const formationWorld = 0x70000

var formationInitialTypes = [6]byte{1, 1, 3, 3, 2, 2}

func formationWord(mem []byte, at int) uint16 {
	return binary.LittleEndian.Uint16(mem[at : at+2])
}

func formationPutWord(mem []byte, at int, value uint16) {
	binary.LittleEndian.PutUint16(mem[at:at+2], value)
}

// world starts at the original DS:0000, without a physical-memory prefix.
// The filter is the independent 0x176A0 reference, including ruler exclusion.
func formationGeneral(world []byte, player byte) int {
	if len(world) < 0x5240 || player >= 22 {
		panic("formation candidate input shape")
	}
	ruler := int(world[int(player)*64+1])
	for general := 0; general < 128; general++ {
		at := 0x4240 + general*32
		if world[at] >= 0x80 && world[at+0x1C] == player && world[at+0x17] == 0 && general != ruler {
			return general
		}
	}
	return -1
}

func formationPrepare(mem []byte, cs uint16, c resumeCase) {
	if len(mem) < 0x100000 || c.general < 0 || c.general >= 128 || c.player >= 22 {
		panic("formation fixture input shape")
	}
	code := int(cs) * 16
	formationPutWord(mem, code+0x0CFD, uint16(c.player)*64)
	mem[code+0x0CFF] = c.player
	formationPutWord(mem, code+0x9872, 0xB000)
	for at := 0xB0000; at < 0xC8000; at++ {
		mem[at] = 0xFF
	}
	pool := formationWorld + int(c.player)*64 + 4
	for kind, count := range c.reserves {
		formationPutWord(mem, pool+kind*2, count)
	}
	record := formationWorld + 0x2240 + c.general*64
	switch c.group {
	case "init":
		if c.poison {
			// The direct initializer must preserve every non-type byte, including
			// men. The pattern is confined to this controlled corps fixture.
			for off := 0; off < 64; off++ {
				mem[record+off] = byte(0xA5 ^ off)
			}
		}
	case "close":
		// The close helper does not consume slots or any world record.
	case "numbers", "icons", "frame":
		for slot := 0; slot < 6; slot++ {
			if c.slotTypes[slot] < 1 || c.slotTypes[slot] > 4 {
				panic("formation renderer requires slot types 1..4")
			}
			mem[record+0x2A+slot*4] = c.slotTypes[slot]
			mem[record+0x29+slot*4] = c.slotMen[slot]
		}
		if c.group == "numbers" {
			var total uint16
			for _, count := range c.slotMen {
				total += uint16(count)
			}
			formationPutWord(mem, record+4, total)
		}
	case "formation", "caller", "switch", "offpanel":
		general := formationWorld + 0x4240 + c.general*32
		world := mem[formationWorld : formationWorld+0x5240]
		if mem[general] < 0x80 || mem[general+0x1C] != c.player || mem[general+0x17] != 0 ||
			int(mem[formationWorld+int(c.player)*64+1]) == c.general {
			panic("formation full entry requires an original eligible non-ruler")
		}
		if c.group == "caller" {
			if formationGeneral(world, c.player) != c.general || len(c.events) < 4 {
				panic("formation caller requires first eligible general and two list clicks plus inner/outer exits")
			}
			mem[code+0x98AA] = 0
		}
		for slot := 0; slot < 6; slot++ {
			mem[record+0x29+slot*4] = 0
		}
	default:
		panic("unknown formation fixture group")
	}
}

// This model uses integer pool accounting, not guest instructions or C output.
type formationModel struct {
	types    [6]byte
	men      [6]byte
	reserves [3]uint16
}

func (m *formationModel) distribute() {
	var remaining [3]uint16
	for _, kind := range m.types {
		if kind < 1 || kind > 4 {
			panic("formation model type outside 1..4")
		}
		if kind != 4 {
			remaining[kind-1]++
		}
	}
	for slot, kind := range m.types {
		if kind == 4 {
			continue
		}
		index := kind - 1
		count := remaining[index]
		amount := m.reserves[index]/count + m.reserves[index]%count
		remaining[index]--
		if amount > 100 {
			amount = 100
		}
		m.reserves[index] -= amount
		m.men[slot] = byte(amount)
	}
}

func (m *formationModel) returnMen() {
	for slot, kind := range m.types {
		if kind == 4 {
			continue
		}
		index := kind - 1
		amount := uint32(m.reserves[index]) + uint32(m.men[slot])
		if amount > 65500 {
			amount = 65500
		}
		m.reserves[index] = uint16(amount)
		m.men[slot] = 0
	}
}

func (m formationModel) total() uint16 {
	var total uint16
	for _, count := range m.men {
		total += uint16(count)
	}
	return total
}

// Audit the original run first. Equality with the native C run is separate.
func formationAudit(before, after []byte, cs uint16, c resumeCase) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent formation audit %s general=%d: %s", c.group, c.general, detail))
		}
	}
	record := formationWorld + 0x2240 + c.general*64
	pool := formationWorld + int(c.player)*64 + 4
	if c.group == "init" {
		want := append([]byte(nil), before[formationWorld:formationWorld+0x5240]...)
		for slot, kind := range formationInitialTypes {
			want[record-formationWorld+0x2A+slot*4] = kind
		}
		check(bytes.Equal(want, after[formationWorld:formationWorld+0x5240]), "initializer changed bytes beyond six types")
		return checks
	}
	if c.group != "formation" && c.group != "caller" && c.group != "switch" && c.group != "offpanel" {
		check(bytes.Equal(before[formationWorld:formationWorld+0x5240], after[formationWorld:formationWorld+0x5240]), "renderer changed world records")
		return checks
	}
	model := formationModel{types: formationInitialTypes}
	for kind := range model.reserves {
		model.reserves[kind] = formationWord(before, pool+kind*2)
	}
	model.distribute()
	lastTotal := model.total()
	events := c.events
	if c.group == "caller" {
		events = events[2:]
	}
	finished, success, refused := false, false, false
	for _, event := range events {
		if event.button == 1 {
			model.returnMen()
			finished = true
			break
		}
		if event.button != 0 {
			continue
		}
		if event.x >= 280 && event.x < 368 && event.y >= 272 && event.y < 288 {
			if model.men[0] == 0 {
				refused = true
				continue
			}
			finished, success = true, true
			break
		}
		if event.x >= 200 && event.x < 224 && event.y >= 192 && event.y < 288 {
			slot := int(event.y-192) / 16
			model.returnMen()
			model.types[slot]++
			if model.types[slot] >= 4 {
				model.types[slot] = 1
			}
			model.distribute()
			lastTotal = model.total()
		}
	}
	check(finished, "event script never leaves inner formation")
	check(success == (c.mode == 1), "mode disagrees with event-derived success")
	check(c.mode != 2 || refused, "empty-leader case never attempted confirmation")
	for slot := 0; slot < 6; slot++ {
		check(after[record+0x2A+slot*4] == model.types[slot], fmt.Sprintf("slot %d type", slot))
		check(after[record+0x29+slot*4] == model.men[slot], fmt.Sprintf("slot %d men", slot))
		check(after[record+0x28+slot*4] == before[record+0x28+slot*4] &&
			after[record+0x2B+slot*4] == before[record+0x2B+slot*4], fmt.Sprintf("slot %d untouched bytes", slot))
	}
	for kind := range model.reserves {
		check(formationWord(after, pool+kind*2) == model.reserves[kind], fmt.Sprintf("reserve %d", kind))
	}
	check(formationWord(after, record+4) == lastTotal, "last displayed total; cancellation must not recalculate")
	general := formationWorld + 0x4240 + c.general*32
	faction := formationWorld + int(c.player)*64
	wantGeneral := append([]byte(nil), before[general:general+32]...)
	if success {
		wantGeneral[0x17] = 1
	}
	check(bytes.Equal(wantGeneral, after[general:general+32]), "general record beyond confirmed duty write")
	if !success {
		check(after[general+0x17] == before[general+0x17], "cancel changed general duty")
		check(after[record] == before[record], "cancel changed corps activation")
		check(after[faction+0x14] == before[faction+0x14], "cancel changed faction corps count")
		check(bytes.Equal(before[0xB0000:0xC8000], after[0xB0000:0xC8000]), "cancel changed occupancy")
		return checks
	}
	check(model.men[0] != 0, "successful leader has no men")
	check(after[general+0x17] == 1, "successful general duty")
	check(after[record] == 0xC0, "successful corps activation")
	check(after[record+1] == c.player, "successful corps owner")
	check(after[record+2] == byte(c.general), "successful corps general index")
	check(after[record+6] == before[faction+0x1D], "successful morale copies faction byte")
	check(after[record+8] == 4, "successful stationary direction")
	check(after[record+0x23] == 1, "successful corps order")
	wantCount := before[faction+0x14]
	if before[record] < 0x80 {
		wantCount++
	}
	check(after[faction+0x14] == wantCount, "successful faction corps count")
	capital := before[faction+3]
	city := formationWorld + 0x840 + int(capital)*32
	x, y := formationWord(before, city+8), formationWord(before, city+10)
	check(after[record+0x20] == capital, "successful ordered capital")
	for _, off := range []int{0x0E, 0x14} {
		check(formationWord(after, record+off) == uint16(capital)*8, "successful capital node")
	}
	for _, off := range []int{0x10, 0x16} {
		check(formationWord(after, record+off) == x, "successful capital X")
	}
	for _, off := range []int{0x12, 0x18} {
		check(formationWord(after, record+off) == y, "successful capital Y")
	}
	segment := uint16(0xB000 + (y&0xFF)*24)
	check(formationWord(after, record+0x1A) == x, "successful occupancy offset")
	check(formationWord(after, record+0x1C) == segment, "successful occupancy segment")
	cell := int(segment)*16 + int(x)
	check(cell >= 0xB0000 && cell < 0xC8000, "capital outside isolated occupancy bank")
	check(after[cell] == before[cell]+1, "successful occupancy increment including byte wrap")
	wantOccupancy := append([]byte(nil), before[0xB0000:0xC8000]...)
	wantOccupancy[cell-0xB0000]++
	check(bytes.Equal(wantOccupancy, after[0xB0000:0xC8000]), "successful occupancy changed beyond capital cell")
	check(formationWord(before, int(cs)*16+0x9872) == 0xB000, "occupancy bank contract")
	return checks
}
