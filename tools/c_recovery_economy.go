//go:build ignore

// 月結及五個經濟函式的原版/C/Go 局部對照；依 spec/205 在 Docker 內執行。
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/rng.c"
#include "/repo/tools/c_recovery/economy.c"
#include "/repo/tools/c_recovery/economy_fixture.h"
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
	"os"
	"unsafe"

	"github.com/wicanr2/dosgolem/oracle"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
)

func ecSet(m *C.KiMachine16, r oracle.Regs) {
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
func ecRegs(m *C.KiMachine16) oracle.Regs {
	return oracle.Regs{AX: uint16(m.ax), BX: uint16(m.bx), CX: uint16(m.cx), DX: uint16(m.dx), SI: uint16(m.si), DI: uint16(m.di), BP: uint16(m.bp), SP: uint16(m.sp), DS: uint16(m.ds), ES: uint16(m.es), SS: uint16(m.ss), CS: uint16(m.cs), IP: uint16(m.ip), Flags: uint16(m.flags)}
}
func ecRegBytes(r oracle.Regs) []byte {
	b := make([]byte, 28)
	for i, v := range []uint16{r.AX, r.BX, r.CX, r.DX, r.SI, r.DI, r.BP, r.SP, r.DS, r.ES, r.SS, r.CS, r.IP, r.Flags} {
		binary.LittleEndian.PutUint16(b[i*2:], v)
	}
	return b
}
func ecSigned24(b []byte) int {
	v := int(b[0]) | int(b[1])<<8 | int(b[2])<<16
	if v&0x800000 != 0 {
		v -= 1 << 24
	}
	return v
}

func main() {
	out := flag.String("out", "/output/results/O2.json", "收據")
	only := flag.String("group", "", "單一矩陣群組")
	smoke := flag.Bool("smoke", false, "只跑矩陣前段")
	skipGo := flag.Bool("skip-go", false, "原版/C 研究探針")
	flag.Parse()
	raw, err := os.ReadFile("/orig/KI.EXE")
	if err != nil {
		panic(err)
	}
	inputHash := fmt.Sprintf("%x", sha256.Sum256(raw))
	if inputHash != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" {
		panic("input identity")
	}
	o, err := oracle.Load("/orig/KI.EXE", "/orig")
	if err != nil {
		panic(err)
	}
	defer o.Close()
	cs := o.Regs().CS
	functions := map[uint16]int{0x5358: 110, 0x5609: 34, 0x563b: 40, 0x54fc: 54, 0x55ec: 13, 0x5828: 55}
	routineHashes := map[string]string{}
	for target, size := range functions {
		want := raw[int(target)+0x200 : int(target)+0x200+size]
		if !bytes.Equal(o.Bytes(oracle.Addr{Seg: cs, Off: target}, size), want) {
			panic("routine identity")
		}
		routineHashes[fmt.Sprintf("sub_%X", 0x10000+uint32(target))] = fmt.Sprintf("%x", sha256.Sum256(want))
	}
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	mem := unsafe.Slice((*byte)(native), 1<<20)
	copy(mem, o.Bytes(oracle.Phys(0), 1<<20))
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	both := func(seg, off uint16, b []byte) {
		a := oracle.Addr{Seg: seg, Off: off}
		o.WriteBytes(a, b)
		copy(mem[a.Linear():], b)
	}
	word := func(seg, off, v uint16) { var b [2]byte; binary.LittleEndian.PutUint16(b[:], v); both(seg, off, b[:]) }
	bankSeg, bankOff, cityOff := uint16(0x2200), uint16(0x100), uint16(0x800)
	rngAddr := oracle.Addr{Seg: cs, Off: 0xecfc}
	settings := oracle.Addr{Seg: cs, Off: 0xd08}
	table := rng.New(12, 34, 56).Raw()
	var trace []byte
	targets := []uint16{0x5358, 0x5609, 0x563b, 0x54fc, 0x55ec, 0x5828, 0xece0, 0x53c6, 0x5695, 0x585f, 0x55a6, 0x2bd9, 0x5715, 0x578f, 0x22db, 0x2286, 0x57fe, 0x5e80}
	for _, target := range targets {
		t := target
		o.OnCall(oracle.Addr{Seg: cs, Off: t}, func(o *oracle.Oracle) {
			var b [2]byte
			binary.LittleEndian.PutUint16(b[:], t)
			trace = append(trace, b[:]...)
			r := o.Regs()
			trace = append(trace, ecRegBytes(r)...)
			for n := uint16(0); n < 64; n++ {
				trace = append(trace, o.Byte(oracle.Addr{Seg: r.DS, Off: r.SI + n}))
			}
			trace = append(trace, o.Bytes(settings, 8)...)
			trace = append(trace, o.Bytes(rngAddr, 2)...)
		})
	}
	for _, target := range targets[7:] {
		both(cs, target, []byte{0xc3})
	}
	cases, goCases, audits := 0, 0, 0
	groups := map[string]int{}
	originalDigest, cDigest := sha256.New(), sha256.New()
	var failure, goFailure map[string]any
	check := func(group string, target, ax, dx uint16, bank []byte, di uint16, city []byte, income uint32) bool {
		if *only != "" && *only != group {
			return true
		}
		if *smoke && groups[group] >= 32 {
			return true
		}
		trace = nil
		pre := o.Regs()
		word(pre.SS, pre.SP-2, 0xbeef)
		flags := uint16(2)
		if cases&1 != 0 {
			flags |= 0x400
		}
		setup := []byte{0xb8, 0, 0x40, 0x8e, 0xd0, 0xbc, 2, 0x80, 0x68, byte(flags), byte(flags >> 8), 0x9d}
		both(0x7000, 0, setup)
		o.CallNear(oracle.Addr{Seg: 0x7000, Off: 0}, 0xbeef, oracle.CallRegs{})
		if err := o.Run(5); err != nil {
			panic(err)
		}
		word(0x4000, 0x8000, 0x100)
		both(bankSeg, bankOff, bank)
		if target == 0x5358 {
			both(bankSeg, 0, bank)
		}
		if len(city) > 0 {
			both(bankSeg, cityOff, city)
		}
		word(cs, 0xd52, bankSeg)
		both(cs, 0xd08, []byte{1, 0, 2, 0, 3, 0, 4, 0})
		both(cs, 0xd10, []byte{9, 0, 8, 0, 7, 0, 6, 0})
		both(cs, 0x53c6, []byte{0xb8, byte(income), byte(income >> 8), 0xb2, byte(income >> 16), 0xc3})
		initialRNG := append([]byte(nil), table...)
		initialRNG[0] = byte(cases * 19)
		initialRNG[1] = byte(cases * 31)
		both(cs, 0xecfc, initialRNG)
		r := oracle.CallRegs{AX: ax, BX: 0x2345, CX: 0x3456, DX: dx, SI: bankOff, DI: di, BP: 0x6789, DS: bankSeg, ES: 0x789a, SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		if target == 0x5358 {
			r.DS = cs
		}
		o.CallNear(oracle.Addr{Seg: cs, Off: target}, 0x100, r)
		input := o.Regs()
		ecSet(&m, input)
		f := C.KiEconomyFixture{income_low: C.uint16_t(income), income_high: C.uint8_t(income >> 16)}
		if err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x100}), oracle.Budget(20000)); err != nil {
			panic(err)
		}
		C.economy_run_fixture(&m, C.uint16_t(target), &f)
		got, want := ecRegs(&m), o.Regs()
		ct := make([]byte, 0, int(f.calls)*104)
		if int(f.calls) > 256 {
			panic("trace capacity")
		}
		for i := 0; i < int(f.calls); i++ {
			s := f.trace[i]
			var b [2]byte
			binary.LittleEndian.PutUint16(b[:], uint16(s.target))
			ct = append(ct, b[:]...)
			for _, v := range s.regs {
				binary.LittleEndian.PutUint16(b[:], uint16(v))
				ct = append(ct, b[:]...)
			}
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.record[0]), 64)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.settings[0]), 8)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.rng[0]), 2)...)
		}
		cases++
		groups[group]++
		dataAt := oracle.Addr{Seg: bankSeg, Off: bankOff}
		dataSize := len(bank)
		if target == 0x5358 {
			dataAt.Off = 0
		}
		originalState := o.Bytes(dataAt, dataSize)
		cState := mem[int(dataAt.Linear()) : int(dataAt.Linear())+dataSize]
		if got != want || !bytes.Equal(trace, ct) || !bytes.Equal(cState, originalState) || !bytes.Equal(mem[settings.Linear():settings.Linear()+8], o.Bytes(settings, 8)) || !bytes.Equal(mem[rngAddr.Linear():rngAddr.Linear()+258], o.Bytes(rngAddr, 258)) {
			failure = map[string]any{"group": group, "case": cases - 1, "target": target, "input": input, "original": want, "c": got, "original_state": hex.EncodeToString(originalState), "c_state": hex.EncodeToString(cState), "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct)}
			return false
		}
		for off := uint16(0x7fe0); off < 0x8002; off++ {
			a := oracle.Addr{Seg: 0x4000, Off: off}
			if mem[a.Linear()] != o.Byte(a) {
				panic("stack mismatch")
			}
		}
		originalDigest.Write(ecRegBytes(want))
		originalDigest.Write(originalState)
		originalDigest.Write(trace)
		cDigest.Write(ecRegBytes(got))
		cDigest.Write(cState)
		cDigest.Write(ct)
		if !*skipGo {
			goOK, compared := true, false
			var details any
			switch target {
			case 0x5609, 0x563b:
				fund := ecSigned24(bank[0x20:])
				amount := int(ax) | int(byte(dx))<<16
				if fund >= economy.MinFunds && fund <= economy.MaxFunds && amount <= economy.MaxFunds {
					v := fund + amount
					if target == 0x563b {
						v = fund - amount
					}
					expected := economy.ClampFunds(v)
					goOK = expected == ecSigned24(originalState[0x20:])
					compared = true
					details = expected
				}
			case 0x55ec:
				compared = true
				expected := economy.ClampReserve(int(ax) + int(dx))
				goOK = expected == int(want.AX)
				details = expected
			case 0x54fc:
				cap := economy.City{X: int(binary.LittleEndian.Uint16(bank[8:])), Y: int(binary.LittleEndian.Uint16(bank[10:]))}
				c := economy.City{X: int(binary.LittleEndian.Uint16(city[8:])), Y: int(binary.LittleEndian.Uint16(city[10:])), Owner: 0, Production: 120}
				faction := economy.Faction{Capital: cap, TaxRate: 100}
				res := economy.Settle(&faction, []economy.City{c}, 0, rng.NewFixed(0))
				goOK = res.GrossBase == 120/int(want.BX)
				compared = true
				details = res.GrossBase
			case 0x5828:
				fund := ecSigned24(bank[0x20:])
				if fund >= economy.MinFunds && fund <= economy.MaxFunds {
					g, _ := rng.FromRaw(initialRNG)
					faction := economy.Faction{Funds: fund}
					for i := 0; i < 3; i++ {
						faction.Reserves[i] = int(binary.LittleEndian.Uint16(bank[4+i*2:]))
					}
					res := economy.Settle(&faction, nil, 0, g)
					goOK = bytes.Equal(g.Raw(), o.Bytes(rngAddr, 258))
					for i := 0; i < 3; i++ {
						goOK = goOK && faction.Reserves[i] == int(binary.LittleEndian.Uint16(originalState[4+i*2:]))
					}
					compared = true
					details = map[string]any{"faction": faction, "result": res, "initial_fund": fund}
				}
			}
			if compared {
				goCases++
				if !goOK {
					goFailure = map[string]any{"group": group, "case": cases - 1, "details": details, "original_state": hex.EncodeToString(originalState)}
					return false
				}
			}
		}
		if cases%2048 == 0 {
			if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory mismatch")
			}
			audits++
		}
		return true
	}
	money := func(fund int) []byte {
		b := make([]byte, 64)
		v := uint32(fund) & 0xffffff
		b[0x20] = byte(v)
		b[0x21] = byte(v >> 8)
		b[0x22] = byte(v >> 16)
		return b
	}
	for _, group := range []string{"credit", "debit"} {
		target := uint16(0x5609)
		if group == "debit" {
			target = 0x563b
		}
		for _, fund := range []int{-655001, -655000, -1, 0, 655000, 655001} {
			for v := 0; v < 65536; v++ {
				if !check(group, target, uint16(v), 0, money(fund), cityOff, nil, 0) {
					goto finish
				}
			}
		}
		for _, fund := range []int{-8388608, -655001, -655000, -1, 0, 655000, 655001, 8388607} {
			for _, hi := range []byte{1, 9, 0xf6, 0xff} {
				for _, low := range []uint16{0, 1, 0x167, 0x168, 0xfe97, 0xfe98, 0xffff} {
					if !check(group, target, low, uint16(hi), money(fund), cityOff, nil, 0) {
						goto finish
					}
				}
			}
		}
	}
	for _, add := range []uint16{0, 1, 35, 36, 65500, 65535} {
		for ax := 0; ax < 65536; ax++ {
			if !check("reserve", 0x55ec, uint16(ax), add, money(0), cityOff, nil, 0) {
				goto finish
			}
		}
	}
	for dx := 0; dx < 384; dx++ {
		for dy := 0; dy < 256; dy++ {
			for direction := 0; direction < 4; direction++ {
				bank, city := money(0), make([]byte, 64)
				x1, x2, y1, y2 := 0, dx, 0, dy
				if direction&1 != 0 {
					x1, x2 = x2, x1
				}
				if direction&2 != 0 {
					y1, y2 = y2, y1
				}
				binary.LittleEndian.PutUint16(bank[8:], uint16(x1))
				binary.LittleEndian.PutUint16(bank[10:], uint16(y1))
				binary.LittleEndian.PutUint16(city[8:], uint16(x2))
				binary.LittleEndian.PutUint16(city[10:], uint16(y2))
				if !check("distance", 0x54fc, 0x1234, 0x4567, bank, cityOff, city, 0) {
					goto finish
				}
			}
		}
	}
	for high := 0; high < 65536; high++ {
		for _, low := range []byte{0, 1, 255} {
			bank := money(0)
			bank[0x20] = low
			bank[0x21] = byte(high)
			bank[0x22] = byte(high >> 8)
			for i, v := range []uint16{0, 65500, 12345} {
				binary.LittleEndian.PutUint16(bank[4+i*2:], v)
			}
			if !check("deficit", 0x5828, 0x1234, 0x4567, bank, cityOff, nil, 0) {
				goto finish
			}
		}
	}
	for slot := 0; slot < 22; slot++ {
		for status := 0; status < 256; status++ {
			bank := make([]byte, 22*64)
			at := slot * 64
			bank[at] = byte(status)
			copy(bank[at+0x20:], money(-1)[0x20:0x23])
			binary.LittleEndian.PutUint16(bank[at+0x1a:], 500)
			for i := 0; i < 3; i++ {
				binary.LittleEndian.PutUint16(bank[at+4+i*2:], 65500)
			}
			if !check("monthly", 0x5358, 0x1234, 0x4567, bank, cityOff, nil, 123) {
				goto finish
			}
		}
	}
	for pattern := 0; pattern < 16; pattern++ {
		bank := make([]byte, 22*64)
		for slot := 0; slot < 22; slot++ {
			at := slot * 64
			if pattern == 0 || slot%4 == pattern%4 {
				bank[at] = 0x80 | byte(pattern)
			}
			fund := slot*9000 - 65000
			v := uint32(fund) & 0xffffff
			bank[at+0x20] = byte(v)
			bank[at+0x21] = byte(v >> 8)
			bank[at+0x22] = byte(v >> 16)
			binary.LittleEndian.PutUint16(bank[at+0x1a:], uint16(slot*100+500))
			for i := 0; i < 3; i++ {
				binary.LittleEndian.PutUint16(bank[at+4+i*2:], 65500)
			}
		}
		if !check("monthly", 0x5358, 0x1234, 0x4567, bank, cityOff, nil, uint32(50000+pattern*123)) {
			goto finish
		}
	}
finish:
	if failure == nil && goFailure == nil {
		if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
			panic("final memory")
		}
		audits++
	}
	report := map[string]any{"schema": "wolong-c-economy-parity-v1", "input_sha256": inputHash, "routine_sha256": routineHashes, "cases": cases, "go_cases": goCases, "groups": groups, "full_memory_audits": audits, "passed": failure == nil && goFailure == nil, "mismatch": failure, "go_mismatch": goFailure, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "c_machine_code_match": false, "scope": "Six local functions; original money/deficit/RNG calls retained, other monthly callees have explicit MOV/RET fixtures; IF/TF=0; disjoint stack; AF follows dosgolem model"}
	encoded, _ := json.MarshalIndent(report, "", "  ")
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*out, encoded, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("economy original/C %d; Go %d; audits %d; pass=%v\n", cases, goCases, audits, report["passed"])
	if failure != nil || goFailure != nil {
		os.Exit(1)
	}
}
