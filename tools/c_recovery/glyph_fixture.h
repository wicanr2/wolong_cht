#ifndef WOLONG_C_GLYPH_FIXTURE_H
#define WOLONG_C_GLYPH_FIXTURE_H
#include "vga_fixture.h"
typedef struct {KiVgaFixture trace;unsigned unsupported;} KiGlyphFixture;
static void gg_snapshot(KiMachine16 *m,uint16_t t,void *user) {KiGlyphFixture *f=user;vg_snapshot(m,t,&f->trace);}
static void gg_external(KiMachine16 *m,uint16_t t,void *user) {
    KiGlyphFixture *f=user;KiEconomyHooks h={gg_snapshot,gg_external,user};KiVgaHooks vh={gg_snapshot,user};
    if(!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&!ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)) f->unsupported=t;
}
static void glyph_run(KiMachine16 *m,uint16_t t,KiGlyphFixture *f) {KiEconomyHooks h={gg_snapshot,gg_external,f};ki_glyph_invoke(m,t,&h);}
#endif
