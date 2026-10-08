/* Five original DOS/V entries, docs/spec/207. Include after economy.c.
 * Shared 16-bit arithmetic/stack adapter remains unchanged.
 */
#include "settlement.h"
#ifndef KI_SETTLEMENT_MUTATION
#define KI_SETTLEMENT_MUTATION 0
#endif
static uint16_t st_shr(KiMachine16 *m,uint16_t a,unsigned width) {
    uint16_t bits=a&1,sign=width==8?0x80:0x8000,r=a>>1;
    if (a&sign) bits|=0x800;
    m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|bits),r,width);return r;
}
static uint16_t st_dec(KiMachine16 *m,uint16_t a,unsigned width) {
    uint16_t carry=m->flags&1,r=ec_math(m,a,1,width,0,1);
    m->flags=(uint16_t)((m->flags&~1u)|carry);return r;
}
static void st_div(KiMachine16 *m,uint16_t divisor) {
    /* Legal inputs only. DIV flags are overwritten before every observable boundary. */
    uint32_t v=((uint32_t)m->dx<<16)|m->ax;
    m->ax=(uint16_t)(v/divisor);m->dx=(uint16_t)(v%divisor);
}
static void st_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *hooks) {
    ec_push(m,next);m->ip=target;ki_settlement_invoke(m,target,hooks);
}
static void st_income_body(KiMachine16 *m) {
    m->ax=ec_word(m,m->ds,(uint16_t)(m->di+0xe));m->dx=0;ec_logic(m,0,16);st_div(m,m->bx);
    ec_store(m,m->ss,m->bp,ec_math(m,ec_word(m,m->ss,m->bp),m->ax,16,0,0));
    unsigned carry=KI_SETTLEMENT_MUTATION==1?0:m->flags&1;
    ec_put(m,m->ss,(uint16_t)(m->bp+2),(uint8_t)ec_math(m,ec_byte(m,m->ss,(uint16_t)(m->bp+2)),0,8,carry,0));
}
static void st_recruit_body(KiMachine16 *m) {
    ec_push(m,m->bx);ec_push(m,m->cx);
    m->ax=ec_word(m,m->ds,(uint16_t)(m->di+0xe));m->dx=0;ec_logic(m,0,16);st_div(m,m->bx);
    for (unsigned i=0;i<5;++i) m->ax=st_shr(m,m->ax,16);
    uint16_t base=m->ax;m->bx=ec_word(m,m->ds,(uint16_t)(m->di+0xa));ec_cmp(m,m->bx,150,16);
    if (!(m->flags&1)) {
        m->dx=st_shr(m,m->ax,16);m->cx=(uint16_t)((m->cx&0xff00u)|5);
        for (unsigned i=0;i<5;++i) m->ax=st_shr(m,m->ax,16);
        m->cx=m->dx;m->dx=ec_math(m,m->dx,m->ax,16,0,1);
    } else {
        ec_cmp(m,m->bx,80,16);
        if (!(m->flags&1)) {
            m->dx=m->ax;for (unsigned i=0;i<3;++i) m->ax=st_shr(m,m->ax,16);
            m->cx=m->ax;m->dx=ec_math(m,m->dx,m->ax,16,0,1);m->dx=ec_math(m,m->dx,m->cx,16,0,1);
        } else {
            m->dx=m->ax;for (unsigned i=0;i<2;++i) m->dx=st_shr(m,m->dx,16);
            m->cx=st_shr(m,m->dx,16);m->dx=ec_math(m,m->dx,m->cx,16,0,0);
            m->cx=st_shr(m,m->cx,16);m->cx=st_shr(m,m->cx,16);
            m->ax=ec_math(m,m->ax,m->dx,16,0,1);m->ax=ec_math(m,m->ax,m->cx,16,0,1);
            if (KI_SETTLEMENT_MUTATION==2) m->ax=(uint16_t)((uint32_t)base*19/32);
        }
    }
    ec_store(m,m->ss,(uint16_t)(m->bp+4),ec_math(m,ec_word(m,m->ss,(uint16_t)(m->bp+4)),m->ax,16,0,0));
    ec_store(m,m->ss,(uint16_t)(m->bp+6),ec_math(m,ec_word(m,m->ss,(uint16_t)(m->bp+6)),m->cx,16,0,0));
    ec_store(m,m->ss,(uint16_t)(m->bp+8),ec_math(m,ec_word(m,m->ss,(uint16_t)(m->bp+8)),m->dx,16,0,0));
    m->cx=ec_pop(m);m->bx=ec_pop(m);
}
static void st_ai_body(KiMachine16 *m) {
    uint8_t hi=(uint8_t)st_shr(m,ec_byte(m,m->ss,(uint16_t)(m->bp+2)),8);
    ec_put(m,m->ss,(uint16_t)(m->bp+2),hi);
    uint16_t before=ec_word(m,m->ss,m->bp),carry=m->flags&1,r=(uint16_t)((before>>1)|(carry<<15));
    m->flags=(uint16_t)((m->flags&~0x801u)|(before&1)|((((r>>15)^(r>>14))&1)<<11));ec_store(m,m->ss,m->bp,r);
    m->ax=m->si;m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);m->ax=(uint16_t)((m->ax&0xff00u)|0x80);
    m->dx=0;ec_logic(m,0,16);m->cx=127;m->bx=0x2240;
    do {
        ec_cmp(m,ec_byte(m,m->ds,m->bx),(uint8_t)m->ax,8);
        if (!(m->flags&1)) {
            ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->bx+1)),m->ax>>8,8);
            if (m->flags&0x40) m->dx=ec_math(m,m->dx,ec_word(m,m->ds,(uint16_t)(m->bx+4)),16,0,0);
        }
        m->bx=ec_math(m,m->bx,KI_SETTLEMENT_MUTATION==4?0x40:0x20,16,0,0);m->cx=(uint16_t)(m->cx-1);
    } while (m->cx);
    m->ax=m->dx>>8;ec_logic(m,0,8);
    m->ax=ec_math(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+0x1b)),16,0,0);m->ax=ec_shl(m,m->ax);
    ec_cmp(m,m->ax,ec_word(m,m->ss,(uint16_t)(m->bp+1)),16);m->flags=(uint16_t)((m->flags&~1u)|((m->flags&1)?0:1));
}
static void st_player_body(KiMachine16 *m) {
    uint32_t original=(uint32_t)ec_word(m,m->ss,m->bp)|((uint32_t)ec_byte(m,m->ss,(uint16_t)(m->bp+2))<<16);
    ec_push(m,m->si);ec_push(m,m->di);
    m->ax=ec_byte(m,m->cs,0xd08);ec_logic(m,0,8);
    uint32_t product=(uint32_t)m->ax*ec_word(m,m->ss,(uint16_t)(m->bp+1));m->ax=(uint16_t)product;m->dx=(uint16_t)(product>>16);
    m->bx=100;st_div(m,m->bx);m->bx=m->dx;ec_store(m,m->ss,(uint16_t)(m->bp+1),m->ax);
    m->ax=ec_byte(m,m->cs,0xd08);ec_logic(m,0,8);m->ax=(uint16_t)(m->ax*ec_byte(m,m->ss,m->bp));
    m->dx=0;ec_logic(m,0,16);
    uint8_t ah=(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->bx,8,0,0);m->ax=(uint16_t)((m->ax&0xffu)|((uint16_t)ah<<8));
    m->dx=(uint16_t)ec_math(m,(uint8_t)m->dx,m->bx>>8,8,m->flags&1,0);m->bx=100;st_div(m,m->bx);
    ec_put(m,m->ss,m->bp,0);ec_store(m,m->ss,m->bp,ec_math(m,ec_word(m,m->ss,m->bp),m->ax,16,0,0));
    if (KI_SETTLEMENT_MUTATION==3) {
        uint32_t v=original*ec_byte(m,m->cs,0xd08)/100;ec_store(m,m->ss,m->bp,(uint16_t)v);ec_put(m,m->ss,(uint16_t)(m->bp+2),(uint8_t)(v>>16));
    }
    m->si=m->bp;m->si=ec_math(m,m->si,4,16,0,0);m->di=0xd0a;m->cx=3;
    do {
        m->ax=ec_word(m,m->cs,m->di);ec_cmp(m,m->ax,ec_word(m,m->ss,m->si),16);
        if (m->flags&1) ec_store(m,m->ss,m->si,m->ax);
        m->si=ec_inc(m,m->si,16);m->si=ec_inc(m,m->si,16);m->di=ec_inc(m,m->di,16);m->di=ec_inc(m,m->di,16);m->cx=(uint16_t)(m->cx-1);
    } while (m->cx);
    m->di=ec_pop(m);m->si=ec_pop(m);m->ax=ec_word(m,m->ss,m->bp);
    m->dx=(uint16_t)((m->dx&0xff00u)|ec_byte(m,m->ss,(uint16_t)(m->bp+2)));
    ec_store(m,m->cs,0xd02,m->ax);ec_put(m,m->cs,0xd04,(uint8_t)m->dx);
    m->ax=ec_word(m,m->ds,(uint16_t)(m->si+0x1a));m->dx=(uint16_t)((m->dx&0xff00u)|ec_byte(m,m->ds,(uint16_t)(m->si+0x1c)));
    ec_store(m,m->cs,0xd05,m->ax);ec_put(m,m->cs,0xd07,(uint8_t)m->dx);
}
static void st_settle_body(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->di);ec_push(m,m->bp);
    m->sp=ec_math(m,m->sp,10,16,0,1);m->bp=m->sp;ec_push(m,m->si);
    m->bx=(uint16_t)((m->bx&0xffu)|((uint16_t)ec_byte(m,m->ds,(uint16_t)(m->si+3))<<8));
    m->bx&=0xff00;ec_logic(m,0,8);for (unsigned i=0;i<3;++i) m->bx=st_shr(m,m->bx,16);
    m->bx=ec_math(m,m->bx,0x840,16,0,0);m->si=m->bx;m->di=0x840;m->dx=0;ec_logic(m,0,16);
    for (unsigned i=0;i<10;i+=2) ec_store(m,m->ss,(uint16_t)(m->bp+i),0);m->cx=(uint16_t)((m->cx&0xffu)|0xc000);
    do {
        ec_cmp(m,(uint8_t)m->cx,ec_byte(m,m->ds,(uint16_t)(m->di+1)),8);
        int selected=(m->flags&0x40)!=0;if (KI_SETTLEMENT_MUTATION==5) selected=!selected;
        if (selected) {
            st_call(m,0x54fc,0x53ff,hooks);st_call(m,0x5538,0x5402,hooks);st_call(m,0x5547,0x5405,hooks);
            ec_put(m,m->ss,(uint16_t)(m->bp+3),(uint8_t)ec_inc(m,ec_byte(m,m->ss,(uint16_t)(m->bp+3)),8));
        }
        m->di=ec_math(m,m->di,0x20,16,0,0);m->cx=(uint16_t)((m->cx&0xffu)|(st_dec(m,m->cx>>8,8)<<8));
    } while (m->cx>>8);
    m->si=ec_pop(m);ec_cmp(m,m->si,ec_word(m,m->cs,0xcfd),16);
    int recruit=1;
    if (m->flags&0x40) st_call(m,0x548f,0x5421,hooks);
    else { st_call(m,0x5456,0x541a,hooks);recruit=!(m->flags&1); }
    if (recruit) for (unsigned i=0;i<3;++i) {
        m->ax=ec_word(m,m->ss,(uint16_t)(m->bp+4+i*2));m->dx=ec_word(m,m->ds,(uint16_t)(m->si+4+i*2));
        st_call(m,0x55ec,(uint16_t)(0x542a + i*12),hooks);ec_store(m,m->ds,(uint16_t)(m->si+4+i*2),m->ax);
    }
    m->ax=ec_word(m,m->ss,m->bp);m->dx=ec_word(m,m->ss,(uint16_t)(m->bp+2));ec_put(m,m->ds,(uint16_t)(m->si+0x23),(uint8_t)(m->dx>>8));
    m->sp=ec_math(m,m->sp,10,16,0,0);m->bp=ec_pop(m);m->di=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);
}
int ki_settlement_body(KiMachine16 *m,uint16_t target,const KiEconomyHooks *hooks) {
    switch (target) {
    case 0x53c6:st_settle_body(m,hooks);break;
    case 0x5456:st_ai_body(m);break;
    case 0x548f:st_player_body(m);break;
    case 0x5538:st_income_body(m);break;
    case 0x5547:st_recruit_body(m);break;
    default:return 0;
    }
    return 1;
}
void sub_153C6(KiMachine16 *m,const KiEconomyHooks *hooks) {st_settle_body(m,hooks);m->ip=ec_pop(m);}
void sub_15456(KiMachine16 *m) {st_ai_body(m);m->ip=ec_pop(m);}
void sub_1548F(KiMachine16 *m) {st_player_body(m);m->ip=ec_pop(m);}
void sub_15538(KiMachine16 *m) {st_income_body(m);m->ip=ec_pop(m);}
void sub_15547(KiMachine16 *m) {st_recruit_body(m);m->ip=ec_pop(m);}
void ki_settlement_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *hooks) {
    if (target==0x53c6||target==0x5456||target==0x548f||target==0x5538||target==0x5547) {
        if (hooks&&hooks->enter) hooks->enter(m,target,hooks->user);
        ki_settlement_body(m,target,hooks);m->ip=ec_pop(m);
    } else ki_economy_invoke(m,target,hooks);
}
