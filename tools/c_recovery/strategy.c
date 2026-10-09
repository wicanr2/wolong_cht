/* 原政略指令閉包；共用規則保留唯一實作。 */
#include "strategy.h"
#ifndef KI_STRATEGY_MUTATION
#define KI_STRATEGY_MUTATION 0
#endif
static void sty_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_strategy_invoke(m,t,h);
}
#include "strategy_generated.inc"
void ki_strategy_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(sty_known(t)) {
  if(h&&h->enter)h->enter(m,t,h->user);
  int handled=ki_strategy_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_march_invoke(m,t,h);
}
