/* 原非區域退出閉包；共用遊戲常式保留唯一實作。 */
#include "main.h"
#ifndef KI_MAIN_MUTATION
#define KI_MAIN_MUTATION 0
#endif
static void nle_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_main_invoke(m,t,h);
}
#include "main_generated.inc"
void ki_main_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(nle_known(t)) {
  if(h&&h->enter)h->enter(m,t,h->user);
  int handled=ki_main_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_strategy_invoke(m,t,h);
}
