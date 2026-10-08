#ifndef WOLONG_C_DISPLAY_FIXTURE_H
#define WOLONG_C_DISPLAY_FIXTURE_H
#include "vga_fixture.h"
typedef struct {KiVgaFixture trace;unsigned unsupported;} KiDisplayFixture;
static void dd_snapshot(KiMachine16 *m,uint16_t t,void *user) {KiDisplayFixture *f=user;vg_snapshot(m,t,&f->trace);}
static void dd_external(KiMachine16 *m,uint16_t t,void *user) {
    KiDisplayFixture *f=user;KiEconomyHooks h={dd_snapshot,dd_external,user};KiVgaHooks vh={dd_snapshot,user};
    /* Explicit raster boundary. Text parser/width/coordinates remain real C. */
    if(t==0xf75e) {ec_cmp(m,m->cx,m->cx,16);return;}
    if(!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&!ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)) f->unsupported=t;
}
static void display_run(KiMachine16 *m,uint16_t t,KiDisplayFixture *f) {KiEconomyHooks h={dd_snapshot,dd_external,f};ki_display_invoke(m,t,&h);}
#endif
