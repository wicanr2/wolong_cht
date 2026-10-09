//go:build matching_details

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const detailsWorld = 0x70000
const detailsWorldSize = 0x5240

func detailsWord(mem []byte, at int) uint16 {
	return binary.LittleEndian.Uint16(mem[at : at+2])
}

func detailsPrepare(mem []byte, cs uint16, c resumeCase) {
	city := detailsWorld + 0x840 + c.city*32
	mem[city+1] = byte(c.owner)
	binary.LittleEndian.PutUint16(mem[city+0x0E:city+0x10], c.production)
	mem[city+0x10] = c.growth
	mem[city+0x11] = c.prevention
	mem[city+0x13] = c.garrison
	mem[city+0x16] = byte(c.art<<4 | c.kind)
	capital := c.city
	if !c.capital {
		capital = (capital + 1) % 192
	}
	// Owner 24 deliberately addresses DS:0603, as the original capital check
	// does after drawing the neutral-owner string.
	mem[detailsWorld+c.owner*64+3] = byte(capital)
	army := detailsWorld + 0x2240 + c.general*64
	mem[army+1] = 0
	mem[army+2] = byte(c.general)
	binary.LittleEndian.PutUint16(mem[army+4:army+6], c.total)
	mem[army+6] = c.morale
	mem[army+7] = 0xA5
	for slot := 0; slot < 6; slot++ {
		mem[army+0x29+slot*4] = c.slotMen[slot]
		mem[army+0x2A+slot*4] = c.slotTypes[slot]
	}
	mem[int(cs)*16+0x98A6] = c.panel
}

func detailsCalls(observed []registers, ip uint16) []registers {
	var calls []registers
	for _, r := range observed {
		if r.IP == ip {
			calls = append(calls, r)
		}
	}
	return calls
}

