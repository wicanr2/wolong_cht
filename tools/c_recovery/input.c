/* Original two-segment input and mouse control; spec/231. Include after talk.c. */
#include "input.h"
#ifndef KI_INPUT_MUTATION
#define KI_INPUT_MUTATION 0
#endif
static void mx_call(KiMachine16 *m,uint16_t t,uint16_t next,uint32_t base,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_input_invoke(m,base+t,h);
}
static void mx_far(KiMachine16 *m,uint16_t site,uint16_t next,const KiEconomyHooks *h) {
    uint16_t off=dl_read16(m,m->cs,(uint16_t)(site+1)),seg=dl_read16(m,m->cs,(uint16_t)(site+3));
    ec_push(m,m->cs);ec_push(m,next);m->cs=seg;m->ip=off;ki_input_invoke(m,0x20000u+off,h);
}
#include "input_generated.inc"
void ki_input_invoke(KiMachine16 *m,uint32_t loc,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,(uint16_t)loc,h->user);
    uint16_t t=(uint16_t)loc;KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_input_body(m,loc,h)&&!ki_talk_body(m,t,h)&&!ki_numbers_body(m,t,h)&&!ki_glyph_body(m,t,h)&&
       !ki_display_body(m,t,h)&&!ki_rect_body(m,t,h)&&!ki_aligned_body(m,t,h)&&
       !ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&!ki_modal_body(m,t,h)&&h&&h->external)
        h->external(m,t,h->user);
    m->ip=ec_pop(m);
    if(loc==0x20000||loc==0x20101) m->cs=ec_pop(m);
}

#define MX_WRAP(n,k) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_input_body(m,k,h);m->ip=ec_pop(m);if(k==0x20000) m->cs=ec_pop(m);}
MX_WRAP(sub_18810,0x18810)
MX_WRAP(sub_121E7,0x121e7)
MX_WRAP(sub_20000,0x20000)
MX_WRAP(sub_2002E,0x2002e)
MX_WRAP(nullsub_4,0x2006f)
MX_WRAP(sub_20070,0x20070)
MX_WRAP(sub_2009A,0x2009a)
MX_WRAP(sub_200BD,0x200bd)
MX_WRAP(sub_200C0,0x200c0)
MX_WRAP(sub_20137,0x20137)
MX_WRAP(sub_2016B,0x2016b)
MX_WRAP(sub_2019D,0x2019d)
MX_WRAP(sub_201C6,0x201c6)
MX_WRAP(sub_201E4,0x201e4)
MX_WRAP(sub_2020C,0x2020c)
MX_WRAP(sub_20249,0x20249)
MX_WRAP(sub_202A0,0x202a0)
MX_WRAP(sub_202BD,0x202bd)
MX_WRAP(sub_202FE,0x202fe)
#undef MX_WRAP
