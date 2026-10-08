/* Original DOS/V entries, spec/219. Include after hourly.c.
 * Raw byte semantics and original ABI; modal/battle callees remain explicit hooks.
 */
#include "events.h"
#ifndef KI_EVENTS_MUTATION
#define KI_EVENTS_MUTATION 0
#endif
static void ev_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=target;ki_events_invoke(m,target,h);
}
static uint16_t ev_shifts(KiMachine16 *m,uint16_t v,unsigned n,int right) {
    while(n--) v=right?st_shr(m,v,16):ec_shl(m,v);return v;
}
static uint8_t ev_shl8(KiMachine16 *m,uint8_t a) {
    uint8_t r=(uint8_t)(a<<1);uint16_t bits=(a>>7)|((((r>>7)^(a>>7))&1)<<11);
    m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|bits),r,8);return r;
}
static void ev_presence(KiMachine16 *m) {
    gw_al(m,0);ec_logic(m,0,8);m->ax=ev_shifts(m,m->ax,2,1);m->si=m->ax;
    ec_cmp(m,ec_byte(m,m->ds,m->si),KI_EVENTS_MUTATION==1?0x81:0x80,8);
}
static void ev_pair(KiMachine16 *m) {
    gw_bl(m,(uint8_t)m->ax);gw_bh(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,0);m->cx=m->bx;m->bx=ec_shl(m,m->bx);m->bx=ec_math(m,m->bx,m->cx,16,0,0);
    gw_bl(m,(uint8_t)ec_math(m,(uint8_t)m->bx,m->ax>>8,8,0,0));gw_bh(m,(uint8_t)ec_math(m,m->bx>>8,0,8,m->flags&1,0));m->si=m->bx;
    gw_bl(m,(uint8_t)(m->ax>>8));gw_bh(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,0);m->cx=m->bx;m->bx=ec_shl(m,m->bx);m->bx=ec_math(m,m->bx,m->cx,16,0,0);
    gw_bl(m,(uint8_t)ec_math(m,(uint8_t)m->bx,(uint8_t)m->ax,8,0,0));gw_bh(m,(uint8_t)ec_math(m,m->bx>>8,0,8,m->flags&1,0));
}
static void ev_relation(KiMachine16 *m,int peace,const KiEconomyHooks *h) {
    ec_cmp(m,(uint8_t)m->ax,24,8);if(m->flags&0x40) return;ec_cmp(m,m->ax>>8,24,8);if(m->flags&0x40) return;
    ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->si);ev_call(m,0x3697,peace?0x367a:0x364a,h);
    pn_cl(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x600)));pn_ch(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x600)));ec_cmp(m,(uint8_t)m->cx,m->cx>>8,8);if(!(m->flags&(1|0x40))) pn_cl(m,(uint8_t)(m->cx>>8));
    if(peace) {pn_cl(m,(uint8_t)m->cx|0x80);ec_logic(m,(uint8_t)m->cx,8);}
    else {pn_cl(m,(uint8_t)m->cx&0x7f);ec_logic(m,(uint8_t)m->cx,8);if(KI_EVENTS_MUTATION!=2) pn_cl(m,(uint8_t)st_shr(m,(uint8_t)m->cx,8));}
    ec_put(m,m->ds,(uint16_t)(m->si+0x600),(uint8_t)m->cx);ec_put(m,m->ds,(uint16_t)(m->bx+0x600),(uint8_t)m->cx);m->si=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);
}
static void ev_best_general(KiMachine16 *m) {
    ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);m->ax=ev_shifts(m,m->ax,2,0);gw_al(m,0);ec_logic(m,0,8);m->dx=0;ec_logic(m,0,16);m->bx=0x4240;m->cx=127;
    do {ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->bx+0x1c)),8);if(m->flags&0x40) {ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x17)),0,8);if(m->flags&0x40) {ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->bx+0x13)),8);if(m->flags&1) {gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x13)));m->dx=m->bx;}}}m->bx=ec_math(m,m->bx,32,16,0,0);--m->cx;}while(m->cx);
    ec_logic(m,(uint8_t)m->ax,8);int found=!(m->flags&0x40);if(found) m->ax=m->dx;m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->flags=(uint16_t)((m->flags&~1u)|!found);
}
static void ev_official_power(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->bx);ec_push(m,m->dx);gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x2a)));ec_cmp(m,m->bx>>8,0xff,8);
    if(m->flags&0x40) {m->ax=m->di;ev_call(m,0x37f5,0x3780,h);if(m->flags&1) goto done;m->bx=m->ax;}
    else {gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,1);m->bx=ec_math(m,m->bx,0x4240,16,0,0);}
    pn_dh(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x13)));gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->si+1)));gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,2,1);gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x2240)));m->bx=st_shr(m,m->bx,16);
    ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4257)),0,8);
    if(m->flags&0x40) {m->ax=m->si;ev_call(m,0x37f5,0x37b0,h);m->bx=m->ax;}
    else {ec_cmp(m,(uint8_t)m->ax,0x80,8);if(m->flags&1) goto done;m->bx=ec_math(m,m->bx,0x4240,16,0,0);}
    gw_dl(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x13)));ec_cmp(m,(uint8_t)m->dx,m->dx>>8,8);
    if(m->flags&1) goto weaker;
    if(m->flags&0x40) {ev_call(m,0xece0,0x37c8,h);ec_logic(m,(uint8_t)m->ax&1,8);if(!(m->flags&0x40)) goto weaker;}
    goto score;
