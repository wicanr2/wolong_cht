/* DOS/V KI.EXE IDA 0x11D8E..0x11E17; docs/spec/203-c-game-clock.md.
 * Native C semantics with explicit callee and wait inputs; not a matching C binary.
 */
#include "clock.h"
#ifndef KI_CLOCK_MUTATION
#define KI_CLOCK_MUTATION 0
#endif

static uint32_t clock_at(uint16_t segment, uint16_t offset) {
    return (((uint32_t)segment << 4) + offset) & 0xfffffu;
}
static uint8_t clock_byte(KiMachine16 *m, uint16_t offset) { return m->memory[clock_at(m->ds, offset)]; }
static void clock_put(KiMachine16 *m, uint16_t offset, uint8_t value) { m->memory[clock_at(m->ds, offset)] = value; }
static uint16_t clock_word(KiMachine16 *m, uint16_t offset) {
    return (uint16_t)(clock_byte(m, offset) | (uint16_t)clock_byte(m, (uint16_t)(offset + 1)) << 8);
}
static void clock_put_word(KiMachine16 *m, uint16_t offset, uint16_t value) {
    clock_put(m, offset, (uint8_t)value);clock_put(m, (uint16_t)(offset + 1), (uint8_t)(value >> 8));
}
static void clock_stack(KiMachine16 *m, uint16_t value) {
    m->sp = (uint16_t)(m->sp - 2);uint32_t at = clock_at(m->ss, m->sp);
    m->memory[at] = (uint8_t)value;m->memory[(at + 1) & 0xfffffu] = (uint8_t)(value >> 8);
}
static uint16_t clock_pop(KiMachine16 *m) {
    uint32_t at = clock_at(m->ss, m->sp);
    uint16_t value = (uint16_t)(m->memory[at] | (uint16_t)m->memory[(at + 1) & 0xfffffu] << 8);
    m->sp = (uint16_t)(m->sp + 2);return value;
}
static uint16_t clock_szp(uint16_t flags, uint16_t result, unsigned width) {
    uint8_t parity = (uint8_t)result;parity ^= parity >> 4;parity ^= parity >> 2;parity ^= parity >> 1;
    flags &= (uint16_t)~0x00c4u;
    if (!(parity & 1)) flags |= 4;
    if (!result) flags |= 0x40;
    if (result & (width == 8 ? 0x80 : 0x8000)) flags |= 0x80;
    return flags;
}
static void clock_cmp(KiMachine16 *m, uint16_t a, uint16_t b, unsigned width) {
    uint16_t result = width == 8 ? (uint8_t)(a - b) : (uint16_t)(a - b);
    uint16_t bits = 0;unsigned sign = width == 8 ? 0x80 : 0x8000;
    if (a < b) bits |= 1;
    if ((a ^ b ^ result) & 0x10) bits |= 0x10;
    if (((a ^ b) & (a ^ result)) & sign) bits |= 0x800;
    m->flags = clock_szp((uint16_t)((m->flags & (uint16_t)~0x08d5u) | bits), result, width);
}
static uint16_t clock_inc(KiMachine16 *m, uint16_t a, unsigned width) {
    uint16_t result = width == 8 ? (uint8_t)(a + 1) : (uint16_t)(a + 1);
    uint16_t bits = m->flags & 1;
    if ((a ^ 1 ^ result) & 0x10) bits |= 0x10;
    if (a == (width == 8 ? 0x7f : 0x7fff)) bits |= 0x800;
    m->flags = clock_szp((uint16_t)((m->flags & (uint16_t)~0x08d5u) | bits), result, width);return result;
}
static void clock_logic(KiMachine16 *m, uint8_t value) {
    /* AF is undefined for AND/XOR on hardware; zero follows this dosgolem model. */
    m->flags = clock_szp((uint16_t)(m->flags & (uint16_t)~0x08d5u), value, 8);
}
static void clock_call(KiMachine16 *m, const KiClockHooks *hooks, uint16_t target, uint16_t next) {
    clock_stack(m, next);m->ip = target;
    if (hooks && hooks->call) hooks->call(m, target, hooks->user);
    m->ip = clock_pop(m);
}

