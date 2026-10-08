#ifndef WOLONG_C_TALK_FIXTURE_H
#define WOLONG_C_TALK_FIXTURE_H
#include "vga_fixture.h"
typedef struct {KiVgaFixture trace;unsigned unsupported;} KiTalkFixture;
static void tx_snapshot(KiMachine16 *m,uint16_t t,void *user) {
    KiTalkFixture *f=user;vg_snapshot(m,t,&f->trace);
}
static void tx_external(KiMachine16 *m,uint16_t t,void *user) {
    KiTalkFixture *f=user;
    KiEconomyHooks h={tx_snapshot,tx_external,user};KiVgaHooks vh={tx_snapshot,user};
    if(!ki_talk_body(m,t,&h)&&!ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&
       !ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&!ki_aligned_body(m,t,&h)&&
       !ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)) f->unsupported=t;
}
static void talk_run(KiMachine16 *m,uint16_t t,KiTalkFixture *f) {
    KiEconomyHooks h={tx_snapshot,tx_external,f};ki_talk_invoke(m,t,&h);
}
#endif
