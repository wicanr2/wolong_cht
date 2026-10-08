#ifndef WOLONG_C_RECOVERY_SEED_H
#define WOLONG_C_RECOVERY_SEED_H
#include "rng.h"

/* BIOS reply bytes, not decimal time. Original DOS/V IDA 0x1EC9D. */
void sub_1EC82(KiRngState *state, uint8_t hour_bcd, uint8_t minute_bcd,
              uint8_t second_bcd, uint8_t rtc_al);
void sub_1EC82_abi(KiMachine16 *machine, uint8_t hour_bcd, uint8_t minute_bcd,
                  uint8_t second_bcd, uint8_t rtc_al);
#endif
