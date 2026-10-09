//go:build matching_catalog

package main

import "encoding/binary"

// Addresses retain the original IDA database linear address or CS offset.
// The filters below are independent references for the fixed c-catalog probe.
type catalogFamily struct {
	Name                       string
	Caller, Builder, Draw, Row uint32
	Title, Descriptor, State   uint16
	X, Y                       uint16
	Columns                    int
}

var catalogFamilies = []catalogFamily{
	{"corps", 0x1716D, 0x171A8, 0x1724D, 0x1727D, 0x70EC, 0x713D, 0x98A8, 24, 88, 5},
	{"coordinate", 0x171D3, 0x17217, 0x1724D, 0x1727D, 0x70EC, 0x713D, 0x98A8, 24, 88, 5},
	{"general", 0x175FA, 0x1763C, 0x176DC, 0x1770C, 0x7550, 0x75C8, 0x98AA, 24, 88, 6},
	{"eligible", 0x17663, 0x176A0, 0x176DC, 0x1770C, 0x7550, 0x75C8, 0x98AA, 24, 88, 6},
	{"faction", 0x178A7, 0x178E5, 0x1796C, 0x1799C, 0x77EE, 0x7875, 0x98AB, 24, 88, 6},
	{"foreign", 0x17906, 0x17944, 0x1796C, 0x1799C, 0x77EE, 0x7875, 0x98AB, 24, 88, 6},
	{"newgame", 0x17B3C, 0x17B6F, 0x17B90, 0x17BC0, 0x7AC6, 0x7B12, 0, 136, 104, 5},
}

const catalogWorld = 0x70000

func catalogWord(mem []byte, at int) uint16 {
	return binary.LittleEndian.Uint16(mem[at : at+2])
}

func catalogPrepare(mem []byte, cs uint16, c resumeCase) {
	code := int(cs) * 16
	player := int(c.player)
	putWord(mem, cs, 0x0CFD, uint16(player*64))
	mem[code+0x0CFF] = c.player
	putWord(mem, cs, 0x722B, c.queryX)
	putWord(mem, cs, 0x7233, c.queryY)
	mem[code+0x7644] = c.paramFaction
	mem[0xD0000+0x0CFF] = c.paramFaction

	// Keep the scenario's original faction leaders, capitals and general names.
	foreign := byte((player + 1) % 22)
	for i := 0; i < 22; i++ {
		record := catalogWorld + i*64
		if i != player && mem[record] >= 0x80 && mem[record+1] < 128 {
			foreign = byte(i)
			break
		}
	}
	for i := 0; i < 22; i++ {
		mem[catalogWorld+i*64+0x3F] = 0
	}
	capital := mem[catalogWorld+player*64+3]

	// Slot 125 is the final included corps; slot 126 is outside the builder.
	slots := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 125, 126}
	for _, slot := range slots {
		off := 0x2240 + slot*64
		record := catalogWorld + off
		clear(mem[record : record+64])
		mem[record] = 0xC0
		if slot%2 != 0 {
			mem[record] |= 4
		}
		mem[record+1] = c.player
		if slot%2 != 0 {
			mem[record+1] = foreign
		}
		mem[record+2] = byte(slot % 10)
		men := uint16(299 + slot%2)
		putWord(mem, 0x7000, uint16(off+4), men)
		mem[record+6] = byte(99 + slot%2)
		mem[record+7] = 0xA5
		putWord(mem, 0x7000, uint16(off+0x0E), uint16(capital)*8)
		x, y := uint16(300), uint16(100)
		if slot >= 3 {
			x += uint16(slot * 3)
			y += uint16(slot * 2)
		}
		putWord(mem, 0x7000, uint16(off+0x10), x)
		putWord(mem, 0x7000, uint16(off+0x12), y)
		mem[record+0x20] = capital
		for troop := 0; troop < 6; troop++ {
			mem[record+0x28+troop*4] = byte(troop%3 + 1)
			mem[record+0x29+troop*4] = 50
		}
		if men == 299 {
			mem[record+0x29] = 49
		}
	}

	if c.group == "builder" && c.profile == 1 {
		record := catalogWorld + 0x4240 + 127*32
		mem[record] = 0x80
		mem[record+0x17] = 0
		mem[record+0x1C] = c.player
	}
	if c.group == "person" {
		record := catalogWorld + 0x4240
		mem[record] = 0xC0
		mem[record+0x17] = byte(c.profile % 5)
		mem[record+0x1C] = c.player
		mem[record+0x1D] = 0xFF
		if c.profile%2 != 0 {
			mem[record+0x1D] = foreign
		}
	}
	if c.group == "header" && (c.family == 2 || c.family == 3) {
		// Controlled display ownership only: retain eligibility owner and job.
		owner := c.player
		if c.family == 2 && c.altDS {
			owner = c.paramFaction
		}
		changed := 0
		for slot := 0; slot < 128 && changed < 8; slot++ {
			record := catalogWorld + 0x4240 + slot*32
			if mem[record] < 0x80 || mem[record+0x1C] != owner {
				continue
			}
			mem[record+0x1D] = 0xFF
			if changed%2 != 0 {
				mem[record+0x1D] = foreign
			}
			changed++
		}
	}
	if c.group == "army" {
		record := catalogWorld + 0x2240
		men := uint16(299 + c.profile&1)
		putWord(mem, 0x7000, 0x2244, men)
		mem[record+0x29] = byte(men - 250)
		mem[record+6] = byte(99 + (c.profile>>1)&1)
		if c.profile&4 != 0 {
			putWord(mem, 0x7000, 0x224E, 0x800)
		}
		mem[record] &^= 4
		if c.profile&8 != 0 {
			mem[record] |= 4
		}
	}
	if c.group == "cache" {
		// War entries follow peace entries to verify the LOOP target clears DL each time.
		boundary := [...]byte{
			0x00, 0x93, 0x01, 0x94, 0x02, 0xE3, 0x03, 0xE4,
			0x04, 0xE5, 0x05, 0x80, 0x06, 0xFF, 0x07, 0x92,
			0x08, 0x95, 0x09, 0xE2, 0x0A, 0xE6,
		}
		row := catalogWorld + 0x600 + player*24
		for i := range boundary {
			mem[row+i] = boundary[(i+c.profile)%len(boundary)]
		}
	}
}

