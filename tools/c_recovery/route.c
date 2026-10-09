/* 原尋路、退卻與潰散閉包；live CS operands 由原 patch 寫入後消費。 */
#include "route.h"
#ifndef KI_ROUTE_MUTATION
#define KI_ROUTE_MUTATION 0
#endif
static void rt_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=target;ki_route_invoke(m,target,h);
}
#include "route_generated.inc"
void ki_route_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *h) {
 if(rt_known(target)) {
  if(h&&h->enter)h->enter(m,target,h->user);
  int handled=ki_route_body(m,target,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_tick_invoke(m,target,h);
}