// Only original-entry snapshots are consumed here. Native-C equality is a
// separate check and cannot satisfy this independent argument model.
func detailsAudit(before, after []byte, cs uint16, c resumeCase, observed []registers) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent details audit group=%s: %s", c.group, detail))
		}
	}
	check(bytes.Equal(before[detailsWorld:detailsWorld+detailsWorldSize],
		after[detailsWorld:detailsWorld+detailsWorldSize]), "world records changed")
	flag := before[int(cs)*16+0x98A6]
	if c.group == "own-panel" || c.group == "own-frame" {
		flag |= 2
	}
	check(after[int(cs)*16+0x98A6] == flag, "self-panel flag contract")
	numeric := detailsCalls(observed, 0x062F)
	strings := detailsCalls(observed, 0x06FD)
	checkNumbers := func(want []registers) {
		check(len(numeric) == len(want), fmt.Sprintf("numeric call count want=%d got=%d", len(want), len(numeric)))
		for i, r := range want {
			got := numeric[i]
			check(got.AX == r.AX && got.DX == r.DX && got.BX == r.BX && got.DI == r.DI,
				fmt.Sprintf("numeric %d AX/DX/BX/DI want=%04X/%04X/%04X/%04X got=%04X/%04X/%04X/%04X",
					i, r.AX, r.DX, r.BX, r.DI, got.AX, got.DX, got.BX, got.DI))
		}
	}
	checkStrings := func(want []registers) {
		check(len(strings) == len(want), fmt.Sprintf("string call count want=%d got=%d", len(want), len(strings)))
		for i, r := range want {
			got := strings[i]
			check(got.DS == r.DS && got.SI == r.SI && got.DX == r.DX && got.BX == r.BX && got.AX == r.AX,
				fmt.Sprintf("string %d DS:SI/position/attribute", i))
		}
	}
	if c.group == "city-values" || c.group == "city-window" {
		city := detailsWorld + 0x840 + c.city*32
		growth := int(before[city+0x10]) - 100
		growthDX, growthBX := uint16(0), uint16(0x0904)
		if growth < 0 {
			growthDX, growthBX = 0xFFFF, 0x0A04
		}
		checkNumbers([]registers{
			{AX: uint16(before[city+0x13]) * 10, DX: 0, BX: 0x0904, DI: 0x641A},
			{AX: detailsWord(before, city+0x0E), DX: 0, BX: 0x0906, DI: 0x6918},
			{AX: uint16(growth), DX: growthDX, BX: growthBX, DI: 0x6E1A},
			{AX: uint16(before[city+0x11]), DX: 0, BX: 0x0903, DI: 0x731B},
		})
		owner := int(before[city+1])
		lordDS := uint16(0x7000)
		lordSI := uint16(0x4242 + int(before[detailsWorld+owner*64+1])*32)
		if owner == 24 {
			lordDS, lordSI = cs, 0x7E19
		}
		kind := uint16(before[city+0x16]&0x0F) * 6
		if before[detailsWorld+owner*64+3] == byte(c.city) {
			kind = 0x1E
		}
		checkStrings([]registers{
			{DS: 0x7000, SI: uint16(0x842 + c.city*32), DX: 0x0080, BX: 0x0120, AX: 0x0F03},
			{DS: lordDS, SI: lordSI, DX: 0x00C0, BX: 0x0130, AX: 0x0F03},
			{DS: cs, SI: 0x7DF5 + kind, DX: 0x00C0, BX: 0x0120, AX: 0x0F03},
		})
	}
	if c.group == "city-image" || c.group == "city-window" {
		var loads []registers
		for _, r := range detailsCalls(observed, 0xE38C) {
			if r.DS == cs && r.DX == 0x0DFA {
				loads = append(loads, r)
			}
		}
		check(len(loads) == 1, "KYOGRF read call count")
		r := loads[0]
		art := before[detailsWorld+0x840+c.city*32+0x16] >> 4
		check(r.AX == uint16(uint32(art)*0x1200) && r.CX == 0 && r.DI == 0x1200 && r.SI == 0 && r.BX == 0x9000,
			"KYOGRF 16-bit seek and destination")
		check(detailsWord(before, int(cs)*16+0x987C) == 0x9000, "KYOGRF buffer mapping")
	}
	if c.group == "army-values" {
		army := detailsWorld + 0x2240 + c.general*64
		total := uint32(detailsWord(before, army+4)) * 10
		checkNumbers([]registers{
			{AX: uint16(total), DX: uint16(total >> 16), BX: 0x0F04, DI: 0x5543},
			{AX: uint16(before[army+6]), DX: 0, BX: 0x0F03, DI: 0x5549},
		})
		owner := int(before[army+1])
		general := int(before[army+2])
		lord := int(before[detailsWorld+owner*64+1])
		capital := int(before[detailsWorld+owner*64+3])
		checkStrings([]registers{
			{DS: 0x7000, SI: uint16(0x4242 + general*32), DX: 0x0240, BX: 0x00D0, AX: 0x0F03},
			{DS: 0x7000, SI: uint16(0x4242 + lord*32), DX: 0x0240, BX: 0x00E0, AX: 0x0F03},
			{DS: 0x7000, SI: uint16(0x842 + capital*32), DX: 0x0240, BX: 0x00F0, AX: 0x0F03},
		})
	}
	if c.group == "army-slots" || c.group == "army-raw" {
		army := detailsWorld + 0x2240 + c.general*64
		blits := detailsCalls(observed, 0xF9B0)
		check(len(blits) == 6, "six slot blits")
		check(detailsWord(before, int(cs)*16+0x0D50) == 0x229A, "slot icon segment mapping")
		wantNumbers := make([]registers, 6)
		for slot := 0; slot < 6; slot++ {
			kind := before[army+0x2A+slot*4]
			// Byte underflow is intentional for controlled type zero.
			source := uint16(0x12C0 + uint32(byte(kind-1))*192)
			position := uint16(0x5A48 + slot*0x500)
			r := blits[slot]
			check(r.DS == 0x229A && r.SI == source && r.BX == position-7 && r.AX == 0x1003,
				fmt.Sprintf("slot %d raw icon source and position", slot))
			wantNumbers[slot] = registers{AX: uint16(before[army+0x29+slot*4]) * 10, DX: 0, BX: 0x0904, DI: position}
		}
		checkNumbers(wantNumbers)
	}
	return checks
}

// The caller supplies the already identity-checked original KYOGRF bytes.
// File reads remain outside the independent model.
func detailsImageAudit(after, kyo []byte, c resumeCase) int {
	if c.group != "city-image" && c.group != "city-window" {
		return 0
	}
	start := (c.art * 0x1200) & 0xFFFF
	if start+0x1200 > len(kyo) || !bytes.Equal(after[0x90000:0x91200], kyo[start:start+0x1200]) {
		panic(fmt.Sprintf("independent details image audit art=%d offset=%04X", c.art, start))
	}
	return 1
}
