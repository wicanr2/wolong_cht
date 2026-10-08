#ifndef WOLONG_C_NUMBERS_FIXTURE_H
#define WOLONG_C_NUMBERS_FIXTURE_H
#include "vga_fixture.h"
typedef struct {KiVgaFixture trace;unsigned unsupported;} KiNumbersFixture;
static void nr_snapshot(KiMachine16 *m,uint16_t t,void *user) {
    KiNumbersFixture *f=user;vg_snapshot(m,t,&f->trace);
}
static void nr_external(KiMachine16 *m,uint16_t t,void *user) {
    KiNumbersFixture *f=user;f->unsupported=t;
}
static void numbers_run(KiMachine16 *m,uint16_t t,KiNumbersFixture *f) {
    KiEconomyHooks h={nr_snapshot,nr_external,f};ki_numbers_invoke(m,t,&h);
}
#endif
