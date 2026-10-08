//go:build matching_aligned

// Five original aligned blit/caller entries; private module and platform are pinned by the wrapper.
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
#include "/repo/tools/c_recovery/aligned_fixture.h"
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
var functions = map[uint16]int{0x1b4: 33, 0x1db: 51, 0xc14: 76, 0xc60: 23, 0xc77: 53, 0x895d: 71, 0x89de: 18, 0xc673: 59, 0xc7f4: 74, 0xd5d4: 65, 0xe3d7: 68, 0xe453: 38, 0xf888: 176, 0xf938: 97, 0xf999: 23, 0xf9b0: 107, 0xfa1b: 28, 0xfa37: 107, 0xfaa2: 32}

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

type alignedCase struct {
	group                                                 string
	target, size, source, x, y, ax, packed, dest, service uint16
	mask, mode, logic, readMap, seed                      byte
	df, realAsset, keep                                   bool
}

func main() {
	output := flag.String("out", "/output/results/O2.json", "receipt")
	only := flag.String("group", "", "group")
	smoke := flag.Bool("smoke", false, "bounded preview matrix")
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
	if assetHash != "2154782c045b898aa5fafa74a4ff0c3745771ec85799b7882e1c4c009b3f1c1d" || len(icons) < 0x6700 {
		panic("icon inventory")
	}
	original, active := machine.New(), machine.New()
	activeDevice = active
	if err = original.LoadEXE(raw); err != nil {
		panic(err)
	}
	for _, device := range []*machine.Machine{original, active} {
		device.IRQ0Every = 0
		device.KeyEvery = 0
		device.VGAFrameEvery = 0
	}
	cs := original.CPU.Seg[cpu.CS]
	routineHashes := map[string]string{}
	header := int(binary.LittleEndian.Uint16(raw[8:])) * 16
	relocCount, relocTable := int(binary.LittleEndian.Uint16(raw[6:])), int(binary.LittleEndian.Uint16(raw[24:]))
	for target, size := range functions {
		start := header + int(target)
		file := raw[start : start+size]
		expected := append([]byte(nil), file...)
		for i := 0; i < relocCount; i++ {
			at := relocTable + i*4
			off := header + int(binary.LittleEndian.Uint16(raw[at:])) + 16*int(binary.LittleEndian.Uint16(raw[at+2:]))
			if off >= start && off < start+size {
				binary.LittleEndian.PutUint16(expected[off-start:], binary.LittleEndian.Uint16(raw[off:])+cs)
			}
		}
		at := int(cs)*16 + int(target)
		if !bytes.Equal(expected, original.Mem[at:at+size]) {
			panic(fmt.Sprintf("loaded routine %x differs", target))
		}
		routineHashes[fmt.Sprintf("sub_%X", 0x10000+uint32(target))] = fmt.Sprintf("%x", sha256.Sum256(file))
	}
	initialRAM := append([]byte(nil), original.Mem...)
	for _, base := range []int{0x22000, 0x30000, 0x50000, 0x60000} {
		for i := 0; i < 65536; i++ {
			initialRAM[base+i] = byte(i*37 + i/80*13 + base/4096 + 73)
		}
	}
	putWord(initialRAM, cs, 0xd48, 0x2200)
	putWord(initialRAM, cs, 0xd4a, 0x3000)
	putWord(initialRAM, cs, 0xd84e, 0x5000)
	putWord(initialRAM, cs, 0xe479, 0x6000)
	putWord(initialRAM, cs, 0xe47b, 0)
	farSeg := binary.LittleEndian.Uint16(initialRAM[int(cs)*16+0x2082:])
	// Explicit fixed device fixture: cmp AX,2; read three control words; cmp AX,AX; RETF.
	stub := []byte{0x83, 0xf8, 2, 0x75, 0x0e, 0x2e, 0xa1, 0x10, 0x70, 0x2e, 0x8b, 0x16, 0x12, 0x70, 0x2e, 0x8b, 0x1e, 0x14, 0x70, 0x39, 0xc0, 0xcb}
	copy(initialRAM[int(farSeg)*16:], stub)
	putWord(initialRAM, farSeg, 0x7012, 0x3456)
	putWord(initialRAM, farSeg, 0x7014, 0x789a)
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	cMemory := unsafe.Slice((*byte)(native), 1<<20)
	active.Mem = cMemory
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	fixture := C.KiAlignedFixture{}
	groups, entries := map[string]int{}, map[string]int{}
	cases, audits, pixelAudits, buttonQueries := 0, 0, 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	check := func(c alignedCase) bool {
		if *only != "" && c.group != *only {
			return true
		}
		if *smoke && groups[c.group] >= 12 && c.group != "buttons" && c.group != "button-query" {
			return true
		}
		if !c.keep {
			copy(original.Mem, initialRAM)
			copy(cMemory, initialRAM)
			original.SetVideoMode(0x12)
			active.SetVideoMode(0x12)
			for p := 0; p < 4; p++ {
				for i := 0; i < 65536; i++ {
					b := byte(i*37 + i/80*13 + p*61 + int(c.seed)*29 + 73)
					original.VGA.Planes[p][i] = b
					active.VGA.Planes[p][i] = b
				}
			}
			for _, dev := range []*machine.Machine{original, active} {
				dev.Out8(0x3c4, 2)
				dev.Out8(0x3c5, c.mask)
				for _, q := range [][2]byte{{0, 0}, {1, 0}, {3, c.logic << 3}, {4, c.readMap}, {5, c.mode}, {8, 0xff}} {
					dev.Out8(0x3ce, q[0])
					dev.Out8(0x3cf, q[1])
				}
				dev.Read8(0xa4321)
			}
			if c.realAsset {
				copy(original.Mem[0x22000:], icons[0x2800:0x6700])
				copy(cMemory[0x22000:], icons[0x2800:0x6700])
				copy(original.Mem[0x30000:], icons[:0x2800])
				copy(cMemory[0x30000:], icons[:0x2800])
			}
		}
		original.PortLog = nil
		active.PortLog = nil
		putWord(original.Mem, farSeg, 0x7010, c.service)
		putWord(cMemory, farSeg, 0x7010, c.service)
		flags := uint16(2 | uint16(c.seed&1) | uint16(c.seed&2)<<10)
		if c.df {
			flags |= 0x400
		}
		r := registers{AX: c.ax, BX: c.y, CX: c.size, DX: c.x, SI: c.source, DI: 0x101, BP: 0x6789, SP: 0x8000, DS: 0x2200, ES: 0xa0c8, SS: 0x4000, CS: cs, IP: c.target, Flags: flags}
		if c.target == 0xf938 {
			r.BX = c.packed
			r.BP = c.dest
		}
		if c.target == 0xf999 {
			r.DX = c.x
		}
		if c.target == 0xc673 {
			r.DS = 0x3500
		}
		if c.target == 0xe453 {
			r.CX = c.x
			r.DX = c.y
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
			if cur.CS == cs && known || cur.CS == farSeg && cur.IP == 0 {
				name := fmt.Sprintf("sub_%X", 0x10000+uint32(cur.IP))
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
			if err := original.Step(); err != nil {
				panic(err)
			}
		}
		if steps == 3000000 {
			panic("bounded original routine")
		}
		C.aligned_run(&m, C.uint16_t(c.target), &fixture)
		if fixture.trace.calls > 8192 || fixture.unsupported != 0 {
			panic(fmt.Sprintf("trace or unsupported %x", fixture.unsupported))
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
		if c.group == "button-query" {
			slot := int(c.x) / 80
			hit := initialRAM[int(cs)*16+0xd2e4+slot] + 0x15
			if byte(want.AX) != hit {
				panic("raw button hit consumer")
			}
			buttonQueries++
		}
		if cases%32 == 0 || c.realAsset {
			if !bytes.Equal(indexed(original), indexed(active)) {
				panic("cropped indexed content")
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
	all := []alignedCase{}
	add := func(c alignedCase) { all = append(all, c) }
	for a := 0; a < 8; a++ {
		for _, w := range []uint16{1, 7, 8, 9, 16, 24, 31, 32, 63, 248, 249, 255, 0} {
			for _, height := range []uint16{1, 2, 16, 0} {
				for _, x := range []uint16{0, 632, 65528} {
					for _, df := range []bool{false, true} {
						i := len(all)
						add(alignedCase{group: "aligned", target: 0xf888, size: height<<8 | w, source: []uint16{0, 1, 65535}[i%3], x: x + uint16(a), y: []uint16{0, 399, 65535}[i%3], mask: []byte{0, 1, 5, 15}[i%4], seed: byte(i % 8), df: df})
					}
				}
			}
		}
	}
	for _, shift := range []uint16{0, 1, 7, 8, 31, 32, 255} {
		for _, bh := range []uint16{0, 1, 2, 3, 127, 128, 129, 255} {
			for _, bl := range []uint16{0, 1, 3, 255} {
				for _, height := range []uint16{1, 2, 16} {
					for _, mask := range []uint16{0, 0xa55a} {
						i := len(all)
						add(alignedCase{group: "row", target: 0xf938, size: height<<8 | shift, source: []uint16{0, 65535}[i%2], ax: mask, packed: bh<<8 | bl, dest: []uint16{0, 65535}[i%2], mask: 15, mode: 0, logic: byte((i / 2) % 4), seed: byte(i % 8), df: i%2 == 1})
					}
				}
			}
		}
	}
	for _, shift := range []uint16{0, 1, 7, 8, 31, 32, 255} {
		for _, bh := range []uint16{0, 1, 2, 3, 127, 128, 129, 255} {
			add(alignedCase{group: "zero-height", target: 0xf938, size: shift, source: 65535, ax: 0xf00f, packed: bh<<8 | 3, dest: 65535, mask: 15, mode: 0, logic: 2, df: true})
		}
	}
	for esr := 0; esr < 256; esr++ {
		for _, df := range []bool{false, true} {
			add(alignedCase{group: "plane", target: 0xf999, ax: 0xa55a, x: 0x8000 | uint16(esr), mask: 15, seed: byte(esr % 8), df: df})
		}
	}
	for glyph := 0; glyph < 6; glyph++ {
		for a := 0; a < 8; a++ {
			for _, mask := range []byte{5, 15} {
				add(alignedCase{group: "real-assets", target: 0xf888, size: 0x1018, source: uint16(0x3900 + glyph*192), x: uint16(glyph*80 + a), y: 374, mask: mask, seed: byte(glyph), realAsset: true})
			}
		}
	}
	for slot := 0; slot < 6; slot++ {
		for glyph := 0; glyph < 6; glyph++ {
			for _, service := range []uint16{0, 1, 65535} {
				for _, df := range []bool{false, true} {
					add(alignedCase{group: "redraw", target: 0xc673, ax: uint16(slot<<8 | glyph), source: 65535, mask: 15, service: service, seed: byte(glyph), realAsset: true, df: df})
				}
			}
		}
	}
	for _, mask := range []byte{0, 5, 15} {
		for _, df := range []bool{false, true} {
			add(alignedCase{group: "buttons", target: 0xc7f4, ax: 0xabcd, source: 65535, mask: mask, seed: 3, realAsset: true, df: df})
			for slot := 0; slot < 6; slot++ {
				add(alignedCase{group: "button-query", target: 0xe453, x: uint16(slot*80 + 40), y: 384, mask: mask, seed: 3, realAsset: true, df: df, keep: true})
			}
		}
	}
	for _, xy := range [][2]uint16{{0, 0}, {5, 4}, {27, 10}, {39, 24}} {
		for _, style := range []uint16{0, 0x1f, 0x51} {
			for _, size := range []uint16{0x202, 0x406, 0x507, 0xc0f} {
				for _, mask := range []byte{0, 15} {
					for _, df := range []bool{false, true} {
						add(alignedCase{group: "window", target: 0x895d, size: size, ax: style, x: xy[0], y: xy[1], mask: mask, seed: 7, realAsset: true, df: df})
					}
				}
			}
		}
	}
	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{"schema": "wolong-c-aligned-parity-v1", "input_sha256": inputHash, "icon_sha256": assetHash, "routine_sha256": routineHashes, "cases": cases, "groups": groups, "full_ram_plane_audits": audits, "indexed_content_audits": pixelAudits, "content_top": 40, "content_height": 400, "button_queries": buttonQueries, "entries_seen": entries, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "trace_entry_size": 90, "platform": "Independent pinned dosgolem devices; native C game algorithms; fixed FAR device fixture; IF/TF=0; undefined flags follow CPU model; no wall-clock claim", "c_machine_code_match": false}
	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(*output, b, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("Aligned original/C %d; complete RAM/plane %d; content pixels %d; pass=%v\n", cases, audits, pixelAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
