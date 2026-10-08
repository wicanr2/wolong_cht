/* Sixteen DOS/V rectangle/bar/selection entries; spec/225.
 * Include after aligned.c. Original clipping, flags, masks, byte width and ABI.
 */
#include "rect.h"
#ifndef KI_RECT_MUTATION
#define KI_RECT_MUTATION 0
#endif
static void rc_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_rect_invoke(m,t,h);
}
static void rc_test(KiMachine16 *m,uint16_t value,uint16_t mask,unsigned width) {ec_logic(m,value&mask,width);}
static void rc_flag(KiMachine16 *m,uint8_t bit) {uint8_t b=ec_byte(m,m->cs,0xf1a2)|bit;ec_put(m,m->cs,0xf1a2,b);ec_logic(m,b,8);}
static void rc_write(KiMachine16 *m,uint16_t seg,uint16_t off,uint8_t b) {wolong_vga_write(ec_at(seg,off),b);}
static uint8_t rc_read(KiMachine16 *m,uint16_t seg,uint16_t off) {(void)m;return wolong_vga_read(ec_at(seg,off));}
static uint8_t rc_sar8(KiMachine16 *m,uint8_t value,unsigned count) {
    while(count--) {uint8_t result=(value>>1)|(value&0x80);m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|(value&1)),result,8);value=result;}return value;
}
static void rc_out16(KiMachine16 *m) {wolong_vga_out(m->dx,(uint8_t)m->ax);wolong_vga_out((uint16_t)(m->dx+1),(uint8_t)(m->ax>>8));}
static void rc_mode(KiMachine16 *m,int reset) {
    ec_push(m,m->ax);ec_push(m,m->dx);m->dx=0x3ce;m->ax=reset?5:0x305;rc_out16(m);
    if(reset) {m->ax=0;ec_logic(m,0,16);rc_out16(m);m->ax=1;rc_out16(m);}
    m->ax=3;rc_out16(m);m->ax=0xff08;rc_out16(m);
    if(!reset) {vg_al(m,0);ec_logic(m,0,8);vg_out(m);}m->dx=ec_pop(m);m->ax=ec_pop(m);
}
static void rc_setup(KiMachine16 *m) {
    m->dx=0x3ce;vg_al(m,5);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,KI_RECT_MUTATION==6?0:3);vg_out(m);m->dx=st_dec(m,m->dx,16);
    vg_al(m,0);ec_logic(m,0,8);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,(uint8_t)(m->ax>>8));vg_out(m);m->dx=st_dec(m,m->dx,16);
    vg_al(m,3);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,0);ec_logic(m,0,8);vg_out(m);m->dx=st_dec(m,m->dx,16);vg_al(m,8);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,255);vg_out(m);
}
static void rc_vertical(KiMachine16 *m) {
    ec_push(m,m->dx);rc_test(m,m->ax>>8,1,8);
    if(m->flags&0x40) {
        m->cx=m->bp;m->cx=st_dec(m,m->cx,16);ec_push(m,m->di);
        do {if(KI_RECT_MUTATION!=5) vg_dl(m,rc_read(m,m->ds,m->di));rc_write(m,m->ds,m->di,(uint8_t)m->bx);m->di=ec_math(m,m->di,80,16,0,0);--m->cx;}while(m->cx);
        m->di=ec_pop(m);ec_cmp(m,m->si,8,16);if(m->flags&(1|0x40)) {m->dx=ec_pop(m);return;}
    }
    rc_test(m,m->ax>>8,2,8);if(!(m->flags&0x40)) {m->dx=ec_pop(m);return;}
    ec_push(m,m->di);vg_cl(m,(uint8_t)m->ax);vg_ch(m,0);ec_logic(m,0,8);rc_test(m,m->si,7,16);if(m->flags&0x40) m->cx=st_dec(m,m->cx,16);m->di=ec_math(m,m->di,m->cx,16,0,0);
    m->cx=m->bp;m->cx=st_dec(m,m->cx,16);do {if(KI_RECT_MUTATION!=5) vg_dl(m,rc_read(m,m->ds,m->di));rc_write(m,m->ds,m->di,(uint8_t)(m->bx>>8));m->di=ec_math(m,m->di,80,16,0,0);--m->cx;}while(m->cx);m->di=ec_pop(m);m->dx=ec_pop(m);
}
static void rc_horizontal(KiMachine16 *m) {
    ec_push(m,m->ax);ec_push(m,m->di);vg_cl(m,(uint8_t)m->ax);vg_ch(m,0);ec_logic(m,0,8);
    if(KI_RECT_MUTATION!=5) vg_al(m,rc_read(m,m->ds,m->di));rc_write(m,m->ds,m->di,(uint8_t)m->dx);m->di=ec_inc(m,m->di,16);m->cx=st_dec(m,m->cx,16);ec_cmp(m,m->si,8,16);
    if(!(m->flags&(1|0x40))) {
        ec_cmp(m,m->si,16,16);
        if(!(m->flags&(1|0x40))) {vg_al(m,255);while(m->cx) {rc_write(m,m->es,m->di,(uint8_t)m->ax);m->di=(uint16_t)(m->di+((m->flags&0x400)?-1:1));--m->cx;}rc_test(m,m->si,7,16);if(m->flags&0x40) goto done;}
        if(KI_RECT_MUTATION!=5) vg_al(m,rc_read(m,m->ds,m->di));rc_write(m,m->ds,m->di,(uint8_t)(m->dx>>8));
    }
done:m->di=ec_pop(m);m->ax=ec_pop(m);
}
static void rc_outline(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);ec_push(m,m->es);ec_push(m,m->cx);ec_push(m,m->bp);if(KI_RECT_MUTATION!=11) m->flags&=(uint16_t)~0x400u;ec_put(m,m->cs,0xf1a2,0);
    ec_cmp(m,m->si,m->dx,16);if(ec_less(m)) vg_swap(&m->si,&m->dx);
    ec_cmp(m,m->dx,640,16);if(!ec_less(m)) goto finish;ec_cmp(m,m->si,0,16);if(ec_less(m)) goto finish;
    ec_cmp(m,m->di,m->bx,16);if(ec_less(m)) vg_swap(&m->di,&m->bx);ec_cmp(m,m->bx,400,16);if(!ec_less(m)) goto finish;ec_cmp(m,m->di,0,16);if(ec_less(m)) goto finish;
    ec_cmp(m,m->dx,0,16);if(ec_less(m)) {m->dx=0;ec_logic(m,0,16);rc_flag(m,1);}ec_cmp(m,m->si,640,16);if(!(m->flags&1)) {m->si=639;rc_flag(m,2);}
    m->si=ec_math(m,m->si,m->dx,16,0,1);if(KI_RECT_MUTATION!=1) m->si=ec_inc(m,m->si,16);
    ec_cmp(m,m->bx,0,16);if(ec_less(m)) {m->bx=0;ec_logic(m,0,16);rc_flag(m,4);}ec_cmp(m,m->di,400,16);if(!(m->flags&1)) {m->di=399;rc_flag(m,8);}
    m->di=ec_math(m,m->di,m->bx,16,0,1);m->di=ec_inc(m,m->di,16);m->bp=m->di;m->bx=vg_left(m,m->bx,4);m->di=m->bx;m->bx=vg_left(m,m->bx,2);m->di=ec_math(m,m->di,m->bx,16,0,0);
    vg_ch(m,0);ec_logic(m,0,8);vg_cl(m,(uint8_t)m->dx);m->dx=vg_right(m,m->dx,3);m->di=ec_math(m,m->di,m->dx,16,0,0);vg_cl(m,(uint8_t)m->cx&7);ec_logic(m,(uint8_t)m->cx,8);m->si=ec_math(m,m->si,m->cx,16,0,0);
    m->dx=0xffff;m->bx=0x8080;vg_dl(m,(uint8_t)al_shr(m,(uint8_t)m->dx,(uint8_t)m->cx,8));gw_bl(m,(uint8_t)al_shr(m,(uint8_t)m->bx,(uint8_t)m->cx,8));m->cx=m->si;vg_cl(m,(uint8_t)m->cx&7);ec_logic(m,(uint8_t)m->cx,8);
    if(!(m->flags&0x40)) {vg_dh(m,(uint8_t)al_shr(m,m->dx>>8,(uint8_t)m->cx,8));vg_dh(m,(uint8_t)~(m->dx>>8));}
    vg_cl(m,(uint8_t)st_dec(m,(uint8_t)m->cx,8));vg_cl(m,(uint8_t)m->cx&7);ec_logic(m,(uint8_t)m->cx,8);gw_bh(m,(uint8_t)al_shr(m,m->bx>>8,(uint8_t)m->cx,8));ec_cmp(m,m->si,8,16);
    if(m->flags&(1|0x40)) {vg_dl(m,(uint8_t)m->dx&(uint8_t)(m->dx>>8));ec_logic(m,(uint8_t)m->dx,8);rc_test(m,ec_byte(m,m->cs,0xf1a2),2,8);if(m->flags&0x40) {gw_bl(m,(uint8_t)m->bx|(uint8_t)(m->bx>>8));ec_logic(m,(uint8_t)m->bx,8);}}
    m->cx=m->dx;rc_setup(m);m->dx=m->cx;m->ax=0xa0c8;m->ds=m->ax;m->es=m->ax;m->ax=m->si;m->ax=vg_right(m,m->ax,3);vg_ah(m,ec_byte(m,m->cs,0xf1a2));
    if(KI_RECT_MUTATION==2) vg_ah(m,0);
    rc_test(m,m->ax>>8,4,8);if(m->flags&0x40) {rc_call(m,0xf17b,0xf114,h);m->di=ec_math(m,m->di,80,16,0,0);m->bp=st_dec(m,m->bp,16);if(m->flags&0x40) goto finish;}
    ec_cmp(m,m->bp,1,16);if(!(m->flags&0x40)) {rc_test(m,m->ax>>8,8,8);if(!(m->flags&0x40)) m->bp=ec_inc(m,m->bp,16);rc_call(m,0xf140,0xf128,h);m->bp=st_dec(m,m->bp,16);m->bp=vg_left(m,m->bp,4);m->di=ec_math(m,m->di,m->bp,16,0,0);m->bp=vg_left(m,m->bp,2);m->di=ec_math(m,m->di,m->bp,16,0,0);}
    rc_test(m,m->ax>>8,8,8);if(m->flags&0x40) rc_call(m,0xf17b,0xf13b,h);
