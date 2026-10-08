//go:build matching_vga

// Eight original DOS/V VGA entries; the private build module is set by the wrapper.
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/rng.c"
#include "/repo/tools/c_recovery/economy.c"
#include "/repo/tools/c_recovery/settlement.c"
#include "/repo/tools/c_recovery/vga.c"
#include "/repo/tools/c_recovery/vga_fixture.h"
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
	"image"
	"image/color"
	"image/png"
	"os"
	"unsafe"
)

var activeDevice *machine.Machine
var functions = map[uint16]int{0x9796: 45, 0x97c3: 45, 0xf9b0: 107, 0xfa1b: 28, 0xfa37: 107, 0xfaa2: 32, 0xfac2: 79, 0xfb11: 24}

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

type vgaCase struct {
	group                            string
	target, size, source, dest, x, y uint16
	mask, mode, logic, readMap, seed byte
	df                               bool
	keep                             bool
	asset                            []byte
	assetIndex                       int
}

func saveIndexed(path string, px []byte) {
	pal := color.Palette{}
	for i := 0; i < 16; i++ {
		pal = append(pal, color.RGBA{byte(i * 17), byte(i * 17), byte(i * 17), 255})
	}
	im := image.NewPaletted(image.Rect(0, 0, 640, 400), pal)
	copy(im.Pix, px)
	f, e := os.Create(path)
	if e != nil {
		panic(e)
	}
	if e = png.Encode(f, im); e != nil {
		panic(e)
	}
	if e = f.Close(); e != nil {
		panic(e)
	}
}

