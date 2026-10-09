/* 原玩家編成閉包；共用規則保留唯一實作。 */
#include "formation.h"
#ifndef KI_FORMATION_MUTATION
#define KI_FORMATION_MUTATION 0
#endif
static void fc_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_formation_invoke(m,t,h);
}
#include "formation_generated.inc"
void ki_formation_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(fc_known(t)) {
  if(h&&h->enter) h->enter(m,t,h->user);
  int handled=ki_formation_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_catalog_invoke(m,t,h);
}
void sub_16C5E(KiMachine16 *m,const KiEconomyHooks *h) {ki_formation_body(m,0x6c5e,h);m->ip=ec_pop(m);}
void sub_16C92(KiMachine16 *m,const KiEconomyHooks *h) {ki_formation_body(m,0x6c92,h);m->ip=ec_pop(m);}
void sub_16D56(KiMachine16 *m,const KiEconomyHooks *h) {ki_formation_body(m,0x6d56,h);m->ip=ec_pop(m);}
void sub_16D6F(KiMachine16 *m,const KiEconomyHooks *h) {ki_formation_body(m,0x6d6f,h);m->ip=ec_pop(m);}
void sub_16DA8(KiMachine16 *m,const KiEconomyHooks *h) {ki_formation_body(m,0x6da8,h);m->ip=ec_pop(m);}
void sub_16DFD(KiMachine16 *m,const KiEconomyHooks *h) {ki_formation_body(m,0x6dfd,h);m->ip=ec_pop(m);}
void sub_16E80(KiMachine16 *m,const KiEconomyHooks *h) {ki_formation_body(m,0x6e80,h);m->ip=ec_pop(m);}