finish:m->bp=ec_pop(m);m->cx=ec_pop(m);m->es=ec_pop(m);m->ds=ec_pop(m);
}
static void rc_fill(KiMachine16 *m) {
    ec_cmp(m,m->dx,m->si,16);if(ec_greater(m)) vg_swap(&m->dx,&m->si);ec_cmp(m,m->bx,m->di,16);if(ec_greater(m)) vg_swap(&m->bx,&m->di);
    ec_cmp(m,m->dx,0,16);if(ec_less(m)) {ec_cmp(m,m->si,0,16);if(ec_less(m)) return;m->dx=0;ec_logic(m,0,16);}
    ec_cmp(m,m->bx,0,16);if(ec_less(m)) {ec_cmp(m,m->di,0,16);if(ec_less(m)) return;m->bx=0;ec_logic(m,0,16);}
    ec_cmp(m,m->si,640,16);if(!ec_less(m)) {ec_cmp(m,m->dx,640,16);if(!ec_less(m)) return;m->si=639;}
    ec_cmp(m,m->di,400,16);if(!ec_less(m)) {ec_cmp(m,m->bx,400,16);if(!ec_less(m)) return;m->di=399;}
    ec_push(m,m->ds);ec_push(m,m->cx);ec_push(m,m->bp);m->cx=m->dx;rc_setup(m);m->dx=m->cx;m->ax=0xa0c8;m->ds=m->ax;m->di=ec_math(m,m->di,m->bx,16,0,1);m->cx=m->dx;vg_cl(m,(uint8_t)m->cx&7);ec_logic(m,(uint8_t)m->cx,8);vg_al(m,255);vg_al(m,(uint8_t)al_shr(m,(uint8_t)m->ax,(uint8_t)m->cx,8));m->cx=m->si;vg_cl(m,(uint8_t)m->cx&7);ec_logic(m,(uint8_t)m->cx,8);vg_ah(m,128);vg_ah(m,KI_RECT_MUTATION==4?(uint8_t)al_shr(m,m->ax>>8,(uint8_t)m->cx,8):rc_sar8(m,(uint8_t)(m->ax>>8),(uint8_t)m->cx));
    vg_cl(m,3);m->dx=vg_right(m,m->dx,3);m->si=vg_right(m,m->si,3);m->si=ec_math(m,m->si,m->dx,16,0,1);vg_cl(m,4);m->bx=vg_left(m,m->bx,4);m->dx=ec_math(m,m->dx,m->bx,16,0,0);m->bx=vg_left(m,m->bx,2);m->dx=ec_math(m,m->dx,m->bx,16,0,0);if(KI_RECT_MUTATION!=3) m->di=ec_inc(m,m->di,16);m->bp=m->dx;
    do {m->bx=m->bp;m->cx=m->si;ec_cmp(m,(uint8_t)m->cx,0,8);if(m->flags&0x40) {vg_ah(m,(uint8_t)(m->ax>>8)&(uint8_t)m->ax);ec_logic(m,m->ax>>8,8);}else {if(KI_RECT_MUTATION!=5) vg_dl(m,rc_read(m,m->ds,m->bx));rc_write(m,m->ds,m->bx,(uint8_t)m->ax);m->bx=ec_inc(m,m->bx,16);vg_ch(m,255);for(;;) {vg_cl(m,(uint8_t)st_dec(m,(uint8_t)m->cx,8));if(m->flags&0x40) break;rc_write(m,m->ds,m->bx,(uint8_t)(m->cx>>8));m->bx=ec_inc(m,m->bx,16);}}
        if(KI_RECT_MUTATION!=5) vg_dl(m,rc_read(m,m->ds,m->bx));rc_write(m,m->ds,m->bx,(uint8_t)(m->ax>>8));m->bp=ec_math(m,m->bp,80,16,0,0);m->di=st_dec(m,m->di,16);
    }while(!(m->flags&0x40));m->bp=ec_pop(m);m->cx=ec_pop(m);m->ds=ec_pop(m);
}
static void rc_span(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->dx);rc_call(m,0xcac,0xae3,h);m->dx=0x3cf;vg_al(m,(uint8_t)(m->ax>>8));vg_out(m);m->dx=ec_pop(m);vg_al(m,(uint8_t)m->dx);m->dx=vg_right(m,m->dx,3);m->bx=ec_math(m,m->bx,m->dx,16,0,0);vg_al(m,(uint8_t)m->ax&7);ec_logic(m,(uint8_t)m->ax,8);m->dx=0xffff;vg_ah(m,(uint8_t)m->ax);vg_ah(m,(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->cx,8,0,0));vg_cl(m,(uint8_t)m->ax);vg_dl(m,(uint8_t)al_shr(m,(uint8_t)m->dx,(uint8_t)m->cx,8));vg_cl(m,(uint8_t)(m->ax>>8));vg_cl(m,(uint8_t)m->cx&7);ec_logic(m,(uint8_t)m->cx,8);vg_dh(m,(uint8_t)al_shr(m,m->dx>>8,(uint8_t)m->cx,8));vg_dh(m,(uint8_t)~(m->dx>>8));ec_cmp(m,m->ax>>8,8,8);if(m->flags&1) {vg_dl(m,(uint8_t)m->dx&(uint8_t)(m->dx>>8));ec_logic(m,(uint8_t)m->dx,8);}vg_cl(m,(uint8_t)(m->ax>>8));m->ax=0xa0c8;m->ds=m->ax;
    do {vg_al(m,(uint8_t)m->cx);m->si=m->bx;if(KI_RECT_MUTATION!=5) vg_ah(m,rc_read(m,m->ds,m->si));rc_write(m,m->ds,m->si,(uint8_t)m->dx);vg_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,8,8,0,1));
        if(!(m->flags&1)) {for(;;) {vg_ah(m,255);ec_cmp(m,(uint8_t)m->ax,8,8);if(m->flags&1) {m->si=ec_inc(m,m->si,16);if(KI_RECT_MUTATION!=5) vg_ah(m,rc_read(m,m->ds,m->si));rc_write(m,m->ds,m->si,(uint8_t)(m->dx>>8));break;}m->si=ec_inc(m,m->si,16);rc_write(m,m->ds,m->si,(uint8_t)(m->ax>>8));vg_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,8,8,0,1));if(m->flags&1) break;}}
        m->bx=ec_math(m,m->bx,KI_RECT_MUTATION==10?79:80,16,0,0);vg_ch(m,(uint8_t)st_dec(m,m->cx>>8,8));
    }while(!(m->flags&0x40));m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
