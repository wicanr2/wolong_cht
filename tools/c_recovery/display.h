#ifndef WOLONG_C_DISPLAY_H
#define WOLONG_C_DISPLAY_H
#include "rect.h"
int ki_display_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_display_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define DL_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
DL_DECL(sub_1030F) DL_DECL(sub_10337) DL_DECL(sub_1E993) DL_DECL(sub_1E9A7) DL_DECL(sub_1E9C1)
DL_DECL(sub_1FB29) DL_DECL(sub_1FBA7) DL_DECL(sub_1F6DC) DL_DECL(sub_1F878)
#undef DL_DECL
#endif
