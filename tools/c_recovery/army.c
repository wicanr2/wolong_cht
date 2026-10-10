/* 原軍團輪轉／行軍靜態C圖；re/132、spec/252。 */
#include "army.h"
#ifndef KI_ARMY_MUTATION
#define KI_ARMY_MUTATION 0
#endif
#ifndef KI_ARMY_SOURCE_DIGEST
#define KI_ARMY_SOURCE_DIGEST 0ULL
#endif
uint64_t ki_army_compiled_source(void) { return (uint64_t)KI_ARMY_SOURCE_DIGEST; }
unsigned ki_army_mutation(void) { return (unsigned)KI_ARMY_MUTATION; }
#include "army_generated.inc"
