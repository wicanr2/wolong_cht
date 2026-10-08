#ifndef WOLONG_POLITICS_FIXTURE_H
#define WOLONG_POLITICS_FIXTURE_H
#include "politics.h"
#include <string.h>
typedef struct {uint16_t target,regs[14];uint8_t record[64],stack[32],globals[59],rng[2];} KiPoliticsSnapshot;
typedef struct {unsigned calls;KiPoliticsSnapshot trace[8192];} KiPoliticsFixture;
static void pn_snapshot(KiMachine16 *m,uint16_t target,void *user) {
    KiPoliticsFixture *f=user;if(f->calls<8192) {
        KiPoliticsSnapshot *s=&f->trace[f->calls];s->target=target;
        uint16_t regs[]={m->ax,m->bx,m->cx,m->dx,m->si,m->di,m->bp,m->sp,m->ds,m->es,m->ss,m->cs,m->ip,m->flags};memcpy(s->regs,regs,sizeof(regs));
        for(unsigned i=0;i<64;++i) s->record[i]=ec_byte(m,m->ds,(uint16_t)(m->si+i));
        for(unsigned i=0;i<32;++i) s->stack[i]=ec_byte(m,m->ss,(uint16_t)(m->sp+i));
        for(unsigned i=0;i<59;++i) s->globals[i]=ec_byte(m,m->cs,(uint16_t)(0xcf0+i));
        for(unsigned i=0;i<2;++i) s->rng[i]=ec_byte(m,m->cs,(uint16_t)(0xecfc+i));
    }++f->calls;
}
static void pn_external(KiMachine16 *m,uint16_t target,void *user) {
    KiEconomyHooks h={pn_snapshot,pn_external,user};
    if(!ki_politics_body(m,target,&h)&&!ki_world_body(m,target,&h)) (void)ki_settlement_body(m,target,&h);
}
static void politics_run(KiMachine16 *m,uint16_t target,KiPoliticsFixture *f) {KiEconomyHooks h={pn_snapshot,pn_external,f};ki_politics_invoke(m,target,&h);}
#endif
