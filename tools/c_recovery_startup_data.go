//go:build matching_startup

package main

import (
	"fmt"
	"sort"

	"github.com/wicanr2/dosgolem/internal/dos"
	"github.com/wicanr2/dosgolem/internal/machine"
)

func startupSeedName(mem []byte, c startupCase) {
	base := 0x70000 + 0x5221
	switch c.namePat {
	case 0:
		mem[base] = 0xFF
		for i := 1; i < 13; i++ {
			mem[base+i] = 0
		}
	case 1:
		mem[base] = 5
		for i := 1; i < 13; i++ {
			mem[base+i] = byte(0xA0 + i)
		}
	case 2:
		mem[base] = 0x91
		for i := 1; i < 13; i++ {
			mem[base+i] = 0x40
		}
	case 3:
		mem[base] = 0
		for i := 1; i < 13; i++ {
			mem[base+i] = byte(i)
		}
	}
}

var clickCode = map[int]startupEvent{
	0x20: {468, 196, 0}, 0x21: {356, 196, 0}, 0x22: {412, 196, 0}, 0x23: {204, 364, 0},
	0x24: {372, 364, 0}, 0x25: {220, 220, 0}, 0x26: {204, 236, 0}, 0x27: {276, 148, 0},
	0x28: {276, 180, 0},
}

func pressTape(x, y uint16, n int) []startupEvent {
	out := make([]startupEvent, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, startupEvent{x, y, 0})
	}
	return out
}

func rightTape(x, y uint16, n int) []startupEvent {
	out := make([]startupEvent, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, startupEvent{x, y, 1})
	}
	return out
}

// gridTape is a SWEEP (temporary): one press per grid point to discover
// which screen positions yield 18B7C slot codes. Revert after mapping.
func gridTape() []startupEvent {
	var out []startupEvent
	for _, y := range []uint16{60, 100, 140, 180, 220, 260, 320, 360} {
		for _, x := range []uint16{80, 160, 240, 320, 400, 480, 560, 600} {
			out = append(out, startupEvent{x, y, 0})
		}
	}
	return out
}

// rowTape is a MICRO-SWEEP (temporary): probe the 18BD1-painted slot rows
// (BX=0x90/0xC0/0xF0/0x120) for codes 0x20-0x23. Revert after mapping.
func rowTape() []startupEvent {
	var out []startupEvent
	for _, y := range []uint16{144, 148, 152, 192, 196, 200, 240, 244, 248, 288, 292, 296} {
		out = append(out, startupEvent{320, y, 0})
	}
	return out
}

func hoverTape(x, y uint16, n int) []startupEvent {
	out := make([]startupEvent, 0, n+1)
	for i := 0; i < n; i++ {
		out = append(out, startupEvent{x, y, -1})
	}
	return out
}

func clickTape(x, y uint16, hovers int) []startupEvent {
	_ = hovers
	return []startupEvent{{x, y, 0}}
}

