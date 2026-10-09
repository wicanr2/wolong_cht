/* Original world overlay producers; include after mapcells.c. re/115, spec/235. */
#include "overlay.h"
#ifndef KI_OVERLAY_MUTATION
#define KI_OVERLAY_MUTATION 0
#endif
static void ov_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
    ec_push(m,next);m->ip=t;ki_overlay_invoke(m,t,h);
}
#include "overlay_generated.inc"
void ki_overlay_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);
    int handled=ki_overlay_body(m,t,h)||ki_mapcells_body(m,t,h)||ki_rect_body(m,t,h);
    assert(handled);(void)handled;m->ip=ec_pop(m);
}
#define OV_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_overlay_body(m,t,h);m->ip=ec_pop(m);}
OV_WRAP(sub_11CC9,0x1cc9) OV_WRAP(sub_12533,0x2533) OV_WRAP(sub_12AF4,0x2af4)
OV_WRAP(sub_12B2A,0x2b2a) OV_WRAP(sub_12B3C,0x2b3c) OV_WRAP(sub_15D19,0x5d19)
#undef OV_WRAP
