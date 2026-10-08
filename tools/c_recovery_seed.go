//go:build ignore

// 固定 RTC 回覆的 sub_1EC82 原版/C/Go 播種對照。只在 Docker 內執行。
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/seed.c"
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
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"os"
	"path/filepath"
	"unsafe"
)

func seedSetRegs(m *C.KiMachine16, r oracle.Regs) {
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
func seedRegs(m *C.KiMachine16) oracle.Regs {
	return oracle.Regs{AX: uint16(m.ax), BX: uint16(m.bx), CX: uint16(m.cx), DX: uint16(m.dx), SI: uint16(m.si),
		DI: uint16(m.di), BP: uint16(m.bp), SP: uint16(m.sp), DS: uint16(m.ds), ES: uint16(m.es),
		SS: uint16(m.ss), CS: uint16(m.cs), IP: uint16(m.ip), Flags: uint16(m.flags)}
}
func seedBCD(v int) byte { return byte(v/10<<4 | v%10) }

func main() {
	exe := flag.String("exe", "/orig/KI.EXE", "原版")
	out := flag.String("out", "/output/seed.json", "收據")
	max := flag.Int("max", 86400, "正常合法時間數")
	flag.Parse()
	if *max < 1 || *max > 86400 {
		panic("max must be 1..86400")
	}
	original, err := os.ReadFile(*exe)
	if err != nil {
		panic(err)
	}
	sha := sha256.Sum256(original)
	expected := "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868"
	if hex.EncodeToString(sha[:]) != expected {
		panic("wrong original")
	}
	o, err := oracle.Load(*exe, filepath.Dir(*exe))
	if err != nil {
		panic(err)
	}
	defer o.Close()
	cs := o.Regs().CS
	entry := o.IDAIn(cs, 0x1ec82)
	code := o.Bytes(entry, 94)
	if !bytes.Equal(code, original[0xee82:0xeee0]) {
		panic("seed code differs")
	}
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	memory := unsafe.Slice((*byte)(native), 1<<20)
	copy(memory, o.Bytes(oracle.Phys(0), 1<<20))
	both := func(a oracle.Addr, b []byte) { o.WriteBytes(a, b); copy(memory[a.Linear():], b) }
	word := func(seg, off, value uint16) {
		var b [2]byte
		binary.LittleEndian.PutUint16(b[:], value)
		both(oracle.Addr{Seg: seg, Off: off}, b[:])
	}
	stateAddr := oracle.Addr{Seg: cs, Off: 0xecfc}
	stateAt := stateAddr.Linear()
	both(oracle.Phys(0x1a*4), []byte{0x00, 0x01, 0x00, 0x70}) // 明示 RTC fixture ISR。
	originalHash, cHash, goHash := sha256.New(), sha256.New(), sha256.New()
	cases, goCases, audits := 0, 0, 0
	var failure map[string]any
	var active bool
	var stackSP uint16
	var unexpected []oracle.WriteHit
	o.OnWrite(0, 0xfffff, func(_ *oracle.Oracle, h oracle.WriteHit) {
		if !active {
			return
		}
		allowed := h.Addr >= stateAt && h.Addr < stateAt+258
		for n := uint16(0); n < 18; n++ {
			if h.Addr == uint32(0x4000)*16+uint32(stackSP-16+n) {
				allowed = true
			}
		}
		if !allowed && len(unexpected) < 8 {
			unexpected = append(unexpected, h)
		}
	})
	stepCounts := map[uint64]int{}
	check := func(hour, minute, second int, al byte, sp uint16, index int) bool {
		h, mi, s := seedBCD(hour), seedBCD(minute), seedBCD(second)
		pre := o.Regs()
		word(pre.SS, pre.SP-2, 0xbeef)
		inputFlags := uint16(2)
		if index&1 != 0 {
			inputFlags |= 0x400
		}
		setup := []byte{0xb8, 0x00, 0x40, 0x8e, 0xd0, 0xbc, byte(sp + 2), byte((sp + 2) >> 8), 0x68, byte(inputFlags), byte(inputFlags >> 8), 0x9d}
		both(oracle.Addr{Seg: 0x7000, Off: 0}, setup)
		o.CallNear(oracle.Addr{Seg: 0x7000, Off: 0}, 0xbeef, oracle.CallRegs{})
		if err := o.Run(5); err != nil {
			panic(err)
		}
		word(0x4000, sp, inputFlags)
		isr := []byte{0xb9, mi, h, 0xba, 0, s, 0xb0, al, 0xcf}
		both(oracle.Addr{Seg: 0x7000, Off: 0x100}, isr)
		initial := make([]byte, 258)
		for i := range initial {
			initial[i] = byte(i*23 + index)
		}
		both(stateAddr, initial)
		word(0x4000, sp, 0x0100)
		r := oracle.CallRegs{AX: uint16(index*19) ^ 0x1234, BX: uint16(index*23) ^ 0x2345, CX: uint16(index*29) ^ 0x3456,
			DX: uint16(index*31) ^ 0x4567, SI: uint16(index*37) ^ 0x5678, DI: uint16(index*41) ^ 0x6789,
			BP: uint16(index*43) ^ 0x789a, DS: uint16(index*47) ^ 0x89ab, ES: uint16(index*53) ^ 0x9abc,
			SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		o.CallNear(entry, 0x0100, r)
		input := o.Regs()
		seedSetRegs(&m, input)
		stackSP = input.SP
		active = true
		start := o.Steps()
		err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x0100}), oracle.Budget(4000))
		active = false
		if err != nil {
			panic(err)
		}
		steps := o.Steps() - start
		stepCounts[steps]++
		C.sub_1EC82_abi(&m, C.uint8_t(h), C.uint8_t(mi), C.uint8_t(s), C.uint8_t(al))
		got, want := seedRegs(&m), o.Regs()
		state := o.Bytes(stateAddr, 258)
		cases++
		if got != want || !bytes.Equal(memory[stateAt:stateAt+258], state) || len(unexpected) > 0 {
			failure = map[string]any{"case": index, "hour": hour, "minute": minute, "second": second, "rtc_al": al, "input": input, "original": want, "c": got, "original_state": hex.EncodeToString(state), "c_state": hex.EncodeToString(memory[stateAt : stateAt+258]), "unexpected": unexpected}
			return false
		}
		for n := uint16(0); n < 18; n++ {
			at := uint32(0x4000)*16 + uint32(sp-16+n)
			if memory[at] != o.Byte(oracle.Phys(at)) {
				failure = map[string]any{"case": index, "stack_at": at, "c": memory[at], "original": o.Byte(oracle.Phys(at))}
				return false
			}
		}
		originalHash.Write(state)
		cHash.Write(memory[stateAt : stateAt+258])
		if al == 0 {
			goState := rng.New(hour, minute, second).Raw()
			goCases++
			goHash.Write(goState)
			if !bytes.Equal(goState, state) {
				failure = map[string]any{"case": index, "go_state": hex.EncodeToString(goState), "original_state": hex.EncodeToString(state)}
				return false
			}
		}
		if cases%1024 == 0 {
			if !bytes.Equal(memory, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory differs")
			}
			audits++
		}
		return true
	}
	for i := 0; i < *max; i++ {
		if !check(i/3600, (i/60)%60, i%60, 0, 0x8000, i) {
			goto finish
		}
	}
	for n, al := range []byte{1, 0x7f, 0x80, 0xff} {
		for _, second := range []int{0, 27, 59} {
			if !check(23, 59, second, al, 0x8000, 86400+n*3+second) {
				goto finish
			}
		}
	}
	for n, sp := range []uint16{0, 2, 0xfffc, 0xfffe} {
		for _, second := range []int{0, 59} {
			if !check(23, 59, second, 0, sp, 86500+n*2+second) {
				goto finish
			}
		}
	}
