/* DOS/V KI.EXE IDA 0x1EC82..0x1ECE0, docs/spec/202-c-rng-seed.md.
 * Proven data flow from the byte-exact assembly baseline; native semantic C.
 */
#include "seed.h"

#ifndef KI_SEED_MUTATION
#define KI_SEED_MUTATION 0
#endif

void sub_1EC82(KiRngState *state, uint8_t hour_bcd, uint8_t minute_bcd,
              uint8_t second_bcd, uint8_t rtc_al) {
    for (unsigned n = 0; n < 256; ++n) state->table[n] = (uint8_t)n;
    uint8_t i = second_bcd, j = (uint8_t)(second_bcd + 1);
    for (unsigned n = 0; n < 256; ++n) {
        uint8_t value = state->table[i];
        state->table[i] = state->table[j];
        state->table[j] = value;
        i = (uint8_t)(i + (KI_SEED_MUTATION == 1 ? 0x4e : 0x4f));
        j = (uint8_t)(j + 0x89);
    }
    uint8_t hour_term = KI_SEED_MUTATION == 2 ? hour_bcd : (uint8_t)(hour_bcd << 2);
    state->c = (uint8_t)(rtc_al + second_bcd + minute_bcd + hour_term);
    state->s = (uint8_t)(state->c ^ (KI_SEED_MUTATION == 3 ? minute_bcd : second_bcd));
}

static void seed_write_word(KiMachine16 *m, uint16_t offset, uint16_t value) {
    uint32_t at = (((uint32_t)m->ss << 4) + offset) & 0xfffffu;
    m->memory[at] = (uint8_t)value;
    m->memory[(at + 1) & 0xfffffu] = (uint8_t)(value >> 8);
}

static uint16_t seed_pop(KiMachine16 *m) {
    uint32_t at = (((uint32_t)m->ss << 4) + m->sp) & 0xfffffu;
    uint16_t value = (uint16_t)(m->memory[at] | (uint16_t)m->memory[(at + 1) & 0xfffffu] << 8);
    m->sp = (uint16_t)(m->sp + 2);
    return value;
}

static void seed_push(KiMachine16 *m, uint16_t value) {
    m->sp = (uint16_t)(m->sp - 2);
    seed_write_word(m, m->sp, value);
}

void sub_1EC82_abi(KiMachine16 *m, uint8_t hour_bcd, uint8_t minute_bcd,
                  uint8_t second_bcd, uint8_t rtc_al) {
    seed_push(m, m->ds);
    seed_push(m, m->es);
    seed_push(m, m->ax);
    seed_push(m, m->bx);
    seed_push(m, m->dx);
    /* INT 1Ah's FLAGS/CS slots are later overwritten by seed AX/BX.
     * Its persistent return-IP word remains below them. */
    seed_write_word(m, (uint16_t)(m->sp - 6), 0xec9f);
    KiRngState state;
    sub_1EC82(&state, hour_bcd, minute_bcd, second_bcd, rtc_al);
    seed_push(m, state.c);
    seed_push(m, second_bcd);
    (void)seed_pop(m);
    (void)seed_pop(m);
    uint32_t at = (((uint32_t)m->cs << 4) + 0xecfc) & 0xfffffu;
    for (unsigned i = 0; i < sizeof(state); ++i) m->memory[at + i] = ((uint8_t *)&state)[i];
    uint8_t parity = state.s;
    parity ^= parity >> 4; parity ^= parity >> 2; parity ^= parity >> 1;
    uint16_t flags = 0;
    if (!(parity & 1)) flags |= 4;
    if (!state.s) flags |= 0x40;
    if (state.s & 0x80) flags |= 0x80;
    /* XOR's AF is undefined on hardware; zero matches this dosgolem model. */
    m->flags = (uint16_t)((m->flags & (uint16_t)~0x08d5u) | flags);
    m->cx = (uint16_t)((uint16_t)(uint8_t)(hour_bcd << 2) << 8 | minute_bcd);
    m->dx = seed_pop(m);
    m->bx = seed_pop(m);
    m->ax = seed_pop(m);
    m->es = seed_pop(m);
    m->ds = seed_pop(m);
    m->ip = seed_pop(m);
}
