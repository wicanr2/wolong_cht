#ifndef WOLONG_C_RECOVERY_SETTLEMENT_H
#define WOLONG_C_RECOVERY_SETTLEMENT_H
#include "economy.h"
void sub_153C6(KiMachine16 *,const KiEconomyHooks *);
void sub_15456(KiMachine16 *);
void sub_1548F(KiMachine16 *);
void sub_15538(KiMachine16 *);
void sub_15547(KiMachine16 *);
void ki_settlement_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
int ki_settlement_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#endif
