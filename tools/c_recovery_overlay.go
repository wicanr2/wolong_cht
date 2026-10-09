//go:build matching_overlay

// Original army/object producers, matrix clipping and minimap pixels.
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
#include "/repo/tools/c_recovery/mapcells.c"
#include "/repo/tools/c_recovery/overlay.c"
#include "/repo/tools/c_recovery/overlay_fixture.h"
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
	"github.com/wicanr2/dosgolem/internal/machine"
	"os"
	"unsafe"
)

var activeDevice *machine.Machine
var functions = map[uint16]int{0x1cc9: 7, 0x2533: 112, 0x2af4: 54, 0x2b2a: 18, 0x2b3c: 108, 0x5d19: 162, 0xd51f: 181, 0xcac: 23, 0xd46a: 25, 0xd483: 32, 0xd4c7: 88, 0xd615: 85, 0xd66a: 257, 0xd76b: 23, 0xd782: 20, 0xd796: 81, 0xd7e7: 29, 0xd804: 70}

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
			off := y*80 + x/8
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

func name(t uint16) string {
	if t == 0xd51f {
		return "code_1D51F"
	}
	return fmt.Sprintf("sub_%X", 0x10000+uint32(t))
}

type overlayCase struct {
	group                                                          string
	target, x, y                                                   uint16
	id, phase, formation, recordFlags, capacity, layers, slotFlags byte
	pattern, scenario                                              int
	minimap, keep                                                  bool
}

