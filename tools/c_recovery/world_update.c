/* Eleven original DOS/V entries, docs/spec/209; include after settlement.c. */
#include "world_update.h"
#ifndef KI_WORLD_MUTATION
#define KI_WORLD_MUTATION 0
#endif
static void gw_al(KiMachine16 *m,uint8_t v) {m->ax=(uint16_t)((m->ax&0xff00u)|v);}
static void gw_ah(KiMachine16 *m,uint8_t v) {m->ax=(uint16_t)((m->ax&0xffu)|((uint16_t)v<<8));}
static void gw_bl(KiMachine16 *m,uint8_t v) {m->bx=(uint16_t)((m->bx&0xff00u)|v);}
static void gw_bh(KiMachine16 *m,uint8_t v) {m->bx=(uint16_t)((m->bx&0xffu)|((uint16_t)v<<8));}
static void gw_dl(KiMachine16 *m,uint8_t v) {m->dx=(uint16_t)((m->dx&0xff00u)|v);}
static uint8_t gw_shl8(KiMachine16 *m,uint8_t a) {
    uint8_t r=(uint8_t)(a<<1);uint16_t bits=(uint16_t)(a>>7);
    if (((r>>7)^(a>>7))&1) bits|=0x800;
    m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|bits),r,8);return r;
}
static uint16_t gw_sar(KiMachine16 *m,uint16_t a) {
    uint16_t r=(uint16_t)((a>>1)|(a&0x8000));m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|(a&1)),r,16);return r;
}
static void gw_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *hooks) {
    ec_push(m,next);m->ip=target;ki_world_invoke(m,target,hooks);
}
static void gw_growth(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);m->si=0x840;m->cx=192;
    do {
        m->bx=ec_byte(m,m->ds,(uint16_t)(m->si+0x10));ec_logic(m,0,8);m->bx=ec_math(m,m->bx,100,16,0,1);
        gw_ah(m,ec_byte(m,m->cs,0xcff));ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->si+1)),8);
        if (m->flags&0x40) {m->ax=ec_byte(m,m->cs,0xd08);ec_logic(m,0,8);m->ax=ec_math(m,m->ax,30,16,0,1);m->bx=ec_math(m,m->bx,m->ax,16,0,1);}
        m->ax=ec_byte(m,m->ds,(uint16_t)(m->si+0xf));ec_logic(m,0,8);ec_logic(m,(uint8_t)m->ax,8);if (!m->ax) m->ax=1;
        int32_t product=(int32_t)(int16_t)m->ax*(int16_t)m->bx;m->ax=(uint16_t)product;
        m->dx=ec_word(m,m->ds,(uint16_t)(m->si+0xe));ec_cmp(m,m->ax>>8,0,8);
        if (!ec_less(m)) {m->ax=gw_sar(m,m->ax);m->ax=ec_math(m,m->ax,m->dx,16,0,0);}
        else {m->ax=ec_math(m,0,m->ax,16,0,1);uint16_t v=m->ax;m->ax=m->dx;m->dx=v;m->ax=ec_math(m,m->ax,m->dx,16,0,1);if (m->flags&1) {m->ax=0;ec_logic(m,0,16);}}
        if (KI_WORLD_MUTATION==1) {
            int32_t v=ec_word(m,m->ds,(uint16_t)(m->si+0xe));v+=product>=0?product/2:product;
            if (v<0) v=0;if (v>65535) v=65535;m->ax=(uint16_t)v;
        }
        ec_cmp(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+0xc)),16);if (!(m->flags&(1|0x40))) m->ax=ec_word(m,m->ds,(uint16_t)(m->si+0xc));
        ec_store(m,m->ds,(uint16_t)(m->si+0xe),m->ax);gw_call(m,0xece0,0x56ef,hooks);m->ax&=0xf;ec_logic(m,m->ax,16);
        m->bx=ec_math(m,m->bx,m->ax,16,0,1);ec_cmp(m,m->bx,100,16);if (ec_greater(m)) m->bx=100;
        ec_cmp(m,m->bx,0xff9c,16);if (ec_less(m)) m->bx=0xff9c;m->bx=ec_math(m,m->bx,100,16,0,0);ec_put(m,m->ds,(uint16_t)(m->si+0x10),(uint8_t)m->bx);
        m->si=ec_math(m,m->si,32,16,0,0);m->cx=(uint16_t)(m->cx-1);
    } while(m->cx);
    m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void gw_score(KiMachine16 *m) {
    ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);m->dx=127;m->si=0x4240;
    do {
        ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);
        if (!(m->flags&1)) {
            m->bx=ec_math(m,m->si,14,16,0,0);m->ax=0;ec_logic(m,0,16);m->cx=3;
            do {gw_al(m,ec_byte(m,m->ds,m->bx));for(unsigned i=0;i<(KI_WORLD_MUTATION==2?3:4);++i) gw_al(m,(uint8_t)st_shr(m,(uint8_t)m->ax,8));gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->ax,8,0,0));m->bx=ec_inc(m,m->bx,16);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
            m->cx=2;do {gw_al(m,ec_byte(m,m->ds,m->bx));gw_al(m,gw_shl8(m,(uint8_t)m->ax));gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->ax,8,0,0));m->bx=ec_inc(m,m->bx,16);m->cx=(uint16_t)(m->cx-1);} while(m->cx);
            ec_put(m,m->ds,(uint16_t)(m->si+0x1f),(uint8_t)(m->ax>>8));
        }
        m->si=ec_math(m,m->si,32,16,0,0);m->dx=st_dec(m,m->dx,16);
    } while(m->dx);
    m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
