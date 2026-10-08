#ifndef WOLONG_C_ALIGNED_FIXTURE_H
#define WOLONG_C_ALIGNED_FIXTURE_H
#include "vga_fixture.h"
typedef struct {KiVgaFixture trace;unsigned unsupported;} KiAlignedFixture;
static void al_external(KiMachine16 *m,uint16_t t,void *user) {
    KiAlignedFixture *f=user;KiVgaHooks vh={vg_snapshot,&f->trace};
    if(t==0) {
        ec_cmp(m,m->ax,2,16);if(m->flags&0x40) {m->ax=ec_word(m,m->cs,0x7010);m->dx=ec_word(m,m->cs,0x7012);m->bx=ec_word(m,m->cs,0x7014);}ec_cmp(m,m->ax,m->ax,16);return;
    }
    if(!ki_vga_body(m,t,&vh)) f->unsupported=t;
}
static void al_snapshot(KiMachine16 *m,uint16_t t,void *user) {KiAlignedFixture *f=user;vg_snapshot(m,t,&f->trace);}
static void aligned_run(KiMachine16 *m,uint16_t t,KiAlignedFixture *f) {KiEconomyHooks h={al_snapshot,al_external,f};ki_aligned_invoke(m,t,&h);}
#endif
