#ifndef WOLONG_C_INPUT_FIXTURE_H
#define WOLONG_C_INPUT_FIXTURE_H
#include "vga_fixture.h"
void wolong_input_checkpoint(KiMachine16 *,uint16_t);
typedef struct {KiVgaFixture trace;unsigned unsupported;uint16_t core_cs,mouse_cs;} KiInputFixture;
static void mx_snapshot(KiMachine16 *m,uint16_t t,void *user) {
    KiInputFixture *f=user;if(m->cs==f->core_cs&&t==0x2259) wolong_input_checkpoint(m,t);
    vg_snapshot(m,t,&f->trace);
}
static void mx_external(KiMachine16 *m,uint16_t t,void *user) {
    KiInputFixture *f=user;KiEconomyHooks h={mx_snapshot,mx_external,user};KiVgaHooks vh={mx_snapshot,user};
    if(m->cs==f->mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h)) f->unsupported=t;return;}
    if(!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&!ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&
       !ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&!ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&
       !ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&!ki_modal_body(m,t,&h)) f->unsupported=t;
}
static void input_run(KiMachine16 *m,uint32_t loc,KiInputFixture *f) {
    KiEconomyHooks h={mx_snapshot,mx_external,f};ki_input_invoke(m,loc,&h);
}
#endif
