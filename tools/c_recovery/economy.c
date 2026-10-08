/* DOS/V KI.EXE: six original entries; contract docs/spec/205.
 * Structured native C with the original 16-bit ABI; not a matching C binary.
 */
#include "economy.h"
#ifndef KI_ECONOMY_MUTATION
#define KI_ECONOMY_MUTATION 0
#endif
static uint32_t ec_at(uint16_t segment, uint16_t offset) {
    return (((uint32_t)segment << 4) + offset) & 0xfffffu;
}
static uint8_t ec_byte(KiMachine16 *m, uint16_t seg, uint16_t off) { return m->memory[ec_at(seg,off)]; }
static void ec_put(KiMachine16 *m,uint16_t seg,uint16_t off,uint8_t v) { m->memory[ec_at(seg,off)]=v; }
static uint16_t ec_word(KiMachine16 *m,uint16_t seg,uint16_t off) {
    uint32_t at=ec_at(seg,off);
    return (uint16_t)(m->memory[at]|(uint16_t)m->memory[(at+1)&0xfffffu]<<8);
}
static void ec_store(KiMachine16 *m,uint16_t seg,uint16_t off,uint16_t v) {
    uint32_t at=ec_at(seg,off);m->memory[at]=(uint8_t)v;m->memory[(at+1)&0xfffffu]=(uint8_t)(v>>8);
}
static void ec_push(KiMachine16 *m,uint16_t v) { m->sp=(uint16_t)(m->sp-2);ec_store(m,m->ss,m->sp,v); }
static uint16_t ec_pop(KiMachine16 *m) { uint16_t v=ec_word(m,m->ss,m->sp);m->sp=(uint16_t)(m->sp+2);return v; }
static uint16_t ec_szp(uint16_t flags,uint16_t r,unsigned width) {
    uint8_t p=(uint8_t)r;p^=p>>4;p^=p>>2;p^=p>>1;flags&=(uint16_t)~0xc4u;
    if (!(p&1)) flags|=4;if (!r) flags|=0x40;if (r&(width==8?0x80:0x8000)) flags|=0x80;return flags;
}
static uint16_t ec_math(KiMachine16 *m,uint16_t a,uint16_t b,unsigned width,unsigned carry,int subtract) {
    uint32_t mask=width==8?0xff:0xffff,sign=width==8?0x80:0x8000;
    uint32_t full=subtract?(uint32_t)a-b-carry:(uint32_t)a+b+carry;
    uint16_t r=(uint16_t)(full&mask),bits=0;
    if (subtract?(uint32_t)a<(uint32_t)b+carry:full>mask) bits|=1;
    if ((a^b^r)&0x10) bits|=0x10;
    if ((subtract?((a^b)&(a^r)):(~(a^b)&(a^r)))&sign) bits|=0x800;
    m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|bits),r,width);return r;
}
static void ec_cmp(KiMachine16 *m,uint16_t a,uint16_t b,unsigned width) { (void)ec_math(m,a,b,width,0,1); }
static uint16_t ec_inc(KiMachine16 *m,uint16_t a,unsigned width) {
    uint16_t carry=m->flags&1,r=ec_math(m,a,1,width,0,0);m->flags=(uint16_t)((m->flags&~1u)|carry);return r;
}
static void ec_logic(KiMachine16 *m,uint16_t r,unsigned width) { m->flags=ec_szp((uint16_t)(m->flags&~0x8d5u),r,width); }
static uint16_t ec_shl(KiMachine16 *m,uint16_t a) {
    uint16_t r=(uint16_t)(a<<1),bits=(uint16_t)(a>>15);
    if (((r>>15)^(a>>15))&1) bits|=0x800;
    m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|bits),r,16);return r;
}
static int ec_less(KiMachine16 *m) { return ((m->flags>>7)^(m->flags>>11))&1; }
static int ec_greater(KiMachine16 *m) { return !(m->flags&0x40)&&!ec_less(m); }
static void ec_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *hooks) {
    ec_push(m,next);m->ip=target;ki_economy_invoke(m,target,hooks);
}

