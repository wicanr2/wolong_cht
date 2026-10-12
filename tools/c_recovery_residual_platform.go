//go:build matching_residual

package main

/*
#include "/repo/tools/c_recovery/rng.h"
*/
import "C"
import (
	"fmt"
	"sort"
	"strings"

	"github.com/wicanr2/dosgolem/internal/machine"
)

//export wolong_residual_checkpoint
func wolong_residual_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	residualCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_residual_snapshot
func wolong_residual_snapshot(m *C.KiMachine16, t C.uint16_t) {
	residualCCalls++
	cur := cRegs(m)
	cur.IP = uint16(t)
	loc := uint32(cur.IP) + 0x10000
	if cur.CS == residualMouseCS {
		loc = uint32(cur.IP) + 0x20000
	}
	nativeHist[loc]++
	evtEmit(fmt.Sprintf("C native %X/%s", loc, entryName(loc)))
	if residualCCalls > 2000000 {
		type kv struct {
			loc uint32
			n   int
		}
		var top []kv
		for l, n := range nativeHist {
			top = append(top, kv{l, n})
		}
		sort.Slice(top, func(i, j int) bool { return top[i].n > top[j].n })
		var sb strings.Builder
		for i := 0; i < len(top) && i < 8; i++ {
			fmt.Fprintf(&sb, " %X/%s=%d", top[i].loc, entryName(top[i].loc), top[i].n)
		}
		panic(fmt.Sprintf("bounded native residual entry count group=%s target=%X hist:%s", curGroup, curTarget, sb.String()))
	}
	if _, known := functions[loc]; !known {
		return
	}
	residualCTrace = append(residualCTrace, byte(t), byte(uint16(t)>>8))
	residualCTrace = append(residualCTrace, regBytes(cur)...)
	for i := uint32(0); i < 32; i++ {
		residualCTrace = append(residualCTrace, activeDevice.Mem[(uint32(cur.SS)*16+uint32(cur.SP)+i)&0xfffff])
	}
	residualCTrace = append(residualCTrace, deviceBytes(activeDevice)[:28]...)
}

//export wolong_residual_in8
func wolong_residual_in8(port C.uint16_t) C.uint8_t {
	value := activeDevice.In8(uint16(port))
	cPortReads = append(cPortReads, byte(port), byte(uint16(port)>>8), value)
	return C.uint8_t(value)
}

//export wolong_startup_in8
func wolong_startup_in8(port C.uint16_t) C.uint8_t {
	value := activeDevice.In8(uint16(port))
	cPortReads = append(cPortReads, byte(port), byte(uint16(port)>>8), value)
	return C.uint8_t(value)
}

//export wolong_main_in8
func wolong_main_in8(port C.uint16_t) C.uint8_t {
	value := activeDevice.In8(uint16(port))
	cPortReads = append(cPortReads, byte(port), byte(uint16(port)>>8), value)
	return C.uint8_t(value)
}

//export wolong_input_in
func wolong_input_in(port C.uint16_t) C.uint8_t {
	value := activeDevice.In8(uint16(port))
	cPortReads = append(cPortReads, byte(port), byte(uint16(port)>>8), value)
	return C.uint8_t(value)
}

var _ = machine.StubSeg
