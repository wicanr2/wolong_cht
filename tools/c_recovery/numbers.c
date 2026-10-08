/* Original number raster; spec/229. Include after glyph.c. */
#include "numbers.h"
#ifndef KI_NUMBERS_MUTATION
#define KI_NUMBERS_MUTATION 0
#endif
static void nr_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_numbers_invoke(m,t,h);
}
#include "numbers_generated.inc"
void ki_numbers_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);
    KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_numbers_body(m,t,h)&&!ki_glyph_body(m,t,h)&&!ki_display_body(m,t,h)&&
       !ki_rect_body(m,t,h)&&!ki_aligned_body(m,t,h)&&!ki_vga_body(m,t,&vh)&&
       !ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&h&&h->external)
        h->external(m,t,h->user);
    m->ip=ec_pop(m);
}
#define NR_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_numbers_body(m,t,h);m->ip=ec_pop(m);}
NR_WRAP(sub_1062F,0x62f) NR_WRAP(sub_1069A,0x69a) NR_WRAP(sub_106DE,0x6de)
NR_WRAP(sub_11E17,0x1e17)
#undef NR_WRAP
