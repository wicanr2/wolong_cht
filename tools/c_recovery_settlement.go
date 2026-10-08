//go:build ignore

// 據點收入、募兵與月結接線的原版/C/Go 對照；spec/207。
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/rng.c"
#include "/repo/tools/c_recovery/economy.c"
#include "/repo/tools/c_recovery/settlement.c"
#include "/repo/tools/c_recovery/settlement_fixture.h"
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

func stSet(m *C.KiMachine16, r oracle.Regs) {
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
func stRegs(m *C.KiMachine16) oracle.Regs {
	return oracle.Regs{AX: uint16(m.ax), BX: uint16(m.bx), CX: uint16(m.cx), DX: uint16(m.dx), SI: uint16(m.si), DI: uint16(m.di), BP: uint16(m.bp), SP: uint16(m.sp), DS: uint16(m.ds), ES: uint16(m.es), SS: uint16(m.ss), CS: uint16(m.cs), IP: uint16(m.ip), Flags: uint16(m.flags)}
}
func stRegBytes(r oracle.Regs) []byte {
	b := make([]byte, 28)
	for i, v := range []uint16{r.AX, r.BX, r.CX, r.DX, r.SI, r.DI, r.BP, r.SP, r.DS, r.ES, r.SS, r.CS, r.IP, r.Flags} {
		binary.LittleEndian.PutUint16(b[i*2:], v)
	}
	return b
}
func stWord(b []byte, off int) uint16   { return binary.LittleEndian.Uint16(b[off:]) }
func stPut(b []byte, off int, v uint16) { binary.LittleEndian.PutUint16(b[off:], v) }
func st24(b []byte, off int) int        { return int(b[off]) | int(b[off+1])<<8 | int(b[off+2])<<16 }
func stStore24(b []byte, off, v int) {
	b[off] = byte(v)
	b[off+1] = byte(v >> 8)
	b[off+2] = byte(v >> 16)
}
func stSigned(b []byte, off int) int {
	v := st24(b, off)
	if v&0x800000 != 0 {
		v -= 1 << 24
	}
	return v
}
func stCitiesForGross(gross int) []economy.City {
	var cities []economy.City
	for gross > 0 {
		n := gross
		if n > 32767 {
			n = 32767
		}
		cities = append(cities, economy.City{Owner: 0, Production: n * 2})
		gross -= n
	}
	return cities
}
func stTyped(b []byte, owner int, player bool, tax byte, caps [3]uint16) (economy.Faction, []economy.City) {
	at := owner * 64
	capital := 0x840 + int(b[at+3])*32
	f := economy.Faction{Funds: stSigned(b, at+0x20), Expense: st24(b, at+0x1a), Capital: economy.City{X: int(stWord(b, capital+8)), Y: int(stWord(b, capital+10))}, TaxRate: int(tax), AI: !player}
	for i := 0; i < 3; i++ {
		f.Reserves[i] = int(stWord(b, at+4+i*2))
		f.RecruitCap[i] = int(caps[i])
	}
	for off := 0x2240; off < 0x2240+127*32; off += 32 {
		if b[off] >= 0x80 && int(b[off+1]) == owner {
			f.CorpsWeight = (f.CorpsWeight + int(stWord(b, off+4))) & 0xffff
		}
	}
	cities := make([]economy.City, 192)
	for i := range cities {
		at := 0x840 + i*32
		cities[i] = economy.City{Owner: int(b[at+1]), X: int(stWord(b, at+8)), Y: int(stWord(b, at+10)), Production: int(stWord(b, at+14))}
	}
	return f, cities
}

type stCase struct {
	group       string
	target      uint16
	bank, local []byte
	owner       int
	player      bool
	tax         byte
	caps        [3]uint16
	bx          uint16
	goCheck     bool
}

func stBank() []byte {
	b := make([]byte, 0x5240)
	b[0] = 0x80
	stStore24(b, 0x20, 600000)
	for i := 0; i < 192; i++ {
		b[0x840+i*32+1] = 0xff
	}
	return b
}
func main() {
	out := flag.String("out", "/output/results/O2.json", "收據")
	only := flag.String("group", "", "單一群組")
	smoke := flag.Bool("smoke", false, "每群組少量案例")
	skipGo := flag.Bool("skip-go", false, "只比較原版/C")
	probe := flag.Bool("go-probe", false, "稅率與累計進位反例")
	flag.Parse()
	raw, err := os.ReadFile("/orig/KI.EXE")
	if err != nil {
		panic(err)
	}
	exeHash := fmt.Sprintf("%x", sha256.Sum256(raw))
	if exeHash != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" {
		panic("input identity")
	}
	scenarios, err := os.ReadFile("/orig/SINARIO.DAT")
	if err != nil || len(scenarios) != 88832 {
		panic("scenario identity")
	}
	o, err := oracle.Load("/orig/KI.EXE", "/orig")
	if err != nil {
		panic(err)
	}
	defer o.Close()
	cs := o.Regs().CS
	functions := map[uint16]int{0x53c6: 144, 0x5456: 57, 0x548f: 109, 0x5538: 15, 0x5547: 95}
	routines := map[string]string{}
	for target, size := range functions {
		b := raw[int(target)+0x200 : int(target)+0x200+size]
		if !bytes.Equal(o.Bytes(oracle.Addr{Seg: cs, Off: target}, size), b) {
			panic("routine identity")
		}
		routines[fmt.Sprintf("sub_%X", 0x10000+uint32(target))] = fmt.Sprintf("%x", sha256.Sum256(b))
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
	bankSeg, ss, bp := uint16(0x2200), uint16(0x4000), uint16(0x2000)
	display := oracle.Addr{Seg: cs, Off: 0xd02}
	rngAddr := oracle.Addr{Seg: cs, Off: 0xecfc}
	bankAddr := oracle.Addr{Seg: bankSeg, Off: 0}
	localAddr := oracle.Addr{Seg: ss, Off: bp}
	table := rng.New(12, 34, 56).Raw()
	var trace []byte
	targets := []uint16{0x5358, 0x53c6, 0x5456, 0x548f, 0x5538, 0x5547, 0x54fc, 0x55ec, 0x5609, 0x563b, 0x5828, 0xece0, 0x5695, 0x585f, 0x55a6, 0x2bd9, 0x5715, 0x578f, 0x22db, 0x2286, 0x57fe, 0x5e80}
	for _, target := range targets {
		t := target
		o.OnCall(oracle.Addr{Seg: cs, Off: t}, func(o *oracle.Oracle) {
			r := o.Regs()
			var b [2]byte
			binary.LittleEndian.PutUint16(b[:], t)
			trace = append(trace, b[:]...)
			trace = append(trace, stRegBytes(r)...)
			for n := uint16(0); n < 64; n++ {
				trace = append(trace, o.Byte(oracle.Addr{Seg: r.DS, Off: r.SI + n}))
			}
			trace = append(trace, o.Bytes(oracle.Addr{Seg: r.SS, Off: r.BP}, 10)...)
			trace = append(trace, o.Bytes(display, 22)...)
			trace = append(trace, o.Bytes(rngAddr, 2)...)
		})
	}
	for _, target := range targets[12:] {
		both(cs, target, []byte{0xc3})
	}
	f := C.KiSettlementFixture{}
	cases, goCases, audits := 0, 0, 0
	groups, goGroups := map[string]int{}, map[string]int{}
	od, cd := sha256.New(), sha256.New()
	var failure, goFailure map[string]any
	var probeResults []map[string]any
	check := func(c stCase) bool {
		if *only != "" && *only != c.group {
			return true
		}
		if *smoke && groups[c.group] >= 20 {
			return true
		}
		trace = nil
		f.calls = 0
		pre := o.Regs()
		word(pre.SS, pre.SP-2, 0xbeef)
		flags := uint16(2)
		if cases&1 != 0 {
			flags |= 0x400
		}
		both(0x7000, 0, []byte{0xb8, 0, 0x40, 0x8e, 0xd0, 0xbc, 2, 0x80, 0x68, byte(flags), byte(flags >> 8), 0x9d})
		o.CallNear(oracle.Addr{Seg: 0x7000, Off: 0}, 0xbeef, oracle.CallRegs{})
		if err := o.Run(5); err != nil {
			panic(err)
		}
		word(ss, 0x8000, 0x100)
		both(bankSeg, 0, c.bank)
		both(ss, bp, c.local)
		globals := make([]byte, 22)
		globals[6] = c.tax
		for i := 0; i < 3; i++ {
			stPut(globals, 8+i*2, c.caps[i])
			stPut(globals, 14+i*2, c.caps[i]/2)
		}
		both(cs, 0xd02, globals)
		playerOff := uint16(0xffff)
		if c.player {
			playerOff = uint16(c.owner * 64)
		}
		word(cs, 0xcfd, playerOff)
		word(cs, 0xd52, bankSeg)
		initialRNG := append([]byte(nil), table...)
		initialRNG[0] = byte(cases * 19)
		initialRNG[1] = byte(cases * 31)
		both(cs, 0xecfc, initialRNG)
		r := oracle.CallRegs{AX: 0x1234, BX: c.bx, CX: 0x3400 | uint16(c.owner), DX: 0x4567, SI: uint16(c.owner * 64), DI: 0x840, BP: bp, DS: bankSeg, ES: 0x789a, SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		if c.target == 0x5358 {
			r.DS = cs
		}
		o.CallNear(oracle.Addr{Seg: cs, Off: c.target}, 0x100, r)
		input := o.Regs()
		stSet(&m, input)
		if err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x100}), oracle.Budget(160000)); err != nil {
			panic(err)
		}
		C.settlement_run(&m, C.uint16_t(c.target), &f)
		got, want := stRegs(&m), o.Regs()
		ct := make([]byte, 0, int(f.calls)*128)
		if int(f.calls) > 1024 {
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
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.locals[0]), 10)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.display[0]), 22)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.rng[0]), 2)...)
		}
		bankResult, localResult, displayResult := o.Bytes(bankAddr, len(c.bank)), o.Bytes(localAddr, 10), o.Bytes(display, 22)
		cases++
		groups[c.group]++
		if *probe {
			income := st24(localResult, 0)
			if c.target == 0x53c6 {
				income = int(want.AX) | int(byte(want.DX))<<16
			}
			probeResults = append(probeResults, map[string]any{"group": c.group, "original_income": income, "original_locals": hex.EncodeToString(localResult), "original_reserves": []uint16{stWord(bankResult, c.owner*64+4), stWord(bankResult, c.owner*64+6), stWord(bankResult, c.owner*64+8)}})
		}
		if got != want || !bytes.Equal(trace, ct) || !bytes.Equal(mem[bankAddr.Linear():bankAddr.Linear()+uint32(len(c.bank))], bankResult) || !bytes.Equal(mem[localAddr.Linear():localAddr.Linear()+10], localResult) || !bytes.Equal(mem[display.Linear():display.Linear()+22], displayResult) || !bytes.Equal(mem[rngAddr.Linear():rngAddr.Linear()+258], o.Bytes(rngAddr, 258)) {
			failure = map[string]any{"group": c.group, "case": cases - 1, "input": input, "original": want, "c": got, "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct), "original_locals": hex.EncodeToString(localResult), "c_locals": hex.EncodeToString(mem[localAddr.Linear() : localAddr.Linear()+10])}
			return false
		}
		for off := uint16(0x7fc0); off < 0x8002; off++ {
			a := oracle.Addr{Seg: ss, Off: off}
			if mem[a.Linear()] != o.Byte(a) {
				panic("stack mismatch")
			}
		}
		od.Write(stRegBytes(want))
		od.Write(bankResult)
		od.Write(localResult)
		od.Write(displayResult)
		od.Write(trace)
		cd.Write(stRegBytes(got))
		cd.Write(mem[bankAddr.Linear() : bankAddr.Linear()+uint32(len(c.bank))])
		cd.Write(localResult)
		cd.Write(displayResult)
		cd.Write(ct)
		if c.goCheck && !*skipGo {
			ok := true
			var details any
			switch c.target {
			case 0x548f:
				gross := st24(c.local, 0)
				gf := economy.Faction{TaxRate: int(c.tax)}
				res := economy.Settle(&gf, stCitiesForGross(gross), 0, rng.NewFixed(0))
				ok = res.Income == st24(localResult, 0)
				details = map[string]any{"gross": gross, "tax": c.tax, "go_income": res.Income, "original_income": st24(localResult, 0)}
			case 0x5538, 0x5547:
				div := int(c.bx)
				x := 0
				if div == 3 {
					x = 81
				}
				if div == 4 {
					x = 201
				}
				city := economy.City{Owner: 0, X: x, Y: int(stWord(c.bank, 0x84a)), Production: int(stWord(c.bank, 0x84e))}
				if city.Y > 80 && div == 2 {
					city.X = 0
				}
				capital := economy.City{Y: city.Y}
				gf := economy.Faction{Capital: capital, TaxRate: 100, RecruitCap: [3]int{65500, 65500, 65500}}
				res := economy.Settle(&gf, []economy.City{city}, 0, rng.NewFixed(0))
				if c.target == 0x5538 {
					ok = res.GrossBase == st24(localResult, 0)
				} else {
					for i := 0; i < 3; i++ {
						ok = ok && res.Recruited[i] == int(stWord(localResult, 4+i*2))
					}
				}
				details = res
			case 0x53c6:
				gf, cities := stTyped(c.bank, c.owner, c.player, c.tax, c.caps)
				gf.Funds = 600000
				gf.Expense = st24(c.bank, c.owner*64+0x1a)
				gr, _ := rng.FromRaw(initialRNG)
				res := economy.Settle(&gf, cities, c.owner, gr)
				ok = res.Income == int(want.AX)|int(byte(want.DX))<<16
				ok = ok && gf.Cities == int(want.DX>>8)
				for i := 0; i < 3; i++ {
					ok = ok && gf.Reserves[i] == int(stWord(bankResult, c.owner*64+4+i*2))
				}
				details = map[string]any{"go": gf, "result": res, "original_income": int(want.AX) | int(byte(want.DX))<<16}
			case 0x5456:
				gross := st24(c.local, 0)
				gf, _ := stTyped(c.bank, c.owner, false, c.tax, c.caps)
				gf.Funds = 600000
				gf.Expense = st24(c.bank, c.owner*64+0x1a)
				res := economy.Settle(&gf, stCitiesForGross(gross), 0, rng.NewFixed(0))
				ok = res.Income == st24(localResult, 0)
				if gross >= 64 {
					sum := 0
					for _, n := range res.Recruited {
						sum += n
					}
					ok = ok && ((sum == 0) == (want.Flags&1 != 0))
				}
				details = res
			case 0x5358:
				g, _ := rng.FromRaw(initialRNG)
				for owner := 0; owner < 22; owner++ {
					if c.bank[owner*64] < 0x80 {
						continue
					}
					gf, cities := stTyped(c.bank, owner, c.player && owner == c.owner, c.tax, c.caps)
					res := economy.Settle(&gf, cities, owner, g)
					at := owner * 64
					ok = ok && gf.Funds == stSigned(bankResult, at+0x20) && gf.Cities == int(bankResult[at+0x23]) && st24(bankResult, at+0x1a) == 0
					for i := 0; i < 3; i++ {
						ok = ok && gf.Reserves[i] == int(stWord(bankResult, at+4+i*2))
					}
					if !ok {
						details = map[string]any{"owner": owner, "go": gf, "result": res}
						break
					}
				}
				ok = ok && bytes.Equal(g.Raw(), o.Bytes(rngAddr, 258))
			}
			goCases++
			goGroups[c.group]++
			if !ok {
				goFailure = map[string]any{"group": c.group, "case": cases - 1, "details": details}
				return false
			}
		}
		if cases%4096 == 0 {
			if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory mismatch")
			}
			audits++
		}
		return true
	}
	caps := [3]uint16{65500, 65500, 65500}
	if *probe {
		for _, gross := range []int{66303, 131071, 6291264} {
			b := stBank()
			l := make([]byte, 10)
			stStore24(l, 0, gross)
			if !check(stCase{group: "player-tax", target: 0x548f, bank: b, local: l, player: true, tax: 99, caps: caps, goCheck: true}) {
				goto finish
			}
		}
		b := stBank()
		for i := 0; i < 192; i++ {
			at := 0x840 + i*32
			b[at+1] = 0
			stPut(b, at+14, 65535)
		}
		if !check(stCase{group: "city-settlement", target: 0x53c6, bank: b, local: make([]byte, 10), player: true, tax: 100, caps: caps, goCheck: true}) {
			goto finish
		}
		goto finish
	}
	for _, div := range []uint16{2, 3, 4} {
		for _, initial := range []int{0, 65535, 0xffffff} {
			for production := 0; production < 65536; production++ {
				b := stBank()
				stPut(b, 0x84e, uint16(production))
				l := make([]byte, 10)
				stStore24(l, 0, initial)
				if !check(stCase{group: "income", target: 0x5538, bank: b, local: l, bx: div, goCheck: initial == 0}) {
					goto finish
				}
			}
		}
	}
	for _, div := range []uint16{2, 3, 4} {
		for _, y := range []uint16{0, 79, 80, 149, 150, 255} {
			for production := 0; production < 65536; production++ {
				b := stBank()
				stPut(b, 0x84e, uint16(production))
				stPut(b, 0x84a, y)
				if !check(stCase{group: "recruit", target: 0x5547, bank: b, local: make([]byte, 10), bx: div, goCheck: true}) {
					goto finish
				}
			}
		}
	}
	for _, gross := range []int{0, 1, 254, 255, 256, 257, 32767, 65535, 65536, 65537, 66048, 66303, 131071, 131072, 6291264, 0xffffff} {
		for tax := 0; tax <= 100; tax++ {
			for capped := 0; capped < 2; capped++ {
				l := make([]byte, 10)
				stStore24(l, 0, gross)
				for i, v := range []uint16{3000, 6000, 9000} {
					stPut(l, 4+i*2, v)
				}
				limits := caps
				if capped == 1 {
					limits = [3]uint16{123, 456, 789}
				}
				if !check(stCase{group: "player-tax", target: 0x548f, bank: stBank(), local: l, player: true, tax: byte(tax), caps: limits, goCheck: gross <= 6291264}) {
					goto finish
				}
			}
		}
	}
	for _, gross := range []int{0, 64, 255, 256, 65535, 65536, 131071, 6291264} {
		for _, expense := range []int{0, 255, 256, 65535, 65536, 655000} {
			for slot := 0; slot < 128; slot++ {
				b := stBank()
				stStore24(b, 0x1a, expense)
				if slot < 127 {
					off := 0x2240 + slot*32
					b[off] = 0x80
					b[off+1] = 0
					stPut(b, off+4, 65535)
				}
				l := make([]byte, 10)
				stStore24(l, 0, gross)
				if !check(stCase{group: "ai-gate", target: 0x5456, bank: b, local: l, caps: caps, goCheck: true}) {
					goto finish
				}
			}
		}
	}
	for owner := 0; owner < 22; owner++ {
		for _, n := range []int{0, 1, 3, 64, 192} {
			for mode := 0; mode < 4; mode++ {
				b := stBank()
				at := owner * 64
				b[at] = 0x80
				stStore24(b, at+0x20, 600000)
				stPut(b, at+0x1a, 200)
				for i := 0; i < n; i++ {
					off := 0x840 + i*32
					b[off+1] = byte(owner)
					stPut(b, off+8, uint16(i%384))
					stPut(b, off+10, uint16(i%256))
					stPut(b, off+14, uint16(1000+i*317))
				}
				if !check(stCase{group: "city-settlement", target: 0x53c6, bank: b, local: make([]byte, 10), owner: owner, player: mode&1 == 0, tax: byte(30 + mode*23), caps: caps, goCheck: true}) {
					goto finish
				}
			}
		}
	}
	for scenario := 0; scenario < 4; scenario++ {
		base := scenario * 22208
		originalBank := scenarios[base+0x80 : base+0x80+0x5240]
		for _, player := range []int{0, 7, 21} {
			b := append([]byte(nil), originalBank...)
			if !check(stCase{group: "scenario-monthly", target: 0x5358, bank: b, local: make([]byte, 10), owner: player, player: true, tax: 50, caps: [3]uint16{5000, 5000, 5000}, goCheck: true}) {
				goto finish
			}
		}
	}
