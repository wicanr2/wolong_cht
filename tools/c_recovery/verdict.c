/* Original ruler sortie and six-slot formation; include after resume.c. */
#include "verdict.h"
#ifndef KI_VERDICT_MUTATION
#define KI_VERDICT_MUTATION 0
#endif
static void vd_call(KiMachine16 *m,uint16_t t,uint16_t next,const KiEconomyHooks *h) {
 ec_push(m,next);m->ip=t;ki_verdict_invoke(m,t,h);
}
#include "verdict_generated.inc"
#define VD_WRAP(n,t) void n(KiMachine16 *m,const KiEconomyHooks *h) {ki_verdict_body(m,t,h);m->ip=ec_pop(m);}
VD_WRAP(sub_1461D,0x461d) VD_WRAP(sub_14698,0x4698) VD_WRAP(sub_14717,0x4717)
VD_WRAP(sub_15E80,0x5e80) VD_WRAP(sub_15EB7,0x5eb7) VD_WRAP(sub_15F27,0x5f27)
VD_WRAP(sub_15F5D,0x5f5d) VD_WRAP(sub_15F7F,0x5f7f) VD_WRAP(sub_1699E,0x699e)
VD_WRAP(sub_16E8F,0x6e8f) VD_WRAP(sub_16EC9,0x6ec9) VD_WRAP(sub_16F26,0x6f26)
VD_WRAP(sub_16F86,0x6f86) VD_WRAP(sub_16FD2,0x6fd2)
#undef VD_WRAP
void ki_verdict_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
 switch(t) {
 case 0x461d:case 0x4698:case 0x4717:case 0x5e80:case 0x5eb7:case 0x5f27:case 0x5f5d:
 case 0x5f7f:case 0x699e:case 0x6e8f:case 0x6ec9:case 0x6f26:case 0x6f86:case 0x6fd2:
  if(h&&h->enter) h->enter(m,t,h->user);
  {int handled=ki_verdict_body(m,t,h);assert(handled);(void)handled;}
  m->ip=ec_pop(m);break;
 case 0x55ec:ki_economy_invoke(m,t,h);break;
 default:ki_resume_invoke(m,t,h);break;
 }
}
