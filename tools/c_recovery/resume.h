#ifndef WOLONG_C_RESUME_H
#define WOLONG_C_RESUME_H
#include "overlay.h"
int ki_resume_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_resume_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define SR_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
SR_DECL(sub_11F30) SR_DECL(sub_11F5A) SR_DECL(sub_13B08) SR_DECL(sub_15C58)
SR_DECL(sub_19541) SR_DECL(sub_195B1) SR_DECL(sub_196ED) SR_DECL(sub_19752) SR_DECL(sub_1D4A3)
#undef SR_DECL
#endif
