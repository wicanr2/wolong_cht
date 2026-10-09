/* 原據點與動畫閉包；既有規則及平台服務保留唯一實作。 */
#include "tick.h"
#ifndef KI_TICK_MUTATION
#define KI_TICK_MUTATION 0
#endif
static void ct_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=target;ki_tick_invoke(m,target,h);
}
#include "tick_generated.inc"
void ki_tick_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *h) {
 if(ct_known(target)) {
  if(h&&h->enter)h->enter(m,target,h->user);
  int handled=ki_tick_body(m,target,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_interaction_invoke(m,target,h);
}
