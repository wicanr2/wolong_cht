#ifndef WOLONG_C_RESOURCE_FIXTURE_H
#define WOLONG_C_RESOURCE_FIXTURE_H
#include "vga_fixture.h"
void wolong_input_checkpoint(KiMachine16 *,uint16_t);
typedef struct {KiVgaFixture trace;unsigned unsupported;uint16_t core_cs,mouse_cs;} KiResourceFixture;
static void rs_snapshot(KiMachine16 *m,uint16_t t,void *user) {
    KiResourceFixture *f=user;if(m->cs==f->core_cs&&t==0x2259) wolong_input_checkpoint(m,t);
    vg_snapshot(m,t,&f->trace);
}
static void rs_external(KiMachine16 *m,uint16_t t,void *user) {
    KiResourceFixture *f=user;KiEconomyHooks h={rs_snapshot,rs_external,user};KiVgaHooks vh={rs_snapshot,user};
    if(m->cs==f->mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h)) f->unsupported=t;return;}
    if(!ki_resource_body(m,t,&h)&&!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&
       !ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&
       !ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&!ki_modal_body(m,t,&h)) f->unsupported=t;
}
static void resource_run(KiMachine16 *m,uint16_t t,KiResourceFixture *f) {
    KiEconomyHooks h={rs_snapshot,rs_external,f};ki_resource_invoke(m,t,&h);
}
#endif
