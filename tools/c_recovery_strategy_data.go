//go:build matching_strategy

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const strategyWorld = 0x70000
const strategyQueue = 0xD0000
const strategySS = 0xF0000

func strategyWord(mem []byte, at int) uint16 {
	return binary.LittleEndian.Uint16(mem[at : at+2])
}

func strategyPutWord(mem []byte, at int, value uint16) {
	binary.LittleEndian.PutUint16(mem[at:at+2], value)
}

func strategyForeign(world []byte) int {
	for faction := 1; faction < 22; faction++ {
		if world[faction*64] >= 0x80 {
			return faction
		}
	}
	return -1
}

// The ordinary general list includes employed generals and the ruler.
func strategyGeneral(world []byte, player byte) int {
	for general := 0; general < 128; general++ {
		at := 0x4240 + general*32
		if world[at] >= 0x80 && world[at+0x1C] == player {
			return general
		}
	}
	return -1
}

func strategyMoney(mem []byte, at int, value int32) {
	raw := uint32(value) & 0xFFFFFF
	mem[at], mem[at+1], mem[at+2] = byte(raw), byte(raw>>8), byte(raw>>16)
}

func strategyRelation(player, target int) int {
	// sub_130CB -> sub_13119, as independently checked in world_update.c:
	// SI/4 + SI/8 + target + 0600h, with SI=player*64.
	return strategyWorld + 0x600 + player*24 + target
}

func strategyPrepare(mem []byte, cs uint16, c resumeCase) {
	code := int(cs) * 16
	strategyPutWord(mem, code+0x0CFD, uint16(c.player)*64)
	mem[code+0x0CFF] = c.player
	mem[code+0x0D00] = c.trust
	strategyPutWord(mem, code+0x0D56, 0xD000)
	strategyPutWord(mem, code+0x0D20, 0)
	// D22 is not a tail consumed by the original writer; leave it untouched.
	clear(mem[strategyQueue : strategyQueue+0x400])
	for off := 0x98A8; off <= 0x98AB; off++ {
		mem[code+off] = 0
	}
	player, foreign := strategyWorld+int(c.player)*64, strategyWorld+c.owner*64
	lord := strategyWorld + 0x4240 + int(mem[player+1])*32
	mem[lord+0x17] = c.lordDuty
	mem[player+0x28], mem[foreign+0x28] = c.aggression, c.aggression
	mem[player+0x23] = c.ownCities
	if foreign != player {
		mem[foreign+0x23] = c.foreignCities
	}
	mem[foreign+0x19] = c.ally
	strategyMoney(mem, player+0x20, c.ownMoney)
	if foreign != player {
		strategyMoney(mem, foreign+0x20, c.foreignMoney)
	}
	mem[strategyRelation(int(c.player), c.owner)] = c.relation
	if int(c.candidate) != c.owner {
		mem[strategyRelation(int(c.player), int(c.candidate))] = c.candidateRelation
	}
	if c.group == "general" {
		general := strategyWorld + 0x4240 + c.general*32
		mem[general+0x1D] = c.captor
		copy(mem[general+0x0E:general+0x11], c.aptitudes[:])
	}
	strategyPutWord(mem, code+0x0D02, uint16(c.income))
	mem[code+0x0D04] = byte(c.income >> 16)
	strategyPutWord(mem, code+0x0D05, uint16(c.expense))
	mem[code+0x0D07] = byte(c.expense >> 16)
	for index, value := range c.fiscal {
		strategyPutWord(mem, code+0x0D08+index*2, value)
	}
	mem[code+0x020F] = c.priority
	frame := strategySS + int(c.bp)
	if c.group == "proposal" && c.target == 0x166D9 {
		strategyPutWord(mem, frame, uint16(c.owner)*64)
		strategyPutWord(mem, frame+2, uint16(c.candidate)*64)
	}
	if c.group == "reason" || c.group == "reason-loop" {
		// BP+0 is the caller's chosen message base and is not inferred here.
		mem[frame+2], mem[frame+3] = c.verdict, c.budget
		mem[frame+4], mem[frame+5] = c.reasonMask, c.attempted
	}
	// Controlled queue profiles: 0 empty, 1 first slot matches, 2 first slot
	// has a different key, 3 full. These are fixtures, not production policy.
	queueGroup := c.group == "queue-search" || c.group == "queue-write" || c.group == "proposal" || c.group == "advice"
	if queueGroup && c.profile != 0 {
		key := uint16(c.player)<<8 | 1
		data := byte(c.owner)
		if c.group == "queue-search" || c.group == "queue-write" {
			key = uint16(byte(c.ax)) | uint16(byte((c.si<<2)>>8))<<8
			data = byte(c.dx)
		}
		count := 1
		if c.profile == 3 {
			count = 256
			if byte(key) == 0 {
				key |= 1 // A full fixture must contain nonempty low bytes.
			}
		}
		if c.profile == 2 {
			key ^= 0x0100 // Remain unmatched even when both data bytes are wildcards.
		}
		for index := 0; index < count; index++ {
			at := strategyQueue + index*4
			strategyPutWord(mem, at, key)
			mem[at+2], mem[at+3] = data, data
			if c.profile != 1 {
				mem[at+2]++
			}
		}
	}
}

