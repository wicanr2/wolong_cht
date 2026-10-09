/* 原交戰／戰術完整靜態 C 圖；re/131、spec/251。 */
#include "engagement.h"
#ifndef KI_ENGAGEMENT_MUTATION
#define KI_ENGAGEMENT_MUTATION 0
#endif
static uint16_t eg_sar1(KiMachine16 *m,uint16_t v,unsigned width) {
 uint16_t sign=width==8?0x80:0x8000,mask=width==8?0xff:0xffff;
 v&=mask;uint16_t r=(uint16_t)((v>>1)|(v&sign));
 m->flags=ec_szp((uint16_t)((m->flags&~0x8d5u)|(v&1)),r,width);return r;
}
static void eg_imul(KiMachine16 *m,uint16_t value,unsigned width) {
 int32_t product=width==8?(int32_t)(int8_t)m->ax*(int8_t)value:(int32_t)(int16_t)m->ax*(int16_t)value;
 int overflow=width==8?product!=(int8_t)product:product!=(int16_t)product;
 m->ax=(uint16_t)product;if(width==16)m->dx=(uint16_t)((uint32_t)product>>16);
 m->flags=ec_szp((uint16_t)((m->flags&~0x801u)|(overflow?0x801:0)),width==8?(uint16_t)product:m->dx,16);
}
static void eg_string8(KiMachine16 *m,unsigned operation,unsigned repeat) {
 if(repeat&&!m->cx)return;
 do {uint8_t value=operation==2?(uint8_t)m->ax:dl_read8(m,m->ds,m->si);
  if(operation==1)vg_al(m,value);else dl_write8(m,m->es,m->di,value);
  int step=m->flags&0x400?-1:1;
  if(operation!=2)m->si=(uint16_t)(m->si+step);
  if(operation!=1)m->di=(uint16_t)(m->di+step);
 }while(repeat&&--m->cx);
}
#include "engagement_generated.inc"
