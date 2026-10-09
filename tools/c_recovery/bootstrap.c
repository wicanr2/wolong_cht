#include "bootstrap.h"
#ifndef KI_BOOTSTRAP_MUTATION
#define KI_BOOTSTRAP_MUTATION 0
#endif
static void bs_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_bootstrap_invoke(m,t,h);
}
#include "bootstrap_generated.inc"
void ki_bootstrap_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(bs_known(t)) {
  if(h&&h->enter)h->enter(m,t,h->user);
  int handled=ki_bootstrap_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_main_invoke(m,t,h);
}
