/* DOS/V KI.EXE, five original entries; spec/224. Include after hotspot.c and vga.c.
 * Keep byte arithmetic, signed branches, masks, source advance and original ABI.
 */
#include "aligned.h"
#ifndef KI_ALIGNED_MUTATION
#define KI_ALIGNED_MUTATION 0
#endif
static uint16_t al_shr(KiMachine16 *m,uint16_t v,unsigned n,unsigned width) {
    while(n--) v=st_shr(m,v,width);return v;
}
static uint8_t al_read(KiMachine16 *m,uint16_t seg,uint16_t off) {
    (void)m;return wolong_vga_read(ec_at(seg,off));
}
static uint16_t al_word(KiMachine16 *m,uint16_t seg,uint16_t off) {
    (void)m;
    return (uint16_t)(wolong_vga_read(ec_at(seg,off))|((uint16_t)wolong_vga_read(ec_at(seg,(uint16_t)(off+1)))<<8));
}
static void al_write(KiMachine16 *m) {
    if(KI_ALIGNED_MUTATION!=3) vg_dh(m,al_read(m,m->es,m->di));
    wolong_vga_write(ec_at(m->es,m->di),(uint8_t)m->dx);
}
static void al_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=target;ki_aligned_invoke(m,target,h);
}
static void al_plane(KiMachine16 *m) {
    ec_push(m,m->ax);ec_push(m,m->dx);vg_ah(m,(uint8_t)m->dx);vg_al(m,1);m->dx=0x3ce;vg_out(m);
    m->dx=ec_inc(m,m->dx,16);vg_al(m,(uint8_t)(m->ax>>8));vg_out(m);m->dx=0x3ce;vg_al(m,8);vg_out(m);
    m->dx=ec_pop(m);m->ax=ec_pop(m);
}
static void al_row(KiMachine16 *m) {
    ec_push(m,m->cx);ec_push(m,m->bp);
    do {
        ec_push(m,m->si);ec_push(m,m->bx);m->di=m->bp;m->dx=0x3cf;vg_out(m);
        vg_dl(m,al_read(m,m->ds,m->si));vg_dl(m,(uint8_t)al_shr(m,(uint8_t)m->dx,(uint8_t)m->cx,8));al_write(m);
        m->di=ec_inc(m,m->di,16);gw_bh(m,(uint8_t)st_dec(m,m->bx>>8,8));
        int negative=KI_ALIGNED_MUTATION==7?((m->bx>>8)&0x80)!=0:ec_less(m);
        if(!negative) {
            if(!(m->flags&0x40)) {
                m->dx=0x3cf;gw_bl(m,(uint8_t)m->ax);vg_al(m,255);vg_out(m);vg_al(m,(uint8_t)m->bx);
                do {
                    m->dx=al_word(m,m->ds,m->si);
                    if(KI_ALIGNED_MUTATION!=5) m->dx=(uint16_t)((m->dx<<8)|(m->dx>>8));
                    m->dx=al_shr(m,m->dx,(uint8_t)m->cx,16);al_write(m);m->si=ec_inc(m,m->si,16);m->di=ec_inc(m,m->di,16);gw_bh(m,(uint8_t)st_dec(m,m->bx>>8,8));
                }while(m->bx>>8);
            }
            ec_logic(m,m->ax>>8,8);
            if(!(m->flags&0x40)) {
                m->dx=0x3cf;uint8_t a=(uint8_t)m->ax;vg_al(m,(uint8_t)(m->ax>>8));vg_ah(m,a);vg_out(m);a=(uint8_t)m->ax;vg_al(m,(uint8_t)(m->ax>>8));vg_ah(m,a);
                m->dx=al_word(m,m->ds,m->si);if(KI_ALIGNED_MUTATION!=5) m->dx=(uint16_t)((m->dx<<8)|(m->dx>>8));m->dx=al_shr(m,m->dx,(uint8_t)m->cx,16);al_write(m);
            }
        }
        m->bx=ec_pop(m);m->si=ec_pop(m);vg_dl(m,(uint8_t)m->bx);vg_dh(m,0);ec_logic(m,0,8);
        m->si=ec_math(m,m->si,KI_ALIGNED_MUTATION==4?(uint16_t)(m->dx+1):m->dx,16,0,0);m->bp=ec_math(m,m->bp,80,16,0,0);vg_ch(m,(uint8_t)st_dec(m,m->cx>>8,8));
    }while(m->cx>>8);
    m->bp=ec_pop(m);m->cx=ec_pop(m);
}
static void al_blit(KiMachine16 *m,const KiEconomyHooks *h) {
    if(KI_ALIGNED_MUTATION!=10) m->flags&=(uint16_t)~0x400u;
    ec_push(m,m->es);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->di);ec_push(m,m->bp);
    uint16_t first_source=m->si;
    m->bp=vg_right(m,m->dx,3);m->bx=vg_left(m,m->bx,4);m->bp=ec_math(m,m->bp,m->bx,16,0,0);m->bx=vg_left(m,m->bx,2);m->bp=ec_math(m,m->bp,m->bx,16,0,0);
    vg_dl(m,(uint8_t)m->dx&7);ec_logic(m,(uint8_t)m->dx,8);vg_dh(m,(uint8_t)m->cx);vg_dh(m,(uint8_t)ec_math(m,m->dx>>8,(uint8_t)m->dx,8,0,0));
    gw_bl(m,(uint8_t)m->cx);gw_bh(m,(uint8_t)(m->dx>>8));gw_bl(m,(uint8_t)ec_math(m,(uint8_t)m->bx,7,8,0,0));m->bx&=0xf8f8;ec_logic(m,m->bx,16);m->bx=vg_right(m,m->bx,3);
    vg_dh(m,(uint8_t)(m->dx>>8)&7);ec_logic(m,m->dx>>8,8);vg_cl(m,(uint8_t)(m->dx>>8));m->ax=0xffff;vg_ah(m,(uint8_t)al_shr(m,m->ax>>8,(uint8_t)m->cx,8));vg_ah(m,(uint8_t)~(m->ax>>8));
    if(KI_ALIGNED_MUTATION==2) vg_ah(m,255);
    vg_cl(m,(uint8_t)m->dx);vg_al(m,(uint8_t)al_shr(m,(uint8_t)m->ax,KI_ALIGNED_MUTATION==1?0:(uint8_t)m->cx,8));ec_logic(m,m->bx>>8,8);
    if(m->flags&0x40) {vg_al(m,(uint8_t)m->ax&(uint8_t)(m->ax>>8));ec_logic(m,(uint8_t)m->ax,8);}
    m->di=m->ax;m->dx=0x3ce;
    for(unsigned index=0;index<4;++index) {
        vg_al(m,(uint8_t[]){5,0,3,8}[index]);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,index==3?255:0);if(index!=3) ec_logic(m,0,8);vg_out(m);m->dx=st_dec(m,m->dx,16);
    }
    m->ax=m->di;m->dx=0xa0c8;m->es=m->dx;vg_dl(m,14);al_call(m,0xf999,0xf907,h);al_call(m,0xf938,0xf90a,h);
    m->di=m->ax;m->dx=0x3ce;vg_al(m,3);vg_out(m);m->dx=ec_inc(m,m->dx,16);vg_al(m,16);vg_out(m);m->ax=m->di;
    for(unsigned plane=1;plane<=3;++plane) {
        if(KI_ALIGNED_MUTATION==6) m->si=first_source;
        vg_dl(m,(uint8_t[]){0,13,11,7}[plane]);al_call(m,0xf999,(uint16_t)(0xf91d+(plane-1)*8),h);al_call(m,0xf938,(uint16_t)(0xf920+(plane-1)*8),h);
    }
    m->bp=ec_pop(m);m->di=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->es=ec_pop(m);
}
static void al_buttons(KiMachine16 *m,const KiEconomyHooks *h) {
    m->dx=0;ec_logic(m,0,16);m->bp=0xd2e4;m->di=0x7300;
    do {
        ec_push(m,m->dx);m->bx=0x170;m->cx=0x40a;vg_al(m,ec_byte(m,m->cs,m->bp));vg_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,0x15,8,0,0));al_call(m,0xe3d7,0xc80c,h);m->dx=ec_pop(m);
        m->ax=0x2005;m->bx=m->di;m->si=0x3000;al_call(m,0xfa37,0xc818,h);m->bp=ec_inc(m,m->bp,16);m->di=ec_math(m,m->di,10,16,0,0);m->dx=ec_math(m,m->dx,80,16,0,0);ec_cmp(m,m->dx,480,16);
    }while(m->flags&1);
    m->si=0x3900;m->ax=6;m->dx=4;m->bx=0x176;
    do {m->cx=0x1018;al_call(m,0xf888,0xc837,h);if(KI_ALIGNED_MUTATION==8) m->si=0x3900;m->dx=ec_math(m,m->dx,80,16,0,0);m->ax=st_dec(m,m->ax,16);}while(m->ax);
}
static void al_redraw(KiMachine16 *m,const KiEconomyHooks *h) {
    al_call(m,0x1b4,0xc676,h);ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);
    gw_bl(m,(uint8_t)(m->ax>>8));gw_bh(m,0);ec_logic(m,0,8);m->bx=KI_ALIGNED_MUTATION==9?m->bx:ec_shl(m,m->bx);m->dx=ec_word(m,m->cs,(uint16_t)(m->bx-0x2d16));
    vg_ah(m,(uint8_t)m->ax);vg_al(m,0);ec_logic(m,0,8);m->ax=st_shr(m,m->ax,16);m->si=m->ax;m->ax=st_shr(m,m->ax,16);m->si=ec_math(m,m->si,m->ax,16,0,0);m->ds=ec_word(m,m->cs,0xd48);m->dx=ec_math(m,m->dx,54,16,0,0);m->bx=0x176;m->cx=0x1018;al_call(m,0xf888,0xc6a4,h);
    m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);al_call(m,0x1db,0xc6ad,h);
}
int ki_aligned_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {case 0xf888:al_blit(m,h);break;case 0xf938:al_row(m);break;case 0xf999:al_plane(m);break;case 0xc7f4:al_buttons(m,h);break;case 0xc673:al_redraw(m,h);break;default:return 0;}return 1;
}
void ki_aligned_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_aligned_body(m,t,h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&h&&h->external) h->external(m,t,h->user);
    m->ip=ec_pop(m);
}
#define AL_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_aligned_body(m,t,h);m->ip=ec_pop(m);}
AL_WRAP(sub_1F888,0xf888) AL_WRAP(sub_1F938,0xf938) AL_WRAP(sub_1F999,0xf999) AL_WRAP(sub_1C7F4,0xc7f4) AL_WRAP(sub_1C673,0xc673)
#undef AL_WRAP
