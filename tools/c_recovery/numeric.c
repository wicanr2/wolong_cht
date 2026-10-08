/* Original DOS/V numeric editor and finance, spec/221. Include after modal.c. */
#include "numeric.h"
#ifndef KI_NUMERIC_MUTATION
#define KI_NUMERIC_MUTATION 0
#endif
static void nu_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {ec_push(m,next);m->ip=t;ki_numeric_invoke(m,t,h);}
static void nu_div(KiMachine16 *m,uint16_t divisor) {
    uint32_t n=((uint32_t)m->dx<<16)|m->ax;
    /* DIV's undefined FLAGS follow the pinned dosgolem high-half comparison model. */
    ec_cmp(m,m->dx,divisor,16);m->ax=(uint16_t)(n/divisor);m->dx=(uint16_t)(n%divisor);
}
static void nu_action(KiMachine16 *m,uint16_t t) {
    switch(t) {
    case 0x7da5: {
        gw_ah(m,0);ec_logic(m,0,8);pn_swap(&m->ax,&m->si);m->dx=10;uint32_t v=(uint32_t)m->ax*m->dx;m->ax=(uint16_t)v;m->dx=(uint16_t)(v>>16);m->si=ec_math(m,m->si,m->ax,16,0,0);unsigned carry=KI_NUMERIC_MUTATION==1?0:m->flags&1;m->dx=ec_math(m,m->dx,0,16,carry,0);ec_logic(m,m->dx,16);if(!(m->flags&0x40)) m->si=0xffff;
        ec_cmp(m,m->si,ec_word(m,m->ss,m->bp),16);if(!(m->flags&(1|0x40))) m->si=ec_word(m,m->ss,m->bp);m->flags&=(uint16_t)~1u;break;}
    case 0x7dc3: {
        m->ax=m->si;m->dx=KI_NUMERIC_MUTATION==2?10:100;uint32_t v=(uint32_t)m->ax*m->dx;m->ax=(uint16_t)v;m->dx=(uint16_t)(v>>16);m->si=m->ax;ec_logic(m,m->dx,16);if(!(m->flags&0x40)) m->si=0xffff;ec_cmp(m,m->si,ec_word(m,m->ss,m->bp),16);if(!(m->flags&(1|0x40))) m->si=ec_word(m,m->ss,m->bp);m->flags&=(uint16_t)~1u;break;}
    case 0x7ddd:m->ax=m->si;m->dx=0;ec_logic(m,0,16);m->si=KI_NUMERIC_MUTATION==3?100:10;nu_div(m,m->si);m->si=m->ax;m->flags&=(uint16_t)~1u;break;
    case 0x7dea:m->flags|=1;break;
    case 0x7dec:m->si=KI_NUMERIC_MUTATION==4?0:ec_word(m,m->ss,m->bp);m->flags&=(uint16_t)~1u;break;
    case 0x7df1:m->si=0;ec_logic(m,0,16);m->flags&=(uint16_t)~1u;break;
    }
}
static void nu_device(KiMachine16 *m,int restore,const KiEconomyHooks *h) {
    ec_push(m,m->flags);md_save4(m);
    if(!restore) {m->ax=2;md_far(m,0x1c1,h);ec_store(m,m->cs,0x1d5,m->ax);ec_store(m,m->cs,0x1d7,m->dx);ec_store(m,m->cs,0x1d9,m->bx);}
    else {ec_cmp(m,ec_word(m,m->cs,0x1d5),0,16);if(!ec_less(m)) {if(ec_greater(m)) {m->dx=ec_word(m,m->cs,0x1d7);m->bx=ec_word(m,m->cs,0x1d9);m->ax=9;md_far(m,0x208,h);}else {m->ax=1;m->bx=0;ec_logic(m,0,16);md_far(m,0x1f4,h);}}}
    md_restore4(m);uint16_t flags=ec_pop(m);if(KI_NUMERIC_MUTATION!=5) m->flags=flags;
}
static void nu_popup(KiMachine16 *m,const KiEconomyHooks *h) {
    uint8_t v=ec_byte(m,m->cs,0x9441);ec_put(m,m->cs,0x9441,(uint8_t)m->dx);gw_dl(m,v);v=ec_byte(m,m->cs,0x9443);ec_put(m,m->cs,0x9443,(uint8_t)(m->dx>>8));pn_dh(m,v);ec_push(m,m->dx);nu_call(m,0x1b4,0x93f7,h);nu_call(m,0x9409,0x93fa,h);nu_call(m,0x1db,0x93fd,h);m->dx=ec_pop(m);ec_put(m,m->cs,0x9441,(uint8_t)m->dx);ec_put(m,m->cs,0x9443,(uint8_t)(m->dx>>8));
}
static void nu_frame(KiMachine16 *m,int restore,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->dx);gw_al(m,restore?0:0x51);if(restore) ec_logic(m,0,8);pn_cl(m,4);m->bx=ec_math(m,m->bx,32,16,0,1);m->dx=ev_shifts(m,m->dx,4,1);m->bx=ev_shifts(m,m->bx,4,1);m->cx=0x507;nu_call(m,0x895d,restore?0x7d5b:0x7d21,h);
    if(!restore) {m->dx=ec_shl(m,m->dx);pn_cl(m,8);m->bx=ev_shifts(m,m->bx,8,0);m->dx=ec_math(m,m->dx,m->bx,16,0,0);m->bx=ev_shifts(m,m->bx,2,0);m->bx=ec_math(m,m->bx,m->dx,16,0,0);m->bx=ec_math(m,m->bx,0xc81,16,0,0);m->di=m->bx;m->ds=ec_word(m,m->cs,0xd50);m->si=0x600;m->ax=0x4006;nu_call(m,0xfa37,0x7d43,h);}
    m->dx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void nu_keys(KiMachine16 *m,const KiEconomyHooks *h) {
    m->flags&=(uint16_t)~0x400u;ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->dx);m->ax=m->cs;m->ds=m->ax;m->bx=ec_math(m,m->bx,16,16,0,0);m->si=0x7d93;pn_ch(m,3);
    do {pn_cl(m,6);do {ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->bx);gw_al(m,ec_byte(m,m->ds,m->si));m->si=(uint16_t)(m->si+1);m->cx=0x202;nu_call(m,0xe3d7,0x7d7b,h);m->bx=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->dx=ec_math(m,m->dx,KI_NUMERIC_MUTATION==6?8:16,16,0,0);pn_cl(m,(uint8_t)st_dec(m,(uint8_t)m->cx,8));}while((uint8_t)m->cx);
        m->dx=ec_math(m,m->dx,96,16,0,1);m->bx=ec_math(m,m->bx,16,16,0,0);pn_ch(m,(uint8_t)st_dec(m,m->cx>>8,8));}while(m->cx>>8);
    m->dx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void nu_editor(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);ec_push(m,m->es);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);ec_push(m,m->bp);m->sp=ec_math(m,m->sp,6,16,0,1);m->bp=m->sp;
    nu_call(m,0x1b4,0x7c7e,h);m->dx=ec_math(m,m->dx,8,16,0,1);m->bx=ec_math(m,m->bx,8,16,0,1);m->cx=0x5007;nu_call(m,0x9796,0x7c8a,h);m->dx=ec_math(m,m->dx,8,16,0,0);m->bx=ec_math(m,m->bx,8,16,0,0);nu_call(m,0x7d0d,0x7c93,h);nu_call(m,0x7d5f,0x7c96,h);nu_call(m,0x1db,0x7c99,h);ec_store(m,m->ss,m->bp,m->ax);ec_store(m,m->ss,(uint16_t)(m->bp+2),m->dx);ec_store(m,m->ss,(uint16_t)(m->bp+4),m->bx);m->si=0;ec_logic(m,0,16);
    for(;;) {m->ax=m->si;m->dx=0;ec_logic(m,0,16);m->bx=0xf06;nu_call(m,0x62f,0x7cae,h);nu_call(m,0x21e7,0x7cb1,h);if(m->flags&1) {m->flags|=1;break;}ec_push(m,m->cx);ec_push(m,m->dx);nu_call(m,0xe453,0x7cb8,h);m->dx=ec_pop(m);m->cx=ec_pop(m);gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,0x52,8,0,1));if(m->flags&1) continue;gw_bl(m,(uint8_t)m->ax);gw_bl(m,(uint8_t)ec_math(m,(uint8_t)m->bx,KI_NUMERIC_MUTATION==7?10:9,8,0,1));if(m->flags&1) {gw_bl(m,0);ec_logic(m,0,8);}gw_bh(m,0);ec_logic(m,0,8);m->bx=ec_shl(m,m->bx);uint16_t t=ec_word(m,m->cs,(uint16_t)(m->bx+0x7d01));nu_call(m,t,0x7cd0,h);if(m->flags&1) {m->flags&=(uint16_t)~1u;break;}}
    gw_ah(m,(uint8_t)((m->flags&0xd5)|2));nu_call(m,0x1b4,0x7cda,h);m->dx=ec_word(m,m->ss,(uint16_t)(m->bp+2));m->bx=ec_word(m,m->ss,(uint16_t)(m->bp+4));nu_call(m,0x7d47,0x7ce3,h);m->dx=ec_math(m,m->dx,8,16,0,1);m->bx=ec_math(m,m->bx,8,16,0,1);m->cx=0x5007;nu_call(m,0x97c3,0x7cef,h);m->sp=ec_math(m,m->sp,6,16,0,0);nu_call(m,0x1db,0x7cf5,h);if(KI_NUMERIC_MUTATION!=8) m->flags=(uint16_t)((m->flags&~0xd5u)|((m->ax>>8)&0xd5));m->ax=m->si;m->bp=ec_pop(m);m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->es=ec_pop(m);m->ds=ec_pop(m);
}
static void nu_finance(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    unsigned slot=t==0x67cd?0:t==0x67e6?1:t==0x6806?2:3;m->cx=(uint16_t)(17+slot);nu_call(m,0x8853,(uint16_t)(t+6),h);m->dx=296;m->bx=184;m->ax=slot?10000:100;nu_call(m,0x7c6e,(uint16_t)(t+18),h);
    if(!(m->flags&1)||KI_NUMERIC_MUTATION==9) {if(slot) {m->dx=0;ec_logic(m,0,16);m->bx=KI_NUMERIC_MUTATION==10?1:10;nu_div(m,m->bx);ec_store(m,m->cs,(uint16_t)(0xd10+slot*2),m->ax);}else ec_put(m,m->cs,0xd10,(uint8_t)m->ax);}
}
int ki_numeric_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {case 0x7c6e:nu_editor(m,h);break;case 0x7d0d:nu_frame(m,0,h);break;case 0x7d47:nu_frame(m,1,h);break;case 0x7d5f:nu_keys(m,h);break;
    case 0x7da5:case 0x7dc3:case 0x7ddd:case 0x7dea:case 0x7dec:case 0x7df1:nu_action(m,t);break;case 0x1b4:nu_device(m,0,h);break;case 0x1db:nu_device(m,1,h);break;case 0x93e9:nu_popup(m,h);break;
    case 0x67cd:case 0x67e6:case 0x6806:case 0x6826:nu_finance(m,t,h);break;default:return 0;}return 1;
}
#define NU_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_numeric_body(m,t,h);m->ip=ec_pop(m);}
NU_WRAP(sub_17C6E,0x7c6e) NU_WRAP(sub_17D0D,0x7d0d) NU_WRAP(sub_17D47,0x7d47) NU_WRAP(sub_17D5F,0x7d5f)
NU_WRAP(sub_17DA5,0x7da5) NU_WRAP(sub_17DC3,0x7dc3) NU_WRAP(sub_17DDD,0x7ddd) NU_WRAP(sub_17DEA,0x7dea) NU_WRAP(sub_17DEC,0x7dec) NU_WRAP(sub_17DF1,0x7df1)
NU_WRAP(sub_101B4,0x1b4) NU_WRAP(sub_101DB,0x1db) NU_WRAP(sub_193E9,0x93e9) NU_WRAP(sub_167CD,0x67cd) NU_WRAP(sub_167E6,0x67e6) NU_WRAP(sub_16806,0x6806) NU_WRAP(sub_16826,0x6826)
#undef NU_WRAP
void ki_numeric_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {case 0x7c6e:case 0x7d0d:case 0x7d47:case 0x7d5f:case 0x7da5:case 0x7dc3:case 0x7ddd:case 0x7dea:case 0x7dec:case 0x7df1:case 0x1b4:case 0x1db:case 0x93e9:case 0x67cd:case 0x67e6:case 0x6806:case 0x6826:if(h&&h->enter) h->enter(m,t,h->user);ki_numeric_body(m,t,h);m->ip=ec_pop(m);break;default:ki_modal_invoke(m,t,h);break;}
}
