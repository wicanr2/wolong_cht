//go:build matching_startup

// Original startup/main-loop/naming graph; spec253. Independent raw fixtures and original-first audits.
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/rng.c"
#include "/repo/tools/c_recovery/seed.c"
#include "/repo/tools/c_recovery/clock.c"
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
#include "/repo/tools/c_recovery/verdict.c"
#include "/repo/tools/c_recovery/list.c"
#include "/repo/tools/c_recovery/catalog.c"
#include "/repo/tools/c_recovery/formation.c"
#include "/repo/tools/c_recovery/personnel.c"
#include "/repo/tools/c_recovery/details.c"
#include "/repo/tools/c_recovery/march.c"
#include "/repo/tools/c_recovery/strategy.c"
#include "/repo/tools/c_recovery/main.c"
#include "/repo/tools/c_recovery/main_fixture.h"
#include "/repo/tools/c_recovery/bootstrap.c"
#include "/repo/tools/c_recovery/interaction.c"
#include "/repo/tools/c_recovery/tick.c"
#include "/repo/tools/c_recovery/route.c"
#include "/repo/tools/c_recovery/outcome.c"
#include "/repo/tools/c_recovery/engagement.c"
#include "/repo/tools/c_recovery/engagement_fixture.h"
#include "/repo/tools/c_recovery/army.c"
#include "/repo/tools/c_recovery/army_fixture.h"
#include "/repo/tools/c_recovery/startup.c"
#include "/repo/tools/c_recovery/startup_fixture.h"
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
	"sort"
	"strings"
	"unsafe"
)

var startupCoreCS, startupMouseCS uint16
var startupCTrace []byte
var startupCCalls int
var activeDevice *machine.Machine
var cPortReads []byte
var activeTicks int
var currentInput startupCase
var originalPositions, activePositions int
var originalPolls, activePolls int
var originalDOSDevice, activeDOSDevice *dos.DOS

var curSide, curGroup string
var curTarget uint32
var curScenario, curCase int

var snapCase1RAM []byte

var snapB5RAM []byte
var pcRing [96]uint32
var pcRingN uint64
var evtRing [32]string
var evtRingN uint64
var nativeHist = map[uint32]int{}

func evtEmit(s string) {
	evtRing[evtRingN%uint64(len(evtRing))] = s
	evtRingN++
}

func ringPCs(n int) string {
	var out []string
	start := uint64(0)
	if pcRingN > uint64(len(pcRing)) {
		start = pcRingN - uint64(len(pcRing))
	}
	var all []uint32
	for i := start; i < pcRingN; i++ {
		all = append(all, pcRing[i%uint64(len(pcRing))])
	}
	if len(all) > n {
		all = all[len(all)-n:]
	}
	for _, loc := range all {
		out = append(out, fmt.Sprintf("%X/%s", loc, entryName(loc)))
	}
	return strings.Join(out, " ")
}

func ringEvts(n int) string {
	start := uint64(0)
	if evtRingN > uint64(len(evtRing)) {
		start = evtRingN - uint64(len(evtRing))
	}
	var all []string
	for i := start; i < evtRingN; i++ {
		all = append(all, evtRing[i%uint64(len(evtRing))])
	}
	if len(all) > n {
		all = all[len(all)-n:]
	}
	return strings.Join(all, " | ")
}

func firstDiff(a, b []byte) (int, string, string) {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			e := i + 16
			if e > len(a) {
				e = len(a)
			}
			f := i + 16
			if f > len(b) {
				f = len(b)
			}
			return i, hex.EncodeToString(a[i:e]), hex.EncodeToString(b[i:f])
		}
	}
	if len(a) != len(b) {
		return n, fmt.Sprintf("len=%d", len(a)), fmt.Sprintf("len=%d", len(b))
	}
	return -1, "", ""
}

func parseCTraceIPs(t []byte) []registers {
	var out []registers
	for off := 0; off+90 <= len(t); off += 90 {
		out = append(out, registers{IP: uint16(t[off]) | uint16(t[off+1])<<8})
	}
	return out
}

type inputReadBus struct {
	cpu.Bus
	reads *[]byte
}

func (b inputReadBus) In8(port uint16) uint8 {
	v := b.Bus.In8(port)
	*b.reads = append(*b.reads, byte(port), byte(port>>8), v)
	return v
}

var functions = map[uint32]int{}

func init() {
	for target, size := range map[uint32]int{
		0x10000: 67, 0x10043: 24, 0x1005B: 15, 0x1006A: 1, 0x1006B: 116,
		0x10155: 95, 0x10210: 38, 0x10236: 11, 0x1033B: 19, 0x1034E: 8,
		0x10356: 23, 0x10612: 13, 0x11A6E: 64, 0x11AAE: 21, 0x11AC3: 151,
		0x11BE0: 173, 0x11CD0: 118, 0x18B12: 46, 0x18B40: 29, 0x18CAE: 81,
		0x18DC8: 146, 0x18E5A: 70, 0x18EA0: 151, 0x18F37: 54, 0x18F6D: 15,
		0x18F7C: 27, 0x18FC9: 149, 0x1905E: 16, 0x1906E: 17, 0x1907F: 14,
		0x1908D: 22, 0x190A3: 29, 0x190C0: 109, 0x1912D: 23, 0x19144: 23,
		0x1915B: 33, 0x1917C: 152, 0x19214: 15, 0x19223: 103, 0x1928A: 89,
		0x192E3: 11, 0x192EE: 7, 0x192F5: 20, 0x1EB6C: 39, 0x1EB93: 73,
	} {
		functions[target] = size
	}
}

