//go:build matching_main

package main

import (
	"bytes"
	"fmt"

	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/machine"
)

type mainOutcome struct {
	Escape, Carry, CheckCarry, CheckResponse bool
	Trust, Attempted, Budget, Response       byte
	HasFrame                                 bool
	Frame                                    uint16
}

func mainPhysical(segment, offset uint16) int {
	return (int(segment)*16 + int(offset)) & 0xFFFFF
}

func mainPrepare(mem []byte, cs uint16, c resumeCase) {
	code := int(cs) * 16
	if c.returnIP != 0x006A || strategyWord(mem, code+0x9901) != c.savedSP ||
		strategyWord(mem, code+0x9903) != c.savedSS ||
		strategyWord(mem, mainPhysical(c.savedSS, c.savedSP)) != c.returnIP {
		panic("main fixture lacks original CALL10067/prefix-established outer frame")
	}
	prepared := c
	prepared.group = "main-fixture" // Avoid strategy's fixed-F000 reason frame.
	strategyPrepare(mem, cs, prepared)
	copy(mem[code+0x1964:code+0x1994], c.palette[:])
	if c.group == "reason" || c.group == "reason-loop" {
		frame := mainPhysical(c.entrySS, c.bp)
		strategyPutWord(mem, frame, 0x66)
		mem[frame+2], mem[frame+3] = c.verdict, c.budget
		mem[frame+4], mem[frame+5] = c.reasonMask, c.attempted
	}
	// The evidence-bearing saved stack is verified rather than synthesized.
	if strategyWord(mem, code+0x9901) != c.savedSP || strategyWord(mem, code+0x9903) != c.savedSS ||
		strategyWord(mem, mainPhysical(c.savedSS, c.savedSP)) != c.returnIP {
		panic("main preparation overwrote original saved outer frame")
	}
}

func mainReasonStep(reason, mask byte, result mainOutcome) mainOutcome {
	if reason < 4 {
		bit := byte(1 << reason)
		if result.Attempted&bit == 0 && mask&bit == 0 && result.Trust < 20 {
			// 13BC9 stores the attempt bit before 13BDA calls 13DC9. Borrow
			// forces CX=19E and escapes before 13BDF can assign response AH=0.
			result.Attempted |= bit
			result.Trust = 0
			result.Escape = true
			result.CheckCarry, result.CheckResponse = false, false
			return result
		}
	}
	state := strategyReasonStep(reason, mask, strategyReasonOutcome{
		Attempted: result.Attempted, Budget: result.Budget, Trust: result.Trust,
	})
	result.Attempted, result.Budget, result.Trust = state.Attempted, state.Budget, state.Trust
	result.Response, result.Carry = state.Response, state.Carry
	result.CheckCarry, result.CheckResponse = true, true
	return result
}

func mainReasonLoop(choices []int, mask byte, result mainOutcome) mainOutcome {
	result.Attempted = 0
	for _, choice := range choices {
		if choice < 0 {
			continue // Original 13B7E retries its own canceled selector.
		}
		result = mainReasonStep(byte(choice), mask, result)
		if result.Escape {
			return result
		}
		if result.Carry {
			result.Carry = result.Response != 1
			return result
		}
	}
	panic("independent main reason script has no terminal response")
}

