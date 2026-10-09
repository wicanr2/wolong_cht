#ifndef WOLONG_C_MAPCELLS_H
#define WOLONG_C_MAPCELLS_H
#include "display.h"
int ki_mapcells_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_mapcells_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define MC_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
MC_DECL(sub_1D46A) MC_DECL(sub_1D483) MC_DECL(sub_1D4C7) MC_DECL(sub_1D615)
MC_DECL(sub_1D66A) MC_DECL(sub_1D76B) MC_DECL(sub_1D782) MC_DECL(sub_1D796)
MC_DECL(sub_1D7E7) MC_DECL(sub_1D804)
#undef MC_DECL
#endif
