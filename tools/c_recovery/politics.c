/* Monthly politics/prisoners, original DOS/V entries; spec/211. Include after world_update.c. */
#include "politics.h"
#ifndef KI_POLITICS_MUTATION
#define KI_POLITICS_MUTATION 0
#endif
static void pn_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=target;ki_politics_invoke(m,target,h);
}
static void pn_ch(KiMachine16 *m,uint8_t v) {m->cx=(uint16_t)((m->cx&0xffu)|((uint16_t)v<<8));}
static void pn_cl(KiMachine16 *m,uint8_t v) {m->cx=(uint16_t)((m->cx&0xff00u)|v);}
static void pn_dh(KiMachine16 *m,uint8_t v) {m->dx=(uint16_t)((m->dx&0xffu)|((uint16_t)v<<8));}
static void pn_swap(uint16_t *a,uint16_t *b) {uint16_t v=*a;*a=*b;*b=v;}
static void pn_assign(KiMachine16 *m) {
    ec_push(m,m->bx);ec_cmp(m,m->ax>>8,0xff,8);
    if(!(m->flags&0x40)) {m->bx=(uint16_t)((m->ax>>8)<<8);ec_logic(m,0,8);m->bx=st_shr(m,m->bx,16);m->bx=st_shr(m,m->bx,16);ec_put(m,m->ds,(uint16_t)(m->bx+0x18),(uint8_t)st_dec(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x18)),8));}
    ec_cmp(m,(uint8_t)m->ax,0xff,8);
    if(!(m->flags&0x40)) {m->bx=(uint16_t)((uint16_t)(uint8_t)m->ax<<8);ec_logic(m,0,8);m->bx=st_shr(m,m->bx,16);m->bx=st_shr(m,m->bx,16);uint8_t v=(uint8_t)ec_inc(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x18)),8);if(KI_POLITICS_MUTATION==4) ++v;ec_put(m,m->ds,(uint16_t)(m->bx+0x18),v);}
    m->bx=ec_pop(m);
}
static void pn_notice(KiMachine16 *m,const KiEconomyHooks *h) {
    gw_al(m,ec_byte(m,m->cs,0xcff));ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->si+0x1c)),8);
    if(m->flags&0x40) {ec_push(m,m->si);m->di=m->sp;pn_call(m,0xcde,0x599f,h);gw_al(m,0x93);pn_call(m,0x8810,0x59a4,h);m->si=ec_pop(m);}
}
static void pn_writer(KiMachine16 *m) {
    ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);m->cx=m->ax;gw_bh(m,0);ec_logic(m,0,8);
    m->bx=ec_shl(m,m->bx);m->bx=ec_shl(m,m->bx);m->bx=ec_math(m,m->bx,ec_word(m,m->cs,0xd20),16,0,0);m->ds=ec_word(m,m->cs,0xd56);
    for(;;) {
        uint16_t v=KI_POLITICS_MUTATION==1?ec_word(m,m->ds,m->bx):ec_byte(m,m->ds,m->bx);ec_cmp(m,v,0,KI_POLITICS_MUTATION==1?16:8);
        if(m->flags&0x40) {ec_store(m,m->ds,m->bx,m->cx);ec_store(m,m->ds,(uint16_t)(m->bx+2),m->dx);break;}
        m->bx=ec_math(m,m->bx,4,16,0,0);ec_cmp(m,m->bx,0x400,16);if(!(m->flags&1)) break;
    }
    m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
