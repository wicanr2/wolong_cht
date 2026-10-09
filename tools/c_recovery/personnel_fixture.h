#ifndef WOLONG_C_PERSONNEL_FIXTURE_H
#define WOLONG_C_PERSONNEL_FIXTURE_H
#include "formation_fixture.h"
void wolong_personnel_entry(KiMachine16 *,uint16_t);
typedef KiListFixture KiPersonnelFixture;
static void psn_snapshot(KiMachine16 *m,uint16_t t,void *user) {
 wolong_personnel_entry(m,t);ll_snapshot(m,t,user);
}
static void psn_external(KiMachine16 *m,uint16_t t,void *user) {
 KiEconomyHooks h={psn_snapshot,psn_external,user};
 if(!ki_personnel_body(m,t,&h)) fc_external(m,t,user);
}
static void personnel_run(KiMachine16 *m,uint16_t t,KiPersonnelFixture *f) {
 KiEconomyHooks h={psn_snapshot,psn_external,f};ki_personnel_invoke(m,t,&h);
}
#endif
