/* 殘部群 10 函式靜態C圖；re/134、spec/254。 */
#include "residual.h"
#ifndef KI_RESIDUAL_MUTATION
#define KI_RESIDUAL_MUTATION 0
#endif
#ifndef KI_RESIDUAL_SOURCE_DIGEST
#define KI_RESIDUAL_SOURCE_DIGEST 0ULL
#endif
uint64_t ki_residual_compiled_source(void) { return (uint64_t)KI_RESIDUAL_SOURCE_DIGEST; }
unsigned ki_residual_mutation(void) { return (unsigned)KI_RESIDUAL_MUTATION; }
static void rs_interrupt(KiMachine16 *m,uint8_t number,uint16_t next,const KiEconomyHooks *h) {
    uint8_t ah=(uint8_t)(m->ax>>8),al=(uint8_t)m->ax;
    m->ip=next;
    assert(wolong_glyph_interrupt(m,number));
    if(number==0x21&&ah==0x4C) wolong_residual_exit(m,al,h?h->user:0);
}
#include "residual_generated.inc"
