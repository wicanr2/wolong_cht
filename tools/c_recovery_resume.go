//go:build matching_resume

// Original scene entry, camera tail transfer and minimap redraw; independent platform APIs.
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
#include "/repo/tools/c_recovery/choice.c"
#include "/repo/tools/c_recovery/resource.c"
#include "/repo/tools/c_recovery/mapcells.c"
#include "/repo/tools/c_recovery/overlay.c"
#include "/repo/tools/c_recovery/resume.c"
#include "/repo/tools/c_recovery/resume_fixture.h"
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
var currentInput resumeCase
var originalPositions, activePositions int

type inputReadBus struct {
	cpu.Bus
	reads *[]byte
}

func (b inputReadBus) In8(port uint16) uint8 {
	v := b.Bus.In8(port)
	*b.reads = append(*b.reads, byte(port), byte(port>>8), v)
	return v
}

var functions = map[uint32]int{0x11d8e: 137, 0x1ec82: 94, 0x1ece0: 28, 0x15358: 110, 0x15609: 34, 0x1563b: 40, 0x154fc: 54, 0x155ec: 13, 0x15828: 55, 0x153c6: 144, 0x15456: 57, 0x1548f: 109, 0x15538: 15, 0x15547: 95, 0x15695: 128, 0x155a6: 70, 0x12fbf: 79, 0x157fe: 42, 0x15715: 122, 0x1578f: 111, 0x130cb: 8, 0x13119: 31, 0x122db: 163, 0x12286: 85, 0x1237e: 129, 0x1585f: 58, 0x15899: 167, 0x15940: 80, 0x12ad2: 34, 0x15990: 22, 0x1301c: 50, 0x12bd9: 121, 0x12c52: 141, 0x12cdf: 91, 0x12d3a: 30, 0x12fb1: 14, 0x12d58: 96, 0x12db8: 59, 0x12df3: 64, 0x130f0: 26, 0x1310a: 15, 0x13091: 58, 0x12e33: 86, 0x12e89: 114, 0x12efb: 118, 0x12f71: 64, 0x13e11: 84, 0x13e65: 41, 0x13e8e: 111, 0x15673: 34, 0x131ae: 68, 0x13496: 16, 0x13507: 19, 0x13dc9: 72, 0x1320c: 20, 0x13220: 66, 0x13262: 71, 0x132a9: 64, 0x132e9: 62, 0x13327: 97, 0x13388: 98, 0x133ea: 19, 0x133fd: 136, 0x13485: 17, 0x134a6: 11, 0x134b1: 86, 0x1351a: 12, 0x13526: 133, 0x135ab: 66, 0x135ed: 76, 0x13639: 48, 0x13669: 46, 0x13697: 45, 0x136c4: 78, 0x13712: 95, 0x13771: 103, 0x137d8: 29, 0x137f5: 59, 0x13138: 55, 0x14502: 70, 0x16a3d: 94, 0x123ff: 57, 0x12438: 33, 0x150d7: 73, 0x12078: 94, 0x120d6: 123, 0x138c7: 31, 0x138e6: 28, 0x13902: 230, 0x139e8: 288, 0x13c3d: 92, 0x13b7e: 43, 0x13d09: 60, 0x13d45: 35, 0x13c99: 39, 0x13cdc: 45, 0x19321: 21, 0x187ff: 17, 0x13d68: 41, 0x11d46: 72, 0x12216: 21, 0x17c6e: 147, 0x17d0d: 58, 0x17d47: 24, 0x17d5f: 52, 0x17da5: 30, 0x17dc3: 26, 0x17ddd: 13, 0x17dea: 2, 0x17dec: 5, 0x17df1: 4, 0x101b4: 33, 0x101db: 51, 0x193e9: 32, 0x167cd: 25, 0x167e6: 32, 0x16806: 32, 0x16826: 32, 0x1e3c0: 23, 0x1e3d7: 68, 0x1e41b: 56, 0x1e453: 38, 0x1895d: 71, 0x189de: 18, 0x1d5d4: 65, 0x10c14: 76, 0x10c60: 23, 0x10c77: 53, 0x19796: 45, 0x197c3: 45, 0x1f9b0: 107, 0x1fa1b: 28, 0x1fa37: 107, 0x1faa2: 32, 0x1fac2: 79, 0x1fb11: 24, 0x1f888: 176, 0x1f938: 97, 0x1f999: 23, 0x1c7f4: 74, 0x1c673: 59, 0x1f020: 288, 0x1f140: 59, 0x1f17b: 39, 0x1f1a3: 203, 0x10aaa: 47, 0x10ad9: 109, 0x10cac: 23, 0x10cc3: 27, 0x1c61f: 52, 0x1c6bf: 55, 0x1c6ae: 17, 0x1c6f6: 86, 0x1c74c: 41, 0x1c775: 25, 0x1c78e: 27, 0x10bcd: 71, 0x1030f: 40, 0x10337: 4, 0x1e993: 20, 0x1e9a7: 26, 0x1e9c1: 76, 0x1fb29: 126, 0x1fba7: 84, 0x1f6dc: 68, 0x1f878: 16, 0x1f720: 62, 0x1f7a4: 212, 0x106f5: 4, 0x106fd: 4, 0x1062f: 107, 0x1069a: 68, 0x106de: 23, 0x11e17: 47, 0x106f9: 4, 0x1075b: 119, 0x107d2: 115, 0x1084a: 90, 0x18853: 48, 0x189a4: 58, 0x1e38c: 26, 0x1f4df: 73, 0x18810: 67, 0x121e7: 47, 0x20000: 14, 0x2002e: 65, 0x2006f: 1, 0x20070: 42, 0x2009a: 35, 0x200bd: 3, 0x200c0: 54, 0x20137: 52, 0x2016b: 50, 0x2019d: 41, 0x201c6: 30, 0x201e4: 40, 0x2020c: 61, 0x20249: 87, 0x202a0: 29, 0x202bd: 65, 0x202fe: 29, 0x1036f: 84, 0x103c3: 35, 0x103e6: 46, 0x10414: 161, 0x104b5: 74, 0x104ff: 78, 0x1054d: 78, 0x1059b: 82, 0x105ed: 1, 0x1061f: 16, 0x10b46: 105, 0x10baf: 30, 0x19479: 70, 0x194bf: 80, 0x187af: 29, 0x1e378: 20, 0x1f4a2: 59, 0x10241: 129, 0x102c2: 14, 0x100df: 118, 0x1d46a: 25, 0x1d483: 32, 0x1d4c7: 88, 0x1d615: 85, 0x1d66a: 257, 0x1d76b: 23, 0x1d782: 20, 0x1d796: 81, 0x1d7e7: 29, 0x1d804: 70, 0x11cc9: 7, 0x12533: 112, 0x12af4: 54, 0x12b2a: 18, 0x12b3c: 108, 0x15d19: 162, 0x11f30: 42, 0x11f5a: 37, 0x13b08: 82, 0x15c58: 76, 0x19541: 112, 0x195b1: 24, 0x196ed: 101, 0x19752: 68, 0x1d4a3: 36, 0x1d51f: 181, 0x19409: 81, 0x1945a: 31, 0x1222b: 91, 0x20101: 54, 0x10701: 90, 0x108b2: 41, 0x108db: 41, 0x10904: 53, 0x10939: 34, 0x1095b: 35, 0x1097e: 6, 0x10984: 43, 0x1f75e: 67}

