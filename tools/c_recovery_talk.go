//go:build matching_talk

// Original TALK/font/number/portrait control; independent RAM, VGA and DOS services.
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
#include "/repo/tools/c_recovery/talk.c"
#include "/repo/tools/c_recovery/talk_fixture.h"
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
var functions = map[uint16]int{0x6f9: 4, 0x75b: 119, 0x7d2: 115, 0x84a: 90, 0x8853: 48, 0x89a4: 58, 0xe38c: 26, 0xf4df: 73, 0x8b2: 41, 0x8db: 41, 0x904: 53, 0x939: 34, 0x95b: 35, 0x97e: 6, 0x984: 43, 0x62f: 107, 0x69a: 68, 0x6de: 23, 0xf75e: 67, 0xf7a4: 212, 0x701: 90, 0xf878: 16}

func entryName(t uint16) string {
	for _, v := range []uint16{0x8b2, 0x8db, 0x904, 0x939, 0x95b, 0x97e, 0x984, 0x701, 0xf75e} {
		if t == v {
			return fmt.Sprintf("code_%X", 0x10000+uint32(t))
		}
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

type talkCase struct {
	group                              string
	target, ax, bx, cx, dx, si, di, bp uint16
	scenario                           int
	message                            []byte
	parameters                         []uint16
	swapTable, keep, df, badFile       bool
	mask, readMap, seed                byte
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
	if err = active.LoadEXE(raw); err != nil {
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
	talk, err := os.ReadFile("/orig/TALK.DAT")
	if err != nil {
		panic(err)
	}
	worlds, err := os.ReadFile("/orig/SINARIO.DAT")
	if err != nil {
		panic(err)
	}
	portraits, err := os.ReadFile("/orig/KAOGRF.DAT")
	if err != nil {
		panic(err)
	}
	identities := map[string]struct {
		Size int
		Hash string
	}{"TALK.DAT": {34182, "08a22e09791d0a6ec2968e87d8655e12c91b45e00fae460b28593b35ff85e384"}, "SINARIO.DAT": {88832, "21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c"}, "KAOGRF.DAT": {307200, "b9c7745e3ed9b32f0c12003fe81f82756136f738427dc421ec96f6c4ed5c4ac8"}}
	for name, data := range map[string][]byte{"TALK.DAT": talk, "SINARIO.DAT": worlds, "KAOGRF.DAT": portraits} {
		id := identities[name]
		if len(data) != id.Size || fmt.Sprintf("%x", sha256.Sum256(data)) != id.Hash {
			panic(name + " identity")
		}
	}
	var originalAPI, cAPI []byte
	installTrace := func(device *machine.Machine, dst *[]byte) {
		previous := device.CPU.IntHook
		device.CPU.IntHook = func(c *cpu.CPU, n uint8) bool {
			if n == 0x21 {
				*dst = append(*dst, n)
				*dst = append(*dst, regBytes(originalRegs(device))...)
			}
			result := previous(c, n)
			if n == 0x21 {
				*dst = append(*dst, regBytes(originalRegs(device))...)
			}
			return result
		}
	}
	installTrace(original, &originalAPI)
	installTrace(active, &cAPI)
	initialRAM := append([]byte(nil), original.Mem...)
	for _, base := range []int{0x22000, 0x30000, 0x50000, 0x52000, 0x60000, 0x68000} {
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
	putWord(initialRAM, cs, 0xd38, 0x3500)
	putWord(initialRAM, cs, 0xd40, 0x6000)
	putWord(initialRAM, cs, 0xd52, 0x5000)
	putWord(initialRAM, cs, 0xd54, 0x2284)
	putWord(initialRAM, cs, 0xe479, 0x6800)
	copy(initialRAM[0x35000:], talk)
	for _, v := range [][2]uint16{{0xf78a, 0x410}, {0xf78c, 0x80}, {0xf791, 0x414}, {0xf793, 0x80}} {
		putWord(initialRAM, cs, v[0], v[1])
	}
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
	fixture := C.KiTalkFixture{}
	groups, entries := map[string]int{}, map[string]int{}
	corpusSlots, dosAPICalls := map[int]int{}, map[string]int{}
	cacheReads := []int{}
	sampleCounts := map[uint16]int{}
	cases, audits, pixelAudits := 0, 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	check := func(c talkCase) bool {
		if *only != "" && c.group != *only {
			return true
		}
		if *smoke && sampleCounts[c.target] >= 16 {
			return true
		}
		if (c.target == 0x75b || c.target == 0x8853) && c.cx != 0xffff {
			index := uint32(c.cx)
			if index >= 0x196 {
				index = 0x196 + (index-0x196)*8 + uint32(c.ax>>8)
			}
			if index >= 1022 {
				panic(fmt.Sprintf("invalid corpus fixture index=%d input=%+v", index, c))
			}
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
			{
				copy(original.Mem[0x22000:], icons[0x9700:])
				copy(cMemory[0x22000:], icons[0x9700:])
				copy(original.Mem[0x30000:], icons[0x2800:0x6700])
				copy(cMemory[0x30000:], icons[0x2800:0x6700])
			}
		}
		for _, mem := range [][]byte{original.Mem, cMemory} {
			if !c.keep {
				copy(mem[0x50000:], worlds[c.scenario*22208:(c.scenario+1)*22208])
				mem[int(cs)*16+0x845] = 0
				for i := 0; i < 4; i++ {
					mem[int(cs)*16+0x846+i] = 0xff
				}
			}
			putWord(mem, cs, 0xcfd, 0)
			for i, v := range c.parameters {
				putWord(mem, 0x4000, uint16(0x9000+i*2), v)
			}
			if c.message != nil {
				copy(mem[0x35000+0x9000:], append(append([]byte{}, c.message...), 0, 0))
			}
			if c.swapTable {
				a := binary.LittleEndian.Uint16(mem[int(cs)*16+0x8a4:])
				b := binary.LittleEndian.Uint16(mem[int(cs)*16+0x8a6:])
				putWord(mem, cs, 0x8a4, b)
				putWord(mem, cs, 0x8a6, a)
			}
			if c.badFile {
				copy(mem[int(cs)*16+0xd79:], []byte("NOFILE.DAT\x00"))
			}
		}
		originalAPI = nil
		cAPI = nil
		original.PortLog = nil
		active.PortLog = nil
		flags := uint16(2 | uint16(c.seed&1))
		if c.df {
			flags |= 0x400
		}
		r := registers{AX: c.ax, BX: c.bx, CX: c.cx, DX: c.dx, SI: c.si, DI: c.di, BP: c.bp, SP: 0x8000, DS: 0x3500, ES: 0xa0c8, SS: 0x4000, CS: cs, IP: c.target, Flags: flags}
		if c.target == 0xf4df || c.target == 0xe38c {
			r.DS = cs
		}
		if c.target == 0x84a || c.target == 0x6f9 {
			r.DS = 0x3500
			r.SI = 0x9000
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
				entries[name]++
				if c.group == "corpus" && cur.CS == cs && cur.IP == 0x6f9 {
					slot := int(c.cx)
					if slot >= 0x196 {
						slot = 0x196 + (slot-0x196)*8 + int(c.ax>>8)
					}
					if cur.SI != binary.LittleEndian.Uint16(talk[slot*2:]) {
						panic("original selected a different TALK slot")
					}
					corpusSlots[slot]++
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
		C.talk_run(&m, C.uint16_t(c.target), &fixture)
		if fixture.trace.calls > 8192 || fixture.unsupported != 0 {
			panic(fmt.Sprintf("unsupported=%x calls=%d input=%+v original=%+v C=%+v", fixture.unsupported, fixture.trace.calls, r, originalRegs(original), cRegs(&m)))
		}
		var cTrace []byte
		for i := 0; i < int(fixture.trace.calls); i++ {
			s := fixture.trace.trace[i]
			if _, known := functions[uint16(s.target)]; !known && !(s.regs[11] == machine.StubSeg && (s.target == 0x410 || s.target == 0x414)) {
				continue
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
		if want != got || !bytes.Equal(original.Mem, cMemory) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(originalState, cState) || !bytes.Equal(originalPorts, cPorts) || !bytes.Equal(originalTrace, cTrace) || !bytes.Equal(originalAPI, cAPI) {
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": r, "original": want, "c": got, "original_api": hex.EncodeToString(originalAPI), "c_api": hex.EncodeToString(cAPI), "original_trace": hex.EncodeToString(originalTrace), "c_trace": hex.EncodeToString(cTrace), "original_ports": hex.EncodeToString(originalPorts), "c_ports": hex.EncodeToString(cPorts), "original_device": hex.EncodeToString(originalState), "c_device": hex.EncodeToString(cState), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(active.VGA.Raw())), "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(cMemory))}
			return false
		}
		audits++
		if len(originalAPI)%57 != 0 {
			panic("DOS API frame")
		}
		for at := 0; at < len(originalAPI); at += 57 {
			dosAPICalls[fmt.Sprintf("INT21/AH=%02X", originalAPI[at+2])]++
		}
		if c.group == "cache" {
			cacheReads = append(cacheReads, len(originalAPI)/57)
		}
		{
			if !bytes.Equal(indexed(original), indexed(active)) {
				panic("content pixels")
			}
			pixelAudits++
		}
		for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), originalState, originalPorts, originalTrace, originalAPI} {
			originalDigest.Write(b)
		}
		for _, b := range [][]byte{regBytes(got), cMemory, active.VGA.Raw(), cState, cPorts, cTrace, cAPI} {
			cDigest.Write(b)
		}
		return true
	}
	all := []talkCase{}
	add := func(c talkCase) { all = append(all, c) }
	base := func(g string, t uint16) talkCase {
		return talkCase{group: g, target: t, mask: 15, seed: 3, bp: 0x6789, di: 0x9000, ax: 0xff00, bx: 10, dx: 80}
	}
	aliases := []uint16{}
	for i := 0; i < 127; i++ {
		off := 0x4240 + i*32
		if worlds[off+8] >= 0x80 && !bytes.Equal(worlds[off+2:off+8], worlds[off+8:off+14]) {
			aliases = append(aliases, uint16(i))
		}
	}
	if len(aliases) == 0 {
		panic("alias positive control")
	}
	for scenario := 0; scenario < 4; scenario++ {
		for _, marker := range []uint16{0x8b2, 0x8db, 0x904, 0x939, 0x95b, 0x97e, 0x984} {
			for _, mode := range []int{0, 1} {
				c := base("markers", marker)
				c.scenario = scenario
				c.parameters = []uint16{0xff00 | aliases[0]}
				if marker == 0x8db {
					c.parameters[0] = 0xff00 | 12
				}
				if marker == 0x904 {
					c.parameters[0] = 0xff00 | 3
				}
				if marker == 0x984 {
					c.parameters[0] = 0x8000
				}
				if mode == 1 {
					switch marker {
					case 0x8b2:
						c.parameters[0] = 0x4240 + aliases[0]*32
					case 0x8db:
						c.parameters[0] = 0x840 + 12*32
					case 0x904:
						c.parameters[0] = 3 * 64
					case 0x984:
						c.parameters[0] = 12345
					}
				}
				add(c)
			}
		}
	}
	for _, text := range [][]byte{{'A', 'B', 0, 'C', 'D'}, {0xa5, 0x5c, 'A', 'B'}, {'A', 0xa4, 0x40, 'B'}, {'\\', '1', 'A', '\\', '6', 'B', '\\', '7'}, {'\\', '1', '\\', '2', '\\', '3', 0, '\\', '4', '\\', '5', '\\', '6', '\\', '7'}} {
		for _, t := range []uint16{0x84a, 0x6f9} {
			for _, swap := range []bool{false, true} {
				c := base("strings", t)
				c.message = text
				c.parameters = []uint16{0xff00 | aliases[0], 0xff00 | 12, 0xff00 | 3, 777, 12345}
				c.swapTable = swap
				add(c)
			}
		}
	}
	for index := 0; index < 1022; index++ {
		c := base("corpus", 0x75b)
		c.dx = 0
		c.ax = uint16(index % 127)
		c.cx = uint16(index)
		if index >= 0x196 {
			c.cx = 0x196 + uint16((index-0x196)/8)
			c.ax |= uint16((index-0x196)%8) << 8
		}
		off := int(binary.LittleEndian.Uint16(talk[index*2:]))
		end := len(talk)
		if index < 1021 {
			end = int(binary.LittleEndian.Uint16(talk[(index+1)*2:]))
		}
		for i := off; i < end; {
			v := talk[i]
			if v >= 0x81 && v <= 0xfe {
				i += 2
				continue
			}
			if v == '\\' && i+1 < end {
				mark := talk[i+1]
				switch mark {
				case '1':
					c.parameters = append(c.parameters, 0xff00|aliases[0])
				case '2':
					c.parameters = append(c.parameters, 0xff00|12)
				case '3':
					c.parameters = append(c.parameters, 0xff00|3)
				case '6':
					c.parameters = append(c.parameters, 777)
				case '7':
					c.parameters = append(c.parameters, 12345)
				}
				i += 2
				continue
			}
			i++
		}
		add(c)
	}
	for portrait := 0; portrait < 150; portrait++ {
		c := base("portrait", 0x7d2)
		c.ax = uint16(portrait)
		c.bx = 100*80 + 20
		add(c)
	}
	for _, id := range []uint16{0, 1, 2, 3, 0, 4, 1, 5, 2, 6, 3, 7, 7} {
		c := base("cache", 0x7d2)
		c.ax = id
		c.bx = 100*80 + 20
		c.keep = len(all) > 0 && all[len(all)-1].group == "cache"
		add(c)
	}
	for _, index := range []uint16{0, 100, 409, 410, 411, 1021, 0xffff} {
		c := base("status", 0x8853)
		c.ax = 0
		c.cx = index
		if index >= 0x196 && index != 0xffff {
			c.cx = 0x196 + (index-0x196)/8
			c.ax |= ((index - 0x196) % 8) << 8
		}
		c.parameters = []uint16{0xff00 | aliases[0], 0xff00 | 12, 0xff00 | 3, 777, 12345}
		add(c)
	}
	for _, t := range []uint16{0xe38c, 0xf4df} {
		for _, n := range []uint16{0, 1, 2048, 4096} {
			c := base("file", t)
			c.ax = 0x800
			c.cx = 0
			c.bx = 0x6000
			c.dx = 0xd79
			c.si = 0x200
			c.di = n
			add(c)
		}
	}
	c := base("file", 0xf4df)
	c.badFile = true
	c.bx = 0x6000
	c.dx = 0xd79
	c.si = 0x200
	c.di = 2048
	add(c)
	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{"schema": "wolong-c-talk-parity-v1", "input_sha256": inputHash, "icon_sha256": assetHash, "routine_sha256": routineHashes, "cases": cases, "groups": groups, "full_ram_plane_audits": audits, "indexed_content_audits": pixelAudits, "content_top": 40, "content_height": 400, "entries_seen": entries, "corpus_slots_seen": corpusSlots, "dos_api_calls": dosAPICalls, "cache_sequence_api_counts": cacheReads, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "trace_entry_size": 90, "platform": "Independent pinned DOS/font/VGA services; native C TALK/number/portrait control; IF/TF=0; CPU-model undefined flags; no original INT50 UI or wall-clock claim", "asset_banks": "TALK at 3500; SINARIO at 5000; portraits at 6000; UI map at 6800; ICONGRF at 2200/3000", "additional_assets": identities, "original_font_calls": originalDOS.Font.Calls, "c_font_calls": activeDOS.Font.Calls, "original_missing_fonts": originalDOS.Font.Missing, "c_missing_fonts": activeDOS.Font.Missing, "c_machine_code_match": false}
	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if e := os.WriteFile(*output, b, 0644); e != nil {
		panic(e)
	}
	fmt.Printf("TALK original/C %d; complete RAM/plane %d; content %d; pass=%v\n", cases, audits, pixelAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
