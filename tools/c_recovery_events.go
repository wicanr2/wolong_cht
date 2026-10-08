//go:build ignore

// 三十函式事件 handler 與原版/C 對照；spec/219。
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/rng.c"
#include "/repo/tools/c_recovery/economy.c"
#include "/repo/tools/c_recovery/settlement.c"
#include "/repo/tools/c_recovery/world_update.c"
#include "/repo/tools/c_recovery/politics.c"
#include "/repo/tools/c_recovery/hourly.c"
#include "/repo/tools/c_recovery/events.c"
#include "/repo/tools/c_recovery/events_fixture.h"
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
	delay, trust, answer               byte
	hourCursor                         uint16
}

func eventBank() []byte {
	b := pBank()
	for i := 0; i < 22; i++ {
		at := i * 64
		b[at] = 0x80
		b[at+1] = byte(i)
		b[at+0x19] = 255
		b[at+0x2a] = 0
		b[at+0x28] = 8
	}
	for i := 0; i < 192; i++ {
		at := 0x840 + i*32
		b[at+1] = byte(i % 3)
		b[at+0x19] = 0
	}
	for i := 0; i < 127; i++ {
		at := 0x4240 + i*32
		b[at+1] = byte(i)
		b[at+0x13] = byte(1 + i%16)
		b[at+0x17] = 0
		b[at+0x1c] = byte(i % 3)
		b[at+0x1d] = 255
		b[at+0x1e] = byte(i % 3)
		b[at+0x1f] = 1
	}
	return b
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
	funcs := map[uint16]int{0x320c: 20, 0x3220: 66, 0x3262: 71, 0x32a9: 64, 0x32e9: 62, 0x3327: 97, 0x3388: 98, 0x33ea: 19, 0x33fd: 136, 0x3485: 17, 0x34a6: 11, 0x34b1: 86, 0x351a: 12, 0x3526: 133, 0x35ab: 66, 0x35ed: 76, 0x3639: 48, 0x3669: 46, 0x3697: 45, 0x36c4: 78, 0x3712: 95, 0x3771: 103, 0x37d8: 29, 0x37f5: 59, 0x3138: 55, 0x4502: 70, 0x6a3d: 94, 0x23ff: 57, 0x2438: 33, 0x50d7: 73}
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
	entriesSeen := map[string]int{}
	targets := []uint16{0x5358, 0x53c6, 0x5456, 0x548f, 0x5538, 0x5547, 0x54fc, 0x55ec, 0x5609, 0x563b, 0x5828, 0xece0, 0x5695, 0x55a6, 0x2fbf, 0x57fe, 0x5715, 0x578f, 0x30cb, 0x3119, 0x22db, 0x2286, 0x237e}
	for target := range funcs {
		targets = append(targets, target)
	}
	targets = append(targets, 0x5e80, 0xce7, 0xcde, 0x8810, 0x2f5, 0x87ff, 0x1cb1)
	handlers := []uint16{0x320c, 0x3220, 0x3262, 0x32a9, 0x32e9, 0x3327, 0x3388, 0x33ea, 0x3485, 0x3496, 0x34a6, 0x34b1, 0x3507}
	targets = append(targets, handlers...)
	targets = append(targets, 0x310a, 0x301c, 0x2ad2, 0x3091, 0x31ae, 0x3e11, 0x3e65, 0x3e8e, 0x5673, 0x3496, 0x3507, 0x3dc9, 0x38c7, 0x38e6, 0x39e8, 0x2078, 0x20d6, 0x3c3d, 0x45f8, 0x4236, 0x5e60)
	seenTargets := map[uint16]bool{}
	for _, target := range targets {
		if seenTargets[target] {
			continue
		}
		seenTargets[target] = true
		t := target
		o.OnCall(oracle.Addr{Seg: cs, Off: t}, func(o *oracle.Oracle) {
			if _, ok := funcs[t]; ok {
				entriesSeen[fmt.Sprintf("sub_%X", 0x10000+uint32(t))]++
			}
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
			trace = append(trace, o.Byte(oracle.Addr{Seg: cs, Off: 0x31ad}))
		})
	}
	for _, target := range []uint16{0x5e80, 0xce7, 0xcde, 0x8810, 0x2f5, 0x87ff, 0x1cb1, 0x38c7, 0x38e6, 0x39e8, 0x2078, 0x20d6, 0x45f8, 0x4236, 0x5e60} {
		both(cs, target, []byte{0xc3})
	}
	both(cs, 0x3c3d, []byte{0x2e, 0xa0, 0x00, 0x70, 0xc3})

	f := C.KiHourlyFixture{}
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
		g[16] = c.trust
		pPut(g, 44, c.hourCursor)
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
		both(cs, 0x31ad, []byte{c.delay})
		both(cs, 0x7000, []byte{c.answer})
		initial := append([]byte(nil), table...)
		initial[0] = c.counter
		initial[1] = c.state
		both(cs, 0xecfc, initial)
		r := oracle.CallRegs{AX: c.ax, BX: c.bx, CX: 0x3456, DX: c.dx, SI: c.si, DI: c.di, BP: 0x6789, DS: bankSeg, ES: scratchSeg, SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		if c.target == 0x5358 || c.target == 0x3e11 || c.target == 0x31ae {
			r.DS = cs
		}
		o.CallNear(oracle.Addr{Seg: cs, Off: c.target}, 0x100, r)
		input := o.Regs()
		pSet(&m, input)
		if err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x100}), oracle.Budget(400000)); err != nil {
			panic(err)
		}
		C.events_run(&m, C.uint16_t(c.target), &f)
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
			ct = append(ct, byte(s.cadence))
		}
		br, qr, sr, gr, rr := o.Bytes(ba, len(c.bank)), o.Bytes(qa, 1024), o.Bytes(sa, 1056), o.Bytes(ga, 59), o.Bytes(ra, 258)
		cases++
		groups[c.group]++
		if got != want || !bytes.Equal(trace, ct) || !bytes.Equal(mem[ba.Linear():ba.Linear()+uint32(len(br))], br) || !bytes.Equal(mem[qa.Linear():qa.Linear()+1024], qr) || !bytes.Equal(mem[sa.Linear():sa.Linear()+1056], sr) || !bytes.Equal(mem[ga.Linear():ga.Linear()+59], gr) || !bytes.Equal(mem[ra.Linear():ra.Linear()+258], rr) || mem[oracle.Addr{Seg: cs, Off: 0x31ad}.Linear()] != o.Byte(oracle.Addr{Seg: cs, Off: 0x31ad}) {
			failure = map[string]any{"group": c.group, "case": cases - 1, "input": input, "original": want, "c": got, "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct), "original_bank": hex.EncodeToString(br), "c_bank": hex.EncodeToString(mem[ba.Linear() : ba.Linear()+uint32(len(br))]), "original_queue": hex.EncodeToString(qr), "c_queue": hex.EncodeToString(mem[qa.Linear() : qa.Linear()+1024]), "original_scratch": hex.EncodeToString(sr), "c_scratch": hex.EncodeToString(mem[sa.Linear() : sa.Linear()+1056])}
			failure["original_globals"] = hex.EncodeToString(gr)
			failure["c_globals"] = hex.EncodeToString(mem[ga.Linear() : ga.Linear()+59])
			failure["original_cadence"] = o.Byte(oracle.Addr{Seg: cs, Off: 0x31ad})
			failure["c_cadence"] = mem[oracle.Addr{Seg: cs, Off: 0x31ad}.Linear()]
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
		cd.Write(mem[ba.Linear() : ba.Linear()+uint32(len(br))])
		cd.Write(mem[qa.Linear() : qa.Linear()+uint32(len(qr))])
		cd.Write(mem[sa.Linear() : sa.Linear()+uint32(len(sr))])
		cd.Write(mem[ga.Linear() : ga.Linear()+uint32(len(gr))])
		cd.Write(mem[ra.Linear() : ra.Linear()+uint32(len(rr))])
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
		return pCase{group: group, target: target, ax: 0x1234, bx: 0x2345, si: 64, di: 128, dx: 0x0201, bank: eventBank(), player: 0, counter: 77, state: 91, delay: 7, trust: 255, hourCursor: 64}
	}

	// Dense gates, with adjacent raw bytes and all model flags preserved.
	for slot := 0; slot < 22; slot++ {
		for v := 0; v < 256; v++ {
			c := base("presence", 0x351a)
			c.ax = uint16(slot)<<8 | 0x5a
			c.bank[slot*64] = byte(v)
			if !check(c) {
				goto finish
			}
		}
	}
	for _, t := range []uint16{0x3639, 0x3669} {
		for v := 0; v < 256; v++ {
			for _, w := range []byte{0, 1, 63, 64, 99, 100, 127, 128, 255} {
				c := base("relation", t)
				c.ax = 0x0100
				c.bank[0x601] = byte(v)
				c.bank[0x618] = w
				if !check(c) {
					goto finish
				}
			}
		}
	}
	for _, t := range []uint16{0x3639, 0x3669, 0x3697} {
		for _, a := range []uint16{0x0000, 0x1514, 0x1801, 0x0118} {
			c := base("relation", t)
			c.ax = a
			if !check(c) {
				goto finish
			}
		}
	}
	for _, t := range []uint16{0x35ab, 0x3526} {
		for _, old := range []byte{0, 1, 21, 23, 24, 35, 36, 255} {
			for _, to := range []byte{0, 1, 2, 21, 24} {
				for strength := 0; strength < 3; strength++ {
					c := base("target", t)
					c.dx = uint16(to)
					c.bank[64+0x19] = 255
					at := int(to) * 64
					c.bank[at+0x19] = old
					pPut(c.bank, at+8, uint16(strength*2000))
					pPut(c.bank, 64+8, 2000)
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	for seed := 0; seed < 256; seed++ {
		for _, power := range []byte{0, 1, 8, 16, 255} {
			for _, official := range []byte{0, 255} {
				c := base("official", 0x3771)
				c.counter = byte(seed)
				c.state = byte(seed * 37)
				c.bank[64+0x2a] = official
				c.bank[0x4253] = power
				if !check(c) {
					goto finish
				}
			}
		}
	}
	for _, t := range []uint16{0x36c4, 0x3712} {
		for _, rel := range []byte{0, 1, 15, 16, 29, 30, 59, 60, 63, 64, 99, 100, 127, 128, 129, 255} {
			for _, policy := range []byte{0, 1, 8, 30, 127, 255} {
				for _, official := range []byte{0, 255} {
					for seed := 0; seed < 8; seed++ {
						c := base("diplomacy", t)
						c.bank[64+0x28] = policy
						c.bank[64+0x2a] = official
						c.bank[0x61a] = rel
						c.bank[0x630+1] = byte(255 - int(rel))
						c.counter = byte(seed * 31)
						c.state = byte(seed * 37)
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	for _, t := range []uint16{0x23ff, 0x2438} {
		for occupied := 0; occupied <= 16; occupied++ {
			for variant := 0; variant < 4; variant++ {
				c := base("disaster", t)
				c.ax = uint16(variant + 1)
				c.dx = 0x1234
				c.bx = 0x3456
				for k := 0; k < 16; k++ {
					at := 0x2040 + k*16
					if k < occupied {
						c.bank[at] = 0x80
					}
					pPut(c.bank, at+2, uint16(0x1234+(k%2)*variant))
					pPut(c.bank, at+4, uint16(0x3456+(k%3)*variant))
				}
				if !check(c) {
					goto finish
				}
			}
		}
	}
	for occupied := 0; occupied <= 16; occupied++ {
		for code := 0; code < 4; code++ {
			for seed := 0; seed < 256; seed++ {
				c := base("disaster-event", 0x34b1)
				c.ax = uint16(code)<<8 | 12
				c.dx = 0x840
				c.counter = byte(seed)
				c.state = byte(seed * 37)
				for k := 0; k < occupied; k++ {
					c.bank[0x2040+k*16] = 0x80
				}
				if !check(c) {
					goto finish
				}
			}
		}
	}
	for _, t := range []uint16{0x6a3d, 0x4502, 0x33fd} {
		for faction := 0; faction < 22; faction++ {
			for variant := 0; variant < 32; variant++ {
				c := base("capital", t)
				c.si = uint16(faction * 64)
				c.bx = uint16(faction)
				c.ax = uint16(faction)
				if t == 0x4502 {
					c.ax = 0x0302
				}
				if t == 0x33fd {
					c.ax = 0x860
				}
				c.player = byte(variant % 22)
				c.bank[faction*64+3] = byte(variant % 3)
				for k := 0; k < 192; k++ {
					at := 0x840 + k*32
					c.bank[at] = byte(variant & 31)
					c.bank[at+1] = byte((k + faction) % 22)
					c.bank[at+0x16] = byte((k + variant) % 16)
					pPut(c.bank, at+0xe, uint16((k*7919+variant)%65536))
				}
				for k := 0; k < 127; k++ {
					at := 0x2240 + k*64
					c.bank[at] = byte(127 + (k+variant)%3)
					c.bank[at+1] = byte((k + faction) % 22)
					c.bank[at+0x20] = byte(2 + (k+variant)%3)
					pPut(c.bank, at+0x14, uint16((2+(k+variant)%3)*8))
				}
				if !check(c) {
					goto finish
				}
			}
		}
	}
	for _, t := range []uint16{0x37f5, 0x3138, 0x37d8, 0x35ed, 0x50d7} {
		for variant := 0; variant < 128; variant++ {
			c := base("cleanup", t)
			c.ax = uint16(variant % 4)
			c.dx = 500
			c.di = 128
			if t == 0x37f5 {
				c.ax = 64
			}
			if t == 0x50d7 {
				c.di = 0x4240
			}
			for k := 0; k < 127; k++ {
				at := 0x4240 + k*32
				c.bank[at+0x13] = byte((k + variant) % 256)
				c.bank[at+0x17] = byte(k % 3)
				c.bank[at+0x1c] = byte((k + variant) % 3)
				c.bank[at+0x1d] = byte((k + variant + 1) % 3)
			}
			if !check(c) {
				goto finish
			}
		}
	}
	for _, t := range []uint16{0x320c, 0x3220, 0x3262, 0x32a9, 0x32e9, 0x3327, 0x3388, 0x33ea, 0x33fd, 0x3485, 0x34a6, 0x34b1} {
		for _, player := range []byte{0, 1, 21} {
			for _, present := range []byte{0, 127, 128, 255} {
				for answer := 0; answer < 4; answer++ {
					for _, official := range []byte{255, 0, 1} {
						for _, seed := range []byte{0, 77, 255} {
							c := base("handlers", t)
							c.ax = 0x0101
							c.dx = 0x0201
							c.player = player
							c.answer = byte(answer)
							c.counter = seed
							c.state = seed * 37
							c.bank[64] = present
							c.bank[64+0x2a] = official
							c.bank[0x840+0x19] = official
							if t == 0x32a9 || t == 0x3485 {
								c.ax = 1
							}
							if t == 0x33fd {
								c.ax = 0x860
							}
							if t == 0x34b1 {
								c.dx = 0x840
							}
							if !check(c) {
								goto finish
							}
						}
					}
				}
			}
		}
	}
	for scenario := 0; scenario < 4; scenario++ {
		for _, player := range []byte{0, 7, 21} {
			for code := 1; code <= 13; code++ {
				for _, seed := range []byte{0, 77, 255} {
					for answer := 0; answer < 4; answer++ {
						c := base("scenario-dispatch", 0x31ae)
						start := scenario*22208 + 0x80
						c.bank = append([]byte(nil), scenarios[start:start+0x5240]...)
						c.player = player
						c.counter = seed
						c.state = seed * 37
						c.delay = 1
						c.answer = byte(answer)
						c.queue = make([]byte, 1024)
						pPut(c.queue, 0, uint16(0x100|code))
						param := uint16(0x0201)
						if code == 12 {
							param = 0x840
						}
						if code == 13 {
							param = 0xffff
						}
						pPut(c.queue, 2, param)
						if !check(c) {
							goto finish
						}
					}
				}
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
	report := map[string]any{"schema": "wolong-c-events-parity-v1", "input_sha256": eh, "scenario_sha256": fmt.Sprintf("%x", sha256.Sum256(scenarios)), "routine_sha256": rh, "cases": cases, "groups": groups, "full_memory_audits": audits, "passed": failure == nil, "mismatch": failure, "original_state_sha256": hex.EncodeToString(od.Sum(nil)), "c_state_sha256": hex.EncodeToString(cd.Sum(nil)), "c_machine_code_match": false, "rng_fixture": map[string]any{"table_seed_hms": []int{12, 34, 56}, "counter_state": "Explicit per case, identically written to both CS:ECFC before execution", "rerolls": 0}, "scope": "30 event functions; original event 8 fall-through; real politics/disaster/RNG/writer; modal/amount/battle/UI/sound/redraw fixtures; explicit choice 0..3; raw target byte boundary; IF/TF=0"}
	report["entries_seen"] = entriesSeen
	encoded, _ := json.MarshalIndent(report, "", "  ")
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*out, encoded, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("events original/C %d; audits %d; pass=%v\n", cases, audits, report["passed"])
	if failure != nil {
		os.Exit(1)
	}
}
