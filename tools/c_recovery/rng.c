/* C restoration of DOS/V KI.EXE sub_1ECE0, IDA linear 0x1ECE0..0x1ECFC.
 * Input SHA-256: fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868.
 * Evidence: docs/re/10-rng.md; contract: docs/spec/201-c-rng-function.md.
 * Semantic reconstruction, not original source or a matching DOS C build.
 */
#include "rng.h"
#include <stddef.h>
#include <string.h>

#ifndef KI_MUTATION
#define KI_MUTATION 0
#endif
_Static_assert(sizeof(KiRngState) == 258, "RNG raw state size");
_Static_assert(offsetof(KiRngState, table) == 2, "RNG table offset");

uint8_t sub_1ECE0(KiRngState *state) {
    uint8_t value = (uint8_t)(state->table[state->s] + state->c);
    state->c = (uint8_t)(state->c + (KI_MUTATION == 1 ? 0x88 : 0x89));
    state->s = value;
    return value;
}

static uint32_t physical(uint16_t segment, uint16_t offset) {
    return (((uint32_t)segment << 4) + offset) & 0xfffffu;
}

static void write_word(KiMachine16 *m, uint16_t offset, uint16_t value) {
    uint32_t at = physical(m->ss, offset);
    m->memory[at] = (uint8_t)value;
    m->memory[(at + 1) & 0xfffffu] = (uint8_t)(value >> 8);
}

static uint16_t read_word(KiMachine16 *m, uint16_t offset) {
    uint32_t at = physical(m->ss, offset);
    return (uint16_t)(m->memory[at] | (uint16_t)m->memory[(at + 1) & 0xfffffu] << 8);
}

/* Intel SDM Vol. 2A ADD: CF/PF/AF/ZF/SF/OF; no game-specific flag guess. */
static uint16_t add8_flags(uint16_t flags, uint8_t a, uint8_t b) {
    uint16_t sum = (uint16_t)a + b;
    uint8_t result = (uint8_t)sum, parity = result;
    uint16_t bits = 0;
    parity ^= parity >> 4;
    parity ^= parity >> 2;
    parity ^= parity >> 1;
    if (sum > 255) bits |= 0x0001;
    if (!(parity & 1)) bits |= 0x0004;
    if ((a ^ b ^ result) & 0x10) bits |= 0x0010;
    if (!result) bits |= 0x0040;
    if (result & 0x80) bits |= 0x0080;
    if ((~(a ^ b) & (a ^ result)) & 0x80) bits |= 0x0800;
    return (uint16_t)((flags & (uint16_t)~0x08d5u) | bits);
}

void sub_1ECE0_abi(KiMachine16 *m) {
    uint16_t saved_bx = m->bx;
    KiRngState state;
    m->sp = (uint16_t)(m->sp - 2);
    write_word(m, m->sp, m->ds);
    m->sp = (uint16_t)(m->sp - 2);
    write_word(m, m->sp, m->bx);
    uint32_t at = physical(m->cs, 0xecfc);
    for (unsigned i = 0; i < sizeof(state); ++i)
        ((uint8_t *)&state)[i] = m->memory[at + i];
    uint8_t old_counter = state.c, old_sample = state.table[state.s];
    uint8_t value = sub_1ECE0(&state);
    m->memory[at] = state.c;
    m->memory[at + 1] = state.s;
    m->ax = (uint16_t)((KI_MUTATION == 2 ? 0 : m->cs & 0xff00u) | value);
    m->flags = KI_MUTATION == 3 ? add8_flags(m->flags, old_sample, old_counter)
                                : add8_flags(m->flags, old_counter, 0x89);
    m->bx = KI_MUTATION == 4 ? saved_bx : read_word(m, m->sp);
    m->sp = (uint16_t)(m->sp + 2);
    m->ds = read_word(m, m->sp);
    m->sp = (uint16_t)(m->sp + 2);
    m->ip = read_word(m, m->sp);
    m->sp = (uint16_t)(m->sp + 2);
}