func catalogRecords(mem []byte, cs uint16, c resumeCase) []uint16 {
	code := int(cs) * 16
	player := mem[code+0x0CFF]
	playerOffset := catalogWord(mem, code+0x0CFD)
	family := catalogFamilies[c.family].Name
	records := make([]uint16, 0, 128)
	switch family {
	case "corps", "coordinate":
		x := catalogWord(mem, code+0x722B)
		y := catalogWord(mem, code+0x7233)
		for slot := 0; slot < 126; slot++ {
			off := 0x2240 + slot*64
			record := catalogWorld + off
			if mem[record] < 0x80 {
				continue
			}
			if family == "corps" {
				if mem[record+1] != player {
					continue
				}
			} else if catalogWord(mem, record+0x10) != x || catalogWord(mem, record+0x12) != y {
				continue
			}
			records = append(records, uint16(off))
		}
	case "general", "eligible":
		owner := mem[code+0x7644]
		ruler := 0x4240 + int(mem[catalogWorld+int(playerOffset)+1])*32
		if family == "eligible" {
			owner = player
		}
		for slot := 0; slot < 128; slot++ {
			off := 0x4240 + slot*32
			record := catalogWorld + off
			if mem[record] < 0x80 || mem[record+0x1C] != owner {
				continue
			}
			if family == "eligible" && (mem[record+0x17] != 0 || off == ruler) {
				continue
			}
			records = append(records, uint16(off))
		}
	case "faction", "foreign", "newgame":
		for slot := 0; slot < 22; slot++ {
			off := slot * 64
			if mem[catalogWorld+off] < 0x80 || family == "foreign" && uint16(off) == playerOffset {
				continue
			}
			records = append(records, uint16(off))
		}
	default:
		panic("unknown catalog family")
	}
	return records
}

func catalogExpectedCache(mem []byte, cs uint16) []byte {
	playerOffset := catalogWord(mem, int(cs)*16+0x0CFD)
	row := catalogWorld + 0x600 + int(playerOffset>>2) + int(playerOffset>>3)
	expected := make([]byte, 22)
	for i := range expected {
		dl := byte(0)
		relationship := mem[row+i]
		if relationship >= 0x80 {
			low := relationship & 0x7F
			dl = 42
			if low <= 100 {
				if low == 100 {
					low--
				}
				dl = (low/20 + 1) * 7
			}
		}
		expected[i] = dl
	}
	return expected
}