func strategyQueueFound(mem []byte, cs uint16, key uint16, dx uint16) bool {
	base := int(strategyWord(mem, int(cs)*16+0x0D56)) * 16
	start := int(strategyWord(mem, int(cs)*16+0x0D20))
	for offset := start; offset < 0x400; offset += 4 {
		if strategyWord(mem, base+offset) != key {
			continue
		}
		low, high := byte(dx), byte(dx>>8)
		if low != 0xFF && mem[base+offset+2] != low {
			continue
		}
		// Original 1304E compares +3 with DL when DH is not the wildcard.
		if high != 0xFF && mem[base+offset+3] != low {
			continue
		}
		return true
	}
	return false
}

func strategyExpectedQueueSearch(before []byte, cs uint16, c resumeCase) bool {
	key := uint16(byte(c.ax)) | uint16(byte((c.si<<2)>>8))<<8
	return strategyQueueFound(before, cs, key, c.dx)
}

func strategyPower(mem []byte, own, foreign int) (uint16, uint16) {
	ours := uint16(mem[strategyWorld+own*64+0x23]) * uint16(mem[strategyWorld+own*64+0x28]+20)
	theirs := uint16(mem[strategyWorld+foreign*64+0x23]) * 25
	return ours, theirs
}

func strategyExpectedPower(before []byte, cs uint16, c resumeCase) (uint16, uint16) {
	return strategyPower(before, int(c.player), c.owner)
}

