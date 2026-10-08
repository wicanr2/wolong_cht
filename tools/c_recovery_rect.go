//go:build matching_rect

// Sixteen original rectangle/bar/caller entries; source and platform are pinned.
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
#include "/repo/tools/c_recovery/rect_fixture.h"
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
var functions = map[uint16]int{0x1b4: 33, 0x1db: 51, 0xaaa: 47, 0xad9: 109, 0xbcd: 71, 0xc14: 76, 0xc60: 23, 0xc77: 53, 0xcac: 23, 0xcc3: 27, 0x895d: 71, 0x89de: 18, 0xc61f: 52, 0xc673: 59, 0xc6ae: 17, 0xc6bf: 55, 0xc6f6: 86, 0xc74c: 41, 0xc775: 25, 0xc78e: 27, 0xc7f4: 74, 0xd5d4: 65, 0xe3d7: 68, 0xe453: 38, 0xf020: 288, 0xf140: 59, 0xf17b: 39, 0xf1a3: 203, 0xf888: 176, 0xf938: 97, 0xf999: 23, 0xf9b0: 107, 0xfa1b: 28, 0xfa37: 107, 0xfaa2: 32}

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

type rectCase struct {
	group                              string
	target, ax, bx, cx, dx, si, di, bp uint16
	mouseX, mouseY, service            uint16
	values                             [10]uint16
	mask, readMap, seed                byte
	df, realAsset, keep, compareBar    bool
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
	if err = original.LoadEXE(raw); err != nil {
		panic(err)
	}
	for _, d := range []*machine.Machine{original, active} {
		d.IRQ0Every = 0
		d.KeyEvery = 0
		d.VGAFrameEvery = 0
	}
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
		routineHashes[fmt.Sprintf("sub_%X", 0x10000+uint32(t))] = fmt.Sprintf("%x", sha256.Sum256(file))
	}
	initialRAM := append([]byte(nil), original.Mem...)
	for _, base := range []int{0x22000, 0x30000, 0x50000, 0x52000, 0x60000} {
		for i := 0; i < 65536; i++ {
			initialRAM[base+i] = byte(i*37 + i/80*13 + base/4096 + 73)
		}
	}
	putWord(initialRAM, cs, 0xd48, 0x2200)
	putWord(initialRAM, cs, 0xd4a, 0x226c)
	putWord(initialRAM, cs, 0xd84e, 0x5800)
	putWord(initialRAM, cs, 0xe479, 0x6000)
	putWord(initialRAM, cs, 0xe47b, 0)
	putWord(initialRAM, cs, 0xd30a, 0x5000)
	putWord(initialRAM, cs, 0xd30e, 0x5200)
	farSeg := binary.LittleEndian.Uint16(initialRAM[int(cs)*16+0x2082:])
	// Fixed FAR device contract, with explicit CS overrides for query/status controls.
	stub := []byte{0x83, 0xf8, 3, 0x75, 0x0a, 0x2e, 0x8b, 0x0e, 4, 0x70, 0x2e, 0x8b, 0x16, 6, 0x70, 0x83, 0xf8, 2, 0x75, 0x0e, 0x2e, 0xa1, 0x10, 0x70, 0x2e, 0x8b, 0x16, 0x12, 0x70, 0x2e, 0x8b, 0x1e, 0x14, 0x70, 0x39, 0xc0, 0xcb}
	copy(initialRAM[int(farSeg)*16:], stub)
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
	fixture := C.KiRectFixture{}
	groups, entries := map[string]int{}, map[string]int{}
	sampleCounts := map[uint16]int{}
	cases, audits, pixelAudits, goCases := 0, 0, 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	check := func(c rectCase) bool {
		if *only != "" && c.group != *only {
			return true
		}
		if *smoke && sampleCounts[c.target] >= 16 {
			return true
		}
		sampleCounts[c.target]++
		if !c.keep {
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
				for _, q := range [][2]byte{{0, 0}, {1, 0}, {3, 0}, {4, c.readMap}, {5, 0}, {8, 0xff}} {
					d.Out8(0x3ce, q[0])
					d.Out8(0x3cf, q[1])
				}
				d.Read8(0xa4321)
			}
			if c.realAsset {
				copy(original.Mem[0x22000:], icons[0x9700:])
				copy(cMemory[0x22000:], icons[0x9700:])
				copy(original.Mem[0x30000:], icons[0x2800:0x6700])
				copy(cMemory[0x30000:], icons[0x2800:0x6700])
			}
		}
		for _, mem := range [][]byte{original.Mem, cMemory} {
			putWord(mem, farSeg, 0x7004, c.mouseX)
			putWord(mem, farSeg, 0x7006, c.mouseY)
			putWord(mem, farSeg, 0x7010, c.service)
			for i := 0; i < 6; i++ {
				mem[0x50009+i*4] = byte(c.values[i])
			}
			putWord(mem, 0x5000, 4, c.values[6])
			putWord(mem, 0x5000, 36, c.values[7])
			putWord(mem, 0x5200, 3, c.values[8])
			putWord(mem, 0x5200, 0x603, c.values[9])
		}
		original.PortLog = nil
		active.PortLog = nil
		flags := uint16(2 | uint16(c.seed&1))
		if c.df {
			flags |= 0x400
		}
		r := registers{AX: c.ax, BX: c.bx, CX: c.cx, DX: c.dx, SI: c.si, DI: c.di, BP: c.bp, SP: 0x8000, DS: 0x3500, ES: 0xa0c8, SS: 0x4000, CS: cs, IP: c.target, Flags: flags}
		if c.target == 0xf140 || c.target == 0xf17b {
			r.DS = 0xa0c8
		}
		if c.target == 0xc74c {
			r.DS = 0x5000
			r.SI = 9
		}
		if c.target == 0xc7f4 {
			r.DS = 0x3000
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
		originalBar := -1
		steps := 0
		for ; steps < 3000000; steps++ {
			cur := originalRegs(original)
			if cur.CS == cs && cur.IP == 0x100 {
				break
			}
			_, known := functions[cur.IP]
			if cur.CS == cs && known || cur.CS == farSeg && cur.IP == 0 {
				name := fmt.Sprintf("sub_%X", 0x10000+uint32(cur.IP))
				if cur.CS == farSeg {
					name = "fixture-far"
				}
				entries[name]++
				if c.compareBar && cur.CS == cs && cur.IP == 0xaaa {
					originalBar = int(byte(cur.CX))
				}
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
		C.rect_run(&m, C.uint16_t(c.target), &fixture)
		if fixture.trace.calls > 8192 || fixture.unsupported != 0 {
			panic(fmt.Sprintf("unsupported/trace %x", fixture.unsupported))
		}
		var cTrace []byte
		cBar := -1
		for i := 0; i < int(fixture.trace.calls); i++ {
			s := fixture.trace.trace[i]
			if c.compareBar && s.target == 0xaaa {
				cBar = int(byte(s.regs[2]))
			}
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
		goBar := -1
		if c.compareBar {
			men, health := battleSideBarLengths(int(c.ax), int(c.ax))
			goBar = men
			if c.target == 0xc78e {
				goBar = health
			}
			goCases++
		}
		if want != got || !bytes.Equal(original.Mem, cMemory) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(originalState, cState) || !bytes.Equal(originalPorts, cPorts) || !bytes.Equal(originalTrace, cTrace) || c.compareBar && (originalBar != cBar || originalBar != goBar) {
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": r, "original": want, "c": got, "original_bar": originalBar, "c_bar": cBar, "go_bar": goBar, "original_trace": hex.EncodeToString(originalTrace), "c_trace": hex.EncodeToString(cTrace), "original_ports": hex.EncodeToString(originalPorts), "c_ports": hex.EncodeToString(cPorts), "original_device": hex.EncodeToString(originalState), "c_device": hex.EncodeToString(cState), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(active.VGA.Raw())), "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(cMemory))}
			return false
		}
		audits++
		if cases%64 == 0 || c.realAsset {
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
	all := []rectCase{}
	add := func(c rectCase) { all = append(all, c) }
	base := func(group string, t uint16) rectCase {
		return rectCase{group: group, target: t, mask: 15, seed: 3, bp: 0x6789, ax: 0xc00}
	}
	coords := [][4]int{{4, 4, 4, 4}, {5, 5, 6, 6}, {0, 0, 639, 399}, {-1, -1, 8, 8}, {632, 392, 640, 400}, {-10, 100, 20, 110}, {100, -10, 110, 20}, {630, 100, 660, 110}, {100, 390, 110, 410}, {-1, 0, -1, 10}, {640, 0, 640, 10}, {0, 400, 10, 400}, {-32768, -32768, 32767, 32767}, {8, 8, 24, 10}, {7, 7, 8, 8}, {15, 15, 16, 16}}
	for _, t := range []uint16{0xf020, 0xf1a3} {
		for _, xy := range coords {
			for _, reverse := range []bool{false, true} {
				for _, color := range []uint16{0, 1, 7, 12, 15, 255} {
					for _, mask := range []byte{0, 1, 5, 15} {
						for _, df := range []bool{false, true} {
							c := base("rectangle", t)
							c.ax = color<<8 | 0x5a
							c.dx = uint16(xy[0])
							c.bx = uint16(xy[1])
							c.si = uint16(xy[2])
							c.di = uint16(xy[3])
							if reverse {
								c.dx, c.si = c.si, c.dx
								c.bx, c.di = c.di, c.bx
							}
							c.mask = mask
							c.df = df
							c.readMap = byte(len(all) % 4)
							add(c)
						}
					}
				}
			}
		}
	}
	for _, t := range []uint16{0xf140, 0xf17b} {
		for _, w := range []uint16{1, 7, 8, 9, 16, 17, 24, 31, 80} {
			for _, pos := range []uint16{0, 79, 65535} {
				for _, df := range []bool{false, true} {
					c := base("edge", t)
					c.si = w
					c.di = pos
					c.bp = 3
					c.ax = (w / 8) | 0x800
					c.bx = 0x8080
					c.dx = 0xa55a
					c.df = df
					add(c)
				}
			}
		}
	}
	for _, t := range []uint16{0xcac, 0xcc3} {
		for _, ax := range []uint16{0, 1, 0xff08, 0xffff} {
			for _, df := range []bool{false, true} {
				c := base("mode", t)
				c.ax = ax
				c.dx = 0xbeef
				c.df = df
				add(c)
			}
		}
	}
	for n := 0; n < 65536; n++ {
		for _, t := range []uint16{0xc775, 0xc78e} {
			c := base("gauge", t)
			c.ax = uint16(n)
			c.bx = 72
			c.compareBar = true
			c.seed = byte(n % 8)
			add(c)
		}
	}
	for _, filled := range []uint16{0, 1, 7, 8, 9, 76, 124, 255} {
		for _, cap := range []uint16{0, 1, 76, 124, 255} {
			for _, x := range []uint16{0, 7, 498, 65535} {
				c := base("bar", 0xaaa)
				c.bx = 396
				c.dx = x
				c.cx = cap<<8 | filled
				c.ax = 0xc00
				add(c)
			}
		}
	}
	for _, w := range []uint16{0, 1, 7, 8, 9, 16, 17, 255} {
		for _, height := range []uint16{1, 2, 0} {
			for a := 0; a < 8; a++ {
				c := base("span", 0xad9)
				c.dx = uint16(a)
				c.cx = height<<8 | w
				c.bx = 65535
				c.df = a%2 == 1
				add(c)
			}
		}
	}
	for _, t := range []uint16{0xc61f, 0xc6bf, 0xc6ae} {
		limit := 16
		if t == 0xc6bf {
			limit = 6
		}
		if t == 0xc6ae {
			limit = 1
		}
		for slot := 0; slot < limit; slot++ {
			for _, color := range []uint16{0, 10, 12, 15} {
				for _, df := range []bool{false, true} {
					c := base("selection", t)
					c.ax = color<<8 | uint16(slot)
					c.df = df
					c.realAsset = true
					add(c)
				}
			}
		}
	}
	for _, xy := range [][2]uint16{{0, 0}, {5, 4}, {27, 10}} {
		for _, size := range []uint16{0x202, 0x406, 0x507} {
			for _, df := range []bool{false, true} {
				c := base("window-fill", 0xbcd)
				c.dx = xy[0]
				c.bx = xy[1]
				c.cx = size
				c.df = df
				c.realAsset = true
				add(c)
			}
		}
	}
	for _, xy := range [][2]uint16{{479, 382}, {479, 383}, {480, 383}, {0, 399}, {65535, 65535}} {
		for _, service := range []uint16{0, 1, 65535} {
			for _, n := range []uint16{0, 3, 167, 343, 65535} {
				c := base("sidebar", 0xc6f6)
				c.mouseX = xy[0]
				c.mouseY = xy[1]
				c.service = service
				c.realAsset = true
				c.values = [10]uint16{0, 1, 76, 77, 255, 33, n, uint16(n + 1), n, uint16(n + 1)}
				add(c)
			}
		}
	}
	for _, n := range []uint16{0, 1, 76, 77, 255} {
		c := base("waiting", 0xc74c)
		c.values = [10]uint16{n, n, 0, 255, 1, 76}
		c.realAsset = true
		add(c)
	}
	for glyph := 0; glyph < 9; glyph++ {
		for slot := 0; slot < 6; slot++ {
			c := base("actual-command", 0xc673)
			c.ax = uint16(slot<<8 | glyph)
			c.realAsset = true
			add(c)
		}
	}
	for _, mask := range []byte{0, 15} {
		c := base("actual-buttons", 0xc7f4)
		c.mask = mask
		c.realAsset = true
		add(c)
	}
	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{"schema": "wolong-c-rect-parity-v1", "input_sha256": inputHash, "icon_sha256": assetHash, "routine_sha256": routineHashes, "cases": cases, "groups": groups, "full_ram_plane_audits": audits, "indexed_content_audits": pixelAudits, "go_bar_cases": goCases, "content_top": 40, "content_height": 400, "entries_seen": entries, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "trace_entry_size": 90, "platform": "Independent pinned dosgolem devices; native C renderer; fixed FAR controls; IF/TF=0; CPU-model undefined flags; no wall-clock claim", "asset_banks": "ICONGRF segment1 at 3000; whole segment3 at 2200; command 2200, border 226c", "c_machine_code_match": false}
	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if e := os.WriteFile(*output, b, 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Rectangle original/C %d; complete RAM/plane %d; Go bars %d; content %d; pass=%v\n", cases, audits, goCases, pixelAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
