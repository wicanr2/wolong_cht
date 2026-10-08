#ifndef WOLONG_C_NUMERIC_H
#define WOLONG_C_NUMERIC_H
#include "modal.h"
int ki_numeric_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_numeric_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define NU_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
NU_DECL(sub_17C6E) NU_DECL(sub_17D0D) NU_DECL(sub_17D47) NU_DECL(sub_17D5F)
NU_DECL(sub_17DA5) NU_DECL(sub_17DC3) NU_DECL(sub_17DDD) NU_DECL(sub_17DEA)
NU_DECL(sub_17DEC) NU_DECL(sub_17DF1) NU_DECL(sub_101B4) NU_DECL(sub_101DB)
NU_DECL(sub_193E9) NU_DECL(sub_167CD) NU_DECL(sub_167E6) NU_DECL(sub_16806) NU_DECL(sub_16826)
#undef NU_DECL
#endif
