/* 原自動戰鬥與據點易主閉包；非區域返回沿用既有 runner。 */
#include "outcome.h"
#ifndef KI_OUTCOME_MUTATION
#define KI_OUTCOME_MUTATION 0
#endif
static void oc_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=target;ki_outcome_invoke(m,target,h);
}
#include "outcome_generated.inc"
void ki_outcome_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *h) {
 if(oc_known(target)) {
  if(h&&h->enter)h->enter(m,target,h->user);
  int handled=ki_outcome_body(m,target,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_route_invoke(m,target,h);
}
