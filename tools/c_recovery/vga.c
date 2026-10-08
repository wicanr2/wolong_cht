/* Original eight DOS/V VGA entries, spec/223. Include after settlement.c.
 * Platform bus is independent from the original CPU instance. No flat-VRAM shortcut.
 */
#include "vga.h"
#ifndef KI_VGA_MUTATION
#define KI_VGA_MUTATION 0
#endif
static void vg_al(KiMachine16 *m,uint8_t b) {m->ax=(uint16_t)((m->ax&0xff00u)|b);}
static void vg_ah(KiMachine16 *m,uint8_t b) {m->ax=(uint16_t)((m->ax&0xffu)|((uint16_t)b<<8));}
static void vg_cl(KiMachine16 *m,uint8_t b) {m->cx=(uint16_t)((m->cx&0xff00u)|b);}
static void vg_ch(KiMachine16 *m,uint8_t b) {m->cx=(uint16_t)((m->cx&0xffu)|((uint16_t)b<<8));}
static void vg_dl(KiMachine16 *m,uint8_t b) {m->dx=(uint16_t)((m->dx&0xff00u)|b);}
static void vg_dh(KiMachine16 *m,uint8_t b) {m->dx=(uint16_t)((m->dx&0xffu)|((uint16_t)b<<8));}
static void vg_swap(uint16_t *a,uint16_t *b) {uint16_t v=*a;*a=*b;*b=v;}
static uint16_t vg_left(KiMachine16 *m,uint16_t v,unsigned n) {while(n--) v=ec_shl(m,v);return v;}
static uint16_t vg_right(KiMachine16 *m,uint16_t v,unsigned n) {while(n--) v=st_shr(m,v,16);return v;}
static uint8_t vg_shl8(KiMachine16 *m,uint8_t a) {uint8_t r=(uint8_t)(a<<1);uint16_t b=(a>>7)|((((r>>7)^(a>>7))&1)<<11);m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|b),r,8);return r;}
static void vg_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiVgaHooks *h) {ec_push(m,next);m->ip=t;ki_vga_invoke(m,t,h);}
static void vg_out(KiMachine16 *m) {wolong_vga_out(m->dx,(uint8_t)m->ax);}
static void vg_movsb(KiMachine16 *m) {
    uint8_t b=wolong_vga_read(ec_at(m->ds,m->si));wolong_vga_write(ec_at(m->es,m->di),b);int delta=(m->flags&0x400)?-1:1;m->si=(uint16_t)(m->si+delta);m->di=(uint16_t)(m->di+delta);
}
static void vg_draw_row(KiMachine16 *m,int twice) {
    ec_push(m,m->cx);m->dx=m->cx;m->bp=m->bx;vg_ch(m,0);ec_logic(m,0,8);
    do {vg_dl(m,(uint8_t)m->cx);m->di=m->bp;do {
        if(KI_VGA_MUTATION!=1) vg_al(m,wolong_vga_read(ec_at(m->es,m->di)));vg_movsb(m);
        if(twice&&KI_VGA_MUTATION!=2) {if(KI_VGA_MUTATION!=1) vg_al(m,wolong_vga_read(ec_at(m->es,m->di)));vg_movsb(m);}
        vg_dl(m,(uint8_t)st_dec(m,(uint8_t)m->dx,8));
    }while((uint8_t)m->dx);m->bp=ec_math(m,m->bp,KI_VGA_MUTATION==3?79:80,16,0,0);vg_dh(m,(uint8_t)st_dec(m,m->dx>>8,8));}while(m->dx>>8);m->cx=ec_pop(m);
}
static void vg_draw(KiMachine16 *m,int twice,const KiVgaHooks *h) {
    if(KI_VGA_MUTATION!=8) m->flags&=(uint16_t)~0x400u;ec_push(m,m->es);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);ec_push(m,m->bp);m->cx=m->ax;
    m->dx=0x3ce;vg_al(m,5);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,0);ec_logic(m,0,8);vg_out(m);m->dx=st_dec(m,m->dx,16);vg_al(m,0);ec_logic(m,0,8);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,0);ec_logic(m,0,8);vg_out(m);m->dx=st_dec(m,m->dx,16);
    vg_al(m,3);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,0);ec_logic(m,0,8);vg_out(m);m->dx=st_dec(m,m->dx,16);vg_al(m,8);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,255);vg_out(m);m->dx=st_dec(m,m->dx,16);vg_al(m,1);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,14);vg_out(m);m->dx=0xa0c8;m->es=m->dx;
    uint16_t row=twice?0xfaa2:0xfa1b;vg_call(m,row,twice?0xfa72:0xf9eb,h);
    m->dx=0x3ce;vg_al(m,3);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,KI_VGA_MUTATION==4?0:16);vg_out(m);m->dx=st_dec(m,m->dx,16);vg_al(m,1);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,13);vg_out(m);vg_call(m,row,twice?0xfa87:0xfa00,h);
    m->dx=0x3cf;vg_al(m,KI_VGA_MUTATION==5?13:11);vg_out(m);vg_call(m,row,twice?0xfa90:0xfa09,h);m->dx=0x3cf;vg_al(m,7);vg_out(m);vg_call(m,row,twice?0xfa99:0xfa12,h);
    m->ax=m->cx;m->bp=ec_pop(m);m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->es=ec_pop(m);
}
static void vg_save_row(KiMachine16 *m) {
    m->ax=m->cx;m->dx=m->cx;m->bp=m->bx;vg_ch(m,0);ec_logic(m,0,8);
    do {vg_cl(m,(uint8_t)m->ax);m->si=m->bp;while(m->cx) {vg_movsb(m);--m->cx;}m->bp=ec_math(m,m->bp,KI_VGA_MUTATION==3?79:80,16,0,0);vg_dh(m,(uint8_t)st_dec(m,m->dx>>8,8));}while(m->dx>>8);m->cx=m->ax;
}
static void vg_save(KiMachine16 *m,const KiVgaHooks *h) {
    if(KI_VGA_MUTATION!=8) m->flags&=(uint16_t)~0x400u;ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);ec_push(m,m->bp);vg_al(m,vg_shl8(m,(uint8_t)m->ax));m->cx=m->ax;m->dx=0x3ce;m->si=0x3cf;
    vg_al(m,5);vg_out(m);vg_swap(&m->si,&m->dx);vg_al(m,0);ec_logic(m,0,8);vg_out(m);vg_swap(&m->si,&m->dx);vg_al(m,4);vg_out(m);vg_swap(&m->si,&m->dx);vg_al(m,0);ec_logic(m,0,8);vg_out(m);m->dx=0xa0c8;m->ds=m->dx;vg_call(m,0xfb11,0xfaee,h);
    for(unsigned plane=1;plane<=3;++plane) {m->dx=0x3cf;vg_al(m,KI_VGA_MUTATION==6?0:(uint8_t)plane);vg_out(m);vg_call(m,0xfb11,(uint16_t)(0xfaf7+(plane-1)*9),h);}
    m->bp=ec_pop(m);m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
