#ifndef WOLONG_CLOCK_FIXTURE_H
#define WOLONG_CLOCK_FIXTURE_H
#include "clock.h"
#include <string.h>

typedef struct {
    uint16_t target, ax, bx, cx, dx, si, di, bp, sp, ds, es, ss, cs, ip, flags;
    uint8_t date[8];
} KiClockSnapshot;
typedef struct {
    unsigned alter, delayed, ready_delay, count_delay, ready_polls, count_polls, calls;
    uint8_t target_count;
    KiClockSnapshot trace[4];
} KiClockFixture;

static uint32_t fixture_at(KiMachine16 *m, uint16_t offset) {
    return (((uint32_t)m->ds << 4) + offset) & 0xfffffu;
}
static void fixture_call(KiMachine16 *m, uint16_t target, void *user) {
    KiClockFixture *f = user;
    if (f->calls < 4) {
        KiClockSnapshot *s = &f->trace[f->calls];
        s->target=target;s->ax=m->ax;s->bx=m->bx;s->cx=m->cx;s->dx=m->dx;
        s->si=m->si;s->di=m->di;s->bp=m->bp;s->sp=m->sp;s->ds=m->ds;s->es=m->es;
        s->ss=m->ss;s->cs=m->cs;s->ip=m->ip;s->flags=m->flags;
        for (unsigned i=0;i<8;++i) s->date[i]=m->memory[fixture_at(m,(uint16_t)(0xcf0+i))];
    }
    ++f->calls;
    if (f->alter) {
        if (target==0x5358) m->memory[fixture_at(m,0xcf0)]=7;
        if (target==0x9377) m->memory[fixture_at(m,0xcfa)]=2;
        if (target==0x3e11) m->ax=0xa5c3;
        if (target==0x1e17) m->bx=0xb7c0;
    }
}
static void fixture_poll(KiMachine16 *m, uint16_t offset, void *user) {
    KiClockFixture *f = user;
    if (offset==0xd2c) {
        ++f->ready_polls;
        if (f->delayed && f->ready_polls>f->ready_delay) m->memory[fixture_at(m,offset)]=1;
    } else {
        ++f->count_polls;
        if (f->delayed && f->count_polls>f->count_delay) m->memory[fixture_at(m,offset)]=f->target_count;
        if (f->count_polls>64) m->memory[fixture_at(m,offset)]=(uint8_t)(f->target_count+1);
    }
}
static void clock_run_fixture(KiMachine16 *m, KiClockFixture *f) {
    KiClockHooks hooks={fixture_call,fixture_poll,f};sub_11D8E(m,&hooks);
}
#endif
