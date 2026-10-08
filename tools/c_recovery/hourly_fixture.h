#ifndef WOLONG_HOURLY_FIXTURE_H
#define WOLONG_HOURLY_FIXTURE_H
#include "hourly.h"
#include <string.h>
typedef struct {uint16_t target,regs[14];uint8_t record[64],stack[32],globals[59],rng[2],cadence;} KiHourlySnapshot;
typedef struct {unsigned calls;KiHourlySnapshot trace[8192];} KiHourlyFixture;
static void hr_snapshot(KiMachine16 *m,uint16_t target,void *user) {
    KiHourlyFixture *f=user;if(f->calls<8192) {KiHourlySnapshot *s=&f->trace[f->calls];s->target=target;
        uint16_t regs[]={m->ax,m->bx,m->cx,m->dx,m->si,m->di,m->bp,m->sp,m->ds,m->es,m->ss,m->cs,m->ip,m->flags};memcpy(s->regs,regs,sizeof(regs));
        for(unsigned i=0;i<64;++i) s->record[i]=ec_byte(m,m->ds,(uint16_t)(m->si+i));
        for(unsigned i=0;i<32;++i) s->stack[i]=ec_byte(m,m->ss,(uint16_t)(m->sp+i));
        for(unsigned i=0;i<59;++i) s->globals[i]=ec_byte(m,m->cs,(uint16_t)(0xcf0+i));
        for(unsigned i=0;i<2;++i) s->rng[i]=ec_byte(m,m->cs,(uint16_t)(0xecfc+i));s->cadence=ec_byte(m,m->cs,0x31ad);
    }++f->calls;
}
static void hr_external(KiMachine16 *m,uint16_t target,void *user) {
    KiEconomyHooks h={hr_snapshot,hr_external,user};
    if(!ki_hourly_body(m,target,&h)&&!ki_politics_body(m,target,&h)&&!ki_world_body(m,target,&h)) (void)ki_settlement_body(m,target,&h);
}
static void hourly_run(KiMachine16 *m,uint16_t target,KiHourlyFixture *f) {KiEconomyHooks h={hr_snapshot,hr_external,f};ki_hourly_invoke(m,target,&h);}
#endif