finish:
	if failure == nil && goFailure == nil {
		if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
			panic("final memory")
		}
		audits++
	}
	report := map[string]any{"schema": "wolong-c-settlement-parity-v1", "input_sha256": exeHash, "scenario_sha256": fmt.Sprintf("%x", sha256.Sum256(scenarios)), "routine_sha256": routines, "cases": cases, "go_cases": goCases, "groups": groups, "go_groups": goGroups, "full_memory_audits": audits, "passed": failure == nil && goFailure == nil, "mismatch": failure, "go_mismatch": goFailure, "original_state_sha256": hex.EncodeToString(od.Sum(nil)), "c_state_sha256": hex.EncodeToString(cd.Sum(nil)), "c_machine_code_match": false, "rng_fixture": map[string]any{"table_seed_hms": []int{12, 34, 56}, "counter": "uint8(case_index*19)", "state_index": "uint8(case_index*31)", "go": "rng.FromRaw identical initial state", "rerolls": 0}, "scope": "Five local functions and linked economic monthly path; nine tail callees and redraw have explicit RET fixtures; legal indices/tax/divisors; IF/TF=0; disjoint stack"}
	report["probe_results"] = probeResults
	encoded, _ := json.MarshalIndent(report, "", "  ")
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*out, encoded, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("settlement original/C %d; Go %d; audits %d; pass=%v\n", cases, goCases, audits, report["passed"])
	if failure != nil || goFailure != nil {
		os.Exit(1)
	}
}