// Returns the original low bytes AL and DL, including unchanged early-return
// DL. Flags and the other return registers are checked by the harness.
func strategyExpectedProposal(before []byte, cs uint16, c resumeCase) (byte, byte) {
	player := int(strategyWord(before, int(cs)*16+0x0CFD)) / 64
	owner := c.owner
	ours, theirs := strategyPower(before, player, owner)
	aggression := before[strategyWorld+player*64+0x28]
	relation := before[strategyRelation(player, owner)]
	ally := before[strategyWorld+owner*64+0x19]
	// Each original condition inserts 10h before SHR DL,1. The first
	// insertion passes through four shifts to bit 0; the last reaches bit 3.
	mask := byte(0)
	switch c.target {
	case 0x16475:
		if strategyQueueFound(before, cs, uint16(player)<<8|1, 0xFF00|uint16(byte(owner))) {
			return 1, byte(owner)
		}
		if relation < 0x80 {
			return 3, byte(owner)
		}
		low, threshold := relation-0x80, aggression*2+20
		if low >= threshold {
			return 0, byte(owner)
		}
		if low < threshold/2+5 {
			mask |= 1
		}
		if ours > theirs {
			mask |= 2
		}
		if ally != 0xFF && ally != byte(player) {
			mask |= 4
		}
		if before[strategyWorld+owner*64+0x22]&0x80 != 0 {
			mask |= 8
		}
		return 2, mask
	case 0x16577:
		if relation >= 0x80 {
			return 3, byte(c.dx)
		}
		if relation < aggression/2 {
			return 0, byte(c.dx)
		}
		if ours < theirs {
			mask |= 1
		}
		for faction := 0; faction < 22; faction++ {
			at := strategyWorld + faction*64
			if faction != owner && before[at] >= 0x80 && before[at+0x19] == byte(player) {
				mask |= 2
			}
		}
		if ally != 0xFF && ally != byte(player) {
			mask |= 4
		}
		if before[strategyWorld+player*64+0x22]&0x80 != 0 {
			mask |= 8
		}
		return 2, mask
	case 0x166D9:
		frame := strategySS + int(c.bp)
		owner = int(strategyWord(before, frame)) / 64
		candidate := int(strategyWord(before, frame+2)) / 64
		if owner == candidate {
			return 4, byte(c.dx)
		}
		relation = before[strategyRelation(player, owner)]
		threshold := (aggression*4 + 30) | 0x80
		if relation < threshold {
			return 0, byte(c.dx)
		}
		if before[strategyRelation(player, candidate)] >= 0x80 {
			return 3, byte(c.dx)
		}
		ours, theirs = strategyPower(before, player, candidate)
		candidateAlly := before[strategyWorld+candidate*64+0x19]
		if candidateAlly == byte(player) && ours < theirs/2 {
			return 1, byte(c.dx)
		}
		if relation >= threshold+30 {
			mask |= 1
		}
		ours, theirs = strategyPower(before, player, owner)
		if ours < theirs {
			mask |= 2
		}
		ours, theirs = strategyPower(before, player, candidate)
		if ours < theirs {
			mask |= 4
		}
		if candidateAlly == byte(player) {
			mask |= 8
		}
		return 2, mask
	default:
		panic("independent strategy proposal target")
	}
}

func strategyTrustBudget(trust byte) byte {
	if trust >= 0xE0 {
		return 1
	}
	if trust >= 0x90 {
		return 2
	}
	if trust >= 0x20 {
		return 3
	}
	return 4
}

func strategyTrustAdd(trust, amount byte) byte {
	value := uint16(trust) + uint16(amount)
	if value > 255 {
		return 255
	}
	return byte(value)
}

func strategyTrustSubtract(trust, amount byte) byte {
	if amount > trust {
		return 0
	}
	return trust - amount
}

type strategyReasonOutcome struct {
	Response, Attempted, Budget, Trust byte
	Carry                              bool
	TalkOffset                         uint16
}

func strategyReasonStep(reason, mask byte, state strategyReasonOutcome) strategyReasonOutcome {
	if reason == 4 {
		state.Response, state.Carry, state.TalkOffset = 4, true, 0x2A
		return state
	}
	if reason > 3 {
		panic("independent reason fixture outside selector rows 0..4")
	}
	bit := byte(1 << reason)
	if state.Attempted&bit != 0 {
		state.Response, state.Carry, state.TalkOffset = 3, false, 0x2D
		return state
	}
	state.Attempted |= bit
	if mask&bit == 0 {
		if state.Trust < 20 {
			panic("reason trust underflow uses original nonlocal startup-stack restore")
		}
		state.Trust = strategyTrustSubtract(state.Trust, 20)
		state.Response, state.Carry = 0, true
	} else {
		state.Budget-- // Original byte DEC wraps when the input budget is zero.
		state.Response, state.Carry = 2, false
		if state.Budget == 0 {
			state.Response, state.Carry = 1, true
		}
	}
	state.TalkOffset = uint16(reason)*9 + uint16(state.Response)*3 + 6
	return state
}

