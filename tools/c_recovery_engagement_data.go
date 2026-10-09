//go:build matching_engagement

package main

import (
	"bytes"
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

type engagementPatch struct {
	segment, offset uint16
	bytes           []byte
}
type engagementVector struct {
	active       bool
	field, flags byte
	phase        string
	patches      []engagementPatch
	warmup       bool
	worldWarmup  bool
	ds, es       uint16
	postPatches  []engagementPatch
	fullEntrySP  uint16
}

// FFFF selects this side's CS; FFFE selects its own initialized tactical
// unit segment. Neither sentinel reads the other execution's RAM.
func engagementResolveSegment(mem []byte, cs, segment uint16) uint16 {
	switch segment {
	case 0xFFFF:
		return cs
	case 0xFFFE:
		return strategyWord(mem, int(cs)*16+0xD30E)
	default:
		return segment
	}
}

func engagementApplyPatches(mem []byte, cs uint16, patches []engagementPatch) {
	for _, patch := range patches {
		segment := engagementResolveSegment(mem, cs, patch.segment)
		for i, value := range patch.bytes {
			mem[mainPhysical(segment, patch.offset+uint16(i))] = value
		}
	}
}

// Only raw caller input is prepared here.19946/1CC31/19A33 and the tactical
// engine must actually create their own buffers, counters and savedSP.
func engagementPrepare(mem []byte, cs uint16, c resumeCase) {
	if !c.engagement.active {
		return
	}
	code := int(cs) * 16
	my, enemy := c.si, c.di
	// Warmed helpers may use SI/DI for tactical units or UI arguments. The
	// original 11B5A prefix still starts from the two typed world armies.
	if len(c.outcome.armies) >= 2 {
		my = uint16(0x2240 + c.outcome.armies[0].slot*64)
		enemy = uint16(0x2240 + c.outcome.armies[1].slot*64)
	}
	if c.engagement.flags&0x80 != 0 {
		my, enemy = enemy, my
	}
	strategyPutWord(mem, code+0x0D2E, my)
	strategyPutWord(mem, code+0x0D30, enemy)
	city := uint16(0)
	if c.engagement.field < 0xC0 {
		city = uint16(0x840 + c.outcome.city*32)
	}
	strategyPutWord(mem, code+0x0D32, city)
	mem[code+0x0D34], mem[code+0x0D35] = c.engagement.field, c.engagement.flags
	mem[code+0x0CFC] = 0 // Original unthrottled branch, not a fabricated timer tick.
	strategyPutWord(mem, code+0xD318, 0)
	for _, army := range c.outcome.armies {
		a := outcomeWorld + 0x2240 + army.slot*64
		node := strategyWord(mem, a+0x0E)
		if node < 0x600 {
			record := outcomeWorld + 0x840 + int(node)*4
			strategyPutWord(mem, a+0x10, strategyWord(mem, record+8))
			strategyPutWord(mem, a+0x12, strategyWord(mem, record+10))
		}
		x, y := strategyWord(mem, a+0x10), strategyWord(mem, a+0x12)
		segment := strategyWord(mem, code+0x9872) + y*24
		strategyPutWord(mem, a+0x1A, x)
		strategyPutWord(mem, a+0x1C, segment)
		mem[mainPhysical(segment, x)] = 1
	}
	engagementApplyPatches(mem, cs, c.engagement.patches)
	// postPatches belong after each side's independently executed prefix.
}

func engagementClone(c resumeCase) resumeCase {
	c.outcome.armies = append([]outcomeArmy(nil), c.outcome.armies...)
	c.outcome.cities = append([]outcomeCity(nil), c.outcome.cities...)
	c.outcome.generals = append([]outcomeGeneral(nil), c.outcome.generals...)
	c.outcome.factions = append([]outcomeFaction(nil), c.outcome.factions...)
	c.outcome.redirectSlots = append([]int(nil), c.outcome.redirectSlots...)
	c.route.nodes = append([]routeNode(nil), c.route.nodes...)
	c.route.links = append([]routeLink(nil), c.route.links...)
	c.route.visited = append([]routeVisit(nil), c.route.visited...)
	c.route.queue = append([]routeQueue(nil), c.route.queue...)
	c.events = append([]listEvent(nil), c.events...)
	c.mapEvents = append([]strategyEvent(nil), c.mapEvents...)
	c.menuChoices = append([]int(nil), c.menuChoices...)
	c.engagement.patches = append([]engagementPatch(nil), c.engagement.patches...)
	for i := range c.engagement.patches {
		c.engagement.patches[i].bytes = append([]byte(nil), c.engagement.patches[i].bytes...)
	}
	c.engagement.postPatches = append([]engagementPatch(nil), c.engagement.postPatches...)
	for i := range c.engagement.postPatches {
		c.engagement.postPatches[i].bytes = append([]byte(nil), c.engagement.postPatches[i].bytes...)
	}
	return c
}

func engagementCases(base func(string, uint32, int) resumeCase) []resumeCase {
	templates := outcomeCases(base)
	var selected [4]resumeCase
	var found [4]bool
	for _, c := range templates {
		if c.scenario >= 0 && c.scenario < 4 && !found[c.scenario] {
			selected[c.scenario] = c
			found[c.scenario] = true
		}
	}
	var cases []resumeCase
	for scenario := 0; scenario < 4; scenario++ {
		if !found[scenario] {
			panic("engagement lacks original scenario template")
		}
		c := engagementClone(selected[scenario])
		c.group, c.target = "engagement", 0x11B5A
		c.engagement = engagementVector{active: true, field: 0, flags: 0, phase: "retreat-ui"}
		c.outcome.escape = false
		c.outcome.mode = 0
		c.outcome.reportCity = 0
		c.outcome.rngSeed = 3
		c.outcome.customCoefficients = false
		c.outcome.panel = 0
		c.outcome.redirectSlots = nil
		c.si, c.di, c.ax, c.bx, c.cx, c.dx = 0x2240, 0x2280, 0, 0x2280, 0, 0
		c.events = []listEvent{{100, 120, 1}}
		c.mapEvents = nil
		c.menuChoices = nil
		for side := 0; side < 2; side++ {
			a := &c.outcome.armies[side]
			a.flags, a.morale, a.stage = 0xC0, 100, 8
			a.men = [6]byte{1, 0, 0, 0, 0, 0}
			a.types = [6]byte{3, 3, 3, 3, 3, 3}
			a.current, a.target = 0, 0
			if side == 0 {
				a.current = 8
			}
			g := &c.outcome.generals[side]
			g.flags, g.duty, g.captor, g.temperament, g.war, g.lead = 0x80, 1, 255, 3, 8, 8
			g.aptitudes = [3]byte{}
		}
		foreign := c.outcome.armies[1].owner
		c.outcome.cities[0].owner, c.outcome.cities[0].oldOwner = foreign, foreign
		c.outcome.cities[1].owner, c.outcome.cities[1].oldOwner = c.player, c.player
		for i := range c.outcome.factions {
			f := &c.outcome.factions[i]
			if f.index == int(c.player) {
				f.capital = 1
				f.ruler = 0
			} else if f.index == int(foreign) {
				f.capital = 0
				f.ruler = 1
			}
		}
		cases = append(cases, c)
	}
	bases := make([]resumeCase, len(cases))
	for i, c := range cases {
		bases[i] = engagementClone(c)
	}
	cases = append(cases, engagementControlCases(bases)...)
	cases = append(cases, engagementUICases(bases)...)
	cases = append(cases, engagementFrontCases(bases)...)
	cases = append(cases, engagementExtraCases(bases)...)
	cases = append(cases, engagementCombatCases(bases)...)
	return cases
}

// This is a front-end/flow oracle, not a complete independent tactical-math
// model. Original/C full-state comparison remains the main tactical gate.
func engagementAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, s string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent engagement scenario=%d target=%X: %s", c.scenario, c.target, s))
		}
	}
	check(c.engagement.active, "active tactical fixture")
	code := int(cs) * 16
	check(strategyWord(before, code+0x0D46) == strategyWord(before, code+0x0D44), "arena/decoded-map base agrees before full restore")
	check(before[code+0x0CFC] == 0, "original unthrottled input")
	if !c.engagement.warmup {
		check(before[code+0x0D34] == c.engagement.field && before[code+0x0D35] == c.engagement.flags, "raw field/flags input")
		for _, pointer := range []uint16{strategyWord(before, code+0x0D2E), strategyWord(before, code+0x0D30)} {
			a := outcomeWorld + int(pointer)
			sum := uint16(0)
			check(pointer >= 0x2240 && pointer < 0x4200 && (pointer-0x2240)%64 == 0, "legal live world army record")
			for slot := 0; slot < 6; slot++ {
				sum += uint16(before[a+0x29+slot*4])
				kind := before[a+0x2A+slot*4]
				check(kind >= 1 && kind <= 3, "legal tactical troop type")
			}
			check(sum == strategyWord(before, a+4) && before[a+0x29] >= 1 && before[a+6] == 100, "consistent small army and morale")
		}
		var initial [258]byte
		initial[0], initial[1] = c.outcome.rngSeed^0xA5, c.outcome.rngSeed
		for i := 0; i < 256; i++ {
			initial[i+2] = byte(i*73 + i/2 + int(c.outcome.rngSeed))
		}
		check(bytes.Equal(before[code+0xECFC:code+0xEDFE], initial[:]), "fixed raw RNG input")
	}
	// Replay only the already-observed RNG entry count using the old proven
	// recurrence. This does not claim independently derived tactical call count.
	rng := &routeModel{mem: append([]byte(nil), before...), cs: cs, calls: map[uint16]int{}}
	for range strategyObserved(observed, 0xECE0) {
		rng.random()
	}
	check(bytes.Equal(after[code+0xECFC:code+0xEDFE], rng.mem[code+0xECFC:code+0xEDFE]), "observed RNG recurrence/table preservation")
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "strategic date preserved")
	arena := strategyWord(before, code+0x0D46)
	checkMappings := func(restored bool) {
		mappings := [][2]uint16{{0xD2FA, 0}, {0xD2FC, 0x700}, {0xD2FE, 0x900}, {0xD300, 0xB00}, {0xD2F6, 0xF80}, {0xD2F8, 0x1080}, {0xD302, 0x1180}, {0xD304, 0x2100}, {0xD306, 0x3D20}, {0xD30E, 0x459A}, {0xD308, 0x4A1A}, {0xD30A, 0x4A2A}, {0xD30C, 0x4A2E}}
		if restored {
			mappings = append(mappings, [][2]uint16{{0x9872, 0x1800}, {0x9876, 0x3000}, {0x987A, 0x4200}, {0x9878, 0x4610}, {0x0D42, 0x47DC}, {0x9874, 0x48D6}, {0x987C, 0x50D6}}...)
		} else {
			mappings = append(mappings, [2]uint16{0x0D42, 0x44A0})
		}
		for _, mapping := range mappings {
			check(strategyWord(after, code+int(mapping[0])) == arena+mapping[1], fmt.Sprintf("real initializer/restorer pointer%04X", mapping[0]))
		}
	}
	fullSP := c.engagement.fullEntrySP
	if fullSP == 0 {
		fullSP = c.entrySP
	}
	saved := fullSP - 18
	if c.engagement.phase == "warmup-pause" {
		// 11B5A pushes eight words; CALL19FA0 supplies saved return1B76.
		// 19FA6 saves SP, then CALL1A156 at19FBB supplies return9FBE.
		for _, ip := range []uint16{0x1B5A, 0x9946, 0xCC31, 0x9A33, 0x9FA0} {
			check(len(strategyObserved(observed, ip)) == 1, fmt.Sprintf("warmup real entry%04X", ip))
		}
		loop := strategyObserved(observed, 0x9FA0)[0]
		check(loop.SS == c.entrySS && loop.SP == saved, "warmup original saved return frame")
		check(strategyWord(after, code+0xD340) == saved, "warmup19FA6 saved original SP")
		check(device.CPU.Seg[cpu.CS] == cs && device.CPU.IP == 0xA156 && device.CPU.Seg[cpu.SS] == c.entrySS && device.CPU.R[cpu.SP] == saved-2, "pause at first actual1A156 entry")
		check(strategyWord(after, mainPhysical(c.entrySS, saved-2)) == 0x9FBE && strategyWord(after, mainPhysical(c.entrySS, saved)) == 0x1B76, "live original polling/cleanup return words")
		check(len(strategyObserved(observed, 0xA065)) >= 1 && len(strategyObserved(observed, 0xA6FA)) >= 1, "warmup executed real tactical updates")
		for _, ip := range []uint16{0x9FDC, 0x1B76, 0x19CA, 0x1CB1} {
			check(len(strategyObserved(observed, ip)) == 0, fmt.Sprintf("warmup has not escaped%04X", ip))
		}
		checkMappings(false)
		return checks
	}
	if c.engagement.warmup {
		entrySP := saved - 2
		if c.engagement.worldWarmup {
			// A complete 11B5A has consumed the saved tactical frame and
			// restored the world. Its next independent CALL starts at P.
			entrySP = fullSP
			check(strategyWord(before, code+0x9872) == arena+0x1800 && strategyWord(before, code+0x0D42) == arena+0x47DC && strategyWord(before, code+0x987C) == arena+0x50D6, "world helper starts after real main restore")
		}
		check(c.engagement.fullEntrySP != 0 && c.entrySP == entrySP, "direct helper uses the correct prefix return frame")
		check(strategyWord(before, code+0xD340) == saved && strategyWord(after, code+0xD340) == saved, "direct helper saved tactical SP")
		entries := strategyObserved(observed, uint16(c.target))
		check(len(entries) >= 1, "real direct helper entry")
		entry := entries[0]
		entryDS, entryES := cs, cs
		if c.engagement.ds != 0 {
			entryDS = engagementResolveSegment(before, cs, c.engagement.ds)
		}
		if c.engagement.es != 0 {
			entryES = engagementResolveSegment(before, cs, c.engagement.es)
		}
		check(entry.SS == c.entrySS && entry.SP == c.entrySP && entry.DS == entryDS && entry.ES == entryES, "direct helper own-side register/frame input")
		if len(strategyObserved(observed, 0x1CB1)) != 0 {
			check(c.outcome.escape && device.CPU.Seg[cpu.CS] == cs && device.CPU.IP == c.returnIP && device.CPU.Seg[cpu.SS] == c.savedSS && device.CPU.R[cpu.SP] == c.savedSP+2, "direct helper world nonlocal return")
			check(strategyWord(after, code+0x9901) == c.savedSP && strategyWord(after, code+0x9903) == c.savedSS, "world nonlocal saved frame preserved")
		} else if len(strategyObserved(observed, 0x9FDC)) != 0 {
			check(len(strategyObserved(observed, 0x9FDC)) == 1 && len(strategyObserved(observed, 0x1B76)) == 1, "direct helper true tactical nonlocal cleanup")
			cleanup := strategyObserved(observed, 0x1B76)[0]
			check(cleanup.SS == c.entrySS && cleanup.SP == saved+2, "direct1A04A RET consumed original1B76")
			check(device.CPU.Seg[cpu.CS] == cs && device.CPU.IP == 0xF000 && device.CPU.Seg[cpu.SS] == c.entrySS && device.CPU.R[cpu.SP] == fullSP+2, "direct tactical cleanup returns full outer caller")
			for _, ip := range []uint16{0x19CA, 0x533D, 0x89F0} {
				check(len(strategyObserved(observed, ip)) == 1, fmt.Sprintf("direct true world restore%04X", ip))
			}
			checkMappings(true)
		} else {
			check(!c.outcome.escape && len(strategyObserved(observed, 0x1B76)) == 0, "normal direct helper did not discard caller")
			check(device.CPU.Seg[cpu.CS] == cs && device.CPU.IP == 0xF000 && device.CPU.Seg[cpu.SS] == c.entrySS && device.CPU.R[cpu.SP] == c.entrySP+2, "normal direct helper return")
		}
		return checks
	}
	for _, ip := range []uint16{0x1B5A, 0x9946, 0xCC31, 0x9A33, 0x9FA0, 0x9FDC, 0x1B76, 0x19CA, 0x533D, 0x89F0} {
		check(len(strategyObserved(observed, ip)) == 1, fmt.Sprintf("real original entry%04X", ip))
	}
	check(len(strategyObserved(observed, 0xA065)) >= 1 && len(strategyObserved(observed, 0xA6FA)) >= 1, "real tactical update and finish decision")
	check(len(strategyObserved(observed, 0x1CB1)) == 0, "direct tactical return did not enter world nonlocal exit")
	loop, exit, cleanup := strategyObserved(observed, 0x9FA0)[0], strategyObserved(observed, 0x9FDC)[0], strategyObserved(observed, 0x1B76)[0]
	check(loop.SS == c.entrySS && loop.SP == saved, "real CALL19FA0 owns saved return frame")
	check(strategyWord(after, code+0xD340) == saved, "19FA6 saved original SP")
	check(exit.SS == c.entrySS && exit.SP < saved, "19FDC reached through a deeper original caller")
	check(cleanup.SS == c.entrySS && cleanup.SP == saved+2, "1A04A RET consumed saved1B76 return")
	seenExit, seenCleanup := false, false
	for _, r := range observed {
		if r.CS != cs {
			continue
		}
		if r.IP == 0x9FDC {
			seenExit = true
		}
		if r.IP == 0x1B76 {
			check(seenExit, "cleanup follows real nonlocal tactical finish")
			seenCleanup = true
		}
		if seenExit && !seenCleanup {
			check(r.IP != 0xA065 && r.IP != 0xA6FA && r.IP != 0x9FA0, "discarded caller did not resume")
		}
	}
	check(device.CPU.Seg[cpu.SS] == c.entrySS && device.CPU.R[cpu.SP] == c.entrySP+2 && device.CPU.IP == 0xF000, "full11B5A normal outer return")
	check(byte(device.CPU.R[cpu.AX]) <= 1 && byte(device.CPU.R[cpu.AX]>>8) <= 3, "battle AL winner/AH broken bits range")
	if len(strategyObserved(observed, 0xA156)) != 0 && (c.engagement.phase == "retreat-ui" || c.engagement.phase == "world-warmup") {
		check(len(strategyObserved(observed, 0xC21A)) >= 1 && len(strategyObserved(observed, 0xA8F6)) >= 1, "real retreat button input reached original handler")
	}
	checkMappings(true)
	return checks
}