func main() {
	output := flag.String("out", "/output/results/O2.json", "receipt")
	only := flag.String("group", "", "group")
	smoke := flag.Bool("smoke", false, "limited matrix")
	flag.Parse()
	raw, err := os.ReadFile("/orig/KI.EXE")
	if err != nil {
		panic(err)
	}
	inputHash := fmt.Sprintf("%x", sha256.Sum256(raw))
	if inputHash != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" {
		panic("input identity")
	}
	original := machine.New()
	if err = original.LoadEXE(raw); err != nil {
		panic(err)
	}
	activeDevice = machine.New()
	original.IRQ0Every = 0
	original.KeyEvery = 0
	original.VGAFrameEvery = 0
	activeDevice.IRQ0Every = 0
	activeDevice.KeyEvery = 0
	activeDevice.VGAFrameEvery = 0
	portraits, err := os.ReadFile("/orig/KAOGRF.DAT")
	if err != nil || len(portraits) != 150*2048 {
		panic("portrait inventory")
	}
	assetHash := fmt.Sprintf("%x", sha256.Sum256(portraits))
	cs := original.CPU.Seg[cpu.CS]
	routineHashes := map[string]string{}
	for target, size := range functions {
		b := raw[int(target)+512 : int(target)+512+size]
		at := int(cs)*16 + int(target)
		if !bytes.Equal(b, original.Mem[at:at+size]) {
			panic("routine identity")
		}
		routineHashes[fmt.Sprintf("sub_%X", 0x10000+uint32(target))] = fmt.Sprintf("%x", sha256.Sum256(b))
	}
	initialRAM := append([]byte(nil), original.Mem...)
	for i := 0; i < 65536; i++ {
		initialRAM[0x22000+i] = byte(i*37 + i/80*13 + 53)
		initialRAM[0x60000+i] = byte(i*19 + i/64*7 + 91)
	}
	putWord(initialRAM, cs, 0x987c, 0x6000)
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	cMemory := unsafe.Slice((*byte)(native), 1<<20)
	activeDevice.Mem = cMemory
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	fixture := C.KiVgaFixture{}
	groups := map[string]int{}
	entries := map[string]int{}
	cases, audits, pixelAudits := 0, 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	check := func(c vgaCase) bool {
		if *only != "" && c.group != *only {
			return true
		}
		limit := 8
		if c.group == "roundtrip" {
			limit = 9 // Keep every sampled save/draw/restore triplet complete.
		}
		if *smoke && groups[c.group] >= limit {
			return true
		}
		if !c.keep {
			copy(original.Mem, initialRAM)
			copy(cMemory, initialRAM)
			original.SetVideoMode(0x12)
			activeDevice.SetVideoMode(0x12)
			for p := 0; p < 4; p++ {
				for i := 0; i < 65536; i++ {
					b := byte(i*37 + i/80*13 + p*61 + int(c.seed)*29 + 73)
					original.VGA.Planes[p][i] = b
					activeDevice.VGA.Planes[p][i] = b
				}
			}
			for _, device := range []*machine.Machine{original, activeDevice} {
				device.Out8(0x3c4, 2)
				device.Out8(0x3c5, c.mask)
				for _, q := range [][2]byte{{0, 0}, {1, 0}, {3, c.logic << 3}, {4, c.readMap}, {5, c.mode}, {8, 0xff}} {
					device.Out8(0x3ce, q[0])
					device.Out8(0x3cf, q[1])
				}
				device.Read8(0xa4321)
				device.PortLog = nil
			}
		}
		original.PortLog = nil
		activeDevice.PortLog = nil
		if c.asset != nil {
			copy(original.Mem[0x22000:], c.asset)
			copy(cMemory[0x22000:], c.asset)
		}

		flags := uint16(2)
		if c.df {
			flags |= 0x400
		}
		r := registers{AX: c.size, BX: c.dest, CX: c.size, DX: c.y, SI: c.source, DI: 0x101, BP: 0x6789, SP: 0x8000, DS: 0x2200, ES: 0xa0c8, SS: 0x4000, CS: cs, IP: c.target, Flags: flags}
		if c.target == 0xfac2 || c.target == 0xfb11 {
			r.ES = 0x6000
			r.DI = c.source
		}
		if c.target == 0xfb11 {
			r.DS = 0xa0c8
		}
		if c.target == 0x9796 || c.target == 0x97c3 {
			r.DX = c.x
			r.BX = c.y
		}
		cpuRegs := original.CPU
		cpuRegs.R = [8]uint16{r.AX, r.CX, r.DX, r.BX, r.SP, r.BP, r.SI, r.DI}
		cpuRegs.Seg = [4]uint16{r.ES, r.CS, r.SS, r.DS}
		cpuRegs.IP = r.IP
		cpuRegs.SetFlags(r.Flags)
		putWord(original.Mem, r.SS, r.SP, 0x100)
		putWord(cMemory, r.SS, r.SP, 0x100)
		assignC(&m, originalRegs(original))
		fixture.calls = 0
		var originalTrace []byte
		for step := 0; step < 2000000; step++ {
			current := originalRegs(original)
			if current.CS == cs && current.IP == 0x100 {
				break
			}
			if _, ok := functions[current.IP]; current.CS == cs && ok {
				entries[fmt.Sprintf("sub_%X", 0x10000+uint32(current.IP))]++
				originalTrace = append(originalTrace, byte(current.IP), byte(current.IP>>8))
				originalTrace = append(originalTrace, regBytes(current)...)
				for i := uint16(0); i < 32; i++ {
					originalTrace = append(originalTrace, original.Mem[uint32(current.SS)*16+uint32(current.SP+i)])
				}
				d := deviceBytes(original)
				originalTrace = append(originalTrace, d[:28]...)
			}
			if err := original.Step(); err != nil {
				panic(err)
			}
			if step == 1999999 {
				panic("bounded original routine")
			}
		}
		C.vga_run(&m, C.uint16_t(c.target), &fixture)
		if fixture.calls > 8192 {
			panic("trace capacity")
		}
		var cTrace []byte
		for i := 0; i < int(fixture.calls); i++ {
			s := fixture.trace[i]
			cTrace = append(cTrace, byte(s.target), byte(s.target>>8))
			for _, v := range s.regs {
				cTrace = append(cTrace, byte(v), byte(v>>8))
			}
			cTrace = append(cTrace, C.GoBytes(unsafe.Pointer(&s.stack[0]), 32)...)
			cTrace = append(cTrace, C.GoBytes(unsafe.Pointer(&s.device[0]), 28)...)
		}
		want, got := originalRegs(original), cRegs(&m)
		originalPorts, cPorts := portBytes(original), portBytes(activeDevice)
		originalState, cState := deviceBytes(original), deviceBytes(activeDevice)
		cases++
		groups[c.group]++
		if want != got || !bytes.Equal(original.Mem, cMemory) || !bytes.Equal(original.VGA.Raw(), activeDevice.VGA.Raw()) || !bytes.Equal(originalState, cState) || !bytes.Equal(originalPorts, cPorts) || !bytes.Equal(originalTrace, cTrace) {
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": r, "original": want, "c": got, "original_trace": hex.EncodeToString(originalTrace), "c_trace": hex.EncodeToString(cTrace), "original_ports": hex.EncodeToString(originalPorts), "c_ports": hex.EncodeToString(cPorts), "original_device": hex.EncodeToString(originalState), "c_device": hex.EncodeToString(cState), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(activeDevice.VGA.Raw())), "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(cMemory))}
			return false
		}
		audits++
		if c.asset != nil {
			saveIndexed(fmt.Sprintf("/output/results/asset-%d-%d-original.png", c.assetIndex, c.dest), indexed(original))
			saveIndexed(fmt.Sprintf("/output/results/asset-%d-%d-c.png", c.assetIndex, c.dest), indexed(activeDevice))
		}
		if cases%32 == 0 {
			if !bytes.Equal(indexed(original), indexed(activeDevice)) {
				panic("indexed pixels")
			}
			pixelAudits++
		}
		for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), originalState, originalPorts, originalTrace} {
			originalDigest.Write(b)
		}
		for _, b := range [][]byte{regBytes(got), cMemory, activeDevice.VGA.Raw(), cState, cPorts, cTrace} {
			cDigest.Write(b)
		}
		return true
	}
	base := func(group string, target, size uint16) vgaCase {
		return vgaCase{group: group, target: target, size: size, source: 1, dest: 1, x: 88, y: 184, mask: 15, mode: 0, logic: 2, readMap: 0, seed: 3}
	}
	configs := [][2]byte{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 0}, {2, 2}, {3, 3}}
	sizes := []uint16{0x101, 0x102, 0x804, 0x808, 0x4004, 0x100, 0x001}
	for _, target := range []uint16{0xf9b0, 0xfa37} {
		for _, size := range sizes {
			for _, mask := range []byte{0, 1, 3, 5, 15} {
				for _, dest := range []uint16{0, 79, 80, 32000, 65535} {
					for _, source := range []uint16{0, 1, 65500, 65535} {
						for _, df := range []bool{false, true} {
							c := base("draw", target, size)
							c.mask = mask
							c.dest = dest
							c.source = source
							c.df = df
							c.readMap = byte(source % 4)
							if !check(c) {
								goto finish
							}
						}
					}
				}
			}
		}
	}
	for _, target := range []uint16{0xfa1b, 0xfaa2, 0xfb11} {
		for _, size := range sizes {
			for _, cfg := range configs {
				for _, dest := range []uint16{0, 79, 80, 32000, 65535} {
					for _, source := range []uint16{0, 1, 65500, 65535} {
						for _, df := range []bool{false, true} {
							c := base("leaf", target, size)
							c.dest = dest
							c.source = source
							c.mode = cfg[0]
							c.logic = cfg[1]
							c.df = df
							if !check(c) {
								goto finish
							}
						}
					}
				}
			}
		}
	}
	for _, size := range []uint16{0x101, 0x102, 0x804, 0x4004, 0x180} {
		for _, dest := range []uint16{0, 79, 80, 32000, 65535} {
			for _, source := range []uint16{0, 1, 65500, 65535} {
				for _, df := range []bool{false, true} {
					c := base("save", 0xfac2, size)
					c.dest = dest
					c.source = source
					c.df = df
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	for _, target := range []uint16{0x9796, 0x97c3} {
		for _, size := range []uint16{0x101, 0x102, 0x804, 0x808, 0x4004} {
			for _, x := range []uint16{0, 7, 8, 639, 65535} {
				for _, y := range []uint16{0, 1, 184, 399, 65535} {
					for read := 0; read < 4; read++ {
						c := base("wrapper", target, size)
						c.x = x
						c.y = y
						c.readMap = byte(read)
						c.df = read%2 == 1
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	for _, index := range []int{0, 5, 49, 149} {
		for _, dest := range []uint16{0, 1280} {
			c := base("real-assets", 0xfa37, 0x4004)
			c.dest = dest
			c.source = 0
			c.asset = portraits[index*2048 : (index+1)*2048]
			c.assetIndex = index
			if !check(c) {
				goto finish
			}
		}
	}
	for _, size := range []uint16{0x804, 0x808, 0x4004} {
		for _, xy := range [][2]uint16{{0, 0}, {88, 184}, {296, 184}} {
			for _, seed := range []byte{0, 7} {
				c := base("roundtrip", 0x9796, size)
				c.x = xy[0]
				c.y = xy[1]
				c.seed = seed
				if !check(c) {
					goto finish
				}
				before := append([]byte(nil), original.VGA.Raw()...)
				c.target = 0xfa37
				c.keep = true
				c.dest = xy[1]*80 + (xy[0] >> 3)
				c.source = 0
				if !check(c) {
					goto finish
				}
				c.target = 0x97c3
				if !check(c) {
					goto finish
				}
				if !bytes.Equal(before, original.VGA.Raw()) || !bytes.Equal(before, activeDevice.VGA.Raw()) {
					panic("save/draw/restore full plane roundtrip")
				}
			}
		}
	}

finish:
	report := map[string]any{"schema": "wolong-c-vga-parity-v1", "input_sha256": inputHash, "portrait_sha256": assetHash, "routine_sha256": routineHashes, "cases": cases, "groups": groups, "full_ram_plane_audits": audits, "indexed_pixel_audits": pixelAudits, "entries_seen": entries, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "trace_entry_size": 90, "platform": "Two independent pinned dosgolem Machine/VGA instances; original guest CPU and C native algorithm share only mature platform contract; IF/TF=0; no wall-clock claim", "c_machine_code_match": false}
	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if err = os.WriteFile(*output, b, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("VGA original/C %d; full RAM/plane %d; indexed %d; pass=%v\n", cases, audits, pixelAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