static void pn_event(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->cx);m->cx=m->si;m->cx=ec_shl(m,m->cx);m->cx=ec_shl(m,m->cx);gw_ah(m,(uint8_t)(m->cx>>8));pn_call(m,0x2fbf,0x2fbd,h);m->cx=ec_pop(m);
}
static void pn_join(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->di);gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)));ec_cmp(m,m->bx>>8,0xff,8);int enlist=0;
    if(!(m->flags&0x40)) {
        pn_call(m,0xece0,0x58a5,h);ec_cmp(m,(uint8_t)m->ax,64,8);
        if(m->flags&1) {
            ec_put(m,m->ds,(uint16_t)(m->si+0x19),0xff);gw_bl(m,0);ec_logic(m,0,8);m->bx=st_shr(m,m->bx,16);m->bx=st_shr(m,m->bx,16);ec_cmp(m,ec_byte(m,m->ds,m->bx),0x80,8);
            if(!(m->flags&1)) enlist=1;
            else {ec_logic(m,ec_byte(m,m->ds,m->si)&0x20,8);if(!(m->flags&0x40)) ec_put(m,m->ds,m->si,0);}
        }
    } else {
        m->di=0;m->bx=0;m->ax=0xff16;
        do {
            ec_cmp(m,ec_byte(m,m->ds,m->di),0x80,8);
            if(!(m->flags&1)) {ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->di+0x18)),8);if(!(m->flags&1)) {gw_ah(m,ec_byte(m,m->ds,(uint16_t)(m->di+0x18)));m->bx=m->di;}}
            m->di=ec_math(m,m->di,64,16,0,0);gw_al(m,(uint8_t)st_dec(m,(uint8_t)m->ax,8));
        } while((uint8_t)m->ax);
        pn_call(m,0xece0,0x58e3,h);gw_al(m,(uint8_t)m->ax&0x3f);ec_logic(m,(uint8_t)m->ax,8);gw_al(m,(uint8_t)ec_inc(m,(uint8_t)m->ax,8));ec_cmp(m,(uint8_t)m->ax,48,8);
        if(m->flags&1) {
            ec_cmp(m,(uint8_t)m->ax,24,8);if(!(m->flags&1)) gw_al(m,1);
            for(;;) {ec_cmp(m,ec_byte(m,m->ds,m->bx),0x80,8);if(!(m->flags&1)) {gw_al(m,(uint8_t)st_dec(m,(uint8_t)m->ax,8));if(m->flags&0x40) {enlist=1;break;}}
                m->bx=ec_math(m,m->bx,64,16,0,0);ec_cmp(m,m->bx,0x580,16);if(!(m->flags&1)) {m->bx=0;ec_logic(m,0,16);}}
        } else {m->bx=ec_word(m,m->cs,0xcfd);gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x23)));gw_al(m,(uint8_t)st_shr(m,(uint8_t)m->ax,8));gw_al(m,(uint8_t)st_shr(m,(uint8_t)m->ax,8));gw_al(m,(uint8_t)ec_inc(m,(uint8_t)m->ax,8));ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->bx+0x18)),8);enlist=!(m->flags&(1|0x40));}
    }
    if(enlist) {
        ec_cmp(m,m->bx,ec_word(m,m->cs,0xcfd),16);
        if(m->flags&0x40) {ec_push(m,m->si);m->di=m->sp;pn_call(m,0xcde,0x5927,h);gw_al(m,0x93);m->cx=0x29;pn_call(m,0x8810,0x592f,h);m->si=ec_pop(m);}
        m->bx=ec_shl(m,m->bx);m->bx=ec_shl(m,m->bx);ec_put(m,m->ds,(uint16_t)(m->si+0x1c),(uint8_t)(m->bx>>8));m->ax=(uint16_t)(0xff00u|(m->bx>>8));pn_call(m,0x2ad2,0x593e,h);
    }
    m->di=ec_pop(m);
}
static void pn_prisoner(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->dx);pn_call(m,0xece0,0x5944,h);ec_cmp(m,(uint8_t)m->ax,64,8);
    if(m->flags&1) {
        ec_cmp(m,(uint8_t)m->ax,32,8);
        if(!(m->flags&1)) {gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x1c)));ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),8);if(m->flags&0x40) {ec_put(m,m->ds,(uint16_t)(m->si+0x1d),0xff);ec_put(m,m->ds,(uint16_t)(m->si+0x17),0);gw_ah(m,0xff);pn_call(m,0x2ad2,0x5963,h);m->cx=0x42;pn_call(m,0x5990,0x5969,h);}}
        else {gw_al(m,(uint8_t)m->ax&0xf);ec_logic(m,(uint8_t)m->ax,8);gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,8,8,0,0));gw_bl(m,(uint8_t)m->ax);m->ax=ec_math(m,m->si,0x4240,16,0,1);for(unsigned i=0;i<3;++i) m->ax=ec_shl(m,m->ax);gw_al(m,9);m->dx=0xffff;pn_call(m,0x301c,0x5984,h);m->cx=0x41;pn_call(m,0x5990,0x598a,h);ec_put(m,m->ds,(uint16_t)(m->si+0x1c),24);}
    }
    m->dx=ec_pop(m);
}
static void pn_generals(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->si);m->cx=127;m->si=0x4240;
    do {ec_push(m,m->cx);ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);
        if(!(m->flags&1)) {ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x18)),0,8);
            if(!(m->flags&0x40)) {uint8_t v=(uint8_t)st_dec(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x18)),8);if(KI_POLITICS_MUTATION==6) --v;ec_put(m,m->ds,(uint16_t)(m->si+0x18),v);}
            else {ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x1c)),0xff,8);if(m->flags&0x40) pn_call(m,0x5899,0x5883,h);
                else {ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x1d)),0xff,8);if(!(m->flags&0x40)) pn_call(m,0x5940,0x588e,h);}}
        }
        m->si=ec_math(m,m->si,32,16,0,0);m->cx=ec_pop(m);m->cx=(uint16_t)(m->cx-1);
    } while(m->cx);
    m->si=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void pn_decrease_relation(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->bx);ec_push(m,m->cx);pn_call(m,0x3119,0x30f5,h);m->cx=(uint16_t)((uint16_t)ec_byte(m,m->ds,m->bx)*0x101u);m->cx&=0x807f;ec_logic(m,m->cx,16);
    pn_cl(m,(uint8_t)ec_math(m,(uint8_t)m->cx,(uint8_t)m->ax,8,0,1));if(m->flags&1) {pn_cl(m,0);ec_logic(m,0,8);}pn_cl(m,(uint8_t)m->cx|(m->cx>>8));ec_logic(m,(uint8_t)m->cx,8);
    if(KI_POLITICS_MUTATION==2) pn_cl(m,(uint8_t)m->cx&0x7f);ec_put(m,m->ds,m->bx,(uint8_t)m->cx);m->cx=ec_pop(m);m->bx=ec_pop(m);
}
static void pn_pair_address(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);pn_swap(&m->si,&m->di);pn_call(m,0x3119,0x3110,h);m->dx=m->bx;pn_swap(&m->si,&m->di);pn_call(m,0x3119,0x3117,h);m->ax=ec_pop(m);
}
static void pn_power(KiMachine16 *m) {
    ec_push(m,m->dx);m->ax=0;ec_logic(m,0,16);
    for(unsigned i=0;i<3;++i) {m->dx=ec_word(m,m->ds,(uint16_t)(m->bx+4+i*2));m->dx=st_shr(m,m->dx,16);m->dx=st_shr(m,m->dx,16);m->ax=ec_math(m,m->ax,m->dx,16,0,0);}
    ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->bx+0x23)),8);if(!(m->flags&1)) m->ax=KI_POLITICS_MUTATION==5?2001:2000;ec_cmp(m,m->ax,2000,16);if(!(m->flags&(1|0x40))) m->ax=KI_POLITICS_MUTATION==5?2001:2000;
    ec_cmp(m,ec_word(m,m->ds,(uint16_t)(m->bx+0x21)),19,16);if(m->flags&(1|0x40)) {m->ax=0;ec_logic(m,0,16);}m->dx=ec_pop(m);
}
static void pn_neighbor(KiMachine16 *m) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->si);ec_push(m,m->di);
    gw_al(m,ec_byte(m,m->ds,m->di));m->si=ec_math(m,m->di,0x1c,16,0,0);m->cx=4;
    do {
        gw_al(m,(uint8_t)st_shr(m,(uint8_t)m->ax,8));
        if(m->flags&1) {
            m->bx=(uint16_t)((uint16_t)ec_byte(m,m->ds,m->si)<<8);ec_logic(m,0,8);for(unsigned i=0;i<3;++i) m->bx=st_shr(m,m->bx,16);gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x841)));ec_cmp(m,m->bx>>8,m->ax>>8,8);
            if(!(m->flags&0x40)) {
                ec_cmp(m,m->bx>>8,24,8);
                if(m->flags&0x40) {m->di=m->dx;ec_store(m,m->es,m->di,0x600);}
                else {gw_bl(m,0);ec_logic(m,0,8);m->bx=st_shr(m,m->bx,16);m->bx=st_shr(m,m->bx,16);m->bx=ec_math(m,m->bx,0,16,0,0);m->di=ec_inc(m,ec_inc(m,m->dx,16),16);
                    for(;;) {ec_cmp(m,ec_word(m,m->es,m->di),m->bx,16);if(m->flags&0x40) break;ec_cmp(m,ec_word(m,m->es,m->di),0xffff,16);if(m->flags&0x40) {ec_store(m,m->es,m->di,m->bx);break;}m->di=ec_inc(m,m->di,16);m->di=ec_inc(m,m->di,16);}}
            }
        }
        m->si=ec_inc(m,m->si,16);m->cx=(uint16_t)(m->cx-1);
    } while(m->cx);
    m->di=ec_pop(m);m->si=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void pn_frontier(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);if(m->flags&1) return;
    ec_push(m,m->si);ec_push(m,m->di);ec_push(m,m->bp);m->dx=m->di;m->ax=m->si;m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);m->di=0x840;m->cx=192;
    do {ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->di+1)),8);if(m->flags&0x40) pn_call(m,0x2cdf,0x2c71,h);m->di=ec_math(m,m->di,32,16,0,0);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
    m->bx=m->si;m->bx=st_shr(m,m->bx,16);m->bx=st_shr(m,m->bx,16);m->ax=m->bx;m->bx=st_shr(m,m->bx,16);m->bx=ec_math(m,m->bx,m->ax,16,0,0);m->bx=ec_math(m,m->bx,0x600,16,0,0);m->bp=m->dx;m->bp=ec_inc(m,m->bp,16);m->bp=ec_inc(m,m->bp,16);m->cx=0x1515;
    for(;;) {
        ec_cmp(m,ec_word(m,m->es,m->bp),0xffff,16);if(m->flags&0x40) break;m->si=m->bp;gw_dl(m,0xff);
        do {
            m->di=ec_word(m,m->es,m->si);ec_cmp(m,m->di,0xffff,16);
            if(m->flags&0x40) pn_cl(m,1);
            else {m->ax=m->di;gw_ah(m,(uint8_t)(m->ax>>8)&0x7f);ec_logic(m,m->ax>>8,8);m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);gw_al(m,(uint8_t)(m->ax>>8));gw_al(m,ec_byte(m,m->ds,(uint16_t)(m->bx+(uint8_t)m->ax)));ec_cmp(m,(uint8_t)m->dx,(uint8_t)m->ax,8);
                if(KI_POLITICS_MUTATION==7?(m->flags&(1|0x40)):!(m->flags&(1|0x40))) {gw_dl(m,(uint8_t)m->ax);ec_cmp(m,(uint8_t)m->ax,0x80,8);if(m->flags&1) {m->di|=0x8000;ec_logic(m,m->di,16);}uint16_t v=ec_word(m,m->es,m->bp);ec_store(m,m->es,m->bp,m->di);m->di=v;ec_cmp(m,m->si,m->bp,16);if(!(m->flags&0x40)) ec_store(m,m->es,m->si,m->di);}
            }
            m->si=ec_inc(m,m->si,16);m->si=ec_inc(m,m->si,16);pn_cl(m,(uint8_t)st_dec(m,(uint8_t)m->cx,8));
        } while((uint8_t)m->cx);
        pn_ch(m,(uint8_t)st_dec(m,m->cx>>8,8));m->bp=ec_inc(m,m->bp,16);m->bp=ec_inc(m,m->bp,16);pn_cl(m,(uint8_t)(m->cx>>8));ec_cmp(m,m->cx>>8,1,8);if(m->flags&0x40) break;
    }
    m->bp=ec_pop(m);m->di=ec_pop(m);m->si=ec_pop(m);
}
static void pn_capital_event(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->dx);ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),0xff,8);
    if(m->flags&0x40) {pn_call(m,0xece0,0x2d46,h);ec_cmp(m,(uint8_t)m->ax,64,8);if(m->flags&1) {gw_al(m,8);m->dx=0xffff;gw_bl(m,0xff);pn_call(m,0x2fb1,0x2d54,h);}}
    m->dx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void pn_ai_relation(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->di);m->di=ec_word(m,m->es,m->di);ec_cmp(m,m->di,0xffff,16);
    if(!(m->flags&0x40)) {
        m->dx=m->di;m->dx&=0x8000;ec_logic(m,m->dx,16);m->di&=0x7fff;ec_logic(m,m->di,16);pn_call(m,0x310a,0x2dce,h);gw_al(m,ec_byte(m,m->ds,m->bx));ec_cmp(m,(uint8_t)m->ax,0x80,8);
        if(!(m->flags&1)) {gw_al(m,(uint8_t)m->ax&0x7f);ec_logic(m,(uint8_t)m->ax,8);ec_cmp(m,(uint8_t)m->ax,22,8);if(m->flags&1) gw_al(m,22);gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,2,8,0,1));gw_al(m,(uint8_t)m->ax|0x80);ec_logic(m,(uint8_t)m->ax,8);}
        else {ec_cmp(m,m->di,ec_word(m,m->cs,0xcfd),16);if(!(m->flags&0x40)) {ec_cmp(m,(uint8_t)m->ax,50,8);if(m->flags&1) gw_al(m,(uint8_t)ec_inc(m,(uint8_t)m->ax,8));}}
        ec_put(m,m->ds,m->bx,(uint8_t)m->ax);
    }m->di=ec_pop(m);
}
static void pn_player_relation(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->si);ec_push(m,m->di);m->di=ec_word(m,m->es,m->di);ec_cmp(m,m->di,0xffff,16);
    if(!(m->flags&0x40)) {m->bx=m->di;m->di&=0x7fff;ec_logic(m,m->di,16);gw_al(m,1);pn_call(m,0x30f0,0x2e08,h);ec_cmp(m,m->bx>>8,0x80,8);if(m->flags&1) {gw_al(m,7);pn_call(m,0x30f0,0x2e12,h);}}
    m->di=ec_word(m,m->cs,0xcfd);m->si=0;m->cx=22;
    do {ec_cmp(m,m->si,m->di,16);if(!(m->flags&0x40)) {ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);if(!(m->flags&1)) {gw_al(m,1);pn_call(m,0x30f0,0x2e2b,h);}}m->si=ec_math(m,m->si,64,16,0,0);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
    m->di=ec_pop(m);m->si=ec_pop(m);
}
static void pn_cooperate(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->si);ec_push(m,m->di);m->dx=m->si;m->dx=ec_shl(m,m->dx);m->dx=ec_shl(m,m->dx);gw_dl(m,(uint8_t)(m->dx>>8));pn_dh(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)));ec_cmp(m,m->dx>>8,0xff,8);int proceed=!(m->flags&0x40);
    if(proceed) {ec_cmp(m,m->dx>>8,ec_byte(m,m->cs,0xcff),8);proceed=!(m->flags&0x40);}
    if(proceed) {
        m->ax=ec_word(m,m->cs,0xcfd);
        for(;;) {ec_cmp(m,m->ax,ec_word(m,m->es,m->di),16);if(m->flags&0x40) break;ec_cmp(m,ec_word(m,m->es,m->di),0xffff,16);if(m->flags&0x40) {proceed=0;break;}m->di=ec_inc(m,m->di,16);m->di=ec_inc(m,m->di,16);}
        if(proceed) {m->di=m->ax;pn_call(m,0x30cb,0x2e66,h);ec_cmp(m,(uint8_t)m->ax,0x80,8);if(!(m->flags&1)) {m->ax=(uint16_t)((m->dx>>8)<<8);ec_logic(m,0,8);m->ax=st_shr(m,m->ax,16);m->ax=st_shr(m,m->ax,16);m->si=m->ax;pn_call(m,0x30cb,0x2e77,h);ec_cmp(m,(uint8_t)m->ax,KI_POLITICS_MUTATION==8?0xa4:0xa3,8);if(!(m->flags&1)) {pn_swap(&m->si,&m->di);gw_al(m,2);gw_bl(m,0xff);pn_call(m,0x2fb1,0x2e84,h);}}}
    }
    m->di=ec_pop(m);m->si=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void pn_ceasefire(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->di);ec_push(m,m->bp);ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);int found=0;
    if(!(m->flags&0x40)) {
        m->cx=m->si;m->cx=ec_shl(m,m->cx);m->cx=ec_shl(m,m->cx);m->bx=m->si;pn_call(m,0x3091,0x2ea0,h);m->dx=m->ax;pn_cl(m,21);m->bp=ec_word(m,m->es,m->di);
        for(;;) {
            m->bx=ec_word(m,m->es,m->di);ec_cmp(m,m->bx,0xffff,16);if(m->flags&0x40) break;ec_cmp(m,m->bx>>8,0x80,8);if(m->flags&1) break;
            gw_bh(m,(uint8_t)(m->bx>>8)&0x7f);ec_logic(m,m->bx>>8,8);pn_call(m,0x3091,0x2eba,h);m->dx=ec_math(m,m->dx,m->ax,16,0,1);
            if(m->flags&(1|0x40)) {found=1;break;}m->di=ec_inc(m,m->di,16);m->di=ec_inc(m,m->di,16);pn_cl(m,(uint8_t)st_dec(m,(uint8_t)m->cx,8));if(!(uint8_t)m->cx) break;
        }
        if(found) {
            ec_push(m,m->dx);
            for(;;) {m->bx=ec_word(m,m->es,m->di);ec_cmp(m,m->bp,m->bx,16);int is_first=(m->flags&0x40)!=0;
                if(!is_first) {ec_cmp(m,m->bx>>8,0x80,8);if(m->flags&1) break;ec_cmp(m,m->bx,0xffff,16);if(m->flags&0x40) break;gw_bh(m,(uint8_t)(m->bx>>8)&0x7f);ec_logic(m,m->bx>>8,8);m->bx=ec_shl(m,m->bx);m->bx=ec_shl(m,m->bx);m->dx=(uint16_t)(0xff00u|(m->bx>>8));gw_al(m,3);gw_bl(m,0xff);pn_call(m,0x2fb1,0x2eea,h);}
                m->di=ec_inc(m,m->di,16);m->di=ec_inc(m,m->di,16);pn_cl(m,(uint8_t)st_dec(m,(uint8_t)m->cx,8));if(!(uint8_t)m->cx) break;
            }
            m->dx=ec_pop(m);m->flags|=1;
        }
    }
    if(!found) m->flags&=(uint16_t)~1u;m->bp=ec_pop(m);m->di=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void pn_declare_war(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->di);m->ax=ec_word(m,m->es,m->di);gw_ah(m,(uint8_t)(m->ax>>8)&0x7f);ec_logic(m,m->ax>>8,8);m->bx=m->ax;m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),8);int allow=!(m->flags&0x40);
    if(allow) {m->ax=ec_byte(m,m->ds,(uint16_t)(m->si+0x23));ec_logic(m,0,8);for(unsigned i=0;i<4;++i) m->ax=ec_shl(m,m->ax);m->ax=ec_math(m,m->ax,64,16,0,0);ec_cmp(m,m->ax,0x61a,16);if(!(m->flags&(1|0x40))) m->ax=0x61a;ec_cmp(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+0x21)),16);allow=ec_less(m);}
    if(allow) {m->di=m->bx;gw_bl(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x28)));gw_bh(m,(uint8_t)m->bx);gw_bh(m,(uint8_t)st_shr(m,m->bx>>8,8));gw_bl(m,(uint8_t)ec_math(m,(uint8_t)m->bx,m->bx>>8,8,0,0));gw_bl(m,(uint8_t)ec_math(m,(uint8_t)m->bx,20,8,0,0));gw_bl(m,(uint8_t)m->bx|0x80);ec_logic(m,(uint8_t)m->bx,8);pn_call(m,0x30cb,0x2f3e,h);ec_cmp(m,(uint8_t)m->ax,(uint8_t)m->bx,8);allow=(m->flags&(1|0x40))!=0;}
    if(allow) {m->bx=m->si;pn_call(m,0x3091,0x2f47,h);m->cx=m->ax;m->bx=m->di;pn_call(m,0x3091,0x2f4e,h);m->bx=m->ax;m->bx=st_shr(m,m->bx,16);m->bx=st_shr(m,m->bx,16);m->ax=ec_math(m,m->ax,m->bx,16,0,1);ec_cmp(m,m->cx,m->ax,16);allow=!(m->flags&1);}
    if(allow) {m->ax=m->di;m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);m->dx=(uint16_t)(0xff00u|(m->ax>>8));gw_al(m,1);gw_bl(m,0xff);pn_call(m,0x2fb1,0x2f6b,h);m->di=ec_pop(m);m->flags&=(uint16_t)~1u;}
    else {m->di=ec_pop(m);m->flags|=1;}
}
static void pn_neutral(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),24,8);if(m->flags&1) return;ec_cmp(m,ec_word(m,m->es,(uint16_t)(m->di-2)),0xffff,16);int clear=(m->flags&0x40)!=0;
    if(!clear) {m->ax=ec_byte(m,m->ds,(uint16_t)(m->si+0x23));ec_logic(m,0,8);for(unsigned i=0;i<4;++i) m->ax=ec_shl(m,m->ax);m->ax=ec_math(m,m->ax,96,16,0,0);ec_cmp(m,m->ax,0x6dd,16);if(!(m->flags&(1|0x40))) m->ax=0x6dd;ec_cmp(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+0x21)),16);clear=!ec_less(m);
        if(!clear) {ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),24,8);if(!(m->flags&0x40)) {gw_al(m,1);m->dx=0xff18;gw_bl(m,0xff);pn_call(m,0x2fb1,0x2fab,h);}}}
    if(clear) ec_put(m,m->ds,(uint16_t)(m->si+0x19),0xff);
}
static void pn_decisions(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);if(m->flags&1) return;m->es=ec_word(m,m->cs,0x987c);m->ax=m->si;m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);gw_al(m,48);m->ax=(uint16_t)((uint8_t)m->ax*(m->ax>>8));m->ax=ec_inc(m,m->ax,16);m->ax=ec_inc(m,m->ax,16);m->di=m->ax;
    ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);if(m->flags&0x40) pn_call(m,0x2df3,0x2d80,h);else pn_call(m,0x2db8,0x2d7b,h);
    pn_call(m,0x2e33,0x2d83,h);pn_call(m,0x2e89,0x2d86,h);pn_call(m,0x2efb,0x2d89,h);
    if(m->flags&1) {pn_call(m,0x2f71,0x2d8e,h);ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),24,8);if(m->flags&0x40) return;}
    ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),24,8);int clear=(m->flags&0x40)!=0;
    if(!clear) {m->ax=ec_word(m,m->es,m->di);ec_cmp(m,m->ax>>8,0xff,8);clear=(m->flags&0x40)!=0;
        if(!clear) {ec_cmp(m,m->ax>>8,0x80,8);clear=(m->flags&1)!=0;
            if(!clear) {gw_ah(m,(uint8_t)(m->ax>>8)&0x7f);ec_logic(m,m->ax>>8,8);m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->si+0x19)),8);clear=!(m->flags&0x40);}}}
    if(clear) ec_put(m,m->ds,(uint16_t)(m->si+0x19),0xff);
}
static void pn_initialize(KiMachine16 *m,const KiEconomyHooks *h) {
    ec_push(m,m->ds);ec_push(m,m->es);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);
    ec_store(m,m->cs,0xd20,0);ec_put(m,m->cs,0x31ad,7);m->ax=ec_word(m,m->cs,0xd56);m->ds=m->ax;m->es=m->ax;m->ax=0;ec_logic(m,0,16);m->di=0;m->si=KI_POLITICS_MUTATION==3?0:0x100;m->cx=0x180;m->flags&=(uint16_t)~0x400u;
    do {ec_store(m,m->es,m->di,ec_word(m,m->ds,m->si));m->si=(uint16_t)(m->si+2);m->di=(uint16_t)(m->di+2);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
    m->cx=0x80;do {ec_store(m,m->es,m->di,m->ax);m->di=(uint16_t)(m->di+2);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
    m->es=ec_word(m,m->cs,0x987c);m->di=0;ec_logic(m,0,16);m->ax=0xffff;m->cx=0x210;
    do {ec_store(m,m->es,m->di,m->ax);m->di=(uint16_t)(m->di+2);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
    m->cx=22;m->ds=ec_word(m,m->cs,0xd52);m->si=0;m->di=0;ec_logic(m,0,16);
    do {ec_push(m,m->cx);ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);if(!(m->flags&1)) {pn_call(m,0x2c52,0x2c2d,h);pn_call(m,0x2d3a,0x2c30,h);}m->di=ec_math(m,m->di,48,16,0,0);m->cx=ec_pop(m);m->si=ec_math(m,m->si,64,16,0,0);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
    m->si=0;m->cx=22;do {ec_push(m,m->cx);pn_call(m,0x2d58,0x2c43,h);m->cx=ec_pop(m);m->si=ec_math(m,m->si,64,16,0,0);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
    m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->es=ec_pop(m);m->ds=ec_pop(m);
}
int ki_politics_body(KiMachine16 *m,uint16_t target,const KiEconomyHooks *h) {
    switch(target) {
    case 0x585f:pn_generals(m,h);break;case 0x5899:pn_join(m,h);break;case 0x5940:pn_prisoner(m,h);break;
    case 0x2ad2:pn_assign(m);break;case 0x5990:pn_notice(m,h);break;case 0x301c:pn_writer(m);break;
    case 0x2bd9:pn_initialize(m,h);break;case 0x2c52:pn_frontier(m,h);break;case 0x2cdf:pn_neighbor(m);break;
    case 0x2d3a:pn_capital_event(m,h);break;case 0x2fb1:pn_event(m,h);break;case 0x2d58:pn_decisions(m,h);break;
    case 0x2db8:pn_ai_relation(m,h);break;case 0x2df3:pn_player_relation(m,h);break;case 0x30f0:pn_decrease_relation(m,h);break;
    case 0x310a:pn_pair_address(m,h);break;case 0x3091:pn_power(m);break;case 0x2e33:pn_cooperate(m,h);break;
    case 0x2e89:pn_ceasefire(m,h);break;case 0x2efb:pn_declare_war(m,h);break;case 0x2f71:pn_neutral(m,h);break;
    default:return 0;
    }return 1;
}
#define PN_WRAP(name,target) void name(KiMachine16 *m,const KiEconomyHooks *h) {ki_politics_body(m,target,h);m->ip=ec_pop(m);}
PN_WRAP(sub_1585F,0x585f) PN_WRAP(sub_15899,0x5899) PN_WRAP(sub_15940,0x5940)
PN_WRAP(sub_12AD2,0x2ad2) PN_WRAP(sub_15990,0x5990) PN_WRAP(sub_1301C,0x301c)
PN_WRAP(sub_12BD9,0x2bd9) PN_WRAP(sub_12C52,0x2c52) PN_WRAP(sub_12CDF,0x2cdf)
PN_WRAP(sub_12D3A,0x2d3a) PN_WRAP(sub_12FB1,0x2fb1) PN_WRAP(sub_12D58,0x2d58)
PN_WRAP(sub_12DB8,0x2db8) PN_WRAP(sub_12DF3,0x2df3) PN_WRAP(sub_130F0,0x30f0)
PN_WRAP(sub_1310A,0x310a) PN_WRAP(sub_13091,0x3091) PN_WRAP(sub_12E33,0x2e33)
PN_WRAP(sub_12E89,0x2e89) PN_WRAP(sub_12EFB,0x2efb) PN_WRAP(sub_12F71,0x2f71)
#undef PN_WRAP
void ki_politics_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *h) {
    switch(target) {
    case 0x585f:case 0x5899:case 0x5940:case 0x2ad2:case 0x5990:case 0x301c:
    case 0x2bd9:case 0x2c52:case 0x2cdf:case 0x2d3a:case 0x2fb1:case 0x2d58:
    case 0x2db8:case 0x2df3:case 0x30f0:case 0x310a:case 0x3091:case 0x2e33:case 0x2e89:case 0x2efb:case 0x2f71:
        if(h&&h->enter) h->enter(m,target,h->user);ki_politics_body(m,target,h);m->ip=ec_pop(m);break;
    default:ki_world_invoke(m,target,h);break;
    }
}
