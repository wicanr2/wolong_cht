//go:build ignore

// 原版 sub_1ECE0 → C → Go 的受控局部對照。所有執行均由 Docker 包裝器啟動。
package main

/*
#include <stdlib.h>
#include <string.h>
#include "/repo/tools/c_recovery/rng.c"
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
	"hash"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/wicanr2/dosgolem/oracle"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

const originalSHA = "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868"

func regs(m *C.KiMachine16) oracle.Regs {
	return oracle.Regs{AX: uint16(m.ax), BX: uint16(m.bx), CX: uint16(m.cx), DX: uint16(m.dx),
		SI: uint16(m.si), DI: uint16(m.di), BP: uint16(m.bp), SP: uint16(m.sp), DS: uint16(m.ds),
		ES: uint16(m.es), SS: uint16(m.ss), CS: uint16(m.cs), IP: uint16(m.ip), Flags: uint16(m.flags)}
}

func setRegs(m *C.KiMachine16, r oracle.Regs) {
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

func digestCase(h hash.Hash, r oracle.Regs, raw []byte) {
	fields := []uint16{r.AX, r.BX, r.CX, r.DX, r.SI, r.DI, r.BP, r.SP, r.DS, r.ES, r.SS, r.CS, r.IP, r.Flags}
	var buffer [28]byte
	for i, v := range fields {
		binary.LittleEndian.PutUint16(buffer[i*2:], v)
	}
	h.Write(buffer[:])
	h.Write(raw)
}

func main() {
	exe := flag.String("exe", "/orig/KI.EXE", "原始 EXE")
	out := flag.String("out", "/output/report.json", "收據")
	max := flag.Int("max", 65536, "每張表的枚舉數")
	flag.Parse()
	if *max < 1 || *max > 65536 {
		panic("max must be 1..65536")
	}
	rawExe, err := os.ReadFile(*exe)
	if err != nil {
		panic(err)
	}
	digest := sha256.Sum256(rawExe)
	if hex.EncodeToString(digest[:]) != originalSHA {
		panic("original identity differs")
	}
	o, err := oracle.Load(*exe, filepath.Dir(*exe))
	if err != nil {
		panic(err)
	}
	defer o.Close()
	cs := o.Regs().CS
	entry := o.IDAIn(cs, 0x1ece0)
	code := o.Bytes(entry, 28)
	if !bytes.Equal(code, rawExe[0xeee0:0xeefc]) {
		panic("original routine bytes differ")
	}
	codeHash := sha256.Sum256(code)
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	cMemory := unsafe.Slice((*byte)(native), 1<<20)
	copy(cMemory, o.Bytes(oracle.Phys(0), 1<<20))
	phys := func(seg, off uint16) uint32 { return (uint32(seg)*16 + uint32(off)) & 0xfffff }
	writeBoth := func(at oracle.Addr, data []byte) { o.WriteBytes(at, data); copy(cMemory[at.Linear():], data) }
	writeWord := func(seg, off, value uint16) {
		var b [2]byte
		binary.LittleEndian.PutUint16(b[:], value)
		writeBoth(oracle.Addr{Seg: seg, Off: off}, b[:])
	}
	stateAddr := oracle.Addr{Seg: cs, Off: 0xecfc}
	stateLinear := stateAddr.Linear()
	var active bool
	var stackLo uint16
	var stackSeg uint16
	var unexpected []oracle.WriteHit
	writes := uint64(0)
	o.OnWrite(0, 0xfffff, func(_ *oracle.Oracle, hit oracle.WriteHit) {
		if !active {
			return
		}
		writes++
		allowed := hit.Addr == stateLinear || hit.Addr == stateLinear+1
		for n := uint16(0); n < 6; n++ {
			if hit.Addr == phys(stackSeg, stackLo+n) {
				allowed = true
			}
		}
		if !allowed && len(unexpected) < 8 {
			unexpected = append(unexpected, hit)
		}
	})
	originalHash, cHash, goHash, originalGoHash := sha256.New(), sha256.New(), sha256.New(), sha256.New()
	cases, goCases, audits := 0, 0, 0
	var failure map[string]any
	tables := []struct {
		Name     string
		Mul, Add int
	}{{"identity", 1, 0}, {"reverse", 255, 255}, {"affine-197-13", 197, 13}, {"affine-137-39", 137, 39}}
	groups := map[string]int{}

	check := func(group string, table []byte, c, s byte, ss, sp uint16, compareGo bool, index int) bool {
		// 用真實指令初始化 SS、SP、FLAGS；不修改原版 routine，也不存取私有 CPU。
		pre := o.Regs()
		writeWord(pre.SS, pre.SP-2, 0xbeef)
		flags := uint16(2)
		bits := []uint16{1, 4, 0x10, 0x40, 0x80, 0x800}
		for k, b := range bits {
			if index&(1<<k) != 0 {
				flags |= b
			}
		}
		// 局部常式固定 IF/TF=0；DF 與六個運算旗標依案例變化。
		if index&128 != 0 {
			flags |= 0x400
		}
		setup := []byte{0xb8, byte(ss), byte(ss >> 8), 0x8e, 0xd0, 0xbc, byte(sp + 2), byte((sp + 2) >> 8), 0x68, byte(flags), byte(flags >> 8), 0x9d}
		setupAddr := oracle.Addr{Seg: 0x7000, Off: 0}
		writeBoth(setupAddr, setup)
		o.CallNear(setupAddr, 0xbeef, oracle.CallRegs{})
		if err := o.Run(5); err != nil {
			panic(err)
		}
		writeWord(ss, sp, flags) // setup 的 PUSH 留下的 bytes，下一步會被返回位址覆蓋。
		state := append([]byte{c, s}, table...)
		writeBoth(stateAddr, state)
		writeWord(ss, sp, 0x0100)
		inputRegs := oracle.CallRegs{AX: uint16(index*19) ^ 0x5317, BX: uint16(index*23) ^ 0xb000,
			CX: uint16(index*31) ^ 0x1234, DX: uint16(index*37) ^ 0x4567, SI: uint16(index*41) ^ 0x6789,
			DI: uint16(index*43) ^ 0x789a, BP: uint16(index*47) ^ 0x89ab, DS: uint16(index*53) ^ 0x2345, ES: uint16(index*59) ^ 0x3456,
			SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		o.CallNear(entry, 0x0100, inputRegs)
		input := o.Regs()
		initialState := o.Bytes(stateAddr, 258)
		setRegs(&m, input)
		stackLo, stackSeg = input.SP-4, input.SS
		active = true
		before := o.Steps()
		err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x0100}), oracle.Budget(20))
		active = false
		if err != nil || o.Steps()-before != 13 {
			panic(fmt.Sprintf("original routine: %v, steps %d", err, o.Steps()-before))
		}
		C.sub_1ECE0_abi(&m)
		got, want := regs(&m), o.Regs()
		finalState := o.Bytes(stateAddr, 258)
		cases++
		groups[group]++
		if got != want || !bytes.Equal(cMemory[stateLinear:stateLinear+258], finalState) || len(unexpected) > 0 {
			failure = map[string]any{"group": group, "index": index, "initial_registers": input, "original": want, "c": got, "initial_state": hex.EncodeToString(initialState), "unexpected_writes": unexpected}
			return false
		}
		for n := uint16(0); n < 6; n++ {
			at := phys(stackSeg, stackLo+n)
			if cMemory[at] != o.Byte(oracle.Phys(at)) {
				panic("stack bytes differ")
			}
		}
		digestCase(originalHash, want, finalState)
		digestCase(cHash, got, cMemory[stateLinear:stateLinear+258])
		if compareGo {
			g, ok := rng.FromRaw(initialState)
			if !ok {
				panic("Go raw state")
			}
			value := g.Next()
			goState := g.Raw()
			goCases++
			if value != int(want.AX&255) || !bytes.Equal(goState, finalState) {
				failure = map[string]any{"group": group, "index": index, "go_value": value, "original_ax": want.AX}
				return false
			}
			goHash.Write(goState)
			goHash.Write([]byte{byte(value)})
			originalGoHash.Write(finalState)
			originalGoHash.Write([]byte{byte(want.AX)})
		}
		if cases%1024 == 0 {
			if !bytes.Equal(cMemory, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory differs")
			}
			audits++
		}
		return true
	}
	for n, table := range tables {
		data := make([]byte, 256)
		for i := range data {
			data[i] = byte(i*table.Mul + table.Add)
		}
		for i := 0; i < *max; i++ {
			if !check(table.Name, data, byte(i>>8), byte(i), 0x4000, 0x8000, true, i+n*65536) {
				goto finish
			}
		}
	}
	{
		data := make([]byte, 256)
		for i := range data {
			data[i] = byte(i)
		}
		for _, sp := range []uint16{0, 2, 0xfffc, 0xfffe} {
			for i := 0; i < 256; i++ {
				if !check("stack-boundary", data, byte(i), byte(255-i), 0x4000, sp, true, i) {
					goto finish
				}
			}
		}
		for _, sp := range []uint16{0xed00, 0xed04} {
			for i := 0; i < 256; i++ {
				if !check("stack-alias", data, byte(i), byte(255-i), cs, sp, false, i) {
					goto finish
				}
			}
		}
	}
finish:
	if failure == nil {
		if !bytes.Equal(cMemory, o.Bytes(oracle.Phys(0), 1<<20)) {
			panic("final full memory differs")
		}
		audits++
		if !bytes.Equal(o.Bytes(entry, 28), code) {
			panic("original code altered")
		}
	}
	report := map[string]any{"schema": "wolong-c-rng-parity-v1", "input_sha256": originalSHA, "routine_sha256": hex.EncodeToString(codeHash[:]),
		"ida_linear": 0x1ece0, "runtime_cs": cs, "runtime_offset": entry.Off, "fixed_inputs": tables, "cases": cases, "go_cases": goCases,
		"groups": groups, "full_memory_audits": audits, "observed_original_byte_writes": writes, "register_or_flag_mismatch": failure,
		"original_trace_sha256": hex.EncodeToString(originalHash.Sum(nil)), "c_trace_sha256": hex.EncodeToString(cHash.Sum(nil)),
		"go_state_trace_sha256": hex.EncodeToString(goHash.Sum(nil)), "passed": failure == nil,
		"original_go_state_trace_sha256": hex.EncodeToString(originalGoHash.Sum(nil)),
		"input_flags":                    "IF/TF=0; DF and six arithmetic flags vary by case index",
		"normal_stack":                   map[string]int{"SS": 0x4000, "entry_SP": 0x8000},
		"stack_boundary_SP":              []uint16{0, 2, 0xfffc, 0xfffe}, "stack_alias_SP": []uint16{0xed00, 0xed04},
		"scope":                "controlled local near-call; final 16-bit CPU state and memory; no normal player path or clock parity",
		"c_machine_code_match": false, "seed_method": "controlled fixed permutation tables and explicit c/s before each call"}
	encoded, _ := json.MarshalIndent(report, "", "  ")
	if err := os.WriteFile(*out, append(encoded, '\n'), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("C / original: %d cases; Go: %d; full-memory audits: %d; pass=%v\n", cases, goCases, audits, failure == nil)
	if failure != nil {
		os.Exit(1)
	}
}