static void gw_writer(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);m->cx=m->ax;ec_cmp(m,(uint8_t)m->bx,0xff,8);
    if (m->flags&0x40) {gw_call(m,0xece0,0x2fce,hooks);gw_ah(m,0);ec_logic(m,0,8);gw_al(m,(uint8_t)m->ax&0x7c);ec_logic(m,(uint8_t)m->ax,8);m->bx=m->ax;}
    else {gw_bh(m,0);ec_logic(m,0,8);m->bx=ec_shl(m,m->bx);m->bx=ec_shl(m,m->bx);}
    m->bx=ec_math(m,m->bx,ec_word(m,m->cs,0xd20),16,0,0);ec_cmp(m,m->bx,0x100,16);int written=0;
    if (m->flags&1) {
        m->ds=ec_word(m,m->cs,0xd56);
        do {
            uint16_t value=KI_WORLD_MUTATION==3?ec_word(m,m->ds,m->bx):ec_byte(m,m->ds,m->bx);
            ec_cmp(m,value,0,KI_WORLD_MUTATION==3?16:8);
            if (m->flags&0x40) {ec_store(m,m->ds,m->bx,m->cx);ec_store(m,m->ds,(uint16_t)(m->bx+2),m->dx);m->flags&=(uint16_t)~1u;written=1;break;}
            m->bx=ec_math(m,m->bx,4,16,0,0);ec_cmp(m,m->bx,0x100,16);
        } while(m->flags&1);
    }
    if (!written) m->flags|=1;m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
