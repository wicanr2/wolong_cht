#ifndef WOLONG_C_RECOVERY_HOURLY_H
#define WOLONG_C_RECOVERY_HOURLY_H
#include "politics.h"
int ki_hourly_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_hourly_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define HR_DECL(name) void name(KiMachine16 *,const KiEconomyHooks *);
HR_DECL(sub_13E11) HR_DECL(sub_13E65) HR_DECL(sub_13E8E) HR_DECL(sub_15673)
HR_DECL(sub_131AE) HR_DECL(sub_13496) HR_DECL(sub_13507) HR_DECL(sub_13DC9)
#undef HR_DECL
#endif
