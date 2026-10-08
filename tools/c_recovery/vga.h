#ifndef WOLONG_C_VGA_H
#define WOLONG_C_VGA_H
#include "settlement.h"
/* External platform functions use a separate mature VGA state, never guest code. */
uint8_t wolong_vga_read(uint32_t);
void wolong_vga_write(uint32_t,uint8_t);
void wolong_vga_out(uint16_t,uint8_t);
void wolong_vga_state(uint8_t *);
typedef void (*KiVgaEnter)(KiMachine16 *,uint16_t,void *);
typedef struct {KiVgaEnter enter;void *user;} KiVgaHooks;
int ki_vga_body(KiMachine16 *,uint16_t,const KiVgaHooks *);
void ki_vga_invoke(KiMachine16 *,uint16_t,const KiVgaHooks *);
#define VG_DECL(n) void n(KiMachine16 *,const KiVgaHooks *);
VG_DECL(sub_19796) VG_DECL(sub_197C3) VG_DECL(sub_1F9B0) VG_DECL(sub_1FA1B)
VG_DECL(sub_1FA37) VG_DECL(sub_1FAA2) VG_DECL(sub_1FAC2) VG_DECL(sub_1FB11)
#undef VG_DECL
#endif
