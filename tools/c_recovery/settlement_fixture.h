#ifndef WOLONG_SETTLEMENT_FIXTURE_H
#define WOLONG_SETTLEMENT_FIXTURE_H
#include "settlement.h"
#include <string.h>
typedef struct {
    uint16_t target,regs[14];
    uint8_t record[64],locals[10],display[22],rng[2];
} KiSettlementSnapshot;
typedef struct {unsigned calls;KiSettlementSnapshot trace[1024];} KiSettlementFixture;
static void st_snapshot(KiMachine16 *m,uint16_t target,void *user) {
    KiSettlementFixture *f=user;
    if (f->calls<1024) {
        KiSettlementSnapshot *s=&f->trace[f->calls];s->target=target;
        uint16_t regs[]={m->ax,m->bx,m->cx,m->dx,m->si,m->di,m->bp,m->sp,m->ds,m->es,m->ss,m->cs,m->ip,m->flags};memcpy(s->regs,regs,sizeof(regs));
        for (unsigned i=0;i<64;++i) s->record[i]=ec_byte(m,m->ds,(uint16_t)(m->si+i));
        for (unsigned i=0;i<10;++i) s->locals[i]=ec_byte(m,m->ss,(uint16_t)(m->bp+i));
        for (unsigned i=0;i<22;++i) s->display[i]=ec_byte(m,m->cs,(uint16_t)(0xd02+i));
        for (unsigned i=0;i<2;++i) s->rng[i]=ec_byte(m,m->cs,(uint16_t)(0xecfc+i));
    }
    ++f->calls;
}
static void st_external(KiMachine16 *m,uint16_t target,void *user) {
    KiEconomyHooks hooks={st_snapshot,st_external,user};
    /* Existing economy dispatcher owns RET. Execute new body only, without an extra pop. */
    (void)ki_settlement_body(m,target,&hooks);
}
static void settlement_run(KiMachine16 *m,uint16_t target,KiSettlementFixture *f) {
    KiEconomyHooks hooks={st_snapshot,st_external,f};ki_settlement_invoke(m,target,&hooks);
}
#endif
