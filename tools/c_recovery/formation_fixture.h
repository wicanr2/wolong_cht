#ifndef WOLONG_C_FORMATION_FIXTURE_H
#define WOLONG_C_FORMATION_FIXTURE_H
#include "catalog_fixture.h"
typedef KiListFixture KiFormationFixture;
static void fc_external(KiMachine16 *m,uint16_t t,void *user) {
 KiEconomyHooks h={ll_snapshot,fc_external,user};
 if(!ki_formation_body(m,t,&h)) cc_external(m,t,user);
}
static void formation_run(KiMachine16 *m,uint16_t t,KiFormationFixture *f) {
 KiEconomyHooks h={ll_snapshot,fc_external,f};ki_formation_invoke(m,t,&h);
}
#endif
