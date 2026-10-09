/* Original scene/camera/minimap closure; include after resource.c and overlay.c. */
#include "resume.h"
#ifndef KI_RESUME_MUTATION
#define KI_RESUME_MUTATION 0
#endif
static void sr_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_resume_invoke(m,t,h);
}
static void sr_tail(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 m->ip=t;if(h&&h->enter) h->enter(m,t,h->user);
 int handled=ki_resume_body(m,t,h);assert(handled);(void)handled;
}
#include "resume_generated.inc"
void ki_resume_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(h&&h->enter) h->enter(m,t,h->user);KiVgaHooks vh={h?h->enter:0,h?h->user:0};
 int handled=ki_resume_body(m,t,h)||ki_overlay_body(m,t,h)||ki_mapcells_body(m,t,h)||
  ki_resource_body(m,t,h)||ki_choice_body(m,t,h)||ki_input_body(m,0x10000u+t,h)||ki_talk_body(m,t,h)||
  ki_numbers_body(m,t,h)||ki_glyph_body(m,t,h)||ki_display_body(m,t,h)||ki_rect_body(m,t,h)||
  ki_aligned_body(m,t,h)||ki_vga_body(m,t,&vh)||ki_hotspot_body(m,t,h)||ki_numeric_body(m,t,h)||ki_modal_body(m,t,h);
 if(!handled&&h&&h->external) h->external(m,t,h->user);else {assert(handled);}
 m->ip=ec_pop(m);
}
#define SR_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_resume_body(m,t,h);m->ip=ec_pop(m);}
SR_WRAP(sub_11F30,0x1f30) SR_WRAP(sub_11F5A,0x1f5a) SR_WRAP(sub_13B08,0x3b08)
SR_WRAP(sub_15C58,0x5c58) SR_WRAP(sub_19541,0x9541) SR_WRAP(sub_195B1,0x95b1)
SR_WRAP(sub_196ED,0x96ed) SR_WRAP(sub_19752,0x9752) SR_WRAP(sub_1D4A3,0xd4a3)
#undef SR_WRAP
