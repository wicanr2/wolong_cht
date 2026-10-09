//go:build matching_personnel

package main

/*
#include "/repo/tools/c_recovery/rng.h"
*/
import "C"

//export wolong_list_checkpoint
func wolong_list_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	listCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_personnel_entry
func wolong_personnel_entry(m *C.KiMachine16, t C.uint16_t) { personnelEntry(activeDevice, uint16(t)) }
