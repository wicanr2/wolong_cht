/* 原軍團行軍閉包；共用規則保留唯一實作。 */
#include "march.h"
#ifndef KI_MARCH_MUTATION
#define KI_MARCH_MUTATION 0
#endif
static void mch_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_march_invoke(m,t,h);
}
#include "march_generated.inc"
void ki_march_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(t==0xece0) {
  if(h&&h->enter) h->enter(m,t,h->user);
  sub_1ECE0_abi(m);return;
 }
 if(mch_known(t)) {
  if(h&&h->enter) h->enter(m,t,h->user);
  int handled=ki_march_body(m,t,h);assert(handled);(void)handled;m->ip=ec_pop(m);
 } else ki_details_invoke(m,t,h);
}
void sub_11C8D(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x1c8d,h);m->ip=ec_pop(m);}
void sub_11F7F(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x1f7f,h);m->ip=ec_pop(m);}
void sub_12151(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x2151,h);m->ip=ec_pop(m);}
void sub_121B2(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x21b2,h);m->ip=ec_pop(m);}
void sub_14325(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x4325,h);m->ip=ec_pop(m);}
void sub_14370(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x4370,h);m->ip=ec_pop(m);}
void sub_1439D(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x439d,h);m->ip=ec_pop(m);}
void sub_143AF(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x43af,h);m->ip=ec_pop(m);}
void sub_1440F(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x440f,h);m->ip=ec_pop(m);}
void sub_14466(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x4466,h);m->ip=ec_pop(m);}
void sub_14483(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x4483,h);m->ip=ec_pop(m);}
void sub_14499(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x4499,h);m->ip=ec_pop(m);}
void sub_144A9(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x44a9,h);m->ip=ec_pop(m);}
void sub_144D6(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x44d6,h);m->ip=ec_pop(m);}
void sub_14548(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x4548,h);m->ip=ec_pop(m);}
void sub_1463E(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x463e,h);m->ip=ec_pop(m);}
void sub_14651(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x4651,h);m->ip=ec_pop(m);}
void sub_14689(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x4689,h);m->ip=ec_pop(m);}
void sub_159A6(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x59a6,h);m->ip=ec_pop(m);}
void sub_15AB6(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x5ab6,h);m->ip=ec_pop(m);}
void sub_1703C(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x703c,h);m->ip=ec_pop(m);}
void sub_1709E(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x709e,h);m->ip=ec_pop(m);}
void sub_17F90(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x7f90,h);m->ip=ec_pop(m);}
void sub_17FDB(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x7fdb,h);m->ip=ec_pop(m);}
void sub_1804E(KiMachine16 *m,const KiEconomyHooks *h) {ki_march_body(m,0x804e,h);m->ip=ec_pop(m);}
