/* 原世界互動閉包；既有遊戲與平台常式保留唯一實作。 */
#include "interaction.h"
#ifndef KI_INTERACTION_MUTATION
#define KI_INTERACTION_MUTATION 0
#endif
static void ix_call(KiMachine16 *m,uint16_t target,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=target;ki_interaction_invoke(m,target,h);
}
#include "interaction_generated.inc"
void ki_interaction_invoke(KiMachine16 *m,uint16_t target,const KiEconomyHooks *h) {
 if(ix_known(target)) {
  if(h&&h->enter)h->enter(m,target,h->user);
  int handled=ki_interaction_body(m,target,h);assert(handled);(void)handled;
  m->ip=ec_pop(m);
 } else ki_bootstrap_invoke(m,target,h);
}
