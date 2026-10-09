//go:build matching_outcome

// Original automatic-combat and city-transfer closures.
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
#include "/repo/tools/c_recovery/outcome_fixture.h"
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
var routeOriginalPortReads []byte
var outcomeOriginalPortReads []byte
var routeBeforePlanes [4][65536]byte
var outcomeBeforePlanes [4][65536]byte
var mainBeforeDAC [768]byte
var activeTicks int
var currentInput resumeCase
var originalPositions, activePositions int
var originalPolls, activePolls, originalDrags, activeDrags int
var originalDragOn, activeDragOn bool
var originalDOSDevice, activeDOSDevice *dos.DOS

func listCheckpoint(device *machine.Machine, d *dos.DOS, t uint16) {
	polls, drags, on := &activePolls, &activeDrags, &activeDragOn
	if device != activeDevice {
		polls, drags, on = &originalPolls, &originalDrags, &originalDragOn
	}
	if t == 0x84dd {
		*drags = 0
		*on = true
		return
	}
	if t != 0x21e7 {
		return
	}
	event := listEvent{100, 120, 1}
	if *polls < len(currentInput.events) {
		event = currentInput.events[*polls]
	}
	(*polls)++
	d.Mouse.X = event.x
	d.Mouse.Y = event.y
	d.Mouse.ReleaseButton(0)
	d.Mouse.ReleaseButton(1)
	d.Mouse.PressButton(event.button)
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

var functions = map[uint32]int{0x11d8e: 137, 0x1ec82: 94, 0x1ece0: 28, 0x15358: 110, 0x15609: 34, 0x1563b: 40, 0x154fc: 54, 0x155ec: 13, 0x15828: 55, 0x153c6: 144, 0x15456: 57, 0x1548f: 109, 0x15538: 15, 0x15547: 95, 0x15695: 128, 0x155a6: 70, 0x12fbf: 79, 0x157fe: 42, 0x15715: 122, 0x1578f: 111, 0x130cb: 8, 0x13119: 31, 0x122db: 163, 0x12286: 85, 0x1237e: 129, 0x1585f: 58, 0x15899: 167, 0x15940: 80, 0x12ad2: 34, 0x15990: 22, 0x1301c: 50, 0x12bd9: 121, 0x12c52: 141, 0x12cdf: 91, 0x12d3a: 30, 0x12fb1: 14, 0x12d58: 96, 0x12db8: 59, 0x12df3: 64, 0x130f0: 26, 0x1310a: 15, 0x13091: 58, 0x12e33: 86, 0x12e89: 114, 0x12efb: 118, 0x12f71: 64, 0x13e11: 84, 0x13e65: 41, 0x13e8e: 111, 0x15673: 34, 0x131ae: 68, 0x13496: 16, 0x13507: 19, 0x13dc9: 72, 0x1320c: 20, 0x13220: 66, 0x13262: 71, 0x132a9: 64, 0x132e9: 62, 0x13327: 97, 0x13388: 98, 0x133ea: 19, 0x133fd: 136, 0x13485: 17, 0x134a6: 11, 0x134b1: 86, 0x1351a: 12, 0x13526: 133, 0x135ab: 66, 0x135ed: 76, 0x13639: 48, 0x13669: 46, 0x13697: 45, 0x136c4: 78, 0x13712: 95, 0x13771: 103, 0x137d8: 29, 0x137f5: 59, 0x13138: 55, 0x14502: 70, 0x16a3d: 94, 0x123ff: 57, 0x12438: 33, 0x150d7: 73, 0x12078: 94, 0x120d6: 123, 0x138c7: 31, 0x138e6: 28, 0x13902: 230, 0x139e8: 288, 0x13c3d: 92, 0x13b7e: 43, 0x13d09: 60, 0x13d45: 35, 0x13c99: 39, 0x13cdc: 45, 0x19321: 21, 0x187ff: 17, 0x13d68: 41, 0x11d46: 72, 0x12216: 21, 0x17c6e: 147, 0x17d0d: 58, 0x17d47: 24, 0x17d5f: 52, 0x17da5: 30, 0x17dc3: 26, 0x17ddd: 13, 0x17dea: 2, 0x17dec: 5, 0x17df1: 4, 0x101b4: 33, 0x101db: 51, 0x193e9: 32, 0x167cd: 25, 0x167e6: 32, 0x16806: 32, 0x16826: 32, 0x1e3c0: 23, 0x1e3d7: 68, 0x1e41b: 56, 0x1e453: 38, 0x1895d: 71, 0x189de: 18, 0x1d5d4: 65, 0x10c14: 76, 0x10c60: 23, 0x10c77: 53, 0x19796: 45, 0x197c3: 45, 0x1f9b0: 107, 0x1fa1b: 28, 0x1fa37: 107, 0x1faa2: 32, 0x1fac2: 79, 0x1fb11: 24, 0x1f888: 176, 0x1f938: 97, 0x1f999: 23, 0x1c7f4: 74, 0x1c673: 59, 0x1f020: 288, 0x1f140: 59, 0x1f17b: 39, 0x1f1a3: 203, 0x10aaa: 47, 0x10ad9: 109, 0x10cac: 23, 0x10cc3: 27, 0x1c61f: 52, 0x1c6bf: 55, 0x1c6ae: 17, 0x1c6f6: 86, 0x1c74c: 41, 0x1c775: 25, 0x1c78e: 27, 0x10bcd: 71, 0x1030f: 40, 0x10337: 4, 0x1e993: 20, 0x1e9a7: 26, 0x1e9c1: 76, 0x1fb29: 126, 0x1fba7: 84, 0x1f6dc: 68, 0x1f878: 16, 0x1f720: 62, 0x1f7a4: 212, 0x106f5: 4, 0x106fd: 4, 0x1062f: 107, 0x1069a: 68, 0x106de: 23, 0x11e17: 47, 0x106f9: 4, 0x1075b: 119, 0x107d2: 115, 0x1084a: 90, 0x18853: 48, 0x189a4: 58, 0x1e38c: 26, 0x1f4df: 73, 0x18810: 67, 0x121e7: 47, 0x20000: 14, 0x2002e: 65, 0x2006f: 1, 0x20070: 42, 0x2009a: 35, 0x200bd: 3, 0x200c0: 54, 0x20137: 52, 0x2016b: 50, 0x2019d: 41, 0x201c6: 30, 0x201e4: 40, 0x2020c: 61, 0x20249: 87, 0x202a0: 29, 0x202bd: 65, 0x202fe: 29, 0x1036f: 84, 0x103c3: 35, 0x103e6: 46, 0x10414: 161, 0x104b5: 74, 0x104ff: 78, 0x1054d: 78, 0x1059b: 82, 0x105ed: 1, 0x1061f: 16, 0x10b46: 105, 0x10baf: 30, 0x19479: 70, 0x194bf: 80, 0x187af: 29, 0x1e378: 20, 0x1f4a2: 59, 0x10241: 129, 0x102c2: 14, 0x100df: 118, 0x1d46a: 25, 0x1d483: 32, 0x1d4c7: 88, 0x1d615: 85, 0x1d66a: 257, 0x1d76b: 23, 0x1d782: 20, 0x1d796: 81, 0x1d7e7: 29, 0x1d804: 70, 0x11cc9: 7, 0x12533: 112, 0x12af4: 54, 0x12b2a: 18, 0x12b3c: 108, 0x15d19: 162, 0x11f30: 42, 0x11f5a: 37, 0x13b08: 82, 0x15c58: 76, 0x19541: 112, 0x195b1: 24, 0x196ed: 101, 0x19752: 68, 0x1d4a3: 36, 0x1d51f: 181, 0x19409: 81, 0x1945a: 31, 0x1222b: 91, 0x20101: 54, 0x10701: 90, 0x108b2: 41, 0x108db: 41, 0x10904: 53, 0x10939: 34, 0x1095b: 35, 0x1097e: 6, 0x10984: 43, 0x1f75e: 67}

var listFunctions = map[uint32]int{0x16909: 149, 0x17400: 59, 0x1748f: 193, 0x181c0: 78, 0x1820e: 129, 0x18412: 70, 0x18458: 1, 0x18463: 89, 0x184bc: 33, 0x184dd: 61, 0x1851a: 44, 0x18546: 57, 0x18607: 91, 0x18662: 103, 0x186c9: 50, 0x18713: 66, 0x18755: 90, 0x1743b: 36, 0x1745f: 48, 0x1828f: 387, 0x1857f: 51, 0x185b2: 85}

func init() {
	functions[0x15e60] = 32
	functions[0x16909] = 149
	functions[0x17400] = 59
	functions[0x1748f] = 193
	functions[0x181c0] = 78
	functions[0x1820e] = 129
	functions[0x18412] = 70
	functions[0x18458] = 1
	functions[0x18463] = 89
	functions[0x184bc] = 33
	functions[0x184dd] = 61
	functions[0x1851a] = 44
	functions[0x18546] = 57
	functions[0x18607] = 91
	functions[0x18662] = 103
	functions[0x186c9] = 50
	functions[0x18713] = 66
	functions[0x18755] = 90
	functions[0x1461d] = 33
	functions[0x14698] = 127
	functions[0x14717] = 51
	functions[0x15e80] = 55
	functions[0x15eb7] = 112
	functions[0x15f27] = 54
	functions[0x15f5d] = 34
	functions[0x15f7f] = 43
	functions[0x1699e] = 138
	functions[0x16e8f] = 58
	functions[0x16ec9] = 93
	functions[0x16f26] = 96
	functions[0x16f86] = 76
	functions[0x16fd2] = 86
}

func init() {
	for k, v := range listFunctions {
		functions[k] = v
	}
}

func init() {
	functions[0x1716d] = 59
	functions[0x171d3] = 68
	functions[0x175fa] = 66
	functions[0x17663] = 61
	functions[0x1770c] = 226
	functions[0x178a7] = 62
	functions[0x17906] = 62
	functions[0x1799c] = 222
	functions[0x17a7a] = 75
	functions[0x17b3c] = 51
	functions[0x17bc0] = 174
	functions[0x171a8] = 43
	functions[0x17217] = 54
	functions[0x1724d] = 48
	functions[0x1727d] = 250
	functions[0x1763c] = 39
	functions[0x176a0] = 60
	functions[0x176dc] = 48
	functions[0x178e5] = 33
	functions[0x17944] = 40
	functions[0x1796c] = 48
	functions[0x17b6f] = 33
	functions[0x17b90] = 48
}

func entryName(loc uint32) string {
	if loc == 0x11be0 {
		return "sub_11BE0"
	}
	if loc == 0x1491b {
		return "raw_entry_1491B"
	}
	if loc == 0x159d0 {
		return "nullsub_1"
	}
	if loc == 0x18458 {
		return "nullsub_2"
	}
	if loc == 0x105ed {
		return "nullsub_5"
	}
	if loc == 0x2006f {
		return "nullsub_4"
	}
	for _, v := range []uint32{0x171a8, 0x17217, 0x1724d, 0x1727d, 0x1763c, 0x176a0, 0x176dc, 0x178e5, 0x17944, 0x1796c, 0x17b6f, 0x17b90, 0x1743b, 0x1745f, 0x1828f, 0x1857f, 0x185b2, 0x1d51f, 0x19409, 0x1945a, 0x1222b, 0x20101, 0x108b2, 0x108db, 0x10904, 0x10939, 0x1095b, 0x1097e, 0x10984, 0x10701, 0x1f75e} {
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

type listEvent struct {
	x, y   uint16
	button int
}

type bootstrapArmy struct {
	slot                  int
	flags                 byte
	x, y, offset, segment uint16
}

type resumeCase struct {
	outcome                                            outcomeVector
	route                                              routeVector
	interaction                                        interactionVector
	decodedMapBase, occupancyBase                      uint16
	mapX, mapY, year                                   uint16
	month, day, cityOwner, cityKind, centerTile, count byte
	repeat                                             int
	armies                                             []bootstrapArmy
	neighbors                                          [4]byte
	generalFlags, war, lead, scorePoison               byte

	savedSS, savedSP, entrySS, entrySP, returnIP uint16
	palette                                      [48]byte

	trust, lordDuty, aggression, ownCities, foreignCities, ally, relation, candidate, candidateRelation byte
	captor, priority, newPriority, soundFlags, reason, reasonMask, attempted                            byte
	aptitudes                                                                                           [3]byte
	ownMoney, foreignMoney                                                                              int32
	income, expense                                                                                     uint32
	fiscal                                                                                              [8]uint16

	currentCity, targetCity, diCity                                                  int
	armyStage, armyFlags, cityFlags, threat, factionFlags, queue16, queue17, rngSeed uint8
	atTarget                                                                         bool
	mismatch                                                                         int
	cameraPixelX, cameraPixelY, lastMouseX, lastMouseY, screenX, screenY             uint16
	mapEvents                                                                        []strategyEvent
	menuChoices                                                                      []int
	tile                                                                             byte

	owner, kind, art                                                   int
	capital                                                            bool
	production, total                                                  uint16
	garrison, growth, prevention, morale                               byte
	menuRow                                                            int
	liveDispatch                                                       bool
	city, foreign, operation                                           int
	oldOfficer, budget                                                 byte
	general                                                            int
	reserves                                                           [3]uint16
	slotTypes                                                          [6]byte
	slotMen                                                            [6]byte
	poison                                                             bool
	family, profile                                                    int
	paramFaction                                                       byte
	altDS                                                              bool
	queryX, queryY                                                     uint16
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
	player                                                             byte
	top, column                                                        int
	events                                                             []listEvent
	drag                                                               []uint16
	verdict                                                            byte
}

func init() {
	functions[0x16c5e] = 52
	functions[0x16c92] = 196
	functions[0x16d56] = 25
	functions[0x16d6f] = 57
	functions[0x16da8] = 85
	functions[0x16dfd] = 131
	functions[0x16e80] = 15
}

func init() {
	functions[0x16265] = 26
	functions[0x16a9b] = 109
	functions[0x16b08] = 71
	functions[0x16b4f] = 34
	functions[0x16b71] = 114
	functions[0x16be3] = 71
	functions[0x16c2a] = 34
}
func init() {
	functions[0x15e1e] = 15
	functions[0x15e2d] = 31
	functions[0x17e1f] = 43
	functions[0x17e4a] = 208
	functions[0x17f1a] = 71
	functions[0x17f61] = 32
	functions[0x17f81] = 15
	functions[0x1807b] = 175
	functions[0x1812a] = 83
	functions[0x1817d] = 38
}
func init() {
	functions[0x11c8d] = 36
	functions[0x11f7f] = 249
	functions[0x12151] = 97
	functions[0x121b2] = 53
	functions[0x14325] = 51
	functions[0x14370] = 45
	functions[0x1439d] = 18
	functions[0x143af] = 96
	functions[0x1440f] = 87
	functions[0x14466] = 29
	functions[0x14483] = 22
	functions[0x14499] = 16
	functions[0x144a9] = 45
	functions[0x144d6] = 44
	functions[0x14548] = 45
	functions[0x1463e] = 19
	functions[0x14651] = 56
	functions[0x14689] = 15
	functions[0x159a6] = 17
	functions[0x15ab6] = 27
	functions[0x1703c] = 98
	functions[0x1709e] = 77
	functions[0x17f90] = 75
	functions[0x17fdb] = 115
	functions[0x1804e] = 45
}

func init() {
	functions[0x102f5] = 26
	functions[0x1300e] = 14
	functions[0x1304e] = 67
	functions[0x13830] = 151
	functions[0x13b5a] = 36
	functions[0x13ba9] = 117
	functions[0x13c1e] = 31
	functions[0x13d91] = 56
	functions[0x161ca] = 74
	functions[0x16224] = 55
	functions[0x16288] = 7
	functions[0x1628f] = 108
	functions[0x162fb] = 107
	functions[0x16366] = 89
	functions[0x163bf] = 70
	functions[0x16405] = 112
	functions[0x16475] = 124
	functions[0x164f1] = 134
	functions[0x16577] = 120
	functions[0x165ef] = 22
	functions[0x16605] = 30
	functions[0x16623] = 182
	functions[0x166d9] = 150
	functions[0x1676f] = 30
	functions[0x1678d] = 56
	functions[0x16846] = 125
	functions[0x168c3] = 55
	functions[0x168fa] = 15
	functions[0x16a28] = 21
}

func init() {
	functions[0x10a1c] = 73
	functions[0x11be0] = 173
	functions[0x11cb1] = 24
	functions[0x1533d] = 27
	functions[0x159b7] = 25
	functions[0x189f0] = 46
	functions[0x1ebdc] = 82
}

func init() {
	functions[0x10a1c] = 73
	functions[0x11be0] = 173
	functions[0x11cb1] = 24
	functions[0x1533d] = 27
	functions[0x159b7] = 25
	functions[0x189f0] = 46
	functions[0x1ebdc] = 82
}

func init() {
	functions[0x11be0] = 173
	functions[0x1533d] = 27
	functions[0x189f0] = 46
	functions[0x18a1e] = 179
	functions[0x18ad1] = 25
	functions[0x18aea] = 40
}
func init() {
	for target, size := range map[uint32]int{0x109d0: 76, 0x11e46: 200, 0x11f0e: 34, 0x159b7: 25, 0x159d0: 1, 0x15aa2: 20, 0x15e4c: 20, 0x161b6: 20} {
		functions[target] = size
	}
}

func init() {
	for target, size := range map[uint32]int{0x10cde: 9, 0x12459: 49, 0x1248a: 117, 0x124ff: 52, 0x13efd: 119, 0x13f74: 53, 0x13fa9: 127, 0x14028: 47, 0x14057: 92, 0x140b3: 22, 0x140c9: 140, 0x14155: 63, 0x14194: 162, 0x14269: 66, 0x14575: 76, 0x145c1: 55, 0x1eb11: 77, 0x1eb5e: 14} {
		functions[target] = size
	}
}

func init() {
	for target, size := range map[uint32]int{0x1291a: 93, 0x12977: 76, 0x129c3: 187, 0x12ba8: 49, 0x147bb: 192, 0x1487b: 160, 0x1491b: 244, 0x14a0f: 108, 0x19656: 121, 0x196cf: 30} {
		functions[target] = size
	}
}

func init() {
	for target, size := range map[uint32]int{0x10ce7: 9, 0x14236: 51, 0x1474a: 113, 0x14cf3: 112, 0x14d63: 65, 0x14da4: 76, 0x14df0: 108, 0x14fce: 166, 0x15074: 64, 0x150b4: 35, 0x15130: 131, 0x151b3: 210, 0x15285: 82, 0x152d7: 102, 0x15ca4: 34, 0x15ce0: 57, 0x17028: 20, 0x188cc: 62, 0x1890a: 83, 0x195c9: 118, 0x1963f: 23} {
		functions[target] = size
	}
}

type strategyEvent struct {
	x, y    uint16
	buttons uint8
}
type strategyInput struct {
	maps, menus, positions int
	popup                  bool
	choice                 int
}

var originalStrategy, activeStrategy strategyInput

func outcomeCheckpoint(device *machine.Machine, d *dos.DOS, t uint16) {
	state := &originalStrategy
	if device == activeDevice {
		state = &activeStrategy
	}
	if t == 0x21e7 {
		state.popup = false
	}
	if t == 0x1f7f && len(currentInput.mapEvents) != 0 {
		event := strategyEvent{currentInput.mx, currentInput.my, 2}
		if state.maps < len(currentInput.mapEvents) {
			event = currentInput.mapEvents[state.maps]
		}
		state.maps++
		state.popup = false
		d.Mouse.ReleaseButton(0)
		d.Mouse.ReleaseButton(1)
		d.Mouse.X = event.x
		d.Mouse.Y = event.y
		if event.buttons&1 != 0 {
			d.Mouse.PressButton(0)
		}
		if event.buttons&2 != 0 {
			d.Mouse.PressButton(1)
		}
	}
	if t == 0x93e9 {
		state.popup = true
		state.positions = 0
		state.choice = -1
		if state.menus < len(currentInput.menuChoices) {
			state.choice = currentInput.menuChoices[state.menus]
		}
		state.menus++
		d.Mouse.ReleaseButton(0)
		d.Mouse.ReleaseButton(1)
	}
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
	originalDOSDevice = originalDOS
	activeDOSDevice = activeDOS
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
	probeRaw, err := os.ReadFile("/repo/workplace/matching-decompilation/c-outcome/ida/ida-probe.json")
	if err != nil {
		panic(err)
	}
	var probe struct {
		RecoveryTargets   []uint32 `json:"recovery_targets"`
		NavigationTargets []uint32 `json:"navigation_only_targets"`
		Targets           []struct {
			IDALinear uint32 `json:"ida_linear"`
		} `json:"targets"`
	}
	if err = json.Unmarshal(probeRaw, &probe); err != nil {
		panic(err)
	}
	if len(probe.RecoveryTargets) != 21 {
		panic("bootstrap recovery probe target count")
	}
	newRoutineHashes := map[string]string{}
	for _, target := range probe.RecoveryTargets {
		name := entryName(target)
		hash, exists := routineHashes[name]
		if !exists {
			panic("main recovery target absent from runtime identity checks")
		}
		newRoutineHashes[name] = hash
	}

	navigationNames := []string{}
	for _, target := range probe.NavigationTargets {
		navigationNames = append(navigationNames, entryName(target))
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
	if raw[header+0x67] != 0xE8 || uint16(0x6A+int16(binary.LittleEndian.Uint16(raw[header+0x68:]))) != 0x1BE0 {
		panic("real outer CALL identity")
	}
	kyo, err := os.ReadFile("/orig/KYOGRF.DAT")
	if err != nil || len(kyo) != 69120 || fmt.Sprintf("%x", sha256.Sum256(kyo)) != "e086f526bdada5baf751c41d2f73a78a0ba70002f282f63a9f33114542ed933f" {
		panic("KYOGRF identity")
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
	var originalAPI, cAPI, originalPortReads []byte
	originalFontMisses, cFontMisses := map[string]int{}, map[string]int{}
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
				input := &originalStrategy
				if device == activeDevice {
					input = &activeStrategy
				}
				if input.popup {
					input.positions++
					p := input.positions - 2
					if p >= 0 {
						if p < input.choice {
							d.Mouse.X = 0
							d.Mouse.Y = 132
						} else {
							button := 0
							if input.choice < 0 {
								button = 1
							}
							d.Mouse.PressButton(button)
						}
					}
				}

				if len(currentInput.drag) > 0 {
					index, on := &activeDrags, &activeDragOn
					if device == original {
						index, on = &originalDrags, &originalDragOn
					}
					if *on {
						if *index < len(currentInput.drag) {
							d.Mouse.Y = currentInput.drag[*index]
							(*index)++
						} else {
							d.Mouse.ReleaseButton(0)
							*on = false
						}
					}
				}
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
			fontBefore, fontCode := d.Font.Missing, c.R[cpu.CX]
			soundStatus := n == 0x61 && byte(c.R[cpu.AX]>>8) == 0x0A
			result := previous(c, n)
			if soundStatus {
				c.R[cpu.AX] = c.R[cpu.AX]&0xFF00 | uint16(currentInput.soundFlags)
			}
			if d.Font.Missing > fontBefore {
				misses := cFontMisses
				if device == original {
					misses = originalFontMisses
				}
				misses[fmt.Sprintf("%s:%02X:%04X", currentInput.group, n, fontCode)]++
			}
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
		{0xd38, 0x4500}, {0xd40, 0x8000}, {0xd44, 0x5000}, {0xd52, 0x7000}, {0xd56, 0xd000}, {0x9872, 0xb000}, {0xd84a, 0x3000}, {0xd84c, 0x3800}, {0xd84e, 0x2500}, {0xd850, 0x5000},
		{0x987c, 0x9000}, {0x9876, 0xe000}, {0x987a, 0x4200}, {0xe479, 0x7600}, {0xe47b, 0}, {0xd30a, 0x7000}, {0xd30e, 0x5200}} {
		putWord(initialRAM, cs, v[0], v[1])
	}
	for _, v := range [][2]uint16{{0xf78a, 0x410}, {0xf78c, 0x80}, {0xf791, 0x414}, {0xf793, 0x80}} {
		putWord(initialRAM, cs, v[0], v[1])
	}
	// Original sub_1030F registration: code-resident list/text and ICONGRF + 0x9A0.
	for _, v := range [][2]uint16{{0xd50, 0x229a}, {0xeae9, 0xe16}, {0xeaeb, cs}, {0xeaed, cs}, {0xeaef, 0x229a}, {0xeaf1, 0x1800}, {0xeaf3, 0x229a}, {0xeaf5, 0x2020}} {
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

	wholeFixture := C.KiOutcomeFixture{}
	fixture := &wholeFixture.list
	originalBase, activeBase := original.Snapshot(), active.Snapshot()
	groups, entries := map[string]int{}, map[string]int{}
	mousePollCounts, waitCheckpoints := map[string]int{}, map[string]int{}
	sampleCounts := map[uint32]int{}
	cases, audits, pixelAudits, independentAudits := 0, 0, 0, 0
	stageAudits := map[string]int{}
	segmentFixtures := map[string]int{}
	rngInputs := map[string]string{}
	nonlocalCases, normalCases := 0, 0
	originalDigest, cDigest := sha256.New(), sha256.New()
	var mismatch map[string]any
	resetCounters := func(c resumeCase) {
		originalPolls, activePolls, originalDrags, activeDrags = 0, 0, 0, 0
		originalDragOn, activeDragOn = false, false
		originalPositions, activePositions = 0, 0
		originalStrategy, activeStrategy = strategyInput{}, strategyInput{}
		currentInput = c
		activeTicks, originalTicks, originalQueries, activeQueries = 0, 0, 0, 0
		cPortReads, originalPortReads, originalAPI, cAPI = nil, nil, nil, nil
		for _, device := range []*machine.Machine{original, active} {
			device.PortLog = nil
			device.PortsIn = map[uint16]uint64{}
		}
		for _, d := range []*dos.DOS{originalDOS, activeDOS} {
			d.Mouse = dos.Mouse{X: c.mx, Y: c.my, MaxX: 6143, MaxY: 4127, Press: c.presses, Calls: map[uint16]int{}}
			for i := 0; i < 3; i++ {
				d.Mouse.PressAt[i] = [2]uint16{c.mx, c.my}
			}
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
	resetRaw := func(c resumeCase) {
		original.Restore(originalBase)
		active.Restore(activeBase)
		for _, device := range []*machine.Machine{original, active} {
			copy(device.Mem, initialRAM)
			copy(device.Mem[0x22000:], icons[0x9700:])
			device.SetVideoMode(0x12)
			for plane := 0; plane < 4; plane++ {
				copy(device.VGA.Planes[plane][:], planes[3][plane*65536:(plane+1)*65536])
			}
			device.Out8(0x3c4, 2)
			device.Out8(0x3c5, 15)
			for _, q := range [][2]byte{{0, 0}, {1, 0}, {3, 0}, {4, 0}, {5, 0}, {8, 255}} {
				device.Out8(0x3ce, q[0])
				device.Out8(0x3cf, q[1])
			}
			device.Read8(0xa4321)
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
			outcomePrepare(device.Mem, cs, c)
			if c.target == 0x1963f {
				// Original caller sub_195C9 selects GC read-map index4 before this helper.
				device.Out8(0x3ce, 4)
			}
			for i := range device.DAC {
				device.DAC[i] = byte(i*13+int(c.seed)*7) & 63
			}
		}
		if &original.Mem[0] == &active.Mem[0] {
			panic("bootstrap shared RAM")
		}
		for plane := 0; plane < 4; plane++ {
			if &original.VGA.Planes[plane][0] == &active.VGA.Planes[plane][0] {
				panic("bootstrap shared VGA plane")
			}
		}
		if !bytes.Equal(original.Mem, active.Mem) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(paletteBytes(original), paletteBytes(active)) || !bytes.Equal(portMapBytes(original), portMapBytes(active)) {
			panic("independent raw fixture setup mismatch")
		}
	}
	check := func(c resumeCase) bool {
		if *only != "" && c.group != *only {
			return true
		}
		if *smoke && sampleCounts[c.target] >= 16 {
			return true
		}
		sampleCounts[c.target]++
		cases++
		groups[c.group]++
		resetRaw(c)
		segmentKey := fmt.Sprintf("D44=%04X/D850=%04X/occupancy=%04X", strategyWord(original.Mem, int(cs)*16+0xd44), strategyWord(original.Mem, int(cs)*16+0xd850), strategyWord(original.Mem, int(cs)*16+0x9872))
		segmentFixtures[segmentKey]++
		sc := c
		phase := "outcome"
		{
			resetCounters(sc)
			before := append([]byte(nil), original.Mem...)
			rngInputs[fmt.Sprintf("%02X", sc.outcome.rngSeed)] = fmt.Sprintf("%x", sha256.Sum256(before[int(cs)*16+0xecfc:int(cs)*16+0xecfc+258]))
			mainBeforeDAC = original.DAC
			r := registers{AX: sc.ax, BX: sc.bx, CX: sc.cx, DX: sc.dx, SI: sc.si, DI: sc.di, BP: sc.bp, SP: sc.entrySP, SS: sc.entrySS, DS: 0x7000, ES: 0x5000, CS: cs, IP: uint16(sc.target), Flags: uint16(2 | sc.seed&1)}
			switch sc.target {
			case 0x15ca4:
				r.DS = cs
			case 0x1963f:
				r.DS = 0xa0c8
				r.ES = 0x6800
			}
			if sc.df {
				r.Flags |= 0x400
			}
			if sc.altDS {
				r.DS = 0x2600
			}
			for plane := 0; plane < 4; plane++ {
				copy(routeBeforePlanes[plane][:], original.VGA.Planes[plane][:])
				copy(outcomeBeforePlanes[plane][:], original.VGA.Planes[plane][:])
			}
			original.CPU.R = [8]uint16{r.AX, r.CX, r.DX, r.BX, r.SP, r.BP, r.SI, r.DI}
			original.CPU.Seg = [4]uint16{r.ES, r.CS, r.SS, r.DS}
			original.CPU.IP = r.IP
			original.CPU.SetFlags(r.Flags)
			putWord(original.Mem, r.SS, r.SP, 0xf000)
			putWord(cMemory, r.SS, r.SP, 0xf000)
			assignC(&m, originalRegs(original))
			var originalTrace []byte
			var observed []registers
			repeats := sc.repeat
			if repeats < 1 {
				repeats = 1
			}
			for iteration := 0; iteration < repeats; iteration++ {
				if iteration != 0 {
					original.CPU.IP = uint16(sc.target)
					original.CPU.R[cpu.SP] = sc.entrySP
				}
				steps := 0
				for ; steps < 6000000; steps++ {
					cur := originalRegs(original)
					if cur.CS == cs && (cur.IP == 0xf000 || cur.IP == sc.returnIP) {
						break
					}
					loc := uint32(cur.IP) + 0x10000
					if cur.CS == mouseCS {
						loc = uint32(cur.IP) + 0x20000
					}
					if cur.CS == cs {
						outcomeCheckpoint(original, originalDOS, cur.IP)
					}
					if cur.CS == cs && (cur.IP == 0x21e7 || cur.IP == 0x84dd) {
						listCheckpoint(original, originalDOS, cur.IP)
					}
					if cur.CS == cs && cur.IP == 0x2259 {
						originalTicks++
						value := byte(255)
						if originalTicks <= sc.timerHold {
							value = 254
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
				if steps >= 6000000 {
					panic(fmt.Sprintf("bounded bootstrap original group=%s phase=%s input=%+v state=%+v", c.group, phase, sc, originalRegs(original)))
				}
			}
			routeOriginalPortReads = append(routeOriginalPortReads[:0], originalPortReads...)
			outcomeOriginalPortReads = append(outcomeOriginalPortReads[:0], originalPortReads...)
			independentAudits += outcomeAudit(before, original.Mem, cs, sc, observed, original)
			for iteration := 0; iteration < repeats; iteration++ {
				if iteration != 0 {
					m.ip = C.uint16_t(sc.target)
					m.sp = C.uint16_t(sc.entrySP)
				}
				C.outcome_run(&m, C.uint16_t(sc.target), &wholeFixture)
			}
			if fixture.calls > 32768 || fixture.unsupported != 0 {
				panic(fmt.Sprintf("bootstrap unsupported=%x calls=%d group=%s phase=%s", fixture.unsupported, fixture.calls, c.group, phase))
			}
			var cTrace []byte
			for i := 0; i < int(fixture.calls); i++ {
				snap := fixture.blocks[i/8192].trace[i%8192]
				loc := uint32(snap.target) + 0x10000
				if uint16(snap.regs[11]) == mouseCS {
					loc = uint32(snap.target) + 0x20000
				}
				_, known := functions[loc]
				if !known && !(uint16(snap.regs[11]) == cs && snap.target == 0x2259) && !(snap.regs[11] == machine.StubSeg && (snap.target == 0x410 || snap.target == 0x414)) {
					continue
				}
				cTrace = append(cTrace, byte(snap.target), byte(snap.target>>8))
				for _, v := range snap.regs {
					cTrace = append(cTrace, byte(v), byte(v>>8))
				}
				cTrace = append(cTrace, C.GoBytes(unsafe.Pointer(&snap.stack[0]), 32)...)
				cTrace = append(cTrace, C.GoBytes(unsafe.Pointer(&snap.device[0]), 28)...)
			}
			transferOK := wholeFixture.transferred == 0
			if sc.outcome.escape {
				transferOK = wholeFixture.transferred == 1 && uint16(wholeFixture.transfer_ip) == sc.returnIP && uint16(wholeFixture.transfer_ss) == sc.savedSS && uint16(wholeFixture.transfer_sp) == sc.savedSP+2
			}
			want, got := originalRegs(original), cRegs(&m)
			originalPorts, cPorts := portBytes(original), portBytes(active)
			originalState, cState := deviceBytes(original), deviceBytes(active)
			if !transferOK || want != got || !bytes.Equal(original.Mem, cMemory) || !bytes.Equal(original.VGA.Raw(), active.VGA.Raw()) || !bytes.Equal(originalState, cState) || !bytes.Equal(originalPorts, cPorts) || !bytes.Equal(originalTrace, cTrace) || !bytes.Equal(originalAPI, cAPI) || !bytes.Equal(originalPortReads, cPortReads) || !bytes.Equal(mouseBytes(originalDOS), mouseBytes(activeDOS)) || !bytes.Equal(soundBytes(originalDOS), soundBytes(activeDOS)) || !bytes.Equal(paletteBytes(original), paletteBytes(active)) || !bytes.Equal(portMapBytes(original), portMapBytes(active)) || originalTicks != activeTicks {
				mismatch = map[string]any{"group": c.group, "stage": phase, "case": cases - 1, "input": r, "original": want, "c": got, "expected_escape": sc.outcome.escape, "native_transferred": int(wholeFixture.transferred), "native_transfer_ip": uint16(wholeFixture.transfer_ip), "native_transfer_ss": uint16(wholeFixture.transfer_ss), "native_transfer_sp": uint16(wholeFixture.transfer_sp), "original_ram": fmt.Sprintf("%x", sha256.Sum256(original.Mem)), "c_ram": fmt.Sprintf("%x", sha256.Sum256(cMemory)), "original_planes": fmt.Sprintf("%x", sha256.Sum256(original.VGA.Raw())), "c_planes": fmt.Sprintf("%x", sha256.Sum256(active.VGA.Raw())), "original_trace": hex.EncodeToString(originalTrace), "c_trace": hex.EncodeToString(cTrace), "original_ports": hex.EncodeToString(originalPorts), "c_ports": hex.EncodeToString(cPorts), "original_in": hex.EncodeToString(originalPortReads), "c_in": hex.EncodeToString(cPortReads), "original_api": hex.EncodeToString(originalAPI), "c_api": hex.EncodeToString(cAPI), "original_palette_dac": hex.EncodeToString(paletteBytes(original)), "c_palette_dac": hex.EncodeToString(paletteBytes(active))}
				return false
			}
			audits++
			if wholeFixture.transferred == 1 {
				nonlocalCases++
			} else {
				normalCases++
			}
			pixelAudits++
			stageAudits[phase]++
			mousePollCounts[c.group] += originalQueries
			waitCheckpoints[c.group] += originalTicks
			if !bytes.Equal(indexed(original), indexed(active)) {
				panic("bootstrap indexed pixels")
			}
			for _, b := range [][]byte{regBytes(want), original.Mem, original.VGA.Raw(), originalState, originalPorts, originalTrace, originalAPI, originalPortReads, mouseBytes(originalDOS), soundBytes(originalDOS), paletteBytes(original), portMapBytes(original)} {
				originalDigest.Write(b)
			}
			for _, b := range [][]byte{regBytes(got), cMemory, active.VGA.Raw(), cState, cPorts, cTrace, cAPI, cPortReads, mouseBytes(activeDOS), soundBytes(activeDOS), paletteBytes(active), portMapBytes(active)} {
				cDigest.Write(b)
			}
			return true
		}
	}

	base := func(group string, target uint32, scenario int) resumeCase {
		c := resumeCase{group: group, target: target, scenario: scenario, repeat: 1, player: 0, entrySS: 0xf000, entrySP: 0x6ffe, returnIP: 0xf000, ax: 0x5500, bx: 0x1234, cx: 0x56, dx: 15, bp: 0x5000, di: 0x5000, si: 0x4240, seed: 3, mask: 15, acceptButton: -1, syntheticSize: -1, mx: 100, my: 120, oldX: 120, oldY: 80, patch: 0x0feb, expectedRow: -1}
		world := worlds[scenario*22208+0x80 : scenario*22208+0x52c0]
		c.owner, c.general = strategyForeign(world), strategyGeneral(world, c.player)
		if c.owner < 0 || c.general < 0 {
			panic("interaction original foreign faction/general")
		}
		copy(c.palette[:], gamePalette[48:96])
		return c
	}
	all := outcomeCases(base)

	for _, c := range all {
		if !check(c) {
			break
		}
	}
	report := map[string]any{
		"schema": "wolong-c-outcome-parity-v1", "input_sha256": inputHash, "icon_sha256": assetHash,
		"routine_sha256": routineHashes, "new_routine_sha256": newRoutineHashes, "new_code_block_sha256": map[string]string{}, "new_probe_sha256": fmt.Sprintf("%x", sha256.Sum256(probeRaw)),
		"source_manifest": map[string]string{"path": "workplace/matching-decompilation/c-outcome/results/c-source.sha256", "sha256": manifestHash},
		"cases":           cases, "groups": groups, "entries_seen": entries, "full_ram_plane_audits": audits, "indexed_content_audits": pixelAudits, "palette_dac_audits": audits, "independent_outcome_audits": independentAudits, "stage_audits": stageAudits, "segment_fixtures": segmentFixtures,
		"original_state_sha256": hex.EncodeToString(originalDigest.Sum(nil)), "c_state_sha256": hex.EncodeToString(cDigest.Sum(nil)),
		"original_font_misses": originalFontMisses, "c_font_misses": cFontMisses, "original_missing_fonts": originalDOS.Font.Missing, "c_missing_fonts": activeDOS.Font.Missing, "original_font_calls": originalDOS.Font.Calls, "c_font_calls": activeDOS.Font.Calls, "mouse_queries": mousePollCounts, "counter_checkpoints": waitCheckpoints,
		"scope":               "Twenty-one original automatic-combat, retreat and city-transfer functions, including original faction-loss nonlocal exit. Controlled raw world/RNG and independent graphics models; full field/siege callers, tactical/army/main loops and C machine-code matching remain incomplete.",
		"platform":            "Independent pinned DOS/font/VGA services; full RAM/planes/DAC/registers/FLAGS/stack/IN/OUT/API; mature 3DA IN contract; no wall-clock claim",
		"game_palette_sha256": fmt.Sprintf("%x", sha256.Sum256(gamePalette)), "game_palette_size": len(gamePalette), "additional_assets": identities,
		"independent_raw_initialization": true, "cross_machine_snapshot_initialization": false, "fixed_raw_rng_state": true, "actual_nonlocal_exit": true,
		"nonlocal_exit_cases": nonlocalCases, "normal_return_cases": normalCases,
		"rng_initial_state_sha256": rngInputs,
		"rng_initialization":       "Fixed synthetic raw 258-byte RNG state before execution: counter=seed^A5, sample=seed, table[i]=byte(i*73+i/2+seed). Neither side is rerolled; native C consumes its own state.",
		"passed":                   mismatch == nil, "mismatch": mismatch, "c_machine_code_match": false,
	}

	b, _ := json.MarshalIndent(report, "", "  ")
	b = append(b, '\n')
	if e := os.WriteFile(*output, b, 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Outcome original/C %d; complete RAM/plane %d; content %d; pass=%v\n", cases, audits, pixelAudits, mismatch == nil)
	if mismatch != nil {
		os.Exit(1)
	}
}
