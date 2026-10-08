/* Original DOS/V modal control, spec/220. Include after events.c.
 * Graphics/input/device leaves are explicit hooks. Original stack and FAR ABI remain.
 */
#include "modal.h"
#ifndef KI_MODAL_MUTATION
#define KI_MODAL_MUTATION 0
#endif
static void md_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_modal_invoke(m,t,h);
}
static void md_far(KiMachine16 *m,uint16_t next,const KiEconomyHooks *h) {
    uint16_t segment=ec_word(m,m->cs,(uint16_t)(next-2));ec_push(m,m->cs);ec_push(m,next);m->cs=segment;m->ip=0;
    /* Target 0 is distinguished by the actual FAR CS:IP in the complete snapshot. */
    if(h&&h->enter) h->enter(m,0,h->user);if(h&&h->external) h->external(m,0,h->user);
    m->ip=ec_pop(m);m->cs=ec_pop(m);
}
static void md_save4(KiMachine16 *m) {ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);}
static void md_restore4(KiMachine16 *m) {m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);}
static void md_window(KiMachine16 *m,int restore,const KiEconomyHooks *h) {
    md_save4(m);m->ax=restore?3:2;md_far(m,restore?0x20e2:0x2084,h);
    if(restore) {ec_store(m,m->cs,0x9886,m->cx);ec_store(m,m->cs,0x9888,m->dx);}
    m->ax=5;m->bx=0;ec_logic(m,0,16);md_far(m,restore?0x20f6:0x208e,h);m->ax=5;m->bx=1;md_far(m,restore?0x2101:0x2099,h);
    m->ax=7;m->cx=0;m->dx=restore?0x17ff:(KI_MODAL_MUTATION==1?0x27e:0x27f);md_far(m,restore?0x210f:0x20a7,h);
    m->ax=8;m->cx=0;m->dx=restore?0x101f:0x18f;md_far(m,restore?0x211d:0x20b5,h);
    m->cx=ec_word(m,m->cs,0x9886);m->dx=ec_word(m,m->cs,0x9888);
    if(restore) {m->cx=ec_math(m,ec_word(m,m->cs,0x9882),m->cx,16,0,KI_MODAL_MUTATION==10);m->dx=ec_math(m,ec_word(m,m->cs,0x9884),m->dx,16,0,0);}
    m->ax=4;md_far(m,restore?0x2139:0x20c7,h);m->ax=restore?2:1;
    if(!restore) {m->bx=0;ec_logic(m,0,16);}md_far(m,restore?0x2141:0x20d1,h);
    if(restore) {m->ax=0xffff;ec_store(m,m->cs,0x988a,m->ax);ec_store(m,m->cs,0x988c,m->ax);}md_restore4(m);
}
static void md_actor(KiMachine16 *m) {
    m->bx=ec_word(m,m->cs,0xcfd);gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->bx+1)));gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,1);
}
static void md_wait(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,ec_word(m,m->cs,0x2239));ec_store(m,m->cs,0x2239,KI_MODAL_MUTATION==7?0x9091:0x9090);md_call(m,0x222b,0x2225,h);ec_store(m,m->cs,0x2239,ec_pop(m));
}
static void md_speech(KiMachine16 *m,int advisor,const KiEconomyHooks *h) {
    md_save4(m);ec_push(m,m->di);
    if(advisor) {m->bx=ec_word(m,m->cs,0xcfd);gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->bx+2)));gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,1);gw_ah(m,0);ec_logic(m,0,8);gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4241)));m->dx=8;m->bx=18;}
    else {m->dx=0;m->bx=5;gw_ah(m,0);ec_logic(m,0,8);gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x1e)));ec_cmp(m,(uint8_t)m->ax,3,8);if(!(m->flags&1)) {gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,3,8,0,1));if(KI_MODAL_MUTATION==9) gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,3,8,0,1));}}
    if(!advisor) {m->cx=ec_math(m,m->cx,m->ax,16,0,0);gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->si+1)));}
    md_call(m,0x75b,advisor?0x3d00:0x3cb7,h);md_call(m,0x2216,advisor?0x3d03:0x3cba,h);m->di=ec_pop(m);md_restore4(m);
}
static void md_frame(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);gw_bl(m,ec_byte(m,m->cs,0xcf4));gw_bh(m,0);ec_logic(m,0,8);m->bx=st_dec(m,m->bx,16);gw_al(m,ec_byte(m,m->cs,(uint16_t)(m->bx+0x9309)));md_call(m,0x241,0x9333,h);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void md_render(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    ec_push(m,m->ds);md_save4(m);ec_push(m,m->si);if(t!=0x3d68) ec_push(m,m->di);
    if(t==0x3d09) {m->dx=m->cs;m->ds=m->dx;gw_ah(m,0);ec_logic(m,0,8);m->dx=0x6300;uint32_t v=(uint32_t)m->ax*m->dx;m->ax=(uint16_t)v;m->dx=(uint16_t)(v>>16);m->cx=m->dx;m->dx=0xe07;m->bx=ec_word(m,m->cs,0x9876);m->si=0;ec_logic(m,0,16);m->di=0x6300;md_call(m,0xe38c,0x3d2d,h);md_call(m,0x3d68,0x3d30,h);gw_al(m,1);m->dx=0;ec_logic(m,0,16);m->bx=2;m->cx=0x151b;md_call(m,0x89a4,0x3d3d,h);}
    else if(t==0x3d45) {m->ax=m->cs;m->ds=m->ax;md_call(m,0x87af,0x3d53,h);gw_al(m,0);ec_logic(m,0,8);m->dx=0;ec_logic(m,0,16);m->bx=2;m->cx=0x151b;md_call(m,0x89a4,0x3d60,h);}
    else {m->dx=3;m->bx=8;m->cx=0xc13;md_call(m,0xc14,0x3d7a,h);m->ds=ec_word(m,m->cs,0x9876);m->si=0;ec_logic(m,0,16);m->bx=0x2a87;m->ax=0xb012;md_call(m,0xfa37,0x3d8a,h);}
    if(t!=0x3d68) m->di=ec_pop(m);m->si=ec_pop(m);md_restore4(m);m->ds=ec_pop(m);
}
static void md_resume(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);ec_push(m,m->es);md_save4(m);ec_push(m,m->si);ec_push(m,m->di);ec_push(m,m->bp);m->ax=m->cs;m->ds=m->ax;ec_store(m,m->ds,0x9886,m->cx);ec_store(m,m->ds,0x9888,m->dx);ec_push(m,m->cx);ec_push(m,m->dx);md_call(m,0xe453,0x1d60,h);m->dx=ec_pop(m);m->cx=ec_pop(m);ec_logic(m,(uint8_t)m->ax,8);if(m->flags&0x40) md_call(m,0x1f30,0x1d69,h);ec_logic(m,ec_byte(m,m->ds,0x98a6)&4,8);if(!(m->flags&0x40)) md_call(m,0x5c58,0x1d73,h);
    m->dx=ec_word(m,m->ds,0x988e);m->bx=ec_word(m,m->ds,0x9890);md_call(m,0xd615,0x1d7e,h);md_call(m,0x1cc9,0x1d81,h);md_call(m,0xd66a,0x1d84,h);m->bp=ec_pop(m);m->di=ec_pop(m);m->si=ec_pop(m);md_restore4(m);m->es=ec_pop(m);m->ds=ec_pop(m);
}
static void md_selector(KiMachine16 *m,const KiEconomyHooks *h) {
    m->dx=80;m->bx=176;m->cx=0x600a;md_call(m,0x9796,0x3b8a,h);gw_bl(m,(uint8_t)m->ax);
    do {gw_al(m,(uint8_t)m->bx);gw_ah(m,0);ec_logic(m,0,8);m->cx=ec_word(m,m->ss,m->bp);gw_dl(m,5);pn_dh(m,11);md_call(m,0x93e9,0x3b9a,h);}while((m->flags&1)&&KI_MODAL_MUTATION!=6);
    m->dx=80;m->bx=176;m->cx=0x600a;md_call(m,0x97c3,0x3ba8,h);
}
static void md_reply(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->si);ec_push(m,m->di);gw_bl(m,0x93);ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);
    if(!(m->flags&0x40)) {gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x2a)));gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,1);gw_bl(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4241)));m->di=m->si;}
    ec_push(m,m->ax);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->di);ec_cmp(m,(uint8_t)m->ax,3,8);if(m->flags&0x40) gw_al(m,2);m->di=m->sp;gw_ah(m,0);ec_logic(m,0,8);m->cx=ec_math(m,m->cx,m->ax,16,0,0);gw_al(m,(uint8_t)m->bx);md_call(m,0x8810,0x3c71,h);m->di=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->ax=ec_pop(m);ec_cmp(m,(uint8_t)m->ax,3,8);
    if(m->flags&0x40) {gw_al(m,KI_MODAL_MUTATION==8?29:30);m->cx=0x1a5;md_call(m,0x3dc9,0x3c93,h);}
    else {ec_logic(m,m->ax>>8,8);if(!(m->flags&0x40)) {ec_cmp(m,(uint8_t)m->ax,2,8);if(!(m->flags&0x40)) {m->cx=ec_math(m,m->cx,29,16,0,0);gw_al(m,0x93);md_call(m,0x8810,0x3c89,h);}}}
    m->di=ec_pop(m);m->si=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void md_number(KiMachine16 *m,int amount,const KiEconomyHooks *h) {
    m->ax=1;m->bx=0;ec_logic(m,0,16);md_far(m,amount?0x3a6e:0x3965,h);m->dx=88;m->bx=184;m->ax=30000;md_call(m,0x7c6e,amount?0x3a7a:0x3971,h);ec_push(m,m->flags);ec_push(m,m->ax);m->ax=2;md_far(m,amount?0x3a84:0x397b,h);m->ax=ec_pop(m);uint16_t saved=ec_pop(m);if(KI_MODAL_MUTATION!=5) m->flags=saved;
}
static void md_diplomacy(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);ec_push(m,m->es);ec_push(m,m->cx);ec_push(m,m->bp);m->sp=ec_math(m,m->sp,14,16,0,1);m->bp=m->sp;
    ec_store(m,m->ss,m->bp,m->cx);ec_put(m,m->ss,(uint16_t)(m->bp+2),(uint8_t)m->ax);ec_store(m,m->ss,(uint16_t)(m->bp+4),m->di);ec_store(m,m->ss,(uint16_t)(m->bp+6),m->bx);ec_store(m,m->ss,(uint16_t)(m->bp+8),m->dx);ec_store(m,m->ss,(uint16_t)(m->bp+10),m->dx);ec_put(m,m->ss,(uint16_t)(m->bp+12),(uint8_t)(m->ax>>8));
    m->di=m->bp;m->di=ec_math(m,m->di,4,16,0,0);m->ax=2;md_far(m,0x392d,h);md_call(m,0xcde,0x3930,h);gw_al(m,0x93);m->cx=42;md_call(m,0x8810,0x3938,h);md_call(m,0x2c2,0x393b,h);gw_al(m,0);md_call(m,0x3d09,0x3940,h);gw_al(m,6);md_call(m,0x241,0x3945,h);m->cx=ec_word(m,m->ss,m->bp);md_call(m,0x3c99,0x394b,h);ec_store(m,m->ss,m->bp,ec_math(m,ec_word(m,m->ss,m->bp),3,16,0,0));
    for(;;) {gw_al(m,3);md_call(m,0x3b7e,0x3954,h);ec_put(m,m->ss,(uint16_t)(m->bp+3),(uint8_t)m->ax);ec_cmp(m,(uint8_t)m->ax,1,8);if(!(m->flags&0x40)) break;md_number(m,0,h);if(m->flags&1) continue;ec_store(m,m->ss,(uint16_t)(m->bp+8),m->ax);ec_logic(m,m->ax,16);if(m->flags&0x40) ec_put(m,m->ss,(uint16_t)(m->bp+3),0);break;}
    gw_al(m,ec_byte(m,m->ss,(uint16_t)(m->bp+3)));ec_store(m,m->ss,m->bp,ec_inc(m,ec_word(m,m->ss,m->bp),16));m->cx=ec_word(m,m->ss,m->bp);gw_ah(m,0);ec_logic(m,0,8);m->cx=ec_math(m,m->cx,m->ax,16,0,0);md_call(m,0x3cdc,0x399a,h);ec_store(m,m->ss,m->bp,ec_math(m,ec_word(m,m->ss,m->bp),3,16,0,0));md_call(m,0xece0,0x39a1,h);ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->cs,0xd00),8);
    if((m->flags&1)||((m->flags&0x40)&&KI_MODAL_MUTATION!=2)) {ec_store(m,m->ss,m->bp,ec_math(m,ec_word(m,m->ss,m->bp),3,16,0,0));gw_al(m,ec_byte(m,m->ss,(uint16_t)(m->bp+3)));m->dx=ec_word(m,m->ss,(uint16_t)(m->bp+8));ec_cmp(m,m->dx,ec_word(m,m->ss,(uint16_t)(m->bp+10)),16);if(!(m->flags&(1|0x40))) gw_al(m,3);ec_put(m,m->ss,(uint16_t)(m->bp+2),(uint8_t)m->ax);ec_store(m,m->ss,(uint16_t)(m->bp+10),m->dx);}
    m->cx=ec_word(m,m->ss,m->bp);md_call(m,0x3c99,0x39c5,h);md_call(m,0x2c2,0x39c8,h);md_call(m,0x3d45,0x39cb,h);md_call(m,0x1d46,0x39ce,h);gw_al(m,ec_byte(m,m->ss,(uint16_t)(m->bp+2)));m->di=ec_word(m,m->ss,(uint16_t)(m->bp+4));m->bx=ec_word(m,m->ss,(uint16_t)(m->bp+6));m->dx=ec_word(m,m->ss,(uint16_t)(m->bp+10));gw_ah(m,ec_byte(m,m->ss,(uint16_t)(m->bp+12)));m->sp=ec_math(m,m->sp,14,16,0,0);md_call(m,0x9321,0x39e3,h);m->bp=ec_pop(m);m->cx=ec_pop(m);m->es=ec_pop(m);m->ds=ec_pop(m);
}
static void md_amount(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);ec_push(m,m->es);md_save4(m);ec_push(m,m->si);ec_push(m,m->di);ec_push(m,m->bp);ec_logic(m,m->ax,16);
    if(!(m->flags&0x40)) {ec_cmp(m,m->ax,KI_MODAL_MUTATION==3?499:500,16);if(m->flags&1) m->ax=KI_MODAL_MUTATION==3?499:500;}
    m->sp=ec_math(m,m->sp,12,16,0,1);m->bp=m->sp;m->di=m->sp;m->di=ec_math(m,m->di,4,16,0,0);ec_store(m,m->ss,m->bp,m->cx);ec_store(m,m->ss,(uint16_t)(m->bp+4),m->bx);ec_store(m,m->ss,(uint16_t)(m->bp+6),m->si);ec_store(m,m->ss,(uint16_t)(m->bp+8),m->ax);ec_store(m,m->ss,(uint16_t)(m->bp+10),m->ax);m->ax=2;md_far(m,0x3a1e,h);md_call(m,0x2c2,0x3a21,h);gw_al(m,1);md_call(m,0x3d09,0x3a26,h);gw_al(m,6);md_call(m,0x241,0x3a2b,h);ec_cmp(m,ec_word(m,m->ss,(uint16_t)(m->bp+8)),0,16);
    if(m->flags&0x40) {m->cx=ec_math(m,ec_word(m,m->ss,m->bp),30,16,0,0);ec_push(m,m->cx);md_call(m,0x3c99,0x3a3b,h);m->cx=ec_pop(m);m->cx=ec_math(m,m->cx,5,16,0,0);ec_push(m,m->cx);md_call(m,0x3cdc,0x3a43,h);m->cx=ec_pop(m);m->cx=ec_inc(m,m->cx,16);md_call(m,0x3c99,0x3a48,h);}
    else {
        m->cx=ec_word(m,m->ss,m->bp);md_call(m,0x3c99,0x3a51,h);ec_store(m,m->ss,m->bp,ec_math(m,ec_word(m,m->ss,m->bp),5,16,0,0));
        for(;;) {gw_al(m,3);md_call(m,0x3b7e,0x3a5a,h);ec_put(m,m->ss,(uint16_t)(m->bp+2),(uint8_t)m->ax);ec_put(m,m->ss,(uint16_t)(m->bp+3),(uint8_t)m->ax);ec_cmp(m,(uint8_t)m->ax,1,8);if(!(m->flags&0x40)) break;md_number(m,1,h);if(m->flags&1) continue;ec_store(m,m->ss,(uint16_t)(m->bp+8),m->ax);m->bx=ec_word(m,m->ss,(uint16_t)(m->bp+10));ec_logic(m,m->ax,16);
            if(m->flags&0x40) ec_store(m,m->ss,(uint16_t)(m->bp+2),0x202);
            else {ec_cmp(m,m->ax,m->bx,16);if(m->flags&0x40) ec_store(m,m->ss,(uint16_t)(m->bp+2),0);else if(!(m->flags&1)) ec_store(m,m->ss,(uint16_t)(m->bp+2),0x303);}break;}
        ec_store(m,m->ss,m->bp,ec_inc(m,ec_word(m,m->ss,m->bp),16));pn_cl(m,ec_byte(m,m->ss,(uint16_t)(m->bp+2)));pn_ch(m,0);ec_logic(m,0,8);m->cx=ec_math(m,m->cx,ec_word(m,m->ss,m->bp),16,0,0);md_call(m,0x3cdc,0x3ab9,h);ec_store(m,m->ss,m->bp,ec_math(m,ec_word(m,m->ss,m->bp),4,16,0,0));m->cx=ec_word(m,m->ss,m->bp);gw_al(m,ec_byte(m,m->ss,(uint16_t)(m->bp+3)));gw_ah(m,0);ec_logic(m,0,8);m->cx=ec_math(m,m->cx,m->ax,16,0,0);m->ax=ev_shifts(m,m->ax,2,0);m->cx=ec_math(m,m->cx,m->ax,16,0,0);md_call(m,0x3c99,0x3ad0,h);
    }
    md_call(m,0x2c2,0x3ad3,h);md_call(m,0x3d45,0x3ad6,h);md_call(m,0x1d46,0x3ad9,h);ec_cmp(m,ec_byte(m,m->ss,(uint16_t)(m->bp+2)),2,8);
    if(!(m->flags&0x40)||KI_MODAL_MUTATION==4) {m->ax=ec_word(m,m->ss,(uint16_t)(m->bp+8));m->ax=ec_shl(m,m->ax);ec_put(m,m->ds,(uint16_t)(m->si+0x1a),(uint8_t)(m->ax>>8));m->ax=st_shr(m,m->ax,16);gw_dl(m,0);ec_logic(m,0,8);m->si=ec_word(m,m->cs,0xcfd);md_call(m,0x563b,0x3af3,h);gw_al(m,4);md_call(m,0x5e80,0x3af8,h);}
    m->sp=ec_math(m,m->sp,12,16,0,0);md_call(m,0x9321,0x3afe,h);m->bp=ec_pop(m);m->di=ec_pop(m);m->si=ec_pop(m);md_restore4(m);m->es=ec_pop(m);m->ds=ec_pop(m);
}
static void md_wrapper(KiMachine16 *m,int cooperate,const KiEconomyHooks *h) {
    if(!cooperate) ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->si);if(cooperate) ec_push(m,m->bx);md_call(m,0x87ff,cooperate?0x38ec:0x38cd,h);m->bx=ec_math(m,m->bx,0x4240,16,0,0);m->si=m->bx;if(cooperate) m->bx=ec_pop(m);m->cx=cooperate?0x175:0x168;if(!cooperate) m->bx=0xffff;
    md_call(m,0x2078,cooperate?0x38f9:0x38dc,h);md_call(m,0x3902,cooperate?0x38fc:0x38df,h);md_call(m,0x20d6,cooperate?0x38ff:0x38e2,h);m->si=ec_pop(m);m->cx=ec_pop(m);if(!cooperate) m->bx=ec_pop(m);
}
int ki_modal_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {case 0x2078:md_window(m,0,h);break;case 0x20d6:md_window(m,1,h);break;case 0x38c7:md_wrapper(m,0,h);break;case 0x38e6:md_wrapper(m,1,h);break;case 0x3902:md_diplomacy(m,h);break;case 0x39e8:md_amount(m,h);break;case 0x3c3d:md_reply(m,h);break;case 0x3b7e:md_selector(m,h);break;
    case 0x3d09:case 0x3d45:case 0x3d68:md_render(m,t,h);break;case 0x3c99:md_speech(m,0,h);break;case 0x3cdc:md_speech(m,1,h);break;case 0x9321:md_frame(m,h);break;case 0x87ff:md_actor(m);break;case 0x1d46:md_resume(m,h);break;case 0x2216:md_wait(m,h);break;default:return 0;}return 1;
}
#define MD_WRAP(name,t) void name(KiMachine16 *m,const KiEconomyHooks *h) {ki_modal_body(m,t,h);m->ip=ec_pop(m);}
MD_WRAP(sub_12078,0x2078) MD_WRAP(sub_120D6,0x20d6) MD_WRAP(sub_138C7,0x38c7) MD_WRAP(sub_138E6,0x38e6)
MD_WRAP(sub_13902,0x3902) MD_WRAP(sub_139E8,0x39e8) MD_WRAP(sub_13C3D,0x3c3d) MD_WRAP(sub_13B7E,0x3b7e)
MD_WRAP(sub_13D09,0x3d09) MD_WRAP(sub_13D45,0x3d45) MD_WRAP(sub_13C99,0x3c99) MD_WRAP(sub_13CDC,0x3cdc)
MD_WRAP(sub_19321,0x9321) MD_WRAP(sub_187FF,0x87ff) MD_WRAP(sub_13D68,0x3d68) MD_WRAP(sub_11D46,0x1d46) MD_WRAP(sub_12216,0x2216)
#undef MD_WRAP
void ki_modal_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {case 0x2078:case 0x20d6:case 0x38c7:case 0x38e6:case 0x3902:case 0x39e8:case 0x3c3d:case 0x3b7e:case 0x3d09:case 0x3d45:case 0x3c99:case 0x3cdc:case 0x9321:case 0x87ff:case 0x3d68:case 0x1d46:case 0x2216:if(h&&h->enter) h->enter(m,t,h->user);ki_modal_body(m,t,h);m->ip=ec_pop(m);break;default:ki_events_invoke(m,t,h);break;}
}
