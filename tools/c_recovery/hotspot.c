/* Original DOS/V hotspot memory path, spec/222. Include after numeric.c. */
#include "hotspot.h"
#ifndef KI_HOTSPOT_MUTATION
#define KI_HOTSPOT_MUTATION 0
#endif
static void hs_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {ec_push(m,next);m->ip=t;ki_hotspot_invoke(m,t,h);}
static void hs_init(KiMachine16 *m) {
    m->flags&=(uint16_t)~0x400u;ec_push(m,m->ax);ec_push(m,m->cx);ec_store(m,m->cs,0xe479,m->es);ec_store(m,m->cs,0xe47b,m->di);m->cx=KI_HOTSPOT_MUTATION==1?1999:2000;m->ax=0;ec_logic(m,0,16);
    do {ec_put(m,m->es,m->di,(uint8_t)m->ax);ec_put(m,m->es,(uint16_t)(m->di+1),(uint8_t)(m->ax>>8));m->di=(uint16_t)(m->di+2);--m->cx;}while(m->cx);m->cx=ec_pop(m);m->ax=ec_pop(m);
}
static void hs_address(KiMachine16 *m) {
    m->dx=ev_shifts(m,m->dx,3,1);gw_bl(m,(uint8_t)m->bx&0xf8);ec_logic(m,(uint8_t)m->bx,8);if(KI_HOTSPOT_MUTATION==2) gw_bh(m,0);m->bx=ec_shl(m,m->bx);m->dx=ec_math(m,m->dx,m->bx,16,0,0);m->bx=ev_shifts(m,m->bx,2,0);m->bx=ec_math(m,m->bx,m->dx,16,0,0);m->bx=ec_math(m,m->bx,ec_word(m,m->cs,0xe47b),16,0,0);
}
static void hs_rectangle(KiMachine16 *m,int clear) {
    ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->si);m->ds=ec_word(m,m->cs,0xe479);hs_address(m);
    if(clear) {gw_al(m,0);ec_logic(m,0,8);}else {gw_ah(m,0);ec_logic(m,0,8);}
    do {ec_push(m,m->cx);m->si=m->bx;do {if(!clear) {ec_cmp(m,ec_byte(m,m->ds,m->si),0,8);if(!(m->flags&0x40)&&KI_HOTSPOT_MUTATION!=3) gw_ah(m,1);}ec_put(m,m->ds,m->si,(uint8_t)m->ax);m->si=ec_inc(m,m->si,16);pn_cl(m,(uint8_t)st_dec(m,(uint8_t)m->cx,8));}while((uint8_t)m->cx);
        m->cx=ec_pop(m);m->bx=ec_math(m,m->bx,KI_HOTSPOT_MUTATION==4?79:80,16,0,0);pn_ch(m,(uint8_t)st_dec(m,m->cx>>8,8));}while(m->cx>>8);
    if(!clear) {ec_logic(m,m->ax>>8,8);if(!(m->flags&0x40)) m->flags|=1;}m->si=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
