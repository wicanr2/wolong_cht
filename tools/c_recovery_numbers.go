//go:build matching_numbers

// Original numeric raster and two caller entries; source and devices are pinned.
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
#include "/repo/tools/c_recovery/vga.c"
#include "/repo/tools/c_recovery/aligned.c"
#include "/repo/tools/c_recovery/rect.c"
#include "/repo/tools/c_recovery/display.c"
#include "/repo/tools/c_recovery/glyph.c"
#include "/repo/tools/c_recovery/numbers.c"
#include "/repo/tools/c_recovery/numbers_fixture.h"
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
	"github.com/wicanr2/dosgolem/internal/cpu"
	"github.com/wicanr2/dosgolem/internal/dos"
	"github.com/wicanr2/dosgolem/internal/machine"
	"os"
	"unsafe"
)

var activeDevice *machine.Machine
var functions = map[uint16]int{0x62f: 107, 0x69a: 68, 0x6de: 23, 0xcac: 23, 0x984: 43, 0x1e17: 47}

func entryName(t uint16) string {
	if t == 0x984 {
		return "code_10984"
	}
	return fmt.Sprintf("sub_%X", 0x10000+uint32(t))
}

type registers struct{ AX, BX, CX, DX, SI, DI, BP, SP, DS, ES, SS, CS, IP, Flags uint16 }