void sub_15609(KiMachine16 *m) {
    ec_push(m,m->ax);ec_push(m,m->dx);
    m->ax=ec_math(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->si+0x20)),16,0,0);
    uint8_t dl=(uint8_t)ec_math(m,(uint8_t)m->dx,ec_byte(m,m->ds,(uint16_t)(m->si+0x22)),8,m->flags&1,0);
    m->dx=(uint16_t)((m->dx&0xff00u)|dl);ec_cmp(m,dl,9,8);
    if (!ec_less(m)) {
        int above=ec_greater(m);dl=9;
        if (!above) { ec_cmp(m,m->ax,0xfe98,16);above=!(m->flags&1); }
        if (above) m->ax=KI_ECONOMY_MUTATION==1?0xfe97:0xfe98;
    }
    ec_store(m,m->ds,(uint16_t)(m->si+0x20),m->ax);ec_put(m,m->ds,(uint16_t)(m->si+0x22),dl);
    m->dx=ec_pop(m);m->ax=ec_pop(m);m->ip=ec_pop(m);
}
void sub_1563B(KiMachine16 *m) {
    ec_push(m,m->ax);ec_push(m,m->dx);
    ec_store(m,m->ds,(uint16_t)(m->si+0x20),ec_math(m,ec_word(m,m->ds,(uint16_t)(m->si+0x20)),m->ax,16,0,1));
    ec_put(m,m->ds,(uint16_t)(m->si+0x22),(uint8_t)ec_math(m,ec_byte(m,m->ds,(uint16_t)(m->si+0x22)),(uint8_t)m->dx,8,m->flags&1,1));
    m->ax=ec_word(m,m->ds,(uint16_t)(m->si+0x20));uint8_t dl=ec_byte(m,m->ds,(uint16_t)(m->si+0x22));
    m->dx=(uint16_t)((m->dx&0xff00u)|dl);ec_cmp(m,dl,0xf6,8);
    if (!ec_greater(m)) {
        int below=ec_less(m);dl=0xf6;
        if (!below) { ec_cmp(m,m->ax,0x168,16);below=(m->flags&(1|0x40))!=0; }
        if (below) m->ax=KI_ECONOMY_MUTATION==2?0x169:0x168;
    }
    ec_store(m,m->ds,(uint16_t)(m->si+0x20),m->ax);ec_put(m,m->ds,(uint16_t)(m->si+0x22),dl);
    m->dx=ec_pop(m);m->ax=ec_pop(m);m->ip=ec_pop(m);
}
void sub_155EC(KiMachine16 *m) {
    m->ax=ec_math(m,m->ax,m->dx,16,0,0);
    if (m->flags&1) { if (KI_ECONOMY_MUTATION!=3) m->ax=0xffdc; }
    else { ec_cmp(m,m->ax,0xffdc,16);if (!(m->flags&(1|0x40))) m->ax=0xffdc; }
    m->ip=ec_pop(m);
}
void sub_154FC(KiMachine16 *m) {
    m->ax=ec_word(m,m->ds,(uint16_t)(m->si+8));
    m->ax=ec_math(m,m->ax,ec_word(m,m->ds,(uint16_t)(m->di+8)),16,0,1);
    if (m->flags&1) m->ax=ec_math(m,0,m->ax,16,0,1);
    m->bx=ec_word(m,m->ds,(uint16_t)(m->si+10));
    m->bx=ec_math(m,m->bx,ec_word(m,m->ds,(uint16_t)(m->di+10)),16,0,1);
    if (m->flags&1) m->bx=ec_math(m,0,m->bx,16,0,1);
    ec_cmp(m,m->ax,m->bx,16);if (m->flags&1) m->ax=m->bx;
    ec_logic(m,m->ax>>8,8);if (m->ax>>8) m->ax=(uint16_t)((m->ax&0xff00u)|0xff);
    m->bx=0;ec_logic(m,0,16);
    for (;;) {
        m->ax=(uint16_t)(((uint16_t)ec_byte(m,m->cs,(uint16_t)(m->bx+0x5535))<<8)|(m->ax&0xff));
        uint8_t threshold=ec_byte(m,m->cs,(uint16_t)(m->bx+0x5532));
        if (KI_ECONOMY_MUTATION==4&&m->bx==0) threshold=79;
        ec_cmp(m,(uint8_t)m->ax,threshold,8);if (m->flags&(1|0x40)) break;
        m->bx=ec_inc(m,m->bx,16);
    }
    m->bx=m->ax>>8;ec_logic(m,0,8);m->ip=ec_pop(m);
}
void sub_15828(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ax);ec_push(m,m->bx);ec_push(m,m->cx);ec_push(m,m->dx);
    m->dx=ec_word(m,m->ds,(uint16_t)(m->si+0x21));ec_cmp(m,m->dx>>8,0x80,8);
    if (!(m->flags&1)) {
        m->bx=m->si;m->bx=ec_math(m,m->bx,4,16,0,0);m->dx=ec_math(m,0,m->dx,16,0,1);
        for (unsigned i=0;i<(KI_ECONOMY_MUTATION==5?3:4);++i) m->dx=ec_shl(m,m->dx);
        m->cx=3;
        do {
            ec_call(m,0xece0,0x5849,hooks);m->ax&=0x1f;ec_logic(m,m->ax,16);
            m->ax=ec_math(m,m->ax,m->dx,16,0,0);
            uint16_t left=ec_math(m,ec_word(m,m->ds,m->bx),m->ax,16,0,1);ec_store(m,m->ds,m->bx,left);
            if (m->flags&1) ec_store(m,m->ds,m->bx,0);
            m->bx=ec_inc(m,m->bx,16);m->bx=ec_inc(m,m->bx,16);m->cx=(uint16_t)(m->cx-1);
        } while (m->cx);
    }
    m->dx=ec_pop(m);m->cx=ec_pop(m);m->bx=ec_pop(m);m->ax=ec_pop(m);m->ip=ec_pop(m);
}
void sub_15358(KiMachine16 *m,const KiEconomyHooks *hooks) {
    ec_push(m,m->ax);m->ds=ec_word(m,m->ds,0xd52);m->si=0;m->cx&=0xff00;ec_logic(m,0,8);
    do {
        ec_cmp(m,ec_byte(m,m->ds,m->si),0x80,8);
        if (!(m->flags&1)) {
            m->ax=ec_word(m,m->ds,(uint16_t)(m->si+0x1a));m->dx=(uint16_t)((m->dx&0xff00u)|ec_byte(m,m->ds,(uint16_t)(m->si+0x1c)));
            ec_call(m,0x563b,0x5370,hooks);
            if (KI_ECONOMY_MUTATION==6) { ec_call(m,0x5609,0x5376,hooks);ec_call(m,0x53c6,0x5373,hooks); }
            else { ec_call(m,0x53c6,0x5373,hooks);ec_call(m,0x5609,0x5376,hooks); }
            m->ax=0;ec_logic(m,0,16);ec_store(m,m->ds,(uint16_t)(m->si+0x1a),0);ec_put(m,m->ds,(uint16_t)(m->si+0x1c),0);
            ec_call(m,0x5828,0x5381,hooks);
        }
        m->si=ec_math(m,m->si,0x40,16,0,0);m->cx=(uint16_t)((m->cx&0xff00u)|ec_inc(m,(uint8_t)m->cx,8));
        ec_cmp(m,(uint8_t)m->cx,0x16,8);
    } while (m->flags&1);
    static const uint16_t targets[]={0x5695,0x585f,0x55a6,0x2bd9,0x5715,0x578f,0x22db,0x2286,0x57fe};
    for (unsigned i=0;i<9;++i) ec_call(m,targets[i],(uint16_t)(0x538e + i*3),hooks);
    m->si=0xd10;m->di=0xd08;m->cx=4;
    do {
        m->ax=ec_word(m,m->cs,m->si);ec_store(m,m->cs,m->di,m->ax);
        m->si=ec_inc(m,m->si,16);m->si=ec_inc(m,m->si,16);m->di=ec_inc(m,m->di,16);m->di=ec_inc(m,m->di,16);
        m->cx=(uint16_t)(m->cx-1);
    } while (m->cx);
    m->ax=(uint16_t)((m->ax&0xff00u)|0xe);ec_call(m,0x5e80,0x53c0,hooks);
    m->ax=m->cs;m->ds=m->ax;m->ax=ec_pop(m);m->ip=ec_pop(m);
}
void ki_economy_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *hooks) {
    if (hooks&&hooks->enter) hooks->enter(m,target,hooks->user);
    switch (target) {
    case 0x5358:sub_15358(m,hooks);break;
    case 0x5609:sub_15609(m);break;
    case 0x563b:sub_1563B(m);break;
    case 0x54fc:sub_154FC(m);break;
    case 0x55ec:sub_155EC(m);break;
    case 0x5828:sub_15828(m,hooks);break;
    case 0xece0:sub_1ECE0_abi(m);break;
    default:if (hooks&&hooks->external) hooks->external(m,target,hooks->user);m->ip=ec_pop(m);break;
    }
}
