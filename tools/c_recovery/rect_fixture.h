#ifndef WOLONG_C_RECT_FIXTURE_H
#define WOLONG_C_RECT_FIXTURE_H
#include "vga_fixture.h"
typedef struct {KiVgaFixture trace;unsigned unsupported;} KiRectFixture;
static void rc_snapshot(KiMachine16 *m,uint16_t t,void *user) {KiRectFixture *f=user;vg_snapshot(m,t,&f->trace);}
static void rc_external(KiMachine16 *m,uint16_t t,void *user) {
    KiRectFixture *f=user;KiEconomyHooks h={rc_snapshot,rc_external,user};KiVgaHooks vh={rc_snapshot,user};
    if(t==0) {ec_cmp(m,m->ax,3,16);if(m->flags&0x40) {m->cx=ec_word(m,m->cs,0x7004);m->dx=ec_word(m,m->cs,0x7006);}ec_cmp(m,m->ax,2,16);if(m->flags&0x40) {m->ax=ec_word(m,m->cs,0x7010);m->dx=ec_word(m,m->cs,0x7012);m->bx=ec_word(m,m->cs,0x7014);}ec_cmp(m,m->ax,m->ax,16);return;}
    if(!ki_rect_body(m,t,&h)&&!ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)) f->unsupported=t;
}
static void rect_run(KiMachine16 *m,uint16_t t,KiRectFixture *f) {KiEconomyHooks h={rc_snapshot,rc_external,f};ki_rect_invoke(m,t,&h);}
#endif
