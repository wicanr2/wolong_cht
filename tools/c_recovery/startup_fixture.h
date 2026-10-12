#ifndef WOLONG_C_STARTUP_FIXTURE_H
#define WOLONG_C_STARTUP_FIXTURE_H
#include "army_fixture.h"
#include <setjmp.h>
typedef KiMainFixture KiStartupFixture;
void wolong_startup_checkpoint(KiMachine16 *,uint16_t);
void wolong_startup_snapshot(KiMachine16 *,uint16_t);
static jmp_buf startup_pause_jump;
static int startup_pause_armed,startup_paused;
static unsigned startup_pause_target,startup_pause_hits;
static int startup_is_paused(void){return startup_paused;}
static jmp_buf startup_exit_jump;
static int startup_exited;static uint8_t startup_exit_code;
static int startup_did_exit(void){return startup_exited;}
static uint8_t startup_exit_status(void){return startup_exit_code;}
_Noreturn void wolong_startup_exit(KiMachine16 *m,uint8_t code,void *user) {
 KiStartupFixture *f=user;(void)m;(void)f;startup_exited=1;startup_exit_code=code;longjmp(startup_exit_jump,1);
}
static void suf_snapshot(KiMachine16 *m,uint16_t t,void *user) {
 KiStartupFixture *f=user;
 if(startup_pause_armed&&m->cs==f->list.core_cs&&t==0x1C1F){startup_pause_hits++;if(startup_pause_hits>=startup_pause_target){startup_paused=1;longjmp(startup_pause_jump,1);}}
 if(m->cs==f->list.core_cs)wolong_startup_checkpoint(m,t);
 wolong_startup_snapshot(m,t);
}
static void suf_external(KiMachine16 *m,uint16_t t,void *user);
static void suf_clock_call(KiMachine16 *m,uint16_t target,void *user) {
 KiStartupFixture *f=user;
 uint16_t next=dl_read16(m,m->ss,m->sp);
 suf_external(m,target,f);
 ec_push(m,next);
}
static void suf_external(KiMachine16 *m,uint16_t t,void *user) {
 KiStartupFixture *f=user;KiEconomyHooks h={suf_snapshot,suf_external,user};KiVgaHooks vh={suf_snapshot,user};
 if(m->cs==f->list.mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h))f->list.unsupported=t;return;}
  if(t==0x563b) {
   const uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,m->ss,m->sp);
   sub_1563B(m);
   assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
   ec_push(m,ret);return;
  }
  if(t==0x1D8E) {
   /* Game clock (docs/spec/203): hand-written body. The cascade callees
      never fire in-matrix (subtick<8 at entry; proven by locwatch), so the
      call hook is a safety net routing through this dispatcher. Waits exit
      immediately on the pre-seeded D2C/D2D (ticks never fire in-matrix). */
   const uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,m->ss,m->sp);
   const KiClockHooks ch={suf_clock_call,0,f};
   sub_11D8E(m,&ch);
   assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
   ec_push(m,ret);return;
  }
  if(t==0xEC82) {
   /* RNG seed (docs/spec/202): hand-written ABI, reached here via 1006B.
      Golem INT 1Ah AH=02 answers CF=1 with regs untouched, and 1EC82 ignores
      CF, so the shuffle observes the stale entry CX/DX; AL is always 0
      (the 256-fill loop wraps it) by the time of the RTC read. */
   const uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,m->ss,m->sp);
   /* The seed ABI models a pushed INT frame at SP-16; this harness routes
      INT 1Ah through the hooked handler (no frame), so save and restore
      that word around the call. Same read/call/restore shape as 0x563b. */
   uint32_t fr=((uint32_t)m->ss<<4)+(uint16_t)(sp-16);
   uint16_t save=(uint16_t)(m->memory[fr&0xfffffu]|(uint16_t)m->memory[(fr+1)&0xfffffu]<<8);
   sub_1EC82_abi(m,(uint8_t)(m->cx>>8),(uint8_t)m->cx,(uint8_t)(m->dx>>8),0);
   m->memory[fr&0xfffffu]=(uint8_t)save;m->memory[(fr+1)&0xfffffu]=(uint8_t)(save>>8);
   assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
   ec_push(m,ret);return;
  }
 int escaped=0;
 if(ki_list_body(m,t,&h,&escaped)) {assert(!escaped);return;}
 if(!ki_startup_body(m,t,&h)&&!ki_army_body(m,t,&h)&&!ki_engagement_body(m,t,&h)&&!ki_outcome_body(m,t,&h)&&!ki_route_body(m,t,&h)&&!ki_tick_body(m,t,&h)&&!ki_interaction_body(m,t,&h)&&!ki_bootstrap_body(m,t,&h)&&!ki_main_body(m,t,&h)&&!ki_strategy_body(m,t,&h)&&!ki_march_body(m,t,&h)&&
  !ki_details_body(m,t,&h)&&!ki_personnel_body(m,t,&h)&&!ki_formation_body(m,t,&h)&&!ki_catalog_body(m,t,&h)&&
  !ki_verdict_body(m,t,&h)&&!ki_resume_body(m,t,&h)&&!ki_overlay_body(m,t,&h)&&!ki_mapcells_body(m,t,&h)&&
  !ki_resource_body(m,t,&h)&&!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&
  !ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&
  !ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&
  !ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&!ki_politics_body(m,t,&h)&&
  !ki_world_body(m,t,&h))f->list.unsupported=t;
}
static void startup_run(KiMachine16 *m,uint16_t t,KiStartupFixture *f) {
 const KiEconomyHooks h={suf_snapshot,suf_external,f};
 startup_pause_armed=0;startup_paused=0;startup_exited=0;startup_exit_code=0;
 f->armed=1;f->transferred=0;f->transfer_ip=f->transfer_ss=f->transfer_sp=0;
 if(setjmp(f->jump)==0) {
  if(setjmp(startup_exit_jump)==0)ki_startup_invoke(m,t,&h);
 }
 f->armed=0;
}
static void startup_continue(KiMachine16 *m,uint16_t t,KiStartupFixture *f) {
 const KiEconomyHooks h={suf_snapshot,suf_external,f};
 if(setjmp(f->jump)==0) {
  if(setjmp(startup_exit_jump)==0)ki_startup_invoke(m,t,&h);
 }
}
static void startup_run_pause(KiMachine16 *m,uint16_t t,KiStartupFixture *f,unsigned iters) {
 const KiEconomyHooks h={suf_snapshot,suf_external,f};
 startup_pause_armed=1;startup_paused=0;startup_pause_target=iters;startup_pause_hits=0;
 startup_exited=0;f->armed=1;f->transferred=0;
 f->transfer_ip=f->transfer_ss=f->transfer_sp=0;
 if(setjmp(f->jump)==0) {
  if(setjmp(startup_exit_jump)==0) {
   if(setjmp(startup_pause_jump)==0)ki_startup_invoke(m,t,&h);
  }
 }
 startup_pause_armed=0;f->armed=0;
}
static void startup_warmup(KiMachine16 *m,KiStartupFixture *f,unsigned iters) {
 const KiEconomyHooks h={suf_snapshot,suf_external,f};
 startup_pause_armed=1;startup_paused=0;startup_pause_target=iters;startup_pause_hits=0;
 startup_exited=0;f->armed=1;f->transferred=0;
 f->transfer_ip=f->transfer_ss=f->transfer_sp=0;
 if(setjmp(f->jump)==0) {
  if(setjmp(startup_exit_jump)==0) {
   if(setjmp(startup_pause_jump)==0)ki_startup_invoke(m,0x1BE0,&h);
  }
 }
 startup_pause_armed=0;f->armed=0;
}
#endif
