/* Original TALK, portrait cache and DOS read ABI; spec/230. Include after numbers.c. */
#include "talk.h"
#ifndef KI_TALK_MUTATION
#define KI_TALK_MUTATION 0
#endif
static uint16_t tx_rcl(KiMachine16 *m,uint16_t v,unsigned width) {
    uint16_t mask=width==8?255:65535,sign=width==8?128:32768;
    uint16_t carry=(v&sign)!=0,result=(uint16_t)(((v<<1)|(m->flags&1))&mask);
    m->flags=(uint16_t)((m->flags&~0x801u)|carry|((((result&sign)!=0)^carry)<<11));
    return result;
}
static void tx_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_talk_invoke(m,t,h);
}
#include "talk_generated.inc"
void ki_talk_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);
    KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_talk_body(m,t,h)&&!ki_numbers_body(m,t,h)&&!ki_glyph_body(m,t,h)&&
       !ki_display_body(m,t,h)&&!ki_rect_body(m,t,h)&&!ki_aligned_body(m,t,h)&&
       !ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&h&&h->external)
        h->external(m,t,h->user);
    m->ip=ec_pop(m);
}
#define TX_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_talk_body(m,t,h);m->ip=ec_pop(m);}
TX_WRAP(sub_106F9,0x6f9) TX_WRAP(sub_1075B,0x75b) TX_WRAP(sub_107D2,0x7d2) TX_WRAP(sub_1084A,0x84a)
TX_WRAP(sub_18853,0x8853) TX_WRAP(sub_189A4,0x89a4) TX_WRAP(sub_1E38C,0xe38c) TX_WRAP(sub_1F4DF,0xf4df)
#undef TX_WRAP