func strategyExpectedReason(before []byte, cs uint16, c resumeCase) strategyReasonOutcome {
	frame := strategySS + int(c.bp)
	state := strategyReasonOutcome{Attempted: before[frame+5], Budget: before[frame+3], Trust: before[int(cs)*16+0x0D00]}
	return strategyReasonStep(c.reason, before[frame+4], state)
}

func strategyReasonLoop(choices []int, mask byte, state strategyReasonOutcome) (strategyReasonOutcome, bool) {
	state.Attempted = 0
	for _, choice := range choices {
		if choice < 0 {
			continue // 13B7E reopens its selector after cancellation.
		}
		state = strategyReasonStep(byte(choice), mask, state)
		if state.Carry {
			state.Carry = state.Response != 1
			return state, true
		}
	}
	return state, false
}

func strategyExpectedReasonLoop(before []byte, cs uint16, c resumeCase) strategyReasonOutcome {
	frame := strategySS + int(c.bp)
	state, finished := strategyReasonLoop(c.menuChoices, before[frame+4], strategyReasonOutcome{
		Budget: before[frame+3], Trust: before[int(cs)*16+0x0D00],
	})
	if !finished {
		panic("independent reason loop script has no terminal response")
	}
	return state
}

// The original scene explicitly selects carry and trust changes. The final
// AL is not inferred across its drawing/resource callees.
func strategyExpectedScene(before []byte, cs uint16, c resumeCase) (bool, byte) {
	trust := before[int(cs)*16+0x0D00]
	switch c.verdict {
	case 0, 3:
		if trust < 20 {
			panic("scene trust underflow is outside the returning-routine closure")
		}
		return true, strategyTrustSubtract(trust, 20)
	case 1:
		return false, strategyTrustAdd(trust, 20)
	case 2:
		state, finished := strategyReasonLoop(c.menuChoices, c.reasonMask, strategyReasonOutcome{
			Budget: strategyTrustBudget(trust), Trust: trust,
		})
		if !finished {
			panic("independent scene reason script has no terminal response")
		}
		if state.Carry {
			return true, state.Trust
		}
		return false, strategyTrustAdd(state.Trust, 10)
	default:
		return true, trust
	}
}

func strategyObserved(observed []registers, ip uint16) []registers {
	var result []registers
	for _, r := range observed {
		if r.IP == ip {
			result = append(result, r)
		}
	}
	return result
}

