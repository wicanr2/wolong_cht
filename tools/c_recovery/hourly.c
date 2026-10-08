/* Original DOS/V hourly/dispatch core, spec/218; include after politics.c. */
#include "hourly.h"
#ifndef KI_HOURLY_MUTATION
#define KI_HOURLY_MUTATION 0
#endif
static void hr_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *h) {ec_push(m,next);m->ip=target;ki_hourly_invoke(m,target,h);}
static uint16_t hr_rcr(KiMachine16 *m,uint16_t a) {
    uint16_t r=(uint16_t)((a>>1)|((m->flags&1)<<15));m->flags=(uint16_t)((m->flags&~0x801u)|(a&1)|((((r>>15)^(r>>14))&1)<<11));return r;
}
static void hr_expense(KiMachine16 *m) {
    ec_push(m,m->ax);ec_push(m,m->dx);m->ax=ec_math(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+0x1a)),16,0,0);
    gw_dl(m,(uint8_t)ec_math(m,(uint8_t)m->dx,ec_byte(m,m->ds,(uint16_t)(m->si+0x1c)),8,m->flags&1,0));ec_cmp(m,(uint8_t)m->dx,9,8);
    if(!ec_less(m)) {int above=ec_greater(m);gw_dl(m,9);if(!above) {ec_cmp(m,m->ax,0xfe98,16);above=!(m->flags&1);}if(above) m->ax=KI_HOURLY_MUTATION==1?0xfe97:0xfe98;}
    ec_store(m,m->ds,(uint16_t)(m->si+0x1a),m->ax);ec_put(m,m->ds,(uint16_t)(m->si+0x1c),(uint8_t)m->dx);m->dx=ec_pop(m);m->ax=ec_pop(m);
}
static void hr_maintenance(KiMachine16 *m,const KiEconomyHooks *h) {
    gw_dl(m,0);ec_logic(m,0,8);m->ax=ec_word(m,m->ds,(uint16_t)(m->si+4));m->ax=ec_math(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+6)),16,0,0);
    gw_dl(m,(uint8_t)ec_math(m,(uint8_t)m->dx,0,8,KI_HOURLY_MUTATION==2?0:m->flags&1,0));m->ax=ec_math(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+8)),16,0,0);gw_dl(m,(uint8_t)ec_math(m,(uint8_t)m->dx,0,8,m->flags&1,0));
    for(unsigned i=0;i<5;++i) {gw_dl(m,(uint8_t)st_shr(m,(uint8_t)m->dx,8));m->ax=hr_rcr(m,m->ax);}hr_call(m,0x5673,0x3e8d,h);
}
static void hr_diplomat(KiMachine16 *m,const KiEconomyHooks *h) {
    gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x2a)));ec_cmp(m,m->bx>>8,0xff,8);if(m->flags&0x40) return;
    hr_call(m,0xece0,0x3e99,h);ec_cmp(m,(uint8_t)m->ax,32,8);if(!(m->flags&1)) return;
    gw_bl(m,0);ec_logic(m,0,8);for(unsigned i=0;i<3;++i) m->bx=st_shr(m,m->bx,16);ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x425a)),0,8);if(m->flags&0x40) return;
    gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4253)));gw_ah(m,KI_HOURLY_MUTATION==7?24:23);gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->ax,8,0,1));uint16_t budget=(uint16_t)(m->bx+0x425a);
    ec_put(m,m->ds,budget,(uint8_t)ec_math(m,ec_byte(m,m->ds,budget),m->ax>>8,8,0,1));if(m->flags&1) ec_put(m,m->ds,budget,0);
    hr_call(m,0xece0,0x3ec2,h);gw_al(m,(uint8_t)m->ax&0xf);ec_logic(m,(uint8_t)m->ax,8);ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4253)),8);if(!(m->flags&(1|0x40))) return;
    m->di=ec_word(m,m->cs,0xcfd);hr_call(m,0x310a,0x3ed2,h);m->di=m->dx;gw_al(m,ec_byte(m,m->ds,m->bx));gw_dl(m,ec_byte(m,m->ds,m->di));gw_ah(m,(uint8_t)m->ax);m->ax&=0x807f;ec_logic(m,m->ax,16);ec_cmp(m,(uint8_t)m->ax,100,8);
    if(m->flags&1) gw_al(m,(uint8_t)ec_inc(m,(uint8_t)m->ax,8));gw_al(m,(uint8_t)m->ax|(m->ax>>8));ec_logic(m,(uint8_t)m->ax,8);ec_put(m,m->ds,m->bx,(uint8_t)m->ax);ec_cmp(m,(uint8_t)m->ax,(uint8_t)m->dx,8);
    if(!(m->flags&(1|0x40))) {pn_dh(m,(uint8_t)m->dx);m->dx&=0x807f;ec_logic(m,m->dx,16);ec_cmp(m,(uint8_t)m->dx,100,8);if(m->flags&1) gw_dl(m,(uint8_t)ec_inc(m,(uint8_t)m->dx,8));gw_dl(m,(uint8_t)m->dx|(m->dx>>8));ec_logic(m,(uint8_t)m->dx,8);ec_put(m,m->ds,m->di,(uint8_t)m->dx);}
}
static void hr_dispatch(KiMachine16 *m,const KiEconomyHooks *h) {
    uint8_t delay=(uint8_t)st_dec(m,ec_byte(m,m->ds,0x31ad),8);ec_put(m,m->ds,0x31ad,delay);if(!(m->flags&0x40)) return;
    ec_cmp(m,ec_word(m,m->ds,0xd20),0x100,16);if(!(m->flags&1)) return;ec_put(m,m->ds,0x31ad,KI_HOURLY_MUTATION==3?9:10);m->es=ec_word(m,m->ds,0xd56);m->bx=ec_word(m,m->ds,0xd20);m->ax=ec_word(m,m->es,m->bx);m->dx=ec_word(m,m->es,(uint16_t)(m->bx+2));m->bx=ec_math(m,m->bx,KI_HOURLY_MUTATION==4?2:4,16,0,0);ec_store(m,m->ds,0xd20,m->bx);ec_logic(m,(uint8_t)m->ax,8);if(m->flags&0x40) return;
    m->ds=ec_word(m,m->ds,0xd52);gw_bl(m,(uint8_t)m->ax);gw_bh(m,0);ec_logic(m,0,8);m->bx=st_dec(m,m->bx,16);m->bx=ec_shl(m,m->bx);uint16_t handler=ec_word(m,m->cs,(uint16_t)(0x31f2+m->bx));hr_call(m,handler,0x31ed,h);m->ax=m->cs;m->ds=m->ax;
}
static void hr_message(KiMachine16 *m,const KiEconomyHooks *h) {
    gw_al(m,(uint8_t)(m->ax>>8));gw_ah(m,0xff);m->cx=m->dx;ec_push(m,m->ax);m->di=m->sp;gw_al(m,0x93);hr_call(m,0x8810,0x34a4,h);m->ax=ec_pop(m);
}
static void hr_trust(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);m->ds=ec_word(m,m->cs,0xd52);uint8_t value=(uint8_t)ec_math(m,ec_byte(m,m->cs,0xd00),(uint8_t)m->ax,8,0,1);ec_put(m,m->cs,0xd00,value);
    if(m->flags&1) {if(KI_HOURLY_MUTATION!=5) ec_put(m,m->cs,0xd00,0);m->cx=0x19e;}
    gw_al(m,2);hr_call(m,0x2f5,0x3de7,h);ec_cmp(m,m->cx,0xffff,16);
    if(!(m->flags&0x40)) {hr_call(m,0x87ff,0x3def,h);gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x425e)));gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x4241)));hr_call(m,0x8810,0x3dfa,h);ec_cmp(m,ec_byte(m,m->cs,0xd00),0,8);if(m->flags&0x40) {gw_al(m,1);hr_call(m,0x1cb1,0x3e07,h);}}
    gw_al(m,2);hr_call(m,0x5e80,0x3e0c,h);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
