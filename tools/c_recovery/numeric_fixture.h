#ifndef WOLONG_C_NUMERIC_FIXTURE_H
#define WOLONG_C_NUMERIC_FIXTURE_H
#include "modal_fixture.h"
typedef struct {KiModalFixture base;uint8_t input[8192][24];} KiNumericFixture;
static void nu_observed(KiMachine16 *m,KiNumericFixture *f,uint8_t *b) {
    uint16_t cs=f->base.origin_cs,far_cs=ec_word(m,cs,0x2082);for(unsigned i=0;i<8;++i) b[i]=ec_byte(m,cs,(uint16_t)(0x7020+i));for(unsigned i=0;i<6;++i) b[8+i]=ec_byte(m,cs,(uint16_t)(0x1d5+i));b[14]=ec_byte(m,cs,0x9441);b[15]=ec_byte(m,cs,0x9443);for(unsigned i=0;i<6;++i) b[16+i]=ec_byte(m,far_cs,(uint16_t)(0x7010+i));b[22]=ec_byte(m,m->ss,m->bp);b[23]=ec_byte(m,m->ss,(uint16_t)(m->bp+1));
}
static void nu_snapshot(KiMachine16 *m,uint16_t t,void *user) {KiNumericFixture *f=user;unsigned i=f->base.base.calls;md_snapshot(m,t,&f->base);if(i<8192) nu_observed(m,f,f->input[i]);}
static void nu_external(KiMachine16 *m,uint16_t t,void *user) {
    KiNumericFixture *f=user;KiEconomyHooks h={nu_snapshot,nu_external,user};uint16_t cs=f->base.origin_cs;
    if(t==0) {ec_cmp(m,m->ax,3,16);if(m->flags&0x40) {m->cx=ec_word(m,m->cs,0x7004);m->dx=ec_word(m,m->cs,0x7006);}ec_cmp(m,m->ax,2,16);if(m->flags&0x40) {m->ax=ec_word(m,m->cs,0x7010);m->dx=ec_word(m,m->cs,0x7012);m->bx=ec_word(m,m->cs,0x7014);}ec_cmp(m,m->ax,m->ax,16);return;}
    if(t==0x9796) return;
    if(t==0x21e7) {ec_push(m,m->bx);m->bx=ec_word(m,cs,0x7020);m->bx=ec_math(m,m->bx,m->bx,16,0,0);gw_al(m,ec_byte(m,cs,(uint16_t)(m->bx+0x7100)));ec_put(m,cs,0x7022,(uint8_t)m->ax);gw_al(m,ec_byte(m,cs,(uint16_t)(m->bx+0x7101)));ec_put(m,cs,0x7023,(uint8_t)m->ax);ec_store(m,cs,0x7020,ec_inc(m,ec_word(m,cs,0x7020),16));ec_cmp(m,(uint8_t)m->ax,1,8);m->flags^=1;m->bx=ec_pop(m);return;}
    if(t==0xe453) {gw_al(m,ec_byte(m,cs,0x7022));return;}
    if(t==0x9409) {m->ax=ec_word(m,cs,0x7000);ec_cmp(m,ec_byte(m,cs,0x7023),1,8);if(m->flags&0x40) m->ax=ec_word(m,cs,0x700c);uint8_t wait=ec_byte(m,cs,0x7008);ec_cmp(m,wait,0,8);if(wait) {ec_put(m,cs,0x7008,(uint8_t)st_dec(m,wait,8));m->flags|=1;}else m->flags&=(uint16_t)~1u;return;}
    if(t==0x7c6e) {ki_numeric_body(m,t,&h);return;}
    if(!ki_numeric_body(m,t,&h)&&!ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&!ki_politics_body(m,t,&h)&&!ki_world_body(m,t,&h)) (void)ki_settlement_body(m,t,&h);
}
static void numeric_run(KiMachine16 *m,uint16_t t,KiNumericFixture *f) {KiEconomyHooks h={nu_snapshot,nu_external,f};ki_numeric_invoke(m,t,&h);}
#endif
