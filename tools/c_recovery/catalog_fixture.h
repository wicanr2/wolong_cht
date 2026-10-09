#ifndef WOLONG_C_CATALOG_FIXTURE_H
#define WOLONG_C_CATALOG_FIXTURE_H
#include "list_fixture.h"
typedef KiListFixture KiCatalogFixture;
static void cc_external(KiMachine16 *m,uint16_t t,void *user) {
 KiEconomyHooks h={ll_snapshot,cc_external,user};
 if(!ki_catalog_body(m,t,&h)) ll_external(m,t,user);
}
static void catalog_run(KiMachine16 *m,uint16_t t,KiCatalogFixture *f) {
 KiEconomyHooks h={ll_snapshot,cc_external,f};ki_catalog_invoke(m,t,&h);
}
#endif
