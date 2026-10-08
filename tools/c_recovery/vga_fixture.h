#ifndef WOLONG_C_VGA_FIXTURE_H
#define WOLONG_C_VGA_FIXTURE_H
#include "vga.h"
#include <string.h>
typedef struct {uint16_t target,regs[14];uint8_t stack[32],device[28];} KiVgaSnapshot;
typedef struct {unsigned calls;KiVgaSnapshot trace[8192];} KiVgaFixture;
static void vg_snapshot(KiMachine16 *m,uint16_t t,void *user) {
    KiVgaFixture *f=user;unsigned i=f->calls++;if(i>=8192) return;KiVgaSnapshot *s=&f->trace[i];s->target=t;uint16_t r[]={m->ax,m->bx,m->cx,m->dx,m->si,m->di,m->bp,m->sp,m->ds,m->es,m->ss,m->cs,m->ip,m->flags};memcpy(s->regs,r,sizeof(r));for(unsigned j=0;j<32;++j) s->stack[j]=ec_byte(m,m->ss,(uint16_t)(m->sp+j));wolong_vga_state(s->device);
}
static void vga_run(KiMachine16 *m,uint16_t t,KiVgaFixture *f) {KiVgaHooks h={vg_snapshot,f};ki_vga_invoke(m,t,&h);}
#endif
