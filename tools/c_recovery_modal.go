//go:build ignore

// 十七函式外交／金額視窗與原版/C 對照；spec/220。
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
#include "/repo/tools/c_recovery/modal.c"
#include "/repo/tools/c_recovery/modal_fixture.h"
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
	group                                                  string
	target, ax, bx, si, di, dx, cursor                     uint16
	bank, queue, scratch                                   []byte
	player, counter, state                                 byte
	delay, trust, answer                                   byte
	hourCursor                                             uint16
	amount, farX, farY, originX, originY                   uint16
	selectWait, numberWait, initial, query, uiFlags, music byte
	message                                                uint16
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
	funcs := map[uint16]int{0x2078: 94, 0x20d6: 123, 0x38c7: 31, 0x38e6: 28, 0x3902: 230, 0x39e8: 288, 0x3c3d: 92, 0x3b7e: 43, 0x3d09: 60, 0x3d45: 35, 0x3c99: 39, 0x3cdc: 45, 0x9321: 21, 0x87ff: 17, 0x3d68: 41, 0x1d46: 72, 0x2216: 21}
	rh := map[string]string{}
	for target, size := range funcs {
		b := raw[int(target)+0x200 : int(target)+0x200+size]
		runtimeBytes := append([]byte(nil), b...)
		relocCount, relocTable := int(binary.LittleEndian.Uint16(raw[6:])), int(binary.LittleEndian.Uint16(raw[24:]))
		for i := 0; i < relocCount; i++ {
			at := relocTable + i*4
			file := 512 + int(binary.LittleEndian.Uint16(raw[at:])) + 16*int(binary.LittleEndian.Uint16(raw[at+2:]))
			if file >= int(target)+512 && file < int(target)+512+size {
				local := file - int(target) - 512
				v := binary.LittleEndian.Uint16(runtimeBytes[local:])
				binary.LittleEndian.PutUint16(runtimeBytes[local:], v+cs)
			}
		}
		if !bytes.Equal(o.Bytes(oracle.Addr{Seg: cs, Off: target}, size), runtimeBytes) {
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
	farSeg := o.Word(oracle.Addr{Seg: cs, Off: 0x2082})
	if farSeg != cs+0x1000 {
		panic("runtime FAR segment")
	}
	uiBytes := func() []byte {
		b := append([]byte(nil), o.Bytes(oracle.Addr{Seg: cs, Off: 0x7000}, 32)...)
		b = append(b, o.Bytes(oracle.Addr{Seg: farSeg, Off: 0x7004}, 8)...)
		b = append(b, o.Bytes(oracle.Addr{Seg: cs, Off: 0x9882}, 16)...)
		b = append(b, o.Bytes(oracle.Addr{Seg: cs, Off: 0x2239}, 2)...)
		b = append(b, o.Byte(oracle.Addr{Seg: cs, Off: 0x20e}), o.Byte(oracle.Addr{Seg: cs, Off: 0x98a6}))
		b = append(b, o.Bytes(oracle.Addr{Seg: cs, Off: 0x9874}, 4)...)
		return b
	}
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
	targets = append(targets, 0x38c7, 0x38e6, 0x39e8, 0x2078, 0x20d6, 0x3c3d, 0x3b7e, 0x3902, 0x3c99, 0x3cdc, 0x3d09, 0x3d45, 0x3d68, 0x9321, 0x87ff, 0x1d46, 0x2216, 0x9796, 0x93e9, 0x97c3, 0x7c6e, 0x75b, 0x222b, 0xe38c, 0x89a4, 0x87af, 0xc14, 0xfa37, 0x2c2, 0x241, 0xe453, 0x1f30, 0x5c58, 0xd615, 0x1cc9, 0xd66a, 0x351a, 0x3526, 0x35ab, 0x35ed, 0x3639, 0x3669, 0x3697, 0x36c4, 0x3712, 0x3771, 0x37d8, 0x37f5, 0x3138, 0x4502, 0x6a3d, 0x23ff, 0x2438, 0x50d7)
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
			trace = append(trace, uiBytes()...)
		})
	}
	for _, target := range []uint16{0x5e80, 0xce7, 0xcde, 0x8810, 0x2f5, 0x1cb1, 0x45f8, 0x4236, 0x5e60, 0x97c3, 0x75b, 0x222b, 0xe38c, 0x89a4, 0x87af, 0xc14, 0xfa37, 0x2c2, 0x241, 0x1f30, 0x5c58, 0xd615, 0x1cc9, 0xd66a} {
		both(cs, target, []byte{0xc3})
	}
	waitStub := func(value, wait uint16) []byte {
		return []byte{0x2e, 0xa1, byte(value), byte(value >> 8), 0x2e, 0x80, 0x3e, byte(wait), byte(wait >> 8), 0, 0x74, 7, 0x2e, 0xfe, 0x0e, byte(wait), byte(wait >> 8), 0xf9, 0xc3, 0xf8, 0xc3}
	}
	both(cs, 0x93e9, waitStub(0x7000, 0x7008))
	both(cs, 0x7c6e, waitStub(0x7002, 0x7009))
	both(cs, 0x9796, []byte{0x2e, 0xa0, 0x0a, 0x70, 0xc3})
	both(cs, 0xe453, []byte{0x2e, 0xa0, 0x10, 0x70, 0xc3})
	both(farSeg, 0, []byte{0x3d, 3, 0, 0x75, 10, 0x2e, 0x8b, 0x0e, 4, 0x70, 0x2e, 0x8b, 0x16, 6, 0x70, 0x3b, 0xc0, 0xcb})
	o.OnCall(oracle.Addr{Seg: farSeg, Off: 0}, func(o *oracle.Oracle) {
		r := o.Regs()
		trace = append(trace, 0, 0)
		trace = append(trace, pRegBytes(r)...)
		for i := uint16(0); i < 64; i++ {
			trace = append(trace, o.Byte(oracle.Addr{Seg: r.DS, Off: r.SI + i}))
		}
		trace = append(trace, o.Bytes(oracle.Addr{Seg: r.SS, Off: r.SP}, 32)...)
		trace = append(trace, o.Bytes(ga, 59)...)
		trace = append(trace, o.Bytes(ra, 2)...)
		trace = append(trace, o.Byte(oracle.Addr{Seg: cs, Off: 0x31ad}))
		trace = append(trace, uiBytes()...)
	})

	f := C.KiModalFixture{origin_cs: C.uint16_t(cs)}
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
		f.base.calls = 0
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
		g[4] = c.music
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
		ctrl := make([]byte, 32)
		pPut(ctrl, 0, uint16(c.answer))
		pPut(ctrl, 2, c.amount)
		ctrl[8] = c.selectWait
		ctrl[9] = c.numberWait
		ctrl[10] = c.initial
		ctrl[16] = c.query
		both(cs, 0x7000, ctrl)
		word(farSeg, 0x7004, c.farX)
		word(farSeg, 0x7006, c.farY)
		word(cs, 0x9882, c.originX)
		word(cs, 0x9884, c.originY)
		word(cs, 0x9886, c.farY)
		word(cs, 0x9888, c.farX)
		word(cs, 0x988e, c.originY)
		word(cs, 0x9890, c.originX)
		word(cs, 0x9876, bankSeg)
		both(cs, 0x98a6, []byte{c.uiFlags})
		word(ss, 0x6789, c.message)
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
		C.modal_run(&m, C.uint16_t(c.target), &f)
		got, want := pRegs(&m), o.Regs()
		if int(f.base.calls) > 8192 {
			panic("trace capacity")
		}
		ct := make([]byte, 0, int(f.base.calls)*252)
		for i := 0; i < int(f.base.calls); i++ {
			s := f.base.trace[i]
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
			ct = append(ct, C.GoBytes(unsafe.Pointer(&f.ui[i][0]), 64)...)
		}
		br, qr, sr, gr, rr := o.Bytes(ba, len(c.bank)), o.Bytes(qa, 1024), o.Bytes(sa, 1056), o.Bytes(ga, 59), o.Bytes(ra, 258)
		var actualUI [64]C.uint8_t
		C.md_observed_ui(&m, &f, &actualUI[0])
		cUI := C.GoBytes(unsafe.Pointer(&actualUI[0]), 64)
		originalUI := uiBytes()
		cases++
		groups[c.group]++
		if !bytes.Equal(originalUI, cUI) || got != want || !bytes.Equal(trace, ct) || !bytes.Equal(mem[ba.Linear():ba.Linear()+uint32(len(br))], br) || !bytes.Equal(mem[qa.Linear():qa.Linear()+1024], qr) || !bytes.Equal(mem[sa.Linear():sa.Linear()+1056], sr) || !bytes.Equal(mem[ga.Linear():ga.Linear()+59], gr) || !bytes.Equal(mem[ra.Linear():ra.Linear()+258], rr) || mem[oracle.Addr{Seg: cs, Off: 0x31ad}.Linear()] != o.Byte(oracle.Addr{Seg: cs, Off: 0x31ad}) {
			failure = map[string]any{"group": c.group, "case": cases - 1, "input": input, "original": want, "c": got, "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct), "original_bank": hex.EncodeToString(br), "c_bank": hex.EncodeToString(mem[ba.Linear() : ba.Linear()+uint32(len(br))]), "original_queue": hex.EncodeToString(qr), "c_queue": hex.EncodeToString(mem[qa.Linear() : qa.Linear()+1024]), "original_scratch": hex.EncodeToString(sr), "c_scratch": hex.EncodeToString(mem[sa.Linear() : sa.Linear()+1056])}
			failure["original_globals"] = hex.EncodeToString(gr)
			failure["c_globals"] = hex.EncodeToString(mem[ga.Linear() : ga.Linear()+59])
			failure["original_ui"] = hex.EncodeToString(originalUI)
			failure["c_ui"] = hex.EncodeToString(cUI)
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
		od.Write(originalUI)
		cd.Write(pRegBytes(got))
		cd.Write(mem[ba.Linear() : ba.Linear()+uint32(len(br))])
		cd.Write(mem[qa.Linear() : qa.Linear()+uint32(len(qr))])
		cd.Write(mem[sa.Linear() : sa.Linear()+uint32(len(sr))])
		cd.Write(mem[ga.Linear() : ga.Linear()+uint32(len(gr))])
		cd.Write(mem[ra.Linear() : ra.Linear()+uint32(len(rr))])
		cd.Write(ct)
		cd.Write(cUI)
		if cases%128 == 0 {
			if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory")
			}
			audits++
		}
		return true
	}
	base := func(group string, target uint16) pCase {
		return pCase{group: group, target: target, ax: 0x1234, bx: 0x2345, si: 0x4240, di: 128, dx: 0x0201, bank: eventBank(), player: 0, counter: 77, state: 91, delay: 7, trust: 255, hourCursor: 64, amount: 500, farX: 639, farY: 399, originX: 24, originY: 40, music: 1, message: 0x3456}
	}

	for _, t := range []uint16{0x2078, 0x20d6} {
		for _, x := range []uint16{0, 1, 639, 640, 65535} {
			for _, y := range []uint16{0, 1, 399, 400, 65535} {
				for _, a := range []uint16{0, 1, 639, 640, 65535} {
					for _, b := range []uint16{0, 1, 399, 400, 65535} {
						c := base("window", t)
						c.originX = x
						c.originY = y
						c.farX = a
						c.farY = b
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	for _, t := range []uint16{0x3c99, 0x3cdc, 0x3d09, 0x3d45, 0x3d68, 0x9321, 0x87ff, 0x2216, 0x1d46} {
		for variant := 0; variant < 256; variant++ {
			c := base("helper", t)
			c.ax = uint16(variant)
			c.query = byte(variant)
			c.uiFlags = byte(variant)
			c.music = byte(variant)
			c.bank[0x425e] = byte(variant)
			if !check(c) {
				goto finish
			}
		}
	}
	for init := 0; init < 4; init++ {
		for choice := 0; choice < 4; choice++ {
			for wait := 0; wait < 4; wait++ {
				for _, message := range []uint16{0, 500, 65535} {
					c := base("selector", 0x3b7e)
					c.initial = byte(init)
					c.answer = byte(choice)
					c.selectWait = byte(wait)
					c.ax = 3
					c.message = message
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	for answer := 0; answer < 4; answer++ {
		for _, high := range []uint16{0, 1, 2, 3, 255} {
			for _, player := range []byte{0, 1} {
				for _, trust := range []byte{0, 29, 30, 31, 255} {
					c := base("reply", 0x3c3d)
					c.si = 64
					c.ax = high<<8 | uint16(answer)
					c.player = player
					c.trust = trust
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	for mode := 0; mode < 4; mode++ {
		for _, high := range []uint16{0, 1, 3} {
			for _, prior := range []uint16{0, 1, 500, 30000, 65535} {
				for choice := 0; choice < 4; choice++ {
					for _, amount := range []uint16{0, 1, 499, 500, 501, 30000, 65535} {
						for _, retry := range []byte{0, 1, 3} {
							for seed := 0; seed < 8; seed++ {
								c := base("diplomacy", 0x3902)
								c.ax = high<<8 | uint16(mode)
								c.dx = prior
								c.answer = byte(choice)
								c.amount = amount
								c.numberWait = retry
								c.selectWait = byte(seed % 4)
								c.initial = byte(seed % 4)
								c.counter = byte(seed * 31)
								c.state = byte(seed * 37)
								c.trust = []byte{0, 128, 255}[seed%3]
								if !check(c) {
									goto finish
								}
							}
						}
					}
				}
			}
		}
	}
	for seed := 0; seed < 256; seed++ {
		for _, trust := range []byte{0, 128, 255} {
			for _, amount := range []uint16{0, 1, 500, 65535} {
				c := base("trust-rng", 0x3902)
				c.ax = 0x103
				c.dx = 500
				c.answer = 1
				c.amount = amount
				c.trust = trust
				c.counter = byte(seed)
				c.state = byte(seed * 37)
				if !check(c) {
					goto finish
				}
			}
		}
	}
	for _, prior := range []uint16{0, 1, 499, 500, 501, 30000, 65535} {
		for _, amount := range []uint16{0, 1, 499, 500, 501, 30000, 32768, 65535} {
			for choice := 0; choice < 4; choice++ {
				for _, retry := range []byte{0, 1, 3} {
					for _, variant := range []byte{0, 3, 7} {
						c := base("amount", 0x39e8)
						c.ax = prior
						c.amount = amount
						c.answer = byte(choice)
						c.numberWait = retry
						c.selectWait = retry
						c.bank[0x425e] = variant
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	for _, t := range []uint16{0x38c7, 0x38e6} {
		for seed := 0; seed < 8; seed++ {
			for _, prior := range []uint16{0, 1, 500, 30000, 65535} {
				for choice := 0; choice < 4; choice++ {
					for _, amount := range []uint16{0, 1, 499, 500, 501, 30000, 65535} {
						c := base("wrapper", t)
						c.ax = 0x103
						c.dx = prior
						c.amount = amount
						c.answer = byte(choice)
						c.counter = byte(seed * 31)
						c.state = byte(seed * 37)
						c.numberWait = byte(seed % 4)
						c.selectWait = byte(seed % 4)
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	for scenario := 0; scenario < 4; scenario++ {
		for _, player := range []byte{0, 7, 21} {
			for _, t := range []uint16{0x3220, 0x3262, 0x32a9, 0x32e9, 0x3327, 0x3388} {
				for choice := 0; choice < 4; choice++ {
					for _, amount := range []uint16{0, 500, 65535} {
						for wait := 0; wait < 2; wait++ {
							c := base("scenario-modal", t)
							start := scenario*22208 + 0x80
							c.bank = append([]byte(nil), scenarios[start:start+0x5240]...)
							c.ax = 0x0101
							c.dx = 0x0201
							if t == 0x32a9 {
								c.ax = 1
							}
							c.player = player
							c.answer = byte(choice)
							c.amount = amount
							c.numberWait = byte(wait)
							c.selectWait = byte(wait)
							c.counter = byte(scenario * 23)
							c.state = byte(player * 11)
							if !check(c) {
								goto finish
							}
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
	report := map[string]any{"schema": "wolong-c-modal-parity-v1", "input_sha256": eh, "scenario_sha256": fmt.Sprintf("%x", sha256.Sum256(scenarios)), "routine_sha256": rh, "cases": cases, "groups": groups, "full_memory_audits": audits, "passed": failure == nil, "mismatch": failure, "original_state_sha256": hex.EncodeToString(od.Sum(nil)), "c_state_sha256": hex.EncodeToString(cd.Sum(nil)), "c_machine_code_match": false, "rng_fixture": map[string]any{"table_seed_hms": []int{12, 34, 56}, "counter_state": "Explicit per case, identically written to both CS:ECFC before execution", "rerolls": 0}, "scope": "17 modal functions; exact MZ/runtime FAR relocation and return words; real CF retry, FLAGS save, RNG/trust, amount, helper and code patch; deterministic input/device leaves, graphics/audio primitives; persistent original stack slots; IF/TF=0"}
	report["entries_seen"] = entriesSeen
	report["runtime_load_paragraph"] = cs
	report["far_runtime_segment"] = farSeg
	report["trace_entry_size"] = 252
	encoded, _ := json.MarshalIndent(report, "", "  ")
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*out, encoded, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("modal original/C %d; audits %d; pass=%v\n", cases, audits, report["passed"])
	if failure != nil {
		os.Exit(1)
	}
}
