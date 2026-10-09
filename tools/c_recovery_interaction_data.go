//go:build matching_interaction

package main

import (
	"bytes"
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

// Controlled inputs for spec/247 and re/127; no additional production rules.
// IDA Pro 9.4 linear base 10000, KI.EXE SHA-256:
// fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868.
// Probe 328f9f6bd299452d717e92828b11e92de1c17ef49a66b4bca4c55246f9c0e3bf.
// Addresses and table targets below retain original segment offsets.
type interactionArmy struct {
	slot         int
	owner, flags byte
	x, y         uint16
}

type interactionVector struct {
	worldX, worldY, cellX, cellY, mapBase, occupancyBase uint16
	tile, count, panel, index, cityFlags, cityOwner      byte
	city                                                 int
	liveIndex                                            int
	liveTarget                                           uint16
	armies                                               []interactionArmy
	armyChoices                                          []int
}

const interactionWorld = 0x70000
const interactionGridSize = 384 * 256

func interactionPrepare(mem []byte, cs uint16, c resumeCase) {
	code, v := int(cs)*16, c.interaction
	if v.mapBase == 0 || v.occupancyBase == 0 {
		panic("interaction fixture requires explicit map and occupancy segments")
	}
	for _, word := range [][2]uint16{
		{0x0D44, v.mapBase}, {0x0D52, 0x7000}, {0x9872, v.occupancyBase},
		{0x989A, v.worldX}, {0x989C, v.worldY}, {0x9896, v.cellX}, {0x9898, v.cellY},
		{0x0CFD, uint16(c.player) * 64},
	} {
		strategyPutWord(mem, code+int(word[0]), word[1])
	}
	mem[code+0x0CFF], mem[code+0x98A6] = c.player, v.panel
	mem[code+0x98A8], mem[code+0x98A9], mem[code+0x98AA], mem[code+0x98AB] = 0, 0, 0, 0
	copy(mem[code+0x1964:code+0x1994], c.palette[:])
	entry := code
	if c.altDS {
		if c.group != "right" && c.group != "clear" {
			panic("interaction alternate DS fixture is limited to right/clear")
		}
		entry = 0x26000
		copy(mem[entry+0x5A06:entry+0x5A3A], mem[code+0x5A06:code+0x5A3A])
		mem[entry+0x98A6] = v.panel
	}
	if v.liveIndex >= 0 {
		if v.liveIndex >= 26 {
			panic("interaction live table index outside original 26 entries")
		}
		strategyPutWord(mem, entry+0x5A06+v.liveIndex*2, v.liveTarget)
	}
	if c.group != "point" {
		return
	}
	if v.worldX > 383 || v.worldY > 255 || v.city < -1 || v.city >= 192 {
		panic("interaction point fixture outside documented legal coordinate/record range")
	}
	// The eight original locals start after PUSH DS/ES/BP and SUB SP,8.
	// Poison is an input only. Their retired final contents are not audited.
	for delta := uint16(0); delta < 8; delta++ {
		mem[mainPhysical(c.entrySS, c.entrySP-14+delta)] = 0xA5
	}
	mapRow := v.mapBase + v.worldY*24
	mem[mainPhysical(mapRow, v.worldX)] = v.tile
	occupancy := int(v.occupancyBase) * 16
	clear(mem[occupancy : occupancy+interactionGridSize])
	mem[mainPhysical(v.occupancyBase+v.worldY*24, v.worldX)] = v.count
	if v.city >= 0 {
		record := interactionWorld + 0x840 + v.city*32
		mem[record], mem[record+1] = v.cityFlags, v.cityOwner
		strategyPutWord(mem, record+8, v.worldX)
		strategyPutWord(mem, record+10, v.worldY)
		// Retain the scenario's name, portrait, art selector and statistics.
	}
	for slot := 0; slot < 128; slot++ {
		mem[interactionWorld+0x2240+slot*64] = 0
	}
	for _, army := range v.armies {
		if army.slot < 0 || army.slot >= 128 || army.owner >= 22 || c.general < 0 || c.general >= 128 {
			panic("interaction army fixture has invalid record/owner/general")
		}
		if army.owner == c.player || mem[interactionWorld+int(army.owner)*64] < 0x80 {
			panic("interaction point fixture requires an original living foreign owner")
		}
		record := interactionWorld + 0x2240 + army.slot*64
		mem[record], mem[record+1], mem[record+2] = army.flags, army.owner, byte(c.general)
		strategyPutWord(mem, record+4, 600)
		mem[record+6], mem[record+7], mem[record+0x1F] = 70, 0xA5, 0
		capital := mem[interactionWorld+int(army.owner)*64+3]
		strategyPutWord(mem, record+0x0E, uint16(capital)*8)
		strategyPutWord(mem, record+0x10, army.x)
		strategyPutWord(mem, record+0x12, army.y)
		strategyPutWord(mem, record+0x1A, army.x)
		strategyPutWord(mem, record+0x1C, v.occupancyBase+army.y*24)
		mem[record+0x20] = capital
		for troop := 0; troop < 6; troop++ {
			mem[record+0x29+troop*4], mem[record+0x2A+troop*4] = 10, byte(troop%4+1)
		}
	}
}

func interactionArmyRecords(before []byte, cs uint16) []uint16 {
	code := int(cs) * 16
	world := int(strategyWord(before, code+0x0D52)) * 16
	x, y := strategyWord(before, code+0x989A), strategyWord(before, code+0x989C)
	var records []uint16
	// Original 17217 scans 126 records. Slot 126 is already outside it.
	for slot := 0; slot < 126; slot++ {
		offset := uint16(0x2240 + slot*64)
		record := world + int(offset)
		if before[record] >= 0x80 && strategyWord(before, record+0x10) == x && strategyWord(before, record+0x12) == y {
			records = append(records, offset)
		}
	}
	return records
}

func interactionPopupDX(before []byte, cs uint16) uint16 {
	code := int(cs) * 16
	x, y := byte(strategyWord(before, code+0x9896)), byte(strategyWord(before, code+0x9898))
	if x > 35 {
		x = 35
	}
	if y > 22 {
		y = 22
	}
	return uint16(y)<<8 | uint16(x)
}

func interactionAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent interaction audit group=%s target=%X: %s", c.group, c.target, detail))
		}
	}
	code, v := int(cs)*16, c.interaction
	entryDS := cs
	if c.altDS {
		entryDS = 0x2600
	}
	entry := int(entryDS) * 16
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "game date changed")
	check(len(strategyObserved(observed, 0x1D8E)) == 0 && len(strategyObserved(observed, 0x1CD0)) == 0, "interaction entered clock/world update")
	check(bytes.Equal(before[interactionWorld:interactionWorld+0x5240], after[interactionWorld:interactionWorld+0x5240]), "world records changed")
	for _, pointer := range []uint16{0x0D44, 0x9872} {
		segment := strategyWord(before, code+int(pointer))
		start := int(segment) * 16
		check(bytes.Equal(before[start:start+interactionGridSize], after[start:start+interactionGridSize]), fmt.Sprintf("input grid %04X changed", pointer))
		check(strategyWord(after, code+int(pointer)) == segment, fmt.Sprintf("input grid pointer %04X changed", pointer))
	}
	check(bytes.Equal(before[code+0x5A06:code+0x5A3A], after[code+0x5A06:code+0x5A3A]), "live right-button table changed")
	if c.altDS {
		check(bytes.Equal(before[entry+0x5A06:entry+0x5A3A], after[entry+0x5A06:entry+0x5A3A]), "alternate DS live table changed")
		check(after[code+0x98A6] == before[code+0x98A6], "alternate DS action incorrectly wrote CS panel flags")
	}
	check(bytes.Equal(before[code+0x1964:code+0x1994], after[code+0x1964:code+0x1994]), "working palette changed")
	flag := before[entry+0x98A6]
	clearTarget := uint16(0)
	if c.group == "right" {
		index := int(byte(c.ax))
		check(index < 26, "right-button fixture exceeds original finite table")
		clearTarget = strategyWord(before, entry+0x5A06+index*2)
	} else if c.group == "clear" {
		clearTarget = uint16(c.target - 0x10000)
	}
	type rectangle struct{ dx, bx, cx uint16 }
	var rect *rectangle
	switch clearTarget {
	case 0:
	case 0x59D0:
	case 0x5AA2:
		flag &^= 4
		rect = &rectangle{27, 0, 0x0A0D}
	case 0x5E4C:
		flag &^= 2
		rect = &rectangle{27, 10, 0x0D0D}
	case 0x61B6:
		flag &^= 1
		rect = &rectangle{0, 0, 0x021B}
	default:
		panic("independent interaction fixture selected a target outside original closure")
	}
	check(after[entry+0x98A6] == flag, "panel flags differ from selected original action")
	if c.group == "right" || c.group == "clear" {
		calls := strategyObserved(observed, 0x895D)
		want := 0
		if rect != nil {
			want = 1
		}
		check(len(calls) == want, "clear rectangle call count")
		if want == 1 {
			r := calls[0]
			check(byte(r.AX) == 0 && r.DX == rect.dx && r.BX == rect.bx && r.CX == rect.cx && r.DS == entryDS,
				"clear rectangle AL/DX/BX/CX and original DS")
		}
		if c.group == "right" {
			for _, target := range []uint16{0x5AA2, 0x5E4C, 0x61B6} {
				want := 0
				if clearTarget == target {
					want = 1
				}
				check(len(strategyObserved(observed, target)) == want, fmt.Sprintf("live table dispatch target %04X", target))
			}
			nulls := 1 // Final MOV ES,AX falls through the original RET.
			if clearTarget == 0x59D0 {
				nulls++ // A table CALL to nullsub_1 returns before that fallthrough.
			}
			check(len(strategyObserved(observed, 0x59D0)) == nulls, "original null target and final RET fallthrough")
			check(len(strategyObserved(observed, 0x1C8D)) == 1, "right-button redraw call count")
			check(device.CPU.Seg[cpu.DS] == cs && device.CPU.Seg[cpu.ES] == cs, "right-button final DS/ES=CS")
		}
	}
	if c.group == "popup" || c.group == "point" {
		city, count := -1, byte(0)
		if c.group == "point" {
			x, y := strategyWord(before, code+0x989A), strategyWord(before, code+0x989C)
			mapBase, occupancy := strategyWord(before, code+0x0D44), strategyWord(before, code+0x9872)
			tile := before[mainPhysical(mapBase+y*24, x)]
			count = before[mainPhysical(occupancy+y*24, x)]
			if tile >= 0xCB && tile < 0xD4 {
				world := int(strategyWord(before, code+0x0D52)) * 16
				for index := 0; index < 192; index++ {
					record := world + 0x840 + index*32
					if strategyWord(before, record+8) == x && strategyWord(before, record+10) == y {
						city = index
						break // Original search never tests the record's existence flag.
					}
				}
				check(city >= 0, "legal city-tile fixture lacks an original matching record")
			}
		}
		popup := c.group == "popup" || city >= 0 && count != 0
		popups := strategyObserved(observed, 0x93E9)
		wantPopups := 0
		if popup {
			wantPopups = 1
		}
		check(len(popups) == wantPopups, "city/corps popup decision")
		if popup {
			r := popups[0]
			check(r.AX == 0x0102 && r.CX == 0x0050 && r.DX == interactionPopupDX(before, cs), "popup count/width/low-byte clamp")
		}
		if c.group == "popup" {
			check(len(c.menuChoices) == 1, "direct popup requires exactly one terminal input")
			choice := c.menuChoices[0]
			check((device.CPU.Flags&cpu.CF != 0) == (choice < 0), "direct popup carry cancellation")
			if choice >= 0 {
				check(byte(device.CPU.R[cpu.AX]) == byte(choice), "direct popup AL selection")
			}
		} else {
			showCity, showArmies := city >= 0 && count == 0, city < 0 && count != 0
			if popup {
				check(len(c.menuChoices) == 1, "point popup requires exactly one terminal input")
				choice := c.menuChoices[0]
				showCity, showArmies = choice >= 0 && choice != 1, choice == 1
			}
			cityCalls := strategyObserved(observed, 0x7E1F)
			wantCities := 0
			if showCity {
				wantCities = 1
			}
			check(len(cityCalls) == wantCities, "city branch call count")
			if showCity {
				check(cityCalls[0].SI == uint16(city*32), "city SI is original index*32")
			}
			selectors, armies := strategyObserved(observed, 0x71D3), strategyObserved(observed, 0x7F90)
			wantSelectors, selected := 0, []uint16(nil)
			if showArmies {
				records := interactionArmyRecords(before, cs)
				check(len(v.armyChoices) != 0 && v.armyChoices[len(v.armyChoices)-1] == -1, "army script needs final original selector cancellation")
				for _, choice := range v.armyChoices {
					wantSelectors++
					if choice < 0 {
						break
					}
					check(choice < len(records), "army script selected absent original row")
					selected = append(selected, records[choice])
				}
			}
			check(len(selectors) == wantSelectors, "army selector repeats after returning panel")
			for index, r := range selectors {
				check(r.AX == strategyWord(before, code+0x989A) && r.BX == strategyWord(before, code+0x989C),
					fmt.Sprintf("saved world coordinates restored before selector %d", index))
				check(r.BP == c.entrySP-14 && r.SS == c.entrySS, "point local frame contract")
			}
			check(len(armies) == len(selected), "army panel selection count")
			for index, offset := range selected {
				check(armies[index].SI == offset-0x2240, "army controller SI is record offset minus2240")
				check(before[interactionWorld+int(offset)+1] != c.player, "point model requires explicit foreign-panel path")
			}
			check(len(strategyObserved(observed, 0x7FDB)) == 0 && len(strategyObserved(observed, 0x4325)) == 0,
				"foreign-panel fixture entered player command/stage controller")
			check(device.CPU.Seg[cpu.DS] == cs && device.CPU.Seg[cpu.ES] == 0x5000, "point restores original caller DS/ES")
		}
	}
	if c.group == "fade" {
		setters := strategyObserved(observed, 0xEBDC)
		check(len(setters) == 272, "fade-in executes all17 levels and16 colors")
		type dacWrite struct {
			port  uint16
			value byte
		}
		wantDAC := mainBeforeDAC
		var wantWrites []dacWrite
		call := 0
		for level := 0; level <= 16; level++ {
			for index := 0; index < 16; index++ {
				at := code + 0x1964 + index*3
				raw0, raw1, raw2 := before[at], before[at+1], before[at+2]
				r := setters[call]
				check(byte(r.AX) == byte(index) && byte(r.BX) == byte(level) && byte(r.BX>>8) == raw0 &&
					byte(r.DX) == raw1 && byte(r.DX>>8) == raw2 && r.CX == uint16(index<<8|level) && r.SI == uint16(0x1964+index*3),
					fmt.Sprintf("fade-in original raw arguments level=%d color=%d", level, index))
				wantWrites = append(wantWrites, dacWrite{0x3C8, byte(index)})
				for channel, raw := range []byte{raw1, raw2, raw0} {
					value := mainDACValue(raw, byte(level))
					wantWrites = append(wantWrites, dacWrite{0x3C9, value})
					wantDAC[index*3+channel] = value & 63
				}
				call++
			}
		}
		var gotWrites []dacWrite
		for _, write := range device.PortLog {
			if write.Port == 0x3C8 || write.Port == 0x3C9 {
				gotWrites = append(gotWrites, dacWrite{write.Port, write.Val})
			}
		}
		check(len(gotWrites) == 1088, "fade-in complete DAC stream length")
		for index, want := range wantWrites {
			check(gotWrites[index] == want, fmt.Sprintf("fade-in DAC stream write %d", index))
		}
		check(bytes.Equal(device.DAC[:], wantDAC[:]), "fade-in full hardware DAC")
		check(bytes.Equal(device.DAC[48:], mainBeforeDAC[48:]), "fade-in leaves colors16..255 unchanged")
	}
	return checks
}

