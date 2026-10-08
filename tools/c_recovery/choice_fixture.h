#ifndef WOLONG_C_CHOICE_FIXTURE_H
#define WOLONG_C_CHOICE_FIXTURE_H
#include "vga_fixture.h"
void wolong_input_checkpoint(KiMachine16 *,uint16_t);
typedef struct {KiVgaFixture trace;unsigned unsupported;uint16_t core_cs,mouse_cs;} KiChoiceFixture;
static void ch_snapshot(KiMachine16 *m,uint16_t t,void *user) {
    KiChoiceFixture *f=user;if(m->cs==f->core_cs&&t==0x2259) wolong_input_checkpoint(m,t);
    vg_snapshot(m,t,&f->trace);
}
static void ch_external(KiMachine16 *m,uint16_t t,void *user) {
    KiChoiceFixture *f=user;KiEconomyHooks h={ch_snapshot,ch_external,user};KiVgaHooks vh={ch_snapshot,user};
    if(m->cs==f->mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h)) f->unsupported=t;return;}
    if(!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&!ki_numbers_body(m,t,&h)&&
       !ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&!ki_aligned_body(m,t,&h)&&
       !ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&!ki_modal_body(m,t,&h)) f->unsupported=t;
}
static void choice_run(KiMachine16 *m,uint16_t t,KiChoiceFixture *f) {
    KiEconomyHooks h={ch_snapshot,ch_external,f};ki_choice_invoke(m,t,&h);
}
#endif