func startupCases(base func(string, uint32, int) startupCase, nameFont []byte) []startupCase {
	_ = nameFont
	var all []startupCase
	add := func(c startupCase) { all = append(all, c) }

	for _, psp := range []byte{0, 9, 127} {
		c := base("entry", 0x10000, 0)
		c.mouseAbsent = true
		c.pspCmd = psp
		c.expectExit = true
		c.int61AH = nil
		add(c)
	}
	{
		c := base("boot", 0x10000, 0)
		c.int61AH = []uint8{0, 2, 3, 0x0C}
		c.pauseIters = 1
		c.pspCmd = 0
		c.events = append(clickTape(300, 175, 1), pressTape(320, 148, 1)...)
		c.events = append(c.events, pressTape(400, 252, 1)...)
		c.events = append(c.events, pressTape(100, 100, 1)...)
		c.events = append(c.events, pressTape(350, 256, 1)...)
		c.events = append(c.events, pressTape(495, 200, 1)...)
		c.events = append(c.events, pressTape(350, 280, 1)...)
		c.events = append(c.events, pressTape(100, 100, 400)...)
		add(c)
	}
	{
		c := base("boot", 0x10000, 1)
		c.int61AH = []uint8{0, 2, 3, 0x0C}
		c.pauseIters = 1
		c.pspCmd = 12
		c.events = append(clickTape(300, 210, 1), rightTape(320, 148, 1)...)
		c.events = append(c.events, pressTape(300, 175, 1)...)
		c.events = append(c.events, pressTape(320, 148, 1)...)
		c.events = append(c.events, pressTape(400, 252, 1)...)
		c.events = append(c.events, pressTape(100, 100, 1)...)
		c.events = append(c.events, pressTape(350, 256, 1)...)
		c.events = append(c.events, pressTape(495, 200, 1)...)
		c.events = append(c.events, pressTape(350, 280, 1)...)
		c.events = append(c.events, pressTape(100, 100, 400)...)
		add(c)
	}

	for _, row := range []struct {
		y     uint16
		pads  int
		right bool
	}{{196, 1, false}, {244, 4, true}} {
		c := base("boot", 0x10000, 0)
		c.int61AH = []uint8{0, 2, 3, 0x0C}
		c.pauseIters = 1
		c.pspCmd = 0
		c.events = append(clickTape(300, 175, 1), pressTape(320, row.y, 1)...)
		if row.right {
			c.events = append(c.events, pressTape(160, 140, 1)...)
		} else {
			c.events = append(c.events, pressTape(400, 252, 1)...)
		}
		c.events = append(c.events, pressTape(100, 100, row.pads)...)
		c.events = append(c.events, pressTape(350, 256, 1)...)
		c.events = append(c.events, pressTape(495, 200, 1)...)
		c.events = append(c.events, pressTape(350, 280, 1)...)
		c.events = append(c.events, pressTape(100, 100, 400)...)
		add(c)
	}
	for i, seed := range []byte{0, 3, 7} {
		c := base("bootchain", 0x1006B, i%4)
		c.seed = seed
		c.df = i == 2
		add(c)
	}
	{
		c := base("bootchain", 0x1006B, 2)
		c.videoMode = 3
		add(c)
	}
	for i := 0; i < 4; i++ {
		c := base("bootchain", 0x10155, i)
		c.mouseAbsent = i == 3
		add(c)
	}
	for i, ax := range []uint16{0x0000, 0x00FF, 0x5500, 0xFF00} {
		c := base("bootchain", 0x10210, i)
		c.ax = ax
		c.int61AH = []uint8{0, 2, 3}
		add(c)
	}
	for i := 0; i < 3; i++ {
		c := base("bootchain", 0x10236, i)
		c.int61AH = []uint8{8, 4}
		add(c)
	}
	for i := 0; i < 2; i++ {
		c := base("bootchain", 0x1033B, i)
		c.int61AH = []uint8{0x0C}
		add(c)
	}
	for i := 0; i < 2; i++ {
		c := base("bootchain", 0x1034E, i)
		c.int61AH = []uint8{0x0C}
		add(c)
	}
	for i, v := range []struct {
		seg, count, src, dst uint16
		df                   bool
	}{{0x7000, 16, 0x100, 0x200, false}, {0x7000, 1, 0, 0, false}, {0x7000, 64, 0x200, 0x100, false}, {0x7000, 32, 0x100, 0x104, false}, {0x7000, 32, 0x104, 0x100, true}, {0x8000, 8, 0x10, 0x20, true}, {0x7000, 0, 0x100, 0x200, false}, {0x4200, 4, 0, 64, false}} {
		c := base("bootchain", 0x10612, i%4)
		c.cx, c.ax, c.bx, c.dx = v.seg, v.count, v.src, v.dst
		c.df = v.df
		add(c)
	}

	for i := 0; i < 3; i++ {
		c := base("dispatcher", 0x11AAE, i)
		add(c)
	}
	for i := 0; i < 3; i++ {
		c := base("boot", 0x11AC3, i)
		c.int61AH = []uint8{}
		c.bx = uint16(i) * 64
		c.events = append(clickTape(310, 175, 1), pressTape(320, 148, 1)...)
		c.events = append(c.events, pressTape(400, 252, 1)...)
		c.events = append(c.events, pressTape(100, 100, 1)...)
		c.events = append(c.events, pressTape(350, 256, 1)...)
		c.events = append(c.events, pressTape(495, 200, 1)...)
		c.events = append(c.events, pressTape(350, 280, 1)...)
		c.events = append(c.events, pressTape(100, 100, 400)...)
		add(c)
	}
	for i := 0; i < 2; i++ {
		c := base("boot", 0x11A6E, i)
		c.int61AH = []uint8{}
		c.events = append(clickTape(310, 175, 1), pressTape(320, 148, 1)...)
		c.events = append(c.events, pressTape(400, 252, 1)...)
		c.events = append(c.events, pressTape(100, 100, 1)...)
		c.events = append(c.events, pressTape(350, 256, 1)...)
		c.events = append(c.events, pressTape(495, 200, 1)...)
		c.events = append(c.events, pressTape(350, 280, 1)...)
		c.events = append(c.events, pressTape(100, 100, 400)...)
		add(c)
	}

	for i, v := range []struct {
		iters    int
		a3, a5   byte
		d2d      byte
		scenario int
	}{{1, 0, 0, 0, 0}, {1, 0x80, 0, 0, 1}, {2, 0x80, 1, 0x50, 2}, {2, 0xFF, 2, 0xFF, 3}, {3, 0x80, 0, 0, 0}, {1, 0x7F, 3, 0xA0, 1}} {
		c := base("mainloop", 0x11BE0, v.scenario)
		c.pauseIters = v.iters
		c.a3, c.a5, c.d2d = v.a3, v.a5, v.d2d
		c.events = pressTape(320, 200, 20)
		_ = i
		add(c)
	}
	for _, a5 := range []byte{0, 1, 2, 3} {
		// d2d=0/1 hang in 11D8E's deadline spin ([D2D] < [0CFA]=2 waits for
		// ticks that never fire in-matrix); 2 is the exit boundary.
		// a5=1 with d2d>=0xA0 zeroes D2C then calls 11D8E, which spins
		// while D2C==0 (ticks re-set it on hardware); sub-0xA0 only there.
		d2ds := []byte{2, 0x9F, 0xA0, 0xFF}
		if a5 == 1 {
			d2ds = []byte{2, 3, 0x9E, 0x9F}
		}
		for _, d2d := range d2ds {
			c := base("mainloop", 0x11CD0, int(a5))
			c.a5, c.d2d, c.d2a = a5, d2d, 0
			c.d2c = 1
			add(c)
		}
	}
	for _, scen := range []int{0, 3} {
		c := base("transfer", 0x11CD0, scen)
		c.a5, c.d2d, c.d2a = 0, 0xFF, 1
		c.d2c = 1
		c.transfer = true
		c.transferIP = 0xF000
		c.anchorSS, c.anchorSP = 0xF000, 0x6FFE
		c.pollFeed = true
		c.events = append(rightTape(320, 200, 2), pressTape(320, 200, 60)...)
		add(c)
	}
	for i, scen := range []int{1, 2} {
		c := base("transfer", 0x11CD0, scen)
		c.a5 = 1
		c.d2a = 1
		c.d2d = 0xFF
		c.transfer = true
		c.transferIP = 0xF000
		c.anchorSS, c.anchorSP = 0xF000, 0x6FFE
		c.pollFeed = true
		c.events = append(rightTape(320, 200, 2), pressTape(320, 200, 60)...)
		_ = i
		add(c)
	}

	for i := 0; i < 3; i++ {
		c := base("load", 0x18B12, i)
		c.bx = uint16(i) * 64
		c.events = pressTape(320, 200, 60)
		add(c)
	}
	for i := 0; i < 3; i++ {
		c := base("load", 0x18B40, i)
		c.bx = uint16(i) * 64
		// Shipped SAVE.DAT has no saves (slots read D0A1 blank), so
		// left-clicks wait again; the right-click cancels via 121E7 CF=1.
		c.events = append(pressTape(320, 200, 3), rightTape(320, 200, 57)...)
		add(c)
	}
	for _, v := range [][2]byte{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 1}, {1, 2}, {0, 3}, {2, 0}} {
		c := base("load", 0x18CAE, int(v[0]+v[1])%4)
		c.ax = uint16(v[0])<<8 | uint16(v[1])
		c.events = pressTape(320, 200, 60)
		add(c)
	}

	for i, clicks := range [][][2]uint16{
		{{300, 175}}, {{300, 210}}, {{248, 168}}, {{375, 183}},
		{{248, 192}}, {{375, 239}}, {{300, 188}, {300, 175}}, {{200, 175}, {300, 175}},
		{{300, 100}, {300, 175}}, {{300, 210}, {300, 210}},
	} {
		c := base("prompt", 0x18DC8, i%4)
		c.si, c.dx, c.bx, c.di = 0x98C8, 0xD0, 0x80, 0x28
		for _, click := range clicks {
			c.events = append(c.events, clickTape(click[0], click[1], 20)...)
		}
		add(c)
	}

	for i := 0; i < 3; i++ {
		c := base("naming-standalone", 0x18E5A, i)
		c.bx = uint16(i) * 64
		c.si = uint16(i) * 2
		c.events = append(clickTape(332, 276, 20), clickTape(332, 252, 5)...)
		add(c)
	}
	for i := 0; i < 4; i++ {
		c := base("naming-standalone", 0x18EA0, i)
		c.ds = 0x7000
		c.si = uint16(i) * 2
		c.namePat = i
		add(c)
	}
	for i := 0; i < 2; i++ {
		c := base("naming-standalone", 0x18F37, i)
		add(c)
	}
	for i := 0; i < 2; i++ {
		c := base("naming-standalone", 0x18F6D, i)
		add(c)
	}
	for i := 0; i < 4; i++ {
		c := base("naming-standalone", 0x18F7C, i)
		c.ds = 0x7000
		c.si = uint16(i)*64 + 0x100
		add(c)
	}
	for i, si := range []uint16{0, 2, 4, 8, 10} {
		c := base("naming-standalone", 0x1905E, i%4)
		c.ds, c.es = 0x7000, 0x9000
		c.si = si
		c.prealloc = true
		add(c)
	}
	for i, si := range []uint16{0, 2, 8, 10} {
		c := base("naming-standalone", 0x1906E, i%4)
		c.ds, c.es = 0x7000, 0x9000
		c.si = si
		c.prealloc = true
		add(c)
	}
	for i, di := range []uint16{0, 0x40, 0xC0, 0x13BA} {
		c := base("naming-standalone", 0x1907F, i)
		c.ds, c.es = 0x7000, 0x9000
		c.di = di
		c.prealloc = true
		add(c)
	}
	for i, di := range []uint16{0, 0x1000, 0x12FA, 0x13BA} {
		c := base("naming-standalone", 0x1908D, i)
		c.ds, c.es = 0x7000, 0x9000
		c.di = di
		c.prealloc = true
		add(c)
	}
	for i, cx := range []uint16{0, 0xD7, 0xD8, 0x100, 0x27F, 0x1A7 + 0xD8} {
		c := base("naming-standalone", 0x190A3, i%4)
		c.ds, c.es = 0x7000, 0x9000
		c.cx = cx
		c.prealloc = true
		add(c)
	}
	for i, v := range [][3]uint16{{0xD2, 0xEE, 0}, {0xE0, 0x100, 0x200}, {0x100, 0x120, 0x1478}, {0x13F, 0x165, 0x13BA}, {0, 0, 0}, {0xD1, 0xEE, 0}, {0x140, 0xEE, 0}, {0xE0, 0xED, 0}, {0xE0, 0x166, 0}, {0xD2 + 19, 0xEE + 19, 0x100}} {
		c := base("naming-standalone", 0x190C0, i%4)
		c.ds, c.es = 0x7000, 0x9000
		c.cx, c.dx, c.di = v[0], v[1], v[2]
		c.si = 4
		c.prealloc = true
		add(c)
	}
	for i, pat := range []int{0, 1, 2, 3, 1} {
		c := base("naming-standalone", 0x1912D, i%4)
		c.ds, c.es = 0x7000, 0x9000
		c.namePat = pat
		c.prealloc = true
		add(c)
	}
	for i, pat := range []int{0, 1, 2, 3, 2} {
		c := base("naming-standalone", 0x19144, i%4)
		c.ds, c.es = 0x7000, 0x9000
		c.namePat = pat
		c.prealloc = true
		add(c)
	}
	for i := 0; i < 4; i++ {
		c := base("naming-standalone", 0x1915B, i)
		c.ds = 0x7000
		c.si = uint16(i) * 2
		c.namePat = i
		add(c)
	}
	for i := 0; i < 2; i++ {
		c := base("naming-standalone", 0x1917C, i)
		add(c)
	}
	for i := 0; i < 2; i++ {
		c := base("naming-standalone", 0x19214, i)
		add(c)
	}
	for i, si := range []uint16{0, 4, 8, 10} {
		c := base("naming-standalone", 0x19223, i)
		c.ds, c.es = 0x7000, 0x9000
		c.si = si
		c.namePat = i
		c.prealloc = true
		add(c)
	}
	for i, di := range []uint16{0, 0x200, 0x13BA, 0x1478, 0x100, 0x800} {
		c := base("naming-standalone", 0x1928A, i%4)
		c.ds, c.es = 0x7000, 0x9000
		c.di = di
		c.prealloc = true
		add(c)
	}
	for i, v := range [][2]uint16{{0, 0}, {0x1478, 0}, {0x1000, 0x478}, {0x1479, 0}, {0x2000, 0x100}} {
		c := base("naming-standalone", 0x192E3, i%4)
		c.ax, c.cx = v[0], v[1]
		add(c)
	}
	for i, v := range [][2]uint16{{0xF, 1}, {0x10, 0x10}, {5, 6}, {0x100, 0xFF}} {
		c := base("naming-standalone", 0x192EE, i)
		c.ax, c.cx = v[0], v[1]
		add(c)
	}
	for i, al := range []byte{0, 1, 0x40, 0x92, 0x93, 0xFF} {
		c := base("naming-standalone", 0x192F5, i%4)
		c.ds = 0x7000
		c.ax = uint16(al)
		add(c)
	}

	for code := 0x20; code <= 0x28; code++ {
		c := base("naming-integ", 0x18FC9, code%4)
		c.bx = uint16(code%4) * 64
		c.si = 4
		c.events = append(clickTape(clickCode[code].x, clickCode[code].y, 30), clickTape(clickCode[0x20].x, clickCode[0x20].y, 10)...)
		add(c)
	}
	{
		c := base("naming-integ", 0x18FC9, 0)
		c.bx = 0
		c.si = 0
		c.events = clickTape(clickCode[0x20].x, clickCode[0x20].y, 30)
		add(c)
	}

	for i, mode := range []int{3, 0x12, 3, 0x12} {
		c := base("video", 0x1EB6C, i)
		c.videoMode = mode
		c.seed = byte(i * 2)
		add(c)
	}
	for i, ah := range []byte{0, 1, 7, 15, 0xFF, 0x80} {
		c := base("video", 0x1EB93, i%4)
		c.ax = uint16(ah) << 8
		c.videoMode = 0x12
		if i == 5 {
			c.videoMode = 3
		}
		add(c)
	}

	for _, d2d := range []byte{0, 1, 0xFE, 0xFF} {
		for _, d2c := range []byte{0, 1} {
			c := base("tick", 0x10356, int(d2d)%4)
			c.d2c, c.d2d = d2c, d2d
			add(c)
		}
	}

	for i, al := range []byte{0, 1, 0x7F, 0xFF} {
		c := base("exit", 0x10043, i)
		c.ax = uint16(al)
		c.expectExit = true
		c.int61AH = []uint8{9, 8, 4, 0x0C}
		add(c)
	}
	for i := 0; i < 2; i++ {
		c := base("exit", 0x1006A, i)
		add(c)
	}
	return all
}

