#ifndef WOLONG_C_DETAILS_FIXTURE_H
#define WOLONG_C_DETAILS_FIXTURE_H
#include "list_fixture.h"
typedef KiListFixture KiDetailsFixture;
static void dtl_external(KiMachine16 *m,uint16_t t,void *user) {
 KiEconomyHooks h={ll_snapshot,dtl_external,user};
 if(!ki_details_body(m,t,&h)&&!ki_personnel_body(m,t,&h)&&!ki_formation_body(m,t,&h)&&!ki_catalog_body(m,t,&h)) ll_external(m,t,user);
}
static void details_run(KiMachine16 *m,uint16_t t,KiDetailsFixture *f) {
 KiEconomyHooks h={ll_snapshot,dtl_external,f};ki_details_invoke(m,t,&h);
}
#endif
