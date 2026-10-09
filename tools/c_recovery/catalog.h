#ifndef WOLONG_C_CATALOG_H
#define WOLONG_C_CATALOG_H
#include "list.h"
int ki_catalog_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_catalog_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define CC_DECL(n) void n(KiMachine16 *,const KiEconomyHooks *);
CC_DECL(sub_1716D)
CC_DECL(sub_171D3)
CC_DECL(sub_175FA)
CC_DECL(sub_17663)
CC_DECL(sub_1770C)
CC_DECL(sub_178A7)
CC_DECL(sub_17906)
CC_DECL(sub_1799C)
CC_DECL(sub_17A7A)
CC_DECL(sub_17B3C)
CC_DECL(sub_17BC0)
#undef CC_DECL
#endif
