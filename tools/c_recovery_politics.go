//go:build ignore

// 二十一函式政治／俘虜依賴閉包與完整 C 月結規則對照；spec/211。
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/rng.c"
#include "/repo/tools/c_recovery/economy.c"
#include "/repo/tools/c_recovery/settlement.c"
#include "/repo/tools/c_recovery/world_update.c"
#include "/repo/tools/c_recovery/politics.c"
#include "/repo/tools/c_recovery/politics_fixture.h"
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
	"unsafe"
)

func pSet(m *C.KiMachine16, r oracle.Regs) {
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
func pRegs(m *C.KiMachine16) oracle.Regs {
	return oracle.Regs{AX: uint16(m.ax), BX: uint16(m.bx), CX: uint16(m.cx), DX: uint16(m.dx), SI: uint16(m.si), DI: uint16(m.di), BP: uint16(m.bp), SP: uint16(m.sp), DS: uint16(m.ds), ES: uint16(m.es), SS: uint16(m.ss), CS: uint16(m.cs), IP: uint16(m.ip), Flags: uint16(m.flags)}
}
func pRegBytes(r oracle.Regs) []byte {
	b := make([]byte, 28)
	for i, v := range []uint16{r.AX, r.BX, r.CX, r.DX, r.SI, r.DI, r.BP, r.SP, r.DS, r.ES, r.SS, r.CS, r.IP, r.Flags} {
		binary.LittleEndian.PutUint16(b[i*2:], v)
	}
	return b
}
func pPut(b []byte, at int, v uint16) { binary.LittleEndian.PutUint16(b[at:], v) }
func p24(b []byte, at, v int)         { b[at] = byte(v); b[at+1] = byte(v >> 8); b[at+2] = byte(v >> 16) }
func pBank() []byte {
	b := make([]byte, 0x5240)
	for i := 0; i < 22; i++ {
		at := i * 64
		b[at+0x19] = 0xff
		b[at+0x2a] = 0xff
		b[at+0x28] = 32
		p24(b, at+0x20, 300000)
		b[at+0x23] = 10
		if i < 3 {
			b[at] = 0x80
		}
		for j := 0; j < 3; j++ {
			pPut(b, at+4+j*2, 5000)
		}
	}
	for i := 0; i < 192; i++ {
		at := 0x840 + i*32
		b[at+1] = 0xff
		b[at+0x19] = 0xff
		pPut(b, at+12, 30000)
		pPut(b, at+14, 3000)
		b[at+16] = 100
		b[at+17] = 100
	}
	for i := 0; i < 127; i++ {
		at := 0x4240 + i*32
		b[at+0x19] = 0xff
		b[at+0x1c] = 0xff
		b[at+0x1d] = 0xff
	}
	for i := 0; i < 22*24; i++ {
		b[0x600+i] = 0x90
	}
	return b
}
func pScratch() []byte {
	b := make([]byte, 1056)
	for i := range b {
		b[i] = 0xff
	}
	return b
}

type pCase struct {
	group                              string
	target, ax, bx, si, di, dx, cursor uint16
	bank, queue, scratch               []byte
	player, counter, state             byte
}

func main() {
	out := flag.String("out", "/output/results/O2.json", "收據")
	only := flag.String("group", "", "單一群組")
	smoke := flag.Bool("smoke", false, "每群組少量案例")
	flag.Parse()
	raw, err := os.ReadFile("/orig/KI.EXE")
	if err != nil {
		panic(err)
	}
	eh := fmt.Sprintf("%x", sha256.Sum256(raw))
	if eh != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" {
		panic("input identity")
	}
	scenarios, err := os.ReadFile("/orig/SINARIO.DAT")
	if err != nil || len(scenarios) != 88832 {
		panic("scenario")
	}
	o, err := oracle.Load("/orig/KI.EXE", "/orig")
	if err != nil {
		panic(err)
	}
	defer o.Close()
	cs := o.Regs().CS
	funcs := map[uint16]int{0x585f: 58, 0x5899: 167, 0x5940: 80, 0x2ad2: 34, 0x5990: 22, 0x301c: 50, 0x2bd9: 121, 0x2c52: 141, 0x2cdf: 91, 0x2d3a: 30, 0x2fb1: 14, 0x2d58: 96, 0x2db8: 59, 0x2df3: 64, 0x30f0: 26, 0x310a: 15, 0x3091: 58, 0x2e33: 86, 0x2e89: 114, 0x2efb: 118, 0x2f71: 64}
	rh := map[string]string{}
	for target, size := range funcs {
		b := raw[int(target)+0x200 : int(target)+0x200+size]
		if !bytes.Equal(o.Bytes(oracle.Addr{Seg: cs, Off: target}, size), b) {
			panic("routine identity")
		}
		rh[fmt.Sprintf("sub_%X", 0x10000+uint32(target))] = fmt.Sprintf("%x", sha256.Sum256(b))
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
	bankSeg, queueSeg, scratchSeg, ss := uint16(0x2200), uint16(0x6000), uint16(0x6500), uint16(0x4000)
	ba, qa, sa := oracle.Addr{Seg: bankSeg, Off: 0}, oracle.Addr{Seg: queueSeg, Off: 0}, oracle.Addr{Seg: scratchSeg, Off: 0}
	ga, ra := oracle.Addr{Seg: cs, Off: 0xcf0}, oracle.Addr{Seg: cs, Off: 0xecfc}
	table := rng.New(12, 34, 56).Raw()
	var trace []byte
	targets := []uint16{0x5358, 0x53c6, 0x5456, 0x548f, 0x5538, 0x5547, 0x54fc, 0x55ec, 0x5609, 0x563b, 0x5828, 0xece0, 0x5695, 0x55a6, 0x2fbf, 0x57fe, 0x5715, 0x578f, 0x30cb, 0x3119, 0x22db, 0x2286, 0x237e}
	for target := range funcs {
		targets = append(targets, target)
	}
	targets = append(targets, 0x5e80, 0xce7, 0xcde, 0x8810)
	for _, target := range targets {
		t := target
		o.OnCall(oracle.Addr{Seg: cs, Off: t}, func(o *oracle.Oracle) {
			r := o.Regs()
			var b [2]byte
			binary.LittleEndian.PutUint16(b[:], t)
			trace = append(trace, b[:]...)
			trace = append(trace, pRegBytes(r)...)
			for i := uint16(0); i < 64; i++ {
				trace = append(trace, o.Byte(oracle.Addr{Seg: r.DS, Off: r.SI + i}))
			}
			trace = append(trace, o.Bytes(oracle.Addr{Seg: r.SS, Off: r.SP}, 32)...)
			trace = append(trace, o.Bytes(ga, 59)...)
			trace = append(trace, o.Bytes(ra, 2)...)
		})
	}
	for _, target := range []uint16{0x5e80, 0xce7, 0xcde, 0x8810} {
		both(cs, target, []byte{0xc3})
	}
	f := C.KiPoliticsFixture{}
	cases, audits := 0, 0
	groups := map[string]int{}
	od, cd := sha256.New(), sha256.New()
	var failure map[string]any
	check := func(c pCase) bool {
		if *only != "" && *only != c.group {
			return true
		}
		if *smoke && groups[c.group] >= 8 {
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
		q := c.queue
		if q == nil {
			q = make([]byte, 1024)
		}
		both(queueSeg, 0, q)
		scratch := c.scratch
		if scratch == nil {
			scratch = pScratch()
		}
		both(scratchSeg, 0, scratch)
		g := make([]byte, 59)
		g[15] = c.player
		g[24] = 50
		pPut(g, 13, uint16(c.player)*64)
		for i := 0; i < 3; i++ {
			pPut(g, 26+i*2, 5000)
			pPut(g, 32+i*2, 2500)
		}
		pPut(g, 48, c.cursor)
		pPut(g, 50, 0xfff0)
		pPut(g, 52, 0xfff0)
		pPut(g, 54, 400)
		pPut(g, 56, 400)
		both(cs, 0xcf0, g)
		word(cs, 0xd52, bankSeg)
		word(cs, 0xd56, queueSeg)
		word(cs, 0x987c, scratchSeg)
		both(cs, 0x31ad, []byte{123})
		initial := append([]byte(nil), table...)
		initial[0] = c.counter
		initial[1] = c.state
		both(cs, 0xecfc, initial)
		r := oracle.CallRegs{AX: c.ax, BX: c.bx, CX: 0x3456, DX: c.dx, SI: c.si, DI: c.di, BP: 0x6789, DS: bankSeg, ES: scratchSeg, SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		if c.target == 0x5358 {
			r.DS = cs
		}
		o.CallNear(oracle.Addr{Seg: cs, Off: c.target}, 0x100, r)
		input := o.Regs()
		pSet(&m, input)
		if err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x100}), oracle.Budget(400000)); err != nil {
			panic(err)
		}
		C.politics_run(&m, C.uint16_t(c.target), &f)
		got, want := pRegs(&m), o.Regs()
		if int(f.calls) > 8192 {
			panic("trace capacity")
		}
		ct := make([]byte, 0, int(f.calls)*187)
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
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.stack[0]), 32)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.globals[0]), 59)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.rng[0]), 2)...)
		}
		br, qr, sr, gr, rr := o.Bytes(ba, len(c.bank)), o.Bytes(qa, 1024), o.Bytes(sa, 1056), o.Bytes(ga, 59), o.Bytes(ra, 258)
		cases++
		groups[c.group]++
		if got != want || !bytes.Equal(trace, ct) || !bytes.Equal(mem[ba.Linear():ba.Linear()+uint32(len(br))], br) || !bytes.Equal(mem[qa.Linear():qa.Linear()+1024], qr) || !bytes.Equal(mem[sa.Linear():sa.Linear()+1056], sr) || !bytes.Equal(mem[ga.Linear():ga.Linear()+59], gr) || !bytes.Equal(mem[ra.Linear():ra.Linear()+258], rr) || mem[oracle.Addr{Seg: cs, Off: 0x31ad}.Linear()] != o.Byte(oracle.Addr{Seg: cs, Off: 0x31ad}) {
			failure = map[string]any{"group": c.group, "case": cases - 1, "input": input, "original": want, "c": got, "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct), "original_bank": hex.EncodeToString(br), "c_bank": hex.EncodeToString(mem[ba.Linear() : ba.Linear()+uint32(len(br))]), "original_queue": hex.EncodeToString(qr), "c_queue": hex.EncodeToString(mem[qa.Linear() : qa.Linear()+1024]), "original_scratch": hex.EncodeToString(sr), "c_scratch": hex.EncodeToString(mem[sa.Linear() : sa.Linear()+1056])}
			return false
		}
		for off := uint16(0x7f40); off < 0x8002; off++ {
			a := oracle.Addr{Seg: ss, Off: off}
			if mem[a.Linear()] != o.Byte(a) {
				panic("stack mismatch")
			}
		}
		od.Write(pRegBytes(want))
		od.Write(br)
		od.Write(qr)
		od.Write(sr)
		od.Write(gr)
		od.Write(rr)
		od.Write(trace)
		cd.Write(pRegBytes(got))
		cd.Write(br)
		cd.Write(qr)
		cd.Write(sr)
		cd.Write(gr)
		cd.Write(rr)
		cd.Write(ct)
		if cases%128 == 0 {
			if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory")
			}
			audits++
		}
		return true
	}
	base := func(group string, target uint16) pCase {
		return pCase{group: group, target: target, ax: 0x1204, bx: 0x2345, si: 0, di: 2, dx: 0x4567, bank: pBank(), player: 0, counter: 77, state: 91}
	}
	if *only == "" || *only == "assignment" {
		for old := 0; old < 23; old++ {
			for next := 0; next < 23; next++ {
				for _, count := range []byte{0, 1, 255} {
					c := base("assignment", 0x2ad2)
					a, b := byte(old), byte(next)
					if old == 22 {
						a = 255
					}
					if next == 22 {
						b = 255
					}
					c.ax = uint16(a)<<8 | uint16(b)
					for i := 0; i < 22; i++ {
						c.bank[i*64+0x18] = count
					}
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "writer" {
		for _, cursor := range []uint16{0, 252, 256, 1020} {
			for _, slot := range []byte{0, 1, 63, 64, 254, 255} {
				if int(cursor)+int(slot)*4 >= 1024 {
					continue
				}
				for mode := 0; mode < 4; mode++ {
					c := base("writer", 0x301c)
					c.bx = 0x7700 | uint16(slot)
					c.cursor = cursor
					c.queue = make([]byte, 1024)
					for i := 0; i < 256; i++ {
						if mode == 1 || mode == 2 && i%2 == 0 {
							pPut(c.queue, i*4, 0x0101)
						}
						if mode == 3 {
							pPut(c.queue, i*4, 0x1200)
						}
					}
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "power" {
		for _, pool := range []uint16{0, 1, 7999, 8000, 65535} {
			for _, cities := range []byte{0, 1, 7, 127, 255} {
				for _, fundHigh := range []uint16{0, 19, 20, 1000, 0xffff} {
					c := base("power", 0x3091)
					c.bx = 64
					for i := 0; i < 3; i++ {
						pPut(c.bank, 64+4+i*2, pool)
					}
					c.bank[64+0x23] = cities
					pPut(c.bank, 64+0x21, fundHigh)
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "relation" {
		for owner := 0; owner < 22; owner++ {
			for target := 0; target < 22; target++ {
				for _, entry := range []uint16{0x30f0, 0x310a} {
					for _, v := range []byte{0, 1, 127, 128, 129, 255} {
						c := base("relation", entry)
						c.si = uint16(owner * 64)
						c.di = uint16(target * 64)
						c.ax = 7
						c.bank[0x600+owner*24+target] = v
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	if *only == "" || *only == "captives" {
		for seed := 0; seed < 256; seed++ {
			for _, entry := range []uint16{0x5899, 0x5940} {
				for mode := 0; mode < 4; mode++ {
					c := base("captives", entry)
					c.si = 0x4240
					c.bank[0x4240] = 0x80
					c.bank[0x4240+0x19] = byte(mode % 3)
					c.bank[0x4240+0x1c] = byte(mode % 3)
					c.bank[0x4240+0x1d] = 1
					if mode == 3 {
						c.bank[0x4240+0x19] = 255
						c.bank[0x4240+0x1c] = 255
					}
					c.counter = byte(seed)
					c.state = byte(seed * 31)
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "generals" {
		for slot := 0; slot < 128; slot++ {
			for _, status := range []byte{0, 127, 128, 255} {
				for _, timer := range []byte{0, 1, 255} {
					c := base("generals", 0x585f)
					at := 0x4240 + slot*32
					c.bank[at] = status
					c.bank[at+0x18] = timer
					c.bank[at+0x1c] = 1
					c.bank[at+0x1d] = 255
					c.counter = byte(slot)
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "frontier" {
		for owner := 0; owner < 22; owner++ {
			for _, entry := range []uint16{0x2c52, 0x2cdf} {
				for _, value := range []byte{0, 127, 128, 255} {
					for neutral := 0; neutral < 2; neutral++ {
						c := base("frontier", entry)
						c.si = uint16(owner * 64)
						c.di = 0
						c.dx = 0
						c.ax = uint16(owner)<<8 | 1
						c.bank[owner*64] = 0x80
						c.bank[0x840] = 1
						c.bank[0x841] = byte(owner)
						c.bank[0x85c] = 1
						c.bank[0x861] = byte((owner + 1) % 22)
						if neutral == 1 {
							c.bank[0x861] = 24
						}
						c.bank[0x600+owner*24+(owner+1)%22] = value
						if entry == 0x2cdf {
							c.di = 0x840
						}
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	if *only == "" || *only == "politics" {
		for _, entry := range []uint16{0x2d3a, 0x2fb1, 0x2d58, 0x2db8, 0x2df3, 0x2e33, 0x2e89, 0x2efb, 0x2f71, 0x5990} {
			for mode := 0; mode < 8; mode++ {
				c := base("politics", entry)
				c.si = 64
				c.di = 50
				c.player = byte(mode % 3)
				c.scratch = pScratch()
				pPut(c.scratch, 48, 0x600)
				pPut(c.scratch, 50, 0x8000)
				pPut(c.scratch, 52, 0x0080)
				c.bank[64+0x19] = byte(mode % 3)
				c.bank[64+0x28] = byte(mode * 17)
				pPut(c.bank, 64+0x21, uint16(mode*300))
				for i := 0; i < 22*24; i++ {
					c.bank[0x600+i] = byte(mode * 31)
				}
				c.counter = byte(mode * 29)
				if entry == 0x5990 {
					c.si = 0x4240
					c.bank[0x425c] = c.player
				}
				if !check(c) {
					goto finish
				}
			}
		}
	}
	if *only == "" || *only == "initialize" {
		for pattern := 0; pattern < 16; pattern++ {
			c := base("initialize", 0x2bd9)
			c.queue = make([]byte, 1024)
			for i := 0; i < 256; i++ {
				pPut(c.queue, i*4, uint16(i*257+pattern))
				pPut(c.queue, i*4+2, uint16(i*31))
			}
			c.bank[0x840] = 1
			c.bank[0x841] = 0
			c.bank[0x85c] = 1
			c.bank[0x861] = 1
			c.player = byte(pattern % 3)
			c.counter = byte(pattern * 17)
			if !check(c) {
				goto finish
			}
		}
	}
	if *only == "" || *only == "scenario-monthly" {
		for scenario := 0; scenario < 4; scenario++ {
			for _, player := range []byte{0, 7, 21} {
				c := base("scenario-monthly", 0x5358)
				start := scenario*22208 + 0x80
				c.bank = append([]byte(nil), scenarios[start:start+0x5240]...)
				c.player = player
				c.counter = byte(scenario * 23)
				c.state = byte(player * 11)
				if !check(c) {
					goto finish
				}
			}
		}
	}
	if *only == "" || *only == "cooperation" {
		for _, value := range []byte{127, 128, 162, 163, 164, 165, 255, 0} {
			c := base("cooperation", 0x2e33)
			c.si, c.di, c.player = 64, 50, 0
			c.bank[64+0x19] = 2
			c.bank[0x600+24] = 0x90
			c.bank[0x600+48] = value
			c.scratch = pScratch()
			pPut(c.scratch, 50, 0)
			if !check(c) {
				goto finish
			}
		}
	}
finish:
	if failure == nil {
		if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
			panic("final memory")
		}
		audits++
	}
	report := map[string]any{"schema": "wolong-c-politics-parity-v1", "input_sha256": eh, "scenario_sha256": fmt.Sprintf("%x", sha256.Sum256(scenarios)), "routine_sha256": rh, "cases": cases, "groups": groups, "full_memory_audits": audits, "passed": failure == nil, "mismatch": failure, "original_state_sha256": hex.EncodeToString(od.Sum(nil)), "c_state_sha256": hex.EncodeToString(cd.Sum(nil)), "c_machine_code_match": false, "rng_fixture": map[string]any{"table_seed_hms": []int{12, 34, 56}, "counter_state": "Explicit per case, identically written to both CS:ECFC before execution", "rerolls": 0}, "scope": "21 politics/prisoner functions linked to full monthly rules and four scenarios; only UI/sound/redraw RET fixtures; IF/TF=0; legal pointers, indices, frontier rows and at least one faction"}
	encoded, _ := json.MarshalIndent(report, "", "  ")
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*out, encoded, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("politics original/C %d; audits %d; pass=%v\n", cases, audits, report["passed"])
	if failure != nil {
		os.Exit(1)
	}
}
