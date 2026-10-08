#ifndef WOLONG_C_RECOVERY_RNG_H
#define WOLONG_C_RECOVERY_RNG_H
#include <stdint.h>

/* proven: DOS/V KI.EXE IDA 0x1ECFC / 0x1ECFD / 0x1ECFE; docs/re/10. */
typedef struct {
    uint8_t c, s, table[256];
} KiRngState;

/* Research adapter only. These names describe the x86 ABI, not recovered source names. */
typedef struct {
    uint16_t ax, bx, cx, dx, si, di, bp, sp, ds, es, ss, cs, ip, flags;
    volatile uint8_t *memory;
} KiMachine16;

uint8_t sub_1ECE0(KiRngState *state);
void sub_1ECE0_abi(KiMachine16 *machine);
#endif
