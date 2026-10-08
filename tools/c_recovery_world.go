//go:build ignore

// 月結世界更新與事件寫入器的原版/C/Go 對照；spec/209。
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/rng.c"
#include "/repo/tools/c_recovery/economy.c"
#include "/repo/tools/c_recovery/settlement.c"
#include "/repo/tools/c_recovery/world_update.c"
#include "/repo/tools/c_recovery/world_update_fixture.h"
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
	"github.com/wicanr2/dosgolem/oracle"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"os"
	"unsafe"
)

func wSet(m *C.KiMachine16, r oracle.Regs) {
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
func wRegs(m *C.KiMachine16) oracle.Regs {
	return oracle.Regs{AX: uint16(m.ax), BX: uint16(m.bx), CX: uint16(m.cx), DX: uint16(m.dx), SI: uint16(m.si), DI: uint16(m.di), BP: uint16(m.bp), SP: uint16(m.sp), DS: uint16(m.ds), ES: uint16(m.es), SS: uint16(m.ss), CS: uint16(m.cs), IP: uint16(m.ip), Flags: uint16(m.flags)}
}
func wRegBytes(r oracle.Regs) []byte {
	b := make([]byte, 28)
	for i, v := range []uint16{r.AX, r.BX, r.CX, r.DX, r.SI, r.DI, r.BP, r.SP, r.DS, r.ES, r.SS, r.CS, r.IP, r.Flags} {
		binary.LittleEndian.PutUint16(b[i*2:], v)
	}
	return b
}
func wWord(b []byte, at int) uint16   { return binary.LittleEndian.Uint16(b[at:]) }
func wPut(b []byte, at int, v uint16) { binary.LittleEndian.PutUint16(b[at:], v) }
func w24(b []byte, at, v int)         { b[at] = byte(v); b[at+1] = byte(v >> 8); b[at+2] = byte(v >> 16) }
func wBank() []byte {
	b := make([]byte, 0x5240)
	for i := 0; i < 22; i++ {
		b[i*64+0x2a] = 0xff
	}
	b[0] = 0x80
	w24(b, 0x20, 500000)
	b[0x28] = 16
	for i := 0; i < 192; i++ {
		at := 0x840 + i*32
		b[at+1] = 0
		wPut(b, at+8, uint16(i*2%384))
		wPut(b, at+10, uint16(i%256))
		wPut(b, at+12, 60000)
		wPut(b, at+14, 3000)
		b[at+16] = 100
		b[at+17] = 100
		b[at+19] = 10
		b[at+0x19] = 0xff
	}
	for i := 0; i < 22*24; i++ {
		b[0x600+i] = 0x20
	}
	return b
}

type wCase struct {
	group                          string
	target, ax, bx, si, di, cursor uint16
	bank, queue                    []byte
	player                         byte
	tax                            byte
	counter, state                 byte
	rect                           [4]uint16
	goGrowth                       bool
}

func main() {
	out := flag.String("out", "/output/results/O2.json", "收據")
	only := flag.String("group", "", "單一群組")
	smoke := flag.Bool("smoke", false, "每群組 12 案例")
	skipGo := flag.Bool("skip-go", false, "只比原版/C")
	probe := flag.Bool("go-probe", false, "生產力位元寬度反例")
	flag.Parse()
	raw, err := os.ReadFile("/orig/KI.EXE")
	if err != nil {
		panic(err)
	}
	eh := fmt.Sprintf("%x", sha256.Sum256(raw))
	if eh != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" {
		panic("input identity")
	}
	scenarios, err := os.ReadFile("/orig/SINARIO.DAT")
	if err != nil || len(scenarios) != 88832 {
		panic("scenario")
	}
	o, err := oracle.Load("/orig/KI.EXE", "/orig")
	if err != nil {
		panic(err)
	}
	defer o.Close()
	cs := o.Regs().CS
	funcs := map[uint16]int{0x5695: 128, 0x55a6: 70, 0x2fbf: 79, 0x57fe: 42, 0x5715: 122, 0x578f: 111, 0x30cb: 8, 0x3119: 31, 0x22db: 163, 0x2286: 85, 0x237e: 129}
	rh := map[string]string{}
	for target, size := range funcs {
		b := raw[int(target)+0x200 : int(target)+0x200+size]
		if !bytes.Equal(o.Bytes(oracle.Addr{Seg: cs, Off: target}, size), b) {
			panic("routine identity")
		}
		rh[fmt.Sprintf("sub_%X", 0x10000+uint32(target))] = fmt.Sprintf("%x", sha256.Sum256(b))
	}
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	mem := unsafe.Slice((*byte)(native), 1<<20)
	copy(mem, o.Bytes(oracle.Phys(0), 1<<20))
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	both := func(seg, off uint16, b []byte) {
		a := oracle.Addr{Seg: seg, Off: off}
		o.WriteBytes(a, b)
		copy(mem[a.Linear():], b)
	}
	word := func(seg, off, v uint16) { var b [2]byte; binary.LittleEndian.PutUint16(b[:], v); both(seg, off, b[:]) }
	bankSeg, queueSeg, ss := uint16(0x2200), uint16(0x6000), uint16(0x4000)
	ba, qa := oracle.Addr{Seg: bankSeg, Off: 0}, oracle.Addr{Seg: queueSeg, Off: 0}
	globals := oracle.Addr{Seg: cs, Off: 0xcf0}
	ra := oracle.Addr{Seg: cs, Off: 0xecfc}
	table := rng.New(12, 34, 56).Raw()
	var trace []byte
	targets := []uint16{0x5358, 0x53c6, 0x5456, 0x548f, 0x5538, 0x5547, 0x54fc, 0x55ec, 0x5609, 0x563b, 0x5828, 0xece0, 0x5695, 0x55a6, 0x2fbf, 0x57fe, 0x5715, 0x578f, 0x30cb, 0x3119, 0x22db, 0x2286, 0x237e, 0x585f, 0x2bd9, 0x5e80, 0xce7, 0x8810}
	for _, target := range targets {
		t := target
		o.OnCall(oracle.Addr{Seg: cs, Off: t}, func(o *oracle.Oracle) {
			r := o.Regs()
			var b [2]byte
			binary.LittleEndian.PutUint16(b[:], t)
			trace = append(trace, b[:]...)
			trace = append(trace, wRegBytes(r)...)
			for i := uint16(0); i < 64; i++ {
				trace = append(trace, o.Byte(oracle.Addr{Seg: r.DS, Off: r.SI + i}))
			}
			trace = append(trace, o.Bytes(oracle.Addr{Seg: r.SS, Off: r.SP}, 32)...)
			trace = append(trace, o.Bytes(globals, 59)...)
			trace = append(trace, o.Bytes(ra, 2)...)
		})
	}
	for _, target := range targets[23:] {
		both(cs, target, []byte{0xc3})
	}
	f := C.KiWorldFixture{}
	cases, goCases, audits := 0, 0, 0
	groups, goGroups := map[string]int{}, map[string]int{}
	od, cd := sha256.New(), sha256.New()
	var failure, goFailure map[string]any
	var probeResults []map[string]any
	check := func(c wCase) bool {
		if *only != "" && *only != c.group {
			return true
		}
		if *smoke && groups[c.group] >= 12 {
			return true
		}
		trace = nil
		f.calls = 0
		pre := o.Regs()
		word(pre.SS, pre.SP-2, 0xbeef)
		flags := uint16(2)
		if cases&1 != 0 {
			flags |= 0x400
		}
		both(0x7000, 0, []byte{0xb8, 0, 0x40, 0x8e, 0xd0, 0xbc, 2, 0x80, 0x68, byte(flags), byte(flags >> 8), 0x9d})
		o.CallNear(oracle.Addr{Seg: 0x7000, Off: 0}, 0xbeef, oracle.CallRegs{})
		if err := o.Run(5); err != nil {
			panic(err)
		}
		word(ss, 0x8000, 0x100)
		both(bankSeg, 0, c.bank)
		q := c.queue
		if q == nil {
			q = make([]byte, 1024)
		}
		both(queueSeg, 0, q)
		g := make([]byte, 59)
		g[15] = c.player
		g[24] = c.tax
		wPut(g, 13, uint16(c.player)*64)
		for i := 0; i < 3; i++ {
			wPut(g, 26+i*2, 5000)
			wPut(g, 32+i*2, 2500)
		}
		wPut(g, 48, c.cursor)
		for i, v := range c.rect {
			wPut(g, 50+i*2, v)
		}
		both(cs, 0xcf0, g)
		word(cs, 0xd52, bankSeg)
		word(cs, 0xd56, queueSeg)
		initial := append([]byte(nil), table...)
		initial[0] = c.counter
		initial[1] = c.state
		both(cs, 0xecfc, initial)
		r := oracle.CallRegs{AX: c.ax, BX: c.bx, CX: 0x3456, DX: 0x4567, SI: c.si, DI: c.di, BP: 0x6789, DS: bankSeg, ES: 0x789a, SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true}
		if c.target == 0x5358 {
			r.DS = cs
		}
		o.CallNear(oracle.Addr{Seg: cs, Off: c.target}, 0x100, r)
		input := o.Regs()
		wSet(&m, input)
		if err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x100}), oracle.Budget(250000)); err != nil {
			panic(err)
		}
		C.world_run(&m, C.uint16_t(c.target), &f)
		got, want := wRegs(&m), o.Regs()
		if int(f.calls) > 4096 {
			panic("trace capacity")
		}
		ct := make([]byte, 0, int(f.calls)*187)
		for i := 0; i < int(f.calls); i++ {
			s := f.trace[i]
			var b [2]byte
			binary.LittleEndian.PutUint16(b[:], uint16(s.target))
			ct = append(ct, b[:]...)
			for _, v := range s.regs {
				binary.LittleEndian.PutUint16(b[:], uint16(v))
				ct = append(ct, b[:]...)
			}
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.record[0]), 64)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.stack[0]), 32)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.globals[0]), 59)...)
			ct = append(ct, C.GoBytes(unsafe.Pointer(&s.rng[0]), 2)...)
		}
		br, qr, gr, rr := o.Bytes(ba, len(c.bank)), o.Bytes(qa, 1024), o.Bytes(globals, 59), o.Bytes(ra, 258)
		cases++
		groups[c.group]++
		if got != want || !bytes.Equal(trace, ct) || !bytes.Equal(mem[ba.Linear():ba.Linear()+uint32(len(br))], br) || !bytes.Equal(mem[qa.Linear():qa.Linear()+1024], qr) || !bytes.Equal(mem[globals.Linear():globals.Linear()+59], gr) || !bytes.Equal(mem[ra.Linear():ra.Linear()+258], rr) {
			failure = map[string]any{"group": c.group, "case": cases - 1, "input": input, "original": want, "c": got, "original_trace": hex.EncodeToString(trace), "c_trace": hex.EncodeToString(ct), "original_bank": hex.EncodeToString(br), "c_bank": hex.EncodeToString(mem[ba.Linear() : ba.Linear()+uint32(len(br))]), "original_queue": hex.EncodeToString(qr), "c_queue": hex.EncodeToString(mem[qa.Linear() : qa.Linear()+1024])}
			return false
		}
		for off := uint16(0x7f80); off < 0x8002; off++ {
			a := oracle.Addr{Seg: ss, Off: off}
			if mem[a.Linear()] != o.Byte(a) {
				panic("stack mismatch")
			}
		}
		od.Write(wRegBytes(want))
		od.Write(br)
		od.Write(qr)
		od.Write(gr)
		od.Write(rr)
		od.Write(trace)
		cd.Write(wRegBytes(got))
		cd.Write(br)
		cd.Write(qr)
		cd.Write(gr)
		cd.Write(rr)
		cd.Write(ct)
		if *probe {
			probeResults = append(probeResults, map[string]any{"original_production": wWord(br, 0x84e), "original_growth": int(br[0x850]) - 100})
		}
		if c.goGrowth && !*skipGo {
			rngGo, _ := rng.FromRaw(initial)
			ok := true
			var details any
			for i := 0; i < 192; i++ {
				at := 0x840 + i*32
				city := economy.CityState{Owner: int(c.bank[at+1]), Production: int(wWord(c.bank, at+14)), ProductionCap: int(wWord(c.bank, at+12)), Growth: int(c.bank[at+16]) - 100}
				economy.GrowCity(&city, int(c.tax), city.Owner == int(c.player), rngGo)
				if city.Production != int(wWord(br, at+14)) || city.Growth != int(br[at+16])-100 {
					ok = false
					details = map[string]any{"city": i, "go": city, "original_production": wWord(br, at+14), "original_growth": int(br[at+16]) - 100}
					break
				}
			}
			ok = ok && bytes.Equal(rngGo.Raw(), rr)
			goCases++
			goGroups[c.group]++
			if !ok {
				goFailure = map[string]any{"group": c.group, "case": cases - 1, "details": details}
				return false
			}
		}
		if cases%128 == 0 {
			if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
				panic("full memory")
			}
			audits++
		}
		return true
	}
	base := func(group string, target uint16) wCase {
		return wCase{group: group, target: target, ax: 0x1234, bx: 0x2345, bank: wBank(), player: 0, tax: 30, rect: [4]uint16{0xfff0, 0xfff0, 400, 400}}
	}
	if *probe {
		for _, setting := range []struct{ production, growth, tax int }{{65000, 200, 30}, {65535, 200, 0}, {65535, 0, 100}} {
			c := base("growth", 0x5695)
			wPut(c.bank, 0x84e, uint16(setting.production))
			wPut(c.bank, 0x84c, 65535)
			c.bank[0x850] = byte(setting.growth)
			c.tax = byte(setting.tax)
			c.goGrowth = true
			if !check(c) {
				goto finish
			}
		}
		goto finish
	}
	if *only == "" || *only == "growth" {
		for _, prod := range []uint16{0, 1, 255, 256, 32767, 40000, 65000, 65535} {
			for _, growth := range []byte{0, 1, 99, 100, 101, 199, 200} {
				for _, tax := range []byte{0, 29, 30, 50, 100} {
					for player := 0; player < 2; player++ {
						c := base("growth", 0x5695)
						wPut(c.bank, 0x84e, prod)
						wPut(c.bank, 0x84c, 65535)
						c.bank[0x850] = growth
						c.tax = tax
						c.player = byte(player)
						c.counter = byte(prod)
						c.state = growth
						c.goGrowth = true
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	if *only == "" || *only == "score" {
		for attr := 0; attr < 256; attr++ {
			for _, status := range []byte{0, 127, 128, 255} {
				for slot := 0; slot < 128; slot++ {
					c := base("score", 0x55a6)
					at := 0x4240 + slot*32
					c.bank[at] = status
					for i := 0; i < 5; i++ {
						c.bank[at+14+i] = byte(attr + i*17)
					}
					c.bank[at+31] = 0xaa
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "writer" {
		for _, cursor := range []uint16{0, 4, 124, 252, 256} {
			for _, slot := range []byte{0, 1, 31, 63, 64, 255} {
				for mode := 0; mode < 4; mode++ {
					for seed := 0; seed < 8; seed++ {
						c := base("writer", 0x2fbf)
						c.ax = 0x1204
						c.bx = 0x7700 | uint16(slot)
						c.cursor = cursor
						c.counter = byte(seed * 31)
						c.state = byte(seed * 17)
						c.queue = make([]byte, 1024)
						for i := 0; i < 64; i++ {
							if mode == 1 || mode == 2 && i%2 == 0 {
								wPut(c.queue, i*4, 0x0101)
							}
							if mode == 3 {
								wPut(c.queue, i*4, 0x1200)
							}
						}
						if !check(c) {
							goto finish
						}
					}
				}
			}
		}
	}
	if *only == "" || *only == "trust" {
		for _, funds := range []int{0, -1, -9728, -9984, -10240, -655000} {
			for _, threshold := range []byte{0, 1, 15, 16, 255} {
				for seed := 0; seed < 16; seed++ {
					c := base("trust", 0x57fe)
					w24(c.bank, 0x20, funds)
					c.bank[0x28] = threshold
					c.counter = byte(seed * 17)
					c.state = byte(seed * 13)
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "governor" {
		for slot := 0; slot < 192; slot++ {
			for mode := 0; mode < 4; mode++ {
				c := base("governor", 0x5715)
				at := 0x840 + slot*32
				c.bank[at+0x19] = 0
				c.bank[at+16] = byte(mode * 60)
				c.bank[at+17] = byte(mode * 50)
				c.bank[at+18] = 180
				c.bank[at+19] = byte(mode * 20)
				if mode == 3 {
					c.bank[0x425a] = 1
				}
				c.counter = byte(slot)
				if !check(c) {
					goto finish
				}
			}
		}
	}
	if *only == "" || *only == "diplomat" {
		for owner := 0; owner < 22; owner++ {
			for target := 0; target < 22; target++ {
				for _, value := range []byte{0, 1, 99, 100, 127, 128, 255} {
					c := base("diplomat", 0x578f)
					c.player = byte(owner)
					c.bank[target*64+0x2a] = 0
					c.bank[0x600+owner*24+target] = value
					c.bank[0x600+target*24+owner] = byte(255 - value)
					c.counter = byte(target)
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "relation" {
		for owner := 0; owner < 22; owner++ {
			for target := 0; target < 22; target++ {
				for _, entry := range []uint16{0x30cb, 0x3119} {
					c := base("relation", entry)
					c.si = uint16(owner * 64)
					c.di = uint16(target * 64)
					c.bank[0x600+owner*24+target] = byte(owner*11 + target)
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "disaster" {
		for seed := 0; seed < 256; seed++ {
			for _, defence := range []byte{0, 23, 24, 63, 64, 100} {
				c := base("disaster", 0x2286)
				for i := 0; i < 192; i++ {
					c.bank[0x840+i*32+16] = defence
					c.bank[0x840+i*32+17] = defence
				}
				c.counter = byte(seed)
				c.state = byte(seed * 37)
				if !check(c) {
					goto finish
				}
			}
		}
	}
	if *only == "" || *only == "storm" {
		for seed := 0; seed < 256; seed++ {
			for old := 0; old < 2; old++ {
				for _, x := range []uint16{0, 9, 10, 191, 192, 383} {
					c := base("storm", 0x22db)
					if old == 1 {
						c.rect = [4]uint16{10, 20, 30, 40}
					}
					for i := 0; i < 192; i++ {
						wPut(c.bank, 0x840+i*32+8, x)
					}
					c.counter = byte(seed)
					c.state = byte(seed * 17)
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "marker" {
		for _, strength := range []uint16{0, 1, 10, 15} {
			for _, distance := range []uint16{0, 1, 19, 20, 21, 255} {
				for owner := 0; owner < 2; owner++ {
					c := base("marker", 0x237e)
					c.ax = strength
					c.rect = [4]uint16{100, 100, 100, 100}
					for i := 0; i < 192; i++ {
						at := 0x840 + i*32
						wPut(c.bank, at+8, 100+distance)
						wPut(c.bank, at+10, 100)
						c.bank[at+1] = byte(owner)
						c.bank[at+21] = 99
					}
					if !check(c) {
						goto finish
					}
				}
			}
		}
	}
	if *only == "" || *only == "scenario-monthly" {
		for scenario := 0; scenario < 4; scenario++ {
			for _, player := range []byte{0, 7, 21} {
				c := base("scenario-monthly", 0x5358)
				start := scenario*22208 + 0x80
				c.bank = append([]byte(nil), scenarios[start:start+0x5240]...)
				c.player = player
				c.tax = 50
				c.counter = byte(scenario * 23)
				c.state = byte(player * 11)
				if !check(c) {
					goto finish
				}
			}
		}
	}
finish:
	if failure == nil && goFailure == nil {
		if !bytes.Equal(mem, o.Bytes(oracle.Phys(0), 1<<20)) {
			panic("final memory")
		}
		audits++
	}
	report := map[string]any{"schema": "wolong-c-world-parity-v1", "input_sha256": eh, "scenario_sha256": fmt.Sprintf("%x", sha256.Sum256(scenarios)), "routine_sha256": rh, "cases": cases, "go_cases": goCases, "groups": groups, "go_groups": goGroups, "full_memory_audits": audits, "passed": failure == nil && goFailure == nil, "mismatch": failure, "go_mismatch": goFailure, "probe_results": probeResults, "original_state_sha256": hex.EncodeToString(od.Sum(nil)), "c_state_sha256": hex.EncodeToString(cd.Sum(nil)), "c_machine_code_match": false, "rng_fixture": map[string]any{"table_seed_hms": []int{12, 34, 56}, "initial_counter_state": "Explicit per case, written to both CS:ECFC before execution", "go": "rng.FromRaw identical pre-call state", "rerolls": 0}, "scope": "Eleven world-update functions linked to C economic monthly path; sub_1585F/sub_12BD9 and UI/sound/redraw RET fixtures; IF/TF=0; legal pointers/indices/tax"}
	encoded, _ := json.MarshalIndent(report, "", "  ")
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*out, encoded, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("world original/C %d; Go %d; audits %d; pass=%v\n", cases, goCases, audits, report["passed"])
	if failure != nil || goFailure != nil {
		os.Exit(1)
	}
}
