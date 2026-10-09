//go:build matching_personnel

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const personnelWorld = 0x70000
const personnelWorldSize = 0x5240

// These builders return record indices in original scan order. They consume
// world starting at DS:0000, rather than the complete physical-memory image.
func personnelGeneral(world []byte, player byte) int {
	ruler := int(world[int(player)*64+1])
	for general := 0; general < 128; general++ {
		at := 0x4240 + general*32
		if world[at] >= 0x80 && world[at+0x1C] == player && world[at+0x17] == 0 && general != ruler {
			return general
		}
	}
	return -1
}

func personnelCity(world []byte, player byte) int {
	for city := 0; city < 192; city++ {
		if world[0x840+city*32+1] == player {
			return city
		}
	}
	return -1
}

func personnelForeign(world []byte, player byte) int {
	for faction := 0; faction < 22; faction++ {
		if world[faction*64] >= 0x80 && faction != int(player) {
			return faction
		}
	}
	return -1
}

func personnelTarget(c resumeCase) (int, byte) {
	foreign := c.operation >= 2
	if c.group == "helper" {
		foreign = c.target == 0x16C2A
	}
	if foreign {
		return personnelWorld + c.foreign*64 + 0x2A, 3
	}
	return personnelWorld + 0x840 + c.city*32 + 0x19, 2
}

func personnelPrepare(mem []byte, cs uint16, c resumeCase) {
	if c.liveDispatch {
		binary.LittleEndian.PutUint16(mem[int(cs)*16+0x6280+c.menuRow*2:], []uint16{0x6a9b, 0x6b08, 0x6b71, 0x6be3}[c.operation])
	}
	if c.group != "helper" {
		code := int(cs) * 16
		binary.LittleEndian.PutUint16(mem[code+0x0CFD:code+0x0CFF], uint16(c.player)*64)
		mem[code+0x0CFF] = c.player
		mem[code+0x98A9] = 0
		mem[code+0x98AA] = 0
		mem[code+0x98AB] = 0
	}
	if c.group == "menu" && c.operation == -1 {
		return
	}
	field, duty := personnelTarget(c)
	if c.group == "helper" {
		mem[field] = c.oldOfficer
		if c.oldOfficer != 0xFF {
			officer := personnelWorld + 0x4240 + int(c.oldOfficer)*32
			mem[officer+0x17] = duty
			mem[officer+0x1A] = c.budget
		}
		return
	}
	mem[field] = 0xFF
	officer := personnelWorld + 0x4240 + c.general*32
	appoint := c.operation == 0 || c.operation == 2
	if appoint {
		switch c.mode {
		case 1, 3, 5:
			world := mem[personnelWorld : personnelWorld+personnelWorldSize]
			if world[c.general*32+0x4240] < 0x80 || world[c.general*32+0x425C] != c.player ||
				world[c.general*32+0x4257] != 0 || int(world[int(c.player)*64+1]) == c.general {
				panic("personnel fixture requires an original eligible non-ruler")
			}
			mem[officer+0x1A] = c.budget
		case 2:
			mem[field] = byte(c.general)
			mem[officer+0x17] = duty
			mem[officer+0x1A] = c.budget
		case 0:
		default:
			panic("unsupported personnel appointment fixture mode")
		}
	} else {
		switch c.mode {
		case 1, 5:
			mem[field] = byte(c.general)
			mem[officer+0x17] = duty
			mem[officer+0x1A] = c.budget
		case 0, 4:
		default:
			panic("unsupported personnel dismissal fixture mode")
		}
	}
}

// Audit the original guest before comparing it with native C. The model is a
// small field-update contract, independent of the generated C implementation.
func personnelAudit(before, after []byte, cs uint16, c resumeCase) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent personnel audit group=%s operation=%d mode=%d: %s", c.group, c.operation, c.mode, detail))
		}
	}
	if c.group == "menu" && c.operation == -1 {
		check(bytes.Equal(before[personnelWorld:personnelWorld+personnelWorldSize],
			after[personnelWorld:personnelWorld+personnelWorldSize]), "menu cancellation changed world records")
		return checks
	}
	field, duty := personnelTarget(c)
	want := append([]byte(nil), before[personnelWorld:personnelWorld+personnelWorldSize]...)
	appoint := c.operation == 0 || c.operation == 2
	if c.group == "helper" {
		check(before[field] == c.oldOfficer, "helper input officer")
		want[field-personnelWorld] = 0xFF
		if c.oldOfficer != 0xFF {
			officer := 0x4240 + int(c.oldOfficer)*32
			want[officer+0x17] = 0
			want[officer+0x1A] = 0
			check(after[personnelWorld+officer+0x17] == 0, "helper clears duty")
			check(after[personnelWorld+officer+0x1A] == 0, "helper clears budget")
		}
	} else {
		if appoint && (c.mode == 1 || c.mode == 5) {
			check(before[field] == 0xFF, "appointment starts with an empty target")
			want[field-personnelWorld] = byte(c.general)
			want[0x4240+c.general*32+0x17] = duty
			check(after[personnelWorld+0x4240+c.general*32+0x1A] == c.budget, "appointment preserves officer budget")
		} else if !appoint && (c.mode == 1 || c.mode == 5) {
			check(before[field] == byte(c.general), "dismissal starts with the selected officer")
			want[field-personnelWorld] = 0xFF
			want[0x4240+c.general*32+0x17] = 0
			want[0x4240+c.general*32+0x1A] = 0
			check(after[personnelWorld+0x4240+c.general*32+0x17] == 0, "dismissal clears duty")
			check(after[personnelWorld+0x4240+c.general*32+0x1A] == 0, "dismissal clears budget")
		}
	}
	check(after[field] == want[field-personnelWorld], "target officer index")
	// 0x17906 recomputes the 22 display-cache bytes even for canceled lists.
	// Exclude only that documented cache, without treating its values as an
	// independently verified diplomacy model in this personnel receipt.
	excludeCache := c.group != "helper" && c.operation >= 2
	for off, expected := range want {
		if excludeCache && off < 22*64 && off%64 == 0x3F {
			continue
		}
		if after[personnelWorld+off] != expected {
			check(false, fmt.Sprintf("unexpected world write DS:%04X expected=%02X actual=%02X", off, expected, after[personnelWorld+off]))
		}
	}
	check(true, "complete world preservation outside documented writes/cache")
	_ = cs // Register/stack ABI is checked separately by the execution harness.
	return checks
}
