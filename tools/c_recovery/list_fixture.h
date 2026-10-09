#ifndef WOLONG_C_LIST_FIXTURE_H
#define WOLONG_C_LIST_FIXTURE_H
#include "vga_fixture.h"
void wolong_list_checkpoint(KiMachine16 *,uint16_t);
void wolong_input_checkpoint(KiMachine16 *,uint16_t);
typedef struct {unsigned calls;KiVgaFixture blocks[4];unsigned unsupported;uint16_t core_cs,mouse_cs;} KiListFixture;
static void ll_snapshot(KiMachine16 *m,uint16_t t,void *user) {
 KiListFixture *f=user;if(m->cs==f->core_cs&&t==0x2259) wolong_input_checkpoint(m,t);
 if(m->cs==f->core_cs&&(t==0x21e7||t==0x84dd)) wolong_list_checkpoint(m,t);
 unsigned block=f->calls++/8192;assert(block<4);vg_snapshot(m,t,&f->blocks[block]);
}
static void ll_external(KiMachine16 *m,uint16_t t,void *user) {
 KiListFixture *f=user;KiEconomyHooks h={ll_snapshot,ll_external,user};KiVgaHooks vh={ll_snapshot,user};
 if(m->cs==f->mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h)) f->unsupported=t;return;}
 int escaped=0;
 if(ki_list_body(m,t,&h,&escaped)) {assert(!escaped);return;}
 if(!ki_verdict_body(m,t,&h)&&!ki_resume_body(m,t,&h)&&!ki_overlay_body(m,t,&h)&&!ki_mapcells_body(m,t,&h)&&
  !ki_resource_body(m,t,&h)&&!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&
  !ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&
  !ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&
  !ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)) f->unsupported=t;
}
static void list_run(KiMachine16 *m,uint16_t t,KiListFixture *f) {
 KiEconomyHooks h={ll_snapshot,ll_external,f};ki_list_invoke(m,t,&h);
}
#endif
