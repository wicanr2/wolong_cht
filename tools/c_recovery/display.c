/* Original display code closure; spec/227. Include after rect.c.
 * Generated bodies contain concrete C calculations and labels, never opcode fetch.
 */
#include "display.h"
#include <assert.h>
#ifndef KI_DISPLAY_MUTATION
#define KI_DISPLAY_MUTATION 0
#endif
static uint8_t dl_read8(KiMachine16 *m,uint16_t seg,uint16_t off) {return al_read(m,seg,off);}
static uint16_t dl_read16(KiMachine16 *m,uint16_t seg,uint16_t off) {return al_word(m,seg,off);}
static void dl_write8(KiMachine16 *m,uint16_t seg,uint16_t off,uint16_t v) {(void)m;wolong_vga_write(ec_at(seg,off),(uint8_t)v);}
static void dl_write16(KiMachine16 *m,uint16_t seg,uint16_t off,uint16_t v) {dl_write8(m,seg,off,(uint8_t)v);dl_write8(m,seg,(uint16_t)(off+1),v>>8);}
static uint16_t dl_shift(KiMachine16 *m,uint16_t value,unsigned count,unsigned width,unsigned op) {
    uint16_t mask=width==8?255:65535,sign=width==8?128:32768;value&=mask;
    while(count--) {
        if(op==5) value=st_shr(m,value,width);
        else if(op==4) value=width==8?vg_shl8(m,(uint8_t)value):ec_shl(m,value);
        else {uint16_t carry,result,overflow;if(op==0) {carry=(value&sign)!=0;result=(uint16_t)(((value<<1)&mask)|carry);overflow=((result&sign)!=0)^carry;}else {carry=value&1;result=(uint16_t)((value>>1)|(carry?sign:0));overflow=((result&sign)!=0)^((result&(sign>>1))!=0);}m->flags=(uint16_t)((m->flags&~0x801u)|carry|(overflow<<11));value=result;}
    }return value;
}
static void dl_mul(KiMachine16 *m,uint16_t v,unsigned width) {
    uint32_t result=width==8?(uint32_t)(uint8_t)m->ax*(uint8_t)v:(uint32_t)m->ax*v;
    m->ax=(uint16_t)result;if(width==16) m->dx=(uint16_t)(result>>16);
    unsigned high=width==8?result>>8:result>>16;uint16_t bits=high?0x801:0;
    m->flags=ec_szp((uint16_t)((m->flags&~0x801u)|bits),width==8?(uint16_t)result:(uint16_t)(result>>16),16);
    if(result==0) m->flags|=0x40;else m->flags&=(uint16_t)~0x40u;
}
static void dl_div(KiMachine16 *m,uint16_t v,unsigned width) {
    uint32_t n=width==8?m->ax:((uint32_t)m->dx<<16)|m->ax;assert(v!=0);
    ec_cmp(m,width==8?m->ax>>8:m->dx,v,width);uint32_t q=n/v,r=n%v;assert(q<=(width==8?255:65535));
    if(width==8) m->ax=(uint16_t)(q|(r<<8));else {m->ax=(uint16_t)q;m->dx=(uint16_t)r;}
}
static void dl_out(KiMachine16 *m,uint16_t port,uint16_t value,unsigned width) {(void)m;wolong_vga_out(port,(uint8_t)value);if(width==16) wolong_vga_out((uint16_t)(port+1),(uint8_t)(value>>8));}
static void dl_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {ec_push(m,next);m->ip=t;ki_display_invoke(m,t,h);}
#include "display_generated.inc"
void ki_display_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
    if(h&&h->enter) h->enter(m,t,h->user);KiVgaHooks vh={h?h->enter:0,h?h->user:0};
    if(!ki_display_body(m,t,h)&&!ki_rect_body(m,t,h)&&!ki_aligned_body(m,t,h)&&!ki_vga_body(m,t,&vh)&&!ki_hotspot_body(m,t,h)&&!ki_numeric_body(m,t,h)&&h&&h->external) h->external(m,t,h->user);
    m->ip=ec_pop(m);
}
#define DL_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_display_body(m,t,h);m->ip=ec_pop(m);}
DL_WRAP(sub_1030F,0x30f) DL_WRAP(sub_10337,0x337) DL_WRAP(sub_1E993,0xe993) DL_WRAP(sub_1E9A7,0xe9a7) DL_WRAP(sub_1E9C1,0xe9c1)
DL_WRAP(sub_1FB29,0xfb29) DL_WRAP(sub_1FBA7,0xfba7) DL_WRAP(sub_1F6DC,0xf6dc) DL_WRAP(sub_1F878,0xf878)
#undef DL_WRAP