static void hs_query(KiMachine16 *m) {
    ec_push(m,m->ds);ec_push(m,m->bx);m->ds=ec_word(m,m->cs,0xe479);m->cx=ev_shifts(m,m->cx,3,1);gw_dl(m,(uint8_t)m->dx&0xf8);ec_logic(m,(uint8_t)m->dx,8);if(KI_HOTSPOT_MUTATION==5) pn_dh(m,0);m->dx=ec_shl(m,m->dx);m->cx=ec_math(m,m->cx,m->dx,16,0,0);m->dx=ev_shifts(m,m->dx,2,0);m->cx=ec_math(m,m->cx,m->dx,16,0,0);m->bx=m->cx;m->bx=ec_math(m,m->bx,ec_word(m,m->cs,0xe47b),16,0,0);gw_al(m,ec_byte(m,m->ds,m->bx));m->bx=ec_pop(m);m->ds=ec_pop(m);
}
static void hs_flags(KiMachine16 *m) {
    ec_push(m,m->ds);m->bx=ev_shifts(m,m->bx,3,0);m->dx=ec_math(m,m->dx,m->bx,16,0,0);m->bx=ev_shifts(m,m->bx,2,0);m->dx=ec_math(m,m->dx,m->bx,16,0,0);m->dx=ev_shifts(m,m->dx,3,0);m->ds=ec_word(m,m->cs,0xd84e);
    do {gw_al(m,(uint8_t)m->cx);m->bx=m->dx;do {ec_logic(m,m->ax>>8,8);uint8_t b=ec_byte(m,m->ds,m->bx);
        if(m->flags&0x40) {b|=KI_HOTSPOT_MUTATION==6?0:0x10;ec_logic(m,b,8);ec_put(m,m->ds,m->bx,b);b&=0x7f;ec_logic(m,b,8);ec_put(m,m->ds,m->bx,b);}
        else {b&=0xef;ec_logic(m,b,8);ec_put(m,m->ds,m->bx,b);b|=0x60;ec_logic(m,b,8);ec_put(m,m->ds,m->bx,b);}
        m->bx=ec_math(m,m->bx,8,16,0,0);gw_al(m,(uint8_t)st_dec(m,(uint8_t)m->ax,8));}while((uint8_t)m->ax);m->dx=ec_math(m,m->dx,320,16,0,0);pn_ch(m,(uint8_t)st_dec(m,m->cx>>8,8));}while(m->cx>>8);m->ds=ec_pop(m);
}
static void hs_horizontal(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);m->ax=0x801;pn_cl(m,(uint8_t)m->dx);m->bx=ec_inc(m,m->bx,16);
    do {m->si=0;ec_logic(m,0,16);hs_call(m,0xf9b0,0xc6c,h);m->bx=ec_inc(m,m->bx,16);m->si=0;ec_logic(m,0,16);hs_call(m,0xf9b0,0xc72,h);m->bx=ec_inc(m,m->bx,16);--m->cx;}while(m->cx);m->ax=ec_pop(m);
}
static void hs_vertical(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);m->ax=0x801;pn_cl(m,(uint8_t)(m->dx>>8));pn_cl(m,(uint8_t)ec_math(m,(uint8_t)m->cx,1,8,0,1));m->si=32;hs_call(m,0xf9b0,0xc86,h);m->bx=ec_math(m,m->bx,KI_HOTSPOT_MUTATION==7?639:640,16,0,0);
    do {m->si=64;hs_call(m,0xf9b0,0xc90,h);m->bx=ec_math(m,m->bx,KI_HOTSPOT_MUTATION==7?639:640,16,0,0);m->si=96;hs_call(m,0xf9b0,0xc9a,h);m->bx=ec_math(m,m->bx,KI_HOTSPOT_MUTATION==7?639:640,16,0,0);--m->cx;}while(m->cx);
    m->si=32;hs_call(m,0xf9b0,0xca6,h);m->bx=ec_math(m,m->bx,KI_HOTSPOT_MUTATION==7?639:640,16,0,0);m->ax=ec_pop(m);
}
static void hs_border(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);md_save4(m);ec_push(m,m->si);ec_push(m,m->di);ec_push(m,m->bp);m->ds=ec_word(m,m->cs,0xd4a);m->dx=ec_shl(m,m->dx);m->bx=(uint16_t)((m->bx<<8)|(m->bx>>8));m->dx=ec_math(m,m->dx,m->bx,16,0,0);m->bx=ev_shifts(m,m->bx,2,0);m->bx=ec_math(m,m->bx,m->dx,16,0,0);m->bp=m->bx;m->ax=0x801;m->dx=m->cx;gw_dl(m,(uint8_t)st_dec(m,(uint8_t)m->dx,8));pn_ch(m,0);ec_logic(m,0,8);
    m->bx=m->bp;hs_call(m,0xc60,0xc3d,h);m->bx=m->bp;hs_call(m,0xc77,0xc42,h);pn_swap(&m->bx,&m->bp);pn_cl(m,(uint8_t)m->dx);m->cx=ec_shl(m,m->cx);m->cx=ec_inc(m,m->cx,16);m->bx=ec_math(m,m->bx,m->cx,16,0,0);hs_call(m,0xc77,0xc4e,h);m->bx=m->bp;m->bx=ec_math(m,m->bx,640,16,0,1);hs_call(m,0xc60,0xc57,h);
    m->bp=ec_pop(m);m->di=ec_pop(m);m->si=ec_pop(m);md_restore4(m);m->ds=ec_pop(m);
}
static void hs_wrapper(KiMachine16 *m,const KiEconomyHooks *h) {
    md_save4(m);ec_push(m,m->si);ec_push(m,m->di);ec_push(m,m->bp);m->si=m->bx;m->di=m->dx;m->bp=m->cx;gw_ah(m,0);ec_logic(m,0,8);ec_push(m,m->ax);ec_logic(m,(uint8_t)m->ax,8);if(m->flags&0x40) gw_ah(m,(uint8_t)ec_inc(m,m->ax>>8,8));hs_call(m,0xd5d4,0x8976,h);m->ax=ec_pop(m);m->si=ec_math(m,m->si,KI_HOTSPOT_MUTATION==8?1:2,16,0,0);ec_logic(m,(uint8_t)m->ax,8);
    if(!(m->flags&0x40)) {m->bx=m->si;m->dx=m->di;m->cx=m->bp;hs_call(m,0xc14,0x8987,h);}m->bx=m->si;m->dx=m->di;m->cx=m->bp;m->bx=ec_shl(m,m->bx);m->dx=ec_shl(m,m->dx);m->cx=ec_shl(m,m->cx);ec_logic(m,(uint8_t)m->ax,8);if(!(m->flags&0x40)) gw_ah(m,(uint8_t)m->ax);hs_call(m,0x89de,0x899c,h);m->bp=ec_pop(m);m->di=ec_pop(m);m->si=ec_pop(m);md_restore4(m);
}
int ki_hotspot_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {case 0xe3c0:hs_init(m);break;case 0xe3d7:hs_rectangle(m,0);break;case 0xe41b:hs_rectangle(m,1);break;case 0xe453:hs_query(m);break;case 0xd5d4:hs_flags(m);break;case 0x895d:hs_wrapper(m,h);break;
    case 0x89de:m->dx=ev_shifts(m,m->dx,3,0);m->bx=ev_shifts(m,m->bx,3,0);gw_al(m,(uint8_t)(m->ax>>8));hs_call(m,0xe3d7,0x89ef,h);break;case 0xc14:hs_border(m,h);break;case 0xc60:hs_horizontal(m,h);break;case 0xc77:hs_vertical(m,h);break;default:return 0;}return 1;
}
#define HS_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_hotspot_body(m,t,h);m->ip=ec_pop(m);}
HS_WRAP(sub_1E3C0,0xe3c0) HS_WRAP(sub_1E3D7,0xe3d7) HS_WRAP(sub_1E41B,0xe41b) HS_WRAP(sub_1E453,0xe453) HS_WRAP(sub_1895D,0x895d) HS_WRAP(sub_189DE,0x89de) HS_WRAP(sub_1D5D4,0xd5d4) HS_WRAP(sub_10C14,0xc14) HS_WRAP(sub_10C60,0xc60) HS_WRAP(sub_10C77,0xc77)
#undef HS_WRAP
void ki_hotspot_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {case 0xe3c0:case 0xe3d7:case 0xe41b:case 0xe453:case 0x895d:case 0x89de:case 0xd5d4:case 0xc14:case 0xc60:case 0xc77:if(h&&h->enter) h->enter(m,t,h->user);ki_hotspot_body(m,t,h);m->ip=ec_pop(m);break;default:ki_numeric_invoke(m,t,h);break;}
}
