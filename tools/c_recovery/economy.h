#ifndef WOLONG_C_RECOVERY_ECONOMY_H
#define WOLONG_C_RECOVERY_ECONOMY_H
#include "rng.h"
typedef struct {
    void (*enter)(KiMachine16 *, uint16_t, void *);
    void (*external)(KiMachine16 *, uint16_t, void *);
    void *user;
} KiEconomyHooks;
void sub_15358(KiMachine16 *, const KiEconomyHooks *);
void sub_15609(KiMachine16 *);
void sub_1563B(KiMachine16 *);
void sub_154FC(KiMachine16 *);
void sub_155EC(KiMachine16 *);
void sub_15828(KiMachine16 *, const KiEconomyHooks *);
void ki_economy_invoke(KiMachine16 *, uint16_t, const KiEconomyHooks *);
#endif
