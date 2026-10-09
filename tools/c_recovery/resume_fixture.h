#ifndef WOLONG_C_RESUME_FIXTURE_H
#define WOLONG_C_RESUME_FIXTURE_H
#include "resource_fixture.h"
typedef KiResourceFixture KiResumeFixture;
static void resume_external(KiMachine16 *m,uint16_t t,void *user) {
 KiResumeFixture *f=user;KiEconomyHooks h={rs_snapshot,resume_external,user};KiVgaHooks vh={rs_snapshot,user};
 if(m->cs==f->mouse_cs) {if(!ki_input_body(m,0x20000u+t,&h)) f->unsupported=t;return;}
 if(!ki_resume_body(m,t,&h)&&!ki_overlay_body(m,t,&h)&&!ki_mapcells_body(m,t,&h)&&!ki_resource_body(m,t,&h)&&
  !ki_choice_body(m,t,&h)&&!ki_input_body(m,0x10000u+t,&h)&&!ki_talk_body(m,t,&h)&&!ki_numbers_body(m,t,&h)&&
  !ki_glyph_body(m,t,&h)&&!ki_display_body(m,t,&h)&&!ki_rect_body(m,t,&h)&&!ki_aligned_body(m,t,&h)&&
  !ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,&h)&&!ki_numeric_body(m,t,&h)&&!ki_modal_body(m,t,&h)) f->unsupported=t;
}
static void resume_run(KiMachine16 *m,uint16_t t,KiResumeFixture *f) {
 KiEconomyHooks h={rs_snapshot,resume_external,f};ki_resume_invoke(m,t,&h);
}
#endif
