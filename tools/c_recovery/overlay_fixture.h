#ifndef WOLONG_C_OVERLAY_FIXTURE_H
#define WOLONG_C_OVERLAY_FIXTURE_H
#include "mapcells_fixture.h"
static void overlay_external(KiMachine16 *m,uint16_t t,void *user) {
    KiEconomyHooks h={mc_snapshot,overlay_external,user};KiVgaHooks vh={mc_snapshot,user};
    int handled=ki_overlay_body(m,t,&h)||ki_mapcells_body(m,t,&h)||ki_rect_body(m,t,&h)||
        ki_aligned_body(m,t,&h)||ki_vga_body(m,t,&vh)||ki_display_body(m,t,&h)||
        ki_hotspot_body(m,t,&h)||ki_numeric_body(m,t,&h);
    assert(handled);(void)handled;
}
static void overlay_run(KiMachine16 *m,uint16_t t,KiMapcellsFixture *f) {
    KiEconomyHooks h={mc_snapshot,overlay_external,f};ki_overlay_invoke(m,t,&h);
}
#endif