weaker:gw_dl(m,16);gw_dl(m,(uint8_t)ec_math(m,(uint8_t)m->dx,m->dx>>8,8,0,1));
score:gw_al(m,(uint8_t)m->dx);gw_al(m,ev_shl8(m,(uint8_t)m->ax));m->flags&=(uint16_t)~1u;
done:m->dx=ec_pop(m);m->bx=ec_pop(m);
}
static void ev_residual_count(KiMachine16 *m) {
    ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);m->cx=ev_shifts(m,m->si,2,0);m->ax=ev_shifts(m,m->di,2,0);pn_dh(m,(uint8_t)(m->ax>>8));gw_dl(m,(uint8_t)(m->cx>>8));m->bx=0x4240;m->ax=0;ec_logic(m,0,16);pn_cl(m,127);
    do {ec_cmp(m,(uint8_t)m->dx,ec_byte(m,m->ds,(uint16_t)(m->bx+0x1d)),8);if(m->flags&0x40) {ec_cmp(m,m->dx>>8,ec_byte(m,m->ds,(uint16_t)(m->bx+0x1c)),8);if(m->flags&0x40) m->ax=ec_inc(m,m->ax,16);}m->bx=ec_math(m,m->bx,32,16,0,0);pn_cl(m,(uint8_t)st_dec(m,(uint8_t)m->cx,8));}while((uint8_t)m->cx);
    m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);ec_cmp(m,m->ax,KI_EVENTS_MUTATION==7?0:m->ax,16);m->flags=(uint16_t)((m->flags&~1u)|!(m->flags&0x40));
}
static void ev_residual_flags(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->dx);ec_push(m,m->ax);pn_dh(m,0);ec_logic(m,0,8);ev_call(m,0x3138,0x37df,h);if(m->flags&1) {pn_dh(m,(uint8_t)(m->dx>>8)|1);ec_logic(m,m->dx>>8,8);}pn_swap(&m->si,&m->di);ev_call(m,0x3138,0x37e9,h);if(m->flags&1) {pn_dh(m,(uint8_t)(m->dx>>8)|2);ec_logic(m,m->dx>>8,8);}pn_swap(&m->si,&m->di);m->ax=ec_pop(m);gw_ah(m,(uint8_t)(m->dx>>8));m->dx=ec_pop(m);
}
static void ev_diplomacy(KiMachine16 *m,int cooperation,const KiEconomyHooks *h) {
    ec_push(m,m->cx);pn_cl(m,1);ev_call(m,0x3771,cooperation?0x3718:0x36ca,h);if(m->flags&1) goto done;gw_dl(m,(uint8_t)m->ax);
    if(cooperation) {
        pn_swap(&m->di,&m->bx);ev_call(m,0x30cb,0x3721,h);pn_swap(&m->di,&m->bx);pn_dh(m,(uint8_t)m->ax);ev_call(m,0x30cb,0x3728,h);ec_cmp(m,(uint8_t)m->ax,m->dx>>8,8);if(m->flags&1) pn_cl(m,2);
        ec_cmp(m,(uint8_t)m->ax,0x80,8);if(m->flags&1) {gw_al(m,0);ec_logic(m,0,8);}gw_al(m,(uint8_t)m->ax&0x7f);ec_logic(m,(uint8_t)m->ax,8);gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x28)));gw_ah(m,ev_shl8(m,(uint8_t)(m->ax>>8)));gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,40,8,0,0));ec_cmp(m,(uint8_t)m->ax,m->ax>>8,8);if(m->flags&1) pn_cl(m,2);
        gw_ah(m,90);gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->ax,8,0,1));gw_dl(m,(uint8_t)(m->ax>>8));ec_cmp(m,(uint8_t)m->dx,0,8);if(ec_less(m)) {gw_dl(m,0);ec_logic(m,0,8);}ec_cmp(m,(uint8_t)m->dx,KI_EVENTS_MUTATION==8?59:60,8);if(!(m->flags&(1|0x40))) gw_dl(m,KI_EVENTS_MUTATION==8?59:60);
    } else {
        m->ax=ev_shifts(m,m->di,2,0);ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),8);if(m->flags&0x40) pn_cl(m,2);ev_call(m,0x30cb,0x36de,h);gw_al(m,(uint8_t)m->ax&0x7f);ec_logic(m,(uint8_t)m->ax,8);gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x28)));gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,2,8,0,0));gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,m->ax>>8,8,0,1));if(m->flags&1) {gw_al(m,0);ec_logic(m,0,8);}gw_ah(m,30);gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->ax,8,0,1));gw_dl(m,(uint8_t)ec_math(m,(uint8_t)m->dx,m->ax>>8,8,0,0));ec_cmp(m,(uint8_t)m->dx,0,8);if(ec_less(m)) {gw_dl(m,0);ec_logic(m,0,8);}
    }
    gw_dl(m,(uint8_t)st_shr(m,(uint8_t)m->dx,8));pn_dh(m,0);ec_logic(m,0,8);m->ax=1000;uint32_t product=(uint32_t)m->ax*m->dx;m->ax=(uint16_t)product;m->dx=(uint16_t)(product>>16);m->dx=m->ax;gw_al(m,(uint8_t)m->cx);ec_logic(m,m->dx,16);if(m->flags&0x40) {gw_al(m,0);ec_logic(m,0,8);}ev_call(m,0x37d8,cooperation?0x376e:0x370f,h);m->flags&=(uint16_t)~1u;
