#ifndef WOLONG_C_MODAL_H
#define WOLONG_C_MODAL_H
#include "events.h"
int ki_modal_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_modal_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
#define MD_DECL(name) void name(KiMachine16 *,const KiEconomyHooks *);
MD_DECL(sub_12078) MD_DECL(sub_120D6) MD_DECL(sub_138C7) MD_DECL(sub_138E6)
MD_DECL(sub_13902) MD_DECL(sub_139E8) MD_DECL(sub_13C3D) MD_DECL(sub_13B7E)
MD_DECL(sub_13D09) MD_DECL(sub_13D45) MD_DECL(sub_13C99) MD_DECL(sub_13CDC)
MD_DECL(sub_19321) MD_DECL(sub_187FF) MD_DECL(sub_13D68) MD_DECL(sub_11D46)
MD_DECL(sub_12216)
#undef MD_DECL
#endif
