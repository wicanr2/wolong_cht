#ifndef WOLONG_C_RECT_H
#define WOLONG_C_RECT_H
#include "aligned.h"
int ki_rect_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_rect_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define RC_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
RC_DECL(sub_1F020) RC_DECL(sub_1F140) RC_DECL(sub_1F17B) RC_DECL(sub_1F1A3)
RC_DECL(sub_10AAA) RC_DECL(sub_10AD9) RC_DECL(sub_10CAC) RC_DECL(sub_10CC3)
RC_DECL(sub_1C61F) RC_DECL(sub_1C6BF) RC_DECL(sub_1C6AE) RC_DECL(sub_1C6F6)
RC_DECL(sub_1C74C) RC_DECL(sub_1C775) RC_DECL(sub_1C78E) RC_DECL(sub_10BCD)
#undef RC_DECL
#endif