static void rc_bar(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);m->bx=vg_left(m,m->bx,4);m->ax=m->bx;m->ax=vg_left(m,m->ax,2);m->bx=ec_math(m,m->bx,m->ax,16,0,0);m->ax=ec_pop(m);ec_push(m,m->cx);ec_logic(m,(uint8_t)m->cx,8);if(!(m->flags&0x40)) {vg_ch(m,2);rc_call(m,0xad9,0xac6,h);}
    vg_ch(m,0);ec_logic(m,0,8);m->dx=ec_math(m,m->dx,m->cx,16,0,0);m->cx=ec_pop(m);vg_ch(m,(uint8_t)ec_math(m,m->cx>>8,(uint8_t)m->cx,8,0,1));if(m->flags&0x40) return;vg_cl(m,(uint8_t)(m->cx>>8));vg_ch(m,2);vg_ah(m,(uint8_t)m->ax);rc_call(m,0xad9,0xad8,h);
}
static void rc_selection(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(t==0xc61f) {ec_push(m,m->bx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);m->bx=248;rc_test(m,(uint8_t)m->ax,8,8);if(!(m->flags&0x40)) m->bx=ec_math(m,m->bx,16,16,0,0);vg_al(m,(uint8_t)m->ax&7);ec_logic(m,(uint8_t)m->ax,8);for(unsigned i=0;i<4;++i) vg_al(m,vg_shl8(m,(uint8_t)m->ax));vg_dl(m,(uint8_t)m->ax);vg_dh(m,0);ec_logic(m,0,8);m->dx=ec_math(m,m->dx,496,16,0,0);m->dx=ec_inc(m,m->dx,16);m->bx=ec_inc(m,m->bx,16);m->si=m->dx;m->di=m->bx;m->si=ec_math(m,m->si,13,16,0,0);m->di=ec_math(m,m->di,13,16,0,0);rc_call(m,0xf020,0xc64e,h);m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->bx=ec_pop(m);return;}
    if(t==0xc6ae) {ec_push(m,m->ax);ec_push(m,m->cx);m->ax=0;ec_logic(m,0,16);m->cx=6;do {rc_call(m,0xc6bf,0xc6b8,h);vg_al(m,(uint8_t)ec_inc(m,(uint8_t)m->ax,8));--m->cx;}while(m->cx);m->cx=ec_pop(m);m->ax=ec_pop(m);return;}
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);gw_bl(m,(uint8_t)m->ax);gw_bh(m,0);ec_logic(m,0,8);m->bx=ec_shl(m,m->bx);m->dx=ec_word(m,m->cs,(uint16_t)(m->bx-0x2d16));m->si=m->dx;m->dx=ec_inc(m,m->dx,16);m->dx=ec_inc(m,m->dx,16);m->si=ec_math(m,m->si,77,16,0,0);m->bx=372;m->di=392;ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);rc_call(m,0xf020,0xc6e4,h);m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->dx=ec_inc(m,m->dx,16);m->si=st_dec(m,m->si,16);m->bx=ec_inc(m,m->bx,16);m->di=st_dec(m,m->di,16);if(KI_RECT_MUTATION!=12) rc_call(m,0xf020,0xc6f0,h);m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void rc_gauge(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    uint16_t input=m->ax;
    m->ax=st_shr(m,m->ax,16);
    if(t==0xc775) {m->ax=st_shr(m,m->ax,16);vg_cl(m,(uint8_t)m->ax);}else {m->cx=m->ax;m->ax=st_shr(m,m->ax,16);m->cx=ec_math(m,m->cx,m->ax,16,0,0);if(KI_RECT_MUTATION==7) m->cx=(uint16_t)((uint32_t)input*3/4);}
    vg_ch(m,KI_RECT_MUTATION==8?125:124);ec_cmp(m,(uint8_t)m->cx,m->cx>>8,8);if(!(m->flags&(1|0x40))) vg_cl(m,(uint8_t)(m->cx>>8));vg_ah(m,t==0xc775?12:11);vg_al(m,0);m->dx=498;rc_call(m,0xaaa,t==0xc775?0xc78d:0xc7a8,h);
}
static void rc_waiting(KiMachine16 *m,const KiEconomyHooks *h) {
    m->di=0;ec_logic(m,0,16);do {m->dx=ec_word(m,m->cs,(uint16_t)(m->di-0x2d16));m->dx=ec_math(m,m->dx,2,16,0,0);m->bx=396;vg_cl(m,ec_byte(m,m->ds,m->si));vg_ch(m,76);ec_cmp(m,(uint8_t)m->cx,m->cx>>8,8);if(!(m->flags&(1|0x40))) vg_cl(m,(uint8_t)(m->cx>>8));vg_ah(m,12);vg_al(m,0);rc_call(m,0xaaa,0xc76a,h);m->si=ec_math(m,m->si,KI_RECT_MUTATION==9?3:4,16,0,0);m->di=ec_inc(m,m->di,16);m->di=ec_inc(m,m->di,16);ec_cmp(m,m->di,12,16);}while(m->flags&1);
}
static void rc_sidebar(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);m->ds=ec_word(m,m->cs,0xd30a);m->si=9;m->ax=3;md_far(m,0xc707,h);ec_cmp(m,m->cx,480,16);int protect=0;if(m->flags&1) {ec_cmp(m,m->dx,383,16);protect=!(m->flags&1);}
    if(protect) {rc_call(m,0x1b4,0xc716,h);rc_call(m,0xc74c,0xc719,h);rc_call(m,0x1db,0xc71c,h);}else rc_call(m,0xc74c,0xc721,h);
    m->ax=ec_word(m,m->ds,4);m->bx=211;rc_call(m,0xc775,0xc72a,h);m->ax=ec_word(m,m->ds,36);m->bx=72;rc_call(m,0xc775,0xc733,h);m->ds=ec_word(m,m->cs,0xd30e);m->ax=ec_word(m,m->ds,3);m->bx=214;rc_call(m,0xc78e,0xc741,h);m->ax=ec_word(m,m->ds,0x603);m->bx=75;rc_call(m,0xc78e,0xc74a,h);m->ds=ec_pop(m);
}
static void rc_window(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);rc_call(m,0xc14,0xbd6,h);m->dx=vg_left(m,m->dx,4);m->bx=vg_left(m,m->bx,4);m->dx=ec_math(m,m->dx,8,16,0,0);m->bx=ec_math(m,m->bx,8,16,0,0);vg_al(m,(uint8_t)m->cx);vg_cl(m,4);vg_ah(m,0);ec_logic(m,0,8);m->ax=vg_left(m,m->ax,4);m->ax=ec_math(m,m->ax,17,16,0,1);m->si=m->dx;m->si=ec_math(m,m->si,m->ax,16,0,0);vg_al(m,(uint8_t)(m->cx>>8));vg_ah(m,0);ec_logic(m,0,8);m->ax=vg_left(m,m->ax,4);m->ax=ec_math(m,m->ax,17,16,0,1);m->di=m->bx;m->di=ec_math(m,m->di,m->ax,16,0,0);vg_ah(m,0);ec_logic(m,0,8);rc_call(m,0xf1a3,0xc0d,h);m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
int ki_rect_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {case 0xf020:rc_outline(m,h);break;case 0xf140:rc_vertical(m);break;case 0xf17b:rc_horizontal(m);break;case 0xf1a3:rc_fill(m);break;case 0xaaa:rc_bar(m,h);break;case 0xad9:rc_span(m,h);break;case 0xcac:rc_mode(m,0);break;case 0xcc3:rc_mode(m,1);break;case 0xc61f:case 0xc6bf:case 0xc6ae:rc_selection(m,t,h);break;case 0xc6f6:rc_sidebar(m,h);break;case 0xc74c:rc_waiting(m,h);break;case 0xc775:case 0xc78e:rc_gauge(m,t,h);break;case 0xbcd:rc_window(m,h);break;default:return 0;}return 1;
}
void ki_rect_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_rect_body(m,t,h)&&!ki_aligned_body(m,t,h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&h&&h->external) h->external(m,t,h->user);m->ip=ec_pop(m);
}
#define RC_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_rect_body(m,t,h);m->ip=ec_pop(m);}
RC_WRAP(sub_1F020,0xf020) RC_WRAP(sub_1F140,0xf140) RC_WRAP(sub_1F17B,0xf17b) RC_WRAP(sub_1F1A3,0xf1a3)
RC_WRAP(sub_10AAA,0xaaa) RC_WRAP(sub_10AD9,0xad9) RC_WRAP(sub_10CAC,0xcac) RC_WRAP(sub_10CC3,0xcc3)
RC_WRAP(sub_1C61F,0xc61f) RC_WRAP(sub_1C6BF,0xc6bf) RC_WRAP(sub_1C6AE,0xc6ae) RC_WRAP(sub_1C6F6,0xc6f6)
RC_WRAP(sub_1C74C,0xc74c) RC_WRAP(sub_1C775,0xc775) RC_WRAP(sub_1C78E,0xc78e) RC_WRAP(sub_10BCD,0xbcd)
#undef RC_WRAP