done:m->cx=ec_pop(m);
}
static void ev_neutral_target(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->di);ec_cmp(m,(uint8_t)m->dx,ec_byte(m,m->cs,0xcff),8);if(m->flags&0x40) goto done;ec_cmp(m,(uint8_t)m->dx,24,8);if(m->flags&0x40) goto done;
    gw_bh(m,(uint8_t)m->dx);gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,2,1);gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x19)));ec_cmp(m,m->ax>>8,KI_EVENTS_MUTATION==3?0x23:0x24,8);if(!(m->flags&1)) goto replace;
    gw_al(m,0);ec_logic(m,0,8);m->ax=ev_shifts(m,m->ax,2,1);m->di=m->ax;ec_push(m,m->bx);m->bx=m->si;ev_call(m,0x3091,0x35d6,h);m->cx=m->ax;m->bx=m->di;ev_call(m,0x3091,0x35dd,h);m->bx=ec_pop(m);ec_cmp(m,m->ax,m->cx,16);if(!(m->flags&1)) goto done;
replace:m->ax=ev_shifts(m,m->si,2,0);ec_put(m,m->ds,(uint16_t)(m->bx+0x19),(uint8_t)(m->ax>>8));
done:m->di=ec_pop(m);
}
static void ev_target(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->di);ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),24,8);if(m->flags&1) goto done;ec_cmp(m,(uint8_t)m->dx,24,8);if(m->flags&0x40) goto apply;
    gw_ah(m,(uint8_t)m->dx);gw_al(m,0);ec_logic(m,0,8);m->ax=ev_shifts(m,m->ax,2,1);m->di=m->ax;ev_call(m,0x30cb,0x3542,h);ec_cmp(m,(uint8_t)m->ax,0x80,8);if(m->flags&1) goto apply;
    m->cx=0x1a0;ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);if(m->flags&0x40) ev_call(m,0xcde,0x3553,h);
    else {m->cx=0x19f;ec_cmp(m,(uint8_t)m->dx,ec_byte(m,m->cs,0xcff),8);if(!(m->flags&0x40)) goto apply;ec_push(m,m->cx);ec_push(m,m->si);m->di=m->sp;ev_call(m,0xce7,0x3566,h);m->cx=0x3f;gw_al(m,0x93);ev_call(m,0x8810,0x356e,h);m->si=ec_pop(m);m->cx=ec_pop(m);}
    pn_dh(m,0xff);ec_push(m,m->dx);m->di=m->sp;gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->si+1)));gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,1);gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x425e)));gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4241)));ev_call(m,0x8810,0x358b,h);m->dx=ec_pop(m);
