/* Original resource loader, cue selection and segment allocation; spec/233. */
#include "resource.h"
#ifndef KI_RESOURCE_MUTATION
#define KI_RESOURCE_MUTATION 0
#endif
static void rs_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_resource_invoke(m,t,h);
}
#include "resource_generated.inc"
void ki_resource_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_resource_body(m,t,h)&&!ki_choice_body(m,t,h)&&!ki_input_body(m,0x10000u+t,h)&&!ki_talk_body(m,t,h)&&
       !ki_numbers_body(m,t,h)&&!ki_glyph_body(m,t,h)&&!ki_display_body(m,t,h)&&!ki_rect_body(m,t,h)&&
       !ki_aligned_body(m,t,h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&
       !ki_modal_body(m,t,h)&&h&&h->external) h->external(m,t,h->user);
    m->ip=ec_pop(m);
}

#define RS_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_resource_body(m,t,h);m->ip=ec_pop(m);}
RS_WRAP(sub_187AF,0x87af)
RS_WRAP(sub_1E378,0xe378)
RS_WRAP(sub_1F4A2,0xf4a2)
RS_WRAP(sub_10241,0x241)
RS_WRAP(sub_102C2,0x2c2)
RS_WRAP(sub_100DF,0xdf)
#undef RS_WRAP
