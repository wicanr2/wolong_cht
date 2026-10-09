#ifndef WOLONG_C_LIST_H
#define WOLONG_C_LIST_H
#include "verdict.h"
int ki_list_body(KiMachine16 *,uint16_t,const KiEconomyHooks *,int *);
void ki_list_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define LL_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
LL_DECL(sub_15E60)
LL_DECL(sub_16909)
LL_DECL(sub_17400)
LL_DECL(sub_1748F)
LL_DECL(sub_181C0)
LL_DECL(sub_1820E)
LL_DECL(sub_18412)
LL_DECL(nullsub_2)
LL_DECL(sub_18463)
LL_DECL(sub_184BC)
LL_DECL(sub_184DD)
LL_DECL(sub_1851A)
LL_DECL(sub_18546)
LL_DECL(sub_18607)
LL_DECL(sub_18662)
LL_DECL(sub_186C9)
LL_DECL(sub_18713)
LL_DECL(sub_18755)
#undef LL_DECL
#endif