func originalRegs(m *machine.Machine) registers {
	c := m.CPU
	return registers{c.R[cpu.AX], c.R[cpu.BX], c.R[cpu.CX], c.R[cpu.DX], c.R[cpu.SI], c.R[cpu.DI], c.R[cpu.BP], c.R[cpu.SP], c.Seg[cpu.DS], c.Seg[cpu.ES], c.Seg[cpu.SS], c.Seg[cpu.CS], c.IP, c.Flags}
}
func cRegs(m *C.KiMachine16) registers {
	return registers{uint16(m.ax), uint16(m.bx), uint16(m.cx), uint16(m.dx), uint16(m.si), uint16(m.di), uint16(m.bp), uint16(m.sp), uint16(m.ds), uint16(m.es), uint16(m.ss), uint16(m.cs), uint16(m.ip), uint16(m.flags)}
}
func assignC(m *C.KiMachine16, r registers) {
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
func regBytes(r registers) []byte {
	b := make([]byte, 28)
	for i, v := range []uint16{r.AX, r.BX, r.CX, r.DX, r.SI, r.DI, r.BP, r.SP, r.DS, r.ES, r.SS, r.CS, r.IP, r.Flags} {
		binary.LittleEndian.PutUint16(b[i*2:], v)
	}
	return b
}
func putWord(b []byte, seg, off, v uint16) {
	at := int(seg)*16 + int(off)
	b[at] = byte(v)
	b[at+1] = byte(v >> 8)
}
func deviceBytes(m *machine.Machine) []byte {
	g, s, l := m.VGAState()
	b := append([]byte(nil), g[:]...)
	b = append(b, s[:]...)
	b = append(b, l[:]...)
	a, _ := m.VGA.In(0x3ce)
	q, _ := m.VGA.In(0x3c4)
	return append(b, a, q)
}
func portBytes(m *machine.Machine) []byte {
	b := make([]byte, 0, len(m.PortLog)*3)
	for _, w := range m.PortLog {
		b = append(b, byte(w.Port), byte(w.Port>>8), w.Val)
	}
	return b
}
func indexed(m *machine.Machine) []byte {
	b := make([]byte, 640*400)
	for y := 0; y < 400; y++ {
		for x := 0; x < 640; x++ {
			off := (y+40)*80 + x/8
			mask := byte(0x80 >> uint(x%8))
			for p := 0; p < 4; p++ {
				if m.VGA.Planes[p][off]&mask != 0 {
					b[y*640+x] |= 1 << p
				}
			}
		}
	}
	return b
}

type numberCase struct {
	group                              string
	target, ax, bx, cx, dx, si, di, bp uint16
	date                               [3]uint16
	parameter                          uint16
	mask, readMap, seed                byte
	df                                 bool
}

func main() {
	output := flag.String("out", "/output/results/O2.json", "receipt")
	only := flag.String("group", "", "group")
	smoke := flag.Bool("smoke", false, "bounded matrix")
	flag.Parse()
	raw, err := os.ReadFile("/orig/KI.EXE")
	if err != nil {
		panic(err)
	}
	inputHash := fmt.Sprintf("%x", sha256.Sum256(raw))
	if inputHash != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" {
		panic("input identity")
	}
	icons, err := os.ReadFile("/orig/ICONGRF.DAT")
	if err != nil {
		panic(err)
	}
	assetHash := fmt.Sprintf("%x", sha256.Sum256(icons))
	if len(icons) != 47776 || assetHash != "2154782c045b898aa5fafa74a4ff0c3745771ec85799b7882e1c4c009b3f1c1d" {
		panic("icon identity/shape")
	}
	original, active := machine.New(), machine.New()
	activeDevice = active
	originalDOS := dos.New(original, "/orig")
	activeDOS := dos.New(active, "/orig")
	if err = original.LoadEXE(raw); err != nil {
		panic(err)
	}
	for _, d := range []*machine.Machine{original, active} {
		d.IRQ0Every = 0
		d.KeyEvery = 0
		d.VGAFrameEvery = 0
	}
	originalDOS.Install()
	activeDOS.Install()
	cs := original.CPU.Seg[cpu.CS]
	header := int(binary.LittleEndian.Uint16(raw[8:])) * 16
	routineHashes := map[string]string{}
	for t, size := range functions {
		start := header + int(t)
		file := raw[start : start+size]
		expected := append([]byte(nil), file...)
		for i := 0; i < int(binary.LittleEndian.Uint16(raw[6:])); i++ {
			at := int(binary.LittleEndian.Uint16(raw[24:])) + i*4
			off := header + int(binary.LittleEndian.Uint16(raw[at:])) + 16*int(binary.LittleEndian.Uint16(raw[at+2:]))
			if off >= start && off < start+size {
				binary.LittleEndian.PutUint16(expected[off-start:], binary.LittleEndian.Uint16(raw[off:])+cs)
			}
		}
		at := int(cs)*16 + int(t)
		if !bytes.Equal(expected, original.Mem[at:at+size]) {
			panic(fmt.Sprintf("runtime identity %x", t))
		}
		routineHashes[entryName(t)] = fmt.Sprintf("%x", sha256.Sum256(file))
	}
	initialRAM := append([]byte(nil), original.Mem...)
	for _, base := range []int{0x22000, 0x30000, 0x50000, 0x52000, 0x60000} {
		for i := 0; i < 65536; i++ {
			initialRAM[base+i] = byte(i*37 + i/80*13 + base/4096 + 73)
		}
	}
	putWord(initialRAM, cs, 0xd48, 0x2200)
	putWord(initialRAM, cs, 0xd4a, 0x226c)
	putWord(initialRAM, cs, 0xd54, 0x2284)
	putWord(initialRAM, cs, 0xd84e, 0x5800)
	putWord(initialRAM, cs, 0xe479, 0x6000)
	putWord(initialRAM, cs, 0xe47b, 0)
	putWord(initialRAM, cs, 0xd30a, 0x5000)
	putWord(initialRAM, cs, 0xd30e, 0x5200)
	farSeg := binary.LittleEndian.Uint16(initialRAM[int(cs)*16+0x2082:])
	// Fixed FAR device contract, with explicit CS overrides for query/status controls.
	stub := []byte{0x83, 0xf8, 3, 0x75, 0x0a, 0x2e, 0x8b, 0x0e, 4, 0x70, 0x2e, 0x8b, 0x16, 6, 0x70, 0x83, 0xf8, 2, 0x75, 0x0e, 0x2e, 0xa1, 0x10, 0x70, 0x2e, 0x8b, 0x16, 0x12, 0x70, 0x2e, 0x8b, 0x1e, 0x14, 0x70, 0x39, 0xc0, 0xcb}
	copy(initialRAM[int(farSeg)*16:], stub)

	putWord(initialRAM, cs, 0xeae9, 0xe16)
	putWord(initialRAM, cs, 0xeaeb, cs)
	putWord(initialRAM, cs, 0xeaed, cs)
	putWord(initialRAM, cs, 0xeaef, 0x229a)
	putWord(initialRAM, cs, 0xeaf1, 0x1800)
	putWord(initialRAM, cs, 0xeaf3, 0x229a)
	putWord(initialRAM, cs, 0xeaf5, 0x2020)
	putWord(initialRAM, farSeg, 0x7012, 0x3456)
	putWord(initialRAM, farSeg, 0x7014, 0x789a)
	planes := make([][]byte, 8)
	for seed := 0; seed < 8; seed++ {
		b := make([]byte, 4*65536)
		for p := 0; p < 4; p++ {
			for i := 0; i < 65536; i++ {
				b[p*65536+i] = byte(i*37 + i/80*13 + p*61 + seed*29 + 73)
			}
		}
		planes[seed] = b
	}
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	cMemory := unsafe.Slice((*byte)(native), 1<<20)
	active.Mem = cMemory
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	fixture := C.KiNumbersFixture{}
	groups, entries := map[string]int{}, map[string]int{}
	sampleCounts := map[uint16]int{}
	cases, audits, pixelAudits := 0, 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	check := func(c numberCase) bool {
		if *only != "" && c.group != *only {
			return true
		}
		if *smoke && sampleCounts[c.target] >= 16 {
			return true
		}
		sampleCounts[c.target]++
		{
			copy(original.Mem, initialRAM)
			copy(cMemory, initialRAM)
			original.SetVideoMode(0x12)
			active.SetVideoMode(0x12)
			for p := 0; p < 4; p++ {
				copy(original.VGA.Planes[p][:], planes[c.seed%8][p*65536:(p+1)*65536])
				copy(active.VGA.Planes[p][:], planes[c.seed%8][p*65536:(p+1)*65536])
			}
			for _, d := range []*machine.Machine{original, active} {
				d.Out8(0x3c4, 2)
				d.Out8(0x3c5, c.mask)
				for _, q := range [][2]byte{{0, 0}, {1, 0}, {3, 0}, {4, c.readMap}, {5, 3}, {8, 0xff}} {
					d.Out8(0x3ce, q[0])
					d.Out8(0x3cf, q[1])
				}
				d.Out8(0x3ce, 0)
				d.Read8(0xa4321)
			}
			{
				copy(original.Mem[0x22000:], icons[0x9700:])
				copy(cMemory[0x22000:], icons[0x9700:])
				copy(original.Mem[0x30000:], icons[0x2800:0x6700])
				copy(cMemory[0x30000:], icons[0x2800:0x6700])
			}
		}
		for _, mem := range [][]byte{original.Mem, cMemory} {
			putWord(mem, cs, 0xcf6, c.date[0])
			mem[int(cs)*16+0xcf4] = byte(c.date[1])
			mem[int(cs)*16+0xcf0] = byte(c.date[2])
			if c.target == 0x984 {
				putWord(mem, 0x4000, c.di, c.parameter)
			}
		}
		original.PortLog = nil
		active.PortLog = nil
		flags := uint16(2 | uint16(c.seed&1))
		if c.df {
			flags |= 0x400
		}
		r := registers{AX: c.ax, BX: c.bx, CX: c.cx, DX: c.dx, SI: c.si, DI: c.di, BP: c.bp, SP: 0x8000, DS: 0x3500, ES: 0xa0c8, SS: 0x4000, CS: cs, IP: c.target, Flags: flags}
		if c.target == 0x69a {
			r.DS = 0x2284
		}
		if c.target == 0x1e17 {
			r.DS = cs
		}
		original.CPU.R = [8]uint16{r.AX, r.CX, r.DX, r.BX, r.SP, r.BP, r.SI, r.DI}
		original.CPU.Seg = [4]uint16{r.ES, r.CS, r.SS, r.DS}
		original.CPU.IP = r.IP
		original.CPU.SetFlags(r.Flags)
		putWord(original.Mem, r.SS, r.SP, 0x100)
		putWord(cMemory, r.SS, r.SP, 0x100)
		assignC(&m, originalRegs(original))
		fixture.trace.calls = 0
		fixture.unsupported = 0
		var originalTrace []byte
		steps := 0
		for ; steps < 3000000; steps++ {
			cur := originalRegs(original)
			if cur.CS == cs && cur.IP == 0x100 {
				break
			}
			_, known := functions[cur.IP]
			if cur.CS == cs && known || cur.CS == machine.StubSeg && (cur.IP == 0x410 || cur.IP == 0x414) {
				name := entryName(cur.IP)
				if cur.CS == machine.StubSeg {
					if cur.IP == 0x410 {
						name = "dosv-font-full@0080:0410"
					} else {
						name = "dosv-font-half@0080:0414"
					}
				}
				if cur.CS == farSeg {
					name = "fixture-far"
				}
				entries[name]++
				originalTrace = append(originalTrace, byte(cur.IP), byte(cur.IP>>8))
				originalTrace = append(originalTrace, regBytes(cur)...)
				for i := uint16(0); i < 32; i++ {
					originalTrace = append(originalTrace, original.Mem[uint32(cur.SS)*16+uint32(cur.SP+i)])
				}
				d := deviceBytes(original)
				originalTrace = append(originalTrace, d[:28]...)
			}
			if e := original.Step(); e != nil {
				panic(e)
			}
		}
		if steps == 3000000 {
			panic("bounded original routine")
		}
		C.numbers_run(&m, C.uint16_t(c.target), &fixture)
		if fixture.trace.calls > 8192 || fixture.unsupported != 0 {
			panic(fmt.Sprintf("unsupported=%x calls=%d input=%+v original=%+v C=%+v", fixture.unsupported, fixture.trace.calls, r, originalRegs(original), cRegs(&m)))
		}
		var cTrace []byte
		for i := 0; i < int(fixture.trace.calls); i++ {
			s := fixture.trace.trace[i]
			cTrace = append(cTrace, byte(s.target), byte(s.target>>8))
			for _, v := range s.regs {
				cTrace = append(cTrace, byte(v), byte(v>>8))
			}
			cTrace = append(cTrace, C.GoBytes(unsafe.Pointer(&s.stack[0]), 32)...)
			cTrace = append(cTrace, C.GoBytes(unsafe.Pointer(&s.device[0]), 28)...)
		}
		want, got := originalRegs(original), cRegs(&m)
		originalPorts, cPorts := portBytes(original), portBytes(active)
		originalState, cState := deviceBytes(original), deviceBytes(active)
		cases++
		groups[c.group]++
		if want != got || !bytes.Equal(original.Mem, cMemory) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(originalState, cState) || !bytes.Equal(originalPorts, cPorts) || !bytes.Equal(originalTrace, cTrace) {
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": r, "original": want, "c": got, "original_trace": hex.EncodeToString(originalTrace), "c_trace": hex.EncodeToString(cTrace), "original_ports": hex.EncodeToString(originalPorts), "c_ports": hex.EncodeToString(cPorts), "original_device": hex.EncodeToString(originalState), "c_device": hex.EncodeToString(cState), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(active.VGA.Raw())), "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(cMemory))}
			return false
		}
		audits++
		{
			if !bytes.Equal(indexed(original), indexed(active)) {
				panic("content pixels")
			}
			pixelAudits++
		}
		for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), originalState, originalPorts, originalTrace} {
			originalDigest.Write(b)
		}
		for _, b := range [][]byte{regBytes(got), cMemory, active.VGA.Raw(), cState, cPorts, cTrace} {
			cDigest.Write(b)
		}
		return true
	}
	all := []numberCase{}
	add := func(c numberCase) { all = append(all, c) }
	base := func(g string, t uint16) numberCase {
		return numberCase{group: g, target: t, mask: 15, seed: 3, bp: 0x6789, di: 100*80 + 20, bx: 0xf06}
	}
	values := []int32{-655359, -655350, -65536, -32768, -10000, -1000, -100, -10, -9, -1, 0, 1, 9, 10, 99, 100, 999, 1000, 9999, 10000, 32767, 65535, 65536, 100000, 655350, 655359}
	for _, v := range values {
		for _, width := range []uint16{0, 1, 2, 3, 6, 8} {
			for _, df := range []bool{false, true} {
				c := base("number", 0x62f)
				c.ax = uint16(v)
				c.dx = uint16(uint32(v) >> 16)
				c.bx = 0xf00 | width
				c.df = df
				add(c)
			}
		}
	}
	for color := 0; color < 256; color++ {
		for _, v := range []int32{-65536, -1, 0, 10, 32767} {
			c := base("colors", 0x62f)
			c.ax = uint16(v)
			c.dx = uint16(uint32(v) >> 16)
			c.bx = uint16(color<<8) | 6
			c.mask = byte(1 + color%15)
			c.readMap = byte(color % 4)
			c.seed = byte(color % 8)
			add(c)
		}
	}
	for digit := 0; digit < 11; digit++ {
		for color := 0; color < 256; color++ {
			c := base("digit", 0x69a)
			c.ax = 0x7654
			c.dx = uint16(digit)
			c.bx = uint16(color<<8) | 6
			c.readMap = byte(color % 4)
			c.seed = byte(color % 8)
			add(c)
		}
	}
	for color := 0; color < 256; color++ {
		c := base("blank", 0x6de)
		c.ax = uint16(color<<8) | 0x54
		c.dx = 0x7654
		c.readMap = byte(color % 4)
		c.seed = byte(color % 8)
		add(c)
	}
	for _, t := range []uint16{0x69a, 0x6de} {
		for _, off := range []uint16{0, 65535, 65520, 383*80 + 79} {
			for _, count := range []uint16{0, 0x100} {
				for _, df := range []bool{false, true} {
					digits := 11
					if t == 0x6de {
						digits = 1
					}
					for digit := 0; digit < digits; digit++ {
						c := base("primitive-boundary", t)
						c.di = off
						c.cx = count
						c.df = df
						c.ax = 0x5054
						c.dx = uint16(digit)
						add(c)
					}
				}
			}
		}
	}
	for _, date := range [][3]uint16{{196, 1, 1}, {0, 0, 0}, {65535, 255, 255}, {999, 12, 31}} {
		c := base("date-caller", 0x1e17)
		c.date = date
		add(c)
	}
	for _, v := range []uint16{0x8000, 0xfc18, 0xffff, 0, 1, 32767, 32768} {
		for _, xy := range [][2]uint16{{0, 0}, {80, 100}, {639, 383}} {
			c := base("parameter-caller", 0x984)
			c.parameter = v
			c.di = 0x9000
			c.dx = xy[0]
			c.bx = xy[1]
			add(c)
		}
	}
	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{"schema": "wolong-c-numbers-parity-v1", "input_sha256": inputHash, "icon_sha256": assetHash, "routine_sha256": routineHashes, "cases": cases, "groups": groups, "full_ram_plane_audits": audits, "indexed_content_audits": pixelAudits, "content_top": 40, "content_height": 400, "entries_seen": entries, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "trace_entry_size": 90, "platform": "Independent pinned dosgolem VGA devices; native C numbers; IF/TF=0; CPU-model undefined flags; legal DIV quotient; no wall-clock claim", "asset_banks": "ICONGRF segment3 at 2200; digits at 2284", "original_font_calls": originalDOS.Font.Calls, "c_font_calls": activeDOS.Font.Calls, "original_missing_fonts": originalDOS.Font.Missing, "c_missing_fonts": activeDOS.Font.Missing, "c_machine_code_match": false}
	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if e := os.WriteFile(*output, b, 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Numbers original/C %d; complete RAM/plane %d; content %d; pass=%v\n", cases, audits, pixelAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