static void vg_wrapper(KiMachine16 *m,int restore,const KiVgaHooks *h) {
    ec_push(m,restore?m->ds:m->es);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->dx);ec_push(m,restore?m->si:m->di);m->dx=vg_right(m,m->dx,KI_VGA_MUTATION==7?2:3);m->bx=vg_left(m,m->bx,4);m->dx=ec_math(m,m->dx,m->bx,16,0,0);m->bx=vg_left(m,m->bx,2);m->bx=ec_math(m,m->bx,m->dx,16,0,0);
    if(restore) {m->ds=ec_word(m,m->cs,0x987c);m->si=0;ec_logic(m,0,16);}else {m->es=ec_word(m,m->cs,0x987c);m->di=0;ec_logic(m,0,16);}m->ax=m->cx;vg_call(m,restore?0xfa37:0xfac2,restore?0x97ea:0x97bd,h);
    if(restore) m->si=ec_pop(m);else m->di=ec_pop(m);m->dx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);if(restore) m->ds=ec_pop(m);else m->es=ec_pop(m);
}
int ki_vga_body(KiMachine16 *m,uint16_t t,const KiVgaHooks *h) {
    switch(t) {case 0x9796:vg_wrapper(m,0,h);break;case 0x97c3:vg_wrapper(m,1,h);break;case 0xf9b0:vg_draw(m,0,h);break;case 0xfa1b:vg_draw_row(m,0);break;case 0xfa37:vg_draw(m,1,h);break;case 0xfaa2:vg_draw_row(m,1);break;case 0xfac2:vg_save(m,h);break;case 0xfb11:vg_save_row(m);break;default:return 0;}return 1;
}
#define VG_WRAP(n,t) void n(KiMachine16 *m,const KiVgaHooks *h) {ki_vga_body(m,t,h);m->ip=ec_pop(m);}
VG_WRAP(sub_19796,0x9796) VG_WRAP(sub_197C3,0x97c3) VG_WRAP(sub_1F9B0,0xf9b0) VG_WRAP(sub_1FA1B,0xfa1b) VG_WRAP(sub_1FA37,0xfa37) VG_WRAP(sub_1FAA2,0xfaa2) VG_WRAP(sub_1FAC2,0xfac2) VG_WRAP(sub_1FB11,0xfb11)
#undef VG_WRAP
void ki_vga_invoke(KiMachine16 *m,uint16_t t,const KiVgaHooks *h) {if(h&&h->enter) h->enter(m,t,h->user);ki_vga_body(m,t,h);m->ip=ec_pop(m);}
