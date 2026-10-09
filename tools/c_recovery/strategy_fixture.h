#ifndef WOLONG_C_STRATEGY_FIXTURE_H
#define WOLONG_C_STRATEGY_FIXTURE_H
#include "list_fixture.h"
typedef KiListFixture KiStrategyFixture;
void wolong_strategy_checkpoint(KiMachine16 *,uint16_t);
static void sty_snapshot(KiMachine16 *m,uint16_t t,void *user) {
 KiListFixture *f=user;if(m->cs==f->core_cs)wolong_strategy_checkpoint(m,t);ll_snapshot(m,t,user);
}
static void sty_external(KiMachine16 *m,uint16_t t,void *user) {
 KiListFixture *f=user;KiEconomyHooks h={sty_snapshot,sty_external,user};KiVgaHooks vh={sty_snapshot,user};
 if(m->cs==f->mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h))f->unsupported=t;return;}
 int escaped=0;
 if(ki_list_body(m,t,&h,&escaped)) {assert(!escaped);return;}
 if(!ki_strategy_body(m,t,&h)&&!ki_march_body(m,t,&h)&&!ki_details_body(m,t,&h)&&!ki_personnel_body(m,t,&h)&&!ki_formation_body(m,t,&h)&&!ki_catalog_body(m,t,&h)&&!ki_verdict_body(m,t,&h)&&!ki_resume_body(m,t,&h)&&!ki_overlay_body(m,t,&h)&&!ki_mapcells_body(m,t,&h)&&
  !ki_resource_body(m,t,&h)&&!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&
  !ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&
  !ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&
  !ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&
  !ki_politics_body(m,t,&h)&&!ki_world_body(m,t,&h))f->unsupported=t;
}
static void strategy_run(KiMachine16 *m,uint16_t t,KiStrategyFixture *f) {
 KiEconomyHooks h={sty_snapshot,sty_external,f};ki_strategy_invoke(m,t,&h);
}
#endif