func entryName(loc uint32) string {
	if loc == 0x105ed {
		return "nullsub_5"
	}
	if loc == 0x2006f {
		return "nullsub_4"
	}
	for _, v := range []uint32{0x1d51f, 0x19409, 0x1945a, 0x1222b, 0x20101, 0x108b2, 0x108db, 0x10904, 0x10939, 0x1095b, 0x1097e, 0x10984, 0x10701, 0x1f75e} {
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

func soundBytes(d *dos.DOS) []byte {
	b, err := json.Marshal(d.Sound)
	if err != nil {
		panic(err)
	}
	return b
}

type resumeCase struct {
	group                                                              string
	target                                                             uint32
	ax, bx, cx, dx, si, di, bp                                         uint16
	scenario                                                           int
	message                                                            []byte
	parameters                                                         []uint16
	keep, df                                                           bool
	patch                                                              uint16
	cursor                                                             uint8
	oldX, oldY, mx, my                                                 uint16
	presses                                                            [3]uint16
	clickAt, clickButton, timerHold                                    int
	mask, readMap, seed                                                byte
	rows, visible, selected, visibleSelected                           uint8
	width                                                              uint8
	mode                                                               uint8
	moves                                                              []uint16
	acceptButton                                                       int
	liveX, liveY                                                       uint8
	expectedRow                                                        int
	liveString                                                         uint16
	filename                                                           string
	syntheticSize                                                      int
	cue, mute, oldCue                                                  uint8
	cameraX, cameraY, previousX, previousY, worldCursorX, worldCursorY uint16
	panel, hidden                                                      uint8
	loadBytes                                                          int
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
	originalDOS.Scratch = "/tmp/resource-fixtures"
	activeDOS := dos.New(active, "/orig")
	activeDOS.Scratch = "/tmp/resource-fixtures"
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
	mapdata, err := os.ReadFile("/output/fixtures/MMAP.raw")
	if err != nil || len(mapdata) != 98304 {
		panic("map fixture")
	}
	mdl, err := os.ReadFile("/orig/MMAP.MDL")
	if err != nil {
		panic(err)
	}
	mch, err := os.ReadFile("/orig/MMAP.MCH")
	if err != nil {
		panic(err)
	}
	mouseCS := binary.LittleEndian.Uint16(original.Mem[int(cs)*16+0x2082:])
	var originalAPI, cAPI, originalPortReads []byte
	originalQueries, activeQueries, originalTicks := 0, 0, 0
	original.CPU.Bus = inputReadBus{Bus: original.CPU.Bus, reads: &originalPortReads}
	installTrace := func(device *machine.Machine, d *dos.DOS, dst *[]byte, queries *int) {
		previous := device.CPU.IntHook
		device.CPU.IntHook = func(c *cpu.CPU, n uint8) bool {
			if n == 0x33 && c.R[cpu.AX] == 3 {
				positions := &activePositions
				if device == original {
					positions = &originalPositions
				}
				(*positions)++
				poll := *positions - 2
				if poll >= 0 && currentInput.acceptButton >= 0 {
					if poll < len(currentInput.moves) {
						d.Mouse.X = 0
						d.Mouse.Y = currentInput.moves[poll]
					} else {
						d.Mouse.PressButton(currentInput.acceptButton)
					}
				}
			}
			if n == 0x33 && c.R[cpu.AX] == 5 {
				(*queries)++
				if currentInput.clickAt > 0 && *queries == currentInput.clickAt {
					d.Mouse.PressButton(currentInput.clickButton)
				}
			}
			if n == 0x21 || n == 0x33 || n == 0x61 {
				*dst = append(*dst, n)
				*dst = append(*dst, regBytes(originalRegs(device))...)
			}
			result := previous(c, n)
			if n == 0x21 || n == 0x33 || n == 0x61 {
				*dst = append(*dst, regBytes(originalRegs(device))...)
			}
			return result
		}
	}
	installTrace(original, originalDOS, &originalAPI, &originalQueries)
	installTrace(active, activeDOS, &cAPI, &activeQueries)
	if err := os.MkdirAll("/tmp/resource-fixtures", 0755); err != nil {
		panic(err)
	}

	initialRAM := append([]byte(nil), original.Mem...)
	copy(initialRAM[0x22000:], icons[0x9700:])
	copy(initialRAM[0x45000:], talk)
	copy(initialRAM[0x30000:], mdl)
	copy(initialRAM[0x38000:], mch)
	copy(initialRAM[0x50000:], mapdata)
	copy(initialRAM[0x68000:], icons[0x6700:0x9700])
	for _, v := range [][2]uint16{{0xd48, 0x2200}, {0xd4a, 0x226c}, {0xd4c, 0x228f}, {0xd54, 0x2284}, {0xd3a, 0x6800}, {0xd3c, 0x7800},
		{0xd38, 0x4500}, {0xd40, 0x8000}, {0xd52, 0x7000}, {0xd84a, 0x3000}, {0xd84c, 0x3800}, {0xd84e, 0x2500}, {0xd850, 0x5000},
		{0x987c, 0x9000}, {0x9876, 0x3000}, {0x987a, 0x4200}, {0xe479, 0x7600}, {0xe47b, 0}, {0xd30a, 0x7000}, {0xd30e, 0x5200}} {
		putWord(initialRAM, cs, v[0], v[1])
	}
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
	fixture := C.KiResumeFixture{}
	groups, entries := map[string]int{}, map[string]int{}
	mousePollCounts, waitCheckpoints := map[string]int{}, map[string]int{}
	sampleCounts := map[uint32]int{}
	cases, audits, pixelAudits := 0, 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	var roundtripPixels []byte
	var xorPixels []byte
	xorRestored := 0
	roundtripRestored := 0
	selectedRows := map[string][]int{}
	bufferAudits := 0
	check := func(c resumeCase) bool {
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

			}
		}
		if c.group == "roundtrip" && !c.keep {
			roundtripPixels = append([]byte{}, original.VGA.Raw()...)
		}
		if c.group == "xor-roundtrip" && !c.keep {
			xorPixels = append([]byte{}, original.VGA.Raw()...)
		}
		if c.syntheticSize >= 0 {
			b := make([]byte, c.syntheticSize)
			for i := range b {
				b[i] = byte(i*37 + i/257 + 19)
			}
			if err := os.WriteFile("/tmp/resource-fixtures/BOUND.DAT", b, 0644); err != nil {
				panic(err)
			}
		}
		for _, mem := range [][]byte{original.Mem, cMemory} {
			if !c.keep {
				copy(mem[0x70000:], worlds[c.scenario*22208+0x80:c.scenario*22208+0x52c0])
				copy(mem[int(cs)*16+0xcf0:], worlds[c.scenario*22208:c.scenario*22208+59])
				mem[int(cs)*16+0x845] = 0
				for i := 0; i < 4; i++ {
					mem[int(cs)*16+0x846+i] = 0xff
				}
			}

			putWord(mem, cs, 0xcfd, 0)
			if !c.keep {
				for i := 0; i < 920; i++ {
					at := 0x25000 + i*8
					mem[at] = 0x40
					mem[at+1] = 0
					mem[at+2] = 255
					for j := 3; j < 8; j++ {
						mem[at+j] = 255
					}
				}
				for _, v := range [][2]uint16{{0x988e, c.cameraX}, {0x9890, c.cameraY}, {0x9892, c.previousX}, {0x9894, c.previousY}, {0x9896, c.worldCursorX}, {0x9898, c.worldCursorY}, {0x9886, c.cx}, {0x9888, c.dx}, {0x988a, 0xffff}, {0x988c, 0xffff}, {0xd854, c.cameraX}, {0xd856, c.cameraY}} {
					putWord(mem, cs, v[0], v[1])
				}
				mem[int(cs)*16+0x98a6] = c.panel
				mem[int(cs)*16+0x98a2] = c.hidden
				putWord(mem, cs, 0xe479, 0x7600)
				putWord(mem, 0x7600, 0, 0)
			}

			for i, v := range c.parameters {
				putWord(mem, 0x4000, uint16(0x9000+i*2), v)
			}
			if c.filename != "" {
				copy(mem[int(cs)*16+0x7800:], append([]byte(c.filename), 0))
			}
			mem[int(cs)*16+0x20e] = c.oldCue
			mem[int(cs)*16+0xcf9] = c.mute
			putWord(mem, cs, 0x2239, c.patch)
			if !c.keep {
				putWord(mem, mouseCS, 0xfc, c.oldX)
				putWord(mem, mouseCS, 0xfe, c.oldY)
				mem[int(mouseCS)*16+0x100] = c.cursor
			}
		}
		if !c.keep {
			if !bytes.Equal(original.Mem[0x30000:0x38000], mdl) || !bytes.Equal(original.Mem[0x38000:0x38000+len(mch)], mch) || !bytes.Equal(original.Mem[0x50000:0x68000], mapdata) || !bytes.Equal(original.Mem[0x68000:0x6b000], icons[0x6700:0x9700]) {
				panic("prepared asset banks")
			}
		}
		originalPositions = 0
		activePositions = 0
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
		var expectedBuffer, expectedTail []byte
		bufferAt := int(c.bx)*16 + int(c.si)
		if c.group == "stream" || c.group == "assets" {
			path := "/orig/" + c.filename
			if c.syntheticSize >= 0 {
				path = "/tmp/resource-fixtures/" + c.filename
			}
			expectedBuffer, err = os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			expectedTail = append([]byte{}, original.Mem[bufferAt+len(expectedBuffer):bufferAt+len(expectedBuffer)+64]...)
		}
		flags := uint16(2 | uint16(c.seed&1))
		if c.df {
			flags |= 0x400
		}
		r := registers{AX: c.ax, BX: c.bx, CX: c.cx, DX: c.dx, SI: c.si, DI: c.di, BP: c.bp, SP: 0x8000, DS: 0x3500, ES: 0xa0c8, SS: 0xf000, CS: cs, IP: uint16(c.target), Flags: flags}
		if true {
			r.DS = cs
		}
		if c.target == 0x13b08 {
			r.DS = 0x7000
		}
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
		for ; steps < 6000000; steps++ {
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
		if steps == 6000000 {
			_ = os.WriteFile("/output/results/original-failure.ram", original.Mem, 0644)
			_ = os.WriteFile("/output/results/original-failure.trace", originalTrace, 0644)

			panic(fmt.Sprintf("bounded original routine input=%+v registers=%+v positions=%d queries=%d mouse=%s", c, originalRegs(original), originalPositions, originalQueries, mouseBytes(originalDOS)))
		}
		C.resume_run(&m, C.uint16_t(c.target), &fixture)
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
		if want != got || !bytes.Equal(original.Mem, cMemory) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(originalState, cState) || !bytes.Equal(originalPorts, cPorts) || !bytes.Equal(originalTrace, cTrace) || !bytes.Equal(originalAPI, cAPI) || !bytes.Equal(originalPortReads, cPortReads) || !bytes.Equal(mouseBytes(originalDOS), mouseBytes(activeDOS)) || !bytes.Equal(soundBytes(originalDOS), soundBytes(activeDOS)) || originalTicks != activeTicks {
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": r, "original": want, "c": got, "original_api": hex.EncodeToString(originalAPI), "original_sound": hex.EncodeToString(soundBytes(originalDOS)), "c_sound": hex.EncodeToString(soundBytes(activeDOS)), "original_mouse": hex.EncodeToString(mouseBytes(originalDOS)), "c_mouse": hex.EncodeToString(mouseBytes(activeDOS)), "original_in": hex.EncodeToString(originalPortReads), "c_in": hex.EncodeToString(cPortReads), "original_ticks": originalTicks, "c_ticks": activeTicks, "c_api": hex.EncodeToString(cAPI), "original_trace": hex.EncodeToString(originalTrace), "c_trace": hex.EncodeToString(cTrace), "original_ports": hex.EncodeToString(originalPorts), "c_ports": hex.EncodeToString(cPorts), "original_device": hex.EncodeToString(originalState), "c_device": hex.EncodeToString(cState), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(active.VGA.Raw())), "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(cMemory))}
			return false
		}
		if c.group == "roundtrip" && c.target == 0x20000 && c.ax == 2 {
			if !bytes.Equal(original.VGA.Raw(), roundtripPixels) {
				panic("original cursor roundtrip did not restore background")
			}
			roundtripRestored++
		}
		if c.expectedRow >= 0 && int(byte(want.AX)) != c.expectedRow {
			panic(fmt.Sprintf("original expected row %d, got %d", c.expectedRow, byte(want.AX)))
		}
		if c.expectedRow >= 0 {
			selectedRows[c.group] = append(selectedRows[c.group], int(byte(want.AX)))
		}
		if c.group == "xor-roundtrip" && c.keep {
			if !bytes.Equal(original.VGA.Raw(), xorPixels) {
				panic("original XOR band did not restore background")
			}
			xorRestored++
		}
		audits++
		if expectedBuffer != nil {
			if !bytes.Equal(original.Mem[bufferAt:bufferAt+len(expectedBuffer)], expectedBuffer) || !bytes.Equal(original.Mem[bufferAt+len(expectedBuffer):bufferAt+len(expectedBuffer)+64], expectedTail) {
				panic("original complete buffer or tail guard differs from source file")
			}
			bufferAudits++
		}
		mousePollCounts[c.group] += originalQueries
		waitCheckpoints[c.group] += originalTicks
		{
			if !bytes.Equal(indexed(original), indexed(active)) {
				panic("content pixels")
			}
			pixelAudits++
		}
		for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), originalState, originalPorts, originalTrace, originalAPI, originalPortReads, mouseBytes(originalDOS), soundBytes(originalDOS)} {
			originalDigest.Write(b)
		}
		for _, b := range [][]byte{regBytes(got), cMemory, active.VGA.Raw(), cState, cPorts, cTrace, cAPI, cPortReads, mouseBytes(activeDOS), soundBytes(activeDOS)} {
			cDigest.Write(b)
		}
		return true
	}

	all := []resumeCase{}
	add := func(c resumeCase) { all = append(all, c) }
	base := func(g string, target uint32) resumeCase {
		return resumeCase{group: g, target: target, mask: 15, seed: 3, bp: 0x7000, di: 0x9000, ax: 0, bx: 176, cx: 0, dx: 448, si: 0, patch: 0x0feb, oldX: 120, oldY: 80, mx: 200, my: 100, expectedRow: -1, acceptButton: -1, syntheticSize: -1, oldCue: 0xff, cameraX: 170, cameraY: 98, previousX: 170, previousY: 98, worldCursorX: 10, worldCursorY: 12}
	}
	for _, target := range []uint32{0x19541, 0x196ed} {
		for bit := uint16(0); bit < 8; bit++ {
			for _, y := range []uint16{40, 89, 140} {
				for _, df := range []bool{false, true} {
					c := base("minimap", target)
					c.dx = 440 + bit
					c.bx = y
					c.df = df
					add(c)
				}
			}
		}
	}
	for _, x := range []uint16{607, 608, 65535} {
		c := base("minimap-clamp", 0x19541)
		c.dx = x
		c.bx = 140
		add(c)
	}
	for _, y := range []uint16{0, 1, 2, 12, 24} {
		for _, hidden := range []byte{0, 1} {
			c := base("camera", 0x11f30)
			c.worldCursorY = y
			c.hidden = hidden
			add(c)
		}
	}
	for _, target := range []uint32{0x11f5a, 0x1d4a3} {
		for _, x := range []uint16{0, 10, 39} {
			for _, y := range []uint16{0, 12, 22} {
				c := base("camera-helper", target)
				c.dx = x
				c.bx = y
				add(c)
			}
		}
	}
	for _, xy := range [][2]uint16{{170, 98}, {171, 98}, {170, 99}, {0xffff, 0}} {
		c := base("panel", 0x15c58)
		c.previousX = xy[0]
		c.previousY = xy[1]
		c.panel = 4
		add(c)
	}
	for scenario := 0; scenario < 4; scenario++ {
		for _, panel := range []byte{0, 4} {
			for _, xy := range [][2]uint16{{0, 0}, {170, 98}, {344, 233}} {
				c := base("resume", 0x11d46)
				c.scenario = scenario
				c.panel = panel
				c.cameraX = xy[0]
				c.cameraY = xy[1]
				c.previousX = 0xffff
				c.worldCursorY = 12
				c.cx = 160
				c.dx = 160
				add(c)
			}
		}
	}
	for scenario := 0; scenario < 4; scenario++ {
		for _, panel := range []byte{0, 4} {
			for _, item := range [][2]uint16{{0, 0x182}, {1, 0x18c}} {
				c := base("scene", 0x13b08)
				c.scenario = scenario
				c.panel = panel
				c.ax = item[0]
				c.cx = item[1]
				c.si = 0x4240
				c.di = 0x4240
				c.dx = 0x9000
				c.previousX = 0xffff
				c.clickAt = 1
				c.clickButton = 0
				add(c)
			}
		}
	}
	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{"schema": "wolong-c-resume-parity-v1", "complete_buffer_audits": bufferAudits, "input_sha256": inputHash, "icon_sha256": assetHash, "routine_sha256": routineHashes, "cases": cases, "groups": groups, "full_ram_plane_audits": audits, "indexed_content_audits": pixelAudits, "content_top": 40, "content_height": 400, "entries_seen": entries, "xor_roundtrips_restored": xorRestored, "selected_rows": selectedRows, "cursor_roundtrips_restored": roundtripRestored, "mouse_queries": mousePollCounts, "counter_checkpoints": waitCheckpoints, "passed": mismatch == nil, "mismatch": mismatch, "original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)), "trace_entry_size": 90, "platform": "Independent pinned DOS/font/VGA services; native C full-file/cue/allocation control; IF/TF=0; CPU-model undefined flags; no original INT50 UI or wall-clock claim", "asset_banks": "MDL/MCH 3000/3800; world map 5000; world records 7000; ICONGRF 2200/6800; TALK 4500; hotspot 7600; BGM 7800; SS F000:8000", "additional_assets": identities, "original_font_calls": originalDOS.Font.Calls, "c_font_calls": activeDOS.Font.Calls, "original_missing_fonts": originalDOS.Font.Missing, "c_missing_fonts": activeDOS.Font.Missing, "c_machine_code_match": false}
	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if e := os.WriteFile(*output, b, 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Scene resume original/C %d; complete RAM/plane %d; content %d; pass=%v\n", cases, audits, pixelAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