func subFlags16(a, b uint16) (cf, zf, sf, of bool) {
	r := a - b
	cf = a < b
	zf = r == 0
	sf = r&0x8000 != 0
	of = (a^b)&(a^r)&0x8000 != 0
	return cf, zf, sf, of
}

func f5scan(mem []byte, ds uint16, al byte) (cf bool, bx, cx uint16) {
	bx = 0x4240
	for i := 0; i < 0x7E; i++ {
		if mem[int(ds)*16+int(bx)+1] == al {
			return true, bx, uint16(0x7E - i)
		}
		bx += 0x20
	}
	return false, 0x4240 + 0x7E*0x20, 0
}

func startupAudit(before, after []byte, cs uint16, c startupCase, entry registers, observed []registers, device *machine.Machine, d *dos.DOS, soundBefore map[uint8]int, allocsBefore, readsBefore int) (int, []string) {
	n := 0
	var obs []string
	check := func(cond bool, msg string) {
		if !cond {
			obs = append(obs, "FAIL "+msg)
		} else {
			obs = append(obs, "ok "+msg)
		}
		n++
	}
	listed := map[uint8]bool{}
	for _, ah := range c.int61AH {
		listed[ah] = true
		check(d.Sound[ah] > soundBefore[ah], fmt.Sprintf("int61 count AH=%02X", ah))
	}
	var soundKeys []uint8
	for ah := range d.Sound {
		soundKeys = append(soundKeys, ah)
	}
	sort.Slice(soundKeys, func(i, j int) bool { return soundKeys[i] < soundKeys[j] })
	for _, ah := range soundKeys {
		v := d.Sound[ah]
		obs = append(obs, fmt.Sprintf("int61 AH=%02X %d->%d listed=%v", ah, soundBefore[ah], v, listed[ah]))
		n++
	}
	final := originalRegs(device)
	u16 := func(mem []byte, seg, off uint16) uint16 {
		return uint16(mem[int(seg)*16+int(off)]) | uint16(mem[int(seg)*16+int(off)+1])<<8
	}
	switch c.target {
	case 0x10000:
		if c.mouseAbsent {
			check(string(d.Console) == "ERROR: Mouse driver not install ?\n\r", "mouse error text")
			check(d.Exited && d.ExitCode == 0, "mouse error exit")
		} else {
			check(!d.Exited, "boot reaches loop")
		}
	case 0x10043:
		check(d.Exited && d.ExitCode == byte(entry.AX), "tail exit code")
		check(len(d.Console) == 0, "tail silence")
	case 0x1006A:
		check(final.IP == 0xF000 && final.SP == entry.SP+2, "shared retn")
	case 0x10356:
		check(after[int(cs)*16+0xD2C] == 1, "tick flag")
		want := c.d2d + 1
		if c.d2d == 0xFF {
			want = 0xFF
		}
		check(after[int(cs)*16+0xD2D] == want, "tick saturate")
		check(final.AX == entry.AX && final.BX == entry.BX && final.CX == entry.CX && final.DX == entry.DX, "tick regs")
		check(final.SP == entry.SP+4 && final.IP == 0xF000, "tick retf")
	case 0x10612:
		count := int(entry.AX) * 2
		src := before[int(entry.CX)*16+int(entry.BX) : int(entry.CX)*16+int(entry.BX)+count]
		dst := after[int(entry.CX)*16+int(entry.DX) : int(entry.CX)*16+int(entry.DX)+count]
		tmp := append([]byte(nil), src...)
		check(len(dst) == len(tmp), "far copy length")
		farBad := -1
		for i := range tmp {
			if dst[i] != tmp[i] {
				farBad = i
				break
			}
		}
		if farBad >= 0 {
			obs = append(obs, fmt.Sprintf("FAIL far copy bytes off=%d", farBad))
		} else {
			obs = append(obs, "ok far copy bytes")
		}
		n++
		check(final.CX == 0, "far copy cx")
		step := int(entry.AX) * 2
		if c.df {
			check(final.SI == entry.SI-uint16(step) && final.DI == entry.DI-uint16(step), "far copy down")
		} else {
			check(final.SI == entry.SI+uint16(step) && final.DI == entry.DI+uint16(step), "far copy up")
		}
	case 0x18DC8:
		check(final.AX == 0 || final.AX == 1, "prompt result")
	case 0x18FC9:
		newAllocs := d.Allocs[allocsBefore:]
		var seg uint16
		allocOK, freeOK := false, false
		for _, op := range newAllocs {
			if op.Fn == 0x48 && op.Want == 0x148 && op.OK {
				allocOK = true
				seg = op.Seg
			}
			if op.Fn == 0x49 {
				freeOK = true
			}
		}
		check(allocOK && freeOK, "naming alloc pair")
		check(u16(after, cs, 0x8FC7) == seg, "naming seg word")
		newReads := d.Reads[readsBefore:]
		check(len(newReads) == 1 && newReads[0].Want == 0x147A && newReads[0].Got == 0x147A && newReads[0].Seg == seg, "naming file read")
	case 0x1905E:
		check(u16(after, c.ds, entry.SI+0x5222) == 0x40A1, "cursor blank")
		if entry.SI == 0 {
			check(final.SI == 0, "cursor min")
		} else {
			check(final.SI == entry.SI-2, "cursor left")
		}
	case 0x1906E:
		want := entry.SI + 2
		if entry.SI == 0xA {
			want = 0xA
		}
		check(final.SI == want, "cursor right")
		check(u16(after, c.ds, want+0x5222) == 0x40A1, "cursor blank new")
	case 0x1907F:
		if entry.DI >= 0xC0 {
			check(final.DI == entry.DI-0xC0, "page up")
		} else {
			check(final.DI == 0, "page floor")
		}
	case 0x1908D:
		want := entry.DI + 0xC0
		if want > 0x1478 {
			want = 0x1478
		}
		if want > 0x13BA {
			want = 0x13BA
		}
		check(final.DI == want, "page down")
	case 0x190A3:
		seen := false
		for _, o := range observed {
			if o.IP == 0x928A {
				seen = true
			}
		}
		if entry.CX < 0xD8 {
			check(final.Flags&1 == 1 && !seen, "table borrow")
		} else {
			word := u16(before, cs, 0x8FA7+uint16(((entry.CX-0xD8)>>5)*2))
			check(seen, "table dispatch")
			check(final.DI == word, "table offset")
		}
	case 0x190C0:
		cx1 := entry.CX - 0xD2
		if entry.CX < 0xD2 || cx1 >= 0x140 {
			check(final.CX == cx1, "grid cx bound")
			break
		}
		dx1 := entry.DX - 0xEE
		if entry.DX < 0xEE || dx1 >= 0x78 {
			check(final.DX == dx1, "grid dx bound")
			break
		}
		col, colRem := cx1/20, cx1%20
		row, rowRem := dx1/20, dx1%20
		if colRem >= 0x10 || rowRem >= 0x10 {
			check(true, "grid rem bound")
			break
		}
		cell := col + row*16
		off := entry.DI
		if cell != 0 {
			off = entry.DI + cell*2
			if off > 0x1478 {
				off = 0x1478
			}
		}
		raw := u16(before, 0x9000, off)
		swapped := raw>>8 | raw<<8
		v := swapped - 0x1000
		out := v>>8 | v<<8
		check(u16(after, c.ds, entry.SI+0x5222) == out, "grid char")
		if entry.SI == 0xA {
			check(final.SI == 0xA, "grid cursor max")
		} else {
			check(final.SI == entry.SI+2, "grid cursor advance")
		}
	case 0x1912D, 0x19144:
		al := before[int(c.ds)*16+0x5221]
		for i := 0; i < 300; i++ {
			if c.target == 0x1912D {
				al--
				if al >= 0x93 {
					al = 0x92
				}
			} else {
				al++
				if al >= 0x93 {
					al = 0
				}
			}
			if cf, _, _ := f5scan(before, c.ds, al); !cf {
				break
			}
		}
		check(after[int(c.ds)*16+0x5221] == al, "cycle value")
	case 0x192E3:
		sum := entry.AX + entry.CX
		want := sum
		if sum > 0x1478 {
			want = 0x1478
		}
		check(final.AX == want, "clamp high")
		cf, zf, _, _ := subFlags16(sum, 0x1478)
		check((final.Flags&1 == 1) == cf && (final.Flags&0x40 == 0x40) == zf, "clamp flags")
	case 0x192EE:
		if entry.AX >= entry.CX {
			check(final.AX == entry.AX-entry.CX, "clamp sub")
		} else {
			check(final.AX == 0 && final.Flags&0x40 != 0, "clamp floor")
		}
	case 0x192F5:
		cf, bx, cx := f5scan(before, c.ds, byte(entry.AX))
		check((final.Flags&1 == 1) == cf && final.BX == bx && final.CX == cx, "faction scan")
	case 0x18F7C:
		check(after[int(c.ds)*16+int(entry.SI)+2] == before[int(c.ds)*16+int(entry.SI)+0x3F], "name reset copy")
		check(after[int(c.ds)*16+0x5221] == 0xFF, "name reset flag")
		for i := 0; i < 6; i++ {
			check(u16(after, c.ds, 0x5222+uint16(i*2)) == 0xD0A1, "name reset fill")
		}
	case 0x1EB6C:
		var seq []byte
		for _, w := range device.PortLog {
			if w.Port == 0x3C0 {
				seq = append(seq, w.Val)
			}
		}
		want := []byte{0x14, 0}
		for i := 0; i < 16; i++ {
			want = append(want, byte(i), byte(i))
		}
		want = append(want, 0x11, 0, 0x20)
		check(len(seq) == len(want), "palette length")
		palBad := -1
		for i := range want {
			if i >= len(seq) || seq[i] != want[i] {
				palBad = i
				break
			}
		}
		if palBad >= 0 {
			obs = append(obs, fmt.Sprintf("FAIL palette seq i=%d", palBad))
		} else {
			obs = append(obs, "ok palette seq")
		}
		n++
	case 0x1EB93:
		type pair struct {
			port uint16
			val  byte
		}
		var seq []pair
		for _, w := range device.PortLog {
			seq = append(seq, pair{w.Port, w.Val})
		}
		ah := byte(entry.AX >> 8)
		want := []pair{{0x3C0, 0x11}, {0x3C0, ah}, {0x3C0, 0x20}, {0x3CE, 5}, {0x3CF, 3}, {0x3CE, 0}, {0x3CF, ah}, {0x3CE, 3}, {0x3CF, 0}, {0x3CE, 8}, {0x3CF, 0xFF}}
		check(len(seq) == len(want), "fill seq length")
		fillBad := -1
		for i := range want {
			if i >= len(seq) || seq[i] != want[i] {
				fillBad = i
				break
			}
		}
		if fillBad >= 0 {
			obs = append(obs, fmt.Sprintf("FAIL fill seq i=%d", fillBad))
		} else {
			obs = append(obs, "ok fill seq")
		}
		n++
	case 0x18CAE:
		if byte(entry.AX>>8) == 0 {
			check(after[int(cs)*16+0x10CF8] == 0, "block zero flag")
		}
	}
	_ = final
	_ = observed
	return n, obs
}