func strategyAudit(before, after []byte, cs uint16, c resumeCase, observed []registers) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent strategy audit group=%s target=%X: %s", c.group, c.target, detail))
		}
	}
	code := int(cs) * 16
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "date changed")
	if c.group == "proposal" || c.group == "power" {
		check(bytes.Equal(before[strategyWorld:strategyWorld+0x5240], after[strategyWorld:strategyWorld+0x5240]), "pure proposal/power changed world")
		check(bytes.Equal(before[strategyQueue:strategyQueue+0x400], after[strategyQueue:strategyQueue+0x400]), "pure proposal/power changed queue")
	}
	if c.group == "finance-values" {
		calls := strategyObserved(observed, 0x062F)
		check(len(calls) == 11, "finance numeric call count")
		player := strategyWorld + int(c.player)*64
		high := uint16(before[player+0x22])
		color := uint16(0x0907)
		if high&0x80 != 0 {
			high |= 0xFF00
			color = 0x0A07
		}
		want := []registers{
			{AX: strategyWord(before, player+0x20), DX: high, BX: color, DI: 0x230D},
			{AX: strategyWord(before, code+0x0D02), DX: uint16(before[code+0x0D04]), BX: 0x0906, DI: 0x1E24},
			{AX: strategyWord(before, code+0x0D05), DX: uint16(before[code+0x0D07]), BX: 0x0906, DI: 0x2324},
		}
		for index := 0; index < 8; index++ {
			value := strategyWord(before, code+0x0D08+index*2)
			if index%4 != 0 {
				value *= 10 // MUL high word is discarded by the original XOR DX.
			}
			position := uint16(0x320F + index%4*0x500)
			if index >= 4 {
				position = uint16(0x3223 + index%4*0x500)
			}
			want = append(want, registers{AX: value, DX: 0, BX: 0x0905, DI: position})
		}
		for index, r := range want {
			got := calls[index]
			check(got.AX == r.AX && got.DX == r.DX && got.BX == r.BX && got.DI == r.DI, fmt.Sprintf("finance numeric %d", index))
		}
	}
	if c.group == "general" {
		general := strategyWorld + 0x4240 + c.general*32
		message := uint16(0x1A8)
		if before[general+0x1D] == 0xFF {
			best := 0
			for index := 1; index < 3; index++ {
				if before[general+0x0E+index] > before[general+0x0E+best] {
					best = index
				}
			}
			message = uint16(0x1A9 + best)
		}
		count := 0
		for _, r := range strategyObserved(observed, 0x8810) {
			if r.CX >= 0x1A8 && r.CX <= 0x1AB {
				check(r.CX == message && byte(r.AX) == before[general+1] && byte(r.AX>>8) == before[general+0x1E] && r.DS == 0x7000, "general full-byte aptitude/tie/captor/portrait")
				count++
			}
		}
		if c.mode != 0 {
			check(count != 0, "general script produced no self-description")
		}
	}
	if c.group == "bar" {
		table := [...]uint16{0x6224, 0x6265, 0x678D, 0x6288, 0x628F, 0x62FB, 0x6366, 0x63BF}
		for index, handler := range table {
			check(strategyWord(before, code+0x6214+index*2) == handler, "original eight-handler table")
		}
		x := c.cx
		if x < 24 || x >= 408 {
			for _, handler := range table {
				check(len(strategyObserved(observed, handler)) == 0, "out-of-range command dispatched")
			}
		} else {
			index := int(x-24) / 48
			check(len(strategyObserved(observed, table[index])) != 0, "selected command handler missing")
			count := 0
			for _, r := range strategyObserved(observed, 0x0B46) {
				if r.DX == uint16(24+index*48) && r.BX == 40 && r.SI == 48 && r.DI == 16 {
					count++
				}
			}
			check(count == 2, "command highlight rectangle/save/restore pair")
		}
	}
	if c.group == "gate" {
		player := strategyWorld + int(c.player)*64
		lord := strategyWorld + 0x4240 + int(before[player+1])*32
		blocked := before[lord+0x17] != 0
		popups := strategyObserved(observed, 0x93E9)
		if blocked {
			check(len(popups) == 0, "nonzero ruler duty opened an advice popup")
			found := false
			for _, r := range strategyObserved(observed, 0x8810) {
				if r.CX == 0x40 && byte(r.AX) == 0x93 {
					found = true
				}
			}
			check(found, "nonzero ruler duty did not produce TALK 64")
		} else {
			check(len(popups) != 0, "free ruler did not open advice popup")
			r := popups[0]
			check(r.AX == 5 && r.CX == 0x4D && r.DX == 0x0400, "advice popup original arguments")
		}
	}
	if c.group == "entry" {
		popup := map[uint32]registers{
			0x1628F: {AX: 2, CX: 0x4F, DX: 0x040C},
			0x162FB: {AX: 2, CX: 0x52, DX: 0x040F},
		}
		if want, exists := popup[c.target]; exists {
			calls := strategyObserved(observed, 0x93E9)
			check(len(calls) != 0, "upper entry popup absent")
			r := calls[0]
			check(r.AX == want.AX && r.CX == want.CX && r.DX == want.DX, "upper popup count/text/position")
		}
		if c.target == 0x16288 {
			calls := strategyObserved(observed, 0x6C5E)
			check(len(calls) == 1 && calls[0].CX == 1, "formation entry CX")
		}
		for _, r := range strategyObserved(observed, 0x2151) {
			check(r.AX == 20 && r.CX == 12, "upper camera centering arguments")
			if c.target == 0x162FB || c.target == 0x163BF {
				check(r.SI%32 == 0 && r.SI < 192*32, "city card entry takes index times 32")
				city := strategyWorld + 0x840 + int(r.SI)
				check(r.DX == strategyWord(before, city+8) && r.BX == strategyWord(before, city+10), "camera consumes selected city coordinates")
				if c.target == 0x163BF {
					check(r.SI == uint16(before[strategyWorld+c.owner*64+3])*32, "faction entry focuses its original capital")
				}
			} else if c.target == 0x1628F {
				army := strategyWorld + 0x2240 + c.general*64
				check(r.DX == strategyWord(before, army+0x10) && r.BX == strategyWord(before, army+0x12), "location entry consumes selected corps coordinates")
			}
		}
	}
	if c.group == "trust" && c.target == 0x13D91 {
		want := strategyTrustAdd(before[code+0x0D00], byte(c.ax))
		check(after[code+0x0D00] == want, "trust saturating addition")
	}
	if c.group == "reason" || c.group == "reason-loop" {
		frame := strategySS + int(c.bp)
		state := strategyExpectedReason(before, cs, c)
		if c.group == "reason-loop" {
			initial := strategyReasonOutcome{Budget: before[frame+3], Trust: before[code+0x0D00]}
			var finished bool
			state, finished = strategyReasonLoop(c.menuChoices, before[frame+4], initial)
			check(finished, "reason script has no terminal response")
		}
		check(after[frame+5] == state.Attempted, "reason attempted mask")
		check(after[frame+3] == state.Budget, "reason byte budget including wrap")
		check(after[code+0x0D00] == state.Trust, "reason trust penalty")
		calls := strategyObserved(observed, 0x3C99)
		check(len(calls) != 0, "reason response display absent")
		check(calls[len(calls)-1].CX == strategyWord(before, frame)+state.TalkOffset, "reason response message offset")
	}
	if c.group == "scene" {
		calls := strategyObserved(observed, 0x3830)
		check(len(calls) == 1, "scene root entry count")
		check(byte(calls[0].AX) == c.verdict && byte(calls[0].DX) == c.reasonMask, "scene verdict/mask input")
		_, trust := strategyExpectedScene(before, cs, c)
		check(after[code+0x0D00] == trust, "scene explicit trust outcome")
	}
	if c.group == "queue-search" {
		check(bytes.Equal(before[strategyQueue:strategyQueue+0x400], after[strategyQueue:strategyQueue+0x400]), "queue search wrote queue")
	}
	if c.group == "queue-write" {
		want := append([]byte(nil), before[strategyQueue:strategyQueue+0x400]...)
		key := uint16(byte(c.ax)) | uint16(byte((c.si<<2)>>8))<<8
		start := int(strategyWord(before, code+0x0D20)) + int(byte(c.bx))*4
		for offset := start; offset < 0x400; offset += 4 {
			if want[offset] == 0 {
				strategyPutWord(want, offset, key)
				strategyPutWord(want, offset+2, c.dx)
				break
			}
		}
		check(bytes.Equal(want, after[strategyQueue:strategyQueue+0x400]), "queue writer first-free slot and raw words")
		check(strategyWord(after, code+0x0D20) == strategyWord(before, code+0x0D20), "queue writer changed head")
		check(strategyWord(after, code+0x0D22) == strategyWord(before, code+0x0D22), "queue writer invented a tail update")
	}
	if c.group == "sound" {
		priority := before[code+0x020F]
		if priority >= c.newPriority || c.soundFlags&4 == 0 {
			priority = c.newPriority
		}
		check(after[code+0x020F] == priority, "sound priority and explicit INT61 status fixture")
	}
	// Other UI/advice/scene effects need original intermediate-state evidence;
	// this helper does not infer world preservation or scene meaning from C.
	return checks
}
