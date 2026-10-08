#ifndef WOLONG_WORLD_FIXTURE_H
#define WOLONG_WORLD_FIXTURE_H
#include "world_update.h"
#include <string.h>
typedef struct {uint16_t target,regs[14];uint8_t record[64],stack[32],globals[59],rng[2];} KiWorldSnapshot;
typedef struct {unsigned calls;KiWorldSnapshot trace[4096];} KiWorldFixture;
static void w_snapshot(KiMachine16 *m,uint16_t target,void *user) {
    KiWorldFixture *f=user;if(f->calls<4096) {
        KiWorldSnapshot *s=&f->trace[f->calls];s->target=target;
        uint16_t regs[]={m->ax,m->bx,m->cx,m->dx,m->si,m->di,m->bp,m->sp,m->ds,m->es,m->ss,m->cs,m->ip,m->flags};memcpy(s->regs,regs,sizeof(regs));
        for(unsigned i=0;i<64;++i) s->record[i]=ec_byte(m,m->ds,(uint16_t)(m->si+i));
        for(unsigned i=0;i<32;++i) s->stack[i]=ec_byte(m,m->ss,(uint16_t)(m->sp+i));
        for(unsigned i=0;i<59;++i) s->globals[i]=ec_byte(m,m->cs,(uint16_t)(0xcf0+i));
        for(unsigned i=0;i<2;++i) s->rng[i]=ec_byte(m,m->cs,(uint16_t)(0xecfc+i));
    }++f->calls;
}
static void w_external(KiMachine16 *m,uint16_t target,void *user) {
    KiEconomyHooks hooks={w_snapshot,w_external,user};
    if(!ki_world_body(m,target,&hooks)) (void)ki_settlement_body(m,target,&hooks);
}
static void world_run(KiMachine16 *m,uint16_t target,KiWorldFixture *f) {
    KiEconomyHooks hooks={w_snapshot,w_external,f};ki_world_invoke(m,target,&hooks);
}
#endif
