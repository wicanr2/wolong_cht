//go:build matching_input

package main

/*
#include "/repo/tools/c_recovery/rng.h"
*/
import "C"

// Platform bus and explicit timer input only; no guest CPU execution.
//
//export wolong_input_in
func wolong_input_in(port C.uint16_t) C.uint8_t {
	value := activeDevice.In8(uint16(port))
	cPortReads = append(cPortReads, byte(port), byte(uint16(port)>>8), value)
	return C.uint8_t(value)
}

//export wolong_input_checkpoint
func wolong_input_checkpoint(m *C.KiMachine16, target C.uint16_t) {
	activeTicks++
	value := byte(0xff)
	if activeTicks <= currentInput.timerHold {
		value = 0xfe
	}
	activeDevice.Mem[int(m.cs)*16+0xd2d] = value
}
