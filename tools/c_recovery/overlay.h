#ifndef WOLONG_C_OVERLAY_H
#define WOLONG_C_OVERLAY_H
#include "mapcells.h"
int ki_overlay_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_overlay_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define OV_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
OV_DECL(sub_11CC9) OV_DECL(sub_12533) OV_DECL(sub_12AF4) OV_DECL(sub_12B2A)
OV_DECL(sub_12B3C) OV_DECL(sub_15D19)
#undef OV_DECL
#endif