finish:
	if failure == nil {
		if !bytes.Equal(memory, o.Bytes(oracle.Phys(0), 1<<20)) {
			panic("final memory differs")
		}
		audits++
	}
	if !bytes.Equal(o.Bytes(entry, 94), code) {
		panic("original code changed")
	}
	codeSha := sha256.Sum256(code)
	report := map[string]any{"schema": "wolong-c-seed-parity-v1", "input_sha256": expected, "routine_sha256": hex.EncodeToString(codeSha[:]),
		"passed": failure == nil, "cases": cases, "go_cases": goCases, "full_memory_audits": audits, "steps": stepCounts, "mismatch": failure,
		"original_state_sha256": hex.EncodeToString(originalHash.Sum(nil)), "c_state_sha256": hex.EncodeToString(cHash.Sum(nil)), "go_state_sha256": hex.EncodeToString(goHash.Sum(nil)),
		"rtc_fixture":          "IVT 1Ah -> 7000:0100, MOV CX, MOV DX, MOV AL, IRET; fixed before each call",
		"scope":                "legal BCD, IF/TF=0, disjoint stack, final 16-bit CPU/memory; XOR AF follows dosgolem model; no real RTC timing or player-path claim",
		"c_machine_code_match": false}
	data, _ := json.MarshalIndent(report, "", "  ")
	if err := os.WriteFile(*out, append(data, '\n'), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("seed C / original: %d; Go: %d; audits: %d; pass=%v\n", cases, goCases, audits, failure == nil)
	if failure != nil {
		os.Exit(1)
	}
}
