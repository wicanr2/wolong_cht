/* 原開機、主排程與命名UI靜態C圖；re/133、spec/253。 */
#include "startup.h"
#ifndef KI_STARTUP_MUTATION
#define KI_STARTUP_MUTATION 0
#endif
#ifndef KI_STARTUP_SOURCE_DIGEST
#define KI_STARTUP_SOURCE_DIGEST 0ULL
#endif
uint64_t ki_startup_compiled_source(void) { return (uint64_t)KI_STARTUP_SOURCE_DIGEST; }
unsigned ki_startup_mutation(void) { return (unsigned)KI_STARTUP_MUTATION; }
static void su_interrupt(KiMachine16 *m,uint8_t number,uint16_t next,const KiEconomyHooks *h) {
    uint8_t ah=(uint8_t)(m->ax>>8),al=(uint8_t)m->ax;
    m->ip=next;
    assert(wolong_glyph_interrupt(m,number));
    if(number==0x21&&ah==0x4C) wolong_startup_exit(m,al,h?h->user:0);
}
#include "startup_generated.inc"
