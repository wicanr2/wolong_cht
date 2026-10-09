//go:build matching_march

package main

/*
#include "/repo/tools/c_recovery/rng.h"
*/
import "C"

//export wolong_list_checkpoint
func wolong_list_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	listCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_march_checkpoint
func wolong_march_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	marchCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}
