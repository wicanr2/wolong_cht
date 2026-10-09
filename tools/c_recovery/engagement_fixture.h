#ifndef WOLONG_C_ENGAGEMENT_FIXTURE_H
#define WOLONG_C_ENGAGEMENT_FIXTURE_H
#include "main_fixture.h"
typedef KiMainFixture KiEngagementFixture;
void wolong_engagement_checkpoint(KiMachine16 *,uint16_t);
void wolong_engagement_snapshot(KiMachine16 *,uint16_t);
static jmp_buf engagement_pause_jump;
static int engagement_pause_armed,engagement_paused;
static int engagement_is_paused(void){return engagement_paused;}
static void egf_snapshot(KiMachine16 *m,uint16_t t,void *user) {
 KiEngagementFixture *f=user;if(engagement_pause_armed&&m->cs==f->list.core_cs&&t==0xa156){engagement_paused=1;longjmp(engagement_pause_jump,1);}if(m->cs==f->list.core_cs)wolong_engagement_checkpoint(m,t);if(m->cs==f->list.core_cs&&t==0x2259)wolong_input_checkpoint(m,t);if(m->cs==f->list.core_cs&&(t==0x21e7||t==0x84dd))wolong_list_checkpoint(m,t);wolong_engagement_snapshot(m,t);
}
static void egf_external(KiMachine16 *m,uint16_t t,void *user) {
 KiEngagementFixture *f=user;KiEconomyHooks h={egf_snapshot,egf_external,user};KiVgaHooks vh={egf_snapshot,user};
 if(m->cs==f->list.mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h))f->list.unsupported=t;return;}
 int escaped=0;
 if(ki_list_body(m,t,&h,&escaped)) {assert(!escaped);return;}
 if(!ki_engagement_body(m,t,&h)&&!ki_outcome_body(m,t,&h)&&!ki_route_body(m,t,&h)&&!ki_tick_body(m,t,&h)&&!ki_interaction_body(m,t,&h)&&!ki_bootstrap_body(m,t,&h)&&!ki_main_body(m,t,&h)&&!ki_strategy_body(m,t,&h)&&!ki_march_body(m,t,&h)&&
  !ki_details_body(m,t,&h)&&!ki_personnel_body(m,t,&h)&&!ki_formation_body(m,t,&h)&&!ki_catalog_body(m,t,&h)&&
  !ki_verdict_body(m,t,&h)&&!ki_resume_body(m,t,&h)&&!ki_overlay_body(m,t,&h)&&!ki_mapcells_body(m,t,&h)&&
  !ki_resource_body(m,t,&h)&&!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&
  !ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&
  !ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&
  !ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&!ki_politics_body(m,t,&h)&&
  !ki_world_body(m,t,&h))f->list.unsupported=t;
}
static void engagement_run(KiMachine16 *m,uint16_t t,KiEngagementFixture *f) {
 engagement_pause_armed=0;engagement_paused=0;const KiEconomyHooks h={egf_snapshot,egf_external,f};
 f->armed=1;f->transferred=0;f->transfer_ip=f->transfer_ss=f->transfer_sp=0;
 if(setjmp(f->jump)==0)ki_engagement_invoke(m,t,&h);
 f->armed=0;
}
static void engagement_warmup(KiMachine16 *m,KiEngagementFixture *f) {
 const KiEconomyHooks h={egf_snapshot,egf_external,f};
 engagement_pause_armed=1;engagement_paused=0;f->armed=1;f->transferred=0;
 f->transfer_ip=f->transfer_ss=f->transfer_sp=0;
 if(setjmp(f->jump)==0) {
  if(setjmp(engagement_pause_jump)==0)ki_engagement_invoke(m,0x1b5a,&h);
 }
 engagement_pause_armed=0;f->armed=0;assert(engagement_paused&&!f->transferred);
}
#endif