func mainExpected(before []byte, cs uint16, c resumeCase) mainOutcome {
	result := mainOutcome{Trust: before[int(cs)*16+0x0D00]}
	switch c.group {
	case "palette", "fade":
		return result
	case "escape":
		result.Escape = true
		return result
	case "trust":
		switch c.target {
		case 0x13DC9:
			amount := byte(c.ax)
			borrow := result.Trust < amount
			result.Trust = strategyTrustSubtract(result.Trust, amount)
			result.Escape = borrow || c.cx != 0xFFFF && result.Trust == 0
		case 0x13D91:
			result.Trust = strategyTrustAdd(result.Trust, byte(c.ax))
		default:
			panic("independent main trust target")
		}
		return result // No independent final CF contract for these helpers.
	case "reason", "reason-loop":
		frame := mainPhysical(c.entrySS, c.bp)
		result.HasFrame, result.Frame = true, c.bp
		result.Attempted, result.Budget = before[frame+5], before[frame+3]
		if c.group == "reason" {
			return mainReasonStep(byte(c.ax), before[frame+4], result)
		}
		return mainReasonLoop(c.menuChoices, before[frame+4], result)
	case "scene":
		switch byte(c.ax) {
		case 0, 3:
			result.Escape = result.Trust < 20
			result.Trust = strategyTrustSubtract(result.Trust, 20)
			result.Carry, result.CheckCarry = true, !result.Escape
		case 1:
			result.Trust = strategyTrustAdd(result.Trust, 20)
			result.Carry, result.CheckCarry = false, true
		case 2:
			result.Frame = c.entrySP - 12
			result.Budget = strategyTrustBudget(result.Trust)
			result = mainReasonLoop(c.menuChoices, byte(c.dx), result)
			// Normal scene cleanup frees these six bytes at 138BC. The CALL
			// at 138BF overwrites +4/+5, and 19321 PUSH AX overwrites +2/+3.
			// Only a nonlocal escape bypasses that cleanup and retains them.
			result.HasFrame = result.Escape
			result.CheckResponse = false // Scene does not return reason AH.
			if !result.Escape && !result.Carry {
				result.Trust = strategyTrustAdd(result.Trust, 10)
			}
		default:
			result.Carry, result.CheckCarry = true, true
		}
		return result
	default:
		panic("unknown independent main group")
	}
}

func mainDACValue(component, level byte) byte {
	return byte(4 * ((uint16(component&15)*uint16(level) + 8) / 16))
}

