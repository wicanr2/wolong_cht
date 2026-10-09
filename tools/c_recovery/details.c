/* 原資訊介面閉包；共用規則保留唯一實作。 */
#include "details.h"
#ifndef KI_DETAILS_MUTATION
#define KI_DETAILS_MUTATION 0
#endif
static void dtl_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_details_invoke(m,t,h);
}
#include "details_generated.inc"
void ki_details_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(dtl_known(t)) {
  if(h&&h->enter) h->enter(m,t,h->user);
  int handled=ki_details_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_personnel_invoke(m,t,h);
}
void sub_15E1E(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x5e1e,h);m->ip=ec_pop(m);}
void sub_15E2D(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x5e2d,h);m->ip=ec_pop(m);}
void sub_17E1F(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x7e1f,h);m->ip=ec_pop(m);}
void sub_17E4A(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x7e4a,h);m->ip=ec_pop(m);}
void sub_17F1A(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x7f1a,h);m->ip=ec_pop(m);}
void sub_17F61(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x7f61,h);m->ip=ec_pop(m);}
void sub_17F81(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x7f81,h);m->ip=ec_pop(m);}
void sub_1807B(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x807b,h);m->ip=ec_pop(m);}
void sub_1812A(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x812a,h);m->ip=ec_pop(m);}
void sub_1817D(KiMachine16 *m,const KiEconomyHooks *h) {ki_details_body(m,0x817d,h);m->ip=ec_pop(m);}
