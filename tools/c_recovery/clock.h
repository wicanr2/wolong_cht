#ifndef WOLONG_C_RECOVERY_CLOCK_H
#define WOLONG_C_RECOVERY_CLOCK_H
#include "rng.h"

typedef struct {
    void (*call)(KiMachine16 *, uint16_t, void *);
    void (*poll)(KiMachine16 *, uint16_t, void *);
    void *user;
} KiClockHooks;
void sub_11D8E(KiMachine16 *machine, const KiClockHooks *hooks);
#endif
