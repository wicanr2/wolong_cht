/* Original city list, live sorting fields and relocation; include after verdict.c. */
#include "list.h"
#ifndef KI_LIST_MUTATION
#define KI_LIST_MUTATION 0
#endif
static int ll_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_list_invoke(m,t,h);
 /* sub_18463 removes one return address when accepting a row. */
 return m->ip!=next;
}
static void ll_sort_load(KiMachine16 *m,uint16_t site) {
 uint16_t address=(uint16_t)(m->bx+m->si);
 if(dl_read8(m,m->cs,site)&1) m->ax=dl_read16(m,m->ds,address);
 else vg_al(m,dl_read8(m,m->ds,address));
}
static void ll_sort_compare(KiMachine16 *m) {
 uint16_t address=(uint16_t)(m->bx+m->si);
 if(dl_read8(m,m->cs,0x85ef)&1) ec_cmp(m,m->ax,dl_read16(m,m->ds,address),16);
 else ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,address),8);
}
static int ll_sort_skip(KiMachine16 *m) {
 uint8_t branch=dl_read8(m,m->cs,0x85f1);
 assert(branch==0x73||branch==0x76);
 return branch==0x73?!(m->flags&1):(m->flags&(1|0x40))!=0;
}
#include "list_generated.inc"
void ki_list_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 if(ll_known(t)) {
  if(h&&h->enter) h->enter(m,t,h->user);
  int escaped=0;int handled=ki_list_body(m,t,h,&escaped);assert(handled);(void)handled;
  if(!escaped) m->ip=ec_pop(m);
 } else if(t==0x33fd) ki_events_invoke(m,t,h);
 else ki_verdict_invoke(m,t,h);
}

#define LL_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {int escaped=0;ki_list_body(m,t,h,&escaped);if(!escaped)m->ip=ec_pop(m);}
LL_WRAP(sub_15E60,0x5e60)
LL_WRAP(sub_16909,0x6909)
LL_WRAP(sub_17400,0x7400)
LL_WRAP(sub_1748F,0x748f)
LL_WRAP(sub_181C0,0x81c0)
LL_WRAP(sub_1820E,0x820e)
LL_WRAP(sub_18412,0x8412)
LL_WRAP(nullsub_2,0x8458)
LL_WRAP(sub_18463,0x8463)
LL_WRAP(sub_184BC,0x84bc)
LL_WRAP(sub_184DD,0x84dd)
LL_WRAP(sub_1851A,0x851a)
LL_WRAP(sub_18546,0x8546)
LL_WRAP(sub_18607,0x8607)
LL_WRAP(sub_18662,0x8662)
LL_WRAP(sub_186C9,0x86c9)
LL_WRAP(sub_18713,0x8713)
LL_WRAP(sub_18755,0x8755)
#undef LL_WRAP
