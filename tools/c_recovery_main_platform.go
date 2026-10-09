//go:build matching_main

package main

/*
#include "/repo/tools/c_recovery/rng.h"
*/
import "C"

//export wolong_list_checkpoint
func wolong_list_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	listCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_main_checkpoint
func wolong_main_checkpoint(m *C.KiMachine16, t C.uint16_t) {
	mainCheckpoint(activeDevice, activeDOSDevice, uint16(t))
}

//export wolong_main_in8
func wolong_main_in8(port C.uint16_t) C.uint8_t {
	value := activeDevice.In8(uint16(port))
	cPortReads = append(cPortReads, byte(port), byte(uint16(port)>>8), value)
	return C.uint8_t(value)
}
