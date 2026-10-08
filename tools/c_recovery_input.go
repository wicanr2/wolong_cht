//go:build matching_input

// Original two-segment mouse/input control; independent DOS/VGA and explicit counter input.
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
#include "/repo/tools/c_recovery/input.c"
#include "/repo/tools/c_recovery/input_fixture.h"
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
var cPortReads []byte
var activeTicks int
var currentInput inputCase

type inputReadBus struct {
	cpu.Bus
	reads *[]byte
}

func (b inputReadBus) In8(port uint16) uint8 {
	v := b.Bus.In8(port)
	*b.reads = append(*b.reads, byte(port), byte(port>>8), v)
	return v
}

var functions = map[uint32]int{0x18810: 67, 0x121e7: 47, 0x20000: 14, 0x2002e: 65, 0x2006f: 1, 0x20070: 42, 0x2009a: 35, 0x200bd: 3, 0x200c0: 54, 0x20137: 52, 0x2016b: 50, 0x2019d: 41, 0x201c6: 30, 0x201e4: 40, 0x2020c: 61, 0x20249: 87, 0x202a0: 29, 0x202bd: 65, 0x202fe: 29, 0x1222b: 91, 0x20101: 54, 0x106f9: 4, 0x1075b: 119, 0x107d2: 115, 0x1084a: 90, 0x18853: 48, 0x189a4: 58, 0x1e38c: 26, 0x1f4df: 73, 0x108b2: 41, 0x108db: 41, 0x10904: 53, 0x10939: 34, 0x1095b: 35, 0x1097e: 6, 0x10984: 43, 0x1062f: 107, 0x1069a: 68, 0x106de: 23, 0x1f75e: 67, 0x1f7a4: 212, 0x10701: 90, 0x1f878: 16}

