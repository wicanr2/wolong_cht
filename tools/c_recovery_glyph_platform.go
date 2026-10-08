//go:build matching_glyph

package main

/*
#include "/repo/tools/c_recovery/rng.h"
*/
import "C"

// Platform register carrier only: no C-side CPU.Step or guest instruction fetch.
//
//export wolong_glyph_interrupt
func wolong_glyph_interrupt(m *C.KiMachine16, number C.uint8_t) C.int {
	r := cRegs(m)
	c := activeDevice.CPU
	c.R = [8]uint16{r.AX, r.CX, r.DX, r.BX, r.SP, r.BP, r.SI, r.DI}
	c.Seg = [4]uint16{r.ES, r.CS, r.SS, r.DS}
	c.IP = r.IP
	c.SetFlags(r.Flags)
	if c.IntHook == nil || !c.IntHook(c, uint8(number)) {
		return 0
	}
	assignC(m, originalRegs(activeDevice))
	return 1
}
