#ifndef WOLONG_C_VERDICT_H
#define WOLONG_C_VERDICT_H
#include "resume.h"
int ki_verdict_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_verdict_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define VD_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
VD_DECL(sub_1461D) VD_DECL(sub_14698) VD_DECL(sub_14717) VD_DECL(sub_15E80)
VD_DECL(sub_15EB7) VD_DECL(sub_15F27) VD_DECL(sub_15F5D) VD_DECL(sub_15F7F)
VD_DECL(sub_1699E) VD_DECL(sub_16E8F) VD_DECL(sub_16EC9) VD_DECL(sub_16F26)
VD_DECL(sub_16F86) VD_DECL(sub_16FD2)
#undef VD_DECL
#endif
