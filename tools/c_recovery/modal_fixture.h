#ifndef WOLONG_C_MODAL_FIXTURE_H
#define WOLONG_C_MODAL_FIXTURE_H
#include "hourly_fixture.h"
typedef struct {KiHourlyFixture base;uint16_t origin_cs;uint8_t ui[8192][64];} KiModalFixture;
static void md_observed_ui(KiMachine16 *m,KiModalFixture *f,uint8_t *b) {
    uint16_t cs=f->origin_cs,far_cs=ec_word(m,cs,0x2082);
    for(unsigned i=0;i<32;++i) b[i]=ec_byte(m,cs,(uint16_t)(0x7000+i));
    for(unsigned i=0;i<8;++i) b[32+i]=ec_byte(m,far_cs,(uint16_t)(0x7004+i));
    for(unsigned i=0;i<16;++i) b[40+i]=ec_byte(m,cs,(uint16_t)(0x9882+i));
    b[56]=ec_byte(m,cs,0x2239);b[57]=ec_byte(m,cs,0x223a);b[58]=ec_byte(m,cs,0x20e);b[59]=ec_byte(m,cs,0x98a6);
    for(unsigned i=0;i<4;++i) b[60+i]=ec_byte(m,cs,(uint16_t)(0x9874+i));
}
static void md_snapshot(KiMachine16 *m,uint16_t t,void *user) {
    KiModalFixture *f=user;unsigned i=f->base.calls;hr_snapshot(m,t,&f->base);
    if(i<8192) {KiHourlySnapshot *s=&f->base.trace[i];for(unsigned j=0;j<59;++j) s->globals[j]=ec_byte(m,f->origin_cs,(uint16_t)(0xcf0+j));for(unsigned j=0;j<2;++j) s->rng[j]=ec_byte(m,f->origin_cs,(uint16_t)(0xecfc+j));s->cadence=ec_byte(m,f->origin_cs,0x31ad);md_observed_ui(m,f,f->ui[i]);}
}
static void md_external(KiMachine16 *m,uint16_t t,void *user) {
    KiModalFixture *f=user;KiEconomyHooks h={md_snapshot,md_external,user};uint16_t cs=f->origin_cs;
    if(t==0) {ec_cmp(m,m->ax,3,16);if(m->flags&0x40) {m->cx=ec_word(m,m->cs,0x7004);m->dx=ec_word(m,m->cs,0x7006);}ec_cmp(m,m->ax,m->ax,16);return;}
    if(t==0x9796) {gw_al(m,ec_byte(m,cs,0x700a));return;}
    if(t==0xe453) {gw_al(m,ec_byte(m,cs,0x7010));return;}
    if(t==0x93e9||t==0x7c6e) {
        uint16_t at=t==0x93e9?0x7008:0x7009;m->ax=ec_word(m,cs,t==0x93e9?0x7000:0x7002);uint8_t wait=ec_byte(m,cs,at);ec_cmp(m,wait,0,8);
        if(wait) {ec_put(m,cs,at,(uint8_t)st_dec(m,wait,8));m->flags|=1;}else m->flags&=(uint16_t)~1u;return;
    }
    if(!ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&!ki_politics_body(m,t,&h)&&!ki_world_body(m,t,&h)) (void)ki_settlement_body(m,t,&h);
}
static void modal_run(KiMachine16 *m,uint16_t t,KiModalFixture *f) {KiEconomyHooks h={md_snapshot,md_external,f};ki_modal_invoke(m,t,&h);}
#endif
