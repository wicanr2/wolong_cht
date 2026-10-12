#ifndef WOLONG_C_RESIDUAL_FIXTURE_H
#define WOLONG_C_RESIDUAL_FIXTURE_H
#include "startup_fixture.h"
#include <setjmp.h>
typedef KiMainFixture KiResidualFixture;
void wolong_residual_checkpoint(KiMachine16 *,uint16_t);
void wolong_residual_snapshot(KiMachine16 *,uint16_t);
static jmp_buf residual_exit_jump;
static int residual_exited;static uint8_t residual_exit_code;
static int residual_did_exit(void){return residual_exited;}
static uint8_t residual_exit_status(void){return residual_exit_code;}
_Noreturn void wolong_residual_exit(KiMachine16 *m,uint8_t code,void *user) {
 KiResidualFixture *f=user;(void)m;(void)f;residual_exited=1;residual_exit_code=code;longjmp(residual_exit_jump,1);
}
static void rsf_snapshot(KiMachine16 *m,uint16_t t,void *user) {
 KiResidualFixture *f=user;
 if(m->cs==f->list.core_cs)wolong_residual_checkpoint(m,t);
 wolong_residual_snapshot(m,t);
}
static void rsf_external(KiMachine16 *m,uint16_t t,void *user);
static void rsf_clock_call(KiMachine16 *m,uint16_t target,void *user) {
 KiResidualFixture *f=user;
 uint16_t next=dl_read16(m,m->ss,m->sp);
 rsf_external(m,target,f);
 ec_push(m,next);
}
static void rsf_external(KiMachine16 *m,uint16_t t,void *user) {
 KiResidualFixture *f=user;KiEconomyHooks h={rsf_snapshot,rsf_external,user};KiVgaHooks vh={rsf_snapshot,user};
 if(m->cs==f->list.mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h))f->list.unsupported=t;return;}
  if(t==0x563b) {
   const uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,m->ss,m->sp);
   sub_1563B(m);
   assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
   ec_push(m,ret);return;
  }
  if(t==0x1D8E) {
   const uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,m->ss,m->sp);
   const KiClockHooks ch={rsf_clock_call,0,f};
   sub_11D8E(m,&ch);
   assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
   ec_push(m,ret);return;
  }
  if(t==0xEC82) {
   const uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,m->ss,m->sp);
   uint32_t fr=((uint32_t)m->ss<<4)+(uint16_t)(sp-16);
   uint16_t save=(uint16_t)(m->memory[fr&0xfffffu]|(uint16_t)m->memory[(fr+1)&0xfffffu]<<8);
   sub_1EC82_abi(m,(uint8_t)(m->cx>>8),(uint8_t)m->cx,(uint8_t)(m->dx>>8),0);
   m->memory[fr&0xfffffu]=(uint8_t)save;m->memory[(fr+1)&0xfffffu]=(uint8_t)(save>>8);
   assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
   ec_push(m,ret);return;
  }
 int escaped=0;
 if(ki_list_body(m,t,&h,&escaped)) {assert(!escaped);return;}
 if(!ki_residual_body(m,t,&h)&&!ki_startup_body(m,t,&h)&&!ki_army_body(m,t,&h)&&!ki_engagement_body(m,t,&h)&&!ki_outcome_body(m,t,&h)&&!ki_route_body(m,t,&h)&&!ki_tick_body(m,t,&h)&&!ki_interaction_body(m,t,&h)&&!ki_bootstrap_body(m,t,&h)&&!ki_main_body(m,t,&h)&&!ki_strategy_body(m,t,&h)&&!ki_march_body(m,t,&h)&&
  !ki_details_body(m,t,&h)&&!ki_personnel_body(m,t,&h)&&!ki_formation_body(m,t,&h)&&!ki_catalog_body(m,t,&h)&&
  !ki_verdict_body(m,t,&h)&&!ki_resume_body(m,t,&h)&&!ki_overlay_body(m,t,&h)&&!ki_mapcells_body(m,t,&h)&&
  !ki_resource_body(m,t,&h)&&!ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&
  !ki_numbers_body(m,t,&h)&&!ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&
  !ki_aligned_body(m,t,&h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&
  !ki_modal_body(m,t,&h)&&!ki_events_body(m,t,&h)&&!ki_hourly_body(m,t,&h)&&!ki_politics_body(m,t,&h)&&
  !ki_world_body(m,t,&h))f->list.unsupported=t;
}
static void residual_run(KiMachine16 *m,uint16_t t,KiResidualFixture *f) {
 const KiEconomyHooks h={rsf_snapshot,rsf_external,f};
 residual_exited=0;residual_exit_code=0;
 f->armed=1;f->transferred=0;f->transfer_ip=f->transfer_ss=f->transfer_sp=0;
 if(setjmp(f->jump)==0) {
  if(setjmp(residual_exit_jump)==0)ki_residual_invoke(m,t,&h);
 }
 f->armed=0;
}
#endif