// The real list uses two left clicks to accept a row. Foreign panels close on
// right click, then the original point controller opens that same list again.
func interactionListEvents(choices []int) []listEvent {
	var events []listEvent
	for _, choice := range choices {
		if choice < 0 {
			return append(events, listEvent{72, 112, 1})
		}
		event := listEvent{72, uint16(112 + choice*16), 0}
		events = append(events, event, event, listEvent{72, 112, 1})
	}
	panic("interaction list case lacks final right-click cancellation")
}

func interactionCases(base func(string, uint32, int) resumeCase) []resumeCase {
	var cases []resumeCase
	makeCase := func(group string, target uint32, scenario int) resumeCase {
		c := base(group, target, scenario)
		c.interaction = interactionVector{
			worldX: 100, worldY: 100, cellX: 12, cellY: 8,
			mapBase: 0x5000, occupancyBase: 0xB000, city: 0,
			cityOwner: c.player, cityFlags: 0x80, tile: 0xCA, liveIndex: -1,
		}
		c.events, c.mapEvents, c.menuChoices = nil, nil, nil
		c.ax = 0xA500
		return c
	}
	point := func(c resumeCase) {
		v := c.interaction
		if v.tile >= 0xCB && v.tile < 0xD4 && v.count == 0 || len(c.menuChoices) != 0 && c.menuChoices[0] == 0 {
			c.events = []listEvent{{72, 112, 1}}
		} else if v.count != 0 && (len(c.menuChoices) == 0 || c.menuChoices[0] == 1) {
			c.events = interactionListEvents(v.armyChoices)
		}
		cases = append(cases, c)
	}
	for scenario := 0; scenario < 4; scenario++ {
		for profile := 0; profile < 3; profile++ {
			c := makeCase("fade", 0x109D0, scenario)
			for index := range c.palette {
				switch profile {
				case 0:
					c.palette[index] = byte(index % 16)
				case 1:
					c.palette[index] = 0xFF
				case 2:
					c.palette[index] = byte(index*37 + 11)
				}
			}
			cases = append(cases, c)
		}
		for _, target := range []uint32{0x159D0, 0x15AA2, 0x15E4C, 0x161B6} {
			for _, panel := range []byte{0, 1, 2, 4, 7, 255} {
				c := makeCase("clear", target, scenario)
				c.interaction.panel = panel
				cases = append(cases, c)
			}
		}
		for index := 0; index < 26; index++ {
			for _, panel := range []byte{0, 7, 255} {
				c := makeCase("right", 0x159B7, scenario)
				c.ax |= uint16(index)
				c.interaction.index, c.interaction.panel = byte(index), panel
				cases = append(cases, c)
			}
		}
		original := func(index int) uint16 {
			switch index {
			case 1, 11, 12:
				return 0x61B6
			case 2, 16:
				return 0x5E4C
			case 3, 21, 22, 23, 24, 25:
				return 0x5AA2
			default:
				return 0x59D0
			}
		}
		for _, index := range []int{0, 1, 2, 3, 25} {
			for _, target := range []uint16{0x59D0, 0x5AA2, 0x5E4C, 0x61B6} {
				if target == original(index) {
					continue
				}
				c := makeCase("right", 0x159B7, scenario)
				c.ax |= uint16(index)
				c.interaction.index, c.interaction.panel = byte(index), 0xFF
				c.interaction.liveIndex, c.interaction.liveTarget = index, target
				cases = append(cases, c)
			}
		}
		for index := 0; index < 26; index++ {
			c := makeCase("right", 0x159B7, scenario)
			c.ax |= uint16(index)
			c.altDS, c.interaction.index, c.interaction.panel = true, byte(index), 0xFF
			cases = append(cases, c)
		}
		for _, target := range []uint32{0x15AA2, 0x15E4C, 0x161B6} {
			c := makeCase("clear", target, scenario)
			c.altDS, c.interaction.panel = true, 0xFF
			cases = append(cases, c)
		}
		for _, target := range []uint16{0x5AA2, 0x5E4C, 0x61B6} {
			c := makeCase("right", 0x159B7, scenario)
			c.altDS, c.interaction.panel = true, 0xFF
			c.interaction.liveIndex, c.interaction.liveTarget = 0, target
			cases = append(cases, c)
		}
		// Every possible low byte occurs on each axis, with distinct poisoned
		// high bytes. The original clamps AL/DL before overwriting AH/DH.
		for low := 0; low < 256; low++ {
			c := makeCase("popup", 0x11F0E, scenario)
			c.interaction.cellX, c.interaction.cellY = 0xA500|uint16(low), 0x5A00|uint16(255-low)
			c.menuChoices = []int{low%3 - 1}
			cases = append(cases, c)
		}
		for _, xy := range [][2]uint16{{0, 0}, {35, 22}, {36, 23}, {255, 255}, {34, 21}} {
			for _, high := range []uint16{0, 0xA500} {
				for choice := -1; choice <= 1; choice++ {
					c := makeCase("popup", 0x11F0E, scenario)
					c.interaction.cellX, c.interaction.cellY = high|xy[0], high|xy[1]
					c.menuChoices = []int{choice}
					cases = append(cases, c)
				}
			}
		}
		army := func(slot int, flags byte, x, y uint16) interactionArmy {
			c := makeCase("point", 0x11E46, scenario)
			return interactionArmy{slot, byte(c.owner), flags, x, y}
		}
		for _, tile := range []byte{0, 0xCA, 0xD4, 0xFF} {
			for _, count := range []byte{0, 1} {
				c := makeCase("point", 0x11E46, scenario)
				c.interaction.tile, c.interaction.count = tile, count
				c.interaction.armies = []interactionArmy{army(0, 0xC0, 100, 100)}
				c.interaction.armyChoices = []int{-1}
				point(c)
			}
		}
		for tile := byte(0xCB); tile < 0xD4; tile++ {
			for _, flags := range []byte{0, 0x80} {
				for _, ownerKind := range []int{0, 1, 2} {
					c := makeCase("point", 0x11E46, scenario)
					owner := c.player
					if ownerKind == 1 {
						owner = byte(c.owner)
					} else if ownerKind == 2 {
						owner = 24
					}
					c.interaction.tile, c.interaction.cityFlags, c.interaction.cityOwner = tile, flags, owner
					point(c)
				}
			}
		}
		for _, tile := range []byte{0xCB, 0xD3} {
			for _, count := range []byte{1, 2, 255} {
				for choice := -1; choice <= 1; choice++ {
					c := makeCase("point", 0x11E46, scenario)
					c.interaction.tile, c.interaction.count, c.interaction.cityFlags = tile, count, 0
					c.interaction.armies = []interactionArmy{army(0, 0xC0, 100, 100), army(1, 0xC0, 100, 100)}
					c.interaction.armyChoices, c.menuChoices = []int{1, 0, -1}, []int{choice}
					point(c)
				}
			}
		}
		for _, count := range []byte{1, 2, 255} {
			for _, choices := range [][]int{{-1}, {0, -1}, {1, 0, -1}} {
				if count == 1 && choices[0] == 1 {
					continue
				}
				c := makeCase("point", 0x11E46, scenario)
				c.interaction.count, c.interaction.armyChoices = count, choices
				c.interaction.armies = []interactionArmy{army(0, 0xC0, 100, 100)}
				if count != 1 {
					c.interaction.armies = append(c.interaction.armies, army(1, 0xC0, 100, 100))
				}
				point(c)
			}
		}
		for _, a := range []interactionArmy{
			army(125, 0xC0, 100, 100), army(126, 0xC0, 100, 100), army(0, 0x7F, 100, 100),
			army(0, 0xC0, 101, 100), army(0, 0xC0, 100, 101),
		} {
			c := makeCase("point", 0x11E46, scenario)
			c.interaction.count, c.interaction.armies, c.interaction.armyChoices = 1, []interactionArmy{a}, []int{-1}
			if a.slot == 125 {
				c.interaction.armyChoices = []int{0, -1}
			}
			point(c)
		}
		for _, xy := range [][2]uint16{{0, 0}, {383, 255}, {300, 200}, {100, 100}} {
			c := makeCase("point", 0x11E46, scenario)
			c.interaction.worldX, c.interaction.worldY, c.interaction.tile = xy[0], xy[1], 0xCB
			point(c)
		}
		c := makeCase("point", 0x11E46, scenario)
		c.interaction.mapBase, c.interaction.tile = 0x5100, 0xCB
		point(c)
		c = makeCase("point", 0x11E46, scenario)
		c.interaction.occupancyBase, c.interaction.count = 0xC000, 2
		c.interaction.armies = []interactionArmy{army(0, 0xC0, 100, 100), army(1, 0xC0, 100, 100)}
		c.interaction.armyChoices = []int{1, 0, -1}
		point(c)
	}
	return cases
}