apply:ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);if(!(m->flags&0x40)) ec_put(m,m->ds,(uint16_t)(m->si+0x19),(uint8_t)m->dx);ev_call(m,0x35ab,0x3599,h);m->ax=ev_shifts(m,m->si,2,0);gw_al(m,(uint8_t)(m->ax>>8));gw_ah(m,(uint8_t)m->dx);ev_call(m,0x3639,0x35a6,h);
done:m->di=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void ev_unassign(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);gw_al(m,0xff);ec_put(m,m->ds,(uint16_t)(m->di+0x17),0);uint8_t old=ec_byte(m,m->ds,(uint16_t)(m->di+0x1d));ec_put(m,m->ds,(uint16_t)(m->di+0x1d),(uint8_t)m->ax);gw_al(m,old);gw_bh(m,(uint8_t)m->ax);gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,2,1);ec_cmp(m,ec_byte(m,m->ds,m->bx),0x80,8);if(m->flags&1) gw_al(m,0xff);ec_put(m,m->ds,(uint16_t)(m->di+0x1c),(uint8_t)m->ax);gw_ah(m,0xff);ev_call(m,0x2ad2,0x50fa,h);ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->cs,0xcff),8);
    if(m->flags&0x40) {ec_push(m,m->di);m->di=m->sp;ev_call(m,0xcde,0x5107,h);m->cx=0x25;gw_al(m,0x93);ev_call(m,0x8810,0x510f,h);m->di=ec_pop(m);m->cx=0x199;gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->di+0x1e)));gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->di+1)));ev_call(m,0x8810,0x511c,h);}
    m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void ev_cleanup(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->cx);ec_push(m,m->dx);ec_cmp(m,(uint8_t)m->ax,1,8);if(m->flags&0x40) {m->ax=m->dx;gw_dl(m,0);ec_logic(m,0,8);ev_call(m,0x5609,0x35fb,h);pn_swap(&m->si,&m->di);ev_call(m,0x563b,0x3600,h);pn_swap(&m->si,&m->di);}
    ec_push(m,m->di);m->ax=ev_shifts(m,m->si,2,0);pn_dh(m,(uint8_t)(m->ax>>8));m->ax=ev_shifts(m,m->di,2,0);gw_dl(m,(uint8_t)(m->ax>>8));m->ax=m->dx;m->ax=(uint16_t)((m->ax<<8)|(m->ax>>8));m->cx=127;m->di=0x4240;
    do {ec_cmp(m,m->dx,ec_word(m,m->ds,(uint16_t)(m->di+0x1c)),16);int match=m->flags&0x40;if(!match) {ec_cmp(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->di+0x1c)),16);match=m->flags&0x40;}if(match) ev_call(m,0x50d7,0x362a,h);m->di=ec_math(m,m->di,32,16,0,0);--m->cx;}while(m->cx);
    m->di=ec_pop(m);gw_al(m,4);ev_call(m,0x5e80,0x3635,h);m->dx=ec_pop(m);m->cx=ec_pop(m);m->ax=ec_pop(m);
}
static void ev_disaster_alloc(KiMachine16 *m) {
    ec_push(m,m->ax);ec_push(m,m->si);m->si=0x2040;
    for(;;) {ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);if(m->flags&1) break;m->si=ec_math(m,m->si,16,16,0,0);ec_cmp(m,m->si,KI_EVENTS_MUTATION==4?0x2130:0x2140,16);if(!(m->flags&1)) {m->flags|=1;goto done;}}
    ec_put(m,m->ds,m->si,0x80);ec_put(m,m->ds,(uint16_t)(m->si+0xe),(uint8_t)m->ax);ec_put(m,m->ds,(uint16_t)(m->si+0xd),KI_EVENTS_MUTATION==6?15:16);ec_store(m,m->ds,(uint16_t)(m->si+2),m->dx);ec_store(m,m->ds,(uint16_t)(m->si+4),m->bx);m->ax=1;
    for(unsigned i=0;i<4;++i) ec_put(m,m->ds,(uint16_t)(m->si+(unsigned[]){6,7,12,15}[i]),1);m->flags&=(uint16_t)~1u;
