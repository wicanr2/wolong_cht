//go:build ignore

// 十函式熱區與原版/C 真實 pixel query；spec/222。
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
#include "/repo/tools/c_recovery/numeric.c"
#include "/repo/tools/c_recovery/hotspot.c"
#include "/repo/tools/c_recovery/hotspot_fixture.h"
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
	"github.com/wicanr2/wolong_cht/internal/state"
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
	cap, saved, pointerDX, pointerBX                       uint16
	keys                                                   []byte
	cx, mapBase                                            uint16
	mapBytes, flagBytes                                    []byte
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

func numericKeys(plan string) []byte {
	var b []byte
	for _, ch := range plan {
		raw, cancel := byte(0), byte(0)
		switch {
		case ch >= '0' && ch <= '9':
			raw = 0x52 + byte(ch-'0')
		case ch == 'F':
			raw = 0x60
		case ch == 'X':
			cancel = 1
		case ch == 'M':
			raw = 0x5f
		case ch == 'Z':
			raw = 0x5c
		case ch == 'C':
			raw = 0x5e
		case ch == 'B':
			raw = 0x5d
		case ch == '!':
			raw = 0x51
		}
		b = append(b, raw, cancel)
	}
	return b
}
func numericEdit(raw byte) (state.AmountEdit, int, bool) {
	switch {
	case raw >= 0x52 && raw <= 0x5b:
		return state.AmountAppendDigit, int(raw - 0x52), true
	case raw == 0x5c:
		return state.AmountAppendHundred, 0, true
	case raw == 0x5d:
		return state.AmountDeleteDigit, 0, true
	case raw == 0x5e:
		return state.AmountClear, 0, true
	case raw == 0x5f:
		return state.AmountSetMax, 0, true
	case raw == 0x60:
		return state.AmountFinishInput, 0, true
	}
	return 0, 0, false
}
func numericGo(c pCase) (int, bool) {
	if c.target == 0x7c6e {
		v := 0
		keys := c.keys
		if keys == nil {
			keys = []byte{0x60, 0}
		}
		for i := 0; i < len(keys); i += 2 {
			if keys[i+1] != 0 {
				break
			}
			e, d, ok := numericEdit(keys[i])
			if !ok {
				continue
			}
			v, _ = state.EditAmountValue(v, int(c.ax), e, d)
			if keys[i] == 0x60 {
				break
			}
		}
		return v, true
	}
	var e state.AmountEdit
	d := 0
	switch c.target {
	case 0x7da5:
		e = state.AmountAppendDigit
		d = int(byte(c.ax))
	case 0x7dc3:
		e = state.AmountAppendHundred
	case 0x7ddd:
		e = state.AmountDeleteDigit
	case 0x7dea:
		e = state.AmountFinishInput
	case 0x7dec:
		e = state.AmountSetMax
	case 0x7df1:
		e = state.AmountClear
	default:
		return 0, false
	}
	if c.si > c.cap {
		return 0, false
	}
	v, ok := state.EditAmountValue(int(c.si), int(c.cap), e, d)
	if !ok {
		panic("Go legal action")
	}
	return v, true
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
	mapSeg, flagSeg := uint16(0x6800), uint16(0x7800)
	funcs := map[uint16]int{0xe3c0: 23, 0xe3d7: 68, 0xe41b: 56, 0xe453: 38, 0x895d: 71, 0x89de: 18, 0xd5d4: 65, 0xc14: 76, 0xc60: 23, 0xc77: 53}
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
	targets = append(targets, 0x1b4, 0x1db, 0x7c6e, 0x7d0d, 0x7d47, 0x7d5f, 0x7da5, 0x7dc3, 0x7ddd, 0x7dea, 0x7dec, 0x7df1, 0x93e9, 0x67cd, 0x67e6, 0x6806, 0x6826, 0x8853, 0x62f, 0x21e7, 0x895d, 0xe3d7, 0x9409)
	inputBytes := func(r oracle.Regs) []byte {
		b := append([]byte(nil), o.Bytes(oracle.Addr{Seg: cs, Off: 0x7020}, 8)...)
		b = append(b, o.Bytes(oracle.Addr{Seg: cs, Off: 0x1d5}, 6)...)
		b = append(b, o.Byte(oracle.Addr{Seg: cs, Off: 0x9441}), o.Byte(oracle.Addr{Seg: cs, Off: 0x9443}))
		b = append(b, o.Bytes(oracle.Addr{Seg: farSeg, Off: 0x7010}, 6)...)
		b = append(b, o.Bytes(oracle.Addr{Seg: r.SS, Off: r.BP}, 2)...)
		return b
	}
	targets = append(targets, 0xe3c0, 0xe3d7, 0xe41b, 0xe453, 0x895d, 0x89de, 0xd5d4, 0xc14, 0xc60, 0xc77, 0xf9b0)
	hotspotBytes := func() []byte {
		b := append([]byte(nil), o.Bytes(oracle.Addr{Seg: cs, Off: 0xe479}, 4)...)
		b = append(b, o.Bytes(oracle.Addr{Seg: cs, Off: 0xd84e}, 2)...)
		b = append(b, o.Bytes(oracle.Addr{Seg: cs, Off: 0xd4a}, 2)...)
		seg := o.Word(oracle.Addr{Seg: cs, Off: 0xe479})
		off := o.Word(oracle.Addr{Seg: cs, Off: 0xe47b})
		for i := uint16(0); i < 8; i++ {
			b = append(b, o.Byte(oracle.Addr{Seg: seg, Off: off + i}))
		}
		return b
	}
	seenTargets := map[uint16]bool{}
	targets = append(targets, 0x38c7, 0x38e6, 0x39e8, 0x2078, 0x20d6, 0x3c3d, 0x3b7e, 0x3902, 0x3c99, 0x3cdc, 0x3d09, 0x3d45, 0x3d68, 0x9321, 0x87ff, 0x1d46, 0x2216, 0x9796, 0x93e9, 0x97c3, 0x7c6e, 0x75b, 0x222b, 0xe38c, 0x89a4, 0x87af, 0xfa37, 0x2c2, 0x241, 0xe453, 0x1f30, 0x5c58, 0xd615, 0x1cc9, 0xd66a, 0x351a, 0x3526, 0x35ab, 0x35ed, 0x3639, 0x3669, 0x3697, 0x36c4, 0x3712, 0x3771, 0x37d8, 0x37f5, 0x3138, 0x4502, 0x6a3d, 0x23ff, 0x2438, 0x50d7)
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
			trace = append(trace, inputBytes(r)...)
			trace = append(trace, hotspotBytes()...)
		})
	}
	for _, target := range []uint16{0x5e80, 0xce7, 0xcde, 0x8810, 0x2f5, 0x1cb1, 0x45f8, 0x4236, 0x5e60, 0x97c3, 0x75b, 0x222b, 0xe38c, 0x89a4, 0x87af, 0xfa37, 0x2c2, 0x241, 0x1f30, 0x5c58, 0xd615, 0x1cc9, 0xd66a, 0x9796, 0x8853, 0x62f, 0xf9b0} {
		both(cs, target, []byte{0xc3})
	}
	// Poll consumes exactly one predeclared event; query returns its latched raw key.
	both(cs, 0x21e7, []byte{0x53, 0x2e, 0x8b, 0x1e, 0x20, 0x70, 0xd1, 0xe3, 0xd1, 0xe3, 0x2e, 0x03, 0x1e, 0x20, 0x70, 0x2e, 0x8b, 0x8f, 0, 0x71, 0x2e, 0x8b, 0x97, 2, 0x71, 0x2e, 0x8a, 0x87, 4, 0x71, 0x2e, 0xa2, 0x23, 0x70, 0x2e, 0xff, 0x06, 0x20, 0x70, 0x3c, 1, 0xf5, 0x5b, 0xc3})

	both(cs, 0x9409, []byte{0x2e, 0xa1, 0, 0x70, 0x2e, 0x80, 0x3e, 0x23, 0x70, 1, 0x75, 4, 0x2e, 0xa1, 0x0c, 0x70, 0x2e, 0x80, 0x3e, 8, 0x70, 0, 0x74, 7, 0x2e, 0xfe, 0x0e, 8, 0x70, 0xf9, 0xc3, 0xf8, 0xc3})
	both(farSeg, 0, []byte{0x3d, 3, 0, 0x75, 10, 0x2e, 0x8b, 0x0e, 4, 0x70, 0x2e, 0x8b, 0x16, 6, 0x70, 0x3d, 2, 0, 0x75, 14, 0x2e, 0xa1, 0x10, 0x70, 0x2e, 0x8b, 0x16, 0x12, 0x70, 0x2e, 0x8b, 0x1e, 0x14, 0x70, 0x3b, 0xc0, 0xcb})
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
		trace = append(trace, inputBytes(r)...)
		trace = append(trace, hotspotBytes()...)
	})

	f := C.KiHotspotFixture{base: C.KiNumericFixture{base: C.KiModalFixture{origin_cs: C.uint16_t(cs)}}}
	cases, audits := 0, 0
	goCases := 0
	goGroups := map[string]int{}
	groups := map[string]int{}
	od, cd := sha256.New(), sha256.New()
	var failure map[string]any
	defaultBank := eventBank()
	defaultMap := make([]byte, 65536)
	defaultFlags := make([]byte, 65536)
	for i := range defaultMap {
		defaultMap[i] = byte(i*37 + i/80*13 + 73)
		defaultFlags[i] = byte(i*11 + 91)
	}
	check := func(c pCase) bool {
		if *only != "" && *only != c.group {
			return true
		}
		if *smoke && groups[c.group] >= 8 {
			return true
		}
		trace = nil
		f.base.base.base.calls = 0
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
		word(ss, 0x6789, c.cap)
		maps := c.mapBytes
		if maps == nil {
			maps = defaultMap
		}
		flagMap := c.flagBytes
		if flagMap == nil {
			flagMap = defaultFlags
		}
		both(mapSeg, 0, maps)
		both(flagSeg, 0, flagMap)
		word(cs, 0xe479, mapSeg)
		word(cs, 0xe47b, c.mapBase)
		word(cs, 0xd84e, flagSeg)
		word(cs, 0xd4a, bankSeg)
		keyBuffer := make([]byte, 1024)
		keyStream := c.keys
		if keyStream == nil {
			keyStream = []byte{0x60, 0}
		}
		if len(keyStream)%2 != 0 || len(keyStream) > len(keyBuffer) {
			panic("input stream")
		}
		anchorX, anchorY := uint16(88), uint16(184)
		if c.target == 0x7c6e {
			anchorX, anchorY = c.dx, c.bx
		}
		if c.target == 0x67cd || c.target == 0x67e6 || c.target == 0x6806 || c.target == 0x6826 {
			anchorX = 296
		}
		tableKeys := raw[0x7f93:0x7fa5]
		events := []byte{}
		for i := 0; i < len(keyStream); i += 2 {
			x, y := uint16(0), uint16(0)
			if keyStream[i+1] == 0 {
				for cell, v := range tableKeys {
					if v == keyStream[i] {
						x = anchorX + uint16(cell%6*16)
						y = anchorY + 16 + uint16(cell/6*16)
						break
					}
				}
			}
			var e [5]byte
			binary.LittleEndian.PutUint16(e[:], x)
			binary.LittleEndian.PutUint16(e[2:], y)
			e[4] = keyStream[i+1]
			events = append(events, e[:]...)
		}
		events = append(events, 0, 0, 0, 0, 1)
		if len(events) > len(keyBuffer) {
			panic("bounded pixel stream")
		}
		copy(keyBuffer, events)
		both(cs, 0x7100, keyBuffer)
		both(cs, 0x7020, make([]byte, 8))

		word(cs, 0x1d5, c.saved)
		word(cs, 0x1d7, c.pointerDX)
		word(cs, 0x1d9, c.pointerBX)
		word(farSeg, 0x7010, c.saved)
		word(farSeg, 0x7012, c.pointerDX)
		word(farSeg, 0x7014, c.pointerBX)
		word(cs, 0x700c, 2)
		initial := append([]byte(nil), table...)
		initial[0] = c.counter
		initial[1] = c.state
		both(cs, 0xecfc, initial)
		r := oracle.CallRegs{AX: c.ax, BX: c.bx, CX: c.cx, DX: c.dx, SI: c.si, DI: c.di, BP: 0x6789, DS: bankSeg, ES: scratchSeg, SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		if c.target == 0xe3c0 {
			r.ES = mapSeg
		}
		if c.target == 0x5358 || c.target == 0x3e11 || c.target == 0x31ae {
			r.DS = cs
		}
		o.CallNear(oracle.Addr{Seg: cs, Off: c.target}, 0x100, r)
		input := o.Regs()
		pSet(&m, input)
		if err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x100}), oracle.Budget(2000000)); err != nil {
			panic(err)
		}
		C.hotspot_run(&m, C.uint16_t(c.target), &f)
		got, want := pRegs(&m), o.Regs()
		if int(f.base.base.base.calls) > 8192 {
			panic("trace capacity")
		}
		ct := make([]byte, 0, int(f.base.base.base.calls)*292)
		for i := 0; i < int(f.base.base.base.calls); i++ {
			s := f.base.base.base.trace[i]
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
			ct = append(ct, C.GoBytes(unsafe.Pointer(&f.base.base.ui[i][0]), 64)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&f.base.input[i][0]), 24)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&f.hotspot[i][0]), 16)...)
		}
		br, qr, sr, gr, rr := o.Bytes(ba, len(c.bank)), o.Bytes(qa, 1024), o.Bytes(sa, 1056), o.Bytes(ga, 59), o.Bytes(ra, 258)
		var actualUI [64]C.uint8_t
		C.md_observed_ui(&m, &f.base.base, &actualUI[0])
		cUI := C.GoBytes(unsafe.Pointer(&actualUI[0]), 64)
		originalUI := append(uiBytes(), inputBytes(o.Regs())...)
		originalUI = append(originalUI, hotspotBytes()...)
		var extra [24]C.uint8_t
		C.nu_observed(&m, &f.base, &extra[0])
		cUI = append(cUI, C.GoBytes(unsafe.Pointer(&extra[0]), 24)...)
		var hsExtra [16]C.uint8_t
		C.hs_observed(&m, &f, &hsExtra[0])
		cUI = append(cUI, C.GoBytes(unsafe.Pointer(&hsExtra[0]), 16)...)
		if value, eligible := numericGo(c); eligible {
			goCases++
			goGroups[c.group]++
			actual := int(want.SI)
			if c.target == 0x7c6e {
				actual = int(want.AX)
			}
			if value != actual {
				failure = map[string]any{"group": c.group, "case": cases, "input": input, "go_value": value, "original_value": actual}
				return false
			}
		}
		mapOriginal := o.Bytes(oracle.Addr{Seg: mapSeg, Off: 0}, 65536)
		flagOriginal := o.Bytes(oracle.Addr{Seg: flagSeg, Off: 0}, 65536)
		cases++
		groups[c.group]++
		if !bytes.Equal(mem[uint32(mapSeg)*16:uint32(mapSeg)*16+65536], mapOriginal) || !bytes.Equal(mem[uint32(flagSeg)*16:uint32(flagSeg)*16+65536], flagOriginal) || !bytes.Equal(originalUI, cUI) || got != want || !bytes.Equal(trace, ct) || !bytes.Equal(mem[ba.Linear():ba.Linear()+uint32(len(br))], br) || !bytes.Equal(mem[qa.Linear():qa.Linear()+1024], qr) || !bytes.Equal(mem[sa.Linear():sa.Linear()+1056], sr) || !bytes.Equal(mem[ga.Linear():ga.Linear()+59], gr) || !bytes.Equal(mem[ra.Linear():ra.Linear()+258], rr) || mem[oracle.Addr{Seg: cs, Off: 0x31ad}.Linear()] != o.Byte(oracle.Addr{Seg: cs, Off: 0x31ad}) {
			failure = map[string]any{"group": c.group, "case": cases - 1, "input": input, "original": want, "c": got, "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct), "original_bank": hex.EncodeToString(br), "c_bank": hex.EncodeToString(mem[ba.Linear() : ba.Linear()+uint32(len(br))]), "original_queue": hex.EncodeToString(qr), "c_queue": hex.EncodeToString(mem[qa.Linear() : qa.Linear()+1024]), "original_scratch": hex.EncodeToString(sr), "c_scratch": hex.EncodeToString(mem[sa.Linear() : sa.Linear()+1056])}
			failure["original_globals"] = hex.EncodeToString(gr)
			failure["c_globals"] = hex.EncodeToString(mem[ga.Linear() : ga.Linear()+59])
			failure["original_map"] = hex.EncodeToString(mapOriginal)
			failure["c_map"] = hex.EncodeToString(mem[uint32(mapSeg)*16 : uint32(mapSeg)*16+65536])
			failure["original_flags_map"] = hex.EncodeToString(flagOriginal)
			failure["c_flags_map"] = hex.EncodeToString(mem[uint32(flagSeg)*16 : uint32(flagSeg)*16+65536])
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
		od.Write(mapOriginal)
		od.Write(flagOriginal)
		cd.Write(pRegBytes(got))
		cd.Write(mem[ba.Linear() : ba.Linear()+uint32(len(br))])
		cd.Write(mem[qa.Linear() : qa.Linear()+uint32(len(qr))])
		cd.Write(mem[sa.Linear() : sa.Linear()+uint32(len(sr))])
		cd.Write(mem[ga.Linear() : ga.Linear()+uint32(len(gr))])
		cd.Write(mem[ra.Linear() : ra.Linear()+uint32(len(rr))])
		cd.Write(ct)
		cd.Write(cUI)
		cd.Write(mem[uint32(mapSeg)*16 : uint32(mapSeg)*16+65536])
		cd.Write(mem[uint32(flagSeg)*16 : uint32(flagSeg)*16+65536])
		if cases%128 == 0 {
			if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory")
			}
			audits++
		}
		return true
	}

	base := func(group string, target uint16) pCase {
		return pCase{group: group, target: target, ax: 0x1234, bx: 0x2345, si: 0x4240, di: 128, dx: 0x0201, bank: defaultBank, player: 0, counter: 77, state: 91, delay: 7, trust: 255, hourCursor: 64, amount: 500, farX: 639, farY: 399, originX: 24, originY: 40, music: 1, message: 0x3456, cap: 65535, saved: 2, pointerDX: 0x2222, pointerBX: 0x3333, cx: 0x3456, mapBase: 0x100}
	}

	plans := []string{"F", "X", "123F", "123X", "999999F", "M0F", "12ZF", "12BF", "12CF", "!123F", "00100F", "12M0F", "12300BF", "999999X", "M!F", "1C2B3F"}
	for y := 0; y < 400; y++ {
		for x := 0; x < 640; x++ {
			c := base("query-pixels", 0xe453)
			c.cx = uint16(x)
			c.dx = uint16(y)
			if !check(c) {
				goto finish
			}
		}
	}
	for _, baseOffset := range []uint16{0, 1, 32768, 65000, 65535} {
		for _, x := range []uint16{0, 7, 8, 639, 640, 65535} {
			for _, y := range []uint16{0, 7, 8, 255, 256, 399, 400, 65535} {
				c := base("query-word", 0xe453)
				c.cx = x
				c.dx = y
				c.mapBase = baseOffset
				if !check(c) {
					goto finish
				}
			}
		}
	}
	for _, baseOffset := range []uint16{0, 1, 32768, 65000, 65535} {
		c := base("init", 0xe3c0)
		c.di = baseOffset
		c.mapBase = baseOffset
		if !check(c) {
			goto finish
		}
	}
	for _, t := range []uint16{0xe3d7, 0xe41b} {
		for _, width := range []uint16{0, 1, 2, 79, 80, 255} {
			for _, height := range []uint16{0, 1, 2, 49, 50, 255} {
				for _, baseOffset := range []uint16{0, 65000} {
					for _, code := range []uint16{0, 1, 255} {
						c := base("rectangle", t)
						c.cx = height<<8 | width
						c.ax = 0x1200 | code
						c.dx = 639
						c.bx = 399
						c.mapBase = baseOffset
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	for y := 0; y < 400; y += 8 {
		for x := 0; x < 640; x += 8 {
			c := base("register-cells", 0xe3d7)
			c.cx = 0x101
			c.ax = 0x1253
			c.dx = uint16(x)
			c.bx = uint16(y)
			c.mapBytes = make([]byte, 65536)
			if !check(c) {
				goto finish
			}
		}
	}
	for _, t := range []uint16{0x895d, 0x89de, 0xd5d4, 0xc14, 0xc60, 0xc77} {
		for _, width := range []uint16{2, 3, 7, 40, 255} {
			for _, height := range []uint16{2, 3, 5, 25, 255} {
				for _, style := range []uint16{0, 1, 81} {
					c := base("wrapper", t)
					c.cx = height<<8 | width
					c.ax = style
					c.dx = 2
					c.bx = 3
					if t == 0x89de {
						c.ax = style << 8
					}
					if t == 0xc60 || t == 0xc77 {
						c.cx = 0
						c.dx = height<<8 | width
					}
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	for _, cap := range []uint16{0, 100, 10000, 30000, 65535} {
		for _, plan := range plans {
			c := base("numeric-pixels", 0x7c6e)
			c.ax = cap
			c.keys = numericKeys(plan)
			c.dx = 88
			c.bx = 184
			c.mapBytes = make([]byte, 65536)
			if !check(c) {
				goto finish
			}
		}
	}
	for _, t := range []uint16{0x67cd, 0x67e6, 0x6806, 0x6826} {
		for _, plan := range plans {
			c := base("finance-pixels", t)
			c.keys = numericKeys(plan)
			c.mapBytes = make([]byte, 65536)
			if !check(c) {
				goto finish
			}
		}
	}
	for scenario := 0; scenario < 4; scenario++ {
		for _, player := range []byte{0, 7, 21} {
			for _, t := range []uint16{0x3902, 0x39e8} {
				for _, plan := range []string{"123F", "999999F", "12ZF", "123X"} {
					c := base("scenario-hotspot", t)
					start := scenario*22208 + 0x80
					c.bank = append([]byte(nil), scenarios[start:start+0x5240]...)
					c.ax = 500
					c.dx = 500
					c.si = 0x4240
					c.player = player
					c.answer = 1
					c.keys = numericKeys(plan)
					c.mapBytes = make([]byte, 65536)
					if !check(c) {
						goto finish
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
	report := map[string]any{"schema": "wolong-c-hotspot-parity-v1", "input_sha256": eh, "scenario_sha256": fmt.Sprintf("%x", sha256.Sum256(scenarios)), "routine_sha256": rh, "cases": cases, "groups": groups, "full_memory_audits": audits, "passed": failure == nil, "mismatch": failure, "original_state_sha256": hex.EncodeToString(od.Sum(nil)), "c_state_sha256": hex.EncodeToString(cd.Sum(nil)), "c_machine_code_match": false, "rng_fixture": map[string]any{"table_seed_hms": []int{12, 34, 56}, "counter_state": "Explicit per case, identically written to both CS:ECFC before execution", "rerolls": 0}, "scope": "10 hotspot/map/query/window memory functions; true pixel input to numeric/finance; whole map/flags and ABI; VGA/drawing/popup/device leaves remain fixtures; IF/TF=0; original zero byte loops and word wrap"}
	report["entries_seen"] = entriesSeen
	report["runtime_load_paragraph"] = cs
	report["far_runtime_segment"] = farSeg
	report["trace_entry_size"] = 292
	report["go_cases"] = goCases
	report["go_groups"] = goGroups
	encoded, _ := json.MarshalIndent(report, "", "  ")
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*out, encoded, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("hotspot original/C %d; audits %d; pass=%v\n", cases, audits, report["passed"])
	if failure != nil {
		os.Exit(1)
	}
}
