#ifndef WOLONG_C_MAIN_H
#define WOLONG_C_MAIN_H
#include "strategy.h"
uint8_t wolong_main_in8(uint16_t);
/* Native ABI bridge, invoked only after the original virtual RET completed. */
_Noreturn void wolong_main_transfer(KiMachine16 *,void *);
int ki_main_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_main_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void sub_11CB1(KiMachine16 *,const KiEconomyHooks *);
void sub_10A1C(KiMachine16 *,const KiEconomyHooks *);
void sub_1EBDC(KiMachine16 *,const KiEconomyHooks *);
#endif
