//go:build ignore

// 原始 sub_11D8E、C 時鐘與 Go 日期規則的局部對照；Docker only。
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/clock.c"
#include "/repo/tools/c_recovery/clock_fixture.h"
*/
import "C"
import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wicanr2/dosgolem/oracle"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"os"
	"path/filepath"
	"unsafe"
)

type ClockTrace struct {
	Target uint16
	Regs   oracle.Regs
	Date   string
}

func clockSet(m *C.KiMachine16, r oracle.Regs) {
	m.ax = C.uint16_t(r.AX)
	m.bx = C.uint16_t(r.BX)
	m.cx = C.uint16_t(r.CX)
	m.dx = C.uint16_t(r.DX)
	m.si = C.uint16_t(r.SI)
	m.di = C.uint16_t(r.DI)
	m.bp = C.uint16_t(r.BP)
	m.sp = C.uint16_t(r.SP)
	m.ds = C.uint16_t(r.DS)
	m.es = C.uint16_t(r.ES)
	m.ss = C.uint16_t(r.SS)
	m.cs = C.uint16_t(r.CS)
	m.ip = C.uint16_t(r.IP)
	m.flags = C.uint16_t(r.Flags)
}
func clockRegs(m *C.KiMachine16) oracle.Regs {
	return oracle.Regs{AX: uint16(m.ax), BX: uint16(m.bx), CX: uint16(m.cx), DX: uint16(m.dx), SI: uint16(m.si), DI: uint16(m.di), BP: uint16(m.bp), SP: uint16(m.sp), DS: uint16(m.ds), ES: uint16(m.es), SS: uint16(m.ss), CS: uint16(m.cs), IP: uint16(m.ip), Flags: uint16(m.flags)}
}
func cClockTrace(s C.KiClockSnapshot) ClockTrace {
	return ClockTrace{Target: uint16(s.target), Regs: oracle.Regs{AX: uint16(s.ax), BX: uint16(s.bx), CX: uint16(s.cx), DX: uint16(s.dx), SI: uint16(s.si), DI: uint16(s.di), BP: uint16(s.bp), SP: uint16(s.sp), DS: uint16(s.ds), ES: uint16(s.es), SS: uint16(s.ss), CS: uint16(s.cs), IP: uint16(s.ip), Flags: uint16(s.flags)}, Date: hex.EncodeToString(C.GoBytes(unsafe.Pointer(&s.date[0]), 8))}
}
func main() {
	exe := flag.String("exe", "/orig/KI.EXE", "original")
	out := flag.String("out", "/output/clock.json", "receipt")
	max := flag.Int("max", 65536, "year matrix")
	flag.Parse()
	if *max < 1 || *max > 65536 {
		panic("max")
	}
	original, e := os.ReadFile(*exe)
	if e != nil {
		panic(e)
	}
	digest := sha256.Sum256(original)
	expected := "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868"
	if hex.EncodeToString(digest[:]) != expected {
		panic("original identity")
	}
	o, e := oracle.Load(*exe, filepath.Dir(*exe))
	if e != nil {
		panic(e)
	}
	defer o.Close()
	cs := o.Regs().CS
	entry := o.IDAIn(cs, 0x11d8e)
	code := o.Bytes(entry, 137)
	if !bytes.Equal(code, original[0x1f8e:0x2017]) {
		panic("routine identity")
	}
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	memory := unsafe.Slice((*byte)(native), 1<<20)
	copy(memory, o.Bytes(oracle.Phys(0), 1<<20))
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	both := func(a oracle.Addr, b []byte) { o.WriteBytes(a, b); copy(memory[a.Linear():], b) }
	word := func(seg, off, value uint16) {
		var b [2]byte
		binary.LittleEndian.PutUint16(b[:], value)
		both(oracle.Addr{Seg: seg, Off: off}, b[:])
	}
	stateAddr := oracle.Addr{Seg: cs, Off: 0xcf0}
	stateAt := stateAddr.Linear()
	months := o.Bytes(oracle.Addr{Seg: cs, Off: 0x98ab}, 13)
	for n := 1; n <= 12; n++ {
		if int(months[n]) != clock.DaysInMonth(n) {
			panic("month table")
		}
	}
	var trace []ClockTrace
	var readyPolls, countPolls, readyDelay, countDelay int
	var delayed bool
	var targetCount byte
	for _, target := range []uint16{0x5358, 0x9377, 0x3e11, 0x1e17} {
		t := target
		o.OnCall(oracle.Addr{Seg: cs, Off: t}, func(o *oracle.Oracle) {
			trace = append(trace, ClockTrace{Target: t, Regs: o.Regs(), Date: hex.EncodeToString(o.Bytes(stateAddr, 8))})
		})
	}
	o.OnCall(oracle.Addr{Seg: cs, Off: 0x1dff}, func(o *oracle.Oracle) {
		readyPolls++
		if delayed && readyPolls > readyDelay {
			o.WriteU8(oracle.Addr{Seg: cs, Off: 0xd2c}, 1)
		}
	})
	o.OnCall(oracle.Addr{Seg: cs, Off: 0x1e06}, func(o *oracle.Oracle) {
		countPolls++
		if delayed && countPolls > countDelay {
			o.WriteU8(oracle.Addr{Seg: cs, Off: 0xd2d}, targetCount)
		}
		if countPolls > 64 {
			o.WriteU8(oracle.Addr{Seg: cs, Off: 0xd2d}, targetCount+1)
		}
	})
	hashOriginal, hashC := sha256.New(), sha256.New()
	cases, goCases, audits := 0, 0, 0
	groups := map[string]int{}
	var failure map[string]any
	check := func(group string, year uint16, month, day, hour, sub, speed byte, alter, wait bool, index int) bool {
		trace = nil
		readyPolls = 0
		countPolls = 0
		delayed = wait
		readyDelay = 3
		countDelay = 5
		targetCount = speed
		if alter {
			targetCount = 2
		}
		pre := o.Regs()
		word(pre.SS, pre.SP-2, 0xbeef)
		ss, sp := uint16(0x4000), uint16(0x8000)
		flags := uint16(2)
		if index&1 != 0 {
			flags |= 0x400
		}
		setup := []byte{0xb8, 0, 0x40, 0x8e, 0xd0, 0xbc, 2, 0x80, 0x68, byte(flags), byte(flags >> 8), 0x9d}
		both(oracle.Addr{Seg: 0x7000, Off: 0}, setup)
		o.CallNear(oracle.Addr{Seg: 0x7000, Off: 0}, 0xbeef, oracle.CallRegs{})
		if e := o.Run(5); e != nil {
			panic(e)
		}
		word(ss, sp, flags)
		date := []byte{day, months[month], sub, hour, month, 0, byte(year), byte(year >> 8)}
		both(stateAddr, date)
		both(oracle.Addr{Seg: cs, Off: 0xcfa}, []byte{speed})
		both(oracle.Addr{Seg: cs, Off: 0xd2c}, []byte{1, targetCount})
		if wait {
			both(oracle.Addr{Seg: cs, Off: 0xd2c}, []byte{0, 0})
		}
		for _, t := range []uint16{0x5358, 0x9377, 0x3e11, 0x1e17} {
			stub := []byte{0xc3}
			if alter {
				switch t {
				case 0x5358:
					stub = []byte{0xc6, 0x06, 0xf0, 0x0c, 7, 0xc3}
				case 0x9377:
					stub = []byte{0xc6, 0x06, 0xfa, 0x0c, 2, 0xc3}
				case 0x3e11:
					stub = []byte{0xb8, 0xc3, 0xa5, 0xc3}
				case 0x1e17:
					stub = []byte{0xbb, 0xc0, 0xb7, 0xc3}
				}
			}
			both(oracle.Addr{Seg: cs, Off: t}, stub)
		}
		word(ss, sp, 0x0100)
		r := oracle.CallRegs{AX: uint16(index*19) ^ 0x1234, BX: uint16(index*23) ^ 0x2345, CX: uint16(index*29) ^ 0x3456, DX: uint16(index*31) ^ 0x4567, SI: uint16(index*37) ^ 0x5678, DI: uint16(index*41) ^ 0x6789, BP: uint16(index*43) ^ 0x789a, DS: cs, ES: uint16(index*47) ^ 0x89ab, SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		o.CallNear(entry, 0x0100, r)
		input := o.Regs()
		clockSet(&m, input)
		fixture := C.KiClockFixture{}
		if alter {
			fixture.alter = 1
		}
		if wait {
			fixture.delayed = 1
		}
		fixture.ready_delay = 3
		fixture.count_delay = 5
		fixture.target_count = C.uint8_t(targetCount)
		if e := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x0100}), oracle.Budget(512)); e != nil {
			panic(e)
		}
		C.clock_run_fixture(&m, &fixture)
		got, want := clockRegs(&m), o.Regs()
		result := o.Bytes(stateAddr, 8)
		cases++
		groups[group]++
		ctrace := []ClockTrace{}
		for n := 0; n < int(fixture.calls) && n < 4; n++ {
			ctrace = append(ctrace, cClockTrace(fixture.trace[n]))
		}
		ot, _ := json.Marshal(trace)
		ct, _ := json.Marshal(ctrace)
		if len(trace) == 0 {
			ot = []byte("[]")
		}
		if got != want || !bytes.Equal(memory[stateAt:stateAt+8], result) || !bytes.Equal(ot, ct) || readyPolls != int(fixture.ready_polls) || countPolls != int(fixture.count_polls) {
			failure = map[string]any{"group": group, "case": index, "input": input, "initial_date": hex.EncodeToString(date), "original": want, "c": got, "original_date": hex.EncodeToString(result), "c_date": hex.EncodeToString(memory[stateAt : stateAt+8]), "original_trace": trace, "c_trace": ctrace, "original_polls": []int{readyPolls, countPolls}, "c_polls": []int{int(fixture.ready_polls), int(fixture.count_polls)}}
			return false
		}
		for _, off := range []uint16{0xd2c, 0xd2d} {
			at := oracle.Addr{Seg: cs, Off: off}
			if memory[at.Linear()] != o.Byte(at) {
				panic("wait bytes")
			}
		}
		for n := uint16(0); n < 4; n++ {
			at := oracle.Addr{Seg: ss, Off: sp - 2 + n}
			if memory[at.Linear()] != o.Byte(at) {
				panic("stack bytes")
			}
		}
		hashOriginal.Write(result)
		hashOriginal.Write(ot)
		hashC.Write(memory[stateAt : stateAt+8])
		hashC.Write(ct)
		if !alter {
			g := clock.Clock{Year: int(year), Month: int(month), Day: int(day), Hour: int(hour), Subtick: int(sub)}
			ev := g.Advance()
			goCases++
			goResult := []byte{byte(g.Day), byte(clock.DaysInMonth(g.Month)), byte(g.Subtick), byte(g.Hour), byte(g.Month), 0, byte(g.Year), byte(g.Year >> 8)}
			if !bytes.Equal(goResult, result) || ev.Hour != (len(trace) > 0) || ev.Month != (len(trace) == 4) {
				failure = map[string]any{"group": group, "case": index, "go": g, "original_date": hex.EncodeToString(result), "go_date": hex.EncodeToString(goResult), "go_event": ev}
				return false
			}
		}
		if cases%1024 == 0 {
			if !bytes.Equal(memory, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory")
			}
			audits++
		}
		return true
	}
	for _, year := range []uint16{196, 999, 1000} {
		for month := byte(1); month <= 12; month++ {
			for day := byte(1); day <= months[month]; day++ {
				for hour := byte(1); hour <= 23; hour++ {
					for sub := byte(0); sub <= 8; sub++ {
						if !check("legal-date", year, month, day, hour, sub, 0, false, false, cases) {
							goto finish
						}
					}
				}
			}
		}
	}
	for y := 0; y < *max; y++ {
		if !check("year-boundary", uint16(y), 12, 31, 23, 8, 0, false, false, cases) {
			goto finish
		}
	}
	for speed := byte(1); speed <= 8; speed++ {
		for _, month := range []byte{1, 2, 12} {
			for _, wait := range []bool{false, true} {
				for _, alter := range []bool{false, true} {
					if !check("wait-and-callback", 999, month, months[month], 23, 8, speed, alter, wait, cases) {
						goto finish
					}
				}
			}
		}
	}
