//go:build matching_residual

package main

import (
	"fmt"
	"sort"

	"github.com/wicanr2/dosgolem/internal/dos"
	"github.com/wicanr2/dosgolem/internal/machine"
)

func residualSeedName(mem []byte, c residualCase) {
	base := 0x70000 + 0x5221
	switch c.namePat {
	case 0:
		mem[base] = 0xFF
		for i := 1; i < 13; i++ {
			mem[base+i] = 0
		}
	case 1:
		mem[base] = 5
		for i := 1; i < 13; i++ {
			mem[base+i] = byte(0xA0 + i)
		}
	case 2:
		mem[base] = 0x91
		for i := 1; i < 13; i++ {
			mem[base+i] = 0x40
		}
	case 3:
		mem[base] = 0
		for i := 1; i < 13; i++ {
			mem[base+i] = byte(i)
		}
	}
}

func csPatch(off uint16, data ...byte) residualPatch {
	return residualPatch{seg: 0, off: off, data: data}
}

func dsPatch(off uint16, data ...byte) residualPatch {
	return residualPatch{seg: 0x8000, off: off, data: data}
}

func residualCases(base func(string, uint32, int) residualCase, nameFont []byte) []residualCase {
	_ = nameFont
	var all []residualCase
	add := func(c residualCase) { all = append(all, c) }

	// diplomacy: sub_145F8, ax = A | B<<8, records at DS:F*64.
	for _, row := range []struct {
		a, b     byte
		pa, pb   byte
		note     string
	}{
		{2, 5, 5, 2, "mutual"},
		{2, 5, 5, 9, "single"},
		{2, 5, 9, 9, "none"},
		{3, 3, 3, 3, "same"},
		{0, 21, 21, 0, "ends"},
	} {
		c := base("diplomacy", 0x145F8, 0)
		c.ax = uint16(row.a) | uint16(row.b)<<8
		c.ds = 0x8000
		c.patches = []residualPatch{
			dsPatch(uint16(row.a)*64+0x19, row.pa),
			dsPatch(uint16(row.b)*64+0x19, row.pb),
		}
		add(c)
	}

	// picker-row: sub_15C14, al=row ah=color, slot head at DS:row*64.
	for _, row := range []struct {
		r     byte
		color byte
		head  byte
	}{
		{0, 0x90, 0x80}, {0, 0x93, 0x80}, {0, 0x9A, 0x80}, {0, 0x90, 0x00},
		{10, 0x90, 0x80}, {10, 0x93, 0x00},
		{11, 0x90, 0x80}, {11, 0x93, 0x00},
		{21, 0x90, 0x80}, {21, 0x93, 0x00},
		{5, 0x9A, 0xFF}, {16, 0x90, 0x7F},
	} {
		c := base("picker-row", 0x15C14, 0)
		c.ax = uint16(row.r) | uint16(row.color)<<8
		c.ds = 0x8000
		c.dx = 0x20C
		c.bx = 0x0D8
		c.patches = []residualPatch{dsPatch(uint16(row.r)*64, row.head)}
		add(c)
	}

	// picker-loop: sub_15AFC. DS records all present; 10CFF skip row, 198A7 selected.
	loopRecs := func(heads map[int]byte) []residualPatch {
		var out []residualPatch
		for r := 0; r < 22; r++ {
			h := byte(0x80)
			if v, ok := heads[r]; ok {
				h = v
			}
			out = append(out, dsPatch(uint16(r)*64, h))
		}
		return out
	}
	{
		c := base("picker-loop", 0x15AFC, 0)
		c.ds = 0x8000
		c.player = 0xFF
		c.patches = append(loopRecs(nil), csPatch(0x98A7, 7))
		c.pollFeed = true
		c.events = []residualEvent{{100, 100, 1}}
		add(c)
	}
	{
		c := base("picker-loop", 0x15AFC, 1)
		c.ds = 0x8000
		c.player = 9
		c.patches = append(loopRecs(nil), csPatch(0x98A7, 4))
		c.pollFeed = true
		c.events = []residualEvent{{530, 272, 0}, {100, 100, 1}}
		add(c)
	}
	{
		c := base("picker-loop", 0x15AFC, 0)
		c.ds = 0x8000
		c.player = 0xFF
		c.patches = append(loopRecs(nil), csPatch(0x98A7, 0))
		c.pollFeed = true
		c.events = []residualEvent{{600, 288, 0}, {100, 100, 1}}
		add(c)
	}
	{
		c := base("picker-loop", 0x15AFC, 2)
		c.ds = 0x8000
		c.player = 0xFF
		c.patches = append(loopRecs(map[int]byte{6: 0x00}), csPatch(0x98A7, 1))
		c.pollFeed = true
		c.events = []residualEvent{{530, 320, 0}, {530, 224, 0}, {100, 100, 1}}
		add(c)
	}
	{
		c := base("picker-loop", 0x15AFC, 3)
		c.ds = 0x8000
		c.player = 3
		c.patches = append(loopRecs(nil), csPatch(0x98A7, 20))
		c.pollFeed = true
		c.events = []residualEvent{{530, 272, 0}, {600, 288, 0}, {100, 100, 1}}
		add(c)
	}

	// picker-entry: sub_15AD1; 98A6 bit1 selects the 15E1E path.
	for _, row := range []struct {
		scen   int
		player byte
		a6, a7 byte
		click  residualEvent
	}{
		{0, 0xFF, 0x00, 2, residualEvent{100, 100, 1}},
		{1, 6, 0x02, 9, residualEvent{100, 100, 1}},
		{2, 0xFF, 0x00, 2, residualEvent{100, 100, 1}},
		{1, 12, 0x00, 9, residualEvent{100, 100, 1}},
	} {
		c := base("picker-entry", 0x15AD1, row.scen)
		c.ds = 0x8300
		c.player = row.player
		c.patches = append(loopRecs(nil), csPatch(0x98A6, row.a6), csPatch(0x98A7, row.a7))
		c.pollFeed = true
		c.events = []residualEvent{row.click, {100, 100, 1}}
		add(c)
	}

	// season: sub_19377 via CF0/CF3/CF4, plus direct sub_193D7.
	for _, cf3 := range []byte{0, 2} {
		c := base("season", 0x19377, 0)
		c.patches = []residualPatch{csPatch(0xCF0, 2), csPatch(0xCF3, cf3), csPatch(0xCF4, 3)}
		add(c)
	}
	{
		c := base("season", 0x19377, 1)
		c.patches = []residualPatch{csPatch(0xCF0, 0x11), csPatch(0xCF3, 1), csPatch(0xCF4, 6)}
		add(c)
	}
	for _, month := range []byte{3, 6, 9, 12, 1, 5, 13} {
		for _, day := range []byte{1, 2, 7} {
			c := base("season", 0x19377, int(month)%4)
			c.patches = []residualPatch{csPatch(0xCF0, day), csPatch(0xCF3, 1), csPatch(0xCF4, month)}
			add(c)
		}
	}
	for _, ah := range []byte{0, 1, 2, 3, 0xFF} {
		for _, al := range []byte{2, 5} {
			c := base("season", 0x193D7, 0)
			c.ax = uint16(al) | uint16(ah)<<8
			add(c)
		}
	}

	// tactical-cap: sub_1AB7C + raw 0x1AB48 tail entry. Record at DS:0x100.
	for _, row := range []struct {
		cap, typ byte
	}{
		{0, 0}, {39, 35}, {40, 0}, {40, 36}, {40, 37},
		{41, 0}, {41, 36}, {41, 37}, {41, 0x48}, {0xFF, 36}, {0xFF, 0xFF},
	} {
		c := base("tactical-cap", 0x1AB7C, 0)
		c.ds = 0x8000
		c.si = 0x100
		c.patches = []residualPatch{dsPatch(0x119, row.cap), dsPatch(0x104, row.typ)}
		add(c)
	}
	for _, row := range []struct {
		cap, typ byte
	}{
		{41, 0}, {30, 37},
	} {
		c := base("tactical-cap", 0x1AB48, 1)
		c.ds = 0x8000
		c.si = 0x100
		c.patches = []residualPatch{dsPatch(0x119, row.cap), dsPatch(0x104, row.typ)}
		add(c)
	}

	// sidebar: sub_1C300; adc bl,dl is CF-sensitive. DS:0x200 cell.
	for _, row := range []struct {
		cf     bool
		bl, dl byte
		dh     byte
		di     uint16
		cell   byte
	}{
		{false, 0x00, 0x01, 0x00, 0x200, 0xFF},
		{true, 0x00, 0x01, 0x00, 0x200, 0xFF},
		{false, 0xFF, 0x01, 0xF0, 0x200, 0xFF},
		{true, 0xFF, 0xFF, 0xF0, 0x200, 0x3C},
		{false, 0x12, 0x34, 0x0F, 0x400, 0xA5},
		{true, 0x12, 0x34, 0xF0, 0x400, 0xA5},
		{false, 0x80, 0x80, 0xFF, 0x1FFE, 0x00},
		{true, 0x80, 0x80, 0xFF, 0x1FFE, 0x55},
	} {
		c := base("sidebar", 0x1C300, 0)
		c.ds = 0x8000
		c.cf = row.cf
		c.bx = uint16(row.bl)
		c.dx = uint16(row.dl) | uint16(row.dh)<<8
		c.di = row.di
		c.patches = []residualPatch{dsPatch(row.di+2, row.cell)}
		add(c)
	}

	// map-node / draw-queue: 8-byte record at DS:SI + counter to ES:DI (es=dx).
	for _, grp := range []struct {
		name   string
		target uint32
	}{
		{"map-node", 0x1E567}, {"draw-queue", 0x1E6FF},
	} {
		for i, row := range []struct {
			si, di, bp, cx uint16
			ax             uint16
			counter        byte
		}{
			{0x300, 0x400, 0x1111, 0x2222, 0x3344, 0xC0},
			{0x308, 0x410, 0x0000, 0xFFFF, 0x00FF, 0x00},
			{0x310, 0x420, 0xABCD, 0x1234, 0xFF00, 0xFF},
			{0x7F0, 0x7F8, 0x5555, 0xAAAA, 0x5A5A, 0x3F},
		} {
			c := base(grp.name, grp.target, i%4)
			c.ds = 0x8000
			c.dx = 0x8100
			c.si, c.di, c.bp, c.cx, c.ax = row.si, row.di, row.bp, row.cx, row.ax
			off := uint16(0xE47E)
			if grp.target == 0x1E6FF {
				off = 0xE47F
			}
			c.patches = []residualPatch{csPatch(off, row.counter)}
			add(c)
		}
	}
	return all
}