done:m->si=ec_pop(m);m->ax=ec_pop(m);
}
static void ev_disaster_clear(KiMachine16 *m) {
    ec_push(m,m->si);m->si=0x2040;do {ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);if(!(m->flags&1)) {ec_cmp(m,m->dx,ec_word(m,m->ds,(uint16_t)(m->si+2)),16);if(m->flags&0x40) {ec_cmp(m,m->bx,ec_word(m,m->ds,(uint16_t)(m->si+4)),16);if(m->flags&0x40) ec_put(m,m->ds,m->si,0);}}m->si=ec_math(m,m->si,16,16,0,0);ec_cmp(m,m->si,0x2140,16);}while(m->flags&1);m->si=ec_pop(m);
}
static void ev_corps_capital(KiMachine16 *m) {
    ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);m->si=0x2240;m->cx=0;ec_logic(m,0,16);m->dx=m->cx;pn_cl(m,(uint8_t)m->ax);gw_dl(m,(uint8_t)(m->ax>>8));m->cx=ev_shifts(m,m->cx,3,0);m->dx=ev_shifts(m,m->dx,3,0);gw_bh(m,127);
    do {ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+1)),(uint8_t)m->bx,8);if(m->flags&0x40) {ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);if(!(m->flags&1)) {ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->si+0x20)),8);if((m->flags&0x40)&&KI_EVENTS_MUTATION!=5) {ec_put(m,m->ds,(uint16_t)(m->si+0x20),(uint8_t)m->ax);ec_cmp(m,ec_word(m,m->ds,(uint16_t)(m->si+0x14)),m->cx,16);if(m->flags&0x40) {ec_store(m,m->ds,(uint16_t)(m->si+0x14),m->dx);uint8_t v=ec_byte(m,m->ds,m->si)|2;ec_logic(m,v,8);ec_put(m,m->ds,m->si,v);}}}}m->si=ec_math(m,m->si,64,16,0,0);gw_bh(m,(uint8_t)st_dec(m,m->bx>>8,8));}while(m->bx>>8);
    m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);
}
static void ev_best_city(KiMachine16 *m) {
    ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);m->cx=192;m->si=0x840;m->dx=0;ec_logic(m,0,16);m->bx=0xff;m->di=0xffff;
    do {ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->si+1)),8);if(!(m->flags&0x40)) goto next;gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x16)));gw_ah(m,(m->ax>>8)&15);ec_logic(m,m->ax>>8,8);ec_cmp(m,(uint8_t)m->bx,m->ax>>8,8);if(m->flags&1) goto next;ec_cmp(m,m->dx,ec_word(m,m->ds,(uint16_t)(m->si+0xe)),16);if(!(m->flags&(1|0x40))) goto next;
        ec_logic(m,m->bx>>8,8);if(m->flags&0x40) {gw_bl(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x16)));gw_bl(m,(uint8_t)m->bx&15);ec_logic(m,(uint8_t)m->bx,8);m->dx=ec_word(m,m->ds,(uint16_t)(m->si+0xe));m->di=m->si;}
        ec_logic(m,ec_byte(m,m->ds,m->si)&0x1f,8);if(!(m->flags&0x40)) goto next;gw_bh(m,1);gw_bl(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x16)));gw_bl(m,(uint8_t)m->bx&15);ec_logic(m,(uint8_t)m->bx,8);m->dx=ec_word(m,m->ds,(uint16_t)(m->si+0xe));m->di=m->si;
next:m->si=ec_math(m,m->si,32,16,0,0);--m->cx;}while(m->cx);
    m->ax=m->di;ec_cmp(m,m->ax,0xffff,16);m->flags=(uint16_t)((m->flags&~1u)|((m->flags&0x40)?1:0));m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);
}
static void ev_capital_tail(KiMachine16 *m,const KiEconomyHooks *h) {
    m->ax=ec_math(m,m->ax,0x840,16,0,1);m->ax=ev_shifts(m,m->ax,3,0);gw_al(m,(uint8_t)(m->ax>>8));uint8_t old=ec_byte(m,m->ds,(uint16_t)(m->si+3));ec_put(m,m->ds,(uint16_t)(m->si+3),(uint8_t)(m->ax>>8));gw_ah(m,old);ec_cmp(m,(uint8_t)m->ax,m->ax>>8,8);if(m->flags&0x40) return;
    m->bx=ev_shifts(m,m->si,2,0);gw_bl(m,(uint8_t)(m->bx>>8));ev_call(m,0x4502,0x341a,h);ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);
    if(m->flags&0x40) {gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->si+3)));gw_ah(m,0xff);ec_push(m,m->ax);m->di=m->sp;gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->si+1)));gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,1);gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x425e)));gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4241)));m->cx=0x1a4;ev_call(m,0x8810,0x3442,h);m->sp=ec_math(m,m->sp,2,16,0,0);ev_call(m,0x5e60,0x3448,h);return;}
    pn_ch(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x2a)));ec_cmp(m,m->cx>>8,0xff,8);if(m->flags&0x40) return;pn_cl(m,0);ec_logic(m,0,8);m->cx=ev_shifts(m,m->cx,3,1);m->cx=ec_math(m,m->cx,0x4240,16,0,0);gw_ah(m,0xff);ec_push(m,m->ax);ec_push(m,m->cx);ec_push(m,m->si);m->di=m->sp;ev_call(m,0xcde,0x3467,h);m->cx=0x39;gw_al(m,0x93);ev_call(m,0x8810,0x346f,h);m->si=ec_pop(m);m->di=ec_pop(m);ec_push(m,m->di);ec_push(m,m->si);m->cx=0x1a4;gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->di+0x1e)));gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->di+1)));m->di=m->sp;ev_call(m,0x8810,0x3481,h);m->si=ec_pop(m);m->cx=ec_pop(m);m->ax=ec_pop(m);
}
static void ev_event1(KiMachine16 *m,const KiEconomyHooks *h) {
    ev_call(m,0x351a,0x320f,h);if(m->flags&1) return;m->bx=m->si;gw_ah(m,(uint8_t)m->dx);ev_call(m,0x351a,0x3218,h);if(m->flags&1) return;pn_swap(&m->bx,&m->si);ev_call(m,0x3526,0x321f,h);
}
static void ev_event2(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->dx);ev_call(m,0x351a,0x3224,h);if(m->flags&1) goto done;m->bx=m->si;gw_ah(m,(uint8_t)m->dx);ev_call(m,0x351a,0x322d,h);if(m->flags&1) goto done;pn_swap(&m->bx,&m->si);m->di=m->si;gw_ah(m,(uint8_t)(m->dx>>8));ev_call(m,0x351a,0x3238,h);if(m->flags&1) goto done;pn_swap(&m->si,&m->di);ev_call(m,0x3712,0x323f,h);if(m->flags&1) goto done;ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);
    if(m->flags&0x40) {ev_call(m,0x38e6,0x324b,h);ev_call(m,0xcde,0x324e,h);m->cx=0x2f;ev_call(m,0x3c3d,0x3254,h);}ec_cmp(m,(uint8_t)m->ax,2,8);if(!(m->flags&1)) goto done;ev_call(m,0x35ed,0x325b,h);m->dx=ec_pop(m);ec_push(m,m->dx);ev_call(m,0x3526,0x3260,h);