func main() {
	output := flag.String("out", "/output/results/O2.json", "receipt")
	only := flag.String("group", "", "group")
	smoke := flag.Bool("smoke", false, "small matrix")
	flag.Parse()
	read := func(p string) []byte {
		b, e := os.ReadFile(p)
		if e != nil {
			panic(e)
		}
		return b
	}
	raw, mdl, mch, world, scenarios := read("/orig/KI.EXE"), read("/orig/MMAP.MDL"), read("/orig/MMAP.MCH"), read("/output/fixtures/MMAP.raw"), read("/orig/SINARIO.DAT")
	digest := func(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
	if digest(raw) != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" || len(scenarios) != 88832 || len(world) != 98304 || len(mch) != 43058 || len(mdl) != 32768 {
		panic("input identity/inventory")
	}
	patterns := []int{}
	for i := 0; i < 64; i++ {
		at := 0xa000 + i*4
		w, h := int(mch[at]), int(mch[at+1])
		off := int(binary.LittleEndian.Uint16(mch[at+2:]))
		if w > 0 && h > 0 {
			if 0xa100+off+w*h > len(mch) {
				panic("pattern bounds")
			}
			patterns = append(patterns, i)
		}
	}
	original := machine.New()
	if e := original.LoadEXE(raw); e != nil {
		panic(e)
	}
	activeDevice = machine.New()
	for _, d := range []*machine.Machine{original, activeDevice} {
		d.IRQ0Every = 0
		d.KeyEvery = 0
		d.VGAFrameEvery = 0
	}
	cs := original.CPU.Seg[cpu.CS]
	hashes := map[string]string{}
	for at, size := range functions {
		data := raw[int(at)+512 : int(at)+512+size]
		if !bytes.Equal(data, original.Mem[int(cs)*16+int(at):int(cs)*16+int(at)+size]) {
			panic("routine identity")
		}
		hashes[name(at)] = digest(data)
	}
	initial := append([]byte(nil), original.Mem...)
	copy(initial[0x30000:], mdl)
	copy(initial[0x38000:], mch)
	copy(initial[0x50000:], world)
	for _, v := range [][2]uint16{{0xd84a, 0x3000}, {0xd84c, 0x3800}, {0xd84e, 0x2500}, {0xd850, 0x5000}, {0xd854, 100}, {0xd856, 90}, {0xd52, 0x7000}, {0x987a, 0x4200}} {
		putWord(initial, cs, v[0], v[1])
	}
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	memory := unsafe.Slice((*byte)(native), 1<<20)
	activeDevice.Mem = memory
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	fixture := C.KiMapcellsFixture{}
	cases, audits := 0, 0
	sourceAudits := 0
	groups, entries, maxSteps := map[string]int{}, map[string]int{}, map[string]int{}
	odigest, cdigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	check := func(c overlayCase) bool {
		if *only != "" && c.group != *only {
			return true
		}
		if *smoke && groups[c.group] >= 8 {
			return true
		}
		if !c.keep {
			copy(original.Mem, initial)
			copy(memory, initial)
			for _, d := range []*machine.Machine{original, activeDevice} {
				d.SetVideoMode(0x12)
				for p := 0; p < 4; p++ {
					for i := 0; i < 65536; i++ {
						d.VGA.Planes[p][i] = byte(i*19 + i/80*7 + p*53 + 47)
					}
				}
				d.Out8(0x3c4, 2)
				d.Out8(0x3c5, 15)
				for _, q := range [][2]byte{{0, 0}, {1, 0}, {3, 0}, {4, 0}, {5, 0}, {8, 255}} {
					d.Out8(0x3ce, q[0])
					d.Out8(0x3cf, q[1])
				}
			}
			for _, b := range [][]byte{original.Mem, memory} {
				for i := 0; i < 920; i++ {
					at := 0x25000 + i*8
					b[at] = c.slotFlags
					b[at+1] = c.layers
					b[at+2] = byte(i*13 + 7)
					for j := 3; j < 8; j++ {
						b[at+j] = byte(i*29 + j*17)
					}
				}
				start := c.scenario * 22208
				copy(b[0x70000:], scenarios[start+0x80:start+0x52c0])
				copy(b[int(cs)*16+0xcf0:], scenarios[start:start+59])
				if c.group != "corpus" && c.group != "pipeline" {
					clear(b[0x70000+0x2040 : 0x70000+0x2240])
					clear(b[0x70000+0x2240 : 0x70000+0x4200])
					at := 0x70000 + 0x2240
					b[at] = c.recordFlags
					b[at+3] = c.phase
					b[at+7] = 211
					b[at+8] = c.phase
					b[at+9] = c.id * 5
					b[at+0x21] = c.formation
					putWord(b, 0x7000, 0x2250, c.x)
					putWord(b, 0x7000, 0x2252, c.y)
					at = 0x70000 + 0x2040
					b[at] = c.recordFlags
					b[at+0xe] = c.id
					b[at+0xf] = c.phase
					putWord(b, 0x7000, 0x2042, c.x)
					putWord(b, 0x7000, 0x2044, c.y)
				}
				if c.group == "corpus" {
					putWord(b, cs, 0xd854, c.x)
					putWord(b, cs, 0xd856, c.y)
				}
				b[int(cs)*16+0x98a6] = 0
				if c.minimap {
					b[int(cs)*16+0x98a6] = 4
				}
			}
		}
		original.PortLog = nil
		activeDevice.PortLog = nil
		rr := registers{AX: 0x0701, BX: c.y, CX: 1, DX: c.x, SI: 0x2240, DI: 0, BP: 0x6789, SP: 0x8000, DS: 0x7000, ES: 0x4200, SS: 0x9000, CS: cs, IP: c.target, Flags: 2}
		if c.target == 0x1cc9 || c.target == 0x2533 || c.target == 0x2af4 {
			rr.DS = cs
		}
		if c.target == 0x5d19 {
			rr.AX = uint16(c.id) << 8
		}
		if c.target == 0xd51f {
			at := 0xa000 + c.pattern*4
			rr.CX = uint16(mch[at])
			rr.AX = uint16(c.capacity)<<8 | uint16(mch[at+1])
			rr.DI = 0x100 + binary.LittleEndian.Uint16(mch[at+2:])
		}
		if c.target == 0xd46a {
			rr.AX = 0x3000
			rr.BX = 0x2500
			rr.CX = 0x5000
		}
		if c.target == 0xd615 {
			rr.DS = cs
		}
		var expectedGrid []byte
		if c.target == 0xd51f {
			expectedGrid = append([]byte(nil), original.Mem[0x25000:0x25000+7360]...)
			at := 0xa000 + c.pattern*4
			w, h := int(mch[at]), int(mch[at+1])
			source := 0xa100 + int(binary.LittleEndian.Uint16(mch[at+2:]))
			originX, originY := int(int16(c.x-100)), int(int16(c.y-90))
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					px, py := originX+x, originY+y
					if px < 0 || px >= 40 || py < 0 || py >= 23 {
						continue
					}
					id := mch[source+y*w+x]
					slot := (py*40 + px) * 8
					n := int(expectedGrid[slot+1])
					if id == 255 || expectedGrid[slot]&16 != 0 || n >= int(c.capacity) {
						continue
					}
					old := expectedGrid[slot+3+n]
					expectedGrid[slot+3+n] = id
					expectedGrid[slot+1] = byte(n + 1)
					if old != id {
						expectedGrid[slot] |= 32
					}
				}
			}
		}
		o := original.CPU
		o.R = [8]uint16{rr.AX, rr.CX, rr.DX, rr.BX, rr.SP, rr.BP, rr.SI, rr.DI}
		o.Seg = [4]uint16{rr.ES, rr.CS, rr.SS, rr.DS}
		o.IP = rr.IP
		o.SetFlags(rr.Flags)
		putWord(original.Mem, rr.SS, rr.SP, 0x100)
		putWord(memory, rr.SS, rr.SP, 0x100)
		assignC(&m, originalRegs(original))
		fixture.calls = 0
		fixture.blocks[0].calls = 0
		fixture.blocks[1].calls = 0
		var trace []byte
		for step := 0; step < 4000000; step++ {
			now := originalRegs(original)
			if now.CS == cs && now.IP == 0x100 {
				if step > maxSteps[c.group] {
					maxSteps[c.group] = step
				}
				break
			}
			if _, ok := functions[now.IP]; ok && now.CS == cs {
				entries[name(now.IP)]++
				trace = append(trace, byte(now.IP), byte(now.IP>>8))
				trace = append(trace, regBytes(now)...)
				for j := uint16(0); j < 32; j++ {
					trace = append(trace, original.Mem[uint32(now.SS)*16+uint32(uint16(now.SP+j))])
				}
				trace = append(trace, deviceBytes(original)[:28]...)
			}
			if e := original.Step(); e != nil {
				panic(e)
			}
			if step == 3999999 {
				panic(fmt.Sprintf("bounded original: %+v; regs %+v", c, now))
			}
		}
		C.overlay_run(&m, C.uint16_t(c.target), &fixture)
		if fixture.calls > 16384 {
			panic("trace capacity")
		}
		var ct []byte
		for i := 0; i < int(fixture.calls); i++ {
			s := fixture.blocks[i/8192].trace[i%8192]
			ct = append(ct, byte(s.target), byte(s.target>>8))
			for _, v := range s.regs {
				ct = append(ct, byte(v), byte(v>>8))
			}
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.stack[0]), 32)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.device[0]), 28)...)
		}
		want, got := originalRegs(original), cRegs(&m)
		op, cp := portBytes(original), portBytes(activeDevice)
		od, cd := deviceBytes(original), deviceBytes(activeDevice)
		cases++
		groups[c.group]++
		if want != got || !bytes.Equal(original.Mem, memory) || !bytes.Equal(original.VGA.Raw(), activeDevice.VGA.Raw()) || !bytes.Equal(op, cp) || !bytes.Equal(od, cd) || !bytes.Equal(trace, ct) {
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": rr, "original": want, "c": got, "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct), "original_ports": hex.EncodeToString(op), "c_ports": hex.EncodeToString(cp), "original_device": hex.EncodeToString(od), "c_device": hex.EncodeToString(cd), "original_ram": digest(original.Mem), "c_ram": digest(memory), "original_planes": digest(original.VGA.Raw()), "c_planes": digest(activeDevice.VGA.Raw())}
			return false
		}
		audits++
		if expectedGrid != nil {
			if !bytes.Equal(original.Mem[0x25000:0x25000+7360], expectedGrid) {
				panic("independent MCH matrix/slot mismatch")
			}
			sourceAudits++
		}
		for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), op, od, trace} {
			odigest.Write(b)
		}
		for _, b := range [][]byte{regBytes(got), memory, activeDevice.VGA.Raw(), cp, cd, ct} {
			cdigest.Write(b)
		}
		return true
	}
	base := func(g string, t uint16) overlayCase {
		return overlayCase{group: g, target: t, x: 120, y: 100, recordFlags: 0xc0, capacity: 5, slotFlags: 0x40}
	}
	all := []overlayCase{}
	add := func(c overlayCase) { all = append(all, c) }
	coords := [][2]uint16{{99, 90}, {100, 89}, {100, 90}, {139, 112}, {140, 90}, {100, 113}, {120, 100}, {120, 356}, {90, 80}}
	for faction := byte(0); faction < 22; faction++ {
		for dir := byte(0); dir < 5; dir++ {
			for _, xy := range coords {
				c := base("point", 0x2b2a)
				c.id = faction
				c.phase = dir
				c.x = xy[0]
				c.y = xy[1]
				add(c)
			}
		}
	}
	for formation := byte(0); formation < 5; formation++ {
		for phase := byte(0); phase < 4; phase++ {
			for _, xy := range coords[:7] {
				for _, mini := range []bool{false, true} {
					c := base("army", 0x2b3c)
					c.formation = formation
					c.phase = phase
					c.x = xy[0]
					c.y = xy[1]
					c.minimap = mini
					add(c)
				}
			}
		}
	}
	for typ := byte(0); typ < 4; typ++ {
		for phase := byte(0); phase < 8; phase++ {
			for _, f := range []byte{0x7f, 0x80, 0x81, 0xc1} {
				for _, xy := range coords[:7] {
					c := base("objects", 0x2533)
					c.id = typ
					c.phase = phase
					c.recordFlags = f
					c.x = xy[0]
					c.y = xy[1]
					add(c)
				}
			}
		}
	}
	for _, pattern := range patterns {
		for _, xy := range coords[:7] {
			for _, cap := range []byte{0, 3, 4, 5} {
				for _, slots := range []byte{0, 3, 4, 5} {
					for _, protect := range []byte{0x40, 0x50} {
						c := base("matrix", 0xd51f)
						c.pattern = pattern
						c.x = xy[0]
						c.y = xy[1]
						c.capacity = cap
						c.layers = slots
						c.slotFlags = protect
						add(c)
					}
				}
			}
		}
	}
	for bit := uint16(0); bit < 8; bit++ {
		for _, color := range []byte{0x3f, 0xfa} {
			for _, y := range []uint16{38, 89, 166} {
				c := base("minimap", 0x5d19)
				c.x = 438 + bit
				c.y = y
				c.id = color
				add(c)
			}
		}
	}
	for _, f := range []byte{0, 0x7f, 0x80, 0xbf, 0xc0, 0xdf, 0xe0, 0xff} {
		for _, xy := range coords[:7] {
			c := base("scan", 0x2af4)
			c.recordFlags = f
			c.x = xy[0]
			c.y = xy[1]
			add(c)
		}
	}
	for scenario := 0; scenario < 4; scenario++ {
		for _, xy := range [][2]uint16{{0, 0}, {100, 90}, {170, 98}, {344, 233}} {
			c := base("corpus", 0x1cc9)
			c.scenario = scenario
			c.x = xy[0]
			c.y = xy[1]
			add(c)
		}
	}
	for scenario := 0; scenario < 4; scenario++ {
		c := base("pipeline", 0xd46a)
		c.scenario = scenario
		add(c)
		c = base("pipeline", 0xd615)
		c.keep = true
		c.x = 100
		c.y = 90
		add(c)
		c = base("pipeline", 0x1cc9)
		c.keep = true
		add(c)
		c = base("pipeline", 0xd66a)
		c.keep = true
		add(c)
	}
	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{"schema": "wolong-c-overlay-parity-v1", "input_sha256": digest(raw), "routine_sha256": hashes, "groups": groups, "cases": cases, "full_ram_plane_audits": audits, "source_matrix_audits": sourceAudits, "entries_seen": entries, "valid_metadata_patterns": patterns, "max_original_steps": maxSteps, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(odigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cdigest.Sum(nil)), "asset_sha256": map[string]string{"MMAP.MDL": digest(mdl), "MMAP.MCH": digest(mch), "MMAP.raw": digest(world), "SINARIO.DAT": digest(scenarios)}, "c_machine_code_match": false}
	b, e := json.MarshalIndent(report, "", "  ")
	if e != nil {
		panic(e)
	}
	if e = os.WriteFile(*output, append(b, '\n'), 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Overlay original/C %d; whole RAM/plane %d; pass=%v\n", cases, audits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
