/* Original corps/general/faction/new-game list callers; include after list.c. */
#include "catalog.h"
#ifndef KI_CATALOG_MUTATION
#define KI_CATALOG_MUTATION 0
#endif
static void cc_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_catalog_invoke(m,t,h);
}
#include "catalog_generated.inc"
void ki_catalog_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(cc_known(t)) {
  if(h&&h->enter) h->enter(m,t,h->user);
  int handled=ki_catalog_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_list_invoke(m,t,h);
}

#define CC_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_catalog_body(m,t,h);m->ip=ec_pop(m);}
CC_WRAP(sub_1716D,0x716d)
CC_WRAP(sub_171D3,0x71d3)
CC_WRAP(sub_175FA,0x75fa)
CC_WRAP(sub_17663,0x7663)
CC_WRAP(sub_1770C,0x770c)
CC_WRAP(sub_178A7,0x78a7)
CC_WRAP(sub_17906,0x7906)
CC_WRAP(sub_1799C,0x799c)
CC_WRAP(sub_17A7A,0x7a7a)
CC_WRAP(sub_17B3C,0x7b3c)
CC_WRAP(sub_17BC0,0x7bc0)
#undef CC_WRAP
