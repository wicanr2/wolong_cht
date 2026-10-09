/* 原人事任免閉包；共用規則保留唯一實作。 */
#include "personnel.h"
#ifndef KI_PERSONNEL_MUTATION
#define KI_PERSONNEL_MUTATION 0
#endif
static void psn_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_personnel_invoke(m,t,h);
}
#include "personnel_generated.inc"
void ki_personnel_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(psn_known(t)) {
  if(h&&h->enter) h->enter(m,t,h->user);
  int handled=ki_personnel_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_formation_invoke(m,t,h);
}
void sub_16265(KiMachine16 *m,const KiEconomyHooks *h) {ki_personnel_body(m,0x6265,h);m->ip=ec_pop(m);}
void sub_16A9B(KiMachine16 *m,const KiEconomyHooks *h) {ki_personnel_body(m,0x6a9b,h);m->ip=ec_pop(m);}
void sub_16B08(KiMachine16 *m,const KiEconomyHooks *h) {ki_personnel_body(m,0x6b08,h);m->ip=ec_pop(m);}
void sub_16B4F(KiMachine16 *m,const KiEconomyHooks *h) {ki_personnel_body(m,0x6b4f,h);m->ip=ec_pop(m);}
void sub_16B71(KiMachine16 *m,const KiEconomyHooks *h) {ki_personnel_body(m,0x6b71,h);m->ip=ec_pop(m);}
void sub_16BE3(KiMachine16 *m,const KiEconomyHooks *h) {ki_personnel_body(m,0x6be3,h);m->ip=ec_pop(m);}
void sub_16C2A(KiMachine16 *m,const KiEconomyHooks *h) {ki_personnel_body(m,0x6c2a,h);m->ip=ec_pop(m);}
