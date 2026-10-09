//go:build matching_engagement

package main

/*
#include "/repo/tools/c_recovery/rng.h"
*/
import "C"
import "github.com/wicanr2/dosgolem/internal/machine"

//export wolong_list_checkpoint
func wolong_list_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	listCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_main_checkpoint
func wolong_main_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	engagementCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_main_in8
func wolong_main_in8(port C.uint16_t) C.uint8_t {
	value := activeDevice.In8(uint16(port))
	cPortReads = append(cPortReads, byte(port), byte(uint16(port)>>8), value)
	return C.uint8_t(value)
}

//export wolong_bootstrap_checkpoint
func wolong_bootstrap_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	engagementCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_interaction_checkpoint
func wolong_interaction_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	engagementCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_tick_checkpoint
func wolong_tick_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	engagementCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_route_checkpoint
func wolong_route_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	engagementCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_outcome_checkpoint
func wolong_outcome_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	engagementCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_engagement_checkpoint
func wolong_engagement_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	engagementCheckpoint(activeDevice, activeDOSDevice, uint16(t))
	engagementObserveSave(activeDevice, uint16(t), cRegs(m))
}

//export wolong_engagement_snapshot
func wolong_engagement_snapshot(m *C.KiMachine16, t C.uint16_t) {
	engagementCCalls++
	if engagementCCalls > 2000000 {
		panic("bounded native engagement entry count")
	}
	cur := cRegs(m)
	cur.IP = uint16(t)
	loc := uint32(cur.IP) + 0x10000
	if cur.CS == engagementMouseCS {
		loc = uint32(cur.IP) + 0x20000
	}
	_, known := functions[loc]
	if !known && !(cur.CS == engagementCoreCS && cur.IP == 0x2259) && !(cur.CS == machine.StubSeg && (cur.IP == 0x410 || cur.IP == 0x414)) {
		return
	}
	engagementCTrace = append(engagementCTrace, byte(t), byte(uint16(t)>>8))
	engagementCTrace = append(engagementCTrace, regBytes(cur)...)
	for i := uint32(0); i < 32; i++ {
		engagementCTrace = append(engagementCTrace, activeDevice.Mem[(uint32(cur.SS)*16+uint32(cur.SP)+i)&0xfffff])
	}
	engagementCTrace = append(engagementCTrace, deviceBytes(activeDevice)[:28]...)
}