done:m->dx=ec_pop(m);
}
static void ev_event3(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->bp);gw_al(m,(uint8_t)m->dx);m->bp=m->ax;gw_bh(m,(uint8_t)(m->ax>>8));gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,2,1);m->di=m->bx;gw_bh(m,(uint8_t)m->dx);gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,2,1);m->si=m->bx;ev_call(m,0x36c4,0x327e,h);if(m->flags&1) goto done;ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);
    if(m->flags&0x40) {ev_call(m,0x38c7,0x328a,h);ev_call(m,0xcde,0x328d,h);m->cx=0x2b;ev_call(m,0x3c3d,0x3293,h);}ec_cmp(m,(uint8_t)m->ax,2,8);if(!(m->flags&1)) goto done;ev_call(m,0x35ed,0x329a,h);m->ax=m->bp;m->ax=(uint16_t)((m->ax<<8)|(m->ax>>8));ev_call(m,0x45f8,0x32a1,h);ev_call(m,0x4236,0x32a4,h);ev_call(m,0x3669,0x32a7,h);
done:m->bp=ec_pop(m);
}
static void ev_report(KiMachine16 *m,int diplomatic,const KiEconomyHooks *h) {
    gw_al(m,0);ec_logic(m,0,8);m->ax=ev_shifts(m,m->ax,diplomatic?2:3,1);m->ax=ec_math(m,m->ax,diplomatic?0:0x840,16,0,0);m->bx=m->ax;gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->bx+(diplomatic?0x2a:0x19))));ec_cmp(m,m->ax>>8,0xff,8);if(m->flags&0x40) return;gw_al(m,0);ec_logic(m,0,8);m->ax=ev_shifts(m,m->ax,3,1);m->ax=ec_math(m,m->ax,0x4240,16,0,0);ec_push(m,m->ax);ec_push(m,m->bx);m->di=m->sp;ev_call(m,0xcde,diplomatic?0x330e:0x32d0,h);m->cx=diplomatic?0x39:0x38;gw_al(m,0x93);ev_call(m,0x8810,diplomatic?0x3316:0x32d8,h);m->bx=ec_pop(m);m->si=ec_pop(m);m->ax=m->dx;m->cx=diplomatic?0x13f:0x116;ev_call(m,0x2078,diplomatic?0x3320:0x32e2,h);ev_call(m,0x39e8,diplomatic?0x3323:0x32e5,h);ev_call(m,0x20d6,diplomatic?0x3326:0x32e8,h);
}
static void ev_diplomat_event(KiMachine16 *m,int cooperation,const KiEconomyHooks *h) {
    ev_call(m,0x351a,cooperation?0x338b:0x332a,h);if(m->flags&1) return;
    if(cooperation) {m->bx=m->si;gw_ah(m,(uint8_t)m->dx);ev_call(m,0x351a,0x3394,h);if(m->flags&1) return;pn_swap(&m->bx,&m->si);}
    gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x2a)));ec_cmp(m,(uint8_t)m->ax,0xff,8);if(m->flags&0x40) return;gw_ah(m,0xff);ec_push(m,m->ax);ec_push(m,m->si);m->di=m->sp;ev_call(m,0xcde,cooperation?0x33a8:0x333c,h);m->cx=0x39;gw_al(m,0x93);ev_call(m,0x8810,cooperation?0x33b0:0x3344,h);m->si=ec_pop(m);m->ax=ec_pop(m);m->di=ec_word(m,m->cs,0xcfd);ev_call(m,cooperation?0x3712:0x36c4,cooperation?0x33ba:0x334e,h);
    if(m->flags&1) {gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x2a)));gw_bl(m,0);ec_logic(m,0,8);m->bx=ev_shifts(m,m->bx,3,1);gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4241)));m->cx=0x3a;ev_call(m,0x8810,cooperation?0x33d1:0x3365,h);return;}
    m->cx=cooperation?0x2f:0x2b;ev_call(m,0x3c3d,cooperation?0x33d9:0x336d,h);ec_cmp(m,(uint8_t)m->ax,2,8);if(!(m->flags&1)) return;ev_call(m,0x35ed,cooperation?0x33e0:0x3374,h);
    if(cooperation) {m->bx=ev_shifts(m,m->bx,2,0);gw_dl(m,(uint8_t)(m->bx>>8));ev_call(m,0x3526,0x33e9,h);}
    else {m->ax=ev_shifts(m,m->si,2,0);gw_al(m,ec_byte(m,m->cs,0xcff));ev_call(m,0x45f8,0x3381,h);ev_call(m,0x4236,0x3384,h);ev_call(m,0x3669,0x3387,h);}
}
static void ev_capital_event(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_cmp(m,m->ax>>8,ec_byte(m,m->cs,0xcff),8);if(m->flags&0x40) return;gw_bh(m,(uint8_t)(m->ax>>8));ev_call(m,0x351a,0x33f6,h);if(m->flags&1) return;gw_al(m,(uint8_t)(m->bx>>8));ev_call(m,0x6a3d,0x33fd,h);
    /* Fall-through into the separately indexed original function, without CALL/RET. */
    m->ip=0x33fd;if(h&&h->enter) h->enter(m,0x33fd,h->user);ev_capital_tail(m,h);
}
static void ev_event12(KiMachine16 *m,const KiEconomyHooks *h) {
    m->si=m->dx;m->dx=ec_word(m,m->ds,(uint16_t)(m->si+8));m->bx=ec_word(m,m->ds,(uint16_t)(m->si+10));ec_logic(m,m->ax>>8,8);if(m->flags&0x40) {ec_put(m,m->ds,(uint16_t)(m->si+0x15),(uint8_t)(m->ax>>8));ev_call(m,0x2438,0x34c3,h);return;}gw_al(m,(uint8_t)(m->ax>>8));ev_call(m,0x23ff,0x34c9,h);if(m->flags&1) return;pn_cl(m,ec_byte(m,m->cs,0xcff));ec_cmp(m,(uint8_t)m->cx,ec_byte(m,m->ds,(uint16_t)(m->si+1)),8);
    if(m->flags&0x40) {ec_push(m,m->ax);ec_push(m,m->si);m->di=m->sp;pn_cl(m,(uint8_t)(m->ax>>8));pn_ch(m,0);ec_logic(m,0,8);ev_call(m,0xce7,0x34e1,h);m->cx=ec_math(m,m->cx,0x46,16,0,0);gw_al(m,0x93);ev_call(m,0x8810,0x34e9,h);m->si=ec_pop(m);m->ax=ec_pop(m);}
    ev_call(m,0xece0,0x34ee,h);gw_al(m,(uint8_t)m->ax&7);ec_logic(m,(uint8_t)m->ax,8);gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,4,8,0,0));ec_put(m,m->ds,(uint16_t)(m->si+0x15),(uint8_t)m->ax);ev_call(m,0xece0,0x34f8,h);gw_al(m,(uint8_t)m->ax&7);ec_logic(m,(uint8_t)m->ax,8);gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,6,8,0,0));gw_bl(m,(uint8_t)m->ax);m->ax=12;m->dx=m->si;ev_call(m,0x301c,0x3506,h);
}
int ki_events_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {
    case 0x320c:ev_event1(m,h);break;case 0x3220:ev_event2(m,h);break;case 0x3262:ev_event3(m,h);break;
    case 0x32a9:ev_report(m,0,h);break;case 0x32e9:ev_report(m,1,h);break;case 0x3327:ev_diplomat_event(m,0,h);break;case 0x3388:ev_diplomat_event(m,1,h);break;
    case 0x33ea:ev_capital_event(m,h);break;case 0x33fd:ev_capital_tail(m,h);break;
    case 0x3485:gw_al(m,0);ec_logic(m,0,8);m->ax=ev_shifts(m,m->ax,3,1);m->ax=ec_math(m,m->ax,0x4240,16,0,0);m->di=m->ax;ev_call(m,0x50d7,0x3495,h);break;
    case 0x34a6:ev_call(m,0xece0,0x34a9,h);gw_al(m,(uint8_t)m->ax&15);ec_logic(m,(uint8_t)m->ax,8);gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,24,8,0,0));ev_call(m,0x237e,0x34b0,h);break;
    case 0x34b1:ev_event12(m,h);break;case 0x351a:ev_presence(m);break;case 0x3526:ev_target(m,h);break;case 0x35ab:ev_neutral_target(m,h);break;case 0x35ed:ev_cleanup(m,h);break;
    case 0x3639:ev_relation(m,0,h);break;case 0x3669:ev_relation(m,1,h);break;case 0x3697:ev_pair(m);break;case 0x36c4:ev_diplomacy(m,0,h);break;case 0x3712:ev_diplomacy(m,1,h);break;
    case 0x3771:ev_official_power(m,h);break;case 0x37d8:ev_residual_flags(m,h);break;case 0x37f5:ev_best_general(m);break;case 0x3138:ev_residual_count(m);break;
    case 0x4502:ev_corps_capital(m);break;case 0x6a3d:ev_best_city(m);break;case 0x23ff:ev_disaster_alloc(m);break;case 0x2438:ev_disaster_clear(m);break;case 0x50d7:ev_unassign(m,h);break;
    default:return 0;
    }return 1;
}
#define EV_WRAP(name,target) void name(KiMachine16 *m,const KiEconomyHooks *h) {ki_events_body(m,target,h);m->ip=ec_pop(m);}
EV_WRAP(sub_1320C,0x320c) EV_WRAP(sub_13220,0x3220) EV_WRAP(sub_13262,0x3262)
EV_WRAP(sub_132A9,0x32a9) EV_WRAP(sub_132E9,0x32e9) EV_WRAP(sub_13327,0x3327)
EV_WRAP(sub_13388,0x3388) EV_WRAP(sub_133EA,0x33ea) EV_WRAP(sub_133FD,0x33fd)
EV_WRAP(sub_13485,0x3485) EV_WRAP(sub_134A6,0x34a6) EV_WRAP(sub_134B1,0x34b1)
EV_WRAP(sub_1351A,0x351a) EV_WRAP(sub_13526,0x3526) EV_WRAP(sub_135AB,0x35ab)
EV_WRAP(sub_135ED,0x35ed) EV_WRAP(sub_13639,0x3639) EV_WRAP(sub_13669,0x3669)
EV_WRAP(sub_13697,0x3697) EV_WRAP(sub_136C4,0x36c4) EV_WRAP(sub_13712,0x3712)
EV_WRAP(sub_13771,0x3771) EV_WRAP(sub_137D8,0x37d8) EV_WRAP(sub_137F5,0x37f5)
EV_WRAP(sub_13138,0x3138) EV_WRAP(sub_14502,0x4502) EV_WRAP(sub_16A3D,0x6a3d)
EV_WRAP(sub_123FF,0x23ff) EV_WRAP(sub_12438,0x2438) EV_WRAP(sub_150D7,0x50d7)
#undef EV_WRAP
void ki_events_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    switch(t) {
    case 0x320c:case 0x3220:case 0x3262:case 0x32a9:case 0x32e9:case 0x3327:case 0x3388:case 0x33ea:case 0x33fd:case 0x3485:case 0x34a6:case 0x34b1:
    case 0x351a:case 0x3526:case 0x35ab:case 0x35ed:case 0x3639:case 0x3669:case 0x3697:case 0x36c4:case 0x3712:case 0x3771:case 0x37d8:case 0x37f5:case 0x3138:case 0x4502:case 0x6a3d:case 0x23ff:case 0x2438:case 0x50d7:
        if(h&&h->enter) h->enter(m,t,h->user);ki_events_body(m,t,h);m->ip=ec_pop(m);break;
    default:ki_hourly_invoke(m,t,h);break;
    }
}
