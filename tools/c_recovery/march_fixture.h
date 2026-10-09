#ifndef WOLONG_C_MARCH_FIXTURE_H
#define WOLONG_C_MARCH_FIXTURE_H
#include "list_fixture.h"
typedef KiListFixture KiMarchFixture;
void wolong_march_checkpoint(KiMachine16 *,uint16_t);
static void mch_snapshot(KiMachine16 *m,uint16_t t,void *user) {
 KiListFixture *f=user;if(m->cs==f->core_cs)wolong_march_checkpoint(m,t);ll_snapshot(m,t,user);
}
static void mch_external(KiMachine16 *m,uint16_t t,void *user) {
 KiListFixture *f=user;KiEconomyHooks h={mch_snapshot,mch_external,user};KiVgaHooks vh={mch_snapshot,user};
 if(m->cs==f->mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h)) f->unsupported=t;return;}
 int escaped=0;
 if(ki_list_body(m,t,&h,&escaped)) {assert(!escaped);return;}
 if(!ki_march_body(m,t,&h)&&!ki_details_body(m,t,&h)&&!ki_personnel_body(m,t,&h)&&!ki_formation_body(m,t,&h)&&!ki_catalog_body(m,t,&h)&&!ki_verdict_body(m,t,&h)&&!ki_resume_body(m,t,&h)&&!ki_overlay_body(m,t,&h)&&!ki_mapcells_body(m,t,&h)&&
  !ki_resource_body(m,t,&h)&&!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&
  !ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&
  !ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&
  !ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)) f->unsupported=t;
}
static void march_run(KiMachine16 *m,uint16_t t,KiMarchFixture *f) {
 KiEconomyHooks h={mch_snapshot,mch_external,f};ki_march_invoke(m,t,&h);
}
#endif