static void gw_trust(KiMachine16 *m,const KiEconomyHooks *hooks) {
    m->bx=ec_word(m,m->cs,0xcfd);m->ax=ec_word(m,m->ds,(uint16_t)(m->bx+0x21));ec_cmp(m,m->ax>>8,0x80,8);
    if (m->flags&1) return;m->ax=ec_math(m,0,m->ax,16,0,1);ec_cmp(m,m->ax,39,16);if (m->flags&1) return;
    gw_call(m,0xece0,0x5815,hooks);gw_al(m,(uint8_t)m->ax&0xf);ec_logic(m,(uint8_t)m->ax,8);ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->bx+0x28)),8);
    if (!(m->flags&1)) return;m->ax=13;m->dx=0x196;gw_bl(m,0xff);gw_call(m,0x2fbf,0x5827,hooks);
}
static void gw_governor(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);
    m->si=ec_word(m,m->cs,0xcfd);m->cx=(uint16_t)(((uint16_t)ec_byte(m,m->cs,0xcff)<<8)|192);m->di=0x840;
    do {
        ec_cmp(m,m->cx>>8,ec_byte(m,m->ds,(uint16_t)(m->di+1)),8);
        if (m->flags&0x40) {
            gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->di+0x19)));ec_cmp(m,m->bx>>8,0xff,8);
            if (!(m->flags&0x40)) {
                gw_bl(m,0);ec_logic(m,0,8);for(unsigned i=0;i<3;++i) m->bx=st_shr(m,m->bx,16);ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x425a)),0,8);
                if (m->flags&0x40) {
                    m->dx=0;ec_logic(m,0,16);m->ax=180;
                    for(unsigned i=0;i<3;++i) {gw_al(m,i==2?ec_byte(m,m->ds,(uint16_t)(m->di+0x12)):180);uint8_t b=ec_byte(m,m->ds,(uint16_t)(m->di+(i==0?0x10:i==1?0x11:0x13)));gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,b,8,0,1));if (!(m->flags&(1|0x40))) m->dx=ec_math(m,m->dx,m->ax,16,0,0);}
                    m->dx=st_shr(m,m->dx,16);uint32_t product=(uint32_t)m->dx*(KI_WORLD_MUTATION==4?51:50);m->ax=(uint16_t)product;m->dx=m->ax;
                    m->ax=ec_math(m,m->di,0x840,16,0,1);for(unsigned i=0;i<3;++i) m->ax=ec_shl(m,m->ax);gw_al(m,4);gw_bl(m,0xff);gw_call(m,0x2fbf,0x5781,hooks);
                }
            }
        }
        m->di=ec_math(m,m->di,32,16,0,0);m->cx=(uint16_t)((m->cx&0xff00u)|st_dec(m,(uint8_t)m->cx,8));
    } while((uint8_t)m->cx);
    m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void gw_relation_address(KiMachine16 *m) {
    ec_push(m,m->ax);m->bx=m->si;m->bx=st_shr(m,m->bx,16);m->bx=st_shr(m,m->bx,16);m->ax=m->bx;m->bx=st_shr(m,m->bx,16);m->bx=ec_math(m,m->bx,m->ax,16,0,0);
    m->ax=m->di;m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);m->ax=m->ax>>8;ec_logic(m,0,8);m->bx=ec_math(m,m->bx,m->ax,16,0,0);m->bx=ec_math(m,m->bx,0x600,16,0,0);m->ax=ec_pop(m);
}
static void gw_relation(KiMachine16 *m,const KiEconomyHooks *hooks) {ec_push(m,m->bx);gw_call(m,0x3119,0x30cf,hooks);gw_al(m,ec_byte(m,m->ds,m->bx));m->bx=ec_pop(m);}
static void gw_diplomat(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);ec_push(m,m->di);m->si=ec_word(m,m->cs,0xcfd);m->cx=(uint16_t)(((uint16_t)ec_byte(m,m->cs,0xcff)<<8)|22);m->di=0;
    do {
        ec_cmp(m,m->si,m->di,16);
        if (!(m->flags&0x40)) {
            gw_bh(m,ec_byte(m,m->ds,(uint16_t)(m->di+0x2a)));ec_cmp(m,m->bx>>8,0xff,8);
            if (!(m->flags&0x40)) {
                gw_bl(m,0);ec_logic(m,0,8);for(unsigned i=0;i<3;++i) m->bx=st_shr(m,m->bx,16);ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x425a)),0,8);
                if (m->flags&0x40) {
                    gw_call(m,0x30cb,0x57c2,hooks);gw_dl(m,(uint8_t)m->ax);uint16_t v=m->si;m->si=m->di;m->di=v;gw_call(m,0x30cb,0x57c9,hooks);v=m->si;m->si=m->di;m->di=v;
                    ec_cmp(m,(uint8_t)m->ax,(uint8_t)m->dx,8);if (!(m->flags&(1|0x40))) gw_al(m,(uint8_t)m->dx);
                    gw_ah(m,100);ec_cmp(m,(uint8_t)m->ax,0x80,8);if (m->flags&1) gw_ah(m,125);
                    gw_al(m,(uint8_t)m->ax&0x7f);ec_logic(m,(uint8_t)m->ax,8);gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->ax,8,0,1));gw_al(m,200);
                    m->ax=(uint16_t)((uint8_t)m->ax*(m->ax>>8));m->dx=m->ax;if (KI_WORLD_MUTATION==5) m->dx=(uint16_t)(m->dx+1);
                    m->ax=m->di;m->ax=ec_shl(m,m->ax);m->ax=ec_shl(m,m->ax);gw_al(m,5);gw_bl(m,0xff);gw_call(m,0x2fbf,0x57f0,hooks);
                }
            }
        }
        m->di=ec_math(m,m->di,64,16,0,0);m->cx=(uint16_t)((m->cx&0xff00u)|st_dec(m,(uint8_t)m->cx,8));
    } while((uint8_t)m->cx);
    m->di=ec_pop(m);m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void gw_disaster(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->si);m->ds=ec_word(m,m->cs,0xd52);m->si=0x840;m->cx=192;
    do {
        int fire=0;gw_call(m,0xece0,0x229a,hooks);ec_cmp(m,(uint8_t)m->ax,KI_WORLD_MUTATION==6?25:24,8);
        if (m->flags&1) {gw_call(m,0xece0,0x22a1,hooks);gw_al(m,(uint8_t)m->ax&0x3f);ec_logic(m,(uint8_t)m->ax,8);ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->si+0x11)),8);if (!(m->flags&1)) {m->dx=m->si;m->ax=0x10c;gw_bl(m,0xff);gw_call(m,0x2fbf,0x22b2,hooks);fire=1;}}
        if (!fire) {gw_call(m,0xece0,0x22b7,hooks);ec_cmp(m,(uint8_t)m->ax,24,8);if (m->flags&1) {gw_call(m,0xece0,0x22be,hooks);gw_al(m,(uint8_t)m->ax&0x3f);ec_logic(m,(uint8_t)m->ax,8);ec_cmp(m,(uint8_t)m->ax,ec_byte(m,m->ds,(uint16_t)(m->si+0x10)),8);if (!(m->flags&1)) {m->dx=m->si;m->ax=0x20c;gw_bl(m,0xff);gw_call(m,0x2fbf,0x22cf,hooks);}}}
        m->si=ec_math(m,m->si,32,16,0,0);m->cx=(uint16_t)(m->cx-1);
    } while(m->cx);
    m->si=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
