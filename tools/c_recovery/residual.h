#ifndef WOLONG_C_RESIDUAL_H
#define WOLONG_C_RESIDUAL_H
#include "startup.h"
uint8_t wolong_residual_in8(uint16_t);
_Noreturn void wolong_residual_exit(KiMachine16 *,uint8_t,void *);
int ki_residual_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_residual_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void sub_145F8(KiMachine16 *,const KiEconomyHooks *);
void sub_15AD1(KiMachine16 *,const KiEconomyHooks *);
void sub_15AFC(KiMachine16 *,const KiEconomyHooks *);
void sub_15C14(KiMachine16 *,const KiEconomyHooks *);
void sub_19377(KiMachine16 *,const KiEconomyHooks *);
void sub_193D7(KiMachine16 *,const KiEconomyHooks *);
void sub_1AB7C(KiMachine16 *,const KiEconomyHooks *);
void sub_1C300(KiMachine16 *,const KiEconomyHooks *);
void sub_1E567(KiMachine16 *,const KiEconomyHooks *);
void sub_1E6FF(KiMachine16 *,const KiEconomyHooks *);
void raw_entry_1AB48(KiMachine16 *,const KiEconomyHooks *);
#endif
