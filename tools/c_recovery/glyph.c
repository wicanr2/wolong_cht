/* Real original glyph raster, spec/228. Include after display.c. */
#include "glyph.h"
#ifndef KI_GLYPH_MUTATION
#define KI_GLYPH_MUTATION 0
#endif
static void gl_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {ec_push(m,next);m->ip=t;ki_glyph_invoke(m,t,h);}
static void gl_interrupt(KiMachine16 *m,uint8_t number,uint16_t next) {m->ip=next;assert(wolong_glyph_interrupt(m,number));}
static void gl_far(KiMachine16 *m,uint16_t site,uint16_t next,const KiEconomyHooks *h) {
    uint16_t off=dl_read16(m,m->cs,(uint16_t)(site+1)),seg=dl_read16(m,m->cs,(uint16_t)(site+3));
    ec_push(m,m->cs);ec_push(m,next);m->cs=seg;m->ip=off;
    if(h&&h->enter) h->enter(m,off,h->user);
    /* Pinned DOS/V font service stub: INT F7/F8; RETF. No C-side guest CPU step. */
    assert(seg==0x80&&(off==0x410||off==0x414));gl_interrupt(m,off==0x410?0xf7:0xf8,(uint16_t)(off+2));
    m->ip=ec_pop(m);m->cs=ec_pop(m);
}
#include "glyph_generated.inc"
void ki_glyph_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_glyph_body(m,t,h)&&!ki_display_body(m,t,h)&&!ki_rect_body(m,t,h)&&!ki_aligned_body(m,t,h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&h&&h->external) h->external(m,t,h->user);
    m->ip=ec_pop(m);
}
#define GL_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_glyph_body(m,t,h);m->ip=ec_pop(m);}
GL_WRAP(sub_1F720,0xf720) GL_WRAP(sub_1F7A4,0xf7a4) GL_WRAP(sub_106F5,0x6f5) GL_WRAP(sub_106FD,0x6fd)
#undef GL_WRAP