finish:
	if failure == nil {
		if !bytes.Equal(memory, o.Bytes(oracle.Phys(0), 1<<20)) {
			panic("final memory")
		}
		audits++
	}
	if !bytes.Equal(o.Bytes(entry, 137), code) {
		panic("code changed")
	}
	codeSha := sha256.Sum256(code)
	report := map[string]any{"schema": "wolong-c-clock-parity-v1", "input_sha256": expected, "routine_sha256": hex.EncodeToString(codeSha[:]), "passed": failure == nil, "cases": cases, "go_cases": goCases, "groups": groups, "full_memory_audits": audits, "mismatch": failure, "original_trace_sha256": hex.EncodeToString(hashOriginal.Sum(nil)), "c_trace_sha256": hex.EncodeToString(hashC.Sum(nil)), "scope": "controlled near-call, legal dates and u16 year boundary, explicit callee/poll fixtures, IF/TF=0; no full callee or hardware timing claim", "c_machine_code_match": false}
	b, _ := json.MarshalIndent(report, "", "  ")
	if e := os.WriteFile(*out, append(b, '\n'), 0644); e != nil {
		panic(e)
	}
	fmt.Printf("clock C / original %d; Go %d; audits %d; pass=%v\n", cases, goCases, audits, failure == nil)
	if failure != nil {
		os.Exit(1)
	}
}
