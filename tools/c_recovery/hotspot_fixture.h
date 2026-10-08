#ifndef WOLONG_C_HOTSPOT_FIXTURE_H
#define WOLONG_C_HOTSPOT_FIXTURE_H
#include "numeric_fixture.h"
typedef struct {KiNumericFixture base;uint8_t hotspot[8192][16];} KiHotspotFixture;
static void hs_observed(KiMachine16 *m,KiHotspotFixture *f,uint8_t *b) {
    uint16_t cs=f->base.base.origin_cs;for(unsigned i=0;i<4;++i) b[i]=ec_byte(m,cs,(uint16_t)(0xe479 + i));for(unsigned i=0;i<2;++i) {b[4+i]=ec_byte(m,cs,(uint16_t)(0xd84e + i));b[6+i]=ec_byte(m,cs,(uint16_t)(0xd4a + i));}uint16_t seg=ec_word(m,cs,0xe479),off=ec_word(m,cs,0xe47b);for(unsigned i=0;i<8;++i) b[8+i]=ec_byte(m,seg,(uint16_t)(off+i));
}
static void hs_snapshot(KiMachine16 *m,uint16_t t,void *user) {KiHotspotFixture *f=user;unsigned i=f->base.base.base.calls;nu_snapshot(m,t,&f->base);if(i<8192) hs_observed(m,f,f->hotspot[i]);}
static void hs_external(KiMachine16 *m,uint16_t t,void *user) {
    KiHotspotFixture *f=user;KiEconomyHooks h={hs_snapshot,hs_external,user};uint16_t cs=f->base.base.origin_cs;
    if(t==0) {ec_cmp(m,m->ax,3,16);if(m->flags&0x40) {m->cx=ec_word(m,m->cs,0x7004);m->dx=ec_word(m,m->cs,0x7006);}ec_cmp(m,m->ax,2,16);if(m->flags&0x40) {m->ax=ec_word(m,m->cs,0x7010);m->dx=ec_word(m,m->cs,0x7012);m->bx=ec_word(m,m->cs,0x7014);}ec_cmp(m,m->ax,m->ax,16);return;}
    if(t==0x9796) return;
    if(t==0x21e7) {ec_push(m,m->bx);m->bx=ec_word(m,cs,0x7020);m->bx=ev_shifts(m,m->bx,2,0);m->bx=ec_math(m,m->bx,ec_word(m,cs,0x7020),16,0,0);m->cx=ec_word(m,cs,(uint16_t)(m->bx+0x7100));m->dx=ec_word(m,cs,(uint16_t)(m->bx+0x7102));gw_al(m,ec_byte(m,cs,(uint16_t)(m->bx+0x7104)));ec_put(m,cs,0x7023,(uint8_t)m->ax);ec_store(m,cs,0x7020,ec_inc(m,ec_word(m,cs,0x7020),16));ec_cmp(m,(uint8_t)m->ax,1,8);m->flags^=1;m->bx=ec_pop(m);return;}
    if(t==0x9409) {m->ax=ec_word(m,cs,0x7000);ec_cmp(m,ec_byte(m,cs,0x7023),1,8);if(m->flags&0x40) m->ax=ec_word(m,cs,0x700c);uint8_t wait=ec_byte(m,cs,0x7008);ec_cmp(m,wait,0,8);if(wait) {ec_put(m,cs,0x7008,(uint8_t)st_dec(m,wait,8));m->flags|=1;}else m->flags&=(uint16_t)~1u;return;}
    if(!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&!ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&!ki_politics_body(m,t,&h)&&!ki_world_body(m,t,&h)) (void)ki_settlement_body(m,t,&h);
}
static void hotspot_run(KiMachine16 *m,uint16_t t,KiHotspotFixture *f) {KiEconomyHooks h={hs_snapshot,hs_external,f};ki_hotspot_invoke(m,t,&h);}
#endif
