//go:build matching_mapcells

// Original world display cells; independent VGA and real MMAP assets.
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
#include "/repo/tools/c_recovery/mapcells_fixture.h"
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
var functions = map[uint16]int{0xd46a: 25, 0xd483: 32, 0xd4c7: 88, 0xd615: 85, 0xd66a: 257, 0xd76b: 23, 0xd782: 20, 0xd796: 81, 0xd7e7: 29, 0xd804: 70}

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

type cellsCase struct {
	group                          string
	target, ax, bx, cx, dx, si, di uint16
	flags, layers, seed            byte
	df, keep, full                 bool
	audit                          string
}

func main() {
	output := flag.String("out", "/output/results/O2.json", "receipt")
	only := flag.String("group", "", "case group")
	smoke := flag.Bool("smoke", false, "small matrix")
	flag.Parse()
	read := func(path string) []byte {
		b, e := os.ReadFile(path)
		if e != nil {
			panic(e)
		}
		return b
	}
	raw := read("/orig/KI.EXE")
	inputHash := fmt.Sprintf("%x", sha256.Sum256(raw))
	if inputHash != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" {
		panic("input identity")
	}
	mdl, mch, world := read("/orig/MMAP.MDL"), read("/orig/MMAP.MCH"), read("/output/fixtures/MMAP.raw")
	if len(mdl) != 32768 || len(mch) != 43058 || len(world) != 98304 {
		panic("asset inventory")
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
	for t, size := range functions {
		b := raw[int(t)+512 : int(t)+512+size]
		if !bytes.Equal(b, original.Mem[int(cs)*16+int(t):int(cs)*16+int(t)+size]) {
			panic("routine identity")
		}
		hashes[fmt.Sprintf("sub_%X", 0x10000+uint32(t))] = fmt.Sprintf("%x", sha256.Sum256(b))
	}
	initial := append([]byte(nil), original.Mem...)
	for _, base := range []int{0x25000, 0x30000, 0x40000, 0x70000} {
		for i := 0; i < 65536; i++ {
			initial[base+i] = byte(i*37 + i/257 + base/4096)
		}
	}
	copy(initial[0x30000:], mdl)
	copy(initial[0x38000:], mch)
	copy(initial[0x50000:], world)
	putWord(initial, cs, 0xd84a, 0x3000)
	putWord(initial, cs, 0xd84c, 0x3800)
	putWord(initial, cs, 0xd84e, 0x2500)
	putWord(initial, cs, 0xd850, 0x5000)
	putWord(initial, cs, 0xd854, 100)
	putWord(initial, cs, 0xd856, 90)
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	memory := unsafe.Slice((*byte)(native), 1<<20)
	activeDevice.Mem = memory
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	fixture := C.KiMapcellsFixture{}
	groups, entries := map[string]int{}, map[string]int{}
	cases, audits, sourceAudits := 0, 0, 0
	tileAudits := 0
	maxSteps := map[string]int{}
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	check := func(c cellsCase) bool {
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
						d.VGA.Planes[p][i] = byte(i*19 + i/80*7 + p*53 + int(c.seed)*31)
					}
				}
				d.Out8(0x3c4, 2)
				d.Out8(0x3c5, 15)
				for _, q := range [][2]byte{{0, 0}, {1, 0}, {3, 0}, {4, byte(c.seed % 4)}, {5, 0}, {8, 255}} {
					d.Out8(0x3ce, q[0])
					d.Out8(0x3cf, q[1])
				}
			}
			if c.group == "flags" || c.group == "render" {
				for _, b := range [][]byte{original.Mem, memory} {
					for i := 0; i < 920; i++ {
						at := 0x25000 + i*8
						b[at] = 0x10
						b[at+1] = 0
						b[at+2] = byte(i*13 + 7)
						for j := 3; j < 8; j++ {
							b[at+j] = byte(i*31 + j*17)
						}
						if c.full || i == 417 {
							b[at] = c.flags
							b[at+1] = c.layers
						}
					}
				}
			}
			if c.group == "push" {
				for _, b := range [][]byte{original.Mem, memory} {
					for i := 0; i < 920; i++ {
						b[0x25000+i*8] = c.flags
						b[0x25001+i*8] = c.layers
					}
				}
			}
		}
		original.PortLog = nil
		activeDevice.PortLog = nil
		flags := uint16(2)
		if c.df {
			flags |= 0x400
		}
		r := registers{AX: c.ax, BX: c.bx, CX: c.cx, DX: c.dx, SI: c.si, DI: c.di, BP: 0x6789, SP: 0x8000, DS: 0x3000, ES: 0xa0c8, SS: 0x9000, CS: cs, IP: c.target, Flags: flags}
		if c.target == 0xd804 {
			r.DS = 0x3800
		}
		if c.target == 0xd615 {
			r.DS = cs
		}
		if c.target == 0xd76b {
			r.DS = cs
		}
		if c.target == 0xd46a {
			r.AX = 0x3000
			r.BX = 0x2500
			r.CX = 0x5000
		}
		var expected []byte
		if c.audit == "background" && !c.df {
			expected = append([]byte(nil), mdl[int(c.ax>>8)*128:int(c.ax>>8)*128+128]...)
		}
		if c.target == 0xd804 && !c.keep {
			at := int(cs)*16 + 0xd858
			base := ((int(c.ax>>8)*53 + 7) & 255) * 128
			copy(original.Mem[at:at+128], mdl[base:base+128])
			copy(memory[at:at+128], mdl[base:base+128])
		}
		if c.audit == "overlay" && !c.df {
			at := int(cs)*16 + 0xd858
			expected = append([]byte(nil), original.Mem[at:at+128]...)
			tile := mch[int(c.ax>>8)*160 : int(c.ax>>8)*160+160]
			for p := 0; p < 4; p++ {
				for i := 0; i < 32; i++ {
					expected[p*32+i] = (expected[p*32+i] &^ tile[i]) | tile[32+p*32+i]
				}
			}
		}
		o := original.CPU
		o.R = [8]uint16{r.AX, r.CX, r.DX, r.BX, r.SP, r.BP, r.SI, r.DI}
		o.Seg = [4]uint16{r.ES, r.CS, r.SS, r.DS}
		o.IP = r.IP
		o.SetFlags(r.Flags)
		putWord(original.Mem, r.SS, r.SP, 0x100)
		putWord(memory, r.SS, r.SP, 0x100)
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
				entries[fmt.Sprintf("sub_%X", 0x10000+uint32(now.IP))]++
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
				panic("bounded original routine")
			}
		}
		C.mapcells_run(&m, C.uint16_t(c.target), &fixture)
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
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": r, "original": want, "c": got, "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct), "original_ports": hex.EncodeToString(op), "c_ports": hex.EncodeToString(cp), "original_device": hex.EncodeToString(od), "c_device": hex.EncodeToString(cd), "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(memory)), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(activeDevice.VGA.Raw()))}
			return false
		}
		audits++
		if c.target == 0xd782 && !c.df {
			tile := mdl[int(c.ax>>8)*128 : int(c.ax>>8)*128+128]
			for plane := 0; plane < 4; plane++ {
				for row := 0; row < 16; row++ {
					at := 0xc80 + int(c.di) + row*80
					if !bytes.Equal(original.VGA.Planes[plane][at:at+2], tile[plane*32+row*2:plane*32+row*2+2]) {
						panic("independent tile blit")
					}
				}
			}
			tileAudits++
		}

		if expected != nil {
			at := int(cs)*16 + 0xd858
			if !bytes.Equal(original.Mem[at:at+128], expected) {
				panic("independent source composite")
			}
			sourceAudits++
		}
		for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), op, od, trace} {
			originalDigest.Write(b)
		}
		for _, b := range [][]byte{regBytes(got), memory, activeDevice.VGA.Raw(), cp, cd, ct} {
			cDigest.Write(b)
		}
		return true
	}
	base := func(g string, t uint16) cellsCase {
		return cellsCase{group: g, target: t, ax: 0x0701, di: 0xa00, bx: 170, dx: 98, seed: 3}
	}
	all := []cellsCase{}
	add := func(c cellsCase) { all = append(all, c) }
	for _, t := range []uint16{0xd46a, 0xd483} {
		for _, df := range []bool{false, true} {
			for seed := byte(0); seed < 3; seed++ {
				c := base("init", t)
				c.df = df
				c.seed = seed
				add(c)
			}
		}
	}
	for _, xy := range [][2]uint16{{0, 0}, {170, 98}, {344, 233}, {344, 0}, {0, 233}} {
		for _, df := range []bool{false, true} {
			c := base("copy", 0xd615)
			c.dx = xy[0]
			c.bx = xy[1]
			c.df = df
			add(c)
		}
	}
	for id := 0; id < 256; id++ {
		for _, t := range []uint16{0xd782, 0xd7e7} {
			for _, df := range []bool{false, true} {
				if t == 0xd7e7 && df {
					continue
				} // Reverse copy would overwrite the original executing code.
				c := base("background", t)
				c.ax = uint16(id)<<8 | 0x71
				c.df = df
				if t == 0xd7e7 {
					c.audit = "background"
				}
				add(c)
			}
		}
		for _, df := range []bool{false, true} {
			c := base("overlay", 0xd804)
			c.ax = uint16(id)<<8 | 0x63
			c.df = df
			c.audit = "overlay"
			add(c)
		}
	}
	for _, t := range []uint16{0xd76b, 0xd796} {
		for _, dest := range []uint16{0, 79, 2560, 65535} {
			for _, df := range []bool{false, true} {
				c := base("flush", t)
				c.di = dest
				c.si = 1
				c.df = df
				add(c)
			}
		}
	}
	for bits := 0; bits < 256; bits++ {
		for _, n := range []byte{0, 1, 4, 5} {
			c := base("flags", 0xd66a)
			c.flags = byte(bits)
			c.layers = n
			add(c)
		}
	}
	for _, xy := range [][2]uint16{{99, 90}, {100, 89}, {100, 90}, {139, 112}, {140, 90}, {100, 113}, {120, 100}, {65535, 90}} {
		for _, n := range []byte{0, 3, 4, 5} {
			for _, f := range []byte{0, 16} {
				for _, id := range []byte{0, 1, 254, 255} {
					c := base("push", 0xd4c7)
					c.dx = xy[0]
					c.bx = xy[1]
					c.layers = n
					c.flags = f
					c.ax = uint16(id)
					add(c)
				}
			}
		}
	}
	for _, n := range []byte{0, 1, 4, 5} {
		for _, df := range []bool{false, true} {
			c := base("render", 0xd66a)
			c.flags = 0xe0
			c.layers = n
			c.full = true
			c.df = df
			add(c)
		}
	}
	for _, xy := range [][2]uint16{{0, 0}, {170, 98}, {344, 233}, {64, 100}} {
		c := base("pipeline", 0xd46a)
		add(c)
		c = base("pipeline", 0xd615)
		c.keep = true
		c.dx = xy[0]
		c.bx = xy[1]
		add(c)
		c = base("pipeline", 0xd4c7)
		c.keep = true
		c.ax = 255
		c.dx = xy[0] + 20
		c.bx = xy[1] + 12
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
	report := map[string]any{"schema": "wolong-c-mapcells-parity-v1", "input_sha256": inputHash, "routine_sha256": hashes, "groups": groups, "cases": cases, "full_ram_plane_audits": audits, "source_composite_audits": sourceAudits, "tile_blit_audits": tileAudits, "max_original_steps": maxSteps, "original_step_bound": 4000000, "entries_seen": entries, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "asset_sha256": map[string]string{"MMAP.MDL": fmt.Sprintf("%x", sha256.Sum256(mdl)), "MMAP.MCH": fmt.Sprintf("%x", sha256.Sum256(mch)), "MMAP.raw": fmt.Sprintf("%x", sha256.Sum256(world))}, "c_machine_code_match": false}
	b, e := json.MarshalIndent(report, "", "  ")
	if e != nil {
		panic(e)
	}
	if e = os.WriteFile(*output, append(b, '\n'), 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Map cells original/C %d; whole RAM/plane %d; source composites %d; pass=%v\n", cases, audits, sourceAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
