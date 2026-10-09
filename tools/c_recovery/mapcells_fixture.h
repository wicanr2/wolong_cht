#ifndef WOLONG_C_MAPCELLS_FIXTURE_H
#define WOLONG_C_MAPCELLS_FIXTURE_H
#include "vga_fixture.h"
typedef struct {unsigned calls;KiVgaFixture blocks[2];} KiMapcellsFixture;
static void mc_snapshot(KiMachine16 *m,uint16_t t,void *user) {
    KiMapcellsFixture *f=user;unsigned block=f->calls++/8192;
    assert(block<2);vg_snapshot(m,t,&f->blocks[block]);
}
static void mapcells_run(KiMachine16 *m,uint16_t t,KiMapcellsFixture *f) {
    KiEconomyHooks h={mc_snapshot,0,f};ki_mapcells_invoke(m,t,&h);
}
#endif