func entryName(loc uint32) string {
	if loc == 0x2006f {
		return "nullsub_4"
	}
	for _, v := range []uint32{0x1222b, 0x20101, 0x108b2, 0x108db, 0x10904, 0x10939, 0x1095b, 0x1097e, 0x10984, 0x10701, 0x1f75e} {
		if loc == v {
			return fmt.Sprintf("code_%X", loc)
		}
	}
	return fmt.Sprintf("sub_%X", loc)
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

func mouseBytes(d *dos.DOS) []byte {
	m := d.Mouse
	b, err := json.Marshal(map[string]any{"position": [2]uint16{m.X, m.Y}, "buttons": m.Buttons, "press": m.Press, "release": m.Release, "press_at": m.PressAt, "release_at": m.ReleaseAt, "mickey": [2]int16{m.MickeyX, m.MickeyY}, "range": [4]uint16{m.MinX, m.MaxX, m.MinY, m.MaxY}, "handler": m.Handler, "queries": m.PressQ, "calls": m.Calls})
	if err != nil {
		panic(err)
	}
	return b
}

type inputCase struct {
	group                           string
	target                          uint32
	ax, bx, cx, dx, si, di, bp      uint16
	scenario                        int
	message                         []byte
	parameters                      []uint16
	keep, df                        bool
	patch                           uint16
	cursor                          uint8
	oldX, oldY, mx, my              uint16
	presses                         [3]uint16
	clickAt, clickButton, timerHold int
	mask, readMap, seed             byte
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
		start := header + int(t-0x10000)
		file := raw[start : start+size]
		expected := append([]byte(nil), file...)
		for i := 0; i < int(binary.LittleEndian.Uint16(raw[6:])); i++ {
			at := int(binary.LittleEndian.Uint16(raw[24:])) + i*4
			off := header + int(binary.LittleEndian.Uint16(raw[at:])) + 16*int(binary.LittleEndian.Uint16(raw[at+2:]))
			if off >= start && off < start+size {
				binary.LittleEndian.PutUint16(expected[off-start:], binary.LittleEndian.Uint16(raw[off:])+cs)
			}
		}
		at := int(cs)*16 + int(t-0x10000)
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
	mouseCS := binary.LittleEndian.Uint16(original.Mem[int(cs)*16+0x2082:])
	var originalAPI, cAPI, originalPortReads []byte
	originalQueries, activeQueries, originalTicks := 0, 0, 0
	original.CPU.Bus = inputReadBus{Bus: original.CPU.Bus, reads: &originalPortReads}
	installTrace := func(device *machine.Machine, d *dos.DOS, dst *[]byte, queries *int) {
		previous := device.CPU.IntHook
		device.CPU.IntHook = func(c *cpu.CPU, n uint8) bool {
			if n == 0x33 && c.R[cpu.AX] == 5 {
				(*queries)++
				if currentInput.clickAt > 0 && *queries == currentInput.clickAt {
					d.Mouse.PressButton(currentInput.clickButton)
				}
			}
			if n == 0x21 || n == 0x33 {
				*dst = append(*dst, n)
				*dst = append(*dst, regBytes(originalRegs(device))...)
			}
			result := previous(c, n)
			if n == 0x21 || n == 0x33 {
				*dst = append(*dst, regBytes(originalRegs(device))...)
			}
			return result
		}
	}
	installTrace(original, originalDOS, &originalAPI, &originalQueries)
	installTrace(active, activeDOS, &cAPI, &activeQueries)
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
	fixture := C.KiInputFixture{}
	groups, entries := map[string]int{}, map[string]int{}
	mousePollCounts, waitCheckpoints := map[string]int{}, map[string]int{}
	sampleCounts := map[uint32]int{}
	cases, audits, pixelAudits := 0, 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	var roundtripPixels []byte
	roundtripRestored := 0
	check := func(c inputCase) bool {
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
			{
				copy(original.Mem[0x22000:], icons[0x9700:])
				copy(cMemory[0x22000:], icons[0x9700:])
				copy(original.Mem[0x30000:], icons[0x2800:0x6700])
				copy(cMemory[0x30000:], icons[0x2800:0x6700])
			}
		}
		if c.group == "roundtrip" && !c.keep {
			roundtripPixels = append([]byte{}, original.VGA.Raw()...)
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
			putWord(mem, cs, 0x2239, c.patch)
			if !c.keep {
				putWord(mem, mouseCS, 0xfc, c.oldX)
				putWord(mem, mouseCS, 0xfe, c.oldY)
				mem[int(mouseCS)*16+0x100] = c.cursor
			}
		}
		currentInput = c
		activeTicks = 0
		originalTicks = 0
		originalQueries = 0
		activeQueries = 0
		cPortReads = nil
		originalPortReads = nil
		for _, d := range []*dos.DOS{originalDOS, activeDOS} {
			d.Mouse = dos.Mouse{X: c.mx, Y: c.my, MaxX: 639, MaxY: 399, Press: c.presses, Calls: map[uint16]int{}}
			for i := 0; i < 3; i++ {
				d.Mouse.PressAt[i] = [2]uint16{c.mx, c.my}
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
		r := registers{AX: c.ax, BX: c.bx, CX: c.cx, DX: c.dx, SI: c.si, DI: c.di, BP: c.bp, SP: 0x8000, DS: 0x3500, ES: 0xa0c8, SS: 0x4000, CS: cs, IP: uint16(c.target), Flags: flags}
		stopCS := cs
		if c.target >= 0x20000 {
			stopCS = mouseCS
			r.CS = mouseCS
			r.DS = 0xa0c8
		}
		farRoot := c.target == 0x20000 || c.target == 0x20101
		if farRoot {
			stopCS = cs
		}
		original.CPU.R = [8]uint16{r.AX, r.CX, r.DX, r.BX, r.SP, r.BP, r.SI, r.DI}
		original.CPU.Seg = [4]uint16{r.ES, r.CS, r.SS, r.DS}
		original.CPU.IP = r.IP
		original.CPU.SetFlags(r.Flags)
		putWord(original.Mem, r.SS, r.SP, 0xf000)
		if farRoot {
			putWord(original.Mem, r.SS, r.SP+2, cs)
		}
		putWord(cMemory, r.SS, r.SP, 0xf000)
		if farRoot {
			putWord(cMemory, r.SS, r.SP+2, cs)
		}
		assignC(&m, originalRegs(original))
		fixture.core_cs = C.uint16_t(cs)
		fixture.mouse_cs = C.uint16_t(mouseCS)
		fixture.trace.calls = 0
		fixture.unsupported = 0
		var originalTrace []byte
		steps := 0
		for ; steps < 3000000; steps++ {
			cur := originalRegs(original)
			if cur.CS == stopCS && cur.IP == 0xf000 {
				break
			}
			loc := uint32(cur.IP) + 0x10000
			if cur.CS == mouseCS {
				loc = uint32(cur.IP) + 0x20000
			}
			if cur.CS == cs && cur.IP == 0x2259 {
				originalTicks++
				value := byte(0xff)
				if originalTicks <= c.timerHold {
					value = 0xfe
				}
				original.Mem[int(cs)*16+0xd2d] = value
			}
			_, known := functions[loc]
			checkpoint := cur.CS == cs && cur.IP == 0x2259
			if (cur.CS == cs || cur.CS == mouseCS) && known || checkpoint || cur.CS == machine.StubSeg && (cur.IP == 0x410 || cur.IP == 0x414) {
				name := entryName(loc)
				if checkpoint {
					name = "checkpoint_12259"
				}
				if cur.CS == machine.StubSeg {
					name = fmt.Sprintf("font@0080:%04X", cur.IP)
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
		C.input_run(&m, C.uint32_t(c.target), &fixture)
		if fixture.trace.calls > 8192 || fixture.unsupported != 0 {
			panic(fmt.Sprintf("unsupported=%x calls=%d input=%+v original=%+v C=%+v", fixture.unsupported, fixture.trace.calls, r, originalRegs(original), cRegs(&m)))
		}
		var cTrace []byte
		for i := 0; i < int(fixture.trace.calls); i++ {
			s := fixture.trace.trace[i]
			loc := uint32(s.target) + 0x10000
			if uint16(s.regs[11]) == mouseCS {
				loc = uint32(s.target) + 0x20000
			}
			_, known := functions[loc]
			if !known && !(uint16(s.regs[11]) == cs && s.target == 0x2259) && !(s.regs[11] == machine.StubSeg && (s.target == 0x410 || s.target == 0x414)) {
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
		if want != got || !bytes.Equal(original.Mem, cMemory) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(originalState, cState) || !bytes.Equal(originalPorts, cPorts) || !bytes.Equal(originalTrace, cTrace) || !bytes.Equal(originalAPI, cAPI) || !bytes.Equal(originalPortReads, cPortReads) || !bytes.Equal(mouseBytes(originalDOS), mouseBytes(activeDOS)) || originalTicks != activeTicks {
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": r, "original": want, "c": got, "original_api": hex.EncodeToString(originalAPI), "original_mouse": hex.EncodeToString(mouseBytes(originalDOS)), "c_mouse": hex.EncodeToString(mouseBytes(activeDOS)), "original_in": hex.EncodeToString(originalPortReads), "c_in": hex.EncodeToString(cPortReads), "original_ticks": originalTicks, "c_ticks": activeTicks, "c_api": hex.EncodeToString(cAPI), "original_trace": hex.EncodeToString(originalTrace), "c_trace": hex.EncodeToString(cTrace), "original_ports": hex.EncodeToString(originalPorts), "c_ports": hex.EncodeToString(cPorts), "original_device": hex.EncodeToString(originalState), "c_device": hex.EncodeToString(cState), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(active.VGA.Raw())), "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(cMemory))}
			return false
		}
		if c.group == "roundtrip" && c.target == 0x20000 && c.ax == 2 {
			if !bytes.Equal(original.VGA.Raw(), roundtripPixels) {
				panic("original cursor roundtrip did not restore background")
			}
			roundtripRestored++
		}
		audits++
		mousePollCounts[c.group] += originalQueries
		waitCheckpoints[c.group] += originalTicks
		{
			if !bytes.Equal(indexed(original), indexed(active)) {
				panic("content pixels")
			}
			pixelAudits++
		}
		for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), originalState, originalPorts, originalTrace, originalAPI, originalPortReads, mouseBytes(originalDOS)} {
			originalDigest.Write(b)
		}
		for _, b := range [][]byte{regBytes(got), cMemory, active.VGA.Raw(), cState, cPorts, cTrace, cAPI, cPortReads, mouseBytes(activeDOS)} {
			cDigest.Write(b)
		}
		return true
	}
	all := []inputCase{}
	add := func(c inputCase) { all = append(all, c) }
	base := func(g string, t uint32) inputCase {
		return inputCase{group: g, target: t, mask: 15, seed: 3, bp: 0x7000, di: 0x9000, ax: 0x9001, bx: 100*80 + 20, cx: 200, dx: 100, si: 0x31b, patch: 0x0feb, oldX: 120, oldY: 80, mx: 200, my: 100}
	}
	for service := 0; service < 16; service++ {
		for _, shown := range []uint8{0, 1, 2} {
			c := base("services", 0x20000)
			c.ax = uint16(service)
			c.cursor = shown
			c.bx = 1
			add(c)
		}
	}
	for _, t := range []uint32{0x2002e, 0x2006f, 0x20070, 0x2009a, 0x200bd, 0x200c0, 0x20137, 0x2016b, 0x2019d, 0x201c6, 0x201e4, 0x2020c, 0x20249, 0x202a0, 0x202bd, 0x202fe, 0x20101} {
		for _, shown := range []uint8{0, 1, 2} {
			c := base("cursor", t)
			c.cursor = shown
			if t == 0x200bd {
				c.ax = 3
			}
			if t == 0x20137 {
				c.cx = 0xffff
			}
			add(c)
		}
	}
	for _, xy := range [][2]uint16{{0, 0}, {1, 1}, {7, 7}, {8, 8}, {616, 383}, {624, 384}, {631, 399}, {632, 399}, {639, 399}} {
		for _, t := range []uint32{0x201e4, 0x2020c, 0x20249, 0x202bd, 0x20101} {
			c := base("edges", t)
			c.cursor = 1
			c.cx = xy[0]
			c.dx = xy[1]
			c.di = xy[0]
			c.bx = xy[1]*80 + xy[0]/8
			c.cx = (c.cx & 0xff00) | (xy[0] & 7)
			if t == 0x20101 {
				c.cx = xy[0]
			}
			add(c)
		}
	}
	for align := 0; align < 8; align++ {
		c := base("alignment", 0x201e4)
		c.cx = uint16(align)
		c.di = 200 + uint16(align)
		c.bx = 100*80 + 25
		add(c)
	}
	for _, patch := range []uint16{0x0feb, 0x9090} {
		for _, button := range []int{-1, 0, 1} {
			for _, delay := range []int{0, 3, 7} {
				c := base("wait", 0x1222b)
				c.patch = patch
				c.clickButton = button
				if button >= 0 {
					c.clickAt = delay + 1
				}
				c.timerHold = delay
				add(c)
			}
		}
	}
	for _, button := range []int{0, 1} {
		for delay := 1; delay <= 7; delay++ {
			c := base("poll", 0x121e7)
			c.clickButton = button
			c.clickAt = delay
			add(c)
		}
	}
	for _, slot := range []uint16{0, 62, 100, 407, 1021} {
		for _, button := range []int{-1, 0, 1} {
			c := base("message", 0x18810)
			c.ax = 0
			c.cx = slot
			c.bx = 8
			c.dx = 10
			if slot >= 0x196 {
				c.cx = 0x196 + (slot-0x196)/8
				c.ax = ((slot - 0x196) % 8) << 8
			}
			c.clickButton = button
			if button >= 0 {
				c.clickAt = 3
			}
			c.parameters = []uint16{0xff00, 0xff0c, 0xff03, 777, 12345}
			add(c)
		}
	}
	for _, patch := range []uint16{0x0feb, 0x9090} {
		c := base("patched-caller", 0x12216)
		c.patch = patch
		c.clickAt = 3
		c.clickButton = 0
		add(c)
	}
	for _, xy := range [][2]uint16{{200, 100}, {631, 383}, {639, 399}} {
		c := base("roundtrip", 0x20000)
		c.ax = 1
		c.mx = 120
		c.my = 80
		add(c)
		c = base("roundtrip", 0x20101)
		c.keep = true
		c.cx = xy[0]
		c.dx = xy[1]
		add(c)
		c = base("roundtrip", 0x20000)
		c.keep = true
		c.ax = 2
		add(c)
	}
	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{"schema": "wolong-c-input-parity-v1", "input_sha256": inputHash, "icon_sha256": assetHash, "routine_sha256": routineHashes, "cases": cases, "groups": groups, "full_ram_plane_audits": audits, "indexed_content_audits": pixelAudits, "content_top": 40, "content_height": 400, "entries_seen": entries, "cursor_roundtrips_restored": roundtripRestored, "mouse_queries": mousePollCounts, "counter_checkpoints": waitCheckpoints, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "trace_entry_size": 90, "platform": "Independent pinned DOS/font/VGA services; native C TALK/number/portrait control; IF/TF=0; CPU-model undefined flags; no original INT50 UI or wall-clock claim", "asset_banks": "TALK at 3500; SINARIO at 5000; portraits at 6000; UI map at 6800; ICONGRF at 2200/3000", "additional_assets": identities, "original_font_calls": originalDOS.Font.Calls, "c_font_calls": activeDOS.Font.Calls, "original_missing_fonts": originalDOS.Font.Missing, "c_missing_fonts": activeDOS.Font.Missing, "c_machine_code_match": false}
	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if e := os.WriteFile(*output, b, 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Input original/C %d; complete RAM/plane %d; content %d; pass=%v\n", cases, audits, pixelAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