static void hr_event13(KiMachine16 *m,const KiEconomyHooks *h) {
    hr_call(m,0xce7,0x350a,h);gw_al(m,0x93);m->cx=0x33;hr_call(m,0x8810,0x3512,h);gw_al(m,50);m->cx=m->dx;hr_call(m,0x3dc9,0x3519,h);
}
static void hr_world(KiMachine16 *m,const KiEconomyHooks *h) {
    hr_call(m,0x31ae,0x3e14,h);m->si=ec_word(m,m->ds,0xd1c);m->ds=ec_word(m,m->cs,0xd52);ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);
    if(!(m->flags&1)) {uint8_t flags=ec_byte(m,m->ds,m->si)&0xbf;ec_logic(m,flags,8);ec_put(m,m->ds,m->si,flags);m->ax=ec_byte(m,m->ds,(uint16_t)(m->si+0x23));ec_logic(m,0,8);for(unsigned i=0;i<3;++i) m->ax=ec_shl(m,m->ax);m->ax=ec_math(m,m->ax,24,16,0,0);ec_cmp(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+0x21)),16);
        if(!ec_less(m)) {ec_put(m,m->ds,(uint16_t)(m->si+0x19),0xff);m->ax=st_shr(m,m->ax,16);ec_cmp(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+0x21)),16);if(!ec_less(m)) {flags=ec_byte(m,m->ds,m->si)|0x40;ec_logic(m,flags,8);ec_put(m,m->ds,m->si,flags);}}}
    if(KI_HOURLY_MUTATION==6) {hr_call(m,0x3e8e,0x3e4c,h);hr_call(m,0x3e65,0x3e49,h);}else {hr_call(m,0x3e65,0x3e49,h);hr_call(m,0x3e8e,0x3e4c,h);}
    m->si=ec_math(m,m->si,64,16,0,0);ec_cmp(m,m->si,0x580,16);if(!(m->flags&1)) {m->si=0;ec_logic(m,0,16);}m->ax=m->cs;m->ds=m->ax;ec_store(m,m->ds,0xd1c,m->si);gw_al(m,4);hr_call(m,0x5e80,0x3e64,h);
}
int ki_hourly_body(KiMachine16 *m,uint16_t target,const KiEconomyHooks *h) {
    switch(target) {case 0x3e11:hr_world(m,h);break;case 0x3e65:hr_maintenance(m,h);break;case 0x3e8e:hr_diplomat(m,h);break;case 0x5673:hr_expense(m);break;case 0x31ae:hr_dispatch(m,h);break;case 0x3496:hr_message(m,h);break;case 0x3507:hr_event13(m,h);break;case 0x3dc9:hr_trust(m,h);break;default:return 0;}return 1;
}
#define HR_WRAP(name,target) void name(KiMachine16 *m,const KiEconomyHooks *h) {ki_hourly_body(m,target,h);m->ip=ec_pop(m);}
HR_WRAP(sub_13E11,0x3e11) HR_WRAP(sub_13E65,0x3e65) HR_WRAP(sub_13E8E,0x3e8e) HR_WRAP(sub_15673,0x5673)
HR_WRAP(sub_131AE,0x31ae) HR_WRAP(sub_13496,0x3496) HR_WRAP(sub_13507,0x3507) HR_WRAP(sub_13DC9,0x3dc9)
#undef HR_WRAP
void ki_hourly_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *h) {
    switch(target) {case 0x3e11:case 0x3e65:case 0x3e8e:case 0x5673:case 0x31ae:case 0x3496:case 0x3507:case 0x3dc9:if(h&&h->enter) h->enter(m,target,h->user);ki_hourly_body(m,target,h);m->ip=ec_pop(m);break;default:ki_politics_invoke(m,target,h);break;}
}
