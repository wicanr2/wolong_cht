#ifndef WOLONG_C_EVENTS_H
#define WOLONG_C_EVENTS_H
#include "hourly.h"
int ki_events_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_events_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define EV_DECL(name) void name(KiMachine16 *,const KiEconomyHooks *);
EV_DECL(sub_1320C) EV_DECL(sub_13220) EV_DECL(sub_13262) EV_DECL(sub_132A9)
EV_DECL(sub_132E9) EV_DECL(sub_13327) EV_DECL(sub_13388) EV_DECL(sub_133EA)
EV_DECL(sub_133FD) EV_DECL(sub_13485) EV_DECL(sub_134A6) EV_DECL(sub_134B1)
EV_DECL(sub_1351A) EV_DECL(sub_13526) EV_DECL(sub_135AB) EV_DECL(sub_135ED)
EV_DECL(sub_13639) EV_DECL(sub_13669) EV_DECL(sub_13697) EV_DECL(sub_136C4)
EV_DECL(sub_13712) EV_DECL(sub_13771) EV_DECL(sub_137D8) EV_DECL(sub_137F5)
EV_DECL(sub_13138) EV_DECL(sub_14502) EV_DECL(sub_16A3D) EV_DECL(sub_123FF)
EV_DECL(sub_12438) EV_DECL(sub_150D7)
#undef EV_DECL
#endif