func entryName(loc uint32) string {
	switch loc {
	case 0x10000:
		return "start"
	case 0x10043:
		return "raw_entry_10043"
	case 0x1006A:
		return "raw_entry_1006A"
	case 0x10356:
		return "raw_entry_10356"
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
func paletteBytes(m *machine.Machine) []byte {
	b := append([]byte(nil), m.DAC[:]...)
	palette := m.Palette()
	for _, color := range palette {
		b = append(b, color[:]...)
	}
	return b
}
func portMapBytes(m *machine.Machine) []byte {
	b, err := json.Marshal(struct {
		Out map[uint16]uint8
		In  map[uint16]uint64
	}{m.Ports, m.PortsIn})
	if err != nil {
		panic(err)
	}
	return b
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

func consoleBytes(d *dos.DOS) []byte { return append([]byte(nil), d.Console...) }

func exitBytes(d *dos.DOS) []byte {
	b := []byte{0, d.ExitCode}
	if d.Exited {
		b[0] = 1
	}
	return b
}

func periodicBytes(m *machine.Machine) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, m.PeriodicCalls())
	return b
}

func allocBytes(d *dos.DOS) []byte {
	type alloc struct {
		Fn   uint8
		Want uint16
		Seg  uint16
		OK   bool
	}
	type read struct {
		Name   string
		Handle uint16
		Seg    uint16
		Off    uint16
		Want   uint16
		Got    int
	}
	aa := make([]alloc, 0, len(d.Allocs))
	for _, op := range d.Allocs {
		aa = append(aa, alloc{op.Fn, op.Want, op.Seg, op.OK})
	}
	rr := make([]read, 0, len(d.Reads))
	for _, op := range d.Reads {
		rr = append(rr, read{op.Name, op.Handle, op.Seg, op.Off, op.Want, op.Got})
	}
	b, err := json.Marshal(map[string]any{"allocs": aa, "reads": rr})
	if err != nil {
		panic(err)
	}
	return b
}

type startupEvent struct {
	x, y   uint16
	button int
}

func startupCheckpoint(device *machine.Machine, d *dos.DOS, t uint16) {
	if t != 0x21e7 && t != 0x84dd && !(t == 0x222b && currentInput.pollFeed) {
		return
	}
	polls := &activePolls
	if device != activeDevice {
		polls = &originalPolls
	}
	event := startupEvent{100, 120, -1}
	if len(currentInput.events) > 0 {
		if *polls < len(currentInput.events) {
			event = currentInput.events[*polls]
		} else {
			event = currentInput.events[len(currentInput.events)-1]
		}
	}
	(*polls)++
	d.Mouse.ReleaseButton(0)
	d.Mouse.ReleaseButton(1)
	d.Mouse.X = event.x
	d.Mouse.Y = event.y
	if event.button >= 0 {
		d.Mouse.PressButton(event.button)
	}
	side := "orig"
	if device == activeDevice {
		side = "C"
	}
	evtEmit(fmt.Sprintf("%s wait t=%04X poll=%d ev=(%d,%d,%d)", side, t, *polls, event.x, event.y, event.button))
}

type startupCase struct {
	group                                     string
	target                                    uint32
	ax, bx, cx, dx, si, di, bp                uint16
	ds, es                                    uint16
	entrySS, entrySP, returnIP                uint16
	df                                        bool
	seed                                      byte
	scenario                                  int
	player                                    byte
	pspCmd, pspFill                           byte
	mouseAbsent                               bool
	prePress                                  bool
	pollFeed                                  bool
	events                                    []startupEvent
	pauseIters                                int
	transfer                                  bool
	transferIP                                uint16
	chain                                     bool
	expectExit                                bool
	d2c, d2d, d2a, a3, a5                   byte
	anchorSS, anchorSP                         uint16
	videoMode                                 int
	prealloc                                  bool
	namePat                                   int
	int61AH                                   []uint8
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
	originalDOSDevice = originalDOS
	activeDOSDevice = activeDOS
	if err = original.LoadEXE(raw); err != nil {
		panic(err)
	}
	if err = active.LoadEXE(raw); err != nil {
		panic(err)
	}

	for _, d := range []*machine.Machine{original, active} {
		// IRQ0 unreachable (1<<50) but NOT 0: int61 AH=0Ch computes
		// every = IRQ0Every/16 clamped to >= 1, so 0 machine-guns the
		// periodic callback every step. The stale nextIRQ0 (165000) is
		// burned off by the pre-roll below instead.
		d.IRQ0Every = 1 << 50
		d.KeyEvery = 0
		d.VGAFrameEvery = 0
	}
	originalDOS.Install()
	activeDOS.Install()
	fmt.Fprintf(os.Stderr, "ident pre-roll orig1100=%x\n", original.Mem[0x1100:0x1104])
	// Burn past the stale nextIRQ0 (165000, seeded at construction): run
	// both machines 200k steps on a scratch jmp-self with IF=1 so the one
	// stale IRQ0 delivers here, then reset BDA. After this nextIRQ0 and
	// periodic every (1<<46) are both unreachable for the whole matrix.
	for _, device := range []*machine.Machine{original, active} {
		device.Mem[0x90000] = 0xEB
		device.Mem[0x90001] = 0xFE
		savedR, savedSeg, savedIP, savedFlags := device.CPU.R, device.CPU.Seg, device.CPU.IP, device.CPU.Flags
		device.CPU.Seg[cpu.CS] = 0x9000
		device.CPU.Seg[cpu.SS] = 0x9000
		device.CPU.R[cpu.SP] = 0x100
		device.CPU.IP = 0
		device.CPU.SetFlags(0x202)
		for i := 0; i < 200000; i++ {
			if err := device.Step(); err != nil {
				panic(err)
			}
		}
		device.Mem[0x46C], device.Mem[0x46D], device.Mem[0x46E], device.Mem[0x46F] = 0, 0, 0, 0
		device.CPU.R, device.CPU.Seg, device.CPU.IP = savedR, savedSeg, savedIP
		device.CPU.SetFlags(savedFlags)
	}
	fmt.Fprintf(os.Stderr, "ident post-roll orig1100=%x active1100=%x\n", original.Mem[0x1100:0x1104], active.Mem[0x1100:0x1104])
	cs := original.CPU.Seg[cpu.CS]
	loaderSS, loaderSP := original.CPU.Seg[cpu.SS], original.CPU.R[cpu.SP]
	fmt.Fprintf(os.Stderr, "loader cs=%04X ss=%04X sp=%04X\n", original.CPU.Seg[cpu.CS], loaderSS, loaderSP)
	header := int(binary.LittleEndian.Uint16(raw[8:])) * 16
	probeRaw, err := os.ReadFile("/repo/workplace/matching-decompilation/c-startup/ida/ida-probe.json")
	if err != nil {
		panic(err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(probeRaw)) != "3d3761a281fa2ef1b50d1e8088e29053335ec7317473448e5624681e87cd8ea7" {
		panic("startup probe identity")
	}
	var probe struct {
		RecoveryTargets []uint32 `json:"recovery_targets"`
		Targets         []struct {
			IDALinear uint32 `json:"ida_linear"`
			Name      string `json:"name"`
			FileSHA   string `json:"file_sha256"`
			Chunks    []struct {
				Start, End uint32
				FileBytes  string `json:"file_bytes"`
			} `json:"chunks"`
		} `json:"targets"`
		Blocks []struct {
			IDALinear uint32 `json:"ida_linear"`
			End       uint32 `json:"end_ida_linear"`
			FileBytes string `json:"file_bytes"`
			Hash      string `json:"file_sha256"`
		} `json:"decoded_blocks"`
		Tables []struct {
			Call  uint32   `json:"call_ida_linear"`
			Words []uint16 `json:"words"`
		} `json:"fixed_dispatch_tables"`
		FarSites  []uint32 `json:"far_call_sites"`
		FarTarget uint32   `json:"far_call_target"`
	}
	if err = json.Unmarshal(probeRaw, &probe); err != nil {
		panic(err)
	}
	if len(probe.RecoveryTargets) != 42 || len(probe.Blocks) != 3 || len(probe.Tables) != 1 || len(probe.FarSites) != 15 || probe.FarTarget != 0x20000 {
		panic("startup source partition")
	}
	verifyChunk := func(start, end uint32, encoded string) []byte {
		file, e := hex.DecodeString(encoded)
		if e != nil {
			panic(e)
		}
		at := int(start-0x10000) + header
		if len(file) != int(end-start) || !bytes.Equal(file, raw[at:at+len(file)]) {
			panic("startup file bytes")
		}
		expected := append([]byte(nil), file...)
		count, table := int(binary.LittleEndian.Uint16(raw[6:])), int(binary.LittleEndian.Uint16(raw[24:]))
		for i := 0; i < count; i++ {
			off := header + int(binary.LittleEndian.Uint16(raw[table+i*4:])) + 16*int(binary.LittleEndian.Uint16(raw[table+i*4+2:]))
			if off >= at && off < at+len(file) {
				binary.LittleEndian.PutUint16(expected[off-at:], binary.LittleEndian.Uint16(raw[off:])+cs)
			}
		}
		physical := int(cs)*16 + int(start-0x10000)
		if !bytes.Equal(expected, original.Mem[physical:physical+len(expected)]) {
			off, a, b := firstDiff(expected, original.Mem[physical:physical+len(expected)])
			panic(fmt.Sprintf("startup runtime identity %X off=%d exp=%s got=%s", start, off, a, b))
		}
		return file
	}
	routineHashes := map[string]string{}
	selected := map[uint32]bool{}
	for _, t := range probe.RecoveryTargets {
		selected[t] = true
	}
	for _, t := range probe.Targets {
		if !selected[t.IDALinear] {
			continue
		}
		var file []byte
		for _, chunk := range t.Chunks {
			file = append(file, verifyChunk(chunk.Start, chunk.End, chunk.FileBytes)...)
		}
		hash := fmt.Sprintf("%x", sha256.Sum256(file))
		if hash != t.FileSHA {
			panic("startup primary routine hash")
		}
		routineHashes[t.Name] = hash
	}
	if len(routineHashes) != 42 {
		panic("startup primary must contain exactly 42 recovered named entries")
	}
	blockHashes := map[string]string{}
	for _, block := range probe.Blocks {
		file := verifyChunk(block.IDALinear, block.End, block.FileBytes)
		hash := fmt.Sprintf("%x", sha256.Sum256(file))
		if hash != block.Hash {
			panic("startup raw identity")
		}
		blockHashes[fmt.Sprintf("raw_entry_%X", block.IDALinear)] = hash
	}
	for _, table := range probe.Tables {
		for _, word := range table.Words {
			if !selected[0x10000+uint32(word)] {
				panic("startup dispatch word outside recovery")
			}
		}
	}
	manifestRaw, err := os.ReadFile("/output/results/c-source.sha256")
	if err != nil {
		panic(err)
	}
	manifestHash := fmt.Sprintf("%x", sha256.Sum256(manifestRaw))
	gamePalette, err := os.ReadFile("/orig/GAMEPAL.BRG")
	if err != nil || len(gamePalette) != 384 || fmt.Sprintf("%x", sha256.Sum256(gamePalette)) != "1f0119c75ea5cd333bd3ac75ef92030f93924f011728edc9b1f727c483263708" {
		panic("GAMEPAL identity")
	}
	kyo, err := os.ReadFile("/orig/KYOGRF.DAT")
	if err != nil || len(kyo) != 69120 || fmt.Sprintf("%x", sha256.Sum256(kyo)) != "e086f526bdada5baf751c41d2f73a78a0ba70002f282f63a9f33114542ed933f" {
		panic("KYOGRF identity")
	}
	_ = kyo
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
	nameFont, err := os.ReadFile("/orig/END_S15.DAT")
	if err != nil {
		panic(err)
	}
	identities := map[string]struct {
		Size int
		Hash string
	}{"TALK.DAT": {34182, "08a22e09791d0a6ec2968e87d8655e12c91b45e00fae460b28593b35ff85e384"}, "SINARIO.DAT": {88832, "21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c"}, "KAOGRF.DAT": {307200, "b9c7745e3ed9b32f0c12003fe81f82756136f738427dc421ec96f6c4ed5c4ac8"}, "END_S15.DAT": {5242, fmt.Sprintf("%x", sha256.Sum256(nameFont))}}
	for name, data := range map[string][]byte{"TALK.DAT": talk, "SINARIO.DAT": worlds, "KAOGRF.DAT": portraits} {
		id := identities[name]
		if len(data) != id.Size || fmt.Sprintf("%x", sha256.Sum256(data)) != id.Hash {
			panic(name + " identity")
		}
	}
	if len(nameFont) != 5242 {
		panic("END_S15 size")
	}
	mapdata, err := os.ReadFile("/output/fixtures/MMAP.raw")
	if err != nil || len(mapdata) != 98304 || fmt.Sprintf("%x", sha256.Sum256(mapdata)) != "740708c27a89db0a7f82865be623b5732099ec97aabf399e7dbde1450c87c861" {
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
	fmt.Fprintf(os.Stderr, "segs cs=%04X mouseCS=%04X\n", cs, mouseCS)
	startupCoreCS, startupMouseCS = cs, mouseCS
	var originalAPI, cAPI, originalPortReads []byte
	originalQueries, activeQueries := 0, 0
	original.CPU.Bus = inputReadBus{Bus: original.CPU.Bus, reads: &originalPortReads}
	installTrace := func(device *machine.Machine, d *dos.DOS, dst *[]byte, queries *int) {
		previous := device.CPU.IntHook
		device.CPU.IntHook = func(c *cpu.CPU, n uint8) bool {
			if n == 0x33 && c.R[cpu.AX] == 0 && currentInput.mouseAbsent {
				*dst = append(*dst, n)
				*dst = append(*dst, regBytes(originalRegs(device))...)
				c.R[cpu.AX] = 0
				c.R[cpu.BX] = 0
				*dst = append(*dst, regBytes(originalRegs(device))...)
				return true
			}
			if n == 0x33 && (c.R[cpu.AX] == 3 || c.R[cpu.AX] == 5) {
				(*queries)++
			}
			if n == 0x21 || n == 0x33 || n == 0x61 || n == 0x10 {
				*dst = append(*dst, n)
				*dst = append(*dst, regBytes(originalRegs(device))...)
				if n != 0x33 || (c.R[cpu.AX] != 3 && c.R[cpu.AX] != 5) {
					side := "orig"
					if device == active {
						side = "C"
					}
					evtEmit(fmt.Sprintf("%s INT%02X AX=%04X", side, n, c.R[cpu.AX]))
				}
			}
			result := previous(c, n)
			if n == 0x21 || n == 0x33 || n == 0x61 || n == 0x10 {
				*dst = append(*dst, regBytes(originalRegs(device))...)
			}
			return result
		}
	}
	installTrace(original, originalDOS, &originalAPI, &originalQueries)
	installTrace(active, activeDOS, &cAPI, &activeQueries)
	initialRAM := append([]byte(nil), original.Mem...)
	copy(initialRAM[0x22000:], icons[0x9700:])
	copy(initialRAM[0x45000:], talk)
	copy(initialRAM[0x30000:], mdl)
	copy(initialRAM[0x38000:], mch)
	copy(initialRAM[0x50000:], mapdata)
	copy(initialRAM[0x68000:], icons[0x6700:0x9700])
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

	wholeFixture := C.KiStartupFixture{}
	fixture := &wholeFixture.list
	originalBase, activeBase := original.Snapshot(), active.Snapshot()
	groups, entries := map[string]int{}, map[string]int{}
	mousePollCounts := map[string]int{}
	sampleCounts := map[uint32]int{}
	cases, audits, pixelAudits, independentAudits := 0, 0, 0, 0
	stageAudits := map[string]int{}
	exits, transfers, pauses := 0, 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	auditSameFails := map[string][]string{}
	resetCounters := func(c startupCase) {
		startupCTrace = nil
		startupCCalls = 0
		originalPolls, activePolls = 0, 0
		originalPositions, activePositions = 0, 0
		currentInput = c
		activeTicks, originalQueries, activeQueries = 0, 0, 0
		cPortReads, originalPortReads, originalAPI, cAPI = nil, nil, nil, nil
		// Fresh DOS per case: each case is a separate boot, so the arena,
		// handles, alloc/read logs and sound map must not leak across cases.
		// (Entry case 2 starved on case 1's 461KB of dead allocs.)
		originalDOS = dos.New(original, "/orig")
		activeDOS = dos.New(active, "/orig")
		originalDOSDevice, activeDOSDevice = originalDOS, activeDOS
		originalDOS.Install()
		activeDOS.Install()
		installTrace(original, originalDOS, &originalAPI, &originalQueries)
		installTrace(active, activeDOS, &cAPI, &activeQueries)
		for _, device := range []*machine.Machine{original, active} {
			device.PortLog = nil
			device.PortsIn = map[uint16]uint64{}
		}
		for _, d := range []*dos.DOS{originalDOS, activeDOS} {
			d.Mouse = dos.Mouse{X: 100, Y: 120, MaxX: 6143, MaxY: 4127, Press: [3]uint16{}, Calls: map[uint16]int{}}
			for i := 0; i < 3; i++ {
				d.Mouse.PressAt[i] = [2]uint16{100, 120}
			}
			if c.prePress && len(c.events) > 0 {
				d.Mouse.X = c.events[0].x
				d.Mouse.Y = c.events[0].y
				if c.events[0].button >= 0 {
					d.Mouse.PressButton(c.events[0].button)
				}
			}
			d.Console = nil
			d.Exited = false
			d.ExitCode = 0
		}
		fixture.core_cs = C.uint16_t(cs)
		fixture.mouse_cs = C.uint16_t(mouseCS)
		fixture.calls = 0
		fixture.unsupported = 0
		for i := range fixture.blocks {
			fixture.blocks[i].calls = 0
		}
		wholeFixture.transferred = 0
		wholeFixture.transfer_ip = 0
		wholeFixture.transfer_ss = 0
		wholeFixture.transfer_sp = 0
	}
	resetRaw := func(c startupCase) {
		original.Restore(originalBase)
		active.Restore(activeBase)
		var carve0, carve1 uint16
		if c.group == "bootchain" {
			carve0 = original.FreeSeg + 1
			carve1 = carve0 + 1 + 0x1D5E
		}
		for _, device := range []*machine.Machine{original, active} {
			copy(device.Mem, initialRAM)
			copy(device.Mem[0x22000:], icons[0x9700:])
			device.SetVideoMode(uint8(c.videoMode))
			for plane := 0; plane < 4; plane++ {
				copy(device.VGA.Planes[plane][:], planes[int(c.seed)&7][plane*65536:(plane+1)*65536])
			}
			device.Out8(0x3c4, 2)
			device.Out8(0x3c5, 15)
			for _, q := range [][2]byte{{0, 0}, {1, 0}, {3, 0}, {4, 0}, {5, 0}, {8, 255}} {
				device.Out8(0x3ce, q[0])
				device.Out8(0x3cf, q[1])
			}
			device.Read8(0xa4321)
			for _, v := range [][2]uint16{{0xd48, 0x2200}, {0xd4a, 0x226c}, {0xd4c, 0x228f}, {0xd54, 0x2284}, {0xd3a, 0x6800}, {0xd3c, 0x7800},
				{0xd38, 0x4500}, {0xd40, 0x8000}, {0xd44, 0x5000}, {0xd52, 0x7000}, {0xd56, 0xd000}, {0x9872, 0xb000}, {0xd84a, 0x3000}, {0xd84c, 0x3800}, {0xd84e, 0x2500}, {0xd850, 0x5000},
				{0x987c, 0x9000}, {0x9876, 0xe000}, {0x987a, 0x4200}, {0xe479, 0x7600}, {0xe47b, 0}, {0xd30a, 0x7000}, {0xd30e, 0x5200}} {
				putWord(device.Mem, cs, v[0], v[1])
			}
			for _, v := range [][2]uint16{{0xf78a, 0x410}, {0xf78c, 0x80}, {0xf791, 0x414}, {0xf793, 0x80}} {
				putWord(device.Mem, cs, v[0], v[1])
			}
			for _, v := range [][2]uint16{{0xd50, 0x229a}, {0xeae9, 0xe16}, {0xeaeb, cs}, {0xeaed, cs}, {0xeaef, 0x229a}, {0xeaf1, 0x1800}, {0xeaf3, 0x229a}, {0xeaf5, 0x2020}} {
				putWord(device.Mem, cs, v[0], v[1])
			}
			copy(device.Mem[0x70000:], worlds[c.scenario*22208+0x80:c.scenario*22208+0x52c0])
			copy(device.Mem[int(cs)*16+0xcf0:], worlds[c.scenario*22208:c.scenario*22208+59])
			putWord(device.Mem, cs, 0xcfd, uint16(c.player)*64)
			device.Mem[int(cs)*16+0xcff] = c.player
			device.Mem[int(cs)*16+0x845] = 0
			for i := 0; i < 4; i++ {
				device.Mem[int(cs)*16+0x846+i] = 255
			}
			for i := 0; i < 920; i++ {
				at := 0x25000 + i*8
				device.Mem[at] = 0x40
				device.Mem[at+1] = 0
				for j := 2; j < 8; j++ {
					device.Mem[at+j] = 255
				}
			}
			for _, v := range [][2]uint16{{0x988e, 0}, {0x9890, 0}, {0x9892, 0}, {0x9894, 0}, {0xd854, 0}, {0xd856, 0}} {
				putWord(device.Mem, cs, v[0], v[1])
			}
			device.Mem[int(cs)*16+0xd2c] = c.d2c
			device.Mem[int(cs)*16+0xd2d] = c.d2d
			device.Mem[int(cs)*16+0xd2a] = c.d2a
			if c.group == "bootchain" {
				seg := carve0
				putWord(device.Mem, cs, 0xd38, seg)
				seg += 0xd00
				putWord(device.Mem, cs, 0xd3a, seg)
				seg += 0x300
				putWord(device.Mem, cs, 0xd3c, seg)
				seg += 0x1c0
				putWord(device.Mem, cs, 0xd3e, seg)
				seg += 0x200
				putWord(device.Mem, cs, 0xd56, seg)
				seg += 0x40
				putWord(device.Mem, cs, 0xd40, seg)
				seg += 0x200
				putWord(device.Mem, cs, 0xd52, seg)
				seg += 0x524
				putWord(device.Mem, cs, 0xd48, seg)
				seg += 0x6c
				putWord(device.Mem, cs, 0xd4a, seg)
				seg += 0x18
				putWord(device.Mem, cs, 0xd54, seg)
				seg += 0xb
				putWord(device.Mem, cs, 0xd4c, seg)
				seg += 0xb
				putWord(device.Mem, cs, 0xd50, seg)
				putWord(device.Mem, cs, 0xd44, carve1)
				putWord(device.Mem, cs, 0xd46, carve1)
			}
			if c.anchorSS != 0 || c.anchorSP != 0 {
				putWord(device.Mem, cs, 0x9901, c.anchorSP)
				putWord(device.Mem, cs, 0x9903, c.anchorSS)
				putWord(device.Mem, c.anchorSS, c.anchorSP, c.transferIP)
			}
			device.Mem[int(cs)*16+0x98a3] = c.a3
			device.Mem[int(cs)*16+0x98a5] = c.a5
			device.Mem[machine.PSPSeg*16+0x80] = c.pspCmd
			for i := 0; i < 8; i++ {
				device.Mem[machine.PSPSeg*16+0x81+i] = c.pspFill
			}
			startupSeedName(device.Mem, c)
			if c.prealloc {
				copy(device.Mem[0x90000:], nameFont)
				putWord(device.Mem, cs, 0x8fc7, 0x9000)
			}
			for i := range device.DAC {
				device.DAC[i] = byte(i*13+int(c.seed)*7) & 63
			}
		}
		if &original.Mem[0] == &active.Mem[0] {
			panic("startup shared RAM")
		}
		for plane := 0; plane < 4; plane++ {
			if &original.VGA.Planes[plane][0] == &active.VGA.Planes[plane][0] {
				panic("startup shared VGA plane")
			}
		}
		if !bytes.Equal(original.Mem, active.Mem) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(paletteBytes(original), paletteBytes(active)) || !bytes.Equal(portMapBytes(original), portMapBytes(active)) {
			panic("independent raw fixture setup mismatch")
		}
	}
	check := func(c startupCase) bool {
		if *only != "" && c.group != *only {
			return true
		}
		if *smoke && sampleCounts[c.target] >= 4 {
			return true
		}
		sampleCounts[c.target]++
		cases++
		groups[c.group]++
		curGroup, curTarget, curScenario, curCase = c.group, c.target, c.scenario, cases
		curSide = "original"
		pcRingN, evtRingN = 0, 0
		nativeHist = map[uint32]int{}
		defer func() {
			if e := recover(); e != nil {
				panic(fmt.Sprintf("startup %s group=%s target=%X scenario=%d caseidx=%d: %v\npc: %s\nevt: %s", curSide, curGroup, curTarget, curScenario, curCase, e, ringPCs(96), ringEvts(16)))
			}
		}()
		resetRaw(c)
		resetCounters(c)
		progress := func(side string) {
			data, e := json.Marshal(map[string]any{"case": cases, "group": c.group, "target": fmt.Sprintf("%X", c.target), "scenario": c.scenario, "side": side, "completed_stages": audits})
			if e != nil {
				panic(e)
			}
			if e = os.WriteFile(*output+".progress.json", append(data, '\n'), 0600); e != nil {
				panic(e)
			}
		}
		progress("original")
		before := append([]byte(nil), original.Mem...)
		soundBefore := map[uint8]int{}
		for k, v := range originalDOS.Sound {
			soundBefore[k] = v
		}
		allocsBefore, readsBefore := len(originalDOS.Allocs), len(originalDOS.Reads)
		r := registers{AX: c.ax, BX: c.bx, CX: c.cx, DX: c.dx, SI: c.si, DI: c.di, BP: c.bp, SP: c.entrySP, SS: c.entrySS, DS: cs, ES: cs, CS: cs, IP: uint16(c.target), Flags: uint16(2 | c.seed&1)}
		if c.target == 0x10000 {
			r.SS, r.SP = loaderSS, loaderSP
		}
		if c.ds != 0 {
			r.DS = c.ds
		}
		if c.es != 0 {
			r.ES = c.es
		}
		if c.df {
			r.Flags |= 0x400
		}
		original.CPU.R = [8]uint16{r.AX, r.CX, r.DX, r.BX, r.SP, r.BP, r.SI, r.DI}
		original.CPU.Seg = [4]uint16{r.ES, r.CS, r.SS, r.DS}
		original.CPU.IP = r.IP
		original.CPU.SetFlags(r.Flags)
		// Entry (start) runs on the loader stack SS:SP=10D0:0400, whose top word
		// IS the mouse-driver head (linear 0x11100 = file 0x10200: FA 56).
		// Writing the 0xF000 return stub there decapitates the driver's cli/
		// push si, unbalances every far call and walks into F000 zeros.
		// Real boot never touches [SP] itself (pushes grow down into BSS),
		// and start never returns to the stub anyway (exit/pause ends it),
		// so entry skips the stub write on both sides.
		if c.target == 0x10356 {
			putWord(original.Mem, r.SS, r.SP, 0xf000)
			putWord(original.Mem, r.SS, r.SP+2, cs)
			putWord(cMemory, r.SS, r.SP, 0xf000)
			putWord(cMemory, r.SS, r.SP+2, cs)
		} else if c.target != 0x10000 {
			putWord(original.Mem, r.SS, r.SP, 0xf000)
			putWord(cMemory, r.SS, r.SP, 0xf000)
		}
		assignC(&m, r)
		var originalTrace []byte
		var observed []registers
		pauseHits := 0
		steps := 0
		trapRetf := os.Getenv("WOLONG_STARTUP_TRAPRETF") == "1"
		allEntries := map[uint32]int{}
		watchFired := false
		watchTickFired := false
		slotHits := 0
		_ = slotHits
		watch2239Fired := false
		waitLogN := 0
		if trapRetf {
			fmt.Fprintf(os.Stderr, "watch driver head at case start: %x\n", original.Mem[0x11100:0x11104])
		}
		firstExit := ""
		prevPC := ""
		for ; steps < 100000000; steps++ {
			cur := originalRegs(original)
			thisPC := fmt.Sprintf("%04X:%04X", cur.CS, cur.IP)
			if cur.CS != cs && cur.CS != mouseCS && firstExit == "" {
				firstExit = prevPC + " -> " + thisPC
			}
			if trapRetf && !watchFired && (original.Mem[0x11100] != 0xFA || original.Mem[0x11101] != 0x56) {
				watchFired = true
				panic(fmt.Sprintf("watch driver clobber at %04X:%04X bytes=%x regs=%+v", cur.CS, cur.IP, original.Mem[0x11100:0x11104], cur))
			}
			if trapRetf && !watchTickFired && (original.Mem[0x46C] != 0 || original.Mem[0x46D] != 0 || original.Mem[0x46E] != 0 || original.Mem[0x46F] != 0) {
				watchTickFired = true
				panic(fmt.Sprintf("watch tick write at %04X:%04X bytes=%x regs=%+v", cur.CS, cur.IP, original.Mem[0x46C:0x470], cur))
			}
			if trapRetf && c.group == "transfer" && cur.CS == cs && cur.IP == 0x1CD0 && waitLogN < 220 {
				waitLogN++
				fmt.Fprintf(os.Stderr, "tentry n=%d case=%d press=%v buttons=%d x=%d y=%d\n", waitLogN, curCase, originalDOS.Mouse.Press, originalDOS.Mouse.Buttons, originalDOS.Mouse.X, originalDOS.Mouse.Y)
			}
			if trapRetf && c.group == "transfer" && cur.CS == cs && (cur.IP == 0x2216 || cur.IP == 0x222B) && waitLogN < 220 {
				waitLogN++
				ret := uint16(original.Mem[int(cur.SS)*16+int(cur.SP)]) | uint16(original.Mem[int(cur.SS)*16+int(cur.SP)+1])<<8
				fmt.Fprintf(os.Stderr, "t2216 n=%d case=%d polls=%d ip=%04X ret=%04X press=%v calls5=%d calls6=%d pr=%d\n", waitLogN, curCase, originalPolls, cur.IP, ret, originalDOS.Mouse.Press, originalDOS.Mouse.Calls[5], originalDOS.Mouse.Calls[6], len(originalDOS.Mouse.PressReads))
			}
			if trapRetf && c.group == "transfer" && cur.CS == cs && cur.IP == 0x2239 && !watch2239Fired {
				watch2239Fired = true
				fmt.Fprintf(os.Stderr, "rt2239 bytes=%x\n", original.Mem[int(cs)*16+0x2229:int(cs)*16+0x2270])
				fmt.Fprintf(os.Stderr, "rt2000 bytes=%x\n", original.Mem[0x20000:0x20100])
				fmt.Fprintf(os.Stderr, "rt1110 bytes=%x\n", original.Mem[0x11100:0x11200])
				fmt.Fprintf(os.Stderr, "rt2270 bytes=%x\n", original.Mem[int(cs)*16+0x2270:int(cs)*16+0x22D0])
			}
			if trapRetf && c.group == "boot" && cur.CS == cs && cur.IP == 0x8BA4 && waitLogN < 80 {
				waitLogN++
				fmt.Fprintf(os.Stderr, "code18 n=%d case=%d polls=%d al=%02X x=%d y=%d\n", waitLogN, curCase, originalPolls, cur.AX&0xFF, originalDOS.Mouse.X, originalDOS.Mouse.Y)
			}
			if trapRetf && c.group == "boot" && cur.CS == cs && cur.IP == 0x8E7A && waitLogN < 80 {
				waitLogN++
				fmt.Fprintf(os.Stderr, "code8E n=%d case=%d polls=%d al=%02X x=%d y=%d\n", waitLogN, curCase, originalPolls, cur.AX&0xFF, originalDOS.Mouse.X, originalDOS.Mouse.Y)
			}
			if trapRetf && c.group == "boot" && cur.CS == cs && cur.IP == 0x9029 && waitLogN < 100 {
				waitLogN++
				fmt.Fprintf(os.Stderr, "code90 n=%d case=%d polls=%d al=%02X x=%d y=%d\n", waitLogN, curCase, originalPolls, cur.AX&0xFF, originalDOS.Mouse.X, originalDOS.Mouse.Y)
			}
			if trapRetf && c.group == "boot" && cur.CS == cs && cur.IP == 0x843B && waitLogN < 160 {
				waitLogN++
				fmt.Fprintf(os.Stderr, "code84 n=%d case=%d polls=%d al=%02X x=%d y=%d\n", waitLogN, curCase, originalPolls, cur.AX&0xFF, originalDOS.Mouse.X, originalDOS.Mouse.Y)
			}
			if trapRetf && c.group == "boot" && cur.CS == cs && cur.IP == 0x8E19 && waitLogN < 200 {
				waitLogN++
				b0 := uint16(original.Mem[int(cur.SS)*16+int(cur.BP)]) | uint16(original.Mem[int(cur.SS)*16+int(cur.BP)+1])<<8
				b2 := uint16(original.Mem[int(cur.SS)*16+int(cur.BP)+2]) | uint16(original.Mem[int(cur.SS)*16+int(cur.BP)+3])<<8
				fmt.Fprintf(os.Stderr, "box8E n=%d case=%d polls=%d x0=%d y0=%d mx=%d my=%d\n", waitLogN, curCase, originalPolls, b0, b2, originalDOS.Mouse.X, originalDOS.Mouse.Y)
			}
			if trapRetf && c.group == "boot" && (curCase == 1 || curCase == 3) && cur.CS == cs && cur.IP == 0x005B && waitLogN < 140 {
				waitLogN++
				if curCase == 1 {
					snapB5RAM = append(snapB5RAM[:0], original.Mem...)
					fmt.Fprintf(os.Stderr, "b5entry n=%d bp=%04X sp=%04X ss=%04X snap=%d\n", waitLogN, cur.BP, cur.SP, cur.SS, len(snapB5RAM))
				} else if len(snapB5RAM) == len(original.Mem) {
					nr := 0
					for i := 0; i < len(original.Mem) && nr < 30; {
						if original.Mem[i] == snapB5RAM[i] {
							i++
							continue
						}
						j := i
						for j < len(original.Mem) && original.Mem[j] != snapB5RAM[j] && j-i < 32 {
							j++
						}
						fmt.Fprintf(os.Stderr, "b5diff off=%06X len=%d c1=%x c3=%x\n", i, j-i, snapB5RAM[i:j], original.Mem[i:j])
						nr++
						i = j
					}
				}
			}
			if trapRetf && c.group == "boot" && cur.CS == cs && cur.IP == 0x8E9F && waitLogN < 120 {
				waitLogN++
				if curCase == 1 && originalPolls == 7 {
					fmt.Fprintf(os.Stderr, "bootdrv entry=%x\n", original.Mem[0x10110:0x10130])
					fmt.Fprintf(os.Stderr, "bootdrv f906=%x\n", original.Mem[0x1FA16:0x1FB16])
				}
				if (curCase == 1 || curCase == 3) && originalPolls == 7 {
					if curCase == 1 {
						snapCase1RAM = append(snapCase1RAM[:0], original.Mem...)
					} else if len(snapCase1RAM) == len(original.Mem) {
						nr := 0
						for i := 0; i < len(original.Mem) && nr < 25; {
							if original.Mem[i] == snapCase1RAM[i] {
								i++
								continue
							}
							j := i
							for j < len(original.Mem) && original.Mem[j] != snapCase1RAM[j] && j-i < 64 {
								j++
							}
							fmt.Fprintf(os.Stderr, "ramdiff off=%06X len=%d c1=%x c3=%x\n", i, j-i, snapCase1RAM[i:j], original.Mem[i:j])
							nr++
							i = j
						}
					}
				}
				st := ""
				for i := 0; i < 8; i++ {
					a := int(cur.SS)*16 + int(cur.SP) + i*2
					st += fmt.Sprintf(" %04X", uint16(original.Mem[a])|uint16(original.Mem[a+1])<<8)
				}
				fmt.Fprintf(os.Stderr, "m2exit n=%d case=%d polls=%d ax=%04X bx=%04X cx=%04X dx=%04X si=%04X di=%04X bp=%04X sp=%04X ss=%04X ds=%04X es=%04X ram=%x stk=%s\n", waitLogN, curCase, originalPolls, cur.AX, cur.BX, cur.CX, cur.DX, cur.SI, cur.DI, cur.BP, cur.SP, cur.SS, cur.DS, cur.ES, sha256.Sum256(original.Mem), st)
			}
			if trapRetf && c.group == "boot" && cur.CS == cs && cur.IP == 0x21E7 && waitLogN < 40 {
				ret := uint16(original.Mem[int(cur.SS)*16+int(cur.SP)]) | uint16(original.Mem[int(cur.SS)*16+int(cur.SP)+1])<<8
				waitLogN++
				fmt.Fprintf(os.Stderr, "waitnav n=%d case=%d polls=%d ret=%04X\n", waitLogN, curCase, originalPolls, ret)
			}
			if trapRetf && (c.group == "load" || c.group == "boot") && cur.CS == cs && cur.IP == 0x8B9A && slotHits < 80 {
				ret := uint16(original.Mem[int(cur.SS)*16+int(cur.SP)]) | uint16(original.Mem[int(cur.SS)*16+int(cur.SP)+1])<<8
				slotHits++
				if slotHits <= 4 || slotHits%20 == 0 {
					ds := uint16(original.Mem[int(cs)*16+0x987C]) | uint16(original.Mem[int(cs)*16+0x987D])<<8
					var sb strings.Builder
					for n := 0; n < 4; n++ {
						a := int(ds)*16 + n*128 + 0x40
						fmt.Fprintf(&sb, " s%d=%04X", n, uint16(original.Mem[a])|uint16(original.Mem[a+1])<<8)
					}
					fmt.Fprintf(os.Stderr, "wait18 hit=%d case=%d polls=%d ret=%04X bp=%04X ds=%04X%s\n", slotHits, curCase, originalPolls, ret, cur.BP, ds, sb.String())
				}
				_ = slotHits
			}
			if trapRetf && c.group == "load" && cur.CS == cs && cur.IP == 0x8BB3 {
				at := int(cur.DS)*16 + int(cur.BX) + 0x40
				slotHits++
				fmt.Fprintf(os.Stderr, "slot hit=%d case=%d bx(entry)=%04X polls=%d ax=%04X bx=%04X ds=%04X mem=%x\n", slotHits, curCase, c.bx, originalPolls, cur.AX, cur.BX, cur.DS, original.Mem[at:at+4])
				if slotHits == 1 {
					var sb strings.Builder
					for n := 0; n < 8; n++ {
						a := int(cur.DS)*16 + n*128 + 0x40
						fmt.Fprintf(&sb, " s%d=%04X", n, uint16(original.Mem[a])|uint16(original.Mem[a+1])<<8)
					}
					fmt.Fprintf(os.Stderr, "slots case=%d target=%X:%s\n", curCase, c.target, sb.String())
				}
				if slotHits >= 6 {
					panic(fmt.Sprintf("slot trap done case=%d", curCase))
				}
			}
			if trapRetf && c.group == "entry" && cur.CS == mouseCS && cur.IP == 0xD {
				base := int(cur.SS)*16 + int(cur.SP)
				dbase := int(mouseCS)*16
				panic(fmt.Sprintf("retf trap ss:sp=%04X:%04X stack=%x driver=%x regs=%+v", cur.SS, cur.SP, original.Mem[base:base+16], original.Mem[dbase:dbase+32], cur))
			}
			prevPC = thisPC
			if originalDOS.Exited {
				break
			}
			if cur.CS == cs && cur.IP == 0xf000 {
				break
			}
			if c.pauseIters > 0 && cur.CS == cs && cur.IP == 0x1c1f {
				pauseHits++
				if pauseHits >= c.pauseIters {
					break
				}
			}
			loc := uint32(cur.IP) + 0x10000
			if cur.CS == mouseCS {
				loc = uint32(cur.IP) + 0x20000
			}
			ringLoc := loc
			if cur.CS != cs && cur.CS != mouseCS {
				ringLoc = uint32(cur.CS)<<16 | uint32(cur.IP)
			}
			pcRing[pcRingN%uint64(len(pcRing))] = ringLoc
			pcRingN++
			if cur.CS == cs {
				startupCheckpoint(original, originalDOS, cur.IP)
			}
			_, known := functions[loc]
			if trapRetf && (cur.CS == cs || cur.CS == mouseCS) {
				allEntries[loc]++
			}
			if (cur.CS == cs || cur.CS == mouseCS) && known {
				entries[entryName(loc)]++
				originalTrace = append(originalTrace, byte(cur.IP), byte(cur.IP>>8))
				originalTrace = append(originalTrace, regBytes(cur)...)
				for i := uint16(0); i < 32; i++ {
					originalTrace = append(originalTrace, original.Mem[(uint32(cur.SS)*16+uint32(cur.SP+i))&0xfffff])
				}
				originalTrace = append(originalTrace, deviceBytes(original)[:28]...)
			}
			if cur.CS == cs && known {
				observed = append(observed, cur)
			}
			if err := original.Step(); err != nil {
				panic(err)
			}
		}
		if steps >= 100000000 {
			cur := originalRegs(original)
			stuck := original.Mem[int(cur.CS)*16+int(cur.IP) : int(cur.CS)*16+int(cur.IP)+8]
			panic(fmt.Sprintf("bounded startup original group=%s target=%X polls=%d queries=%d state=%+v bytes=%x firstexit=%s", c.group, c.target, originalPolls, originalQueries, cur, stuck, firstExit))
		}
		if trapRetf && len(allEntries) > 0 {
			type kv struct {
				loc uint32
				n   int
			}
			var top []kv
			for l, n := range allEntries {
				top = append(top, kv{l, n})
			}
			sort.Slice(top, func(i, j int) bool { return top[i].n > top[j].n })
			var sb strings.Builder
			for i := 0; i < len(top) && i < 30; i++ {
				fmt.Fprintf(&sb, " %X:%d", top[i].loc, top[i].n)
			}
			fmt.Fprintf(os.Stderr, "lochist group=%s target=%X case=%d steps=%d top:%s\n", c.group, c.target, curCase, steps, sb.String())
			var wl strings.Builder
			for _, w := range []uint32{0x11D8E, 0x11E06, 0x11E17, 0x15358, 0x19377, 0x13E11} {
				fmt.Fprintf(&wl, " %X:%d", w, allEntries[w])
			}
			fmt.Fprintf(os.Stderr, "locwatch group=%s target=%X case=%d:%s\n", c.group, c.target, curCase, wl.String())
		}
		if trapRetf {
			return true
		}
		nOrig, origObs := startupAudit(before, original.Mem, cs, c, r, observed, original, originalDOS, soundBefore, allocsBefore, readsBefore)
		cSoundBefore := map[uint8]int{}
		for k, v := range activeDOS.Sound {
			cSoundBefore[k] = v
		}
		cAllocsBefore, cReadsBefore := len(activeDOS.Allocs), len(activeDOS.Reads)
		curSide = "native-c"
		progress("native-c")
		if c.pauseIters > 0 {
			if pauseHits != c.pauseIters {
				panic(fmt.Sprintf("pause count group=%s want=%d got=%d", c.group, c.pauseIters, pauseHits))
			}
			C.startup_run_pause(&m, C.uint16_t(c.target), &wholeFixture, C.uint(pauseHits))
			if C.startup_is_paused() != 1 {
				panic("native warmup did not pause at loop entry")
			}
			pauses++
		} else if c.transfer {
			C.startup_run(&m, C.uint16_t(c.target), &wholeFixture)
			if wholeFixture.transferred != 1 {
				panic("native transfer did not fire")
			}
			transfers++
			if c.chain {
				C.startup_continue(&m, 0x006A, &wholeFixture)
				C.startup_continue(&m, 0x0043, &wholeFixture)
			}
		} else {
			C.startup_run(&m, C.uint16_t(c.target), &wholeFixture)
		}
		if startupCCalls > 2000000 || fixture.unsupported != 0 {
			panic(fmt.Sprintf("startup unsupported=%X calls=%d group=%s", fixture.unsupported, startupCCalls, c.group))
		}
		cFinal := cRegs(&m)
		active.CPU.R = [8]uint16{cFinal.AX, cFinal.CX, cFinal.DX, cFinal.BX, cFinal.SP, cFinal.BP, cFinal.SI, cFinal.DI}
		active.CPU.Seg = [4]uint16{cFinal.ES, cFinal.CS, cFinal.SS, cFinal.DS}
		active.CPU.IP = cFinal.IP
		active.CPU.SetFlags(cFinal.Flags)
		cTrace := startupCTrace
		nC, cObs := startupAudit(before, cMemory, cs, c, r, parseCTraceIPs(cTrace), active, activeDOS, cSoundBefore, cAllocsBefore, cReadsBefore)
		independentAudits += nOrig + nC
		auditDiv := ""
		if len(origObs) != len(cObs) {
			auditDiv = fmt.Sprintf("startup audit stream length group=%s target=%X orig=%d c=%d origObs=%q cObs=%q", c.group, c.target, len(origObs), len(cObs), origObs, cObs)
		} else {
			for i := range origObs {
				if origObs[i] != cObs[i] {
					auditDiv = fmt.Sprintf("startup audit stream div group=%s target=%X i=%d orig=%q c=%q", c.group, c.target, i, origObs[i], cObs[i])
					break
				}
			}
		}
		if auditDiv != "" {
			key := fmt.Sprintf("%s/%X", c.group, c.target)
			auditSameFails[key] = append(auditSameFails[key], fmt.Sprintf("caseidx=%d div %s", cases, auditDiv))
		}
		for _, o := range origObs {
			if strings.HasPrefix(o, "FAIL") {
				key := fmt.Sprintf("%s/%X", c.group, c.target)
				auditSameFails[key] = append(auditSameFails[key], fmt.Sprintf("caseidx=%d %s", cases, o))
			}
		}
		transferOK := wholeFixture.transferred == 0
		if c.transfer {
			transferOK = wholeFixture.transferred == 1 && uint16(wholeFixture.transfer_ip) == c.transferIP
		}
		exitOK := C.startup_did_exit() == 0 && !originalDOS.Exited
		if c.expectExit {
			exitOK = C.startup_did_exit() == 1 && originalDOS.Exited && uint8(C.startup_exit_status()) == originalDOS.ExitCode
			exits++
		}
		want, got := originalRegs(original), cRegs(&m)
		ramOff, ramOrig, ramC := firstDiff(original.Mem, cMemory)
		plnOff, plnOrig, plnC := firstDiff(original.VGA.Raw(), active.VGA.Raw())
		originalPorts, cPorts := portBytes(original), portBytes(active)
		originalState, cState := deviceBytes(original), deviceBytes(active)
		if !transferOK || !exitOK || want != got || !bytes.Equal(original.Mem, cMemory) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(originalState, cState) || !bytes.Equal(originalPorts, cPorts) || !bytes.Equal(originalTrace, cTrace) || !bytes.Equal(originalAPI, cAPI) || !bytes.Equal(originalPortReads, cPortReads) || !bytes.Equal(mouseBytes(originalDOS), mouseBytes(activeDOS)) || !bytes.Equal(soundBytes(originalDOS), soundBytes(activeDOS)) || !bytes.Equal(consoleBytes(originalDOS), consoleBytes(activeDOS)) || !bytes.Equal(exitBytes(originalDOS), exitBytes(activeDOS)) || !bytes.Equal(periodicBytes(original), periodicBytes(active)) || !bytes.Equal(allocBytes(originalDOS), allocBytes(activeDOS)) || !bytes.Equal(paletteBytes(original), paletteBytes(active)) || !bytes.Equal(portMapBytes(original), portMapBytes(active)) || originalQueries != activeQueries || originalPolls != activePolls {
			mismatch = map[string]any{"group": c.group, "case": cases - 1, "input": r, "original": want, "c": got, "expected_exit": c.expectExit, "native_exited": int(C.startup_did_exit()), "native_exit_code": uint8(C.startup_exit_status()), "native_transferred": int(wholeFixture.transferred), "native_transfer_ip": uint16(wholeFixture.transfer_ip), "ram_first_diff_off": ramOff, "ram_first_diff_orig": ramOrig, "ram_first_diff_c": ramC, "planes_first_diff_off": plnOff, "planes_first_diff_orig": plnOrig, "planes_first_diff_c": plnC, "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(cMemory)), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(active.VGA.Raw())), "original_trace": hex.EncodeToString(originalTrace), "c_trace": hex.EncodeToString(cTrace), "original_ports": hex.EncodeToString(originalPorts), "c_ports": hex.EncodeToString(cPorts), "original_in": hex.EncodeToString(originalPortReads), "c_in": hex.EncodeToString(cPortReads), "original_api": hex.EncodeToString(originalAPI), "c_api": hex.EncodeToString(cAPI), "original_console": hex.EncodeToString(consoleBytes(originalDOS)), "c_console": hex.EncodeToString(consoleBytes(activeDOS)), "original_allocs": string(allocBytes(originalDOS)), "c_allocs": string(allocBytes(activeDOS)), "original_queries": originalQueries, "c_queries": activeQueries, "original_polls": originalPolls, "c_polls": activePolls}
			return false
		}
		if auditDiv != "" {
			panic("startup audit nondeterminism (no state div): " + auditDiv)
		}
		audits++
		if wholeFixture.transferred == 1 {
			stageAudits["transfer"]++
		} else if C.startup_did_exit() == 1 {
			stageAudits["exit"]++
		} else if c.pauseIters > 0 {
			stageAudits["pause"]++
		} else {
			stageAudits["return"]++
		}
		pixelAudits++
		mousePollCounts[c.group] += originalQueries
		if !bytes.Equal(indexed(original), indexed(active)) {
			panic("startup indexed pixels")
		}
		for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), originalState, originalPorts, originalTrace, originalAPI, originalPortReads, mouseBytes(originalDOS), soundBytes(originalDOS), consoleBytes(originalDOS), exitBytes(originalDOS), periodicBytes(original), allocBytes(originalDOS), paletteBytes(original), portMapBytes(original)} {
			originalDigest.Write(b)
		}
		for _, b := range [][]byte{regBytes(got), cMemory, active.VGA.Raw(), cState, cPorts, cTrace, cAPI, cPortReads, mouseBytes(activeDOS), soundBytes(activeDOS), consoleBytes(activeDOS), exitBytes(activeDOS), periodicBytes(active), allocBytes(activeDOS), paletteBytes(active), portMapBytes(active)} {
			cDigest.Write(b)
		}
		return true
	}

	base := func(group string, target uint32, scenario int) startupCase {
		return startupCase{group: group, target: target, scenario: scenario, player: 0, entrySS: 0xf000, entrySP: 0x6ffe, returnIP: 0xf000, ax: 0x5500, bx: 0x1234, cx: 0x56, dx: 15, bp: 0x5000, di: 0x5000, si: 0x4240, seed: 3, videoMode: 0x12, pspCmd: 0, pspFill: 0}
	}
	all := startupCases(base, nameFont)

	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{
		"schema": "wolong-c-startup-parity-v1", "input_sha256": inputHash, "icon_sha256": assetHash,
		"routine_sha256": routineHashes, "new_routine_sha256": routineHashes, "new_code_block_sha256": blockHashes, "new_probe_sha256": fmt.Sprintf("%x", sha256.Sum256(probeRaw)),
		"source_manifest": map[string]string{"path": "workplace/matching-decompilation/c-startup/results/c-source.sha256", "sha256": manifestHash},
		"cases": cases, "groups": groups, "entries_seen": entries, "full_ram_plane_audits": audits, "indexed_content_audits": pixelAudits, "palette_dac_audits": audits, "independent_startup_audits": independentAudits, "stage_audits": stageAudits,
		"original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)),
		"mouse_queries": mousePollCounts,
		"scope":               "Original DOS entry, command dispatch, new-game/load, idle main loop, naming UI and video init with real C callees; fixed raw inputs. Complete loader chain, normal player paths and C machine-code matching remain incomplete.",
		"platform":            "Independent pinned DOS/font/VGA services; full RAM/planes/DAC/registers/FLAGS/stack/IN/OUT/API/console/exit/allocs; IRQ0Every unreachable; no wall-clock claim",
		"game_palette_sha256": fmt.Sprintf("%x", sha256.Sum256(gamePalette)), "game_palette_size": len(gamePalette), "additional_assets": identities,
		"independent_raw_initialization": true, "cross_machine_snapshot_initialization": false,
		"exit_cases": exits, "transfer_cases": transfers, "pause_cases": pauses,
		"passed":     mismatch == nil, "mismatch": mismatch, "audit_same_fails": auditSameFails, "c_machine_code_match": false,
	}

	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if e := os.WriteFile(*output, b, 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Startup original/C %d; complete RAM/plane %d; content %d; pass=%v\n", cases, audits, pixelAudits, mismatch == nil)
	if mismatch != nil {
		fmt.Printf("mismatch group=%s case=%v ramoff=%v ram=%v/%v plnoff=%v orig=%+v c=%+v q=%v/%v p=%v/%v\n", mismatch["group"], mismatch["case"], mismatch["ram_first_diff_off"], mismatch["ram_first_diff_orig"], mismatch["ram_first_diff_c"], mismatch["planes_first_diff_off"], mismatch["original"], mismatch["c"], mismatch["original_queries"], mismatch["c_queries"], mismatch["original_polls"], mismatch["c_polls"])
	}
	if mismatch != nil {
		os.Exit(1)
	}
}