static void gw_storm_marker(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);ec_push(m,m->di);
    m->dx=ec_word(m,m->cs,0xd26);m->cx=ec_word(m,m->cs,0xd22);m->dx=ec_math(m,m->dx,m->cx,16,0,1);m->dx=st_shr(m,m->dx,16);m->dx=ec_math(m,m->dx,m->cx,16,0,0);
    m->bx=ec_word(m,m->cs,0xd28);m->cx=ec_word(m,m->cs,0xd24);m->bx=ec_math(m,m->bx,m->cx,16,0,1);m->bx=st_shr(m,m->bx,16);m->bx=ec_math(m,m->bx,m->cx,16,0,0);m->di=0x840;m->cx=192;
    do {
        ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);m->dx=ec_math(m,m->dx,ec_word(m,m->ds,(uint16_t)(m->di+8)),16,0,1);if (m->flags&1) m->dx=ec_math(m,0,m->dx,16,0,1);
        m->bx=ec_math(m,m->bx,ec_word(m,m->ds,(uint16_t)(m->di+10)),16,0,1);if (m->flags&1) m->bx=ec_math(m,0,m->bx,16,0,1);ec_cmp(m,m->dx,m->bx,16);if (m->flags&1) m->dx=m->bx;ec_cmp(m,m->dx,KI_WORLD_MUTATION==7?19:20,16);
        if (m->flags&(1|0x40)) {
            gw_dl(m,(uint8_t)st_shr(m,(uint8_t)m->dx,8));gw_ah(m,(uint8_t)m->ax);gw_ah(m,(uint8_t)ec_math(m,m->ax>>8,(uint8_t)m->dx,8,0,1));if (m->flags&1) {gw_ah(m,0);ec_logic(m,0,8);}ec_put(m,m->ds,(uint16_t)(m->di+0x15),(uint8_t)(m->ax>>8));ec_logic(m,m->ax>>8,8);
            if (m->ax>>8) {gw_ah(m,ec_byte(m,m->cs,0xcff));ec_cmp(m,m->ax>>8,ec_byte(m,m->ds,(uint16_t)(m->di+1)),8);if (m->flags&0x40) {ec_push(m,m->ax);ec_push(m,m->di);m->di=m->sp;gw_call(m,0xce7,0x23e7,hooks);gw_al(m,0x93);m->cx=0x46;gw_call(m,0x8810,0x23ef,hooks);m->di=ec_pop(m);m->ax=ec_pop(m);}}
        }
        m->di=ec_math(m,m->di,32,16,0,0);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->cx=(uint16_t)(m->cx-1);
    } while(m->cx);
    m->di=ec_pop(m);m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);
}
static void gw_storm(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ds);ec_push(m,m->ax);ec_push(m,m->bx);m->ds=ec_word(m,m->cs,0xd52);ec_cmp(m,ec_word(m,m->cs,0xd22),0xfff0,16);
    if (!(m->flags&0x40)) {gw_al(m,0);ec_logic(m,0,8);gw_call(m,0x237e,0x22f0,hooks);ec_store(m,m->cs,0xd22,0xfff0);ec_store(m,m->cs,0xd24,0xfff0);ec_store(m,m->cs,0xd26,400);ec_store(m,m->cs,0xd28,400);}
    gw_call(m,0xece0,0x230f,hooks);ec_logic(m,(uint8_t)m->ax&1,8);
    if (!(m->flags&0x40)) {
        gw_call(m,0xece0,0x2316,hooks);ec_cmp(m,(uint8_t)m->ax,192,8);
        if (m->flags&1) {
            m->bx=(uint16_t)((uint16_t)(uint8_t)m->ax<<8);ec_logic(m,0,8);for(unsigned i=0;i<3;++i) m->bx=st_shr(m,m->bx,16);ec_cmp(m,ec_byte(m,m->ds,(uint16_t)(m->bx+0x848)),192,8);int allowed=1;
            if (m->flags&1) {gw_call(m,0xece0,0x232e,hooks);ec_logic(m,(uint8_t)m->ax&1,8);allowed=!(m->flags&0x40);}
            if (allowed) {
                ec_push(m,m->bx);gw_call(m,0xece0,0x2336,hooks);gw_al(m,(uint8_t)m->ax&7);ec_logic(m,(uint8_t)m->ax,8);gw_al(m,(uint8_t)ec_math(m,(uint8_t)m->ax,8,8,0,0));gw_al(m,gw_shl8(m,(uint8_t)m->ax));gw_al(m,gw_shl8(m,(uint8_t)m->ax));gw_bl(m,(uint8_t)m->ax);m->ax=11;ec_logic(m,0,8);m->dx=0;ec_logic(m,0,16);gw_call(m,0x2fbf,0x2349,hooks);m->bx=ec_pop(m);
                if (!(m->flags&1)) for(unsigned i=0;i<2;++i) {m->ax=ec_word(m,m->ds,(uint16_t)(m->bx+0x848+i*2));ec_cmp(m,m->ax,10,16);if (!ec_less(m)) m->ax=ec_math(m,m->ax,5,16,0,1);ec_store(m,m->cs,(uint16_t)(0xd22+i*2),m->ax);m->ax=ec_math(m,m->ax,10,16,0,0);ec_store(m,m->cs,(uint16_t)(0xd26+i*2),m->ax);}
            }
        }
    }
    m->bx=ec_pop(m);m->ax=ec_pop(m);m->ds=ec_pop(m);
}
int ki_world_body(KiMachine16 *m,uint16_t target,const KiEconomyHooks *hooks) {
    switch(target) {
    case 0x5695:gw_growth(m,hooks);break;case 0x55a6:gw_score(m);break;case 0x2fbf:gw_writer(m,hooks);break;
    case 0x57fe:gw_trust(m,hooks);break;case 0x5715:gw_governor(m,hooks);break;case 0x578f:gw_diplomat(m,hooks);break;
    case 0x30cb:gw_relation(m,hooks);break;case 0x3119:gw_relation_address(m);break;
    case 0x22db:gw_storm(m,hooks);break;case 0x2286:gw_disaster(m,hooks);break;case 0x237e:gw_storm_marker(m,hooks);break;
    default:return 0;
    }return 1;
}
#define GW_WRAPPER(name,target) void name(KiMachine16 *m,const KiEconomyHooks *h) {ki_world_body(m,target,h);m->ip=ec_pop(m);}
GW_WRAPPER(sub_15695,0x5695) GW_WRAPPER(sub_12FBF,0x2fbf) GW_WRAPPER(sub_157FE,0x57fe)
GW_WRAPPER(sub_15715,0x5715) GW_WRAPPER(sub_1578F,0x578f) GW_WRAPPER(sub_130CB,0x30cb)
GW_WRAPPER(sub_122DB,0x22db) GW_WRAPPER(sub_12286,0x2286) GW_WRAPPER(sub_1237E,0x237e)
void sub_155A6(KiMachine16 *m) {gw_score(m);m->ip=ec_pop(m);}
void sub_13119(KiMachine16 *m) {gw_relation_address(m);m->ip=ec_pop(m);}
void ki_world_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *hooks) {
    if (target==0x5695||target==0x55a6||target==0x2fbf||target==0x57fe||target==0x5715||target==0x578f||target==0x30cb||target==0x3119||target==0x22db||target==0x2286||target==0x237e) {
        if(hooks&&hooks->enter) hooks->enter(m,target,hooks->user);ki_world_body(m,target,hooks);m->ip=ec_pop(m);
    } else ki_settlement_invoke(m,target,hooks);
}
