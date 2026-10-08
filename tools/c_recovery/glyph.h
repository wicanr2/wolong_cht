#ifndef WOLONG_C_GLYPH_H
#define WOLONG_C_GLYPH_H
#include "display.h"
int ki_glyph_body(KiMachine16 *,uint16_t,const KiEconomyHooks *);
void ki_glyph_invoke(KiMachine16 *,uint16_t,const KiEconomyHooks *);
int wolong_glyph_interrupt(KiMachine16 *,uint8_t);
void sub_1F720(KiMachine16 *,const KiEconomyHooks *);
void sub_1F7A4(KiMachine16 *,const KiEconomyHooks *);
void sub_106F5(KiMachine16 *,const KiEconomyHooks *);
void sub_106FD(KiMachine16 *,const KiEconomyHooks *);
#endif