void sub_11D8E(KiMachine16 *m, const KiClockHooks *hooks) {
    uint8_t subtick = clock_byte(m, 0xcf2);clock_cmp(m, subtick, 8, 8);
    if (subtick < 8) {
        clock_put(m, 0xcf2, (uint8_t)clock_inc(m, subtick, 8));
    } else {
        uint8_t hour = clock_byte(m, 0xcf3);
        clock_cmp(m, hour, KI_CLOCK_MUTATION == 1 ? 24 : 23, 8);
        if (hour >= (KI_CLOCK_MUTATION == 1 ? 24 : 23)) {
            m->ax = clock_word(m, 0xcf0);
            uint8_t day = (uint8_t)m->ax, days = (uint8_t)(m->ax >> 8);
            clock_cmp(m, day, days, 8);
            if (day >= days) {
                uint8_t month = clock_byte(m, 0xcf4);clock_cmp(m, month, 12, 8);
                if (month >= 12) {
                    uint16_t year = clock_word(m, 0xcf6);
                    clock_cmp(m, year, KI_CLOCK_MUTATION == 2 ? 999 : 1000, 16);
                    if (year >= (KI_CLOCK_MUTATION == 2 ? 999 : 1000)) clock_put_word(m, 0xcf6, 998);
                    clock_put_word(m, 0xcf6, clock_inc(m, clock_word(m, 0xcf6), 16));
                    clock_put(m, 0xcf4, 0);
                }
                clock_put(m, 0xcf4, (uint8_t)clock_inc(m, clock_byte(m, 0xcf4), 8));
                m->bx = clock_byte(m, 0xcf4);clock_logic(m, 0);
                m->ax = (uint16_t)clock_byte(m, (uint16_t)(m->bx + 0x98ab)) << 8;
                clock_logic(m, 0);clock_put_word(m, 0xcf0, m->ax);
                clock_call(m, hooks, 0x5358, 0x1dd7);
            }
            clock_put(m, 0xcf0, (uint8_t)clock_inc(m, clock_byte(m, 0xcf0), 8));
            clock_put(m, 0xcf3, 0);
        }
        clock_put(m, 0xcf3, (uint8_t)clock_inc(m, clock_byte(m, 0xcf3), 8));
        clock_put(m, 0xcf2, 0);
        if (KI_CLOCK_MUTATION == 3) {
            clock_call(m, hooks, 0x3e11, 0x1def);clock_call(m, hooks, 0x9377, 0x1dec);
        } else {
            clock_call(m, hooks, 0x9377, 0x1dec);clock_call(m, hooks, 0x3e11, 0x1def);
        }
        clock_call(m, hooks, 0x1e17, 0x1df2);
    }
    uint8_t speed = clock_byte(m, 0xcfa);m->ax = (uint16_t)((m->ax & 0xff00u) | speed);clock_logic(m, speed);
    if (speed) {
        do {
            m->ip = 0x1dff;if (hooks && hooks->poll) hooks->poll(m, 0xd2c, hooks->user);
            clock_cmp(m, clock_byte(m, 0xd2c), 0, 8);
        } while (!clock_byte(m, 0xd2c));
        for (;;) {
            m->ip = 0x1e06;if (hooks && hooks->poll) hooks->poll(m, 0xd2d, hooks->user);
            uint8_t count = clock_byte(m, 0xd2d);clock_cmp(m, count, speed, 8);
            if (KI_CLOCK_MUTATION == 4 ? count > speed : count >= speed) break;
        }
        clock_put(m, 0xd2d, 0);clock_put(m, 0xd2c, 0);
    }
    m->ip = clock_pop(m);
}
