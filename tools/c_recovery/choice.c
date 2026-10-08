/* Original popup selector and live row callback; spec/232. Include after input.c. */
#include "choice.h"
#ifndef KI_CHOICE_MUTATION
#define KI_CHOICE_MUTATION 0
#endif
static void ch_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_choice_invoke(m,t,h);
}
#include "choice_generated.inc"
void ki_choice_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_choice_body(m,t,h)&&!ki_input_body(m,0x10000u+t,h)&&!ki_talk_body(m,t,h)&&!ki_numbers_body(m,t,h)&&
       !ki_glyph_body(m,t,h)&&!ki_display_body(m,t,h)&&!ki_rect_body(m,t,h)&&!ki_aligned_body(m,t,h)&&
       !ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&!ki_modal_body(m,t,h)&&h&&h->external)
        h->external(m,t,h->user);
    m->ip=ec_pop(m);
}

#define CH_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_choice_body(m,t,h);m->ip=ec_pop(m);}
CH_WRAP(sub_1036F,0x36f)
CH_WRAP(sub_103C3,0x3c3)
CH_WRAP(sub_103E6,0x3e6)
CH_WRAP(sub_10414,0x414)
CH_WRAP(sub_104B5,0x4b5)
CH_WRAP(sub_104FF,0x4ff)
CH_WRAP(sub_1054D,0x54d)
CH_WRAP(sub_1059B,0x59b)
CH_WRAP(nullsub_5,0x5ed)
CH_WRAP(sub_1061F,0x61f)
CH_WRAP(sub_10B46,0xb46)
CH_WRAP(sub_10BAF,0xbaf)
CH_WRAP(sub_19479,0x9479)
CH_WRAP(sub_194BF,0x94bf)
#undef CH_WRAP
