#ifndef WOLONG_C_ROUTE_FIXTURE_H
#define WOLONG_C_ROUTE_FIXTURE_H
#include "main_fixture.h"
typedef KiMainFixture KiRouteFixture;
void wolong_route_checkpoint(KiMachine16 *,uint16_t);
static void rt_snapshot(KiMachine16 *m,uint16_t t,void *user) {
 KiRouteFixture *f=user;if(m->cs==f->list.core_cs)wolong_route_checkpoint(m,t);ll_snapshot(m,t,&f->list);
}
static void rt_external(KiMachine16 *m,uint16_t t,void *user) {
 KiRouteFixture *f=user;KiEconomyHooks h={rt_snapshot,rt_external,user};KiVgaHooks vh={rt_snapshot,user};
 if(m->cs==f->list.mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h))f->list.unsupported=t;return;}
 int escaped=0;
 if(ki_list_body(m,t,&h,&escaped)) {assert(!escaped);return;}
 if(!ki_route_body(m,t,&h)&&!ki_tick_body(m,t,&h)&&!ki_interaction_body(m,t,&h)&&!ki_bootstrap_body(m,t,&h)&&!ki_main_body(m,t,&h)&&!ki_strategy_body(m,t,&h)&&!ki_march_body(m,t,&h)&&
  !ki_details_body(m,t,&h)&&!ki_personnel_body(m,t,&h)&&!ki_formation_body(m,t,&h)&&!ki_catalog_body(m,t,&h)&&
  !ki_verdict_body(m,t,&h)&&!ki_resume_body(m,t,&h)&&!ki_overlay_body(m,t,&h)&&!ki_mapcells_body(m,t,&h)&&
  !ki_resource_body(m,t,&h)&&!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&
  !ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&
  !ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&
  !ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&!ki_politics_body(m,t,&h)&&
  !ki_world_body(m,t,&h))f->list.unsupported=t;
}
static void route_run(KiMachine16 *m,uint16_t t,KiRouteFixture *f) {
 const KiEconomyHooks h={rt_snapshot,rt_external,f};
 f->armed=1;f->transferred=0;f->transfer_ip=f->transfer_ss=f->transfer_sp=0;
 if(setjmp(f->jump)==0)ki_route_invoke(m,t,&h);
 f->armed=0;
}
#endif
