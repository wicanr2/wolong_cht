//go:build ignore

// 原版/C/Go 完整月結 raw-state audit；正式引擎不含 matching 直接入口。
package main

/*
#include <stdlib.h>
#include "/repo/tools/c_recovery/rng.c"
#include "/repo/tools/c_recovery/economy.c"
#include "/repo/tools/c_recovery/settlement.c"
#include "/repo/tools/c_recovery/world_update.c"
#include "/repo/tools/c_recovery/politics.c"
#include "/repo/tools/c_recovery/politics_fixture.h"
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
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/state"
	"os"
	"path/filepath"
	"unsafe"
)

func mcSet(m *C.KiMachine16, r oracle.Regs) {
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
func mcRegs(m *C.KiMachine16) oracle.Regs {
	return oracle.Regs{AX: uint16(m.ax), BX: uint16(m.bx), CX: uint16(m.cx), DX: uint16(m.dx), SI: uint16(m.si), DI: uint16(m.di), BP: uint16(m.bp), SP: uint16(m.sp), DS: uint16(m.ds), ES: uint16(m.es), SS: uint16(m.ss), CS: uint16(m.cs), IP: uint16(m.ip), Flags: uint16(m.flags)}
}
func mcPut(b []byte, at int, v uint16) { binary.LittleEndian.PutUint16(b[at:], v) }
func mcDiff(a, b []byte) []map[string]any {
	var rows []map[string]any
	for i, v := range a {
		if v != b[i] {
			domain := "unloaded"
			switch {
			case i < 59:
				domain = "globals"
			case i >= 0x80 && i < 0x80+0x580:
				domain = "factions"
			case i >= 0x8c0 && i < 0x20c0:
				domain = "cities"
			case i >= 0x22c0 && i < 0x4280:
				domain = "corps"
			case i >= 0x42c0 && i < 0x52c0:
				domain = "generals"
			case i >= 0x52c0:
				domain = "queue"
			case i >= 0x80:
				domain = "other-world"
			}
			rows = append(rows, map[string]any{"file_offset": i, "domain": domain, "original": v, "go": b[i]})
		}
	}
	return rows
}
func main() {
	out := flag.String("out", "/output/audit.json", "收據")
	dir := flag.String("vectors", "/output/vectors", "向量輸出")
	flag.Parse()
	raw, err := os.ReadFile("/orig/KI.EXE")
	if err != nil {
		panic(err)
	}
	eh := fmt.Sprintf("%x", sha256.Sum256(raw))
	if eh != "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868" {
		panic("input identity")
	}
	scenarioBytes, err := os.ReadFile("/orig/SINARIO.DAT")
	if err != nil || len(scenarioBytes) != 88832 {
		panic("scenario")
	}
	if err := os.MkdirAll(*dir, 0755); err != nil {
		panic(err)
	}
	native := C.malloc(1 << 20)
	if native == nil {
		panic("malloc")
	}
	defer C.free(native)
	cm := unsafe.Slice((*byte)(native), 1<<20)
	m := C.KiMachine16{memory: (*C.uint8_t)(native)}
	f := C.KiPoliticsFixture{}
	var records []map[string]any
	allCDigest, allGoDigest := sha256.New(), sha256.New()
	passed, diverged := 0, 0
	for scenario := 0; scenario < 4; scenario++ {
		for _, player := range []int{0, 7, 21} {
			for _, seed := range []int{0, 77, 255} {
				name := fmt.Sprintf("s%d-p%d-seed%d", scenario+1, player, seed)
				initial := append([]byte(nil), scenarioBytes[scenario*22208:(scenario+1)*22208]...)
				mcPut(initial, 13, uint16(player*64))
				initial[15] = byte(player)
				initial[24] = 50
				initial[25] = 0
				for i := 0; i < 3; i++ {
					mcPut(initial, 26+i*2, 5000)
					mcPut(initial, 34+i*2, 2500)
				}
				mcPut(initial, 32, 50)
				mcPut(initial, 48, 0)
				for i, v := range []uint16{0xfff0, 0xfff0, 400, 400} {
					mcPut(initial, 50+i*2, v)
				}
				originalRNG := rng.New(12, 34, 56).Raw()
				originalRNG[0] = byte(seed)
				originalRNG[1] = byte(seed * 37)
				w, err := state.LoadBlock(initial)
				if err != nil {
					panic(err)
				}
				w.EnableStrategicAI()
				w.SetApproximateEvent10(false)
				beforeGo := w.Bytes()
				preDiff := mcDiff(initial, beforeGo)
				o, err := oracle.Load("/orig/KI.EXE", "/orig")
				if err != nil {
					panic(err)
				}
				cs := o.Regs().CS
				for _, target := range []uint16{0x5e80, 0xce7, 0xcde, 0x8810} {
					o.WriteBytes(oracle.Addr{Seg: cs, Off: target}, []byte{0xc3})
				}
				ss := uint16(0x4000)
				o.WriteBytes(oracle.Addr{Seg: 0x7000, Off: 0}, []byte{0xb8, 0, 0x40, 0x8e, 0xd0, 0xbc, 2, 0x80, 0x68, 2, 0, 0x9d})
				o.CallNear(oracle.Addr{Seg: 0x7000, Off: 0}, 0xbeef, oracle.CallRegs{})
				if err := o.Run(5); err != nil {
					panic(err)
				}
				o.WriteBytes(oracle.Addr{Seg: cs, Off: 0xcf0}, initial[:59])
				o.WriteBytes(oracle.Addr{Seg: 0x2200, Off: 0}, initial[0x80:0x52c0])
				o.WriteBytes(oracle.Addr{Seg: 0x6000, Off: 0}, initial[0x52c0:])
				scratch := make([]byte, 1056)
				for i := range scratch {
					scratch[i] = 255
				}
				o.WriteBytes(oracle.Addr{Seg: 0x6500, Off: 0}, scratch)
				o.WriteU16(oracle.Addr{Seg: cs, Off: 0xd52}, 0x2200)
				o.WriteU16(oracle.Addr{Seg: cs, Off: 0xd56}, 0x6000)
				o.WriteU16(oracle.Addr{Seg: cs, Off: 0x987c}, 0x6500)
				o.WriteU8(oracle.Addr{Seg: cs, Off: 0x31ad}, 7)
				o.WriteBytes(oracle.Addr{Seg: cs, Off: 0xecfc}, originalRNG)
				o.CallNear(oracle.Addr{Seg: cs, Off: 0x5358}, 0x100, oracle.CallRegs{AX: 0x1234, BX: 0x2345, CX: 0x3456, DX: 0x4567, SI: 0x5678, DI: 0x6789, BP: 0x789a, DS: cs, ES: 0x6500, SetAX: true, SetBX: true, SetCX: true, SetDX: true, SetSI: true, SetDI: true, SetBP: true, SetDS: true, SetES: true})
				input := o.Regs()
				copy(cm, o.Bytes(oracle.Phys(0), 1<<20))
				mcSet(&m, input)
				f.calls = 0
				if err := o.RunUntil(oracle.At(oracle.Addr{Seg: cs, Off: 0x100}), oracle.Budget(400000)); err != nil {
					panic(err)
				}
				C.politics_run(&m, 0x5358, &f)
				if mcRegs(&m) != o.Regs() || !bytes.Equal(cm, o.Bytes(oracle.Phys(0), 1<<20)) {
					panic("C/original drift")
				}
				original := append([]byte(nil), initial...)
				copy(original[:59], o.Bytes(oracle.Addr{Seg: cs, Off: 0xcf0}, 59))
				copy(original[0x80:0x52c0], o.Bytes(oracle.Addr{Seg: 0x2200, Off: 0}, 0x5240))
				copy(original[0x52c0:], o.Bytes(oracle.Addr{Seg: 0x6000, Off: 0}, 1024))
				finalRNG := o.Bytes(oracle.Addr{Seg: cs, Off: 0xecfc}, 258)
				g, _ := rng.FromRaw(originalRNG)
				ev := w.MatchingMonthly(g)
				goBlock := w.Bytes()
				diff := mcDiff(original, goBlock)
				rngEqual := bytes.Equal(g.Raw(), finalRNG)
				if len(diff) == 0 && rngEqual {
					passed++
				} else {
					diverged++
				}
				snaps := w.TakeSnapshot()
				record := map[string]any{"name": name, "scenario": scenario + 1, "player": player, "counter": byte(seed), "state_index": byte(seed * 37), "input_sha256": fmt.Sprintf("%x", sha256.Sum256(initial)), "precondition_diffs": preDiff, "postcondition_diffs": diff, "rng_equal": rngEqual, "original_rng": hex.EncodeToString(finalRNG), "go_rng": hex.EncodeToString(g.Raw()), "go_runtime": snaps, "settled": ev.Settled, "original_event_cursor": o.Word(oracle.Addr{Seg: cs, Off: 0xd20}), "original_event_delay": o.Byte(oracle.Addr{Seg: cs, Off: 0x31ad}), "original_go_equal": len(diff) == 0 && rngEqual}
				records = append(records, record)
				allCDigest.Write(original)
				allCDigest.Write(finalRNG)
				allGoDigest.Write(goBlock)
				allGoDigest.Write(g.Raw())
				for suffix, data := range map[string][]byte{"input.block": initial, "original.block": original, "go.block": goBlock, "rng-before.bin": originalRNG, "rng-original.bin": finalRNG, "rng-go.bin": g.Raw()} {
					if err := os.WriteFile(filepath.Join(*dir, name+"."+suffix), data, 0644); err != nil {
						panic(err)
					}
				}
				fmt.Printf("%s: initial diffs %d; monthly diffs %d; RNG=%v\n", name, len(preDiff), len(diff), rngEqual)
				o.Close()
				_ = ss
			}
		}
	}
	report := map[string]any{"schema": "wolong-c-go-monthly-audit-v1", "input_sha256": eh, "scenario_sha256": fmt.Sprintf("%x", sha256.Sum256(scenarioBytes)), "cases": len(records), "original_c_cases": len(records), "go_passed": passed, "go_diverged": diverged, "vectors": records, "original_state_sha256": hex.EncodeToString(allCDigest.Sum(nil)), "go_state_sha256": hex.EncodeToString(allGoDigest.Sum(nil)), "scope": "Local sub_15358 rules, four scenarios, UI/sound/redraw RET fixtures, matching-only Go monthly entry; strategic AI enabled, approximate event10 disabled; IF/TF=0", "c_machine_code_match": false}
	encoded, _ := json.MarshalIndent(report, "", "  ")
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*out, encoded, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("original/C %d; Go pass %d diverged %d\n", len(records), passed, diverged)
}
