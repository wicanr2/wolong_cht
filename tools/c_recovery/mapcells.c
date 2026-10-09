/* Original map-cell closure; include after display.c. re/114, spec/234. */
#include "mapcells.h"
#ifndef KI_MAPCELLS_MUTATION
#define KI_MAPCELLS_MUTATION 0
#endif
static void mc_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_mapcells_invoke(m,t,h);
}
static void mc_string(KiMachine16 *m,unsigned operation,unsigned repeat) {
    if(repeat&&!m->cx) return;
    do {
        uint16_t value=operation==2?m->ax:dl_read16(m,m->ds,m->si);
        if(operation==1) m->ax=value;else dl_write16(m,m->es,m->di,value);
        int step=m->flags&0x400?-2:2;
        if(operation!=2) m->si=(uint16_t)(m->si+step);
        if(operation!=1) m->di=(uint16_t)(m->di+step);
    }while(repeat&&--m->cx);
}
#include "mapcells_generated.inc"
void ki_mapcells_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);
    int handled=ki_mapcells_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
}
#define MC_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_mapcells_body(m,t,h);m->ip=ec_pop(m);}
MC_WRAP(sub_1D46A,0xd46a) MC_WRAP(sub_1D483,0xd483) MC_WRAP(sub_1D4C7,0xd4c7)
MC_WRAP(sub_1D615,0xd615) MC_WRAP(sub_1D66A,0xd66a) MC_WRAP(sub_1D76B,0xd76b)
MC_WRAP(sub_1D782,0xd782) MC_WRAP(sub_1D796,0xd796) MC_WRAP(sub_1D7E7,0xd7e7)
MC_WRAP(sub_1D804,0xd804)
#undef MC_WRAP
