//go:build matching_vga

package main

/*
#include <stdint.h>
*/
import "C"
import "unsafe"

// activeDevice is the C side's independent dosgolem machine, not the original CPU.
//
//export wolong_vga_read
func wolong_vga_read(address C.uint32_t) C.uint8_t {
	return C.uint8_t(activeDevice.Read8(uint32(address)))
}

//export wolong_vga_write
func wolong_vga_write(address C.uint32_t, value C.uint8_t) {
	activeDevice.Write8(uint32(address), uint8(value))
}

//export wolong_vga_out
func wolong_vga_out(port C.uint16_t, value C.uint8_t) {
	activeDevice.Out8(uint16(port), uint8(value))
}

//export wolong_vga_state
func wolong_vga_state(buffer *C.uint8_t) {
	gc, seq, latch := activeDevice.VGAState()
	b := unsafe.Slice((*byte)(unsafe.Pointer(buffer)), 28)
	copy(b, gc[:])
	copy(b[16:], seq[:])
	copy(b[24:], latch[:])
}
