#ifndef WOLONG_C_BOOTSTRAP_H
#define WOLONG_C_BOOTSTRAP_H
#include "main.h"
int ki_bootstrap_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_bootstrap_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void raw_entry_11BE0_prefix(KiMachine16 *,const KiEconomyHooks *);
void sub_1533D(KiMachine16 *,const KiEconomyHooks *);
void sub_189F0(KiMachine16 *,const KiEconomyHooks *);
void sub_18A1E(KiMachine16 *,const KiEconomyHooks *);
void sub_18AD1(KiMachine16 *,const KiEconomyHooks *);
void sub_18AEA(KiMachine16 *,const KiEconomyHooks *);
#endif
