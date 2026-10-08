#ifndef WOLONG_C_HOTSPOT_H
#define WOLONG_C_HOTSPOT_H
#include "numeric.h"
int ki_hotspot_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_hotspot_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define HS_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
HS_DECL(sub_1E3C0) HS_DECL(sub_1E3D7) HS_DECL(sub_1E41B) HS_DECL(sub_1E453)
HS_DECL(sub_1895D) HS_DECL(sub_189DE) HS_DECL(sub_1D5D4)
HS_DECL(sub_10C14) HS_DECL(sub_10C60) HS_DECL(sub_10C77)
#undef HS_DECL
#endif