func mainAudit(before, after []byte, cs uint16, c resumeCase, observed []registers, device *machine.Machine) int {
	checks := 0
	check := func(ok bool, detail string) {
		checks++
		if !ok {
			panic(fmt.Sprintf("independent main audit group=%s target=%X: %s", c.group, c.target, detail))
		}
	}
	code := int(cs) * 16
	result := mainExpected(before, cs, c)
	check(strategyWord(before, code+0x9901) == c.savedSP && strategyWord(before, code+0x9903) == c.savedSS,
		"original outer frame input")
	check(strategyWord(after, code+0x9901) == c.savedSP && strategyWord(after, code+0x9903) == c.savedSS,
		"saved SS/SP changed")
	check(strategyWord(before, mainPhysical(c.savedSS, c.savedSP)) == c.returnIP &&
		strategyWord(after, mainPhysical(c.savedSS, c.savedSP)) == c.returnIP, "outer return word")
	check(after[code+0x0D00] == result.Trust, "trust outcome")
	check(bytes.Equal(before[code+0x0CF0:code+0x0CF8], after[code+0x0CF0:code+0x0CF8]), "date changed")
	if result.HasFrame {
		frame := mainPhysical(c.entrySS, result.Frame)
		check(after[frame+5] == result.Attempted, "attempt bit before possible nonlocal escape")
		check(after[frame+3] == result.Budget, "reason budget")
	}
	if result.Escape {
		check(device.CPU.Seg[cpu.SS] == c.savedSS && device.CPU.R[cpu.SP] == c.savedSP+2 &&
			device.CPU.Seg[cpu.CS] == cs && device.CPU.IP == c.returnIP, "nonlocal saved-frame return")
	} else {
		check(device.CPU.Seg[cpu.SS] == c.entrySS && device.CPU.R[cpu.SP] == c.entrySP+2 &&
			device.CPU.Seg[cpu.CS] == cs && device.CPU.IP == 0xF000, "normal caller return")
	}
	if result.CheckCarry {
		check((device.CPU.Flags&cpu.CF != 0) == result.Carry, "normal caller carry")
	}
	if result.CheckResponse {
		check(byte(device.CPU.R[cpu.AX]>>8) == result.Response, "normal reason response AH")
	}
	if c.group == "scene" {
		var minus, plus byte
		switch byte(c.ax) {
		case 0, 3:
			minus = 20
		case 1:
			plus = 20
		case 2:
			if result.Escape || result.Response == 0 {
				minus = 20
			} else if !result.Carry {
				plus = 10
			}
		}
		for _, expected := range []struct {
			ip     uint16
			amount byte
		}{{0x3DC9, minus}, {0x3D91, plus}} {
			calls := strategyObserved(observed, expected.ip)
			count := 0
			if expected.amount != 0 {
				count = 1
			}
			check(len(calls) == count, fmt.Sprintf("scene trust decision calls %04X", expected.ip))
			if count != 0 {
				check(byte(calls[0].AX) == expected.amount && calls[0].CX == 0xFFFF, "scene trust decision arguments")
			}
		}
		for _, ip := range []uint16{0x3D45, 0x1D46, 0x9321} {
			count := 1
			if result.Escape {
				count = 0
			}
			check(len(strategyObserved(observed, ip)) == count, fmt.Sprintf("scene suffix lifetime %04X", ip))
		}
	}
	escapes := strategyObserved(observed, 0x1CB1)
	if result.Escape {
		check(len(escapes) == 1, "nonlocal exit entry count")
		seenExit := false
		for _, r := range observed {
			if r.CS != cs {
				continue // Existing mouse far-service entries use another CS.
			}
			if r.IP == 0x1CB1 {
				seenExit = true
				continue
			}
			if seenExit {
				check(r.IP == 0x0A1C || r.IP == 0xEBDC, "caller suffix resumed after nonlocal exit")
			}
		}
	} else {
		check(len(escapes) == 0, "normal caller entered nonlocal exit")
	}
	fade := result.Escape || c.group == "fade"
	setters := strategyObserved(observed, 0xEBDC)
	wantDAC := mainBeforeDAC
	type dacWrite struct {
		port  uint16
		value byte
	}
	var wantWrites []dacWrite
	appendSetter := func(index, level, raw0, raw1, raw2 byte) {
		values := [...]byte{mainDACValue(raw1, level), mainDACValue(raw2, level), mainDACValue(raw0, level)}
		wantWrites = append(wantWrites, dacWrite{0x3C8, index})
		for channel, value := range values {
			wantWrites = append(wantWrites, dacWrite{0x3C9, value})
			wantDAC[int(index)*3+channel] = value & 63
		}
	}
	if fade {
		check(len(strategyObserved(observed, 0x0A1C)) == 1, "real fade entry count")
		check(len(setters) == 17*16, "fade must execute all 272 DAC setter calls")
		call := 0
		for level := 16; level >= 0; level-- {
			for index := 0; index < 16; index++ {
				at := code + 0x1964 + index*3
				raw0, raw1, raw2 := before[at], before[at+1], before[at+2]
				r := setters[call]
				check(byte(r.AX) == byte(index) && byte(r.BX) == byte(level) && byte(r.BX>>8) == raw0 &&
					byte(r.DX) == raw1 && byte(r.DX>>8) == raw2 && r.CX == uint16(index<<8|level) &&
					r.SI == uint16(0x1964+index*3), fmt.Sprintf("fade raw arguments level=%d color=%d", level, index))
				appendSetter(byte(index), byte(level), raw0, raw1, raw2)
				call++
			}
		}
		check(bytes.Equal(wantDAC[:48], make([]byte, 48)), "fade modeled first sixteen colors are black")
	} else if c.group == "palette" {
		check(len(setters) == 1, "direct DAC setter entry count")
		r := setters[0]
		check(byte(r.AX) == byte(c.ax) && r.BX == c.bx && r.DX == c.dx, "direct setter raw inputs")
		appendSetter(byte(c.ax), byte(c.bx), byte(c.bx>>8), byte(c.dx), byte(c.dx>>8))
	} else {
		check(len(setters) == 0, "normal caller unexpectedly faded the DAC")
	}
	var gotWrites []dacWrite
	for _, w := range device.PortLog {
		if w.Port == 0x3C8 || w.Port == 0x3C9 {
			gotWrites = append(gotWrites, dacWrite{w.Port, w.Val})
		}
	}
	check(len(gotWrites) == len(wantWrites), "complete DAC write stream length")
	for index, want := range wantWrites {
		check(gotWrites[index] == want, fmt.Sprintf("DAC stream write %d", index))
	}
	check(bytes.Equal(device.DAC[:], wantDAC[:]), "full DAC including untouched colors")
	if fade {
		check(bytes.Equal(device.DAC[:48], make([]byte, 48)), "hardware DAC first sixteen colors are black")
		check(bytes.Equal(device.DAC[48:], mainBeforeDAC[48:]), "hardware DAC colors >=16 unchanged")
	}
	return checks
}
