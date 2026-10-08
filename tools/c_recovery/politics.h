#ifndef WOLONG_C_RECOVERY_POLITICS_H
#define WOLONG_C_RECOVERY_POLITICS_H
#include "world_update.h"
int ki_politics_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_politics_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define PN_DECL(name) void name(KiMachine16 *,const KiEconomyHooks *);
PN_DECL(sub_1585F) PN_DECL(sub_15899) PN_DECL(sub_15940) PN_DECL(sub_12AD2)
PN_DECL(sub_15990) PN_DECL(sub_1301C) PN_DECL(sub_12BD9) PN_DECL(sub_12C52)
PN_DECL(sub_12CDF) PN_DECL(sub_12D3A) PN_DECL(sub_12FB1) PN_DECL(sub_12D58)
PN_DECL(sub_12DB8) PN_DECL(sub_12DF3) PN_DECL(sub_130F0) PN_DECL(sub_1310A)
PN_DECL(sub_13091) PN_DECL(sub_12E33) PN_DECL(sub_12E89) PN_DECL(sub_12EFB) PN_DECL(sub_12F71)
#undef PN_DECL
#endif