func residualAudit(before, after []byte, cs uint16, c residualCase, entry registers, observed []registers, device *machine.Machine, d *dos.DOS, soundBefore map[uint8]int, allocsBefore, readsBefore int) (int, []string) {
	n := 0
	var obs []string
	check := func(cond bool, msg string) {
		if !cond {
			obs = append(obs, "FAIL "+msg)
		} else {
			obs = append(obs, "ok "+msg)
		}
		n++
	}
	u8 := func(mem []byte, seg, off uint16) byte { return mem[int(seg)*16+int(off)] }
	_ = u8
	var ips []uint16
	for _, o := range observed {
		ips = append(ips, o.IP)
	}
	obs = append(obs, fmt.Sprintf("entries %d first=%04X last=%04X", len(ips), func() uint16 {
		if len(ips) > 0 {
			return ips[0]
		}
		return 0
	}(), func() uint16 {
		if len(ips) > 0 {
			return ips[len(ips)-1]
		}
		return 0
	}()))
	n++
	final := originalRegs(device)
	switch c.group {
	case "diplomacy":
		a, b := byte(entry.AX), byte(entry.AX>>8)
		for _, half := range []struct {
			self, peer byte
		}{{a, b}, {b, a}} {
			at := int(entry.DS)*16 + int(half.self)*64 + 0x19
			want := before[at]
			if before[at] == half.peer {
				want = 0xFF
			}
			check(after[at] == want, fmt.Sprintf("diplo %02X/%02X", half.self, half.peer))
		}
		_ = final
	case "tactical-cap":
		at := int(entry.DS)*16 + int(entry.SI) + 0x19
		want := before[at]
		if want > 40 {
			want = 40
		}
		check(after[at] == want, "cap clamp")
		typ := before[int(entry.DS)*16+int(entry.SI)+4]
		path := "abd2"
		if typ == 36 {
			path = "abff"
		} else if typ > 36 {
			path = "ac55"
		}
		obs = append(obs, "dispatch "+path)
		n++
	case "sidebar":
		sum := uint16(byte(entry.BX)) + uint16(byte(entry.DX))
		if entry.Flags&1 != 0 {
			sum++
		}
		check(byte(final.BX) == byte(sum), "adc bl")
		cell := before[int(entry.DS)*16+int(entry.DI)+2] & byte(entry.DX>>8)
		check(after[int(entry.DS)*16+int(entry.DI)+2] == cell, "and cell")
		check(final.AX == 5, "al=5 ah=0")
	case "map-node", "draw-queue":
		base := int(entry.DS)*16 + int(entry.SI)
		check(after[base] == byte(entry.DI) && after[base+1] == byte(entry.DI>>8), "rec di")
		check(after[base+2] == byte(entry.BP) && after[base+3] == byte(entry.BP>>8), "rec bp")
		check(after[base+4] == byte(entry.CX) && after[base+5] == byte(entry.CX>>8), "rec cx")
		check(after[base+6] == byte(entry.AX>>8) && after[base+7] == byte(entry.AX), "rec ax")
		off := uint16(0xE47E)
		if c.target == 0x1E6FF {
			off = 0xE47F
		}
		check(after[int(entry.DX)*16+int(entry.DI)] == before[int(cs)*16+int(off)], "counter copy")
	case "season":
		if c.target == 0x19377 {
			month := before[int(cs)*16+0xCF4]
			day := before[int(cs)*16+0xCF0]
			gate := before[int(cs)*16+0xCF3] == 1 && day <= 0x10
			path := "ret"
			switch {
			case !gate:
				path = "ret"
			case month == 3:
				path = "al2"
			case month == 6:
				path = "al3"
			case month == 9:
				path = "al4"
			case month == 12:
				path = "al5"
			}
			obs = append(obs, "season "+path)
			n++
		} else {
			ah := byte(entry.AX >> 8)
			path := "ret"
			if ah == 1 {
				path = "stop"
			} else if ah == 2 {
				path = "play"
			}
			obs = append(obs, "bgm "+path)
			n++
		}
		var soundKeys []uint8
		for ah := range d.Sound {
			soundKeys = append(soundKeys, ah)
		}
		sort.Slice(soundKeys, func(i, j int) bool { return soundKeys[i] < soundKeys[j] })
		for _, ah := range soundKeys {
			obs = append(obs, fmt.Sprintf("int61 AH=%02X %d->%d", ah, soundBefore[ah], d.Sound[ah]))
			n++
		}
	case "picker-row":
		row := byte(entry.AX)
		right := row >= 11
		r := row
		if right {
			r -= 11
		}
		obs = append(obs, fmt.Sprintf("cell right=%v y=%d", right, int(r)*16+0xD8))
		n++
		empty := before[int(entry.DS)*16+int(row)*64] < 0x80
		obs = append(obs, fmt.Sprintf("slot empty=%v", empty))
		n++
	case "picker-loop", "picker-entry":
		obs = append(obs, fmt.Sprintf("mouse queries=%d", len(d.Mouse.PressReads)))
		n++
		obs = append(obs, fmt.Sprintf("sel=%02X", after[int(cs)*16+0x98A7]))
		n++
	}
	_ = allocsBefore
	_ = readsBefore
	return n, obs
}
