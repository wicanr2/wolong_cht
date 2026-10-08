#ifndef WOLONG_ECONOMY_FIXTURE_H
#define WOLONG_ECONOMY_FIXTURE_H
#include "economy.h"
#include <string.h>

typedef struct {
    uint16_t target,regs[14];
    uint8_t record[64],settings[8],rng[2];
} KiEconomySnapshot;
typedef struct {
    unsigned calls;
    uint16_t income_low;
    uint8_t income_high;
    KiEconomySnapshot trace[256];
} KiEconomyFixture;
static void economy_snapshot(KiMachine16 *m,uint16_t target,void *user) {
    KiEconomyFixture *f=user;
    if (f->calls<256) {
        KiEconomySnapshot *s=&f->trace[f->calls];s->target=target;
        uint16_t regs[]={m->ax,m->bx,m->cx,m->dx,m->si,m->di,m->bp,m->sp,m->ds,m->es,m->ss,m->cs,m->ip,m->flags};
        memcpy(s->regs,regs,sizeof(regs));
        for (unsigned i=0;i<64;++i) s->record[i]=m->memory[ec_at(m->ds,(uint16_t)(m->si+i))];
        for (unsigned i=0;i<8;++i) s->settings[i]=m->memory[ec_at(m->cs,(uint16_t)(0xd08+i))];
        for (unsigned i=0;i<2;++i) s->rng[i]=m->memory[ec_at(m->cs,(uint16_t)(0xecfc+i))];
    }
    ++f->calls;
}
static void economy_external(KiMachine16 *m,uint16_t target,void *user) {
    KiEconomyFixture *f=user;
    if (target==0x53c6) { m->ax=f->income_low;m->dx=(uint16_t)((m->dx&0xff00u)|f->income_high); }
}
static void economy_run_fixture(KiMachine16 *m,uint16_t target,KiEconomyFixture *f) {
    KiEconomyHooks hooks={economy_snapshot,economy_external,f};ki_economy_invoke(m,target,&hooks);
}
#endif
